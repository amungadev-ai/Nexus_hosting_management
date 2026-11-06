// Package agentproto defines the host-agent's single-request wire protocol and
// the typed command allowlist (ADR-R0-003). Frames are 4-byte big-endian
// length followed by UTF-8 JSON, bounded by MaxRequestSize / MaxResponseSize
// so a hostile peer cannot exhaust memory through the decoder.
package agentproto

import (
	"encoding/binary"
	"encoding/json"
	"errors"
	"fmt"
	"io"
)

const (
	// MaxRequestSize bounds one request frame (64 KiB is generous for the v1 commands).
	MaxRequestSize = 64 << 10
	// MaxResponseSize bounds one response frame (detect reports stay well under 1 MiB).
	MaxResponseSize = 1 << 20
)

// Commands are the complete typed allowlist. The agent executes nothing that
// is not enumerated here; there is no shell and no passthrough exec.
const (
	CmdPing   = "ping"
	CmdStatus = "status"
	CmdDetect = "detect"
)

var allowlist = map[string]bool{
	CmdPing:   true,
	CmdStatus: true,
	CmdDetect: true,
}

// Allowed reports whether cmd is in the typed allowlist.
func Allowed(cmd string) bool { return allowlist[cmd] }

// Request is one agent command.
type Request struct {
	ID      string   `json:"id"`
	Command string   `json:"command"`
	Args    []string `json:"args,omitempty"`
}

// Response carries the result or an error for Request.ID.
type Response struct {
	ID     string          `json:"id"`
	OK     bool            `json:"ok"`
	Result json.RawMessage `json:"result,omitempty"`
	Error  string          `json:"error,omitempty"`
}

func writeFrame(w io.Writer, body []byte, limit int) error {
	if len(body) > limit {
		return fmt.Errorf("agentproto: frame of %d bytes exceeds limit of %d", len(body), limit)
	}
	var lenbuf [4]byte
	binary.BigEndian.PutUint32(lenbuf[:], uint32(len(body)))
	if _, err := w.Write(lenbuf[:]); err != nil {
		return err
	}
	_, err := w.Write(body)
	return err
}

func readFrame(r io.Reader, limit int) ([]byte, error) {
	var lenbuf [4]byte
	if _, err := io.ReadFull(r, lenbuf[:]); err != nil {
		return nil, fmt.Errorf("agentproto: read length: %w", err)
	}
	n := binary.BigEndian.Uint32(lenbuf[:])
	if n > uint32(limit) {
		return nil, fmt.Errorf("agentproto: frame of %d bytes exceeds limit of %d", n, limit)
	}
	body := make([]byte, n)
	if _, err := io.ReadFull(r, body); err != nil {
		return nil, fmt.Errorf("agentproto: read body: %w", err)
	}
	return body, nil
}

// EncodeRequest writes one request frame.
func EncodeRequest(w io.Writer, req Request) error {
	body, err := json.Marshal(req)
	if err != nil {
		return err
	}
	return writeFrame(w, body, MaxRequestSize)
}

// DecodeRequest reads one request frame. It validates framing and JSON shape
// only; allowlist enforcement is the agent's policy decision, not the codec's.
func DecodeRequest(r io.Reader) (Request, error) {
	var req Request
	body, err := readFrame(r, MaxRequestSize)
	if err != nil {
		return req, err
	}
	if err := json.Unmarshal(body, &req); err != nil {
		return req, fmt.Errorf("agentproto: bad JSON request: %w", err)
	}
	if req.Command == "" {
		return req, errors.New("agentproto: request has empty command")
	}
	return req, nil
}

// EncodeResponse writes one response frame.
func EncodeResponse(w io.Writer, resp Response) error {
	body, err := json.Marshal(resp)
	if err != nil {
		return err
	}
	return writeFrame(w, body, MaxResponseSize)
}

// DecodeResponse reads one response frame.
func DecodeResponse(r io.Reader) (Response, error) {
	var resp Response
	body, err := readFrame(r, MaxResponseSize)
	if err != nil {
		return resp, err
	}
	if err := json.Unmarshal(body, &resp); err != nil {
		return resp, fmt.Errorf("agentproto: bad JSON response: %w", err)
	}
	return resp, nil
}
