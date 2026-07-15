---
type: Review
title: Intent-faithful evaluator selection — review ledger
description: Requirement, discovery, receipt, authentication, skill, and local-gate evidence.
tags: [evaluation, evaluator, selection, discovery, skill, review]
timestamp: 2026-07-15T00:00:00Z
---

# Intent-faithful evaluator selection — review ledger

This ledger closes implementation for Change Case 0206 with application,
service, receipt, skill-contract, integration, and repository-gate evidence.

## Requirement status

| Requirement                                        | Status | Evidence                                                                                                                                                                                                                                                                                                                                                                            |
| -------------------------------------------------- | ------ | ----------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------- |
| R1 — Probe-all candidate reporting                 | Passed | `selectEvaluator` builds both built-in candidates before selecting the first usable one and attaches the full ordered array to every `auto` result. Dry-run receipts already project that array; `reportEvaluatorSelection` now adds it to returned run receipts. The selection test and fake-runtime CLI integration test both prove that Claude is reported when Codex is usable. |
| R2 — Ordering-decided selection reason             | Passed | When both candidates are usable, `selectEvaluator` names Codex as selected, states that deterministic discovery order decided, and names Claude as usable but not selected. Unit, run-receipt, and CLI dry-run assertions prove the winner, ordering basis, and usable-but-unselected candidate remain visible.                                                                     |
| R3 — Authentication basis is structured            | Passed | Every `EvaluatorCandidate` has `authenticationBasis: verified \| assumed \| unchecked`, separate from its human-readable evidence and with no credential field. Tests prove verified Codex, assumed Claude, and verified Claude receipt shapes; the CLI integration proves both verified values survive JSON serialization.                                                         |
| R4 — Documented Claude authentication probe        | Passed | `HostRuntimeLive` invokes `claude auth status --json`, parses only the boolean `loggedIn` field, returns `null` for an unavailable/unparseable probe, and never includes command output in candidate data. Selection rejects verified `false`, accepts verified `true`, and marks `null` assumed. The probe was also verified for login-backed and API-key-only authentication.     |
| R5 — Provider-named intent is disambiguated        | Passed | The durable evaluation and evaluate-workflow specs and their bundled runtime files require a single-select choice when a provider name could mean the same-provider current harness or SDK evaluator. They explain in-session versus fresh independent execution, put the SDK path first when independence is implied, and name explicit and configured paths for both.             |
| R6 — Default-harness explanation names alternative | Passed | The durable skill specs and bundled runtime require the pre-mutation explanation to identify current-session judgment and authentication, name the fresh independent SDK alternative, and state how to request it for the run or set `evaluation.evaluator` durably.                                                                                                                |
| R7 — Deterministic coverage                        | Passed | `test/application/evaluator-selection.test.ts` covers two candidates after a first-ready Codex probe, ordering reason, verified and assumed bases, Claude `true`/`false`/`null` handling, and run-receipt decoration. `test/integration/cli.test.ts` exercises the live host-runtime probe path with deterministic fake runtimes.                                                   |

## Durable-artifact rollup

- `specs/cli/evaluation-run.md` owns probe-all receipt reporting, structured
  authentication basis, and ordering-decided reasons.
- `specs/evaluation/agent-evaluators.md` owns the documented non-interactive
  authentication-probe policy and credential boundary.
- `specs/skills/quality-skill/evaluation.md` and
  `specs/skills/quality-skill/workflows/evaluate.md` own transport
  disambiguation and default-harness explanation; the bundled runtime mirrors
  those contracts.
- `CHANGELOG.md` records the user-visible CLI receipt and skill behavior. The
  evaluator contract, format specification, project model, schema, adapters,
  CLI grammar, dependencies, installer, scaffold, and README remain unchanged.

## Local gate

- `mise run check` passed warning-free typecheck, lint, 19 test files and 86
  tests, TypeScript and Markdown formatting, schema and generated-output drift,
  npm package checks, and Mintlify link checks.
- Focused evaluator-selection, CLI, provider-execution, and deterministic-byte
  integration tests passed before the full gate.
- `git diff --check` passed before review.

All R1–R7 requirements passed. The implementation is ready for review.
