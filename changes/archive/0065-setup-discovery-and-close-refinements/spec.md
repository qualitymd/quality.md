---
type: Functional Specification
title: Setup discovery and close refinements
description: Agent-agnostic discovery presentation, read-before-author, model-maturity vocabulary, and the modes/ to workflows/ rename for /quality setup.
tags: [skill, setup, ux]
timestamp: 2026-06-23T00:00:00Z
---

# Setup discovery and close refinements

This Change Case spec defines four deltas to the `/quality setup` workflow:
present the discovery questions in an agent-agnostic way, read the scaffolded
model before authoring it, disentangle the skill's model-maturity judgment from
the CLI's lifecycle readiness, and rename the skill's `modes/` folder to
`workflows/`. It does not change which questions are asked, their defaults, the
confidence vocabulary, the QUALITY.md format, or any `qualitymd` CLI behavior.

This document uses BCP 14 keywords only for testable conformance requirements.
The key words "MUST", "MUST NOT", "SHOULD", and "MAY" are to be interpreted as
described in [RFC 2119](../../docs/reference/rfc2119.md) and
[RFC 8174](../../docs/reference/rfc8174.md) when, and only when, they appear in
all capitals.

## Background / Motivation

A first end-to-end `/quality setup` run against an external monorepo completed
without blocking errors but exposed four avoidable frictions, detailed in the
[parent change case](../0065-setup-discovery-and-close-refinements.md):
the discovery prompt is implicitly shaped to one agent's question UI; authoring
after `qualitymd init` costs a wasted read-before-write round-trip; the close
step's "readiness" word collides with the CLI `status` lifecycle `readiness`;
and the `modes/` folder name predates the "workflow" terminology the setup
guidance now uses.

The discovery questions carry both critical setup context and a pedagogical
purpose — they teach the user the dimensions a quality model captures — so the
presentation must surface all of them regardless of the agent's interaction
surface, never silently dropping or merging one to fit a tool's caps.

## Scope

Covered:

- Agent-agnostic presentation of the existing ten discovery questions.
- One-question-at-a-time iteration when no structured question affordance exists.
- Paging all questions through a structured question tool within its caps.
- Early-exit escapes on user request.
- A read-before-author step after `qualitymd init`.
- Renaming the skill's model-maturity judgment off the word "readiness" and
  separating maturity labels from CLI lifecycle labels.
- A condensed close checklist so setup need not read the full top-10 guide every
  run.
- Renaming the skill `modes/` folder to `workflows/` in the runtime skill and
  the `specs/` mirror, with all live path references updated.

Deferred / non-goals:

- No QUALITY.md format change.
- No `qualitymd` CLI or Go code change. The CLI `status` lifecycle `readiness`
  field and values are unchanged.
- No change to which discovery questions are asked, their option sets, their
  recommended defaults, or the `strongly inferred`/`weakly inferred`/`assumed`
  confidence vocabulary.
- No new public skill workflow, and no change to setup's mutation boundary
  (`QUALITY.md` only).

## Requirements

### Discovery presentation

The setup workflow **MUST** present all ten discovery questions every run,
unless the user explicitly asks to accept all inferred defaults. It **MUST NOT**
drop, merge, or silently default away a question to fit an interaction surface's
limits.

The workflow **MUST** choose how to present the questions from the agent's own
interaction capabilities. The guidance **MUST NOT** assume or name a specific
agent's question tool.

When the agent has a structured question affordance with item or option limits,
the workflow **MUST** page all ten questions through it across as many rounds as
the limits require, and **MUST** keep open-ended questions (primary users,
maintainers and collaborators, other stakeholders, missing context) as free
text rather than forcing them into fixed options.

When the agent has no structured question affordance, the workflow **MUST**
iterate the questions one at a time. Each step **MUST** carry that question's
recommended default and confidence signal so the user can confirm or correct it
and advance, and **MUST NOT** require a full prose answer.

Open-ended questions **MUST NOT** be shoehorned into fixed options on any
surface.

The workflow **MUST NOT** re-ask context the user has already supplied earlier
in the interaction.

### Discovery escapes

The workflow **MUST** honor an explicit user request to accept all inferred
defaults and skip the remaining questions.

