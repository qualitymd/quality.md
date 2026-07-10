---
type: Change Case
title: Deterministic evaluation runner
description: Move `/quality evaluate` from skill-orchestrated judgment to a CLI-owned deterministic evaluation runner with pluggable evaluators.
status: Done
---

# Deterministic evaluation runner

## Motivation

The current evaluation workflow relies on the `/quality` skill as the primary
orchestrator. The CLI owns mechanical artifact writes, but the skill owns the
work graph, source routing, evaluator fanout, QC loop, and final synthesis. That
keeps evaluation portable across agent harnesses only as long as each harness can
faithfully execute the same prompt protocol.

`qualitymd` should instead own the deterministic evaluation runner: workflow
state, dependencies, validation, persistence, logging, and report generation.
Codex, Claude, direct provider APIs, shell checks, and human-mediated work should
all be evaluator transports for bounded judgment work units, not separate
evaluation engines. The runner should also own the choice of execution strategy:
serial calls, concurrent workers, native subagents, and provider context or
prompt-cache reuse are optimizations under the runner's contract, not alternate
orchestrators.

## Scope

This case specifies and designs the long-term evaluation runner architecture:

- `qualitymd evaluation run`;
- `.quality/config.yaml` evaluator selection;
- runner-owned execution strategy selection for sequential work, parallel work,
  subagent fanout, and evaluator context reuse;
- pluggable evaluator profiles for CLI, API, shell, and manual runtimes;
- a single authoritative `evaluation.json` run artifact;
- run-local execution and evaluator-call logs;
- `/quality evaluate` as an agent-mediated wrapper around the CLI runner.

Implementation is complete: the first slice ships `auto`/`sequential`
execution, CLI-backed `codex`/`claude` and API-backed `openai`/`anthropic`
evaluators, with `shell`/`manual` reserved.

## Affected artifacts

- **Code:** new `internal/runner/` (work graph, engine, run store,
  `evaluation.json` artifact, logs, dry run) and `internal/evaluator/`
  (contract, selection, `codex`/`claude` subprocess adapters,
  `openai`/`anthropic` API adapters) packages; `internal/cli/`
  (`evaluation_run.go`, `evaluation.go` dispatch); `internal/evaluation/`
  (`runner_support.go` report bridge and run preparation, artifact-aware
  `verifyRun`, manifest and status dispatch); `internal/workspace/` (the
  `evaluation` and `evaluators` config surface).
- **Durable specs:** added `specs/cli/evaluation-run.md`,
  `specs/evaluation/runner.md`, `specs/evaluation/evaluator-contract.md`, and
  `specs/evaluation/evaluation-json.md`; modified `specs/cli.md`,
  `specs/evaluation/{evaluation,orchestration,protocol}.md`,
  `specs/evaluation/records/{data-layout,payload-kinds}.md`,
  `specs/evaluation/reports/report-tree.md`, and
  `specs/skills/quality-skill/{quality-skill,evaluation,reporting}.md` with
  `workflows/evaluate.md` and `workflows/evaluate/feedback-log.md`, plus the
  affected bundle indexes and logs.
- **Bundled skill:** `skills/quality/SKILL.md`,
  `skills/quality/workflows/evaluate.md`, and
  `skills/quality/resources/{cli-workflow-conventions,output-policy}.md` now
  wrap `qualitymd evaluation run` instead of orchestrating evaluation.
- **Docs:** `README.md`, `install.md`, `docs/guides/reporting-design.md`,
  `mintlify/skill.mdx`, and bundle logs. `CONTRIBUTING.md`, `docs/reference/`,
  and the remaining Mintlify pages needed no changes.
- **Generated artifacts:** `mintlify/cli.mdx` regenerated for the new command.
  Report gallery outputs are unchanged: historical multi-file runs keep the
  existing report path.

## Children

- [Functional spec](0192-deterministic-evaluation-runner/spec.md)
- [Design doc](0192-deterministic-evaluation-runner/design.md)
