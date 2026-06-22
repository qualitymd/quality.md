---
type: Design Doc
title: Companion JSON Schema
description: How the companion JSON Schema is generated, bundled, and emitted, including terminal highlighting for qualitymd schema.
tags: [schema, cli, format, tooling]
timestamp: 2026-06-22T00:00:00Z
---

# Companion JSON Schema

How the [functional spec](spec.md) is delivered: how the schema artifact is
generated from `internal/schema`, bundled into the binary, and emitted by
`qualitymd schema`. Behavior and requirements live in the spec; this doc records
the *how* and the *why*.

## Context

The spec requires a structural JSON Schema for QUALITY.md frontmatter that is
**derived from the same source the linter uses**
([`internal/schema`](../../../internal/schema/schema.go)) so it cannot drift, and
a verbatim-artifact `qualitymd schema` command that emits it — plain bytes when
redirected, optionally highlighted on a terminal.

## Approach

### Terminal highlighting — chroma (decided)

`qualitymd schema` follows the `spec` command's two-branch shape: a `writeJSON`
helper that, like `writeMarkdown` in
[`internal/cli/spec.go`](../../../internal/cli/spec.go), checks `colorEnabled(w)`.

- **Plain branch** (non-TTY, `NO_COLOR`, redirect/pipe): write the verbatim
  indented JSON bytes — nothing else. This is what satisfies the spec's
  byte-stability requirement; `qualitymd schema > quality.schema.json` reproduces
  the artifact exactly.
