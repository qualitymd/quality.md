---
type: Change Case
title: Publish the example evaluation to the docs site
description: Surface the report-gallery example quality evaluation from the docs site through a single external nav link to the report on GitHub.
status: In-Review
tags: [docs, mintlify, report-gallery]
timestamp: 2026-06-29T00:00:00Z
---

# Publish the example evaluation to the docs site

The [report gallery](../examples/report-gallery/software-service/README.md)
already generates a complete, illustrative quality evaluation for the fictional
LedgerLite service, browsable as Markdown files in the repo. This case surfaces
that evaluation from the Mintlify docs site so a reader can reach a real,
rendered report in one click.

- [Functional spec](0188-publish-example-evaluation-to-docs/spec.md) — what the
  case must do.

No design doc: the case is a single navigation link to an existing file; the
spec carries the whole contract.

## Motivation

The example evaluation is the clearest demonstration of what running QUALITY.md
produces — the roll-up report, ranked findings, recommendations, and the
drill-down area/factor/requirement pages. Today a reader has to find it in the
repo tree. Linking to it from the docs sidebar puts a real evaluation one click
from the docs, while pointing directly at the gallery's source so the docs
cannot drift from it.

## Scope

**Covered:** a single sidebar navigation link in `mintlify/docs.json` that opens
the gallery's report (`report.md`) on GitHub in a new tab.

**Deferred:** generating, embedding, or mirroring any gallery pages inside the
docs site; porting any second gallery example; site search indexing of the
evaluation.

## History

An earlier revision of this case generated the gallery evaluation into Mintlify
`.mdx` pages under `mintlify/examples/software-service/` via a
`scripts/report-docs.mjs` converter, wired into a `report-docs` /
`report-docs-check` gate and the pre-commit hook. That machinery was removed in
favor of linking out: the generated pages duplicated the gallery, added a
converter and a CI gate to keep in sync, and embedded a long evaluation tree the
nav never surfaced. A direct link to the source report is simpler and cannot
drift.

## Affected artifacts

**Durable docs**

- `mintlify/docs.json` — adds an `Example report` sidebar link (under a
  `Documentation` anchor) pointing to the gallery `report.md` on GitHub.
- `CONTRIBUTING.md` — drops the `report-docs` / `report-docs-check` task entries
  and the generated-Mintlify-pages mention in the git-hooks section.

**Build / scaffold**

- `mise.toml` — removes the `report-docs` and `report-docs-check` tasks and drops
  `report-docs-check` from the `check` gate.
- `.githooks/run-check` — removes the example-page regeneration step from the
  staged-commit path.

**Removed**

- `scripts/report-docs.mjs` and the generated `mintlify/examples/**` tree.

**Durable specs:** none. The link owns no format contract; the report structure
it points at is already governed by the report specs that drive the gallery
generator.

## Status

`In-Review`. The nav link is in place, the generator and generated pages are
removed, and the `mise run check` gate no longer references `report-docs`. See
the [status lifecycle](index.md#status-lifecycle).
