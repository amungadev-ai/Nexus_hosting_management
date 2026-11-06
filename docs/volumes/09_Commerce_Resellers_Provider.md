# NexusPanel — Master Product & Engineering Specification
## Volume 09 — Commerce, Resellers & Provider Operations

**Version:** Consolidated draft 1.0 · **Status:** Proposed, not approved or deployed · **Scope:** Commercial spine, licensing, metering, hosting packages, delegated reseller administration, provider fleet operations and commercial support.  
**Authority:** Volume 01 defines product requirements; Volumes 03–08 define architecture, persistence, API, agent, service and security boundaries. This volume supplies implementation-grade commercial behavior and provider operating contracts, without changing those boundaries.  
**Provenance:** Consolidated from NexusPanel Master Blueprint v2.0, Stage 2 PRD v2.0, 29-area register v2.0, and master Volumes 01–08. Proposed elaborations are explicitly described as design contracts rather than validated implementation facts.

## 1. Executive contract and invariants

NexusPanel's **commercial spine exists from v1**, even though reseller self-service and full provider fleet commerce are phased later. The commercial data model is **gateway independent**: `Plan → Subscription → Invoice → Payment → Entitlement`. A gateway adapter handles provider-specific interactions and cannot be the source of truth for local entitlements without verified, reconciled events. **M-Pesa Daraja (STK Push and C2B) is the proposed first adapter; Stripe is the proposed second adapter.** Currency and tax configuration are data, not hard-coded behavior. The Kenya-oriented reference configuration (KES, VAT and invoice format) requires current legal/tax validation before production use.

**Non-negotiable invariants:**
1. Every commercial object belongs to a scoped organization, with provider/reseller/customer ownership and delegated authority enforced on every request and background job.
2. The billing ledger and invoice records are not editable by changing a payment gateway's callback or a browser redirect. Financial transitions are append-only or reversal/adjustment-based and idempotent.
3. A successful payment does not mean hosting is already provisioned; entitlement activation and host provisioning are separately observable durable operations.
4. Account suspension never silently deletes hosting data. Termination/purge is a distinct policy-controlled workflow with retention, export, legal hold and audit.
5. Package CPU/memory/disk limits are **ceilings, not dedicated hardware reservations**, unless a specific separately priced reservation service explicitly guarantees capacity.
6. Capability constraints override a paid entitlement: purchasing a plan cannot make a Nano host support Docker, databases, PHP or other unavailable local services.
7. A reseller cannot delegate more privileges, resources, quota or feature access than the provider granted. Branding is presentation only, not a trust boundary.
8. Billing outages and license-server outages must not unexpectedly destroy availability of existing customer workloads; offline grace, least-disruptive fallback and explicit policy govern behavior.
9. Provider-wide aggregate reporting must not grant routine access to customer secrets, message contents or private data.
10. Commercial operations are observable, reconciled and recoverable across process restarts and retries.

## 2. Commercial product model and release ladder

| Profile / release | Commercial and provider behavior | Explicit exclusions or gates |
|---|---|---|
| Nano v1 | Standalone identity, plan/license status, lightweight admin billing, M-Pesa reference checkout, resource/capability warnings | No implied local mail/database/Docker; no provider fleet or reseller storefront |
| Starter v1 | Same commercial contracts, eligible lightweight hosting features | Availability determined by measured host capacity |
| Standard v1 | Same commercial spine, hosted services where supported, provider-operated admin billing | Self-service reseller portal is not v1 |
| Provider v2+ | Delegated reseller hierarchy, hosting packages, customer lifecycle, billing portal, scoped support and white-label | Must pass isolation, metering and reconciliation tests |
| Provider v3+ | Fleet placement, wholesale/retail settlement, multi-node operations, richer automation, HA and broader provider ecosystem | Requires multi-node trust, availability and financial assurance evidence |

The release ladder is not a substitute for implementation acceptance: an item is **not shipped** until tests, documentation, rollout and support policies are complete.

## 3. Authoritative requirements preserved from the upgraded PRD

The following NP-* identifiers are inherited and remain binding proposed requirements. The explanatory v2 notes are incorporated as requirements rather than editorial annotations.


## 6. Hosting accounts, packages, subscriptions, quotas, and reseller operations

- **NP-ACC-001:** Providers define package templates with disk, traffic, sites, domains, databases, mailboxes, backups, runtime and compute limits where supported.
- **NP-ACC-002:** Resellers can define sub-packages only within delegated provider ceilings and permitted feature sets.
- **NP-ACC-003:** Create hosting accounts with unique identity, resource allocation, owner, selected package, lifecycle state and placement.
- **NP-ACC-004:** Support pending/active/suspended/grace/terminating/terminated states with explicit allowed operations per state.
- **NP-ACC-005:** Meter usage at account and provider level with transparent units, sampling intervals and billing-grade reconciliation where required.
- **NP-ACC-006:** Enforce hard limits for provisionable resources and configurable soft warnings for bandwidth/storage thresholds.
- **NP-ACC-007:** Package changes preview effects, quota violations, proration handoff and migration requirements; never silently delete excess resources.
- **NP-ACC-008:** Account suspension defines behavior for websites, apps, mail, SSH, API, cron, data retention and backup; avoid unspecified partial suspension.
- **NP-ACC-009:** Termination supports retention windows, export, legal holds where applicable, approval and final purge tracking.
- **NP-ACC-010:** Support account transfer between authorized owners/resellers with domain, billing, data and permission reconciliation.
- **NP-ACC-011:** Provide reseller branding, delegated customer management, plan visibility, reporting and support tooling.
- **NP-ACC-012:** Provider retains visibility into fleet health and abuse signals while enforcing tenant data-access boundaries.

