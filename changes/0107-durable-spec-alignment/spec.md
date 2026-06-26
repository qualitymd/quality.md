---
type: Functional Specification
title: Durable spec alignment - functional spec
description: Requirements for aligning durable specs with the latest functional-spec and OKF guidance.
tags: [specs, okf, documentation]
timestamp: 2026-06-26T00:00:00Z
---

# Durable spec alignment - functional spec

Companion to the [Durable spec alignment](../0107-durable-spec-alignment.md)
change case. This spec states what the alignment pass must do. No design doc is
needed because the change is a documentation/specification maintenance pass.

The key words **MUST**, **MUST NOT**, **SHOULD**, and **MAY** are to be
interpreted as described in IETF RFC 2119 and RFC 8174 when, and only when, they
appear in all capitals.

## Background / Motivation

The durable spec bundle is the current source of truth for `qualitymd` behavior,
while archived Change Cases explain past deltas. After the functional-spec guide
was strengthened, durable specs need a pass that checks the new guidance against
the actual bundle and applies only earned fixes. The goal is not to stamp every
spec with the same headings; it is to remove concrete drift that weakens future
maintenance.

## Scope

Covered: the repository-root `SPECIFICATION.md`, durable concepts under
`specs/`, and the OKF index/log/schema files that describe the `specs/` bundle.

Not covered: code, CLI behavior, bundled runtime skill files, generated report
behavior, archived Change Cases, historical log wording, and broad prose
normalization that does not change conformance clarity.

## Requirements

### Audit

- The alignment pass **MUST** audit every durable spec concept under `specs/` and
  the repository-root `SPECIFICATION.md` for current-guidance issues:
  frontmatter/type presence, registered concept types, enclosing index entries,
  BCP 14 declaration coverage, change-case-only sections, stale placeholders,
  vocabulary/register drift, and obvious parent/child contract mismatch.

- The audit **MUST NOT** treat historical `log.md` entries or archived Change
  Cases as live spec prose to be rewritten.

### Alignment edits

- Any live durable spec concept that uses all-caps BCP 14 requirement keywords
  **MUST** declare the convention near the top of the spec.

- Short runtime-artifact guide specs **SHOULD** identify the concrete runtime
  guide they govern before stating requirements, so the companion relationship is
  visible without reading the enclosing index.

- Durable specs **MUST NOT** carry the change-case-only `Durable spec changes`
  section; that section belongs only in Change Case functional specs.

- Alignment edits **MUST** preserve the guide's palette rule: do not add
  Background, Scope, Assumptions, or other headings to a durable spec unless the
  spec earns that heading.

- Alignment edits **MUST NOT** change deterministic CLI behavior, Evaluation v2
  runtime contracts, the QUALITY.md file format, bundled skill runtime guidance,
  or generated artifacts.

### OKF and records

- Every touched durable spec concept **MUST** remain listed in its enclosing
  `index.md`.

- The alignment pass **MUST** update relevant `specs/` bundle logs for edited
  durable spec files.

- The Change Case parent **MUST** reconcile its affected-artifacts checklist
  before the case reaches `In-Review`.

### Verification

- The alignment pass **MUST** run the repository Markdown formatting check.

- After the work first appears complete, the implementer **MUST** perform two
  follow-up audit passes against the current tree before marking the goal
  complete.

## Durable spec changes

### To add

None.

### To modify

- [`specs/skills/quality-skill/guides/authoring/agent-harness.md`](../../specs/skills/quality-skill/guides/authoring/agent-harness.md)
  - add missing BCP 14 declaration and clarify the governed guide note. Driven
    by [Alignment edits](#alignment-edits).
- [`specs/skills/quality-skill/guides/authoring/agent-harnessability.md`](../../specs/skills/quality-skill/guides/authoring/agent-harnessability.md)
  - add missing BCP 14 declaration and clarify the governed guide note. Driven
    by [Alignment edits](#alignment-edits).
- [`specs/skills/quality-skill/guides/authoring/body.md`](../../specs/skills/quality-skill/guides/authoring/body.md)
  - add missing BCP 14 declaration and clarify the governed guide note. Driven
    by [Alignment edits](#alignment-edits).
- [`specs/skills/quality-skill/guides/authoring/factors.md`](../../specs/skills/quality-skill/guides/authoring/factors.md)
  - add missing BCP 14 declaration and clarify the governed guide note. Driven
    by [Alignment edits](#alignment-edits).
- [`specs/skills/quality-skill/guides/authoring/model-structure.md`](../../specs/skills/quality-skill/guides/authoring/model-structure.md)
  - add missing BCP 14 declaration and clarify the governed guide note. Driven
    by [Alignment edits](#alignment-edits).
- [`specs/skills/quality-skill/guides/authoring/quality-log.md`](../../specs/skills/quality-skill/guides/authoring/quality-log.md)
  - add missing BCP 14 declaration and clarify the governed guide note. Driven
    by [Alignment edits](#alignment-edits).
- [`specs/skills/quality-skill/guides/authoring/rating-scale.md`](../../specs/skills/quality-skill/guides/authoring/rating-scale.md)
  - add missing BCP 14 declaration and clarify the governed guide note. Driven
    by [Alignment edits](#alignment-edits).
- [`specs/skills/quality-skill/guides/authoring/requirements.md`](../../specs/skills/quality-skill/guides/authoring/requirements.md)
  - add missing BCP 14 declaration and clarify the governed guide note. Driven
    by [Alignment edits](#alignment-edits).
- [`specs/skills/quality-skill/guides/log.md`](../../specs/skills/quality-skill/guides/log.md)
  - record the guide-spec alignment edits. Driven by
    [OKF and records](#okf-and-records).
- [`specs/log.md`](../../specs/log.md) - record the bundle-level durable spec
  alignment pass. Driven by [OKF and records](#okf-and-records).

### To rename

None.

### To delete

None.

## Validation check

If every requirement above is satisfied, the durable spec bundle will be aligned
with the latest functional-spec and OKF guidance where the current tree shows
concrete drift, without changing runtime behavior or flattening each spec into a
template. That achieves the motivation and leaves any future substantive
contract reshaping to its own earned Change Case.
