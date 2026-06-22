---
type: Functional Specification
title: Recommendation follow-up guide
description: Contract for the skill's runtime guide for applying or handing off evaluation recommendations.
tags: [skill, quality, guide, recommendation]
timestamp: 2026-06-22T00:00:00Z
---

# Recommendation follow-up guide

This spec governs the **Recommendation Follow-Up** guide the
[`/quality` skill](../quality-skill.md) ships at
[`skills/quality/guides/recommendation-follow-up.md`](../../../../skills/quality/guides/recommendation-follow-up.md).
The guide is the runtime procedure for acting on evaluation recommendations.

This document uses BCP 14 keywords only for testable conformance requirements.
The key words "MUST" and "MUST NOT" are to be interpreted as described in
[RFC 2119](../../../../docs/reference/rfc2119.md) and
[RFC 8174](../../../../docs/reference/rfc8174.md) when, and only when, they
appear in all capitals.

## Requirements

The guide **MUST** state that recommendation follow-up is not a `/quality` mode.

The guide **MUST** offer only two explicit outcomes: apply a confirmed
recommendation option now, or hand off the recommendation to an issue tracker.

The guide **MUST NOT** present defer, skip, or keep open as formal follow-up
options.

The guide **MUST** require explicit confirmation before editing evaluated
source, editing `QUALITY.md`, writing the quality log, or creating an external
issue.

The guide **MUST** state that issue-tracker handoff alone does not edit
evaluated source, `QUALITY.md`, or `.quality/log/`.

The guide **MUST** name the issue-ready content fields needed for handoff:
recommendation identity, source run, affected model target, current rating when
known, target or done criterion, evidence locators, suggested option,
verification path, and report/recommendation artifact paths.
