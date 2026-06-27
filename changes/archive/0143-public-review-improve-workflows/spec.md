---
type: Functional Specification
title: Public Review and Improve Workflows - functional spec
description: What /quality must do to make review and improve public workflows routed by focus, and align README framing around evaluate, review, and improve.
tags: [skill, workflows, docs, ux]
timestamp: 2026-06-27T00:00:00Z
---

# Public Review and Improve Workflows - functional spec

Companion to the
[Public Review and Improve Workflows](../0143-public-review-improve-workflows.md)
change case. This spec states *what* the change must do.

The key words **MUST**, **MUST NOT**, **SHOULD**, and **MAY** are to be
interpreted as described in BCP 14 when, and only when, they appear in all
capitals.

## Background / Motivation

The public `/quality` surface should match the simplest explanation of what
QUALITY.md enables: evaluate, review, and improve quality from a shared model of
what good means. The skill already has strong evaluation bounds. The next public
surface should let users enter review or improvement work without memorizing
subcommands, old internal paths, or implementation distinctions.

The user-facing concept for routing review and improve is **focus**: where the
skill should spend attention. Focus is understandable for both new and returning
users, and it gives the agent a fast way to confirm intent before doing deep
analysis or mutating anything.

## Scope

Covered: public `/quality review` and `/quality improve` invocation, focus
routing and confirmation, runtime workflow stubs, durable skill workflow specs,
README framing, runtime/spec logs, and Markdown verification.

Deferred:

- full deep workflow design for review;
- full deep workflow design for improve;
- new CLI commands or new deterministic artifact types for review/improve;
- automatic recurring review/improve runs; and
- recommendation generation changes in `evaluate`.

## Requirements

### Public framing

- `README.md` **MUST** introduce QUALITY.md primarily as a way for agents and
  teams to evaluate, review, and improve quality using a shared model of what
  good means.

  > Rationale: the README should teach the product shape before specialized
  > concepts such as quality debt, Agent Harnessability, or CLI support tooling.
  >
  > Durable spec: none — README-only durable docs change tracked in the parent
  > Affected artifacts index.

- `README.md` **MUST** present `/quality evaluate`, `/quality review`, and
  `/quality improve` as the primary value workflows.

  > Durable spec: none — README-only durable docs change tracked in the parent
  > Affected artifacts index.

- `README.md` **MUST** present `/quality setup` as the preparatory workflow for
  creating or updating the model, and `/quality update` as maintenance for the
  skill/CLI pair.

  > Durable spec: none — README-only durable docs change tracked in the parent
  > Affected artifacts index.

- `README.md` **SHOULD** simplify the quality-loop explanation to: model what
  good means, evaluate the work, review the evidence and model fit, then improve
  the work or the model.

  > Durable spec: none — README-only durable docs change tracked in the parent
  > Affected artifacts index.

- README usage examples **MUST** include `/quality review [focus]` and
  `/quality improve [focus]`, and **SHOULD** describe focus examples as the
  latest evaluation, `QUALITY.md` model, and a specific quality concern.

  > Durable spec: none — README-only durable docs change tracked in the parent
  > Affected artifacts index.

### Invocation and routing

- `/quality review` and `/quality improve` **MUST** be public workflows in the
  root skill invocation contract.

  > Durable spec: modify `specs/skills/quality-skill/quality-skill.md` — add
  > `review` and `improve` to the public invocation contract and workflow
  > dispatch list; add `specs/skills/quality-skill/workflows/review.md` and
  > `specs/skills/quality-skill/workflows/improve.md` — create the workflow
  > behavioral component specs.

- The root skill dispatcher **MUST** read the matching workflow file before
  executing `review` or `improve`.

  > Durable spec: modify `specs/skills/quality-skill/quality-skill.md` — require
  > the dispatcher to read `workflows/review.md` and `workflows/improve.md`
  > before executing those workflows.

- `review` and `improve` **MUST** use **focus** as the user-facing routing term
  for what the workflow should attend to.

  > Rationale: focus names the user's attention target without exposing an
  > implementation taxonomy such as lane, mode, or source.
  >
  > Durable spec: modify `specs/skills/quality-skill/quality-skill.md` — define
  > focus as the shared review/improve routing concept; add
  > `specs/skills/quality-skill/workflows/review.md` and
  > `specs/skills/quality-skill/workflows/improve.md` — require focus-driven
  > workflow routing in each workflow spec.

