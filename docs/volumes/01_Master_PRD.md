# NexusPanel — Master Product Requirements Document

**Volume:** 01 of 10 — Product Vision, Scope & Master PRD  
**Edition:** Consolidated baseline 1.0 (based on uploaded v2.0 source set)  
**Date:** 6 November 2026  
**Status:** **Consolidated draft — not approved for implementation**  
**Product name:** NexusPanel (working title)  
**Owner:** Product owner, with engineering/security approval required for technical gates

> **Authority and intent.** This is the canonical proposed **product behavior contract** for the 29-area NexusPanel vision. It consolidates the uploaded Master Blueprint v2.0, Stage 2 PRD v2.0, the authoritative 29-area scope register v2.0, and the compatible decisions of Stages 3–7. Where a source expresses an unapproved decision, this document labels it **proposed**; it does not invent sign-off. The nine companion volumes will define UX, implementation architecture, data, APIs, host operations, platform services, security, commerce, and testing. All original NP-* requirement IDs are retained wherever possible.

## Reader's guide and governance

- **MUST** means a requirement that must be demonstrated before its assigned module is declared production-ready; it does **not** mean every module ships in v1 or runs on Nano.
- **SHOULD** requires a documented rationale for exceptions. **MAY** is optional.
- **SPECIFIED** means drafted, **APPROVED** means product-owner sign-off, **VERIFIED** means supported by tests/evidence, and **RELEASED** means deployed and supported. None of these documents alone establishes VERIFIED or RELEASED.
- **Source precedence:** Approved decisions and this volume's explicit resolution rules take precedence once approved; until then, material disagreements remain **open gates**. The original v2.0 Master Blueprint is the reference for intent, not an automatic substitute for approval.
- **Traceability:** Every `NP-*` requirement shall map to area ID (01–29), release phase, UX flow, API/data/agent contract where applicable, authorization, tests, and evidence. No silent renumbering.
- **Product scope vs installed capability:** All 29 areas are in product scope. A particular server enables only the subset supported by resources, topology, permissions, policy and release phase.
- **No infrastructure authorization:** This document does not authorize modifying or deploying to the Leaseweb VPS.

## Executive product contract

NexusPanel is an integrated, AI-assisted hosting and infrastructure management platform for personal VPS owners, hosting customers, developers, resellers, and commercial providers. It offers beginner-first guided onboarding without sacrificing advanced controls, and a coherent control plane for websites, apps, databases, DNS, email, files, deployments, backups, security, storage services, self-hosted software and provider operations. The same resource ownership and policy engine governs UI, API, CLI, automation, extensions and AI-initiated proposals.

**Four independent product spines:** (1) Commercial: tenants, plans, entitlements, invoicing and payment; (2) Service: catalog, templates, dependencies and lifecycle; (3) Execution: durable operations, policy, typed agent, verification and recovery; (4) Surface: guided journeys, portfolios, provider/reseller workspaces and APIs. The v1 architecture must leave explicit interfaces for all four even when a feature is scheduled for v2/v3.

**Primary first-run experience:** “What do you want to run?” The system offers a guided website, app, self-hosted service, storage, infrastructure or exploration journey according to detected capabilities and user permissions. Returning users see their service portfolio. Navigation menus remain secondary, not the default beginner experience.

**Deployment ladder:** Nano (proposed 512 MB class; static hosting and lightweight control), Starter (selected local services), Standard (Docker-backed apps/services), Provider (fleet/reseller/commercial). These are capability profiles, not universal performance guarantees. A 512 MB machine must not be represented as able to host Docker, databases, self-hosted mail and multiple dynamic apps simultaneously. Fractional-vCPU and 10 GB storage hosts require explicit install preflights and benchmark evidence.

**Technology direction (proposed, approval-gated):** Go control plane, Go typed privileged agent, static SPA, SQLite metadata for Nano and PostgreSQL for Standard+, cgroup v2 controls, Docker isolation where feasible, and signed/versioned ServiceTemplates. These are implementation constraints derived from the source set, not proof that the 512 MB footprint has been achieved.

## Complete authoritative platform scope (29 areas)

The table below preserves the uploaded scope register's delivery and Nano mapping. `v1`, `v2`, `v3`, `v4` are **proposed release phases**, not dates or approvals. The original register's documentation-state field is a historical review judgment, not a production readiness certification.

| # | Platform area | Doc state | Capability | Docker | Phase | Nano |
|---|---|---|---|---|---|---|
| 01 | Dashboard | Foundation (draft) | infrastructure | n/a | v1 | available |
| 02 | Server management | Partial | infrastructure | n/a | v1 | available (read-only) |
| 03 | Websites | Foundation | serving | container-preferred | v1 | available (static only) |
| 04 | DNS | Foundation → Partial (v1 is guidance) | infrastructure | n/a | v1 (guidance) / v2 (authoritative) | available (guidance) |
| 05 | SSL/TLS | Foundation (draft) | infrastructure | n/a | v1 (HTTP-01) / v2 (DNS-01, wildcard) | available |
| 06 | File manager | Foundation | n/a | container-preferred | v1 | available |
| 07 | Application manager | Foundation (draft) | serving (+ batch build) | container-only | v1 | unavailable |
| 08 | Runtime manager | Foundation (draft) | serving | container-only | v1 | unavailable |
| 09 | Databases | Foundation | serving | container-only | v2 | unavailable |
| 10 | Email hosting | Partial | serving | container-preferred | v1 (external SMTP) / v2 (self-hosted) | unavailable |
| 11 | Git and deployments | Foundation (draft) | serving + batch | container-only | v1 | unavailable |
| 12 | Backups | Foundation | batch | n/a | v1 | available (remote-only) |
| 13 | Security | Foundation | infrastructure | container-preferred | v1 | available |
| 14 | Monitoring | Foundation (draft) | infrastructure | n/a | v1 | available (minimal) |
| 15 | Terminal/developer tools | Foundation (draft) | infrastructure | container-only | v1 | unavailable (terminal), available (SFTP) |
| 16 | User/account management | Foundation | infrastructure | n/a | v1 | available |
| 17 | AI assistant | Foundation (draft) | batch | n/a | v1 (read-only) / v4 (actions) | available (read-only) |
| 18 | Automation/API | Foundation | serving + batch | n/a | v1 | partial |
| 19 | Panel administration | Foundation (draft) | infrastructure | n/a | v1 | available |
| 20 | Extensions | Foundation (draft) | n/a | n/a | v2 | unavailable |
| 21 | Commerce/reseller | Foundation (draft) | infrastructure | n/a | v1 (M-Pesa) / v2 (reseller) / v3 (full) | unavailable |
| 22 | Migration/staging | Major gap → Partial | batch | container-preferred | v1 (VPS-to-VPS) / v2 (cPanel/DA) | partial |
| 23 | Containers/orchestration | Foundation (draft) | infrastructure | **internal primitive** | v1 (as primitive) / v2 (user-facing) | unavailable |
| 24 | Performance/caching | Major gap → Partial | serving + batch | container-preferred | v1 (basic) / v2 (full) | unavailable |
| 25 | Email security/deliverability | Major gap → Partial | serving | container-preferred | v2 | unavailable |
| 26 | Analytics/site quality | Major gap → Partial | batch | n/a | v2 | unavailable |
| 27 | Audit/compliance/support | Foundation (draft) | infrastructure | n/a | v1 | available |
| **28** | **Platform Services** *(new)* | **Foundation (v2.0)** | **serving + batch** | **container-only** | **v1 (storage/transform/delivery) / v2 (cache/queue/search)** | **unavailable** |
| **29** | **Self-hosted Templates** *(new)* | **Foundation (v2.0)** | **serving** | **container-only** | **v1 (5–10 apps) / v2 (20–30 apps)** | **unavailable** |

### Feature-completeness rule

An area is **specified** only when its requirements, ownership, permissions, capability constraints, failure states, UX states, lifecycle and acceptance criteria are recorded. It is **implemented** only after the corresponding contracts and code exist. It is **release-ready** only after security, operational, recovery, resource-budget and E2E evidence pass. The 29-area scope is never equated with installing all 29 services on a small VPS.

## Cross-volume contract and acceptance discipline

| Volume | Authority and content | Requirement links |
|---|---|---|
| 01 (this volume) | Product behavior, area scope, release gates, business rules, acceptance outcomes | `NP-*` |
| 02 | Guided router, screens, interactions, accessibility, mobile and journeys | `UX-*`, `FLOW-*` → `NP-*` |
| 03 | Control plane, topology, capabilities, cgroups, Docker, service boundaries | `ARCH-*`, `ADR-*` → `NP-*` |
| 04 | Logical entities, tenant isolation, dual storage backends, migration and integrity | Entities → `NP-*` |
| 05 | REST/API, agent RPC, webhooks, payment callbacks, automation and errors | Operations/endpoints → `NP-*` |
| 06 | Typed Go agent, host adapters, preflight, execution, reconciliation, rollback | `AGT-*` → `NP-*` |
| 07 | ServiceTemplates, platform services, catalog, signed artifacts, lifecycle exceptions | Template contracts → `NP-*` |
| 08 | Security, privacy, abuse, audit integrity, incident response and recovery | Controls/tests → `NP-*` |
| 09 | Commercial ledger, entitlements, M-Pesa, resellers and provider fleet operations | Business invariants → `NP-*` |
| 10 | Delivery phases, benchmarks, test plans, acceptance gates and runbooks | Evidence → `NP-*` |

## Cross-cutting normative clarifications

1. **Router ownership:** The backend is authoritative for permissions and capability data; the frontend may compose the landing experience from a versioned, server-supplied persona/surface matrix. `GET /` is the browser entry route, not an implicit public API authorization endpoint. The exact wire contract belongs to Volume 05.
2. **Data engines:** SQLite on Nano and PostgreSQL on Standard+ are distinct supported deployments of one logical management model. PostgreSQL-only behavior cannot be a hidden prerequisite for Nano. Cross-engine migration and parity tests are required; automatic migration is not presumed safe.
3. **ServiceTemplates:** A common discovery, provisioning and instance model is mandatory where applicable; databases, stateful storage, email and security-sensitive services retain typed domain-specific lifecycle and recovery rules. A template record alone is not permission to run arbitrary OCI code.
4. **Resource ceilings:** cgroup v2 `cpu.weight` enables relative fairness; `memory.max` is a limit, not a guarantee. `cpu.max` may be used for bounded batch work. Host/kernel and container compatibility must be verified. Disk quotas and bandwidth controls require separate enforcement mechanisms.
5. **Nano support:** The proposed minimum is a **validation target**, not a proven production floor. Benchmark total resident memory, installation peak, update peak, recovery headroom, disk footprint and CPU pressure. If a profile fails, the installer must decline or offer remote control-plane deployment.
6. **Payment integrity:** Asynchronous state, deduplication, reconciliation and ledger consistency are mandatory. Provider-specific webhook authenticity must be verified against actual gateway capabilities; do not assume all gateways provide equivalent signed callbacks. Financial and tax configuration requires jurisdiction-specific review.
7. **Agent trust:** v1 local agent communication uses Unix sockets and OS peer-credential validation as proposed; remote agent mTLS and key lifecycle are later-phase contracts. No browser, plugin or AI prompt may directly invoke unrestricted root commands.
8. **AI boundary:** Read-only diagnostics and explicit proposed actions are v1 scope; no autonomous privileged execution in v1. Human approval does not replace deterministic policy authorization or agent validation.
9. **Immutable audit:** Tamper-evident hash chaining is a design requirement, not a guarantee of immutability against an administrator who controls the database and keys. Off-host anchoring, retention and access controls must be specified and tested.
10. **Release interpretation:** v1 includes a production-quality foundation and explicitly assigned feature slices; complete 29-area coverage is the long-term product scope. A deferred area remains in the specification but cannot be advertised as shipped.

