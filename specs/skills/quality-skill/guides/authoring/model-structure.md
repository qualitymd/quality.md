---
type: Functional Specification
title: QUALITY.md model-structure authoring guide
description: Contract for area, source, decomposition, traceability, constituent, and use-context guidance.
tags: [skill, quality, guide, authoring]
timestamp: 2026-06-24T00:00:00Z
---

# QUALITY.md model-structure authoring guide

This spec governs the runtime
[`skills/quality/guides/authoring/model-structure.md`](../../../../../skills/quality/guides/authoring/model-structure.md)
guide.

This document uses BCP 14 keywords only for testable conformance requirements.
The key words "MUST", "MUST NOT", and "SHOULD" are to be interpreted as
described in [RFC 2119](../../../../../docs/reference/rfc2119.md) and
[RFC 8174](../../../../../docs/reference/rfc8174.md) when, and only when, they
appear in all capitals.

The guide **MUST** cover areas, sources, grouping areas, child area split rules,
source inheritance, authored/inspectable entities, primary-subject/collection/
composite decomposition, traceability graph edges, entities as both areas and
assessment references, normative artifacts, high-leverage concerns, domain
constituent kinds, stewardship concern projection boundaries, and recurring
use-context constituents.

The guide **MUST** route detailed agent-harness area guidance to
`agent-harness.md` while preserving the self-check area guidance and the general
use-context constituent rule. It **MUST NOT** define the agent harness as
instructions alone; it **MUST** refer to the harness system or project-owned
harness artifacts consistently with `agent-harness.md`.
