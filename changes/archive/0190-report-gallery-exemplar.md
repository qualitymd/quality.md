---
type: Change Case
title: Make the report gallery an authoring exemplar
description: Expand the generated LedgerLite gallery into a best-practice QUALITY.md exemplar with a realistic evaluation report, matured-model history, and sensor-grounded assessments.
status: Done
tags: [examples, report-gallery, authoring, docs]
timestamp: 2026-07-02T00:00:00Z
---

# Make the report gallery an authoring exemplar

The [report gallery](../../examples/report-gallery/software-service/README.md) is
the repository's one browsable, end-to-end QUALITY.md example: a fictional
LedgerLite service model plus a generated evaluation run that the Mintlify docs
link to directly. Today the model is deliberately minimal — four flat areas, six
requirements, a three-paragraph disclaimer body — so it demonstrates report
_structure_ but not authoring _practice_.

This case expands the gallery into an exemplar: a QUALITY.md that visibly
follows the `/quality` skill's authoring guide family, reads like a model that
went through the skill's setup workflow and then matured through the quality
loop, and whose single generated evaluation report reads like real skill
output — including assessments that lean on guides, specs, docs, and
computational sensors.

- [Functional spec](0190-report-gallery-exemplar/spec.md) — what the exemplar
  must demonstrate.
- [Design doc](0190-report-gallery-exemplar/design.md) — how the generator is
  restructured to carry it.

## Motivation

The gallery is the first (often only) full artifact a prospective adopter
inspects, and the docs site links straight into it. A minimal model teaches the
report format but silently mis-teaches authoring: no body sections, no
unknowns/open questions, no review provenance, no model-wide factors, no
normative-artifact area, no rating overrides or veto, no not-assessed handling,
no quality changelog. An exemplar gallery makes the authoring guides concrete,
shows what the skill's setup + evaluate workflows actually produce, and
demonstrates the harness-engineering posture (feedforward guides, feedback
sensors) that the skill's assessments are written to use.

## Scope

**Covered:** the LedgerLite model and body, the generated single evaluation run
`0001-full-eval` (payloads, findings, recommendations, rankings), a synthetic
quality changelog under the example's `.quality/changelog/`, the gallery README,
and the `scripts/report-gallery` generator restructure that carries the content.

**Deferred:** a second non-software gallery example; including a fictional
source tree behind the `synthetic-source:*` references; a second evaluation run
demonstrating a closed loop across runs.

**Non-goals:** changes to the report format, evaluation record schema, CLI
behavior, or the `/quality` skill; the run stays generated through the real
`internal/evaluation` pipeline.

## Affected artifacts

**Code**

- `scripts/report-gallery/` — generator restructure (embedded content files,
  generalized payload tables, changelog emission).

**Generated example (regenerated, not hand-edited)**

- `examples/report-gallery/software-service/QUALITY.md`
- `examples/report-gallery/software-service/README.md`
- `examples/report-gallery/software-service/.quality/changelog/*`
- `examples/report-gallery/software-service/.quality/evaluations/0001-full-eval/**`

**Durable docs**

- None beyond the generated gallery README; `mintlify/index.mdx` links keep the
  stable `0001-full-eval` path and need no edit.

**Durable specs:** none — the report format specs that govern the generated
artifacts are unchanged; the gallery's content has no durable spec owner (the
gallery itself remains governed by its generator plus the `report-gallery-check`
gate).

## Status

`Done`. Implemented and archived. The generator restructure, exemplar model
and body, synthetic changelog, payload tables, and regenerated gallery landed;
the repo gates pass, regeneration is byte-stable across runs, and the emitted
`QUALITY.md` passes `qualitymd lint`. See the
[status lifecycle](../index.md#status-lifecycle).
