---
type: Functional Specification
title: View command — functional spec
description: Behavior for the qualitymd view presentation command — default document render and the outline lens on text, json, and mermaid surfaces.
tags: [cli, view, presentation]
timestamp: 2026-06-24T00:00:00Z
---

# View command — functional spec

Companion to the [View command](../0078-view-command.md) change case. This spec
states *what* the first slice of `qualitymd view` must do; the design doc covers
*how* the layered surface is built and why. It inherits the cross-cutting CLI
contract from the [CLI spec](../../specs/cli.md) and defers the QUALITY.md format
itself to [`SPECIFICATION.md`](../../SPECIFICATION.md).

The key words "MUST", "MUST NOT", "SHOULD", and "MAY" are to be interpreted as
described in BCP 14 when, and only when, they appear in all capitals.

## Background / Motivation

A `QUALITY.md` workspace has no presentation surface today: a reader either parses
raw YAML/Markdown or scrapes `status`, which is deliberately a compact routing
snapshot (it MUST stay route-oriented and MUST NOT print every source-coverage
row). `view` is the missing presentation home — a read-only command that renders
the model (and, in later slices, its evaluation state) for humans and agents.

`view` is organized along two orthogonal axes so it can grow without re-shaping:
a **lens** selects *what* is rendered (a subcommand: the document, `outline`, and
later `health`/`coverage`/`trends`/`recommendations`), and a **surface** selects
*how* (`--format`: `text`, `json`, `mermaid`, and later `dot`/`html`/`--serve`).
Every surface renders the same deterministic workspace data, so the data is the
durable asset and each renderer is a consumer of it. This slice builds the default
document render and the `outline` lens on `text`, `json`, and `mermaid`.

## Scope

Covered: `qualitymd view [path]`; the default (no-lens) document render; the
`outline` lens; the `text`, `json`, and `mermaid` surfaces; `--json` as an alias;
read-only and deterministic behavior; and missing/invalid-model handling.

Deferred (designed, not built here): the `health`, `coverage`, `trends`, and
`recommendations` lenses; the `dot`, `html`, and `--serve` surfaces; reading the
model from stdin (`view -`); and a standalone durable spec for the workspace-graph
data contract.

Non-goals: `view` never judges quality, recomputes ratings, scrapes report bodies,
writes files, or changes the QUALITY.md format. Operational evaluation work stays
with `evaluation`.

## Requirements

### Invocation and safety

`view` **MUST** inspect `QUALITY.md` in the current working directory by default,
and `view <path>` **MUST** inspect the file named by `<path>` instead.

`view` **MUST** be read-only: it **MUST NOT** write, create, repair, or delete
files.

`view` output **MUST** be deterministic: an unchanged model file produces
byte-equivalent output for a given lens and surface.

> Rationale: determinism is what lets a `mermaid` or `json` surface be embedded in
> a report or diffed in review, and is the property every later renderer relies on.

`view` **MUST** be read-only presentation: it **MUST NOT** judge quality or
recompute ratings. It renders the model (and, in later slices, evaluation records)
as data.

### Lens selection

`view` with no lens subcommand **MUST** render the **document** lens. `view
outline` **MUST** render the **outline** lens. An unrecognized lens **MUST** fail
with a usage error.

### Document lens (default)

The document lens **MUST** pretty-render the model's Markdown body using the shared
glamour-based renderer that `qualitymd spec` uses, and **MUST** present the
frontmatter model as a readable header (at least the model title and rating-scale
levels). It **MUST NOT** emit the raw YAML frontmatter as undecorated text.

### Outline lens

The outline lens **MUST** render the model hierarchy from the root area through its
areas and factors, showing per-area and per-factor requirement counts, nesting
child areas, and including a rating-scale legend and document totals (areas,
factors, requirements). It renders at a counts-plus-top-levels depth and **MUST**
derive its counts from the same deterministic model-shape computation `status`
performs, so the two never disagree.

### Surfaces

`view` **MUST** support `--format` with the values `text` (default), `json`, and
`mermaid`. `--json` **MUST** be accepted as an alias for `--format json`, matching
the established CLI idiom (`status --json`, `version --json`).

- `text` is the human terminal surface.
- `json` emits the deterministic data the lens renders, with no terminal styling or
  control sequences, suitable for agents and for future renderers.
- `mermaid` emits a Mermaid diagram source for the lens (for the `outline` lens, a
  diagram of the area/factor hierarchy) that renders inline in Markdown.

An unrecognized `--format` value **MUST** fail with a usage error.

> Rationale: `mermaid` is the cheap, Markdown-native surface the agent can paste
> into a report or PR; `json` is the stable contract every later surface (html,
> serve) consumes. Both are chosen ahead of a heavier interactive tier.

### Missing and invalid models

When the model file is absent or unreadable, `view` **MUST** fail with a clear
error and a non-zero exit per the shared CLI exit-status contract; there is nothing
to render.

When the model file exists but is not lint-valid, the `outline` lens (and the
model-shape portion of any surface) **MUST NOT** derive partial model-shape from
the invalid model; `view` **MUST** report the invalidity instead. The document lens
**MAY** still render the Markdown body when it is separable from the invalid
frontmatter, but **MUST** surface that the model is invalid.

## Durable spec changes

Durable **specs** this case rewrites — the [`specs/`](../../specs/index.md) bundle
and [`SPECIFICATION.md`](../../SPECIFICATION.md). See
[Writing functional specs](../../docs/guides/write-functional-specs.md#durable-spec-changes).

### To add

- `specs/cli/view.md` - new durable command spec for `view`, carrying the two-axis
  lens/surface model, the document and outline lenses, the text/json/mermaid
  surfaces, and the read-only/determinism contract (per all requirements above).

### To modify

- `specs/cli/index.md` - register the new `view` spec.
- `specs/cli.md` - list `view` in the command set if it enumerates commands (per
  the invocation requirement above).
- `specs/cli/status.md` - cross-reference the status (routing) vs. view
  (presentation) boundary so neither grows the other's job (per the document/outline
  lens requirements above).
- `specs/skills/quality-skill/workflows/setup.md` - specify that setup closeout
  renders the model outline via `view` (per the outline lens requirement above).

### To rename

None.

### To delete

None.
