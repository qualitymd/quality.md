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
[QUALITY.md format specification](../SPECIFICATION.md): where this document
constrains the *tool*, the format spec constrains the *file*.

This document uses BCP 14 keywords only for testable conformance requirements.
The key words "MUST", "MUST NOT", "SHOULD", and "MAY" are to be interpreted as
described in [RFC 2119](../docs/reference/rfc2119.md) and
[RFC 8174](../docs/reference/rfc8174.md) when, and only when, they appear in all
capitals.

## Scope

This phase covers the deterministic CLI surface: format-tooling commands that
operate on QUALITY.md files, plus evaluation-run commands that scaffold,
validate, and render runtime records. The CLI never calls a model; skills carry
judgment and pass judgment payloads to the deterministic surface.

| Command                                           | Purpose                                                  |
| ------------------------------------------------- | -------------------------------------------------------- |
| [`init`](./cli/init.md)                           | Scaffold a starter `QUALITY.md` to fill in.              |
| [`lint`](./cli/lint.md)                           | Validate a file's structure against the format spec.     |
| [`spec`](./cli/spec.md)                           | Emit the QUALITY.md format specification.                |
| [`status`](./cli/status.md)                       | Emit a deterministic project-state snapshot.             |
| [`evaluation create`](./cli/evaluation-create.md) | Create a numbered evaluation run folder.                 |
| [`evaluation data`](./cli/evaluation-data.md)     | Persist and inspect Evaluation structured data.          |
| [`evaluation list`](./cli/evaluation-list.md)     | List evaluation runs.                                    |
| [`evaluation status`](./cli/evaluation-status.md) | Inspect whether a run can be rendered.                   |
| [`evaluation report`](./cli/evaluation-report.md) | Build evaluation reports.                                |
| [`version`](./cli/version.md)                     | Show structured CLI and bundled spec version metadata.   |
| [`update`](./cli/update.md)                       | Apply or check for CLI updates through managed channels. |

## Commands

- [`init`](./cli/init.md) — scaffold a starter `QUALITY.md`.
- [`lint`](./cli/lint.md) — validate a file's structure.
- [`spec`](./cli/spec.md) — emit the format specification.
- [`status`](./cli/status.md) — emit a deterministic project-state snapshot.
- [`evaluation create`](./cli/evaluation-create.md) — create a numbered
  evaluation run folder.
- [`evaluation data`](./cli/evaluation-data.md) — persist and inspect Evaluation
  structured data.
- [`evaluation list`](./cli/evaluation-list.md) — list evaluation runs.
- [`evaluation status`](./cli/evaluation-status.md) — inspect renderability.
- [`evaluation report`](./cli/evaluation-report.md) — build evaluation reports.
- [`version`](./cli/version.md) — show structured CLI version metadata.
- [`update`](./cli/update.md) — apply or check for CLI updates.

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
- **MUST** emit plain human output with no color or terminal escape sequences
  when the destination stream is not a terminal, honoring `NO_COLOR` across
  stdout and stderr human surfaces. The Fang / Lip Gloss stack's idiomatic
  non-TTY behavior is the expected mechanism for this.

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
baseline. `--json` is a near-universal default across commands, with the detailed
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

**Structured input.** Commands that read an author-supplied structured payload
from stdin or a documented path option **MUST** make that payload contract
discoverable from inside the tool. Command help **MUST** document every payload
field, its JSON type, whether it is required, and allowed values for enum-like
fields, and **MUST** include at least one complete valid payload example.

Side-effecting structured-input commands **MUST** offer `-n/--dry-run`, which
validates the payload and reports the records or files that would be written
without creating, replacing, numbering, or otherwise persisting records. Under
`--json`, a valid dry-run **MUST** emit the same receipt shape as the real write
with a dry-run marker and intended paths. Invalid dry-runs use the usage-error
exit category.

Structured-input validation **MUST** aggregate every problem the command can
detect in one pass. Each problem **MUST** name the offending field in caller
vocabulary: the JSON key path from the payload, not an internal Go struct field
or type name. Where useful, the diagnostic should state the expected JSON type
and allowed enum values.

**`--json` for machine-readable output.** Commands should offer
machine-readable output, spelled `--json` — never `--format json` or a
per-command variant. Machine-readable output is the broad default wherever it is
meaningful, so an agent can reach for `--json` without reasoning per command
about whether the flag exists.

Commands default to human-readable terminal output. There is no format
auto-detection: passing `--json` is the only way to switch a command to emitting
a JSON document on stdout, the form agents and CI consume.

A command can omit `--json` only when its output is a verbatim artifact that
*is* the payload and is meant to be redirected, so wrapping it adds nothing. For
example, [`spec`](./cli/spec.md) emits the format specification itself.

Commands whose stdout is already a JSON artifact, such as a JSON Schema, example
payload, or stored JSON record, **MUST NOT** define a second JSON result-wrapper
mode. Such commands **MAY** recognize `--json` only to fail with a targeted usage
error explaining that the command already emits JSON on stdout and should be
rerun without `--json`.

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

**Human output styling.** A command's default human output is styled for the
terminal — color, weight, and a leading status glyph — drawn from a single brand
palette shared with the harness's help, error, and version rendering, so the
whole tool reads as one program. Styling is a terminal-only convenience layered
over a canonical plain form:

- The plain form is the source of truth. Styling **MUST NOT** change the words,
  order, or facts the plain form carries; it only adds color and glyphs.
- Styling applies **only** when the destination stream is a terminal and
  `NO_COLOR` is unset, per the [baseline](#baseline). Piped, redirected, and
  agent-driven output is the plain form, byte-for-byte.
- A command can render a long verbatim artifact through the user's pager
  (`$PAGER`, else a system default) when stdout is a terminal. Paging is never
  load-bearing: it is skipped when stdout is not a terminal or no pager is
  available, and the artifact is written directly instead.

**Binary version.** `qualitymd --version` reports the version stamped into the
binary at release. For a binary built without that stamp — `go install` or a
local build — it **MUST** fall back to the module build information the Go
toolchain embeds (the module version for an installed release, otherwise a
development label carrying the VCS revision). If a local development invocation
does not expose the revision through embedded build information but can resolve
it from the local VCS checkout, it should use that revision rather than
reporting a bare placeholder.

`qualitymd version --json` **MUST** expose the CLI version, commit when known,
development-build state, and bundled `SPECIFICATION.md` version as structured
metadata. `qualitymd update --check` **MUST** be the explicit non-mutating
update-check surface. Ordinary commands **MAY** refresh a local update cache in a
bounded, best-effort background process and **MAY** show a cached
update-available notice on stderr only, per
[`update notice`](./cli/update-notice.md); they **MUST NOT** block on a network
fetch for update discovery or include update notices in stdout or
machine-readable output.

## To be specified

- The shared invocation form and the file / stdin argument convention.
- Global flags common to every command.
- Machine-readable output stability across commands, and the verbosity/quiet
  control that governs human noise.
