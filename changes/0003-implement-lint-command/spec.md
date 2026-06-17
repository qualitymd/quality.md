---
type: Functional Specification
title: lint command implementation — functional spec
description: Implement qualitymd lint according to the durable lint sub-spec.
tags: [cli, command, lint]
timestamp: 2026-06-17T00:00:00Z
---

# lint command implementation — functional spec

Companion to the [Implement the lint command](../0003-implement-lint-command.md)
change. The durable [`qualitymd lint` sub-spec](../../specs/cli/lint.md) is the
complete functional contract for command-specific behavior; this change spec
states the implementation delta and does not restate the rule catalog or output
schema.

The key words **MUST**, **MUST NOT**, **SHOULD**, and **MAY** are to be
interpreted as described in IETF RFC 2119.

## Scope

Covered: add an executable `qualitymd lint` command that validates a
`QUALITY.md` file according to [`specs/cli/lint.md`](../../specs/cli/lint.md),
including the initial rule set, findings, locations, human output, `--json`
output, `--fix` repair behavior, deterministic ordering, and
blocked-downstream-rule behavior.

Deferred: every feature deferred by the durable lint sub-spec, including
suppression directives, rule selection, severity overrides, patch/full-file
repair output modes, and a lint-emitted rule catalog.

## Requirements

- The implementation **MUST** satisfy the command-specific behavior in the
  durable [`qualitymd lint` sub-spec](../../specs/cli/lint.md).
- `lint` **MUST** be wired into the `qualitymd` CLI as a first-class command
  using the established Cobra/Fang stack.
- `lint` **MUST** validate the same structural model that other local tooling
  uses to load `QUALITY.md`, so parser behavior, YAML shape checks, and rule
  behavior do not drift apart.
- `lint` **MUST** include focused tests for every initial error and warning rule
  listed in the durable lint sub-spec.
- `lint` **MUST** include tests for the `--json` document shape, finding
  location shape, valid-file output, deterministic ordering, and non-zero exit
  behavior on errors.
- `lint --fix` **MUST** apply fixable findings in place according to the durable
  lint sub-spec, then report the post-repair lint result.
- `lint --fix` **MUST** include tests for applied repairs, no-op repair runs,
  repair failure without partial writes, symlink refusal, Markdown body
  preservation, and JSON repair reporting.
- The implementation **MUST NOT** add lint-specific flags or behavior that the
  durable lint sub-spec defers.

## Cross-cutting dependency

The lint sub-spec is complete for lint-specific behavior. The parent
[`qualitymd` CLI spec](../../specs/cli.md) still has the shared invocation form
and file/stdin argument convention in its "To be specified" list. This change
settles only the minimum invocation shape needed to implement `lint` now:
`qualitymd lint [path]`, defaulting to `QUALITY.md`. Stdin handling and the
shared file argument convention remain parent-CLI work.

## Done criteria

- `qualitymd lint` validates valid and invalid `QUALITY.md` files according to
  the durable lint sub-spec.
- `qualitymd lint --fix` repairs fixable findings in place according to the
  durable lint sub-spec.
- Human-readable output and `--json` output are covered by tests.
- The README no longer describes `lint` as planned.
- The change is moved to **Done** and archived according to the
  [changes process](../index.md#status-lifecycle).
