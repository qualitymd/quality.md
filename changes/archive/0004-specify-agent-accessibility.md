---
type: Change Case
title: Specify and enforce agent accessibility
description: Add the cross-cutting agent-accessibility contract to the qualitymd CLI spec, broaden --json to a near-universal SHOULD, and bring the existing commands into conformance.
status: Done
tags: [cli, specification, agent-accessibility]
timestamp: 2026-06-17T00:00:00Z
---

# Specify and enforce agent accessibility

The [`qualitymd` CLI spec](../../specs/cli.md) names *agent accessibility* in its
header as a contract it will define, but the requirement is still a placeholder
on that spec's **To be specified** list. This change settles it, sharpens the
`--json` convention so machine-readable output is the broad default rather than
the exception, **and** brings the shipped commands into conformance — so the
contract is real rather than aspirational. The detail lives in its children:

- [Functional spec](0004-specify-agent-accessibility/spec.md) - the durable
  contract and the conformance delta the commands must meet.
- [Design doc](0004-specify-agent-accessibility/design.md) - how the baseline
  (notably categorized exit codes) is threaded through the Cobra / Fang stack.
  Added when the change advances to **Design**.

## Motivation

The bet behind `qualitymd` is that the CLI is the deterministic, mechanical
surface that skills and CI drive — *CLI is deterministic/mechanical; skills judge
through it*. That only holds if every command is reliably consumable by a
non-interactive caller. Today the CLI spec has agent-facing **conventions**
(`--json` output and `nextActions`) but no cross-cutting statement of what agent
accessibility *is*, and two concrete gaps:

- **Exit codes collapse.** Every error path in `internal/cli` ends in
  `os.Exit(1)`, so a caller cannot tell a `lint` that *found problems* from a
  usage error or an internal failure.
- **`--json` reads as off-by-default.** The convention offers `--json` only when
  a tight three-part test holds, so an agent must reason per command about
  whether the flag exists — and a wrong guess errors. Surveying agent-first CLIs
  (e.g. the Basecamp CLI, where nearly every command accepts `--json`) suggests
  the better default is *broad availability*: machine-readable output everywhere
  it is meaningful, so an agent can always reach for it.

Specifying the contract once, broadening `--json`, then closing the gaps gives
agents and CI a stable, uniform surface to branch on.

## Scope

Covered, in two parts:

1. **Specify.** Add an **Agent accessibility** section to
   [`specs/cli.md`](../../specs/cli.md) framing the contract as two tiers — a
   **baseline** binding every in-scope command (non-interactivity,
   stdout-is-payload / stderr-is-everything-else, deterministic output,
   categorized exit codes, plain output when stdout is not a terminal) and the
   **opt-in capabilities** a caller activates per invocation (`--json`,
   `nextActions`, and the parked quiet/verbosity control). As part of this:
   - settle the **exit-code categories** (success, ran-but-found-problems, usage
     error, internal error) concretely enough to implement; and
   - **revise the `--json` convention** in [Conventions](../../specs/cli.md#conventions)
     from a narrow gate to a SHOULD: commands SHOULD offer `--json`, the default
     rendering stays human (never auto-JSON), a command MAY omit it only when its
     output is a verbatim artifact that *is* the payload, and under `--json` a
     side-effecting command emits a **result receipt** (what it did plus
     `nextActions`) rather than its human prose.

   This removes *Agent-accessibility and CI requirements* and *Exit-code
   semantics* from the **To be specified** list.

2. **Enforce.** Audit the shipped commands (`init`, `lint`) against the baseline
   and the revised convention, and bring them into compliance, with tests:
   - **exit codes** — map the four categories to distinct, documented codes,
     replacing the blanket `os.Exit(1)`; and
   - **`init --json`** — add a result-receipt JSON mode to `init`, which under
     the old gate offered no `--json`. `lint` already satisfies the convention;
     `spec` is unimplemented and stays the deliberate verbatim-artifact carve-out
     when it lands.

Deferred: the global flag set and the quiet/verbosity control (still their own
**To be specified** items); the shared invocation form and file/stdin argument
convention; CI-specific requirements beyond what the baseline implies; and
implementing the [`spec`](../../specs/cli/spec.md) command — a separate change that
inherits this baseline (and the carve-out) by reference. No command's *feature*
behavior changes here beyond adding `init`'s receipt mode.

## Affected specs & docs

Updated before this change reaches **Done**:

- [x] [`specs/cli.md`](../../specs/cli.md) - add the **Agent accessibility**
      section (baseline, opt-in capabilities, concrete exit-code categories),
      revise the `--json` convention to a SHOULD with the verbatim-artifact
      carve-out and receipt framing, and remove the *Agent-accessibility and CI
      requirements* and *Exit-code semantics* entries from the **To be
      specified** list.
- [x] [`specs/cli/init.md`](../../specs/cli/init.md) - replace the "offers no
      `--json`" statement with the `init --json` receipt contract, and reconcile
      it with the `-` stdout-passthrough and stderr-confirmation behavior.

The conformance work in `internal/cli` (exit-code mapping and `init --json`,
plus tests) is this change's implementation, tracked by the functional spec and
design doc rather than listed here. No README status changes — no command becomes
newly available.

## Status

`Done`. Implemented and archived after adding the durable agent-accessibility
contract, documenting the exit-code categories and broadened `--json`
convention, adding `init --json`, mapping exit categories through the Fang
boundary, and pinning the behavior with tests.
