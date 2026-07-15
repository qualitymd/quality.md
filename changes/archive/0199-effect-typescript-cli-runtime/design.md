---
type: Design Doc
title: Effect TypeScript CLI runtime design
description: Replace the Go CLI with a layered Effect v4 TypeScript application, keep a pure deterministic core, isolate provider agent contexts behind capability-aware adapters, and ship the same multi-platform install surfaces as standalone executables.
tags: [cli, evaluation, evaluator, typescript, effect, distribution]
timestamp: 2026-07-14T00:00:00Z
---

# Effect TypeScript CLI runtime design

## Context

This design answers the
[Effect TypeScript CLI runtime functional spec](spec.md). The present CLI is a
single Go executable: Cobra/Fang provides the command surface; Go packages own
model loading, evaluation orchestration, evaluator selection, artifact
persistence, reporting, update behavior, and generated docs; GoReleaser builds
six archives; npm and Homebrew redistribute those same executables.

The evaluator boundary is where that architecture is now expensive. The direct
API adapters can make one structured request, and the CLI adapters can invoke
`codex exec` or `claude --print`, but inferred source selection needs a real
agent loop with tool calls and context management. Adding that loop as a Node
sidecar would make Go responsible for supervising another application runtime
and translating cancellation, errors, streams, schemas, and configuration over
a custom protocol.

The chosen approach is one TypeScript application built on Effect v4. “One
runtime” means one project-owned application and orchestration model. The Codex
and Claude SDKs still launch their documented provider runtimes: the Codex
TypeScript SDK wraps and spawns Codex over JSONL, and the Claude Agent SDK
launches Claude Code. Those are replaceable evaluator workers, not a second
`qualitymd` service or sidecar protocol.

## Decision summary

- Use strict TypeScript and pin `effect@4.0.0-beta.98` and the matching Effect
  platform packages, the version published from the current Effect `main`
  branch at implementation start.
- Use Effect for I/O boundaries, services, resources, structured concurrency,
  cancellation, retry, timeouts, and typed failures.
- Keep model projection, reference resolution, graph construction, hashes,
  roll-up, sorting, and rendering as plain pure TypeScript.
- Wrap `effect/unstable/cli` behind a small local command adapter so Effect v4
  CLI churn does not spread through command handlers.
- Preserve the runner's work graph and artifact contracts. Provider SDKs execute
  bounded work units; they never schedule evaluation or write run artifacts.
- Resolve a frozen `AreaContext` once, then create one fresh provider session
  per requirement using a stable area prefix and a small requirement delta.
- Integrate `@openai/codex-sdk` and `@anthropic-ai/claude-agent-sdk` directly.
  Retain the harness, OpenAI API, and Anthropic API evaluators.
- Do not adopt Mastra in this cutover. Its new `OpenAISDKAgent` uses the OpenAI
  Agents SDK, not the Codex SDK, and the deterministic runner already supplies
  the orchestration layer this application needs.
- Build standalone executables with Bun 1.3.14. The mandatory compatibility
  spike passed; Node's single-executable application path remains the proven
  no-sidecar fallback without changing the application/service design.
- Preserve existing release asset names and install owners, and add two
  libc-specific musl Linux assets because native testing proved one Linux
  executable cannot retain current glibc and musl coverage.
- Develop against the Go implementation as a differential oracle, then remove
  Go, GoReleaser, and `go install` in one cutover release.

## Architecture

### Layer map

```text
Effect CLI adapter
    │ parses args / selects command / renders terminal result
    ▼
Application use cases
    │ lint, model, evaluation, report, update
    ▼
Pure domain core ──────────────── Effect service ports
    │ model / graph / hashes          │ fs / clock / http / process / tty
    │ roll-up / ordering / render     │ artifact store / evaluator / logging
    └──────────────────────────────────┘
                         │
                         ▼
              Evaluator adapter registry
        harness │ codex │ claude │ openai │ anthropic
                         │
                         ▼
             Runner-owned acceptance/persistence
```

Dependencies point inward. Domain modules do not import Effect, Bun, Node,
provider SDKs, or CLI types. Application modules compose pure decisions with
Effect services. Runtime modules implement those services for the standalone
executable and for tests.

### Proposed source layout

```text
src/
  main.ts                    executable entry point
  cli/
    app.ts                   local wrapper over effect/unstable/cli
    commands/                one command module per durable CLI surface
    output.ts                stdout/stderr, JSON, color, next actions
    errors.ts                CLI error and exit-category mapping
  domain/
    model/                   QUALITY.md types, projection, refs, traversal
    evaluation/              work graph, payloads, roll-up, hashes
    reports/                 pure report model and renderers
    update/                  pure channel/readiness decisions
  application/
    lint.ts
    model.ts
    evaluation/              create/run/resume/status/report use cases
    update.ts
  services/
    file-system.ts
    process.ts
    clock.ts
    terminal.ts
    http.ts
    artifact-store.ts
    evaluator.ts
    logger.ts
  runtime/
    bun/                     production Effect service layers
    test/                    in-memory/fake layers
  evaluation/
    runner/                  Effect orchestration around pure graph logic
    source/                  classifier, deterministic resolvers, packaging
    evaluators/              registry and built-in adapters
      harness.ts
      codex.ts
      claude.ts
      openai.ts
      anthropic.ts
    schemas/                 Effect Schema definitions for persisted/wire data
test/
  contract/                  CLI, artifact, evaluator, report contracts
  fixtures/                  shared Go/TypeScript differential corpus
scripts/
  build-release.ts
  build-npm.ts
  cli-docs.ts
  report-gallery.ts
```

