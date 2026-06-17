---
type: Assessment Record
title: Sparrow Payments — Assessment: daily reconciliation (0001, not assessed)
description: Reference Assessment Record — the write-once record for a requirement recorded as not assessed for want of evidence.
tags: [skill, quality, evaluation, example]
timestamp: 2026-06-17T00:00:00Z
---

> **Reference instance — non-normative.** A captured example of a **write-once
> assessment record** for a requirement recorded as *not assessed* — kept distinct
> from a rated outcome, never assigned a level to fill the gap. See
> [Reporting](../../../quality-skill.md#reporting) and the
> [examples index](../../index.md) for the facts shared across this bundle.

# Assessment: *Reconciliation runs daily and flags drift*

**Target / factor:** Ledger → Correctness
**Requirement:** *Reconciliation runs daily and flags drift*
**Rating:** **Not assessed**

## Findings

- None — no reconciliation job output, log, or report was available to assess
  against.
- A `reconcile` entrypoint exists (`ledger/reconcile.go:31`), but no evidence of a
  scheduled run was found.

## Rationale

Insufficient evidence to rate against the scale, so recorded as *not assessed*
rather than assigned a level (per the *not assessed*, not a guess rule of
[Grounding judgment](../../../quality-skill.md#grounding-judgment)). The Ledger's
local rating therefore rests on its one assessed requirement and is noted as
incomplete. See recommendation
[002](../recommendations/002-produce-reconciliation-evidence.md).
