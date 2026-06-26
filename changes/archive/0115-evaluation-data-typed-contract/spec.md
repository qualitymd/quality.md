---
type: Functional Specification
title: Type-safe, model-bound Evaluation v2 data — functional spec
description: Requirements for typed per-kind validation, model-binding, generated schema and examples, and the data schema/verify commands in qualitymd evaluation data.
tags: [cli, evaluation, data, validation]
timestamp: 2026-06-26T00:00:00Z
---

# Type-safe, model-bound Evaluation v2 data — functional spec

Companion to the
[Type-safe, model-bound Evaluation v2 data](../0115-evaluation-data-typed-contract.md)
change case. This spec states *what* the change must do; the
[design doc](design.md) covers *how*.

Normative durable specs this delta lands in: the
[`qualitymd evaluation data`](../../../specs/cli/evaluation-data.md) command spec
and the Evaluation v2
[payload kinds](../../../specs/evaluation-v2/records/payload-kinds.md) and
[JSON conventions](../../../specs/evaluation-v2/records/json-conventions.md). The
existing behavior they describe is the baseline these requirements tighten.

The key words **MUST**, **MUST NOT**, **SHOULD**, and **MAY** are to be
interpreted as described in BCP 14 when, and only when, they appear in all
capitals.

## Background / Motivation

Today `data set` decodes into `map[string]any` and validates a few named fields,
so unknown/misspelled fields and references to non-existent model nodes are
accepted and persisted; the failure only surfaces later (a blank report) or never.
There is no typed source of truth for the data kinds, so an agent cannot discover
the real payload sub-shapes and the `data example` output leaves the key arrays
empty. The fix binds one typed definition per kind to validation, schema
generation, and examples, and adds a semantic pass that resolves model references
against the run's snapshot — the same structure-vs-semantics split the frontmatter
schema already draws (structure in the schema, semantics in `qualitymd lint`).

## Scope

Covered: structural strictness, typed/enum and required-field validation,
model-binding, dry-run parity, generated schema and populated examples from one
source, the `data schema` and `data verify` commands, and pointing the `/quality`
skill's payload discovery at `data schema` (with `--dry-run` reframed as
validation, not shape-sniffing).

