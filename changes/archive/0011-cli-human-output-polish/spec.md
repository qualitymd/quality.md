---
type: Functional Specification
title: CLI human output polish - functional spec
description: Remaining human-output requirements for styled model listings, lint next actions, dev version output, and gate coverage.
tags: [cli, dx, output]
timestamp: 2026-06-18T00:00:00Z
---

# CLI human output polish - functional spec

This spec states the delta for [CLI human output polish](../0011-cli-human-output-polish.md).
It builds on the durable [CLI spec](../../../specs/cli.md), especially the
agent-accessibility baseline, `--json` convention, next-action convention, human
output styling, and binary version rules.

The key words **MUST**, **MUST NOT**, **SHOULD**, and **MAY** are to be
interpreted as described in IETF RFC 2119.

## Requirements

### Styled models list

`qualitymd models list` **MUST** preserve the current plain table when stdout is
not a terminal or `NO_COLOR` is set. When stdout is a terminal and styling is
enabled, it **SHOULD** render the same table facts with the shared CLI palette:
visually distinct headers and accented model names.

`qualitymd models list --json` **MUST NOT** change.

### Lint next actions

When `qualitymd lint` emits an invalid result, its JSON document **MUST** include
deterministic `nextActions`. If any remaining finding is fixable, the preferred
action **SHOULD** be a `qualitymd lint --fix <path>` command; otherwise it
**SHOULD** be a rerun command for the same path.

Human lint failure output **MUST** render the same next-action command as a
footer on stderr, leaving stdout as the lint finding payload.

### Dev version output

An unstamped local build with a VCS revision available **MUST NOT** report a bare
`dev` version. It **MUST** include the short revision in the reported version or
commit form. If the Go toolchain's embedded build information omits the revision
for a local development invocation, the CLI **SHOULD** resolve the local VCS
revision best-effort before falling back to bare `dev`.

### Output gates

The implementation **MUST** keep the existing non-interactive baseline: no
prompts, spinners, fuzzy pickers, or blocking TTY-only flows. Tests **MUST** cover
the styling gates and prove JSON output remains free of terminal escape
sequences.
