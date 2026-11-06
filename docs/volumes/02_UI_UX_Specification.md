# NexusPanel — Volume 02: Complete UI/UX & User Journeys

**Edition:** Consolidated baseline 1.0 | **Status:** Proposed specification — not implementation-approved
**Date:** 6 November 2026 | **Working product name:** NexusPanel (unconfirmed)
**Parent contract:** `01_Master_PRD.md` (authoritative proposed product requirements); 29-area product scope.
**Scope:** Standalone VPS owner, hosting customer, developer, reseller, provider operator, billing/support/auditor; desktop, tablet and mobile.
**Design target:** Premium, calm, fast, approachable hosting UX with deep expert controls. *Inspired by the clarity of modern deployment platforms, without copying any proprietary interface.*
**Governance:** No production server changes are authorized by this document. Requirements marked **PROPOSED** are subject to approval. Volume 01 `NP-*` identifiers remain authoritative.



> **Authority.** This volume is the proposed UX behavior contract subordinate to `01_Master_PRD.md`. It integrates the uploaded UI/UX v2.0 material into the ten-volume master specification. A design is not a deployed feature. No change to the Leaseweb VPS is authorized.

> **Product scope.** All **29 platform areas** are in scope, even when unavailable on a particular server or deferred to a later release. The Nano profile must never promise dynamic hosting or containers that the actual host cannot support. Hardware and UX performance targets require measured validation.

## Volume 02 — Governing decisions and cross-volume alignment

1. **Primary landing.** For first-time and create-oriented experiences, `GET /` serves the static SPA and the SPA composes a guided landing from authenticated capability, permission, and journey-preference data. The HTTP server does not infer a user's persona and redirect them to an arbitrary UI route. Returning users may select the service portfolio as their preferred landing. A reseller/provider landing uses the same permission model with a different task-focused presentation. This resolves conflicting descriptions of server-side versus client-side router ownership.
2. **Security and persona.** Persona and Guided/Standard/Expert are UX preferences, never authorization principals. Every action uses the same server-side policy and capability enforcement described in Volume 01. Forbidden resources are not exposed merely to explain unavailability.
3. **Capability presentation.** The UI distinguishes `supported`, `not_installed`, `insufficient_resources`, `unsupported_os`, `blocked_by_policy`, `not_authorized`, `temporarily_unhealthy`, and `not_yet_released`. It must not misrepresent a future-release capability as a resource shortage. Where permitted, show prerequisite and resource delta with provenance and confidence; otherwise show an appropriate non-disclosing explanation.
4. **Service state.** Provisioning and deployment progress derives from durable operation records; successful HTTP submission is not successful service activation. The UI differentiates desired, pending, observed, failed, and degraded states, with timestamp and next action.
5. **Resource metrics.** `current`, `typical` (7-day p50), and `peak` (24-hour p99) display only when enough data exists; otherwise use `insufficient history` rather than fabricating statistics. Allocation is a ceiling, not a reservation.
6. **Template boundaries.** A unified ServiceTemplate catalog and provisioning flow can host specialized service detail screens. Database integrity, backup consistency, queue persistence, and mail delivery must not be reduced to generic start/stop controls.
7. **Release scope.** Every screen in the complete catalog is a target specification, not a v1 delivery promise. Navigation and entry points are gated by actual phase, topology, resources, entitlement and policy.
8. **Testing.** Every `UX-*` route and `FLOW-*` journey requires requirement traceability, authorized API calls, accessibility, success/error states, and automated/manual evidence before its module ships.

## UX traceability rules

| Contract | Source of truth | UX obligation | Verification |
|---|---|---|---|
| Functional behavior | Volume 01 `NP-*` | Screen and journey acceptance | UX-to-PRD mapping |
| API behavior | Volume 05 | No optimistic claims of completed mutations | Async-state E2E |
| Data and ownership | Volume 04 | Visible tenant/account/server context | Cross-tenant E2E |
| Privileged operations | Volume 06 | Preview, approval, timeline and recovery | Agent failure E2E |
| Templates and platform services | Volume 07 | Template-aware but lifecycle-specific screens | Template conformance tests |
| Security | Volume 08 | Step-up auth, sensitive data masking, audit | Threat and accessibility tests |
| Billing/resellers | Volume 09 | Pending payment, entitlement, delegation | Payment/reseller E2E |
| Release quality | Volume 10 | Responsive/accessibility/performance gates | Release evidence |

---

## 0. Purpose, traceability and design rules

This is a **behavioral UI specification**, not pixel-perfect Figma artwork or production frontend code. It defines information architecture, screen inventory, role-aware states, interactions, forms, validation, responsive behavior, permissions, accessibility, failure recovery and testable design acceptance. Companion Volumes 03–06 define the corresponding implementation architecture, data, API, and host-operation contracts.

- **UI identifiers:** `UX-<AREA>-###` for screens; `FLOW-###` for multi-screen journeys; `DS-###` for design-system components; `UXR-###` for cross-cutting requirements.
- **Traceability:** Every screen must map to Volume 01 requirement(s), authorized API actions, success/error states, analytics events and E2E test cases. Screens are never evidence that a backend feature works.
- **MUST:** Required for the module to be complete; **SHOULD:** preferred; **MAY:** optional. Proposed defaults are not approved product decisions.
- **Core safety rule:** UI does not bypass server-side authorization, policy, quota, approval or agent controls. Hidden buttons are not a security boundary.
- **Universal action sequence:** inspect → configure → validate → review impact → authorize/approve → submit operation → track progress → verify actual result → audit/recover.
- **Explicit distinction:** A complete UX catalog does not mean all services are enabled on a single small VPS. Unsupported, unavailable, disabled and unauthorized are different states with different explanations.

> **Integrated requirement:** Additional design rules from the deliberations:

- **Guided-default rule:** `GET /` resolves the default landing page per persona (jobs-to-be-done, Volume 01 personas and journeys). The **guided router is the primary IA** for personas whose job is "get something new working." The **service portfolio is the primary IA** for personas whose job is "manage what I already have." The **menu is the secondary IA**, reachable via "Advanced tools." This is an architectural constraint, not a UI preference. If the default landing is the menu, the lay person drowns — that is cPanel's failure mode.
- **Persona ≠ permission.** Guided / Standard / Expert are presentation modes. Permissions are identical across modes. The API enforces permissions; the UI reflects them. If you find yourself writing `if (role === 'reseller') showCustomerMenu()`, stop — that is the permission model leaking into the UI.
- **Capability display is honest.** When the installer determines a capability is unavailable, the UI shows the reason and the resource delta — never a greyed-out button without explanation. "Your VPS cannot run WordPress (needs 150 MB more)" is honest. A disabled menu item with a tooltip is not.
- **Allocation is a ceiling, not a reservation.** Every workload's dashboard shows current / typical (7-day p50) / peak (24h p99). The upgrade recommendation is only possible because allocation is separate from usage.
- **Batch workloads are visible as batch.** Image transforms, backups, builds, AI diagnostics and migrations are shown with their queue state, not as errors. "Queued, will run when server is idle" is the honest state, not "failed" or "slow."
- **The Nano experience is complete, not cut-down.** A 512 MB VPS gets the full panel with an honest capability set. No greyed-out menus. No fake features. Just the honest picture.

## 1. Experience principles and target outcomes

| Principle | UX obligation | Evidence |
|---|---|---|
| Beginner-first, expert-capable | Plain-language guided path with optional advanced section | New owner creates site without shell; expert inspects exact effects |
| Context over clutter | Organization, hosting account, server and resource are always evident | No accidental cross-tenant edits |
| Trustworthy state | Distinguish desired, pending and observed actual state; show last checked | Failed jobs never display as successful |
| Safe operations | Impact preview, approval, backups and rollback where supported | Destructive tests cover irreversible paths |
| No dead ends | Every failure offers retry, repair, support or diagnostic path | All catalogued states have next action |
| Accessibility by default | WCAG 2.2 AA target, keyboard and screen reader | Automated and manual audit |
| Fast perception | Skeletons, optimistic UI only for non-authoritative state | No fabricated success during async work |
| Consistent interactions | Shared primitives, form grammar, job timeline and audit | Design system adoption measured |
| Multi-tenant clarity | Scope and permissions surfaced before actions | Reseller/customer isolation E2E |
| AI with consent | AI drafts plans, never silently runs privileged actions | Approvals and action logs visible |

**Proposed measurable UX targets (validate in research):** first site completion without CLI after installation; no critical task blocked at 320px viewport; 100% of privileged mutations have visible progress and audit links; 100% of destructive operations explain retention and recovery; zero cross-tenant resources in search; all primary workflows keyboard-operable. Do not publish performance numbers before measured testing.

> **Integrated requirement:** Additional measurable UX targets from the deliberations:

- **Router-to-live in under 10 minutes** for a non-technical user on a fresh install (heuristic evaluation or moderated usability test).
- **Zero dead-end states** when a capability is unavailable — every unavailable capability shows the reason and an upgrade path.
- **Nano honest display:** a 512 MB VPS user can articulate what their server can and cannot run, from the UI alone, without reading docs.
- **Batch visibility:** every batch workload (transform, backup, build) shows its queue state, not an error, when the server is under contention.
- **Allocation/usage clarity:** every workload shows current / typical / peak; the user can answer "am I close to my limits?" from the dashboard.

## 2. Personas, contexts and tenancy UX

> **Integrated requirement:** The persona table is preserved. What changes is the frame: personas are **jobs-to-be-done**, and a user can have several simultaneously. The persona determines the **default landing**, not a permission identity. The tenant view (organization → membership → account) answers "what is this user permitted to do?" The persona view answers "what is this user currently trying to accomplish?" and drives the router.

| Persona | Landing context | Primary navigation | Boundaries |
|---|---|---|---|
| Solo VPS owner | Installation overview | Sites, Apps, Domains, Files, Databases, Backups, Server, Security | Can administer owned installation; expert privileges require step-up |
| Hosting customer | Assigned hosting account | Overview, Sites, Apps, Domains, Email, Databases, Files, Backups | No infrastructure/provider operations unless delegated |
| Developer | Authorized project/workspace | Deployments, Files, Logs, Env, Databases | Scoped to granted resources/actions |
| Reseller | Reseller organization | Customers, Packages, Accounts, Branding, Usage, Support | Cannot exceed upstream quotas/permissions |
| Provider operator | Provider organization/fleet | Fleet, Accounts, Resellers, Packages, Billing, Security, Incidents | High-risk operations gated and audited |
| Support/billing/auditor | Explicit function | Tickets/Grants, Invoices, Audit | Read-only or time-limited support as assigned |

> **Integrated requirement:** The **jobs-to-be-done** view that drives the router and the default landing:

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

**Context model:** top-left organization switcher; secondary account/project/server selector; breadcrumb identifies current scope. Switch must clear scoped search, filters, cached details and active selections. If user has multiple memberships, do not assume same role across organizations. Show a persistent banner for temporary support access, maintenance, impersonation-like delegated view or read-only mode.

> **Integrated requirement:** Additional context rules:

- **The router is context-free.** The first-run router appears before any organization/account context is established. It asks "what do you want to run?" and the answer determines the scope for the next screen.
- **The service portfolio replaces the dashboard** as the primary surface for personas whose job is "manage what I already have." It lists everything the user runs (websites, applications, self-hosted services, storage buckets, databases), each with health, usage, and quick actions. The menu (27-area IA) is the **advanced surface**, reachable from "Advanced tools."
- **The capability banner is persistent** on installations where any capability is unavailable. It is honest, not alarming: "Your VPS cannot run: WordPress (needs 150 MB more), Databases (needs 250 MB more), Image transformation (needs 2 vCPU)." Clicking expands the explanation and the upgrade path.

## 3. Global information architecture

```text
NexusPanel
├─ Overview (service portfolio)
├─ Websites & Applications
│  ├─ Websites / Staging / Redirects
│  ├─ Applications / Deployments / Processes / Workers
│  └─ Runtime & PHP/Python/Node Settings
├─ Platform Services  
│  ├─ Storage buckets / Access keys / Usage
│  ├─ Image transformation / Transform pipelines
│  ├─ Delivery / CDN integration
│  └─ Cache / Queue / Search
├─ Self-hosted Templates  
│  ├─ Template catalog
│  ├─ My self-hosted services
│  └─ Service health / backups
├─ Domains & DNS
│  ├─ Domains / Ownership / Delegation
│  ├─ DNS Zones / Records / DNSSEC
│  └─ Certificates / HTTPS
├─ Email
│  ├─ Mailboxes / Forwarders / Aliases / Lists
│  ├─ Webmail / Routing / Deliverability
│  └─ Mail Queue / Logs / Policies
├─ Data & Files
│  ├─ File Manager / Code Editor / Terminal / SFTP
│  └─ Databases / Users / Backups / Admin Tools
├─ Operations
│  ├─ Backups / Restores / Schedules
│  ├─ Monitoring / Logs / Alerts / Incidents
│  ├─ Security / Firewall / SSH / Isolation
│  └─ Jobs / Cron / Automation / Webhooks
├─ Commercial (permission-gated)
│  ├─ Hosting Accounts / Packages / Quotas
│  ├─ Customers / Resellers / Branding
│  └─ Billing / Licensing / Usage / Invoices
├─ Infrastructure (permission-gated)
│  ├─ Servers / Fleet / Service Roles / Placement
│  ├─ Agents / Updates / Maintenance
│  └─ DNS/Email/Backup/DB service nodes
├─ AI Assistant
├─ Activity & Audit
├─ Support & Documentation
└─ Settings / Users / API Tokens / Integrations / Plugins
```

> **Integrated requirement:** The information architecture includes: **Platform Services** (§3.16) and **Self-hosted Templates** (§3.17). Both sit between the primary surfaces (websites, apps) and the advanced surfaces (domains, email, data, operations), because they are the bridge between "I want to run this" and "here's the tool that makes it happen." The router (FLOW-013) routes to them.

**Navigation rules:** desktop collapsible left sidebar with pinned favorites; tablet compact rail + drilldown; mobile drawer and bottom shortcuts (Overview, Resources, Activity, Search, More). Menu labels adapt to persona but resource identity and URL structure remain stable. Unavailable module displays disabled explanatory state only when discovery is useful; prohibited resources must not be discoverable. Provide breadcrumbs, page title, resource status, scope badge and primary action in consistent locations.

> **Integrated requirement:** The menu is the **secondary** IA for personas whose job is "get something new working." On first login, the router is the primary surface. The menu is reachable via "Advanced tools" on the router or from the portfolio view. This is the single most important UX decision in the product.

## 4. Global shell and design system

### 4.1 Layout anatomy

