# /quality Skill Workflows Update Log

## 2026-06-23

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
