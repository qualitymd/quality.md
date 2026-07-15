---
type: Review
title: Transport-aware evaluator concurrency — review ledger
description: Requirement, scheduling, persistence, harness delegation, schema, and local-gate evidence.
tags: [evaluation, evaluator, concurrency, agents, runner, review]
timestamp: 2026-07-15T00:00:00Z
---

# Transport-aware evaluator concurrency — review ledger

This archived ledger closes Change Case 0204 with deterministic implementation,
contract, workflow, and local-gate evidence. No live provider run is required:
the change governs runner scheduling and transport declarations, and its
acceptance properties are observable through controlled evaluator services and
persisted artifacts.

## Requirement status

| Requirement                                                 | Status | Evidence                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                            |
| ----------------------------------------------------------- | ------ | --------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------- |
| R1 — Runner scheduling authority                            | Passed | The runner alone creates and advances the graph, request identities, retries, validation, evidence sealing, persistence, ordering, and reports. Direct workers receive one reconstructed `EvaluationRequest`, return one completion through a queue, and never receive an artifact or scheduler callback. The harness workflow now gives a delegated worker exactly one receipt request and explicitly withholds the full frontier, artifact writes, quality control, and recursive delegation.                                                     |
| R2 — Honest transport capabilities                          | Passed | `EvaluatorDispatchCapability` separately declares `concurrentCalls`, `delegatedRequests`, positive `automaticConcurrency`, and optional `maxConcurrency`. Built-ins declare harness `(false, true, 4)`, Codex `(true, false, 4)`, and Claude `(false, false, 1, max 1)`. Selection tests prove a configured Claude profile inherits the Claude declaration. Validation rejects malformed maximums and sequential declarations.                                                                                                                      |
| R3 — Transport-aware resolution                             | Passed | The pure resolution matrix proves automatic and configured harness/Codex behavior, Claude clamping and sequential resolution, source provenance, and rejection of zero, negative, fractional, string, and null configured values. `HostRuntime.hardwareConcurrency` and every evaluation use of CPU count are removed. A CLI integration test proves `--concurrency` remains unrecognized.                                                                                                                                                          |
| R4 — Selected-evaluator-first creation                      | Passed | `executeEvaluationRun` receives the already selected evaluator, validates and resolves its dispatch declaration before staging requests, and writes the evaluator, kind, capability record, and cap together in the initial schema-9 manifest. The controlled direct-pool test reads the artifact before any completion and proves the selected provider identity, exact capability declaration, cap two, `running` state, and exactly two pending requests. No post-creation evaluator patch path remains.                                         |
| R5 — Completion-driven direct dispatch                      | Passed | A deferred evaluator test proves cap two starts exactly two calls; completing the second call persists it and starts the third while the first remains active. The observed peak stays two, the accepted mid-run unit is durable before top-up, every call is a fresh evaluator invocation, and deliberately out-of-order completion still projects payloads in graph order. A sequential provider test proves concurrency-one completion and report parity.                                                                                        |
| R6 — Bounded harness delegation                             | Passed | Existing CLI coverage proves partial submission preserves unsubmitted requests without retry cost and tops the bounded window up. Artifact and receipt tests prove delegated status, cap semantics, and peak outstanding diagnostics. Skill and workflow contracts now constrain native workers to one request and use “outstanding of up to N” wording without claiming active subagents.                                                                                                                                                          |
| R7 — Per-result durability, retry, resume, and cancellation | Passed | Shared resume acceptance serializes each result, atomically writes `evaluation.json`, and only then lets the direct coordinator free the slot. A retry test proves only the failed requirement consumes a second attempt. An interruption test proves the accepted sibling remains completed, the active worker is interrupted, the run becomes `cancelled` with one pending request, and resume completes without re-evaluating the accepted requirement.                                                                                          |
| R8 — Artifact, observability, and verification              | Passed | Dry-run JSON proves harness automatic concurrency four and the structured capability declaration; receipts and artifacts retain resolved `concurrency` as a cap. Initial events record source, requested, automatic, optional maximum, resolved, clamped, and dispatch mode; activity events record direct peak active or delegated peak outstanding. Artifacts and report rebuild use schema 9, and CLI coverage proves an in-flight schema-8 run is refused with a start-new-run remedy. The complete acceptance matrix and repository gate pass. |

## Durable-artifact rollup

- `specs/evaluation/runner.md`, `orchestration.md`, `evaluator-contract.md`,
  `agent-evaluators.md`, and `evaluation-json.md` carry the cumulative runner,
  transport, persistence, and schema contract with the 0204 rationale.
- `specs/cli/evaluation-run.md` retains workspace configuration as the only
  concurrency override and describes cap semantics in dry-run and receipts.
- `specs/skills/quality-skill/evaluation.md`, its evaluate workflow spec, and
  the runtime `/quality` skill constrain one-request native delegation and
  honest progress wording. Their evaluation, workflow, and skill logs record
  the revision.
- `CHANGELOG.md` records the selected-transport defaults, completion-driven
  pool, schema-9 clean break, and skill behavior. No README, Mintlify, format
  specification, project model, dependency, installer, scaffold, JSON Schema,
  or report-gallery contract changed.

## Local gate

- `mise run check` passed warning-free typecheck, lint, 19 test files and 80
  tests, TypeScript and Markdown formatting, schema drift, report-gallery drift,
  specification and CLI documentation drift, npm package checks, and Mintlify
  link checks.
- `git diff --check` passed.
- Repository searches found no live `hardwareConcurrency`, legacy evaluator
  `concurrent`/`subagents` capability fields, version-8 current artifact writer,
  or concurrency flag. The remaining schema-8 references are the intentional
  refusal fixture and historical records.

All R1–R8 requirements passed. The change is ready for `Done` and archival.
