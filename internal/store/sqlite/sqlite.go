// Package sqlite is the SQLite Store driver (modernc.org/sqlite, pure Go —
// no cgo). It must remain semantically equivalent to the postgres driver;
// both run the same conformance suite in CI (DEC-001).
package sqlite

import (
	"context"
	"database/sql"
	"fmt"
	"strings"
	"sync"
	"time"

	_ "modernc.org/sqlite"

	"github.com/amungadev-ai/Nexus_hosting_management/internal/audit"
	"github.com/amungadev-ai/Nexus_hosting_management/internal/store"
)

// Open opens (creating if needed) the SQLite database at dsn and applies the
// pragmas the audit chain relies on.
func Open(dsn string) (store.Store, error) {
	db, err := sql.Open("sqlite", dsn)
	if err != nil {
		return nil, fmt.Errorf("sqlite: open %s: %w", dsn, err)
	}
	for _, pragma := range []string{
		"PRAGMA busy_timeout = 5000",
		"PRAGMA journal_mode = WAL",
		"PRAGMA synchronous = NORMAL",
		"PRAGMA foreign_keys = ON",
	} {
		if _, err := db.Exec(pragma); err != nil {
			db.Close()
			return nil, fmt.Errorf("sqlite: %s: %w", pragma, err)
		}
	}
	schema := `
CREATE TABLE IF NOT EXISTS audit (
	seq       INTEGER PRIMARY KEY,
	ts        TEXT NOT NULL,
	actor     TEXT NOT NULL,
	action    TEXT NOT NULL,
	payload   TEXT NOT NULL DEFAULT '',
	prev_hash TEXT NOT NULL,
	hash      TEXT NOT NULL
);
CREATE TABLE IF NOT EXISTS sites (
	id         TEXT PRIMARY KEY,
	hostname   TEXT NOT NULL UNIQUE,
	root_path  TEXT NOT NULL,
	state      TEXT NOT NULL,
	created_at TEXT NOT NULL,
	updated_at TEXT NOT NULL
);`
	if _, err := db.Exec(schema); err != nil {
		db.Close()
		return nil, fmt.Errorf("sqlite: schema: %w", err)
	}
	return &sqliteStore{db: db}, nil
}

type sqliteStore struct {
	db *sql.DB
	mu sync.Mutex // single-writer discipline for the audit chain
}

func (s *sqliteStore) Ping(ctx context.Context) error { return s.db.PingContext(ctx) }
func (s *sqliteStore) Close() error                   { return s.db.Close() }

func (s *sqliteStore) AppendAudit(ctx context.Context, actor, action, payload string) (audit.Record, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return audit.Record{}, err
	}
	defer tx.Rollback()

	var seq int64
	var prev string
	err = tx.QueryRowContext(ctx, `SELECT COALESCE(MAX(seq),0), COALESCE(
		(SELECT hash FROM audit ORDER BY seq DESC LIMIT 1), '') FROM audit`).Scan(&seq, &prev)
	if err != nil && err != sql.ErrNoRows {
		return audit.Record{}, err
	}
	if seq == 0 {
		prev = audit.GenesisPrevHash
	}
	seq++

	rec := audit.Record{
		Seq:      seq,
		Time:     time.Now().UTC(),
		Actor:    actor,
		Action:   action,
		Payload:  payload,
		PrevHash: prev,
	}
	rec.Hash = audit.ComputeHash(rec)
	_, err = tx.ExecContext(ctx,
		`INSERT INTO audit (seq, ts, actor, action, payload, prev_hash, hash) VALUES (?,?,?,?,?,?,?)`,
		rec.Seq, rec.Time.Format(time.RFC3339Nano), rec.Actor, rec.Action, rec.Payload, rec.PrevHash, rec.Hash)
	if err != nil {
		return audit.Record{}, err
	}
	if err := tx.Commit(); err != nil {
		return audit.Record{}, err
	}
	return rec, nil
}

func scanRecord(scan func(dest ...any) error) (audit.Record, error) {
	var r audit.Record
	var ts string
	if err := scan(&r.Seq, &ts, &r.Actor, &r.Action, &r.Payload, &r.PrevHash, &r.Hash); err != nil {
		return audit.Record{}, err
	}
	t, err := time.Parse(time.RFC3339Nano, ts)
	if err != nil {
		return audit.Record{}, fmt.Errorf("sqlite: bad timestamp %q: %w", ts, err)
	}
	r.Time = t.UTC()
	return r, nil
}

