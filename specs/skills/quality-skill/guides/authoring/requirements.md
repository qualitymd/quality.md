---
type: Functional Specification
title: QUALITY.md requirement authoring guide
description: Contract for requirement names and titles, assessments, factor connections, and rating overrides.
tags: [skill, quality, guide, authoring]
timestamp: 2026-06-24T00:00:00Z
---

# QUALITY.md requirement authoring guide

This spec governs the runtime
[`skills/quality/guides/authoring/requirements.md`](../../../../../skills/quality/guides/authoring/requirements.md)
guide.

This document uses BCP 14 keywords only for testable conformance requirements.
The key words "MUST", "MUST NOT", and "SHOULD" are to be interpreted as
described in [RFC 2119](../../../../../docs/reference/rfc2119.md) and
[RFC 8174](../../../../../docs/reference/rfc8174.md) when, and only when, they
appear in all capitals.

The guide **MUST** cover requirements as assessable quality expectations, stable
requirement names, natural-language requirement titles, primary and secondary
factor connections, area-local `factors`, concrete observable properties, scale
and meter preflight, ratable statements, behavioral trigger conditions,
verifiable assessments, risk-weighted detail, exactly one assessment per
requirement, external assessment references, one referenced assessment feeding
several same-area factors, splitting by assessable claim rather than by factor,
conformance vs fitness separation, rating overrides, and closing validation of
the requirement set. The guide **MUST NOT** tell agents that `factors` entries
can resolve to ancestor area factors.
