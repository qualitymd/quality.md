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
targets:
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
          what the thing is, who it is for, see it work, reach a first result,
          and trust that the example reflects reality.
        requirements:
          "the README says what QUALITY.md is and who it's for":
            assessment: >
              A first-time reader learns within the opening lines what a
              QUALITY.md file is, what problem it solves, and who it is for.
          "the README shows the format and its payoff by example":
            assessment: >
              The README shows a realistic QUALITY.md excerpt and what running
              qualitymd against it produces, before reference detail.
          "the README gets a newcomer to a first result quickly":
            assessment: >
              A newcomer can copy a short install-then-one-command sequence and
              see a real result, with representative output and CI exit-code
              behavior where relevant.
          "the README reflects what the CLI and spec actually provide":
            assessment: >
              Every command, flag, and capability shown matches what the CLI
              provides today; planned work is marked planned and placeholders
              are marked provisional.

  cli:
    title: qualitymd CLI
    source: ./internal/cli
    factors:
      usability:
        title: Usability
        description: >
          Can people discover, understand, and use commands from terminal help,
          examples, output, and error messages? The CLI is a working interface,
          so it must make the common path clear and recovery from mistakes
          practical.
      automation-compatibility:
        title: Automation compatibility
        description: >
          Can agents, CI, and scripts drive commands without hidden human
          assumptions? The CLI is the deterministic surface skills and
          automation call, so streams, exit codes, prompts, and structured output
          must remain safe to compose.
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
---

# Quality model - QUALITY.md

## Overview

This model governs the QUALITY.md project itself: the format
([`SPECIFICATION.md`](./SPECIFICATION.md)), the README that introduces it, and
the deterministic `qualitymd` CLI that implements the mechanical surface. At
this pre-1.0 stage quality rests on the project's design, so the model covers the
maturity of the format specification, the README's job of orienting newcomers,
and the CLI's ability to serve both humans and automation.

## Scope

The deliverables are modeled as target nodes: `format-spec`, `readme`, and
`cli`. Each deliverable carries a direct "does it do its job" requirement where
appropriate. The format spec declares Clarity, Consistency, Verifiability,
Extensibility, and Usability factors; the README declares Approachability; and
the CLI declares Usability, Automation Compatibility, Consistency, and
Determinism. Applicability is structural: factors apply where they are declared
and below.

The CLI functional specs (`specs/cli.md` and `specs/cli/`) are binding for what
commands do. The CLI design guide (`docs/guides/cli-design.md`) supplies the
quality expectations for how command arguments, flags, help, output, errors,
interactivity, exit codes, and next actions should work.

Out of scope by design: dependencies the project does not own, including the Go
toolchain, Cobra/Fang, and release tooling.

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
- Coding agents can evaluate a subject from the model alone.
- Newcomers can tell from the README what QUALITY.md is and reach a first result.
- Humans can understand and recover from CLI behavior through help, output, and
  errors.
- Agents, CI, and scripts can drive the CLI deterministically without prompts or
  polluted streams.

## Risks

An ambiguous or incomplete format spec is the worst outcome because
implementations diverge and the format stops being portable. A README that
overstates what exists turns newcomers away or misleads them. A CLI that is
inconsistent, noisy on the wrong stream, interactive by surprise, or
non-deterministic undermines the automation role the project depends on.

## Targets and factors

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
the README's job is newcomer orientation, not format precision.

### cli

The CLI is the deterministic tool surface for humans, agents, and CI. Its
quality depends on satisfying the binding CLI functional specs while following
the CLI design guide's expectations for usable, scriptable, consistent, and
deterministic command behavior.

*Unknowns* — none known.
*Open questions* — none.

[triple-debt]: https://arxiv.org/abs/2603.22106 "Margaret-Anne Storey, The Triple Debt of Software Development (arXiv:2603.22106)"