1. **Global top bar:** brand, organization switcher, server/account context, search/command palette, activity, help, theme, user menu.
2. **Navigation:** grouped sections, active item, collapse, unread incident indicator, permission filtering.
3. **Page header:** breadcrumb, title, short description, observed status + freshness, one primary CTA, contextual actions.
4. **Main region:** summary cards, resource tables, detail tabs, wizards and editors; avoid dashboard-only dead ends.
5. **Context drawer:** non-destructive inspection, logs, inline help, AI explanation; deep links supported.
6. **Operation tray:** persistent progress across navigation, approval requests, failed jobs and recovery actions.
7. **Toast:** acknowledgement only; durable result belongs in activity center. Never use toast as sole error presentation.

> **Integrated requirement:** Additional shell components:

8. **Capability banner:** persistent on installations where any capability is unavailable. Honest, not alarming. Expands to explanation and upgrade path.
9. **Batch workload indicator:** small queue indicator in the shell showing batch jobs (transforms, backups, builds) that are queued or running. "3 jobs queued, 1 running" — clickable to the batch jobs view. This is how batch workloads are made visible without being alarming.
10. **Allocation/usage strip:** a compact strip below the page header on workload detail pages, showing allocated / typical / peak / current. This is the "am I close to my limits?" answer.

### 4.2 Visual tokens — proposed starting palette, not approved brand identity

| Token | Light | Dark | Rule |
|---|---|---|---|
| Canvas | `#F7F9FC` | `#0B1020` | Neutral backgrounds |
| Surface | `#FFFFFF` | `#121A2B` | Layered surfaces |
| Primary text | `#172033` | `#F3F6FC` | Contrast validated |
| Secondary text | `#526079` | `#B0BCD0` | Never convey status by color alone |
| Accent | `#465CF5` | `#8794FF` | Brand action; validate contrast |
| Success | `#177A53` | `#55D7A1` | Pair with icon and text |
| Warning | `#946200` | `#F4C46A` | Pair with icon and text |
| Danger | `#C23545` | `#FF8793` | Reserved for risk and failure |
| Border | `#DDE4EE` | `#2A354A` | Subtle structural boundaries |

**Typography:** modern system/Inter-like sans; clear 12/14/16/20/24/32 scale; body minimum 14px where feasible; monospaced font for logs, paths, keys, environment variables and code. **Spacing:** 4px base grid; 8/12/16/24/32/48 steps. **Radius:** 6px compact controls, 10–12px cards, 16px modal; avoid excessive decorative shadows. **Density:** comfortable default, compact for expert tables. **Motion:** 120–240ms subtle transitions; respect reduced-motion. **Icons:** consistent stroke family, always text-label destructive actions.

### 4.3 Component contract

| ID | Component | States / behavior |
|---|---|---|
| DS-001 | App shell / nav | Expanded/collapsed/mobile, active, keyboard navigation |
| DS-002 | Organization/context switcher | Search, recent, permission-aware, pending-work warning |
| DS-003 | Resource table | Sort/filter/search, cursor pagination, row selection, bulk scope validation |
| DS-004 | Status pill | Healthy/degraded/pending/offline/unknown, freshness tooltip |
| DS-005 | Capacity meter | Used/reserved/available, units, threshold, accessible text |
| DS-006 | Multi-step wizard | Stepper, back, draft, preflight, review, job handoff |
| DS-007 | Form field | Label, hint, validation, error summary, keyboard focus |
| DS-008 | Code/secret field | Mask by default, copy permissions, no analytics capture |
| DS-009 | Impact review | Before/after, affected resources, downtime, reversibility, approval |
| DS-010 | Confirmation dialog | Exact target, consequences, typed confirmation for high risk |
| DS-011 | Job timeline | Queued/running/waiting/partial/success/failed/cancelled, logs |
| DS-012 | Diff viewer | Syntax highlighting, changed lines, safe secret redaction |
| DS-013 | File tree/editor | Expand, rename, upload, tabs, dirty indicator, conflict resolution |
| DS-014 | Log viewer | Live/pause, severity, time, filter, copy/export, redact |
| DS-015 | Command palette | Scoped search, keyboard shortcuts, no permission leaks |
| DS-016 | Empty/disabled state | Reason, prerequisites, next action, no misleading upsell |
| DS-017 | AI plan card | Evidence, suggested actions, risks, approvals, execution state |
| DS-018 | Notification center | Severity, read state, deep links, per-channel preferences |
| DS-019 | Access grant banner | Scope, actor, reason, expiry, revoke link |
| DS-020 | Responsive tabs/drawers | Deep-linkable, accessible focus, preserved selection |

> **Integrated requirement:** Additional components:

| ID | Component | States / behavior |
|---|---|---|
| DS-021 | **Capability banner** | Unavailable capabilities with reason and resource delta; upgrade link; dismissible but persistent on overview |
| DS-022 | **Allocation/usage strip** | Allocated / typical / peak / current, with sparkline; per-workload and per-server |
| DS-023 | **Batch queue indicator** | Queued/running counts; clickable to batch jobs view; informational, not error |
| DS-024 | **Router card** | Six-card layout: icon, title, one-line description, "Get started" — used in the first-run router |
| DS-025 | **Live status line** | "Last checked 12 seconds ago" — visible on every resource page; distinguishes observed from desired |
| DS-026 | **Pending payment banner** | Async payment state (M-Pesa pending, Stripe processing); shows state + retry/wait action |
| DS-027 | **Template card** | ServiceTemplate catalog item: icon, name, category, resource profile, deploy button |
| DS-028 | **Upgrade recommendation** | "Consider upgrading" with evidence (peak RAM > 80%, OOM kills, CPU contention); link to upgrade flow |

### 4.4 Content and terminology

Use **Website**, **Application**, **Hosting Account**, **Server**, **DNS Zone**, **Mailbox**, **Backup**, **Restore**, **Deployment**, **Job**, **Approval** consistently. Do not use "domain" interchangeably with website. Use plain language first, then expert terms in expandable help. Every destructive label must name its target: "Delete website" rather than "Delete". Distinguish **saved**, **submitted**, **running**, **verified**. Show absolute timestamp and relative age where useful; dates displayed in chosen timezone, audits in UTC.

> **Integrated requirement:** Additional terminology:

- **"Platform Service"** — storage, transform, delivery, cache, queue, search. These are services the user's other services consume, provisioned from `ServiceTemplate` records.
- **"Self-hosted Template"** — a `ServiceTemplate` instance (Nextcloud, Gitea, Ghost, n8n, Uptime Kuma, etc.) that the user provisions as a self-hosted app.
- **"Batch workload"** — image transform, backup, build, AI diagnostics, migration. Runs low-priority, queues under contention, completes eventually.
- **"Allocation"** — the ceiling, not a reservation. Say "this workload may use up to X" not "this workload has X."
- **"Typical" / "Peak"** — measured, not assumed. Say "typical 0.4 vCPU (7-day p50), peak 1.5 vCPU (24h p99)."
- **"Capability unavailable"** — the installer determined this VPS cannot run this. Not "disabled," not "not available on your plan." Say "your VPS cannot run this because it needs 150 MB more RAM."

## 5. Global interaction and system states

| State | Presentation | Available actions |
|---|---|---|
| Initial loading | Skeleton matching content, accessible loading announcement | Cancel navigation where safe |
| Empty | Explain why empty, contextual first action | Create/import or learn |
| Loading more | Inline progress, preserve previous rows | Cancel filter or retry |
| Pending | Visible job ID and stage, no false green success | View job, safe cancel if supported |
| Awaiting approval | Approver, reason, expiry, proposed impact | Approve/reject only with permission |
| Success (verified) | Observed result and timestamp | View resource, audit, next step |
| Partial success | Per-step outcomes, remaining risk | Retry failed steps, compensating plan |
| Failed | Safe error code, reason, correlation ID | Retry when safe, diagnostics, support |
| Offline/stale | Last seen, no live-health claim | Refresh, inspect last state |
| Unauthorized | Generic denial, no resource existence leak | Request access when permitted |
| Unsupported | Capability mismatch and prerequisites | Change topology/adapter |
| Quota exceeded | Current/reserved/limit, upgrade/request flow | Free capacity or request quota |
| Conflict/concurrent edit | Show versions/diff and owner if allowed | Reload, merge, force only with policy |
| Maintenance | Banner, estimated window if known | View status, schedule permitted action |
| Recovery required | Prominent incident status and out-of-band runbook | Guided recovery, support |

**Mutations:** disable duplicate submit; use idempotency key; retain draft during navigation; browser refresh must rehydrate operation status. Form validation is both inline and summarized; server validation remains authoritative. Bulk operations show count, scope and per-item results. High-risk actions use step-up MFA and/or second-person approval when policy requires. Avoid optimistic "success" for DNS propagation, TLS issuance, deployment health or restore.

> **Integrated requirement:** Additional states from the deliberations:

| State | Presentation | Available actions |
|---|---|---|
| **Capability unavailable** | Honest explanation with resource delta and upgrade link | View upgrade path, choose alternative |
| **Batch queued** | "Will run when server is idle" with queue position | View queue, remove from queue if supported |
| **Batch running** | Progress bar with low priority indicator | Pause if supported (rare) |
| **Batch paused** | Reason (serving workload contention) | Resume when idle, or force (with warning) |
| **Memory pressure** | "This workload is close to its memory limit; consider upgrading" | View metrics, adjust workload, upgrade |
| **Throttled** | "This batch workload is yielding to site traffic" — informational, not error | View metrics, no action required |
| **Pending payment** | Async payment state (M-Pesa pending, Stripe processing) | Wait, retry, cancel |
| **Payment callback unmatched** | "We received a payment notification we couldn't match to an invoice" | Manual reconciliation (admin only) |
| **Router no-match** | "Let's see what your VPS can do first" | Fall back to capability display + guided tour |

## 6. Search, notifications, help and command palette

- **Global search:** keyboard shortcut, typeahead grouped by sites/apps/domains/accounts/servers/jobs/docs; search only authorized scopes; result contains type, name, context and status; no hidden tenant hints.
- **Command palette:** "Create website", "View failed jobs", "Open server logs", "Add DNS record" filtered by capability and permission. A command requiring mutation opens normal guarded wizard; no silent execution.
- **Notifications:** in-app plus optional email/webhook; severity and deduplication; incident and job links; user preferences with mandatory critical notifications policy.
- **Help:** contextual inline explanations, docs links, CLI equivalents when safe, glossary, troubleshooting steps and support ticket entry.
- **Activity:** user-visible operation stream with actor, scope, submitted/observed state, approvals, correlation ID and safe logs.

> **Integrated requirement:** Additional capabilities:

- **Global search** includes **platform services** (storage buckets, self-hosted services) and **templates** ("deploy Nextcloud").
- **Command palette** includes **capability-aware commands**: "Deploy Nextcloud" appears only if the capability is available; if unavailable, the command shows the reason and the upgrade path.
- **Notifications** include **batch completion** ("Your backup completed successfully") and **batch queued** ("Your backup is queued, will run when the server is idle").
- **Help** includes a **"What can my VPS run?"** article that reads the capability set and explains it in plain language. This is the primary lay-person help entry point.

## 7. Screen specification template

Each catalogued screen below uses this minimum implementation contract: **route**, **audience and permission**, **purpose**, **data required**, **primary/secondary actions**, **field validation**, **states**, **audit/analytics**, **responsive and accessibility behavior**, **Stage 2 mapping**. Screen IDs are stable; routes are *proposed patterns* pending API/URL architecture approval. All screens inherit the global state matrix above. "Delete" and other dangerous actions inherit impact-review rules.

## 8. Exhaustive screen catalog and per-module interaction contracts

All routes below are proposed. Each screen inherits global loading, empty, error, stale, permission and job states. Detailed flows and field specifications follow the catalog.

> **Integrated requirement:** The screen catalog is preserved in full. **New screens** are added for Platform Services, Self-hosted Templates, the Router, the Service Portfolio, and the Allocation/Usage views. The **screen count grows from 138 to approximately 165**. The lay-person journey (FLOW-001) is reframed as the Router (§9), and the Nano experience (FLOW-013) is added.

### 8.1 Identity, onboarding and settings

**Stage 2 traceability:** `NP-IAM, NP-UX`. **Permission rule:** audience labels are descriptive; server policy resolves the actual actor/resource/action/scope at runtime.

| Screen ID | Screen / route | Audience | Required content and primary interactions |
|---|---|---|---|
| UX-IAM-001 | **Welcome / installation claim**  `/setup/welcome` | Owner | Start one-time enrollment; show system prerequisites and claim expiry. **Actions:** Start setup, view docs. |
| UX-IAM-002 | **Owner account and MFA**  `/setup/owner` | Owner | Create owner identity, password/passkey and MFA. **Actions:** Verify, continue. |
| UX-IAM-003 | **First-run router**  `/onboarding` | Owner, Customer, Reseller | **NEW in v2:** "What do you want to run?" — six-card router. See §9 FLOW-001 for detail. **Actions:** Pick a card, or "I'm not sure yet." |
| UX-IAM-004 | **Sign-in and recovery**  `/auth/sign-in` | All | Secure authentication, MFA, recovery without exposing accounts. **Actions:** Sign in, recover. |
| UX-IAM-005 | **Organizations**  `/settings/organizations` | Owner, Provider | List owned organizations and memberships. **Actions:** Create, switch. |
| UX-IAM-006 | **Organization detail**  `/settings/organizations/:org` | Authorized members | Manage name, roles, limits, policy. **Actions:** Edit, invite. |
| UX-IAM-007 | **Users and invitations**  `/settings/members` | Owner, Provider, Reseller | Member list, invitations, expiry and revocation. **Actions:** Invite, revoke. |
| UX-IAM-008 | **Roles and access policies**  `/settings/access` | Authorized admins | Permissions, inherited scopes and effective access preview. **Actions:** Create role, test permission. |
| UX-IAM-009 | **Sessions and devices**  `/settings/sessions` | All | Session inventory, last active and revoke. **Actions:** Revoke session. |
| UX-IAM-010 | **API tokens**  `/settings/api-tokens` | Authorized | Token scopes, expiry, last used; one-time reveal. **Actions:** Create, rotate, revoke. |
| UX-IAM-011 | **Support access grants**  `/settings/support-access` | Owner, Provider, Support | Time-limited scope and justification. **Actions:** Request, approve, revoke. |
| UX-IAM-012 | **Personal preferences**  `/settings/profile` | All | Theme, language, timezone, density, notifications. **Actions:** Save preferences. |

**Module UX acceptance:** Each screen must have a named owner, authorized route guard, scope breadcrumb, read-only and mutating variants where applicable, specific empty/unsupported/failed states, keyboard test, responsive test, audit correlation and E2E scenario. Resource actions launch impact review and durable jobs; UI cannot mark completion before verification. The module's screens must be cross-linked to `NP-IAM, NP-UX` requirement IDs individually during implementation planning.

> **Integrated requirement:** UX-IAM-003 is reframed from "First-run checklist" to "First-run router." The checklist exists as a follow-up after the router routes to a path; it is not the first screen.

