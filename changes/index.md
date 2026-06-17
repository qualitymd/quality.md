---
okf_version: "0.1"
---

# Changes

Incremental work on the `qualitymd` repo, as an OKF knowledge bundle. Each
**Change** is a unit of work: a parent concept (`type: Change`) that records the
motivation and status and links to its **Functional Specification** (what to
build) and **Design Doc** (how, and why that way). A change that needs no design
doc simply omits it.

The `changes/` and enduring [`specs/`](../specs/index.md) bundles play different
roles. A change's **Functional Specification** states the *delta* — what this one
unit of work must do — and is archived with the change once it lands. The
enduring `specs/` bundle and the repository-root
[`SPECIFICATION.md`](../SPECIFICATION.md) hold the *cumulative* source of truth
for the tool's current behavior. A change bridges the two: it records the durable
specs and docs it creates or updates in an **Affected specs & docs** section, and
those enduring artifacts are brought into sync **before** the change reaches
**Done**. Completed changes then move into [`archive/`](archive/); the enduring
specs carry the result forward.

## Status lifecycle

A Change's `status` frontmatter moves through, in order:

- **Draft** — writing up the functional spec (the *what*).
- **Design** — working out the technical design (the *how*).
- **In-Progress** — implementing it.
- **Done** — landed.

When a change reaches **Done**, move it (and its child folder) into
[`archive/`](archive/) in the same change.

# Open changes

- [0008 — Describe targets with title and description](0008-target-display-fields.md)
  (`status: Design`) — give every target a recommended `title` and optional
  `description`, and reframe the root as a Model (`ratingScale` + Target
  properties) so `ratingScale` is the one Model-only key by type.

Completed changes live in [`archive/`](archive/); copy
[`archive/0001-example-change`](archive/0001-example-change.md) as a starting
template for a new one.

# Bundle

- [schema.md](schema.md) - the concept types used in this bundle.

# Subfolders

- [archive/](archive/) - completed changes.
