// Package capability detects the host facts the panel gates on: cgroup v2
// (DEC-015), supported OS (ADR-R0-002), RAM/disk floors (ADR-R0-001),
// virtualization, Docker presence and which well-known ports are taken.
package capability

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"os"
	"runtime"
	"strconv"
	"strings"
	"syscall"
)

// cgroup2Magic is CGROUP2_SUPER_MAGIC from linux/magic.h.
const cgroup2Magic = 0x63677270

// WellKnownPorts are ports the panel, its sites or well-known services may claim.
var WellKnownPorts = []int{22, 25, 80, 443, 8080, 8081}

// Report is the machine-readable capability snapshot of one host.
type Report struct {
	Hostname       string   `json:"hostname"`
	OSID           string   `json:"os_id"`
	OSVersion      string   `json:"os_version"`
	OSPrettyName   string   `json:"os_pretty_name"`
	Kernel         string   `json:"kernel"`
	Virt           string   `json:"virt"`
	CPUCores       int      `json:"cpu_cores"`
	MemTotalMB     int64    `json:"mem_total_mb"`
	MemAvailableMB int64    `json:"mem_available_mb"`
	DiskFreeGB     int64    `json:"disk_free_gb"`
	CgroupV2       bool     `json:"cgroup_v2"`
	DockerSocket   bool     `json:"docker_socket_present"`
	PortsInUse     []int    `json:"ports_in_use"`
	Warnings       []string `json:"warnings"`
	Errors         []string `json:"errors"`
}

// Detect gathers host facts. It never fails: unreadable facts are left at
// their zero value so Evaluate can flag them.
func Detect() Report {
	r := Report{}
	if h, err := os.Hostname(); err == nil {
		r.Hostname = h
	}
	if data, err := os.ReadFile("/etc/os-release"); err == nil {
		m := parseOSRelease(bytes.NewReader(data))
		r.OSID = m["ID"]
		r.OSVersion = m["VERSION_ID"]
		r.OSPrettyName = m["PRETTY_NAME"]
	}
	if b, err := os.ReadFile("/proc/sys/kernel/osrelease"); err == nil {
		r.Kernel = strings.TrimSpace(string(b))
	}
	r.Virt = detectVirt()
	r.CPUCores = runtime.NumCPU()
	if total, avail, err := memInfoMB(); err == nil {
		r.MemTotalMB, r.MemAvailableMB = total, avail
	}
	if free, err := diskFreeGB("/"); err == nil {
		r.DiskFreeGB = free
	}
	r.CgroupV2 = cgroupV2Active()
	r.DockerSocket = fileExists("/var/run/docker.sock")

	listening := map[int]bool{}
	for _, f := range []string{"/proc/net/tcp", "/proc/net/tcp6"} {
		if data, err := os.ReadFile(f); err == nil {
			for p := range parseProcNetTCP(data) {
				listening[p] = true
			}
		}
	}
	for _, p := range WellKnownPorts {
		if listening[p] {
			r.PortsInUse = append(r.PortsInUse, p)
		}
	}
	return r
}

// Evaluate applies the v1 install gates to a report, fail-closed: any entry in
// Errors blocks installation; Warnings are advisory.
func Evaluate(r Report) Report {
	if !r.CgroupV2 {
		r.Errors = append(r.Errors, "cgroup v2 not active on /sys/fs/cgroup (DEC-015 requires the unified hierarchy)")
	}
	if r.MemTotalMB > 0 && r.MemTotalMB < 512 {
		r.Errors = append(r.Errors, fmt.Sprintf("RAM %d MB is below the Nano floor of 512 MB (ADR-R0-001)", r.MemTotalMB))
	}
	if r.DiskFreeGB > 0 && r.DiskFreeGB < 5 {
		r.Errors = append(r.Errors, fmt.Sprintf("free disk %d GB is below the 5 GB minimum (ADR-R0-001)", r.DiskFreeGB))
	}
	switch {
	case r.OSID == "ubuntu" && sameOSMinor(r.OSVersion, "24.04"):
		// First-class supported target (ADR-R0-002).
	case r.OSID == "ubuntu" && sameOSMinor(r.OSVersion, "26.04"):
		r.Warnings = append(r.Warnings, "Ubuntu 26.04 is experimental until profile tests pass (ADR-R0-002)")
	default:
		r.Errors = append(r.Errors, fmt.Sprintf(
			"unsupported OS %s %s (ADR-R0-002: Ubuntu 24.04 LTS first-class, 26.04 experimental)",
			r.OSID, r.OSVersion))
	}
	if !r.DockerSocket {
		r.Warnings = append(r.Warnings, "Docker socket not present — static-site provisioning only (ADR-005c commitment gate)")
	}
	return r
}

