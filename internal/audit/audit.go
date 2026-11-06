// Package audit implements the hash-chained audit record format shared by all
// storage drivers. The chain is append-only: every record commits to the hash
// of its predecessor, so any retroactive edit breaks verification.
package audit

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"strconv"
	"strings"
	"time"
)

// GenesisPrevHash is the prev_hash of the very first record in a store.
const GenesisPrevHash = "0000000000000000000000000000000000000000000000000000000000000000"

// Record is one audit entry. Store drivers persist it as-is; Hash is computed
// over the other fields plus PrevHash. Time is always UTC.
type Record struct {
	Seq      int64     `json:"seq"`
	Time     time.Time `json:"time"`
	Actor    string    `json:"actor"`
	Action   string    `json:"action"`
	Payload  string    `json:"payload,omitempty"`
	PrevHash string    `json:"prev_hash"`
	Hash     string    `json:"hash"`
}

// ComputeHash derives the chain hash of r. The \x1f separator keeps fields
// unambiguous; actors, actions and payloads must not contain control bytes
// (payload is canonical JSON produced by the callers).
func ComputeHash(r Record) string {
	var b strings.Builder
	b.WriteString(strconv.FormatInt(r.Seq, 10))
	b.WriteByte(0x1f)
	b.WriteString(r.Time.UTC().Format(time.RFC3339Nano))
	b.WriteByte(0x1f)
	b.WriteString(r.Actor)
	b.WriteByte(0x1f)
	b.WriteString(r.Action)
	b.WriteByte(0x1f)
	b.WriteString(r.Payload)
	b.WriteByte(0x1f)
	b.WriteString(r.PrevHash)
	sum := sha256.Sum256([]byte(b.String()))
	return hex.EncodeToString(sum[:])
}

// VerifyChainFull validates a complete chain read from seq 1: continuity from
// the genesis hash, strictly incrementing sequence numbers, and every hash.
// Store.VerifyAuditChain implementations feed it the full record set.
func VerifyChainFull(records []Record) error {
	prev := GenesisPrevHash
	for i, r := range records {
		if r.Seq != int64(i+1) {
			return fmt.Errorf("record %d: seq %d out of order (expected %d)", i, r.Seq, i+1)
		}
		if r.PrevHash != prev {
			return fmt.Errorf("record %d (seq %d): prev_hash mismatch", i, r.Seq)
		}
		if got := ComputeHash(r); got != r.Hash {
			return fmt.Errorf("record %d (seq %d): hash mismatch (record tampered)", i, r.Seq)
		}
		prev = r.Hash
	}
	return nil
}

// VerifyChainWindow validates linkage inside a window of records that does not
// necessarily start at genesis (e.g. a tail query): every record after the
// first must reference its predecessor, and every hash must recompute.
func VerifyChainWindow(records []Record) error {
	for i := 1; i < len(records); i++ {
		if records[i].PrevHash != records[i-1].Hash {
			return fmt.Errorf("record %d (seq %d): prev_hash does not match predecessor", i, records[i].Seq)
		}
	}
	for i, r := range records {
		if got := ComputeHash(r); got != r.Hash {
			return fmt.Errorf("record %d (seq %d): hash mismatch (record tampered)", i, r.Seq)
		}
	}
	return nil
}
