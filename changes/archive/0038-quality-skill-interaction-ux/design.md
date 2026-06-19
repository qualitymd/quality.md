---
type: Design Doc
title: /quality skill interaction UX - design doc
description: How the /quality skill interaction contract is added to the durable spec and runtime skill prompts without changing CLI mechanics or evaluation records.
tags: [skill, quality, evaluation, ux]
timestamp: 2026-06-19T00:00:00Z
---

# /quality skill interaction UX - design doc

Design behind the
[/quality skill interaction UX](../0038-quality-skill-interaction-ux.md) change
case and its [functional spec](spec.md). This change is prompt/spec work: it
adds an interaction contract around the existing `/quality` modes without
changing CLI commands, evaluation record shape, or `QUALITY.md` format semantics.

## Context

The current `/quality` skill already has the important hard boundaries:
`wizard` is read-only, `evaluate` records only through `qualitymd evaluation ...`,
`improve` and `upgrade` require confirmation before mutation, and source under
evaluation is treated as data. What is missing is a consistent user-facing frame
around those boundaries.

The design therefore keeps the skill architecture intact:

```text
durable skill spec
        |
        v
root SKILL.md shared interaction rules
        |
        v
mode files apply the rules where they matter
        |
        v
qualitymd CLI remains the mechanical artifact writer
```

## Approach

### Durable spec section

Add a new `## User interaction contract` section to
`specs/skills/quality-skill/quality-skill.md`, after invocation/mode concepts
and before detailed workflow sections. The section becomes the cumulative source
of truth for the six interaction behaviors in the functional spec:

- run frames;
- decision briefs;
- stop rules and rerouting;
- history-aware operation;
- improvement delta reports;
- voice and status posture.

The durable spec should carry only the stable contract and rationale. It should
not include long templates for every mode; those belong in the runtime skill
files where the agent needs executable guidance.

### Root prompt as the shared UX router

Add a compact `## User Interaction Contract` section to `skills/quality/SKILL.md`.
It should define shared shapes once:

```text
Run frame fields:
- Mode
- Target file
- Scope
- Effort, when applicable
- Mutation policy
- Expected artifacts
- Next gate

Decision brief fields:
- Proposed action
- Changes
- Evidence/reason
- Recommended option
- Alternatives
- Done criterion / verification
```

Keep this section short. The root prompt is already the always-loaded contract,
so it should orient the agent without becoming a second mode procedure. Each
mode file remains responsible for when and how to apply the shared shapes.

### Mode-file integration

Apply the contract at the mode boundaries instead of scattering it through every
step.

- `wizard.md`: keep the existing status-first shape, but make readiness output
  explicitly separate CLI readiness, model validity/usefulness, subject
  readiness, and evaluation history. It should recommend one next step and then
  list concrete alternatives.
- `evaluate.md`: emit the run frame after resolving mode, target file, scope,
  effort, and config. Before creating a run, apply the stop rules. Add a short
  history-context step that inspects prior runs and recommendations when present.
- `improve.md`: use the same evaluate preflight, then present a decision brief
  before any subject or `QUALITY.md` edit. After applying an approved option,
  re-evaluate the affected scope and emit the improvement delta report.
- `setup.md`: use decision briefs only when changing an existing `QUALITY.md` or
  choosing among materially different setup paths. A first-time scaffold that the
  user explicitly requested can stay direct, provided the run frame names the
  mutation.
- `upgrade.md`: replace any bare confirmation with a decision brief that names
  whether the skill, CLI, or both are being changed, which command/tool performs
  the mutation, and how compatibility will be verified after.

### Stop-rule handling

Implement stop rules as preflight checks in the mode files, not as a new mode.
The response shape is:

```text
Stopped: <reason>
- What blocked rating:
- Why it matters:
- Best next step:
- Options:
  1. <runnable workflow>
  2. <runnable workflow>
```

This preserves `/quality wizard` as the wayfinder while letting `evaluate` and
`improve` reroute cleanly when they discover non-rateable scope.

### History context

Use existing CLI/status surfaces and evaluation files rather than adding storage.
The skill can gather:

- latest run and status;
- incomplete or stale-looking runs;
- open recommendations;
- prior rating for the same resolved scope when available.

History is summarized in the run output as context only. The current model and
fresh evidence still control the current rating. This avoids turning history into
a hidden cache of truth.

### Improvement delta report

Add the delta report as final prose output from `improve.md`; do not add a new
evaluation record field. The report should link the recommendation, applied
option, changed artifacts, evidence before/after, verification, rating movement,
and remaining limits. If later this proves useful as a machine artifact, it can
be specified separately.

## Alternatives

- **Import a gstack-style shared runtime platform.** Rejected. The useful lesson
  is interaction shape, not telemetry prompts, preference learning, routing
  injection, or model-specific prompt patches. Bringing those in would fight this
  repo's concise public skill contract.
- **Add new CLI commands for interaction UX.** Rejected. Run frames, decision
  briefs, and stop/reroute behavior are judgment/presentation work owned by the
  skill. The CLI should stay focused on deterministic mechanics.
- **Persist decision briefs or delta reports as new evaluation records.**
  Rejected for this case. The current need is human UX; changing the artifact
  schema would broaden the work and require downstream compatibility decisions.
- **Create a new `/quality plan` or `/quality interact` mode.** Rejected. The
  interaction contract should improve every existing mode rather than add another
  concept users must choose.
- **Leave all wording to mode files with no durable spec change.** Rejected. The
  behavior is cross-cutting and user-visible enough to belong in the durable
  `/quality` skill spec.

## Trade-offs & risks

- The root skill prompt gets slightly longer. Keep shared templates compact and
  put mode-specific sequencing in mode files.
- Stop rules include judgment. The design lowers risk by requiring concrete
  reasons and runnable next steps rather than vague refusal.
- History-aware output can overfit to stale runs. The contract explicitly treats
  history as context, not authority.
- Decision briefs can slow simple flows if overused. The design limits them to
  user-confirmed mutation or materially different setup choices.
- Delta reports may be verbose. Keep them scoped to the affected recommendation
  and avoid restating the full evaluation report.

## Open questions

None blocking. During implementation, wording should be tightened against the
actual mode files so the prompt remains concise and does not duplicate the
functional spec line-for-line.
