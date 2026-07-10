---
type: Functional Specification
title: Harness-native evaluator dispatch
description: Requirements for using the invoking agent harness as an evaluation runner transport while preserving runner-owned orchestration and supporting unattended automations.
tags: [evaluation, evaluator, agents, automation]
timestamp: 2026-07-10T00:00:00Z
---

# Harness-native evaluator dispatch

This spec governs the delta from subprocess-first `auto` discovery to a
first-class harness evaluator for `qualitymd evaluation run`. It inherits the
binding work graph, result-validation, persistence, and report semantics from
the durable [evaluation runner](../../specs/evaluation/runner.md),
[evaluator contract](../../specs/evaluation/evaluator-contract.md), and
[orchestration contract](../../specs/evaluation/orchestration.md). It does not
change evaluation judgment semantics or result payload schemas.

The key words **MUST**, **MUST NOT**, **SHOULD**, and **MAY** are to be
interpreted as described in BCP 14 when, and only when, they appear in all
capitals.

## Background / motivation

`auto` currently prefers an installed Codex CLI over an installed Claude CLI.
That makes a run launched by Claude Code silently dispatch judgment to Codex
when both commands are present. Reversing the order would produce the mirror
failure, and environment-based parent detection is not portable across agent
harnesses. Subscription-backed routines and scheduled tasks should use the
authenticated agent already running the quality skill without needing a nested
agent process or an additional API credential.

The needed boundary is not a second orchestrator. The runner remains the sole
owner of evaluation state and emits one bounded evaluator work request; the
invoking harness supplies a typed result that enters the same validation and
persistence path as every other evaluator result.

## Scope

This change covers harness-backed evaluator selection, request/result
checkpointing, identity and attribution, resume and failure behavior, direct-CLI
fallback readiness, structured CLI fallback output, and automation guidance.

It defers parallel harness dispatch, native subagent scheduling, `shell` and
human-mediated `manual` evaluators, provider API modernization, and automation
trigger creation. Token and context reuse remain owned by 0193.

## Assumptions and dependencies

- The invoking harness can run the quality skill, execute successive CLI
  commands, consume JSON receipts, and produce a JSON object matching a supplied
  schema.
- The runner can rebuild a pending work request from its model snapshot,
  work-graph state, and current source package and compare its input hash with
  the persisted pending-call metadata.
- 0193 lands before implementation so this case builds on the combined
  requirement judgment unit and stable/delta request shape.
- Host automation capabilities and credential names come from current Claude
  Code and Codex documentation and may require future documentation updates;
  evaluator correctness does not depend on those product-specific details.

## Requirements

### Harness evaluator

`qualitymd evaluation run` **MUST** accept `harness` as a reserved built-in
evaluator that delegates only bounded evaluator work requests to the invoking
agent harness.

> Rationale: the current harness is a distinct transport, not an installed CLI
> candidate; naming it makes the intended routing explicit and attributable.
>
> Durable spec: modify `specs/cli/evaluation-run.md` and
> `specs/evaluation/evaluator-contract.md` to add the built-in harness
> evaluator and its boundary.

The harness evaluator **MUST NOT** own run creation, scope expansion, work-graph
ordering, retry policy, result validation, persistence, report generation, or
final authority outside the result envelope it submits.

> Rationale: making the outer agent the workflow orchestrator would recreate
> the harness-dependent evaluation behavior that 0192 removed.
>
> Durable spec: modify `specs/evaluation/evaluator-contract.md` and
> `specs/evaluation/runner.md` to preserve runner ownership across harness
> checkpoints.

### Checkpointed dispatch

When a harness-backed evaluator work unit becomes ready, the runner **MUST**
atomically persist an awaiting-evaluator checkpoint before returning a receipt
that contains the run reference, request identity, work-unit identity and kind,
subject, instructions, context, bounded source package, expected result schema,
input hash, and correlation ID.

