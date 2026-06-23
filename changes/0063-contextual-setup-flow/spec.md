---
type: Functional Specification
title: Contextual setup flow
description: Rework /quality setup into a short context-informed discovery flow that writes only QUALITY.md and routes next steps.
tags: [skill, setup, ux]
timestamp: 2026-06-23T00:00:00Z
---

# Contextual setup flow

This Change Case spec defines the delta for `/quality setup`: make setup a
short, context-informed discovery flow that writes a useful first `QUALITY.md`
and then offers next-step choices without running evaluation, configuring
automation, creating issues, or writing other workspace artifacts.

This document uses BCP 14 keywords only for testable conformance requirements.
The key words "MUST", "MUST NOT", "SHOULD", and "MAY" are to be interpreted as
described in [RFC 2119](../../docs/reference/rfc2119.md) and
[RFC 8174](../../docs/reference/rfc8174.md) when, and only when, they appear in
all capitals.

## Background / Motivation

Setup is where the user first turns project context into an explicit quality
bar. A good setup experience should be fast and pedagogical without explaining
the model in long prose: the questions themselves should illuminate what matters
for QUALITY.md. The repo can supply many clues, but it cannot reliably know the
project's lifecycle, risk tolerance, collaboration model, stakeholder needs, or
missing context. Those are the judgments setup needs from the user.

Setup should not blur into the rest of the quality loop. Evaluation, recurring
quality review, recommendation handoff, and automation setup are important
follow-on workflows, but setup's mutation boundary should stay narrow: create or
update `QUALITY.md`, validate it, and route the next step.

## Scope

Covered:

- Runtime `/quality setup` behavior.
- Missing and existing `QUALITY.md` setup paths.
- Context analysis and discovery questions.
- Model authoring expectations for the generated `QUALITY.md`.
- Setup validation and model-readiness inspection.
- Setup next-step options.
- Runtime skill guidance, durable skill specs, public docs, and quality-log
  writer contracts affected by setup's mutation boundary.

Deferred / non-goals:

- No QUALITY.md format change.
- No CLI-owned interactive setup wizard.
- No evaluation run, recommendation record, issue creation, or automation
  configuration during setup.
- No recommendation that CI or release gating is the default recurring quality
  loop.
- No design of recurring review automation itself; setup only records posture
  and offers the follow-on workflow.

## Requirements

### Context Analysis

`/quality setup` **MUST** inspect available repository context before asking
setup questions. Relevant context includes README and docs, repository
structure, package metadata, tests, contributor docs, existing agent
instructions, and visible workflow or work-management signals.

Setup **MUST** treat the current directory as the default root Area convention
unless the user supplied an explicit model path or repository context strongly
indicates a narrower root.

Setup **MUST** present inferred defaults with a confidence signal. The confidence
vocabulary **SHOULD** distinguish at least strongly inferred, weakly inferred,
and assumed defaults.

> Rationale: defaults are useful only when the user can see whether they came
> from evidence or from a fallback assumption. — 0063

### Discovery Questions

Setup **MUST** ask a small set of high-leverage discovery questions before
writing `QUALITY.md`, unless the user explicitly asks to accept all inferred
defaults.

The discovery questions **MUST** cover:

- lifecycle: exploratory, pre-release, active production, maintenance, or
  sunset;
- risk tolerance: high, moderate, or low;
- modeling rigor: lightweight, standard, or high-assurance;
- collaboration context, assuming agent-heavy development and asking who else
  must align with the quality bar; and
- needs map: primary user needs, collaborator/maintainer needs, other
  stakeholder needs, and missing context.

Setup **MUST** keep lifecycle, risk tolerance, and modeling rigor separate.

> Rationale: a production internal tool can need lightweight modeling, and a
> pre-release library can need high-assurance modeling. Collapsing the concepts
> makes the quality bar harder to calibrate. — 0063

The missing-context prompt **MUST** be seeded from analysis, for example
identifying likely missing customer context, release expectations, support
workflow, private operational constraints, or inaccessible evidence before
asking what else is missing.

Setup **MAY** ask optional workflow-context questions about future work
management and recurring quality review. Those questions **MUST** be framed as
preparation and context capture, not as permission to create external artifacts
or automations.

If setup asks about recurring quality review, it **MUST** treat ad hoc
`/quality evaluate` as always available rather than as a selectable option.
Recurring review options **SHOULD** be cadence-oriented, such as weekly, per
sprint or iteration, monthly, before major planning or review sessions, custom,
or not now.

Setup **MUST NOT** recommend CI or release gating as the default quality loop.

### Mutation Boundary

Setup **MUST** mutate only the selected `QUALITY.md` file.

Setup **MUST NOT** run `/quality evaluate`, create evaluation artifacts, write
the quality log, create external issues, configure issue trackers, create CI or
release workflows, create scheduled automations, configure Codex automations, or
configure Claude Code routines.

If the target model file is missing, setup **MUST** still delegate deterministic
scaffolding to `qualitymd init` before model authoring.

