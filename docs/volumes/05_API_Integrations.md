# NexusPanel — Master Product & Engineering Specification
## Volume 05 — Complete API, Integrations & Automation

**Version:** 1.0 consolidated | **Status:** Proposed engineering contract; review required | **Date:** 2026-11-06  
**Authority:** Volume 01 (Master PRD); Volumes 02–04; authoritative 29-area scope.  
**Source basis:** Uploaded Stage 6 API & Integration Specification v2.0, Master Blueprint v2.0, Stages 2–7 v2.0 and completed master volumes 01–04.  
**Operational boundary:** Documentation only; no services deployed, APIs published, credentials configured or VPS changes made.

### 0. Document contract and consolidated normative decisions

This volume is the single API/integration specification for NexusPanel. It preserves the substantive endpoint inventories, contracts, workflows, and verification coverage from Stage 6 v2.0 while resolving duplicated or conflicting architecture statements. All endpoint paths and payload examples are **proposed contracts**, not an existing production API. **MUST** denotes a release obligation for the corresponding feature; **VERIFY** denotes an assumption requiring proof; **APPROVAL** denotes a policy requiring explicit product/security signoff. Where this volume and Volume 01 differ on product behavior, Volume 01 governs; an engineering conflict must be documented and corrected, not silently ignored.

#### 0.1 Binding principles

1. **One governed public API:** Browser SPA, CLI, SDK, plugins and AI invoke the same tenant-scoped authorization and entitlement policy. The privileged agent is never exposed as a public generic shell API.
2. **Two independent gates:** Authentication/authorization is evaluated without leaking resource existence; resource feasibility/capability is separately evaluated for authorized requests. An unavailable capability must not disclose private resource information to an unauthorized caller.
3. **Nano is a deployment profile, not a second product:** API shapes and authorization rules stay consistent across SQLite and PostgreSQL deployments; unavailable local services return reasoned capability outcomes and may offer a remote alternative.
4. **Async state is explicit:** Infrastructure mutations and payment initiation produce durable operation/payment resources. HTTP 202 means accepted for processing, **not** completed. Read-only and genuinely synchronous endpoints may use 200/201/204.
5. **Shared ServiceTemplate envelope, specialized lifecycles:** Generic template/instance provisioning coexists with typed stateful-service APIs (database users, backups, mail delivery, storage object policy, DNS zones). A template cannot override tenant, security or lifecycle policies.
6. **Local and remote agents differ by transport, not authority:** Standalone v1 uses a protected Unix socket with peer credentials and command binding. Remote agents require mutually authenticated encrypted transport, enrollment and rotation. A local Unix socket does not automatically need a remote-style signature, but every command needs an authenticated origin, bounded authority and replay defense.
7. **Router contract:** The SPA route `/` composes the landing from authenticated capabilities, permissions and the approved jobs-to-be-done matrix. A server may serve the static SPA for `/`; the API supplies the authoritative data, not an unauthenticated persona decision. UI persona is not an authorization identity.
8. **Payments are provider-specific at the trust boundary:** M-Pesa Daraja is the reference integration and Stripe follows the same business-state model. Do not assume Daraja callbacks have a Stripe-like cryptographic signature. Use provider-supported authenticity checks, reference/amount/currency correlation, authoritative status queries where available, idempotent reconciliation, and manual review on ambiguity.
9. **Resource telemetry is measured:** `allocated` is a limit, not consumption; `current`, `typical` and `peak` include sample time and measurement provenance. Missing data is `unknown`, not zero.
10. **Workload classes are governed:** Serving, batch and infrastructure classes determine scheduling and resource policy. Changing a class is a separate migration/approval operation, not an ordinary patch or hidden rewrite.
11. **Secrets and personal data:** Secret write inputs are write-only, encrypted at rest, never echoed in audit/events/logs. Provider callback bodies are minimized and retained by policy.
12. **No fake implementation evidence:** OpenAPI, generated SDKs, signature verification, SLOs and conformance tests remain release gates until actually implemented and tested.

#### 0.2 Contract priority and response vocabulary

| Situation | Public response | Required server behavior |
|---|---|---|
| Caller not authenticated | `401 AUTHENTICATION_REQUIRED` | No private resource or capability details |
| Caller lacks scope | `403 FORBIDDEN` or privacy-preserving `404` | Audit denial as policy permits |
| Authorized request but infeasible on target node | `409 CAPABILITY_UNAVAILABLE` | Return machine-readable reason, missing prerequisites, safe alternatives and resource delta where trustworthy |
| Mutating request accepted | `202` with operation URL | Persist operation and idempotency binding before dispatch |
| Precondition/desired revision mismatch | `412 PRECONDITION_FAILED` or `409 REVISION_CONFLICT` | Do not dispatch host action |
| Resource is being reconciled | `423 RESOURCE_LOCKED` or `409 OPERATION_CONFLICT` | Link current operation, do not duplicate |
| External payment not yet settled | `202` on initiation; `200` on GET with `pending` | No entitlement grant before verified settlement |
| Batch work yielded by contention | `200` on operation GET with `queued`/`yielded` | Not an HTTP error or fake failure |
| Agent disconnected | Operation `waiting_agent`/`unknown` | Reconcile actual host state before retry |

#### 0.3 Cross-volume invariants

- Volume 01 defines requirement identifiers and product release scope; Volume 02 defines the UX contract; Volume 03 defines trust and runtime topology; Volume 04 defines persistence and transaction semantics.
- All 29 areas remain in product scope, but only features scheduled for the current release and supported by the detected profile are exposed as actionable.
- API routes are versioned under `/api/v1`; examples in subsequent tables are relative to that prefix unless specified otherwise.
- API and worker operations are idempotent with bounded retention; replay and out-of-order events must not double-post money, overwrite newer revisions, or bypass approvals.
- The final API specification must be emitted as machine-readable OpenAPI and conformance suites before implementation signoff; this Markdown document is the human-readable source contract.

---

## 1. API design principles

1. **One governed API:** Web, mobile, CLI, SDK, automation, plugins and AI use the same authorization and resource APIs. Privileged agent actions are never public endpoints.
2. **Declarative intent:** Clients request a desired outcome; orchestration plans and reconciles host state.
3. **Asynchronous by default for infrastructure changes:** Long-running or potentially disruptive operations return an operation resource rather than pretending completion.
4. **Resource ownership:** Every resource is scoped to an organization, hosting account, project, server or other explicit ownership boundary.
5. **Least privilege:** Scopes and policies are evaluated on each resource, including reseller delegation and support impersonation.
6. **Safe retries:** Mutations are idempotent where feasible and expose explicit conflict semantics.
7. **Observable outcomes:** Operations, status, steps, logs, audit records and reconciliation results can be inspected.
8. **Compatibility:** Public APIs are versioned and support documented deprecation windows.
9. **No secrets in URLs or logs:** Secret values use controlled write-only inputs and masked read representations.
10. **Separation of control plane and data plane:** Customer database contents, email bodies and website files do not belong in management metadata.

Additional principles:

11. **Async-first payments.** Every asynchronous payment initiation endpoint returns `202 Accepted` with a `pending` operation/payment resource. The state machine is identical for synchronous and asynchronous gateways. This is the only way to support both M-Pesa and Stripe without two code paths.
12. **Capability gate is an independent feasibility check, not a substitute for authorization.** Before policy evaluation, the capability record is checked. A workload that cannot run on this server returns `CAPABILITY_UNAVAILABLE` only after non-disclosing access checks; unauthorized callers receive a non-leaking authorization response.
13. **Workload category is a first-class API concept.** Every operation, resource, and background job declares its category (`serving` | `batch` | `infrastructure`). The category determines cgroup configuration and queue priority. Changes require a separately authorized reclassification workflow.
14. **ServiceTemplate is the universal service API.** Service provisioning uses the common instance envelope; specialized database, object, and mail lifecycle endpoints remain available when their distinct invariants require them. There is `POST /service-instances` with a `template_id`. New services are data, not endpoints.
15. **The router is a client-side concern.** The SPA route `/` resolves the default landing using the persona matrix data (fetched from `GET /personas/{persona}/surface`). The API does not have a "router" endpoint. The API has the capability set, the permission set, and the persona matrix — the SPA composes the landing.
16. **Allocation is queryable, not inferred.** Every resource exposes `GET /resources/{id}/usage` with current/typical/peak. Clients must not derive "am I close to my limits?" from allocation alone.
17. **Batch is visible.** Every batch operation is visible in `GET /operations?workload_category=batch` with queue state. Clients can show "queued, will run when server is idle" without a special endpoint.
18. **Payment callbacks are operations.** A Daraja or Stripe callback creates an `operations.operations` row of kind `payment.resolve`. It is idempotent, audited, and retryable. There is no "silent side effect" path.

## 2. API surfaces and trust boundaries

| Surface | Consumer | Transport | Trust boundary | Notes |
|---|---|---|---|---|
| Public REST API | Browser, CLI, SDK, automation | HTTPS JSON | Internet → control plane | `/api/v1` |
| Real-time events | Dashboard, CLI | Authenticated SSE (initial) | Internet → control plane | Scoped, resumable streams |
| Agent RPC | Control plane ↔ Go agents | **Unix socket (v1) / mTLS (v2+)** | Control plane → privileged host | Typed operations only |
| Internal event bus | Control-plane modules/workers | Durable DB outbox initially | Service boundaries | At-least-once delivery |
| Outbound webhooks | NexusPanel → external endpoints | HTTPS signed JSON | Platform → customer endpoint | Retry and replay protection |
| Inbound provider callbacks | Billing, DNS, domain, identity integrations | HTTPS verified signatures | Third party → integration boundary | Never trust client-supplied state |
| Plugin API | Approved extensions | Capability-scoped SDK | Plugin sandbox → platform | No implicit root privileges |
| **M-Pesa Daraja callback** *(new)* | **Daraja → control plane** | **HTTPS signed callback** | **Third party → payment boundary** | **Idempotent, signature-verified, reconciliation queue** |
| **Stripe webhook** *(new)* | **Stripe → control plane** | **HTTPS signed webhook** | **Third party → payment boundary** | **Idempotent, event ID dedup** |
| **ServiceTemplate artifact fetch** *(new)* | **Control plane → OCI registry** | **HTTPS signed OCI** | **Control plane → external registry** | **Digest-pinned, signature-verified** |

**Topology:** A single VPS may host UI, API, worker, metadata database and local agent. A provider deployment may place control plane and agents on different hosts. Their public contracts must remain the same.

in v1. Not mTLS. Not network. The agent is on the same machine as the control plane (for Nano and Starter) or on a managed node reached via the control plane's worker (Standard and Provider). For multi-node (v2+), mTLS is added. The typed operation catalog is identical; only the transport changes.

## 3. URI conventions and versioning

- Canonical prefix: `/api/v1`.
- Collections use plural kebab-case nouns: `/organizations/{organizationId}/hosting-accounts`.
- Resource IDs are opaque stable identifiers; never assume sortable numeric IDs or encode tenant ownership into IDs.
- Query parameters: `limit`, `cursor`, `sort`, `filter[...]`, `include`, `fields[...]` with an allowlist per endpoint.
- Use `GET` for retrieval, `POST` for creation/actions, `PATCH` for partial updates, `DELETE` for deletion requests. Do not overload GET with mutations.
- Changes that break clients require a new major API version; additive fields and endpoints may remain within v1.
- API schema, docs and generated SDKs MUST be versioned together.
- Clients MUST ignore unknown response fields but MUST NOT ignore unrecognized operation states.

