---
type: Runtime Guide
title: Recommendation Follow-Up
description: Guide for applying or handing off evaluation recommendations.
---

# Recommendation Follow-Up

Use recommendation follow-up when the user asks to apply, act on, improve from,
or hand off an active evaluation recommendation. This is not a public `/quality`
workflow: it is a post-evaluation follow-up over recommendation records that
already exist or have just been produced by `evaluate`.

## Opening

Emit the follow-up frame as the first user-visible output, before recommendation
inspection, history inspection, outcome selection, issue-ready drafting, local
apply, issue creation, quality-log writes, or any other tool-dependent work. Use
`resolving…` for fields that need inspection:

```text
**QUALITY.md · recommendation follow-up**
- **Recommendation:** <id/title | resolving…>
- **Outcome:** <apply locally | hand off to issue tracker | resolving…>
- **Mutation:** <evaluated source | QUALITY.md | quality log | external issue | read-only until confirmed>
- **Artifacts:** <changed files/log entry/issue-ready text/issue link | none yet>
- **Next gate:** recommendation selection, outcome choice, decision brief, or verification
```

This frame does not introduce a new public invocation. It is the status-first
opening for the selected follow-up path, and all mutation still waits for the
confirmation gates below.

## Outcomes

Offer only two explicit productive outcomes:

1. Apply a confirmed recommendation option now.
2. Hand off the recommendation to an issue tracker.

If the user has not already chosen one, this is a single-select closed-choice
intent: render it through an option picker when one is fit-for-purpose, otherwise
use the numbered text fallback below (recommended option first). Keep the teaching
in the message either way.

```text
**What should happen with this recommendation?**

1. Apply the confirmed recommendation option now. **Recommended** when the option
   is clear and safe to perform locally.
2. Hand off the recommendation to an issue tracker.

**Answer:** Reply `1` or `2`.
```

Only mark option `1` as recommended when recommendation evidence supports local
apply; otherwise list the two options without inventing a recommendation.

If the user does not choose one of those outcomes, stop without changing
evaluated source, `QUALITY.md`, `.quality/log/`, or external systems. Do not
present defer, skip, or keep open as formal options.

## Apply Now

Before editing anything, present a decision brief:

```text
**Apply <recommendation>?**
<evaluated source | QUALITY.md | both | quality log when model changes>

  [y] Apply the confirmed recommendation now  — recommended
  [n] Stop, or hand off to an issue tracker instead

Reason: <evidence the recommendation rests on>
Done when: <verification of the done criterion>
```

Do not treat an obvious recommendation as consent. Apply only the confirmed
recommendation option and mutation surface.

After applying, verify the done criterion with the narrowest useful evidence. If
the done criterion is rating-bound or depends on the QUALITY.md model, run a
scoped re-evaluation in a new numbered folder and report the before/after delta:

```text
**Recommendation result:** applied — <recommendation>

**Applied option:**
**Changed artifacts:**        (name the quality log entry when the model changed)
**Verification:**

Rating movement: <when known>
Remaining gaps / limits:
Not done: <when a boundary matters>
Next:
```

Lead with the outcome line; keep the strongly-emphasized labels to the few the
user acts on (applied option, changed artifacts, verification) and let the rest
read as plain `label:` detail.

If verification is incomplete, label the result as limited rather than fully
confirmed.

When a confirmed apply changes the QUALITY.md model, append one quality log
entry under `.quality/log/` for the coherent model change. Cross-link the source
evaluation run and recommendation when present. Evaluated-source fixes that do
not change the model get no quality log entry.

Before applying a model-changing recommendation, read
[`authoring.md`](authoring.md), the routed authoring sub-guide for the model
element being changed, and [`authoring/quality-log.md`](authoring/quality-log.md).

## Issue-Tracker Handoff

Issue handoff does not edit evaluated source, `QUALITY.md`, or `.quality/log/`.

Prepare issue-ready text with:

- recommendation ID and title;
- source evaluation run;
- affected area, factor, or requirement;
- current rating when known;
- target or done criterion;
- evidence summary with locators;
- suggested implementation option;
- verification path;
- risk or priority notes when useful;
- links or paths to the generated report and recommendation artifact.

Creating an external issue requires explicit user confirmation and available
issue-tracker tooling. Use a decision brief before creating it. Where the
issue-tracker tooling will itself prompt to authorize the creation, render the
confirmation through that native gate and keep the brief's teaching in the
preceding message rather than stacking a second text `y`/`n` gate; this never
removes the confirmation requirement, only the redundant gate:

```text
**Create external issue?**
Create <issue tracker item> from the prepared issue-ready text.

  [y] Create the issue now  — recommended when tooling and confirmation are available
  [n] Stop after issue-ready text, or apply locally instead

Reason: recommendation evidence and issue-ready handoff text are prepared.
Done when: issue link or ID is returned.
Not changed: evaluated source, `QUALITY.md`, and `.quality/log/`.
```

If tooling is unavailable or the user has not confirmed creation, stop after
producing the issue-ready text.
