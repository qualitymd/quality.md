---
date: 2026-06-26
kind: assessment-change
target: quality-md/the-model-follows-the-authoring-guide-family
---

Changed the model self-check from purely inferential body review toward a named
drift detector. Earlier evaluations relied on a reviewer noticing whether body
unknowns, open questions, and review provenance had fallen behind the
frontmatter. That remained useful judgment, but it was too easy to miss after a
model expansion.

The revised assessment still reviews the authoring guide family, but it now
expects a detector that compares factor and requirement changes against body
sections and review lines. That sensor does not replace judgment; it gives the
quality loop a repeatable signal for the stale-context failure mode.