---

# Detailed functional and nonfunctional requirements

The following requirements preserve the substantive Stage 2 v2.0 functional contract and its stable `NP-*` identifiers, while eliminating editorial v1-to-v2 provenance markers. Where an earlier section describes a proposed implementation choice as final, the clarifications above and the decision gates at the end of this volume govern.

## 0. How to read and govern this PRD

This is the *product behavior contract*, not a database schema, API OpenAPI file, or implementation runbook. It turns the master blueprint into traceable functional and nonfunctional requirements. An item marked **MUST** is mandatory for its module's production completion; **SHOULD** is strongly preferred; **MAY** is optional. A complete product requirement does **not** imply that its service must be installed on every VPS. A module may be disabled because of topology, resources, permissions, or policy, with a clear explanation and a supported alternative.

**Statuses:** `SPECIFIED` = intended requirement awaiting approval; `DECISION` = unresolved policy/architecture choice; `DEPENDENCY` = requires another capability; `RELEASE GATE` = cannot claim production readiness without evidence. All requirements below are `SPECIFIED` unless otherwise noted. Proposed defaults are recommendations, not agreed facts.

**Traceability convention:** `NP-<DOMAIN>-###` identifies a stable requirement. Each feature must eventually map to a UI screen, API contract, authorization rule, job/agent operation (if applicable), test case, telemetry and release evidence. Changing a requirement requires an impact review.

**Cross-cutting transaction lifecycle:** Discover/collect → validate → preview impact → authorize → approve if high risk → enqueue durable operation → execute least-privileged typed agent action → verify resulting state → audit and notify. Failed operations expose status, retry eligibility, compensating actions and safe recovery guidance.


- **Capability gate.** Every module declares its `floor` / `typical` / `ceiling` per resource, its workload category (serving / batch / infrastructure), and its prerequisites. The installer (§5.1) evaluates these against physical capacity and produces an honest capability set. A module the installer cannot enable is displayed as unavailable with a reason — not greyed out silently.
- **Allocation is a ceiling, not a reservation.** Every workload uses cgroup v2 `cpu.weight` for fairness under contention and `memory.max` for the RAM ceiling. `cpu.max` is set only on batch workloads. `cpuset` (core pinning) is never used. Idle CPU is never wasted.
- **Docker is the primary isolation primitive** where Docker fits (Standard / Provider profiles). OS-user-per-account is the fallback for workloads that cannot containerize (Nano static sites).
- **Guided-default rule.** `GET /` resolves the default landing page per persona. The guided router is the primary IA for personas whose job is "get this working." The menu is the secondary IA. This is an architectural constraint, not a UI preference.
- **AI is never the authority.** The policy engine is deterministic and separate from the LLM. AI proposes; policy disposes. v1 ships read-only diagnostics + proposed-action approval; autonomous execution is v4+.

## 1. Product objectives, boundaries, and personas

### 1.1 Objectives

- **NP-STR-001:** Provide a safe, understandable first-website/app deployment journey for a non-expert VPS owner.
- **NP-STR-002:** Provide complete hosting-account and reseller operations without conflating a customer account with a Linux root account.
- **NP-STR-003:** Enable providers to operate a fleet of servers, service roles, packages, quotas, subscriptions, incidents and delegated support.
- **NP-STR-004:** Offer the same authorized capabilities through UI, versioned API, CLI and AI tools; no special AI bypass.
- **NP-STR-005:** Run control plane and managed agent on one machine for a standalone installation; support remote agents and role-specific nodes without rewriting domain contracts.
- **NP-STR-006:** Treat security, backups, upgrades, disaster recovery, and auditable change management as product features.
- **NP-STR-007:** Support an extensible ecosystem of adapters, integrations, webhooks and provider-approved plugins.


- **NP-STR-008:** Support VPS infrastructure from 512 MB to multi-node clusters via a capability ladder (Nano / Starter / Standard / Provider) with honest detection at install time and honest re-evaluation on hardware change.
- **NP-STR-009:** Provide a fully commercial commercial spine (billing, licensing, reseller, provider) from v1 with a **pluggable, async-first payment adapter interface**, using M-Pesa as the reference implementation and Stripe as the second.
- **NP-STR-010:** Treat the guided router as the primary information architecture for personas whose job is "get something new working" — not a presentation mode, not a later feature, but the first release gate.
- **NP-STR-011:** Unify storage, image transformation, delivery, self-hosted apps, cache, queue, and search under a single `ServiceTemplate` primitive so that adding a new self-hostable service is a data operation, not a code change.

### 1.2 Personas and jobs-to-be-done


| Persona | Primary jobs | Typical scope |
|---|---|---|
| Solo VPS owner | Install, secure, create sites/apps, manage databases/email, diagnose issues | Owned installation and resources |
| Hosting customer | Operate assigned sites, domains, files, databases, mail, backups | One or more hosting accounts |
| Developer/team member | Deploy, view logs, manage authorized resources and secrets | Delegated resource scopes |
| Reseller | Create packages/accounts, onboard customers, brand portal, support clients | Delegated subtree and commercial limits |
| Provider operator | Manage fleet, placement, plans, abuse, health, upgrades, disaster recovery | Provider organization and permitted servers |
| Billing/support/auditor | Invoices, support actions, compliance evidence | Least-privilege functional scopes |


| Job-to-be-done | What they want | Default landing |
|---|---|---|
| Website owner | Get a site online, securely | Website flow |
| App developer | Deploy a backend, see logs | Application flow |
| Self-hoster | Run Nextcloud / Gitea / n8n / my own thing | Self-hosted template flow |
| Storage user | Store photos/files, serve images | Storage flow |
| Infrastructure user | VPN, proxy, mail, DNS | Infrastructure flow |
| Explorer | "I bought a VPS and I'm not sure what to do" | Guided tour |
| Reseller | Sell hosting to others | Delegated dashboard |
| Provider admin | Run the fleet | Fleet dashboard |

### 1.3 Out of scope for this PRD (not necessarily the product)

Exact SQL DDL, concrete API payload schemas, pixel-perfect screens, provider-specific payment adapters, installer commands, package compatibility evidence, and production capacity benchmarks belong to later artifacts. No requirement grants permission to modify the user's current VPS.


## 2. Global rules and system-wide invariants

- **NP-CORE-001:** Every resource has an immutable identifier, owning organization/account, lifecycle state, creator and audit metadata; transfers require explicit authorization.
- **NP-CORE-002:** Authorization evaluates actor, role, organization, delegated scope, resource, action, plan entitlement, server capability, policy and relevant approvals. Deny by default.
- **NP-CORE-003:** Resource provisioning uses idempotent, durable jobs with unique operation IDs; refreshes or retries cannot create duplicates.
- **NP-CORE-004:** Every mutating action has preflight checks, a human-readable impact preview, structured result and auditable record.
- **NP-CORE-005:** Sensitive values are masked in UI, logs, events, exports, errors, AI context and audit; secrets may be revealed only via explicitly authorized workflows with expiration where feasible.
- **NP-CORE-006:** A server, service or feature in degraded/offline state cannot be shown as healthy based only on a submitted request.
- **NP-CORE-007:** Destructive actions distinguish reversible disable/suspend, soft delete with retention, and irreversible purge; confirmation describes actual data impact.
- **NP-CORE-008:** Quota enforcement is concurrency-safe and covers reserved, pending and active usage; over-quota operations return actionable errors.
- **NP-CORE-009:** Multi-step changes expose partial success and compensating recovery; never claim atomicity across external services unless guaranteed.
- **NP-CORE-010:** Resource names and paths are validated and normalized consistently; avoid cross-tenant name collisions and unsafe Unicode/path ambiguities.
- **NP-CORE-011:** Modules declare compatibility, dependencies, resource budgets, health checks, ownership of configuration and safe uninstall semantics.
- **NP-CORE-012:** All timestamps use UTC in storage and user-configured time zones for presentation; audit records retain UTC.
- **NP-CORE-013:** The UI provides consistent empty/loading/success/error/partial/permission-denied/quota-exceeded/offline states.
- **NP-CORE-014:** Admin and agent APIs are versioned; breaking changes require compatibility and migration policies.
- **NP-CORE-015:** Audit and critical operation logs survive normal tenant resource deletion under defined retention policies.
- **NP-CORE-016:** Every service has a declared backup and recovery policy or an explicit statement that it is ephemeral.


- **NP-CORE-017:** Every module declares a **capability record** with `floor` / `typical` / `ceiling` per resource, a workload category (serving / batch / infrastructure), and prerequisites. The installer reads this from a data file, not from code. See §5.1.
- **NP-CORE-018:** **Allocation is a ceiling, not a reservation.** Workloads use cgroup v2 `cpu.weight` for fair-share CPU and `memory.max` for the RAM ceiling. `cpu.max` is set only on batch workloads. `cpuset` is never used. RAM is never overcommitted beyond 1.2× ceiling; CPU may overcommit 2–4× typical.
- **NP-CORE-019:** Batch workloads (image transforms, backups, builds, AI diagnostics, migrations) run at low `cpu.weight` and yield to serving workloads under contention. They are queued, not rejected.
- **NP-CORE-020:** The **guided router** (`GET /`) is the default landing page for personas whose job is "get something new working." The menu is the secondary IA. Both share identical permissions; the difference is presentation, not authorization.
- **NP-CORE-021:** Every workload's cgroup metrics are polled, aggregated, and surfaced as current / typical (7-day p50) / peak (24h p99). Allocation is separate from typical usage and separate from peak. The dashboard shows all three.
- **NP-CORE-022:** The audit chain starts at v1 and is hash-chained. It cannot be added later.
- **NP-CORE-023:** The payment adapter interface is **async-first**. Every payment is a state machine with a `pending` state. Card gateways resolve `pending` in <1s; mobile money resolves it in 10s–2min. Both fit the same interface.
- **NP-CORE-024:** Docker is the primary isolation primitive where Docker fits. OS-user-per-account is the fallback for workloads that cannot containerize.
- **NP-CORE-025:** Every service is expressible as a `ServiceTemplate` record (identity, source, resources, network, credentials, data, health, configuration, metering, export) unless there is a documented reason it cannot be. A bespoke service that cannot be templated is a design smell requiring review.

