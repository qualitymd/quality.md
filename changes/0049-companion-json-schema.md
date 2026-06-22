---
type: Change Case
title: Companion JSON Schema
description: Publish a companion JSON Schema for QUALITY.md frontmatter, derived from the structural schema, and add a qualitymd schema command that emits it.
status: Design
tags: [schema, cli, format, tooling]
timestamp: 2026-06-22T00:00:00Z
---

# Companion JSON Schema

A **Change Case** capturing the *why* and *status* for publishing a companion
JSON Schema for QUALITY.md frontmatter and a `qualitymd schema` command that
emits it. The detail lives in its
[functional spec](0049-companion-json-schema/spec.md).

> **Design.** The [functional spec](0049-companion-json-schema/spec.md) and
> [design doc](0049-companion-json-schema/design.md) are settled: chroma terminal
> highlighting, repo-root `quality.schema.json` embedded via a new root
> `schema.go`, `go:generate` generation guarded by a consistency test, JSON
> Schema draft 2020-12, and an unversioned project-domain `$id`. The only pending
> external input is confirming `quality.md` as the `$id` domain. Ready to advance
> to **In-Progress**; no code is touched until then.

## Motivation

[`SPECIFICATION.md`](../../SPECIFICATION.md) is the normative source of truth for
QUALITY.md, and [`internal/schema`](../../internal/schema/schema.go) already
expresses the structural frontmatter schema in Go for the linter to consume.
What is missing is a portable, machine-readable schema artifact that external
tooling can use without running the Go linter — most concretely, editor
validation and autocomplete for hand-authors (via a
`# yaml-language-server: $schema=…` reference) and third-party tools that
exchange QUALITY.md documents.

A JSON Schema can capture the *structural* layer — shapes, required/optional
presence, recursion, `minItems`, map-keyed entries, and the "at least one of
factors/requirements/areas" rule — but it cannot express the *semantic* layer
the linter owns (factor-reference resolution, rating-override keys matching
declared levels, the placement-dependent factor-connection rule, level ordering
and uniqueness). The companion schema must therefore be explicitly
**non-normative and subordinate** to `SPECIFICATION.md`, structural-only, and
**derived from the same source the linter uses** so it cannot silently drift.

## Scope

Covered: a companion JSON Schema artifact for QUALITY.md frontmatter; a durable
artifact-spec governing it; a new `qualitymd schema` command that emits the
artifact verbatim (a verbatim-artifact carve-out like
[`spec`](../../specs/cli/spec.md)); and the durable CLI spec for that command.

Deferred / non-goals:

- No change to the normative format or to evaluation/report semantics.
- The JSON Schema does **not** replace or duplicate the linter's semantic
  checks; structural-only by design.
- No JSON form of the specification *prose* — that is a separate deferral noted
  in [`spec`](../../specs/cli/spec.md) and is not this case.
- No schema-driven rewrite of the linter; `internal/schema` stays the structural
  source of truth and the schema is derived from it.

## Affected artifacts

### Code

- `internal/schema/` - expose the structural schema as JSON Schema, derived from
  the existing `Node`/`Property`/`RequiredAny` model so the artifact cannot drift
  from what the linter enforces (generation mechanism is a design decision).
- `schema.go` (new root file in package `qualitymd`, mirroring
  [`specification.go`](../../specification.go)) - `//go:embed` the generated
  repo-root `quality.schema.json` and expose `func Schema() []byte` next to
  `Specification()`.
- `internal/cli/schema.go` (+ `schema_test.go`) - the `qualitymd schema` command
  emitting the embedded artifact verbatim to stdout, with a `writeJSON` helper
  mirroring `spec.go`'s `colorEnabled`/`page` two-branch shape.
- `internal/cli/root.go` - register the new command.
- `go.mod` - promote `github.com/alecthomas/chroma/v2` from an indirect to a
  direct dependency for terminal JSON highlighting (already in the module tree
  via glamour; no new download).
- Generation/consistency check (e.g. a `go:generate` step or a test) ensuring the
  embedded `quality.schema.json` matches the structural schema, so the two cannot
  fall out of sync.

### Durable specs

- `specs/quality-schema-json.md` - **new** 1:1 artifact-spec (artifact-spec
  filename convention for `quality.schema.json`) governing the companion schema:
  what it covers, its non-normative/subordinate status, and the no-drift
  property.
- `specs/cli/schema.md` - **new** durable CLI spec for `qualitymd schema`.
- `specs/cli/index.md` - register the `schema` command.
- `specs/index.md` and `specs/log.md` - register and log the new
  `quality-schema-json.md` durable spec.
- `SPECIFICATION.md` - add a short non-normative note pointing to the companion
  JSON Schema and stating it is structural-only and subordinate to this spec.
- `specs/cli/spec.md` - clarify that its deferred "JSON form of the
  specification" is the prose-as-JSON idea, distinct from the new structural
  `schema` command, so the two are not conflated.

### Durable docs, skill, and examples

- `README.md` - mention the `schema` command / companion schema only if the
  README enumerates commands or authoring aids.
- `docs/guides/cli-design.md` and `docs/guides/design-go-packages.md` - update
  only if the generation approach establishes a pattern worth documenting.
- `skills/quality/` - update only if hand-authoring guidance should point at the
  schema reference; otherwise no impact.
- Scaffold/install files - no expected impact.

## Children

- [Functional spec](0049-companion-json-schema/spec.md) - what the companion
  schema and `schema` command must do.
- [Design doc](0049-companion-json-schema/design.md) - how the schema is
  generated, bundled, and emitted, including chroma terminal highlighting.

## Status

`Design`. See the [status lifecycle](../index.md#status-lifecycle). The
functional spec and design doc are settled: chroma highlighting, repo-root
`quality.schema.json` embedded via a new root `schema.go`, `go:generate`
generation guarded by a consistency test, JSON Schema draft 2020-12, and an
unversioned project-domain `$id` (pending confirmation of the `quality.md`
domain). No open design questions remain; ready to advance to **In-Progress**.
No code touched until then.
