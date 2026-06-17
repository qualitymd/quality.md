---
type: Analysis Record
title: Sparrow Payments — Analysis: Webhooks (0001)
description: Reference Analysis Record — the write-once Analyze-phase record for the Webhooks node, whose aggregate (Minimum) differs from its local rating (Target) because a child subtree pulls it down.
tags: [skill, quality, evaluation, example]
timestamp: 2026-06-17T00:00:00Z
---

> **Reference instance — non-normative.** A captured example of a **write-once
> analysis record** for an internal target node whose **aggregate differs from its
> local rating** — the case where a child subtree binds. See
> [Reporting](../../../quality-skill.md#reporting) and the
> [examples index](../../index.md) for the facts shared across this bundle.

# Analysis: Webhooks *(child of root)*

**Aggregate: Minimum**
**Local: Target**

## Local rating

**Target** — over the Webhooks node's one own requirement.

Derived from the Webhooks assessment record:

- *Every outbound webhook is signed so merchants can verify its origin* —
  **Target** (record 006) — full signing coverage with origin and replay
  protection; per-merchant signing-secret rotation is not enforced, so short of
  Outstanding.

## Factors

- **Security — Target.** A refinement of the root Security factor for the
  outbound-delivery context: proving an event's origin to a merchant. Every emit
  path signs and no unsigned path was found.

## Aggregate rating

**Minimum** — the node's own work is *rated* **Target**, but its **Delivery** child
subtree is *rated* **Minimum** and pulls the aggregate below this node's local
rating. The binding constraint is the webhook-delivery deduplication gap
([008](../assessments/008-delivery-redelivery-dedup.md)); see recommendation
[003](../recommendations/003-bound-webhook-dedup-window.md).

Derived from this node's local rating and the child analysis record: Delivery
(Minimum).
