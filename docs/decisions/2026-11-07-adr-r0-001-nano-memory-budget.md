# ADR-R0-001 — Canonical Nano memory budget (512 MB profile)

- **ID:** ADR-R0-001 (working ID; renumber into the formal ADR register per Volume 10 §1.3)
- **Date:** 2026-11-07
- **Status:** Proposed for owner ratification — derived from the owner ruling of 2026-11-07 §2.1
- **Supersedes:** the four conflicting Nano memory statements listed below

## Context

Four places in the ingested material state incompatible Nano (512 MB) memory figures:

| Source | Claim | Problem |
|---|---|---|
| Volume 06 footprint table | panel ~125 MB total, headroom ~387 MB | Ignores the OS entirely (512 − 125 = 387) |
| Volume 01 NP-INF-015 / NP-NFR-019 / DEC-004 | 380 MB floor, panel + OS | Leaves ~132 MB for workloads — conflicts with the 387 MB figure above |
| Volume 03 §17 | "40 MB panel target is a hypothesis" | The 40 MB figure is the **agent alone**, not the panel |
| Stage 1 §10.2 (register) | 100 MB panel floor | Not present in the consolidated volumes |

Additionally, Volume 05's Nano example reports `available_ram_mb: 120`, inconsistent with every
other figure, and Volume 08 §32 gives the agent an idle RSS ceiling of ≤150 MiB versus the
40 MB footprint row (floor vs ceiling — both are now recorded).

## Decision

The canonical Nano budget until WP-04 benchmarks say otherwise:

| Segment | Budget | Basis |
|---|---:|---|
| **Total (Nano profile)** | **512 MB** | Profile definition |
| OS reserve (kernel + minimal server userland) | 280 MB | Derived: DEC-004 380 MB "panel + OS" floor − 100 MB panel floor |
| Control-plane (panel) floor | 100 MB | Stage 1 §10.2 |
| Host agent floor (idle) | 40 MB | Volume 06 footprint row (agent alone; confirmed by Volume 03 §17) |
| **Workload headroom** | **92 MB** | Derived: 512 − 380 − 40 |

Rules:

1. DEC-004's **380 MB panel + OS floor** is the conservative operating assumption (owner ruling §2.1).
2. The Volume 06 figures "panel ~125 MB / headroom ~387 MB" are **withdrawn** (they ignored the OS).
3. The Volume 05 example value `available_ram_mb: 120` for Nano is **withdrawn**.
4. Agent memory is recorded as **floor 40 MB (idle) / ceiling 150 MiB (Volume 08 §32)** until benchmarked.
5. **Consequence:** with 92 MB of headroom, running tenant containers on Nano is **not a supported
   v1 scenario**; Nano is a control-plane bring-up and test target.

## Consequences

- WP-04 (installer/benchmarks) must validate or revise the 280/100/40 splits; this ADR is
  re-ratified after those benchmarks.
- Memory-pressure benchmarks (QA-NANO-014/015) must state the swap policy per run — the owner's
  test VPS has no swap configured.
- Any spec text contradicting this table is corrected via `docs/volumes/CORRECTIONS.md`.
