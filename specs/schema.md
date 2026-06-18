---
type: Schema
title: specs/ concept types
description: Concept types used in the specs/ OKF bundle.
types:
  - name: Schema
    description: This file — the registry of concept types used in a bundle.
  - name: Functional Specification
    description: A spec for what the qualitymd tooling or runtime artifact contracts must do — the surface as a whole, a shared contract, or a single subcommand.
  - name: Example Model
    description: A worked reference instance of a QUALITY.md model, reproduced so an example evaluation's findings are traceable to declared requirements.
---

This bundle's recommended concept-type vocabulary lives in the `types`
frontmatter above. `type` is a free-form OKF string and consumers tolerate
unknown values, so it is a recommendation, not a closed schema: reuse a listed
type when it fits, or coin a new descriptive one and add it here.
