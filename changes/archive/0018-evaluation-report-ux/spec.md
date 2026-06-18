---
type: Functional Specification
title: Evaluation report UX - functional spec
description: Summary-first report.md and clearer report.json output for generated evaluation reports.
tags: [evaluation, report, cli, skill]
timestamp: 2026-06-18T00:00:00Z
---

# Evaluation report UX - functional spec

Companion to [Evaluation report UX](../0018-evaluation-report-ux.md). This spec
states the report-output delta for `qualitymd evaluation build-report`.

The key words **MUST**, **MUST NOT**, **SHOULD**, and **MAY** are to be
interpreted as described in IETF RFC 2119.

## Scope

Covered: generated `report.md` structure, generated `report.json` clarity,
explicit scope and limitation rendering, grouping-target display, empty
recommendation rendering, and the skill/run metadata needed to support those
outputs.

Deferred: rating inference, changed roll-up semantics, alternate output formats,
and interactive report viewing.

## Report rendering

`build-report` **MUST** keep rendering recorded judgment only. It **MUST NOT**
infer, recompute, or alter ratings, rationales, findings, limitations, or
recommendations.

`report.md` **MUST** begin with a summary layer before detailed target,
requirement, finding, and advice sections. The summary layer **MUST** include:

- Summary: subject or run label, altitude, effort when recorded, headline
  rating, not-assessed state, and headline rationale.
- Scope: what was evaluated, any narrowing, in-scope areas, and out-of-scope or
  deferred areas when recorded.
- Top Risks and Limitations: the most important failing findings, evidence
  limitations, not-assessed reasons, or confidence boundaries recorded in the
  run.
- Evidence Basis: the commands, searches, source surfaces, and execution status
  that support the headline.
- Next Action: the highest-priority recommendation when one exists, or an
  explicit no-remediation statement when no recommendation records exist.
- Target Summary: one row per in-scope target, root first, with local rating,
  aggregate rating, covered requirement count, and a short note.

After the summary layer, `report.md` **MUST** preserve the complete detailed
content required by the Evaluation Report contract: target details, requirement
results, factor roll-ups, findings, and advice.

## Scope and limitations

The report **MUST** make narrowed scope visible. A scoped evaluation **MUST NOT**
read like a whole-model verdict.

When the skill records exclusions, deferred areas, missing evidence, static-only
validation, or command execution limits, `report.md` **MUST** surface those items
in the summary layer. When no such metadata is recorded, the report **MUST**
state that scope or limitation metadata was not recorded rather than emitting
`null` or silently omitting the section.

The summary layer **MUST NOT** repeat equivalent limitation statements when the
same constraint appears in planned limitations and recorded rationales. The
renderer **MAY** normalize limitation text for deduplication, but it **MUST**
preserve the first displayed statement.

When deriving limitation summaries from prose, the renderer **MUST NOT** split
or corrupt dotted file paths or other locator-like text, such as
`docs/production-telemetry.md`.

## Grouping targets

A target with no local requirements and only child targets **MUST** render as a
structural grouping target. Its local rating display **MUST** be distinct from a
not-assessed local rating caused by missing evidence.

## report.json

`report.json` **MUST** carry the same summary-layer data in machine-readable
form. It **MUST** include:

- a non-null `scope` object;
- `recommendations: []` when no recommendation records exist;
- explicit rating objects for null or not-assessed ratings;
- grouping-target metadata when a target is structural; and
- evidence-basis and limitation summaries when recorded.

The JSON report **MUST** remain deterministic and idempotent.
Limitation summaries **MUST** use the same equivalent-statement deduplication as
`report.md`.
They **MUST** preserve locator-like text in the same way.

## Skill metadata

The `/quality` skill **MUST** record enough run metadata in `design.md`,
`plan.md`, or CLI-supported structured fields for `build-report` to render the
summary layer. At minimum, a reportable run should identify:

- effort;
- altitude;
- subject or run label;
- in-scope requirement set;
- out-of-scope or deferred areas;
- evidence commands or searches used for headline support; and
- limitations that constrain the rating.

The skill **SHOULD** keep this metadata concise and auditable. It **MUST NOT**
duplicate secret values into report metadata.