**Acceptance scenarios:** Concurrent account creation cannot oversubscribe hard package limits; reseller cannot allocate beyond delegated ceiling; suspension/restore produces predictable service states.


- **NP-ACC-013:** Package limits are **ceilings**, not reservations. A package that says "4 vCPU" means the account's workloads share up to 4 vCPU under contention, not that 4 vCPU is reserved. The API surfaces both allocation (ceiling) and typical usage (measured) for every account.
- **NP-ACC-014:** The Nano profile is a **package type** with a specific capability set (static sites + storage client + remote backup + AI diagnostics). Nano accounts cannot be allocated PHP/Node/database/Docker/transform resources because the physical VPS cannot host them. The installer enforces this at install time; the API enforces it at allocation time.
- **NP-ACC-015:** The **reseller surface** is deferred to v2+ (Master Blueprint §12.2). The **data model** (organizations, memberships, delegations) exists from v1. The lay person is an organization with one member; the reseller is an organization with delegated children; the provider is an organization with fleet scope. Same model at every tier.


## 18. Billing, licensing, payments and commercial operations

- **NP-BILL-001:** Provider defines catalog of products, plans, billing periods, setup fees, metered dimensions and entitlements.
- **NP-BILL-002:** Subscription lifecycle covers trial, active, grace, past due, suspended, cancelled and expired with configurable policy.
- **NP-BILL-003:** Integrate payment processors/invoicing/tax services through replaceable adapters; do not store raw payment card data.
- **NP-BILL-004:** Maintain reconciliation between provisioning state, subscription entitlements and billable usage.
- **NP-BILL-005:** Support invoices, credits, refunds, disputes, receipts, taxes and exportable financial records via integration or internal ledger as decided.
- **NP-BILL-006:** Support reseller wholesale pricing, retail price visibility, margin/commission models and delegated billing responsibility.
- **NP-BILL-007:** Payment failures trigger policy-driven notifications and grace/suspension actions with human override and audit.
- **NP-BILL-008:** Billing actions do not delete hosting data automatically; termination follows explicit retention policy.
- **NP-BILL-009:** Licensing distinguishes self-hosted standalone installations, provider deployments and add-on entitlements according to chosen commercial model.
- **NP-BILL-010:** Billing data and payment credentials have stricter access, retention and export controls.

**Acceptance scenarios:** Provider provisions plan and usage reconciles; past-due state follows configured grace policy; reseller cannot see provider wholesale prices unless permitted.


- **NP-BILL-011:** The **ledger is domain-shaped, not gateway-shaped.** The objects are `Plan → Subscription → Invoice → Payment → Entitlement`, each with its own state machine and independent lifecycle. The payment gateway is a plugin at the edge. This is the only way to support both card-first (Stripe) and mobile-money-first (M-Pesa) without two code paths.
- **NP-BILL-012:** The **payment adapter interface is async-first.** Every payment is a state machine: `initiated → pending → {confirmed | failed | expired | reversed}`. Card gateways (Stripe) resolve `pending` in <1s. Mobile money (M-Pesa) resolves in 10s–2min. Both fit the same interface. Building the interface synchronous-first and bolting async on later produces two code paths forever — that is the failure mode to avoid.
- **NP-BILL-013:** **M-Pesa is the reference implementation**, not Stripe. Specifically: Daraja STK Push + C2B; KES currency; 16% VAT; KRA-compliant invoice format (PIN, VAT number, line-item tax codes); callback URL registration per tenant; sandbox + production credentials separately stored; mobile-money-first UX (phone number input + push notification, not card form). This is DEC-006 in Master Blueprint §15.
- **NP-BILL-014:** **Stripe is the second adapter**, implemented after M-Pesa. It proves the adapter interface for the international market. Card-first is a special case of async-first where `pending` resolves in <1s.
- **NP-BILL-015:** **Licensing is a tenant of your own commercial spine, not a separate system.** If a Kenyan reseller buys a NexusPanel license, the license purchase supports M-Pesa. **License enforcement: online check + offline grace (7–14 days).** A VPS panel that stops working because the license server is down is not commercially viable.
- **NP-BILL-016:** **Currency, tax, and invoice format are configurable from v1.** Currency is a first-class field on plans, invoices, payments. Tax rules are configuration, not code. Invoice templates are data, not code. Language is locale-ready. This is what enables East African deployment without a rewrite.
- **NP-BILL-017:** **v1 scope of the commercial spine:** data model complete; async state machine complete; M-Pesa adapter complete; license activation + offline grace complete; admin-only billing surface complete; KES + Kenya VAT + KRA invoice format complete. **v2+ scope:** self-service tenant billing portal; multi-currency invoicing at scale; reseller billing markup; automated dunning; Stripe adapter.


## 4. Domain aggregates and data ownership

