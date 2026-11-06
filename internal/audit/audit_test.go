package audit

import (
	"strings"
	"testing"
	"time"
)

func chainOf(t *testing.T, n int) []Record {
	t.Helper()
	recs := make([]Record, 0, n)
	prev := GenesisPrevHash
	for i := 1; i <= n; i++ {
		r := Record{
			Seq:      int64(i),
			Time:     time.Date(2026, 11, 7, 12, 0, i, 123456789, time.UTC),
			Actor:    "owner",
			Action:   "test.action",
			Payload:  `{"n":` + itoa(i) + `}`,
			PrevHash: prev,
		}
		r.Hash = ComputeHash(r)
		recs = append(recs, r)
		prev = r.Hash
	}
	return recs
}

func itoa(i int) string {
	if i == 0 {
		return "0"
	}
	digits := ""
	for i > 0 {
		digits = string(rune('0'+i%10)) + digits
		i /= 10
	}
	return digits
}

func TestVerifyChainFullValid(t *testing.T) {
	recs := chainOf(t, 5)
	if err := VerifyChainFull(recs); err != nil {
		t.Fatalf("expected valid chain, got %v", err)
	}
	if err := VerifyChainWindow(recs); err != nil {
		t.Fatalf("expected valid window, got %v", err)
	}
}

func TestVerifyChainFullEmpty(t *testing.T) {
	if err := VerifyChainFull(nil); err != nil {
		t.Fatalf("empty chain should verify, got %v", err)
	}
}

func TestVerifyChainFullDetectsPayloadTamper(t *testing.T) {
	recs := chainOf(t, 4)
	recs[1].Payload = `{"n":999}`
	err := VerifyChainFull(recs)
	if err == nil {
		t.Fatal("expected tamper detection")
	}
	if !strings.Contains(err.Error(), "seq 2") {
		t.Fatalf("error should name the tampered seq, got: %v", err)
	}
}

func TestVerifyChainFullDetectsHashTamper(t *testing.T) {
	recs := chainOf(t, 3)
	recs[2].Hash = strings.Repeat("ab", 32)
	if err := VerifyChainFull(recs); err == nil {
		t.Fatal("expected hash tamper detection")
	}
}

func TestVerifyChainFullDetectsPrevHashBreak(t *testing.T) {
	recs := chainOf(t, 3)
	recs[1].PrevHash = GenesisPrevHash // rewrote record 2, re-hashed it, but did not re-link
	recs[1].Hash = ComputeHash(recs[1])
	if err := VerifyChainFull(recs); err == nil {
		t.Fatal("expected prev_hash break detection")
	}
}

func TestVerifyChainFullDetectsSeqGap(t *testing.T) {
	recs := chainOf(t, 3)
	recs[2].Seq = 4
	recs[2].Hash = ComputeHash(recs[2])
	if err := VerifyChainFull(recs); err == nil {
		t.Fatal("expected seq gap detection")
	}
}

func TestVerifyChainWindowDetectsBrokenLink(t *testing.T) {
	recs := chainOf(t, 5)
	recs[3].PrevHash = "tampered"
	recs[3].Hash = ComputeHash(recs[3])
	if err := VerifyChainWindow(recs); err == nil {
		t.Fatal("expected broken link detection")
	}
}

func TestComputeHashDeterministicAndTimeSensitive(t *testing.T) {
	r1 := Record{Seq: 1, Time: time.Date(2026, 11, 7, 12, 0, 0, 123456789, time.UTC), Actor: "a", Action: "x"}
	r2 := r1
	r2.Time = time.Date(2026, 11, 7, 12, 0, 0, 123456788, time.UTC)
	if ComputeHash(r1) != ComputeHash(r1) {
		t.Fatal("hash must be deterministic")
	}
	if ComputeHash(r1) == ComputeHash(r2) {
		t.Fatal("nanosecond timestamp must affect the hash")
	}
	// Non-UTC times must hash identically to their UTC equivalent.
	r3 := r1
	r3.Time = r1.Time.In(time.FixedZone("EAT", 3*3600))
	if ComputeHash(r1) != ComputeHash(r3) {
		t.Fatal("timezone conversion must not change the hash")
	}
}
