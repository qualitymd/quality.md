# /quality Skill Workflows Update Log

## 2026-06-29

- **Revision**: Updated the [`evaluate`](evaluate.md) workflow spec for
  [0163 - Report Artifact IDs](../../../../changes/archive/0163-report-artifact-ids.md).
  Evaluate now lets `qualitymd evaluation data set` assign public
  recommendation and ranked finding artifact IDs before report generation.

## 2026-06-27

- **Revision**: Updated the setup, evaluate, review, improve, and update
  workflow specs for
  [0145 - Scannable Skill Output](../../../../changes/archive/0145-scannable-skill-output.md).
  Workflow contracts now require labeled opening, review, and closeout templates
  for scan-friendly runtime output.

- **Revision**: Updated workflow specs for
  [0146 - Changelog Directory](../../../../changes/archive/0146-changelog-directory.md).
  Workflow mutation boundaries now name the quality changelog, and the shared
  feedback-log contract keeps workflow logs flat under `.quality/logs/`.

- **Addition**: Added the [`review`](review.md) and [`improve`](improve.md)
  workflow specs for
  [0143 - Public Review and Improve Workflows](../../../../changes/archive/0143-public-review-improve-workflows.md).
  The new specs define focus-routed public workflow stubs, review's read-only
  boundary, and improve's focus plus mutation-surface confirmation.

- **Revision**: Updated the [`evaluate`](evaluate.md) workflow spec for
  [0142 - Requirement Findings Only](../../../../changes/archive/0142-requirement-findings-only.md).
  Evaluate now records findings only on Requirement assessments, requires rated
  Requirements to have findings and rating drivers, and requires Factor/Area
  roll-up ratings to cite lower-level outputs through rating drivers instead of
  authoring new findings.

- **Revision**: Updated the [`evaluate`](evaluate.md) workflow spec for
  [0136 - Candidate Actions Payload](../../../../changes/archive/0136-candidate-actions-payload.md).
  Evaluate now writes Requirement Finding `candidateActions` with local IDs and
  keeps candidate actions out of Area Findings and recommendation-like
  presentation.

## 2026-06-26

- **Revision**: Updated the [`evaluate`](evaluate.md) workflow spec for
  [0132 - Remove info finding severity](../../../../changes/archive/0132-remove-info-finding-severity.md).
  Area Finding severity is now limited to `critical`, `high`, `medium`, and
  `low`, with informational observations represented as `type: note`.

- **Revision**: Updated the [`evaluate`](evaluate.md) workflow spec for
  [0131 - Area findings in evaluation reports](../../../../changes/archive/0131-area-findings.md).
  Evaluate now records material Area Findings on `AreaAnalysisResult.findings`
  during Area analysis, before report generation, and keeps recommendations,
  impact, priority, effort, benefit, ROI, and global rankings out of that field.

- **Revision**: Updated the [`evaluate`](evaluate.md) workflow spec for
  [0129 - Evaluation orchestration overhaul](../../../../changes/archive/0129-evaluation-orchestration-overhaul.md).
  The evaluate feedback-log sub-spec no longer permits a `rigor` frontmatter
  field, matching the single-workflow evaluation contract.

- **Revision**: Updated the [`setup`](setup.md) workflow spec for
  [0128 - Agent-mediated skill alignment](../../../../changes/archive/0128-agent-mediated-skill-alignment.md).
  Setup now requires the run frame as the first element in the opening block,
  clarifies that `QUALITY.md` changes wait for review while a local workflow
  feedback log may be created after the preview, and keeps feedback-log timing
  aligned with runtime guidance.

- **Revision**: Updated the [`setup`](setup.md) and [`update`](update.md)
  workflow specs for
  [0123 - Render interactions through native affordances](../../../../changes/archive/0123-native-interaction-affordances.md).
  `setup` defers presentation to the shared progressive-enhancement contract, tags
  the lifecycle/risk-tolerance/Rating-Scale questions as single-select intents,
  and locks the human context checkpoint as free text; `update` renders its
  apply-plan confirmation through a harness authorization gate when present rather
  than stacking a second text gate.

- **Revision**: Updated the [`setup`](setup.md) workflow spec for early-alpha
  compatibility cleanup.
  Setup now must report old `harnessability` factors as stale legacy naming to
  fix, not as current Agent Harnessability coverage.

- **Revision**: Updated the [`setup`](setup.md) and [`evaluate`](evaluate.md)
  workflow specs for
  [0121 - Scannable interaction hierarchy](../../../../changes/0121-scannable-interaction-hierarchy.md).
  Setup discovery blocks now require the rationale and recommendation to precede
  the answer line and per-option explanation to be capped. Evaluate now requires
  re-emitted factual progress beats before run creation and before the
  per-requirement loop, not only in the opening frame.