| Aggregate | Owner and key fields | Invariant and event |
|---|---|---|
| Organization | provider/reseller/customer/standalone; status; parent delegation edge | Membership never implies universal parent access; `organization.updated` |
| Delegation | grantor, grantee, scopes, ceilings, expiry, revocation | Grant is subset of grantor; `delegation.revoked` |
| Hosting package | provider owner, immutable published revision, feature ceilings, price references | Existing accounts remain on explicit revision until migration; `package.published` |
| Hosting account | tenant owner, package revision, server placement, lifecycle state | Account identity is not root OS identity; `account.state_changed` |
| Plan | merchant/provider, currency, billing period, pricing/tax configuration, revision | Published price revision immutable; `plan.published` |
| Subscription | subscriber, plan revision, term, state, renewal/grace policy | Entitlement changes via state machine; `subscription.transitioned` |
| Invoice | invoice number, issuer, buyer, lines, tax snapshots, currency, totals | Issued invoice is immutable; corrections use credit note; `invoice.issued` |
| Payment | invoice, adapter, merchant reference, amount, currency, state | One verified settlement may not credit twice; `payment.confirmed` |
| Ledger entry | journal, debit/credit account, amount minor, currency, posting ref | Journal balances; reversals preserve history; `ledger.posted` |
| Entitlement | subject, product/capability, effective interval, source, version | Effective grants intersect authorization and host capability; `entitlement.changed` |
| License | installation ID, signed grant/lease, feature set, expiry/grace, revocation | Offline grace policy explicit; `license.revalidated` |
| Usage record | source resource, meter, interval, units, provenance, immutable version | Recomputable and deduplicated; `usage.finalized` |
| Provider server | provider org, agent identity, capabilities, health, placement status | Quarantined nodes cannot receive placements; `server.state_changed` |
| Support session | actor, target, reason, approval, TTL, scope, audit ID | Time-limited, least privilege; `support.session_ended` |

**Persistence:** On Nano, use SQLite with transactional uniqueness and a durable outbox; on Standard/Provider, PostgreSQL with equivalent logical constraints and transactional outbox. Never implement different financial semantics per database engine. Customer-hosted databases remain outside the control-plane store. Financial amounts are **integer minor units**, currencies are ISO codes, quantities are exact decimal representations where needed, and timestamps use UTC plus explicit display locale. No floating-point currency arithmetic. Each invoice preserves the pricing, currency, tax rule version and customer/seller identity snapshot effective at issue time.

## 5. Tenant hierarchy, delegation and role boundaries

A user may belong to multiple organizations. The **acting organization** is explicit and checked on every request; UI persona (provider, reseller, customer) is not an authorization grant. Provider control over infrastructure does not automatically authorize viewing customer secrets. Support access requires separately scoped elevation, reason, approval/consent policy, expiry and full audit.

| Action | Provider operator | Reseller delegate | Hosting customer | Billing specialist | Auditor |
|---|---|---|---|---|---|
| Publish provider packages | With `package.publish` | No; may publish delegated sub-packages | No | No | Read only |
| Create hosting account | Provider scope + quota | Delegated subtree + quota | Only explicit self-service entitlement | No | No |
| Change account package | Scoped with impact preview | Only permitted sub-packages | Request or authorized self-service | No | Read only |
| Read wholesale rates | Explicit finance scope | Only if explicitly granted | No | Authorized scope | Approved audit scope |
| Issue credit/refund | Finance + approval policy | Only if delegated finance policy permits | No | Finance + separation of duties | No |
| View customer content/secrets | No by default | No by default | Own resources | No | No by default |
| Suspend abusive account | Authorized provider policy | Delegated scope only | No | No | No |
| Configure merchant adapter | Merchant admin + step-up | Only if merchant-of-record model permits | No | Limited finance scope | Read metadata only |
| View aggregate fleet health | Authorized provider scope | Own delegated nodes/accounts | Own service status | Financial summaries only | Scoped reports |

**Delegation rules:** Explicitly represent `grantor_org_id`, `grantee_org_id`, resource subtree, permission set, ceilings, effective interval and revocation epoch. Validate graph acyclicity, prevent cross-provider implicit inheritance, recheck grants on background execution, and invalidate sessions/streams when a grant is revoked. On transfer of an account, re-evaluate DNS, TLS, secrets, invoices, subscriptions, data export rights and ownership; transfer must not move unrelated historical invoices between merchants.

## 6. Package catalog and capacity management

**Package schema:** `package_id`, `revision`, `provider_org_id`, display name, capability flags, allowed runtimes, maximum sites/domains/subdomains, database/mailbox count, disk bytes, traffic bytes per billing period, CPU share/ceiling, memory max, backup retention, scheduled job limits, feature restrictions, suspension policy, price reference, effective dates and publication status. A reseller sub-package references its provider ceiling and must pass subset checks at creation, modification and provisioning.

**Three distinct resource concepts:** (a) **allocated ceiling**: maximum permitted use; (b) **measured usage**: observed current/typical/peak consumption with window and provenance; (c) **available physical headroom**: host capacity minus measured and safety reserves. Do not conflate these with a guarantee of reserved compute. cgroup v2 `cpu.weight` supplies proportional share; `cpu.max` supplies hard batch caps when configured; `memory.max` supplies memory ceilings. Disk quotas require a verified filesystem enforcement adapter; UI labels quota enforcement accurately when unsupported. The scheduler must evaluate actual host headroom and package policy before placement.

**Admission algorithm:** Validate tenant and subscription; verify package revision and entitlement; detect host capabilities and measured headroom; evaluate dependencies and isolation; reserve quota atomically; create durable provisioning operation; execute typed agent action; verify outcome; commit allocated resource and meter origin; release reservation on confirmed failure. Concurrent requests cannot double-spend a quota. Reconciliation detects host-side drift and expired reservations.

