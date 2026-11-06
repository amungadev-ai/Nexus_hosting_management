# NexusPanel — Master Product & Engineering Specification
## Volume 08 — Security, Compliance & Disaster Recovery

**Version:** 1.0 consolidated | **Status:** Proposed engineering specification, approval and security validation pending | **Date:** 2026-11-06  
**Authority:** Volume 01 Master PRD and authoritative 29-area scope; engineering dependencies in Volumes 02–07.  
**Source basis:** User-provided Stage 1–7 v2.0 documents, especially security, backups, privileged agent, data, and API sections.  
**Operational boundary:** Documentation only. No commands executed on the Leaseweb VPS.

## 0. Purpose, precedence and scope

This volume defines the security and recovery contract for NexusPanel across all **29 platform areas**, every tenant and role, and every capability profile (Nano, Starter, Standard, Provider). It is a product and engineering specification, **not** an assertion of deployed controls, certification, tested performance, or legal compliance. The v2.0 source requirements are preserved below with their original `NP-*`, `ARCH-*`, and `AGT-*` identifiers; new detail uses `SEC8-*` identifiers. Source contract sections are included to retain the full underlying traceability, while the integrated clauses at the start resolve their implementation boundaries.

**Precedence:** Volume 01 controls product commitments. Volume 03 controls architecture, Volume 04 data contracts, Volume 05 APIs, Volume 06 privileged host behavior and Volume 07 platform-service lifecycles. This volume adds mandatory security/recovery gates. If a source gives an exact target without measurement, it remains a **proposed target**. Contradictory source statements must be decided at the release gate rather than silently interpreted as implementation fact.

**Terms:** `MUST` is mandatory for a claimed production-ready module; `SHOULD` allows a documented exception; `MAY` is optional; `VERIFY` demands test evidence; `APPROVAL` demands an explicit product/technical decision. A local backup is not offsite disaster recovery; a successful backup process is not a verified restore; a signed image is not a proof of safe runtime behavior.

## 1. Security objectives, boundaries and invariants

| ID | Mandatory invariant | Enforcement boundary | Evidence |
|---|---|---|---|
| SEC8-001 | No user, tenant, reseller or support principal gains authority merely by changing a frontend persona or route | Server-side policy on every request | Role × tenant × object negative tests |
| SEC8-002 | No raw shell, arbitrary command string, arbitrary host path or unchecked container configuration crosses into the privileged agent | Typed agent allowlist, schema, path ownership, runtime adapter | Agent fuzzing, injection tests |
| SEC8-003 | Every high-risk mutation binds actor, target, immutable plan digest, approval and idempotency to one durable operation | API + operation ledger + agent | Replay, stale approval, plan-substitution tests |
| SEC8-004 | Secret plaintext is not logged, sent to analytics or included in AI context by default | Secret vault, redaction, data egress policy | Canary leak tests |
| SEC8-005 | Tenant suspension, role revocation and compromised-credential revocation invalidate pending authorization, including queued work | Policy check at enqueue and immediately before privileged execution | Revocation race test |
| SEC8-006 | Service health is never inferred solely from a successful command exit or stale heartbeat | Verified observations with timestamps and confidence | Offline and failed-probe simulations |
| SEC8-007 | Backups are encrypted, independent of the source failure domain, integrity checked and periodically restored to disposable targets | Backup engine, storage and restore verification | Destructive recovery drill |
| SEC8-008 | Audit is tamper-evident, retention-controlled, access-restricted and exported outside tenant-writable storage | Hash chain, append-only policy, external sink | Tamper detection, chain reconstruction |
| SEC8-009 | Firewall/SSH changes cannot remove the last verified management path without a timed revert and console recovery | Host agent safe-apply workflow | Lockout simulation |
| SEC8-010 | AI cannot autonomously perform privileged execution in v1, even when a prompt claims urgency or admin authority | AI tool gateway + policy + human approval | Prompt-injection red team |
| SEC8-011 | Nano resource savings never remove core security boundaries or imply unsupported recovery guarantees | Capability gating + admission preflight | 512 MB target benchmark and failure tests |
| SEC8-012 | Production launch requires documented threat model, dependency provenance, penetration test, incident runbooks and recovery evidence | Release governance | Signed release gate checklist |

### 1.1 Trust zones and principal classes

**Zones:** untrusted browser/mobile client; public TLS ingress; Go control plane; durable metadata/operation store; privileged local Unix socket; host agent; tenant workload/process/container; local secrets; backup/offsite object store; external gateways (M-Pesa, Stripe), DNS registrars, mail and CDN; AI provider/tool gateway; provider support; remote multi-server agent in v2+. Each transition requires an authenticated principal, explicit scope, schema validation, authorization, bounded resources and an audit correlation identifier.

**Principal classes:** anonymous visitor; tenant owner; tenant admin; tenant member; developer/operator; reseller; provider operator; time-limited support delegate; automation token; scheduled job; control-plane service; enrolled agent; third-party integration. **Persona and UI display mode are not authorization roles.** Tenant ownership, resource ownership, delegated scope, capability, entitlement, action sensitivity and approval status are all evaluated server-side.

### 1.2 Security modes across release phases

| Deployment | Trust/channel | Isolation | Data and recovery | Non-negotiable limitation |
|---|---|---|---|---|
| Nano (target 512 MB) | Local control plane ↔ agent Unix socket | OS-user separation, file ownership and cgroup v2 | SQLite, encrypted remote backup, no local-only recovery | No default local Docker/DB/mail/platform-service promise |
| Starter | Local socket | OS principals; stronger sandbox when measured feasible | Encrypted offsite, tested restores | No claim of isolation equivalence with Docker |
| Standard | Local socket | Docker for supported workloads plus seccomp/AppArmor, network and cgroup constraints | PostgreSQL and remote WAL/PITR where configured | Docker is defense in depth, not a security boundary by itself |
| Provider v2 | Enrolled remote agents using mTLS with rotation/revocation | Per-node policy, per-tenant isolation, scheduler constraints | Coordinated offsite backup and control-plane recovery | Remote agents cannot trust client-provided scope |
| HA v3+ | Multi-control-plane with fencing and quorum-aware state | Same per-node invariants | Replication/failover, PITR and verified fencing | HA is not a v1 feature |

## 2. Authentication, sessions and access governance

**SEC8-013:** Admin enrollment MUST require verified ownership of the installation and a strong initial credential. Privileged roles MUST support MFA and recovery codes; release policy determines which roles require MFA at first login. Recovery codes are single-use and never stored reversibly.  
**SEC8-014:** Browser authentication MUST use secure, HttpOnly, SameSite-appropriate cookies, CSRF defenses for state-changing requests, short-lived sessions where risk warrants and server-side invalidation. Token-in-localStorage is not the default design.  
**SEC8-015:** API keys MUST be scoped, expiring or revocable, stored hashed where possible, and displayed only at issuance. Rotation and last-used metadata MUST be exposed.  
**SEC8-016:** Invite acceptance MUST bind the invited address/principal, organization, role and single-use token; reissued invites invalidate prior tokens.  
**SEC8-017:** Provider support access MUST require reason, scope, expiry, consent or policy basis, visible banner, independent audit and emergency revocation. Provider status does not confer silent access to customer secrets or content.  
**SEC8-018:** Break-glass access MUST be isolated from normal login, time-bounded, independently logged, reviewed after use and tested before a real incident.  
**SEC8-019:** Sensitive actions (delete, export, DNS transfer, firewall/SSH, restore over production, change billing destination, secret reveal, enrollment/rotation) MUST require contextual confirmation and step-up authentication according to risk.

### 2.1 Authorization decision inputs

`decision = policy(subject, tenant, resource, action, ownership, delegation, entitlements, capability, approval, operation_state, environment)`; the result is **allow** or a stable deny reason. The frontend may present `unauthorized`, `unavailable`, `disabled`, `not shipped`, `degraded` or `offline`, but only backend policy can grant execution. Reseller scope is bounded by upstream organization grants and quotas. A suspended tenant's existing hosted sites may remain running per policy while management mutations are denied; this state must be explicit.

### 2.2 Authorization matrix (minimum negative tests)

| Actor | Read own site | Read other tenant site | Change own DNS | Execute root shell | Approve own high-risk plan | Access provider audit |
|---|---|---|---|---|---|---|
| Tenant member (read only) | Scoped yes | No | No | No | No | No |
| Tenant admin | Yes | No | If policy permits | No | If approval policy permits | No |
| Reseller | Delegated only | No outside delegation | Delegated only | No | Policy-dependent | No |
| Provider operator | Metadata per duty | Not automatically content | Governed support path | No public API shell | Separation of duties | Scoped yes |
| Agent service identity | Target-bound operation only | No free browsing | Typed approved mutation | No arbitrary shell | Never self-approves | Append operation evidence |
| AI assistant | Read-only scoped diagnostics | No | Proposal only in v1 | Never | Never | No by default |

## 3. Privileged agent, host and workload hardening

