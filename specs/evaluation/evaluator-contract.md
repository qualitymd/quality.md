---
type: Functional Specification
title: Evaluator contract
description: Capability, inspection request, result, harness, and agent-runtime profile contract for evaluation evaluators.
tags: [evaluation, evaluator, agents]
timestamp: 2026-07-15T00:00:00Z
---

# Evaluator contract

An evaluator is the runtime the [evaluation runner](runner.md) uses for bounded
evaluation judgment. This spec defines the contract every evaluator kind honors
and the workspace configuration that names evaluators.

The key words **MUST**, **MUST NOT**, **SHOULD**, and **MAY** are to be
interpreted as described in BCP 14 when, and only when, they appear in all
capitals.

## Background / motivation

Requirement evidence cannot be selected responsibly before the requirement is
interpreted. Codex and Claude coding-agent SDKs already provide iterative
search, inspection, context management, cancellation, and tool policy. The CLI
therefore consumes those agent loops instead of recreating one around a direct
model API. The runner remains the deterministic evaluation and artifact harness;
the evaluator owns requirement-specific discovery and judgment. — 0201

## Boundaries

The runner **MUST** own run state, model and scope validation, work-graph
construction, dependency ordering, scheduling, retry, cancellation, result and
evidence validation, persistence, artifact paths, report generation, and output
ordering.

An evaluator **MUST** own iterative requirement-specific context discovery,
inspection, evidence selection, assessment, rating, and synthesis judgment
inside the typed work request it receives. It **MUST NOT** expand scope, write
run artifacts, or become a second workflow engine.

## Capability declaration

Every evaluator **MUST** declare structured output, workspace inspection,
instruction isolation, mediated verification, network policy, tool use,
concurrent calls, nested subagents, fresh-session isolation, cancellation,
usage reporting, turn limits, token or cost limits, context-window visibility,
compaction control, sandbox control, and executable override. Unsupported
controls **MUST** be represented as unsupported rather than silently ignored.

Requirement inspection requires structured output, workspace inspection,
instruction isolation, a fresh session, cancellation, read-only workspace
policy, disabled network, and non-interactive approval policy. A runtime that
cannot establish these boundaries **MUST** fail readiness or dispatch as
`evaluator_incompatible`.

Verification is optional and separately declared. An evaluator without a
mediated verification path remains usable, but its request **MUST** say
verification is unavailable and judgment **MUST** record that limit when it
matters. Unmediated host shell access is never a fallback.

## Work-unit envelope

Every evaluator kind **MUST** consume the same envelope for a work-unit kind.
The envelope carries:

- run, work-unit, subject, and correlation identity;
- runner-owned instructions;
- shared accepted context and per-work-unit context;
- applicable QUALITY.md body guidance;
- the expected JSON Schema;
- a request input hash; and
- for requirement judgment, an inspection context containing the absolute
  authorized workspace root, effective source selector and form, and workspace,
  network, approval, verification, and repository-instruction policy.

The request **MUST NOT** contain a runner-selected source bundle. The effective
source identifies the judged subject, not all context the evaluator may inspect.
The evaluator **MAY** inspect supporting context elsewhere in the authorized
workspace but **MUST** classify it separately and **MUST NOT** widen the subject.

`assessRateRequirement` **MUST** use one fresh session and return one object
containing `assessment`, `rating`, and an evidence proposal. Other judgment
units receive accepted structured results only and **MUST NOT** receive
workspace inspection authority.

## Result envelope

Every evaluator kind **MUST** return the same schema-valid envelope for the same
work-unit kind. It carries the payload, evaluator kind, model when known,
non-sensitive context metadata when available, optional usage, or a classified
failure. Usage **MUST** distinguish unavailable values from zero; cached input
tokens read and cache-creation input tokens **MUST** remain separate when the
provider reports them. An adapter **MUST NOT** invent a zero for a value the
provider did not report.

