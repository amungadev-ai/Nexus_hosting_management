# ADR-R0-002 — OS support matrix (resolves the Ubuntu 22.04 claim)

- **ID:** ADR-R0-002 (working ID; renumber into the formal ADR register per Volume 10 §1.3)
- **Date:** 2026-11-07
- **Status:** Accepted — owner DEC-015 ruling: fix the OS story before baselining

## Decision

| OS | Support level | Notes |
|---|---|---|
| Ubuntu 24.04 LTS | **First-class supported** | Target of the planned Nano/Starter test profiles |
| Ubuntu 26.04 / 26.04.1 | Dev reference; **experimental until profile tests pass** | Per ADR-017 / SEC8-D13; the owner's test VPS runs 26.04.1 (probed 2026-11-07, `cgroup2fs`) |
| Ubuntu 22.04 | **Dropped** | Claimed in NP-SEC-014 (Volume 08) but present in no test profile; the cgroup v2 policy was never exercised on it |

1. All v1 targets require **unified cgroup v2** (`stat -fc %T /sys/fs/cgroup` → `cgroup2fs`) per
   DEC-015; 24.04 and 26.04 both default to unified v2.
2. NP-SEC-014's "Ubuntu 22.04+ / 24.04+" baseline is corrected to "Ubuntu 24.04 LTS;
   26.04 experimental" (see `docs/volumes/CORRECTIONS.md`).
3. Production support claims still require the clean-room OS/profile matrices of Volume 10 §4 —
   the owner's VPS ruling covers testing only, not production claims.

## Evidence

- Leaseweb test VPS probe (2026-11-07): Ubuntu 26.04.1, kernel 7.0.0-38-generic, KVM,
  `cgroup2fs` with controllers cpuset/cpu/io/memory/hugetlb/pids/rdma/misc/dmem.
- The panel's capability detector enforces the cgroup v2 gate at install time (walking skeleton),
  so an unsupported host fails closed instead of silently degrading.
