---
type: Functional Specification
title: /quality setup
description: Behavioral component spec for the /quality setup workflow.
tags: [skill, quality, mode, setup, workflow]
timestamp: 2026-06-23T00:00:00Z
---

# /quality setup

`setup` is the `/quality` workflow that creates or updates a useful first
`QUALITY.md` through context inspection, a concrete discovery prompt, model
authoring, validation, and readiness routing. It implements the shared
contracts in the parent [/quality skill](../quality-skill.md) spec and owns
only the setup-specific behavior below.

The runtime procedure lives at
[`skills/quality/modes/setup.md`](../../../../skills/quality/modes/setup.md).
The file path remains under `modes/` for dispatch compatibility, but user-facing
setup behavior is described as a workflow.

This document uses BCP 14 keywords only for testable conformance requirements.
The key words "MUST", "MUST NOT", "SHOULD", and "MAY" are to be interpreted as
described in [RFC 2119](../../../../docs/reference/rfc2119.md) and
[RFC 8174](../../../../docs/reference/rfc8174.md) when, and only when, they
appear in all capitals.

## Purpose and routing

`setup` is selected when no model file is present, when the user explicitly asks
to create or initialize a QUALITY.md file, or when read-only orientation routes
to bootstrap or first-population work.

The setup workflow's purpose is to produce or improve a valid, useful
project-specific model while keeping setup's mutation boundary narrow. It is not
an evaluation workflow and does not rate evaluated source.

## Mutation surface and artifacts

`setup` may mutate only the target `QUALITY.md` model file.

`setup` **MUST NOT** run evaluation, create evaluation artifacts, write the
quality log, create external issues, configure issue trackers, create CI or
release workflows, create scheduled automations, configure Codex automations, or
configure Claude Code routines.

## Workflow structure

Runtime setup guidance **MUST** read as an operator playbook with ordered steps,
not only as conformance requirements.

The setup workflow **MUST** include these stages, in order:

1. Resolve the target `QUALITY.md` and verify setup prerequisites.
2. Inspect repository context for setup signals.
3. Build a setup brief with inferred defaults, confidence, and evidence.
4. Ask concrete discovery questions.
5. Run `qualitymd init [path]` when the target model is missing.
6. Synthesize or update `QUALITY.md`.
7. Run lint and inspect model readiness.
8. Report completion and next-step choices.

The workflow **MUST NOT** ask the user to design Factors, child Areas,
Requirements, or Rating Levels cold. The skill derives model shape from the
setup brief, discovery answers, authoring guide, and repository context.

## Context analysis and setup brief

`setup` **MUST** inspect available repository context before asking setup
questions. Relevant context includes README and docs, repository structure,
package metadata, tests, contributor docs, existing agent instructions, and
visible workflow or work-management signals. This inspection **MUST** stay
bounded to setup signals and **MUST NOT** become source-quality evaluation.

Setup **MUST** treat the current directory as the default root Area convention
unless the user supplied an explicit model path or repository context strongly
indicates a narrower root.

When the root spans multiple workspaces, packages, or services, setup **SHOULD**
delegate a bounded component census — one bounded exploration pass per component
capturing purpose, entry points, external systems, risk surfaces, and test or CI
coverage — to produce structured Area candidates. This census **MUST** remain
optional and proportional: a small single-package root does not require it.

Before asking discovery questions, setup **MUST** build a concise setup brief
containing root Area, domain, lifecycle, risk tolerance, modeling rigor,
collaboration context, inferred primary users and outcomes, inferred maintainer
or collaborator needs, inferred other stakeholder needs, missing or
non-agent-accessible context, review posture when visible, and candidate model
shape.

Every inferred setup brief item **MUST** include a recommended default,
confidence signal, and short evidence note when evidence exists.

The confidence vocabulary **MUST** be fixed to:

- `strongly inferred`;
- `weakly inferred`; and
- `assumed`.

## Discovery questions

Setup **MUST** ask or present the following discovery questions before writing
`QUALITY.md`, unless the user explicitly asks to accept all inferred defaults:

1. **Root area.** Should this `QUALITY.md` model the whole current project, or a
   narrower Area?
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
9. **Missing context.** The skill thinks these important inputs are not visible:
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
issue-tracker artifacts. Ad hoc `/quality evaluate` **MUST** be treated as
always available rather than as a selectable automation option. Setup **MUST
NOT** recommend CI or release gating as the default quality loop.

Setup **MAY** ask an additional work-handoff question about where future
evaluation recommendations should usually go. If asked, it **MUST** say setup
will not create issues or configure integrations.

## Prompt form

The discovery questions **MAY** be presented as one compact prompt or as a short
sequence when the interaction surface works better that way.

A compact prompt **MUST** ask the user to confirm or correct the setup
assumptions instead of requiring full prose answers to every question.

A compact prompt **MUST** preserve all required questions, defaults, confidence
signals, and seeded missing-context content.

A short sequence **SHOULD** group closely related questions together, such as
root Area plus domain, lifecycle plus risk plus rigor, stakeholders plus needs,
and missing context plus review posture.

## Model authoring

`setup` **MUST** drive `qualitymd init` for deterministic scaffolding when the
model file is absent. It **MUST NOT** reimplement scaffolding, validation, CLI
installation tooling, or source-driven authoring judgment.

After discovery and scaffolding when needed, `setup` **MUST** write a model that
follows the authoring guide and active specification. The model **MUST** address
the Markdown body's Overview, Scope, Needs, and Risks, including each section's
unknowns, open questions, and any material support that is not
agent-accessible. The body **MUST** preserve setup assumptions where they shape
the model: root Area, domain, lifecycle, risk tolerance, modeling rigor,
collaboration context, stakeholder needs, important missing or
non-agent-accessible context, and review posture when it affects model use.

Setup-authored Factors **MUST** derive from project-specific needs, risks,
stakeholder concerns, component boundaries, and available evidence, not generic
quality labels or setup-question labels alone. Child Areas **SHOULD** be added
only when they represent distinct evaluated entities. Starter Requirements
**MUST** be concrete and assessable from agent-accessible evidence or explicitly
name missing evidence or assessment constraints.

Setup **SHOULD** use the standard Rating Scale unless discovery shows a real
need to customize it.

Setup **SHOULD** include a `quality-md` Area that evaluates the `QUALITY.md`
artifact itself against the active authoring guide unless the user declines or
the model file is not in the root Area it governs. The Area **SHOULD** use the
key `quality-md`, a title of the form `<Root Title> QUALITY.md`, an Area
`description`, and an explicit path-based `source` such as `./QUALITY.md`. It
**MUST NOT** use prose aliases such as `(this file)` for `source`.

When setup adds that Area, it **SHOULD** include concise YAML comments that
distinguish the Area `source` from the Requirement `assessment`. It **SHOULD**
use one Area-level Requirement with `factors` when the active authoring guide
defines one coherent judgment across multiple Factors.

## Stop conditions

`setup` **MUST** stop before CLI-dependent work when the `qualitymd` CLI is
missing, outside the released-install SemVer range declared by the skill, or a
local development build lacks required commands.

`setup` **MUST** run `qualitymd lint` after writing `QUALITY.md`. It **MUST**
report lint failures before offering evaluation as a next step.

`setup` **MUST** inspect the resulting model against the bundled Top 10
QUALITY.md checks before reporting completion. This inspection **MUST** remain a
model-readiness inspection and **MUST NOT** evaluate root Area source quality.

## Completion criteria

`setup` is complete when the target model exists, lint has run, the model has
received context-informed authoring or a clearly reported user-deferred
authoring step, and setup has reported model readiness. Completion output
**MUST** summarize the `QUALITY.md` change, lint result, readiness
classification, important remaining model gaps, and next-step choices.

Next-step choices **SHOULD** include continuing to iterate on `QUALITY.md`,
running evaluation, setting up a recurring quality review loop, setting up
recommendation handoff, and stopping. `setup` **MUST NOT** automatically take
any next-step action.
