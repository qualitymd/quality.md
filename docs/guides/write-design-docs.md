---
type: How-to Guide
title: Writing design docs
description: How to write a design doc for a piece of qualitymd.
tags: [design, contributing]
timestamp: 2026-06-18T00:00:00Z
---

# Writing design docs

A **design doc** describes *how* a piece of `qualitymd` is built — the technical
approach behind a [functional spec](write-functional-specs.md). Where a spec says what
must hold, a design doc says how the code makes it so, and why that way.

## Shape

Keep it short and decision-focused:

- **Context** — what we're building and the spec or problem it answers (link it).
- **Approach** — the design: key components, data flow, and the shape of the
  solution. Diagrams and code sketches welcome.
- **Spec response** — how the approach satisfies the major requirement groups,
  plus any requirement whose verification depends on a non-obvious design
  choice. Link to the spec instead of restating its normative text.
- **Alternatives** — options considered and why they lost. This is the part worth
  writing — it captures the reasoning a reader would otherwise have to reconstruct.
- **Trade-offs & risks** — what this choice costs, and what could go wrong.
- **Open questions** — unresolved decisions, so they're visible rather than buried.

## Conventions

- **Design, not spec.** Behavior and requirements belong in the
  [functional spec](write-functional-specs.md) — link to it, don't restate it.
- **Answer the spec.** A design doc should make it clear that the chosen approach
  can satisfy the spec's requirement groups and verification-sensitive edge
  cases. If a requirement is hard to satisfy, partially satisfied, or relies on a
  risky assumption, name that as a trade-off or open question rather than
  weakening the spec by implication.
- **Durable impact is the spec's job.** Which durable specs the change rewrites —
  and what they must say — go in the functional spec's
  [Durable spec changes](write-functional-specs.md#durable-spec-changes) section;
  code and durable docs go in the change case's parent **Affected artifacts**
  index. A
  design doc does **not** carry its own list of durable specs or docs to edit. It
  says *how* the code delivers the change; the *what* stays out of it.
- **Record the why.** The durable value is the rationale and rejected
  alternatives, not a re-description of the final code. When the Change Case lands,
  that value is *promoted*, not abandoned: the enduring rationale is lifted into
  the [functional spec](write-functional-specs.md)'s **Background / motivation**
  and per-requirement annotations, while this doc stays in the archive as the
  fuller record of alternatives and trade-offs.
- **Match the scope.** A one-paragraph note for a small change; a fuller doc for
  anything cross-cutting. Don't over-document.
- **Use sentence-case headings.** Design doc headings use sentence case while
  preserving proper nouns, formal type names, and cited source casing.
- **Let it age gracefully.** A design doc reflects a decision at a point in time;
  supersede it with a new one rather than silently rewriting history.