Additional conventions:

- **The API is Go-idiomatic.** JSON field names are `snake_case` (e.g., `created_at`, `workload_category`). Go structs use `json:"created_at"` tags. The OpenAPI is generated from the Go types.
- **All IDs are UUIDv7 strings.** Sortable by time, opaque to clients, safe for URLs.
- **All timestamps are RFC 3339 UTC** (`2026-11-06T15:00:00Z`). Clients convert to local timezone for display.
- **Money is integer minor units** with an explicit `currency` field. `{"amount_minor": 15000, "currency": "KES"}` = 150.00 KES.
- **Capability and persona-matrix data is versioned.** `GET /capabilities/data-version` returns the current version and digest. Clients that cache the matrix should check for updates.

## 4. Authentication, sessions and machine credentials

### 4.1 Interactive identity

Support secure web sessions with HttpOnly, Secure, SameSite cookies, CSRF protection, session rotation and explicit logout/revocation. Require MFA for privileged actions based on policy. Support recovery codes and step-up authentication for sensitive operations. Do not store bearer access tokens in browser localStorage.

### 4.2 API access

Personal and service tokens MUST be hashed at rest, scoped to named capabilities and organization/resource constraints, have expiration/revocation, and show last-used metadata. A token secret is returned only once. Provider and reseller tokens MUST not exceed their owner's delegated privileges. OAuth/OIDC integrations are optional adapters with explicit mapping of external identity to internal principals.

### 4.3 Authentication endpoints

| Method | Path | Purpose | Guardrails |
|---|---|---|---|
| POST | `/auth/login` | Interactive login | Rate limiting, anti-enumeration |
| POST | `/auth/logout` | End session | CSRF protection |
| POST | `/auth/mfa/challenge` | Begin MFA | Short-lived challenge |
| POST | `/auth/mfa/verify` | Complete MFA | Attempt limits |
| GET | `/me` | Current principal and effective context | No sensitive credentials |
| GET | `/me/sessions` | List sessions | Device metadata redacted |
| DELETE | `/me/sessions/{id}` | Revoke session | Audit event |
| GET/POST | `/api-tokens` | List/create tokens | Write-only secret |
| DELETE | `/api-tokens/{id}` | Revoke token | Audit event |

**Open decision:** precise cookie/session implementation, OIDC strategy and MFA factors. These must be decided before auth code is finalized.

Additional endpoints:

| Method | Path | Purpose | Guardrails |
|---|---|---|---|
| GET | `/me/capabilities` | The capabilities available to this actor on this installation | Derived from server capability set × actor's permission set; no secret exposure |
| GET | `/me/persona-matrix` | The persona × capability matrix for the SPA's router | Signed data; clients verify the digest |
| GET | `/me/default-landing` | The resolved default landing (`router` / `portfolio` / `fleet`) | Derived from `me.capabilities` and `tenancy.organizations.primary_persona` |
| POST | `/me/step-up` | Step-up authentication for sensitive actions | Short-lived token; binds to action + resource |

The `GET /me/default-landing` endpoint is what the SPA calls on first load to decide whether to render the router or the portfolio. The router itself is a client-side composition.

## 5. Authorization and tenancy

Every request MUST resolve a principal, acting organization, resource ownership chain, explicit capability, subscription/entitlement and any required step-up/approval policy. Resource existence must not be leaked to unauthorized tenants. Provider support access requires justification, time bounds, explicit audit and appropriate consent/policy. Resellers may manage only delegated customer resources. Role labels (owner, provider, reseller, customer, developer, support, auditor) are UI concepts; authorization is capability and relationship based.

**Representative capabilities:** `server.read`, `server.manage`, `hosting-account.create`, `website.deploy`, `domain.dns.write`, `database.provision`, `email.mailbox.manage`, `backup.restore`, `security.firewall.write`, `billing.read`, `reseller.delegate`, `audit.read`, `agent.enroll`, `approval.decide`.

**Authorization decision record:** principal ID, resource ID, action, effective organization, policy version, decision and reason code; never record secret values.

Additional authorization rules:

- **The capability gate precedes the policy engine.** Sequence: (1) authenticate, (2) resolve server capability record, (3) if capability is unavailable → return `CAPABILITY_UNAVAILABLE` with reason and resource delta, (4) evaluate policy, (5) evaluate quota/entitlement, (6) approve if required, (7) enqueue operation. This means the policy engine never has to reason about "can this server do this?" — that is data.
- **Personas are not permissions.** `GET /me/capabilities` returns what the actor can do (authorization). `GET /me/persona-matrix` returns what the UI should show (presentation). The API never denies a request because of the persona. It denies because of capability or policy.
- **The agent re-evaluates.** Before executing an operation, the agent re-checks: capability still available, resource ownership unchanged, policy still permits. A revoked permission during a queued job is caught at execution time, not just at submission.

**Representative capabilities (v2 additions):**

- `service.create`, `service.read`, `service.manage`, `service.export`, `service.import`
- `capability.read` (view capability set), `capability.rescan` (trigger hardware re-detection)
- `payment.read`, `payment.reconcile` (admin only)
- `template.deploy`, `template.export`
- `workload.category.read`, `workload.usage.read`

**Authorization decision record (v2 additions):**

- `capability_decision`: `{capability_key, state, reason, source_digest}` — recorded alongside the policy decision
- `workload_category`: recorded on every operation

## 6. Standard response and error contracts

### 6.1 Successful resource response

```json
{
  "data": {
    "id": "web_01EXAMPLE",
    "type": "website",
    "organization_id": "org_01EXAMPLE",
    "name": "example-site",
    "status": "active",
    "version": 7,
    "created_at": "2026-11-06T15:00:00Z"
  },
  "meta": {"request_id": "req_01EXAMPLE"}
}
```

### 6.2 Error response

Use `application/problem+json` (RFC 9457 style), with stable machine-readable codes:

```json
{
  "type": "https://docs.nexuspanel.example/problems/quota-exceeded",
  "title": "Quota exceeded",
  "status": 409,
  "code": "QUOTA_EXCEEDED",
  "detail": "The hosting account has reached its website limit.",
  "request_id": "req_01EXAMPLE",
  "errors": [{"field": "hosting_account_id", "code": "LIMIT_REACHED"}]
}
```

The documentation URL above is illustrative, not an existing website. Never expose stack traces, internal filesystem paths, SQL queries or secrets in client errors.

| HTTP | Typical meaning | Examples |
|---|---|---|
| 200 | Synchronous success | Read, preview |
| 201 | Resource created | Token, project |
| 202 | Operation accepted | Provision, deploy, restore, **payment initiate** |
| 204 | Success without body | Revoke, remove |
| 400 | Invalid request syntax | Malformed JSON |
| 401 | Authentication required | Missing/expired credentials |
| 403 | Forbidden | Policy denies action |
| 404 | Not found or intentionally concealed | Cross-tenant resource |
| 409 | Conflict | Version, quota, desired state |
| 412 | Failed precondition | ETag mismatch |
| 422 | Semantically invalid | Unsupported runtime version |
| 423 | Locked | Resource under conflicting operation |
| 429 | Rate limited | Token/tenant/IP thresholds |
| 500/502/503/504 | Platform/dependency failure | Include safe retry guidance |

Additional error codes:

| HTTP | Code | Meaning | Example |
|---|---|---|---|
| 409 | `CAPABILITY_UNAVAILABLE` | The requested workload cannot run on this server | Nano VPS attempting to create a Node app |
| 409 | `WORKLOAD_CATEGORY_MISMATCH` | The operation's category doesn't match the resource's | Attempting to run a batch operation as serving |
| 409 | `PAYMENT_PENDING` | A payment for this invoice is already pending | Second attempt while M-Pesa STK push is awaiting |
| 409 | `PAYMENT_ALREADY_RESOLVED` | A payment was already confirmed/failed/reversed | Duplicate callback with different outcome |
| 422 | `CAPABILITY_DEPENDENCY_MISSING` | A prerequisite (port 25, Docker, public IP) is missing | Email on a VPS with port 25 blocked |
| 422 | `TEMPLATE_SIGNATURE_INVALID` | The template's signed image reference is invalid | Malicious or corrupted template |
| 422 | `CGROUP_CONFIG_INVALID` | The cgroup config violates the workload category's constraints | Attempt to set `cpuset` or RAM ceiling above physical × 1.2 |
| 423 | `BATCH_QUEUE_FULL` | The batch queue has reached its concurrency limit | Too many concurrent transforms |
| 429 | `BATCH_YIELDED` | The operation was yielded under contention | A build waiting for site traffic to subside |
| 503 | `AGENT_UNAVAILABLE` | The local agent is not reachable | Agent crashed; operation cannot be dispatched |

## 7. Pagination, concurrency, rate limits and idempotency

- Use cursor pagination with deterministic order and bounded `limit` (default 25, maximum 100 unless documented otherwise).
- List responses include `data[]` and `page.next_cursor` (nullable); cursors are opaque, signed/validated and scoped to filters/tenant.
- Mutable resources SHOULD expose `ETag`; updates require `If-Match` when lost updates would be dangerous.
- For non-idempotent POST mutations, clients MUST supply `Idempotency-Key` (UUID or high-entropy token). Server stores key + principal/tenant + request digest + response/operation for a documented retention window. Reuse with different payload yields 409.
- Rate limits apply by IP, principal, organization and endpoint sensitivity; include `Retry-After` on 429.
- Resource names and unique constraints must be validated transactionally; a successful preflight is not a reservation.

Additional rules:

- **Payment idempotency is scoped to actor + invoice + amount + provider.** `Idempotency-Key` alone is not enough for payments; two different actors could submit the same key. The server composes `(actor_id, invoice_id, amount_minor, provider, idempotency_key)` as the deduplication key.
- **Payment callbacks are deduplicated by provider reference.** Daraja `CheckoutRequestID` and Stripe `evt_...` are the dedup keys. A duplicate callback returns `200 OK` with the existing state, not an error.
- **Batch operations have a separate rate limit.** Serving operations are limited by actor; batch operations are limited by queue depth. A user can enqueue many batch jobs (they queue), but only a bounded number run concurrently.
- **Capability data version check.** Clients that cache `/capabilities` should send `If-None-Match` with the digest. The server returns `304 Not Modified` if unchanged.

## 8. Asynchronous operations: mandatory lifecycle

Infrastructure mutations generally return `202 Accepted`, `Location: /api/v1/operations/{id}` and an operation object. An operation references the requested resource, actor, desired-state revision, policy/approval decision, execution plan, steps, timestamps and final result.

**State machine:** `pending` → `validating` → (`awaiting_approval` | `queued`) → `running` → (`verifying` | `compensating`) → (`succeeded` | `failed` | `partially_succeeded` | `cancelled`). Transitions must be explicitly validated; retries create attempts, not duplicate resources. `partially_succeeded` requires actionable remediation details.

**Endpoints:**

| Method | Path | Purpose |
|---|---|---|
| GET | `/operations/{id}` | Status, steps, outcome, safe error |
| GET | `/operations/{id}/events` | Resumable SSE status stream |
| POST | `/operations/{id}/cancel` | Request cancellation if safe |
| POST | `/operations/{id}/retry` | Retry eligible failed operation |
| GET | `/operations?resource_id=...` | Filtered operation history |
| GET | `/operations/{id}/logs` | Redacted execution diagnostics |

**Invariants:** Work must survive API restarts; workers claim tasks using leases/fencing; every external effect is tracked; cancellation is best-effort and cannot promise reversal of already committed changes; compensations are domain-specific; postcondition verification is required before reporting success.

