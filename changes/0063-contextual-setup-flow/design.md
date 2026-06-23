---
type: Design Doc
title: Contextual setup flow - design doc
description: How /quality setup becomes a context-informed discovery flow that writes only QUALITY.md.
tags: [skill, setup, ux]
timestamp: 2026-06-23T00:00:00Z
---

# Contextual setup flow - design doc

## Context

This design answers the
[Contextual setup flow functional spec](spec.md). The change keeps setup as a
skill-driven workflow: the `qualitymd` CLI still owns deterministic
scaffolding, linting, status, and format grounding, while the `/quality` skill
owns judgment, elicitation, model authoring, and routing.

The main design constraint is the new setup mutation boundary. Setup writes
only the selected `QUALITY.md`; evaluation artifacts, quality-log entries,
external issues, recurring-review automation, CI, release gates, Codex
automations, and Claude Code routines are all follow-on workflows.

## Approach

Treat setup as a four-step runtime procedure inside
`skills/quality/modes/setup.md`:

```text
1. Preflight and context analysis
2. Discovery questions
3. QUALITY.md authoring
4. Validation, readiness inspection, and next-step routing
```

### Preflight and context analysis

Keep the current CLI prerequisite check from `SKILL.md`: setup still stops when
`qualitymd` is missing, stale, or outside the supported range. Resolve the model
file the same way the skill does today: explicit path when supplied, otherwise
`QUALITY.md` in the current working directory.

Before asking the user questions, inspect available repository context for
setup signals:

- README, install docs, contributor docs, and user-facing guides;
- package metadata, test layout, build scripts, release/version files, and
  workflow files;
- existing agent instructions such as `AGENTS.md` and skill docs;
- repository structure and obvious component boundaries;
- visible work-management or recurring-review signals.

This is not an evaluation pass. The read is broad enough to infer setup defaults
and missing context, not to rate the project.

Represent the analysis internally as a short setup fact sheet:

```text
Root area: <path/title/default>
Lifecycle: <default> (<confidence>, evidence)
Risk tolerance: <default> (<confidence>, evidence)
Modeling rigor: <default> (<confidence>, evidence)
Collaboration: <default> (<confidence>, evidence)
Needs: <candidate user/maintainer/other needs>
Missing context: <specific gaps>
Workflow posture: <work management / recurring review hints>
Candidate model shape: <factors and child areas>
```

The fact sheet is not a new file. It is working context used to ask better
questions and author `QUALITY.md`.

### Discovery questions

Ask one compact setup prompt, or a short sequence if the interaction surface
handles that better. Each question should show the recommended default and its
confidence. Use a small fixed confidence vocabulary:

- `strongly inferred`
- `weakly inferred`
- `assumed`

The question set stays stable:

1. Lifecycle.
2. Risk tolerance.
3. Modeling rigor.
4. Collaboration context, with agent-heavy development assumed.
5. Needs map and missing context.

Work-management and recurring-review posture can be included as optional
questions when repo context suggests they are relevant or when the user has
already brought them up. Frame them as context capture. Do not ask for permission
to create issues or automations.

The needs question should be seeded, not blank. For example:

```text
From the repo, I think user context and support workflow are missing. What else
important is not visible here?
```

Do not ask users to design factors and areas cold. Use discovery answers plus
repo context to propose the model shape during authoring.

### QUALITY.md authoring

After discovery, synthesize directly into `QUALITY.md`.

For a missing model, run `qualitymd init [path]` after discovery and before
authoring content. That preserves CLI ownership of the starter scaffold while
letting the discovery answers shape the first useful model immediately.

For an existing model, present the required decision brief before editing. The
brief should summarize the intended model changes, the evidence/reason from
context analysis and discovery, and the verification step. This is the only
write gate; do not add a separate alignment phase.

Author the body first, then the frontmatter model:

1. Overview and Scope establish the root area, default current-directory scope,
   lifecycle, risk tolerance, modeling rigor, and key boundaries.
2. Needs and Risks capture user needs, maintainer/collaborator needs, other
   stakeholder needs, and the failure modes that matter.
