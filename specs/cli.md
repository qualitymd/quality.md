---
type: Functional Specification
title: qualitymd CLI
description: High-level requirements for the deterministic qualitymd command-line surface.
tags: [cli, specification]
timestamp: 2026-06-16T00:00:00Z
---

# qualitymd CLI

This document specifies the high-level requirements for the `qualitymd`
command-line interface — invocation-wide behavior, output conventions, exit
codes, and agent accessibility independent of any one command's behavior.
Per-command behavior lives in the command sub-specs under [Commands](#commands).

It is a companion to the
[`QUALITY.md` format specification](../SPECIFICATION.md): where this document
constrains the *tool*, the format spec constrains the *file*.

The key words **MUST**, **MUST NOT**, **SHOULD**, and **MAY** are to be
interpreted as described in IETF RFC 2119.

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

## Agent accessibility

Agent accessibility is the property that lets a non-interactive caller — an
agent, CI, or another tool — drive a command reliably. This is the CLI's role in
`qualitymd`: the deterministic, mechanical surface that skills and automation
evaluate through.

The contract has two tiers:

- **Baseline.** Every command within this spec's [Scope](#scope) owes these
  invariants. The baseline is not opt-in and has no per-command exemption.
- **Opt-in capabilities.** A caller activates these per invocation, when a
  command offers them: `--json`, `nextActions`, and the still-to-be-specified
  quiet/verbosity control.

### Baseline

Every in-scope command:

- **MUST** run non-interactively. It never blocks on a prompt and never assumes a
  TTY. When required input is absent, it fails with a diagnostic and a usage exit
  code rather than waiting.
- **MUST** keep stdout as the payload and stderr for diagnostics, progress, and
  human next-action footers, so redirecting or piping stdout is never polluted.
- **MUST** be deterministic: the same input and file state produce the same
  output, with no timestamps, ordering jitter, or sampling in the payload.
- **MUST** signal its outcome through the exit-code categories below.
- **MUST** emit plain output with no color or terminal escape sequences when
  stdout is not a terminal, honoring `NO_COLOR`. The Fang / Lip Gloss stack's
  idiomatic non-TTY behavior is the expected mechanism for this.

### Exit Codes

`qualitymd` exits `0` only on success. Non-zero outcomes use these stable,
documented categories so callers can branch without parsing output:

| Code | Category               | Meaning                                                                                                                           |
| ---- | ---------------------- | --------------------------------------------------------------------------------------------------------------------------------- |
| `0`  | Success                | The command did its job and found nothing to report.                                                                              |
| `1`  | Ran but found problems | The command completed normally but its result is a reportable negative, such as `lint` finding error-severity findings.           |
| `2`  | Usage error            | The invocation was malformed: unknown flag, bad argument, incompatible options, or unknown command.                               |
| `70` | Internal error         | The command could not complete the requested action: I/O failure, unmet precondition such as guarded overwrite refusal, or a bug. |

### Opt-in Capabilities

Agent-facing enrichments are opt-in per invocation and do not weaken the
baseline. `--json` is a near-universal SHOULD across commands, with the detailed
rules in [Conventions](#conventions). `nextActions` are offered only when a
command has a useful follow-up; their rendering is also defined in
[Conventions](#conventions). The quiet/verbosity control that governs human
noise remains deferred.

A command that omits `--json` because its output is a verbatim artifact still
owes the full baseline, including categorized exit codes and stderr diagnostics
on failure.

## Conventions

Conventions that hold wherever they apply, so flags and output behave the same
across commands.

**`--json` for machine-readable output.** Commands **SHOULD** offer
machine-readable output, spelled `--json` — never `--format json` or a
per-command variant. Machine-readable output is the broad default wherever it is
meaningful, so an agent can reach for `--json` without reasoning per command
about whether the flag exists.

Commands default to human-readable terminal output. There is no format
auto-detection: passing `--json` is the only way to switch a command to emitting
a JSON document on stdout, the form agents and CI consume.

A command **MAY** omit `--json` only when its output is a verbatim artifact that
*is* the payload and is meant to be redirected, so wrapping it adds nothing. For
example, [`spec`](./cli/spec.md) emits the format specification itself.

Under `--json`, a side-effecting command **MUST** emit a result receipt — a JSON
document describing what it did and carrying its `nextActions` in-band — rather
than its human prose. This makes `--json` meaningful for commands whose human
output is a confirmation rather than a result.

Among the current commands, [`init`](./cli/init.md) offers a receipt,
[`lint`](./cli/lint.md) offers its finding result, and [`spec`](./cli/spec.md)
is the verbatim-artifact carve-out.

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
- Versioning of the binary and the format spec it targets.
