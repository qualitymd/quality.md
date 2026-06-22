---
type: Design Doc
title: Prospective evaluation plan artifacts - design doc
description: How the evaluation workflow and artifact specs make design.md and the initial plan.md prospective planning artifacts without changing CLI scaffolding.
tags: [evaluation, skill, records]
timestamp: 2026-06-22T00:00:00Z
---

# Prospective evaluation plan artifacts - design doc

Design behind the
[Prospective evaluation plan artifacts](../0056-prospective-evaluation-plan-artifacts.md)
change case and its [functional spec](spec.md). The spec says *what* must hold:
`design.md` and the initial `plan.md` are authored before assessment starts.
This doc says how to land that behavior across the durable specs and runtime
skill without changing the CLI's mechanical run scaffolding.

## Context

The current evaluation workflow already has the right physical artifact order:
`qualitymd evaluation create` creates the run folder and seeds `design.md` and
`plan.md` before records exist. The drift risk is procedural. The runtime skill
can still fill those files too late, and phrases like "executed or inspected
evidence basis" blur planning content with after-the-fact evidence.

This change is therefore a contract and prompt repair, not a new CLI feature.
The CLI keeps creating stubs. The skill gets a clearer planning checkpoint
between run creation and assessment. Durable artifact specs make the distinction
between prospective plan, process log, and formal judgment hard to lose again.

## Approach

### Add a planning checkpoint to the workflow

Update the durable skill evaluation workflow and runtime evaluate mode around a
single checkpoint:

```text
create run -> author design.md and initial plan.md -> add settled coverage if useful -> assess
```

The key implementation detail is ordering language. The prompt should say the
skill completes the initial `design.md` and `plan.md` before evidence collection
or record writes begin. Assessment can still discover surprises, but those
surprises update the plan explicitly rather than rewriting the initial plan into
a retrospective narrative.

### Split design, plan, debug log, and records by job

Use the artifact specs and skill prompt to give each file one job:

- `design.md` records the resolved evaluation frame: model snapshot, scope,
  in-scope areas, exclusions, and known method limits.
- `plan.md` records intended execution: rigor, requirement coverage, intended
  evidence basis, planned checks, and planned limitations.
- `debug-log.md` records process events and recovery notes.
- assessment, analysis, recommendation, and report artifacts carry actual
  findings, evidence, rating rationale, advice, and final presentation.

That split removes the ambiguous "executed or inspected evidence basis" wording
from planning instructions. The plan can name intended commands or source reads;
the records cite what actually happened.

### Add a `design.md` artifact spec

`plan.md` already has a dedicated artifact spec. `design.md` does not, even
though the workflow depends on it. Add `specs/evaluation-records/design-md.md`
as a small 1:1 artifact spec and link it from the evaluation-records parent,
child index, and run-folder spec.

The design artifact spec should stay narrow: it owns purpose, timing, and
content categories. It should not prescribe a rigid template or make reports
parse arbitrary prose.

### Treat later changes as amendments

The runtime guidance should allow the skill to adapt when evaluation reveals a
scope or coverage mismatch. The design choice is to use a plainly named section,
for example `## Plan updates`, rather than a separate machine artifact.

If `coverage:` frontmatter exists, the skill updates it alongside the amendment
when the planned assessment or analysis set changes. That keeps resume
diagnostics aligned with the human-readable explanation without inventing a
second coverage lifecycle.

### Keep CLI behavior unchanged

No changes are needed to `qualitymd evaluation create`: it already writes the
files at the right time as stubs. The CLI cannot author scope judgment, rigor
trade-offs, or evidence strategy, so moving more content into the create command
would put judgment in the wrong layer.

## Alternatives

- **Have the CLI generate fuller templates.** Rejected for this case. Templates
  might help later, but they do not solve the timing problem by themselves, and a
  rigid template risks turning evaluator judgment into form filling.
- **Add a separate `plan-updates.md` artifact.** Rejected as unnecessary. The
  problem is provenance, not storage capacity; a clear section in `plan.md`
  preserves the initial plan and keeps the update beside the thing it changes.
- **Put actual coverage summaries in `debug-log.md`.** Rejected because the
  debug log is process-only. Using it for coverage or rating context would blur
  the same boundary this case is trying to sharpen.
- **Rely on report outputs to reconstruct the plan.** Rejected. Reports describe
  final judgment; they cannot prove what the evaluator intended before evidence
  collection began.

## Trade-offs & risks

- **More explicit prompt language can feel procedural.** Acceptable here because
  the ordering is the whole quality bar; loose language is what allowed the
  after-the-fact pattern.
- **Amendments are still hand-authored.** The skill can forget to mark a change
  unless the runtime prompt is direct. The durable spec makes the expectation
  reviewable, but it is not mechanically enforced.
- **A new `design.md` spec adds one more artifact contract.** Worth it because
  `design.md` is already a required run artifact and the absence of a contract is
  part of why its role is easier to blur.

## Open questions

- Should the runtime skill include a compact heading skeleton for `design.md` and
  `plan.md`, or should the specs define content categories only and leave the
  prose shape flexible?
- Should `qualitymd evaluation status` ever warn when `plan.md` has `coverage:`
  but no visible amendment after coverage changes, or is that too heuristic for
  a prose artifact?
