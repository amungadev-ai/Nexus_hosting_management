// Command nexus-installer implements the installer's privileged steps:
// fail-closed preflight checks and writing the systemd units + env file.
// Orchestration (downloading, building, systemctl) lives in deploy/install.sh.
package main

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"

	"github.com/amungadev-ai/Nexus_hosting_management/internal/capability"
)

const version = "0.1.0-skeleton"

func main() {
	if len(os.Args) < 2 {
		usage()
		os.Exit(2)
	}
	switch os.Args[1] {
	case "preflight":
		if err := cmdPreflight(os.Args[2:]); err != nil {
			fmt.Fprintln(os.Stderr, "preflight FAILED:", err)
			os.Exit(1)
		}
	case "write-units":
		if err := cmdWriteUnits(os.Args[2:]); err != nil {
			fmt.Fprintln(os.Stderr, "write-units FAILED:", err)
			os.Exit(1)
		}
	case "version":
		fmt.Println("nexus-installer", version)
	default:
		usage()
		os.Exit(2)
	}
}

func usage() {
	fmt.Fprint(os.Stderr, `nexus-installer — NexusPanel installer steps (walking skeleton)

usage:
  nexus-installer preflight [-strict=true]
      Detect host capabilities and apply the v1 install gates
      (DEC-015 cgroup v2, ADR-R0-001 memory floor, ADR-R0-002 OS matrix).

  nexus-installer write-units [-bin-dir /usr/local/bin] [-state-dir /var/lib/nexuspanel]
                             [-socket /run/nexuspanel/agent.sock] [-driver sqlite] [-dsn ""]
      Write /etc/nexuspanel/nexuspanel.env and the two systemd units.
      Run `+"`systemctl daemon-reload`"+` afterwards (install.sh does this).

  nexus-installer version
`)
}

func cmdPreflight(args []string) error {
	fs := flag.NewFlagSet("preflight", flag.ExitOnError)
	strict := fs.Bool("strict", true, "exit non-zero when gate errors are present")
	if err := fs.Parse(args); err != nil {
		return err
	}
	report := capability.Evaluate(capability.Detect())

	fmt.Printf("host            %s\n", report.Hostname)
	fmt.Printf("os              %s %s (%s)\n", report.OSID, report.OSVersion, report.OSPrettyName)
	fmt.Printf("kernel          %s\n", report.Kernel)
	fmt.Printf("virt            %s\n", report.Virt)
	fmt.Printf("cpu/mem/disk    %d vCPU / %d MB total / %d GB free\n", report.CPUCores, report.MemTotalMB, report.DiskFreeGB)
	fmt.Printf("cgroup v2       %v (DEC-015)\n", report.CgroupV2)
	fmt.Printf("docker socket   %v\n", report.DockerSocket)
	fmt.Printf("ports in use    %v\n", report.PortsInUse)
	for _, warning := range report.Warnings {
		fmt.Printf("warning: %s\n", warning)
	}
	for _, err := range report.Errors {
		fmt.Printf("ERROR: %s\n", err)
	}
	if len(report.Errors) > 0 && *strict {
		return fmt.Errorf("%d gate error(s)", len(report.Errors))
	}
	fmt.Println("preflight: PASS")
	return nil
}

func cmdWriteUnits(args []string) error {
	fs := flag.NewFlagSet("write-units", flag.ExitOnError)
	binDir := fs.String("bin-dir", "/usr/local/bin", "directory holding the nexus binaries")
	stateDir := fs.String("state-dir", "/var/lib/nexuspanel", "panel state directory")
	socket := fs.String("socket", "/run/nexuspanel/agent.sock", "agent unix socket")
	driver := fs.String("driver", "sqlite", "storage driver: sqlite | postgres")
	dsn := fs.String("dsn", "", "driver DSN (default: <state-dir>/nexus.db for sqlite)")
	if err := fs.Parse(args); err != nil {
		return err
	}

	resolvedDSN := *dsn
	if resolvedDSN == "" {
		if *driver == "sqlite" {
			resolvedDSN = filepath.Join(*stateDir, "nexus.db")
		} else {
			return fmt.Errorf("-dsn is required for the postgres driver")
		}
	}

	if err := os.MkdirAll("/etc/nexuspanel", 0o700); err != nil {
		return err
	}
	envPath := "/etc/nexuspanel/nexuspanel.env"
	env := fmt.Sprintf(`# Written by nexus-installer %s
NEXUS_DB_DRIVER=%s
NEXUS_DB_DSN=%s
NEXUS_LISTEN_ADDR=127.0.0.1:8080
NEXUS_SITE_ADDR=127.0.0.1:8081
NEXUS_DATA_DIR=%s
NEXUS_AGENT_SOCKET=%s
`, version, *driver, resolvedDSN, *stateDir, *socket)
	if err := os.WriteFile(envPath, []byte(env), 0o600); err != nil {
		return err
	}
	fmt.Println("wrote", envPath)

	controlUnit := fmt.Sprintf(`[Unit]
Description=NexusPanel control plane
After=network-online.target
Wants=network-online.target

[Service]
ExecStart=%s/nexus-control-plane
EnvironmentFile=/etc/nexuspanel/nexuspanel.env
Restart=on-failure
RestartSec=2

[Install]
WantedBy=multi-user.target
`, *binDir)
	if err := writeUnit("/etc/systemd/system/nexuspanel-control-plane.service", controlUnit); err != nil {
		return err
	}

	agentUnit := fmt.Sprintf(`[Unit]
Description=NexusPanel host agent
After=network.target

[Service]
ExecStart=%s/nexus-host-agent -socket %s
RuntimeDirectory=nexuspanel
RuntimeDirectoryMode=0755
Restart=on-failure
RestartSec=2

[Install]
WantedBy=multi-user.target
`, *binDir, *socket)
	if err := writeUnit("/etc/systemd/system/nexuspanel-host-agent.service", agentUnit); err != nil {
		return err
	}

	if err := os.MkdirAll(filepath.Join(*stateDir, "sites"), 0o750); err != nil {
		return err
	}
	fmt.Println("wrote /etc/systemd/system/nexuspanel-control-plane.service")
	fmt.Println("wrote /etc/systemd/system/nexuspanel-host-agent.service")
	fmt.Println("state dir ready:", *stateDir)
	return nil
}

func writeUnit(path, content string) error {
	return os.WriteFile(path, []byte(content), 0o644)
}