**Package change:** Show old/new ceilings, usage, blocked services, price difference, proration method, affected billing date, potential downtime and any needed migration. Upgrades must not provision until payment/credit policy allows; downgrades with over-limit usage enter an explicit remediation state and must not delete resources silently. Pricing changes do not retroactively rewrite issued invoices.

## 7. Subscription, invoice and payment state machines

### 7.1 Subscription lifecycle

`trial → active → past_due → grace → suspended → cancelled/expired`; policy may allow `past_due → active` after verified payment, `grace → active` after settlement, or a scheduled cancellation at term end. Every transition requires an authorized actor or a verified system trigger, reason, timestamps, entitlement effects and audit. A suspended subscription may preserve control-panel login and data-export rights. Renewal creates a new invoice with a stable period identity; concurrent workers cannot issue duplicates.

### 7.2 Invoice lifecycle

`draft → issued → partially_paid/paid/overdue → voided` subject to jurisdiction and issue-state restrictions. An issued invoice is not edited in place; adjustments use credit notes or approved corrective documents. Tax rounding and line allocations must be deterministic. An invoice's `paid` state is derived from verified settlement allocations net of refunds/chargebacks, not from a client flag.

### 7.3 Payment lifecycle

`initiated → pending → confirmed | failed | expired | reversed`; reversals may occur after confirmation and must create separate ledger movements. Use provider-specific status queries and reconcile ambiguous events. `pending` is a first-class visible state and may last minutes or longer. Do not hard-code Stripe-like timing guarantees into the generic adapter. Handle partial payments, overpayments, duplicate references, callbacks before local response, late confirmations, reversals and manual bank/mobile-money reconciliation.

### 7.4 Payment state transition table

| Current | Trigger | Guard | Next | Financial effect |
|---|---|---|---|---|
| initiated | Adapter accepts request | Stable invoice/merchant reference | pending | None |
| initiated/pending | Verified settlement | Amount, currency, merchant, reference, uniqueness | confirmed | Exactly one settlement posting |
| pending | Verified decline | Provider evidence | failed | No settlement posting |
| pending | Expiry policy | Confirm authoritative status before terminal decision | expired | No settlement posting |
| expired/failed | Late provider confirmation | Evidence plus exception workflow | confirmed or manual_review | Exactly once, never discard funds |
| confirmed | Refund/chargeback/reversal | Verified provider event + authorization | reversed or partially_reversed | Separate compensating posting |
| any | Duplicate callback | Unique provider event/reference | unchanged | No duplicate effect |
| any | Conflicting evidence | Mismatched amount/merchant or state | manual_review | Hold entitlements pending resolution |

## 8. M-Pesa Daraja reference adapter

**Proposed scope:** STK Push, C2B, KES, sandbox and production credentials stored separately; phone-number-first checkout; pending status screen; verified callback ingestion; transaction-status reconciliation where provider APIs and account permissions allow; manual exception queue. A checkout request uses invoice/merchant identifiers, normalized Kenyan MSISDN, amount and reference. Do not expose API credentials, passkeys or full phone numbers in logs or general support screens.

**End-to-end sequence:** (1) server creates a pending payment intent and durable idempotency record; (2) UI requests STK Push and displays a prompt to confirm on the customer's phone; (3) adapter sends request to Daraja; (4) adapter stores request and checkout references; (5) callback is ingested into a restricted endpoint; (6) authenticity and merchant/reference/amount/currency consistency are checked using **actual Daraja-supported controls**, not an invented Stripe-style callback signature; (7) authoritative transaction status is queried when needed/available; (8) settlement posts once to ledger and updates invoice; (9) entitlement/provisioning operations are queued; (10) UI receives status by polling/event stream and shows activation separately.

**Security boundary:** A callback URL or possession of a transaction ID is not proof of payment. Where Daraja does not provide a cryptographic callback signature, combine network/application controls, callback registration, strict correlation, provider-side status verification where supported, anomaly detection and manual review. Do not describe a configured `callback_secret` as a native Daraja signature unless independently supported by the chosen integration mode. Sandbox success does not validate production merchant configuration.

**C2B:** Registration and validation/confirmation URLs must be governed by merchant account capabilities. Unmatched receipts go to suspense/reconciliation, not automatic activation. Manual payment matching requires maker-checker approval above configured thresholds. Repeated STK prompts require rate limits and user consent.

## 9. Stripe second adapter and generic payment interface

The proposed interface includes `createIntent`, `getStatus`, `verifyInboundEvent`, `refund`, `listReconciliationItems`, and `normalizeProviderEvent`; exact signatures belong to Volume 05. Stripe Checkout/PaymentIntents and verified signed webhook events map into the **same** payment/invoice/entitlement transitions as M-Pesa. The adapter may resolve `pending` rapidly but must support delayed methods, disputes and asynchronous outcomes. Store only provider tokens/IDs; no raw card numbers or CVV. Verify signature with the provider's documented signing scheme, reject replay, and scope events to merchant account and environment. A second adapter must pass the same conformance suite as the reference adapter.

## 10. Ledger, tax, pricing and settlement

**Double-entry principle:** Every posted financial event creates a balanced journal within a currency. Example account categories include accounts receivable, gateway clearing, cash/bank, deferred revenue, recognized revenue, tax payable, reseller payable, discounts, refunds and suspense. The actual chart of accounts and revenue recognition require finance approval. Gateway fees are separately booked, not silently subtracted from gross customer payment. No cross-currency journal balancing without explicit FX legs and a documented rate source.