If the target model file exists and setup would change it, setup **MUST** use a
decision brief before editing.

### Model Authoring

The generated `QUALITY.md` **MUST** follow the bundled authoring guide and
conform to the active QUALITY.md specification exposed to the skill.

The generated `QUALITY.md` **MUST** capture enough body context to explain why
the initial model has its shape. At minimum, the body should include Overview,
Scope, Needs, and Risks, with each section's unknowns and open questions.

The generated `QUALITY.md` **MUST** record setup assumptions where they shape
the model: lifecycle, risk tolerance, modeling rigor, collaboration context,
stakeholder needs, and important missing or non-agent-accessible context.

Root quality factors **MUST** derive from project-specific needs and risks, not
from generic quality labels alone.

Child Areas **SHOULD** be added only when they represent distinct evaluated
entities that make the model clearer or more useful.

Starter Requirements **MUST** be concrete and assessable from agent-accessible
evidence or explicitly name missing evidence or assessment constraints.

Setup **SHOULD** use the standard rating scale unless discovery shows a real
need to customize it.

Setup **SHOULD** include a `quality-md` self-check Area when appropriate under
the existing setup-quality-md-area contract, unless the user declines or the
model file is not in the root Area it governs.

### Validation and Readiness

After writing `QUALITY.md`, setup **MUST** run `qualitymd lint [path]`.

Setup **MUST** inspect the resulting model against the bundled Top 10
QUALITY.md checks before reporting completion. This inspection **MUST** remain a
model-readiness inspection and **MUST NOT** evaluate root Area source quality.

Setup **MUST** report lint failures before offering evaluation as a next step.

Setup **MUST NOT** present the model as ready to evaluate unless body context,
rating scale fit, root/scope alignment, factor coverage, Requirement
assessability, and assessment evidence are good enough to bind evidence to
ratings.

### Completion and Next Steps

Setup completion output **MUST** summarize what changed in `QUALITY.md`, the
lint result, and the model-readiness classification.

Setup completion **MUST** offer next-step choices. The choices **SHOULD**
include:

- continue iterating on `QUALITY.md`;
- run a quality evaluation;
- set up a recurring quality review loop;
- set up recommendation handoff; and
- stop here.

Setup **MUST NOT** automatically take any next-step action.

## Acceptance Criteria

- Setup inspects available repository context before asking questions.
- Each setup question includes a recommended default and confidence signal.
- The current directory is assumed as the root Area convention unless there is
  an explicit or strongly inferred reason to narrow it.
- Agent-heavy development is assumed; collaboration questions ask who else must
  align with agents.
- Lifecycle, risk tolerance, and modeling rigor are separate setup concepts.
- Missing-context prompts are seeded with specific inferred gaps.
- Optional workflow-context prompts distinguish future work-management and
  recurring quality-review posture from setup actions.
- Ad hoc `/quality evaluate` is treated as always available, not as a selectable
  automation option.
- Recurring review is framed around a maintainer-chosen cadence, not default
  CI/release gating.
- Setup mutates only `QUALITY.md`.
- Setup does not write `.quality/log/`.
- Setup does not create evaluation runs, issue-tracker artifacts, CI workflows,
  release gates, scheduler config, Codex automations, or Claude Code routines.
- Generated `QUALITY.md` follows the authoring guide and active specification.
- Generated `QUALITY.md` captures lifecycle, risk tolerance, modeling rigor,
  collaboration context, stakeholder needs, and missing context when they shape
  the model.
- Generated Requirements are assessable from agent-accessible evidence or name
  missing evidence explicitly.
- Setup runs `qualitymd lint` after writing.
- Setup inspects the resulting model with the Top 10 QUALITY.md checks without
  treating that inspection as source evaluation.
- Final setup output reports changed `QUALITY.md`, lint result, model-readiness
  classification, and next-step options.
- No change is made to `SPECIFICATION.md`.

## Durable spec changes

### To add

None.

### To modify

- `specs/skills/quality-skill/quality-skill.md` - update setup's summary,
  mutation surfaces, run-frame expectations, and quality-log writer contract so
  setup writes only `QUALITY.md` and routes next steps.
- `specs/skills/quality-skill/modes/setup.md` - replace the current
  skeleton-plus-guided-population procedure with context analysis, discovery,
  model authoring, lint/readiness inspection, and next-step routing.
- `specs/skills/quality-skill/quality-log.md` - remove setup-created inaugural
  entries from the quality-log writer contract; keep the log limited to
  confirmed model-authoring and recommendation-apply workflows.
- `specs/skills/quality-skill/guides/getting-started-md.md` - align first-run
  model iteration guidance with setup's new completion and next-step behavior.
- `specs/skills/quality-skill/guides/top-10-quality-md-checks-md.md` - confirm
  or refine setup-readiness inspection wording so setup can use the checklist
  after writing without re-running setup or evaluating source.
- Relevant OKF logs and indexes under `specs/` - record durable spec updates
  when they land.

### To rename

None.

### To delete

None.
