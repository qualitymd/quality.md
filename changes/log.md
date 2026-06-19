# Changes Update Log

## 2026-06-19

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
  [0030 — CLI status snapshot command](archive/0030-cli-status-command.md) from
  `Design` to `In-Progress` so the settled status-command spec and design can be
  implemented and synced into the durable CLI docs and `/quality` skill
  consumers.

- **Design**: Advanced change
  [0030 — CLI status snapshot command](archive/0030-cli-status-command.md) from
  `Draft` to `Design` and added its
  [design doc](archive/0030-cli-status-command/design.md). The design introduces
  an `internal/status` snapshot assembler, keeps CLI rendering thin, reuses lint
  and evaluation mechanics, compares run `model.md` snapshots for staleness,
  counts active recommendations through evaluation-owned helpers, and keeps
  report-body scraping out of the command.

- **Draft**: Replaced the placeholder for
  [0030 — CLI status snapshot command](archive/0030-cli-status-command.md) with a
  full [functional spec](archive/0030-cli-status-command/spec.md). The spec
  defines the read-only `qualitymd status [path] [--json]` invocation, lint
  validity and model-shape snapshot, source coverage, evaluation history and
  staleness signals, active recommendation counts, readiness states,
  deterministic next-action data, and exit behavior. Updated the case and bundle
  listings.

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
  [quality-md-guide.md](../skills/quality/resources/quality-md-guide.md),
  removing the bundled `models` CLI/package, making evaluation run creation
  subject-only, syncing durable specs and docs, and verifying the Go test suite.

- **Implementation**: Added change
  [0026 — Authoring guide replaces meta-model workflow](archive/0026-authoring-guide-remove-meta-model.md)
  in `In-Progress` with its
  [functional spec](archive/0026-authoring-guide-remove-meta-model/spec.md) and
  [design doc](archive/0026-authoring-guide-remove-meta-model/design.md). The change
  replaces the skill-facing meta-model reference with a `quality-md-guide.md`
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
  one per change). The designs settle the *how* against the settled specs:
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
    Updated the bundle [index](index.md). Each spec is the change's *delta*;
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
  management of a project/entity or one of its components/targets. Trigger for
  requests about quality factors, characteristics, attributes, criteria, Targets,
  Factors, Requirements, improving a quality factor such as
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
  **blocking** open items at their recommended resolutions and worked out the *how*:
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
  adding the *rating scale and any overrides are well-formed and meaningful*
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
  *is*), and adds `description` to the model root. It also reframes the root as a
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
  internal, broadened *internal error* to "could not complete the requested
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
  commands. `internal/cli` exits `1` on every error path, so a `lint` that *found
  problems* is indistinguishable from a usage or internal error — a baseline
  violation. The change now settles the **exit-code categories** (success,
  ran-but-found-problems, usage error, internal error) concretely, removes both
  *Agent-accessibility and CI requirements* and *Exit-code semantics* from the
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
  the *Agent-accessibility and CI requirements* item on the
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
  note) — a Change Case states a *delta* and is archived, while `specs/` and
  [`SPECIFICATION.md`](../SPECIFICATION.md) hold the *cumulative* source of
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
