# NexusPanel — Master Product & Engineering Specification
## Volume 07 — Platform Services & Self-hosted Ecosystem

**Version:** 1.0 consolidated | **Status:** Proposed product and engineering contract; technical validation and approvals pending | **Date:** 2026-11-06  
**Authority:** Volume 01 Master PRD; Volumes 02–06; authoritative 29-area scope.  
**Source basis:** User-supplied Master Blueprint v2.0, Stage 2 PRD v2.0 (§21-bis), Stages 3–7 v2.0, 29-area register and consolidated Volumes 01–06.  
**Operational boundary:** Documentation only; no commands have been executed on the Leaseweb VPS.

## 0. Scope, terminology and precedence

This volume is the canonical detailed contract for **Area 28 — Platform Services** and **Area 29 — Self-hosted Templates**, and their dependencies on websites, applications, DNS, TLS, storage, backups, databases, containers, security, monitoring, billing and migration. It describes the common catalog/provisioning layer and the specialized service behavior that must not be flattened into a generic container lifecycle. Volume 01 governs *what* the product promises; this volume governs how these two areas are designed, surfaced, secured and verified. Conflicts are resolved against Volume 01 and recorded for approval, not hidden.

**MUST** is mandatory for a module to claim production completion. **SHOULD** is strongly recommended, with documented exceptions. **MAY** is optional. **VERIFY** means measured evidence is required; **APPROVAL** means the choice is not yet final. A template listed in the catalog is not proof of a production-ready installation. Product scope is broader than any one hardware profile or release.

### 0.1 Consolidated binding principles

1. `ServiceTemplate` is the shared *catalog and provisioning abstraction*, not a claim that all services have identical storage, credentials, consistency, or upgrade semantics.
2. The source v2.0 requirements **NP-SVC-001 through NP-SVC-006** remain traceable; expanded detailed obligations use `SVC-*` and `TPL-*` identifiers in this volume.
3. Storage, image transformation and delivery are three separately operated services with distinct cost, performance, and failure boundaries. A reverse-proxy cache is **not** a global CDN; third-party CDN integrations must be labeled as such.
4. v1 reference targets are S3-compatible storage (MinIO or Garage), image transformation (imgproxy or thumbor), and delivery (local caching plus optional Cloudflare/Bunny). Specific vendor selection remains an adapter decision, subject to maintenance and security review.
5. v1 self-hosted template targets are WordPress, Nextcloud, Gitea, Ghost, n8n and Uptime Kuma. Their presence in the *target catalog* does not override per-profile capability gates or release-quality gates.
6. Nano (proposed 512 MB target) does **not** run local Docker-based platform services. It may consume remotely hosted S3 storage, CDN and backup services using scoped credentials. Capability availability must be measured, not assumed.
7. Standard+ may run containerized services if prerequisites, reserved system headroom, cgroup policies, storage durability, security and backup coverage pass preflight. A nominal 4 GB server is not guaranteed to host every catalog item concurrently.
8. All services are tenant-scoped and governed by the same policy, durable operation, quota, audit, health and recovery model as the rest of NexusPanel.
9. Templates are signed and versioned data plus reviewed runtime adapters. Adding a catalog record may not require a control-plane release **only if** the requested capabilities and operations already exist in an allowlisted runtime.
10. Stateful data export/import is a verified migration workflow, not merely a tarball or OCI image export. Recovery, secrets, ownership and integrity must be independently verified.
11. Batch transforms, builds, migrations and backups yield to serving and infrastructure workloads under contention. Backpressure and fair queueing must be observable.
12. The product must never imply that a template's successful container start means the application is configured, secure, backed up, or ready for production.

## 1. Product outcomes, personas and user journeys

### 1.1 Jobs-to-be-done

| Persona/job | Desired outcome | Core experience | Guardrail |
|---|---|---|---|
| Beginner self-hoster | Run a useful application without shell expertise | Guided router → template → prerequisites → deploy → verify | No hidden root operations or silent unsupported features |
| Website owner | Attach media storage to an existing website | Create/select bucket → grant scoped access → connect | Never reveal reusable secrets in logs |
| Developer | Deploy repeatable stacks | Versioned templates, env schema, volumes, logs, revisions | Immutable image digests and controlled rollback |
| Small VPS owner | Know what fits before committing | Honest capability and resource-delta display | Do not install Docker or databases on Nano |
| Reseller | Offer approved templates to customers | Curated catalog, packages, quotas, delegated access | No upstream quota bypass |
| Provider | Operate service fleets and pricing | Catalog policy, image allowlists, placement, metering | Tenant isolation and audited changes |

### 1.2 Required end-to-end journeys

**FLOW-SVC-01 — Discover and provision a service:** Enter guided router → choose service goal → browse catalog filtered by entitlement and host capability → inspect data/privacy and estimated resource/cost impact → configure hostname, storage, admin identity, backup target and optional integrations → preview effective plan → authorize and approve where necessary → create durable operation → track steps → verify health and externally observable behavior → show credentials once → record audit and recovery options.

**FLOW-SVC-02 — Connect storage to an existing site:** Select site and bucket → verify common tenant/delegation → choose access level and prefix → issue a site-specific short-lived or revocable credential → inject only into managed secret configuration → test object PUT/GET under intended policy → display endpoint, rotation, last-used metadata and disconnection controls.

**FLOW-SVC-03 — Publish transformed media:** Choose source bucket → configure accepted formats, maximum dimensions and transform presets → enable signed transform URLs and cache keys → configure local delivery → optionally attach external CDN → validate TLS and cache behavior → measure origin hits, transform CPU, errors and egress.

**FLOW-SVC-04 — Upgrade a stateful app:** Review compatibility and release notes → check backup freshness and restore point → schedule maintenance/approval → snapshot and quiesce if required → apply pinned image/config migration → run application health and data integrity probes → complete or roll forward/restore using the service-specific runbook.

**FLOW-SVC-05 — Export/import to a larger VPS:** Preflight destination capacity and compatibility → take application-consistent backup → export manifest and encrypted data → transfer securely → provision destination without publishing traffic → restore and verify → switch DNS/routes with rollback window → revoke source credentials → decommission source after confirmation.

**FLOW-SVC-06 — Handle unavailable capability:** The catalog remains discoverable where policy permits, but shows the exact unmet dependency, approximate resource delta, supported remote alternative, and upgrade path; creation APIs reject unsupported local deployment without side effects.

## 2. ServiceTemplate manifest and registry contract

### 2.1 Canonical manifest fields

