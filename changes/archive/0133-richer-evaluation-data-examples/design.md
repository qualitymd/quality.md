---
type: Design Doc
title: Richer evaluation data examples — design
description: How generated Evaluation data examples become fuller representative artifacts.
tags: [cli, evaluation, examples]
timestamp: 2026-06-26T00:00:00Z
---

# Richer evaluation data examples — design

## Context

Answers the [functional spec](spec.md) for change case
[0133](../0133-richer-evaluation-data-examples.md). The current examples are
built by small constructor functions in `internal/evaluation/data.go`, and their
structural contracts live in `internal/evaluation/data_contract.go`.

The implementation already has one contract registry that drives validation,
schema generation, and example kind resolution. The gap is not a missing schema;
it is that some examples contain empty arrays where the CLI spec asks for
representative nested object entries, and there is no all-kind test protecting
that example surface.

## Approach

### Enrich examples at the constructor source

Keep one constructor per data kind and add small shared helpers for repeated
example concepts:

- canonical model refs: root Area, child Area, root Factor, child Factor, root
  Requirement, and Rating Level;
- representative limits and stop conditions;
- routine refs and report refs that carry the correct typed subject IDs.

Then update the frame examples so currently empty `expectedEvaluationLimits` and
`stopConditions` arrays carry one representative object. Keep scalar strings
short and domain-neutral. Enrich `EvaluationOutputResult` so `areaOutputs[]`,
`reportRefs[]`, and `reportOutputs[]` demonstrate Area, Factor, and Requirement
report references rather than only the root Area report.

### Keep examples representative, not exhaustive

Do not create a combinatorial example corpus. One payload per kind remains the
contract because `data example <kind>` is a quick authoring aid. Exhaustive
constraints stay in `data schema <kind>` and validation tests.

### Add all-kind tests

Add an Evaluation-package test that loops through `dataContractOrder`, calls
`DataExample(kind)`, decodes the JSON, checks the emitted `kind`, and validates
the payload with `validateDataPayload`. This catches drift between the kind
registry, example constructors, and structural contract.

Add a focused assertion over the generated examples for the reference strings
the examples must teach:

- `area:root` and a non-root Area reference;
- `factor:root::verification` and a nested Factor reference;
- `requirement:root::has-tests`;
- `rating:target`;
- report refs for `area`, `factor`, and `requirement`.

Model-bound validation is not part of `DataExample` because examples are emitted
without a run's `model-snapshot.md`; the structural validator is the right
in-package check for the standalone examples.

## Spec response

- **Generated example content** — satisfied by enriching the constructors and
  shared helper output.
- **Typed ID demonstration** — satisfied by using canonical model refs in
  subjects, input refs, relationship refs, and report refs.
- **Verification** — satisfied by all-kind structural validation plus focused
  reference coverage tests.
- **Durable spec clarity** — satisfied by updating the CLI data spec to
  distinguish representative examples from exhaustive constraints.

## Alternatives

- **Generate a complete run folder as the example corpus.** Rejected for this
  case. It belongs to the future checked-in Evaluation example corpus and would
  make `data example <kind>` heavier than a one-artifact command.
- **Add one example per enum/status variant.** Rejected because the schemas are
  already the machine-readable source for enum sets, and examples would become
  noisy quickly.
- **Model-bind examples against a hidden example model.** Rejected because the
  command emits standalone artifacts without a run snapshot. Tests can use the
  structural contract directly without inventing hidden command state.

## Trade-offs & risks

- Fuller examples are longer. That is acceptable because `data example` is an
  authoring reference, and complete representative objects reduce guesswork.
- A representative example can still be mistaken for the only valid shape. The
  durable CLI spec now calls out that schemas carry the exhaustive constraint
  surface.

## Open questions

None for this design.
