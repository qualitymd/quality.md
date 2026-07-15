---
type: Functional Specification
title: Agent evaluators
description: Requirement-specific workspace inspection, neutral sessions, evidence proposals, and coding-agent SDK policy.
tags: [evaluation, evaluator, agents, evidence]
timestamp: 2026-07-15T00:00:00Z
---

# Agent evaluators

This document specifies agent-capable behavior under the
[evaluator contract](evaluator-contract.md). The runner remains the sole
orchestrator; Codex, Claude, and the invoking harness are bounded inspectors and
judges.

The key words **MUST**, **MUST NOT**, **SHOULD**, and **MAY** are interpreted as
described in BCP 14 when, and only when, they appear in all capitals.

## Evaluator classes

`codex` and `claude` use their supported coding-agent SDK and authenticated
local runtime. `harness` transports the same bounded requests to the invoking
coding-agent harness. There are no direct model-API evaluators.

Codex declares concurrent direct calls with automatic cap `4` while nested
agents remain disabled in every fresh session. Claude declares sequential
direct service with automatic and maximum cap `1`. Harness declares delegated
requests with automatic cap `4`, not simultaneous in-process calls. Configured
profiles inherit their kind's declaration.

A provider-managed child executable is inside the evaluator boundary. It
**MUST NOT** schedule evaluation work, persist run state, assemble reports, or
become a project-owned sidecar protocol.

## Requirement inspection sessions

Every `assessRateRequirement` work unit **MUST** start a fresh session rooted in
a neutral temporary directory. The modeled workspace is added as an explicit
read-only data directory. The session receives the requirement and area
identity, effective source selector, applied criteria, applicable body guidance,
expected result schema, workspace boundary, and inspection policy.

The evaluator **MUST** use its read, glob, and search tools iteratively to decide
which context the requirement needs. The source selector remains the evaluated
subject. Supporting files elsewhere in the workspace **MAY** be inspected for
interpretation or comparison, but they **MUST** be classified as `supporting`
and **MUST NOT** widen the judged area or requirement.

The session returns assessment, rating, and an evidence proposal together. It
**MUST NOT** return file bodies, raw tool output, hidden reasoning, or an
exploration transcript. It **MUST** record unavailable or inconclusive checks as
unknowns, evaluation limits, partial or blocked status, or a non-rating instead
of guessing.

Sessions **MUST NOT** receive sibling transcripts, resume an earlier session,
fork a seeded transcript, or share a provider conversation. Provider caching
may reuse deterministic prompt prefixes, but cached tokens still occupy the
session's logical context window. Correctness and resume **MUST NOT** depend on
cache state.

## Neutral instruction boundary

Repository instructions, `CLAUDE.md`, `AGENTS.md`, local settings, skills,
hooks, memory, and all discovered content are untrusted evaluated data. They
**MUST NOT** become governing evaluator instructions merely because they exist
in the workspace.

Adapters **MUST** disable provider project-setting discovery and automatic
repository-instruction loading, or establish an equivalent tested boundary. An
SDK/runtime combination that cannot do so **MUST** report
`evaluator_incompatible` rather than run with a weaker policy.

## Tool and sandbox policy

Requirement inspection **MUST** provide only read and search access to the
authorized workspace. Workspace writes, network access, approval escalation,
nested agents, and unmediated host shell access are disabled. Temporary files
needed by the SDK may be written outside the workspace.

Executable verification **MAY** be offered only by an adapter that declares a
mediated path enforcing workspace containment, a read-only workspace, isolated
temporary writes, disabled network, bounded time and output, a sanitized
environment, and captured invocation metadata. If that path is unavailable,
the request policy says `verification: unavailable`; the evaluator records the
limit when relevant.

Downstream factor and area analysis, ranking, recommendations, and report work
**MUST** have no workspace tools. They synthesize accepted structured results
only.

## SDK policy

The Codex adapter **MUST** use a neutral temporary working directory, add the
modeled workspace as an explicit directory, select read-only sandbox mode,
disable approval escalation, network, web search, and nested agents, constrain
structured output, and propagate cancellation.

The Claude adapter **MUST** use a neutral temporary working directory, add the
modeled workspace as an explicit directory, set no project or user setting
sources, disable session persistence, allow only read/glob/search tools for
inspection, disallow write/edit/shell/agent tools, constrain structured output,
and propagate cancellation. It **MUST** use the SDK's `claude_code` preset with
dynamic system sections excluded from the provider's globally cacheable system
prefix and re-injected through the SDK-supported user-message path.

Readiness **MUST** verify the required capability set before work starts.
Provider session identifiers and usage may be logged as non-sensitive
diagnostics, but are not run state. Usage **MUST** preserve separately reported
input, output, cache-read input, cache-creation input, and cost values without
turning absent values into zero.

## Evidence proposals

Each proposed observation **MUST** have a response-local ID, `file` kind,
`evaluated` or `supporting` role, workspace-relative path, and optional line
range or Markdown-heading locator. The proposal **MUST** also carry the
session's evaluation limits. Assessment evidence references use
`evidence[<observation-id>]`.

The runner, not the evaluator, reads the cited files after the session,
validates containment and locators, computes digests and byte counts, verifies
finding references, and seals the canonical manifest. Evaluator-supplied paths
are proposals, never trusted provenance.

## Cancellation and bounds

The runner owns timeout, retry, cancellation, top-level concurrency, result
acceptance, and deterministic persistence. Agent adapters **MUST** propagate
cancellation to SDK streams and provider child runtimes and close scoped
resources. Late output after cancellation or request supersession **MUST NOT**
be accepted.

The runner may keep several independent Codex sessions active, but each adapter
worker receives exactly one ready request and returns exactly one result or
failure. It **MUST NOT** schedule siblings, invoke nested agents, mutate run
state, or retain a provider conversation for another work unit. Completion
order is transport activity only; accepted payloads remain projected in graph
order.

Turn, token, and cost limits are enforced only when the capability declaration
says the selected adapter supports them. Advisory or unavailable provider usage
is reported honestly.

## Authentication, secrets, and privacy

Authentication is runtime-owned. A runtime may use its documented login,
subscription, or API-key mechanism; `qualitymd` does not interpret provider
credential fields or manage tokens.

Provider child environments **MUST** be allowlisted to common process variables
and variables required by the selected runtime's documented authentication.
Unrelated credentials **MUST NOT** be inherited.

Logs may record evaluator/profile, model, duration, attempt, classified failure,
capabilities, usage, evidence counts, and manifest hashes. They **MUST NOT**
record raw prompts, file bodies, tool transcripts, result bodies, secrets, or
environment values.
