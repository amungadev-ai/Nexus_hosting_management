# NexusPanel — Master Product & Engineering Specification
## Volume 10 — Implementation Roadmap, Testing, Release Engineering & Operations

**Status:** Proposed implementation baseline; formal approval, threat review, hardware benchmarks and production validation outstanding.  
**Scope:** Delivery plan across all 29 product areas and all four architectural spines (Commercial, Service, Execution, Surface).  
**Authority:** Volume 01 is authoritative for product requirements. Volumes 02–09 define the detailed domain contracts. This volume defines sequencing, evidence, operational ownership and release gates; it does not silently alter their requirements.  
**Deployment constraint:** No commands in this document authorize modification of the existing Leaseweb development VPS.  
**Normative vocabulary:** MUST / MUST NOT = release requirement; SHOULD = preferred with documented exception; MAY = optional. All numerical targets are proposed until measured and approved.

## 0. Executive delivery contract

NexusPanel is a guided hosting-control-panel product for individual site owners, developers, agencies, resellers and commercial hosting providers. A single-user Nano installation must not be architecturally forced to run a provider-grade cluster. The initial release must deliver an understandable, secure, recoverable hosting journey before advanced features expand the attack surface. Scope coverage in documentation does not imply all features ship in the same release.

**Non-negotiable engineering invariants**

1. A static frontend speaks to a Go control plane; only a narrowly privileged Go agent performs allowlisted host mutations.
2. The frontend is never the authority for authorization, licensing, capacity or entitlements. Persona and UI presentation modes do not confer permissions.
3. Every meaningful mutation has an idempotency policy, durable operation record, observable state and explicit failure semantics.
4. Nano uses the lowest feasible baseline with SQLite and optional OS-user isolation; Docker, heavyweight services and PostgreSQL must not be silently installed on Nano.
5. Standard+ may use PostgreSQL and container isolation when capability and security preflight pass. Provider fleet control is phased; local Unix-socket agent is v1, remote mTLS is later.
6. The ServiceTemplate catalog is a common interface, not a promise that all stateful services share identical lifecycle semantics.
7. Money is recorded in an independent ledger; a payment callback never independently authorizes provisioning without verified, idempotent settlement logic.
8. AI v1 is diagnostic and proposal-only. It cannot execute privileged actions autonomously.
9. Installation, upgrade, backup and restore are product features subject to destructive-path tests.
10. A declared capability must have an owner, tests, documentation, support status and evidence. A passing happy-path demo is insufficient.

## 1. Source authority, change control and requirement traceability

### 1.1 Source hierarchy

| Order | Source | Use |
|---|---|---|
| 1 | Volume 01 — Master PRD | Product behavior, 29-area scope, entitlement and release intent |
| 2 | Volumes 02–09 | UI, architecture, data, API, agent, platform services, security and commerce contracts |
| 3 | Approved architecture decision records (ADRs) | Explicitly ratified deviations and platform choices |
| 4 | Versioned OpenAPI, RPC schemas, migrations, template manifests and policy fixtures | Executable contracts that MUST trace to approved requirements |
| 5 | Tests, operational evidence and telemetry | Verification, not independent authority to redefine product scope |

Where sources disagree, record a discrepancy and obtain an owner decision. Do not let a developer resolve a pricing, permission, security or capacity contradiction by implication.

### 1.2 Traceability records

Each requirement record SHOULD include: stable `NP-*` identifier; source volume/section; area (1–29); release phase; implementation epic; API or RPC contract; data entities; frontend surface; security controls; automated test IDs; owner; verification evidence; decision status; documentation links. Additional local IDs use `DEL-*`, `QA-*`, `OPS-*`, `REL-*`, `MIG-*` and `CAP-*` without replacing existing NP identifiers.

A traceability build fails if a required NP identifier is orphaned, an API operation has no policy test, or a privileged agent action has no allowlist and abuse-case test.

### 1.3 Approval and versioning

Maintain `decisions/ADR-NNN-*.md` with context, options, chosen decision, consequences, security review and approvers. The previously proposed Go stack, AGPL/open-core licensing, Nano resource floor, Docker policy, M-Pesa initial gateway, guided router, Unix-socket agent and cgroup controls remain proposals until explicitly ratified. Use semantic document versions; immutable tagged release specifications and migration plans accompany shipped builds.

## 2. Release train, scope ladder and dependencies

### 2.1 Phase definitions

| Phase | Intended outcome | Primary acceptance | Explicit exclusions / constraints |
|---|---|---|---|
| R0 — Architecture and feasibility | Approved ADRs, executable skeleton, resource baseline and threat model | Signed architecture, repeatable dev environment, baseline benchmarks | No production customer hosting |
| R1 — Secure single-server hosting MVP | Guided install, login, site/domain, Nginx, TLS, file/deploy path, logs, backups and recovery | Fresh-user site live, verified rollback and restoration | No provider HA, no autonomous AI |
| R2 — Hosting completeness | Runtimes, databases, DNS/email as gated modules, Git deploy, monitoring, accounts, packages, quota and CLI | Supported modules pass integration and recovery tests | Heavy modules unavailable on unsupported Nano hosts |
| R3 — Commercial and ecosystem | M-Pesa billing, licensing, reseller workflows, ServiceTemplate catalog, platform services, self-hosted apps | Reconciled ledger, tenant isolation, signed templates, upgrade and restore evidence | Stripe only if independently qualified; remote agent may remain gated |
| R4 — Provider and advanced scale | Remote mTLS agents, fleet operations, advanced provider capabilities, optional HA | Chaos/failover, fleet tenancy, support and compliance review | HA is never claimed before measured failover |

Release numbers describe dependencies rather than calendar promises. An individual module may be gated behind a capability flag, commercial entitlement and support tier. A roadmap item is not equivalent to a generally available product promise.

### 2.2 Dependency graph

`ADRs + threat model + capability schema → control-plane identity/policy → durable operation engine → local agent allowlist → ingress/site/TLS → backup/restore → hosting modules → metering/ledger → commercial workflows → signed template ecosystem → remote agents/provider scale`.

Cross-cutting work (observability, audit, migrations, docs, test harnesses, UX accessibility and supply-chain security) starts at R0 and is required for every subsequent phase.

### 2.3 Stop/go checkpoints

