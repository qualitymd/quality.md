---
type: How-to Guide
title: Working with change cases
description: How to create a change case and move it through its lifecycle in the changes/ bundle.
tags: [changes, contributing]
timestamp: 2026-06-18T00:00:00Z
---

# Working with change cases

Incremental work on the repo is tracked in the
[`changes/`](../../changes/index.md) OKF bundle. A **Change Case** is a formal
unit of work: a parent concept that records _why_, _what state_ it's in, and
_which artifacts it touches_ — code, specs, docs, the bundled skill — over a
[functional spec](write-functional-specs.md) (the _what_) and an optional
[design doc](write-design-docs.md) (the _how_). This guide covers creating one
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
   **motivation**, **scope**, and **affected artifacts** (see
   [below](#account-for-the-artifacts-it-touches)), then links to its
   children.
3. **Add the child folder** `changes/NNNN-<slug>/` with an `index.md` listing and
   the detail concepts (below). The parent-concept-plus-folder shape mirrors the
   rest of the repo (e.g. `specs/cli.md` + `specs/cli/`).

Copy [`0001-example-change`](../../changes/archive/0001-example-change.md) and its
folder as a starting template.

## Consult applicable guidance

At the start of a Change Case, and before each status advance, identify the
repository guidance that applies to the work now being planned, specified,
designed, implemented, or reviewed. Read the applicable guides before making or
reviewing the corresponding artifact, and let that guidance shape the Change
Case artifacts, affected-artifact accounting, implementation process, and review
notes.

Some guides are phase-baseline guidance: functional spec work follows
[Writing functional specs](write-functional-specs.md), design work follows
[Writing design docs](write-design-docs.md), and OKF edits follow
[Working with OKF](work-with-okf.md).

Other guides apply based on the work's substance. Start from the
[guide index](index.md), then follow links only for guidance that matches the
work's affected artifacts, user experience, implementation area, or review risk.
Use the affected-artifact sweep to find relevant guidance too; do not treat any
example list as complete. If a guide appears relevant but is not applicable, note
that briefly when the case is broad, risky, or likely to be reviewed on that
basis.

## Write the spec, then the design

Write them in order — settle _what_ the change case must do before working out
_how_.
The two stages line up with the `status` lifecycle below.

1. **Functional spec first** (`spec.md`, `type: Functional Specification`), while
   the change case is **Draft**. Read
   [Writing functional specs](write-functional-specs.md) before creating or
   updating any change case spec — it is required reading, not an optional
   reference. Then follow it: state clear, testable requirements; use BCP 14
   keywords only where they change conformance meaning; specify behavior, not
   implementation; and don't invent requirements. Before advancing to
   **Design**, review the spec against the functional-spec
   [requirement quality bar](write-functional-specs.md#requirement-quality-bar):
   the requirements must be earned, observable, bounded, and verifiable enough
   to design against. Then validate the set as a whole: run the
   [requirement-set check](write-functional-specs.md#the-requirement-set)
   (consistent, complete, able to be validated) and confirm that satisfying every
   requirement would achieve the case's **motivation** — no more, no less. This
   is the validation question ("did we scope the right thing?") that the
   per-requirement verification path ("can we prove each one?") does not ask. Pin
   down the _what_ before moving on.
2. **Design doc next** (`design.md`, `type: Design Doc`), once the spec is
   settled and the change case moves to **Design**. Follow
   [Writing design docs](write-design-docs.md): it answers the spec, so don't
   start it until the spec is stable. Omit it for a case that needs no design
   discussion.

## Account for the artifacts it touches

A Change Case rarely touches one kind of file. It changes **code**, and usually
durable artifacts alongside it: the enduring [`specs/`](../../specs/index.md)
bundle, the repository-root [`SPECIFICATION.md`](../../SPECIFICATION.md), the
[`README.md`](../../README.md), the [`docs/`](../index.md) guides, the bundled
[`skills/`](../../skills/), and install/scaffold files. Its
[functional spec](write-functional-specs.md) states the _delta_ — what this unit
of work must do — and is archived with the case once it lands; the durable
artifacts hold the _cumulative_ source of truth for the tool's current behavior.
A Change Case is the bridge between them: it proposes a delta, and the enduring
artifacts absorb it.

Account for **every** kind of artifact, and don't let the section's name
pre-decide its shape. Framing the impact as "specs and docs" is exactly how code
edits and skill or scaffold files go unlisted: a surface redesign that names only
the spec files it rewrites — and silently leaves the implementation and the
bundled skill that drives it off the list — has under-accounted for its own work.
The question is _"what does this work change?"_, answered across all artifact
kinds, not _"which specs and docs does it touch?"_

**Find them by analysis, not recall.** The list is the _output_ of a deliberate
sweep, not a memory of what you think you touched — the artifact that gets missed
is always the one you forgot existed. Don't author the list from intuition;
_derive_ it. Start from the concrete things the work changes — command and flag
names, type and function names, file and artifact paths, and the names it
_renames away from_ — and search the whole repo for each one across every bundle:
the code (`src/`, `test/`), the [`specs/`](../../specs/index.md) bundle and
[`SPECIFICATION.md`](../../SPECIFICATION.md), the [`docs/`](../index.md) guides,
the bundled [`skills/`](../../skills/), `README.md`, install/scaffold files, and
`CHANGELOG.md`. Follow inbound cross-references too: a renamed or deleted spec
breaks every link to it, and an old command name lingers in examples and `--help`
blocks. Then triage what the sweep returns — separate live artifacts that must
change from historical ones that must _not_ (append-only `log.md` files and
[`archive/`](../../changes/archive/index.md) cases record past state and stay
frozen). For a wide-reaching change, a search agent over the repo is the fast way
to run this sweep. It is cheap; a contract change that lands with its skill or
implementation unlisted is not.

So, up front — alongside the motivation and scope, not discovered at the end —
record that footprint in the parent concept's **Affected artifacts** section,
grouped by kind so it stays skimmable:

- A Change Case must list every artifact its work impacts and that needs
  updating — **code**, durable specs, the format spec, durable docs, the bundled
  skill, install/scaffold files. Keep it honest: an empty kind (and an empty
  whole) must read as a deliberate "no impact of this kind," not an oversight.
- A Change Case that creates a new public, durable, named artifact should either
  add a durable spec for that artifact or explicitly state which existing
  durable spec owns the complete artifact contract and why a new spec is not
  warranted. Workspace-root or root-adjacent documentation artifacts used by
  generated outputs, such as a shared glossary, count as durable named
  artifacts; do not hide their contracts under a neighboring workflow or
  renderer spec by default. For other contracts the work reveals as
  under-specified or worth lifting out of code, raise
  suggestions for any _new_ durable specs worth creating. Suggesting those
  non-blocking specs is enough; actually creating them is not a precondition for
  the case to land. New specs may be 1:1 artifact specs or behavioral component
  specs. Name 1:1 artifact specs with the artifact-spec convention in
  [Writing functional specs](write-functional-specs.md#conventions); name
  behavioral component specs after the capability or component they govern.
  For a large existing spec, do not stop at the first obvious split: inventory
  its major headings and classify each by durable contract so workflows,
  lifecycles, generated artifacts, examples, and shared invariants each land in
  the right home. For example, a fictional "subscription portal" parent spec
  might keep shared account vocabulary while child specs own plan-change
  routing, renewal reconciliation, and the `renewal-summary.json` artifact.
- A Change Case should bring the durable docs it lists into sync before it
  reaches **In-Review**, so the enduring source of truth matches the landed
  behavior. Because durable docs may also be edited independently, the case is
  not the only path to that update — but it must not leave a doc it listed
  silently stale.

Listing an artifact is _accounting_, not license to edit it now. The two move on
different clocks (see [the lifecycle below](#move-it-through-the-lifecycle)):
**code** is gated to **In-Progress**, while durable specs and docs are not locked
to the lifecycle at all — they may be edited in any phase, within a Change Case
or on their own with no case, whenever the current source of truth is wrong,
stale, or incomplete. What a Change Case owes is honest accounting for the full
footprint of its work, not exclusive custody of the edits.

The parent's **Affected artifacts** is the _index_ — the skimmable list of
every artifact the case touches, and the checklist reconciled before
**In-Review**. The _substance_ of a durable **spec** change lives one level down,
in the [functional spec](write-functional-specs.md): which durable specs the
change rewrites, and what they must say differently, belong with the requirements
that drive them — not in the design doc, and not discovered at the end. So each
**requirement** in a change case's functional spec carries its own durable-spec
impact in a subordinate annotation — the named **add** / **modify** / **rename** /
**delete**, or an explicit `none` — so it is traceable which requirements result
in a spec change and which deliberately do not. Those per-requirement annotations
then roll up into the spec's
[`## Durable spec changes`](write-functional-specs.md#durable-spec-changes)
section — **To add** / **To modify** / **To rename** / **To delete**, each
reading a list or an explicit `None` — which gives the at-a-glance footprint and
catches any change no single requirement owns. Both cover the
[`specs/`](../../specs/index.md) bundle and the format spec
[`SPECIFICATION.md`](../../SPECIFICATION.md). A spec rename belongs in **To
rename** (`old → new`), not split across add and delete.

The other artifact kinds — durable _docs_ (README, guides, scaffold) and _code_ —
need no such breakdown; the parent's index is their only home. The **bundled
skill is the exception that hides**: its runtime content under
[`skills/`](../../skills/) is a durable doc, but its functional spec under
[`specs/skills/`](../../specs/index.md) is part of the `specs/` bundle — so a
change to skill _behavior_ belongs in the **Durable spec changes** breakdown like
any other durable spec, even when the runtime skill files change alongside it.

What gets absorbed is not only the _what_ but the **enduring _why_**. A Change
Case's **motivation** and its [design doc](write-design-docs.md)'s durable
rationale are the reasons its requirements exist — and they archive with the
case. Unless they are promoted into the durable spec, the spec inherits the rule
and loses the reason, and a later editor re-litigates a settled lesson or
"simplifies" a rule back into the bug it fixed. So when a Change Case updates a
durable spec, it should lift that rationale into the spec's
**Background / motivation** (the
big-picture _why_) and its per-requirement **annotations** (the fine-grained
_why_), exactly as the [functional-spec guide](write-functional-specs.md)
describes — not just port the functional delta. The design doc stays in the
archive as the fuller record of alternatives and trade-offs; the spec carries
forward the intent and the lessons.

## Move it through the lifecycle

A Change Case's `status` frontmatter advances, in order:

**Draft → Design → In-Progress → In-Review → Done**

- **Draft** while you write up and review the
  [functional spec](write-functional-specs.md) (the _what_); **Design** while you
  work out the
  [design doc](write-design-docs.md) (the _how_); **In-Progress** while it's
  being implemented; **In-Review** after implementation is complete and ready for
  review; **Done** only after the work lands.
- **Modify nothing outside what the current phase authorizes.** This is a
  whitelist, not a blacklist: the question is never "is this file on the Affected
  list?" but "does this phase permit changing it?" The gate governs **code**.
  Before **In-Progress**, the Change Case's own work stays inside its folder,
  `changes/NNNN-<slug>/` — **Draft** produces its `spec.md`, **Design** produces
  its `design.md`, and either can add supplementary files (sketches, examples,
  data fixtures) that support those, so long as they live within that folder.
  **In-Progress** is the first phase that authorizes touching the **code**; all
  code stays untouched until then. **Implementation does not begin until the
  Change Case is In-Progress**: advance the `status` first, then implement. This
  is what keeps _what_ and _how_ settled before any code is written. Durable
  specs and docs are the exception to the gate — they track the _current_ source
  of truth, not this case's unimplemented delta, so they may be edited in any
  phase (or with no Change Case at all; see
  [Account for the artifacts it touches](#account-for-the-artifacts-it-touches)).
  The Change Case's parent concept, its `index.md`, and the bundle
  [`log.md`](../../changes/log.md)/`index.md` are updated in every phase to
  record the step. Because `status` names the **code** clock only, while durable
  specs and docs run on their own (they may advance in any phase), a single
  status word can understate a case mid-flight. When durable artifacts have moved
  ahead of code, say so in the status block — name both clocks (e.g. "Design;
  code not started; durable report specs scaffolded") so a reviewer isn't
  surprised by spec edits the tree already carries.
- Before setting **In-Review**, reconcile the
  [Affected artifacts](#account-for-the-artifacts-it-touches) list with
  reality: every durable spec and doc the case listed should be updated (or
  already independently up to date), so the enduring source of truth matches the
  landed behavior and is ready for review with the implementation. When a Change
  Case updates a durable spec, it should also promote the
  [enduring _why_](#account-for-the-artifacts-it-touches) — the Change Case's
  motivation and the design doc's durable rationale — into the spec's
  **Background / motivation** and per-requirement annotations, so the reason a
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

Use sentence case for active Change Case headings and current log headings.
Preserve proper nouns, formal type names, and archived historical headings.
