---
title: LedgerLite Service
description: >
  Illustrative quality model for LedgerLite, a fictional ledger and payments
  service: its public API, service contract, persistence, operations, codebase,
  agent harness, and this model itself.
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
      output, and stay safely bounded - distinct from the agent-harness area,
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
              Run the recorded sensor commands - contract tests, ledger
              invariant suite, and lint - and confirm each returns a
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
      reliability:
        title: Reliability
        description: >
          The degree to which API behavior remains dependable when dependencies
          time out, callers retry, or traffic shifts during peak periods.
      security:
        title: Security
        description: >
          The degree to which API access prevents unauthorized money movement,
          tenant data exposure, and privilege escalation.
      compatibility:
        title: Compatibility
        description: >
          The degree to which the API preserves documented v1 caller behavior
          while deprecations move through an explicit window.
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
      testability:
        title: Testability
        description: >
          The degree to which API expectations can be verified by contract
          tests and seeded replay scenarios rather than manual inspection.
    requirements:
      idempotent-mutations:
        title: mutation endpoints are idempotent under retry
        factors: [correctness, reliability]
        assessment: >
          Run the contract tests for mutation endpoints and compare their replay
          cases with the retry and idempotency section of the service contract,
          including duplicate-key and interrupted-write replay behavior.
      ledger-entry-signs-match-intent:
        title: ledger entry signs match caller intent
        factors: [correctness]
        assessment: >
          Run the contract tests over debit, credit, refund, and reversal cases
          and confirm recorded signs match the operation intent in the service
          contract.
      dependency-timeouts-return-safe-results:
        title: downstream dependency timeouts return safe results
        factors: [reliability, operability]
        assessment: >
          Inspect timeout-injection contract tests and the runbook to confirm a
          bank-connector timeout returns a retryable failure without writing a
          partial ledger entry.
      tenant-access-is-enforced:
        title: tenant access is enforced for every money-moving endpoint
        factors: [security]
        assessment: >
          Run authorization matrix tests for every money-moving endpoint and
          inspect handler coverage for cross-tenant account identifiers.
      sensitive-fields-stay-out-of-error-responses:
        title: sensitive fields stay out of error responses
        factors: [security, operability]
        assessment: >
          Run error-envelope contract tests and inspect sampled failure payloads
          for account numbers, bank tokens, internal IDs, and tenant secrets.
      v1-error-envelope-remains-compatible:
        title: v1 error-envelope behavior remains compatible during deprecation
        factors: [compatibility, operability]
        assessment: >
          Compare the compatibility matrix, deprecation notice, and handler
          matrix to confirm v1 callers still receive documented fields until
          the public deprecation window closes.
      predictable-error-contracts:
        title: error responses are predictable for callers
        factors: [operability]
        assessment: >
          Compare the error-envelope section of the service contract with
          handler behavior for validation, authorization, conflict, and timeout
          cases across the endpoint index.
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
      contract-tests-cover-public-endpoints:
        title: contract tests cover every public endpoint
        factors: [testability]
        assessment: >
          Compare the endpoint index against the contract-test manifest and
          confirm every public endpoint has at least one success and one failure
          case.
  service-contract:
    title: Service contract
    description: >
      The normative API contract that defines endpoint, retry, idempotency,
      compatibility, and error semantics; other areas are judged against it,
      and this area judges the contract itself.
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
      currentness:
        title: Currentness
        description: >
          The degree to which contract claims match the latest shipped behavior,
          deprecation state, and handler matrix.
      understandability:
        title: Understandability
        description: >
          The degree to which integrators can interpret endpoint semantics,
          retries, errors, and compatibility promises without private support.
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
        factors: [consistency, currentness]
        assessment: >
          Compare contract semantics against the handler matrix and the latest
          contract-test sensor results for the same endpoints.
      deprecation-window-is-current:
        title: the v1 deprecation window is current and visible
        factors: [currentness, understandability]
        assessment: >
          Compare the compatibility appendix, changelog, and handler matrix to
          confirm the v1 error-envelope deprecation date and shipped fields are
          aligned.
      examples-explain-retry-and-error-semantics:
        title: examples explain retry and error semantics for integrators
        factors: [understandability, completeness]
        assessment: >
          Review contract examples for duplicate retries, interrupted writes,
          authorization failures, and validation errors, and confirm each
          example names the caller action.
  persistence:
    title: Ledger persistence
    description: >
      The stored ledger state, reconciliation traces, and migrations that evolve
      the schema.
    source: synthetic-source:persistence
    factors:
      integrity:
        title: Integrity
        description: >
          The degree to which stored ledger state preserves accounting
          invariants under success, failure, and concurrency; a silent
          violation is the failure the whole model exists to prevent.
      auditability:
        title: Auditability
        description: >
          The degree to which a stored balance remains explainable through
          durable, ordered, tamper-evident records.
      security:
        title: Security
        description: >
          The degree to which stored ledger data and migration access are
          protected from unauthorized reads, writes, and privilege expansion.
      recoverability:
        title: Recoverability
        description: >
          The degree to which data changes can be reversed or recovered when a
          release fails.
      durability:
        title: Durability
        description: >
          The degree to which committed ledger entries remain present and
          reproducible across process restarts, migrations, and restore drills.
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
          job's drift report for the most recent four-week window.
        ratings:
          unacceptable: >
            Any reproducible invariant violation or unexplained reconciliation
            drift, however small.
      reconciliation-explains-balance-changes:
        title: reconciliation explains every balance change
        factors: [integrity, auditability]
        assessment: >
          Compare the nightly reconciliation job with the audit-event stream and
          confirm every sampled balance delta traces to one ordered ledger event.
      audit-events-are-ordered-and-tamper-evident:
        title: audit events are ordered and tamper-evident
        factors: [auditability]
        assessment: >
          Inspect the audit-log schema and append-path tests for monotonic
          sequence numbers, hash chaining, and immutable event writes.
      persistence-access-is-least-privilege:
        title: persistence access is least-privilege
        factors: [security]
        assessment: >
          Run the dependency audit and inspect database role manifests to
          confirm service, migration, and analytics roles cannot write outside
          their documented scope.
      migration-rollback:
        title: migrations have rollback paths rehearsed against the current schema
        factors: [recoverability, durability]
        assessment: >
          Inspect the migration runbook for rollback steps and confirm the most
          recent rehearsal record is newer than the latest schema change.
      restore-drills-replay-current-backups:
        title: restore drills replay current backups without ledger loss
        factors: [durability, recoverability]
        assessment: >
          Inspect the latest restore drill, backup manifest, and reconciliation
          result to confirm restored balances match the pre-drill ledger state.
  operations:
    title: Operations
    description: >
      The runbooks, telemetry, capacity plans, access controls, and drills that
      keep the service supportable in production.
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
      security:
        title: Security
        description: >
          The degree to which operational access, incident tools, and emergency
          paths avoid unnecessary privilege and leave reviewable records.
      capacity:
        title: Capacity
        description: >
          The degree to which operational evidence shows the service can absorb
          projected peak traffic without violating latency or recovery targets.
    requirements:
      customer-impact-telemetry:
        title: health signals explain customer impact
        factors: [observability]
        assessment: >
          Inspect the dashboards-as-code definitions and alert rules for signals
          that connect service symptoms to failed customer actions, and confirm
          the definitions match the deployed dashboards.
      recovery-drill-ownership:
        title: recovery drills have current owners and recent practice records
        factors: [recoverability]
        assessment: >
          Inspect the recovery calendar and incident playbooks for the named
          current owner, the last drill date, and unresolved drill follow-up.
      break-glass-access-is-reviewed:
        title: break-glass access is reviewed after use
        factors: [security, recoverability]
        assessment: >
          Inspect break-glass logs for the latest quarter and confirm every use
          has a named approver, incident link, and post-use access review.
      holiday-peak-capacity-is-supported-by-load-evidence:
        title: holiday-peak capacity is supported by load evidence
        factors: [capacity]
        assessment: >
          Compare the capacity plan with the latest load-test rollup and
          production traffic forecast for the next holiday peak.
  codebase:
    title: Codebase
    description: >
      The implementation that realizes the API, persistence, operations hooks,
      and contract conformance.
    source: synthetic-source:codebase
    factors:
      maintainability:
        title: Maintainability
        description: >
          The degree to which maintainers and coding agents can understand,
          change, and verify the implementation safely during weekly work.
        factors:
          analyzability:
            title: Analyzability
            description: >
              The degree to which code structure, names, and boundaries reveal
              how money movement flows through the service.
            requirements:
              money-flow-is-analyzable:
                title: money movement flow is analyzable from entry point to ledger write
                assessment: >
                  Inspect the handler-to-ledger trace and complexity check to
                  confirm a maintainer can follow debit, credit, refund, and
                  reversal flow without crossing unrelated modules.
          modifiability:
            title: Modifiability
            description: >
              The degree to which localized API, contract, and persistence
              changes can be made without broad unintended edits.
            requirements:
              changes-remain-local-to-owned-boundaries:
                title: changes remain local to owned architecture boundaries
                factors: [consistency]
                assessment: >
                  Run structural import-boundary tests and inspect the
                  changed-module matrix for recent API and persistence changes.
          testability:
            title: Testability
            description: >
              The degree to which implementation changes can be checked by
              focused tests and sensors at the right boundary.
            requirements:
              implementation-has-focused-tests-for-risky-branches:
                title: risky implementation branches have focused tests
                assessment: >
                  Compare the branch inventory for retry, authorization,
                  rollback, and reconciliation paths against focused unit or
                  contract tests.
      consistency:
        title: Consistency
        description: >
          The degree to which implementation boundaries conform to the service
          architecture and the contract they realize.
      security:
        title: Security
        description: >
          The degree to which implementation dependencies, secret handling, and
          authorization paths resist preventable exposure or misuse.
    requirements:
      architecture-boundaries-match-the-service-contract:
        title: architecture boundaries match the service contract
        factors: [consistency]
        assessment: >
          Run structural import-boundary tests and compare handler ownership
          against the service contract's endpoint families.
      dependency-and-secret-handling-stay-within-policy:
        title: dependencies and secret handling stay within policy
        factors: [security]
        assessment: >
          Run the dependency audit and lint rules for secrets, then inspect any
          suppressions against the recorded security policy.
  # This area rates the checked-in harness artifacts themselves; the
  # agent-harnessability factor above rates the equipping capability those
  # artifacts help create.
  agent-harness:
    title: Agent harness
    description: >
      The checked-in artifacts that steer and check agent work on LedgerLite -
      the entry point, routed guidance, and the recorded sensor catalog - rated
      for their own quality, distinct from the agent harnessability capability
      they support.
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
        factors: [completeness, coherence, currentness, assessability]
        assessment: >
          Assess the agent entry point, routed guidance, and sensor catalog for
          coverage of setup, scoped work, verification, and handoff; agreement
          with the service contract and runbooks; current command names; and
          sensors an agent can actually run.
      sensor-catalog-names-reusable-sensors:
        title: the sensor catalog names reusable sensors consistently
        factors: [completeness, currentness, assessability]
        assessment: >
          Compare the sensor catalog against requirement assessments in this
          QUALITY.md and confirm each catalog sensor name is used consistently
          by at least two assessments.
  # This area evaluates the concrete QUALITY.md artifact; its requirements
  # reference the authoring guide family as the assessment.
  quality-md:
    title: LedgerLite Service QUALITY.md
    description: >
      This model itself - the structured frontmatter and the Markdown judgment
      context - assessed as an artifact.
    source: ./QUALITY.md
    factors:
      credibility:
        title: Credibility
        description: >
          The body grounds model claims in specific LedgerLite needs, risks,
          evidence limits, and review provenance.
      assessability:
        title: Assessability
        description: >
          Requirements are assessable from recorded evidence and distinguish
          rating levels.
      currentness:
        title: Currentness
        description: >
          The model and its quality changelog stay aligned with LedgerLite's
          current scope, risks, factors, and sensors.
    requirements:
      the-model-follows-the-authoring-guide-family:
        title: the quality model follows its authoring guide family
        factors: [credibility, assessability, currentness]
        assessment: >
          Assess ./QUALITY.md against the /quality skill's authoring guide
          family, especially whether the body credibly supports the model,
          factors trace to visible needs and risks, requirements are
          assessable, the projection boundary between agent harnessability and
          the agent-harness area is encoded, and unknowns and open questions
          are current.
      the-quality-changelog-explains-model-growth:
        title: the quality changelog explains meaningful model growth
        factors: [credibility, currentness]
        assessment: >
          Inspect .quality/changelog entries for meaningful model changes,
          including the codebase and security expansion and the sensor
          maturation from inferential review to named sensors.
