---
type: Functional Specification
title: Remove wizard mode
description: Remove wizard from the /quality skill's public contract while preserving safe read-only orientation for ambiguous requests.
tags: [skill, ux, contract]
timestamp: 2026-06-23T00:00:00Z
---

# Remove wizard mode

This Change Case spec defines the delta for `/quality` skill guidance and
user-facing docs: remove `wizard` as a public mode and invocation without
turning `status`, `next`, or review-oriented aliases into replacement public
contract.

This document uses BCP 14 keywords only for testable conformance requirements.
The key words "MUST", "MUST NOT", "SHOULD", and "MAY" are to be interpreted as
described in [RFC 2119](../../docs/reference/rfc2119.md) and
[RFC 8174](../../docs/reference/rfc8174.md) when, and only when, they appear in
all capitals.

## Background / Motivation

The current `/quality wizard` behavior is a read-only wayfinder: it inspects
local lifecycle state, classifies readiness, and recommends a next workflow.
That behavior is useful, but the word `wizard` is vague beside task-oriented
public modes such as `setup`, `evaluate`, and `update`. It also creates a naming
collision with a likely future use of "setup wizard" to describe guided setup.

The public contract should shrink rather than be renamed. Users should not have
to learn a fourth mode for "tell me what to do," and the skill should not expose
`status` or `next` as new public commands just to replace `wizard`. Ambiguous or
bare input still needs the current safety property: no evaluation, file
mutation, tooling update, or artifact write happens by surprise.

## Scope

Covered:

- Public `/quality` mode vocabulary and invocation examples.
- Safe read-only behavior for bare or ambiguous `/quality` requests.
- Runtime skill instructions and durable skill specs that currently name
  `wizard`.
- Setup handoff wording that currently routes to wizard.
- Quality-log and checklist wording that currently makes wizard the named
  reconciliation or inspection surface.
- Public docs that advertise `/quality wizard`, `/quality status`, or
  `/quality next`.

Deferred / non-goals:

- No setup-wizard behavior is specified here.
- No new public replacement mode such as `status`, `next`, `review`, or
  `orient`.
- No change to `qualitymd status` as a CLI support command.
- No change to the QUALITY.md format specification.
- No requirement to preserve `/quality wizard` as a compatibility alias beyond a
  possible deprecation response.

## Requirements

The `/quality` skill public modes **MUST** be limited to `setup`, `evaluate`,
and `update`.

> Rationale: this keeps the public skill surface task-oriented and removes the
> vague `wizard` mode without replacing it with another non-task command. — 0062

Recommendation follow-up **MUST** remain a non-mode workflow triggered by
natural requests to apply, act on, improve from, or hand off recommendations.

The public documentation **MUST NOT** advertise `/quality wizard`,
`/quality status`, `/quality next`, `/quality review model`, or
`/quality review history` as supported invocations.

Bare `/quality` **MUST** remain valid and read-only.

Bare or ambiguous `/quality` requests **MUST NOT** evaluate source, create
evaluation artifacts, mutate `QUALITY.md`, write the quality log, apply
recommendations, or update tooling.

When handling bare or ambiguous `/quality` requests, the skill **MAY** inspect
local lifecycle state and recommend a next action, but the recommended actions
**MUST** be limited to public workflows: `setup`, `evaluate`, `update`, or
recommendation follow-up.

Read-only orientation output **MUST NOT** identify itself as `wizard`, `status`,
`next`, or any other public replacement mode.

Run frames **MUST NOT** emit `Mode: wizard`.

Setup completion and getting-started guidance **MUST NOT** route users to
wizard. They should recommend the next public workflow directly.

Quality-log reconciliation guidance **MUST NOT** name wizard as the
reconciliation surface. Read-only orientation or model-review language may
surface drift, but quality-log writes remain limited to setup and confirmed
model-authoring or recommendation-apply workflows.

The top-10 QUALITY.md checklist **MUST NOT** be described as wizard-specific.
It may remain available to read-only orientation, setup, model review, or other
bounded model-lifecycle inspection.

If `/quality wizard` is handled for compatibility, it **MUST** produce a brief
deprecation/orientation response that points to public workflows and **MUST NOT**
be documented as part of the public contract.

## Acceptance Criteria

- `README.md` no longer mentions `/quality wizard`, `/quality status`, or
  `/quality next` as invocations.
- `install.md` no longer uses `/quality wizard` in bootstrap guidance.
- `docs/guides/use-quality-skill.md` and `docs/guides/index.md` describe the
  public skill surface without wizard.
- `skills/quality/SKILL.md` frontmatter and body no longer advertise wizard
  advice, `wizard` mode parsing, `wizard` hard rules, or `Mode: wizard`.
- `skills/quality/SKILL.md` lists only `setup`, `evaluate`, and `update` as
  public modes.
- Bare `/quality` remains read-only and does not default to `evaluate`.
- Ambiguous requests may produce orientation, but the response recommends only
  `setup`, `evaluate`, `update`, or recommendation follow-up.
- `/quality status`, `/quality next`, `/quality review model`, and
  `/quality review history` are not documented as public invocations.
- `skills/quality/modes/setup.md` and getting-started guidance no longer route
  to wizard after setup.
- Runtime and durable quality-log guidance no longer names wizard as the
  reconciliation surface.
- Runtime and durable top-10 checklist guidance no longer describes the
  checklist as wizard-specific.
- Durable skill specs remove or replace the wizard mode component so no public
  mode spec remains for wizard.
- Mode indexes and skill indexes contain no public wizard listing or broken
  wizard link.
- `/quality wizard`, if supported as a compatibility alias, is documented only
  as deprecated behavior in runtime guidance, not in public docs.
- No change is made to `SPECIFICATION.md`.

## Durable spec changes

### To add

None.

### To modify

- `specs/skills/quality-skill/quality-skill.md` - remove wizard from public mode
  vocabulary, default routing, invocation examples, and mode dispatch; preserve
  safe read-only handling for bare or ambiguous input according to the
  requirements above.
- `specs/skills/quality-skill/modes/setup.md` - stop routing setup completion to
  wizard and recommend public follow-on workflows directly.
- `specs/skills/quality-skill/guides/getting-started-md.md` - remove wizard as
  the normal next step after first useful model population.
- `specs/skills/quality-skill/guides/top-10-quality-md-checks-md.md` - make the
  checklist support bounded read-only orientation/model-review behavior rather
  than wizard specifically.
- `specs/skills/quality-skill/quality-log.md` - replace wizard-specific
  reconciliation wording with non-mode read-only orientation/model-review
  wording.
- `specs/skills/quality-skill/recommendation-follow-up.md` - remove wizard as a
  named router to recommendation follow-up.
- `specs/skills/quality-skill/index.md`,
  `specs/skills/quality-skill/modes/index.md`,
  `specs/skills/quality-skill/guides/index.md`, and relevant OKF logs - update
  listings and logs to match the removed wizard public contract.

### To rename

None.

### To delete

- `specs/skills/quality-skill/modes/wizard.md` - remove the public wizard mode
  component spec, unless the settled design replaces it with a non-public
  orientation component under a different path.