**SEC8-020:** v1 control-plane↔agent communication uses a local Unix socket with restrictive ownership/mode, peer credential verification and authenticated, typed operation envelopes. A local socket is not inherently sufficient authorization. Remote mTLS is a later phase and must not be presented as v1 deployment.  
**SEC8-021:** The agent MUST reject unknown operations, unsupported fields, arbitrary executable paths, unsafe archive extraction, symlink escapes, path traversal, host root bind mounts, privileged Docker mode, host PID/network namespaces, dangerous capabilities and unrestricted device access.  
**SEC8-022:** All managed file operations MUST be rooted in an approved tenant-owned directory using race-resistant, descriptor-relative path resolution where supported; never rely on string prefix checks alone.  
**SEC8-023:** Workloads MUST run as least-privilege identities; process/container resource accounting uses cgroup v2. `memory.max` MUST enforce ceilings, `cpu.weight` shares CPU, and `cpu.max` is reserved for batch runaway protection under the proposed policy. The agent MUST reject `cpuset`/core-pinning configurations.  
**SEC8-024:** The agent MUST re-evaluate plan digest, target ownership, approval, policy version, server generation, idempotency and capability before applying a privileged step; stale queued work is denied or re-planned.  
**SEC8-025:** Managed configuration changes MUST use validate→stage→atomic apply→health check→commit or rollback, with drift markers and journaled recovery. Manual edits are detected and surfaced; NexusPanel must not silently overwrite unmanaged state.  
**SEC8-026:** Host hardening MUST include patch policy, least-privilege sudo, secure SSH keys, firewall safe-apply, service isolation, restrictive filesystem permissions and monitoring of suspicious process/permission changes.

### 3.1 SSH and firewall lockout-safe workflow

1. Identify current SSH session, provider console route, existing firewall policy and new rule effect; block if the independent recovery route cannot be established.
2. Create/verify non-root sudo account and key login in a separate session before restricting root or password login.
3. Present exact changes, exposed ports, source IP limitations, service dependencies and timed rollback deadline.
4. Persist a recovery timer **outside** the connection-dependent request path, apply atomically where possible and test connectivity from a separate vantage point.
5. Commit only after explicit verified confirmation; otherwise revert automatically and record audit.
6. If verification is inconclusive, show `needs-attention` and keep provider console instructions available. Never claim success from an API 200 alone.

### 3.2 Container and template policy

Require pinned image digest, signature/provenance policy, SBOM or documented exception, reviewed environment schema, non-root execution where feasible, seccomp/AppArmor, read-only filesystem where compatible, restricted mounts/network, bounded pids and memory, health checks and backups. Third-party templates are untrusted until reviewed. `ServiceTemplate` manifests cannot introduce new privileged agent verbs by adding JSON fields. Stateful application rollback requires service-specific data compatibility checks; replacing a container image alone does not restore the application.

## 4. Application, API and integration security

**SEC8-027:** API entrypoints MUST validate request size, content type, schema, canonical identifiers and tenant/resource ownership; use prepared SQL statements and output encoding.  
**SEC8-028:** SSRF controls MUST cover URL fetches, image transforms, webhook probes, importers, Git endpoints and AI tools; deny loopback, link-local, metadata services, private networks and DNS-rebinding targets except for explicitly reviewed internal service adapters.  
**SEC8-029:** Webhooks MUST verify provider authenticity using the supported provider-specific mechanism, event identity, merchant/account identity, amount/currency, invoice association and replay protection; an HTTP callback alone cannot settle a ledger.  
**SEC8-030:** Outbound webhooks MUST be signed, scoped, retried with backoff and idempotency guidance, and protected against leaking secrets or internal network addresses.  
**SEC8-031:** Rate limiting, abuse quotas and admission control MUST apply to login, token creation, DNS changes, builds, uploads, mail, AI requests, transforms, backup and restore.  
**SEC8-032:** Every privileged mutation MUST be an idempotent durable operation with step-level verification and no duplicate billing/provisioning on retry.  
**SEC8-033:** CORS/CSP/TLS settings, security headers and safe redirect rules MUST be versioned and tested; browser-route `GET /` is not an authorization API.  
**SEC8-034:** Payment and commercial secrets belong in the control-plane integration boundary, never in the host agent or tenant app by default.

### 4.1 M-Pesa and payment-state threat scenarios

Reject forged or replayed callbacks, mismatched merchant IDs, incorrect currency/amount, cross-tenant invoice references and unauthorized state transitions. Preserve raw gateway evidence in restricted storage, correlate provider transaction IDs, enforce unique settlement semantics and support explicit reversal/refund events. Gateway authentication specifics are **integration-specific validation items**; do not claim every provider supports the same signature scheme. See Volume 09 for commercial accounting and entitlements.

## 5. Secrets, encryption, data privacy and retention

**SEC8-035:** Secrets MUST use envelope encryption or equivalent documented design; key-encryption keys, data-encryption keys, backup keys and agent enrollment credentials have separate lifecycles and access policies.  
**SEC8-036:** Document local key custody for Nano and optional external KMS for larger deployments; loss of the only local master key must be described as a recovery risk, not concealed.  
**SEC8-037:** Rotation MUST support overlap windows, dependent-service reconfiguration, revocation and recovery after partial failure.  
**SEC8-038:** Secrets MUST never appear in telemetry, support exports, audit before/after payloads, screenshots, browser analytics, AI prompts or error traces; redaction is tested with canary values.  
**SEC8-039:** Define data classification for public service metadata, tenant private content, identity records, payment records, operational logs, secret material and audit evidence. Retention is per category, tenant and legal basis.  
**SEC8-040:** Support tenant export, deletion, retention holds and access requests without violating audit/legal retention or exposing other tenants' data. Deletion MUST propagate to caches, indexes and backup expiry according to policy, with documented exceptions.  
**SEC8-041:** Cross-border storage/processing, third-party subprocessors and breach notifications require a legal/compliance review before production. For Kenya-focused operations, assess applicable data-protection obligations; do not assert certification or legal compliance based solely on this document.

### 5.1 Data-classification matrix

| Class | Examples | Access | Encryption | Retention/disposal |
|---|---|---|---|---|
| Public | Published website assets | Intended public scope | TLS in transit | Per hosting lifecycle |
| Tenant confidential | Files, mail, databases, backups | Tenant-scoped roles, controlled support | At rest where supported + TLS | Customer policy, deletion and backup expiry |
| Highly restricted | Passwords, API tokens, private keys, backup keys | Just-in-time least privilege | Envelope encryption; one-time display | Rotation/revocation, secure deletion |
| Regulated/commercial | Invoices, settlement references, identity data | Role-specific, purpose-limited | Encrypted storage and transport | Applicable retention schedule |
| Security evidence | Audit chain, incident artifacts | Security operators, independent review | Integrity protected, encrypted exports | Non-tenant-writable protected retention |
| Operational telemetry | Metrics, logs, traces | Scoped operators and tenant owners | Transport + storage controls | Bounded by privacy and cost policy |

## 6. Audit, observability and evidence integrity

**SEC8-042:** Audit events MUST capture UTC timestamp, actor/service identity, tenant, target, action, decision, approval, immutable plan hash, operation ID, before/after summary (redacted), outcome, agent identity and policy version where applicable.  
**SEC8-043:** v1 audit MUST be hash-chained with `previous_hash` and `record_hash`, canonical serialization, chain verification, checkpoints and independent export. Hash chaining makes tampering **detectable** when anchors are protected; it is not equivalent to immutable WORM storage.  
**SEC8-044:** Define a single authoritative audit writer/serialization and concurrency policy; do not permit the control plane and agent to append independently to the same global chain without a sequencer or partition strategy.  
**SEC8-045:** Missing audit delivery, broken chains, clock anomalies, log redaction failures and backup export failures are visible security incidents. Sensitive operations MAY fail closed when mandatory evidence cannot be recorded.  
**SEC8-046:** Tenant deletion MUST NOT erase retained security evidence; support export MUST redact unrelated tenants and secrets.  
**SEC8-047:** Monitoring must distinguish `healthy`, `degraded`, `unknown`, `offline`, `disabled`, `unsupported` and `not shipped`, and display last-observed timestamps.  
**SEC8-048:** Resource metrics include cgroup current/typical/peak, OOM and throttling events; capacity recommendations are advisory and measured, not promises.

### 6.1 Audit record shape (illustrative)

```json
{
  "event_id": "aud_01EXAMPLE",
  "occurred_at": "2026-11-06T15:02:31Z",
  "actor_type": "tenant_admin",
  "actor_id": "usr_example",
  "tenant_id": "org_example",
  "server_id": "srv_example",
  "action": "firewall.rule.apply",
  "target_id": "fw_example",
  "operation_id": "op_example",
  "plan_digest": "sha256:<digest>",
  "policy_version": "v1",
  "decision": "allow",
  "outcome": "verified",
  "before_digest": "sha256:<redacted-state-digest>",
  "after_digest": "sha256:<redacted-state-digest>",
  "previous_hash": "sha256:<previous>",
  "record_hash": "sha256:<current>"
}
```

