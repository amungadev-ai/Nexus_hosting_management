# Walking skeleton — thin v1 slice (R0 exit target)

Owner ruling §2.4: v1 as drafted is too large. This slice sequences everything else behind it.

## Scope (exactly five legs)

1. **Install** — one timed command (`deploy/install.sh`) that installs dependencies, builds,
   installs systemd units, starts the control plane and host agent, and prints the elapsed time
   (DEC-011 condition: automated and timed, < 10 min target).
2. **Detect capabilities** — `internal/capability`: cgroup v2 gate (DEC-015), OS match
   (ADR-R0-002), RAM/disk, virtualization, Docker presence, port availability; exposed via
   `npctl detect` and `GET /api/v1/capabilities`; enforced fail-closed by `installer preflight`.
3. **Guided router (API)** — `GET /api/v1/router/steps` +
   `POST /api/v1/router/steps/{id}/complete`: the five skeleton steps as state; the SPA wires
   the IA later (DEC-011).
4. **One static site** — `POST /api/v1/sites` provisions `index.html` under the data dir and
   serves it on the site listener (host-routed) and via `/s/{hostname}/` on the control plane;
   no Docker (ADR-005c gate), no ACME yet.
5. **Audit chain** — SHA-256 hash-chained, single-writer audit records in the store; every
   skeleton state change appends; `GET /api/v1/audit/verify` and `npctl audit verify`
   re-validate the chain.

## Acceptance criteria

- [ ] `deploy/install.sh` completes on a clean Ubuntu 24.04/26.04 VM-class host in < 600 s
      (timed output printed by the script).
- [ ] `installer preflight` fails closed when cgroup v2 is absent or RAM < 512 MB.
- [ ] Both SQLite and PostgreSQL drivers pass the shared conformance suite in CI (DEC-001
      condition) from the first code commit.
- [ ] A provisioned site is fetchable; provisioning is recorded in the audit chain; the chain
      verifies.
- [ ] The agent accepts only its typed built-in commands over the Unix socket; unknown commands
      are rejected.

## Explicitly deferred (sequenced behind the skeleton)

- **HTTPS/ACME**: needs a throwaway domain per the test-safety rules; the site listener keeps a
  TLS-ready seam so this is additive.
- **Docker-based site provisioning**: gated by ADR-005c + ADR-R0-003 review.
- **Auth/multi-user, M-Pesa/invoicing, backups, AI diagnostics, email, DNS guidance**: later
  iterations (Volume 10).
- Router step state is in-memory in the skeleton; durable state = audit chain + site records
  only (documented limitation).
