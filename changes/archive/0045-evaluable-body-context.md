---
type: Change Case
title: Evaluable body context
description: Clarify the Markdown body as concise, self-explanatory, agent-accessible judgment context for building, justifying, and evaluating model quality.
status: Done
tags: [authoring, body, skill, accessibility]
timestamp: 2026-06-21T00:00:00Z
---

# Evaluable body context

A Change Case for tightening the authoring guide's Markdown-body guidance so
the body is treated as evaluable judgment context, not merely explanatory prose.
The body should be concise and progressively disclosed, while still rigorous
enough for a later human or agent to evaluate the body's own quality and the
model quality it supports.

> **Done.** Implementation, durable spec/doc updates, scaffold coverage, and
> verification (`go test ./...`, `mise run check`) are complete, and the change
> is archived for release.

## Motivation

The current authoring guide already says the Markdown body provides judgment
context, records per-section unknowns and open questions, and should be written
before the model tree. It still understates two important jobs.

First, the body is not only context for interpreting the model. It is context
for building the model, understanding the model's purpose, evaluating whether
the model is still fit for that purpose, and judging the quality of the body
itself: completeness, thoroughness, recency, specificity, grounding, and
accessibility.

Second, accessibility needs to be explicitly agent-accessible. If a useful
supporting source exists but an evaluating agent cannot inspect it through the
repository, configured tools, cited public links, or explicitly provided
context, that limitation is first-class model context. It should be captured
where it bears on the body, usually as a section unknown or open question, not
hidden as an evaluator surprise later.

## Scope

Covered: authoring guidance for the Markdown body as evaluable judgment context;
body-section expectations for concise, self-explanatory, progressively disclosed
content; grounding in cited support; and explicit treatment of inaccessible
support as an agent-accessibility limitation.

Deferred / non-goals: no new frontmatter schema, no required body section names,
no required `Access gaps` line in every section, and no lint or CLI behavior that
parses or validates agent-accessibility prose. This case does not redesign the
per-section unknowns/open-questions convention from 0044; it clarifies how
support accessibility fits into that convention.

## Affected artifacts

### Code

- [x] `internal/scaffold/skeleton.md` — review whether the starter body should
      mention agent-accessible support gaps in its unknowns/open-questions
      prompts.
- [x] `internal/scaffold/scaffold_test.go` — update only if scaffold markers or
      expected starter text change.

### Durable specs

See the functional spec's
[Durable spec changes](0045-evaluable-body-context/spec.md#durable-spec-changes)
for the per-requirement breakdown.

- [x] [`SPECIFICATION.md`](../../SPECIFICATION.md) — body semantics may need to say
      the body supports building, using, and evaluating the Model, not only
      interpreting it.
- [x] [`specs/skills/quality-skill/guides/authoring.md`](../../specs/skills/quality-skill/guides/authoring.md)
      — contract for the runtime authoring guide's body guidance.
- [x] [`specs/skills/quality-skill/guides/getting-started.md`](../../specs/skills/quality-skill/guides/getting-started.md)
      — first-pass body outcome and checks.
- [x] [`specs/skills/quality-skill/guides/top-10-quality-md-checks.md`](../../specs/skills/quality-skill/guides/top-10-quality-md-checks.md)
      — body-context and model-quality inspection checks.
- [x] [`specs/skills/quality-skill/quality-skill.md`](../../specs/skills/quality-skill/quality-skill.md)
      — setup/wizard guidance if the mode contract needs explicit
      agent-accessibility routing.
- [x] [`specs/cli/init.md`](../../specs/cli/init.md) — starter-body contract if the
      scaffold prompt changes.
- [x] [`specs/log.md`](../../specs/log.md) — record durable spec updates when they
      are made.

### Durable docs and bundled skill

- [x] `skills/quality/guides/authoring.md` — primary runtime guide update. The
      tracked `.agents/skills/quality` path is a symlink to this skill.
- [x] `skills/quality/guides/getting-started.md` — first-pass body checks if
      they need to mention evaluability or agent-accessible support.
- [x] `skills/quality/guides/top-10-quality-md-checks.md` — model-quality
      inspection should flag material support that is not agent-accessible.
- [x] `skills/quality/modes/setup.md` — setup guided population if it should ask
      for or record important inaccessible support.
- [x] [`README.md`](../../README.md) — public format summary if the Markdown-body
      description needs the broader judgment-context framing.

### Fixtures and dogfood instances

- [x] `internal/scaffold/skeleton.md` starter content, if updated above.
- [x] `QUALITY.md` and active evaluation model reviewed; update only if the new
      guidance reveals a concrete body-context access gap, not merely because the
      guide wording changes.

## Children

- [Functional spec](0045-evaluable-body-context/spec.md) — what the change must
  do.
