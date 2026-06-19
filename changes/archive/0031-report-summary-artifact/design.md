---
type: Design Doc
title: Evaluation report summary artifact - design doc
description: How build-report generates report-summary.md from the existing report model.
tags: [evaluation, report, cli, design]
timestamp: 2026-06-19T00:00:00Z
---

# Evaluation report summary artifact - design doc

Design behind the
[Evaluation report summary artifact](../0031-report-summary-artifact.md) change
and its [functional spec](spec.md). The spec fixes *what* the generated summary
artifact must provide; this doc covers the intended approach.

## Context

The current evaluation renderer already assembles one shared `ReportJSON` value
in `internal/evaluation/report.go`. That value carries the summary layer used by
both `report.md` and `report.json`: root rating, scope, evidence basis,
limitations, next action, target summary, finding summaries, and recommendation
states. `report.md` is summary-first, but it still includes the full target,
requirement, finding, and advice detail in the same file.

`report-summary.md` should therefore be a third projection of the same assembled
report model, not a new model and not a new evaluation record.

## Approach

Keep `Run.Report()` as the single place that converts loaded run records into
the report model. Add a small Markdown renderer:

```go
func renderReportSummaryMarkdown(report ReportJSON) []byte
```

`BuildReport` should assemble all render bytes before writing:

1. Load and validate the run as reportable.
2. Build `ReportJSON` with `run.Report()`.
3. Render `report-summary.md` from `ReportJSON`.
4. Render `report.md` from `ReportJSON`.
5. Marshal `report.json` from `ReportJSON`.
6. Write the three artifacts through the existing report-file write path.

This keeps all three generated artifacts tied to the same recorded judgment and
the same active/superseded recommendation state.

### Summary renderer

The summary renderer should use only fields already exposed by `ReportJSON`:

- `Summary`, `Rating`, and `Scope` for run identity, effort, narrowing, root
  rating, and boundaries;
- `FindingSummaries` and `Limitations` for top risks and confidence limits;
- `TargetSummary` for the compact rating table;
- `NextAction` and `Recommendations` for links to active recommendations; and
- the same display helpers used by `report.md` for ratings, paths, lists, and
  table escaping.

The renderer should intentionally omit target details, per-requirement details,
the full findings list, and superseded recommendation detail. It can still link
to `report.md` and `report.json` at the top, and it should link active
recommendation records in the next-action section.

Top risks should reuse the same `riskFindings`, `firstFindingSummaries`, and
`firstStrings` helpers that `report.md` already uses. That keeps truncation,
severity handling, and limitation display consistent between the summary and
the full report.

### Build result and CLI output

Extend `BuildResult` with `ReportSummaryMD string` using a stable JSON field such
as `reportSummaryMd`. The human `build-report` output should report all three
written paths.

The gate result should remain unchanged: `--fail-at-or-below` still gates on the
root rating from the same `BuildResult`, independent of how many artifacts were
written.

### Artifact writes

The current write path already uses deterministic bytes and replacement writes
for existing report files. The implementation should compute all render bytes
before writing any file so render or JSON marshal failures cannot leave a newly
rendered subset.

Writing three sibling artifacts sequentially still means a filesystem failure
can leave an older artifact beside newer ones. That is not new to this command:
the existing two-file writer already has the same cross-file atomicity limit.
The practical mitigation is to keep each individual file write atomic and make
the command idempotent so a rerun restores consistency.

### Tests

Extend the existing `internal/evaluation` report tests instead of creating a
separate fixture family:

- assert `BuildReport` writes `report-summary.md`;
- assert the summary contains run identity, root rating, links to `report.md`
  and `report.json`, top risks/limitations, target summary, and next action;
- assert the summary omits detailed sections such as full requirements and the
  full findings audit trail;
- assert a second render leaves `report-summary.md` byte-identical; and
- assert superseded recommendations are not selected as next action.

Update `internal/cli` tests around `evaluation build-report` output so both
human and `--json` forms expose the new path.

## Alternatives

**Slice the first sections out of `report.md`.** Rejected. It would couple the
summary artifact to Markdown section order and make future report-body edits
riskier. Rendering from `ReportJSON` is simpler and avoids parsing generated
Markdown.

**Add `report-summary.json`.** Rejected for this change. `report.json` already
contains the machine-readable summary layer; adding a second JSON artifact would
duplicate the same data without a demonstrated consumer.

**Make the skill author the summary.** Rejected. The artifact is mechanical
rendering of recorded judgment, so the CLI should own it. A skill-authored
summary would drift from `report.md` and `report.json`, and would weaken the
existing CLI-writes / skill-judges boundary.

**Create a separate summary report model.** Rejected. A second model would need
its own active recommendation selection, risk truncation, limitation handling,
and rating display logic. Reusing `ReportJSON` keeps those choices in one place.

## Trade-offs and risks

The main trade-off is one more generated file per run. That is acceptable because
the run folder is already a durable artifact bundle, and the new file is small.

The main risk is drift between the compact summary and the full report. The
design controls that by deriving both from `ReportJSON` and sharing display
helpers where possible.

A second risk is overloading `report-summary.md` until it becomes another full
report. Keep it narrow: headline, risks, limitations, target summary, next
action, and links out.

## Open questions

- Whether to cap top risks at the same count as `report.md` or a smaller count
  tuned for PR comments. Start with the shared helper and adjust only if real
  review surfaces show the artifact is still too long.
- Whether example bundles should include a full checked-in `report-summary.md`
  or document the expected shape in the example index. The functional spec allows
  either, but renderer regression tests should cover the real generated bytes.
