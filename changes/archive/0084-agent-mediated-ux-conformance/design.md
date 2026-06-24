---
type: Design Doc
title: Agent-mediated UX conformance — design doc
description: Design for applying the agent-mediated UX guide across live /quality skill guidance and durable skill specs.
tags: [skill, ux, agents, workflows]
timestamp: 2026-06-24T00:00:00Z
---

# Agent-mediated UX conformance — design doc

Design for
[Agent-mediated UX conformance](../0084-agent-mediated-ux-conformance.md) and
its [functional spec](spec.md).

## Context

The new
[Designing agent-mediated UX](../../../docs/guides/agent-mediated-ux.md) guide
defines how workflows should feel when the agent is the interface. The repo
already has many of the right primitives — run frames, decision briefs,
discovery questions, checkpoints, closeouts — but they are specified unevenly
across the parent `/quality` skill contract and the individual setup, evaluate,
update, and recommendation-follow-up workflows.

The conformance pass should turn the guide into the shared interaction contract
for live workflow guidance without changing quality semantics, CLI behavior, or
historical records.

## Approach

Make the parent `/quality` interaction contract the shared source of truth for
agent-mediated presentation. It should name the guide, require status-first
output, require primary question/call-to-action emphasis, keep supporting labels
scannable, and preserve the existing mutation and confirmation boundaries.

Then update each workflow only where it needs workflow-specific application:

- `setup` owns discovery questions, human-context checkpoints, review gates, and
  setup closeout examples.
- `evaluate` owns long-running progress, stop responses, evaluation summaries,
  limitations, recommendations, and next-action closeout.
- `update` owns version/update plans, mutation confirmation, verification, and
  restart/reload guidance.
- recommendation follow-up owns apply/handoff decision briefs and completion
  reporting.

Keep examples concrete but avoid over-templating. The implementation should add
small exemplar blocks where they clarify the contract, not build a rigid
renderer that every agent must reproduce byte-for-byte.

Mark code and CLI output as verified no-impact unless the audit finds a live
agent-mediated output path in `cmd/` or `internal/`. CLI human output remains
governed by the CLI design guide.

## Alternatives

One option was to put all formatting rules directly in each workflow. That would
make each workflow self-contained, but it would repeat the same output hierarchy
and emphasis rules four times and invite drift.

Another option was to introduce a strict template for every agent response. That
would improve uniformity, but it would make natural agent output brittle and
could crowd out workflow-specific judgment. The guide should govern hierarchy and
scanning, not require identical prose.

## Trade-offs & risks

The main risk is over-formatting: bold labels and semantic emoji can improve
scan quality, but too much emphasis becomes decoration. The implementation must
keep emphasis structural and preserve the existing concise, evidence-led style.

The second risk is confusing CLI UX with agent-mediated UX. The conformance pass
should not redesign deterministic CLI output unless a code audit finds output
that is explicitly a live agent-mediated workflow template.

## Open questions

None. If implementation discovers a reusable rendering abstraction or real code
impact, update this design or add a follow-up design before changing code.