### 8.2 Overview and navigation

**Stage 2 traceability:** `NP-UX, NP-CORE`. **Permission rule:** audience labels are descriptive; server policy resolves the actual actor/resource/action/scope at runtime.

| Screen ID | Screen / route | Audience | Required content and primary interactions |
|---|---|---|---|
| UX-OV-001 | **Service portfolio**  `/overview` | All | **NEW in v2:** the primary surface for personas whose job is "manage what I already have." Lists everything the user runs (websites, apps, self-hosted services, storage buckets, databases), each with health, usage, quick actions. **Actions:** Open, create, upgrade. |
| UX-OV-002 | **Capability overview**  `/overview/capabilities` | All | **NEW in v2:** honest display of what this VPS can and cannot run, with resource deltas and upgrade link. Replaces the v1 "Overview" dashboard for Nano/Starter installations. **Actions:** Upgrade, learn more. |
| UX-OV-003 | **Resource search**  `/search` | All | Permission-scoped search across resources (including platform services and templates). **Actions:** Open result. |
| UX-OV-004 | **Activity and jobs**  `/activity` | All | Operations, approvals, outcomes and progress (including batch jobs). **Actions:** Inspect, retry. |
| UX-OV-005 | **Job detail**  `/activity/jobs/:job` | Authorized | Timeline, logs, impact, retry safety. **Actions:** Retry, cancel when supported. |
| UX-OV-006 | **Approvals inbox**  `/approvals` | Approvers | Pending approvals with impact and expiry. **Actions:** Approve, reject. |
| UX-OV-007 | **Notifications**  `/notifications` | All | Alerts, read status, preferences. **Actions:** Open, mark read. |
| UX-OV-008 | **Help center**  `/help` | All | Contextual docs, glossary, troubleshooting, **"What can my VPS run?"** article. **Actions:** Search, contact support. |
| UX-OV-009 | **Allocation and usage**  `/overview/usage` | All | **NEW in v2:** server-level view of allocation / typical / peak / current across all workloads. **Actions:** Drill down, adjust ceilings, upgrade. |
| UX-OV-010 | **Upgrade recommendations**  `/overview/recommendations` | All | **NEW in v2:** auto-generated upgrade suggestions (peak RAM > 80% for 7 days, OOM kills, CPU contention). **Actions:** View evidence, upgrade. |

**Module UX acceptance:** Each screen must have a named owner, authorized route guard, scope breadcrumb, read-only and mutating variants where applicable, specific empty/unsupported/failed states, keyboard test, responsive test, audit correlation and E2E scenario. Resource actions launch impact review and durable jobs; UI cannot mark completion before verification. The module's screens must be cross-linked to `NP-UX, NP-CORE` requirement IDs individually during implementation planning.

> **Integrated requirement:** The "Role-aware overview" (v1 UX-OV-001) becomes the **service portfolio** (UX-OV-001). This is the primary surface for personas whose job is "manage what I already have." The v1 "Overview" dashboard becomes the **capability overview** (UX-OV-002), which is only shown on Nano/Starter installations or when a user navigates to "What can my VPS run?" It is not the default landing for any persona.

### 8.3 Servers, fleet and service roles

**Stage 2 traceability:** `NP-INF`. **Permission rule:** audience labels are descriptive; server policy resolves the actual actor/resource/action/scope at runtime.

| Screen ID | Screen / route | Audience | Required content and primary interactions |
|---|---|---|---|
| UX-INF-001 | **Server inventory**  `/infrastructure/servers` | Owner, Provider | Health, version, capacity, roles, region. **Actions:** Enroll, filter. |
| UX-INF-002 | **Server detail**  `/infrastructure/servers/:server` | Authorized admins | CPU/RAM/disk/IPs, heartbeat, modules, incidents. **Actions:** Manage, diagnose. |
| UX-INF-003 | **Enroll server wizard**  `/infrastructure/servers/new` | Owner, Provider | Enrollment token, prerequisites, trust verification. **Actions:** Generate token, approve. |
| UX-INF-004 | **Service roles**  `/infrastructure/roles` | Provider | Role assignment, conflicts, capacity. **Actions:** Assign, review. |
| UX-INF-005 | **Fleet topology**  `/infrastructure/topology` | Provider | Nodes, zones, roles, failure domains. **Actions:** Inspect, plan. |
| UX-INF-006 | **Placement planner**  `/infrastructure/placement` | Provider | Eligible servers and capacity reservation. **Actions:** Preview, approve. |
| UX-INF-007 | **Maintenance and drain**  `/infrastructure/maintenance` | Provider | Maintenance windows, workload impact, drain progress. **Actions:** Schedule, cancel. |
| UX-INF-008 | **Agent health and upgrades**  `/infrastructure/agents` | Provider | Versions, trust, drift, upgrade jobs. **Actions:** Plan upgrade, recover. |
| UX-INF-009 | **Installed modules**  `/infrastructure/modules` | Owner, Provider | Adapters, dependencies, ownership, install status. **Actions:** Install, disable. |
| UX-INF-010 | **Server diagnostics**  `/infrastructure/servers/:server/diagnostics` | Authorized admins | Read-only health evidence and runbook. **Actions:** Run safe checks. |

> **Integrated requirement:** Additional screens:

| Screen ID | Screen / route | Audience | Required content and primary interactions |
|---|---|---|---|
| UX-INF-011 | **Capability detection**  `/infrastructure/capabilities` | Owner, Provider | **NEW in v2:** detailed view of the installer's 5-step detection. Shows detected hardware, floor/typical/ceiling per capability, overcommit warnings, and blockers (port 25, PTR, no public IP). **Actions:** Rescan, view evidence, export. |
| UX-INF-012 | **Hardware rescan**  `/infrastructure/rescan` | Owner, Provider | **NEW in v2:** post-upgrade re-detection. Shows new capabilities that became available after hardware change. **Actions:** Rescan, apply new capability set. |
| UX-INF-013 | **cgroup metrics**  `/infrastructure/servers/:server/cgroups` | Authorized admins | **NEW in v2:** per-workload cgroup metrics (cpu.stat, memory.current, memory.peak, memory.events, io.stat, pids.current). Current / typical / peak. **Actions:** Drill down, adjust ceilings. |

**Module UX acceptance:** Each screen must have a named owner, authorized route guard, scope breadcrumb, read-only and mutating variants where applicable, specific empty/unsupported/failed states, keyboard test, responsive test, audit correlation and E2E scenario. Resource actions launch impact review and durable jobs; UI cannot mark completion before verification. The module's screens must be cross-linked to `NP-INF` requirement IDs individually during implementation planning.

### 8.4 Hosting accounts, reseller and commercial operations

**Stage 2 traceability:** `NP-ACC, NP-BILL`. **Permission rule:** audience labels are descriptive; server policy resolves the actual actor/resource/action/scope at runtime.

| Screen ID | Screen / route | Audience | Required content and primary interactions |
|---|---|---|---|
| UX-COM-001 | **Hosting accounts**  `/commercial/accounts` | Provider, Reseller | Owned accounts, state, package, limits. **Actions:** Create, suspend. |
| UX-COM-002 | **Account detail**  `/commercial/accounts/:account` | Authorized | Sites, users, quotas, subscriptions, health. **Actions:** Manage, transfer. |
| UX-COM-003 | **Account creation wizard**  `/commercial/accounts/new` | Provider, Reseller | Customer, package, region, limits, approval. **Actions:** Review, provision. |
| UX-COM-004 | **Packages**  `/commercial/packages` | Provider, Reseller | Resource quotas, features, pricing rules. **Actions:** Create, clone. |
| UX-COM-005 | **Package editor**  `/commercial/packages/:package` | Authorized | Limits, overrides, versions, entitlements. **Actions:** Save, publish. |
| UX-COM-006 | **Quota and usage**  `/commercial/usage` | Authorized | Metered/reserved/actual usage, source freshness. **Actions:** Inspect, export. |
| UX-COM-007 | **Customers**  `/commercial/customers` | Provider, Reseller | Customer organizations, contacts and services. **Actions:** Invite, assign. |
| UX-COM-008 | **Resellers**  `/commercial/resellers` | Provider | Reseller subtree, grants, limits, branding. **Actions:** Create, suspend. |
| UX-COM-009 | **Reseller portal branding**  `/commercial/branding` | Reseller, Provider | Logo, domain, colors, email templates preview. **Actions:** Preview, publish. |
| UX-COM-010 | **Subscriptions**  `/commercial/subscriptions` | Billing authorized | Plans, renewals, state, grace periods. **Actions:** Adjust, cancel. |
| UX-COM-011 | **Invoices and payments**  `/commercial/invoices` | Billing authorized | Invoices, payments, refunds, reconciliation. **Actions:** View, export. |
| UX-COM-012 | **Licensing**  `/commercial/licenses` | Provider | License state, seats/nodes, activation. **Actions:** Activate, renew. |
| UX-COM-013 | **Commerce integrations**  `/commercial/integrations` | Provider | Gateway adapters, credentials, test mode. **Actions:** Connect, test. |

> **Integrated requirement:** Significant refinement — the commercial UI reflects the async-first payment state machine:

| Screen ID | Screen / route | Audience | Required content and primary interactions |
|---|---|---|---|
| UX-COM-014 | **Payment methods**  `/commercial/payment-methods` | Billing authorized | **NEW in v2:** managed payment methods (M-Pesa phone, Stripe card). Not raw card data — tokenized references. **Actions:** Add, remove, set default. |
| UX-COM-015 | **Payment pending**  `/commercial/payments/:id/pending` | Billing authorized | **NEW in v2:** the async payment waiting state. Shows "Awaiting M-Pesa confirmation — check your phone" or "Processing card payment" with appropriate timing. **Actions:** Wait, retry, cancel. |
| UX-COM-016 | **Payment history**  `/commercial/payments` | Billing authorized | **NEW in v2:** list of payments with state (pending/confirmed/failed/reversed), provider reference, reconciliation status. **Actions:** Filter, export, reconcile unmatched. |
| UX-COM-017 | **Payment reconciliation**  `/commercial/payments/reconcile` | Admin only | **NEW in v2:** manual reconciliation for unmatched callbacks (webhook received but no matching invoice, or invoice marked paid but no callback). **Actions:** Match, refund, dispute. |
| UX-COM-018 | **M-Pesa settings**  `/commercial/integrations/mpesa` | Admin only | **NEW in v2:** Daraja credentials (sandbox + production), callback URL, test transaction. **Actions:** Test, save, enable. |

**Module UX acceptance:** Each screen must have a named owner, authorized route guard, scope breadcrumb, read-only and mutating variants where applicable, specific empty/unsupported/failed states, keyboard test, responsive test, audit correlation and E2E scenario. Resource actions launch impact review and durable jobs; UI cannot mark completion before verification. The module's screens must be cross-linked to `NP-ACC, NP-BILL` requirement IDs individually during implementation planning.

> **Integrated requirement:** The commercial surface is **deferred to v2+ for customer-facing UI**, but the **data model and API are v1**. A v1 installation with a single owner sees no commercial surface; the owner uses the license activation flow (UX-COM-012) to activate the panel. When v2 ships the reseller/provider surface, these screens become active without architecture change.

### 8.5 Websites, domains and application deployments

**Stage 2 traceability:** `NP-WEB, NP-APP`. **Permission rule:** audience labels are descriptive; server policy resolves the actual actor/resource/action/scope at runtime.

| Screen ID | Screen / route | Audience | Required content and primary interactions |
|---|---|---|---|
| UX-WEB-001 | **Websites list**  `/websites` | Owner, Customer, Developer | Domain, runtime, SSL, health, account. **Actions:** Create, filter. |
| UX-WEB-002 | **Create website wizard**  `/websites/new` | Authorized | Domain, root, runtime, placement, TLS, preview. **Actions:** Validate, provision. |
| UX-WEB-003 | **Website overview**  `/websites/:site` | Authorized | Observed health, URLs, resource usage, recent jobs. **Actions:** Open, manage. |
| UX-WEB-004 | **Website configuration**  `/websites/:site/settings` | Authorized | Document root, redirects, PHP/HTTP settings, staging. **Actions:** Preview, apply. |
| UX-WEB-005 | **Staging environments**  `/websites/:site/staging` | Authorized | Staging copies, sync, promotion impact. **Actions:** Create, promote. |
| UX-WEB-006 | **Redirects and aliases**  `/websites/:site/redirects` | Authorized | Rules, loops, conflicts, test results. **Actions:** Create, test. |
| UX-WEB-007 | **Applications list**  `/apps` | Owner, Customer, Developer | Runtime, process, deployments, health. **Actions:** New app. |
| UX-WEB-008 | **New app wizard**  `/apps/new` | Authorized | Git/upload/existing source, runtime, build/start, env, URL. **Actions:** Preflight, deploy. |
| UX-WEB-009 | **Application overview**  `/apps/:app` | Authorized | URL, build, runtime, process, recent deploys. **Actions:** Deploy, restart. |
| UX-WEB-010 | **Deployment history**  `/apps/:app/deployments` | Authorized | Commit, actor, logs, health, rollback eligibility. **Actions:** Compare, rollback. |
| UX-WEB-011 | **Deployment detail**  `/apps/:app/deployments/:deploy` | Authorized | Build/run timeline, errors, artifact, verification. **Actions:** Inspect, retry. |
| UX-WEB-012 | **Environment variables**  `/apps/:app/environment` | Authorized | Scoped secrets, paste/import .env, duplicate validation. **Actions:** Import, edit, rotate. |
| UX-WEB-013 | **Runtime configuration**  `/apps/:app/runtime` | Authorized | Supported versions, startup command, ports, limits. **Actions:** Validate, save. |
| UX-WEB-014 | **Processes and workers**  `/apps/:app/processes` | Authorized | Process states, restart policy, queue workers. **Actions:** Scale, restart. |
| UX-WEB-015 | **Git connections**  `/apps/:app/git` | Authorized | Repo, branch, webhook, deploy key permissions. **Actions:** Connect, disconnect. |
| UX-WEB-016 | **Application logs**  `/apps/:app/logs` | Authorized | Structured logs, streaming, filters, redact. **Actions:** Search, export. |
| UX-WEB-017 | **Domains attached to resource**  `/apps/:app/domains` | Authorized | Hostnames, routing, TLS status. **Actions:** Attach, detach. |

> **Integrated requirement:** Significant additions and reframings:

