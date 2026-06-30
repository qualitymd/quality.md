---
type: Change Case
title: Evaluation orchestration overhaul
description: Replace evaluation rigor levels with one best-quality evaluate workflow — exhaustive coverage, parallel-by-default collection, and an always-on two-pronged QC phase.
status: Done
tags: [skill, evaluation, orchestration, qc]
timestamp: 2026-06-26T00:00:00Z
---

# Evaluation orchestration overhaul

Make `evaluate` a single best-quality workflow. Remove the variable evaluation
rigor dial (`quick`/`standard`/`deep`), default to exhaustive coverage, fan out
collection and verification to subagents when the harness supports them, and
promote verification into an always-on, two-pronged **QC phase** that guards both
false positives (unreal findings) and false negatives (missed findings).

- [Functional spec](0129-evaluation-orchestration-overhaul/spec.md) — what the
  change must do.
- [Design doc](0129-evaluation-orchestration-overhaul/design.md) — how the skill
  and durable specs make it so.

## Motivation

Evaluation rigor (`quick`/`standard`/`deep`) is a coverage-and-cost dial. Its
only reason to ever trade away coverage was wall-clock and token cost. Parallel
subagent fan-out dissolves that tradeoff: exhaustive per-area/per-requirement
collection and verification run concurrently, so thoroughness stops being slow.
Once that holds, the dial's job evaporates and it becomes a footgun — a `quick`
pass silently skips requirements, which is exactly the failure mode the spec
already warns about ("a shallow pass never reads as whole coverage").

The real knob is **scope** (which areas/factors/requirements), already handled
cleanly by model references and `--narrowing`. Users should go faster by
narrowing scope, never by lowering quality.

Verification today is a conditional step tangled into the procedure: `standard`
re-checks only the headline-binding findings; adversarial fan-out is gated behind
`deep`. That conflates two distinct failure modes and only addresses one.
Re-checking guards **false positives** ("is the finding I surfaced real?") but
not **false negatives** ("what did I fail to surface?"). The highest-risk place
for a missed finding is any area or requirement the first pass returned only
strengths or nothing for — usually "didn't look hard enough," not "nothing
there." A single best workflow needs a first-class QC phase that guards both.

## Scope

Covered:

- Remove evaluation rigor: the `quick`/`standard`/`deep` levels, the `Rigor:`
  run-frame field, the `--rigor`/`deep` invocation forms, and the feedback-log
  `rigor:` field.
- Make exhaustive coverage and the two-pronged QC phase the mandatory evaluation
  contract; make parallel subagent fan-out the default execution strategy when
  the harness exposes a subagent capability, with an identical serial fallback
  when it does not.
- Promote verification into an always-on QC phase with two prongs — **verify**
  (re-run the command/search for every roll-up-binding and low-confidence
  finding) and **completeness sweep** (coverage ledger, quiet-zone re-examination,
  thin-evidence escalation), with a bounded re-collection loop and an explicit
  stop condition.
- Preserve the orchestrator-owned invariant: roll-up judgment and ratings stay
  with the orchestrating skill; subagents return structured findings only.

Deferred / non-goals:

- **Modeling rigor** (a setup-discovery concept) and per-requirement
  **assessment rigor** (a model-authoring guide concept) are unrelated and
  unchanged.
- No CLI changes. Rigor is a skill-side concept; `qualitymd evaluation` records,
  schemas, and report shapes do not move.
- No recommendation/advise behavior change beyond what already exists.

## Affected artifacts

Derived by repo sweep for the evaluation-rigor tokens (`rigor`, `--rigor`,
`evaluate deep`, the `Rigor:` frame field), excluding the unrelated _modeling
rigor_ / _assessment rigor_ senses and frozen history (`changes/archive/`,
append-only `log.md` files, prior `.quality/` run records). Reconciled before
archival.

**Code:**

- [x] None — rigor is a skill-prompt concept; no `cmd/` or `internal/` change.

**Bundled skill (durable docs):**

- [x] [`skills/quality/SKILL.md`](../../skills/quality/SKILL.md) — remove the Rigor
      Levels table, the `Rigor:` run-frame field, the `/quality evaluate deep`
      invocation; replace the deep-only subagent guidance with the
      parallel-by-default + QC contract.
- [x] [`skills/quality/workflows/evaluate.md`](../../skills/quality/workflows/evaluate.md)
      — remove the `## Rigor` section and the `**Rigor:**` frame field; drop the
      feedback-log `rigor:` field; thread the QC phase and fan-out into the
      procedure; drop `rigor` from the `EvaluationFrame` contents.
- [x] [`docs/guides/agent-mediated-ux.md`](../../docs/guides/agent-mediated-ux.md)
      — remove the `Rigor: standard` line from the example evaluate run frame
      (line ~175). The unrelated "modeling rigor" mention stays.

**Durable specs:** see the functional spec's
[Durable spec changes](0129-evaluation-orchestration-overhaul/spec.md#durable-spec-changes).
Indexed here:

- [x] [`specs/skills/quality-skill/evaluation.md`](../../specs/skills/quality-skill/evaluation.md)
- [x] [`specs/skills/quality-skill/quality-skill.md`](../../specs/skills/quality-skill/quality-skill.md)
- [x] [`specs/skills/quality-skill/index.md`](../../specs/skills/quality-skill/index.md)
- [x] [`specs/skills/quality-skill/workflows/evaluate/feedback-log.md`](../../specs/skills/quality-skill/workflows/evaluate/feedback-log.md)

**Changelog:**

- [x] [`CHANGELOG.md`](../../CHANGELOG.md) — add the release entry. Prior rigor
      mentions are historical and stay.

## Status

`Done`. Functional spec and design doc are settled; implementation is complete
across the bundled evaluate workflow, durable skill specs, agent-mediated UX
guide, logs, and changelog. Verified with `mise run check`; archived with the
parent and child folder under `changes/archive/`.
