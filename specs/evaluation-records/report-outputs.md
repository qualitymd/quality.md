---
type: Functional Specification
title: Evaluation report outputs
description: Shared generated report-output invariants and report artifact relationships.
tags: [evaluation, records, cli, skill]
timestamp: 2026-06-22T00:00:00Z
---

# Evaluation report outputs

This spec is part of the [Evaluation records](../evaluation-records.md) contract.
It owns the independently reviewable runtime contract below; shared
responsibility, runtime-not-OKF status, schema-version, and historical-record
compatibility rules live in the parent spec.

This document uses BCP 14 keywords only for testable conformance requirements.
The key words "MUST", "MUST NOT", "SHOULD", and "MAY" are to be interpreted as
described in [RFC 2119](../../docs/reference/rfc2119.md) and
[RFC 8174](../../docs/reference/rfc8174.md) when, and only when, they appear in
all capitals.

## Report Outputs

`report-summary.md`, `report.md`, and `report.json` are generated artifacts, not
input records, judgment records, or OKF concepts. Their artifact-specific
contracts live in [`report-summary.md`](../reports/report-summary-md.md),
[`report.md`](../reports/report-md.md), and [`report.json`](../reports/report-json.md).
This spec owns the shared run-record inputs, assembled report-model invariants,
and report trust boundary.

All report outputs are projections of one assembled Evaluation Report model. The
renderer **MUST** assemble that model from the run's `model.md`, `plan.md`,
assessment result records, analysis records, recommendation records, and run
metadata before rendering any output. `report-summary.md`, `report.md`,
`report.json`, and gate behavior **MUST NOT** be composed independently from run
files in ways that can diverge on scope, headline rating, active advice,
limitations, or record references.

Report rendering **MUST** preserve the meaning of the underlying records rather
than adding a new judgment layer. In particular:

- The headline verdict is the in-scope root Area's aggregate
  `ratingResult`.
- A local Area rating **MUST NOT** be presented as the headline verdict unless
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
  documented selection rules. It **MUST NOT** reread evaluated source, recompute
  ratings, invent findings, or choose new recommendations by evaluator
  judgment.

`report.json` is the canonical serialized form of that assembled report model.
`report.md` is the complete human Evaluation Report. `report-summary.md` is the
concise human triage projection. They can differ in density, but they **MUST**
agree on the in-scope root verdict, scope, active recommendations, limitations,
typed rating states, and record references.

## report.json

`report.json` is the machine rendering of the same Evaluation Report as
`report.md`. Its artifact-specific contract lives in
[`report.json`](../reports/report-json.md). It **MUST** present the in-scope root
Area's aggregate verdict and rationale, scope, per-area results, and advice. It
**MUST** reference
findings by assessment result record; full finding detail stays in
`assessments/*.json`.
Referencing a finding by record (rather than inlining its raw
`observation`/`evidence`) also keeps the deterministic renderer from echoing
secret values or prompt-injection text into the report artifact, the same
trust-boundary rule the
[`report build`](../cli/evaluation-report.md) renderer follows.

`report.json` **MUST** carry a leading report-summary layer equivalent to the
leading sections of `report.md`: verdict, scope, evidence basis, limitations,
next action, and area summary. Collections **MUST** render as arrays,
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

Area summaries and details **MUST** include `localRating`, a typed local rating
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

Area local and aggregate ratings **MUST** render as explicit rating objects.
Human Markdown reports resolve Model, Area, Factor, and Rating Level display
labels from the run's `model.md` snapshot. They may use unqualified references
where the surrounding report context fixes the reference type, such as the
Area-specific Area Breakdown `Path` column. `report.json` preserves the
structured Area IDs and Factor IDs from assessment and analysis records through
`areaPath` and `factorPath` arrays; derived qualified references may be added
only without replacing those arrays, and unqualified references must not be
persisted there.

## report.md

`report.md` is the complete human Evaluation Report. Its artifact-specific
contract lives in [`report.md`](../reports/report-md.md). It **MUST** be derived
from the same report model as `report-summary.md` and `report.json`, preserve the
full human audit trail, and preserve the typed state distinctions required by
this shared report-output contract.

## report-summary.md

`report-summary.md` is the concise human triage artifact generated beside the
full report. Its artifact-specific contract lives in
[`report-summary.md`](../reports/report-summary-md.md). It **MUST** be derived from the
same report model as `report.md` and `report.json`, link to both, preserve typed
rating and lifecycle states, and preserve the same active/superseded
recommendation distinction. It **MUST NOT** replace `report.md` as the complete
human Evaluation Report.
