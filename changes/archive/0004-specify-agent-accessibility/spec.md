---
type: Functional Specification
title: agent accessibility — functional spec
description: The agent-accessibility contract the durable qualitymd CLI spec must define, the broadened --json convention, and the conformance the shipped commands must meet.
tags: [cli, specification, agent-accessibility]
timestamp: 2026-06-17T00:00:00Z
---

# agent accessibility — functional spec

Companion to the [Specify and enforce agent accessibility](../0004-specify-agent-accessibility.md)
change. This spec states two things: the delta the durable
[`qualitymd` CLI spec](../../../specs/cli.md) must absorb (a cross-cutting
**Agent accessibility** section and a revised `--json` convention), and the
conformance the shipped commands must meet once it lands.

The key words **MUST**, **MUST NOT**, **SHOULD**, and **MAY** are to be
interpreted as described in IETF RFC 2119.

## Scope

Covered: the content of the **Agent accessibility** section, the concrete
exit-code categories, the revised `--json` convention, the removal of the
matching **To be specified** entries, and the conformance changes (with tests)
needed to make the shipped `init` and `lint` commands meet the result.

Deferred: the global flag set, the quiet/verbosity control, the shared
invocation form, and the file/stdin argument convention — each remains its own
**To be specified** item. The unimplemented [`spec`](../../../specs/cli/spec.md)
command is out of scope; it inherits the baseline and the verbatim-artifact
carve-out by reference when a later change implements it.

## Requirements — the specification

### Framing

- The durable CLI spec **MUST** gain an **Agent accessibility** section that
  defines the contract as two tiers: a **baseline** every in-scope command owes,
  and **opt-in capabilities** a caller activates per invocation.
- The section **MUST** state that agent accessibility is the property that lets a
  non-interactive caller — an agent, CI, or another tool — drive a command
  reliably, consistent with the spec's premise that the CLI is the deterministic
  surface skills evaluate through.
- The section **MUST** make clear that the baseline is not opt-in: every command
  within the spec's [Scope](../../../specs/cli.md#scope) owes it, with no
  per-command exemption.

### Baseline (every command)

The section **MUST** state, as invariants binding every in-scope command, that a
command:

- **MUST** run non-interactively: it never blocks on a prompt and never assumes a
  TTY. When required input is absent it **MUST** fail with a diagnostic and a
  usage exit code rather than wait.
- **MUST** keep stdout as the payload and emit diagnostics, progress, and the
  next-actions footer on stderr, so redirecting or piping stdout is never
  polluted.
- **MUST** be deterministic: the same input and file state produce the same
  output, with no timestamps, ordering jitter, or sampling in the payload.
- **MUST** signal its outcome through the exit-code categories below.
- **MUST** emit plain output with no color or other terminal escape sequences
  when stdout is not a terminal, honoring `NO_COLOR`. The section **SHOULD** note
  this is satisfied by the idiomatic behavior of the Fang / Lip Gloss stack
  rather than hand-rolled detection.

### Exit-code categories

- The section **MUST** define a stable, documented set of exit-code categories
  that distinguishes at least:
  - **success** — the command did its job and found nothing to report;
  - **ran-but-found-problems** — the command completed normally but its result
    is a reportable negative (e.g. `lint` found findings);
  - **usage error** — the invocation was malformed (unknown flag, bad argument);
  - **internal error** — the command could not complete the requested action (an
    I/O failure, an unmet precondition such as a guarded overwrite refusal, or a
    bug).
- The section **MUST** assign each category a concrete, distinct exit code, so a
  caller can branch on the category without parsing output. *Ran-but-found-problems*
  **MUST** be distinguishable from *usage error* and *internal error*.
- The section **MUST** state that a command exits `0` only on *success*.

### The `--json` convention

