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
authoring, validation, and important-gap reporting. It implements the shared
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

`setup` may mutate the target `QUALITY.md` model file and **MUST** additionally
write and update a workflow feedback log under `.quality/logs/`, creating that
directory on demand (see the shared
[workflow feedback log](../workflow-feedback-log.md) and setup-specific
[Setup feedback log](setup/feedback-log.md)).

`setup` **MUST NOT** run evaluation, create evaluation artifacts, write the
quality log under `.quality/log/`, create external issues, configure issue
trackers, create CI or release workflows, create scheduled automations, configure
Codex automations, or configure Claude Code routines.

> Annotation: the feedback log is the only widening of setup's mutation boundary.
> It is kept narrow — the current run's `.quality/logs/` feedback file only —
> and every other prohibition above stays in force, so the feedback artifact
> cannot become a back door for the writes setup still forbids. The
> `.quality/logs/` directory (plural) is distinct from the quality log's
> `.quality/log/` (singular), which setup still must not write. — 0066, 0068

## Workflow structure

Runtime setup guidance **MUST** read as an operator playbook with ordered steps,
not only as conformance requirements.

The setup workflow **MUST** include these stages, in order:

1. As the first user-visible output, before any tool call, present a short
   educational opening — a warm welcome, what QUALITY.md gives teams and agents, a
   short phase roadmap, the read-only-scan-first and review-before-changes
   boundary — and emit the run frame in the same block.
2. Verify setup prerequisites as a fail-fast gate, and confirm the resolved
   target `QUALITY.md` the opening's run frame already named.
3. Inspect repository context for setup signals.
4. Build a setup brief with inferred defaults, confidence, and evidence.
5. Present a concise setup preview distilled from the setup brief, then create
   the current run's setup feedback log if the run continues into discovery.
6. Ask concrete discovery questions and present the human context checkpoint.
7. Present a final review recap of the question/answer set and wait for an
   explicit review-gate response before authoring.
8. Run `qualitymd init [path]` when the target model is missing.
9. Synthesize or update `QUALITY.md`.
10. Run lint and identify important model gaps.
11. Report completion and one immediate next step, and finalize the setup feedback
    log.

The workflow **MUST NOT** ask the user to design Factors, child Areas,
Requirements, or Rating Levels cold. The skill derives model shape from the
setup brief, discovery answers, authoring guide, and repository context.

## Context analysis and setup brief

Setup **MUST** present the opening orientation as the first user-visible output of
the run, before any tool call — including before the CLI prerequisite check,
repository inspection, and any filesystem read or write. The opening's content
has no tool dependency. The opening **MUST** explain the setup value proposition
in operational terms: QUALITY.md gives AI assistants, coding agents, and teams a
holistic definition of quality tailored to their project, so they can stay
aligned, identify critical risks and issues, and keep improving. The opening
**MUST** also state that the first phase is a read-only context scan and that the
user will review before changes are written.

The opening **MUST** include a short phase roadmap of what setup will do —
read-only scan, calibration questions, review, write, verify — so the read-only
scan that follows reads as an expected step rather than a hang. The opening
**SHOULD** include a brief cue that the read-only scan may take a moment on a
large repository.

Setup **MUST** emit the run frame as part of this first-output block, alongside
the opening and before any tool call. Run-frame emission **MUST NOT** be gated on
CLI prerequisite verification or any other tool result; the run frame's resolved
model path is derived from the invocation (the explicit path when supplied,
otherwise `QUALITY.md` in the current working directory), not from a filesystem
probe. After the first-output block, setup **MUST** run the CLI prerequisite
check as a fail-fast gate before the read-only context scan, stopping with a clear
message when the CLI is missing or unsupported.

The opening **SHOULD** stay short: a value-proposition sentence, the phase
roadmap, immediate workflow status, and the review-before-changes boundary. It
**MUST NOT** become a marketing splash screen or replace the run frame, the setup
preview, the discovery questions, the human context checkpoint, or the review
gate.

> Annotation: 0096 added the opening but a field run showed it did not reach the
> user *first* — the agent front-loaded CLI checks, repository scans, and even the
> feedback-log write before flushing any text, so the welcome, run frame, and
> preview all arrived after 1–2 minutes of silence. The fix is ordering: the
> opening and run frame have no tool dependency (the frame's only variable, the
> resolved model path, is known from the invocation), so they become genuine
> first output before any tool call, and the CLI gate runs after the opening
> rather than gating it. The phase roadmap makes the subsequent silent scan read
> as an expected step. — 0098