## 3. Identity, onboarding, organizations, and permissions

**Goal:** One identity model for standalone operators, hosting customers, resellers and providers.

- **NP-IAM-001:** Support first-owner bootstrap with one-time enrollment, secure session issuance, and a recovery path that does not depend on an already healthy UI.
- **NP-IAM-002:** Support users, organizations, memberships, invitations, role assignment, revocation, and membership expiration.
- **NP-IAM-003:** Support password and passkey/TOTP MFA strategies according to approved identity design; require MFA for privileged provider/reseller actions in production.
- **NP-IAM-004:** Support session inventory, device/session revocation, idle and absolute expiry, account lockout/rate limits, and suspicious-login alerts.
- **NP-IAM-005:** Provide fine-grained permissions by resource and action (view/create/edit/deploy/delete/export/approve/impersonate/billing).
- **NP-IAM-006:** A user may belong to multiple organizations; context switching must not leak resources or cached data.
- **NP-IAM-007:** Delegation from provider to reseller is explicitly scoped, revocable, auditable, and cannot exceed grantor privileges.
- **NP-IAM-008:** Support support-access grants with justification, expiry, resource limits, visible access banner and full audit. No silent impersonation.
- **NP-IAM-009:** Expose access-denied reasons without revealing another tenant's resource existence.
- **NP-IAM-010:** API tokens/service identities support least-privilege scopes, expiry, rotation, last-used metadata and revocation.

**Acceptance scenarios:** (a) New standalone owner enrolls and signs in with MFA; (b) reseller cannot see unrelated provider customers; (c) revoked team member loses access promptly, including websocket/stream subscriptions; (d) support access expires and is audited.


## 4. Information architecture, dashboard and universal UX

- **NP-UX-001:** Role-aware navigation covers Overview, Websites, Applications, Domains/DNS, SSL, Email, Databases, Files, Backups, Security, Monitoring, Automation, Accounts, Resellers, Billing, Servers, Settings and AI.
- **NP-UX-002:** Global search and command palette find permitted resources, actions and documentation without leaking restricted resources.
- **NP-UX-003:** Overview summarizes health, pending operations, active incidents, capacity, usage, expiring certificates, failed backups and actionable recommendations.
- **NP-UX-004:** Onboarding checklist adapts to persona and installation topology; it never suggests unsafe operations as mandatory shortcuts.
- **NP-UX-005:** Wizards show progress, dependencies, preflight failures, estimated impact, irreversible consequences and recoverability.
- **NP-UX-006:** Job/activity center shows queued/running/waiting-for-approval/failed/partially-succeeded/succeeded/cancelled states and links to logs.
- **NP-UX-007:** Support accessible keyboard operation, screen-reader semantics, visible focus, color contrast and WCAG 2.2 AA as target.
- **NP-UX-008:** Support responsive desktop/tablet/mobile layouts, light/dark themes, time zones, localization-ready copy and consistent formatting.
- **NP-UX-009:** Show a clear novice/expert distinction through progressive disclosure, not different security rules.
- **NP-UX-010:** Unsaved changes, stale data, concurrent edits and disconnections are communicated clearly.

**Acceptance scenarios:** A novice can complete a guided website journey without CLI use after installation; expert users can inspect exact configuration implications; keyboard-only user can finish a provisioning wizard.


- **NP-UX-011:** **`GET /` resolves the default landing page per persona** (jobs-to-be-done, §1.2). The guided router is the primary IA for lay persons, the delegated dashboard for resellers, and the fleet dashboard for provider admins. The menu (NP-UX-001) is the **secondary** IA, reachable from "Advanced tools" on the router, or from the portfolio view.
- **NP-UX-012:** The router is a single screen with six cards: **A website / An application / A self-hosted app / Data & storage / Infrastructure / I'm not sure yet**. Each card routes to its own guided flow with its own "you're live" moment.
- **NP-UX-013:** The **service portfolio** replaces the "Overview dashboard" as the primary surface for personas whose job is "manage what I already have." It lists everything the user runs (websites, applications, self-hosted services, storage buckets, databases), each with health, usage, and quick actions. This is closer to a projects list than to a tools menu.
- **NP-UX-014:** **Guided / Standard / Expert are presentation modes**, not separate permission systems. Permissions are identical across modes. The persona determines the default landing, but the user can always switch modes.
- **NP-UX-015:** The **allocation-vs-usage view** is a first-class dashboard element. Each workload shows allocated (ceiling), typical (7-day p50), peak (24h p99), and current usage. Server-level view shows aggregate across all workloads, with headroom.
- **NP-UX-016:** The **capability set** is displayed on the overview for installations where any capability is unavailable: "Your VPS cannot run: WordPress (needs 150 MB more), Databases (needs 250 MB more)..." with reasons and an upgrade link. This is the Nano layman experience (§9.3).

**Additional acceptance scenarios:**

- **NP-UX-017:** A lay person landing on `GET /` for the first time sees the router, not the menu. The menu is reachable via "Advanced tools" but is not the default.
- **NP-UX-018:** A user who switches from guided mode to expert mode sees the same resources with identical permissions; the only difference is presentation.
- **NP-UX-019:** When the installer determines a capability is unavailable, the corresponding menu entry shows an honest explanation with the required resource delta — not a greyed-out button.

## 5. Servers, fleet inventory, agents, and infrastructure

- **NP-INF-001:** Enroll standalone/local or remote servers using short-lived enrollment credentials and explicit operator approval.
- **NP-INF-002:** Record OS, architecture, kernel, IPs, hostname, resources, roles, service inventory, agent version and capabilities.
- **NP-INF-003:** Continuously report heartbeat, health, capacity, connectivity and version drift with clear freshness timestamps.
- **NP-INF-004:** Assign server roles (control plane, web/app, database, mail, DNS primary/secondary, backup, monitoring) with compatibility checks.
- **NP-INF-005:** Provide placement policies based on available capacity, supported services, tenant restrictions, location and failure domains.
- **NP-INF-006:** Allow controlled cordon/drain/maintenance, workload migration where supported, and agent upgrade with rollback strategy.
- **NP-INF-007:** Agent exposes only authenticated, typed, allowlisted operations with scoped resources and argument validation; never a generic root command endpoint.
- **NP-INF-008:** Agent supports secure identity rotation, lost-agent recovery, command replay protection and bounded execution.
- **NP-INF-009:** A failed/unreachable server queues or rejects operations according to explicit policy; no false-positive provisioning status.
- **NP-INF-010:** Service installation/uninstallation is module-aware, checks conflicts and dependencies, and preserves unrelated operator-managed configuration.
- **NP-INF-011:** Fleet actions require scope and approval proportional to blast radius; show affected servers/accounts.
- **NP-INF-012:** Provide inventory export and topology view; surface unsupported OS/service combinations.

**Acceptance scenarios:** Single VPS operates without a separate broker/Kubernetes; remote server can enroll securely; compromised/expired enrollment token cannot add a node; maintenance mode prevents new placements.


- **NP-INF-013:** The installer runs **five-step capability detection** (detect hardware and blockers → floor check → typical fit → CPU overcommit → ceiling fit → disk fit) and produces an honest capability set. The installer logic is in Master Blueprint §9.1.
- **NP-INF-014:** The **capability record** for each module is a data file (YAML or JSON), not code. It declares `floor` / `typical` / `ceiling` per resource, workload category, and prerequisites. The installer reads it, evaluates it, and displays the result.
- **NP-INF-015:** The installer **refuses below floor** (RAM below ~380 MB panel+OS baseline; disk below ~3 GB). Every other constraint produces a warning and an honest capability set. Refusal reasons are surfaced to the user, not logged silently.
- **NP-INF-016:** The installer detects **non-resource blockers**: port 25 blocked → email unavailable regardless of RAM; no public IP → hosting impossible; no PTR → email functional but deliverability will fail. These are marked unavailable with the actual reason.
- **NP-INF-017:** The panel supports **hardware re-detection** (via heartbeat or "Rescan hardware" button). After a VPS upgrade, installer logic re-runs and shows new capabilities that just became available. This is a first-class UI flow, not a CLI command.
- **NP-INF-018:** The capability set is **re-evaluated on every boot**. If physical capacity shrank (e.g., disk full, memory reduced by provider), the panel surfaces the honest delta and refuses to start workloads that no longer fit.


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

## 7. Websites, virtual hosts, domains, and staging

- **NP-WEB-001:** Create website with owner/account, domain or temporary preview hostname, runtime type, document root and server placement.
- **NP-WEB-002:** Support primary domains, subdomains, aliases, redirects and ownership verification where needed; prevent unauthorized domain takeover.
- **NP-WEB-003:** Generate isolated Nginx virtual hosts from managed templates, validate configuration and reload safely with rollback.
- **NP-WEB-004:** Support per-site PHP selection, static serving and reverse proxy to managed applications; Apache compatibility via declared adapter strategy.
- **NP-WEB-005:** Provide staging environments, preview URLs, environment separation and promotion with data-risk warnings.
- **NP-WEB-006:** Manage redirects, canonical host, custom error pages, caching headers and security headers through safe UI controls.
- **NP-WEB-007:** Show DNS readiness, TLS status, service health, logs, traffic and deployment history for each site.
- **NP-WEB-008:** Domain removal or website deletion must preview effect on aliases, certificates, email, DNS and data.
- **NP-WEB-009:** Detect manual configuration drift and offer reconcile, adopt, or explicitly unsupported options; do not overwrite blindly.
- **NP-WEB-010:** Validate that paths and runtime identities remain inside the authorized account isolation boundary.

**Acceptance scenarios:** New static site serves HTTPS; invalid config never takes down existing sites; conflicting domain assignment is rejected; rollback restores prior working config.


- **NP-WEB-011:** On Standard and Provider profiles, websites are **containerized**. Each site runs in its own Docker container with an Nginx/PHP-FPM/Node/Python image. Caddy (or Nginx) is the reverse proxy in front of the containers. The container is the isolation primitive.
- **NP-WEB-012:** On Nano and Starter profiles where Docker does not fit, static sites are served directly by Nginx from a per-site document root. OS-user isolation applies. No PHP, Node, or database.
- **NP-WEB-013:** The user-facing behavior is identical across profiles: "create website, choose hostname, deploy, live." The implementation differs (container vs. static root) but the UX does not.
- **NP-WEB-014:** The **domain connection wizard** is a first-class flow, not a settings page. It shows the exact A/AAAA/CNAME record to add at the registrar, polls for propagation, confirms when live, and only then offers SSL issuance. It does not claim "live" before the DNS has actually propagated. This is the killer feature for the lay person (Master Blueprint §4).

## 8. Application runtimes, Git, deployments, processes and workers

