---
title: QUALITY.md
ratingScale:
  - level: outstanding
    title: 🟢 Outstanding
    description: "The stretch band: the artifact exceeds the quality requirement with meaningful margin."
    criterion: "Exceeds the requirement; satisfies it with margin to spare."
  - level: target
    title: 🔵 Target
    description: "The expected good state: the artifact satisfies the quality requirement."
    criterion: "Satisfies the requirement."
  - level: minimum
    title: 🟡 Minimum
    description: "The acceptable floor: the artifact falls short of the goal but remains good enough to proceed."
    criterion: "Falls short of the goal but holds the acceptable floor."
  - level: unacceptable
    title: 🔴 Unacceptable
    description: "Below the floor: the artifact does not satisfy the quality requirement acceptably."
    criterion: "Falls below the acceptable floor."
areas:
  format-spec:
    title: Format specification
    source: ./SPECIFICATION.md
    requirements:
      "the format specification is complete":
        factors:
          - clarity
          - consistency
          - verifiability
          - extensibility
          - usability
        assessment: >
          Every frontmatter field and recommended body section has its shape,
          allowed values, requiredness, cardinality, and any default; the spec
          states how a conforming reader treats malformed or omitted content and
          edge cases. An implementer could build a parser, and an author write a
          valid file, from the spec alone.
    factors:
      clarity:
        title: Clarity
        description: >
          Can each rule be read in only one way? A spec governs independent
          implementations that never confer, so every obligation must land with
          one settled meaning and force.
        requirements:
          "the format specification admits a single interpretation":
            assessment: >
              Each rule admits one reading. Obligation strength is explicit
              where it affects conformance, BCP 14 keywords are used sparingly
              and consistently, and no normative statement leans on a vague
              quantifier without a stated bound.
          "the format specification separates rules from rationale":
            assessment: >
              A reader can always tell whether a sentence states a binding rule
              or merely explains one, and a rule never appears only inside an
              example or aside.
          "the format specification defines its terms before use":
            assessment: >
              Every technical term used in a rule is defined before, or at, the
              point it is first used.
      consistency:
        title: Consistency
        description: >
          Does the document agree with itself? One concept keeps one name, no two
          statements contradict, and every illustration tracks the rule it
          illustrates.
        requirements:
          "the format specification is internally consistent":
            assessment: >
              No two statements contradict each other. One term denotes one
              concept throughout, and every example agrees with the rule it
              illustrates.
      verifiability:
        title: Verifiability
        description: >
          Can conformance be decided rather than argued? Each rule turns on
          something a reader can observe or test.
        requirements:
          "each rule is observable or testable":
            assessment: >
              Each rule maps to something a reader could observe or test about a
              file or implementation, so independent readers decide conformance
              the same way.
          "the format's constructs are shown with valid and invalid examples":
            assessment: >
              Constructs are shown with worked examples that include both valid
              cases and invalid counter-examples.
      extensibility:
        title: Extensibility
        description: >
          Can the format grow without breaking what exists? A stable minimal core,
          defined room to extend, and a versioning path that does not strand
          earlier files.
        requirements:
          "the format specifies its core and how it extends and evolves":
            assessment: >
              The spec names the minimal core every file must have, says how
              authors may add factors, keys, or sections, how a reader treats
              unrecognized content, and how the format versions forward.
      usability:
        title: Usability
        description: >
          Can a reader find and follow what they need? Navigability and
          readability: logical order, scannable tables, and copy-and-adapt
          examples.
        requirements:
          "the format specification is well-structured and readable":
            assessment: >
              Sections introduce a concept before dependent detail; field tables
              make structure scannable; prose is plain; and the document carries
              minimal and realistic examples.

  readme:
    title: README
    source: ./README.md
    factors:
      approachability:
        title: Approachability
        description: >
          Does the front door bring a newcomer in? A first-time reader can grasp
          what the thing is, who it is for, why the agent skill is the primary
          workflow, see it work, reach a first result, and trust that the example
          reflects reality.
        requirements:
          "the README says what QUALITY.md is and who it's for":
            assessment: >
              A first-time reader learns within the opening lines what a
              QUALITY.md file is, what problem it solves, and who it is for.
          "the README foregrounds the agent-first workflow":
            assessment: >
              The README presents `/quality`, `QUALITY.md`, and thoughtful direct
              edits as the normal user workflow, and positions the CLI as support
              tooling for validation, status, and evaluation records rather than
              as the main product surface.
          "the README shows the format and its payoff by example":
            assessment: >
              The README shows a realistic QUALITY.md excerpt and what running
              an evaluation against it produces, before reference detail.
          "the README gets a newcomer to a first result quickly":
            assessment: >
              A newcomer can install the skill and CLI, invoke `/quality`, and
              understand the shortest path to a useful first result.
          "the README reflects what the skill, CLI, and spec provide":
            assessment: >
              Every skill mode, command, flag, and capability shown matches what
              the skill, CLI, and spec provide today; planned work is marked
              planned and placeholders are marked provisional.

  quality-skill:
    title: /quality skill
    source: ./skills/quality
    factors:
      judgment-grounding:
        title: Judgment grounding
        description: >
          Does the skill keep evaluative judgment tied to the active model,
          current evidence, and the skill's functional contract? The skill is the
          primary experience, so its advice and ratings must be explainable,
          repeatable, and grounded.
        requirements:
          "the skill grounds evaluation in the active model and evidence":
            assessment: >
              Assess the runtime skill files under skills/quality/ against
              specs/skills/quality-skill/quality-skill.md. The skill resolves
              mode, model file, scope, and rigor from the user's request; reads
              the required resources before acting; uses QUALITY.md terms
              consistently; distinguishes model validity, model usefulness,
              evaluated-source quality, tooling readiness, and evaluation history; and
              stops when evidence cannot support a rating.
      mutation-safety:
        title: Mutation safety
        description: >
          Does the skill make mutation explicit, confirmed, and verifiable? Users
          should know whether the skill may change source, QUALITY.md, evaluation
          artifacts, or installed tooling before any edit happens.
        requirements:
          "the skill gates and explains mutating actions":
            assessment: >
              Assess setup, improve, and upgrade guidance in skills/quality/
              against specs/skills/quality-skill/quality-skill.md. Mutating
              workflows name the artifact class being changed, present a decision
              brief before mutation, preserve read-only wizard behavior, and
              verify the result after applying an approved change.
      cli-orchestration:
        title: CLI orchestration
        description: >
          Does the skill delegate deterministic mechanics to the CLI while
          retaining judgment in the skill? This keeps validation, status, version
          checks, scaffolding, and evaluation records predictable for agents and
          automation.
        requirements:
          "the skill delegates deterministic mechanics to qualitymd":
            assessment: >
              Assess the runtime skill prompt, mode files, and resource guidance
              against specs/skills/quality-skill/quality-skill.md and
              skills/quality/resources/cli-quick-reference.md. Mechanical steps
              use `qualitymd` commands, CLI compatibility is checked before
              CLI-dependent workflows, structured output is consumed when
              available, and the skill does not reimplement deterministic CLI
              behavior in prompt instructions.
      lifecycle-guidance:
        title: Lifecycle guidance
        description: >
          Does the skill help users move from setup to evaluation to improvement
          without losing context? The skill should make the next useful action
          obvious while keeping evaluation history and active recommendations in
          view.
        requirements:
          "the skill guides setup wizard evaluate improve and upgrade workflows":
            assessment: >
              Assess skills/quality/SKILL.md, skills/quality/modes/, and bundled
              guides against specs/skills/quality-skill/quality-skill.md. The
              skill routes ambiguous requests through wizard, inspects status and
              history when relevant, records recommendations as evaluation
              output, applies recommendations only through confirmed improve
              flows, and handles paired skill/CLI upgrade checks.

  cli:
    title: qualitymd CLI
    source: ./internal/cli
    factors:
      usability:
        title: Usability
        description: >
          Can people discover, understand, and use commands from terminal help,
          examples, output, and error messages? The CLI is support tooling, so it
          must make validation, status, evaluation record, and recovery paths
          clear without becoming the primary user workflow.
      automation-compatibility:
        title: Automation compatibility
        description: >
          Can agents, CI, and scripts drive commands without hidden human
          assumptions? The CLI is the deterministic surface that the `/quality`
          skill and automation call, so streams, exit codes, prompts, and
          structured output must remain safe to compose.
      consistency:
        title: Consistency
        description: >
          Does the command tree feel like one program? Arguments, flags, output,
          help, errors, and next actions should follow shared conventions so
          users and callers can transfer expectations between commands.
      determinism:
        title: Determinism
        description: >
          Does the same input and file state produce the same CLI result? Stable
          output, ordering, and next actions let agents and CI diff, cache, and
          assert against command behavior without flaky variation.
      maintainability:
        title: Maintainability
        description: >
          Is the Go implementation readable and idiomatic beyond what the
          deterministic check gate enforces? Naming, error handling, interface
          and API shape, concurrency discipline, value and state hygiene, doc
          comments, and test style determine how safely contributors and agents
          can change the code.
    requirements:
      "the CLI follows its functional specifications":
        factors:
          - consistency
          - determinism
          - automation-compatibility
        assessment: >
          Assess the implementation against specs/cli.md and the command
          sub-specifications under specs/cli/. Commands provide the behavior,
          flags, output, exit codes, and records those specifications require.
      "the CLI follows the project CLI design guide":
        factors:
          - usability
          - automation-compatibility
          - consistency
          - determinism
        assessment: >
          Assess command arguments, flags, help, output, errors, interactivity,
          exit codes, and next actions against docs/guides/cli-design.md.
      "the Go implementation follows the project Go style guide":
        factors:
          - maintainability
          - consistency
        assessment: >
          Assess the Go code under cmd/qualitymd/ and internal/ against
          docs/guides/go-style.md: naming, error wrapping and handling,
          interface and API shape, concurrency lifetimes, value and state
          hygiene, doc comments, and test style. Judge only the conventions the
          guide covers; the deterministic check gate (gofmt, go vet,
          staticcheck, the complexity linters, and the rest of `mise run check`)
          owns formatting, correctness, and size, so do not re-litigate those
          here.
