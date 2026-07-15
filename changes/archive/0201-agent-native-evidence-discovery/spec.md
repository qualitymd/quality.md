---
type: Functional Specification
title: Agent-native evidence discovery — functional spec
description: Requirements for evaluator-owned context discovery, runner-validated evidence provenance, isolated coding-agent sessions, and an SDK-only evaluator surface.
tags: [evaluation, agents, cli, skill, evidence, security]
timestamp: 2026-07-14T00:00:00Z
---

# Agent-native evidence discovery — functional spec

Archived companion to the
[Agent-native evidence discovery](../0201-agent-native-evidence-discovery.md)
Change Case. This spec defines what must change in the evaluation contract; the
[design doc](design.md) owns the implementation shape.

The key words "MUST", "MUST NOT", "SHOULD", and "MAY" are to be interpreted
as described in BCP 14 when, and only when, they appear in all capitals.

## Background / motivation

A quality requirement determines which context is relevant. A deterministic
runner can resolve paths and package bytes, but it cannot know in advance which
implementation, documentation, tests, configuration, history, or cross-reference
will prove a particular requirement. Reusing one capped area bundle for every
requirement makes incidental traversal order part of judgment and prevents the
evaluator from recovering omitted evidence.

Coding-agent runtimes already provide the iterative search, inspection, context
management, and tool loop this work needs. The CLI should coordinate those
runtimes, constrain their authority, validate their results, and preserve the
evidence of record. It should not recreate their harness around a raw model API.

## Requirements

### R1 — Evaluation ownership boundary

The CLI runner **MUST** own model and scope validation, work-graph construction,
scheduling, evaluator selection and invocation, retry and cancellation,
structured-result validation, atomic persistence, resume, synthesis inputs,
report generation, and deterministic output ordering.

> Rationale: these are mechanical and artifact-integrity responsibilities that
> remain stable across evaluator implementations. — 0201
>
> Durable spec: modify `specs/evaluation/evaluation.md`,
> `specs/evaluation/runner.md`, and `specs/evaluation/orchestration.md`.

For requirement judgment, the selected evaluator **MUST** own iterative context
discovery, inspection, evidence selection, assessment, findings, and rating. The
runner **MUST NOT** preselect files, construct an area-wide source bundle, or
truncate candidate evidence before the evaluator sees the requirement.

> Rationale: evidence relevance is part of requirement judgment, not a
> deterministic preprocessing operation. — 0201
>
> Durable spec: modify `specs/evaluation/runner.md`,
> `specs/evaluation/protocol.md`, and
> `specs/evaluation/agent-evaluators.md`.

### R2 — Source and supporting context

An area's effective `source` **MUST** identify the subject or starting boundary
being evaluated. It **MUST NOT** mean “all context the evaluator may inspect,” a
precomputed prompt payload, or a file-read permission boundary.

> Durable spec: modify `SPECIFICATION.md`,
> `specs/evaluation/evaluation.md`, and
> `specs/evaluation/agent-evaluators.md`.

An evaluator **MAY** inspect supporting context elsewhere in the authorized
workspace when a requirement needs comparison, interpretation, or verification.
It **MUST** distinguish evidence about the evaluated subject from supporting
context and verification observations, and **MUST NOT** silently widen the area,
entity, or requirement it judges.

> Rationale: an API implementation may need to be compared with its design or
> user documentation without turning those documents into the evaluated API
> area. — 0201
>
> Durable spec: modify `SPECIFICATION.md`,
> `specs/evaluation/evaluator-contract.md`, and
> `specs/evaluation/records/payload-kinds.md`.

Path, glob, and prose source selectors **MUST** remain valid source forms. The
runner **MAY** parse them for model validation, display, and workspace-safety
checks, but the evaluator **MUST** interpret and investigate the selector for
the requirement rather than receive a runner-resolved bundle.

> Durable spec: modify `SPECIFICATION.md`,
> `specs/evaluation/runner.md`, and `specs/cli/evaluation-run.md`.

### R3 — Requirement inspection sessions

