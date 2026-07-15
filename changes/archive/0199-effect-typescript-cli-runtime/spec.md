---
type: Functional Specification
title: Effect TypeScript CLI runtime
description: Requirements for replacing the Go CLI with one TypeScript runtime while preserving command, evaluation, artifact, and distribution contracts and adding bounded SDK-backed agent evaluation.
tags: [cli, evaluation, evaluator, typescript, effect, distribution]
timestamp: 2026-07-14T00:00:00Z
---

# Effect TypeScript CLI runtime

This spec governs the behavioral delta of replacing the Go implementation of
`qualitymd` with a TypeScript runtime and making Codex and Claude agent runtimes
first-class evaluators. It inherits the current **normative** contracts from the
[CLI specification](../../specs/cli.md),
[evaluation runner](../../specs/evaluation/runner.md),
[evaluator contract](../../specs/evaluation/evaluator-contract.md),
[orchestration contract](../../specs/evaluation/orchestration.md),
[`evaluation.json`](../../specs/evaluation/evaluation-json.md), and
[QUALITY.md specification](../../SPECIFICATION.md). Those contracts remain
binding unless a requirement below explicitly changes them.

The key words **MUST**, **MUST NOT**, **SHOULD**, and **MAY** are to be
interpreted as described in BCP 14 when, and only when, they appear in all
capitals.

## Background / motivation

The CLI currently combines a deterministic runner with provider-specific HTTP
and subprocess evaluators. Agent inference is qualitatively different from one
model call: it needs a bounded loop that can inspect the workspace, call tools,
and decide what context answers a prose source selector. Continuing to implement
that loop inside Go, or adding a JavaScript sidecar for provider SDKs, would
create a second orchestration and failure boundary around every evaluation.

One TypeScript runtime can use the provider agent SDKs without giving them
authority over evaluation state. The runner still detects and walks globs and
paths, freezes agent-collected source into the same bounded evidence bundle,
builds the dependency graph, schedules work, validates results, and persists
artifacts. Provider agents perform bounded collection or judgment inside that
contract. The rewrite is successful only if it gains that control without
breaking the commands and install mechanisms users already rely on.

## Scope and dependencies

The cutover covers the complete production CLI, evaluation runtime, tests,
generators, developer toolchain, packaging, and release workflow. It depends on
the current evaluator-dispatched selector contract from
[0197](../archive/0197-resolver-dispatched-source-selectors.md) and the
runner-owned concurrent/checkpoint orchestration from
[0198](../archive/0198-batched-harness-checkpoints.md).

The current Go implementation may run in development as a differential oracle
until parity gates pass. It is not a shipped fallback. No requirement below
authorizes a compatibility wrapper, sidecar, dual writer, or legacy reader.

## Requirements

### R1 — One production runtime

At cutover, the repository **MUST** build, test, generate, package, and release
the production `qualitymd` CLI without Go source, a Go toolchain, a Go-compiled
executable, or a Node sidecar owned by this project.

> Rationale: the rewrite removes the two-runtime boundary; leaving Go or a
> project-owned sidecar in the production path would retain the complexity that
> motivates the change.
>
> Durable spec: add `specs/cli/runtime-distribution.md` — define the one-runtime,
> standalone executable and no-sidecar distribution boundary.

Provider SDKs **MAY** launch their own documented native agent runtime when an
SDK-backed evaluator is selected. That child runtime **MUST** remain an evaluator
implementation detail behind the same runner-owned work contract and **MUST NOT**
be required by commands that do not select that evaluator.

> Durable spec: modify `specs/evaluation/evaluator-contract.md` — distinguish a
> provider-managed agent runtime from a project-owned sidecar and keep it inside
> the evaluator boundary.

### R2 — CLI behavioral compatibility

For every command and invocation covered by the active CLI specs, the cutover
CLI **MUST** preserve the command and flag names, defaults, argument validation,
help-visible semantics, stdout/stderr placement, human and JSON payloads,
interactive/TTY rules, next actions, and exit categories `0`, `1`, `2`, and `70`
unless this Change Case explicitly changes that behavior.

> Rationale: the implementation language is not a reason to invalidate scripts,
> agent workflows, docs, or user muscle memory.
>
> Durable spec: modify `specs/cli.md` — state that the common CLI contract is
> runtime-independent and include the cutover parity boundary.

