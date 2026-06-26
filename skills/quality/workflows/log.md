# /quality Runtime Workflows Update Log

## 2026-06-26

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