> Rationale: the CLI command returns control to the invoking agent between
> request and result, so the run must be resumable before the work request
> leaves stdout.
>
> Durable spec: modify `specs/cli/evaluation-run.md`,
> `specs/evaluation/runner.md`, `specs/evaluation/evaluation-json.md`, and
> `specs/evaluation/orchestration.md` to define the checkpoint receipt and
> persisted pending-call metadata.

An awaiting-evaluator receipt **MUST** use the stable status
`awaiting_evaluator`, **MUST** exit successfully, and **MUST** distinguish the
normal checkpoint from `failed` and `cancelled` outcomes.

> Rationale: awaiting harness judgment is expected progress, not a failure that
> automation should retry from the beginning.
>
> Durable spec: modify `specs/cli/evaluation-run.md` and
> `specs/evaluation/evaluation-json.md` to add the lifecycle status and exit
> behavior.

`qualitymd evaluation status`, `qualitymd evaluation list`, and workspace
`qualitymd status` **MUST** identify an awaiting-evaluator run as resumable and
incomplete, report that harness judgment is the pending action, and distinguish
it from a failed, cancelled, malformed, or generically incomplete run.

> Rationale: unattended or interrupted workflows use status as the recovery
> surface; describing a normal checkpoint as an undifferentiated gap loses the
> exact action needed to continue.
>
> Durable spec: modify `specs/cli/evaluation-status.md`,
> `specs/cli/evaluation-list.md`, and `specs/cli/status.md` to surface the
> awaiting lifecycle and continuation action.

`qualitymd evaluation run --resume <run> --evaluator-result <path|->` **MUST**
accept exactly one harness result envelope from a file or stdin, and it **MUST**
advance deterministic work until the next evaluator checkpoint or the terminal
run receipt.

> Durable spec: modify `specs/cli/evaluation-run.md` and
> `specs/evaluation/runner.md` to define result submission and continuation.

If an awaiting run is resumed without `--evaluator-result`, the runner **MUST**
re-emit the same pending request when its rebuilt input hash matches the
checkpoint. If it cannot rebuild the same request, it **MUST** fail with
`run_state_invalid` and recommend a new run rather than accepting judgment for
changed input.

> Rationale: an interrupted agent must be able to recover the work request, but
> a result must never be attached to evidence other than the evidence the
> evaluator saw.
>
> Durable spec: modify `specs/evaluation/evaluation-json.md`,
> `specs/evaluation/runner.md`, and `specs/evaluation/orchestration.md` to define
> pending-request reconstruction and stale-input behavior.

### Result validation, retry, and identity

A submitted harness result envelope **MUST** identify the pending request and
input hash and carry either a result payload or a classified evaluator failure.
The runner **MUST** reject a mismatched, duplicate, or unsolicited result
without advancing the work graph.

> Durable spec: modify `specs/evaluation/evaluator-contract.md`,
> `specs/evaluation/runner.md`, and `specs/evaluation/orchestration.md` to define
> harness result correlation and rejection behavior.

The runner **MUST** normalize, schema-validate, retry, accept, log, and persist a
harness payload through the same paths used for CLI- and API-backed evaluator
payloads. Invalid output **MUST NOT** be repaired or persisted by the skill.

> Durable spec: modify `specs/evaluation/evaluator-contract.md` and
> `specs/evaluation/runner.md` to make transport-independent validation explicit.

The first accepted harness result **MUST** bind the run to a stable harness
runtime identity supplied in the result envelope; later results from a different
runtime identity **MUST** be refused with `run_state_invalid`. The runner
**SHOULD** record the model per call when the harness reports it without making
model metadata a correctness dependency.

> Rationale: a resumable run must not silently mix Claude, Codex, or another
> harness after accepted judgments already exist, while model metadata is not
> uniformly available across surfaces.
>
> Durable spec: modify `specs/evaluation/evaluator-contract.md`,
> `specs/evaluation/evaluation-json.md`, and `specs/evaluation/runner.md` to add
> harness attribution and resume compatibility.

Raw harness requests, source contents, result bodies, credentials, and tokens
**MUST NOT** be written to evaluator-call logs. Existing hashes, identities,
durations, attempt state, usage when available, and failure categories remain
the observable call metadata.

