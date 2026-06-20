# Wizard Mode

Use wizard as the quality wayfinder: a read-only coach that meets the user where
they are in the `QUALITY.md` lifecycle, summarizes current state, recommends the
next useful workflow, and offers concrete alternatives. Use it when the user is
unsure what to run next, asks for status/next steps, asks to review the model or
history, or sends a bare `/quality`.

Wizard must be **fast**. Its job is to route, not to audit: use the CLI status
snapshot, run the bounded
[`../guides/top-10-quality-md-checks.md`](../guides/top-10-quality-md-checks.md)
inspection when the model is valid, classify readiness, and offer concrete next
steps. Do not perform an unbounded parse of `QUALITY.md`, read evaluation report
bodies, count requirements by hand, enumerate recommendations, or resolve
build/install paths. Use `qualitymd status --json` for mechanical signals.

Run frame:

```text
/quality run
- Mode: wizard
- Target file: <resolved path>
- Scope: model/history readiness only
- Mutation: read-only
- Artifacts: none
- Next gate: choose a concrete workflow
```

## Decision Tree

```text
Probe state (CLI version plus status JSON)
- CLI missing or below the prerequisite range? recommend /quality upgrade
- QUALITY.md missing? classify no setup; recommend setup
- lint errors? classify invalid model; recommend repair
- lint valid? run top-10 QUALITY.md checks for model/lifecycle findings
- lint valid, no blocking checklist findings, no runs? classify ready to evaluate
- lint valid, no blocking checklist findings, runs present? classify evaluation history or reconciliation needs

Classify readiness (from `qualitymd status --json` plus checklist findings)
- no setup            (QUALITY.md missing)
- invalid model       (lint fails)
- starter model       (valid but skeleton/body placeholders dominate)
- immature model      (valid but model-usefulness findings block fair evaluation)
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

   These signals establish the mechanical lifecycle state.
2. If the model exists and lint passes, read
   [`../guides/top-10-quality-md-checks.md`](../guides/top-10-quality-md-checks.md)
   and inspect the target `QUALITY.md` against it. Keep findings bounded to the
   checklist. Do **not** inspect subject source files, read reports, or create
   evaluation artifacts.
3. Classify readiness from the `readiness` field, supporting status counts, and
   checklist findings. This is routing judgment, not an evaluation rating. When
   a signal is genuinely needed to choose between options and is cheap, gather
   just that one; otherwise state the open question and let the user's choice
   resolve it.
4. Only read [`../guides/authoring.md`](../guides/authoring.md) if the user
   explicitly asks for authoring/model help — not to classify readiness. If the
   user has just initialized a skeleton or asks how to start from one, read
   [`../guides/getting-started.md`](../guides/getting-started.md).
5. Report in this shape:

   ```text
   Status
   - CLI:                (version; in range or stale)
   - QUALITY.md:         (present/absent at target path)
   - Model:              (validity and usefulness/readiness; no shape audit)
   - Subject:            (ready/blocked/unknown from source signals)
   - Evaluation history: (runs/recommendations present; no body reads)
   - Readiness:

   QUALITY.md inspection findings
   - <top finding or "none blocking">

   Recommended next step
   - <workflow> because <observed reason>

   Options
   1. <concrete workflow>
   2. <concrete workflow>
   3. <concrete workflow>
   ```

6. Offer only concrete workflows the user can choose next, such as setup, model
   repair, model review/improvement, whole-subject evaluation, scoped
   target/factor evaluation, recommendation review/improvement, or evaluation
   history review. Include `/quality upgrade` when the CLI is missing, below the
   prerequisite range, or the skill/CLI pair appears incompatible. Wizard stays
   offline: it judges staleness from `qualitymd --version` against the
   prerequisite range only and does not probe the network. `/quality upgrade`
   owns the network check (`qualitymd upgrade --check`), so it is also the path
   to discover a newer-but-compatible release — surface it when the user asks
   whether a newer CLI is available.

Wizard is read-only, shallow, and status-first. It does not edit `QUALITY.md`,
create evaluation records, build reports, or rate the subject. It may produce
checklist findings about the model and judge readiness to route to work; the work
happens in the mode or confirmed workflow it hands off to.
