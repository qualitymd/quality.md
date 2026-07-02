---
title: LedgerLite Service
description: >
  Illustrative quality model for LedgerLite, a fictional ledger and payments
  service: its public API, service contract, persistence, operations, agent
  harness, and this model itself.
ratingScale:
  - level: outstanding
    title: 🟢 Outstanding
    description: "The stretch band: the requirement is exceeded with meaningful margin."
    criterion: Exceeds the requirement with margin a maintainer can verify.
  - level: target
    title: 🔵 Target
    description: "The expected good state: the requirement is satisfied."
    criterion: Satisfies the requirement with evidence a maintainer can verify.
  - level: minimum
    title: 🟡 Minimum
    description: "The acceptable floor: short of the goal but still safe to rely on."
    criterion: Falls short of target with visible gaps or limited evidence, while remaining acceptable.
  - level: unacceptable
    title: 🔴 Unacceptable
    description: "Below the floor: not good enough to rely on."
    criterion: Does not meet the requirement to an acceptable degree.
factors:
  # Agent harnessability rates how well the whole project equips an agent to
  # work; the agent-harness area below rates the checked-in harness artifacts
  # themselves. Keep the two projections distinct.
  agent-harnessability:
    title: Agent Harnessability
    description: >
      The degree to which LedgerLite's checked-in materials, tools, sensors,
      standards, and action limits equip an AI agent to understand the service,
      take scoped work, operate the environment, preserve state, verify its
      output, and stay safely bounded — distinct from the agent-harness area,
      which rates the harness artifacts' own quality.
    factors:
      agent-accessibility:
        title: Agent Accessibility
        description: >
          Decision-relevant service knowledge is reachable, selective, and
          intelligible to an agent working in context.
        requirements:
          agents-reach-service-context-from-a-stable-entry-point:
            title: a fresh agent reaches service context and deeper guidance from a stable entry point
            assessment: >
              Inspect the agent entry point and its routed guidance
              (synthetic-source:agent-harness) to confirm an agent can find the
              service purpose, the service contract, operational runbooks, and
              the recorded sensors without private context.
      task-specifiability:
        title: Task Specifiability
        description: >
          Work can be handed to an agent as a scoped assignment with explicit
          success and done criteria.
        requirements:
          quality-loop-work-items-carry-done-criteria:
            title: quality-loop work items carry scoped goals and done criteria
            assessment: >
              Sample the most recent quality-loop handoffs in the team tracker
              for stated goals, non-goals, done criteria, and the sensor that
              confirms completion.
      agent-operability:
        title: Agent Operability
        description: >
          A fresh agent session can establish and operate the working
          environment from recorded materials.
        requirements:
          a-fresh-session-reaches-a-ready-to-work-environment:
            title: a fresh agent session reaches a ready-to-work environment from recorded setup
            assessment: >
              Follow the recorded setup steps from a clean checkout and confirm
              the service builds, the sensors run, and required access is
              documented in agent-accessible materials.
      continuity:
        title: Continuity
        description: >
          Agent work preserves state and resumes across interruption, handoff,
          and fresh sessions.
        requirements:
          handoffs-survive-session-loss:
            title: in-flight work survives session loss through durable handoff records
            assessment: >
              Inspect recent long-running work for durable progress records
              that capture decisions, remaining work, and verification status
              without depending on chat history.
      self-verifiability:
        title: Self-Verifiability
        description: >
          Recorded sensors let an agent confirm on demand whether its own work
          is correct.
        requirements:
          sensors-return-pass-fail-with-remediation:
            title: recorded sensors return objective pass/fail with remediation-bearing output
            assessment: >
              Run the recorded sensor commands — contract tests, ledger
              invariant suite, and lint — and confirm each returns a
              deterministic result whose failures name the violated expectation
              and point toward the fix.
      enforcement-of-standards:
        title: Enforcement of Standards
        description: >
          Stated standards hold through gates rather than advisory prose.
        requirements:
          standards-gate-nonconforming-changes:
            title: core service standards are enforced by merge gates, not advisory prose
            assessment: >
              Inspect the merge pipeline configuration to confirm the contract
              tests, invariant suite, and lint sensors block nonconforming
              changes or route them through reviewable exceptions.
      containment-of-action:
        title: Containment of Action
        description: >
          Agent-permitted actions are confined by enforced limits and approval
          gates.
        requirements:
          consequential-actions-require-approval:
            title: money-moving and schema-changing actions require human approval
            assessment: >
              Inspect sandbox policy, permission allowlists, and deploy
              configuration to confirm an unattended agent run cannot move
              money, alter production schemas, or widen its own permissions
              without an approval gate.
