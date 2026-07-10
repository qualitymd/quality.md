# /quality runtime workflows update log

## 2026-07-10

- **Revision**: Updated the [evaluate](evaluate.md) workflow for 0195 -
  Evaluation concurrency. Dry-run selection guidance now describes resolved
  concurrency instead of execution strategy.

## 2026-07-09

- **Revision**: Rewrote the [evaluate](evaluate.md) workflow for 0192 -
  Deterministic evaluation runner. Evaluate now wraps
  `qualitymd evaluation run`: it frames the run, verifies CLI support, resolves
  scope, lints, opens the feedback log, explains evaluator selection (with an
  optional `--dry-run --json` preview), invokes the runner, handles
  failure/resume, and summarizes the receipt and generated reports instead of
  orchestrating assessment, QC, roll-up, and payload writes itself.

## 2026-06-29

- **Revision**: Updated runtime workflow headings for 0189 - Heading Sentence
  Case. Active setup/evaluate/review/improve/update headings now use sentence
  case while preserving formal workflow and model concept names.

- **Revision**: Updated the [evaluate](evaluate.md) workflow for 0181 -
  Evaluation Identity Manifest. Evaluate now names
  `data/evaluation-manifest.json` as the CLI-owned manifest artifact for scoped
  Evaluation runs.

- **Revision**: Updated the [evaluate](evaluate.md) workflow for 0180 - Finding
  Taxonomy and Report Details. Evaluate now writes only `strength`, `gap`,
  `risk`, and `note` Finding types, uses `gap` for ambiguous current-state
  rating constraints, and records missing evidence through non-finding fields.

- **Revision**: Updated the [evaluate](evaluate.md) workflow for 0165 - Run IDs
  and Artifact Numbering. Evaluation Advice now writes recommendation payloads
  before ranking when numbers are assigned, references recommendations by
  number, and leaves finding ranking entries without artifact IDs.

- **Revision**: Updated the [setup](setup.md) workflow for 0164 - Agent
  Instruction Init Pointer. Setup now scaffolds missing model files with
  `qualitymd init --no-agent-instructions` and classifies existing model files
  as scaffold-only, partially authored, or mature before planning edits.

## 2026-06-27

- **Revision**: Updated setup, evaluate, review, improve, and update workflows
  for 0145 - Scannable Skill Output.
  Runtime workflow guidance now uses labeled opening, review, and closeout
  templates so status, artifacts, boundaries, verification, and next actions are
  visible in a five-second scan.

- **Revision**: Updated setup, evaluate, review, and improve workflow guidance
  for 0146 - Changelog Directory.
  Workflow mutation boundaries now name the quality changelog, while setup and
  evaluate feedback logs remain flat files under `.quality/logs/`.

- **Addition**: Added [review](review.md) and [improve](improve.md) runtime
  workflow stubs for 0143 - Public Review and Improve Workflows.
  Review is read-only and focus-routed; improve starts read-only, confirms focus
  and mutation surface, then delegates to existing safe routes or stops at the
  stub boundary.

- **Revision**: Updated the [evaluate](evaluate.md) workflow for 0142 -
  Requirement Findings Only.
  Evaluate now authors findings only during Requirement assessment, requires
  rated Requirements to have findings and rating drivers, and keeps Factor/Area
  roll-up analysis driver-backed without creating new findings.

- **Revision**: Updated the [evaluate](evaluate.md) workflow for 0136 -
  Candidate Actions Payload.
  Evaluate now writes finding-local `candidateActions` with `id`, `description`,
  and optional `rationale` when a local remediation lead is evident, and keeps
  them out of Area Findings and recommendation-like presentation.

## 2026-06-26

- **Revision**: Updated the [evaluate](evaluate.md) workflow for 0132 - Remove
  info finding severity.
  Evaluate now tells agents to use only `critical`, `high`, `medium`, and `low`
  for Area Finding severity and to use `type: note` for informational
  observations.

- **Revision**: Updated the [evaluate](evaluate.md) workflow for 0131 - Area
  findings in evaluation reports.
  Evaluate now records material Area Findings on `AreaAnalysisResult.findings`
  during Area analysis, before persistence and report generation, and keeps
  advice/ranking fields out of that data.

- **Revision**: Updated the [evaluate](evaluate.md) workflow for 0130 -
  Self-contained per-kind data schema.
  Before authoring a payload kind, Evaluate now reads
  `qualitymd evaluation data schema <kind>` for required fields and allowed enum
  values, then treats `qualitymd evaluation data example <kind>` as one concrete
  instance rather than a closed value-set source.

