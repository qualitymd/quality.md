---
type: Change Case
title: Report empty-cell marker and legend
description: Render every empty scalar cell in an Evaluation report as a single em dash (`—`) instead of a blank, and add one static per-report legend that defines the marker — so an absent value is unambiguous and distinct from a not-assessed outcome.
status: Done
tags: [evaluation, reports, markdown]
timestamp: 2026-06-26T00:00:00Z
---

# Report empty-cell marker and legend

This parent concept captures the why and status; the detail lives in its
children:

- [Functional spec](0118-report-empty-cell-legend/spec.md) - what the change
  must do.
- [Design doc](0118-report-empty-cell-legend/design.md) - how the renderer
  delivers it, and the alternatives weighed.

## Motivation

Evaluation reports today render an absent scalar value as a blank table cell. A
blank is ambiguous: a reader cannot tell "no value was recorded" from "not
applicable to this row" from a rendering bug. The codebase has also accumulated
several uncoordinated empty-value treatments — a blank cell (the default through
`markdownCell`), a bare `/` for an all-empty confidence pair, an em dash for
the descendant roll-up column, and parenthetical `(no …)` rows for empty
sections — so the same "nothing here" reads differently in different places.

The fix is one consistent, visible marker for an empty *cell* — the em dash
(`—`), already used for the roll-up column — plus a single static legend per
report that says what it means. The em dash is chosen over a literal `N/A`
because "not applicable" is usually the wrong claim (the value is most often
simply *not recorded*); a neutral dash states absence without overclaiming. A
static legend is chosen over per-cell asterisks because reports are written to
disk and diffed: an asterisk anchored to the "first instance" migrates as data
changes and there is no natural "first" cell in a Markdown table, whereas a
fixed legend stays put across re-runs.

This case touches only empty *scalar cells*. The not-assessed / not-rated /
not-analyzed / blocked **outcomes** are a different thing — an evaluation result
the format requires be shown distinctly — and they keep their existing status
labels untouched.

## Scope

Covered:

- a single em-dash (`—`) marker for any empty scalar cell in a generated report,
  including the composite Confidence and analysis-status pair cells (so an absent
  component renders `— / —`, not a bare `/`);
- one static legend per generated report (root Area report, Area, Factor, and
  Requirement reports) that defines the em-dash marker; and
- the durable report-tree Rendering Rules describing both, reconciled with the
  existing roll-up em-dash rule and the status-label distinctness rule.

Deferred / non-goals:

- no change to the not-assessed / not-rated / not-analyzed / blocked status
  labels or to Rating Level labels — they stay distinct from the em dash;
- no change to the parenthetical empty whole-section rows (`(no findings)`,
  `(no rating drivers)`, `(none recorded)`) — they remain as worded empty-state
  rows and the legend does not cover them;
- the legend defines the em dash only; status labels and section rows are
  self-describing English and are not legended;
- no literal `N/A` text — the marker is the em dash;
- no change to structured routine data, `EvaluationOutputResult`, JSON field
  names or raw values, report paths, filenames, navigation trails, or links; and
- no migration of existing completed evaluation runs.

## Affected artifacts

### Code

- [x] `internal/evaluation/report_tree.go` - render empty scalar cells as `—`
      (the `markdownCell` guard plus the composite `evaluationConfidencePair` /
      `evaluationAnalysisStatusPair` builders, which assemble the pair before the
      cell guard sees it), reconciled with the existing `evaluationSubRatingCell`
      roll-up `—`; and append one static legend per report in the report
      renderers (`renderEvaluationAreaReport`, `renderEvaluationFactorReport`,
      `renderEvaluationRequirementReport`, driven by
      `renderEvaluationReportTree`).
- [x] `internal/evaluation/display.go` - decide whether `humanizeEnum`'s empty
      branch routes absent enum-like values to `—` here or leaves `""` for the
      cell layer to mark; either way an empty enum cell must end up as `—` once.
      (Implementation detail, settled at In-Progress.)
- [x] `internal/evaluation/evaluation_test.go` - update report assertions to
      expect `—` for empty cells and `— / —` for all-empty pair cells, expect the
      per-report legend, and keep the assertions that not-assessed/not-rated
      status labels remain distinct from the em dash.

### Format spec

- [x] None - `SPECIFICATION.md` does not govern Evaluation generated report
      presentation (consistent with 0117). Its `not assessed` distinctness rule
      concerns the rating *outcome*, which this case does not touch.
      (Deliberate.)

### Durable specs

- [x] `specs/evaluation/reports/report-tree.md` - Rendering Rules: add the
      empty-scalar-cell em-dash rule (including composite pair cells) and the
      per-report static-legend requirement, and reconcile both with the existing
      roll-up em-dash rule (0097/0111) and the status-label distinctness rule.

### Durable docs / bundled skill

- [x] `specs/skills/quality-skill/reporting.md` - reviewed; it defers cell-level
      rendering to the report tree (its Report Tree section). No change expected.

### Suggested new durable specs

- None.

## Status

`Done`. Implemented, verified, and archived with the report renderer, focused tests, and durable report-tree spec in sync.
