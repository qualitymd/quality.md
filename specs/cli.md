---
type: Functional Specification
title: qualitymd CLI
description: High-level requirements for the deterministic qualitymd command-line surface.
tags: [cli, specification]
timestamp: 2026-06-16T00:00:00Z
---

# qualitymd CLI

> 🚧 **Placeholder.** This spec is a stub. It will specify the cross-cutting
> contract every `qualitymd` command shares — invocation, flags, output, exit
> codes, and agent accessibility — independent of any one command's behavior.
> Per-command behavior lives in the command sub-specs under [Commands](#commands).

This document will specify the high-level requirements for the `qualitymd`
command-line interface. It is a companion to the
[`QUALITY.md` format specification](../SPECIFICATION.md): where this document
constrains the *tool*, the format spec constrains the *file*.

## Scope

This phase covers three commands — the **format-tooling layer**: commands that
operate on a single `QUALITY.md` file, hold no evaluation state, and never call a
model.

| Command                 | Purpose                                              |
| ----------------------- | ---------------------------------------------------- |
| [`init`](./cli/init.md) | Scaffold a starter `QUALITY.md` to fill in.          |
| [`lint`](./cli/lint.md) | Validate a file's structure against the format spec. |
| [`spec`](./cli/spec.md) | Emit the `QUALITY.md` format specification.          |

**Deferred.** The deeper surface that records per-target verdicts, rolls them up
the target tree, and gates CI on the outcome is out of scope for this phase. The
judgment-based evaluation defined in the format spec's
[Evaluation](../SPECIFICATION.md#evaluation) is carried by skills that orchestrate
the model, not by these commands.

## Commands

- [`init`](./cli/init.md) — scaffold a starter `QUALITY.md`.
- [`lint`](./cli/lint.md) — validate a file's structure.
- [`spec`](./cli/spec.md) — emit the format specification.

## Technical requirements

These constraints bind every command regardless of its behavior.

**Idiomatic to the tech stack.** Every functional requirement in this spec and
its command sub-specs must be satisfiable through the idiomatic capabilities of
the chosen CLI stack — Go, with [Cobra](https://github.com/spf13/cobra) for
command and flag structure, [Charm Fang](https://github.com/charmbracelet/fang)
for the invocation harness, and [Lip Gloss](https://github.com/charmbracelet/lipgloss)
for terminal output. A requirement that can only be met by working against the
grain of these libraries — hand-rolling flag parsing, replacing the framework's
help and error rendering, or reimplementing what the harness already provides —
is a signal to reshape the requirement, not the stack. Where a command
deliberately diverges from a stack default, its sub-spec should say so and why.

## Conventions

Conventions that hold wherever they apply, so flags and output behave the same
across commands. They do not force a behavior onto every command — they fix the
form a behavior takes once a command opts into it.

**`--json` for machine-readable output.** Not every command offers a
machine-readable form, but where one does, it is spelled `--json` — never
`--format json` or a per-command variant. Commands default to human-readable
terminal output; passing `--json` switches the command to emitting a single JSON
document on stdout, the form agents and CI consume.

A command should offer `--json` when all of the following hold:

- Its primary job is to *produce results a caller reads*, rather than to perform
  a side effect or emit a verbatim artifact.
- Those results have, or can have, a stable documented schema worth committing to.
- A non-interactive consumer — an agent, CI, another tool — would plausibly
  branch on the *contents* of the output, not just the exit code.

A command should not offer `--json` when its output is already a
machine-consumable artifact that *is* the payload (so wrapping it in JSON adds
nothing), when the command is purely side-effecting and its outcome is fully
carried by the exit code, or when the output is free-form human guidance with no
stable schema worth freezing. Among the current commands, [`lint`](./cli/lint.md)
offers `--json` (its findings); [`spec`](./cli/spec.md) and [`init`](./cli/init.md)
do not — the spec text and the scaffolded file are themselves the payload.

**Suggested next actions.** A command may close its response with a short list of
*next actions* — the commands a caller would most plausibly run next, given what
just happened. They are advisory: they never change behavior or the exit code,
and a command with no useful follow-up omits them.

- *Concrete over vague.* Each action is, where possible, a runnable command the
  caller can copy (`qualitymd lint QUALITY.md`), not prose like "consider
  validating your file."
- *Deterministic.* The same outcome yields the same suggestions; they are derived
  from the command's result, not ranked or sampled.
- *Subordinate to the payload.* In human output they render as a distinct footer
  on stderr, after the primary output, so redirecting or piping stdout
  (`qualitymd spec > SPECIFICATION.md`) is never polluted by them. Under `--json`
  they instead appear in-band as a `nextActions` array in the document, so agents
  receive them as data; a command that does not offer `--json` still shows the
  human footer.
- *Useful on success and failure alike.* After a successful `init`, the next
  action points to linting or editing the new file; after a `lint` failure, to
  the command that re-runs the check or opens the offending file.

Each action carries a stable machine `id`, a human-readable `label`, and — when
the action is a single command — the exact `command` string:

```json
"nextActions": [
  { "id": "lint", "label": "Validate the scaffolded file", "command": "qualitymd lint QUALITY.md" }
]
```

Suppressing the footer rides on the global verbosity/quiet convention (still to
be specified); the `nextActions` data under `--json` is always present when a
command has next actions to offer.

## To be specified

- The shared invocation form and the file / stdin argument convention.
- Global flags common to every command.
- Output formats (human and machine-readable) and their stability.
- Exit-code semantics.
- Agent-accessibility and CI requirements.
- Versioning of the binary and the format spec it targets.