| Screen ID | Screen / route | Audience | Required content and primary interactions |
|---|---|---|---|
| UX-WEB-018 | **Domain connection wizard**  `/domains/connect` | Authorized | **NEW in v2 (promoted from settings):** the killer feature for the lay person. Shows exact A/AAAA/CNAME record to add at registrar, polls for propagation, confirms when live, then offers SSL. **Actions:** Copy record, check propagation, continue. |
| UX-WEB-019 | **Website capability strip**  `/websites/:site/capabilities` | Authorized | **NEW in v2:** shows what the site can and cannot do (PHP? Node? DB? Transform?) based on installer capability set + container fit. **Actions:** Upgrade, adjust. |
| UX-WEB-020 | **Batch jobs for site**  `/websites/:site/batch` | Authorized | **NEW in v2:** shows the site's batch jobs (builds, backups, transforms queued or running). **Actions:** View queue, cancel if supported. |

**Module UX acceptance:** Each screen must have a named owner, authorized route guard, scope breadcrumb, read-only and mutating variants where applicable, specific empty/unsupported/failed states, keyboard test, responsive test, audit correlation and E2E scenario. Resource actions launch impact review and durable jobs; UI cannot mark completion before verification. The module's screens must be cross-linked to `NP-WEB, NP-APP` requirement IDs individually during implementation planning.

> **Integrated requirement:** On Nano and Starter profiles, PHP/Node/DB screens are **not shown**. The website creation wizard adapts: on Nano, "Create website" only offers static; on Starter, it offers static and one of {PHP, Node, small DB}. The capability strip (UX-WEB-019) explains why.

### 8.6 Files, editor and developer access

**Stage 2 traceability:** `NP-FILE`. **Permission rule:** audience labels are descriptive; server policy resolves the actual actor/resource/action/scope at runtime.

| Screen ID | Screen / route | Audience | Required content and primary interactions |
|---|---|---|---|
| UX-FILE-001 | **File manager**  `/files` | Authorized | Scoped tree, breadcrumbs, permissions, usage. **Actions:** Upload, new, rename. |
| UX-FILE-002 | **File detail and preview**  `/files/view` | Authorized | Metadata, MIME preview, safe download. **Actions:** Download, inspect. |
| UX-FILE-003 | **Code editor**  `/files/edit` | Authorized | Tabs, syntax, save, lint if supported, conflict diff. **Actions:** Save, revert. |
| UX-FILE-004 | **Bulk file operations**  `/files/operations` | Authorized | Multi-select, archive, move, delete impact. **Actions:** Review, submit. |
| UX-FILE-005 | **Upload manager**  `/files/uploads` | Authorized | Drag-drop, resumable status, conflicts, size limits. **Actions:** Upload, retry. |
| UX-FILE-006 | **File permissions**  `/files/permissions` | Authorized | Ownership, ACL/mode preview, security implications. **Actions:** Review, apply. |
| UX-FILE-007 | **SFTP/SSH access**  `/developer/access` | Authorized | Keys, access scopes, expiry and revocation. **Actions:** Add key, revoke. |
| UX-FILE-008 | **Web terminal**  `/developer/terminal` | Explicit privileged scope | Session warning, scope, audit, inactivity. **Actions:** Open, close. |
| UX-FILE-009 | **Cron and scheduled tasks**  `/developer/schedules` | Authorized | Expression, timezone, last/next run, logs. **Actions:** Create, pause. |

> **Integrated requirement:** Additional clarifications and screens:

| Screen ID | Screen / route | Audience | Required content and primary interactions |
|---|---|---|---|
| UX-FILE-010 | **Terminal unavailable**  `/developer/terminal` (state) | Authorized | **NEW in v2:** honest state on Nano — "Terminal is not available on this VPS because sites are not containerized. SFTP works for file transfer." **Actions:** View SFTP, upgrade. |
| UX-FILE-011 | **SFTP setup**  `/developer/sftp` | Authorized | **NEW in v2:** panel-managed SFTP endpoint, per-site credentials, key upload. This is the primary file-transfer path on Nano. **Actions:** Copy SFTP details, add key. |

**Module UX acceptance:** Each screen must have a named owner, authorized route guard, scope breadcrumb, read-only and mutating variants where applicable, specific empty/unsupported/failed states, keyboard test, responsive test, audit correlation and E2E scenario. Resource actions launch impact review and durable jobs; UI cannot mark completion before verification. The module's screens must be cross-linked to `NP-FILE` requirement IDs individually during implementation planning.

> **Integrated requirement:** Terminal is **disabled by default** for customer roles (NP-FILE-013). On Nano, terminal is not available at all. SFTP is the primary file-transfer path on Nano; on Standard+, both SFTP and container shell are available.

### 8.7 Databases and administration

**Stage 2 traceability:** `NP-DB`. **Permission rule:** audience labels are descriptive; server policy resolves the actual actor/resource/action/scope at runtime.

| Screen ID | Screen / route | Audience | Required content and primary interactions |
|---|---|---|---|
| UX-DB-001 | **Databases list**  `/databases` | Authorized | Engine, version, host, size, backups. **Actions:** Create. |
| UX-DB-002 | **Create database wizard**  `/databases/new` | Authorized | Engine, name, owner, server, quotas. **Actions:** Review, create. |
| UX-DB-003 | **Database overview**  `/databases/:db` | Authorized | Connection, users, size, health, recent operations. **Actions:** Manage, backup. |
| UX-DB-004 | **Database users and grants**  `/databases/:db/users` | Authorized | Least-privilege grants, rotation, host restrictions. **Actions:** Create, rotate. |
| UX-DB-005 | **Database administration**  `/databases/:db/admin` | Authorized | Safe admin-tool launch, restricted scopes. **Actions:** Launch tool. |
| UX-DB-006 | **Import/export jobs**  `/databases/:db/transfers` | Authorized | File/source, size, checks, status, impact. **Actions:** Import, export. |
| UX-DB-007 | **Database backups and restores**  `/databases/:db/backups` | Authorized | Snapshots, consistency, restore target. **Actions:** Restore preview. |
| UX-DB-008 | **Database service settings**  `/infrastructure/databases` | Provider | Engine modules, capacity, maintenance, upgrades. **Actions:** Plan change. |

> **Integrated requirement:** Significant scope change — on Nano databases are unavailable, on Starter they are borderline, on Standard+ they ship in v2:

| Screen ID | Screen / route | Audience | Required content and primary interactions |
|---|---|---|---|
| UX-DB-009 | **Databases unavailable**  `/databases` (state) | Authorized | **NEW in v2:** honest state on Nano — "Databases are not available on this VPS because the panel's metadata DB already consumes the available SQLite budget. Upgrade to Standard to enable MariaDB/Postgres." **Actions:** Upgrade, learn more. |
| UX-DB-010 | **Database as ServiceTemplate**  `/databases/:db/service` | Authorized | **NEW in v2:** shows the database as a service instance (container, volume, credentials, health check, backup policy). This unifies the UI with the ServiceTemplate primitive. **Actions:** Manage as service, export. |

**Module UX acceptance:** Each screen must have a named owner, authorized route guard, scope breadcrumb, read-only and mutating variants where applicable, specific empty/unsupported/failed states, keyboard test, responsive test, audit correlation and E2E scenario. Resource actions launch impact review and durable jobs; UI cannot mark completion before verification. The module's screens must be cross-linked to `NP-DB` requirement IDs individually during implementation planning.

### 8.8 Domains, DNS and SSL

**Stage 2 traceability:** `NP-DNS, NP-TLS`. **Permission rule:** audience labels are descriptive; server policy resolves the actual actor/resource/action/scope at runtime.

| Screen ID | Screen / route | Audience | Required content and primary interactions |
|---|---|---|---|
| UX-DNS-001 | **Domains inventory**  `/domains` | Authorized | Ownership, registrar, delegation, SSL, expiration. **Actions:** Add domain. |
| UX-DNS-002 | **Domain onboarding wizard**  `/domains/new` | Authorized | Domain, ownership, DNS mode, zone, nameservers. **Actions:** Preflight, activate. |
| UX-DNS-003 | **Domain detail**  `/domains/:domain` | Authorized | DNS status, verification, attached sites/apps. **Actions:** Manage. |
| UX-DNS-004 | **Ownership verification**  `/domains/:domain/verify` | Authorized | TXT challenge, TTL, polling, errors. **Actions:** Verify, retry. |
| UX-DNS-005 | **Delegation checklist**  `/domains/:domain/delegation` | Authorized | Registrar NS/glue, observed NS, propagation. **Actions:** Recheck. |
| UX-DNS-006 | **DNS zone editor**  `/dns/zones/:zone` | Authorized | Record table, TTL, conflicts, syntax, diff. **Actions:** Add, edit, delete. |
| UX-DNS-007 | **DNS zone change review**  `/dns/zones/:zone/review` | Authorized | Before/after, affected records, propagation. **Actions:** Approve, apply. |
| UX-DNS-008 | **DNSSEC**  `/dns/zones/:zone/dnssec` | Authorized | Keys, DS records, registrar instructions, rollover. **Actions:** Enable, rotate. |
| UX-DNS-009 | **DNS infrastructure**  `/infrastructure/dns` | Provider | Primary/secondary, health, replication, zones. **Actions:** Manage roles. |
| UX-DNS-010 | **SSL certificates**  `/certificates` | Authorized | Issued/pending/expiring, SANs, renewal. **Actions:** Request, inspect. |
| UX-DNS-011 | **Certificate issuance wizard**  `/certificates/new` | Authorized | Domains, ACME challenge, prerequisites, validation. **Actions:** Issue. |
| UX-DNS-012 | **Certificate detail**  `/certificates/:cert` | Authorized | Chain, issuer, expiry, renewals, assignments. **Actions:** Renew, revoke. |
| UX-DNS-013 | **HTTPS policy**  `/websites/:site/https` | Authorized | Redirect, HSTS, TLS modes and warnings. **Actions:** Preview, apply. |

> **Integrated requirement:** Significant scope change — v1 is DNS guidance mode, not authoritative hosting:

| Screen ID | Screen / route | Audience | Required content and primary interactions |
|---|---|---|---|
| UX-DNS-014 | **DNS guidance**  `/domains/:domain/dns-guidance` | Authorized | **NEW in v2:** the v1 DNS mode. Shows exact records to add at registrar, propagation polling, error diagnosis. This is the primary DNS experience for v1. **Actions:** Copy records, poll propagation, fix errors. |
| UX-DNS-015 | **Authoritative DNS unavailable**  `/dns/zones` (state) | Authorized | **NEW in v2:** honest state on v1 — "Authoritative DNS hosting is coming in v2. For now, use the DNS guidance flow to point your domain at this server." **Actions:** Use guidance, learn more. |

**Module UX acceptance:** Each screen must have a named owner, authorized route guard, scope breadcrumb, read-only and mutating variants where applicable, specific empty/unsupported/failed states, keyboard test, responsive test, audit correlation and E2E scenario. Resource actions launch impact review and durable jobs; UI cannot mark completion before verification. The module's screens must be cross-linked to `NP-DNS, NP-TLS` requirement IDs individually during implementation planning.

> **Integrated requirement:** Authoritative DNS hosting is deferred to v2+. v1 uses **DNS guidance mode** (UX-DNS-014) as the primary flow: "point this A record at this IP, we'll poll until it propagates." This is the killer feature for the lay person, and it does not require authoritative hosting. Zone editor, DNSSEC, and secondary DNS are v2+.

### 8.9 Email hosting and deliverability

**Stage 2 traceability:** `NP-MAIL`. **Permission rule:** audience labels are descriptive; server policy resolves the actual actor/resource/action/scope at runtime.

| Screen ID | Screen / route | Audience | Required content and primary interactions |
|---|---|---|---|
| UX-MAIL-001 | **Email domain setup**  `/email/domains/new` | Authorized | Mail role readiness, MX/SPF/DKIM/DMARC checks. **Actions:** Configure, verify. |
| UX-MAIL-002 | **Email domains**  `/email/domains` | Authorized | Domains, health, DNS, deliverability. **Actions:** Add, manage. |
| UX-MAIL-003 | **Mailboxes**  `/email/mailboxes` | Authorized | Accounts, quota, status, access. **Actions:** Create, suspend. |
| UX-MAIL-004 | **Mailbox detail**  `/email/mailboxes/:box` | Authorized | Quota, password, aliases, forwarding, logs. **Actions:** Edit, reset. |
| UX-MAIL-005 | **Aliases and forwarders**  `/email/aliases` | Authorized | Source/target, loops, restrictions. **Actions:** Create, test. |
| UX-MAIL-006 | **Mailing lists**  `/email/lists` | Authorized | Members, moderation, delivery policy. **Actions:** Create, edit. |
| UX-MAIL-007 | **Webmail launcher**  `/email/webmail` | Authorized | SSO or separate login, safe redirect. **Actions:** Open webmail. |
| UX-MAIL-008 | **Mail DNS health**  `/email/deliverability` | Authorized | SPF/DKIM/DMARC/PTR, reputation checks. **Actions:** Diagnose. |
| UX-MAIL-009 | **Mail queue and delivery logs**  `/email/queue` | Authorized | Deferred/bounced messages, redacted logs. **Actions:** Retry, inspect. |
| UX-MAIL-010 | **Mail server configuration**  `/infrastructure/mail` | Provider | Ports, TLS, spam, quotas, backups. **Actions:** Plan changes. |

> **Integrated requirement:** Significant scope change — v1 is external SMTP, self-hosted mail is v2+:

| Screen ID | Screen / route | Audience | Required content and primary interactions |
|---|---|---|---|
| UX-MAIL-011 | **External SMTP settings**  `/email/external-smtp` | Authorized | **NEW in v2:** configure Resend / Postmark / Mailgun / SES. Used for panel notifications and WordPress outbound mail. **Actions:** Add credentials, test, enable. |
| UX-MAIL-012 | **Mailbox recommendation**  `/email/recommendation` | Authorized | **NEW in v2:** honest guidance — "For mailboxes, we recommend Zoho or Google Workspace. Self-hosted mail is coming in v2." **Actions:** Open Zoho/Google, learn more. |

**Module UX acceptance:** Each screen must have a named owner, authorized route guard, scope breadcrumb, read-only and mutating variants where applicable, specific empty/unsupported/failed states, keyboard test, responsive test, audit correlation and E2E scenario. Resource actions launch impact review and durable jobs; UI cannot mark completion before verification. The module's screens must be cross-linked to `NP-MAIL` requirement IDs individually during implementation planning.

> **Integrated requirement:** Self-hosted email is deferred to v2+. v1 uses external SMTP (UX-MAIL-011) for outbound mail and recommends Zoho/Google Workspace for mailboxes (UX-MAIL-012). On Nano, email is **unavailable** — the UI does not offer it.

### 8.10 Backups, monitoring, incidents and security

**Stage 2 traceability:** `NP-BACKUP, NP-SEC, NP-OBS`. **Permission rule:** audience labels are descriptive; server policy resolves the actual actor/resource/action/scope at runtime.

