# NexusPanel — Master Product & Engineering Specification
## Volume 06 — Privileged Server Agent & Hosting Service Engineering

**Version:** 1.0 consolidated | **Status:** Proposed engineering contract; review and validation required | **Date:** 2026-11-06  
**Authority:** Volume 01 Master PRD; Volumes 02–05; authoritative 29-area scope.  
**Source basis:** Uploaded Stage 7 Privileged Server Agent & Service Management v2.0, Stages 1–6 v2.0 and consolidated Volumes 01–05.  
**Operational boundary:** Documentation only. No commands have been executed on the Leaseweb VPS and no hosting service has been installed.

## 0. Consolidated engineering contract

This is the single normative specification for the NexusPanel host-side execution boundary and hosting service adapters. It integrates the substantive design and acceptance catalog from Stage 7 v2.0, replacing earlier deliberation markers with the decisions below. **MUST** means a release obligation for the feature concerned, **VERIFY** means an assumption requiring measured evidence, and **APPROVAL** means a decision requiring product/security signoff. A proposed contract is not an implemented capability.

### 0.1 Binding decisions and reconciled boundaries

1. **Agent and control plane:** Both are written in Go; the web UI is a static SPA. The privileged agent executes only typed, allowlisted operations, never unrestricted shell from the browser, extension or AI.
2. **Privilege separation:** The agent supervisor runs without root where feasible. A narrow privileged helper owns only the host mutations requiring elevation; local IPC uses a protected Unix socket and verified peer credentials. Socket file permissions alone are insufficient authorization.
3. **Transport profiles:** Local standalone v1 uses authenticated Unix-socket RPC with operation binding, replay protection, authorization rechecks and audit. Remote agent support requires enrollment, mutual authentication and encrypted transport; it must not inherit implicit trust from the local transport.
4. **Policy ownership:** The control plane authenticates and authorizes the user, validates entitlements and requested capability, and issues a bounded operation. The agent independently validates origin, host scope, expiry, revisions, action schema and local safety conditions before execution. Capability checks must not leak resource details to unauthorized users.
5. **Resource-aware Nano:** 512 MB RAM is a **target profile to benchmark**, not a guarantee that Docker, databases, mail, PHP or all services can run. Nano's baseline is static serving, essential management, minimal monitoring and remote backups where feasible. Sub-1-vCPU environments are best-effort subject to provider throttling.
6. **Workload scheduling:** cgroup v2 resource controls are profile-aware. CPU weight provides fairness; batch CPU caps and memory ceilings are policy tools, not a promise that every process has an independent cgroup or that a cap reserves RAM. Protect kernel, SSH and the control plane against resource exhaustion.
7. **Isolation:** Docker is preferred for eligible tenant workloads on capable profiles; user namespaces, capability drops, seccomp, mount/network policy and resource limits remain mandatory. Docker alone is not a complete multi-tenant security boundary. Nano's OS-user fallback requires explicit symlink, path, permissions and process-isolation tests.
8. **Universal service abstraction with specialized adapters:** `ServiceTemplate` standardizes catalog, prerequisites, deployment, health and inventory; databases, mail, DNS, object storage and backups still have typed lifecycle, integrity, access and recovery contracts. A signed template is not permission to execute arbitrary privileged code.
9. **Durable operations:** Every host mutation binds a tenant, resource, desired revision, operation ID, idempotency key, policy context, approval evidence where required, execution deadline and observable result. Unknown completion triggers reconciliation rather than blind retry.
10. **DNS/email phased scope:** v1 DNS guidance and external SMTP must not be presented as local authoritative DNS or self-hosted mail. These remain in the 29-area product scope for subsequent releases.
11. **Payment boundary:** The host agent does not process gateway credentials or infer subscription settlement. Commercial entitlements arrive as authenticated, scoped control-plane decisions; payment verification is governed by Volume 05.
12. **AI and terminal:** AI cannot bypass policy or invoke a generic root shell. Scoped developer execution is distinct from privileged host administration and is unavailable on Nano unless a safe profile is proven.
13. **ServiceTemplate category:** `serving`, `batch`, and `infrastructure` are governed classifications. Reclassification requires a reviewed migration rather than an ordinary untracked edit.
14. **Operational honesty:** Service status distinguishes requested, pending, observed, failed, degraded and unknown. A successful RPC is not proof that a website, certificate, restore or database is healthy.
15. **All 29 areas:** The agent may supply a host adapter, a metrics collector, an external-service boundary, or no local action for an area; this is documented rather than implying all features install on each VPS.

### 0.2 Mandatory operation contract

| Field | Purpose | Validation |
|---|---|---|
| `operation_id` | Durable control-plane operation | Unique and traceable |
| `idempotency_key` | Deduplicate retries | Bound to actor, target, action and payload digest |
| `tenant_id`, `account_id` | Ownership and delegation | Resolve against authorized target scope |
| `host_id`, `resource_id` | Exact execution target | Reject mismatch, stale ownership or unmanaged resource |
| `action`, `schema_version` | Typed allowlisted method | Reject unknown verbs and incompatible versions |
| `desired_generation` | Optimistic concurrency | Reject stale revisions; reconcile on ambiguous completion |
| `capability_decision` | Prerequisites and profile | Recheck current host conditions and limits |
| `approval_reference` | Elevated-risk approval | Bind approver, scope, action, expiry and payload |
| `issued_at`, `expires_at`, `nonce` | Freshness and replay defense | Enforce skew tolerance, expiry and replay cache |
| `payload_digest` | Integrity of exact requested change | Verify before preflight and commit |
| `audit_correlation_id` | Evidence linkage | Emit start, checkpoint, completion or failure |

### 0.3 Host operation lifecycle

`requested → validated → authorized → queued → dispatched → preflight → executing → verifying → succeeded` is the nominal path. Explicit alternative states include `approval_required`, `waiting_agent`, `yielded`, `retryable_failure`, `permanent_failure`, `partially_applied`, `compensating`, `cancelled`, and `unknown`. Any post-dispatch timeout is **unknown**, not automatically failed. Only a verified observed outcome may be reported as successful.

### 0.4 Host security invariants

- No direct browser-to-root transport; no arbitrary `bash -c` escape hatch in typed operations.
- No tenant-controlled file path may escape a validated ownership root through `..`, symlinks, bind mounts or race conditions.
- Secrets never appear in command arguments, process listings, debug logs, audit payloads or status responses.
- Configuration changes use parse/validate/stage/atomic activate/health-check/rollback where safe; irreversible effects require explicit compensating plans.
- Agent upgrades must be signed, compatible, reversible where possible and protected from downgrade attacks.
- Security-sensitive operations preserve an independent emergency recovery path for SSH and networking.
- Every module's resource budgets and dependency claims must be verified on the supported OS/kernel/service versions.

---

## 1. Product and engineering goals

1. An ordinary VPS owner can deploy and manage a website without understanding systemd, Linux permissions, TLS, or reverse proxies.
2. A hosting provider can manage thousands of segregated resources across multiple servers through the same governed operation model.
3. Resellers and customers receive delegated capabilities without acquiring operating-system root privileges.
4. All mutating host actions are typed, policy-checked, bounded, observable, auditable, and retry-safe.
5. An interrupted control plane, disconnected server, failed service reload, or expired certificate is diagnosable and recoverable.
6. Existing unmanaged applications and configuration are detected and protected from accidental takeover.
7. AI may recommend and submit operations through normal authorization and approval flows, but cannot bypass policy or run arbitrary privileged shell commands.


8. **Run on VPS infrastructure from 512 MB to multi-node clusters.** The agent adapts its adapter set to the detected capability profile. On Nano, it does not install or manage PHP/Node/MariaDB/Docker; it manages Nginx static serving and capability detection.
9. **Honest capability reporting.** The agent's capability detection is the source of truth for what the installer and UI display. A capability that is unavailable is reported with reason and resource delta, not silently disabled.
10. **Resource enforcement via cgroup v2.** Every workload has a declared `cpu.weight`, `memory.max`, and (for batch) `cpu.max`. `cpuset` is never used. Idle CPU is never wasted.
11. **Docker is the primary isolation primitive where Docker fits.** The agent validates container safety and rejects privileged mode, host PID/network, dangerous capabilities.
12. **ServiceTemplate is the universal service shape.** The agent does not have a `nextcloud.*` operation catalog. It has a `service.*` catalog. Adding a new template is a database insert, not an agent upgrade.
13. **The agent never sees payment credentials.** Payment is a control-plane state machine. The agent has no `payment.*` operations. This is a key security boundary.

### 1.1 Non-goals and prohibited shortcuts

- No root SSH key in the browser, database, or AI prompt context.
- No generic `execute(command: string)` privileged RPC endpoint.
- No `sudo bash -c` generated from user-controlled strings.
- No silent overwrite of administrator-owned Nginx, SSH, DNS, or firewall configuration.
- No assumption that all server services share the same trust boundary.
- No claim that cgroups alone provide strong hostile-tenant isolation.
- No assumption that backups or rollbacks are atomic across unrelated systems.


- **No `cpuset` (core pinning).** It reserves cores, which wastes idle CPU. This is ADR-005a. The agent rejects any operation that specifies `cpuset`.
- **No `--privileged`, host PID/network, or dangerous Linux capabilities in containers.** The agent rejects these configurations at validation time. This is AGT-108/AGT-109 restated as an agent contract.
- **No unbounded `cpu.max` on batch workloads.** Batch workloads get `cpu.max = 50% of cores` to prevent runaway. Serving workloads get no `cpu.max`.
- **No payment credential in agent environment.** The agent does not read M-Pesa or Stripe credentials. Payment is control-plane only.
- **No bespoke per-service operation endpoints.** All services use `service.*`. New templates are data.

## 2. Reference deployment topologies

### 2.1 Single-server topology (first development/integration environment)

```text
Browser / API client
       |
       v HTTPS
NexusPanel Static SPA (Nginx)
       |
       v HTTPS
Control Plane API (Go)
       |
       v Unix socket + SO_PEERCRED
Go Host Agent (root)
       |
       +--------------+------------------+--------------+
       |              |                  |              |
   Nginx adapter  Docker adapter   Service adapter  cgroup adapter
       |              |                  |              |
   Sites / TLS   Containers/cgroups  DB/DNS/Mail    /sys/fs/cgroup/
```


### 2.2 Multi-server provider topology

```text
UI / Public API
       |
Control Plane + Identity + Policy + Job Scheduler (Go)
       |
       +--> Server A Agent [web/app]        (Unix socket on each node)
       +--> Server B Agent [databases]
       +--> Server C Agent [mail]
       +--> Server D Agent [authoritative DNS]
       +--> Server E Agent [backup/restore]
       |
       +--> External provider adapters [DNS, CDN, PTR, billing]
```

Server roles are capabilities, not hard-coded classes. A node may support several roles; placement must account for capacity, failure domains, compliance, service compatibility, and tenant policy. Independent secondary authoritative DNS requires an independently reachable failure domain; the first VPS alone cannot provide meaningful DNS redundancy.


### 2.3 Connectivity model

