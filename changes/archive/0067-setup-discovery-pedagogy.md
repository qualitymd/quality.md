---
type: Change Case
title: Setup discovery pedagogy
description: Make /quality setup discovery teaching-first — authored per-question background and how-to-change-later copy inline in the skill, ask every question even at high confidence, relabel confidence to Low/Med/High, and add a final review recap with a last-chance comment before writing QUALITY.md.
status: Done
tags: [skill, setup, ux, pedagogy]
timestamp: 2026-06-23T00:00:00Z
---

# Setup discovery pedagogy

A **Change Case** repositioning the `/quality setup` discovery step from
efficient inference-confirmation toward a guided teaching walkthrough. It adds
authored per-question copy to the runtime skill, makes setup ask every discovery
question even when its inferred default is high-confidence, simplifies the
confidence vocabulary to Low/Med/High, and adds a final review recap before
authoring. Detail lives in its
[functional spec](0067-setup-discovery-pedagogy/spec.md) and
[design doc](0067-setup-discovery-pedagogy/design.md).

## Motivation

The discovery questions carry a dual purpose the current workflow names but does
not fully serve: they capture the context the model needs *and* teach the user
the dimensions a quality model spans. Today the teaching half is left to the
agent to improvise per run — the skill only asserts the questions "do double
duty" — so the pedagogy is inconsistent across runs and impossible for a human
to tune. A field run logged in
[`acquire-roi-next/setup-observations.md`](../../../acquire-roi/acquire-roi-next/setup-observations.md)
also surfaced two specific rough edges: the inferred/assumed confidence
vocabulary reads awkwardly next to recommended defaults, and once inferences are
strong the per-question paging feels like pure overhead rather than instruction.

This case deliberately resolves that tension in favor of teaching. Setup runs
roughly once per project; spending extra interaction to make each dimension
legible — and to leave the user knowing *why* each answer shapes the model and
*how* to change it later — is worth more than minimizing round-trips. Four
deltas implement that stance:

1. **Authored per-question pedagogy.** Each discovery question gains authored
   background copy in the runtime skill: why the dimension matters and what it
   shapes in `QUALITY.md`, plus how the user can change it later. Authoring the
   copy in the skill (not improvising it) makes the teaching reproducible and
   human-tunable over time. The copy lives inline in the setup workflow so it
   travels with the procedure.
2. **Ask every question.** Because the pedagogical purpose applies to every
   dimension, setup asks each question even when its inferred default is
   high-confidence. The current "accept all defaults and skip the remaining
   questions" escape directly contradicts this and is removed or revised; the
   "show all at once" escape is retained.
3. **Low/Med/High confidence.** The `strongly inferred` / `weakly inferred` /
   `assumed` vocabulary is replaced with `Low` / `Med` / `High`. The evidence
   note is retained, so the "no evidence, pure default" meaning that `assumed`
   carried is preserved as e.g. `Low (no signal in repo)`.
4. **Final review recap.** After all questions, setup presents the full set of
   questions and final answers and invites one last free-text comment before
   writing `QUALITY.md`, giving the user a consolidated confirmation and a place
   for cross-cutting remarks.

These recur every setup run. Together they make setup's instruction value
explicit and tunable, and make its confidence signal legible.

## Scope

Covered:

- Add authored per-question copy for each of the ten discovery questions in the
  runtime setup workflow: a short background (why the dimension matters / what it
  shapes in `QUALITY.md`) and a how-to-change-later note, written as copy the
  agent presents, not as widget-field text.
- Establish, in the workflow framing, that setup optimizes for teaching over
  round-trip count, so the per-question pedagogy is not "optimized" back out.
- Require setup to ask every discovery question regardless of inferred-default
  confidence. Remove or revise the "accept all defaults → skip the remaining
  questions" escape and any other guidance that contradicts asking every
  question.
- Replace the confidence vocabulary `strongly inferred` / `weakly inferred` /
  `assumed` with `Low` / `Med` / `High`, retaining the per-item evidence note.
- Add a final review step after discovery: present the full question/answer set
  and invite one last free-text comment before writing `QUALITY.md`.
- Update the runtime skill and durable setup spec to match.

Deferred / non-goals:

- No QUALITY.md format change.
- No `qualitymd` CLI or Go code change. The CLI `status` lifecycle `readiness`
  field is unchanged.
- No change to **which** ten discovery questions are asked, their option sets, or
  their recommended defaults — only how they are presented, confirmed, and
  labeled.
- No change to the agent-agnostic presentation contract from 0065 (structured
  paging vs one-at-a-time iteration); this case layers pedagogy and the recap
  onto it.
