---
type: Schema
title: specs/ concept types
description: Concept types used in the specs/ OKF bundle.
types:
  - name: Schema
    description: This file — the registry of concept types used in a bundle.
  - name: Functional Specification
    description: A spec for what the qualitymd tooling must do — the surface as a whole or a single subcommand.
  - name: Example Model
    description: A worked reference instance of a QUALITY.md model, reproduced so an example evaluation's findings are traceable to declared requirements.
  - name: Evaluation Design
    description: A worked reference instance of the design artifact an evaluation records — its resolved parameters and the model snapshot it is bound to.
  - name: Evaluation Plan
    description: A worked reference instance of the plan artifact an evaluation records — how the in-scope source is covered at the chosen effort.
  - name: Assessment Record
    description: A worked reference instance of a single write-once assessment record — one in-scope requirement's findings, inferred rating, and rationale — that the report rolls up.
  - name: Analysis Record
    description: A worked reference instance of a single write-once analysis record — one target node's inferred local, aggregate, and factor ratings (the roll-up), citing the assessment and child analysis records it derives from.
  - name: Evaluation Report
    description: A worked reference instance of the Evaluation Report a skill produces, used to make a spec's reporting contract concrete.
  - name: Recommendation
    description: A worked reference instance of a single triageable recommendation artifact emitted alongside an Evaluation Report.
---

This bundle's recommended concept-type vocabulary lives in the `types`
frontmatter above. `type` is a free-form OKF string and consumers tolerate
unknown values, so it is a recommendation, not a closed schema: reuse a listed
type when it fits, or coin a new descriptive one and add it here.