areas:
  api:
    title: Public API
    description: >
      The HTTP surface integrators call to record and query ledger activity.
    source: synthetic-source:api
    factors:
      correctness:
        title: Correctness
        description: >
          The degree to which requests preserve ledger intent under retries,
          concurrency, and partial failure; integrators depend on it to trust
          recorded money movement.
      operability:
        title: Operability
        description: >
          The degree to which API behavior is understandable and diagnosable
          for callers and operators when requests fail.
      performance:
        title: Performance
        description: >
          The degree to which the API stays within its latency budget under
          representative production load; checkout flows time out beyond it.
    requirements:
      idempotent-mutations:
        title: mutation endpoints are idempotent under retry
        factors: [correctness]
        assessment: >
          Run the recorded contract-test suite for mutation endpoints and
          compare its replay cases with the retry and idempotency section of
          the service contract (synthetic-source:service-contract), including
          duplicate-key and partial-write replay behavior.
      predictable-error-contracts:
        title: error responses are predictable for callers
        factors: [operability]
        assessment: >
          Compare the error-envelope section of the service contract
          (synthetic-source:service-contract) with handler behavior for
          validation, authorization, and conflict cases across the endpoint
          index.
      p99-latency-within-budget:
        title: p99 mutation latency stays within budget
        factors: [performance]
        assessment: >
          Query the latency telemetry rollup for mutation endpoints over the
          most recent representative four-week production window and read p99
          against the rating bands.
        ratings:
          outstanding: p99 at or under 200 ms over the window.
          target: p99 at or under 300 ms over the window.
          minimum: p99 at or under 500 ms over the window.
          unacceptable: p99 above 500 ms over the window.
  service-contract:
    title: Service contract
    description: >
      The normative API contract that defines endpoint, retry, idempotency,
      and error semantics; other areas are judged against it, and this area
      judges the contract itself.
    source: synthetic-source:service-contract
    factors:
      completeness:
        title: Completeness
        description: >
          The degree to which the contract defines behavior for every endpoint
          and failure mode integrators can hit.
      consistency:
        title: Consistency
        description: >
          The degree to which the contract agrees with shipped behavior and
          with itself.
    requirements:
      contract-covers-mutation-semantics:
        title: the contract defines retry, idempotency, and error semantics for every mutation endpoint
        factors: [completeness]
        assessment: >
          Enumerate mutation endpoints from the contract's endpoint index and
          confirm each defines retry, idempotency, and error semantics; the
          index is the population, so absences are countable.
      contract-matches-shipped-behavior:
        title: contract semantics match shipped handler behavior
        factors: [consistency]
        assessment: >
          Compare contract semantics against the handler matrix and the latest
          contract-test sensor results for the same endpoints.
  persistence:
    title: Ledger persistence
    description: >
      The stored ledger state and the migrations that evolve its schema.
    source: synthetic-source:persistence
    factors:
      integrity:
        title: Integrity
        description: >
          The degree to which stored ledger state preserves accounting
          invariants under success, failure, and concurrency; a silent
          violation is the failure the whole model exists to prevent.
      recoverability:
        title: Recoverability
        description: >
          The degree to which data changes can be reversed or recovered when a
          release fails.
    requirements:
      # Veto requirement: an unacceptable result here caps the whole model
      # regardless of sibling ratings. The body's model-shape section names
      # this role.
      balance-invariants:
        title: ledger mutations preserve balance invariants
        factors: [integrity]
        assessment: >
          Run the property-based invariant suite across debit, credit, failed
          write, and concurrent paths, and inspect the nightly reconciliation
          sensor's drift report for the most recent four-week window.
        ratings:
          unacceptable: >
            Any reproducible invariant violation or unexplained reconciliation
            drift, however small.
      migration-rollback:
        title: migrations have rollback paths rehearsed against the current schema
        factors: [recoverability]
        assessment: >
          Inspect the migration runbook for rollback steps and confirm the
          most recent rehearsal record is newer than the latest schema change.
  operations:
    title: Operations
    description: >
      The runbooks, telemetry, and drills that keep the service supportable in
      production.
    source: synthetic-source:operations
    factors:
      observability:
        title: Observability
        description: >
          The degree to which operators can read customer impact from service
          health signals during an incident.
      recoverability:
        title: Recoverability
        description: >
          The degree to which incident recovery is practiced and owned rather
          than assumed.
    requirements:
      customer-impact-telemetry:
        title: health signals explain customer impact
        factors: [observability]
        assessment: >
          Inspect the dashboards-as-code definitions and alert rules for
          signals that connect service symptoms to failed customer actions,
          and confirm the definitions match the deployed dashboards.
      recovery-drill-ownership:
        title: recovery drills have current owners and recent practice records
        factors: [recoverability]
        assessment: >
          Inspect the recovery calendar and incident playbooks for the named
          current owner, the last drill date, and unresolved drill follow-up.
  # This area rates the checked-in harness artifacts themselves; the
  # agent-harnessability factor above rates the equipping capability those
  # artifacts help create.
  agent-harness:
    title: Agent harness
    description: >
      The checked-in artifacts that steer and check agent work on LedgerLite —
      the entry point, routed guidance, and the recorded sensor catalog —
      rated for their own quality, distinct from the agent harnessability
      capability they support.
    source: synthetic-source:agent-harness
    factors:
      completeness:
        title: Completeness
        description: >
          The harness covers the guides and sensors agents need for setup,
          scoped work, verification, and handoff.
      coherence:
        title: Coherence
        description: >
          Harness guidance does not contradict itself or the contract,
          runbooks, and sensors it routes to.
      currentness:
        title: Currentness
        description: >
          Harness guidance matches the current service layout, command
          surfaces, and sensor names.
      assessability:
        title: Assessability
        description: >
          Harness quality can be checked through inspectable artifacts and
          runnable sensors rather than tribal knowledge.
    requirements:
      harness-orients-agents-and-routes-to-sensors:
        title: the harness orients agents and routes them to runnable sensors
        factors:
          - completeness
          - coherence
          - currentness
          - assessability
        assessment: >
          Assess the agent entry point, routed guidance, and sensor catalog
          (synthetic-source:agent-harness) for coverage of setup, scoped work,
          verification, and handoff; agreement with the service contract and
          runbooks; current command names; and sensors an agent can actually
          run.
  # This area evaluates the concrete QUALITY.md artifact; its requirement
  # references the authoring guide family as the assessment.
  quality-md:
    title: LedgerLite Service QUALITY.md
    description: >
      This model itself — the structured frontmatter and the Markdown judgment
      context — assessed as an artifact.
    source: ./QUALITY.md
    factors:
      context-grounding:
        title: Context Grounding
        description: >
          The body explains the service, scope, needs, risks, unknowns, and
          review state well enough for a later human or agent to judge the
          model.
      evaluability:
        title: Evaluability
        description: >
          Requirements are assessable from recorded evidence and distinguish
          rating levels.
      lifecycle-maintenance:
        title: Lifecycle Maintenance
        description: >
          The model evolves with the service, and the quality changelog
          preserves why it changed.
    requirements:
      the-model-follows-the-authoring-guide-family:
        title: the quality model follows its authoring guide family
        factors:
          - context-grounding
          - evaluability
          - lifecycle-maintenance
        assessment: >
          Assess ./QUALITY.md against the /quality skill's authoring guide
          family, especially whether the body credibly supports the model,
          factors trace to visible needs and risks, requirements are
          assessable, the projection boundary between agent harnessability and
          the agent-harness area is encoded, and unknowns and open questions
          are current.