- **NP-APP-001:** Support managed Node.js, Python, PHP and static apps with adapter-defined version matrices and runtime health checks.
- **NP-APP-002:** Node.js wizard accepts Git repository, uploaded archive, or existing authorized directory; detects `package.json`, scripts, lockfiles and package manager.
- **NP-APP-003:** User configures environment (development/production), runtime version, install/build/start commands, app root, startup entrypoint, working directory and proxy hostname.
- **NP-APP-004:** Provide bulk environment-variable paste and `.env` upload parsing with validation, duplicate detection, masked secrets and no secret exposure in job logs.
- **NP-APP-005:** Support secret references, environment scoping, rotation and explicit redeploy/restart semantics.
- **NP-APP-006:** Allocate internal ports safely and configure HTTPS reverse proxy; applications must not need direct public port exposure by default.
- **NP-APP-007:** Manage process start/stop/restart, health checks, resource limits, restart policy, stdout/stderr logs and uptime.
- **NP-APP-008:** Deployments track source revision, build artifact, environment, actor, start/end time, status, health verification and rollback target.
- **NP-APP-009:** Failed builds do not replace a healthy production deployment; use staged releases where supported.
- **NP-APP-010:** Support Git webhook deploy triggers with signature verification, branch policy and deduplication.
- **NP-APP-011:** Support background workers, scheduled jobs, lifecycle controls and execution history.
- **NP-APP-012:** Prevent cross-tenant access to build artifacts, caches, processes, environment variables and filesystem paths.
- **NP-APP-013:** Show compatibility warnings for unsupported runtime/OS versions; never silently substitute a version.
- **NP-APP-014:** Optional container execution requires explicit capability and stronger isolation controls; do not assume containers are installed.

**Acceptance scenarios:** Import Node repo → detect scripts → bulk paste env → deploy → TLS → logs → rollback; invalid `.env` input yields field-level errors without exposing values; one tenant cannot read another's process secrets.


- **NP-APP-015:** On Standard and Provider profiles, applications are **containerized by default**. The build runs in a constrained build container with CPU/memory/disk limits and no host access. The runtime runs in a separate container. This is the isolation primitive for applications.
- **NP-APP-016:** The application workload is a **serving** category. It uses `cpu.weight = 100` (fair-share under contention) and **no `cpu.max`** (so idle CPU is available). RAM has a hard `memory.max`. This is the correction that prevents the "Node app throttled at 1 vCPU while 3 cores idle" failure mode.
- **NP-APP-017:** **Builds are batch workloads.** They use `cpu.weight = 10` and `cpu.max = 50% of cores`, so a runaway build cannot starve serving workloads. The build queue is durable; a failed build does not affect running apps.
- **NP-APP-018:** **`.env` parsing is a bulk-paste feature**, not just key/value rows. The user can paste a full `.env` file, the panel parses it client-side (or in a safe backend flow), validates each key against a schema, deduplicates, masks values, and only then stores encrypted. The values never appear in job logs, telemetry, or AI context.
- **NP-APP-019:** **The deployment flow is Vercel-like**: pick source (Git / upload / existing) → detect framework → confirm build and start commands → configure environment → deploy → view logs → see health → rollback if needed. The same flow works for Node, Python, Go, Ruby, and static sites.

## 9. File manager, code editor, SSH/SFTP, and terminal

- **NP-FILE-001:** Tree-based file explorer supports breadcrumbs, folders, search, sorting, filters, multi-select, drag/drop upload and progress.
- **NP-FILE-002:** Support create, rename, move, copy, delete/trash, restore, archive/extract, download, permissions and ownership within authorized roots.
- **NP-FILE-003:** Integrated editor supports syntax highlighting, line numbers, search/replace, diff, encoding detection, safe save and conflict detection.
- **NP-FILE-004:** Every path operation canonicalizes paths, blocks traversal, validates symlinks and prevents crossing tenant boundaries.
- **NP-FILE-005:** Uploads enforce size, quotas, content handling policy, archive-bomb protection and malware checks where configured.
- **NP-FILE-006:** Large files and binary files display appropriate previews/limits instead of attempting unsafe editing.
- **NP-FILE-007:** Dangerous permissions and ownership changes require warnings or elevated approval; no indiscriminate `chmod 777` shortcuts.
- **NP-FILE-008:** SSH/SFTP identities and access are managed per account with key rotation, revocation and activity logging.
- **NP-FILE-009:** Browser terminal, if enabled, requires a separate explicit permission and audited constrained session; no automatic root access.
- **NP-FILE-010:** File modifications caused by deployments or backups show provenance where feasible.

**Acceptance scenarios:** Edit a config file with concurrent-change protection; extract archive cannot write outside account root; user without file permission cannot download a guessed path.


- **NP-FILE-011:** On Standard and Provider profiles, the file manager operates on a **virtual filesystem scoped to the site's container volume**. The user sees a tree; the agent maps it to the container's bind-mounted volume. Path resolution uses descriptor-relative operations (`openat2` with `RESOLVE_BENEATH`) to prevent symlink escape races.
- **NP-FILE-012:** On Nano and Starter profiles where sites are not containerized, the file manager operates directly on the site's document root, scoped by the site's OS user.
- **NP-FILE-013:** **Terminal access is disabled by default** for customer roles. On Standard and Provider profiles, a scoped container shell is available when explicitly enabled by the account owner. On Nano, terminal is not available at all (there is no container to exec into). The distinction is displayed honestly in the UI.
- **NP-FILE-014:** **SSH access is a separate permission** from file manager access. The default is **SFTP-only** (file transfer, no shell). SSH shell requires explicit opt-in and audit. Nano does not offer SSH at all; SFTP via the panel-managed SFTP service is the only option.

## 10. Databases and database administration

- **NP-DB-001:** Provision MariaDB/MySQL and PostgreSQL databases with account ownership, placement, quotas and managed credentials.
- **NP-DB-002:** Create/revoke database users and grants with least privilege and tenant-safe naming/authorization.
- **NP-DB-003:** Offer password rotation, connection details, TLS guidance, host access policies and secret-safe copy flows.
- **NP-DB-004:** Database service ports remain private by default; remote access requires explicit allowlist, network review and audit.
- **NP-DB-005:** Support import/export with size limits, job progress, validation and safe error reporting.
- **NP-DB-006:** Support scheduled/on-demand logical backups and restore into original or new target with conflict preview.
- **NP-DB-007:** Database administration UI is authenticated, isolated and scoped to a selected account/database.
- **NP-DB-008:** Expose database health, storage, active connections, supported version and maintenance status.
- **NP-DB-009:** Version upgrades and migrations require compatibility checks, backups and rollback/recovery plans.
- **NP-DB-010:** Delete operations distinguish removing database users, dropping databases and deleting backup copies.

**Acceptance scenarios:** Customer creates DB and user, connects from authorized app, rotates credentials and restores backup; cannot enumerate another customer's databases.


- **NP-DB-011:** Managed databases are **deferred to v2+** (Master Blueprint §12.2). On Nano, databases are **unavailable** — the panel itself uses SQLite, and there is no room for MariaDB/Postgres alongside the panel. On Starter, a small database is **borderline** — allowed only if the installer confirms floor + typical fit. On Standard+, managed databases ship in v2.
- **NP-DB-012:** The **database service** (MariaDB or Postgres) is a `ServiceTemplate` instance, not a bespoke feature. Provisioning a database is provisioning a service from a template with a database-specific health check and backup policy. This is what makes it consistent with storage, apps, cache, and queue.
- **NP-DB-013:** The **database admin UI** (a phpMyAdmin/pgAdmin-like interface) is also a `ServiceTemplate` instance — a container with authenticated access, scoped to a selected database. It is not embedded in the control plane.
- **NP-DB-014:** The **panel's own metadata DB** is SQLite (Nano) or PostgreSQL (Standard+). It is **not** the customer database. The customer's MariaDB/Postgres is a separate container, even on Standard+.

## 11. DNS, domain verification, delegation and DNSSEC

- **NP-DNS-001:** Support external manual DNS, external API adapter mode (when provider integration exists), and local authoritative DNS mode.
- **NP-DNS-002:** DNS onboarding performs preflight, ownership/delegation checks, nameserver/glue guidance, record preview, approval, apply and post-change verification.
- **NP-DNS-003:** Zone editor supports A, AAAA, CNAME, MX, TXT, CAA, SRV, NS, SOA and other adapter-supported types with semantic validation.
- **NP-DNS-004:** Prevent invalid record combinations, apex CNAME conflicts where applicable, unsafe SOA/NS changes and malformed TTLs.
- **NP-DNS-005:** Support zone templates, import/export, change history, diff, dry-run and rollback where technically feasible.
- **NP-DNS-006:** Ownership of a domain is not inferred from entering a hostname; verification method and expiry are tracked.
- **NP-DNS-007:** Show current delegation, authoritative responses, propagation caveats and difference between desired vs observed records.
- **NP-DNS-008:** Local authoritative DNS must support secure primary/secondary replication and independently available secondary topology for production resilience.
- **NP-DNS-009:** Support DNSSEC key lifecycle, DS guidance, rollover and safe disable workflows when adapter supports DNSSEC.
- **NP-DNS-010:** DNS records created by email, websites and ACME carry ownership/provenance; deleting one service must not remove shared records without impact analysis.
- **NP-DNS-011:** DNS operations are not claimed complete merely because a zone file was written; verify authoritative service and intended response.

**Acceptance scenarios:** Customer uses external registrar DNS without transferring registrar; DNS preflight identifies wrong delegation; local zone update propagates to secondary; unsafe record conflict is blocked.


- **NP-DNS-012:** **Authoritative DNS hosting is deferred to v2+** (Master Blueprint §12.2). The v1 default is **guidance mode**: the panel shows the exact A/AAAA/CNAME/MX/TXT record the user needs to add at their registrar, polls for propagation, and confirms when live. The panel does **not** host authoritative DNS in v1. It does **not** claim the DNS change succeeded before `dig` confirms it.
- **NP-DNS-013:** The **domain connection wizard** (§7 NP-WEB-014) is the primary UX for DNS in v1. It combines: "what DNS mode do you want?" (external manual vs. future external API vs. future local authoritative) → "add this record at your registrar" → "we'll poll until it propagates" → "it's live." This is a first-class flow, not a settings page.
- **NP-DNS-014:** **Let's Encrypt HTTP-01** challenge is the primary SSL path. It requires the A record to point at the VPS first, so the domain connection wizard precedes SSL issuance. **DNS-01** challenge is available for wildcard certificates when the user has authoritative DNS — but since authoritative DNS is deferred to v2, wildcard SSL is also deferred to v2.
- **NP-DNS-015:** **PTR (reverse DNS) coordination** is a provider integration concern. Most VPS providers require a support ticket to set PTR. The panel detects this and guides the user to the provider's control panel or support flow. It never claims it set PTR itself.

## 12. TLS/SSL certificates and HTTPS