- **Revision**: Updated the [`evaluate`](evaluate.md) workflow spec for
  [0116 - Drop the "Evaluation v2" naming](../../../../changes/archive/0116-drop-evaluation-v2-naming.md).
  The workflow spec now links to the renamed Evaluation bundle and uses plain
  "Evaluation" for the active protocol.

- **Revision**: Updated the [`evaluate`](evaluate.md) workflow spec for
  [0114 - Run frame as first output](../../../../changes/archive/0114-run-frame-first-output.md).
  Its Required flow now requires the run frame as the first output before tool
  inspection, with a provisional scope value allowed.

- **Revision**: Unified workflow vocabulary in the [`setup`](setup.md),
  [`evaluate`](evaluate.md), [`update`](update.md), and [index](index.md) workflow
  specs for
  [0110 - Run frame title and workflow vocabulary](../../../../changes/0110-run-frame-and-workflow-vocabulary.md).
  Each spec now calls the dispatched procedure a "workflow" rather than a "mode",
  drops the `mode` frontmatter tag, and rewords "dispatched as a mode" framing.

- **Revision**: Updated the [`setup`](setup.md) and [`update`](update.md)
  workflow specs for
  [0106 - Binary confirmation UX](../../../../changes/archive/0106-binary-confirmation-ux.md).
  True binary mutation gates now require visible `y`/`n` answer paths, while
  numbered responses remain for multi-option discovery and routing prompts.

- **Revision**: Updated [`setup`](setup.md), [`evaluate`](evaluate.md), and
  [`update`](update.md) workflow specs for
  [0101 - Quality skill UX action clarity](../../../../changes/archive/0101-quality-skill-ux-action-clarity.md).
  Setup now requires explicit answer paths for open-ended discovery, a
  correction-first human context checkpoint, and a decision brief before writing
  `QUALITY.md`; evaluate now requires numbered ambiguity prompts and stop-response
  answer paths; update now emits its run frame before tool inspection and reports
  concise version-inspection progress before the mutation gate.

- **Revision**: Updated the [`setup`](setup.md) workflow spec for
  [0099 - Closed-choice setup UX](../../../../changes/archive/0099-closed-choice-setup-ux.md).
  Setup discovery closed-choice prompts now use numbered options with the
  recommended answer first and `1` as the shortest confirmation; risk discovery
  presents cost labels while mapping them to the existing risk-tolerance meaning.

- **Revision**: Updated the [`setup`](setup.md) workflow spec and
  [setup feedback log](setup/feedback-log.md) spec for
  [0096 - Setup intro preview](../../../../changes/archive/0096-setup-intro-preview.md).
  Setup now opens with short educational orientation, presents a project-specific
  setup preview after the read-only context scan and before discovery, and allows
  setup feedback-log creation after that preview when the run continues.

## 2026-06-25

- **Revision**: Updated the [`evaluate`](evaluate.md) workflow spec and
  [evaluate feedback log](evaluate/feedback-log.md) spec for
  [0095 - Evaluate feedback log outcomes](../../../../changes/archive/0095-evaluate-feedback-log-outcomes.md).
  Evaluate feedback logs remain under `.quality/logs/`, and their `outcome`
  field now records workflow terminal states such as `completed-reportable`,
  `stopped-model`, or `interrupted`, not report, rating, or recommendation state.

- **Revision**: Updated the [`setup`](setup.md) workflow spec and
  [setup feedback log](setup/feedback-log.md) spec for
  [0092 - Setup workflow scope trim](../../../../changes/archive/0092-setup-workflow-scope-trim.md).
  Setup no longer asks about future recommendation handling, review cadence,
  recurring review, or automation posture, and its completion contract now reports
  lint validation plus important model gaps. Setup feedback-log `outcome` values
  now describe workflow results instead of maturity or evaluation-readiness.

## 2026-06-24

- **Revision**: Updated the [`setup`](setup.md) workflow spec for
  [0091 - Agent-harness holistic definition](../../../../changes/archive/0091-agent-harness-holistic-definition.md).
  Setup must now check for project-owned runtime harness machinery, scope
  generated agent-harness Areas to checked-in steering and owned-control artifacts,
  distinguish that Area from the broader Agent Harnessability factor in comments,
  descriptions, and recap, and avoid dropping owned controls silently.

- **Revision**: Updated the [`setup`](setup.md) workflow spec for
  [0090 - Skill-content OKF authoring split](../../../../changes/archive/0090-skill-content-okf-authoring-split.md).
  Setup now points at the current seven-sub-factor Agent Harnessability shape,
  including `continuity`, matching the routed authoring guide family.
- **Revision**: Updated the [`setup`](setup.md) workflow for
  [0085 - Agent Harnessability naming](../../../../changes/archive/0085-agent-harnessability-naming.md).
  Setup now proposes `agent-harnessability` / Agent Harnessability by default for
  agent-collaborated composite roots, uses the accountability-preserving
  definition, and treats legacy `harnessability` with the expected six sub-factors
  as semantic coverage while recommending the new name during model-authoring work.
