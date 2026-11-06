# ADR-R0-003 — Host-agent security review is a release gate

- **ID:** ADR-R0-003 (working ID; renumber into the formal ADR register per Volume 10 §1.3)
- **Date:** 2026-11-07
- **Status:** Accepted — owner ruling: agent fuzzing + external security review are a
  **release gate**, not a later item

## Context

The v1 agent is a root process holding the Docker socket, reachable only through a local Unix
socket with `SO_PEERCRED` peer verification. The typed command allowlist is the **entire**
boundary. ADR-005c (Docker as isolation primitive) is approved only *with* this security review
before commitment.

## Plan

1. **Fuzzing (continuous).** Go native fuzz targets over the agent's request decoder: malformed
   JSON, oversized frames, unknown/aliased commands, peer-credential edge cases (agent restart
   while a client is connected). Minimum 1 hour of CI fuzzing before any release; any crash
   blocks the release.
2. **Boundary review (external, before R1 commercial exposure).** Checklist:
   - Unix socket path, ownership and mode (0o660); socket directory not world-writable;
   - `SO_PEERCRED` → uid mapping and the policy table (who may call what);
   - allowlist completeness against the documented capability set — every accepted command
     enumerated and typed, no shell interpolation;
   - Docker socket blast radius and why the agent needs it;
   - failure modes: agent crash/restart, socket recreation, partially-written state.
3. **Docker commitment gate.** No Docker-based site provisioning code ships until review
   findings are closed; the walking skeleton deliberately provisions static sites without Docker.

## Exit criteria

- Zero open high/critical findings; fuzz corpus green for the full run; review report stored
  under `docs/security/`.
- The skeleton's agent implements built-in, typed commands only (no `exec`), keeping the
  reviewed surface minimal.