### 8.1 Operation response example

```json
{
  "data": {
    "id": "op_01EXAMPLE",
    "kind": "website.provision",
    "status": "queued",
    "resource": {"type": "website", "id": "web_01EXAMPLE"},
    "requested_by": "usr_01EXAMPLE",
    "workload_category": "serving",
    "progress": {"completed_steps": 0, "total_steps": 5},
    "created_at": "2026-11-06T15:00:00Z"
  }
}
```

Additional operation rules:

- **Every operation has `workload_category`.** It is set at creation and is immutable. The `GET /operations` endpoint supports `filter[workload_category]=batch` for the batch queue UI.
- **Batch operations expose queue position.** `GET /operations/{id}` for a batch job returns `{queue_position: 3, estimated_start: "...", yielding_reason: "serving_workload_contention"}`.
- **cgroup changes are operations.** `POST /resources/{id}/limits` creates an `operation` of kind `resource.limit.update`. It verifies the new cgroup value is applied before succeeding. The response includes `{before: {...}, after: {...}}`.
- **Payment resolution is an operation.** A callback creates an operation of kind `payment.resolve`. The operation's state machine mirrors the payment's state machine. This means every payment state change is auditable through the operations endpoint.

### 8.2 Batch operation endpoints

| Method | Path | Purpose |
|---|---|---|
| GET | `/operations?workload_category=batch` | List all batch operations with queue state |
| GET | `/batch-queue` | Current batch queue with positions and estimated start times |
| POST | `/operations/{id}/yield` | Request the operation to yield (returns to queue) |
| POST | `/operations/{id}/prioritize` | Move a batch operation earlier in the queue (admin only, bounded) |

Batch operations cannot be promoted to serving. `prioritize` reorders within the batch queue; it does not change the category.

## 9. Preflight, plan, approval, apply and verify

For impactful changes (firewall, DNS delegation, destructive database actions, restore, account suspension, bulk operations), the UI and API MUST support an explicit workflow:

1. **Preflight:** Validate permissions, ownership, service availability, quotas, conflicts, backups and dependencies.
2. **Plan:** Produce human-readable diff, exact targets, impact/risk, estimated downtime class and rollback/compensation strategy.
3. **Approve:** Bind approval to plan hash, actor, scope, expiry and step-up auth; policy may require a second approver.
4. **Apply:** Execute only the approved plan version with an idempotency key.
5. **Verify:** Test observable postconditions; report any unresolved differences.

Proposed endpoints: `POST /plans`, `GET /plans/{id}`, `POST /plans/{id}/approve`, `POST /plans/{id}/apply`. Plans expire on relevant state changes; approval cannot be reused for modified payloads. Read-only previews may use module-specific `/preview` endpoints.

Additional preflight content:

- **Capability preflight is automatic.** Before plan generation, the capability gate runs. The plan includes `{capability_decision: {state, reason}}`. If the capability is unavailable, no plan is generated; the request returns `CAPABILITY_UNAVAILABLE` directly.
- **Resource fit preflight is included in the plan.** The plan includes `{resource_fit: {ram_floor, ram_typical, ram_ceiling, cpu_weight, cpu_max, current_usage, delta_after}}`. This is what makes "will this fit?" answerable before the operation runs.
- **Batch plans include queue estimate.** A batch operation's plan includes `{estimated_queue_time, yielding_policy}` so the user knows "this build will queue behind 2 running builds and yield to site traffic."

## 10. Identity, organizations, teams, roles and invitations

| Method | Endpoint family | Operations |
|---|---|---|
| GET/POST | `/organizations` | List/create organizations |
| GET/PATCH | `/organizations/{id}` | Read/update organization |
| GET/POST | `/organizations/{id}/members` | List/add membership |
| PATCH/DELETE | `/organizations/{id}/members/{id}` | Change/revoke membership |
| GET/POST | `/organizations/{id}/invitations` | Invite, list |
| POST | `/invitations/{token}/accept` | Accept bounded invitation |
| GET | `/organizations/{id}/effective-permissions` | UI permission discovery |
| GET/POST | `/organizations/{id}/teams` | Teams |
| GET/POST | `/organizations/{id}/service-accounts` | Automation identities |

Invitations expire, are single-use, and never grant privileges above the inviter's delegable scope. Membership changes revoke or invalidate derived sessions/tokens as policy requires.

Additional identity endpoints:

| Method | Endpoint | Purpose |
|---|---|---|
| GET/PATCH | `/organizations/{id}/persona` | The organization's default persona (for router landing) |
| PATCH | `/organizations/{id}/members/{id}/persona` | Per-membership persona override |
| GET/POST | `/organizations/{id}/branding` | White-label branding (v2+) |

The persona is a UX hint, not a permission. Setting `primary_persona = "app_developer"` changes the router's default landing for new members; it does not grant or deny any capability.

## 11. Provider, reseller, customer and hosting-account APIs

| Method | Endpoint family | Operations |
|---|---|---|
| GET/POST | `/providers/{id}/resellers` | Delegate reseller relationships |
| GET/POST | `/providers/{id}/hosting-packages` | Manage package templates |
| GET/POST | `/resellers/{id}/customers` | Manage delegated customers |
| GET/POST | `/hosting-accounts` | Create/list hosting accounts |
| GET/PATCH | `/hosting-accounts/{id}` | Read/update account |
| POST | `/hosting-accounts/{id}/suspend` | Suspend with plan/approval |
| POST | `/hosting-accounts/{id}/unsuspend` | Restore service |
| GET | `/hosting-accounts/{id}/quotas` | Limits, reservations, usage |
| POST | `/hosting-accounts/{id}/transfer` | Governed ownership transfer |

Package changes must define whether existing accounts inherit changes. Quotas must be enforced at transaction time and during reconciliation; oversubscription and grace periods are explicit provider policies. Provider/reseller branding must not alter authorization boundaries.

No substantive change. The commercial endpoints are defined but return `403 Forbidden` on Nano installations where the commercial surface is not enabled. The data model is complete; the API surface is gated.

## 12. Servers, clusters, agents and capacity

| Method | Endpoint family | Operations |
|---|---|---|
| GET/POST | `/servers` | Register/list managed nodes |
| GET/PATCH | `/servers/{id}` | Metadata and scheduling labels |
| GET | `/servers/{id}/inventory` | OS, services, runtime inventory |
| GET | `/servers/{id}/health` | Heartbeat, connectivity, capacity |
| GET/POST | `/server-pools` | Group eligible servers |
| GET | `/server-pools/{id}/capacity` | Aggregated placement capacity |
| POST | `/agents/enrollment-tokens` | Short-lived one-time enrollment |
| POST | `/servers/{id}/drain` | Stop new placements and migrate where supported |
| POST | `/servers/{id}/maintenance` | Schedule maintenance state |

Agents authenticate using mTLS with rotation and revocation. Server deletion requires confirming there are no managed workloads or an explicit governed evacuation plan. Scheduling must consider labels, capability compatibility, resource reservations, locality, isolation and health.

Significant addition — the capability endpoints:

| Method | Endpoint family | Operations |
|---|---|---|
| GET | `/servers/{id}/capabilities` | The server's capability set (from `inventory.capabilities`) |
| POST | `/servers/{id}/capabilities/rescan` | Trigger hardware re-detection (after VPS upgrade) |
| GET | `/servers/{id}/capabilities/data-version` | The version/digest of the capability data file |
| GET | `/servers/{id}/capabilities/{key}` | Single capability detail (floor/typical/ceiling, category, prerequisites, state, reason) |
| GET | `/servers/{id}/usage` | Server-level allocation/typical/peak/current |
| GET | `/resources/{id}/usage` | Per-workload usage metrics |
| GET | `/servers/{id}/cgroup-metrics` | Raw cgroup metrics (time-series) |
| POST | `/resources/{id}/limits` | Adjust the cgroup ceiling (creates an operation) |

**Agent transport in v1 is a Unix socket.** The enrollment endpoint is simplified: the local agent auto-enrolls on install. For remote agents (v2+), the full enrollment flow with mTLS applies.

## 13. Websites, domains and routing

| Method | Endpoint family | Operations |
|---|---|---|
| GET/POST | `/websites` | List/provision sites |
| GET/PATCH/DELETE | `/websites/{id}` | Inspect/change/delete via operation |
| POST | `/websites/{id}/preview` | Validate web server changes |
| POST | `/websites/{id}/publish` | Apply routing/configuration |
| GET/POST | `/domains` | Register domain resource/ownership |
| POST | `/domains/{id}/verify` | Domain ownership verification |
| GET/POST | `/websites/{id}/bindings` | Hostnames and routing |
| POST | `/websites/{id}/redirects` | Redirect rules |
| GET | `/websites/{id}/traffic` | Authorized metrics |

Enforce domain ownership and uniqueness across tenants; verify wildcard and apex cases separately. Routing changes must be tested before reload, preserve last-known-good configuration, and report TLS readiness independently of DNS readiness.

Additional website endpoints:

| Method | Endpoint | Purpose |
|---|---|---|
| GET | `/websites/{id}/capabilities` | The site's capability subset (PHP? Node? DB? Docker?) |
| GET | `/websites/{id}/batch-jobs` | The site's batch jobs (builds, backups, transforms) |
| POST | `/websites` with `runtime` param | The runtime is validated against the server's capability set |

The website creation endpoint rejects with `CAPABILITY_UNAVAILABLE` if the requested runtime isn't available on the server.

## 14. Applications, runtimes, Git, deployments and environment variables

| Method | Endpoint family | Operations |
|---|---|---|
| GET/POST | `/applications` | List/create applications |
| GET/PATCH | `/applications/{id}` | Settings, runtime, start command |
| GET | `/runtimes/catalog` | OS/server-compatible runtime versions |
| POST | `/applications/{id}/deployments` | Deploy Git ref/upload/artifact |
| GET | `/applications/{id}/deployments` | History and statuses |
| POST | `/deployments/{id}/rollback` | Revert to retained release |
| GET | `/applications/{id}/logs` | Scoped/redacted logs |
| POST | `/applications/{id}/actions/restart` | Restart managed process |
| GET/PUT | `/applications/{id}/environment` | Masked metadata/write-only values |
| POST | `/applications/{id}/environment/parse` | Parse pasted `.env` text for preview |
| POST | `/applications/{id}/environment/import` | Apply validated bulk variables |
| GET/POST | `/applications/{id}/git-connections` | Authorized repository connections |

Deployment input MUST specify source, revision, build strategy, runtime and target environment. Avoid allowing arbitrary root-level startup commands. `.env` parsing must preserve quotes, escaped newlines and comments according to a documented parser dialect; reject duplicate or invalid keys; never return secret values after persistence. Secret rotations trigger only approved redeploy/restart actions. Logs must mask credentials.

Additional application rules:

- **The runtime is a serving workload; the build is a batch workload.** `POST /applications` creates the runtime (serving). `POST /applications/{id}/deployments` creates a build (batch) followed by a release (serving). The build appears in `GET /operations?workload_category=batch`.
- **The container is the isolation primitive.** On Standard+, the runtime is a container; the build is a separate container with no host access. The API returns `{container_id, build_container_id}`.
- **Deployment rollback is immutable.** A rollback references a specific artifact digest. The original artifact must still exist. If it has been garbage-collected, the rollback fails with `ARTIFACT_EXPIRED`.
- **The `runtime` field is capability-checked.** `POST /applications` with `runtime: "node"` on a Nano VPS returns `CAPABILITY_UNAVAILABLE`. The error includes `{required_ram_mb: 200, available_ram_mb: 120}`.