The workflow **MUST** honor an explicit user request to see all questions at
once instead of iterating.

The workflow **MUST NOT** lead with these escapes; iteration is the default and
the escapes are offered or honored on request.

### Read before author

After running `qualitymd init [path]` to scaffold a missing model, the workflow
**MUST** read the scaffolded file before authoring it with a file write, so a
single authoring pass does not fail a read-before-write guard.

### Model maturity vs lifecycle readiness

The close step's model-maturity judgment **MUST** use a term distinct from
"readiness," because the CLI `status` command reports a separate lifecycle
`readiness`.

The skill **MUST NOT** present the model-maturity classification and the CLI
lifecycle states as one blended list. Maturity levels (such as starter,
immature, and an evaluation-ready level) describe how developed the model is;
lifecycle states (missing, invalid, no runs yet, has history, needs
reconciliation) describe where the model sits in the evaluation lifecycle and
are owned by the CLI.

The close step **MAY** lean on the CLI `status` readiness signal plus a
condensed checklist to classify maturity, and **MUST** read the full top-10
guide only when the maturity call is borderline.

Setup completion output **MUST** report the model-maturity classification under
the maturity term, not under "readiness."

### Folder rename

The skill's `modes/` folder **MUST** be renamed to `workflows/` in both the
runtime skill (`skills/quality/`) and the spec mirror
(`specs/skills/quality-skill/`).

Every live reference to a `modes/<file>` path **MUST** be updated to the
`workflows/<file>` path across the runtime skill, the durable specs, the format
and project docs, and packaging files.

Append-only `log.md` files (for example [`specs/log.md`](../../specs/log.md) and
[`changes/log.md`](../../changes/log.md)) **MUST NOT** have their historical
`modes/` references rewritten; they record past state and stay frozen.

Runtime and durable text that refers to these files **SHOULD** use "workflow"
terminology consistent with the renamed folder. "Mode" **MAY** remain only where
it names internal dispatch state rather than the files themselves.

## Acceptance Criteria

- Setup presents all ten discovery questions on every run unless the user opts
  out, with none dropped or merged.
- Discovery presentation is keyed to the agent's own capabilities and names no
  specific question tool.
- With no structured affordance, setup iterates one question at a time, each
  carrying its default and confidence.
- With a structured tool, setup pages all ten through it and keeps open-ended
  items as free text.
- Accept-all and show-all-at-once escapes are honored on request and are not led
  with.
- Setup reads the `qualitymd init` scaffold before authoring it.
- The close step reports maturity under a term distinct from "readiness."
- Maturity labels and CLI lifecycle states are no longer presented as one list.
- The condensed close checklist exists and is what setup reads unless the call
  is borderline.
- The `modes/` folder is renamed to `workflows/` in the runtime skill and spec
  mirror, with all live path references updated and append-only logs left frozen.

## Durable spec changes

### To add

None.

### To modify

- `specs/skills/quality-skill/workflows/setup.md` (renamed from
  `modes/setup.md`) - replace the prompt-form contract with the agent-agnostic
  tiered-iteration and escape requirements, add the read-before-author
  requirement, and update the close contract to the maturity vocabulary and
  condensed checklist.
- `specs/skills/quality-skill/guides/top-10-quality-md-checks-md.md` - separate
  the model-maturity classification from CLI lifecycle states and specify the
  condensed close checklist.
- `specs/skills/quality-skill/quality-skill.md` - rename "readiness" to the
  maturity term only where it means the model-maturity judgment (keep CLI/tooling
  readiness intact), and update `modes/` path references to `workflows/`.
- `specs/skills/quality-skill/index.md` and
  `specs/skills/quality-skill/evaluation.md` - update `modes/` path references to
  `workflows/`.
- Relevant OKF logs and indexes under `specs/` - record durable spec updates
  when they land; the append-only `specs/log.md` keeps historical `modes/`
  references frozen.

### To rename

- `specs/skills/quality-skill/modes/` → `specs/skills/quality-skill/workflows/`
  (folder, including `setup.md`, `evaluate.md`, `update.md`, `log.md`, and
  `index.md`).

### To delete

None.