**Gate G0:** approve architectural decisions and threat boundaries. **G1:** complete Nano feasibility and host compatibility matrix. **G2:** pass local single-server security, site lifecycle and disaster-restore tests. **G3:** qualify each additional service adapter. **G4:** reconcile payments and reseller boundaries. **G5:** prove remote-agent/provider failover. No marketing claim may exceed the evidence from the relevant gate.

## 3. Repository organization and module ownership

A proposed Go monorepo permits shared types without granting the public API direct root privileges:

```text
nexuspanel/
  README.md  LICENSES/  SECURITY.md  CONTRIBUTING.md
  docs/volumes/  docs/adr/  docs/runbooks/  docs/support-matrix/
  cmd/control-plane/  cmd/host-agent/  cmd/npctl/  cmd/installer/
  internal/auth/  internal/policy/  internal/tenancy/
  internal/capabilities/  internal/operations/  internal/audit/
  internal/sites/  internal/domains/  internal/tls/
  internal/runtimes/  internal/databases/  internal/email/
  internal/backups/  internal/commerce/  internal/ledger/
  internal/templates/  internal/provider/  internal/observability/
  internal/agentclient/  internal/hostops/
  api/openapi/  api/events/  api/agent-rpc/
  web/  packages/design-system/
  migrations/sqlite/  migrations/postgres/
  templates/manifests/  templates/verification/
  packaging/systemd/  packaging/deb/  packaging/rpm/
  scripts/dev/  scripts/ci/  scripts/bench/
  tests/contracts/  tests/e2e/  tests/security/  tests/chaos/
  tests/fixtures/  tests/performance/
  deployments/local/  deployments/provider/
  .github/workflows/
```

The names are illustrative, not an instruction to scaffold files before approvals. Enforce domain boundaries: commerce cannot issue privileged host commands directly; UI cannot read the agent socket; agent cannot mint permissions; migrations cannot perform network calls; templates cannot access raw host credentials. Define package dependency checks and lint rules to prevent accidental cyclic coupling.

**Ownership model:** each module has a primary engineering owner, backup reviewer, security reviewer, test owner and on-call runbook owner. A solo developer may occupy multiple roles but must still use independent checklists and staged sign-offs for destructive changes.

### 3.1 Code review policy

Require review for privileged commands, auth, payment state transitions, schema migrations, package scripts, template trust policy, firewall changes and backup/restore. For a solo operating model, use a written self-review checklist and a deliberate review interval for high-risk changes; automated checks are necessary but not a substitute for independent external security assessment before commercial launch.

## 4. Environment strategy and reproducible development

| Environment | Purpose | Data policy | Promotion |
|---|---|---|---|
| Developer workstation / disposable VM | Fast iteration, local mock agent, fake payment adapters | Synthetic only | Unit/contract tests |
| Ephemeral CI VM | Install from clean supported OS, upgrade and teardown | Synthetic only | Automated integration gates |
| Nano benchmark VM | Resource/CPU pressure and supported-minimum validation | Synthetic only | Signed benchmark reports |
| Standard integration VM | Docker, PostgreSQL, email, DB, service-template testing | Synthetic only | Module acceptance |
| Preproduction | End-to-end provider, TLS, DNS and payment sandbox | Masked/synthetic only | Release-candidate gates |
| Production | Supported paying customers | Governed retention, encrypted backups | Controlled canary and rollback |

**Host matrix:** test every advertised OS/version, architecture, kernel, cgroup mode, init system, firewall backend, package manager and virtualization mode. Ubuntu 26.04.1 on the current development VPS is an observed context, not proof of a certified support matrix. A clean-room CI image MUST match the advertised OS support. Record unsupported kernel features and provider restrictions as capability failures, not installation bugs to conceal.

Use mock payment webhooks and fake DNS providers in ordinary CI; dedicated sandbox credentials in restricted integration jobs. Never store real payment tokens, production secrets, customer data or private signing keys in git or CI logs.

## 5. Engineering work packages and implementation sequence

### WP-01 — Foundations and decisions (R0)

- Ratify Go version policy, license model, Linux support matrix, Nano installation floor, process topology and data-store selection.
- Define capability schema, feature status taxonomy, organization/account ownership model, role/permission matrix and authoritative `GET /api/...` bootstrap contracts. `GET /` is the browser route, not a privileged API.
- Publish architecture diagram, threat model, signed-asset trust roots, secrets strategy, package update policy and support boundaries.
- Create the reproducible dev harness, CI checks, test fixtures, ADR process and requirements traceability manifest.
- **Exit:** reviewed ADRs, compile/test pipeline, empty agent attack-surface inventory, documented unresolved decisions.

### WP-02 — Identity, tenants and policy (R0–R1)

- Implement authentication, session rotation, MFA recovery, user/organization membership and server-side policy enforcement.
- Separate persona, presentation mode, role, capability, commercial entitlement and host resource availability in both API and UI.
- Make denial reasons machine-readable without leaking privileged internal details.
- **Exit:** cross-tenant authorization tests, session invalidation, permission-change propagation, audit correlation.

### WP-03 — Operations and agent boundary (R0–R1)

- Define typed agent RPC with peer identity, protocol negotiation, replay protection, command allowlists, path validation and timeout budgets.
- Implement durable operation queue and explicit states: `queued`, `preflight`, `running`, `verifying`, `succeeded`, `failed`, `cancelled`, `rollback_pending`, `rolled_back`, `needs_attention` (exact enum requires cross-volume ratification).
- Enforce operation ownership, idempotency, locks, leases, compensations and crash recovery; support read-only mode during incompatible upgrades.
- **Exit:** malicious-input harness, forced process death, duplicate request and interrupted-reconciliation tests.

### WP-04 — Nano install, capability detection and resource policy (R1)

- Installer preflight checks disk, RAM, CPU, swap, cgroup v2, permissions, package locks, network ports and existing services.
- Separate idle footprint, steady-state footprint, peak provisioning and customer-workload budgets. Resource ceiling is not actual usage.
- Use cgroup v2 `cpu.weight` for proportional sharing, `memory.max` for enforceable limits and `cpu.max` only where a hard cap is warranted; verify behavior under contention.
- Exclude Docker and PostgreSQL from default Nano footprint. Require a documented graceful refusal for unsupported services.
- **Exit:** 512 MB / fractional-vCPU / 10 GB disk target tested with representative workloads, OOM/pressure recovery and installer rollback. **Target remains unproven until measurements exist.**

