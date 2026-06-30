---
type: Change Case
title: Publish the example evaluation to the docs site
description: Generate Mintlify pages for the report-gallery example quality evaluation, surfaced under a single "Example quality evaluation" nav entry.
status: In-Review
tags: [docs, mintlify, report-gallery]
timestamp: 2026-06-29T00:00:00Z
---

# Publish the example evaluation to the docs site

The [report gallery](../examples/report-gallery/software-service/README.md)
already generates a complete, illustrative quality evaluation for the fictional
LedgerLite service, but it is only browsable as Markdown files in the repo. This
case publishes that evaluation to the Mintlify docs site so the report design is
visible to readers of the docs.

- [Functional spec](0188-publish-example-evaluation-to-docs/spec.md) — what the
  case must do.

No design doc: the approach reuses the established generated-page pattern
(`scripts/sync-spec-docs.mjs`, `scripts/cli-docs`) closely enough that the spec
carries the whole contract.

## Motivation

The example evaluation is the clearest demonstration of what running QUALITY.md
produces — the roll-up report, ranked findings, recommendations, and the
drill-down area/factor/requirement pages. Today a reader has to find it in the
repo tree. Publishing it to the docs site puts a real, navigable evaluation one
click from the docs, while keeping the gallery as the single generated source so
the two cannot drift.

## Scope

**Covered:** a generated converter that renders every Markdown page of the
gallery evaluation into Mintlify pages, rewrites cross-links to internal docs
routes (and out-of-site targets to GitHub), and adds a single
"Example quality evaluation" entry to the docs navigation. The converter is
wired into the same CI gate and pre-commit hook as the other generated docs
pages, running after the gallery itself regenerates.

**Deferred:** surfacing more than the report in the sidebar (findings,
recommendations, and the drill-down pages are published but reached by link, not
listed in the nav); porting any second gallery example; and site search indexing
of the hidden pages.

## Affected artifacts

**Code**

- `scripts/report-docs.mjs` *(new)* — the converter: gallery Markdown →
  `mintlify/examples/software-service/**` `.mdx`, with link rewriting, MDX-safety
  checks, and the docs.json nav splice.

**Durable docs**

- `mintlify/examples/software-service/**` *(new, generated)* — the published
  evaluation pages.
- `mintlify/docs.json` — adds the `Examples` nav group with the single
  `examples/software-service/report` page.
- `CONTRIBUTING.md` — documents `mise run report-docs` / `report-docs-check`
  alongside the other generated-docs tasks.

**Build / scaffold**

- `mise.toml` — adds `report-docs` and `report-docs-check` tasks; adds
  `report-docs-check` to the `check` gate.
- `.githooks/run-check` — regenerates the example pages after the report gallery
  in the staged-commit path.

**Durable specs:** none. The converter owns no format contract of its own; its
input is the already-generated gallery and its output is presentation. See the
spec's [Durable spec changes](0188-publish-example-evaluation-to-docs/spec.md#durable-spec-changes).

## Status

`In-Review`. Converter, wiring, and generated pages are complete; the full
`mise run check` gate and Mintlify's `mint broken-links` both pass. Durable docs
(`CONTRIBUTING.md`) are in sync. See the
[status lifecycle](index.md#status-lifecycle).
