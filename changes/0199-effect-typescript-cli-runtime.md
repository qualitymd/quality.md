---
type: Change Case
title: Effect TypeScript CLI runtime
description: Replace the Go implementation with one Effect v4 TypeScript runtime, add SDK-backed agent evaluators with bounded context isolation, and preserve the CLI's behavior and multi-platform distribution contracts through a clean cutover.
status: In-Review
tags: [cli, evaluation, evaluator, typescript, effect, distribution]
timestamp: 2026-07-14T00:00:00Z
---

# Effect TypeScript CLI runtime

## Motivation

The deterministic `qualitymd` runner and the agent-capable evaluators now need
to cooperate more closely. Graceful source selection is already the right
contract — interpret valid globs as globs, existing paths as paths, and
otherwise use an agent to infer and collect the intended material — but the Go
CLI can only reach Codex and Claude through hand-built subprocess adapters. That
forces `qualitymd` to reproduce provider agent behavior and makes a JavaScript
sidecar the likely next step for first-class SDK integration.

A sidecar would leave two runtimes, two process protocols, two error models, and
two packaging systems around one evaluation run. Rewriting the CLI in
TypeScript lets the deterministic runner use provider SDKs directly while
Effect supplies one typed model for I/O, resources, concurrency, cancellation,
retry, and dependency injection. It also makes context ownership explicit:
resolve an area's source once, freeze it, then assess each requirement in a
fresh provider thread with the same bounded area context rather than sharing a
growing conversation.

The rewrite must not make users relearn the product. The command tree, output
contracts, evaluation artifacts, standalone archives, npm launcher, managed
installers, Homebrew cask, and install-aware updater are valuable interfaces
independent of Go. This case therefore makes a clean implementation break while
preserving the current behavioral and distribution contracts.

## Scope

Covered:

- replace all production Go source, tests, generators, build tasks, and release
  machinery with TypeScript built on a pinned Effect v4 release;
- keep pure projection, graph, roll-up, and rendering logic deterministic while
  expressing effects, services, typed failures, cancellation, retry, and
  bounded concurrency through Effect;
- preserve the current CLI command/flag grammar, stdout/stderr split, JSON
  shapes, exit categories, TTY behavior, evaluation protocol, persisted
  artifacts, and report outputs unless an active durable spec changes them;
- retain deterministic path and glob collection, and use a source-resolution
  agent only for selectors that are neither a valid glob nor an existing path;
- integrate Codex and Claude through their TypeScript agent SDKs, retain direct
  OpenAI and Anthropic API evaluators and the harness evaluator, and describe
  provider capabilities honestly rather than normalize unsupported controls;
- resolve an immutable area context once and run each requirement assessment in
  an isolated provider context with runner-owned parallelism, budgets,
  cancellation, validation, and persistence;
- keep current GitHub release archives, npm/npx commands, managed shell and
  PowerShell installers, Homebrew installation, checksums, update detection,
  and release-repair workflows working across macOS, Linux, and Windows on
  arm64 and x64; and
- use differential testing during implementation, then delete Go and cut over
  every release channel in one release without shipping a compatibility shim,
  Go wrapper, or Node sidecar.

Deferred:

- adopting Mastra as the evaluator or runner framework; provider adapters stay
  behind a local interface so a later case can add it if its harness,
  observability, or memory features become necessary;
- changing the QUALITY.md schema, rating semantics, evaluation protocol moves,
  report information architecture, or recommendation workflow;
- adding a remote evaluation service, managed-agent control plane, or daemon;
- making nested provider subagents the default for judgment; and
- selecting a long-term stable Effect v4 CLI API after the v4 line exits beta.

Non-goals:

- preserving Go as a fallback implementation, accepting both old and new
  runtime configuration, or maintaining parallel Go and TypeScript command
  trees;
- embedding the full Codex and Claude native runtimes in each `qualitymd`
  release artifact; and
- treating prompt caching, provider session history, or nested subagents as a
  correctness boundary.

