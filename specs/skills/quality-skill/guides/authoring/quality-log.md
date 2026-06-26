---
type: Functional Specification
title: QUALITY.md quality-log authoring guide
description: Contract for meaningful model-change and quality-log authoring guidance.
tags: [skill, quality, guide, authoring]
timestamp: 2026-06-24T00:00:00Z
---

# QUALITY.md quality-log authoring guide

This spec governs the runtime
[`skills/quality/guides/authoring/quality-log.md`](../../../../../skills/quality/guides/authoring/quality-log.md)
guide.

This document uses BCP 14 keywords only for testable conformance requirements.
The key words "MUST", "MUST NOT", and "SHOULD" are to be interpreted as
described in [RFC 2119](../../../../../docs/reference/rfc2119.md) and
[RFC 8174](../../../../../docs/reference/rfc8174.md) when, and only when, they
appear in all capitals.

The guide **MUST** cover when to update QUALITY.md, model drift vs root-area
quality, body/frontmatter synchronization, defect-backlog boundaries,
recalibration vs drift correction, criteria gaming, missing Requirements surfaced
by findings, and whether the Requirement set still delivers the body's Needs.

The guide **MUST** define what counts as a meaningful quality-log entry while
leaving the quality-log format contract in `SKILL.md`.