**Integrity workflow:** canonicalize → sequence → hash → append durably → periodically anchor/export to independently controlled storage → verify full chain and anchors → alert on gap/reorder/mismatch. Clock sync is mandatory; use server-generated timestamps and durable sequence IDs. The independent export needs access separation, retention and tested reconstruction. A missing agent event must not be quietly replaced by a synthetic success event.

## 7. Backup architecture and recovery-point taxonomy

**SEC8-049:** Define separate backup sets for (a) control-plane metadata, (b) control-plane keys/enrollment, (c) host-managed configuration, (d) tenant websites/files, (e) databases, (f) mail, (g) ServiceTemplate volumes/state, (h) audit chain, and (i) external integrations/configuration. No single database snapshot restores all these categories.  
**SEC8-050:** Every recovery point records scope, tenant, server, consistency method, start/end, status, destination, checksum/manifest, encryption key reference, retention/hold, verification status, dependencies and restore prerequisites.  
**SEC8-051:** Backups MUST be encrypted in transit and at rest, transferred to independent failure domains, and protected from tenant or compromised source-host deletion where target supports it.  
**SEC8-052:** A backup is `completed` only when transfer and manifest verification pass; it is `restore_verified` only after a successful representative restore drill. The UI MUST show both independently.  
**SEC8-053:** Scheduled backups are batch workloads, subject to fair queueing and resource ceilings, and MUST yield to serving workloads. Backup failure and lag MUST alert.  
**SEC8-054:** Retention policies MUST enforce legal holds, immutability constraints, object versioning where supported and a minimum recoverability window; deletions require authorization and audit.  
**SEC8-055:** A restore over production MUST preview conflicts, overwrites, downtime, DNS implications, credential changes and permission effects; require explicit approval and pre-restore checkpoint where feasible.  
**SEC8-056:** Restore into an isolated new destination SHOULD be the default for verification and incident recovery.  
**SEC8-057:** Recovery drills MUST exercise keys and identities, not merely file extraction, and MUST validate application behavior and authorization boundaries.

### 7.1 Backup consistency by asset

| Asset | Minimum consistent capture | Verification | Recovery caveat |
|---|---|---|---|
| SQLite metadata | Transactionally consistent backup API or safe checkpoint/copy strategy | Open DB, integrity check, replay where applicable | WAL files must not be naively copied mid-write |
| PostgreSQL metadata | Base backup + WAL archiving for PITR where configured | Restore to isolated DB and query invariants | WAL retention and key availability govern RPO |
| Website files | Snapshot or staged copy with change tracking | Hashes, ownership and serving probe | Dynamic uploads may need quiesce |
| Customer DB | Engine-native consistent dump or snapshot | Import and integrity query | Logical dump alone may omit external files |
| Mail | Mailbox state plus account/config metadata | Message count and sample retrieval | Delivery during restore needs cutover handling |
| Stateful app | Application-aware quiesce + volumes + DB + secrets references | App login/read/write health | OCI image is not user data backup |
| DNS/TLS | Zone/config export plus key custody | Zone validation, certificate/key match | External registrar/provider may need separate restore |
| Audit | Chain export and protected anchors | Chain recomputation | Primary DB restore cannot establish integrity without anchors |

### 7.2 Nano versus Standard recovery policies

The source v2.0 data specification proposes **Nano SQLite remote copy hourly (RPO target 1 hour)** and **Standard+ PostgreSQL WAL/PITR (RPO target 5 minutes)**. These are **design objectives**, not guarantees; measured RPO depends on backup frequency, WAL shipping, network, durability and successful restore. **RTO is not established by these numbers** and must be measured per workload. Nano remote-only backup policy is preserved; no same-disk backup is presented as disaster recovery. On Standard+ local backups may be used for fast recovery, but an offsite copy is required to claim host-loss recoverability.

### 7.3 Recovery-point state machine

`scheduled → running → uploaded → manifest_verified → restore_verification_pending → restore_verified`; failure branches include `failed`, `partial`, `corrupt`, `expired`, `verification_failed`, `key_unavailable`. A backup can be `uploaded` but not `restore_verified`. Restore verification may expire after a policy window, transitioning to `verification_stale`; this does not imply data loss but prevents a green recovery assurance badge.

## 8. Disaster recovery plans and recovery runbooks

**SEC8-058:** Define incident tiers: individual service failure, tenant data corruption, agent failure, control-plane metadata loss, host loss, region/provider outage, credential/key compromise and malicious deletion/ransomware. Each requires a different restoration sequence.  
**SEC8-059:** Publish a dependency graph for restoration: secure infrastructure → identity/keys → metadata → agent enrollment → managed config → tenant data → application verification → DNS/traffic → audits/notifications.  
**SEC8-060:** Before restoring, preserve forensic evidence where safe, isolate compromised credentials and prevent the recovery process from replaying malicious automation.  
**SEC8-061:** After restoring metadata, reconcile desired state against actual host observations without blindly deleting orphaned resources.  
**SEC8-062:** Restore operations MUST be idempotent, journaled and resumable, with step-level evidence and clear manual intervention states.  
**SEC8-063:** RPO/RTO reports MUST state measured recovery window, actual last valid restore point, unrecoverable interval and known exceptions; never claim zero loss without proof.

### RUNBOOK-01 — Compromised panel account

1. Trigger incident and preserve audit/session evidence.
2. Revoke active sessions and tokens, suspend sensitive operations and step up administrator authentication.
3. Review affected tenant/resource graph, queued operations and agent journal; block stale approvals.
4. Rotate impacted secrets and integration credentials, inspect unexpected DNS, mail, firewall, backups and billing changes.
5. Restore known-good state only after validating provenance; verify tenant isolation and user access.
6. Notify affected parties according to incident policy; document timeline, blast radius and corrective action.

### RUNBOOK-02 — Agent offline or corrupted

1. Show `offline` rather than `healthy`; hosted services continue where possible.
2. Check socket ownership, agent service, disk, journal and protocol compatibility via independent console path.
3. Restart/reinstall only a signed compatible binary; retain journal and previous version.
4. Reconcile unfinished operations with fencing/lease checks and service observations.
5. Validate one read-only probe, one safe mutation and audit export before declaring recovered.

### RUNBOOK-03 — SQLite metadata loss on Nano

1. Secure host and determine whether corruption, disk loss or credential compromise occurred.
2. Retrieve latest **verified remote** metadata backup and matching encryption/key material.
3. Restore transactionally consistent DB, replay available WAL only when supported and validated.
4. Re-enroll/reconcile agent identity as required; do not assume restored token is safe after compromise.
5. Reconcile websites, DNS, TLS and jobs with actual host state; quarantine unexpected objects.
6. Validate authentication, tenant ownership, service probes and audit chain export; publish measured RPO/RTO.

### RUNBOOK-04 — PostgreSQL metadata loss on Standard+

1. Freeze new mutations and preserve WAL/forensic state.
2. Restore encrypted base backup to isolated PostgreSQL instance; replay verified WAL to approved recovery point.
3. Validate schema version, ledger invariants, authorization graph, operation journal and audit anchors.
4. Switch control plane only after consistency and policy checks; reconcile agents with fencing tokens.
5. Re-enable writes gradually and run synthetic hosting, payment and backup verification checks.

### RUNBOOK-05 — Complete VPS loss

1. Declare disaster and choose a new trusted host with compatible OS, architecture, capacity, DNS and networking.
2. Bootstrap from signed artifacts with independent administrator enrollment and verified console access.
3. Restore keys and metadata under controlled access; re-enroll agent and restore managed config.
4. Restore tenant data/service volumes with application-specific consistency methods.
5. Verify sites, databases, mail, DNS, certificates, backup targets and isolation before traffic cutover.
6. Change DNS/routes with monitored rollback window; revoke obsolete host credentials and document recovery gap.

### RUNBOOK-06 — Ransomware or backup deletion attack

1. Isolate compromised host/credentials and disable destructive automation.
2. Protect immutable/offsite copies and audit anchors with independent operator access.
3. Identify clean restore point using malware and timeline evidence; do not restore compromised secrets unchanged.
4. Restore into clean infrastructure and validate security controls, tenant data and applications.
5. Rotate all exposed keys, enforce new recovery policies and record incident review.

### RUNBOOK-07 — ServiceTemplate stateful recovery

1. Confirm template revision, OCI digest, manifest, volume map, DB engine and compatible application version.
2. Provision isolated destination with appropriate quotas and networking; import **data and metadata**, not image only.
3. Restore application-consistent DB/volumes, rebind secrets and domains, run migration/health probes.
4. Test representative user flows and rollback/cutover plan; record new restore baseline.

## 9. Incident detection, response and communication