- **NP-TLS-001:** Request, renew, install, revoke and inventory certificates with domain/account ownership and issuer/expiry metadata.
- **NP-TLS-002:** Support ACME HTTP-01 and DNS-01 challenges when prerequisites are satisfied; expose validation requirements.
- **NP-TLS-003:** Renewal scheduler retries with backoff, respects CA rate limits and alerts before expiry/failure.
- **NP-TLS-004:** Certificate deployment validates private-key/certificate pairing, chain and target configuration before reload.
- **NP-TLS-005:** Wildcard certificates require compatible DNS challenge and authorized DNS access.
- **NP-TLS-006:** Prevent private-key exposure through API, UI, logs, backup exports and support sessions.
- **NP-TLS-007:** Show HTTPS redirect, protocol/cipher policy and certificate coverage for aliases/subdomains.

**Acceptance scenarios:** HTTPS enabled without downtime; failed renewal alerts with root cause; invalid certificate never replaces working certificate.


- **NP-TLS-008:** In v1, **HTTP-01 is the primary challenge** because authoritative DNS hosting is deferred. The flow is: domain points at the VPS → HTTP-01 challenge → certificate issued → HTTPS live. The user sees "obtaining certificate" as a first-class state, not an error message.
- **NP-TLS-009:** **DNS-01 and wildcard certificates are deferred to v2+**, when authoritative DNS hosting ships. The panel UI shows this honestly: "Wildcard SSL requires authoritative DNS hosting, which is not available on your plan. Standard SSL works now."

## 13. Email hosting and deliverability

- **NP-MAIL-001:** Provision mail domains, mailboxes, aliases, forwarders, distribution lists and configurable catch-all behavior.
- **NP-MAIL-002:** Support authenticated SMTP submission, IMAP, webmail integration and mailbox quota management.
- **NP-MAIL-003:** Generate and diagnose SPF, DKIM, DMARC, MX and reverse-DNS prerequisites; distinguish what the panel can set from what requires ISP/provider action.
- **NP-MAIL-004:** Provide anti-spam/anti-malware integration, outbound rate limits, queue inspection, bounce visibility and abuse controls.
- **NP-MAIL-005:** Support mailbox credential resets, MFA where supported, access revocation, retention and backup/restore.
- **NP-MAIL-006:** Expose service prerequisites including public IP reputation, PTR, port-25 restrictions, DNS topology and resource budgets.
- **NP-MAIL-007:** Suspend compromised mailboxes or domains with scope-specific impact previews and clear recovery.
- **NP-MAIL-008:** Keep mailbox data isolated from other accounts and prevent unauthorized cross-domain alias/forwarder creation.
- **NP-MAIL-009:** Track delivery errors and service health without promising inbox placement or deliverability guarantees.
- **NP-MAIL-010:** Support separate mail-server role; never require full mail stack on a small standalone VPS.

**Acceptance scenarios:** Create mailbox, configure DNS, send/receive through supported clients; mail queue and errors are visible; provider blocks abuse without deleting unrelated mail data.


- **NP-MAIL-011:** **Self-hosted email stack is deferred to v2+** (Master Blueprint §12.2). The v1 default is **external SMTP integration**: panel notifications and WordPress outbound mail go through Resend / Postmark / Mailgun / SES. The panel UI recommends Zoho or Google Workspace for actual mailboxes. The `mail adapter` interface exists from v1; the self-hosted mail adapter ships later.
- **NP-MAIL-012:** The reason for this deferral: a fresh VPS with no PTR and no sender reputation will fail deliverability regardless of how good the mail stack is. 30% of the engineering effort and 90% of the operational pain of a hosting panel is email, and most of that pain is deliverability, not software. The v1 default avoids it honestly.
- **NP-MAIL-013:** On Nano, email is **unavailable** — the panel does not offer it. On Starter, email is **borderline** — allowed only if port 25 is open, PTR is set, and the installer confirms floor + typical fit. On Standard+, email ships in v2 with the honest prerequisite checks from NP-MAIL-006.

## 14. Backups, restore, snapshots, retention and disaster recovery

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


- **NP-BKP-011:** Backups are a **batch** workload. They run with `cpu.weight = 10` and `cpu.max = 50% of cores`, so a scheduled backup does not starve websites. They are queued, not rejected. On Nano, backups use restic to an external S3-compatible target (Backblaze B2, R2, etc.) — the panel does not host backup storage on the Nano VPS itself.
- **NP-BKP-012:** **Backup integrity is verified by restore test**, not by process exit code. The panel schedules periodic restore drills to a disposable target and records the result. "Backup succeeded" ≠ "backup is restorable." A backup without a recent verified restore is marked `unverified`, not `healthy`.
- **NP-BKP-013:** On Nano, backups are **remote-only**. There is no local backup destination — the disk is too small. The panel requires a remote backup target during onboarding. This is honest: a Nano VPS cannot backup to itself.
- **NP-BKP-014:** On Standard+, backups may be local or remote. Remote is strongly recommended; the panel shows the storage tier ("local" vs. "offsite") and warns if only local backups exist. "A backup on the same VPS alone is not disaster recovery."

## 15. Security, firewall, isolation, certificates, and audit

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


- **NP-SEC-013:** **Docker is the primary isolation primitive** where Docker fits (Standard / Provider profiles). Container isolation provides filesystem, process, network, and resource separation. OS-user-per-account is the fallback for workloads that cannot containerize (Nano static sites).
- **NP-SEC-014:** **cgroup v2 is required.** Ubuntu 22.04+ / 24.04+ is the supported OS baseline. The installer refuses to install on systems without cgroup v2. This is DEC-015 in Master Blueprint §15.
- **NP-SEC-015:** **`cpuset` (core pinning) is never used.** It reserves cores, which wastes idle CPU. The scenario where a Node app wants 1.5 vCPU on an idle 4-vCPU server but gets throttled to 1.0 because transforms have "1 vCPU reserved" is a failure mode the architecture prevents. `cpu.weight` is the mechanism.
- **NP-SEC-016:** **`cpu.max` is set only on batch workloads** as a runaway protection. Serving workloads (websites, apps, databases) have no `cpu.max` — they can burst to all cores under idle, and share fairly by weight under contention.
- **NP-SEC-017:** The **audit chain is hash-chained from v1**. It cannot be retrofitted. Audit records survive tenant resource deletion under retention policy. The chain is exported to a protected storage boundary not writable by the tenant.
- **NP-SEC-018:** **Container escape prevention.** The agent validates every container configuration before applying: rejects `--privileged`, host PID/network namespace, host root bind mounts, dangerous Linux capabilities, and unrestricted device access. This is Stage 7 AGT-108/AGT-109 — restated here as a product invariant.

## 16. Monitoring, logs, incidents, analytics and notifications

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


- **NP-OBS-010:** Metrics come from **cgroup v2** (`/sys/fs/cgroup/<workload>/cpu.stat`, `memory.current`, `memory.peak`, `memory.events`, `io.stat`, `pids.current`), polled every 15–30 s. This is the same source the installer uses to compute capability fit. The dashboard shows current / typical (7-day p50) / peak (24h p99) / OOM events / throttled seconds / contention signal.
- **NP-OBS-011:** On Nano, monitoring is **minimal**: cgroup stats only, no Prometheus, no Grafana, no Loki. On Standard+, monitoring uses **VictoriaMetrics** (much lighter than Prometheus) plus optional log aggregation. The full Prometheus+Grafana+Loki stack is optional, not default, because the RAM cost is real on a 4 GB VPS.
- **NP-OBS-012:** The **upgrade recommendation logic** is a first-class feature, only possible because allocation is separate from usage (§5.5): peak RAM > 80% for 7+ days → "Consider upgrading"; OOM kills in audit log → "A workload was killed for memory"; CPU contention events → "Workloads are competing for CPU." The recommendations link to the upgrade flow (§9, Master Blueprint §10.5).

## 17. Automation, scheduled jobs, webhooks, APIs, CLI and SDK

- **NP-AUTO-001:** Expose versioned REST API with stable resource IDs, typed errors, pagination, filtering and documented scopes.
- **NP-AUTO-002:** Mutating requests support idempotency keys and return operation/job IDs for asynchronous work.
- **NP-AUTO-003:** Support signed webhooks, event filtering, retries, delivery history, replay and receiver deduplication guidance.
- **NP-AUTO-004:** Support scheduled tasks/cron with account identity, environment, execution limits, logs and disable controls.
- **NP-AUTO-005:** Offer automation rules with trigger, condition, action, dry-run, approval policy and loop prevention.
- **NP-AUTO-006:** CLI and SDKs use public authorized APIs; no privileged undocumented backdoor.
- **NP-AUTO-007:** Rate limits, quotas, abuse prevention and audit apply equally to human/API/AI callers.
- **NP-AUTO-008:** Provide API documentation, examples, change log, deprecation and backward compatibility policy.

**Acceptance scenarios:** Repeated create request with same idempotency key creates one resource; invalid webhook signature rejected; automation cannot exceed creator's delegated permissions.


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

## 19. AI assistant, diagnostics, planning and guarded actions

- **NP-AI-001:** Assistant explains permitted server/account state and provides troubleshooting using authorized, redacted context.
- **NP-AI-002:** Assistant can generate a structured plan for supported tasks (deploy, DNS, TLS, backup, remediation) without directly executing privileged commands.
- **NP-AI-003:** Every AI tool invocation passes through the same identity, scope, policy, quota, approval and audit checks as normal API actions.
- **NP-AI-004:** Read-only diagnosis is distinct from mutation; risky actions require explicit human confirmation and possibly elevated approval.
- **NP-AI-005:** Display planned changes, affected resources, confidence/uncertainty, rollback strategy and post-action verification.
- **NP-AI-006:** Treat logs, repositories, websites and user-provided documents as untrusted; prevent prompt-injection instructions from becoming tool permissions.
- **NP-AI-007:** Do not expose secrets, cross-tenant data, raw privileged shell or unapproved external transmission of sensitive logs.
- **NP-AI-008:** Support cancellation, job tracking, tool error reporting and evidence-backed summaries; do not claim success before verification.
- **NP-AI-009:** Provide AI enable/disable, provider/model configuration, data-retention and consent policies per organization.
- **NP-AI-010:** Audit model/tool identity, requester, context scope, approval and resulting jobs without storing sensitive prompts indiscriminately.

**Acceptance scenarios:** Assistant proposes Nginx repair and awaits approval; malicious text in an application log cannot cause unauthorized server action; AI never reads another tenant's environment variables.


- **NP-AI-011:** The **policy engine is deterministic and separate from the AI** (Stage 4 ARCH-044–047). The LLM proposes; the policy engine disposes. **The AI is never the authority that decides whether its own actions are safe.** This is the single most important architectural constraint in the AI module.
- **NP-AI-012:** **v1 ships read-only diagnostics + proposed-action approval.** The AI can read logs, metrics, state; it can propose a typed operation with a plan; the user approves; the policy engine validates again at execution. **Autonomous action execution is v4+** (Master Blueprint §12). This is the boundary that makes "fully commercial with AI" honest.
- **NP-AI-013:** **The AI has no `exec` access, ever.** It uses the typed agent operation catalog through the standard API. There is no special AI path, no privileged AI role, no AI backdoor. The policy engine evaluates AI actions with the same rules as human actions — the only difference is the `requested_by` principal kind (`ai` vs. `human`), which is recorded in audit.
- **NP-AI-014:** **Prompt content is untrusted.** Logs, website content, repository files, and user-provided documents can contain prompt injection. The AI treats them as data, never as instructions with execution authority. This is enforced by the tool gateway, not by the model's good behavior.
- **NP-AI-015:** **Nano includes AI diagnostics.** The AI is a client of the external LLM API — it does not consume meaningful panel RAM. This is one of the few "advanced" capabilities available on Nano, because it costs ~30 MB and runs on demand.

