---
type: Runtime Guide
title: Recommendation Follow-Up
description: Guide for applying or handing off evaluation recommendations.
---

# Recommendation Follow-Up

Use recommendation follow-up when the user asks to apply, act on, improve from,
or hand off an active evaluation recommendation. This is not a `/quality` mode:
it is a post-evaluation workflow over recommendation records that already exist
or have just been produced by `evaluate`.

## Outcomes

Offer only two explicit productive outcomes:

1. Apply a confirmed recommendation option now.
2. Hand off the recommendation to an issue tracker.

If the user does not choose one of those outcomes, stop without changing
evaluated source, `QUALITY.md`, `.quality/log/`, or external systems. Do not
present defer, skip, or keep open as formal options.

## Apply Now

Before editing anything, present a decision brief:

```text
**Apply <recommendation>?**

**Changes:** <evaluated source | QUALITY.md | both | quality log when model changes>
**Evidence/reason:**
**Recommended option:**
**Alternatives:** hand off to issue tracker
**Done criterion / verification:**
```

Do not treat an obvious recommendation as consent. Apply only the confirmed
recommendation option and mutation surface.

After applying, verify the done criterion with the narrowest useful evidence. If
the done criterion is rating-bound or depends on the QUALITY.md model, run a
scoped re-evaluation in a new numbered folder and report the before/after delta:

```text
**Recommendation result**

**Recommendation:**
**Outcome:** applied
**Applied option:**
**Changed artifacts:**        (name the quality log entry when the model changed)
**Verification:**
**Rating movement:**          (when known)
**Remaining gaps / limits:**
```

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
issue-tracker tooling. If tooling is unavailable or the user has not confirmed
creation, stop after producing the issue-ready text.