Unexpected runtime defects **MUST NOT** expose a raw JavaScript stack trace in
ordinary output; they **MUST** map to the existing internal-error presentation
and exit category while retaining debug detail through the existing diagnostic
path.

> Durable spec: modify `specs/cli.md` — make runtime-defect normalization part of
> the common error contract.

### R3 — Model, artifact, and report continuity

The cutover CLI **MUST** accept every `QUALITY.md`, workspace configuration, and
completed or resumable evaluation artifact accepted by the immediately prior
release when its schema version is current, and **MUST** produce artifacts and
reports that conform to the same active schemas and deterministic ordering.

> Rationale: current artifacts are the current contract, not a legacy format;
> changing implementation does not earn a schema bump by itself.
>
> Durable spec: modify `specs/evaluation/evaluation-json.md` — record that a
> runtime cutover alone does not invalidate compatible run state.

When a real artifact-shape change is required during implementation, the runner
**MUST** use the existing early-alpha schema-version rule: make one explicit
clean break, document it, and refuse incompatible resume rather than add a dual
reader or migration shim.

> Durable spec: none — the existing `evaluation.json` compatibility contract
> already owns schema-version breaks.

### R4 — Graceful source-selector dispatch

For each effective area source selector, the runner **MUST** dispatch a string
that contains supported glob metacharacters and parses under the resolver's
supported glob grammar to deterministic glob collection, an existing filesystem
entry to deterministic path collection, and only a selector that is neither to
evaluator-backed inference. A string that contains glob metacharacters but is
not a valid supported glob **MUST NOT** be classified as a glob. The detected
kind remains pinned for the lifetime of the run.

> Rationale: common, explicit selectors should stay fast and reproducible while
> natural-language intent remains usable without inventing a second source
> field.
>
> Durable spec: modify `SPECIFICATION.md` and `specs/evaluation/runner.md` —
> replace metacharacter-only detection with supported-glob parsing while
> retaining path-before-prose fallback and pinned kind.

An absolute selector or one that lexically escapes the workspace **MUST** stay a
filesystem selector and fail the workspace-containment rule; it **MUST NOT** be
reinterpreted as prose merely because it is not an allowed path.

> Rationale: inference must not become a bypass around filesystem containment.
>
> Durable spec: modify `SPECIFICATION.md` and `specs/evaluation/runner.md` —
> retain unsafe filesystem intent as a filesystem error under the tightened
> classifier.

An empty glob result, unreadable existing path, or deterministic collection
failure **MUST** remain a filesystem-source failure and **MUST NOT** fall through
to inference.

> Durable spec: none — the existing source-selector divergence rule remains
> binding.

### R5 — Agent source resolution

For an inferred selector, an agent-capable evaluator **MUST** be allowed to take
multiple read-only tool actions to identify relevant workspace material, and
its accepted output **MUST** be a finite set of workspace-relative source files
rather than an assessment, rating, or free-form context transcript.

> Rationale: inference must be genuinely agentic to understand prose such as
> “the public API,” but collection and judgment must stay separate so the
> evidence boundary remains inspectable and reproducible.
>
> Durable spec: add `specs/evaluation/agent-evaluators.md` — define bounded
> source-collection sessions and their output boundary; modify
> `specs/evaluation/evaluator-contract.md` — add the agent source-resolution
> capability.

The runner **MUST** validate, bound, hash, and persist agent-selected files
through the same source-package contract as deterministic collection before any
dependent judgment can run.

> Durable spec: modify `specs/evaluation/runner.md` — make provider-agent
> collection an explicit resolver implementation under the existing capture
> boundary.

### R6 — Immutable area context

The runner **MUST** resolve one immutable context package for each in-scope area
before dispatching that area's requirement assessments. The package **MUST**
identify the captured source bundle, area evaluation frame, applicable rating
criteria and body guidance, and a content hash that changes when any of those
inputs change.

> Rationale: requirements in one area should share the same evidence and frame,
> not independently rediscover or silently widen the area.
>
> Durable spec: add `specs/evaluation/agent-evaluators.md` — define the area
> context package and its lifecycle; modify `specs/evaluation/orchestration.md`
> — add it as the dependency of local requirement work.

