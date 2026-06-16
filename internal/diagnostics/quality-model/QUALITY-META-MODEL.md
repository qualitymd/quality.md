---
title: Quality meta-model
source: ./QUALITY.md
ratingScale:
  - level: outstanding
    title: Outstanding
    criterion: "Exceeds the diagnostic requirement; satisfies it with margin to spare."
  - level: target
    title: Target
    criterion: "Satisfies the diagnostic requirement; no material gaps."
  - level: minimum
    title: Minimum
    criterion: "Satisfies the core of the requirement but falls short of the goal; minor or scoped gaps remain at the acceptable floor."
  - level: unacceptable
    title: Unacceptable
    criterion: "Does not satisfy the diagnostic requirement; falls below the acceptable floor."
factors:
  functionality:
    description: >
      Does the quality model do its job: produce a meaningful, trustworthy verdict
      about whether its subject meets the quality expectations it declares?
    requirements:
      "model fulfills its declared functional purpose":
        assessment: >
          Taken as a whole, the model serves the purpose it declares. Evaluating it
          yields a meaningful verdict about whether the subject meets the quality
          expectations the model sets, supporting a real accept, reject, or
          improvement decision.
      "requirements collectively realize the model's declared needs":
        assessment: >
          Every need and risk the model declares is guarded by at least one
          requirement whose failure would surface a violation of it, and no
          requirement stands without a need behind it.
      "the model yields correct, trustworthy verdicts":
        assessment: >
          For a subject that does or does not satisfy a requirement, the model's
          assessments and rating scale produce the verdict a knowledgeable reviewer
          would reach.
      "each factor is individually well-formed":
        assessment: >
          Each factor is a quality attribute, relevant to the subject, distinct from
          sibling factors, grounded in a recognizable quality vocabulary where one
          fits, and operationalized by requirements. A factor that names a component,
          feature, or activity rather than a quality attribute is unacceptable.
      "the factor set is well-formed as a whole":
        assessment: >
          The factor set is complete, non-overlapping, coherent, and appropriately
          decomposed across the target tree. Concerns intentionally left out are
          recorded in Scope or Known gaps rather than silently omitted.
      "each requirement is individually well-formed":
        assessment: >
          Each requirement is necessary, appropriate, unambiguous, complete,
          concise, singular, feasible, verifiable, correct, and conforming to the
          model's style. It carries one real assessment and localizes any shortfall.
      "the requirement set is well-formed as a whole":
        assessment: >
          The requirement set is complete, consistent, lean, feasible,
          comprehensible, and able to be validated; it has no unresolved
          placeholders, dangling references, conflicts, or redundant checks.
      "the Overview body section frames the subject":
        assessment: >
          Where present, the Overview section establishes what the subject is, who
          depends on it, and what "good" means for it, without doing the Scope
          section's boundary work.
      "the Scope body section draws the model's boundary":
        assessment: >
          Where present, the Scope section states what the model covers and what it
          deliberately leaves out. Out-of-scope items are framed as exclusions by
          design, not deferred in-scope gaps.
      "the Needs body section states what matters and to whom":
        assessment: >
          Where present, the Needs section states stakeholder outcomes in plain
          language: who depends on what, and how they suffer if it is unmet or met
          poorly.
      "the Risks body section states what failure costs and to whom":
        assessment: >
          Where present, the Risks section says what goes wrong, for whom, and with
          what relative severity if a need is not met.
      "the Factors body section explains every declared factor":
        assessment: >
          The body gives every factor declared in the target tree matching prose
          that characterizes the quality attribute itself, why it matters, and how
          it differs from siblings, without merely restating requirements.
      "the Known gaps body section records deferred concerns with reasons":
        assessment: >
          Where the model defers quality concerns inside its scope, the Known gaps
          section records each with a reason. Concerns outside the model's remit
          belong in Scope instead.
      "the Markdown body earns its length":
        assessment: >
          Across its sections, the Markdown body supplies subject-specific reasoning
          a reader could not supply themselves and avoids padding, generic advice,
          or obvious narration.
      "the model passes structural lint":
        assessment: >
          Running `qualitymd lint` on the model reports no error-level findings:
          the target tree parses, sources and assessment references resolve, scoped
          factor references are valid, and the rating scale is well-shaped.
      "the model correctly applies the QUALITY.md format spec":
        assessment: >
          The model is a correct application of the QUALITY.md format, not merely a
          file that parses. Target nodes, scoped factors, direct and lensed
          requirements, assessment fields, source bindings, secondary factors,
          and rating criteria are used for their intended purposes.
      "the model includes everything the QUALITY.md spec prescribes":
        assessment: >
          The model includes the structured target tree and the recommended
          Markdown rationale needed to make it useful. Required assessment fields
          are present, recommended body sections that apply are present, and every
          declared target/factor has enough prose context to be understood.
  usability:
    description: >
      Can the agents and developers who consume, maintain, and act on the quality
      model use it without guesswork?
    requirements:
      "an agent can interpret each requirement and assessment unambiguously":
        assessment: >
          A coding agent arrives at a single reading of each requirement and
          assessment. Statements supply the scope, terms, and instruction needed to
          render the verdict.
      "an agent can execute every assessment the model declares":
        assessment: >
          Each assessment can be carried out by an agent as written. The source or
          referenced artifacts it must read can be located and retrieved, and no
          assessment requires unavailable access or capability.
      "an agent can use the model as working context":
        assessment: >
          Beyond formal evaluation, an agent can read the model and build an
          accurate picture of the subject's quality expectations to guide work.
      "a developer can understand the model's intent":
        assessment: >
          A developer can grasp what the model expects and why from the Markdown
          body without reverse-engineering the frontmatter.
      "a developer can act on the model's verdicts":
        assessment: >
          A result points clearly enough at what is deficient, and where, that a
          developer knows what to change.
      "a developer can extend and maintain the model":
        assessment: >
          A developer can add, refine, or remove Targets, factors, and
          requirements through the format's extension points while keeping the model
          coherent.