---

# Quality model - QUALITY.md

## Overview

This model governs the QUALITY.md project itself: the open format
([`SPECIFICATION.md`](./SPECIFICATION.md)), the `/quality` agent skill that
turns the format into a working quality-management experience, the README that
introduces that experience, and the deterministic `qualitymd` CLI that supports
validation, status, and evaluation records.

Good quality here means a team or agent can make quality expectations explicit,
evaluate current evidence against them, and improve the evaluated source without losing
the intent behind the model. The governing sense of "good" is fitness for
purpose first, backed by conformance where the project defines a normative
contract: the format specification, skill functional spec, CLI specs, and design
guides are the sources of truth for the artifacts they govern.

*Unknowns* — none known.
*Open questions* — none.

*Agent-reviewed — Claude, 2026-06.*

## Scope

The deliverables are modeled as area nodes: `format-spec`, `readme`,
`quality-skill`, and `cli`. Each area carries the requirements that make its
own job assessable. The format spec declares Clarity, Consistency,
Verifiability, Extensibility, and Usability factors; the README declares
Approachability; the `/quality` skill declares Judgment Grounding, Mutation
Safety, CLI Orchestration, and Lifecycle Guidance; and the CLI declares
Usability, Automation Compatibility, Consistency, Determinism, and
Maintainability. Applicability is structural: factors apply where they are
declared and below.

