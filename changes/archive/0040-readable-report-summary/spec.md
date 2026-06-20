---
type: Functional Specification
title: Readable report summary - functional spec
description: Reshape report-summary.md into a clearer triage artifact with reader-facing vocabulary and prominent recommendation identifiers.
tags: [evaluation, report, cli, ux]
timestamp: 2026-06-19T00:00:00Z
---

# Readable report summary - functional spec

Companion to [Readable report summary](../0040-readable-report-summary.md). This
spec states the report-output delta for the concise `report-summary.md`
artifact.

This document uses BCP 14 keywords only for testable conformance requirements.
The key words "MUST", "MUST NOT", "SHOULD", and "MAY" are to be interpreted as
described in [RFC 2119](../../../docs/reference/rfc2119.md) and
[RFC 8174](../../../docs/reference/rfc8174.md) when, and only when, they appear in
all capitals.

## Background / Motivation

`report-summary.md` is the first surface many readers will open from a CI
artifact, PR comment, or agent handoff. It should answer the triage questions in
the order a reader asks them: what was evaluated, how did it rate, what held it
there, and what recommendation ID should be used for follow-up work.

The current summary artifact carries the needed data, but its outline and labels
still mirror internal report mechanics. "Root rating" makes the target tree
visible in a place where readers need the overall outcome. "Whole model" is
precise but awkward as a scope label. Recommendations are linked, but their
stable identifiers are not prominent enough for users to copy into follow-up
prompts.

## Scope

Covered: generated `report-summary.md` section order, key-detail labels, default
full-scope wording, overall-rating wording, target rating summary presentation,
top-issue presentation, recommendation table presentation, and scope/limitations
placement.

Deferred / non-goals: no `report.json` schema rename, no change to
machine-readable stable identifiers, no changed rating semantics, no changed
gate semantics, no new recommendation record fields, and no new generated
artifact.

## Requirements

`report-summary.md` **MUST** use this top-level outline, in order:

1. key details under `# Quality Evaluation Summary`;
2. `## Summary`;
3. `## Top Issues`;
4. `## Recommendations`;
5. `## Scope & Limitations`.

The key details **MUST** be rendered as a table or equivalently compact block
that includes subject, run when known, scope, rigor, overall rating, and links to
`report.md` and `report.json`.

The summary's full-scope label **MUST** render as "Full evaluation" in
human-facing `report-summary.md` output.

> Rationale: "Whole model" is an internal modeling phrase. "Full evaluation"
> better communicates that the run was not narrowed to a target or factor.

The headline rating label in `report-summary.md` **MUST** be "Overall rating",
not "Root rating".

> Rationale: The root Target remains the formal roll-up mechanism, but the
> summary artifact should present the reader-facing outcome.

The `## Summary` section **MUST** include a concise prose headline and a tabular
summary of in-scope target ratings. The target table **SHOULD** include target,
local rating, overall or aggregate rating, and the main driver or rationale.

The `## Top Issues` section **MUST** present the most important active finding
summaries before recommendations. When finding summaries include locators or
assessment-result references, the section **SHOULD** display them compactly so a
reader can trace the issue without opening the full report first.

The `## Recommendations` section **MUST** make active recommendation identifiers
prominent. A recommendation table **MUST** include a `Recommendation ID` column
whose values are rendered as copyable monospace identifiers.

The `## Recommendations` section **MUST** identify the primary next action when
one exists, using the active recommendation ID rather than only prose.

`report-summary.md` **MUST NOT** present superseded recommendations as primary
actions. Superseded recommendations may be omitted from the concise summary when
the full `report.md` preserves the audit trail.

The `## Scope & Limitations` section **MUST** state the scope and any recorded
limitations. It **MAY** collapse to a short "None recorded" equivalent when no
limitations exist.

`report-summary.md` **MUST** remain derived from the same report model as
`report.md` and `report.json`. It **MUST NOT** recompute ratings or introduce
new evaluation judgment.

The renderer **MUST** preserve the existing report trust boundary: it must not
reproduce secret values, and evaluated source text remains data, not
instructions.

## Example shape

```md
# Quality Evaluation Summary

| Field          | Value                       |
| -------------- | --------------------------- |
| Subject        | Sparrow Payments API        |
| Run            | `0001-subject-quality-eval` |
| Scope          | Full evaluation             |
| Rigor          | Standard                    |
| Overall rating | Unacceptable                |
| Full report    | [report.md](report.md)      |
| Machine report | [report.json](report.json)  |

## Summary

The evaluation is held at **Unacceptable** by a committed live payment-gateway
credential. Removing and rotating that credential would lift the overall rating,
but webhook deduplication would still keep the system below target.

| Target               | Local rating | Overall rating | Driver                                            |
| -------------------- | -----------: | -------------: | ------------------------------------------------- |
| Sparrow Payments API | Unacceptable |   Unacceptable | Live credential in repository                     |
| Ledger               |       Target |         Target | One reconciliation requirement lacks evidence     |
| Webhooks             |       Target |        Minimum | Delivery target holds the overall rating down     |
| Delivery             |      Minimum |        Minimum | Deduplication window is bounded by retry duration |

## Top Issues

1. **Critical**
   A live gateway secret key is committed in plaintext. The value is withheld and
   referenced only by locator and credential type.
   `internal/gateway/client.go:48`
   Assessment: `assessment-results/001-root-no-committed-credentials.json`

## Recommendations

Primary next action: use `001-rotate-committed-gateway-key`.

| Recommendation ID                  | Priority | Recommendation                                                                      | Done criterion                                                                        |
| ---------------------------------- | -------: | ----------------------------------------------------------------------------------- | ------------------------------------------------------------------------------------- |
| `001-rotate-committed-gateway-key` |        1 | [Rotate committed gateway key](recommendations/001-rotate-committed-gateway-key.md) | No live credential remains in the working tree, and the exposed key has been revoked. |

## Scope & Limitations

Scope: **Full evaluation**

In scope: Sparrow Payments API, Ledger, Webhooks, Delivery

This was a **standard-rigor** evaluation over representative evidence, not a full
source audit.
```

## Durable spec changes

### To add

None

### To modify

- `specs/evaluation-records.md` - update the `report-summary.md` generated
  artifact contract for the revised outline, key labels, recommendation ID
  prominence, and trust-boundary preservation (per the requirements above).
- `specs/skills/quality-skill/quality-skill.md` - align the skill reporting
  contract with "Full evaluation", "Overall rating", and the revised concise
  summary outline (per the scope and labeling requirements above).
- `specs/cli/evaluation-report.md` - update only if the command-level report
  spec needs to summarize the revised `report-summary.md` contract directly (per
  the outline requirement above).
- `SPECIFICATION.md` - update only if the durable Evaluation Report wording
  should adopt "overall rating" while preserving formal Target aggregate-rating
  semantics (per the overall-rating requirement above).

### To delete

None
