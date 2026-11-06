# Corrections applied to the ingested volumes (2026-11-07)

The volumes in `docs/volumes/` are the ingested baseline and are kept verbatim. Normative
corrections live in working ADRs (`docs/decisions/`); the mapping below is the authoritative
list. The 29-area register document was **not** among the uploads; its two fixes (chat-style
tail, stale "six available documents" sourcing note) apply when it is ingested
(owner decision record §2 items 2–3).

| # | Location | Issue | Correction |
|---|---|---|---|
| 1 | Volume 06 footprint table | "panel ~125 MB, headroom ~387 MB" ignores the OS | Superseded by ADR-R0-001 canonical budget |
| 2 | Volume 05 Nano example | `available_ram_mb: 120` | Withdrawn; use ADR-R0-001 |
| 3 | Volume 08 NP-SEC-014 | "Ubuntu 22.04+ / 24.04+" baseline | Superseded by ADR-R0-002 (24.04 first-class; 22.04 dropped) |
| 4 | Volume 03 §17 / Volume 06 | 40 MB attributed ambiguously | 40 MB = agent idle floor (ADR-R0-001) |
| 5 | Volume 08 §32 vs Volume 06 | agent ≤150 MiB vs 40 MB row | Floor 40 MB / ceiling 150 MiB pending WP-04 benchmarks |
| 6 | 29-area register (not yet ingested) | Chat-style tail ("The decisions are yours"); "six available documents" | Pending register ingest (owner record §2.2–2.3) |
