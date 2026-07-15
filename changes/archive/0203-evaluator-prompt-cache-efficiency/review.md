---
type: Review
title: Evaluator prompt cache efficiency — review ledger
description: Requirement, isolation, provider usage, repeated live evaluation, and local-gate evidence.
tags: [evaluation, evaluator, prompt-caching, tokens, codex, claude, review]
timestamp: 2026-07-15T00:00:00Z
---

# Evaluator prompt cache efficiency — review ledger

This ledger closes Change Case 0203 with implementation and live provider
evidence. Live run folders were local verification artifacts only and were
removed after their non-sensitive call metadata was recorded here.

## Requirement status

| Requirement                            | Status | Evidence                                                                                                                                                                                                                                                                                                                                                                                              |
| -------------------------------------- | ------ | ----------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------- |
| R1 — Stable shared prompt prefix       | Passed | `EvaluationPromptParts` exposes a cacheable prefix and work-unit suffix joined by one stable boundary. Policy, kind/task, body guidance, shared context, and inspection policy precede work-unit identity, subject, and context. Focused tests prove byte equality across differing run/work-unit identity and byte divergence when shared inputs change.                                             |
| R2 — Canonical structured blocks       | Passed | Shared, inspection, and work-unit structured blocks use the existing recursively key-sorted and safely escaped `canonicalJson`. A focused test builds equivalent objects in different insertion orders and proves both prompt parts are byte-identical.                                                                                                                                               |
| R3 — Cache usage preservation          | Passed | Provider-neutral usage now distinguishes `cachedInputTokens` from `cacheWriteInputTokens`, retains reported zero, and leaves absent values absent. Codex maps every exposed usage value; Claude maps cache read and creation. The integration test proves all usage fields enter `logs/evaluator-calls.jsonl` and neither cache field enters `evaluation.json`; live Claude logs carried both fields. |
| R4 — Cache-stable Claude system prefix | Passed | Claude uses its supported `claude_code` preset with `excludeDynamicSections: true` beside the unchanged empty setting sources, non-persistence, neutral temporary directory, tool restrictions, environment allowlist, and cancellation. A unit test protects the exact option and two complete live runs exercised it.                                                                               |
| R5 — Independent-session semantics     | Passed | Codex still starts a new thread and Claude still starts a fresh non-persisted query for each request. No session ID, resume, continue, fork, earlier transcript, sibling output, or cache outcome enters another request, acceptance, or artifact. The durable evaluator specs now prohibit resume/fork token optimization and state that cached tokens still occupy the logical context window.      |
| R6 — Verification and live comparison  | Passed | Focused adapter/domain tests, the provider log integration test, the full repository gate, and two complete identical scoped Claude runs passed. The durable-spec rollup updates evaluator-contract, agent-evaluators, runner, and the evaluation spec log. The live comparison below records scope, runtime model handling, concurrency, per-call usage, duration, and attribution limits.           |

## Repeated live Claude comparison

Both accepted runs used `area:format-spec` scoped to
`factor:format-spec::clarity`, the same repository state, source
`./SPECIFICATION.md`, built-in `claude` evaluator, Claude Code 2.1.205, its
unchanged runtime-default model (the SDK did not surface an identifier), and
resolved concurrency 1. Each completed 12 work units, including six evaluator
calls. Run B started immediately after run A, inside the provider cache
lifetime.

| Run | Calls |  Duration | Input | Cache read | Cache creation | Output |    Cost |
| --- | ----: | --------: | ----: | ---------: | -------------: | -----: | ------: |
| A   |     6 | 428,820ms |    20 |     92,780 |        102,798 | 39,606 | $2.0830 |
| B   |     6 | 380,078ms |    18 |     83,529 |         67,101 | 33,022 | $1.5571 |

| Work-unit kind         | A input/read/write | A duration | B input/read/write | B duration |
| ---------------------- | -----------------: | ---------: | -----------------: | ---------: |
| Requirement judgment   |    8/92,780/40,648 |  281,385ms |    6/52,510/35,376 |  234,124ms |
| Factor analysis        |         2/0/11,759 |   29,551ms |      2/8,187/3,256 |   31,490ms |
| Area analysis          |         2/0/14,455 |   29,564ms |      2/8,101/5,536 |   24,991ms |
| Finding ranking        |          2/0/5,058 |   16,187ms |      2/2,718/2,965 |   21,492ms |
| Recommendation         |         4/0/22,918 |   51,725ms |     4/9,116/14,625 |   44,652ms |
| Recommendation ranking |          2/0/7,960 |   20,408ms |      2/2,897/5,343 |   23,329ms |

The strongest cross-run signal is the five corresponding synthesis calls. Run
A reported zero cache reads for all five; run B reported 31,019 cache-read input
tokens across them. Their cache creation fell from 62,150 to 31,725 tokens. The
inspection call reported substantial cache reads in both runs because its
multi-turn agent loop also reuses prior-turn prefixes.

Aggregate duration fell 11.4%, cache creation fell 34.7%, and reported cost
fell 25.2%. These are directional observations, not isolated causal estimates:
SDK totals combine within-session and cross-session caching, while output
tokens, tool choices, provider load, and model judgment varied. The evidence
therefore supports cache reuse and telemetry correctness, not a guaranteed
latency or cost percentage.

A preliminary CLI-maintainability scope reached the runtime but exhausted the
existing eight-turn inspection bound and was excluded from the accepted pair.
It exposed no cache-configuration or isolation failure.

## Local gate and rollup

- `mise run check` passed warning-free lint, typecheck, 18 test files and 65
  tests, schema/spec/CLI-doc/report-gallery drift checks, npm bundle checks,
  Mintlify links, and Markdown formatting.
- `git diff --check` passed.
- `specs/evaluation/evaluator-contract.md`, `agent-evaluators.md`, `runner.md`,
  and `log.md` carry the complete current contract. No requirement was renamed,
  added, or deleted during implementation, and the R1–R6 set-level review found
  no ambiguity, duplication, or missing dependency.
- No README, Mintlify, format specification, project model, dependency,
  generated artifact, or `/quality` workflow behavior changed.

All R1–R6 requirements passed. The change is ready for `Done` and archival.