- `review` **MUST** support at least these focus values: latest or selected
  Evaluation result, the `QUALITY.md` model, and a specific quality concern.

  > Durable spec: add `specs/skills/quality-skill/workflows/review.md` — define
  > the minimum review focus set.

- `improve` **MUST** support at least these focus values: Evaluation result or
  finding/candidate action, the `QUALITY.md` model, a specific quality concern,
  and an existing recommendation artifact when one is present.

  > Durable spec: add `specs/skills/quality-skill/workflows/improve.md` — define
  > the minimum improve focus set.

- When the user supplies an explicit focus, the workflow **MUST** use it unless
  the target is impossible or unsafe; when focus is absent or ambiguous, the
  workflow **MUST** infer the likely focus from user text and local lifecycle
  state before asking.

  > Durable spec: add `specs/skills/quality-skill/workflows/review.md` and
  > `specs/skills/quality-skill/workflows/improve.md` — require explicit-focus
  > handling and absent/ambiguous-focus inference.

- When inference is not strong enough to proceed, the workflow **MUST** ask a
  single-select closed-choice question with the recommended focus first and an
  explicit shortest answer path.

  > Durable spec: add `specs/skills/quality-skill/workflows/review.md` and
  > `specs/skills/quality-skill/workflows/improve.md` — define the fallback
  > focus-choice prompt when inference is insufficient.

### Review workflow stub

- `review` **MUST** emit a public run frame as its first output, before tool
  calls, and the frame **MUST** include the resolved or provisional focus.

  > Durable spec: add `specs/skills/quality-skill/workflows/review.md` — require
  > the review run frame and its focus field.

- `review` **MUST** be read-only by default. It **MUST NOT** edit evaluated
  source, edit `QUALITY.md`, write evaluation records, write the quality log,
  create external issues, or update tooling.

  > Durable spec: add `specs/skills/quality-skill/workflows/review.md` — define
  > review's default read-only mutation boundary.

- `review` **MUST** quickly confirm or correct the inferred focus before deep
  review work. For the initial stub, it **MAY** stop after focus confirmation and
  a clear preview of what the deeper review would inspect.

  > Rationale: making the workflow public now should improve routing without
  > pretending the deeper review design is complete.
  >
  > Durable spec: add `specs/skills/quality-skill/workflows/review.md` — require
  > early focus confirmation and allow the initial stub to stop after preview.

- A review of an Evaluation result **SHOULD** inspect the latest reportable run
  when no run is named, summarize its available rating/findings/evidence-limit
  surface, and identify likely next improve focuses without mutating anything.

  > Durable spec: add `specs/skills/quality-skill/workflows/review.md` — define
  > the Evaluation-result review stub behavior.

- A model review **SHOULD** route through existing authoring and
  top-10-check guidance, inspect `QUALITY.md` for usefulness, clarity, coverage,
  assessability, and stale assumptions, and report suggested improve focuses
  without changing the model.

  > Durable spec: add `specs/skills/quality-skill/workflows/review.md` — define
  > the model review stub behavior and its use of existing model-review guidance.

- An ad hoc concern review **SHOULD** inspect the named concern against the
  model and available project context, and should recommend either a scoped
  evaluation, model improvement, or work improvement as the next action.

  > Durable spec: add `specs/skills/quality-skill/workflows/review.md` — define
  > the ad hoc concern review stub behavior.

### Improve workflow stub

- `improve` **MUST** emit a public run frame as its first output, before tool
  calls, and the frame **MUST** include the resolved or provisional focus and a
  mutation surface of `read-only until confirmed` when the mutation is not yet
  resolved.

  > Durable spec: add `specs/skills/quality-skill/workflows/improve.md` —
  > require the improve run frame, focus field, and unresolved mutation boundary.

- `improve` **MUST** confirm both focus and mutation surface before editing
  evaluated source, editing `QUALITY.md`, writing the quality log, creating an
  external issue, or updating tooling.

  > Durable spec: add `specs/skills/quality-skill/workflows/improve.md` —
  > require focus and mutation-surface confirmation before any improve mutation.

- `improve` **MUST** use the existing decision-brief or review-gate shapes for
  mutation confirmation, with `y`/`n` as the visible shortest responses for true
  binary mutation gates.

  > Durable spec: add `specs/skills/quality-skill/workflows/improve.md` —
  > require improve mutation gates to reuse existing review-gate and decision-brief
  > shapes.

- For model-focused improvement, `improve` **MUST** delegate to the existing
  direct model-authoring guidance after focus and mutation surface are
  confirmed.

  > Durable spec: modify `specs/skills/quality-skill/quality-skill.md` — route
  > model-improvement intent through public `/quality improve` while preserving
  > the direct model-authoring contract; add
  > `specs/skills/quality-skill/workflows/improve.md` — require the workflow to
  > delegate model-focused improvement to that existing contract.

