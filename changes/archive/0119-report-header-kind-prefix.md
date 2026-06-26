---
type: Change Case
title: Report header kind prefix and title-first layout
description: Prefix each Evaluation report's H1 with its kind (`Area:` / `Factor:` / `Requirement:`), move the title above the navigation trails, drop the now-redundant `Path:` / `Name:` identifier line, and lock the `Area:` trail's root element to the model title.
status: Done
tags: [evaluation, reports, markdown]
timestamp: 2026-06-26T00:00:00Z
---

# Report header kind prefix and title-first layout

This parent concept captures the why and status; the detail lives in its
children:

- [Functional spec](0119-report-header-kind-prefix/spec.md) - what the change
  must do.
- [Design doc](0119-report-header-kind-prefix/design.md) - how the renderers
  deliver it, and the alternatives weighed.

## Motivation

Every generated report opens with its navigation trail and then a **bare** H1 —
just the subject's display title, with no statement of what *kind* of thing the
report covers. The kind is only inferable: from the trail label, and from a
secondary `Path:` (Area/Factor) or `Name:` (Requirement) identifier line below
the title. Because an Area, a Factor, and a Requirement can share a display
title, a report opened cold — or linked from outside the run — does not say
plainly whether it is an Area, Factor, or Requirement report.

Two coupled changes fix this:

- **Kind in the title.** Prefix each report's H1 with its kind —
  `# Area: <title>`, `# Factor: <title>`, `# Requirement: <title>` — so the
  report self-identifies in its most prominent line. This puts the subject kind
  where the subject's *identity* already lives (the title), which is consistent
  with the existing rule that the header summary *table* should not repeat the
  kind as metadata ([0104](archive/0104-evaluation-v2-report-header-navigation.md));
  the title is the right home for the kind, the state table is not.
- **Title first.** Render the H1 as the report's first line, with the `Area:`
  trail (and the Factor `Factor:` trail / Requirement `Factors:` line) following
  it. Trail-first makes navigation chrome the opening line, which reads
  backwards; a kind-prefixed title-first header states *what this is* before
  *where it sits*.

Together these make the `Path:` / `Name:` identifier line redundant: the trail
already carries **location** and the kind prefix now carries **kind**. The line's
only remaining job was surfacing the canonical structural ID — which is not lost,
because report filenames and `Data` links already derive from structural IDs, not
titles ([Report Paths](../specs/evaluation/reports/report-tree.md)). So the case
drops the `Path:` / `Name:` line.

Finally, the case locks one existing-but-unspecified behavior the new title
surfaces more prominently: the `Area:` trail's root element renders the **model
title**, not a generic "Root". The renderer already does this
(`areaTitle(spec, nil)` returns the Model `title`); promoting it to the durable
contract guards it against regression now that the model title appears in both
the root report's kind-prefixed H1 and its trail.

## Scope

Covered:

- a kind prefix (`Area:` / `Factor:` / `Requirement:`) on every generated
  report's H1, before the subject display title;
- title-first header ordering — the H1 renders first, and the `Area:` trail, the
  Factor `Factor:` trail, and the Requirement `Factors:` context line all follow
  it, keeping their existing relative order and content;
- removal of the `Path:` (Area, Factor) and `Name:` (Requirement) identifier
  lines from report headers;
- locking the `Area:` trail's root element to the Model `title`; and
- the durable report-tree Navigation rules, per-report MUST-include lists, and
  Rendering Rules rationale, reconciled with the new layout.

Deferred / non-goals:

- **the `Area:` prefix / `Area:` trail echo is accepted, not fixed.** On Area
  reports the kind prefix and the trail label are the same word, so `# Area: X`
  sits directly above `Area: <model> / X`. This case does **not** relabel the
  `Area:` trail (e.g. to `In:` / `Location:`) — the prefix is a kind label and
  the trail is navigation; their coincidence for areas is tolerated here and may
  be revisited separately;
- no change to the no-model-title fallback for the root trail element (the
  renderer's existing `Root Area` fallback is out of scope);
- no change to which trails, Factor links, or `Factors:` set a report renders, or
  to their separators, ordering, or link targets — only the title's position
  relative to them;
- no change to the header summary tables' columns, to ratings, assessments,
  findings, confidence, summaries, drivers, limits, or any other report content;
- no change to structured routine data, `EvaluationOutputResult`, JSON field
  names or raw values, report paths, filenames, or `Data` links; and
- no migration of existing completed evaluation runs.

## Affected artifacts

### Code

- [x] `internal/evaluation/report_tree.go` - in `renderEvaluationAreaReport`,
      `renderEvaluationFactorReport`, and `renderEvaluationRequirementReport`:
      write the kind-prefixed H1 (`# Area:` / `# Factor:` / `# Requirement:`
      + title) as the first line, move the `writeEvaluationAreaTrail` /
      `writeEvaluationFactorTrail` / `writeEvaluationRequirementFactorsLine`
      calls to render *after* the title, and delete the `Path:` / `Name:` lines.
      The root-element model title already resolves through `areaTitle(spec,
      nil)`; no change there beyond the test lock.
- [x] `internal/evaluation/evaluation_test.go` - update report header
      assertions: expect the kind-prefixed H1 as the first line, expect the
      `Area:` trail (and Factor `Factor:` / Requirement `Factors:` lines) after
      the title, assert no `Path:` / `Name:` line is rendered, and assert the
      `Area:` trail's root element renders the Model `title`.

### Format spec

- [x] None - `SPECIFICATION.md` does not govern Evaluation generated report
      presentation (consistent with 0117 and 0118). (Deliberate.)

### Durable specs

- [x] `specs/evaluation/reports/report-tree.md` - Navigation: replace the
      "Every report MUST start with an `Area:` navigation trail" rule with the
      title-first, kind-prefixed-H1 rule and require the trails to follow the
      title; add the root-element-renders-Model-`title` rule; forbid the
      `Path:` / `Name:` identifier line. Per-report MUST-include lists: drop
      "and path" / "and name" from the Area, Factor, and Requirement entries and
      note the kind prefix. Rendering Rules: reconcile the 0104 rationale (the
      `path/name line` it cites is gone; the kind now rides the title, still not
      the table).

### Durable docs / bundled skill

- [x] `specs/skills/quality-skill/reporting.md` - reviewed; it defers navigation
      trails to the report tree and does not redescribe the header shape. No
      change expected.

### Suggested new durable specs

- None.

## Status

`Done`. Implemented, verified, and archived with the report renderer, focused tests, and durable report-tree spec in sync. See the [status lifecycle](index.md#status-lifecycle).
