---
type: Functional Specification
title: Publish the example evaluation to the docs site — functional spec
description: Requirements for the generator that renders the report-gallery example quality evaluation into Mintlify pages.
tags: [docs, mintlify, report-gallery]
timestamp: 2026-06-29T00:00:00Z
---

# Publish the example evaluation to the docs site — functional spec

Companion to the
[Publish the example evaluation to the docs site](../0188-publish-example-evaluation-to-docs.md)
change case. This spec states *what* the converter must do.

The key words **MUST**, **SHOULD**, and **MAY** are to be interpreted as
described in IETF RFC 2119.

## Scope

A generator renders the Markdown pages of the report-gallery example evaluation
(`examples/report-gallery/software-service/.quality/evaluations/0001-full-eval/`)
into Mintlify pages under `mintlify/examples/software-service/`. The gallery
remains the single source of truth; the generator only transforms it for
publication. The nav surfaces the report alone; all other pages are published
but reached by link. Site-search behavior of unlisted pages and any second
example are out of scope.

## Requirements

### Page generation

- The generator **MUST** render every `.md` file under the source evaluation
  directory to a corresponding `.mdx` file at the mirrored path under
  `mintlify/examples/software-service/`, so the published route tree matches the
  gallery tree.
- The generator **MUST** treat its output directory as fully generated: a run
  **MUST** reproduce the output from the source alone, leaving no page from a
  source file that no longer exists.
- Each generated page **MUST** carry Mintlify frontmatter with a `title` taken
  from the source page (its frontmatter `title`, falling back to its leading H1),
  and **MUST** drop the source leading H1 so the title is not rendered twice.
- The generated output **MUST** be deterministic — identical source yields
  byte-identical output — so the `report-docs-check` gate is stable.

### Link rewriting

- A link to another generated page (a `.md` target within the source tree)
  **MUST** be rewritten to that page's internal docs route, preserving any
  anchor.
- A link whose target is not a generated page — repository data files
  (`data/**.json`) and the shared `glossary.md` — **MUST** be rewritten to a
  stable external target (the file on GitHub for repo files; the repository-root
  glossary for `glossary.md`), since the gallery's relative `glossary.md` target
  does not resolve.
- In-page anchors and already-absolute URLs **MUST** be left unchanged, and link
  rewriting **MUST NOT** alter text inside code spans.

### Build safety

- The generator **MUST** fail loudly if a page would contain characters MDX
  parses as JSX outside code spans, except the intentional
  `<a id="…"></a>` finding anchors, which it **MUST** preserve so in-page finding
  links resolve.

### Navigation

- The generator **MUST** add exactly one navigation entry — an `Examples` group
  containing the single report page — and that entry **MUST** render in the
  sidebar as "Example quality evaluation". It **MUST NOT** list the other
  published pages in the navigation.
- Updating the navigation **MUST** preserve the rest of `mintlify/docs.json`
  unchanged on re-runs once the entry is present (idempotent).

### Integration

- A `report-docs` task **MUST** run the generator, and a `report-docs-check`
  task **MUST** fail when the committed output is stale; `report-docs-check`
  **MUST** be part of the `check` gate.
- The pre-commit hook **MUST** regenerate the example pages *after* it
  regenerates the report gallery, since the converter reads the gallery output.

## Durable spec changes

### To add

None. The converter consumes the already-generated gallery and emits
presentation-only docs pages; it introduces no new format or behavioral contract
that a durable spec must own. The report structure it renders is already
governed by the report specs that drive the gallery generator.

### To modify

None.

### To delete

None.