**SEC8-064:** Maintain severity levels `SEV1` (platform-wide compromise/data loss), `SEV2` (major tenant/service outage), `SEV3` (limited degraded operation), `SEV4` (minor issue); severity criteria and response ownership require formal approval.  
**SEC8-065:** Alert sources include auth anomalies, policy denies, agent enrollment changes, audit-chain break, secret-access spikes, unexpected firewall/DNS modifications, backup failure, OOM/resource pressure, malware signals and suspicious payment events.  
**SEC8-066:** Incident lifecycle: `detected → triaged → contained → eradicated → recovered → verified → closed`; reopening and postmortem are supported.  
**SEC8-067:** Evidence collection MUST use access-controlled immutable or tamper-evident storage with time synchronization, chain of custody and privacy minimization.  
**SEC8-068:** Notifications MUST distinguish confirmed impact from investigation, identify affected scope and actions, and avoid disclosing other tenants' information.  
**SEC8-069:** Provider, reseller and tenant incident views MUST reflect role scope; support cannot silently inspect customer data.  
**SEC8-070:** Postmortems MUST document detection delay, root cause, control failures, measured RPO/RTO, customer impact, fixes, owners and retest evidence.

### 9.1 Operational escalation matrix

| Severity | Example | Immediate containment | Recovery proof |
|---|---|---|---|
| SEV1 | Confirmed root/agent compromise, widespread data loss | Freeze privileged writes, revoke credentials, isolate affected nodes | Clean rebuild, security validation, independent audit |
| SEV2 | Control-plane outage, offsite backup unavailable, multiple tenants down | Stop unsafe retries, preserve data, activate recovery runbook | API/agent reconciliation, representative service probes |
| SEV3 | One workload offline, partial delivery failure | Isolate service and queue safe retry | Health probes and tenant confirmation |
| SEV4 | UI warning, low-risk observability issue | Triage and monitor | Verified fix and no security regression |

## 10. Supply-chain, patch and vulnerability management

**SEC8-071:** Produce versioned SBOMs for control plane, agent, frontend, templates and container images; pin dependencies and record provenance.  
**SEC8-072:** Verify signed releases before install/upgrade; reject untrusted artifacts, signature mismatch and downgrade outside an approved recovery path.  
**SEC8-073:** Maintain vulnerability triage by exploitability, exposure, tenant impact and remediation availability; publish supported OS/runtime/service matrix.  
**SEC8-074:** Updates MUST have migration dry-run, pre-upgrade backup, compatibility check, staged rollout, health gates and rollback plan; data-schema rollback is not assumed reversible.  
**SEC8-075:** Plugin/template publisher trust, manifest signatures, digest pinning and runtime capability review are required.  
**SEC8-076:** Security patches SHOULD prioritize exposed critical components, but exact SLA must be approved and operationally demonstrated before marketing claims.  
**SEC8-077:** The product MUST publish a vulnerability reporting process and security-contact channel before public launch.

## 11. Compliance and legal review gates

The product targets Kenya initially but may serve customers elsewhere. The following are **review checkpoints, not claims of legal advice or completed compliance**.

| Gate | Required work product | Responsible owner | Release evidence |
|---|---|---|---|
| Privacy governance | Data inventory, purpose, lawful basis, rights handling, privacy notice | Product/legal | Approved policies and workflows |
| Kenya data protection | Applicability assessment, registration/role assessment where relevant, cross-border processing | Qualified counsel/privacy lead | Written legal review |
| Security governance | Threat model, access reviews, vulnerability process, incident response | Security engineering | Test reports and runbooks |
| Payment security | M-Pesa/Stripe contractual, webhook and reconciliation controls | Commerce/security | Sandbox/production acceptance |
| Email abuse | Anti-spam, bounce handling, provider restrictions and customer acceptable use | Operations/legal | Abuse tests and escalation path |
| Data retention | Category-based schedule, legal hold and deletion workflows | Privacy/security | Retention tests |
| Processor/subprocessor | Vendor inventory, access, region, retention and contractual terms | Legal/procurement | Reviewed agreements |
| Backups and residency | Destination selection, key custody, regional considerations | Infrastructure/privacy | Restore drills and destination inventory |
| Customer transparency | Accurate tier capability statements, backup/RPO/RTO disclosures | Product | Reviewed UX and terms |

Do not state that NexusPanel is ISO 27001, SOC 2, GDPR compliant, Kenya DPA compliant, PCI certified or otherwise certified without an actual completed independent assessment where required. Define incident-notification deadlines only after legal review of applicable obligations and contracts.

## 12. Security UX and operator controls

**SEC8-078:** Security center MUST present MFA, SSH, firewall, patches, certificates, backup assurance, active sessions, API keys, agent identity and audit health as separate status items.  
**SEC8-079:** Critical changes MUST show target, blast radius, dependency impact, downtime, reversible/irreversible steps and rollback deadline.  
**SEC8-080:** Distinguish `permission denied`, `hardware unsupported`, `not installed`, `not shipped`, `temporarily degraded`, `offline` and `requires approval`; provide next action without exposing sensitive system details.  
**SEC8-081:** Backup UI MUST show last attempted backup, last valid uploaded point, last restore-verified point, target region/failure domain and next scheduled verification.  
**SEC8-082:** Security notifications MUST be actionable and avoid green success states until verification completes.  
**SEC8-083:** Mobile emergency actions (revoke token, suspend account, view incident, safe restart) MUST remain accessible without exposing high-risk workflows to accidental taps.  
**SEC8-084:** Audit browsing and support exports MUST be role-filtered, paginated, searchable and privacy-redacted.

## 13. Cross-area security and recovery ownership matrix

| Area | Primary security obligation | Recovery dependency |
|---|---|---|
| 01 Dashboard | Scoped observability, freshness and safe summary | Cached status and unknown/offline distinction |
| 02 Websites | Tenant file roots, safe Nginx config | Files + config + TLS |
| 03 Domains | Ownership proof, DNS change approval | Zones, external registrar access |
| 04 Email | Mailbox isolation, abuse prevention | Mail stores, routing and credentials |
| 05 Databases | Tenant DB principals, injection defenses | Engine-native backups, consistency |
| 06 File manager | Path safety, upload scanning | File snapshots, permissions |
| 07 Applications | Sandboxed build and runtime | Code, env refs, persistent data |
| 08 Runtimes | Pinned packages, restricted execution | Runtime compatibility and configs |
| 09 SSL | Private-key custody, renewal safety | Key/cert recovery and issuance |
| 10 DNS | Zone authorization, anti-takeover | Zone export and propagation verification |
| 11 Monitoring | Scoped logs, redaction | Alert history and evidence |
| 12 Backups | Encryption, immutable/offsite copies | Restore verification |
| 13 Security | MFA, firewall, SSH, audit | Break-glass and clean host rebuild |
| 14 Automation | Approval, scope, loop prevention | Job journal and safe retry |
| 15 Users/RBAC | Membership, delegated scopes | Identity and access graph |
| 16 Resellers | Quota and delegation bounds | Customer hierarchy recovery |
| 17 Billing | Ledger integrity, webhook verification | Financial reconciliation |
| 18 Subscriptions | Entitlement state and idempotency | Event replay/reconciliation |
| 19 Licensing | Signed entitlement/capability claims | License validation recovery |
| 20 AI | Prompt injection isolation, approved plans | Human-controlled remediation |
| 21 API | Auth, schema, rate limits, idempotency | Versioned contract and replay |
| 22 Migration | Integrity, secret rotation, safe cutover | Application-consistent transfer |
| 23 Server management | Privileged operation policy | Agent journal and server recovery |
| 24 Multi-server | Agent identity, fencing | Fleet placement and node failover |
| 25 Marketplace | Signed provenance and sandbox | Template revisions and state |
| 26 Support | Consent, expiry, redacted evidence | Incident case continuity |
| 27 Provider ops | Separation of duties, abuse controls | Control-plane and provider DR |
| 28 Platform services | Storage ACLs, transform SSRF, queue isolation | Object data, cache rebuild, search index |
| 29 Self-hosted templates | Signed images, sandbox, secrets | App-aware DB/volume restoration |

*Area names in this matrix are descriptive groupings for traceability; the authoritative 29-area register and Volume 01 govern exact area labels and IDs.*

## 14. Verification and adversarial acceptance suite

