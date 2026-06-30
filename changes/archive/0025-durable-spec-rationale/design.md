---
type: Design Doc
title: Durable spec rationale - design doc
description: Why durable rationale lives in the spec as a two-layer, co-located convention rather than a separate doc or a full ADR.
tags: [specs, docs, contributing, design]
timestamp: 2026-06-18T00:00:00Z
---

# Durable spec rationale - design doc

Design behind the [Durable spec rationale](../0025-durable-spec-rationale.md)
change and its [functional spec](spec.md). The spec fixes _what_ the three
guides must require; this doc covers _why that shape_, and what was rejected.

## Context

The repo already separates the _delta_ (a Change Case and its archived spec/design)
from the _cumulative_ source of truth (the enduring `specs/` bundle). The
mechanism is sound for the _what_ but lossy for the _why_: motivation lives in
the Change Case, deeper rationale lives in the design doc, and both archive with
the case — so the durable spec keeps the rule and loses the reason. The
functional-spec guide compounds it by capping all rationale at a clause or
`Note:`, leaving genuine durable intent no home in the spec.

This is a guides-only change. Nothing in `specs/` is rewritten here; the
convention takes effect as specs are next edited.

## Approach

Two layers of durable rationale, both **co-located in the spec**:

1. **Background / Motivation** — one short prose section near the top of a spec,
   for the spec-scale _why_ (the problem or failure-mode the capability
   addresses). It sits alongside Scope and the companion note, and is explicitly
   not either of them.
2. **Per-requirement annotation** — an optional subordinate blockquote led by
   `Rationale:` directly beneath a requirement, one or two sentences, optionally
   carrying the originating change id. The requirement stays the lead, testable
   sentence; the annotation never wraps around it.

A **litmus** governs the fine grain: annotate only where a future editor would
otherwise repeat a mistake or be misled. Full decision records and dead-end
alternatives stay in the design doc. A **say-it-once** rule keeps the two layers
from restating each other, and authors supersede stale rationale rather than
letting it pile up.

The change-cases guide closes the loop: absorbing a landing case into the
durable specs now includes promoting its _why_ (the Change Case's motivation and
the design doc's durable rationale) into the spec's Background and annotations.
Because durable specs may be edited at any time — with or without a case — this
promotion is a **SHOULD** encouraged whenever a case updates a spec, not a gate
on the **Before setting In-Review** step. The design-doc guide records the same
promotion from its side, and keeps the design doc as the fuller archived record.

The spec for this change dogfoods all of the above, so it doubles as a worked
example the guides can point at.

## Alternatives

**A separate Diátaxis _explanation_ doc per capability.** Rejected. Diátaxis
would file the _why_ as explanation, away from the normative spec. That keeps the
spec lean but reopens the same gap this change exists to close: the rationale
sits in a different file from the requirement it justifies, so an editor changing
the rule need never see it. Co-location is the whole point.

**Design-intent depth only — promote a one-line intent, nothing more.** Rejected
as too thin. A bare intent line ("be strict about exit codes") does not carry the
failure-mode ("a zero exit on lint failure green-lit broken files"), which is
exactly the lesson that stops the rule from being undone. The litmus, not a
length cap, is the right control.

**A full ADR embedded in the spec.** Rejected as too heavy. An
Architecture-Decision-Record per requirement — context, options, decision,
consequences — would bury the requirement under its own justification, the very
smell the guide warns against, and duplicate the design doc. The design doc
already is the decision record; the spec promotes only durable intent and
lessons.

In-spec two-layer rationale wins because it puts the _why_ where the _what_
lives (unlike the explanation doc), carries enough to prevent regressions (unlike
intent-only), and stays subordinate to the rule (unlike an ADR).

## Trade-offs and risks

The main risk is **spec bloat** — rationale crowding out requirements until the
spec is no longer skimmable, which is the failure the old "asides only" rule was
guarding against. Three mitigations, all in the convention: the requirement stays
the lead testable sentence and the annotation is subordinate; the litmus admits
only regression-preventing context, not running commentary; and the say-it-once
rule plus supersede-don't-accrete keep Background and annotations from
duplicating each other or growing stale.

A secondary risk is **drift between the spec's promoted rationale and the
archived design doc**. This is acceptable by design: the spec carries durable
intent and lessons, the design doc carries the full point-in-time decision
record, and the design-doc guide already says to supersede a design rather than
rewrite it.

## Open questions

- Whether a future change should retrofit Background sections and annotations
  into the highest-traffic existing specs in one pass, rather than waiting for
  each to be next edited.
- Whether a lightweight check should ever flag a requirement whose originating
  change is known but which carries no annotation. For this change, the litmus is
  applied by authors, not tooling.
