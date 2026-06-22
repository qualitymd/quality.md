---
type: Functional Specification
title: quality.schema.json
description: Companion structural JSON Schema for QUALITY.md frontmatter, derived from the linter's schema.
tags: [schema, format, json, tooling]
timestamp: 2026-06-22T00:00:00Z
---

# quality.schema.json

`quality.schema.json` is the companion **JSON Schema** for QUALITY.md
frontmatter, shipped at the repository root beside
[`SPECIFICATION.md`](../SPECIFICATION.md) and emitted by
[`qualitymd schema`](cli/schema.md).

The normative format is [`SPECIFICATION.md`](../SPECIFICATION.md); this artifact
is a derived, machine-readable description of the frontmatter's *structural*
layer and is **subordinate** to it. This spec owns the artifact's content and
guarantees; the command that emits it is specified in
[`qualitymd schema`](cli/schema.md).

This document uses BCP 14 keywords only for testable conformance requirements.
The key words "MUST", "MUST NOT", "SHOULD", and "MAY" are to be interpreted as
described in [RFC 2119](../docs/reference/rfc2119.md) and
[RFC 8174](../docs/reference/rfc8174.md) when, and only when, they appear in all
capitals.

## Background / Motivation

[`internal/schema`](../internal/schema/schema.go) already expresses the
structural frontmatter schema in Go for the linter to consume. External consumers
— editor validation and autocomplete for hand-authors (via a
`# yaml-language-server: $schema=…` reference), and third-party tools that
exchange QUALITY.md documents — have no portable schema to point at without
running the Go linter. `quality.schema.json` fills that gap.

It can only fill the *structural* layer: shapes, required/optional presence,
recursion, `minItems`, map-keyed entries, strict model-name patterns where JSON
Schema can express them, and the "at least one of factors/requirements/areas"
rule. The *semantic* layer the linter owns
(factor-reference resolution, rating-override keys matching declared levels, the
placement-dependent factor-connection rule, level ordering and uniqueness) cannot
be expressed in JSON Schema. Drawing that line explicitly is what keeps the
artifact honest: a consumer must not read passing validation as full conformance,
and the artifact must not pretend to a check it cannot make. Because the value of
a companion schema is that it agrees with the tool, it is derived from the same
source the linter enforces rather than hand-maintained, so the two cannot
silently drift.

## Scope

Covered: the artifact's content (what structural shape and constraints it
encodes), its non-normative/subordinate status, its no-drift property, and the
schema-level declarations (dialect, `$id`, open-extension rule) it carries.

Deferred / non-goals:

- No encoding of the linter's semantic rules; structural-only by design. A
  document this schema accepts can still fail [`lint`](cli/lint.md).
- No JSON form of the specification *prose* — that is a separate deferral noted
  in [`spec`](cli/spec.md), distinct from this structural artifact.
- No qualitymd-only tooling conventions such as the root `config` pointer; those
  are lint/CLI behavior, not normative Model structure.
- No published hosted-URL `$id` registry or versioned schema-hosting scheme
  beyond the identifier the artifact declares.
- The command surface that emits the artifact is specified in
  [`qualitymd schema`](cli/schema.md), not here.

## Requirements

`quality.schema.json` **MUST** be a valid JSON Schema describing the structural
shape of QUALITY.md frontmatter: the Model root, and recursively Areas, Factors,
Requirements, and Rating Levels, matching the structure in
[`SPECIFICATION.md`](../SPECIFICATION.md#frontmatter-schema).

`quality.schema.json` **MUST** encode the structural constraints the linter
enforces from [`internal/schema`](../internal/schema/schema.go): per-property
shape (scalar / map / sequence), required vs. optional presence, the rating-scale
`minItems` of 2, map-keyed entries for `factors` / `requirements` / `areas`,
strict `propertyNames` patterns for `factors` and `areas`, the strict scalar
pattern for `ratingScale[].level`, and the "at least one of factors,
requirements, or areas" rule on the Model.

> Rationale: the value of a companion schema is that it agrees with the tool; a
> schema that accepts what the linter rejects (or vice versa) is worse than none.
> — 0049

`quality.schema.json` **MUST** be derived from the same structural schema the
linter uses, such that it cannot encode a structural rule the linter does not, or
omit one it does. The repository **MUST** carry a check that fails when the
shipped artifact and the structural schema disagree.

> Rationale: two hand-maintained copies drift; binding the artifact to
> `internal/schema` keeps the structural contract single-sourced. — 0049

`quality.schema.json` **MUST** identify itself as structural-only and
non-normative, and point to [`SPECIFICATION.md`](../SPECIFICATION.md) as the
normative source — for example through its `title`, `description`, or `$comment`
— so a consumer does not mistake passing validation for full conformance.

`quality.schema.json` **MUST NOT** claim to validate semantic rules it cannot
express (factor-reference resolution, rating-override keys, the factor-connection
rule, level ordering and uniqueness, and the area subtree
requirement-reachability check behind the warning-level `empty-area` rule).
Where a structural approximation would falsely imply such a check, the artifact
omits it rather than encoding a misleading one.

> Rationale: a near-miss structural check that looks semantic teaches consumers
> to trust the schema for guarantees only the linter makes. — 0049

`quality.schema.json` **MUST NOT** constrain Requirement statement keys or
Requirement `ratings` override keys with the strict model-name pattern.

`quality.schema.json` **MUST** declare its JSON Schema dialect (`$schema`) as
draft 2020-12 and **MUST** declare its `$id` as
`https://getquality.md/quality.schema.json`. The `$id` is an identifier, not a
fetch URL, so it remains stable independent of where the artifact is hosted.

`quality.schema.json` **MUST NOT** forbid unknown extension properties (it does
not set `additionalProperties: false`), consistent with the format's
[Extensions](../SPECIFICATION.md#extensions) rule.

> Rationale: the format permits extension frontmatter; a closed schema would
> reject conforming documents that use it. — 0049

`quality.schema.json` **MUST NOT** describe the qualitymd-only root `config`
convention as a normative Model property.

> Rationale: `config` selects qualitymd tooling configuration for a workspace;
> including it in the companion schema would blur tooling behavior with the
> QUALITY.md Model. — 0057
