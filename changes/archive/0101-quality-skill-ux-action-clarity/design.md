---
type: Design Doc
title: Quality skill UX action clarity — design
description: Implementation approach for aligning /quality prompt shapes with the agent-mediated UX guide.
tags: [skill, quality, ux, agent-mediated]
timestamp: 2026-06-26T00:00:00Z
---

# Quality skill UX action clarity — design

Companion to the
[Quality skill UX action clarity](../0101-quality-skill-ux-action-clarity.md)
change case and its [functional spec](spec.md).

## Context

This change is a documentation-and-runtime-skill alignment pass. The product
behavior is the agent output contract, not CLI logic. The implementation must
make the durable `/quality` skill specs and the bundled runtime guidance say the
same thing at the concrete prompt shapes that failed review.

## Approach

Use the durable specs as the lead contract, then mirror their exact interaction
shape into the runtime skill files:

1. Patch the shared `/quality` interaction contract for explicit shortest answer
   paths, run frames for public workflows, decision-brief coverage, and code-span
   precision.
2. Patch each workflow-specific durable spec where the shared contract is too
   general to prevent the observed failure.
3. Patch the corresponding runtime skill files with concrete prompt templates or
   procedural instructions.
4. Update append-only spec/skill logs for the changed durable surfaces.
5. Verify with targeted search for the required prompt affordances, then run the
   Markdown formatting check.

Keep the agent-mediated UX guide unchanged unless the implementation uncovers a
guide ambiguity. The reviewed failures are skill-conformance gaps, not guide
gaps.

## Spec response

The shared-contract requirements are satisfied by updating
`quality-skill.md` and `SKILL.md` once, then avoiding duplicated explanation in
workflow files except where a specific prompt shape needs local enforcement.

Setup needs the most concrete runtime template changes: the root/domain
questions gain `Answer` lines, the checkpoint becomes a correction-first table,
and the final review becomes a decision brief before `QUALITY.md` writes.

Evaluate scope ambiguity and stop responses get numbered options and explicit
`Answer` lines. The durable evaluation scope spec owns the ambiguity contract;
the evaluate workflow spec owns the stop-response shape.

Update gains an initial run frame before inspection and a progress/status update
after inspection. The runtime workflow should not imply project files are
mutated; tooling is the mutation surface.

Recommendation follow-up gains numbered apply-vs-handoff selection, a decision
brief for external issue creation, and a result closeout with `Next` and
boundary-sensitive `Not done`.

## Alternatives

- **Patch only runtime skill files.** Rejected because the durable specs would
  remain weaker than the executable skill and future runtime edits could regress.
- **Patch only the shared interaction contract.** Rejected because the failures
  occurred at concrete workflow prompt shapes; broad principles did not prevent
  them.
- **Rewrite the agent-mediated UX guide.** Rejected because the guide already
  states the relevant principles. The implementation needs to apply it locally.

## Trade-offs & risks

More explicit prompt templates make the skill slightly longer. That cost is
acceptable because these are user-facing workflow contracts, and the extra lines
remove ambiguity at decision and correction points.

The final setup gate changes from a casual confirmation prompt to a decision
brief. It should still preserve a fast confirmation path, but it will look more
formal because it authorizes writing `QUALITY.md`.

## Open questions

None.