`setup` **MUST** inspect available repository context before asking setup
questions. Relevant context includes README and docs, repository structure,
package metadata, tests, contributor docs, existing agent instructions, and
visible workflow signals that affect first-model authoring. This inspection
**MUST** stay bounded to setup signals and **MUST NOT** become source-quality
evaluation.

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
rating-scale recommendation, collaboration context, inferred primary users and
outcomes, inferred maintainer or collaborator needs, inferred other stakeholder
needs, missing or non-agent-accessible context, and candidate model shape.

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

Before discovery questions, setup **MUST** present a user-facing setup preview
distilled from the setup brief. The preview **MUST** include:

- likely root Area or boundary;
- likely domain or quality context;
- visible evidence used for the inference;
- likely candidate model shape, such as initial Areas, Factors, or quality
  concerns;
- missing or not-agent-accessible context;
- the next user action.

The setup preview **MUST** be framed as draft context for correction, not as
confirmed fact. It **MUST** distinguish high-confidence evidence-backed
inferences from low-confidence or missing context. It **MUST NOT** replace the
required discovery questions, human context checkpoint, or final review gate.

> Annotation: the preview is the first project-specific value the user sees. It
> turns discovery into calibration of an evidence-backed draft rather than a cold
> questionnaire, while still preserving the teaching beats and confirmation
> boundaries that prevent low-confidence inferences from becoming facts. — 0096

## Discovery inputs

Setup **MUST** ask or present the following discovery inputs before writing
`QUALITY.md`:

1. **Root area.** Should this `QUALITY.md` model the whole current project, or a
   narrower Area?
2. **Domain.** What kind of thing is this model evaluating?
3. **Lifecycle.** Which stage best fits: exploratory, pre-release, active
   production, maintenance, or sunset?
4. **Risk tolerance.** How costly is poor quality here: high cost, moderate
   cost, or low cost? Setup maps the visible cost answer to the internal
   risk-tolerance meaning: high cost maps to low tolerance, moderate cost maps to
   moderate tolerance, and low cost maps to high tolerance.
5. **Rating Scale.** Should the model use the recommended
   `outstanding`, `target`, `minimum`, `unacceptable` Rating Scale?
6. **Human context checkpoint.** A draft for correction that covers primary
   users and outcomes, maintainers and collaborators, other stakeholders, and
   missing or not-agent-accessible context.

Each discovery question **MUST** include a recommended answer and confidence
signal. Each human context checkpoint item **MUST** include a recommended answer
and confidence signal when setup has inferred one from repository evidence.

When a discovery question is rendered as a small closed-choice set, setup
**MUST** present numbered options, put the recommended option first, mark it as
recommended, and make `1` the shortest confirmation. Closed-choice option labels
**MUST** match the question's visible axis; setup **MUST** map any different
internal value behind the scenes instead of making the user translate while
answering.

> Rationale: setup questions teach the quality-model dimension while collecting
> input. A recommended default is the easiest path only when it is option `1`,
> and a cost question should be answered with cost labels even when the model
> records risk tolerance. — 0099

The Rating Scale question **MUST** explain that Rating Levels are configurable in
`QUALITY.md` and are not baked into the format. It **MUST** recommend the
standard four-level scale for most first models, keep its stable Rating Level IDs
as `outstanding`, `target`, `minimum`, and `unacceptable`, and use the default
display titles `🟢 Outstanding`, `🔵 Target`, `🟡 Minimum`, and
`🔴 Unacceptable` unless the project context calls for plain or custom titles. It
**MUST** explain each recommended level's role: `outstanding` is a stretch band
where further investment may need ROI justification, `target` is the expected
good-enough bar, `minimum` is the acceptable floor that still warrants
improvement, and `unacceptable` is below the floor. It **MUST NOT** treat emoji
markers as rating identity, ordering, or semantics.

Setup **MUST NOT** ask the user to invent custom Rating Level names during
discovery. When the user rejects the recommended scale, setup **MAY** choose a
simple alternate scale only when project context clearly supports it, such as a
pass/fail gate; otherwise setup **SHOULD** use the recommended scale and record
the scale decision as an open question or assumption in the model body.

### Discovery pedagogy

