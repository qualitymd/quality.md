---
type: Design Doc
title: Evaluation report UX - design doc
description: How generated reports add a summary layer without changing recorded evaluation judgment.
tags: [evaluation, report, cli, design]
timestamp: 2026-06-18T00:00:00Z
---

# Evaluation report UX - design doc

Design behind the [Evaluation report UX](../0018-evaluation-report-ux.md)
change and its [functional spec](spec.md). The spec fixes _what_ the generated
reports must expose; this doc covers the intended approach.

## Context

The existing report renderer is mechanically correct enough to produce reports,
but the experiment pass found repeated usability gaps:

- `report.md` lacks front-loaded scope, limitations, top risks, evidence basis,
  and target summary sections.
- `report.json` emits `scope: null` and `recommendations: null` in cases where
  tools need stable objects or arrays.
- Structural grouping targets can look like not-assessed evidence gaps.
- Larger reports, especially ESLint, require too much scanning to answer basic
  reviewer questions.

The design goal is a better rendering of the same recorded judgment, not a new
judgment path.

## Approach

Keep one in-memory report model, but split rendering into two layers:

1. **Summary layer** - derived from recorded run metadata, assessments,
   analyses, recommendations, and render-time classification.
2. **Detail layer** - the existing complete target, requirement, finding, factor,
   and advice rendering.

`build-report` should assemble the summary layer before writing either
`report.md` or `report.json`. Both renderers consume the same assembled report
value so Markdown and JSON stay aligned.

### Structured report context

Add a small report-context structure to the loaded run. The implementation reads
bounded, conventional sections from the skill-authored `design.md` and
`plan.md`, then falls back to folder naming and recorded rating rationales for
older or sparse runs.

The parsed sections are deliberately narrow:

- `design.md` -> `Resolved parameters` for altitude, narrowing slug, effort, and
  scope description.
- `design.md` or `plan.md` -> `Out of scope` / `Deferred areas` for explicit
  scoped exclusions.
- `plan.md` -> `Effort` for effort/scope summary when not already recorded.
- `plan.md` -> `Planned limitations` for out-of-scope areas and summary
  limitations.

The parser only reads headings, bullets, and short paragraphs under those
bounded conventional sections, matching heading text case-insensitively so
title-case variants such as `Out of Scope` and `Deferred Areas` behave like the
documented forms. It does not treat arbitrary prose or evaluated source content
as instructions.

The structure should represent:

- subject or run label;
- altitude and effort;
- narrowing label;
- in-scope and out-of-scope areas;
- evidence basis entries;
- limitations; and
- next-action summary.

The renderer should not parse arbitrary prose for secrets or hidden meaning. If
metadata is absent, the summary says it was not recorded.

Limitations can come from planned limitations, analysis rationales, assessment
rationales, and factor rationales. The renderer should deduplicate equivalent
summary limitations with a deterministic normalized key, while preserving the
first recorded display text.

The limitation extractor should treat locator-like text as ordinary content, not
sentence punctuation. In particular, dotted paths such as
`docs/production-telemetry.md` must remain intact in summary limitations.

### Grouping-target classification

Classify a target as structural when it has no local requirements and has child
targets. Render local rating as `n/a` for these nodes while preserving the
aggregate rating. This is a render-time classification unless the record
contract later needs a persisted field.

### Machine-readable defaults

Normalize report JSON so consumers do not need to special-case absent
collections or implicit rating states:

- empty recommendation set -> `recommendations: []`;
- known scope with missing details -> non-null `scope` object with empty arrays
  or explicit "not recorded" fields;
- not-assessed/null rating -> explicit rating object carrying `rating: null` and
  `notAssessed: true`;
- structural local rating -> explicit non-rating state such as `kind:
"structural"`.

### Markdown order

Render `report.md` in this order:

1. Summary
2. Scope
3. Top Risks and Limitations
4. Evidence Basis
5. Next Action
6. Target Summary
7. Detailed target and requirement results
8. Findings
9. Advice
10. Audit Trail, when available

The first six sections answer the reviewer walkthrough. The remaining sections
preserve auditability.

## Alternatives

**Only update the skill prompt.** Rejected. The report is CLI-generated, so
prompt guidance alone cannot fix JSON shape, grouping-target display, or
deterministic section order.

**Add a second report command.** Rejected for now. The V1 shape should become
the default generated report if it preserves all existing detail and improves
first-read usability.

**Leave limitations in rationales only.** Rejected. It is mechanically simple
but failed the real-repo and ESLint reviewer walkthroughs.

## Trade-offs and risks

The main risk is mixing structured metadata with judgment prose. Keep the
summary layer shallow: it can surface recorded scope, limitations, evidence
basis, and recommendations, but ratings and rationales still come from the
assessment and analysis records.

The second risk is over-growing `report.json`. Keep full finding detail in
assessment records; `report.json` should carry summaries and references.

## Open questions

- Whether a future change should replace Markdown-section parsing with an
  explicit `report-context.json` or CLI-owned context writer.
- Whether `show-status` should eventually warn on missing report-context
  metadata. For this change, missing metadata renders as "not recorded" and does
  not block reportability.
