---
type: Change Case
title: Derive-values conformance refactor
description: Bring src/ into derive-values conformance and reshape evaluation-execute around explicit Effect composition and pure domain functions while preserving public CLI behavior.
status: Done
tags: [refactor, effect, typescript, evaluation, cli]
timestamp: 2026-07-14T00:00:00Z
---

# Derive-values conformance refactor

Status note: case is **Done**. Implementation, durable specs, regression tests,
byte-level artifact comparison, and the full local gate passed; acceptance
evidence is recorded in the child review ledger. The style-guide sections this
case enforces landed independently in
[`docs/guides/effect-typescript-style.md`](../../docs/guides/effect-typescript-style.md)
("Compose Effect workflows explicitly" and "Derive values; do not accumulate
them").

## Motivation

A readability review of
[`src/application/evaluation-execute.ts`](../../src/application/evaluation-execute.ts)
found one 225-line `Effect.gen` interleaving parsing, pure planning math,
filesystem writes, and output rendering, held together by shared mutable
accumulators (`payloads`, `workUnits`, `areaSources`, `pending`, `requests`)
that a reader must replay temporally to understand. The same review surfaced a
codebase-wide pattern: collections built by threading mutable accumulators
across loops rather than derived as expressions, which the style guide now
prohibits in its "Derive values; do not accumulate them" section.

The affected-artifact sweep also found real drift that the accumulator style
concealed: run-folder enumeration and next-number logic is copied five times
(`evaluation-create`, `evaluation-execute`, `evaluation-run`, `status`,
`evaluation-inspect`) with inconsistent manifest precedence and an undocumented
exclusion of folder slugs containing `quality` that `evaluation-create` lacks
but the other four apply. No durable spec states the fallback rule, and the
exclusion arrived wholesale in the Go-to-TypeScript port with no recorded
reason.

This case brings `src/` into derive-values conformance, applies the adjacent
Effect-composition rules to the operations it reshapes, decomposes the flagship
file into pure, unit-tested domain functions, and replaces the five enumeration
copies with one shared, spec-backed implementation — all without changing public
CLI behavior beyond the one decided enumeration edge case.

A comparison with current Effect and alchemy-effect practice also exposed a
problem in the first design draft: batching hashes with
`Effect.promise(() => Promise.all(...))` would leave Effect's structured
concurrency, while reusable operations would remain anonymous wrappers around
`Effect.gen`. The revised guide and design keep those operations named and keep
batch traversal, failure, interruption, and ordering inside Effect.

## Scope

Covered:

- a conformance sweep of `src/` against the guide's "Derive values; do not
  accumulate them" section, including `ReadonlyArray`/`readonly` collection
  fields on exported types;
- conformance with the guide's Effect-composition rules for every reusable
  Effect operation introduced or materially reshaped by this case, including
  named `Effect.fn` operations, typed failure handling, and Effect-native batch
  traversal;
- decomposition of `executeHarnessRun` so frame construction, work-unit
  derivation, request assembly, artifact assembly, and summary rendering are
  pure functions in `src/domain/evaluation/`;
- ready-unit scheduling as a pure function in `src/domain/evaluation/graph.ts`;
- one shared run-folder enumeration implementation, with the recognition and
  numbering-fallback rule decided and stated in the durable CLI spec; and
- unit tests for every extracted domain function.

Deferred:

- refining `evaluation-execute`'s collapsed `FileSystemFailure` error channel
  into per-operation typed failures (behavior-adjacent; separate case);
- a full-`src/` conformance sweep for the guide's newly added Effect-composition
  rules outside operations this case materially reshapes;
- a matching conformance sweep of the `test/` tree; and
- lint automation for accumulator patterns.

Non-goals:

- any public CLI grammar, output, or artifact change beyond the enumeration
  edge case R6 decides;
- performance optimization; and
- reshaping the services/adapters architecture.

## Affected artifacts

Derived from a repo sweep for accumulator patterns (`let`/`.push` per file),
the five enumeration call sites, and the durable homes of run recognition.

- **Durable docs:** `docs/guides/effect-typescript-style.md` — updated with the
  derive-values rule plus named Effect composition, structured batch/repetition,
  service lifetime, and deterministic-time test guidance.
- **Durable specs:** `specs/cli.md` gains the shared run-folder recognition and
  numbering rule; `specs/cli/evaluation-create.md` references it for
  next-number computation. No other durable spec changes.
- **Domain code:** new `src/domain/evaluation/frames.ts`; `graph.ts`
  (ready-unit selection); `protocol.ts` (request assembly and optional fields
  instead of sentinels); `run.ts` (artifact/summary assembly and run-folder name
  classification); conformance sweep over `report.ts`, `plan.ts`, `data.ts`,
  `result.ts`, `src/domain/lint/lint.ts`, `src/domain/model/model.ts`, and
  `document.ts`.
- **Application code:** `evaluation-execute.ts` (flagship decomposition and
  Effect-native hashing traversal);
  `evaluation-create.ts`, `evaluation-run.ts`, `status.ts`, and
  `evaluation-inspect.ts` (shared enumeration); conformance sweep over
  `evaluation-data.ts`, `evaluation-resume.ts`, `evaluation-report.ts`,
  `evaluation-provider.ts`, `update.ts`, `model.ts`, `init.ts`, and `lint.ts`.
- **Services and adapters:** conformance sweep over `src/services/source.ts`,
  `src/services/workspace.ts`, and `src/adapters/evaluator.ts`.
- **CLI code:** no changes planned; `src/cli/` showed no accumulator findings.
- **Tests:** new unit tests under `test/domain/evaluation/` for extracted
  functions and the enumeration classification; existing contract tests pass
  unweakened; a shared-enumeration test replaces per-command duplicates where
  they exist.
- **User guidance:** `CHANGELOG.md` notes the removed `quality`-slug
  enumeration exclusion. No README, mintlify, or generated-artifact changes.
- **Skill:** no changes; no skill-visible behavior moves.
- **Project model:** no `QUALITY.md` change.
- **Install/scaffold and dependencies:** no changes.

## Children

- [Functional spec](0202-derive-values-refactor/spec.md) — behavior
  preservation, derivation and Effect-workflow conformance, decomposition,
  scheduling, and enumeration requirements.
- [Design doc](0202-derive-values-refactor/design.md) — module map, derivation
  idioms, structured hashing, unified enumeration, and the exclusion decision.
- [Review ledger](0202-derive-values-refactor/review.md) — requirement status,
  byte-preservation evidence, enumeration coverage, and local gate results.