### WP-05 — Website lifecycle (R1)

- Site create/edit/delete, domain verification, Nginx config generation and validation, atomic reload, ACME issuance/renewal, redirects and log visibility.
- Enforce account ownership, directory confinement, safe permissions, rollback on invalid config and certificate expiry alerts.
- Guided UI routes from desired workload to live HTTPS site with progress, failure recovery and next steps.
- **Exit:** working site, valid TLS, no cross-tenant path access, safe rollback, recoverable site deletion.

### WP-06 — Backup, restore and observability (R1)

- Backup catalog, encryption, integrity manifest, retention, offsite target and restore to a clean host.
- Operational metrics, structured logs, audit events, notifications, health checks, disk-pressure behavior and diagnostic bundles.
- **Exit:** independently verified backup restore, RPO/RTO measurement, alert routing, support-ready runbooks.

### WP-07 — Hosting breadth (R2)

- PHP/Node/Python runtimes, Git deployments, SSH/SFTP, file editor, databases, DNS, email and scheduled jobs as separate capability-gated adapters.
- Test per-site runtime selection, process isolation, database privilege boundaries, DNS drift and mail deliverability dependencies.
- **Exit:** module-specific smoke, negative, migration, restore and uninstall tests; no unsupported feature appears as ready.

### WP-08 — Accounts, quota and API/CLI (R2)

- Hosting accounts, packages, reseller-ready ownership, quota reservation, public API, CLI and webhook signature framework.
- Test consistency between UI, CLI and API enforcement; distinguish quota from physical capacity and paid entitlement.
- **Exit:** concurrency-safe reservations, policy parity, contract tests and documented API versioning.

### WP-09 — Commercial operations (R3)

- Independent double-entry-style ledger semantics as defined in Volume 09, invoices, subscription state, verified M-Pesa Daraja callbacks, reconciliation and dispute/refund workflows.
- Licensing outage grace policy, reseller delegation, white-label scoping and customer suspension/reactivation safety.
- **Exit:** duplicate/reordered callbacks, pending/failed/reversed payments, reconciliation drift, multi-tenant leakage and accounting evidence tests.

### WP-10 — Platform services and self-hosted catalog (R3)

- Signed `ServiceTemplate` manifests, pinned artifacts, SBOM/provenance, compatibility checks, deployment plans, upgrade policy and rollback.
- Object storage, media transforms, delivery, cache/queue/search and named self-hosted applications have dedicated stateful lifecycle adapters.
- **Exit:** supply-chain verification, persistent-data preservation, backup/restore and resource/capability gating per template.

### WP-11 — Remote agents and provider fleet (R4)

- mTLS identity and rotation, remote agent enrollment/revocation, fleet desired-state reconciliation, partial-failure handling and provider support workflows.
- Introduce HA only after consensus, database failover, lease ownership, split-brain prevention and recovery tests are defined.
- **Exit:** network partition tests, revoked-cert behavior, rolling upgrade compatibility and measured failover.

### WP-12 — AI diagnostics and extension governance (phased)

- Read-only diagnostics, scoped evidence retrieval, privacy controls, action proposals with human approval and auditable outcomes.
- Extension permission manifests, signature policy, sandboxing, version compatibility and revocation.
- **Exit:** prompt-injection resistance, secrets redaction, no unapproved privileged execution and extension escape tests.

## 6. Work breakdown, prioritization and scheduling

### 6.1 Dependency-aware backlog

Represent each epic as: `Epic ID | NP requirement IDs | user value | phase | dependencies | estimate range | risk | security impact | acceptance tests | owner | status`. Prioritize first by critical-path and security/recovery obligations, then customer value. A high-revenue feature does not bypass the host-agent safety gate.

**Suggested initial 12 iterations (planning units, not committed weeks):**

| Iteration | Primary objective | Demo/evidence |
|---|---|---|
| 01 | Decisions, repository, CI and test fixtures | ADRs, reproducible build |
| 02 | Capability and permission models | Policy fixtures, capability API |
| 03 | Control plane, SQLite and migrations | Startup, recovery, migration test |
| 04 | Local agent RPC and allowlist | Privilege-boundary adversarial tests |
| 05 | Durable operation orchestration | Kill/retry/idempotency demonstration |
| 06 | Nano installer and resource preflight | Clean VM install/uninstall and metrics |
| 07 | Site, account and Nginx lifecycle | Isolated HTTPS site provision |
| 08 | Domain, ACME and DNS workflows | Renewal and failure-recovery tests |
| 09 | File/deploy workflow and guided UX | Beginner end-to-end usability run |
| 10 | Backup, restore and operational alerts | Clean-host restoration evidence |
| 11 | Security hardening and performance | Threat tests and measured budgets |
| 12 | Release-candidate audit and documentation | R1 go/no-go package |

If an iteration fails its exit gate, reschedule downstream dependent work rather than weakening the gate. R2–R4 should be planned after actual R1 throughput and support burden are measured.

### 6.2 Definition of ready

An item is ready when requirement IDs, user outcome, policy/capability behavior, data impact, API/RPC contracts, rollback plan, negative tests, ownership and documentation are clear. If it introduces a new privilege, secret, network exposure or billable event, security and commerce review are prerequisites.

### 6.3 Definition of done

Merged code; static and dynamic checks pass; schema/API compatibility verified; accessibility tested; logs/metrics/audit wired; backup/restore impact assessed; runbook and support notes written; documentation updated; traceability evidence attached; rollout and rollback procedures reviewed; no untriaged critical or high vulnerabilities.

## 7. Quality strategy and test pyramid

| Level | Required coverage | Primary evidence |
|---|---|---|
| Unit/property | Parsers, permissions, quotas, state machines, path validation, money math | Deterministic tests, fuzz seeds |
| Contract | OpenAPI, RPC versioning, event schema, webhook signatures, templates | Consumer/provider contract fixtures |
| Integration | SQLite/PostgreSQL parity, agent↔control plane, Nginx, ACME, DNS, DB, M-Pesa sandbox | Ephemeral VM logs, trace IDs |
| End-to-end | Guided onboarding, HTTPS site, deploy, backup/restore, billing and reseller flows | Browser/API replay and artifacts |
| Security | Cross-tenant access, privilege escalation, SSRF, traversal, injection, supply-chain | Threat-case reports and remediation |
| Reliability | Crashes, retries, partial writes, disk full, clock skew, network partition, expired cert | Fault-injection timeline |
| Performance | Idle, steady, burst, memory pressure, concurrent sites, cold start, restore time | Reproducible benchmark reports |
| Accessibility | Keyboard, screen reader, contrast, responsive 320px, error clarity | WCAG 2.2 AA review |
| Disaster recovery | Fresh-host restore, encrypted backup, key rotation, payment reconciliation | Timed exercise and signed evidence |

