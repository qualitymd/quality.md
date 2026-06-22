---
type: Functional Specification
title: /quality setup
description: Behavioral component spec for bootstrapping and first-populating a QUALITY.md model through the /quality skill.
tags: [skill, quality, mode, setup]
timestamp: 2026-06-22T00:00:00Z
---

# /quality setup

`setup` is the `/quality` skill mode that bootstraps a missing model and guides
the first useful population of `QUALITY.md`. It implements the shared contracts
in the parent [/quality skill](../quality-skill.md) spec and owns only the
setup-specific behavior below.

The runtime procedure lives at
[`skills/quality/modes/setup.md`](../../../../skills/quality/modes/setup.md).

This document uses BCP 14 keywords only for testable conformance requirements.
The key words "MUST", "MUST NOT", and "SHOULD" are to be interpreted as
described in [RFC 2119](../../../../docs/reference/rfc2119.md) and
[RFC 8174](../../../../docs/reference/rfc8174.md) when, and only when, they
appear in all capitals.

## Purpose and routing

`setup` is selected when no model file is present, when the user explicitly asks
to create or initialize a QUALITY.md file, or when wizard routes to bootstrap or
first-population work.

The mode's purpose is to produce a valid starter model and immediately begin
turning it into a useful project-specific model. It is not an evaluation mode
and does not rate evaluated source.

## Mutation surface and artifacts

`setup` may mutate:

- the target `QUALITY.md` model file;
- the quality log under `.quality/log/` after guided first population; and
- no evaluated source other than the model file itself.

`setup` **MUST** drive `qualitymd init` for deterministic scaffolding when the
model file is absent, then run `qualitymd lint`. It **MUST NOT** reimplement
scaffolding, validation, CLI installation tooling, or source-driven authoring
judgment.

After creating or validating the skeleton, `setup` **MUST** read the authoring
guide and getting-started guide before guided population. Guided population
**MUST** address the Markdown body's Overview, Scope, Needs, and Risks,
including each section's unknowns, open questions, and any material support that
is not agent-accessible.

Guided population **SHOULD** propose project-specific factors and requirements
to replace placeholders.

Guided population **SHOULD** include a `quality-md` area that evaluates the
`QUALITY.md` artifact itself against the active authoring guide unless the user
declines or the model file is not in the root area it governs. The area
**SHOULD** use the key `quality-md`, a title of the form `<Root Title>
QUALITY.md`, an area `description`, and an explicit path-based `source` such as
`./QUALITY.md`. It **MUST NOT** use prose aliases such as `(this file)` for
`source`.

When setup adds that area, it **SHOULD** include concise YAML comments that
distinguish the area `source` from the requirement `assessment`. It **SHOULD**
use one area-level requirement with `factors` when the active authoring guide
defines one coherent judgment across multiple factors.

## Stop conditions

`setup` **MUST** stop before CLI-dependent work when the `qualitymd` CLI is
missing, outside the released-install SemVer range declared by the skill, or a
local development build lacks required commands.

`setup` **MUST** stop before guided population when `qualitymd lint` reports
errors that make the model structurally invalid.

## Completion criteria

`setup` is complete when the target model exists, passes lint, has received
guided first-population work or a clearly reported user-deferred population step,
and any meaningful model creation or model-shape change has a corresponding
quality-log entry. Follow-on routing belongs to [`wizard`](wizard.md).