- **Revision**: Updated the [`setup`](setup.md), [`evaluate`](evaluate.md), and
  [`update`](update.md) workflow specs for
  [0084 - Agent-mediated UX conformance](../../../../changes/archive/0084-agent-mediated-ux-conformance.md).
  Runtime workflow output now follows the shared agent-mediated UX contract:
  status-first presentation, visually emphasized primary questions and calls to
  action, scannable labels, progress at useful phase boundaries, decision briefs
  with explicit mutation boundaries, and closeouts with clear next actions.
- **Revision**: Updated the [`setup`](setup.md) workflow for
  [0081 - Harnessability factor](../../../../changes/archive/0081-harnessability-factor.md).
  Setup now proposes harnessability by default as a model-wide umbrella factor for
  agent-collaborated composite roots, decomposed into the six authoring-guide
  sub-factors, and treats a thin or absent harness as rating evidence rather than
  an omission reason. The maturity bar now caps an agent-collaborated composite
  root below `evaluation-ready` when harnessability coverage is missing without a
  clear not-germane boundary.
- **Revision**: Updated the [`setup`](setup.md) workflow for
  [0080 - Model constituents by default](../../../../changes/0080-model-constituents-by-default.md).
  Setup's maturity classification now bars `evaluation-ready` for a composite root
  that leaves a germane, non-disqualified constituent unmodeled or recorded only as
  a deferral; a bare deferral or Scope note does not satisfy constituent coverage.
- **Revision**: Updated the [`setup`](setup.md) workflow for
  [0075 - Rating title emoji defaults](../../../../changes/0075-rating-title-emoji-defaults.md).
  Setup's recommended standard Rating Scale now keeps stable IDs plain while
  using `🟢 Outstanding`, `🔵 Target`, `🟡 Minimum`, and `🔴 Unacceptable` as
  default display titles, with emoji treated as a human scanning aid rather than
  rating semantics.

## 2026-06-23

- **Creation + Revision**: Added the
  [evaluate feedback log](evaluate/feedback-log.md) sub-spec for
  [0073 - Evaluation feedback log](../../../../changes/archive/0073-evaluation-feedback-log.md)
  and listed the new [evaluate/](evaluate/index.md) subfolder. Updated
  [`evaluate`](evaluate.md) so the workflow writes
  `.quality/logs/<timestamp>-evaluate-feedback-log.md`, and refactored
  [setup feedback log](setup/feedback-log.md) to adopt the shared
  [workflow feedback log](../workflow-feedback-log.md) contract.

- **Revision**: Updated the [`setup`](setup.md) workflow for
  [0072 - Setup context checkpoint](../../../../changes/archive/0072-setup-context-checkpoint.md).
  Human-context discovery now uses a correction-oriented checkpoint covering
  primary users/outcomes, maintainers/collaborators, other stakeholders, and
  missing or not-agent-accessible context instead of four separate open-ended
  prompts. The checkpoint keeps provenance visible and records omitted
  low-confidence or not-visible material context as Unknown, open questions, or
  low-confidence inference rather than confirmed fact.

- **Revision**: Updated the [`setup`](setup.md) workflow for
  [0071 - Setup open-ended review gate](../../../../changes/archive/0071-setup-open-ended-review-gate.md).
  The final review recap now uses friendly, open-ended wording that preserves the
  `"looks good"` fast path while inviting priorities, worries, wording, edge
  cases, repo-invisible context, and other useful last-call input before
  authoring.

- **Revision**: Updated the [`setup`](setup.md) workflow for
  [0070 - Setup missing-context provenance](../../../../changes/archive/0070-setup-missing-context-provenance.md).
  Missing-context discovery now treats material context as agent-accessible only
  when it is available through repository/tool/source evidence or explicit
  setup-provided context, prohibits generated choices that assume low/no-evidence
  project facts are sufficiently understood, and requires setup-provided context
  to retain visible provenance in the authored `QUALITY.md` body.

- **Revision**: Updated the [`setup`](setup.md) workflow for
  [0069 - Setup review gate and discovery trim](../../../../changes/archive/0069-setup-review-gate-and-pedagogy-trim.md).
  Setup discovery now removes the modeling-rigor and review-posture questions,
  adds a Rating Scale confirmation question, carries purpose/context teaching
  copy without per-question how-to-change-later guidance, and treats the final
  recap as a review gate that waits for an explicit user response before
  authoring.

