---
type: Functional Specification
title: Agent-harness Area authoring guide
description: Contract for modeling the agent harness as a QUALITY.md Area/constituent.
tags: [skill, quality, guide, authoring]
timestamp: 2026-06-24T00:00:00Z
---

# Agent-harness Area authoring guide

This spec governs the runtime
[`skills/quality/guides/authoring/agent-harness.md`](../../../../../skills/quality/guides/authoring/agent-harness.md)
guide.

This document uses BCP 14 keywords only for testable conformance requirements.
The key words "MUST", "MUST NOT", and "SHOULD" are to be interpreted as
described in [RFC 2119](../../../../../docs/reference/rfc2119.md) and
[RFC 8174](../../../../../docs/reference/rfc8174.md) when, and only when, they
appear in all capitals.

The guide **MUST** cover the agent harness as a recurring use-context constituent
and partly normative Area. It **MUST** define the agent harness holistically as
the whole engineered system around the model - everything that is not the model
itself, including code, configuration, and execution logic. It **MUST** name
feedforward controls that steer the agent before action and feedback controls
that catch and correct it after. It **MUST NOT** define the harness as
instructions, steering prose, or guidance files alone.

The guide **MUST** derive the agent-harness Area as the checked-in,
project-owned governing-artifacts projection of that holistic harness: agent
entry points, guidance files, skills, prompts, owned hooks, tool/MCP definitions,
sandbox or permission policy, orchestration config, and similar agent-governing
controls. Project-owned runtime harness machinery **MUST** be surfaced in this
Area or, when large enough to warrant distinct factors, in its own Area; it
**MUST NOT** be silently folded into prose instructions or dropped.

The guide **MUST** preserve the dual Area/assessment-reference role and the
projection boundary against Agent Harnessability and the agent audience. It
**MUST** state that Agent Harnessability rates the capability of equipping an
agent, the agent-harness Area rates checked-in project-owned governing artifacts,
and verification corpus or runtime environment artifacts belong to tests or
operations constituents unless their primary job is to govern or equip agent work.

The guide **MUST** include a scoping decision rule for mixed artifacts: rate an
artifact in the agent-harness Area when its primary job is to govern or equip the
agent's work and the project owns it; cede it to a domain constituent when it is
primarily a product artifact the agent merely also uses; when one artifact does
both, rate the agent-governing quality in the harness Area and cross-reference the
domain constituent under the no-double-count rule.

The guide **MUST** preserve the domain-agnostic factor-family prompt and
served-domain guardrail while expanding requirement shapes across feedforward
guidance, feedback routing and traces, and owned controls. Requirements **MUST**
rate the artifact's own quality, not the capability it confers.
