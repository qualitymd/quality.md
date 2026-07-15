---
type: Design Doc
title: Derive-values conformance refactor — design
description: Module map for the evaluation-execute decomposition, derivation idioms for the sweep, a single hashing boundary, and one run-folder enumeration.
tags: [refactor, effect, typescript, evaluation, cli]
timestamp: 2026-07-14T00:00:00Z
---

# Derive-values conformance refactor — design

## Context

Answers the [functional spec](spec.md). The flagship file,
`src/application/evaluation-execute.ts`, holds five function-scoped mutable
accumulators, an async `complete` closure crossing the Effect boundary once
per frame, an inline scheduler loop, a hand-rolled frontmatter scan, and a
copy of run-folder enumeration. A `let`/`.push` census over `src/` puts the
sweep at roughly twenty files, concentrated in `src/application/` and
`src/domain/evaluation/report.ts`, `plan.ts`, and `src/domain/lint/lint.ts`.

## Approach

### Module map for the decomposition

| Concern                                               | From (inline in `evaluation-execute.ts`)         | To                                                               |
| ----------------------------------------------------- | ------------------------------------------------ | ---------------------------------------------------------------- |
| Evaluation/area/requirement frames, work-unit records | `evaluationFrame`, `requirementFrame`, loop body | new `src/domain/evaluation/frames.ts`                            |
| Harness request assembly, run-artifact assembly       | request/artifact literals                        | `frames.ts` (or a sibling if it grows past one screenful per fn) |
| Ready-unit selection                                  | scheduler `for` loop                             | `readyUnits(graph, completed, limit)` in `graph.ts`              |
| Run summary stderr text, `events.jsonl` entries       | string concatenation                             | pure renderers in `src/domain/evaluation/run.ts`                 |
| Run-folder name/number classification                 | `nextNumber` (×5 across commands)                | pure classifier in `src/domain/evaluation/run.ts`                |
| Directory scan feeding the classifier                 | `nextNumber` (×5)                                | one effectful helper, `src/application/evaluation-runs.ts`       |

`executeHarnessRun` then reads as: resolve workspace → parse → plan → build
frames → validate sources → select ready units → build requests → write
artifacts → render result. Each step is a named call; the only remaining
mutation is none.

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

### One hashing boundary

`hashJson` is async (Web Crypto), which is what forced the per-frame
`Effect.promise(() => complete(...))` pattern. All frame payload hashes are
computed in a single `Effect.promise(() => Promise.all(payloads.map(...)))`
step; request-id hashes join the same batch. One boundary crossing replaces
about eight.

### Unified enumeration and the exclusion decision

The pure classifier implements the R6 rule (manifest number from
`evaluation.json`, else `data/evaluation-manifest.json`, else `NNNN` folder
prefix, else not a run). The application helper scans the directory and maps
names through it; `evaluation-create`, `evaluation-execute`,
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
survives. Enforcement of the style going forward is the guide plus review plus
`readonly` types — no new lint toolchain (see Alternatives).

## Spec response

- **R1:** byte identity is protected by constructing literals in today's key
  order and by the untouched contract suite; artifact byte comparison over a
  fixture run is the acceptance check.
- **R2:** the census is the worklist; the confined-local escape hatch keeps
  the sweep from contorting genuinely local reduction loops.
- **R3–R4:** the module map above; every extracted function gets a
  `test/domain/evaluation/` unit test.
- **R5–R6:** one classifier, one scanner, five consumers; the enumeration
  test covers manifest-backed, fallback-named, formerly excluded, and
  unrecognized folder names.

## Alternatives

### Add ESLint functional/immutability rules

Rejected for now: the repo has no ESLint toolchain, and adopting one to police
a rule the guide states adds churn disproportionate to the risk. Revisit if
accumulator regressions show up in review.

### Make `hashJson` synchronous via `node:crypto`

Would eliminate the hashing boundary entirely, not just batch it. Rejected:
domain code stays on portable Web Crypto rather than Node-specific APIs, and
batching captures most of the readability win. May be revisited independently.

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
- **The removed exclusion is a behavior change**, confined to manifest-less
  folders with `quality` in the slug. Accepted deliberately (R6) and noted in
  the changelog.
- **Blame churn** on moved functions. Mitigated by separating move-only from
  edit commits.
- **`readonly` ripple**: annotating exported fields can propagate through
  signatures. Bounded to the sweep list; an SDK-required mutable shape is the
  spec'd escape hatch.

## Open questions

None blocking. The deferred items (error-channel refinement, `test/` sweep,
lint automation) are recorded in the case's Scope.