---

# Quality model: LedgerLite Service

## Overview

LedgerLite is a fictional double-entry ledger and payments service that small
product teams embed to record money movement. This model exists so its
maintainers - and the AI agents they work with - share one explicit definition
of what "good" means for the service and can evaluate it against recorded
evidence rather than intuition.

Integrators call the public API to record and query ledger activity; finance
users reconcile against stored balances; operators keep the service supportable
in production; maintainers and coding agents change it weekly. The governing
sense of good is fitness for purpose - money is recorded accurately and remains
explainable - backed by conformance to the service contract, which is the
normative artifact other areas are judged against. Where conformance and
fitness diverge, a conformant-but-unfit result is recorded as a contract problem
rather than silently passed.

_Unknowns_ - real integrator retry behavior is inferred from support tickets,
which are not agent-accessible.
_Open questions_ - should LedgerLite publish tenant-specific replay limits for
high-volume integrators?

_Reviewed - Rosa Delgado, 2026-06; agent-reviewed - Codex (GPT-5), 2026-07._

## Scope

Scope matters because LedgerLite sits between systems the team does not own:
this model covers what the team can actually change. It includes the public
API, the service contract, ledger persistence and migrations, operations
material (runbooks, dashboards-as-code, access logs, capacity plans, drills),
the implementation codebase, the checked-in agent harness, and this
`QUALITY.md` itself.

