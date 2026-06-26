---
type: Functional Specification
title: Drop the "Evaluation v2" naming - functional spec
description: Requirements for retiring the vestigial "v2" qualifier from the live evaluation surface (specs, format spec, skill, code) and renaming it to plain "Evaluation".
tags: [evaluation, naming, specs, cli, skill]
timestamp: 2026-06-26T00:00:00Z
---

# Drop the "Evaluation v2" naming - functional spec

Companion to the
[Drop the "Evaluation v2" naming](../0116-drop-evaluation-v2-naming.md) change
case. This spec states what the rename must do.

The key words **MUST**, **MUST NOT**, **SHOULD**, and **MAY** are to be
interpreted as described in IETF RFC 2119 and RFC 8174 when, and only when, they
appear in all capitals.

## Background / Motivation

"Evaluation v2" named the replacement evaluation workflow while it was understood
in contrast to the workflow it superseded. The
[clean break](../archive/0097-evaluation-v2-clean-break.md) removed the
predecessor, making v2 the only runtime evaluation workflow. A version qualifier
with no other version to distinguish from is noise: it implies a v1 still exists
and adds a token a reader has to mentally discard everywhere the workflow is
named. Plain "Evaluation" is already the name of the package
(`internal/evaluation`), the command group (`qualitymd evaluation`), and the run
folders. This change makes the prose, spec paths, and code symbols agree.

This is a naming change with no behavioral delta. The risk to manage is collateral
damage: rewriting frozen history, or confusing the cosmetic name "v2" with the
load-bearing `schemaVersion` payload marker. The requirements below both perform
the rename and fence those two hazards.

## Scope

Covered: every occurrence of the "Evaluation v2" / "v2-as-evaluation-name" token
on the *live* evaluation surface — the durable [`specs/`](../../specs/index.md)
bundle, [`SPECIFICATION.md`](../../SPECIFICATION.md), the bundled
[`/quality` skill](../../skills/quality/), the Go implementation, and the
superseded root sketch.

Non-goals (out of scope by design, not omission):

- **No behavioral change.** Run-folder layout, `data/` payloads, report content,
  CLI exit codes and output, and skill workflow steps are unchanged.
- **No `schemaVersion` change.** The JSON `schemaVersion` field is a payload-shape
  marker, not the "v2" name; it is not touched (see the guard requirement below).
- **No history rewrite.** Archived change cases, `log.md` files, and existing
  `CHANGELOG.md` entries stay as written.
- **No spec restructuring.** The evaluation spec tree's contents and boundaries
  are preserved; only its folder/file names and "v2" prose change.

## Assumptions & dependencies

- The predecessor evaluation workflow no longer exists in any live artifact, so
  "Evaluation" is unambiguous without a qualifier.
- `specs/evaluation/` does not currently exist, so the folder rename has no
  collision target.
- The Go package is already named `evaluation` (not `evaluationv2`); only
  in-package identifiers and one filename carry the `v2`/`V2` qualifier.

## Requirements

### Naming on the live surface

- After this change, no live evaluation artifact — durable spec, format spec,
  bundled skill, or Go source — **MUST** refer to the evaluation workflow, its
  protocol, its data, or its reports as "Evaluation v2", "v2 evaluation", or an
  equivalent version-qualified name. Each such occurrence **MUST** be replaced by
  the unqualified "Evaluation" (or the bare noun the sentence needs).

  > Rationale: the qualifier's only meaning was the contrast with the retired
  > predecessor; with that gone it falsely implies a coexisting v1. — 0116

- The replacement **MUST** preserve the existing capitalization register: the
  formal model/workflow term stays "Evaluation" where the source said "Evaluation
  v2", and ordinary prose stays lowercase "evaluation" where the source was
  lowercase.

### Durable spec folder rename

- The `specs/evaluation-v2/` folder **MUST** be renamed to `specs/evaluation/`,
  and its parent concept `specs/evaluation-v2/evaluation-v2.md` **MUST** be
  renamed to `specs/evaluation/evaluation.md`. The folder's internal structure
  (child concepts, `records/`, `reports/`, `routines/`, `index.md`, `log.md`)
  **MUST** be preserved unchanged except for the "v2" naming.

