package capability

import (
	"encoding/json"
	"reflect"
	"strings"
	"testing"
)

const sampleOSRelease = `NAME="Ubuntu"
ID=ubuntu
VERSION_ID="24.04"
PRETTY_NAME="Ubuntu 24.04 LTS"
# a comment

HOME_URL="https://www.ubuntu.com/"
`

func TestParseOSRelease(t *testing.T) {
	m := parseOSRelease(strings.NewReader(sampleOSRelease))
	if m["ID"] != "ubuntu" || m["VERSION_ID"] != "24.04" || m["PRETTY_NAME"] != "Ubuntu 24.04 LTS" {
		t.Fatalf("unexpected parse result: %v", m)
	}
	if _, ok := m["HOME_URL"]; !ok {
		t.Fatal("non-quoted values must parse too")
	}
}

func TestSameOSMinor(t *testing.T) {
	cases := []struct {
		versionID, want string
		expect          bool
	}{
		{"24.04", "24.04", true},
		{"24.04.1", "24.04", true},
		{"24.10", "24.04", false},
		{"26.04", "26.04", true},
		{"2", "24.04", false},
		{"", "24.04", false},
	}
	for _, c := range cases {
		if got := sameOSMinor(c.versionID, c.want); got != c.expect {
			t.Errorf("sameOSMinor(%q, %q) = %v, want %v", c.versionID, c.want, got, c.expect)
		}
	}
}

func TestParseMemInfoKB(t *testing.T) {
	sample := `MemTotal:       6062616 kB
MemFree:          987680 kB
MemAvailable:    5610512 kB
Buffers:          104857 kB
`
	total, avail := parseMemInfoKB(strings.NewReader(sample))
	if total != 6062616 || avail != 5610512 {
		t.Fatalf("got total=%d avail=%d", total, avail)
	}
}

func TestParseProcNetTCP(t *testing.T) {
	sample := `  sl local_address rem_address   st tx_queue rx_queue tr tm->when retrnsmt   uid timeout inode
   0: 0100007F:0016 00000000:0000 0A 00000000:00000000 00:00000000 00000000     0        0 12345 1 ...
   1: 00000000:1F90 00000000:0000 0A 00000000:00000000 00:00000000 00000000     0        0 23456 1 ...
   2: 0100007F:8AE2 0100007F:1F90 01 00000000:00000000 00:00000000 00000000  1000        0 34567 1 ...
`
	ports := parseProcNetTCP([]byte(sample))
	want := map[int]bool{22: true, 8080: true}
	if !reflect.DeepEqual(ports, want) {
		t.Fatalf("got %v, want %v", ports, want)
	}
}

func TestEvaluateGates(t *testing.T) {
	base := Report{
		OSID:       "ubuntu",
		OSVersion:  "24.04",
		MemTotalMB: 5920,
		DiskFreeGB: 93,
		CgroupV2:   true,
	}

	r := Evaluate(base)
	if len(r.Errors) != 0 {
		t.Fatalf("healthy host must pass, got errors: %v", r.Errors)
	}
	if !contains(r.Warnings, "Docker socket not present — static-site provisioning only (ADR-005c commitment gate)") {
		t.Fatalf("missing docker must warn: %v", r.Warnings)
	}

	noCgroup := base
	noCgroup.CgroupV2 = false
	r = Evaluate(noCgroup)
	if !contains(r.Errors, "cgroup v2 not active on /sys/fs/cgroup (DEC-015 requires the unified hierarchy)") {
		t.Fatalf("cgroup v2 gate must fail closed: %v", r.Errors)
	}

	nano := base
	nano.MemTotalMB = 511
	r = Evaluate(nano)
	if len(r.Errors) != 1 || !strings.Contains(r.Errors[0], "ADR-R0-001") {
		t.Fatalf("RAM floor must reference ADR-R0-001: %v", r.Errors)
	}

	experimental := base
	experimental.OSVersion = "26.04"
	r = Evaluate(experimental)
	if len(r.Errors) != 0 || !contains(r.Warnings, "Ubuntu 26.04 is experimental until profile tests pass (ADR-R0-002)") {
		t.Fatalf("26.04 must warn, not error: %v / %v", r.Errors, r.Warnings)
	}

	unsupported := base
	unsupported.OSID = "debian"
	r = Evaluate(unsupported)
	if len(r.Errors) != 1 || !strings.Contains(r.Errors[0], "ADR-R0-002") {
		t.Fatalf("unsupported OS must fail closed: %v", r.Errors)
	}
}

func contains(list []string, s string) bool {
	for _, v := range list {
		if v == s {
			return true
		}
	}
	return false
}

func TestDetectRunsOnLinux(t *testing.T) {
	r := Evaluate(Detect())
	b, err := json.Marshal(r)
	if err != nil {
		t.Fatalf("report must be JSON-serializable: %v", err)
	}
	if len(b) == 0 {
		t.Fatal("empty report")
	}
}
