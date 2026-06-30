---
type: Change Case
title: Durable spec rationale
description: Make functional specs capture durable rationale — big-picture and per-requirement — so the lessons that motivate requirements stop dying in archived change folders.
status: Done
tags: [specs, docs, contributing]
timestamp: 2026-06-18T00:00:00Z
---

# Durable spec rationale

A unit of work on the contributor guides: teach the enduring
[`specs/`](../../specs/index.md) bundle to carry the _why_ behind its requirements,
not just the requirements themselves. Detail lives in the children:

- [Functional spec](0025-durable-spec-rationale/spec.md) - what the three durable
  guides must now require and teach.
- [Design doc](0025-durable-spec-rationale/design.md) - why a two-layer,
  in-spec rationale convention, and the alternatives rejected.

## Motivation

A [Change Case](../../docs/guides/work-with-change-cases.md) states a _delta_
and is archived once it lands; the enduring `specs/` bundle carries the
_cumulative_ source of truth forward. Today that hand-off keeps the _what_ but
drops the _why_: a case's motivation and its design doc's hard-won rationale go
into [`archive/`](.) with the case, while the durable spec inherits only the bare
requirement. A later editor reading the spec sees the rule but not the
failure-mode that produced it — so settled lessons get re-litigated, and rules
that look arbitrary get "simplified" away, reintroducing the very bug they fixed.

The functional-spec guide currently makes this worse: it caps _all_ rationale at
a clause or a `Note:` and frames any longer _why_ as a sign the _what_ has been
buried. That guidance protects skimmability but leaves no durable home for
genuine intent, so rationale that deserves to survive has nowhere to live in the
spec and dies in the archive instead.

## Scope

Covered: the contributor guides only. Add a two-layer rationale convention to the
functional-spec guide (a spec-level **Background / Motivation** section for the
big-picture _why_, and a per-requirement subordinate annotation for the
fine-grained _why_); rewrite the guide's "Motivation in asides" convention and
its rationale-related smells to match; teach the changes-workflow guide to
_absorb_ a landing case's enduring _why_ into the durable spec, not just its
functional delta; and note in the design-doc guide that durable rationale is
_promoted_ into the spec when the case lands.

Deferred: any change to the `specs/` bundle contents themselves (the new
convention is applied to existing specs as they are next touched, not in a sweep
here), any tooling or linter that detects un-annotated requirements, and any
schema change — the three existing concept types
([`schema.md`](../schema.md)) already cover this work.

## Affected specs & docs

Created or updated during `In-Progress`, before this change reaches
`In-Review`. This change's durable artifacts are themselves contributor guides:

- [x] [`docs/guides/write-functional-specs.md`](../../docs/guides/write-functional-specs.md) - add the **Background / Motivation** shape entry and the per-requirement
      annotation convention; rewrite the "Motivation in asides" bullet as the
      two-whys split; refine the rationale smells.
- [x] [`docs/guides/work-with-change-cases.md`](../../docs/guides/work-with-change-cases.md) - require absorbing the case's enduring _why_ (its motivation and the
      design doc's durable rationale) into the durable spec's Background and
      per-requirement annotations, tied to the existing **Before setting
      In-Review** gate.
- [x] [`docs/guides/write-design-docs.md`](../../docs/guides/write-design-docs.md) - note that the design doc's enduring rationale is promoted into the
      functional spec on landing, while the design doc remains the archived
      record of alternatives and trade-offs.

## Status

`Done`. Implemented and archived after teaching the three durable contributor guides to keep rationale in the spec: a Background/Motivation shape entry and a per-requirement `Rationale:` annotation convention.
