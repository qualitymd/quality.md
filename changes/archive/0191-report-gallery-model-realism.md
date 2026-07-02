---
type: Change Case
title: Deepen the report gallery exemplar model
description: Expand the LedgerLite gallery to a realistically scaled, guide-conformant factor family with sensor-grounded assessments, a re-derived body, and workflow-faithful report prose.
status: Done
tags: [examples, report-gallery, authoring, docs]
timestamp: 2026-07-02T00:00:00Z
---

# Deepen the report gallery exemplar model

Change Case [0190](0190-report-gallery-exemplar.md) made the
[report gallery](../../examples/report-gallery/software-service/README.md) an
authoring exemplar in structure: a composite root, a model-wide
agent-harnessability factor, a veto requirement, a measured override, a
synthetic changelog. But its model stops short of the authoring guides' own
coverage bar — most factors carry exactly one requirement, a payments ledger
carries no security or auditability lens, no area covers the implementation
itself, and nearly every assessment reads as a bespoke inferential
investigation.

This case deepens the exemplar to realistic scale and grounding: factor
families per constituent that meet the factors guide's stable-stakes bar, a
small named sensor catalog that assessments visibly reuse, a deliberate spread
of computational, sensor-plus-guide, and inferential assessment kinds
(the harness-engineering posture of feedforward guides and feedback sensors),
a Markdown body re-derived so the fictional scenario earns every factor, and
generated report prose consistent with what the `/quality` setup and evaluate
workflows would produce.

- [Functional spec](0191-report-gallery-model-realism/spec.md) — the delta the
  deepened gallery must satisfy.
- Design doc — expected to be omitted: the generator keeps the shape 0190 gave
  it (embedded content files, generalized payload tables); if the payload
  tables need restructuring beyond scale, a design doc is added at the Design
  phase.

## Motivation

The gallery is the first full QUALITY.md a prospective adopter inspects, and
it now claims to demonstrate authoring practice. A model that under-delivers
on the factors guide's own coverage aim (roughly ten factors per
primary-subject constituent), omits the lenses any reader of the payments
domain expects (security, auditability), has no home for the maintainability
and architecture-fitness concerns of harness engineering, and names two
factors the factors guide explicitly calls out as practice names
(`context-grounding`, `lifecycle-maintenance`) teaches the wrong lessons at
the exact moment it claims to teach the right ones. Deepening the model also
makes the assessment economics visible — which requirements a recorded sensor
settles, which need a sensor read against a guide, and which remain judgment —
so the gallery demonstrates how assessments lean on guides and computational
sensors rather than merely asserting that they do.

## Scope

**Covered:** the LedgerLite `QUALITY.md` frontmatter and body, the generated
`0001-full-eval` run content (payloads, findings, recommendations, rankings,
coverage), the synthetic quality changelog, the gallery README, and the
`scripts/report-gallery` content and payload tables that carry all of it.

**Deferred:** a second non-software gallery example; a checked-in fictional
source tree behind `synthetic-source:*`; a second evaluation run closing the
loop across runs; a model-wide root `currentness` factor demonstrating the
composite-root recurring-factor move.

**Non-goals:** changes to the report format, evaluation record schema, CLI,
or `/quality` skill behavior; edits to the authoring guides themselves — if
the exemplar and a guide conflict, the exemplar bends here and any guide
amendment is its own scoped change.

## Affected artifacts

**Code**

- `scripts/report-gallery/content/` — exemplar `QUALITY.md`, README, and
  changelog content files.
- `scripts/report-gallery/` — payload tables (requirement cases, factor
  tables, findings, rankings, recommendations, coverage) scaled to the
  deepened model.

**Generated example (regenerated, not hand-edited)**

- `examples/report-gallery/software-service/QUALITY.md`
- `examples/report-gallery/software-service/README.md`
- `examples/report-gallery/software-service/.quality/changelog/*`
- `examples/report-gallery/software-service/.quality/evaluations/0001-full-eval/**`

**Durable docs**

- `CHANGELOG.md` — note the gallery deepening under Unreleased.
- None otherwise: `mintlify/index.mdx` and `docs/guides/reporting-design.md`
  keep the stable `0001-full-eval` path and need no edit.

**Durable specs:** none — the gallery's content has no durable spec owner
(per 0190, it remains governed by its generator plus the
`report-gallery-check` gate), and no report-format or skill behavior changes.

## Status

`Done`. Implemented and archived. The generator-owned LedgerLite gallery now
uses a 39-requirement model with a `codebase` area, expanded per-area factor
families, reused sensor-grounded assessments, refreshed body context, new
synthetic changelog entries, and workflow-faithful generated report prose.
