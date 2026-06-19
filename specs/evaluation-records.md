---
type: Functional Specification
title: Evaluation records
description: The deterministic on-disk contract for QUALITY.md evaluation run records.
tags: [evaluation, records, cli, skill]
timestamp: 2026-06-18T00:00:00Z
---

# Evaluation records

This spec defines the runtime record contract for a `QUALITY.md` evaluation run:
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
planned-coverage.json
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

`assessments/` and `recommendations/` each use their own local `NNN` sequence.

`planned-coverage.json` is optional.

## Schema Version

Every JSON record (`planned-coverage.json`, `assessments/*.json`,
`analysis/*.json`, `report.json`) **MUST** carry top-level `schemaVersion: 1`.

Every CLI-written recommendation Markdown record **MUST** carry runtime YAML
frontmatter with `schemaVersion: 1`.

## Assessment Record

An assessment record is one JSON file per assessed requirement. Required fields:

- `schemaVersion`
- `target`
- `targetPath`
- `requirement`
- `factors`
- `rating`, or `null`
- `notAssessed`
- `criterionSource`
- `findings`
- `rationale`
- `recommendations`
- `supersedes`, optional

When `notAssessed` is `true`, `rating` **MUST** be `null`. Each finding **MUST**
carry `locator`, `observation`, and `category`; it **MAY** carry `severity`,
`evidence`, and `attributes`. Evidence verification and locator rigor ride on
these existing fields deliberately, with no new schema field; this keeps
`schemaVersion` stable and the record mechanically gate-able. Add a dedicated
field only when repeated real-repo use shows the existing fields insufficient.

A run **MUST NOT** contain more than one active assessment record for the same
ordered `targetPath` and `requirement`. Duplicate active assessment records make
the run non-reportable. This uniqueness rule exists because a correction or
resume workflow that re-adds an assessment would otherwise append a conflicting
second record while the run still reported as renderable, producing a report
whose requirement entry and roll-up disagree. A corrected assessment **MAY**
include `supersedes`, a list of earlier assessment IDs or paths for the same
ordered `targetPath` and `requirement`; superseded assessment records remain
part of the run but are not active. Analysis records **MUST** reference active
assessment records. This is stricter than recommendation superseding (below):
because analysis ratings bind to assessment references, a corrected assessment
**MUST** be paired with an updated analysis, or roll-ups would silently inherit
stale judgment.

## Analysis Record

An analysis record is one JSON file per target. Required fields:

- `schemaVersion`, `target`, `targetPath`
- `localRating`, or `null` for a grouping target with no own requirements
- `factorRatings`
- `aggregateRating`
- `assessmentRecords`
- `childAnalysisRecords`

Every rating result **MUST** record `notAssessed` distinctly from a rating level.

## Planned Coverage

`planned-coverage.json` is an optional run-root JSON artifact that lists the
assessment requirements and analysis targets intended for the run. It exists to
support deterministic resume diagnostics; it **MUST NOT** replace `design.md` or
`plan.md`.

When present, it **MUST** contain:

- `schemaVersion`
- `assessments`
- `analyses`

Each assessment entry **MUST** contain ordered `targetPath` and `requirement`.
Each analysis entry **MUST** contain ordered `targetPath`.

The CLI **MUST** write `planned-coverage.json` through
`qualitymd evaluation set-planned-coverage`; the skill **MUST NOT** hand-author
or hand-repair it when that command is available.

A planned assessment key is ordered `targetPath` plus `requirement`. A planned
analysis key is ordered `targetPath`. Duplicate planned assessment or analysis
keys are invalid.

When `planned-coverage.json` is absent, the run keeps the same status and
reportability behavior it would have without planned coverage metadata.

## Recommendation Record

A recommendation record is one Markdown file per key gap. Its runtime
frontmatter **MUST** carry:

- `schemaVersion`
- `title`
- `gap`
- `evidenceLocators`
- `assessmentRecords`
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

## report.json

`report-summary.md`, `report.md`, and `report.json` are generated artifacts, not
input records, judgment records, or OKF concepts.

`report.json` is the machine rendering of the same Evaluation Report as
`report.md`. It **MUST** present the in-scope root rating and rationale, scope,
per-target results, and advice. It **MUST** reference findings by assessment
record; full finding detail stays in `assessments/*.json`. Referencing a finding
by record (rather than inlining its raw `observation`/`evidence`) also keeps the
deterministic renderer from echoing secret values or prompt-injection text into
the report artifact, the same trust-boundary rule the
[`build-report`](cli/evaluation-build-report.md) renderer follows.

`report.json` **MUST** carry a summary layer equivalent to the leading sections
of `report.md`: summary, scope, evidence basis, limitations, next action, and
target summary. Collections **MUST** render as arrays, including empty arrays;
they **MUST NOT** render as `null`.

Recommendation summaries **MUST** indicate whether each recommendation is active
or superseded. The report Next Action **MUST** choose from active
recommendations only.

Assessment summaries **MUST** indicate whether each assessment is active or
superseded. Superseded assessments remain visible in the report audit trail but
must not be treated as active judgment.

Equivalent limitation summaries **MUST** be deduplicated across recorded context
and rationale-derived constraints. Deduplication **MUST** be deterministic and
preserve the first displayed wording.
Derived limitation summaries **MUST** preserve locator-like text such as dotted
file paths.

Target local and aggregate ratings **MUST** render as explicit rating objects.

## report-summary.md

`report-summary.md` is the concise human triage artifact generated beside the
full report. It **MUST** be derived from the same report model as `report.md` and
`report.json`, link to both, and preserve the same active/superseded
recommendation distinction. It **MUST NOT** replace `report.md` as the complete
human Evaluation Report.
A not-assessed rating uses `rating: null` and `notAssessed: true`. A structural
grouping target with no local requirements **MUST** render a distinct structural
local-rating state rather than looking like a missing-evidence not-assessed
rating.
