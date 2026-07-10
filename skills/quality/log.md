# /quality runtime skill update log

## 2026-07-09

- **Revision**: Updated [`SKILL.md`](SKILL.md) for 0192 - Deterministic
  evaluation runner. The CLI now owns evaluation orchestration and judgment
  dispatch through `qualitymd evaluation run`; evaluate writes evaluation
  records only through the runner, the payload-batch authoring guidance is
  removed for new runs, workspace config gains `evaluation.evaluator`,
  `evaluation.executionStrategy`, and the `evaluators` profile map, and the
  evaluate feedback log stays skill-owned.

## 2026-07-02

- **Revision**: Updated [`SKILL.md`](SKILL.md) prerequisites for the
  document-only `0.10 (Draft)` format specification. Runtime guidance now points
  agents to Model semantics for source resolution, Requirement scope, Factor
  connection, and Rating Scale meaning.

- **Release prep**: Bumped the runtime skill metadata to `0.27.0` with
  `qualitymd >=0.27.0 <0.28.0` compatibility for the `v0.27.0` release.

## 2026-06-30

- **Revision**: Updated runtime skill headings for 0189 - Heading sentence case.
  Active `/quality` skill, guide, resource, and workflow headings now follow
  sentence case while preserving formal QUALITY.md concept names.

- **Revision**: Updated [`SKILL.md`](SKILL.md) for 0185 - Concern Finding
  Severity.
  Runtime evaluation guidance now tells agents to write `severity` only for
  `gap` and `risk` Findings, omitting it from `strength` and `note` Findings.

## 2026-06-29

- **Revision**: Updated [`SKILL.md`](SKILL.md),
  [`workflows/evaluate.md`](workflows/evaluate.md), and
  [`resources/cli-workflow-conventions.md`](resources/cli-workflow-conventions.md)
  for 0181 - Evaluation Identity Manifest. Runtime evaluation guidance now names
  the CLI-owned Evaluation manifest and its `data/evaluation-manifest.json`
  artifact.

- **Revision**: Updated [`SKILL.md`](SKILL.md) and
  [`workflows/evaluate.md`](workflows/evaluate.md) for 0180 - Finding Taxonomy
  and Report Details. Runtime evaluation guidance now treats ambiguous
  current-state shortfalls as `gap` findings and routes missing evidence that
  prevents rating to not-assessed/not-rated status, `unknowns`, or
  `missingEvidence`.

- **Revision**: Updated [`SKILL.md`](SKILL.md) and
  [`workflows/evaluate.md`](workflows/evaluate.md) for 0176 - Recommendation IDs
  and Numbers. Runtime evaluation guidance now uses `qrec_...` recommendation
  IDs in structured ranking and coverage while reserving recommendation numbers
  for ranked report order.

- **Revision**: Updated [`SKILL.md`](SKILL.md) for 0173 - Evaluation Enum
  Catalogs. Runtime evaluation guidance now tells the skill to write canonical
  fixed enum values rather than generated-report labels or markers.

- **Revision**: Updated [`SKILL.md`](SKILL.md) and
  [`resources/cli-workflow-conventions.md`](resources/cli-workflow-conventions.md)
  for 0172 - Workspace Status Contract. Runtime guidance now treats
  `qualitymd status` as workspace status while preserving project wording for
  modeled value and setup intent.

- **Release prep**: Bumped the runtime skill metadata to `0.26.1` with
  `qualitymd >=0.26.0 <0.27.0` compatibility for the `v0.26.1` release.

- **Release prep**: Bumped the runtime skill metadata to `0.26.0` with
  `qualitymd >=0.26.0 <0.27.0` compatibility for the `v0.26.0` release.

- **Revision**: Updated [`SKILL.md`](SKILL.md) and
  [`workflows/evaluate.md`](workflows/evaluate.md) for 0165 - Run IDs and
  Artifact Numbering. Runtime evaluation guidance now lets the CLI assign
  recommendation numbers, references recommendations by number in ranking and
  coverage, and keeps finding ranking entries free of synthetic artifact IDs.

- **Revision**: Updated [`workflows/setup.md`](workflows/setup.md) for 0164 -
  Agent Instruction Init Pointer. Runtime setup guidance now opts out of
  init-time agent instruction edits and handles existing `QUALITY.md` maturity
  explicitly.

## 2026-06-27

- **Revision**: Updated [`workflows/evaluate.md`](workflows/evaluate.md) for
  0151 - Evaluation Report CTA.
  Runtime evaluation closeouts now point first to `<run>/report.md` as the
  decision-ready evaluation result and describe `<run>/recommendations.md` as
  the action-planning report, instead of presenting a generic report list.