## 15. Files, code editor, uploads, archives and terminal

File operations are scoped to an explicit hosting account and permitted root, not arbitrary host paths.

| Method | Endpoint | Purpose |
|---|---|---|
| GET | `/hosting-accounts/{id}/files?path=...` | Directory listing |
| GET | `/hosting-accounts/{id}/file-content?path=...` | Read authorized file |
| PUT | `/hosting-accounts/{id}/file-content?path=...` | Atomic write with version precondition |
| POST | `/hosting-accounts/{id}/file-operations` | Move/copy/rename/delete/chmod/archive |
| POST | `/hosting-accounts/{id}/uploads` | Initiate bounded upload |
| POST | `/uploads/{id}/complete` | Validate checksum and finalize |
| POST | `/hosting-accounts/{id}/terminal-sessions` | Governed interactive terminal |

**Critical safeguards:** canonicalize and validate paths on the agent; defend against symlink races, traversal, hardlinks and TOCTOU; enforce tenant filesystem ownership, file-size and extension policies; scan uploads where policy requires; preserve editor conflict versions; require explicit terminal permission and time-limited session with audit. Terminal access is a separate privileged feature, not an implementation shortcut for normal API operations.

Additional file rules:

- **The file manager is container-aware.** The path is resolved inside the site's container volume. `GET /hosting-accounts/{id}/files` returns paths as the user sees them in the container (e.g., `/app/wp-content/uploads`), not the host path. The agent maps between them.
- **Terminal is disabled by default.** `POST /hosting-accounts/{id}/terminal-sessions` returns `403 FEATURE_DISABLED` unless the account owner explicitly enabled it.
- **Terminal is unavailable on Nano.** The endpoint returns `409 CAPABILITY_UNAVAILABLE` with reason "No container to exec into."
- **SFTP is the primary file-transfer path on Nano.** The panel manages the SFTP endpoint (`GET /hosting-accounts/{id}/sftp-credentials`) and per-site credentials.

## 16. Customer databases and database users

| Method | Endpoint family | Operations |
|---|---|---|
| GET/POST | `/database-instances` | Provider-managed database services |
| GET/POST | `/databases` | List/create tenant databases |
| GET/PATCH/DELETE | `/databases/{id}` | Metadata, limits, governed deletion |
| GET/POST | `/databases/{id}/users` | Scoped database principals |
| POST | `/database-users/{id}/rotate-password` | Write-only new credential |
| POST | `/databases/{id}/backups` | Backup request |
| POST | `/databases/{id}/restore` | Governed restore |
| GET | `/databases/{id}/connection-info` | Non-secret host/port metadata |

Never expose a global database superuser credential. Per-tenant grants must be least privilege; database deletion/restore requires an impact preview and backup policy checks. Control-plane PostgreSQL is never treated as a customer database instance.

Significant reframe — databases are `ServiceTemplate` instances:

- **`/databases` is a read-only convenience view.** It returns a filtered list of `/service-instances` where `template.category = "database"`. Creation uses `POST /service-instances` with `template_id = "database.mariadb"` or `"database.postgres"`.
- **Each database is its own container.** The `service_instances` row has `container_id`, `volume_id`, `cgroup_path`. The database-specific extensions live in `data_service.*` tables referencing `service_instances.id`.
- **On Nano, `/databases` returns `[]` with a capability explanation.** The endpoint returns `{data: [], meta: {capability_unavailable: true, reason: "No room for MariaDB/Postgres container"}}`.
- **The database admin UI is a `ServiceTemplate` instance.** `POST /service-instances` with `template_id = "phpmyadmin"` or `"pgadmin"` creates an authenticated admin tool.

## 17. DNS, nameservers, domain delegation and TLS

| Method | Endpoint family | Operations |
|---|---|---|
| GET/POST | `/dns/zones` | List/create managed zones |
| GET/POST | `/dns/zones/{id}/records` | List/add records |
| PATCH/DELETE | `/dns/records/{id}` | Change/remove record |
| POST | `/dns/zones/{id}/preview` | Validate proposed zone change |
| POST | `/dns/zones/{id}/publish` | Versioned zone publication |
| GET | `/domains/{id}/delegation-check` | External registrar NS preflight |
| GET | `/dns/zones/{id}/health` | SOA/NS/secondary propagation checks |
| GET/POST | `/certificates` | Certificates and issuance requests |
| POST | `/certificates/{id}/renew` | Controlled renewal |
| GET | `/certificates/{id}/status` | Expiry and validation status |

Support `external_manual`, future `external_api` and `local_authoritative` DNS modes. Domain ownership verification is separate from registrar delegation and from certificate validation. Do not claim independent authoritative secondary DNS when both instances share one failure domain. DNSSEC keys, glue records and zone serials need explicit ownership and lifecycle rules. Certificate private keys never appear in ordinary API reads.

Scope change — v1 is DNS guidance, not authoritative hosting:

- **`/dns/zones` is deferred to v2+.** On a v1 installation, the endpoint returns `409 CAPABILITY_UNAVAILABLE` with reason "Authoritative DNS hosting is coming in v2. Use DNS guidance."
- **`/domains/{id}/dns-guidance` is the v1 DNS API.** It returns the exact records to add at the registrar, polls propagation via `GET /domains/{id}/dns-guidance/check`, and confirms when live. It is the source of truth for UX-DNS-014.
- **`/certificates` uses HTTP-01 in v1.** DNS-01 and wildcard are deferred to v2+ with authoritative DNS.
- **PTR coordination is a provider integration.** `POST /domains/{id}/ptr-request` creates a ticket-style record; the panel guides the user to the provider's control panel. The panel never claims it set PTR itself.

**Additional v2 endpoint:**

| Method | Endpoint | Purpose |
|---|---|---|
| GET | `/domains/{id}/dns-guidance` | The exact records to add at the registrar, plus instructions |
| POST | `/domains/{id}/dns-guidance/check` | Poll for propagation; returns the observed state |

## 18. Email, mailboxes, aliases, forwarding and deliverability

| Method | Endpoint family | Operations |
|---|---|---|
| GET/POST | `/mail-domains` | Provision mail domain |
| GET/POST | `/mail-domains/{id}/mailboxes` | Create/list mailbox |
| PATCH/DELETE | `/mailboxes/{id}` | Quotas, status, deletion |
| POST | `/mailboxes/{id}/reset-password` | Credential rotation |
| GET/POST | `/mail-domains/{id}/aliases` | Aliases and forwarding |
| GET | `/mail-domains/{id}/deliverability` | SPF/DKIM/DMARC/PTR diagnostics |
| POST | `/mail-domains/{id}/dkim/rotate` | Controlled key rotation |
| GET | `/mail-domains/{id}/queue` | Authorized queue summary |

Do not expose email message bodies through ordinary hosting-admin APIs. Mailbox administration and mailbox content access are separate consent/security domains. Outbound mail capability must account for IP reputation, port restrictions, abuse controls, bounce handling and reverse DNS provider dependencies.

Scope change — v1 is external SMTP:

- **`/mail-domains`, `/mailboxes`, `/mail-domains/{id}/queue` are deferred to v2+.** On a v1 installation, these return `409 CAPABILITY_UNAVAILABLE` with reason "Self-hosted mail is coming in v2. Use external SMTP."
- **`/email/external-smtp` is the v1 path.** `POST /email/external-smtp` with provider (Resend/Postmark/Mailgun/SES), API key (write-only), and from address. `POST /email/external-smtp/test` verifies the credentials. This is used for panel notifications and WordPress outbound mail.
- **`GET /email/recommendation`** returns the honest guidance: "For mailboxes, we recommend Zoho or Google Workspace. Self-hosted mail is coming in v2."
- **On Nano, `/email/*` returns `CAPABILITY_UNAVAILABLE`.** No email of any kind (not even external SMTP) on the smallest profile.

## 19. Backups, restore, snapshots and disaster recovery

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

Additional backup rules:

- **Backups are batch operations.** `POST /backups` creates an operation with `workload_category = "batch"`. It appears in the batch queue.
- **Backup verification is a first-class endpoint.** `POST /backups/{id}/verify` runs a restore test to a disposable target and records the result in `storage.backup_verifications`. A backup without a recent verification is `unverified`, not `healthy`.
- **On Nano, backups are remote-only.** `POST /backups` with `destination = "local"` returns `409 CAPABILITY_UNAVAILABLE` with reason "Local backups require more disk. Use remote backup."
- **`GET /backups/{id}` includes `last_verified_restore_at`.** This is what the UI shows as "last verified restore: 3 days ago."

## 20. Security, firewall, SSH, access and incidents

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

Additional security rules:

- **The audit chain is queryable.** `GET /audit-logs/{id}` includes `{previous_hash, record_hash}`. The chain is verifiable: `POST /audit-logs/verify` returns `{chain_intact: true, last_verified_at: "..."}`. A broken chain returns `{chain_intact: false, broken_at: "..."}` and triggers a critical alert.
- **Container safety validation is at the API layer.** `POST /service-instances` with a template that requests `--privileged` or host PID is rejected with `CGROUP_CONFIG_INVALID`. The agent enforces this too; defense in depth.
- **Firewall lockout prevention is enforced in the plan.** A plan that would remove the current SSH source IP is rejected unless the actor has verified an alternate access path. This is a policy rule, not a UI warning.

## 21. Monitoring, metrics, logs, alerts and usage

| Method | Endpoint family | Operations |
|---|---|---|
| GET | `/metrics/servers/{id}` | CPU, RAM, disk, bandwidth time series |
| GET | `/metrics/resources/{id}` | Tenant-authorized resource metrics |
| GET | `/logs` | Filtered scoped logs |
| GET/POST | `/alert-rules` | Define threshold and service alerts |
| GET | `/alerts` | Alert lifecycle |
| POST | `/alerts/{id}/acknowledge` | Acknowledge alert |
| GET | `/usage` | Metered usage with time windows |
| GET | `/service-health` | Service state and dependencies |

Metrics must distinguish collected vs missing samples, include timestamps/timezones and document aggregation windows. Logs must enforce tenant scoping and redaction. Billing-grade metering requires reproducible immutable usage records and reconciliation with quota counters.

Additional metrics rules:

- **`GET /metrics/resources/{id}` returns the allocation/usage shape.** `{allocated: {ram_max_bytes, cpu_weight, cpu_max_millicores}, current: {ram_bytes, cpu_percent}, typical_7d: {ram_bytes, cpu_percent}, peak_24h: {ram_bytes, cpu_percent}, throttled_seconds_24h, oom_kills_7d}`. This is the source of truth for the UX allocation/usage strip.
- **`GET /metrics/servers/{id}/cgroup` returns raw cgroup data.** For power users and debugging; the aggregated view is `GET /metrics/resources/{id}`.
- **The upgrade recommendation is a computed endpoint.** `GET /recommendations/upgrade` returns auto-generated suggestions based on 7-day peak RAM > 80%, OOM kills, and CPU contention. This is what powers Stage 3 UX-OV-010.

## 22. Commercial APIs: subscriptions, licensing, invoices and payment providers

| Method | Endpoint family | Operations |
|---|---|---|
| GET/POST | `/plans` | Commercial product plans |
| GET/POST | `/subscriptions` | Tenant subscription lifecycle |
| GET | `/invoices` | Scoped invoices |
| GET | `/usage/charges` | Metered charges |
| GET | `/licenses/current` | License entitlements |
| POST | `/billing/checkout-sessions` | Provider-mediated checkout |
| POST | `/integrations/billing/{provider}/callbacks` | Provider-authenticated event ingress |

