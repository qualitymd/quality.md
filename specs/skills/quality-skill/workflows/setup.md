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
authoring, validation, and maturity routing. It implements the shared
contracts in the parent [/quality skill](../quality-skill.md) spec and owns
only the setup-specific behavior below.

The runtime procedure lives at
[`skills/quality/workflows/setup.md`](../../../../skills/quality/workflows/setup.md).
`setup` is dispatched as a mode, but its files live under `workflows/` and its
user-facing behavior is described as a workflow.

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

`setup` may mutate the target `QUALITY.md` model file and **MAY** additionally
write a workflow feedback log under `.quality/logs/`, creating that directory on
demand (see [Workflow feedback log](setup/feedback-log.md)).

`setup` **MUST NOT** run evaluation, create evaluation artifacts, write the
quality log under `.quality/log/`, create external issues, configure issue
trackers, create CI or release workflows, create scheduled automations, configure
Codex automations, or configure Claude Code routines.

> Annotation: the feedback log is the only widening of setup's mutation boundary.
> It is kept narrow — `.quality/logs/` only — and every other prohibition above
> stays in force, so the feedback artifact cannot become a back door for the
> writes setup still forbids. The `.quality/logs/` directory (plural) is distinct
> from the quality log's `.quality/log/` (singular), which setup still must not
> write. — 0066

## Workflow structure

Runtime setup guidance **MUST** read as an operator playbook with ordered steps,
not only as conformance requirements.

The setup workflow **MUST** include these stages, in order:

1. Resolve the target `QUALITY.md` and verify setup prerequisites.
2. Inspect repository context for setup signals.
3. Build a setup brief with inferred defaults, confidence, and evidence.
4. Ask concrete discovery questions.
5. Present a final review recap of the question/answer set and invite a last
   comment or correction before authoring.
6. Run `qualitymd init [path]` when the target model is missing.
7. Synthesize or update `QUALITY.md`.
8. Run lint and classify model maturity.
9. Report completion and next-step choices, and author a workflow feedback log
   when the run had notable experience events.

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

The confidence vocabulary **MUST** be fixed to `Low`, `Med`, and `High`. The
workflow **MUST NOT** use the prior `strongly inferred`, `weakly inferred`, or
`assumed` labels. Each inferred setup-brief item and each recommended discovery
default **MUST** carry one of `Low`, `Med`, or `High`, plus the short evidence
note when evidence exists. A default with no supporting evidence **MUST** be
labeled `Low` and **SHOULD** name the absence of evidence in its note (e.g.
`Low (no signal in repo)`), preserving the "no-evidence, pure default" meaning
the prior `assumed` label carried.

> Annotation: the three-level scale replaces the `strongly/weakly inferred` and
> `assumed` vocabulary, which read awkwardly next to a recommended default in a
> real field run. The per-item evidence note is what lets the simpler scale carry
> the old meaning — including the distinct "pure default, no evidence" signal,
> which now lives in a `Low` label plus an explicit no-evidence note rather than a
> separate word. — 0067

## Discovery questions

Setup **MUST** ask or present the following discovery questions before writing
`QUALITY.md`:

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

### Per-question pedagogy

The setup workflow **MUST** carry authored teaching copy for each of the ten
discovery questions in the runtime skill. For each question, that copy **MUST**
state the purpose of the question — why the dimension matters and what it shapes
in `QUALITY.md` — and **MUST** state how the user can change that answer later.

The teaching copy **MUST** be authored in the workflow itself, not left to
per-run agent improvisation, and **MUST** be written as copy the agent presents
to the user (prose around the question), not as text confined to a structured
tool's option or description fields. The workflow **MUST** present a question's
purpose and how-to-change-later context to the user before or together with that
question, on whatever presentation surface the agent uses.

The workflow framing **MUST** state that setup optimizes for teaching the user
the quality-model dimensions over minimizing interaction round-trips, so the
per-question pedagogy is preserved rather than treated as removable overhead.

> Annotation: the discovery questions always did double duty — capturing context
> and teaching the dimensions a quality model spans — but the teaching half was
> left to per-run improvisation, so it was inconsistent and not human-tunable.
> Authoring the copy in the skill makes the teaching reproducible and tunable, and
> the framing line is a guardrail so a later contributor reading an efficiency
> complaint does not "optimize" the pedagogy back out. Setup runs ~once per
> project, so the extra interaction is worth the legibility. — 0067

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

`setup` **MUST** ask every one of the ten discovery questions on every run,
including questions whose inferred default is high-confidence. High confidence in
an inferred default **MUST NOT** be a reason to skip a question. `setup`
**MUST NOT** drop, merge, or silently default away a question to fit an
interaction surface's limits.