**Invoice calculation pipeline:** Resolve immutable price revision → determine usage/quantity → calculate discounts → calculate taxable base → apply jurisdiction-configured tax rules → round per approved policy → compute line totals and invoice totals → snapshot all inputs → issue invoice with unique sequence. **Kenya reference values mentioned in the source (including 16% VAT)** are proposed configuration examples, not legal advice or proof of current applicability. Validate seller registration status, VAT applicability, eTIMS/KRA electronic invoicing obligations, invoice field requirements and data retention with qualified local professionals before launch. Do not presume that a printable PDF alone establishes tax compliance.

**Proration:** Define `effective_at`, billing period boundaries, timezone, policy (`none`, daily, second-based), credit treatment, rounding and maximum charge. Simulate before confirming. Refund and credit notes never delete the original invoice or payment. Revenue/margin reports identify whether amounts are gross/net of VAT, gateway fees, refunds and wholesale cost.

**Reconciliation:** Run periodic provider settlement matching and daily exception reporting. Match by merchant account, provider transaction reference, invoice/payment reference, amount, currency and event timestamps; never by phone number alone. Flag unmatched, duplicate, amount-mismatch, late, reversed and orphan events. Record reconciliation batch ID, source file/API cursor, counts, totals, exceptions and reviewer decisions. Close periods only after unresolved material discrepancies are approved and tracked.

## 11. Licensing and offline continuity

**Installation identity:** Each deployment has a stable installation ID and public/private identity material held in secure local storage. A license grant binds installation or approved provider scope, licensed capabilities, version policy, validity and renewal terms. Signed lease/grant verification must work offline. Do not encode secrets in license tokens. Reissue/rehost requires documented transfer and abuse controls.

**Proposed state machine:** `unactivated → active → expiring → grace → expired/revoked`; network failures must be distinguished from a verified revocation. Proposed offline grace is **7–14 days**, subject to approval and measured failure tests. Clock skew and rollback are detected with monotonic observations and signed server timestamps where possible; do not brick the panel solely because a wall clock is wrong.

**Degraded license behavior:** During verified license-server outage, continue existing workloads, backups, exports, essential security fixes and administrative access under grace. After grace expiry, restrict paid feature creation according to an approved policy while retaining emergency access and data export. Explicitly test for license service compromise, leaked signing keys, expired grants, repeated offline activation and split-brain renewals. License checks must not be on the critical path of every host operation.

## 12. Reseller wholesale, retail and settlement models

Three models require explicit merchant-of-record selection: **provider-billed customer** (provider invoices customer; reseller earns commission), **reseller-billed customer** (reseller invoices customer; provider invoices reseller wholesale), and **agency/referral** (reseller introduces customer without provisioning/billing authority). Do not mix tax liability, customer support obligation or refund authority across models. Each customer account records its merchant of record, reseller attribution, provider infrastructure owner and support responsibility.

**Pricing:** Wholesale price, retail price, discounts, commissions, tax and gateway fees are separate records with scoped visibility. A reseller may set retail price only within provider policy; provider can enforce floor/ceiling if lawful and approved. Changes to wholesale rates require notice, effective date and handling of existing subscriptions. Commission accrues only on qualifying settled payments; chargebacks reverse or adjust accrued commission, with immutable audit.

**Reseller account limits:** A reseller can create accounts only inside a delegated resource subtree. A delegated package may not exceed provider feature ceilings. An account may be suspended for billing, abuse or infrastructure incidents; reasons and reactivation paths differ. Reseller termination must not silently orphan or purge end-customer data; define transfer-to-provider and retention workflows.

**White-label:** Brand name, logo, colors, support URL, domain and email templates may be configured per reseller. Security notices, legal issuer information, tax seller identity, payment merchant name and audit actor must remain truthful. Custom domains require verified ownership and TLS. White-label cannot hide mandatory security warnings or misrepresent a provider outage.

## 13. Provider fleet and customer operations

**Server lifecycle:** `discovered/enrolling → quarantined → healthy/active → degraded → draining → maintenance → retired`. Enrollment requires authenticated agent identity and explicit provider approval. A node in quarantine, maintenance or draining cannot accept new placements. Provider fleet dashboard shows capability profile, physical CPU/RAM/disk, measured headroom, placement, backups, agent version, security posture, incidents and maintenance windows.

**Placement:** Filter by region/jurisdiction, capacity, service compatibility, runtime support, isolation policy, anti-affinity and provider policy. Reserve capacity atomically and create durable placement operations. Existing workload migrations must consider DNS TTL, TLS, stateful data, downtime, snapshots, rollback and billing continuity. Provider resource pool oversubscription policy must be explicit and cannot imply dedicated guarantees.

**Provider customer lifecycle:** Invitation/signup → organization verification → plan selection → checkout → verified entitlement → account provisioning → onboarding → renewals/usage notifications → support → upgrades/downgrades → grace/suspension → reactivation/export/termination. Every stage must have recoverable partial-failure behavior. Providers can view operational metadata, not tenant secret contents by default.

**Abuse operations:** Complaints and automated alerts create case records; quarantine/suspension requires scoped authority, evidence and reason codes. High-risk interventions are auditable and reversible where safe. Support impersonation is time-bound and visibly indicated; no silent login as a customer. For legal process and data access requests, use a separate approval and evidence chain.

## 14. Commercial user journeys and UI contracts

