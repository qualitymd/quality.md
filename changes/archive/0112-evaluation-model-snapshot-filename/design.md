---
type: Design Doc
title: Evaluation model snapshot filename - design
description: Design for renaming the Evaluation v2 run-folder model snapshot to model-snapshot.md via a single shared constant, a clean break with no old-name reader, and a rename of this repo's own tracked dogfood runs.
tags: [cli, evaluation, reports]
timestamp: 2026-06-26T00:00:00Z
---

# Evaluation model snapshot filename - design

This design answers the
[Evaluation model snapshot filename functional spec](spec.md).

## Context

The snapshot filename `model.md` is currently a bare string literal repeated at
five call sites across two packages:

- `internal/evaluation/create.go` writes it (the only writer);
- `internal/evaluation/path.go` validates a run folder by its presence;
- `internal/evaluation/load.go` parses it for the model load path;
- `internal/evaluation/report_v2.go` parses it for report rating-label
  resolution; and
- `internal/status/status.go` reads it for staleness comparison.

There is no shared constant, so the literal — and the rename — has to stay
consistent across both `internal/evaluation` and `internal/status`. `status`
already imports `internal/evaluation` (it consumes `evaluation.RunDir`), so the
two packages can share one source of truth.

## Approach

Introduce a single exported constant in `internal/evaluation` — for example
`ModelSnapshotFile = "model-snapshot.md"` — and have all five call sites
reference it instead of a literal. `internal/status` references the same
constant through its existing dependency on `internal/evaluation`. This removes
the repeated literal that made the current name easy to drift, so a future
rename touches one declaration.

Update the two operator-facing error messages — the `path.go` "missing" message
and the `status.go` "reading" message — to name `model-snapshot.md`. The
`create.go` write error message likewise names the new file.

Update `internal/evaluation/evaluation_test.go` so the seed-layout assertion
checks for `model-snapshot.md` (and, by the existing not-seeded list pattern,
that nothing writes `model.md`).

Rename this repository's two tracked active dogfood snapshot files in place —
`.quality/evaluations/0005-subject-quality-eval/model.md` and
`.quality/evaluations/0006-quality-eval/model.md` — and update the run-local
prose in `0006`'s `design.md` and `plan.md` that names `model.md`. These runs
match the run-name regex and are enumerated by `status`, so leaving them on the
old name would make this repo's own runs fail validation. Frozen runs under
`.quality/evaluations/archive/` are not enumerated (their folder name does not
match the run-name pattern), so they keep their historical `model.md` files
untouched.

Durable specs and the bundled skill workflow get wording-only updates so the
public contract and skill artifact summary name the snapshot consistently.

## Spec Response

The snapshot-filename contract is satisfied by routing the single write and all
four reads through one constant set to `model-snapshot.md`. Because the writer
and every reader resolve the same constant, a freshly created run is immediately
readable, validatable, and stale-comparable.

The no-compatibility requirement is satisfied by not adding any fallback Stat or
Parse on the old name: validation, load, report build, and staleness check Stat
or read exactly one filename.

The operator-message requirements are satisfied by updating the `path.go` and
`status.go` error strings (and the `create.go` write error) to the new name.

The dogfood-run boundary is satisfied by renaming the two tracked active runs
and their run-local prose, while the archive boundary is satisfied because the
CLI never enumerates `archive/` and this change does not touch it.

## Alternatives

### Keep `model.md` and rely on prose

Rejected. The durable specs already gloss it as "the model snapshot," but the
filename is what shows in a file tree, an editor tab, and a `Data`/link target.
The bare name is the one place the snapshot/live distinction is invisible, which
is exactly the ambiguity the change removes.

### Dual-name reader for back-compat

Rejected. A reader that accepts `model.md` or `model-snapshot.md` would let the
old name persist indefinitely and erode the deterministic run layout, against the
Evaluation v2 clean-break stance (0097). Runs are cheap to recreate and the only
on-disk runs that matter — this repo's own — are renamed in this change.

### Leave existing dogfood runs on `model.md`

Rejected. Those two runs are tracked and enumerated by `status`; under a clean
rename they would report as "not an evaluation run folder," breaking this repo's
own dogfooding. Renaming them in the same change keeps the repository
self-consistent. This is distinct from runtime migration of runs in other
repositories, which stays a non-goal.

### Inline the new literal at each site instead of a constant

Rejected. Five literals across two packages are what let the name drift in the
first place. One exported constant makes the contract single-sourced and the next
rename a one-line change.

## Trade-offs & Risks

Existing runs in other repositories keep their `model.md` files and will no
longer validate or report staleness under the new CLI. That is the accepted
clean-break cost; Evaluation v2 does not promise cross-version run
compatibility, and recreating a run is cheap.

Introducing an exported constant slightly widens the `internal/evaluation`
package surface, but the constant is a genuine shared contract between
`internal/evaluation` and `internal/status` and is preferable to a repeated
literal.

## Open Questions

None.
