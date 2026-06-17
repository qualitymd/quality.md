---
type: How-to Guide
title: Proposing a change
description: How to propose and track a unit of incremental work in the changes/ bundle.
tags: [changes, contributing]
timestamp: 2026-06-17T00:00:00Z
---

# Proposing a change

Incremental work on the repo is tracked in the
[`changes/`](../../changes/index.md) OKF bundle. A **Change** is one unit of work:
a parent concept that records *why*, *what state* it's in, and *which durable
specs and docs it touches*, over a [functional spec](write-functional-specs.md)
(the *what*) and an optional [design doc](write-design-docs.md) (the *how*). This
guide covers proposing one and moving it through to done.

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

Copy [`0001-example-change`](../../changes/0001-example-change.md) and its folder
as a starting template.

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
created or updated *as part of the change*, before it reaches **Done**. This is
what keeps the enduring specs from drifting as the change's own spec is archived.

## Move it through the lifecycle

A Change's `status` frontmatter advances, in order:

**Draft → Design → In-Progress → Done**

- **Draft** while you write up the [functional spec](write-functional-specs.md)
  (the *what*); **Design** while you work out the
  [design doc](write-design-docs.md) (the *how*); **In-Progress** while it's
  being implemented.
- Before setting **Done**, create or update every durable spec and doc listed in
  the change's [Affected specs & docs](#account-for-the-specs-and-docs-it-touches)
  section, so the enduring source of truth reflects the landed work.
- When the work lands, set **Done** and **move the change into
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
