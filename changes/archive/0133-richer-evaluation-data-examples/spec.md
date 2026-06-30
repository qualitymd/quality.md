---
type: Functional Specification
title: Richer evaluation data examples — functional spec
description: What the change must do to make generated Evaluation data examples fuller reference artifacts.
tags: [cli, evaluation, examples]
timestamp: 2026-06-26T00:00:00Z
---

# Richer evaluation data examples — functional spec

Companion to the
[Richer evaluation data examples](../0133-richer-evaluation-data-examples.md)
change case. This spec states _what_ the change must do.

The key words **MUST**, **MUST NOT**, **SHOULD**, and **MAY** are to be
interpreted as described in IETF RFC 2119.

## Background / Motivation

`qualitymd evaluation data example <kind>` is the concrete example surface for
Evaluation payload authoring. Agents should be able to inspect a generated
example and see the payload's normal object shape, representative repeated
nested entries, and canonical model-reference strings. The examples should teach
shape, but schemas and validation tests should remain the exhaustive constraint
source for every enum value, status value, and invalid case.

## Scope

Covered: the generated Evaluation data examples for the current data kinds, the
tests that validate those examples, and the durable CLI spec for the example
surface.

Deferred: a checked-in complete Evaluation run corpus, every enum/status
permutation, and any `qualitymd init` scaffold-template change.

## Requirements

### Generated example content

- `qualitymd evaluation data example <kind>` **MUST** continue to emit one
  complete valid JSON artifact for every supported Evaluation data kind.

  > Durable spec: modify `specs/cli/evaluation-data.md` — preserve the one-kind,
  > one-artifact command contract.

- Generated Evaluation data examples **MUST** include representative entries for
  repeated nested object fields that teach payload shape, including findings,
  rating drivers, unknown or missing-evidence entries, analysis input refs,
  stop conditions, and evaluation limits where those fields appear in the
  requested kind.

  > Rationale: empty arrays validate but do not teach agents how to author the
  > objects the field accepts.
  >
  > Durable spec: modify `specs/cli/evaluation-data.md` — clarify the
  > representative-entry expectation without turning examples into an exhaustive
  > corpus.

- Generated Evaluation data examples **MUST** demonstrate canonical Area,
  Factor, Requirement, and Rating Level reference strings in the relevant
  subject, input, rating, finding relationship, and report-reference fields.

  > Rationale: the examples are the easiest place for agents to see that IDs are
  > typed reference strings such as `area:root`, `factor:root::verification`,
  > `requirement:root::has-tests`, and `rating:target`.
  >
  > Durable spec: modify `specs/cli/evaluation-data.md` — name typed model
  > references as part of what representative examples demonstrate.

### Verification

- Tests **MUST** exercise every supported generated example kind and verify that
  each emitted artifact satisfies the typed structural contract for its `kind`.

  > Durable spec: none.

- Tests **MUST** verify that the generated example set includes Area, Factor,
  Requirement, and Rating Level references in representative places, including
  routine refs and report refs.

  > Durable spec: none.

- The durable CLI spec **MUST** state that generated examples are representative
  complete artifacts and **MUST NOT** imply they cover every valid enum, status,
  or error case.

  > Durable spec: modify `specs/cli/evaluation-data.md` — keep `data schema` as
  > the exhaustive machine-readable constraint surface and examples as concrete
  > instances.

## Durable spec changes

Rollup of the per-requirement `Durable spec:` annotations above (those are
authoritative) — the [`specs/`](../../../specs/index.md) bundle and
[`SPECIFICATION.md`](../../../SPECIFICATION.md).

### To add

None.

### To modify

- `specs/cli/evaluation-data.md` — clarify that `data example <kind>` emits one
  representative complete artifact per kind, includes representative nested
  objects and typed model references, and is not an exhaustive enum/status/error
  corpus.

### To rename

None.

### To delete

None.