## 20. Extensions, marketplace, integrations and white-label branding

- **NP-EXT-001:** Define signed/versioned extension manifests with declared capabilities, permissions, lifecycle and compatibility.
- **NP-EXT-002:** Plugins run under bounded identities and cannot access privileged agent operations outside declared grants.
- **NP-EXT-003:** Provide UI extension points, service adapters, event subscriptions, settings and health indicators.
- **NP-EXT-004:** Support install, upgrade, disable, rollback and uninstall with dependency and data-retention handling.
- **NP-EXT-005:** Provider approval and trust policy govern third-party extensions and marketplace distribution.
- **NP-EXT-006:** Reseller/provider branding includes logo, theme, portal domain, support contact and permitted customer-facing copy.
- **NP-EXT-007:** White-label customization must not hide legally required notices or security/incident warnings.

**Acceptance scenarios:** Plugin denied undeclared server permission; incompatible plugin blocked during upgrade; reseller branding does not alter provider security controls.


- **NP-EXT-008:** The **plugin runtime and marketplace are deferred to v2+** (Master Blueprint §12.2). The v1 extension mechanism is the **`ServiceTemplate` primitive** (§21 below and Master Blueprint §5.6) — a signed container image plus a template record is already an extension without the plugin runtime complexity. The `extension manifest` schema exists from v1; the plugin runtime does not.
- **NP-EXT-009:** **White-label branding ships in v2+** with the reseller surface. The data model (organizations with branding fields) exists from v1. The v1 UI has no branding overrides because the reseller surface is not shipped until v2.

## 21. Settings, support, migration, import and export

- **NP-OPS-001:** Separate installation-wide, provider, organization, reseller, account and personal settings with clear inheritance/override rules.
- **NP-OPS-002:** Provide health diagnostics, support bundles with redaction, installation ID/version and configuration provenance.
- **NP-OPS-003:** Support migration/import from supported panels and conventional hosting layouts through adapters, dry-run and validation; do not promise universal migration fidelity.
- **NP-OPS-004:** Support export of account assets, DNS zones, database dumps, mailbox data and metadata as allowed by permissions and formats.
- **NP-OPS-005:** Track migration conflicts, unsupported features, estimated downtime, cutover and rollback strategy.
- **NP-OPS-006:** Provide documented recovery mode when UI, API, metadata DB or agent is unavailable.
- **NP-OPS-007:** Support uninstall/repair/upgrade with backups, preflight, maintenance windows and audit.

**Acceptance scenarios:** Operator generates sanitized support bundle; migration dry-run reports unsupported mail settings; failed upgrade can recover control-plane access.


- **NP-OPS-008:** **`ServiceTemplate` export/import is a first-class primitive.** A storage bucket, a self-hosted app, a database, or a cache instance is exportable as an OCI-compatible artifact plus a config manifest. This is what makes "move my VPS to a bigger server" and "migrate between hosts" possible without a bespoke migration tool. The `ServiceTemplate` primitive from §21 (see below) does the work; this requirement is the product-level commitment that migration is designed in, not an afterthought.
- **NP-OPS-009:** **Migration from cPanel / DirectAdmin / Plesk is deferred to v2+.** The v1 focus is internal migration (VPS-to-VPS) via the `ServiceTemplate` export path. Import from other panels ships when there are customers who need it.

## 21A. Platform Services and ServiceTemplate


**NP-SVC-001:** Every service is expressible as a `ServiceTemplate` record with the shape from Master Blueprint §5.6: identity, source (container image + signature), resources (floor/typical/ceiling + category), network (hostname + SSL + exposure), credentials (generated-once, scoped, revocable), data (volumes + backup + export), health (probe + expected response), configuration (env schema + UI form), metering (metrics to collect), export (target + format).

**NP-SVC-002:** Platform services ship in v1: **Storage** (MinIO or Garage, S3-compatible), **Image transformation** (`imgproxy` or `thumbor`, batch workload), **Delivery** (Caddy cache + optional Cloudflare/Bunny integration). These are the reference implementations that prove the `ServiceTemplate` primitive.

**NP-SVC-003:** The **storage + transform + delivery** trio is three separate services, not one "CDN." Each has its own resource profile, its own failure mode, its own billing unit. The panel uses honest naming: "Media storage" or "Asset storage + transformation," not "CDN." The optional Cloudflare/Bunny integration is the CDN.

**NP-SVC-004:** The **value proposition is integration depth, not storage.** Storage itself is undifferentiated. The differentiator is that the user's websites, apps, and self-hosted services are provisioned by the same panel, in the same dashboard, with the same access model, and can consume the bucket with one click: "Use my storage bucket for media" injects the endpoint and scoped credentials into the site's config automatically.

**NP-SVC-005:** **Self-hosted app templates** ship in v1: WordPress, Nextcloud, Gitea, Ghost, n8n, Uptime Kuma. Expansion in v2+: Meilisearch, Plausible, MinIO, imgproxy, PostgreSQL, MariaDB, Redis — 20–30 apps total. Templates are **data, not code** — adding a template is an admin action, not a release.

**NP-SVC-006:** New data entities: `service_templates`, `service_instances`, `service_credentials`, `service_volumes`, `service_health_checks`, `service_metrics`, `service_exports`. New agent operations: `service.provision`, `service.configure`, `service.start`, `service.stop`, `service.health_check`, `service.backup`, `service.restore`, `service.export`, `service.import`, `service.uninstall`. New module in the 29-area scope: **Platform Services**.

**Acceptance scenarios:** Storage bucket created with S3 keys shown once; website consumes bucket with one-click integration; image transform runs as batch workload (yields under load, runs flat-out when idle); Nextcloud template provisions, health-checks, backs up, exports, uninstalls cleanly; template catalog expansion is a database insert, not a release.

## 22. Detailed end-to-end journeys

### J1 — First VPS owner website

1. Install panel with supported OS preflight and backup/recovery notice.
2. Enroll owner, configure MFA, verify SSH/admin recovery and host connectivity.
3. Dashboard detects server capacity and installed service modules.
4. Choose Create Website → static/PHP/app → account/hostname/document root.
5. Validate domain authorization, quota, port/service dependencies and config preview.
6. Provision site through durable job; validate and reload web service.
7. Upload or deploy content; test preview URL.
8. Guide DNS configuration (external or authoritative) and verify observed records.
9. Issue TLS when challenges are reachable; verify HTTPS.
10. Configure backup policy and verify recovery point; view health and logs.
**Failure branches:** missing DNS, failed ACME, unsupported PHP, full disk, Nginx syntax failure, agent offline. Each must preserve existing healthy services and show recovery steps.


### J2 — Developer Node.js deployment

1. Select Git/upload/existing files → detect package.json and scripts.
2. Choose runtime version, install/build/start commands, production mode, app root and hostname.
3. Paste or upload `.env`; classify secrets and validate syntax without echoing values.
4. Preflight package support, account quotas, internal port allocation and Git access.
5. Build isolated release; store logs and artifact; deploy behind proxy.
6. Run health checks; publish success only after verification.
7. View process logs and revisions; rollback to a known-good release.
**Failure branches:** dependency install failure, secret missing, unhealthy process, Git webhook replay, port collision, failed rollback.


### J3 — Provider provisions reseller and customer

1. Provider creates reseller delegation and ceilings.
2. Reseller configures permitted brand and package.
3. Reseller invites customer or creates hosting account within delegated quota.
4. Placement engine selects compatible server; provision OS/service isolation.
5. Customer receives secure onboarding; can manage only owned resources.
6. Provider sees aggregate usage; reseller sees delegated accounts; customer sees own account.
7. Suspension, upgrade, transfer and termination follow lifecycle policy and retention.
**Failure branches:** over-allocation, insufficient capacity, delegated permission revoked, partial provisioning, billing failure.


### J4 — Email domain activation

1. Check public IP, port availability, reverse DNS, spam controls and server role.
2. Add verified mail domain and mailbox policy.
3. Provision mail services and publish or guide MX/SPF/DKIM/DMARC records.
4. Verify DNS and authentication; test submission/IMAP and inspect queue.
5. Configure mailbox backups, quotas, rate limits and abuse alerts.
**Failure branches:** provider blocks SMTP, missing PTR, DKIM mismatch, queue buildup, blacklisted IP.


### J5 — Backup disaster recovery

1. Operator identifies failed site/database or server.
2. Inspect latest valid offsite recovery point and restore test evidence.
3. Select restore target; preview overwrite and downtime.
4. Approve; restore to staging/new target where possible; verify integrity and service health.
5. Cut over deliberately; log incident, recovery time and outcome.
**Failure branches:** corrupt backup, missing encryption key, insufficient disk, target incompatibility, incomplete restore.


### J6 — AI-guided incident remediation

1. User asks why website is down; AI retrieves only permitted status/logs.
2. AI correlates recent deployment, DNS/TLS status, agent/service health.
3. AI presents evidence, uncertainty and proposed typed remediation plan.
4. User approves specific action; policy engine validates again at execution.
5. Job executes; verification determines actual result; audit/notification follows.
**Failure branches:** no permission, stale data, approval timeout, untrusted prompt injection, remediation failed.


### J7 — First-run router **(new in v2)**