Folder names may compress during implementation, but the three boundaries —
pure domain, Effect application/services, provider/runtime adapters — must
remain visible in imports and tests.

## Effect application model

### Version and unstable API containment

Pin the exact `effect@4.0.0-beta.<n>` version in `package.json` and the lockfile;
do not use the moving `beta` tag in CI. Effect v4 consolidates platform modules
under `effect/unstable/*`, including CLI and process facilities. All imports
from unstable CLI modules live under `src/cli/app.ts`; command modules consume a
project-owned descriptor API instead of Effect CLI types directly.

The wrapper is deliberately narrow:

```ts
export interface CliCommand<A> {
  readonly name: string
  readonly description: string
  readonly input: CliInput<A>
  readonly run: (input: A) => Effect.Effect<CommandResult, CommandError, AppServices>
}
```

It owns help formatting, unknown-command/flag normalization, version flags,
argument ordering, and Effect CLI validation-error translation. This lets the
differential suite pin exact current help and usage behavior even if an
upstream beta changes its renderer.

### Services and layers

Define project services with Effect service tags and provide production or test
layers at the executable boundary:

```ts
type AppServices =
  | FileSystem
  | ProcessRunner
  | Clock
  | Terminal
  | HttpClient
  | ArtifactStore
  | EvaluatorRegistry
  | AppLogger
```

Use scoped resources for temporary directories, open streams, signal handlers,
locks, and provider child runtimes. The root program installs SIGINT/SIGTERM
handlers, interrupts the root fiber, waits for scoped finalizers, maps the final
typed error to an exit category, and calls the runtime's exit once.

Retries are explicit schedules attached only to typed retryable errors.
Cancellation is never represented as an ordinary provider failure. Defects are
caught at the command boundary, logged to the diagnostic sink, and rendered as
the existing internal-error response.

### Typed error taxonomy

Represent expected failures as tagged values, not thrown `Error` strings:

```ts
type QualitymdError =
  | UsageError
  | ModelInvalid
  | SourceUnavailable
  | SelectorUnsupported
  | EvaluatorUnavailable
  | EvaluatorPolicyUnsupported
  | EvaluatorOutputInvalid
  | RateLimited
  | TimedOut
  | RunStateInvalid
  | FileSystemFailure
  | NetworkFailure
  | InternalFailure
```

Each tag maps in one place to the current runner failure taxonomy, CLI exit
category, human message, JSON error shape, retryability, and safe diagnostic
fields. Provider adapters translate SDK events/errors into this union before
the runner sees them.

### Schemas and serialization

Use Effect Schema as the executable definition for external and persisted
shapes: QUALITY.md frontmatter, workspace config, evaluator profiles,
`evaluation.json`, JSONL event records, evaluator work requests/results, command
JSON output, and update metadata. YAML parsing remains a separate syntax step;
the decoded value then passes through the Schema decoder so field names and
aggregated errors stay in caller vocabulary.

Pure TypeScript interfaces may be derived from schemas, but no handwritten
interface may independently redefine a persisted shape. Encoders must control
key omission, array order, and newline behavior so output stays byte-stable.

## CLI and deterministic core

Each command handler returns a `CommandResult` containing exit category,
stdout payload, stderr payload, and optional next actions. It does not write to
the terminal directly. One renderer applies JSON/human, color, TTY, `--no-input`,
and broken-pipe policies after the use case completes.

The following remain pure functions and receive explicit inputs:

- model/default projection and canonical reference parsing;
- scope selection and dependency graph construction;
- selector classification (successful parse under the supported glob grammar,
  then existing path result supplied by the caller, otherwise prose);
- source-package ordering, truncation decisions, and hashes;
- work-unit IDs and input hashes;
- rating roll-up, finding/recommendation coverage, and deterministic ranking
  validation;
- report view models and Markdown/JSON renderers; and
- update SemVer/readiness/action decisions.

Time, filesystem state, environment, process discovery, HTTP, and randomness
enter only through services. Tests can therefore hold the pure results against
the Go golden corpus without emulating an Effect runtime.

## Evaluation runner

### State machine and scheduling

Port the existing work graph without redesigning it. A pure planner creates the
ordered unit graph and hashes from the model snapshot and resolved plan. The
Effect runner repeatedly:

1. loads and schema-validates current state;
2. determines complete, stale, and dependency-ready units purely;
3. runs deterministic ready units inline;
4. dispatches ready evaluator units under a bounded Effect semaphore;
5. validates each result and commits it through the single artifact writer;
6. recomputes readiness until complete, awaiting a harness checkpoint,
   cancelled, or failed; and