// sameOSMinor reports whether versionID matches want at major.minor precision
// ("24.04.1" matches "24.04"; "24.10" does not).
func sameOSMinor(versionID, want string) bool {
	vp := strings.Split(versionID, ".")
	wp := strings.Split(want, ".")
	for i := range wp {
		if i >= len(vp) || vp[i] != wp[i] {
			return false
		}
	}
	return true
}

func parseOSRelease(rd io.Reader) map[string]string {
	m := map[string]string{}
	sc := bufio.NewScanner(rd)
	for sc.Scan() {
		line := strings.TrimSpace(sc.Text())
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}
		k, v, ok := strings.Cut(line, "=")
		if !ok {
			continue
		}
		v = strings.Trim(v, `"'`)
		m[k] = v
	}
	return m
}

func parseMemInfoKB(rd io.Reader) (total, avail int64) {
	sc := bufio.NewScanner(rd)
	for sc.Scan() {
		key, val, ok := strings.Cut(sc.Text(), ":")
		if !ok {
			continue
		}
		fields := strings.Fields(val)
		if len(fields) == 0 {
			continue
		}
		kb, err := strconv.ParseInt(fields[0], 10, 64)
		if err != nil {
			continue
		}
		switch strings.TrimSpace(key) {
		case "MemTotal":
			total = kb
		case "MemAvailable":
			avail = kb
		}
	}
	return total, avail
}

// parseProcNetTCP returns the set of TCP ports appearing in /proc/net/tcp{,6}
// content. Only rows in LISTEN state (0A) are considered.
func parseProcNetTCP(data []byte) map[int]bool {
	ports := map[int]bool{}
	sc := bufio.NewScanner(bytes.NewReader(data))
	for sc.Scan() {
		fields := strings.Fields(sc.Text())
		if len(fields) < 4 || !strings.HasSuffix(fields[0], ":") {
			continue // header line
		}
		if fields[3] != "0A" { // TCP_LISTEN
			continue
		}
		local := fields[1]
		colon := strings.LastIndex(local, ":")
		if colon < 0 {
			continue
		}
		p, err := strconv.ParseUint(local[colon+1:], 16, 16)
		if err != nil {
			continue
		}
		ports[int(p)] = true
	}
	return ports
}

func memInfoMB() (total, avail int64, err error) {
	f, err := os.Open("/proc/meminfo")
	if err != nil {
		return 0, 0, err
	}
	defer f.Close()
	t, a := parseMemInfoKB(f)
	return t / 1024, a / 1024, nil
}

func diskFreeGB(path string) (int64, error) {
	var st syscall.Statfs_t
	if err := syscall.Statfs(path, &st); err != nil {
		return 0, err
	}
	return int64(st.Bavail) * int64(st.Bsize) / (1 << 30), nil
}

func cgroupV2Active() bool {
	var st syscall.Statfs_t
	if err := syscall.Statfs("/sys/fs/cgroup", &st); err != nil {
		return false
	}
	return uint64(st.Type) == cgroup2Magic
}

func detectVirt() string {
	if fileExists("/.dockerenv") {
		return "docker"
	}
	if fileExists("/run/.containerenv") {
		return "podman"
	}
	if b, err := os.ReadFile("/sys/class/dmi/id/sys_vendor"); err == nil {
		v := strings.ToLower(strings.TrimSpace(string(b)))
		switch {
		case strings.Contains(v, "qemu"):
			return "kvm"
		case strings.Contains(v, "vmware"):
			return "vmware"
		case strings.Contains(v, "microsoft"):
			return "hyper-v"
		case strings.Contains(v, "innotek"):
			return "virtualbox"
		case strings.Contains(v, "xen"):
			return "xen"
		}
	}
	if fileExists("/proc/xen") {
		return "xen"
	}
	return "bare/unknown"
}

func fileExists(path string) bool {
	_, err := os.Stat(path)
	return err == nil
}