Each requirement judgment **MUST** run in a fresh isolated evaluator session.
Its initial request **MUST** include the model and area identity, requirement,
effective source selector, applied rating criteria, applicable body guidance,
expected result schema, authorized workspace boundary, tool policy, and
evaluation limits.

> Durable spec: modify `specs/evaluation/evaluator-contract.md`,
> `specs/evaluation/agent-evaluators.md`, and
> `specs/evaluation/protocol.md`.

The session **MUST** be able to search and read the authorized workspace
iteratively. It **MAY** perform explicitly permitted verification when static
inspection is insufficient. One session **MUST** return the requirement
assessment and rating together so the rating cannot rely on evidence absent
from its paired assessment.

> Durable spec: modify `specs/evaluation/protocol.md` and
> `specs/evaluation/records/payload-kinds.md`.

When relevant evidence cannot be accessed, inspected, or safely verified, the
evaluator **MUST** record the resulting unknowns, missing evidence, evaluation
limits, and non-rated or blocked status instead of guessing. Ordinary evidence
insufficiency **MUST** be a judgment outcome, not an evaluator infrastructure
failure.

> Durable spec: modify `specs/evaluation/evaluation.md`,
> `specs/evaluation/runner.md`, and
> `specs/evaluation/records/payload-kinds.md`.

### R4 — Inspection authority and isolation

Requirement sessions **MUST** receive read-only access to the modeled workspace,
no network access by default, no approval escalation, no workspace writes, and
only the minimum environment needed to start the authenticated agent runtime.
Temporary files required by an isolated sandbox **MAY** be written outside the
workspace.

> Rationale: delegating evidence selection does not delegate authorization. —
> 0201
>
> Durable spec: modify `specs/evaluation/evaluator-contract.md` and
> `specs/evaluation/agent-evaluators.md`.

Repository instruction files, local settings, skills, hooks, and discovered
content **MUST** be treated as untrusted evaluated data. They **MUST NOT** become
governing instructions for the evaluator session merely because they exist in
the workspace. An SDK path that cannot establish this neutral instruction
boundary **MUST** be reported as `evaluator_incompatible`.

> Durable spec: modify `specs/evaluation/agent-evaluators.md` and
> `specs/evaluation/evaluator-contract.md`.

Executable verification **MUST** pass through a runner-selected mediated path
provided by the agent SDK or its thin adapter. That path **MUST** enforce
workspace containment, a read-only workspace, isolated temporary writes,
disabled network, bounded time and output, a sanitized environment, and
captured invocation metadata. Unmediated host shell access **MUST NOT** be
available to a requirement session.

> Durable spec: modify `specs/evaluation/agent-evaluators.md` and
> `specs/evaluation/runner.md`.

### R5 — Evidence manifest and validation

Every successful requirement response **MUST** carry a distinct evidence
manifest alongside its assessment and rating. The manifest **MUST** identify the
effective source selector, the requirement, each material observation used,
each observation's `evaluated` or `supporting` role, stable workspace-relative
locator, capture time, content digest, and the session's stated limits. A future
adapter that declares mediated verification **MUST** extend the shared schema
with runner-observed command provenance before command evidence is accepted;
the current built-ins declare verification unavailable.

> Durable spec: modify `specs/evaluation/evaluator-contract.md`,
> `specs/evaluation/evaluation-json.md`, and
> `specs/evaluation/records/payload-kinds.md`.

Before accepting a requirement result, the runner **MUST** validate manifest
shape, workspace containment, locator readability, observation hashes,
verification provenance, and every finding evidence reference. A finding
**MUST NOT** cite a source that is absent from its accepted evidence manifest.

> Rationale: the evaluator selects evidence; the runner decides whether that
> evidence is safe, internally consistent, and durable enough to accept. — 0201
>
> Durable spec: modify `specs/evaluation/runner.md`,
> `specs/evaluation/evaluation-json.md`, and
> `specs/evaluation/records/payload-kinds.md`.

The run artifact **MUST** persist accepted evidence manifests atomically with
their requirement results. It **MUST NOT** persist raw prompts, hidden reasoning,
complete tool transcripts, credentials, or full source bodies by default.
Reports **MAY** render validated locators, concise evidence statements, hashes,
and recorded limits.

