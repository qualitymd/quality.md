---
schemaVersion: 1
title: Produce reconciliation evidence so the requirement can be rated
gap: No reconciliation job output, log, or report was available, so the requirement could not be assessed.
evidenceLocators:
  - ledger/reconcile.go:31
assessmentResultRecords:
  - assessment-results/005-ledger-reconciliation.json
remediationOptions:
  - Stand up the scheduled reconciliation run and emit a durable report.
  - Surface existing reconciliation output if it already runs elsewhere.
  - Narrow or retire the requirement.
recommendedOption: Stand up the scheduled reconciliation run and emit a durable report.
doneCriterion: The reconciliation requirement becomes assessable and reaches at least the acceptable floor.
---

# Produce reconciliation evidence so the requirement can be rated

**Target / factor:** Ledger → Correctness
**In-scope requirement:** *Reconciliation runs daily and flags drift*
**Current rating:** Not assessed — no evidence was available to rate against.

## Gap

The requirement could not be assessed: no reconciliation job output, log, or
report was available, so there is no evidence to rate against the scale. This is
an evaluator-determined *not assessed* outcome for this run — distinct from a
standing **Known gap** the author has declared in the model body. The Ledger
subtree's Correctness rating currently rests on the double-entry invariant alone
and is noted as incomplete.

**Evidence:**

- `ledger/reconcile.go:31` — a `reconcile` entrypoint exists, but no evidence of
  a scheduled or recent run (no job log, report, or run record) was found at
  `standard` rigor.

## Options

- **(a) Stand up the scheduled reconciliation run and emit a durable report.**
  Run `reconcile` on a daily schedule and write a report (drift count, window,
  timestamp) the evaluator can assess against.
- **(b) Surface existing reconciliation output if it already runs elsewhere.**
  If reconciliation already runs outside the repo, route its report to a location
  the evaluation can resolve, so the existing run becomes assessable evidence.
- **(c) Narrow or retire the requirement.** If daily reconciliation is not
  actually a requirement for the Ledger, remove or rescope it in the model so the
  model stops measuring something that is not expected.

## Recommended

**(a) Stand up the scheduled reconciliation run and emit a durable report.** The
requirement reflects a real correctness need for a ledger (catching drift), so
producing the evidence is preferable to retiring it. Option (b) is equivalent if
the run already exists; option (c) is only right if the need itself does not.

## Done-criterion

The requirement *Reconciliation runs daily and flags drift* moves from *not
assessed* to a rated level — at least **Target** against its criterion: a daily
run exists and its report shows drift is detected and flagged within the expected
window. A later `improve` re-evaluates the Ledger scope to confirm the
requirement is now assessable and that the subtree rating reflects full coverage
rather than resting on a single requirement.
