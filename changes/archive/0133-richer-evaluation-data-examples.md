---
type: Change Case
title: Richer evaluation data examples
description: Make `qualitymd evaluation data example` payloads fuller and more useful as reference artifacts.
status: Done
tags: [cli, evaluation, examples]
timestamp: 2026-06-26T00:00:00Z
---

# Richer evaluation data examples

A **Change Case** capturing the *why* and *status*; the detail lives in its
children:

- [Functional spec](0133-richer-evaluation-data-examples/spec.md) — what the case must do.
- [Design doc](0133-richer-evaluation-data-examples/design.md) — how it's built, and why.

## Motivation

`qualitymd evaluation data example <kind>` is currently the active CLI source
for concrete Evaluation payload examples. It emits one artifact for every
supported kind, but several frame examples leave representative repeated object
fields such as limits and stop conditions empty, and the tests only check a small
subset of generated examples. That makes the examples weaker than the durable
contract says they should be and leaves agents to infer some Area, Factor, and
Requirement ID/reference shapes from schemas alone.

This case makes each generated example a fuller representative artifact: every
kind still emits one valid JSON payload, repeated nested object fields carry
representative entries where the contract calls for them, and the example set
demonstrates Area, Factor, and Requirement identifiers in subjects, input refs,
and report refs.

## Scope

Covered:

- Enrich `qualitymd evaluation data example <kind>` outputs for the current
  Evaluation data kinds.
- Ensure generated examples demonstrate canonical Area, Factor, Requirement, and
  Rating Level reference strings where those fields appear.
- Add tests that exercise every generated example kind, not only one or two
  spot checks.
- Update the durable CLI spec so examples are explicitly representative, not an
  exhaustive corpus of every enum value or status state.

Deferred:

- A checked-in full Evaluation example corpus. The skill examples index already
  reserves that for future work.
- Exhaustive positive examples for every enum/status combination. Schemas and
  validation tests remain the exhaustive constraint surface.
- Alternate `qualitymd init` scaffolds or domain profiles.

## Affected artifacts

Derived by sweeping for `data example`, `DataExample`, Evaluation data example
constructors, and generated example references across code, specs, and skill
content.

**Code**

- [x] `internal/evaluation/data.go` — enrich generated Evaluation data examples
      and reusable example reference helpers.
- [x] `internal/evaluation/evaluation_test.go` — add all-kind example validation
      and ID/reference coverage tests.
- [x] `internal/cli/evaluation_test.go` — existing command coverage remains
      aligned; no output expectation changes were needed.

**Durable specs** (substance in the [functional spec](0133-richer-evaluation-data-examples/spec.md))

- [x] `specs/cli/evaluation-data.md` — clarify that `data example` emits a
      representative complete artifact per kind and is not the exhaustive
      constraint corpus.

**Durable docs / bundled skill runtime**

- No planned runtime skill changes. The bundled skill already routes agents to
  `data schema` for constraints and `data example` for concrete instances.

No planned impact: `SPECIFICATION.md`, `README.md`, install/scaffold files,
`CHANGELOG.md`, or generated JSON Schema artifacts.

## Status

`Done`. See the [status lifecycle](../index.md#status-lifecycle). Implemented and
archived. Evaluation data examples now emit fuller representative payloads with
canonical Area, Factor, Requirement, and Rating Level references; all generated
example kinds are structurally validated in tests.