| Test ID | Scenario | Expected result |
|---|---|---|
| SEC8-T001 | Guess another tenant's resource ID | 403/404 per disclosure policy; no data leak |
| SEC8-T002 | Change UI persona to Provider | No additional API permissions |
| SEC8-T003 | Revoke role after job queued | Agent rejects stale privileged step |
| SEC8-T004 | Replay approved plan with modified target | Digest mismatch; no mutation |
| SEC8-T005 | Inject shell metacharacters in hostname/path | Schema rejection; no shell execution |
| SEC8-T006 | Symlink swap during file write | Resolver stays in authorized root |
| SEC8-T007 | Submit Docker privileged/host-network/host-root mount | Rejected before runtime apply |
| SEC8-T008 | Set cgroup cpuset or invalid memory limit | Rejected and audited |
| SEC8-T009 | Force OOM in tenant workload | Other tenants and control plane remain within tested isolation envelope |
| SEC8-T010 | Prompt-inject AI tool output | AI cannot bypass scoped tools or approval |
| SEC8-T011 | Insert secret canary into logs/errors | Redacted from all outputs and support exports |
| SEC8-T012 | Steal/replay API token after revocation | Request and queued work denied |
| SEC8-T013 | Forge/replay payment callback | No double credit, no unauthorized ledger state |
| SEC8-T014 | SSRF through image URL/webhook/import | Private/metadata destinations blocked |
| SEC8-T015 | Tamper with one audit record | Chain verification alerts |
| SEC8-T016 | Remove audit DB but keep independent anchor | Reconstruction detects missing segment |
| SEC8-T017 | Disconnect SSH after firewall apply | Timed rollback restores access |
| SEC8-T018 | Kill agent during multi-step mutation | Journal recovery; no silent partial success |
| SEC8-T019 | Delete local disk with only local backup | DR gate fails; no recoverability claim |
| SEC8-T020 | Corrupt offsite backup manifest | Restore point invalid; alert emitted |
| SEC8-T021 | Restore SQLite snapshot without WAL consistency | Verification fails rather than declaring recovered |
| SEC8-T022 | PITR PostgreSQL to chosen timestamp | Tenant, ledger, job and schema invariants hold |
| SEC8-T023 | Restore stateful app from OCI image only | Validation rejects incomplete data restore |
| SEC8-T024 | Restore with missing encryption key | Explicit blocked/key-unavailable state |
| SEC8-T025 | Expire restore verification window | UI downgrades confidence to unverified/stale |
| SEC8-T026 | Compromised agent enrollment credential | Revoke, rotate and quarantine old agent |
| SEC8-T027 | Provider support without consent/valid policy | Denied and audited |
| SEC8-T028 | Export audit as tenant | No other-tenant or secret exposure |
| SEC8-T029 | Tamper signed template/image | Signature/digest rejection |
| SEC8-T030 | Deploy on target 512 MB Nano | Measured resource floor and core controls pass; otherwise unsupported |
| SEC8-T031 | Block port 25 / remove PTR | Email capability correctly unavailable/degraded |
| SEC8-T032 | Backup batch under live site load | Serving latency remains within measured SLO, batch yields |
| SEC8-T033 | Full host-loss restore drill | Document actual RPO/RTO and failed steps |
| SEC8-T034 | Key rotation with in-flight backup | New and old encrypted points recover according to policy |
| SEC8-T035 | Provider outage during payment settlement | Async state reconciles without double posting |
| SEC8-T036 | Unsupported Ubuntu 26.04 package on installer | Preflight blocks or clearly gates module; no silent assumption |

**Release acceptance:** Tests are executed against each supported OS/profile and release tier with documented environment, reproducible steps, owner, observed results, defects and signed approval. This document lists intended tests; it does not assert they have passed.

## 15. Decision and validation register

| Gate | Proposed default | Must be decided/verified before |
|---|---|---|
| SEC8-D01 | MFA mandatory for provider/admin and high-risk actions | Public beta |
| SEC8-D02 | Hash-chained audit + independent protected anchors | First production write |
| SEC8-D03 | Nano hourly remote SQLite, 1h target RPO | Nano production claim |
| SEC8-D04 | Standard+ Postgres PITR, 5m target RPO | Standard production claim |
| SEC8-D05 | RTO by tier remains uncommitted until drills | Any SLA publication |
| SEC8-D06 | Offsite immutable/WORM-compatible backup adapter | DR assurance marketing |
| SEC8-D07 | Local key custody with tested escrow/recovery | Production installer |
| SEC8-D08 | Docker default Standard+; OS user Nano | Public multi-tenant launch |
| SEC8-D09 | cgroup v2 required; cpu.weight and memory.max | Nano and Standard benchmark |
| SEC8-D10 | No arbitrary shell or privileged container options | Agent security review |
| SEC8-D11 | Exact incident severity SLAs and legal notice windows | Customer contracts |
| SEC8-D12 | Kenya and cross-border privacy legal assessment | Commercial launch |
| SEC8-D13 | Ubuntu 26.04.1 adapter/package compatibility | User's VPS installation |
| SEC8-D14 | Audit writer ownership and canonicalization algorithm | Schema freeze |
| SEC8-D15 | Support consent and break-glass controls | Provider access launch |
| SEC8-D16 | Recovery key escrow, deletion and legal hold | Backup launch |
| SEC8-D17 | Third-party template trust policy and SBOM minimum | Marketplace launch |
| SEC8-D18 | Agent enrollment revocation and remote mTLS | v2 multi-server |
| SEC8-D19 | HA quorum/fencing and split-brain recovery | v3 HA |
| SEC8-D20 | Penetration test scope and independent review | Public hosting launch |

## 16. Preserved v2.0 source contracts (requirements and engineering details)

The following sections retain the original source requirement identifiers and engineering detail so implementers can trace each integrated obligation back to the user-provided v2.0 specifications. These sections are not additional documents or new architectural approvals. Editorial v2 labels are normalized; the source's proposed values remain proposed.

### Source contract PRD-14: Backups, restore, snapshots, retention and disaster recovery

- **NP-BKP-001:** Create policies for websites, files, databases, mail and account-level backup sets with schedules and retention.
- **NP-BKP-002:** Support encrypted offsite targets via adapters, credentials separation, integrity manifests and transfer verification.
- **NP-BKP-003:** Show restore-point inventory, size, status, scope, encryption and last verified restore test.
- **NP-BKP-004:** Restore wizard previews overwritten paths, databases, permissions, domains, mailboxes and downtime risk.
- **NP-BKP-005:** Support partial restore to new location where service adapter permits, with isolation from active production.
- **NP-BKP-006:** Failed backups are visible and alertable; do not count them as valid recovery points.
- **NP-BKP-007:** Implement retention and deletion protection including legal hold/immutable targets where supported.
- **NP-BKP-008:** Support backup of control-plane metadata, secrets recovery material and agent enrollment state with secure handling.
- **NP-BKP-009:** Define RPO/RTO objectives per service tier; test restore procedures regularly and record evidence.
- **NP-BKP-010:** Backup and restore jobs are bounded by disk space, bandwidth, quotas and failure handling.

**Acceptance scenarios:** Restore deleted site/database from tested offsite backup; interrupted backup is marked invalid; retention job cannot delete protected recovery point.

**Consolidated requirement:** Additional requirements:

- **NP-BKP-011:** Backups are a **batch** workload. They run with `cpu.weight = 10` and `cpu.max = 50% of cores`, so a scheduled backup does not starve websites. They are queued, not rejected. On Nano, backups use restic to an external S3-compatible target (Backblaze B2, R2, etc.) — the panel does not host backup storage on the Nano VPS itself.
- **NP-BKP-012:** **Backup integrity is verified by restore test**, not by process exit code. The panel schedules periodic restore drills to a disposable target and records the result. "Backup succeeded" ≠ "backup is restorable." A backup without a recent verified restore is marked `unverified`, not `healthy`.
- **NP-BKP-013:** On Nano, backups are **remote-only**. There is no local backup destination — the disk is too small. The panel requires a remote backup target during onboarding. This is honest: a Nano VPS cannot backup to itself.
- **NP-BKP-014:** On Standard+, backups may be local or remote. Remote is strongly recommended; the panel shows the storage tier ("local" vs. "offsite") and warns if only local backups exist. "A backup on the same VPS alone is not disaster recovery."

### Source contract PRD-15: Security, firewall, isolation, certificates, and audit

- **NP-SEC-001:** Provide secure setup checklist for admin account, MFA, SSH keys, OS updates, firewall, recovery console and backup.
- **NP-SEC-002:** Firewall changes are previewed and guarded against locking out active administration; rollback/recovery path is required.
- **NP-SEC-003:** Support account isolation via distinct OS identities and file/process boundaries; stronger sandboxing for untrusted workloads according to threat model.
- **NP-SEC-004:** Enforce input validation and safe command execution in every service adapter; no user-controlled arbitrary shell interpolation.
- **NP-SEC-005:** Protect against CSRF, XSS, SSRF, SQL injection, path traversal, insecure file upload, token replay and privilege escalation.
- **NP-SEC-006:** Store secrets encrypted with documented key management and rotation; never log plaintext credentials.
- **NP-SEC-007:** Security events include failed login, role change, privileged action, secret access, agent enrollment and firewall modification.
- **NP-SEC-008:** Audit trail captures actor, tenant, action, target, decision, before/after summary, job reference and timestamp; restrict deletion/tampering.
- **NP-SEC-009:** Provide vulnerability and patch visibility with supported OS/service compatibility matrix and safe maintenance workflow.
- **NP-SEC-010:** Support incident containment (disable token, revoke session, isolate service, suspend tenant) with controlled blast radius.
- **NP-SEC-011:** Provider support access requires justification, consent/policy where applicable, expiry and audit.
- **NP-SEC-012:** Production release requires documented threat model, dependency scans, penetration testing and incident runbooks.

