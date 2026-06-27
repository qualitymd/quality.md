---
type: Design Doc
title: Requirement Findings Only - design
description: How Evaluation schema, validation, reports, and skill guidance remove Area Findings and require finding-backed ratings.
tags: [evaluation, findings, reports, skill, cli]
timestamp: 2026-06-27T00:00:00Z
---

# Requirement Findings Only - design

## Context

Answers the [functional spec](spec.md) for change case
[0142](../0142-requirement-findings-only.md). The change makes Requirement
Findings the only finding layer, removes `AreaAnalysisResult.findings`, and
requires rated outputs to carry traceable rating drivers. It is a clean break:
Evaluation data moves to schema version 3, and schema version 2 runs remain
incompatible historical data.

## Approach

### Schema v3 and payload contracts

Bump `internal/evaluation.SchemaVersion` from `2` to `3`. Keep the existing
version gate strict: `data set`, examples, generated schemas, receipts, status,
and report output now speak only schema version 3.

Remove `AreaAnalysisResult.findings` from `areaAnalysisResultContract()`. Delete
the Area Finding and Factor relationship contracts and the duplicate/relationship
semantic validator. The generic Finding Core remains, but it is attached only to
`RequirementAssessmentResult.findings[]`.

### Effective-run validation on data set

`data set` already validates a whole incoming batch before writing. Add a second
cross-payload validation pass after candidate paths are derived and before any
write occurs:

```text
effective run data = existing persisted payloads overlaid with incoming batch
```

The validator indexes payloads by derived path and by routine subject. For every
effective payload:

- rated Requirement results require a paired assessed/partially assessed
  Requirement Assessment with at least one finding;
- rated Requirement results require non-empty `ratingDrivers`;
- analyzed Factor/Area scopes with `ratingLevelId` require non-empty
  `ratingDrivers`;
- rating driver `inputRefs` resolve to an effective payload of the referenced
  kind and subject.

This keeps batch order irrelevant and allows manual incremental writes as long
as the referenced paired payload already exists. It deliberately does not decide
whether a finding semantically proves the configured Rating Level; that remains
skill judgment.

`data verify` reuses the same validation over persisted payloads by treating the
run as an effective set with no candidate overlay.

### Reports

Delete Area/Factor Findings rendering from the report tree:

- remove `## Findings` from Area reports;
- remove `## Findings` from Factor reports;
- keep Requirement report Finding summaries and details unchanged;
- keep `## Rating Drivers` on Area and Factor reports as the roll-up explanation
  surface.

Delete helper types and sort functions that exist only for Area Findings and
Factor relationships. Driver tables continue to render their `inputRefs` as
compact JSON for now; richer driver-link rendering can be a later report UX case.

### Runtime skill and durable specs

Update `/quality evaluate` guidance so collection produces Requirement Findings
for every rated Requirement and roll-up analysis produces rating drivers instead
of Area Findings. Keep the QC sweep, but narrow the "no findings" adversarial
re-check to Requirements rather than Areas.

Update durable specs to carry the enduring rationale: findings are evidence;
rating drivers are synthesis; recommendations are the later action layer.

### Release notes

Record this as a breaking Evaluation data/report change. Because the CLI and
skill depend on the same v3 contract, the release should bump to the next minor
compatibility line and update `/quality` skill metadata accordingly.

## Spec response

- **Requirement Findings and ratings** - satisfied by skill/runtime guidance plus
  CLI cross-payload validation for the deterministic "rated with no findings"
  case.
- **Roll-up analysis** - satisfied by deleting Area Findings and keeping
  `ratingDrivers` required for rated roll-up outputs.
- **CLI validation and schema** - satisfied by schema version 3, strict version
  rejection, effective-run validation, and regenerated schema/examples.
- **Reports and skill behavior** - satisfied by removing Area/Factor Findings
  sections and preserving Requirement Finding rendering.

## Alternatives

- **Keep Area Findings but make them rare.** Rejected. The second finding layer
  still duplicates Requirement evidence and invites agents to synthesize
  findings where drivers or recommendations belong.
- **Enforce only in the skill.** Rejected. Skill judgment owns sufficiency, but
  the CLI can deterministically reject the known-bad shape: rated Requirement,
  no paired findings.
- **Validate semantic sufficiency in the CLI.** Rejected. Rating Scales are
  configurable, so mechanical validation should not infer whether a finding
  proves a Rating Level.
- **Keep schema version 2 and just reject Area Findings.** Rejected. Removing a
  persisted payload field is a payload-shape break, and existing v2 runs should
  be clearly incompatible rather than half-readable.

## Trade-offs & risks

- Old Evaluation runs cannot be rebuilt with the new CLI. This is acceptable
  under early-alpha clean-break policy, and status/report tooling already treats
  schema-incompatible historical runs as history rather than current evidence.
- Requiring `ratingDrivers` makes agent-authored payloads more verbose. The
  payoff is that roll-up ratings remain auditable after removing Area Findings.
- Driver references are still rendered as compact JSON in reports. That is less
  polished than clickable driver links but keeps this change scoped to the
  evidence-model cleanup.

## Open questions

None.
