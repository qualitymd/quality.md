---
type: Functional Specification
title: Companion JSON Schema
description: A structural JSON Schema for QUALITY.md frontmatter, derived from the linter's schema, and a qualitymd schema command that emits it.
tags: [schema, cli, format, tooling]
timestamp: 2026-06-22T00:00:00Z
---

# Companion JSON Schema

This spec governs a companion **JSON Schema** for QUALITY.md frontmatter and the
`qualitymd schema` command that emits it. The normative format remains
[`SPECIFICATION.md`](../../../SPECIFICATION.md); this spec adds a derived,
machine-readable structural artifact and the command that ships it.

The `schema` command inherits invocation-wide behavior, stdout/stderr
separation, determinism, plain-output rules, and exit-code categories from the
[CLI spec](../../../specs/cli.md). This spec states only the behavior particular
to the schema artifact and the `schema` command.

The key words "MUST", "MUST NOT", "SHOULD", "SHOULD NOT", and "MAY" are to be
interpreted as described in [RFC 2119](../../../docs/reference/rfc2119.md) and
[RFC 8174](../../../docs/reference/rfc8174.md) when, and only when, they appear
in all capitals.

## Background / Motivation

`SPECIFICATION.md` is the normative source of truth, and
[`internal/schema`](../../../internal/schema/schema.go) already expresses the
structural frontmatter schema in Go for the linter. External consumers — editor
validation/autocomplete for hand-authors, and third-party tools that exchange
QUALITY.md documents — have no portable schema to point at without running the
Go linter. A JSON Schema fills that gap, but only for the *structural* layer:
shapes, presence, recursion, `minItems`, map-keyed entries, and the
"at least one of factors/requirements/areas" rule. The *semantic* layer
(factor-reference resolution, rating-override keys matching declared levels, the
placement-dependent factor-connection rule, level ordering and uniqueness) stays
with the linter and cannot be expressed in JSON Schema. The artifact must
therefore be non-normative, subordinate to `SPECIFICATION.md`, and derived from
the same source the linter enforces so it cannot drift.

## Scenario

A hand-author opens `QUALITY.md` in an editor with a YAML language server and a
`# yaml-language-server: $schema=…` reference (or a tool configured to validate
QUALITY.md). They get inline structural diagnostics and completion for
frontmatter keys without installing or running `qualitymd`. To obtain or vendor
the schema, they run `qualitymd schema > quality.schema.json`.

## Scope

Covered: the companion JSON Schema artifact (its content and guarantees) and the
`qualitymd schema` command surface and output.

Deferred / non-goals:

- No JSON form of the specification *prose* — that is a separate deferral in
  [`spec`](../../../specs/cli/spec.md) and is not this artifact.
- The schema does not encode the linter's semantic rules; structural-only by
  design. A document the schema accepts can still fail `lint`.
- No `schema`-specific flags beyond the cross-cutting CLI flags.
- No published hosted-URL `$id` registry or versioned schema-hosting scheme
  beyond what the artifact itself declares.

## Requirements

### The schema artifact

- The artifact **MUST** be a valid JSON Schema describing the structural shape of
  QUALITY.md frontmatter: the Model root, and recursively Areas, Factors,
  Requirements, and Rating Levels, matching the structure in
  [`SPECIFICATION.md`](../../../SPECIFICATION.md#frontmatter-schema).
- The artifact **MUST** encode the structural constraints the linter enforces
  from `internal/schema`: per-property shape (scalar / map / sequence),
  required vs. optional presence, the rating-scale `minItems` of 2, map-keyed
  entries for `factors` / `requirements` / `areas`, and the "at least one of
  factors, requirements, or areas" rule on the Model and Areas.

  >> Rationale: the value of a companion schema is that it agrees with the tool;
  >> a schema that accepts what the linter rejects (or vice versa) is worse than
  >> none. — 0049

- The artifact **MUST** be derived from the same structural schema the linter
  uses, such that it cannot encode a structural rule the linter does not, or omit
  one it does. The repository **MUST** carry a check (generation step or test)
  that fails when the shipped artifact and the structural schema disagree.

  >> Rationale: two hand-maintained copies drift; binding the artifact to
  >> `internal/schema` keeps the structural contract single-sourced. — 0049

- The artifact **MUST** identify itself as structural-only and non-normative,
  and point to `SPECIFICATION.md` as the normative source — e.g. via the
  schema's `$comment`/`description` and `title` — so a consumer does not mistake
  passing validation for full conformance.
- The artifact **MUST NOT** claim to validate semantic rules it cannot express
  (factor-reference resolution, rating-override keys, the factor-connection
  rule, level ordering/uniqueness). Where a structural approximation would
  falsely imply such a check, the artifact omits it rather than encoding a
  misleading one.
- The artifact **SHOULD** declare a stable JSON Schema dialect (`$schema`) and an
  `$id`, and **SHOULD** allow unknown extension properties consistent with the
  format's [Extensions](../../../SPECIFICATION.md#extensions) rules rather than
  forbidding all additional keys.

  >> Rationale: the format permits extension frontmatter; a closed schema would
  >> reject conforming documents that use it. — 0049

### The `qualitymd schema` command

- `schema` **MUST** emit the bundled companion JSON Schema to stdout, sourced
  from the schema the binary was built with, so the command is self-contained and
  needs no file on disk or network access.
- When output must be plain — stdout is not a terminal, or `NO_COLOR` is set —
  `schema` **MUST** write the schema as verbatim JSON and nothing else.
  `qualitymd schema > quality.schema.json` reproduces the artifact
  byte-for-byte.
- When stdout is a terminal, `schema` **MAY** syntax-highlight the JSON and page
  it for readability per the [paging convention](../../../specs/cli.md#conventions).
  This rendering is a human convenience; it **MUST NOT** change the bytes written
  when output must be plain.
- `schema` **MUST NOT** require any argument, and **MUST** treat an unexpected
  argument or flag as a usage error.
- `schema` **MUST NOT** offer `--json`. It is a verbatim-artifact carve-out per
  the [`--json` convention](../../../specs/cli.md#conventions); its output is
  already the JSON payload.
- On success `schema` **MUST** exit `0`. If it cannot emit the schema, it
  **MUST** exit `70`.

## Durable spec changes

### To add

- `specs/quality-schema-json.md` — new 1:1 artifact-spec (artifact-spec filename
  convention for `quality.schema.json`) carrying the durable contract for the
  companion schema: its structural-only/non-normative/subordinate status and the
  no-drift property (per the schema-artifact requirements above).
- `specs/cli/schema.md` — new durable CLI spec for `qualitymd schema` (per the
  command requirements above).

### To modify

- `specs/cli/index.md` — register the `qualitymd schema` command (per the command
  requirements above).
- `specs/index.md` and `specs/log.md` — register and log the new
  `quality-schema-json.md` durable spec.
- `SPECIFICATION.md` — add a non-normative note pointing to the companion JSON
  Schema and stating it is structural-only and subordinate (per the
  identification requirement above).
- `specs/cli/spec.md` — clarify that its deferred "JSON form of the
  specification" is the prose-as-JSON idea, distinct from the structural `schema`
  command (per the deferral in [Scope](#scope) above).

### To rename

None.

### To delete

None.
