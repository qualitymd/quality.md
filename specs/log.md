# Specs Update Log

## 2026-06-27

- **Revision**: Implemented durable spec changes for
  [0147 - Report Descendant Terms](../changes/archive/0147-report-descendant-terms.md).
  Evaluation report specs now require generated reports to use Child Areas for
  immediate Area descendants and Sub-Factors for immediate Factor descendants.

- **Revision**: Implemented durable spec changes for
  [0145 - Scannable Skill Output](../changes/archive/0145-scannable-skill-output.md).
  The `/quality` skill, workflow, reporting, recommendation follow-up, and guide
  specs now require labeled, five-second-scan output templates for multi-fact
  review gates, summaries, closeouts, and next-workflow prompts.

- **Revision**: Implemented durable spec changes for
  [0146 - Changelog Directory](../changes/archive/0146-changelog-directory.md).
  The `/quality` skill specs now name the model-change history as the quality
  changelog under `.quality/changelog/`, require timestamped changelog entry
  filenames, and keep `.quality/logs/` as a flat workflow-log directory.

- **Revision**: Implemented durable spec changes for
  [0144 - Pointed Review Gates](../changes/archive/0144-pointed-review-gates.md).
  The `/quality` skill specs now require direct model-authoring checkpoints to
  state inferred purpose and ask the user to react to the consequential
  assumption instead of defaulting to a generic adjustment prompt.

- **Revision**: Implemented durable spec changes for
  [0143 - Public Review and Improve Workflows](../changes/archive/0143-public-review-improve-workflows.md).
  The `/quality` skill specs now make `review` and `improve` public workflows,
  route them by focus, keep review read-only by default, and require improve to
  confirm focus plus mutation surface before delegating to existing safe routes.

- **Revision**: Implemented durable spec changes for
  [0142 - Requirement Findings Only](../changes/archive/0142-requirement-findings-only.md).
  The Evaluation specs now make Requirement Findings the only finding layer,
  remove Area and Factor Finding report sections, define schema version 3 data
  validation for finding-backed Requirement ratings and non-empty rating drivers,
  and align `/quality evaluate` with driver-only Factor/Area roll-up analysis.

- **Revision**: Implemented durable spec changes for
  [0141 - Area-local Factor References](../changes/archive/0141-area-local-factor-references.md).
  The format and CLI lint specs now make Requirement Factor references
  Area-local, preserve same-named Factors across different Areas, warn when a
  Factor name repeats inside one Area's Factor tree, and remove descendant-Area
  roll-up semantics for ancestor Factors. The `/quality` skill guide specs now
  teach the same Area-local authoring boundary.

- **Revision**: Implemented durable spec changes for
  [0140 - Casual Review Gate Wording](../changes/archive/0140-casual-review-gate-wording.md).
  The `/quality` skill specs now require direct model-authoring checkpoints to
  state the planned change and value prop in simple prose, keep important
  boundaries and quality-log routing visible, and invite welcoming feedback
  before waiting to mutate.

- **Revision**: Implemented durable spec changes for
  [0139 - Real Review Gates](../changes/archive/0139-real-review-gates.md).
  The `/quality` skill specs now require direct model authoring to acknowledge
  long model/guidance reads before they begin, treat intent checkpoints as review
  gates that wait for a response, and keep judgment-shaping `QUALITY.md` model
  edits reviewable even when intent is clear.

- **Revision**: Implemented durable spec changes for
  [0138 - Lightweight Authoring Checkpoint](../changes/archive/0138-lightweight-authoring-checkpoint.md).
  The `/quality` skill specs now define direct model authoring for existing
  `QUALITY.md` edits, require infer-first material follow-up, add the
  conversational intent checkpoint with `looks good` as clear confirmation when
  the mutation is named, and keep meaningful direct model-authoring changes in
  the quality-log contract.

- **Revision**: Implemented durable spec changes for
  [0137 - Run Report Entrypoint](../changes/archive/0137-run-report-entrypoint.md).
  The Evaluation report tree, payload-kind, JSON convention, orchestration,
  protocol, CLI report, and `/quality` skill specs now make `report.md` the
  run-level report, move the root Area detail report to `root-area.md`, and
  index explicit run/headline report refs in `EvaluationOutputResult`.

- **Revision**: Implemented durable spec changes for
  [0134 - Model-relative workspace paths](../changes/archive/0134-model-relative-workspace-paths.md).
  The CLI, status, lint, install-adjacent workflow, and `/quality` skill specs
  now treat workspace config, Evaluation history, quality logs, and workflow
  feedback logs as model-relative to the selected `QUALITY.md`, while keeping
  the Git repository root as the containment boundary.

- **Revision**: Implemented durable spec changes for
  [0136 - Candidate Actions Payload](../changes/archive/0136-candidate-actions-payload.md).
  The Evaluation payload, JSON convention, routine, report-tree, and `/quality`
  skill specs now use `candidateActions` for finding-local remediation leads,
  require candidate action IDs to be local to the containing Finding, reject the
  legacy `actions` field, and keep candidate actions out of v0 reports and
  closeouts.

- **Revision**: Implemented durable spec changes for
  [0135 - Structured Finding Core](../changes/archive/0135-structured-finding-core.md).
  The Evaluation payload, routine, JSON convention, report-tree, and `/quality`
  skill specs now require a shared Finding Core with statement, condition,
  criteria, cause, effect, and evidence; finding IDs are payload-local; and
  reports use one findings table/detail structure across Requirement, Area, and
  Factor reports.

## 2026-06-26

- **Revision**: Implemented durable spec changes for
  [0132 - Remove info finding severity](../changes/archive/0132-remove-info-finding-severity.md).
  Evaluation finding severity is now limited to `critical`, `high`, `medium`,
  and `low`; informational observations use finding `type: note`. The Evaluation
  data, report-tree, and `/quality` skill specs now reflect the reduced severity
  vocabulary and schema/validation behavior.

- **Revision**: Implemented durable spec changes for
  [0130 - Self-contained per-kind data schema](../changes/archive/0130-self-contained-data-schema.md).
  The Evaluation data CLI spec now requires `data schema <kind>` to emit a
  self-contained schema whose required fields and enum value sets are legible
  without following a top-level `$ref` into `$defs`, while allowing the
  no-argument full-surface schema to keep `$defs`/`$ref`. The `/quality` skill
  spec now names `data schema <kind>` as the required-fields and enum-value
  source, with examples as concrete instances only.

- **Revision**: Implemented durable spec changes for
  [0129 - Evaluation orchestration overhaul](../changes/archive/0129-evaluation-orchestration-overhaul.md).
  The `/quality` evaluation spec now removes the evaluation rigor selector,
  makes exhaustive in-scope Requirement coverage mandatory, defaults independent
  collection and QC work to subagent fan-out when available, and requires an
  always-on two-pronged QC phase before roll-up. The parent skill spec and
  evaluate feedback-log spec no longer carry `Rigor:` / `rigor` fields or
  `--rigor` examples.

- **Revision**: Implemented durable spec changes for
  [0128 - Agent-mediated skill alignment](../changes/archive/0128-agent-mediated-skill-alignment.md).
  The `/quality` skill spec now gives read-only orientation a standard
  status-first output shape and requires a non-public recommendation-follow-up
  frame before inspection or mutation. The setup workflow spec now requires the
  setup run frame to be the first element in the opening block and distinguishes
  `QUALITY.md` changes from local workflow feedback-log writes.

- **Revision**: Implemented durable spec changes for
  [0123 - Render interactions through native affordances](../changes/archive/0123-native-interaction-affordances.md).
  The shared [user interaction contract](skills/quality-skill/quality-skill.md)
  now treats each interaction as an intent rendered through a fit-for-purpose
  native affordance when present, with the numbered-option and `y`/`n` renderings
  reframed as the text fallback and a no-double-gate rule for harness-authorized
  mutations. The [`setup`](skills/quality-skill/workflows/setup.md),
  [`update`](skills/quality-skill/workflows/update.md),
  [`evaluation`](skills/quality-skill/evaluation.md), and
  recommendation-follow-up
  ([behavior](skills/quality-skill/recommendation-follow-up.md),
  [guide contract](skills/quality-skill/guides/recommendation-follow-up-md.md))
  specs inherit it, tag their closed choices as single-select intents, and lock
  the human context checkpoint as free text.

- **Revision**: Implemented durable spec changes for
  [0122 - Finding-level candidate actions](../changes/archive/0122-finding-candidate-actions.md).
  The [`/quality` evaluation spec](skills/quality-skill/evaluation.md) now allows
  the skill to record non-binding candidate actions on `gap`/`risk` findings as
  raw material for a later Advise phase, and the
  [reporting spec](skills/quality-skill/reporting.md) states that candidate actions
  are not recommendations and that the v0 report and closeout exclude them.

