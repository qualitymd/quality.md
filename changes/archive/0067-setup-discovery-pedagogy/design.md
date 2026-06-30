---
type: Design Doc
title: Setup discovery pedagogy
description: How the per-question teaching copy, ask-every-question rule, Low/Med/High confidence, and final review recap are shaped for /quality setup discovery.
tags: [skill, setup, ux, pedagogy]
timestamp: 2026-06-23T00:00:00Z
---

# Setup discovery pedagogy - design doc

## Context

The [functional spec](spec.md) settles _what_ must change in the `/quality setup`
discovery step; this doc settles _how_, and records the alternatives weighed in
discussion. The change repositions discovery from efficient
inference-confirmation toward a guided teaching walkthrough, motivated by a field
run logged in `acquire-roi-next/setup-observations.md` and by the observation
that the discovery questions' pedagogical purpose is asserted but not actually
served by the current skill.

The four deltas all land in one surface — the discovery and write-prep portion of
`skills/quality/workflows/setup.md`, mirrored in
`specs/skills/quality-skill/workflows/setup.md`:

1. Per-question teaching copy — Ask Discovery Questions.
2. Ask every question / escape revision — Ask Discovery Questions (How to present
   them, Escapes).
3. `Low`/`Med`/`High` confidence — Read Context (setup brief) and Ask Discovery
   Questions (the per-question default labels).
4. Final review recap — a new step between Ask Discovery Questions and Write
   QUALITY.md.

## Approach

### Per-question teaching copy — inline, authored, presented as prose

The ten-question block in `workflows/setup.md` today gives each question a one-line
stem plus `Recommended: <default> (<confidence>)`. This case expands each entry to
carry two authored fields:

- **Why it matters** — the purpose of the dimension and what it shapes in
  `QUALITY.md` (e.g. risk tolerance drives modeling rigor and which factors earn
  requirements; root area sets the model's boundary and the `quality-md`
  self-check Area).
- **How to change it later** — the concrete path to revising that answer after
  setup (edit the relevant body section / re-run `/quality setup` / adjust the
  factor set), reinforcing that `QUALITY.md` is a living document.