- **AGT-001:** Agent connectivity MUST use authenticated, encrypted transport; local Unix sockets MAY be used with OS peer-credential validation.
- **AGT-002:** Remote agents MUST use mutually authenticated TLS or an equivalent authenticated transport with rotating, revocable machine identities.
- **AGT-003:** An agent MUST NOT accept privileged operations from arbitrary internet clients; network exposure and allowlisting are explicit deployment decisions.
- **AGT-004:** Offline agents MUST be shown as unavailable; queued operations MUST not be represented as completed.
- **AGT-005:** Enrollment MUST bind agent identity to server inventory and record operator approval, expiry, and revocation.


- **AGT-001a (v1):** The local agent listens on `/run/nexuspanel/agent.sock` (mode 0660, owner `nexuspanel:nexuspanel`). It validates `SO_PEERCRED` on every connection — the connecting process must be the panel API's UID. No other process can connect.
- **AGT-002a (v1):** There is no mTLS in v1. The agent auto-enrolls on install. There is no enrollment token, no CSR, no certificate. The peer credential check is the authentication.
- **AGT-005a (v1):** The agent's identity is its OS process UID. The panel API is the only process that runs as that UID. Peer credential check is sufficient for a single-machine deployment.

## 3. Trust boundaries and principal model

| Principal | Capabilities | Must not be able to do |
|---|---|---|
| Browser/user | Authenticate, request authorized actions, inspect permitted status | Reach root agent or submit raw privileged commands |
| Control plane | Authorize, plan, approve, schedule, record operations | Bypass agent validation or directly edit host-managed configs |
| Job worker | Deliver signed/authorized operation envelopes | Invent permissions or substitute resource ownership |
| Agent supervisor | Verify envelope, dispatch typed operation, report outcome | Accept arbitrary commands or cross-tenant path requests |
| Privileged helper | Perform narrowly scoped host action | Perform unrelated OS administration |
| Workload owner | Operate own files/apps within quota | Access other tenants or control-plane secrets |
| External integration | Act under explicit scoped token/credential | Inherit full provider/admin privileges |
| AI assistant | Diagnose, draft plan, request approved action | Bypass approvals, retrieve raw secrets, or run privileged shell |

**AGT-006:** Every operation MUST carry an authenticated initiating principal, organization/tenant, resource ID, effective permissions, policy decision reference, and immutable audit correlation ID. **AGT-007:** Authorization MUST be enforced in the control plane and independently validated by the agent for target, operation type, bounds, and expected resource ownership; agent-side checks do not replace full RBAC evaluation. **AGT-008:** All effective privileges MUST be least-privilege, including service-specific OS identities and delegated file access.


- **The capability gate is a first-class authorization step at the agent.** Before executing any workload-creating operation, the agent re-checks the capability record. If the capability has become unavailable (e.g., hardware downgraded, disk full), the operation is rejected with `CAPABILITY_UNAVAILABLE` — even if the control plane approved it. This is defense in depth.
- **The agent never sees payment credentials.** Payment is control-plane only. The agent has no M-Pesa/Stripe access.
- **The agent's Unix socket has no authentication beyond peer credentials.** This is deliberate — the local agent trusts the local panel process. If the panel process is compromised, the agent is compromised. This is an accepted v1 tradeoff; mTLS is v2+.
- **The `cgroup` adapter is a trusted boundary.** It writes to `/sys/fs/cgroup/`. It validates that the requested config is compatible with the workload category and rejects anything else.

## 4. Agent process design and privilege separation

### 4.1 Recommended components

- **Agent supervisor:** Go binary running under a dedicated unprivileged service account; handles transport, registration, heartbeats, request validation, task orchestration, and telemetry.
- **Privileged executor:** Minimal root-owned helper, invoked through a constrained local IPC interface with operation-specific schemas; no arbitrary shell feature. The helper must validate again rather than trust the supervisor blindly.
- **Service adapters:** Versioned modules for Nginx, systemd, runtimes, files, databases, DNS, mail, firewall, backup, containers, caching, monitoring, and migrations.
- **Local operation journal:** Durable record of accepted task IDs, transitions, applied revisions, and completion proofs, allowing safe restart and duplicate detection.
- **Reconciler:** Periodically compares approved desired state with discovered actual state, reports drift, and performs only permitted low-risk repairs automatically.
- **Inventory collector:** Reports installed packages, supported runtime versions, service health, disk mounts, network addresses, capacity, and module capabilities.
- **Capability detector (v2):** Runs the 5-step detection logic, produces the capability record, signs it, and reports it to the control plane.

### 4.2 Service account and sandboxing

**AGT-009:** Supervisor MUST run with minimal Linux capabilities and a hardened systemd unit. **AGT-010:** Privileged executor MUST expose a strict allowlist of typed operations; avoid persistent unrestricted root RPC. **AGT-011:** Every adapter MUST declare required capabilities, OS/package dependencies, paths owned, rollback support, supported Ubuntu releases, and a health probe. **AGT-012:** Sensitive process environment variables and log messages MUST be redacted. **AGT-013:** Agent startup MUST fail closed when its identity or trusted policy material is invalid.

Proposed hardening controls include `NoNewPrivileges`, restricted writable paths, protected home/system directories, private temporary directories, syscall restrictions, and per-service Linux capabilities. These settings require compatibility testing; a blanket setting can break package management, network configuration, or service supervision. The privileged executor needs a separately reviewed confinement profile.


- **The agent runs as a dedicated UID** (`nexuspanel-agent`). The panel API runs as a separate UID (`nexuspanel-api`). The Unix socket is mode 0660 owned by `nexuspanel:nexuspanel-agent`, so only those two UIDs can connect. This is the authentication mechanism.
- **The agent has no network access.** It communicates only over the local Unix socket. No outbound HTTP, no inbound TCP. This is enforced by the systemd unit: `PrivateNetwork=yes` with the socket bind-mounted.
- **The agent's filesystem is read-mostly.** It can read `/etc/nginx`, `/etc/caddy`, `/var/lib/docker`, `/sys/fs/cgroup`, and the workload directories. It can write only to its own managed paths. `ProtectSystem=strict`, `ReadWritePaths=` for specific directories.

### 4.3 Filesystem layout (proposed)

```text
/opt/nexuspanel/agent/                 # versioned, signed agent binaries
/etc/nexuspanel/agent/                 # root-controlled agent config / trust
/var/lib/nexuspanel/agent/             # journal, checkpoints, inventory cache, capability record
/var/log/nexuspanel/agent/             # sanitized logs (or journald)
/run/nexuspanel/agent.sock             # protected local RPC socket (mode 0660)
/etc/nginx/nexuspanel-managed/         # agent-owned configuration fragments
/srv/nexuspanel/tenants/<tenant-id>/   # policy-managed hosting data roots
/sys/fs/cgroup/nexuspanel.slice/       # cgroup root for all workloads
```

Paths are design proposals; the installer must determine actual compatible locations. No tenant-controlled input is interpolated into privileged path names.


- **`/sys/fs/cgroup/nexuspanel.slice/<workload-id>.scope/`** — the cgroup for each workload. The agent creates these on workload creation and validates the applied values by reading them back.

## 5. Host capability discovery and compatibility

**AGT-014:** Agent MUST publish an inventory document with OS release/kernel/architecture, package manager, disk layout, networking, available services, adapter versions, and capacity. **AGT-015:** Each capability MUST be reported as `available`, `installed`, `unsupported`, `degraded`, or `requires_setup`, with a reason and prerequisite list. **AGT-016:** The control plane MUST reject unsupported operations before scheduling them. **AGT-017:** OS upgrades and module upgrades MUST use a tested compatibility matrix and preflight checks.

Initial target: Ubuntu 26.04.1 LTS on the Leaseweb VPS supplied by the user; compatibility with specific versions of Go, Node.js, Nginx, PostgreSQL, MariaDB, PHP, Docker, PowerDNS, Postfix/Dovecot, and supporting packages must be verified rather than assumed. Additional distributions may be added through explicitly tested adapters.


**AGT-014a:** The agent runs the following detection on install and on `POST /servers/{id}/capabilities/rescan`:

**Step 1 — Detect hardware and blockers.**
- CPU: cores, model, flags (AVX2 matters for some workloads), architecture
- RAM: total, available after OS baseline, swap
- Disk: total, available, filesystem, I/O speed
- Network: public IPv4, IPv6, reverse DNS, port 25 reachability, port 53, port 80/443
- OS: distro, version, kernel, cgroup v2 availability
- Existing state: other panels, other web servers, existing databases, existing Docker

**Step 2 — Compute floor check.** Sum floors of the panel + the workloads the user selects. If sum > available RAM → refuse.

**Step 3 — Compute typical fit.** Sum typicals. If > 80% of available RAM → warn. If > 100% → warn harder.

**Step 4 — Compute CPU overcommit.** Sum typical CPU. Divide by physical cores. If ratio > 3× → warn. If > 5× → refuse.

**Step 5 — Compute ceiling fit.** Sum RAM ceilings. If > physical RAM × 1.2 → refuse.

**Step 6 — Compute disk.** Per-user quotas. If sum of disk quotas > physical disk × 0.9 → reduce quotas or refuse.

**AGT-014b:** Non-resource blockers are detected independently. Port 25 blocked → email unavailable regardless of RAM. No public IP → hosting impossible. No PTR → email functional but deliverability fails. No cgroup v2 → agent refuses to install.

**AGT-014c:** The capability record is signed by the agent (HMAC with a key in `/etc/nexuspanel/agent/trust.key`). The control plane verifies the signature before displaying it. A tampered capability record is rejected.

**AGT-014d:** The capability record is stored in `/var/lib/nexuspanel/agent/capability.json` with a digest. The control plane stores the same digest in `inventory.capabilities.source_digest`. If the agent's digest and the control plane's digest diverge, a reconciliation job triggers a rescan.

## 6. Canonical operation envelope

```json
{
  "schema_version": "agent.operation.v1",
  "operation_id": "op_01JEXAMPLE",
  "idempotency_key": "org_42:create-site:request_7",
  "correlation_id": "req_01JEXAMPLE",
  "requested_by": {"principal_id": "usr_42", "kind": "human"},
  "tenant_id": "org_42",
  "resource": {"kind": "website", "id": "site_93", "server_id": "srv_1"},
  "action": "website.provision",
  "desired_revision": 7,
  "policy_decision_id": "pdp_173",
  "approval_id": "appr_55",
  "plan_digest": "sha256:<canonical-plan-hash>",
  "not_before": "2026-11-06T15:00:00Z",
  "expires_at": "2026-11-06T15:10:00Z",
  "preconditions": {"current_revision": 6, "server_capability": "nginx"},
  "parameters": {"site_template_id": "static-site-v1", "hostname_ref": "host_2"},
  "credential_ref": null,
  "authenticity": {"method": "mTLS+signed-envelope", "key_id": "control-2026-11"}
}
```

Illustrative only; `authenticity` must be implemented using a canonical signed serialization with anti-replay protections and binding to the actual connection identity. Never trust a user-provided `requested_by` field without control-plane verification. **AGT-018:** Agent MUST validate schema version, server binding, allowed action, expiry, replay nonce/ID, plan digest, expected revision, ownership constraints, and operation-specific parameter bounds. **AGT-019:** Agent MUST reject unrecognized fields where they could alter security-sensitive interpretation. **AGT-020:** The same idempotency key and payload MUST produce the same recorded outcome or resume the original operation; same key with changed payload MUST fail. **AGT-021:** Secrets MUST be referenced through short-lived, purpose-bound credentials; do not embed secrets in ordinary operation records.