## Affected artifacts

Derived by sweeping the Go command/runtime packages, evaluator and runner
contracts, generators, install owners, release target tables, and user-facing
toolchain instructions. The functional spec owns the durable-spec delta.

- **Code and tests:** remove `cmd/qualitymd/`, `internal/`, `go.mod`, `go.sum`,
  and all remaining `.go` production/test/generator files; add the TypeScript
  application under `src/`, TypeScript tests and contract fixtures, root
  `package.json`, lockfile, strict TypeScript configuration, and local adapters
  around Effect's unstable CLI surface and provider SDKs.
- **Toolchain and generated artifacts:** replace Go, Cobra/Fang, `gofmt`,
  `go vet`, `golangci-lint`, and Go generator tasks in `mise.toml`; rewrite
  `scripts/cli-docs/` and `scripts/report-gallery/` in TypeScript; keep generated
  `mintlify/cli.mdx` and `examples/report-gallery/` byte-stable where their
  contracts do not change.
- **Evaluation implementation:** replace the runner, artifact store, source
  packaging, work graph, evaluator selection, direct API evaluators, CLI
  subprocess evaluators, harness checkpoints, status/report loading, and logging
  packages with Effect services and pure TypeScript domain modules; add Codex
  SDK and Claude Agent SDK adapters plus capability and context-policy mapping.
- **Durable specs:** add `specs/cli/runtime-distribution.md` and
  `specs/evaluation/agent-evaluators.md`; modify `specs/cli.md`,
  `specs/cli/evaluation-run.md`, `specs/cli/update.md`,
  `specs/evaluation/evaluator-contract.md`, `specs/evaluation/runner.md`,
  `specs/evaluation/orchestration.md`, `specs/evaluation/evaluation-json.md`,
  `specs/skills/quality-skill/quality-skill.md`,
  `specs/skills/quality-skill/evaluation.md`, and
  `specs/skills/quality-skill/workflows/evaluate.md`, with their bundle indexes
  and logs.
- **Format spec:** `SPECIFICATION.md` — tighten source-kind detection from any
  string containing glob metacharacters to a syntactically valid supported glob,
  then an existing path, then prose inference. `quality.schema.json` is unchanged
  because `source` remains a scalar string.
- **Bundled skill:** `skills/quality/SKILL.md`,
  `skills/quality/workflows/evaluate.md`, evaluator-selection and CLI workflow
  resources, compatibility metadata, and skill logs; other workflows should
  need only generated/help reference updates.
- **Distribution and release:** replace `.goreleaser.yaml`; update
  `.github/workflows/ci.yml`, `.github/workflows/release.yml`,
  `.github/workflows/install-smoke.yml`, `scripts/build-npm.mjs`, release and
  repair scripts, `scripts/lib/release.mjs`, `npm/quality.md/`,
  `npm/platforms/`, `install/install.sh`, `install/install.ps1`,
  `install/install.cmd`, and the Homebrew cask update path while preserving
  current install commands and release asset names. Add Linux libc-specific
  artifacts only where compatibility testing proves they are needed.
- **Durable docs:** `README.md`, `CONTRIBUTING.md`, `docs/install.md`,
  `docs/guides/cli-design.md`, `docs/guides/cut-a-release.md`,
  `docs/reference/versioning.md`, related Mintlify pages, and any live source
  install instructions; remove `go install` and document the TypeScript
  development workflow.
- **Release communication:** `CHANGELOG.md` and the release notes for the
  cutover. Historical Change Cases, archived evaluations, and prior changelog
  entries stay unchanged.

## Children

- [Functional spec](0199-effect-typescript-cli-runtime/spec.md) — clean runtime
  cutover, compatibility, source resolution and context isolation, evaluator
  control, packaging, install/update continuity, and acceptance gates.
- [Design doc](0199-effect-typescript-cli-runtime/design.md) — Effect service
  architecture, pure core boundaries, provider adapters, area/requirement
  context model, compiled executable packaging, release migration, and test
  strategy.