- **Revision**: Removed active legacy `harnessability` coverage from the
  `/quality` skill specs.
  The Agent Harnessability guide, Top 10 checks, and setup workflow specs now
  require old `harnessability` factors to be reported as stale legacy naming, not
  accepted as current coverage.

- **Revision**: Implemented durable spec changes for
  [0121 - Scannable interaction hierarchy](../changes/0121-scannable-interaction-hierarchy.md).
  The [`/quality` skill](skills/quality-skill/quality-skill.md) interaction
  contract and Decision briefs section now require gates to lead with the
  question, render choices as a visually separated block, fold the non-mutating
  alternative into the stop choice, cap supporting fields at about three, and
  carry hierarchy by position rather than bold alone (surviving bold-stripping).
  The [recommendation follow-up](skills/quality-skill/recommendation-follow-up.md)
  result-report requirement now forbids stacking equally-weighted bold labels and
  requires a primary outcome line.

- **Revision**: Implemented durable spec changes for
  [0120 - String model-identity fields in evaluation data](../changes/archive/0120-string-model-identity-fields.md).
  `SPECIFICATION.md` now reserves `root` as an Area name and carries
  specification version `0.5 (Draft)`. The Evaluation overview and record specs
  now define persisted model identities as canonical qualified-reference strings;
  the lint specs include the `reserved-area-name` rule; the companion
  `quality.schema.json` spec covers the expressible `root` reservation; and the
  `/quality` skill spec describes qualified references in Evaluation artifacts.

- **Revision**: Implemented durable spec changes for
  [0119 - Report header kind prefix and title-first layout](../changes/archive/0119-report-header-kind-prefix.md).
  The Evaluation report-tree spec now requires title-first, kind-prefixed H1s,
  trails below the title, the Model `title` as the root trail element, and no
  standalone `Path:` / `Name:` header lines.

- **Revision**: Implemented durable spec changes for
  [0118 - Report empty-cell marker and legend](../changes/archive/0118-report-empty-cell-legend.md).
  The Evaluation report-tree spec now requires empty scalar cells and paired
  empty components to render as `—`, requires a static per-report legend, and
  keeps worded empty-section rows and not-assessed-style outcomes distinct from
  the marker.

- **Revision**: Implemented durable spec changes for
  [0117 - Requirement report Factors line](../changes/archive/0117-requirement-report-factors-line.md).
  The Evaluation report-tree spec now requires a plural `Factors:` context line
  on Requirement reports and removes `Factors` from the Requirement summary table
  columns.

- **Revision**: Implemented durable spec changes for
  [0116 - Drop the "Evaluation v2" naming](../changes/archive/0116-drop-evaluation-v2-naming.md).
  Renamed the active Evaluation spec bundle from `evaluation-v2/` to
  `evaluation/`, renamed the parent concept to `evaluation.md`, and updated live
  spec links/prose to use plain "Evaluation".

- **Revision**: Implemented durable spec changes for
  [0115 - Type-safe, model-bound Evaluation v2 data](../changes/archive/0115-evaluation-data-typed-contract.md).
  The CLI data and Evaluation record specs now require strict field/type/enum
  validation, model-binding against the run snapshot, dry-run parity, populated
  examples, the `data schema` discovery command, and `data verify`.

- **Revision**: Implemented durable spec changes for
  [0114 - Run frame as first output](../changes/archive/0114-run-frame-first-output.md).
  The `/quality` skill spec's Run frames section now requires the frame as the
  workflow's first output before any tool call, forbids gating emission on a tool
  result, and allows a provisional `resolving…` value for a tool-dependent field —
  carrying the 0096 ordering lesson into the shared contract.

- **Revision**: Implemented durable spec changes for
  [0113 - Evaluation run folder naming](../changes/0113-evaluation-run-folder-naming.md).
  `qualitymd evaluation create` now documents the `NNNN-full-eval` /
  `NNNN-<scope-path>-eval` run-folder grammar, the reserved `quality` slug
  segment, and the `/quality evaluate` convention that scoped runs pass the
  Area/Factor full structural path through `--narrowing`.

