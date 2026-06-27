---
type: Design Doc
title: Finding Basis Design
description: Implementation approach for renaming finding-local cause posture to basis.
tags: [evaluation, records, reports, terminology]
timestamp: 2026-06-27T00:00:00Z
---

# Finding Basis Design

## Context

This design answers the [Finding Basis functional spec](spec.md). The change is
a clean evaluation-record vocabulary rename: `cause` becomes `basis` in
Requirement Finding data, reports, examples, and skill guidance while preserving
the nested status model.

## Approach

Treat the rename as a source-level contract change, not as a translation layer.
The data contract should require `basis` directly, examples should emit `basis`,
and report rendering should read `basis` directly. Existing `cause` helper names
should be renamed where they describe the Finding Core field, so code review and
future searches point at the current contract.

The implementation has four coordinated edits:

1. Update durable text surfaces: the root format spec, bundled spec copy,
   evaluation record/routine/report specs, skill evaluation spec, runtime skill,
   scaffold comment, and release notes.
2. Update the data contract and examples in `internal/evaluation` so schema and
   example commands expose `basis`.
3. Update report rendering to render `Basis` and `Basis Evidence`, and to read
   `finding["basis"]`.
4. Update tests and fixtures to assert the new field and labels.

No parser, migration, or aliasing path is added. Historical archived change
cases and append-only historical logs keep their original wording.

## Spec Response

The data contract edit satisfies the required-field and no-`cause` contract.
Because the generated JSON schema is derived from the same contract shape in
source, updating the source contract and regenerated/checked schema keeps the
machine-readable artifact aligned.

Report rendering changes are literal label and field-source changes: summary
tables use the `Basis` column, detail sections use `Basis`, and nested evidence
uses `Basis Evidence`.

Skill and routine guidance preserve the existing no-overclaim rule by changing
only the field owner: `cause.status` becomes `basis.status`, with the same enum
and evidence requirement.

## Alternatives

**Keep `cause` in data and render `Basis` in reports.** Rejected because it
would leave the schema vocabulary awkward for agents authoring structured
findings. The problem is not only presentation; the field name itself guides
judgment.

**Accept both `cause` and `basis` temporarily.** Rejected because QUALITY.md is
early alpha and the repository guidance prefers clean breaks over compatibility
shims unless an active release task requires them.

**Use `attribution` instead of `basis`.** Rejected for this case because
`attribution` is precise but longer and heavier. `basis` is short enough for
reports and still neutral across strengths, gaps, risks, unknowns, and notes.

**Rename nested `status` values.** Rejected because the existing enum describes
support posture well; only the owning field name is wrong.

## Trade-offs & Risks

This is a breaking change for any hand-authored or external evaluation data that
still writes `cause`. That is acceptable under the early-alpha compatibility
policy, but the release note should name it clearly.

`basis` already appears in prose such as "evidence basis". The durable spec and
skill guidance need to distinguish finding-local `basis` from report-level
evidence summaries so agents do not duplicate evidence provenance inside the
field.

Search results for ordinary English "because" and historical `cause` records
will remain noisy. Completion should be verified by searching active surfaces
for Finding Core uses of `cause`, not by deleting every English occurrence.

## Open Questions

None.
