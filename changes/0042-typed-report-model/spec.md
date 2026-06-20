---
type: Functional Specification
title: Typed report model - functional spec
description: Requirements for replacing implicit and stringly typed evaluation-report states with explicit typed model concepts.
tags: [evaluation, report, records, types]
timestamp: 2026-06-20T00:00:00Z
---

# Typed report model - functional spec

Companion to [Typed report model](../0042-typed-report-model.md).

This document uses BCP 14 keywords only for testable conformance requirements.
The key words "MUST", "MUST NOT", "SHOULD", and "MAY" are to be interpreted as
described in [RFC 2119](../../../docs/reference/rfc2119.md) and
[RFC 8174](../../../docs/reference/rfc8174.md) when, and only when, they appear in
all capitals.

## Background / Motivation

Evaluation reports are consumed by humans and agents. Several current report
concepts are semantically important but structurally weak: rating state is a raw
string, structural local ratings are represented by `nil` plus a boolean, record
lifecycle is compressed into `active`, and reportability routing compares string
gap names. These encodings let invalid states drift into reports and make later
changes brittle. The report model should carry the distinctions the report
contract already requires.

## Scope

Covered: typed rating-result kinds, explicit local rating state, typed next step,
typed run gaps, typed lifecycle state, typed rigor and evaluation level,
structured missing metadata, and target/factor path helper types.

Deferred / non-goals: no QUALITY.md frontmatter schema change; no new rating
aggregation behavior; no new report sections; no compatibility aliases for
invalid values.

## Requirements

### Rating results

Rating results **MUST** use a typed kind with the canonical values `rated` and
`not-assessed`. Validation **MUST** reject unknown kinds, a rated result without
a level, a not-assessed result with a level, and an empty rationale.

### Local rating state

Report target details and summaries **MUST** represent local rating as an
explicit state with canonical kinds `rated`, `not-assessed`, and `structural`.
The structural state **MUST** remain distinct from a not-assessed rating in JSON
and Markdown.

### Next step

The report next action **MUST** be modeled internally as a typed next step with
canonical kinds `recommendation` and `none`. Recommendation next steps **MUST**
carry a recommendation id and path; none next steps **MUST NOT** carry
recommendation locators.

### Run gaps

Evaluation run gaps **MUST** use typed gap kinds. Routing policy such as whether
a gap requires review **MUST** live on the gap kind rather than in ad hoc string
lists.

### Record lifecycle

Assessment-result and recommendation report digests **MUST** expose a lifecycle
state with canonical values `active` and `superseded`. `supersedes` remains audit
data, but consumers must not infer lifecycle only from a boolean.

### Evidence kind

Evidence kind remains an open classification string. The renderer **MUST** treat
it as display/grouping metadata only and **MUST NOT** branch on undocumented
evidence kind values.

### Rigor and evaluation level

Rigor **MUST** use a typed value with canonical levels `quick`, `standard`, and
`deep`; missing rigor must remain distinct from an unknown value. Evaluation
level **MUST** use a typed value for the canonical `subject` level and deliberately
handle historical `model` runs.

### Missing metadata

Report missing metadata **MUST** be structured with a stable field id and a
display title instead of bare prose strings.

### Paths

Target and factor paths **SHOULD** use helper types for identity keys and display
fallbacks while preserving the existing JSON array representation.

## Durable spec changes

### To add

None

### To modify

- `specs/evaluation-records.md` — require typed rating states, local rating
  states, lifecycle states, next-step states, missing metadata, and gap behavior.
- `specs/skills/quality-skill/quality-skill.md` — align the skill contract with
  explicit report states and typed record validation.
- `specs/cli/evaluation-status.md` — describe typed run gaps when reportability
  status reports invalid records.
- `specs/cli/evaluation-report.md` — describe the report JSON state objects and
  Markdown rendering distinctions.

### To delete

None
