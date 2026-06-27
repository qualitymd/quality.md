---
type: Functional Specification
title: Report Gallery — functional spec
description: Delta contract for a generated example Evaluation report gallery.
tags: [evaluation, reports, examples]
timestamp: 2026-06-27T00:00:00Z
---

# Report Gallery — functional spec

Companion to the [Report Gallery](../0156-report-gallery.md). This spec states
*what* the case must do; the [design doc](design.md) covers *how*.

The key words **MUST**, **MUST NOT**, **SHOULD**, and **MAY** are to be
interpreted as described in BCP 14 when, and only when, they appear in all
capitals.

## Background / Motivation

Report layout work needs a loop measured in seconds, not a release plus a fresh
agent-authored evaluation. A generated gallery can reuse the current Evaluation
record and report-builder contracts to make report changes visible quickly while
also giving readers a browsable example of a non-trivial Evaluation report.

## Scope

Covered:

- a generated software-service report gallery example;
- deterministic regeneration of the sample `QUALITY.md`, persisted Evaluation
  payload graph, and report tree;
- freshness checking for the generated example; and
- clear disclosure that the Evaluation judgment data and evidence references are
  synthetic.

Deferred:

- a concrete fictional source system or source-code tree;
- non-software gallery examples;
- historical payload compatibility fixtures; and
- a public `qualitymd` gallery command.

## Requirements

1. The repository **MUST** include a report-gallery example under
   `examples/report-gallery/software-service/` with a sample `QUALITY.md` and a
   generated `.quality/evaluations/0001-full-eval/` Evaluation run.

   > Durable spec: none.

2. The generated gallery run **MUST** include the persisted Evaluation `data/`
   graph and generated human reports needed to browse the sample without running
   `/quality evaluate`.

   > Rationale: The example is useful as public documentation only if readers can
   > inspect the current report tree directly.
   >
   > Durable spec: none.

3. The gallery generator **MUST** validate the synthetic payload batch through
   the current Evaluation data validation path before persisting payloads and
   building reports.

   > Rationale: The gallery should fail when current payload contracts change,
   > rather than silently preserving stale hand-authored JSON.
   >
   > Durable spec: none.

4. The report-gallery check **MUST** fail when regenerating the gallery changes
   files under `examples/report-gallery/`.

   > Rationale: CI needs a direct stale-output signal when report rendering,
   > sample model content, or payload contracts change.
   >
   > Durable spec: none.

5. The example documentation **MUST** state that the gallery uses synthetic
   Evaluation routine outputs and intentionally omits a concrete source
   system/code tree.

   > Rationale: The sample should demonstrate report structure without implying
   > that the synthetic judgment is an actual product assessment.
   >
   > Durable spec: none.

6. The gallery example **MUST NOT** be framed as a default QUALITY.md model or
   default software-quality factor set.

   > Rationale: QUALITY.md is domain agnostic; a software-service gallery is one
   > illustrative example, not the default modeled domain.
   >
   > Durable spec: none.

## Durable spec changes

### To add

None

### To modify

None

### To rename

None

### To delete

None