**Coverage policy:** require 100% of privileged operation types to have at least one positive, one authorization-negative, one malformed-input and one interruption test. Coverage percentage alone is not a release gate; mutation and fault tests must exercise invariants. Each supported installation profile gets an end-to-end test suite.

### 7.1 Test data discipline

Synthetic tenants and payment references; seeded deterministic clock; fake DNS provider; controlled ACME staging; sandbox M-Pesa; mock email SMTP; local object store; test signing root. Avoid real customer identifiers and never replay production secrets. Use fixtures for unauthorized, capability-unavailable, not-shipped, degraded and offline states separately.

### 7.2 Test identifiers

Use `QA-AUTH-*`, `QA-AGENT-*`, `QA-OPS-*`, `QA-NANO-*`, `QA-SITE-*`, `QA-BACKUP-*`, `QA-PAY-*`, `QA-TEMPLATE-*`, `QA-FLEET-*`, `QA-UX-*`. Each maps to the corresponding `NP-*` requirement and volume-level test where available. A test ID without an asserted invariant is not meaningful evidence.

## 8. Concrete conformance matrix

Every row is a **minimum** required scenario; the detailed test harness must add setup, fixture, teardown, exact expected API/RPC codes and evidence pointers. Failures are release blockers where they violate security, integrity or data recovery.

| ID | Scenario | Required result |
|---|---|---|

| QA-AUTH-001 | Tenant A token reads Tenant B site | 403/404 policy-safe response; no existence leak; audit event |

| QA-AUTH-002 | Role revoked during active session | Next privileged operation denied; session policy refreshed |

| QA-AUTH-003 | Persona set to Expert | No change in API permissions or entitlements |

| QA-AUTH-004 | MFA recovery attempted with stale code | Denied; replay prevented; event recorded |

| QA-AGENT-005 | Non-allowlisted method invoked | Denied before host mutation; audit and RPC error |

| QA-AGENT-006 | Symlink swaps during file write | Path confinement maintained; no cross-tenant write |

| QA-AGENT-007 | Control plane restarted mid-operation | Operation reconciles without duplicate host side effect |

| QA-AGENT-008 | Agent socket accessible to unprivileged tenant | Access denied by filesystem/peer checks |

| QA-AGENT-009 | Malicious Nginx configuration payload | Config validation fails and active config remains |

| QA-OPS-010 | Identical idempotency key replayed | One logical operation and one host effect |

| QA-OPS-011 | Queue worker killed after side effect | Reconciliation detects effect; no double application |

| QA-OPS-012 | Cancel requested during irreversible stage | Honest state and documented compensation |

| QA-OPS-013 | Database transaction fails after reservation | Reservation released or reconciled deterministically |

| QA-NANO-014 | Install on 512 MB target | Measure footprint and success/failure; no unverified support claim |

| QA-NANO-015 | Memory pressure during certificate issuance | Bounded peak, recovery and no silent OOM loop |

| QA-NANO-016 | Docker-only template selected | Unavailable with precise reason and alternatives |

| QA-NANO-017 | CPU contention between two workloads | Weighted sharing measured; batch caps obeyed |

| QA-NANO-018 | 10 GB disk nearly full | Preflight refuses unsafe install/upgrade; backup protected |

| QA-SITE-019 | Create site with existing document root | No destructive overwrite; explicit conflict |

| QA-SITE-020 | TLS renewal fails | Warning, retry, escalation; old cert preserved until expiry |

| QA-SITE-021 | Nginx reload fails | Atomic rollback and service remains serving |

| QA-SITE-022 | Domain ownership proof absent | Public exposure/issuance denied as policy requires |

| QA-SITE-023 | Git deploy interrupted | Previous deploy retained or documented recovery |

| QA-DATA-024 | SQLite migration interrupted | Transactional recovery or safe rollback; backup verified |

| QA-DATA-025 | PostgreSQL migration against old app | Compatibility matrix blocks unsafe pairing |

| QA-DATA-026 | Tenant delete with retained billing ledger | Retention rules enforced; PII minimization verified |

| QA-DATA-027 | Secret appears in backup metadata | Rejected/redacted; encrypted backup verified |

| QA-BACKUP-028 | Restore from clean host | Site, cert, data, metadata and ownership verified |

| QA-BACKUP-029 | Backup checksum corrupted | Restore refuses corruption and alerts |

| QA-BACKUP-030 | Backup key rotated | Old backups decrypt under documented key retention |

| QA-BACKUP-031 | Disk full during backup | Incomplete artifact not advertised as restorable |

| QA-PAY-032 | M-Pesa callback delivered twice | Exactly one settled ledger effect |

| QA-PAY-033 | Callback with invalid signature/origin proof | Rejected without granting entitlement |

| QA-PAY-034 | Out-of-order reversal after settlement | Compensating ledger entries; no silent deletion |

| QA-PAY-035 | Payment pending and user retries | No double invoice or premature provisioning |

| QA-PAY-036 | Provider API outage | Pending state and reconciliation; no invented success |

| QA-PAY-037 | Reseller credits another tenant | Denied; ledger/tenant isolation preserved |

| QA-TEMPLATE-038 | Manifest signature invalid | Install refused before download/execution |

| QA-TEMPLATE-039 | Upgrade stateful app with failed health check | Data preserved; safe rollback/recovery instructions |

| QA-TEMPLATE-040 | Template requests unsupported memory | Capability refusal with resource delta |

| QA-TEMPLATE-041 | Image transform decompression bomb | Bounded CPU/memory/pixel limits; safe error |

| QA-TEMPLATE-042 | Object storage key exposed | Scope, rotation, audit and revocation tested |

| QA-FLEET-043 | Remote agent certificate revoked | Subsequent connections denied and surfaced |

