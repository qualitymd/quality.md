---
type: Design Doc
title: Harness-native evaluator dispatch design
description: Checkpointed runner and skill design for using the invoking agent harness as an evaluator without introducing a second orchestrator.
tags: [evaluation, evaluator, agents, automation]
timestamp: 2026-07-10T00:00:00Z
---

# Harness-native evaluator dispatch design

## Context

This design answers the
[harness-native evaluator dispatch functional spec](spec.md). The current
runner synchronously calls a Go `Evaluator` implementation. That works for CLI
and API transports, but an agent skill cannot safely call back into its own
parent harness through that interface. Installed-command `auto` discovery
therefore chose Codex during a run launched by Claude Code.

The solution must let the current agent supply judgment while keeping every
orchestration decision and authoritative artifact inside the runner. It must
also work in unattended Claude Code routines and Codex scheduled tasks, where a
second CLI login or provider key may be unavailable or undesirable.

### Research basis

The product-facing parts of the design use current first-party documentation:

- [Claude Code routines](https://code.claude.com/docs/en/routines) run as
  autonomous cloud sessions with repositories, environments, connectors, and
  schedule/API/GitHub triggers. The
  [cloud environment reference](https://code.claude.com/docs/en/claude-code-on-the-web)
  documents fresh clones, committed project skills, environment variables,
  setup scripts, network policy, and `CLAUDE_CODE_REMOTE`; the
  [environment-variable reference](https://code.claude.com/docs/en/env-vars)
  documents `CLAUDECODE` in spawned shells. Current nested-session behavior is
  also visible in the
  [Anthropic Claude Code tracker](https://github.com/anthropics/claude-code/issues/29029),
  so the design does not depend on launching `claude -p` from an active Claude
  session.
- [Codex scheduled tasks](https://learn.chatgpt.com/docs/automations?surface=app)
  distinguish local projects/worktrees from web tasks and recommend explicit
  skill invocation for repeatable workflows. The
  [Codex environment-variable reference](https://learn.chatgpt.com/docs/config-file/environment-variables)
  documents `CODEX_API_KEY` and `CODEX_ACCESS_TOKEN` but no public parent-agent
  marker; the
  [cloud environment reference](https://learn.chatgpt.com/docs/environments/cloud-environment)
  distinguishes task-wide environment variables from setup-only secrets.
- The [Claude non-interactive reference](https://code.claude.com/docs/en/headless)
  and [Codex non-interactive reference](https://learn.chatgpt.com/docs/non-interactive-mode)
  document native JSON Schema output, no-persistence/ephemeral execution, and
  automation credential behavior used by fallback adapters.

These are documentation inputs, not evaluator correctness dependencies. The
core protocol relies only on the typed runner request/result contract.

## Approach

### A checkpointed evaluator transport

Add `harness` as a reserved evaluator whose `Evaluate` operation is split
across CLI invocations. The engine runs normally until it reaches a harness
work unit, then:

1. builds the ordinary typed `WorkRequest`;
2. derives a request ID and input hash;
3. persists a `PendingEvaluatorCall` in `evaluation.json` and sets the run to
   `awaiting_evaluator`;
4. atomically writes the artifact; and
5. returns an awaiting receipt carrying the complete bounded request.

The command exits zero. No evaluator failure occurred; control intentionally
returned to the caller.

Evaluation status, list, and workspace status read the pending-call metadata to
render an explicit awaiting-harness state and continuation command. They keep
the run in the incomplete count, but do not convert the checkpoint into a
failure or a generic missing-data gap.

The pending state stores correlation metadata, not raw prompt, source, or result
bodies:

```go
type PendingEvaluatorCall struct {
    RequestID     string
    WorkUnitID    string
    InputHash     string
    CorrelationID string
    Attempt       int
}
```

On `--resume` without a result, the engine rebuilds the request, compares its
hash with `PendingEvaluatorCall.InputHash`, and re-emits it. A mismatch is
`run_state_invalid`; this avoids persisting a second source snapshot while
preventing an old judgment from being attached to changed input.

### Result submission and continuation

`--evaluator-result <path|->` is valid only with `--resume` for a run awaiting a
harness result. The input is a transport envelope around the same payload the
Go `Evaluator` interface returns:

```go
type HarnessResultEnvelope struct {
    RequestID string
    InputHash string
    Evaluator HarnessIdentity
    Payload   map[string]any
    Failure   FailureCategory
    Detail    string
    Usage     *Usage
}
```

The runner verifies request correlation before passing the payload through the
existing normalization, schema validation, retry, accepted-result merge, log,
and atomic persistence paths. A normal submission call continues through any
deterministic units and stops only at the next harness checkpoint or terminal
receipt. This keeps the skill loop small and prevents deterministic work from
leaking into the agent interface.

Invalid envelopes do not advance the graph. Schema-invalid payloads use the
existing retry budget: the runner persists the attempt, emits a classified
receipt, and returns the next attempt's request. Once the retry budget is
exhausted, the run fails under the existing evaluator failure rules.

### Harness identity binding

The manifest continues to record `evaluator: harness` and `evaluatorKind:
harness`. The first accepted result binds `state.harnessIdentity.runtime` to a
stable surface identifier such as `claude-code` or `codex`. Later submissions
must match that runtime. Per-call model identifiers remain optional log metadata
because some harnesses do not expose a stable model string and a host may update
models without changing the evaluator transport.

Binding after the first accepted result avoids adding a provider-specific flag
to `evaluation run`, while still preventing a resumed run from mixing accepted
Claude and Codex judgments. A run with no accepted judgment may be resumed by
another harness because no evaluator attribution has yet been established.

### Skill-owned transport loop, runner-owned evaluation

The evaluate workflow changes from one long runner call to a receipt loop:

```text
resolve explicit/configured evaluator
  -> otherwise select harness
  -> run or resume qualitymd evaluation run
  -> terminal receipt? summarize and stop
  -> awaiting_evaluator? judge only the supplied request
  -> submit HarnessResultEnvelope
  -> repeat
```

The skill may use the harness's ordinary reasoning or a native subagent to
answer one request, but it receives no authority to schedule a different unit,
read broader source, persist results, or alter accepted output. The runner's
request is the entire evaluation boundary for that turn.

The workflow records evaluator-selection or checkpoint friction in its existing
feedback log. In unattended automation it does not add interactive gates: a
requested evaluation either advances, returns a report, or stops with the
runner's classified remedy.

### Selection and readiness

The skill resolves evaluator intent before invoking the runner:

1. explicit user evaluator request;
2. non-`auto` workspace `evaluation.evaluator`;
3. `harness` when the current agent can service checkpoints;
4. CLI `auto` only when no harness transport is available.

The standalone CLI never auto-discovers `harness`; it cannot portably know a
parent agent exists. Its existing `auto` order remains deterministic, but each
CLI candidate gains a readiness probe for executable presence, authentication,
and required structured-output flags. Dry run exposes candidate readiness and
the chosen reason without credential values.

CLI evaluators use their native schema flags (`--json-schema` for Claude Code,
`--output-schema` for Codex) and their no-persistence controls where supported.
Codex receives a temporary schema file removed after the call; both adapters
still return through the runner's independent schema validator. Adapter version
or capability mismatch is detected before the first judgment call.

### Automation documentation

Replace the two existing Mintlify placeholders with product-specific,
agent-first guides:

- **Claude Code:** distinguish cloud routines, Desktop scheduled tasks, and
  session-local loops; make the project skill and compatible CLI available in
  the cloud environment; invoke `/quality evaluate`; explain fresh-clone and
  branch/PR persistence; and reserve API/OAuth credentials for an explicitly
  selected non-harness fallback.
- **Codex:** create the scheduled task from ChatGPT web or the desktop app,
  select the local project or worktree, invoke `$quality` explicitly, keep
  workspace-write access for evaluation artifacts, and explain that local tasks
  require the machine and app to remain available.

Both guides say that the harness evaluator uses the automation's current agent
authentication and does not need a provider key. A separate fallback section
maps direct OpenAI/Anthropic API profiles and non-interactive Codex/Claude CLI
auth to their documented environment-variable names, with secrets supplied by
the host's secret facility and never committed to `.quality/config.yaml`.

## Spec response

- **Harness evaluator and ownership:** the checkpoint adapter preserves the
  existing work request/result contract and routes all accepted output back
  through the runner.
- **Checkpoint, resume, and stale input:** atomic pending metadata plus hash
  reconstruction makes every emitted request recoverable without persisting raw
  source twice.
- **Inspection and recovery:** existing status/list surfaces expose the awaiting
  lifecycle and exact continuation action from the same pending metadata.
- **Validation and attribution:** correlation checks precede the existing
  validation path; first-result runtime binding prevents mixed-harness runs.
- **Skill routing:** the skill uses information it actually has and names
  `harness` explicitly, while explicit and configured choices retain priority.
- **Fallbacks:** direct CLI discovery proves readiness rather than equating
  installation with usability, and native schema controls reduce malformed
  output without weakening runner validation.
- **Automation:** the same receipt loop runs interactively or unattended, and
  the durable guides make artifact and credential boundaries explicit.

## Alternatives

### Prefer the CLI matching the parent harness

Rejected. Claude Code can be detected in spawned shells, but launching nested
Claude Code is not a safe general evaluator transport. Codex has no equivalent
documented parent marker. Matching subprocess names would still create a second
agent session, authentication context, and usage stream rather than use the
agent already running the skill.

### Reverse or configure the `auto` command order

Rejected as the primary fix. It can choose a preferred fallback for a bare CLI
run, but it cannot make a provider-neutral default correct for both Claude and
Codex harnesses. Workspace pinning remains supported for teams that deliberately
want one external evaluator.

### Require direct provider API profiles for automation

Rejected. It would make API keys mandatory for subscription-backed agent users,
duplicate the host automation's authentication and billing path, and fit poorly
with cloud environments that intentionally restrict secret exposure during the
agent phase.

### Keep one long-running JSONL or JSON-RPC subprocess

Rejected for the first slice. Some harnesses expose durable PTY sessions while
others expose only discrete command calls. Receipt checkpoints compose with
both and reuse the runner's existing resume/persistence model. A streaming
transport may be added later as a performance optimization under the same
envelopes.

### Implement the reserved `manual` evaluator instead

Rejected. Human-mediated task export and agent-harness dispatch share a
checkpoint shape but have different defaults, attribution, latency, and UX.
Keeping `harness` explicit prevents a normal agent run from being described as
manual evaluation and leaves the human workflow free to earn its own contract.

## Trade-offs and risks

- **More CLI round trips.** Each evaluator work unit requires a result-submission
  command. 0193 reduces the number and size of those units; future streaming or
  batching can optimize transport without changing the contract.
- **Long agent turns.** Full evaluations may span many checkpoints. The workflow
  needs factual progress and durable resume rather than depending on one
  uninterrupted context window.
- **Source drift between checkpoints.** Rebuilding and hash-checking the pending
  request refuses stale results, but a changed source can force a new run. This
  is safer than silently accepting judgment against different evidence.
- **Harness identity granularity.** Runtime names are stable enough to prevent
  cross-provider mixing, while exact model strings may not be. Keeping model
  attribution per call avoids making unavailable metadata a resume dependency.
- **Docs track preview products.** Claude routines and Codex scheduled-task
  surfaces can change. The guides should link to current upstream docs and keep
  product details out of the core evaluator correctness contract.
- **Large awaiting receipts.** A bounded source package can still be large.
  Output remains machine-oriented JSON and must follow existing source limits;
  future file-backed or streaming request transport can optimize size without
  changing request identity.

## Open questions

None blocking. The status name, submission flag, first-result identity binding,
and sequential checkpoint shape are settled for the first implementation slice.
