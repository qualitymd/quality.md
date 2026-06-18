---
type: Change Case
title: CLI human output polish
description: Finish the remaining human-output polish for models list, lint next actions, dev version reporting, and output-gate coverage.
status: Done
tags: [cli, dx, output]
timestamp: 2026-06-18T00:00:00Z
---

# CLI human output polish

This change follows up on the delightful CLI output work by closing the remaining
gaps where the current command surface does not yet use the existing stack and
output conventions fully.

- [Functional spec](0011-cli-human-output-polish/spec.md) - what the change must
  do.
- [Design doc](0011-cli-human-output-polish/design.md) - how it is implemented
  without weakening the agent-facing output contract.

## Motivation

`qualitymd` already has a shared Fang/Lip Gloss/Glamour presentation stack, but a
few human responses still fall short of the durable CLI conventions:
`models list` is an unstyled tabwriter table, invalid lint results carry no
useful next action, local dev version output can still be bare `dev`, and the
gate between terminal polish and plain/JSON output needs broader coverage.

## Scope

Covered:

- Styled terminal rendering for `qualitymd models list`.
- Deterministic lint next actions in JSON and a human stderr footer on lint
  failures.
- Dev version output that includes a short VCS revision when available.
- Focused tests for terminal styling gates, plain output, JSON cleanliness, and
  version fallback behavior.

Deferred:

- Interactive prompts, spinners, fuzzy pickers, or TUI flows.
- A general table-rendering abstraction before more table surfaces exist.
- Global quiet/verbosity flags.

## Affected specs & docs

- [x] [`specs/cli.md`](../../specs/cli.md) - clarify that `NO_COLOR` suppresses
      styling for stdout and stderr human surfaces.
- [x] [`specs/cli/lint.md`](../../specs/cli/lint.md) - require useful lint
      `nextActions` on invalid results and the matching human stderr footer.
- [x] [`specs/cli/models.md`](../../specs/cli/models.md) - specify terminal-only
      styling for `models list`.

## Status

`Done`. Implemented and archived after styling `models list`, adding shared lint
next actions for JSON and human output, making dev version output include a
short VCS revision when available, updating the durable CLI specs, and adding
focused output-gate tests.
