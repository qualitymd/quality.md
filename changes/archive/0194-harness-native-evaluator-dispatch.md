---
type: Change Case
title: Harness-native evaluator dispatch
description: Make /quality evaluation use the invoking agent harness for bounded judgment by default, with checkpointed runner dispatch, readiness-aware fallbacks, and unattended automation guidance.
status: Done
tags: [evaluation, evaluator, agents, automation]
timestamp: 2026-07-10T00:00:00Z
---

# Harness-native evaluator dispatch

## Motivation

The first full-model run through the deterministic evaluation runner was
invoked from Claude Code, but `auto` selected the Codex CLI because discovery
uses a fixed installed-command order. The outer Claude session therefore became
only a launcher for 111 separate Codex subprocess calls. The choice was
deterministic, but it was not native to the invoking harness and silently
crossed provider, authentication, usage, and model boundaries.

Changing the installed-command order is not enough. Claude Code exposes a
documented subprocess signal, but nested Claude Code sessions are not a safe
general transport; Codex does not document an equivalent parent-harness marker.
Cloud routines and scheduled tasks add a second constraint: they should be able
to use the authenticated agent already running the skill without provisioning a
second provider credential merely so `qualitymd` can call another agent.

The runner needs a first-class harness evaluator: the runner still owns the
work graph, bounded request, validation, retries, persistence, and reports, while
the invoking agent harness supplies only the judgment result for the current
work request.

## Scope

Covered:

- a built-in `harness` evaluator and a checkpointed request/result protocol for
  agent-mediated judgment without a nested evaluator subprocess;
- runner state, receipts, resume, identity binding, validation, and failure
  behavior for harness-backed work;
- `/quality evaluate` selecting the current harness by default when no explicit
  user or workspace evaluator overrides it;
- readiness-aware `auto` fallback discovery for direct CLI use, without trying
  to infer an invoking harness from undocumented environment variables;
- native structured-output and no-persistence controls for CLI-backed fallback
  evaluators where the installed CLI supports them; and
- user guidance for recurring evaluation in Claude Code routines and Codex
  scheduled tasks, including installation, artifact persistence, permissions,
  network access, and documented credential options for non-harness fallbacks.

Deferred:

- parallel or batched harness work requests and native subagent fanout;
- `shell` and human-mediated `manual` evaluator implementations;
- direct OpenAI adapter migration from Chat Completions to Responses and general
  provider-model default policy;
- automation creation or management by `qualitymd`; this case makes evaluation
  runnable inside host automations but does not create schedules or triggers;
- token and context reuse work owned by
  [0193](0193-evaluation-token-efficiency.md).

Non-goals:

- requiring API keys from Claude Code or Codex subscription users;
- changing evaluation semantics, the work graph, rating roll-up, accepted
  result payload kinds, or generated reports; or
- allowing the skill or harness evaluator to write run artifacts, expand scope,
  choose dependencies, or alter accepted runner results.

Implementation depends on 0193 landing first because both cases touch evaluator
request construction and runner execution. The functional contract and design
can settle while 0193 is in review; code remains unstarted.

## Affected artifacts

- **Code:** `internal/cli/evaluation_run.go` and tests (harness result input and
  checkpoint receipts); `internal/evaluator/` (reserved name, harness result
  envelope, readiness probes, and hardened CLI adapters); `internal/runner/`
  (pending evaluator-call state, checkpoint/resume, validation, logging, and
  receipts); `internal/evaluation/` and `internal/status/` (run inspection and
  history summaries for awaiting runs); `internal/cli/evaluation.go` and tests
  (evaluation status/list presentation); and `internal/workspace/workspace.go`
  and tests (accepted evaluator configuration).
- **Durable specs:** modify `specs/cli/evaluation-run.md`,
  `specs/cli/evaluation-status.md`, `specs/cli/evaluation-list.md`,
  `specs/cli/status.md`,
  `specs/evaluation/evaluator-contract.md`, `specs/evaluation/runner.md`,
  `specs/evaluation/evaluation-json.md`, `specs/evaluation/orchestration.md`,
  `specs/skills/quality-skill/quality-skill.md`,
  `specs/skills/quality-skill/evaluation.md`, and
  `specs/skills/quality-skill/workflows/evaluate.md`. See the functional spec's
  [durable spec changes](0194-harness-native-evaluator-dispatch/spec.md#durable-spec-changes).
- **Bundled skill:** update `skills/quality/SKILL.md`,
  `skills/quality/workflows/evaluate.md`,
  `skills/quality/resources/cli-workflow-conventions.md`, and
  `skills/quality/resources/output-policy.md` for harness selection, the
  checkpoint loop, and new receipts; update their local logs/indexes where
  required.
- **Durable docs:** replace the placeholders in `mintlify/claude-code.mdx` and
  `mintlify/codex.mdx` with runnable automation guidance; update `install.md`
  for evaluator defaults and fallback credentials; regenerate
  `mintlify/cli.mdx` and `mintlify/skill.mdx`, and update Mintlify navigation
  only if the page structure changes. `README.md` needs no change unless
  implementation introduces a new top-level invocation rather than keeping
  `/quality evaluate` unchanged.
- **Release record:** add the user-visible evaluator-routing and automation
  support to `CHANGELOG.md` when implementation begins.
- **Format specification and scaffold:** no impact. The QUALITY.md schema,
  model semantics, and starter content do not change.
- **Generated evaluation reports:** no content or layout change; existing
  report-gallery artifacts should remain byte-stable for unchanged evaluation
  data.

## Children

- [Functional spec](0194-harness-native-evaluator-dispatch/spec.md)
- [Design doc](0194-harness-native-evaluator-dispatch/design.md)

## Status

`Done`. Implemented and landed: the reserved `harness` evaluator with
checkpointed dispatch (`awaiting_evaluator` receipts, `--evaluator-result`
submission, deterministic pending-request reconstruction, first-result runtime
binding), awaiting-aware status/list/workspace-status surfaces,
readiness-aware `auto` discovery with dry-run candidate evidence,
capability-gated native CLI schema and no-persistence flags, the skill's
harness-native selection precedence and checkpoint loop, the Claude Code and
Codex automation guides, and all listed durable spec amendments. Verified with
the full check suite plus a live end-to-end checkpoint loop against the built
binary.
