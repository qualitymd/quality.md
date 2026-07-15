---
type: Design Doc
title: Evaluator prompt cache efficiency — design
description: Explicit prompt parts, canonical structured blocks, provider usage adapters, and cache-stable Claude session configuration.
tags: [evaluation, evaluator, prompt-caching, tokens, codex, claude]
timestamp: 2026-07-15T00:00:00Z
---

# Evaluator prompt cache efficiency — design

Design behind the
[Evaluator prompt cache efficiency](../0203-evaluator-prompt-cache-efficiency.md)
Change Case and its [functional spec](spec.md).

## Context

The runner already constructs one immutable `EvaluationRequest` per work unit,
and both SDK adapters call the same pure prompt renderer before opening a fresh
session. The optimization therefore belongs at that renderer/adapter boundary:
make repeated request bytes identical as far into the provider input as their
semantics allow, then preserve provider usage as diagnostics. The runner's work
graph, acceptance, persistence, and report paths do not need a caching concept.

The output schema is supplied through each SDK's structured-output option rather
than in the user prompt. It is already stable for a work-unit kind, as are the
tool set and sandbox options. The project-owned prompt must align with that
provider prefix instead of introducing work-unit identity before its own shared
material.

## Approach

### Explicit prompt parts

`src/domain/evaluator/context.ts` introduces a readonly
`EvaluationPromptParts` value:

```ts
interface EvaluationPromptParts {
  readonly cacheablePrefix: string
  readonly workUnitSuffix: string
}
```

`renderEvaluationPromptParts(request)` derives both strings. The existing
`renderEvaluationPrompt(request)` remains the adapter-facing entry point and
joins them with one exported, stable `Work-unit delta:` boundary. Tests compare
the prefix directly; runtime logs never store either raw part.

The prefix order is:

1. evaluator isolation and JSON-only output policy;
2. work-unit kind and its runner-owned instructions;
3. applicable QUALITY.md body guidance;
4. shared accepted model/area context; and
5. inspection availability, boundary, and read/search policy.

The suffix carries work-unit ID, subject, and work-unit context. `runId` remains
transport identity only and is not rendered, matching current behavior.

This makes every requirement in one area share the longest available prefix:
the same task, body, area frame, source selector, workspace root, and policy.
Requirements in different areas still share the earlier task/body portion until
their shared area context diverges. Synthesis kinds form their own stable
prefixes with tools unavailable.

### Canonical structured rendering

Every structured block uses the existing `canonicalJson` domain function. It
recursively sorts object keys and applies the repository's safe JSON escaping,
while compact rendering removes formatting-only tokens. The request hash remains
unchanged because it already hashes structured values canonically; this change
only aligns provider-facing bytes with that determinism.

### Provider usage adapters

`EvaluationUsage` gains optional `cacheWriteInputTokens`. Small pure functions
in the adapter translate the two pinned SDK usage shapes:

- Codex: input, output, and cached input tokens exposed by the SDK;
- Claude: input, output, cache-read input, cache-creation input, and total cost.

The functions copy reported numeric values, including zero, and do not fabricate
fields absent from an SDK. `evaluation-provider` already forwards the returned
usage object, and `evaluation-resume` already writes it to the evaluator-call
log, so no application orchestration change is needed.

The integration test gives the mock provider all fields, then proves the call
log contains them and `evaluation.json` does not. This keeps provider cache state
diagnostic and non-authoritative.

### Claude system-prompt shaping

The Claude query options add:

```ts
systemPrompt: {
  type: "preset",
  preset: "claude_code",
  excludeDynamicSections: true,
}
```

The preset preserves Claude Code's provider-owned agent/tool instructions.
`excludeDynamicSections` moves the random neutral working directory and other
session-specific runtime material after the SDK's global cache boundary through
its supported path. All existing isolation options stay adjacent and unchanged:
empty setting sources, no persistence, read/glob/grep only when inspecting,
explicitly disallowed mutation/shell/agent tools, sanitized environment, and
cancellation.

A readonly exported constant supplies this option so a unit test can protect it
without mocking or running the SDK.

### Live evidence

After focused and full gates pass, run the same small factor scope twice through
Claude within the provider cache lifetime. Keep the model, workspace, source,
prompt inputs, and concurrency configuration fixed. Extract only non-sensitive
call metadata from each run's `logs/evaluator-calls.jsonl`: work-unit kind,
duration, and usage counts. Record the comparison and attribution limits in
`review.md`; do not copy prompts, tool transcripts, file bodies, or model output.

If Claude is unavailable at execution time, the case cannot satisfy R4/R6 and
does not advance to `In-Review`. A zero cache read is not a blocker by itself,
but it triggers a source/options review before acceptance.

## Spec response

- R1–R2 are answered by explicit prompt parts, a stable boundary, fixed field
  order, and canonical JSON, with pure byte-equality tests.
- R3 is answered by the expanded provider-neutral usage value, pure SDK mappers,
  unchanged forwarding, and the run-log integration assertion.
- R4 is answered by the cache-stable Claude preset constant beside the unchanged
  isolation options and by live Claude execution.
- R5 is answered by retaining `startThread`/fresh `query` calls and
  `persistSession: false`; no session ID becomes input to a later call.
- R6 is answered by focused tests, the existing full gate, and the repeated live
  comparison ledger.

## Alternatives

### Fork or resume one seeded session

Rejected. Codex's pinned public SDK has start/resume but no fork, while Claude's
fork copies a transcript into each branch. Either approach makes shared model
conversation part of judgment input, violates the fresh non-persisted session
contract, and leaves the copied transcript inside the logical context window.
Provider key/value caching gives the economic reuse without semantic coupling.

### Add direct OpenAI/Anthropic API evaluators for explicit cache controls

Rejected. The coding-agent SDK owns iterative inspection, context management,
and tool policy. Reintroducing direct APIs would reverse the 0201 boundary and
create a second agent loop merely to expose cache keys or breakpoints.

### Reorder the flat prompt without exposing parts

Rejected. It could improve today's prefix, but tests would need fragile substring
positions and a later edit could silently move a varying field earlier. An
explicit value makes the intended compatibility property directly testable.

### Warm one request before releasing each concurrency cohort

Deferred. A warm-up barrier trades latency for a possible cache hit and could be
counterproductive when provider/runtime system prefixes are already warm.
Prompt shaping and telemetry must land first; scheduling needs its own evidence
and contract if cold bursts remain material.

### Replace the Claude system prompt with a minimal custom prompt

Rejected. A custom prompt could be shorter but would discard provider-owned
coding-agent instructions and create a project-owned sidecar policy. The preset's
supported dynamic boundary preserves behavior while improving cacheability.

## Trade-offs and risks

- Canonical compact JSON is less pleasant to read in a raw prompt, but raw
  prompts are neither user output nor logged artifacts; exact and smaller bytes
  are more valuable here.
- Moving Claude dynamic system sections into its first user message makes the
  neutral working-directory context marginally less authoritative. That path is
  not evaluator policy, while the project-owned prompt repeats the actual
  inspection and instruction boundaries explicitly.
- Aggregate SDK usage may include within-session tool-loop cache reads as well
  as cross-session prefix reuse. Repeated runs can show reported reuse and
  directional change, not isolate causal savings.
- Prefixes still diverge by model body, workspace path, area frame, source, and
  schema/tool configuration. That is intentional: making unlike judgment input
  identical would be a correctness defect.

## Open questions

None. Cache-aware scheduling remains deferred rather than unresolved inside
this design.
