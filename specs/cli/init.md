---
type: Functional Specification
title: qualitymd init
description: Scaffold a starter QUALITY.md to fill in.
tags: [cli, command, init]
timestamp: 2026-06-17T00:00:00Z
---

# qualitymd init

`qualitymd init` scaffolds a starter `QUALITY.md` for an author to fill in. It
inherits the cross-cutting CLI contract from the [CLI spec](../cli.md); this
sub-spec covers only what is particular to `init`.

The key words **MUST**, **MUST NOT**, **SHOULD**, and **MAY** are to be
interpreted as described in IETF RFC 2119.

## Scope

Covered: the output target and stdout piping, overwrite protection, and the
contents of the scaffold `init` produces, including its `--json` result receipt.

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

- When the target file already exists, `init` **MUST NOT** overwrite it: it
  **MUST** leave the file untouched and exit non-zero, reporting that the file
  already exists.
- A `--force` flag **MUST** permit overwriting an existing target file.

### Scaffold contents

- The scaffold **MUST** be a structurally valid `QUALITY.md` per the
  [format specification](../../SPECIFICATION.md): a freshly initialized file
  **MUST** pass `qualitymd lint` with no errors, placeholder text and all.
- The frontmatter **MUST** seed the suggested four-level rating scale
  (`outstanding` / `target` / `minimum` / `unacceptable`) from the
  [format spec](../../SPECIFICATION.md#model), which a scaffolding tool MAY seed.
- The frontmatter **MUST** include a minimal skeleton that leads to at least one
  requirement: a placeholder factor carrying a `description`, and beneath it a
  placeholder requirement carrying a single non-empty `assessment`.
- Frontmatter placeholders **SHOULD** be accompanied by inline YAML comments that
  tell the author what to replace, so the file reads as a guided skeleton rather
  than a filled-in example.
- The Markdown body **MUST** include the
  [recommended body sections](../../SPECIFICATION.md#markdown-body) — Overview,
  Scope, Needs, Risks, and Known gaps — as headed stubs, each with a brief prompt
  of what it captures.
- The scaffold **SHOULD** seed a placeholder `title`, and the body **SHOULD**
  open with a top-level heading naming the subject, matching that `title`.

### Reporting

- On successfully writing a file, `init` **MUST** report the created path. This
  confirmation is written to standard error, so standard output carries only the
  scaffold itself when piping with `-`.
- `init` **SHOULD** close with
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
  a success receipt. It **SHOULD** emit a JSON error object on stderr with
  `schemaVersion`, `path`, and `reason`; the exit code remains the CLI spec's
  internal-error category.

## Deferred

- Interactive prompting for the title, factors, or requirements.
- Selectable scaffold templates, including alternate rating scales or domain
  profiles.
- Directory-target semantics. The path argument names the file to write; how
  `init` behaves when handed a directory is left for later.