| Group | Required fields | Contract |
|---|---|---|
| Identity | `template_id`, `slug`, `display_name`, `version`, `schema_version`, publisher | Immutable version identity; collision-safe naming |
| Provenance | OCI digest(s), signature/attestation, license, SBOM reference | Digest pinning and trust verification before use |
| Category | `workload_category`, `service_kind`, `capability_requirements` | Serving/batch/infrastructure immutable for a deployed revision |
| Resources | floor/typical/ceiling for RAM, CPU, disk, I/O, network; headroom | Estimates separate from measured consumption |
| Runtime | adapter ID/version, image, command entrypoint profile, security context | No arbitrary privileged command interpolation |
| Network | listener ports, internal-only/exposed, hostname, TLS, egress policy | Exposure must be explicit and policy-approved |
| Configuration | typed fields, defaults, validators, sensitive flags, generated secrets | No secrets embedded in template manifest |
| Data | volumes, persistence class, consistency mode, retention, export format | Stateful upgrades/restore require specialized handlers |
| Health | startup/readiness/liveness probes, timeouts, grace periods | App readiness distinct from container status |
| Operations | supported actions, risk levels, migration hooks, rollback support | Actions map to typed allowlisted agent operations |
| Integration | provided/consumed service contracts, credential scopes, DNS bindings | Cross-service binding is explicit and revocable |
| Metering | CPU/RAM/storage/egress/requests/transform metrics | Tenant-attributable, versioned measurement units |
| Governance | compatible panel/agent versions, OS/arch, phase, approval state | Catalog availability never bypasses release gates |

### 2.2 Illustrative non-executable manifest

```yaml
schema_version: nexuspanel.io/servicetemplate/v1
template_id: media-storage
version: 1.0.0
service_kind: object-storage
workload_category: serving
source:
  image: registry.example.invalid/storage@sha256:<reviewed-digest>
  signature_required: true
runtime:
  adapter: object-storage.v1
  privileged: false
resources:
  floor: {memory_mib: 512, disk_gib: 10}
  typical: {memory_mib: 1024, disk_gib: 30}
  ceiling: {memory_mib: 4096}
capabilities: [docker, persistent-volume, tls]
network:
  exposure: internal-by-default
  tls_required_if_public: true
data:
  persistence: required
  backup_consistency: application-aware
health:
  readiness: adapter-defined
operations: [provision, start, stop, health_check, backup, restore, export, import, uninstall]
```

The digest placeholder above is intentionally invalid. Production manifests MUST contain a real, verified immutable digest. Actual resource floors must be established through measurements and vendor-specific benchmarks, not copied from this example.

### 2.3 Catalog lifecycle

`draft → validated → security_reviewed → published → deprecated → retired`; revoked/compromised is a separate blocking state. Publication requires schema validation, source verification, adapter compatibility, sandbox tests, backup/restore evidence for stateful templates, vulnerability review and approved ownership. Deprecation prevents new installs according to policy; retirement cannot silently delete active instances. Emergency revocation blocks unsafe new installs and invokes an incident workflow for deployed instances.

**SVC-001:** Catalog reads MUST filter visibility by permissions and entitlements without leaking private template metadata.  
**SVC-002:** Publishing MUST be a separately authorized provider action; customer-authored manifests cannot acquire privileged capabilities.  
**SVC-003:** Template version upgrades MUST be explicit, reviewable operations; no unreviewed `latest` tags.  
**SVC-004:** Schema evolution MUST preserve compatibility or supply a migration with rollback and versioned error semantics.  
**SVC-005:** Each published template MUST declare a supported matrix of OS, CPU architecture, panel/agent version, deployment profile and storage backend.  
**SVC-006:** Every catalog entry MUST disclose its license and material external dependencies; do not imply third-party software is owned by NexusPanel.  
**SVC-007:** The catalog MUST distinguish *available*, *eligible but not installed*, *blocked by policy*, *insufficient resources*, *unsupported OS*, and *unverified*.

## 3. Capability detection, placement and resource policy

### 3.1 Capability evaluation pipeline

1. Authenticate requester and establish resource scope without exposing other tenants' inventory.
2. Load host facts (CPU allocation/throttling, RAM, swap, free disk/inodes, filesystem, kernel, cgroup v2, Docker, network, OS/arch, DNS, ports and connectivity).
3. Evaluate template prerequisites and policy allowlists.
4. Calculate *incremental* impact: control-plane overhead, existing workloads, resource ceilings, realistic observed usage, safety headroom and expected data growth.
5. Evaluate data residency, durability, public exposure and backup prerequisites.
6. Return a signed/traceable capability decision with reasons, alternatives and an expiry; recheck immediately before execution.
7. Refuse unsafe placement atomically and do not reserve capacity indefinitely for abandoned plans.

**SVC-008:** Capability floor, typical and ceiling values MUST have a measurement source, version and timestamp.  
**SVC-009:** CPU allocation is a scheduling policy, not a guarantee of dedicated cores; `cpu.weight` and workload class govern contention.  
**SVC-010:** `memory.max` is an enforced ceiling, not reserved physical memory; admission control MUST separately protect host headroom.  
**SVC-011:** Batch workers MUST have queue admission, concurrency limits, cancellation, deadline and contention yielding.  
**SVC-012:** Storage admission MUST account for both bytes and inodes, snapshot/backup overhead and filesystem growth.  
**SVC-013:** Provider placement MUST consider anti-affinity and fault domains where the deployment topology supports them.  
**SVC-014:** Nano MUST expose remote consumption/integration alternatives, not claim local support for container-only services.

### 3.2 Profile expectations

| Profile | Local platform services | Remote consumption | Required user-facing behavior |
|---|---|---|---|
| Nano (512 MB target) | No local Docker-based storage/transform/self-hosted stacks | S3 client, external CDN, remote backups where feasible | Explain unsupported local services and resource delta |
| Starter (1–2 GB) | Selected low-footprint services only after benchmark | All supported external adapters | Show concurrency/storage tradeoffs |
| Standard (4 GB+) | Container services as resources permit | External adapters optional | Enforce headroom and stateful backup preflight |
| Provider/fleet | Placement across eligible nodes | External and local | Show placement, quotas, fault domains and per-tenant metering |

**VERIFY:** Exact RAM floors, Docker overhead, per-service concurrency, image transformation throughput and minimum disk footprints are release-blocking benchmark questions; profile labels alone do not establish compatibility.

## 4. Deployment, reconciliation and operation state machine

A service instance has distinct **desired**, **operation**, and **observed** states. Desired lifecycle: `absent | configured | running | stopped | suspended | deleting`. Observed lifecycle: `unknown | provisioning | starting | ready | degraded | unhealthy | stopped | failed | orphaned`. Operation lifecycle follows Volume 05's durable state machine. An operation success requires the relevant health/consistency verification; Docker returning success is insufficient.

### 4.1 Typed actions and responsibilities

