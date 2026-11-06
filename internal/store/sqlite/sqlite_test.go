package sqlite

import (
	"context"
	"database/sql"
	"path/filepath"
	"testing"

	_ "modernc.org/sqlite"

	"github.com/amungadev-ai/Nexus_hosting_management/internal/store"
	"github.com/amungadev-ai/Nexus_hosting_management/internal/store/storetest"
)

func newSQLiteStore(t *testing.T) (store.Store, func()) {
	t.Helper()
	dsn := filepath.Join(t.TempDir(), "nexus.db")
	st, err := Open(dsn)
	if err != nil {
		t.Fatalf("open: %v", err)
	}
	return st, func() { st.Close() }
}

// TestConformance runs the shared DEC-001 suite against the SQLite driver.
func TestConformance(t *testing.T) {
	storetest.RunConformance(t, newSQLiteStore)
}

// TestTamperDetection proves the chain verifier catches a retroactive edit.
func TestTamperDetection(t *testing.T) {
	dsn := filepath.Join(t.TempDir(), "nexus.db")
	st, err := Open(dsn)
	if err != nil {
		t.Fatal(err)
	}
	defer st.Close()
	ctx := context.Background()
	for i := 0; i < 3; i++ {
		if _, err := st.AppendAudit(ctx, "t", "tick", ""); err != nil {
			t.Fatal(err)
		}
	}
	// Tamper through a second raw connection, as an attacker would.
	raw, err := sql.Open("sqlite", dsn)
	if err != nil {
		t.Fatal(err)
	}
	defer raw.Close()
	if _, err := raw.Exec(`UPDATE audit SET payload = '{"evil":true}' WHERE seq = 2`); err != nil {
		t.Fatal(err)
	}
	if err := st.VerifyAuditChain(ctx); err == nil {
		t.Fatal("tampered chain must fail verification")
	}
}
