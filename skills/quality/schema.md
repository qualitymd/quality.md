---
type: Runtime Schema
title: /quality runtime concept types
description: Concept types used in the /quality runtime skill bundle.
types:
  - name: Runtime Index
    description: A progressive-disclosure listing for a runtime skill directory.
  - name: Runtime Log
    description: An update log for runtime skill content.
  - name: Runtime Schema
    description: This file — the registry of concept types used in the runtime skill bundle.
  - name: Runtime Guide
    description: A runtime guide resource read by the /quality skill.
  - name: Runtime Workflow
    description: Runtime workflow instructions for a /quality public workflow.
  - name: Runtime Resource
    description: A bundled runtime reference resource.
---

This bundle's recommended concept-type vocabulary lives in the `types`
frontmatter above. `type` is a free-form OKF string and consumers tolerate
unknown values, so it is a recommendation, not a closed schema.
