---
type: Recommendation
title: Bound the webhook redelivery-deduplication window
description: Reference Recommendation — close the bounded-deduplication gap that becomes the binding constraint on Sparrow Payments once the committed credential is rotated.
tags: [skill, quality, evaluation, example, reliability]
timestamp: 2026-06-17T00:00:00Z
---

> **Reference instance — non-normative.** A captured example of a single
> recommendation artifact the [`/quality` skill](../../../quality-skill.md)
> emits alongside its [report](../report.md). It is written to stand on its own:
> a reader can triage and route it without the report or the session.

# Bound the webhook redelivery-deduplication window

**Target / factor:** Webhooks → Delivery → Reliability
**In-scope requirement:** *A redelivery of an already-acknowledged event is suppressed for that endpoint*
**Current rating:** Minimum — the **next binding constraint**: it holds the Webhooks subtree at Minimum and, once the committed credential is rotated, holds the whole-model rating below Target.

## Gap

The delivery engine deduplicates redeliveries by event id, but retains
suppression state for only 24 hours — the same bound as the retry window. A
delivery re-enqueued late in that window can fire after its dedup record has
expired, so a merchant can observe a duplicate event. This meets the Minimum
criterion (acceptable, consciously bounded) rather than Target ("not delivered
again to the same endpoint").

**Evidence:**

- `webhooks/delivery/dedup.go:52` — suppression keyed by delivery id with a
  fixed 24-hour retention; a redelivery after the window is not suppressed.
- `webhooks/delivery/retry.go:37` — the retry window is also 24 hours, so a
  delivery re-enqueued near its end can fire after the dedup record has expired.
- Merchant docs note "consume events idempotently", confirming duplicates are
  possible by design rather than prevented.

## Options

- **(a) Tie the dedup-retention window to the retry window plus a margin.**
  Retain suppression state for at least the maximum time a retry can fire, so no
  in-flight retry can outlive its dedup record; the engine then never redelivers
  an acknowledged event on the retry path.
- **(b) Make suppression durable and unbounded per event id.** Persist
  acknowledged-delivery ids so a redelivery is always suppressed, at the cost of
  unbounded dedup storage to manage.
- **(c) Leave dedup best-effort and document at-least-once.** Keep the bounded
  window and rely on merchants consuming idempotently — the current Minimum
  posture, made explicit rather than improved.

## Recommended

**(a) Tie the dedup-retention window to the retry window plus a margin.** It
directly closes the reachable duplicate path — no acknowledged event is
redelivered while a retry could still fire — without the unbounded storage of
(b). Option (c) only documents the gap, so it does not move the rating off
Minimum.

## Done-criterion

The requirement *A redelivery of an already-acknowledged event is suppressed for
that endpoint* reaches **Target** against its criterion — an acknowledged event
is not delivered again to the same endpoint for the full duration any retry
could fire. A later `improve` re-evaluates the Delivery scope to confirm the
rating moved off Minimum; clearing this gap (after the committed credential is
rotated) is expected to lift the Webhooks subtree and the whole-model rating to
**Target**.
