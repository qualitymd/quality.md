---
type: Functional Specification
title: Effect runtime acceptance — functional spec
description: Requirements for behavioral traceability, idiomatic Effect tests and boundaries, complete toolchain validation, agent-operable guidance, native distribution evidence, and the v0.31.0 release.
tags: [cli, typescript, effect, testing, agents, distribution, release]
timestamp: 2026-07-14T00:00:00Z
---

# Effect runtime acceptance — functional spec

Companion to the
[Effect runtime acceptance](../0200-effect-runtime-acceptance.md) Change Case.
This spec states what must be true before the Effect TypeScript cutover is
complete; the later design doc owns how the repository reaches and proves it.

The key words "MUST", "SHOULD", and "MAY" are to be interpreted as described
in BCP 14 when, and only when, they appear in all capitals.

## Background / motivation

A green suite proves only what that suite exercises. The Go-to-TypeScript port
removed a much larger implementation and test corpus while preserving public
commands, artifacts, evaluation behavior, installers, and releases. Acceptance
therefore needs explicit behavioral disposition and platform evidence, not an
inference from successful compilation. The resulting repository must also make
the intended Effect boundaries and checks easy for future agents to follow, or
the architecture will drift as soon as the migration record is archived.

## Requirements

### R1 — Honest migration state and acceptance evidence

Until every 0199 cutover gate has authoritative evidence, the active Change
Case indexes **MUST** distinguish a completed implementation port from completed
acceptance and release readiness.

> Rationale: implementation presence and release proof are separate claims;
> collapsing them hides missing coverage and platform evidence. — 0200
>
> Durable spec: none.

The acceptance work **MUST** maintain a review ledger that maps each 0199
cutover gate to its proving test, command output, workflow result, source
inspection, or explicit failed/pending state.

> Durable spec: none.

### R2 — Former behavior disposition

Every behavior represented by the removed Go test suites **MUST** be classified
as covered by a TypeScript test, intentionally retired with a reason and any
governing spec delta, or missing and therefore blocking acceptance.

> Rationale: line-count parity is not the goal, but silently dropping a tested
> contract is not a clean cutover. — 0200
>
> Durable spec: none.

Representative CLI and artifact differential checks **MUST** compare the final
Go release or pre-cutover worktree with the TypeScript executable for command
grammar, exit categories, stdout, stderr, machine output, generated files, and
report trees; each observed difference **MUST** be explained or corrected.

> Durable spec: none.

### R3 — Idiomatic and deterministic tests

Runtime tests **MUST** remain in the package-local `test/` tree and **SHOULD**
mirror the corresponding `src/` boundary; cross-boundary executable tests
**MUST** be distinguishable from unit and service tests.

> Rationale: the Effect ecosystem's package-local sibling tree keeps production
> modules clean while preserving a deterministic source-to-test mapping. — 0200
>
> Durable spec: none.

Tests returning an Effect **MUST** use the repository's Effect-aware test
integration, while pure synchronous domain tests **MUST** remain ordinary test
runner tests. Time, filesystem, HTTP, process, output, and evaluator behavior
**MUST** use deterministic test services except where the test is explicitly an
integration or native-executable check.

> Durable spec: none.

The TypeScript suite **MUST** cover the cutover's high-risk behavior: model and
schema decoding, CLI compatibility, workspace/source containment, evaluation
planning and ordering, evaluator selection and failure mapping, artifact
validation and reports, bounded concurrency, cancellation, accepted-result
durability, resume, updater selection, and generated-artifact drift.

> Durable spec: none.

### R4 — Effect architecture and resource safety

Production domain modules **MUST NOT** perform filesystem, network, process,
environment, clock, terminal, or provider SDK I/O. Application workflows
**MUST** receive those capabilities through replaceable services or adapters.

> Rationale: deterministic domain logic and injectable boundaries are what make
> the new runtime testable without rebuilding hand-written process protocols. —
> 0200
>
> Durable spec: none.

