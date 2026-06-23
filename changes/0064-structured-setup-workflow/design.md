---
type: Design Doc
title: Structured setup workflow - design doc
description: How /quality setup becomes an operator workflow with concrete discovery prompts.
tags: [skill, setup, ux]
timestamp: 2026-06-23T00:00:00Z
---

# Structured setup workflow - design doc

## Context

This design answers the
[Structured setup workflow functional spec](spec.md). The change is a guidance
rewrite, not a new CLI feature: `qualitymd` still owns deterministic mechanics
such as `init` and linting, while the `/quality` skill owns context analysis,
elicitation, model authoring, readiness inspection, and next-step routing.

The key design move is to make the runtime setup document executable by an
agent. Instead of mostly describing setup's constraints, it should give an
ordered workflow, a setup brief template, concrete prompt text, branching rules,
and completion output.

## Approach

Rewrite `skills/quality/modes/setup.md` around the workflow the agent runs:

```text
1. Preflight
2. Context read
3. Setup brief
4. Discovery prompt
5. Scaffold when missing
6. Author QUALITY.md
7. Validate and inspect readiness
8. Report completion
```

Keep the file path and dispatch wiring unchanged. The heading and prose should
say "Setup Workflow"; the parent `SKILL.md` can still use "mode" where it
describes argument parsing and dispatch.

### Runtime structure

Make the runtime setup guide a playbook with short imperative sections:

- **Preflight** — verify CLI support, resolve the model file, emit the run frame.
- **Read context** — inspect bounded setup signals and avoid source-quality
  evaluation.
- **Build the setup brief** — record defaults, confidence, and evidence.
- **Ask discovery questions** — use the compact prompt by default; split into a
  short sequence when needed.
- **Write the model** — run `qualitymd init` for missing files, use a decision
  brief before changing an existing file, then synthesize body and frontmatter.
- **Verify** — run lint and Top 10 readiness inspection.
- **Close** — report changed artifact, validation, readiness, important gaps,
  and next-step choices.

This is still text guidance, not a generated script. The purpose is to make
agent behavior repeatable without pretending setup is deterministic.

### Setup brief

Represent the brief as a fenced template in the runtime guide:

```text
Root area: <default> (<confidence>, <evidence>)
Domain: <default> (<confidence>, <evidence>)
Lifecycle: <default> (<confidence>, <evidence>)
Risk tolerance: <default> (<confidence>, <evidence>)
Modeling rigor: <default> (<confidence>, <evidence>)
Collaboration: <default> (<confidence>, <evidence>)
Primary users/outcomes: <default> (<confidence>, <evidence>)
Maintainers/collaborators: <default> (<confidence>, <evidence>)
Other stakeholders: <default> (<confidence>, <evidence>)
Missing context: <specific gaps>
Review posture: <default when visible> (<confidence>, <evidence>)
Candidate model shape: <factors and child areas>
```

The brief stays working context. It is shown to the user only through the
discovery prompt, then distilled into `QUALITY.md` where it affects the model.

### Discovery prompt

Default to one compact prompt that asks the user to confirm or correct inferred
assumptions. This keeps setup fast while still surfacing every required
question:

```text
I inspected the repo for setup signals. Please confirm or correct these
assumptions before I write QUALITY.md:

1. Root area: <default>. Recommended: <answer> (<confidence>).
2. Domain: <default>. Recommended: <answer> (<confidence>).
3. Lifecycle: <options>. Recommended: <answer> (<confidence>).
4. Risk tolerance: <options>. Recommended: <answer> (<confidence>).
5. Modeling rigor: <options>. Recommended: <answer> (<confidence>).
6. Primary users/outcomes: <answer> (<confidence>).
7. Maintainers/collaborators: <answer> (<confidence>).
8. Other stakeholders: <answer> (<confidence>).
9. Missing context: I think <specific gaps> are not visible. What else should
   the model record as unknown or not agent-accessible?
10. Review posture: <options>. Recommended: <answer> (<confidence>).
```

Use a short sequence when the compact prompt would be too large or the
interaction surface benefits from grouping. The recommended grouping is:

1. Root area and domain.
2. Lifecycle, risk tolerance, and modeling rigor.
3. Users, maintainers/collaborators, and other stakeholders.
4. Missing context and review posture.

The runtime guide should state that the agent may accept terse corrections. The
user does not need to answer every item in prose if the defaults are right.

### Model authoring

Keep the synthesis approach from the contextual setup flow: body first, then
frontmatter. The new design sharpens the input contract rather than changing the
authoring output.

Map answers into `QUALITY.md` like this:

- root area and domain shape Overview and Scope;
- lifecycle, risk tolerance, and modeling rigor calibrate the rating bar and
  requirement depth;
- user, maintainer, collaborator, and stakeholder answers shape Needs and Risks;
- missing context becomes Unknowns or Open questions in the relevant sections;
- review posture is recorded only when it affects how the model should be used;
- candidate factors and child Areas are derived after discovery, not chosen by
  the user directly.

### Durable specs and public docs

Update durable setup specs to describe the concrete workflow and question set.
Keep public docs lightweight: they should say setup runs a guided setup
workflow with recommended defaults, not reproduce all ten questions.

Use "mode" only where the document is discussing dispatch or internal routing.
Use "workflow" where the document is describing what the user or agent does.

## Alternatives

**Rename `modes/` to `workflows/`.** Rejected for this change. The terminology
problem is in prose and operator guidance, not in the file path. Renaming now
would create link churn and distract from the prompt contract.

**Keep five abstract discovery categories.** Rejected because it leaves too
much discretion in the runtime guide. The earlier categories were right, but
setup needs concrete questions so different agents ask the same thing.

**Ask one question at a time.** Rejected as the default because setup should be
quick and context-informed. A compact confirmation prompt makes defaults visible
and lets the user correct only what is wrong. A short sequence remains available
when needed.

**Make review posture optional only.** Rejected for the workflow prompt. Review
posture is cheap to ask and helps record how the model should be used, as long
as it is framed as context capture rather than automation setup.

**Ask for factors and Areas directly.** Rejected because it pushes model design
onto the user. The agent should infer model shape from needs, risks,
stakeholders, and evidence.

## Trade-offs & Risks

A ten-item prompt can feel long. The mitigation is to present inferred defaults
and ask for corrections, not full answers. The short-sequence fallback handles
surfaces where one compact prompt is unwieldy.

Concrete prompt text improves consistency but may become stale if setup
concepts change. Keep the durable spec as the contract and the runtime guide as
the operational rendering of that contract.

Review posture could be mistaken for automation setup. The prompt and closeout
must keep saying setup records posture only; recurring review, issue handoff,
and automations are follow-on workflows.

Leaving `modes/` paths in place means terminology will not be perfectly pure.
That is acceptable because the user-facing and runtime prose carry the
workflow concept, while dispatch internals remain stable.

## Open Questions

- Should the compact prompt include the full evidence note for every default, or
  keep evidence in the setup brief and show only confidence in the prompt?
- Should the work-handoff question be part of the default compact prompt, or
  remain an optional add-on only when repo context suggests an issue tracker or
  handoff process?