An assessment agent **MUST NOT** read workspace material outside its supplied
area context package. If the package is insufficient, the result **MUST** report
an evidence limitation through the evaluation schema rather than collect new
source during judgment.

> Rationale: a judge that can silently explore beyond captured evidence breaks
> resume hashes, reproducibility, and the distinction between data and
> instructions.
>
> Durable spec: add `specs/evaluation/agent-evaluators.md` — define the judgment
> evidence boundary.

### R7 — Isolated requirement assessments

Each requirement assessment **MUST** run in a fresh evaluator session or thread
whose initial context consists of the stable area context plus only that
requirement's identity, criteria, and required result schema.

> Rationale: sharing an accumulating conversation makes later requirements
> depend on model order, earlier conclusions, and context compaction; isolation
> makes parallel and sequential runs comparable.
>
> Durable spec: add `specs/evaluation/agent-evaluators.md` — define fresh
> requirement contexts; modify `specs/evaluation/orchestration.md` — require
> observational equivalence independent of thread scheduling.

Requirement assessment sessions **MUST NOT** receive the source-resolution
agent's exploration transcript or another requirement's assessment transcript.
Provider prompt caching **MAY** reuse the stable area prefix, but correctness
**MUST NOT** depend on cache availability or a provider-specific cache hit.

> Durable spec: add `specs/evaluation/agent-evaluators.md` — define transcript
> isolation and cache independence.

### R8 — Runner-owned control and bounded concurrency

The runner **MUST** remain the sole owner of the work graph, dependency
readiness, top-level concurrency, timeout, cancellation, retry budget, accepted
result validation, deterministic persistence order, and final report assembly
for every evaluator kind.

> Rationale: provider agent loops are workers inside an evaluation, not a second
> evaluation orchestrator.
>
> Durable spec: modify `specs/evaluation/orchestration.md` and
> `specs/evaluation/evaluator-contract.md` — apply runner ownership explicitly to
> SDK-backed agents.

Nested provider subagents **MUST** be disabled for requirement judgment by
default. A source-resolution profile **MAY** enable bounded nested delegation
only when the selected provider declares it and the resolved profile supplies
an explicit depth and concurrency cap.

> Rationale: nested delegation can help broad collection, but invisible fan-out
> during every requirement would defeat cost, latency, and context controls.
>
> Durable spec: add `specs/evaluation/agent-evaluators.md` — define the default
> and bounded opt-in; modify `specs/evaluation/evaluator-contract.md` — declare
> nested-agent capability.

### R9 — Honest evaluator capabilities

Every evaluator selection **MUST** expose a capability record covering at least
structured output, source resolution, tool use, concurrent calls, nested
subagents, session isolation, cancellation, usage reporting, turn limits, token
or cost limits, context-window visibility, compaction control, sandbox control,
and executable override, with unsupported controls represented as unsupported
rather than silently ignored.

> Rationale: Codex, Claude, direct APIs, and the invoking harness do not expose
> identical controls; pretending they do creates false safety and budget
> guarantees.
>
> Durable spec: modify `specs/evaluation/evaluator-contract.md` — expand the
> capability declaration; add `specs/evaluation/agent-evaluators.md` — define
> capability-to-policy resolution.

When a requested policy depends on an unsupported capability, planning **MUST**
either choose a documented safe fallback that preserves the policy or fail
before evaluator work with a classified, actionable error.

> Durable spec: modify `specs/evaluation/evaluator-contract.md` — add capability
> negotiation failure behavior.

### R10 — Evaluator selection and provider continuity

The built-in evaluator names `harness`, `codex`, `claude`, `openai`, and
`anthropic`, configured evaluator profiles, explicit selection, and standalone
`auto` discovery **MUST** remain available with their current precedence and
authentication ownership.

> Durable spec: modify `specs/evaluation/evaluator-contract.md` and
> `specs/cli/evaluation-run.md` — map existing evaluator names to the new
> runtime without changing selection semantics.

The `codex` and `claude` evaluators **MUST** use their provider's supported agent
SDK interface and authenticated local runtime. The `openai` and `anthropic`
evaluators **MUST** remain direct API evaluators that do not require the coding
agent runtime.