| QA-FLEET-044 | Control plane loses agent connectivity | Offline state; no false completed operation |

| QA-FLEET-045 | Network split during leader change | No concurrent conflicting host mutation |

| QA-UX-046 | Unauthenticated vs unavailable vs unshipped | Distinct states and next actions; no policy leakage |

| QA-UX-047 | Mobile 320px emergency certificate renewal | Usable controls, focus order and feedback |

| QA-UX-048 | Beginner guided site onboarding | Task completion and error recovery measured |

| QA-UX-049 | Long-running job stalls | Progress remains honest; retry/support action shown |

| QA-SEC-050 | Prompt injection requests shell access | AI returns diagnostic/proposal only; no privileged call |

| QA-SEC-051 | Extension attempts secret read | Permission/sandbox enforcement and audit |

| QA-SEC-052 | Unauthorized API token export | Denied and secret not logged |

| QA-SEC-053 | Offsite backup destination compromised | Encryption/integrity protects recoverable data |

## 9. Performance and Nano benchmark protocol

### 9.1 Benchmark method

Publish immutable test configuration: host provider/type, actual RAM and CPU quota, disk type and free capacity, OS/kernel, swap, cgroup settings, Go build, agent/control-plane versions, services enabled, load generator, workload mix and sample duration. Record warm and cold start, p50/p95/p99 latencies, CPU usage, RSS and cgroup memory, disk growth, IO wait, OOM kills, queue depth, failure rate and restoration time. Compare at least three independent clean installs, then repeat after upgrade.

### 9.2 Test tiers

- **Nano candidate:** 512 MB RAM, fractional/1 vCPU, 10 GB disk, SQLite, local agent, static SPA and minimal ingress. No default Docker, PostgreSQL, heavy email stack or search engine.
- **Starter:** more headroom, additional optional runtimes; exact published floor must be measured.
- **Standard:** container-enabled workloads, PostgreSQL where selected, richer service adapters.
- **Provider:** fleet and reseller workload, multiple agents, isolation and capacity oversubscription controls.

### 9.3 Proposed acceptance measurements

| Measure | Evidence requirement | Decision rule |
|---|---|---|
| Idle panel RAM | RSS and cgroup memory after warmup | Publish measured baseline and reserve margin |
| Peak install/upgrade RAM | Max memory under constrained host | No repeated OOM, failed install or hidden swap dependency |
| Basic site provisioning | p50/p95, retries and completion | Approve target after actual runs |
| TLS issuance/renewal | Success, retry/backoff and CPU peak | Safe under ACME limits and pressure |
| Backup and restore | Bytes, elapsed time, checksums, RPO/RTO | Targets approved by data class |
| CPU contention | `cpu.weight` distribution under simultaneous loads | Fairness and starvation behavior documented |
| Disk pressure | Remaining free bytes at failure threshold | Installer refuses unsafe state and explains recovery |
| Guided onboarding | Timed novice task study | Under-10-minute aspiration validated or revised |

**No blanket guarantee of 512 MB compatibility** is permitted before the full Nano suite passes on representative hardware and published workload constraints. A CPU cap is not equivalent to a performance guarantee; quotas are ceilings and utilization is measured separately.

## 10. CI/CD, artifacts and software supply chain

### 10.1 Pull-request pipeline

1. Format, lint, Go vet/staticcheck, TypeScript checks, dependency license and vulnerability scan.
2. Unit/property/fuzz regression, API schema compatibility, RPC and event-contract tests.
3. SQLite/PostgreSQL migration test, policy fixtures, security regression and agent allowlist tests.
4. Frontend accessibility smoke, browser critical path and UI state snapshots.
5. Reproducible build, SBOM, provenance statement and signed artifacts in isolated signing environment.
6. Ephemeral VM install/upgrade/rollback smoke for each release-supported OS/profile.

### 10.2 Release pipeline

`commit → candidate tag → reproducible build → signed binaries/packages/manifests → staging install → upgrade from previous supported versions → backup/restore → security scan → canary → progressive rollout → final sign-off`.

Publish checksums, signing identity, SBOM, dependency manifest, release notes, API/migration compatibility table, known limitations, recovery instructions and support window. Never require an installer to execute an unverified remote script as root. Package repositories and update metadata must be authenticated and rollback-safe.

### 10.3 Version compatibility

| Component | Compatibility rule |
|---|---|
| Frontend ↔ API | Negotiated API version; older browser assets handled explicitly |
| Control plane ↔ agent | Protocol version handshake; reject unsafe mismatch |
| Application ↔ schema | Expand/migrate/contract phases; no blind downgrade |
| Template manifest ↔ runtime | Schema version and capability checks before install |
| Extension ↔ panel | Permission and API compatibility review |
| License grant ↔ feature | Signed claim validation; outage grace per policy |

Release notes MUST identify irreversible schema changes and required backups. Rolling upgrades are only allowed where mixed-version behavior has been tested.

## 11. Installer, uninstaller and host compatibility

### 11.1 Installer lifecycle

`discover → explain → preflight → preview plan → confirm → snapshot/backup → install packages → create service accounts → configure agent socket → initialize datastore → configure ingress → health verify → handoff`.

The installer MUST enumerate port conflicts, existing Nginx/Apache, DNS and mail services, firewall state, existing databases, resource deficits and privileged changes before applying them. Require explicit user confirmation for potentially disruptive actions. Never delete preexisting websites or databases to satisfy a default install plan.

### 11.2 Safe uninstall

Offer keep-data vs purge-data paths; default to retaining customer content and encrypted backups. Explain which accounts, packages, firewall rules, services and certificates will be removed. Record a removal manifest, require explicit destructive confirmation and verify that unrelated system services remain unchanged.

### 11.3 Upgrade and downgrade

Upgrade preflight validates free space, backups, protocol versions, schema migration path, service impact and rollback feasibility. For irreversible migrations, offer restore-from-backup rather than promising in-place downgrade. On failed verification, leave the panel in a known safe maintenance state and provide a runbook.

## 12. Data migration, schema evolution and recovery