| Operation | Control plane responsibility | Agent/adapter responsibility | Evidence |
|---|---|---|---|
| `service.provision` | Scope, quota, capability, plan, approval, reservation | Create sandbox, volumes, config, network, process | Reconciled instance ID and readiness |
| `service.configure` | Validate version and policy; redact secrets | Atomic config write; safe reload | Observed revision and probe |
| `service.start/stop` | Check dependencies and downtime impact | Start/stop within bounded timeout | Process and application state |
| `service.health_check` | Schedule and surface checks | Probe app-specific readiness | Timestamped check results |
| `service.backup` | Retention, encryption, target, quota | App-consistent snapshot/export | Manifest and verified object |
| `service.restore` | Approval, target, recovery point | Verify and restore safely | Integrity and functional probe |
| `service.export/import` | Ownership, portability, destination validation | Serialize/restore data and config via adapter | Signed manifest and migration report |
| `service.uninstall` | Deletion policy, retention, consent | Stop, revoke, detach, remove managed assets | Deletion receipt and residual inventory |

### 4.2 Reconciliation and failure semantics

- Persist desired revision and operation intent before host mutation; use an idempotency key and expected resource version.
- Recheck policy and host capability at execution time, including destination scope and immutable template digest.
- Validate paths, mounts, DNS, ports, secrets and resource budgets before creating externally visible resources.
- Write checkpoints for irreversible transitions; never claim atomic rollback across databases, DNS, external CDN and object stores.
- A timed-out RPC is **unknown**, not failed; inspect host journal and actual state before retry.
- On crash/restart, reconcile managed resources by immutable instance ID and ownership metadata; do not adopt unmanaged containers or files automatically.
- On repeated failure, apply bounded retry/backoff and a visible repair or restore path; never spin indefinitely under low memory.

**SVC-015:** Every mutation MUST link operation ID, tenant, instance, template version, desired revision, approval and audit event.  
**SVC-016:** Agent operations MUST be typed, bounded, authenticated and idempotent; templates cannot inject arbitrary root commands.  
**SVC-017:** Preflight failures MUST be side-effect-free or produce a documented compensating action.  
**SVC-018:** Unknown execution outcomes MUST reconcile before retry.  
**SVC-019:** The UI MUST distinguish installed, running, ready and publicly reachable.  
**SVC-020:** Uninstall MUST show which data, snapshots, DNS records, credentials and billing entitlements survive.

## 5. Area 28 — Object storage service

### 5.1 Product capabilities

NexusPanel MUST support an S3-compatible storage service through a reviewed adapter (proposed initial engines: MinIO or Garage) on capable nodes, plus **external S3-compatible endpoints** for users whose VPS cannot host storage. The UI exposes buckets, object counts, quota, lifecycle policies, access keys, endpoint/TLS, versioning (where supported), usage, backup/replication status, and integration bindings. Do not imply that S3 API compatibility equals feature parity across providers.

**SVC-021:** Bucket creation MUST validate DNS-safe name, tenant ownership, engine capabilities, placement, quota and encryption policy.  
**SVC-022:** Buckets MUST be private by default. Public read requires explicit impact review and policy permission.  
**SVC-023:** Credentials MUST be scoped to bucket/prefix/actions where supported, displayed only once, rotated/revoked and never logged.  
**SVC-024:** Server-side encryption and TLS MUST be configurable and required according to provider policy; key custody and rotation are governed by Volume 08.  
**SVC-025:** Versioning, retention, object lock and lifecycle MUST be reported as supported/unsupported per backend rather than simulated.  
**SVC-026:** Storage metering MUST distinguish logical stored bytes, physical used bytes, requests, transfer/egress and replication overhead where observable.  
**SVC-027:** Bucket deletion MUST block on nonempty state unless an explicit authorized recursive deletion workflow is invoked.  
**SVC-028:** External storage integrations MUST test endpoint authenticity, credentials, TLS and minimal required permissions without persisting test secrets in logs.

### 5.2 Storage durability and backups

A single-node object store is not inherently durable or highly available. The service card must identify the fault domain and whether copies exist off-node. Snapshotting a live filesystem is not a substitute for backend-consistent backup. An object store used as a **backup destination** cannot be treated as its own sole disaster recovery copy. Recovery must be exercised with data-integrity verification, including multipart uploads, object metadata, version history (where supported), and access policies.

### 5.3 Site-to-bucket integration contract

- Verify site and bucket belong to the same tenant or an explicitly authorized delegated relationship.
- Provision a distinct integration credential; never reuse the bucket owner's master key.
- Store the secret in the managed secret store; inject via a narrowly scoped environment/config binding, with file permissions and reload semantics appropriate to the runtime.
- Provide endpoint, region/path-style options and SDK compatibility guidance.
- Test an authorized write/read/delete in a dedicated test prefix, then clean up.
- Revoke integration independently of bucket deletion; show consumers before rotation.

**SVC-029:** A one-click binding MUST create, verify, inventory and revoke the connection as one auditable operation family.  
**SVC-030:** Credential rotation MUST support a safe overlap window or explicitly document service interruption.

## 6. Area 28 — Image transformation service

### 6.1 Functional contract

An image transformation service (proposed imgproxy or thumbor) consumes original objects and produces resized, cropped, format-converted, optimized derivatives. It is a **batch or bounded on-demand transformation workload**, not a substitute for storage or a CDN. Supported input formats, output formats, maximum pixels, maximum output bytes, animation behavior, metadata stripping, and transform presets must be explicitly declared.

**SVC-031:** Transform requests MUST require signed URLs or an equivalent scoped authorization policy; anonymous arbitrary remote-source fetches are prohibited by default.  
**SVC-032:** The service MUST enforce maximum input bytes, dimensions, output bytes, decode time, CPU time and concurrency.  
**SVC-033:** Remote URL fetching MUST prevent SSRF, metadata endpoint access, internal DNS/IP targets, redirect bypass and DNS rebinding. Prefer managed-bucket sources.  
**SVC-034:** Transformation queues MUST expose queued/running/completed/failed/cancelled and yield under host contention.  
**SVC-035:** Cache keys MUST bind source object version/ETag, transform parameters, tenant, policy and output format; do not leak private derivatives across tenants.  
**SVC-036:** Corrupt/malicious media and decompression bombs MUST fail safely with bounded resource consumption.  
**SVC-037:** Transform metrics MUST include queue wait, execution latency, cache hit/miss, output bytes, failures and resource usage.  
**SVC-038:** Deleting or revoking source content MUST have an explicit derivative invalidation policy.

### 6.2 Pipeline and failure behavior

`authorized request → canonical parameter validation → source authorization → deduplicated cache lookup → bounded queue → isolated decode/transform → validation → cache write → delivery`. A cache miss during saturation may return a documented pending response, a controlled error, or an existing stale derivative under policy; it must not overload the host. Failed transformations must not poison the cache. No transform engine receives host-root credentials.

## 7. Area 28 — Asset delivery, caching and external CDN

### 7.1 Local delivery

