---
type: Design Doc
title: Concern Finding Severity
description: Technical design for concern-only Finding severity validation, schemas, reports, examples, and guidance.
tags: [evaluation, findings, schema, reports]
timestamp: 2026-06-30T00:00:00Z
---

# Concern Finding Severity

## Context

The functional spec makes `severity` a concern-only Finding field. The
implementation must keep `data set`, `data verify`, `data schema`, examples,
reports, runtime guidance, and generated gallery outputs aligned without
compatibility shims.

## Approach

The structural Finding contract keeps `severity` in the allowed property set but
makes it optional. Type-specific validation then enforces the semantic rule for
each Finding in `RequirementAssessmentResult.findings`:

- `gap` and `risk`: `severity` must be present and valid;
- `strength` and `note`: `severity` must be absent.

The existing enum validator continues to validate a present `severity` value
against `critical`, `high`, `medium`, and `low`.

Schema generation gains object-level validation rules for Finding objects. The
Finding object schema keeps `severity` in `properties` but omits it from the
base `required` list, then adds conditional JSON Schema clauses:

```json
{
  "allOf": [
    {
      "if": { "properties": { "type": { "enum": ["gap", "risk"] } }, "required": ["type"] },
      "then": { "required": ["severity"] }
    },
    {
      "if": { "properties": { "type": { "enum": ["strength", "note"] } }, "required": ["type"] },
      "then": { "not": { "required": ["severity"] } }
    }
  ]
}
```

Generated reports use a shared helper to display severity:

- concern Findings render the known severity label;
- non-concern Findings render `—`.

The run report keeps its existing `Finding Summary` table near Key Details, but
its severity counts remain concern-only. The full Findings report link receives
a summary string generated from ranked Findings, with type segments in Finding
catalog order and severity counts in severity catalog order.

Examples and gallery generation stop writing `severity` on strengths and notes.
Existing synthetic gaps and risks keep their severity.

## Spec response

The validator and schema changes satisfy the data-shape requirements together:
validation is authoritative at write/verify time, and schema output teaches the
same contract to agents and automation.

The report helper keeps all severity display paths consistent: Top Findings,
`findings.md`, and Requirement Finding tables do not need separate special-case
logic.

The examples and gallery changes prove the new shape across copied example
payloads and generated report artifacts.

## Alternatives

**Keep `severity` required and hide it for strengths/notes in reports.** Rejected
because structured data would still encode a false concern signal.

**Allow but ignore `severity` on strengths/notes.** Rejected because it creates a
soft legacy path and leaves agents uncertain about whether the field matters.

**Split Finding objects into separate Go contracts per type.** Rejected because
the shared Finding shape is still useful; only one field is conditional.
Semantic validation plus schema conditionals keeps the implementation smaller.

**Rename `note` to `info`.** Rejected as out of scope. The current taxonomy uses
`note`, and the user request is about severity, not type names.

## Trade-offs & risks

The schema generator becomes slightly more expressive because it needs object
conditionals. That is worth the complexity because `data schema` is a public
discovery surface and must match validation.

This is a breaking data-shape change for existing runs whose strength Findings
include `severity`. The repo's early-alpha policy allows clean breaks, and this
case intentionally does not add migration or fallback behavior.

## Open questions

None.
