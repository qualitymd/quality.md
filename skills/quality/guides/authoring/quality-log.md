---
type: Runtime Guide
title: Authoring Quality-Log Changes
description: Guidance for deciding when QUALITY.md model changes require quality-log entries.
tags: [quality, authoring, guide]
---

# Authoring Quality-Log Changes

Read this when:

- changing what a QUALITY.md model is or how it judges;
- applying a recommendation that changes QUALITY.md;
- deciding whether a quality-log entry is meaningful.

Depends on:

- `../authoring.md`

---

## When to update QUALITY.md

A `QUALITY.md` is expected to evolve when evaluation or authoring reveals that
the model no longer represents what quality means for the root area. The
`quality-md` area is evaluated and rolled up like any other in-scope area; the
distinct follow-up behavior is mutation history, not evaluation semantics. When a
confirmed change meaningfully alters the model, record why it changed in the
quality log.

- **Do** revise when a discovery changes the context or content of the
  evaluation — a new factor that matters, a requirement whose assessment changed,
  a scope that shifted.
- **Do** update the model when an evaluation finding shows the model no longer
  reflects the root area's real scope, risks, or decision needs. *That is model
  drift, not merely a weak root area rating.*
- **Do** keep the body current with the frontmatter. *A model whose body no
  longer explains its factors misleads the next evaluator.*
- **Avoid** using `QUALITY.md` as a defect backlog. *Evaluated-source defects belong in
  the root area's normal planning system unless they also change what quality
  means or how it should be assessed.*
- **Do** distinguish *recalibration* (a deliberate decision to reset a criterion
  because you have learned what is achievable) from *drift* (the model silently
  falling out of step with the root area). *Recalibration is healthy: after a
  breakthrough, raise `minimum` so the new floor sticks; after hitting a real
  constraint, lower a `target` consciously and say why in the body.*
- **Avoid** sharpening criteria only to keep ratings green. *The review's job is to
  keep the rubric valid, not passing; locked baselines and an honest "not assessed"
  guard against gaming it.*
- **Do** treat a finding that no existing requirement anticipated as a signal to
  add a requirement or factor. *A real weakness your model could not express is the
  strongest evidence the model is incomplete — the model improves by being used,
  not only by being authored.*
- **Do** periodically check that satisfying the requirement set would actually
  deliver the body's Needs. *If the model can be fully green while the root area
  still fails its purpose, the requirement set is incomplete — that is model drift,
  not a strong root area.*

### Logging a model change

When a confirmed recommendation follow-up or direct model-authoring edit actually
changes the model, record it in the **quality log**. Its format contract lives in
[`SKILL.md`](../../SKILL.md); this guide covers what counts as meaningful:

- **Do** log a change that alters what the model *is* or *how it judges*: adding,
  removing, or renaming an Area, Factor, or Requirement; changing the rating
  scale, a criterion, or a relative weight; shifting scope; changing the apex or
  required margin; or applying an evaluation recommendation.
- **Do** state whether a criterion move is deliberate *recalibration* or a *drift
  correction*, and cross-link the evaluation run and recommendation behind it when
  the change came from one.
- **Do** write **one entry per coherent change** — a confirmed recommendation
  apply, a direct model-authoring change, or the initial population — not one per
  field touched. *The unit of record is the decision, not the edit.*
- **Avoid** logging Markdown-body wording, typo, or formatting changes, or
  evaluated-source fixes that leave the model unchanged. *Those are not model
  changes; git already records them, and logging them turns a curated timeline
  into noise.*
- **Avoid** treating the log as a second evaluation record or a defect backlog.
  *It references evaluation runs; it never copies them.*
