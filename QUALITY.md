---
title: QUALITY.md
ratingScale:
  - level: outstanding
    title: Outstanding
    criterion: "Exceeds the requirement; satisfies it with margin to spare."
  - level: target
    title: Target
    criterion: "Satisfies the requirement."
  - level: minimum
    title: Minimum
    criterion: "Falls short of the goal but holds the acceptable floor."
  - level: unacceptable
    title: Unacceptable
    criterion: "Falls below the acceptable floor."
targets:
  format-spec:
    source: ./SPECIFICATION.md
    requirements:
      "the format specification is complete":
        assessment: >
          Every frontmatter field and recommended body section has its shape,
          allowed values, requiredness, cardinality, and any default; the spec
          states how a conforming reader treats malformed or omitted content and
          edge cases. An implementer could build a parser, and an author write a
          valid file, from the spec alone.
    factors:
      clarity:
        description: >
          Can each rule be read in only one way? A spec governs independent
          implementations that never confer, so every obligation must land with
          one settled meaning and force.
        requirements:
          "the format specification admits a single interpretation":
            assessment: >
              Each rule admits one reading. Obligations use a consistent
              must/should/may vocabulary and no normative statement leans on a
              vague quantifier without a stated bound.
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
    source: ./README.md
    factors:
      approachability:
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
---

# Quality model - QUALITY.md

## Overview

This model governs the QUALITY.md project itself: the format
([`SPECIFICATION.md`](./SPECIFICATION.md)) and the README that introduces it. At
this pre-1.0 stage quality rests on the project's design, so the model covers the
maturity of the format specification and the README's job of orienting newcomers.

## Scope

The deliverables are modeled as target nodes: `format-spec` and `readme`. Each
deliverable carries a direct "does it do its job" requirement where appropriate.
Only the format spec declares Clarity, Consistency, Verifiability, Extensibility,
and Usability factors; only the README declares Approachability. Applicability is
structural: factors apply where they are declared and below.

Out of scope by design: dependencies the project does not own, including the Go
toolchain, Cobra/Fang, and release tooling.

## Needs

- Format implementers can build a parser and evaluator from the specification.
- Authors can write a valid model without reverse-engineering implementation.
- Coding agents can evaluate a subject from the model alone.
- Newcomers can tell from the README what QUALITY.md is and reach a first result.

## Risks

An ambiguous or incomplete format spec is the worst outcome because
implementations diverge and the format stops being portable. A README that
overstates what exists turns newcomers away or misleads them.

## Targets and factors

### format-spec

The format spec is the contract for every reader, author, implementation, and
file. It carries the most detailed factors: Clarity, Consistency, Verifiability,
Extensibility, and Usability.

### readme

The README is the project's front door. Approachability is scoped here because
the README's job is newcomer orientation, not format precision.

## Known gaps

- The CLI and skills are not modeled. Their specifications have been removed from
  the repo, and the implementation's runtime quality — reliability, performance,
  packaging, and test coverage — is deferred while the implementation is nascent.
- No structural self-lint requirement yet. Once `qualitymd lint` ships, add a
  direct apex requirement asking that this model lint cleanly.