---

# Quality model: LedgerLite Service

## Overview

LedgerLite is a fictional double-entry ledger and payments service that small
product teams embed to record money movement. This model exists so its
maintainers — and the AI agents they work with — share one explicit definition
of what "good" means for it and can evaluate the service against that
definition instead of against intuition.

Integrators call the public API to record and query ledger activity; finance
users reconcile against its stored balances; operators keep it supportable in
production; maintainers and coding agents change it weekly. The governing sense
of good is fitness for purpose — money is recorded accurately and remains
explainable — backed by conformance to the service contract, which is the
normative artifact other areas are judged against. Where the two diverge, a
conformant-but-unfit result is recorded as a contract problem, not silently
passed.

_Unknowns_ — real integrator retry behavior is inferred from support tickets,
which are not agent-accessible.
_Open questions_ — should the contract commit to a public deprecation window
for the v1 error envelope?

_Reviewed — Rosa Delgado, 2026-06; agent-reviewed — Claude Code (Opus 4.8),
2026-06._

## Scope

Scope matters here because LedgerLite sits between systems the team does not
own: this model covers what the team can actually change. It includes the
public API, the service contract, ledger persistence and migrations,
operations material (runbooks, dashboards-as-code, drills), the checked-in
agent harness, and this `QUALITY.md` itself.