The local delivery adapter (proposed Caddy/Nginx caching) provides TLS, hostname binding, cache policy, origin access controls, content-type correctness, range requests where supported, and logging. **It is not globally distributed delivery.** Optional external CDN adapters (Cloudflare or Bunny proposed) add separate DNS, caching, purge, edge TLS, and provider-specific billing/quotas.

**SVC-039:** Public hostname activation MUST validate domain ownership/authorization and TLS readiness before advertising ready state.  
**SVC-040:** Cache policy MUST distinguish immutable public assets, mutable public assets and private/authenticated objects. Private objects MUST NOT be cached in a shared public cache without a verified isolation scheme.  
**SVC-041:** Origin protection MUST support signed URLs or scoped origin credentials and prevent bypass of intended access controls.  
**SVC-042:** Cache invalidation MUST support targeted purge with operation status, provider-specific eventual consistency and visible failure.  
**SVC-043:** External CDN configuration MUST clearly disclose that it is a third-party service and may incur separate fees.  
**SVC-044:** Delivery analytics MUST distinguish origin traffic, edge egress when available, cache hit ratio, errors and estimated costs.  
**SVC-045:** DNS changes and CDN activation MUST use staged verification and rollback guidance; do not assume instant global propagation.

### 7.2 Security and reliability

- Enforce HTTPS and correct host routing; HSTS changes require domain-level safety review.
- Never use customer-controlled URL fragments to construct unrestricted filesystem paths.
- Protect origin from request floods, abusive transforms, cache-key explosion and hotlinking according to policy.
- Expose origin-down and CDN-down states separately; a healthy edge cache does not imply origin recovery.
- On external CDN API outage, preserve the last known configuration, show degraded state, and retry only idempotent operations.

## 8. Area 28 — Cache, queue and search services

These are part of the target 29-area scope, but the source v2.0 release mapping places much of their catalog expansion in v2+. They share the template catalog, not necessarily identical runtime or data semantics.

### 8.1 Redis/Memcached-style caching

**SVC-046:** Cache instances MUST be private-network by default, tenant-scoped, password/ACL protected where supported and bounded by memory/eviction policy.  
**SVC-047:** The UI MUST disclose whether data is ephemeral or persisted; cache eviction is not an integrity failure unless the application contract says otherwise.  
**SVC-048:** Cache credentials and namespace bindings MUST be revocable independently per consumer.  
**SVC-049:** Cache metrics MUST include hit/miss, evictions, memory, connections and latency where supported.

### 8.2 Queue services

**SVC-050:** Queue services MUST expose delivery semantics (at-most-once, at-least-once, ordering limits), retries, dead-letter handling, visibility timeout and retention.  
**SVC-051:** Queue consumers MUST be idempotent or implement documented deduplication.  
**SVC-052:** Queue instances MUST have resource caps and overload behavior that protect serving workloads.  
**SVC-053:** Queue depth, oldest message age, redelivery and poison-message rate MUST be visible to authorized users.

### 8.3 Search services

**SVC-054:** Search indexing MUST separate tenant indexes, access policies and encryption boundaries.  
**SVC-055:** Search services MUST declare indexing memory/disk overhead, rebuild time and recovery semantics.  
**SVC-056:** Reindex jobs MUST be batch-governed and not silently consume the host's serving headroom.  
**SVC-057:** Search APIs MUST enforce query cost/rate limits and prevent unauthorized cross-tenant document retrieval.

## 9. Area 29 — Self-hosted application catalog

### 9.1 Common installation contract

Each app template MUST specify its image provenance, configuration schema, ports, domain/TLS needs, persistent volumes, database/cache dependencies, admin bootstrap procedure, email/SMTP dependencies, upgrade strategy, backup/restore procedure, security hardening, readiness checks, and minimum resource envelope. One-click means **guided and automated**, not “install with no prerequisites or risk.”

**TPL-001:** The catalog MUST expose vendor/project name, license, maintained versions, known prerequisites, support level and last validation date.  
**TPL-002:** Admin credentials MUST be generated or supplied through secure write-only flows; default vendor credentials are prohibited.  
**TPL-003:** Dependencies MUST be explicit and deployable in order, with an auditable dependency graph and cleanup rules.  
**TPL-004:** App templates MUST pin images by digest and declare supported upgrade paths; no silent major upgrades.  
**TPL-005:** Persistent volumes MUST be separated from immutable image layers and included in backup/restore plans.  
**TPL-006:** Public exposure MUST be explicit, TLS-protected and protected by application-specific authentication.  
**TPL-007:** App readiness MUST include an application-level check and, for stateful apps, a data-integrity or login/bootstrap probe.  
**TPL-008:** Uninstall MUST offer retain/export/delete options for user data and clearly state what remains.  
**TPL-009:** Templates MUST be tested on each claimed OS/arch/profile combination; “runs in Docker” alone is insufficient.  
**TPL-010:** Each template MUST document unsupported operations and known data-loss risks.

### 9.2 v1 reference application profiles

| Template | Key dependencies | Special risks and required checks | Suggested release evidence |
|---|---|---|---|
| WordPress | PHP/web stack, MariaDB/MySQL, SMTP optional | File ownership, plugin/theme security, DB migration, uploads | Fresh install, admin login, media upload, backup/restore, safe update |
| Nextcloud | Web/PHP stack, DB, persistent storage, background jobs, cache recommended | File locking, data directory isolation, cron, app-version migrations | Upload/share/download, user login, consistent backup and restore |
| Gitea | Persistent repos, DB (profile dependent), SSH/HTTP endpoints | Git SSH key custody, repo permissions, webhook SSRF, backup consistency | Push/pull, auth, repo restore, upgrade |
| Ghost | Node runtime, DB as required by supported release, SMTP | Version compatibility, DB migrations, admin bootstrap | Publish post, media, email config, restore |
| n8n | DB/persistent state, webhook ingress, encryption key | Workflow secrets, arbitrary outbound requests, public webhook abuse | Run workflow, restart persistence, secret preservation, restore |
| Uptime Kuma | Persistent data, monitoring egress | SSRF/network scans, monitor credential protection | Create monitor, alert, restart persistence, restore |

The exact dependencies and minimum resource floors are **VERIFY** items against the specific vendor releases selected at implementation time. Some templates may be unavailable on Starter despite appearing in the v1 target catalog.

### 9.3 Application-specific acceptance obligations

**TPL-011 (WordPress):** A new instance MUST finish installation with a reachable HTTPS site and an authenticated admin account; backup and restore MUST preserve content, uploads and DB state.  
**TPL-012 (Nextcloud):** A restored instance MUST preserve files, metadata, users and app-specific configuration; locking and scheduled jobs MUST be healthy.  
**TPL-013 (Gitea):** Git-over-SSH/HTTPS access MUST respect repo ownership; export/import MUST preserve repositories, users and permissions as documented.  
**TPL-014 (Ghost):** Content and media MUST survive a supported major-version migration with explicit rollback/recovery instructions.  
**TPL-015 (n8n):** Encryption keys and stored workflow credentials MUST survive backup/restore without exposure in exported plaintext.  
**TPL-016 (Uptime Kuma):** Monitors MUST not permit unauthorized probing of protected host/network endpoints.