> Rationale: agent SDKs supply the tool loop and context lifecycle needed for
> inference; direct APIs remain useful for controlled, credential-based
> evaluation without local agent installation.
>
> Durable spec: modify `specs/evaluation/evaluator-contract.md` and add
> `specs/evaluation/agent-evaluators.md` — distinguish agent and direct-API
> evaluator classes.

Mastra **MUST NOT** be a required runtime dependency for this cutover.

> Rationale: Mastra's `OpenAISDKAgent` wraps the OpenAI Agents SDK, not the Codex
> SDK, and the existing runner already owns the deterministic orchestration that
> would otherwise justify a framework layer.
>
> Durable spec: none — this is an implementation boundary, not a public
> evaluator behavior.

### R11 — Cancellation, retry, and resume

Ctrl-C, SIGINT, SIGTERM, evaluator timeout, retryable provider failure, and
resume **MUST** preserve the current evaluation state guarantees: cancel active
provider work, persist every previously accepted result, leave a valid resumable
artifact, and never accept a late result after its work unit has been cancelled
or superseded.

> Durable spec: modify `specs/evaluation/orchestration.md` — apply cancellation
> and late-result rules to SDK streams and provider child runtimes.

Provider session or thread identifiers **MAY** be recorded as diagnostic
metadata, but replay **MUST** be determined by runner input hashes and accepted
artifacts, not by resuming provider conversation state.

> Rationale: provider histories may be pruned, unavailable, or contain tool
> transcript state that is outside the captured evidence contract.
>
> Durable spec: add `specs/evaluation/agent-evaluators.md` — make provider
> session continuity non-authoritative.

### R12 — Safety and secret boundaries

All captured source and tool output supplied to an evaluator **MUST** be treated
as untrusted data, not executable instructions, and agent evaluators **MUST** run
with the least workspace, tool, network, approval, and environment access needed
for their declared work kind.

> Durable spec: add `specs/evaluation/agent-evaluators.md` and modify
> `specs/evaluation/evaluator-contract.md` — define source-as-data and
> least-capability policy for provider agents.

Evaluator configuration **MUST** continue to store secret locators rather than
secret values. A provider child runtime **MUST** receive a filtered environment
that excludes unrelated credentials and inherits only variables required by the
selected evaluator and its documented authentication mechanism.

> Durable spec: modify `specs/evaluation/evaluator-contract.md` — extend the
> existing secret ownership boundary to provider child environments.

### R13 — Standalone platform distribution

Each release **MUST** continue to provide a directly executable `qualitymd`
archive, with `qualitymd.exe` on Windows, for Darwin arm64/x64, Linux arm64/x64,
and Windows arm64/x64 under the current six GitHub asset names, plus
`checksums.txt`.

> Rationale: changing the compiler must not break managed installers, npm
> platform packages, Homebrew, or saved release URLs.
>
> Durable spec: add `specs/cli/runtime-distribution.md` — define the supported
> target and asset-name matrix.

The standalone executable **MUST NOT** require Node.js, Bun, Go, or project
dependencies to be installed at runtime. Linux distribution **MUST** cover both
glibc and musl environments currently served by the static Go release; the
release **MAY** add libc-specific assets while preserving the existing asset
names and installer behavior.

> Durable spec: add `specs/cli/runtime-distribution.md` — define standalone and
> Linux compatibility requirements.

### R14 — Install and update continuity

The commands `npm install -g quality.md`, `npx quality.md`, the `qualitymd` and
`quality.md` npm bins, the managed shell and PowerShell installers, and the
documented Homebrew command **MUST** keep installing the platform executable
through their current ownership model.

> Durable spec: add `specs/cli/runtime-distribution.md` and modify
> `specs/cli/update.md` — preserve install owners and commands.

`qualitymd update`, `qualitymd update --check`, and ambient notices **MUST**
continue to detect and update npm, Homebrew, and managed standalone installs.
The cutover **MUST** remove Go/source as an install owner and replace every live
`go install` recovery instruction with the appropriate npm, managed-installer,
Homebrew, archive, or developer-build guidance.

> Durable spec: modify `specs/cli/update.md` — delete Go/source ownership and
> source guidance; add `specs/cli/runtime-distribution.md` — own supported
> install channels.