### Journey C-01: Buy NexusPanel license with M-Pesa
1. User selects eligible plan; UI shows exact currency, taxes, term and license entitlements.
2. User enters phone number and confirms STK prompt request; UI shows `pending` with cancel/help/retry guidance.
3. Callback or reconciliation confirms settlement; UI shows `payment confirmed` but **not** `installation activated` until license grant is issued.
4. Signed license is installed and verified; entitlement cache updates; UI shows activation, expiry and offline-grace policy.
5. If payment is confirmed but license issuance fails, the system retries durable grant issuance and opens an exception; it never asks user to pay again without verification.

### Journey C-02: Reseller creates customer hosting account
1. Reseller chooses delegated package, customer, domain and target region.
2. System previews quota headroom, price, support ownership and capability restrictions.
3. Authorization, reseller delegation, provider ceiling, subscription and host capacity are rechecked server-side.
4. Atomic quota reservation and provisioning operation are created; progress shows queued/running/verifying.
5. Customer receives scoped access only after verification; failed provisioning releases reservation and preserves an actionable error.

### Journey C-03: Overdue invoice and controlled suspension
1. Invoice becomes overdue based on a scheduled, idempotent evaluation.
2. Notifications show amount, merchant, due date, payment methods and grace deadline.
3. At policy threshold, subscription enters grace, then suspension; API access and hosting services change according to an explicit service-by-service matrix.
4. Verified payment restores entitlement; restore operation validates host health and reports any incomplete recovery.
5. Termination and data purge remain separate approved workflows.

### Journey C-04: Package downgrade with excess usage
1. User requests downgrade; system computes current and peak usage and new ceilings.
2. UI identifies resources exceeding new limits and whether service interruption or migration is required.
3. User resolves excess resources or schedules a supported transition; no automatic deletion.
4. Billing/proration preview and authorized confirmation are recorded.
5. Reconciliation verifies resource and invoice states converge.

### Journey C-05: Provider responds to fleet incident
1. Alert identifies affected servers and impacted customer accounts without exposing private content.
2. Operator marks nodes draining, pauses placements and reviews health/backups.
3. Authorized recovery/migration operations execute with tracked progress and rollback.
4. Customer/reseller notices reflect actual service impact and incident timeline.
5. Incident closure requires restoration evidence, financial adjustments if policy warrants and postmortem actions.

## 15. API and event contracts

Volume 05 is authoritative for endpoint names and request/response formats. Commercial APIs include plan/subscription/invoice/payment/license resources, provider/reseller/customer package operations, metering exports, merchant adapter configuration and scoped support actions. **All writes** require organization context, resource policy, entitlements, capability checks where relevant, idempotency for retries and durable audit.

| Command/event | Trigger | Consumer and required behavior |
|---|---|---|
| `invoice.issue` / `invoice.issued` | Renewal, one-off charge | Create immutable invoice and notification once |
| `payment.intent.create` | Customer checkout | Create pending intent; never provision yet |
| `payment.provider_event.received` | Verified/ingested callback | Deduplicate, reconcile and post settlement once |
| `payment.confirmed` | Settlement verified | Update receivable, enqueue entitlement evaluation |
| `entitlement.changed` | Plan/subscription/license transition | Invalidate cache, notify provisioning/policy engine |
| `account.provision.requested` | Authorized order | Reserve quota, select node, execute durable agent job |
| `account.suspended` | Authorized policy transition | Apply documented service-specific suspension behavior |
| `usage.finalized` | Metering period close | Generate reproducible billable usage snapshots |
| `license.revalidated` | Scheduled/online verification | Update signed lease/grace; do not interrupt existing services on transient outage |
| `delegation.revoked` | Provider authorization | Revoke affected sessions/tokens/operations and audit |

**Outbox/inbox:** Persist domain state change and event in one transaction; consumers use stable event IDs and idempotent handlers; retry with backoff and dead-letter exception queues. Exactly-once network delivery is not assumed. Long-running agent operations use correlation IDs and verification before success is displayed.

## 16. Metering, usage accounting and quotas

Meter dimensions include provisioned disk, measured disk consumption, network egress/ingress, account count, domains, sites, databases, mailboxes, backup bytes, compute usage and add-on service units where supported. Every meter declares **unit, source, precision, interval, aggregation, late-arrival policy, billable/nonbillable flag and accuracy class**. Metrics collected for observability are not automatically billing-grade.

**Billing-grade pipeline:** Agent or service emits signed/attributed samples → collector validates tenant/resource/time → deduplicates by meter/source/window → records immutable usage facts → aggregates into a period snapshot → flags gaps/clock skew/resets → reconciles against service and provider counters → freezes period for invoicing → records correction entries separately. Missing samples must not be silently estimated as fact. Estimates require policy, labeling and dispute process. Rate cards reference meter versions and effective dates.

**Quota enforcement:** Enforce count quotas transactionally before resource creation; enforce runtime limits with host adapters only where proven; warn on soft limits; support burst policy explicitly; reconcile actual host state and accounting daily. `used`, `reserved`, `ceiling`, `measured_peak` and `headroom` are separate fields. Traffic cap handling (throttle, charge, suspend or notify) must be explicitly selected and communicated to customers.

## 17. Support, dispute, refunds and operator controls

Support tickets link customer, account, invoice/payment or incident with sensitivity tags and retention policy. Support personnel may see only the metadata needed for the case; credential disclosure and secret retrieval are forbidden by default. Financial disputes require immutable timeline of invoices, checkout requests, provider responses, verified settlement, refunds and entitlement decisions. Refunds require scope and approval thresholds; refund initiation is not the same as confirmed provider refund. Support adjustments and credits are ledger entries with actor/reason.