| Screen ID | Screen / route | Audience | Required content and primary interactions |
|---|---|---|---|
| UX-OPS-001 | **Backup overview**  `/backups` | Authorized | Coverage, last verified restore, failed jobs. **Actions:** Create, configure. |
| UX-OPS-002 | **Backup policies**  `/backups/policies` | Authorized | Schedules, scope, encryption, retention, destinations. **Actions:** Edit, test. |
| UX-OPS-003 | **Backup inventory**  `/backups/sets` | Authorized | Snapshots, size, integrity, expiry. **Actions:** Inspect, restore. |
| UX-OPS-004 | **Restore wizard**  `/backups/restore` | Authorized | Source, target, overwrite impact, dry-run, approval. **Actions:** Review, restore. |
| UX-OPS-005 | **Backup destinations**  `/backups/destinations` | Authorized | Offsite providers, credentials, health. **Actions:** Connect, test. |
| UX-OPS-006 | **Monitoring overview**  `/monitoring` | Authorized | CPU/RAM/disk, service health, latency, trends. **Actions:** Drill down. |
| UX-OPS-007 | **Metrics explorer**  `/monitoring/metrics` | Authorized | Timeseries, comparison, filters, freshness. **Actions:** Chart, export. |
| UX-OPS-008 | **Log explorer**  `/monitoring/logs` | Authorized | Search, time range, severity, tenant scopes. **Actions:** Stream, export. |
| UX-OPS-009 | **Alerts and policies**  `/monitoring/alerts` | Authorized | Thresholds, routing, silences, test delivery. **Actions:** Create, silence. |
| UX-OPS-010 | **Incidents**  `/incidents` | Authorized | Severity, ownership, timeline, status. **Actions:** Acknowledge, assign. |
| UX-OPS-011 | **Incident detail**  `/incidents/:incident` | Authorized | Evidence, runbook, actions, postmortem. **Actions:** Resolve, escalate. |
| UX-OPS-012 | **Security overview**  `/security` | Authorized | Threats, risky settings, MFA, patch status. **Actions:** Review findings. |
| UX-OPS-013 | **Firewall rules**  `/security/firewall` | Authorized admins | Rules, order, lockout preview, IPv4/IPv6. **Actions:** Dry-run, approve. |
| UX-OPS-014 | **SSH security**  `/security/ssh` | Authorized admins | Keys, access, safe change workflow. **Actions:** Add key, test. |
| UX-OPS-015 | **Isolation and policies**  `/security/isolation` | Provider | Tenant boundaries, process users, quotas, exceptions. **Actions:** Inspect, plan. |
| UX-OPS-016 | **Audit trail**  `/audit` | Authorized auditors | Actor/action/target, before-after, export controls. **Actions:** Filter, export. |
| UX-OPS-017 | **Updates and patches**  `/operations/updates` | Authorized admins | Available updates, compatibility, maintenance, rollback. **Actions:** Plan, schedule. |

> **Integrated requirement:** Additional and reframed screens:

| Screen ID | Screen / route | Audience | Required content and primary interactions |
|---|---|---|---|
| UX-OPS-018 | **Allocation and usage**  `/monitoring/usage` | Authorized | **NEW in v2:** per-workload current / typical / peak / throttled view. **Actions:** Drill down, adjust. |
| UX-OPS-019 | **Batch queue**  `/monitoring/batch` | Authorized | **NEW in v2:** all batch jobs (transforms, backups, builds, AI diagnostics, migrations) with queue state. Shows "3 queued, 1 running, will yield to site traffic." **Actions:** Inspect, cancel if supported. |
| UX-OPS-020 | **Memory pressure**  `/monitoring/pressure` | Authorized | **NEW in v2:** workloads near memory.max, OOM kill history, upgrade recommendations. **Actions:** Adjust ceiling, upgrade. |
| UX-OPS-021 | **Backup verification**  `/backups/verify` | Authorized | **NEW in v2:** scheduled restore test evidence. Shows "last verified restore: 3 days ago" for each backup policy. **Actions:** Run verification, view history. |

**Module UX acceptance:** Each screen must have a named owner, authorized route guard, scope breadcrumb, read-only and mutating variants where applicable, specific empty/unsupported/failed states, keyboard test, responsive test, audit correlation and E2E scenario. Resource actions launch impact review and durable jobs; UI cannot mark completion before verification. The module's screens must be cross-linked to `NP-BACKUP, NP-SEC, NP-OBS` requirement IDs individually during implementation planning.

> **Integrated requirement:** On Nano, backups are **remote-only** (UX-OPS-005 is the primary flow; local backup destination is hidden). On Standard+, both local and remote are available; the UI warns if only local backups exist. Backup verification (UX-OPS-021) is a v1 feature, not deferred — "backup succeeded" is not the same as "backup is restorable."

### 8.11 Automation, integrations and developer APIs

**Stage 2 traceability:** `NP-AUTO, NP-EXT`. **Permission rule:** audience labels are descriptive; server policy resolves the actual actor/resource/action/scope at runtime.

| Screen ID | Screen / route | Audience | Required content and primary interactions |
|---|---|---|---|
| UX-AUTO-001 | **Automation overview**  `/automation` | Authorized | Jobs, triggers, rules, failures. **Actions:** Create. |
| UX-AUTO-002 | **Workflow builder**  `/automation/workflows/new` | Authorized | Triggers, steps, approvals, dry-run. **Actions:** Test, publish. |
| UX-AUTO-003 | **Scheduled jobs**  `/automation/schedules` | Authorized | Next/last runs, timezone, ownership. **Actions:** Create, pause. |
| UX-AUTO-004 | **Webhooks**  `/automation/webhooks` | Authorized | Endpoints, signing secrets, delivery retries. **Actions:** Create, test. |
| UX-AUTO-005 | **API explorer**  `/developer/api` | Authorized | Versioned endpoints, examples, scopes. **Actions:** Try safe requests. |
| UX-AUTO-006 | **Integrations**  `/integrations` | Authorized | Connected services, health, permissions. **Actions:** Connect, revoke. |
| UX-AUTO-007 | **Plugin marketplace**  `/plugins` | Authorized | Verified modules, publisher, compatibility. **Actions:** Inspect, install. |
| UX-AUTO-008 | **Plugin detail and consent**  `/plugins/:plugin` | Authorized | Scopes, data access, upgrade policy, warnings. **Actions:** Approve, install. |
| UX-AUTO-009 | **White-label settings**  `/settings/branding` | Reseller, Provider | Brand tokens, assets, domain and preview. **Actions:** Publish. |

**Module UX acceptance:** Each screen must have a named owner, authorized route guard, scope breadcrumb, read-only and mutating variants where applicable, specific empty/unsupported/failed states, keyboard test, responsive test, audit correlation and E2E scenario. Resource actions launch impact review and durable jobs; UI cannot mark completion before verification. The module's screens must be cross-linked to `NP-AUTO, NP-EXT` requirement IDs individually during implementation planning.

> **Integrated requirement:** Plugin marketplace and white-label branding are **deferred to v2+**. The v1 extension mechanism is the `ServiceTemplate` primitive. The `ServiceTemplate` catalog is the v1 "marketplace."

### 8.12 AI assistant and diagnostics

**Stage 2 traceability:** `NP-AI`. **Permission rule:** audience labels are descriptive; server policy resolves the actual actor/resource/action/scope at runtime.

| Screen ID | Screen / route | Audience | Required content and primary interactions |
|---|---|---|---|
| UX-AI-001 | **AI assistant workspace**  `/assistant` | Authorized opt-in | Chat, scoped context, evidence citations, history. **Actions:** Ask, attach resource. |
| UX-AI-002 | **AI diagnostic report**  `/assistant/diagnostics/:report` | Authorized | Read-only evidence, uncertainty, likely causes. **Actions:** Review, export. |
| UX-AI-003 | **AI proposed change**  `/assistant/plans/:plan` | Authorized | Exact actions, risks, impact, policy checks. **Actions:** Edit, request approval. |
| UX-AI-004 | **AI execution timeline**  `/assistant/runs/:run` | Authorized | Approval, job stages, observed results, rollback. **Actions:** Monitor, stop if safe. |
| UX-AI-005 | **AI settings and consent**  `/settings/ai` | Owner, Provider | Provider selection, retention, scopes, opt-in. **Actions:** Enable, revoke. |

> **Integrated requirement:** v1 ships read-only diagnostics + proposed-action approval (NP-AI-012). Autonomous action execution is v4+:

| Screen ID | Screen / route | Audience | Required content and primary interactions |
|---|---|---|---|
| UX-AI-006 | **AI diagnostic first-run**  `/assistant` (state) | Authorized | **NEW in v2:** the v1 default. AI reads logs, metrics, state; proposes typed actions with approval; user approves; policy engine validates; job executes. Clear language: "AI can diagnose and propose; you approve; the policy engine decides." **Actions:** Ask, approve, reject. |
| UX-AI-007 | **AI unavailable**  `/assistant` (state) | Authorized | **NEW in v2:** honest state on Nano if no LLM API key configured — "AI diagnostics require an LLM API key. Add one in Settings → AI." **Actions:** Configure, disable AI. |

**Module UX acceptance:** Each screen must have a named owner, authorized route guard, scope breadcrumb, read-only and mutating variants where applicable, specific empty/unsupported/failed states, keyboard test, responsive test, audit correlation and E2E scenario. Resource actions launch impact review and durable jobs; UI cannot mark completion before verification. The module's screens must be cross-linked to `NP-AI` requirement IDs individually during implementation planning.

> **Integrated requirement:** The AI **does not decide** — the policy engine does. This is displayed honestly in the UI: the AI proposes; the user approves; the policy engine validates again at execution; the job runs; verification determines actual result. The AI never claims success before verification (NP-AI-008 preserved).

### 8.13 Support, migration and lifecycle

**Stage 2 traceability:** `NP-SET, NP-INF`. **Permission rule:** audience labels are descriptive; server policy resolves the actual actor/resource/action/scope at runtime.

| Screen ID | Screen / route | Audience | Required content and primary interactions |
|---|---|---|---|
| UX-SUP-001 | **Support center**  `/support` | All | Guides, tickets, diagnostics, safe data sharing. **Actions:** Open ticket. |
| UX-SUP-002 | **Migration hub**  `/migrations` | Authorized | Source platforms, compatibility, import status. **Actions:** Start migration. |
| UX-SUP-003 | **Migration wizard**  `/migrations/new` | Authorized | Inventory, dry-run, conflicts, cutover, rollback. **Actions:** Validate, import. |
| UX-SUP-004 | **Import/export center**  `/settings/data-transfer` | Authorized | Config exports, redaction, access control. **Actions:** Export, import. |
| UX-SUP-005 | **System settings**  `/settings/system` | Owner, Provider | Locale, URL, module policies, branding, email. **Actions:** Review, save. |
| UX-SUP-006 | **Installer status**  `/system/install` | Owner | Preflight, stages, recovery, logs. **Actions:** Resume, repair. |
| UX-SUP-007 | **Upgrade center**  `/system/updates` | Owner, Provider | Versions, compatibility, backup requirement. **Actions:** Plan, upgrade. |
| UX-SUP-008 | **Recovery center**  `/system/recovery` | Owner, Provider | Offline-safe runbooks and rescue modes. **Actions:** View, verify. |

> **Integrated requirement:** Additional screens:

| Screen ID | Screen / route | Audience | Required content and primary interactions |
|---|---|---|---|
| UX-SUP-009 | **VPS upgrade**  `/system/upgrade` | Owner, Provider | **NEW in v2:** the ladder-climbing flow. Shows current capability set, available capabilities after hardware upgrade, migration options, and the "rescan hardware" button. **Actions:** Rescan, apply, view plan. |
| UX-SUP-010 | **ServiceTemplate export/import**  `/settings/templates/transfer` | Authorized | **NEW in v2:** export a ServiceTemplate instance (storage bucket, self-hosted app, database, cache) as an OCI-compatible artifact + manifest. Import on a new VPS. **Actions:** Export, import. |
| UX-SUP-011 | **Migration conflict resolution**  `/migrations/:id/conflicts` | Authorized | **NEW in v2:** show conflicts between source and destination (name clashes, port conflicts, cert issues), with resolution options. **Actions:** Resolve, rollback. |

**Module UX acceptance:** Each screen must have a named owner, authorized route guard, scope breadcrumb, read-only and mutating variants where applicable, specific empty/unsupported/failed states, keyboard test, responsive test, audit correlation and E2E scenario. Resource actions launch impact review and durable jobs; UI cannot mark completion before verification. The module's screens must be cross-linked to `NP-SET, NP-INF` requirement IDs individually during implementation planning.

> **Integrated requirement:** Migration from cPanel/DirectAdmin/Plesk is **deferred to v2+**. The v1 migration path is **VPS-to-VPS via ServiceTemplate export/import** (UX-SUP-009, UX-SUP-010). The "VPS upgrade" flow (UX-SUP-009) is the ladder paying off: rescan hardware after upgrade, see new capabilities, apply.

### 8.14 Platform Services 

> **Integrated clarification — new module]** This is a new screen family for the `ServiceTemplate` primitive. It covers storage, transform, delivery, cache, queue, search.

| Screen ID | Screen / route | Audience | Required content and primary interactions |
|---|---|---|---|
| UX-SVC-001 | **Platform services overview**  `/services` | Authorized | List of provisioned platform services (storage buckets, transform pipelines, delivery endpoints, cache, queue, search). Health, usage, quick actions. **Actions:** Create, manage. |
| UX-SVC-002 | **Create storage bucket wizard**  `/services/storage/new` | Authorized | Bucket name, visibility (public/private), region, quota, access keys. **Actions:** Create, copy endpoint/keys. |
| UX-SVC-003 | **Storage bucket detail**  `/services/storage/:id` | Authorized | S3 endpoint, access keys (shown once), usage (storage/requests/bandwidth), consumed-by (which sites/apps use it), lifecycle. **Actions:** Rotate keys, view usage, export. |
| UX-SVC-004 | **Storage access keys**  `/services/storage/:id/keys` | Authorized | Scoped keys, permissions, last used, revocation. **Actions:** Create, revoke. |
| UX-SVC-005 | **Transform pipeline**  `/services/transform/:id` | Authorized | Resize/convert/crop rules, pre-transform on upload vs on-the-fly, cache policy, CPU allocation. **Actions:** Configure, apply, view batch queue. |
| UX-SVC-006 | **Delivery configuration**  `/services/delivery/:id` | Authorized | Caddy cache config, optional Cloudflare/Bunny integration, bandwidth usage. **Actions:** Enable integration, test. |
| UX-SVC-007 | **Cache instance**  `/services/cache/:id` | Authorized | Redis/Memcached instance, memory limit, persistence, connection details. **Actions:** Manage, flush, export. |
| UX-SVC-008 | **Queue instance**  `/services/queue/:id` | Authorized | RabbitMQ-style queue, max depth, connection details. **Actions:** Manage, purge. |
| UX-SVC-009 | **Search instance**  `/services/search/:id` | Authorized | Meilisearch-style search, index size, connection details. **Actions:** Manage, export. |
| UX-SVC-010 | **Consume service**  `/services/consume/:id` | Authorized | One-click integration: attach storage bucket to a website/app, inject credentials, update config. **Actions:** Select site, inject, preview diff. |