The skill functional spec (`specs/skills/quality-skill/quality-skill.md`) is
binding for `/quality` behavior. The CLI functional specs (`specs/cli.md` and
`specs/cli/`) are binding for what commands do. The CLI design guide
(`docs/guides/cli-design.md`) supplies the quality expectations for how command
arguments, flags, help, output, errors, interactivity, exit codes, and next
actions should work. The Go style guide (`docs/guides/go-style.md`) supplies the
judgment-based expectations for the Go implementation under `cmd/qualitymd/` and
`internal/`, complementing the deterministic check gate that owns formatting,
correctness, and size.

Ratings should be read with a worst-of bias for contract failures: a single
unacceptable finding in the format spec, skill safety rules, or CLI automation
contract can cap its rating, even if surrounding requirements are strong.
Usability and approachability findings can compensate more often, but not when
they mislead users about the supported workflow. `Minimum` is acceptable for
early pre-1.0 work only when the body names the gap and the remaining risk is
contained.

Out of scope by design: dependencies the project does not own, including the Go
toolchain, Cobra/Fang, Agent Skills installers, and release tooling.

*Unknowns* — none known.
*Open questions* — whether to encode per-requirement rating overrides for the
most important safety and conformance requirements, given the model's stated
worst-of bias and veto behavior.

