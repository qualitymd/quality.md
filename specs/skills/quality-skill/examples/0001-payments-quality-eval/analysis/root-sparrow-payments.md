---
type: Analysis Record
title: Sparrow Payments — Analysis: root (0001)
description: Reference Analysis Record — the write-once Analyze-phase record (the inferred roll-up) for the root target, held at Unacceptable by its own local rating.
tags: [skill, quality, evaluation, example]
timestamp: 2026-06-17T00:00:00Z
---

> **Reference instance — non-normative.** A captured example of one **write-once
> analysis record** the [`/quality` skill](../../../quality-skill.md) produces
> during the Analyze phase (see [Reporting](../../../quality-skill.md#reporting)).
> One record per target node, holding the inferred roll-up; it cites the
> assessment records behind its local rating and its children's analysis records
> behind its aggregate, so the chain leaf → node → root is explicit. The
> [report](../report.md) renders over these. See the
> [examples index](../../index.md) for the facts shared across this bundle.

# Analysis: Sparrow Payments API *(root)*

**Aggregate: Unacceptable**
**Local: Unacceptable**

## Local rating

**Unacceptable** — over the root's three own requirements. The committed-credential
shortfall is security-critical, and a single such finding holds the local rating
at the floor regardless of the two requirements *rated* **Target**.

Derived from the root assessment records:

- *No credentials are committed to the repository* — **Unacceptable**
  ([001](../assessments/001-root-no-committed-credentials.md)) — the binding
  constraint.
- *Every money-moving endpoint enforces authentication* — **Target** (record 002).
- *Transfers are idempotent on retry* — **Target** (record 003).

## Factors

- **Security — Unacceptable.** Bound by the **Secrets handling** sub-factor (the
  committed live credential); **Access control** is *rated* **Target**.
- **Reliability — Target.** Transfer idempotency holds, and the Ledger's
  double-entry invariant — which tags Reliability as a secondary factor — adds
  evidence; this lens ranges into the Ledger subtree.

## Aggregate rating

**Unacceptable** — the root's own local rating binds. Neither child subtree
changes it: the **Ledger** analysis is *rated* **Target** and the **Webhooks**
analysis ([webhooks](webhooks.md)) is *rated* **Minimum**. Were the local rating
cleared, the Webhooks subtree would be the next binding constraint and the
aggregate would rise only to **Minimum** — reaching **Target** once the
webhook-delivery deduplication gap is closed too.

Derived from this node's local rating and the child analysis records: Ledger
(Target) and [Webhooks](webhooks.md) (Minimum).