- Maintain ordered, checksum-verified SQLite and PostgreSQL migration histories with matching domain semantics, not necessarily byte-identical SQL.
- Each migration specifies online/offline mode, lock behavior, estimated data growth, rollback strategy, backup requirement and compatibility range.
- Use expand/contract for zero-downtime-capable changes; never drop a column in the same release that first stops reading it unless mixed-version deployment is prohibited and tested.
- Backfill jobs must be idempotent, resumable, bounded and observable; tenant ownership remains explicit.
- Test corrupted SQLite WAL, failed filesystem sync, PostgreSQL connection exhaustion, schema drift and migration interruption.
- Validate recovery of financial ledger and audit chains independently of ordinary service data.

**MIG-001 acceptance:** install from empty schema; migrate from each supported prior release; interrupt each destructive step; restore snapshot; reconcile operation and ledger states; verify data invariants and UI/API behavior. A migration that passes only on empty databases is not qualified.

## 13. Operational model, SLOs and on-call

### 13.1 Service indicators

| Service | Proposed SLI | Initial SLO approach |
|---|---|---|
| Panel API | Successful authenticated requests excluding client errors | Set after baseline and dependency analysis |
| Privileged operations | Correctly completed/failed with reconciled terminal state | Track latency and stuck operation rate |
| Site ingress | HTTPS availability per site | Separate panel outage from hosted-site outage |
| TLS | Certificates renewed before expiry | Alert well before expiration; retry safely |
| Backups | Scheduled jobs completed and verified restorable | Restore drills determine effective reliability |
| Payment | Verified event-to-ledger reconciliation delay | Never confuse provider callback receipt with settlement |
| Fleet | Connected agents and desired-state convergence | Explicit offline/partition states |

Avoid arbitrary 99.9% claims until the architecture, maintenance windows, dependencies and monitoring support them. Document service ownership, alert thresholds, escalation path, maintenance communication and error budgets.

### 13.2 Observability baseline

Every operation has trace ID, actor/tenant ID (minimized), target resource, operation state, duration, safe error code, retry count and audit link. Never log secrets, raw ID documents, payment tokens or sensitive file contents. Metrics distinguish panel overhead from customer workload usage. Monitor agent liveness, queue depth, memory pressure, disk pressure, database health, ACME errors, backup freshness, SMTP queue and payment reconciliation backlog.

### 13.3 Incident severities

- **SEV-1:** widespread hosted-site outage, active compromise, irreversible data loss, systemic payment misposting.
- **SEV-2:** critical function unavailable, tenant isolation suspected, backups unavailable, prolonged payment reconciliation.
- **SEV-3:** partial degradation with workaround, noncritical adapter failure, delayed notifications.
- **SEV-4:** cosmetic defects, low-impact documentation issues.

Severity assignment is provisional and must be calibrated to actual customer contracts. Every SEV-1/2 requires an incident commander, timeline, evidence preservation, customer communication decision, root-cause analysis and tracked corrective actions.

## 14. Runbook library and operational procedures

| Runbook | Trigger | Required safe actions | Recovery evidence |
|---|---|---|---|
| OPS-001 Agent unreachable | Agent heartbeat or RPC failures | Check socket/service identity, logs, protocol; do not weaken permissions | Agent handshake and operation reconciliation |
| OPS-002 Control plane crash loop | Repeated restart | Preserve logs, check migrations/secrets/disk, safe mode | Health checks and no orphan mutations |
| OPS-003 Disk full | Free-space threshold | Pause nonessential writes, preserve backups, identify owner of growth | Restored margin, integrity check |
| OPS-004 OOM/memory pressure | OOM kill or high PSI | Stop batch work, inspect cgroups, downgrade workload or refuse capacity | Stable memory and service health |
| OPS-005 TLS renewal failure | Renewal deadline alert | Check DNS/ACME reachability, rate limits, ownership; safe retry | Valid certificate served |
| OPS-006 Nginx invalid config | Reload verification fails | Keep last known good config, inspect generated diff | Config test and site availability |
| OPS-007 Failed deployment | Health check regression | Roll back application artifact, preserve persistent data | Previous release healthy |
| OPS-008 Backup corruption | Checksum/restore test fails | Quarantine artifact, use alternate backup, alert | Full clean restore |
| OPS-009 Database migration failure | Migration error or checksum mismatch | Enter maintenance, snapshot, follow compatibility plan | Data invariants and application readiness |
| OPS-010 M-Pesa discrepancy | Ledger/provider mismatch | Freeze affected settlement action, reconcile verified provider records | Balanced ledger and case closure |
| OPS-011 Suspicious privileged action | Agent/audit anomaly | Isolate access, rotate credentials, preserve forensic evidence | Security sign-off and replay check |
| OPS-012 Tenant boundary incident | Unauthorized data access suspected | Disable affected path, preserve audit, assess notification duties | Negative test and incident review |
| OPS-013 Template supply-chain alert | Compromised image/signature | Revoke catalog entry, block new installs, assess running instances | Verified replacement and integrity |
| OPS-014 Remote agent partition | Fleet loses connectivity | Mark offline, avoid conflicting desired-state writes | Rejoin/reconcile evidence |
| OPS-015 Payment provider outage | Gateway unavailable | Queue/reconcile, show pending state, avoid premature activation | Verified settlement after recovery |
| OPS-016 Disaster restore | Host unrecoverable | Provision clean supported host, restore encrypted backup, rebind secrets/DNS | Timed full service recovery |
| OPS-017 Lost signing key | Signing trust compromise | Revoke key, rotate trust roots, quarantine releases | Verified new chain and customer advisory |
| OPS-018 License validation outage | Licensing endpoint unavailable | Apply documented offline grace, preserve existing hosting | Reconciliation without service disruption |

Every runbook MUST include: preconditions, risk, required privilege, exact observable indicators, safe stop condition, owner, rollback, verification and escalation. Shell commands are to be developed and tested on disposable hosts before inclusion; do not paste unreviewed root commands into production procedures.

## 15. Backup, disaster recovery and resilience exercises

**Data classes:** control-plane metadata, customer files, service databases, DNS and TLS state, payment ledger, signing/secret material, template manifests, agent identities, audit records and monitoring configuration. Each class needs an owner, encryption policy, backup cadence, retention, geographic placement and recovery objective.

**Drills:** monthly sample restore (proposed), quarterly clean-host full restore (proposed), and pre-release restore for every migration or backup format change. Simulate loss of the primary VPS, loss of the primary datastore, corrupt backup, leaked key, unavailable DNS and inaccessible payment gateway. Verify customer-visible behavior and not just database process startup.

