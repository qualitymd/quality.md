---
type: How-to Guide
title: Writing design docs
description: How to write a design doc for a piece of qualitymd.
tags: [design, contributing]
timestamp: 2026-06-16T00:00:00Z
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
- **Alternatives** — options considered and why they lost. This is the part worth
  writing — it captures the reasoning a reader would otherwise have to reconstruct.
- **Trade-offs & risks** — what this choice costs, and what could go wrong.
- **Open questions** — unresolved decisions, so they're visible rather than buried.

## Conventions

- **Design, not spec.** Behavior and requirements belong in the
  [functional spec](write-functional-specs.md) — link to it, don't restate it.
- **Record the why.** The durable value is the rationale and rejected
  alternatives, not a re-description of the final code. When the Change Case lands,
  that value is *promoted*, not abandoned: the enduring rationale is lifted into
  the [functional spec](write-functional-specs.md)'s **Background / Motivation**
  and per-requirement annotations, while this doc stays in the archive as the
  fuller record of alternatives and trade-offs.
- **Match the scope.** A one-paragraph note for a small change; a fuller doc for
  anything cross-cutting. Don't over-document.
- **Let it age gracefully.** A design doc reflects a decision at a point in time;
  supersede it with a new one rather than silently rewriting history.
