---
type: Design Doc
title: CLI human output polish - design doc
description: Implement the remaining CLI human-output polish through the existing Cobra, Fang, Lip Gloss, and Glamour stack.
tags: [cli, dx, output, design]
timestamp: 2026-06-18T00:00:00Z
---

# CLI human output polish - design doc

This design answers the [functional spec](spec.md) with the existing CLI stack.

## Approach

- Keep `tabwriter` as the `models list` alignment engine and add only a styled
  branch behind `colorEnabled(stdout)`. This avoids a table abstraction for a
  single table surface.
- Derive lint next actions in the lint result construction path so JSON and human
  output share one source. Human rendering writes the footer to stderr after the
  stdout findings.
- Normalize the dev-version fallback before handing values to Fang: when the
  embedded build info has a VCS revision but no module version, return a version
  string that already includes the short revision. If a contributor path such as
  `go run` omits embedded VCS data, use a best-effort `git rev-parse HEAD`
  fallback before accepting bare `dev`.
- Add focused tests around facts and semantic markers, not exact ANSI snapshots.

## Constraints

The plain and JSON paths remain the contract. Terminal styling stays a
convenience layer over the canonical output and is disabled for non-terminal
writers and `NO_COLOR`.
