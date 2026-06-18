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
for the tool's current behavior. A change bridges the two: durable specs and docs
**MAY** be edited at any time (within a change or on their own), but a change
**MUST** record every durable spec or doc its work impacts in an **Affected specs
& docs** section and **SHOULD** suggest any new durable specs worth creating. It
**SHOULD** bring the durable docs it lists into sync **before** reaching
**In-Review** so the source of truth is not left stale. Completed changes then
move into [`archive/`](archive/); the enduring specs carry the result forward.

## Status lifecycle

A Change's `status` frontmatter moves through, in order:

- **Draft** — writing up the functional spec (the *what*).
- **Design** — working out the technical design (the *how*).
- **In-Progress** — implementing it.
- **In-Review** — implementation complete and ready for review.
- **Done** — landed.

Durable specs and docs **MAY** be edited at any time, with or without a change;
before setting **In-Review**, a change **SHOULD** bring every durable spec and doc
it listed in **Affected specs & docs** into sync. When a change reaches **Done**,
move it (and its child folder) into [`archive/`](archive/) in the same change.

# Open changes

No changes are open right now. The coordinated evaluation-workflow set (the
deterministic CLI writes the records, the skill judges — `0012`–`0016`), the
independent skill rigor pass (`0017`), the experiment-backed report/status
follow-ups (`0018`–`0024`), and the contributor-guide change teaching durable
specs to carry their rationale (`0025`) have all landed and moved to
[`archive/`](archive/).

Completed changes live in [`archive/`](archive/); copy
[`archive/0001-example-change`](archive/0001-example-change.md) as a starting
template for a new one.

# Bundle

- [schema.md](schema.md) - the concept types used in this bundle.

# Subfolders

- [archive/](archive/) - completed changes.