> Durable spec: modify `specs/evaluation/runner.md` to extend the existing
> evaluator-call logging boundary to harness dispatch.

### Skill routing and boundaries

When `/quality evaluate` runs in a harness capable of servicing harness
checkpoints, and neither the user nor `evaluation.evaluator` selects another
evaluator, the skill **MUST** invoke the runner with `--evaluator harness`
instead of relying on CLI `auto` discovery.

> Rationale: the skill knows the active harness; the standalone CLI does not
> have a portable, documented way to infer it.
>
> Durable spec: modify `specs/skills/quality-skill/evaluation.md` and
> `specs/skills/quality-skill/workflows/evaluate.md` to define the harness-native
> default and override precedence.

The skill **MUST** service each harness checkpoint only from the runner-supplied
request, submit the typed envelope through the CLI, and treat the runner's
accepted state and terminal receipt as authoritative. It **MUST NOT** construct
its own work graph, widen source, write evaluation records, or adjust accepted
results.

> Rationale: the harness provides judgment, not a second evaluation workflow.
>
> Durable spec: modify `specs/skills/quality-skill/quality-skill.md`,
> `specs/skills/quality-skill/evaluation.md`, and
> `specs/skills/quality-skill/workflows/evaluate.md` to replace the current
> blanket prohibition on skill judgment with this bounded evaluator role.

The skill **MUST** preserve explicit user evaluator choices and a non-`auto`
workspace evaluator configuration. It **MUST** explain the selected transport
before the first evaluation mutation and **MUST NOT** silently cross to a
different provider after harness selection or failure.

> Durable spec: modify `specs/skills/quality-skill/evaluation.md` and
> `specs/skills/quality-skill/workflows/evaluate.md` to define selection
> precedence and failure presentation.

### Direct CLI fallback selection

For direct CLI use, `auto` **MUST** consider a CLI-backed evaluator usable only
after verifying that its executable, authentication state, and required
non-interactive structured-output capabilities are available. It **MUST** skip
an unusable candidate and report the readiness evidence and final selection
reason in dry-run JSON.

> Rationale: command presence alone currently lets dry run describe an
> unauthenticated or incompatible evaluator as ready.
>
> Durable spec: modify `specs/cli/evaluation-run.md` and
> `specs/evaluation/evaluator-contract.md` to define readiness-aware discovery
> and preview evidence.

`auto` **MUST NOT** infer a parent agent harness from undocumented environment
variables. A harness-backed run is selected explicitly by the skill or caller.

> Rationale: Claude and Codex expose different subprocess environments, and an
> internal variable is not a cross-harness compatibility contract.
>
> Durable spec: modify `specs/cli/evaluation-run.md` to keep direct CLI
> discovery distinct from harness routing.

Where a supported CLI exposes JSON Schema output enforcement and ephemeral or
no-session-persistence controls, its built-in evaluator adapter **MUST** use
those controls for bounded work requests. The runner still **MUST** validate the
returned payload independently.

> Rationale: prompt-only JSON and retained one-off sessions add avoidable output
> failures and local state without replacing runner validation.
>
> Durable spec: modify `specs/evaluation/evaluator-contract.md` to require
> native structured output and non-persistence where supported.

### Automation guidance

The Claude Code and Codex automation guides **MUST** provide runnable,
agent-first paths for recurring `/quality evaluate` use that invoke the quality
skill explicitly, make the compatible `qualitymd` CLI and project skill
available, and state how evaluation artifacts persist beyond one run.

> Durable spec: none. This is durable user documentation tracked by the parent
> Change Case.

The automation guides **MUST** distinguish local and cloud execution,
permissions, network access, repository/worktree behavior, and unattended-run
limitations that affect evaluation. They **MUST** explain that harness-backed
evaluation uses the automation's authenticated agent and needs no provider API
key.

> Durable spec: none. This is durable user documentation tracked by the parent
> Change Case.

Fallback credential guidance **MUST** name only documented credential sources
for the selected evaluator kind, distinguish CLI subscription/access-token auth
from direct provider API keys, reference secrets by environment-variable name
only, and warn against exposing job-wide credentials to repository-controlled
code.