Out of scope by design: the third-party bank connectors (owned by the payments
platform team), the billing UI that consumes the API, mobile client usability,
packaging portability beyond the hosted service runtime, and the cloud
provider's own availability. Their failures reach LedgerLite users, but this
model judges only how LedgerLite behaves when those systems misbehave.

_Unknowns_ - none known.
_Open questions_ - none.

_Reviewed - Rosa Delgado, 2026-06; agent-reviewed - Codex (GPT-5), 2026-07._

## Needs

Needs are ranked because the roll-up is not one-vote-each: balance integrity
outranks every other concern this model judges.

Finance users need stored balances they can reconcile without surprises - the
benefit the service exists to deliver. Integrators need mutation semantics they
can retry against safely, v1 compatibility they can plan around, and error
responses they can branch on. Operators need to read customer impact from
telemetry during an incident, know who owns the next recovery drill, review
break-glass access, and trust that holiday-peak capacity has been tested rather
than hoped for. Maintainers and coding agents need architecture boundaries,
focused tests, dependency checks, contract tests, runbooks, and the sensor
catalog to be current enough that a scoped change can be made, verified, and
handed off without private context.

The team works agent-first: routine changes and quality evaluations are run by
coding agents under the checked-in harness, with humans directing and reviewing.
That is a fact about how LedgerLite is maintained, not about what every
QUALITY.md model judges.