### 9.4 Catalog expansion

Later templates may include Meilisearch, Plausible, MinIO, imgproxy, PostgreSQL, MariaDB and Redis, with the source roadmap targeting roughly 20–30 templates in v2+. Each addition must pass the same manifest, security, operational and acceptance gates; it must not be treated as a simple unreviewed database insert.

## 10. Secrets, identities and service-to-service bindings

- Every service instance receives an immutable instance ID and an explicit tenant, account, project, host and policy scope.
- A service credential has a principal, scope, creation method, expiry/rotation policy, last-used metadata where supported, revocation state and audit trail.
- Credential values are encrypted at rest, never returned after creation unless an explicit approved recovery mechanism exists; generated-once display is preferred.
- Cross-service bindings use dedicated scoped credentials and are tracked as first-class resources; deleting a consumer revokes its binding without deleting shared infrastructure.
- Template-defined configuration is validated against a typed schema; sensitive values cannot appear in shell arguments, process listings, event payloads, metrics labels or support bundles.
- Secret injection must avoid exposing values through world-readable environment files or container inspect permissions; deployment method must be threat-modeled.
- External connectors use least-privileged API tokens and verify callback authenticity according to the provider's documented mechanism.

**SVC-058:** Cross-tenant binding is denied unless explicit delegated authorization covers both sides.  
**SVC-059:** Revocation MUST propagate to the consuming service or report a degraded/manual-intervention state.  
**SVC-060:** Rotation MUST be safe under retries and agent disconnects; stale credentials MUST be identified.  
**SVC-061:** Secret-bearing exports MUST be encrypted and must not be exposed to untrusted catalog publishers.

## 11. Container and supply-chain security

### 11.1 Runtime restrictions

- Default rootless/non-root process, read-only root filesystem where supported, dropped Linux capabilities, seccomp/AppArmor/SELinux policy, bounded mounts, no host PID/IPC/network namespaces unless specifically approved and isolated.
- No privileged containers, host Docker socket mounts, arbitrary host path mounts, or unreviewed device access in the general template catalog.
- Separate internal service networks and explicit public ingress; restrict outbound connections for services with SSRF or data-exfiltration risk.
- Validate image signatures, digests, SBOM and known critical vulnerabilities at publication and deployment time; define exceptions with expiry and owner.
- Isolate template build hooks from the privileged agent and control-plane credentials; never run catalog-supplied scripts as host root.
- Protect against malicious publishers, compromised registries, dependency confusion, typosquatting and rollback to revoked images.

**SVC-062:** Container execution MUST use the approved security profile for the service class.  
**SVC-063:** An untrusted template MUST NOT obtain the host's root privileges, panel database, agent socket or cloud credentials.  
**SVC-064:** Images MUST be pinned by immutable digest and checked against the published manifest.  
**SVC-065:** Vulnerability policy MUST define blocking thresholds, temporary exceptions and re-evaluation cadence.  
**SVC-066:** Template execution MUST enforce egress, ingress and volume scopes at the actual runtime boundary, not only in the UI.

### 11.2 Template trust levels

`official_reviewed`, `provider_approved`, `community_unreviewed`, `private_custom` are **catalog trust labels**, not permission grants. General availability of unreviewed templates is deferred until a defensible sandbox and review process exists. The v1 catalog should use official/provider-approved templates only. Custom template imports are subject to an administrator approval and risk warning; public customer installation of arbitrary privileged templates is out of scope for v1.

## 12. Persistent data, backups, export/import and migration

### 12.1 Data classification and consistency

Each instance declares `stateless`, `ephemeral_cache`, `stateful_files`, `stateful_database`, or `multi_component_stateful`. Backup and restore methods vary accordingly. For multi-component apps, application-consistent recovery may require quiescing writers, coordinating database and file snapshots, or using vendor-native exports. A container image alone is never a backup of persistent data.

**SVC-067:** Every stateful template MUST publish a tested backup/restore plan before production release.  
**SVC-068:** Backups MUST include necessary application config, volume data, database state, secret references and template/image versions without exposing plaintext credentials.  
**SVC-069:** Export MUST include a signed manifest of schema, origin, timestamps, checksums, data formats, ownership and compatibility constraints.  
**SVC-070:** Import MUST validate the manifest, destination capability, disk headroom, policy and integrity before changing active routing.  
**SVC-071:** Restores MUST run in an isolated target or explicit destructive replacement workflow, with impact preview and approvals.  
**SVC-072:** Data migration MUST verify application functionality, not only archive extraction or container readiness.  
**SVC-073:** Exported credentials MUST be reissued or safely rewrapped for destination identity; source credentials must not remain indefinitely valid.  
**SVC-074:** Offsite backups and periodic restoration drills MUST be supported, with evidence retained under policy.

### 12.2 Portability boundaries

OCI artifacts are appropriate for **images and manifests**; application state is separately exported using service-aware formats. A template can be portable while its data backend is not. Export/import must declare unsupported migration paths (for example, engine-specific storage features or database major-version incompatibility) and block unsupported transitions with actionable explanations.

### 12.3 Recovery objective declarations

Each template may declare target RPO/RTO classes for supported topology and backup frequency, but no RPO/RTO is guaranteed until measured restoration tests support it. Nano remote-service integrations depend on external-provider recovery characteristics; NexusPanel must not claim control of their durability.

## 13. Network, DNS, TLS and public exposure

1. Service templates declare internal ports and optional public endpoints. Host port allocation is centrally arbitrated to avoid collisions and unauthorized binding.
2. Public services require domain authorization, DNS validation, certificate readiness and approved ingress policies; private services must remain private by default.
3. Reverse-proxy routes bind immutable service IDs, tenant scopes, desired revisions and health states; stale routes are reconciled after service deletion.
4. Webhooks and externally accessible admin interfaces require explicit risk disclosure and appropriate authentication/rate limits.
5. Outbound requests from automation tools, monitoring and media transforms are SSRF-sensitive and must respect egress policy and DNS rebinding defenses.
6. Provider-controlled network policy may prohibit SMTP, torrenting, open proxies, public databases or arbitrary ports; the installer must surface such restrictions.
7. External DNS/CDN providers may be eventually consistent. Operation completion is based on verified readiness criteria, not just accepted API calls.

**SVC-075:** Public exposure MUST be denied unless the caller has both service and domain/network permissions.  
**SVC-076:** TLS termination and origin TLS policies MUST be explicit for each public route.  
**SVC-077:** Route removal MUST be included in uninstall, rollback and compromised-service quarantine.  
**SVC-078:** Public service health MUST be verified from an appropriate network perspective, not only localhost.