7. builds reports only from accepted persisted payloads.

Worker fibers return values to one coordinator fiber. They never write the run
folder. The coordinator applies completed results in deterministic graph order;
out-of-order results wait in a bounded in-memory map. The artifact is persisted
after each accepted unit using the current temp-write/sync/rename discipline,
with a platform implementation that preserves Windows replacement behavior.

The harness adapter retains the persisted rolling-window checkpoint from 0198.
It does not create fibers or child processes: dispatch returns an `Awaiting`
decision, the coordinator persists it, and the command emits the same receipt.

### Source classification and resolution

Tighten the existing classifier at its syntax boundary:

```text
has glob metacharacters and parses successfully -> glob
else filesystem entry exists               -> path
else absolute or lexically escapes workspace -> path (containment failure)
else                                       -> prose
```

The glob parser is the same parser/grammar used by expansion; metacharacters
alone are insufficient. The classifier is a pure decision over its parse result,
path shape, and a service-provided existence result. A glob or path resolver
walks deterministically and packages files directly. A prose resolver creates a
`resolveSource` work unit. It never falls through from a failed deterministic
resolver.

The agent resolver receives:

- the prose selector and its pinned `prose` kind;
- the canonical area frame and workspace root;
- read-only collection instructions;
- a minimal tool surface for path listing, globbing, searching, and reading;
- source-package caps and prohibited paths; and
- a JSON schema whose only useful result is an ordered candidate path set with
  brief selection rationale for diagnostics.

The SDK may perform many tool turns, but the adapter returns only the candidate
paths and usage metadata. The runner canonicalizes and de-duplicates paths,
rejects escapes and non-files, applies caps, reads and hashes content itself,
marks truncation, computes the bundle hash, and persists provenance before
unblocking requirement work.

### Area context and requirement threads

After source resolution, build one immutable value per area:

```ts
interface AreaContext {
  readonly areaRef: CanonicalAreaRef
  readonly evaluationFrame: AreaEvaluationFrame
  readonly sourceBundle: CapturedSourceBundle
  readonly ratingScale: ResolvedRatingScale
  readonly bodyGuidance: string | undefined
  readonly modelSnapshotHash: Sha256
  readonly contextHash: Sha256
}
```

`contextHash` covers the canonical encodings of every preceding field. The
value is persisted or deterministically reconstructible from persisted run
state; it contains no SDK thread ID or transcript.

For every requirement, build:

```text
stable prefix = evaluation rules + AreaContext + output schema
delta         = requirement ref + statement + evidence question + criteria
```

Then start a fresh SDK session/thread and submit exactly that prompt. Freshness
is enforced by adapter construction: no thread object is pooled or shared
across work units. Stable byte ordering makes the prefix eligible for provider
prompt caching. The coordinator may run multiple sessions concurrently, capped
by resolved evaluation concurrency and any stricter evaluator limit.

The assessment tool set is empty by default because all evidence is embedded in
the area context. If a provider requires a read tool to handle a package too
large to embed, expose a project-owned, read-only `getCapturedSource` tool keyed
only by evidence reference; it reads from the frozen bundle, never from the
workspace. That is a transport optimization, not source expansion.

### Context budgets and compaction

The model's real context-window size is a provider/model property, not a knob the
runner can enlarge. A profile may choose a model and the adapter may report a
known window, but a local value that merely tells a provider its window size is
not exposed as an enforcement control.

The runner instead controls the material it supplies: source-bundle byte caps,
per-tool-result caps, stable-prefix size, requirement-delta size, and maximum
structured-result size. Planning estimates the complete prompt against a
provider-reported window when available. If it cannot fit, it fails before
judgment or uses `getCapturedSource` retrieval over the already frozen bundle;
it never silently drops evidence to satisfy a token estimate.

Fresh requirement sessions are expected to finish without compaction. If a
provider compacts anyway, the result remains valid only because it is checked
against the same output schema and captured-evidence boundary; no later
requirement inherits that compacted state. Source-resolution sessions may need
several turns, so their tool outputs, wall time, and provider-supported turn or
task budgets are capped. Their transcript and any compaction summary are
discarded after the runner accepts the selected file set.

## Evaluator architecture

### Local contract

```ts
interface EvaluatorCapabilities {
  readonly structuredOutput: boolean
  readonly sourceResolution: boolean
  readonly toolUse: boolean
  readonly concurrentCalls: boolean
  readonly nestedSubagents: boolean
  readonly isolatedSessions: boolean
  readonly cancellation: boolean
  readonly usage: boolean
  readonly maxTurns: "supported" | "unsupported"
  readonly tokenBudget: "supported" | "advisory" | "unsupported"
  readonly costBudget: "supported" | "advisory" | "unsupported"
  readonly contextWindow: "reported" | "configured" | "unknown"
  readonly compaction: "observable" | "configurable" | "opaque"
  readonly sandbox: "provider" | "host" | "unsupported"
  readonly executableOverride: boolean
}

interface Evaluator {
  readonly name: string
  readonly kind: EvaluatorKind
  readonly capabilities: EvaluatorCapabilities
  readonly readiness: Effect.Effect<Readiness, EvaluatorError>
  readonly execute: (
    request: WorkRequest,
    policy: ResolvedEvaluatorPolicy,
  ) => Effect.Effect<WorkResult, EvaluatorError>
}
```

