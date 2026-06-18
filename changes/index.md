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
**In-Review**. Completed changes then move into [`archive/`](archive/); the
enduring specs carry the result forward.

## Status lifecycle

A Change's `status` frontmatter moves through, in order:

- **Draft** — writing up the functional spec (the *what*).
- **Design** — working out the technical design (the *how*).
- **In-Progress** — implementing it.
- **In-Review** — implementation complete and ready for review.
- **Done** — landed.

Before setting **In-Review**, update every durable spec and doc listed in the
change's **Affected specs & docs** section. When a change reaches **Done**, move
it (and its child folder) into [`archive/`](archive/) in the same change.

# Open changes

A coordinated set sharpening the evaluation workflow — the deterministic CLI
writes the records, the skill judges (`0012`–`0016`) — plus an independent skill
rigor pass (`0017`). All `In-Review`.

- [0012 — Evaluation record format](0012-evaluation-record-format.md) - lift the
  evaluation artifact contract out of the skill prompt into an enduring `specs/`
  spec the CLI writes and the skill consumes. Keystone for `0013`–`0016`.
- [0013 — Evaluation run scaffold](0013-evaluation-run-scaffold.md) -
  `qualitymd evaluation create-run`: create and number a run folder deterministically.
  Depends on `0012`.
- [0014 — Evaluation record write](0014-evaluation-record-write.md) -
  `qualitymd evaluation add-record`: write schema-conformant assessment, analysis, and
  recommendation records from skill-supplied judgment. Depends on `0012`, `0013`.
- [0015 — Evaluation status and report build](0015-evaluation-report-build.md) -
  `qualitymd evaluation show-status` and `qualitymd evaluation build-report`:
  inspect run renderability, derive `report.md`/`report.json` from the records,
  and gate CI with `--fail-at-or-below`. Depends on `0014`.
- [0016 — Skill consumes evaluation CLI](0016-skill-consume-eval-cli.md) - drive
  the CLI for scaffolding, record writes, and reports instead of hand-authoring
  them. Depends on `0013`–`0015`; sibling of `0017`.
- [0017 — Skill rigor and efficiency](0017-skill-rigor-efficiency.md) -
  operationalize effort levels, require verified evidence and pinned locators,
  re-check rating-binding findings, batch writes, allow deep fan-out. Independent
  of the CLI work.

Completed changes live in [`archive/`](archive/); copy
[`archive/0001-example-change`](archive/0001-example-change.md) as a starting
template for a new one.

# Bundle

- [schema.md](schema.md) - the concept types used in this bundle.

# Subfolders

- [archive/](archive/) - completed changes.
