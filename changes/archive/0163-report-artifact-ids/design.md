---
type: Design Doc
title: Report Artifact IDs — design doc
description: Design for CLI-assigned Evaluation report artifact IDs.
tags: [evaluation, reports, advice, identifiers]
timestamp: 2026-06-29T00:00:00Z
---

# Report Artifact IDs — design doc

Design for [Report Artifact IDs](../0163-report-artifact-ids.md) and its
[functional spec](spec.md).

## Context

`qualitymd evaluation data set` already validates agent-written payloads,
derives storage paths, and rejects effective-data inconsistencies before writing.
That makes it the right place to assign mechanical artifact IDs: it has access
to the run manifest, existing persisted data, and the full candidate batch.

The main asymmetry is recommendation versus finding identity. Recommendations
are their own payloads and need an ID before a data path can be derived. Findings
remain local objects inside Requirement assessment payloads; their public
artifact IDs belong in the run-level `FindingRankingResult`.

## Approach

Add an assignment normalization step inside `SetData` before payload validation
and path derivation. The normalizer loads `RunManifest.number`, clones the input
payloads, and mutates only the artifact-ID fields the CLI owns.

For `RecommendationResult`:

- If `id` is absent, assign the next unused `QREC-<run>-<seq>` for that run.
- If `id` is present, validate the `QREC` grammar and run segment, then preserve
  it so an existing recommendation can be rewritten deliberately.
- Derive the recommendation data path from the persisted `id`.
- Leave ranking and coverage references as assigned IDs; a workflow that omits
  recommendation IDs writes recommendations first, then uses the receipt/data
  list to author ranking and coverage against assigned IDs.

For `FindingRankingResult`:

- Build a map from any existing persisted finding ranking entries:
  `findingRef` key -> `QFIND` ID.
- For each new ranking entry, preserve the existing ID for the same `findingRef`
  when present.
- Otherwise assign the next unused `QFIND-<run>-<seq>`.
- Keep `findingRef` as the exact trace to the Requirement Finding. The new ID is
  report/handoff identity, not a selector replacement.

The normalizer runs for dry-run too, so `data set --dry-run` receipts show the
same derived paths that a real write would use. It does not reserve IDs on
dry-run.

Reports add a thin rendering layer:

- `reportRunLine` renders `QEVAL-<run>` instead of only `#<run>`.
- Top and full finding/recommendation tables get an `ID` column.
- Recommendation detail summaries add `ID` and `Reference`.
- Finding detail sections render the ranked `QFIND` ID when ranked while leaving
  anchors based on the payload-local finding ID.

## Spec response

- CLI assignment is centralized in `SetData`, where run state and existing data
  are available.
- Recommendation IDs are persisted before path derivation, so data paths remain
  deterministic.
- Ranked finding IDs are stable across ranking rewrites because assignment is
  keyed by `findingRef`, not current rank.
- Existing Model reference parsing and rendering is untouched.
- Generated reports show citable IDs without changing source-data links or
  generated report path rules.

## Alternatives

- **Agent-authored IDs with CLI validation.** Rejected because users will treat
  the public token as the ID and duplicate prevention is mechanical CLI work.
- **UUIDs or ULIDs.** Rejected because they solve global uniqueness while making
  chat, oral reference, and report scanning worse.
- **Use title slugs as IDs.** Rejected because title wording changes and path
  friendliness should not define identity.
- **Assign finding IDs inside Requirement assessment payloads.** Rejected
  because those findings are payload-local observations; the public handoff
  surface is the ranked finding artifact.
- **Support temporary recommendation aliases in one batch.** Deferred. It would
  add another authoring identity surface. The simpler workflow writes
  recommendations first, then uses assigned IDs in ranking and coverage.

## Trade-offs & risks

Auto-assignment makes `data set` a normalizing write surface for a small set of
fields, not a pure byte-preserving persistence command. That is acceptable
because the fields are CLI-owned artifact IDs, but it should stay narrow and
explicit in specs and receipts.

Recommendation ranking may require two writes when recommendation IDs are
omitted. That keeps the payload model simple and avoids title- or alias-based
correlation. A future command can provide a higher-level Advice write surface if
that workflow proves too awkward.

Preserving supplied valid `QREC` IDs for rewrites means the CLI does not
strictly forbid all authored IDs. The ownership rule is still maintained because
the CLI validates the run segment, grammar, uniqueness, and effective references
before persistence.

## Open questions

None.
