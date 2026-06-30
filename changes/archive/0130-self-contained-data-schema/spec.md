---
type: Functional Specification
title: Self-contained per-kind data schema — functional spec
description: What the change must do to make `evaluation data schema <kind>` self-contained and align the skill.
tags: [cli, evaluation, schema, skill]
timestamp: 2026-06-26T00:00:00Z
---

# Self-contained per-kind data schema — functional spec

Companion to the
[Self-contained per-kind data schema](../0130-self-contained-data-schema.md)
change case. This spec states _what_ the change must do.

The key words **MUST**, **MUST NOT**, **SHOULD**, and **MAY** are to be
interpreted as described in IETF RFC 2119.

## Scope

Covered: the output of `qualitymd evaluation data schema <kind>`, its terminal
rendering, the source-of-truth invariant it shares with `data set` and
`data example`, the preservation of the two distinct `schema` command surfaces,
and the `/quality` skill guidance that consumes the schema.

Deferred: a unified `qualitymd schema <kind>` command, a `data describe` command,
and any annotation of `data example` output — all rejected in the
[change case](../0130-self-contained-data-schema.md#scope).

## Requirements

### Self-contained single-kind schema

- `evaluation data schema <kind>` **MUST** emit a schema from which the requested
  kind's **required fields** and each closed field's **allowed enum values** can
  be read at the document root, without dereferencing a `$ref` into a separate
  `$defs` definition for the requested kind. An author reading the top level
  **MUST** be able to see, for example, that analysis `status` is one of
  `analyzed | empty | not_analyzed | blocked` and that `areaId` takes the
  `area:<id>` reference form.

  > Durable spec: modify `specs/cli/evaluation-data.md` — replace "rooted at that
  > kind" with the self-contained-legibility requirement.

- The single-kind schema **MUST** impose the same constraints as before — the
  flattening is presentational, not a contract change. Any payload that `data set`
  accepts or rejects for that kind **MUST** validate or fail identically against
  the emitted single-kind schema.

  > Durable spec: none — restates the existing validation contract; no durable
  > spec text changes.

- The no-argument `data schema` form **MUST** continue to emit a schema covering
  the full data surface (every kind). It **MAY** retain `$defs`/`$ref` structure
  for that multi-kind document.

  > Durable spec: none — preserves the current no-argument `data schema` behavior
  > already specified.

- The emitted schema **MUST** remain derived from the one typed contract that
  backs `data set` validation and `data example`. Enum value sets and required
  fields in the schema **MUST** stay identical to those the validator enforces;
  there **MUST NOT** be a second hand-maintained source.

  > Durable spec: none — already required by
  > `specs/evaluation/records/payload-kinds.md` (one typed source of truth for
  > schema output and examples); relied on unchanged.

### Command surface

- `qualitymd schema` (QUALITY.md frontmatter) and `qualitymd evaluation data
schema` (evaluation payloads) **MUST** remain separate commands. This change
  **MUST NOT** introduce a unified `qualitymd schema <kind>` whose kind set spans
  both the model frontmatter and the evaluation payload kinds.

  > Durable spec: none — the two commands already have separate sub-specs
  > (`specs/cli/schema.md`, `specs/cli/evaluation-data.md`); this keeps that
  > boundary and adds no text.

- `evaluation data schema` **MUST NOT** offer a second JSON result-wrapper mode;
  its stdout is already the JSON artifact (unchanged from the current contract).

  > Durable spec: none — already specified in `specs/cli/evaluation-data.md` (the
  > `--json` carve-out).

### Terminal rendering (optional alignment)

- When output must be plain — stdout is not a terminal, or `NO_COLOR` is set —
  `evaluation data schema` **MUST** write the schema as verbatim JSON and nothing
  else, so a redirect reproduces the artifact byte-for-byte.
- On a terminal, `evaluation data schema` **MAY** syntax-highlight and page the
  JSON for readability, consistent with `qualitymd schema`. This rendering
  **MUST NOT** change the bytes written when output must be plain.

  > Durable spec: modify `specs/cli/evaluation-data.md` — add the terminal-rendering
  > allowance (verbatim on plain; optional highlight/page on a terminal) consistent
  > with `specs/cli/schema.md`.

### Skill guidance

- The bundled `/quality` skill **MUST** direct agents to read
  `data schema <kind>` as the source for a kind's required fields and allowed enum
  values, and to treat `data example <kind>` as a concrete valid instance only —
  so closed value sets are not discovered by triggering `data set --dry-run`
  validation failures.

  > Durable spec: modify
  > `specs/skills/quality-skill/quality-skill.md` — sharpen the payload-discovery
  > `MUST` (which already designates `data schema [<kind>]` "the authoritative
  > payload contract" and bars dry-run shape-discovery) to name _required fields
  > and enum value sets_ as the schema's job and `data example` as a concrete
  > instance only. The CLI legibility fix above is what makes it followable.

## Durable spec changes

Rollup of the per-requirement `Durable spec:` annotations above (those are
authoritative) — the [`specs/`](../../../specs/index.md) bundle and
[`SPECIFICATION.md`](../../../SPECIFICATION.md). No change here is cross-cutting;
every entry traces to a single requirement.

### To add

None.

### To modify

- `specs/cli/evaluation-data.md` — self-contained `data schema <kind>` legibility,
  plus the terminal-rendering allowance (per the self-contained-schema and
  terminal-rendering requirements above). Promote the
  [motivation](../0130-self-contained-data-schema.md#motivation) — constraints were
  discovered via dry-run because the schema was unreadable — into the requirement's
  rationale.
- `specs/skills/quality-skill/quality-skill.md` — name `data schema <kind>` as the
  source for required fields and enum values, `data example` as an instance (per
  the skill-guidance requirement above).

### To rename

None.

### To delete

None.
