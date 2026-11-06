package postgres

import (
        "context"
        "database/sql"
        "os"
        "testing"

        "github.com/amungadev-ai/Nexus_hosting_management/internal/store"
        "github.com/amungadev-ai/Nexus_hosting_management/internal/store/storetest"
)

// testDSN comes from the environment; CI provides it with a postgres:16
// service container. Locally the suite skips unless TEST_POSTGRES_DSN is set.
func testDSN(t *testing.T) string {
        t.Helper()
        dsn := os.Getenv("TEST_POSTGRES_DSN")
        if dsn == "" {
                t.Skip("TEST_POSTGRES_DSN not set; PostgreSQL conformance runs in CI (DEC-001)")
        }
        return dsn
}

func newPostgresStore(t *testing.T) (store.Store, func()) {
        t.Helper()
        dsn := testDSN(t)
        st, err := Open(dsn)
        if err != nil {
                t.Fatalf("open: %v", err)
        }
        // Start each subtest from empty tables (schema is idempotent).
        clean := func() {
                ctx := context.Background()
                _, _ = st.(*pgStore).db.ExecContext(ctx, `TRUNCATE audit RESTART IDENTITY`)
                _, _ = st.(*pgStore).db.ExecContext(ctx, `TRUNCATE sites`)
        }
        clean()
        return st, func() { st.Close() }
}

// TestConformance runs the shared DEC-001 suite against the PostgreSQL driver.
func TestConformance(t *testing.T) {
        storetest.RunConformance(t, newPostgresStore)
}

func openRaw(dsn string) (*sql.DB, error) {
        return sql.Open("pgx", dsn)
}

// TestTamperDetection proves the chain verifier catches a retroactive edit.
func TestTamperDetection(t *testing.T) {
        dsn := testDSN(t)
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
        raw, err := openRaw(dsn)
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
