---
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
    title: Security
    description: >
      Security is the degree to which Sparrow Payments keeps money movement and
      account data out of unauthorized hands; it matters because a leaked
      credential or an unauthenticated transfer endpoint lets an attacker move
      customer money. It is the "who may act and what may be exposed" lens,
      distinct from Reliability's "does a legitimate request behave correctly".
    factors:
      secrets handling:
        title: Secrets handling
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
              credentials. Classify each match as a live secret or a non-secret
              (test fixture, publishable key, placeholder). The requirement
              fails if any live credential is present in the working tree.
              Reference every match by file:line and credential type only; never
              reproduce the value, and recommend rotation for any live secret
              found, since committed secrets persist in history.
            ratings:
              outstanding: "No credentials in the tree; secrets load from a secret manager and a scan is enforced in CI."
              target: "No live credential is present in the working tree, and any previously exposed credential has been rotated."
              minimum: "No live credential is present in the working tree, but a previously committed secret has not yet been rotated."
              unacceptable: "A live credential is present in the working tree (e.g. committed in plaintext)."
      access control:
        title: Access control
        description: >
          Access control is the degree to which only authenticated, authorized
          callers can invoke money-moving operations; it matters because an
          unauthenticated transfer or refund endpoint lets anyone move funds.
          Distinct from Secrets handling, which concerns stored credential
          material rather than request-time authorization.
        requirements:
          "Every money-moving endpoint enforces authentication":
            assessment: >
              Enumerate the registered HTTP routes and identify the money-moving
              ones (transfers, refunds, payouts). Confirm each resolves through
              the authentication middleware before its handler runs. The
              requirement is met when no money-moving route is reachable without
              authentication; report any exceptions by route and location.
  reliability:
    title: Reliability
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
          required idempotency key) and confirm, by test or trace, that replaying
          the same request returns the original result without a second debit.
          The requirement is met when a retried transfer produces no duplicate
          money movement.
areas:
  ledger:
    title: Ledger
    description: >
      The double-entry accounting ledger under `./ledger` — the system of record
      that books every customer money movement, evaluated as its own subtree
      distinct from the API service that writes to it.
    source: ./ledger
    factors:
      correctness:
        title: Correctness
        description: >
          Correctness is the degree to which the ledger records every transfer as
          a balanced, reconcilable set of entries; it matters because a ledger
          that drifts from actual money movement cannot be trusted for settlement
          or audit.
        requirements:
          "Every transfer debits and credits to a net zero (double-entry invariant)":
            assessment: >
              Inspect the posting routine for enforcement that a transfer's
              entries sum to zero before commit, and confirm by test that no
              imbalanced posting is accepted. The requirement is met when the
              double-entry invariant holds on every observed path.
            factors:
              - reliability
          "Reconciliation runs daily and flags drift":
            assessment: >
              Confirm a reconciliation job runs on a daily schedule and emits a
              durable report showing whether ledger balances match the gateway
              and flagging any drift within the expected window. Assess against
              that report's output; if no run, log, or report is available,
              record the requirement as not assessed rather than assigning a
              level.
  webhooks:
    title: Webhooks
    description: >
      The outbound webhook subsystem under `./webhooks` — the service that
      notifies merchants of events affecting them (a completed transfer, a
      settled payout) through signed HTTP callbacks, evaluated as its own subtree
      distinct from the API that emits the events and the ledger that records
      them.
    source: ./webhooks
    factors:
      security:
        title: Security
        description: >
          Security here is the degree to which a merchant can trust that a
          received webhook genuinely originated from Sparrow and was not forged;
          it matters because a merchant that acts on a spoofed event ships goods
          or releases funds against money that never moved. A refinement of the
          root Security factor for the outbound-delivery context — the concern is
          authenticating Sparrow to the merchant, not authenticating callers to
          Sparrow.
        requirements:
          "Every outbound webhook is signed so merchants can verify its origin":
            assessment: >
              Confirm every outbound webhook path signs the payload with the
              recipient merchant's signing secret before send, and that the
              signature covers the body and a timestamp so a merchant can verify
              origin and reject replays. The requirement is met when no code path
              emits an unsigned webhook; report any unsigned path by location.
    areas:
      delivery:
        title: Delivery
        description: >
          The delivery engine under `./webhooks/delivery` — the queue, retry, and
          acknowledgement machinery that gets a signed webhook to the merchant's
          endpoint, evaluated as its own subtree distinct from the signing
          concern on the parent target.
        source: ./webhooks/delivery
        factors:
          reliability:
            title: Reliability
            description: >
              Reliability here is the degree to which an event that should reach
              a merchant does reach it, exactly once in effect; it matters
              because a silently dropped event leaves a merchant's view of money
              movement stale, and a duplicate event can trigger a double
              fulfillment. A refinement of the root Reliability factor for the
              delivery context — retry on failure and suppression of duplicate
              redeliveries.
            requirements:
              "Failed deliveries retry with exponential backoff until acknowledged or the retry window expires":
                assessment: >
                  Confirm a failed delivery is retried on an exponential backoff
                  schedule until the endpoint acknowledges or a bounded retry
                  window expires, and that an exhausted delivery is recorded as
                  failed rather than dropped silently. The requirement is met
                  when no transient failure results in a silently lost event.
              "A redelivery of an already-acknowledged event is suppressed for that endpoint":
                assessment: >
                  Confirm each event carries a stable delivery id and that the
                  engine suppresses a redelivery to an endpoint that has already
                  acknowledged that id, so retries do not surface as duplicate
                  events at the merchant. The requirement is met when an
                  acknowledged event is not delivered again to the same endpoint;
                  note the bound on any dedup-retention window within which this
                  holds.
---

# Sparrow Payments API — Quality model

## Overview

Sparrow Payments moves money between customer accounts and an external payment
gateway, records every movement in a double-entry ledger, and notifies merchants
of those movements through signed webhooks. Good means money moves only on
authenticated requests, never duplicates on retry, the ledger stays balanced and
reconcilable, and merchants are notified verifiably and exactly once.

## Scope

This model covers the API service (`./`), the ledger (`./ledger`), and the
outbound webhook subsystem (`./webhooks`). The external payment gateway and the
banks behind it are out of scope — Sparrow does not own them — and neither are
the merchant endpoints that receive webhooks.

## Needs

- Customers' funds move only on authenticated, authorized requests.
- A retried transfer never moves money twice.
- The ledger stays balanced and can be reconciled against actual settlement.
- Credentials never leak through the source repository.
- Merchants are notified of events affecting them, verifiably and exactly once.

## Risks

A leaked gateway credential or an unauthenticated money-moving endpoint is the
worst outcome: either lets an attacker move customer funds. A double-debit on
retry, or a ledger that silently drifts, is unrecoverable without manual
investigation. A spoofed or unsigned webhook can make a merchant act on money
that never moved; a silently dropped notification leaves the merchant's view
stale.

*Unknowns* — sustained peak-load behavior is in scope but not yet modeled.
*Open questions* — none.

*Reviewed — Margaret Hamilton, 2026-05; agent-reviewed — Claude, 2026-06.*
