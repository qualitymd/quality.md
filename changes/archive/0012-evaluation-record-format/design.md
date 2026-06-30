---
type: Design Doc
title: Evaluation record format — design doc
description: Where the enduring evaluation-record contract lives in specs/ and how the CLI and skill consume it.
tags: [evaluation, specs, design]
timestamp: 2026-06-17T00:00:00Z
---

# Evaluation record format — design doc

Design behind the
[Evaluation record format](../0012-evaluation-record-format.md) change and its
[functional spec](spec.md). This is a spec-only change — no code — so the design
is just the placement and registration decisions; the contract itself is in the
spec and is not restated here.

## Context

The evaluation artifact contract currently lives as prose in the
[`/quality` skill prompt](../../../skills/quality/SKILL.md#artifact-contract) — the
folder layout, record fields, and `schemaVersion`. We are moving to a model where
the deterministic CLI **writes** every record and the skill only supplies
judgment, so two surfaces (the CLI implementation and the skill prompt) must read
the same contract. Prose in a prompt cannot be that shared source.

## Approach

Establish the contract as one enduring concept, `specs/evaluation-records.md`, in
the [`specs/`](../../../specs/index.md) bundle alongside the other deterministic-
surface specs (`cli.md`, the `cli/` sub-specs). It sits at the bundle root, not
under `cli/`: it describes the on-disk artifact contract, not one command's
behavior — the commands that produce these records (changes 0013–0015) are the
consumers and will live under `cli/`.

Register it like any other concept in the bundle:

- add a concept type to [`specs/schema.md`](../../../specs/schema.md) — the existing
  `Functional Specification` type fits ("the surface as a whole or a single
  subcommand"), broadened or supplemented to cover the artifact contract if its
  description reads too command-specific; and
- list it under **Specs** in [`specs/index.md`](../../../specs/index.md).

Both edits land when this change reaches **In-Progress** (Design authorizes only
this change's own folder), and are tracked in the parent concept's
[Affected specs & docs](../0012-evaluation-record-format.md#affected-specs--docs).

**The skill seam.** The skill's **Artifact Contract** section becomes a reference
to `specs/evaluation-records.md` rather than an inlined copy. That edit is
deferred to In-Progress and is noted here only to mark where the two surfaces
join: the spec becomes the single source, the CLI implements it, and the prompt
points at it. No skill changes happen in this spec-only change.

## Alternatives

- **Keep the contract as prose in the skill prompt.** Rejected — the contract now
  has two consumers (CLI writer, skill reader). Prompt prose cannot be cited by
  CLI code or tested for drift, and duplicating it invites the prompt and
  implementation to diverge.
- **Put it under `specs/cli/`.** Rejected — it is not a single command's spec.
  The record format is shared across the run lifecycle; the commands that write
  it are the per-command specs that belong under `cli/`.
- **Define it as an OKF bundle / schema for the run folder.** Rejected, and the
  spec says so explicitly: a run's records are raw runtime outputs, not OKF
  concepts. The enduring _spec_ is an OKF concept; the _records it governs_ are
  not.
- **Encode the contract as Go types now.** Rejected for this change — that is
  implementation, which the spec-only scope and the lifecycle defer to 0013–0015.
  An enduring prose-plus-schema spec is the source those changes implement
  against.

## Trade-offs & risks

- A new spec concept without code yet means the contract and its implementation
  can drift until 0013–0015 land. Mitigated by those changes being scoped against
  this spec, and by the skill reference (In-Progress) making the spec the single
  cited source.
- The contract is authored prose, not a generated/machine schema. A second
  divergence-prone consumer could later justify a data-file schema; not warranted
  now.

## Open questions

- Whether the existing `Functional Specification` type in `specs/schema.md`
  should cover the artifact contract as-is or get a dedicated type — resolved when
  the registration edit lands in In-Progress.