**New fields:**
- **`workload_category`** — required. `serving` | `batch` | `infrastructure`. Set at submission; immutable. The agent uses it to derive cgroup configuration if the command creates or modifies a workload.
- **`cgroup_config`** — for workload-creating commands. Validated by the agent. `{cpu_weight, cpu_max_millicores, memory_low_bytes, memory_max_bytes, pids_max}`. `cpuset` is rejected.
- **`capability_decision_id`** — the decision ID from the capability gate. The agent re-checks the capability at execution time.

**New rules:**
- **AGT-018a (v2):** The agent MUST re-evaluate the capability gate at execution time, not just submission time. A capability that was available at submission but is now unavailable (hardware downgraded, disk full) MUST cause the operation to be rejected.
- **AGT-021a (v2):** The `credential_ref` field points to the secret store. The agent resolves it at execution time. For workload operations, the credentials are injected into the workload's environment, never logged, never returned.
- **AGT-021b (v2):** The `authenticity` field is `null` in v1 (Unix socket peer credential check is the authentication). It is populated in v2+ (mTLS).

### 6.1 Operation response

```json
{
  "operation_id": "op_01JEXAMPLE",
  "state": "verifying",
  "progress": {"phase": "nginx_health_check", "completed_steps": 4, "total_steps": 5},
  "result": null,
  "error": null,
  "agent_id": "agent_1",
  "observed_revision": 7,
  "last_update_at": "2026-11-06T15:03:00Z",
  "workload_category": "serving",
  "cgroup_applied": {"cpu_weight": 100, "memory_max_bytes": 419430400, "cpu_max_millicores": null}
}
```


- **`workload_category`** — mirrors the request.
- **`cgroup_applied`** — the values actually applied, read back from `/sys/fs/cgroup/`. This is how the control plane verifies the cgroup config was applied.
- **`capability_decision`** — `{state: "available", reason: null}` or `{state: "unavailable", reason: "..."}`. The agent reports its re-check.
- **`queue_position`** — for batch operations, the position in the batch queue.

## 7. Operation lifecycle and state machine

```text
requested -> validating -> planned -> awaiting_approval -> queued
                                                       |
                                                       v
                                                     leased
                                                       |
                                                       v
                                                    preparing
                                                       |
                                                       v
                                                     applying
                                                       |
                                                       v
                                                    verifying
                                                    /      \
                                              succeeded   failed
                                                        /     \
                                                  rolling_back  manual_recovery
                                                      |
                                               rolled_back / rollback_failed
```

Other terminal states: `rejected`, `cancelled_before_apply`, `expired`, and `unsupported`. `interrupted` is a recoverable nonterminal condition, not proof of failure or success.

**AGT-022:** Every mutation MUST support a preflight and human-readable plan, including affected resources, downtime risk, permissions, expected changes, and rollback limitations. **AGT-023:** Plans MUST be invalidated when relevant preconditions or resource revisions change. **AGT-024:** High-risk actions MUST require approval according to policy, including root SSH changes, firewall changes, mass deletion, DNSSEC key changes, bulk mailbox deletion, restore-overwrite, and server removal. **AGT-025:** Operations MUST record progress and terminal evidence. **AGT-026:** Cancellation MUST only be offered at defined safe points; once external side effects begin, cancellation may become a compensating operation. **AGT-027:** Failed verification MUST not be labeled successful merely because a command exited with code zero.


- **Batch operations are a distinct state.** Batch operations add `queued → waiting → running → yielding → running → verifying → succeeded`. The `yielding` state is when the operation pauses to give CPU to serving workloads. This is not failure; it's the expected behavior.
- **The batch state machine is:**
  ```
  queued → leased (when the agent is ready)
        → waiting (if the batch queue has higher-priority items)
        → running
        → yielding (if cpu.stat.nr_throttled elevated on serving workloads)
        → waiting (yielded, will resume)
        → running (resumed)
        → verifying → succeeded
  ```
- **AGT-027a:** The agent MUST verify cgroup applied values by reading them back. `cpu.weight` and `memory.max` are read from the cgroup filesystem after apply. If they don't match, the operation fails.
- **AGT-027b:** The agent MUST verify container safety (no `--privileged`, no host PID/network) by inspecting the running container's config before considering it healthy. A container that somehow got a privileged flag is stopped and the operation fails.

### 7.1 Transaction and compensation model

The agent uses a saga-like sequence of checkpoints rather than claiming distributed ACID transactions across Nginx, systemd, DNS, and storage. Each step records `prepare`, `apply`, `verify`, and, where possible, `compensate`. Rollback may be partial and must be represented honestly. For irreversible actions, require explicit backup or snapshot and a recovery plan before applying.

### 7.2 Operation leases and crash recovery

**AGT-028:** Leases MUST expire and be reclaimable only after checking the local journal and observed host state. **AGT-029:** After agent restart, incomplete operations MUST be reconciled rather than blindly replayed. **AGT-030:** Resource-scoped locks MUST prevent conflicting changes (e.g., two deployments to one site or simultaneous firewall policy updates). **AGT-031:** Long-running operations MUST send progress heartbeats and preserve resumable checkpoints where technically feasible.

## 8. Desired-state reconciliation and drift

- **Desired state:** approved control-plane resource revision.
- **Observed state:** live host state from adapter inspection.
- **Last applied state:** signed, journaled revision and manifest.
- **Drift:** observed state differs from the approved manifest or expected service health.

**AGT-032:** Drift MUST be classified as benign, recoverable, destructive, externally managed, or security-sensitive. **AGT-033:** The agent MUST NOT overwrite unmanaged or externally owned configurations without explicit adoption. **AGT-034:** Low-risk repair MAY be automatic only if policy permits; destructive or security-sensitive drift requires review. **AGT-035:** Reconciliation MUST be rate-limited and stable under repeated failure (backoff, circuit breakers, quarantine).


- **cgroup drift detection.** The reconciler reads `/sys/fs/cgroup/<workload>/` values and compares to the approved config. A workload whose `memory.max` was changed outside the panel is flagged as `security-sensitive drift` (someone bypassed the panel).
- **Capability drift detection.** The reconciler re-checks capability floor conditions. If the physical RAM has been reduced (VPS provider downgrade), the capability record is re-evaluated. A capability that became unavailable is flagged; running workloads that exceed the new capability set are flagged for upgrade or migration.
- **Batch queue reconciliation.** After agent restart, the batch queue is reconstructed from the journal. Operations that were `running` when the agent crashed are marked `interrupted` and re-queued if safe, or flagged for manual review if not.

## 9. Service adapter interface

```text
Adapter:
  Describe() -> capability, version, prerequisites, owned_paths
  Discover() -> installed_services, instances, health
  Validate(request, context) -> violations, warnings
  Plan(request, observed_state) -> steps, diff, risk, digest
  Prepare(plan) -> checkpoint, acquired_locks
  Apply(plan, checkpoint) -> step_receipts
  Verify(plan) -> health_evidence, observed_revision
  Compensate(plan, checkpoint) -> recovery_receipt
  Reconcile(desired, observed) -> drift_report, safe_repair_plan
  Uninstall(managed_resource) -> explicit cleanup_plan
```

**AGT-036:** Adapters MUST have explicit version compatibility and contract tests. **AGT-037:** Adapters MUST declare whether they support dry run, reversible apply, online reload, tenant isolation, backup, and external-provider coordination. **AGT-038:** All managed files MUST have ownership markers, checksums, and atomic write behavior where supported. **AGT-039:** Configuration validation MUST run before reload/restart. **AGT-040:** Service restarts MUST be impact-assessed; prefer graceful reload when safe.


Every adapter MUST declare its **resource profile**: floor / typical / ceiling per resource, workload category, prerequisites. This is what the installer reads to compute capability fit.

| Field | Type | Meaning |
|---|---|---|
| `ram_floor_bytes` | int64 | Minimum RAM to function |
| `ram_typical_bytes` | int64 | Expected steady-state |
| `ram_ceiling_bytes` | int64 | Maximum under peak (enforced by `memory.max`) |
| `cpu_weight` | int | Fair-share under contention |
| `cpu_max_millicores` | int (nullable) | Only for batch; runaway protection |
| `disk_floor_bytes` | int64 | Minimum disk |
| `workload_category` | enum | `serving` / `batch` / `infrastructure` |
| `prerequisites` | map | `{docker: true, port_25: true, ...}` |

Adapters that cannot declare this are not shippable.

### 9.1 Managed configuration policy

Generate new configs into a temporary root-controlled staging directory; validate syntax and dependencies; atomically replace only files owned by NexusPanel; reload; run health checks; commit journal; preserve last known good manifest. For services that cannot support atomic replace, declare that limitation and use maintenance windows or compensating procedures.


- **The `service.*` adapter is the generic adapter.** It covers storage, transform, delivery, self-hosted apps, databases, cache, queue, search. The template record supplies the image, resources, ports, health check, credentials, export path. The adapter supplies the lifecycle.
- **The `cgroup.*` adapter is new.** It manages `/sys/fs/cgroup/` for every workload. Operations: `cgroup.create`, `cgroup.update`, `cgroup.verify`, `cgroup.remove`. It rejects `cpuset` and out-of-range values.
- **The `container.*` adapter is Docker-specific.** It creates, starts, stops, removes containers. It validates safety (rejects `--privileged`, host PID/network, dangerous capabilities). It reads container health checks and logs.

## 10. Host OS, package and service management (Area 02, 19)

Capabilities: package inventory, trusted repositories, service discovery, restart/reload, systemd unit templates, OS updates, kernel/reboot scheduling, hostname/time sync, capacity, network interfaces, IPAM coordination, PTR request workflow, and multi-server placement.

**AGT-041:** Package operations MUST use signed repositories and an allowlisted package set with version pinning and rollback feasibility recorded. **AGT-042:** OS upgrades MUST be separate high-risk workflows with snapshot/backup checks, maintenance windows, and console recovery instructions. **AGT-043:** The agent MUST NOT change SSH/firewall/network routes without a timed revert mechanism and a proven alternative management path. **AGT-044:** IP assignment and PTR changes MUST distinguish local OS configuration from provider-owned network/rDNS APIs; never pretend local DNS changes alter Leaseweb reverse DNS. **AGT-045:** Traffic shaping and bandwidth limits MUST identify interface, tenant attribution model, and interference risk.


- **cgroup v2 is required.** Ubuntu 22.04+ / 24.04+ is the OS baseline. The installer refuses to install on systems without cgroup v2. This is DEC-015.
- **The agent creates `/sys/fs/cgroup/nexuspanel.slice/` on install.** All workloads live under this slice. The panel and agent themselves live under `nexuspanel-api.slice` and `nexuspanel-agent.slice` with `cpu.weight=200` (infrastructure category).
- **Container runtime is Docker on Standard+.** The agent uses the Docker socket. It does not have direct access to `containerd` or `runc`. Docker is the abstraction.

## 11. Web server, websites and routing (Areas 03, 05, 24)

Nginx is the proposed initial reverse proxy; Apache may be an optional adapter. Managed resources include website, virtual host, hostname, redirect, document root, TLS policy, app upstream, cache policy, WAF attachment, and staging route.

**AGT-046:** Every hostname MUST be normalized, ownership-checked, and conflict-checked across all tenants before activation. **AGT-047:** Generated web-server configs MUST pass syntax validation, be tenant-isolated, and support graceful reload and rollback. **AGT-048:** Document roots MUST remain within assigned tenant storage and must reject traversal and unsafe symlinks. **AGT-049:** Reverse-proxy upstreams MUST use explicitly registered ports/sockets, with collision checks and health probes. **AGT-050:** Staging/cloning MUST keep secrets, cookies, indexing, mail delivery, and production data access segregated by policy.


