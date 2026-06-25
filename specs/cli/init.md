---
type: Functional Specification
title: qualitymd init
description: Scaffold a starter `QUALITY.md` to fill in.
tags: [cli, command, init]
timestamp: 2026-06-17T00:00:00Z
---

# qualitymd init

`qualitymd init` scaffolds a starter `QUALITY.md` for an author to fill in. It
inherits the cross-cutting CLI contract from the [CLI spec](../cli.md); this
sub-spec covers only what is particular to `init`.

This document uses BCP 14 keywords only for testable conformance requirements.
The key words "MUST", "MUST NOT", "SHOULD", and "MAY" are to be interpreted as
described in [RFC 2119](../../docs/reference/rfc2119.md) and
[RFC 8174](../../docs/reference/rfc8174.md) when, and only when, they appear in
all capitals.

## Scope

Covered: the output area and stdout piping, overwrite protection, the scaffold
variant selection, and the contents of the scaffold `init` produces, including
its `--json` result receipt.

## Requirements

### Output target

- `init` **MUST** write the scaffold to a file named `QUALITY.md` in the current
  working directory by default.
- `init` **MUST** accept an optional path argument naming the file to write
  instead, so an author can scaffold to a chosen location.
- A path argument of `-` **MUST** write the scaffold to standard output rather
  than a file. When emitting to standard output, overwrite protection does not
  apply; nothing on disk is read or replaced.
- `init --json -` **MUST** fail with a usage error because `-` selects the raw
  scaffold on stdout while `--json` selects a receipt on stdout.

### Overwrite protection

- When the model file already exists, `init` **MUST NOT** overwrite it: it
  **MUST** leave the file untouched and exit non-zero, reporting that the file
  already exists.
- A `--force` flag **MUST** permit overwriting an existing model file.

### Scaffold variant

- By default `init` writes the *guided* scaffold: a structurally valid skeleton
  whose frontmatter placeholders carry inline YAML comments and whose body
  includes the recommended headed sections, so the file reads as a guided
  starting point for an author filling it in by hand.
- A `--minimal` flag **MUST** write a *minimal* scaffold instead: the same
  structurally valid frontmatter skeleton — seeded rating scale, placeholder
  factor, requirement, and assessment, with non-empty placeholder `title`
  values — but **without** the guided inline-comment prose and **without** the
  headed body-section stubs, keeping only the top-level heading naming the root
  area. This is for an authoring pass — agent or human — that will replace the
  frontmatter wholesale and would otherwise discard the guided prose.
- Both variants **MUST** satisfy every requirement in
  [Scaffold contents](#scaffold-contents) that does not concern the guided
  inline comments or the headed body sections; in particular both **MUST** pass
  `qualitymd lint` with no errors and seed the suggested four-level rating
  scale.
- The variant selection applies equally to a file write and to the `-`
  stdout form.

### Scaffold contents

The requirements in this section describe the guided scaffold. The minimal
scaffold relaxes only the inline-comment and headed-body-section requirements,
per [Scaffold variant](#scaffold-variant).

- The scaffold **MUST** be a structurally valid QUALITY.md file per the
  [format specification](../../SPECIFICATION.md): a freshly initialized file
  **MUST** pass `qualitymd lint` with no errors, placeholder text and all.
- The frontmatter **MUST** seed the suggested four-level rating scale
  (`outstanding` / `target` / `minimum` / `unacceptable`) from the
  [format spec](../../SPECIFICATION.md#model), which a scaffolding tool can seed.
- The scaffold **MUST** seed non-empty placeholder `title` values wherever the
  format requires them, including the model root, every seeded rating level, the
  placeholder factor, the placeholder requirement, and any commented area/factor
  or requirement examples.
- The frontmatter **MUST** include a minimal skeleton that leads to at least one
  requirement: a placeholder factor carrying a `description`, and beneath it a
  placeholder requirement carrying a `title` and a single non-empty
  `assessment`.
- Frontmatter placeholders should be accompanied by inline YAML comments that
  tell the author what to replace, so the file reads as a guided skeleton rather
  than a filled-in example.
- The Markdown body **MUST** include the
  [recommended body sections](../../SPECIFICATION.md#body-semantics) — Overview,
  Scope, Needs, and Risks — as headed stubs, each with a brief prompt of what it
  captures and a place for the section's unknowns, open questions, and review
  state line. The starter body should teach that material support unavailable to
  the evaluating agent belongs in the relevant section's unknowns or open
  questions.
- The body should open with a top-level heading naming the root area, matching the
  scaffolded model `title`.

### Reporting

- On successfully writing a file, `init` **MUST** report the created path. This
  confirmation is written to standard error, so standard output carries only the
  scaffold itself when piping with `-`.
- `init` should close with
  [next actions](../cli.md#conventions) pointing the author to validate and then
  edit the new file, for example `qualitymd lint QUALITY.md`.
- Under `--json`, `init` **MUST** emit a single result receipt on stdout instead
  of the human confirmation. The receipt **MUST** include:
  - `schemaVersion` — `1` until the receipt shape changes incompatibly;
  - `path` — the file path written;
  - `created` — `true` when a new file was created and `false` when `--force`
    overwrote an existing file; and
  - `nextActions` — the same next-action data that would render as the human
    footer.
- When `init --json` refuses to overwrite an existing file, it **MUST NOT** emit
  a success receipt. It should emit a JSON error object on stderr with
  `schemaVersion`, `path`, and `reason`; the exit code remains the CLI spec's
  internal-error category.

## Deferred

- Interactive prompting for the title, factors, or requirements.
- Selectable scaffold templates, including alternate rating scales or domain
  profiles.
- Directory-target semantics. The path argument names the file to write; how
  `init` behaves when handed a directory is left for later.
