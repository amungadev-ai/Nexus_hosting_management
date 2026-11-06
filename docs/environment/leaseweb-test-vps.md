# Test Environment Facts — Leaseweb VPS (owner-designated disposable test host)

**Probed:** 2026-11-07 (read-only reconnaissance only — no state changed; only `ss`, `df`, `free`, `stat` and version queries were run)
**Authorization:** Owner designated this VPS as the test environment on 2026-11-07. SSH credentials live in the local sandbox secrets store, never in this repository.

## Facts

| Item | Observed value | Spec reference |
|---|---|---|
| Hostname | `VM-97057ab9-6a82-4554-969f-e3ab5db01a39` | — |
| OS | Ubuntu **26.04.1 LTS** | Dev reference node (Vol 03 §16, Vol 06 §5); SEC8-D13 compat gate |
| Kernel | 7.0.0-38-generic | — |
| Virtualization | KVM | Real-kernel VM host — suitable for nested-VM profile tests if ever needed |
| CPU | 4 vCPU, AMD EPYC | Dev node reference: 4 vCPU |
| RAM | 5920 MB total, 5478 MB available | Dev node reference: ~6 GiB |
| Swap | **None configured** | Note: Nano OOM/pressure tests need deliberate swap policy per run |
| Disk | 99 GB total, 2.2 GB used, 93 GB free | Dev node reference: ~100 GB |
| **cgroup mode** | **`cgroup2fs` (unified v2)** — controllers: `cpuset cpu io memory hugetlb pids rdma misc dmem` | **Exit criterion PASSED** (DEC-015: cgroup v2 required) |
| systemd | 259 (259.5-0ubuntu3.4) | — |
| Existing stack | docker: none · nginx: none · caddy: none · go: none · node: none · php: none · psql/mysql/mariadb: none | Clean host — no conflicts for installer preflight (WP-04) |
| Listening ports | 22 only | No port 80/443/25/53 conflicts |
| Network | Public IPv4 95.211.164.33/23 + global IPv6 on eth0 | Public IP present → hosting possible; PTR/port-25 checks still open |
| Uptime / load | 4.5 h, load ~0.1 | Idle |

## Status against the "environment established" exit criteria

| Criterion (owner ruling §4) | Status |
|---|---|
| cgroup v2 on the host (`stat -fc %T` → `cgroup2fs`) | ✅ **PASSED** on this host |
| One-command destroy/rebuild per profile | ⬜ Pending — profiles (Nano/Starter/Standard as code) not yet defined |
| Fault injection runnable (process kill, disk full, network partition) | ⬜ Pending — harness not yet built |
| Golden snapshot per profile | ⬜ Pending |

## Notes for the plan

1. The host matches the consolidated spec's dev reference node **exactly** (Ubuntu 26.04.1, 4 vCPU, ~6 GiB, ~100 GB) — it can serve as the Standard-class integration host immediately.
2. Ubuntu **26.04.1 remains "experimental until tested"** per ADR-017/SEC8-D13; the 24.04 profile (Nano/Starter VM definitions) is still the first-class tested target per the test-environment ruling.
3. No swap is configured — deliberate choice to record: memory-pressure benchmarks (QA-NANO-014/015) must set swap policy explicitly per run or results are meaningless.
4. Port 25 reachability and PTR record were **not** probed (outbound-impact checks deferred); they gate email capability detection (AGT-014b) later.
