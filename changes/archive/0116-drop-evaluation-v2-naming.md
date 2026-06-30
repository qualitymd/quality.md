---
type: Change Case
title: Drop the "Evaluation v2" naming
description: Retire the vestigial "v2" suffix from the live evaluation surface now that v2 is the only evaluation workflow, renaming it to plain "Evaluation".
status: Done
tags: [evaluation, naming, specs, cli, skill]
timestamp: 2026-06-26T00:00:00Z
---

# Drop the "Evaluation v2" naming

A **Change Case** to retire the `v2` qualifier from the active evaluation
surface. Detail lives in the child:

- [Functional spec](0116-drop-evaluation-v2-naming/spec.md) — what the case must do.

No design doc yet; the rename is mechanical enough that the spec may carry it,
but the code grain (symbol and filename renames) may warrant one at **Design**.

## Motivation

"Evaluation v2" was the name for the replacement evaluation workflow while it
coexisted, conceptually, with the workflow it superseded. The
[clean break](archive/0097-evaluation-v2-clean-break.md) made v2 the _only_
runtime evaluation workflow — the predecessor is gone. With nothing to
distinguish it from, the `v2` suffix no longer carries meaning: it is vestigial
naming that reads as "there is also a v1" to anyone new to the code, specs, and
skill. Plain "Evaluation" is unambiguous and is already how the package
(`internal/evaluation`), the command group (`qualitymd evaluation`), and the run
folders are named. This case aligns the prose, spec paths, and code symbols with
that reality.

## Scope

Covered: the live evaluation surface — the durable
[`specs/`](../specs/index.md) bundle (including the `specs/evaluation-v2/`
folder and the CLI and skill specs that link into it), the format spec
[`SPECIFICATION.md`](../SPECIFICATION.md), the bundled
[`/quality` skill](../skills/quality/), the Go implementation, and the superseded
root sketch.

Deferred / out of scope:

- **The `schemaVersion` JSON field and its values.** `schemaVersion` is a
  payload-shape marker for the evaluation JSON contract, not an instance of the
  "v2" _name_; this case does not change the field, its values, or the payload
  shape.
- **Frozen historical records.** Archived change cases
  ([`changes/archive/**`](archive/index.md)), append-only `log.md` files, and
  existing `CHANGELOG.md` release entries record past state under the name in
  force at the time and stay as written. (A new `CHANGELOG.md` entry for _this_
  change is added when it lands.)
- **Restructuring the evaluation specs.** This is a rename, not a reshaping of
  the spec tree's contents or boundaries.

## Affected artifacts

The footprint, by kind. The substance of the durable **spec** edits — what each
spec must say differently, and the folder rename — lives in the spec's
[Durable spec changes](0116-drop-evaluation-v2-naming/spec.md#durable-spec-changes)
section; this index is the skimmable checklist.

### Code

- [ ] `internal/evaluation/report_v2.go` — file rename, plus the `v2`/`V2`
      symbols it defines (`buildV2Report`, `v2RenderableGaps`, `collectV2Artifacts`,
      `renderV2ReportTree`, `v2OutputResult`, `v2ReceiptRating`, …).
- [ ] `internal/evaluation/report.go` — `BuildReportReceipt` doc comment,
      `buildV2Report` call site.
- [ ] `internal/evaluation/{create,data,load,path,types}.go` — "Evaluation v2"
      prose in doc comments and labels.
- [ ] `internal/cli/evaluation.go` — "Evaluation v2" in command `Short` help and
      error strings.

### Durable specs (see spec's Durable spec changes)

- [ ] `specs/evaluation-v2/**` — folder rename to `specs/evaluation/**` and
      in-body name updates.
- [ ] `specs/cli.md`, `specs/cli/{evaluation-create,evaluation-data,evaluation-list,evaluation-report,evaluation-status,status,index}.md`,
      `specs/index.md` — prose and inbound links to the renamed folder.
- [ ] `specs/skills/quality-skill/{evaluation,reporting,quality-skill,examples/index,workflows/evaluate}.md`
      — prose and inbound links.

### Format spec

- [ ] `SPECIFICATION.md` — the "current Evaluation v2 workflow" reference and its
      link to the renamed spec folder; the "Evaluation v2 v0" advice/recommendations
      note.

### Durable docs

- [ ] `docs/guides/write-functional-specs.md` — the "Evaluation v2 routine-output
      examples" reference and its link into the renamed folder.
- [ ] `evaluation-v2-sketch.md` — superseded root sketch; remove it (and any live
      inbound references), since the durable `specs/evaluation/**` specs replaced it.
- [ ] `CHANGELOG.md` — _no edits to historical entries_; add a new entry when the
      change lands.

### Bundled skill

- [ ] `skills/quality/{SKILL.md,workflows/evaluate.md,resources/cli-quick-reference.md}`
      — "Evaluation v2" prose in the runtime skill.

### New durable specs worth creating

- None. The rename does not reveal an under-specified contract.

## Status

`Done`. Implemented, verified, and archived. The live surface now uses plain
"Evaluation"; the old spec folder path, superseded sketch, and private `v2` Go
report symbols were removed from live artifacts.