- No new public skill workflow, and no change to setup's mutation boundary
  (`QUALITY.md` only; the 0066 feedback-log write, if it lands, is unaffected).
- No change to the `quality-md` self-check Area or the close/maturity step.

## Relationship to prior cases

This case revisits two things
[0065 — Setup discovery and close refinements](archive/0065-setup-discovery-and-close-refinements.md)
deliberately left alone. 0065 established the agent-agnostic presentation
contract but explicitly scoped out "no change to which questions are asked,
their defaults, or their confidence vocabulary" and preserved the accept-all
escape. This case keeps 0065's presentation tiers and changes what 0065 froze:
it relabels the confidence vocabulary and revises the accept-all escape, and it
adds the per-question pedagogy and final recap on top.

It is independent of
[0066 — Setup feedback log](0066-setup-feedback-log.md), which explicitly lists
"no change to which discovery questions setup asks, their defaults, or confidence
vocabulary" as a non-goal. The two cases touch the same setup files but
different concerns and do not conflict.

## Affected artifacts

Durable **spec** changes are itemized in the functional spec's
[Durable spec changes](0067-setup-discovery-pedagogy/spec.md#durable-spec-changes)
section. The index below is the full skimmable list to reconcile before
In-Review.

### Code

- [x] **None.** Setup discovery is skill-driven; no `qualitymd` CLI or Go code
      change. The CLI `status` lifecycle `readiness` field is intentionally
      unchanged. Confirmed no generated example or CLI output encodes the old
      confidence vocabulary.

### Durable specs

- [x] [`specs/skills/quality-skill/workflows/setup.md`](../specs/skills/quality-skill/workflows/setup.md)
      - required authored per-question pedagogy, required asking every question and
      removed the accept-all escape, changed the confidence-vocabulary
      `MUST` to `Low`/`Med`/`High`, and added the final-review-recap requirement.
- [x] [`specs/skills/quality-skill/quality-skill.md`](../specs/skills/quality-skill/quality-skill.md)
      - reviewed the "confidence-labeled defaults" framing (line ~384); it is
      generic and does not pin the old terms, so left unchanged.
- [x] OKF logs and indexes under [`specs/`](../specs/log.md) - recorded durable
      spec updates in [`specs/log.md`](../specs/log.md) and the
      [workflows log](../specs/skills/quality-skill/workflows/log.md). Append-only
      `log.md` files keep their historical confidence-vocabulary references frozen.

### Durable docs

- [x] `docs/guides/use-quality-skill.md` - **superseded.** This guide was removed
      independently (concurrent docs cleanup folding its operational guidance into
      the README and the skill itself), so 0067 lands with no durable-docs delta.
- [x] [`docs/log.md`](../docs/log.md) - no 0067 docs entry needed; the only
      candidate doc was removed independently.

### Bundled skill

- [x] [`skills/quality/workflows/setup.md`](../skills/quality/workflows/setup.md)
      - the primary change: added authored per-question background and
      how-to-change-later copy inline, added the teaching-first framing note,
      required asking every question and removed the accept-all escape,
      relabeled confidence to `Low`/`Med`/`High` in the setup brief and presentation
      guidance, and added the final-review-recap step before authoring.
- [x] [`skills/quality/SKILL.md`](../skills/quality/SKILL.md) - reviewed setup
      dispatch/run-frame wording; no behavioral change needed. Skill
      `metadata.version` bump deferred to the release cut (consistent with 0065/0066).

### Release

- [ ] [`CHANGELOG.md`](../CHANGELOG.md) - add the user-facing entry when
      implementation lands.

## Children

- [Functional spec](0067-setup-discovery-pedagogy/spec.md) - what the
  per-question pedagogy, ask-every-question rule, confidence relabel, and final
  review recap must do.
- [Design doc](0067-setup-discovery-pedagogy/design.md) - how the pedagogy copy,
  escape revision, confidence relabel, and recap step are shaped, and the
  alternatives weighed.

## Status

`Done`. Implemented skill-only with no code change: the runtime setup workflow
gained authored per-question pedagogy and the teaching-first framing, asks every
question every run (accept-all-and-skip escape removed; per-question fast confirm
and show-all-at-once retained), uses `Low`/`Med`/`High` confidence with the
evidence note, and recaps answers before authoring. The durable setup spec was
synced with promoted rationale annotations; the parent skill spec's generic
"confidence-labeled defaults" phrasing was reviewed and left unchanged. The one
listed public doc (`use-quality-skill.md`) was removed independently by
concurrent docs cleanup, so this case lands with no durable-docs delta.
`CHANGELOG.md` and the OKF logs/indexes are updated. The In-Review review gate
was collapsed at the user's explicit direction. Moved to `archive/`.
