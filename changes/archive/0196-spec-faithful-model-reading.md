---
type: Change Case
title: Spec-faithful model reading
description: Correct CLI assumptions that narrow abstract QUALITY.md concepts — source selectors, extension frontmatter, scalar shapes — back to what SPECIFICATION.md defines.
status: Done
---

# Spec-faithful model reading

## Motivation

Two evaluation runs (in external projects on `qualitymd` 0.29.0) recorded a
cascade of failed runs and empty-evidence judgments. Diagnosis traced them to a
family of places where the CLI silently narrows an **abstract** concept the
format spec defines to a **concrete, closed** assumption the spec does not
state — then fails, or degrades evidence, when a real document does not fit that
assumption.

The pattern, with `SPECIFICATION.md` as the conformance target:

- **Source is a _selector_** (`SPECIFICATION.md` §Terminology; §Area:
  "paths and globs resolve relative to the containing QUALITY.md file"), but the
  runner treats it as a single, present, literal filesystem path. So:
  a source-less root area resolves to an **empty bundle** instead of the file's
  directory (contradicts §Document structure); a source-less child area does
  **not inherit** its nearest ancestor's source (contradicts §Source
  resolution); globs are never expanded; and a source that resolves to nothing
  is judged against **zero evidence** with no error. A directory walk also
  crashes (`is a directory`) on committed symlinked directories, which is the
  bug that made the runs non-resumable and forced the duplicate run
  directories.
- **Documents MAY carry extension frontmatter** (`SPECIFICATION.md`
  §Extensions, §Conformance), but `qualitymd lint` reports any unknown key as an
  `invalid-frontmatter` **error**, marking a spec-conforming document
  non-conforming and blocking `model.Load`. This also disagrees with the CLI's
  own companion JSON schema, which deliberately stays open
  (`specs/quality-schema-json.md`).
- **Scalar-shaped properties accept any non-empty scalar** (`assessment: 42`
  is valid), but the emitted `quality.schema.json` forces them to JSON
  `string`, and a `$comment` claims an ordering check that no tool performs.

None of these are evaluation-method choices (which the spec leaves open); each
is a claim about what a conforming _document_ means. Restoring fidelity removes
a class of silent evidence loss and spurious invalidity, and it is the
groundwork for treating `source` as a first-class selector later.

## Scope

Covered — bring the CLI's document reading back in line with `SPECIFICATION.md`:

- runner source resolution: root default source, ancestor inheritance, glob
  expansion, an explicit unresolved-selector failure instead of silent empty
  evidence, and symlink-safe directory walking;
- `qualitymd lint`: stop classifying spec-permitted extension frontmatter as a
  conformance error;
- companion JSON schema: accept non-string scalars; remove the inaccurate
  enforcement claim.

Deferred (see [considerations](0196-spec-faithful-model-reading/considerations.md)):

- **Typed / resolver-dispatched selectors** for non-path sources (a prose
  selector like "all specs", a query, a live system). This case makes the
  _path/glob_ selector spec-faithful and stops silent empty evidence; the
  general resolver architecture is a follow-up.
- **Out-of-tree source refs** (`../shared`). The runner deliberately contains
  resolution within the workspace for the source-as-data safety boundary
  (`specs/evaluation/runner.md` §Source packaging). The format permits `../`;
  the runner's narrowing is intentional and stays. Noted, not changed.

## Affected artifacts

- **Code:** `internal/runner/source.go` (symlink-safe walk, glob expansion,
  unresolved-selector signalling), `internal/runner/graph.go` (`areaSource`:
  root default + ancestor inheritance), `internal/lint/rules.go`
  (and `result.go` if severity moves) for the unknown-key rule,
  `internal/schema/jsonschema.go` (scalar shape, `$comment`). Surfacing
  `source_unavailable` may touch `internal/runner/` and
  `internal/evaluation/`. Tests: `internal/runner` source tests,
  `internal/lint/rules_test.go` (`TestSchemaDrivenUnknownKeys` currently
  asserts the behavior being changed), `internal/schema` tests.
- **Durable specs:** modify `specs/evaluation/runner.md` (§Source packaging),
  `specs/cli/lint-rules.md` (unknown-key rule), and
  `specs/quality-schema-json.md` (scalar shape, `$comment` accuracy). See the
  per-requirement annotations in the [functional spec](0196-spec-faithful-model-reading/spec.md).
  `SPECIFICATION.md` is **not** changed: it is the conformance target this case
  restores fidelity to, and it already states each rule correctly.
- **Bundled skill:** none anticipated — evaluation is CLI-owned; confirm at
  In-Progress that no `skills/quality/` source-packaging or lint guidance
  references the old behavior.
- **Docs / generated artifacts:** regenerate `quality.schema.json`; regenerate
  `mintlify/cli.mdx` only if `lint` help text changes. Sweep lint reference
  docs at In-Progress for wording that calls extension frontmatter invalid.

## Children

- [Functional spec](0196-spec-faithful-model-reading/spec.md)
- [Design doc](0196-spec-faithful-model-reading/design.md)
- [Considerations](0196-spec-faithful-model-reading/considerations.md)
