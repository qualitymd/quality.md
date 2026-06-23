---
type: Functional Specification
title: Structured setup workflow
description: Turn /quality setup guidance into an explicit workflow with concrete discovery questions and prompt framing.
tags: [skill, setup, ux]
timestamp: 2026-06-23T00:00:00Z
---

# Structured setup workflow

This Change Case spec defines the delta for `/quality setup`: keep the
context-informed setup contract, but express it as a concrete setup workflow
with a stable discovery question set, explicit prompt framing, and clear
operator steps.

This document uses BCP 14 keywords only for testable conformance requirements.
The key words "MUST", "MUST NOT", "SHOULD", and "MAY" are to be interpreted as
described in [RFC 2119](../../docs/reference/rfc2119.md) and
[RFC 8174](../../docs/reference/rfc8174.md) when, and only when, they appear in
all capitals.

## Background / Motivation

The previous contextual setup work identified the right setup judgments but left
the runtime guidance too close to a spec. An agent running setup needs a
repeatable workflow: what to inspect, what brief to build, exactly what to ask,
how to frame defaults, what to synthesize, how to validate, and how to stop.

The user-facing concept is also better described as a workflow than a mode.
"Mode" remains useful for internal dispatch across `setup`, `evaluate`, and
`update`, but it should not be the dominant way the setup process explains
itself.

## Scope

Covered:

- Runtime `/quality setup` workflow structure.
- Concrete setup discovery questions and option sets.
- Recommended-default and confidence framing for every setup question.
- Grouped-prompt and short-sequence interaction forms.
- Root area, domain, lifecycle, risk tolerance, modeling rigor, collaboration,
  stakeholder needs, missing context, and review posture as explicit setup
  inputs.
- Public/runtime terminology that calls setup a workflow while preserving
  internal dispatch vocabulary where needed.

Deferred / non-goals:

- No QUALITY.md format change.
- No CLI-owned setup wizard or interactive CLI workflow.
- No evaluation, recommendation handoff, quality-log writing, issue creation,
  or automation configuration during setup.
- No filesystem rename from `modes/` to `workflows/` in this case unless a later
  design explicitly decides the rename is worth the compatibility and link
  churn.
- No user-facing requirement to choose factors, Areas, Requirements, or rating
  criteria directly.

## Requirements

### Terminology

User-facing setup guidance **MUST** describe `/quality setup` as starting or
running a setup workflow.

Runtime guidance **SHOULD** use "workflow" for the operator procedure and
reserve "mode" for internal dispatch, routing, or compatibility language.

Existing `modes/` paths **MAY** remain unchanged when the durable and runtime
text makes the workflow terminology clear.

> Rationale: the user chooses a task, not an implementation state. Keeping path
> names stable avoids churn while still making the setup procedure easier to
> understand. — 0064

### Workflow Shape

The runtime setup guidance **MUST** read as an operator playbook with ordered
steps, not only as conformance requirements.

The setup workflow **MUST** include these stages, in order:

1. Resolve the target `QUALITY.md` and verify setup prerequisites.
2. Inspect repository context for setup signals.
3. Build a setup brief with inferred defaults, confidence, and evidence.
4. Ask the concrete discovery questions.
5. Run `qualitymd init [path]` when the target model is missing.
6. Synthesize or update `QUALITY.md`.
7. Run lint and inspect model readiness.
8. Report completion and next-step choices.

The workflow **MUST NOT** ask the user to design factors, child Areas,
Requirements, or rating levels cold. The agent derives model shape from the
setup brief, discovery answers, authoring guide, and repository context.

### Setup Brief

Before asking discovery questions, the workflow **MUST** build a concise setup
brief containing:

- root area;
- domain;
- lifecycle;
- risk tolerance;
- modeling rigor;
- collaboration context;
- inferred primary users and outcomes;
- inferred maintainer or collaborator needs;
- inferred other stakeholder needs;
- missing or non-agent-accessible context;
- review posture when visible; and
- candidate model shape.

Every inferred setup brief item **MUST** include a recommended default,
confidence signal, and short evidence note when evidence exists.

The confidence vocabulary **MUST** be fixed to:

- `strongly inferred`;
- `weakly inferred`; and
- `assumed`.

### Discovery Questions

The setup workflow **MUST** ask or present the following discovery questions
before writing `QUALITY.md`, unless the user explicitly asks to accept all
inferred defaults:

1. **Root area.** Should this `QUALITY.md` model the whole current project, or a
   narrower area?
2. **Domain.** What kind of thing is this model evaluating?
3. **Lifecycle.** Which stage best fits: exploratory, pre-release, active
   production, maintenance, or sunset?
4. **Risk tolerance.** How costly is poor quality here: high tolerance,
   moderate tolerance, or low tolerance?
5. **Modeling rigor.** How detailed should the first quality model be:
   lightweight, standard, or high-assurance?
6. **Primary users and outcomes.** Who needs the evaluated thing to work, and
   what outcomes matter most?
7. **Maintainers and collaborators.** Who has to change, operate, review, or
   rely on this work?
8. **Other stakeholders.** Are there customers, operators, compliance, support,
   data, security, business, or other stakeholders not visible in the repo?
9. **Missing context.** The agent thinks these important inputs are not visible:
   `<specific gaps>`. What else should the model record as unknown or not
   agent-accessible?
10. **Review posture.** Should the model record a recurring review expectation:
    none for now, per sprint or iteration, monthly, before major releases or
    planning, custom, or another cadence?

Each question **MUST** include a recommended answer and confidence signal.

The missing-context question **MUST** be seeded from repository analysis rather
than phrased as a blank "anything else?" prompt.

The collaboration question **MUST** assume agent-heavy development and ask which
human collaborators, reviewers, maintainers, or stakeholders also need to align
with the quality bar.

The review-posture question **MUST** be framed as context capture, not as
permission to create automations, CI gates, release gates, calendar events, or
issue-tracker artifacts.

Setup **MAY** ask an additional work-handoff question about where future
evaluation recommendations should usually go. If asked, it **MUST** say setup
will not create issues or configure integrations.

### Prompt Form

The discovery questions **MAY** be presented as one compact prompt or as a short
sequence when the interaction surface works better that way.

A compact prompt **MUST** ask the user to confirm or correct the setup
assumptions instead of requiring full prose answers to every question.

A compact prompt **MUST** preserve all required questions, defaults, confidence
signals, and seeded missing-context content.

A short sequence **SHOULD** group closely related questions together, such as
root area plus domain, lifecycle plus risk plus rigor, stakeholders plus needs,
and missing context plus review posture.

### Model Synthesis

Setup-authored `QUALITY.md` body content **MUST** preserve the confirmed or
corrected setup assumptions where they shape the model.

The generated model **MUST** derive factors and child Areas from project needs,
risks, stakeholder concerns, component boundaries, and available evidence rather
than from the question labels alone.

The generated model **MUST** record missing or non-agent-accessible context as
unknowns or open questions in the relevant body sections.

### Completion

Setup completion output **MUST** report:

- changed artifact;
- lint result;
- model-readiness classification;
- important remaining model gaps; and
- next-step choices.

Setup completion **MUST NOT** automatically run evaluation, write the quality
log, create issues, or configure recurring review.

## Acceptance Criteria

- Runtime setup guidance is organized as an ordered workflow.
- Runtime setup guidance includes a concrete setup brief template.
- Runtime setup guidance includes all ten discovery questions.
- Each question has required framing for recommended defaults and confidence.
- Root area and domain are first-class setup questions.
- Lifecycle, risk tolerance, and modeling rigor remain separate questions.
- Missing context is seeded from analysis.
- Review posture is framed as context capture, not automation setup.
- Users are not asked to design factors, Areas, Requirements, or rating criteria
  directly.
- Public/runtime prose uses setup workflow language where user-facing.
- Internal dispatch can continue to use mode terminology.
- Existing `modes/` file paths are not renamed by default.
- Setup's mutation boundary remains `QUALITY.md` only.

## Durable spec changes

### To add

None.

### To modify

- `specs/skills/quality-skill/quality-skill.md` - align setup summary and
  routing language with workflow terminology while preserving internal dispatch
  semantics.
- `specs/skills/quality-skill/modes/setup.md` - replace abstract setup-mode
  wording with the structured workflow, setup brief, concrete question set,
  prompt forms, and completion contract specified above.
- `specs/skills/quality-skill/guides/getting-started-md.md` - align post-setup
  iteration guidance with the setup assumptions and review posture captured by
  the structured workflow.
- Relevant OKF logs and indexes under `specs/` - record durable spec updates
  when they land.

### To rename

None.

### To delete

None.
