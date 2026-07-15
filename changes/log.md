# Changes update log

## 2026-07-15

- **Done**: Completed and archived
  [0205 - Evaluation executive summary](archive/0205-evaluation-executive-summary.md).
  All R1-R8 requirements passed the review ledger: advice ordering and context,
  payload shape and persistence, stakeholder-facing register, deterministic run
  report rendering, status/reportability coverage, schema-version-10 refusal,
  durable contracts, generated gallery output, and the complete repository gate.

- **In-Review**: Completed implementation for
  [0205 - Evaluation executive summary](archive/0205-evaluation-executive-summary.md).
  Added `summarizeEvaluation`, `EvaluationSummaryResult`, deterministic run-report
  rendering, schema-version-10 refusal, status/reportability coverage, durable
  contracts, generated gallery output, release notes, and an R1–R8 review ledger.
  The full `mise run check` gate passes 19 test files and 82 tests.

- **In-Progress**: Advanced
  [0205 - Evaluation executive summary](archive/0205-evaluation-executive-summary.md) to
  implementation after settling the R1–R8 functional spec and design.

- **Design**: Opened
  [0205 - Evaluation executive summary](archive/0205-evaluation-executive-summary.md) with
  a functional spec (R1–R8) and design doc, and set `status: Design`. Adds an
  advice-phase `summarizeEvaluation` unit and `EvaluationSummaryResult` payload
  that synthesize a stakeholder-facing executive summary from the ranked findings,
  ranked recommendations, and root roll-up, rendered deterministically into the
  run report `## Summary` in place of the roll-up rationale. Records the
  orchestration, protocol, payload-kinds, data-layout, report-tree, and
  evaluation.json durable-spec impact. No production code or durable spec changed
  in this phase. Updated the bundle [index](index.md).

- **Done**: Completed and archived
  [0204 - Transport-aware evaluator concurrency](archive/0204-transport-aware-evaluator-concurrency.md).
  All R1-R8 requirements passed the review ledger: transport capability and
  resolution matrices, selected-first initialization, bounded completion-driven
  dispatch, deterministic ordering, harness partial submission, per-result
  durability, isolated retry, cancellation and resume, honest progress,
  schema-version refusal, and the complete repository gate.

- **In-Review**: Implemented
  [0204 - Transport-aware evaluator concurrency](archive/0204-transport-aware-evaluator-concurrency.md).
  Added structured dispatch capabilities, selected-evaluator-first schema-9 run
  creation, transport-derived cap resolution, completion-driven direct worker
  scheduling, serialized per-result persistence, resumable cancellation, honest
  harness delegation guidance, durable-spec and skill rollups, and an R1-R8
  review ledger. The complete repository gate passes with 19 test files and 80
  tests.

- **In-Progress**: Advanced
  [0204 - Transport-aware evaluator concurrency](archive/0204-transport-aware-evaluator-concurrency.md)
  from `Design` after confirming the selected-evaluator-first preparation and
  completion-driven coordinator design against the current runner. Runtime,
  durable-spec, skill, and test implementation is underway.

- **Design**: Advanced
  [0204 - Transport-aware evaluator concurrency](archive/0204-transport-aware-evaluator-concurrency.md)
  from `Draft` to `Design` after R1–R8 passed individual requirement-quality
  and set-level review. Added the design doc: explicit direct/delegated dispatch
  capabilities, transport defaults, selected-evaluator-first run preparation,
  a coordinator-owned completion-driven provider pool, bounded harness
  delegation, per-result persistence, schema version 9, and deterministic
  concurrency verification. No production code or durable current-behavior
  spec changed in this phase.

- **Draft**: Opened
  [0204 - Transport-aware evaluator concurrency](archive/0204-transport-aware-evaluator-concurrency.md)
  with a functional spec for runner scheduling authority, honest transport
  capabilities, transport-derived automatic caps, selected-evaluator-first run
  creation, completion-driven direct dispatch, bounded harness delegation,
  per-result durability, observability, and verification. No production code
  or durable current-behavior spec changed.

- **Done**: Completed and archived
  [0203 - Evaluator prompt cache efficiency](archive/0203-evaluator-prompt-cache-efficiency.md).
  All R1–R6 requirements passed the review ledger: exact prompt-part tests,
  provider usage and call-log coverage, the full local gate, unchanged fresh
  sessions, and two repeated scoped Claude runs with cache-read evidence.

- **In-Review**: Implemented
  [0203 - Evaluator prompt cache efficiency](archive/0203-evaluator-prompt-cache-efficiency.md).
  Added explicit canonical shared prompt prefixes, work-unit suffixes, Claude's
  supported cache-stable system preset, separate cache-read/cache-creation
  usage mapping, run-log coverage, durable evaluation contracts, and a review
  ledger. The full gate and two identical scoped Claude runs pass; the repeated
  synthesis calls reported 31,019 cache-read tokens after zero on the first
  run, without session resume or fork context.

- **In-Progress**: Advanced
  [0203 - Evaluator prompt cache efficiency](archive/0203-evaluator-prompt-cache-efficiency.md)
  from `Design` to `In-Progress` after settling prompt-part boundaries,
  canonical rendering, usage mapping, Claude system-prefix shaping, isolation,
  and verification. Implementation may now update the affected runtime, tests,
  durable evaluation specs, and release notes.

- **Design**: Advanced
  [0203 - Evaluator prompt cache efficiency](archive/0203-evaluator-prompt-cache-efficiency.md)
  from `Draft` to `Design` after R1–R6 passed individual requirement-quality
  and set-level review. Added the design doc: explicit prompt parts with a stable
  boundary, canonical JSON blocks, pure provider usage mapping, Claude's preset
  dynamic-section split, unchanged fresh-session isolation, and repeated scoped
  live evidence. No production code changed in this phase.

- **Draft**: Opened
  [0203 - Evaluator prompt cache efficiency](archive/0203-evaluator-prompt-cache-efficiency.md)
  with a functional spec for layered deterministic prompt prefixes, canonical
  structured blocks, separate cache-read/cache-write telemetry, Claude's
  cache-stable preset system prefix, fresh-session independence, and bounded
  repeated live verification. No production code changed.

## 2026-07-14

- **Done**: Completed and archived
  [0202 - Derive-values conformance refactor](archive/0202-derive-values-refactor.md).
  All R1–R6 requirements passed the review ledger: the full local gate is green,
  deterministic pre/post checkpoint artifacts are byte-identical, pure domain
  functions have focused coverage, and all run-history consumers share the
  durable recognition rule.

- **In-Review**: Implemented
  [0202 - Derive-values conformance refactor](archive/0202-derive-values-refactor.md).
  Replaced shared collection accumulators across `src/`, decomposed harness
  execution into pure frame/protocol/artifact/scheduling functions, kept hashing
  in ordered Effect traversal, unified run enumeration, updated durable CLI
  specs and release notes, and added preservation and enumeration regression
  tests. `mise run check` passes.

- **In-Progress**: Advanced
  [0202 - Derive-values conformance refactor](0202-derive-values-refactor.md)
  from `Design` to `In-Progress` after settling the derivation, Effect
  composition, domain-decomposition, structured-hashing, shared-enumeration,
  and behavior-preservation design. Implementation may now update runtime code,
  tests, durable CLI specs, and release notes.

- **Design refinement**: Compared the Effect TypeScript guide and
  [0202 - Derive-values conformance refactor](0202-derive-values-refactor.md)
  with current Effect (`8ce4795c`) and alchemy-effect (`573bb646`) practice.
  Added project rules for named `Effect.fn` operations, explicit terminal and
  typed-error handling, structured batch/repetition, service lifetimes, and
  deterministic-time tests. Revised 0202 so changed Effect operations conform,
  hashing uses ordered `Effect.forEach` instead of a nested `Promise.all`, and
  the module map assigns requests and run artifacts to their existing domain
  owners. No production code changed.

- **Design**: Advanced
  [0202 - Derive-values conformance refactor](0202-derive-values-refactor.md)
  from `Draft` to `Design` after reviewing R1–R6 individually and as a set.
  Added the design doc: a domain module map for the `evaluation-execute`
  decomposition, derivation idioms for the sweep, a single batched hashing
  boundary, one run-folder enumeration classifier shared by five commands, and
  the decision to drop the unspec'd `quality`-slug exclusion. No production
  code changed in this phase.

- **Draft**: Opened
  [0202 - Derive-values conformance refactor](0202-derive-values-refactor.md)
  with a functional spec: byte-identical public behavior (R1), `src/`
  conformance with the guide's new "Derive values; do not accumulate them"
  section (R2), pure-domain decomposition of `evaluation-execute` (R3),
  ready-unit selection in the graph domain (R4), one shared run-folder
  enumeration
  (R5), and a spec-backed recognition/numbering rule replacing five diverging
  copies (R6). The style-guide section itself landed independently ahead of
  this case. Also removed a leftover empty `0201-*` folder from `changes/`.

- **Done**: Completed and archived
  [0201 - Agent-native evidence discovery](archive/0201-agent-native-evidence-discovery.md).
  All R1–R10 requirements passed the review ledger: the full local gate is
  green, harness checkpoint and evidence-retry tests pass, and live Codex and
  Claude evaluations discovered requirement-specific evaluated/supporting
  evidence without obeying a hostile repository instruction. The clean
  schema-version-8 and SDK-only evaluator contract is ready for release.

- **In-Review**: Implemented
  [0201 - Agent-native evidence discovery](0201-agent-native-evidence-discovery.md)
  and set it `In-Review`. The graph now opens fresh requirement inspection
  sessions instead of source-resolution units; the runner seals atomic
  per-requirement evidence manifests in schema-version-8 artifacts; only
  harness, Codex, and Claude evaluators remain. Durable specs, the runtime
  skill, generated docs, and tests are synchronized. The review ledger records
  51 green tests plus successful live Codex and Claude runs over supporting
  evidence and a hostile repository-instruction fixture.

- **In-Progress**: Advanced
  [0201 - Agent-native evidence discovery](0201-agent-native-evidence-discovery.md)
  from `Design` to `In-Progress` after settling the requirement-inspection,
  evidence-sealing, SDK-isolation, evaluator-selection, artifact, resume, and
  clean-cutover design. Implementation may now update runtime code, tests,
  durable specs, the bundled skill, generated artifacts, and user guidance.

- **Design**: Advanced
  [0201 - Agent-native evidence discovery](0201-agent-native-evidence-discovery.md)
  from `Draft` to `Design` after reviewing R1–R10 individually and as a set.
  Added the design doc: a graph without `resolveSource`; fresh requirement
  inspection sessions; neutral, read-only Codex and Claude SDK policies;
  runner-sealed per-requirement evidence manifests; tools-off downstream
  synthesis; schema-version-8 persistence; and a clean removal of direct API
  evaluators and API-specific configuration. No production code changed in this
  phase.

- **Draft**: Opened
  [0201 - Agent-native evidence discovery](0201-agent-native-evidence-discovery.md)
  to move requirement-specific context discovery from deterministic area source
  packages into coding-agent evaluators while retaining runner-owned safety,
  validation, persistence, resume, and reporting. The functional spec defines
  `source` as the evaluated subject, separates authentication from evaluator
  method, keeps `harness`, `codex`, and `claude`, and removes raw OpenAI and
  Anthropic API evaluators. No production code changed in this phase.

- **Done**: Published v0.31.0 and archived
  [0199 - Effect TypeScript CLI runtime](archive/0199-effect-typescript-cli-runtime.md)
  and [0200 - Effect runtime acceptance](archive/0200-effect-runtime-acceptance.md).
  Exact release-commit CI, preflight, the release workflow, independent public
  channel verification, and the full hosted install matrix passed. The first
  install smoke exposed a Windows cmd-shim compatibility issue; the fix and its
  exact-commit CI passed, followed by a green managed-installer, npm, pnpm, and
  Homebrew rerun on macOS, Linux, and Windows. Both cases reached `Done` with
  the completed acceptance ledger archived under 0200.

- **In-Review**: Implemented
  [0200 - Effect runtime acceptance](archive/0200-effect-runtime-acceptance.md) and set
  it `In-Review`. Added mirrored Effect/Vitest coverage, deterministic host and
  updater services, typed expected failures, scoped evaluator resources,
  provider-schema adaptation, honest Claude concurrency, static architecture
  enforcement, full TypeScript checking, native CI coverage, and active Effect
  TypeScript guidance. The review ledger now accounts for every former Go suite
  and records v0.30.0 differential, eight-target snapshot/npm, Darwin,
  glibc/musl, npm install, interruption/resume, and full external Codex and
  Claude evidence. Local acceptance is green; hosted release-commit and
  published-channel evidence remain before `Done` and archival.

- **In-Progress**: Advanced
  [0200 - Effect runtime acceptance](archive/0200-effect-runtime-acceptance.md) from
  `Design` to `In-Progress` after settling its acceptance-ledger, test,
  architecture, toolchain, native-validation, and release design. Implementation
  may now update code, tests, workflows, and active guidance.

- **Design**: Advanced
  [0200 - Effect runtime acceptance](archive/0200-effect-runtime-acceptance.md) from
  `Draft` to `Design` after reviewing its R1–R8 requirements individually and
  as a set. Added the design doc: one acceptance ledger and former-behavior
  inventory; mirrored package-local tests with `@effect/vitest`; deterministic
  Layers and an `UpdateRuntime` boundary for the clearest application I/O leak;
  complete TypeScript checking and active Effect/agent guidance; native,
  provider, install, updater, and repair evidence; and an ordered v0.31.0 landing
  and publication sequence. No production code changed in this phase.

- **Draft**: Opened
  [0200 - Effect runtime acceptance](archive/0200-effect-runtime-acceptance.md) to close
  the proof gap between 0199's completed TypeScript port and a release-ready
  cutover. The functional spec requires former Go behavior disposition,
  representative differential checks, idiomatic package-local Effect tests,
  deterministic services and scoped resources, complete TypeScript toolchain
  coverage, active agent guidance, native/provider/distribution evidence, and a
  verified v0.31.0 release before 0199 and 0200 archive. No production code
  changed in this phase.

- **In-Review**: Implemented
  [0199 - Effect TypeScript CLI runtime](archive/0199-effect-typescript-cli-runtime.md)
  and set it `In-Review`. Replaced all Go source, tests, generators, toolchain,
  and GoReleaser configuration with one layered Effect v4 TypeScript runtime;
  ported the complete command and evaluation surfaces; added SDK-backed Codex
  and Claude evaluators, direct API and harness adapters, immutable area
  context, fresh requirement sessions, capability/readiness records, typed
  cancellation/retry, bounded source capture, and classified failures. Replaced
  schema, CLI-doc, and report-gallery generators with TypeScript, retained
  byte-stable report gallery output, and moved builds, installers, npm,
  Homebrew metadata, updater, repair, and CI/release workflows to the eight Bun
  targets spanning Darwin, Windows, Linux glibc, and Linux musl on arm64/x64.
  Added durable runtime/distribution and agent-evaluator specs, synchronized the
  project model, skill, docs, and release notes, and verified TypeScript checks,
  real CLI subprocess tests, native smoke behavior, and cross-target builds.

- **In-Progress**: Advanced
  [0199 - Effect TypeScript CLI runtime](archive/0199-effect-typescript-cli-runtime.md)
  from `Design` to `In-Progress` after completing the mandatory runtime spike.
  Pinned the implementation target to Effect `4.0.0-beta.98` from the current
  Effect `main` branch and Bun `1.3.14`; verified the Effect CLI, Schema,
  filesystem, HTTP, process, and signal surfaces in a compiled executable;
  verified schema-constrained streaming and cancellation through external Codex
  `0.144.4` and Claude Agent SDK `0.3.209` runtimes with no orphaned provider
  process; and compiled the full Darwin, Linux, and Windows arm64/x64 matrix.
  The spike selected Bun over the working Node SEA fallback, selected separate
  glibc and musl Linux assets because neither executable crosses libc families,
  and accepted measured caps of 100 MiB executable size, 40 MiB compressed
  archive size, 100 ms no-op cold start, and 128 MiB maximum RSS. The Darwin
  arm64 spike measured 63 MiB raw, 23 MiB compressed, 79.0 ms mean cold start,
  and about 70 MiB maximum RSS versus the Go release's 16 MiB, 11.4 ms, and
  about 24 MiB. Repeated Bun builds were not byte-identical, so release repair
  must preserve the originally published checksum rather than rebuild an
  existing target. Cloned Effect at `8ce4795c` and alchemy-effect at `573bb646`
  under `tmp/0199-reference/` for implementation reference. No production code
  changed before this phase advance.

- **Done**: Landed
  [0198 - Batched harness checkpoints](archive/0198-batched-harness-checkpoints.md),
  set it `Done`, and moved the case parent and child folder into
  [`archive/`](archive/index.md). Updated the archive index, removed the case
  from the bundle-root index, and retargeted the open 0199 dependency link to
  the archived case.

- **Design**: Opened
  [0199 - Effect TypeScript CLI runtime](archive/0199-effect-typescript-cli-runtime.md)
  and advanced it from `Draft` to `Design` after reviewing its functional
  [spec](archive/0199-effect-typescript-cli-runtime/spec.md) against the requirement
  quality and set-level bars. The case specifies one clean Go-to-TypeScript
  cutover: preserve the command, artifact, evaluation, install, update, and
  release contracts; detect a syntactically valid glob, then an existing path,
  then fall back to a bounded source-resolution agent; freeze one area context and
  assess each requirement in a fresh provider session; keep the runner in sole
  control of scheduling, budgets, cancellation, validation, and persistence;
  and expose provider capabilities without pretending Codex, Claude, direct
  APIs, and harness checkpoints have identical controls. Added the design doc:
  a pure TypeScript domain core inside an Effect v4 application, direct Codex
  and Claude Agent SDK adapters using external authenticated provider runtimes,
  Bun standalone executables behind a mandatory compatibility spike with Node
  SEA as the no-sidecar fallback, differential Go/TypeScript acceptance, and one
  release cutover that removes Go, GoReleaser, and `go install`. The three
  implementation measurements — Bun/Codex compatibility, Linux glibc/musl
  artifact selection, and acceptable binary size/startup — must be resolved by
  the runtime spike before `In-Progress`; no production code changed.

## 2026-07-11

- **In-Review**: Implemented
  [0198 - Batched harness checkpoints](0198-batched-harness-checkpoints.md)
  and set it `In-Review`. The runner's harness path now runs an
  apply/top-up/emit loop over the concurrent scheduler's ready frontier:
  `evaluation.json` bumps to schema 7 with plural `pendingEvaluatorCalls`,
  awaiting receipts carry the outstanding `evaluatorRequests` set (each
  request complete, with `lastFailure` on retries), `--evaluator-result`
  accepts one envelope or a JSON array covering any subset (valid members
  accepted independently and durably; unmatched members rejected without
  touching them; not-submitted members stay outstanding at no retry cost),
  the harness evaluator declares subagent delegation, and
  `resolveConcurrency` no longer clamps harness runs to 1. Verified with new
  runner tests (window cap, rolling top-up, partial/mixed/duplicate
  submission, concurrency-parity of accepted results, single-request
  degeneration at concurrency 1), the full gate, and a real-CLI end-to-end
  harness run at window 3. Durable specs, the skill runtime docs, generated
  CLI docs, and `CHANGELOG.md` (Unreleased) were brought into sync; the
  Affected-artifacts index was reconciled — `specs/evaluation/protocol.md`,
  `internal/runner/concurrent.go`, `internal/runner/dryrun.go`,
  `internal/evaluator/evaluator.go`, and `internal/status/` deliberately
  needed no edit, and
  `specs/skills/quality-skill/workflows/evaluate.md` was added by the
  implementation sweep.

- **In-Progress**: Advanced
  [0198 - Batched harness checkpoints](0198-batched-harness-checkpoints.md) from
  `Design` to `In-Progress` to begin implementation of the rolling checkpoint
  window per the settled [spec](0198-batched-harness-checkpoints/spec.md) and
  [design doc](0198-batched-harness-checkpoints/design.md). Updated the
  bundle-root [index](index.md).

- **Design**: Advanced
  [0198 - Batched harness checkpoints](0198-batched-harness-checkpoints.md) from
  `Draft` to `Design` after a pass on the functional
  [spec](0198-batched-harness-checkpoints/spec.md) (tightened R1-R7 to the
  requirement quality bar, split independent-durability from parity, parked usage
  aggregation under Open questions) and drafted the
  [design doc](0198-batched-harness-checkpoints/design.md). On review the
  transport model was settled as a **rolling window** (streaming) rather than a
  per-checkpoint barrier: pending state becomes a set (artifact schema `6 → 7`),
  the harness path reuses the concurrent scheduler's `nextReadyEvaluationStep`
  frontier to keep up to `concurrency` requests outstanding — topping up as
  members are accepted — and `--evaluator-result` accepts one or more envelopes
  per resume call, bound per member by `(RequestID, InputHash)` through the
  existing validation/persistence path. Partial submission is the normal path, so
  not-submitted members stay outstanding at no retry cost while failed members
  re-emit and consume retry budget. Alternatives (barrier-per-checkpoint, a
  persistent JSON-RPC/JSONL transport, a separate harness cap,
  singular-plus-side-table, not-submitted-as-failed) recorded as rejected or
  deferred. No code touched; case work stays within its folder.

- **Design**: Resolved the open questions on
  [0198 - Batched harness checkpoints](0198-batched-harness-checkpoints.md).
  Usage aggregation needs no change — usage is already logged per evaluator call
  with no run total, so each accepted member preserves per-member usage logging
  and parity extends to it; a run-total rollup is a separate observability
  change. Default concurrency keeps 0195's shared default (no harness-specific
  default) with resolved concurrency and window width surfaced in receipts and
  dry-run. The only remaining open item is a soft implementation question — an
  optional runner-internal receipt-size clamp below the concurrency cap, left
  unclamped until a fixture shows it matters. The case is ready to advance to
  `In-Progress`.

- **Draft**: Opened
  [0198 - Batched harness checkpoints](0198-batched-harness-checkpoints.md) to
  design the multi-checkpoint harness work deferred by
  [0194](archive/0194-harness-native-evaluator-dispatch.md) and
  [0195](archive/0195-evaluation-concurrency.md). The case proposes emitting a
  bounded batch of dependency-ready judgment requests per `awaiting_evaluator`
  checkpoint — bounded by the run's resolved concurrency and the work-graph DAG —
  and accepting a batch of correlated result envelopes on resume, so the
  invoking harness can fan independent requests out to subagents while the runner
  keeps sole ownership of scheduling, validation, and persistence. Added the
  parent concept, child [index](0198-batched-harness-checkpoints/index.md), and a
  functional-spec [sketch](0198-batched-harness-checkpoints/spec.md) (R1-R7 with
  per-requirement durable-spec impact), and listed it in the bundle-root
  [index](index.md).

- **Done**: Landed
  [0197 - Resolver-dispatched source selectors](archive/0197-resolver-dispatched-source-selectors.md)
  in the v0.30.0 release, set it `Done`, and moved the case parent and its
  child folder into [`archive/`](archive/index.md). Updated the archive
  [index](archive/index.md) and cleared the entry from the bundle-root
  [index](index.md).

- **In-Review**: Implemented
  [0197 - Resolver-dispatched source selectors](0197-resolver-dispatched-source-selectors.md)
  and advanced it from `In-Progress` to `In-Review`. The runner now detects
  each effective selector's kind (glob metacharacters → glob, existing
  filesystem entry → path, otherwise prose) and pins it at run creation in
  `evaluation.json`'s new per-area `sources` provenance record (artifact
  schema version 6); prose selectors resolve through evaluator-backed
  `resolveSource` work units on the existing harness checkpoint transport,
  whose returned files are validated, capped, hashed, and captured into the
  bounded bundle — persisted atomically with unit completion — before any
  dependent judgment (R1–R3, R6); a selector kind the selected evaluator
  cannot serve fails the run at plan time with the new `selector_unsupported`
  category, distinct from `source_unavailable`, with a per-unit guard as
  backstop (R4); and dry-run previews plus run receipts surface the per-area
  source dispatch plan while the provenance record carries resolver and
  harness-runtime attribution (R5). Durable specs (`SPECIFICATION.md`, the
  evaluation [runner](../specs/evaluation/runner.md),
  [protocol](../specs/evaluation/protocol.md),
  [orchestration](../specs/evaluation/orchestration.md),
  [evaluator contract](../specs/evaluation/evaluator-contract.md),
  [evaluation.json](../specs/evaluation/evaluation-json.md),
  [evaluation run](../specs/cli/evaluation-run.md), and the
  [skill evaluation workflow](../specs/skills/quality-skill/evaluation.md)),
  the bundled skill's evaluate workflow, the regenerated Mintlify
  specification page, and release notes were updated; the affected-artifacts
  list was reconciled (orchestration and evaluation-run specs added;
  `quality.schema.json` confirmed unchanged). Verified with the full local
  gate (`vet`, `lint`, tests, docs generators and link check) and an
  end-to-end CLI smoke run: prose dry-run plan, `resolveSource` checkpoint,
  captured provenance, and the dependent judgment request carrying the
  captured bundle. Updated the bundle [index](index.md).

- **In-Progress**: Advanced
  [0197 - Resolver-dispatched source selectors](0197-resolver-dispatched-source-selectors.md)
  from `Design` to `In-Progress`. The functional spec (R1–R6) and design doc
  are settled; implementation begins: runner-owned kind detection pinned at
  run creation in the per-area `sources` record, `resolveSource` work units
  riding the harness checkpoint transport with capture-before-dependents, the
  `SourceResolution` capability with the plan-time `selector_unsupported`
  failure, and the durable spec updates. Updated the bundle [index](index.md).

- **Design**: Advanced
  [0197 - Resolver-dispatched source selectors](0197-resolver-dispatched-source-selectors.md)
  from `Draft` to `Design`. Settled the two gating questions in the
  [functional spec](0197-resolver-dispatched-source-selectors/spec.md):
  `SPECIFICATION.md` commits to non-filesystem selectors (Q2), and a selector
  stays a bare string with detected kind — glob metacharacters → glob,
  existing filesystem entry → path, otherwise prose — so the frontmatter shape
  and `quality.schema.json` do not move (Q1). Added the
  [design doc](0197-resolver-dispatched-source-selectors/design.md): kind
  detection is a runner-owned classifier in front of the untouched 0196
  packaging path, pinned at run creation in a new per-area `sources` record in
  `evaluation.json`; prose resolution rides the existing harness checkpoint
  transport as first-class `resolveSource` work units that capture returned
  material (caps, hashes, provenance) before dependent judgment; a
  `SourceResolution` evaluator capability plus a plan-time
  `selector_unsupported` failure keeps unsupported selectors loud, early, and
  distinct from `source_unavailable`. Recorded the accepted Q1 typo hazard and
  its mitigations as design trade-offs. Reconciled the parent's affected
  artifacts (skill resolution step confirmed; `internal/model` and the
  companion schema deliberately unchanged) and updated the bundle
  [index](index.md). No code changed (gated to In-Progress).

- **Draft**: Opened
  [0197 - Resolver-dispatched source selectors](0197-resolver-dispatched-source-selectors.md)
  (`status: Draft`), the follow-up 0196 deferred: resolve every effective
  source selector through a per-kind resolver feeding the bounded, hashed
  evidence bundle, dispatch non-deterministic kinds to the invoking harness
  through the existing checkpoint transport, fail unsupported kinds with a
  classification distinct from `source_unavailable`, and record bundle
  provenance in the run artifact. Added the parent concept and child folder
  with the [functional spec](0197-resolver-dispatched-source-selectors/spec.md)
  (R1–R6, gated by two open questions on selector typing and format
  commitment), and listed it in the bundle [index](index.md). Motivated by
  0196's [considerations](archive/0196-spec-faithful-model-reading/considerations.md)
  and the 2026-07-11 re-validation, which confirmed the bundle contract's
  guarantees (input-hash guard, evidence-bound re-judgment, resumability) as
  the seam resolvers must feed and surfaced default-selector granularity as
  design pressure.

- **Done**: Implemented and archived
  [0196 - Spec-faithful model reading](archive/0196-spec-faithful-model-reading.md).
  The shared `model.EffectiveSource` resolver now serves both the evaluation
  runner and `status`, so a source-less root area packages the document's
  directory and source-less child areas inherit the nearest ancestor's selector
  (R1–R2); source packaging expands glob selectors through the same sorted,
  hashed, capped bundling via a small in-tree matcher (R3), skips symlinked and
  other non-regular entries instead of crashing (R5), and fails work whose
  selector packages zero readable files with a classified `source_unavailable`
  naming the selector across the sequential, concurrent, and harness dispatch
  paths (R4). `qualitymd lint` reports extension frontmatter as the new
  warning-severity `unknown-key` advisory instead of an `invalid-frontmatter`
  error, so conforming extension documents lint valid and load (R6). The
  companion `quality.schema.json` accepts any non-empty scalar for content
  scalars and drops the inaccurate ordering-enforcement `$comment` claim
  (R7–R8). Durable specs ([runner](../specs/evaluation/runner.md),
  [lint rules](../specs/cli/lint-rules.md),
  [quality.schema.json](../specs/quality-schema-json.md)), the regenerated
  schema artifact, and release notes were updated; the bundled skill needed no
  change. Verified with `mise run check`. Moved the case into
  [`archive/`](archive/index.md) and updated the bundle [index](index.md).

- **In-Progress**: Advanced
  [0196 - Spec-faithful model reading](archive/0196-spec-faithful-model-reading.md) from
  `Design` to `In-Progress`. The functional spec (R1–R8) and design doc are
  settled; implementation of the shared `EffectiveSource` resolver, glob and
  symlink-safe source packaging with the `source_unavailable` guard, the
  warning-severity `unknown-key` lint reclassification, the companion-schema
  scalar widening, and the durable spec updates begins. Updated the bundle
  [index](index.md).

- **Design**: Advanced
  [0196 - Spec-faithful model reading](archive/0196-spec-faithful-model-reading.md) from
  `Draft` to `Design` and added its
  [design doc](archive/0196-spec-faithful-model-reading/design.md). The design extracts a
  shared `effectiveSource` resolver so `runner` and `status` cannot diverge on
  root-default and ancestor-inherited source (R1–R2), adds glob and symlink-safe
  packaging plus a `source_unavailable` guard against silent empty evidence
  (R3–R5), reclassifies unknown frontmatter keys from an `invalid-frontmatter`
  error to a warning-severity `unknown-key` advisory so spec-permitted extensions
  lint valid (R6), and widens the companion schema's content scalars beyond
  `string` while correcting an inaccurate ordering-enforcement claim (R7–R8).
  Resolved the four design questions (D1–D4): `**` via a small in-tree matcher
  with no new dependency; globs honor the skip list except a literal-segment
  opt-in (`vendor/**`); binary-only sources are `source_unavailable`; and the
  shared `EffectiveSource` resolver lives in `internal/model`. No code changed
  (gated to In-Progress).

- **Draft**: Opened
  [0196 - Spec-faithful model reading](archive/0196-spec-faithful-model-reading.md)
  (`status: Draft`). The case corrects a family of places where the CLI narrows
  an abstract `SPECIFICATION.md` concept to a closed assumption: the runner
  treats `source` as a single present literal path (source-less root/child areas
  get empty evidence, globs are not expanded, unresolved selectors are judged
  against nothing, and the directory walk crashes on committed symlinked skill
  directories), `lint` rejects spec-permitted extension frontmatter as invalid,
  and the companion JSON schema forces scalars to strings and overclaims an
  ordering check. Diagnosed from duplicated failed evaluation runs in two
  external projects. Added the parent concept, child folder with
  [functional spec](archive/0196-spec-faithful-model-reading/spec.md) (R1–R8) and a
  [considerations](archive/0196-spec-faithful-model-reading/considerations.md) sketch on
  resolver-dispatched selectors, and listed it in the bundle
  [index](index.md). `SPECIFICATION.md` is the conformance target and is
  intentionally unchanged.

## 2026-07-10

- **Done**: Implemented and archived
  [0195 - Evaluation concurrency](archive/0195-evaluation-concurrency.md).
  The runner now resolves a single `evaluation.concurrency` setting (omitted
  default `max(2, runtime.NumCPU()*2)`) instead of the public
  `executionStrategy` enum, records only resolved `concurrency` in dry-run JSON,
  run receipts, logs, and schema-version-5 `evaluation.json`, resolves
  harness-backed runs to `concurrency: 1`, and runs dependency-ready
  evaluator-backed steps concurrently through a coordinator-owned scheduler that
  keeps artifact mutation and deterministic result ordering inside the runner.
  Durable evaluation/CLI specs, runtime skill guidance, and release notes were
  updated. Verified with `mise run check` and `go test -race ./internal/runner`.
  Moved the case into [`archive/`](archive/index.md) and updated the bundle
  [index](index.md).

- **In-Progress**: Advanced
  [0195 - Evaluation concurrency](archive/0195-evaluation-concurrency.md) from
  `Design` to `In-Progress`. The functional spec and design doc are settled;
  implementation of `evaluation.concurrency`, concurrency resolution, the
  coordinator-owned runner scheduler, and durable spec/skill updates begins.
  Updated the bundle [index](index.md).

- **Design**: Created
  [0195 - Evaluation concurrency](archive/0195-evaluation-concurrency.md) with
  its [functional spec](archive/0195-evaluation-concurrency/spec.md) and
  [design doc](archive/0195-evaluation-concurrency/design.md). The case removes
  the public `executionStrategy` runner surface in favor of a single optional
  `evaluation.concurrency` setting, defaults omitted concurrency to
  `max(2, runtime.NumCPU()*2)`, keeps harness-backed runs sequential, and
  designs a coordinator-owned scheduler so concurrent evaluator calls cannot
  mutate `evaluation.json` directly or change deterministic result/report
  ordering. Code remains untouched until **In-Progress**. Updated the bundle
  [index](index.md).

- **Done**: Implemented, landed, and archived
  [0194 - Harness-native evaluator dispatch](archive/0194-harness-native-evaluator-dispatch.md).
  The runner gains the reserved `harness` evaluator: a ready judgment work
  unit persists an `awaiting_evaluator` checkpoint (`pendingEvaluatorCall`
  correlation metadata, exit `0`) and emits the complete bounded work request;
  `--resume <run> --evaluator-result <path|->` submits the typed result
  envelope, which is correlation-checked, validated, retried, logged, and
  persisted through the shared evaluator paths, with stale rebuilt input and
  mixed harness runtimes refused as `run_state_invalid` and the first accepted
  result binding `state.harnessIdentity`. `evaluation status`,
  `evaluation list` (new `--state awaiting`), and workspace `status` surface
  the awaiting lifecycle with its continuation action. CLI `auto` discovery
  becomes readiness-aware (executable, authentication where probeable, and
  structured-output capability, with dry-run evidence) and the codex/claude
  adapters use native schema and no-persistence flags when advertised. The
  runtime skill selects `harness` by default in capable agents and services
  the checkpoint loop; the Claude Code and Codex automation guides replace
  their placeholders; `install.md` documents evaluator defaults; and the
  durable CLI, evaluation, and quality-skill specs carry the contract (see the
  [specs log](../specs/log.md)). Verified by the full check suite and a live
  end-to-end checkpoint loop against the built binary. Moved the case into
  [`archive/`](archive/index.md) and updated the bundle [index](index.md).

- **In-Progress**: Advanced
  [0194 - Harness-native evaluator dispatch](0194-harness-native-evaluator-dispatch.md)
  from `Design` to `In-Progress` now that 0193 has landed; implementation of the
  checkpointed harness transport, status surfaces, readiness-aware fallbacks,
  skill routing, and automation guides is underway. Updated the bundle
  [index](index.md).

- **Done**: Landed and archived
  [0193 - Evaluation runner token efficiency](archive/0193-evaluation-token-efficiency.md):
  moved the parent concept and its folder into [`archive/`](archive/index.md)
  and updated the bundle [index](index.md).

- **Design**: Created
  [0194 - Harness-native evaluator dispatch](0194-harness-native-evaluator-dispatch.md)
  with its
  [functional spec](0194-harness-native-evaluator-dispatch/spec.md) and
  [design doc](0194-harness-native-evaluator-dispatch/design.md). The case adds
  a reserved `harness` evaluator whose runner persists an
  `awaiting_evaluator` checkpoint, emits one bounded typed work request, accepts
  a correlated result through resume, and validates and persists it through the
  existing evaluator path. The quality skill selects this transport when no
  explicit or configured evaluator overrides it, while direct CLI `auto`
  discovery gains authentication/capability readiness checks and CLI adapters
  use native schema and no-persistence controls. The case also replaces the
  placeholder Claude Code and Codex automation pages with runnable routine and
  scheduled-task guidance, including artifact, permission, network, and
  credential boundaries. Code remains untouched until **In-Progress**, after
  0193 lands. Updated the bundle [index](index.md).

- **In-Review**: Implemented
  [0193 - Evaluation runner token efficiency](archive/0193-evaluation-token-efficiency.md)
  and advanced it from `In-Progress` to `In-Review`. `evaluator.BuildPrompt`
  now renders a stable prefix (task, expected schema, packaged source, shared
  area context) before the per-work-unit delta and returns the boundary; the
  Anthropic adapter marks the system prompt and stable prefix with
  `cache_control` and folds cache reads/writes into total input; the OpenAI,
  claude, and codex adapters record provider-reported cached input tokens; the
  work graph replaces the `assessRequirement`/`rateRequirement` pair with one
  `assessRateRequirement` evaluator unit whose composite result is split,
  validated, and persisted as the two unchanged payload kinds (a partial
  composite fails as retryable evaluator output); an area's source bundle is
  packaged once per run and memoized; and `logs/evaluator-calls.jsonl` records
  `cachedInputTokens`. Amended the durable
  [evaluator contract](../specs/evaluation/evaluator-contract.md),
  [runner](../specs/evaluation/runner.md),
  [orchestration](../specs/evaluation/orchestration.md), and
  [payload kinds](../specs/evaluation/records/payload-kinds.md) specs; the
  dry-run CLI contract needed no edit. Per-area CLI session reuse stays a
  permitted, unimplemented `MAY`. Verified with `go build`, `go vet`,
  `golangci-lint` (0 issues), the full test suite, and a live
  `qualitymd evaluation run --dry-run --json` (one judgment unit per
  requirement). Updated the bundle [index](index.md).

- **In-Progress**: Advanced
  [0193 - Evaluation runner token efficiency](archive/0193-evaluation-token-efficiency.md)
  from `Design` to `In-Progress` to implement the settled design: stable/delta
  prompt split with provider caching, the combined `assessRateRequirement`
  evaluator unit, per-area source reuse, and cached-input-token logging.
  Updated the bundle [index](index.md).

## 2026-07-09

- **Design**: Advanced
  [0193 - Evaluation runner token efficiency](archive/0193-evaluation-token-efficiency.md)
  from `Draft` to `Design` and added its
  [design doc](archive/0193-evaluation-token-efficiency/design.md). The design keeps the
  three levers independent but landing together: `evaluator.BuildPrompt` renders
  stable-first (source and area frame in a cacheable prefix, requirement and its
  frame in a trailing delta) and returns the boundary; the Anthropic adapter
  marks the prefix with `cache_control` and both API adapters record cached
  input tokens; a single `assessRateRequirement` evaluator unit returns a
  composite that the engine splits back into the unchanged
  `RequirementAssessmentResult` and `RequirementRatingResult` payloads (a partial
  result fails the unit under the retry policy); and an area's source is packaged
  once per run with optional per-area CLI session reuse. Observational
  equivalence is preserved throughout. Updated the bundle [index](index.md). Code
  remains untouched until **In-Progress**.

- **Creation**: Added
  [0193 - Evaluation runner token efficiency](archive/0193-evaluation-token-efficiency.md)
  (`status: Draft`) with its [functional spec](archive/0193-evaluation-token-efficiency/spec.md)
  and [index](archive/0193-evaluation-token-efficiency/index.md). A measured full-model
  run spends ~8.8M input tokens against ~0.3M output (96.7% input) because the
  per-area source bundle is re-uploaded per requirement and sits past the prompt
  cache boundary, every requirement costs two evaluator calls, and the default
  CLI path reuses no context. The spec requires cache-friendly prompt layout
  (source in the stable prefix) with API evaluators applying provider caching,
  a single evaluator call that both assesses and rates each requirement while
  persisting the same two payload kinds, per-area source reuse with optional
  per-area CLI session reuse, and cached-input-token observability — all under
  the existing observational-equivalence invariant so ratings, payloads, and
  reports are unchanged. Design doc deferred to the **Design** phase. Updated
  the bundle [index](index.md).

- **Done**: Implemented and archived
  [0192 - Deterministic evaluation runner](archive/0192-deterministic-evaluation-runner.md).
  `qualitymd evaluation run` is now the evaluation engine: a new
  `internal/runner` package owns the deterministic work graph, execution
  strategy (first slice: `auto` resolving to sequential), evaluator retries,
  atomic per-result persistence into one authoritative `evaluation.json`
  (schema version 4), run-local `logs/events.jsonl` and
  `logs/evaluator-calls.jsonl`, cancellation/resume, and report generation
  through the existing renderer; a new `internal/evaluator` package provides
  the capability-declaring evaluator contract with CLI-backed `codex`/`claude`
  and API-backed `openai`/`anthropic` adapters (`shell`/`manual` reserved),
  selected via `--evaluator`, `evaluation.evaluator`, or `auto` discovery, with
  `evaluators` profiles that reference API keys by environment-variable name
  only. Added [`specs/cli/evaluation-run.md`](../specs/cli/evaluation-run.md)
  and the durable [runner](../specs/evaluation/runner.md),
  [evaluator contract](../specs/evaluation/evaluator-contract.md), and
  [`evaluation.json`](../specs/evaluation/evaluation-json.md) specs; revised
  the evaluation, orchestration, protocol, data-layout, payload-kind, and
  report-tree specs; rewrote the quality-skill evaluation specs and the
  runtime `/quality evaluate` workflow as an agent-mediated wrapper around the
  runner; and updated README, install, reporting-design, and Mintlify skill
  docs plus the generated CLI reference. Updated the bundle
  [index](index.md) and [archive index](archive/index.md).

- **In-Progress**: Advanced
  [0192 - Deterministic evaluation runner](0192-deterministic-evaluation-runner.md)
  from `Design` to `In-Progress`. The functional spec and design doc are
  settled; implementation of `qualitymd evaluation run`, the evaluator
  adapters, the `evaluation.json` run store, and the durable spec/doc/skill
  updates begins. Updated the bundle [index](index.md).

- **Refinement**: Tightened
  [0192 - Deterministic evaluation runner](0192-deterministic-evaluation-runner.md)
  after review, closing six gaps in the
  [functional spec](0192-deterministic-evaluation-runner/spec.md) and
  [design doc](0192-deterministic-evaluation-runner/design.md): evaluators now
  declare execution capabilities the planner reads before scheduling (with a
  `Capabilities()` method on the Go interface and recorded fallbacks);
  `evaluation.json` persists atomically after each accepted work-unit result,
  with serialized writes under parallel execution; resume compatibility is
  defined (supported schema version, same workspace/model path, matching
  evaluator profile — a conflicting `--evaluator` is refused); user
  cancellation leaves the run valid and resumable and reports `cancelled`;
  CLI-backed evaluators must honor non-interactive structured invocation or
  fail selection with `evaluator_incompatible` (both added to the failure
  taxonomy); and the overlapping open questions are resolved — the first slice
  ships `auto`/`sequential` only, `shell`/`manual` stay reserved names, state
  compaction and positional scope are deferred, and provider context
  identifiers live only in the evaluator-call log.

## 2026-07-08

- **Design**: Created
  [0192 - Deterministic evaluation runner](0192-deterministic-evaluation-runner.md)
  with its
  [functional spec](0192-deterministic-evaluation-runner/spec.md) and
  [design doc](0192-deterministic-evaluation-runner/design.md). The case moves
  evaluation toward a CLI-owned deterministic runner with pluggable evaluators,
  `qualitymd evaluation run`, `evaluation.evaluator` config, a single
  authoritative `evaluation.json`, structured run logs, and `/quality evaluate`
  as a wrapper around the runner. Code is not started.

## 2026-07-02

- **Done**: Implemented and archived
  [0191 - Deepen the report gallery exemplar model](archive/0191-report-gallery-model-realism.md).
  The LedgerLite gallery now carries a 39-requirement model with a `codebase`
  area, expanded API/security/auditability/capacity/compatibility coverage,
  quality-named `quality-md` factors, reused sensor-catalog assessments, new
  synthetic changelog entries, and a generated recommendation that matures
  inferential body review into a body-drift detector. The case moved into
  [`archive/`](archive/index.md), and the bundle/archive indexes were updated.

- **In-Progress**: Advanced
  [0191 - Deepen the report gallery exemplar model](archive/0191-report-gallery-model-realism.md).
  The functional spec is settled; no separate design doc is needed because the
  0190 generator shape remains in place and the work is a content/table
  expansion. Implementation begins on the generator-owned LedgerLite
  `QUALITY.md`, synthetic changelog, README, and evaluation payload tables.

- **Done**: Archived
  [0190 - Make the report gallery an authoring exemplar](archive/0190-report-gallery-exemplar.md)
  — set `status: Done` and moved the parent concept and child folder into
  [`archive/`](archive/index.md). The exemplar gallery landed on `main`.
  Updated the bundle [index](index.md), the archive index, and inbound links
  (including 0191's references to the 0190 spec).

- **Done**: Archived
  [0188 - Publish the example evaluation to the docs site](archive/0188-publish-example-evaluation-to-docs.md)
  — set `status: Done` and moved the parent concept and child folder into
  [`archive/`](archive/index.md). The docs sidebar link landed on `main`.
  Updated the bundle [index](index.md) and the archive index.

- **Creation**: Added change
  [0191 - Deepen the report gallery exemplar model](archive/0191-report-gallery-model-realism.md)
  (`status: Draft`) with its
  [functional spec](archive/0191-report-gallery-model-realism/spec.md): deepen the
  LedgerLite gallery beyond 0190's structural exemplar to realistic scale and
  grounding — a new `codebase` constituent, expanded per-area factor families
  (security, auditability, compatibility, capacity, contract currentness),
  quality-named `quality-md` factors, a named sensor catalog with a deliberate
  computational / sensor-plus-guide / inferential assessment spread, a
  re-derived body that earns every factor family from the fictional scenario,
  conformance with the authoring guide family, and run prose consistent with
  the setup and evaluate workflows. Depends on 0190 landing first. Updated the
  bundle [index](index.md).

- **In-Review**: Completed implementation for
  [0190 - Make the report gallery an authoring exemplar](archive/0190-report-gallery-exemplar.md).
  The generator now embeds the exemplar `QUALITY.md`, README, and a synthetic
  quality changelog from content files, builds the expanded evaluation (18
  requirements across a composite root, all four finding types, a not-assessed
  result, six ranked recommendations) through the real pipeline at the stable
  `0001-full-eval` path, and regenerates byte-stably. `CHANGELOG.md` notes the
  gallery upgrade under Unreleased.

- **In-Progress**: Advanced
  [0190 - Make the report gallery an authoring exemplar](archive/0190-report-gallery-exemplar.md).
  Spec and design are settled; implementation begins on the generator
  restructure and exemplar content.

- **Design**: Added the
  [design doc](archive/0190-report-gallery-exemplar/design.md) for 0190: embedded
  content files via `go:embed`, generalized payload tables (multi-factor,
  multi-finding, not-assessed, umbrella-factor roll-up), and generator-owned
  changelog emission, keeping the byte-stable regeneration gate.

- **Draft**: Created
  [0190 - Make the report gallery an authoring exemplar](archive/0190-report-gallery-exemplar.md)
  with its [functional spec](archive/0190-report-gallery-exemplar/spec.md). The case
  expands the LedgerLite gallery into a best-practice QUALITY.md exemplar — full
  body sections with unknowns and review provenance, model-wide
  agent-harnessability with seven sub-factors, a normative service-contract
  area, a veto requirement, a measured rating override, sensor-grounded
  assessments, a not-assessed result, a synthetic quality changelog, and
  realistic finding/recommendation copy — regenerated at the stable
  `0001-full-eval` path. Updated the bundle [index](index.md).

## 2026-06-30

- **Done**: Implemented and archived
  [0189 - Heading sentence case](archive/0189-heading-sentence-case.md).
  Active Markdown/MDX headings, generated Evaluation report titles/headings,
  Contents labels, Mintlify specification output, runtime skill docs, and
  checked-in report-gallery artifacts now use sentence case while preserving
  proper names, formal model vocabulary, model-provided titles, and historical
  archive records.

- **In-Review**: Completed implementation for
  [0189 - Heading sentence case](archive/0189-heading-sentence-case.md). Active
  docs, durable specs, runtime skill content, generated report labels, generated
  docs, tests, and checked-in gallery output are aligned for review.

- **In-Progress**: Advanced
  [0189 - Heading sentence case](archive/0189-heading-sentence-case.md). Functional
  requirements and design are settled; implementation is beginning across active
  docs, specs, runtime skill files, generated report labels, generated docs, and
  checked-in gallery output.

- **Design**: Created
  [0189 - Heading sentence case](archive/0189-heading-sentence-case.md) with its
  [functional spec](archive/0189-heading-sentence-case/spec.md) and
  [design doc](archive/0189-heading-sentence-case/design.md). The case scopes
  sentence-case heading alignment across active docs, specs, Mintlify output,
  runtime skill files, generated reports, and checked-in report-gallery output
  while preserving proper names, model titles, and historical archives. It
  explicitly excludes heading-case automation.

- **In-Review**: Reworked
  [0188 - Publish the example evaluation to the docs site](archive/0188-publish-example-evaluation-to-docs.md)
  to link out instead of generate. Removed `scripts/report-docs.mjs` and the
  generated `mintlify/examples/**` tree, dropped the `report-docs` /
  `report-docs-check` tasks from `mise.toml` (and the `check` gate) and the
  example-page regeneration step from `.githooks/run-check`, and pruned the task
  docs from [`CONTRIBUTING.md`](../CONTRIBUTING.md). `mintlify/docs.json` now
  surfaces the gallery report through a single `Example report` sidebar link that
  opens `report.md` on GitHub. Updated the change case and its functional spec.

- **In-Review**: Created and implemented
  [0188 - Publish the example evaluation to the docs site](archive/0188-publish-example-evaluation-to-docs.md).
  Added `scripts/report-docs.mjs`, which renders the report-gallery LedgerLite
  evaluation into Mintlify pages under `mintlify/examples/software-service/`,
  rewriting cross-links to internal docs routes (and data files / the glossary to
  GitHub) and adding a single `Examples` nav group surfaced as
  "Example quality evaluation". Wired `report-docs` / `report-docs-check` into
  `mise.toml` (the `check` gate) and the `.githooks/run-check` staged path after
  the gallery regenerates; documented the tasks in
  [`CONTRIBUTING.md`](../CONTRIBUTING.md). The full `mise run check` gate and
  `mint broken-links` pass. Updated the bundle [index](index.md).

- **Done**: Implemented and archived
  [0187 - Remove Run Finding Summary](archive/0187-remove-run-finding-summary.md).
  Generated run reports no longer render the standalone `Finding Summary` table
  near `## Key Details`; the total remains in Key Details and the complete
  Finding breakdown remains beside the full `findings.md` link.

- **In-Review**: Completed implementation for
  [0187 - Remove Run Finding Summary](archive/0187-remove-run-finding-summary.md).
  Report rendering, durable specs, docs, changelog, tests, and report-gallery
  output are aligned for review.

- **In-Progress**: Created
  [0187 - Remove Run Finding Summary](archive/0187-remove-run-finding-summary.md)
  with
  its
  [functional spec](archive/0187-remove-run-finding-summary/spec.md) and
  [design doc](archive/0187-remove-run-finding-summary/design.md). Functional
  requirements and design are settled; implementation is beginning across report
  rendering, durable specs, docs, tests, release notes, and generated examples.

- **Done**: Implemented and archived
  [0186 - Report Count Summary Hierarchy](archive/0186-report-count-summary-hierarchy.md).
  Generated run reports now place both full-list links above capped preview
  tables, use marked inline count summaries, label recommendation count groups
  by impact, and render Strength Findings as `💪 Strength`.

- **In-Review**: Completed implementation for
  [0186 - Report Count Summary Hierarchy](archive/0186-report-count-summary-hierarchy.md).
  Report rendering, durable specs, docs, glossary, changelog, tests, and
  report-gallery output are aligned for review.

- **In-Progress**: Advanced
  [0186 - Report Count Summary Hierarchy](archive/0186-report-count-summary-hierarchy.md).
  Functional requirements and design are settled; implementation is beginning
  across report rendering, durable specs, docs, tests, glossary, changelog, and
  generated examples.

- **Design**: Created
  [0186 - Report Count Summary Hierarchy](archive/0186-report-count-summary-hierarchy.md)
  with its
  [functional spec](archive/0186-report-count-summary-hierarchy/spec.md) and
  [design doc](archive/0186-report-count-summary-hierarchy/design.md). The case
  covers moving the full Recommendations report link above its capped table,
  adding semantic markers and clearer grouping grammar to full-list count
  summaries, and changing the Strength Finding type display marker to `💪`.

- **Done**: Implemented and archived
  [0185 - Concern Finding Severity](archive/0185-concern-finding-severity.md).
  Finding `severity` is now concern-only in validation, schema, examples,
  generated reports, runtime guidance, and durable specs.

- **In-Progress**: Created
  [0185 - Concern Finding Severity](archive/0185-concern-finding-severity.md)
  with its
  [functional spec](archive/0185-concern-finding-severity/spec.md) and
  [design doc](archive/0185-concern-finding-severity/design.md). Functional
  requirements and design are settled; implementation is beginning across
  Finding data validation, schemas, reports, runtime guidance, durable specs,
  tests, and generated examples.

- **Done**: Implemented and archived
  [0184 - Evaluation Links Blockquote](archive/0184-evaluation-links-blockquote.md).
  Generated reports now render `Evaluation links:` as a blockquote immediately
  below each H1 while preserving the overview, Findings, Recommendations, and
  glossary links.

- **Done**: Implemented and archived
  [0183 - Glossary and Report Links](archive/0183-glossary-report-links.md).
  Generated reports now link to `report.md`, `findings.md`,
  `recommendations.md`, and the workspace-root `glossary.md`, while local
  `Legend` blocks are removed.

- **In-Review**: Completed implementation for
  [0183 - Glossary and Report Links](archive/0183-glossary-report-links.md).
  The workspace-root glossary, generated `Evaluation links:` navigation,
  removed report legends, durable specs, guidance, tests, and report-gallery
  output are aligned for review.

- **In-Progress**: Advanced
  [0183 - Glossary and Report Links](archive/0183-glossary-report-links.md). Functional
  requirements and design are settled; implementation is beginning across the
  glossary artifact, durable specs, generated report rendering, docs, tests, and
  generated examples.

- **Design**: Added the
  [design doc](archive/0183-glossary-report-links/design.md) for
  [0183 - Glossary and Report Links](archive/0183-glossary-report-links.md).
  The design keeps the glossary authored for this slice, gives `glossary.md` a
  1:1 durable spec, and replaces report legends with one `Evaluation links:`
  line.

- **Draft**: Created
  [0183 - Glossary and Report Links](archive/0183-glossary-report-links.md) with
  its [functional spec](archive/0183-glossary-report-links/spec.md). The case
  defines a shared workspace-root glossary, generated `Evaluation links:`
  navigation, removal of per-report legends, and cleanup for the stale `unknown`
  Finding type spec entry.

- **Done**: Implemented and archived
  [0182 - Finding Summary Display Order](archive/0182-finding-summary-display-order.md).
  Generated run reports now render `Finding Summary`, show all Finding types in
  concern-first order including zero-count rows, and label sparse gap/risk
  severity counts as `Severity`.

- **In-Review**: Completed implementation for
  [0182 - Finding Summary Display Order](archive/0182-finding-summary-display-order.md).
  Report rendering, durable specs, reporting guidance, release notes, tests,
  generated schema, and report-gallery output are aligned for review.

- **In-Progress**: Created
  [0182 - Finding Summary Display Order](archive/0182-finding-summary-display-order.md)
  with its
  [functional spec](archive/0182-finding-summary-display-order/spec.md) and
  [design doc](archive/0182-finding-summary-display-order/design.md). Functional
  requirements and design are settled; implementation is beginning across
  generated report rendering, durable specs, docs, tests, and report-gallery
  output.

## 2026-06-29

- **Done**: Implemented and archived
  [0181 - Evaluation Identity Manifest](archive/0181-evaluation-identity-manifest.md).
  Evaluation manifests now use `EvaluationManifest`, `evaluationId`, nested
  local `run` metadata, and `data/evaluation-manifest.json`; generated run
  report frontmatter now carries `evaluationId`, `created`, `model`, and `run`
  without duplicated scope fields.

- **In-Review**: Completed implementation for
  [0181 - Evaluation Identity Manifest](archive/0181-evaluation-identity-manifest.md).
  Manifest data, report frontmatter, durable specs, runtime skill guidance,
  tests, generated schema, and report-gallery output are aligned for review.

- **In-Progress**: Created
  [0181 - Evaluation Identity Manifest](archive/0181-evaluation-identity-manifest.md)
  with its
  [functional spec](archive/0181-evaluation-identity-manifest/spec.md) and
  [design doc](archive/0181-evaluation-identity-manifest/design.md). Functional
  requirements and design are settled; implementation is beginning across
  Evaluation manifest data, report frontmatter, durable specs, runtime skill
  guidance, tests, and generated examples.

- **Done**: Implemented and archived
  [0180 - Finding Taxonomy and Report Details](archive/0180-finding-taxonomy-report-details.md).
  Evaluation Finding types now use `strength`, `gap`, `risk`, and `note`;
  generated reports render `🚩 Gap`, add a run Finding Breakdown, use a single
  Overall Rating confidence, and include confidence in run Top Recommendations.

- **In-Review**: Completed implementation for
  [0180 - Finding Taxonomy and Report Details](archive/0180-finding-taxonomy-report-details.md).
  Finding type validation, generated reports, durable specs, runtime skill
  guidance, tests, and report-gallery output are aligned for review.

- **In-Progress**: Created
  [0180 - Finding Taxonomy and Report Details](archive/0180-finding-taxonomy-report-details.md)
  with its
  [functional spec](archive/0180-finding-taxonomy-report-details/spec.md) and
  [design doc](archive/0180-finding-taxonomy-report-details/design.md). Functional
  requirements and design are settled; implementation is beginning across
  Finding type validation, generated reports, skill guidance, durable specs,
  tests, and generated examples.

- **Done**: Implemented and archived
  [0179 - Enum Catalog Metadata](archive/0179-enum-catalog-metadata.md).
  Fixed Evaluation enum catalogs now carry type-level labels and descriptions
  plus value descriptions, and generated report local keys render catalog labels.

- **In-Review**: Completed implementation for
  [0179 - Enum Catalog Metadata](archive/0179-enum-catalog-metadata.md).
  Enum catalog metadata, generated report keys, durable report specs, docs,
  tests, changelog, and report-gallery output are aligned for review.

- **In-Progress**: Created
  [0179 - Enum Catalog Metadata](archive/0179-enum-catalog-metadata.md) with its
  [functional spec](archive/0179-enum-catalog-metadata/spec.md) and
  [design doc](archive/0179-enum-catalog-metadata/design.md). Functional
  requirements and design are settled; implementation is beginning across enum
  catalog metadata, generated report keys, durable specs, docs, tests, and
  generated examples.

- **Done**: Implemented and archived
  [0178 - Recommendation Number Columns](archive/0178-recommendation-number-columns.md).
  Generated recommendation list tables now use `#` as the single visible
  recommendation number column and no longer duplicate it with `Rank`.

- **In-Review**: Completed implementation for
  [0178 - Recommendation Number Columns](archive/0178-recommendation-number-columns.md).
  Report rendering, durable report specs, tests, docs, changelog, and
  report-gallery output are aligned for review.

- **In-Progress**: Advanced
  [0178 - Recommendation Number Columns](archive/0178-recommendation-number-columns.md).
  Functional spec and design are settled; implementation is beginning across
  generated recommendation table rendering, durable report specs, tests, docs,
  and generated examples.

- **Design**: Added the
  [design doc](archive/0178-recommendation-number-columns/design.md) for
  [0178 - Recommendation Number Columns](archive/0178-recommendation-number-columns.md).
  The design removes duplicate recommendation ordinal columns in Markdown table
  rendering without changing structured ranking data.

- **Draft**: Created
  [0178 - Recommendation Number Columns](archive/0178-recommendation-number-columns.md)
  with its
  [functional spec](archive/0178-recommendation-number-columns/spec.md). The
  case keeps `#` as the single visible recommendation number in list tables.

- **Done**: Implemented and archived
  [0177 - Report Full List Links](archive/0177-report-full-list-links.md).
  Generated run reports now emphasize the full Findings and Recommendations
  report links below capped top-list tables and show complete ranked counts for
  each linked report.

- **In-Review**: Completed implementation for
  [0177 - Report Full List Links](archive/0177-report-full-list-links.md).
  Renderer code, durable report specs, report design guidance, tests, changelog,
  and report-gallery output are aligned for review.

- **In-Progress**: Advanced
  [0177 - Report Full List Links](archive/0177-report-full-list-links.md).
  Functional spec and design are settled; implementation is beginning across
  generated report rendering, durable report specs, tests, and generated
  examples.

- **Design**: Added the
  [design doc](archive/0177-report-full-list-links/design.md) for
  [0177 - Report Full List Links](archive/0177-report-full-list-links.md). The
  design keeps the change local to run-report rendering, deriving total counts
  from complete ranked Finding and Recommendation artifacts.

- **Draft**: Created
  [0177 - Report Full List Links](archive/0177-report-full-list-links.md) with
  its [functional spec](archive/0177-report-full-list-links/spec.md). The case
  makes `report.md` full-list links for Findings and Recommendations more
  scannable and adds total ranked counts beside those links.

- **Done**: Implemented and archived
  [0176 - Recommendation IDs and Numbers](archive/0176-recommendation-ids-and-numbers.md).
  Evaluation recommendations now persist opaque `qrec_...` IDs for JSON data
  paths and structured ranking/coverage refs; generated reports derive
  user-facing recommendation numbers from ranked order.

- **In-Review**: Completed implementation for
  [0176 - Recommendation IDs and Numbers](archive/0176-recommendation-ids-and-numbers.md).
  Evaluation data contracts, report rendering, durable specs, runtime skill
  guidance, tests, changelog, and report-gallery output are aligned for review.

- **In-Progress**: Advanced
  [0176 - Recommendation IDs and Numbers](archive/0176-recommendation-ids-and-numbers.md).
  Functional spec and design are settled; implementation is beginning across
  Evaluation data contracts, reports, durable specs, runtime skill guidance,
  tests, and generated examples.

- **Design**: Added the
  [design doc](archive/0176-recommendation-ids-and-numbers/design.md) for
  [0176 - Recommendation IDs and Numbers](archive/0176-recommendation-ids-and-numbers.md).
  The design keeps recommendation payloads writable before ranking by assigning
  opaque `qrec_...` IDs, while recommendation numbers are derived from ranked
  order in generated reports and follow-up.

- **Draft**: Created
  [0176 - Recommendation IDs and Numbers](archive/0176-recommendation-ids-and-numbers.md)
  with its
  [functional spec](archive/0176-recommendation-ids-and-numbers/spec.md). The
  case separates non-user-facing recommendation IDs from user-facing
  recommendation numbers.

- **Done**: Implemented and archived
  [0175 - Report Contents Sections](archive/0175-report-contents-sections.md).
  Generated Evaluation reports now use standard Contents sections for
  multi-section report artifacts and no longer render compact `Jump to:` local
  navigation lines.

- **In-Review**: Completed implementation for
  [0175 - Report Contents Sections](archive/0175-report-contents-sections.md).
  Renderer code, durable report specs, report design guidance, release notes,
  tests, and report-gallery output are aligned for review.

- **In-Progress**: Advanced
  [0175 - Report Contents Sections](archive/0175-report-contents-sections.md).
  Functional spec and design are settled; implementation is beginning across
  generated report rendering, durable specs, report design guidance, tests, and
  generated examples.

- **Design**: Added the
  [design doc](archive/0175-report-contents-sections/design.md) for
  [0175 - Report Contents Sections](archive/0175-report-contents-sections.md).
  The design renders deterministic shallow Contents sections from each report
  renderer's known top-level sections and removes compact `Jump to:` output.

- **Draft**: Created
  [0175 - Report Contents Sections](archive/0175-report-contents-sections.md)
  with its [functional spec](archive/0175-report-contents-sections/spec.md).
  The case restores generated report Contents sections for multi-section report
  artifacts and removes compact `Jump to:` local navigation lines.

- **Done**: Implemented and archived
  [0174 - Report Local Keys and Navigation](archive/0174-report-local-keys.md).
  Generated Evaluation reports now use compact local navigation and
  notation-only local keys instead of generated Contents and bottom Legend
  sections, while table cells keep marker-plus-text display labels.

- **In-Review**: Completed implementation for
  [0174 - Report Local Keys and Navigation](archive/0174-report-local-keys.md).
  Renderer code, durable report specs, skill reporting guidance, report design
  guidance, release notes, tests, and report-gallery output are aligned for
  review.

- **In-Progress**: Advanced
  [0174 - Report Local Keys and Navigation](archive/0174-report-local-keys.md).
  Functional spec and design are settled; implementation is beginning across
  generated report rendering, durable specs, report design guidance, tests, and
  generated examples.

- **Design**: Added the
  [design doc](archive/0174-report-local-keys/design.md) for
  [0174 - Report Local Keys and Navigation](archive/0174-report-local-keys.md).
  The design replaces generated contents and legend sections with compact jump
  links and notation-only local keys sourced from the model snapshot and display
  catalogs.

- **Draft**: Created
  [0174 - Report Local Keys and Navigation](archive/0174-report-local-keys.md)
  with its [functional spec](archive/0174-report-local-keys/spec.md). The case
  defines compact report navigation, local indicator keys, detail-report opening
  alignment, and marker/text-label accessibility rules.

- **Done**: Implemented and archived
  [0173 - Evaluation Enum Catalogs](archive/0173-evaluation-enum-catalogs.md).
  Fixed Evaluation enum values now use typed catalogs for validation, generated
  schemas, display labels, markers, and report ordering; durable specs, skill
  guidance, release notes, tests, and report-gallery output are aligned.

- **In-Review**: Completed implementation for
  [0173 - Evaluation Enum Catalogs](archive/0173-evaluation-enum-catalogs.md).
  Evaluation enum catalogs, validation, schema output, report rendering, durable
  specs, skill guidance, release notes, tests, and report-gallery output are
  aligned for review.

- **In-Progress**: Advanced
  [0173 - Evaluation Enum Catalogs](archive/0173-evaluation-enum-catalogs.md).
  Functional spec and design are settled; implementation is beginning across
  Evaluation enum catalogs, validation, schemas, report rendering, durable
  specs, and tests.

- **Design**: Added the
  [design doc](archive/0173-evaluation-enum-catalogs/design.md) for
  [0173 - Evaluation Enum Catalogs](archive/0173-evaluation-enum-catalogs.md).
  The design keeps persisted values stable while deriving validation, schema
  enums, labels, markers, and report ordering from typed Evaluation catalogs.

- **Draft**: Created
  [0173 - Evaluation Enum Catalogs](archive/0173-evaluation-enum-catalogs.md)
  with its [functional spec](archive/0173-evaluation-enum-catalogs/spec.md). The
  case centralizes fixed Evaluation enum values, labels, markers, schema
  validation, and report rendering metadata.

- **Done**: Implemented and archived
  [0172 - Workspace Status Contract](archive/0172-workspace-status-contract.md).
  `qualitymd status` now reports workspace status, emits JSON
  `schemaVersion: 2` with relative workspace metadata, and aligns CLI specs,
  docs, and skill guidance around project/workspace terminology.

- **In-Review**: Completed implementation for
  [0172 - Workspace Status Contract](archive/0172-workspace-status-contract.md).
  Status JSON/help/code, durable CLI and skill specs, README/install guidance,
  runtime skill guidance, tests, and release notes are aligned for review.

- **In-Progress**: Advanced
  [0172 - Workspace Status Contract](archive/0172-workspace-status-contract.md).
  Functional spec and design are settled; implementation is beginning across
  `qualitymd status` JSON/help/code, durable CLI and skill specs, README/install
  guidance, runtime skill guidance, tests, and release notes.

- **Design**: Added the
  [design doc](archive/0172-workspace-status-contract/design.md) for
  [0172 - Workspace Status Contract](archive/0172-workspace-status-contract.md). The
  design exposes relative workspace metadata in status JSON v2 while preserving
  existing status fields and project value-proposition language.

- **Draft**: Created
  [0172 - Workspace Status Contract](archive/0172-workspace-status-contract.md)
  with its [functional spec](archive/0172-workspace-status-contract/spec.md). The case aligns
  status workspace terminology across CLI help, JSON, specs, docs, and skill
  guidance.

- **Done**: Implemented and archived
  [0171 - Run Report Simplification](archive/0171-run-report-simplification.md).
  Generated run-level `report.md` now opens around Summary, Key Details, Model
  Evaluation, ranked findings, ranked recommendations, Legend, and Primary
  Source Data.

- **In-Review**: Completed implementation for
  [0171 - Run Report Simplification](archive/0171-run-report-simplification.md).
  Renderer code, durable report specs, report design guidance, release notes,
  tests, and report-gallery output are aligned for review.

- **In-Progress**: Advanced
  [0171 - Run Report Simplification](archive/0171-run-report-simplification.md).
  Functional spec and design are settled; implementation is beginning across
  generated run report rendering, durable specs, report design guidance, tests,
  release notes, and generated examples.

- **Design**: Added the
  [design doc](archive/0171-run-report-simplification/design.md) for
  [0171 - Run Report Simplification](archive/0171-run-report-simplification.md).
  The design keeps detail report navigation intact while simplifying
  `report.md`, moving the model table under `Model Evaluation`, and narrowing
  source-data lists to report-local primary inputs.

- **Draft**: Created
  [0171 - Run Report Simplification](archive/0171-run-report-simplification.md)
  with its
  [functional spec](archive/0171-run-report-simplification/spec.md). The case
  removes duplicated run-report opening chrome and provenance sections, renames
  the run report model table to `Model Evaluation`, and introduces
  `Primary Source Data`.

- **Done**: Implemented and archived
  [0170 - Report Visual Markers](archive/0170-report-visual-markers.md).
  Generated Evaluation reports now use subtle shape markers for recommendation
  impact, restore impact to run-report Top Recommendations, and use box-like
  Area / Factor Breakdown row markers.

- **In-Review**: Completed implementation for
  [0170 - Report Visual Markers](archive/0170-report-visual-markers.md).
  Renderer code, durable report specs, report design guidance, release notes,
  tests, and report-gallery output are aligned for review.

- **In-Progress**: Advanced
  [0170 - Report Visual Markers](archive/0170-report-visual-markers.md).
  Functional spec and design are settled; implementation is beginning across
  generated report rendering, durable report specs, report design guidance,
  tests, release notes, and generated examples.

- **Design**: Added the
  [design doc](archive/0170-report-visual-markers/design.md) for
  [0170 - Report Visual Markers](archive/0170-report-visual-markers.md). The
  design centralizes recommendation impact display labels and localizes quiet
  Area / Factor row markers to the breakdown table.

- **Draft**: Created
  [0170 - Report Visual Markers](archive/0170-report-visual-markers.md) with
  its [functional spec](archive/0170-report-visual-markers/spec.md). The case
  adds subtle recommendation impact shape markers, restores impact to run-report
  Top Recommendations, and replaces the puzzle-piece breakdown marker with
  box-like Area and Factor markers.

- **Done**: Implemented and archived
  [0169 - Run Report Opening](archive/0169-run-report-opening.md). Generated
  run-level `report.md` now carries non-judgmental run metadata in
  frontmatter, opens with `Summary`, `Key Details`, and `Contents`, moves
  traceability to `Report Details`, and omits the visible limits section for now.

- **In-Review**: Completed implementation for
  [0169 - Run Report Opening](archive/0169-run-report-opening.md). Renderer
  code, durable report specs, report design guidance, release notes, tests, and
  report-gallery output are aligned for review.

- **In-Progress**: Advanced
  [0169 - Run Report Opening](archive/0169-run-report-opening.md). Functional
  spec and design are settled; implementation is beginning across generated run
  report rendering, durable specs, report design guidance, tests, and generated
  examples.

- **Design**: Added the
  [design doc](archive/0169-run-report-opening/design.md) for
  [0169 - Run Report Opening](archive/0169-run-report-opening.md). The design
  keeps the run-report opening human-first while moving routing metadata into
  frontmatter and lower report details.

- **Draft**: Created
  [0169 - Run Report Opening](archive/0169-run-report-opening.md) with its
  [functional spec](archive/0169-run-report-opening/spec.md). The case
  reshapes generated `report.md` around `Summary`, `Key Details`, `Contents`,
  and non-judgmental frontmatter metadata while removing visible limits for now.

- **Done**: Implemented and archived
  [0168 - Quality Evaluation Title](archive/0168-quality-evaluation-title.md).
  Generated run report titles now use `Quality Evaluation - <Area>` with
  factor filters appended in parentheses, durable specs, report design guidance,
  and report-gallery output are aligned, and focused tests cover area-only,
  single-factor, and multiple-factor title content.

- **In-Review**: Completed implementation for
  [0168 - Quality Evaluation Title](archive/0168-quality-evaluation-title.md).
  Renderer code, durable report specs, report design guidance, release notes,
  and focused tests are aligned for review.

- **In-Progress**: Advanced
  [0168 - Quality Evaluation Title](archive/0168-quality-evaluation-title.md).
  Functional spec and design are settled; implementation is beginning across
  generated run report rendering, durable specs, report design guidance, release
  notes, and tests.

- **Design**: Added the
  [design doc](archive/0168-quality-evaluation-title/design.md) for
  [0168 - Quality Evaluation Title](archive/0168-quality-evaluation-title.md).
  The design derives the run report heading from resolved Area and Factor
  labels while preserving stable references in the Scope section.

- **Draft**: Created
  [0168 - Quality Evaluation Title](archive/0168-quality-evaluation-title.md)
  with its
  [functional spec](archive/0168-quality-evaluation-title/spec.md). The case
  renames generated run report titles to `Quality Evaluation - <Area>` and
  appends factor filters in parentheses.

- **Done**: Implemented and archived
  [0167 - Report Frontmatter H1 Titles](archive/0167-report-frontmatter-h1-titles.md).
  Generated Evaluation Markdown report frontmatter `title` values now match the
  visible H1 document title, durable specs and report design guidance are
  aligned, tests cover the title contract, and report-gallery output has been
  regenerated.

- **In-Review**: Completed implementation for
  [0167 - Report Frontmatter H1 Titles](archive/0167-report-frontmatter-h1-titles.md).
  Generated report frontmatter titles now match visible H1 document titles;
  durable report specs, report design guidance, tests, and report-gallery output
  are aligned for review.

- **In-Progress**: Advanced
  [0167 - Report Frontmatter H1 Titles](archive/0167-report-frontmatter-h1-titles.md).
  Functional spec and design are settled; implementation is beginning across
  generated report rendering, durable report specs, report design guidance,
  tests, and report-gallery output.

- **Design**: Added the
  [design doc](archive/0167-report-frontmatter-h1-titles/design.md) for
  [0167 - Report Frontmatter H1 Titles](archive/0167-report-frontmatter-h1-titles.md).
  The design derives frontmatter `title` from the shared report `Heading` value
  so metadata and the visible H1 use one document-title source.

- **Draft**: Created
  [0167 - Report Frontmatter H1 Titles](archive/0167-report-frontmatter-h1-titles.md)
  with its
  [functional spec](archive/0167-report-frontmatter-h1-titles/spec.md). The
  case aligns generated report frontmatter `title` values with visible H1
  document titles while preserving `type` as the report artifact taxonomy and
  keeping report frontmatter tiny.

- **Done**: Implemented and archived
  [0166 - Setup Factor Proposal Checkpoint](archive/0166-setup-factor-proposal-checkpoint.md).
  Setup now teaches and applies Factor-set and individual-Factor desiderata,
  presents a draft factor proposal checkpoint before final review, carries the
  reviewed proposal into authoring, and reports factor coverage/depth issues as
  first-model usefulness gaps.

- **In-Review**: Completed implementation for
  [0166 - Setup Factor Proposal Checkpoint](archive/0166-setup-factor-proposal-checkpoint.md).
  Runtime guidance, durable skill specs, README framing, and OKF logs are aligned
  for review.

- **In-Progress**: Advanced
  [0166 - Setup Factor Proposal Checkpoint](archive/0166-setup-factor-proposal-checkpoint.md).
  Functional spec and design are settled; implementation is beginning across
  setup workflow guidance, factor authoring guidance, durable skill specs, README
  framing, and logs.

- **Design**: Added the
  [design doc](archive/0166-setup-factor-proposal-checkpoint/design.md) for
  [0166 - Setup Factor Proposal Checkpoint](archive/0166-setup-factor-proposal-checkpoint.md).
  The design inserts a factor proposal checkpoint between human-context discovery
  and final review, keeps the checkpoint correction-oriented, and binds the
  reviewed proposal into authoring and closeout checks.

- **Draft**: Created
  [0166 - Setup Factor Proposal Checkpoint](archive/0166-setup-factor-proposal-checkpoint.md)
  with its [functional spec](archive/0166-setup-factor-proposal-checkpoint/spec.md).
  The case captures factor-set desiderata, individual-Factor desiderata, setup
  teaching copy, targeted factor proposal corrections, durable specs, runtime
  guidance, and README framing.

- **Done**: Implemented and archived
  [0165 - Run IDs and Artifact Numbering](archive/0165-run-id-artifact-numbering.md).
  Evaluation runs now persist a globally-unique `RunManifest.id`,
  recommendations use per-run `number` values and numeric ranking/coverage refs,
  finding ranking entries no longer carry synthetic artifact IDs, reports render
  run IDs and typed recommendation references, and durable specs, runtime
  guidance, schema, tests, and the report gallery are aligned.

- **In-Review**: Completed implementation for
  [0165 - Run IDs and Artifact Numbering](archive/0165-run-id-artifact-numbering.md).
  The implementation passed focused Evaluation/CLI tests and the full
  `mise run check` gate before archival.

- **In-Progress**: Advanced
  [0165 - Run IDs and Artifact Numbering](archive/0165-run-id-artifact-numbering.md).
  Functional spec and design are settled; implementation is beginning across
  Evaluation data identity, generated reports, durable specs, runtime guidance,
  tests, and regenerated gallery fixtures.

- **Design**: Added the
  [design doc](archive/0165-run-id-artifact-numbering/design.md) for
  [0165 - Run IDs and Artifact Numbering](archive/0165-run-id-artifact-numbering.md).
  The design generates the run ID at creation from the `createdAt` instant plus
  a `crypto/rand` base32 tail (no scan, opaque on read), renames the
  recommendation `id`→`number` with a zero-padded data path, deletes the
  `QFIND` finding ID and its stability-by-`findingRef` machinery, and pins the
  run ID in the report gallery.

- **Draft**: Created
  [0165 - Run IDs and Artifact Numbering](archive/0165-run-id-artifact-numbering.md)
  with its [functional spec](archive/0165-run-id-artifact-numbering/spec.md). The case
  revises [0163](archive/0163-report-artifact-ids.md): a globally-unique
  `RunManifest.id` replaces the locally-scoped `QEVAL` run identity, recommendations
  carry a per-run `number` instead of `QREC`, and the `QFIND` finding artifact ID
  is removed in favor of the finding's requirement-scoped `findingRef`. Records
  the affected code, durable specs, runtime guidance, and regenerated gallery.

- **Done**: Implemented and archived
  [0164 - Agent Instruction Init Pointer](archive/0164-agent-instruction-init-pointer.md).
  `qualitymd init` now writes a concise, idempotent pointer to local agent
  instruction files by default, supports `--no-agent-instructions`, keeps
  stdout scaffold mode pure, reports pointer effects in human and JSON output,
  points direct init users to `/quality setup`, and setup uses the opt-out while
  handling existing model maturity explicitly.

- **In-Progress**: Advanced
  [0164 - Agent Instruction Init Pointer](archive/0164-agent-instruction-init-pointer.md).
  Functional spec and design are settled; implementation is beginning across
  the CLI init path, scaffold templates, durable specs, runtime setup guidance,
  docs, and tests.

- **Design**: Added the
  [design doc](archive/0164-agent-instruction-init-pointer/design.md) for
  [0164 - Agent Instruction Init Pointer](archive/0164-agent-instruction-init-pointer.md).
  The design keeps pointer update logic in a focused internal package, runs it
  only after successful file scaffolding, preserves stdout purity, and has setup
  opt out with `--no-agent-instructions`.

- **Draft**: Created
  [0164 - Agent Instruction Init Pointer](archive/0164-agent-instruction-init-pointer.md)
  with its
  [functional spec](archive/0164-agent-instruction-init-pointer/spec.md). The
  case adds a concise, idempotent init-time QUALITY.md pointer to agent
  instruction files and routes direct init users toward `/quality setup`.

- **Done**: Implemented and archived
  [0163 - Report Artifact IDs](archive/0163-report-artifact-ids.md).
  `qualitymd evaluation data set` now assigns `QREC` and `QFIND` artifact IDs,
  generated reports render `QEVAL`, `QREC`, and `QFIND` identifiers for handoff
  use, durable specs and runtime guidance are aligned, and the report-gallery
  output has been regenerated.

- **In-Review**: Completed implementation for
  [0163 - Report Artifact IDs](archive/0163-report-artifact-ids.md). Evaluation
  data assignment, schema output, generated reports, durable specs, runtime
  skill guidance, tests, and report-gallery examples are aligned with
  handoff-ready run, recommendation, and ranked finding IDs.

- **In-Progress**: Advanced
  [0163 - Report Artifact IDs](archive/0163-report-artifact-ids.md). Functional spec and
  design are settled; implementation is beginning across Evaluation data
  assignment, generated reports, durable specs, runtime skill guidance, tests,
  and schema output.

- **Design**: Added the
  [design doc](archive/0163-report-artifact-ids/design.md) for
  [0163 - Report Artifact IDs](archive/0163-report-artifact-ids.md). The design makes
  `qualitymd evaluation data set` the narrow assignment surface for `QREC` and
  `QFIND` IDs while preserving Model references and payload-local finding
  selectors.

- **Draft**: Created
  [0163 - Report Artifact IDs](archive/0163-report-artifact-ids.md) with its
  [functional spec](archive/0163-report-artifact-ids/spec.md). The case separates
  citable Evaluation report artifact IDs from Model references, title slugs, and
  generated paths.

- **Done**: Implemented and archived
  [0162 - Report Source Data Section](archive/0162-report-source-data-section.md).
  Generated Evaluation Markdown reports now use identity-only frontmatter and
  end with report-local `Source Data` sections; durable specs, report design
  guidance, focused tests, report-gallery output, and `mise run check` are
  aligned.

- **In-Review**: Completed implementation for
  [0162 - Report Source Data Section](archive/0162-report-source-data-section.md).
  Generated report frontmatter is now identity-only; every generated report ends
  with a report-local `Source Data` section; durable report specs, report design
  guidance, tests, and report-gallery output are aligned.

- **In-Progress**: Advanced
  [0162 - Report Source Data Section](archive/0162-report-source-data-section.md).
  Functional spec and design are settled; implementation is beginning across
  generated report rendering, durable report specs, focused tests, and
  report-gallery output.

- **Design**: Added the
  [design doc](archive/0162-report-source-data-section/design.md) for
  [0162 - Report Source Data Section](archive/0162-report-source-data-section.md).
  The design keeps existing report-local source path computation, removes YAML
  `data` emission, and appends a shared bottom `Source Data` section to every
  generated report.

- **Draft**: Created
  [0162 - Report Source Data Section](archive/0162-report-source-data-section.md)
  with its
  [functional spec](archive/0162-report-source-data-section/spec.md). The case
  moves generated report source-data pointers from noisy YAML frontmatter into a
  stable bottom `Source Data` section.

- **Done**: Implemented and archived
  [0161 - Area / Factor Breakdown Reports](archive/0161-area-factor-breakdown.md).
  Generated run and Area reports now render Area / Factor Breakdown tables with
  Area/Factor links, overall and local ratings, ranked finding counts, and
  ranked recommendation counts; the old `Subject Reports`, Area `Factors`, and
  Area `Child Areas` body sections are removed.

- **In-Review**: Completed implementation for
  [0161 - Area / Factor Breakdown Reports](archive/0161-area-factor-breakdown.md).
  The renderer, durable report specs, report design guide, focused tests, and
  generated report-gallery examples are aligned with the shared Area / Factor
  Breakdown report section.

- **In-Progress**: Advanced
  [0161 - Area / Factor Breakdown Reports](archive/0161-area-factor-breakdown.md).
  Functional spec and design are settled; implementation is beginning across
  generated report rendering, durable report specs, tests, and report-gallery
  output.

- **Design**: Added the
  [design doc](archive/0161-area-factor-breakdown/design.md) for
  [0161 - Area / Factor Breakdown Reports](archive/0161-area-factor-breakdown.md). The
  design renders a shared breakdown from Evaluation artifacts rather than the
  generated report-output manifest, with ranked finding and recommendation
  counts resolved per Area or Factor row.

- **Draft**: Created
  [0161 - Area / Factor Breakdown Reports](archive/0161-area-factor-breakdown.md)
  with its
  [functional spec](archive/0161-area-factor-breakdown/spec.md). The case
  replaces flat generated subject-report tables and split Area structure tables
  with a compact Area / Factor breakdown in generated Evaluation reports.

- **Done**: Implemented and archived
  [0160 - Report Body Rating Drivers](archive/0160-report-body-rating-drivers.md).
  Generated run, Area, and Factor reports no longer render standalone
  rating-driver body sections; structured rating drivers remain in source
  payloads; focused Go tests, full Go tests, and Markdown formatting pass.

- **In-Review**: Completed implementation for
  [0160 - Report Body Rating Drivers](archive/0160-report-body-rating-drivers.md).
  Generated run, Area, and Factor reports no longer render standalone
  rating-driver body sections; durable report specs, tests, and report-gallery
  output are aligned.

- **In-Progress**: Advanced
  [0160 - Report Body Rating Drivers](archive/0160-report-body-rating-drivers.md).
  Functional spec and design are settled; implementation is beginning across
  durable report specs, generated report rendering, focused tests, and
  report-gallery output.

- **Design**: Added the
  [design doc](archive/0160-report-body-rating-drivers/design.md) for
  [0160 - Report Body Rating Drivers](archive/0160-report-body-rating-drivers.md).
  The design removes only the Markdown projection of rating drivers while
  preserving structured payload validation and source-data traceability.

- **Draft**: Created
  [0160 - Report Body Rating Drivers](archive/0160-report-body-rating-drivers.md)
  with its
  [functional spec](archive/0160-report-body-rating-drivers/spec.md). The case
  keeps `ratingDrivers` in structured Evaluation data while removing standalone
  Markdown rating-driver sections from generated report bodies.

- **Done**: Implemented and archived
  [0159 - Report source-data frontmatter](archive/0159-report-source-data-frontmatter.md).
  Generated report frontmatter `data` now lists the structured source payloads
  used to render each report artifact, visible report bodies no longer duplicate
  report-level `Data` columns, the report gallery is regenerated, and
  `mise run check` passes.

- **In-Review**: Completed implementation for
  [0159 - Report source-data frontmatter](archive/0159-report-source-data-frontmatter.md).
  Durable report specs and design guidance now define `data` frontmatter as a
  source-data manifest; the renderer emits report-local source payload lists and
  removes body `Data` summary columns; tests and gallery output are updated.

- **In-Progress**: Advanced
  [0159 - Report source-data frontmatter](archive/0159-report-source-data-frontmatter.md).
  Functional spec and design are settled; implementation is beginning across
  durable report specs, the report design guide, renderer tests, and generated
  report-gallery output.

- **Design**: Added the
  [design doc](archive/0159-report-source-data-frontmatter/design.md) for
  [0159 - Report source-data frontmatter](archive/0159-report-source-data-frontmatter.md).
  The design keeps `data` frontmatter as the only report-level source-data
  pointer surface, removes visible `Data` summary columns, and builds
  deterministic report-local source payload lists.

- **Draft**: Created
  [0159 - Report source-data frontmatter](archive/0159-report-source-data-frontmatter.md)
  with its
  [functional spec](archive/0159-report-source-data-frontmatter/spec.md). The case makes
  generated report `data` frontmatter the report-local source-data manifest and
  removes duplicate source-data links from visible report bodies.

- **Done**: Implemented and archived
  [0158 - OKF-compatible report headers](archive/0158-okf-report-headers.md).
  Generated Evaluation Markdown reports now carry pointer-only OKF-compatible
  frontmatter, subject-shaped report types, visible overview/findings/
  recommendations navigation, run freshness from `RunManifest.createdAt`, and
  regenerated report-gallery examples.

- **In-Review**: Completed implementation for
  [0158 - OKF-compatible report headers](archive/0158-okf-report-headers.md).
  Added shared report frontmatter/header rendering, RunManifest `createdAt`,
  durable report specs, report design guidance, updated tests, and regenerated
  report-gallery output.

- **In-Progress**: Advanced
  [0158 - OKF-compatible report headers](archive/0158-okf-report-headers.md). Functional
  spec and design are settled; implementation is beginning across generated
  Evaluation reports, durable report specs, focused tests, and report-gallery
  examples.

- **Design**: Added the
  [design doc](archive/0158-okf-report-headers/design.md) for
  [0158 - OKF-compatible report headers](archive/0158-okf-report-headers.md).
  The design centralizes report frontmatter, report type taxonomy, navigation,
  and header rendering while keeping report-specific summary data near existing
  report logic.

- **Draft**: Created
  [0158 - OKF-compatible report headers](archive/0158-okf-report-headers.md) with its
  [functional spec](archive/0158-okf-report-headers/spec.md). The case adds
  OKF-compatible pointer frontmatter and consistent human navigation headers to
  generated Evaluation reports without making report Markdown a new source of
  Evaluation result data.

## 2026-06-27

- **Done**: Implemented and archived
  [0157 - Report Markdown Authoring](archive/0157-report-markdown-authoring.md).
  Generated Evaluation reports now use shared Markdown helpers for table rows,
  report links, data links, code spans, and empty cells; table cells escape
  Markdown separators and normalize multiline scalar content; durable report
  specs and release notes are aligned; and `mise run check` passes.

- **In-Review**: Completed implementation for
  [0157 - Report Markdown Authoring](archive/0157-report-markdown-authoring.md).
  Added `internal/markdown`, migrated Evaluation report dynamic rows and links
  through shared helpers, documented table-cell hygiene in durable report specs,
  and verified with `go test ./internal/markdown ./internal/evaluation`,
  `mise run report-gallery-check`, and `mise run check`.

- **In-Progress**: Advanced
  [0157 - Report Markdown Authoring](archive/0157-report-markdown-authoring.md).
  Functional spec and design are settled; implementation is beginning across the
  internal Markdown writer, Evaluation report rendering, focused tests, and
  durable report hygiene docs.

- **Design**: Advanced
  [0157 - Report Markdown Authoring](archive/0157-report-markdown-authoring.md) with its
  [design doc](archive/0157-report-markdown-authoring/design.md). The design adds a
  small standard-library `internal/markdown` writer for generated report
  primitives while keeping Evaluation semantics in the report renderer.

- **Draft**: Created
  [0157 - Report Markdown Authoring](archive/0157-report-markdown-authoring.md) with its
  [functional spec](archive/0157-report-markdown-authoring/spec.md). The case captures
  the Markdown-authoring research conclusion: keep deterministic Go report
  rendering, add a focused internal writer for report primitives, centralize
  table/link/empty-cell escaping, and avoid broad parser or builder dependencies
  unless design proves a concrete benefit.

- **Done**: Implemented and archived
  [0156 - Report Gallery](archive/0156-report-gallery.md). Added a generated
  software-service report gallery under `examples/`, with a sample `QUALITY.md`,
  synthetic current-contract Evaluation payloads, generated human reports,
  `mise run report-gallery`, and `mise run report-gallery-check`.

- **Done**: Implemented and archived
  [0155 - Recommendation Result Shape](archive/0155-recommendation-result-shape.md).
  `RecommendationResult` now uses `description`, `background`, `expectedValue`,
  and `doneCriterion`; generated recommendation reports render from persisted
  Advice data only; Top Recommendations shows `Rank`, `Recommendation`, `Area /
Factors`, and `Reason`; and `mise run check` passes.

- **In-Review**: Completed implementation for
  [0155 - Recommendation Result Shape](archive/0155-recommendation-result-shape.md).
  `RecommendationResult` now uses `description`, `background`, `expectedValue`,
  and `doneCriterion`; generated recommendation reports render from persisted
  Advice data only; Top Recommendations shows `Rank`, `Recommendation`, `Area /
Factors`, and `Reason`; and `mise run check` passes.

- **In-Progress**: Created
  [0155 - Recommendation Result Shape](archive/0155-recommendation-result-shape.md)
  with its [functional spec](archive/0155-recommendation-result-shape/spec.md)
  and [design doc](archive/0155-recommendation-result-shape/design.md). The case
  replaces the old recommendation prose fields with `description`, `background`,
  `expectedValue`, and `doneCriterion`, keeps recommendation reports data-only,
  and reshapes Top Recommendations around `Rank`, `Recommendation`, `Area /
Factors`, and `Reason`.

- **Done**: Implemented and archived
  [0154 - Ranked Findings Reports](archive/0154-ranked-findings-reports.md).
  Evaluation report build now writes a full ranked `findings.md` index, reshapes
  `report.md` Top Findings around linked Finding, Area, and Factor names, links
  finding statements to stable Requirement detail anchors, and renders Advice
  rank/tier/rationale context in Requirement finding details. `go test ./...`
  and `mise run fmt-md-check` pass.

- **In-Progress**: Advanced
  [0154 - Ranked Findings Reports](archive/0154-ranked-findings-reports.md). Functional
  spec and design are settled; implementation is beginning across generated
  report rendering, durable specs, tests, and lifecycle artifacts.

- **Design**: Advanced
  [0154 - Ranked Findings Reports](archive/0154-ranked-findings-reports.md) with its
  [design doc](archive/0154-ranked-findings-reports/design.md). The design uses one
  ranked-finding row model for `report.md` and `findings.md`, stable
  `finding-<id>` anchors on Requirement reports, and a compact Advice ranking
  block inside each finding detail.

- **Draft**: Created
  [0154 - Ranked Findings Reports](archive/0154-ranked-findings-reports.md) with its
  [functional spec](archive/0154-ranked-findings-reports/spec.md) and
  [design doc](archive/0154-ranked-findings-reports/design.md). The case reshapes
  `report.md` Top Findings, adds a full ranked `findings.md` index, links ranked
  finding statements to stable Requirement finding-detail anchors, and adds
  Advice rank/tier/rationale context to Requirement finding details.

- **Done**: Implemented and archived
  [0153 - Evaluate ranking and umbrella-factor authoring clarity](archive/0153-evaluate-ranking-and-umbrella-clarity.md).
  The `evaluate` finding-ranking guidance, example, and protocol/payload-kinds
  rules now lead with completeness (tiers express priority, no finding dropped),
  and the single `FactorAnalysisResult` example models the umbrella case (empty
  `localAnalysis`, child-Factor roll-up in `localAndDescendantAnalysis`) earned by
  a new Factor Traversal rule. `go test ./...` passes; examples verified via the
  CLI. Moved the case into [`archive/`](archive/index.md) and removed it from the
  bundle [index](index.md).

- **Design**: Wrote the
  [design doc](archive/0153-evaluate-ranking-and-umbrella-clarity/design.md) for
  [0153 - Evaluate ranking and umbrella-factor authoring clarity](archive/0153-evaluate-ranking-and-umbrella-clarity.md)
  and advanced it to `Design`. The single `FactorAnalysisResult` example becomes
  the worked umbrella case; the finding-ranking fix combines a reordered skill
  instruction with a low-value tail entry in the example.

- **Draft**: Created
  [0153 - Evaluate ranking and umbrella-factor authoring clarity](archive/0153-evaluate-ranking-and-umbrella-clarity.md)
  with its [functional spec](archive/0153-evaluate-ranking-and-umbrella-clarity/spec.md)
  to remove two `evaluate` authoring ambiguities a feedback log surfaced:
  finding-ranking completeness (tiers express priority, no finding omitted) and
  umbrella-factor analysis scopes (empty `localAnalysis`, child-Factor roll-up in
  `localAndDescendantAnalysis`).

- **Done**: Implemented and archived
  [0152 - Release Reliability](archive/0152-release-reliability.md). Release
  publishing now has credential preflight, draft-first GitHub Releases,
  independent Homebrew and npm jobs, channel verification before publication,
  idempotent repair helpers, updated runbook guidance, and a successful
  `v0.25.4` proof release.

- **In-Review**: Repaired
  [0152 - Release Reliability](archive/0152-release-reliability.md) after the
  `v0.25.3` release workflow verified all release channels but failed the final
  publish job because `gh release edit` had no checked-out repository context.
  The workflow now passes `--repo qualitymd/quality.md`; `v0.25.3` was manually
  published after verification passed.

- **In-Review**: Repaired
  [0152 - Release Reliability](archive/0152-release-reliability.md) after the
  `v0.25.2` release workflow proved that the workflow token cannot read the
  Homebrew tap through the cross-repo contents API during verification. The
  verifier now reads the public raw cask URL, and `v0.25.2` was manually
  verified and published after the Homebrew/npm jobs succeeded.

- **In-Review**: Repaired
  [0152 - Release Reliability](archive/0152-release-reliability.md) after the
  `v0.25.1` release workflow proved that GitHub's tag-specific release endpoint
  can return 404 for draft releases. Release scripts now fall back to listing
  releases by `tag_name`, the Homebrew updater avoids token-dependent clone, and
  `v0.25.1` was manually verified and published after the tap repair.

- **In-Review**: Completed implementation for
  [0152 - Release Reliability](archive/0152-release-reliability.md). The release workflow
  now runs credential preflight, creates draft GitHub releases, publishes
  Homebrew and npm through independently visible jobs, verifies all release
  channels before publication, and provides verifier/repair tasks. `node
scripts/release-verify.mjs v0.25.0` passes against the current public release.

- **Done**: Implemented and archived
  [0151 - Evaluation Report CTA](archive/0151-evaluation-report-cta.md).
  `/quality evaluate` closeouts now point users first to `<run>/report.md` as
  the decision-ready evaluation result, describe `<run>/recommendations.md` as
  the action-planning report, keep machine report indexes out of the primary CTA,
  and preserve evidence, limitations, not-done, and next-action signals.

- **In-Review**: Completed implementation for
  [0151 - Evaluation Report CTA](archive/0151-evaluation-report-cta.md). Durable
  skill specs, runtime evaluate guidance, release notes, and logs now align
  around a value-oriented human report CTA. `mise run fmt-md-check` passes.

- **Design**: Created [0152 - Release Reliability](archive/0152-release-reliability.md)
  with its [functional spec](archive/0152-release-reliability/spec.md) and
  [design doc](archive/0152-release-reliability/design.md). The case hardens release
  preflight credential checks, draft-first GitHub releases, independent
  Homebrew/npm publication, idempotent repair, and post-release verification
  after the `v0.25.0` partial-publish failure.

- **In-Progress**: Advanced
  [0151 - Evaluation Report CTA](archive/0151-evaluation-report-cta.md).
  Functional spec and design are settled; implementation is beginning across
  durable skill specs, runtime guidance, release notes, and Change Case lifecycle
  artifacts.

- **Design**: Advanced
  [0151 - Evaluation Report CTA](archive/0151-evaluation-report-cta.md) with its
  [design doc](archive/0151-evaluation-report-cta/design.md). The design keeps
  generated reports and CLI receipts unchanged, uses `Open next` for the primary
  `report.md` CTA, and keeps `recommendations.md` on the Recommendations line as
  the action-planning report.

- **Draft**: Created
  [0151 - Evaluation Report CTA](archive/0151-evaluation-report-cta.md) with its
  [functional spec](archive/0151-evaluation-report-cta/spec.md). The case makes
  `/quality evaluate` closeouts point users to the human report files with
  value-oriented CTAs, keeps machine report indexes out of the primary closeout
  action, and preserves the evaluate workflow's evidence and mutation-boundary
  signals.

- **Done**: Implemented
  [0150 - Evaluation Advice](archive/0150-evaluation-advice.md). Evaluation data now
  accepts and validates `FindingRankingResult`, `RecommendationResult`, and
  `RecommendationRankingResult`; reportability requires complete Advice and
  finding coverage; reports render Top Findings, Top Recommendations,
  `recommendations.md`, and recommendation detail files; durable specs, runtime
  skill guidance, release notes, and tests are aligned. `go test ./...` and
  `mise run fmt-md-check` pass.

- **Design**: Advanced [0150 - Evaluation Advice](archive/0150-evaluation-advice.md)
  with its [design doc](archive/0150-evaluation-advice/design.md). The design adds three
  Advice payload kinds, keeps coverage accounting on
  `RecommendationRankingResult` as the final Advice closure payload, renders
  human-first recommendation artifacts, and uses progressive QC to catch issues
  at each Advice artifact boundary.

- **Draft**: Created [0150 - Evaluation Advice](archive/0150-evaluation-advice.md) with
  its [functional spec](archive/0150-evaluation-advice/spec.md). The case makes Advice a
  required late evaluation phase with `FindingRankingResult`,
  `RecommendationResult`, and `RecommendationRankingResult` payloads, full
  finding coverage accounting, Top Findings and Top Recommendations in
  `report.md`, and generated recommendation handoff pages.

- **Done**: Implemented and archived
  [0149 - Scope-driven evaluation runs](archive/0149-scope-driven-evaluation-runs.md).
  `qualitymd evaluation create` now writes `RunManifest` scope data, accepts
  `--area`/`--factor` instead of `--narrowing`, report generation renders
  `report.md` as the scoped Area report without positional headline selection,
  Evaluation Frame scope fields were removed, durable specs and skill guidance
  are aligned, and `mise run check` passes.

- **Draft**: Created [0149 - Scope-driven evaluation runs](archive/0149-scope-driven-evaluation-runs.md)
  with [functional spec](archive/0149-scope-driven-evaluation-runs/spec.md) and
  [design doc](archive/0149-scope-driven-evaluation-runs/design.md). Captures
  requested/planned scope as a CLI-owned `RunManifest` at `create` (replacing the
  `--narrowing` slug with `--area`/`--factor`), renders the run report as the
  scoped area report, and removes the positional headline concept — fixing full
  runs that headlined an arbitrary first-listed factor (e.g. Security) instead of
  the overall result.

- **Done**: Implemented and archived
  [0148 - Finding Basis](archive/0148-finding-basis.md). Evaluation data
  contracts, generated schema/examples, reports, durable specs, runtime skill
  guidance, scaffold comments, tests, and release notes now use Finding Core
  `basis` instead of `cause`. `go test ./...` and `mise run fmt-md-check` pass.

- **In-Review**: Completed implementation for
  [0148 - Finding Basis](archive/0148-finding-basis.md). Evaluation data
  contracts, generated schema/examples, reports, durable specs, runtime skill
  guidance, scaffold comments, tests, and release notes now use Finding Core
  `basis` instead of `cause`. `go test ./...` and `mise run fmt-md-check` pass.

- **In-Progress**: Advanced
  [0148 - Finding Basis](archive/0148-finding-basis.md). Functional spec and
  design are settled; implementation is beginning across evaluation records,
  reports, durable specs, skill guidance, examples, tests, scaffold comments,
  and release notes.

- **Design**: Advanced [0148 - Finding Basis](archive/0148-finding-basis.md)
  with its [design doc](archive/0148-finding-basis/design.md). The design
  implements a clean `cause` to `basis` schema rename with no compatibility
  alias, keeping the existing nested status model and report structure.

- **Draft**: Created [0148 - Finding Basis](archive/0148-finding-basis.md) with
  its [functional spec](archive/0148-finding-basis/spec.md). The case renames
  Finding Core `cause` to `basis` across evaluation data, generated reports,
  runtime skill guidance, durable specs, examples, and tests.

- **Done**: Implemented and archived
  [0147 - Report Descendant Terms](archive/0147-report-descendant-terms.md).
  Generated Evaluation reports now use Child Areas and Sub-Factors for immediate
  descendants, durable report specs and release notes are aligned, and
  `go test ./...` plus `mise run fmt-md-check` pass.

- **In-Review**: Completed implementation for
  [0147 - Report Descendant Terms](archive/0147-report-descendant-terms.md). Generated
  Evaluation reports now use Child Areas and Sub-Factors for immediate
  descendants, tests cover the new labels and old-label exclusions, durable
  report specs and release notes are aligned, and `go test ./...` plus
  `mise run fmt-md-check` pass.

- **In-Progress**: Advanced
  [0147 - Report Descendant Terms](archive/0147-report-descendant-terms.md). Functional
  spec and design are settled; implementation is beginning across generated
  report labels, tests, durable report specs, and release notes.

- **Design**: Advanced
  [0147 - Report Descendant Terms](archive/0147-report-descendant-terms.md) with its
  [design doc](archive/0147-report-descendant-terms/design.md). The design keeps
  structured data and traversal helpers unchanged while updating generated
  Markdown report labels and durable report specs.

- **Draft**: Created
  [0147 - Report Descendant Terms](archive/0147-report-descendant-terms.md) with its
  [functional spec](archive/0147-report-descendant-terms/spec.md). The case aligns
  generated report terminology around Child Areas for immediate Area descendants
  and Sub-Factors for immediate Factor descendants.

- **Done**: Implemented and archived
  [0145 - Scannable Skill Output](archive/0145-scannable-skill-output.md).
  The shared agent-mediated UX guide now defines five-second-scan output, while
  `/quality` runtime guidance, durable skill specs, workflow and guide templates,
  logs, and release notes are aligned around labeled, scan-friendly review
  gates, summaries, closeouts, and next-workflow prompts. `mise run
fmt-md-check` and `mise run check` pass.

- **Done**: Implemented and archived
  [0146 - Changelog Directory](archive/0146-changelog-directory.md).
  Model-change history now lives under `.quality/changelog/` with timestamped
  entry filenames, `.quality/logs/` remains a flat workflow/process log
  directory, active skill/spec/docs guidance is aligned, and dogfood data moved
  to the new path. `go test ./...`, `qualitymd status --json`,
  `mise run fmt-md-check`, and `mise run check` pass.

- **Done**: Implemented and archived
  [0144 - Pointed Review Gates](archive/0144-pointed-review-gates.md).
  Shared UX guidance, `/quality` runtime guidance, durable skill specs, authoring
  guide specs, Change Case process guidance, logs, and release notes now require
  purpose-first review gates that ask the user to react to the consequential
  assumption. `mise run fmt-md-check` and `mise run check` pass.

- **Done**: Implemented and archived
  [0143 - Public Review and Improve Workflows](archive/0143-public-review-improve-workflows.md).
  `/quality review` and `/quality improve` are now public, focus-routed workflow
  stubs; README, runtime skill guidance, durable skill specs, guide routing,
  release notes, and logs are aligned. `mise run fmt-md-check` passes.

- **In-Progress**: Advanced
  [0143 - Public Review and Improve Workflows](archive/0143-public-review-improve-workflows.md).
  Functional spec and design doc are settled; implementation is beginning across
  runtime skill guidance, durable skill specs, README framing, and logs.

- **Design**: Advanced
  [0143 - Public Review and Improve Workflows](archive/0143-public-review-improve-workflows.md)
  with its
  [design doc](archive/0143-public-review-improve-workflows/design.md). The design keeps
  review and improve as focus-routed public stubs, delegates safe mutations to
  existing direct-authoring and recommendation-follow-up paths, and leaves deeper
  workflow behavior deferred. Implementation has not started.

- **Done**: Implemented and archived
  [0142 - Requirement Findings Only](archive/0142-requirement-findings-only.md).
  Evaluation data schema version 3 removes Area-analysis findings, `data set`
  and `data verify` enforce finding-backed rated Requirements and non-empty
  rating drivers, Area/Factor reports no longer render Finding sections,
  `/quality evaluate` is aligned with Requirement-only findings, and release
  notes/spec logs are updated. `go test ./...`, `mise run fmt-md-check`, and
  `mise run check` pass.

- **In-Review**: Completed implementation for
  [0142 - Requirement Findings Only](archive/0142-requirement-findings-only.md).
  Evaluation schema, validation, reports, durable specs, runtime skill guidance,
  generated schema, tests, and release notes now make Requirement Findings the
  only Evaluation finding layer and keep Factor/Area roll-up explanation in
  rating drivers and rationale.

- **Draft**: Created
  [0143 - Public Review and Improve Workflows](archive/0143-public-review-improve-workflows.md)
  with its
  [functional spec](archive/0143-public-review-improve-workflows/spec.md). The case makes
  `/quality review` and `/quality improve` public workflows, routes them by
  focus, and aligns README framing around evaluate, review, and improve.

- **In-Progress**: Created
  [0142 - Requirement Findings Only](archive/0142-requirement-findings-only.md)
  with its
  [functional spec](archive/0142-requirement-findings-only/spec.md) and
  [design doc](archive/0142-requirement-findings-only/design.md). The case makes
  Requirement Findings the only Evaluation finding layer, removes
  `AreaAnalysisResult.findings`, requires rated results to be backed by findings
  and rating drivers, and moves Evaluation data to schema version 3.

- **Done**: Implemented and archived
  [0141 - Area-local Factor References](archive/0141-area-local-factor-references.md).
  Requirement Factor references now resolve only within the declaring Area,
  `unknown-factor` rejects ancestor-only references, lint warns on duplicate
  Factor names within one Area's Factor tree, and format/spec/skill/docs
  guidance is aligned. `mise run check` passes.

- **In-Review**: Completed implementation for
  [0141 - Area-local Factor References](archive/0141-area-local-factor-references.md).
  Requirement Factor references now resolve only within the declaring Area,
  `unknown-factor` rejects ancestor-only references, lint warns on duplicate
  Factor names within one Area's Factor tree, and format/spec/skill/docs
  guidance is aligned. `go test ./...`, `go run ./cmd/qualitymd lint QUALITY.md
--json`, `mise run fmt-md-check`, and `mise run check` pass.

- **In-Progress**: Advanced
  [0141 - Area-local Factor References](archive/0141-area-local-factor-references.md).
  Functional spec and design doc are settled; implementation is beginning across
  lint resolution, warning rules, durable format and CLI specs, runtime skill
  guidance, scaffold comments, shared UX examples, logs, and release notes.

- **Design**: Advanced
  [0141 - Area-local Factor References](archive/0141-area-local-factor-references.md)
  with its
  [design doc](archive/0141-area-local-factor-references/design.md). The design
  makes the lint Factor resolver Area-local, adds a same-Area duplicate
  Factor-name warning, and aligns format/spec/skill guidance around local Factor
  ownership. Code changes have not started.

- **Draft**: Created
  [0141 - Area-local Factor References](archive/0141-area-local-factor-references.md)
  with its
  [functional spec](archive/0141-area-local-factor-references/spec.md). The case
  makes Requirement Factor references Area-local, preserves same-named Factors
  across different Areas, and adds a lint warning for duplicate Factor names
  within one Area's Factor tree. Design and implementation have not started.

- **Done**: Implemented and archived
  [0140 - Casual Review Gate Wording](archive/0140-casual-review-gate-wording.md).
  Direct `QUALITY.md` review gates now use a casual "Here's what I'm planning to
  do" checkpoint with a visible value prop, welcoming feedback wording, and the
  same real wait-before-mutation rule. `mise run fmt-md-check` and
  `mise run check` pass.

- **In-Review**: Completed implementation for
  [0140 - Casual Review Gate Wording](archive/0140-casual-review-gate-wording.md).
  Shared UX guidance and `/quality` direct authoring now use a casual "Here's
  what I'm planning to do" checkpoint with a visible value prop, a welcoming
  feedback invitation, and the real wait-before-mutation rule. Durable skill
  specs, authoring guides, logs, and release notes are aligned.
  `mise run fmt-md-check` passes.

- **In-Progress**: Created and advanced
  [0140 - Casual Review Gate Wording](archive/0140-casual-review-gate-wording.md) with
  its
  [functional spec](archive/0140-casual-review-gate-wording/spec.md) and
  [design doc](archive/0140-casual-review-gate-wording/design.md). The case updates
  shared UX guidance and `/quality` direct authoring so review gates state
  "Here's what I'm planning to do", name the value prop, and invite feedback in a
  more conversational way while still waiting before mutation.

- **Done**: Implemented and archived
  [0139 - Real Review Gates](archive/0139-real-review-gates.md).
  The shared agent-mediated UX guide now distinguishes informational previews,
  review gates, and decision gates; requires quick acknowledgement before long
  non-workflow reads; and names feedback-while-proceeding as a false affordance.
  Direct `QUALITY.md` authoring now acknowledges long model/guidance reads before
  they begin, waits for `looks good` or corrections after its checkpoint, and
  keeps judgment-shaping model edits reviewable even when intent is clear.
  `mise run fmt-md-check` passes.

- **In-Review**: Completed implementation for
  [0139 - Real Review Gates](archive/0139-real-review-gates.md).
  The shared UX guide now treats feedback invitations as real review gates and
  adds quick acknowledgement before long non-workflow reads. Direct `QUALITY.md`
  authoring now acknowledges long model/guidance reads, waits after its
  intended-edit checkpoint, and avoids feedback-while-proceeding wording.
  `mise run fmt-md-check` passes.

- **In-Progress**: Advanced [0139 - Real Review Gates](archive/0139-real-review-gates.md).
  Functional spec and design doc are settled; implementation is beginning across
  shared UX guidance, `/quality` runtime skill guidance, durable skill specs,
  authoring guides, logs, and release notes.

- **Design**: Created [0139 - Real Review Gates](archive/0139-real-review-gates.md)
  with its [functional spec](archive/0139-real-review-gates/spec.md) and
  [design doc](archive/0139-real-review-gates/design.md). The case makes feedback
  invitations real review gates, adds quick acknowledgement before long
  non-workflow reads, and aligns direct `QUALITY.md` authoring so it waits for
  review before mutation. Runtime/spec implementation is beginning next.

- **Done**: Implemented and archived
  [0138 - Lightweight Authoring Checkpoint](archive/0138-lightweight-authoring-checkpoint.md).
  Direct `QUALITY.md` authoring now routes separately from tooling `update`,
  infers intent before follow-up, asks only material follow-up questions, states
  the intended edit in a conversational checkpoint, accepts `looks good` as clear
  confirmation when the mutation is named, escalates high-impact model changes to
  decision briefs, and aligns quality-log guidance for meaningful direct
  model-authoring changes. `mise run fmt-md-check` passes.

- **In-Review**: Completed implementation for
  [0138 - Lightweight Authoring Checkpoint](archive/0138-lightweight-authoring-checkpoint.md).
  Direct `QUALITY.md` authoring now has infer-first routing, material follow-up
  thresholds, a conversational intent checkpoint with `looks good` confirmation,
  high-impact decision-brief escalation, and aligned quality-log guidance.
  `mise run fmt-md-check` passes.

- **In-Progress**: Advanced
  [0138 - Lightweight Authoring Checkpoint](archive/0138-lightweight-authoring-checkpoint.md).
  Functional spec and design doc are settled; implementation is beginning across
  durable skill specs, runtime skill guidance, authoring guides, and release
  notes.

- **Design**: Created
  [0138 - Lightweight Authoring Checkpoint](archive/0138-lightweight-authoring-checkpoint.md)
  with its
  [functional spec](archive/0138-lightweight-authoring-checkpoint/spec.md) and
  [design doc](archive/0138-lightweight-authoring-checkpoint/design.md). The case adds a
  lightweight direct-authoring path for `QUALITY.md` edits: infer intent, ask
  follow-up only when material, state the intended edit, accept `looks good` as
  clear confirmation when the mutation is named, and escalate to decision briefs
  for high-impact model changes. Code and durable spec/doc updates have not
  started.

- **Done**: Implemented and archived
  [0137 - Run Report Entrypoint](archive/0137-run-report-entrypoint.md).
  `qualitymd evaluation report build` now writes `report.md` as the run-level
  report, writes the root Area detail as `root-area.md` when present, supports
  scoped Area/Factor reportability without root Area analysis, and emits
  `runReportRef`, `headlineResultRef`, and `headlineReportRef` in
  `EvaluationOutputResult`. Durable specs, runtime skill guidance, generated
  schema, changelog, and tests are aligned. `mise run check` passes.

- **In-Progress**: Advanced
  [0137 - Run Report Entrypoint](archive/0137-run-report-entrypoint.md). Functional
  spec and design doc are settled; implementation is beginning across
  Evaluation report rendering, output contracts, durable specs, runtime skill
  guidance, generated schema, and tests.

- **Design**: Created
  [0137 - Run Report Entrypoint](archive/0137-run-report-entrypoint.md) with its
  [functional spec](archive/0137-run-report-entrypoint/spec.md) and
  [design doc](archive/0137-run-report-entrypoint/design.md). The case separates the
  run-level `report.md` from root Area details, moves the root Area report to
  `root-area.md`, and makes scoped Area/Factor reportability use recorded
  Evaluation Frame scope instead of requiring root Area analysis.

- **Done**: Implemented and archived
  [0136 - Candidate Actions Payload](archive/0136-candidate-actions-payload.md).
  Requirement Finding payloads now use `candidateActions` instead of legacy
  `actions`; each candidate action carries a local `id`, `description`, and
  optional `rationale`; validation rejects legacy `actions`, duplicate
  candidate action IDs within a Finding, missing fields, and candidate actions
  on Area Findings. Generated examples, schema, durable specs, and `/quality`
  runtime guidance are aligned.

- **Design**: Created
  [0136 - Candidate Actions Payload](archive/0136-candidate-actions-payload.md)
  with its
  [functional spec](archive/0136-candidate-actions-payload/spec.md) and
  [design doc](archive/0136-candidate-actions-payload/design.md). The case
  renames finding-local `actions` to `candidateActions`, adds a local candidate
  action ID, preserves the non-rendered raw-material boundary, and defers
  recommendation/advice modeling.

- **Done**: Implemented and archived
  [0135 - Structured Finding Core](archive/0135-structured-finding-core.md).
  Requirement and Area Findings now share a structured Finding Core with
  statement, condition, criteria, cause, effect, and evidence. Evaluation data
  validation rejects legacy finding `description`, `summary`, and top-level
  `rationale` fields; examples and generated schema are updated; Requirement,
  Area, and Factor reports render one findings table/detail structure; durable
  specs and `/quality` runtime guidance now require type-specific structured
  finding analysis. `go test ./...` passes.

- **In-Review**: Completed implementation for
  [0135 - Structured Finding Core](archive/0135-structured-finding-core.md).
  Updated the Evaluation data contract, examples, generated schema, report
  renderer, tests, format spec, durable Evaluation specs, `/quality` skill spec,
  runtime skill guidance, and scaffold comment. `go test ./...` passes.

- **Done**: Implemented and archived
  [0134 - Model-relative workspace paths](archive/0134-model-relative-workspace-paths.md).
  Workspace config and artifacts now resolve relative to the selected
  `QUALITY.md`, while the Git root remains the containment boundary. Evaluation
  list/latest/data/status/report paths now accept `--model` for nested
  workspaces, status and lint use the model-relative contract, and durable CLI,
  skill, and install docs are aligned. Focused verification passes for
  `go test ./internal/status ./internal/cli ./internal/lint` and
  `go test ./internal/evaluation -run 'TestCreateRun|TestListRunsIgnoresUnrecognizedRunFolders|TestRunNameRecognition'`.

- **In-Progress**: Advanced
  [0135 - Structured Finding Core](archive/0135-structured-finding-core.md). Functional
  spec and design doc are settled; implementation is beginning across the
  Evaluation data contract, examples, generated schema, reports, durable specs,
  tests, and bundled skill runtime guidance.

- **Design**: Created
  [0135 - Structured Finding Core](archive/0135-structured-finding-core.md) with its
  [functional spec](archive/0135-structured-finding-core/spec.md) and
  [design doc](archive/0135-structured-finding-core/design.md). The case aligns
  Requirement and Area Findings around a shared Finding Core: statement,
  condition, criteria, cause, effect, and evidence; preserves payload-local
  finding IDs; unifies Requirement, Area, and Factor report rendering; and
  updates `/quality` skill guidance for type-specific finding analysis. Code and
  durable spec/doc updates have not started.

- **Design**: Created
  [0134 - Model-relative workspace paths](archive/0134-model-relative-workspace-paths.md)
  with its
  [functional spec](archive/0134-model-relative-workspace-paths/spec.md) and
  [design doc](archive/0134-model-relative-workspace-paths/design.md). The case
  makes the selected `QUALITY.md` the anchor for workspace config and artifact
  paths, retains the Git repository root as a containment boundary, and adds
  `--model` to evaluation history/latest commands so nested workspaces do not
  depend on cwd or repo-root defaults. Code and durable spec/doc updates had not
  started.

## 2026-06-26

- **Done**: Implemented and archived
  [0133 - Richer evaluation data examples](archive/0133-richer-evaluation-data-examples.md).
  `qualitymd evaluation data example <kind>` now emits fuller representative
  payloads across all supported kinds, including populated representative limits
  and stop conditions, canonical Area/Factor/Requirement/Rating references, and
  Area/Factor/Requirement report refs. Added all-kind structural validation and
  reference-shape assertions; updated the durable CLI spec to keep examples
  representative rather than exhaustive. `go test ./...` and
  `mise run fmt-md-check` pass.

- **In-Review**: Completed implementation for
  [0133 - Richer evaluation data examples](archive/0133-richer-evaluation-data-examples.md).
  `qualitymd evaluation data example <kind>` now emits fuller representative
  artifacts with populated limits/stop conditions where relevant, canonical
  Area/Factor/Requirement/Rating references in subjects, routine refs, and
  report refs, and a valid string `areaId` for `AreaAnalysisResult`. Added an
  all-kind structural example test and reference-shape assertions, and updated
  the durable CLI spec to keep examples representative rather than exhaustive.
  `go test ./...` and `mise run fmt-md-check` pass.

- **In-Progress**: Advanced
  [0133 - Richer evaluation data examples](archive/0133-richer-evaluation-data-examples.md).
  Functional spec and design doc are settled; implementation is beginning across
  `internal/evaluation/data.go`, Evaluation example tests, and the durable CLI
  data spec.

- **Design**: Advanced
  [0133 - Richer evaluation data examples](archive/0133-richer-evaluation-data-examples.md)
  to `Design` with its
  [design doc](archive/0133-richer-evaluation-data-examples/design.md). The design
  enriches the existing one-artifact-per-kind example constructors with shared
  reference and limit helpers, demonstrates Area/Factor/Requirement/Rating
  references in subjects, routine refs, and report refs, and adds all-kind
  structural validation tests. Code not started.

- **Draft**: Created
  [0133 - Richer evaluation data examples](archive/0133-richer-evaluation-data-examples.md)
  with its
  [functional spec](archive/0133-richer-evaluation-data-examples/spec.md), at `Draft`.
  The case makes `qualitymd evaluation data example <kind>` payloads fuller
  representative artifacts, demonstrates Area, Factor, Requirement, and Rating
  Level reference IDs in examples, clarifies that examples are not exhaustive
  enum/status/error corpora, and adds all-kind example validation. Code not
  started. Added the case to the open [index](index.md).

- **Done**: Implemented and archived
  [0132 - Remove info finding severity](archive/0132-remove-info-finding-severity.md).
  Removed `info` from the Evaluation finding severity enum for Requirement and
  Area Findings, regenerated `evaluation-data.schema.json`, updated report
  severity display/sort helpers, durable specs, bundled skill guidance, and
  release notes for `v0.21.0`. `qualitymd evaluation data schema
area-analysis-result` now exposes only `critical`, `high`, `medium`, and
  `low`; `go test ./internal/evaluation` passes.

- **In-Progress**: Advanced
  [0132 - Remove info finding severity](archive/0132-remove-info-finding-severity.md).
  Functional spec and design doc are settled; implementation is beginning across
  the Evaluation data contract, generated schema, report sort/display helpers,
  durable specs, bundled skill runtime guidance, tests, and release notes.

- **Design**: Advanced
  [0132 - Remove info finding severity](archive/0132-remove-info-finding-severity.md)
  to `Design` with its
  [design doc](archive/0132-remove-info-finding-severity/design.md). The design removes
  `info` at the typed finding severity enum source, regenerates the emitted data
  schema, removes `info` from known report display/sort helpers while keeping
  defensive unknown-value fallbacks, and updates specs/skill guidance to route
  informational observations to finding `type: note`. Code not started.

- **Draft**: Created
  [0132 - Remove info finding severity](archive/0132-remove-info-finding-severity.md)
  with its
  [functional spec](archive/0132-remove-info-finding-severity/spec.md), at `Draft`. The
  case removes `info` from the Evaluation finding severity vocabulary so
  severity stays an adverse-finding scale (`critical`, `high`, `medium`, `low`)
  and informational observations use finding `type: note`. It covers data
  validation/schema output, report sort/display helpers, durable specs, bundled
  skill guidance, and release notes. Conditional severity applicability by
  finding type is explicitly deferred. Code not started. Added the case to the
  open [index](index.md).

- **Done**: Implemented and archived
  [0131 - Area findings in evaluation reports](archive/0131-area-findings.md).
  `AreaAnalysisResult.findings` now carries traceable Area Findings with closed
  type/severity/confidence and Factor relationship enums, same-Area Factor
  relationship validation, and closed object rejection of advice/ranking fields.
  Area reports render all Area Findings; Factor reports render matching Area
  Findings; durable specs, runtime skill guidance, release notes, generated
  schema, and focused tests are updated. `mise run check` passes.

- **In-Progress**: Advanced
  [0131 - Area findings in evaluation reports](archive/0131-area-findings.md).
  Functional spec and design doc are settled; implementation is beginning across
  the evaluation data contract, report tree, durable evaluation/skill specs, and
  runtime skill guidance.

- **Design update**:
  [0131 - Area findings in evaluation reports](archive/0131-area-findings.md) now
  explicitly requires Area Finding and Factor relationship closed vocabularies to
  appear as JSON Schema enums, requires `data set` to reject out-of-set values,
  and clarifies that forbidden advice/ranking fields are rejected by closed
  object validation rather than ignored.

- **Design**: Advanced
  [0131 - Area findings in evaluation reports](archive/0131-area-findings.md) to
  `Design` with its [design doc](archive/0131-area-findings/design.md). The design keeps
  Area Findings inside `AreaAnalysisResult.findings`, adds targeted validation
  for local IDs and same-Area Factor relationships, renders Area-owned findings
  directly in Area reports and filtered findings in Factor reports, and leaves
  recommendations plus global top-finding synthesis deferred. Code not started.

- **Draft**: Created
  [0131 - Area findings in evaluation reports](archive/0131-area-findings.md) with its
  [functional spec](archive/0131-area-findings/spec.md). The case adds analysis-phase
  Area Findings on `AreaAnalysisResult.findings`, projects matching findings into
  Factor reports, uses shared finding type/severity/confidence vocabulary, and
  defers recommendations, global top-finding synthesis, and impact/priority
  ranking. Code and durable-spec implementation have not started.

- **Done**: Implemented and archived
  [0130 - Self-contained per-kind data schema](archive/0130-self-contained-data-schema.md).
  `evaluation data schema <kind>` now emits the requested kind's object schema at
  the document root with `$schema` and a kind-qualified `$id`, so required fields,
  enum sets, and reference patterns are legible without following a `$ref` into
  `$defs`; the no-argument full-surface schema remains `$defs` + `oneOf`.
  `data schema` now uses the shared schema presenter, durable CLI/skill specs and
  runtime skill guidance identify schema as the required-fields/enum source, and
  examples as concrete instances. Added focused schema/CLI tests and release
  notes; `go test ./internal/evaluation` and `go test ./internal/cli` pass.

- **In-Progress**: Advanced
  [0130 - Self-contained per-kind data schema](archive/0130-self-contained-data-schema.md).
  Functional spec and design doc are settled; implementation is beginning across
  `internal/evaluation/data_contract.go`, `internal/cli/evaluation.go`, tests,
  durable CLI/skill specs, and the bundled skill runtime guidance.

- **Design**: Advanced
  [0130 - Self-contained per-kind data schema](archive/0130-self-contained-data-schema.md)
  to `Design` with its
  [design doc](archive/0130-self-contained-data-schema/design.md). The design found the
  evaluation data kind objects carry no inter-kind `$ref` — `schemaForObject`
  inlines everything — so "flattening" the single-kind form is just shedding the
  top-level `$defs`/`$ref` envelope and emitting the kind's object schema with the
  document header; the generic `internal/schema` flattener and `cli/schema.go`
  extraction scoped in Draft are unnecessary and were dropped from the footprint.
  Terminal rendering reuses the existing in-package `writeSchema`. Corrected the
  parent's [Affected artifacts](archive/0130-self-contained-data-schema.md) and the open
  [index](index.md). Code not started.

- **Draft**: Created
  [0130 - Self-contained per-kind data schema](archive/0130-self-contained-data-schema.md)
  with its
  [functional spec](archive/0130-self-contained-data-schema/spec.md), at `Draft`. The case
  makes `evaluation data schema <kind>` self-contained — required fields and enum
  value sets legible at the root, no `$ref`/`$defs` dereference for the requested
  kind — so closed sets stop being discovered via `data set --dry-run` failures
  (per the 0.18.0 evaluate feedback log). Keeps `qualitymd schema` and
  `evaluation data schema` as separate commands (rejecting a unified
  `qualitymd schema <kind>`), adds no `data describe` and no `data example`
  annotation, and sharpens the `/quality` skill to read the schema as the
  constraint source. Records [`specs/cli/evaluation-data.md`](../specs/cli/evaluation-data.md),
  the `internal/schema`/`internal/evaluation` code, and the bundled skill as
  affected. Code not started. Added the case to the open [index](index.md).

- **Done**: Landed and archived
  [0129 - Evaluation orchestration overhaul](archive/0129-evaluation-orchestration-overhaul.md).
  Removed the evaluate rigor dial from the bundled skill and durable skill specs,
  made exhaustive in-scope Requirement coverage mandatory, added
  parallel-by-default collection/QC with serial fallback, and promoted
  verification into an always-on two-pronged QC phase. Updated runtime/durable
  logs, the agent-mediated UX guide, and changelog. `mise run check` passes.

- **In-Progress**: Advanced
  [0129 - Evaluation orchestration overhaul](archive/0129-evaluation-orchestration-overhaul.md).
  Functional spec and design doc are settled; implementation is beginning across
  the bundled evaluate workflow, durable skill specs, and agent-mediated UX guide
  to remove evaluation rigor, make exhaustive coverage mandatory, default
  collection/QC to subagent fan-out when available, and add the always-on
  two-pronged QC phase.

- **Design**: Created
  [0129 - Evaluation orchestration overhaul](archive/0129-evaluation-orchestration-overhaul.md)
  with its
  [functional spec](archive/0129-evaluation-orchestration-overhaul/spec.md) and
  [design doc](archive/0129-evaluation-orchestration-overhaul/design.md), and
  advanced it to `Design`. The case removes the evaluation rigor dial
  (`quick`/`standard`/`deep`, `--rigor`, `/quality evaluate deep`, the `Rigor:`
  run-frame field, the feedback-log `rigor:` field), makes exhaustive coverage
  and an always-on two-pronged QC phase (verify ∥ completeness sweep) the
  mandatory evaluate contract, and makes parallel subagent fan-out the default
  execution strategy with an identical serial fallback. No CLI change; **modeling
  rigor** and **assessment rigor** are unrelated and untouched. Code not started.
  Added the case to the open [index](index.md).

- **Done**: Landed and archived
  [0128 - Agent-mediated skill alignment](archive/0128-agent-mediated-skill-alignment.md).
  Completed the remaining `/quality` skill alignment fixes against the
  agent-mediated UX guide, moved the parent and child folder into
  [`archive/`](archive/), updated the archive [index](archive/index.md), and
  removed the case from the open [index](index.md). `mise run check` passes.

- **In-Review**: Completed implementation for
  [0128 - Agent-mediated skill alignment](archive/0128-agent-mediated-skill-alignment.md).
  Setup now renders the run frame first and separates `QUALITY.md` mutation from
  workflow feedback-log creation, setup feedback-log timing is consistent across
  runtime and workflow conventions, recommendation follow-up opens with a
  non-public follow-up frame, and read-only orientation has a status-first output
  shape. Updated durable skill specs, runtime skill files, logs, and changelog.
  `mise run check` passes.

- **In-Progress**: Created
  [0128 - Agent-mediated skill alignment](archive/0128-agent-mediated-skill-alignment.md)
  with its
  [functional spec](archive/0128-agent-mediated-skill-alignment/spec.md) and
  [design doc](archive/0128-agent-mediated-skill-alignment/design.md), then advanced it
  to implementation. The case closes the remaining `/quality` skill
  agent-mediated UX gaps from the audit: setup's run frame must be first in the
  first-output block, setup must disclose workflow feedback-log writes separately
  from `QUALITY.md` mutation, setup feedback-log timing must be consistent,
  recommendation follow-up needs an opening frame, and read-only orientation
  needs a status-first output shape. No Go code change is planned.

- **Done**: Landed and archived
  [0126 - Bulk data set](archive/0126-bulk-data-set.md). Shipped the array-only
  `evaluation data set` batch contract with all-or-nothing validation/write
  behavior, indexed diagnostics, duplicate-derived-path rejection, staged writes
  with rollback, and batch receipts. Updated the durable CLI/evaluation/skill
  specs and the evaluate workflow so routine payloads are persisted with one
  whole-batch dry-run and one write. `mise run check` passes.

- **Done**: Landed and archived
  [0127 - Introspection-first CLI workflow conventions](archive/0127-introspection-first-cli-reference.md).
  Renamed/refocused the bundled skill resource to
  `cli-workflow-conventions.md`, removed embedded command/flag listings,
  preserved non-introspectable workflow conventions, and routed command, flag,
  and payload discovery to CLI introspection. `mise run check` passes.

- **In-Review**: Completed implementation for
  [0126 - Bulk data set](archive/0126-bulk-data-set.md). `evaluation data set` now accepts
  a non-empty JSON array, rejects bare objects and empty arrays, validates every
  payload against one loaded `model-snapshot.md`, aggregates indexed diagnostics,
  rejects duplicate derived paths, stages writes before committing them with
  rollback, and emits a batch receipt with `count` plus input-order `writes`.
  Updated CLI human output, next actions, tests, durable CLI/evaluation/skill
  specs, and the evaluate workflow to batch routine payloads with one dry-run and
  one write. `mise run check` passes.

- **In-Progress**: Advanced
  [0126 - Bulk data set](archive/0126-bulk-data-set.md). Functional spec and design doc
  are settled; implementation is beginning across `internal/evaluation/data.go`,
  `internal/cli/evaluation.go`, tests, durable CLI specs, and the evaluate
  workflow.

- **Design**: Added the
  [0126 - Bulk data set](archive/0126-bulk-data-set.md)
  [design doc](archive/0126-bulk-data-set/design.md) and advanced the case to `Design`.
  The design splits `SetData` into array decoding, candidate validation, duplicate
  path preflight, dry-run receipt generation, and staged write commit with
  best-effort rollback. The batch receipt is a summary object with `count` and
  input-order `writes[]` entries (`index`, `kind`, `path`).

- **In-Review**: Completed implementation for
  [0127 - Introspection-first CLI workflow conventions](archive/0127-introspection-first-cli-reference.md).
  Renamed `skills/quality/resources/cli-quick-reference.md` to
  `cli-workflow-conventions.md`, removed the embedded command/flag listing
  including the 0125 `model` rows and the stale single-payload `data set` row,
  kept the workspace-artifact, feedback-log, narrowing-slug, command-rule, and
  orchestration conventions, and updated `SKILL.md`, the resource index/log, and
  `specs/skills/quality-skill/quality-skill.md`. No CLI code changed.
  `mise run fmt-md-check` and `mise run npm-pack-check` pass.

- **In-Progress**: Advanced
  [0127 - Introspection-first CLI workflow conventions](archive/0127-introspection-first-cli-reference.md).
  Functional spec is settled and no design doc is planned; implementation is the
  bundled skill resource/refocus plus skill-spec/reference updates, with no CLI
  code change.

- **Done**: Landed and archived
  [0125 - Model query commands](archive/0125-model-query-commands.md). Shipped the
  read-only `qualitymd model` group (`tree`/`list`/`get`) projecting a model's
  elements, canonical reference IDs, and containment. Moved the reference grammar
  (path types, encoders, parsers, existence helpers) and a new shared projection
  walk into `internal/model` (`reference.go`, `projection.go`), so `model` depends
  on neither `status` nor `evaluation`; rating references stay in `evaluation`.
  Retired the duplicate model-tree walk in `status` by deriving its shape counts
  from the shared projection. Wired the evaluate workflow to query in-scope
  canonical IDs from the run's `model-snapshot.md` via `model list --json` instead
  of hand-deriving them. Durable specs (`specs/cli/model.md` new, `specs/cli.md`,
  `specs/cli/index.md`, the evaluate skill spec) and docs (bundled skill quick
  reference, README) brought in sync. Moved the parent and folder into
  [`archive/`](archive/); updated the archive [index](archive/index.md) and the
  bundle [index](index.md). Full build and test suite green.

- **Draft**: Created
  [0127 - Introspection-first CLI workflow conventions](archive/0127-introspection-first-cli-reference.md)
  with its [functional spec](archive/0127-introspection-first-cli-reference/spec.md).
  Renames/refocuses the skill's former `cli-quick-reference.md` as
  `cli-workflow-conventions.md`: strips the duplicated command/flag listings the
  CLI's own `--help` and discovery commands (`evaluation data
kinds`/`example`/`schema`) already provide, retains the content that is
  non-introspectable skill convention (workspace-artifact layout, feedback-log
  sequencing, narrowing-slug rule, do/don't rules, orchestration sequences), and
  routes command/flag/payload discovery to the CLI's structured introspection
  channels (`--json`/schema/example, not human tables). Refreshed after landing
  [0125](archive/0125-model-query-commands.md): 0127 now also strips the
  `model tree`/`model list`/`model get` rows 0125 added to the embedded listing.
  Resolves the skill spec's contradiction — `quality-skill.md:146-147` prescribes
  reading the embedded reference while `:589-594` mandates introspection over an
  embedded list that drifts. No code change (the introspection surface already
  exists); no CLI behavior change. Durable spec edit:
  `specs/skills/quality-skill/quality-skill.md`. Touches the same resource as
  [0126](archive/0126-bulk-data-set.md); 0127 absorbs the shared-file edit. Listed in the
  bundle [index](index.md).

- **Update**: Strengthened
  [0126 - Bulk data set](archive/0126-bulk-data-set.md) — added the normative Skill
  integration requirements to its [functional spec](archive/0126-bulk-data-set/spec.md)
  (SK1: persist a scope's payloads in one batched `data set` rather than per
  element; SK2: one whole-batch `--dry-run` instead of per-payload), making skill
  adoption a conformance requirement rather than only an acceptance checkbox, in
  line with 0125's precedent. Refreshed after landing
  [0125](archive/0125-model-query-commands.md) and advancing
  [0127](archive/0127-introspection-first-cli-reference.md): 0125 is now the committed
  baseline, and 0127 absorbs the old quick-reference edit by removing the stale
  single-payload listing wholesale.

- **Draft**: Created
  [0126 - Bulk data set](archive/0126-bulk-data-set.md) with its
  [functional spec](archive/0126-bulk-data-set/spec.md). Replaces `evaluation data set`'s
  single-object stdin contract (and the v0 "MUST NOT accept batch payloads" rule)
  with an **array-only** batch contract: one invocation reads a JSON array of
  payloads, validates the whole batch first, and writes **all-or-nothing** — one
  bad element rejects the batch with per-index diagnostics and nothing is written.
  A one-payload write becomes a one-element array (clean break, not a second
  path); intra-batch duplicate derived paths are rejected; `--json` emits a batch
  receipt; `--dry-run` validates the whole batch without persisting. The driving
  need is agent round-trips: the evaluate workflow loops one `data set` per
  Requirement/Factor/Area (~115 invocations for the cited acquire-roi-next run);
  batching collapses that to one. Write-side companion to
  [0125](archive/0125-model-query-commands.md) (read-side ID friction); 0126 now
  builds on 0125's committed `model list --json` snapshot-query workflow text. No
  format-spec or payload-schema change (the array is a transport envelope).
  Durable spec edits: `specs/cli/evaluation-data.md`, `specs/cli.md`,
  `specs/skills/quality-skill/workflows/evaluate.md`. Design doc deferred
  (atomic write-staging mechanism and receipt schema are design questions).
  Listed the case in the bundle [index](index.md). Code not started.

- **Done**: Implemented and archived
  [0124 - Constrain reference kind fields to closed kind vocabularies](archive/0124-reference-kind-enum.md).
  Enum-constrained the `kind` member of both Evaluation reference shapes from
  typed sources: `routineRefContract` (used for `inputRefs[]` / `*Ref`) against
  `supportedDataKinds` — all ten payload kinds, including the CLI-owned
  `EvaluationOutputResult`, with the agent-writable `acceptedDataKinds` derived
  from it — and `reportRefContract` against `reportKinds` (`area`, `factor`,
  `requirement`). Implementation corrected the design's premise that both `kind`
  fields name payload kinds: a report reference `kind` names a report kind, so it
  is pinned to the report-kind set instead. Regenerated
  `evaluation-data.schema.json`, added a routine-reference rejection test, and
  recorded the rule in `json-conventions.md`. `mise run check` passes.

- **Done**: Implemented and archived
  [0123 - Render interactions through native affordances](archive/0123-native-interaction-affordances.md).
  Made progressive enhancement the default agent-mediated interaction contract:
  each interaction is an intent rendered through a fit-for-purpose native affordance
  (option picker, confirm/approve gate, the harness's own authorization prompt)
  when present, with the numbered-option and `y`/`n` forms reframed as the text
  fallback, the teaching kept in the message not the widget, and a
  harness-authorized mutation no longer double-gated — all stated tool-agnostically.
  Lifted the pattern from the `setup` workflow into the shared `/quality`
  interaction contract; brought the agent-mediated UX guide, `SKILL.md`, the
  `setup`/`update`/`evaluate`/recommendation-follow-up runtime files, the durable
  skill specs, changelog, and logs into conformance. `mise run check` passes.

- **Draft**: Created
  [0125 - Model query commands](0125-model-query-commands.md) with its
  [functional spec](0125-model-query-commands/spec.md). Adds a read-only
  `qualitymd model` group — `tree` (hierarchical view), `list` (flat, filterable
  enumeration), `get <id>` (single-element detail) — that projects a model's
  elements, their canonical reference IDs (`area:…`/`factor:…::…`/`requirement:…::…`),
  labels, and containment, with a human default form and `--json` everywhere. The
  driving need is payload authoring: nothing emits canonical IDs today, so agents
  hand-derive tens of them from `QUALITY.md`. Bounded to _structure and identity_
  to avoid overlap — `status` keeps state/readiness/source-coverage/counts,
  `evaluation` keeps runs/payloads/reports/snapshots; `model` has no `--run` flag,
  no provenance, no validation (defers to `lint`), no ratings. `--area` accepts
  the canonical `area:<path>` form only (one addressing vocabulary). New durable
  spec `specs/cli/model.md`; `specs/cli.md` and `specs/cli/index.md` gain the
  command. Because the driving need is the evaluate workflow itself, the case now
  also covers wiring `/quality` evaluate to consume `model list --json` (against
  the run's `model-snapshot.md`) instead of hand-deriving IDs — a required
  outcome, captured as spec requirements SK1–SK4. Added the
  [design doc](0125-model-query-commands/design.md): the shared projection and
  the canonical-reference grammar (encoders plus parsers) move into
  `internal/model` so `model` depends on neither `status` nor `evaluation`,
  `status` folds onto the same projection walk, and the three verbs land in one
  slice. Listed the case in the bundle [index](index.md). Code not started.

- **Design**: Created
  [0124 - Constrain reference kind fields to the payload-kind vocabulary](0124-reference-kind-enum.md)
  with its [functional spec](0124-reference-kind-enum/spec.md) and
  [design doc](0124-reference-kind-enum/design.md). The `kind` member of
  Evaluation reference objects (`routineRefContract` / `reportRefContract`) is the
  one required reference field still typed as a free `dataString`, so a misspelled
  or invented kind persists silently while every other closed vocabulary
  (`finding.type`, `severity`, `confidence`, `status`) is enum-validated and the
  reference `subject` is resolved against the model snapshot. The case
  enum-constrains reference `kind` to the full supported payload-kind set
  (including the CLI-owned `EvaluationOutputResult`, distinct from the
  agent-writable `acceptedDataKinds` gate), sourced from one typed list, so an
  out-of-vocabulary kind is rejected at `evaluation data set` time and surfaces as
  an `enum` in `evaluation data schema`. Durable change lands in
  [`json-conventions.md`](../specs/evaluation/records/json-conventions.md); no
  `SchemaVersion` bump. Listed the case in the bundle [index](index.md). Code not
  started.

- **Design**: Created
  [0123 - Render interactions through native affordances](archive/0123-native-interaction-affordances.md)
  with its [functional spec](archive/0123-native-interaction-affordances/spec.md)
  and [design doc](archive/0123-native-interaction-affordances/design.md). Lifts
  the `setup`
  workflow's agent-agnostic "render through your own interaction capabilities"
  three-tier pattern into the shared `/quality` interaction contract so every
  workflow renders questions, closed choices, and confirmation gates through a
  fit-for-purpose native affordance when present — reframing the numbered-option
  and `[y]`/`[n]` prose blocks as the text fallback, keeping the teaching in the
  message not the widget, staying agent-agnostic, and not stacking a prose gate on
  a mutation the harness already authorizes. The
  [agent-mediated UX guide](../docs/guides/agent-mediated-ux.md) gained a
  "Channels and progressive enhancement" section (already applied); the shared
  contract, skill workflows, and durable skill specs are not yet reshaped. Listed
  the case in the bundle [index](index.md). No code touched.

- **Done**: Implemented and archived
  [0122 - Finding-level candidate actions](archive/0122-finding-candidate-actions.md).
  Typed the finding `actions` field as an array of non-binding candidate-action
  objects (`description` + optional `rationale`) with CLI validation and an
  example payload, dropped the report's per-finding `Actions` row, and had the
  evaluate workflow record candidate actions on `gap`/`risk` findings (omitting
  `strength`) as raw material for a future Advise phase — kept out of the
  Evaluation v0 report and closeout. Brought the evaluation and quality-skill
  durable specs, `SPECIFICATION.md`, the bundled skill, changelog, and logs into
  conformance. `mise run check` passes.

- **Draft**: Created
  [0122 - Finding-level candidate actions](archive/0122-finding-candidate-actions.md)
  with its [functional spec](archive/0122-finding-candidate-actions/spec.md) and an
  [advise sketch](archive/0122-finding-candidate-actions/advise-sketch.md). Option **A**
  from the recommendation-infrastructure discussion: turn the untyped,
  never-populated `findings[].actions` stub into a typed, non-binding
  candidate-action field (`description` + optional `rationale`) that the assessor
  records on `gap`/`risk` findings as raw material for a future Advise phase,
  validated by the CLI and kept out of the Evaluation v0 report and closeout to
  stay inside the v0 "no recommendations" boundary. The advise sketch captures
  option **B** (the dedicated Advise phase) as forward-looking, non-binding
  context. Registered a `Sketch` concept type in [schema.md](schema.md) and
  listed the case in the bundle [index](index.md). No code touched.

- **Done**: Implemented and archived
  [0121 - Scannable interaction hierarchy](archive/0121-scannable-interaction-hierarchy.md).
  Fixed the _flat-wall_ failure in agent-mediated output — decision gates,
  discovery questions, progress, and result blocks that stacked the question and
  supporting labels at equal weight so the call to action was buried and collapsed
  when bold was stripped. The driving guide
  [`docs/guides/agent-mediated-ux.md`](../docs/guides/agent-mediated-ux.md) now
  teaches a question → separated choices → demoted-rationale gate shape, an
  Emphasis "don't rely on bold alone" rule, an Avoid counter-example, and updated
  Checklist items. The shared `/quality` interaction contract, every mutation gate
  in the bundled skill (`setup` write/update, `update` apply-plan, recommendation
  apply and issue creation), setup discovery ordering, evaluate progress beats,
  and the recommendation result block were reshaped to match, with durable skill
  specs aligned. Verified with `mise run fmt-md-check`.

- **Done**: Implemented and archived
  [0120 - String model-identity fields in evaluation data](archive/0120-string-model-identity-fields.md).
  Evaluation routine and report JSON now persist Area, Factor, Requirement, and
  Rating Level identities as canonical qualified-reference strings, with
  `schemaVersion` bumped to 2 and old structured identity shapes rejected.
  `qualitymd lint` now rejects `root` as a reserved Area name, the companion
  schemas are regenerated, `/quality` compatibility moves to the `0.16` CLI
  line, and the QUALITY.md specification version is `0.5 (Draft)`. Verified with
  `go test ./...` and `mise run fmt-md-check`.

- **Done**: Implemented and archived
  [0119 - Report header kind prefix and title-first layout](archive/0119-report-header-kind-prefix.md).
  Generated Evaluation reports now render a kind-prefixed H1 as the first line,
  move navigation/context trails below the title, drop the redundant `Path:` /
  `Name:` header lines, and keep the root `Area:` trail element tied to the Model
  `title`. Verified with `go test ./...` and `mise run fmt-md-check`.

- **Done**: Implemented and archived
  [0118 - Report empty-cell marker and legend](archive/0118-report-empty-cell-legend.md).
  Generated Evaluation reports now render empty scalar table cells as `—`,
  render empty paired Confidence/Status components as `—`, and append one static
  legend defining the marker to every report while leaving status labels and
  worded empty-section rows distinct. Verified with `go test ./...` and
  `mise run fmt-md-check`.

- **Creation**: Added
  [0120 - String model-identity fields in evaluation data](0120-string-model-identity-fields.md)
  (`Draft`) with its child [index](0120-string-model-identity-fields/index.md)
  and [functional spec](0120-string-model-identity-fields/spec.md). The case
  collapses the structured `areaId` (array), `factorId`, and `requirementId`
  (objects) identity fields in persisted Evaluation routine and report JSON into
  single canonical qualified model-reference strings — the lossless string
  encoding [`SPECIFICATION.md`](../SPECIFICATION.md) already defines — keeps the
  `*Id` names (not `*Ref`, which denotes routine-record links), and reserves
  `root` as a forbidden Area name to close the one rendering ambiguity. It
  reverses the structured-shape mandate in
  [`json-conventions.md`](../specs/evaluation/records/json-conventions.md) (the
  rule [0058](archive/0058-model-reference-identifiers.md) /
  [0059](archive/0059-unqualified-model-references.md) established) and bumps the
  data `SchemaVersion`. Affected artifacts derived by a repo-wide sweep. No code
  touched (gated to In-Progress). Updated the bundle [index](index.md).

- **Design**: Advanced
  [0119 - Report header kind prefix and title-first layout](0119-report-header-kind-prefix.md)
  to `Design`. The functional spec passed the Draft→Design gate (each requirement
  single-obligation and observable; the set, satisfied, achieves the motivation).
  Added the child [design doc](0119-report-header-kind-prefix/design.md): the
  change is a reorder of the three report renderers (kind-prefixed H1 first,
  existing trail writers moved below it, `Path:` / `Name:` lines deleted); the
  root-element Model `title` is already satisfied by `areaTitle(spec, nil)` and
  only locked by a test; and the one real decision — keeping the trails as
  separate paragraphs rather than hard-break tight-stacking — is recorded with
  the rejected alternatives (kicker line, trail-first, `Area:` trail relabel). No
  code touched (gated to In-Progress).

- **Design**: Advanced
  [0118 - Report empty-cell marker and legend](0118-report-empty-cell-legend.md)
  to `Design`. The spec passed the Draft→Design quality bar (each requirement has
  one obligation and an observable result; the set cleanly partitions empty cell
  vs. empty section vs. not-assessed outcome). Added the child
  [design doc](0118-report-empty-cell-legend/design.md): the empty-cell policy
  lives in the `markdownCell` chokepoint (one line) plus a shared `emDashIfEmpty`
  helper for the three composite pair cells, leaving `humanizeEnum` a pure
  humanizer; an unconditional `writeEvaluationLegend` appends one static `## Legend`
  to every report. Records why `N/A`, per-cell asterisks, a `humanizeEnum`
  chokepoint, and a conditional legend were rejected. No code touched (gated to
  In-Progress).

- **Creation**: Drafted
  [0119 - Report header kind prefix and title-first layout](0119-report-header-kind-prefix.md)
  (`status: Draft`). Generated reports open with a navigation trail and then a
  bare H1 — just the subject title, with no statement of the report's kind, which
  matters because an Area, Factor, and Requirement can share a title. The case
  prefixes each report's H1 with its kind (`# Area:` / `# Factor:` /
  `# Requirement:`), renders that title first with the `Area:` trail (and the
  Factor `Factor:` trail / Requirement `Factors:` line) following it, drops the
  now-redundant `Path:` / `Name:` identifier line (the canonical ID still lives in
  filenames and `Data` links), and locks the `Area:` trail's root element to the
  Model `title`. Added the child [spec](0119-report-header-kind-prefix/spec.md);
  the durable contract lands in
  [`specs/evaluation/reports/report-tree.md`](../specs/evaluation/reports/report-tree.md)
  at In-Progress.

- **Creation**: Drafted
  [0118 - Report empty-cell marker and legend](0118-report-empty-cell-legend.md)
  (`status: Draft`). Reports render absent scalar values as blank cells, which a
  reader cannot tell apart from "not applicable" or a rendering fault, and the
  builder carries several uncoordinated empty-value treatments. The case settles
  on one em-dash (`—`) marker for any empty scalar cell — including the composite
  Confidence and analysis-status pair cells (`— / —`) — plus one static legend
  per generated report that defines it, while leaving the not-assessed status
  labels, Rating Level labels, and parenthetical empty-section rows untouched.
  Added the child [spec](0118-report-empty-cell-legend/spec.md); the durable
  contract lands in
  [`specs/evaluation/reports/report-tree.md`](../specs/evaluation/reports/report-tree.md)
  at In-Progress.

- **Done**: Implemented and archived
  [0117 - Requirement report Factors line](archive/0117-requirement-report-factors-line.md).
  Requirement reports now render a plural `Factors:` context line after the
  `Area:` trail, with attached Factor links joined as a flat set, and their
  summary table now contains `Rating`, `Assessment`, `Confidence`, and `Data`
  only. The report-tree spec carries the plural-line rule while preserving the
  prohibition on singular `Factor:` breadcrumbs. Verified with
  `go test ./internal/evaluation ./internal/cli`.

- **Done**: Implemented and archived
  [0116 - Drop the "Evaluation v2" naming](archive/0116-drop-evaluation-v2-naming.md).
  The live spec bundle is now `specs/evaluation/` with
  `specs/evaluation/evaluation.md` as the parent concept; live docs, specs,
  runtime skill text, CLI help strings, and private Go report symbols use plain
  "Evaluation". Removed the superseded `evaluation-v2-sketch.md` and renamed
  `internal/evaluation/report_v2.go` to `report_tree.go`. Verification searches
  for the retired name/path now return only unrelated dependency/tool version
  numbers outside the evaluation surface. Verified with
  `go test ./internal/evaluation ./internal/cli`.

- **Done**: Implemented and archived
  [0115 - Type-safe, model-bound Evaluation v2 data](archive/0115-evaluation-data-typed-contract.md).
  `qualitymd evaluation data set` now validates payloads against a single typed
  data-kind contract, rejects unknown fields and unresolved model/rating IDs
  against the run snapshot, and applies the same checks under `--dry-run`.
  Added populated examples, the generated
  `internal/evaluation/evaluation-data.schema.json` artifact with a drift test,
  kebab-case kind resolution, and the `data schema` / `data verify` commands.
  The CLI data specs and `/quality` skill resources now direct agents to schema
  first, examples second, and dry-run for authored-payload validation. Verified
  with `go test ./internal/evaluation ./internal/cli`.

- **Creation**: Drafted
  [0117 - Requirement report Factors line](0117-requirement-report-factors-line.md)
  (`status: Draft`). A Requirement sits along two axes — its declaring Area and
  the set of Factors it serves — but the report only surfaces the Area axis as a
  header trail; the Factor axis is buried in a summary-table column. The case adds
  a plural `Factors:` context line to the Requirement report header (a flat,
  `;`-separated set of attached-Factor links, parallel to the `Area:` trail and
  mirroring the Factor report's stacked trails) and drops the now-redundant
  `Factors` table column. It reconciles the existing report-tree rule forbidding a
  singular `Factor:` breadcrumb: the new line is plural and parent-free, so the
  prohibition stands. Added the child
  [spec](0117-requirement-report-factors-line/spec.md) and updated the bundle
  [index](index.md).

- **Creation**: Drafted
  [0116 - Drop the "Evaluation v2" naming](0116-drop-evaluation-v2-naming.md)
  (`status: Draft`). The [clean break](archive/0097-evaluation-v2-clean-break.md)
  made v2 the only runtime evaluation workflow, so the `v2` qualifier is now
  vestigial — it falsely implies a coexisting v1. The case retires it from the
  live surface (the `specs/evaluation-v2/` → `specs/evaluation/` folder rename and
  inbound links, CLI and skill specs, `SPECIFICATION.md`, the bundled `/quality`
  skill, and Go symbols/`report_v2.go`), removes the superseded
  `evaluation-v2-sketch.md`, and fences off the load-bearing `schemaVersion`
  marker and frozen history (archive, `log.md`, existing `CHANGELOG.md` entries)
  as non-goals. Added the child [spec](0116-drop-evaluation-v2-naming/spec.md)
  and updated the bundle [index](index.md).

- **Creation**: Drafted
  [0115 - Type-safe, model-bound Evaluation v2 data](0115-evaluation-data-typed-contract.md)
  (`status: Draft`). Motivated by a field evaluation where `qualitymd evaluation
data set` accepted findings with wrong field names (rendered blank in the
  report) and references to non-existent model nodes, because v2 data has no typed
  source of truth. The case defines each kind as typed Go definitions feeding
  strict decode, model-binding validation against the run snapshot, and generated
  schema + populated examples, plus `data schema` and `data verify` commands.
  Added the child [spec](0115-evaluation-data-typed-contract/spec.md) and
  [design](0115-evaluation-data-typed-contract/design.md) and updated the bundle
  [index](index.md).

- **Done**: Implemented and archived
  [0114 - Run frame as first output](archive/0114-run-frame-first-output.md).
  The run-frame "first output, before any tool call" timing rule and a provisional
  / `resolving…` value allowance for tool-dependent fields are now stated in the
  shared homes — the `SKILL.md` dispatcher run-frame instruction and the durable
  `Run frames` spec section — rather than re-derived per workflow. The `evaluate`
  runtime procedure now emits the frame first (invocation-derived model path,
  provisional `Scope: resolving…`) before workspace resolution, and the `evaluate`
  durable spec's Required flow states the requirement; `setup` and `update` were
  already compliant. The companion guide `docs/guides/agent-mediated-ux.md` was
  synced ahead with the Opening section. Verified with `mise run fmt-md-check`.
  Moved the parent and child folder into [`archive/`](archive/), updated the
  [archive index](archive/index.md); the bundle [index](index.md) already lists no
  open cases.

- **Done**: Implemented and archived
  [0113 - Evaluation run folder naming](archive/0113-evaluation-run-folder-naming.md).
  New Evaluation v2 run folders are now named `NNNN-full-eval` for full runs or
  `NNNN-<scope-path>-eval` for narrowed runs. The parser recognizes only the
  current `NNNN-<scope>-eval` grammar. The `/quality evaluate` workflow and
  durable skill spec now pass `--narrowing` as the Area/Factor full structural
  path, and the CLI spec, quick reference, changelog, and v2 sketch are aligned.
  Verified with `mise run check`. Moved the parent and child folder into
  [`archive/`](archive/), updated the [archive index](archive/index.md), and
  removed the entry from the bundle [index](index.md).

- **In-Progress**: Advanced
  [0113 - Evaluation run folder naming](archive/0113-evaluation-run-folder-naming.md).
  Added the design doc and began implementation of the `NNNN-full-eval` /
  `NNNN-<scope-path>-eval` run-folder grammar.

- **Done**: Created, implemented, and archived
  [0110 - Run frame title and workflow vocabulary](archive/0110-run-frame-and-workflow-vocabulary.md).
  The `/quality` run-frame header is now `**Quality · <workflow>**` instead of
  the fake-command `**/quality run**`, the `Mode:` field is dropped (the workflow
  name moves into the header), and the durable Run frames spec now forbids a
  command-style header or a `Mode:` field label. "Mode" is retired in favor of
  "workflow" for the public-surface concept across the durable skill specs,
  bundled runtime skill, and two durable docs, with recommendation follow-up
  described as a post-evaluation follow-up that is not a public workflow. Verified
  with `mise run fmt-md-check`. Moved the parent and child folder into
  [`archive/`](archive/), updated the [archive index](archive/index.md), and
  removed the entry from the bundle [index](index.md).

- **In-Review**: Implemented
  [0112 - Evaluation model snapshot filename](0112-evaluation-model-snapshot-filename.md).
  The Evaluation v2 run-folder model snapshot is now written and read as
  `model-snapshot.md` through a single `evaluation.ModelSnapshotFile` constant —
  `create.go` (write), `path.go` (run-folder validation), `load.go` and
  `report_v2.go` (parse), and `status.go` (staleness) — with no old-name
  compatibility path and operator messages naming the new file. Renamed this
  repo's two tracked active dogfood snapshots (`0005`, `0006`) and their
  run-local prose; left frozen `archive/` runs untouched. Aligned the durable
  CLI, report-tree, and `/quality` skill specs plus the bundled evaluate
  workflow, and updated the seed-layout test. `mise run check` passes. Not yet
  committed or archived. Updated the bundle [index](index.md).

- **Draft**: Created
  [0113 - Evaluation run folder naming](archive/0113-evaluation-run-folder-naming.md).
  Shortens the Evaluation v2 run-folder tag from `quality-eval` to `eval` (the
  `.quality/evaluations/` parent already says "evaluation") and gives every run a
  uniform `NNNN-<scope>-eval` shape: full runs carry an explicit reserved `full`
  marker (`NNNN-full-eval`), and a narrowed run's `--narrowing` slug is the
  scope's full structural path — the Area path from the root plus the Factor path
  when scoping to a Factor, hyphen-joined, with no kind marker or boundary
  separator. The run number stays the identity and `model.md` the structural
  source of truth. Added the parent concept, child folder, functional spec, and
  bundle [index](index.md) entry.

- **Draft**: Created
  [0112 - Evaluation model snapshot filename](0112-evaluation-model-snapshot-filename.md).
  Renames the Evaluation v2 run-folder model snapshot from `model.md` to
  `model-snapshot.md` so the filename signals a frozen point-in-time copy of the
  working-tree model rather than the live model. A clean break with no old-name
  reader; this repo's own two tracked active dogfood runs are renamed so they keep
  validating, while frozen `archive/` runs are left untouched. Added the parent
  concept, child folder, functional spec, design doc, and bundle
  [index](index.md) entry.

- **Done**: Created, implemented, and archived
  [0111 - Evaluation report rating labels](archive/0111-evaluation-report-rating-labels.md).
  Evaluation v2 Area and Factor report header tables now read `Overall Rating` /
  `Local Rating`, and the Area report Factors and Sub-Areas tables and the Factor
  report child Factors table now render a `Local Rating` column from
  `localAnalysis` beside a descendant-inclusive `+ Sub-Factors Rating` /
  `+ Sub-Areas Rating` column from `localAndDescendantAnalysis` (em dash when the
  node has no descendants), replacing the prior aggregate-in-`Rating` column and
  `Yes`/`No` boolean — satisfying the unmet distinction clean-break case 0097
  required. Removed the now-unused `v2BoolLabel` wrapper. Aligned the durable
  [report-tree spec](../specs/evaluation-v2/reports/report-tree.md) and the
  [v2 sketch](../evaluation-v2-sketch.md), and updated focused report tests
  (including a second navigation-fixture rating level to verify the
  local-vs-roll-up split). Verified with `mise run check`. Updated the
  [archive index](archive/index.md).

- **Done**: Created, implemented, and archived
  [0111 - Evaluation report rating labels](archive/0111-evaluation-report-rating-labels.md).
  Evaluation v2 Area and Factor report header tables now read `Overall Rating` /
  `Local Rating`, and the Area report Factors and Sub-Areas tables and the Factor
  report child Factors table now render a `Local Rating` column from
  `localAnalysis` beside a descendant-inclusive `+ Sub-Factors Rating` /
  `+ Sub-Areas Rating` column from `localAndDescendantAnalysis` (em dash when the
  node has no descendants), replacing the prior aggregate-in-`Rating` column and
  `Yes`/`No` boolean — satisfying the unmet distinction clean-break case 0097
  required. Removed the now-unused `v2BoolLabel` wrapper. Aligned the durable
  [report-tree spec](../specs/evaluation-v2/reports/report-tree.md) and the
  [v2 sketch](../evaluation-v2-sketch.md), and updated focused report tests
  (including a second navigation-fixture rating level to verify the
  local-vs-roll-up split). Verified with `mise run check`. Updated the
  [archive index](archive/index.md).

- **Done**: Created, implemented, and archived
  [0109 - Filename text for evaluation data links](archive/0109-evaluation-data-link-filenames.md).
  Evaluation v2 `Data`-column links now render their payload base filename as
  link text (`area-analysis-result.json`, `factor-analysis-result.json`,
  `requirement-assessment-result.json`, `requirement-rating-result.json`)
  instead of the generic words `analysis`/`assessment`/`rating`; link targets
  and structured `data/` paths are unchanged. Aligned the durable
  [report-tree spec](../specs/evaluation-v2/reports/report-tree.md) and updated
  focused report tests. Verified with `mise run check`. Updated the bundle
  [index](index.md) and [archive index](archive/index.md).

- **Done**: Implemented and archived
  [0108 - Short evaluation report filenames](archive/0108-short-evaluation-report-filenames.md).
  Evaluation v2 now keeps the root Area report at `report.md`, writes descendant
  Area, Factor, and Requirement reports with short subject-aware filenames,
  carries those paths through generated links and `EvaluationOutputResult`
  report refs, and keeps generated descendant reports out of repository Markdown
  formatting. Durable specs, the v2 sketch, and runtime `/quality evaluate`
  guidance are aligned. Verified with `mise run check`.

- **In-Review**: Advanced
  [0108 - Short evaluation report filenames](archive/0108-short-evaluation-report-filenames.md)
  to `In-Review`. Evaluation v2 now keeps only the root Area report at
  `report.md`; descendant Area, Factor, and Requirement reports use short
  subject-aware filenames in generated links and report refs. Updated durable
  report layout specs, `/quality` reporting guidance, the v2 sketch, dprint
  generated-report excludes, and focused report tests. Verified with
  `mise run check`.

- **In-Progress**: Advanced
  [0108 - Short evaluation report filenames](archive/0108-short-evaluation-report-filenames.md)
  to `In-Progress`; implementation and durable contract updates are now
  underway.

- **Design**: Advanced
  [0108 - Short evaluation report filenames](archive/0108-short-evaluation-report-filenames.md)
  to `Design` and added its
  [design doc](archive/0108-short-evaluation-report-filenames/design.md). The design
  keeps report path derivation centralized in the Evaluation v2 report helpers:
  root Area remains `report.md`; descendant Area, Factor, and Requirement paths
  use local structural IDs plus kind suffixes; existing relative-link and report
  ref generation continues to consume those helpers.

- **Refinement**: Updated
  [0108 - Short evaluation report filenames](archive/0108-short-evaluation-report-filenames.md)
  so the root Area remains the run entrypoint at `report.md`, while descendant
  Area, Factor, and Requirement reports use short subject-aware filenames. The
  functional spec now treats old descendant `report.md` paths as the retired
  compatibility shape.

- **Draft**: Created
  [0108 - Short evaluation report filenames](archive/0108-short-evaluation-report-filenames.md)
  with its child
  [functional spec](archive/0108-short-evaluation-report-filenames/spec.md). The case
  keeps the root Area report at `report.md` while replacing repeated descendant
  Evaluation v2 generated Markdown `report.md` filenames with short
  subject-aware filenames such as `<area>-area.md`, `<factor>-factor.md`, and
  `<requirement>-requirement.md`, while preserving structural directories,
  generated report navigation semantics, structured data paths, and completed
  run history. Design and implementation have not started.

- **Done**: Archived
  [0107 - Durable spec alignment](archive/0107-durable-spec-alignment.md).
  The durable spec alignment pass is complete and moved to `changes/archive/`;
  the open-case index is clear again.

- **In-Review**: Advanced
  [0107 - Durable spec alignment](archive/0107-durable-spec-alignment.md) to
  `In-Review`. Audited `SPECIFICATION.md` and 49 active durable spec concepts;
  fixed the routed authoring guide specs that used BCP 14 keywords without a
  convention declaration; added companion-note wording for those runtime guide
  contracts; updated spec logs; and verified with `mise run fmt-md-check` plus
  two follow-up audit passes.

- **Draft**: Created
  [0107 - Durable spec alignment](archive/0107-durable-spec-alignment.md) with
  its child [index](archive/0107-durable-spec-alignment/index.md) and
  [functional spec](archive/0107-durable-spec-alignment/spec.md). The case audits
  `SPECIFICATION.md` and durable `specs/` concepts against the latest
  functional-spec and OKF guidance, fixes concrete current-guidance drift, and
  requires Markdown formatting plus two follow-up audit passes before completion.
  No CLI, Go, bundled skill runtime, format-schema, rating, roll-up,
  evaluation-record, or report-rendering behavior change is expected.

- **Done**: Implemented and archived
  [0106 - Binary confirmation UX](archive/0106-binary-confirmation-ux.md).
  `/quality` interaction guidance now distinguishes non-binary closed choices,
  which keep `1` as the shortest accept path, from true binary mutation gates,
  which show visible `y`/`n` answer paths. Updated durable skill specs, runtime
  skill guidance, and skill/spec logs; the agent-mediated UX guide clarification
  is included in the landed change. No CLI, Go, format-schema, rating, roll-up,
  evaluation-record, report-rendering, or `QUALITY.md` format behavior changed.

- **Design**: Created
  [0106 - Binary confirmation UX](archive/0106-binary-confirmation-ux.md) with
  its child [index](archive/0106-binary-confirmation-ux/index.md),
  [functional spec](archive/0106-binary-confirmation-ux/spec.md), and
  [design doc](archive/0106-binary-confirmation-ux/design.md). The case applies the
  newly clarified agent-mediated UX rule that true binary confirmations use
  visible `y`/`n`, while multi-option choices keep numbered responses. The
  shared UX guide clarification already exists; durable `/quality` skill specs
  and bundled runtime skill guidance remain to be aligned.

- **Done**: Implemented and archived
  [0105 - Evaluation v2 report subject links](archive/0105-evaluation-v2-report-subject-links.md).
  Evaluation v2 Area and Factor report tables now link generated human report
  targets from the row subject cells instead of repeated `Details` columns, while
  preserving explicit data links, report paths, machine-readable outputs, and
  existing summary/status columns. Updated the durable Evaluation v2 report-tree
  spec, reconciled the v2 sketch examples, and added focused report tests.

- **Done**: Implemented and archived
  [0104 - Evaluation v2 report header navigation](archive/0104-evaluation-v2-report-header-navigation.md).
  Evaluation v2 Markdown reports now start with labeled `Area:` trails, Factor
  reports add Factor-only `Factor:` trails, standalone `Breadcrumb:` and parent
  header links are removed, and report headers use compact report-specific
  summary tables. Updated the durable Evaluation v2 report-tree spec, reconciled
  the v2 sketch examples, and added focused report tests.

- **Done**: Archived
  [0100 - Strengthen spec requirement standards (29148 + EARS)](archive/0100-strengthen-spec-requirement-standards.md).
  The guide changes were already implemented and verified; moved the parent case
  and child folder into `changes/archive/` and removed it from the open-case
  index.

- **Design**: Created
  [0105 - Evaluation v2 report subject links](archive/0105-evaluation-v2-report-subject-links.md)
  with its child [index](archive/0105-evaluation-v2-report-subject-links/index.md),
  [functional spec](archive/0105-evaluation-v2-report-subject-links/spec.md),
  and [design doc](archive/0105-evaluation-v2-report-subject-links/design.md).
  The case
  requires Evaluation v2 Area and Factor report tables to move generated human
  report links from repeated `Details` columns into the row subject cells while
  preserving data links, report paths, machine-readable outputs, existing
  summary/status columns, and the separate header-navigation scope of 0104.
  Implementation has not started.

- **Design**: Created
  [0104 - Evaluation v2 report header navigation](archive/0104-evaluation-v2-report-header-navigation.md)
  with its child [index](archive/0104-evaluation-v2-report-header-navigation/index.md),
  [functional spec](archive/0104-evaluation-v2-report-header-navigation/spec.md),
  and [design doc](archive/0104-evaluation-v2-report-header-navigation/design.md).
  The case
  replaces generic Evaluation v2 report `Breadcrumb:` and `Parent` header links
  with labeled `Area:` trails and Factor-only `Factor:` trails, keeps
  Requirement reports free of Factor breadcrumbs, and specifies compact
  report-specific summary headers. Implementation and durable report-spec
  updates have not started.

- **Done**: Created, specced, designed, implemented, and archived
  [0103 - Evaluation v2 report enum display titles](archive/0103-v2-report-enum-display-titles.md).
  Evaluation v2 Markdown reports now render CLI-owned statuses, confidence
  levels, booleans, report kinds, limits, unknowns, and known finding
  classifications through typed human display titles with semantic emoji and
  fallback title-casing for unknown strings. Rating Levels still resolve through
  the run's model snapshot, and routine JSON, `EvaluationOutputResult`, and build
  receipts keep stable raw values. Updated the durable Evaluation v2 report-tree
  spec and focused report tests. Verified with targeted evaluation tests, full Go
  tests, and `mise run fmt-md-check`.

- **Done**: Created, specced, designed, implemented, and archived
  [0102 - Evaluation v2 report rating titles](archive/0102-v2-report-rating-titles.md).
  Evaluation v2 Markdown reports now resolve selected Rating Levels through the
  run's `model.md` Rating Scale `title` values, falling back to stable IDs when
  needed, while routine JSON, `EvaluationOutputResult`, and build receipts keep
  stable Rating Level IDs. Updated the durable Evaluation v2 report-tree spec and
  focused report tests. Verified with targeted evaluation tests, full Go tests,
  and `mise run fmt-md-check`.

- **Done**: Implemented and archived
  [0101 - Quality skill UX action clarity](archive/0101-quality-skill-ux-action-clarity.md);
  set `status` to `Done` and moved the parent concept and child folder into
  [`archive/`](archive/), updating the archive index and removing the entry from
  the bundle-root index. Updated the parent `/quality` interaction contract,
  setup, evaluate, update, and recommendation follow-up durable specs, bundled
  runtime skill guidance, and skill/spec logs so prompts expose explicit CTAs,
  shortest answer paths, numbered ambiguity and outcome choices, decision briefs,
  progress, closeout `Next` fields, and code-span precision. No CLI, Go,
  format-schema, rating, roll-up, evaluation-record, or report behavior changed.
  Verified with `mise run fmt-md-check`.

- **In-Progress**: Advanced
  [0101 - Quality skill UX action clarity](archive/0101-quality-skill-ux-action-clarity.md)
  through `Design` to `In-Progress`; added the
  [design doc](archive/0101-quality-skill-ux-action-clarity/design.md). The
  design patches durable skill specs first, mirrors the concrete prompt shapes
  into runtime guidance, and verifies with targeted searches plus Markdown
  formatting.

- **Draft**: Created
  [0101 - Quality skill UX action clarity](archive/0101-quality-skill-ux-action-clarity.md)
  with its child [index](archive/0101-quality-skill-ux-action-clarity/index.md)
  and [functional spec](archive/0101-quality-skill-ux-action-clarity/spec.md).
  The case
  turns the reviewed agent-mediated UX gaps in `/quality` setup, evaluate,
  update, and recommendation follow-up into implementation requirements:
  explicit checkpoint and prompt CTAs, shortest answer paths, numbered ambiguity
  and outcome choices, decision briefs before `QUALITY.md` writes and external
  issue creation, update run-frame/progress sequencing, result closeout `Next`
  fields, and code-span precision for literal artifacts. It accounts for the
  current agent-mediated UX, functional-spec, and change-case guide updates and
  does not start implementation.

- **In-Review**: Advanced
  [0100 - Strengthen spec requirement standards (29148 + EARS)](archive/0100-strengthen-spec-requirement-standards.md)
  to `In-Review` (no design doc needed). Edited
  [Writing functional specs](../docs/guides/write-functional-specs.md) — added a
  set-level **The requirement set** check (§5.2.6), an **Assumptions &
  dependencies** Shape element (§9.6.8), normative/informational reference
  classification (§9.2.4), an **Unambiguous** quality-bar item plus
  active-voice/subject guidance and a BCP 14-vs-`shall` note (§5.2.4), and an
  optional EARS **statement template**, each with brief provenance — and
  [Working with change cases](../docs/guides/work-with-change-cases.md) — added
  the Draft→Design validation check. No code/CLI/skill/format change. Reconciled
  the Affected artifacts list and updated the bundle [index](index.md). Verified
  with `mise run fmt-md-check`. Not committed or archived.

- **Draft**: Created
  [0100 - Strengthen spec requirement standards (29148 + EARS)](archive/0100-strengthen-spec-requirement-standards.md)
  (`status: Draft`) with its [functional spec](archive/0100-strengthen-spec-requirement-standards/spec.md)
  and [index](archive/0100-strengthen-spec-requirement-standards/index.md). The case
  patches [Writing functional specs](../docs/guides/write-functional-specs.md) and
  [Working with change cases](../docs/guides/work-with-change-cases.md) to close
  requirement-quality gaps found against ISO/IEC/IEEE 29148 — a set-level
  requirement check, a Draft→Design validation gate, an Assumptions &
  dependencies element, normative/informational reference classification, an
  Unambiguous bar item plus language guidance — and to add an optional EARS
  statement template. `SPECIFICATION.md` and the QUALITY.md format are out of
  scope; no code/CLI/skill change. Added the entry to the bundle
  [index](index.md). No durable docs edited yet.

- **Done**: Implemented and archived
  [0099 - Closed-choice setup UX](archive/0099-closed-choice-setup-ux.md);
  set `status` to `Done` and moved the parent concept and child folder into
  [`archive/`](archive/), updating the archive index and removing the entry from
  the bundle-root index. Updated the agent-mediated UX guide, parent `/quality`
  interaction contract, setup runtime guidance, and durable setup workflow specs
  so small closed-choice prompts use numbered options with the recommended
  answer first and `1` as the shortest confirmation. Setup risk discovery now
  presents cost labels and maps them to risk-tolerance values internally. No CLI,
  Go, format-schema, rating, roll-up, evaluation-record, or report behavior
  changed. Verified with `mise run fmt-md-check`.

- **In-Progress**: Advanced
  [0099 - Closed-choice setup UX](archive/0099-closed-choice-setup-ux.md) to
  `In-Progress`; the functional spec and design are settled, and implementation
  is beginning across `docs/guides/agent-mediated-ux.md`, the durable
  `/quality` skill specs, and the bundled runtime skill guidance. No CLI, Go,
  format-schema, rating, roll-up, evaluation-record, or report behavior change
  is expected.

- **Design**: Advanced
  [0099 - Closed-choice setup UX](archive/0099-closed-choice-setup-ux.md) to
  `Design` and added the
  [design doc](archive/0099-closed-choice-setup-ux/design.md). The design
  keeps setup's internal risk-tolerance meaning while presenting cost-labeled
  choices to the user, and updates both durable and runtime skill guidance rather
  than relying on the general UX guide alone.

- **Draft**: Created
  [0099 - Closed-choice setup UX](archive/0099-closed-choice-setup-ux.md) with
  its child [index](archive/0099-closed-choice-setup-ux/index.md) and
  [functional spec](archive/0099-closed-choice-setup-ux/spec.md). The case makes
  `/quality setup` closed-choice discovery prompts use numbered options with the
  recommended answer first and `1` as the shortest confirmation, while mapping
  user-facing cost labels to the existing risk-tolerance model meaning. No CLI,
  Go, format-schema, rating, roll-up, evaluation-record, or report behavior
  change is expected.

- **Done**: Implemented and archived
  [0097 - Evaluation v2 clean break](archive/0097-evaluation-v2-clean-break.md).
  Evaluation v2 is now the only active runtime evaluation path: new runs seed
  `model.md` and `data/`, previous assessment/analysis/recommendation run shapes
  are rejected as unsupported, `evaluation data set` is stdin-only, status/list
  surfaces report v2 data artifacts and gaps, report build writes
  `EvaluationOutputResult` plus the v2 report tree, and old active
  record/report specs and examples were removed. Prepared `v0.13.0` release
  notes, skill compatibility metadata, and the `0.4 (Draft)` specification
  version.

- **Done**: Created, specced, designed, implemented, and archived
  [0098 - Setup opening as first output](archive/0098-setup-opening-first-output.md).
  `/quality setup` now emits its warm welcome, phase roadmap, and run frame as the
  first output before any tool call; the run frame is decoupled from the CLI
  prerequisite check, which becomes a fail-fast gate after the opening and before
  the read-only scan. Updated the runtime
  [`skills/quality/workflows/setup.md`](../skills/quality/workflows/setup.md) and
  the durable
  [`specs/skills/quality-skill/workflows/setup.md`](../specs/skills/quality-skill/workflows/setup.md)
  (Workflow structure + Context-analysis opening requirements, plus a stale
  feedback-log timing preamble), and recorded the entry in the archive
  [index](archive/index.md). Docs/spec/skill guidance only, no CLI/Go or
  format-spec change.

- **In-Progress**: Advanced
  [0097 - Evaluation v2 clean break](0097-evaluation-v2-clean-break.md) to
  `In-Progress`; implementation is beginning across the v2 evaluation CLI,
  runtime graph/status/report paths, durable specs, bundled skill guidance, and
  tests. The case remains scoped to the clean break from legacy
  assessment/analysis/recommendation records and `--file` data-set input.

- **Revision**: Tightened
  [0097 - Evaluation v2 clean break](0097-evaluation-v2-clean-break.md) with an
  explicit `evaluation-v2-sketch.md` reconciliation ledger. The ledger maps each
  sketch heading to preserve, defer, or supersede disposition; the spec now also
  names v2 data inspection/discovery commands and complete-output-plan report
  writes as part of the clean-break punchlist.

- **Revision**: Expanded
  [0097 - Evaluation v2 clean break](0097-evaluation-v2-clean-break.md) after a
  full pass over `evaluation-v2-sketch.md`. Added explicit coverage for prompt
  contracts, shared JSON conventions and refs, synthesis defaults,
  empty-signal policy, evidence-target coverage, criterion results,
  `EvaluationOutputResult` shape, report navigation and empty-state rules,
  resume/failure behavior, settled decisions, and deferred QC/schema/bulk-import
  items so the sketch no longer has untracked active requirements.

- **Design**: Advanced
  [0097 - Evaluation v2 clean break](0097-evaluation-v2-clean-break.md) to
  `Design` and added the
  [design doc](0097-evaluation-v2-clean-break/design.md). The design chooses a
  removal-first implementation sequence: one v2 run/graph validation path,
  stdin-only data persistence, v2-only status/list/report build, explicit
  unsupported old-run diagnostics, and deletion of legacy record/report APIs
  once commands no longer reference them.

- **Draft**: Created
  [0097 - Evaluation v2 clean break](0097-evaluation-v2-clean-break.md) with
  its child [index](0097-evaluation-v2-clean-break/index.md) and
  [functional spec](0097-evaluation-v2-clean-break/spec.md). The case records
  the clean-break punchlist for removing legacy evaluation compatibility,
  changing `qualitymd evaluation data set` to stdin-only input, and closing the
  v2 status/report gaps exposed by the first real generated run.

- **Done**: Landed and archived
  [0096 - Setup intro preview](archive/0096-setup-intro-preview.md); set
  `status` to `Done` and moved the parent concept and child folder into
  [`archive/`](archive/), updating the archive index and removing the entry from
  the bundle-root index. Implemented an educational `/quality setup` opening,
  a project-specific setup preview before discovery questions, and delayed setup
  feedback-log creation until after that preview when the run continues. No CLI,
  Go, format-schema, rating, roll-up, report, or evaluation-record behavior
  changed. Verified with `mise run fmt-md-check`, `git diff --check`, and
  `mise run check`.

- **In-Progress**: Advanced
  [0096 - Setup intro preview](archive/0096-setup-intro-preview.md) through
  `Design` to `In-Progress`; the functional spec and
  [design doc](archive/0096-setup-intro-preview/design.md) are settled, and
  implementation is beginning across bundled setup guidance, durable setup
  workflow specs, workflow logs, and CHANGELOG. No CLI, Go, format-schema,
  rating, roll-up, evaluation-record, or report behavior change is expected.

- **Draft**: Created
  [0096 - Setup intro preview](archive/0096-setup-intro-preview.md) with its
  child [index](archive/0096-setup-intro-preview/index.md) and
  [functional spec](archive/0096-setup-intro-preview/spec.md). The case makes
  `/quality setup` open with short educational orientation, show a
  project-specific setup preview before discovery, and move setup feedback-log
  creation after that preview when the run continues. No CLI, Go,
  format-schema, rating, roll-up, evaluation-record, or report behavior change is
  expected.

## 2026-06-25

- **Done**: Landed and archived
  [0095 - Evaluate feedback log outcomes](archive/0095-evaluate-feedback-log-outcomes.md);
  set `status` to `Done` and moved the parent concept and child folder into
  [`archive/`](archive/), updating the archive index and removing the entry from
  the bundle-root index. Implemented workflow-process outcome values for
  `/quality evaluate` feedback logs across runtime skill guidance, durable skill
  spec mirrors, logs, and CHANGELOG: `outcome` now uses values such as
  `completed-reportable`, `stopped-model`, and `interrupted`, and remains separate
  from report, rating, and recommendation semantics. No CLI, Go, format-schema,
  rating, roll-up, report, or evaluation-record behavior changed. Verified with
  `mise run check`.

- **In-Progress**: Created
  [0095 - Evaluate feedback log outcomes](archive/0095-evaluate-feedback-log-outcomes.md)
  with its child [index](archive/0095-evaluate-feedback-log-outcomes/index.md),
  [functional spec](archive/0095-evaluate-feedback-log-outcomes/spec.md), and
  [design doc](archive/0095-evaluate-feedback-log-outcomes/design.md), then
  advanced it to `In-Progress`. The case keeps `/quality evaluate` feedback
  logging under `.quality/logs/<timestamp>-evaluate-feedback-log.md` and tightens
  `outcome` to workflow-process terminal states such as `completed-reportable`,
  `stopped-model`, and `interrupted`, not report or rating semantics. No CLI, Go,
  format-schema, rating, roll-up, report, or evaluation-record behavior change is
  expected.

- **Done**: Landed and archived
  [0094 - Replace evaluation workflow](archive/0094-replace-evaluation-workflow.md);
  set `status` to `Done` and moved the parent concept and child folder into
  [`archive/`](archive/), updating the archive index and removing the entry from
  the bundle-root index. Implemented the Evaluation v2 protocol/spec bundle,
  structured `data/` routine payloads, `qualitymd evaluation data` commands,
  v2 status and report build support, deterministic Area/Factor/Requirement
  report rendering, superseded legacy report/record specs, and updated bundled
  `/quality` runtime guidance. Verified with `mise run fmt`,
  `mise run fmt-md-check`, `go test ./...`, and `git diff --check`.

- **In-Progress**: Advanced
  [0094 - Replace evaluation workflow](archive/0094-replace-evaluation-workflow.md) from
  `Design` to `In-Progress`; the functional spec and design are drafted, and
  implementation is beginning across durable evaluation specs, CLI command
  specs and implementation, report contracts and rendering, and bundled
  `/quality` skill guidance.

- **Design**: Advanced
  [0094 - Replace evaluation workflow](archive/0094-replace-evaluation-workflow.md) to
  `Design` and added the
  [design doc](archive/0094-replace-evaluation-workflow/design.md). The design records
  the replacement strategy, skill/CLI responsibility split, `data/` run model,
  validation layers, shared status/report gap collection, deterministic report
  rendering, skill runtime shape, orchestration approach, rejected alternatives,
  and risks. No implementation has started.

- **Draft**: Created
  [0094 - Replace evaluation workflow](archive/0094-replace-evaluation-workflow.md)
  with its child
  [index](archive/0094-replace-evaluation-workflow/index.md) and
  [functional spec](archive/0094-replace-evaluation-workflow/spec.md). The case replaces
  the current evaluation workflow, records, report contract, CLI surface, and
  skill runtime with the Evaluation v2 protocol captured in
  `evaluation-v2-sketch.md`: frames before judgment, Requirement assessment and
  rating, bottom-up Factor and Area analysis, structured JSON under `data/`,
  deterministic Area/Factor/Requirement reports, and agent-agnostic
  orchestration. No implementation has started.

- **Done**: Landed and archived
  [0093 - Named Requirement identity](archive/0093-requirement-identity.md);
  set `status` to `Done` and moved the parent concept and child folder into
  [`archive/`](archive/), updating the archive index and removing the entry from
  the bundle-root index. Implemented stable Requirement names, required
  Requirement titles, optional descriptions, retained assessments, qualified
  `requirement:<area-path>::<requirement-name>` references, lint/schema
  enforcement, generated schema updates, scaffolds, docs, examples, and bundled
  skill guidance. This is a breaking format change for legacy statement-key
  Requirements. Verified with `mise run check`.

- **In-Progress**: Advanced
  [0093 - Named Requirement identity](archive/0093-requirement-identity.md) through
  `Design` to `In-Progress`; the functional spec and
  [design doc](archive/0093-requirement-identity/design.md) are settled, and
  implementation is beginning across the format spec, schema, lint/model code,
  scaffolds, bundled skill guidance, docs, and examples. The compatibility
  decision is explicit: legacy statement-key Requirements are rejected rather
  than normalized invisibly.

- **Draft**: Created
  [0093 - Named Requirement identity](archive/0093-requirement-identity.md)
  with its child [index](archive/0093-requirement-identity/index.md) and
  [functional spec](archive/0093-requirement-identity/spec.md). The case gives
  Requirements stable id-like names, required titles, retained assessment text,
  optional descriptions, and qualified Requirement references, while deferring
  compatibility and migration decisions before implementation. No code changes
  yet.

- **Done**: Landed and archived
  [0092 - Setup workflow scope trim](archive/0092-setup-workflow-scope-trim.md);
  set `status` to `Done` and moved the parent concept and child folder into
  [`archive/`](archive/), updating the archive index and removing the entry from
  the bundle-root index. Implemented the setup scope trim across runtime skill
  guidance, durable skill spec mirrors, docs/logs, and CHANGELOG: setup no longer
  asks about future recommendation handling, handoff destination, review cadence,
  recurring review, or automation preferences; setup closeout reports lint
  validation plus important model gaps instead of maturity/evaluation-readiness
  labels; and setup feedback logs record workflow outcomes. No CLI, Go,
  format-schema, rating, roll-up, or evaluation-record behavior changed. Verified
  with `mise run check`.

- **In-Progress**: Advanced
  [0092 - Setup workflow scope trim](archive/0092-setup-workflow-scope-trim.md)
  through `Design` to `In-Progress`; the functional spec and
  [design doc](archive/0092-setup-workflow-scope-trim/design.md) are settled, and
  implementation is beginning across bundled skill setup guidance, Top 10 and
  getting-started guides, setup feedback-log rules, durable skill spec mirrors,
  public docs, logs, and changelog. No CLI, Go, format-schema, rating, roll-up,
  or evaluation-record behavior change is expected.

- **Draft**: Created
  [0092 - Setup workflow scope trim](archive/0092-setup-workflow-scope-trim.md)
  with its child [index](archive/0092-setup-workflow-scope-trim/index.md) and
  [functional spec](archive/0092-setup-workflow-scope-trim/spec.md). The case trims
  `/quality setup` back to first-model authoring: remove future recommendation
  handling and work-handoff destination discovery, remove review cadence and
  recurring automation posture from setup, and replace setup's
  maturity/evaluation-ready closeout with validation plus important model gaps.
  It deliberately keeps the pedagogical setup questions and workflow feedback log.
  No CLI, Go, format-schema, rating, roll-up, or evaluation-record behavior
  change is expected.

## 2026-06-24

- **Done**: Landed and archived
  [0091 — Agent-harness holistic definition](archive/0091-agent-harness-holistic-definition.md);
  set `status` to `Done` and moved the parent concept and child folder into
  [`archive/`](archive/), updating the archive index and removing the entry from
  the bundle-root index. Implemented the holistic harness definition across the
  bundled skill guidance, setup workflow, Top 10 checks, durable skill spec
  mirrors, docs/logs, and CHANGELOG: the agent harness is now defined as the whole
  engineered system around the model, while the agent-harness area is scoped to
  checked-in steering and project-owned control artifacts with an explicit
  mixed-artifact decision rule and checks for owned runtime harness machinery. No
  CLI, Go, format-schema, rating, roll-up, or evaluation behavior changed.
  Verified with `mise run check`.

- **In-Progress**: Advanced
  [0091 — Agent-harness holistic definition](archive/0091-agent-harness-holistic-definition.md)
  through `Design` to `In-Progress`; the functional spec and design are settled,
  and implementation is beginning across the bundled skill authoring guidance,
  setup workflow, Top 10 checks, durable skill spec mirrors, logs, and changelog.
  No CLI, Go, format-schema, rating, roll-up, or evaluation behavior change is
  expected.

- **Draft**: Created
  [0091 — Agent-harness holistic definition](archive/0091-agent-harness-holistic-definition.md)
  with its child
  [index](archive/0091-agent-harness-holistic-definition/index.md),
  [functional spec](archive/0091-agent-harness-holistic-definition/spec.md), and
  [design doc](archive/0091-agent-harness-holistic-definition/design.md). The case
  corrects a part-for-whole error in the `/quality` authoring guidance: the agent
  harness is defined as "the instructions that steer the agent," but the harness is
  the whole engineered system around the model (feedforward guides plus feedback
  sensors, everything-but-the-model). It redefines the harness holistically,
  derives the agent-harness area as one projection with explicit boundaries against
  the Agent Harnessability factor and the tests/operations constituents, adds an
  area-scoping decision rule and a never-define-as-instructions doctrine, brings
  project-owned runtime harness machinery into scope, and aligns `model-structure.md`,
  `setup.md`, and Top 10 check 8 — with durable spec mirrors and logs. Builds on
  [0089](archive/0089-agent-harness-modeling-guidance.md) and
  [0087](archive/0087-encode-projection-boundaries.md) without reopening their
  projection boundary; redefines rather than renames the area. Documentation,
  doctrine, and bundled-skill guidance only; no CLI, Go, format-schema, rating,
  roll-up, or evaluation behavior change expected. Listed it in the bundle
  [index](index.md).

- **Done**: Landed and archived
  [0090 — Skill-content OKF authoring split](archive/0090-skill-content-okf-authoring-split.md);
  set `status` to `Done` and moved the parent concept and child folder into
  [`archive/`](archive/), updating the archive index and removing the entry from
  the bundle-root index. Implemented the runtime `/quality` OKF shape with root,
  guide, workflow, and resource indexes/logs; split `guides/authoring.md` into an
  entry/router plus concern-specific sub-guides; mirrored that authoring tree in
  durable skill specs; updated `SKILL.md`, setup, recommendation follow-up, and
  supporting specs for routed reads; and recorded the change in CHANGELOG. No
  CLI, Go, format-schema, rating, roll-up, or evaluation behavior changed.
  Verified with `mise run check`.

- **In-Progress**: Advanced
  [0090 — Skill-content OKF authoring split](archive/0090-skill-content-okf-authoring-split.md)
  from `Design` to `In-Progress`; spec and design are settled and implementation
  is beginning across the bundled skill guides, mirrored durable specs, reading
  contracts, indexes, logs, and changelog. No code, CLI, format-schema, rating,
  roll-up, or evaluation behavior change is expected.

- **Design**: Created
  [0090 — Skill-content OKF authoring split](archive/0090-skill-content-okf-authoring-split.md)
  with its child
  [index](archive/0090-skill-content-okf-authoring-split/index.md),
  [functional spec](archive/0090-skill-content-okf-authoring-split/spec.md), and
  [design doc](archive/0090-skill-content-okf-authoring-split/design.md). The case splits
  the long `/quality` authoring guide into routed OKF-shaped runtime sub-guides,
  keeps `guides/authoring.md` as the mandatory entry point, mirrors the runtime
  authoring tree in durable sub-specs, starts `skills/quality/` as an OKF-shaped
  runtime bundle, and makes agent read obligations explicit in `SKILL.md` and
  relevant workflows. Documentation, durable specs, and bundled-skill guidance
  structure only; no CLI, Go, format-schema, rating, roll-up, or evaluation
  behavior change expected. Listed it in the bundle [index](index.md).

- **Done**: Landed and archived
  [0089 — Agent-harness modeling guidance](archive/0089-agent-harness-modeling-guidance.md);
  set `status` to `Done` and moved the parent concept and child folder into
  [`archive/`](archive/), updating the archive index and removing the entry from the
  bundle-root index. Implemented documentation and skill-guidance changes only:
  added `continuity` to the current Agent Harnessability decomposition, strengthened
  sub-factor guidance for progressive disclosure, task completion discipline,
  tool affordances, good verification signals, guardrails, containment, and
  trace/run evidence; added a domain-agnostic agent-harness area template for
  steering materials; added the use-context-constituent served-domain guardrail to
  the doctrine guide and `AGENTS.md`; mirrored the contract in durable skill specs;
  extended Top 10 check 8 for thinly factored and software-leaking harness areas;
  and recorded the work in docs/spec logs and `CHANGELOG.md`. No CLI/Go,
  `SPECIFICATION.md`, schema, rating, or roll-up change.

- **Refinement (Design)**: Corrected cross-references and tightened the spec/design
  for
  [0089 — Agent-harness modeling guidance](archive/0089-agent-harness-modeling-guidance.md)
  after a review pass. Repointed the stale `0088` links to `archive/` and fixed the
  "0083 is In-Review" status note (both 0083 and 0088 have since landed and
  archived); reconciled the parent Scope's good-sensor property list with the
  spec/design (fast, actionable, grounded, context-aware, suppressible-not-binary);
  and added spec/design requirements to (a) renumber the now-stale "seventh
  sub-factor / beside the six" language and reconcile the legacy six-sub-factor
  recognition notes once `continuity` is added, (b) update the umbrella
  `agent-harnessability` `description` for the state-preservation capability, (c)
  state the `continuity` boundary against `agent-operability` and
  `agent-accessibility`, (d) define "served domain" on first use and carry the
  guardrail by reference (say-it-once), and (e) keep the self-verifiability
  "suppressible" / enforcement "constrained suppression" boundary explicit.
  Change-case documents only; still `Design`.

- **Cleanup**: Removed the dangling 0078 — View command entry from the bundle
  [index](index.md); its parent concept and child folder were deleted from the
  working tree and the case was never archived, so the open-cases list no longer
  points at missing files.

- **Done**: Landed and archived
  [0088 — Domain-agnostic corpus alignment](archive/0088-domain-agnostic-corpus-alignment.md);
  set `status` to `Done` and moved the parent concept and child folder into
  [`archive/`](archive/), updating the archive index and removing the entry from the
  bundle-root index. Implemented the `0002-city-bike-stations-quality-eval`
  non-software data-product reference fixture with the same reportable runtime
  artifact shape as `0001`; marked and cross-linked the reference example corpus as
  domain-illustrative; added the earned Factors rule to `AGENTS.md`; re-scoped the
  README modeled-domain framing while preserving the agent-first use context; added
  the Appendix B invariance note and lineage clause in `SPECIFICATION.md`; reinforced
  the Top 10, reporting, and report-summary specs; and recorded the work in the docs
  and specs logs plus `CHANGELOG.md`. Reconciled the Affected artifacts list; no CLI
  or Go behavior change.

- **Done**: Landed and archived
  [0083 — Quality-domain agnosticism guide and secondary illustrations](archive/0083-quality-domain-agnosticism.md);
  set `status` to `Done` and moved the parent concept and child folder into
  [`archive/`](archive/), updating the archive index and removing the entry from the
  bundle-root index. Documentation-only change: added the domain-agnostic modeling
  guide, routed agents and contributors to it from `AGENTS.md`, aligned public and
  bundled skill domain-list wording, recorded the durable docs/spec-guide updates,
  and added the CHANGELOG note. No code, format-schema, or CLI behavior changed.

- **Design**: Created
  [0089 — Agent-harness modeling guidance](archive/0089-agent-harness-modeling-guidance.md)
  with its child [index](archive/0089-agent-harness-modeling-guidance/index.md),
  [functional spec](archive/0089-agent-harness-modeling-guidance/spec.md), and
  [design doc](archive/0089-agent-harness-modeling-guidance/design.md). The case fixes an
  asymmetry in the `/quality` authoring guidance: the model-wide Agent
  Harnessability factor is richly specified (six sub-factors with boundaries and
  example requirements) while the agent-harness _area_ — the steering-materials
  constituent the same guidance says to model by default — gets no factor family or
  requirement template, so generated harness areas come out thinly factored.
  Grounded in harness-engineering practice (Fowler's _Harness Engineering_ and
  _Sensors for Coding Agents_; LangChain's _Anatomy of an Agent Harness_), which
  confirms the six sub-factors already map onto the field's harness anatomy. The
  case (a) gives the harness area an illustrative, domain-agnostic factor and
  requirement template at parity with the self-check template, (b) adds a
  served-domain guardrail so harness requirements never assume a software toolchain,
  (c) sharpens `self-verifiability` with good-sensor properties, (d) extends Top 10
  check 8 for thinly-factored and software-leaking harness areas, and (e) adds a
  doctrine principle distinguishing use-context constituents (harness, self-check —
  explicit guidance licensed) from modeled domains (never privileged). Documentation,
  doctrine, and bundled-skill guidance content only; no `SPECIFICATION.md` normative,
  CLI, or Go change. Builds on the 0083 guide and the 0087 projection-boundary rule.
  Listed it in the bundle [index](index.md).

- **Design**: Created
  [0088 — Domain-agnostic corpus alignment](archive/0088-domain-agnostic-corpus-alignment.md)
  with its child [index](archive/0088-domain-agnostic-corpus-alignment/index.md),
  [functional spec](archive/0088-domain-agnostic-corpus-alignment/spec.md), and
  [design doc](archive/0088-domain-agnostic-corpus-alignment/design.md). The case closes the
  residual gaps a multi-agent content audit found against the
  [0083](archive/0083-quality-domain-agnosticism.md) domain-agnosticism guide: add a complete
  non-software worked example (a data product) to the `/quality` reference example
  set, give the software example corpus a domain-illustrative marking and guide
  cross-links, add the factors-earned-per-Model rule to `AGENTS.md`, and re-scope the
  README modeled domain while preserving the agent-first use context. Documentation,
  doctrine, bundled-skill, and spec-example content only; no `SPECIFICATION.md`
  normative, CLI, or Go change. Listed it in the bundle [index](index.md).

- **Done**: Landed and archived
  [0087 — Encode projection boundaries in the model](archive/0087-encode-projection-boundaries.md);
  set `status` to `Done` and moved the parent concept and child folder into
  [`archive/`](archive/), updating the archive index and removing the entry from the
  bundle-root index. Implemented across the bundled authoring guide and Top 10
  checks, their durable spec mirrors, the guides log, and CHANGELOG: the general
  projection-boundary rule lands at the three-projections rule (a YAML comment per
  projection node, plus a disambiguating `description` clause when both projections
  are rated nodes that surface in a report), with the Agent Harnessability factor
  vs. the agent-harness area as the canonical instance and a matching Top 10
  missing-boundary-note check. Reconciled the Affected artifacts list; no
  `SPECIFICATION.md`, CLI schema, or Go code change. Verified with `mise run check`.

- **Design**: Created
  [0087 — Encode projection boundaries in the model](archive/0087-encode-projection-boundaries.md)
  with its child [index](archive/0087-encode-projection-boundaries/index.md),
  [functional spec](archive/0087-encode-projection-boundaries/spec.md), and
  [design doc](archive/0087-encode-projection-boundaries/design.md). The case adds a
  general authoring rule that, when a model carries two or more projections of one
  concern (factor / constituent-area / audience), the boundary is encoded in the
  emitted model — a YAML comment on each node and, when both projections are rated
  nodes that surface in a report, a disambiguating `description` clause — with the
  Agent Harnessability factor vs. the agent-harness area as the canonical instance
  and a matching Top 10 readiness check. Skill-guidance and spec-mirror only; no
  `SPECIFICATION.md`, CLI schema, or Go code change expected. Listed it in the
  bundle [index](index.md).

- **Done**: Landed and archived
  [0086 — Umbrella factor roll-up framing](archive/0086-umbrella-factor-rollup-framing.md);
  set `status` to `Done` and moved the parent concept and child folder into
  [`archive/`](archive/), updating the archive index and removing the entry from
  the bundle-root index. Bundled authoring guide and its spec mirror only; no code,
  format-spec, or CLI change.

- **In-Review**: Created and implemented
  [0086 — Umbrella factor roll-up framing](archive/0086-umbrella-factor-rollup-framing.md)
  in one step. Corrected the bundled
  [authoring guide](../skills/quality/guides/authoring.md) and its
  [spec mirror](../specs/skills/quality-skill/guides/authoring-md.md), which had
  overstated the Agent Harnessability umbrella factor as "do not rate the parent
  directly" / "does not roll up directly"; both now say the umbrella carries no
  requirements of its own and is rated by rolling up its sub-factors, anchored to
  the grouping-area rule in
  [`SPECIFICATION.md`](../SPECIFICATION.md). Recorded the spec-mirror revision in
  the [guides log](../specs/skills/quality-skill/guides/log.md). A full-repo sweep
  confirmed no other live file repeats the framing; the frozen
  `changes/archive/0081-harnessability-factor/{spec.md,design.md}` keep the
  original wording as history. No code, format-spec, or CLI change. Added the
  parent concept, [child folder](archive/0086-umbrella-factor-rollup-framing/index.md), and
  [functional spec](archive/0086-umbrella-factor-rollup-framing/spec.md); updated the
  bundle [index](index.md).

- **Done**: Landed and archived
  [0085 — Agent Harnessability naming](archive/0085-agent-harnessability-naming.md);
  set `status` to `Done` and moved the parent concept and child folder into
  [`archive/`](archive/), updating the archive index and removing the entry from
  the bundle-root index. Skill-guidance, spec-mirror, README, and CHANGELOG only;
  no `SPECIFICATION.md`, CLI schema, or Go code change.

- **In-Review**: Implemented
  [0085 — Agent Harnessability naming](archive/0085-agent-harnessability-naming.md)
  across the bundled authoring guide, setup workflow, Top 10 checks, durable skill
  spec mirrors, guide/workflow logs, README surfaces, and CHANGELOG. New and
  revised models now use `agent-harnessability` / Agent Harnessability with an
  accountability-preserving definition; legacy `harnessability` with the expected
  six sub-factors remains semantic coverage. Historical Change Cases and
  append-only 0081 log entries were left unchanged. Reconciled the affected
  artifacts list and verified with focused `dprint check`, `git diff --check`, and
  stale-wording searches. No `SPECIFICATION.md`, CLI schema, or Go code change.

- **In-Progress**: Advanced
  [0085 — Agent Harnessability naming](archive/0085-agent-harnessability-naming.md)
  from `Design` to `In-Progress`; spec and design are settled and implementation
  is beginning across the bundled skill guidance, durable spec mirrors, README
  surfaces, and CHANGELOG. No `SPECIFICATION.md`, CLI schema, or Go code change
  is expected.

- **Design**: Created
  [0085 — Agent Harnessability naming](archive/0085-agent-harnessability-naming.md)
  with its child [index](archive/0085-agent-harnessability-naming/index.md),
  [functional spec](archive/0085-agent-harnessability-naming/spec.md), and
  [design doc](archive/0085-agent-harnessability-naming/design.md). The case
  renames the 0081 harnessability factor guidance to Agent Harnessability,
  recommends the `agent-harnessability` key for new and revised models, replaces
  attention-scarcity shorthand with human direction/review/accountability wording,
  and keeps legacy `harnessability` models recognizable as semantic coverage.
  Skill-guidance, spec-mirror, README, generated package README verification, and
  CHANGELOG only; no `SPECIFICATION.md`, CLI schema, or Go code change expected.
  Listed it in the bundle [index](index.md).

- **Done**: Landed and archived
  [0084 — Agent-mediated UX conformance](archive/0084-agent-mediated-ux-conformance.md);
  set `status` to `Done` and moved the parent concept and child folder into
  [`archive/`](archive/), updating the archive index and removing the entry from
  the bundle-root index. Skill-guidance, spec-mirror, AGENTS, and CHANGELOG only;
  no `SPECIFICATION.md`, CLI schema, or Go code change. Verified with
  `mise run fmt-md-check`.

- **In-Review**: Implemented
  [0084 — Agent-mediated UX conformance](archive/0084-agent-mediated-ux-conformance.md)
  across the bundled `/quality` skill, runtime setup/evaluate/update/
  recommendation-follow-up guidance, durable skill specs, AGENTS guide routing,
  spec/workflow/guide logs, and CHANGELOG. The shared interaction contract now
  requires status-first agent-mediated output, visually emphasized primary
  questions or calls to action, scannable labels, semantic emoji only, explicit
  mutation boundaries, and closeouts with clear next actions. Reconciled the
  affected-artifact list; `cmd/` and `internal/` were verified no-impact because
  their output is CLI UX rather than agent-mediated workflow output. Verified
  with `mise run fmt-md-check`.

- **In-Progress**: Advanced
  [0084 — Agent-mediated UX conformance](archive/0084-agent-mediated-ux-conformance.md)
  from `Design` to `In-Progress`; spec and design are settled and implementation
  is beginning across the parent `/quality` interaction contract, setup/evaluate/
  update/recommendation-follow-up runtime guidance, durable skill specs, AGENTS
  guide routing, README/install verification, and CHANGELOG.

- **Design**: Advanced
  [0084 — Agent-mediated UX conformance](archive/0084-agent-mediated-ux-conformance.md)
  from `Draft` to `Design` and added its
  [design doc](archive/0084-agent-mediated-ux-conformance/design.md). The design keeps
  the parent `/quality` interaction contract as the shared source of truth,
  applies workflow-specific guidance in setup/evaluate/update/recommendation
  follow-up, and treats CLI output as no-impact unless implementation finds a
  live agent-mediated output path.

- **Draft**: Created
  [0084 — Agent-mediated UX conformance](archive/0084-agent-mediated-ux-conformance.md)
  to bring live agent-mediated workflow guidance and durable skill specs into
  conformance with the new
  [Designing agent-mediated UX](../docs/guides/agent-mediated-ux.md) guide.
  Added the parent concept and child
  [folder](archive/0084-agent-mediated-ux-conformance/index.md) with its
  [functional spec](archive/0084-agent-mediated-ux-conformance/spec.md); no design doc
  yet because the expected work is editorial/spec alignment unless implementation
  discovers a reusable rendering abstraction or code path. Listed it in the
  bundle [index](index.md).

- **In-Review**: Implemented
  [0083 — Quality-domain agnosticism guide and secondary illustrations](archive/0083-quality-domain-agnosticism.md).
  Added the `docs/guides/model-quality-across-domains.md` contributor guide (stress
  axes, the canonical secondary-domain set, range-finder illustrations, and a full
  worked documentation example), registered it in the guides index and the
  `AGENTS.md` Guides table with a doctrine pointer, aligned the four domain
  enumerations (`SPECIFICATION.md` Lineage, `AGENTS.md`, the authoring guide, and
  the setup workflow), added the README secondary-domain note, recorded the
  skill-guide alignment in the [guides log](../specs/skills/quality-skill/guides/log.md),
  and added the CHANGELOG note. Documentation-only; no `SPECIFICATION.md` normative
  or code change. Reconciled the Affected artifacts list; the worked example lints
  clean and `mise run fmt-md-check` passes.
- **Creation**: Opened
  [0083 — Quality-domain agnosticism guide and secondary illustrations](archive/0083-quality-domain-agnosticism.md)
  (`status: Draft`) to make the project's domain agnosticism demonstrated rather
  than only asserted: a contributor-doctrine guide
  (`docs/guides/model-quality-across-domains.md`) defining the stress axes and a
  canonical secondary-domain set, one full worked non-software example, and
  alignment of the four drifting domain enumerations. Added the parent concept and
  child [folder](archive/0083-quality-domain-agnosticism/index.md) with its
  [functional spec](archive/0083-quality-domain-agnosticism/spec.md); no design doc
  (editorial change). Listed it in the bundle [index](index.md).
- **Done**: Landed and archived
  [0082 — Normalize QUALITY.md self-check roll-up](archive/0082-normalize-quality-md-rollup.md);
  set `status` to `Done` and moved the parent concept and child folder into
  [`archive/`](archive/), updating the archive index and removing the entry from
  the bundle-root index. Skill-guidance, spec-mirror, and CHANGELOG only, with
  `QUALITY.md` assessed no-change; no `SPECIFICATION.md`, CLI schema, or Go code
  change.
- **In-Review**: Implemented
  [0082 — Normalize QUALITY.md self-check roll-up](archive/0082-normalize-quality-md-rollup.md)
  across the bundled authoring guide and evaluate workflow, durable authoring and
  evaluation spec mirrors, guide-spec log, and CHANGELOG. `quality-md` is now an
  ordinary in-scope area for assessment, analysis, reporting, and roll-up, while
  quality-log writing remains limited to meaningful confirmed model changes.
  Reconciled the Affected artifacts list and verified with
  `qualitymd lint QUALITY.md --json`, `qualitymd status QUALITY.md --json`, and
  `mise run fmt-md-check`.
- **In-Progress**: Advanced
  [0082 — Normalize QUALITY.md self-check roll-up](archive/0082-normalize-quality-md-rollup.md)
  from `Design` to `In-Progress`; spec and design are settled and implementation
  is beginning on the runtime authoring/evaluation guidance, durable spec
  mirrors, dogfooded `QUALITY.md` wording, and CHANGELOG. Skill/spec/docs only;
  no `SPECIFICATION.md`, CLI schema, or Go code change is expected.
- **Done**: Landed and archived
  [0081 — Harnessability factor](archive/0081-harnessability-factor.md); set
  `status` to `Done` and moved the parent concept and child folder into
  [`archive/`](archive/), updating the archive index and removing the entry from
  the bundle-root index. Skill-guidance, spec-mirror, README, and CHANGELOG only;
  no `SPECIFICATION.md` or CLI/Go change. Verified with `mise run check`.
- **In-Review**: Implemented
  [0081 — Harnessability factor](archive/0081-harnessability-factor.md) across the
  bundled skill (authoring guide, setup workflow, Top 10 checks), durable spec
  mirrors (authoring-md, setup, top-10-quality-md-checks-md, plus guide/workflow
  logs), README, and CHANGELOG. Harnessability is now a model-wide umbrella factor
  proposed by default for agent-collaborated composite roots, decomposed into six
  sub-factors, kept distinct from the agent-harness constituent, and routed by the
  Top 10 check when missing. Reconciled the Affected artifacts list and verified
  with `mise run check`.
- **Design**: Created
  [0082 — Normalize QUALITY.md self-check roll-up](archive/0082-normalize-quality-md-rollup.md)
  with its child [index](archive/0082-normalize-quality-md-rollup/index.md),
  [functional spec](archive/0082-normalize-quality-md-rollup/spec.md), and
  [design doc](archive/0082-normalize-quality-md-rollup/design.md). The case removes the
  `quality-md` self-check's special out-of-roll-up treatment and makes it an
  ordinary in-scope area for evaluation, analysis, reporting, and roll-up while
  preserving quality-log behavior for meaningful model changes. Skill-guidance
  and spec-mirror only; no `SPECIFICATION.md`, CLI schema, or code change.
- **Refinement (Draft)**: Generalized
  [0081 — Harnessability factor](0081-harnessability-factor.md)'s sub-factors to be
  quality-domain agnostic and renamed `operational-readiness` to
  **agent-operability**. The operate sub-factor no longer assumes a software
  development environment — it is the agent establishing and operating its working
  environment from recorded materials (a dev environment, a budgeting project's
  connected accounts, or a legal case's loaded matter) — and the rename names the
  quality (operability) while completing the agent-/self-scoped pattern with
  `agent-accessibility` and `self-verifiability`. Held the other five definitions
  domain-neutral with software as one illustrative instance, and broadened
  `containment-of-action` with a consequential real-world action example. Updated
  the [spec](0081-harnessability-factor/spec.md) and
  [design doc](0081-harnessability-factor/design.md) only; still `Draft`.
- **Draft**: Created [0081 — Harnessability factor](0081-harnessability-factor.md)
  to add harnessability — how well a project equips an agent to work on it — as a
  model-wide umbrella factor with six sub-factors (agent-accessibility,
  task-specifiability, operational-readiness, self-verifiability,
  enforcement-of-standards, containment-of-action), proposed by default for an
  agent-collaborated composite root. Gives the agent-collaboration concern its
  factor projection alongside the agent-harness constituent (0080); relates to 0077. Added the parent concept and child folder
  ([index](0081-harnessability-factor/index.md)), and registered it in the
  bundle-root [index](index.md). Skill-guidance, spec-mirror, and README only; no
  `SPECIFICATION.md` or code change.
- **Done**: Landed and archived
  [0080 - Model constituents by default](archive/0080-model-constituents-by-default.md);
  set `status` to `Done` and moved the parent concept and its child folder into
  [`archive/`](archive/), updating the archive index and removing the entry from
  the bundle-root index. Skill-guidance and spec-mirror only; no `SPECIFICATION.md`
  or code change.
- **In-Review → In-Progress**: Implemented 0080 across the bundled skill (authoring
  guide, setup workflow, Top 10 checks, getting-started) and its spec mirror
  (authoring-md, setup, top-10-quality-md-checks-md, with 0080 rationales), plus the
  guide/workflow spec logs and the CHANGELOG. Constituent coverage is now
  model-by-default with two disqualifiers, a no-silent-omission rule, demoted
  deferral, and a maturity completeness bar. Reconciled the Affected artifacts list;
  verified with `mise run check` and `qualitymd lint`.
- **Design**: Authored the [0080](0080-model-constituents-by-default.md) functional
  [spec](0080-model-constituents-by-default/spec.md) and
  [design doc](0080-model-constituents-by-default/design.md): model-by-default with
  two disqualifiers (no-distinct-concerns → fold; not-germane → out of scope), the
  no-silent-omission rule with (a) minimal-area vs (b) requirement-elsewhere
  routing, deferral demoted to a blocker-only exception, and a first-pass
  completeness bar. Recorded the SPECIFICATION.md exclusion and the dogfood
  re-check as a follow-up.
- **Draft**: Created [0080 - Model constituents by default](0080-model-constituents-by-default.md)
  to flip the skill's constituent-coverage guidance from earn-it/defer-freely to
  model-by-default with a short "don't model" list and a no-silent-omission rule.
  Added the parent concept and child folder ([index](0080-model-constituents-by-default/index.md)),
  and registered it in the bundle-root [index](index.md). Skill-guidance and
  spec-mirror only; no `SPECIFICATION.md` or code change.
- **Done**: Landed and archived
  [0079 - Stewardship vocabulary discipline](archive/0079-stewardship-vocabulary-discipline.md);
  set `status` to `Done` and moved the parent concept and its child folder into
  [`archive/`](archive/), updating the archive index and removing the entry from
  the bundle-root index. Documentation-only; no `SPECIFICATION.md` or code change.
- **In-Review**: Completed implementation of
  [0079 - Stewardship vocabulary discipline](0079-stewardship-vocabulary-discipline.md).
  Added the AGENTS.md "Keep the motivation and taxonomy registers distinct"
  subsection; rephrased the two authoring-guide "stewardship lenses" fusions and
  added an Avoid bullet at the three-projections rule; added the register clause
  and a 0079 rationale to the durable authoring-guide spec; added an operational
  setup-workflow guard; and recorded the guides-log and CHANGELOG notes. The
  durable setup spec was assessed no-change. Documentation-only; reconciled the
  Affected artifacts list and verified with `mise run check` (markdown format,
  bundle link resolution, Go vet/lint/test all pass). Not committed, not archived.
- **In-Progress**: Advanced
  [0079 - Stewardship vocabulary discipline](0079-stewardship-vocabulary-discipline.md)
  from `Design` to `In-Progress`; spec and design are settled and implementation is
  beginning on the AGENTS.md register rule, the two authoring-guide "stewardship
  lenses" rephrasings plus a three-projections Avoid bullet, the durable
  authoring-guide spec clause, and an operational setup guard. Documentation-only;
  the durable setup spec is assessed no-change.
- **Design**: Added the
  [design doc](0079-stewardship-vocabulary-discipline/design.md) for
  [0079 - Stewardship vocabulary discipline](0079-stewardship-vocabulary-discipline.md)
  and moved it from `Draft` to `Design`. Records the register rule's three homes
  (AGENTS.md canonical; the authoring guide and its durable spec operational), the
  exact rephrasing of the two "stewardship lenses" fusions, the decision to leave
  the singular "quality lens" gloss intact, and the conclusion that the setup guard
  is operational rather than a new spec contract.
- **Creation**: Added
  [0079 - Stewardship vocabulary discipline](0079-stewardship-vocabulary-discipline.md)
  (`status: Draft`) with its child [index](0079-stewardship-vocabulary-discipline/index.md)
  and [functional spec](0079-stewardship-vocabulary-discipline/spec.md). The case
  confines the stewardship/care core language to its motivation register so it
  never modifies or replaces taxonomy nouns (factor, area, requirement,
  constituent, audience): it removes the two "stewardship lenses" fusions in the
  authoring guide, guards the setup model summary against "stewardship factors,"
  and records the register rule in the guide, its durable spec, and `AGENTS.md`.
  Documentation-only; no format, schema, CLI, or evaluation change, and no
  retraction of the 0076/0077 stewardship grounding.
- **Done**: Landed
  [0077 - Care-grounded stewardship concerns](archive/0077-stewardship-care-grounding.md)
  and moved the parent concept and its child folder into
  [`archive/`](archive/index.md). Removed it from the bundle [index](index.md) and
  recorded it in the archive [index](archive/index.md). Care-grounding of the
  stewardship-concern generator; documentation-only, no format or code change.
- **In-Review**: Completed implementation of
  [0077 - Care-grounded stewardship concerns](archive/0077-stewardship-care-grounding.md).
  Reframed the "Cover the domain's constituent kinds" authoring subsection
  (steward-as-holding-in-trust, claim-from-a-Need-or-Risk, artifact-as-evidence,
  recurring lifecycle band, vulnerability-grounded protective pair, and the new
  constituent-vs-stewardship-quality "Do" bullet), aligned the durable authoring
  guide spec and promoted the rationale, updated the guides log, and added the
  CHANGELOG note. Documentation-only; reconciled the Affected artifacts list (setup
  workflow and Top 10 checks assessed as no-change) and verified with
  `mise run check`.
- **In-Progress**: Advanced
  [0077 - Care-grounded stewardship concerns](archive/0077-stewardship-care-grounding.md)
  from `Design` to `In-Progress`; spec and design are settled and implementation
  is beginning on the "Cover the domain's constituent kinds" authoring subsection
  and its durable guide spec.
- **Creation**: Added [0078 - View command](0078-view-command.md) (`status: Draft`)
  with its child [index](0078-view-command/index.md) and
  [functional spec](0078-view-command/spec.md). The case adds `qualitymd view`, a
  read-only presentation surface structured along two axes — a lens (document,
  `outline`, later ratings/coverage/trends/recommendations) and a surface
  (`text`/`json`/`mermaid`, later `dot`/`html`/`--serve`) — over shared
  deterministic workspace data. The first slice builds the default document render
  and the `outline` lens on text/json/mermaid, and wires `/quality setup` closeout
  to render the outline. Motivated by setup having no way to show the model it just
  authored, and by `status` being a routing snapshot rather than a presentation
  home. Design doc and implementation not started. Listed it in the bundle
  [index](index.md).
- **Advance**: Moved
  [0077 - Care-grounded stewardship concerns](archive/0077-stewardship-care-grounding.md)
  to `Design` and authored its
  [design doc](archive/0077-stewardship-care-grounding/design.md). The design settles the
  five care refinements (artifact-as-evidence, claim-from-a-Need-or-Risk,
  protective-pair-under-vulnerability, steward-as-holding-in-trust, recurring
  lifecycle) and the new constituent-vs-stewardship-quality "Do" bullet, maps each
  to where it lands in the authoring subsection, and records why the phenomenology
  is taken as framing rather than a category-theory formalism. Implementation not
  started.
- **Creation**: Added
  [0077 - Care-grounded stewardship concerns](archive/0077-stewardship-care-grounding.md)
  (`status: Draft`) with its child
  [index](archive/0077-stewardship-care-grounding/index.md) and
  [functional spec](archive/0077-stewardship-care-grounding/spec.md). The case refines
  [0076](archive/0076-domain-constituent-kinds.md)'s stewardship-concern generator
  with the phenomenology of care (Rousse & Spivak): the artifact is the _trace_ of
  tending rather than the care itself, the claim that earns a constituent comes
  from a Need or Risk rather than the list, the protective pair answers to
  vulnerability, stewardship is defined as holding-in-trust, and the lifecycle band
  is recurring rather than a once-through pipeline. No format, schema, CLI, or
  evaluation change. Listed it in the bundle [index](index.md). Design and
  implementation not started.
- **Creation**: Added
  [0076 - Domain constituent kinds and stewardship concerns](archive/0076-domain-constituent-kinds.md)
  (`status: Draft`) with its child [index](archive/0076-domain-constituent-kinds/index.md).
  The case closes the domain-constituent coverage gap 0074 left open: it teaches
  the authoring guidance to enumerate a composite root's constituents by
  stewardship concerns (a lifecycle band plus the secure/safeguard protective
  pair) and an audience×purpose axis, with a three-projections rule that keeps
  factors, constituents, and audiences from double-counting. Design and
  implementation not started. Listed it in the bundle [index](index.md).
- **Advance**: Moved
  [0076 - Domain constituent kinds and stewardship concerns](archive/0076-domain-constituent-kinds.md)
  to `Design` and authored its
  [design doc](archive/0076-domain-constituent-kinds/design.md). The design settles the
  two-axis taxonomy (stewardship concerns × audience/purpose), the secure vs.
  safeguard direction-of-harm split, the three-projections rule, and where each
  edit lands across the authoring guide, setup workflow, and Top 10 checks.
  Implementation not started.
- **In-Progress**: Advanced
  [0076 - Domain constituent kinds and stewardship concerns](archive/0076-domain-constituent-kinds.md)
  from `Design` to `In-Progress`; spec and design are settled and implementation
  is beginning across the authoring guide, the setup workflow, the Top 10 checks,
  and the durable guide specs.
- **In-Review**: Completed implementation of
  [0076 - Domain constituent kinds and stewardship concerns](archive/0076-domain-constituent-kinds.md).
  Added the "Cover the domain's constituent kinds" authoring subsection and the
  use-context cross-reference, the composite-root constituent-enumeration step in
  the setup workflow, the Top 10 Check 8 missing-domain-constituent finding and
  the condensed-checklist line, aligned the authoring and top-10 guide specs and
  the guides log, and added the CHANGELOG note. Documentation-only; reconciled the
  Affected artifacts list and verified with `mise run check`.
- **Done**: Landed
  [0076 - Domain constituent kinds and stewardship concerns](archive/0076-domain-constituent-kinds.md)
  and moved the parent concept and its child folder into
  [`archive/`](archive/index.md). Removed it from the bundle [index](index.md)
  and recorded it in the archive [index](archive/index.md).

- **Creation**: Added
  [0075 - Rating title emoji defaults](archive/0075-rating-title-emoji-defaults.md)
  (`status: Draft`) with its child
  [spec](archive/0075-rating-title-emoji-defaults/spec.md). The case makes
  emoji-prefixed Rating Level titles the default starter and `/quality setup`
  display convention while preserving plain stable `level` identifiers and
  keeping the formal format neutral. Design and implementation not started.
  Listed it in the bundle [index](index.md).
- **Advance**: Moved
  [0075 - Rating title emoji defaults](archive/0075-rating-title-emoji-defaults.md) to
  `Design` and authored its
  [design doc](archive/0075-rating-title-emoji-defaults/design.md). The design
  treats emoji as a Rating Level title display convention, not rating semantics:
  scaffold and setup defaults change, stable `level` IDs stay plain, reports
  continue to resolve model titles, and the formal format stays neutral.
  Implementation not started.
- **In-Progress**: Advanced
  [0075 - Rating title emoji defaults](archive/0075-rating-title-emoji-defaults.md)
  from `Design` to `In-Progress`; spec and design are settled and implementation
  is beginning across scaffold templates, setup guidance, authoring guidance, and
  release notes.
- **In-Review**: Completed implementation of
  [0075 - Rating title emoji defaults](archive/0075-rating-title-emoji-defaults.md).
  Updated both `qualitymd init` scaffolds and tests, setup runtime guidance, the
  authoring and getting-started guides, durable setup/authoring guide specs and
  logs, the README example, and release notes. Left `SPECIFICATION.md` unchanged
  so the formal format remains neutral. Verified with
  `go test ./internal/scaffold` and `mise run check`.
- **Implementation + Archival**: Advanced
  [0075 - Rating title emoji defaults](archive/0075-rating-title-emoji-defaults.md)
  through `In-Review` to `Done`, moving it (parent and folder) into
  [`archive/`](archive/). The landed behavior makes emoji-prefixed Rating Level
  titles the default for CLI scaffolds and `/quality setup`, keeps stable IDs
  plain, updates README and skill guidance/specs/logs, and leaves
  `SPECIFICATION.md` neutral. Updated the bundle [index](index.md) and the
  [archive index](archive/index.md).

- **Creation**: Added
  [0074 - Composite root areas and use-context constituents](0074-composite-root-areas.md)
  (`status: Draft`) with its child [spec](0074-composite-root-areas/spec.md) and
  [design](0074-composite-root-areas/design.md). The case names three root shapes
  (primary-subject, collection, composite) in the authoring guidance, separates
  domain constituents from recurring use-context constituents (the agent harness
  and the QUALITY.md self-check), states the self-check's out-of-roll-up
  learn-loop axis and the harness's dual normative role, and corrects the
  root-factor-count heuristic to apply per primary subject / per constituent.
  Documentation-only; no `SPECIFICATION.md` or code change. Listed it in the
  bundle [index](index.md).
- **Refinement (Draft)**: Captured that the three shapes are recursive and
  composable, not mutually exclusive root choices — a composite area may hold
  collection or composite children to any depth (e.g. a composite root with
  harness and quality-md constituents alongside an `apps` collection of
  `apps/product1` and `apps/product2`). Updated
  [0074](archive/0074-composite-root-areas.md) scope, added spec requirements
  R17–R18, and added a worked example to the
  [design doc](archive/0074-composite-root-areas/design.md).
- **Implementation + Archival**: Advanced
  [0074 - Composite root areas and use-context constituents](archive/0074-composite-root-areas.md)
  through `In-Progress` and to `Done`, moving it (parent and folder) into
  [`archive/`](archive/). Added the "Choose the decomposition shape" and "Carry
  the recurring use-context constituents" subsections to the authoring guide,
  scoped the factor-coverage aim per primary-subject node, extended Top 10 Check 8
  with the composite-flattening and missing-constituent findings, aligned both
  guide specs and the [guides log](../specs/skills/quality-skill/guides/log.md),
  and added the CHANGELOG note. Documentation-only; no `SPECIFICATION.md` or code
  change. Verified with `mise run check`. Updated the bundle [index](index.md) and
  the [archive index](archive/index.md).

## 2026-06-23

- **Implementation + Archival**: Implemented and advanced
  [0073 - Evaluation feedback log](archive/0073-evaluation-feedback-log.md)
  through `In-Progress` and `In-Review` to `Done`, moving it (parent and folder)
  into [`archive/`](archive/). Added the shared workflow feedback-log durable
  spec, the evaluate adopter spec, runtime evaluate feedback-log guidance under
  `.quality/logs/<timestamp>-evaluate-feedback-log.md`, CLI scaffold changes so
  new evaluation runs no longer seed `debug-log.md`, legacy compatibility docs,
  example updates, tests, and release notes. Verified with `mise run check`.

- **In-Review**: Completed implementation of
  [0073 - Evaluation feedback log](archive/0073-evaluation-feedback-log.md). Added the
  shared workflow feedback-log spec and evaluate adopter spec, updated setup to
  reference the shared contract, moved current evaluate workflow feedback to
  `.quality/logs/<timestamp>-evaluate-feedback-log.md`, removed `debug-log.md`
  from new evaluation run scaffolds, updated specs/runtime skill/examples/release
  notes, and verified with `mise run check`.

- **In-Progress**: Advanced
  [0073 - Evaluation feedback log](archive/0073-evaluation-feedback-log.md) from
  `Design` to `In-Progress`; spec and design are settled and implementation is
  beginning across durable specs, runtime skill instructions, CLI scaffold tests,
  examples, and release notes.

- **Advance**: Moved
  [0073 - Evaluation feedback log](archive/0073-evaluation-feedback-log.md) to `Design`
  and authored its
  [design doc](archive/0073-evaluation-feedback-log/design.md). The design makes
  `.quality/logs/<timestamp>-evaluate-feedback-log.md` the current evaluation
  workflow feedback artifact, extracts a shared feedback-log durable spec,
  treats setup and evaluate as adopters, stops new evaluation scaffolds from
  seeding `debug-log.md`, and keeps historical `debug-log.md` artifacts
  compatible. Implementation not started.

- **Creation**: Added
  [0073 - Evaluation feedback log](archive/0073-evaluation-feedback-log.md) at `Draft`
  with its child folder
  ([index](archive/0073-evaluation-feedback-log/index.md),
  [functional spec](archive/0073-evaluation-feedback-log/spec.md)). The case aligns
  `/quality evaluate` with setup feedback logging by introducing a shared
  workflow feedback-log durable spec, adding an evaluate-specific feedback-log
  spec under the evaluate workflow, moving new evaluate feedback to
  `.quality/logs/<timestamp>-evaluate-feedback-log.md`, and treating
  `debug-log.md` as historical compatibility for old evaluation runs. Design and
  implementation not started.

- **Creation + Implementation + Archival**: Created
  [0072 - Setup context checkpoint](archive/0072-setup-context-checkpoint.md)
  with its child folder
  ([index](archive/0072-setup-context-checkpoint/index.md),
  [functional spec](archive/0072-setup-context-checkpoint/spec.md),
  [design doc](archive/0072-setup-context-checkpoint/design.md)), advanced it
  through `Draft`, `Design`, `In-Progress`, and `In-Review` to `Done`, and
  archived it. Implemented skill-only with **no CLI/Go change**: setup now
  presents primary users/outcomes, maintainers/collaborators, other stakeholders,
  and missing or not-agent-accessible context as a compact correction-oriented
  checkpoint instead of four separate open-ended prompts, preserving provenance
  and recording omitted low-confidence gaps honestly. Synced the durable setup
  workflow spec, workflow/spec logs, runtime setup workflow, `CHANGELOG.md`, and
  archive index.

- **Creation + Implementation + Archival**: Created
  [0071 - Setup open-ended review gate](archive/0071-setup-open-ended-review-gate.md)
  with its child folder
  ([index](archive/0071-setup-open-ended-review-gate/index.md),
  [functional spec](archive/0071-setup-open-ended-review-gate/spec.md),
  [design doc](archive/0071-setup-open-ended-review-gate/design.md)), advanced
  it through `Draft`, `Design`, `In-Progress`, and `In-Review` to `Done`, and
  archived it. Implemented skill-only with **no CLI/Go change**: setup's final
  review recap now uses friendly, open-ended wording that preserves the
  `"looks good"` fast path while inviting priorities, worries, wording, edge
  cases, repo-invisible context, and other useful last-call input before
  authoring. Synced the durable setup workflow spec, workflow/spec logs, runtime
  setup workflow, `CHANGELOG.md`, and archive index.

- **Creation + Implementation + Archival**: Created
  [0070 - Setup missing-context provenance](archive/0070-setup-missing-context-provenance.md)
  with its child folder
  ([index](archive/0070-setup-missing-context-provenance/index.md),
  [functional spec](archive/0070-setup-missing-context-provenance/spec.md),
  [design doc](archive/0070-setup-missing-context-provenance/design.md)),
  advanced it through `Draft`, `Design`, `In-Progress`, and `In-Review` to
  `Done`, and archived it. Implemented skill-only with **no CLI/Go change**:
  setup missing-context discovery now treats material context as
  agent-accessible only when supported by repository/tool/source evidence or
  explicit setup-provided context, prohibits choices that assume low/no-evidence
  project facts are sufficiently understood, and preserves setup-provided
  provenance in authored `QUALITY.md` body context. Synced the durable setup
  workflow spec, workflow/spec logs, runtime setup workflow, `CHANGELOG.md`, and
  archive index.

- **Implementation + Archival**: Implemented and advanced
  [0069 - Setup review gate and discovery trim](archive/0069-setup-review-gate-and-pedagogy-trim.md)
  through `Design`, `In-Progress`, and `In-Review` to `Done`, moving it (parent
  and folder) into [`archive/`](archive/). Added the
  [design doc](archive/0069-setup-review-gate-and-pedagogy-trim/design.md), then
  implemented skill-only with **no CLI/Go change**: setup discovery now asks nine
  questions, removes modeling rigor and review posture as user-facing discovery
  questions, adds a Rating Scale confirmation question, trims per-question
  pedagogy to purpose/context only, and treats the final recap as a hard review
  gate before authoring. Synced the runtime setup workflow, durable setup spec,
  spec logs, and `CHANGELOG.md`. Removed it from the open-cases [index](index.md)
  and added it to the [archive index](archive/index.md).

- **Update**: Amended
  [0069 - Setup review gate and discovery trim](0069-setup-review-gate-and-pedagogy-trim.md)
  while still in `Draft` to add a rating-scale confirmation question. The new
  question teaches that Rating Levels are configurable model vocabulary, not
  baked into QUALITY.md, while recommending the standard
  `outstanding`/`target`/`minimum`/`unacceptable` scale and explaining the
  decision role of each level. Setup still must not ask the user to invent
  custom Rating Level names during discovery.

- **Implementation + Archival**: Implemented and advanced
  [0068 - Always-on setup feedback log](archive/0068-always-on-setup-feedback-log.md)
  through `In-Progress` and `In-Review` to `Done`, moving it (parent and folder)
  into [`archive/`](archive/). Implemented skill-only with **no CLI/Go change**:
  setup now creates the current run's feedback log during preflight after CLI
  support and the run frame, updates the current file for material
  workflow-experience events, and finalizes it at close with stable frontmatter
  metadata, lifecycle status, a timeline, and explicit no-notable-content notes.
  Synced the durable workflow feedback-log sub-spec, setup workflow spec, parent
  `/quality` skill spec, runtime skill files, CLI quick reference, spec logs, and
  `CHANGELOG.md`; no public durable docs changed. Removed it from the open-cases
  [index](index.md) and added it to the [archive index](archive/index.md).

- **In-Review**: Completed implementation of
  [0068 - Always-on setup feedback log](0068-always-on-setup-feedback-log.md).
  Implemented skill-only with **no CLI/Go change**: setup now creates the
  current run's feedback log during preflight after CLI support and the run
  frame, updates the current file for material workflow-experience events, and
  finalizes it at close. Synced the durable workflow feedback-log sub-spec,
  setup workflow spec, parent `/quality` skill spec, runtime skill files, CLI
  quick reference, spec logs, and `CHANGELOG.md`. No public durable docs changed.

- **In-Progress**: Advanced
  [0068 - Always-on setup feedback log](0068-always-on-setup-feedback-log.md)
  from `Design` to `In-Progress`; spec and design were settled and
  implementation began across durable specs, the runtime skill, and release
  notes.

- **Update**: Amended
  [0069 - Setup review gate and discovery trim](0069-setup-review-gate-and-pedagogy-trim.md)
  while still in `Draft` to include removing the modeling-rigor and
  review-posture discovery questions. Modeling rigor may remain an internal
  setup-brief inference, and review/loop expectations move to setup closeout
  next-step routing rather than discovery. No replacement question added.

- **Creation**: Added
  [0069 - Setup review gate and discovery trim](0069-setup-review-gate-and-pedagogy-trim.md)
  at `Draft` with its child folder
  ([index](0069-setup-review-gate-and-pedagogy-trim/index.md),
  [functional spec](0069-setup-review-gate-and-pedagogy-trim/spec.md)). The case
  makes `/quality setup` stop after discovery, present the final recap, and wait
  for an explicit user response before authoring `QUALITY.md`; structured
  question-tool completion does not satisfy that review gate. It also trims
  per-question teaching copy to purpose/context only, removing repeated
  "how to change it later" guidance while allowing one general living-document
  note. No CLI/Go change expected; skill + durable setup spec + changelog only.
  Design and implementation not started. Listed it under open cases in the
  bundle [index](index.md).

- **Advance**: Moved
  [0068 - Always-on setup feedback log](0068-always-on-setup-feedback-log.md) to
  `Design` and authored its
  [design doc](0068-always-on-setup-feedback-log/design.md). The design keeps the
  artifact skill-only and local, emits the run frame before the first feedback-log
  write, creates `.quality/logs/<started-at>-setup-feedback-log.md` immediately
  after preflight has CLI/model metadata, updates the current run's file in place
  for material workflow-experience events, and finalizes it at close. It records
  the current-run overwrite boundary, the timeline/body split, stop-handling edge
  cases, and the rejected alternatives of close-only creation, append-only events,
  a JSONL sidecar, and a CLI helper. Implementation not started.

- **Creation**: Added
  [0068 - Always-on setup feedback log](0068-always-on-setup-feedback-log.md) at
  `Draft` with its child folder
  ([index](0068-always-on-setup-feedback-log/index.md),
  [functional spec](0068-always-on-setup-feedback-log/spec.md)). The case changes
  `/quality setup` feedback logging from optional close-step authoring to an
  always-created run artifact under `.quality/logs/` that is created during
  preflight, updated as the workflow progresses, and finalized at close with
  stable frontmatter metadata and body sections. It remains skill-only, local,
  never transmitted, and bounded to workflow-experience feedback rather than
  `QUALITY.md` model rationale or evaluation records. Design and implementation
  not started. Listed it under open cases in the bundle [index](index.md).

- **Implementation + Archival**: Implemented and advanced
  [0067 - Setup discovery pedagogy](archive/0067-setup-discovery-pedagogy.md)
  through `In-Progress` and `In-Review` to `Done`, moving it (parent and folder)
  into [`archive/`](archive/). Implemented skill-only with **no code change**:
  added authored per-question background and how-to-change-later copy inline in
  the runtime [setup workflow](../skills/quality/workflows/setup.md) with a
  teaching-over-round-trips framing, made setup ask every one of the ten questions
  every run (removed the accept-all-and-skip escape; kept a per-question fast
  confirm and show-all-at-once), relabeled confidence from
  `strongly inferred`/`weakly inferred`/`assumed` to `Low`/`Med`/`High` (evidence
  note retained, no-evidence → `Low`), and added a final review recap before
  authoring. Synced the durable
  [setup spec](../specs/skills/quality-skill/workflows/setup.md) with promoted
  rationale annotations (the parent skill spec's generic "confidence-labeled
  defaults" phrasing reviewed and left unchanged), `CHANGELOG.md`, and the
  spec log. Reconciled the **Affected artifacts** list: no `qualitymd` CLI/Go
  change and the `status` lifecycle `readiness` field is unchanged; the listed
  public `use-quality-skill.md` guide was removed independently by concurrent docs
  cleanup, so 0067 lands with no durable-docs delta. Skill `metadata.version`
  bump deferred to the release cut. The In-Review gate was
  collapsed at the user's explicit direction. Removed it from the open-cases
  [index](index.md) and added it to the [archive index](archive/index.md).

- **Implementation + Archival**: Implemented and advanced
  [0066 - Setup feedback log](archive/0066-setup-feedback-log.md) to `Done`,
  moving it (parent and folder) into [`archive/`](archive/). Implemented
  skill-only with **no code change**: added the new durable
  [workflow feedback log](../specs/skills/quality-skill/workflows/setup/feedback-log.md)
  sub-spec under a new `workflows/setup/` folder, made
  [`setup`](../specs/skills/quality-skill/workflows/setup.md) its parent with an
  amended (narrowly widened) mutation surface and close-step authoring, and
  recorded the shared artifact plus its redaction/no-transmission boundary in the
  parent [`/quality` skill](../specs/skills/quality-skill/quality-skill.md) spec.
  Synced the runtime skill
  ([`SKILL.md`](../skills/quality/SKILL.md),
  [`workflows/setup.md`](../skills/quality/workflows/setup.md),
  [`cli-quick-reference.md`](../skills/quality/resources/cli-quick-reference.md)),
  the [use-the-skill guide](../docs/guides/use-quality-skill.md), `CHANGELOG.md`,
  and the spec/doc logs. Skill `metadata.version` bump deferred to the release
  cut. Removed it from the open-cases [index](index.md) and added it to the
  [archive index](archive/index.md).

- **Advance**: Moved
  [0067 - Setup discovery pedagogy](0067-setup-discovery-pedagogy.md) to `Design`
  and authored its [design doc](0067-setup-discovery-pedagogy/design.md): inline
  per-question teaching copy composed with the 0065 presentation tiers, the
  accept-all-and-skip escape removed in favor of a per-question fast confirm,
  the `assumed`→`Low (no signal)` mapping that lets `Low`/`Med`/`High` plus the
  evidence note carry the old vocabulary's meaning, and the final recap inserted
  between discovery and authoring. No code involved.

- **Creation**: Added
  [0067 - Setup discovery pedagogy](0067-setup-discovery-pedagogy.md) at `Draft`
  with its child folder ([index](0067-setup-discovery-pedagogy/index.md),
  [functional spec](0067-setup-discovery-pedagogy/spec.md)). The case
  repositions `/quality setup` discovery as teaching-first: authored
  per-question background and how-to-change-later copy inline in the skill, asks
  every discovery question even at high confidence (removing/revising the
  accept-all-defaults-and-skip escape), relabels the confidence vocabulary to
  `Low`/`Med`/`High` (retaining the evidence note), and adds a final review
  recap of the full question/answer set with a last-chance comment before
  writing `QUALITY.md`. No CLI/Go change; skill + specs + docs only. Revisits
  the confidence vocabulary and accept-all escape that archived
  [0065](archive/0065-setup-discovery-and-close-refinements.md) deliberately
  left alone, and is independent of [0066](0066-setup-feedback-log.md). Spec
  authored; design doc deferred until the spec settles. Listed it under open
  cases in the bundle [index](index.md).

- **Creation**: Added
  [0066 - Setup feedback log](0066-setup-feedback-log.md) at `Design` with its
  child folder ([index](0066-setup-feedback-log/index.md),
  [functional spec](0066-setup-feedback-log/spec.md),
  [design doc](0066-setup-feedback-log/design.md)). The case adds a
  hand-authored, **skill-only** workflow feedback log under `.quality/logs/`
  (`<timestamp>-setup-feedback-log.md`) that records setup-experience friction,
  UX/AX, and efficiency signals for improving the skill — distinct from
  evaluation's per-run `debug-log.md` and from the quality log under
  `.quality/log/`. No CLI/Go change: the directory is created on demand (as
  evaluation already does), the log is recorded locally and never transmitted, so
  no opt-in is needed; secrets and raw prompt-injection text are never written
  and sensitive project context is sanitized. Spec and design authored; no code
  involved. Listed it under open cases in the bundle [index](index.md).

- **Archival**: Advanced
  [0065 - Setup discovery and close refinements](archive/0065-setup-discovery-and-close-refinements.md)
  to `Done` and moved it (parent and folder) into [`archive/`](archive/).
  Implementation landed: made setup discovery agent-agnostic (present all ten
  questions, iterate one at a time without a structured question affordance, page
  through one when available, escapes on request), added the
  read-the-`qualitymd init`-scaffold-before-authoring step, disentangled model
  maturity (`starter`/`immature`/`evaluation-ready`) from the CLI's lifecycle
  `readiness` in the setup close and the Top 10 guide, and renamed the skill
  `modes/` folder to `workflows/` across the runtime skill and `specs/` mirror
  with all live path references updated and append-only logs left frozen. Updated
  durable specs, public docs, and `CHANGELOG.md`; verified `mise run check`. The
  In-Review review gate was collapsed at the user's explicit direction. Updated
  the bundle [index](index.md) and [archive index](archive/index.md).

- **In-Progress**: Advanced
  [0065 - Setup discovery and close refinements](archive/0065-setup-discovery-and-close-refinements.md)
  from `Design` to `In-Progress`; spec and design were settled and implementation
  began across the runtime skill, durable specs, docs, and the folder rename.

- **Archival**: Advanced
  [0062 - Remove wizard mode](archive/0062-remove-wizard-mode.md),
  [0063 - Contextual setup flow](archive/0063-contextual-setup-flow.md), and
  [0064 - Structured setup workflow](archive/0064-structured-setup-workflow.md)
  from `In-Review` to `Done` and moved each parent concept and child folder into
  [`archive/`](archive/index.md). Added their archive [index](archive/index.md)
  entries and removed the open-cases entries from the bundle
  [index](index.md). Repointed the live [0065](0065-setup-discovery-and-close-refinements.md)
  "Relationship to 0064" link into `archive/` and updated its now-stale "0064 is
  In-Review" note. Append-only `log.md` references under
  [`specs/`](../specs/log.md), `specs/skills/quality-skill/`, and
  [`docs/`](../docs/log.md) stay frozen at their original paths as historical
  record.

- **Creation**: Added change
  [0065 - Setup discovery and close refinements](0065-setup-discovery-and-close-refinements.md)
  (`status: Design`) with its
  [functional spec](0065-setup-discovery-and-close-refinements/spec.md),
  [design doc](0065-setup-discovery-and-close-refinements/design.md), and
  [index](0065-setup-discovery-and-close-refinements/index.md). The case captures
  four frictions from a first field run of `/quality setup`: make discovery
  agent-agnostic and present all ten questions (iterating one at a time when no
  structured question affordance exists), read the `qualitymd init` scaffold
  before authoring it, disentangle the skill's model-maturity judgment from the
  CLI's lifecycle `readiness`, and take up the `modes/` → `workflows/` folder
  rename 0064 deferred. Records the affected runtime skill, durable specs, docs,
  and packaging; notes that append-only `log.md` files keep historical `modes/`
  references frozen. Spec and design are settled; no code, runtime, or durable
  spec edits made yet. Updated the open-cases entry in the bundle
  [index](index.md).

- **In-Review**: Completed implementation of
  [0064 - Structured setup workflow](0064-structured-setup-workflow.md) and
  advanced it from `In-Progress` to `In-Review`. Rewrote runtime setup guidance
  as an operator workflow, added a setup brief and ten concrete discovery
  questions, aligned durable setup and parent skill specs, updated
  getting-started guidance, public README copies, specs logs, and changelog, and
  preserved setup's `QUALITY.md`-only mutation boundary. Verified
  `mise run check`.

- **In-Progress**: Advanced
  [0064 - Structured setup workflow](0064-structured-setup-workflow.md) from
  `Design` to `In-Progress`. The functional spec and design doc are settled;
  implementation begins across runtime setup guidance, durable skill specs,
  public docs, logs, and changelog.

- **Design**: Advanced
  [0064 - Structured setup workflow](0064-structured-setup-workflow.md) from
  `Draft` to `Design` and added its
  [design doc](0064-structured-setup-workflow/design.md). The design keeps
  existing dispatch paths stable, rewrites runtime setup guidance as an operator
  playbook, adds a setup brief template, defaults to one compact confirmation
  prompt with all ten discovery questions, and preserves the `QUALITY.md`-only
  mutation boundary. Updated the open-cases entry in the bundle
  [index](index.md). Code not started.

- **Draft**: Created
  [0064 - Structured setup workflow](0064-structured-setup-workflow.md)
  (`Draft`) with its
  [functional spec](0064-structured-setup-workflow/spec.md) and
  [child index](0064-structured-setup-workflow/index.md). The case turns
  `/quality setup` guidance into an explicit setup workflow with a setup brief,
  ten concrete discovery questions, confidence-labeled defaults, prompt framing,
  workflow terminology, and the existing `QUALITY.md`-only mutation boundary.
  Added the case to the open-cases list in the bundle [index](index.md). Design
  and code not started.

- **In-Review**: Completed implementation of
  [0063 - Contextual setup flow](0063-contextual-setup-flow.md) and advanced it
  from `In-Progress` to `In-Review`. Updated runtime setup guidance, durable
  setup and quality-log contracts, getting-started and Top 10 checklist
  guidance, public docs, changelog, and OKF logs so setup analyzes context,
  asks confidence-labeled discovery questions, writes only `QUALITY.md`,
  validates/readiness-checks the model, and offers next-step choices without
  running evaluation, writing the quality log, creating issues, or configuring
  integrations. Verified `mise run check`.

- **In-Progress**: Advanced
  [0063 - Contextual setup flow](0063-contextual-setup-flow.md) from `Design` to
  `In-Progress`. The functional spec and design doc are settled;
  implementation begins across runtime setup guidance, durable skill specs,
  quality-log contracts, public docs, logs, and changelog.

- **Design**: Advanced
  [0063 - Contextual setup flow](0063-contextual-setup-flow.md) from `Draft` to
  `Design` and added its
  [design doc](0063-contextual-setup-flow/design.md). The design keeps setup
  skill-driven, uses a bounded context-analysis fact sheet, asks compact
  discovery questions with confidence-labeled defaults, writes only
  `QUALITY.md`, validates with lint plus Top 10 readiness inspection, and offers
  next-step choices without running evaluation or configuring integrations.
  Updated the open-cases entry in the bundle [index](index.md). Code not
  started.

- **Draft**: Created
  [0063 - Contextual setup flow](0063-contextual-setup-flow.md) (`Draft`) with
  its [functional spec](0063-contextual-setup-flow/spec.md) and
  [child index](0063-contextual-setup-flow/index.md). The case reworks
  `/quality setup` into a context-informed discovery flow that writes only
  `QUALITY.md`, validates/readiness-checks the model, and offers next-step
  choices without running evaluation, writing the quality log, creating issues,
  or configuring recurring review automation. Added the case to the open-cases
  list in the bundle [index](index.md). Design and code not started.

- **In-Review**: Completed implementation of
  [0062 - Remove wizard mode](0062-remove-wizard-mode.md) and advanced it from
  `In-Progress` to `In-Review`. Removed runtime and durable wizard mode files,
  folded bare/ambiguous `/quality` handling into read-only orientation, removed
  wizard from public docs and setup handoffs, updated quality-log/checklist
  wording, reconciled indexes/logs/changelog, and verified `mise run check`.

- **In-Progress**: Advanced
  [0062 - Remove wizard mode](0062-remove-wizard-mode.md) from `Design` to
  `In-Progress`. The functional spec and design doc are settled; implementation
  begins across runtime `/quality` guidance, durable skill specs, public docs,
  indexes, logs, and changelog.

- **Design**: Advanced
  [0062 - Remove wizard mode](0062-remove-wizard-mode.md) from `Draft` to
  `Design` and added its
  [design doc](0062-remove-wizard-mode/design.md). The design treats this as a
  surface reduction rather than a rename, absorbs safe read-only orientation
  into the parent skill routing contract, deletes public wizard mode files, and
  keeps `/quality status` and `/quality next` out of the public contract. Updated
  the open-cases entry in the bundle [index](index.md). Code not started.

- **Draft**: Created
  [0062 - Remove wizard mode](0062-remove-wizard-mode.md) (`Draft`) with its
  [functional spec](0062-remove-wizard-mode/spec.md) and
  [child index](0062-remove-wizard-mode/index.md). The case removes `wizard`
  from the `/quality` public contract without promoting `status` or `next` as
  replacement modes, while preserving read-only orientation for bare or
  ambiguous requests. Added the case to the open-cases list in the bundle
  [index](index.md). Design and code not started.

- **Done**: Landed and archived
  [0061 - Natural scope labels](archive/0061-natural-scope-labels.md) —
  advanced it through implementation to `Done`, moved the parent concept and
  child folder into [`archive/`](archive/index.md), added it to the archive
  [index](archive/index.md), and removed it from the open-cases [index](index.md).
  The implementation updates README examples, runtime `/quality` scope
  resolution guidance, durable `/quality` skill specs, specs logs, and the
  changelog so natural Area and Factor labels are the primary documented scoped
  input while qualified references remain exact-addressing syntax.

- **Design**: Advanced
  [0061 - Natural scope labels](archive/0061-natural-scope-labels.md) from
  `Draft` to `Design` and added its
  [design doc](archive/0061-natural-scope-labels/design.md). The design keeps
  label resolution in the `/quality` skill, treats natural labels as human-edge
  input only, preserves qualified references and stable artifact identifiers,
  and records exact matching plus targeted ambiguity prompts as the
  implementation shape. Updated the open-cases entry in the bundle
  [index](index.md). Code not started.

- **Draft**: Created
  [0061 - Natural scope labels](archive/0061-natural-scope-labels.md) (`Draft`)
  with its [functional spec](archive/0061-natural-scope-labels/spec.md) and
  [child index](archive/0061-natural-scope-labels/index.md). The case makes
  natural Area and Factor labels the primary documented scoped-evaluation input
  for `/quality evaluate`, keeps qualified model references as exact-addressing
  syntax, and preserves stable IDs in durable evaluation artifacts. Added the
  case to the open-cases list in the bundle [index](index.md). Design and code
  not started.

## 2026-06-22

- **Done**: Landed and archived
  [0059 - Unqualified model references](archive/0059-unqualified-model-references.md)
  and [0060 - Friendly path display](archive/0060-friendly-path-display.md) —
  advanced both to `Done`, moved their parent concepts and child folders into
  [`archive/`](archive/index.md), added them to the archive [index](archive/index.md),
  and removed them from the open-cases [index](index.md). Both cases are part of
  the v0.9.0 release state.

- **In-Review**: Completed implementation of
  [0060 - Friendly path display](archive/0060-friendly-path-display.md) and advanced it
  from `In-Progress` to `In-Review`. Separated Area/Factor/Rating display
  helpers from reference helpers; rendered `/` for root Area paths in human
  Markdown report path fields; preserved `area:root`, `root`, `root::factor`,
  and structured `report.json` identifiers; aligned durable specs, runtime
  `/quality` guidance, generated examples, logs, and changelog. Verified
  `go test ./internal/evaluation` and `mise run check`.

- **In-Progress**: Advanced
  [0060 - Friendly path display](archive/0060-friendly-path-display.md) from `Design`
  to `In-Progress`. The functional spec and design doc are settled;
  implementation begins across display/reference helper separation, report path
  rendering, durable specs, runtime skill guidance, generated examples, and
  changelog.

- **Design**: Advanced
  [0060 - Friendly path display](archive/0060-friendly-path-display.md) from `Draft` to
  `Design` and added its
  [design doc](archive/0060-friendly-path-display/design.md). The design separates
  display helpers from qualified and unqualified reference helpers, keeps `/`
  out of reference parsing, and updates report rendering to use display values
  in human path fields. Updated the open-cases entry in the bundle
  [index](index.md). Code not started.

- **Draft**: Created
  [0060 - Friendly path display](archive/0060-friendly-path-display.md) (`Draft`) with
  its [functional spec](archive/0060-friendly-path-display/spec.md) and
  [child index](archive/0060-friendly-path-display/index.md). The case separates
  display values from qualified and unqualified model references, keeps
  reference grammar stable, and proposes rendering the root Area path as `/` in
  human report display contexts. Added the case to the open-cases list in the
  bundle [index](index.md). Design and code not started.

- **In-Review**: Completed implementation of
  [0059 - Unqualified model references](archive/0059-unqualified-model-references.md)
  and advanced it from `In-Progress` to `In-Review`. Added unqualified Area,
  Factor, and Rating reference helpers; type-specific unqualified parsers while
  preserving strict qualified parsing; unqualified Area Breakdown `Path`
  rendering in `report-summary.md` and `report.md`; generated example updates;
  durable spec alignment; runtime `/quality` guidance; and changelog. Verified
  `go test ./...` and `mise run check`.

- **In-Progress**: Advanced
  [0059 - Unqualified model references](archive/0059-unqualified-model-references.md)
  from `Design` to `In-Progress`. The functional spec and design doc are
  settled; implementation begins across unqualified reference helpers,
  type-specific parsing, shared Area Breakdown rendering, durable specs, runtime
  skill guidance, generated examples, and changelog.

- **Done**: Landed and archived
  [0058 - Model reference identifiers](archive/0058-model-reference-identifiers.md)
  — advanced it to `Done` and moved the parent concept and its
  [folder](archive/0058-model-reference-identifiers/index.md) into
  [`archive/`](archive/index.md). The case defined strict Area name, Factor
  name, and Rating Level ID grammar; canonical typed model references;
  edge-only shorthand boundaries; lint diagnostics; JSON Schema patterns;
  revised Area Breakdown columns; durable specs; runtime `/quality` guidance;
  scaffold updates; docs; and changelog. Updated the archive [index](archive/index.md)
  and removed the open-cases entry from the bundle [index](index.md).

- **Refinement**: Updated
  [0059 - Unqualified model references](archive/0059-unqualified-model-references.md)
  to explicitly include the durable `/quality` reporting spec in the affected
  artifacts, so agent-facing report guidance is updated alongside
  `report-summary.md`, `report.md`, `report.json`, generated examples, and the
  shared report renderer.

- **Design**: Created
  [0059 - Unqualified model references](archive/0059-unqualified-model-references.md)
  (`Design`) with its
  [functional spec](archive/0059-unqualified-model-references/spec.md),
  [design doc](archive/0059-unqualified-model-references/design.md), and
  [child index](archive/0059-unqualified-model-references/index.md). The case defines
  unqualified references as a bounded fixed-type form for Areas, Factors, and
  Rating Levels; preserves qualified references for mixed-reference and
  machine-readable surfaces; and plans named helper functions plus Area
  Breakdown rendering updates. Added the case to the open-cases list in the
  bundle [index](index.md). Code not started.

- **In-Review**: Completed implementation of
  [0058 - Model reference identifiers](archive/0058-model-reference-identifiers.md) and
  advanced it from `In-Progress` to `In-Review`. Added strict Area name, Factor
  name, and Rating Level ID validation; generated JSON Schema patterns where
  JSON Schema can express the strict grammar; canonical model-reference
  render/parse helpers; revised Area Breakdown columns; updated generated
  example reports, scaffold placeholders, durable specs, runtime and durable
  `/quality` guidance, authoring guidance, README, and changelog. Verified
  `mise run check` and `go run ./cmd/qualitymd lint --json QUALITY.md`.

- **In-Progress**: Advanced
  [0058 - Model reference identifiers](archive/0058-model-reference-identifiers.md)
  from `Design` to `In-Progress`. The functional spec and design doc are
  settled; implementation begins across strict model-name lint rules, structural
  schema support, canonical model-reference helpers, Area Breakdown rendering,
  durable specs, runtime skill files, docs, and changelog.

- **Design**: Advanced
  [0058 - Model reference identifiers](archive/0058-model-reference-identifiers.md)
  from `Draft` to `Design` and added its
  [design doc](archive/0058-model-reference-identifiers/design.md). The design keeps
  `areaPath` and `factorPath` arrays as durable machine data, adds canonical
  typed reference helpers at human/tool boundaries, enforces strict local names
  through named lint rules with JSON Schema pattern support where structural
  support is possible,
  and updates the shared Area Breakdown renderer to separate Area titles from
  stable Area references. Updated the open-cases entry in the bundle
  [index](index.md). Code not started.

- **Done**: Landed and archived
  [0057 - Quality data directory](archive/0057-quality-data-directory.md) —
  advanced it to `Done` and moved the parent concept and its
  [folder](archive/0057-quality-data-directory/index.md) into
  [`archive/`](archive/index.md). The case defined shared QUALITY.md workspace
  resolution, moved evaluation and quality-log defaults under `.quality/`,
  added the root `config` tooling pointer with lint validation, updated durable
  specs, runtime skill guidance, docs, and release notes, and moved existing
  project quality data into `.quality/`.

- **In-Progress**: Advanced
  [0057 - Quality data directory](archive/0057-quality-data-directory.md) from `Design`
  to `In-Progress`. The functional spec and design doc are settled;
  implementation begins across shared workspace resolution, evaluation/status
  path defaults, lint handling for the root `config` tooling key, durable specs,
  runtime skill guidance, docs, and changelog.

- **Refinement**: Updated
  [0058 - Model reference identifiers](archive/0058-model-reference-identifiers.md) to
  explicitly list [`specs/cli/lint.md`](../specs/cli/lint.md) alongside
  [`specs/cli/lint-rules.md`](../specs/cli/lint-rules.md), so the lint command
  contract and lint rule catalog both account for strict Area name, Factor name,
  and Rating Level ID validation with named rule IDs.

- **Refinement**: Updated
  [0058 - Model reference identifiers](archive/0058-model-reference-identifiers.md) to
  use `rating:` as the canonical model-reference prefix for Rating Level IDs,
  while keeping the formal identifier term "Rating Level ID". This keeps the
  CLI/user-facing reference vocabulary aligned with Area and Factor references
  without renaming the underlying `ratingScale[].level` field.

- **Draft**: Created
  [0058 - Model reference identifiers](archive/0058-model-reference-identifiers.md)
  (`Draft`) with its
  [functional spec](archive/0058-model-reference-identifiers/spec.md) and
  [child index](archive/0058-model-reference-identifiers/index.md). The case defines
  strict Area names, Factor names, and Rating Level IDs; formal Area, Factor,
  and Rating Level IDs; canonical typed model references; edge-only shorthand;
  and clearer report summary Area Breakdown columns that separate Area title,
  stable Area reference, Area-only rating, aggregate rating, and compact Factor
  ratings. Added the case to the open-cases list in the bundle
  [index](index.md).

- **Design**: Advanced
  [0057 - Quality data directory](archive/0057-quality-data-directory.md) from `Draft`
  to `Design` and added its
  [design doc](archive/0057-quality-data-directory/design.md). The design introduces
  `internal/workspace` as the shared resolver for selected model path,
  repository root, config file, quality data directory, evaluation directory,
  and quality log directory; keeps `config` out of the normative Model and JSON
  Schema; and makes lint's unknown-key handling internally rule-configurable
  while defaulting to strict errors except for root `config`. Updated the
  open-cases entry in the bundle [index](index.md). Code not started.

- **Draft**: Created
  [0057 - Quality data directory](archive/0057-quality-data-directory.md) (`Draft`) with
  its [functional spec](archive/0057-quality-data-directory/spec.md) and
  [child index](archive/0057-quality-data-directory/index.md). The case defines the
  QUALITY.md workspace as the resolved operating context for one model file,
  uses `.quality/` as the quality data directory, moves default evaluation runs
  and the quality log under that directory, adds a root `config` pointer for
  tooling config resolution, and keeps lint strict through internally
  configurable unknown-key rule options. Added the case to the open-cases list
  in the bundle [index](index.md).

- **Done**: Landed and archived
  [0056 - Prospective evaluation plan artifacts](archive/0056-prospective-evaluation-plan-artifacts.md) —
  advanced it to `Done` and moved the parent concept and its
  [folder](archive/0056-prospective-evaluation-plan-artifacts/index.md) into
  [`archive/`](archive/index.md). The case made `design.md` and the initial
  `plan.md` prospective `/quality evaluate` artifacts authored before assessment
  begins, with later scope, coverage, rigor, or evidence-strategy changes
  recorded as plan amendments. Updated the archive [index](archive/index.md) and
  removed the entry from the open-cases list in the bundle [index](index.md).

- **In-Review**: Completed implementation of
  [0056 - Prospective evaluation plan artifacts](archive/0056-prospective-evaluation-plan-artifacts.md)
  and advanced it from `In-Progress` to `In-Review`. Added the
  [`design.md`](../specs/evaluation-records/design-md.md) artifact spec,
  clarified [`plan.md`](../specs/evaluation-records/plan-md.md) as prospective
  execution planning with explicit amendments, aligned durable `/quality`
  evaluation workflow, evaluate mode, and reporting specs, updated the runtime
  evaluate prompt and quick reference, and added the unreleased changelog entry.
  Verified `mise run fmt-md-check`, `git diff --check`, and
  `mise run npm-pack-check`.

- **In-Progress**: Advanced
  [0056 - Prospective evaluation plan artifacts](archive/0056-prospective-evaluation-plan-artifacts.md)
  from `Design` to `In-Progress`. The functional spec and design doc are
  settled; implementation begins across the durable evaluation-record specs,
  durable `/quality` evaluation workflow specs, and runtime `/quality evaluate`
  guidance.

- **Done**: Landed and archived
  [0055 - Self-describing evaluation record input](archive/0055-evaluation-input-ergonomics.md) —
  advanced it to `Done` and moved the parent concept and its
  [folder](archive/0055-evaluation-input-ergonomics/index.md) into
  [`archive/`](archive/index.md). The case made evaluation record payloads
  discoverable and validatable through payload-documenting help, no-persist
  dry-runs, aggregated key-named validation for record payloads and `plan.md`
  coverage, synced runtime and durable `/quality` skill surfaces, and added a
  published-skill relative-link package guard. Updated the archive [index](archive/index.md)
  and removed the entry from the open-cases list in the bundle [index](index.md).

- **Design**: Advanced
  [0056 - Prospective evaluation plan artifacts](archive/0056-prospective-evaluation-plan-artifacts.md)
  from `Draft` to `Design` and added its
  [design doc](archive/0056-prospective-evaluation-plan-artifacts/design.md). The design
  lands the change as a contract and prompt repair rather than a CLI behavior
  change: add a planning checkpoint after run creation and before assessment,
  split `design.md`, `plan.md`, `debug-log.md`, and formal records by job, add a
  small `design.md` artifact spec, and keep later scope or coverage changes as
  explicit `plan.md` amendments. Updated the open-cases entry in the bundle
  [index](index.md). Code not started.

- **Draft**: Created
  [0056 - Prospective evaluation plan artifacts](archive/0056-prospective-evaluation-plan-artifacts.md)
  (`Draft`) with its
  [functional spec](archive/0056-prospective-evaluation-plan-artifacts/spec.md)
  and [child index](archive/0056-prospective-evaluation-plan-artifacts/index.md). The case
  tightens `/quality evaluate` so `design.md` and the initial `plan.md` are
  authored immediately after run creation and before assessment begins, separates
  intended evidence planning from actual findings and rating rationale, and
  requires later scope or coverage changes to be explicit amendments. Added the
  case to the open-cases list in the bundle [index](index.md).

- **In-Review**: Completed implementation of
  [0055 - Self-describing evaluation record input](archive/0055-evaluation-input-ergonomics.md)
  and advanced it from `In-Progress` to `In-Review`. Added payload-documenting
  help and canonical examples for evaluation write commands, `-n/--dry-run`
  receipts that do not persist records, aggregated JSON-key validation for record
  payloads and `plan.md` coverage, seeded planned-coverage shape guidance, updated
  the runtime and durable `/quality` skill surfaces, and added a published-skill
  relative-link package guard. Verified with `mise run check`.

- **In-Progress**: Advanced
  [0055 - Self-describing evaluation record input](archive/0055-evaluation-input-ergonomics.md)
  from `Design` to `In-Progress`. The functional spec and design doc are settled;
  implementation begins across the evaluation write commands, validation path,
  status coverage checks, runtime skill surfaces, package guard, and durable
  specs/docs listed in the case.

- **Design**: Advanced
  [0055 - Self-describing evaluation record input](archive/0055-evaluation-input-ergonomics.md)
  from `Draft` to `Design` and added its
  [design doc](archive/0055-evaluation-input-ergonomics/design.md). The design splits each
  evaluation write into decode → validate → plan → commit (the seam that lets
  `-n/--dry-run` validate without persisting and report intended paths), replaces
  first-error validation with a key-named accumulator that also folds the decoder's
  unknown-field and type-mismatch errors into JSON-key vocabulary, and drift-proofs
  the surfaced payloads with one golden-tested canonical example per kind embedded
  in both help and the skill quick reference. Rejected a standalone `schema`
  command, a full example generator, and a descriptor-table validator rewrite.
  Updated the open-cases entry in the bundle [index](index.md). Code not started.

- **Draft**: Created
  [0055 - Self-describing evaluation record input](archive/0055-evaluation-input-ergonomics.md)
  (`Draft`) with its [functional spec](archive/0055-evaluation-input-ergonomics/spec.md).
  Motivated by field run `0001-quality-eval`, the case makes the `qualitymd
evaluation` record-writing surface self-describing — payload-documenting help, a
  no-persist `-n/--dry-run`, aggregated key-named validation — and repairs the
  `/quality` skill surfaces that drifted from the binary (the unshipped
  source-of-truth citation, the stale quick-reference payloads) plus a published-
  bundle link guard. The motivating CLI design-guide additions (a new **Structured
  input** section and Help/Documentation/Errors edits) landed alongside as durable
  docs running ahead of code. Added the case to the open-cases list in the bundle
  [index](index.md).

- **Done**: Landed and archived
  [0054 - Remove improve mode](archive/0054-remove-improve-mode.md) — advanced
  it to `Done` and moved the parent concept and its
  [folder](archive/0054-remove-improve-mode/index.md) into
  [`archive/`](archive/index.md). The case removed `/quality improve` as a public
  mode, added recommendation follow-up with apply-now and issue-tracker handoff
  outcomes, updated runtime skill guidance and durable skill specs, and removed
  the improve mode files. Updated the archive [index](archive/index.md) and
  emptied the open-cases list in the bundle [index](index.md).

- **In-Review**: Completed implementation of
  [0054 - Remove improve mode](archive/0054-remove-improve-mode.md) and advanced
  it from `In-Progress` to `In-Review`. Removed runtime and durable `improve`
  mode files, added recommendation follow-up runtime and durable guidance,
  updated wizard/evaluate/update routing and quality-log ownership, updated user
  docs and examples, and verified `mise run fmt-md-check`, `git diff --check`,
  and targeted stale-reference searches.

- **In-Progress**: Advanced
  [0054 - Remove improve mode](archive/0054-remove-improve-mode.md) from `Design` to
  `In-Progress`. No design doc is required: the implementation is a mechanical
  skill/spec/doc surface change that removes the public `improve` mode and keeps
  recommendation follow-up as a non-mode workflow.

- **Design**: Advanced
  [0054 - Remove improve mode](archive/0054-remove-improve-mode.md) from `Draft` to
  `Design`. The functional spec is settled enough to work through design for
  removing the public `/quality improve` mode while preserving recommendation
  follow-up and issue-tracker handoff.

- **Creation**: Added
  [0054 - Remove improve mode](archive/0054-remove-improve-mode.md) (`status:
Draft`) with its
  [functional spec](archive/0054-remove-improve-mode/spec.md) and
  [child index](archive/0054-remove-improve-mode/index.md). The case simplifies
  the `/quality` skill surface by removing the separate improve mode while
  preserving recommendation follow-up with apply-now and issue-tracker handoff
  outcomes. Updated the bundle [index](index.md).

- **Done**: Landed and archived
  [0053 - Align remaining durable specs](archive/0053-align-remaining-durable-specs.md) —
  advanced it to `Done` and moved the parent concept and its
  [folder](archive/0053-align-remaining-durable-specs/index.md) into
  [`archive/`](archive/index.md). The case split remaining large durable specs
  for evaluation records, lint, and ambient update notices into parent and
  component/artifact contracts. Updated the archive [index](archive/index.md) and
  emptied the open-cases list in the bundle [index](index.md).

- **In-Review**: Completed implementation of
  [0053 - Align remaining durable specs](archive/0053-align-remaining-durable-specs.md)
  and advanced it from `In-Progress` to `In-Review`. Split evaluation-records
  runtime contracts into child specs for the run folder, records, artifacts, and
  report outputs; split the lint command from lint rules and output schema; and
  split ambient update-notice behavior from the explicit update command. Updated
  affected links, indexes, and spec logs. `mise run fmt-md-check` and
  `git diff --check` pass.

- **Creation**: Added
  [0053 - Align remaining durable specs](archive/0053-align-remaining-durable-specs.md)
  (`status: In-Progress`) with its
  [functional spec](archive/0053-align-remaining-durable-specs/spec.md) and
  [child index](archive/0053-align-remaining-durable-specs/index.md). The case applies
  the revised durable-spec granularity guidance to evaluation records, lint, and
  ambient update notice behavior while keeping `SPECIFICATION.md` out of scope as
  the single primary format deliverable. Updated the bundle [index](index.md).

- **Done**: Landed and archived
  [0052 - Durable spec alignment](archive/0052-durable-spec-alignment.md) —
  advanced it to `Done` and moved the parent concept and its
  [folder](archive/0052-durable-spec-alignment/index.md) into [`archive/`](archive/index.md).
  The case aligned durable specs with artifact-spec versus behavioral-component
  guidance, added `/quality` child specs for modes, evaluation workflow,
  reporting, and quality log, narrowed the parent skill spec to shared contracts
  and links, and strengthened the general spec-splitting guidance with a heading
  inventory and fictional examples. Updated the archive [index](archive/index.md)
  and emptied the open-cases list in the bundle [index](index.md).

- **Review correction**: Reopened
  [0052 - Durable spec alignment](archive/0052-durable-spec-alignment.md) from
  `In-Review` to `In-Progress` after review found the parent `/quality` skill
  spec still retained large independently reviewable contracts. Extended the
  functional spec and affected artifacts to split the evaluation workflow,
  reporting contract, and quality log into child component specs before archive.

- **In-Review**: Completed implementation of
  [0052 - Durable spec alignment](archive/0052-durable-spec-alignment.md) and advanced it
  from `In-Progress` to `In-Review`. Added behavioral component specs under
  [`specs/skills/quality-skill/modes/`](../specs/skills/quality-skill/modes/index.md)
  for setup, wizard, evaluate, improve, and update; narrowed the parent
  [`/quality` skill spec](../specs/skills/quality-skill/quality-skill.md) to
  shared contracts plus mode summaries; updated the skill-spec
  [index](../specs/skills/quality-skill/index.md), mode [index](../specs/skills/quality-skill/modes/index.md),
  mode [log](../specs/skills/quality-skill/modes/log.md), and
  [`specs/log.md`](../specs/log.md). Reconciled the affected-artifacts list:
  no code, format spec, runtime skill files, install/scaffold files, or generated
  artifact formats changed. `mise run fmt-md-check` passes.

- **In-Progress**: Advanced
  [0052 - Durable spec alignment](archive/0052-durable-spec-alignment.md) from `Draft` to
  `In-Progress`. Its functional spec is settled and no design doc is required:
  the implementation is a mechanical durable-spec restructuring that adds
  behavioral component specs for the `/quality` modes and narrows the parent
  skill spec to shared contracts plus mode links. Updated the bundle
  [index](index.md).

- **Creation**: Added
  [0052 - Durable spec alignment](archive/0052-durable-spec-alignment.md) (`status: Draft`)
  with its [functional spec](archive/0052-durable-spec-alignment/spec.md) and
  [child index](archive/0052-durable-spec-alignment/index.md). The case aligns durable
  specs with the updated artifact-spec versus behavioral-component guidance,
  starting with child specs for the `/quality` modes while keeping 1:1 artifact
  specs named after their artifacts. Updated the bundle [index](index.md).

- **Done**: Landed and archived
  [0051 - Setup quality-md Area](archive/0051-setup-quality-md-area.md) —
  advanced from `In-Review` to `Done` and moved the parent concept and its
  [folder](archive/0051-setup-quality-md-area/index.md) into [`archive/`](archive/index.md).
  The case added the setup-authored `quality-md` Area pattern, kept
  `qualitymd init` and the CLI scaffold generic, strengthened the authoring guide
  and guide spec around quality-attribute Factor names and one referenced
  assessment across multiple Factors, and synced setup mode plus the durable
  skill spec with the concrete Area shape. Updated the archive [index](archive/index.md)
  and emptied the open-cases list in the bundle [index](index.md).

- **In-Review**: Completed implementation of
  [0051 - Setup quality-md Area](archive/0051-setup-quality-md-area.md) and advanced it
  from `In-Progress` to `In-Review`. Synced the durable skill spec and runtime
  setup mode on the concrete `quality-md` Area shape (`quality-md` key,
  `<Root Title> QUALITY.md` title, Area `description`, path-based `source`, YAML
  comments, and one guide-backed Requirement across Factors). Synced the
  authoring-guide durable spec and runtime guide on quality-attribute Factor names
  and single referenced assessments across multiple Factors. Verified the
  affected-artifacts list: no Go code, CLI scaffold, format spec, or durable docs
  were changed. `mise run fmt-md-check` passes.

- **In-Progress**: Advanced
  [0051 - Setup quality-md Area](archive/0051-setup-quality-md-area.md) from `Design` to
  `In-Progress`. Implementation is limited to the durable skill specs and bundled
  skill prompt/guide files: the CLI scaffold and Go code remain out of scope.

- **Design**: Advanced
  [0051 - Setup quality-md Area](archive/0051-setup-quality-md-area.md) from `Draft` to
  `Design` and added its [design doc](archive/0051-setup-quality-md-area/design.md).
  The design keeps `qualitymd init` generic, puts the `quality-md` Area in skill
  setup's guided population phase, uses normal path-based `source`, adds concise
  YAML comments to distinguish `source` from `assessment`, and records why one
  authoring-guide Requirement can feed multiple Factor roll-ups. Updated the
  [child index](archive/0051-setup-quality-md-area/index.md).

- **Creation**: Added
  [0051 - Setup quality-md Area](archive/0051-setup-quality-md-area.md) (`status: Draft`)
  with its [functional spec](archive/0051-setup-quality-md-area/spec.md) and
  [child index](archive/0051-setup-quality-md-area/index.md). The case proposes a
  setup-authored `quality-md` Area that evaluates the active `QUALITY.md` artifact
  itself against the active authoring guide, keeps `qualitymd init` generic, and
  strengthens the authoring guide around quality-attribute Factor names plus one
  referenced assessment connected to multiple Factors. Functional spec lists
  `specs/skills/quality-skill/quality-skill.md` and
  `specs/skills/quality-skill/guides/authoring-md.md` under **To modify**.
  Updated the bundle [index](index.md).

- **Done**: Landed and archived
  [0050 - Quality log](archive/0050-quality-log.md) — advanced from `In-Progress`
  through `In-Review` to `Done` and moved the parent concept and its
  [folder](archive/0050-quality-log/index.md) into [`archive/`](archive/index.md).
  The case added the convention-first quality log: dated `quality/log/` entries the
  `/quality` skill writes (`setup` seeds an inaugural entry, `improve` appends one
  per confirmed model change), with the format contract in `skills/quality/SKILL.md`,
  the meaningful-change taxonomy in `skills/quality/guides/authoring.md`, the
  inaugural-seed step in `modes/setup.md`, the model-change entry in `modes/improve.md`,
  and the read-only model-history/reconciliation surface in `modes/wizard.md`. Synced
  the durable `/quality` skill spec with a new `## Quality log` section and a
  deferred-CLI bullet, logged it in `specs/log.md`, and added a quality-log mention
  to `docs/guides/use-quality-skill.md`. No Go code: the `qualitymd log` CLI
  command is explicitly deferred. Updated the archive [index](archive/index.md) and
  emptied the open-cases list in the bundle [index](index.md).

- **In-Review**: Completed implementation of
  [0050 - Quality log](archive/0050-quality-log.md) and advanced it from
  `In-Progress` to `In-Review`. Reconciled the Affected artifacts list with reality
  — the only doc beyond the listed durable spec and bundled skill files was
  `docs/guides/use-quality-skill.md`, which already enumerated skill outputs and
  would have gone stale.

- **In-Progress**: Advanced [0050 - Quality log](archive/0050-quality-log.md) from
  `Draft` to `In-Progress`. The functional spec is settled and needs no design
  doc, so implementation of the convention-first quality log begins: the durable
  quality-skill spec subsection plus the bundled skill edits (`SKILL.md`,
  `guides/authoring.md`, `modes/setup.md`, `modes/improve.md`, `modes/wizard.md`).

- **Done**: Landed and archived
  [0049 - Companion JSON Schema](archive/0049-companion-json-schema.md) — advanced
  from `In-Review` to `Done` and moved the parent concept and its
  [folder](archive/0049-companion-json-schema/index.md) into
  [`archive/`](archive/index.md). The case published a structural, non-normative
  JSON Schema for QUALITY.md frontmatter (`quality.schema.json`), generated from
  `internal/schema` by `GenerateJSON()` and guarded against drift by a
  consistency test, embedded via a new root `schema.go`, and emitted by the new
  `qualitymd schema` command; it added the durable `specs/quality-schema-json.md`
  and `specs/cli/schema.md` specs. Updated the archive [index](archive/index.md)
  and removed the entry from the bundle [index](index.md).

- **In-Review**: Completed implementation of
  [0049 - Companion JSON Schema](0049-companion-json-schema.md) and advanced it
  from `In-Progress` to `In-Review`. Landed `internal/schema/jsonschema.go`
  (`GenerateJSON()`), the `internal/schema/gen` `go:generate` entrypoint, the
  generated repo-root `quality.schema.json`, the root `schema.go` embed +
  `Schema()`, the `qualitymd schema` command (verbatim plain output;
  chroma-highlighted + paged on a TTY), command registration, the chroma direct
  dependency, a root no-drift consistency test, and generator unit tests. Synced
  the durable specs/docs: new `specs/quality-schema-json.md` and
  `specs/cli/schema.md`; registered in `specs/index.md`, `specs/cli/index.md`,
  and `specs/log.md`; added the non-normative note in `SPECIFICATION.md`, the
  deferral clarification in `specs/cli/spec.md`, and the `README.md`
  quick-reference row. Build, tests, `go vet`, and `gofmt` clean; `go generate`
  idempotent; redirect round-trip byte-identical. Corrected the functional spec
  and the new artifact-spec to scope the "at least one of" `anyOf` rule to the
  Model only (Area emptiness is the warning-level, semantic `empty-area` check).
  Not committed or archived. Updated the bundle [index](index.md).

- **In-Progress**: Advanced
  [0049 - Companion JSON Schema](0049-companion-json-schema.md) from `Design` to
  `In-Progress` and began implementation. Resolved the one pending external
  input — the schema `$id` domain is `getquality.md`
  (`https://getquality.md/quality.schema.json`), matching the live docs site the
  CLI links to. Work spans `internal/schema` (JSON Schema generation), the
  generated repo-root `quality.schema.json` embedded via a new root `schema.go`,
  the `qualitymd schema` command, and the durable specs/docs being synced in
  parallel. Updated the bundle [index](index.md).

- **Done**: Landed and archived
  [0048 - Area factor report breakdown](archive/0048-area-factor-report-breakdown.md)
  — advanced from `In-Review` to `Done` and moved the parent concept and its
  [folder](archive/0048-area-factor-report-breakdown/index.md) into
  [`archive/`](archive/index.md). The case exposed a compact Area-by-Factor
  breakdown from a first-class report model across `report-summary.md`,
  `report.md`, and `report.json`, renamed the Area rating fields, and landed the
  durable `specs/reports/` artifact specs. Updated the archive
  [index](archive/index.md) and removed the entry from the bundle [index](index.md).

- **Creation**: Added [0050 - Quality log](0050-quality-log.md) (`status: Draft`)
  with its [functional spec](0050-quality-log/spec.md) and
  [child index](0050-quality-log/index.md). The case proposes a curated,
  evidence-linked **quality log** — dated entries under `quality/log/` recording
  meaningful changes to a QUALITY.md model, written by the `/quality` skill
  (`setup` seeds, `improve` appends, `wizard` reconciles drift). Convention-first:
  no `qualitymd log` CLI command or standalone artifact-spec yet. Functional spec
  lists `specs/skills/quality-skill/quality-skill.md` under **To modify**. Updated
  the bundle [index](index.md).

- **Creation**: Added
  [0049 - Companion JSON Schema](0049-companion-json-schema.md) (`status: Draft`)
  with its [functional spec](0049-companion-json-schema/spec.md) and
  [child index](0049-companion-json-schema/index.md). The case proposes a
  structural JSON Schema for QUALITY.md frontmatter — derived from
  [`internal/schema`](../internal/schema/schema.go) so it cannot drift,
  non-normative and subordinate to [`SPECIFICATION.md`](../SPECIFICATION.md) —
  plus a `qualitymd schema` verbatim-artifact command that emits it. Functional
  spec lists new durable specs `specs/quality-schema-json.md` and
  `specs/cli/schema.md`. Updated the bundle [index](index.md).

- **Design**: Advanced
  [0049 - Companion JSON Schema](0049-companion-json-schema.md) from `Draft` to
  `Design` and added its [design doc](0049-companion-json-schema/design.md).
  Decided terminal JSON highlighting uses
  [`chroma`](https://github.com/alecthomas/chroma) directly (promoted from an
  indirect dep, byte-safe on the redirect path), rejecting a glamour code-fence
  (reflows content) and a hand-rolled lipgloss tokenizer (reinvents a lexer).
  Decided the artifact lives at the repo root (`quality.schema.json`, a sibling
  of [`SPECIFICATION.md`](../SPECIFICATION.md)), embedded via a new root
  `schema.go` mirroring [`specification.go`](../specification.go) — over
  co-locating under `internal/schema/` or a dedicated `schema/` dir. Updated the
  parent artifacts (added `go.mod` chroma promotion) and the bundle
  [index](index.md).

- **Design**: Closed the remaining
  [0049 - Companion JSON Schema](0049-companion-json-schema.md) open questions.
  Generation is a `go:generate` tool writing the committed root file (the embed
  _is_ the golden, guarded by a consistency test re-running an exported
  `GenerateJSON()`) over runtime generation — keeping schema changes visible as a
  reviewable diff. The schema declares JSON Schema draft 2020-12 and an
  unversioned `$id` of `https://quality.md/quality.schema.json` (identity, not
  hosting; GitHub raw-root URL as fallback if `quality.md` is not the canonical
  domain). No design questions remain; the case is ready for **In-Progress**.

- **Implementation**: Completed
  [0048 - Area factor report breakdown](0048-area-factor-report-breakdown.md)
  code in `internal/evaluation/` and advanced it to `In-Review`. Renamed the
  report-model Area rating fields to
  `areaRatingState` / `areaRatingResult` / `areaWithDescendantsRatingResult` on
  `areaSummary` and `areas`, dropped the `structural` bool and the
  structural-grouping note, added `factorRatingResults` to the compact
  `areaSummary` layer, and rendered a shared `## Area Breakdown` table (absolute
  Area display paths, path-aware Factor labels, `(area group)`, `not assessed`,
  and empty-Factor states) in `report-summary.md` and `report.md`. Strengthened
  analysis-write validation to reject duplicate and vocabulary-unresolvable
  Factor paths, added regression tests, and regenerated the three
  `0001-quality-eval` golden report fixtures. `go test ./...` green.

- **Status**: Advanced
  [0048 - Area factor report breakdown](0048-area-factor-report-breakdown.md)
  from `Design` to `In-Progress` to begin code implementation in
  `internal/evaluation/`. Durable specs and guide-spec renames already landed.

- **Design refinement**: Sharpened
  [0048 - Area factor report breakdown](0048-area-factor-report-breakdown.md)
  for long-term structure: adopted renaming the opaque Area rating fields to
  `areaRatingResult` / `areaWithDescendantsRatingResult` / `areaRatingState` and
  collapsing the redundant `structural` bool and derived note into the typed
  state (spec and design); asserted that `report.json` element arrays are the
  canonical identifiers while display paths are derived (with a separator-escaping
  non-goal); moved the guide-spec renames into a new `To rename` durable-spec
  subsection; and stated the parent status as two clocks (code vs. durable specs).

- **Durable specs**: Applied the artifact-spec filename convention to existing
  `/quality` runtime guide contracts:
  [`authoring.md`](../specs/skills/quality-skill/guides/authoring-md.md),
  [`getting-started.md`](../specs/skills/quality-skill/guides/getting-started-md.md),
  and
  [`top-10-quality-md-checks.md`](../specs/skills/quality-skill/guides/top-10-quality-md-checks-md.md).
  Updated spec links while leaving the runtime guide artifacts themselves
  unchanged.

- **Durable specs**: Added 1:1 report artifact specs for
  [0048 - Area factor report breakdown](0048-area-factor-report-breakdown.md):
  [`report-summary.md`](../specs/reports/report-summary-md.md),
  [`report.md`](../specs/reports/report-md.md), and
  [`report.json`](../specs/reports/report-json.md). Updated the specs index,
  specs log, evaluation-report command contract, and shared evaluation-records
  report-output contract to point at the new artifact specs.

- **Design refinement**: Updated
  [0048 - Area factor report breakdown](0048-area-factor-report-breakdown.md)
  to use `(area group)` as the human Markdown label for Areas with child Areas
  but no direct requirements, while preserving the typed structural/grouping
  state in machine-readable report data.

- **Design refinement**: Updated
  [0048 - Area factor report breakdown](0048-area-factor-report-breakdown.md)
  to use absolute Area display paths in the example and requirements: the root
  renders as `/ (<root title>)`, descendants render with a leading `/`, and the
  breakdown table's first column is labeled `Path`.

- **Design refinement**: Clarified the rating vocabulary for
  [0048 - Area factor report breakdown](0048-area-factor-report-breakdown.md).
  The functional spec and design now distinguish Area-only ratings from
  Area-with-descendants ratings and recommend concise Markdown labels of `Area`
  and `+ Sub-Areas` when both ratings appear in the breakdown.

- **Design refinement**: Tightened
  [0048 - Area factor report breakdown](0048-area-factor-report-breakdown.md)
  while keeping it in `Design`. The functional spec now explicitly requires
  path-aware Area and Factor labels and keeps detailed rationales in the full
  report; the design narrows cleanup to stale summary-basis helpers and records
  the `areaSummary` naming trade-off.

- **Design**: Advanced
  [0048 - Area factor report breakdown](0048-area-factor-report-breakdown.md)
  from `Draft` to `Design` and added its
  [design doc](0048-area-factor-report-breakdown/design.md). The design keeps
  `areaSummary` as the canonical compact report layer, adds Factor ratings to
  that shape, reuses the same breakdown rendering in `report-summary.md` and
  `report.md`, and strengthens new analysis writes with duplicate and
  model-aware Factor path validation. Updated the [child index](0048-area-factor-report-breakdown/index.md)
  and bundle [index](index.md).

- **Creation**: Opened
  [0048 - Area factor report breakdown](0048-area-factor-report-breakdown.md)
  (`status: Draft`) with its
  [functional spec](0048-area-factor-report-breakdown/spec.md) and
  [child index](0048-area-factor-report-breakdown/index.md). The case strengthens
  generated evaluation reports so `report-summary.md`, `report.md`, and
  `report.json` expose an at-a-glance Area-by-Factor breakdown from the
  assembled report model, with nested Area and Factor paths, structural and
  not-assessed states, and tests for stable machine identifiers. Updated the
  bundle [index](index.md).

## 2026-06-21

- **Done**: Archived
  [0047 - Area terminology changeover](archive/0047-area-terminology.md) after
  implementation and verification. Moved the parent concept and child folder
  into [`archive/`](archive/), set status to `Done`, added it to the
  [archive index](archive/index.md), and removed it from the open
  [changes index](index.md).

- **In-Review**: Completed implementation for
  [0047 - Area terminology changeover](archive/0047-area-terminology.md). The live
  schema, typed model, lint/status surfaces, evaluation records, reports, CLI
  create flag and run naming, durable specs, `/quality` skill guidance,
  scaffold, dogfood model, README/npm README, changelog, and maintained Sparrow
  example bundle now use Area/`areas:`/`areaPath` terminology while preserving
  the default `target` / `Target` rating level. Verified `go test ./...`.

- **In-Progress**: Advanced
  [0047 - Area terminology changeover](archive/0047-area-terminology.md) from `Design`
  to `In-Progress` to implement the full no-compatibility Target to Area
  changeover across the schema, evaluation records, reports, CLI, `/quality`
  skill, scaffold, dogfood model, maintained examples, and docs. Updated the
  bundle [index](index.md).

- **Design refinement**: Updated
  [0047 - Area terminology changeover](archive/0047-area-terminology.md) to keep
  `source` as the Area selector property. The
  [functional spec](archive/0047-area-terminology/spec.md) now explicitly rejects
  renaming `source` and asks prose to distinguish `source` from source code; the
  [design doc](archive/0047-area-terminology/design.md) records the rejected alternatives.

- **Design**: Advanced
  [0047 - Area terminology changeover](archive/0047-area-terminology.md) from `Draft` to
  `Design` and added its
  [design doc](archive/0047-area-terminology/design.md). The design uses a big-bang
  schema/type/record rename from Target to Area, replaces user-facing Subject
  labels with root area or model-file wording, renames evaluation-create
  `--subject` to `--model`, drops the subject altitude from new run folders, and
  guards record decoding so legacy `targetPath` records cannot be mistaken for
  root-area records. Updated the [child index](archive/0047-area-terminology/index.md)
  and bundle [index](index.md).

- **Creation**: Opened
  [0047 - Area terminology changeover](archive/0047-area-terminology.md)
  (`status: Draft`) with its
  [functional spec](archive/0047-area-terminology/spec.md) and
  [child index](archive/0047-area-terminology/index.md). The case replaces the formal
  Target model-node vocabulary with Area, introduces root area as the formal root
  descriptor, rejects legacy `targets:` / `targetPath` compatibility, and scopes
  the change across schema, records, reports, CLI, skill, scaffold, examples,
  and docs. Updated the bundle [index](index.md).

- **Done**: Archived
  [0046 - Evaluation debug log](archive/0046-evaluation-debug-log.md) after
  implementation and verification. Moved the parent concept and child folder
  into [`archive/`](archive/), set status to `Done`, added it to the
  [archive index](archive/index.md), and removed it from the open
  [changes index](index.md).

- **In-Review**: Completed implementation and durable artifact synchronization
  for [0046 - Evaluation debug log](archive/0046-evaluation-debug-log.md). New
  evaluation runs seed a process-only `debug-log.md`; the record specs, CLI
  create contract, `/quality` skill guidance, reference fixture, docs,
  changelog, and skill compatibility metadata now preserve the boundary between
  evaluation-process events and formal subject-quality evidence. Verified
  `go test ./...` and `mise run check`. Updated the bundle [index](index.md).

- **In-Progress**: Advanced
  [0046 - Evaluation debug log](archive/0046-evaluation-debug-log.md) from `Design` to
  `In-Progress`. Implementation will seed `debug-log.md` in evaluation runs,
  update the runtime and CLI contracts, align the `/quality` skill guidance, and
  refresh tests, examples, and release notes. Updated the bundle
  [index](index.md).

- **Design**: Advanced
  [0046 - Evaluation debug log](archive/0046-evaluation-debug-log.md) from `Draft` to
  `Design` and added its
  [design doc](archive/0046-evaluation-debug-log/design.md). The design seeds a small
  run-root `debug-log.md` through `qualitymd evaluation create`, keeps report
  assembly independent of debug prose, and puts the process-only boundary in the
  `/quality` skill guidance. Updated the bundle [index](index.md).

- **Creation**: Opened
  [0046 - Evaluation debug log](archive/0046-evaluation-debug-log.md) (`status: Draft`)
  with its [functional spec](archive/0046-evaluation-debug-log/spec.md) and
  [child index](archive/0046-evaluation-debug-log/index.md). The case adds a
  process-only `debug-log.md` artifact to evaluation runs while keeping
  assessments, analysis, recommendations, and reports authoritative for
  subject-quality judgment. Updated the bundle [index](index.md).

- **Done**: Archived completed change cases
  [0042 - Typed report model](archive/0042-typed-report-model.md),
  [0043 - Evaluation history compatibility](archive/0043-evaluation-history-compatibility.md),
  [0044 - Section unknowns and open questions](archive/0044-section-unknowns-open-questions.md),
  and [0045 - Evaluable body context](archive/0045-evaluable-body-context.md)
  after review. Moved each parent concept and child folder into
  [`archive/`](archive/), set their statuses to `Done`, added them to the
  [archive index](archive/index.md), and removed them from the open
  [changes index](index.md).

- **In-Review**: Completed implementation and durable artifact synchronization
  for [0045 - Evaluable body context](archive/0045-evaluable-body-context.md). The
  authoring guide, format spec, guide contracts, getting-started/top-10/setup
  guidance, README summary, and scaffold now treat the Markdown body as
  evaluable, agent-accessible judgment context; the scaffold test asserts the
  new marker. Reviewed the dogfood `QUALITY.md` and active eval model for
  concrete access-gap fallout. Verified `go test ./...` and `mise run check`.
  Updated the [child index](archive/0045-evaluable-body-context/index.md) and bundle
  [index](index.md).

- **In-Progress**: Advanced
  [0045 - Evaluable body context](archive/0045-evaluable-body-context.md) from `Draft`
  through `Design` (no design doc needed) to `In-Progress`. Implementation will
  update the authoring guide, its durable spec contract, body-context checks,
  and any scaffold or setup guidance needed to treat the Markdown body as
  evaluable, agent-accessible judgment context. Updated the [child index](archive/0045-evaluable-body-context/index.md)
  and bundle [index](index.md).

- **Creation**: Opened
  [0045 - Evaluable body context](archive/0045-evaluable-body-context.md)
  (`status: Draft`) with its
  [functional spec](archive/0045-evaluable-body-context/spec.md) and
  [child index](archive/0045-evaluable-body-context/index.md). The case clarifies that
  the Markdown body is evaluable judgment context for building, justifying, and
  evaluating model quality; that body sections should be concise,
  self-explanatory, and progressively disclosed; and that material support that
  is not agent-accessible is a first-class limitation captured in the relevant
  section's unknowns or open questions. Updated the bundle [index](index.md).

## 2026-06-20

- **In-Review**: Completed
  [0044 - Section unknowns and open questions](archive/0044-section-unknowns-open-questions.md).
  Retired the standalone Known gaps body section in favor of per-section unknowns,
  open questions, and a human/agent review state line across the format spec,
  authoring guide, `init` scaffold (and its test), skill setup/getting-started/
  top-10 checks, the durable specs, the example fixtures, and the dogfood
  `QUALITY.md` and active eval model. Verified `go test ./...` and `mise run check`.
  Updated the [child index](archive/0044-section-unknowns-open-questions/index.md) and
  bundle [index](index.md).

- **In-Progress**: Created and advanced
  [0044 - Section unknowns and open questions](archive/0044-section-unknowns-open-questions.md)
  from `Draft` through `Design` (no design doc needed) to `In-Progress`. The case
  replaces the standalone Known gaps body section with a common per-section shape,
  per-section unknowns and open questions, and a human/agent review state line,
  propagating across the format spec, authoring guide, scaffold, skill checks, and
  dogfood instances. Added the [functional spec](archive/0044-section-unknowns-open-questions/spec.md),
  [child index](archive/0044-section-unknowns-open-questions/index.md), and bundle
  [index](index.md) entry.

- **In-Review**: Completed implementation and durable artifact synchronization
  for [0043 - Evaluation history compatibility](archive/0043-evaluation-history-compatibility.md).
  Evaluation history inspection now surfaces malformed, unsupported, and
  incomplete historical records as typed non-reportable gaps; list/status/latest
  workflows remain usable; report build/gate refuse incompatible selected runs
  with status-oriented diagnostics; and the `/quality` skill guidance treats
  incompatible records as history status rather than subject quality evidence.
  Verified `go test ./...` and `mise run check`. Updated the [child index](archive/0043-evaluation-history-compatibility/index.md)
  and bundle [index](index.md).

- **In-Progress**: Advanced
  [0043 - Evaluation history compatibility](archive/0043-evaluation-history-compatibility.md)
  from `Design` to `In-Progress` to implement tolerant evaluation-history
  inspection, compatibility gaps, and graceful report/list/status behavior.
  Updated the [child index](archive/0043-evaluation-history-compatibility/index.md) and
  bundle [index](index.md).

- **Design**: Advanced
  [0043 - Evaluation history compatibility](archive/0043-evaluation-history-compatibility.md)
  from `Draft` to `Design` and added its
  [design doc](archive/0043-evaluation-history-compatibility/design.md). The design uses
  a tolerant run-inspection layer for status/list/history commands, records
  incompatible files as reportability gaps, keeps record writers strict, and
  gates report build/gate through compatibility status before trusted report
  assembly. Updated the [child index](archive/0043-evaluation-history-compatibility/index.md)
  and bundle [index](index.md).

- **Creation**: Opened
  [0043 - Evaluation history compatibility](archive/0043-evaluation-history-compatibility.md)
  (`status: Draft`) with its
  [functional spec](archive/0043-evaluation-history-compatibility/spec.md) and
  [child index](archive/0043-evaluation-history-compatibility/index.md). The case
  captures the strict-writer / tolerant-reader posture for evaluation history:
  historical or hand-edited runs can become non-reportable compatibility gaps
  without breaking ordinary status, list, latest-run, or fresh-evaluation
  workflows. Updated the bundle [index](index.md).

- **In-Review**: Completed implementation and durable artifact synchronization
  for [0042 - Typed report model](archive/0042-typed-report-model.md). Evaluation reports
  now use typed rating-result, local-rating, next-step, lifecycle,
  missing-metadata, rigor, evaluation-level, path, and gap concepts; report JSON
  exposes explicit state objects; existing invalid rating/severity records become
  non-reportable gaps; and the Sparrow fixture reports were regenerated.
  Verified `mise run check`. Updated the bundle [index](index.md).

- **In-Progress**: Opened
  [0042 - Typed report model](archive/0042-typed-report-model.md) to replace stringly
  typed and implicit evaluation-report states with explicit typed concepts for
  rating results, local target ratings, next steps, lifecycle state, run gaps,
  rigor, evaluation level, missing metadata, and path identities. Added the
  parent case, functional spec, design doc, and child index; updated the bundle
  [index](index.md).

- **Done**: Set status `Done` and archived
  [0041 - Update command and improvements](archive/0041-update-command.md) after
  publishing and verifying `v0.5.0`. Moved the parent concept and child folder
  into [`archive/`](archive/), added the entry to the [archive index](archive/index.md),
  and removed it from the open [changes index](index.md).

- **In-Review**: Completed implementation and durable artifact synchronization
  for [0041 - Update command and improvements](archive/0041-update-command.md). The CLI
  now exposes apply-by-default `qualitymd update` with `--check`, readiness and
  release-notes fields, managed standalone apply, post-apply version
  verification, update-check opt-out, and a cached ambient notice. Renamed the
  `/quality` maintenance mode and durable skill/spec/docs references to
  `update`. Verified `mise run check`, a Windows compile-only check for
  `internal/cli`, and CLI smoke checks for `update` and removed `upgrade`.
  Updated the bundle [index](index.md).

- **In-Progress**: Advanced
  [0041 - Update command and improvements](archive/0041-update-command.md) from `Design`
  to `In-Progress` to begin implementation of the apply-by-default
  `qualitymd update` command, ambient cached update notice, and paired
  `/quality update` skill-mode rename. Updated the bundle [index](index.md).

- **Re-characterization**: Re-characterized
  [0041 - Update command and improvements](archive/0041-update-command.md) as the
  upgrade→update rename plus its improvements, dropping the earlier framing and
  renaming the case from slug `0041-codex-aligned-update` to `0041-update-command`
  (parent, child folder, and the same-day entries below repointed to the new
  path). Expanded scope to rename the paired `/quality upgrade` skill mode to
  `/quality update`: the [functional spec](archive/0041-update-command/spec.md) gains a
  paired skill-mode-rename requirement and a durable-spec change for
  `specs/skills/quality-skill/quality-skill.md`, and the parent's
  **Affected artifacts** now lists the skill spec, the runtime
  `skills/quality/modes/upgrade.md` → `update.md` rename, `SKILL.md` routing,
  `wizard.md`, the CLI quick reference, and the top-10-checks route token. Updated
  the bundle [index](index.md).

- **Redesign**: Reshaped 0041 to an apply-by-default `update` command and renamed
  the case to [0041 - Update command and improvements](archive/0041-update-command.md)
  (slug `0041-upgrade-apply-and-readiness` → `0041-update-command`; earlier
  entries below repointed to the new path). Per the chosen direction, the
  [functional spec](archive/0041-update-command/spec.md) and
  [design doc](archive/0041-update-command/design.md) now rename `upgrade`→`update`
  with apply-by-default and a `--check` advisory (deprecated `upgrade` alias for
  one cycle), and add an ambient cached update notice on ordinary commands. The
  notice deliberately reverses 0032's "ordinary commands MUST NOT check the
  network" rule; it is fenced by strict rails — stderr only, never in
  stdout/`--json`, suppressed off a terminal, in CI, behind
  `QUALITYMD_NO_UPDATE_CHECK`, and for dev builds — served from a cache under
  `$QUALITYMD_HOME` refreshed by a detached, non-blocking subprocess. The managed
  standalone self-apply, readiness gating, and release-notes reference carry
  forward onto the new command shape. Expands the affected-artifact footprint
  (the durable `specs/cli/upgrade.md` is renamed to `specs/cli/update.md`;
  `specs/cli.md`, versioning docs, and the `/quality` skill files all change).
  Updated the bundle [index](index.md).

- **Design**: Advanced
  [0041 - Upgrade self-apply, readiness, and release notes](archive/0041-update-command.md)
  from `Draft` to `Design` and added its
  [design doc](archive/0041-update-command/design.md). The design lands all
  three deltas inside the existing `internal/cli/upgrade.go` seams: widen
  `latestVersionProvider` to return a `{version, ready, releaseNotesURL}` struct
  (so readiness and notes ride the single injectable, offline-testable network
  call); resolve `Ready` from the `assets[]`/`html_url` already in the GitHub
  `releases/latest` response (npm's registry latest is ready by definition); gate
  reported availability and `--apply` on readiness; and add managed standalone to
  the `applySupported`/`upgradeCommand` tables, invoking the existing idempotent
  installer non-interactively via `QUALITYMD_NO_INPUT=1`. Records the Homebrew
  latest-provider quirk and a possible `releaseReady` JSON field as open
  questions. Updated the [child index](archive/0041-update-command/index.md)
  and the bundle [index](index.md).

- **Creation**: Opened
  [0041 - Upgrade self-apply, readiness, and release notes](archive/0041-update-command.md)
  (`status: Draft`) with its [functional spec](archive/0041-update-command/spec.md)
  and [child index](archive/0041-update-command/index.md). The case captures
  three improvements drawn from comparing `qualitymd upgrade` with a conventional
  CLI update flow: extend `--apply` to self-update managed standalone installs (the
  channel the project owns yet 0032 left unable to apply), gate "update
  available" and `--apply` on the target release actually being downloadable, and
  surface a release-notes reference in advisory and `--json` output. Records
  [`internal/cli/upgrade.go`](../internal/cli/upgrade.go),
  [`specs/cli/upgrade.md`](../specs/cli/upgrade.md), the `/quality` upgrade-mode
  skill files, and versioning docs as affected. Added the open-case entry to the
  bundle [index](index.md).

## 2026-06-19

- **Done**: Set status `Done` and archived
  [0040 - Readable report summary](archive/0040-readable-report-summary.md).
  Moved the parent concept and child folder into [`archive/`](archive/), added
  the entry to the [archive index](archive/index.md), and removed it from the
  open [changes index](index.md). Verified `go test ./...`, targeted
  `dprint check`, and `git diff --check`.

- **In-Review**: Completed implementation and durable artifact synchronization
  for [0040 - Readable report summary](archive/0040-readable-report-summary.md). The
  summary renderer now emits the key-details, Summary, Top Issues,
  Recommendations, and Scope & Limitations outline; uses "Full evaluation" and
  "Overall rating" in human-facing summary output; surfaces copyable
  Recommendation IDs; and keeps `report.json` unchanged. Updated durable report
  specs, the `/quality` skill contract/runtime wording, tests, and the worked
  summary example. Verified `go test ./...` and targeted `dprint check`.

- **Implementation**: Advanced
  [0040 - Readable report summary](archive/0040-readable-report-summary.md)
  from `Draft` to `In-Progress` and added its
  [design doc](archive/0040-readable-report-summary/design.md). The design keeps the
  existing `EvaluationReportDocument` and JSON schema, reshaping only the concise
  Markdown renderer into a decision-brief outline with display-time wording for
  "Full evaluation" and "Overall rating".

- **Creation**: Opened
  [0040 - Readable report summary](archive/0040-readable-report-summary.md)
  (`status: Draft`) with its child
  [index](archive/0040-readable-report-summary/index.md) and
  [functional spec](archive/0040-readable-report-summary/spec.md). The spec proposes the
  revised `report-summary.md` outline: key details, Summary, Top Issues,
  Recommendations, and Scope & Limitations; updates human-facing labels to
  "Full evaluation" and "Overall rating"; and makes active Recommendation IDs
  prominent for follow-up prompts.

- **Done**: Set status `Done` and archived
  [0039 — Evaluation command surface redesign](archive/0039-evaluation-command-surface.md).
  Moved the parent concept and child folder into [`archive/`](archive/), fixed
  archive-relative links, added the entry to the [archive index](archive/index.md),
  and removed it from the open [changes index](index.md). Verified `go test
./...`, `go vet ./...`, targeted `dprint check`, and CLI smoke checks for the
  new and removed evaluation command surfaces.

- **Design**: Reconciled the
  [0039 — Evaluation command surface redesign](archive/0039-evaluation-command-surface.md)
  impact list and renamed its section from "Affected specs & docs" to **Affected
  artifacts**. Added an **Affected code** subsection (the `internal/cli` command
  tree, the `internal/evaluation/*` backends incl. `planned_coverage.go`, and
  `internal/status`) and the previously-missing artifacts: the skill spec
  `specs/skills/quality-skill/quality-skill.md`, plus `skills/quality/SKILL.md`,
  `install.md`, `docs/guides/use-quality-skill.md`, `docs/guides/cli-design.md`,
  `docs/guides/write-functional-specs.md`, and `CHANGELOG.md`.

- **Design**: Advanced
  [0039 — Evaluation command surface redesign](archive/0039-evaluation-command-surface.md)
  from `Draft` to `Design` and added its
  [design doc](archive/0039-evaluation-command-surface/design.md): the cobra command
  tree, shared `resolveRun`/payload-batching/output-stream helpers, the
  `plan.md`-folded coverage with read-time validation, the report `build`/`gate`
  split, and the altitude removal — with rejected alternatives and three open
  design calls (flat vs subfolder specs, `list --state` scope, malformed-coverage
  gap name). Updated the [changes index](index.md) and child
  [index](archive/0039-evaluation-command-surface/index.md).

- **Creation**: Opened
  [0039 — Evaluation command surface redesign](archive/0039-evaluation-command-surface.md)
  (`status: Draft`) with its child [index](archive/0039-evaluation-command-surface/index.md)
  and [functional spec](archive/0039-evaluation-command-surface/spec.md). The spec sets a
  single noun/verb rule for the `qualitymd evaluation` surface, renames the
  run-lifecycle verbs, promotes the record kinds and the report to nouns with
  honest verbs, adds run/record `list`, folds planned coverage into `plan.md`
  frontmatter (deleting `set-planned-coverage` and `planned-coverage.json`),
  separates `report gate` from `report build`, and removes the altitude residue.
  Added the case to the open [changes index](index.md).

- **Done**: Set status `Done` and archived
  [0038 — /quality skill interaction UX](archive/0038-quality-skill-interaction-ux.md).
  Moved the parent concept and child folder into [`archive/`](archive/), fixed
  repo-relative links for the deeper path, added the entry to the
  [archive index](archive/index.md), and removed it from the open
  [changes index](index.md).

- **In-Review**: Completed implementation and durable artifact synchronization
  for [0038 — /quality skill interaction UX](archive/0038-quality-skill-interaction-ux.md).
  Added the durable `User interaction contract` section to the `/quality` skill
  spec, added compact shared interaction rules to root `SKILL.md`, and updated
  wizard/evaluate/improve/setup/upgrade mode prompts for run frames, decision
  briefs, stop/reroute behavior, history context, improvement delta reporting,
  and status-first output. Verified targeted Markdown formatting with
  `dprint check`. The case is ready to archive per the requested goal.

- **Implementation**: Advanced change
  [0038 — /quality skill interaction UX](archive/0038-quality-skill-interaction-ux.md)
  from `Design` to `In-Progress` so the settled interaction contract can be
  implemented in the durable `/quality` skill spec and runtime skill files.

- **Design**: Advanced change
  [0038 — /quality skill interaction UX](archive/0038-quality-skill-interaction-ux.md)
  from `Draft` to `Design` and added its
  [design doc](archive/0038-quality-skill-interaction-ux/design.md). The design adds a
  durable `User interaction contract` section to the `/quality` skill spec,
  keeps shared run-frame and decision-brief shapes compact in root `SKILL.md`,
  applies the behavior at mode boundaries, uses existing status/evaluation
  history surfaces rather than new storage, and keeps improvement delta reports
  as human output rather than a new evaluation artifact.

- **Creation**: Added change
  [0038 — /quality skill interaction UX](archive/0038-quality-skill-interaction-ux.md)
  in `Draft` with its
  [functional spec](archive/0038-quality-skill-interaction-ux/spec.md). The change
  proposes a durable interaction contract for the `/quality` skill covering run
  frames, decision briefs, stop/reroute behavior, history-aware operation,
  improvement delta reports, and status-first output, while keeping the existing
  skill/CLI boundary and evaluation artifact format intact. Updated the bundle
  [index](index.md).

- **Done**: Set status `Done` and archived
  [0036 — Harden install scripts and upgrade idiomatics](archive/0036-harden-install-scripts.md)
  and
  [0037 — Render rating-level titles in evaluation reports](archive/0037-report-rating-titles.md).
  Moved both parent concepts and child folders into [`archive/`](archive/),
  fixed repo-relative links for the deeper path, added the entries to the
  [archive index](archive/index.md), and removed them from the open
  [changes index](index.md).

- **In-Review**: Completed implementation and durable artifact synchronization
  for [0037 — Render rating-level titles in evaluation reports](archive/0037-report-rating-titles.md).
  Human Markdown reports now resolve rating labels through the run's rating-scale
  titles with a level-id fallback, while `report.json`, `BuildResult`, and
  `--fail-at-or-below` continue using level ids. Added emoji rating titles to
  `QUALITY.md`, clarified the durable build-report spec, and updated evaluation
  tests for title rendering, JSON id preservation, fallback behavior, and
  non-rating states. Verified `go test ./...`, `go vet ./...`, and
  `dprint check`. The case remains open in [`changes/`](index.md) for review; it
  is not archived until it lands.

- **Implementation**: Advanced change
  [0037 — Render rating-level titles in evaluation reports](archive/0037-report-rating-titles.md)
  from `Draft` to `In-Progress`. No design doc is needed for this localized
  renderer change, so implementation can begin from the settled
  [functional spec](archive/0037-report-rating-titles/spec.md).

- **Creation**: Added change
  [0037 — Render rating-level titles in evaluation reports](archive/0037-report-rating-titles.md)
  in `Draft` with its
  [functional spec](archive/0037-report-rating-titles/spec.md). The change makes the human
  reports (`report.md`, `report-summary.md`) display each rating level's `title`
  instead of its `level` id — bringing the renderer into conformance with the
  existing build-report SHOULD so emoji-bearing titles read in reports — while
  keeping `level` ids in `report.json`, `BuildResult`, and the
  `--fail-at-or-below` gate, and dogfoods emoji titles in `QUALITY.md`. Omits a
  design doc. Updated the bundle [index](index.md).

- **In-Review**: Completed implementation and durable artifact synchronization
  for [0036 — Harden install scripts and upgrade idiomatics](archive/0036-harden-install-scripts.md).
  Added the three-tool checksum fallback with a non-silent skip and the
  print-the-export-line PATH guidance to `install/install.sh`; added the TLS 1.2
  shim, per-user PATH mutation, and `-NonInteractive` gating to
  `install/install.ps1`; made `updateAvailable` SemVer-correct via
  `golang.org/x/mod/semver` in `internal/cli/upgrade.go` with regression
  coverage; commented the intentional Homebrew cask in `.goreleaser.yaml`; and
  synced the durable upgrade spec, install docs, contributor guide, and release
  guide. Verified `go test ./...`, `go vet`, `golangci-lint`, `gofmt`,
  `shellcheck install/install.sh`, `dprint check`, and `goreleaser check`. The
  case remains open in [`changes/`](index.md) for review; it is not archived
  until it lands.

- **Implementation**: Advanced change
  [0036 — Harden install scripts and upgrade idiomatics](archive/0036-harden-install-scripts.md)
  from `Design` to `In-Progress` so the settled installer and upgrade-check fixes
  can be implemented and synced into the durable upgrade spec and install docs.

- **Design**: Advanced change
  [0036 — Harden install scripts and upgrade idiomatics](archive/0036-harden-install-scripts.md)
  from `Draft` to `Design` and added its
  [design doc](archive/0036-harden-install-scripts/design.md). The design settles a
  three-tool checksum fallback with a non-silent skip, a `-bor` TLS 1.2 shim for
  PowerShell 5.1, a deliberately asymmetric PATH model (per-user PATH mutation on
  Windows, print-the-export-line on Unix), `--non-interactive` as a verbosity
  gate rather than a phantom prompt, SemVer-correct update detection via
  `golang.org/x/mod/semver`, and keeping the Homebrew cask with documented
  rationale (rejecting the deprecated formula path).

- **Creation**: Added change
  [0036 — Harden install scripts and upgrade idiomatics](archive/0036-harden-install-scripts.md)
  in `Draft` with its
  [functional spec](archive/0036-harden-install-scripts/spec.md). The change fixes five
  portability/convention gaps surfaced by an install-surface review — dead
  checksum verification on stock Linux, a missing TLS 1.2 pin on Windows
  PowerShell, absent/asymmetric PATH integration, a string-compare standing in
  for a SemVer update check, and a no-op `--non-interactive` flag — and records
  that the Homebrew **cask** (not a formula) is the idiomatic distribution path
  after the "convert to formula" review item was investigated and reversed.
  Updated the bundle [index](index.md).

- **Done**: Set status `Done` and archived
  [0035 — /quality upgrade mode](archive/0035-quality-skill-upgrade-mode.md).
  Moved the parent concept and child folder into [`archive/`](archive/), fixed
  repo-relative links for the deeper path, added the entry to the
  [archive index](archive/index.md), and removed it from the open
  [changes index](index.md).

- **In-Review**: Completed implementation and durable artifact synchronization
  for [0035 — /quality upgrade mode](archive/0035-quality-skill-upgrade-mode.md). Added
  the runtime `skills/quality/modes/upgrade.md` procedure, routed `upgrade` from
  `SKILL.md`, taught wizard to recommend it for stale/incompatible skill/CLI
  state, updated the durable `/quality` skill spec, documented the existing
  install maintenance flow, and verified targeted Markdown formatting. The case
  remains open in [`changes/`](index.md) for review; it is not archived until it
  lands.

- **Implementation**: Advanced change
  [0035 — /quality upgrade mode](archive/0035-quality-skill-upgrade-mode.md) from
  `Design` to `In-Progress` so the settled upgrade-mode spec and design can be
  implemented and synced into the durable skill spec, runtime skill files, and
  install/versioning docs.

- **Design**: Advanced change
  [0035 — /quality upgrade mode](archive/0035-quality-skill-upgrade-mode.md) from
  `Draft` to `Design` and added its
  [design doc](archive/0035-quality-skill-upgrade-mode/design.md). The design adds a
  mode-specific upgrade procedure that snapshots skill and CLI versions, builds a
  plan before mutation, delegates CLI changes to `qualitymd upgrade`, delegates
  skill changes to the Agent Skills installer when available, verifies the
  visible CLI afterward, and warns that skill upgrades may require a restarted
  agent session.

- **Creation**: Added change
  [0035 — /quality upgrade mode](archive/0035-quality-skill-upgrade-mode.md) in `Draft`
  with its [functional spec](archive/0035-quality-skill-upgrade-mode/spec.md). The
  change proposes a skill mode that checks the installed `/quality` skill and
  `qualitymd` CLI pair, diagnoses compatibility and available updates, plans
  skill and CLI upgrade actions, asks before mutation, delegates mechanics to
  the Agent Skills installer and `qualitymd upgrade`, and reports any required
  agent restart or reload. Updated the bundle [index](index.md).

- **Done**: Set status `Done` and archived
  [0031 — Evaluation report summary artifact](archive/0031-report-summary-artifact.md),
  [0032 — CLI managed upgrades](archive/0032-cli-managed-upgrades.md),
  [0033 — Required display titles](archive/0033-required-display-titles.md), and
  [0034 — Skill release metadata](archive/0034-skill-release-metadata.md).
  Moved each parent concept and child folder into [`archive/`](archive/), updated
  the root and archive indexes, and left no open change cases.

- **In-Review**: Completed implementation and durable artifact synchronization
  for [0034 — Skill release metadata](archive/0034-skill-release-metadata.md). The case
  remains open in [`changes/`](index.md) for review; it is not archived until it
  lands.

- **Implementation**: Advanced
  [0034 — Skill release metadata](archive/0034-skill-release-metadata.md) from `Design`
  to `In-Progress` so its settled metadata and release-check design can be
  implemented.

- **In-Review**: Completed implementation and durable artifact synchronization
  for [0031 — Evaluation report summary artifact](archive/0031-report-summary-artifact.md),
  [0032 — CLI managed upgrades](archive/0032-cli-managed-upgrades.md), and
  [0033 — Required display titles](archive/0033-required-display-titles.md). The cases
  remain open in [`changes/`](index.md) for review; they are not archived until
  they land.

- **Design**: Advanced change
  [0034 — Skill release metadata](archive/0034-skill-release-metadata.md) from `Draft`
  to `Design` and added its
  [design doc](archive/0034-skill-release-metadata/design.md). The design uses
  Agent Skills `metadata.version` and `metadata.requires-qualitymd-cli`, mirrors
  the range in `compatibility`, adds release-check validation against the tag and
  changelog, updates runtime/docs wording, and leaves installer enforcement for a
  future package contract.

- **Creation**: Added change
  [0034 — Skill release metadata](archive/0034-skill-release-metadata.md) in `Draft`
  with its [functional spec](archive/0034-skill-release-metadata/spec.md). The change
  proposes project-owned Agent Skills metadata in `skills/quality/SKILL.md` for
  the `/quality` skill SemVer and required `qualitymd` CLI range, mirrored by
  `compatibility` prose and curated release notes, with release-check validation
  and installer enforcement explicitly deferred. Updated the bundle
  [index](index.md).

- **Implementation**: Advanced changes
  [0031 — Evaluation report summary artifact](archive/0031-report-summary-artifact.md),
  [0032 — CLI managed upgrades](archive/0032-cli-managed-upgrades.md), and
  [0033 — Required display titles](archive/0033-required-display-titles.md) to
  `In-Progress` so their settled specs/designs can be implemented and synced
  into durable specs, docs, tests, and examples. Updated the bundle
  [index](index.md).

- **Creation**: Added change
  [0033 — Required display titles](archive/0033-required-display-titles.md) in `Draft`
  with its [functional spec](archive/0033-required-display-titles/spec.md). The change
  proposes required `title` properties on the Model, every Target, every Factor,
  and every Rating Level; adds `Factor.title`; keeps Requirements title-free;
  makes `missing-title` an error across those nodes; and records the affected
  format, lint, init, report, status, skill, README, guide, scaffold, and example
  updates. Updated the bundle [index](index.md).

- **Done**: Set status `Done` and archived
  [0030 — CLI status snapshot command](archive/0030-cli-status-command.md).
  Moved the parent concept and child folder into [`archive/`](archive/), fixed
  repo-relative links for the deeper path, added the entry to the
  [archive index](archive/index.md), and removed it from the open
  [changes index](index.md).

- **Status**: Advanced change
  [0030 — CLI status snapshot command](archive/0030-cli-status-command.md) to
  `In-Review` after adding `qualitymd status [path] [--json]`, the
  `internal/status` snapshot assembler, evaluation helpers for run listing and
  active recommendation counts, CLI tests and status-package tests, durable CLI
  specs, README command docs, and `/quality` skill updates. Verified targeted
  Markdown formatting, `go test ./...`, `mise run vet`, and a smoke run of
  `go run ./cmd/qualitymd status --json`.

- **Implementation**: Advanced change
  [0030 — CLI status snapshot command](archive/0030-cli-status-command.md) from `Design`
  to `In-Progress` so the settled status-command spec and design can be
  implemented and synced into the durable CLI docs and `/quality` skill
  consumers.

- **Design**: Advanced change
  [0030 — CLI status snapshot command](archive/0030-cli-status-command.md) from `Draft`
  to `Design` and added its
  [design doc](archive/0030-cli-status-command/design.md). The design introduces an
  `internal/status` snapshot assembler, keeps CLI rendering thin, reuses lint and
  evaluation mechanics, compares run `model.md` snapshots for staleness, counts
  active recommendations through evaluation-owned helpers, and keeps report-body
  scraping out of the command.

- **Draft**: Replaced the placeholder for
  [0030 — CLI status snapshot command](archive/0030-cli-status-command.md) with a full
  [functional spec](archive/0030-cli-status-command/spec.md). The spec defines the
  read-only `qualitymd status [path] [--json]` invocation, lint validity and
  model-shape snapshot, source coverage, evaluation history and staleness
  signals, active recommendation counts, readiness states, deterministic
  next-action data, and exit behavior. Updated the case and bundle listings.

- **Design**: Advanced change
  [0031 — Evaluation report summary artifact](archive/0031-report-summary-artifact.md)
  from `Draft` to `Design` and added its
  [design doc](archive/0031-report-summary-artifact/design.md). The design reuses the
  existing `ReportJSON` assembly as the single report model, adds a
  `renderReportSummaryMarkdown` projection, extends `BuildResult` with the
  summary path, and keeps `report-summary.md` generated from the same recorded
  run data as `report.md` and `report.json`. Updated the bundle [index](index.md).

- **Design**: Advanced change
  [0032 — CLI managed upgrades](archive/0032-cli-managed-upgrades.md) from `Draft` to
  `Design` and added its
  [design doc](archive/0032-cli-managed-upgrades/design.md). The design stages the work
  through structured version metadata, install-context detection, explicit
  upgrade checks, guarded apply behavior, and GitHub-hosted managed installer
  entrypoints under top-level `install/`.

- **Creation**: Added change
  [0032 — CLI managed upgrades](archive/0032-cli-managed-upgrades.md) in `Draft` with
  its [functional spec](archive/0032-cli-managed-upgrades/spec.md). The change proposes
  structured version metadata, explicit upgrade checks, safe install-method
  detection, advisory output by default, guarded `--apply` behavior, npm launcher
  marking, and a long-term managed standalone installer path. Records affected
  CLI specs, install/versioning docs, release guidance, npm launcher, and
  `/quality` skill consumers. Updated the bundle [index](index.md).

- **Creation**: Added change
  [0031 — Evaluation report summary artifact](archive/0031-report-summary-artifact.md)
  in `Draft` with its
  [functional spec](archive/0031-report-summary-artifact/spec.md). The change proposes
  generating `report-summary.md` beside `report.md` and `report.json` during
  `qualitymd evaluation build-report`, derived from the same recorded run data
  and summary layer, for PR/CI/stakeholder triage without replacing the full
  Evaluation Report. Records affected CLI specs, evaluation record contract,
  skill reporting spec, README, and example bundles. Updated the bundle
  [index](index.md).

- **Done**: Set status `Done` and archived
  [0028 — Require factor references](archive/0028-require-characterized-requirements.md)
  and
  [0029 — Sharpen assessment references and traceability](archive/0029-sharpen-assessment-references.md).
  Moved each parent concept and child folder into [`archive/`](archive/),
  fixed repo-relative links for the deeper path, added both entries to the
  [archive index](archive/index.md), and left
  [0030 — CLI status snapshot command](archive/0030-cli-status-command.md) as the only
  open change case.

- **Creation**: Queued change
  [0030 — CLI status snapshot command](archive/0030-cli-status-command.md) in `Draft`
  with a placeholder [functional spec](archive/0030-cli-status-command/spec.md). The case
  proposes a read-only `qualitymd status [--json]` command that emits a
  deterministic project-state snapshot (model validity and shape, evaluation run
  history, open recommendation counts) so the `/quality` wizard routes from
  structured data instead of hand-parsing `QUALITY.md` and reading report bodies —
  restoring the CLI-owns-mechanical-work split. Records the affected CLI specs,
  README, and `/quality` skill consumers. Spec is a placeholder until the case is
  picked up. Updated the bundle [index](index.md).

- **Status**: Advanced change
  [0029 — Sharpen assessment references and traceability](archive/0029-sharpen-assessment-references.md)
  to `In-Review` after extending `SPECIFICATION.md`'s **Assessment** terminology
  and Requirement section with the inline-or-reference framing and a non-normative
  traceability note, applying the six authoring-guide edits (reserve "source" for
  `Target.source`, traceability-graph job, entity gloss, target/assessment-
  reference duality, split-by-claim job, and the renamed "Reference an external
  assessment" job), nudging the scaffold to "reference" wording, and verifying the
  Go test suite, vet, and Markdown formatting.

- **Implementation**: Advanced change
  [0029 — Sharpen assessment references and traceability](archive/0029-sharpen-assessment-references.md)
  from `Draft` to `In-Progress` (no design doc) so the durable `SPECIFICATION.md`,
  authoring guide, and scaffold edits can be made from the settled spec.

- **Creation**: Added change
  [0029 — Sharpen assessment references and traceability](archive/0029-sharpen-assessment-references.md)
  in `Draft` with its
  [functional spec](archive/0029-sharpen-assessment-references/spec.md). The change frames
  a requirement's `assessment` as either stated inline or a reference to another
  entity, reserves "source" for `Target.source`, extends the "reference"
  terminology 0028 set for factors to the requirement→entity edge, and makes the
  model's traceability graph an authoring concern — across `SPECIFICATION.md`, the
  authoring guide, and the scaffold, with no schema or lint change. Omits a design
  doc. Updated the bundle [index](index.md).

## 2026-06-18

- **Status**: Advanced change
  [0028 — Require factor references](archive/0028-require-characterized-requirements.md)
  to `In-Review` after adding the `missing-factor-reference` lint error,
  updating factor-reference terminology, syncing durable specs/docs/scaffold, and
  verifying the Go test suite.

- **Implementation**: Advanced change
  [0028 — Require factor references](archive/0028-require-characterized-requirements.md)
  from `Design` to `In-Progress` so the settled lint rule, terminology updates,
  durable specs, README, and scaffold guidance can be implemented.

- **Alignment**: Brought change
  [0028 — Require factor references](archive/0028-require-characterized-requirements.md)
  into alignment with the current change-case guides by adding the required
  [Durable spec changes](archive/0028-require-characterized-requirements/spec.md#durable-spec-changes)
  section to its functional spec and moving durable-doc accounting out of its
  design doc.

- **Done**: Set status `Done` and archived
  [0026 — Authoring guide replaces meta-model workflow](archive/0026-authoring-guide-remove-meta-model.md)
  and [0027 — Modularize quality skill modes](archive/0027-modularize-quality-skill.md).
  Moved each parent concept and child folder into [`archive/`](archive/),
  updated repo-relative links for the deeper path, added both entries to the
  [archive index](archive/index.md), and left
  [0028 — Require factor references](archive/0028-require-characterized-requirements.md)
  as the only open change case.

- **Template**: Added a required `## Durable spec changes` section (**To add** /
  **To modify** / **To delete**, each a list or `None`) to the example template
  [spec](archive/0001-example-change/spec.md), so copies account for the durable
  specs a change rewrites. See
  [Writing functional specs](../docs/guides/write-functional-specs.md#durable-spec-changes).

- **Design**: Advanced change
  [0028 — Require factor references](archive/0028-require-characterized-requirements.md)
  from `Draft` to `Design` and added its
  [design doc](archive/0028-require-characterized-requirements/design.md). The design
  keeps `factors` structurally optional, adds a context-sensitive
  `missing-factor-reference` lint error for direct target-level requirements
  without factor references, renames secondary-factor internals to neutral
  factor-reference wording, and records why `missing-factor` and a schema-level
  required `factors` property were rejected.

- **Creation**: Added change
  [0028 — Require factor references](archive/0028-require-characterized-requirements.md)
  in `Draft` with its
  [functional spec](archive/0028-require-characterized-requirements/spec.md). The change
  makes requirements invalid unless they reference at least one factor, keeps
  "lens" available as shorthand, and distinguishes direct target-level
  `factors` references from secondary factors on requirements already nested
  under a factor. Updated the bundle [index](index.md).

- **Schema migration**: Renamed the `changes/` parent concept type from
  `Change` to `Change Case`, updated existing parent concepts and the
  [changes schema](schema.md), renamed the contributor guide to
  [Working with change cases](../docs/guides/work-with-change-cases.md), and
  narrowed `AGENTS.md` so routine prompted edits do not require a Change Case.

- **Status**: Advanced change
  [0027 — Modularize quality skill modes](archive/0027-modularize-quality-skill.md) to
  `In-Review` after keeping `SKILL.md` as the root router/global contract,
  adding setup, wizard, evaluate, and improve mode files under
  `skills/quality/modes/`, renaming supporting skill docs to `resources/`,
  syncing the durable skill spec, and verifying the test suite.

- **Implementation**: Added change
  [0027 — Modularize quality skill modes](archive/0027-modularize-quality-skill.md) in
  `In-Progress` with its
  [functional spec](archive/0027-modularize-quality-skill/spec.md).
  The change keeps `SKILL.md` as the `/quality` router and moves setup, wizard,
  evaluate, and improve procedures into separate files under
  `skills/quality/modes/`, with supporting docs under `skills/quality/resources/`.

- **Status**: Advanced change
  [0026 — Authoring guide replaces meta-model workflow](archive/0026-authoring-guide-remove-meta-model.md)
  to `In-Review` after replacing the skill-facing meta-model reference with
  [authoring.md](../skills/quality/guides/authoring.md),
  removing the bundled `models` CLI/package, making evaluation run creation
  subject-only, syncing durable specs and docs, and verifying the Go test suite.

- **Implementation**: Added change
  [0026 — Authoring guide replaces meta-model workflow](archive/0026-authoring-guide-remove-meta-model.md)
  in `In-Progress` with its
  [functional spec](archive/0026-authoring-guide-remove-meta-model/spec.md) and
  [design doc](archive/0026-authoring-guide-remove-meta-model/design.md). The change
  replaces the skill-facing meta-model reference with an authoring guide
  authoring guide, removes the public `qualitymd models` / model-altitude
  workflow, and syncs durable specs and docs around subject-only evaluation.

- **Done**: Set status `Done` and archived the full in-review set —
  [0012 — Evaluation record format](archive/0012-evaluation-record-format.md),
  [0013 — Evaluation run scaffold](archive/0013-evaluation-run-scaffold.md),
  [0014 — Evaluation record write](archive/0014-evaluation-record-write.md),
  [0015 — Evaluation status and report build](archive/0015-evaluation-report-build.md),
  [0016 — Skill consumes evaluation CLI](archive/0016-skill-consume-eval-cli.md),
  [0017 — Skill rigor and efficiency](archive/0017-skill-rigor-efficiency.md),
  [0018 — Evaluation report UX](archive/0018-evaluation-report-ux.md),
  [0019 — Duplicate assessment status](archive/0019-duplicate-assessment-status.md),
  [0020 — Planned coverage status](archive/0020-planned-coverage-status.md),
  [0021 — Recommendation superseding](archive/0021-recommendation-superseding.md),
  [0022 — Create-run subject validation](archive/0022-create-run-subject-validation.md),
  [0023 — Assessment superseding](archive/0023-assessment-superseding.md),
  [0024 — Report regression coverage](archive/0024-report-regression-coverage.md),
  and [0025 — Durable spec rationale](archive/0025-durable-spec-rationale.md).
  Moved each parent concept and its child folder into
  [`archive/`](archive/), fixed their repo-relative links for the deeper path,
  added them to the [archive index](archive/index.md), and emptied the
  open-changes [index](index.md).

- **Status**: Advanced change
  [0025 — Durable spec rationale](archive/0025-durable-spec-rationale.md) to `In-Review`
  after teaching the three durable contributor guides to keep rationale in the
  spec: a **Background / Motivation** shape entry and per-requirement `Rationale:`
  annotation convention (with litmus and say-it-once rule) in
  [write-functional-specs.md](../docs/guides/write-functional-specs.md), the
  rewritten two-whys convention and refined rationale smells there, the
  absorb-the-why step gated on **Before setting In-Review** in
  [work-with-change-cases.md](../docs/guides/work-with-change-cases.md), and the
  rationale-is-promoted note in
  [write-design-docs.md](../docs/guides/write-design-docs.md). Recorded the guide
  edits in the [docs log](../docs/log.md).

- **Implementation**: Advanced change
  [0025 — Durable spec rationale](archive/0025-durable-spec-rationale.md) from `Design`
  to `In-Progress` so the three durable contributor guides
  ([write-functional-specs.md](../docs/guides/write-functional-specs.md),
  [work-with-change-cases.md](../docs/guides/work-with-change-cases.md),
  [write-design-docs.md](../docs/guides/write-design-docs.md)) can be edited from
  the settled spec and design.

- **Design**: Advanced change
  [0025 — Durable spec rationale](archive/0025-durable-spec-rationale.md) from `Draft`
  to `Design` and added its
  [design doc](archive/0025-durable-spec-rationale/design.md). The design settles a
  two-layer, co-located in-spec rationale convention — a Background/Motivation
  section plus subordinate per-requirement `Rationale:` annotations governed by
  a litmus and a say-it-once rule — over the rejected alternatives (a separate
  Diátaxis explanation doc, design-intent-only depth, and a full ADR embedded in
  the spec), with spec bloat the headline risk mitigated by keeping the
  requirement the lead sentence.

- **Creation**: Added change
  [0025 — Durable spec rationale](archive/0025-durable-spec-rationale.md) in `Draft`
  with its [functional spec](archive/0025-durable-spec-rationale/spec.md). The change
  targets the contributor guides: durable specs inherit a requirement when a
  case archives but lose the case's motivation and the design doc's
  rationale, so editors re-litigate settled lessons and "simplify" rules back
  into the bugs they fixed. The spec states a two-layer in-spec rationale
  convention — a spec-level Background/Motivation section and per-requirement
  `Rationale:` annotations — plus the litmus for when to annotate and an
  absorb-the-why step on landing, and dogfoods the convention itself. Updated
  the bundle [index](index.md).

- **Refinement**: Folded the E49 TypeScript SDK recommendation-quality finding
  into [0017 — Skill rigor and efficiency](archive/0017-skill-rigor-efficiency.md):
  recommendations should name inferable route hints such as affected package,
  path, workflow, maintainer surface, or verification route in existing text
  fields rather than adding a schema field.

- **Status**: Advanced change
  [0024 — Report regression coverage](archive/0024-report-regression-coverage.md) to
  `In-Review` after adding focused temp-run tests for secret-style,
  prompt-injection-style, not-assessed, dotted-path, structural-root, and
  empty-recommendation report behavior.

- **Implementation**: Advanced change
  [0024 — Report regression coverage](archive/0024-report-regression-coverage.md) from
  `Design` to `In-Progress` to turn repeated report-rendering experiment findings
  into focused automated tests without committing benchmark fixture snapshots.

- **Design**: Advanced change
  [0024 — Report regression coverage](archive/0024-report-regression-coverage.md) from
  `Draft` to `Design` and added its
  [design doc](archive/0024-report-regression-coverage/design.md). The design builds
  temporary evaluation runs in tests and asserts high-risk rendered `report.md`
  and `report.json` properties without committing benchmark fixture snapshots.

- **Creation**: Added change
  [0024 — Report regression coverage](archive/0024-report-regression-coverage.md) in
  `Draft` after the experiment program repeatedly found report-rendering
  regressions around seeded safety cases, prompt-injection handling,
  not-assessed propagation, dotted-path limitation extraction, structural roots,
  and empty recommendation arrays.

- **Status**: Advanced change
  [0023 — Assessment superseding](archive/0023-assessment-superseding.md) to
  `In-Review` after implementing assessment `supersedes` metadata,
  superseding status gaps, active/superseded report rendering, durable specs,
  and skill guidance.

- **Implementation**: Advanced change
  [0023 — Assessment superseding](archive/0023-assessment-superseding.md) from `Design`
  to `In-Progress` to close the remaining correction-workflow gap after
  recommendation superseding. The change requires analyses to reference active
  corrected assessments rather than superseded records.

- **Design**: Advanced change
  [0023 — Assessment superseding](archive/0023-assessment-superseding.md) from `Draft`
  to `Design` and added its
  [design doc](archive/0023-assessment-superseding/design.md). The design adds an
  optional `supersedes` list to assessment records, validates superseding and
  stale-analysis references in status, and renders active versus superseded
  assessments while requiring analyses to reference active records.

- **Creation**: Added change
  [0023 — Assessment superseding](archive/0023-assessment-superseding.md) in `Draft`
  after recommendation superseding (E28) left no ergonomic way to correct an
  assessment inside a run while keeping analysis roll-ups bound to active
  judgment.

- **Status**: Advanced change
  [0022 — Create-run subject validation](archive/0022-create-run-subject-validation.md)
  to `In-Review` after validating subject paths before run-folder creation,
  syncing the durable create-run spec, and verifying the failed `--subject .`
  scenario leaves no partial evaluation artifacts.

- **Implementation**: Advanced change
  [0022 — Create-run subject validation](archive/0022-create-run-subject-validation.md)
  from `Design` to `In-Progress` after the E14/E29 CLI UX finding showed that
  invalid `create-run --subject` values can fail after creating an empty run
  skeleton.

- **Design**: Advanced change
  [0022 — Create-run subject validation](archive/0022-create-run-subject-validation.md)
  from `Draft` to `Design` and added its
  [design doc](archive/0022-create-run-subject-validation/design.md). The design
  validates the subject path before creating the evaluation directory or run
  folder so invalid subjects consume no run number and leave no partial
  artifacts.

- **Creation**: Added change
  [0022 — Create-run subject validation](archive/0022-create-run-subject-validation.md)
  in `Draft` after the E14 improve/re-evaluate experiment found that `qualitymd
evaluation create-run --subject .` failed after creating an empty run
  skeleton.

- **Status**: Advanced change
  [0021 — Recommendation superseding](archive/0021-recommendation-superseding.md) to
  `In-Review` after implementing recommendation `supersedes` metadata,
  dangling-reference status gaps, active/superseded report rendering, durable
  specs, and skill guidance.

- **Implementation**: Advanced change
  [0021 — Recommendation superseding](archive/0021-recommendation-superseding.md) from
  `Design` to `In-Progress` so the CLI record schema, status checks, report
  rendering, durable specs, and skill guidance can be updated from the settled
  spec and design.

- **Design**: Advanced change
  [0021 — Recommendation superseding](archive/0021-recommendation-superseding.md) from
  `Draft` to `Design` and added its
  [design doc](archive/0021-recommendation-superseding/design.md). The design uses an
  optional `supersedes` list on recommendation records, validates dangling
  references in status, and keeps superseded recommendations visible while
  choosing Next Action from active recommendations.

- **Creation**: Added change
  [0021 — Recommendation superseding](archive/0021-recommendation-superseding.md) in
  `Draft` after the E15 recommendation-correction trial showed that append-only
  correction records can leave the report's primary Next Action pointing at a
  stale recommendation.

- **Status**: Advanced change
  [0020 — Planned coverage status](archive/0020-planned-coverage-status.md) to
  `In-Review` after implementing `qualitymd evaluation set-planned-coverage`,
  planned-coverage status gaps, durable specs/docs, and skill prompt guidance.

- **Implementation**: Advanced change
  [0020 — Planned coverage status](archive/0020-planned-coverage-status.md) from
  `Design` to `In-Progress` so the CLI writer, status checks, durable specs,
  and skill prompt can be implemented from the settled spec and design.

- **Design**: Advanced change
  [0020 — Planned coverage status](archive/0020-planned-coverage-status.md) from
  `Draft` to `Design` and added its
  [design doc](archive/0020-planned-coverage-status/design.md). The design uses an
  optional CLI-owned `planned-coverage.json` artifact plus status set
  comparisons so interrupted runs can identify missing planned assessments,
  missing planned analyses, and unexpected records without changing current
  behavior for runs that omit the artifact.

- **Creation**: Added change
  [0020 — Planned coverage status](archive/0020-planned-coverage-status.md) in `Draft`
  after the E11 interruption/resume experiment and planned-coverage prototype
  showed that `show-status` can report missing analysis but cannot name missing
  planned assessments without structured planned coverage metadata.

- **Status**: Advanced change
  [0019 — Duplicate assessment status](archive/0019-duplicate-assessment-status.md) to
  `In-Review` after implementing duplicate-assessment renderability checks,
  syncing durable specs, and verifying the command-boundary duplicate trial.

- **Implementation**: Advanced change
  [0019 — Duplicate assessment status](archive/0019-duplicate-assessment-status.md) from
  `Design` to `In-Progress` to implement the `duplicate-assessment` renderability
  gap, sync durable specs, and update skill guidance from the settled design.

- **Design**: Advanced change
  [0019 — Duplicate assessment status](archive/0019-duplicate-assessment-status.md) from
  `Draft` to `Design` and added its
  [design doc](archive/0019-duplicate-assessment-status/design.md). The design detects
  assessment records that share a target path and requirement, reports them as
  `duplicate-assessment` gaps, and fails report rendering through the existing
  renderability gate.

- **Creation**: Added change
  [0019 — Duplicate assessment status](archive/0019-duplicate-assessment-status.md) in
  `Draft` after the experiment program found that re-adding a corrected
  assessment appends a conflicting duplicate record while status still reports
  the run as reportable.

- **Refinement**: Tightened reportability for change
  [0015 — Evaluation status and report build](archive/0015-evaluation-report-build.md).
  `show-status` now reports `missing-root-analysis` unless exactly one analysis
  record has an empty `targetPath`, and `build-report` refuses to render instead
  of silently using a child target as the headline.

- **Refinement**: Extended change
  [0018 — Evaluation report UX](archive/0018-evaluation-report-ux.md) to read report
  context from bounded `design.md` and `plan.md` sections before falling back to
  folder names or rationale text. This makes scope descriptions and planned
  limitations populate `report.md` and `report.json` directly when the skill
  records them.

- **Refinement**: Updated change
  [0018 — Evaluation report UX](archive/0018-evaluation-report-ux.md) after the ESLint
  standard-run experiment to deduplicate equivalent summary limitations and
  accept the skill's `Scope description` / `Narrowing` resolved-parameter
  labels.

- **Refinement**: Updated change
  [0018 — Evaluation report UX](archive/0018-evaluation-report-ux.md) after seeded
  stability repeats to preserve dotted file paths when deriving limitation
  summaries from recorded rationales.

- **Status**: Advanced change
  [0018 — Evaluation report UX](archive/0018-evaluation-report-ux.md) from
  `In-Progress` to `In-Review` after implementing summary-first report
  rendering, syncing durable specs/docs and the skill prompt, and verifying the
  renderer on copied ESLint and DataLoader runs.

- **Implementation**: Advanced change
  [0018 — Evaluation report UX](archive/0018-evaluation-report-ux.md) from `Design` to
  `In-Progress` so the report renderer, durable specs/docs, and skill prompt can
  be updated from the settled functional spec and design.

- **Design**: Created change
  [0018 — Evaluation report UX](archive/0018-evaluation-report-ux.md) in `Design`.
  The change turns the experiment-backed V1 report-shape recommendation into a
  functional spec and design doc for summary-first `report.md`, clearer
  `report.json`, explicit scope/limitations/evidence basis, grouping-target
  rendering, and stable empty recommendation arrays.

- **Status**: Advanced changes
  [0012](archive/0012-evaluation-record-format.md),
  [0013](archive/0013-evaluation-run-scaffold.md),
  [0014](archive/0014-evaluation-record-write.md),
  [0015](archive/0015-evaluation-report-build.md),
  [0016](archive/0016-skill-consume-eval-cli.md), and
  [0017](archive/0017-skill-rigor-efficiency.md) from `In-Progress` to `In-Review`
  after implementing the evaluation CLI surface and syncing the listed durable
  specs/docs.

- **Status**: Advanced changes
  [0012](archive/0012-evaluation-record-format.md),
  [0013](archive/0013-evaluation-run-scaffold.md),
  [0014](archive/0014-evaluation-record-write.md),
  [0015](archive/0015-evaluation-report-build.md),
  [0016](archive/0016-skill-consume-eval-cli.md), and
  [0017](archive/0017-skill-rigor-efficiency.md) from `Design` to `In-Progress` so code
  and durable specs/docs are now phase-authorized.

- **Refinement**: Adopted implementation-readiness review fixes across the open
  changes. Updated lifecycle wording to gate durable spec/doc sync before
  `In-Review`, aligned affected-doc sections with the current `Design` phase,
  renamed the planned durable `show-status` spec path to
  `specs/cli/evaluation-show-status.md`, clarified `show-status` gap semantics,
  kept recommendation Markdown bodies CLI-rendered from structured payload fields,
  expanded durable-doc coverage for the `/quality` skill spec and reference
  examples, and fixed the design-guide link below.

- **Design**: Advanced the evaluation-workflow set and the skill rigor pass from
  `Draft` to `Design`, adding a
  [design doc](../docs/guides/write-design-docs.md) to each (drafted in parallel,
  one per change). The designs settle the _how_ against the settled specs:
  - [0012 — Evaluation record format](archive/0012-evaluation-record-format.md) - the
    contract lives as one enduring bundle-root concept `specs/evaluation-records.md`
    (not under `cli/`), registered the normal OKF way; the skill's switch from
    inlined prose to a reference is deferred to In-Progress.
  - [0013 — Evaluation run scaffold](archive/0013-evaluation-run-scaffold.md) - a new
    `evaluation` Cobra group with `create-run`, a thin CLI over a new
    `internal/evaluation` package, collision-fixing numbering by a single
    directory scan (max+1, create-or-fail), and three-level
    `--evaluation-dir` → `.quality/config.yaml` → `quality/evaluations/`
    resolution.
  - [0014 — Evaluation record write](archive/0014-evaluation-record-write.md) - the
    `internal/evaluation` package owns the record layer; three subcommands share
    one decode→validate→place→write pipeline with strict `DisallowUnknownFields`
    rejection of CLI-owned fields, stateless scan-based numbering with one retry,
    and deterministic recommendation Markdown.
  - [0015 — Evaluation status and report build](archive/0015-evaluation-report-build.md) -
    a typed run-read layer with a shared `Renderable` predicate so green
    `show-status` guarantees `build-report`; one in-memory report renders both
    byte-deterministic files (Glamour kept out of the write path), and
    `--fail-at-or-below` reuses the existing coded-exit mechanism.
  - [0016 — Skill consumes evaluation CLI](archive/0016-skill-consume-eval-cli.md) - the
    skill's evaluation flow maps onto the four commands, judgment JSON piped over
    stdin (CLI stamps/numbers), the inlined Artifact Contract replaced by a
    reference, and the prerequisite probe extended to the evaluation commands.
  - [0017 — Skill rigor and efficiency](archive/0017-skill-rigor-efficiency.md) - rigor
    rules enforced as durable artifact constraints (applied breadth recorded in
    `plan.md`, evidence/locator rigor on existing fields with no schema bump,
    a rating-binding re-check that re-runs the verifying command, compute-then-batch
    writes, and confined `deep` fan-out).
    Updated the bundle [index](index.md). Durable `specs/`/skill/README sync still
    happens per change at In-Progress.

- **Refinement**: Adopted the full explicit verb-object evaluation command
  surface and the CLI-guideline follow-ups. Renamed run creation to
  `qualitymd evaluation create-run`, report rendering to
  `qualitymd evaluation build-report`, added
  `qualitymd evaluation show-status`, renamed the gate flag to
  `--fail-at-or-below`, added `--evaluation-dir` to run creation, added `--file`
  input support to `add-record`, specified atomic numbered writes with one
  recompute retry, and fixed deterministic recommendation rendering order.

- **Refinement**: Renamed the record writer command to
  `qualitymd evaluation add-record assessment|analysis|recommendation`, keeping
  record writes in the evaluation namespace while making the object explicit in
  the command name. Updated dependent wording in
  [0013](archive/0013-evaluation-run-scaffold.md),
  [0014](archive/0014-evaluation-record-write.md),
  [0015](archive/0015-evaluation-report-build.md), and
  [0016](archive/0016-skill-consume-eval-cli.md), plus the bundle index.

- **Refinement**: Kept recommendation artifacts as human-readable Markdown while
  making them first-class CLI-written runtime records. The evaluation-record
  contract now says `recommendations/*.md` carries runtime YAML frontmatter with
  `schemaVersion: 1` and machine-readable metadata, without making the run folder
  an OKF bundle. Change
  [0014 — Evaluation record write](archive/0014-evaluation-record-write.md) now includes
  `qualitymd evaluation add-record recommendation <run>`, and dependent report/skill
  wording reads recommendation records mechanically instead of hand-authoring
  them.

- **Refinement**: Revised change
  [0014 — Evaluation record write](archive/0014-evaluation-record-write.md) and its
  dependents to place record writes under the evaluation namespace:
  `qualitymd evaluation add-record assessment|analysis|recommendation` instead of a separate
  `qualitymd result add` top-level command. Updated dependent wording in
  [0013](archive/0013-evaluation-run-scaffold.md),
  [0015](archive/0015-evaluation-report-build.md), and
  [0016](archive/0016-skill-consume-eval-cli.md), plus the bundle index.

- **Creation**: Added a coordinated set of six changes to sharpen the evaluation
  workflow, drafted in parallel and consolidated here. The seam: the
  deterministic `qualitymd` CLI **writes** every evaluation record (numbering,
  schema stamping, report rendering) while the skill supplies **judgment**.
  - [0012 — Evaluation record format](archive/0012-evaluation-record-format.md)
    (`Draft`) - the keystone: move the artifact contract from the skill prompt
    into an enduring `specs/` spec both the CLI and skill consume.
  - [0013 — Evaluation run scaffold](archive/0013-evaluation-run-scaffold.md) (`Draft`) -
    `qualitymd evaluation create-run`, deterministic shared run numbering (fixes
    a real collision), seeds `model.md`/`design.md`/`plan.md`.
  - [0014 — Evaluation record write](archive/0014-evaluation-record-write.md) (`Draft`) -
    `qualitymd evaluation add-record assessment|analysis|recommendation`,
    JSON judgment payload from `--file` or stdin, inherent validation, atomic
    rejection.
  - [0015 — Evaluation status and report build](archive/0015-evaluation-report-build.md) (`Draft`) -
    `qualitymd evaluation show-status` and `qualitymd evaluation build-report`,
    idempotent rendering of `report.md`/`report.json` from records,
    `--fail-at-or-below` CI gate; renders ratings, never infers them.
  - [0016 — Skill consumes evaluation CLI](archive/0016-skill-consume-eval-cli.md)
    (`Draft`) - the skill drives the CLI instead of hand-authoring run folders,
    records, or reports; replaces its inlined Artifact Contract with a reference.
  - [0017 — Skill rigor and efficiency](archive/0017-skill-rigor-efficiency.md)
    (`Draft`) - operationalized effort levels, verified evidence and pinned
    locators, rating-binding re-check, batched writes, optional deep fan-out;
    independent of the CLI work.
    Updated the bundle [index](index.md). Each spec is the change's _delta_;
    durable `specs/`/skill/README sync happens per change at In-Progress.

- **Done**: Archived change
  [0011 — CLI human output polish](archive/0011-cli-human-output-polish.md)
  after styling `models list`, adding lint next actions to JSON and human
  output, making dev version output include a short VCS revision when available,
  syncing the durable CLI specs, and adding focused output-gate tests.

- **Implementation**: Created and advanced change
  [0011 — CLI human output polish](archive/0011-cli-human-output-polish.md) to
  `In-Progress`. The change covers terminal styling for `models list`, lint
  next actions in JSON and human output, dev version output that includes a VCS
  revision, and broader output-gate tests, with durable CLI specs to be synced
  during implementation.

## 2026-06-17

- **Done**: Archived change
  [0010 — Implement the /quality skill](archive/0010-implement-quality-skill.md)
  after landing the skill artifact, `qualitymd models` CLI surface, durable
  specs/docs, raw runtime example bundle, and verification. Removed it from the
  open-changes index and added it to the archive index.

- **Implementation**: Implemented change
  [0010 — Implement the /quality skill](archive/0010-implement-quality-skill.md):
  added
  `skills/quality/SKILL.md`, implemented `qualitymd models list/view` with Markdown
  and JSON output plus `--source`, moved the bundled quality meta-model under
  `internal/models`, added skill-first install/docs, synced the durable CLI and
  skill specs, re-captured the example as raw runtime artifacts with JSON
  assessment/analysis/report files, ignored default dogfood runs under
  `quality/evaluations/`, and verified the skill/CLI surfaces locally.

- **Implementation**: Advanced change
  [0010 — Implement the /quality skill](archive/0010-implement-quality-skill.md) from
  `Design` to `In-Progress`. The functional spec and design doc are settled, so
  implementation files and durable specs/docs can now be updated: the
  `skills/quality/SKILL.md` artifact, `qualitymd models` CLI surface,
  `.quality/config.yaml` behavior, raw JSON example artifacts, and related durable
  documentation.

- **Refinement**: Added a comprehensive acceptance checklist to change
  [0010 — Implement the /quality skill](archive/0010-implement-quality-skill.md), covering
  skill packaging/install, CLI prerequisite handling including dev builds,
  `qualitymd models` Markdown/JSON behavior, `.quality/config.yaml` validation,
  default dogfood-output ignoring, quick model-altitude dogfooding, JSON artifact
  parsing and shape, minimal `report.json` finding summaries, example re-capture,
  and durable spec/doc sync before **Done**. Optional installer/UI metadata such as
  `agents/openai.yaml` is explicitly deferred until the installer or target agent
  docs require it.

- **Refinement**: Settled the final `SKILL.md` description text for change
  [0010 — Implement the /quality skill](archive/0010-implement-quality-skill.md): "Use when
  a user wants setup, wizard guidance, evaluation, or improvement for quality
  management of a project/entity or one of its components/areas. Trigger for
  requests about quality factors, characteristics, attributes, criteria, areas,
  factors, requirements, improving a quality factor such as
  security/reliability/usability, evaluating a subject against quality criteria, or
  evaluating/improving the QUALITY.md model itself."

- **Refinement**: Added evaluation-directory configuration to change
  [0010 — Implement the /quality skill](archive/0010-implement-quality-skill.md). The skill
  now reads repository-local `.quality/config.yaml` with `evaluationDir` to choose
  the parent directory for numbered evaluation runs, defaulting to
  `quality/evaluations/` when absent; the config is framed as shared qualitymd
  system config that future CLI evaluation commands should also honor. The path must
  be repository-relative, normalized, and unable to escape the repository; unknown
  keys are warned and ignored. Broader configuration ideas (default target file,
  effort, output formats, thresholds, retention, install commands) are deferred
  until they have a concrete need.

- **Refinement**: Added trigger-description requirements for change
  [0010 — Implement the /quality skill](archive/0010-implement-quality-skill.md). The skill
  description must now cover quality management/evaluation/improvement prompts even
  when the user does not mention `QUALITY.md` (for example, improving security
  quality), include mode trigger terms (`setup`, `wizard`, `evaluate`, `improve`)
  and generic quality vocabulary (factors, characteristics, attributes, criteria),
  and stay bounded away from generic copyediting or one-off "higher quality"
  requests. The design records initial `SKILL.md` description text that names
  Targets, Factors, Requirements, subject evaluation, model evaluation/improvement,
  and broad project/entity plus component/target quality, while leaving CLI
  implementation details to the skill body. The design now records the criteria
  used to derive that description and the durable spec sync explicitly carries those
  criteria into `quality-skill.md`'s Frontmatter and metadata section.

- **Refinement**: Added dogfooding guidance to change
  [0010 — Implement the /quality skill](archive/0010-implement-quality-skill.md). The design
  now requires an In-Progress verification pass that installs the skill from the
  working tree, accepts a local development `qualitymd` binary when it exposes the
  required commands, runs a quick model-altitude evaluation against this repo's
  `QUALITY.md`, checks the generated artifact shape, and avoids committing ad hoc
  `quality/evaluations/` output unless deliberately re-captured as a durable
  example.

- **Refinement**: Resolved the remaining open questions for change
  [0010 — Implement the /quality skill](archive/0010-implement-quality-skill.md). Root
  `install.md` now uses a verification-first install flow (`qualitymd --version`,
  documented CLI install/upgrade, verify again, `npx skills add qualitymd/quality.md`,
  `npx skills list`) with the exact package-manager command filled in when the first
  public CLI release channel exists. `report.json` now inlines only minimal generic
  finding summaries by reference for single-file gate/dashboard consumers, while
  full finding detail remains in `assessments/*.json`. Future
  `qualitymd outline QUALITY.md --json`, CLI evaluation record/gate commands,
  and deep-run subagent fanout are explicitly deferred rather than open design
  questions.

- **Refinement**: Generalized the structured finding shape for change
  [0010 — Implement the /quality skill](archive/0010-implement-quality-skill.md). Replaced
  the sample's bespoke top-level `credentialType` with a generic finding object:
  `locator`, `observation`, open `category`, optional `severity`, supporting
  `evidence`, and optional `attributes` for domain-specific metadata. Added the
  broader JSON-shape rule that public top-level fields stay tied to the evaluation
  workflow, while factor- or requirement-specific details live under
  `attributes`.

- **Refinement**: Updated change
  [0010 — Implement the /quality skill](archive/0010-implement-quality-skill.md) so
  `qualitymd models view <name>` supports `--json`. The default output remains
  Markdown with the same terminal-rendered vs plain/verbatim split as
  `qualitymd spec`, preserving byte-for-byte `model.md` snapshots while giving
  humans a readable TTY view. `--json` emits the same source-rewritten bundled model
  as structured data (`model` plus `bodyMarkdown`) for agents and gates that should
  not have to reparse Markdown/YAML. Updated the functional spec, design doc, and
  design log wording.

- **Refinement**: Corrected the onboarding model for change
  [0010 — Implement the /quality skill](archive/0010-implement-quality-skill.md). The skill
  is now the primary entry point, installed from this repo with
  `npx skills add qualitymd/quality.md`; the `qualitymd` CLI is a prerequisite that
  `setup` and `wizard` detect, version-check, and help install or upgrade before
  running CLI-dependent work. Updated the functional spec, design doc, parent
  change, and indexes to remove the plugin-first assumption, added root
  `install.md` to affected docs, and kept Claude plugin/marketplace packaging as a
  possible secondary channel rather than this change's primary distribution path.

- **Design**: Advanced change
  [0010 — Implement the /quality skill](archive/0010-implement-quality-skill.md) from
  `Draft` to `Design` and added its
  [design doc](archive/0010-implement-quality-skill/design.md). Confirmed the three
  **blocking** open items at their recommended resolutions and worked out the _how_:
  the skill ships from `skills/quality/SKILL.md` as an Agent Skills artifact
  installable with `npx skills add qualitymd/quality.md`, while `setup`/`wizard`
  verify the separate `qualitymd` CLI prerequisite before doing CLI-dependent work.
  Specified the `qualitymd models` surface (`list` with `--json`; `view <name>
[--source]` as verbatim Markdown by default and structured JSON with `--json`,
  reusing the `lint --fix` node-rewrite to re-point the meta-model's apex
  `source`), homed in a new `internal/models` package for the bundled-model
  catalog. Settled the raw, non-OKF evaluation
  artifacts — JSON `assessments/`/`analysis/` source-of-record records and a
  `report.json` rendered over them beside the human `report.md`, altitude-first
  folder naming, and `improve`'s new-folder re-evaluation. Recorded the alternatives
  (plugin-marketplace-first/CLI-installer distribution, inline vs referencing
  `report.json`, meta-model embed home) and planned the In-Progress durable sync.
  Updated the change [index](archive/0010-implement-quality-skill/index.md), bundle
  [index](index.md), and the parent concept's status.

- **Refinement**: Tightened change
  [0010 — Implement the /quality skill](archive/0010-implement-quality-skill.md) after
  review. Reconciled the open items' conflicting lifecycle timing — they are now
  **surfaced in Draft**, with the **blocking** ones resolved before **Design** and
  the rest during **In-Progress**, all before **Done** (replacing contradictory
  "settle before Design" / "while Draft" / "before Done" wording across the spec and
  parent). Settled that **evaluation artifacts are raw runtime outputs, not OKF
  concepts**: JSON assessment/analysis records plus a `report.json`, with the worked
  example re-captured raw and its now-unused concept types retired from
  [`specs/schema.md`](../specs/schema.md). Brought the deterministic
  `qualitymd models` command into the change's **Covered** scope, since the model
  altitude depends on it and the skill drives the CLI rather than reimplementing it.
  Renamed the open-items anchor and aligned the parent's `Q1` references to the
  spec's item numbering. Updated the [log](log.md).

- **Creation**: Added change
  [0010 — Implement the /quality skill](archive/0010-implement-quality-skill.md)
  (`status: Draft`) with its
  [functional spec](archive/0010-implement-quality-skill/spec.md) to build the
  specified-but-unimplemented
  [`/quality` skill](../specs/skills/quality-skill/quality-skill.md). The spec
  **defers the behavioral contract** to the durable skill spec and states only the
  delta — package an invocable skill that conforms to it and drives the `qualitymd`
  CLI for every mechanical step — plus the open items and gaps a review of the
  skill spec surfaced for this change to settle: where the skill is packaged; where the
  **model** altitude draws its criteria (the built-in
  [meta-model](../internal/models/quality-meta-model.md) is
  neither referenced nor CLI-exposed); what `setup` does beyond `init`; how the
  default target file resolves (a CLI convention still "to be specified"); and
  whether `improve`'s post-apply re-evaluation writes a new evaluation folder.
  Records [`quality-skill.md`](../specs/skills/quality-skill/quality-skill.md) and
  user docs as affected durable artifacts, to sync once the questions resolve.
  Updated the bundle [index](index.md).

- **Completion**: Implemented and archived
  [0009 — Diagnose rating-scale soundness in the meta-model](archive/0009-rating-scale-diagnostic.md),
  adding the _rating scale and any overrides are well-formed and meaningful_
  requirement to the [meta-model](../internal/models/quality-meta-model.md)'s
  Functionality factor. The meta-model previously assessed the rating scale only
  structurally (lint's "well-shaped" check) and as one clause in a conformance
  list, despite the scale being what turns assessments into verdicts and despite
  per-requirement `ratings` overrides giving authors room to miscalibrate a
  threshold or quietly redefine a level. Synced the Functionality summary and the
  diagnostic coverage checklist, and confirmed `qualitymd lint` still reports the
  model valid. No durable specs/docs were affected — the requirement traces to the
  rating-scale semantics already in [`SPECIFICATION.md`](../SPECIFICATION.md),
  which is unchanged.

- **Creation**: Added change
  [0009 — Diagnose rating-scale soundness in the meta-model](archive/0009-rating-scale-diagnostic.md)
  (`status: Draft`) with its
  [functional spec](archive/0009-rating-scale-diagnostic/spec.md). Prompted by the
  meta-model's coverage asymmetry — six requirements for the Markdown body
  sections, but the rating scale assessed only for structural shape. The change
  adds a single Functionality requirement covering level meaning, band
  separability, floor placement against the model's needs and risks, and sound
  per-requirement criterion overrides, written to pass trivially for a model that
  inherits the suggested scale unchanged. Omits a design doc as a one-requirement
  content change; records no affected durable specs/docs.

- **Completion**: Implemented and archived
  [0008 — Describe targets with title and description](archive/0008-target-display-fields.md),
  adding `title` and `description` to Target, `description` to Model, and the
  matching durable [`SPECIFICATION.md`](../SPECIFICATION.md) prose and
  [`lint`](../specs/cli/lint.md) rule row. The structural schema now accepts
  target display fields, `misplaced-root-key` flags only nested `ratingScale`,
  and focused tests cover accepted nested target `title`/`description` plus the
  still-rejected nested `ratingScale`.

- **Implementation**: Advanced change
  [0008 — Describe targets with title and description](archive/0008-target-display-fields.md)
  from `Design` to `In-Progress` so the schema, linter, and durable
  specification updates can land.

- **Design**: Refined change
  [0008 — Describe targets with title and description](archive/0008-target-display-fields.md):
  made `Model.description` **Optional** (was `Recommended`), matching
  `Target.description`, so `description` reads uniformly across the tree. Updated
  the [functional spec](archive/0008-target-display-fields/spec.md) (Model schema now shows
  `description` as Optional) and the
  [design doc](archive/0008-target-display-fields/design.md): the `OptionalPresence`
  addition, the `# Optional` Model snippet, the composition alternative (now closer
  since `title`/`description` presence matches, still rejected on the `model-content`
  `RequiredAny` group and the mid-list `ratingScale` splice), and the trade-off note.

- **Design**: Advanced change
  [0008 — Describe targets with title and description](archive/0008-target-display-fields.md)
  from `Draft` to `Design` and added its
  [design doc](archive/0008-target-display-fields/design.md). Reading the code showed the
  change is almost entirely schema + prose: three property additions in
  [`internal/schema`](../internal/schema/schema.go) (`Target` gains `title`/
  `description`, `Model` gains `description`) drive everything, because
  `misplaced-root-key` already fires on "a key valid on `Model` but not on this
  target" and so **self-narrows** to `ratingScale` with no rule-logic change, and
  the data-driven spec↔schema consistency test passes once the declarations and
  `SPECIFICATION.md` snippets move in lockstep. The only test edit flips
  `rules_test.go`'s "nested target title" case from a finding to a valid model.
  Recorded the decisions: keep `Model`/`Target` as explicit property lists rather
  than composing `Model` from `Target` (their `title` presence diverges —
  `Recommended` on the root, `Optional` on a target — and the consistency test
  already guards drift); and **not** add a `missing-target-description` warning
  (`RecommendedPresence` is documentary, not auto-enforced), leaving it as a noted
  follow-up. Updated the change [index](archive/0008-target-display-fields/index.md) and
  bundle [index](index.md).

- **Creation**: Added change
  [0008 — Describe targets with title and description](archive/0008-target-display-fields.md)
  (`status: Draft`) with its
  [functional spec](archive/0008-target-display-fields/spec.md). A target's only
  human-facing label today is its map key; the change lets every target carry an
  optional `title` (display name) and a recommended `description` (what the target
  _is_), and adds `description` to the model root. It also reframes the root as a
  **Model** — the model-wide `ratingScale` plus all **Target** properties — so
  the difference from a nested target is a type distinction (`ratingScale` is the
  one Model-only property) rather than the awkward "a non-root target MUST NOT
  declare …" prohibition, matching how `internal/schema` already models the two
  as distinct nodes. Records the
  [`SPECIFICATION.md`](../SPECIFICATION.md) and
  [`lint` sub-spec](../specs/cli/lint.md) deltas plus the `internal/schema` and
  `internal/lint` conformance (the `misplaced-root-key` rule and its
  "nested target title" test case) as affected. Updated the bundle
  [index](index.md). A design doc follows at **Design** for how
  `misplaced-root-key` narrows to `ratingScale` alone.

- **Completion**: Implemented and archived
  [0007 — Delightful human CLI output](archive/0007-delightful-cli-output.md),
  giving the human surface a single brand palette shared with the Fang harness, a
  styled [`lint`](../specs/cli/lint.md) finding list (severity glyphs, color,
  clickable `file:line`, colored summary) and [`init`](../specs/cli/init.md)
  confirmation, runnable `--help` examples on all three commands,
  [`spec`](../specs/cli/spec.md) paging through `$PAGER`/`less`, and an
  informative `--version` recovered from the Go toolchain's embedded build info.
  All of it sits behind the TTY/`NO_COLOR` gate, so the agent-facing plain and
  `--json` paths are byte-for-byte unchanged. Added the **Human output styling**
  and **Binary version** conventions to the [CLI spec](../specs/cli.md), the
  paging clause to the [`spec` sub-spec](../specs/cli/spec.md), the shared
  `internal/cli/style.go`, and focused tests; the styling consolidates onto one
  `colorEnabled` gate, retiring `spec`'s bespoke `shouldRenderSpec`.

- **Completion**: Implemented and archived
  [0006 — Specify and implement the spec command](archive/0006-spec-command.md),
  replacing the placeholder [`spec` sub-spec](../specs/cli/spec.md), adding the
  [design doc](archive/0006-spec-command/design.md), registering
  `qualitymd spec`, embedding [`SPECIFICATION.md`](../SPECIFICATION.md) in the
  binary, rendering Markdown for TTY output with Glamour while preserving
  byte-for-byte Markdown for redirected/plain output, and updating the
  [`README.md`](../README.md) command status.

- **Completion**: Implemented and archived
  [0005 — Single source of truth for the structural schema](archive/0005-schema-source-of-truth.md),
  adding `internal/schema` as the typed structural schema declaration, deriving
  lint's unknown-key, shape, required-property, model-content, and rating-scale
  minimum checks from it, and adding tests that compare the declaration against
  [`SPECIFICATION.md`](../SPECIFICATION.md). Reconciled the public format
  snippet's `title` presence to `Recommended`, updated
  [`lint`](../specs/cli/lint.md) to record schema-derived structural validation,
  and added the [design doc](archive/0005-schema-source-of-truth/design.md).

- **Refinement**: Added a human-readable rendering to change
  [0006 — Specify and implement the spec command](0006-spec-command.md). `spec`
  now **SHOULD** render the specification formatted (via the stack's terminal
  renderer) when stdout is a terminal, while still writing **verbatim Markdown**
  when output must be plain (non-terminal or `NO_COLOR`) — so a redirect still
  reproduces the artifact byte-for-byte and the `--json` carve-out is unaffected.
  The rendered/verbatim split needs no flag: it rides the
  [baseline](../specs/cli.md#baseline)'s existing terminal-detection rule, exactly
  as color does. Updated the [functional spec](0006-spec-command/spec.md) and the
  change's scope.

- **Creation**: Added change
  [0006 — Specify and implement the spec command](0006-spec-command.md)
  (`status: Draft`) with its
  [functional spec](0006-spec-command/spec.md). Picks up the `spec` command that
  [0004](archive/0004-specify-agent-accessibility.md) deferred as "a separate
  change that inherits this baseline." The change settles the still-stub
  [`specs/cli/spec.md`](../specs/cli/spec.md) — whose open questions predate and
  now conflict with the agent-accessibility work (it floats a `--format json`
  form the settled [`--json` convention](../specs/cli.md#conventions) forbids) —
  and lands the command: emit the bundled format specification verbatim as
  Markdown to stdout, no arguments, no `--json` (the verbatim-artifact carve-out),
  full baseline conformance. Records [`specs/cli/spec.md`](../specs/cli/spec.md)
  and [`README.md`](../README.md) as affected; structured forms, sub-views, and
  `spec`-specific flags are deferred. A design doc follows at **Design** for how
  the root-level specification is embedded. Updated the bundle [index](index.md).

- **Refinement**: Recorded the schema-source direction for change
  [0005 — Single source of truth for the structural schema](archive/0005-schema-source-of-truth.md):
  a **typed Go declaration** the linter derives from directly (over an embedded
  data file or a `specs/` concept), with spec/linter consistency enforced by a
  test checking [`SPECIFICATION.md`](../SPECIFICATION.md) against it rather than by
  generating docs — lowest-machinery path that meets the spec's requirements, with
  data-file/generation revisited only if a second consumer appears. Left
  unknown-key typo suggestions as a deferred follow-up, untouched. Detail lands in
  the design doc at **Design**.

- **Creation**: Added change
  [0005 — Single source of truth for the structural schema](archive/0005-schema-source-of-truth.md)
  (`status: Draft`) with its
  [functional spec](archive/0005-schema-source-of-truth/spec.md). Prompted by reviewing
  design.md's linter, which derives its structural rules from one schema artifact:
  our structural schema is encoded twice — implicitly in `internal/lint/rules.go`
  and again in prose in [`SPECIFICATION.md`](../SPECIFICATION.md) and
  [`specs/`](../specs/index.md) — so the two can drift. The change requires a
  single authoritative definition of valid keys, required properties, and the
  rating-scale shape that the linter derives from, as a behavior-preserving
  refactor; records [`SPECIFICATION.md`](../SPECIFICATION.md) and
  [`specs/cli/lint.md`](../specs/cli/lint.md) as affected durable docs; and defers
  doc generation, runtime configuration, and unknown-key typo suggestions. Updated
  the bundle [index](index.md).

- **Completion**: Implemented and archived
  [0004 — Specify and enforce agent accessibility](archive/0004-specify-agent-accessibility.md),
  adding the durable [CLI spec](../specs/cli.md) agent-accessibility contract,
  categorized exit codes (`0`/`1`/`2`/`70`), the broadened `--json` convention,
  and the [`init --json`](../specs/cli/init.md) receipt contract. The
  implementation maps exit categories through Fang, suppresses duplicate stderr
  for already-reported lint findings, adds the neutral `internal/receipt` action
  type, and tests the exit categories plus receipt and overwrite-refusal shapes.
- **Design**: Advanced change
  [0004 — Specify and enforce agent accessibility](archive/0004-specify-agent-accessibility.md)
  from `Draft` to `Design` and added its
  [design doc](archive/0004-specify-agent-accessibility/design.md). Reading
  `fang@v1.0.0` confirmed `fang.Execute` returns the command error and never
  exits, so categorized exit codes ride Fang's intended model: a thin boundary
  mapping in `cli.Execute` (`errors.As` → category, default `ExitInternal`), a
  `CodedError` carrying the category plus a `Silent` flag, a `WithErrorHandler`
  that suppresses the already-reported `lint` found-problems error, and
  Cobra-native `FlagErrorFunc`/`Args` hooks to tag usage errors at their source —
  with only an unknown-subcommand string fallback left as an explicit open
  decision. Picked `0`/`1`/`2`/`70` for success / found-problems / usage /
  internal, broadened _internal error_ to "could not complete the requested
  action" so guarded refusals (e.g. `init` overwrite without `--force`) have a
  home, and ruled `-` plus `--json` a usage error. Settled the one open
  design decision — keep a thin, owned prefix check so an unknown subcommand maps
  to `ExitUsage` (option a), failing safe to `ExitInternal` and pinned by a test.
  Reconciled the
  [functional spec](archive/0004-specify-agent-accessibility/spec.md)'s internal-error
  definition with that broadening, and updated the change
  [index](archive/0004-specify-agent-accessibility/index.md) and bundle
  [index](index.md).

- **Scope**: Broadened the `--json` convention within change
  [0004 — Specify and enforce agent accessibility](archive/0004-specify-agent-accessibility.md)
  after comparing against agent-first CLIs (e.g. Basecamp, where nearly every
  command accepts `--json`). The change now revises the
  [CLI spec](../specs/cli.md)'s `--json` convention from a narrow should-offer
  gate to a SHOULD-by-default: commands SHOULD offer `--json`, human rendering
  stays the default (no auto-JSON), a command MAY omit it only when its output is
  a verbatim artifact that is the payload (e.g. `spec`), and under `--json` a
  side-effecting command emits a **result receipt** rather than human prose. The
  conformance work gains `init --json` (a receipt of the written path / created
  flag / `nextActions`), and [`specs/cli/init.md`](../specs/cli/init.md) joins
  the affected durable docs — its "offers no `--json`" statement is replaced by
  the receipt contract. `spec` stays the deliberate carve-out. Updated the
  functional [spec](archive/0004-specify-agent-accessibility/spec.md), the change
  [index](archive/0004-specify-agent-accessibility/index.md), and the bundle
  [index](index.md).

- **Scope**: Expanded change
  [0004 — Specify and enforce agent accessibility](archive/0004-specify-agent-accessibility.md)
  from a spec-only change to spec **plus** conformance after auditing the shipped
  commands. `internal/cli` exits `1` on every error path, so a `lint` that _found
  problems_ is indistinguishable from a usage or internal error — a baseline
  violation. The change now settles the **exit-code categories** (success,
  ran-but-found-problems, usage error, internal error) concretely, removes both
  _Agent-accessibility and CI requirements_ and _Exit-code semantics_ from the
  [CLI spec](../specs/cli.md)'s "To be specified" list, and brings `init` and
  `lint` into compliance with tests. Threading distinct codes through Fang is a
  real design question, so the change now carries a forthcoming
  [design doc](archive/0004-specify-agent-accessibility/design.md) (added at **Design**);
  the unimplemented [`spec`](../specs/cli/spec.md) command is scoped out as a
  separate change that inherits the baseline. Retitled accordingly and updated
  the change [index](archive/0004-specify-agent-accessibility/index.md) and bundle
  [index](index.md).

- **Creation**: Added change
  [0004 — Specify agent accessibility](archive/0004-specify-agent-accessibility.md)
  (`status: Draft`) with its
  [functional spec](archive/0004-specify-agent-accessibility/spec.md). The change settles
  the _Agent-accessibility and CI requirements_ item on the
  [CLI spec](../specs/cli.md)'s "To be specified" list by adding an **Agent
  accessibility** section framed as two tiers: a baseline every in-scope command
  owes (non-interactivity, stdout-is-payload/stderr-is-everything-else,
  determinism, categorized exit codes, plain non-TTY output) and the opt-in
  capabilities (`--json`, `nextActions`, quiet/verbosity) gated by criteria and
  cross-referenced to the existing conventions. Records
  [`specs/cli.md`](../specs/cli.md) as the only affected durable doc; no command
  behavior changes. Omits a design doc as spec-only work. Updated the bundle
  [index](index.md).

- **Completion**: Implemented and archived
  [0003 — Implement the lint command](archive/0003-implement-lint-command.md),
  adding `qualitymd lint`, the shared lint rule catalog, JSON and human output,
  deterministic finding ordering, in-place `--fix` repair for fixable findings,
  parser/render/write support in `internal/spec`, and focused tests for the rule
  set, output shape, exit behavior, and repair behavior. Updated the README
  status and moved the change into [`archive/`](archive/).

- **Revision**: Hardened [0003's design doc](archive/0003-implement-lint-command/design.md)
  after review. Gave `internal/spec` a one-way dependency — it owns the document
  layer and no longer imports `internal/lint`, which now owns the rule catalog and
  the valid-model convenience (`lint.Load` replacing `spec.Load`) — removing a
  `spec`↔`lint` import cycle and a duplicate validator. Routed misplaced root-only
  keys (`title`/`ratingScale` on a non-root target) to `misplaced-root-key`
  instead of `invalid-frontmatter`; added the original Markdown body to the
  `Document` model so `Render` preserves it byte-for-byte; noted that `Load`'s
  acceptance tightens under the required-frontmatter parser; had `lint` reject a
  bare `-` this phase; clarified that post-repair `summary` counts (including
  `fixable`) reflect the re-lint; and reframed Resolved Questions as Open
  questions with the parent-CLI invocation as the one genuinely-open item.
  Recorded the provisional `lint [path]` shape as deliberately not durably specced
  in the [change](archive/0003-implement-lint-command.md).
- **Revision**: Worked down the open questions and risks in
  [0003's design doc](archive/0003-implement-lint-command/design.md): kept the shared
  document/model code in `internal/spec`, assigned rule-level repair operations
  to `internal/lint` and rendering/atomic replacement to `internal/spec`,
  resolved unknown frontmatter keys as `invalid-frontmatter` in this phase,
  confirmed `lint [path]` as the local invocation shape, and added mitigations
  for YAML round-tripping, deterministic ordering, atomic replacement, and
  symlink paths.
- **Revision**: Scoped `--fix` into change
  [0003 — Implement the lint command](archive/0003-implement-lint-command.md) after
  reviewing fixable-rule behavior. The durable lint spec, implementation spec,
  and design now require deterministic in-place repair of fixable findings,
  transactional per-file writes, post-repair linting, and JSON repair reporting,
  while keeping suppression, rule selection, and patch/full-file repair output
  modes deferred.
- **Design**: Advanced change
  [0003 — Implement the lint command](archive/0003-implement-lint-command.md) from
  `Draft` to `Design` and added its
  [design doc](archive/0003-implement-lint-command/design.md): `lint` parses once into a
  shared document/model graph with stable `modelPath` locations and optional
  source positions, runs narrow rule visitors from `internal/lint`, exposes the
  traversal primitives needed by current rules and future query commands, and
  adds a narrow repair writer for `lint --fix`. The design uses `lint [path]`,
  defaulting to `QUALITY.md`, as the minimum invocation shape while the parent
  CLI spec continues to own the broader file/stdin convention. Updated the
  change [index](archive/0003-implement-lint-command/index.md).

- **Creation**: Added change
  [0003 — Implement the lint command](archive/0003-implement-lint-command.md)
  (`status: Draft`) with a child
  [functional spec](archive/0003-implement-lint-command/spec.md). The change defers
  command-specific behavior to the completed durable
  [`lint` sub-spec](../specs/cli/lint.md), records README status updates as the
  durable docs work before Done, and calls out the remaining cross-cutting CLI
  invocation/file-argument convention as a dependency to settle before Design.
  Updated the bundle [index](index.md).

- **Archival**: Retired the placeholder [0001 — Example change](archive/0001-example-change.md)
  into [`archive/`](archive/) now that the bundle has real changes to follow,
  keeping it as the reference template the
  [propose-a-change guide](../docs/guides/work-with-change-cases.md) points to. Set its
  status to `Done`, fixed the relative links for the deeper path, and updated the
  bundle [index](index.md) and the [archive index](archive/index.md).

- **Completion**: Implemented and archived
  [0002 — Specify the init command](archive/0002-init-command.md), adding
  `qualitymd init`, replacing the durable [`init` sub-spec](../specs/cli/init.md),
  and updating the README status.

- **Refinement**: Tightened change [0002 — Specify the init command](archive/0002-init-command.md)
  after review: framed implementation as the change's own **In-Progress** phase
  rather than deferred work, specified that a successful `init` writes its
  confirmation to standard error (keeping stdout clean for `-` piping), recorded
  the non-atomic `--force` overwrite as a [design](archive/0002-init-command/design.md)
  risk, and trimmed the `--json` note in the
  [functional spec](archive/0002-init-command/spec.md) to a pointer to the
  [CLI spec](../specs/cli.md) convention.

- **Design**: Advanced change [0002 — Specify the init command](archive/0002-init-command.md)
  from `Draft` to `Design` and added its [design doc](archive/0002-init-command/design.md):
  the scaffold ships as a static `//go:embed` asset (comments and body prose can't
  round-trip through YAML struct marshalling), overwrite protection rides on an
  atomic `O_CREATE|O_EXCL` open, and a conformance test runs the embedded skeleton
  through `spec.Load`. Updated the change [index](archive/0002-init-command/index.md).

- **Creation**: Added change [0002 — Specify the init command](archive/0002-init-command.md)
  (`status: Draft`) with its [functional spec](archive/0002-init-command/spec.md), settling
  the "To be specified" list on the [`init` sub-spec](../specs/cli/init.md): the
  scaffold contents (seeded rating scale, a commented target → factor → requirement
  skeleton, recommended body sections as headed stubs), the output target and
  stdout (`-`) piping, and `--force` overwrite protection. Records
  [`specs/cli/init.md`](../specs/cli/init.md) and [`README.md`](../README.md) as
  affected. Updated the bundle [index](index.md).

- **Process**: Defined the relationship between `changes/` and the enduring
  [`specs/`](../specs/index.md) bundle (replacing the "independent for now"
  note) — a Change Case states a _delta_ and is archived, while `specs/` and
  [`SPECIFICATION.md`](../SPECIFICATION.md) hold the _cumulative_ source of
  truth. Added an **Affected specs & docs** section to the
  [Change Case concept](archive/0001-example-change.md) so each change records the durable
  specs and docs it creates or updates, brought into sync before `Done`.

## 2026-06-16

- **Initialization**: Created the `changes/` OKF bundle — a home for incremental
  work, independent of [`specs/`](../specs/index.md) for now. Added the bundle
  [index](index.md), [`schema.md`](schema.md) (`type: Schema`) registering the
  `Change Case`, `Functional Specification`, and `Design Doc` types, and an
  [`archive/`](archive/) folder for completed changes.
- **Creation**: Added a placeholder [Example change](archive/0001-example-change.md)
  (`status: Draft`) with child [spec](archive/0001-example-change/spec.md) and
  [design](archive/0001-example-change/design.md) concepts showing the intended shape.
