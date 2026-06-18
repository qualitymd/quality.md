---
type: Functional Specification
title: Authoring guide replaces meta-model workflow
description: Requirements for replacing the quality meta-model workflow with a QUALITY.md authoring guide.
tags: [skill, cli, evaluation, authoring-guide]
timestamp: 2026-06-18T00:00:00Z
---

# Authoring guide replaces meta-model workflow

The key words **MUST**, **MUST NOT**, **SHOULD**, and **MAY** are to be
interpreted as described in IETF RFC 2119.

## Background / Motivation

The quality meta-model encoded model-review criteria as another `QUALITY.md`,
but the artifact users need first is author guidance: how to make good target,
factor, requirement, assessment, rating, and rationale choices. Treating
authoring advice as a bundled diagnostic model made the public CLI carry a
model-altitude workflow that is harder to explain and less useful than a direct
guide.

## Requirements

The skill-facing reference formerly named `quality-meta-model.md` **MUST** be
renamed to `quality-md-guide.md` and rewritten as a prose authoring guide.

The guide **MUST** explain how to author a useful `QUALITY.md`: choosing scope,
shaping the target tree, selecting factors, writing requirements, defining
assessments, designing rating scales, writing body rationale, and reviewing the
file before evaluation.

The installable `/quality` skill **MUST** read the authoring guide for setup,
wizard, authoring, and improvement work where practical authoring judgment is
needed.

The `/quality` skill **MUST NOT** offer or depend on model-altitude evaluation,
the bundled `quality-meta-model`, or `qualitymd models`.

The public CLI **MUST NOT** expose `qualitymd models`.

`qualitymd evaluation create-run` **MUST NOT** create model-altitude runs.
Run creation is for subject evaluation only.

Durable specs and user docs **MUST** stop presenting the meta-model or
model-altitude evaluation as current behavior.

Archived changes and already-written evaluation artifacts **MAY** keep their
historical references.