- **On Standard+, websites are containerized.** The web server (Nginx) is the reverse proxy; each site runs in a container with its own web server (Nginx + PHP-FPM or Node or Python). The reverse proxy routes by hostname to the container. Docker is the isolation primitive.
- **On Nano, websites are static-only.** The agent configures Nginx to serve the site's document root directly. No containers, no PHP, no Node. This is honest: the Nano VPS cannot run a PHP stack.
- **The `container.*` adapter is used to create the site's container.** The website resource has `container_id` and `volume_id`. The reverse proxy has a route to `container_ip:port`.
- **The `cgroup.*` adapter configures the container's cgroup.** `cpu.weight=100` (serving), no `cpu.max`, `memory.max` set to the site's ceiling.

### 11.1 Example website provision plan

1. Validate tenant entitlement, hostname control, quota, capability fit, and server capacity.
2. Reserve resource identity, storage root, runtime binding, and proxy route.
3. Create site OS identity and directory with safe ownership/permissions.
4. **(v2)** Create the container with the correct image (static / PHP / Node), volume, cgroup config. Validate safety.
5. Stage Nginx configuration for the reverse proxy route and optional TLS challenge route.
6. Validate `nginx -t` or equivalent in an isolated plan step.
7. Commit config, graceful reload, probe HTTP/TLS, and record revision.
8. **(v2)** Verify the container health check passes, the cgroup values match, and the reverse proxy resolves.
9. On failure, restore previous managed config and remove only newly created owned resources; preserve diagnostics.

## 12. TLS/ACME certificate management (Area 05)

**AGT-051:** ACME challenges MUST support HTTP-01 and DNS-01 through approved adapters, with DNS-01 required for wildcard issuance. **AGT-052:** Private keys MUST remain access-controlled and never appear in logs, API payloads, or UI exports. **AGT-053:** Renewal MUST have jittered scheduling, retry/backoff, expiry thresholds, and failure alerts. **AGT-054:** HTTPS redirect and HSTS changes MUST show a preflight impact assessment; HSTS preload must never be enabled implicitly. **AGT-055:** Revocation, replacement, multi-host deployment, and certificate rollback MUST be explicit operations.


- **v1 uses HTTP-01 only.** Caddy's automatic HTTPS is the primary path. The certificate is issued by Let's Encrypt via Caddy's built-in ACME client.
- **DNS-01 and wildcard are deferred to v2+.** They require authoritative DNS, which is deferred.
- **PTR coordination is a provider integration.** The agent does not set PTR. It surfaces the requirement to the control plane, which guides the user to the provider's control panel.

## 13. File manager, editor, SFTP and FTP (Area 06)

Features: tree explorer, streaming upload/download, chunked transfer, atomic editor saves, revision/ETag checks, archive inspection/extraction, bulk operations, ACL/ownership management, SFTP identities, optional FTP compatibility adapter, and quota accounting.

**AGT-056:** All file operations MUST be scoped to an authorized root using directory-file-descriptor-relative operations and symlink-safe resolution; string-prefix path checks are insufficient. **AGT-057:** Archives MUST be screened for traversal, symlinks, decompression bombs, executable metadata, and quota overrun before extraction. **AGT-058:** Editor saves MUST be conflict-aware and atomic where possible. **AGT-059:** SFTP identities MUST be constrained to assigned tenant roots and authenticated with scoped credentials; plain FTP is disabled by default. **AGT-060:** Bulk delete/permission changes MUST provide preview, ownership boundaries, and recovery policy. **AGT-061:** File scanning/quarantine integrations MUST not expose quarantined content to unauthorized users.


- **The file manager operates on the site's container volume on Standard+.** The agent execs into the container (`docker exec`) or uses the volume mount. Path resolution uses `openat2` with `RESOLVE_BENEATH`.
- **On Nano, the file manager operates on the document root directly** (no container). OS-user isolation.
- **SFTP is the primary file-transfer path on Nano.** The panel manages the SFTP endpoint and per-site credentials. SSH shell is unavailable on Nano.
- **Terminal access is disabled by default** on Standard+. It's a separate permission. On Nano, terminal is unavailable (no container to exec into).

## 14. Application manager and language runtimes (Areas 07–08)

Initial adapter families: static, Node.js, Python, PHP-FPM, and one-click CMS templates. Version discovery is dynamic from approved runtime repositories. Workload execution uses tenant-specific identities, resource quotas, process supervision, controlled environment variables, and isolated build directories.

**AGT-062:** Applications MUST have explicit states: `draft`, `provisioning`, `running`, `degraded`, `stopped`, `updating`, `failed`, `removing`. **AGT-063:** Runtime versions and extensions MUST be validated against a compatibility catalog, including deprecation and security status. **AGT-064:** Build and start commands MUST be treated as tenant workload execution, not as privileged agent commands. **AGT-065:** Environment variables MUST support bulk paste/file import at the UI layer, but secrets MUST be parsed, redacted, encrypted at rest, and injected only into the target workload. **AGT-066:** Each app MUST have bounded CPU/memory/process/disk limits and non-root execution; stronger isolation profiles are required for untrusted tenants. **AGT-067:** Deployment rollback MUST reference an immutable prior artifact and explicitly account for non-reversible database migrations. **AGT-068:** One-click installers MUST use versioned, verified manifests with dependency checks and secure initial credential delivery.


- **The app runtime is a serving workload.** `cpu.weight=100`, no `cpu.max`, `memory.max` set. The runtime runs in a container on Standard+.
- **The build is a batch workload.** `cpu.weight=10`, `cpu.max=50% of cores`, `memory.high` throttling, `memory.max` hard cap. The build runs in a separate container with no host access.
- **The build queue is visible.** `POST /applications/{id}/deployments` creates a batch operation. It appears in `GET /operations?workload_category=batch`. The UI shows "Build queued (position 2), will run when server is idle."
- **The build yields to serving workloads.** Under site traffic, the build pauses between items. This is not failure; it's expected.

### 14.1 Node.js UX-to-agent mapping

A user's choices—production/development, app root, Node version, startup file, build command, repository branch, environment variables, and hostname—become a validated `application.deploy` plan. The agent creates a scoped runtime environment, installs dependencies as the tenant identity, registers a supervised process, configures an Nginx upstream, runs a health check, and publishes logs and deployment status. The agent never executes arbitrary build commands as root.


- **`.env` bulk paste is a first-class agent operation.** `app.environment.import` takes a full `.env` text, parses it (respecting quotes, escapes, comments), validates keys, stores secrets encrypted, injects only at container start. The parsed keys are visible to the user in the UI; the values are masked. The agent never logs values.

## 15. Git, CI/builds and deployments (Area 11)

**AGT-069:** Git credentials MUST be scoped, encrypted, revocable, and isolated from other tenants. **AGT-070:** Webhook events MUST be signature-verified and replay-protected at the control plane. **AGT-071:** Builds MUST run with constrained network/CPU/memory/disk/time permissions and must not mount host root or privileged sockets. **AGT-072:** Deployment plans MUST include branch-to-environment mapping, artifact digest, health gate, approval policy, and rollback target. **AGT-073:** Push-to-deploy MUST be opt-in with protected branch and deployment policy support. **AGT-074:** Logs and artifacts MUST have retention, access control, and secret redaction.


- **Builds run in a build container with `--network=none` by default.** If the build needs network (npm install, pip install), it's explicitly enabled for the build container only, with an egress allowlist.
- **The build container is `--rm`** — removed after the build. No build cache leaks between tenants.
- **Builds are batch.** Even if the build is fast, it's classified as batch and yields under contention. This is what keeps site traffic responsive during a deploy.

## 16. Customer database service adapters (Area 09)

MariaDB/MySQL and PostgreSQL adapters provision logical databases, users, scoped privileges, backups, restores, connection policies, and quotas. The NexusPanel metadata PostgreSQL database is **not** a customer database.

**AGT-075:** Database users MUST receive only the privileges required for their database; avoid shared superuser credentials in applications. **AGT-076:** Database name/user allocation MUST be collision-safe and tenant-bound. **AGT-077:** Remote access MUST be disabled by default and explicitly controlled by network policy, TLS, and source restrictions. **AGT-078:** Admin UI integrations (phpMyAdmin/pgAdmin-like) MUST be isolated and access-brokered; never expose database admin services unauthenticated. **AGT-079:** Backup/restore MUST declare consistency method, recovery point, schema compatibility, and overwrite risk. **AGT-080:** Quotas MUST define measurement and enforcement method; where a database engine cannot hard-enforce a quota, the UI MUST disclose that limitation.


- **v1: databases are unavailable.** On Nano, the panel's SQLite metadata DB uses the disk. On Starter, a small database is borderline. On Standard+, managed databases ship in v2.
- **When v2 ships, each database is a container.** `POST /service-instances` with `template_id = "database.mariadb"` or `"database.postgres"`. The database-specific extensions (`logical_databases`, `database_principals`, `database_grants`) reference `service_instances.id`.
- **The database container has its own cgroup.** Serving workload. `cpu.weight=100`, `memory.max` set.
- **The database admin UI (phpMyAdmin-like) is also a `ServiceTemplate` instance.** Not embedded in the control plane.

## 17. DNS and domain service adapters (Area 04)

Three operating modes: `external_manual`, `external_api`, and `local_authoritative`. Local authoritative DNS is an optional role, not a prerequisite to hosting a website. Support zone lifecycle, record sets, DNSSEC, nameserver delegation, primary/secondary synchronization, glue checks, propagation verification, and PTR coordination.

**AGT-081:** DNS edits MUST use preflight validation for record syntax, conflicts, TTL, CNAME exclusivity, zone authority, and DNSSEC implications. **AGT-082:** External-manual mode MUST never claim to have applied registrar or external DNS changes. **AGT-083:** DNSSEC signing/rollover MUST use a dedicated high-risk workflow with registrar DS coordination and recovery steps. **AGT-084:** Primary/secondary health MUST be verified independently; one-server dual nameserver aliases do not satisfy redundancy. **AGT-085:** PTR changes MUST route to the owner of the IP reverse zone/provider integration. **AGT-086:** DNS change verification MUST distinguish accepted configuration, authoritative response, and recursive propagation.


- **Authoritative DNS is v2+.** The agent has no `dns.zone.*` operations in v1.
- **v1 agent operation:** `dns.guidance.check` — the agent runs `dig` against public DNS resolvers and reports the observed state. This is what the control plane uses for propagation polling.
- **`PTR` coordination:** the agent reports the observed PTR (via `dig -x`) and the control plane compares to what the user requested. The agent does not set PTR.
- **HTTP-01 SSL validation requires the A record to point at the VPS.** The agent checks `dig +short <domain>` before attempting HTTP-01. If the record doesn't match the VPS IP, the SSL operation is deferred.

## 18. Email hosting and deliverability (Areas 10, 25)

Optional mail role: SMTP ingress/egress, IMAP, mailbox storage, webmail, aliases, forwarding, catch-all, quotas, migration, queue management, DKIM signing, spam/virus filtering, outbound throttling, reputation monitoring, and per-domain policies. Candidate components include Postfix, Dovecot, Rspamd, and webmail integration, subject to compatibility and licensing review.

