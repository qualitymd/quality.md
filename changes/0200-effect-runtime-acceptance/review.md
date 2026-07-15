---
type: Review
title: Effect runtime acceptance — review ledger
description: Requirement, cutover-gate, former-suite, differential, native, provider, and release evidence for the Effect TypeScript runtime.
tags: [cli, typescript, effect, testing, review, release]
timestamp: 2026-07-14T00:00:00Z
---

# Effect runtime acceptance — review ledger

This ledger is the authoritative acceptance record for Change Cases 0199 and 0200. `Passed` means the named evidence exists; `Pending` blocks archival and
release. Temporary worktrees and command captures under `tmp/` are removed
before landing; this file retains the reproducible command and finding.

## Requirement status

| Requirement                                  | Status  | Evidence                                                                                                                                                                               |
| -------------------------------------------- | ------- | -------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------- |
| R1 — honest state and ledger                 | Passed  | 0199 and 0200 remain `In-Review`; implementation is complete while hosted and published-channel acceptance remains open.                                                               |
| R2 — former behavior disposition             | Passed  | Suite-by-suite disposition and v0.30.0 differential results below.                                                                                                                     |
| R3 — Effect-native deterministic tests       | Passed  | Mirrored `test/{domain,application,services,integration}` tree; `@effect/vitest`; model, source, updater, provider, partial-resume, cancellation, CLI, schema, and architecture tests. |
| R4 — Effect architecture and resources       | Passed  | `HostRuntime` and `UpdateRuntime` services; domain error ownership; static architecture contract; provider temp files and sessions scoped; interruption finalizer test.                |
| R5 — complete toolchain validation           | Passed  | TypeScript includes `scripts/**/*.ts`; `mise run check` passes typecheck, lint, 13 test files/49 tests, generated-artifact checks, docs links, formatting, and npm bundle checks.      |
| R6 — agent-operable conventions              | Passed  | `docs/guides/effect-typescript-style.md`, `AGENTS.md`, and `CONTRIBUTING.md`.                                                                                                          |
| R7 — native/provider/distribution acceptance | Pending | Eight targets compile; local Darwin, Debian/glibc, Alpine/musl, npm, Codex, Claude, cancellation, and resume evidence passes. Hosted Windows and published-channel evidence remain.    |
| R8 — release and archival                    | Pending | Requires exact release-commit CI, archive, v0.31.0 publish, and channel verification.                                                                                                  |

## 0199 cutover gates

| Gate                                                        | Status  | Evidence                                                                                                                                                                   |
| ----------------------------------------------------------- | ------- | -------------------------------------------------------------------------------------------------------------------------------------------------------------------------- |
| Former Go tests have a disposition                          | Passed  | Former suite table below.                                                                                                                                                  |
| Differential tests have no unexplained differences          | Passed  | v0.30.0 differential table below; only renderer whitespace and generator-owner comment changed intentionally.                                                              |
| Every supported target passes native command/install checks | Pending | All eight targets compile; Darwin arm64 and Linux arm64 glibc/musl run locally. Hosted macOS, Linux, and Windows matrix plus post-publish results remain.                  |
| Compiled Codex and Claude adapters use external runtimes    | Passed  | Release-candidate binaries completed 12/12-unit evaluations through explicit external Codex and Claude executables; interruption finalization remains in the normal suite. |
| Linux glibc and musl coverage                               | Passed  | Debian Bookworm ran the glibc arm64 archive; Alpine 3.22 ran the musl arm64 archive; target-selection tests and hosted Alpine coverage remain active.                      |
| Ctrl-C cleans up and leaves a resumable artifact            | Passed  | `test/integration/evaluation-provider.test.ts` interrupts an in-flight evaluator, observes its finalizer, and reads the pending resumable artifact.                        |
| Generated CLI docs and report gallery are current           | Passed  | `mise run check` passed CLI-doc, schema, specification sync, and report-gallery drift checks; the v0.30.0 report tree is byte-identical.                                   |
| Archives, npm, Homebrew, updater, and repair pass           | Pending | Eight archives and eight npm platform packages build; local npm launcher install passes. Homebrew, updater, repair, and published channels await release.                  |
| No live Go production/toolchain path                        | Passed  | No `.go`, `go.mod`, Go build task, GoReleaser config, or active Go contributor guidance remains; historical records are preserved.                                         |

