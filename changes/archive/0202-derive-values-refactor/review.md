---
type: Review
title: Derive-values conformance refactor — review ledger
description: Requirement, byte-preservation, enumeration, source-conformance, and local-gate evidence.
tags: [refactor, effect, typescript, evaluation, cli, review]
timestamp: 2026-07-14T00:00:00Z
---

# Derive-values conformance refactor — review ledger

This ledger closes Change Case 0202 with its acceptance evidence. The temporary
detached worktree used for the before/after artifact comparison was removed
after the hashes matched; the deterministic fixture remains an integration
test.

## Requirement status

| Requirement                            | Status | Evidence                                                                                                                                                                                                                                                                                                                                                                                |
| -------------------------------------- | ------ | --------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------- |
| R1 — Public behavior preservation      | Passed | The existing contract suite passes unweakened. A deterministic harness fixture run against the pre-refactor commit and final implementation produced identical SHA-256 hashes for stdout, stderr, `evaluation.json`, `events.jsonl`, `evaluator-calls.jsonl`, and `model-snapshot.md`; the fixture is retained in the integration suite.                                                |
| R2 — Derivation and Effect conformance | Passed | The `src/` census replaced cross-loop collection mutation with mapped, filtered, flattened, reduced, or Effect-traversed expressions. The lint visitor confines its emission buffer to one small collector helper. Newly introduced and reshaped Effect operations use named `Effect.fn`, typed result handling, and ordered Effect traversal; exported collection fields are readonly. |
| R3 — Evaluation-execute decomposition  | Passed | Frame construction, protocol request/receipt assembly, artifact/event/receipt assembly, and summary rendering are pure domain functions with focused tests. Execution uses the parsed document body and named Effect hashing adapters.                                                                                                                                                  |
| R4 — Ready-unit graph selection        | Passed | `readyUnits` is pure and its test verifies dependency eligibility, concurrency bounding, order, and non-mutation.                                                                                                                                                                                                                                                                       |
| R5 — Single run-folder enumeration     | Passed | Creation, execution, dry-run, status, inspection/listing, and evaluation-data latest selection consume one application scanner backed by one domain classifier.                                                                                                                                                                                                                         |
| R6 — Recognition and numbering rule    | Passed | Pure and filesystem tests cover current and historical manifests, precedence, invalid and unreadable documents, folder fallback, non-directories, positive-number validation, and the formerly excluded `quality` slug. The shared rule is durable in the CLI spec.                                                                                                                     |

## Local gate

- `mise run check` passed warning-free lint, typecheck, 18 test files and 60
  tests, schema/spec/CLI-doc/report-gallery drift checks, npm bundle checks,
  Mintlify links, and Markdown formatting.
- `git diff --check` passed.
- The focused deterministic artifact-byte test passed after being promoted to
  the permanent integration suite.

All R1–R6 requirements passed. The change is ready for `Done` and archival.
