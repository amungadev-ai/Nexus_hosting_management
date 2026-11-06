# NexusPanel — Baseline Status

**Updated:** 2026-11-07
**Repo:** github.com/amungadev-ai/Nexus_hosting_management (this repository)
**Stage:** R0 — specification ingested and baselined; environment verified; **no panel code yet**.

## What exists now

| Item | Location | State |
|---|---|---|
| Consolidated 10-volume specification (8,555 lines) | `docs/volumes/01–10_*.md` | Ingested end to end 2026-11-07 (full-text digests + targeted greps) |
| Owner decision record (DEC-001/011/015, ADR-005c, agent boundary, DEC-006; register fixes; test-env ruling; credential policy) | `docs/decisions/2026-11-07-owner-decision-record.md` | Captured verbatim from owner rulings |
| Ingestion anomaly log (9 cross-volume inconsistencies beyond the owner's 4) | same file, §3 | Recorded for the pre-baseline fix pass |
| Test VPS facts (Ubuntu 26.04.1, 4 vCPU, 5.8 GiB, 99 GB, **cgroup2fs ✅**, clean stack) | `docs/environment/leaseweb-test-vps.md` | Verified by read-only probe 2026-11-07 |
| Credentials (VPS SSH, GitHub PAT, Daraja sandbox) | **Local sandbox store only, outside this repo** | Never committed; pre-push scan enforces; revoke at project completion |

## Approved with conditions (owner, 2026-11-07)

- **DEC-001 stack** — approved; SQLite driver must run in CI from the first commit (repository interface, Postgres-first).
- **DEC-011 guided router** — approved; "live 1-page site in <10 min" becomes an automated timed script.
- **DEC-015 cgroup v2** — approved; fix the Ubuntu 22.04 claim (drop it or add it to the matrix).
- **ADR-005c Docker v1** — approved with security review before commitment; reseller/provider tiers need stronger isolation.
- **Agent boundary (Unix socket + SO_PEERCRED)** — accepted for v1 single-host; agent fuzzing + external security review are a **release gate**.
- **DEC-006 M-Pesa** — approved; open item: eTIMS requirement confirmation with a Kenyan tax advisor (COM-D02).

## Next steps (R0 → walking skeleton)

1. **WP-01 ADRs:** canonical Nano memory budget (fix the 125/380/100/40 MB conflict); OS support matrix (22.04 decision); agent security-review plan.
2. **Scaffold** the Go monorepo (`cmd/control-plane`, `cmd/host-agent`, `cmd/npctl`, `cmd/installer`) + CI that runs **both SQLite and PostgreSQL drivers from the first commit** (DEC-001 condition).
3. **Walking skeleton** (thin v1 slice, per owner ruling §2.4): install → capability detection → guided router → one static HTTPS site → audit chain. Everything else sequences behind it.
4. **Test-env tooling:** profile definitions as code (Nano/Starter/Standard), one-command rebuild, fault-injection harness, golden snapshots.
5. **Ingest & clean the 29-area register** (remove chat tail, update sourcing) before it becomes the canonical scope baseline.
6. **External confirmations:** eTIMS/Kenyan tax advisor (DEC-006); nothing else blocks R0 exit.
