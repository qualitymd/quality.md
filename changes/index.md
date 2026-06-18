---
okf_version: "0.1"
---

# Changes

Incremental work on the `qualitymd` repo, as an OKF knowledge bundle. Each
**Change Case** is a formal unit of work: a parent concept
(`type: Change Case`) that records the motivation and status and links to its
**Functional Specification** (what to build) and **Design Doc** (how, and why
that way). A case that needs no design doc simply omits it.

The `changes/` and enduring [`specs/`](../specs/index.md) bundles play different
roles. A Change Case's **Functional Specification** states the *delta* — what
this one unit of work must do — and is archived with the case once it lands. The
enduring `specs/` bundle and the repository-root
[`SPECIFICATION.md`](../SPECIFICATION.md) hold the *cumulative* source of truth
for the tool's current behavior. A Change Case bridges the two: durable specs
and docs **MAY** be edited at any time (within a case or on their own), but a
Change Case **MUST** record every durable spec or doc its work impacts in an
**Affected specs & docs** section and **SHOULD** suggest any new durable specs
worth creating. It **SHOULD** bring the durable docs it lists into sync
**before** reaching **In-Review** so the source of truth is not left stale.
Completed cases then move into [`archive/`](archive/); the enduring specs carry
the result forward.

## Status lifecycle

A Change Case's `status` frontmatter moves through, in order:

- **Draft** — writing up the functional spec (the *what*).
- **Design** — working out the technical design (the *how*).
- **In-Progress** — implementing it.
- **In-Review** — implementation complete and ready for review.
- **Done** — landed.

Durable specs and docs **MAY** be edited at any time, with or without a Change
Case; before setting **In-Review**, a case **SHOULD** bring every durable spec
and doc it listed in **Affected specs & docs** into sync. When a case reaches
**Done**, move it (and its child folder) into [`archive/`](archive/) in the same
edit.

# Open change cases

- [0026 — Authoring guide replaces meta-model workflow](0026-authoring-guide-remove-meta-model.md)
  - replaces the bundled quality meta-model workflow with a practical
    `QUALITY.md` authoring guide and removes the public model-altitude surfaces.
- [0027 — Modularize quality skill modes](0027-modularize-quality-skill.md)
  - splits setup, wizard, evaluation, and improve procedures into separate skill
    reference files while keeping `SKILL.md` as the router.
- [0028 — Require characterized requirements](0028-require-characterized-requirements.md)
  - requires every requirement to be characterized by at least one factor and
    aligns terminology for direct versus secondary factor associations.

Completed change cases live in [`archive/`](archive/); copy
[`archive/0001-example-change`](archive/0001-example-change.md) as a starting
template for a new one.

# Bundle

- [schema.md](schema.md) - the concept types used in this bundle.

# Subfolders

- [0026-authoring-guide-remove-meta-model/](0026-authoring-guide-remove-meta-model/) - spec and design for replacing the meta-model workflow.
- [0027-modularize-quality-skill/](0027-modularize-quality-skill/) - spec for the skill mode-file refactor.
- [0028-require-characterized-requirements/](0028-require-characterized-requirements/) - spec for mandatory requirement characterization.
- [archive/](archive/) - completed change cases.
