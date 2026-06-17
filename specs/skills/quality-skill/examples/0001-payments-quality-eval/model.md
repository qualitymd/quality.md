---
type: Example Model
title: Sparrow Payments — the QUALITY.md evaluated (0001)
description: Reference QUALITY.md — the fictional model the 0001 Evaluation Report evaluates, reproduced so the report's evidence is traceable.
tags: [skill, quality, evaluation, example, model]
timestamp: 2026-06-17T00:00:00Z
---

> **Reference instance — non-normative.** This reproduces the `QUALITY.md` the
> [0001 report](report.md) evaluates, so its findings trace to declared
> requirements and `source` selectors. At runtime this file lives at the
> **evaluated** repository root as `./QUALITY.md` — *not* inside the skill's
> `quality/evaluations/…` output folder; it is bundled here only so the example
> is self-contained. The subject ("Sparrow Payments") and its `file:line`
> locators are fictional. The model uses the suggested four-level scale.

# The model evaluated

A whole view of the `QUALITY.md` the report assesses. Frontmatter and body below
together form the single file `./QUALITY.md`.

## Frontmatter

```yaml
title: Sparrow Payments API
ratingScale:
  - level: outstanding
    title: Outstanding
    description: "The stretch band — reached only with significant extra effort."
    criterion: "Exceeds the requirement; satisfies it with margin to spare."
  - level: target
    title: Target
    description: "The level to aim for — achievable at reasonable cost and effort."
    criterion: "Satisfies the requirement."
  - level: minimum
    title: Minimum
    description: "The acceptable floor — less than you'd aim for, but consciously agreed as good enough to ship."
    criterion: "Falls short of the target but remains acceptable."
  - level: unacceptable
    title: Unacceptable
    description: "Below the floor — not good enough to ship."
    criterion: "Does not meet the requirement to an acceptable degree."

factors:
  security:
    description: >
      Security is the degree to which Sparrow Payments keeps money movement and
      account data out of unauthorized hands; it matters because a leaked
      credential or an unauthenticated transfer endpoint lets an attacker move
      customer money. It is the "who may act and what may be exposed" lens,
      distinct from Reliability's "does a legitimate request behave correctly".
    factors:
      secrets handling:
        description: >
          Secrets handling is the degree to which credentials — gateway keys,
          tokens, database passwords — are kept out of source and managed so an
          exposed secret can be contained; it matters because a committed live
          key is an immediate, history-persistent exposure. Distinct from Access
          control, which governs request-time authorization, not stored secret
          material.
        requirements:
          "No credentials are committed to the repository":
            assessment: >
              Scan the tracked source and configuration for committed
              credentials. Classify each match as a live secret or a
              non-secret (test fixture, publishable key, placeholder). The
              requirement fails if any live credential is present in the working
              tree. Reference every match by file:line and credential type
              only; never reproduce the value, and recommend rotation for any
              live secret found, since committed secrets persist in history.
            ratings:
              outstanding: "No credentials in the tree; secrets load from a secret manager and a scan is enforced in CI."
              target: "No live credential is present in the working tree, and any previously exposed credential has been rotated."
              minimum: "No live credential is present in the working tree, but a previously committed secret has not yet been rotated."
              unacceptable: "A live credential is present in the working tree (e.g. committed in plaintext)."
      access control:
        description: >
          Access control is the degree to which only authenticated, authorized
          callers can invoke money-moving operations; it matters because an
          unauthenticated transfer or refund endpoint lets anyone move funds.
          Distinct from Secrets handling, which concerns stored credential
          material rather than request-time authorization.
        requirements:
          "Every money-moving endpoint enforces authentication":
            assessment: >
              Enumerate the registered HTTP routes and identify the
              money-moving ones (transfers, refunds, payouts). Confirm each
              resolves through the authentication middleware before its handler
              runs. The requirement is met when no money-moving route is
              reachable without authentication; report any exceptions by route
              and location.
  reliability:
    description: >
      Reliability is the degree to which Sparrow Payments performs money
      movements correctly under retry and partial failure, without losing or
      duplicating effects; it matters because a retried transfer that
      double-debits is unrecoverable for the customer. It is the "does a
      legitimate request behave correctly" lens, distinct from Security's "who
      may act".
    requirements:
      "Transfers are idempotent on retry":
        assessment: >
          Inspect the transfer entrypoint for an idempotency mechanism (e.g. a
          required idempotency key) and confirm, by test or trace, that
          replaying the same request returns the original result without a
          second debit. The requirement is met when a retried transfer produces
          no duplicate money movement.

targets:
  ledger:
    source: ./ledger
    factors:
      correctness:
        description: >
          Correctness is the degree to which the ledger records every transfer
          as a balanced, reconcilable set of entries; it matters because a
          ledger that drifts from actual money movement cannot be trusted for
          settlement or audit.
        requirements:
          "Every transfer debits and credits to a net zero (double-entry invariant)":
            assessment: >
              Inspect the posting routine for enforcement that a transfer's
              entries sum to zero before commit, and confirm by test that no
              imbalanced posting is accepted. The requirement is met when the
              double-entry invariant holds on every observed path.
          "Reconciliation runs daily and flags drift":
            assessment: >
              Confirm a reconciliation job runs on a daily schedule and emits a
              durable report showing whether ledger balances match the gateway
              and flagging any drift within the expected window. Assess against
              that report's output; if no run, log, or report is available,
              record the requirement as not assessed rather than assigning a
              level.
```

## Body

```markdown
# Sparrow Payments API — Quality model

## Overview

Sparrow Payments moves money between customer accounts and an external payment
gateway, and records every movement in a double-entry ledger. Good means money
moves only on authenticated requests, never duplicates on retry, and the ledger
stays balanced and reconcilable.

## Scope

This model covers the API service (`./`) and the ledger (`./ledger`). The
external payment gateway and the banks behind it are out of scope — Sparrow does
not own them.

## Needs

- Customers' funds move only on authenticated, authorized requests.
- A retried transfer never moves money twice.
- The ledger stays balanced and can be reconciled against actual settlement.
- Credentials never leak through the source repository.

## Risks

A leaked gateway credential or an unauthenticated money-moving endpoint is the
worst outcome: either lets an attacker move customer funds. A double-debit on
retry, or a ledger that silently drifts, is unrecoverable without manual
investigation.

## Known gaps

- Sustained peak-load behavior is in scope but not yet modeled.
```