**Acceptance scenarios:** Unauthorized agent operation rejected; path traversal cannot escape tenant root; root SSH not disabled until alternate access verified; audit remains available after account suspension.

**Consolidated requirement:** Additional requirements:

- **NP-SEC-013:** **Docker is the primary isolation primitive** where Docker fits (Standard / Provider profiles). Container isolation provides filesystem, process, network, and resource separation. OS-user-per-account is the fallback for workloads that cannot containerize (Nano static sites).
- **NP-SEC-014:** **cgroup v2 is required.** Ubuntu 22.04+ / 24.04+ is the supported OS baseline. The installer refuses to install on systems without cgroup v2. This is DEC-015 in Master Blueprint §15.
- **NP-SEC-015:** **`cpuset` (core pinning) is never used.** It reserves cores, which wastes idle CPU. The scenario where a Node app wants 1.5 vCPU on an idle 4-vCPU server but gets throttled to 1.0 because transforms have "1 vCPU reserved" is a failure mode the architecture prevents. `cpu.weight` is the mechanism.
- **NP-SEC-016:** **`cpu.max` is set only on batch workloads** as a runaway protection. Serving workloads (websites, apps, databases) have no `cpu.max` — they can burst to all cores under idle, and share fairly by weight under contention.
- **NP-SEC-017:** The **audit chain is hash-chained from v1**. It cannot be retrofitted. Audit records survive tenant resource deletion under retention policy. The chain is exported to a protected storage boundary not writable by the tenant.
- **NP-SEC-018:** **Container escape prevention.** The agent validates every container configuration before applying: rejects `--privileged`, host PID/network namespace, host root bind mounts, dangerous Linux capabilities, and unrestricted device access. This is Stage 7 AGT-108/AGT-109 — restated here as a product invariant.

### Source contract PRD-16: Monitoring, logs, incidents, analytics and notifications

- **NP-OBS-001:** Show CPU, RAM, disk, I/O, network, process, service and per-account resource metrics where instrumentation supports them.
- **NP-OBS-002:** Offer searchable, access-controlled logs for web, apps, deployments, databases, mail, jobs and agent events.
- **NP-OBS-003:** Support alert thresholds, routing, deduplication, suppression, maintenance windows and acknowledgement.
- **NP-OBS-004:** Show service status with last-observed timestamp, source and confidence; distinguish unknown from healthy.
- **NP-OBS-005:** Correlate user actions, deployment IDs, job IDs and incidents for root-cause analysis.
- **NP-OBS-006:** Support notifications in panel and optional email/webhook channels, with delivery failure visibility.
- **NP-OBS-007:** Metered usage and operational metrics are separate datasets with appropriate retention and reconciliation.
- **NP-OBS-008:** Enforce log retention, cardinality/volume limits, PII and secret redaction, and tenant isolation.
- **NP-OBS-009:** Provide capacity forecasts and top resource consumers with transparent estimation limitations.

**Acceptance scenarios:** Failed deployment creates traceable event; high disk usage alerts without exposing other tenants' data; unavailable agent shows unknown/offline rather than green.

**Consolidated requirement:** Significant refinement:

- **NP-OBS-010:** Metrics come from **cgroup v2** (`/sys/fs/cgroup/<workload>/cpu.stat`, `memory.current`, `memory.peak`, `memory.events`, `io.stat`, `pids.current`), polled every 15–30 s. This is the same source the installer uses to compute capability fit. The dashboard shows current / typical (7-day p50) / peak (24h p99) / OOM events / throttled seconds / contention signal.
- **NP-OBS-011:** On Nano, monitoring is **minimal**: cgroup stats only, no Prometheus, no Grafana, no Loki. On Standard+, monitoring uses **VictoriaMetrics** (much lighter than Prometheus) plus optional log aggregation. The full Prometheus+Grafana+Loki stack is optional, not default, because the RAM cost is real on a 4 GB VPS.
- **NP-OBS-012:** The **upgrade recommendation logic** is a first-class feature, only possible because allocation is separate from usage (§5.5): peak RAM > 80% for 7+ days → "Consider upgrading"; OOM kills in audit log → "A workload was killed for memory"; CPU contention events → "Workloads are competing for CPU." The recommendations link to the upgrade flow (§9, Master Blueprint §10.5).

### Source contract ARCH-13: Backups, restore and disaster recovery

**ARCH-038:** Backup policy declares protected resources, consistency model, schedule, destination, encryption key ownership, retention, immutability/offsite requirement, restore dependencies and verification. Local backup on the same disk is not disaster recovery.

**ARCH-039:** Restore is a first-class durable operation with dry-run/impact preview, conflict policy, target mapping, pre-restore snapshot where possible, approval and post-restore verification. A restore must not silently overwrite unrelated tenant resources.

**ARCH-040:** Define RPO/RTO per tier and validate through periodic full restores on disposable infrastructure. Single-VPS mode explicitly documents that a destroyed host requires external backup and reinstallation.

**Consolidated requirement:** Additional rules:

- **Backups are a batch workload.** They run with `cpu.weight=10`, `cpu.max=50%`, `memory.high` throttling, `memory.max` hard cap. A scheduled backup does not starve websites. This is what makes backups safe to schedule on a small VPS.
- **Backup integrity is verified by restore test.** The panel schedules periodic restore drills to a disposable target and records the result. "Backup succeeded" ≠ "backup is restorable." A backup without a recent verified restore is marked `unverified`, not `healthy`.
- **On Nano, backups are remote-only.** The disk is too small for local backups. The panel requires a remote backup target during onboarding.

### Source contract ARCH-14: Security architecture and threat model

#### 14.1 Threat actors and assets

Threat actors: anonymous Internet clients, compromised tenant apps, malicious tenants, stolen API credentials, compromised plugin/adapter, compromised support identity, network attacker, malicious artifact and insider with delegated privileges. Assets: root control, tenant data, domain/DNS control, email contents, backups, secrets, billing records and audit evidence.

| Threat | Primary controls | Test evidence |
|---|---|---|
| Cross-tenant resource access | Object-level authorization, scoped queries, isolated OS principals, Docker | Negative authorization matrix and IDOR tests |
| Root compromise through agent | Typed allowlist, strict schemas, privilege separation, mTLS | Agent fuzzing, adversarial inputs, external review |
| Malicious upload/build | Sandboxed build, size/CPU limits, scan policy, safe archive extraction | Escape and resource exhaustion tests |
| Secret disclosure | Envelope encryption, masking, scoped reveal, no secret logs | DLP/redaction tests and key rotation drill |
| SSH/firewall lockout | Safe apply, independent access verification, timed rollback | Simulated failed connectivity test |
| DNS hijack | Ownership proof, approval, registrar integration scope, audit | Unauthorized transfer and mutation tests |
| Backup tampering | Offsite encrypted immutable copies, restore drills | Destruction and corruption recovery exercises |
| Supply chain compromise | Pinned dependencies, SBOM, signatures, provenance, scanning | Reproducible builds and signature verification |
| AI prompt injection | Untrusted context boundaries, plan approval, same policy engine | Tool-injection adversarial tests |
| Provider insider misuse | Least privilege, support consent, audit, separation of duties | Break-glass and abuse simulations |

**Consolidated requirement:** Additional threats:

| Threat | Primary controls | Test evidence |
|---|---|---|
| **Container escape** | Reject `--privileged`, host PID/network, dangerous capabilities; use rootless where possible | Malicious Compose tests, capability fuzzing |
| **cgroup misconfiguration** | Agent validates all cgroup values; rejects `cpuset` and unknown controllers | cgroup behavior tests |
| **cgroup escape via workload** | Non-root workload, no host PID, seccomp/AppArmor, read-only root filesystem | Container penetration tests |
| **Payment callback forgery** | Verify M-Pesa/Stripe signature, event ID, merchant identity, amount, currency | Negative-signature tests, replay tests |
| **Idempotency bypass via payment retry** | Idempotency key scoped to actor + invoice + amount; same key returns same result | Duplicate payment tests |
| **Capability spoofing** | Capability records read from signed data files, not client input | Installer tamper tests |

**ARCH-041:** Secrets are never stored in plaintext application tables; design an envelope encryption scheme with separate master-key lifecycle and per-tenant/per-resource authorization. External KMS may be optional; local key custody must document backup and compromise recovery.

**ARCH-042:** Public ingress uses TLS, secure cookies, CSP, CSRF defenses, input validation, request limits and audit. Security updates, incident response, vulnerability disclosure and dependency tracking are release requirements.

**ARCH-043:** All high-risk actions have explicit approval policy, contextual impact summary, re-authentication where justified and durable audit. Approval is bound to the exact immutable plan hash, not an editable future request.

