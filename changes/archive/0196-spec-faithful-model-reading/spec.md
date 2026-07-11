---
type: Functional Specification
title: Spec-faithful model reading — functional spec
description: Requirements to restore CLI fidelity to SPECIFICATION.md for source selectors, extension frontmatter, and scalar shapes.
---

# Spec-faithful model reading — functional spec

Companion to the [Spec-faithful model reading](../0196-spec-faithful-model-reading.md)
change case. This spec states _what_ the change must do; a design doc will cover
_how_ once the case reaches Design.

The key words **MUST**, **MUST NOT**, **SHOULD**, and **MAY** are to be
interpreted as described in IETF RFC 2119.

Throughout, "the file's directory" means the directory containing the QUALITY.md
document, which `SPECIFICATION.md` §Document structure defines as the root area's
default source.

## Scope

Covered: runner source resolution (root default, inheritance, globs, unresolved
signalling, symlink-safe walking), lint's treatment of extension frontmatter,
and companion-JSON-schema scalar fidelity. Deferred: typed / resolver-dispatched
non-path selectors, and out-of-tree source refs — see the change case
[scope](../0196-spec-faithful-model-reading.md#scope) and
[considerations](considerations.md).

## Requirements

### Source resolution

**R1 — Root default source.** When the root area declares no `source`, the
runner **MUST** resolve its source to the file's directory and its descendants,
not to an empty bundle.

> _Why:_ `SPECIFICATION.md` §Document structure: the document's location defines
> the default source for the root area. `internal/status` already reports this as
> `SourceStateDefault`; the runner must resolve it.
> _Durable spec:_ modify `specs/evaluation/runner.md` §Source packaging.

**R2 — Ancestor source inheritance.** When a non-root area declares no `source`,
the runner **MUST** resolve its source to the nearest ancestor area that
declares one, falling back to the root default (R1) when no ancestor declares
one.

> _Why:_ `SPECIFICATION.md` §Source resolution. The runner currently uses only
> the area's own declared value.
> _Durable spec:_ modify `specs/evaluation/runner.md` §Source packaging.

**R3 — Glob expansion.** A `source` value that contains glob metacharacters
**MUST** be expanded relative to the file's directory, and its matches packaged
by the same deterministic, sorted, hashed process as a directory walk.

> _Why:_ `SPECIFICATION.md` §Area: "paths and globs resolve relative to the
> containing QUALITY.md file." The runner currently `Stat`s the glob literally
> and reports it missing.
> _Durable spec:_ modify `specs/evaluation/runner.md` §Source packaging.

**R4 — No silent empty evidence.** When a declared `source` selector resolves to
zero readable material, the runner **MUST** surface it as a `source_unavailable`
outcome for the affected work, and **MUST NOT** package an empty bundle and
proceed as though the requirement were assessed against present evidence.

> _Why:_ the failure mode behind the reported runs — areas judged against nothing
> with no signal. `source_unavailable` already exists as a failure category in
> `specs/evaluation/runner.md`.
> _Durable spec:_ modify `specs/evaluation/runner.md` §Source packaging.

**R5 — Symlink-safe walking.** The directory walk **MUST NOT** treat a symlinked
directory (or other non-regular filesystem entry) as a readable file. Non-regular
entries **MUST** be skipped, and the walk **MUST NOT** error on them.

> _Why:_ the `is a directory` crash on committed symlinked skill directories
> (`.claude/skills/quality` → …) that made runs non-resumable. Skipping
> non-regular entries also keeps packaging within the workspace boundary the
> runner already enforces.
> _Durable spec:_ modify `specs/evaluation/runner.md` §Source packaging.

### Extension frontmatter

**R6 — Extension frontmatter is not invalidity.** `qualitymd lint` **MUST NOT**
classify a document as invalid or non-conforming, and **MUST NOT** block
`model.Load`, solely because the document carries extension frontmatter
properties that `SPECIFICATION.md` §Extensions permits. Lint **MAY** still warn
about an unrecognized key as an authoring aid, but an advisory **MUST NOT** be a
conformance error.

> _Why:_ `SPECIFICATION.md` §Extensions ("a document MAY include frontmatter
> properties beyond those defined") and §Conformance. Distinguishing a genuine
> extension from a typo of a known key is a design decision for the design doc;
> the observable requirement is that a conforming extension document lints as
> valid. This aligns lint with the already-open companion JSON schema.
> _Durable spec:_ modify `specs/cli/lint-rules.md` (unknown-key rule default).

### Companion JSON schema

**R7 — Scalar shape fidelity.** The emitted `quality.schema.json` **MUST NOT**
constrain scalar-shaped model properties (including `assessment`, rating-level
`criterion`, and `ratings` override values) to JSON `string`; it **MUST** accept
any non-empty scalar, matching `SPECIFICATION.md` and the linter.

> _Why:_ `SPECIFICATION.md` §Requirement types `assessment` as "a single
> non-empty scalar"; the linter accepts any non-empty YAML scalar. The schema is
> stricter than both, so `assessment: 42` passes `lint` but fails the schema.
> _Durable spec:_ modify `specs/quality-schema-json.md` (scalar shape).

**R8 — Accurate schema claims.** Generated schema annotations (`$comment`,
`title`, `description`) **MUST NOT** claim an enforcement — such as rating-level
ordering — that no tool actually performs.

> _Why:_ the schema `$comment` states ordering is "enforced by `qualitymd
lint`," but no ordering check exists (the rule is semantic and not mechanically
> checkable). The annotation must not overstate.
> _Durable spec:_ modify `specs/quality-schema-json.md` (`$comment` accuracy).

## Durable spec changes

Covers the [`specs/`](../../../specs/index.md) bundle and
[`SPECIFICATION.md`](../../../SPECIFICATION.md). `SPECIFICATION.md` is
deliberately **unchanged** — it is the conformance target, and each rule above
already reads correctly there.

### To add

None.

### To modify

- `specs/evaluation/runner.md` (§Source packaging) — root default source (R1),
  ancestor inheritance (R2), glob expansion (R3), unresolved-selector signalling
  as `source_unavailable` (R4), symlink-safe walking (R5).
- `specs/cli/lint-rules.md` — the unknown-key rule default must not treat
  spec-permitted extension frontmatter as a conformance error (R6).
- `specs/quality-schema-json.md` — scalar-shaped properties accept any non-empty
  scalar, not only strings (R7); remove the inaccurate ordering-enforcement
  claim (R8).

### To rename

None.

### To delete

None.