> Rationale: `OPENAI_API_KEY`, `CODEX_API_KEY`, `ANTHROPIC_API_KEY`, and Claude
> OAuth tokens serve different transports and are not interchangeable.
>
> Durable spec: modify `specs/evaluation/evaluator-contract.md` for evaluator
> credential semantics; the product-specific setup remains in durable user
> documentation.

## Example protocol

An initial harness-backed command reaches its first judgment checkpoint:

```text
qualitymd evaluation run --evaluator harness --json
```

```json
{
  "status": "awaiting_evaluator",
  "run": ".quality/evaluations/0002-full-eval",
  "evaluator": "harness",
  "evaluatorRequest": {
    "requestId": "req_...",
    "workUnitId": "assessRateRequirement:requirement:root::example",
    "kind": "assessRateRequirement",
    "subject": "requirement:root::example",
    "inputHash": "...",
    "correlationId": "...",
    "expectedSchema": {}
  }
}
```

The skill supplies the result and receives either the next checkpoint or the
terminal receipt:

```text
qualitymd evaluation run --resume 0002-full-eval \
  --evaluator-result - --json
```

```json
{
  "requestId": "req_...",
  "inputHash": "...",
  "evaluator": {
    "runtime": "claude-code",
    "model": "<reported model when available>"
  },
  "payload": {}
}
```

The complete work request and payload schemas remain owned by the evaluator
contract and existing payload-kind specs; this example fixes the routing and
correlation shape, not the domain result schema.

## Durable spec changes

### To add

None. The existing evaluator, runner, orchestration, artifact, CLI, and skill
specs own the complete harness transport contract; no new durable named artifact
is introduced.

### To modify

- `specs/cli/evaluation-run.md` - add the harness evaluator, checkpoint/result
  command surface, awaiting receipt and exit behavior, readiness-aware direct
  CLI discovery, and dry-run evidence (per the harness, checkpoint, and fallback
  requirements).
- `specs/cli/evaluation-status.md` - report awaiting harness judgment as a
  distinct resumable, incomplete state with the exact continuation action (per
  the status requirement).
- `specs/cli/evaluation-list.md` - preserve and filter the awaiting lifecycle in
  run listings rather than collapsing it into a generic incomplete entry (per
  the status requirement).
- `specs/cli/status.md` - include awaiting harness runs and their continuation
  action in workspace evaluation history and routing (per the status
  requirement).
- `specs/evaluation/evaluator-contract.md` - add harness work/result envelopes,
  runtime attribution, credential semantics, and hardened CLI fallback behavior
  (per the harness, validation, identity, fallback, and credential requirements).
- `specs/evaluation/runner.md` - define runner ownership across checkpoints,
  pending-call handling, transport-independent validation, identity enforcement,
  and logging boundaries (per the checkpoint, validation, identity, and logging
  requirements).
- `specs/evaluation/evaluation-json.md` - add `awaiting_evaluator`, pending-call
  metadata, harness identity state, and resume compatibility (per the checkpoint,
  lifecycle, reconstruction, and identity requirements).
- `specs/evaluation/orchestration.md` - add harness checkpoint scheduling,
  stale/mismatched result handling, and retry behavior (per the checkpoint,
  reconstruction, and result-correlation requirements).
- `specs/skills/quality-skill/quality-skill.md` - permit the skill's bounded
  harness-evaluator role while preserving runner ownership (per the skill
  boundary requirement).
- `specs/skills/quality-skill/evaluation.md` - make harness-native dispatch the
  default in capable agent surfaces and define the checkpoint loop and override
  precedence (per the skill routing requirements).
- `specs/skills/quality-skill/workflows/evaluate.md` - specify the runtime
  evaluate procedure, progress, failures, and unattended behavior for harness
  checkpoints (per the skill routing requirements).

### To rename

None.

### To delete

None.

## Open questions

None blocking. The first slice is deliberately sequential: one persisted
pending request, one submitted result, then deterministic advancement to the
next checkpoint.
