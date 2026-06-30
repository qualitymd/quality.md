---
type: Change Case
title: Companion JSON Schema
description: Publish a companion JSON Schema for QUALITY.md frontmatter, derived from the structural schema, and add a qualitymd schema command that emits it.
status: Done
tags: [schema, cli, format, tooling]
timestamp: 2026-06-22T00:00:00Z
---

# Companion JSON Schema

A **Change Case** capturing the _why_ and _status_ for publishing a companion
JSON Schema for QUALITY.md frontmatter and a `qualitymd schema` command that
emits it. The detail lives in its
[functional spec](0049-companion-json-schema/spec.md).

> **Done.** Landed and archived. The companion schema is generated from
> `internal/schema` by an exported `GenerateJSON()`, written to the repo-root
> `quality.schema.json` by a `go:generate` entrypoint, embedded via the new root
> `schema.go`, and emitted by `qualitymd schema`. A consistency test guards
> no-drift. The `$id` is `https://getquality.md/quality.schema.json` (draft
> 2020-12). `chroma` is a direct dependency for terminal highlighting. Build,
> full test suite, `go vet`, and `gofmt` are clean; `go generate` is idempotent
> and the redirect round-trip is byte-identical.

## Motivation

[`SPECIFICATION.md`](../../SPECIFICATION.md) is the normative source of truth for
QUALITY.md, and [`internal/schema`](../../internal/schema/schema.go) already
expresses the structural frontmatter schema in Go for the linter to consume.
What is missing is a portable, machine-readable schema artifact that external
tooling can use without running the Go linter — most concretely, editor
validation and autocomplete for hand-authors (via a
`# yaml-language-server: $schema=…` reference) and third-party tools that
exchange QUALITY.md documents.

A JSON Schema can capture the _structural_ layer — shapes, required/optional
presence, recursion, `minItems`, map-keyed entries, and the "at least one of
factors/requirements/areas" rule — but it cannot express the _semantic_ layer
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
- No JSON form of the specification _prose_ — that is a separate deferral noted
  in [`spec`](../../specs/cli/spec.md) and is not this case.
- No schema-driven rewrite of the linter; `internal/schema` stays the structural
  source of truth and the schema is derived from it.

## Affected artifacts

### Code

- `internal/schema/jsonschema.go` - exported `GenerateJSON()` rendering the
  structural schema as JSON Schema, derived from the existing
  `Node`/`Property`/`RequiredAny` model so the artifact cannot drift from what the
  linter enforces. The "at least one of factors/requirements/areas" `anyOf` lands
  on the Model only, faithfully matching `internal/schema` (Areas carry no such
  hard rule — their emptiness is the warning-level `empty-area` reachability
  check, which is semantic and not JSON-Schema-expressible).
- `internal/schema/gen/` - thin `package main` `go:generate` entrypoint that calls
  `GenerateJSON()` and writes the committed repo-root `quality.schema.json`.
- `schema.go` (new root file in package `qualitymd`, mirroring
  [`specification.go`](../../specification.go)) - `//go:embed` the generated
  repo-root `quality.schema.json` and expose `func Schema() []byte` next to
  `Specification()`.
- `internal/cli/schema.go` (+ `schema_test.go`) - the `qualitymd schema` command
  emitting the embedded artifact verbatim to stdout, with a `writeSchema` helper
  mirroring `spec.go`'s `colorEnabled`/`page` two-branch shape. (Named
  `writeSchema`, not `writeJSON`: the evaluation surface already has a
  value-encoding `writeJSON`.) Terminal highlighting uses `chroma`'s `quick`
  package; the formatter is chosen by `COLORTERM` (truecolor → `terminal16m`,
  else `terminal256`) and the style is isolated in one `schemaStyle` constant.
- `internal/cli/root.go` - register the new command in the Common Tasks group,
  immediately after `spec`.
- `go.mod` - promote `github.com/alecthomas/chroma/v2` from an indirect to a
  direct dependency for terminal JSON highlighting (already in the module tree
  via glamour; no new download).
- `schema_test.go` (new root test) - consistency check: the embedded
  `quality.schema.json` must equal a fresh `GenerateJSON()` render, so the
  bundled artifact cannot fall out of sync with the structural schema.
  `internal/schema/jsonschema_test.go` adds generator unit tests (determinism,
  dialect/`$id`, required/anyOf, `minItems`, recursive `$defs`, open extensions).

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

- `README.md` - **updated.** The CLI Quick Reference enumerates commands, so a
  `Show frontmatter schema | qualitymd schema` row was added beside `spec`.
- `docs/guides/cli-design.md` and `docs/guides/design-go-packages.md` - **no
  change.** The generate-embed-with-consistency-test approach follows the repo's
  existing canonical-public-artifact precedent (`specification.go` embedding
  `SPECIFICATION.md`); documenting a second instance of an established,
  already-undocumented pattern would pad the guides beyond what this case
  established.
- `skills/quality/` - **no change.** The skill is agent-first; the editor
  `# yaml-language-server: $schema=` workflow targets human hand-authors and is
  served by the README row and the `SPECIFICATION.md` note. No natural
  hand-authoring touchpoint exists in the skill.
- Scaffold/install files - no impact, as expected.

## Children

- [Functional spec](0049-companion-json-schema/spec.md) - what the companion
  schema and `schema` command must do.
- [Design doc](0049-companion-json-schema/design.md) - how the schema is
  generated, bundled, and emitted, including chroma terminal highlighting.

## Status

`Done`. See the [status lifecycle](../index.md#status-lifecycle). Landed after
review and archived. The companion schema is generated from `internal/schema` by
`GenerateJSON()`, written to the repo-root `quality.schema.json` by a
`go:generate` entrypoint, embedded via the new root `schema.go`, and emitted
verbatim by `qualitymd schema`; a root consistency test guards no-drift, and
generator unit tests cover the encoded structure. The `$id` is
`https://getquality.md/quality.schema.json` (JSON Schema draft 2020-12). Build,
`go test ./...`, `go vet`, and `gofmt` are clean; `go generate ./...` is
idempotent and `qualitymd schema > quality.schema.json` round-trips
byte-for-byte.