The setup workflow **MUST** carry authored teaching copy for each discovery
question and human context checkpoint item in the runtime skill. For each
question or checkpoint item, that copy **MUST** state the purpose of the input —
why the dimension matters and what it shapes in `QUALITY.md`.

The teaching copy **MUST** be authored in the workflow itself, not left to
per-run agent improvisation, and **MUST** be written as copy the agent presents
to the user (prose around the question or checkpoint), not as text confined to a
structured tool's option or description fields. The workflow **MUST** present a
question's or checkpoint item's purpose/context to the user before or together
with that input, on whatever presentation surface the agent uses.

The workflow **MAY** state once outside the individual question copy that
`QUALITY.md` is a living document and that setup answers can be revised later.
It **MUST NOT** include per-question "how to change it later" copy or equivalent
per-question lifecycle guidance.

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

### Human context checkpoint

After the first five setup discovery questions, `setup` **MUST** present a human
context checkpoint instead of asking separate open-ended questions for primary
users/outcomes, maintainers/collaborators, other stakeholders, and missing
context.

The checkpoint **MUST** present repository-inferred human context as a draft for
confirmation or correction, with confidence labels and evidence notes where
useful. It **MUST** let the user confirm, correct, fill in terse fragments, or
point to agent-accessible evidence the setup pass missed.

The checkpoint **MUST** make the correction action the strongest visible element
and **MUST** include an `Answer` line, or materially equivalent wording, that
states how to accept the draft or provide terse corrections. The unresolved-item
warning **MUST** support that action rather than precede it as the first or
strongest element.

The checkpoint **MUST** cover all of these dimensions:

- primary users and outcomes;
- maintainers and collaborators;
- other stakeholders;
- missing or not-agent-accessible context.

The checkpoint **MUST** state that unanswered low-confidence or not-visible
items will be recorded as Unknown, open questions, or low-confidence inference as
appropriate, not treated as confirmed facts.

The checkpoint **MUST** end with a short, prioritized list of the highest-value
corrections: who the evaluated thing is for, what outcome matters most, and
whether data, compliance, availability, or business-criticality constraints
exist.

The checkpoint **MUST NOT** end human-context discovery with a broad catch-all
question that can obscure the primary users/outcomes, maintainer/collaborator,
and other stakeholder dimensions.

> Annotation: the human context dimensions are still required because they ground
> Needs, Risks, Factors, and Unknowns. The change is their response shape:
> reviewing and correcting a draft is easier than answering four essay prompts,
> especially when the final prompt is the broadest one. Silence on a
> low-confidence item is not evidence, so setup records unresolved material
> context honestly instead of treating omission as confirmation. — 0072

The missing-context checkpoint item **MUST** be seeded from repository analysis
rather than phrased as a blank "anything else?" prompt.

For missing-context discovery, setup **MUST** treat material context as
agent-accessible only when it is available through repository content, cited
local paths, configured tools, linked public sources, or explicit user-provided
setup context. Missing-context discovery **MUST** distinguish context that is
visible from agent-accessible evidence, context that should be recorded as
unknown or not agent-accessible, and context the user provides during setup.

If any generated missing-context choices are used, they **MUST NOT** invite the
user to assume that product purpose, operational context, stakeholder needs,
telemetry, security/compliance posture, incident history, SLAs, production
metrics, or similar material project-specific facts are understood when the setup
brief marks them `Low` confidence or not visible from evidence. When an option
excludes an identified material gap from unknowns, the option **MUST** make the
provenance explicit: either the user is providing the missing context during
setup, or the user is pointing setup to agent-accessible evidence it missed. For
material gaps with low or no evidence, the recommended option **SHOULD** record
the gaps as unknowns or open questions.

> Annotation: the missing-context checkpoint item exists to prevent guessing. If
> setup already marked a project-specific fact as low/no-evidence, an "assume it
> is understood" option silently turns tacit maintainer knowledge into evidence
> the agent cannot inspect. — 0070

The maintainer/collaborator checkpoint item **MUST** assume agent-heavy
development and ask which human collaborators, reviewers, maintainers, or
stakeholders also need to align with the quality bar.

Setup **MUST NOT** ask a review-posture discovery question. Setup **MUST NOT**
ask for review cadence, recurrence, quality-loop options, recommendation
handling, work-handoff destination, issue tracker, or automation preferences. Ad
hoc `/quality evaluate` **MUST** be treated as always available after setup
rather than as a selectable automation option. Setup **MUST NOT** recommend CI or
release gating as the default quality loop.

