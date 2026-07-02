---
title: LedgerLite Service
description: Illustrative quality model for a fictional ledger and payments API.
ratingScale:
  - level: outstanding
    title: 🟢 Outstanding
    description: The service clearly exceeds the shared quality bar.
    criterion: Consistently exceeds the requirement with clear operational margin.
  - level: target
    title: 🔵 Target
    description: The service meets the shared quality bar.
    criterion: Meets the expected quality bar with evidence a maintainer can verify.
  - level: minimum
    title: 🟡 Minimum
    description: The service is usable, but quality gaps need attention.
    criterion: Meets the lowest acceptable bar, with visible gaps or limited evidence.
  - level: unacceptable
    title: 🔴 Unacceptable
    description: The service is below the shared quality bar.
    criterion: Falls below the lowest acceptable bar.
areas:
  api:
    title: Public API
    source: synthetic-source:api
    factors:
      correctness:
        title: Correctness
        description: Requests have clear semantics and preserve ledger intent.
      operability:
        title: Operability
        description: API behavior is understandable to callers and operators.
    requirements:
      idempotent-mutations:
        title: mutation endpoints are idempotent under retry
        factors: [correctness]
        assessment: >
          Inspect the mutation contract and retry tests for idempotency keys,
          replay behavior, and duplicate-write prevention.
      predictable-error-contracts:
        title: error responses are predictable for callers
        factors: [operability]
        assessment: >
          Compare documented error responses with handler behavior for common
          validation, authorization, and conflict cases.
  persistence:
    title: Ledger Persistence
    source: synthetic-source:persistence
    factors:
      integrity:
        title: Integrity
        description: Stored ledger state preserves accounting invariants.
      recoverability:
        title: Recoverability
        description: Data changes can be reversed or recovered when releases fail.
    requirements:
      balance-invariants:
        title: ledger mutations preserve balance invariants
        factors: [integrity]
        assessment: >
          Inspect mutation tests and reconciliation checks for conservation of
          balance across debits, credits, and failed writes.
      migration-rollback:
        title: migrations have rehearsed rollback paths
        factors: [recoverability]
        assessment: >
          Inspect migration runbooks and release notes for rollback instructions,
          rehearsal evidence, and known irreversible changes.
  operations:
    title: Operations
    source: synthetic-source:operations
    factors:
      observability:
        title: Observability
        description: Operators can understand health and customer impact quickly.
      recoverability:
        title: Recoverability
        description: Incidents have clear ownership and practiced recovery paths.
    requirements:
      customer-impact-telemetry:
        title: health signals explain customer impact
        factors: [observability]
        assessment: >
          Inspect dashboards and alerts for signals that connect service health to
          failed customer actions.
      recovery-drill-ownership:
        title: recovery drills have current owners
        factors: [recoverability]
        assessment: >
          Inspect the recovery calendar and incident playbooks for named owners,
          recency, and unresolved drill follow-up.
  agent-harness:
    title: Agent Harness
    source: synthetic-source:agent-harness
    factors:
      agent-accessibility:
        title: Agent Accessibility
        description: Agent-facing instructions expose context, checks, and limits.
    requirements:
      evaluation-entrypoint:
        title: agent guidance routes quality evaluation work
        factors: [agent-accessibility]
        assessment: >
          Inspect agent guidance for a stable entry point to the quality model,
          evaluation checks, generated reports, and permitted mutation boundaries.
---

# Quality model: LedgerLite Service

This illustrative model describes a fictional service that accepts ledger
mutation requests, persists accounting state, operates under production support,
and exposes an agent harness for quality-management work.

The example is intentionally software-specific so report design can exercise a
familiar, non-trivial service shape. It is not a default QUALITY.md template and
does not imply these factors are universal defaults.

The generated evaluation report uses synthetic routine outputs and synthetic
evidence references. The concrete source system is omitted.