- **Revision**: Updated the [`setup`](setup.md) workflow and
  [workflow feedback log](setup/feedback-log.md) sub-spec for
  [0068 - Always-on setup feedback log](../../../../changes/archive/0068-always-on-setup-feedback-log.md).
  Setup now creates the current run's feedback log during preflight after CLI
  support and the run frame, updates that file for material workflow-experience
  events, and finalizes it at close. The feedback-log sub-spec now defines the
  required frontmatter metadata, lifecycle status values, `Timeline` section,
  current-run update rule, and no-notable-content close behavior.

- **Revision**: Repositioned the [`setup`](setup.md) discovery contract as
  teaching-first for
  [0067 - Setup discovery pedagogy](../../../../changes/archive/0067-setup-discovery-pedagogy.md).
  Required authored per-question pedagogy (purpose + how-to-change-later) and the
  teaching-over-round-trips framing; changed the prompt-form contract to ask every
  one of the ten questions every run (removing the accept-all-defaults-and-skip
  escape, keeping a per-question fast confirm and show-all-at-once); changed the
  confidence vocabulary `MUST` from `strongly inferred`/`weakly inferred`/`assumed`
  to `Low`/`Med`/`High` with the evidence note retained (no-evidence → `Low` plus
  an explicit note); and added the final-review-recap stage between discovery and
  authoring. Promoted the rationale into per-requirement annotations. Append-only
  entries below keep their old confidence-vocabulary references frozen.

- **Creation + Revision**: Added the workflow feedback-log artifact contract for
  [0066 - Setup feedback log](../../../../changes/archive/0066-setup-feedback-log.md).
  Created the [`setup/`](setup/index.md) sub-folder with the
  [workflow feedback log](setup/feedback-log.md) sub-spec (location/naming under
  `.quality/logs/`, environment header, body schema, redaction, no-transmission
  posture), making [`setup`](setup.md) a parent concept. Amended the
  [`setup`](setup.md) mutation surface to permit the `.quality/logs/` write while
  keeping every other setup prohibition in force, and added the close-step
  feedback-log authoring. Updated the [index](index.md) to list the sub-folder.

- **Rename + Revision**: Renamed this folder from `modes/` to `workflows/` and
  updated its [index](index.md) for
  [0065 - Setup discovery and close refinements](../../../../changes/0065-setup-discovery-and-close-refinements.md).
  Replaced the [`setup`](setup.md) prompt-form contract with agent-agnostic
  presentation — present all ten discovery questions, iterate one at a time when
  there is no structured question affordance, page through a structured tool when
  there is, keep open-ended questions free text, and honor accept-all/show-all
  escapes on request. Added the read-the-scaffold-before-authoring requirement,
  and reframed the close contract to classify model maturity (`starter`,
  `immature`, `evaluation-ready`) as distinct from the CLI's lifecycle
  `readiness`. Updated runtime path references in [`evaluate`](evaluate.md) and
  [`update`](update.md) to the `workflows/` path. Historical entries below keep
  their `modes/` references frozen.

- **Revision**: Updated the [`setup`](setup.md) workflow spec for
  [0064 - Structured setup workflow](../../../../changes/0064-structured-setup-workflow.md).
  Setup now has an explicit workflow structure, setup brief, concrete
  discovery-question set, prompt-form contract, model synthesis expectations,
  and completion output while keeping the existing `modes/` path as dispatch
  vocabulary.

- **Revision**: Updated the [`setup`](setup.md) mode spec for
  [0063 - Contextual setup flow](../../../../changes/0063-contextual-setup-flow.md).
  Setup now analyzes repository context, asks compact discovery questions with
  confidence-labeled defaults, writes only `QUALITY.md`, validates with lint,
  inspects readiness with the Top 10 checklist, and offers next-step choices
  without writing the quality changelog or configuring integrations.

- **Mode removal**: Removed the `wizard` mode spec from this folder for
  [0062 - Remove wizard mode](../../../../changes/0062-remove-wizard-mode.md).
  Bare and ambiguous `/quality` requests are now governed by the parent
  [`/quality` skill](../quality-skill.md) spec as read-only orientation, not as a
  public mode.

## 2026-06-22

- **Revision**: Updated the [`evaluate`](evaluate.md) mode spec for
  [0056 - Prospective evaluation plan artifacts](../../../../changes/archive/0056-prospective-evaluation-plan-artifacts.md)
  so `design.md` and the initial `plan.md` are authored before assessment
  evidence collection or record writes, later plan changes are explicit
  amendments, and `debug-log.md` remains process-only.

- **Mode removal**: Removed the `improve` mode spec from this folder. Applying or
  handing off evaluation recommendations is now governed by the non-mode
  [`recommendation follow-up`](../recommendation-follow-up.md) spec.

- **Creation**: Originally added behavioral component specs for the `/quality` runtime modes:
  [`setup`](setup.md), `wizard`, [`evaluate`](evaluate.md), `improve`, and
  [`update`](update.md). The parent
  [`/quality` skill](../quality-skill.md) spec keeps shared contracts and links
  to these mode-specific contracts.