**Module UX acceptance:** Each screen must have a named owner, authorized route guard, scope breadcrumb, read-only and mutating variants where applicable, specific empty/unsupported/failed states, keyboard test, responsive test, audit correlation and E2E scenario. Resource actions launch impact review and durable jobs; UI cannot mark completion before verification. The module's screens must be cross-linked to `NP-SVC` requirement IDs individually during implementation planning.

> **Integrated requirement:** The storage service's **value proposition is integration depth**: UX-SVC-010 lets the user attach the bucket to a site with one click, injecting the endpoint and scoped credentials into the site's config automatically. No `.env` editing, no copying keys between SaaS dashboards.

### 8.15 Self-hosted Templates 

> **Integrated clarification — new module]** This is the catalog of `ServiceTemplate` instances the user can provision. It turns "one VPS panel" into "one-click marketplace of self-hostable software."

| Screen ID | Screen / route | Audience | Required content and primary interactions |
|---|---|---|---|
| UX-TPL-001 | **Template catalog**  `/templates` | Authorized | Grid of available self-hosted apps: icon, name, category, resource profile, deploy button. Categories: CMS, files, dev, analytics, monitoring, database, cache. **Actions:** Browse, search, deploy. |
| UX-TPL-002 | **Template detail**  `/templates/:id` | Authorized | Description, resource profile (RAM floor/typical/ceiling, disk), prerequisites, screenshots, deploy button. **Actions:** Review, deploy. |
| UX-TPL-003 | **Deploy template wizard**  `/templates/:id/deploy` | Authorized | Hostname, admin credentials (generated), resource allocation, backup policy. Preflight: capability check, port conflict, disk quota. **Actions:** Review, deploy. |
| UX-TPL-004 | **Self-hosted service detail**  `/services/template/:id` | Authorized | App URL, admin login, health, resource usage, backups, export. **Actions:** Open app, manage, export, uninstall. |
| UX-TPL-005 | **Template unavailable**  `/templates` (state) | Authorized | **NEW in v2:** honest state on Nano — "Self-hosted templates require Docker, which is not available on this VPS. Upgrade to Standard to enable the catalog." **Actions:** Upgrade, learn more. |

**Module UX acceptance:** Each screen must have a named owner, authorized route guard, scope breadcrumb, read-only and mutating variants where applicable, specific empty/unsupported/failed states, keyboard test, responsive test, audit correlation and E2E scenario. Resource actions launch impact review and durable jobs; UI cannot mark completion before verification. The module's screens must be cross-linked to `NP-SVC-005` requirement IDs individually during implementation planning.

> **Integrated requirement:** The template catalog is **data, not code**. Adding a template is a database insert (from the admin side) with a signed container image reference. The v1 catalog ships with 5–10 templates (WordPress, Nextcloud, Gitea, Ghost, n8n, Uptime Kuma); v2+ expands to 20–30.

**Screen catalog total:** approximately **165 proposed screen definitions** across **15 module groups**. (The v1.0 count was 138 across 13 groups; the additions are the Router, Service Portfolio, Capability Overview, Allocation/Usage views, Platform Services, Self-hosted Templates, and several honest-state screens for unavailable capabilities.) A screen is a logical route/view, not necessarily a separate frontend bundle or page file.

## 9. Detailed end-to-end journeys (wireflow contracts)

> **Current workflow rule:** FLOW-013 is the primary first-run router. FLOW-001 is its website branch; FLOW-002 through FLOW-012 are contextual journeys accessible from the router, service portfolio or advanced navigation. **FLOW-014 (the Nano experience)** and **FLOW-015 (Platform Services)** are new. The remaining flows are preserved with v2 refinements.

### FLOW-001 — First website deployment (sub-flow of FLOW-013)

> **Integrated requirement:** FLOW-013 is the primary first-run journey; FLOW-001 is the website sub-flow. The reasoning: not all VPS buyers want a website. Some want an app, a self-hosted service, storage, or infrastructure. The Router asks "what do you want to run?" and routes accordingly. The website path is **one of six** paths from the Router, not the default.

The website path, reached from the Router, is preserved as a **sub-flow** of FLOW-013:

**Entry:** owner bootstrap complete, server enrolled, capability scan available. **Steps:** (1) pick "A website" from the router → (2) enter domain or choose temporary preview hostname → (3) choose static/PHP/app type (constrained by capability set — Nano only offers static) → (4) choose source (upload/Git/starter) → (5) select eligible server/automatic placement → (6) inspect document root, disk, account, DNS/TLS prerequisites → (7) impact preview → (8) submit durable provisioning job → (9) timeline with actual checks → (10) verified preview and DNS delegation instructions → (11) optional SSL issuance once DNS validates. **Failure:** no eligible server, domain conflict, quota, DNS not delegated, install error, partial provisioning. **Exit:** accessible site or clear blocked state and recovery link. **Maps:** J1, NP-STR-001, NP-WEB, NP-CORE.

### FLOW-002 — Node.js application: Vercel-like ease with hosting-depth controls

**Entry:** app creation permission, compatible Node adapter. **Wizard steps:** source Git/upload/existing files → repository/branch or file selection → auto-detect package.json and suggested framework → runtime version selected from actual installed/supported versions → build command and working directory → start command/startup file and app port → environment variables via add rows, multiline bulk paste, or `.env` upload → domains/preview URL → resources, process count, restart policy → review secret names only, build/start commands, proxy/TLS and quota → deploy job. **Critical validation:** reject duplicate keys, malformed multiline input, invalid port, path escape, missing start script, unsupported Node version, secret accidentally pasted into non-secret fields; redact values from job logs. **Lifecycle:** queued/building/releasing/starting/health-checking/verified, with log stream and rollback where supported. **Maps:** J2, NP-APP, NP-FILE, NP-TLS.

> **Integrated requirement:** Additions:

- **The build is a batch workload** — it appears in the batch queue indicator, runs with `cpu.weight=10` and `cpu.max=50% of cores`, so it does not starve running sites. The UI shows "Build queued (batch)" not "Build failed (timeout)."
- **The runtime is a serving workload** — it appears in the allocation/usage strip with current / typical / peak.
- **On Nano, this flow is not offered** — the router does not show the "Application" card, or shows it with an honest unavailable state. On Starter, it is borderline — the installer confirms fit before the flow starts.

### FLOW-003 — Provider creates reseller, package and customer

Provider selects organization → creates scoped reseller grant → defines package with quotas, feature entitlements and pricing → previews inherited limits → invites reseller → reseller creates hosting account/customer within delegated capacity → system reserves capacity → provisions resources → customer receives onboarding invitation → provider and reseller see only authorized subtree and billing data. **Failure:** quota exhausted, package changed during creation, invalid delegated permission, server placement failure, customer email already registered. **Maps:** J3, NP-IAM, NP-ACC, NP-BILL.

> **Integrated requirement:** This flow is **deferred to v2+**. The v1 product does not ship the reseller or provider surface. The data model is complete from v1, but the UI is not. When v2 ships, this flow works without architecture change.

### FLOW-004 — DNS and domain onboarding (nine-screen guided sequence)

1. **Intent:** domain name and intended use (site, app, mail, DNS only).
2. **Control model:** external/manual DNS, future external API, or local authoritative DNS; never assume registrar ownership.
3. **Ownership proof:** TXT token and validation status; do not confuse ownership with delegation.
4. **Preflight:** collisions, existing MX/TXT, nameserver availability, TTL, DNSSEC, email dependencies.
5. **Zone proposal:** preview records, detect conflicts, preserve existing email records unless explicitly changed.
6. **Nameserver delegation:** registrar-specific instructions, glue only when required, independent secondary warning.
7. **Security:** DNSSEC choices and DS instructions with risk explanation.
8. **Apply:** diff, approval, operation ID, authoritative server health.
9. **Verify and recover:** public resolution checks, propagation uncertainty, rollback guidance and record-level diagnostics.

**Guardrails:** do not promise global propagation completion; local authoritative DNS needs correct port reachability and independent secondary for production resilience. **Maps:** NP-DNS, NP-TLS, J1/J4.

> **Integrated requirement:** This flow is reframed as **DNS guidance mode** (UX-DNS-014) for v1. The nine-screen sequence is preserved as a future flow when authoritative DNS ships in v2. The v1 flow is simpler: "enter domain" → "add this A record at your registrar" → "we'll poll until it propagates" → "SSL issued, site live." It is a first-class flow, not a settings page.

### FLOW-005 — Email domain activation

Choose domain → check mail role availability and IP reputation/PTR prerequisites → propose MX/SPF/DKIM/DMARC changes with diff → create mailbox policy → request DNS changes → verify external observations → activate mailboxes/webmail → send/receive diagnostics. **Failures:** blocked outbound port 25, missing PTR, stale DNS, DKIM mismatch, unavailable mail node, partial configuration. Never display "deliverability guaranteed." **Maps:** J4, NP-MAIL, NP-DNS.

> **Integrated requirement:** This flow is **deferred to v2+**. The v1 flow is external SMTP (UX-MAIL-011): configure Resend/Postmark/Mailgun/SES credentials, set From address, done. Mailboxes are recommended via Zoho/Google Workspace (UX-MAIL-012).

### FLOW-006 — Database creation, grant and safe import

Select account and engine → compatible node → database name → least-privilege user/grants → quota/backup policy → review → provision job → masked connection string with controlled copy → optional import with schema/size checks and progress. **Failures:** duplicate name, engine missing, capacity, grant conflict, oversized import, partial restore. **Maps:** NP-DB.

> **Integrated requirement:** This flow is **deferred to v2+**. On Nano, databases are unavailable. On Starter, borderline. On Standard+, ships in v2 as a `ServiceTemplate` instance.

### FLOW-007 — File manager and code editor

Open account-scoped root → browse tree with breadcrumb and virtualized rows → select file → preview metadata and contents → edit in tabs with dirty markers → save with version check → show merge/diff on concurrent modification → observe write job and status. Upload supports drag/drop, progress, duplicate conflict choice, limits, resumability when supported. **Security:** no traversal outside permitted root, no unmasked secrets in previews/telemetry, symlink escape detection, dangerous chmod/chown preview. **Maps:** NP-FILE.

> **Integrated requirement:** Preserved. Additional note: on Standard+, the file manager operates on the site's container volume (UX-FILE-001); on Nano, on the site's document root directly.

### FLOW-008 — Backup restore and disaster recovery

Select restore source → validate integrity, encryption keys and target compatibility → choose original/new destination → enumerate overwritten resources and downtime → offer verified backup checkpoint when possible → step-up MFA/approval → restore job with per-step status → integrity and application health checks → incident record and rollback options. **Failures:** offsite unavailable, corrupt backup, incompatible engine, missing key, insufficient disk, partial restore. Never equate completed file copy with verified service recovery. **Maps:** J5, NP-BACKUP.

> **Integrated requirement:** Preserved. Additional note: backup verification is scheduled (UX-OPS-021), not ad-hoc. The restore wizard shows "last verified restore: 3 days ago" as a trust signal.

### FLOW-009 — Firewall or SSH change without lockout

Inspect existing connectivity and session → stage new rules/keys → validate IPv4 and IPv6 and explicit management access → show diff and lockout warning → require alternate verified access path where feasible → apply with timed automatic rollback mechanism if supported → test connectivity from independent session → commit only after verification. If timed rollback unavailable, label operation higher risk and require out-of-band recovery. **Maps:** NP-SEC.

> **Integrated requirement:** Preserved.

### FLOW-010 — AI-guided incident remediation

Open incident → AI reads only permitted, redacted evidence → report observations vs hypotheses with uncertainty → propose typed operations with affected resources, expected effect, prerequisites and rollback → user edits/approves → policy engine evaluates → job executes → AI displays actual observed results and unresolved symptoms. No arbitrary root-shell execution; no automatic irreversible actions. **Maps:** J6, NP-AI, NP-CORE.

> **Integrated requirement:** Preserved. Additional note: this is **read-only diagnosis + proposed-action approval** in v1 (NP-AI-012). Autonomous multi-step remediation is v4+. The AI never claims success before verification.

### FLOW-011 — Multi-server role deployment and maintenance

Provider opens fleet → enrolls agent using short-lived token → verifies server identity → inspects capability/resource profile → assigns role and placement eligibility → previews conflicts/impact → provisions module → verifies health → schedules maintenance/cordon/drain → tracks workload movement and failures. Never call one-node deployment highly available. **Maps:** NP-INF.

> **Integrated requirement:** Deferred to v2+ with the Provider profile.

### FLOW-012 — Reseller-branded customer onboarding

Reseller configures branding preview → confirms domain/TLS → creates package → invites customer → customer receives branded signup with clear operator identity → sees only assigned hosting account and enabled features → creates site/mail/database within quota → support routes to reseller/provider per policy. **Maps:** NP-ACC, NP-IAM, NP-EXT.

> **Integrated requirement:** Deferred to v2+.

### FLOW-013 — The Router: "What do you want to run?" (primary first-run)

**Entry:** owner bootstrap complete (UX-IAM-002), server enrolled, capability scan available.

**The screen:** a single page with six cards. Each card has an icon, a title, a one-line description, and a "Get started" button. Above the cards: "Welcome to NexusPanel. What do you want to run?" Below the cards: "I'm not sure yet →".

**The six cards:**

```
🌐  A website              WordPress, static, or your own
⚙️  An application         Node, Python, Go, Ruby, Docker
📦  A self-hosted app      Nextcloud, Gitea, Ghost, n8n…
🗄️  Data & storage         Databases, object storage, images
🔒  Infrastructure         VPN, proxy, DNS, mail
🤔  I'm not sure yet       Show me what's possible
```

**The routing:**

| Card | Route | "You're live" moment | Failure mode |
|---|---|---|---|
| A website | FLOW-001 website sub-flow | Site reachable at `name.panel-host` with SSL | DNS not propagated → guidance + retry; capability unavailable → honest display |
| An application | FLOW-002 (if capability available) | App responds at `name.panel-host`; logs visible | Capability unavailable (Nano) → honest display + upgrade |
| A self-hosted app | FLOW-015 (template catalog → deploy) | App's own login screen loads; admin credentials shown | Capability unavailable (Nano) → honest display + upgrade |
| Data & storage | FLOW-015 (storage) or FLOW-006 (database, if v2+) | Bucket exists + S3 endpoint + keys shown; first upload succeeds | Capability unavailable (Nano) → honest display + upgrade |
| Infrastructure | Depends on sub-choice | Service running + connection details | Capability unavailable → honest display |
| I'm not sure yet | Guided tour | Suggestion + "try this first" | Fallback to capability display |

