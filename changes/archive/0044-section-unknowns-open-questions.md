---
type: Change Case
title: Section unknowns and open questions
description: Replace the standalone Known gaps body section with per-section unknowns, open questions, and a human/agent review state line across the format spec, skill, scaffold, checks, specs, and dogfood instances.
status: Done
tags: [authoring, body, skill, scaffold, spec]
timestamp: 2026-06-20T00:00:00Z
---

# Section unknowns and open questions

A Change Case that makes the Markdown body more rigorous and consistent: every
body section follows one shape, carries its own unknowns and open questions, and
closes with a review-provenance state line. This retires the separate **Known
gaps** body section, whose content now lives in the section it concerns.

> **Done.** Implementation, durable spec/doc updates, fixture and dogfood
> migration, and verification (`go test ./...`, `mise run check`) are complete,
> and the change is archived for release.

## Motivation

The authoring guidance treated the body as a loose set of suggested sections,
with a single catch-all **Known gaps** section sitting far from the content it
qualified. On a high-leverage, largely agent-authored file that is too weak:
gaps get skimmed, and a reader cannot tell which context a human has actually
vetted versus what an agent drafted.

This change gives each section a common shape (purpose, contents, its own
unknowns and open questions) and a state line that records the last *human*
review (cite a person) distinctly from the last *agent* review — so unreviewed
agent edits are visible. Unknowns and open questions are scoped to the section's
own topic and are kept distinct: an unknown is a broad area of uncertainty that
may not resolve to a single answer; an open question is a specific question with
one particular answer, still unresolved. Both are context that feeds the model,
not commentary on the model.

## Scope

Covered: the authoring guide and its functional spec; the format spec's body
guidance; the `qualitymd init` scaffold and its test; the init CLI spec; the
skill's setup mode, getting-started guide, and top-10 checks (runtime and spec);
the bundled example fixtures; and the repo's own dogfood `QUALITY.md` and active
evaluation model.

Deferred / non-goals: no new frontmatter schema, lint rule, or CLI behavior — the
state line and per-section unknowns/open questions remain freeform body prose, not
structured fields the tooling parses or validates. Historical records
(`archive/**`, `log.md` files) stay frozen.

## Affected artifacts

### Code

- [x] `internal/scaffold/skeleton.md` — drop the `## Known gaps` section; give
      each body section the new shape (unknowns, open questions, state line).
- [x] `internal/scaffold/scaffold_test.go` — stop asserting `## Known gaps`;
      assert the new per-section markers instead.

### Durable specs

See the functional spec's
[Durable spec changes](0044-section-unknowns-open-questions/spec.md#durable-spec-changes)
for the per-requirement breakdown.

- [x] [`SPECIFICATION.md`](../../SPECIFICATION.md) — body context names unknowns and
      open questions instead of known gaps.
- [x] [`specs/cli/init.md`](../../specs/cli/init.md) — recommended body sections.
- [x] [`specs/skills/quality-skill/guides/authoring-md.md`](../../specs/skills/quality-skill/guides/authoring-md.md)
      — recommended-section MUST and the body best-practice coverage.
- [x] [`specs/skills/quality-skill/guides/top-10-quality-md-checks-md.md`](../../specs/skills/quality-skill/guides/top-10-quality-md-checks-md.md)
      — body-context check.
- [x] [`specs/skills/quality-skill/quality-skill.md`](../../specs/skills/quality-skill/quality-skill.md)
      — setup first-population body list.
- [x] [`specs/log.md`](../../specs/log.md) — record the spec changes.

### Durable docs (bundled skill)

- [x] `skills/quality/guides/authoring.md` — body shape, state line, unknowns and
      open questions (already revised).
- [x] `skills/quality/guides/getting-started.md` — first-pass body outcome and
      checks.
- [x] `skills/quality/guides/top-10-quality-md-checks.md` — body-context check.
- [x] `skills/quality/modes/setup.md` — guided first-population body list.

### Fixtures and dogfood instances

- [x] `specs/skills/quality-skill/examples/0001-subject-quality-eval/model.md`
      and `.../recommendations/002-produce-reconciliation-evidence.md`.
- [x] `QUALITY.md` — the repo's own model body.
- [x] `quality/evaluations/0005-subject-quality-eval/model.md` — active eval model.

## Children

- [Functional spec](0044-section-unknowns-open-questions/spec.md) — what the
  change must do.
