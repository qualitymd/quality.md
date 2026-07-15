---
type: Change Case
title: Agent-native evidence discovery
description: Move requirement-specific context discovery from deterministic source packaging into coding-agent evaluators while preserving runner safety, provenance, and reproducibility.
status: Done
tags: [evaluation, agents, cli, skill, evidence, security]
timestamp: 2026-07-14T00:00:00Z
---

# Agent-native evidence discovery

## Motivation

The evaluation runner currently resolves each area's `source`, captures one
bounded source bundle, and gives every requirement in the area the same package.
That makes evidence collection deterministic, but it asks the CLI to predict
the context a coding agent needs before the requirement is known. In a recent
large-repository evaluation, alphabetical traversal and bundle caps consumed the
root-area budget on incidental files while omitting the instructions, docs, and
implementation surfaces needed by most requirements. The isolated judge could
not recover because it had no workspace inspection tools.

This is the wrong ownership boundary. The runner should remain the deterministic
orchestration and artifact harness. A coding-agent evaluator should decide what
to inspect for each requirement, using the area source as the evaluated subject
boundary and following requirement-specific references as supporting context.
The runner should then validate and persist what the evaluator actually used.

The same correction removes a second accidental responsibility. Direct
OpenAI- and Anthropic-API evaluators require the CLI to implement an agent loop,
tool routing, sandboxing, context management, provider behavior, and credential
handling. Supported coding-agent SDKs already provide that harness. API keys may
authenticate those runtimes, but an API key is not an evaluator method.

## Scope

Covered:

- redefine the runner/evaluator boundary around evaluator-owned,
  requirement-specific evidence discovery;
- treat area `source` as the evaluated subject selector, not a precomputed
  prompt payload or a boundary on all supporting context;
- run each requirement judgment in a fresh coding-agent session with read-only
  workspace inspection and optional mediated, sandboxed verification;
- replace area source bundles and `resolveSource` work units with validated,
  per-requirement evidence manifests;
- preserve deterministic planning, scheduling, schema validation, atomic
  persistence, resume, synthesis, report generation, and failure handling in
  the CLI;
- support `harness`, `codex`, and `claude` evaluator methods, with `auto`
  selecting a ready coding-agent runtime;
- remove raw `openai` and `anthropic` direct-API evaluators and their
  API-specific configuration and failure paths;
- keep authentication a property of the selected agent runtime rather than the
  evaluator grammar; and
- synchronize the format semantics, evaluation specs, CLI, `/quality` skill,
  generated artifacts, user docs, tests, and release notes.

Deferred:

- external evidence connectors and network-enabled evaluation;
- evaluator writes to the modeled workspace;
- general custom evaluator/plugin protocols beyond configured Codex and Claude
  agent-runtime profiles; and
- distributed or remote execution beyond the existing harness checkpoint
  transport.

Non-goals:

- making requirement judgment deterministic across models or repeated runs;
- turning the CLI into a general coding-agent harness;
- letting an evaluator expand the entity or area being judged merely because it
  may inspect supporting context;
- auto-loading repository instructions, skills, or settings as evaluator
  authority; and
- preserving schema-version-7 source packages, direct-API evaluator aliases,
  or other compatibility paths after the clean early-alpha cutover.

## Affected artifacts

The inventory follows the current source-package flow from model semantics
through planning, evaluator invocation, persistence, reporting, skill guidance,
and generated documentation.

- **Format semantics:** `SPECIFICATION.md` and generated
  `mintlify/specification.mdx` clarify that `source` identifies the evaluated
  subject while supporting evidence may be discovered elsewhere within the
  authorized workspace.
- **Evaluation specs:** `specs/evaluation/evaluation.md`, `runner.md`,
  `protocol.md`, `orchestration.md`, `evaluator-contract.md`,
  `agent-evaluators.md`, `evaluation-json.md`, and
  `records/payload-kinds.md` replace gather-before-judge source packages with
  requirement inspection sessions and evidence manifests.
- **CLI spec:** `specs/cli/evaluation-run.md` narrows evaluator selection and
  dry-run output to `harness`, `codex`, `claude`, and `auto`, and describes the
  new request and resume contract.
- **Skill specs and runtime:** `specs/skills/quality-skill/evaluation.md`,
  `specs/skills/quality-skill/workflows/evaluate.md`,
  `skills/quality/SKILL.md`, and `skills/quality/workflows/evaluate.md` teach
  harness agents to inspect requirement-specific context and return evidence
  provenance rather than serving `resolveSource` checkpoints.
- **Domain and application code:** `src/domain/evaluator/`,
  `src/domain/evaluation/graph.ts`, `protocol.ts`, and `result.ts`, plus
  `src/application/evaluation-run.ts`, `evaluation-execute.ts`,
  `evaluation-provider.ts`, and `evaluation-resume.ts` change work units,
  evaluator envelopes, evidence validation, selection, and persistence.
- **Services and adapters:** `src/adapters/evaluator.ts`,
  `src/services/source.ts`, and `src/services/workspace.ts` remove direct API
  adapters and static packaging, add isolated inspection-session policy, and
  seal evaluator evidence into runner-owned provenance.
- **CLI and configuration:** `src/cli/app.ts` and the workspace configuration
  schema remove direct API kinds, `apiKeyEnv`, and `baseUrl`; configured agent
  profiles retain only agent-runtime options such as kind, command, and model.
- **Artifacts and reports:** `src/assets/evaluation-data.schema.json`,
  `src/assets/evaluation-examples.json`, report rendering where evidence
  locators are shown in `src/domain/evaluation/report.ts`, and associated
  generated galleries move to the new `evaluation.json` schema.
- **Tests:** evaluator selection, graph, protocol, provider, source/workspace,
  resume, artifact-schema, CLI, security-boundary, evidence-provenance, report,
  and generated-artifact tests under `test/`.
- **User guidance:** `README.md`, `mintlify/codex.mdx`,
  `mintlify/claude-code.mdx`, generated `mintlify/cli.mdx`, and `CHANGELOG.md`.
- **Project model:** no planned `QUALITY.md` factor or requirement change; the
  project model already evaluates agent and CLI quality at the appropriate
  level.
- **Dependencies:** `package.json` and `bun.lock` retain the supported Codex and
  Claude agent SDKs; they change only if implementation needs a version with
  the required isolation or tool-policy controls.
- **Install and scaffold:** no generated workspace or installer contract change
  is planned.

## Children

- [Functional spec](0201-agent-native-evidence-discovery/spec.md) — ownership,
  source, inspection, provenance, evaluator, safety, and resume requirements.
- [Design doc](0201-agent-native-evidence-discovery/design.md) — the revised
  work graph, inspection-session boundary, evidence ledger, SDK policies, and
  clean cutover.
