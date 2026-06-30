---
type: Design Doc
title: Heading sentence case - design
description: Implementation approach for aligning active headings with sentence case without adding automation.
tags: [docs, specs, reports, skill, mintlify]
timestamp: 2026-06-30T00:00:00Z
---

# Heading sentence case - design

Answers the [functional spec](spec.md) for change case
[0189 - Heading sentence case](../0189-heading-sentence-case.md).

## Context

The work is mostly editorial, but one part is generated: Evaluation report
section headings and Contents labels come from the Go report renderer and tests.
`mintlify/specification.mdx` is also generated from `SPECIFICATION.md`.

The design needs to avoid two mistakes:

- hand-editing generated artifacts as if they were source; and
- adding automated heading-case enforcement after the user explicitly rejected
  automation.

## Approach

### Update the rule first

Start by updating repository guidance so the target rule is unambiguous before
mass editing. The primary rule lives in `AGENTS.md`; supporting guides receive
small local notes and examples where authors create headings:

- OKF concepts and logs;
- functional specs;
- design docs;
- Change Cases;
- generated report design; and
- release notes.

The guidance should use examples that demonstrate both sides of the rule:

- generic section names become sentence case, such as `Primary source data`;
- true names stay capitalized, such as QUALITY.md, CLI, OKF, `SPECIFICATION.md`,
  Agent Harnessability, Rating Scale, and user/model display titles.

### Align authored active files by surface

Apply the cleanup in stable groups:

1. Top-level docs and current guides.
2. `SPECIFICATION.md`, then regenerate `mintlify/specification.mdx`.
3. Active `specs/` files.
4. Runtime `skills/quality/` files and matching skill specs.
5. Active Mintlify pages not generated from the formal spec.
6. Current active `.quality` artifacts only if they are maintained as current
   project records.

Do not sweep `changes/archive/` or `.quality/evaluations/archive/`. For
append-only logs, update current or newly added headings only; do not rewrite
old dated entries just to normalize casing.

### Change generated report labels at the renderer

Generated report headings and Contents labels should be changed in
`internal/evaluation/report_tree.go`, not by editing generated output first.
Likely renderer labels include:

- `Key Details` -> `Key details`
- `Primary Source Data` -> `Primary source data`
- `Model Evaluation` -> `Model evaluation`
- `Top Findings` -> `Top findings`
- `Top Recommendations` -> `Top recommendations`
- `Area / Factor Breakdown` -> `Area / Factor breakdown`
- `Sub-Factors` -> `Sub-factors`
- `Findings Summary` -> `Findings summary`
- `Finding Details` -> `Finding details`
- `Unknowns & Missing Evidence` -> `Unknowns and missing evidence`
- `Ranked Findings` -> `Ranked findings`
- `Ranked Recommendations` -> `Ranked recommendations`

The H1 prefixes `Area:`, `Factor:`, `Requirement:`, and `Recommendation:` stay
as prefix labels, and the following model or recommendation title is preserved
exactly.

After the renderer changes, update affected tests that assert exact Markdown
strings, then regenerate checked-in report-gallery output.

### Keep non-heading surfaces stable

Do not recase:

- report frontmatter `type` values;
- enum display labels;
- table headers;
- JSON/YAML field names or values;
- model `title` values;
- command names; or
- filenames and literal values.

Some table headers are title-cased today, but changing them is a separate
compatibility/readability decision. This case only changes Markdown headings and
generated Contents labels.

### Validate without new automation

Use ad hoc inspection and existing project checks:

- `rg` or a temporary local one-off scan may be used during the change, but no
  script or check is committed.
- Run `mise run sync-spec-docs` after `SPECIFICATION.md` edits.
- Run `mise run report-gallery` after report renderer edits.
- Run focused Go tests around report rendering after updating expected strings.
- Run `mise run fmt-md-check` for docs-only phases and `mise run check` before
  the case moves to review.

## Spec response

- **Convention scope:** Updating `AGENTS.md` and relevant authoring guides gives
  the rule a single clear repository scope and puts reminders where headings are
  authored.
- **Preserved casing:** Manual review preserves proper names and model-provided
  display titles; renderer changes only touch fixed generated labels.
- **Active authored content:** Surface-by-surface editing keeps generated files
  synchronized with their sources.
- **Generated reports:** Report headings change in the Go renderer, tests, and
  report-gallery output, with durable report specs updated to match.
- **Historical boundary:** Archive folders and historical log entries are left
  untouched unless a current edit already touches them for another reason.
- **No automation:** Validation uses existing checks and throwaway local scans
  only; no new committed enforcement is added.

## Alternatives

### Add a heading-case checker

Rejected. The user explicitly approved the cleanup except for automation. A
checker would also need a growing allowlist for proper nouns, acronyms, model
titles, generated report prefixes, and quoted source casing.

### Rewrite every Markdown heading in the repository

Rejected. It would churn archived Change Cases and historical evaluation records
without improving current behavior. The historical boundary keeps the diff tied
to active sources and generated examples.

### Treat generated report section names as proper names

Rejected. Labels such as `Key details` and `Primary source data` are generic
section headings, not formal named artifacts. Formal report `type` values and
model titles still keep their current casing.

### Normalize table headers too

Rejected for this case. Table headers and report metadata are different
surfaces, often compact labels rather than prose headings. Changing them should
be driven by readability or compatibility, not by the heading convention.

## Trade-offs and risks

- Manual review can miss a heading. That is acceptable for this case because the
  user rejected automation; existing review and formatting checks remain the
  guardrail.
- Some headings contain formal QUALITY.md concepts. Over-lowercasing those would
  reduce precision, so edits should prefer preserving term-of-art casing when in
  doubt.
- Generated report string changes touch many test expectations and gallery
  files, so the implementation should keep renderer changes grouped and
  regenerate examples once after tests are updated.

## Open questions

None. The user accepted the scope decisions: sentence-case active headings and
generated Contents labels, preserve model titles, leave archives historical, and
do not add automation.
