---
type: Assessment Record
title: Sparrow Payments — Assessment: no committed credentials (0001)
description: Reference Assessment Record — the write-once assess→finding→rating record behind the root committed-credential rating, held at Unacceptable.
tags: [skill, quality, evaluation, example, security]
timestamp: 2026-06-17T00:00:00Z
---

> **Reference instance — non-normative.** A captured example of one **write-once
> assessment record** the [`/quality` skill](../../../quality-skill.md) produces
> during the Evaluate phase (see [Reporting](../../../quality-skill.md#reporting)).
> One record per in-scope requirement; the [report](../report.md) is the roll-up
> rendered over them. Written atomically and never mutated — a re-assessment
> produces a new evaluation folder. See the [examples index](../../index.md) for
> the facts shared across this bundle.

# Assessment: *No credentials are committed to the repository*

**Target / factor:** Sparrow Payments API (root) → Security → Secrets handling
**Requirement:** *No credentials are committed to the repository*
**Rating:** **Unacceptable**

## Findings

Credential referenced by location and type only; the value is not reproduced here.

- `internal/gateway/client.go:48` — a live payment-gateway **API secret key**
  (credential type: gateway secret key) committed in plaintext. Matches the format
  of an active key; not a placeholder.
- `internal/gateway/client_test.go:12` — a **test publishable key**, non-secret by
  design. Not counted against this requirement.

## Rationale

Rated against this requirement's `ratings` overrides, not the scale's default
criteria — Unacceptable: "a live credential is present in the working tree";
Minimum and Target both require a clean tree (Target also requires any prior
exposure rotated). The one committed live key satisfies the Unacceptable
criterion, and a single live secret lands here regardless of the clean findings
elsewhere. Reported by `file:line` and credential type only. See recommendation
[001](../recommendations/001-rotate-committed-gateway-key.md) — rotation is
required, not merely removal, because the key is already exposed in history.
