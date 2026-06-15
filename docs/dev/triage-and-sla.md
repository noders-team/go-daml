# Issue Triage & Response SLA

How issues reported against `github.com/noders-team/go-daml` are categorized,
prioritized, and responded to. It sets **target response times** per severity so
reporters know what to expect and maintainers have a consistent process.

> **These are Service Level *Expectations* (SLEs), not contractual SLAs.**
> go-daml is an open-source SDK maintained on a best-effort basis. The targets
> below are goals the maintainers aim for during normal operation — they are not
> a legal guarantee. Teams that need contractual, enforceable SLAs should arrange
> [commercial support](#support-tiers--escalation). This distinction follows the
> SLE-vs-SLA guidance for open-source projects (see [Sources](#sources)).

- [Scope](#scope)
- [Severity levels](#severity-levels)
- [Response SLE matrix](#response-sle-matrix)
- [Resolution targets](#resolution-targets)
- [Issue type & status labels](#issue-type--status-labels)
- [The triage workflow](#the-triage-workflow)
- [Security issues](#security-issues)
- [Stale & inactive issues](#stale--inactive-issues)
- [Support tiers & escalation](#support-tiers--escalation)
- [Definitions](#definitions)
- [Sources](#sources)

## Scope

Applies to issues opened on the GitHub tracker for this repository: bugs,
integration failures, code-generation defects, documentation problems, and
feature requests. Security vulnerabilities follow a separate, private path —
see [Security issues](#security-issues).

"Integration issue" here means a problem connecting this SDK to a DAML / Canton
participant: client/auth failures, service call errors, codegen output that
won't compile or round-trip, and protocol/version mismatches — the failure modes
catalogued in [troubleshooting.md](./troubleshooting.md).

## Severity levels

Severity describes **impact**, independent of how many users are affected. It is
assigned by a maintainer during triage and may change as understanding improves.
The model mirrors the widely used S1–S4 scheme (GitLab, SEV taxonomy — see
[Sources](#sources)), specialized for an integration SDK.

| Severity | Label | Meaning | Integration examples |
| --- | --- | --- | --- |
| **S1 — Critical** | `severity/S1` | SDK is unusable for a core workflow with **no workaround**; data integrity or security at risk. | Cannot establish a client connection at all against a supported Canton/Ledger-API-v2 participant; `SubmitAndWait` silently drops or corrupts commands; panic on a core code path; generated code is wrong in a way that produces incorrect ledger writes; any security vulnerability (→ private path). |
| **S2 — High** | `severity/S2` | A major feature is broken or returns wrong results; workaround is missing or impractical. | An entire service returns incorrect data (e.g. `StateService` omits active contracts); `godaml` emits non-compiling Go for a valid `.dar`; auth refresh fails for a whole provider; a streaming call never delivers events for valid filters. |
| **S3 — Medium** | `severity/S3` | A feature is broken or degraded but has a **reasonable workaround**, or the problem affects a narrow edge case. | One optional field mis-marshals but a manual map works; codegen mangles an uncommon type; an error is misclassified; flaky behavior under specific timing. |
| **S4 — Low** | `severity/S4` | Minimal functional impact. | Documentation gaps, typos, cosmetic log wording, non-blocking warnings, nice-to-have ergonomics. |

If severity is unclear at intake, the issue is labeled `severity/needs-triage`
and treated **as the next level up** until confirmed (err toward urgency).

## Response SLE matrix

"**Initial response**" = a maintainer has triaged the issue: acknowledged it,
applied severity/type labels, and either started investigating, asked for the
specific missing information, or set expectations. It is **not** a promise of a
fix in that window.

Targets are measured in **business days** (see [Definitions](#definitions)) from
the moment the issue is opened with enough detail to triage.

| Severity | Initial response (target) | Update cadence while open |
| --- | --- | --- |
| **S1 — Critical** | within **3 business days** | every **3 business days** |
| **S2 — High** | within **5 business days** | weekly |
| **S3 — Medium** | within **10 business days** | as progress occurs |
| **S4 — Low** | best-effort | none committed |

> The S1 target of **3 business days for critical integration issues** is the
> headline commitment of this policy. Other tiers scale from it. These are
> deliberately community-OSS-appropriate windows; a [commercial tier](#support-tiers--escalation)
> can offer materially tighter, hour-scale response.

A few realistic published points of comparison: GitLab Support commits to 4h
(S1) / 8h (S3) / 24h (S4) **business-hour** responses on paid plans; generic
commercial SLAs often use 5 business hours (Critical) → 5 business days (Low).
Community/best-effort projects commonly publish multi-business-day windows like
the ones above (see [Sources](#sources)).

## Resolution targets

Resolution is best-effort and depends on root cause, upstream dependencies
(DAML/Canton, `dazl-client`), and reproducibility. We commit to **transparency**,
not to a fix date.

| Severity | Resolution posture |
| --- | --- |
| **S1** | Worked actively until a fix, a documented workaround, or a downgrade path is available; prioritized for the **next patch release**. |
| **S2** | Targeted for an upcoming minor/patch release; a workaround is documented as soon as one is known. |
| **S3** | Scheduled into the normal backlog; may wait for a related change. |
| **S4** | Addressed opportunistically or accepted as a good first issue / community PR. |

When a fix can't be committed (out of scope, won't-fix, upstream), the issue is
closed with a written rationale and, where possible, a pointer to the upstream
tracker.

## Issue type & status labels

Type and status are orthogonal to severity. Severity says *how bad*; type says
*what kind*; status says *where it is in the pipeline*.

| Category | Labels |
| --- | --- |
| **Type** | `type/bug`, `type/integration`, `type/codegen`, `type/docs`, `type/enhancement`, `type/question` |
| **Severity** | `severity/S1` … `severity/S4`, `severity/needs-triage` |
| **Status** | `status/needs-info`, `status/confirmed`, `status/in-progress`, `status/blocked-upstream`, `status/wontfix`, `status/duplicate` |
| **Area** | `area/client`, `area/auth`, `area/ledger`, `area/admin`, `area/topology`, `area/codec`, `area/codegen` |
| **Community** | `good-first-issue`, `help-wanted` |

A well-triaged issue ends up with **one** type, **one** severity, **one or more**
area labels, and a status.

## The triage workflow

For every newly opened issue a maintainer:

1. **Acknowledges & timestamps** — the [response SLE](#response-sle-matrix) clock
   starts when the issue has enough detail to act on.
2. **Checks completeness** — for integration/bug reports we need: SDK version
   (module version / commit), Go version, target Canton/Ledger-API version
   (`GetLedgerAPIVersion` output), the failing call, the **full error** (run it
   through `errors.AsDamlError` — see [troubleshooting.md](./troubleshooting.md#classifying-any-error)),
   and a minimal repro. If missing, apply `status/needs-info` and ask **once,
   specifically**.
3. **De-duplicates** — link to an existing issue and apply `status/duplicate`
   if already tracked.
4. **Classifies** — apply type, area, and severity labels per the tables above.
5. **Confirms or reproduces** — reproduce where feasible; apply
   `status/confirmed` or request more evidence.
6. **Routes** — assign an owner / milestone for S1–S2; backlog S3–S4. If the
   root cause is in DAML/Canton or `dazl-client`, apply `status/blocked-upstream`
   and link the upstream issue.
7. **Communicates** — post the initial response and keep the
   [update cadence](#response-sle-matrix) until the issue is resolved or closed.

## Security issues

**Do not open public issues for vulnerabilities.** Report privately via GitHub
Security Advisories ("Report a vulnerability" on the repository's **Security**
tab) or by emailing the maintainers. Security reports get the **fastest** handling
regardless of the matrix above:

| Stage | Target |
| --- | --- |
| Acknowledge receipt | within **2 business days** |
| Initial assessment & severity | within **5 business days** |
| Fix / coordinated disclosure | scheduled with the reporter; embargoed until a release is available |

This repository also runs automated dependency/vulnerability scanning in CI
(`.github/workflows/vuln.yml`); advisories it surfaces are triaged as S1/S2 as
appropriate.

## Stale & inactive issues

To keep the queue meaningful:

- An issue in `status/needs-info` with **no reply for 21 calendar days** is
  marked `stale`; after a further **14 days** it may be closed as
  `not-reproducible`. It can be reopened the moment the requested detail arrives.
- `wontfix` / `duplicate` issues are closed immediately with a rationale.
- Closing for inactivity is **not** a severity judgment — a clear repro reopens
  it and restarts the SLE clock.

(This mirrors the lifecycle-bot pattern used by large OSS projects such as
Kubernetes: *stale → rotten → closed*, always reversible.)

## Support tiers & escalation

| Tier | Channel | Response |
| --- | --- | --- |
| **Community** (default) | GitHub issues on `noders-team/go-daml` | Best-effort, per the [SLE matrix](#response-sle-matrix). |
| **Priority / commercial** | Arranged directly with the maintainers (noders-team) | Negotiated, contractual SLA with tighter (hour-scale) windows and named contacts. |

Escalate an existing community issue by commenting with new evidence of impact
(production outage, data integrity, security) — a maintainer will re-evaluate the
severity. Severity can move **up or down** as facts change.

## Definitions

- **Business day** — Monday–Friday, excluding public holidays observed by the
  maintaining team. Targets exclude weekends and holidays.
- **Initial response** — first substantive maintainer reply that triages the
  issue (labels + acknowledgement + next step). Automated bot comments don't
  count.
- **Workaround** — a documented way to achieve the intended outcome without the
  broken path (e.g. building command maps by hand instead of via generated
  helpers).
- **Resolution** — a merged fix, a documented permanent workaround, or a
  reasoned won't-fix/closed-as-upstream.

## Sources

Open-source and industry references this policy draws on:

- [A guide to SLEs and SLAs for open source projects — Opensource.com](https://opensource.com/article/23/2/sle-sla-open-source-projects) — the SLE-vs-SLA distinction for community projects.
- [Issue Triage — The GitLab Handbook](https://handbook.gitlab.com/handbook/engineering/infrastructure/engineering-productivity/issue-triage/) and [GitLab Support definitions](https://about.gitlab.com/support/definitions/) — S1–S4 severity model and business-hour response targets.
- [Tutorial: Set up a project for issue triage — GitLab Docs](https://docs.gitlab.com/tutorials/issue_triage/) — labelling/triage workflow.
- [Severity Levels Explained (SEV1–SEV5) — UptimeRobot](https://uptimerobot.com/knowledge-hub/monitoring/severity-levels-explained/) — severity taxonomy.
- [Response and Resolution Times in SLA — Indeed](https://www.indeed.com/career-advice/career-development/response-and-resolution-times-in-sla) and [SLA Severity Levels — Atlas Systems](https://www.atlassystems.com/blog/sla-severity-levels) — response/resolution-time framing and example windows.
