---
type: Change Case
title: Setup discovery and close refinements
description: Make setup discovery agent-agnostic with one-question-at-a-time iteration, read the scaffolded model before authoring, disentangle model maturity from CLI lifecycle readiness, and rename the skill modes/ folder to workflows/.
status: Done
tags: [skill, setup, ux]
timestamp: 2026-06-23T00:00:00Z
---

# Setup discovery and close refinements

A **Change Case** refining three points in the `/quality setup` workflow that a
field test against a real repository surfaced: the discovery prompt assumes a
particular question UI, authoring after `qualitymd init` costs a wasted
read-before-write round-trip, and the close step's "readiness" vocabulary
collides with the CLI's lifecycle readiness. Detail lives in its
[functional spec](0065-setup-discovery-and-close-refinements/spec.md) and
[design doc](0065-setup-discovery-and-close-refinements/design.md).

## Motivation

A first end-to-end run of `/quality setup` against an external monorepo
completed with no blocking errors, but logged three avoidable frictions:

1. **Discovery assumes a question UI.** The workflow prescribes a ten-item
   compact prompt (and a four-group short sequence). When the agent reached for a
   structured question tool, that tool's item/option caps forced the agent to
   compress and drop questions, re-improvised each run. The prompt form is
   implicitly shaped to one agent's UI rather than being agent-agnostic, and the
   ten questions carry both critical setup context and a pedagogical purpose, so
   silently dropping or merging them is costly.
2. **Read-before-write round-trip.** Setup scaffolds with `qualitymd init` (a
   tool-created file), then authors with a file write. The harness's
   "read this file before writing" guard rejects the first write, forcing a
   no-op read every run.
3. **"Readiness" collision at close.** The close step classifies model
   "readiness" as `starter | immature | ready to evaluate`, while the CLI
   `status` command emits a *lifecycle* `readiness` (`ready-to-evaluate` meaning
   only "valid, no runs yet"). The two are different axes under one word, and the
   top-10 checklist blends maturity labels with lifecycle labels in one list. An
   agent can read the CLI's `ready-to-evaluate` as the model being mature, which
   it does not mean.

These are small, but they recur every setup run, and the third can mislead
routing. Fixing them makes setup portable across coding agents, cheaper per run,
and unambiguous at close.

## Scope

Covered:

- Replace the prompt-form guidance with an agent-agnostic rule keyed to the
  agent's own interaction capabilities, not a named tool. Always present all ten
  discovery questions; never drop or merge one away.
- When no structured question affordance exists, iterate the discovery questions
  one at a time, each carrying its recommended default and confidence.
- When a structured question tool exists, page all ten through it within its
  caps; keep open-ended questions as free text.
- Preserve early-exit escapes (accept all defaults; show all at once) on user
  request, without leading with them.
- Add an explicit step to read the `qualitymd init` scaffolded model before
  authoring it with a file write.
- Disentangle the skill's model-maturity judgment from the CLI's lifecycle
  readiness: rename the maturity axis off the word "readiness," and stop blending
  maturity and lifecycle labels in one classification list.
- Let the close step lean on the CLI's readiness signal plus a condensed
  checklist, reading the full top-10 guide only when the maturity call is
  borderline.
- Rename the skill's `modes/` folder to `workflows/` in both the runtime skill
  and the `specs/` mirror, update every live path reference, and align the
  surrounding vocabulary to "workflow" for these files. Append-only `log.md`
  files keep their historical `modes/` references frozen as past-state record.

Deferred / non-goals:

- No QUALITY.md format change.
- No `qualitymd` CLI or Go code change. In particular, the CLI `status`
  lifecycle `readiness` field and its values are unchanged.
- No new CLI command or interactive CLI workflow.
- No change to which questions are asked, their defaults, or their confidence
  vocabulary — only how they are presented.
- No automatic evaluation, quality-log writing, issue creation, or automation
  configuration during setup.

## Relationship to 0064

This case modifies the **Prompt Form** requirements introduced by
[0064 — Structured setup workflow](archive/0064-structured-setup-workflow.md)
(the "compact prompt or short sequence" framing) and takes up the `modes/` →
`workflows/` filesystem rename 0064 explicitly deferred to "a later design." 0064
has landed and archived; this case must reconcile the durable setup spec so it
carries the agent-agnostic tiered-iteration contract from this case (not 0064's
compact/short-sequence wording) and the renamed `workflows/` paths.