Deferred / non-goals: no wire-format change (`schemaVersion` stays `1`; a valid
payload's on-disk JSON is unchanged); no change to the accepted kind set, the
protocol order, or routine prompt contracts; no retroactive rewrite of persisted
data; no runtime JSON-Schema engine.

## Assumptions & dependencies

- Every evaluation run folder contains a parseable `model-snapshot.md` (written
  at run creation); the model-binding requirements rest on it. If it is absent or
  unparseable, [R8](#r8) governs.
- The Area / Factor / Requirement / Rating Level existence helpers in
  `internal/evaluation/model_reference.go` resolve references against a
  `*model.Spec`; the model-binding pass depends on them.

## Requirements

### Validation — `data set`

<a id="r1"></a>**R1 — Reject unknown fields.** When `qualitymd evaluation data set`
receives a payload containing an object field not defined for its `kind`, it
**MUST** reject the payload without writing, and the error **MUST** name the
offending field and the kind.

> Rationale: findings written with `title`/`summary` instead of
> `description`/`evidence` were accepted and rendered blank in the report. Silent
> acceptance of unknown fields is the failure this case exists to end. — 0115

<a id="r2"></a>**R2 — Enforce field types and enumerations.** For each accepted
kind, `data set` **MUST** validate that every present field has the type the kind
defines, and that an enumerated field (for example a status, confidence, finding
`type`, or `severity`) holds a value the kind permits, and **MUST** reject a
payload that violates either.

<a id="r3"></a>**R3 — Enforce required fields.** `data set` **MUST** reject a
payload that omits a field its kind requires, naming the missing field.

<a id="r4"></a>**R4 — Bind references to the model.** `data set` **MUST** resolve
every model-identity reference in a payload — Area IDs, Factor IDs, Requirement
IDs, and Rating Level IDs — against the run's model snapshot, and **MUST** reject
a payload that references an Area, Factor, Requirement, or Rating Level absent
from that snapshot, naming the unresolved reference.

> Rationale: `data set` previously checked identity fields only for path-safe
> characters, so invented IDs were persisted into the data graph. Existence
> checkers already existed but were wired only into tests. — 0115

<a id="r5"></a>**R5 — Dry-run parity.** Under `--dry-run`, `data set` **MUST**
perform the same structural, typed, required-field, and model-binding validation
as a real write and report the same errors, and **MUST NOT** persist anything.

> Rationale: the agent discovers payload shape by probing with `--dry-run`; a
> dry run that validates less than a real write teaches the wrong contract. — 0115

<a id="r6"></a>**R6 — Single source of truth.** The structural and typed
validation `data set` applies, the schema `data schema` emits, and the examples
`data example` emits **MUST** derive from one typed definition per kind, such that
none can diverge from the others. A drift check **MUST** fail if they do.

<a id="r7"></a>**R7 — Snapshot integrity.** If the run's model snapshot is missing
or cannot be parsed, then `data set` **MUST** fail with an error identifying the
snapshot, rather than skipping the model-binding validation in [R4](#r4).

> Rationale: fail closed — a run that cannot resolve its own model must not accept
> data as if every reference were valid. — 0115

### Discovery — `data schema` and `data example`

<a id="r8"></a>**R8 — Schema command.** `qualitymd evaluation data schema`
**MUST** emit a JSON Schema describing the accepted payload kinds, and
`qualitymd evaluation data schema <kind>` **MUST** emit the schema for the named
kind. The emitted schema **MUST** be the same definition [R6](#r6) validates
against.

<a id="r9"></a>**R9 — Populated examples.** `qualitymd evaluation data example
<kind>` **MUST** emit an example whose repeated and nested object fields (for
example `findings`, `ratingDrivers`, `unknowns`, and the analysis-scope arrays)
each contain at least one representative element, and the emitted example **MUST**
validate against that kind's schema and pass `data set --dry-run` against a run
whose model contains the referenced nodes.

> Rationale: the prior `data example` hard-coded empty arrays for exactly the
> sub-shapes an evaluator needed, forcing reverse-engineering from the binary. — 0115

<a id="r10"></a>**R10 — Kind-argument forms.** A command that takes a `<kind>`
argument **MUST** accept the kind's canonical name and **SHOULD** also accept its
kebab-case form (for example `area-analysis-result` for `AreaAnalysisResult`).

### Re-validation — `data verify`

<a id="r11"></a>**R11 — Verify a run.** `qualitymd evaluation data verify <run>`
**MUST** validate every persisted Evaluation v2 payload in the run against the
full [R1](#r1)–[R4](#r4) contract and report each payload that fails, identifying
the file and the reason. It **MUST NOT** modify any payload, and **MUST** exit
non-zero when any payload fails.

> Rationale: tightening acceptance leaves older runs potentially non-conformant;
> verify is the on-demand migration and self-check path, and the read-side
> counterpart to write-time validation. — 0115

### Skill discovery

<a id="r13"></a>**R13 — Skill discovers shape from the schema.** The `/quality`
skill's structured-payload discovery instruction **MUST** direct the agent to the
authoritative payload contract via `qualitymd evaluation data schema` (and the
populated `qualitymd evaluation data example <kind>`) as the source of payload
shape, and **MUST** frame `qualitymd evaluation data set --dry-run` as validation
of an authored payload rather than the means of discovering its shape.

> Rationale: the skill spec already requires `data kinds` → `data example` →
> `--dry-run`, and the field agent followed it — but with empty examples and no
> schema, `--dry-run` became a shape-sniffing loop. Adding the real discovery
> surface (R8, R9) without pointing the skill at it would leave the skill's
> discovery requirement stale. — 0115

### Compatibility

<a id="r12"></a>**R12 — No wire-format change.** This change **MUST NOT** alter
`schemaVersion` or the on-disk JSON shape of a payload that is already valid; it
**MUST** only tighten which payloads are accepted.

## Durable spec changes

### To add

None. (A 1:1 artifact spec for the generated JSON Schema file is *suggested* in
the parent's Affected artifacts, not required to land this case.)

### To modify

- [`specs/cli/evaluation-data.md`](../../../specs/cli/evaluation-data.md) —
  strengthen the `data set` validation requirements to cover unknown-field
  rejection, typed/enum and required-field checks, model-binding, and dry-run
  parity (per [R1](#r1)–[R5](#r5), [R7](#r7)); document the `data schema` (per
  [R8](#r8)) and `data verify` (per [R11](#r11)) commands; require `data example`
  to be populated and schema-valid (per [R9](#r9)); record the kebab-case
  `<kind>` allowance (per [R10](#r10)).
- [`specs/evaluation-v2/records/payload-kinds.md`](../../../specs/evaluation-v2/records/payload-kinds.md)
  — restate "the CLI **MUST** validate each accepted kind" as structural + typed +
  required-field + model-bound validation from one source of truth (per
  [R1](#r1)–[R4](#r4), [R6](#r6)).
- [`specs/evaluation-v2/records/json-conventions.md`](../../../specs/evaluation-v2/records/json-conventions.md)
  — note that persisted identity fields are resolved against the run's model
  snapshot and that unknown fields are rejected (per [R1](#r1), [R4](#r4)).
- [`specs/skills/quality-skill/quality-skill.md`](../../../specs/skills/quality-skill/quality-skill.md)
  — add `data schema` to the structured-payload discovery requirement and reframe
  `data set --dry-run` as payload validation rather than shape discovery (per
  [R13](#r13)).

### To rename

None.

### To delete

None.
