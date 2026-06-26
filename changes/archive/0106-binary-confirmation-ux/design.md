---
type: Design Doc
title: Binary confirmation UX — design
description: Implementation approach for applying y/n binary confirmation guidance to the /quality skill.
tags: [skill, quality, ux, agent-mediated]
timestamp: 2026-06-26T00:00:00Z
---

# Binary confirmation UX — design

Companion to the [Binary confirmation UX](../0106-binary-confirmation-ux.md)
change case and its [functional spec](spec.md).

## Context

This is a documentation-and-runtime-skill alignment change. The behavior under
design is the agent's interaction contract: which terse response the prompt
visibly asks for when a user is authorizing a mutating action.

The shared UX guide has already been clarified. The remaining work is to apply
that rule in the durable `/quality` skill specs and the bundled runtime skill
files that agents actually read.

## Approach

Patch the shared contract first, then patch concrete gates:

1. Update the durable parent `/quality` skill spec and runtime `SKILL.md` to say
   that true binary confirmations use visible `y`/`n`, while non-binary
   closed-choice prompts keep the `1`-first numbered pattern.
2. Extend the shared decision-brief guidance so a binary decision brief includes
   an explicit `Answer` line or equivalent wording.
3. Patch the update workflow's update-plan brief with a concrete `y`/`n` answer
   path.
4. Patch recommendation follow-up's local apply and external issue creation
   briefs with concrete `y`/`n` answer paths.
5. Patch setup's fallback existing-file update brief with a concrete `y`/`n`
   answer path.
6. Update append-only skill/spec logs for the durable and runtime files touched.
7. Verify with targeted searches for stale `reply 1 to apply`-style wording and
   for the new `y`/`n` answer guidance, then run the Markdown formatting check.

The implementation should not alter setup discovery, evaluation ambiguity, or
apply-vs-handoff prompts except to avoid accidentally treating them as binary.

## Spec response

The shared interaction requirements are satisfied by changing
`specs/skills/quality-skill/quality-skill.md` and `skills/quality/SKILL.md`.
Those files own the umbrella rule and prevent future workflow guidance from
reverting to all-closed-choices-are-`1` wording.

The workflow-gate requirements are satisfied by local templates in
`skills/quality/workflows/update.md`,
`skills/quality/guides/recommendation-follow-up.md`, and
`skills/quality/workflows/setup.md`, with matching durable spec changes under
`specs/skills/quality-skill/`.

The preserved-numbered-prompts requirements are satisfied mostly by leaving
existing numbered prompt guidance intact. Where broad wording is edited, keep it
scoped so it does not imply setup discovery choices, evaluation routing, or
apply-vs-handoff selection should become yes/no prompts.

## Alternatives

- **Patch only `skills/quality/workflows/update.md`.** Rejected because the same
  binary confirmation shape appears in recommendation follow-up and setup, and
  the shared contract would still imply `1` for every small closed-choice prompt.
- **Use `1`/`2` for binary gates.** Rejected because the clarified UX guide treats
  true binary confirmations as yes/no questions. Numbering is still better for
  non-binary choices.
- **Convert setup final review to `y`/`n`.** Rejected because that gate is a
  review-and-correction interaction, not a plain yes/no confirmation. The
  existing `looks good` fast path carries the intended affordance.

## Trade-offs & risks

The main risk is over-applying the exception. A prompt is not binary just because
one path is recommended: lifecycle, risk, rating scale, scope, and handoff
choices all carry more than yes/no information. The implementation should make
the exception visible but narrow.

Accepting aliases while displaying `y`/`n` may make the runtime guidance slightly
longer. That is acceptable because it keeps the visible prompt simple without
making the agent brittle when users answer naturally.

## Open questions

None.