## Former Go suite disposition

Each row classifies every test in the named former suite. `Covered` means the
current contract is exercised at the listed TypeScript boundary. `Retired`
means the test asserted Go/Cobra/Fang rendering, build-info, or internal
representation that the clean runtime break deliberately removed; the current
public contract is still tested.

| Former suite                          | Disposition       | Current evidence or reason                                                                                                                                                                      |
| ------------------------------------- | ----------------- | ----------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------- |
| `internal/agentinstructions`          | Covered           | CLI init default/opt-out/force tests cover create, update, deduplication, and idempotence through observable files.                                                                             |
| `internal/cli/evaluation_run`         | Covered           | Harness failure/checkpoint/retry/partial-resume tests and evaluator-selection tests.                                                                                                            |
| `internal/cli/evaluation`             | Covered           | CLI command/removed-alias contracts plus schema/example/provider lifecycle checks.                                                                                                              |
| `internal/cli/init`                   | Covered           | Default, minimal, stdout, JSON, refusal, force, and agent-instruction modes in `test/integration/cli.test.ts`.                                                                                  |
| `internal/cli/lint`                   | Covered           | Valid, invalid human/JSON, usage, exit-category, deterministic repair, and idempotence contracts.                                                                                               |
| `internal/cli/model`                  | Covered           | Projection unit tests and executable tree/list/get canonical-ID contracts.                                                                                                                      |
| `internal/cli/root`                   | Covered / retired | Help, version, unknown command, stdout/stderr, and exit categories covered; Go build-info fallback internals retired for tag-injected Bun constants.                                            |
| `internal/cli/schema`                 | Covered           | Executable artifact/JSON/no-argument contracts and generated schema drift test.                                                                                                                 |
| `internal/cli/spec`                   | Covered           | Executable bundled artifact/no-argument contract; Go writer/style internals retired.                                                                                                            |
| `internal/cli/status`                 | Covered           | Valid workspace JSON snapshot plus source/evaluation state exercised by full-project check.                                                                                                     |
| `internal/cli/style`                  | Retired / covered | Fang-specific color scheme and exact padded rendering retired; stable plain help and lint contracts remain executable tests.                                                                    |
| `internal/cli/version_update`         | Covered           | Pure version/updater selection, cache, opt-out, platform/libc archive, and injected runtime tests; removed `upgrade` alias is rejected.                                                         |
| `internal/document`                   | Covered           | Frontmatter failures, body preservation, rendering, and map-entry unit tests.                                                                                                                   |
| `internal/evaluation`                 | Covered           | Bundled example/schema semantic validation, provider end-to-end report build, partial persistence, and byte-stable report-gallery drift/differential checks.                                    |
| `internal/evaluation/model_reference` | Covered           | Model projection/scope/graph contracts and evaluation data reference validation.                                                                                                                |
| `internal/evaluator`                  | Covered           | Auto/profile selection, capability metadata, provider end-to-end, structured payload validation, and interruption; compiled SDK override is a separate pending native gate.                     |
| `internal/lint`                       | Covered           | Executable valid/invalid/fix contracts and current project lint under `mise run check`.                                                                                                         |
| `internal/lint/rules`                 | Covered           | Schema-driven lint implementation is exercised by invalid frontmatter, empty-property repair, live `QUALITY.md`, schema examples, and generated schema consistency; Go table mechanics retired. |
| `internal/markdown`                   | Covered / retired | Public report Markdown is byte-identical to v0.30.0 gallery; Go writer helper representation retired.                                                                                           |
| `internal/model/model`                | Covered           | Current model decodes under CLI and full gate.                                                                                                                                                  |
| `internal/model/projection`           | Covered           | Deterministic ordering, IDs, find, and depth tests.                                                                                                                                             |
| `internal/model/reference`            | Covered           | Canonical projection, area resolution, scope, and CLI round-trip tests.                                                                                                                         |
| `internal/model/source`               | Covered           | Declared/inherited/default source unit tests.                                                                                                                                                   |
| `internal/runner/runner`              | Covered           | Graph order, provider completion, retry, same-request resume, partial accepted-result durability, and interruption/resumability.                                                                |
| `internal/runner/source`              | Covered           | Path/glob/prose detection, default/vendor walk, containment, symlink policy, bounds, and unresolved selector tests.                                                                             |
| `internal/scaffold`                   | Covered           | Default/minimal/stdout/force executable tests plus lint/schema validation.                                                                                                                      |
| `internal/schema/jsonschema`          | Covered           | `test/integration/quality-schema.test.ts` verifies deterministic generated schema and validation examples.                                                                                      |
| `internal/schema/schema`              | Covered           | Specification/schema generated consistency gate.                                                                                                                                                |
| `internal/status`                     | Covered           | Executable status snapshot plus provider-created evaluation state and workspace/source tests.                                                                                                   |
| root `schema_test.go`                 | Covered           | Generated schema drift test and `mise run schema-check`.                                                                                                                                        |