1. User completes owner enrollment (J1 steps 1–2).
2. `GET /` returns the router: "What do you want to run?" with six cards.
3. User picks a card (Website / Application / Self-hosted app / Data & storage / Infrastructure / I'm not sure).
4. The panel routes to the corresponding guided flow.
5. Each flow ends at its own "you're live" moment (site live / app responds / app login / bucket exists / service running / suggested first step).
6. The user lands on the **service portfolio** view — everything they now run, with health and next actions.
7. The menu (29-area IA) is reachable via "Advanced tools" but is not the default landing.

**Failure branches:** none of the six cards fits → "I'm not sure" → guided tour → suggested first step; the user's VPS cannot run the chosen card → honest capability display → upgrade path.

### J8 — Nano layman **(new in v2)**

1. User installs NexusPanel on a 512 MB VPS.
2. Installer detects 512 MB / 1 vCPU / 10 GB / no port 25.
3. Installer computes capability set: static sites + storage client + remote backup + AI diagnostics. Refuses: PHP, Node, databases, transforms, email.
4. User sees the honest capability display: "Your VPS can run: static sites, storage client, remote backups, AI troubleshooting. Your VPS cannot run: WordPress, Node apps, databases, transforms, email."
5. User picks "A website" from the router → static site flow → live in 60 seconds with SSL.
6. User sees the service portfolio: one static site, live, with SSL.
7. No greyed-out menus. No fake features. Just the honest picture.

## 23. Initial authorization matrix (illustrative; subject to policy approval)

Legend: `O` own/delegated resources, `P` provider scope, `R` read-only in granted scope, `—` denied unless separate grant, `A` approval required for risky actions. Roles are *capability bundles*, not permanent hardcoded hierarchy.

| Capability | VPS owner | Provider admin | Reseller | Hosting customer | Developer | Auditor |
|---|---|---|---|---|---|---|
| Manage servers/agents | O | P | — | — | — | R |
| Create packages/accounts | O | P | O | — | — | R |
| Manage sites/apps | O | P (policy) | O (policy) | O | O | R |
| Manage databases/files | O | P (policy) | O (policy) | O | O (grant) | R |
| Manage DNS/certificates | O | P (policy) | O (policy) | O | O (grant) | R |
| Manage mailboxes | O | P (policy) | O (policy) | O | — | R |
| Restore/delete backups | O/A | P/A | O/A | O/A | O (grant)/A | R |
| Change firewall/SSH | O/A | P/A | — | — | — | R |
| Billing and invoices | O | P | O (delegated) | O (own) | — | R (grant) |
| View audit/security | O | P | O (delegated) | O (own) | O (grant) | R |
| Use AI for changes | O/A | P/A | O/A | O/A | O/A | — |
| Impersonate/support access | A (grant) | A (grant) | A (grant) | — | — | — |

**Policy invariant:** A provider's operational visibility does not imply unrestricted reading of customer files, email contents, or secrets. All privileged support operations require explicit scopes and audit.


## 24. Nonfunctional requirements and service-level targets

These are **proposed engineering targets**, not measured performance claims or contractual SLAs. Validate on a defined reference topology and revise before release.

| ID | Category | Requirement / target |
|---|---|---|
| NP-NFR-001 | Security | No known critical vulnerabilities at production release; privileged paths independently reviewed |
| NP-NFR-002 | Availability | Single-node mode documents its single point of failure; HA claims only for tested multi-node topology |
| NP-NFR-003 | Reliability | Jobs durable across API/worker restart; idempotent retries; no silent partial success |
| NP-NFR-004 | Performance | Common UI actions target responsive feedback under 2 seconds excluding async provisioning, subject to benchmark |
| NP-NFR-005 | Scalability | Define and benchmark tenant/server/job limits by topology; no unsupported "unlimited" claims |
| NP-NFR-006 | Resource budget | Every module declares minimum RAM/CPU/disk and measured idle/peak overhead |
| NP-NFR-007 | Recovery | Document RPO/RTO per service tier; demonstrate periodic restore tests |
| NP-NFR-008 | Observability | Every async job has correlation ID, progress, structured error, audit and final verification |
| NP-NFR-009 | Compatibility | Supported OS/runtime/service version matrix with explicit tested status |
| NP-NFR-010 | Accessibility | WCAG 2.2 AA target with automated and manual audits |
| NP-NFR-011 | Internationalization | Localizable copy, UTC storage, locale/time-zone formatting and extensible currencies |
| NP-NFR-012 | Data privacy | Defined data classes, retention, export, deletion, processor boundaries and breach response |
| NP-NFR-013 | Upgradeability | Signed updates, preflight, migration backups, compatibility checks and recovery/rollback plan |
| NP-NFR-014 | Extensibility | Versioned plugin/adapter API and permissions; incompatible modules fail safely |
| NP-NFR-015 | Testability | Reproducible disposable test hosts and contract/integration/security/end-to-end suites |
| NP-NFR-016 | Operability | Documented install, repair, emergency access, backup, restore, upgrade and incident runbooks |


| ID | Category | Requirement / target |
|---|---|---|
| NP-NFR-017 | **Capability model** | The installer computes floor / typical / ceiling for every requested workload and refuses overcommit violations. Every workload's cgroup metrics are polled, aggregated, and surfaced. Allocation is separate from usage. |
| NP-NFR-018 | **cgroup v2** | Ubuntu 22.04+ / 24.04+ is the baseline. `cpu.weight` is used for fair-share CPU; `cpu.max` only on batch; `memory.max` on all; `cpuset` never used. Verified by test. |
| NP-NFR-019 | **Nano floor** | Panel + OS baseline ≤ 380 MB on a 512 MB VPS. Static sites, storage client, remote backup, AI diagnostics available. Everything else honestly unavailable. |
| NP-NFR-020 | **Guided-default** | `GET /` resolves per persona. Lay person lands on the router. This is a release gate, not a later feature. |
| NP-NFR-021 | **Payment adapter** | Async-first state machine supports both synchronous (card, <1s) and asynchronous (mobile money, 10s–2min) gateways without code forks. M-Pesa is the reference implementation. |
| NP-NFR-022 | **ServiceTemplate** | Every service is expressible as a `ServiceTemplate` record. Adding a template is a database insert, not a release. Verified by test. |
| NP-NFR-023 | **Audit chain** | Hash-chained from v1. Cannot be added later. Records survive tenant resource deletion under retention policy. |
| NP-NFR-024 | **Migration path** | `ServiceTemplate` export/import supports VPS-to-VPS migration. Migrating a Nano to Standard is a documented, tested flow. |

## 25. Global error and edge-case catalog

| Code | Situation | Required behavior |
|---|---|---|
| AUTH_DENIED | Actor lacks scope | Deny; avoid tenant enumeration; audit sensitive denial |
| MFA_REQUIRED | Sensitive action without step-up | Prompt secure step-up; do not queue operation |
| QUOTA_EXCEEDED | Account/reseller/provider ceiling reached | Explain limiting dimension; show upgrade/request path |
| CAPACITY_UNAVAILABLE | No eligible node | Do not overprovision; show constraints and alternatives |
| SERVER_OFFLINE | Agent unreachable | Preserve desired state; queue only if safe; report unknown outcome |
| CONFLICT | Resource/domain/path/version conflict | Show conflicting owned resource where authorized; safe retry |
| VALIDATION_FAILED | Invalid DNS, env, config, import | Field-level errors; do not modify active service |
| APPROVAL_REQUIRED | High-risk operation | Show diff/blast radius; expire pending approval safely |
| PARTIAL_FAILURE | Multi-step job partially completed | Identify completed steps; retry/compensation plan; audit |
| DRIFT_DETECTED | Manual change differs from desired state | Show diff and safe reconcile/adopt choice |
| SECRET_UNAVAILABLE | Vault/key failure | Stop dependent job; never fall back to plaintext |
| BACKUP_INVALID | Corrupt/incomplete recovery point | Block destructive restore; show other recovery points |
| RATE_LIMITED | Abuse or service throttling | Clear retry guidance; no duplicate operations |
| VERSION_UNSUPPORTED | Unsupported OS/runtime/module | Block incompatible install; explain tested alternatives |
| BILLING_HOLD | Payment/entitlement state restricts action | Apply policy without deleting data unexpectedly |
| RECOVERY_REQUIRED | Control plane/agent unavailable | Provide documented out-of-band recovery workflow |


| Code | Situation | Required behavior |
|---|---|---|
| CAPABILITY_UNAVAILABLE | Installer determined a capability is unavailable on this VPS | Explain the resource delta required; offer upgrade path; do not grey out silently |
| MEMORY_PRESSURE | `memory.events.high` or `oom_kill` observed on a workload | Surface in UI; suggest upgrade; do not silently retry |
| CPU_THROTTLED | Batch workload throttled by `cpu.max` | Expected behavior; surface as informational, not error |
| PENDING_PAYMENT | M-Pesa or async gateway awaiting callback | Show pending state; allow retry; timeout per gateway policy |
| PAYMENT_CALLBACK_UNMATCHED | Payment callback does not match an invoice | Hold for manual reconciliation; do not credit or reject automatically |
| TEMPLATE_UNHEALTHY | A `ServiceTemplate` instance health check fails | Show unhealthy state; offer restart, reinstall, or restore from backup |
| ROUTER_NO_MATCH | User picks "I'm not sure" and the guided tour has no suggestion | Fall back to "here's what your VPS can run" (capability display) |

## 26. Module definition of done and release evidence

For each module, require: (1) approved requirement IDs and journeys; (2) screen catalog and all UI states; (3) permission and quota matrix; (4) API schemas and errors; (5) data model/migrations; (6) typed agent operations and config ownership if applicable; (7) durable jobs and rollback/compensation; (8) audit and observability; (9) backup/recovery impact; (10) unit/integration/E2E/security tests; (11) performance/resource benchmark; (12) upgrade and compatibility tests; (13) documentation and support runbook; (14) product-owner acceptance. No module is "done" because a button or demo works.


- **(15) Capability declaration:** the module declares its `floor` / `typical` / `ceiling` per resource, its workload category, and its prerequisites. The installer reads it from a data file, not code. Tested by the capability-model test suite.
- **(16) ServiceTemplate binding:** if the module is a service (storage, transform, delivery, self-hosted app, cache, queue, search), it is expressible as a `ServiceTemplate` record. If it can only be a bespoke feature, that is a design smell and must be reviewed.
- **(17) Usage surfacing:** the module's cgroup metrics are surfaced in the dashboard as current / typical / peak. Allocation without measurement is not complete.
- **(18) Upgrade path:** the module declares whether its resources can be upgraded in place or require migration. The upgrade flow is tested, not assumed.
- **(19) Async adapter (for commercial modules):** payment adapters implement the async-first interface and are tested against both a synchronous stub and an asynchronous stub.
- **(20) Audit chain integrity:** every mutation writes to the hash chain. Verified by test.

### Release readiness gates

- **Foundation gate:** IAM, policy, secure agent enrollment, job durability, audit and recovery are proven.
- **Hosting gate:** site/app/domain/TLS/files/backups work end-to-end on supported OS with isolation tests.
- **Services gate:** DB/mail/DNS adapters pass compatibility, failure and recovery tests; mail/DNS topology prerequisites documented.
- **Commercial gate:** account/reseller/package/billing lifecycle, usage accuracy and suspension/retention pass E2E tests.
- **Fleet gate:** remote agent trust, placement, failure domains, maintenance, upgrades and recovery tested.
- **Public launch gate:** threat model, independent security review, backups/restore, upgrade recovery, support docs, licensing and privacy controls approved.


- **Capability gate:** installer correctly detects Nano / Starter / Standard / Provider, refuses below floor, computes floor/typical/ceiling for every requested workload, and produces honest capability displays. Tested on 512 MB, 1 GB, 2 GB, 4 GB, 8 GB reference VPSes.
- **Guided-default gate:** `GET /` returns the router for a new owner. The router routes to each of the six paths. Each path ends at a "you're live" moment. Tested end-to-end with a non-technical user (heuristic evaluation or moderated usability test).
- **cgroup gate:** serving workloads are never throttled on idle CPU; batch workloads yield to serving under contention; `memory.low` protects serving workloads from batch memory pressure; OOM kills are surfaced honestly. Tested with synthetic load.
- **Commercial gate (v2 addition):** M-Pesa sandbox STK push → callback → confirmation; timeout → retry; reversal → audit; reconciliation between M-Pesa transaction records and panel invoices. Tested end-to-end against Daraja sandbox.
- **ServiceTemplate gate:** storage, transform, delivery, and at least one self-hosted template (Nextcloud) provision, configure, health-check, back up, export, and uninstall cleanly. Tested end-to-end.
- **Nano gate:** the panel installs on a 512 MB VPS, runs static sites, refuses unavailable capabilities with honest reasons, and the upgrade detection works. Tested on a real 512 MB VPS.

## 27. Dependencies and delivery mapping


| Delivery phase | Modules introduced/validated | Not a permanent scope limit |
|---|---|---|
| 0: Specification | This PRD, architecture, policy matrix, screen inventory, ADRs | All domains documented now |
| 1: Foundation | IAM, organizations, jobs, audit, API, agent, installer preflight | Full resource model designed |
| 2: First complete journey | Websites, apps, files, domains, TLS, backups, monitoring | Customer workflow production-quality |
| 3: Hosting services | PHP/Python, databases, DNS modes, email, restores, isolation | Service role compatibility proven |
| 4: Commercial operations | Provider/reseller, packages, quotas, usage, billing adapters, branding | Full tenant hierarchy operational |
| 5: Fleet and ecosystem | Multi-server placement, plugins, AI actions, HA patterns | Distributed topology validated |
| 6: Hardening/release | Security, performance, recovery, accessibility, documentation | Public launch only after gates |

**Dependency rules:** Secure agent before privileged hosting mutations; ownership/policy before customer resources; durable jobs before provisioning; backups before destructive migration; DNS readiness before ACME issuance; metering before usage-based billing; permissions before AI actions; compatible service adapter before enabling its UI.


**Spines complete:** Go control plane + agent + static SPA (Standard profile); Postgres + SQLite adapter; Caddy + Docker (Standard); Nginx (Nano); full typed operation catalog; full job state machine; full policy engine; full audit chain; full secrets management; full capability model (floor/typical/ceiling, cgroup); full installer logic (5 steps, overcommit rules); full commercial ledger + M-Pesa adapter; full persona matrix; AI read-only diagnostics.

**Journey complete:** First-run router ("what do you want to run?"); static site end-to-end; WordPress end-to-end; Node app end-to-end; ServiceTemplate end-to-end (Nextcloud as reference); storage bucket end-to-end (MinIO + transform + delivery); backup + restore (restic to S3); domain connection wizard; SSL via Let's Encrypt.

**Installable profiles:** Nano (static sites + storage client + remote backup + AI diagnostics); Standard (full capability).


Provider/fleet profile (Standard on multiple nodes); reseller surface; Stripe adapter; email stack (self-hosted, with honest prerequisites); authoritative DNS; plugin runtime; full AI action execution (approval-bound); migration from cPanel / DirectAdmin; self-hosted template catalog expansion (20–30 apps).


Full billing automation (dunning, proration, refunds); multi-jurisdiction tax; white-label; marketplace; cross-region backups; HA topology (Postgres replication, leader election).

## 28. Decisions requiring product-owner approval


| ID | Decision | Recommended starting position | Blocking? |
|---|---|---|---|
| **DEC-001** | Control plane + agent + frontend stack | **Go + static SPA + Go agent.** Requires Go rewrite from Stage 4 ADR-002. | **Yes** — everything downstream assumes this |
| DEC-002 | Product name and trademark | NexusPanel working name; finalize before public launch | No |
| DEC-003 | Licensing model | Open core (AGPL-3.0 community + commercial edition) | No |
| **DEC-004** | Nano VPS minimum floor | 380 MB panel footprint; refuse install below this | **Yes** — installer logic |
| **DEC-005** | Docker on Nano or no | **No Docker on Nano** — static sites only | **Yes** — capability model |
| **DEC-006** | M-Pesa adapter scope | Daraja STK Push + C2B; KES + 16% VAT + KRA invoice | **Yes** — commercial v1 |
| DEC-007 | Storage engine | MinIO (familiar) or Garage (lighter); MinIO recommended | No — swappable |
| DEC-008 | Transform engine | `imgproxy` (Rust, fast, smaller) recommended | No — swappable |
| DEC-009 | Delivery integration | Cloudflare or Bunny first; Cloudflare recommended | No |
| DEC-010 | AI provider in v1 | Single provider behind abstraction; no router | No |
| **DEC-011** | Guided-default enforcement | **Router is the primary IA; menu is secondary.** Non-negotiable. | **Yes** — lay-person gate |
| DEC-012 | Self-hosted template catalog v1 | 5–10 apps (Nextcloud, Gitea, Ghost, n8n, WordPress) | No |
| DEC-013 | Remote agent / mTLS PKI | Deferred to v2+; local Unix socket in v1 | No |
| DEC-014 | HA topology | Deferred to v3; single node in v1, design interfaces for later | No |
| **DEC-015** | cgroup v2 requirement | **Yes** — Ubuntu 22.04+ / 24.04+ | **Yes** — capability model depends on it |

## 29. Next specification deliverables (Stage 3 onward)

1. **Screen-by-screen UX specification:** navigation map, screen catalog, wizard states, validation and error copy, design system and responsive flows.
2. **System architecture and ADRs:** component boundaries, sequence diagrams, trust model, topology, orchestration, service adapters, failure handling.
3. **Full authorization matrix:** exact actor/action/resource rules, delegation, step-up and approval conditions.
4. **API contract:** OpenAPI endpoints, payloads, events, errors, pagination, idempotency and webhook signatures.
5. **Database schema:** SQL entities, relationships, constraints, tenancy boundaries, migrations, retention.
6. **Installer/operations specification:** OS compatibility, preflight, safe bootstrap, upgrade/repair, rollback, security and recovery runbooks.
7. **Testing matrix:** requirement-to-test traceability, automated acceptance suites, security and recovery scenarios.


8. **Capability data file:** the YAML/JSON artifact that declares every module's floor/typical/ceiling/category/prerequisites. This is the source of truth for the installer.
9. **Persona × capability matrix data file:** the source of truth for the router and for menu visibility.
10. **`ServiceTemplate` schema and reference implementation:** the storage, transform, delivery, and Nextcloud templates — end-to-end.
11. **M-Pesa Daraja adapter specification:** the async payment state machine against the Daraja sandbox, including reconciliation.
12. **cgroup configuration spec:** the exact `cpu.weight`, `cpu.max`, `memory.low`, `memory.max` values per category, with test cases for idle/contention/batch behavior.
13. **Guided router flow spec:** the six-card screen, the routing logic, the "I'm not sure" fallback, and the "you're live" moment for each path.

## 30. Approval checklist

- [ ] Confirm product name and ownership.
- [ ] Confirm all user audiences and functional modules are in scope.
- [ ] Review lifecycle and destructive-action rules.
- [ ] Approve baseline authorization/delegation policy.
- [ ] Confirm target OS/runtime compatibility strategy.
- [ ] Resolve or explicitly defer each `DEC-*` item with owner/date.
- [ ] Approve stage sequencing and production release gates.
- [ ] Sign off Stage 2 PRD as the baseline for Stage 3 UX specification.


- [ ] **Approve DEC-001 (Go stack), DEC-011 (guided-default rule), DEC-015 (cgroup v2).** These are the gate.
- [ ] **Approve the four profiles** (Nano / Starter / Standard / Provider) as narrative bundles with the capability model as mechanism.
- [ ] **Approve the capability model** (floor/typical/ceiling, asymmetric overcommit, batch/serving/infrastructure categories).
- [ ] **Approve Docker-per-site as the primary isolation primitive** where Docker fits.
- [ ] **Approve the `ServiceTemplate` primitive** and its reference instances (storage, transform, delivery, self-hosted apps).
- [ ] **Approve the async-first payment adapter** and M-Pesa as the reference implementation.
- [ ] **Approve the "What do you want to run?" router** as the primary first-run IA.
- [ ] **Approve the release gates** including the capability gate, guided-default gate, cgroup gate, commercial gate, ServiceTemplate gate, and Nano gate.

---

**Change control:** This document is a detailed draft derived from the approved direction of the v1.0 PRD. It adds proposed product rules and acceptance scenarios for review; those additions are not claimed as already approved. Keep the VPS unchanged until architecture and installation controls are reviewed.

**End of Stage 2 — PRD v2.0.**

---



---

# Consolidated decision gates and sign-off

**Approval status: OPEN.** No product-owner approval is inferred from an uploaded draft or a request to consolidate documentation.

| Gate | Proposed baseline | What must be verified or approved |
|---|---|---|
| DEC-001 | Go control plane + Go agent + static SPA | Product/engineering sign-off and measured Nano footprint |
| DEC-003 | Open-core AGPL community + commercial edition | Legal/license compatibility and commercial strategy approval |
| DEC-004 | Nano minimum resource floor | Define whether 380 MB refers to panel footprint, total available memory or minimum free headroom; test complete host |
| DEC-005 | No Docker on Nano | Verify supported Nano website/SFTP workflows and safe fallback |
| DEC-006 | M-Pesa STK Push/C2B reference | Gateway eligibility, callback authenticity, VAT/tax/invoicing and refund requirements |
| DEC-011 | Guided router is default beginner IA | UX tests, fallback and permission/capability contract |
| DEC-013 | Unix socket local agent in v1 | Peer identity, file permissions, threat model and remote topology boundary |
| DEC-015 | cgroup v2 as resource enforcement foundation | Kernel/OS compatibility, nested/containerized VPS behavior and resource measurement |
| Data portability | SQLite Nano / PostgreSQL Standard+ | Logical schema parity, concurrency, backups, migrations and upgrade path |
| Template extensibility | Signed catalog plus typed lifecycle extensions | Supply-chain verification, compatibility, rollback and stateful-service recovery |
| Release scope | 29 areas, phased v1–v4 | Approve exact v1 launch acceptance matrix and dependencies |

## Final release approval checklist

- [ ] Product owner approves the 29-area canonical scope and phased delivery slices.
- [ ] Every retained `NP-*` requirement has a stable owner, area ID, phase and acceptance evidence reference.
- [ ] All P0 security, tenant isolation, payment integrity and backup/restore requirements have verified tests.
- [ ] Nano/Starter/Standard/Provider profiles have reproducible resource and installation benchmarks.
- [ ] No unsupported capability is advertised as available or installed.
- [ ] No v1 user journey ends in an unexplained disabled state.
- [ ] UX, architecture, data, API, agent and security contracts pass cross-volume consistency checks.
- [ ] All approved exceptions, deferred features and operational limits are published in release notes.
- [ ] The initial Leaseweb server remains untouched until explicit deployment authorization.

**End of Volume 01.**