## 14. Metering, quotas, pricing and commercial entitlements

A template may declare units such as stored GiB-month, object operations, outbound GiB, transform count/CPU time, provisioned RAM ceiling, active instances, backup storage and requests. These are **metering dimensions**, not automatic billable items. Billing policy, subscriptions, entitlements, taxation and gateway settlement are governed by Volume 09 and Volume 05.

**SVC-079:** Metering MUST preserve tenant attribution and distinguish estimated from measured usage.  
**SVC-080:** Quota checks MUST occur at plan time and execution time; concurrent provisioning uses atomic reservations or equivalent reconciliation.  
**SVC-081:** Usage reporting MUST identify collection lag, sampling and backend limitations.  
**SVC-082:** Paid activation MUST depend on confirmed entitlement state, never on an unverified payment callback or client-supplied success flag.  
**SVC-083:** Service suspension MUST define data retention, ingress, background jobs, backups and reactivation behavior.  
**SVC-084:** Provider catalog policy MUST be able to limit template visibility, versions, instance counts and resource profiles per plan without bypassing tenant authorization.

A Nano user consuming a third-party storage/CDN service may incur external charges. The UI must disclose which charges are NexusPanel-managed and which belong to external vendors.

## 15. Observability, troubleshooting and support

Per instance, the UI exposes desired/observed state, last reconciliation, image digest, template version, deployment history, readiness checks, bounded logs, resource usage, backup freshness, certificate state, active bindings, incidents, and recent audit events. Batch workloads expose queue wait and progress. Support bundles redact secrets and tenant data by default.

**SVC-085:** All critical operations MUST emit structured audit and diagnostic events with correlation and operation IDs.  
**SVC-086:** Service logs MUST be tenant-scoped, bounded by retention/volume quotas, and redacted before export.  
**SVC-087:** Health checks MUST distinguish startup delay, dependency outage, application failure, ingress failure and unknown/offline agent.  
**SVC-088:** Alerts MUST be deduplicated and actionable, with safe retry/recovery guidance and no false success claims.  
**SVC-089:** Troubleshooting recommendations generated by AI MUST be read-only by default and governed by the normal approval path for any proposed action.

### 15.1 Suggested operational metrics

| Layer | Metrics | Notes |
|---|---|---|
| Template registry | publish/revoke events, validation failures, version adoption | Supply-chain oversight |
| Provisioning | plan latency, queue wait, retries, completion, rollback | Operation ID trace |
| Runtime | CPU weight/usage, memory current/max/OOM, disk/inodes, restart count | Measured vs allocated |
| Storage | bytes, objects, PUT/GET failures, egress, durability state | Logical vs physical bytes |
| Transform | queue age, decode failures, CPU seconds, output bytes, cache hits | Prevent resource abuse |
| Delivery | origin errors, edge status, cache hit ratio, egress | Separate local from CDN |
| Stateful apps | backup age, restore test age, DB/volume integrity, health | Production release evidence |

## 16. API, data and UI cross-volume bindings

### 16.1 API surfaces (Volume 05)

- `GET /api/v1/service-templates` — catalog filtered by visibility and compatibility; includes capability reasons.
- `GET /api/v1/service-templates/{id}/versions/{version}` — reviewed manifest, provenance and compatibility.
- `POST /api/v1/service-instances/plan` — validate desired deployment, produce scoped plan, resource estimate, risks and expiration.
- `POST /api/v1/service-instances` — create instance via idempotent durable operation; `202` until verified.
- `GET /api/v1/service-instances/{id}` — desired, observed, readiness, bindings, usage and operation links.
- `POST /api/v1/service-instances/{id}/operations` — typed lifecycle operations with authorization and approvals.
- `POST /api/v1/service-instances/{id}/bindings` — scoped consumer integration.
- `GET /api/v1/service-instances/{id}/usage` — measured and allocated values with timestamps.
- `POST /api/v1/service-instances/{id}/exports` and `/imports` — specialized data-aware workflows.

These are **conceptual routes** requiring alignment with Volume 05's final OpenAPI contract; this volume does not unilaterally replace previously specified endpoint names. Every mutation carries idempotency, expected revision, tenant scope, operation link, audit and consistent error semantics.

### 16.2 Data model (Volume 04)

Required logical entities: `service_templates`, `service_template_versions`, `service_instances`, `service_credentials`, `service_bindings`, `service_volumes`, `service_health_checks`, `service_metrics`, `service_exports`, `service_backups`, `service_deployments`, `service_dependencies`, `capabilities`, `operations`, `audit_events`, and applicable tenant/entitlement/resource entities. Where Volume 04 uses different physical names, its schema is authoritative; names here denote conceptual relationships. Cross-tenant foreign keys and credential disclosure are prohibited.

### 16.3 UI states (Volume 02)

Catalog cards and instance screens must expose: capability available/unavailable with reason, trust/review status, estimated resources, actual usage, version, health, backup readiness, cost ownership, and safe next action. The router recommends by **job-to-be-done**, not by unrestricted role-specific menus. No disabled item without explanation. No success screen until the relevant verification gate passes.

## 17. Release phases and delivery dependencies

| Phase | Platform Services | Self-hosted Templates | Required foundations |
|---|---|---|---|
| v1 target | Storage + transform + local delivery, optional external CDN integration, remote service consumption | WordPress, Nextcloud, Gitea, Ghost, n8n, Uptime Kuma as proposed catalog targets | Go control plane/agent, typed ops, Docker on eligible hosts, secrets, TLS, backups, audit, capability gate |
| v2+ | Cache/queue/search catalog expansion; improved provider adapters, multi-node options | Broader 20–30 template catalog, richer integrations | Plugin governance, specialized adapters, verified migration and provider policy |
| v3+ | Fleet placement, replication/fault-domain policy, commercial scaling | Provider-managed template ecosystem | Fleet controls, HA and cross-node recovery evidence |

**Reconciliation note:** Source documents describe some platform services as v1 targets while other area-level matrices mark related infrastructure as v2. The release contract is therefore: *the named v1 reference workflows must be proven on eligible Standard+ configurations*, while general-purpose self-service orchestration and broader catalogs may ship later. Any template whose mandatory dependency (such as a database lifecycle adapter) is not release-ready cannot be advertised as v1 production-ready. Resolve exact release packaging during Volume 10 approval.

### 17.1 Dependency graph

`Identity/tenancy → policy/entitlements → host capability → template registry/signature → durable operations → secure runtime/network/volumes → secrets → health → backup/restore → integration bindings → metering → public catalog publication`.

A module may be developed earlier for testing, but must not bypass its security and recovery dependencies at release.

## 18. Acceptance tests and release gates

The following scenarios are mandatory for the affected feature's production readiness. Test IDs are stable within this volume.

**TEST-SVC-CAT-01:** Reject unsigned or digest-mismatched template; no host mutation. **Evidence:** automated integration/E2E result, operation/audit trace, and reproducible environment record; manual security review where applicable.