- **Revision**: Updated the [evaluate](evaluate.md) workflow for 0129 -
  Evaluation orchestration overhaul.
  Evaluate now removes the `Rigor:` run-frame field and feedback-log `rigor`
  field, assesses every in-scope Requirement against the full in-scope source,
  uses subagent fan-out for independent collection/QC when available, and runs
  verify plus completeness-sweep QC before roll-up.

- **Revision**: Updated the [setup](setup.md) workflow for 0128 -
  Agent-mediated skill alignment.
  Setup now renders the run frame as the first element in its opening block,
  clarifies that `QUALITY.md` changes wait for review while a local workflow
  feedback log may be created after the preview, and keeps feedback-log timing
  consistent with shared workflow conventions.

- **Revision**: Updated the [setup](setup.md) workflow for early-alpha
  compatibility cleanup.
  Setup now reports old `harnessability` factors as stale legacy naming to fix,
  not as current Agent Harnessability coverage.

- **Revision**: Reshaped interaction output for 0121 - Scannable interaction
  hierarchy.
  [setup](setup.md): write/update gates lead with the question and use a separated
  choice block; discovery inputs now place rationale and recommendation before the
  answer line, the rating-scale glosses are capped and offered on demand, and Q1/Q2
  accept `y`. [update](update.md): the apply-plan gate uses the separated-choice
  shape. [evaluate](evaluate.md): added re-emitted progress beats before run
  creation and before the per-requirement loop.

- **Revision**: Retitled the rendered run-frame headers in [setup](setup.md),
  [evaluate](evaluate.md), and [update](update.md) from `Quality` to `QUALITY.md`.
  Each frame header is now `**QUALITY.md · <workflow>**`.

- **Revision**: Updated the [evaluate](evaluate.md) procedure for 0115 -
  Type-safe, model-bound Evaluation v2 data.
  Evaluate now discovers payload shape with `qualitymd evaluation data schema
<kind>`, uses populated examples for samples, and treats `data set --dry-run`
  as validation of authored payloads.

- **Revision**: Reordered the [evaluate](evaluate.md) procedure for 0114 - Run
  frame as first output.
  The run frame is now emitted as the first output before workspace resolution or
  any other tool call, using the invocation-derived model path, with a provisional
  `Scope: resolving…` confirmed once the workspace and model are read.

- **Revision**: Retitled the rendered run frames in [setup](setup.md),
  [evaluate](evaluate.md), and [update](update.md) for 0110 - Run frame title and
  workflow vocabulary.
  Each frame header is now `**Quality · <workflow>**` and drops the `Mode:` line;
  the [evaluate](evaluate.md) and [update](update.md) H1 titles change from
  "… Mode" to "… Workflow".

- **Revision**: Updated [setup](setup.md) and [update](update.md) workflows for
  0106 - Binary confirmation UX.
  Runtime guidance now shows `y`/`n` answer paths for true binary mutation gates
  while preserving numbered responses for non-binary choices.

- **Revision**: Updated [setup](setup.md), [evaluate](evaluate.md), and
  [update](update.md) workflows for 0101 - Quality skill UX action clarity.
  Setup now gives open-ended discovery questions explicit answer paths, presents
  the human context checkpoint as a correction-first table, and uses a decision
  brief before writing `QUALITY.md`; evaluate stop and ambiguity prompts now
  include numbered response paths; update now starts with a run frame and reports
  version-inspection progress before confirmation.

- **Revision**: Updated the [setup](setup.md) workflow for 0096 - Setup intro
  preview.
  Setup now opens with short educational orientation, runs the first context pass
  as an explicit read-only scan, presents a project-specific setup preview before
  discovery questions, and creates the setup feedback log after that preview when
  the run continues.

## 2026-06-25

- **Revision**: Updated the [evaluate](evaluate.md) workflow for 0095 - Evaluate
  feedback log outcomes.
  Evaluate feedback logs remain under `.quality/logs/` as process-only workflow
  artifacts, and their `outcome` field now records workflow terminal states such
  as `completed-reportable`, `stopped-model`, or `interrupted`, not report or
  rating semantics.

- **Revision**: Updated the [setup](setup.md) workflow for 0092 - Setup workflow
  scope trim.
  Setup no longer asks about recommendation handling, handoff destination, review
  cadence, recurring review, or automation preferences. Its closeout now reports
  lint validation plus important model gaps, and setup feedback logs record
  workflow-result outcomes instead of maturity or evaluation-readiness labels.

## 2026-06-24

- **Revision**: Updated the [setup](setup.md) workflow for 0091 - Agent-harness
  holistic definition.
  Setup now actively checks for project-owned runtime harness machinery and scopes
  generated agent-harness Areas as checked-in steering and owned-control artifacts
  distinct from the broader Agent Harnessability factor.

- **Restructure**: Added this workflow index/log as part of the runtime skill
  OKF shape.