**Consolidated requirement:** Additional security rules:

- **cgroup v2 is required.** Ubuntu 22.04+ / 24.04+ is the OS baseline. The installer refuses to install on systems without cgroup v2.
- **`cpuset` (core pinning) is never used.** It reserves cores, which wastes idle CPU. `cpu.weight` is the mechanism.
- **`cpu.max` is set only on batch workloads** for runaway protection. Serving workloads can burst to all cores under idle.
- **RAM gets a hard `memory.max` on every workload.** RAM cannot be overcommitted safely. The ceiling is enforced.
- **Docker is the primary isolation primitive** where Docker fits. The agent validates every container configuration and rejects privileged mode, host namespace, dangerous capabilities.

### Source contract ARCH-18: Observability, support and incident management

**ARCH-054:** Instrument API, worker, agent and data-plane services with correlated IDs, structured logs, metrics and distributed traces where supported. Apply secret and personal-data redaction at collection and export.

**ARCH-055:** Health states distinguish `unknown`, `healthy`, `degraded`, `unhealthy`, `offline`, `stale` and `maintenance`. Monitoring intervals and staleness thresholds are explicit and configurable.

**ARCH-056:** Alert rules include deduplication, routing, severity, acknowledgement, escalation and resolution evidence. Customer views are scoped; providers can inspect service health without blanket access to customer content.

**ARCH-057:** Support bundles are opt-in, time-limited and sanitized. Diagnostic export inventories files/fields and requires consent where customer data could be included.

**Consolidated requirement:** Additional observability rules:

- **Metrics come from cgroup v2.** `/sys/fs/cgroup/<workload>/cpu.stat`, `memory.current`, `memory.peak`, `memory.events`, `io.stat`, `pids.current`. Polled every 15–30 s. Stored in SQLite (Nano) or VictoriaMetrics (Standard+). Surfaced as current / typical (7-day p50) / peak (24h p99).
- **On Nano, monitoring is minimal.** cgroup stats only. No Prometheus, no Grafana, no Loki.
- **On Standard+, VictoriaMetrics is the default.** Much lighter than Prometheus. Full Prometheus+Grafana+Loki is optional, not default.
- **The upgrade recommendation logic is a first-class feature.** Peak RAM > 80% for 7+ days → "Consider upgrading." OOM kills in audit → "A workload was killed for memory." CPU contention events → "Workloads are competing for CPU." Only possible because allocation is separate from usage.

### Source contract ARCH-23: Critical failure scenarios and recovery behavior

1. **Agent unreachable:** keep requested state queued or blocked; do not claim applied; retry with bounded backoff; expose last observation timestamp.
2. **Worker crash mid-deploy:** lease expires; fencing/idempotency prevents double-apply; reconcile actual host state before continuing.
3. **Invalid web configuration:** syntax test fails; previous configuration remains active; show diff and safe repair.
4. **DNS delegation mismatch:** mark `pending_external`; show exact expected NS/glue/DS records and verification evidence; never silently mutate registrar.
5. **ACME rate limit:** preserve existing certificate; schedule retry with backoff and expiry alert.
6. **Disk full during backup:** backup fails without false success; preserve prior valid backup and notify.
7. **Firewall change drops connectivity:** staged rollback activates using host-local timer; operator verifies from independent connection.
8. **Database import interrupted:** retain import status and staging metadata; define restart or restore rather than assuming transactional rollback.
9. **Provider payment webhook duplicated:** idempotent ledger entry and entitlement update; no duplicate account provisioning.
10. **AI proposes unsafe command:** tool gateway rejects non-allowlisted action; audit denied attempt; offer manual safe workflow.
11. **Tenant deleted while job queued:** execution-time policy recheck cancels or compensates operation.
12. **Control-plane DB unavailable:** existing hosted sites should continue where possible; new management mutations fail safely; recover from documented backups.

**Consolidated requirement:** Additional failure scenarios:

13. **M-Pesa callback arrives late:** payment state remains `pending`; timeout per Daraja policy; on callback, resolve to confirmed/failed. If both callback and timeout fire, idempotency key ensures single resolution.
14. **M-Pesa callback never arrives:** scheduled reconciliation job queries Daraja Transaction Status API; resolves or escalates to manual reconciliation (Stage 3 UX-COM-017).
15. **Container fails to start:** `service_instances` state becomes `degraded`; health check surfaces; user can restart, inspect logs, or uninstall.
16. **cgroup write fails:** agent retries; if persistent, workload is marked `degraded`; user sees "memory limit could not be applied" with correlation ID.
17. **OOM kill:** workload restarts per its restart policy; audit records `oom_kill`; upgrade recommendation triggered if repeated.
18. **Template image signature invalid:** deployment rejected; audit records; user sees "image verification failed — contact template maintainer."
19. **Storage bucket quota exceeded:** writes fail with `QUOTA_EXCEEDED`; reads continue; user sees "bucket full — 5 GB used of 5 GB."
20. **Router path unavailable:** user picks a card the VPS cannot run; honest capability display; upgrade path offered.
21. **Upgrade detection fails after hardware bump:** user can trigger manual "Rescan hardware"; panel re-runs installer logic; shows delta.

### Source contract DATA-12: Audit, events and observability data

Audit events capture actor, impersonation/support context, organization, action, target, policy decision, request correlation, operation ID, timestamp, outcome and tamper-evident integrity linkage. Store redacted metadata, not full request bodies containing secrets. Export append-only records off-host with restricted deletion rights.

Transactional outbox rows are committed alongside resource intent. Publishers deliver at least once, record attempts and deduplicate downstream. Webhook deliveries have per-endpoint retry budgets, exponential backoff, replay controls and signature verification.

Metrics and high-cardinality traces should live in observability stores, not grow unbounded in control-plane transactional tables. Define rollups and retention for usage records, agent heartbeats, DNS checks and mail events.

**Consolidated requirement:** Additional rules:

- **The audit chain is hash-chained.** Each `audit_events` row has `previous_hash` and `record_hash`. `record_hash = sha256(previous_hash || actor_id || action || resource_id || occurred_at || payload_digest)`. A scheduled job verifies the chain hourly; a broken chain is a critical incident. The chain cannot be retrofitted.
- **cgroup metrics are stored in `inventory.cgroup_metrics` (hourly rollup) and VictoriaMetrics (raw series on Standard+).** On Nano, only the hourly rollup is stored (SQLite is too small for the raw series). The retention is: raw series 30 days in VictoriaMetrics, hourly rollup 1 year in the control plane DB, daily rollup 5 years.
- **Payment events are append-only.** `commercial.payment_events` records every state transition. The `payments` table has the current state. The events table is the audit trail. You cannot update a `payments` row without recording the transition in `payment_events` — this is enforced by a trigger (Postgres) or application-layer transaction (SQLite).
- **`operations.batch_queue`** is drained in priority order. The worker claims the next batch job; if the server is under contention (`cpu.stat.nr_throttled` elevated), it yields. This is how batch workloads stay low-priority.

### Source contract DATA-16: Backup, PITR and disaster recovery

Control-plane PostgreSQL requires encrypted full backups plus WAL archiving for point-in-time recovery where configured. Backups must leave the host. Define recovery point objective (RPO) and recovery time objective (RTO) by deployment tier; do not claim zero data loss. Include encryption key escrow and verified recovery access.

Recovery drill: restore isolated control-plane database → validate schema and tenant ownership → re-register/reauthenticate agents → reconcile desired/observed generations → verify a sample website, DNS zone, mailbox and database → compare audit chain and backup manifest → record achieved RPO/RTO.

Customer data backups are separate from control-plane PITR and must be recoverable independently. Restoring metadata does not magically restore deleted filesystem or mailbox payloads.

**Consolidated requirement:** Additional backup rules:

- **On Nano, there is no local PITR.** SQLite is copied to remote backup hourly by a batch job. Recovery is "restore the latest SQLite file, replay the WAL if available, accept the RPO." RPO is documented as 1 hour.
- **On Standard+, Postgres PITR is the default.** WAL is archived to remote object storage. RPO is documented as 5 minutes.
- **ServiceTemplate instances have their own backup policies.** A Nextcloud instance has its own volume backup; a database has its own dump schedule. These are independent of the control-plane PITR.
- **Backup verification is scheduled.** `storage.backup_verifications` records the result. A backup without a verification within the policy window is `unverified`, not `healthy`. The UI shows "last verified restore: 3 days ago."
- **The audit chain is backed up separately.** It is exported to a write-protected storage boundary. If the primary DB is lost, the chain can be reconstructed from the export.

### Source contract DATA-18: Privacy, compliance and threat cases

Threat scenarios to test: guessed resource UUID; cross-tenant joins; forged organization scope; reseller over-delegation; stale queued operation after revocation; idempotency replay; compromised API token; SQL injection; RLS bypass by privileged DB role; tenant export of another tenant's files; compromised agent; leaked backup key; forged payment webhook; malicious plugin; log/AI secret exfiltration.