Planning intersects requested policy with capabilities. It records the resolved
policy in the run manifest/dry run. Unsupported hard controls fail as
`evaluator_policy_unsupported`; advisory controls are labeled advisory in
diagnostics. Runtime version probes populate readiness and are cached only for
the command invocation.

### Capability matrix

The initial matrix is intentionally asymmetric:

| Control                   | Harness            | Codex SDK                  | Claude Agent SDK             | Direct APIs            |
| ------------------------- | ------------------ | -------------------------- | ---------------------------- | ---------------------- |
| Structured output         | caller contract    | per-turn schema            | `outputFormat` schema        | provider schema        |
| Workspace tools           | harness-owned      | Codex tools                | allowed/disallowed tools     | project tool loop only |
| Fresh requirement context | request envelope   | new thread                 | new query/session            | new request            |
| Cancellation              | checkpoint         | abort/child termination    | `AbortController`/close      | HTTP abort             |
| Hard max turns            | harness-owned      | unsupported by TS SDK      | `maxTurns`                   | project loop only      |
| Cost budget               | harness-owned      | unsupported by TS SDK      | `maxBudgetUsd` estimate      | project accounting     |
| Token/task budget         | harness-owned      | context/config only        | `taskBudget` (alpha)         | request limits         |
| Context/compaction        | harness-owned      | metadata; no hard turn cap | usage/config visibility      | request metadata       |
| Nested agents             | harness capability | Codex multi-agent config   | `Agent` tool/definitions     | none by default        |
| Executable override       | n/a                | `codexPathOverride`        | `pathToClaudeCodeExecutable` | n/a                    |

The implementation test suite reads this matrix as data; adapters cannot claim
a capability without a contract test. Provider additions extend the union and
matrix instead of growing `if` statements through the runner.

### Codex adapter

Use `@openai/codex-sdk`. At readiness time, find an authenticated `codex`
executable using explicit profile override then `PATH`, and pass the result as
`codexPathOverride`; do not bundle the `@openai/codex` native platform package
into `qualitymd` artifacts.

Instantiate a client with a filtered environment and config overrides. For
ordinary requirement work:

- start a new thread with the workspace root and read-only sandbox;
- set the selected model/reasoning values when configured;
- set `features.multi_agent: false`, disable network/connectors, and supply the
  structured output schema;
- consume `runStreamed()` events into bounded diagnostics and usage; and
- connect Effect interruption to the SDK/child process abort path, then reject
  all subsequent events for that request ID.

For a source-resolution profile that explicitly enables nested delegation, set
`features.multi_agent: true` and apply `agents.max_threads` and
`agents.max_depth` from the bounded policy. Because the TypeScript SDK does not
offer direct programmatic spawn control, Codex decides whether to use its agent
tool within those caps. Top-level requirement concurrency still belongs to
Effect.

The adapter never resumes a requirement thread. A reported thread ID is logged
as redacted diagnostic metadata only. Codex currently exposes no first-class
hard turn or cost limit in the TypeScript SDK, so those capabilities remain
unsupported; Effect timeout and cancellation bound wall time but are not
misrepresented as token/cost enforcement.

### Claude adapter

Use `@anthropic-ai/claude-agent-sdk` and its streaming `query()` interface.
Resolve an installed `claude` executable and pass
`pathToClaudeCodeExecutable`; omit the SDK's optional native platform packages
from release artifacts rather than embedding hundreds of megabytes per target.

For requirement work, pass a fresh query with:

- `cwd` at the workspace root, `persistSession: false`, and
  `settingSources: []` so project/user agent customization cannot widen the
  evaluation tool surface (managed restrictions may still tighten it);
- a filtered `env`, explicit model/effort, JSON-schema `outputFormat`, and
  adapter-created `AbortController`;
- assessment tools empty or limited to `getCapturedSource`, with write, shell,
  web, MCP, skill, plugin, and `Agent` tools disallowed; and
- `maxTurns`, `maxBudgetUsd`, and `taskBudget` when the resolved profile requests
  and the installed SDK/runtime reports support.

Source resolution enables only read/search/glob tools under a read-only
sandbox. Nested delegation is off because `Agent` is absent. An explicit
bounded source-resolution profile may add programmatic `AgentDefinition`s and
the `Agent` tool with a per-agent `maxTurns`; no other work kind enables them.

Consume SDK messages until the one result message, validate the structured
payload, map terminal reason and provider usage, and close the query in a scoped
finalizer. The adapter treats `maxBudgetUsd` as the SDK's client-side estimate,
not exact billing, and labels it accordingly.

### Direct API and harness adapters

Port the OpenAI and Anthropic direct evaluators with their official TypeScript
HTTP SDKs or the shared Effect HTTP service. They retain current profile names,
API-key locator semantics, base URL, structured-output validation, usage, and
retry mapping. They do not gain workspace tools or inferred source resolution
unless a future spec explicitly adds a project-owned tool loop.