> Durable spec: modify `specs/evaluation/evaluation-json.md`,
> `specs/evaluation/runner.md`, and
> `specs/evaluation/reports/report-tree.md`.

### R6 — Work graph and downstream synthesis

The work graph **MUST NOT** contain `resolveSource` work units. Each
`assessRateRequirement` unit **MUST** depend on its requirement frame and begin
its own inspection session; requirements in the same area **MUST NOT** share an
immutable source package or provider conversation.

> Durable spec: modify `specs/evaluation/orchestration.md` and
> `specs/evaluation/protocol.md`.

Factor analysis, area analysis, finding ranking, recommendation generation, and
recommendation ranking **MUST** remain tools-off synthesis over validated,
accepted run results. Those units **MUST NOT** discover new workspace evidence.

> Durable spec: modify `specs/evaluation/evaluation.md`,
> `specs/evaluation/orchestration.md`, and
> `specs/evaluation/protocol.md`.

Bounded parallelism, deterministic persisted order, cancellation, retry,
accepted-result durability, and harness batching **MUST** continue to apply to
the simplified graph.

> Durable spec: modify `specs/evaluation/runner.md` and
> `specs/evaluation/orchestration.md`.

### R7 — Evaluator methods and authentication

The runnable evaluator methods **MUST** be `harness`, `codex`, and `claude`.
`auto` **MUST** select a ready Codex agent runtime first, then a ready Claude
agent runtime, and **MUST** fail clearly when neither is usable.

> Durable spec: modify `specs/evaluation/evaluator-contract.md`,
> `specs/evaluation/agent-evaluators.md`, and
> `specs/cli/evaluation-run.md`.

The direct `openai` and `anthropic` evaluator methods, inactive reserved
`shell` and `manual` evaluator names, direct HTTP adapters, API-profile fallback,
and API-specific evaluator failures **MUST** be removed in the clean cutover.
Configured profiles **MUST** use only the `codex` or `claude` agent-runtime kind.

> Rationale: a raw model API plus tools would make the CLI responsible for an
> agent loop it should consume from a supported SDK. Unimplemented reserved
> names add grammar without capability. — 0201
>
> Durable spec: modify `specs/evaluation/evaluator-contract.md` and
> `specs/cli/evaluation-run.md`.

Authentication **MUST** remain distinct from evaluator method. The CLI **MUST**
verify runtime readiness but **MUST NOT** define API-key evaluators, interpret
`apiKeyEnv` or `baseUrl` profile fields, or manage provider tokens. A selected
agent runtime **MAY** authenticate through login, subscription, or its own
documented API-key mechanism.

> Durable spec: modify `specs/evaluation/evaluator-contract.md`,
> `specs/evaluation/agent-evaluators.md`, and
> `specs/skills/quality-skill/evaluation.md`.

### R8 — Harness evaluator flow

The `harness` evaluator **MUST** remain an explicit transport selected by the
invoking `/quality` skill or caller, never by CLI `auto` discovery. A harness
checkpoint request **MUST** carry the same requirement inspection request and
tool policy as an SDK-backed session, without a source bundle or a preceding
resolution request.

> Durable spec: modify `specs/evaluation/evaluator-contract.md`,
> `specs/evaluation/runner.md`, and
> `specs/skills/quality-skill/workflows/evaluate.md`.

The invoking agent **MUST** inspect the workspace with its available authorized
tools and submit the combined judgment and evidence manifest. The CLI remains
the checkpoint, validation, persistence, and resume harness; it **MUST NOT**
start a nested agent or request an API key for harness judgment.

> Durable spec: modify `specs/skills/quality-skill/evaluation.md` and
> `specs/skills/quality-skill/workflows/evaluate.md`.

### R9 — Artifact, resume, and determinism contract

The authoritative run artifact **MUST** make one clean schema-version break that
replaces per-area `sources` records with per-requirement evidence manifests.
The implementation **MUST NOT** add a migration, fallback reader, dual writer,
or compatibility alias for the replaced evaluator and source-package shapes.

> Durable spec: modify `specs/evaluation/evaluation-json.md`.