- Every live inbound link to the old `specs/evaluation-v2/**` paths — from the
  [`specs/`](../../specs/index.md) bundle, [`SPECIFICATION.md`](../../SPECIFICATION.md),
  the [`docs/`](../../docs/index.md) guides, and the bundled skill — **MUST** be
  updated to the new path so no live artifact links to a path that no longer
  exists.

  > Rationale: a folder rename is a delete-plus-add; the hazard is dangling
  > inbound links, which the verification path below checks for explicitly. — 0116

### Code symbols

- The Go file `internal/evaluation/report_v2.go` **SHOULD** be renamed to drop the
  `_v2` qualifier, and the package-internal identifiers carrying `v2`/`V2`
  (e.g. `buildV2Report`, `v2RenderableGaps`, `collectV2Artifacts`,
  `renderV2ReportTree`, `v2OutputResult`, `v2ReceiptRating`) **SHOULD** be renamed
  to drop it, provided the result stays unambiguous within the package and does
  not collide with an existing identifier.

  > These are private to `internal/evaluation`, so the rename is not a
  > compatibility surface; `SHOULD`, not `MUST`, leaves the exact target names to
  > implementation where a bare `report` would collide with the existing
  > `report.go`. — 0116

### Guards (no collateral change)

- This change **MUST NOT** alter the `schemaVersion` JSON field, its values, or
  any evaluation payload shape.

- This change **MUST NOT** modify archived change cases under
  [`changes/archive/`](../archive/index.md), any append-only `log.md` file, or
  existing `CHANGELOG.md` release entries. A new `CHANGELOG.md` entry describing
  this change **MAY** be added when it lands.

- The superseded root sketch `evaluation-v2-sketch.md` **MUST** be removed, along
  with any live (non-frozen) reference to it, since the durable
  `specs/evaluation/**` specs are now its source of truth.

### No behavioral regression

- After the change, `go build ./...` and `go test ./...` **MUST** pass, and the
  `qualitymd evaluation` command surface **MUST** behave identically to before
  (same subcommands, flags, exit codes, and output structure) — the rename is
  cosmetic.

  > Verification path: `go build`/`go test`; a repo-wide search for the retired
  > token across live artifacts (see below) returning only frozen-history hits;
  > and a link check that no live file references `evaluation-v2-sketch.md` or an
  > `specs/evaluation-v2/**` path.

## Verification

- A case-insensitive repo search for `evaluation[ _-]?v2` and `v2`-as-name
  tokens, excluding `changes/archive/**`, `**/log.md`, and historical
  `CHANGELOG.md` entries, returns no matches.
- No live file links to `specs/evaluation-v2/**` or to `evaluation-v2-sketch.md`.
- `go build ./...` and `go test ./...` succeed.

## Durable spec changes

### To add

None.

### To modify

- `SPECIFICATION.md` — update the "current Evaluation v2 workflow" reference, its
  link to the renamed spec folder, and the "Evaluation v2 v0" advice note (per the
  naming and folder-rename requirements above).
- `specs/index.md`, `specs/cli.md`, and
  `specs/cli/{evaluation-create,evaluation-data,evaluation-list,evaluation-report,evaluation-status,status,index}.md`
  — update "Evaluation v2" prose and inbound links to the renamed folder (per the
  naming and folder-rename requirements above).
- `specs/skills/quality-skill/{evaluation,reporting,quality-skill,examples/index,workflows/evaluate}.md`
  — update "Evaluation v2" prose and inbound links (per the naming requirement
  above).
- The renamed `specs/evaluation/**` concepts (formerly `specs/evaluation-v2/**`)
  — update in-body "Evaluation v2" titles, prose, and self-references (per the
  naming requirement above).

### To rename

- `specs/evaluation-v2/` → `specs/evaluation/` (whole folder), including
  `specs/evaluation-v2/evaluation-v2.md` → `specs/evaluation/evaluation.md` (per
  the folder-rename requirement above).

### To delete

None. (`evaluation-v2-sketch.md` is a repo-root doc, not a `specs/` concept; its
removal is tracked in the parent's **Affected artifacts** index.)