The harness adapter stays an Effect-free transport value behind the same
interface. Its “execute” path yields a checkpoint request to the coordinator;
the invoking skill/harness supplies correlated envelopes on resume. This
preserves the active-session evaluator without coupling `qualitymd` to any one
agent SDK.

### Why not Mastra now

Mastra's SDK-agent surface provides useful stream composition, workflows,
observability, and subagent delegation. Its `ClaudeSDKAgent` wraps the Claude
Agent SDK. Its `OpenAISDKAgent`, however, accepts `@openai/agents` `Agent`, so it
is an OpenAI Agents SDK adapter rather than a Codex SDK adapter. Adopting Mastra
would therefore not normalize the two coding harnesses we need, and would put a
second workflow abstraction above a runner that must continue to own the graph,
hashes, artifacts, resume, and deterministic ordering.

Keep the local `Evaluator` contract narrow enough that a later Mastra adapter
can implement it. Revisit Mastra only if the product needs its persistent
harness sessions, Studio/trace surface, cross-workflow handoffs, or observational
memory outside an individual evaluation work unit.

## Safety model

Agent source resolution is read-only exploration. Requirement judgment is
normally tool-less. Provider adapters start from a deny-by-default environment:
`PATH`, home/config locations needed for provider authentication, locale, and
explicit provider variables; unrelated `*_KEY`, `*_SECRET`, and `*_TOKEN`
variables are removed. User-selected API key locators add only the named value.

Prompts delimit source bodies and tool output as untrusted evidence. File paths
are normalized under the workspace root, symlink escapes are rejected under the
same policy as the current packager, and prohibited directories remain
excluded. Neither adapter can write `.quality/` or the evaluation run folder;
only the artifact service holds that authority.

Provider network behavior follows the chosen evaluator: direct APIs require
provider HTTP; local agent authentication may require provider endpoints.
Workspace browsing does not imply general web access. Planning reports the
resolved sandbox/network capability and fails when a hard requested boundary
cannot be enforced.

## Packaging and distribution

### Runtime spike result

The 2026-07-14 spike used Effect `4.0.0-beta.98` from Effect `main` at
`8ce4795c`, with alchemy-effect `main` at `573bb646` as an additional idiom
reference. It pinned Bun `1.3.14`, `@openai/codex-sdk@0.144.4`, and
`@anthropic-ai/claude-agent-sdk@0.3.209`. The throwaway spike remains under
`tmp/0199-spike/` during implementation and is not a repository test fixture.

The compiled Darwin arm64 executable exercised Effect CLI, Schema, filesystem,
HTTP, child-process, scoped runtime, and signal facilities. It used external
authenticated Codex and Claude executables through the SDK path overrides,
received schema-constrained streamed results, and cancelled each provider with
no matching child process left behind. The provider packages installed 297 MiB
and 229 MiB native payloads respectively, while the compiled application was
63 MiB, confirming those native payloads were not embedded.

Bun 1.3.14 compiled all Darwin, Linux, and Windows arm64/x64 targets, including
Windows arm64. Native Darwin execution and Dockerized arm64 Debian glibc and
Alpine musl execution passed; Alpine requires its normal `libstdc++` and
`libgcc` runtime packages. The glibc executable does not run on musl, and the
musl executable does not run on glibc, so releases add the two musl assets and
select by libc. The fallback Node SEA path also ran Effect and both external
provider SDKs on Darwin, but it was larger and slower and is not selected.

The full-provider Darwin arm64 spike measured 63 MiB raw and 23 MiB in a gzip
archive. The completed eight-target build produced 23–37 MiB archives, with the
larger Linux and Windows runtime bases establishing the cross-target bound. A
no-op invocation averaged 79.0 ms and about 70 MiB maximum RSS, compared with
11.4 ms, 16 MiB raw, and about 24 MiB maximum RSS for the current Go release.
The accepted cutover gates are at most 100 MiB raw per executable, 40 MiB per
compressed archive, 100 ms mean no-op cold start on the Darwin arm64 reference
host, and 128 MiB maximum RSS for that invocation. These are material
regressions but keep ordinary command startup sub-100-ms while buying the one
SDK-capable runtime this case requires. Release builds enforce and record the
raw and archive size budgets; the reference-host spike records startup and RSS.

Two identical Bun inputs produced different executable checksums even after
removing the Darwin ad-hoc signature. Release repair therefore treats the first
published checksum as target identity and reuses that artifact; it never
silently rebuilds an already published target. Reproducibility remains a
measured builder limitation rather than a claimed property.

### Primary builder: Bun compile

Use pinned Bun 1.3.14 to compile `src/main.ts` into a standalone executable.
Bun documents cross-compilation for Darwin, Linux, and Windows on arm64 and x64,
including baseline x64 and Linux musl targets. Disable runtime loading of
`.env`, `bunfig.toml`, `package.json`, and `tsconfig.json` so the released
binary's behavior is not changed by ambient build-tool files. Inject version
and commit as build constants.

Build this matrix:

| Release target    | Bun target                 | Existing asset                  |
| ----------------- | -------------------------- | ------------------------------- |
| Darwin arm64      | `bun-darwin-arm64`         | `qualitymd_darwin_arm64.tar.gz` |
| Darwin x64        | `bun-darwin-x64-baseline`  | `qualitymd_darwin_amd64.tar.gz` |
| Linux arm64 glibc | `bun-linux-arm64`          | `qualitymd_linux_arm64.tar.gz`  |
| Linux x64 glibc   | `bun-linux-x64-baseline`   | `qualitymd_linux_amd64.tar.gz`  |
| Windows arm64     | `bun-windows-arm64`        | `qualitymd_windows_arm64.zip`   |
| Windows x64       | `bun-windows-x64-baseline` | `qualitymd_windows_amd64.zip`   |

Also build and test `bun-linux-arm64-musl` and
`bun-linux-x64-musl-baseline`. Keep the existing names for glibc and add
`qualitymd_linux_arm64_musl.tar.gz` and
`qualitymd_linux_amd64_musl.tar.gz`; managed installers and npm platform
selection detect libc and choose them. Alpine jobs install the system
`libstdc++` and `libgcc` packages required by Bun's musl executable. The release
cannot cut over until Alpine and representative glibc distributions pass.

The executable bundles application JavaScript and the Bun runtime, not provider
agent binaries. Selecting `codex` or `claude` discovers the user's installed
runtime and returns the existing evaluator-unavailable remedy when absent.
Provider SDK adapter code is bundled only after a dependency/license and binary
size audit confirms optional native packages were excluded.

### Mandatory packaging spike

Before the application port proceeds past its first vertical slice, compile a
throwaway executable (kept only under `tmp/`) that proves:

1. Effect v4 CLI, filesystem, HTTP, process, signals, and Schema work under each
   Bun target;
2. `@openai/codex-sdk` can use `codexPathOverride`, stream structured output,
   cancel a turn, preserve cwd/env/sandbox settings, and exclude its bundled
   native runtime;
3. the Claude Agent SDK can use `pathToClaudeCodeExecutable`, stream structured
   output, enforce tool restrictions, cancel, and exclude its optional native
   runtime;
4. Ctrl-C terminates all child processes without orphaning them on Unix and
   Windows;
5. archive size, cold start, RSS, filesystem semantics, and generated output are
   acceptable against the Go baseline; and
6. native binaries pass every OS/arch/libc runner in the release matrix.

Failure of one provider adapter is a hard packaging failure, not permission to
ship a sidecar. The fallback is Node SEA.

### Fallback builder: Node SEA

If a future pinned Bun fails a hard gate, bundle the same ESM application to one JavaScript
file and create per-target Node single-executable applications in native CI
jobs. Node SEA supports embedded ESM and assets and runs without a separately
installed Node runtime, but is marked active development and has more
module/asset constraints. Disable snapshots/code cache for cross-platform
artifacts and keep provider runtimes external. The release artifact names and
install topology remain unchanged.

The fallback changes only `runtime/builder` and production service layers. It
does not authorize a Node script distribution, runtime dependency, or sidecar.

### npm, managed installers, and Homebrew

Keep `quality.md` as the small npm launcher with both bin names and optional
per-platform packages. Replace only the executable placed under each platform
package. If separate musl packages are required, add libc-qualified optional
packages and make the launcher select by platform/arch/libc; its missing-binary
message points to the managed installer/archive instead of `go install`.

Managed installers keep their current URLs, ownership metadata, checksum
verification, staging, atomic replacement, PATH behavior, and post-install
`--version` verification. Their target resolver gains libc detection only if
the release has separate Linux assets.

Homebrew remains a cask over the standalone Darwin archives, so users do not
inherit a Node/Bun formula dependency. The updater keeps GitHub readiness for
managed installs, npm registry readiness for npm, and tap readiness for
Homebrew. Remove the Go/source detection branch and manual action.

### Release build and repair

Replace GoReleaser with `scripts/build-release.ts`:

1. validate a clean tag/version/lockfile and run the full gate;
2. compile every target with pinned builder and build constants;
3. run native smoke jobs and collect the exact executables;
4. archive under current names with normalized metadata;
5. compute sorted SHA-256 `checksums.txt`;
6. create/update the draft GitHub release;
7. assemble npm platform packages from the same executables;
8. publish npm, update Homebrew, then verify channel readiness; and
9. record artifact checksums in repair state so retry reuses verified files.

Release scripts remain Node/TypeScript tooling run by mise; they are not part of
the shipped binary. Archive construction must normalize timestamps, owner/group,
file modes, and entry order. If the compiler embeds nondeterministic bytes,
repair keys identity by the original release checksum and never rebuilds a
published target silently.

## Migration and cutover

### Vertical slices

1. **Runtime spike:** prove packaging/provider compatibility and select Bun or
   Node SEA.
2. **Foundation:** add TypeScript/Effect toolchain, schemas, service layers,
   typed errors, CLI adapter, and pure-domain ports for `version`, `schema`, and
   `lint`.
3. **Read-only commands:** port model/status/list/report loading and establish
   output/help differential fixtures.
