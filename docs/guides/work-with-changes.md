---
type: How-to Guide
title: Working with changes
description: How to create a change and move it through its lifecycle in the changes/ bundle.
tags: [changes, contributing]
timestamp: 2026-06-17T00:00:00Z
---

# Working with changes

Incremental work on the repo is tracked in the
[`changes/`](../../changes/index.md) OKF bundle. A **Change** is one unit of work:
a parent concept that records *why*, *what state* it's in, and *which durable
specs and docs it touches*, over a [functional spec](write-functional-specs.md)
(the *what*) and an optional [design doc](write-design-docs.md) (the *how*). This
guide covers creating one and moving it through its full lifecycle to review and
landing.

## Create the change

1. **Pick the next number.** Changes are numbered sequentially —
   `NNNN-<slug>`, zero-padded to four digits. Look at the highest existing
   `changes/NNNN-*` (and [`archive/`](../../changes/archive/index.md)) and add one.
2. **Add the parent concept** at `changes/NNNN-<slug>.md` with `type: Change`.
   Give it `status: Draft`, a one-line `description`, and a body that states the
   **motivation**, **scope**, and **affected specs & docs** (see
   [below](#account-for-the-specs-and-docs-it-touches)), then links to its
   children.
3. **Add the child folder** `changes/NNNN-<slug>/` with an `index.md` listing and
   the detail concepts (below). The parent-concept-plus-folder shape mirrors the
   rest of the repo (e.g. `specs/cli.md` + `specs/cli/`).

Copy [`0001-example-change`](../../changes/archive/0001-example-change.md) and its
folder as a starting template.

## Write the spec, then the design

Write them in order — settle *what* the change must do before working out *how*.
The two stages line up with the `status` lifecycle below.

1. **Functional spec first** (`spec.md`, `type: Functional Specification`), while
   the change is **Draft**. Follow
   [Writing functional specs](write-functional-specs.md): state requirements with
   RFC 2119 keywords, specify behavior not implementation, and don't invent
   requirements. Pin down the *what* before moving on.
2. **Design doc next** (`design.md`, `type: Design Doc`), once the spec is settled
   and the change moves to **Design**. Follow
   [Writing design docs](write-design-docs.md): it answers the spec, so don't
   start it until the spec is stable. Omit it for a change that needs no design
   discussion.

## Account for the specs and docs it touches

A change usually changes more than code. Its
[functional spec](write-functional-specs.md) states the *delta* — what this unit
of work must do — and is archived with the change once it lands. The enduring
[`specs/`](../../specs/index.md) bundle, the repository-root
[`SPECIFICATION.md`](../../SPECIFICATION.md), the [`README.md`](../../README.md),
and the [`docs/`](../index.md) guides hold the *cumulative* source of truth for
the tool's current behavior. A change is the bridge between them: it proposes a
delta, and the enduring artifacts absorb it.

List the durable specs and docs the change creates or updates in the parent
concept's **Affected specs & docs** section, decided up front alongside the
motivation and scope — not discovered at the end. Keep it honest: an empty list
must read as a deliberate "no durable changes," and every listed artifact is
created or updated *as part of the change*, before it reaches **In-Review**. This
is what keeps the enduring specs from drifting while the change waits for review.

## Move it through the lifecycle

A Change's `status` frontmatter advances, in order:

**Draft → Design → In-Progress → In-Review → Done**

- **Draft** while you write up the [functional spec](write-functional-specs.md)
  (the *what*); **Design** while you work out the
  [design doc](write-design-docs.md) (the *how*); **In-Progress** while it's
  being implemented; **In-Review** after implementation is complete and ready for
  review; **Done** only after the work lands.
- **Modify nothing outside what the current phase authorizes.** This is a
  whitelist, not a blacklist: the question is never "is this file on the Affected
  list?" but "does this phase permit changing it?" Before **In-Progress**, all
  work stays inside the change's own folder, `changes/NNNN-<slug>/` — **Draft**
  produces its `spec.md`, **Design** produces its `design.md`, and either MAY add
  supplementary files (sketches, examples, data fixtures) that support those, so
  long as they live within that folder. **In-Progress** is the first phase that
  authorizes touching anything *outside* the change folder — the code and the
  durable specs and docs. The change's parent concept, its `index.md`, and the
  bundle [`log.md`](../../changes/log.md)/`index.md` are updated in every phase to
  record the step. Everything else — all code, every durable spec and doc,
  whether or not it appears in
  [Affected specs & docs](#account-for-the-specs-and-docs-it-touches) — stays
  untouched until the phase that authorizes it. **Implementation does not begin
  until the change is In-Progress**: advance the `status` first, then implement.
  This is what keeps *what* and *how* settled before any code is written.
- Before setting **In-Review**, create or update every durable spec and doc listed in
  the change's [Affected specs & docs](#account-for-the-specs-and-docs-it-touches)
  section, so the enduring source of truth is ready for review with the
  implementation.
- When implementation is complete, set **In-Review**. Do **not** archive the
  change at this point, and agents must not commit the changes. Leave the parent
  concept and child folder in [`changes/`](../../changes/index.md), and update
  the bundle-root [index](../../changes/index.md) and
  [`log.md`](../../changes/log.md).
- When the work lands after review, set **Done** and **move the change into
  [`archive/`](../../changes/archive/index.md)** — both the
  `NNNN-<slug>.md` parent and its `NNNN-<slug>/` folder — in the same change.
  Update the archive's `index.md` and remove the entry from the bundle-root
  [index](../../changes/index.md).

## Keep OKF tidy

Every concept needs parseable frontmatter and a non-empty `type` — reuse a type
from [`changes/schema.md`](../../changes/schema.md) or add a new one there. After
any add, move, or status change, update the enclosing `index.md` and add a
[`log.md`](../../changes/log.md) entry in the same change. See
[Working with OKF](work-with-okf.md) for the full editing contract.