**AGT-087:** Mail services MUST not be auto-enabled on every server. **AGT-088:** Domain onboarding MUST validate MX, SPF, DKIM, DMARC, TLS, reverse DNS, and provider port restrictions with clear ownership boundaries. **AGT-089:** Mailbox access MUST be tenant-scoped with secure credential lifecycle and storage encryption policy. **AGT-090:** Outbound mail MUST support rate limits, abuse detection, queue inspection, and emergency suspension. **AGT-091:** Catch-all and forwarding MUST be explicit opt-in due to abuse/deliverability risk. **AGT-092:** IMAP migration MUST use checkpoints, duplicate detection, progress reporting, and post-copy verification. **AGT-093:** Spam/quarantine access MUST be restricted and retention-controlled. **AGT-094:** Deliverability checks MUST distinguish configuration correctness from actual inbox placement, which cannot be guaranteed.


- **Self-hosted email is v2+.** The agent has no `mail.*` operations in v1.
- **v1 external SMTP is a control-plane configuration.** `POST /email/external-smtp` stores the provider credentials (Resend/Postmark/Mailgun/SES) in the encrypted secrets store. The agent never touches mail.
- **The agent may need to configure WordPress outbound mail** if the user has a WordPress site. This uses the WordPress site's `wp-config.php` — the agent updates the SMTP settings via the file manager API. The credentials are injected from the secrets store, never logged.

## 19. Backups, restores and migration (Areas 12, 22)

Support scheduled and on-demand snapshots, per-site files, database dumps/physical backups, mailbox backups, panel metadata, offsite storage, retention, encryption, integrity checks, restore drills, cross-server migration, panel-to-panel import, staging/cloning, and cutover verification.

**AGT-095:** Backup jobs MUST declare included resources, consistency guarantees, encryption, retention, and destination failure domain. **AGT-096:** Backups MUST be encrypted in transit and at rest, with independently managed restore credentials. **AGT-097:** Restore MUST default to non-destructive preview or new-target restore; overwrite requires explicit approval and a recoverable pre-restore point where feasible. **AGT-098:** Backup success MUST include artifact integrity verification, not merely successful process exit. **AGT-099:** Offsite backup and restore drills MUST be supported; a backup on the same VPS alone is not disaster recovery. **AGT-100:** Migrations MUST report preflight compatibility, transferred resources, data consistency, DNS/mail cutover risks, and post-migration validation. **AGT-101:** Large transfers MUST support resumable, bandwidth-controlled, checksummed transport.


- **Backups are batch workloads.** `cpu.weight=10`, `cpu.max=50% of cores`. A scheduled backup does not starve websites.
- **Backup verification is first-class.** Every backup run produces a `storage.backup_verifications` record. The agent performs the restore test to a disposable target and reports success/failure. A backup without a verification within policy window is `unverified`, not `healthy`.
- **`restic` is the v1 backup engine.** The agent runs `restic backup` and `restic restore`. Remote destination is S3-compatible (Backblaze B2, R2, etc.).
- **On Nano, backups are remote-only.** The disk is too small for local backups. The installer requires a remote target during onboarding.
- **`ServiceTemplate` export is a backup-like operation.** A Nextcloud instance is exported as an OCI artifact + manifest. The agent handles the export.

## 20. Firewall, SSH, WAF, malware and security response (Area 13)

**AGT-102:** Firewall changes MUST include a dry-run diff, remote-access impact assessment, timed rollback, and connectivity confirmation. **AGT-103:** SSH key, root login, and password-authentication changes MUST require an independently verified alternative management path. **AGT-104:** WAF policies MUST be scoped per site, versioned, tested, and reversible. **AGT-105:** Malware scanning MUST have clear scope, false-positive review, quarantine, and release workflows. **AGT-106:** Security incident actions MUST preserve forensic evidence and emit high-priority audit events. **AGT-107:** Tenant isolation MUST combine OS identity, path confinement, network/resource policy, and stronger workload sandboxing where needed; no single control is sufficient.


- **The `container.validate_safety` operation is mandatory** for every container creation. The agent inspects the container config and rejects:
  - `--privileged`
  - host PID namespace (`--pid=host`)
  - host network namespace (`--network=host`)
  - host IPC namespace (`--ipc=host`)
  - dangerous Linux capabilities (ALL, SYS_ADMIN, SYS_MODULE, SYS_PTRACE, SYS_RAWIO, NET_ADMIN, NET_RAW, etc.)
  - host root bind mounts (`/`, `/etc`, `/var`, `/proc`, `/sys`, `/dev`, `/boot`)
  - unrestricted device access (`--device=/dev/...` for anything other than explicitly allowed)
- **The agent has its own allowlist** of safe capabilities and mounts. Anything outside the allowlist is rejected with `CONTAINER_SAFETY_VIOLATION`.
- **This is defense in depth.** The API also validates the template's container config. The agent re-validates because a compromised control plane could bypass the API validation.

## 21. Containers and orchestration (Area 23)

Optional Docker/Compose adapter supports approved image registries, image digest pinning, Compose validation, templates, networks, volumes, secrets, environment variables, logs, lifecycle, resource limits, health checks, and controlled upgrades. Kubernetes is a separate future adapter rather than a dependency of the first VPS deployment.

**AGT-108:** Tenant users MUST NOT receive direct Docker daemon socket access; membership in the Docker group is effectively host-root-equivalent. **AGT-109:** Container definitions MUST reject privileged mode, host PID/network namespaces, unrestricted host mounts, device access, and dangerous Linux capabilities by default. **AGT-110:** Image pulls MUST support registry allowlists, digest pinning, vulnerability scanning, and provenance policy. **AGT-111:** Volumes and networks MUST be tenant-scoped with quota and backup semantics. **AGT-112:** Compose deployments MUST have plan/diff, health verification, rollout strategy, and rollback constraints. **AGT-113:** Container resource and egress limits MUST be enforced and monitored; oversubscription policy belongs to the provider.


- **Every website, app, database, cache, and self-hosted service on Standard+ is a container.** The container is created by the `service.*` adapter or the `container.*` adapter.
- **The agent creates the container with `cpu.weight`/`memory.max` at creation time.** These are cgroup values, not Docker flags, but Docker applies them via the container's cgroup.
- **Container networks are per-site or per-organization.** Inter-container communication is explicit.
- **The agent rejects `--privileged` etc.** (see AGT-108/AGT-109 above).
- **Image pulls are digest-pinned.** A tag is resolved to a digest at first use; subsequent pulls use the digest. A mutable tag is rejected for templates.

## 22. Performance, caching and CDN (Area 24)

Optional adapters: Nginx proxy cache, Varnish, PHP OPcache, Redis, Memcached, and external CDN. The UI distinguishes per-site caching from shared caches and application data stores.

**AGT-114:** Cache provisioning MUST declare ownership, memory limits, persistence, eviction, and authentication policy. **AGT-115:** Cache purge operations MUST be scoped by tenant/site and must not accidentally flush shared caches. **AGT-116:** Redis/Memcached instances MUST not be exposed publicly by default; remote access requires explicit secure networking policy. **AGT-117:** CDN integrations MUST distinguish external API changes from local origin changes and record propagation uncertainty. **AGT-118:** Performance recommendations MUST not silently alter application behavior or bypass authenticated/private-content caching rules.


- **`POST /service-instances` with `template_id = "redis"`** creates a Redis instance. The agent uses the `service.*` adapter.
- **Storage/transform/delivery are also `ServiceTemplate` instances.** They share the same adapter contract, health check pattern, credential generation, export path.
- **The agent's adapter for storage is MinIO or Garage.** Both are S3-compatible containers. The template record supplies the choice.
- **The agent's adapter for transform is imgproxy or thumbor.** Both are CPU-batch workloads.

## 23. Monitoring, analytics, resource accounting and site quality (Areas 01, 14, 26)

Agent collectors: CPU/memory/disk/process health, service state, per-tenant cgroup usage, filesystem quotas, network counters, proxy access logs, TLS status, uptime probes, and deployment/backup outcomes. Higher-level traffic analytics, bot classification, SEO auditing, synthetic performance tests, and reporting are control-plane or external-worker features.

**AGT-119:** Metrics MUST include source, timestamp, server, resource scope, and attribution confidence. **AGT-120:** Per-site CPU/memory/network usage MUST not be presented as precise if shared processes prevent accurate attribution; disclose estimated/unknown values. **AGT-121:** Logs MUST be access-controlled, sanitized, bounded, and retained according to policy. **AGT-122:** Alert evaluation MUST support thresholds, deduplication, escalation, maintenance windows, and channel delivery status. **AGT-123:** SEO and performance scans MUST use rate limits and avoid leaking unpublished or private site content. **AGT-124:** Emergency/mobile status views MUST remain useful in degraded mode and not expose privileged controls without strong authentication.


- **The `metrics.collect` operation reads cgroup v2 metrics:**
  ```
  /sys/fs/cgroup/<workload>/cpu.stat          → usage_usec, throttled_usec, nr_throttled
  /sys/fs/cgroup/<workload>/memory.current    → current bytes
  /sys/fs/cgroup/<workload>/memory.peak       → peak bytes
  /sys/fs/cgroup/<workload>/memory.events     → oom_kill count, high events
  /sys/fs/cgroup/<workload>/io.stat           → rbytes, wbytes
  /sys/fs/cgroup/<workload>/pids.current      → process count
  ```
- **Polling interval:** 15–30 seconds. Aggregated to hourly buckets.
- **Storage:** SQLite on Nano, VictoriaMetrics on Standard+.
- **The dashboard shows current / typical (7-day p50) / peak (24h p99) / throttled seconds / OOM events.**
- **The upgrade recommendation is derived from these metrics.** Peak RAM > 80% for 7+ days → "Consider upgrading." OOM kills → "A workload was killed for memory." CPU contention → "Workloads are competing for CPU."

## 24. Terminal, cron, automation, extensions and AI (Areas 15, 17, 18, 20)

### 24.1 Terminal and developer tools

**AGT-125:** Interactive terminals MUST run as a scoped workload identity within the permitted site/container, never as unrestricted host root. **AGT-126:** Session recording, duration, resource limits, file access, and command audit policy MUST be defined; sensitive input must be redacted appropriately. **AGT-127:** Cron jobs MUST have tenant ownership, quotas, safe execution context, logs, and conflict detection. **AGT-128:** WP-CLI, Composer, npm, and equivalent tooling MUST execute in the tenant's sandbox with explicit runtime selection.


- **Terminal on Standard+ runs inside the container** (`docker exec -it <container> /bin/bash`). The user is the container's non-root user. No host root.
- **Terminal on Nano is unavailable.** The agent returns `CAPABILITY_UNAVAILABLE` with reason "No container to exec into."
- **Cron jobs are classified by workload category.** A `restic backup` cron is batch; a curl cron is serving. The category is derived from the command and stored in `web.scheduled_tasks.workload_category`.

### 24.2 Automation and extension model

**AGT-129:** Public API, CLI, webhooks, IaC, and scheduled automation MUST use the same policy and operation engine as UI requests. **AGT-130:** Extensions MUST declare permissions, supported agent/API versions, resource scope, network access, and upgrade compatibility. **AGT-131:** Third-party extension code MUST NOT run inside the privileged executor or gain unrestricted host access. **AGT-132:** Extension failures MUST be isolated and visible without compromising core agent health.