_Unknowns_ - holiday-peak integrator load is projected from last year's figures
and a sales forecast, not from committed customer traffic.
_Open questions_ - should high-volume integrators receive a separate replay
test pack before the next peak?

_Reviewed - Rosa Delgado, 2026-06; agent-reviewed - Codex (GPT-5), 2026-07._

## Risks

The failure that would end trust in LedgerLite is silent ledger corruption: a
balance that drifts without an error. It is low-likelihood and
catastrophic-impact, which is why `balance-invariants` is this model's veto
requirement. Its assessment leans on computational sensors - the invariant
suite and the nightly reconciliation job - rather than review alone.

Second-order risks: replayed mutations double-charging under retry; unauthorized
money movement across tenant boundaries; error-envelope drift during the v1
deprecation; migrations that cannot be rolled back once a release is half-out;
audit records that cannot explain a balance; break-glass access that is used but
not reviewed; capacity that looks healthy on ordinary load but fails during peak
traffic; and architecture drift that makes agent-authored changes sprawl across
module boundaries. For the agent-first workflow, the standing risk is guidance
that sounds right but is stale - sensors renamed, commands moved, or done
criteria missing - which turns bounded work into improvised work.

_Unknowns_ - full-region outage behavior is still owned by the platform disaster
recovery plan and is not assessed here.
_Open questions_ - should replay protection be load-tested at holiday-peak
volume before the next peak?

