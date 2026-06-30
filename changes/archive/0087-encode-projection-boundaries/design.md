---
type: Design Doc
title: Encode projection boundaries in the model — design doc
description: Why the projection boundary is encoded comment-primary with an optional description clause, why the rule is general, and how it stays inside the existing description discipline.
tags: [skill, authoring, factors, projections, harnessability]
timestamp: 2026-06-24T00:00:00Z
---

# Encode projection boundaries in the model — design doc

Design behind the [Encode projection boundaries in the
model](../0087-encode-projection-boundaries.md) Change Case and its [functional spec](spec.md).

## Context

The three-projections rule (0076) and the Agent Harnessability guidance (0081,
renamed 0085) are sound: a concern projects as a factor, a constituent/area, and an
audience, and the author is told to _name the projection meant and model it once_.
What is missing is a step _after_ that reasoning — leaving the boundary legible in
the model that ships. This case adds that step at the general rule, with Agent
Harnessability / agent-harness as the worked instance, plus a readiness check.

This is the same move this repo already makes for durable specs: rationale must be
_promoted into_ the artifact (Background / per-requirement annotations), because
reasoning that lives only in the change archives dies there and a later editor
re-litigates it. Here the "artifact" is the generated `QUALITY.md`, and the
"editor" is whoever reads or evaluates it.

## Approach

A general authoring `Do` at the three-projections rule: when a model carries two or
more projections of one concern, encode the boundary at the point of definition.
Two mechanisms, ranked:

1. **YAML comment on each projection's node (primary).** This is _modeling
   rationale_ — "why two nodes for one concern" — which belongs in a comment, not in
   the quality definition. It is exactly what the field model did by hand; this case
   makes it required rather than lucky. Comments are the natural home and keep the
   `description` clean per [Name the quality, not the practice](../../../skills/quality/guides/authoring.md).
2. **Disambiguating clause in each `description` (only when both projections are
   rated nodes that appear in a report).** Comments do not survive rendering: an
   evaluation report shows the Agent Harnessability factor result and the
   agent-harness area result with no comments, so a reader of the _report_ could be
   confused again. The description is the only carrier that survives, so it earns a
   one-clause distinction — but only where it is needed (both sides rated and
   surfaced), to avoid bloating descriptions on the audience projection or on a
   concern whose second projection is not a rated node.

The clause stays inside the existing description discipline: [Write a description
that distinguishes, not
enumerates](../../../skills/quality/guides/authoring.md) already says a description
states what the entity _is_ and how it differs from siblings, and avoids restating
factors/requirements. Distinguishing a same-rooted sibling projection _is_ that
rule, so no new tension is introduced — the clause names the distinction, nothing
more.

## Alternatives

- **Description clause as the primary (or only) mechanism.** Rejected: it mixes
  modeling rationale into the quality definition, fights "name the quality, not the
  practice," and couples a node's description to whether its sibling projection
  happens to exist. The comment is the right home for the _why_; the description
  carries only the minimal distinction, and only when rendering would otherwise
  lose it.
- **Special-case it on the Agent Harnessability guidance only.** Rejected on the
  agreed framing: the confusion is structural to _any_ concern with multiple
  projections (the guide's own `secure` example), so the rule belongs at the
  general three-projections rule with Agent Harnessability as one instance. Bolting
  it only onto the harness guidance would leave the next multi-projection concern
  to re-create the confusion.
- **A deterministic CLI lint rule.** Rejected (and out of scope): whether a comment
  or clause genuinely distinguishes two projections is a judgment, not a
  machine-checkable property. The Top 10 readiness check carries the enforcement at
  the skill layer, where judgment lives; the CLI format stays neutral.
- **Mandate a retrofit of existing models.** Rejected: existing models remain
  structurally valid, and the field model that motivated this already carries the
  boundary. Authoring work updates others opportunistically.

## Trade-offs & risks

- **A check that reads comments.** The Top 10 finding inspects for a boundary note
  (comment or description clause) between two same-rooted projections. It is a
  heuristic — it cannot judge whether the note is _good_ — but it reliably catches
  the failure mode that motivated the case (no note at all). False positives are
  unlikely because it only fires when two same-rooted projections are both present.
- **Description-clause restraint.** The "only when both are rated nodes in a
  report" condition keeps descriptions from accreting boilerplate. The risk is an
  author over-applying it to the audience projection; the spec scopes it explicitly
  to rated, report-surfaced nodes to prevent that.

## Open questions

None. The mechanisms, their ranking, and the report-surfaced condition are settled;
the CLI-lint and retrofit options are deliberately out of scope.