### 24.3 AI actions

**AGT-133:** AI-generated plans MUST be validated as if submitted by any other untrusted client. **AGT-134:** AI MUST NOT call raw shell/root operations or read raw credentials. **AGT-135:** Every AI-initiated action MUST record the initiating user, model/tool context reference, policy decision, approval, operation ID, and final evidence. **AGT-136:** High-risk AI actions MUST require explicit human approval and cannot self-approve.


- **The agent is a batch workload when running inference.** But actually — the agent doesn't run inference. The AI is in the control plane. The agent sees only `POST /ai/action-proposals` results as approved operations. The agent has no `ai.*` operations.
- **The AI's approved action is a normal operation.** It flows through the same envelope, same policy, same agent. The `requested_by.kind = "ai"` field is the only difference, and it's recorded in audit.

## 25. Complete 29-area traceability matrix


| # | Platform area | Agent responsibility | Control-plane / external responsibility | Principal adapter / contract |
|---|---|---|---|---|
| 01 | Dashboard | Health and resource telemetry | Role-specific dashboard, uptime/activity aggregation | `inventory`, `metrics.collect` |
| 02 | Server management | Packages, services, SSH/firewall/network safely | Fleet scheduling, IPAM, provider PTR requests | `host`, `systemd`, `network` |
| 03 | Websites | Roots, vhosts, redirects, staging hooks, **containers** | Ownership, templates, website lifecycle | `nginx.website`, `container.website` |
| 04 | DNS | **DNS guidance check (v1); authoritative DNS (v2+)** | Registrar/external DNS, delegation, propagation | `dns.guidance.check` |
| 05 | SSL/TLS | ACME HTTP-01 challenges, keys, deploy/reload | Certificate policy, expiry tracking | `tls.acme` |
| 06 | File manager | Safe file ops, SFTP isolation, archives, **container volume scoping** | UI editor, file permissions policy | `files.scoped`, `sftp` |
| 07 | Application manager | Process and app lifecycle, templates, **build container / runtime container** | App configuration, deployments | `app.node`, `app.python`, `app.php` |
| 08 | Runtime manager | Runtime install/version/extension isolation in containers | Compatibility catalog, lifecycle policies | `runtime.*` |
| 09 | Databases | **`ServiceTemplate` instances (v2+)** | Metadata, permissions, admin UI brokering | `service.*` (template=`database.*`) |
| 10 | Email hosting | **None (v1 external SMTP); self-hosted (v2+)** | Mail account UX, domain policy | `mail.external` (v1), `mail.smtp`, `mail.imap` (v2+) |
| 11 | Git and deployments | Sandboxed builds (batch), artifacts, health/rollback | Git OAuth/webhooks, approvals, history | `build` (batch), `deploy` (serving) |
| 12 | Backups | Consistent capture, encrypt, transfer, restore, **verification** | Retention, offsite orchestration, recovery UI | `backup` (batch), `restore`, `backup.verify` |
| 13 | Security | Firewall/SSH/WAF/quarantine actions, **container safety validation**, **cgroup validation** | Policy, approvals, incident handling | `security.*`, `container.validate_safety`, `cgroup.validate` |
| 14 | Monitoring | Host/site telemetry, logs, **cgroup metrics** | Alerts, trends, notifications | `metrics`, `logs`, `cgroup.metrics.read` |
| 15 | Terminal/developer tools | Scoped shell (container), cron, package tools | Access grants, session UX, audit | `terminal.scoped`, `cron` |
| 16 | User/account management | Provision OS identities, enforce quotas, **cgroup limits** | Users, SSO/2FA, RBAC, suspension | `identity.local`, `quota`, `cgroup.*` |
| 17 | AI assistant | **Agent has no AI operations** — approved actions flow through standard operations | AI reasoning, tool routing, consent/audit | Normal operation contract |
| 18 | Automation and API | Execute typed jobs and report events | API, CLI, webhooks, IaC orchestration | Operation/event contract |
| 19 | Panel administration | Agent update, local health, recovery, **capability rescan** | Panel upgrades, license, maintenance | `agent.lifecycle`, `capability.rescan` |
| 20 | Extensions | **`ServiceTemplate` is the v1 extension mechanism; plugin runtime is v2+** | Marketplace, signing, permissions | `service.*` |
| 21 | Commerce/reseller | Enforce assigned limits and suspensions | Billing, invoicing, reseller hierarchy/branding | `quota`, `account.lifecycle` |
| 22 | Migration/staging | Transfer/import, clone, verify, **`ServiceTemplate` export/import** | Mapping, cutover, reports, approvals | `migration`, `staging`, `service.export`, `service.import` |
| 23 | Container/orchestration | **Docker is primary primitive** — create, validate, cgroup, health | Templates, policy, cluster scheduling (v2+) | `container.*`, `cgroup.*` |
| 24 | Performance/caching | Nginx/Caddy cache, MinIO, imgproxy, Redis, Memcached | CDN API, caching policies, UX | `cache.*`, `cdn.external`, `service.*` |
| 25 | Email security/deliverability | **None in v1; v2+ self-hosted** | SPF/DMARC guidance, reputation reporting | v2+ |
| 26 | Analytics/site quality | Access/usage signals, safe probes, **cgroup metrics** | Traffic reports, bot/SEO/performance analysis | `analytics.collect`, `site.scan`, `cgroup.metrics` |
| 27 | Audit/compliance/support | Signed operation receipts, local change evidence, **hash chain verification** | Audit retention/exports, helpdesk, emergency UI | `audit.receipt`, `support.external`, `audit.chain.verify` |
| **28** | **Platform Services** *(new)* | **`ServiceTemplate` operations for storage/transform/delivery/cache/queue/search** | Service catalog, metering, billing units | `service.provision`, `service.*` |
| **29** | **Self-hosted Templates** *(new)* | **`ServiceTemplate` operations for Nextcloud/Gitea/Ghost/n8n/etc.** | Template catalog (data), signature verification | `service.*` (template-specific) |

**Traceability rule:** Every feature in each area MUST ultimately have linked requirement ID, UX screen/flow, API operation, agent adapter (or explicit `not_agent_responsibility`), data entity, permission rule, threat control, and acceptance test. A row in this table is a mapping, not proof of implementation completeness.

## 26. Tenant and resource isolation strategy

| Layer | Baseline control | Stronger isolation for untrusted hosting |
|---|---|---|
| Identity | Separate OS user/group per hosting account | User namespaces and dedicated workload sandbox |
| Files | Owner/mode/ACL, safe path resolution, quotas | Mount namespaces, isolated volumes, mandatory access controls |
| Processes | systemd units, cgroup v2 resource accounting | Container or stronger sandbox with seccomp/AppArmor |
| Network | Per-app sockets/ports and firewall policy | Network namespaces, egress policy, tenant segmentation |
| Databases | Dedicated DB users and privileges | Dedicated DB instances/hosts where policy requires |
| Secrets | Scoped encrypted references and injection | Short-lived credentials, dedicated secret delivery |
| Logs | Per-resource permissions and redaction | Separate retention/storage domains |
| Backups | Tenant-scoped encrypted artifacts | Separate keys/buckets and restore authorizations |

**AGT-137:** Isolation profile MUST be explicit per workload/plan. **AGT-138:** A reseller's delegated rights MUST never imply OS-level access to other resellers or provider infrastructure. **AGT-139:** Resource quotas MUST be enforced where possible and otherwise labeled advisory. **AGT-140:** Suspension MUST specify what is blocked (web, login, mail, deploy, background jobs) and preserve data per retention policy.


| Layer | Nano (no Docker) | Starter (borderline) | Standard+ (Docker) |
|---|---|---|---|
| Identity | Per-site OS user | Per-site OS user | Container user (non-root) |
| Files | Per-user directory + ACL | Per-user directory + ACL | Container volume with bind mount |
| Processes | systemd unit or direct | systemd unit | Docker container |
| Network | Shared Nginx | Shared Nginx | Per-site bridge network |
| Database | N/A | Shared MariaDB container (borderline) | Per-database container |
| Secrets | Panel secrets store | Panel secrets store | Docker secrets + panel secrets store |
| cgroup | Per-workload cgroup v2 | Per-workload cgroup v2 | Per-container cgroup v2 |

## 27. Secrets, credentials, enrollment and rotation

**AGT-141:** Secrets MUST be encrypted at rest, protected in transit, and delivered to the minimum required workload. **AGT-142:** Agent identity and service credentials MUST support rotation, revocation, expiration, and audit. **AGT-143:** Secrets MUST never be included in plan diffs, operation logs, crash reports, analytics, or AI transcripts. **AGT-144:** Enrollment MUST require a short-lived, single-use bootstrap authorization and server identity confirmation. **AGT-145:** Credential exposure incidents MUST have a defined revocation and re-enrollment workflow.

### 27.1 Example enrollment flow

1. Operator registers a server in the control plane and receives an expiring one-time enrollment token.
2. Installer verifies signed agent artifact, OS compatibility, and transport trust.
3. Agent generates its own private key locally and submits a CSR/attestation under the one-time token.
4. Control plane verifies enrollment context, records approval, and issues a bounded identity certificate.
5. Agent connects and publishes capability inventory; control plane marks it `pending_validation` until a health check succeeds.
6. Token is consumed and cannot be replayed. Rotation and revocation are tested.


- **v1 has no enrollment flow.** The agent is installed by the panel installer on the same machine. The agent runs as `nexuspanel-agent` UID. The panel API runs as `nexuspanel-api` UID. The Unix socket at `/run/nexuspanel/agent.sock` is mode 0660 owned by `nexuspanel:nexuspanel-agent`. Both UIDs are in the `nexuspanel` group.
- **Peer credential check on every connection.** `getsockopt(SO_PEERCRED)` returns the connecting process's UID. The agent verifies it is the panel API's UID. No other process can connect.
- **The agent's secret store access is via the Unix socket.** The agent does not read the secrets DB directly. The panel API injects the resolved secret into the operation envelope at dispatch time.
- **v2+ uses mTLS.** The enrollment flow above applies for remote agents. Rotation and revocation are tested.

## 28. Upgrade, rollback and disaster recovery

**AGT-146:** Agent updates MUST use signed, verified artifacts and version compatibility gates. **AGT-147:** Upgrade MUST preserve last-known-good binary/config and provide health-gated rollback where safe. **AGT-148:** Control-plane and agent protocol MUST have a documented version compatibility window and graceful feature negotiation. **AGT-149:** An unavailable control plane MUST not interrupt already-running hosted sites solely because management connectivity is lost. **AGT-150:** Agent recovery MUST distinguish control-plane metadata loss, local journal loss, service corruption, credential compromise, and full server loss. **AGT-151:** Disaster-recovery drills MUST validate restoration of metadata, managed configs, customer data, credentials, and authorization boundaries.


- **The agent is a single static Go binary.** Upgrade is: stop agent, replace binary, start agent, verify. Rollback is: stop agent, replace with previous binary, start. The journal survives.
- **Panel upgrade is: stop API, run migrations, start API, verify.** Agent upgrade is: stop agent, replace binary, start agent. On Nano, this happens in seconds.
- **`ServiceTemplate` export/import is the migration path.** A Nextcloud instance is exported as an OCI artifact + manifest. On the new VPS, `POST /service-instances/import` recreates it.
- **The audit chain is backed up separately.** Even if the primary DB is lost, the chain can be reconstructed.

