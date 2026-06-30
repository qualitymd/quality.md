---
type: Design Doc
title: Report Gallery — design doc
description: Design for the generated Evaluation report gallery.
tags: [evaluation, reports, examples]
timestamp: 2026-06-27T00:00:00Z
---

# Report Gallery — design doc

Design behind the [Report Gallery](../0156-report-gallery.md) and its
[functional spec](spec.md).

## Context

The report builder already renders deterministic Markdown from persisted
Evaluation data. The gallery uses that separation directly: generate a realistic
current-contract run once, check in the resulting payload graph and reports, and
make regeneration cheap enough to use during report design.

## Approach

Add a Go generator under `scripts/report-gallery/` that rewrites
`examples/report-gallery/software-service/` from deterministic fixture content:

1. write the example `README.md` and `QUALITY.md`;
2. remove any previous `.quality/evaluations/` tree;
3. create a fresh full Evaluation run through `evaluation.CreateRun`;
4. assemble a synthetic payload batch for a realistic software-service model;
5. call `evaluation.SetData` with `DryRun: true`;
6. call `evaluation.SetData` again to persist the payloads; and
7. call `evaluation.BuildReport` to write the report tree.

The sample model uses a fictional software service with multiple Areas,
cross-cutting Factors, mixed Finding types, ranked findings, and ranked
recommendations. Evidence references use a `synthetic-source:` prefix because
the concrete source system is intentionally deferred.

Two `mise` tasks wrap the generator:

- `report-gallery` regenerates the checked-in example.
- `report-gallery-check` regenerates it and runs `git diff --exit-code
examples/report-gallery`.

`mise run check` depends on `report-gallery-check` so stale examples fail in the
normal local/CI gate.

## Spec response

The generated example lives in `examples/report-gallery/software-service/` and
uses the normal `.quality/evaluations/0001-full-eval/` workspace layout. Because
the generator creates the run and writes data through the evaluation package,
schema drift fails during generation rather than after publication. The checked
task turns any regenerated file diff into a failing gate. The README carries the
synthetic-data disclosure and frames the model as illustrative.

## Alternatives

**Check in hand-authored JSON payloads.** Rejected because stale payloads would
survive data model changes until a report build or reader noticed.

**Keep the gallery under `internal/evaluation/testdata`.** Rejected because the
example should be browsable as documentation, not only as a test fixture.

**Add a public `qualitymd report-gallery` command.** Rejected because this is a
repository development and documentation aid, not a user-facing CLI capability.

**Put reports outside `.quality/` for discoverability.** Rejected because the
example should match normal QUALITY.md workspace behavior. The README can link
directly to the hidden-folder reports.

## Trade-offs & risks

The generator is deterministic but still produces a large checked-in example
tree. That is acceptable because public report examples need to be browsable
without running tooling. The synthetic evidence references are less satisfying
than real source files, so the README calls out the omission explicitly.

Wiring the check into `mise run check` makes report changes carry example diffs.
That is intentional, but it means contributors changing reports must regenerate
the gallery as part of normal verification.

## Open questions

- Which non-software example should pair with this software-service gallery
  later?
- When the gallery grows, should edge-case report fixtures stay public or move
  into a test-only fixture set?
