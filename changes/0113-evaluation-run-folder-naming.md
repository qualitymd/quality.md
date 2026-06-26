---
type: Change Case
title: Evaluation run folder naming
description: Shorten the Evaluation v2 run-folder tag from quality-eval to eval and make narrowed runs carry the scope's full structural path, while still recognizing existing legacy run folders.
status: Draft
tags: [cli, evaluation, naming, skill]
timestamp: 2026-06-26T00:00:00Z
---

# Evaluation run folder naming

This parent concept captures the why and status; the detail lives in its
children:

- [Functional spec](0113-evaluation-run-folder-naming/spec.md) - what the change
  must do.

## Motivation

Evaluation v2 names each run folder `NNNN-quality-eval`, or
`NNNN-<narrowing>-quality-eval` when the run is narrowed. Two things grate:

1. The constant `quality-eval` tag echoes its own parent directory
   (`.quality/evaluations/`), so in a folder listing it is pure redundancy. The
   tag only earns its keep when the run name travels alone as a handle — for
   example `run: 0003-quality-eval` in a quality-log entry — where a bare number
   would be ambiguous against recommendation and change ids.
2. When a run is narrowed, the scope slug is whatever the caller passed, with no
   convention. A narrowed run reads more descriptively than a full one, and the
   narrowing segment is the genuinely informative part of the name.

This case shortens the tag to `eval` — keeping a self-identifying noun on the
bare handle without the redundant `quality-` — and establishes that a narrowed
run's slug is the scope's **full structural path** (the Area path from the root,
plus the Factor path when scoping to a Factor), hyphen-joined. The run number
stays the identity and `model.md` stays the structural source of truth; the slug
is a human mnemonic that now carries maximal scope context.

This continues the trajectory of the earlier rename from
`NNNN-subject[-<narrowing>]-quality-eval` to `NNNN[-<narrowing>]-quality-eval`.

## Scope

Covered:

- the Evaluation v2 run-folder name grammar produced by
  `qualitymd evaluation create` (the `eval` tag);
- the `--narrowing` slug convention: the scope's full structural path;
- recognition of run folders across `evaluation create`, `evaluation list`,
  `evaluation status`, and next-run-number computation, including continued
  recognition of existing legacy `-quality-eval` (and legacy `subject`/`model`
  prefixed) folders; and
- the durable CLI and `/quality` contracts and bundled-skill guidance that
  describe run-folder names and narrowing slugs.

Deferred / non-goals:

- no migration, rename, or rewrite of existing completed run folders;
- no Area-vs-Factor kind marker or boundary separator inside the slug (path-safe
  slugs offer no boundary character; a marker would only be another ambiguous
  token);
- no change to the run number as the run's identity, nor to `nextRunNumber`
  monotonicity across the format change; and
- no change to structured data layout under `data/`, generated report
  filenames, or report content.

Full-scope runs carry an explicit `full` scope marker (`NNNN-full-eval`), giving
every run a uniform `NNNN-<scope>-eval` shape; `full` is reserved.

## Affected artifacts

### Code

- [ ] `internal/evaluation/create.go` - `nextRunName`: produce `NNNN-full-eval`
      (no narrowing) and `NNNN-<narrowing>-eval` (drop the `quality-` segment).
- [ ] `internal/evaluation/path.go` - `runNameRE`: match the new
      `NNNN-<scope>-eval` shape while still matching legacy `-quality-eval` (and
      legacy `subject`/`model` prefixes); keep `nextRunNumber` counting both old
      and new folders.
- [ ] `internal/evaluation/list.go` - `narrowingFromRunName`: treat `full` as no
      narrowing and return the narrowing segment for both new and legacy run
      names.
- [ ] `internal/cli/evaluation.go` - `--narrowing` flag help: note it takes the
      scope's full structural path.
- [ ] `internal/evaluation/evaluation_test.go` and any `internal/cli` tests that
      assert run names - update `quality-eval` expectations to `eval` and add
      full-path narrowing coverage.

### Format spec

- [ ] None - `SPECIFICATION.md` does not govern Evaluation v2 run-folder names.
      (Deliberate.)

### Durable specs

- [ ] `specs/cli/evaluation-create.md` - state the run-folder grammar
      `NNNN-<scope>-eval` (`full` or the narrowing path), document the
      `--narrowing` full-structural-path convention, and note continued legacy
      recognition.
- [ ] `specs/skills/quality-skill/evaluation.md` - in the create-run workflow
      step, record that the skill passes `--narrowing` the scope's full
      structural path.

Swept and unchanged (recognition is abstracted, no literal `quality-eval` in
spec prose, and narrowing scope semantics are unchanged):
`specs/cli/evaluation-list.md`, `specs/cli/status.md`,
`specs/skills/quality-skill/reporting.md`,
`specs/skills/quality-skill/quality-skill.md`.

### Durable docs / bundled skill

- [ ] `skills/quality/SKILL.md` - update the `run: 0003-quality-eval` handle
      example to `run: 0003-eval`, and the narrowing note to name the full-path
      slug convention.
- [ ] `skills/quality/workflows/evaluate.md` - in the create-run step, construct
      `--narrowing` as the scope's full structural path; refresh any run-name
      example.
- [ ] `skills/quality/resources/cli-quick-reference.md` - reconcile any
      run-name example and the `--narrowing` description.
- [ ] `CHANGELOG.md` - add an Unreleased entry recording the run-folder name
      change (mirroring the prior naming-change note).
- [ ] `evaluation-v2-sketch.md` - reconcile any run-folder examples that show the
      old tag.

### Suggested new durable specs

- The run-folder name grammar is now enforced and parsed across three code sites
  (`create.go`, `path.go` regex, `list.go`). If it grows further, lifting it into
  a small shared run-folder-name contract may be worthwhile. Suggestion only; not
  required for this case.

## Status

`Draft`. Writing the functional spec. Code untouched; no durable specs edited
yet. See the [status lifecycle](index.md#status-lifecycle).