## 29. Observability, audit and support evidence

Every operation emits structured events: `operation.accepted`, `plan.generated`, `approval.granted`, `step.started`, `step.completed`, `step.failed`, `verification.passed`, `verification.failed`, `rollback.started`, `rollback.completed`, and `operation.closed`.

```json
{
  "event_type": "operation.step.completed",
  "operation_id": "op_01JEXAMPLE",
  "correlation_id": "req_01JEXAMPLE",
  "tenant_id": "org_42",
  "server_id": "srv_1",
  "adapter": "nginx.website",
  "step": "config_validation",
  "outcome": "passed",
  "duration_ms": 148,
  "config_digest": "sha256:<digest>",
  "timestamp": "2026-11-06T15:02:31Z",
  "workload_category": "serving",
  "cgroup_applied": {"cpu_weight": 100, "memory_max_bytes": 419430400}
}
```

**AGT-152:** Audit records MUST be tamper-evident and shipped to a storage/retention boundary not writable by the workload tenant; do not claim mathematical immutability without WORM or equivalent guarantees. **AGT-153:** Support exports MUST redact secrets and personal data according to policy. **AGT-154:** Each operation MUST have a human-readable incident/recovery summary. **AGT-155:** Log ingestion failure MUST be visible; security-sensitive operations may fail closed when mandatory audit recording is unavailable.


- **The audit chain is hash-chained.** Each `audit_events` row has `previous_hash` and `record_hash`. The agent writes to the chain on every operation. The control plane verifies hourly. A broken chain is a critical incident.
- **cgroup metrics are audited.** Every `cgroup.update` operation records the before/after values in the audit chain.
- **Batch yield events are audited.** Every `batch.yielded` and `batch.resumed` is recorded. This is how you prove the batch job yielded to site traffic, not failed.

## 30. Threat model and defensive controls

| Threat | Primary controls | Verification |
|---|---|---|
| Forged/replayed agent request | mTLS, signed envelope, nonce, expiry, idempotency | Replay and signature-negative tests |
| Tenant path traversal | FD-relative open, symlink-safe resolver, owned roots | Traversal/symlink race tests |
| Privilege escalation via Docker | No socket exposure, reject privileged/host mounts | Malicious Compose tests |
| Shell injection via build command | Tenant sandbox; no privileged shell interpolation | Metacharacter and escape tests |
| Unauthorized cross-tenant action | Control-plane policy + agent target binding | Cross-tenant fuzz/property tests |
| Firewall/SSH lockout | Timed revert, console recovery, confirmation | Remote access rollback drill |
| Poisoned extension or package | Signing, allowlist, isolation, capability review | Supply-chain verification |
| Secret exfiltration via logs | Redaction, scoped injection, log access control | Canary secret leak tests |
| Host configuration drift | Ownership markers, checksums, adoption policy | Manual edit/drift tests |
| Partial operation and retry | Journal, idempotency, locks, checkpoints | Kill/restart fault injection |
| Compromised agent identity | Revocation, rotation, bounded scope | Credential revocation exercise |
| Malicious customer workload | Sandboxing, quotas, network controls | Isolation penetration tests |


| Threat | Primary controls | Verification |
|---|---|---|
| **cgroup escape via workload** | Non-root workload, no host PID, seccomp/AppArmor, read-only root filesystem, cgroup v2 | Container penetration tests |
| **cgroup misconfiguration** | Agent validates all cgroup values; rejects `cpuset` and unknown controllers; verifies applied values match | cgroup behavior tests |
| **Container escape** | Agent rejects `--privileged`, host PID/network, dangerous capabilities, host root mounts | Malicious Compose tests, fuzzing |
| **Capability spoofing** | Capability records signed by agent (HMAC); control plane verifies digest | Installer tamper tests |
| **Persona matrix tampering** | Persona matrix signed; control plane verifies | Data file integrity tests |
| **Idempotency bypass via payment retry** | Idempotency key scoped to actor + invoice + amount + provider | Duplicate payment tests |
| **Payment credential exposure** | Agent has no payment credentials; payment is control-plane only | Agent inspection test |
| **ServiceTemplate image poisoning** | Digest-pinned image; signature verification before pull | Malicious image tests |

## 31. Detailed acceptance test catalog

| Test ID | Requirement | Pass condition |
|---|---|---|
| AT-001 | Agent enrollment | Single-use token cannot enroll twice; revoked identity rejected |
| AT-002 | Typed RPC | Unknown action and extra security-sensitive fields rejected |
| AT-003 | Authorization | Cross-tenant resource IDs rejected even with syntactically valid envelope |
| AT-004 | Replay | Duplicate exact operation returns original status, not duplicate side effects |
| AT-005 | Revision conflict | Stale plan rejected and replanning required |
| AT-006 | Service configuration | Invalid Nginx config does not replace active config |
| AT-007 | Host restart | Interrupted operation reconciles against journal and observed state |
| AT-008 | Agent disconnect | Site stays available; dashboard shows agent offline |
| AT-009 | SSH lockout prevention | Failed confirmation triggers timed firewall/SSH rollback |
| AT-010 | File traversal | `../`, symlink and rename races cannot escape tenant root |
| AT-011 | Archive safety | Zip-slip and decompression bomb rejected without tenant corruption |
| AT-012 | Runtime sandbox | Build command cannot write outside assigned workload sandbox |
| AT-013 | Secret safety | Secrets absent from logs, plans, traces, and AI outputs |
| AT-014 | Database scope | Customer DB credentials cannot access another customer's DB |
| AT-015 | DNS | Invalid DNSSEC/record operation blocked with recovery guidance |
| AT-016 | TLS | Renewal failure alerts before expiry; private keys remain inaccessible |
| AT-017 | Mail | Mailbox quota and outbound throttling enforced; queue state visible |
| AT-018 | Backup | Restore to new target verifies integrity and permissions |
| AT-019 | Migration | Resume after interruption without duplicate or silent data loss |
| AT-020 | Containers | Privileged Compose/host socket mount rejected |
| AT-021 | Cache | Tenant purge cannot flush another tenant's cache |
| AT-022 | Metrics | Attribution displays estimated/unknown when precision unsupported |
| AT-023 | AI | AI cannot self-approve a high-risk action or call root shell |
| AT-024 | Reseller | Child quota cannot exceed assigned delegation; isolation holds |
| AT-025 | Panel update | Bad signed-artifact verification aborts; previous agent remains healthy |
| AT-026 | Audit | Tamper attempt detected; support export excludes secrets |
| AT-027 | Unmanaged config | Existing admin-owned Nginx config is never silently overwritten |
| AT-028 | Performance | Concurrent job limits preserve hosting service availability |
| AT-029 | Recovery | Corrupted local journal triggers safe read-only/recovery mode |
| AT-030 | Compatibility | Unsupported Ubuntu/package combo yields explicit preflight failure |


| Test ID | Requirement | Pass condition |
|---|---|---|
| **AT-031** | **Capability detection** | 5-step detection on 512 MB, 1 GB, 2 GB, 4 GB, 8 GB reference VPSes produces correct capability sets |
| **AT-032** | **Capability gate** | `CAPABILITY_UNAVAILABLE` returned with reason and resource delta when workload exceeds capability |
| **AT-033** | **Capability rescan** | After hardware bump, rescan enables new capabilities; no reinstall required |
| **AT-034** | **cgroup applied** | `cgroup_applied` in operation response matches `/sys/fs/cgroup/` values |
| **AT-035** | **cgroup validation** | `cpuset` rejected; `memory.max` > physical × 1.2 rejected |
| **AT-036** | **Serving no cpu.max** | A serving workload can use all idle CPU; no `cpu.max` applied |
| **AT-037** | **Batch cpu.max** | A batch workload has `cpu.max = 50% of cores`; yields under contention |
| **AT-038** | **Batch queue** | `GET /batch-queue` shows queued jobs with positions |
| **AT-039** | **Batch yield** | A batch job pauses under site traffic; resumes when idle; not marked failed |
| **AT-040** | **Container safety** | `--privileged`, host PID/network, dangerous capabilities all rejected |
| **AT-041** | **ServiceTemplate provision** | Nextcloud instance provisions, health-checks, reachable with SSL |
| **AT-042** | **ServiceTemplate export** | Instance exports to OCI artifact + manifest; digest matches |
| **AT-043** | **ServiceTemplate import** | Artifact imports on new VPS; instance reproduces |
| **AT-044** | **Unix socket peer check** | Non-panel process cannot connect to agent socket |
| **AT-045** | **Audit chain** | `POST /audit-logs/verify` returns `chain_intact: true`; tampered chain returns false |
| **AT-046** | **Nano install** | Panel installs on 512 MB VPS; static sites work; unsupported capabilities show honest reason |
| **AT-047** | **M-Pesa callback security** | Tampered callback rejected; valid callback resolves; duplicate returns 200 |
| **AT-048** | **Payment state machine** | `pending → confirmed`, `pending → failed`, `pending → expired → reconciled` all work |
| **AT-049** | **Payment credential isolation** | Agent has no payment credentials; payment is control-plane only |
| **AT-050** | **Guided router landing** | First-time user on a Nano VPS sees the router with honest capability display |

### 31.1 Test environment matrix

- Single-server Ubuntu 26.04.1 LTS with 4 vCPU / ~6 GiB RAM / 100 GB disk.
- **Nano test:** Ubuntu 24.04 with 512 MB / 1 vCPU / 10 GB.
- **Starter test:** Ubuntu 24.04 with 2 GB / 1 vCPU / 25 GB.
- At least two Linux test hosts for multi-agent transport, disconnection, and failover tests.
- Dedicated secondary DNS node for redundancy validation.
- Isolated container/build test host for adversarial tenant workloads.
- External mail/DNS/provider test environments for integrations that cannot be safely emulated locally.
- Fault injection: process kill, network partition, disk full, DNS propagation delay, expired certificates, unavailable backup store, conflicting operator edits.
- **cgroup fault injection:** synthetic load to trigger yielding, OOM kills, memory.high throttling.
- **Payment fault injection:** Daraja sandbox failure modes, callback replay, timeout.

## 32. Performance and capacity budgets (proposed; validate by measurement)

| Metric | Proposed initial engineering target | Caveat |
|---|---|---|
| Agent idle RSS | <= 150 MiB | Excludes service processes and sidecars |
| Heartbeat | Every 30 s, configurable | Jitter and batching required at scale |
| Offline detection | 2–3 missed heartbeats | Not a service availability guarantee |
| Inventory refresh | 5–15 min + on-demand | Avoid expensive full scans |
| Simple operation dispatch | p95 < 2 s when connected | Excludes operation execution time |
| Local journal durability | Survives agent crash/restart | Requires fsync/storage testing |
| Concurrent host mutations | Resource-scoped bounded worker pool | Avoid host saturation and deadlocks |
| Per-site metrics | Configurable 30–60 s collection | Cardinality and cost must be bounded |

These are **design targets, not measured performance claims**. Test results must replace them before release.