Never treat a browser redirect as proof of payment. Billing providers are pluggable; verify callback signature, event ID, merchant identity, amount/currency, order linkage and state transitions. Entitlement revocation/suspension must follow explicit grace and data-retention policies. Avoid storing raw card data.

Significant reframe — async-first, M-Pesa reference:

**Payment endpoints:**

| Method | Endpoint | Purpose |
|---|---|---|
| POST | `/payments` | Initiate a payment (returns `202 Accepted` with `pending`) |
| GET | `/payments/{id}` | Payment state, provider reference, reconciliation state |
| GET | `/payments/{id}/events` | Append-only state transition log |
| POST | `/payments/{id}/retry` | Retry a failed/expired payment (new idempotency key) |
| POST | `/payments/{id}/cancel` | Cancel a pending payment (best-effort) |
| GET | `/payments?state=pending` | List pending payments |
| GET | `/payments?reconciliation_state=unmatched` | Admin: reconciliation queue |

**Payment callback endpoints (provider-specific):**

| Method | Endpoint | Purpose |
|---|---|---|
| POST | `/integrations/payments/mpesa/callback` | M-Pesa Daraja callback (provider-specific verification and reconciliation) |
| POST | `/integrations/payments/stripe/webhook` | Stripe webhook (event ID dedup) |

**Reconciliation endpoints (admin only):**

| Method | Endpoint | Purpose |
|---|---|---|
| POST | `/payments/{id}/reconcile` | Manually match an unmatched callback to an invoice |
| POST | `/payments/{id}/reverse` | Manual reversal (with reason) |
| GET | `/payments/reconciliation` | The reconciliation queue |

**M-Pesa config endpoints (admin only):**

| Method | Endpoint | Purpose |
|---|---|---|
| GET/POST | `/integrations/payments/mpesa/config` | Daraja credentials (sandbox + production) |
| POST | `/integrations/payments/mpesa/test` | Test transaction against Daraja sandbox |

**Payment state machine:**

```
initiated → pending → {confirmed | failed | expired | reversed}
                ↓
          manual_review (if callback unmatched)
```

Every transition is recorded in `commercial.payment_events`. `GET /payments/{id}/events` returns the full state history.

**The same payment initiation envelope works for Stripe.** `POST /payments` with `provider: "stripe"` returns `202 Accepted` with `{state: "pending"}`. The Stripe webhook or provider-status reconciliation resolves it when authoritative status is available; latency is not guaranteed. The API contract is identical.

**Licensing:**

- `GET /licenses/current` returns `{status, expires_at, offline_grace_until, features: [...]}`.
- The license is activated via `POST /licenses/activate` with a license key.
- Offline grace: if the license server is unreachable, the panel continues functioning for a configurable, approved grace window. This is critical — a VPS panel that stops working because the license server is down is not commercially viable.

**Commercial capability is deployment-dependent, not inherently disabled by Nano hardware.** Hosted billing may run remotely. If no billing service is configured, return a typed `FEATURE_UNAVAILABLE` response with a supported alternative.

## 23. AI assistant, guided diagnostics and controlled actions

| Method | Endpoint family | Operations |
|---|---|---|
| POST | `/ai/conversations` | Start scoped assistance |
| POST | `/ai/conversations/{id}/messages` | Ask or provide context |
| POST | `/ai/action-proposals` | Propose typed platform action |
| GET | `/ai/action-proposals/{id}` | Diff, risk, permission checks |
| POST | `/ai/action-proposals/{id}/submit` | Route through plan/approval/apply |
| GET | `/ai/diagnostics/{id}` | Read-only evidence-backed diagnosis |

AI cannot gain permissions beyond the initiating principal, issue unrestricted root shell commands or silently apply high-impact changes. Tool calls must use allowlisted schemas, include resource scope and be auditable. Retrieved logs, websites and files are untrusted content; prompt injection must not override system policy. Secret-bearing material must be minimized/redacted before model access.

Additional AI rules:

- **v1 ships read-only diagnostics + proposed-action approval.** `POST /ai/action-proposals` creates a proposal, not an execution. The user must approve via `POST /ai/action-proposals/{id}/submit`, which routes through the standard plan/approval/apply workflow. **Autonomous multi-step remediation is v4+.**
- **The AI has no privileged API.** It uses `POST /ai/action-proposals` for proposals. It does not call `/agents/*` or `/servers/*` directly. The policy engine evaluates the proposal as if it were a human request.
- **The AI's context is scoped.** `POST /ai/conversations` requires `context_scope: {resource_ids: [...], log_window: "1h"}`. The AI sees only what the initiating principal could see.
- **Prompt content is untrusted.** The API does not distinguish between "system prompt" and "user content" at the transport level. The AI gateway applies the boundary. This is enforced server-side, not client-side.
- **The AI is a batch workload.** When the AI runs inference, it is classified as batch. It runs low-priority, yields to serving workloads. On a 1 GB VPS, the AI is available but may take longer to respond under load.

## 24. Automations, cron, webhooks and extension SDK

| Method | Endpoint family | Operations |
|---|---|---|
| GET/POST | `/automations` | Define scheduled/event-triggered jobs |
| GET/POST | `/webhooks/endpoints` | Configure delivery destinations |
| GET | `/webhooks/deliveries` | Delivery attempts and status |
| POST | `/webhooks/deliveries/{id}/redeliver` | Governed replay |
| GET/POST | `/extensions` | List/install vetted extensions |
| GET | `/extensions/{id}/permissions` | Granted capabilities |

Automation triggers use explicit timezone, missed-run behavior, concurrency policy and owner identity. Cron tasks for hosted users execute within tenant isolation, not as root. Extension manifests declare capabilities, network access, compatibility, update and removal hooks; provider approval is required for privileged capabilities.

Additional rules:

- **Cron jobs are classified by command.** A cron job running `restic backup` is batch. A cron job hitting an external API is serving. The classification is derived from the command pattern; the user can override with `workload_category`. Misclassification is visible in `GET /automations/{id}`.
- **The plugin SDK is v2+.** `GET /extensions` returns an empty list on v1. The `ServiceTemplate` catalog is the v1 extension mechanism.
- **Webhooks include new event types:** `service.instance.running`, `service.instance.degraded`, `service.instance.uninstalled`, `payment.pending`, `payment.confirmed`, `payment.failed`, `payment.reversed`, `payment.unmatched`, `capability.changed`, `batch.queued`, `batch.started`, `batch.yielded`, `batch.completed`.

## 25. Outbound event/webhook contract

**Envelope:**

```json
{
  "id": "evt_01EXAMPLE",
  "type": "website.deployment.succeeded",
  "version": "1",
  "occurred_at": "2026-11-06T15:05:00Z",
  "organization_id": "org_01EXAMPLE",
  "resource": {"type": "deployment", "id": "dep_01EXAMPLE"},
  "operation_id": "op_01EXAMPLE",
  "data": {"status": "succeeded"}
}
```

Delivery uses HTTPS, a configured signing secret, HMAC signature over timestamp + raw payload, unique event ID and delivery ID. Consumers should reject stale timestamps and deduplicate event IDs. Retry on network errors/5xx and documented 429 with exponential backoff + jitter; stop or dead-letter after bounded attempts. Support replay through authenticated API. No secrets, private customer data or internal stack traces in event payloads. Event ordering is not globally guaranteed; consumers must tolerate duplicate/out-of-order delivery.

**Initial event taxonomy:** `organization.member.invited`, `server.agent.connected`, `server.agent.disconnected`, `website.created`, `website.deployment.started`, `website.deployment.succeeded`, `website.deployment.failed`, `domain.verified`, `dns.zone.published`, `certificate.renewal.failed`, `database.created`, `mailbox.created`, `backup.completed`, `backup.failed`, `restore.completed`, `security.alert.created`, `hosting_account.suspended`, `subscription.updated`, `operation.failed`.

Expanded event taxonomy:

**Capability events:**
- `capability.changed` — a capability became available/unavailable (hardware rescan, hardware change)
- `capability.unavailable_blocking` — a running workload now exceeds the current capability set

**Service events:**
- `service.template.deployed`
- `service.instance.running`
- `service.instance.degraded`
- `service.instance.uninstalled`
- `service.instance.exported`

**Payment events:**
- `payment.pending` — payment initiated, awaiting callback
- `payment.confirmed`
- `payment.failed`
- `payment.expired`
- `payment.reversed`
- `payment.unmatched` — callback received but no matching invoice

**Batch events:**
- `batch.queued`
- `batch.started`
- `batch.yielded` — batch job paused under contention
- `batch.resumed`
- `batch.completed`

**Workload events:**
- `workload.limit.updated`
- `workload.oom_kill` — a workload was killed by the OOM killer
- `workload.throttled` — a batch workload is yielding

**Additional payload fields:**

- Every event includes `workload_category` (serving / batch / infrastructure).
- Payment events include `provider` (`mpesa` | `stripe`), `provider_reference`, and `reconciliation_state`.
- Service events include `template_key` and `instance_state`.

## 26. Agent RPC and privileged execution contract

### 26.1 Enrollment and identity

Enrollment token is single-use, short-lived and bound to server registration and provider scope. Agent generates or securely receives a distinct identity and obtains a short-lived certificate; mTLS is mandatory. Certificates rotate automatically; compromised agents can be revoked. Agent outbound connectivity may be preferred for NAT/firewall compatibility, but exact transport must be decided and threat-modeled.

Transport change — Unix socket in v1:

- **Standalone v1 uses a Unix domain socket** (`/run/nexuspanel/agent.sock`) with peer-credential checks (`SO_PEERCRED`). The local Unix transport requires OS peer authentication, file permissions and per-operation policy binding; remote enrollment and certificate rotation are required only when remote transport is enabled. The installer provisions a locally scoped agent identity and explicitly verifies socket ownership and permissions. This is DEC-013.
- **Remote deployments use authenticated mTLS** when remote agent connectivity is enabled. The enrollment flow is the v1.0 flow (single-use token, CSR, certificate, rotation).
- **The typed operation catalog is identical.** Only the transport changes. The v1.0 contract (below) applies to both.

### 26.2 Typed command envelope (illustrative)

```json
{
  "command_id": "cmd_01EXAMPLE",
  "operation_id": "op_01EXAMPLE",
  "server_id": "srv_01EXAMPLE",
  "action": "nginx.site.apply",
  "schema_version": 1,
  "resource_id": "web_01EXAMPLE",
  "desired_revision": 8,
  "deadline": "2026-11-06T15:10:00Z",
  "policy_decision_id": "pdc_01EXAMPLE",
  "workload_category": "serving",
  "payload": {"site_config_ref": "cfg_01EXAMPLE"}
}
```

Agent MUST validate authenticated command origin and authorization binding (and signatures when supported by the chosen transport), server target, action allowlist, payload schema, tenant/resource mapping, expected revision, expiry, replay protection and local safety constraints. It MUST NOT accept generic `exec` as a normal management primitive. Privileged adapters run with the narrowest available OS capabilities and service-specific identities.

Additional envelope fields:

- **`workload_category`** is required. The agent uses it to derive the cgroup configuration if the command creates or modifies a workload.
- **`cgroup_config`** is included in the payload for workload-creating commands. The agent validates it against the workload category's constraints. Rejects `cpuset`, unknown controllers, out-of-range values.
- **`capability_decision_id`** is included. The agent verifies the capability was available at submission time and re-checks at execution.

### 26.3 Agent response and reconciliation