Out of scope by design: the third-party bank connectors (owned by the payments
platform team), the billing UI that consumes the API, and the cloud provider's
own availability. Their failures reach LedgerLite users, but this model can
only judge how LedgerLite behaves when they misbehave.

_Unknowns_ — none known.
_Open questions_ — none.

_Reviewed — Rosa Delgado, 2026-06; agent-reviewed — Claude Code (Opus 4.8),
2026-06._

## Needs

Needs are ranked here because the roll-up is not one-vote-each: balance
integrity outranks everything else this model judges.

Finance users need stored balances they can reconcile without surprises — the
benefit the service exists to deliver. Integrators need mutation semantics they
can retry against safely, and error responses they can branch on. Operators
need to read customer impact from telemetry during an incident and to know who
runs the next recovery drill. Maintainers and coding agents need the contract,
runbooks, and sensors to be current enough that a scoped change can be made,
verified, and handed off without private context.

The team works agent-first: routine changes and quality evaluations are run by
coding agents under the checked-in harness, with humans directing and
reviewing. That is a fact about how LedgerLite is maintained, not about what
this model can judge.

_Unknowns_ — holiday-peak integrator load is projected from last year's
figures.
_Open questions_ — none.

_Reviewed — Rosa Delgado, 2026-06; agent-reviewed — Claude Code (Opus 4.8),
2026-06._

## Risks

The failure that would end trust in LedgerLite is silent ledger corruption: a
balance that drifts without an error. It is low-likelihood and
catastrophic-impact, which is why `balance-invariants` is this model's veto
requirement — one reproducible violation makes the whole service unacceptable
regardless of sibling ratings, and its assessment leans on computational
sensors (the property-based invariant suite and the nightly reconciliation
job) rather than review alone.

Second-order risks: replayed mutations double-charging under retry (the
`correctness` factor traces to this, and `idempotent-mutations` makes it
inspectable — the model's worked example of a Need-to-requirement trace);
migrations that cannot be rolled back once a release is half-out; incidents
where telemetry shows symptoms but not customer impact; and drift between the
contract and shipped behavior, which quietly invalidates every judgment made
against the contract. For the agent-first workflow, the standing risk is
guidance that sounds right but is stale — sensors renamed, commands moved —
which turns agent work from bounded to improvised.

_Unknowns_ — failure behavior under a full region outage is untested.
_Open questions_ — should replay protection be load-tested at holiday-peak
volume before the next peak?

_Reviewed — Rosa Delgado, 2026-06; agent-reviewed — Claude Code (Opus 4.8),
2026-06._

## Model shape and how to read ratings

The root is composite: the API, contract, persistence, operations, harness,
and this model each carry their own factor family, so no single factor list
sits at the root. The service contract is normative — `api` requirements are
assessed against it by the same `synthetic-source:service-contract` selector
that is the contract area's own source, which is the traceability edge between
the two areas. Agent harnessability is the model-wide factor (how well the
whole project equips an agent); the `agent-harness` area rates the checked-in
harness artifacts themselves. The factors and requirements here are earned
from LedgerLite's needs and risks; they are not a default set for other
models.

Read ratings worst-of within an area: one weak requirement caps its area,
because the areas group concerns that do not compensate for each other.
`balance-invariants` is the named veto — `unacceptable` there makes the root
unacceptable outright. The required margin: money-touching areas (`api`,
`service-contract`, `persistence`) must land at `target` or better for the
service to be considered good enough to rely on; an all-`minimum` result is a
stop-and-fix signal, not a pass. Supporting areas may sit at `minimum`
temporarily with a named follow-up. A `not assessed` result stays visible as
missing evidence and is never read as a low rating.

_Unknowns_ — none known.
_Open questions_ — none.

_Reviewed — Rosa Delgado, 2026-06; agent-reviewed — Claude Code (Opus 4.8),
2026-06._