| Metric | Nano (512 MB) | Starter (2 GB) | Standard (4–8 GB) | Provider |
|---|---|---|---|---|
| Panel API RSS | 60 MB | 80 MB | 120 MB | 150 MB |
| Agent RSS | 40 MB | 50 MB | 80 MB | 100 MB |
| Static SPA (Nginx) | 10 MB | 10 MB | 15 MB | 20 MB |
| SQLite (Nano) | 15 MB | — | — | — |
| PostgreSQL | — | 200 MB | 400 MB | 800 MB |
| Caddy | — | 30 MB | 50 MB | 80 MB |
| Docker daemon | — | 80 MB | 150 MB | 200 MB |
| VictoriaMetrics | — | — | 150 MB | 300 MB |
| **Total panel footprint** | **~125 MB** | **~450 MB** | **~965 MB** | **~1.5 GB** |
| **Headroom for workloads** | **~387 MB** | **~1.5 GB** | **~3–7 GB** | Scales |

These are design targets. The installer validates against actual measurements and refuses to install below floor.

## 33. Implementation sequence and deliverables

**Phase A — Agent skeleton and security:** Go service, typed RPC, identity enrollment, journaling, audit, policy binding, capability inventory, signed updates, sandboxed helper, test harness. Deliverable: a no-op read-only agent that safely reports inventory and rejects unauthorized mutations.

**Phase B — First real operations:** Managed Nginx site lifecycle, scoped file operations, static app deploy, ACME certificate flow, process/log health, idempotent rollback. Deliverable: verified static website provisioning without raw root access.

**Phase C — Runtime and data services:** Node/Python/PHP, MariaDB/PostgreSQL, Git builds, backups/restores, quotas, tenant isolation. Deliverable: repeatable application deployment and restore tests.

**Phase D — Hosting-provider services:** Local/external DNS, email, SFTP, IPAM/PTR coordination, reseller enforcement, multi-server scheduling, migrations. Deliverable: provider/reseller lifecycle with cross-server tests.

**Phase E — Advanced modules:** Containers/Compose, caching/CDN, WAF, malware scanning, analytics/site quality, extension SDK, AI-controlled operations. Deliverable: full 29-area traceability with acceptance evidence.

Implementation phases describe dependency order, **not exclusions from the product**. Features must not be advertised as complete until tests and operational readiness are met.


**v1 (Foundation + first-complete-journey):**
- Go agent with Unix socket (`SO_PEERCRED`)
- 5-step capability detection
- cgroup v2 adapter (serving/batch/infrastructure)
- Nginx static site serving (Nano) + Caddy reverse proxy + Docker containers (Standard+)
- `ServiceTemplate` adapter (storage, transform, delivery, Nextcloud)
- `container.*` adapter with safety validation
- Backup via restic (remote-only on Nano)
- DNS guidance check (no authoritative DNS)
- External SMTP
- Audit chain
- Batch queue and yield
- All 50 acceptance tests pass

**v2:**
- mTLS for remote agents
- Authoritative DNS
- Self-hosted email
- Plugin runtime
- Reseller/provider surface
- Stripe adapter
- Migration from cPanel/DirectAdmin
- 20–30 template catalog

**v3:**
- HA topology (Postgres replication, leader election)
- Multi-server scheduling
- Cross-region backups
- Marketplace

## 34. Release gates

1. All security-critical `AGT-` requirements are implemented or formally waived with a documented risk owner.
2. Threat model reviewed; penetration and cross-tenant isolation tests passed.
3. Operation schemas versioned; API and database references reconciled with Stages 5–6.
4. Safe update and rollback demonstrated on supported Ubuntu releases.
5. Restore drills completed for panel metadata, sites, databases, and secrets.
6. No unrestricted root shell, Docker socket, or host filesystem access available to customers, AI, or extensions.
7. Monitoring, alerting, audit, and incident runbooks validated.
8. Every shipped adapter has installation, compatibility, ownership, health, upgrade, rollback, and uninstall contracts.
9. Documentation traceability links user-visible features to tests and evidence.
10. A documented capacity profile exists for small VPS, medium hosting node, and provider cluster.


11. **Capability gate:** installer correctly detects Nano/Starter/Standard/Provider on reference VPSes; refuses below floor; produces honest capability displays with resource delta.
12. **cgroup gate:** serving workloads never throttled on idle CPU; batch workloads yield under contention; `memory.low` protects serving from batch memory pressure; `cpuset` rejected; applied values read back match.
13. **Container safety gate:** all AT-040 checks pass; malicious Compose rejected; container escape tests pass.
14. **ServiceTemplate gate:** storage/transform/delivery/Nextcloud provision, health-check, backup, export, import, uninstall cleanly.
15. **Async payment gate:** M-Pesa sandbox STK push → callback → confirmation; timeout → retry; duplicate callback idempotent; unmatched callback → reconciliation queue; Stripe webhook equivalent.
16. **Audit chain gate:** hash chain intact after 10,000 operations; tamper detected; off-host export verified.
17. **Nano gate:** panel installs on 512 MB VPS; static site works in 60 seconds; unsupported capabilities shown with reason.
18. **Guided router gate:** first-time user sees router; each of six paths ends at "you're live" moment; capability display is honest; **release blocker** — the lay-person experience is the reason the product exists.



## 35. Consolidated 29-area host responsibility map

| Area | Host responsibility | Profile / release boundary |
|---|---|---|
| 01 Dashboard | Supply bounded, timestamped health and resource measurements | Minimal Nano sampling |
| 02 Server management | Package, systemd, SSH, network and service adapters | Dangerous mutations approval-gated |
| 03 Websites | Site directories, web server routes, ownership and health | Nano static; advanced runtimes larger profiles |
| 04 DNS | Guided validation; optional authoritative adapters | v1 guidance, authoritative later |
| 05 SSL/TLS | ACME challenges, secure key custody, renewal and reload | HTTP-01 first; DNS-01 capability-dependent |
| 06 File manager | Confined file operations, SFTP, archive and edit | Path and ownership isolation required |
| 07 Application manager | App lifecycle, builds, health, logs | Docker-capable profiles |
| 08 Runtime manager | Versioned runtimes, isolation and extensions | Docker-capable profiles |
| 09 Databases | Provision, users, backups, integrity and upgrade | Specialized stateful adapter; later release |
| 10 Email hosting | Mailbox/transport operations and migrations | External SMTP v1; local hosting later |
| 11 Git/deployments | Fetch, build, stage, deploy and rollback | Resource-bounded batch builds |
| 12 Backups | Snapshots, encrypted transport, restore verification | Remote backup support on Nano where feasible |
| 13 Security | Firewall/SSH/WAF controls, quarantine and malware hooks | Emergency recovery mandatory |
| 14 Monitoring | cgroup and service metrics, logs and event capture | Low-overhead Nano mode |
| 15 Terminal/developer | Scoped, time-limited execution and cron | Not unrestricted privileged shell |
| 16 User/account management | Apply ownership, OS identities and quotas | Authorization remains control-plane owned |
| 17 AI assistant | Expose safe read-only diagnostics and governed actions | No privileged AI bypass |
| 18 Automation/API | Execute durable typed jobs, emit events | Same policy as interactive requests |
| 19 Panel administration | Installer, update, rollback and recovery | Signed and compatibility checked |
| 20 Extensions | Sandboxed extension operations and permissions | Not implicit host-root access |
| 21 Commerce/reseller | Enforce scoped entitlement decisions and quotas | No payment secrets on host |
| 22 Migration/staging | Verified import, clone, staged cutover and rollback | Resource-aware batch scheduling |
| 23 Containers/orchestration | Container lifecycle, image policy, volumes and limits | Standard+; not Nano baseline |
| 24 Performance/caching | Reverse proxy, caching, TLS and CDN adapters | Dependency and resource gated |
| 25 Email security | SPF/DKIM/DMARC guidance and mail filtering adapters | Local MTA only when installed |
| 26 Analytics/site quality | Access logs and bounded local collection | Minimize PII and storage on Nano |
| 27 Audit/compliance/support | Emit signed/correlated evidence and support diagnostics | Never disable audit for performance silently |
| 28 Platform Services | Storage, transform, delivery, cache, queue, search adapters | Template plus typed specialized lifecycle |
| 29 Self-hosted Templates | Signed catalog recipe, install, update, health and uninstall | Resource-aware, image-pinned deployment |

## 36. Cross-volume verification and release evidence

Each shipping adapter MUST produce an evidence record linking `NP-*` product requirements, UI flows, API operation contracts, data entities, agent operation schemas, threat controls and automated test IDs. The following scenarios are mandatory release gates for the relevant features:

1. Local agent rejects an untrusted peer even when the socket path is discovered.
2. A stale operation, replayed nonce or mismatched payload digest cannot mutate the host.
3. A tenant cannot target another tenant's file, process, container, volume, database or network namespace.
4. Symlink swaps and path traversal are rejected during file creation, editing, archive extraction and restore.
5. Failed Nginx configuration validation leaves the previous live configuration intact.
6. Interrupted certificate renewal does not destroy the last valid certificate or private key.
7. A post-dispatch timeout produces `unknown` and reconciles actual state before retry.
8. Idempotent retries do not duplicate sites, mailboxes, database users or paid entitlements.
9. A batch build yields under CPU/memory contention without starving SSH or the control plane.
10. A disk-full condition fails safely with bounded logs and recovery instructions.
11. A container recipe requesting privileged mode, host PID, unsafe bind mount or unapproved capability is rejected.
12. An untrusted ServiceTemplate image digest or incompatible schema is rejected.
13. A database restore is verified for integrity and tenant ownership before promotion.
14. An external SMTP configuration cannot be mistaken for local mail-hosting availability.
15. Nano capability detection gives correct, understandable reasons for unavailable services.
16. A disconnected agent cannot falsely report a completed host mutation.
17. Rollback and compensation tests distinguish reversible from irreversible changes.
18. Audit records connect request, approval, execution, observed state and actor without exposing secrets.
19. Emergency SSH recovery remains possible after a failed firewall or access-control update.
20. All 29 areas are traceable to a host adapter or explicit non-host boundary.

### 36.1 Performance and capacity verification matrix

| Test profile | What must be measured | Release claim policy |
|---|---|---|
| 512 MB / fractional or 1 vCPU / 10 GB | Installer peak, idle RSS/PSS, static site throughput, update/restore disk headroom, SSH survival | Supported only after full repeatable testing |
| 1 GB / 1 vCPU / 10 GB | Same plus concurrency, memory-pressure behavior and service feasibility | No untested full-stack claim |
| 2 GB / 1–2 vCPU | Starter services, contention, queue fairness and rollback | Publish tested combinations only |
| 4–8 GB / 2–4 vCPU | Docker isolation, app/database lifecycle, builds and backups | Publish per-module resource budgets |
| Multi-node provider | Enrollment, mTLS, partitions, delegated tenancy, placement, recovery | Requires distributed conformance suite |

### 36.2 Open implementation approvals

- Confirm agent helper privilege model, local socket identity checks and enrollment authority.
- Confirm supported Ubuntu/kernel/cgroup versions with tested package availability.
- Approve signed ServiceTemplate schema, image trust policy and lifecycle adapter interfaces.
- Define measured minimum footprint and whether 512 MB is an official support tier.
- Finalize per-adapter backup, restore, update and rollback semantics.
- Approve release-specific DNS, mail, database and container capabilities.
- Specify durable operation envelope in machine-readable protobuf/JSON schema and compatibility matrix.

**End of Volume 06.** The remaining security threat model and incident playbooks are consolidated in Volume 08; product commercial rules in Volume 09; final benchmark evidence and rollout gates in Volume 10.
