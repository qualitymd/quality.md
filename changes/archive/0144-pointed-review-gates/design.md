---
type: Design Doc
title: Pointed Review Gates - design
description: How shared UX guidance and direct QUALITY.md authoring adopt purpose-first, assumption-focused review gates.
tags: [docs, skill, ux, authoring]
timestamp: 2026-06-27T00:00:00Z
---

# Pointed Review Gates - design

## Context

Answers the [functional spec](spec.md) for change case
[0144](../0144-pointed-review-gates.md). The change builds on
[0139](../0139-real-review-gates.md) and
[0140](../0140-casual-review-gate-wording.md): feedback invitations already
wait, and review gates already state the planned change and value prop. This case
sharpens what the gate asks the user to react to.

## Approach

### Make inferred purpose explicit

Update the shared UX guide and `/quality` direct-authoring guidance so a
content/model review gate says not only what the agent plans to change, but why
the agent thinks the change is needed. This keeps the user's correction path
about intent, not only mechanics.

### Ask for reaction to one consequential assumption

Replace the Security example's broad "what should I adjust or watch out for?"
ending with a narrower scope choice: broad Security versus narrower appsec. The
runtime skill template uses the same shape:

```text
One <scope/risk/naming> choice before I edit: <consequential assumption>. If
<alternative meaning>, say so. Otherwise say `go`.
```

This is still a review gate, not a binary decision brief. The user can answer
with the short approval path or provide a correction in free text.

### Keep catch-all prompts as fallback

The guide does not ban broad feedback invitations. It demotes them to a fallback
when the agent cannot infer a narrower steering axis. That preserves flexibility
for genuinely open-ended edits without letting routine model edits default to
weak "anything else?" wording.

### Align durable specs and logs

Update the durable `/quality` skill spec and durable authoring guide spec to make
the purpose and consequential-assumption fields part of the direct-authoring
contract. Update runtime and durable logs plus release notes so the change is
traceable after the Change Case archives.

### Record applicable guide consultation

Keep the already-started `work-with-change-cases.md` improvement in scope. It
adds a process rule requiring Change Case authors to identify and read applicable
guides before phase work and status advances, which supports this case's own
cross-guide alignment.

## Spec response

- **Shared UX guidance** - satisfied by updating the Review Gates section,
  Security example, and checklist.
- **Direct QUALITY.md authoring** - satisfied by updating `SKILL.md` and the
  durable skill spec.
- **Authoring guide alignment** - satisfied by updating the runtime authoring
  guide and durable authoring guide spec.
- **Change Case guidance** - satisfied by adding applicable-guide consultation to
  `work-with-change-cases.md`.
- **Verification** - satisfied by source inspection and Markdown formatting.

## Alternatives

- **Keep "what should I adjust?" and add examples.** Rejected. It improves the
  wording but still leaves the user to infer which assumption matters most.
- **Make every review gate a closed choice.** Rejected. A pointed review gate is
  still conversational; a forced picker would be too rigid when the user wants to
  correct the premise.
- **Only update the shared UX guide.** Rejected. The `/quality` skill prompt and
  durable skill specs are the executable contract agents follow.

## Trade-offs & risks

- A pointed assumption can be wrong. That is acceptable because the gate exists
  to expose wrong inference before mutation, and the response path explicitly
  invites correction.
- Too much rationale can bloat the gate. The example stays concise and asks for
  one scope decision rather than an exhaustive questionnaire.

## Open questions

None.
