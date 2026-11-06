// Command nexus-npctl is the operator CLI for the walking skeleton:
// local capability detection, agent status over the Unix socket, and
// audit-chain verification through the control-plane API.
package main

import (
        "encoding/json"
        "flag"
        "fmt"
        "io"
        "net"
        "net/http"
        "os"
        "time"

        "github.com/amungadev-ai/Nexus_hosting_management/internal/agentproto"
        "github.com/amungadev-ai/Nexus_hosting_management/internal/capability"
)

func envOr(key, def string) string {
        if v := os.Getenv(key); v != "" {
                return v
        }
        return def
}

func main() {
        if len(os.Args) < 2 {
                usage()
                os.Exit(2)
        }
        var err error
        switch os.Args[1] {
        case "detect":
                err = cmdDetect()
        case "agent":
                err = cmdAgent(os.Args[2:])
        case "audit-verify":
                err = cmdAuditVerify(os.Args[2:])
        case "version":
                fmt.Println("nexus-npctl 0.1.0-skeleton")
        default:
                usage()
                os.Exit(2)
        }
        if err != nil {
                fmt.Fprintln(os.Stderr, "error:", err)
                os.Exit(1)
        }
}

func usage() {
        fmt.Fprint(os.Stderr, `nexus-npctl — NexusPanel operator CLI (walking skeleton)

usage:
  nexus-npctl detect                       detect and gate-check this host
  nexus-npctl agent ping|status [-socket PATH]
                                           query the host agent over its Unix socket
  nexus-npctl audit-verify [-api URL]      verify the audit chain via the control plane
  nexus-npctl version
`)
}

func cmdDetect() error {
        report := capability.Evaluate(capability.Detect())
        out, err := json.MarshalIndent(report, "", "  ")
        if err != nil {
                return err
        }
        fmt.Println(string(out))
        if len(report.Errors) > 0 {
                return fmt.Errorf("%d gate error(s) — installation would fail closed", len(report.Errors))
        }
        return nil
}

func cmdAgent(args []string) error {
        fs := flag.NewFlagSet("agent", flag.ExitOnError)
        socket := fs.String("socket", envOr("NEXUS_AGENT_SOCKET", "/run/nexuspanel/agent.sock"), "agent unix socket")
        if err := fs.Parse(args); err != nil {
                return err
        }
        if fs.NArg() != 1 {
                usage()
                os.Exit(2)
        }
        command := fs.Arg(0)
        switch command {
        case agentproto.CmdPing, agentproto.CmdStatus, agentproto.CmdDetect:
        default:
                return fmt.Errorf("unknown agent command %q (ping|status|detect)", command)
        }

        conn, err := net.DialTimeout("unix", *socket, 5*time.Second)
        if err != nil {
                return fmt.Errorf("dial agent socket %s: %w", *socket, err)
        }
        defer conn.Close()
        _ = conn.SetDeadline(time.Now().Add(10 * time.Second))

        req := agentproto.Request{ID: fmt.Sprintf("npctl-%d", time.Now().UnixNano()), Command: command}
        if err := agentproto.EncodeRequest(conn, req); err != nil {
                return err
        }
        resp, err := agentproto.DecodeResponse(conn)
        if err != nil {
                return err
        }
        if !resp.OK {
                return fmt.Errorf("agent error: %s", resp.Error)
        }
        var pretty any
        if err := json.Unmarshal(resp.Result, &pretty); err == nil {
                out, err := json.MarshalIndent(pretty, "", "  ")
                if err == nil {
                        fmt.Println(string(out))
                        return nil
                }
        }
        fmt.Println(string(resp.Result))
        return nil
}

func cmdAuditVerify(args []string) error {
        fs := flag.NewFlagSet("audit-verify", flag.ExitOnError)
        api := fs.String("api", envOr("NEXUS_API", "http://127.0.0.1:8080"), "control-plane base URL")
        if err := fs.Parse(args); err != nil {
                return err
        }
        client := &http.Client{Timeout: 10 * time.Second}
        resp, err := client.Get(*api + "/api/v1/audit/verify")
        if err != nil {
                return fmt.Errorf("control plane at %s: %w", *api, err)
        }
        defer resp.Body.Close()
        body, err := io.ReadAll(io.LimitReader(resp.Body, 1<<20))
        if err != nil {
                return err
        }
        var result struct {
                Valid bool   `json:"valid"`
                Error string `json:"error"`
        }
        if err := json.Unmarshal(body, &result); err != nil {
                return fmt.Errorf("unexpected response: %s", string(body))
        }
        if result.Valid {
                fmt.Println("audit chain: VALID")
                return nil
        }
        fmt.Printf("audit chain: INVALID — %s\n", result.Error)
        return fmt.Errorf("audit chain invalid: %s", result.Error)
}
