---
type: Change Case
title: User-facing evaluation progress
description: Keep evaluator protocol and worker orchestration behind the agent interface while presenting evaluation choices, phases, coverage, and recovery in user-facing language.
status: Done
tags: [evaluation, skill, agent-mediated-ux, progress]
timestamp: 2026-07-15T00:00:00Z
---

# User-facing evaluation progress

Status note: case is **Done** and archived. Runtime guidance, durable specs, UX
docs, release notes, R1-R5 review evidence, and the full repository gate are
complete.

## Motivation

A real `/quality evaluate` run exposed its implementation protocol as the user
experience: deterministic and evaluator "work units," outstanding request
windows, concurrency caps, payload schemas, resume loops, subagent fan-out, and
orchestration scaffolding. Those details were correct but forced the user to
reverse-engineer runner mechanics to understand the simple state: preflight had
passed and evidence review was underway.

The same run said it would use the current session "unless you'd prefer
otherwise" and then continued immediately. That made a preference invitation
look like a choice while giving the user no actual gate. The evaluator-selection
contract also required the independent alternative to be offered for the
current run, while workflow ordering opened the feedback log before claiming
that the first mutation was still ahead.

The agent is the interface. Ordinary progress should tell the user where the
evaluation is, what has been covered, whether anything needs attention, and
what comes next. Protocol details should surface only when one is necessary to
make a decision or recover a stopped run.

## Scope

Covered:

- add an explicit implementation-boundary rule to shared agent-mediated UX;
- translate evaluate progress into user-facing phases and useful coverage
  counts rather than runner request, payload, concurrency, or worker mechanics;
- keep evaluator selection informational when default precedence has already
  settled it, while preserving a real wait-for-answer choice for ambiguous
  provider-named intent or any offered current-run change;
- order evaluator selection before the first workflow write and make the
  pre-mutation progress beat accurate; and
- align durable skill specs, runtime skill guidance, logs, and release notes.

Non-goals:

- changing the runner work graph, evaluator protocol, concurrency, receipt
  schema, evaluator precedence, or recovery semantics;
- hiding diagnostic details that a user actually needs to choose a remedy or
  recover a run;
- changing the CLI's direct human stderr contract; and
- adding a new progress artifact, command, configuration key, or interaction
  affordance.

## Affected artifacts

Derived from a repository sweep for evaluator-selection explanations,
`awaiting_evaluator` progress, outstanding-request caps, pre-mutation ordering,
and the shared user-interaction contract.

- **Change record:** this parent, `spec.md`, `design.md`, and `review.md` under
  `changes/0207-user-facing-evaluation-progress/`; `changes/index.md` and
  `changes/log.md`.
- **Durable skill specs:**
  `specs/skills/quality-skill/quality-skill.md` adds the shared implementation
  boundary; `specs/skills/quality-skill/evaluation.md` aligns evaluator
  selection, workflow order, and harness progress;
  `specs/skills/quality-skill/workflows/evaluate.md` owns the workflow-specific
  phase and choice behavior. `specs/log.md` and
  `specs/skills/quality-skill/workflows/log.md` record the revisions.
- **Durable UX guide:** `docs/guides/agent-mediated-ux.md` adds protocol-to-UX
  translation guidance and a false-affordance example; `docs/log.md` records
  the revision.
- **Bundled skill runtime:** `skills/quality/SKILL.md` aligns the shared
  interaction and evaluator-selection rules;
  `skills/quality/workflows/evaluate.md` changes the selection/write order and
  presents phase-based progress while retaining the internal checkpoint loop;
  `skills/quality/log.md` and `skills/quality/workflows/log.md` record the
  revisions.
- **Release surfaces:** `CHANGELOG.md`, `package.json`,
  `npm/quality.md/package.json`, and `skills/quality/SKILL.md` release metadata
  advance for the patch release.
- **Code and tests:** no `src/` or `test/` behavior changes; the runner and CLI
  contracts remain unchanged. Verification is durable/runtime contract review,
  Markdown and package drift checks, and the full release gate.
- **Format specification, project model, generated reports, scaffold, install,
  and dependencies:** no `SPECIFICATION.md`, `QUALITY.md`, report-gallery,
  scaffold, installer, or dependency change.

## Children

- [Functional spec](0207-user-facing-evaluation-progress/spec.md) — the
  implementation-boundary, progress, choice, and sequencing requirements.
- [Design doc](0207-user-facing-evaluation-progress/design.md) — the shared
  implementation boundary, phase translation, evaluator-choice posture, and
  pre-mutation sequence.
- [Review ledger](0207-user-facing-evaluation-progress/review.md) — R1-R5
  durable/runtime contract, wording-sweep, and repository-gate evidence.

## Status

`Done`. Every R1-R5 requirement has direct implementation and verification
evidence in the
[review ledger](0207-user-facing-evaluation-progress/review.md). The case and
its child bundle are archived together.
