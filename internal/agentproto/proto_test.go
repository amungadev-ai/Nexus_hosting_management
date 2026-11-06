package agentproto

import (
        "bytes"
        "encoding/binary"
        "encoding/json"
        "strings"
        "testing"
)

func TestRequestRoundTrip(t *testing.T) {
        var buf bytes.Buffer
        req := Request{ID: "abc123", Command: CmdDetect, Args: []string{"--json"}}
        if err := EncodeRequest(&buf, req); err != nil {
                t.Fatal(err)
        }
        got, err := DecodeRequest(&buf)
        if err != nil {
                t.Fatal(err)
        }
        if got.ID != req.ID || got.Command != req.Command || len(got.Args) != 1 {
                t.Fatalf("round trip mismatch: %+v", got)
        }
}

func TestResponseRoundTrip(t *testing.T) {
        var buf bytes.Buffer
        resp := Response{ID: "abc123", OK: true, Result: json.RawMessage(`{"pong":true}`)}
        if err := EncodeResponse(&buf, resp); err != nil {
                t.Fatal(err)
        }
        got, err := DecodeResponse(&buf)
        if err != nil {
                t.Fatal(err)
        }
        if got.ID != resp.ID || !got.OK || string(got.Result) != `{"pong":true}` {
                t.Fatalf("round trip mismatch: %+v", got)
        }
}

func TestDecodeRejectsOversizedFrame(t *testing.T) {
        huge := make([]byte, MaxRequestSize+1)
        var buf bytes.Buffer
        var lenbuf [4]byte
        binary.BigEndian.PutUint32(lenbuf[:], uint32(len(huge)))
        buf.Write(lenbuf[:])
        if _, err := DecodeRequest(&buf); err == nil {
                t.Fatal("oversized frame must be rejected")
        } else if !strings.Contains(err.Error(), "exceeds limit") {
                t.Fatalf("wrong error: %v", err)
        }
}

func TestDecodeRejectsEmptyCommand(t *testing.T) {
        var buf bytes.Buffer
        if err := EncodeRequest(&buf, Request{ID: "x"}); err != nil {
                t.Fatal(err)
        }
        if _, err := DecodeRequest(&buf); err == nil {
                t.Fatal("empty command must be rejected")
        }
}

func TestDecodeRejectsBadJSON(t *testing.T) {
        var buf bytes.Buffer
        buf.Write([]byte{0, 0, 0, 4})
        buf.WriteString("{!?!")
        if _, err := DecodeRequest(&buf); err == nil {
                t.Fatal("bad JSON must be rejected")
        }
}

func TestEncodeRejectsOversizedBody(t *testing.T) {
        req := Request{ID: strings.Repeat("a", MaxRequestSize)}
        if err := EncodeRequest(&bytes.Buffer{}, req); err == nil {
                t.Fatal("encoder must reject oversized bodies")
        }
}

func TestAllowlistComplete(t *testing.T) {
        for _, cmd := range []string{CmdPing, CmdStatus, CmdDetect} {
                if !Allowed(cmd) {
                        t.Fatalf("command %q must be allowed", cmd)
                }
        }
        for _, cmd := range []string{"", "exec", "sh", "docker", "rm -rf /", "../escape"} {
                if Allowed(cmd) {
                        t.Fatalf("command %q must NOT be allowed", cmd)
                }
        }
}

// FuzzDecodeRequest guards the agent boundary (ADR-R0-003): the decoder must
// never panic on hostile input, and must enforce the frame limit.
func FuzzDecodeRequest(f *testing.F) {
        f.Add([]byte{0, 0, 0, 2}, []byte(`{}`))
        f.Add([]byte{0, 0, 255, 255}, []byte(`{"id":"x","command":"ping"}`))
        f.Add([]byte{0, 0, 0, 9}, []byte(`{"id":"x"`))
        f.Fuzz(func(t *testing.T, lenBytes, body []byte) {
                var buf bytes.Buffer
                buf.Write(lenBytes)
                buf.Write(body)
                req, err := DecodeRequest(&buf)
                if err == nil && req.Command == "" {
                        t.Fatal("decoded request with empty command must error")
                }
        })
}