The copy lives **inline** in the workflow (the user's decision), so the teaching
text travels with the procedure and a single read loads both the question and its
instruction. It is authored prose the agent presents, _not_ widget text: the
spec forbids confining it to a structured tool's option/description fields, which
are too small to teach in.

This composes with the 0065 presentation tiers rather than replacing them:

- **One-at-a-time (default):** emit the question's "why it matters" before the
  question and the "how to change later" with or after it, then take the answer.
  This is the tier the pedagogy serves best — each dimension gets a real beat.
- **Structured question tool:** present the teaching prose in the message that
  precedes the tool call (or between paging rounds), and keep the tool itself for
  the choice capture. The prose carries the pedagogy; the widget carries the
  selection.
- **Show-all-at-once (on request):** the batch prompt includes each question's
  teaching copy inline.

A framing line in the workflow states that setup optimizes for teaching the
quality-model dimensions over minimizing round-trips. This is a guardrail: it
exists so a later contributor reading the efficiency complaint in the
observations report does not "optimize" the pedagogy back out.

### Ask every question — remove the skip escape, keep the cheap paths

The spec requires asking all ten questions every run regardless of
inferred-default confidence. The 0065 "accept all defaults and skip the remaining
questions" escape directly contradicts that and is removed. What stays:

- **Per-question fast confirm.** The user can accept a question's recommended
  default and advance without writing prose. This keeps an expert from being
  trapped in a ten-turn interrogation while still presenting every question and
  its teaching copy — the pedagogy survives a terse "yes."
- **Show-all-at-once.** Honored on request, not led with. It still surfaces every
  question, so it does not conflict with ask-every-question.

So the cost control shifts from "skip questions" to "answer them cheaply." High
confidence in a default becomes a reason to _recommend_ it firmly, never to hide
the question.

### `Low`/`Med`/`High` confidence

Replace the `strongly inferred` / `weakly inferred` / `assumed` vocabulary with
`Low` / `Med` / `High` in the setup brief and in each question's recommended
default. The per-item evidence note is retained, which is what lets the simpler
scale carry the full meaning:

- `strongly inferred` → `High (<evidence>)`
- `weakly inferred` → `Med (<evidence>)`
- `assumed` → `Low (no signal in repo)` — the spec requires a no-evidence default
  to be labeled `Low` and to name the absence of evidence, so the distinct
  "pure default, no evidence" signal `assumed` carried is preserved in the note
  rather than the label.

The old vocabulary is a normative `MUST` in the setup spec; the only other live
reference (`quality-skill.md` "confidence-labeled defaults") is generic and
survives unchanged. The top-10 guides' "assumed" is the unrelated adjective and is
not touched.

### Final review recap

Insert a recap step between Ask Discovery Questions and Write QUALITY.md. It:

- lists every question with its final answer (a single consolidated screen);
- invites one last free-text comment or correction;
- applies any correction before authoring; and
- does not require a comment to proceed.

This is a confirmation-and-teaching-recap moment, not a replacement for asking
each question. It also gives cross-cutting remarks — ones that did not fit a
single question — a home, which the mid-flow missing-context question (9) cannot
fully serve.

## Alternatives

- **Pedagogy copy in a `resources/` file.** Considered; the user chose inline so
  the teaching text travels with the procedure and is human-tunable in place.
  Inline costs a larger question block (~3× the current section) but keeps one
  source of truth for the step.
- **Improvise teaching per run (status quo).** Rejected: inconsistent across runs
  and impossible for a human to tailor, which is the whole point of authoring it.
- **Cram pedagogy into structured-tool option descriptions.** Rejected: those
  fields are too small and vary per agent; the spec mandates prose presentation
  instead.
- **Keep the accept-all-and-skip escape.** Rejected: it directly contradicts
  ask-every-question. The per-question fast confirm preserves its speed benefit
  without skipping the teaching.
- **Drop questions the agent is highly confident about.** Rejected for the same
  reason 0065 rejected it — confidence in a default justifies recommending it, not
  hiding the dimension. This case strengthens that from "present all ten" to "ask
  all ten."
- **Numeric or percentage confidence.** Rejected: false precision for an inferred
  judgment; `Low`/`Med`/`High` plus an evidence note is legible and honest.
- **Make the recap the only confirmation.** Rejected: it would turn discovery back
  into a single batch screen and lose the per-question teaching beats. The recap
  supplements iteration.

## Trade-offs & Risks

- **More round-trips.** One-at-a-time with no skip can mean up to ten beats plus
  the recap on a plain surface. This is the intentional trade — setup runs ~once
  per project and the framing line makes the choice explicit. The per-question
  fast confirm and show-all-at-once keep an expert from being trapped.
- **Larger workflow file.** The inline copy roughly triples the discovery section.
  Accepted for single-source-of-truth and human tunability.
- **Authored copy quality.** The teaching value now depends on the wording in the
  skill. The exact prose is left to implementation and expected to be
  human-tailored over time; the spec binds only the presence and content shape
  (purpose + how-to-change-later), not the final wording.
- **Structured-surface fit.** Agents whose only affordance is a constrained
  question widget must emit the pedagogy and recap as surrounding prose. The spec
  already requires this; the risk is an agent that under-emits the prose, which is
  a conformance matter, not a design gap.

## Open Questions

- Whether the recap should also restate each answer's confidence label or only the
  answers. Lean: answers only, since the user has already seen and confirmed each
  default; confidence is a discovery-time aid, not a confirmation-time one.
- Exact authored wording of the per-question copy. Deferred to implementation and
  ongoing human tailoring.
- If 0066 (setup feedback log) lands, whether the recap is a natural prompt point
  to mention the feedback log. Out of scope here; noted so the two steps can be
  sequenced sensibly later.
