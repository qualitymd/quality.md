---
type: Design Doc
title: Effect runtime acceptance — design
description: Acceptance ledger, test architecture, service-boundary hardening, complete toolchain checks, native release evidence, and the v0.31.0 release sequence.
tags: [cli, typescript, effect, testing, agents, distribution, release]
timestamp: 2026-07-14T00:00:00Z
---

# Effect runtime acceptance — design

## Context

This design answers the [functional spec](spec.md) for turning 0199's completed
TypeScript implementation into a proved, idiomatic, agent-operable, published
cutover. The repository already has the target runtime, exact dependency pins,
release builders, and a green local gate. The missing work is evidence density,
boundary enforcement, active guidance, and native/distribution proof.

## Approach

### One acceptance ledger

The case folder gains `review.md`, a review-time ledger rather than a permanent
runtime artifact. It has one row per R1–R8 requirement and every 0199 cutover
gate. Each row names the authoritative evidence and remains `Pending` until that
evidence exists. Commands may be summarized there, while detailed transient
logs and probes remain under `tmp/` and are removed before landing.

The former Go suite is classified by behavior family from the pre-cutover tree:
CLI, model/document, lint, schema, status/workspace, evaluation data, reports,
runner/concurrency/harness, source resolution, evaluators, scaffold, updater,
and release tooling. A family cannot close merely because a similarly named
TypeScript module exists; its ledger entry links to contract tests or records an
obsolete behavior and the contract that removed it.

Differential acceptance uses a temporary worktree at the pre-port commit or the
v0.30.0 release binary. A throwaway driver exercises representative commands in
fresh temporary workspaces and captures status/stdout/stderr/files. The driver
is not committed. Stable fixtures discovered by the comparison become normal
TypeScript contract tests only when they protect current behavior.

### Package-local test architecture

The repository remains one package with sibling `src/` and `test/` trees. Tests
move into mirrored directories as coverage is added:

```text
test/
  cli/
  application/
  domain/
  services/
  adapters/
  integration/
  fixtures/
```

Pure functions use Vitest's ordinary `it`. Effect-returning tests use
`@effect/vitest` and return the Effect directly. Shared test Layers provide
filesystem/path, output capture, evaluator responses, clock, HTTP, and process
behavior. Real filesystem, subprocess, compiled executable, and provider
runtime checks stay visibly under `integration/` or hosted workflow steps.

The first conversion targets the current manual `Effect.runPromise` wrappers in
source, version, and provider-run tests. New coverage then follows risk rather
than deleted line count: parsing and projection; source containment; artifact
validation; graph order and concurrency; harness partial submission and resume;
evaluator failure categories; update selection; CLI output and generated trees.

### Effect boundary hardening

The existing directory architecture remains. The boundary audit uses these
rules:

- `domain` may transform values and use deterministic hashing/schema logic but
  does not read process state or external systems;
- `application` composes use cases from Effect services and pure domain logic;
- `services` define repository-owned capabilities and live/test Layers;
- `adapters` translate provider SDK, HTTP, process, and host behavior into those
  capabilities; and
- `cli` maps all failures and results to terminal output and exit codes.

The update workflow is the clearest current boundary leak because it reads
environment/process state, performs HTTP requests, and launches installers
inside the application module. It is split behind an `UpdateRuntime` service
whose live Layer owns environment, platform/libc detection, clock, HTTP, and
command execution. Pure version, archive, and readiness decisions stay exported
for ordinary tests. The application consumes the service, and tests supply a
deterministic Layer.

Provider adapters keep SDK-required async iteration and direct provider HTTP
inside `adapters/`; they already return typed Effect failures. Their resource
acquisition is reviewed for finalizers and interruption. New tests assert abort
and late-event behavior without credentials; credentialed compiled-runtime
checks remain hosted/manual release evidence.

### Toolchain and active guidance

`tsconfig.json` includes every authored `.ts` file under `src/`, `test/`, and
`scripts/`. Stable `.mjs` release helpers remain Node JavaScript; release-critical
ones receive syntax/runtime coverage from existing checks rather than a
gratuitous language conversion. `package.json` keeps code-level commands;
`mise run check` stays the sole documented full repository gate.

A new `docs/guides/effect-typescript-style.md` replaces the deleted Go-specific
style/package guides. It records only active decisions: layer ownership,
namespace imports, Effect versus pure code, services/Layers, typed failures,
scoped resources, test placement, `@effect/vitest`, integration boundaries,
generated artifacts, and targeted/full validation. `AGENTS.md` routes to it and
names generated files tersely; `CONTRIBUTING.md` carries concrete commands and
the mirrored test layout.

### Native and distribution proof