**TEST-SVC-CAT-02:** Publish approved template without code release when only existing allowlisted adapters are used. **Evidence:** automated integration/E2E result, operation/audit trace, and reproducible environment record; manual security review where applicable.

**TEST-SVC-CAT-03:** Prevent unreviewed catalog author from requesting privileged mounts or host root. **Evidence:** automated integration/E2E result, operation/audit trace, and reproducible environment record; manual security review where applicable.

**TEST-SVC-CAT-04:** Show accurate unavailable capability on 512 MB Nano; reject local Docker deployment. **Evidence:** automated integration/E2E result, operation/audit trace, and reproducible environment record; manual security review where applicable.

**TEST-SVC-CAT-05:** Install supported template on eligible Standard host and verify app-level readiness. **Evidence:** automated integration/E2E result, operation/audit trace, and reproducible environment record; manual security review where applicable.

**TEST-SVC-CAT-06:** Prevent oversubscription and protect SSH/control-plane headroom under concurrent installs. **Evidence:** automated integration/E2E result, operation/audit trace, and reproducible environment record; manual security review where applicable.

**TEST-SVC-CAT-07:** Preserve operation state after control-plane restart and agent disconnect. **Evidence:** automated integration/E2E result, operation/audit trace, and reproducible environment record; manual security review where applicable.

**TEST-SVC-CAT-08:** Resolve unknown outcome without duplicate volumes, users or public routes. **Evidence:** automated integration/E2E result, operation/audit trace, and reproducible environment record; manual security review where applicable.

**TEST-SVC-CAT-09:** Deny cross-tenant service/bucket binding and conceal private inventory. **Evidence:** automated integration/E2E result, operation/audit trace, and reproducible environment record; manual security review where applicable.

**TEST-SVC-CAT-10:** Create S3 bucket private by default; scoped credentials work and revoke. **Evidence:** automated integration/E2E result, operation/audit trace, and reproducible environment record; manual security review where applicable.

**TEST-SVC-CAT-11:** Connect website to bucket, rotate credential and verify no secret leakage. **Evidence:** automated integration/E2E result, operation/audit trace, and reproducible environment record; manual security review where applicable.

**TEST-SVC-CAT-12:** Restore bucket/object data and metadata with checksum evidence. **Evidence:** automated integration/E2E result, operation/audit trace, and reproducible environment record; manual security review where applicable.

**TEST-SVC-CAT-13:** Reject transform SSRF including redirect, DNS rebinding and private IP. **Evidence:** automated integration/E2E result, operation/audit trace, and reproducible environment record; manual security review where applicable.

**TEST-SVC-CAT-14:** Bound decompression bomb and oversized media; host remains responsive. **Evidence:** automated integration/E2E result, operation/audit trace, and reproducible environment record; manual security review where applicable.

**TEST-SVC-CAT-15:** Yield transform batch under serving load and expose queue state. **Evidence:** automated integration/E2E result, operation/audit trace, and reproducible environment record; manual security review where applicable.

**TEST-SVC-CAT-16:** Invalidate transformed derivative after source revocation. **Evidence:** automated integration/E2E result, operation/audit trace, and reproducible environment record; manual security review where applicable.

**TEST-SVC-CAT-17:** Public delivery validates TLS/domain and never exposes private cached object. **Evidence:** automated integration/E2E result, operation/audit trace, and reproducible environment record; manual security review where applicable.

**TEST-SVC-CAT-18:** External CDN purge failure remains visible and retry-safe. **Evidence:** automated integration/E2E result, operation/audit trace, and reproducible environment record; manual security review where applicable.

**TEST-SVC-CAT-19:** WordPress backup/restore preserves posts, database and uploads. **Evidence:** automated integration/E2E result, operation/audit trace, and reproducible environment record; manual security review where applicable.

**TEST-SVC-CAT-20:** Nextcloud backup/restore preserves file and user integrity. **Evidence:** automated integration/E2E result, operation/audit trace, and reproducible environment record; manual security review where applicable.

**TEST-SVC-CAT-21:** Gitea restore preserves repository data and authorized Git access. **Evidence:** automated integration/E2E result, operation/audit trace, and reproducible environment record; manual security review where applicable.

**TEST-SVC-CAT-22:** n8n restore preserves workflow encryption without plaintext export. **Evidence:** automated integration/E2E result, operation/audit trace, and reproducible environment record; manual security review where applicable.

**TEST-SVC-CAT-23:** Uptime Kuma prevents unauthorized internal network probing. **Evidence:** automated integration/E2E result, operation/audit trace, and reproducible environment record; manual security review where applicable.

**TEST-SVC-CAT-24:** Stateful major upgrade has compatible rollback/restore evidence. **Evidence:** automated integration/E2E result, operation/audit trace, and reproducible environment record; manual security review where applicable.

**TEST-SVC-CAT-25:** Export/import across VPS verifies manifest, ownership, integrity and cutover. **Evidence:** automated integration/E2E result, operation/audit trace, and reproducible environment record; manual security review where applicable.

**TEST-SVC-CAT-26:** Uninstall retains or deletes data exactly as approved and revokes credentials. **Evidence:** automated integration/E2E result, operation/audit trace, and reproducible environment record; manual security review where applicable.

**TEST-SVC-CAT-27:** Quota enforcement and billing entitlement survive concurrent callback and provisioning. **Evidence:** automated integration/E2E result, operation/audit trace, and reproducible environment record; manual security review where applicable.

**TEST-SVC-CAT-28:** Logs, support bundle and metrics expose no secrets or cross-tenant content. **Evidence:** automated integration/E2E result, operation/audit trace, and reproducible environment record; manual security review where applicable.

**TEST-SVC-CAT-29:** Resource metrics show measured usage vs ceiling with sample freshness. **Evidence:** automated integration/E2E result, operation/audit trace, and reproducible environment record; manual security review where applicable.

**TEST-SVC-CAT-30:** Service unhealthy state offers actionable recovery and never reports success. **Evidence:** automated integration/E2E result, operation/audit trace, and reproducible environment record; manual security review where applicable.

**TEST-SVC-CAT-31:** Template revocation blocks new deployment and flags existing instances. **Evidence:** automated integration/E2E result, operation/audit trace, and reproducible environment record; manual security review where applicable.

**TEST-SVC-CAT-32:** Cache service enforces private network, memory limit and auth. **Evidence:** automated integration/E2E result, operation/audit trace, and reproducible environment record; manual security review where applicable.

**TEST-SVC-CAT-33:** Queue redelivery follows documented semantics and dead-letter policy. **Evidence:** automated integration/E2E result, operation/audit trace, and reproducible environment record; manual security review where applicable.

**TEST-SVC-CAT-34:** Search indexes remain isolated across tenants and rate-limited. **Evidence:** automated integration/E2E result, operation/audit trace, and reproducible environment record; manual security review where applicable.