Pending harness requests **MUST** remain reconstructible from the model
snapshot, work graph, requirement context, evaluator policy, and workspace
boundary. Accepted requirement results and their evidence manifests **MUST** be
immutable resume inputs, and retries of unaccepted requirement work **MUST** use
a fresh inspection session.

> Durable spec: modify `specs/evaluation/orchestration.md`,
> `specs/evaluation/evaluation-json.md`, and
> `specs/cli/evaluation-run.md`.

Documentation **MUST** describe determinism as common orchestration, validation,
and artifact guarantees across evaluator methods, not as identical evidence,
judgments, or ratings across agent runtimes or repeated runs.

> Durable spec: modify `specs/evaluation/evaluation.md` and user guidance.

### R10 — Skill and user experience alignment

The `/quality` skill **MUST** foreground harness-backed agent judgment when it
already runs inside a capable coding-agent session, and **MUST** explain Codex
or Claude agent-runtime fallback without presenting direct API evaluators or
CLI-managed credentials.

> Durable spec: modify `specs/skills/quality-skill/evaluation.md`,
> `specs/skills/quality-skill/workflows/evaluate.md`, and the matching runtime
> skill files.

Dry-run and progress output **MUST** show the selected evaluator, readiness and
isolation capabilities, requirement work count, concurrency, source selectors,
and inspection policy. They **MUST NOT** claim source resolution, bundle sizes,
or static per-area dispatch plans.

> Durable spec: modify `specs/cli/evaluation-run.md`.

## Requirement-set review

R1–R3 put semantic context discovery with the requirement judge while retaining
mechanical authority in the runner. R4 constrains that delegation, and R5 gives
the runner a concrete acceptance and audit boundary. R6 removes the obsolete
gathering graph without weakening tools-off synthesis. R7–R8 leave exactly two
agent SDK methods plus the explicit outer-harness transport and separate
authentication from method selection. R9 preserves honest resume and
reproducibility claims through a clean artifact break. R10 makes the agent-first
experience teach the same contract.

The set covers the motivating failures without granting broader workspace,
network, or mutation authority. Each requirement is verifiable through graph
tests, adapter capability checks, sandbox and instruction-isolation tests,
malicious-workspace fixtures, evidence-manifest validation, resume tests,
generated-schema checks, CLI contract tests, and authenticated end-to-end Codex
and Claude acceptance runs.

## Durable spec changes

### To add

None. The existing evaluation and agent-evaluator specs are the correct durable
homes for the revised contract.

### To modify

- `SPECIFICATION.md` — make `source` the evaluated subject selector and permit
  explicitly classified supporting context within the authorized workspace.
- `specs/evaluation/evaluation.md` — state the new ownership and determinism
  invariants.
- `specs/evaluation/runner.md` — replace source packaging with inspection policy,
  evidence validation, and mediated verification.
- `specs/evaluation/protocol.md` — remove `resolveSource` and define the
  inspection-bearing requirement move.
- `specs/evaluation/orchestration.md` — simplify the graph, dependencies,
  persistence, retry, and resume rules.
- `specs/evaluation/evaluator-contract.md` — revise capabilities, request/result
  envelopes, supported methods, profiles, and harness checkpoints.
- `specs/evaluation/agent-evaluators.md` — define fresh neutral SDK sessions,
  workspace access, tool mediation, and runtime-owned authentication.
- `specs/evaluation/evaluation-json.md` — replace `sources` with accepted
  per-requirement evidence manifests in a new schema version.
- `specs/evaluation/records/payload-kinds.md` — bind finding evidence references
  to the accepted evidence manifest and distinguish evidence roles.
- `specs/evaluation/reports/report-tree.md` — render accepted evidence locators
  and limits without introducing or regathering evidence.
- `specs/cli/evaluation-run.md` — revise evaluator selection, configuration,
  dry-run, harness requests, resume, output, and failures.
- `specs/skills/quality-skill/evaluation.md` and
  `specs/skills/quality-skill/workflows/evaluate.md` — teach harness agents to
  discover and return requirement evidence directly.

### To rename

None.

### To delete

None. Obsolete sections and contracts are removed within the modified durable
specs rather than deleting an entire durable spec.
