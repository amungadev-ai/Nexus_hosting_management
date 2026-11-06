// Package store defines the repository interface all NexusPanel storage
// drivers implement (DEC-001: build against this interface, Postgres first,
// SQLite kept equivalent in CI from the first commit).
package store

import (
	"context"
	"errors"
	"time"

	"github.com/amungadev-ai/Nexus_hosting_management/internal/audit"
)

var (
	// ErrNotFound is returned when a requested record does not exist.
	ErrNotFound = errors.New("store: not found")
	// ErrDuplicate is returned when a uniqueness constraint is violated.
	ErrDuplicate = errors.New("store: duplicate")
	// ErrInvalidLimit is returned for AuditTail limits <= 0.
	ErrInvalidLimit = errors.New("store: invalid limit")
)

// Site lifecycle states (v1).
const (
	SiteActive   = "active"
	SiteDisabled = "disabled"
)

// Site is one provisioned static site.
type Site struct {
	ID        string    `json:"id"`
	Hostname  string    `json:"hostname"`
	RootPath  string    `json:"root_path"`
	State     string    `json:"state"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// Store is the single repository interface for the walking skeleton: the
// hash-chained audit log plus site records.
type Store interface {
	Ping(ctx context.Context) error
	Close() error

	// AppendAudit appends one record to the hash chain, filling Seq, Time
	// (UTC), PrevHash and Hash. It is safe for concurrent use.
	AppendAudit(ctx context.Context, actor, action, payload string) (audit.Record, error)
	// AuditTail returns the n most recent records in ascending Seq order.
	// n must be > 0, else ErrInvalidLimit.
	AuditTail(ctx context.Context, n int) ([]audit.Record, error)
	// VerifyAuditChain re-reads every record and validates the full chain
	// from genesis. A nil return means the chain is intact.
	VerifyAuditChain(ctx context.Context) error

	CreateSite(ctx context.Context, s Site) error
	GetSite(ctx context.Context, id string) (Site, error)
	GetSiteByHostname(ctx context.Context, hostname string) (Site, error)
	ListSites(ctx context.Context) ([]Site, error)
	SetSiteState(ctx context.Context, id, state string) error
}