`setup` **MUST** choose the presentation form from the agent's own interaction
capabilities. This guidance **MUST NOT** assume or name a specific agent's
question tool.

When the agent has a structured question affordance with item or option limits,
`setup` **MUST** page all ten questions through it across as many rounds as the
limits require, and **MUST** keep open-ended questions (primary users,
maintainers and collaborators, other stakeholders, missing context) as free text
rather than forcing them into fixed options.

When the agent has no structured question affordance, `setup` **MUST** iterate
the questions one at a time. Each step **MUST** carry that question's recommended
default and confidence signal so the user can confirm or correct it and advance,
and **MUST NOT** require a full prose answer. One-at-a-time iteration is the
default presentation form.

`setup` **MUST NOT** offer an escape that accepts all inferred defaults and skips
the remaining questions. Any prior guidance permitting "accept all defaults and
skip the remaining questions" **MUST** be removed or revised so it does not
contradict asking every question. A per-question fast confirm — the user accepts
the recommended default for a single question and advances without writing prose
— **MAY** remain, because it still presents that question and its teaching copy.
`setup` **MUST** honor an explicit user request to see all ten questions at once
instead of iterating, and **MUST NOT** lead with that escape.

> Annotation: 0065 established the agent-agnostic presentation tiers but kept an
> "accept all defaults and skip the rest" escape. That escape directly
> contradicts the teaching purpose — it drops the per-question beats wholesale —
> so this case removes it. The per-question fast confirm preserves the speed
> benefit for an expert without skipping any question's instruction; show-all-at-
> once still surfaces every question with its teaching copy. — 0067

`setup` **MUST NOT** re-ask context the user has already supplied earlier in the
interaction.

## Final review recap

After all ten discovery questions are answered and before writing `QUALITY.md`,
`setup` **MUST** present a final review recap that lists every discovery question
with its final answer.

The recap **MUST** invite the user to add a last free-text comment or correct any
answer before authoring proceeds. `setup` **MUST** incorporate corrections the
user makes at this step before authoring, and **MUST NOT** require the user to add
a comment to proceed.

The recap **MUST NOT** be the only place a question is surfaced; it supplements,
and does not replace, asking each question during discovery.

> Annotation: the recap is a consolidated confirmation-and-teaching moment, not a
> replacement for per-question iteration — making it the only confirmation would
> collapse discovery back into a single batch screen and lose the teaching beats.
> It also gives cross-cutting remarks a home that the mid-flow missing-context
> question cannot fully serve. — 0067

## Model authoring

`setup` **MUST** drive `qualitymd init` for deterministic scaffolding when the
model file is absent. It **MUST NOT** reimplement scaffolding, validation, CLI
installation tooling, or source-driven authoring judgment.

When `setup` scaffolds with `qualitymd init`, it **MUST** read the scaffolded
file before authoring it, so a single authoring pass does not fail a
read-before-write guard.

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

`setup` **MUST** classify the resulting model's maturity against the bundled Top
10 QUALITY.md checks before reporting completion, using that guide's condensed
close checklist and reading the full checks only when the maturity call is
borderline. This inspection **MUST** remain a model-maturity inspection and
**MUST NOT** evaluate root Area source quality. The model-maturity classification
(`starter`, `immature`, `evaluation-ready`) **MUST** be reported as distinct from
the lifecycle `readiness` that `qualitymd status` owns.

## Completion criteria

`setup` is complete when the target model exists, lint has run, the model has
received context-informed authoring or a clearly reported user-deferred
authoring step, and setup has reported model maturity. Completion output
**MUST** summarize the `QUALITY.md` change, lint result, maturity
classification, important remaining model gaps, and next-step choices.

Next-step choices **SHOULD** include continuing to iterate on `QUALITY.md`,
running evaluation, setting up a recurring quality review loop, setting up
recommendation handoff, and stopping. `setup` **MUST NOT** automatically take
any next-step action.

## Feedback log

At the close of the run, `setup` **SHOULD** author a workflow feedback log under
`.quality/logs/` capturing notable experience events from that run, and **MAY**
omit it when nothing notable occurred. The artifact contract — location, naming
(`<timestamp>-setup-feedback-log.md`), environment header, body schema,
redaction, and no-transmission posture — is owned by the
[Workflow feedback log](setup/feedback-log.md) sub-spec. Writing or omitting a
feedback log **MUST NOT** change setup's completion criteria, maturity
classification, or next-step routing.

> Annotation: the feedback log records the *experience* of running setup so the
> skill, CLI, and prompts can improve from real runs — distinct from the user-
> facing completion summary, which stays terse. The [0065 setup refinements](../../../../changes/archive/0065-setup-discovery-and-close-refinements.md)
> existed only because a human hand-captured such friction once; the feedback log
> makes that loop durable. — 0066
