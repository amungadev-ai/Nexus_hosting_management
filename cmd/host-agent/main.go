// Command nexus-host-agent is the v1 host agent: a root process reachable only
// through a local Unix socket with SO_PEERCRED peer verification, serving a
// typed, closed allowlist of built-in commands (ADR-R0-003). There is no
// shell, no passthrough exec and no network listener.
package main

import (
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"log"
	"net"
	"os"
	"path/filepath"
	"time"

	"golang.org/x/sys/unix"

	"github.com/amungadev-ai/Nexus_hosting_management/internal/agentproto"
	"github.com/amungadev-ai/Nexus_hosting_management/internal/capability"
)

const version = "0.1.0-skeleton"

var startTime = time.Now()

func envOr(key, def string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return def
}

func main() {
	socket := flag.String("socket", envOr("NEXUS_AGENT_SOCKET", "/run/nexuspanel/agent.sock"), "Unix socket path")
	flag.Parse()

	if err := os.MkdirAll(filepath.Dir(*socket), 0o755); err != nil {
		log.Fatalf("create socket dir: %v", err)
	}
	if _, err := os.Stat(*socket); err == nil {
		log.Printf("removing stale socket %s", *socket)
		if err := os.Remove(*socket); err != nil {
			log.Fatalf("remove stale socket: %v", err)
		}
	}
	ln, err := net.Listen("unix", *socket)
	if err != nil {
		log.Fatalf("listen: %v", err)
	}
	if err := os.Chmod(*socket, 0o660); err != nil {
		log.Fatalf("chmod socket: %v", err)
	}
	log.Printf("host agent %s listening on %s (allowlist: ping, status, detect)", version, *socket)

	for {
		conn, err := ln.Accept()
		if err != nil {
			log.Printf("accept: %v", err)
			continue
		}
		go serve(conn)
	}
}

func serve(conn net.Conn) {
	defer conn.Close()
	_ = conn.SetDeadline(time.Now().Add(30 * time.Second))

	uid, err := peerUID(conn)
	if err != nil {
		log.Printf("rejecting connection without peer credentials: %v", err)
		return
	}

	req, err := agentproto.DecodeRequest(conn)
	if err != nil {
		log.Printf("uid=%d bad request: %v", uid, err)
		return
	}
	log.Printf("uid=%d command=%s id=%s", uid, req.Command, req.ID)

	resp := dispatch(req, uid)
	if err := agentproto.EncodeResponse(conn, resp); err != nil {
		log.Printf("uid=%d write response: %v", uid, err)
	}
}

func dispatch(req agentproto.Request, uid uint32) agentproto.Response {
	// The typed allowlist is the entire boundary (ADR-R0-003).
	if !agentproto.Allowed(req.Command) {
		return errorResp(req.ID, "command not in allowlist")
	}
	switch req.Command {
	case agentproto.CmdPing:
		return okResp(req.ID, mustJSON(map[string]any{"pong": true, "uid": uid}))
	case agentproto.CmdStatus:
		return okResp(req.ID, mustJSON(map[string]any{
			"agent_pid": os.Getpid(),
			"version":   version,
			"uptime_s":  int(time.Since(startTime).Seconds()),
		}))
	case agentproto.CmdDetect:
		return okResp(req.ID, mustJSON(capability.Detect()))
	default:
		// Unreachable while the allowlist and this switch agree; kept as a
		// fail-closed backstop so the two can never drift apart silently.
		return errorResp(req.ID, "command not in allowlist")
	}
}

func okResp(id string, result json.RawMessage) agentproto.Response {
	return agentproto.Response{ID: id, OK: true, Result: result}
}

func errorResp(id, msg string) agentproto.Response {
	return agentproto.Response{ID: id, OK: false, Error: msg}
}

func mustJSON(v any) json.RawMessage {
	b, err := json.Marshal(v)
	if err != nil {
		return json.RawMessage(fmt.Sprintf(`{"error":%q}`, err.Error()))
	}
	return b
}

// peerUID returns the effective uid of the process at the other end of the
// socket via SO_PEERCRED. Connections that cannot present credentials are
// rejected — the panel never serves an unauthenticated peer.
func peerUID(conn net.Conn) (uint32, error) {
	uc, ok := conn.(*net.UnixConn)
	if !ok {
		return 0, errors.New("not a unix socket connection")
	}
	raw, err := uc.SyscallConn()
	if err != nil {
		return 0, err
	}
	var cred *unix.Ucred
	var ctlErr error
	if ctlErr = raw.Control(func(fd uintptr) {
		cred, ctlErr = unix.GetsockoptUcred(int(fd), unix.SOL_SOCKET, unix.SO_PEERCRED)
	}); ctlErr != nil {
		return 0, ctlErr
	}
	if cred == nil {
		return 0, errors.New("no peer credentials")
	}
	return cred.Uid, nil
}