Every child process, provider session, stream, and temporary resource acquired
by an Effect workflow **MUST** have scoped cleanup, and interruption **MUST NOT**
leave an orphan or discard already accepted evaluation results.

> Durable spec: none.

Failures crossing a service or adapter boundary **MUST** enter the application's
typed failure model with a stable category and actionable message; the CLI
**MUST** remain the sole owner of exit-code and stdout/stderr mapping.

> Durable spec: none.

### R5 — Complete and unambiguous toolchain validation

The repository's full gate **MUST** typecheck every production, test,
generator, and release-critical TypeScript file; lint and formatting coverage
**MUST** include the same authored TypeScript surfaces.

> Durable spec: none.

Contributor and agent guidance **MUST** name one canonical full validation
command and the narrowest commands for pure tests, Effect tests, CLI behavior,
generated artifacts, and release changes.

> Durable spec: none.

The repository **MUST** retain exact pins for the pre-stable runtime, compiler,
Effect packages, provider SDKs, formatter, linter, and test runner through this
release.

> Durable spec: none.

### R6 — Agent-operable repository conventions

Repository guidance **MUST** state the source-layer responsibilities, the
source-to-test mapping, Effect testing conventions, generated files and their
owning commands, and phase-appropriate validation without requiring an agent to
read the archived migration design.

> Rationale: archived decisions help explain the past; active guidance must make
> the correct next action discoverable in the present. — 0200
>
> Durable spec: none.

Rules whose violation can be detected reliably **SHOULD** be enforced by the
normal gate rather than left only as prose, unless enforcement would add a
broader framework than the observed risk earns.

> Durable spec: none.

### R7 — Native, provider, and distribution acceptance

Before release, compiled executables **MUST** pass representative native checks
for supported Darwin, Windows, Linux glibc, and Linux musl families, and every
declared release target **MUST** compile into the documented asset name.

> Durable spec: none; `specs/cli/runtime-distribution.md` already owns the
> supported target and asset contract.

Compiled Codex and Claude evaluator paths **MUST** demonstrate external runtime
override/discovery, structured result streaming, cancellation, and absence of
an embedded provider executable; credential-independent adapter contracts
**MUST** remain part of the normal suite.

> Durable spec: none; `specs/evaluation/agent-evaluators.md` already owns the
> provider boundary.

The release candidate **MUST** demonstrate interruption and resume, Linux libc
selection, checksum validation, npm launcher selection, managed installers,
Homebrew metadata, updater readiness, and repair reuse of published checksums.

> Durable spec: none.

### R8 — Release and archival completion

The project **MUST** run the full local gate, generated-artifact checks,
standalone snapshot, npm package build, release preflight, and hosted CI for the
exact release commit before creating the tag.

> Durable spec: none.

The v0.31.0 release **MUST** publish and verify GitHub archives, checksums, npm
packages, Homebrew, version metadata, updater output, bundled specification, and
curated release notes through the repository's release workflow.

> Durable spec: none; the release version is already selected by the breaking
> runtime and skill-facing surface recorded under `Unreleased`.

Change Cases 0199 and 0200 **MUST** reach `Done` and move to `changes/archive/`
only after their work has landed and their acceptance evidence is complete.

> Durable spec: none.

## Requirement-set review

The requirements are consistent: behavioral disposition proves preserved
contracts, Effect-aware tests and boundaries make those contracts repeatable,
toolchain and active guidance keep them operable, and native/distribution checks
prove the shipped form rather than only source execution. The set is complete
for the motivation because it covers the gap between an implemented port and a
verified release without adding new public behavior. Each requirement has a
direct verification path through the review ledger, tests, source inspection,
generated checks, build outputs, hosted workflows, or published-channel
receipts.

## Durable spec changes

### To add

None.

### To modify

None planned. Existing runtime-distribution, agent-evaluator, CLI, evaluation,
and skill specs already own the behavior being verified; an implementation
finding that reveals a contract correction must be recorded here before the
case advances to `In-Review`.

### To rename

None.

### To delete

None.
