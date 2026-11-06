# NexusPanel — Baseline Status

**Updated:** 2026-11-07 (R0 exit)
**Repo:** github.com/amungadev-ai/Nexus_hosting_management (this repository)
**Stage:** R0 EXITED — walking skeleton implemented, deployed to the test VPS, CI green on both drivers.

## What exists now

| Item | Location | State |
|---|---|---|
| Consolidated 10-volume specification (8,555 lines) | `docs/volumes/01–10_*.md` | Ingested end to end 2026-11-07 |
| Owner decision record | `docs/decisions/2026-11-07-owner-decision-record.md` | Captured verbatim from owner rulings |
| WP-01 working ADRs (R0-001 Nano budget, R0-002 OS matrix, R0-003 agent gate) | `docs/decisions/2026-11-07-adr-r0-00*.md` | Applied; owner ratification pending for R0-001 |
| Corrections mapping into the volumes | `docs/volumes/CORRECTIONS.md` | 6 items mapped; register fixes pending ingest |
| Walking-skeleton spec + acceptance criteria | `docs/specs/walking-skeleton.md` | All criteria verified on the test VPS 2026-11-07 |
| Go monorepo (`nexus-control-plane`, `nexus-host-agent`, `nexus-npctl`, `nexus-installer`) | `cmd/`, `internal/` | vet + tests green; shared dual-driver conformance suite |
| CI — DEC-001 both drivers from first code commit | `.github/workflows/ci.yml` | sqlite ✅ · postgres ✅ · secret-pattern-scan ✅ |
| Test-env tooling (owner ruling §4) | `scripts/testenv/` | reset-baseline / snapshot / fault-inject — all exercised on the VPS |
| Test VPS facts + exit criteria | `docs/environment/leaseweb-test-vps.md` | All four criteria PASSED 2026-11-07 |

## Walking skeleton — verified on the owner's VPS (2026-11-07)

- `deploy/install.sh` timed: **48 s** on the VPS (target < 600 s) — DEC-011 condition met.
- Preflight verified fail-closed (gates fire on unsupported OS / missing cgroup v2 — confirmed on a Debian container).
- Capability report on the VPS: Ubuntu 26.04.1, virt kvm, `cgroup_v2: true`, 4 vCPU, 5920 MB, 90 GB free.
- One static site provisioned and served host-routed (`:8081`) and via the `/s/{hostname}/` path; duplicate and hostile hostnames rejected (409 / 400).
- Audit chain: SHA-256 hash-chained; **VALID across SIGKILL crash-restart** (systemd auto-restart) and across binary upgrades; single-writer via SQLite `BEGIN IMMEDIATE` / Postgres advisory lock.
- Guided router: five steps exposed via API, completed end to end, every state change audited.
- Environment-established criteria: **all four PASSED** — `cgroup2fs`; one-command reset; fault injection; golden snapshot.
- Credentials on the VPS: none (source tarball only, built on the host; Daraja/PAT never transferred).

## Lesson recorded

- The first CI run failed the postgres job: the driver used SQLite's `?` placeholders.
  DEC-001's dual-driver CI did exactly its job — a driver that compiled and passed its own
  unit semantics locally (no local server) failed against a real PostgreSQL 16. Fixed in
  `776ae80`, with a CI regression guard grepping the driver for `?`.

## Next steps (post-skeleton sequencing)

1. Owner ratification of ADR-R0-001 (canonical Nano budget) — one-line approval.
2. Ingest the 29-area register document (not among the uploads) → apply owner-record §2.2–2.3
   fixes (chat tail, stale sourcing) → becomes the canonical scope baseline.
3. HTTPS/ACME leg: needs a throwaway domain per the test-safety rules → then the site listener
   gains TLS at the seam already reserved for it.
4. Docker-based site provisioning: only after the ADR-R0-003 external security review
   (release gate before R1 commercial exposure).
5. eTIMS confirmation with a Kenyan tax advisor (DEC-006 / COM-D02).
6. WP-04 benchmarks validate or revise the Nano splits (QA-NANO-014/015, explicit swap policy).
