# /quality Runtime Guides Update Log

## 2026-06-27

- **Revision**: Updated [Authoring QUALITY.md](authoring.md) and
  [Authoring Quality-Log Changes](authoring/quality-log.md) for 0138 -
  Lightweight Authoring Checkpoint.
  Direct `QUALITY.md` edits now preserve intent, target, rationale, judgment
  effect, unresolved unknowns, and quality-log routing; confirmed direct
  model-authoring changes write the quality log only when they meaningfully alter
  model judgment.

## 2026-06-26

- **Revision**: Updated
  [Recommendation Follow-Up](recommendation-follow-up.md) for 0128 -
  Agent-mediated skill alignment.
  Recommendation follow-up now opens with a frame before recommendation
  inspection, outcome selection, local apply, issue creation, or quality-log
  writes.

- **Revision**: Updated [Agent Harnessability](authoring/agent-harnessability.md)
  and [Top 10 checks](top-10-quality-md-checks.md) for early-alpha compatibility
  cleanup.
  Runtime guidance now treats old `harnessability` factors as stale legacy naming
  that must be corrected, not as current Agent Harnessability coverage.

- **Revision**: Reshaped the gates and result block in
  [Recommendation Follow-Up](recommendation-follow-up.md) for 0121 - Scannable
  interaction hierarchy.
  The apply and issue-creation gates now lead with the question and use a
  separated `[y]`/`[n]` choice block with the alternative folded into the stop
  choice; the recommendation result block leads with a primary outcome line and
  de-stacks its labels.

- **Revision**: Updated
  [Recommendation Follow-Up](recommendation-follow-up.md) for 0110 - Run frame
  title and workflow vocabulary.
  The guide now describes recommendation follow-up as "not a public `/quality`
  workflow" instead of "not a `/quality` mode".

- **Revision**: Updated
  [Recommendation Follow-Up](recommendation-follow-up.md) for 0106 - Binary
  confirmation UX.
  Runtime guidance now shows `y`/`n` answer paths for true binary local-apply and
  external issue-creation gates.

- **Revision**: Updated [Recommendation Follow-Up](recommendation-follow-up.md)
  for 0101 - Quality skill UX action clarity.
  The guide now gives apply-vs-handoff outcome selection a numbered answer path,
  requires a decision brief before external issue creation, and includes `Next`
  plus boundary-sensitive `Not done` in recommendation result closeouts.

## 2026-06-25

- **Revision**: Updated [Authoring QUALITY.md](authoring.md), the
  [authoring index](authoring/index.md), and
  [Requirements](authoring/requirements.md) for 0093 - Named Requirement
  identity.
  Runtime guidance now teaches stable Requirement names as durable identity,
  natural-language Requirement titles for report display, and named Requirement
  YAML examples.

- **Revision**: Updated [Getting started](getting-started.md) and
  [Top 10 checks](top-10-quality-md-checks.md) guides for 0092 - Setup workflow
  scope trim.
  Getting-started guidance now focuses on first-model gaps instead of
  starter/immature labels or follow-on handoff/review setup. Top 10 now reports
  lifecycle state and model-usefulness findings instead of maturity or
  evaluation-readiness classifications.

## 2026-06-24

- **Revision**: Updated the
  [agent-harness Area](authoring/agent-harness.md),
  [Agent Harnessability](authoring/agent-harnessability.md),
  [model-structure](authoring/model-structure.md), and
  [Top 10 checks](top-10-quality-md-checks.md) guides for 0091 - Agent-harness
  holistic definition.
  The runtime guidance now defines the agent harness as the whole engineered
  system around the model, scopes the agent-harness Area to checked-in steering
  and owned-control artifacts, adds the mixed-artifact decision rule, expands
  requirement shapes across feedforward, feedback, and owned controls, and flags
  instructions-only or unmodeled runtime-harness gaps.

- **Restructure**: Split [Authoring QUALITY.md](authoring.md) into a compact
  entry guide plus routed sub-guides under [authoring/](authoring/index.md), and
  added this guide index/log for runtime OKF navigation.