Kenya-focused deployments should assess applicable data protection obligations and cross-border processing arrangements before production. This is an engineering compliance checkpoint, not legal advice.

**Consolidated requirement:** Additional threat cases:

- **M-Pesa callback forgery:** signature verification, event ID dedup, merchant identity check, amount/currency match. Tested with negative signatures and replay.
- **Payment state transition bypass:** a `pending` payment cannot jump to `confirmed` without a callback; a `confirmed` payment cannot be reversed without an explicit reversal event. Enforced by state machine.
- **cgroup misconfiguration:** a workload cannot be assigned `cpuset` or a memory ceiling above physical RAM × 1.2. Enforced by API and agent.
- **Container escape:** the agent rejects `--privileged`, host PID/network, dangerous capabilities. Tested with malicious Compose.
- **Capability spoofing:** capability records are loaded from signed data files. A tampered file is rejected at install time. The installer signature is verified.
- **Persona matrix tampering:** the persona matrix is data, but it is signed and verified on load. The API rejects unsigned data.

### Source contract API-19: Backups, restore, snapshots and disaster recovery

| Method | Endpoint family | Operations |
|---|---|---|
| GET/POST | `/backup-policies` | Retention, schedule, destination |
| GET/POST | `/backups` | List/request backup |
| GET | `/backups/{id}` | Integrity, scope, encryption, location |
| POST | `/backups/{id}/verify` | Integrity test |
| POST | `/restores/preview` | Scope/conflict/impact plan |
| POST | `/restores` | Approved restore operation |
| GET | `/restore-points` | Available recovery points |

Backups must define consistency level, encryption, key custody, retention, restore scope and recovery objectives. Restore MUST identify overwritten resources and provide pre-restore safety measures. Offsite backup support is part of the platform; local-only backups are not a complete disaster recovery strategy.

**Consolidated requirement:** Additional backup rules:

- **Backups are batch operations.** `POST /backups` creates an operation with `workload_category = "batch"`. It appears in the batch queue.
- **Backup verification is a first-class endpoint.** `POST /backups/{id}/verify` runs a restore test to a disposable target and records the result in `storage.backup_verifications`. A backup without a recent verification is `unverified`, not `healthy`.
- **On Nano, backups are remote-only.** `POST /backups` with `destination = "local"` returns `409 CAPABILITY_UNAVAILABLE` with reason "Local backups require more disk. Use remote backup."
- **`GET /backups/{id}` includes `last_verified_restore_at`.** This is what the UI shows as "last verified restore: 3 days ago."

### Source contract API-20: Security, firewall, SSH, access and incidents

| Method | Endpoint family | Operations |
|---|---|---|
| GET | `/servers/{id}/firewall` | Effective firewall policy |
| POST | `/servers/{id}/firewall/preview` | Lockout and conflict preflight |
| POST | `/servers/{id}/firewall/apply` | Approved policy change |
| GET/POST | `/ssh-keys` | Public-key registry |
| GET/POST | `/security-policies` | Organization/provider policies |
| GET | `/security/events` | Authorized security timeline |
| GET | `/audit-logs` | Immutable/redacted audit records |
| POST | `/security/incident-reports` | Incident tracking |

Firewall updates MUST support temporary rollback/recovery if management connectivity fails; avoid removing the only viable administration path. SSH private keys are never uploaded as general account credentials. Audit logs must be append-only in application semantics and exported to protected storage as policy requires.

**Consolidated requirement:** Additional security rules:

- **The audit chain is queryable.** `GET /audit-logs/{id}` includes `{previous_hash, record_hash}`. The chain is verifiable: `POST /audit-logs/verify` returns `{chain_intact: true, last_verified_at: "..."}`. A broken chain returns `{chain_intact: false, broken_at: "..."}` and triggers a critical alert.
- **Container safety validation is at the API layer.** `POST /service-instances` with a template that requests `--privileged` or host PID is rejected with `CGROUP_CONFIG_INVALID`. The agent enforces this too; defense in depth.
- **Firewall lockout prevention is enforced in the plan.** A plan that would remove the current SSH source IP is rejected unless the actor has verified an alternate access path. This is a policy rule, not a UI warning.

### Source contract API-29: Threat model and integration security controls

- **Cross-tenant IDOR:** verify ownership and capability on every lookup, mutation and event stream.
- **Agent compromise:** mTLS, scoped identity, revocation, command allowlist, host isolation and tamper-evident audit.
- **SSRF:** block private/metadata networks for user-provided webhook URLs and remote fetches unless explicitly allowed; DNS rebinding protections.
- **Webhook spoofing/replay:** signed payloads, bounded timestamp skew and deduplication.
- **Supply-chain risk:** signed releases, SBOM, dependency scanning, pinned checksums, staged rollouts.
- **Secrets leakage:** encryption, scoped access, log masking, write-only APIs, rotation.
- **Abusive hosting:** outbound email limits, resource quotas, abuse reports, suspension workflows and evidence retention policy.
- **Confused deputy:** AI, plugins, resellers and service accounts cannot exceed caller/delegation scope.
- **Races and duplicate jobs:** transactional reservations, idempotency, worker leases and fencing tokens.
- **Dangerous changes:** plan hashes, approvals, backup checks, rollback/recovery paths.

**Consolidated requirement:** Additional threats:

- **M-Pesa callback forgery:** verify Daraja signature, event ID, merchant identity, amount/currency match. Tested with negative signatures and replay.
- **Payment state transition bypass:** a `pending` payment cannot jump to `confirmed` without a callback. A `confirmed` payment cannot be reversed without an explicit reversal event. Enforced by state machine.
- **Idempotency bypass via payment retry:** the idempotency key is composed of actor + invoice + amount + provider. Two different actors cannot share the same key. Duplicate submissions return the existing payment.
- **cgroup misconfiguration:** the agent validates all cgroup values; rejects `cpuset`, unknown controllers, out-of-range values.
- **Container escape:** the agent rejects `--privileged`, host PID/network, dangerous capabilities. Tested with malicious Compose.
- **Capability spoofing:** capability records are loaded from signed data files. A tampered file is rejected at install time. The installer signature is verified.
- **Persona matrix tampering:** the persona matrix is data, but it is signed and verified on load. The API rejects unsigned data.

### Source contract AGENT-28: Upgrade, rollback and disaster recovery

**AGT-146:** Agent updates MUST use signed, verified artifacts and version compatibility gates. **AGT-147:** Upgrade MUST preserve last-known-good binary/config and provide health-gated rollback where safe. **AGT-148:** Control-plane and agent protocol MUST have a documented version compatibility window and graceful feature negotiation. **AGT-149:** An unavailable control plane MUST not interrupt already-running hosted sites solely because management connectivity is lost. **AGT-150:** Agent recovery MUST distinguish control-plane metadata loss, local journal loss, service corruption, credential compromise, and full server loss. **AGT-151:** Disaster-recovery drills MUST validate restoration of metadata, managed configs, customer data, credentials, and authorization boundaries.

**Consolidated requirement:** Additional rules:

- **The agent is a single static Go binary.** Upgrade is: stop agent, replace binary, start agent, verify. Rollback is: stop agent, replace with previous binary, start. The journal survives.
- **Panel upgrade is: stop API, run migrations, start API, verify.** Agent upgrade is: stop agent, replace binary, start agent. On Nano, this happens in seconds.
- **`ServiceTemplate` export/import is the migration path.** A Nextcloud instance is exported as an OCI artifact + manifest. On the new VPS, `POST /service-instances/import` recreates it.
- **The audit chain is backed up separately.** Even if the primary DB is lost, the chain can be reconstructed.

### Source contract AGENT-29: Observability, audit and support evidence

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

**Consolidated requirement:** Additional rules:

- **The audit chain is hash-chained.** Each `audit_events` row has `previous_hash` and `record_hash`. The agent writes to the chain on every operation. The control plane verifies hourly. A broken chain is a critical incident.
- **cgroup metrics are audited.** Every `cgroup.update` operation records the before/after values in the audit chain.
- **Batch yield events are audited.** Every `batch.yielded` and `batch.resumed` is recorded. This is how you prove the batch job yielded to site traffic, not failed.

### Source contract AGENT-30: Threat model and defensive controls

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

**Consolidated requirement:** Additional threats:

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

## 17. Completion criteria and next-volume handoff

Volume 08 is ready for engineering review when each mandatory `SEC8-*` clause is linked to an owner, test case, implementation component and release tier; the inherited `NP-*`, `ARCH-*`, `AGT-*` contracts are reconciled with their owning volumes; key custody, audit canonicalization, DR targets and legal gates have explicit decisions; and test evidence exists for every marketed security or recovery claim. **Document completeness is not implementation readiness.**

**Volume 09 dependency:** commerce, resellers and provider operations must adopt the authentication, support-access, payment authenticity, ledger audit, revocation and disaster-recovery contracts defined here, without moving payment secrets into the privileged agent.

