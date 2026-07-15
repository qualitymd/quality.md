---
type: Design Doc
title: Derive-values conformance refactor — design
description: Module map for the evaluation-execute decomposition, derivation and Effect-composition idioms, structured hashing, and one run-folder enumeration.
tags: [refactor, effect, typescript, evaluation, cli]
timestamp: 2026-07-14T00:00:00Z
---

# Derive-values conformance refactor — design

Archived after implementation and acceptance review.

## Context

Answers the [functional spec](spec.md). The flagship file,
`src/application/evaluation-execute.ts`, holds five function-scoped mutable
accumulators, an async `complete` closure crossing the Effect boundary once
per frame, an inline scheduler loop, a hand-rolled frontmatter scan, and a
copy of run-folder enumeration. A `let`/`.push` census over `src/` puts the
sweep at roughly twenty files, concentrated in `src/application/` and
`src/domain/evaluation/report.ts`, `plan.ts`, and `src/domain/lint/lint.ts`.

The first design draft proposed a single
`Effect.promise(() => Promise.all(...))` hashing boundary. A comparison with
current Effect and alchemy-effect guidance changed that decision: reusable
Effect operations should be named, and batch work should retain Effect's
structured failure, interruption, and ordering rather than optimize for fewer
adapter calls.

## Approach

### Module map for the decomposition

| Concern                                               | From (inline in `evaluation-execute.ts`)         | To                                                         |
| ----------------------------------------------------- | ------------------------------------------------ | ---------------------------------------------------------- |
| Evaluation/area/requirement frames, work-unit records | `evaluationFrame`, `requirementFrame`, loop body | new `src/domain/evaluation/frames.ts`                      |
| Harness request assembly                              | request literals                                 | `src/domain/evaluation/protocol.ts`                        |
| Run-artifact assembly                                 | artifact literals                                | `src/domain/evaluation/run.ts`                             |
| Ready-unit selection                                  | scheduler `for` loop                             | `readyUnits(graph, completed, limit)` in `graph.ts`        |
| Run summary stderr text, `events.jsonl` entries       | string concatenation                             | pure renderers in `src/domain/evaluation/run.ts`           |
| Run-folder name/number classification                 | `nextNumber` (×5 across commands)                | pure classifier in `src/domain/evaluation/run.ts`          |
| Directory scan feeding the classifier                 | `nextNumber` (×5)                                | one effectful helper, `src/application/evaluation-runs.ts` |

`executeHarnessRun` then reads as: resolve workspace → parse → plan → build
frames → validate sources → select ready units → build requests → write
artifacts → render result. Each step is a named call; the only remaining
mutation is none.

Reusable Effect-returning operations introduced or reshaped here use
`Effect.fn("qualitymd.<operation>")`. One-off inline composition may stay
`Effect.gen`. Terminal failures use `return yield*`; manifest read failures are
handled in the Effect error channel, while JSON syntax classification stays in
a pure function. No `try`/`catch` spans a `yield*`.

### Derivation idioms for the sweep

- Frames and work units: `plan.areas.map(...)`/`plan.requirements.map(...)`;
  the `workUnits` record derives via `Object.fromEntries` from the frame list
  plus the pending list — it is never mutated in two places.
- Effectful per-item work (source validation): `Effect.forEach`, whose result
  array is the derivation; records via `Object.fromEntries` after.
- `nextNumber`-style maxima: `Math.max(0, ...numbers) + 1` over a mapped list.
- Body guidance: `parseQualityDocument(...).body`, deleting the
  `raw.indexOf("\n---", 4)` scan.
- Protocol optionals: `buildProtocolRequest` returns optional fields; request
  assembly becomes a plain literal. Emitted JSON is unchanged because the
  conditional spreads stripped exactly these sentinels.
- Key-order discipline: object literals are constructed with the same key
  insertion order as today, since `evaluation.json` and JSON output are
  serialized in insertion order and R1 demands byte identity.

### Structured hashing traversal

`hashJson` is async because it uses Web Crypto. A named application-level
adapter wraps each call with `Effect.promise`; the fixed algorithm and valid
input make rejection a defect, matching current behavior. Frame payload hashes
derive through one sequential `Effect.forEach`, which preserves input/result
order and keeps interruption in the parent workflow. Request-id hashes use the
same adapter when each protocol request becomes available. No raw
`Promise.all` or Promise chain orchestrates workflow work.

Sequential traversal is intentional: performance is a non-goal, the payload
set is bounded by the evaluation plan, and the existing implementation is
sequential. If profiling later earns parallel hashing, the declaration can add
an explicit Effect concurrency option without changing failure semantics.

