# Wizard Mode

Use wizard as the quality wayfinder: a read-only coach that meets the user where
they are in the `QUALITY.md` lifecycle, summarizes current state, recommends the
next useful workflow, and offers concrete alternatives. Use it when the user is
unsure what to run next, asks for status/next steps, asks to review the model or
history, or sends a bare `/quality`.

## Decision Tree

```text
Probe state
- CLI missing/stale? recommend install/upgrade
- QUALITY.md missing? classify no setup; recommend setup
- QUALITY.md present? run lint
  - lint errors? classify invalid model; recommend repair
  - lint valid? inspect model shape and evaluation history

Classify readiness
- no setup
- invalid model
- starter/skeleton model
- usable but immature model
- ready to evaluate
- has evaluation history
- mature but needs maintenance/reconciliation

Recommend one next step
- name the workflow
- explain why it is the best next step from observed state

Offer concrete alternatives
- create setup
- repair model
- review/improve QUALITY.md
- evaluate subject
- improve subject from recommendations
- review evaluation history
```

## Procedure

1. Verify the CLI prerequisite from `SKILL.md`.
2. Resolve the target file.
3. Resolve `.quality/config.yaml` and the evaluation directory from `SKILL.md`.
4. Probe state:
   - CLI readiness.
   - `QUALITY.md` presence.
   - `qualitymd lint [path]` result when a model exists.
   - model shape when lint passes: targets, factors, requirement count, and
     visible source coverage.
   - evaluation history when present: latest run, incomplete/stale-looking runs,
     active recommendations, and report/status files already available.
5. Classify readiness using the lifecycle states in the decision tree. This is a
   readiness judgment, not an evaluation rating.
6. If the model exists and the user asks for authoring/model help, or if the
   readiness classification depends on model quality, read
   [`../resources/quality-md-guide.md`](../resources/quality-md-guide.md).
7. Report in this shape:

   ```text
   Status
   - CLI:
   - QUALITY.md:
   - Model:
   - Evaluation history:
   - Readiness:

   Recommended next step
   - <workflow> because <observed reason>

   Options
   1. <concrete workflow>
   2. <concrete workflow>
   3. <concrete workflow>
   ```

8. Offer only concrete workflows the user can choose next, such as setup, model
   repair, model review/improvement, whole-subject evaluation, scoped
   target/factor evaluation, recommendation review/improvement, or evaluation
   history review.

Wizard is read-only and shallow. It does not edit `QUALITY.md`, create
evaluation records, build reports, or rate the subject. It may judge readiness
and route to work; the work happens in the mode or confirmed workflow it hands
off to.
