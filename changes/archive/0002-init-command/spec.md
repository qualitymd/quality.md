---
type: Functional Specification
title: init command — functional spec
description: What qualitymd init scaffolds, where it writes, and how it refuses to overwrite.
tags: [cli, command, init]
timestamp: 2026-06-17T00:00:00Z
---

# init command — functional spec

Companion to the [Specify the init command](../0002-init-command.md) change.
This spec states _what_ `qualitymd init` must do; it is the delta that the
durable [`specs/cli/init.md`](../../../specs/cli/init.md) sub-spec absorbs when the
change lands.

`init` inherits the cross-cutting CLI contract — invocation, global flags,
output conventions, and exit-code semantics — from the
[CLI spec](../../../specs/cli.md); this spec covers only what is particular to
`init`. The shape of a valid `QUALITY.md` is fixed by the
[format specification](../../../SPECIFICATION.md) and is not restated here.

The key words **MUST**, **MUST NOT**, **SHOULD**, and **MAY** are to be
interpreted as described in IETF RFC 2119.

## Scope

Covered: the output target and stdout piping, overwrite protection, and the
contents of the scaffold `init` produces. Cross-cutting CLI behavior is deferred
to the [CLI spec](../../../specs/cli.md); further `init`-specific behavior is listed
under [Deferred](#deferred).

Per the [CLI spec's `--json` convention](../../../specs/cli.md#conventions), `init`
offers no `--json` output.

## Requirements

### Output target

- `init` **MUST** write the scaffold to a file named `QUALITY.md` in the current
  working directory by default.
- `init` **MUST** accept an optional path argument naming the file to write
  instead, so an author can scaffold to a chosen location.
- A path argument of `-` **MUST** write the scaffold to standard output rather
  than a file, so it can be piped or redirected. When emitting to standard
  output, the [overwrite protection](#overwrite-protection) does not apply —
  nothing on disk is read or replaced.

### Overwrite protection

- When the target file already exists, `init` **MUST NOT** overwrite it: it
  **MUST** leave the file untouched and exit non-zero, reporting that the file
  already exists. _(An author's existing model is never silently clobbered.)_
- A `--force` flag **MUST** permit overwriting an existing target file.

### Scaffold contents

- The scaffold **MUST** be a structurally valid `QUALITY.md` per the
  [format specification](../../../SPECIFICATION.md): a freshly initialized file
  **MUST** pass [`qualitymd lint`](../../../specs/cli/lint.md) with no errors,
  placeholder text and all.
- The frontmatter **MUST** seed the suggested four-level rating scale
  (`outstanding` / `target` / `minimum` / `unacceptable`) from the
  [format spec](../../../SPECIFICATION.md#model), which a scaffolding tool MAY seed.
- The frontmatter **MUST** include a minimal skeleton that leads to at least one
  requirement — a placeholder factor carrying a `description`, and beneath it a
  placeholder requirement carrying a single non-empty `assessment` — so the
  scaffold is valid and illustrates the target → factor → requirement shape.
- Frontmatter placeholders **SHOULD** be accompanied by inline YAML comments that
  tell the author what to replace, so the file reads as a guided skeleton rather
  than a filled-in example.
- The Markdown body **MUST** include the
  [recommended body sections](../../../SPECIFICATION.md#markdown-body) — Overview,
  Scope, Needs, Risks, and Known gaps — as headed stubs, each with a brief prompt
  of what it captures.
- The scaffold **SHOULD** seed a placeholder `title`, and the body **SHOULD**
  open with a top-level heading naming the subject, matching that `title`.

### Reporting

- On successfully writing a file, `init` **MUST** report the created path. This
  confirmation is written to standard error, so that standard output carries only
  the scaffold itself when piping with `-`.
- `init` **SHOULD** close with [next actions](../../../specs/cli.md#conventions)
  pointing the author to validate and then edit the new file — for example
  `qualitymd lint QUALITY.md`.

## Deferred

- **Cross-cutting CLI behavior** — global flags, exit-code semantics, and quiet
  verbosity — specified in the [CLI spec](../../../specs/cli.md).
- **Interactive prompting** for the title, factors, or requirements: the scaffold
  is produced non-interactively in this phase.
- **Selectable scaffold templates** — alternate rating scales or domain profiles
  chosen by flag. One commented skeleton is seeded; variants wait for demand.
- **Directory-target semantics** — the path argument names the file to write; how
  `init` behaves when handed a directory is left for later.