**Separation of duties:** Merchant credential changes, bulk refunds, payout destination changes, wholesale price changes and mass suspension should require step-up authentication and, for configured thresholds, a second approver. The actor who requests a sensitive change should not approve their own request. Emergency break-glass is time-limited, reason-coded, alerting and post-reviewed.

## 18. Failure catalog and recovery rules

| Failure | Required behavior | Never do |
|---|---|---|
| Daraja callback absent | Keep pending, query status where available, escalate reconciliation | Assume failed or paid from UI timeout |
| Duplicate settlement callback | Return successful acknowledgement, no new journal posting | Credit twice |
| Payment verified but entitlement job failed | Retry idempotent grant/provisioning and show exception | Demand second payment |
| Callback amount/merchant mismatch | Quarantine event, manual review | Activate service |
| License server unreachable | Use signed cached grant and approved grace | Kill hosted workloads immediately |
| Agent unavailable after checkout | Retain order and entitlement, queue provisioning, notify | Show website live |
| Reseller ceiling reduced | Block new allocations, show remediation | Delete existing customer data |
| Customer moves reseller | Reconcile ownership, billing and access | Transfer historical invoice issuer silently |
| Host runs out of disk | Block risky new work, preserve recoverability | Claim package ceiling equals capacity |
| Late usage data | Record correction and review policy | Rewrite issued invoice silently |
| Gateway refund initiated but pending | Show pending refund and reconcile | Mark refunded immediately |
| Reseller branding misconfigured | Fall back to safe verified provider branding | Hide merchant identity or TLS warning |
| Payment database unavailable | Reject/queue safely; preserve provider event evidence | Accept untracked paid status |
| Provider fleet node compromised | Quarantine, rotate agent identity, preserve audit | Continue privileged provisioning |

## 19. Security, privacy and compliance boundaries

Apply Volume 08 threat model: strict tenant scoping, MFA for privileged roles, secret vault references, immutable audit, encryption in transit/at rest, redacted payment event logs, signed provider webhooks where supported, correlation verification where signatures are not available, webhook replay protection and restricted refund authority. PCI scope should be minimized by using hosted card interfaces. Phone numbers, financial identities and invoices require purpose limitation, retention, access/export controls and incident response. Kenyan data protection and tax compliance must be reviewed for the actual merchant/operator structure; this specification does not assert legal certification.

**Threats to test:** reseller horizontal access, forged entitlement, duplicated callback, webhook replay, gateway credential leakage, offline license rollback, package quota race, provider-to-customer secret exposure, forged metering, refund fraud, merchant-account substitution, support-session abuse and node placement into compromised hosts.

## 20. Operational dashboards, reports and SLO proposals

| Dashboard | Minimum indicators | Actionability |
|---|---|---|
| Billing | pending/confirmed/failed payment counts, reconciliation exceptions, overdue aging, tax liabilities | Retry verified status checks, assign exceptions |
| Licensing | active/grace/expired installs, verification failures, signing-key version | Renew/reissue with audit, inspect offline state |
| Reseller | delegated quota used/reserved, customer lifecycle, package uptake, support tickets | Provision, remediate, request capacity |
| Provider fleet | server health, headroom, placements, incident impact, agent versions | Drain, schedule maintenance, restore |
| Metering | missing samples, late events, frozen periods, billable usage by dimension | Reconcile, annotate corrections |
| Finance | gross/net sales, credits, refunds, fees, reseller liabilities, settlement differences | Export auditable journal and exception report |

Proposed SLOs (approval required): 100% duplicate callback deduplication in conformance tests; zero unbalanced ledger postings; 100% authorization and delegation checks on protected actions; zero irreversible data deletions caused solely by billing state transitions; all pending payments older than configured threshold appear in an exception queue; reconciliation runs on an approved cadence with tracked backlog; all provider-impacting operations produce an audit correlation ID. Performance/availability targets must be measured on representative Nano, Standard and Provider installations before public claims.

## 21. Acceptance and conformance test matrix