## v0.30.0 differential evidence

The comparison built tag `v0.30.0` with pinned Go 1.26.4 in a detached
`tmp/0200-go-baseline` worktree and built the current standalone executable with
`mise run build`. Both ran with `NO_COLOR=1` and
`QUALITYMD_NO_UPDATE_CHECK=1` against the same current model.

| Contract                | Result               | Finding                                                                                                                                          |
| ----------------------- | -------------------- | ------------------------------------------------------------------------------------------------------------------------------------------------ |
| Root help               | Explained difference | Same status, stderr, command tree, prose, flags, and examples. Effect renderer removes Fang padding/wrapping whitespace.                         |
| Valid lint JSON         | Byte-identical       | Same status, stdout, and stderr.                                                                                                                 |
| Companion schema        | Explained difference | Same status and schema bytes except `$comment` now names `src/domain/model/schema.ts` and `mise run schema` instead of the deleted Go generator. |
| Model factor list JSON  | Byte-identical       | Same status, stdout, and stderr.                                                                                                                 |
| Minimal scaffold stdout | Byte-identical       | Same status, stdout, and stderr.                                                                                                                 |
| Evaluation report tree  | Byte-identical       | `diff -qr` reports no differences for `examples/report-gallery/`.                                                                                |

Executable contract tests additionally fix unknown-command/flag categories,
spec/schema arguments, invalid lint output, removed aliases, evaluator checkpoint
correlation, retry, partial submission, and resume behavior in the active suite.

## Final command and release evidence

Local release-candidate evidence on 2026-07-14:

- `mise run check` passed typecheck, warning-free lint, 13 test files and 49
  tests, schema/spec/CLI-doc/report-gallery drift checks, formatting, npm bundle
  checks, and Mintlify link validation.
- `mise run snapshot` compiled all eight documented targets. Archives range
  from 23–37 MiB, below the accepted 40 MiB compressed cap; checksums and
  metadata were generated.
- The Darwin arm64 archive ran version, `spec`, `schema`, and project `lint`;
  Debian Bookworm ran the Linux arm64 glibc archive; Alpine 3.22 ran the Linux
  arm64 musl archive with its declared runtime libraries.
- `mise run npm-build -- 0.31.0` assembled all eight platform packages. Packing
  and installing the launcher plus Darwin arm64 package in an empty prefix ran
  version `0.31.0` and `schema` successfully. Shell installer syntax and
  PowerShell parsing passed.
- A compiled Codex profile using an explicit external executable completed all
  12 work units. A compiled Claude profile using an explicit external
  executable completed all 12 work units at its honestly resolved concurrency
  of 1. These probes added strict outbound provider-schema adaptation, bounded
  Claude stderr evidence, and a regression test that provider execution—not
  only dry-run—honors non-concurrent capabilities.
- The v0.30.0 differential worktree found byte-identical lint JSON, model JSON,
  minimal scaffold output, and report-gallery output; help differs only in
  intentional renderer whitespace, and the schema differs only in its current
  generator-owner comment.

Still required before `Done`: exact release-commit hosted CI and preflight,
hosted native Windows coverage, published GitHub/npm/Homebrew/updater/repair
verification, final URLs and checksums, and Change Case archival.
