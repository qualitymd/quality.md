# /quality Skill Workflows Update Log

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
  without writing the quality log or configuring integrations.

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