- For Evaluation-result or recommendation-focused improvement, `improve` **MAY**
  route to existing recommendation-follow-up guidance when a compatible
  recommendation artifact exists, and **MUST** otherwise avoid inventing a
  recommendation-generation step.

  > Durable spec: modify `specs/skills/quality-skill/quality-skill.md` — route
  > recommendation and Evaluation-result improvement through public
  > `/quality improve`; modify
  > `specs/skills/quality-skill/recommendation-follow-up.md` — keep the existing
  > recommendation follow-up contract as the route used when compatible
  > recommendation artifacts exist; add
  > `specs/skills/quality-skill/workflows/improve.md` — prohibit inventing
  > recommendation generation when no compatible artifact exists.

- For ad hoc concern improvement, `improve` **MUST** identify whether the likely
  mutation surface is evaluated source, `QUALITY.md`, external issue handoff, or
  no mutation yet, then confirm that surface before acting.

  > Durable spec: add `specs/skills/quality-skill/workflows/improve.md` — define
  > ad hoc concern mutation-surface inference and confirmation.

- The initial `improve` stub **MAY** stop after focus and mutation-surface
  confirmation with a clear explanation that deeper workflow design is deferred,
  unless the requested work safely routes to an existing confirmed model-authoring
  or recommendation-follow-up path.

  > Durable spec: add `specs/skills/quality-skill/workflows/improve.md` — allow
  > the initial improve stub to stop at a clear deferred-design boundary unless
  > an existing confirmed route applies.

### Backward compatibility and boundaries

- The change **MUST NOT** add legacy aliases, fallback readers, or CLI shims for
  review/improve.

  > Durable spec: none — this follows the repository early-alpha compatibility
  > policy and creates no new durable spec delta beyond the public workflow specs.

- `evaluate` **MUST** remain the only workflow that creates numbered Evaluation
  records and reports.

  > Durable spec: modify `specs/skills/quality-skill/quality-skill.md` — preserve
  > `evaluate` as the only workflow that creates Evaluation records and reports.

- `setup` and `update` **MUST** remain public support workflows, not part of the
  primary evaluate/review/improve value-prop trio.

  > Durable spec: modify `specs/skills/quality-skill/quality-skill.md` — frame
  > `setup` and `update` as support workflows relative to the primary
  > evaluate/review/improve value proposition.

- The change **MUST NOT** alter `SPECIFICATION.md`, Evaluation data schema, lint
  rules, model parsing, scaffold output, or install mechanics.

  > Durable spec: none — this requirement explicitly excludes non-skill durable
  > specs and implementation surfaces.

## Durable spec changes

### To add

- `specs/skills/quality-skill/workflows/review.md` — create the review workflow
  behavioral component spec, including focus values, run frame, read-only
  boundary, confirmation behavior, and initial stub behaviors (per the public
  workflow, focus, and review-stub requirements above).
- `specs/skills/quality-skill/workflows/improve.md` — create the improve
  workflow behavioral component spec, including focus values, run frame,
  mutation-surface confirmation, delegation to existing routes, and initial stub
  behaviors (per the public workflow, focus, and improve-stub requirements
  above).

### To modify

- `specs/skills/quality-skill/quality-skill.md` — add public `review` and
  `improve` invocation, focus routing, workflow dispatch, support-workflow
  framing for `setup`/`update`, Evaluation-record boundaries, and improve
  routing to existing direct authoring and recommendation-follow-up contracts
  (per the invocation, improve delegation, and boundary requirements above).
- `specs/skills/quality-skill/workflows/index.md` — list review and improve as
  public workflows alongside setup, evaluate, and update (per the public workflow
  requirement above).
- `specs/skills/quality-skill/recommendation-follow-up.md` — describe
  recommendation follow-up as the existing implementation route used by
  `/quality improve` when compatible recommendation artifacts exist, rather than
  as a separate public workflow (per the Evaluation-result or
  recommendation-focused improve requirement above).

### To rename

None.

### To delete

None.

## Verification

- `mise run fmt-md-check`
- Search active docs/specs/runtime skill files for stale claims that `setup`,
  `evaluate`, and `update` are the only public workflows.
- Search active README and skill guidance for internal-only wording such as
  "review model" or "recommendation follow-up" where public `/quality review` or
  `/quality improve` should now be named.
