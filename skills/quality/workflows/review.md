---
type: Runtime Workflow
title: Review workflow
description: Runtime workflow stub for reviewing an evaluation result, QUALITY.md model, or quality concern.
---

# Review workflow

Use review to inspect what a quality signal means before deciding whether to
act. Review is read-only by default.

## Decision tree

```text
Emit run frame first

Resolve focus
- explicit evaluation/latest/report/run? focus: evaluation
- explicit model/QUALITY.md? focus: model
- named quality concern? focus: concern
- unclear? infer from lifecycle state, then ask a single-select focus question

Confirm focus
- confident? state the inferred focus and continue
- ambiguous? ask before deeper inspection

Inspect shallowly
- evaluation? inspect latest or selected reportable run when available
- model? inspect QUALITY.md using authoring and top-10-check guidance
- concern? inspect the concern against the model and available context

Close read-only
- recommend evaluate, improve, setup, update, or stop
```

## Procedure

1. Emit the run frame before tool inspection:

   ```text
   **QUALITY.md · review**
   - **Model file:** <invocation-derived path>
   - **Scope:** n/a
   - **Focus:** <evaluation | model | concern | resolving…>
   - **Mutation:** read-only
   - **Artifacts:** none
   - **Next gate:** focus confirmation, then read-only review summary
   ```

2. Resolve the model file and workspace only after the frame is emitted. If no
   model is present and the user did not provide a path, route to `setup` as the
   likely next workflow rather than inventing model findings.
3. Resolve focus from explicit user text first. Use:
   - `evaluation` for latest or selected evaluation result, report, run, rating,
     finding, evidence limit, or candidate action review;
   - `model` for `QUALITY.md` structure, usefulness, clarity, coverage,
     assessability, stale assumptions, or authoring quality;
   - `concern` for a named quality concern that is not clearly a model element or
     evaluation artifact.
4. When focus is absent or ambiguous, inspect only enough lifecycle state to make
   a recommendation: model presence, latest reportable evaluation run, and
   obvious model gaps. If confidence is still low, ask:

   ```text
   **What do you want to review?**

   1. <recommended focus> - <why this is likely>
   2. <alternate focus>
   3. <alternate focus>

   **Answer:** Reply `1`, `2`, or `3`, or name another focus.
   ```

5. For `evaluation` focus, inspect the latest reportable run when no run is
   named. Summarize available ratings, findings, evidence limits, incomplete
   inputs, and likely improve focuses. Do not create a new evaluation run.
6. For `model` focus, read `guides/authoring.md`, routed authoring sub-guides
   relevant to observed model elements, and
   `guides/top-10-quality-md-checks.md`. Report model usefulness, clarity,
   coverage, assessability, stale assumptions, and suggested improve focuses.
   Do not edit `QUALITY.md`.
7. For `concern` focus, inspect the named concern against the model and available
   project context. Recommend one next action: scoped `evaluate`, model-focused
   `improve`, work-focused `improve`, or stop.
8. Close status-first with what was reviewed, evidence limits, recommended next
   action, and the read-only boundary:

   ```text
   **Review complete**

   **Reviewed:** <evaluation result | QUALITY.md model | concern>
   **Signal:** <main judgment or routing signal>
   **Evidence limits:** <limits or none observed>
   **Recommended next:** <one workflow/action>
   **Alternatives:** <other concrete workflows/actions, when useful>
   **Not changed:** no files, evaluation records, reports, tooling, quality changelog, or external issues.
   ```

Review does not edit evaluated source, edit `QUALITY.md`, write evaluation
records, write the quality changelog, create external issues, update tooling, or create
workflow feedback logs.