*Agent-reviewed — Claude, 2026-06.*

## Needs

- Software development teams can hold the line against the three debts that
  erode quality over time: technical debt (code drifting from where it should
  be), cognitive debt (the mental burden of understanding complex or
  under-documented systems), and intent debt (software diverging from what
  stakeholders actually need) ([Storey][triple-debt]). QUALITY.md makes
  those expectations explicit and checkable so the gaps stay visible and
  addressable.
- Format implementers can build a parser and evaluator from the specification.
- Authors can write a valid model without reverse-engineering implementation.
- Users can manage quality through their coding agent and the `/quality` skill
  without learning the CLI first.
- Authors can improve the Markdown body manually or with thoughtful AI
  assistance because it carries the model's purpose, scope, needs, risks, and the
  unknowns and open questions behind them.
- Coding agents can discover the model from `AGENTS.md` guidance and evaluate a
  root area from the model alone.
- The `/quality` skill can safely set up, inspect, evaluate, improve, and
  upgrade the quality-management workflow while making mutation and evidence
  boundaries explicit.
- Newcomers can tell from the README what QUALITY.md is, why `/quality` is the
  normal workflow, and how to reach a first result.
- Humans can understand and recover from CLI behavior through help, output, and
  errors when they need the support tooling directly.
- Agents, CI, scripts, and the `/quality` skill can drive the CLI
  deterministically without prompts or polluted streams.

*Unknowns* — none known.
*Open questions* — none.

*Agent-reviewed — Claude, 2026-06.*

## Risks

An ambiguous or incomplete format spec is the worst outcome because
implementations diverge and the format stops being portable. A skill that rates
without enough evidence, follows instructions from evaluated content, hides
mutation, or drifts from the CLI can damage trust in the whole quality workflow.
A README that drifts back to CLI-first framing turns newcomers away or teaches
the wrong path. A stale Markdown body can make future agents optimize the wrong
surface. A CLI that is inconsistent, noisy on the wrong stream, interactive by
surprise, or non-deterministic undermines the automation role the project
depends on.

*Unknowns* — none known.
*Open questions* — none.

*Agent-reviewed — Claude, 2026-06.*

## Areas and factors

### format-spec

The format spec is the contract for every reader, author, implementation, and
file. It carries the most detailed factors: Clarity, Consistency, Verifiability,
Extensibility, and Usability.

*Unknowns* — the BCP 14 references (`docs/reference/rfc2119.md` and
`docs/reference/rfc8174.md`) are not yet cited as reference standards for
requirements that depend on normative vocabulary.
*Open questions* — none.

### readme

The README is the project's front door. Approachability is scoped here because
the README's job is newcomer orientation, not format precision. It should make
the agent-first path clear: install the skill, install the CLI support tooling,
invoke `/quality`, and treat direct edits to `QUALITY.md` as a useful companion
workflow.

*Unknowns* — the README/AGENTS agent-first positioning is represented in this
model, but the readme requirements have not yet been evaluated against the
current artifacts.
*Open questions* — none.

### quality-skill

The `/quality` skill is the primary user experience for setup, evaluation,
improvement, and upgrade guidance. Its quality depends on grounded judgment,
safe mutation gates, correct CLI delegation, and lifecycle guidance that keeps
status, history, recommendations, and next actions visible.

*Unknowns* — none known.
*Open questions* — whether the skill requirements need sharper per-level criteria
once the first evaluation distinguishes adjacent ratings.

### cli

The CLI is deterministic support tooling for humans, agents, CI, and the
`/quality` skill. Its quality depends on satisfying the binding CLI functional
specs while following the CLI design guide's expectations for usable,
scriptable, consistent, and deterministic command behavior. The CLI should be
excellent at validation, status, version checks, scaffolding, and evaluation
records without becoming the default workflow users must learn first. The Go
implementation under `cmd/qualitymd/` and `internal/` is also held to the Go
style guide (`docs/guides/go-style.md`) for maintainability, which the
deterministic check gate does not assess.

*Unknowns* — none known.
*Open questions* — none.

[triple-debt]: https://arxiv.org/abs/2603.22106 "Margaret-Anne Storey, The Triple Debt of Software Development (arXiv:2603.22106)"