- **Revision**: Updated [`SKILL.md`](SKILL.md) and
  [`workflows/evaluate.md`](workflows/evaluate.md) for 0150 - Evaluation
  Advice. Runtime evaluation guidance now requires Advice after roll-up, ranks
  Findings, generates domain-agnostic Recommendations, accounts for every
  Finding, ranks Recommendations, and reports `recommendations.md` alongside
  `report.md`.

- **Revision**: Updated [`SKILL.md`](SKILL.md),
  [`workflows/evaluate.md`](workflows/evaluate.md), and
  [`resources/cli-workflow-conventions.md`](resources/cli-workflow-conventions.md)
  for 0149 - Scope-driven evaluation runs. Runtime evaluation guidance now
  resolves natural scope to canonical `area:`/`factor:` IDs before create, uses
  `--area`/`--factor`, relies on CLI-owned `RunManifest` scope, and treats
  `report.md` as the scoped Area report.

- **Revision**: Updated [`SKILL.md`](SKILL.md) for 0145 - Scannable Skill
  Output.
  Runtime interaction guidance now requires labeled blocks, bullets, or numbered
  lists for multi-fact outputs, and direct model-authoring checkpoints use
  `Planned edit`, `Why`, `Approach`, `Boundary`, `Log`, and `Answer` labels.

- **Revision**: Updated [`SKILL.md`](SKILL.md) and workflow guidance for
  0146 - Changelog Directory.
  Runtime guidance now writes meaningful model-change history to the quality
  changelog under `.quality/changelog/`, uses timestamped changelog entry names,
  and keeps `.quality/logs/` as the flat workflow-log directory for feedback
  logs and future process logs.

- **Revision**: Updated [`SKILL.md`](SKILL.md) for 0144 - Pointed Review Gates.
  Direct model authoring checkpoints now state inferred purpose and ask users to
  react to the consequential assumption instead of defaulting to a generic
  adjustment prompt.

- **Revision**: Updated [`SKILL.md`](SKILL.md) for 0143 - Public Review and
  Improve Workflows.
  Runtime dispatch now treats `review` and `improve` as public workflows, uses
  focus for review/improve routing, and delegates model improvement and
  compatible recommendation artifacts to existing safe routes.

- **Revision**: Updated [`SKILL.md`](SKILL.md) and
  [`workflows/evaluate.md`](workflows/evaluate.md) for 0142 - Requirement
  Findings Only.
  Runtime evaluation guidance now makes Requirement Findings the only finding
  layer, requires rated Requirements to be finding-backed and driver-backed, and
  uses rating drivers/rationale for Factor and Area roll-ups instead of Area
  Findings.

- **Revision**: Updated authoring guides for 0141 - Area-local Factor
  References.
  Runtime guidance now keeps Requirement Factor references local to the
  Requirement's declaring Area and directs cross-Area root judgment toward root
  Requirements rather than descendant opt-in Factor references.

- **Revision**: Updated [`SKILL.md`](SKILL.md) for 0140 - Casual Review Gate
  Wording.
  Direct model authoring checkpoints now use a plain "Here's what I'm planning
  to do" plan, name the value prop, invite feedback with "what should I adjust or
  watch out for", and keep `looks good` as the short approval path before
  mutation.

- **Revision**: Updated [`SKILL.md`](SKILL.md) for 0139 - Real Review Gates.
  Direct model authoring now acknowledges long model/guidance reads before they
  begin, phrases the intended-edit checkpoint as a review gate, waits for
  `looks good` or corrections before mutation, and avoids asking for feedback
  while proceeding in the same turn.

- **Revision**: Updated [`SKILL.md`](SKILL.md) and authoring guides for 0138 -
  Lightweight Authoring Checkpoint.
  Runtime guidance now routes direct existing `QUALITY.md` edit requests to
  direct model authoring, infers intent before asking follow-up, uses a
  conversational checkpoint that accepts `looks good` when the mutation is named,
  and writes the quality changelog only for meaningful model changes.

- **Revision**: Updated [`SKILL.md`](SKILL.md) and
  [`workflows/evaluate.md`](workflows/evaluate.md) for 0137 - Run Report
  Entrypoint.
  Runtime guidance now treats `report.md` as the run-level Evaluation report,
  names `headlineReportMd` separately, and uses `root-area.md` only when the
  root Area was evaluated.

- **Revision**: Updated [`SKILL.md`](SKILL.md) for 0134 - Model-relative
  workspace paths.
  Runtime workspace guidance now resolves config, `.quality/`, Evaluation
  history, quality changelogs, and workflow feedback logs relative to the selected
  `QUALITY.md`, with the Git root retained as the containment boundary.

