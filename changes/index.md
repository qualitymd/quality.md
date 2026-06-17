---
okf_version: "0.1"
---

# Changes

Incremental work on the `qualitymd` repo, as an OKF knowledge bundle. Each
**Change** is a unit of work: a parent concept (`type: Change`) that records the
motivation and status and links to its **Functional Specification** (what to
build) and **Design Doc** (how, and why that way). A change that needs no design
doc simply omits it.

This bundle is **independent** of the enduring [`specs/`](../specs/index.md)
bundle for now; the relationship between the two may be revisited. Completed
changes move into [`archive/`](archive/).

## Status lifecycle

A Change's `status` frontmatter moves through, in order:

- **Draft** — writing up the functional spec (the *what*).
- **Design** — working out the technical design (the *how*).
- **In-Progress** — implementing it.
- **Done** — landed.

When a change reaches **Done**, move it (and its child folder) into
[`archive/`](archive/) in the same change.

# Open changes

- [0001 — Example change](0001-example-change.md) - placeholder showing the intended shape (`Draft`).

# Bundle

- [schema.md](schema.md) - the concept types used in this bundle.

# Subfolders

- [archive/](archive/) - completed changes.
