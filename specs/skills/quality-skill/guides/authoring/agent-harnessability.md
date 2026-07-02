---
type: Functional Specification
title: Agent harnessability authoring guide
description: Contract for modeling agent harnessability as a model-wide factor.
tags: [skill, quality, guide, authoring]
timestamp: 2026-06-24T00:00:00Z
---

# Agent harnessability authoring guide

This spec governs the runtime
[`skills/quality/guides/authoring/agent-harnessability.md`](../../../../../skills/quality/guides/authoring/agent-harnessability.md)
guide.

This document uses BCP 14 keywords only for testable conformance requirements.
The key words "MUST", "MUST NOT", and "SHOULD" are to be interpreted as
described in [RFC 2119](../../../../../docs/reference/rfc2119.md) and
[RFC 8174](../../../../../docs/reference/rfc8174.md) when, and only when, they
appear in all capitals.

The guide **MUST** cover agent harnessability as the model-wide factor projection
of the agent-collaboration concern for agent-collaborated composite roots, using
`agent-harnessability` as the recommended key. It **MUST** keep the factor distinct
from the agent-harness area and agent audience. That boundary **MUST** present
the agent-harness area as the checked-in, project-owned governing-artifacts
projection, not as the whole equipping capability and not as instructions alone.

The guide **MUST** present agent harnessability as an umbrella carrying no direct
requirements, rated by rolling up independently assessable sub-factors:
`agent-accessibility`, `task-specifiability`, `agent-operability`, `continuity`,
`self-verifiability`, `enforcement-of-standards`, and `containment-of-action`. It
**MUST** preserve boundary guidance and the rule that harness thinness is rating
evidence rather than an omission reason. It **MUST** treat an existing
`harnessability` factor as stale legacy naming, not current agent harnessability
coverage, and route authoring work toward the current `agent-harnessability`
shape.
