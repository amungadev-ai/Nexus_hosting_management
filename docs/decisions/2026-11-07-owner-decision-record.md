# Owner Decision Record — 2026-11-07

**Source:** Owner rulings delivered in chat on 2026-11-07 (review of the NexusPanel decision register).
**Effect:** The rulings below upgrade the listed items from "proposed" to "owner-approved (with conditions)". Formal ADR files must still be drafted per Volume 10 §1.3; nothing here bypasses the security review gates named below.

## 1. Approved decisions

| ID | Subject | Ruling | Conditions / follow-ups |
|---|---|---|---|
| DEC-001 | Go control plane, Go agent, static SPA, SQLite/PostgreSQL dual driver | **Approved** | Build against a repository interface with Postgres first, **but run the SQLite driver in CI from the first commit** — this is the dual-engine equivalence check the spec itself requires. |
| DEC-011 | Guided router as primary IA | **Approved** | "Install → live 1-page site in under 10 minutes" becomes an **automated or timed script**, not a stated goal. |
| DEC-015 | cgroup v2 required | **Approved** | Fix the OS story before baselining: either drop the Ubuntu 22.04 claim or add 22.04 to the test matrix. Currently 24.04 + 26.04 are covered and 26.04.1 is only the dev reference (ADR-017 calls 26.04 "experimental until tested"). |
| ADR-005c | Docker as the isolation primitive (owner-operated single-tenant v1) | **Approved with security review before commitment** | Shared-kernel Docker is weak isolation for reseller/provider tiers; those tiers move toward rootless/untrusted containers (per Stage 7 OD-707). |
| Agent boundary (v1) | Unix socket + `SO_PEERCRED`; agent holds the Docker socket and runs as root | **Accepted for v1 single-host only** | The typed allowlist is the entire boundary → **agent fuzzing + external security review are a release gate** (before R1 commercial exposure), not a later item. |
| DEC-006 | M-Pesa Daraja + KES/VAT/KRA invoicing (async payment state machine) | **Approved** | Open item: confirm with a Kenyan tax advisor whether **eTIMS** integration is required (spec currently describes only PIN, VAT number and tax codes). Tracks COM-D02 (Volume 09 §10 already flags eTIMS as unvalidated). |

## 2. Inconsistencies to fix before the register becomes the baseline

1. **Nano memory arithmetic.** Reconcile the four conflicting statements:
   - Volume 06 footprint table: agent RSS 40 MB, **total panel ~125 MB, headroom ~387 MB** (= 512 − 125, ignores the OS entirely);
   - Volume 01 NP-INF-015 / NP-NFR-019 / DEC-004: **380 MB floor, panel + OS**, which leaves ~132 MB for workloads;
   - Volume 03 §17: "40 MB panel target is a hypothesis" (the 40 MB figure is the **agent alone**, not the panel);
   - Stage 1 §10.2 lists **100 MB as the panel floor** (not present in the consolidated volumes — apply when the stage docs/register are ingested).
   **Action:** one canonical Nano budget table (512 MB total → OS reserve / panel floor / agent / SQLite / workload headroom) drafted as an ADR before WP-04 benchmarking. Until then the 380 MB panel+OS floor (DEC-004) is the conservative operating assumption.
2. **Register tail.** Remove chat-style compilation commentary ("The decisions are yours") before the 29-area v2.0 functional-scope register becomes the canonical baseline. *(The register file was not among the 10 uploaded volumes — verified by search on 2026-11-07; apply the fix when it is ingested.)*
3. **Register sourcing.** Update the stale "assessment based on six available documents" note to reference the full consolidated set (now `docs/volumes/01–10` in this repo). *(Same: apply on ingest.)*
4. **v1 scope.** v1 as currently specified is very large (Go agent, installer, ServiceTemplate with three reference services, M-Pesa + invoicing, router, audit chain, verified backups, Docker sites, builds, AI diagnostics). Define a **thin walking skeleton** and sequence the rest behind it: **install → detect capabilities → guided router → one static site (HTTPS) → audit chain.** Consistent with Volume 10 R1 and iterations 01–12.

## 3. Additional anomalies found during full end-to-end ingestion (2026-11-07)

