---
type: Design Doc
title: Evaluation run folder naming - design doc
description: How Evaluation v2 run creation and discovery move to NNNN-full-eval and NNNN-<scope-path>-eval while recognizing legacy run folders.
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
`NNNN-full-eval` / `NNNN-<scope-path>-eval` grammar, while existing
`-quality-eval` folders remain valid historical runs that are listed, inspected,
and counted when allocating the next number.

## Approach

Keep run naming inside `internal/evaluation`, where the run-folder concept
already lives. Replace the single legacy-only matcher with a small parser over
two explicit grammars:

- current: `NNNN-<scope>-eval`;
- legacy: `NNNN[-<legacy-scope>]-quality-eval`, including the older
  `subject`/`model` prefixed forms.

The parser returns the run number and the normalized narrowing value. For current
runs, scope `full` normalizes to no narrowing; any other scope is the narrowing
slug. For legacy runs, bare `quality-eval`, `subject-quality-eval`, and
`model-quality-eval` normalize to no narrowing, while `subject-<slug>`,
`model-<slug>`, and unprefixed legacy narrowed names normalize to the slug.

`CreateRun` keeps validating `--narrowing` as a path-safe slug, then builds the
new name by choosing `full` when no narrowing is supplied and appending the
`eval` tag. The CLI does not derive or validate structural model paths; that
remains the `/quality` workflow's responsibility because the narrowing slug is a
mnemonic and the model snapshot remains the structural source of truth.

## Spec Response

- The new name builder satisfies the run-folder grammar requirements for
  full-scope and narrowed runs.
- The parser keeps recognition centralized for next-number computation, run
  listing, and latest-run resolution, so both current and legacy folders are
  counted and discoverable through the same path.
- `narrowingFromRunName` delegates to the parser, so list output treats `full` as
  no narrowing and strips legacy `subject`/`model` prefixes.
- The bundled `/quality` evaluation workflow carries the full structural scope
  path into `--narrowing`; the CLI continues to treat the slug as opaque after
  path-safe validation.

## Alternatives

- **Extend the old regular expression and keep submatch indexing.** Rejected.
  A single expression that accepts both grammars makes the capture groups
  fragile and easy to misuse. Two named matchers plus parser logic keep the
  compatibility cases readable.
- **Reject `--narrowing full`.** Rejected by the functional spec. The run number
  remains the identity, and adding a special usage error for a vanishingly rare
  cosmetic collision is not worth the extra branch.
- **Rename existing run folders.** Rejected. Existing runs are historical
  artifacts; compatibility belongs in recognition, not migration.

## Trade-offs & Risks

- The CLI cannot prove that a caller's narrowing slug really came from a full
  Area/Factor structural path. That is intentional: the skill owns semantic
  scoping, while the CLI owns safe naming and deterministic run allocation.
- Current and legacy grammars share one number sequence. That preserves
  monotonicity but means parser mistakes affect both creation and listing, so the
  tests need explicit mixed-format coverage.

## Open Questions

None.
