# Wizard Mode

Use wizard as the quality wayfinder: a read-only coach that meets the user where
they are in the `QUALITY.md` lifecycle, summarizes current state, recommends the
next useful workflow, and offers concrete alternatives. Use it when the user is
unsure what to run next, asks for status/next steps, asks to review the model or
history, or sends a bare `/quality`.

Wizard must be **fast**. Its job is to route, not to audit: use the CLI status
snapshot, classify readiness, and offer concrete next steps. Do not hand-parse
`QUALITY.md` frontmatter, read evaluation report bodies, count requirements,
enumerate recommendations, or resolve build/install paths. Use
`qualitymd status --json` for those mechanical signals.

## Decision Tree

```text
Probe state (CLI version plus status JSON)
- CLI missing/stale? recommend /quality upgrade
- QUALITY.md missing? classify no setup; recommend setup
- lint errors? classify invalid model; recommend repair
- lint valid, no runs? classify ready to evaluate
- lint valid, runs present? classify evaluation history or reconciliation needs

Classify readiness (from `qualitymd status --json`)
- no setup            (QUALITY.md missing)
- invalid model       (lint fails)
- ready to evaluate   (lint passes, no runs yet)
- has evaluation history (lint passes, runs present)
- needs reconciliation (stale/incomplete/malformed runs or active recommendations)

Recommend one next step
- name the workflow
- explain why it is the best next step from observed state

Offer concrete alternatives
- upgrade skill/CLI pair
- create setup
- repair model
- review/improve QUALITY.md
- evaluate subject
- improve subject from recommendations
- review evaluation history
```

## Procedure

1. Run a shallow CLI probe:
   - `qualitymd --version` and whether it satisfies the prerequisite range from
     `SKILL.md`.
   - `qualitymd status --json [path]` for the resolved target path (do not walk
     parents).

   These signals are sufficient to route. Do **not** open the model, reports, or
   recommendation files to classify — defer that to the chosen mode.
2. Classify readiness from the `readiness` field and supporting status counts.
   This is a readiness judgment from mechanical signals, not an evaluation
   rating. When a signal is genuinely needed to choose between options and is
   cheap, gather just that one; otherwise state the open question and let the
   user's choice resolve it.
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
   history review. Include `/quality upgrade` when the CLI is missing, stale, or
   outside the prerequisite range, or when the skill/CLI pair appears
   incompatible.

Wizard is read-only and shallow. It does not edit `QUALITY.md`, create
evaluation records, build reports, or rate the subject. It may judge readiness
and route to work; the work happens in the mode or confirmed workflow it hands
off to.
