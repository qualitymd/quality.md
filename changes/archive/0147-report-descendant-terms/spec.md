---
type: Functional Specification
title: Report descendant terms
description: Requirements for generated report labels for immediate descendant Areas and Factors.
tags: [evaluation, reports, terminology]
timestamp: 2026-06-27T00:00:00Z
---

# Report Descendant Terms

This change-case spec governs generated Evaluation report terminology for
immediate descendant Areas and Factors. It updates the durable report contracts
without changing Evaluation data identity fields or roll-up semantics.

The key words **MUST** and **MUST NOT** are to be interpreted as described in
BCP 14 when, and only when, they appear in all capitals.

## Motivation

The generated Markdown report tree is a human review surface. It should use the
same relationship names readers see in the Model vocabulary: Child Areas for
immediate Area descendants and Sub-Factors for immediate Factor descendants.
Mixing `Sub-Areas` with `Child Factors` makes equivalent report relationships
look like different concepts.

## Requirements

Generated Area reports **MUST** label the immediate descendant-Area section
`Child Areas`.

> Rationale: Areas decompose into Child Areas in the Model vocabulary; report
> readers should see that term at the section boundary.
>
> Durable spec: modify `specs/evaluation/reports/report-tree.md` - require Area
> reports to name direct descendant Areas as Child Areas.

Generated Area report Child Areas tables **MUST** label the
descendant-inclusive Area roll-up column `+ Child Areas Rating`.

> Rationale: the column is only distinct when an Area row has Child Areas, so the
> column label should name that relationship directly.
>
> Durable spec: modify `specs/evaluation/reports/report-tree.md` - replace the
> `+ Sub-Areas Rating` report table contract with `+ Child Areas Rating`.

Generated Factor reports **MUST** label the immediate descendant-Factor section
`Sub-Factors`.

> Rationale: Factors decompose into sub-factors, not child Factors, in
> human-facing Model vocabulary.
>
> Durable spec: modify `specs/evaluation/reports/report-tree.md` - require Factor
> reports to name direct descendant Factors as Sub-Factors.

Generated report empty-state rows for these sections **MUST** use the same
relationship terms as their section labels: `(no Child Areas)` for Area reports
and `(no Sub-Factors)` for Factor reports.

> Durable spec: modify `specs/evaluation/reports/report-tree.md` - align empty
> table row wording with the report labels.

Generated Markdown report labels **MUST NOT** use `Sub-Areas` or
`Child Factors` for immediate descendant report sections, table headings, rating
columns, or empty-state rows.

> Durable spec: modify `specs/evaluation/reports/report-tree.md` - add the
> negative label boundary to the rendering rules.

Skill reporting specs **MUST** use Child Areas and Sub-Factors when describing
generated report links for immediate descendants.

> Durable spec: modify `specs/skills/quality-skill/reporting.md` - align report
> navigation terminology with the report tree contract.

This change **MUST NOT** rename structured Evaluation data fields, generated file
paths, model-reference strings, or internal implementation helpers that are not
rendered as human report labels.

> Rationale: the user-facing problem is report vocabulary, not machine contract
> shape; preserving machine names keeps the change a compatible report wording
> fix.
>
> Durable spec: none.

## Durable spec changes

### To add

None

### To modify

- `specs/evaluation/reports/report-tree.md` - align generated report label,
  column, empty-state, and negative-boundary requirements with Child Areas and
  Sub-Factors.
- `specs/skills/quality-skill/reporting.md` - align skill report navigation
  terminology with the report tree contract.

### To rename

None

### To delete

None