_Reviewed - Rosa Delgado, 2026-06; agent-reviewed - Codex (GPT-5), 2026-07._

## Sensor catalog and assessment posture

Assessment evidence is deliberately economical: a small set of named sensors is
reused across many requirements, and the model leaves judgment visible where no
sensor yet exists.

The current catalog is: contract tests, invariant suite, reconciliation job,
latency rollup, load-test rollup, lint, complexity check, dependency audit,
drift detector, structural import-boundary tests, authorization matrix, restore
drill, and break-glass log review. Assessments prefer these recorded
computational sensors where they exist. Some requirements read a sensor result
against a guide, runbook, contract, or policy. A smaller set remains
inferential, especially where the evidence is example quality, review
provenance, or whether a body section earns a factor. The quality loop's job is
to shrink that inferential set by adding sensors only where the sensor would
make a real judgment more repeatable.

_Unknowns_ - the drift detector's contract-vs-handler comparison does not yet
cover every deprecated field.
_Open questions_ - should the sensor catalog become a checked-in source file
rather than a harness section?

_Reviewed - Rosa Delgado, 2026-06; agent-reviewed - Codex (GPT-5), 2026-07._

## Model shape and how to read ratings

The root is composite: API, contract, persistence, operations, codebase,
harness, and this model each carry their own factor family, so no single factor
list sits at the root. The service contract is normative - `api` requirements
are assessed against it by the same `synthetic-source:service-contract`
selector that is the contract area's own source, which is the traceability edge
between the two areas. Agent harnessability is the model-wide factor (how well
the whole project equips an agent); the `agent-harness` area rates the
checked-in harness artifacts themselves. The `codebase` area carries the
implementation concerns that were previously homeless: maintainability through
analyzability, modifiability, and testability; architecture consistency; and
implementation security.

Read ratings worst-of within an area: one weak requirement caps its area,
because the areas group concerns that do not compensate for each other.
`balance-invariants` is the named veto - `unacceptable` there makes the root
unacceptable outright. The required margin: money-touching areas (`api`,
`service-contract`, `persistence`) and the implementation paths that realize
money movement must land at `target` or better for the service to be considered
good enough to rely on; an all-`minimum` result is a stop-and-fix signal, not a
pass. Supporting areas may sit at `minimum` temporarily with a named follow-up.
A `not assessed` result stays visible as missing evidence and is never read as a
low rating.

The model omits `usability` because the billing UI and client onboarding flows
are outside this service boundary. It omits broad `portability` because
LedgerLite is operated as a hosted service in one supported runtime; restore and
deployment concerns are modeled under persistence durability and operations
recoverability instead. Those omissions are deliberate scope choices, not
claims that the concerns never matter.

_Unknowns_ - none known.
_Open questions_ - none.

_Reviewed - Rosa Delgado, 2026-06; agent-reviewed - Codex (GPT-5), 2026-07._
