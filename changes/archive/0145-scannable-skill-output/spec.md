---
type: Functional Specification
title: Scannable Skill Output - functional spec
description: What /quality runtime guidance must do to adopt labeled, five-second-scan output shapes.
tags: [docs, skill, ux, workflows]
timestamp: 2026-06-27T00:00:00Z
---

# Scannable Skill Output - functional spec

Companion to the
[Scannable Skill Output](../0145-scannable-skill-output.md) change case. This
spec states _what_ the change must do.

The key words **MUST**, **MUST NOT**, **SHOULD**, and **MAY** are to be
interpreted as described in BCP 14 when, and only when, they appear in all
capitals.

## Background / Motivation

The shared agent-mediated UX guide now defines scannable output as a first-class
interaction property: results, importance, boundaries, and answer paths should be
visible in a five-second scan. Several `/quality` runtime templates still rely on
paragraph prose or field lists without concrete output shapes. This change makes
the runtime skill and durable contracts absorb that guidance so user-facing
outputs consistently use labels, bullets, numbered options, and explicit answer
paths.

## Scope

Covered: runtime `/quality` prompt guidance, setup/evaluate/review/improve/update
workflow output templates, recommendation follow-up result templates, top-10
check and getting-started guide summaries, the shared UX guide's new scannability
section, matching durable specs, logs, changelog, and Markdown verification.

Deferred:

- native UI rendering details;
- CLI human-output formatting;
- generated Evaluation report rendering;
- new workflow surfaces; and
- changes to Evaluation data schemas.

## Requirements

### Shared UX guidance

- `docs/guides/agent-mediated-ux.md` **MUST** define scannable output as
  five-second-scan output where result, importance, and next action are visible
  without reading every sentence.

  > Rationale: the guide is the shared contract; runtime skill templates should
  > not be the only place this design rule exists.
  >
  > Durable spec: none - this is durable docs guidance, not a `specs/` bundle
  > contract.

- `docs/guides/agent-mediated-ux.md` **MUST** tell authors to prefer labels,
  bullets, or numbered lists over dense paragraphs when output contains multiple
  independent facts.

  > Durable spec: none - this is durable docs guidance.

### Runtime interaction contract

- `skills/quality/SKILL.md` **MUST** carry the scannable-output rule in the
  global user interaction contract: multi-fact user-facing outputs use labels,
  bullets, or numbered lists, and answer paths appear on their own line or
  label.

  > Durable spec: modify `specs/skills/quality-skill/quality-skill.md` - add the
  > scannable-output rule to the user interaction contract.

- Direct model authoring in `skills/quality/SKILL.md` **MUST** show its
  pre-mutation checkpoint as labeled fields for planned edit, why, approach,
  boundary, log, and answer path when those fields apply.

  > Durable spec: modify `specs/skills/quality-skill/quality-skill.md` - require
  > labeled direct-authoring review gates.

### Workflow output templates

- `skills/quality/workflows/setup.md` **MUST** keep its first output scannable by
  separating the run frame from labeled `Why`, `Plan`, and `Boundary` blocks.

  > Durable spec: modify `specs/skills/quality-skill/workflows/setup.md` - define
  > the labeled setup opening shape.

- `skills/quality/workflows/setup.md` **MUST** provide a labeled final setup
  review recap template that surfaces root, domain, lifecycle, risk, rating
  scale, human context, open gaps, and answer path before writing.

  > Durable spec: modify `specs/skills/quality-skill/workflows/setup.md` - define
  > the final review recap shape.

- `skills/quality/workflows/evaluate.md` **MUST** provide a concrete labeled
  evaluation closeout template with rating, scope, evidence basis, known
  limitations, changed artifacts, not-done boundary, reports, and next action.

  > Durable spec: modify `specs/skills/quality-skill/workflows/evaluate.md` and
  > `specs/skills/quality-skill/reporting.md` - define the labeled evaluation
  > closeout shape.

- `skills/quality/workflows/review.md` **MUST** provide a labeled read-only
  closeout template with reviewed subject, signal, evidence limits, recommended
  next action, alternatives, and not-changed boundary.

  > Durable spec: modify `specs/skills/quality-skill/workflows/review.md` -
  > define the labeled review closeout shape.

- `skills/quality/workflows/improve.md` **MUST** provide a labeled improve
  closeout template with focus, changed artifacts, verification, remaining
  limits, not-changed boundary, and next action.

  > Durable spec: modify `specs/skills/quality-skill/workflows/improve.md` -
  > define the labeled improve closeout shape.

- `skills/quality/workflows/update.md` **MUST** provide a labeled update closeout
  template with inspected versions, applied actions, verification, restart/reload
  note, not-changed boundary, and next action.

  > Durable spec: modify `specs/skills/quality-skill/workflows/update.md` -
  > define the labeled update closeout shape.

### Runtime guide output templates

- `skills/quality/guides/recommendation-follow-up.md` **MUST** normalize the
  recommendation result closeout to labeled fields for recommendation, applied
  option, changed artifacts, verification, rating movement, remaining gaps, not
  done, and next action.

  > Durable spec: modify
  > `specs/skills/quality-skill/guides/recommendation-follow-up-md.md` and
  > `specs/skills/quality-skill/recommendation-follow-up.md` - define the labeled
  > recommendation result shape.

- `skills/quality/guides/top-10-quality-md-checks.md` **MUST** provide a labeled
  model-review signal template with lifecycle state, model-usefulness findings,
  recommended next action, alternatives, and not-assessed boundary.

  > Durable spec: modify
  > `specs/skills/quality-skill/guides/top-10-quality-md-checks-md.md` - define
  > the labeled top-10 summary shape.

- `skills/quality/guides/getting-started.md` **MUST** make its next-workflow
  guidance a labeled recommendation plus numbered options and explicit answer
  path.

  > Durable spec: modify
  > `specs/skills/quality-skill/guides/getting-started-md.md` - define the
  > labeled next-workflow choice shape.

## Durable spec changes

### To add

None

### To modify

- `specs/skills/quality-skill/quality-skill.md` - add scannable-output contract
  and labeled direct-authoring checkpoint requirements.
- `specs/skills/quality-skill/workflows/setup.md` - add labeled setup opening and
  final review recap shapes.
- `specs/skills/quality-skill/workflows/evaluate.md` - add labeled evaluation
  closeout requirements.
- `specs/skills/quality-skill/workflows/review.md` - add labeled review closeout
  requirements.
- `specs/skills/quality-skill/workflows/improve.md` - add labeled improve
  closeout requirements.
- `specs/skills/quality-skill/workflows/update.md` - add labeled update closeout
  requirements.
- `specs/skills/quality-skill/guides/recommendation-follow-up-md.md` - add
  labeled recommendation result requirements.
- `specs/skills/quality-skill/guides/top-10-quality-md-checks-md.md` - add
  labeled model-review signal requirements.
- `specs/skills/quality-skill/guides/getting-started-md.md` - add labeled
  next-workflow choice requirements.
- `specs/skills/quality-skill/recommendation-follow-up.md` - align
  recommendation follow-up behavior with the guide contract.
- `specs/skills/quality-skill/reporting.md` - align evaluation closeout behavior
  with the workflow contract.

### To rename

None

### To delete

None

## Verification

- Source inspection **MUST** show concrete labeled templates in every runtime
  file named in the requirements.
- Source inspection **MUST** show matching durable spec requirements.
- Markdown formatting checks **SHOULD** pass for touched Markdown files.
- The repository check **SHOULD** pass before commit.
