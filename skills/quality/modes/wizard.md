# Wizard Mode

Use wizard as the quality wayfinder: a read-only coach that meets the user where
they are in the `QUALITY.md` lifecycle, summarizes current state, recommends the
next useful workflow, and offers concrete alternatives. Use it when the user is
unsure what to run next, asks for status/next steps, asks to review the model or
history, or sends a bare `/quality`.

Wizard must be **fast**. Its job is to route, not to audit: probe a handful of
cheap signals in one shot, classify readiness, and offer concrete next steps.
Do not hand-parse `QUALITY.md` frontmatter, read evaluation report bodies, count
requirements, enumerate recommendations, or resolve build/install paths — all of
that is deferred to the mode the user picks. Aim for one batched probe and a fast
answer.

## Decision Tree

```text
Probe state (one batched probe of cheap signals only)
- CLI missing/stale? recommend install/upgrade
- QUALITY.md missing? classify no setup; recommend setup
- QUALITY.md present? run lint
  - lint errors? classify invalid model; recommend repair
  - lint valid? check only whether the evaluation dir has runs

Classify readiness (from the cheap signals only)
- no setup            (QUALITY.md missing)
- invalid model       (lint fails)
- ready to evaluate   (lint passes, no runs yet)
- has evaluation history (lint passes, runs present)
Finer maturity (skeleton vs. mature, stale runs, reconciliation needs) is judged
in the mode the user hands off to, not eagerly here.

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

1. Run **one batched probe** of cheap signals — prefer a single shell call that
   collects all four at once:
   - CLI version (`qualitymd --version`) and whether it satisfies the prerequisite
     range from `SKILL.md`.
   - `QUALITY.md` presence at the resolved target path (do not walk parents).
   - `qualitymd lint [path]` pass/fail when the file exists.
   - whether the evaluation directory (from `.quality/config.yaml`, default
     `quality/evaluations/`) contains any run folders.

   These four signals are sufficient to route. Do **not** open the model,
   reports, or recommendation files to classify — defer that to the chosen mode.
2. Classify readiness using the lifecycle states in the decision tree. This is a
   readiness judgment from the probe signals, not an evaluation rating. When a
   signal is genuinely needed to choose between options and is cheap, gather just
   that one; otherwise state the open question and let the user's choice resolve
   it.
3. Only read [`../resources/quality-md-guide.md`](../resources/quality-md-guide.md)
   if the user explicitly asks for authoring/model help — not to classify
   readiness.
4. Report in this shape:

   ```text
   Status
   - CLI:                (version; in range or stale)
   - QUALITY.md:         (present/absent at target path)
   - Model:              (lint pass/fail — not a shape breakdown)
   - Evaluation history: (runs present/absent; no body reads)
   - Readiness:

   Recommended next step
   - <workflow> because <observed reason>

   Options
   1. <concrete workflow>
   2. <concrete workflow>
   3. <concrete workflow>
   ```

5. Offer only concrete workflows the user can choose next, such as setup, model
   repair, model review/improvement, whole-subject evaluation, scoped
   target/factor evaluation, recommendation review/improvement, or evaluation
   history review.

Wizard is read-only and shallow. It does not edit `QUALITY.md`, create
evaluation records, build reports, or rate the subject. It may judge readiness
and route to work; the work happens in the mode or confirmed workflow it hands
off to.
