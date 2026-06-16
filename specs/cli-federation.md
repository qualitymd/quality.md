# Federation: multiple `QUALITY.md` models

Federation operationalizes the format's containment rule across files. An
in-file `targets:` tree and a discovered tree of `QUALITY.md` files are the same
primitive at different altitudes: a child model is grafted as a target subtree.

## Purpose

Large repositories need multiple quality models without losing a single quality
picture. Federation supports:

- decomposition into component-local models;
- progressive disclosure for agents and humans;
- monorepos with system-level and package-level concerns.

A repository with one `QUALITY.md` behaves exactly like the single-model specs. A
repository with several exposes a tree-shaped view.

## Composition Rules

The normative rules live in [`../SPECIFICATION.md`](../SPECIFICATION.md). This
doc applies them to CLI behavior:

- **Ownership stays at the defining file.** Each model owns its own declarations
  and result files. Federation is a reading of the set, not a merge written back.
- **Containment spans downward.** Ancestor target nodes and ancestor models apply
  to descendants unless overridden with rationale.
- **Factor identity is scoped.** A factor declared by an ancestor is visible in
  descendants. A child may refine it by adding requirements, not redefine it.
- **Overrides cross file boundaries.** A child model may suppress or replace an
  inherited requirement at its graft point only with rationale. A stale override
  is diagnostic.
- **Baseline is the rolling root ancestor.** Shipped baseline target trees sit
  outside the repository tree. They are always visible and evaluated; projects do
  not pin away improvements.
- **Shared scale is by reference.** A model that omits `ratings` uses the default
  scale, not an ancestor's scale. A federation that wants one scale references one
  shared scale file.

## Discovery

- A model is a file named exactly `QUALITY.md`.
- The scan root is the invocation directory or repository root, according to the
  command's resolution rules.
- Discovery walks downward, honoring `.gitignore`, configured ignore globs, and
  conventional vendor directories.
- Directory nesting forms the file-level target tree.
- `-f <file>` opts out and operates on exactly one model.

`model show` without `-f` lists the discovered tree. `model show -f <file>` shows
one model's recursive in-file target tree.

## Grafting

A discovered child model is grafted below the nearest ancestor model as a target
subtree. Its apex target becomes the child target at that path. From that graft
point it inherits ancestor factors and requirements, plus any rolling baseline
content.

The graft is not a physical merge. Commands can always report which file owns a
target, factor, requirement, override, result, and rating scale.

## Runs

Federation does not create one giant run. It creates one living run per
`(model, CLI run target)` pair, exactly as [`cli-evaluate.md`](./cli-evaluate.md)
defines.

- Each model's run lives in that model's `.quality/evaluations/<slug>/`.
- A run records results for requirements owned by that model, addressed by
  target-tree locators.
- Inherited ancestor requirements are recorded under their owning model, not
  duplicated into descendants.
- Skills orchestrate the normal loop once per discovered model.

This keeps component evaluations reviewable on their own while the report can
still render the effective inherited context.

## Tree Report

`evaluation report` without `-f` reads each living run and renders a tree:

- each model node shows its own rolled-up verdict;
- inherited ancestor requirements that cover the node are shown as inherited
  context, attributed to the owner;
- secondary-factor appearances are shown without duplicating result records;
- pending or stale results make that node incomplete;
- by default there is no single collapsed federation rating.

`--fail-on <level>` gates per model. The command exits `1` when any model's own
rollup lands at or below the threshold. A spanning ancestor requirement gates the
ancestor model that owns it.

## Output Layout

Runs decentralize; configuration centralizes:

```text
repo/
  QUALITY.md
  .quality/
    config.yaml
    ratings.yaml
    evaluations/<slug>/
  packages/api/
    QUALITY.md
    .quality/evaluations/<slug>/
  packages/web/
    QUALITY.md
    .quality/evaluations/<slug>/
```

The scan-root `.quality/` holds federation config and shared scales. Each model's
own `.quality/` holds its runs.

## `lint` Over A Federation

`lint` without `-f` validates every discovered model, then runs cross-file rules:
mixed rating scales, graft conflicts, cross-file factor redefinitions,
cross-file stale overrides, and source overlap. See
[`cli-lint.md`](./cli-lint.md#federation-rules).

## Open Questions

- Whether `-f` should gain an `--effective` mode that includes inherited ancestor
  context while still writing results to owner runs.
- Whether the aggregated tree report is persisted at the scan root or rendered
  only.
- Whether to offer an opt-in aggregate federation rating.
- Exact scan-root resolution in multi-root workspaces.