## Prompt form

`setup` **MUST** ask every one of the first five discovery questions on every run
and **MUST** present the human context checkpoint on every run. High confidence
in an inferred default **MUST NOT** be a reason to skip a question or checkpoint
item. `setup` **MUST NOT** drop or silently default away a discovery dimension to
fit an interaction surface's limits.

`setup` **MUST** choose the presentation form from the agent's own interaction
capabilities. This guidance **MUST NOT** assume or name a specific agent's
question tool.

When the agent has a structured question affordance with item or option limits,
`setup` **MUST** page the first five discovery questions through it across as
many rounds as the limits require, then present the human context checkpoint as
free text rather than forcing it into fixed options.

When the agent has no structured question affordance, `setup` **MUST** iterate
the questions one at a time. Each step **MUST** carry that question's recommended
default and confidence signal so the user can confirm or correct it and advance,
and **MUST NOT** require a full prose answer. After the fifth question, `setup`
**MUST** present the human context checkpoint. One-at-a-time iteration is the
default presentation form.

Each discovery interaction block **MUST** visually emphasize the primary
question or checkpoint call to action, preferably with bold Markdown. The
question's purpose, recommended answer, confidence/evidence, and shortest
acceptable response **MUST** stay adjacent to the question. When the surface
supports Markdown, labels such as `Why it matters`, `Recommended`, `Confidence`,
and `Answer` **SHOULD** be bold for scanning.

For closed-choice discovery questions, the shortest acceptable response
**MUST** be `1` when accepting the recommendation. Other numbered responses
**MUST** remain available for the visible alternatives.

The human context checkpoint **MUST** make the correction action explicit. When
reviewing several inferred values together, it **SHOULD** use a table or
similarly compact structure unless the interaction surface makes that less
usable.

For the open-ended root-area and domain questions, setup **MUST** still use the
shared agent-mediated question shape: visually emphasized question, adjacent
`Why it matters`, `Recommended`, `Confidence`, and `Answer` labels when Markdown
is available, and a shortest accept path such as replying `yes` to accept the
recommendation or naming the correction.

`setup` **MUST NOT** offer an escape that accepts all inferred defaults and skips
the remaining questions or the human context checkpoint. Any prior guidance
permitting "accept all defaults and skip the remaining questions" **MUST** be
removed or revised so it does not contradict presenting every discovery input. A
per-question fast confirm — the user accepts the recommended default for a single
question and advances without writing prose — **MAY** remain, because it still
presents that question and its teaching copy. `setup` **MUST** honor an explicit
user request to see all discovery inputs at once instead of iterating, and
**MUST NOT** lead with that escape.

> Annotation: 0065 established the agent-agnostic presentation tiers but kept an
> "accept all defaults and skip the rest" escape. That escape directly
> contradicts the teaching purpose — it drops the per-question beats wholesale —
> so this case removes it. The per-question fast confirm preserves the speed
> benefit for an expert without skipping any question's instruction; show-all-at-
> once still surfaces every question with its teaching copy. — 0067

`setup` **MUST NOT** re-ask context the user has already supplied earlier in the
interaction.

## Final review recap

After all discovery questions and the human context checkpoint are answered and
before writing `QUALITY.md`, `setup` **MUST** stop for a review gate and present
a final review recap that lists every asked discovery question and checkpoint
item with its final answer.

`setup` **MUST** wait for a user response to the final recap before writing or
editing `QUALITY.md`. Completing the final discovery question or a structured
question-tool page **MUST NOT** satisfy this review gate.

Because this gate authorizes creating or editing `QUALITY.md`, the final review
prompt **MUST** use a decision-brief shape, whether the file is newly scaffolded
or already exists. The brief **MUST** name the proposed `QUALITY.md` change, the
evidence or reason for writing it, the recommended option, at least one
non-mutating alternative, and the done criterion or verification.

The recap **MAY** receive a correction, a cross-cutting comment, broader
last-call setup context, or an explicit confirmation such as "looks good",
"continue", "write it", or an equivalent phrase. `setup` **MUST** incorporate
corrections and additional review-gate context before authoring, and **MUST NOT**
require the user to add a substantive comment to proceed.

The final review prompt **MUST** use friendly, colloquial, open-ended wording
that keeps the confirmation or correction action visually primary while inviting
useful context beyond corrections as secondary material: priorities, worries,
wording, edge cases, repo-invisible context, or anything else the user considers
important. The runtime setup workflow **SHOULD** use a decision brief whose
primary action has this meaning:

```text
Reply `looks good` to write `QUALITY.md`, or send corrections.
```

When setup presents a separate fallback decision brief for updating an existing
`QUALITY.md` after the review gate and that brief is a true binary update-or-stop
confirmation, it **MUST** show `y`/`n` as the visible shortest answer path. This
does not replace the final review gate's `looks good` fast path when that gate is
framed as a review-and-correction interaction.

The recap **MUST NOT** be the only place a question or checkpoint item is
surfaced; it supplements, and does not replace, asking or presenting each input
during discovery.

> Annotation: the recap is a consolidated confirmation-and-teaching moment, not a
> replacement for per-question iteration — making it the only confirmation would
> collapse discovery back into a single batch screen and lose the teaching beats.
> It also gives cross-cutting remarks a home that the mid-flow missing-context
> question cannot fully serve. The prompt names examples of useful late context
> so users do not have to frame every non-confirmation as a correction. — 0067,
> 0071

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
rating-scale choice, collaboration context, stakeholder needs, and important
missing or non-agent-accessible context.

When a user provides missing context during setup, the generated body **MUST**
preserve that provenance clearly enough that a later reader can tell the context
came from explicit setup input rather than repository inspection. Setup **MUST NOT**
treat tacit maintainer, operator, contributor, or stakeholder knowledge as
available evidence unless it has been explicitly provided or cited.

Setup-authored Factors **MUST** derive from project-specific needs, risks,
stakeholder concerns, component boundaries, and available evidence, not generic
quality labels or setup-question labels alone. Child Areas **SHOULD** be added
only when they represent distinct evaluated entities. Starter Requirements
**MUST** be concrete and assessable from agent-accessible evidence or explicitly
name missing evidence or assessment constraints.

For an agent-collaborated composite root, setup **MUST** propose
`agent-harnessability` / Agent Harnessability as a model-wide umbrella factor by
default, decomposed into the authoring guide's seven sub-factors:
`agent-accessibility`, `task-specifiability`, `agent-operability`,
`continuity`, `self-verifiability`, `enforcement-of-standards`, and
`containment-of-action`.
Setup **MUST** define it as the degree to which the project's own materials and
tooling equip an AI agent to understand the project, take scoped work, operate the
environment, preserve and resume state, verify its output, and stay safely bounded
while preserving clear human direction, review, and accountability. Setup **MUST
NOT** omit Agent
Harnessability because the project's harness is thin or absent; that condition is
a rating gap and finding, not a factor-omission reason. Setup **MUST** keep Agent
Harnessability distinct from the agent harness constituent: the factor rates how
each constituent equips an agent, while the harness Area rates the checked-in,
project-owned governing artifacts themselves.

For an agent-collaborated composite root whose agent harness is germane, setup
**MUST** actively check for project-owned runtime harness machinery such as
hooks, tool/MCP definitions, sandbox or permission policy, orchestration config,
and subagent config. Setup **MUST** scope the generated agent-harness Area as the
checked-in steering and owned-control artifacts, not instructions alone. Setup
**MUST** align the generated Area `description` and projection-boundary comment
with that scope and distinguish the Area from the broader Agent Harnessability
factor. Setup **MUST** surface project-owned runtime harness machinery in the
agent-harness Area or, when large enough to carry distinct factors, in its own
Area; it **MUST NOT** silently fold that machinery into prose guidance or drop
it. Setup **SHOULD** include in its user-facing recap one line naming that the
agent-harness Area is the checked-in steering and owned-control artifacts while
Agent Harnessability rates the broader equipping capability.

> Annotation: Agent Harnessability is the factor-side projection of the same
> agent-collaboration concern that 0080 made a modeled-by-default constituent.
> Leading with the factor prevents a thin harness from being deferred out of the
> model and makes the agent-equipping quality visible across every constituent. —
> 0081, refined by 0085

Setup **SHOULD** use the standard Rating Scale, including its emoji-prefixed
human display titles, unless discovery and the body show a real need to
customize it. If the user rejects the recommended scale but the project context
does not clearly support an alternate scale, setup **SHOULD** use the
recommended scale and record the scale decision as an open question or assumption
in the model body.

