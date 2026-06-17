---
type: Change
title: Example change
description: Placeholder change demonstrating the Change concept shape. Replace with a real change.
status: Draft
tags: [example]
timestamp: 2026-06-16T00:00:00Z
---

# Example change

> 🚧 **Placeholder.** This change exists to show the intended shape of a `Change`
> concept. Replace it with a real change, or delete it once the bundle has real
> content.

A **Change** is a unit of incremental work on the repo. This parent concept
captures the *why* and the *status*; the detail lives in its children:

- [Functional spec](0001-example-change/spec.md) — what the change must do.
- [Design doc](0001-example-change/design.md) — how it's built, and why.

A change that needs no design doc omits `design.md`.

## Motivation

Why this change is worth making — the problem it solves or the value it adds.

## Scope

What's covered, and what's intentionally deferred.

## Status

`Draft`. See the [status lifecycle](index.md#status-lifecycle). When this change
reaches **Done**, move it and its [child folder](0001-example-change/) into
[`archive/`](archive/).