const auditCols = `seq, ts, actor, action, payload, prev_hash, hash`

func (s *sqliteStore) AuditTail(ctx context.Context, n int) ([]audit.Record, error) {
	if n <= 0 {
		return nil, store.ErrInvalidLimit
	}
	rows, err := s.db.QueryContext(ctx,
		`SELECT `+auditCols+` FROM audit ORDER BY seq DESC LIMIT ?`, n)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var recs []audit.Record
	for rows.Next() {
		r, err := scanRecord(rows.Scan)
		if err != nil {
			return nil, err
		}
		recs = append(recs, r)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	// reverse into ascending order
	for i, j := 0, len(recs)-1; i < j; i, j = i+1, j-1 {
		recs[i], recs[j] = recs[j], recs[i]
	}
	return recs, nil
}

func (s *sqliteStore) VerifyAuditChain(ctx context.Context) error {
	rows, err := s.db.QueryContext(ctx, `SELECT ` + auditCols + ` FROM audit ORDER BY seq ASC`)
	if err != nil {
		return err
	}
	defer rows.Close()
	var recs []audit.Record
	for rows.Next() {
		r, err := scanRecord(rows.Scan)
		if err != nil {
			return err
		}
		recs = append(recs, r)
	}
	if err := rows.Err(); err != nil {
		return err
	}
	return audit.VerifyChainFull(recs)
}

const siteCols = `id, hostname, root_path, state, created_at, updated_at`

func scanSite(scan func(dest ...any) error) (store.Site, error) {
	var s store.Site
	var created, updated string
	if err := scan(&s.ID, &s.Hostname, &s.RootPath, &s.State, &created, &updated); err != nil {
		return store.Site{}, err
	}
	var err error
	if s.CreatedAt, err = time.Parse(time.RFC3339Nano, created); err != nil {
		return store.Site{}, fmt.Errorf("sqlite: bad created_at: %w", err)
	}
	if s.UpdatedAt, err = time.Parse(time.RFC3339Nano, updated); err != nil {
		return store.Site{}, fmt.Errorf("sqlite: bad updated_at: %w", err)
	}
	return s, nil
}

func (s *sqliteStore) CreateSite(ctx context.Context, site store.Site) error {
	_, err := s.db.ExecContext(ctx,
		`INSERT INTO sites (id, hostname, root_path, state, created_at, updated_at) VALUES (?,?,?,?,?,?)`,
		site.ID, site.Hostname, site.RootPath, site.State,
		site.CreatedAt.UTC().Format(time.RFC3339Nano), site.UpdatedAt.UTC().Format(time.RFC3339Nano))
	if err != nil && strings.Contains(err.Error(), "UNIQUE constraint failed") {
		return store.ErrDuplicate
	}
	return err
}

func (s *sqliteStore) GetSite(ctx context.Context, id string) (store.Site, error) {
	row := s.db.QueryRowContext(ctx, `SELECT `+siteCols+` FROM sites WHERE id = ?`, id)
	site, err := scanSite(row.Scan)
	if err == sql.ErrNoRows {
		return store.Site{}, store.ErrNotFound
	}
	return site, err
}

func (s *sqliteStore) GetSiteByHostname(ctx context.Context, hostname string) (store.Site, error) {
	row := s.db.QueryRowContext(ctx, `SELECT `+siteCols+` FROM sites WHERE hostname = ?`, hostname)
	site, err := scanSite(row.Scan)
	if err == sql.ErrNoRows {
		return store.Site{}, store.ErrNotFound
	}
	return site, err
}

func (s *sqliteStore) ListSites(ctx context.Context) ([]store.Site, error) {
	rows, err := s.db.QueryContext(ctx, `SELECT `+siteCols+` FROM sites ORDER BY created_at ASC`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var sites []store.Site
	for rows.Next() {
		site, err := scanSite(rows.Scan)
		if err != nil {
			return nil, err
		}
		sites = append(sites, site)
	}
	return sites, rows.Err()
}

func (s *sqliteStore) SetSiteState(ctx context.Context, id, state string) error {
	res, err := s.db.ExecContext(ctx,
		`UPDATE sites SET state = ?, updated_at = ? WHERE id = ?`,
		state, time.Now().UTC().Format(time.RFC3339Nano), id)
	if err != nil {
		return err
	}
	if n, err := res.RowsAffected(); err == nil && n == 0 {
		return store.ErrNotFound
	}
	return nil
}