The normal CI gate remains fast and deterministic on Linux. A second native
matrix builds or consumes the standalone artifacts and runs representative CLI
commands on Linux, macOS, and Windows. Alpine supplies the musl-native check.
Cross-architecture targets that hosted runners cannot execute are still compiled
and structurally inspected; native execution evidence is required where a
runner exists and is recorded honestly where it does not.

The existing install-smoke workflow remains the post-publication channel check.
Release preparation adds local snapshot/npm build evidence and hosted workflow
evidence for:

- all eight target asset names and checksum manifest entries;
- glibc/musl launcher and installer selection;
- compiled `--version`, `spec`, `lint`, and a small evaluation lifecycle;
- interruption with durable accepted state followed by resume;
- Codex/Claude external executable overrides and cancellation, using available
  authenticated hosted runtimes without embedding provider binaries;
- npm launcher packages, managed installers, Homebrew cask metadata, update
  readiness, and repair checksum reuse.

No secret-dependent test blocks the ordinary pull-request gate. Credentialed
provider checks run at release acceptance and record the exact environment and
result without prompts, source bodies, transcripts, or secret values.

### Landing and release sequence

Implementation first reaches `In-Review` with the local ledger complete except
for hosted/published evidence. The code and Change Cases land on `main`; hosted
CI for that exact commit must pass. Release preparation then moves `Unreleased`
notes into `v0.31.0`, keeps the skill at `0.31.0` with
`>=0.31.0 <0.32.0`, runs release preflight/check, and lands as a separate clean
commit. After hosted CI passes, tag `v0.31.0` triggers the repository release
workflow. Published verification closes the final ledger rows. Both cases then
move to `Done` and archive in a post-release closeout commit, followed by a fresh
`Unreleased` section.

Because the user explicitly requested completion through archival and release,
the archive may occur immediately before the release only if all non-published
acceptance is complete and the repository convention requires completed work to
land before tagging. Any published-channel failure keeps the ledger open and is
repaired or fixed forward under the release guide rather than hidden by archive
status.

## Spec response

- **State and behavioral traceability (R1–R2):** the ledger and pre-cutover
  behavior-family inventory prevent implementation presence from standing in
  for contract proof; temporary differential checks protect current output
  without reintroducing Go.
- **Tests and architecture (R3–R4):** mirrored package-local tests,
  `@effect/vitest`, deterministic Layers, explicit integration tests, and the
  update-runtime boundary make service behavior repeatable and resource cleanup
  observable.
- **Toolchain and agents (R5–R6):** complete TypeScript inclusion plus one
  active Effect guide makes the intended path both machine-checked and easy to
  discover.
- **Distribution and release (R7–R8):** native workflows, compiled-provider
  checks, existing channel verifiers, and the ordered v0.31.0 release sequence
  bind source acceptance to the artifacts users install.

## Alternatives

### Treat the existing green gate as sufficient

Rejected because 25 tests cannot prove a cutover from the larger Go suite, and
the local Linux source run does not exercise compiled native/provider/install
contracts.

### Co-locate tests inside `src/`

Rejected because the canonical Effect repositories use package-local sibling
`test/` and `typetest/` trees. Mirroring provides proximity without mixing test
dependencies and fixtures into production modules.

### Convert every test and helper to Effect

Rejected because pure deterministic logic is clearer under ordinary Vitest.
Effect-aware tests are reserved for workflows that actually need services,
scope, time, concurrency, or typed failure behavior.

### Convert all Node `.mjs` release tooling immediately

Rejected as unrelated churn. The shipped runtime is the standalone Bun
executable; stable repository release helpers may remain Node scripts when the
full release gate exercises them.

### Put native smoke logic in repository scripts

Rejected by the repository's explicit smoke-helper rule. Existing workflow
steps and temporary probes cover native checks without adding permanent smoke
utilities.

## Trade-offs and risks

- Rebuilding coverage takes longer than polishing the existing small suite, but
  it is the only evidence that addresses behavioral loss.
- Moving tests while adding coverage can obscure blame; moves therefore happen
  by boundary as the corresponding tests are strengthened.
- Hosted arm64 and provider-runtime availability may limit native evidence. The
  ledger records compiled-only versus executed evidence instead of overstating
  support.
- `@effect/vitest` is pinned to the same Effect beta and adds another pre-stable
  dependency; exact pins and the full gate contain that risk.
- Release publication depends on configured GitHub secrets and external npm,
  Homebrew, and provider services. Preflight runs before tagging, and partial
  publication follows repair/fix-forward rules.

## Open questions

None that changes the implementation approach. The acceptance ledger may reveal
specific behavior or platform defects; those are implementation findings to fix,
not reasons to weaken the requirements.
