---
type: Design Doc
title: Setup discovery and close refinements
description: How agent-agnostic discovery, read-before-author, the maturity vocabulary, and the modes/ to workflows/ rename are shaped.
tags: [skill, setup, ux]
timestamp: 2026-06-23T00:00:00Z
---

# Setup discovery and close refinements - design doc

## Context

A first field run of `/quality setup` (recorded in an external setup-observations
report) completed end to end but logged four recurring frictions. None block,
but each costs every run, and one can mislead routing. The
[functional spec](spec.md) settles _what_ must change; this doc settles _how_,
and records the alternatives weighed in discussion.

The four items, and the durable surfaces they touch:

1. Discovery presentation — `workflows/setup.md` (Ask Discovery Questions).
2. Read-before-author — `workflows/setup.md` (Write QUALITY.md).
3. Maturity vs lifecycle readiness — `workflows/setup.md` (Verify and Close) and
   the top-10 guide.
4. `modes/` → `workflows/` rename — the runtime skill and the spec mirror.

## Approach

### Discovery presentation — three tiers keyed to agent capability

Replace the current "compact prompt or short sequence" guidance (introduced by 0064) with one agent-agnostic rule. The agent picks a tier from its _own_
interaction capabilities, so no specific tool is ever named:

- **Tier 1 — structured question tool available.** Page all ten questions
  through it across as many rounds as the tool's item/option caps require. The
  four open-ended questions (primary users, maintainers/collaborators, other
  stakeholders, missing context) stay free text; they are never coerced into
  fixed options.
- **Tier 2 — no structured affordance (the default).** Iterate the questions one
  at a time. Each step carries the question's recommended default and confidence
  so the user confirms or terse-corrects and advances. This is the most
  digestible form and best serves the pedagogical goal — each dimension gets real
  attention instead of being skimmed in a wall of ten.
- **Tier 3 — escapes, honored on request, not led with.** The user may say
  "accept all defaults" to stop early, or "show me all at once" for the batch
  form. Iteration stays the default so the teaching value is not lost to a
  reflexive "looks fine."

The invariant across tiers: all ten questions are surfaced; none is dropped,
merged, or silently defaulted to fit a surface.

### Read-before-author

In the Write step, after `qualitymd init` scaffolds a missing model, add an
explicit instruction to read the scaffolded file before authoring it with a file
write. `init` creates the file via a tool, which does not satisfy the harness's
read-before-write guard, so without this the first authoring write fails and
costs a no-op read. One line removes a guaranteed round-trip.

### Maturity vs lifecycle readiness

The close step today classifies model "readiness" as
`starter | immature | ready to evaluate`, while CLI `status` emits a _lifecycle_
`readiness` (`ready-to-evaluate` = valid, no runs yet). Same word, different
axes. The top-10 guide compounds this by listing maturity labels and lifecycle
labels (`missing`, `invalid`, `has evaluation history`, `needs reconciliation`)
in one classification.

Resolution:

- Rename the skill's model-maturity judgment off "readiness" to a distinct term
  (working name: **maturity**), used in the close step and completion output.
- Keep CLI lifecycle `readiness` exactly as is — it is the CLI's, and other skill
  surfaces correctly cite it.
- In the top-10 guide, separate the two axes: maturity describes how developed
  the model is; lifecycle states are owned by CLI `status`.
- Let the close step lean on the CLI `status` readiness signal plus a condensed
  checklist, reading the full top-10 guide only when the maturity call is
  borderline. This removes the routine full-guide read when the CLI already
  classifies lifecycle state.

### `modes/` → `workflows/` rename

0064 deferred this rename to "a later design"; this is that design. Rename the
folder in both `skills/quality/` and `specs/skills/quality-skill/`, and update
every _live_ path reference (`SKILL.md`, `quality-skill.md`, `index.md`,
`evaluation.md`, and intra-folder links). Align surrounding text to "workflow"
terminology; "mode" survives only where it names internal dispatch state, not the
files.

Append-only `log.md` files keep their historical `modes/` references frozen, per
the OKF rule that logs record past state. Their links to the old path are
accepted as historical record rather than rewritten — the alternative (rewriting
history to keep links live) violates the append-only contract.

## Alternatives

- **Map the prompt explicitly to a 4×4 question tool.** Rejected: this was the
  original instinct, but it bakes one agent's UI into the skill — the opposite of
  the agent-agnostic goal.
- **Topic-grouped iteration (four steps) as the no-UI default.** Considered and
  set aside in favor of one-at-a-time. Grouping trades fewer round-trips for less
  pedagogical attention per dimension; the user chose one-at-a-time as the
  default, with grouping not needed once tiers are capability-keyed.
- **Drop the questions the agent is highly confident about.** Rejected: the
  questions serve a pedagogical purpose and capture critical context; confidence
  in a default is a reason to _recommend_ it, not to hide the question.
- **Keep "readiness" and just annotate the difference.** Rejected: one word for
  two axes is the root confusion; a distinct term is cleaner than a footnote.
- **Rewrite `modes/` links in historical logs so nothing breaks.** Rejected:
  violates the append-only log contract; frozen historical paths are correct
  past-state record.

## Trade-offs & Risks

- One-at-a-time iteration can mean up to ten turns on a plain surface. Mitigated
  by per-step defaults (terse confirm) and the accept-all escape, so an expert is
  never trapped in a ten-turn interrogation.
- The exact maturity label set (whether the top level stays "ready to evaluate"
  or becomes something like "evaluation-ready"/"substantive") is left to
  implementation; the binding requirement is only that maturity is distinct from
  lifecycle readiness and not blended into one list.
- The rename causes broad link churn and lands concurrently with 0064. Whichever
  case lands second must reconcile the durable setup spec to carry both the
  tiered-iteration contract and the `workflows/` paths. Risk is mechanical, not
  semantic.

## Open Questions

- Final maturity term and label set (maturity vs another word; whether to rename
  the "ready to evaluate" level). Deferred to implementation.
- Whether `SKILL.md` dispatch language should retire "mode" entirely or keep it
  for internal routing state only. Lean: keep "mode" solely for internal dispatch
  state, use "workflow" for the files.