- **Revision**: Updated [`SKILL.md`](SKILL.md) for 0136 - Candidate Actions
  Payload.
  Runtime Area Finding guidance now excludes `candidateActions` by name, keeping
  candidate actions Requirement Finding-local raw material for later Advice.

## 2026-06-26

- **Revision**: Updated [`SKILL.md`](SKILL.md) for 0132 - Remove info finding
  severity.
  Runtime Area Finding guidance now limits severity to `critical`, `high`,
  `medium`, and `low`, and routes informational observations to `type: note`
  rather than `severity: info`.

- **Revision**: Updated [`SKILL.md`](SKILL.md) for 0131 - Area findings in
  evaluation reports.
  Runtime evaluation guidance now tells Area analysis to synthesize
  `AreaAnalysisResult.findings` for material Area/Factor report observations and
  keeps recommendations, impact, priority, effort, benefit, ROI, actions, and
  global top-finding rankings out of Area Findings.

- **Revision**: Updated [`SKILL.md`](SKILL.md) for 0130 - Self-contained
  per-kind data schema.
  Artifact Contract guidance now treats `qualitymd evaluation data schema <kind>`
  as the source for required fields and allowed enum values, with
  `data example <kind>` as one concrete instance and `data set --dry-run` as
  authored-payload validation only.

- **Revision**: Updated [`SKILL.md`](SKILL.md) for 0129 - Evaluation
  orchestration overhaul.
  Runtime guidance now removes the evaluation rigor argument and run-frame field,
  makes exhaustive in-scope coverage the only evaluate workflow, defaults
  collection/QC to subagent fan-out when available, and requires the two-pronged
  QC phase before roll-up.

- **Revision**: Updated [`SKILL.md`](SKILL.md) for 0128 - Agent-mediated skill
  alignment.
  Runtime guidance now gives read-only orientation a status-first shape and
  requires recommendation follow-up to emit a frame before inspection, outcome
  selection, or mutation.

- **Revision**: Reshaped the [`SKILL.md`](SKILL.md) User Interaction Contract for
  0121 - Scannable interaction hierarchy.
  The decision-brief template now leads with the question, renders choices as a
  separated `[y]`/`[n]` block with the alternative folded in, and demotes a capped
  set of `label:` rationale lines; the contract prose now requires hierarchy by
  position rather than bold alone, surviving bold-stripping.

- **Revision**: Retitled the run-frame header label from `Quality` to
  `QUALITY.md`.
  The run-frame template header is now `**QUALITY.md · <workflow>**` across the
  dispatcher and the setup/evaluate/update workflows.

- **Revision**: Updated the Evaluation data contract guidance for 0115 -
  Type-safe, model-bound Evaluation v2 data.
  The dispatcher now points agents to `qualitymd evaluation data schema [<kind>]`
  as the authoritative payload-shape source, examples as populated samples, and
  `data set --dry-run` as authored-payload validation.

- **Revision**: Renamed live Evaluation wording for 0116 - Drop the "Evaluation
  v2" naming.
  Runtime skill prose now uses plain "Evaluation" for the active workflow and
  data surface.

- **Revision**: Tightened the run-frame instruction for 0114 - Run frame as
  first output.
  The dispatcher now requires the run frame as the workflow's first output before
  any tool call, forbids gating emission on a tool result, and allows a
  best-known or `resolving…` value for a field (such as a many-Area scope) that
  needs a tool to resolve.

- **Revision**: Retitled the run frame and unified workflow vocabulary for 0110 -
  Run frame title and workflow vocabulary.
  The run-frame template header is now `**Quality · <workflow>**` instead of
  `**/quality run**`, the `Mode:` field is dropped (the workflow name moves into
  the header), and Arguments/Workflow Dispatch wording now says "workflow" rather
  than "mode".

- **Revision**: Updated the root runtime interaction contract for 0106 - Binary
  confirmation UX.
  Runtime guidance now distinguishes non-binary closed choices, which keep `1` as
  the shortest accept path, from true binary mutation confirmations, which show
  `y`/`n`.

- **Revision**: Updated the root runtime interaction contract for 0101 - Quality
  skill UX action clarity.
  Runtime guidance now requires explicit shortest-answer paths for user
  interactions, code spans for concrete operational examples, and numbered
  ambiguity choices for scoped evaluation prompts.

## 2026-06-24

- **Restructure**: Started the runtime skill content as an OKF-shaped bundle with
  root `index.md`, `schema.md`, and `log.md`; added guide indexes/logs; and split
  authoring guidance into routed sub-guides.
