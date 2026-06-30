---
type: Runtime Workflow
title: Improve Workflow
description: Runtime workflow stub for improving work or QUALITY.md after confirming focus and mutation surface.
---

# Improve workflow

Use improve to act on quality judgment. Improve starts read-only, confirms focus
and mutation surface, then delegates to the existing safe route or stops at the
stub boundary.

## Decision tree

```text
Emit run frame first

Resolve focus
- model/QUALITY.md? focus: model
- recommendation id/artifact? focus: recommendation
- evaluation finding/candidate action/latest run? focus: evaluation
- named quality concern? focus: concern
- unclear? infer from lifecycle state, then ask a single-select focus question

Resolve mutation surface
- model focus? QUALITY.md + quality changelog when meaningful
- recommendation focus? evaluated source | QUALITY.md | quality changelog | external issue
- evaluation focus without recommendation? likely read-only until a finding/action is selected
- concern focus? evaluated source | QUALITY.md | external issue | none yet

Confirm before mutation
- confirmed model mutation? use direct model authoring route
- compatible recommendation artifact? use recommendation follow-up guide
- otherwise stop with the deferred deeper-workflow boundary
```

## Procedure

1. Emit the run frame before tool inspection:

   ```text
   **QUALITY.md · improve**
   - **Model file:** <invocation-derived path>
   - **Scope:** n/a
   - **Focus:** <evaluation | model | concern | recommendation | resolving…>
   - **Mutation:** read-only until confirmed
   - **Artifacts:** none yet
   - **Next gate:** focus and mutation-surface confirmation
   ```

2. Resolve the model file and workspace only after the frame is emitted. If no
   model is present and the user did not provide a path, route to `setup` unless
   the user named a non-model evaluated-source improvement that can proceed
   safely without QUALITY.md context.
3. Resolve focus from explicit user text first. Use:
   - `model` for improving `QUALITY.md` structure, body, Areas, Factors,
     Requirements, Rating Scale, coverage, clarity, or assessability;
   - `recommendation` for a compatible existing recommendation artifact or
     explicit recommendation ID;
   - `evaluation` for improving from a latest or selected evaluation result,
     finding, evidence limit, candidate action, or run when no compatible
     recommendation artifact has been selected;
   - `concern` for a named quality concern that is not clearly an evaluation
     artifact, recommendation artifact, or model element.
4. When focus is absent or ambiguous, inspect only enough lifecycle state to make
   a recommendation: model presence, latest reportable evaluation run, compatible
   recommendation artifacts, and obvious model gaps. If confidence is still low,
   ask:

   ```text
   **What do you want to improve?**

   1. <recommended focus> - <why this is likely>
   2. <alternate focus>
   3. <alternate focus>

   **Answer:** Reply `1`, `2`, or `3`, or name another focus.
   ```

5. Identify the likely mutation surface before acting:
   - `model`: `QUALITY.md`, and the quality changelog for meaningful model changes;
   - `recommendation`: evaluated source, `QUALITY.md`, quality changelog, or external
     issue, resolved by the recommendation-follow-up guide;
   - `evaluation`: read-only until a finding, candidate action, model gap, or
     work target is selected;
   - `concern`: evaluated source, `QUALITY.md`, external issue, or none yet.
6. Confirm both focus and mutation surface before mutating. Use the direct
   model-authoring review gate for model changes and the decision-brief shape for
   true binary mutation gates.
7. For `model` focus, read `guides/authoring.md` and the routed sub-guides
   relevant to the model elements being changed. Then follow the direct
   model-authoring route in `SKILL.md`.
8. For `recommendation` focus, read
   `guides/recommendation-follow-up.md` and follow that route. Do not treat
   recommendation follow-up as a separate public workflow.
9. For `evaluation` focus without a compatible recommendation artifact, do not
   synthesize a recommendation. Help select a finding, candidate action, model
   gap, or work target, then either route to model-focused/direct improvement,
   recommend a scoped `evaluate`, or stop with the deferred deeper-workflow
   boundary.
10. For `concern` focus, confirm whether the user wants to improve the work,
    improve `QUALITY.md`, hand off an issue, or first run/review an evaluation.
    Mutate only after the specific surface is confirmed.
11. Close status-first with changed artifacts, verification performed, remaining
    limits, and what was not changed. If the workflow stops at the stub boundary,
    say which deeper improve behavior is deferred and offer the nearest runnable
    next workflow. Use this shape:

    ```text
    **Improve status**

    **Focus:** <model | recommendation | evaluation | concern>
    **Changed:** <artifacts changed, or none>
    **Verification:** <checks run, or not run>
    **Remaining limits:** <limits or deferred behavior>
    **Not changed:** <boundaries that matter>
    **Next:** <nearest runnable workflow/action>
    ```

Improve never creates numbered evaluation records itself. If verification needs
a fresh rating, route to `evaluate` for the affected scope.
