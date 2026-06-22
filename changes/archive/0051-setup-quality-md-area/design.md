---
type: Design Doc
title: Setup quality-md Area — design doc
description: Design for adding the setup-authored quality-md Area through skill guidance rather than CLI scaffolding.
tags: [skill, setup, quality-model]
timestamp: 2026-06-22T00:00:00Z
---

# Setup quality-md Area — design doc

Design behind [Setup quality-md Area](../0051-setup-quality-md-area.md) and its
[functional spec](spec.md).

## Context

The change adds a setup-time modeling pattern, not a deterministic file format or
CLI primitive. `/quality setup` already runs `qualitymd init` for the valid
skeleton, then reads the authoring and getting-started guides before applying
project-specific judgment. The `quality-md` Area belongs in that second phase:
the skill can name the root artifact, choose useful Factors, cite the active
authoring guide, and explain the pattern to first-time users.

## Approach

Keep `qualitymd init` and `internal/scaffold/skeleton.md` unchanged. The CLI
continues to emit a generic, valid starter that works without knowing which
agent skill is installed or where its guide files live.

Update the setup mode instructions so guided population normally proposes or
inserts a `quality-md` Area after lint succeeds. The Area uses the normal model
shape:

- key: `quality-md`;
- title: `<Root Title> QUALITY.md`;
- description: one sentence identifying the concrete model artifact;
- source: a path selector for the model file, normally `./QUALITY.md`;
- factors: project-model qualities such as decision suitability, credibility,
  assessability, traceability, and maintainability;
- requirement: one Area-level Requirement referencing the active authoring guide
  and listing the relevant Factors through `factors`.

Add concise YAML comments with the Area when setup writes it. The comments
explain that `source` points at the evaluated `QUALITY.md` file, while
`assessment` names the guide used to judge it. That avoids introducing a prose
source alias such as `(this file)` while still making the self-referential-looking
pattern understandable.

Update the authoring guide rather than only the setup prompt for the two reusable
lessons this pattern exposes:

- Factor names should be quality attributes the Area exhibits, not practices or
  workflows.
- A single guide/spec/checklist can define one Requirement that contributes to
  several Factor roll-ups through `factors`; authors should split only when the
  claims could rate differently.

## Alternatives

**Put the Area in `qualitymd init`.** Rejected. The CLI scaffold is deterministic
and context-light, while this Area depends on authoring judgment, the root title,
and an agent-accessible guide path. Baking it into the scaffold would either
hard-code an installation-specific path or emit generic placeholders that setup
would still need to repair.

**Use `(this file)` as `source`.** Rejected. `source` is a machine-resolvable
selector. A prose alias would require special-case evaluator behavior and weaken
the existing path-based source model.

**Repeat the authoring-guide assessment under each Factor.** Rejected. Repeating
the same assessment fragments one coherent judgment and makes the model noisier.
One Area-level Requirement with multiple Factor references preserves a single
assessment while still contributing to each relevant roll-up.

**Use abstract names like `quality-model` or `lifecycle-stewardship`.** Rejected.
The Area should identify the concrete artifact (`quality-md`), and Factors should
name qualities the artifact can exhibit to a degree, not stewardship practices.

## Trade-offs & risks

The setup prompt will carry more modeling opinion. That is acceptable because
setup already owns guided population after the CLI-created skeleton; the
functional spec keeps the behavior a `SHOULD` so users and unusual repository
layouts can opt out.

The active authoring-guide path may vary by runtime. The skill should cite the
most agent-accessible locator available in the current environment rather than
pretending the CLI can know it.

The `quality-md` Area adds one more Area to first-run models. The comments and
single-requirement shape keep the extra surface small while making model quality
visible from the beginning.

## Open questions

None.