Return accepted/rejected status, command ID, observed revision, safe error code, structured results and progress events. Long-running commands report heartbeats; on disconnect, control plane marks state unknown until reconciliation, not automatically failed or successful. Desired and observed states are compared, drift recorded, and repairs are policy-governed. Host-level operations use atomic config writes, syntax validation, reload checks and last-known-good rollback where feasible.

### 26.4 Initial typed action catalog

`system.inventory.read`, `system.service.status`, `nginx.site.validate`, `nginx.site.apply`, `runtime.install`, `app.release.deploy`, `app.process.restart`, `files.list`, `files.read`, `files.write_atomic`, `database.provision`, `database.user.rotate`, `dns.zone.publish`, `tls.certificate.install`, `mail.mailbox.provision`, `backup.create`, `backup.restore`, `firewall.preview`, `firewall.apply`, `metrics.collect`.

Each action requires its own schema, validation rules, privilege declaration, idempotency behavior, rollback strategy and test suite.

Expanded action catalog:

**Workload / cgroup:**
- `workload.create` — create a workload with cgroup config
- `workload.update_limits` — update `cpu.weight`/`cpu.max`/`memory.max`
- `workload.remove`
- `workload.verify_cgroup` — read back the applied cgroup values
- `workload.metrics.read`

**Container (Docker):**
- `container.create`
- `container.start` / `container.stop` / `container.restart`
- `container.remove`
- `container.health_check`
- `container.logs`
- `container.validate_safety` — rejects `--privileged`, host PID/network, dangerous capabilities

**ServiceTemplate:**
- `service.provision` — pull image, create container/volume, configure proxy/SSL, run health check
- `service.configure` — update env vars, restart
- `service.start` / `service.stop`
- `service.health_check`
- `service.backup`
- `service.restore`
- `service.export` — create OCI artifact + manifest
- `service.import` — restore from OCI artifact
- `service.uninstall` — stop, remove container/volume, cleanup

**Capability:**
- `capability.detect` — run the 5-step detection
- `capability.rescan` — re-detect after hardware change

**Batch:**
- `batch.yield` — check if the server is under contention; if so, pause
- `batch.resume`

**Payment:**
- `payment.verify_callback` — verify Daraja/Stripe signature and event ID
- `payment.reconcile` — query Daraja Transaction Status API