4. **Evaluation core:** port work graph, artifact store, deterministic
   resolvers, resume/retry/cancellation, roll-up, and reporting.
5. **Evaluators:** port harness and direct APIs; add Codex/Claude SDK adapters,
   `AreaContext`, and fresh requirement threads.
6. **Generators and updater:** port report gallery, CLI docs, update/install
   detection, and release tooling.
7. **Distribution:** build native matrix, npm packages, managed installers,
   Homebrew, upgrade and repair smokes.
8. **Clean cutover:** update durable specs/docs/skill, remove every Go file and
   Go task, delete `.goreleaser.yaml`, regenerate artifacts, and cut the release.

Slices may coexist on the development branch, but no release artifact invokes
both implementations. The TypeScript command cannot shell out to Go for an
unfinished path. Until cutover, the normal released CLI remains Go; at cutover,
all commands switch together.

### Data compatibility

Port current schemas first and use the Go fixtures as immutable examples. Do not
bump `evaluation.json` merely to record JavaScript types or provider thread IDs.
If the new area-context contract needs persisted fields, derive them from
already persisted source/frame/rating data where possible. If a new field is
required for correctness, make the normal early-alpha schema break explicit in
durable specs and changelog; never read both shapes.

### Go removal sweep

The final cutover commit verifies that `rg` finds no live references to
`go install`, `go run`, `go test`, `gofmt`, `go vet`, `golangci-lint`, Cobra,
Fang, GoReleaser, `cmd/qualitymd`, or `internal/` outside historical records.
Then remove `go.mod`, `go.sum`, `.goreleaser.yaml`, Go source/test files, and Go
tool pins. Archived Change Cases and historical changelog/log entries remain
untouched.

## Verification strategy

### Test layers

- **Pure unit tests:** Vitest over parsers, projection, graph planning, hashes,
  roll-up, ordering, rendering, and update decisions.
- **Effect service tests:** deterministic TestClock/in-memory filesystem/fake
  process and HTTP layers for timeouts, retries, cancellation, atomic writes,
  broken streams, and signal cleanup. Use `@effect/vitest` only where it adds
  leak/fiber assertions; ordinary pure tests stay ordinary Vitest.
- **Schema/golden tests:** decode and re-encode every current YAML/JSON fixture,
  command JSON payload, evaluation artifact, JSONL event, and report tree.
- **Evaluator contracts:** one shared suite for readiness, capability
  negotiation, structured output, failure mapping, environment filtering,
  abort, late events, usage, and no direct artifact writes; SDK live tests run
  only in credentialed provider jobs.
- **Differential CLI tests:** execute final Go and TypeScript binaries against
  the same temporary workspaces and compare argv behavior, stdout bytes, stderr
  bytes, exit code, files, and reports with an allowlist linked to spec deltas.
- **Concurrency tests:** compare concurrency 1 and N over deterministic fake
  evaluators, random completion order, cancellation points, and resume after
  every accepted unit.
- **Native distribution tests:** each archive, npm package, managed installer,
  Homebrew cask, update check/apply, and repair path across the target matrix.

### Cutover gates

The release is blocked until:

- all active Go tests have an equivalent TypeScript contract test or an explicit
  recorded reason they were obsolete;
- differential tests have no unexplained differences;
- every supported target passes native command and install smoke tests;
- both SDK adapters pass the compiled-executable spike using external provider
  runtimes;
- glibc and musl Linux coverage is demonstrated;
- Ctrl-C leaves no child process and a valid resumable artifact;
- generated CLI docs and report gallery are current;
- artifact archives, npm packages, Homebrew, updater, and repair verification
  pass from a draft release; and
- the live source tree and docs contain no Go production/toolchain path.

## Spec response

