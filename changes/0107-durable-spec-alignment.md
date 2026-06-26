---
type: Change Case
title: Durable spec alignment
description: Align durable specs with the latest functional-spec and OKF guidance.
status: In-Review
tags: [specs, okf, documentation]
timestamp: 2026-06-26T00:00:00Z
---

# Durable spec alignment

This parent concept captures the why and status; the detail lives in its
children:

- [Functional spec](0107-durable-spec-alignment/spec.md) - what the change must
  do.

This case needs no design doc: it changes durable specification text and bundle
metadata only.

## Motivation

The functional-spec guidance now asks durable specs to be cumulative current
truth with explicit requirement subjects, clear BCP 14 usage, source-of-truth
boundaries, and tidy OKF listings/logs. The durable spec bundle has accumulated
new files quickly, and at least one family of short guide specs uses all-caps
requirement keywords without declaring the BCP 14 convention.

This case audits the durable spec surface and applies the concrete alignment
edits needed now, while avoiding template churn in specs that already satisfy the
current guidance.

## Scope

Covered:

- audit `SPECIFICATION.md` and all durable spec concepts under `specs/`;
- fix concrete current-guidance defects found by the audit;
- update relevant spec-bundle logs and change-bundle metadata; and
- verify with formatting plus two follow-up audit passes.

Deferred / non-goals:

- no CLI, Go, bundled skill runtime, format-schema, rating, roll-up,
  evaluation-record, or report-rendering behavior change;
- no broad heading/template normalization where a spec already reads clearly;
- no archival rewrite of historical `log.md` prose or archived Change Cases; and
- no new durable spec split unless the audit finds a contract that cannot remain
  reviewable in its current parent.

## Affected artifacts

### Code

- [x] None - no Go, CLI, or generated report implementation change expected.

### Format spec

- [x] `SPECIFICATION.md` - audit only unless the pass finds current-guidance
      defects in the format spec itself.

### Durable specs

- [x] `specs/**/*.md` - audit the durable spec bundle for OKF mechanics,
      BCP 14 declarations, change-case-only sections, vocabulary/register drift,
      stale placeholders, and index/log alignment.
- [x] `specs/skills/quality-skill/guides/authoring/*.md` - add missing BCP 14
      convention declarations and align the companion note wording where needed.
- [x] `specs/skills/quality-skill/guides/log.md` - record guide-spec alignment
      edits.
- [x] `specs/log.md` - record bundle-level durable spec alignment.

### Durable docs / bundled skill

- [x] None - this case updates durable specs and bundle metadata only.

### Suggested new durable specs

- None expected. The audit will record any split candidate if one is found, but
  the current task is alignment of existing durable specs.

## Status

`In-Review`. Durable spec audit and alignment edits are complete; verified with
`mise run fmt-md-check` and two follow-up audit passes. See the
[status lifecycle](index.md#status-lifecycle).
