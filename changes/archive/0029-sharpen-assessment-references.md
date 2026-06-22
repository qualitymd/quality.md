---
type: Change Case
title: Sharpen assessment references and traceability
description: Clarify that a requirement's assessment is stated inline or as a reference to another entity, and make the model's traceability graph an explicit authoring concern.
status: Done
tags: [specification, terminology, guide]
timestamp: 2026-06-19T00:00:00Z
---

# Sharpen assessment references and traceability

A **Change Case** capturing the *why* and *status*; the detail lives in its
[functional spec](0029-sharpen-assessment-references/spec.md). No design doc — the
change needs no separate design discussion.

## Motivation

A requirement's `assessment` can either state the means of judging inline or
point at where those means are defined — a spec, a guide, a checklist. The format
treats both as the same scalar, but the durable artifacts never name the
distinction, so authors copy criteria into requirements (where they drift from
their origin) and slice one assessable claim into several per-factor requirements
that each re-cite the same artifact.

The deeper point is that the entity an assessment references is usually itself a
target in the model — a spec is a target *and* the criteria the code that
implements it is judged against; a how-to guide is a target *and* the criteria
for the docs it governs. The connections between these targets — which entity is
the criteria for which — form a traceability graph that is among the most
valuable things a model records, yet nothing in the guidance helps an author make
those edges visible.

The word **source** makes this worse: the model already binds it to
`Target.source` (the entities a target evaluates), so calling the referenced
criteria a second "source" overloads the term. This change reserves "source" for
`Target.source`, frames an assessment as **inline or a reference**, and uses
**reference** for the requirement→entity edge — extending the "reference"
terminology that [0028](0028-require-characterized-requirements.md) settled for
requirement→factor relationships.

## Scope

Covered: the format spec's description of `assessment`; the authoring guide's
treatment of assessments, the target/assessment-reference duality, the
traceability graph, and requirement granularity; and the scaffold's assessment
guidance.

Deferred / non-goals: no format **schema** change (an `assessment` stays a single
scalar) and no new lint rule — the inline-vs-reference distinction is not
mechanically enforced. No first-class machine-readable assessment→target link
type; traceability rides on shared canonical references. Reworking this repo's
own `QUALITY.md` model to model the doc/spec/code graph is separate work.

Implementation begins only after the change advances to `In-Progress`.

## Affected specs & docs

- [x] [`SPECIFICATION.md`](../../SPECIFICATION.md) — describe an `assessment` as
      stated inline or as a reference to an entity that defines the means, and add a
      non-normative note that a referenced entity may itself be a Target reachable by
      its `source`. (Durable spec; detailed in the functional spec's
      [Durable spec changes](0029-sharpen-assessment-references/spec.md#durable-spec-changes).)
- [x] [`skills/quality/guides/authoring.md`](../../skills/quality/guides/authoring.md)
      — reserve "source" for `Target.source`; add the inline/reference framing, the
      "Make the traceability graph visible" job, the entity gloss, the
      target/assessment-reference duality job, and the "Split by assessable claim,
      not by factor" job; rename "Reference assessment sources; don't copy them".
- [x] [`internal/scaffold/skeleton.md`](../../internal/scaffold/skeleton.md) — use
      "reference" rather than "defer to" for the assessment, and note a referenced
      entity can be a target in its own right.
- Reviewed, no change:
  [`specs/skills/quality-skill/guides/authoring-md.md`](../../specs/skills/quality-skill/guides/authoring-md.md)
  — its conformance and single-level "Working with…" structure already admit
  these as authoring jobs; the guide's conformance duty carries the new
  `SPECIFICATION.md` wording forward.

## Children

- [Functional spec](0029-sharpen-assessment-references/spec.md)

## Status

`Done`. No design doc was needed; the change advanced from `Draft` straight to
implementation. All durable
artifacts in **Affected specs & docs** are synced; the repo test suite, vet, and
Markdown format checks pass. (Two unrelated pre-existing `golangci-lint` findings
in `internal/evaluation/report.go` and `internal/cli/style.go` are untouched by
this change.)