## Affected artifacts

Durable **spec** changes are itemized in the functional spec's
[Durable spec changes](0065-setup-discovery-and-close-refinements/spec.md#durable-spec-changes)
section. The index below is the full skimmable list to reconcile before
In-Review.

Code:

- [x] No planned `qualitymd` CLI or Go code changes; implementation analysis
      should confirm setup remains skill-driven and that no generated examples or
      CLI output encode the old prompt-form wording. The CLI `status` lifecycle
      `readiness` field is intentionally unchanged.

Specs:

- [x] [`specs/skills/quality-skill/modes/setup.md`](../specs/skills/quality-skill/modes/setup.md)
      → `specs/skills/quality-skill/workflows/setup.md` - rename, then replace
      the prompt-form contract with agent-agnostic tiered iteration, add the
      read-before-author step, and update the close contract to the maturity
      vocabulary.
- [x] [`specs/skills/quality-skill/modes/`](../specs/skills/quality-skill/modes/index.md)
      → `specs/skills/quality-skill/workflows/` - rename the whole folder
      (`evaluate.md`, `update.md`, `log.md`, `index.md`) and update its
      `index.md`/intra-folder links.
- [x] [`specs/skills/quality-skill/guides/top-10-quality-md-checks-md.md`](../specs/skills/quality-skill/guides/top-10-quality-md-checks-md.md)
      - separate the model-maturity classification from CLI lifecycle states and
      specify the condensed close checklist.
- [x] [`specs/skills/quality-skill/quality-skill.md`](../specs/skills/quality-skill/quality-skill.md)
      - rename "readiness" to the maturity term only where it means the skill's
      model-maturity judgment (keep CLI/tooling readiness wording intact), and
      update its `modes/` path references to `workflows/`.
- [x] [`specs/skills/quality-skill/index.md`](../specs/skills/quality-skill/index.md)
      and [`specs/skills/quality-skill/evaluation.md`](../specs/skills/quality-skill/evaluation.md)
      - update `modes/` path references to `workflows/`.
- [x] OKF logs and indexes under [`specs/`](../specs/log.md) - record durable
      spec updates when they land. The append-only [`specs/log.md`](../specs/log.md)
      keeps its historical `modes/` references frozen; do not rewrite past
      entries.

Runtime skill and docs:

- [x] [`skills/quality/modes/`](../skills/quality/modes/) →
      `skills/quality/workflows/` - rename the folder (`setup.md`,
      `evaluate.md`, `update.md`).
- [x] `skills/quality/workflows/setup.md` (renamed from `modes/setup.md`) -
      rewrite the discovery section to the tiered-iteration rule, add the
      read-before-author step, and update the close step to the maturity
      vocabulary and condensed checklist.
- [x] [`skills/quality/guides/top-10-quality-md-checks.md`](../skills/quality/guides/top-10-quality-md-checks.md)
      - disentangle maturity from lifecycle labels and add the condensed
      checklist setup reads at close.
- [x] [`skills/quality/SKILL.md`](../skills/quality/SKILL.md) - update `modes/`
      path references to `workflows/`, align dispatch wording, and review close /
      routing wording for the maturity-vs-lifecycle distinction; keep CLI status
      `readiness` references intact.
- [x] [`docs/guides/use-quality-skill.md`](../docs/guides/use-quality-skill.md)
      - align the "ready to evaluate" close wording with the maturity term.
- [x] [`README.md`](../README.md) and
      [`npm/quality.md/README.md`](../npm/quality.md/README.md) - update setup
      wording only if public phrasing changes.
- [x] [`CHANGELOG.md`](../CHANGELOG.md) - add the unreleased entry when the
      implementation lands.

## Children

- [Functional spec](0065-setup-discovery-and-close-refinements/spec.md) - what
  the discovery presentation, read-before-author step, and maturity vocabulary
  must do.
- [Design doc](0065-setup-discovery-and-close-refinements/design.md) - how the
  tiered iteration, escapes, maturity/lifecycle disentanglement, and `modes/` →
  `workflows/` rename are shaped, and the alternatives weighed.

## Status

`Done`. Implementation landed across the runtime skill, durable specs, public
docs, changelog, and the `modes/` → `workflows/` rename. The In-Review review
gate was collapsed at the user's explicit direction. Verified with `mise run
check` (Go lint/vet/test, markdown format, npm bundle link check). Archived.