The agent does **not** have direct payment operations. Payment resolution happens in the control plane (it's a state machine, not a host operation). The agent has no access to M-Pesa credentials. This is a key security boundary.

## 27. Service adapter and external integration contracts

An adapter declares: `name`, `version`, `capabilities`, `supported_os`, `supported_service_versions`, `health_check`, `validate`, `plan`, `apply`, `verify`, `rollback/compensate`, `inventory`, `uninstall`, `observability`. Adapter selection is capability-based; UI must not assume a service is available on every node.

| Domain | Initial adapter candidate | Alternative/future |
|---|---|---|
| HTTP reverse proxy | Nginx | Apache, Caddy |
| Application runtime | Node.js systemd/process adapter | PHP-FPM, Python, containers |
| Metadata database | PostgreSQL | Managed PostgreSQL for control plane |
| Customer SQL | MariaDB, PostgreSQL | Remote managed database services |
| Authoritative DNS | PowerDNS Authoritative | BIND, external DNS APIs |
| TLS | ACME client integration | Enterprise CA/provider certificates |
| Mail | Postfix + Dovecot stack | External mail provider integrations |
| Storage | Local filesystem | S3-compatible object storage |
| Source control | Git over SSH/HTTPS | GitHub/GitLab OAuth adapters |
| Billing | Provider-neutral interface | Selected gateway adapters |

These are design candidates, not installation instructions. Ubuntu 26.04 compatibility, packaging and security maintenance must be validated for each service/version before declaring supported.

Updated adapter table:

| Domain | v1 adapter | v2+ adapter | Notes |
|---|---|---|---|
| HTTP reverse proxy | **Caddy (Standard+) / Nginx (Nano)** | Apache | Caddy's automatic HTTPS is the simplification |
| Container runtime | **Docker** (Standard+) | Rootless Podman | Docker is the isolation primitive |
| Application runtime | **Node/Python/PHP in containers** | Container orchestration | |
| Metadata database | **SQLite (Nano) / PostgreSQL (Standard+)** | Managed Postgres | Dual-driver repository abstraction |
| Customer SQL | **MariaDB/Postgres containers** (v2+) | Remote managed | v1: unavailable |
| Authoritative DNS | **None (guidance mode)** | PowerDNS / BIND | v1: DNS guidance only |
| TLS | **ACME HTTP-01 via Caddy** | DNS-01, wildcard | v1: HTTP-01 only |
| Mail | **External SMTP (Resend/Postmark/Mailgun/SES)** | Postfix + Dovecot | v1: external only |
| Storage (platform service) | **MinIO or Garage** | S3-compatible | `ServiceTemplate` instance |
| Image transform | **imgproxy or thumbor** | Alternative | Batch workload |
| Delivery | **Caddy cache + Cloudflare/Bunny integration** | Alternative CDN | |
| Source control | **Git over HTTPS** | GitHub/GitLab OAuth | |
| Payments | **M-Pesa Daraja (reference) + Stripe** | More regional gateways | Async-first state machine |
| Backup | **restic → S3-compatible** | Borg, custom | Remote-only on Nano |
| ServiceTemplate | **Container + volume + proxy + health** | | Generic primitive for all services |

## 28. End-to-end reference workflows

### 28.1 Provision a Node.js website

1. User submits app source, runtime, hostname, environment and hosting account.
2. API authorizes actor, account, quota and domain ownership; validates runtime availability.
3. API creates desired-state resource and operation atomically with outbox event.
4. Scheduler reserves target capacity; agent validates directory isolation and runtime prerequisites.
5. Deployment adapter fetches immutable source revision, builds in constrained context and prepares release.
6. Proxy adapter validates routing configuration and TLS state.
7. Health checks verify process and route; operation marks succeeded only after verification.
8. On failure, preserve previous healthy release and provide redacted step logs and retry options.

Additions:

- **Capability check first.** Before step 2, the API checks `capabilities[runtime.node].state`. If unavailable, return `CAPABILITY_UNAVAILABLE` with resource delta.
- **Build is batch.** Step 5 creates a batch operation. The build queues, runs when idle, yields under contention.
- **Runtime is serving.** Step 7 creates a serving workload with `cpu.weight=100`, no `cpu.max`, `memory.max` set.
- **Container isolation.** Step 5 is a build container; step 7 is a runtime container. Separate cgroups, separate filesystems.
- **`.env` bulk paste is supported.** The API accepts a full `.env` file, parses it, validates, stores encrypted. Never echoes values.

### 28.2 Change DNS records

Read zone version → propose record diff → validate syntax/conflicts/ownership → plan → approval if policy requires → publish new zone revision → verify authoritative nameservers → emit event. Propagation to third-party recursive resolvers is not guaranteed and must not be misreported as immediate.

. The v1 workflow is DNS guidance:

1. User enters domain, intended use.
2. API returns `{records: [{type: "A", name: "@", value: "<vps-ip>"}], instructions: "..."}`.
3. User copies the records to their registrar.
4. `POST /domains/{id}/dns-guidance/check` polls. Returns `{propagated: false}` until resolved.
5. Once propagated, SSL issuance proceeds.

No authoritative zone editing. No propagation guarantee. The panel never claims "live" before the check succeeds.

### 28.3 Restore a customer database

Select recovery point → verify permissions and backup integrity → preview overwrite targets/downtime → capture safety backup if required → approve exact plan → restore in isolated operation → verify consistency and application connectivity → audit outcome. Partial failure must not be presented as successful restore.

Additions:

- **Restore is a batch operation.** It queues, runs when idle, yields under contention.
- **Database is a `ServiceTemplate` instance.** The restore targets the specific instance's volume.
- **Verification includes a query against the restored DB.** Not just "the restore command exited 0."
- **"Backup succeeded" ≠ "backup is restorable."** The restore is verified by a real query.

### 28.4 Provider creates reseller and customer

Provider delegates reseller scope and limits → reseller creates customer and hosting account within delegation → account receives quota reservation and placement → customer sees only owned resources → provider retains governed oversight. Permission escalation and quota bypass must fail at API and database boundaries.

. The v1 product does not ship the reseller/provider surface. The data model is complete; the API endpoints exist but return `403 FEATURE_DISABLED` on a v1 installation.

### 28.5 AI diagnoses failed deployment

User asks for diagnosis → AI requests scoped deployment logs/metrics through read-only APIs → proposes remediation and cites evidence → authorization and plan checks run → user approves → normal operation engine executes → AI reports observed result, not assumed success.

Additional rules:

- **v1 is read-only diagnostics + approval-gated actions.** The AI can diagnose, propose, and wait for approval. It cannot execute autonomously.
- **The AI's proposals are `POST /ai/action-proposals`.** The user approves via `POST /ai/action-proposals/{id}/submit`, which routes through the standard plan/approval/apply workflow.
- **The AI is a batch workload.** Under load, AI responses queue. The UI shows "AI is queued behind site traffic."

### 28.6 M-Pesa payment flow **(new in v2)**

1. User creates an invoice (or the invoice is auto-generated at subscription renewal).
2. User initiates payment: `POST /payments` with `{invoice_id, provider: "mpesa", phone_number: "+254712345678", amount_minor: 15000, currency: "KES"}`.
3. API returns `202 Accepted` with `{id, state: "pending", provider_reference: "<CheckoutRequestID>"}`.
4. API calls Daraja STK Push. User receives prompt on phone, enters PIN.
5. Daraja calls the callback URL: `POST /integrations/payments/mpesa/callback` with the result.
6. API verifies signature, deduplicates by `CheckoutRequestID`, resolves the payment state to `confirmed` or `failed`.
7. If confirmed: invoice marked paid, entitlement updated, subscription active. If failed: payment marked failed, retry allowed.
8. If no callback within timeout (e.g., 5 minutes): the payment is marked `expired`. A reconciliation job polls Daraja Transaction Status API. If Daraja reports success but no callback was received, the payment is resolved and flagged for review.

**Failure branches:** callback signature invalid → rejected, audit; duplicate callback → idempotent 200; unmatched callback → reconciliation queue; timeout → expired, reconciliation job; reversal → audit.

### 28.7 ServiceTemplate deployment flow **(new in v2)**

1. User picks a template: `POST /service-instances` with `{template_id: "nextcloud", organization_id, hostname, credentials: {generated}}`.
2. API validates: capability fit, disk quota, hostname uniqueness, port conflict, template signature.
3. API returns `202 Accepted` with operation ID.
4. Agent: pulls image (digest-pinned, signature-verified), creates container with cgroup limits, creates volume, configures reverse proxy + SSL, generates credentials, runs health check.
5. If healthy: state `running`. If not: state `degraded`, logs surfaced, user can restart/reinstall/uninstall.
6. Post-provision: `GET /service-instances/{id}` returns `{state: "running", hostname: "cloud.name.panel-host", admin_credentials_shown_once: true}`.

**Failure branches:** image signature invalid → rejected, audit; health check fails → degraded; disk full → rejected; port conflict → rejected.

### 28.8 Hardware rescan after VPS upgrade **(new in v2)**

1. User's provider bumps the VPS to 4 GB RAM.
2. User triggers: `POST /servers/{id}/capabilities/rescan`.
3. API returns `202 Accepted`. Agent runs 5-step detection.
4. New capabilities computed. Diff shown: `{new_available: ["app.node", "db.mariadb", "service.docker"], new_borderline: [...]}`.
5. User confirms: `POST /servers/{id}/capabilities/apply` with the accepted diff.
6. Panel now allows the new workloads. No reinstall. No migration.

This is the ladder climbing.

## 29. Threat model and integration security controls

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

Additional threats:

- **M-Pesa callback forgery:** verify Daraja signature, event ID, merchant identity, amount/currency match. Tested with negative signatures and replay.
- **Payment state transition bypass:** a `pending` payment cannot jump to `confirmed` without a callback. A `confirmed` payment cannot be reversed without an explicit reversal event. Enforced by state machine.
- **Idempotency bypass via payment retry:** the idempotency key is composed of actor + invoice + amount + provider. Two different actors cannot share the same key. Duplicate submissions return the existing payment.
- **cgroup misconfiguration:** the agent validates all cgroup values; rejects `cpuset`, unknown controllers, out-of-range values.
- **Container escape:** the agent rejects `--privileged`, host PID/network, dangerous capabilities. Tested with malicious Compose.
- **Capability spoofing:** capability records are loaded from signed data files. A tampered file is rejected at install time. The installer signature is verified.
- **Persona matrix tampering:** the persona matrix is data, but it is signed and verified on load. The API rejects unsigned data.

## 30. Observability, SLIs, SLOs and diagnostics

Every request gets a request ID and trace context; operations link API request, orchestration job, agent command and audit event. Track API latency/error rate, operation completion latency, queue age, agent heartbeat freshness, reconciliation drift, deployment success rate, certificate renewal success, backup verification and webhook delivery success. Define SLO targets after measuring baseline; do not invent contractual uptime guarantees. Metrics must be privacy-safe and cardinality-bounded.

Additional observability:

- **Capability gate metrics:** how often is `CAPABILITY_UNAVAILABLE` returned? Which capabilities are most often requested but unavailable? This tells you what to build next or what to enable on smaller profiles.
- **Batch queue depth:** the number of queued batch jobs. Sustained depth > 10 indicates a capacity problem.
- **Payment state metrics:** how many payments are `pending` for > 5 minutes? How many reach `reconciliation_state = unmatched`? How many are reversed?
- **Router conversion metrics:** which persona paths are chosen? Where do users drop off? (Note: this is client-side telemetry, not API.)
- **cgroup throttling metrics:** how often are serving workloads throttled? If > 5% of time, the CPU overcommit is too aggressive.
- **OOM kill rate:** how many OOM kills per day? If > 0 for serving workloads, the RAM ceiling is too tight.

## 31. OpenAPI, SDK, documentation and developer experience

Deliver a machine-readable OpenAPI 3.1 specification with schemas, error responses, security requirements, idempotency headers, pagination, example payloads and operation links. Generate TypeScript and Go clients with pinned versions and contract tests. Provide API explorer, sandbox/test fixtures, CLI commands, webhook signature verification examples, migration guides, rate-limit documentation and deprecation notices. Internal agent protocol uses separately versioned protobuf or equivalent IDL with backward compatibility tests.

Additional DX:

- **The OpenAPI is generated from Go types.** `oapi-codegen` (or equivalent) produces the spec from the source of truth. This prevents drift.
- **The SDK is generated.** TypeScript and Go SDKs are generated from the OpenAPI. Hand-written SDKs drift; generated SDKs don't.
- **The sandbox environment is the dev VPS.** For v1, there is no hosted sandbox. Developers run `nexuspanel install --dev` on a disposable VPS or local VM.
- **The CLI is `nexus`.** `nexus websites create`, `nexus apps deploy`, `nexus services provision`, `nexus payments list`. It uses the same API.
- **The API explorer is embedded in the panel.** `/developer/api` (Stage 3 UX-AUTO-005) shows the OpenAPI spec, allows authenticated requests, and displays the response.

## 32. API conformance and acceptance test matrix

| Test ID | Scenario | Pass criterion |
|---|---|---|
| API-001 | Cross-tenant GET | No data or existence leakage |
| API-002 | Cross-tenant mutation | Denied; no host side effects |
| API-003 | Reseller exceeds delegation | Denied with audited reason |
| API-004 | Duplicate POST with same key/body | One resource/operation, replayed result |
| API-005 | Same key with different body | 409 conflict |
| API-006 | Concurrent quota reservation | No over-allocation |
| API-007 | ETag mismatch | 412; unchanged resource |
| API-008 | Worker crash mid-operation | Recovery/reconciliation; no false success |
| API-009 | Agent offline | Operation pending/unknown with clear status |
| API-010 | Invalid agent action | Rejected before execution |
| API-011 | Expired/replayed agent command | Rejected and audited |
| API-012 | Firewall lockout risk | Plan blocks or requires safe recovery |
| API-013 | Invalid Nginx config | Existing healthy config preserved |
| API-014 | Secret in environment | Never returned in read/list/logs |
| API-015 | File traversal/symlink escape | Blocked on agent |
| API-016 | DNS zone change | Versioned publish and authoritative verification |
| API-017 | Backup restore failure | Partial state surfaced; remediation offered |
| API-018 | Duplicate webhook | Consumer-safe deduplication support |
| API-019 | Spoofed webhook signature | Rejected |
| API-020 | AI unauthorized action | No privilege escalation or host effect |
| API-021 | API version compatibility | Existing supported clients pass contract suite |
| API-022 | Server capability mismatch | Placement rejected or alternative selected |
| API-023 | Payment callback replay | No double credit/charge |
| API-024 | Audit export | Actor, action, target, result and timestamps preserved |
| API-025 | Cancel after irreversible step | Honest partial/cannot-cancel result |

Additional API conformance tests:

| Test ID | Scenario | Pass criterion |
|---|---|---|
| API-026 | Capability unavailable | `CAPABILITY_UNAVAILABLE` returned with reason and resource delta |
| API-027 | Capability rescan | New capabilities visible after rescan; no reinstall |
| API-028 | Allocation vs usage | `GET /resources/{id}/usage` returns allocated/typical/peak/current, not just allocated |
| API-029 | Workload category immutability | Attempting to reclassify returns `WORKLOAD_CATEGORY_MISMATCH` |
| API-030 | cgroup config validation | `cpuset` rejected; RAM ceiling > physical × 1.2 rejected |
| API-031 | Batch queue visibility | `GET /batch-queue` returns queued jobs with positions |
| API-032 | Batch yield | Batch job yields under contention; resumes when idle |
| API-033 | ServiceTemplate deploy | `POST /service-instances` provisions; health check passes |
| API-034 | ServiceTemplate export/import | Export artifact; import on new VPS; instance reproduces |
| API-035 | M-Pesa payment initiate | `POST /payments` returns 202 with pending; STK push sent |
| API-036 | M-Pesa callback | Signature verified; state resolves; invoice paid |
| API-037 | M-Pesa duplicate callback | Idempotent 200; no double credit |
| API-038 | M-Pesa callback unmatched | Payment marked unmatched; reconciliation queue populated |
| API-039 | M-Pesa timeout | Payment expired; reconciliation job polls Daraja Transaction Status |
| API-040 | Stripe webhook | Signature verified; state resolves; no double charge |
| API-041 | Async payment state machine | Both M-Pesa and Stripe use the same pending→confirmed→failed states |
| API-042 | Audit chain verification | `POST /audit-logs/verify` returns `chain_intact: true`; tampered chain returns false |
| API-043 | Persona matrix signing | Tampered matrix rejected at load |
| API-044 | Capability file signing | Tampered capability file rejected at install |
| API-045 | Agent Unix socket peer check | Non-panel process cannot connect to agent socket |
| API-046 | Router landing | `GET /me/default-landing` returns correct landing for each persona |
| API-047 | Nano install | Panel installs on 512 MB; `/capabilities` returns honest set; unsupported endpoints return `CAPABILITY_UNAVAILABLE` |
| API-048 | Nano → Standard upgrade | After hardware bump, `POST /capabilities/rescan` enables new capabilities; no reinstall |

## 33. Implementation sequencing without feature exclusion

**Foundation:** identity, tenancy, authorization, OpenAPI skeleton, request/error standards, resource registry, operation engine, audit, agent enrollment, inventory.
**Core hosting:** websites, routing, Node.js/PHP/Python runtimes, files, databases, TLS, DNS.
**Commercial hosting:** accounts, packages, resellers, quotas, billing, customer isolation, email, backups.
**Distributed operations:** scheduling, multi-server placement, agent upgrades, disaster recovery, cross-node migration.
**AI and ecosystem:** diagnostics, action proposals, plugin SDK, marketplace integrations, advanced automation.

This is a dependency-aware delivery sequence, not a restriction of the final architecture. Each module requires schema, API, adapter, permissions, UX, tests and operations runbook before release.

Aligned with the v1/v2/v3 phases from Master Blueprint §12:

**v1 (Foundation + first-complete-journey):** identity, tenancy, authorization, capability model, installer, operation engine, audit chain, secrets, agent (Unix socket), websites (static + WordPress + Node), ServiceTemplate primitive (storage, transform, delivery, Nextcloud), backups (restic), domain connection wizard, SSL (HTTP-01), guided router, persona matrix, allocation/usage endpoints, batch queue, async payment (M-Pesa), read-only AI diagnostics, external SMTP.

**v2 (Expand the surface):** provider/fleet, reseller, Stripe adapter, email stack (self-hosted), authoritative DNS, plugin runtime, full AI action execution, migration from cPanel/DirectAdmin, template catalog expansion (20–30 apps), mTLS remote agent.

**v3 (Commercial maturity):** full billing automation, multi-jurisdiction tax, white-label, marketplace, cross-region backups, HA topology.

## 34. Open decisions requiring approval

1. Product/API naming and public documentation domain.
2. Primary authentication/session implementation and MFA factors.
3. OAuth/OIDC and enterprise SSO release requirements.
4. Organization hierarchy and delegated reseller semantics.
5. API version compatibility/deprecation support duration.
6. Operation retention and idempotency-key TTL.
7. SSE vs WebSocket for terminal/real-time streams.
8. Agent transport (gRPC, HTTP/2, outbound connection topology).
9. Agent certificate authority, rotation and revocation policy.
10. Agent command schema versioning and signing/fencing strategy.
11. Initial supported Ubuntu versions and service versions.
12. First-class Nginx vs Apache support sequencing.
13. Runtime isolation mechanism for Node/Python/PHP apps.
14. File editor maximum size and archive/upload policy.
15. Authoritative DNS adapter and independent secondary strategy.
16. Mail stack, abuse controls and outbound port constraints.
17. Backup object storage provider and encryption key custody.
18. Billing gateway(s), currency/tax handling and licensing model.
19. Webhook delivery retry window and event retention.
20. Plugin sandbox and extension signing policy.
21. AI model/data residency and redaction policy.
22. SLOs, support levels and audit log retention.
23. API pagination/filter naming standard and default limits.
24. Quota grace periods and subscription enforcement policy.
25. Approval thresholds for destructive or network-impacting actions.

Additional decisions:

26. **Dual-driver repository strategy** — shared DSL + CI equivalence check (DATA-021).
27. **Capability data file format and signing** (DATA-022).
28. **Persona matrix data file format and signing** (DATA-023).
29. **ServiceTemplate artifact format** — OCI image + signed JSON manifest (DATA-024).
30. **Payment idempotency scope** — actor + invoice + amount + provider (DATA-025).
31. **M-Pesa callback secret** — per-tenant, verified on every callback (DATA-026).
32. **cgroup metric retention** — raw 30 days, hourly 1 year, daily 5 years (DATA-027).
33. **Audit chain export cadence** — hourly append to write-protected storage (DATA-028).
34. **Batch queue priority derivation** — from workload_category + operation_kind; user cannot override (DATA-029).
35. **Nano → Standard migration path** — stop → export SQLite → start Postgres → import → resume (DATA-030).
36. **Guided router as primary IA** — release gate, not feature (UXD-021/022/026/031).
37. **cgroup v2 requirement** — Ubuntu 22.04+ baseline (DEC-015).
38. **Go stack + static SPA + dual DB** (DEC-001).

## 35. Stage 6 exit criteria and next artifacts

- [ ] Every Stage 2 module has mapped endpoints and authorization requirements.
- [ ] Every Stage 3 workflow has corresponding API/operation states.
- [ ] Every Stage 5 entity has an owner and API exposure policy.
- [ ] Operation state machine and idempotency semantics approved.
- [ ] Agent protocol, trust boundary and action allowlist approved.
- [ ] Error schema, pagination, concurrency and versioning approved.
- [ ] Webhook security and retry contracts approved.
- [ ] Threat model and cross-tenant tests reviewed.
- [ ] Open decisions have named owners and resolution dates.
- [ ] OpenAPI 3.1 and agent IDL generation tasks created.

Additional exit criteria:

- [ ] **Capability-aware endpoints are specified for every workload-creating operation.** `CAPABILITY_UNAVAILABLE` is defined, tested, and rendered honestly.
- [ ] **`ServiceTemplate` endpoints are the single API for all services.** No bespoke endpoints.
- [ ] **Async payment state machine is specified** for both M-Pesa and Stripe with the same API contract.
- [ ] **Allocation/usage endpoints are specified** for every resource. `GET /resources/{id}/usage` returns allocated/typical/peak/current.
- [ ] **Batch queue endpoints are specified.** `GET /batch-queue`, `POST /operations/{id}/yield`.
- [ ] **The persona matrix is data, not code.** `GET /me/persona-matrix` is signed and verified.
- [ ] **Audit chain verification endpoint** is specified and tested.
- [ ] **Unix socket agent transport** is specified for v1; mTLS is documented as v2+.
- [ ] **The v2 conformance test matrix (API-026 through API-048) is adopted** as the release gate.
- [ ] **The guided router is a release gate, not a later feature.** This is the one thing to hold the line on.

**Recommended next stage:** Stage 7 — Privileged Server Agent & Service Management v2.0: detailed Go agent design, host-level operation schemas, privilege separation, OS adapters, configuration validation, reconciliation, safe rollback, installer integration and agent conformance tests.

---

**Document status reminder:** This is a complete architectural API specification draft across the product domains. It is not a substitute for the endpoint-by-endpoint OpenAPI definitions, generated SDKs, or agent IDL required before production implementation.

**End of Stage 6 — API & Integration Specification v2.0.**

---



---

## 36. Consolidated cross-area API traceability (29 areas)

| Area | Primary API family | Critical gate / contract |
|---|---|---|
| 01 Dashboard | `/overview`, `/resources/{id}/usage` | Scoped metrics, staleness labels |
| 02 Server management | `/servers`, `/services`, `/operations` | Host ownership, typed agent operations |
| 03 Websites | `/websites`, `/domains` | DNS/TLS and document-root validation |
| 04 DNS | `/dns/zones`, `/dns/records` | Zone ownership, external/manual mode |
| 05 SSL/TLS | `/certificates`, `/tls` | ACME challenge and expiry lifecycle |
| 06 File manager | `/files`, `/uploads` | Path containment and quota |
| 07 Application manager | `/applications`, `/deployments` | Image/runtime capability, isolation |
| 08 Runtime manager | `/runtimes`, `/runtime-versions` | Supported version/extension compatibility |
| 09 Databases | `/databases`, `/database-users` | Tenant DB separation, restore integrity |
| 10 Email hosting | `/mailboxes`, `/mail-domains` | Provider/self-hosted mode and deliverability |
| 11 Git/deployments | `/repositories`, `/deployments` | Secrets, build isolation, rollback |
| 12 Backups | `/backups`, `/restores` | Retention, verified restoration |
| 13 Security | `/security`, `/firewall`, `/ssh` | Step-up and high-risk approval |
| 14 Monitoring | `/metrics`, `/alerts`, `/logs` | Access scope, bounded retention |
| 15 Terminal/developer | `/terminals`, `/cron` | Session scope, audited shell limits |
| 16 User/account | `/organizations`, `/memberships`, `/roles` | Delegation, tenant isolation |
| 17 AI assistant | `/ai/diagnostics`, `/ai/proposals` | Read-only v1; policy-controlled approval |
| 18 Automation/API | `/api-keys`, `/webhooks`, `/operations` | Scoped keys, signatures, replay defense |
| 19 Panel administration | `/panel`, `/updates`, `/licenses` | Safe update/rollback and break-glass |
| 20 Extensions | `/extensions`, `/plugins` | Signed packages and capability sandbox |
| 21 Commerce/resellers | `/subscriptions`, `/invoices`, `/payments` | Financial ledger and provider reconciliation |
| 22 Migration/staging | `/migrations`, `/staging` | Source ownership and safe cutover |
| 23 Containers | `/containers`, `/service-instances` | Deny privileged/host namespace defaults |
| 24 Performance/caching | `/cache`, `/performance` | Tenant cache boundaries and invalidation |
| 25 Email security | `/mail/security`, `/mail/deliverability` | SPF/DKIM/DMARC validation, reputation |
| 26 Analytics/site quality | `/analytics`, `/site-quality` | Privacy, sampling, bot attribution |
| 27 Audit/compliance/support | `/audit`, `/support` | Immutable evidence, scoped support grants |
| 28 Platform Services | `/service-templates`, `/service-instances` | Storage/transform/delivery/queue/search specialization |
| 29 Self-hosted Templates | `/service-templates`, `/service-instances` | Signed image, dependencies, upgrade/backup contract |

## 37. Mandatory machine-readable contract examples

The following are **illustrative** request/response shapes. Exact names, enums and field constraints must be stabilized in OpenAPI before coding.

### 37.1 Capability-aware service installation

```json
{
  "template_id": "tmpl_example",
  "target_server_id": "srv_example",
  "hosting_account_id": "acct_example",
  "configuration": {"hostname": "app.example.test"},
  "requested_limits": {"memory_bytes": 268435456, "cpu_weight": 100},
  "expected_capability_revision": 17
}
```

On acceptance, return HTTP 202:

```json
{
  "operation_id": "op_example",
  "state": "queued",
  "resource_id": "svc_example",
  "links": {"operation": "/api/v1/operations/op_example"}
}
```

On infeasibility, return a non-leaking problem response after authorization:

```json
{
  "type": "about:blank",
  "status": 409,
  "code": "CAPABILITY_UNAVAILABLE",
  "detail": "This server cannot currently host the selected service.",
  "capability": "container.runtime",
  "reason_codes": ["DOCKER_NOT_AVAILABLE", "MEMORY_HEADROOM_INSUFFICIENT"],
  "alternatives": ["select_another_server", "use_remote_service"]
}
```

### 37.2 Agent command authority envelope

```json
{
  "command_id": "cmd_example",
  "operation_id": "op_example",
  "server_id": "srv_example",
  "resource_id": "svc_example",
  "action": "service.provision",
  "schema_version": 1,
  "desired_revision": 4,
  "policy_decision_id": "policy_example",
  "capability_decision_id": "cap_example",
  "deadline": "2026-11-06T21:30:00Z",
  "workload_category": "serving",
  "payload_ref": "encrypted_config_example"
}
```

The agent must independently verify the authenticated transport peer, the bound resource and action, command deadline, revision, local allowlist, host safety, and replay/idempotency journal. For remote transport, the mTLS identity must match the intended agent. No browser, plugin or LLM may submit this envelope directly.

### 37.3 Payment verification record

```json
{
  "payment_id": "pay_example",
  "provider": "mpesa",
  "provider_reference": "reference_example",
  "invoice_id": "inv_example",
  "expected_amount_minor": 50000,
  "currency": "KES",
  "state": "pending",
  "verification_state": "awaiting_provider_confirmation",
  "reconciliation_state": "open"
}
```

Callback receipt is **not** payment confirmation. Validate provider-specific callback semantics and correlation, then reconcile status using supported provider mechanisms. Persist provider event IDs/references, reject duplicate ledger effects, quarantine ambiguous or conflicting outcomes, and provide audited manual review. Do not assert a cryptographic Daraja callback signature unless the actual product/provider contract supports it.

## 38. API conformance release gates

| ID | Gate | Evidence required |
|---|---|---|
| API-GATE-001 | OpenAPI coverage | Every shipped endpoint has request/response/error schemas and security scopes |
| API-GATE-002 | Tenant isolation | Cross-tenant read/write fuzz and E2E tests show no resource disclosure |
| API-GATE-003 | Durable mutations | Crash/retry tests demonstrate operation idempotency and reconciliation |
| API-GATE-004 | Agent boundary | Local socket peer-credential and remote mTLS tests where applicable |
| API-GATE-005 | Capability gating | Nano/Starter/Standard/Provider fixtures produce accurate reason codes |
| API-GATE-006 | Payment reconciliation | Duplicate, delayed, spoofed, reordered, reversed and unmatched events tested |
| API-GATE-007 | ServiceTemplate lifecycle | Signed artifact, dependency, rollback and specialized-stateful tests |
| API-GATE-008 | Webhooks | Signature validation, rotation, retry/backoff and replay protection |
| API-GATE-009 | Secret redaction | No secrets in API read responses, logs, traces or webhook payloads |
| API-GATE-010 | SDK compatibility | Generated client contract tests match the shipped OpenAPI |
| API-GATE-011 | Resource metrics | Current/typical/peak semantics include units, timestamps and unknown states |
| API-GATE-012 | Authorization consistency | UI, CLI, SDK and AI use identical permission decisions |
| API-GATE-013 | Profile portability | Equivalent API state transitions on SQLite and PostgreSQL |
| API-GATE-014 | Degraded dependencies | Agent offline, DNS unavailable, payment timeout, full disk and queue saturation |
| API-GATE-015 | Release traceability | Each shipped `NP-*` requirement links to endpoint, scope, tests and evidence |

## 39. Approval and verification register

- **APPROVAL-API-01:** Confirm exact router matrix endpoint and ownership of jobs-to-be-done preferences.
- **APPROVAL-API-02:** Finalize stable endpoint naming for generic service instances versus specialized stateful resources.
- **APPROVAL-API-03:** Finalize provider-specific Daraja verification and settlement/reversal behavior against the selected account/product contract.
- **APPROVAL-API-04:** Define agent remote deployment release and mTLS enrollment/rotation schedule.
- **APPROVAL-API-05:** Approve idempotency retention, rate limits, webhook delivery windows and API deprecation periods.
- **VERIFY-API-01:** Prove that the local Nano control plane, durable queue and telemetry meet memory/storage budgets under concurrent activity.
- **VERIFY-API-02:** Prove SQLite/PostgreSQL transactional equivalence for all shipped API workflows, not only schema shape.
- **VERIFY-API-03:** Prove that cgroup controller availability and nesting match supported Linux distributions and provider virtualization.
- **VERIFY-API-04:** Validate gateway callback authenticity controls against current provider documentation before production.

**End of Volume 05.** Volume 06 owns privileged-agent implementation and host-service adapters; Volume 08 owns full security policy; Volume 09 owns detailed commercial rules. This volume defines the integration contract without claiming those later volumes are complete.
