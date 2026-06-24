---
type: Design Doc
title: Normalize QUALITY.md self-check roll-up - design
description: Why normal area semantics replace the special learn-loop roll-up exception for the `quality-md` self-check area.
tags: [skill, authoring, evaluation, roll-up, quality-md]
timestamp: 2026-06-24T00:00:00Z
---

# Normalize QUALITY.md self-check roll-up - design

## Context

Answers the [functional spec](spec.md) for the
[0082 change case](../0082-normalize-quality-md-rollup.md). The work changes
guidance, not the format or CLI. The design question is where the real distinction
around `quality-md` belongs: in evaluation semantics, or in follow-up behavior
when the model changes.

## Approach

### 1. Delete the evaluation exception

The implementation should remove live guidance that says the QUALITY.md
self-check is out of the root roll-up. The replacement rule is direct: if
`quality-md` is a modeled area and the evaluation scope includes it, it receives
assessment records, an analysis record, and aggregate contribution like any other
area.

This matches the current mechanical model. `qualitymd status` already reports
`quality-md` as an ordinary area with a declared source, factors, and requirements.
The report builder renders whatever analysis records the skill writes; it has no
`quality-md` exclusion switch. The simplest implementation is therefore to align
the guidance with the tools rather than adding tool support for a special case.

### 2. Keep the area pattern

The setup-time `quality-md` pattern remains useful. It names the concrete artifact
being evaluated (`./QUALITY.md`) and the guide used as assessment criteria. That
pattern does not require special roll-up semantics; it only requires the evaluator
to remember that the area source and the active model file are the same artifact.

### 3. Put the remaining distinction in recommendation follow-up

The only distinct operational behavior is after evaluation: applying a
recommendation against `quality-md` may change the model itself. That is already
covered by the quality-log rule for meaningful model changes. The implementation
should preserve that rule and, if needed, sharpen wording so it is clearly about
mutation history rather than evaluation exclusion.

## Alternatives

- **Keep the learn-loop roll-up exception.** Rejected. It keeps an area visible
  but not consequential, and it asks every evaluator to remember a convention that
  neither the model nor the CLI records can express.
- **Add a schema field such as `rollup: false`.** Rejected. The user is
  specifically questioning special treatment, and this would turn an implicit
  exception into a permanent format feature. It is unnecessary for the immediate
  simplification.
- **Produce two default verdicts, subject quality and model quality.** Rejected.
  That preserves the conceptual split and keeps full evaluation ambiguous. A user
  can still request a scoped evaluation of `quality-md` when they want only model
  quality.
- **Remove the `quality-md` area.** Rejected. The model artifact is a real,
  inspectable project constituent. Removing it would hide model quality rather
  than normalizing it.
- **Keep "learn loop" everywhere but redefine it.** Partly rejected. The term can
  remain where it describes model-change learning and quality-log rationale, but
  live guidance must not use it to imply a separate roll-up axis.

## Trade-offs & risks

- **Aggregate ratings may move downward.** That is intended. A weak
  `QUALITY.md` means the project's quality-management surface is weak, so the
  aggregate should feel that pressure.
- **Self-reference can surprise readers.** The area evaluates the file that
  contains the model. The existing setup comments and assessment wording are the
  mitigation: `source` points at the artifact, while `assessment` names the guide
  used to judge it.
- **Model edits can invalidate prior runs quickly.** This is already true for any
  model change. The run snapshot and quality-log entry preserve the audit trail.
- **Old archived rationale will disagree.** Archived Change Cases should remain
  historical. The new durable guidance and this case supersede the earlier
  rationale without rewriting it.

## Open questions

- Whether implementation should also rename "QUALITY.md self-check" to a less
  exceptional phrase in live guidance, such as "`quality-md` area" or "model
  artifact area", is left to the implementation pass.