**TEST-SVC-CAT-35:** External storage/CDN outage does not corrupt control-plane desired state. **Evidence:** automated integration/E2E result, operation/audit trace, and reproducible environment record; manual security review where applicable.

**TEST-SVC-CAT-36:** No unsupported vendor version is claimed compatible without test evidence. **Evidence:** automated integration/E2E result, operation/audit trace, and reproducible environment record; manual security review where applicable.

## 19. Failure and abuse scenario matrix

| Scenario | Required response | Recovery/verification |
|---|---|---|
| Host OOM or disk full | Reject new workload, preserve control plane/SSH, report resource shortage | Release capacity, reconcile volumes and restart safely |
| Compromised template image | Revoke catalog version, quarantine affected instances under policy | Signed replacement, incident evidence, secret rotation |
| Broken app migration | Halt activation, preserve backup, mark degraded | Service-specific rollback or restore and integrity check |
| Object-store corruption | Stop destructive reconciliation, mark durability incident | Restore from independent verified copy |
| External CDN provider outage | Preserve origin config and desired state, show partial outage | Retry signed idempotent provider calls, verify propagation |
| Transformation abuse | Enforce signed URLs, rate limits, CPU/time/memory bounds | Quarantine abusive client and purge affected queue |
| Stolen integration key | Revoke scoped key and all active bindings | Issue replacement, review audit and object access |
| Tenant deletes shared service | Deny unless all consumers are handled and scope authorized | Dependency-aware removal plan |
| Agent unreachable mid-install | Mark outcome unknown, never report success | Journal inspection and reconciliation |
| Payment pending/reversed | Do not grant unconfirmed entitlement | Reconcile with authoritative gateway state |
| DNS/TLS not ready | Do not expose service as publicly healthy | Recheck ownership, cert and ingress |
| Registry unreachable | No unverified mutable-tag fallback | Retry or use verified local digest if policy permits |

## 20. Requirements traceability and 29-area coverage

| Area | Relevance to this volume | Owning detailed volume |
|---|---|---|
| 01 Dashboard | Portfolio, usage, readiness and alerts | 02, 07 |
| 02 Server management | Capability, ports, storage, resource limits | 03, 06 |
| 03 Websites | Object storage/media bindings, TLS, routes | 06, 07 |
| 04 DNS | Public hostname verification | 06, 08 |
| 05 SSL/TLS | Secure endpoints and CDN origin | 06, 08 |
| 06 File manager | Volumes, user data, permissions | 06, 07 |
| 07 Application manager | App deployment, lifecycle, logs | 06, 07 |
| 08 Runtime manager | Compatibility and runtime dependencies | 06, 07 |
| 09 Databases | Stateful app dependencies and recovery | 04, 06, 07 |
| 10 Email hosting | SMTP dependencies and notification flows | 06, 07 |
| 11 Git/deployments | Image provenance and app deployment | 05, 06, 07 |
| 12 Backups | App-consistent recovery and export | 06, 07, 08 |
| 13 Security | Container sandbox, secrets, SSRF | 07, 08 |
| 14 Monitoring | Per-instance metrics, alerts | 02, 07 |
| 15 Terminal/developer | Scoped diagnostics only | 06, 08 |
| 16 Users/accounts | Ownership and service delegation | 01, 07, 08 |
| 17 AI assistant | Guided diagnostics without policy bypass | 02, 07, 08 |
| 18 Automation/API | Service lifecycle and integrations | 05, 07 |
| 19 Panel administration | Registry approvals, compatibility | 03, 07 |
| 20 Extensions | Trusted templates vs plugin runtime | 07, 08 |
| 21 Commerce/reseller | Catalog policy, entitlements, metering | 07, 09 |
| 22 Migration/staging | Service export/import and cutover | 06, 07 |
| 23 Containers/orchestration | Secure runtime and placement | 03, 06, 07 |
| 24 Performance/caching | Local cache and capacity controls | 07 |
| 25 Email security | External SMTP and webhook abuse | 08, 09 |
| 26 Analytics/site quality | Traffic/asset delivery metrics | 02, 07 |
| 27 Audit/compliance/support | Tamper evidence and incident handling | 07, 08 |
| **28 Platform Services** | **Storage, transforms, delivery, cache, queue, search** | **07 (primary)** |
| **29 Self-hosted Templates** | **Catalog, WordPress, Nextcloud, Gitea, Ghost, n8n, Uptime Kuma** | **07 (primary)** |

## 21. Open approvals and validation backlog

These are not retroactively treated as user-approved decisions. The baseline architecture is integrated above; the following still require a named owner and evidence before a production release.

| ID | Decision / validation | Proposed default | Required evidence |
|---|---|---|---|
| DEC-SVC-01 | Initial object storage engine | MinIO or Garage adapter, choose by benchmark | Maintenance, licensing, memory, durability and S3 compatibility |
| DEC-SVC-02 | Image transformer | imgproxy preferred | Security, memory, format coverage and load test |
| DEC-SVC-03 | External CDN adapter | Cloudflare first, Bunny optional | API, billing, token scopes, purge tests |
| DEC-SVC-04 | v1 template lineup | Six reference apps | Dependency readiness and per-app restore evidence |
| DEC-SVC-05 | Starter service eligibility | No blanket guarantee | 1 GB/2 GB load and contention benchmarks |
| DEC-SVC-06 | Template publishing policy | Official/provider-approved only in v1 | Signature and sandbox review |
| DEC-SVC-07 | Template data migration format | Signed manifest + service-aware encrypted data | Cross-host restore drill |
| DEC-SVC-08 | External service cost visibility | Distinguish third-party fees | Product and legal review |
| DEC-SVC-09 | Storage replication policy | Single-node labeled non-HA | Failure-domain and durability testing |
| DEC-SVC-10 | Service health contract | Readiness + functional probe | Service-specific tests |
| DEC-SVC-11 | Secrets injection mechanism | Scoped encrypted store + runtime adapter | Threat model and leak tests |
| DEC-SVC-12 | Queue/search/cache delivery phase | Broader v2+ | Roadmap dependency signoff |

## 22. Definition of done for Volume 07

Volume 07 is **documentation complete** when its requirements are reconciled with Volumes 01–06, the two new platform areas are fully mapped, all interfaces and error states are specified, and open decisions are visible. It is **implementation ready** only after the manifest schema, agent adapter contracts, threat model, OpenAPI, physical data schema and resource benchmarks are reviewed and approved. It is **production ready** only when every claimed reference template and platform service passes its acceptance tests, including real backup/restore and security tests, on supported hardware. None of these implementation or production claims is made by this document.

**Next volume:** 08 — Security, Compliance & Disaster Recovery. It must further specify the template trust model, tenant boundaries, key custody, signed updates, incident response and backup integrity.

