---
type: Schema
title: docs/ concept types
description: Concept types used in the docs/ OKF bundle — the four Diátaxis modes.
types:
  - name: Schema
    description: This file — the registry of concept types used in a bundle.
  - name: Tutorial
    description: Learning-oriented lesson that walks a newcomer through a task end to end. The Diátaxis tutorial mode.
  - name: How-to Guide
    description: Task-oriented steps for accomplishing a specific goal. The Diátaxis how-to mode.
  - name: Reference
    description: Information-oriented, lookup material — factual and complete. The Diátaxis reference mode.
  - name: Explanation
    description: Understanding-oriented discussion of background, concepts, and rationale. The Diátaxis explanation mode.
---

The recommended concept-type vocabulary lives in the `types` frontmatter above:
the four [Diátaxis](https://diataxis.fr/) modes, plus `Schema`. Each concept's
`type` names the mode it belongs to, and concepts of a mode are grouped in the
matching subfolder (`guides/`, `reference/`, and — as they appear — `tutorials/`
and `explanation/`).

`type` is a free-form OKF string and consumers tolerate unknown values, so this
is a recommendation, not a closed schema: reuse a listed type when it fits, or
coin a new descriptive one and add it here.
