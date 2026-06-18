---
type: Change Case
title: Example change
description: Placeholder change demonstrating the Change Case concept shape, retired as a reference template.
status: Done
tags: [example]
timestamp: 2026-06-17T00:00:00Z
---

# Example change

> 🗂️ **Reference template.** This placeholder showed the intended shape of a
> `Change Case` concept. It's archived now that the bundle has real cases — copy
> it and its folder as the starting point for a new case.

A **Change Case** is a unit of incremental work on the repo. This parent concept
captures the *why* and the *status*; the detail lives in its children:

- [Functional spec](0001-example-change/spec.md) — what the case must do.
- [Design doc](0001-example-change/design.md) — how it's built, and why.

A case that needs no design doc omits `design.md`.

## Motivation

Why this case is worth making — the problem it solves or the value it adds.

## Scope

What's covered, and what's intentionally deferred.

## Affected specs & docs

The durable specs and docs this case creates or updates — the enduring
[`specs/`](../../specs/index.md) bundle, the repository-root
[`SPECIFICATION.md`](../../SPECIFICATION.md), the [`README.md`](../../README.md),
and the [`docs/`](../../docs/index.md) guides. Decide this up front, alongside the
motivation and scope. List each artifact with what changes; an empty list must
read as a deliberate "no durable changes," not an oversight. Every listed
artifact is created or updated **before** this change reaches `Done`.

- [ ] `path/to/spec.md` — *what changes, and why*.

## Status

`Done`. See the [status lifecycle](../index.md#status-lifecycle). Retired as a
reference template once the bundle had real changes to follow.
