---
type: Design Doc
title: Evaluation run folder naming - design doc
description: How Evaluation v2 run creation and discovery move to NNNN-full-eval and NNNN-<scope-path>-eval.
tags: [cli, evaluation, naming, design]
timestamp: 2026-06-26T00:00:00Z
---

# Evaluation run folder naming - design doc

Design behind the
[Evaluation run folder naming](../0113-evaluation-run-folder-naming.md) change
and its [functional spec](spec.md).

## Context

Evaluation v2 run names are currently produced in `internal/evaluation/create.go`
and recognized through one regular expression in `internal/evaluation/path.go`.
That recognition feeds next-run-number computation, run listing, latest-run
resolution, and the narrowing value shown by `qualitymd evaluation list`.

The change is deliberately small: new runs should use the shorter
`NNNN-full-eval` / `NNNN-<scope-path>-eval` grammar.

## Approach

Keep run naming inside `internal/evaluation`, where the run-folder concept
already lives. Replace the old matcher with a parser over the current
`NNNN-<scope>-eval` grammar.

The parser returns the run number and the normalized narrowing value. Scope
`full` normalizes to no narrowing; any other accepted scope is the narrowing
slug. The `quality` segment is reserved so old `-quality-eval` names do not
match the current grammar.

`CreateRun` keeps validating `--narrowing` as a path-safe slug, then builds the
new name by choosing `full` when no narrowing is supplied and appending the
`eval` tag. The CLI does not derive or validate structural model paths; that
remains the `/quality` workflow's responsibility because the narrowing slug is a
mnemonic and the model snapshot remains the structural source of truth.

## Spec Response

- The new name builder satisfies the run-folder grammar requirements for
  full-scope and narrowed runs.
- The parser keeps recognition centralized for next-number computation, run
  listing, and latest-run resolution.
- `narrowingFromRunName` delegates to the parser, so list output treats `full` as
  no narrowing.
- The bundled `/quality` evaluation workflow carries the full structural scope
  path into `--narrowing`; the CLI continues to treat the slug as opaque after
  path-safe validation.

## Alternatives

- **Extend the old regular expression and keep submatch indexing.** Rejected.
  The new grammar is clearer as its own parser.
- **Reject `--narrowing full`.** Rejected by the functional spec. The run number
  remains the identity, and adding a special usage error for a vanishingly rare
  cosmetic collision is not worth the extra branch.
- **Rename existing run folders.** Rejected. Existing runs are historical
  artifacts; this case does not rewrite them.

## Trade-offs & Risks

- The CLI cannot prove that a caller's narrowing slug really came from a full
  Area/Factor structural path. That is intentional: the skill owns semantic
  scoping, while the CLI owns safe naming and deterministic run allocation.
- Parser mistakes affect both creation and listing, so tests need explicit
  current-format coverage.

## Open Questions

None.
