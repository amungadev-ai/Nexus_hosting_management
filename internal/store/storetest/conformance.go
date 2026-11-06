// Package storetest is the driver-agnostic conformance suite every Store
// implementation must pass (DEC-001 dual-driver equivalence check).
package storetest

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/amungadev-ai/Nexus_hosting_management/internal/audit"
	"github.com/amungadev-ai/Nexus_hosting_management/internal/store"
)

// Factory opens a Store against a fresh, empty database. The returned cleanup
// function is called after each subtest.
type Factory func(t *testing.T) (store.Store, func())

// RunConformance runs the full suite against one driver instance.
func RunConformance(t *testing.T, newStore Factory) {
	t.Run("ping", func(t *testing.T) {
		st, done := newStore(t)
		defer done()
		if err := st.Ping(context.Background()); err != nil {
			t.Fatalf("ping: %v", err)
		}
	})

	t.Run("audit chain append and verify", func(t *testing.T) {
		st, done := newStore(t)
		defer done()
		ctx := context.Background()
		for i := 0; i < 5; i++ {
			payload := `{"i":` + itoa(i) + `}`
			if _, err := st.AppendAudit(ctx, "tester", "test.event", payload); err != nil {
				t.Fatalf("append %d: %v", i, err)
			}
		}
		if err := st.VerifyAuditChain(ctx); err != nil {
			t.Fatalf("chain must verify after appends: %v", err)
		}
		recs, err := st.AuditTail(ctx, 100)
		if err != nil {
			t.Fatalf("tail: %v", err)
		}
		if len(recs) != 5 {
			t.Fatalf("expected 5 records, got %d", len(recs))
		}
		for i, r := range recs {
			if r.Seq != int64(i+1) {
				t.Fatalf("record %d has seq %d, want %d", i, r.Seq, i+1)
			}
			if r.Time.Location() != time.UTC {
				t.Fatal("record times must be UTC")
			}
		}
	})

	t.Run("audit records carry actor action payload", func(t *testing.T) {
		st, done := newStore(t)
		defer done()
		ctx := context.Background()
		rec, err := st.AppendAudit(ctx, "owner", "site.provisioned", `{"hostname":"demo.test"}`)
		if err != nil {
			t.Fatalf("append: %v", err)
		}
		if rec.Seq != 1 || rec.Actor != "owner" || rec.Action != "site.provisioned" {
			t.Fatalf("record fields lost: %+v", rec)
		}
		if rec.Payload != `{"hostname":"demo.test"}` {
			t.Fatalf("payload lost: %q", rec.Payload)
		}
		if rec.PrevHash != audit.GenesisPrevHash {
			t.Fatal("first record prev_hash must be genesis")
		}
		tail, err := st.AuditTail(ctx, 1)
		if err != nil {
			t.Fatalf("tail: %v", err)
		}
		if len(tail) != 1 || tail[0].Hash != rec.Hash {
			t.Fatal("tail must return the appended record")
		}
	})

	t.Run("audit tail invalid limit", func(t *testing.T) {
		st, done := newStore(t)
		defer done()
		if _, err := st.AppendAudit(context.Background(), "t", "x", ""); err != nil {
			t.Fatal(err)
		}
		for _, n := range []int{0, -1} {
			if _, err := st.AuditTail(context.Background(), n); !errors.Is(err, store.ErrInvalidLimit) {
				t.Fatalf("AuditTail(%d) must return ErrInvalidLimit, got %v", n, err)
			}
		}
	})

	t.Run("audit tail ordering", func(t *testing.T) {
		st, done := newStore(t)
		defer done()
		ctx := context.Background()
		for i := 0; i < 5; i++ {
			if _, err := st.AppendAudit(ctx, "t", "tick", ""); err != nil {
				t.Fatal(err)
			}
		}
		tail, err := st.AuditTail(ctx, 2)
		if err != nil {
			t.Fatal(err)
		}
		if len(tail) != 2 {
			t.Fatalf("expected 2 records, got %d", len(tail))
		}
		if tail[0].Seq != 4 || tail[1].Seq != 5 {
			t.Fatalf("tail must return the n most recent in ascending order, got seq %d,%d", tail[0].Seq, tail[1].Seq)
		}
		if err := audit.VerifyChainWindow(tail); err != nil {
			t.Fatalf("tail window must be internally linked: %v", err)
		}
	})

	t.Run("sites CRUD", func(t *testing.T) {
		st, done := newStore(t)
		defer done()
		ctx := context.Background()
		now := time.Now().UTC().Truncate(time.Second)
		s1 := store.Site{ID: "id-1", Hostname: "one.test", RootPath: "/srv/one", State: store.SiteActive, CreatedAt: now, UpdatedAt: now}
		s2 := store.Site{ID: "id-2", Hostname: "two.test", RootPath: "/srv/two", State: store.SiteActive, CreatedAt: now, UpdatedAt: now}
		if err := st.CreateSite(ctx, s1); err != nil {
			t.Fatalf("create s1: %v", err)
		}
		if err := st.CreateSite(ctx, s2); err != nil {
			t.Fatalf("create s2: %v", err)
		}
		got, err := st.GetSite(ctx, "id-1")
		if err != nil {
			t.Fatalf("get s1: %v", err)
		}
		if got.Hostname != "one.test" || got.State != store.SiteActive {
			t.Fatalf("s1 mismatch: %+v", got)
		}
		byHost, err := st.GetSiteByHostname(ctx, "two.test")
		if err != nil {
			t.Fatalf("get by hostname: %v", err)
		}
		if byHost.ID != "id-2" {
			t.Fatalf("by-hostname mismatch: %+v", byHost)
		}
		sites, err := st.ListSites(ctx)
		if err != nil {
			t.Fatal(err)
		}
		if len(sites) != 2 {
			t.Fatalf("expected 2 sites, got %d", len(sites))
		}
		if err := st.SetSiteState(ctx, "id-2", store.SiteDisabled); err != nil {
			t.Fatalf("set state: %v", err)
		}
		got, err = st.GetSite(ctx, "id-2")
		if err != nil {
			t.Fatal(err)
		}
		if got.State != store.SiteDisabled {
			t.Fatalf("state not persisted: %+v", got)
		}
	})

	t.Run("sites not found and duplicate hostname", func(t *testing.T) {
		st, done := newStore(t)
		defer done()
		ctx := context.Background()
		if _, err := st.GetSite(ctx, "missing"); !errors.Is(err, store.ErrNotFound) {
			t.Fatalf("missing site must be ErrNotFound, got %v", err)
		}
		now := time.Now().UTC()
		s := store.Site{ID: "a", Hostname: "dupe.test", RootPath: "/x", State: store.SiteActive, CreatedAt: now, UpdatedAt: now}
		if err := st.CreateSite(ctx, s); err != nil {
			t.Fatal(err)
		}
		s.ID = "b"
		if err := st.CreateSite(ctx, s); !errors.Is(err, store.ErrDuplicate) {
			t.Fatalf("duplicate hostname must be ErrDuplicate, got %v", err)
		}
	})
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