### R15 — Release, repair, and provenance continuity

The release workflow **MUST** continue to derive the CLI version from the tag,
embed version and commit provenance, create draft GitHub releases, publish
checksums and platform archives, publish npm launcher/platform packages, update
Homebrew, verify all channels, and repair a partially published release without
rebuilding already verified artifacts.

> Durable spec: add `specs/cli/runtime-distribution.md` — define release
> provenance and cross-channel readiness; existing release docs remain the
> procedural source.

Builds for the same tag and target **SHOULD** be reproducible enough that repair
can verify and reuse an existing artifact by checksum. Any unavoidable runtime
builder nondeterminism **MUST** be measured and documented before cutover.

> Durable spec: add `specs/cli/runtime-distribution.md` — define repair identity
> and documented nondeterminism.

### R16 — Differential parity and platform acceptance

Before cutover, the TypeScript CLI **MUST** pass a differential suite against
the final Go implementation over command help, validation errors, exit codes,
human and JSON output, generated schemas/docs, model projection, evaluation
planning, source packaging, resume/retry/cancellation, accepted evaluation
payloads, and report trees, with every intentional difference linked to an
active spec change.

> Durable spec: none — this is the migration acceptance method for existing
> contracts.

Release candidates **MUST** pass native smoke tests for every supported OS/arch
and install channel, including Ctrl-C, filesystem edge cases, provider runtime
discovery/override, structured output streaming, and Linux glibc/musl coverage.

> Durable spec: add `specs/cli/runtime-distribution.md` — define platform and
> channel conformance gates.

### R17 — TypeScript developer workflow

After cutover, a contributor **MUST** be able to install the pinned toolchain,
run the CLI from source, type-check, format, lint, test, regenerate artifacts,
build every release target, and execute the full local gate without Go or
GoReleaser.

> Durable spec: none — the contributor and build documentation own the
> developer workflow rather than a public product contract.

The project **MUST** pin the Effect v4, TypeScript, package-manager, runtime
builder, provider SDK, formatter, linter, and test-runner versions used by CI and
release builds.

> Rationale: Effect v4 and its CLI modules are pre-stable; unpinned upgrades
> would make the migration target move underneath parity work.
>
> Durable spec: none — dependency pins are implementation and repository policy.

### R18 — Skill-mediated experience

The `/quality` skill **MUST** continue to foreground the skill workflow and
`QUALITY.md`, use the CLI as its deterministic execution surface, select the
active harness when appropriate, and service harness checkpoints without asking
users to understand the implementation runtime.

> Durable spec: modify `specs/skills/quality-skill/quality-skill.md`,
> `specs/skills/quality-skill/evaluation.md`, and
> `specs/skills/quality-skill/workflows/evaluate.md` — align evaluator selection,
> inferred source collection, and checkpoint behavior with the new runtime.

The skill **MUST** explain an agent-evaluator capability or installation failure
with a concrete remedy and **MUST NOT** silently switch evaluators after a run
has accepted results.

> Durable spec: modify `specs/skills/quality-skill/evaluation.md` and
> `specs/skills/quality-skill/workflows/evaluate.md` — preserve explicit
> evaluator ownership and failure recovery.

### R19 — Diagnostics, usage, and privacy

Evaluator-call logs and run metadata **MUST** identify evaluator kind/profile,
provider model when reported, duration, attempts, classified failure, declared
capabilities, and provider-reported usage when available, while excluding raw
prompts, captured source bodies, tool transcripts, result bodies, and secrets.

> Durable spec: modify `specs/evaluation/evaluator-contract.md` — extend the
> no-raw-bodies observability contract to SDK streams and capability decisions.

The cutover CLI **MUST NOT** add project telemetry or enable provider telemetry
on the user's behalf. Provider data behavior that is inherent to a selected SDK
or authenticated runtime **MUST** be documented with the evaluator rather than
represented as `qualitymd` telemetry.

> Durable spec: add `specs/cli/runtime-distribution.md` and modify
> `specs/evaluation/evaluator-contract.md` — preserve opt-in/no-project-telemetry
> posture and disclose provider boundaries.

## Acceptance scenarios

The requirement set is complete only when these representative paths pass:

1. A user upgrades a managed standalone, npm, and Homebrew install from the last
   Go release to the TypeScript release and every current command remains
   callable under the same name.
2. A CI fixture runs the old and new CLIs on the same model and receives the same
   exit category, stdout/stderr payloads, dry-run plan, accepted evaluation
   results, and report tree.
3. An area's selector is a matching glob, an existing directory, a prose phrase,
   an empty glob, and a mistyped path-like phrase; each follows the pinned
   dispatch and failure rules without opportunistic fallback.
4. A prose selector causes a Codex or Claude source-resolution session to take
   several read-only tool actions, after which the runner captures a bounded,
   hashed file set before any requirement judgment.
5. Two requirements in one area run concurrently in fresh sessions with the same
   area context hash and produce the same persisted result ordering as
   concurrency 1; neither sees the resolver or sibling transcript.
6. Ctrl-C interrupts several active SDK-backed requirements; accepted work stays
   durable, late events are ignored, and resume schedules only incomplete or
   stale units.
7. A requested nested-agent, sandbox, turn, or cost control is unavailable in a
   provider adapter; planning reports the unsupported capability or a documented
   safe fallback instead of claiming enforcement.
8. Every release archive and npm platform package runs natively on its target;
   managed installers select a compatible Linux libc artifact and
   `qualitymd --version` reports the tag and commit.
9. A release is interrupted after GitHub assets but before npm/Homebrew publish;
   repair verifies the existing checksums, finishes missing channels, and the
   updater reports readiness only after its owning channel is live.
10. The repository's full gate, generated docs, report gallery, and distribution
    checks pass after all Go source and Go toolchain configuration are removed.

## Requirement-set review

The requirements are consistent: the runner always owns deterministic state;
provider agents act only inside source-resolution or judgment work units; shared
area evidence is immutable while requirement conversations are isolated; and
distribution compatibility is independent of implementation language. They are
complete for the stated motivation: one runtime replaces Go and the sidecar
pressure, inference gains a real agent loop, evaluator control is bounded and
honest, and existing users retain the CLI and install surfaces. Each group has a
direct verification path through differential fixtures, artifact-schema tests,
provider adapter contract tests, interruption tests, or native release-channel
smokes. No unresolved choice changes the functional outcome.

## Durable spec changes

### To add

- `specs/cli/runtime-distribution.md` — standalone runtime, supported platform
  and asset matrix, install-channel ownership, release provenance/readiness,
  platform gates, and telemetry posture (R1, R13–R15, R16, R19).
- `specs/evaluation/agent-evaluators.md` — source-resolution sessions, immutable
  area context, isolated requirement sessions, provider capability policy,
  nested-agent defaults, safety boundaries, and non-authoritative provider
  session state (R5–R12).

### To modify

- `specs/cli.md` — runtime-independent compatibility and defect normalization
  (R2).
- `specs/cli/evaluation-run.md` — preserve evaluator names and selection through
  the SDK-backed cutover (R10).
- `specs/cli/update.md` — preserve npm/Homebrew/managed owners and remove
  Go/source install ownership (R14).
- `specs/evaluation/evaluator-contract.md` — provider-managed runtimes, agent
  and API evaluator classes, expanded capabilities and negotiation, runner
  authority, environment/secret filtering, and SDK observability (R1, R5,
  R8–R12, R19).
- `specs/evaluation/runner.md` — provider-agent source resolution under the
  existing bounded capture contract and valid-glob/path/prose classification
  (R4–R5).
- `specs/evaluation/orchestration.md` — area-context dependency, isolated
  requirement sessions, SDK cancellation/late-result behavior, and runner-owned
  control (R6–R8, R11).
- `specs/evaluation/evaluation-json.md` — runtime-independent compatibility of
  current run state (R3).
- `SPECIFICATION.md` — detect a syntactically valid supported glob before an
  existing path and prose fallback, while retaining filesystem containment
  failures (R4).
- `specs/skills/quality-skill/quality-skill.md`,
  `specs/skills/quality-skill/evaluation.md`, and
  `specs/skills/quality-skill/workflows/evaluate.md` — skill-mediated evaluator
  selection, inference, checkpoint, and recovery behavior (R18).

### To rename

None

### To delete

None