- **Runtime and CLI parity ([R1–R3](spec.md#r1--one-production-runtime)):** one
  executable entry point, centralized output/error mapping, schema-derived
  persisted types, and differential fixtures preserve existing contracts while
  the final sweep deletes Go and the sidecar path.
- **Source and context ([R4–R7](spec.md#r4--graceful-source-selector-dispatch)):**
  the tightened classifier feeds deterministic resolvers or a bounded read-only
  agent resolver; the runner captures one hashed `AreaContext`; fresh
  requirement sessions receive only that stable prefix and their local delta.
- **Control and providers ([R8–R12](spec.md#r8--runner-owned-control-and-bounded-concurrency)):**
  Effect owns the graph, fibers, resources, cancellation, retries, and writes;
  capability-aware SDK adapters expose provider differences and enforce a
  least-capability tool/environment policy.
- **Distribution ([R13–R15](spec.md#r13--standalone-platform-distribution)):**
  Bun compile targets the current six archives, with musl variants gated by
  native tests; existing npm, managed, Homebrew, update, checksum, release, and
  repair flows reuse those artifacts.
- **Acceptance and stewardship ([R16–R19](spec.md#r16--differential-parity-and-platform-acceptance)):**
  layered tests, native install smokes, pinned tooling, skill updates, and
  redacted diagnostics carry the product through one clean cutover.

## Alternatives

### Keep Go and add a Node sidecar

Rejected. It gives immediate SDK access but makes every agent evaluation cross a
project-owned process protocol. Schema evolution, cancellation, signal routing,
child cleanup, logging, packaging, and debugging would exist twice, and the Go
runner would still need abstractions similar to the Effect services proposed
here.

### Reimplement agent loops in Go

Rejected. Tool loops, provider session/context management, structured streams,
subagents, and permission controls are provider harness behavior. Reproducing
them would track undocumented semantics and still expose fewer controls than the
supported SDKs.

### Use Mastra as the primary runner

Rejected for this cutover. Mastra adds valuable agent composition, but
`OpenAISDKAgent` is backed by the OpenAI Agents SDK rather than Codex, and its
workflow/memory layer would compete with the deterministic work graph and
artifact resume contract. A later adapter remains possible.

### Use Effect v3 and stable `@effect/cli`

Rejected. The intended destination is Effect v4 and an early-alpha project can
absorb its clean APIs now. The cost is beta churn, mitigated through exact pins
and the local unstable-CLI wrapper.

### Ship JavaScript and require Node or Bun

Rejected. It would simplify builds but break managed archive and Homebrew
semantics, add a runtime prerequisite, complicate update ownership, and weaken
offline/CI installation.

### Embed Codex and Claude native runtimes

Rejected as the default. It greatly increases each platform artifact, couples
`qualitymd` releases to provider runtime redistribution and license terms, and
duplicates runtimes users of these evaluators generally already install. SDK
path overrides preserve direct integration without that payload.

### Node SEA as the primary builder

Not chosen, retained as the hard fallback. It best matches the Codex SDK's
documented Node requirement, but Node marks SEA active development and its
module/assets/cross-build path is more involved. Bun has direct TypeScript
compilation and the complete target/libc matrix; the spike decides with evidence.

### Share one provider thread across an area

Rejected. It can reduce repeated prefix tokens but makes requirement results
order-dependent, shares earlier judgment and tool noise, and delegates context
compaction to the provider. Fresh threads plus a stable cacheable prefix retain
isolation and allow runner-owned parallelism.

## Trade-offs and risks

- **Effect v4 is beta.** Exact pins and boundary wrappers reduce churn, but
  upgrades require deliberate work. Do not upgrade Effect during parity porting
  unless required for a proven blocker.
- **Bun compatibility is not guaranteed by every provider SDK.** The Codex SDK
  documents Node 18+, while Claude documents a Bun compile caveat for its bundled
  runtime. The mandatory spike and SEA fallback make this a gate, not a hope.
- **Standalone binaries will grow.** Bundling Bun or Node is larger than the Go
  binary. Excluding provider runtimes, measuring startup/RSS, and retaining
  checksummed archives bound the cost.
- **Linux portability is the sharpest packaging risk.** Go's `CGO_ENABLED=0`
  artifact hid libc differences. Native glibc/musl tests and optional extra
  assets are mandatory.
- **Provider SDKs still use subprocesses.** The rewrite removes the
  project-owned sidecar, not provider runtimes. Scoped cleanup, path overrides,
  capability probes, and redacted diagnostics must make that boundary explicit.
- **Capabilities drift.** SDK updates may add/remove controls. Exact pins and
  executable-version probes protect a release; upgrades must update the matrix
  and contract tests together.
- **Prefix repetition costs tokens.** Fresh requirement contexts repeat area
  evidence. Provider cache eligibility and captured-source tools mitigate cost;
  correctness and isolation win over conversational reuse.
- **A clean cutover concentrates release risk.** Differential testing and
  vertical slices reduce implementation risk, while a full draft-release and
  upgrade/repair rehearsal reduces channel risk. There is no runtime fallback
  after release by design.

## Resolved spike decisions

- Bun 1.3.14 is the selected builder; the compiled external-runtime Codex and
  Claude paths passed, and Node SEA remains the no-sidecar fallback.
- Releases add distinct arm64 and x64 musl assets; libc detection chooses them.
- The accepted size, cold-start, and RSS gates are recorded in
  [Runtime spike result](#runtime-spike-result) and enforced by distribution
  checks.

## Informational references

- [Effect CLI API](https://effect-ts.github.io/effect/docs/cli) and the v4
  `effect/unstable/cli` surface (version pinned in implementation).
- [Bun standalone executable and cross-compilation](https://bun.com/docs/bundler/executables).
- [Codex TypeScript SDK](https://github.com/openai/codex/tree/main/sdk/typescript)
  — wraps the Codex CLI and exposes threads, structured output, streaming,
  environment/config controls, and executable override.
- [Claude Agent SDK TypeScript reference](https://code.claude.com/docs/en/agent-sdk/typescript)
  — tools, sessions, structured output, budgets, cancellation, subagents,
  executable override, and Bun compile caveats.
- [Mastra SDK subagents announcement](https://mastra.ai/blog/introducing-sdk-subagents)
  — confirms `OpenAISDKAgent` accepts an OpenAI Agents SDK agent.
- [Node single-executable applications](https://nodejs.org/api/single-executable-applications.html)
  — fallback builder and its active-development caveat.
