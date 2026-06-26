# /quality Runtime Workflows Update Log

## 2026-06-26

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