**RPO/RTO:** choose targets per data class after cost/risk approval. Financial ledger integrity and audit continuity may demand stronger guarantees than cached metrics. A successful backup job does not establish a successful recovery; only a verified restore does.

## 16. Security, privacy and compliance release gates

- Maintain threat model for internet API, agent socket, host root boundary, customer workloads, templates, payment webhooks, AI prompts, extensions and provider fleet.
- Require independent penetration test before commercial multi-tenant/provider launch; remediate critical/high findings or document an explicit release block.
- Inventory personal data, retention/deletion behavior, encryption, processor/subprocessor responsibilities and breach response obligations. Kenya-specific legal compliance requires qualified legal review; documentation is not certification.
- Confirm default-deny network exposure, secure SSH practices, package patching, session protections, MFA options, least-privilege secrets, tamper-evident audit and backup key recovery.
- Exercise emergency revocation of tokens, templates, signing keys, remote agents and reseller access.
- AI diagnostic outputs must redact secrets and require approval for any separately authorized action; prompt content never grants permissions.
- Use a vulnerability disclosure policy, patch severity windows and security advisories; do not promise timelines without staffing and response capability.

## 17. Documentation, support and developer experience

**User-facing:** installation prerequisites, Nano/Starter/Standard/Provider capability matrix, guided setup, common error explanations, restore instructions, pricing and licensing terms, API reference, CLI reference and accessible UX guidance.

**Operator-facing:** supported OS matrix, deployment architecture, migration runbooks, agent troubleshooting, security incident playbook, billing reconciliation, reseller support and fleet maintenance.

**Developer-facing:** build prerequisites, architecture ADRs, schema and RPC conventions, test fixtures, template authoring/signing, code ownership and contribution process.

Support responses must not imply a feature is available when it is merely planned, disabled by capability, restricted by role or degraded by an outage. Publish release status separately from product-area scope. Maintain an issue-to-requirement mapping for customer-reported defects.

## 18. Canonical 29-area delivery traceability and release ownership

The area labels below are taken from the authoritative v2.0 functional-scope register. The proposed phase assignments are delivery-planning interpretations; they **do not override** the source register’s v1–v4 roadmap or Volume 01 requirements. The actual mapping from each area to `NP-*` identifiers must be generated and approved before implementation.

| Area | Canonical scope label | Proposed delivery window | Key verification |
|---|---|---|---|
| 01 | Dashboard | R1 | Dashboard metric provenance and emergency mobile |
| 02 | Server management | R1 | Service operations, network changes and recovery |
| 03 | Websites | R1 | Website isolation and safe Nginx reload |
| 04 | DNS | R1 guidance / R2 authoritative | Delegation, DNS drift and zone recovery |
| 05 | SSL/TLS | R1 basic / R2 advanced | ACME renewal and HSTS safety |
| 06 | File manager | R1 | Path confinement, uploads and safe editor saves |
| 07 | Application manager | R2 | Application lifecycle, secrets and health checks |
| 08 | Runtime manager | R2 | Runtime version compatibility and rollback |
| 09 | Databases | R2 | Database user isolation, quota and restore |
| 10 | Email hosting | R2 external / later self-hosted | Mailbox isolation, queues and delivery |
| 11 | Git and deployments | R2 | Webhook signatures, build isolation and rollback |
| 12 | Backups | R1 | Offsite integrity and clean-host restore |
| 13 | Security | R1 | Threat model, firewall and MFA recovery |
| 14 | Monitoring | R1 | Alerting, per-site attribution and retention |
| 15 | Terminal/developer tools | R2 | Terminal scope, SSH keys and session audit |
| 16 | User/account management | R1 | Role boundaries, account lifecycle and recovery |
| 17 | AI assistant | R1 diagnostics / later actions | Read-only AI, redaction and prompt injection |
| 18 | Automation/API | R1 core / R2 full | API policy parity, job idempotency and webhooks |
| 19 | Panel administration | R1 | Admin actions, update safety and recovery |
| 20 | Extensions | R3 | Extension signing, permissions and sandbox |
| 21 | Commerce/reseller | R3 commercial / R4 provider | Ledger, payment, reseller delegation and licensing |
| 22 | Migration/staging | R2 | Migration fidelity, staging isolation and rollback |
| 23 | Containers/orchestration | R2 internal / R4 user-facing | Docker isolation, image trust and capacity |
| 24 | Performance/caching | R2 | Cache invalidation, pressure and site latency |
| 25 | Email security/deliverability | R3 | SPF/DKIM/DMARC and abuse handling |
| 26 | Analytics/site quality | R3 | Metrics accuracy, privacy and report attribution |
| 27 | Audit/compliance/support | R1 audit / R3 full | Tamper-evident audit, support access and retention |
| 28 | Platform Services | R3 | Storage, transforms, delivery, cache/queue/search |
| 29 | Self-hosted Templates | R3 | Template signing, persistent data and upgrade restore |

**Traceability rule:** Every canonical area maps to its requirements, modules, UI states, API/RPC, data objects, security controls, acceptance tests and release evidence. Shared architectural spines apply across all areas; no row by itself promises availability on Nano or in the first release.

## 19. Commercial launch, support readiness and customer protection

Before charging customers, approve: legal entity and invoicing disclosures; privacy policy and terms; product support tiers; payment reconciliation; refund/dispute procedure; customer data export; suspension policy; disaster recovery capability; vulnerability disclosure; and operational escalation contacts. Provider claims such as uptime, backup guarantee, data residency or HA must be evidence-backed and contractually reviewed.

For Kenya-first billing, test M-Pesa Daraja sandbox flows and independently verify production callback security and settlement reconciliation. Stripe is an additional adapter, not a substitute for gateway-independent accounting. Tax and consumer-protection decisions require jurisdiction-specific review.

**Customer protection principle:** license validation, payment gateway outage or provider API downtime must not immediately delete hosted content. Suspension and reactivation must preserve recoverability and follow an approved policy.

## 20. Release go/no-go checklist

