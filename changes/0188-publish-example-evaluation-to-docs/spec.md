---
type: Functional Specification
title: Publish the example evaluation to the docs site — functional spec
description: Requirements for the docs-site navigation link that surfaces the report-gallery example quality evaluation.
tags: [docs, mintlify, report-gallery]
timestamp: 2026-06-29T00:00:00Z
---

# Publish the example evaluation to the docs site — functional spec

Companion to the
[Publish the example evaluation to the docs site](../0188-publish-example-evaluation-to-docs.md)
change case. This spec states _what_ the navigation link must do.

The key words **MUST**, **SHOULD**, and **MAY** are to be interpreted as
described in IETF RFC 2119.

## Scope

The docs site surfaces the report-gallery example evaluation through a single
navigation link to the report on GitHub. The gallery remains the single source
of truth; the docs site does not copy, generate, or embed any of its pages.
Site search indexing of the gallery and any second example are out of scope.

## Requirements

### Navigation link

- `mintlify/docs.json` **MUST** present exactly one navigation entry for the
  example evaluation, labeled `Example report`, in the left sidebar.
- The entry **MUST** link to the gallery report
  (`examples/report-gallery/software-service/.quality/evaluations/0001-full-eval/report.md`)
  on GitHub, and **MUST** open in a new tab (an external `href`).
- The entry **SHOULD** carry an icon indicating it leaves the site.

### No in-site generation

- The docs site **MUST NOT** contain generated or hand-authored copies of the
  gallery evaluation pages.
- No build task, CI gate, or git hook **MUST** depend on rendering the gallery
  into docs pages.

## Durable spec changes

### To add

None. The link consumes the already-generated gallery and introduces no new
format or behavioral contract that a durable spec must own.

### To modify

None.

### To delete

None.