---

# Quality meta model

## Overview

This is the CLI's built-in diagnostic model for evaluating a project's own
`QUALITY.md`. Its apex target source is that file. Keeping these criteria in a
versioned, inspectable model rather than an opaque evaluator prompt makes the bar
for a good quality model reviewable.

The model has two factors. **Functionality** asks whether the model produces a
meaningful, trustworthy verdict about its subject. **Usability** asks whether the
agents and developers who consume and maintain it can use it.

## Scope

The subject is the `QUALITY.md` artifact: its target tree, factors,
requirements, rating criteria, references, and Markdown body. The quality of the
software governed by that file is out of scope; that belongs to evaluating the
project's model against the project. The CLI and skills that run evaluations are
also out of scope.

## Needs

- A project can tell whether its `QUALITY.md` is a useful quality model, not only
  a syntactically valid file.
- An agent can produce concrete defects and coverage gaps that improve the model
  before using it to evaluate the subject.

## Risks

A false acceptance is the worst outcome because a deficient model then governs
the subject and gives the project false confidence. A false rejection wastes
effort. A verdict that does not localize the defect leaves authors unable to act.

## Targets and factors

### apex target

The apex target is the user's `QUALITY.md`. All diagnostics are assessed against
that artifact and any referenced assessment or source material needed to judge it.

### Functionality

Functionality covers purpose fit, model well-formedness, body-section quality,
format conformance, and format completeness. These concerns are expressed as
requirements under the Functionality factor because decomposition now happens
through Targets.

### Usability

Usability covers both agent usability and developer usability as requirements
under one factor. The distinction is retained in requirement wording rather than
as nested factors.

## Diagnostic coverage checklist

The schema migration preserves the previous diagnostic requirement set:

- model fulfills its declared functional purpose
- requirements collectively realize the model's declared needs
- the model yields correct, trustworthy verdicts
- each factor is individually well-formed
- the factor set is well-formed as a whole
- each requirement is individually well-formed
- the requirement set is well-formed as a whole
- the Overview / Scope / Needs / Risks / Factors / Known gaps body sections do
  their prescribed jobs
- the Markdown body earns its length
- the model passes structural lint
- the model correctly applies the QUALITY.md format spec
- the model includes everything the QUALITY.md spec prescribes
- an agent can interpret, execute, and use the model as working context
- a developer can understand, act on, extend, and maintain the model

## Known gaps

- The model's own maintainability is judged through Developer usability rather
  than as a standalone factor.
- The structural-lint floor depends on `qualitymd lint`, which is specified but
  not yet implemented.