| Gate | Required evidence | Blocking condition |
|---|---|---|
| REL-001 Product | Signed scope, phase matrix and open-decision register | Critical requirement ambiguous |
| REL-002 Architecture | Approved ADRs and trust-boundary review | Unreviewed privilege escalation path |
| REL-003 Compatibility | Certified OS/profile test matrix | Unsupported host advertised as supported |
| REL-004 Nano | Reproducible resource benchmark | Unmeasured 512 MB compatibility claim |
| REL-005 Security | Auth/agent/tenant tests and vulnerability review | Unresolved critical/high risk |
| REL-006 Data | Migration and restore drill | Irreversible unbacked schema change |
| REL-007 Operations | Idempotency, queue crash recovery and audit | Duplicate destructive side effects |
| REL-008 Hosting | HTTPS site create/edit/delete, rollback | Host service outage from invalid config |
| REL-009 UX | Guided flow, error states, WCAG review | Critical inaccessible action or misleading status |
| REL-010 Payments | Ledger reconciliation and callback replay tests | Incorrect entitlement or financial balance |
| REL-011 Templates | Signature, compatibility, upgrade/restore tests | Unsigned/unbounded privileged execution |
| REL-012 Provider | Remote mTLS and split-brain tests | Unsupported HA/fleet claim |
| REL-013 Observability | Alerts, runbooks, incident contact and logs | No way to diagnose or recover |
| REL-014 Supply chain | Signed artifacts, SBOM, provenance, rollback | Unverified update distribution |
| REL-015 Documentation | User/operator/API docs, support matrix | Hidden unsupported behavior |
| REL-016 Governance | Owner sign-off, known limitations and release notes | Open blocker concealed or misclassified |

Release approval requires a recorded date, artifact hashes, test environment versions, risk acceptance owners and linkable evidence. A release cannot be approved solely by a percentage complete.

## 21. Cross-volume interface reconciliation checklist

| Boundary | Must match | Verification |
|---|---|---|
| Vol 01 ↔ Vol 02 | Entitlements, personas, guided router, states | UI-to-policy tests |
| Vol 01 ↔ Vol 09 | Pricing, subscriptions, suspension and refunds | Commerce acceptance fixtures |
| Vol 03 ↔ Vol 06 | Agent RPC, cgroups, process and privilege boundaries | Host integration and security tests |
| Vol 03 ↔ Vol 04 | SQLite/PostgreSQL choices and durable operations | Migration/crash consistency |
| Vol 04 ↔ Vol 05 | Data ownership, idempotency and API entities | Contract and transactional tests |
| Vol 05 ↔ Vol 06 | RPC methods, status codes, replay handling | Consumer/provider schema tests |
| Vol 05 ↔ Vol 09 | Payment events, ledger and entitlements | Duplicate/reversal callback tests |
| Vol 07 ↔ Vol 06 | Template host permissions, persistence and rollback | Signed-template sandbox tests |
| Vol 07 ↔ Vol 08 | Supply-chain, secret, isolation and backup rules | Adversarial and restore tests |
| Vol 08 ↔ all | Threat model, retention, incident and audit | Security traceability matrix |
| Vol 10 ↔ all | Release phases, test IDs, decision owners | Machine-checked requirement matrix |

Reconciliation must compare exact identifiers and definitions, not merely section titles. Where a contradiction remains, record `XR-NNN` with source excerpts, impact, options and approver; leave implementation blocked until resolved. This volume deliberately does not invent an approval on the user's behalf.

## 22. Decision register and evidence backlog

| ID | Decision / validation | Required owner | Blocking milestone |
|---|---|---|---|
| DEL-001 | Confirm Go versions and package/update policy | Technical owner | G0 |
| DEL-002 | Approve open-core license and third-party obligations | Product/legal | Public distribution |
| DEL-003 | Establish real Nano minimum and service limits | Performance lead | G1 |
| DEL-004 | Confirm SQLite/PostgreSQL migration and operational split | Data owner | G0/G2 |
| DEL-005 | Approve cgroup policy, CPU sharing and memory enforcement | Host/security | G1/G2 |
| DEL-006 | Approve exact agent RPC schema and Unix-socket permissions | Security/agent | G2 |
| DEL-007 | Approve feature/capability status taxonomy | Product/API/UX | G0 |
| DEL-008 | Approve supported Linux distributions and architectures | Release engineering | G1 |
| DEL-009 | Define backup RPO/RTO and offsite key custody | Operations/security | G2 |
| DEL-010 | Confirm M-Pesa production verification and reconciliation | Finance/security | G4 |
| DEL-011 | Approve reseller delegation and suspension safeguards | Product/security | G4 |
| DEL-012 | Approve template signing roots, trust policy and review | Security/ecosystem | G4 |
| DEL-013 | Decide whether email stack is supported on each profile | Product/host | R2 |
| DEL-014 | Confirm remote mTLS identity lifecycle and fleet semantics | Security/provider | G5 |
| DEL-015 | Define measured HA topology, failover and split-brain policy | Architecture/operations | G5 |
| DEL-016 | Confirm legal/privacy requirements and support obligations | Legal/operations | Commercial launch |
| DEL-017 | Approve SLOs and on-call escalation coverage | Operations | Commercial launch |
| DEL-018 | Approve payment/provider outage grace and data retention | Product/finance | G4 |
| DEL-019 | Validate guided onboarding time with target users | UX research | R1 release |
| DEL-020 | Complete 29-area canonical ID mapping and NP traceability | Product/release | Final specification sign-off |

All unresolved entries remain explicitly open. Their presence does not authorize changes to the development VPS or constitute an engineering sign-off.

## 23. First engineering kickoff after documentation approval

1. Freeze and version the ten-volume specification, recording a SHA-256 manifest and open-decision register.
2. Approve or amend the foundational ADRs, especially Nano floor, OS support, license, process topology and agent privilege boundary.
3. Create a requirements-to-tests machine-readable index from existing `NP-*` IDs, without renumbering requirements.
4. Build a disposable VM laboratory, not the existing Leaseweb host, for installer, agent and resource experiments.
5. Implement a no-privilege skeleton and mocked end-to-end flow; measure actual process memory before enabling real host mutations.
6. Run a security architecture review and publish a go/no-go decision for the first host-operation adapter.
7. Begin R1 work only after the above evidence exists; schedule later phases from observed throughput and risk.

**Final acceptance of the master specification:** all ten volumes present; links and filenames verified; all 29 canonical areas mapped; source precedence documented; conflicting decisions explicitly tracked; every planned release has clear entry/exit gates; security, payment, data recovery and Nano claims are qualified; ZIP integrity and per-file checksums validated. A complete document set is not a claim that the software has been implemented.

