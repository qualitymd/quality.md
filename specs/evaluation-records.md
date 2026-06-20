---
type: Functional Specification
title: Evaluation records
description: The deterministic on-disk contract for QUALITY.md evaluation run records.
tags: [evaluation, records, cli, skill]
timestamp: 2026-06-18T00:00:00Z
---

# Evaluation records

This spec defines the runtime record contract for a QUALITY.md evaluation run:
folder names, record names, record schemas, `schemaVersion`, and the division of
responsibility between the deterministic `qualitymd` CLI and the judging skill.
The evaluation semantics are defined by
[`SPECIFICATION.md` → Evaluation](../SPECIFICATION.md#evaluation).

This contract is a standalone spec — not prompt prose, and not nested under
`cli/` — because it has two consumers: the CLI that writes records and the skill
that supplies judgment. A single cited source of truth keeps the two surfaces
from drifting in a way duplicated prose could not.

This document uses BCP 14 keywords only for testable conformance requirements.
The key words "MUST", "MUST NOT", "SHOULD", and "MAY" are to be interpreted as
described in [RFC 2119](../docs/reference/rfc2119.md) and
[RFC 8174](../docs/reference/rfc8174.md) when, and only when, they appear in all
capitals.

## Responsibility

The CLI **MUST** own file creation, serialization, run-folder numbering,
record numbering, `model.md` snapshotting, `schemaVersion` stamping, and
`report-summary.md` / `report.md` / `report.json` rendering. The skill supplies judgment content:
findings, ratings, rationales, roll-up judgment, and recommendations. The CLI
owns the `model.md` snapshot — rather than the skill — because the snapshot is
mechanically resolvable content, and keeping the record of *what was evaluated*
off the judging skill preserves the CLI-writes / skill-judges division.

The skill **MUST NOT** hand-author, number, serialize, or stamp evaluation
records when the corresponding CLI command exists. A separately tracked
numbering counter previously drifted and produced a real run-number collision;
numbering is therefore CLI-owned and derived from a single directory scan
(one past the highest present), so the on-disk folders are the single source of
truth and two writers cannot claim the same number from a stale counter.

## Runtime, Not OKF

Evaluation records are raw runtime outputs, not OKF concepts. A run folder
**MUST NOT** be treated as an OKF bundle and **MUST NOT** contain OKF
`index.md`, `log.md`, or `schema.md` semantics. Runtime Markdown frontmatter in
recommendation records is machine metadata, not OKF concept frontmatter.

## Run Folder

Each run folder is named:

```text
NNNN-<altitude>[-<narrowing>]-quality-eval
```

`<altitude>` is `subject`. Historical runs may use `model`, but new run creation
is subject-only. `NNNN` is one zero-padded sequence across the evaluation
directory, computed by directory scan immediately before creation (one past the
highest present). An existing folder at the computed name **MUST** fail rather
than be overwritten or silently retried at the next number: under a deterministic
shared sequence a collision signals concurrent or corrupt state that must
surface, not be papered over. The scan-then-create step is safe but not
serializable — concurrent invocations are resolved by the loser failing on an
existing-folder error — and a lock is deferred until concurrent use is a real
requirement.

A run folder contains:

```text
model.md
design.md
plan.md
assessments/
  NNN-<target>-<requirement>.json
analysis/
  <target>.json
recommendations/
  NNN-<slug>.md
report.md
report-summary.md
report.json
```

`assessments/` and `recommendations/` each use their own local `NNN`
sequence.

`plan.md` may be body-only, or it may carry YAML frontmatter with optional
planned coverage metadata.

## Schema Version

Every JSON record (`assessments/*.json`, `analysis/*.json`,
`report.json`)
**MUST** carry top-level `schemaVersion: 1`.

Every CLI-written recommendation Markdown record **MUST** carry runtime YAML
frontmatter with `schemaVersion: 1`.

## Historical and Non-CLI Records

The current CLI writer is strict: new records **MUST** satisfy the active
contract and carry the active `schemaVersion`. Readers that inspect evaluation
history are tolerant: historical, partial, hand-edited, copied, or non-CLI
records can be present in a run folder without making ordinary history
inspection fail.

An individual record that cannot be trusted under the current contract makes
only that run non-reportable. Status/list readers **SHOULD** surface it as a
run gap that names the record path and reason, preserving any run metadata and
record-file counts that can be determined without trusting the malformed
payload. Tools **MUST NOT** migrate, rewrite, silently skip, or reinterpret old
record shapes as a compatibility mechanism. A fresh evaluation or explicit
correction through the current CLI is the forward path.

At minimum, incompatible-record gaps distinguish malformed JSON or runtime
frontmatter, unreadable records, missing `schemaVersion`, unsupported
`schemaVersion`, and structurally incomplete current-schema records.

## Assessment Result Record

An assessment result record is one JSON file per evaluated requirement. It is
the result of carrying out the Requirement's authored `assessment` instruction
against run evidence. Required fields:

- `schemaVersion`
- `targetPath`
- `requirement`
- `factorPaths`
- `ratingResult`
- `criterionSource`
- `findings`
- `recommendations`
- `supersedes`, optional

`targetPath`, `factorPaths`, and `ratingResult.level` values are stable model
identifiers: ordered Target paths, ordered Factor paths, and rating `level` ids.
They are not human display titles.

`ratingResult` **MUST** be an object with:

- `kind`, either `rated` or `not-assessed`
- `level`, required when `kind` is `rated` and omitted when `kind` is
  `not-assessed`
- `rationale`

The `kind` value is a typed rating-result state. A rated result without a
`level`, a not-assessed result with a `level`, an unknown `kind`, or an empty
`rationale` makes the record invalid for reporting.

Each finding **MUST** carry `locator`, `observation`, `category`, and
`severity`; it **MAY** carry `evidence` and `attributes`.

`severity` **MUST** be one of the canonical severity levels below. The `level`
is the stable record value; the `title` is the human display label used by
reports.

| level      | title    |
| ---------- | -------- |
| `critical` | Critical |
| `high`     | High     |
| `medium`   | Medium   |
| `low`      | Low      |
| `info`     | Info     |

Findings with `severity: "info"` are neutral evidence or supporting
observations. Findings with `critical`, `high`, `medium`, or `low` are risk
findings eligible for selected-finding summaries.

Evidence verification and locator rigor ride on these existing fields
deliberately, with no new schema field; this keeps `schemaVersion` stable and
the record mechanically gate-able. Add a dedicated field only when repeated
real-repo use shows the existing fields insufficient.

`evidence[].kind` is an intentionally open classification string. Report
renderers can display or group it, but they must not assign special semantics to
undocumented kind values.

A run **MUST NOT** contain more than one active assessment result record for the
same ordered `targetPath` and `requirement`. Duplicate active assessment result
records make the run non-reportable. This uniqueness rule exists because a
correction or resume workflow that re-adds an assessment result would otherwise
append a conflicting second record while the run still reported as renderable,
producing a report whose requirement entry and roll-up disagree. A corrected
assessment result **MAY** include `supersedes`, a list of earlier assessment
result IDs or paths for the same ordered `targetPath` and `requirement`;
superseded assessment result records remain part of the run but are not active.
Analysis records **MUST** reference active assessment result records. This is
stricter than recommendation superseding (below): because analysis ratings bind
to assessment result references, a corrected assessment result **MUST** be paired
with an updated analysis, or roll-ups would silently inherit stale judgment.

## Analysis Record

An analysis record is one JSON file per target. Required fields:

- `schemaVersion`, `targetPath`
- `localRatingResult`, or `null` for a grouping target with no own requirements
- `factorRatingResults`
- `aggregateRatingResult`
- `assessmentResultRecords`
- `childAnalysisRecords`
- `ratingConstraints`, optional

Every rating result **MUST** use the explicit `ratingResult` object shape
defined above. `targetPath`, `factorRatingResults[].factorPath`, and rating
values are stable model identifiers, not human display titles. A
`localRatingResult: null` on a target with child analyses and no local
assessment result records represents a structural grouping target; report outputs
must render that as a distinct structural local-rating state, not as a missing
not-assessed rating.
When present, each `ratingConstraints` entry **SHOULD** identify the binding
`assessmentResultRecord`, `requirement`, and constrained `level`.

## Planned Coverage

`plan.md` is a YAML-frontmatter + Markdown-body artifact. Its Markdown body is
the run's prose plan. Optional frontmatter `coverage:` lists the assessment
result requirements and analysis targets intended for the run. Planned coverage
exists to support deterministic resume diagnostics; it **MUST NOT** replace the
prose plan or the evaluation design.

When `coverage:` is present, it **MUST** contain:

- `assessmentResults`
- `analyses`

Each assessment result entry **MUST** contain ordered `targetPath` and
`requirement`.
Each analysis entry **MUST** contain ordered `targetPath`.

Coverage frontmatter is hand-authored as part of `plan.md`; there is no separate
CLI write command for planned coverage.

A planned assessment result key is ordered `targetPath` plus `requirement`. A
planned analysis key is ordered `targetPath`. Duplicate planned assessment
result or analysis keys are invalid.

When `coverage:` is absent, the run keeps the same status and reportability
behavior it would have without planned coverage metadata. Malformed `coverage:`
frontmatter makes the run non-reportable through an `invalid-plan-coverage`
status gap rather than making the run unloadable.

## Recommendation Record

A recommendation record is one Markdown file per key gap. Its runtime
frontmatter **MUST** carry:

- `schemaVersion`
- `title`
- `gap`
- `evidenceLocators`
- `assessmentResultRecords`
- `remediationOptions`
- `recommendedOption`
- `doneCriterion`
- `supersedes`, optional

The Markdown body **MUST** state the gap, evidence locators, remediation options,
recommended option, and done criterion in stable human-readable sections.

When a recommendation corrects an earlier recommendation while preserving the
audit trail, it can include `supersedes`, a list of earlier recommendation
IDs or paths. Superseded recommendation records remain part of the run, but
reports treat them as inactive advice. Active selection is driven by explicit
`supersedes` intent, not by record numbering or recency: appending a corrected
recommendation without `supersedes` leaves the run reportable and renders both
files, so the report's Next Action can still point at the stale original — a
silent error. Requiring explicit superseding makes the active advice unambiguous
without making report output depend on which correction happened to be written
last. Route hints that help a reader act (affected package, path, workflow,
maintainer surface, verification command) belong in the existing recommendation
text fields rather than a dedicated schema field, for the same schema-stability
reason as assessment evidence above.

## Report Outputs

`report-summary.md`, `report.md`, and `report.json` are generated artifacts, not
input records, judgment records, or OKF concepts.

All report outputs are projections of one assembled Evaluation Report model. The
renderer **MUST** assemble that model from the run's `model.md`, `plan.md`,
assessment result records, analysis records, recommendation records, and run
metadata before rendering any output. `report-summary.md`, `report.md`,
`report.json`, and gate behavior **MUST NOT** be composed independently from run
files in ways that can diverge on scope, headline rating, active advice,
limitations, or record references.

Report rendering **MUST** preserve the meaning of the underlying records rather
than adding a new judgment layer. In particular:

- The headline verdict is the in-scope root Target's aggregate
  `ratingResult`.
- A local Target rating **MUST NOT** be presented as the headline verdict unless
  the active analysis record also makes it the aggregate rating.
- Labels **MUST NOT** imply whole-model coverage, ranking, priority, causality,
  or actionability unless that meaning is present in the records or in a
  documented deterministic renderer rule.
- Scoped evaluations **MUST** be explicitly labeled as scoped wherever their
  verdict is displayed.
- `rated`, `not-assessed`, active, superseded, and structural/no-local-rating
  states **MUST** remain distinct.
- Report JSON **MUST** expose those states as typed fields rather than requiring
  consumers to infer them from `null`, absent fields, or booleans alone.
- Findings are evidence. A report output **MUST NOT** turn findings into actions
  unless they are connected to active recommendation records.
- Recommendation-facing action surfaces **MUST** use active recommendations only;
  superseded recommendations remain audit/detail data.
- Rendering may resolve labels, sort, deduplicate, link, summarize, and apply
  documented selection rules. It **MUST NOT** reread subject source, recompute
  ratings, invent findings, or choose new recommendations by evaluator
  judgment.

`report.json` is the canonical serialized form of that assembled report model.
`report.md` is the complete human Evaluation Report. `report-summary.md` is the
concise human triage projection. They can differ in density, but they **MUST**
agree on the in-scope root verdict, scope, active recommendations, limitations,
typed rating states, and record references.

## report.json

`report.json` is the machine rendering of the same Evaluation Report as
`report.md`. It **MUST** present the in-scope root Target's aggregate verdict
and rationale, scope, per-target results, and advice. It **MUST** reference
findings by assessment result record; full finding detail stays in
`assessments/*.json`.
Referencing a finding by record (rather than inlining its raw
`observation`/`evidence`) also keeps the deterministic renderer from echoing
secret values or prompt-injection text into the report artifact, the same
trust-boundary rule the
[`report build`](cli/evaluation-report.md) renderer follows.

`report.json` **MUST** carry a leading report-summary layer equivalent to the
leading sections of `report.md`: verdict, scope, evidence basis, limitations,
next action, and target summary. Collections **MUST** render as arrays,
including empty arrays; they **MUST NOT** render as `null`.

Recommendation summaries **MUST** indicate whether each recommendation is active
or superseded through a typed lifecycle `state`. The legacy convenience `active`
boolean can be present, but it must agree with `state`. The report Next Action
**MUST** choose from active recommendations only.

Assessment result summaries **MUST** indicate whether each assessment result is
active or superseded through the same typed lifecycle `state`. Superseded
assessment results remain visible in the report audit trail but must not be
treated as active judgment.

Finding summaries **MUST** preserve the canonical severity `level` and expose
the corresponding display `title`.

Target summaries and details **MUST** include `localRating`, a typed local rating
state with `kind: rated`, `kind: not-assessed`, or `kind: structural`. Rated and
not-assessed local states include the underlying `ratingResult`; structural local
states do not. The historical `localRatingResult` field can remain as a
compatibility convenience, but consumers should use `localRating`.

The report next action **MUST** be a typed next-step object under `nextAction`
with `kind: recommendation` or `kind: none`. Recommendation next steps include
the recommendation id and path; none next steps do not.

Missing report metadata **MUST** render as objects with a stable `field` and a
human `title`, not as prose-only strings.

Equivalent limitation summaries **MUST** be deduplicated across recorded context
and rationale-derived constraints. Deduplication **MUST** be deterministic and
preserve the first displayed wording.
Derived limitation summaries **MUST** preserve locator-like text such as dotted
file paths.

Target local and aggregate ratings **MUST** render as explicit rating objects.
Human Markdown reports resolve Model, Target, Factor, and Rating Level display
labels from the run's `model.md` snapshot. `report.json` preserves the stable
identifiers from assessment and analysis records.

## report-summary.md

`report-summary.md` is the concise human triage artifact generated beside the
full report. It **MUST** be derived from the same report model as `report.md` and
`report.json`, link to both, and preserve the same active/superseded
recommendation distinction. It **MUST NOT** replace `report.md` as the complete
human Evaluation Report.

`report-summary.md` **MUST** use a decision-brief outline for human readers:
key details under `# Quality Evaluation Summary`, then Verdict, Target Ratings,
Selected Findings, Recommended Actions, and Scope & Limitations. The key details
**MUST** use reader-facing labels, including "Full evaluation" for an
unnarrowed run and "Evaluation verdict" for the in-scope root Target's aggregate
verdict.

The Target Ratings section **MUST** distinguish local ratings from aggregate
ratings. The selected finding section **MUST** be named for its deterministic
selection rule, such as "Selected Findings" for a severity-filtered subset or
"Rating-Binding Findings" only when the renderer can prove binding status from
recorded constraints.

The Recommended Actions section **MUST** make active recommendation identifiers
prominent for follow-up prompts. When active recommendations exist, the summary
**MUST** render a `Recommendation ID` column with copyable stable identifiers and
**MUST NOT** present superseded recommendations as primary actions.

A not-assessed rating uses `ratingResult.kind: not-assessed` and no
`ratingResult.level`. A structural grouping target with no local requirements
**MUST** render the `localRating.kind: structural` state rather than looking like
a missing-evidence not-assessed rating.