| ID | Scenario | Pass criteria |
|---|---|---|
| COM-001 | Two identical M-Pesa callbacks | One settlement posting and one entitlement event |
| COM-002 | Callback arrives before checkout response | Correlated payment reaches correct terminal state |
| COM-003 | Forged paid browser redirect | No entitlement activation |
| COM-004 | Incorrect callback amount | Manual review, no automatic grant |
| COM-005 | Late confirmation after timeout | Corrected payment state without double charge |
| COM-006 | Payment confirmed, agent offline | Order retained, provisioning queued, status truthful |
| COM-007 | Refund and chargeback | Reversal journal balances and original preserved |
| COM-008 | Reconcile provider statement | All unmatched transactions tracked with provenance |
| COM-009 | Two simultaneous renewals | One invoice per subscription/period |
| COM-010 | Two concurrent quota allocations | Cannot exceed hard ceiling |
| COM-011 | Reseller exceeds delegated limit | Denied before reservation or agent call |
| COM-012 | Revoked reseller token | Denied immediately on new requests/streams |
| COM-013 | Provider support attempts secret access | Denied without explicit policy and audit |
| COM-014 | Package downgrade with excess resources | Blocked or staged; no silent deletion |
| COM-015 | Suspension after past due | Exactly documented service behavior; data preserved |
| COM-016 | Payment during grace | Reactivation job idempotent and verifiable |
| COM-017 | License server down 10 days | Approved offline grace behavior demonstrated |
| COM-018 | Signed license modified | Invalid grant rejected without data destruction |
| COM-019 | Host capability shrinks | New allocations blocked; honest delta surfaced |
| COM-020 | Nano purchases unsupported feature | Checkout/entitlement messaging discloses host limitation; no impossible provision |
| COM-021 | Invoice tax rules change | Old issued invoices remain immutable |
| COM-022 | Duplicate usage samples | Metered totals unchanged |
| COM-023 | Missing usage interval | Explicit gap/exception, no silent invented reading |
| COM-024 | Account ownership transfer | DNS, billing, access and audit reconciled |
| COM-025 | Reseller white-label domain | Verified ownership, valid TLS, truthful merchant identity |
| COM-026 | Provider node quarantine | New placements denied and current incidents visible |
| COM-027 | Agent action partial failure | Retry/rollback state consistent with billing order |
| COM-028 | Merchant credentials rotated | Old credentials revoked, no secret leaks in logs |
| COM-029 | Manual refund above threshold | Second approver required; requester cannot self-approve |
| COM-030 | SQLite/PostgreSQL parity | Identical commercial transitions and uniqueness semantics |
| COM-031 | Outbox replay after crash | Exactly-once business effect, at-least-once event delivery |
| COM-032 | Dispute export | Redacted, scoped, complete immutable evidence chain |
| COM-033 | Cross-tenant invoice enumeration | No existence leakage or unauthorized records |
| COM-034 | Settlement currency mismatch | Quarantined and reconciled, no auto-conversion |
| COM-035 | Provider bulk suspend | Scoped, approved, audit logged and reversible as policy permits |
| COM-036 | Data retention on termination | Export/hold/purge sequence followed with evidence |

## 22. Decision register and sign-off gates

| ID | Decision required | Default proposal / validation |
|---|---|---|
| COM-D01 | Legal merchant of record | Explicit per provider/reseller arrangement; legal review |
| COM-D02 | Kenya tax and e-invoicing obligations | External qualified review before production |
| COM-D03 | M-Pesa merchant configuration | Daraja sandbox/production credentials and live callbacks tested |
| COM-D04 | Payment settlement authority | Provider-confirmed evidence and reconciliation; no callback-only trust |
| COM-D05 | License offline grace | 7–14 days proposed; approval and outage tests |
| COM-D06 | License expiry feature restrictions | Never destroy data; define safe admin/export access |
| COM-D07 | Invoice numbering, credit notes, rounding | Finance/tax policy approved and tested |
| COM-D08 | Proration and refund schedule | Published customer terms and deterministic calculation |
| COM-D09 | Wholesale vs reseller merchant models | Contractual liability and visibility explicitly assigned |
| COM-D10 | Hard/soft quotas and traffic overages | Enforceability benchmarked per host adapter |
| COM-D11 | Support impersonation policy | Consent, scope, TTL and audit requirements approved |
| COM-D12 | Metering accuracy for billing | Source/precision/retention and dispute thresholds validated |
| COM-D13 | Provider oversubscription | Transparent non-guaranteed CPU semantics and risk budgets |
| COM-D14 | Stripe rollout | Second adapter conformance before enabling |
| COM-D15 | Currency/FX handling | No implicit FX; approved rates and settlement rules |
| COM-D16 | Multi-node provider placement | Remote agent mTLS and failure tests in v2+ |
| COM-D17 | Reseller branding/legal disclosures | Verified domain, seller identity and customer notices |
| COM-D18 | Data retention on cancellation | Policy, legal holds, export, purge and backup expiry |
| COM-D19 | Financial report close process | Reviewer, audit trail, reconciliation exceptions |
| COM-D20 | Production release | Security, accounting, usability and end-to-end sign-off |

## 23. Cross-volume consistency and implementation handoff

- **Volume 01:** Requirements `NP-ACC-001..015`, `NP-BILL-001..017` and all platform area scope remain authoritative; this volume elaborates them, not overrides them.
- **Volume 02:** Guided router for lay users; reseller delegated dashboard and provider fleet dashboard; role-aware views, payment `pending`, distinct capability/permission/offline/degraded states.
- **Volume 03:** Go control plane, durable execution, host capacity profiles, local Unix-socket agent v1, remote mTLS v2+, HA v3+.
- **Volume 04:** Commercial entities, immutable ledger, tenancy, SQLite/PostgreSQL parity, migrations and transactional outbox.
- **Volume 05:** REST/RPC endpoints, gateway adapters, webhook ingress, idempotency, scoped authorization, event contracts.
- **Volume 06:** Host provisioning, package quota enforcement, service-specific suspension and resource measurements.
- **Volume 07:** ServiceTemplate feature entitlements, storage/queue/search consumption and self-hosted service billing dimensions.
- **Volume 08:** Payment security, signed license verification, least privilege, privacy, audit, incident response and recovery.
- **Volume 10:** Implementation sequence, ownership, testing fixtures, measured resource benchmarks, staged rollouts, financial reconciliation drills and launch evidence.

**Exit criterion:** Volume 09 is documentation-complete when all 36 acceptance tests have an implementation owner and test fixture plan, the 20 decisions have approvers, API/data mappings are reconciled against Volumes 04–05, and finance/legal/security reviews have signed off applicable production policies. Documentation completion does not mean product implementation or tax compliance has been verified.
