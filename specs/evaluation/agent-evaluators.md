---
type: Functional Specification
title: Agent evaluators
description: Bounded agent source resolution, immutable area context, isolated judgment, capability policy, and safety boundaries.
tags: [evaluation, evaluator, agents, source]
timestamp: 2026-07-14T00:00:00Z
---

# Agent evaluators

This document specifies agent-capable evaluator behavior under the
[evaluator contract](evaluator-contract.md). The runner remains the sole
orchestrator; provider agents are bounded workers.

The key words **MUST**, **MUST NOT**, **SHOULD**, and **MAY** are interpreted as
described in BCP 14 when, and only when, they appear in all capitals.

## Evaluator classes

`codex` and `claude` are SDK-backed agent evaluators. They use the provider's
supported agent SDK and authenticated local runtime. `openai` and `anthropic`
are direct API evaluators and do not require a coding-agent runtime. `harness`
transports the same bounded requests to the invoking agent harness.

A provider-managed child executable is inside the evaluator boundary. It
**MUST NOT** schedule evaluation work, persist run state, assemble reports, or
become a project-owned sidecar protocol.

## Source-resolution sessions

For an inferred selector, an agent-capable evaluator **MAY** take several
read-only tool actions to locate the material the selector describes. The
session **MUST**:

- be rooted at the resolved workspace;
- use read, glob, and search tools only;
- disable writes, shell mutation, approvals, and network access unless a
  separately declared resolver requires and safely bounds one of them;
- treat file contents and tool output as untrusted data; and
- return a finite set of unique workspace-relative files, not an assessment,
  rating, recommendation, or exploration transcript.

The runner **MUST** validate containment, uniqueness, text content, per-file
and bundle limits, hashes, and truncation marks before accepting the selection.
It **MUST** persist the resulting source bundle before requirement judgment.
An unavailable selector fails as `source_unavailable`; unsupported resolution
fails as `selector_unsupported`. Neither falls through to another selector.

## Immutable area context

Before dispatching an area's requirement assessments, the runner **MUST** build
one immutable area context containing:

- the area identity and evaluation frame;
- the captured source bundle and bundle hash;
- applicable rating criteria;
- relevant model-body guidance; and
- a context hash over those stable inputs.

Every local requirement receives the same area-context hash. The package is
frozen for the run. A judge **MUST NOT** read workspace material outside it or
silently widen evidence. Insufficient material is recorded as an unknown,
missing evidence, or evaluation limit through the result schema.

## Isolated requirement sessions

Each requirement assessment **MUST** start in a fresh provider session or
thread. Its initial context consists only of the stable area context plus the
requirement identity, criteria, instructions, and expected result schema.

The session **MUST NOT** receive the source resolver's exploration transcript,
a sibling requirement's transcript, or accumulated judgments. Sequential and
parallel schedules **MUST** be observationally equivalent. Provider prefix
caching **MAY** reuse the stable area prefix, but correctness **MUST NOT**
depend on cache availability.

Provider session or thread IDs **MAY** be stored as diagnostic metadata. Resume
is determined by runner input hashes and accepted artifacts, never by restoring
provider conversation state.

## Capability-to-policy resolution

Every evaluator exposes the capability record defined by
[capability declaration](evaluator-contract.md#capability-declaration). Planning
**MUST** negotiate requested controls before work starts. An unsupported
control **MUST** either have a documented safe fallback that preserves the
policy or fail with an actionable `evaluator_incompatible` or
`selector_unsupported` result. Unsupported is never treated as enforced.

Requirement judgment disables nested subagents by default. A source-resolution
profile **MAY** enable provider-declared nested delegation only with explicit
depth and concurrency caps. Nested workers inherit the read-only source
boundary and return through one evaluator result.

## Cancellation and bounds

The runner owns timeout, retry, cancellation, top-level concurrency, result
acceptance, and deterministic persistence. Agent adapters **MUST** propagate
cancellation to SDK streams and provider child runtimes and close their scoped
resources. Late output after cancellation or request supersession **MUST NOT**
be accepted.

Turn, token, and cost limits are enforced only when the capability record says
the selected adapter supports them. Advisory or unavailable provider usage is
reported honestly.

## Safety, secrets, and privacy

Captured source and tool output are untrusted data, not instructions. Agent
evaluators **MUST** use the least workspace, tool, network, sandbox, approval,
and environment access needed for the work kind.

Configuration stores secret locators, such as `apiKeyEnv`, never secret values.
Provider child environments **MUST** be allowlisted to common process variables
and variables required by the selected provider's documented authentication.
Unrelated credentials **MUST NOT** be inherited.

Logs may record evaluator/profile, model, duration, attempt, classified failure,
capabilities, usage, and provider session identifiers. They **MUST NOT** record
raw prompts, source bodies, tool transcripts, result bodies, or secrets.