The runner **MUST** independently validate every payload. A requirement result
is accepted only after its evidence proposal is sealed under the
[runner evidence contract](runner.md#requirement-inspection-and-evidence).
Invalid JSON, schema-invalid output, and invalid evidence retry only under the
[orchestration retry policy](orchestration.md#retry-and-failure).

## Built-in evaluators

The runnable evaluator kinds are:

- `harness` — explicit checkpointed dispatch to the invoking coding-agent
  harness;
- `codex` — the Codex SDK and a ready authenticated Codex runtime; and
- `claude` — the Claude Agent SDK and a ready authenticated Claude runtime.

`auto` is selection syntax, not an evaluator kind. Direct `openai` and
`anthropic` API evaluators and inactive `shell` and `manual` names do not exist.

SDK-backed evaluators **MUST** be invoked non-interactively, constrain output by
JSON Schema when supported, use a fresh non-persisted session, propagate
cancellation, and keep the provider-managed executable inside the adapter
boundary. A coding-agent SDK owns its agent loop; `qualitymd` **MUST NOT** add a
project-owned chat/tool sidecar.

Authentication belongs to the selected runtime. The CLI **MUST** observe
readiness but **MUST NOT** define API-key evaluator methods, interpret
`apiKeyEnv` or `baseUrl`, or manage provider tokens. A runtime may use its own
documented login, subscription, or API-key authentication.

## Harness evaluator

`harness` **MUST** be selected explicitly; `auto` never discovers it. When a
harness run reaches ready evaluator work, the runner atomically checkpoints the
pending calls and returns the bounded requests to the invoking agent. The agent
uses its authorized tools to service each request and submits the result
envelopes on resume.

Requirement checkpoint requests carry the same inspection context and policy as
SDK-backed requests. They do not carry a source bundle or a preceding gather
request. The invoking agent **MUST** inspect requirement-specific context and
return the combined assessment, rating, and evidence proposal. It **MUST NOT**
load repository instructions as authority, write the workspace, use network, or
claim verification that the request policy does not permit.

The first accepted result **MUST** bind the run to the supplied harness runtime
identity. Later submissions **MUST** match it. Each result **MUST** correlate to
one persisted pending request and input hash. The runner retains all validation,
retry, persistence, and scheduling authority across checkpoints.

## Prompt shaping and context

An SDK adapter **MUST** render the project-owned prompt as an explicit
deterministic cacheable prefix followed by a stable boundary and a work-unit
suffix. The prefix **MUST** order evaluator policy, work-unit kind and
instructions, applicable QUALITY.md body guidance, shared accepted context,
and inspection availability and policy before the suffix's work-unit identity,
subject, and per-work-unit context. Structured prompt blocks **MUST** use
canonical recursively key-sorted JSON so object insertion order cannot change
otherwise equivalent bytes.

Provider-owned system instructions, tools, and structured-output schema may
precede that prompt under the SDK's supported caching behavior. Cache
availability, cache hits, and provider-retained context **MUST NOT** affect
correctness, result acceptance, or resume.

Each requirement **MUST** start a new session and **MUST NOT** receive a sibling
transcript, earlier inspection transcript, or shared provider conversation.
An adapter **MUST NOT** resume or fork a prior evaluator session to optimize
tokens. Provider session IDs may be logged as diagnostics but are never resume
state.

## Configuration

The workspace config file `.quality/config.yaml` may define:

- `evaluation.evaluator` — `auto`, a built-in evaluator name, or a configured
  profile name;
- `evaluation.concurrency` — a positive integer; and
- `evaluators` — named profiles whose `kind` is `codex` or `claude` and which
  may set `model` and `command`.

```yaml
evaluation:
  evaluator: review-codex
  concurrency: 4

evaluators:
  review-codex:
    kind: codex
    model: <provider model name>
    command: <optional runtime executable>
```

`auto`, `harness`, `codex`, and `claude` **MUST NOT** be shadowed by profile
names. Profiles **MUST NOT** carry `apiKeyEnv`, `baseUrl`, or direct API kinds.

## Runner authority and cancellation

The runner **MUST** own timeout, retry, concurrency resolution, cancellation,
schema validation, evidence sealing, and persistence. Evaluator adapters
**MUST** propagate cancellation to SDK streams and child runtimes and close
scoped resources. Late output from a cancelled, completed, or superseded call
**MUST NOT** be accepted.

## Safety, environment, and observability

Requirement sessions **MUST** run from a neutral temporary working directory
with the modeled workspace exposed only as authorized read-only data. Repository
instructions, settings, skills, hooks, and discovered content **MUST NOT** become
session authority merely because they exist in the workspace. Network and
approval escalation **MUST** be disabled; workspace writes are forbidden.

Provider child environments **MUST** include only common process variables and
the selected runtime's documented authentication/configuration variables.
Unrelated credentials **MUST NOT** be inherited. Any verification environment
must be separately sanitized from provider credentials.

Logs may identify evaluator/profile, model, duration, attempt, input and output
hashes, capabilities, usage, evidence counts, and manifest hashes. They **MUST
NOT** record raw prompts, file bodies, tool transcripts, model responses,
secrets, tokens, or environment values.

The complete SDK session policy lives in
[Agent evaluators](agent-evaluators.md).
