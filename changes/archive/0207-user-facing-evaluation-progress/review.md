---
type: Review
title: User-facing evaluation progress — review ledger
description: Requirement, durable/runtime contract, wording-sweep, and repository-gate evidence.
tags: [evaluation, skill, agent-mediated-ux, progress, review]
timestamp: 2026-07-15T00:00:00Z
---

# User-facing evaluation progress — review ledger

This ledger closes Change Case 0207 with durable-contract, runtime-guidance,
wording-sweep, and complete local-gate evidence.

## Requirement status

| Requirement                                             | Status | Evidence                                                                                                                                                                                                                                                                                                                                                                                                                          |
| ------------------------------------------------------- | ------ | --------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------- |
| R1 — Protect the implementation boundary                | Passed | `docs/guides/agent-mediated-ux.md`, `specs/skills/quality-skill/quality-skill.md`, and `skills/quality/SKILL.md` now require ordinary output to present task state, scope or coverage, attention, artifacts or results, and next action. They keep planning, protocol requests, payload/schema mechanics, workers/subagents, concurrency, and resume loops internal unless a decision or recovery action needs a specific detail. |
| R2 — Evaluate progress uses quality-task phases         | Passed | The durable evaluation wrapper and evaluate workflow, plus `skills/quality/workflows/evaluate.md`, use preflight, evidence review, report generation, and completion as healthy-run phases and require an attention-needed statement when unclear. The runtime includes a complete `Ready to evaluate` example.                                                                                                                   |
| R3 — Progress counts describe meaningful coverage       | Passed | Durable and runtime evaluate guidance permits model areas or requirements and prohibits work-unit, outstanding-request, request-window, concurrency, payload, and worker counts as ordinary progress. A repository sweep finds no remaining live skill requirement to narrate an outstanding cap or completed work-unit totals. The direct CLI stderr contract remains deliberately unchanged.                                    |
| R4 — Evaluator alternatives do not create false choices | Passed | Shared and workflow-specific durable specs and both runtime files now make default `harness` selection informational, allow independent-evaluator paths for future invocations, and require any offered current-run change to render a real closed choice and wait before mutation. Provider-named ambiguity still uses the existing wait-for-answer gate. Searches find no live “request it now/current run” default wording.    |
| R5 — Evaluator selection precedes the first write       | Passed | The durable flowchart and procedure now order lint → selection/optional dry-run → pre-mutation beat → feedback-log write → runner. The runtime workflow follows the same order, and the feedback-log section records that selection plus the beat precede log creation. The beat names evaluation artifacts and the local log accurately.                                                                                         |

## Durable-artifact rollup

- The shared `/quality` interaction contract carries the implementation
  boundary and its enduring rationale.
- The durable evaluation wrapper and evaluate workflow carry the phase,
  coverage, evaluator-choice, and mutation-order behavior.
- The agent-mediated UX guide carries the reusable protocol-to-task translation
  with good and avoid examples.
- Runtime `SKILL.md` and `workflows/evaluate.md` carry the complete operating
  procedure plus the user-facing presentation boundary.
- `CHANGELOG.md` records the user-facing skill and documentation outcomes. No
  CLI, runner, QUALITY.md format, project model, scaffold, installer,
  dependency, or generated-report behavior changed.

## Local gate

- `mise run check` passed warning-free typecheck, lint, 19 test files and 86
  tests, TypeScript and Markdown formatting, schema and report-gallery drift,
  specification and CLI documentation drift, npm package checks, and Mintlify
  link checks.
- `git diff --check` passed.
- Targeted searches found no live skill guidance requiring outstanding-cap
  progress, completed work-unit narration, a default evaluator “request it now”
  invitation, or feedback-log-before-selection flow.

All R1–R5 requirements passed. The change is ready for `Done` and archival.