### Unified enumeration and the exclusion decision

The pure classifier implements the R6 precedence (`manifest.run.number` from
`evaluation.json`, then `run.number` from `data/evaluation-manifest.json`, then
the `NNNN` folder prefix, else not a run) and validates positive integers. The
application helper filters child directories, reads candidate manifests in
precedence order, and maps them through the classifier with `Effect.forEach`;
it uses typed Effect handling for unreadable manifests and pure parsing for
malformed JSON or missing fields. `evaluation-create`, `evaluation-execute`,
`evaluation-run`, `status`, and `evaluation-inspect` all consume it.

The `quality`-slug exclusion is removed. Provenance was checked:
`git log -S 'includes("quality")'` shows it arriving only in the
Go-to-TypeScript port commit (3206576) with no dedicated rationale, and
`evaluation-create` never applied it — so today the exclusion cannot even be
described as consistent current behavior. Its only effect is on manifest-less
folders whose slug contains `quality`; the changelog notes the change.

### Sweep method and enforcement

The sweep works from the grep-derived census, file by file: classify each
`let`/`.push` as confined-local (allowed by the guide) or shared accumulator
(refactor). Move-only commits are kept separate from edit commits so blame
survives. Each Effect-returning operation introduced or materially reshaped by
the case also receives a bounded composition review against the guide; this is
not a second whole-codebase sweep. Enforcement going forward is the guide plus
review plus `readonly` types — no new lint toolchain (see Alternatives).

## Spec response

- **R1:** byte identity is protected by constructing literals in today's key
  order and by the untouched contract suite; artifact byte comparison over a
  fixture run is the acceptance check.
- **R2:** the census is the derivation worklist; the confined-local escape hatch
  keeps the sweep from contorting genuinely local reduction loops. The touched
  Effect-operation review supplies the bounded second half.
- **R3–R4:** the module map above; every extracted function gets a
  `test/domain/evaluation/` unit test.
- **R5–R6:** one classifier, one scanner, five consumers; the enumeration
  test covers both manifest shapes, precedence, invalid/unreadable manifests,
  fallback-named and formerly excluded directories, non-directory entries, and
  unrecognized names.

## Alternatives

### Add ESLint functional/immutability rules

Rejected for now: the repo has no ESLint toolchain, and adopting one to police
a rule the guide states adds churn disproportionate to the risk. Revisit if
accumulator regressions show up in review.

### Make `hashJson` synchronous via `node:crypto`

Would eliminate the hashing boundary entirely, not just batch it. Rejected:
domain code stays on portable Web Crypto rather than Node-specific APIs. The
small named Effect adapter already makes the asynchronous boundary explicit.
May be revisited independently.

### Batch hashes with `Promise.all`

Rejected after upstream comparison: wrapping `Promise.all` once reduces adapter
calls but creates a second orchestration model without Effect's structured
interruption and failure semantics. Sequential `Effect.forEach` preserves
current behavior and is simpler to validate; explicit Effect concurrency stays
available if performance later earns it.

### Keep five enumeration copies, cleaned up locally

Rejected: the drift is the bug. Cleaning each copy preserves five places for
the next divergence.

### Keep the `quality`-slug exclusion

Rejected: no recorded reason, no spec, not applied by `evaluation-create`, and
it makes numbering depend on scope-slug vocabulary. If a real need surfaces,
it returns through a spec'd rule, not a resurrected special case.

### Decompose in place without domain moves

Rejected: helper functions inside the application file stay untestable without
Effect layers and leave the owning-layer rule unsatisfied; the moves are what
make the pure logic unit-testable.

## Trade-offs and risks

- **Byte stability under reconstruction.** Deriving objects can silently
  reorder JSON keys. Mitigated by key-order discipline and an
  artifact-byte-comparison acceptance check; any diff fails R1.
- **Additional trace frames.** Named `Effect.fn` operations can alter internal
  tracing but not public output or artifacts. Stable, qualified names make that
  internal change useful and reviewable.
- **The removed exclusion is a behavior change**, confined to manifest-less
  folders with `quality` in the slug. Accepted deliberately (R6) and noted in
  the changelog.
- **Blame churn** on moved functions. Mitigated by separating move-only from
  edit commits.
- **`readonly` ripple**: annotating exported fields can propagate through
  signatures. Bounded to the sweep list; an SDK-required mutable shape is the
  spec'd escape hatch.

## Open questions

None blocking. The deferred items (error-channel refinement, full-source
Effect-composition conformance, `test/` sweep, lint automation) are recorded in
the case's Scope.
