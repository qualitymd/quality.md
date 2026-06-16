---
type: Specification
title: qualitymd CLI
description: High-level requirements for the deterministic qualitymd command-line surface.
tags: [cli, specification]
timestamp: 2026-06-16T00:00:00Z
---

# qualitymd CLI

**Version 0.1 — Draft**

This document specifies the high-level requirements for the `qualitymd`
command-line interface: the contract every command shares — invocation,
flags, output, exit codes, and agent accessibility — independent of any one
command's behavior. Per-command behavior is specified in the command
sub-specs linked under [Commands](#commands).

This is a companion to the [`QUALITY.md` format specification](../SPECIFICATION.md),
which is the source of truth for the file format and evaluation semantics. Where
this document constrains the *tool*, the format spec constrains the *file*.

The key words "MUST", "MUST NOT", "REQUIRED", "SHOULD", "SHOULD NOT",
"RECOMMENDED", "MAY", and "OPTIONAL" are to be interpreted as described in IETF
RFC 2119, and apply here to a conforming implementation of the `qualitymd` CLI.

## Scope

This phase specifies three commands — the **format-tooling layer**: the commands
that operate on a single `QUALITY.md` file, hold no evaluation state, and never
call a model.

| Command                 | Purpose                                              |
| ----------------------- | ---------------------------------------------------- |
| [`init`](./cli/init.md) | Scaffold a starter `QUALITY.md` to fill in.          |
| [`lint`](./cli/lint.md) | Validate a file's structure against the format spec. |
| [`spec`](./cli/spec.md) | Emit the `QUALITY.md` format specification.          |

These three form a closed loop: `init` produces a well-formed file, `lint`
validates one, and `spec` emits the rules a file is validated against.

**Deferred.** The deeper surface that records per-target verdicts, rolls them up
the target tree, and gates CI on the outcome is **out of scope for this phase**
(see [Deferred surface](#deferred-surface)). The deterministic CLI never performs
the judgment-based evaluation itself; that is carried by skills that orchestrate
the model (see the format spec's [Evaluation](../SPECIFICATION.md#evaluation)).

## Design requirements

The CLI contract rests on five properties. Every command MUST uphold them.

### Deterministic and model-free

The CLI MUST be deterministic and MUST NOT call a language model. For a given
input file and flags, a command MUST produce the same result every run, modulo
presentational styling (color, terminal width). This is the hard line that
separates the CLI (mechanics) from skills (judgment): the CLI scaffolds,
validates, and emits; it never judges quality.

### Agent-accessible

The CLI is consumed by coding agents as much as by people. Every command whose
output is structured MUST offer a machine-readable form via `--format json`
([Output formats](#output-formats)), and that form MUST follow a stable,
documented schema. Findings MUST be addressable — each carries a stable
identifier, a severity, and a location into the source file — so an agent can act
on them without parsing prose. See [Agent accessibility](#agent-accessibility).

### CI-friendly

Commands MUST be usable unattended in CI: non-interactive by default, no network
access required, and a meaningful [exit code](#exit-codes) that a gate can branch
on. A command MUST NOT block waiting for input on a terminal; actions that would
otherwise prompt are governed by explicit flags.

### Stateless and single-file

In this phase every command operates on a single `QUALITY.md` file and holds no
state between runs. The CLI reads the file (or stdin), does its work, and writes
to stdout. No database, lockfile, or run directory is read or written.

### Composable

Commands MUST behave as well-mannered Unix filters: read a path or stdin, write
the primary result to stdout, write diagnostics to stderr, and stay quiet on the
parts a caller did not ask for. This lets commands pipe into one another and into
standard tooling (`jq`, CI log collectors, an agent's tool harness).

## Invocation and global conventions

### Command form

```
qualitymd <command> [path] [flags]
```

The binary is named `qualitymd` (no `.md` suffix), so it carries none of the
shell/Windows file-association ambiguity a `.md`-suffixed bin name would. The
first positional argument is the command; the conventions below apply across all
commands unless a command sub-spec overrides them.

### File argument and stdin

Commands that read a `QUALITY.md` take the file as an OPTIONAL positional
argument:

- When omitted, the command defaults to `./QUALITY.md` in the working directory.
- A single `-` reads the file from standard input.

A command that cannot find its input file MUST exit with a usage error
([exit code](#exit-codes) `2`), not a silent success.

### Global flags

These flags MUST be accepted by every command (a command MAY ignore one that does
not apply to it):

| Flag            | Type              | Default | Description                                                             |
| --------------- | ----------------- | ------- | ----------------------------------------------------------------------- |
| `--format`      | `human` \| `json` | `human` | Output format ([Output formats](#output-formats)).                      |
| `-h`, `--help`  | boolean           | `false` | Print usage for the command and exit `0`.                               |
| `--version`     | boolean           | `false` | Print the CLI version and the format-spec version it targets, exit `0`. |
| `-q`, `--quiet` | boolean           | `false` | Suppress non-essential output; emit only the primary result and errors. |
| `--no-color`    | boolean           | `false` | Disable ANSI color in `human` output.                                   |

Commands MAY define additional flags in their sub-specs; those flags MUST NOT
redefine the meaning of a global flag.

### Output formats

Two output formats are defined:

- **`human`** (default) — readable, optionally colorized text for a person
  authoring or reviewing a model at a terminal. Layout and color are
  presentational and MAY change between versions.
- **`json`** — a single JSON document on stdout, following a stable schema
  documented by the command. This is the form agents and CI SHOULD consume. JSON
  output MUST NOT be colorized and MUST be the only thing written to stdout (logs
  and diagnostics go to stderr).

`human` is the default because the first audience for these commands is an author
writing a model by hand; agents and gates opt into `json` explicitly. A command
whose result is a single document (e.g. `spec`) MAY treat `human` as "the raw
document" and define `json` as a structured wrapper.

### Exit codes

| Code | Meaning                                                                                                     |
| ---- | ----------------------------------------------------------------------------------------------------------- |
| `0`  | Success — the command completed and found nothing that should fail a build.                                 |
| `1`  | Command-specific failure — e.g. `lint` found errors. The command ran correctly; the *subject* did not pass. |
| `2`  | Usage error — unknown command or flag, missing/unreadable input, malformed invocation.                      |

Exit codes are part of the contract and MUST be stable. A gate SHOULD treat any
non-zero code as a failure and MAY distinguish `1` (the file failed a check) from
`2` (the tool was invoked wrong).

### Color and terminals

`human` output MAY use ANSI color. Color MUST be disabled when `--no-color` is
set, when the `NO_COLOR` environment variable is present, or when stdout is not a
terminal. `json` output is never colorized.

### Network and filesystem

A command MUST NOT require network access to complete its core function and
SHOULD make no network calls at all. A command MUST NOT write outside paths the
user named (e.g. `init`'s target file); the format-tooling commands read their
input and write their result, nothing more.

## Agent accessibility

Because a primary consumer is a coding agent driving the evaluation loop, the CLI
states agent-facing requirements explicitly:

- **Structured output.** Every command with structured results MUST support
  `--format json` with a stable, documented schema, suitable for direct
  consumption by an agent's tool harness without scraping human text.
- **Addressable findings.** A finding emitted by a command MUST carry a stable
  rule identifier, a severity, a human-readable message, and a location into the
  source file (path plus, where meaningful, a line and/or a pointer into the YAML
  structure) so the agent can navigate to and fix it.
- **Spec on demand.** The [`spec`](./cli/spec.md) command MUST be able to emit the
  `QUALITY.md` format specification to stdout, so an agent can inject the current
  format rules into its context without reaching for an external copy.
- **Non-interactive.** No command may block on an interactive prompt. Any action
  that would otherwise require confirmation (overwriting a file, for instance)
  MUST be governed by an explicit flag, so an agent or CI job can run the command
  unattended and predictably.
- **Stream-friendly.** Commands MUST accept input on stdin (`-`) and write their
  primary result to stdout, so they compose inside an agent's shell tooling.

## Versioning and stability

The CLI reports two versions via `--version`: the version of the `qualitymd`
binary and the version of the `QUALITY.md` format specification it targets. The
`json` output schemas and the [exit codes](#exit-codes) are part of the CLI's
stable contract; presentational `human` output is not and MAY change between
releases. Breaking changes to a `json` schema or to exit-code semantics SHOULD be
reflected in a CLI major-version bump.

## Commands

Each command is specified in its own sub-spec:

- [`init`](./cli/init.md) — scaffold a starter `QUALITY.md`.
- [`lint`](./cli/lint.md) — validate a file's structure.
- [`spec`](./cli/spec.md) — emit the format specification.

## Deferred surface

The following is intentionally **out of scope** for this phase and is recorded
here so its absence reads as deliberate:

- **Evaluation state** — recording per-requirement findings and rating results,
  rolling them up the target tree, and persisting an evaluation run.
- **CI gating on a verdict** — failing a build on a rated outcome (as opposed to
  failing on structural lint errors, which `lint` already covers).

These belong to a later phase. Nothing in this phase records, rolls up, or gates
on a *quality verdict*; the deterministic CLI here is the format-tooling layer
only. The judgment-based evaluation defined in the format spec's
[Evaluation](../SPECIFICATION.md#evaluation) is carried by skills, not by these
commands.
