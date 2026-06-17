---
type: Assessment Record
title: Sparrow Payments — Assessment: redelivery deduplication (0001)
description: Reference Assessment Record — the write-once record behind the Delivery deduplication rating, the next binding constraint at Minimum.
tags: [skill, quality, evaluation, example, reliability]
timestamp: 2026-06-17T00:00:00Z
---

> **Reference instance — non-normative.** A captured example of a **write-once
> assessment record**, here the one the [report](../report.md) rolls up into the
> Delivery subtree's **Minimum** local rating. See
> [Reporting](../../../quality-skill.md#reporting) and the
> [examples index](../../index.md) for the facts shared across this bundle.

# Assessment: *A redelivery of an already-acknowledged event is suppressed*

**Target / factor:** Delivery (child of Webhooks) → Reliability
**Requirement:** *A redelivery of an already-acknowledged event is suppressed for
that endpoint*
**Rating:** **Minimum**

## Findings

- Each event carries a stable delivery id, and the engine suppresses a redelivery
  to an endpoint that already acknowledged it (`webhooks/delivery/dedup.go:52`).
- Suppression state is retained for only 24 hours — the same as the retry window —
  so a delivery re-enqueued late in that window can fire after its dedup record has
  expired and reach the merchant a second time.
- The merchant docs note "consume events idempotently", confirming duplicates are
  possible by design.

## Rationale

Deduplication exists and covers the common case but is bounded and best-effort, so
a duplicate is reachable on a known path — short of the target's "not delivered
again" intent, yet a consciously-bounded floor rather than a failure. Rated
**Minimum**. Once the committed credential (root) is cleared, this is the next
binding constraint on the whole-model rating. See recommendation
[003](../recommendations/003-bound-webhook-dedup-window.md).
