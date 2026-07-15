---
type: Change Case
title: Evaluator prompt cache efficiency
description: Make fresh evaluator sessions share deterministic cacheable prompt prefixes and report provider cache usage without coupling judgment contexts.
status: In-Progress
tags: [evaluation, evaluator, prompt-caching, tokens, codex, claude]
timestamp: 2026-07-15T00:00:00Z
---

# Evaluator prompt cache efficiency

Status note: **In-Progress**; R1–R6 passed individual and set-level review, the
technical design is settled, and implementation may now update the affected
runtime, test, durable-spec, and release-note surfaces.

## Motivation

The evaluator architecture intentionally opens a fresh session for every
requirement so sibling inspection transcripts cannot influence one another.
That isolation is correct, but it makes token efficiency depend on provider
prompt caching: repeated evaluator policy, work-kind instructions, model body,
area frame, inspection policy, tools, and output schema should form an identical
prefix, with only the requirement delta varying.

The current renderer defeats most of that opportunity by putting work-unit ID,
subject, and other varying fields before the shared model and body context. The
Codex adapter reports cache reads, while the Claude SDK reports both cache reads
and cache creation but the adapter drops them. Claude sessions also use a fresh
random working directory inside the default system-prompt context even though
the pinned SDK can move those dynamic sections after its globally cacheable
prefix.

This case makes caching an adapter optimization over the existing independent
session contract. It does not introduce a shared transcript, provider resume
state, direct model API, or correctness dependency on a cache hit.

## Scope

Covered:

- split evaluator prompt rendering into a deterministic cacheable prefix and a
  work-unit-specific suffix;
- order work-kind instructions, body guidance, shared area context, and
  inspection policy before work-unit identity and context;
- serialize structured prompt blocks canonically so object insertion order
  cannot cause a false cache miss;
- configure Claude's existing preset system prompt to keep its dynamic runtime
  sections outside the cross-session cacheable prefix;
- preserve provider-reported cache reads and cache writes as optional evaluator
  usage fields and in run-local evaluator-call logs;
- focused unit and integration coverage; and
- two repeated, small, scoped live evaluations to inspect cache usage and
  latency before review.

Deferred:

- cache-aware warm-up or cohort scheduling; live evidence from this change must
  show a concrete cold-burst problem before scheduling changes earn a case;
- explicit OpenAI prompt-cache keys or breakpoints, which the pinned Codex SDK
  does not expose; and
- run-level cost dashboards or a new public usage-summary command.

Non-goals:

- sharing, resuming, or forking evaluator transcripts;
- reducing logical context-window occupancy through provider caching;
- changing evaluator selection, work-graph dependencies, concurrency defaults,
  judgment semantics, ratings, evidence sealing, reports, or `evaluation.json`;
- adding a direct OpenAI or Anthropic API evaluator; and
- changing the `/quality` workflow.

## Affected artifacts

Derived from searches for prompt rendering, usage fields, evaluator-call
logging, fresh-session rules, provider SDK options, and release/version surfaces.

- **Change record:** this parent, `spec.md`, `design.md`, and `review.md` under
  `changes/0203-evaluator-prompt-cache-efficiency/`; changes bundle indexes and
  log; archived together when done.
- **Durable specs:** `specs/evaluation/evaluator-contract.md` strengthens prompt
  prefix ordering and cache-usage preservation;
  `specs/evaluation/agent-evaluators.md` records provider cache shaping without
  weakening neutral fresh sessions; `specs/evaluation/runner.md` records cache
  read/write call-log fields; `specs/evaluation/log.md` records the revision.
- **Domain code:** `src/domain/evaluator/context.ts` owns prompt parts and
  canonical rendering; `src/domain/evaluator/types.ts` owns the optional cache
  write usage field.
- **Adapter code:** `src/adapters/evaluator.ts` configures Claude's stable system
  prefix and maps provider cache usage without inventing unavailable values.
- **Application code:** no behavior change planned; existing response forwarding
  and evaluator-call logging preserve the expanded usage object unchanged.
- **Tests:** `test/adapters/evaluator.test.ts` covers exact prefix stability,
  suffix variation, ordering, canonicalization, Claude cache shaping, and usage
  mapping; `test/integration/evaluation-provider.test.ts` verifies cache fields
  survive into `logs/evaluator-calls.jsonl` without entering `evaluation.json`.
- **Release notes and metadata:** `CHANGELOG.md`; release preparation advances
  the CLI and skill patch metadata while retaining the `>=0.32.0 <0.33.0`
  compatibility line.
- **Durable docs:** no README, Mintlify, install, scaffold, or user-guide change;
  the behavior is an internal evaluator optimization already covered by the
  evaluation specs.
- **Bundled skill runtime/spec:** no workflow behavior change; release metadata
  only.
- **Format specification and project model:** no `SPECIFICATION.md` or
  `QUALITY.md` change; conforming model/evaluation meaning is unchanged.
- **Dependencies and generated artifacts:** no dependency or generated-artifact
  change planned.

## Children

- [Functional spec](0203-evaluator-prompt-cache-efficiency/spec.md) — cacheable
  prompt, usage, isolation, and verification requirements.
- [Design doc](0203-evaluator-prompt-cache-efficiency/design.md) — explicit
  prompt parts, canonical blocks, provider usage adapters, and cache-stable
  Claude configuration.
- Review ledger — added during implementation and completed before archival.
