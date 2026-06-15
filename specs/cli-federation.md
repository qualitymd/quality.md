# Federation: multiple `QUALITY.md` models

> **Status:** rewritten for the **resource-based** surface (the inversion — see
> [`cli.md`](./cli.md) and [`skills.md`](./skills.md)). This doc replaces the
> superseded "the CLI runs the agent across the tree" framing: there is no
> cross-tree agent run. Federation is now **one living evaluation run per model**
> over the deterministic `model` / `evaluation` / `result` resources, with the
> skills orchestrating per model. The composition *convention* (ownership, scope,
> narrowing) is normative in the format spec
> ([`../SPECIFICATION.md#federation`](../SPECIFICATION.md#federation)); this doc
> **operationalizes** it across the CLI and skill layer. Field names and on-disk
> shapes are illustrative; genuinely open items are under
> [Open questions](#open-questions).

## Purpose

A single `QUALITY.md` models one system or component. Real repositories — and
monorepos especially — hold many components, each with its own needs, risks, and
owners. **Federation** lets a repository carry many `QUALITY.md` files that
compose into one quality picture, without any of them ceasing to be a standalone,
portable model.

Three use cases motivate it:

- **Decomposition** — split one oversized model into per-component models, each
  legible on its own, instead of one sprawling file.
- **Progressive disclosure** — a reader, agent, or skill loads the root model for
  the system-level picture and descends into a node model only when working in
  that part of the tree. The nearest model carries the local detail; the context
  cost stays proportional to the work — the same economy the Agent Skills format
  gets from loading a skill's body only on activation.
- **Monorepos** — the headline case: a root model for system-wide and
  cross-cutting concerns, plus a model per package or service for component-local
  quality. This is the lineage of nested `AGENTS.md` and nested `DESIGN.md`.

> Federation is opt-in **by presence**. A repository with one `QUALITY.md`
> behaves exactly as the rest of these specs describe; federation is simply what
> the surface does when it finds more than one — one run per model, a
> tree-shaped report.

## The composition convention (recap)

The normative rules live in the format spec
([`../SPECIFICATION.md#federation`](../SPECIFICATION.md#federation)); they are
restated here only because the operational behavior below depends on them. Treat
the spec as authoritative — this doc does not re-derive them as if normative.

- **Ownership stays at the defining file.** No file absorbs another's
  requirements; each `QUALITY.md` remains a complete, standalone model on disk.
  Composition happens at *evaluation* time — a reading of the set — never as a
  merge written back to a file.
- **Scope spans downward.** A model's requirements apply over its whole directory
  subtree, *including the subtrees of descendant models*. A descendant model
  **adds** requirements for its own subtree; it never removes the coverage an
  ancestor projects onto it.
- **`target` narrows.** A requirement with no `target` spans the model's whole
  subtree; an explicit `target` narrows it to the matched set — by selection (a
  positive glob) or by exclusion (a `!`-prefixed entry; see
  [`../SPECIFICATION.md#schema`](../SPECIFICATION.md#schema)). A centralized
  exemption is exactly this: an ancestor requirement excludes a subtree from its
  own `target`, so the exemption is visible on the requirement that owns it.
- **Scales are shared by reference, not inherited.** A model that omits `ratings`
  uses the built-in default four-level scale, not the nearest ancestor's; a tree
  reads commensurably only when its models point `ratings` at one shared scale
  file.

The **effective model at a path** is the union of every ancestor requirement
whose scope still covers that path and the node's own requirements. The two
levers are orthogonal and each does one job: **placement composes (adds);
`target` narrows.**

## Discovery

By default the resource commands resolve a **single** model (`-f, --file`,
default `./QUALITY.md`; see [`cli.md`](./cli.md#conventions)). A **federation
view** is what they expose when `-f` is omitted against a directory holding more
than one `QUALITY.md`: the surface discovers the model set rather than resolving
one file, and reports a tree-shaped result.

- **What counts as a model.** A file named exactly `QUALITY.md`.
- **Scan root.** The directory the command is invoked in (or the repository root
  when run there). Discovery walks the scan root's subtree downward.
- **Ignored paths.** Discovery honors `.gitignore` and the `ignore` globs in
  `./.quality/config.yaml` (see [`cli.md`](./cli.md#the-quality-home)), and skips
  `.git/` and conventional vendor directories (`node_modules/`, …). A model under
  an ignored path is not discovered.
- **The discovered set forms a tree** by directory nesting. The shallowest model
  is the root of the federation — no file is otherwise privileged; "root" just
  means "shallowest" — and each model owns its directory subtree.
- **`-f <file>` opts out of discovery** and operates on exactly that one model,
  as the single-model surface describes.

`lint` (no `-f`) discovers and validates the **whole** set by default — it is
cheap (see [`lint` over a federation](#lint-over-a-federation)). The evaluation
resources address the set **per model** (below): discovery enumerates the models,
but the work, state, and reports stay per model.

`model show` (no `-f`) over a discovered set lists the models in the tree, each
with its scan-root-relative path; `model show -f <file>` inspects one node as
usual. The federation view is a tree *of single-model inspections*, never a merged
super-model.

## Ownership and scope on the resource surface

The convention is operationalized, not re-implemented: each model is evaluated
**over its own scope** as an independent run, and spanning is rendered into the
tree report rather than baked into any node's stored verdict.

- **A model owns its node.** `evaluation create -f <model>` (or the per-model
  step of a federation run) enumerates the requirements **that model owns** — its
  own `requirements`, resolved over its subtree — as the run's `result` set. An
  ancestor's requirements are **not** copied into a descendant's run; ownership
  stays at the defining file, so the descendant's `results/` contains only what it
  owns.
- **Scope spans downward at evaluation time.** An ancestor requirement with no
  `target` (or one whose `target` reaches into a descendant subtree) is evaluated,
  in the **ancestor's** run, over files that physically live under a descendant.
  The verdict is recorded once, under the owner — never duplicated into the
  descendant's run.
- **`target` narrows the owner's scope.** Resolved `target` globs (visible via
  `model show --requirement <path>`) are what bound each requirement's reach;
  exclusions exempt a subtree on the requirement that owns it. Narrowing is a
  property of the owning model's requirement, so it stays inspectable on that
  model.
- **The effective model at a path** is assembled by the **report**, not stored:
  the union of every ancestor requirement whose resolved scope still covers the
  path and the node's own requirements. It is a *reading* of the per-model runs
  (see [The per-model-gated tree report](#the-per-model-gated-tree-report)), which
  is why no run needs to hold another model's verdicts.

This is the deterministic analogue of the old framing: where a single agent run
once spanned the tree, the work is now decomposed into per-model runs whose
verdicts the report re-assembles by scope.

## One living run per model

Federation does **not** create a giant cross-tree run. It creates **one living
per-target run per model**, exactly the run specified in
[`cli-evaluate.md`](./cli-evaluate.md#data-model) — discovery just produces a set
of them:

- Each discovered model gets its own living run under **its own**
  `.quality/evaluations/<slug>/` (resolved relative to that model's directory; see
  [`cli-evaluate.md`](./cli-evaluate.md#on-disk-layout)). A node's run lives beside
  the node it evaluates, so a component's evaluation travels with the component.
- The full lifecycle is per model and unchanged: `evaluation create` enumerates
  the model's owned requirements, `result run` / `result set` record verdicts,
  staleness and `--from` carry-forward work per run, and git history is each run's
  timeline (see [`cli-evaluate.md`](./cli-evaluate.md#run-states)).
- **Skills orchestrate per model.** A skill (`evaluate-quality`,
  `improve-quality-md`) runs the canonical loop
  ([`skills.md`](./skills.md#orchestration-contract)) **once per discovered
  model** — `evaluation create` → `result list --status pending,stale` →
  judge/`set` or `run` → `evaluation report`. Across a federation the skill walks
  the model set (typically root-down so spanning context is loaded before node
  detail), driving one loop per node. There is no special federation loop; there
  is the per-model loop, repeated, plus the aggregating report.
- **Progressive disclosure follows from this.** A skill or reader working in one
  subtree loads the **nearest** model and its run; it pulls in an ancestor only
  for the spanning context the report attributes upward. Context cost stays
  proportional to the work, never the size of the whole tree.

The bash-only, skill-free CI path composes the same way — `evaluation create &&
result run --all && evaluation report` per model — so a federation of
`bash`-only models is evaluable with no skill and no model calls at all (see
[`cli.md`](./cli.md#the-split-deterministic-cli-judgment-in-skills)).

## The per-model-gated tree report

`evaluation report` (no `-f`) over a discovered set renders a **tree-shaped
report**: it reads each model's living run and assembles them by directory
nesting into one federation view. The aggregation is deterministic — it composes
stored per-model verdicts; it never re-evaluates and never calls a model.

- **Per-model rollup and verdict.** Each model's overall rating is the rollup of
  the requirements **it owns**, computed exactly as the single-model report
  specifies ([`cli-evaluate.md`](./cli-evaluate.md#report-rollup)). The tree shows
  one verdict per node.
- **Spanning is shown, not re-gated.** Because an ancestor requirement spans into
  descendant subtrees, the report shows, at each node, which ancestor requirements
  cover it and how they rated — attributed to the **owning** model. Those
  requirements gate under their owner, not the node, so a node's own verdict and
  its spanning context stay distinct and no requirement is double-counted. This is
  what makes a node's *effective* quality readable without re-running anything.
- **No single aggregate rating.** A federation is reported as a tree of per-model
  verdicts, not collapsed to one number. (An opt-in rolled-up rating is an
  [open question](#open-questions).)
- **Incompleteness is per model.** A node whose run still has `pending` / `stale`
  / `errored` results reports as incomplete at that node, exactly as a single-model
  report does; the tree states it rather than scoring around it.

### Gating with `--fail-on`

`evaluation report --fail-on <level>` applies the gate **per model**: a model
trips if its own rolled-up rating lands at or below `<level>` on its scale. Over a
federation the run exits **`1` if any model trips** — *fail the build if any
component fails*. The exit-code meanings are the shared three-code convention
verbatim (see [`cli.md`](./cli.md#exit-codes)): `1` is a gate failure (some node's
quality is below the bar), `2` is a tool failure (e.g. a model that does not
parse — `lint` first), `0` otherwise. Because gating is per owner, a spanning
ancestor requirement gates the **ancestor's** node, never the descendant it
reaches into.

### Output

The on-disk layout follows one principle — **runs decentralize; configuration
centralizes.** There is no single federation run directory: each model's runs live
in its **own** `.quality/`, resolved relative to that model's directory, so a
node's home *is* the single-model layout
([`cli-evaluate.md`](./cli-evaluate.md#on-disk-layout)) and a component's
evaluation travels with the component. What stays at the **scan root** is
federation-level configuration — `config.yaml` (the `ignore` globs, `defaults`)
and any [shared rating scale](#shared-rating-scale) — resolved once for the tree.
The scan-root `.quality/` therefore does double duty: the federation's config home
and, when a root model sits there, that root model's own runs — the same
`.quality/` a lone model already keeps. The federation report is an aggregating
**view** over the per-model runs:

```text
repo/
  QUALITY.md
  .quality/
    config.yaml                       # federation-level config (ignore, defaults)
    ratings.yaml                      # the shared rating scale, referenced by every model
    evaluations/<slug>/               # the root model's own living run
  packages/api/
    QUALITY.md
    .quality/evaluations/<slug>/      # packages/api's living run
  packages/web/
    QUALITY.md
    .quality/evaluations/<slug>/      # packages/web's living run
```

`evaluation report` (no `-f`) renders the tree — each node's verdict and its
spanning attribution — to stdout (and the schema-stable rollup under `--json`),
reading those per-model runs in place. Each node's run remains a complete,
standalone bundle identical to what `-f <that model>` produces, so a federation
decomposes cleanly into its parts and any one component's evaluation is reviewable
on its own. `TODO`: whether the aggregated tree report is also persisted (e.g. at
the scan root) or is render-only — leaning render-only, since the per-model runs
are the committed source of truth.

## `lint` over a federation

`lint` (no `-f`) discovers and validates **every** model in the set, emitting one
findings document whose entries carry the offending model's path. The single-model
rules (see [`cli-lint.md`](./cli-lint.md)) run per model unchanged; federation adds
a small class of **cross-file** rules that only a set can violate — e.g. node
models using **different rating scales** (a tree reads commensurably only on a
shared scale; see [Shared rating scale](#shared-rating-scale)), and
ancestor/descendant `target` globs that **overlap** on the same files. `lint`
stays the cheap structural gate; it forms no opinion on whether the composed model
is *good*.

## Shared rating scale

A federation's tree report is only commensurable when its models share a rating
scale. Sharing is **by reference, not inheritance**: a model's `ratings` may be a
path to a shared scale file (see
[`../SPECIFICATION.md#federation`](../SPECIFICATION.md#federation)), so a
repository defines its scale once (e.g. `./.quality/ratings.yaml`) and every model
points at it. Omitting `ratings` falls back to the built-in default four-level
scale — it does **not** inherit the nearest ancestor's scale; federation adds no
implicit inheritance. `lint` warns when models in one federation use different scales (see
[`lint` over a federation](#lint-over-a-federation)).

## Open questions

- **`-f` and spanning context.** A single-model `-f` run evaluates that model's
  owned requirements only; the spanning ancestor requirements appear only in the
  federation report's tree. Whether `-f` should gain an opt-in
  `--with-inherited` / `--effective` mode (evaluate one node *including* the
  ancestor requirements that span it) is open — and how that would interact with
  ownership staying at the defining file.
- **Persisting the aggregated tree report.** Whether the federation view is
  render-only or also written somewhere at the scan root (above).
- **Aggregate federation rating.** v1 reports a tree of per-model verdicts and
  gates per model. Whether to offer an opt-in single rolled-up rating across the
  tree — and how a child verdict would weigh into it — is deferred.
- **Scan root and multi-root.** Whether the scan root is always the invocation
  directory or resolves to the repository root, and whether a federation may span
  more than one root, is unsettled.
- **Cross-file `target` overlap severity.** Whether overlapping ancestor/
  descendant globs are a `warning` (current lean) or sometimes plainly legitimate
  (a cross-cutting ancestor requirement *intends* to overlap node scopes).