The change **MUST** revise the `--json` convention in the CLI spec's
[Conventions](../../../specs/cli.md#conventions) from its current narrow
should-offer gate to a SHOULD-by-default, so that:

- Commands **SHOULD** offer `--json`. Machine-readable output is the broad
  default wherever it is meaningful, so an agent can reach for `--json` without
  reasoning per command about whether the flag exists.
- The default rendering **MUST** remain human; `--json` **MUST NOT** become the
  emitted-by-default format and there is no format auto-detection. A command
  switches to JSON only when `--json` is passed.
- A command **MAY** omit `--json` only when its output is a verbatim artifact
  that *is* the payload and is meant to be redirected (so wrapping it adds
  nothing) — for example [`spec`](../../../specs/cli/spec.md).
- Under `--json`, a side-effecting command **MUST** emit a **result receipt** — a
  JSON document describing what it did, carrying its `nextActions` in-band —
  rather than its human prose. This is what makes `--json` meaningful for a
  command whose human output is a confirmation rather than a result.
- The revision **MUST** keep `--json` (not `--format json`) as the single spelled
  form, and **MUST** keep `nextActions` appearing in-band under `--json`, per the
  existing convention.

### Opt-in capabilities

- The section **MUST** present the agent-facing enrichments — `--json`,
  `nextActions`, and the parked quiet/verbosity control — as capabilities a
  caller activates per invocation, distinct from the baseline.
- It **MUST** distinguish their availability: `--json` is a SHOULD across
  commands (near-universal, per the revised convention), whereas `nextActions`
  are offered only when a command has a useful follow-up.
- It **MUST** cross-reference the [Conventions](../../../specs/cli.md#conventions)
  for the detailed `--json` and `nextActions` rules rather than restating them,
  so there is a single source of truth.
- It **MUST** address the verbatim-artifact case explicitly: a command that omits
  `--json` (e.g. `spec`) still owes the full baseline — categorized exit codes
  and stderr diagnostics on failure.

### Housekeeping

- The change **MUST** remove *Agent-accessibility and CI requirements* and
  *Exit-code semantics* from the CLI spec's **To be specified** list once the
  section lands.
- The new and revised text **MUST NOT** contradict the existing `nextActions` or
  scope text; where it overlaps it **MUST** reconcile by reference.
- The change **MUST** update [`specs/cli/init.md`](../../../specs/cli/init.md):
  replace its "offers no `--json`" statement with the `init --json` receipt
  contract, reconciled with the `-` stdout-passthrough and the stderr
  confirmation.

## Requirements — the conformance

- The shipped `init` and `lint` commands **MUST** be audited against the baseline
  and the revised convention, and every divergence **MUST** be brought into
  compliance under this change.
- `qualitymd` **MUST** exit with the category code matching the outcome, replacing
  the current blanket `os.Exit(1)`:
  - a `lint` run that *finds findings* **MUST** exit with the
    *ran-but-found-problems* code, distinct from the *usage error* and
    *internal error* codes;
  - a malformed invocation **MUST** exit with the *usage error* code;
  - an I/O or unexpected failure **MUST** exit with the *internal error* code;
  - a clean run **MUST** exit `0`.
- `init` **MUST** gain a `--json` mode that emits a result receipt — at minimum
  the written path, whether the file was created, and its `nextActions` — to
  stdout, leaving the human confirmation (stderr) unchanged when `--json` is not
  passed.
- `init --json` **MUST** reconcile with the existing modes: the `-` argument
  still writes the raw scaffold to stdout, and `--json` reports the action rather
  than the scaffold contents; the spec/design **MUST** define the behavior when
  both are combined.
- The exit-code mapping **MUST** be implemented within the idiomatic Cobra / Fang
  stack per the [CLI spec's technical requirements](../../../specs/cli.md#technical-requirements);
  if Fang cannot carry distinct codes idiomatically, the
  [design doc](design.md) **MUST** record how the categories are surfaced.
- The conformance work **MUST** include tests pinning each exit-code category for
  the affected commands and the `init --json` receipt shape, and **MUST NOT**
  change any command's feature behavior or output beyond what the baseline and
  the revised convention require.

## Done criteria

- [`specs/cli.md`](../../../specs/cli.md) carries an **Agent accessibility**
  section with the baseline invariants, the concrete exit-code categories, the
  gated opt-in capabilities, and the revised SHOULD-by-default `--json`
  convention.
- The *Agent-accessibility and CI requirements* and *Exit-code semantics* lines
  no longer appear on that spec's **To be specified** list.
- [`specs/cli/init.md`](../../../specs/cli/init.md) documents the `init --json`
  receipt instead of "offers no `--json`".
- `qualitymd init` and `qualitymd lint` return the specified category exit codes,
  and `qualitymd init --json` emits the receipt — both covered by tests.
- No command feature behavior, README text, or other command sub-spec changes
  beyond the baseline conformance and `init`'s receipt mode.
- The change is moved through the lifecycle and archived according to the
  [changes process](../../index.md#status-lifecycle).
