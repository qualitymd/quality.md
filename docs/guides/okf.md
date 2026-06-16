# Working with OKF

Some of our docs — notably the [`specs/`](../../specs/index.md) bundle — are
authored as **OKF** (Open Knowledge Format): a directory of Markdown files with
YAML frontmatter, where each file is a self-describing *concept*. This guide
covers the essentials for reading and editing those files. For the full format,
see the [OKF specification](https://github.com/GoogleCloudPlatform/knowledge-catalog/blob/main/okf/SPEC.md).

## The essentials

- A **bundle** is a directory tree of Markdown files.
- A **concept** is one Markdown file: YAML frontmatter + a Markdown body.
- Two filenames are reserved at any level: `index.md` (a listing) and `log.md`
  (a change history). Every other `.md` file is a concept.

## Concept frontmatter

Each concept begins with a YAML frontmatter block. Only `type` is required:

```yaml
---
type: Specification          # REQUIRED — what kind of concept this is
title: qualitymd CLI         # recommended — display name
description: One-line summary # recommended — used in index listings
tags: [cli, specification]   # optional
timestamp: 2026-06-16T00:00:00Z  # optional — ISO 8601 last-modified
---
```

`type` is a free-form string — pick something descriptive (`Specification`,
`Command Specification`, `Reference`). There is no central registry. You MAY add
any other keys you need.

## Reserved files

- **`index.md`** — a progressive-disclosure listing of a directory's contents, so
  a reader (or agent) can see what's available before opening files. It has **no
  frontmatter** (the one exception: a bundle-root `index.md` may carry
  `okf_version`). Group entries under headings:

  ```markdown
  # Commands

  - [qualitymd init](init.md) - scaffold a starter QUALITY.md.
  - [qualitymd lint](lint.md) - validate a file's structure.
  ```

- **`log.md`** — an optional change history, newest first, with ISO-8601 date
  headings:

  ```markdown
  # Update Log

  ## 2026-06-16

  - **Initialization**: Created the bundle.
  ```

## Links

Link between concepts with ordinary Markdown links. We use **relative** links
(`./lint.md`, `../cli.md`) so they resolve in the GitHub UI. Broken links are
tolerated — they may just point at not-yet-written concepts.

## When you add or edit a concept

1. Give it parseable frontmatter with a non-empty `type`.
2. Update the enclosing `index.md` so the listing stays accurate.
3. Add a `log.md` entry for the change.

That's the whole contract: parseable frontmatter, a `type`, and tidy `index.md` /
`log.md`. Everything else is soft guidance — consumers won't reject a bundle over
missing optional fields, unknown types, or broken links.