Setup **SHOULD** include a `quality-md` Area that evaluates the `QUALITY.md`
artifact itself against the active authoring guide family unless the user declines or
the model file is not in the root Area it governs. The Area **SHOULD** use the
key `quality-md`, a title of the form `<Root Title> QUALITY.md`, an Area
`description`, and an explicit path-based `source` such as `./QUALITY.md`. It
**MUST NOT** use prose aliases such as `(this file)` for `source`.

When setup adds that Area, it **SHOULD** include concise YAML comments that
distinguish the Area `source` from the Requirement `assessment`. It **SHOULD**
use one Area-level Requirement with `factors` when the active authoring guide family
defines one coherent judgment across multiple Factors.

## Stop conditions

`setup` **MUST** stop before CLI-dependent work when the `qualitymd` CLI is
missing, outside the released-install SemVer range declared by the skill, or a
local development build lacks required commands.

`setup` **MUST** run `qualitymd lint` after writing `QUALITY.md`. It **MUST**
report lint failures before offering evaluation as a next step.

When lint passes, `setup` **MUST** inspect the model for important gaps that
materially affect first-model usefulness. This inspection **MUST** remain a
model-usefulness inspection and **MUST NOT** evaluate root Area source quality.
Important gaps include thin or generic Markdown body context, missing material
unknowns or open questions, factors that do not trace to the body's needs and
risks, vague or unassessable Requirements, missing germane constituent kinds, and
missing Agent Harnessability coverage for an agent-collaborated composite root.
Setup **MUST NOT** collapse those gaps into a `starter`, `immature`,
`evaluation-ready`, `Maturity`, or `Evaluation readiness` label.

For a composite root, `setup` **MUST** report an important gap when a germane
constituent that hits neither disqualifier (no distinct concerns; not germane /
outside the boundary) is left unmodeled or recorded only as a deferral; a bare
deferral or Scope note **MUST NOT** be treated as satisfying constituent coverage.

For an agent-collaborated composite root, `setup` **MUST** report an important
gap when Agent Harnessability or its sub-factors are missing from the model-wide
factors without a clear not-germane boundary. A thin or absent harness **MUST
NOT** satisfy that boundary. An existing `harnessability` factor with the expected
six sub-factors **SHOULD** count as semantic coverage, while setup's authoring
guidance **SHOULD** recommend renaming it to `agent-harnessability` / Agent
Harnessability.

> Rationale: model-by-default is enforceable only if under-coverage is reported
> as a concrete model gap; a deferral note previously passed the close check
> silently. — 0080, 0092

## Completion criteria

`setup` is complete when the target model exists, lint has run, the model has
received context-informed authoring or a clearly reported user-deferred
authoring step, and setup has reported important model gaps. Completion output
**MUST** summarize the `QUALITY.md` change, lint result, important remaining
model gaps, what setup did not do, and one immediate next step.
Completion output **MUST** follow the shared agent-mediated UX contract:
status-first, with scannable labels for changed artifacts, validation, important
gaps, what setup did not do, and the recommended next action.

The recommended next step **SHOULD** be one of: continue iterating on
`QUALITY.md`, run `/quality evaluate`, or stop here. If important model gaps
remain, setup **SHOULD** recommend continued iteration. If no important gaps are
visible and lint passed, setup **SHOULD** recommend either `/quality evaluate` or
stopping based on the user's immediate goal. `setup` **MUST NOT** automatically
take any next-step action.

## Feedback log

After it presents the setup preview, when the run continues into discovery,
`setup` **MUST** create a setup feedback log under `.quality/logs/` for the
current run. As setup progresses, it **MUST** update the
current run's file when material workflow-experience events occur, and at close
it **MUST** finalize the log with terminal status, outcome, effort when
available, and explicit no-notable-content notes for empty sections. The artifact
contract is owned by the shared
[workflow feedback log](../workflow-feedback-log.md) spec; setup-specific
creation, naming (`<timestamp>-setup-feedback-log.md`), and finalization rules
are owned by the [Setup feedback log](setup/feedback-log.md) sub-spec. Writing,
updating, or finalizing a feedback log **MUST NOT** change setup's completion
criteria, important-gap judgment, or next-step routing.

> Annotation: the feedback log records the *experience* of running setup so the
> skill, CLI, and prompts can improve from real runs — distinct from the user-
> facing completion summary, which stays terse. Always creating the log removes
> ambiguity from absence and preserves partial feedback when a setup run is
> interrupted. — 0066, 0068
