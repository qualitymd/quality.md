---
type: Design Doc
title: Render interactions through native affordances — design
description: How the setup workflow's "render through your own interaction capabilities" pattern is lifted into the shared /quality interaction contract, the alternatives weighed, and the risks.
status: Draft
timestamp: 2026-06-26T00:00:00Z
---

# Render interactions through native affordances — design

## Context

This design answers the
[0123 functional spec](spec.md). The spec requires the shared `/quality`
interaction contract to treat each interaction as an intent rendered through a
fit-for-purpose native affordance when present, with the prose renderings as the
text fallback (R1–R5), to stop double-gating harness-authorized mutations (R6),
and to apply both rules through the inheriting workflow specs (R7–R9).

The whole change is to prose contracts and skill instructions: durable specs
under `specs/skills/quality-skill/`, the bundled skill under `skills/quality/`,
and the already-applied guide. There is no Go code: the CLI is non-interactive by
contract, and rendering is the agent's runtime job.

## Approach

**Single home for the pattern.** The `setup` workflow already contains the
correct, agent-agnostic three-tier pattern in its "How to present them" block:
Markdown rendering → structured question tool (page questions through it, keep the
teaching in the message, do not force fixed options) → no-affordance fallback
(iterate one at a time). The core move is to *lift* that text — generalized away
from setup's specific questions — into the shared User Interaction Contract in
`SKILL.md` and its `quality-skill.md` spec mirror, so all workflows inherit it
(R1–R3, R9).

**Organize the contract as intent → affordance → fallback.** In the shared
contract, the existing "number the options / use `[y]`/`[n]`" prose stops being
the headline rule and becomes the *fallback rendering* under a short
progressive-enhancement frame (R4). The frame states: name the interaction
intent; render through the richest fit-for-purpose affordance the runtime
exposes; always author the text fallback; keep the teaching in the message, not
the widget labels.

**Reference the taxonomy, do not re-enumerate it.** The affordance categories
(single-select, multi-select, confirm/approve, plan-or-diff review, harness
authorization prompt, free-text, progress) and the not-fit-for-purpose tests are
defined once in the guide's new "Channels and progressive enhancement" section.
The contract names the categories it acts on and defers to the guide for the full
list (R5), keeping one source of truth.

**Gate policy gains the no-double-gate rule.** The shared decision-brief
requirement adds: where the harness authorizes the mutation itself, render the
confirmation through that gate rather than stacking a prose brief, but never drop
confirmation (R6). This is itself an application of the progressive-enhancement
frame to the authorization-prompt affordance category.

**Inheriting specs get light, explicit edits.** `setup` deduplicates its block to
a reference (R9) and tags its multi-value questions as single-select intents
(R7); recommendation follow-up tags its outcome choice and applies the
no-double-gate rule to issue creation (R6, R7); `update` and `evaluate` inherit
the reshaped contract with at most a one-line note. The human-context checkpoint
is explicitly locked as free-text so the reframe cannot widget-ize it (R8).

## Spec response

- **R1–R3 (progressive-enhancement contract, teaching in the message,
  agent-agnostic):** satisfied by lifting the setup three-tier block into the
  shared contract with capability-conditional phrasing and no named UI.
- **R4 (prose as fallback):** satisfied by demoting the numbered-option and
  `[y]`/`[n]` templates under the progressive-enhancement frame while keeping the
  recommended-first and `1`/`y`/`n` rules intact for the fallback.
- **R5 (taxonomy):** satisfied by naming the acted-on categories and referencing
  the guide's enumeration.
- **R6 (no double-gate):** satisfied by an added clause in the shared
  decision-brief requirement plus its application in the issue-creation and
  tooling-mutation gates, with the "still confirm" guard explicit.
- **R7–R9 (inherited applications):** satisfied by the per-workflow tagging and
  the setup deduplication described above.

## Alternatives

- **Name a specific question tool in the contract.** Rejected: it breaks the
  agent-agnostic posture the project requires, and the skill cannot assume any
  named UI exists at runtime. Capability-conditional phrasing covers the same
  ground without binding to one harness.
- **Leave the prose-only contract and add a one-line "use a picker if you have
  one" note.** Rejected: a note bolted onto a contract that still *mandates* the
  prose form is the status quo that produced the gap — the default stays
  prose-only and workflows keep emitting numbered prose. The reframe has to move
  the prose to the fallback slot to change the default.
- **Copy the three-tier pattern into each workflow instead of sharing it.**
  Rejected by say-it-once: that is how setup ended up as the only conformant
  workflow, and N copies drift. The pattern belongs in the one shared contract.
- **Build a shared renderer or interaction-template library.** Rejected as YAGNI
  and out of layer: rendering happens in the agent at runtime, not in the
  deterministic CLI, so there is no code surface to centralize. The shared
  *contract* is the right abstraction; a code library is not.
- **Convert every closed choice, including the human-context checkpoint, to a
  picker.** Rejected: open-cardinality, confirm-or-correct surfaces are the case
  where free text is the fit-for-purpose affordance. R8 locks this so the reframe
  does not over-correct.

## Trade-offs & risks

- **Over-widgetizing.** An agent could read "render through a native affordance"
  as "always force a picker," including for open-ended input. Mitigated by the
  guide's not-fit-for-purpose tests, the named free-text category, and R8's
  explicit checkpoint lock.
- **Reading R6 as "drop the confirmation."** The no-double-gate rule could be
  misread as permission to skip confirming a mutation. Mitigated by the explicit
  "MUST still not mutate without confirmation" guard in both the spec and the
  contract text.
- **A more abstract contract is harder for a weak agent to render well.** A
  capability-conditional contract asks more of the agent than a fixed template.
  Mitigated because the text fallback remains fully specified as the floor: an
  agent that does nothing clever still produces the previous, known-good prose.
- **Drift between guide and contract.** The contract references the guide's
  taxonomy. Mitigated by keeping the enumeration in the guide only and naming
  just the acted-on categories in the contract.

## Open questions

- Inline vs referenced taxonomy in `quality-skill.md` (mirrors the spec's open
  question). Leaning toward referencing the guide; settle during implementation
  to whichever keeps a single source of truth without forcing the reader to leave
  the contract for the common categories.
