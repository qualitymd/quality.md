---
type: Review
title: Provider-affine SDK evaluator selection — review ledger
description: Requirement, durable/runtime contract, stale-policy sweep, evaluator-test, and repository-gate evidence.
tags: [evaluation, evaluator, selection, skill, sdk, harness, review]
timestamp: 2026-07-15T00:00:00Z
---

# Provider-affine SDK evaluator selection — review ledger

This ledger closes implementation review for Change Case 0208 with
durable-contract, runtime-guidance, stale-policy-sweep, evaluator-test, and
complete local-gate evidence.

## Requirement status

| Requirement                                               | Status | Evidence                                                                                                                                                                                                                                                                                                                                                                                     |
| --------------------------------------------------------- | ------ | -------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------- |
| R1 — Automatic SDK discovery precedes harness fallback    | Passed | `specs/skills/quality-skill/evaluation.md`, the durable evaluate workflow, `skills/quality/SKILL.md`, and the runtime evaluate workflow all place CLI automatic discovery ahead of `harness`. They allow automatic harness selection only when the preview reports `missing_evaluator` and the invoking harness can service checkpoints; other preview failures remain stops.                |
| R2 — Automatic selection prefers the invoking provider    | Passed | Durable and runtime guidance consume the structured automatic candidate receipt, select the usable candidate matching the known Codex or Claude harness, preserve the CLI-selected usable evaluator when no provider match is usable, and pass the resulting concrete evaluator explicitly. The design and runtime example make the provider-affinity reason observable.                     |
| R3 — Provider-named requests resolve to the SDK evaluator | Passed | Both durable selection contracts and both runtime files map bare Codex and Claude evaluator requests directly to `codex` and `claude`, reserve current-session judgment for explicit `harness` or the no-SDK fallback, and prohibit the former same-provider harness-versus-SDK question. A targeted active-surface search finds no remaining provider-name ambiguity or single-select rule. |
| R4 — Selection is visible before mutation                 | Passed | The evaluation wrapper, evaluate workflow, parent skill contract, and runtime guidance require the pre-mutation beat to name the evaluator and whether explicit intent, configuration, provider affinity, CLI discovery, or no-SDK fallback determined it. They distinguish fresh SDK sessions from current-session judgment and prohibit a current-run transport invitation.                |
| R5 — Selected runs remain evaluator-pinned                | Passed | Durable and runtime contracts bound automatic fallback to pre-run determination, require the real invocation to pass the concrete selected evaluator, keep unavailable explicit and configured selections outside the fallback chain, and prohibit switching providers or transports after run creation. Existing resume guidance continues to require a new run for a different evaluator.  |

## Durable-artifact rollup

- The parent `/quality` spec now summarizes SDK discovery before harness,
  provider affinity, determined-selection visibility, and checkpoint servicing
  only for harness-backed runs.
- The durable evaluation wrapper and evaluate workflow carry the complete
  precedence, candidate-consumption, provider-name, fallback, explanation, and
  pinning contract with enduring rationale.
- Runtime `SKILL.md` and `workflows/evaluate.md` carry the executable operating
  procedure and a provider-affine pre-mutation example.
- Bundle logs record the durable and runtime revisions. `CHANGELOG.md` records
  the user-facing selection outcome and unchanged CLI/format compatibility.
- Standalone CLI selection code, CLI/evaluator specs, tests, receipts,
  evaluation artifacts, the QUALITY.md specification, project model, docs,
  scaffold, installer, dependencies, and generated reports remain unchanged by
  design.

## Local gate

- `mise run test -- test/application/evaluator-selection.test.ts` passed all 10
  evaluator-selection tests, proving the unchanged automatic candidate receipt
  and deterministic standalone winner remain available to the skill.
- `mise run check` passed warning-free typecheck, lint, 19 test files and 86
  tests, TypeScript and Markdown formatting, schema and report-gallery drift,
  specification and CLI documentation drift, npm package checks, and Mintlify
  link checks.
- `git diff --check` passed.
- Targeted active-surface searches found no remaining harness-first precedence,
  provider-name ambiguity, same-provider single-select interaction, or default
  in-session selection rule.

All R1-R5 requirements passed. The change is ready for `Done` and archival.
