---
type: Design Doc
title: spec command — design doc
description: How qualitymd embeds and emits the bundled format specification.
tags: [cli, command, spec, design]
timestamp: 2026-06-17T00:00:00Z
---

# spec command — design doc

Design behind the
[Specify and implement the spec command](../0006-spec-command.md) change and its
[functional spec](spec.md).

## Context

`qualitymd spec` needs to emit the repository-root `SPECIFICATION.md` from the
binary itself. The byte-for-byte redirected output is the command's primary
agent-facing payload; terminal rendering is only a human convenience.

## Approach

Add a tiny root package, `github.com/qualitymd/quality.md`, that embeds
`SPECIFICATION.md` with `go:embed` and exposes `Specification() []byte`.
`go:embed` patterns cannot reach parent directories, so the embedding Go source
lives beside the artifact it embeds rather than under `internal/cli`.

`internal/cli` owns the command. `newSpecCmd` accepts no arguments and defines no
flags, so Cobra reports unexpected arguments and unknown flags through the same
usage-error path as the other commands. The command writes to `cmd.OutOrStdout()`
and lets write failures return normally, which the existing CLI boundary maps to
exit `70`.

Output selection is deliberately narrow:

- If `NO_COLOR` is set, or stdout is not an `Fd()`-bearing terminal, write the
  embedded bytes directly.
- If stdout is a terminal, render the Markdown with Charm's Glamour terminal
  renderer and write the rendered bytes.

This keeps the redirect path free of formatting, color, or diagnostic concerns.

## Alternatives

- **Embed from `internal/cli`.** Rejected because `go:embed` cannot embed
  `../../SPECIFICATION.md`.
- **Copy the specification under `internal/`.** Rejected because it creates a
  second artifact that can drift from the public source of truth.
- **Hand-render Markdown with Lip Gloss styles.** Rejected because it would
  duplicate a Markdown renderer poorly. Glamour is the Charm stack's purpose-built
  terminal Markdown renderer.
- **Always write verbatim Markdown.** Rejected because the spec asks for
  formatted terminal output when stdout is a terminal, and Glamour makes that
  cheap without touching redirected output.

## Trade-offs & risks

- The root package exists only to expose the embedded artifact. That is small but
  visible to importers; if the public API grows later, keep it intentional.
- Glamour adds dependencies. The command uses it only on the TTY path, so the
  deterministic redirected payload does not depend on renderer behavior.
- Terminal rendering is best-effort human output. Tests pin the verbatim payload
  and error behavior; renderer styling is left to Glamour.

## Open questions

None.
