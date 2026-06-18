---
type: How-to Guide
title: Working with change cases
description: How to create a change case and move it through its lifecycle in the changes/ bundle.
tags: [changes, contributing]
timestamp: 2026-06-17T00:00:00Z
---

# Working with change cases

Incremental work on the repo is tracked in the
[`changes/`](../../changes/index.md) OKF bundle. A **Change Case** is a formal
unit of work: a parent concept that records *why*, *what state* it's in, and
*which durable specs and docs it touches*, over a
[functional spec](write-functional-specs.md) (the *what*) and an optional
[design doc](write-design-docs.md) (the *how*). This guide covers creating one
and moving it through its full lifecycle to review and landing.

Routine prompted edits do not require a Change Case. Use `changes/` only when
the user asks for a Change Case, when continuing an existing `changes/NNNN-*`
item, or when the work needs durable spec/design/review history.

## Create the change case

1. **Pick the next number.** Change cases are numbered sequentially —
   `NNNN-<slug>`, zero-padded to four digits. Look at the highest existing
   `changes/NNNN-*` (and [`archive/`](../../changes/archive/index.md)) and add one.
2. **Add the parent concept** at `changes/NNNN-<slug>.md` with
   `type: Change Case`.
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

Write them in order — settle *what* the change case must do before working out
*how*.
The two stages line up with the `status` lifecycle below.

1. **Functional spec first** (`spec.md`, `type: Functional Specification`), while
   the change case is **Draft**. You **MUST** read
   [Writing functional specs](write-functional-specs.md) before creating or
   updating any change case spec — it is required reading, not an optional
   reference. Then follow it: state requirements with RFC 2119 keywords, specify
   behavior not implementation, and don't invent requirements. Pin down the
   *what* before moving on.
2. **Design doc next** (`design.md`, `type: Design Doc`), once the spec is
   settled and the change case moves to **Design**. Follow
   [Writing design docs](write-design-docs.md): it answers the spec, so don't
   start it until the spec is stable. Omit it for a case that needs no design
   discussion.

## Account for the specs and docs it touches

A Change Case usually changes more than code. Its
[functional spec](write-functional-specs.md) states the *delta* — what this unit
of work must do — and is archived with the case once it lands. The enduring
[`specs/`](../../specs/index.md) bundle, the repository-root
[`SPECIFICATION.md`](../../SPECIFICATION.md), the [`README.md`](../../README.md),
and the [`docs/`](../index.md) guides hold the *cumulative* source of truth for
the tool's current behavior. A Change Case is the bridge between them: it
proposes a delta, and the enduring artifacts absorb it.

Durable specs and docs are **not** locked to the change lifecycle. They **MAY**
be edited at any time — within a Change Case or on their own, with no case at
all — whenever the current source of truth is wrong, stale, or incomplete. What
a Change Case owes is honest *accounting* for the durable impact of its work,
not exclusive custody of the edits.

So, up front — alongside the motivation and scope, not discovered at the end —
record that impact in the parent concept's **Affected specs & docs** section:

- A Change Case **MUST** list every durable spec or doc its work impacts and that
  needs updating. Keep it honest: an empty list must read as a deliberate "no
  durable impact," not an oversight.
- A Change Case **SHOULD** raise suggestions for any *new* durable specs worth
  creating — a contract its work reveals is under-specified or worth lifting out
  of code. Suggesting is enough; actually creating the new spec is not a
  precondition for the case to land.
- A Change Case **SHOULD** bring the durable docs it lists into sync before it
  reaches **In-Review**, so the enduring source of truth matches the landed
  behavior. Because durable docs may also be edited independently, the case is
  not the only path to that update — but it must not leave a doc it listed
  silently stale.

What gets absorbed is not only the *what* but the **enduring *why***. A Change
Case's **motivation** and its [design doc](write-design-docs.md)'s durable
rationale are the reasons its requirements exist — and they archive with the
case. Unless they are promoted into the durable spec, the spec inherits the rule
and loses the reason, and a later editor re-litigates a settled lesson or
"simplifies" a rule back into the bug it fixed. So when a Change Case updates a
durable spec, it **SHOULD** lift that rationale into the spec's
**Background / Motivation** (the
big-picture *why*) and its per-requirement **annotations** (the fine-grained
*why*), exactly as the [functional-spec guide](write-functional-specs.md)
describes — not just port the functional delta. The design doc stays in the
archive as the fuller record of alternatives and trade-offs; the spec carries
forward the intent and the lessons.

## Move it through the lifecycle

A Change Case's `status` frontmatter advances, in order:

**Draft → Design → In-Progress → In-Review → Done**

- **Draft** while you write up the [functional spec](write-functional-specs.md)
  (the *what*); **Design** while you work out the
  [design doc](write-design-docs.md) (the *how*); **In-Progress** while it's
  being implemented; **In-Review** after implementation is complete and ready for
  review; **Done** only after the work lands.
- **Modify nothing outside what the current phase authorizes.** This is a
  whitelist, not a blacklist: the question is never "is this file on the Affected
  list?" but "does this phase permit changing it?" The gate governs **code**.
  Before **In-Progress**, the Change Case's own work stays inside its folder,
  `changes/NNNN-<slug>/` — **Draft** produces its `spec.md`, **Design** produces
  its `design.md`, and either MAY add supplementary files (sketches, examples,
  data fixtures) that support those, so long as they live within that folder.
  **In-Progress** is the first phase that authorizes touching the **code**; all
  code stays untouched until then. **Implementation does not begin until the
  Change Case is In-Progress**: advance the `status` first, then implement. This
  is what keeps *what* and *how* settled before any code is written. Durable
  specs and docs are the exception to the gate — they track the *current* source
  of truth, not this case's unimplemented delta, so they **MAY** be edited in any
  phase (or with no Change Case at all; see
  [Account for the specs and docs it touches](#account-for-the-specs-and-docs-it-touches)).
  The Change Case's parent concept, its `index.md`, and the bundle
  [`log.md`](../../changes/log.md)/`index.md` are updated in every phase to
  record the step.
- Before setting **In-Review**, reconcile the
  [Affected specs & docs](#account-for-the-specs-and-docs-it-touches) list with
  reality: every durable spec and doc the case listed should be updated (or
  already independently up to date), so the enduring source of truth matches the
  landed behavior and is ready for review with the implementation. When a Change
  Case updates a durable spec, it **SHOULD** also promote the
  [enduring *why*](#account-for-the-specs-and-docs-it-touches) — the Change Case's
  motivation and the design doc's durable rationale — into the spec's
  **Background / Motivation** and per-requirement annotations, so the reason a
  requirement exists survives the change's archival rather than dying with it.
- When implementation is complete, set **In-Review**. Do **not** archive the
  case at this point, and agents must not commit the changes. Leave the parent
  concept and child folder in [`changes/`](../../changes/index.md), and update
  the bundle-root [index](../../changes/index.md) and
  [`log.md`](../../changes/log.md).
- When the work lands after review, set **Done** and **move the Change Case into
  [`archive/`](../../changes/archive/index.md)** — both the
  `NNNN-<slug>.md` parent and its `NNNN-<slug>/` folder — in the same edit.
  Update the archive's `index.md` and remove the entry from the bundle-root
  [index](../../changes/index.md).

## Keep OKF tidy

Every concept needs parseable frontmatter and a non-empty `type` — reuse a type
from [`changes/schema.md`](../../changes/schema.md) or add a new one there. After
any add, move, or status change, update the enclosing `index.md` and add a
[`log.md`](../../changes/log.md) entry in the same edit. See
[Working with OKF](work-with-okf.md) for the full editing contract.