3. Unknowns and open questions capture missing or non-agent-accessible context.
4. The rating scale stays standard unless the body shows a real mismatch.
5. Factors and child areas derive from the needs/risks and obvious component
   boundaries.
6. Requirements are small, concrete, and assessable from agent-accessible
   evidence, or explicitly name missing evidence.
7. Include the `quality-md` self-check area when the existing setup-area
   contract says it is appropriate.

Record work-management and recurring-review posture in the body only when it
matters to the quality loop. Keep it descriptive: where recommendations may be
handed off later, whether maintainers expect a cadence, and what remains
undecided. Do not create config files or automation artifacts.

### Validation and readiness

After writing, run `qualitymd lint [path]`. If lint fails, report the lint
failure and route to continued `QUALITY.md` iteration. Do not offer evaluation
as the recommended next step while the model is invalid.

When lint passes, inspect the resulting model with the Top 10 QUALITY.md checks.
This is a bounded model-readiness inspection, not a project evaluation. Use it
to classify the model as `starter`, `immature`, or `ready to evaluate` and to
surface model-usefulness gaps.

The completion response should be status-first:

```text
Setup complete
- Changed: QUALITY.md
- Validation: lint passed | lint failed
- Readiness: starter | immature | ready to evaluate
- Not done: no evaluation, no issues, no automations
- Next: continue iterating | run evaluation | set up recurring review | set up recommendation handoff | stop here
```

When readiness is not `ready to evaluate`, list the most important model gaps
and make continued iteration the recommended next step.

## Alternatives

**Keep skeleton-first guided population.** Rejected because it makes setup
mechanically valid before it is contextually useful. The new flow asks for the
human judgments first, then uses the scaffold as mechanics.

**Add a public `/quality setup wizard` or revive `/quality wizard`.** Rejected
because the public surface should remain task-oriented. The setup interaction
can be wizard-like without adding a new mode.

**Ask users to choose quality factors and child areas directly.** Rejected
because it asks users to design the schema before seeing the synthesized
project context. Setup should ask about lifecycle, risk, collaboration, and
needs, then propose factors and areas from those answers.

**Use an explicit Alignment phase before writing.** Rejected for this iteration
because it slows setup and creates a second approval moment. Existing-model
edits still use the required decision brief; missing-model setup can proceed
from discovery to writing.

**Seed an inaugural quality-log entry during setup.** Rejected because setup's
mutation boundary is intentionally narrow. The first model's rationale should
live in `QUALITY.md`; the quality log should resume with later confirmed
model-authoring or recommendation-apply changes.

**Recommend CI or release gating as the default recurring loop.** Rejected
because the intended loop is a maintainer-chosen cadence that fits how the team
works. CI/release gating can be requested later, but it is not the default setup
recommendation.

## Trade-offs & Risks

Writing only `QUALITY.md` simplifies setup and makes side effects clear, but it
removes the separate quality-log record of initial model creation. The mitigation
is to put the durable setup rationale directly in the body: assumptions, needs,
risks, missing context, and review state.

The context analysis step could become too broad. The runtime guide should keep
it bounded to setup signals and avoid source-quality judgment. If an agent finds
itself collecting findings about the project implementation, it has crossed into
evaluation and should stop.

Default confidence labels are useful only if they stay terse and evidence-led.
If every default becomes a paragraph, setup stops being fast. Prefer one short
evidence phrase per default.

The no-alignment-gate approach is faster, but it puts more pressure on the
quality of the discovery prompt and final summary. Existing models still have a
decision brief before mutation; new models rely on the user's choice to run
`/quality setup` plus the discovery answers.

Removing setup from the quality-log writer contract touches multiple durable
specs and docs. That is an intentional simplification, but implementation needs
to cleanly reconcile all references to inaugural setup log entries.

## Open Questions

- Should the recurring-review posture be represented in a recommended body
  section, or simply captured wherever it naturally fits in Overview, Scope,
  Needs, Risks, or open questions?
- Should setup always ask optional work-management and recurring-review
  questions, or ask them only when repo context or user input suggests they
  matter?