**After the first path:** the user lands on the **service portfolio** (UX-OV-001), which lists everything they now run. The checklist follows: "☑ Your first site is live. ☐ Connect your domain. ☐ Enable automatic backups. ☐ Invite a collaborator." Each item links to the relevant flow.

**Capability awareness:** the router reads the installer's capability set and shows honest availability:
- On Nano: "An application" is shown with a small note: "Not available on this VPS — needs 200 MB more RAM." Clicking shows the upgrade path.
- On Starter: all cards are shown; borderline ones ("An application," "Data & storage → Databases") show "runs, but tight" and the resource delta.
- On Standard+: all cards are fully available.

**Maps:** NP-STR-010, NP-UX-011, NP-UX-012, NP-CORE-020, NP-STR-001 (lay-person gate).

### FLOW-014 — The Nano layman experience 

**Entry:** a 512 MB VPS install with NexusPanel. Installer ran 5-step detection and produced a capability set.

**The first-run router** shows six cards with honest availability:

```
Welcome to NexusPanel. Your VPS can run:

🌐  A website              ✓ Static sites work
                            (Portfolio, landing page, docs)
⚙️  An application         ✗ Needs 200 MB more RAM — upgrade
📦  A self-hosted app      ✗ Needs Docker — upgrade
🗄️  Data & storage        ✓ File storage client
                            (Connect to external S3/B2)
                            ✗ Databases need 250 MB more
🔒  Infrastructure         ✗ Email needs port 25 + reverse DNS
🤔  I'm not sure yet       Show me what's possible

[ Upgrade my VPS → ]
```

**The user picks "A website"** → static site flow → site live in 60 seconds with SSL.

**The service portfolio** lists one site, live, with SSL. The capability banner persists at the top: "Your VPS can also run: storage client, remote backups, AI diagnostics. Your VPS cannot run: WordPress, Node, databases, transforms, email. [Learn more]"

**If the user tries an unavailable capability** — e.g., clicks "An application" — they see: "Your VPS cannot run Node.js applications because the panel and OS already use most of your 512 MB RAM. Node needs 200 MB more. [Upgrade to a 2 GB VPS →] [Learn what a 2 GB VPS can run →]"

**No greyed-out menus. No fake features. Just the honest picture.** This is the ladder working: the product scales down honestly, not by hiding features.

**Maps:** NP-STR-008, NP-INF-013 through NP-INF-016, NP-NFR-019, J8.

### FLOW-015 — Platform Services 

**Storage bucket flow:**

1. User picks "Data & storage" from the router, or navigates to `/services` from the menu.
2. User picks "Storage → Create bucket."
3. Wizard: bucket name, visibility (public/private), region, quota, access keys.
4. Preflight: disk quota fit, MinIO health, port conflict.
5. Provision job: create bucket, generate scoped keys, set public/private policy.
6. **"You're live" moment:** bucket exists, S3 endpoint shown, access keys shown once with "copy" buttons, first upload succeeds.
7. Usage dashboard: storage used, requests, bandwidth, top assets.
8. **Integration offer:** "Use this bucket for a website's media?" → pick site → inject endpoint and keys into site's config → preview diff → confirm.

**Transform pipeline flow:**

1. From the storage bucket detail: "Enable image transformation."
2. Choose rules: pre-transform predictable sizes (thumbnail, medium, large) on upload vs. on-the-fly transform via URL params.
3. Configure CPU allocation: batch workload with `cpu.weight=10`, `cpu.max=50% of cores`.
4. Confirm: "Transforms will run when the server is idle. Under site traffic, they queue. Your sites stay responsive."
5. **"You're live" moment:** transform pipeline active; upload an image → auto-converted; serve `?w=800` → resized.

**Self-hosted template flow:**

1. From `/templates` or "A self-hosted app" from the router.
2. Browse catalog: WordPress, Nextcloud, Gitea, Ghost, n8n, Uptime Kuma.
3. Pick Nextcloud → deploy wizard → hostname, admin credentials (generated), resource allocation, backup policy.
4. Preflight: Docker capability, port conflict, disk quota, memory fit.
5. Provision job: pull signed image, create container, volume, reverse proxy route, SSL, admin credentials, health check.
6. **"You're live" moment:** Nextcloud's own login screen loads at `cloud.name.panel-host` with SSL; admin credentials shown once.

**Maps:** NP-SVC-001 through NP-SVC-006, NP-TPL-001 through NP-TPL-005, FLOW-001 sub-flow for storage and templates.

## 10. Field-level UX specifications for high-complexity forms

| Form | Field | Input and validation | Error/help behavior |
|---|---|---|---|
| New website | Primary hostname | Normalize IDN, reject invalid labels, check ownership/conflicts | Explain DNS ownership vs routing separately |
| New website | Document root | Account-scoped path picker; no traversal/symlink escape | Show safe root and why rejected |
| New application | Source | Git/upload/existing/starter; mutually exclusive initial source | Show adapter availability |
| New application | Node version | Dropdown from verified server capabilities, no invented options | Show compatibility and install prerequisite |
| New application | Build/start | Detect scripts, allow explicit override with advanced warning | Missing script or invalid working dir |
| New application | Port | Integer within approved range; no reserved conflicts | Explain proxy and internal port distinction |
| Environment | Bulk paste | Parse KEY=VALUE, quotes/comments/newlines; duplicate detection | Preview parsed names and masked values; no accidental analytics |
| Environment | `.env` upload | Local parsing, file size cap, no unapproved secret upload retention | Per-line errors and safe discard |
| DNS record | Type/name/value/TTL | Type-aware validation, duplicate and conflict detection | Show record examples, no hidden auto-rewrite |
| DNSSEC | DS/key algorithm | Registrar compatibility, explicit sequence | Warn of resolution outage from mismatched DS |
| Mailbox | Address/quota/password | Domain ownership, mailbox uniqueness, policy | Quota, password policy, reset workflow |
| Database | Engine/name/user/grants | Server capability, identifier rules, least privilege | Explain host vs DB user and storage quota |
| Backup | Scope/schedule/retention/target | Capacity and offsite readiness; UTC/timezone clarity | Show effective retention and cost implications |
| Restore | Target/overwrite | Resource-specific inventory and backup health | Typed target confirmation for irreversible overwrite |
| Package | Resource limits | Units, reserved vs active, upstream ceiling | Prevent overselling unless approved policy |
| Firewall | Protocol/ports/CIDR/direction | CIDR and IP family validation, shadowing/conflicts | Preview effective rule order and lockout risk |
| Reseller | Scope/limits/brand | Cannot exceed provider delegation; domain validation | Explain inherited vs override settings |
| API token | Name/scopes/expiry | Least privilege, explicit lifetime, one-time secret reveal | Copy once; never display again |

**Form behavior:** Labels never rely on placeholders; required/optional explicit; preserve user input on recoverable errors; disable submission while policy check pending; allow safe draft save; reveal sensitive values only after explicit permission and never in screenshots/telemetry. For multiline bulk paste, show parsed counts, duplicates and warnings before applying.

> **Integrated requirement:** Additional field specifications for the new modules:

| Form | Field | Input and validation | Error/help behavior |
|---|---|---|---|
| Storage bucket | Bucket name | S3-compatible rules (lowercase, hyphens, 3–63 chars), uniqueness per region | Show S3 naming rules; reject silently otherwise |
| Storage bucket | Visibility | Public/private, with clear "public = anyone can read" warning | Preview public URL before confirming |
| Storage access key | Permissions | Scoped (read-only, write-only, read-write), resource scope | Show least-privilege default; explain escalation |
| Transform pipeline | Pre-transform sizes | List of named sizes (thumbnail, medium, large, custom) | Show storage cost estimate per size |
| Transform pipeline | On-the-fly | Toggle with cache warning ("unpredictable sizes will be cached for X days") | Warn about cache invalidation |
| Template deploy | Admin credentials | Generated password shown once, or user-provided | Show "copy password" once; never display again |
| Template deploy | Resource allocation | RAM/CPU/disk within capability set | Refuse if exceeds fit; show delta |
| M-Pesa settings | Daraja credentials | Consumer key/secret, passkey, shortcode, sandbox/production selector | Test transaction required before enabling |
| M-Pesa settings | Callback URL | Read-only, panel-generated, copyable | Show "copy this URL to Daraja dashboard" |
| Payment pending | Phone number | E.164 format for M-Pesa; validate before STK push | Show "check your phone" after push |

## 11. File manager and code editor detailed design

**Desktop layout:** left resource/account selector and file tree (resizable), center tabbed editor/preview, optional right inspector (metadata, Git diff, permissions), bottom terminal/log drawer if allowed. **Mobile:** folder list first, single-file view/edit, limited bulk actions, no assumption of desktop keyboard. **Tree:** lazy loading, virtualized large folders, sort, search, breadcrumbs, drag-drop with confirmation, selection persistence only within authorized root. **Editor:** language detection, line numbers, syntax, find/replace, undo/redo, soft-wrap, encoding/line endings, file size limits, save state, conflict diff, optional formatting and syntax validation. **File operations:** copy/move/rename/delete/archive/extract/chmod/chown with impact review and safe extraction rules. **Permissions:** show effective OS/account permission, explain why editing may fail; never promise arbitrary root access. **Terminal:** optional, permission-gated, recorded according to policy, clear scope, idle timeout and visible exit.

> **Integrated requirement:** Additional design notes:

- **File manager is container-aware** on Standard+ — the tree is scoped to the site's container volume, not the host filesystem. The path shown is the container path (e.g., `/app/wp-content/uploads`), not the host path.
- **Terminal is disabled by default** for customer roles. On Standard+, a scoped container shell is available when explicitly enabled. On Nano, terminal is not available at all (UX-FILE-010).
- **SFTP is the primary file-transfer path on Nano** (UX-FILE-011). The panel manages the SFTP endpoint and per-site credentials.

## 12. Role-specific dashboards and KPIs

| Dashboard | Primary cards | Critical attention area | Primary CTA |
|---|---|---|---|
| Solo owner | Sites/apps, CPU/RAM/disk, domains, SSL, backups | Failed deployments, backup gaps, expiring SSL | Create website/app |
| Customer | Assigned sites, email, databases, disk/bandwidth quota | Account suspension, failed jobs, SSL | Manage website |
| Developer | Deployments, build health, logs, app resource usage | Failed deploy, crash loops, env changes | Deploy app |
| Reseller | Customers, packages, account states, allocated quota | Over-quota, unpaid/suspended accounts, support | Create account |
| Provider | Fleet health, capacity, tenants, incidents, revenue indicators | Offline nodes, abuse/security, backups, billing drift | Add server / investigate |
| Support | Assigned tickets, access grants, incident context | Expiring grants, high severity incidents | Open ticket |
| Auditor | Change history, approvals, access anomalies | Unreviewed privileged changes | Review audit |

Dashboard metrics must show source, measurement period and last refresh. If a metric is unavailable, display "Unavailable" rather than zero. Avoid exposing financial/tenant details to unauthorized roles.

> **Integrated requirement:** The role-specific dashboards are reframed as **portfolio views** for the service portfolio (UX-OV-001). The "Solo owner" dashboard becomes the **service portfolio** — a list of what the user runs, not a KPI dashboard. The **capability overview** (UX-OV-002) is a separate surface for "what can my VPS run?" The **allocation/usage** view (UX-OV-009) answers "am I close to my limits?"

## 13. Responsive, accessibility and localization contract

**Breakpoints (proposed):** small mobile <640px, large mobile/tablet 640–1023px, desktop ≥1024px, wide ≥1440px. Responsive rules must be content-based, not solely device labels. Tables collapse to cards or horizontal scroll with sticky identity column; destructive actions never become hidden behind inaccessible hover. Dialogs become full-height sheets on narrow screens; multi-step wizards retain step title, progress and back controls. Code editor supports reduced capability on mobile with explicit limitations. Touch targets target ≥24×24 CSS px minimum per WCAG 2.2 AA, preferably 40–44px for common actions; keyboard focus visible and never obscured. Support 200% zoom, text resizing, high contrast, reduced motion, screen readers, no color-only status, accessible error summary and focus restoration after modal close. Keyboard shortcuts must not override browser/screen reader shortcuts without user control. Prepare copy for localization and pluralization, IDN/punycode display, 24-hour vs 12-hour time, units and user time zones. Test RTL before claiming support.

> **Integrated requirement:** Additional responsive notes:

- **The router cards** stack vertically on mobile, display as a 2×3 grid on tablet, 3×2 on desktop. Each card remains fully tappable.
- **The capability banner** is collapsible on mobile; the primary text ("Your VPS cannot run WordPress") remains visible without expanding.
- **The allocation/usage strip** becomes a horizontal scroll on mobile, not a wrap; each workload remains reachable.
- **The batch queue indicator** remains in the shell on all breakpoints; it is important for understanding why a build/backup is slow.

## 14. Notification, error and confirmation copy examples

| Situation | Recommended copy | Action |
|---|---|---|
| Job submitted | "Website creation started. We'll verify the result before marking it ready." | View progress |
| DNS not delegated | "Your domain still points to different nameservers. Update the registrar settings, then recheck." | Show nameservers |
| SSL pending | "Certificate issuance is waiting for domain validation." | Check DNS |
| Deployment failed | "The application build failed. No new release was activated." | View build logs |
| Quota exceeded | "This account has 1.2 GB available; the requested backup needs 2 GB." | Review usage |
| Server offline | "Last contact with this server was 12 minutes ago. Current health is unknown." | Diagnostics |
| Partial restore | "Files restored, but database verification failed. The service is not marked recovered." | Recovery plan |
| Firewall risk | "This rule could block SSH access. Verify another access method before applying." | Review changes |
| Unauthorized | "You don't have permission to perform this action." | Request access |
| Secret copied | "Secret copied. It will not be shown again." | Close |
| AI recommendation | "Suggested action based on available logs; review the impact before approval." | Review plan |

Never show raw stack traces, SQL, passwords, tokens, private keys or unredacted environment variables to general users. Provide correlation IDs for support.

> **Integrated requirement:** Additional copy for the new states:

| Situation | Recommended copy | Action |
|---|---|---|
| Capability unavailable | "Your VPS cannot run WordPress because it needs 150 MB more RAM. Upgrade to a 2 GB VPS to enable PHP hosting and WordPress." | Upgrade, learn more |
| Batch queued | "Your backup is queued. It will run when your websites are idle, usually within minutes." | View queue |
| Batch running (throttled) | "Your image transforms are running at reduced speed to keep your sites responsive. This is normal." | View metrics |
| Memory pressure | "Your Node app is using 380 MB of its 400 MB limit. If it hits the limit, it will be restarted. Consider upgrading or reducing usage." | View metrics, adjust |
| Pending M-Pesa | "We've sent a payment request to your phone. Enter your M-Pesa PIN to complete the payment. Waiting for confirmation…" | Wait, retry, cancel |
| Payment confirmed | "Payment received. Your subscription is active." | View invoice |
| Payment callback unmatched | "We received a payment notification we couldn't match to an invoice. An administrator will review it within 24 hours." | Contact support |
| Template deploying | "Deploying Nextcloud. Pulling image, creating container, configuring SSL. Usually takes 2–3 minutes." | View progress |
| Template unhealthy | "Nextcloud's health check is failing. The container may need a restart, or you may need to restore from a backup." | Restart, restore, uninstall |
| Nano capability display | "Your 512 MB VPS can run: static sites, storage client, remote backups, AI diagnostics. It cannot run: WordPress, Node, databases, transforms, email." | Upgrade, learn more |
| Upgrade recommendation | "Your peak RAM has been above 80% for 7 days. Consider upgrading to a 4 GB VPS for more headroom." | View evidence, upgrade |
| Router I'm not sure | "Let's see what your VPS can do first." | Show capability display, guided tour |

## 15. UX instrumentation, privacy and performance

- Track page load and task completion at coarse granularity: `screen_viewed`, `wizard_started`, `wizard_step_completed`, `preflight_failed`, `operation_submitted`, `operation_verified`, `operation_failed`, `help_opened`, `approval_decided`.
- Include event version, screen ID, anonymized actor scope, capability, operation ID and duration where allowed. Never collect raw secrets, full paths containing sensitive names, email contents, SQL queries, terminal input or arbitrary logs by default.
- Measure perceived performance via navigation responsiveness, LCP/INP/CLS where applicable; set budgets only after baseline profiling. Use progressive rendering for large tables/logs, virtualized file trees and cancellable searches.
- Use safe polling/backoff or subscriptions for live jobs; label disconnect and last update; avoid infinite client retries and duplicate operations.
- All telemetry and error-reporting providers require documented data processing/retention review.

> **Integrated requirement:** Additional instrumentation:

- **Router conversion:** which of the six cards does the user pick? Where do they drop off? This is the primary funnel metric for the lay person.
- **Capability display engagement:** does the user click "Upgrade" from the capability banner? Does the banner reduce support tickets about "why can't I install X?"
- **Batch queue visibility:** does the user understand that backups/transforms are queued (not failed)? Measure support tickets about "backup is slow."
- **Payment pending engagement:** does the user complete M-Pesa payments? Where does the flow drop off?
- **Template catalog adoption:** which templates are deployed most? Which are abandoned mid-deployment?

## 16. Design QA and acceptance matrix

| Test ID | Required scenario | Pass criteria |
|---|---|---|
| UXT-001 | First website as novice | Completes without shell after install; status verified |
| UXT-002 | Node app with bulk env paste | Correct parsed keys, masked secrets, no telemetry leak |
| UXT-003 | Customer vs reseller scope | No cross-tenant search, route or cached data exposure |
| UXT-004 | Failed job + browser refresh | Persistent accurate state, safe retry, audit link |
| UXT-005 | DNS not delegated | Honest pending state, precise registrar guidance |
| UXT-006 | Restore overwrite | Exact target impact, MFA/approval, verification |
| UXT-007 | Firewall lockout attempt | Warning, preflight and recovery safeguards |
| UXT-008 | Agent offline | Last seen shown; no fake health or success |
| UXT-009 | Keyboard-only wizard | All fields, errors, steps, confirmations operable |
| UXT-010 | Screen reader dialogs | Correct labels, focus trap/restoration, announcements |
| UXT-011 | 320px and 200% zoom | No loss of critical content/actions |
| UXT-012 | Reduced motion / themes | Motion respects preference; contrast validated |
| UXT-013 | Concurrent editor changes | Version conflict diff; no silent overwrite |
| UXT-014 | Reseller quota edge | Limits and reservation explained; backend enforced |
| UXT-015 | AI proposed privileged action | No execution before explicit approval and policy check |
| UXT-016 | Mobile file upload | Progress, conflict handling, retry, accessible status |
| UXT-017 | Email domain prerequisite failure | Accurate PTR/port/DNS warnings, no false deliverability |
| UXT-018 | Partial multi-server operation | Per-node results, retry/compensation, audit |

> **Integrated requirement:** Additional design QA scenarios:

| Test ID | Required scenario | Pass criteria |
|---|---|---|
| UXT-019 | Router first-run | Non-technical user completes a path in under 10 minutes without reading docs |
| UXT-020 | Router capability honesty | User can articulate what their VPS can/cannot run from the UI alone |
| UXT-021 | Nano layman | Static site live in 60 seconds; capability banner is clear; no dead ends |
| UXT-022 | Batch queue visibility | User understands that backups/transforms are queued, not failed |
| UXT-023 | Allocation vs usage | User can answer "am I close to my limits?" from the dashboard |
| UXT-024 | Pending M-Pesa | Async payment state is clear; user can wait, retry, or cancel |
| UXT-025 | Capability display upgrade | Upgrade link from capability banner routes correctly |
| UXT-026 | Storage bucket integration | One-click "use this bucket for site" injects credentials and works |
| UXT-027 | Template deploy health | Template deploys, health-checks, and is reachable at `name.panel-host` with SSL |
| UXT-028 | ServiceTemplate export | Export a Nextcloud instance, import on another VPS, works end-to-end |
| UXT-029 | Guided → Expert mode | User switches modes; permissions identical; presentation clear |
| UXT-030 | cgroup metrics display | Current/typical/peak visible per workload; matches measured values |

**Release gate:** every shipped screen must pass route authorization, functional tests, accessibility, responsive review, privacy/secret review, error-state review and API contract tests. A visual mockup alone is insufficient.

> **Integrated requirement:** Additional release gate for v2:

- **The guided-default gate is a release gate, not a later feature.** If a non-technical user cannot go from install to a live 1-page site in under 10 minutes without touching a menu, the release does not ship. This is the one thing to hold the line on. If it slips, the product has drifted back toward "cPanel with better docs."

## 17. Prototype and handoff deliverables

1. **Design tokens:** semantic colors, typography, spacing, radii, elevation, icons, dark/light, density, accessibility notes.
2. **Reusable component library:** all DS-001–020 components with story states and interaction tests.
3. **Navigation prototypes:** all personas, context switching, mobile drawer and keyboard search.
4. **High-fidelity prototypes:** owner onboarding, website creation, Node deployment/env import, DNS nine-step onboarding, email setup, file manager/editor, backup restore, provider/reseller workflow, fleet, AI approval.
5. **Screen-level specs:** screen ID, annotated wireframe, data contract, actions, errors, role scopes, breakpoints, copy and analytics.
6. **State matrix:** loading/empty/pending/success/partial/failure/offline/permission/quota for each critical screen.
7. **UX test evidence:** moderated novice/expert testing, accessibility report, task outcomes, prioritized fixes.
8. **Engineering handoff:** route map, component inventory, tokens, interaction diagrams, API needs, test IDs and unresolved decisions.

> **Integrated requirement:** Additional deliverables:

9. **Router prototype:** the six-card screen, the routing logic, the capability-aware availability display, and the "I'm not sure" fallback.
10. **Service portfolio prototype:** the list of what the user runs, with health, usage, quick actions.
11. **Capability display prototype:** honest display of what the VPS can/cannot run, with resource deltas and upgrade path.
12. **Allocation/usage prototype:** current / typical / peak / throttled view per workload, and server-level aggregate.
13. **Batch workload prototype:** queue visibility, informational states, "will yield to site traffic" messaging.
14. **Platform Services prototype:** storage bucket creation, transform pipeline, template catalog and deploy flow.
15. **Async payment prototype:** pending M-Pesa state, callback matching, reconciliation for admins.

## 18. Open design decisions and recommended defaults

| ID | Decision | Proposed default | Approval required |
|---|---|---|---|
| UXD-001 | Product name | NexusPanel working name | Yes |
| UXD-002 | Brand system | Neutral surfaces + restrained indigo accent | Yes |
| UXD-003 | Sidebar IA | Role-aware grouped navigation | Yes |
| UXD-004 | Organization/account context | Explicit switcher + persistent breadcrumbs | Yes |
| UXD-005 | Beginner/expert modes | Progressive disclosure, same permissions | Yes |
| UXD-006 | Homepage by persona | Role-specific overview | Yes |
| UXD-007 | First website flow | Guided wizard with preview hostname option | Yes |
| UXD-008 | Environment variables | Rows + bulk paste + `.env` upload | Yes |
| UXD-009 | DNS onboarding | Nine-step guided process | Yes |
| UXD-010 | File editor | Tree + tabs + diff + safe save | Yes |
| UXD-011 | Terminal | Off by default for customer roles | Yes |
| UXD-012 | Dangerous actions | Impact preview, step-up and policy approvals | Yes |
| UXD-013 | AI assistant | Opt-in, evidence-first, no silent mutations | Yes |
| UXD-014 | Mobile | Fully responsive admin with selective expert limitations | Yes |
| UXD-015 | Themes | Light/dark/system and compact density | Yes |
| UXD-016 | White label | Reseller brand overrides within accessibility constraints | Yes |
| UXD-017 | Localization | English first, translation-ready | Yes |
| UXD-018 | Navigation URL scheme | Scoped resource routes with stable IDs | Yes |
| UXD-019 | Data tables | Cursor pagination and saved filters where appropriate | Yes |
| UXD-020 | Support access | Visible scoped banner with expiry and audit | Yes |

> **Integrated requirement:** The v1.0 open decisions are preserved. Additional decisions from the deliberations:

| ID | Decision | Proposed default | Approval required |
|---|---|---|---|
| UXD-021 | **Default landing** | Router for new users; service portfolio for returning users | **Yes** — release gate |
| UXD-022 | **Router card set** | Six cards: Website / Application / Self-hosted app / Data & storage / Infrastructure / I'm not sure | **Yes** — this is the primary IA |
| UXD-023 | **Capability display style** | Honest availability with resource delta; not greyed-out buttons | **Yes** |
| UXD-024 | **Allocation/usage strip** | Compact strip below page header on every workload page | Yes |
| UXD-025 | **Batch queue visibility** | Queue indicator in shell; informational, not error | Yes |
| UXD-026 | **Service portfolio as primary surface** | Replaces the v1 "Overview" dashboard as the default landing for returning users | **Yes** |
| UXD-027 | **Platform Services UI** | Storage, transform, delivery, cache, queue, search under `/services` | **Yes** |
| UXD-028 | **Self-hosted Templates UI** | Template catalog under `/templates`; deploy flow with preflight | **Yes** |
| UXD-029 | **Async payment UI** | Pending state with clear wait/retry/cancel | **Yes** |
| UXD-030 | **M-Pesa settings** | Daraja credentials, callback URL, test transaction required before enabling | **Yes** |
| UXD-031 | **Nano experience** | Honest capability display; no greyed-out menus; no fake features | **Yes** — release gate |
| UXD-032 | **Upgrade path UI** | "Rescan hardware" and "View upgrade plan" as first-class flows | **Yes** |
| UXD-033 | **Guided → Expert mode** | Toggle; same permissions; presentation only | **Yes** |

## 19. Stage 3 approval checklist and next stage

- [ ] Confirm working product name and design direction.
- [ ] Approve persona-based information architecture and context switching.
- [ ] Review complete screen catalog for missing modules or screens.
- [ ] Approve critical user journeys and high-risk confirmation behavior.
- [ ] Approve file manager, Node.js environment import and DNS onboarding interaction requirements.
- [ ] Confirm accessibility target, responsive requirements and content standards.
- [ ] Assign owners and resolve/defer `UXD-*` decisions.
- [ ] Approve prototype deliverables and handoff criteria.
- [ ] Map each shipped screen to specific `NP-*` requirements, API contracts and tests before coding.

**Next: Stage 4 — Technical System Architecture and Architecture Decision Records (ADRs).** Define control-plane boundaries, agent trust, jobs/reconciliation, database tenancy, service adapters, deployment topology, failure modes, security and install lifecycle. UI/UX specification is a reviewable design contract, not authorization to modify the Leaseweb VPS.

> **Integrated requirement:** Additional approval items:

- [ ] **Approve UXD-021 (default landing = router), UXD-022 (six-card router), UXD-026 (service portfolio as primary surface), UXD-031 (Nano honest experience).** These are the release gates for the lay-person experience.
- [ ] **Approve UXD-023 (capability display style)** — honest, with resource delta, not greyed-out.
- [ ] **Approve UXD-027 (Platform Services UI) and UXD-028 (Self-hosted Templates UI).** These are new modules.
- [ ] **Approve UXD-029 and UXD-030 (async payment UI).** This is where M-Pesa lives.
- [ ] **Approve the additional design QA scenarios (UXT-019 through UXT-030).** These are the tests that prove the lay-person experience works.
- [ ] **Approve the additional prototype deliverables (items 9–15 in §17).** These are what the design team produces next.

---

**Change control:** This document is a comprehensive draft for product and design review, not approved for implementation. Changes to requirements, personas, modules, or journeys require an explicit impact review and cross-reference against Volume 01 Master PRD and the companion engineering volumes.

**End of Volume 02 — Complete UI/UX & User Journeys.**

---


## Volume-level acceptance and handoff gates

| Gate | Pass condition | Evidence required |
|---|---|---|
| UX-GATE-01 Router | New user selects a goal and sees only currently eligible paths, with truthful alternatives | Moderated Nano and Standard usability tests |
| UX-GATE-02 Tenant isolation | No unauthorized resource identifiers, details or actions leak via navigation, search, notifications or deep links | Negative E2E and API authorization tests |
| UX-GATE-03 Operation honesty | All privileged actions display queued/running/verified/failed/recovery states from authoritative operations | Failure injection with delayed and interrupted agents |
| UX-GATE-04 Small-screen | Critical setup, status, security recovery and support tasks function on narrow mobile viewport | 320px and keyboard/screen-reader checks |
| UX-GATE-05 Capability truth | Nano, Starter, Standard and Provider views correctly distinguish not installed, unsupported, unentitled, deferred and unavailable | Capability matrix fixture tests |
| UX-GATE-06 Payment truth | No access is granted merely because a payment request was submitted; pending, confirmed, failed and reversed are distinct | Gateway callback/reconciliation E2E |
| UX-GATE-07 Recovery | Every destructive or high-risk workflow shows impact, approvals, result, audit link and recovery guidance | Tabletop and fault-injection scenarios |
| UX-GATE-08 Performance | UI and static assets fit agreed bandwidth, startup and runtime budgets; do not assume 512 MB feasibility without tests | Repeatable Nano benchmarks |

**Open approval gates:** final brand tokens, persona preference persistence, mobile navigation breakpoints, real Nano performance, precise resource-delta calculations, and end-user language/localization policy. These remain proposed, not silently approved.
