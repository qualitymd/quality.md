---
type: Change Case
title: Effect runtime acceptance
description: Turn the completed TypeScript port into a release-ready, idiomatic, agent-operable Effect codebase with traceable behavioral, native-platform, provider, and distribution evidence.
status: Done
tags: [cli, typescript, effect, testing, agents, distribution, release]
timestamp: 2026-07-14T00:00:00Z
---

# Effect runtime acceptance

## Motivation

Change Case 0199 replaced the Go CLI with an Effect v4 TypeScript runtime, but
the port and the proof of the port are different deliverables. The local gate is
green over the new implementation, while the cutover design also requires
traceability from the former Go behavior, deterministic Effect service tests,
native platform and libc checks, compiled provider-runtime checks, interruption
and resume evidence, and release-channel verification. Calling the migration
complete before those claims are evidenced would turn an implementation status
into a release-readiness claim it cannot yet support.

The new codebase also needs conventions that make its architecture legible to
future people and agents. Tests should follow the package-local Effect layout,
I/O should remain behind replaceable service boundaries, release-critical
TypeScript should be typechecked, generated artifacts should name their owning
commands, and the shortest relevant validation path should be discoverable
without reconstructing the migration.

## Scope

Covered:

- reconcile 0199's implementation status with its still-open acceptance gates;
- account for former Go test behavior as covered, intentionally retired, or
  still needing a TypeScript contract test;
- organize the package-local `test/` tree by source boundary and use
  `@effect/vitest` for Effect-returning tests while keeping pure tests ordinary
  Vitest;
- add deterministic service, cancellation, resume, evaluator, schema/golden,
  and CLI contract coverage for the cutover's high-risk paths;
- keep domain logic runtime-independent and external I/O behind Effect services
  or adapters with scoped resource ownership;
- typecheck release-critical TypeScript and keep one documented full local/CI
  gate plus narrow targeted commands;
- add concise Effect TypeScript, generated-artifact, testing, and validation
  guidance for repository agents and contributors;
- demonstrate compiled executables, provider runtime overrides, interruption,
  Linux glibc/musl selection, install channels, updater, checksums, and release
  repair against the supported matrix; and
- archive the completed runtime cases and cut, publish, and verify CLI v0.31.0.

Deferred:

- a public TypeScript library API and `typetest/` suite;
- a general-purpose architecture enforcement framework beyond focused lint or
  contract checks earned by observed boundary drift;
- replacing stable Node-based release helpers solely for runtime uniformity;
- broad performance optimization beyond enforcing the accepted binary, archive,
  startup, and memory budgets; and
- upgrading Effect, Bun, TypeScript, provider SDKs, or other exact pins as part
  of the acceptance work.

Non-goals:

- moving tests beside production files in `src/`;
- preserving Go code, a Go oracle in the production tree, or a compatibility
  wrapper after cutover;
- changing public CLI commands, evaluation semantics, model semantics, or
  persisted artifact schemas; and
- adding repository-owned smoke-test scripts or fixtures prohibited by
  `AGENTS.md`; temporary probes stay under `tmp/` or in throwaway worktrees.

## Affected artifacts

Derived from the 0199 cutover gates, the former Go test/package inventory, the
current TypeScript dependency and project boundaries, Effect-returning tests,
release workflows, generated artifacts, and contributor/agent entry points.

- **Code:** `src/domain/`, `src/application/`, `src/services/`, and
  `src/adapters/` where the boundary audit finds unmanaged external I/O,
  unscoped resources, or errors that bypass the application failure model.
- **Tests:** `test/` reorganized to mirror source boundaries; additional pure,
  Effect service, CLI contract, schema/golden, evaluator, concurrency,
  cancellation, resume, and integration coverage; temporary differential and
  native probes under `tmp/` only.
- **Toolchain:** `package.json`, `bun.lock`, `tsconfig.json` and any focused
  script TypeScript configuration, `vitest.config.ts`, `.oxfmtrc.json`,
  `mise.toml`, and release/build checks.
- **CI and release:** `.github/workflows/ci.yml`,
  `.github/workflows/release-preflight.yml`,
  `.github/workflows/install-smoke.yml`, `.github/workflows/release.yml`, and
  existing release/build/verify/repair helpers where required to exercise the
  supported native and distribution matrix.
- **Agent and contributor guidance:** `AGENTS.md`, `CONTRIBUTING.md`,
  `docs/guides/index.md`, and a new concise Effect TypeScript guide covering
  boundaries, tests, generated artifacts, and validation.
- **Existing Change Case:** 0199 parent/index status language and acceptance
  evidence; both 0199 and 0200 archive only after their requirements land.
- **Durable specs:** no planned behavioral delta. Existing
  `specs/cli/runtime-distribution.md`, `specs/evaluation/agent-evaluators.md`,
  CLI/evaluation specs, and their logs already own the contracts being proven.
  They change only if implementation review exposes a real contract correction.
- **Format spec and schema:** no semantic change planned; generated schema drift
  checks remain acceptance evidence.
- **Bundled skill:** no behavioral change planned; compatibility metadata and
  release notes remain aligned with CLI v0.31.0.
- **Install/scaffold:** no contract change planned; existing installers and
  scaffold outputs are exercised as release evidence.
- **Release communication:** `CHANGELOG.md` v0.31.0 section and compatibility
  block, plus the normal release artifacts and post-release `Unreleased` reset.

## Children

- [Functional spec](0200-effect-runtime-acceptance/spec.md) — acceptance,
  architecture, testing, toolchain, agent-operability, distribution, and release
  requirements.
- [Design doc](0200-effect-runtime-acceptance/design.md) — acceptance ledger,
  coverage recovery, Effect test and service boundaries, toolchain checks,
  native workflows, and release sequence.
