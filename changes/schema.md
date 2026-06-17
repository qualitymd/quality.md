---
type: Schema
title: changes/ concept types
description: Concept types used in the changes/ OKF bundle.
types:
  - name: Schema
    description: This file — the registry of concept types used in a bundle.
  - name: Change
    description: A unit of incremental work on the repo — a parent concept that records the motivation and status and links to its functional spec and design doc.
  - name: Functional Specification
    description: What a change must do — its requirements, not the implementation. The same type the specs/ bundle uses.
  - name: Design Doc
    description: How a change is built — the technical approach and the rationale behind it.
---

This bundle's recommended concept-type vocabulary lives in the `types`
frontmatter above. `type` is a free-form OKF string and consumers tolerate
unknown values, so it is a recommendation, not a closed schema: reuse a listed
type when it fits, or coin a new descriptive one and add it here.
