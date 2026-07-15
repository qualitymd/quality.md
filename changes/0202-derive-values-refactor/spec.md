---
type: Functional Specification
title: Derive-values conformance refactor — functional spec
description: Requirements for a behavior-preserving refactor to derived-value style, domain decomposition of evaluation-execute, and one spec-backed run-folder enumeration.
tags: [refactor, effect, typescript, evaluation, cli]
timestamp: 2026-07-14T00:00:00Z
---

# Derive-values conformance refactor — functional spec

Companion to the
[Derive-values conformance refactor](../0202-derive-values-refactor.md) Change
Case. This spec defines what the refactor must preserve and what source state
it must reach; the [design doc](design.md) owns the module map and idioms. The
normative style rule itself lives in
[Write Effect TypeScript](../../docs/guides/effect-typescript-style.md)
("Derive values; do not accumulate them") — this spec binds the codebase to it
but does not restate it.

The key words "MUST", "MUST NOT", "SHOULD", and "MAY" are to be interpreted as
described in BCP 14 when, and only when, they appear in all capitals.

## Background / motivation

Shared mutable accumulators make data flow temporal: to know what a collection
contains at a line, a reader replays every loop that touched it. In
`evaluation-execute.ts` five such accumulators span a 225-line workflow, and
the style hid genuine drift — five diverging copies of run-folder enumeration,
one of them applying an undocumented slug exclusion the others skip. Deriving
values as expressions makes each collection's inputs visible at its
declaration, moves pure logic into unit-testable domain functions, and leaves
exactly one implementation for behavior that must agree across commands.

Because this is a refactor, its verification anchor is preservation: public
behavior stays byte-identical, so the existing contract-test suite — not new
assertions — proves the sweep changed structure and nothing else. The one
deliberate behavior decision (R6) is spec'd rather than inherited, per "an
unspecified case is a decision delegated."

## Requirements

### R1 — Public behavior preservation

For identical inputs and injected time and randomness, every `qualitymd`
command MUST produce byte-identical public output before and after this
refactor — command grammar, stdout and stderr text, JSON shapes and values,
exit categories, persisted run artifacts, and generated files — except for the
enumeration edge case R6 decides.

> Rationale: byte-level preservation is what makes a wide mechanical sweep
> reviewable; any observable diff outside R6 is a defect, not a judgment
> call. — 0202
>
> Durable spec: none.

Existing contract and unit tests MUST pass without weakened assertions; tests
MAY move with their subjects and new tests MAY be added, but no assertion on
public behavior may be deleted or loosened.

> Durable spec: none.

### R2 — Derivation conformance

Modules under `src/` MUST conform to the style guide's "Derive values; do not
accumulate them" section: collections are built as expressions, no mutable
accumulator is shared across statements or loops, and a mutable local appears
only when confined to a few lines inside one small function.

> Rationale: the rule already binds new code via the guide; this requirement
> clears the existing stock so the guide describes the codebase rather than an
> aspiration. — 0202
>
> Durable spec: none — the normative rule lives in the guide, which already
> carries it.

Exported types with collection-valued fields MUST declare them `ReadonlyArray`
or `readonly` unless an external SDK contract requires a mutable shape.

> Durable spec: none.

### R3 — Evaluation-execute decomposition

Frame construction, work-unit derivation, harness request assembly,
run-artifact assembly, and run summary rendering currently inlined in
`src/application/evaluation-execute.ts` MUST be pure, unit-tested functions
under `src/domain/evaluation/`, leaving `executeHarnessRun` a linear
composition of named steps.

> Rationale: these are deterministic transformations of the model and plan;
> the owning-layer rule assigns them to domain, and their inline placement is
> what forced the accumulator style. — 0202
>
> Durable spec: none.

`executeHarnessRun` MUST derive body guidance from the parsed quality
document's body field rather than re-scanning the raw model text for a
frontmatter delimiter.

> Durable spec: none.

The protocol request builder MUST represent absent optional request parts as
absent fields rather than sentinel empty-string or null values, so request
assembly needs no conditional-spread stripping.

> Durable spec: none — the emitted request JSON is unchanged (R1); only the
> internal representation of absence changes.

### R4 — Ready-unit selection in the graph domain

Selection of dispatchable work units — evaluator-backed, dependencies
complete, bounded by the resolved concurrency — MUST be a pure, unit-tested
function in `src/domain/evaluation/graph.ts`.

> Rationale: scheduling eligibility is graph logic; inlined in the workflow it
> is expressed as a three-clause negated `continue` a reader must invert. —
> 0202
>
> Durable spec: none — dispatch behavior is already governed by the evaluation
> orchestration specs and does not change.

### R5 — Single run-folder enumeration

One shared implementation MUST replace the five per-command run-folder
enumeration and numbering copies in `evaluation-create`, `evaluation-execute`,
`evaluation-run`, `status`, and `evaluation-inspect`, with the pure
folder-name and manifest-number classification living in `src/domain/`.

> Rationale: the copies have already drifted — different manifest precedence
> and an exclusion applied by four of the five — and every new caller drifts
> further. — 0202
>
> Durable spec: none — the decided behavior is R6; this requirement is the
> single-implementation source state.

### R6 — Decided recognition and numbering rule

When enumerating an evaluation directory, run recognition and numbering MUST
follow one rule: a folder's run number is `run.number` from a readable
`evaluation.json`, else the manifest run number from a readable
`data/evaluation-manifest.json`, else the `NNNN` prefix of a folder named
`NNNN-<slug>-eval`, else the folder is not a run. No slug-content exclusion
applies; the undocumented exclusion of slugs containing `quality` MUST be
removed.

> Rationale: the exclusion arrived wholesale in the runtime port with no
> recorded reason and no spec backing, and it makes numbering depend on the
> words in a scope slug. Removing it is the early-alpha clean break; the rule
> is stated durably so the next implementation cannot re-diverge. — 0202
>
> Durable spec: modify `specs/cli.md` — add the shared run-folder recognition
> and numbering rule; modify `specs/cli/evaluation-create.md` — define
> next-number computation by reference to that shared rule.

## Requirement-set review

R1 fixes the verification anchor: the suite proves preservation, so the sweep
(R2) and the structural moves (R3–R5) can be reviewed as pure refactors. R6 is
the single deliberate behavior decision, carved out of R1 explicitly and given
a durable spec home so the unified implementation (R5) has one contract to
satisfy. Together they achieve the motivation — a codebase that conforms to
the guide rule, a legible flagship workflow, and enumeration that cannot
silently disagree across commands — without granting the case any other
behavior change. Each requirement is verifiable: R1 via the existing contract
suite and artifact byte comparison, R2 via source inspection over the
sweep-derived file list, R3–R5 via source inspection plus new domain unit
tests, and R6 via an enumeration unit test covering manifest, fallback, and
formerly excluded folder names.

## Durable spec changes

### To add

None.

### To modify

- `specs/cli.md` — add the shared run-folder recognition and numbering rule
  (per R6).
- `specs/cli/evaluation-create.md` — define next-number computation by
  reference to the shared rule (per R6).

### To rename

None.

### To delete

None.