Recorded for the pre-baseline fix list beyond the owner's four items:

- **Volume 06:** AGT-014a enumerates **six** detection steps but the text, the v1 list and AT-031 call it the "5-step" detection. §2.1 diagram shows "Go Host Agent (root)" while §0.1/§4 mandate an unprivileged supervisor + narrow helper. Socket ownership stated as `nexuspanel:nexuspanel` (AGT-001a) vs `nexuspanel:nexuspanel-agent` (§4.2). §32 agent idle RSS ≤150 MiB conflicts with the 40 MB footprint row.
- **Volume 05:** line 93 orphaned fragment ("in v1. Not mTLS. Not network…"); commercial-surface gating stated three different ways (403 Forbidden / FEATURE_UNAVAILABLE / FEATURE_DISABLED); §26.4 "5-step detection" vs Volume 06's six; Nano example `available_ram_mb: 120` conflicts with Volume 06's ~387 MB headroom.
- **Volume 09:** section numbering broken/duplicated ("## 6" and "## 18" each appear twice; §3 body empty).
- **Volume 01:** "End of Stage 2 — PRD v2.0" appears mid-volume with a bolted-on decision-gates section; §28 DEC table duplicated; orphan bullets at lines ~852–855 and ~908–916; unsourced "30% of engineering effort / 90% of operational pain is email" claim.
- **Volume 02:** "27-area IA" vs the authoritative 29-area scope; several screens marked "NEW in v2" that describe v1 behavior (DNS guidance, M-Pesa pending screens); FLOW-015 maps to nonexistent `NP-TPL-*` IDs; trailing gates section appears after the "End of Volume 02" marker.
- **Volume 08:** NP-SEC-014 states "Ubuntu 22.04+ / 24.04+" baseline while SEC8-D13/T036 reference 26.04/26.04.1 (the DEC-015 fix above); audit chain alternately agent-appended vs single-writer (reconcilable only via a sequencer).
- **Volume 04:** sections §23–25 appended after the "End of Volume 04" marker; no `MIG-001` identifier in this volume (it lives in Volume 10 §12).
- **Volume 10:** §8 conformance matrix rows separated by blank lines (formatting artifact).
- **Ubuntu support-matrix tension (cross-volume):** 22.04 claimed in baselines, untested anywhere — resolved by the DEC-015 ruling above.

**Register provenance check:** the phrases "The decisions are yours" and "six available documents" do **not** appear anywhere in the 10 consolidated volumes (verified 2026-11-07); they live only in the separate 29-area register document, which must be ingested and cleaned before it becomes the canonical baseline.

## 4. Test environment ruling

- The owner's Leaseweb VPS **is the disposable test environment** ("my VPS is the test environment, and nothing more is needed to start" — owner authorization 2026-11-07 supersedes the drafting-time "no changes" boundary for test purposes only). Production support claims still require clean-room OS/profile matrices per Volume 10 §4.
- Profile definitions (as code) when VM-based tests are added: **Nano** (24.04, 512 MB / 1 vCPU / 10 GB), **Starter** (24.04, 2 GB / 1 / 25 GB), **Standard** (24.04 + a 26.04.1 variant, 6 GB / 4 / 100 GB); plus a second host for multi-agent tests and an isolated host for adversarial tenant workloads. Use real VMs (Multipass / libvirt-QEMU / Vagrant+cloud-init), not Docker-in-Docker — the agent, systemd units and cgroup behavior need their own kernel; snapshot before each destructive test.
- Safety rules: Daraja **sandbox only**; ACME **staging** with a throwaway domain; outbound mail to a sink; no production secrets anywhere.
- Exit criteria for "environment established": one command destroys and rebuilds each profile; `stat -fc %T /sys/fs/cgroup` reports `cgroup2fs` on all of them; fault injection (process kill, disk full, network partition) runnable against any host; a golden snapshot exists per profile.

## 5. Credential handling policy

- All credentials (VPS SSH, GitHub PAT, M-Pesa Daraja sandbox app) live **only** in the local sandbox secrets store *outside this repository* and in the original chat. They are **never committed**; a pre-push scan enforces this. CI uses no secrets until a proper secrets mechanism exists.
- Everything is revoked at project completion (owner's stated intent).