- **Revision**: Implemented durable spec changes for
  [0110 - Run frame title and workflow vocabulary](../changes/0110-run-frame-and-workflow-vocabulary.md).
  The `/quality` skill spec's Run frames section now requires the header to name
  the resolved workflow and forbids rendering a command-style `/quality run`
  header or a `Mode:` field label (extending 0062's `Mode: wizard` constraint).
  Across the quality-skill bundle, "mode" is retired as the name for a `/quality`
  workflow in favor of "workflow", with recommendation follow-up described as a
  post-evaluation follow-up that is not a public workflow.

- **Revision**: Implemented durable spec changes for
  [0108 - Short evaluation report filenames](../changes/archive/0108-short-evaluation-report-filenames.md).
  Evaluation v2 report layout and `/quality` reporting specs now keep the root
  Area report at `report.md` while using short subject-aware filenames for
  descendant Area, Factor, and Requirement Markdown reports. The CLI status spec
  now refers to generated Markdown reports generically instead of `report.md`
  bodies.

- **Revision**: Implemented durable spec alignment for
  [0107 - Durable spec alignment](../changes/archive/0107-durable-spec-alignment.md).
  Audited `SPECIFICATION.md` and the active `specs/` bundle for OKF mechanics,
  BCP 14 declarations, change-case-only sections, vocabulary/register drift, and
  index coverage. Updated the authoring sub-guide specs so all live durable specs
  using BCP 14 requirement keywords declare the convention near the top.

- **Revision**: Implemented durable spec changes for
  [0104 - Evaluation v2 report header navigation](../changes/archive/0104-evaluation-v2-report-header-navigation.md)
  and
  [0105 - Evaluation v2 report subject links](../changes/archive/0105-evaluation-v2-report-subject-links.md).
  Updated the Evaluation v2 report-tree spec so generated Markdown reports use
  labeled `Area:` trails, Factor-only `Factor:` trails, compact report-specific
  headers, and linked subject cells instead of generic breadcrumbs, parent
  header links, and repeated generated-report `Details` columns.

- **Revision**: Implemented durable spec changes for
  [0099 - Closed-choice setup UX](../changes/archive/0099-closed-choice-setup-ux.md).
  Updated the parent
  [`/quality` skill](skills/quality-skill/quality-skill.md) interaction
  contract and the
  [`setup`](skills/quality-skill/workflows/setup.md) workflow spec so small
  closed-choice prompts use numbered options, put the recommended answer first,
  and accept `1` as the shortest confirmation. Setup risk discovery now presents
  cost labels to the user and maps them to the existing risk-tolerance meaning.

- **Revision**: Implemented durable spec changes for
  [0097 - Evaluation v2 clean break](../changes/archive/0097-evaluation-v2-clean-break.md).
  The active specs now treat Evaluation v2 as the only runtime evaluation
  contract, require stdin-only `evaluation data set`, remove previous
  evaluation-record and legacy report specs from the active bundle, remove old
  checked-in skill examples, update status/list/project-status contracts to
  report v2 data artifacts and gaps, and bump the bundled specification version
  to `0.4 (Draft)`.

## 2026-06-25

- **Revision**: Updated durable skill workflow specs for
  [0095 - Evaluate feedback log outcomes](../changes/archive/0095-evaluate-feedback-log-outcomes.md).
  The evaluate feedback log remains the process-only workflow artifact under
  `.quality/logs/<timestamp>-evaluate-feedback-log.md`, and its `outcome` field
  now records workflow terminal states rather than report or rating semantics.

- **Creation**: Added the
  [Evaluation v2](evaluation-v2/evaluation-v2.md) durable spec bundle for
  [0094 - Replace evaluation workflow](../changes/0094-replace-evaluation-workflow.md).
  The new bundle captures shared invariants, protocol order, orchestration,
  routine prompt contracts, JSON conventions, data layout, payload kinds, and
  deterministic report tree rules for the replacement workflow.

- **Revision**: Removed the active CLI specs for the legacy
  `evaluation assessment`, `evaluation analysis`, and `evaluation recommendation`
  commands as part of
  [0094 - Replace evaluation workflow](../changes/0094-replace-evaluation-workflow.md).
  New evaluation data is written through
  [`qualitymd evaluation data`](cli/evaluation-data.md).

- **Revision**: Marked the previous
  [Evaluation records](evaluation-records.md) and [reports](reports/index.md)
  specs as superseded legacy contracts. The active evaluation workflow source of
  truth is now [Evaluation v2](evaluation-v2/evaluation-v2.md).

- **Revision**: Implemented durable spec changes for
  [0093 - Named Requirement identity](../changes/archive/0093-requirement-identity.md).
  Requirements now use stable names as map keys, require `title`, retain
  `assessment`, may carry `description`, and have qualified
  `requirement:<area-path>::<requirement-name>` references. Updated the format
  spec, JSON Schema spec, lint-rule catalog, lint-output placeholder, init
  scaffold spec, durable `/quality` authoring-guide mirrors, and reference
  examples.

## 2026-06-24

- **Creation + Revision**: Implemented durable spec changes for
  [0088 - Domain-agnostic corpus alignment](../changes/archive/0088-domain-agnostic-corpus-alignment.md).
  Added the non-software
  [`0002-city-bike-stations-quality-eval`](skills/quality-skill/examples/0002-city-bike-stations-quality-eval/report.md)
  reference example bundle with a complete data-product model, assessment records,
  analysis records, report renders, and recommendations in the same runtime shape as
  `0001`. Updated the example index to mark the corpus as domain-illustrative,
  updated the reporting spec to cross-reference the paired examples and the
  domain-agnostic guide, paired the `report-summary.md` Area Breakdown example with a
  non-software data-product rendering, and added the non-software Top 10 constituent
  bracket through the quality-skill guides spec mirror. No CLI or Go behavior change.

- **Revision**: Implemented durable spec changes for
  [0084 - Agent-mediated UX conformance](../changes/archive/0084-agent-mediated-ux-conformance.md).
  Updated the parent
  [`/quality` skill](skills/quality-skill/quality-skill.md) interaction
  contract and the setup, evaluate, update, recommendation-follow-up, and
  reporting specs so live workflow output follows
  [Designing agent-mediated UX](../docs/guides/agent-mediated-ux.md): status
  first, visually emphasized primary questions or calls to action, scannable
  labels, semantic emoji only, explicit mutation boundaries, and clear closeouts
  with next actions. No QUALITY.md format or CLI behavior change.

## 2026-06-23

- **Creation + Revision**: Implemented durable spec changes for
  [0073 - Evaluation feedback log](../changes/archive/0073-evaluation-feedback-log.md).
  Added the shared
  [workflow feedback log](skills/quality-skill/workflow-feedback-log.md) spec and
  the evaluate adopter
  [feedback log](skills/quality-skill/workflows/evaluate/feedback-log.md) spec.
  Updated the setup adopter spec, evaluation workflow specs, reporting specs,
  evaluation-record specs, and `qualitymd evaluation create` spec so current
  evaluation workflow feedback lives under
  `.quality/logs/<timestamp>-evaluate-feedback-log.md` and `debug-log.md` is
  legacy-compatible only for historical runs.

- **Revision**: Implemented durable spec changes for
  [0072 - Setup context checkpoint](../changes/archive/0072-setup-context-checkpoint.md).
  Updated the
  [`setup`](skills/quality-skill/workflows/setup.md) workflow spec so the final
  human-context discovery inputs are presented as one correction-oriented
  checkpoint covering primary users/outcomes, maintainers/collaborators, other
  stakeholders, and missing or not-agent-accessible context, with omitted
  low-confidence gaps recorded honestly rather than confirmed by silence.

- **Revision**: Implemented durable spec changes for
  [0071 - Setup open-ended review gate](../changes/archive/0071-setup-open-ended-review-gate.md).
  Updated the
  [`setup`](skills/quality-skill/workflows/setup.md) workflow spec so the final
  review recap uses friendly, open-ended wording that preserves the
  `"looks good"` fast path while inviting priorities, worries, wording, edge
  cases, repo-invisible context, and other useful last-call input before
  authoring.

- **Revision**: Implemented durable spec changes for
  [0070 - Setup missing-context provenance](../changes/archive/0070-setup-missing-context-provenance.md).
  Updated the
  [`setup`](skills/quality-skill/workflows/setup.md) workflow spec so
  missing-context discovery preserves provenance: low/no-evidence project facts
  cannot be presented as assumed-understood, excluded gaps require explicit
  setup-provided context or missed agent-accessible evidence, and
  setup-provided context remains visible in the authored model body.

- **Revision**: Implemented durable spec changes for
  [0069 - Setup review gate and discovery trim](../changes/archive/0069-setup-review-gate-and-pedagogy-trim.md).
  Updated the
  [`setup`](skills/quality-skill/workflows/setup.md) workflow spec to remove the
  modeling-rigor and review-posture discovery questions, add a Rating Scale
  confirmation question, trim per-question pedagogy to purpose/context only, and
  make the final recap a hard review gate that waits for an explicit user
  response before authoring.

- **Revision**: Implemented durable spec changes for
  [0068 - Always-on setup feedback log](../changes/archive/0068-always-on-setup-feedback-log.md).
  Updated the
  [workflow feedback log](skills/quality-skill/workflows/setup/feedback-log.md)
  sub-spec so setup creates a feedback log during preflight after CLI support and
  the run frame, updates the current run's file for material workflow-experience
  events, and finalizes it at close with stable frontmatter metadata, lifecycle
  status, a timeline, and no-notable-content notes. Updated
  [`setup`](skills/quality-skill/workflows/setup.md) and the parent
  [`/quality` skill](skills/quality-skill/quality-skill.md) spec to reflect the
  always-on mutation boundary and current-run update rule.

- **Revision**: Implemented durable spec changes for
  [0067 - Setup discovery pedagogy](../changes/archive/0067-setup-discovery-pedagogy.md)
  in [`setup`](skills/quality-skill/workflows/setup.md): required authored
  per-question pedagogy and the teaching-over-round-trips framing, changed the
  prompt-form contract to ask every one of the ten questions every run (removing
  the accept-all-and-skip escape, keeping a per-question fast confirm and
  show-all-at-once), changed the confidence vocabulary from
  `strongly inferred`/`weakly inferred`/`assumed` to `Low`/`Med`/`High` with the
  evidence note retained, and added the final-review-recap stage. The parent
  [`/quality` skill](skills/quality-skill/quality-skill.md) spec's generic
  "confidence-labeled defaults" phrasing was reviewed and left unchanged.

- **Creation + Revision**: Implemented durable spec changes for
  [0066 - Setup feedback log](../changes/archive/0066-setup-feedback-log.md).
  Added the
  [workflow feedback log](skills/quality-skill/workflows/setup/feedback-log.md)
  sub-spec under a new
  [`workflows/setup/`](skills/quality-skill/workflows/setup/index.md) folder,
  defining a hand-authored, skill-only, never-transmitted feedback artifact under
  `.quality/logs/` (`<timestamp>-<workflow>-feedback-log.md`) with an environment
  header, body schema, and redaction rules. Made
  [`setup`](skills/quality-skill/workflows/setup.md) the parent of that sub-spec,
  amended its mutation surface to permit the `.quality/logs/` write while keeping
  every other prohibition in force, and added the close-step authoring. Recorded
  the shared feedback-log artifact and its redaction/no-transmission boundary in
  the parent [`/quality` skill](skills/quality-skill/quality-skill.md) spec so
  later workflows inherit it.

- **Rename + Revision**: Implemented durable spec changes for
  [0065 - Setup discovery and close refinements](../changes/0065-setup-discovery-and-close-refinements.md).
  Renamed the durable `specs/skills/quality-skill/modes/` folder to `workflows/`
  and updated live `modes/` path references in the
  [`/quality` skill](skills/quality-skill/quality-skill.md) parent spec,
  [index](skills/quality-skill/index.md), and
  [evaluation](skills/quality-skill/evaluation.md) specs. Made the durable setup
  prompt-form contract agent-agnostic, added the read-the-scaffold-before-author
  requirement, reframed setup's close contract and the
  [top-10-checks](skills/quality-skill/guides/top-10-quality-md-checks-md.md)
  guide to separate model maturity (`starter`, `immature`, `evaluation-ready`)
  from the CLI's lifecycle `readiness`, and renamed the model-maturity sense of
  "readiness" in the parent spec while keeping CLI/tooling readiness intact.
  Historical log entries keep their `modes/` references frozen.

- **Revision**: Implemented durable spec alignment for
  [0064 - Structured setup workflow](../changes/0064-structured-setup-workflow.md).
  Updated durable [`/quality` skill](skills/quality-skill/index.md) specs so
  setup is described as a workflow, runtime setup guidance must be an operator
  playbook, the setup brief and ten concrete discovery questions are specified,
  and setup still writes only `QUALITY.md`.

- **Revision**: Clarified agentic use context in the durable
  [`/quality` skill](skills/quality-skill/index.md) specs and installable skill
  metadata. The specs now preserve AI assistant and coding-agent workflow
  language as use context while keeping the modeled quality domain agnostic.

- **Revision**: Implemented durable spec alignment for
  [0063 - Contextual setup flow](../changes/0063-contextual-setup-flow.md).
  Updated durable [`/quality` skill](skills/quality-skill/index.md) specs so
  setup analyzes context, asks confidence-labeled discovery questions, writes
  only `QUALITY.md`, validates and readiness-checks the model, offers next-step
  choices, and no longer writes quality-log entries.

- **Revision**: Implemented durable spec alignment for
  [0062 - Remove wizard mode](../changes/0062-remove-wizard-mode.md). Updated the
  durable [`/quality` skill](skills/quality-skill/index.md) specs to remove
  `wizard` as a public mode, keep bare `/quality` as read-only orientation,
  avoid promoting `status` or `next` as public invocations, and recommend only
  public workflows from ambiguous requests.

- **Revision**: Implemented durable spec alignment for
  [0061 - Natural scope labels](../changes/archive/0061-natural-scope-labels.md).
  Updated the durable [`/quality` skill](skills/quality-skill/index.md) specs to
  make natural Area and Factor labels the primary scoped-evaluation input,
  preserve qualified model references for exact addressing, require targeted
  repeated-Factor clarification, and keep durable evaluation artifacts on stable
  identifiers.

## 2026-06-22

- **Revision**: Implemented durable spec alignment for
  [0060 - Friendly path display](../changes/archive/0060-friendly-path-display.md).
  Updated [`SPECIFICATION.md`](../SPECIFICATION.md),
  [`report-summary.md`](reports/report-summary-md.md),
  [`report.md`](reports/report-md.md), [`report.json`](reports/report-json.md),
  [evaluation report outputs](evaluation-records/report-outputs.md), and durable
  [`/quality` skill](skills/quality-skill/index.md) specs to distinguish human
  display values from model references, render the root Area path as `/` in
  human Markdown reports, and keep report JSON identifiers structured or
  reference-compatible.

- **Revision**: Implemented durable spec alignment for
  [0059 - Unqualified model references](../changes/archive/0059-unqualified-model-references.md).
  Updated [`SPECIFICATION.md`](../SPECIFICATION.md),
  [`report-summary.md`](reports/report-summary-md.md),
  [`report.md`](reports/report-md.md), [`report.json`](reports/report-json.md),
  [evaluation report outputs](evaluation-records/report-outputs.md), and durable
  [`/quality` skill](skills/quality-skill/index.md) specs to define
  unqualified references for fixed-type contexts, render unqualified Area
  Breakdown paths in human reports, and keep durable machine artifacts
  structured.

- **Revision**: Implemented durable spec alignment for
  [0058 - Model reference identifiers](../changes/archive/0058-model-reference-identifiers.md).
  Updated [`SPECIFICATION.md`](../SPECIFICATION.md),
  [`report-summary.md`](reports/report-summary-md.md),
  [`report.md`](reports/report-md.md), [`report.json`](reports/report-json.md),
  [evaluation report outputs](evaluation-records/report-outputs.md),
  [`qualitymd lint`](cli/lint.md), [`qualitymd lint rules`](cli/lint-rules.md),
  the [`quality.schema.json`](quality-schema-json.md) artifact spec, and durable
  [`/quality` skill](skills/quality-skill/index.md) specs to define strict Area
  names, Factor names, Rating Level IDs, canonical typed model references,
  edge-only shorthand, structured path preservation, and the revised Area
  Breakdown columns.

- **Revision**: Implemented durable spec alignment for
  [0057 - Quality data directory](../changes/0057-quality-data-directory.md).
  Updated [`qualitymd evaluation create`](cli/evaluation-create.md),
  [`qualitymd status`](cli/status.md), [`qualitymd lint`](cli/lint.md),
  [`qualitymd lint rules`](cli/lint-rules.md), the
  [`quality.schema.json`](quality-schema-json.md) artifact spec, and durable
  [`/quality` skill](skills/quality-skill/index.md) specs to define the
  QUALITY.md workspace, `.quality/` quality data directory,
  `.quality/evaluations/` default, `.quality/log/` quality log, root `config`
  tooling pointer, and strict lint handling for unknown keys.

- **Revision**: Clarified prospective evaluation design and plan artifacts for
  [0056 - Prospective evaluation plan artifacts](../changes/archive/0056-prospective-evaluation-plan-artifacts.md).
  Added the [design.md](evaluation-records/design-md.md) artifact spec, updated
  [plan.md](evaluation-records/plan-md.md) so the initial plan records intended
  execution before assessment begins, and aligned the durable
  [/quality evaluation workflow](skills/quality-skill/evaluation.md),
  [`evaluate` mode](skills/quality-skill/modes/evaluate.md),
  [reporting](skills/quality-skill/reporting.md), and runtime skill guidance
  with explicit plan amendments for later scope, coverage, rigor, or evidence
  strategy changes.

- **Revision**: Added the structured-input contract for
  [0055 - Self-describing evaluation record input](../changes/0055-evaluation-input-ergonomics.md).
  The CLI parent spec now requires payload-documenting help, `-n/--dry-run`
  validation without persistence, and aggregated JSON-key validation errors for
  structured payload commands. The evaluation assessment, analysis,
  recommendation, create, and status command specs now carry the corresponding
  record-write and `plan.md` coverage-discovery behavior, while the evaluation
  record specs document analysis-set whole-set behavior, surfaced-example drift
  protection, and planned coverage diagnostics.

- **Revision**: Removed `/quality improve` as a durable mode contract for
  [0054 - Remove improve mode](../changes/archive/0054-remove-improve-mode.md).
  Added
  the non-mode
  [`recommendation follow-up`](skills/quality-skill/recommendation-follow-up.md)
  contract for apply-now and issue-tracker handoff outcomes, updated the parent
  [`/quality` skill spec](skills/quality-skill/quality-skill.md), evaluation
  workflow, reporting, quality log, wizard/evaluate/update mode specs, and mode
  index to route active recommendations through follow-up rather than an
  `improve` mode. Added the runtime
  [`recommendation follow-up guide`](skills/quality-skill/guides/recommendation-follow-up-md.md)
  contract.

- **Revision**: Split remaining large durable specs for
  [0053 - Align remaining durable specs](../changes/archive/0053-align-remaining-durable-specs.md).
  `evaluation-records.md` now keeps shared responsibility, runtime-not-OKF,
  schema-version, and historical compatibility rules while child specs under
  [`evaluation-records/`](evaluation-records/index.md) own run folders,
  assessment records, analysis records, `plan.md`, `debug-log.md`,
  recommendation records, and report-output invariants. Split lint's rule system
  into [`lint rules`](cli/lint-rules.md) and its finding/output contract into
  [`lint output`](cli/lint-output.md). Split the cross-command
  [`update notice`](cli/update-notice.md) from the explicit
  [`update`](cli/update.md) command.

- **Revision**: Aligned the durable `/quality` skill specs with the
  repository capitalization convention for QUALITY.md vocabulary. The skill,
  evaluation, guide, setup, and evaluate specs now reserve capitalized model
  vocabulary for formal type-name use and prefer lowercase in ordinary
  instructional and user-facing prose.

- **Revision**: Aligned the durable `/quality` skill specs with the
  artifact-spec versus behavioral-component guidance for
  [0052 - Durable spec alignment](../changes/0052-durable-spec-alignment.md).
  Added mode behavior specs under
  [`skills/quality-skill/modes/`](skills/quality-skill/modes/index.md) for
  setup, wizard, evaluate, improve, and update; added cross-mode component specs
  for the [`evaluation workflow`](skills/quality-skill/evaluation.md),
  [`reporting`](skills/quality-skill/reporting.md), and
  [`quality log`](skills/quality-skill/quality-log.md); narrowed the parent
  [`/quality` skill spec](skills/quality-skill/quality-skill.md) to shared
  contracts plus component summaries; and updated the skill-spec index to
  distinguish the parent spec, component specs, guide artifact specs, examples,
  and the installable prompt.

- **Revision**: Added a `## Quality log` section to the
  [`/quality` skill spec](skills/quality-skill/quality-skill.md) for
  [0050 - Quality log](../changes/archive/0050-quality-log.md). The section specifies the
  convention-first quality log under `quality/log/`: date-named one-file-per-change
  entries, runtime-not-OKF status, the meaningful-change criterion, which modes
  write (`setup` seeds, `improve` appends) and reconcile (`wizard`, read-only), and
  the run-frame mutation accounting. Promoted the curated-not-complete and
  date-naming rationale into the section and recorded the CLI graduation (a
  `qualitymd log` command, a `logDir` config key, a standalone artifact-spec, and a
  queryable index) as deferred.

- **Creation**: Added the companion JSON Schema artifact-spec
  [`quality.schema.json`](quality-schema-json.md) and the
  [`qualitymd schema`](cli/schema.md) command sub-spec for
  [0049 - Companion JSON Schema](../changes/0049-companion-json-schema.md). The
  artifact-spec carries the durable contract for the structural-only,
  non-normative schema derived from `internal/schema` (its no-drift property,
  draft 2020-12 dialect, stable `$id`, and open-extension rule), while the
  command sub-spec specifies the verbatim-artifact `schema` command. Registered
  both in [`specs/index.md`](index.md) and [`specs/cli/index.md`](cli/index.md).

- **Revision**: Applied the 1:1 artifact-spec filename convention to the
  `/quality` runtime guide contracts:
  [`authoring.md`](skills/quality-skill/guides/authoring-md.md),
  [`getting-started.md`](skills/quality-skill/guides/getting-started-md.md), and
  [`top-10-quality-md-checks.md`](skills/quality-skill/guides/top-10-quality-md-checks-md.md).
  Runtime guide artifact filenames remain unchanged under `skills/quality/guides/`.

- **Creation**: Added 1:1 generated report artifact specs for
  [0048 - Area factor report breakdown](../changes/0048-area-factor-report-breakdown.md):
  [`report-summary.md`](reports/report-summary-md.md),
  [`report.md`](reports/report-md.md), and [`report.json`](reports/report-json.md).
  The existing [Evaluation records](evaluation-records.md) spec remains the
  shared report-model and record contract, while the new specs own each
  artifact's shape, labels, and density.

## 2026-06-21

- **Revision**: Added the process-only `debug-log.md` evaluation artifact for
  [0046 - Evaluation debug log](../changes/archive/0046-evaluation-debug-log.md).
  Updated [Evaluation records](evaluation-records.md),
  [`evaluation create`](cli/evaluation-create.md), and the
  [/quality skill contract](skills/quality-skill/quality-skill.md) so new runs
  seed a debug log while reports, ratings, findings, recommendations, and next
  actions remain derived from formal evaluation records.

- **Revision**: Clarified Markdown-body semantics and skill-guide contracts for
  [0045 - Evaluable body context](../changes/0045-evaluable-body-context.md).
  The body is now described as evaluable judgment context for building, using,
  justifying, and evaluating model quality; body sections should be concise,
  self-explanatory, progressively disclosed, and grounded in agent-accessible
  support; material support that is not agent-accessible is captured in the
  relevant section's unknowns or open questions.

## 2026-06-20

- **Revision**: Replaced the standalone Known gaps body section with per-section
  unknowns, open questions, and a human/agent review state line across the
  [format spec](../SPECIFICATION.md), [`init`](cli/init.md),
  the [authoring guide spec](skills/quality-skill/guides/authoring-md.md),
  the [top-10 checks spec](skills/quality-skill/guides/top-10-quality-md-checks-md.md),
  and the [`/quality` skill spec](skills/quality-skill/quality-skill.md). The
  recommended body sections are now Overview, Scope, Needs, and Risks; each
  records its own unknowns (broad uncertainty) and open questions (specific,
  one-answer), distinct from a `not assessed` result. Per
  [0044](../changes/0044-section-unknowns-open-questions.md).

- **Revision**: Updated the [Evaluation records](evaluation-records.md),
  [`evaluation list`](cli/evaluation-list.md),
  [`evaluation status`](cli/evaluation-status.md),
  [`evaluation report`](cli/evaluation-report.md), [`status`](cli/status.md),
  and [`/quality` skill](skills/quality-skill/quality-skill.md) contracts for
  evaluation-history compatibility. Current record writers stay strict, while
  history/status/list readers surface malformed, unsupported, or incomplete
  historical records as non-reportable run gaps instead of breaking ordinary
  workflows.

- **Revision**: Updated the [Evaluation records](evaluation-records.md),
  [`evaluation status`](cli/evaluation-status.md),
  [`evaluation report`](cli/evaluation-report.md), and
  [`/quality` skill](skills/quality-skill/quality-skill.md) contracts for typed
  evaluation-report states. Rating results, local target ratings, record
  lifecycle, next steps, missing metadata, finding severity, and run gaps now
  have explicit typed-state contracts, while evidence kind remains intentionally
  open display/grouping metadata.

- **Revision**: Updated the [Evaluation records](evaluation-records.md),
  [`evaluation report`](cli/evaluation-report.md), and
  [`/quality` skill](skills/quality-skill/quality-skill.md) contracts so all
  report outputs are faithful projections of one assembled report model.
  Clarified headline verdict sourcing, scoped-output labeling,
  active/superseded handling, typed rating states, deterministic finding/action
  labels, and the renderer's no-rejudging boundary.

- **Revision**: Completed the
  [`Sparrow Payments`](skills/quality-skill/examples/0001-subject-quality-eval/report.md)
  example as a reportable generated fixture, with full assessment and analysis
  records for the run and checked-in reports regenerated from those records.

## 2026-06-19

- **Revision**: Added [`qualitymd update`](cli/update.md), replacing
  `qualitymd upgrade` with apply-by-default update behavior, `--check` advisory
  mode, release readiness, release notes, managed standalone self-apply, and the
  cached ambient update notice. Updated the [CLI overview](cli.md), command
  index, and `/quality` skill contract for the `/quality update` mode rename.

- **Revision**: Updated the evaluation report summary contract for
  [0040 - Readable report summary](../changes/0040-readable-report-summary.md):
  `report-summary.md` now uses a decision-brief outline, "Full evaluation",
  "Overall rating", and prominent Recommendation IDs while preserving the
  existing `report.json` schema.

- **Revision**: Clarified guide boundaries for the `/quality` skill guides:
  authoring now carries reusable best practices, getting-started is the
  first-run process/outcomes guide, and the Top 10 checklist routes against that
  split.

- **Creation**: Added the
  [`Top 10 QUALITY.md checks`](skills/quality-skill/guides/top-10-quality-md-checks-md.md)
  guide contract and wired the `/quality` wizard spec to use it for bounded
  model/lifecycle inspection findings.

- **Revision**: Clarified that the
  [`getting-started`](skills/quality-skill/guides/getting-started-md.md) guide's
  Known gaps section includes known unknowns: missing context, unresolved
  questions, and evidence gaps.

- **Revision**: Added desired outcomes for each Markdown body section in the
  [`getting-started`](skills/quality-skill/guides/getting-started-md.md) guide
  contract.

- **Revision**: Updated the
  [`getting-started`](skills/quality-skill/guides/getting-started-md.md) guide
  contract so the rating scale follows the Markdown body before the rest of the
  model tree is expanded.

- **Revision**: Updated the
  [`getting-started`](skills/quality-skill/guides/getting-started-md.md) guide
  contract so first-run authoring fills the Markdown body before building out
  the quality model tree.

- **Revision**: Added runtime guide structure for the
  [`/quality` skill](skills/quality-skill/quality-skill.md): moved the
  authoring guide contract to
  [`guides/authoring`](skills/quality-skill/guides/authoring-md.md), added the
  [`getting-started`](skills/quality-skill/guides/getting-started-md.md) guide
  contract for first-run model population after `qualitymd init`, and updated
  the skill spec to route setup to the new guide.

- **Revision**: Replaced the evaluation CLI command specs with
  [`evaluation create`](cli/evaluation-create.md),
  [`evaluation list`](cli/evaluation-list.md),
  [`evaluation status`](cli/evaluation-status.md),
  [`evaluation assessment`](cli/evaluation-assessment.md),
  [`evaluation analysis`](cli/evaluation-analysis.md),
  [`evaluation recommendation`](cli/evaluation-recommendation.md), and
  [`evaluation report`](cli/evaluation-report.md). Updated the CLI overview,
  command index, [Evaluation records](evaluation-records.md), and the
  [`/quality` skill spec](skills/quality-skill/quality-skill.md) for the
  noun/verb evaluation surface, `plan.md` coverage frontmatter, batched writes,
  `--latest`, and separate report gate.

- **Revision**: Updated [`evaluation build-report`](cli/evaluation-build-report.md),
  [Evaluation records](evaluation-records.md), and the
  [`/quality` skill spec](skills/quality-skill/quality-skill.md) so human
  reports and skill prose use required model titles as primary labels while
  evaluation records, `report.json`, and gates preserve stable identifiers.

- **Revision**: Updated the
  [`/quality` skill spec](skills/quality-skill/quality-skill.md) for
  `/quality upgrade` as a first-class maintenance mode that plans paired
  skill/CLI upgrades, asks before mutation, delegates to owner tools, verifies
  the visible CLI, and warns about skill reload requirements.

- **Revision**: Updated the
  [`/quality` skill spec](skills/quality-skill/quality-skill.md) for
  project-owned Agent Skills metadata: `metadata.version`,
  `metadata.requires-qualitymd-cli`, matching `compatibility` prose, and the
  released-install prerequisite range.

- **Revision**: Added [`qualitymd version`](cli/version.md) and
  [`qualitymd upgrade`](cli/upgrade.md) command specs, updated the
  [CLI overview](cli.md) and [command index](cli/index.md), and specified
  structured local version metadata plus explicit, opt-in upgrade checks.

- **Revision**: Updated [`evaluation build-report`](cli/evaluation-build-report.md),
  [Evaluation records](evaluation-records.md), and the
  [`/quality` skill spec](skills/quality-skill/quality-skill.md) for generated
  `report-summary.md` alongside `report.md` and `report.json`. Added a
  representative summary artifact to the
  [skill examples](skills/quality-skill/examples/index.md).

- **Revision**: Updated the root [`SPECIFICATION.md`](../SPECIFICATION.md),
  [`lint`](cli/lint.md), [`init`](cli/init.md), [`status`](cli/status.md),
  [`evaluation build-report`](cli/evaluation-build-report.md), and
  [`/quality` authoring guide spec](skills/quality-skill/guides/authoring-md.md) for
  required display titles on Models, Targets, Factors, and Rating Levels.

- **Revision**: Updated the durable
  [`/quality` skill spec](skills/quality-skill/quality-skill.md) to align CLI
  prerequisite checks with the new versioning policy: released skill installs
  use the declared `qualitymd` SemVer range, while local development builds can
  still be accepted by command-surface probing.

- **Revision**: Reframed the durable `/quality` wizard contract as a
  read-only quality wayfinder: it now probes setup/model/history state, classifies
  lifecycle readiness, recommends one next workflow, and offers concrete
  alternatives without producing ratings or editing files.

- **Revision**: Updated the durable
  [`/quality` skill spec](skills/quality-skill/quality-skill.md) to settle the
  skill's bundled runtime resources and require deep-evaluation subagent prompts
  to carry scope, relevant requirements, secret-handling, source-as-data, and
  structured-findings-only instructions.

- **Revision**: Normalized active specs to the RFC 8174 / BCP 14 convention and
  reduced uppercase requirement keywords where the text was advisory,
  illustrative, or presentation guidance rather than a conformance switch.
  Updated the root format spec, CLI specs, evaluation record contract, and
  `/quality` skill specs to reserve BCP 14 wording for testable requirements.

## 2026-06-18

- **Creation**: Added the [QUALITY.md authoring guide](skills/quality-skill/guides/authoring-md.md) sub-spec — the 1:1 contract for
  the skill's [`authoring.md`](../skills/quality/guides/authoring.md)
  resource. States its purpose (a single comprehensive guide to understanding and
  working with `QUALITY.md` files), resolves the embed-vs-drift tension by scoping
  the guide as a deliberately self-contained reference resource that restates the
  format (an explicit exception to the prompt's no-drifting-copy rule) paid for by
  a conformance duty to [`SPECIFICATION.md`](../SPECIFICATION.md), and fixes its
  structural conventions (single-level concepts, reference + how-to per chapter,
  directive/job form, authoring order). Promoted the authoring guide out of the
  skill spec's **Deferred** resources note and added the one-line conformance note
  and purpose statement to the guide itself.

- **Revision**: Updated the
  [`lint`](cli/lint.md) spec for mandatory factor references: added
  `missing-factor-reference`, made `unknown-factor` apply to factor references
  generally, and aligned `empty-factor` terminology with requirement
  references.

- **Revision**: Updated the durable
  [`/quality` skill spec](skills/quality-skill/quality-skill.md) to allow the
  installable skill to keep `SKILL.md` as the always-loaded router while moving
  mode-specific procedures into `skills/quality/modes/` and supporting docs into
  `skills/quality/resources/`.

- **Revision**: Removed the public bundled-model workflow. Deleted the
  `qualitymd models` command spec, removed it from the CLI command catalog,
  made `evaluation create-run` subject-only with no `--altitude` flag, and
  updated the evaluation record contract to treat `model` altitude as historical
  only. Synced the durable `/quality` skill spec around subject evaluation plus
  guide-backed QUALITY.md authoring.

- **Rationale capture**: Made durable the motivating learnings behind the
  in-review evaluation changes (0012–0024), so the *why* survives when the change
  records are archived. Added inline rationale to
  [Evaluation records](evaluation-records.md) (standalone-spec reason; CLI-owned
  scan-derived numbering and collision-as-error from the real run-number
  collision; CLI `model.md` snapshot ownership; duplicate-assessment correction
  hazard; assessment-vs-recommendation superseding asymmetry; explicit-intent
  superseding; no-new-schema-field constraint for evidence/route hints),
  [`build-report`](cli/evaluation-build-report.md) (renderer secret /
  prompt-injection trust boundary; summary-first shape provenance; the bounded
  `design.md`/`plan.md` summary-source contract; dotted-path and idempotency
  rationale), [`add-record`](cli/evaluation-add-record.md) (input/decode rules,
  CLI-owned-field rejection, exit-code mapping, run-snapshot rating validation,
  collision retry, deterministic generated rendering, subcommand-per-kind),
  [`create-run`](cli/evaluation-create-run.md) (partial-artifact prevalidation
  reason and prevalidate-over-rollback), and the
  [`/quality` skill spec](skills/quality-skill/quality-skill.md) (re-run vs.
  re-read re-check rationale; superseding correction rationale).

- **Revision**: Added assessment superseding for evaluation runs. Corrected
  assessment result records may now carry `supersedes`; status reports dangling,
  cross-requirement, and stale-analysis superseding references; reports
  distinguish active from superseded assessment results.

- **Revision**: Tightened `evaluation create-run --subject` validation. Invalid
  subject paths now fail before creating the evaluation directory or a numbered
  run folder, and subject paths must resolve to files rather than directories.

- **Revision**: Added recommendation superseding for evaluation runs. Corrected
  recommendation records may now carry `supersedes`, status reports dangling
  superseding references, and reports distinguish active from superseded advice
  when choosing Next Action.

- **Revision**: Added optional planned coverage metadata for evaluation runs.
  `qualitymd evaluation set-planned-coverage` writes run-root
  `planned-coverage.json`, and `show-status` now reports missing planned
  assessment results, missing planned analyses, and unexpected assessment result records when that
  artifact exists.

- **Revision**: Tightened evaluation reportability for duplicate assessment
  records. Runs with more than one assessment for the same ordered `targetPath`
  and `requirement` now report a `duplicate-assessment` gap in `show-status`,
  and `build-report` refuses to render them.

- **Revision**: Tightened evaluation reportability: `show-status` and
  `build-report` now require exactly one in-scope root analysis record with an
  empty `targetPath`, so a child target cannot silently become the report
  headline.

- **Revision**: Updated the evaluation report contract for change
  [0018](../changes/archive/0018-evaluation-report-ux.md). Generated reports are now
  summary-first for human readers and expose non-null scope, target summaries,
  evidence basis, limitations, empty recommendation arrays, explicit rating
  objects, and structural grouping-target states in `report.json`.

- **Refinement**: Clarified that equivalent summary limitations are
  deterministically deduplicated across recorded context and rationale-derived
  constraints while preserving the first displayed wording.

- **Refinement**: Clarified that derived limitation summaries preserve
  locator-like text such as dotted file paths.

- **Creation**: Added the durable
  [Evaluation records](evaluation-records.md) contract and the evaluation CLI
  sub-specs for
  [`create-run`](cli/evaluation-create-run.md),
  [`add-record`](cli/evaluation-add-record.md),
  [`show-status`](cli/evaluation-show-status.md), and
  [`build-report`](cli/evaluation-build-report.md). Updated the CLI overview,
  command index, and schema type description for the shared runtime record
  contract.

- **Revision**: Synced the durable
  [`/quality` skill spec](skills/quality-skill/quality-skill.md), installable
  skill prompt, skill guide, and reference recommendation examples with the new
  evaluation CLI surface. The skill now delegates run scaffolding, record writes,
  renderability checks, and report rendering to `qualitymd evaluation ...`, while
  retaining judgment and the stricter evidence/re-check/rigor rules.

- **CLI output polish**: Updated the CLI baseline to apply the `NO_COLOR` and
  non-terminal plain-output gate across stdout and stderr human surfaces. Updated
  `qualitymd lint` to require deterministic next actions for invalid results and
  a matching stderr footer, and updated `qualitymd models list` to allow
  terminal-only styling over the unchanged plain table.

## 2026-06-17

- **Revision**: Synced the [`/quality` skill spec](skills/quality-skill/quality-skill.md)
  with its implementation path and raw runtime artifact contract. The spec now
  points to [`skills/quality/SKILL.md`](../skills/quality/SKILL.md), uses the
  finalized trigger description, defines setup and model-altitude evaluation via
  `qualitymd models view quality-meta-model --source`, adds `.quality/config.yaml`
  `evaluationDir`, requires JSON assessment/analysis records plus `report.json`,
  and removes runtime artifact concept types from [`schema.md`](schema.md).

- **Creation**: Added the [`qualitymd models`](cli/models.md) command spec and
  listed it in the [CLI spec](cli.md). The command exposes bundled models through
  `models list` and `models view <name>`, with JSON forms for agents/tools and a
  `--source` rewrite for model-altitude evaluation with the bundled
  `quality-meta-model`.

- **Revision**: Strengthened the [`0001` evaluation report](skills/quality-skill/examples/0001-subject-quality-eval/report.md)
  against the elements ISO/IEC 25040:2024 §5.6.3.2 lists for an evaluation
  report, while staying within the format spec's "presents at least" Report
  contract. Added an **Evaluated** provenance line (source commit, model
  revision, evaluator, date, assessment inputs) so the verdict is reproducible; a
  **Limitations** section (rigor ceiling, point-in-time secret scan, single-test
  confidence) qualifying the ratings without changing them, kept distinct from
  Scope exclusions and *not assessed* outcomes; and sharpened the committed-credential
  rationale to trace measure → applied criterion against the requirement's own
  `ratings` overrides. Declared the new `9f2c1ab` commit as a shared fictional
  locator in the [examples index](skills/quality-skill/examples/index.md).

- **Revision**: Extended the [`0001` example bundle](skills/quality-skill/examples/0001-subject-quality-eval/report.md)
  with a **two-level nested area** to exercise multi-level roll-up: a
  **Webhooks** child area (sibling of Ledger, source `./webhooks`, a signing
  requirement under a refined **Security** factor) with a **Delivery**
  sub-area (`./webhooks/delivery`, retry + redelivery-suppression under a
  refined **Reliability** factor) in
  [`model.md`](skills/quality-skill/examples/0001-subject-quality-eval/model.md).
  The Delivery deduplication requirement is *rated* **Minimum**, so the Webhooks
  aggregate (**Minimum**) falls below its own local rating (**Target**) — the
  bundle's first intermediate aggregate that differs from its local — and the
  root's counterfactual now layers: rotating the committed credential lifts the
  root only to **Minimum** (the webhook gap then binds), reaching **Target**
  once that closes too. Updated the
  [report](skills/quality-skill/examples/0001-subject-quality-eval/report.md)
  rationale/scope/advice, revised recommendation
  [001](skills/quality-skill/examples/0001-subject-quality-eval/recommendations/001-rotate-committed-gateway-key.md)'s
  done-criterion to match, added recommendation
  [003](skills/quality-skill/examples/0001-subject-quality-eval/recommendations/003-bound-webhook-dedup-window.md),
  and refreshed the [examples index](skills/quality-skill/examples/index.md).

- **Revision**: Polished the [`0001` example bundle](skills/quality-skill/examples/0001-subject-quality-eval/report.md)
  for readability and to track the format spec's new **target display fields**:
  gave the **Ledger** child target a `title` and `description` (the report
  already displayed "Ledger", which the bare map key did not provide); added an
  "At a glance" condensed model tree above the verbatim YAML in
  [`model.md`](skills/quality-skill/examples/0001-subject-quality-eval/model.md),
  mirroring the specification's sample-report tree; centralized the facts shared
  across the bundle (fictional subject/locators, the suggested four-level scale,
  and that `model.md` is the evaluated file rather than a runtime output) into
  the [examples index](skills/quality-skill/examples/index.md), trimming the
  repeated per-file admonitions to a pointer; stated an explicit
  model → report → recommendations reading order in the index; and added a
  cross-area secondary factor — the Ledger's double-entry requirement now
  tags the root **Reliability** factor, so the example exercises a factor lens
  ranging wider than the local rating (the spec's Analyze/Report secondary-factor
  notes).

- **Revision**: Updated [`lint`](cli/lint.md) for target display fields:
  `misplaced-root-key` now documents only `ratingScale` on a Target, matching the
  format spec's Model/Target distinction after target `title` and `description`
  became valid Target properties.

- **Restructure + Creation**: Gave the [`/quality`](skills/quality-skill/quality-skill.md)
  skill its own OKF folder — moved the spec into
  [`skills/quality-skill/`](skills/quality-skill/index.md) (as a named concept
  beside a new folder `index.md`, keeping `index.md` a frontmatter-free listing
  per OKF) and rewrote its relative links for the new depth. Added an
  [`examples/`](skills/quality-skill/examples/index.md) subfolder with the first
  worked reference instance of the skill's [Reporting](skills/quality-skill/quality-skill.md#reporting)
  contract: a whole-model **subject** evaluation of a fictional "Sparrow
  Payments" service held at **Unacceptable** by a committed live gateway
  credential — a [report](skills/quality-skill/examples/0001-subject-quality-eval/report.md)
  plus two standalone
  [recommendations](skills/quality-skill/examples/0001-subject-quality-eval/recommendations/001-rotate-committed-gateway-key.md)
  with done-criteria, plus the
  [model evaluated](skills/quality-skill/examples/0001-subject-quality-eval/model.md)
  reproduced alongside (an `Example Model` concept embedding the fictional
  `QUALITY.md` — root **Security** with **Secrets handling**/**Access control**
  sub-factors and **Reliability**, a **Ledger** child target with
  **Correctness**, and a `ratings` override on the secrets requirement the
  report's rationale and done-criterion lean on) so every finding traces to a
  declared requirement and `source`. The example exercises `file:line`
  evidence, the
  secret-by-reference rule (credential cited by location and type, value
  withheld, rotation recommended), a prompt-injection comment recorded as a
  finding and treated as data, and a *not assessed* requirement excluded from
  the roll-up but noted. Registered `Evaluation Report` and `Recommendation`
  concept types in [`schema.md`](schema.md) and linked the example from the
  skill's Reporting section.

- **Revision**: Reconciled the structural schema source of truth with
  [`SPECIFICATION.md`](../SPECIFICATION.md) and [`lint`](cli/lint.md): the public
  format spec now notes the typed schema declaration consumed by lint and the
  drift test that compares it to the YAML schema snippets; the model `title`
  snippet now matches its recommended status; and the lint sub-spec records that
  structural validation derives from the single schema declaration rather than a
  second valid-key list.
- **Revision**: Added the [CLI spec](cli.md)'s **Agent accessibility** section:
  a non-opt-in baseline for in-scope commands (non-interactivity,
  stdout-is-payload/stderr-is-diagnostics, determinism, plain non-TTY output,
  and categorized exit codes) plus opt-in capabilities (`--json`, `nextActions`,
  and deferred quiet/verbosity). Documented exit codes `0` success, `1`
  ran-but-found-problems, `2` usage error, and `70` internal/could-not-complete
  error; revised the `--json` convention to a SHOULD-by-default with a
  verbatim-artifact carve-out and result receipts for side-effecting commands;
  and updated [`init`](cli/init.md) to specify its JSON receipt, `--json -`
  usage error, and JSON overwrite-refusal object.
- **Revision**: Reframed the [`/quality`](skills/quality-skill/quality-skill.md) skill's
  relationship to the format spec's evaluation from **deference** to
  **conformance**: the skill now *owns and specifies* its evaluation process
  (this spec, its prompt, and the CLI together) rather than pointing at
  [`SPECIFICATION.md`](../SPECIFICATION.md) to "run." Added a **Conformance to
  the format spec** subsection stating that every evaluation MUST conform to the
  spec's Evaluation contract (assessment → finding → rating, *not assessed* over
  guessing, inferred weighted roll-up, required report contents) while the skill
  remains one *implementation* free to specify its own workflow, ordering,
  heuristics, rigor levels, and artifacts — with the spec as the **conformance
  target** that governs on divergence. Scoped the existing "don't embed the
  format" rule to the *format/schema rules and rating vocabulary* (still grounded
  at runtime from `qualitymd spec`), explicitly excluding the *evaluation
  process*, which the skill carries.
- **Revision**: Tightened [`lint --fix`](cli/lint.md)'s in-place write contract
  to avoid ambiguous symbolic-link replacement: repair should refuse a linted
  symlink path until symlink write semantics are specified.
- **Revision**: Scoped in [`lint --fix`](cli/lint.md) for deterministic in-place
  repair of fixable findings. The durable lint spec now defines the `--fix`
  flag, transactional per-file repair behavior, preservation requirements,
  post-repair linting, exit-status semantics, and JSON repair reporting through
  `summary.fixed` and `repairs`, while leaving patch/full-file repair output
  modes deferred.
- **Revision**: Closed the lint-specific open items in
  [`lint`](cli/lint.md): added scope/deferred boundaries, specified that this
  phase has no command-specific flags beyond inherited `--json`, defined the
  `--json` finding document and per-finding fields, made `modelPath` the stable
  location contract with optional source positions, set the minimum human-output
  content, and fixed deterministic finding ordering plus downstream-rule
  blocking behavior; the sub-spec now carries a neutral inheritance note instead
  of a draft placeholder. Automatic repair, suppression, rule selection, severity
  overrides, and a lint-emitted rule catalog remain deferred.
- **Revision**: Added fixability criteria to the [`lint`](cli/lint.md) rule
  authoring contract — single mechanical edit, no presumed intent,
  content-preserving, strictly improving, idempotent, local/explainable, and
  stable formatting — and expanded the rule inventory with `Fixable` and
  `Fixable rationale` columns. The initial inventory marks only `empty-property`
  as fixable; placeholder insertion remains scaffold behavior owned by
  [`init`](cli/init.md), not a lint repair.
- **Revision**: Fleshed out the stub [`/quality`](skills/quality-skill/quality-skill.md) skill
  spec into a draft, working through the unblocked TODOs (inspired by the
  basecamp and shadcn/improve `SKILL.md` patterns). Promoted the use-case sketch
  to normative requirements — resolving the open question as *evaluate is
  read-only; improve recommends and applies only on explicit confirmation* — and
  added: a **Boundaries and hard rules** section (judgment-vs-CLI division,
  evaluated content is untrusted data not instructions, never reproduce secrets,
  scoped-result-is-not-a-whole-model-verdict, determinism); **Invocation**
  (skill frontmatter/metadata kept in sync with the format via `qualitymd spec`
  rather than hard-coded, plus five resolved arguments — mode
  `evaluate`/`improve`/`setup`/`wizard`, subject-vs-model altitude, target file,
  target/factor scope that composes with either altitude and resolves bare names
  against the grounded model, and rigor — where a bare `/quality` runs a
  read-only `wizard` mode that inspects model state and suggests concrete next
  invocations);
  **Driving the CLI** (`init`/`lint`/`spec` for every mechanical step, introspect
  rather than hard-code the surface); an **Evaluation workflow** that wraps the
  format spec's five Evaluation phases with mechanical read → lint → ground →
  evaluate → report steps; **Grounding judgment** (rate against declared
  criteria, evidence per rating, *not assessed* over guessing, inferred weighted
  roll-up); **Rigor levels** (`quick`/`standard`/`deep` coverage); and
  **Reporting** (the Evaluation Report, scoped, human or `--json`), plus an
  illustrative (non-normative) commented-examples block sketching invocation
  patterns in the style of the shadcn/improve README. Recording
  assessment results through the CLI and bundled `resources/` remain deferred in step
  with the CLI's deferred record/gate surface.
- **Revision**: Firmed up [`SPECIFICATION.md`](../SPECIFICATION.md) to fully
  ground the [`lint`](cli/lint.md) rules — made the per-property YAML shapes
  normative, defined a null/empty *required* value as absent, and stated the
  root-only `title`/`ratingScale` constraint as an explicit MUST NOT — then
  synced `lint`: broadened `invalid-frontmatter` to structural-shape conformance,
  added an "empty is absent" note, renamed `unresolved-factor` to `unknown-factor`
  (verb-consistent with `unknown-rating-key`), added `RECOMMENDED` to the RFC 2119
  keywords, corrected the rule `description` guidance and the
  `empty-model`/`empty-factor`/`empty-target` entries, and recorded the body
  heading under "Not checked".
- **Creation**: Added the [`skills/`](skills/index.md) subfolder with a stub
  [`/quality`](skills/quality-skill/quality-skill.md) functional spec — the companion
  evaluation skill that carries judgment against a `QUALITY.md` and records
  results through the deterministic CLI. The concrete workflow, inputs, and
  recorded-assessment shape remain deferred (in step with the CLI's deferred
  record/log surface).
- **Revision**: Fleshed out the [`lint`](cli/lint.md) sub-spec with a "Rule
  scope" section (four inclusion criteria — conformance-grounded, deterministic,
  format-not-goodness, self-contained — plus the error/warning/info severity
  mapping) and a "Rule naming" section (kebab-case noun-phrase ids in two shapes,
  using QUALITY.md vocabulary, treated as a stable public identifier) and a
  "Rule authoring" section with criteria for the static rule `description`
  (generic, present-tense, severity-signalling) and the per-instance finding
  `message` (names the element and location, contrasts found-vs-required,
  deterministic, tone-matched, actionable, one-per-instance), with the principle
  that description and message must not duplicate each other.
- **Revision**: Added the concrete "Rules" section to [`lint`](cli/lint.md) —
  eleven `error` rules (one per format-spec MUST: `invalid-frontmatter`,
  `missing-rating-scale`, `too-few-levels`, `missing-level-name`,
  `duplicate-level`, `missing-criterion`, `empty-model`, `misplaced-root-key`,
  `invalid-assessment`, `unknown-factor`, `unknown-rating-key`) and six
  `warning` rules (one per mechanically checkable SHOULD/RECOMMENDED:
  `missing-title`, `missing-level-description`, `missing-factor-description`,
  `empty-factor`, `empty-target`, `empty-property`), each citing the clause it
  enforces, plus a "Not checked" note for rating-level order (not mechanically
  verifiable) and the as-yet-unused `info` severity. The finding schema and
  lint-specific flags remain deferred.
- **Revision**: Replaced the placeholder [`init`](cli/init.md) command sub-spec
  with the durable requirements for scaffold contents, output target handling,
  stdout piping, overwrite protection, `--force`, and stderr reporting.
- **Revision**: Added a "Technical requirements" section to the
  [CLI spec](cli.md) requiring that every functional requirement be satisfiable
  through the idiomatic capabilities of the chosen stack (Go + Cobra + Charm
  Fang + Lip Gloss), rather than working against the grain of those libraries.
- **Revision**: Added a "Conventions" section to the [CLI spec](cli.md)
  establishing `--json` as the spelling for machine-readable output wherever a
  command offers one (not a requirement that every command do so), with criteria
  for when a command should offer `--json` and worked examples across the current
  commands. Updated [`lint`](cli/lint.md) to reference `--json` in place of its
  earlier `--format json`.
- **Convention**: Added a "Suggested next actions" convention to the
  [CLI spec](cli.md): commands may close with advisory, deterministic next-action
  suggestions — a stderr footer in human output, an in-band `nextActions` array
  under `--json` — that never affect behavior or the exit code.

## 2026-06-16

- **Convention**: Added a bundle-root [`schema.md`](schema.md) (`type: Schema`)
  registering the bundle's concept types (`Schema`, `Functional Specification`)
  in frontmatter, and listed it from the [index](index.md). Retyped the CLI
  spec and command sub-specs from `Specification` / `Command Specification` to a
  single `Functional Specification` type.

- **Initialization**: Created the `specs/` OKF bundle with a root
  [CLI spec](cli.md) capturing the high-level CLI requirements (design
  properties, global flags, output formats, exit codes, agent accessibility).
- **Creation**: Added placeholder command sub-specs for
  [`init`](cli/init.md), [`lint`](cli/lint.md), and [`spec`](cli/spec.md), plus
  the [`cli/` index](cli/index.md).
- **Revision**: Reduced the [CLI spec](cli.md) to a stub. No requirements had
  actually been requested, so the speculative design properties, global flags,
  output formats, exit codes, and agent-accessibility requirements were stripped
  out, leaving scope, a command list, and a "to be specified" outline — matching
  the placeholder command sub-specs.