- **TTY branch:** syntax-highlight with
  [`github.com/alecthomas/chroma/v2`](https://github.com/alecthomas/chroma), then
  page through the existing `page()` helper.

```go
// TTY branch
quick.Highlight(w, string(jsonBytes), "json", terminalFormatter(), schemaStyle())
```

chroma is already present in the module tree as an indirect dependency (glamour
pulls it in to highlight fenced code blocks); this promotes it to a **direct**
dependency — a one-line `go.mod` change, no new download. It ships a `json`
lexer, terminal formatters (`terminal16m` truecolor, `terminal256`), and a style
registry.

Two reasons this is byte-safe: chroma only injects ANSI escapes *around* tokens,
so the highlighted form differs from the plain form only by escape codes —
indentation and content are untouched — and highlighting runs **only** on the
TTY branch, never on the bytes a redirect captures.

Formatter selection mirrors the stack's existing color handling: pick
`terminal16m` when the terminal advertises truecolor, else `terminal256`.

### Schema generation — derived from `internal/schema`

`internal/schema` already models the structure as `Node` / `Property` /
`RequiredAny` values ([`Model`, `Area`, `Factor`, `Requirement`,
`RatingLevel`](../../../internal/schema/schema.go)). A generator walks `Nodes`
and emits JSON Schema:

- `ScalarShape` → `type: string`; `MapShape` → `object` with
  `additionalProperties: <$ref to element node>` and `patternProperties` for the
  map-keyed entries (`factors` / `requirements` / `areas`); `SequenceShape` →
  `array` with `items: <$ref>` and, for the rating scale, `minItems: 2`.
- `RequiredPresence` → the property is listed in the node's `required`.
- `RequiredAny` → `anyOf` requiring at least one of the named properties.
- Recursion (`Area`, `Factor`) → `$ref` into `$defs` so the schema is finite.

Open-ended frontmatter (the format's
[Extensions](../../../SPECIFICATION.md#extensions) rule) means nodes do **not**
set `additionalProperties: false`.

The generation logic is an exported `GenerateJSON()` in `internal/schema` — one
tested implementation, called by both a thin `//go:generate go run …` entrypoint
that writes the committed root file and the consistency test that guards it (see
[Artifact location and bundling](#artifact-location-and-bundling--repo-root-decided)).

**Dialect and `$id` (decided):** the schema declares `$schema` as **JSON Schema
draft 2020-12** — the current standard, fully supporting the `$defs`/`$ref`,
`patternProperties`, `anyOf`, and `minItems` constructs above, and supported by
the YAML language server that backs the editor scenario. It declares a stable,
**unversioned** `$id` of `https://quality.md/quality.schema.json`. A `$id` is an
identifier, not a fetch URL, so reserving the project-domain identity is
compatible with the spec's no-hosted-registry non-goal — identity now, hosting
later. (If `quality.md` is not the canonical project domain, fall back to the
GitHub raw-root URL.) Schema versioning stays deferred per the spec; a
`$comment` notes the artifact tracks the structural layer of `SPECIFICATION.md`
rather than coupling to a spec-version number.

### Artifact location and bundling — repo root (decided)

`quality.schema.json` lives at the **repo root**, a sibling of
[`SPECIFICATION.md`](../../../SPECIFICATION.md), and is embedded by a new root
`schema.go` in package `qualitymd` that mirrors
[`specification.go`](../../../specification.go) — a `func Schema() []byte`
accessor next to `Specification()`. `qualitymd schema` emits those bytes.

This follows the repo's canonical-public-artifact pattern: the schema is a
first-class public companion to the prose spec, so the two sit together, it is
discoverable, `qualitymd schema > quality.schema.json` round-trips to the same
path, and a root URL is the natural target for a hosted `$id` if one is added
later. The generator itself lives in/near `internal/schema`; only the committed
output and its embed sit at the root.

Rejected: co-locating under `internal/schema/` (the
[`internal/scaffold/skeleton.md`](../../../internal/scaffold/scaffold.go) pattern)
reads as a private implementation detail, undercutting the point of a schema
external editors and tools point at; and a dedicated top-level `schema/` dir is a
premature directory for a single file (YAGNI) that breaks the root-artifact
convention.

Generation is a **`go:generate` tool that writes the committed root file**
(decided over runtime generation): `qualitymd schema` emits the embedded bytes,
exactly as `spec` emits the embedded `SPECIFICATION.md`, so the embedded file
*is* the golden file. A test re-runs `GenerateJSON()` and fails if its output
differs from the committed/embedded file, so the bundled artifact can never
drift from `internal/schema`. This keeps schema changes visible as a reviewable
diff in the committed file — the win that runtime generation would forfeit for a
guarantee the test already provides.

## Alternatives

**Highlighting via a glamour code fence (rejected).** Wrapping the JSON in a
`` ```json `` fence and rendering through the glamour `TermRenderer` `spec` already
uses needs zero new deps and looks identical to `qualitymd spec`. Rejected
because glamour reflows content — `WithWordWrap(100)` wraps long lines and it
adds code-block padding/background — so the rendered view stops being a faithful
picture of the artifact a reader would redirect. For an artifact whose value is
byte-fidelity, that is the wrong trade.

**Hand-rolled lipgloss tokenizer (rejected).** A small JSON tokenizer colored
with the charmtone palette would match the brand exactly and add no highlighting
dep. Rejected as reinventing a lexer chroma already implements correctly, with
ongoing edge-case maintenance (escapes, unicode, nesting). If brand-exact colors
later matter, the cheaper path is a custom `chroma.Style` built from charmtone
tokens — chroma's correct lexing with our palette — kept as a future refinement,
not first-cut scope.

**Hand-authored schema file (rejected).** Maintaining `quality.schema.json` by
hand decouples it from the linter and guarantees eventual drift — the exact
failure the spec's no-drift requirement forbids.

## Trade-offs & risks

- chroma's built-in styles won't exactly match the lipgloss/charmtone brand. Low
  stakes for a developer artifact; the custom-style path above is the escape
  hatch if it matters.
- The generator must track `internal/schema` as that model gains shapes. The
  consistency test makes drift a loud build failure rather than a silent one.
- JSON Schema cannot express the semantic layer; the spec already requires the
  artifact to disclaim it, so the risk is mis-messaging, not mis-validation.

## Open questions

None outstanding. The design decisions — chroma highlighting, repo-root artifact
and embed, `go:generate` generation with a consistency test, draft 2020-12, and
an unversioned project-domain `$id` — are settled. The only external input
pending is confirmation that `quality.md` is the canonical project domain for the
`$id` (otherwise the GitHub raw-root URL is the fallback).
