---
type: Design Doc
title: Spec-faithful model reading — design
description: How the runner, linter, and companion schema are changed to read documents as SPECIFICATION.md defines them.
---

# Spec-faithful model reading — design

## Context

Answers the [functional spec](spec.md) (R1–R8) for the
[Spec-faithful model reading](../0196-spec-faithful-model-reading.md) case. Three
independent surfaces each narrow an abstract concept; the design keeps them
independent so they can land and be reviewed separately, in the order
source → lint → schema.

## Approach

### Source resolution (R1–R5) — `internal/runner`, with a shared resolver

The runner currently returns an area's _declared_ source verbatim
(`graph.go:areaSource`) and packages it (`source.go:packageSource`). Two
resolution concerns are simply absent — the root default and ancestor
inheritance — and `internal/status` already implements both correctly
(`appendAreaCoverage` threads `inheritedSource`; `sourceCoverageRow` names the
`declared / inherited / default` states). The two surfaces diverged precisely
because the logic lived in one and not the other.

**Extract one resolver both callers share.** Add a source-resolution helper to
the model layer that, given the spec and an area path, returns the _effective
selector_ and its _state_:

```
effectiveSource(spec, path):
  sel = spec.Source            # root declared
  if sel == "": sel = "."      # R1: root default = the file's directory
  for area in walk(path):      # root → target
     if area.Source != "": sel = area.Source   # R2: nearest declaring ancestor wins
  return sel, state
```

`status` consumes `(selector, state)`; the runner consumes `selector` and hands
it to `packageSource`. Because the runner's `workspaceRoot` is already the file's
directory (`workspace.go:115`), the `.` default resolves to exactly the location
`SPECIFICATION.md` §Document structure defines. Sharing the function means the
two surfaces cannot silently disagree again — the regression is designed out, not
just patched.

**Packaging (`packageSource`) gains glob and symlink handling:**

- **R3 globs.** When the effective selector contains glob metacharacters
  (`*`, `?`, `[`, or `**`), resolve it as a pattern rather than a literal
  `os.Stat` path: walk from the longest non-glob parent directory and keep each
  file whose file-relative path matches the pattern, using a small in-tree
  matcher (segment-wise `*`/`?`/`[]` plus `**` across segments) — **no new
  dependency** (D1). The walk honors `skippedSourceDirs` (so `**/*.md` does not
  drag in `node_modules`), except that a skipped directory named as a **literal
  prefix segment** of the pattern opts back in (so `vendor/**/*.go` is honored as
  an explicit selection) (D2). Matches feed the same sorted, hashed, capped
  bundling as a directory walk, so determinism and hashing are unchanged.
- **R5 symlinks.** In `walkSourceDir`, skip any entry that is not a regular file
  (`!entry.Type().IsRegular()`) _before_ appending it — symlinked directories,
  symlinked files, sockets, and devices are all skipped, so the walk never hands
  a non-file path to `os.ReadFile`. `packageSourceFile` keeps a defensive skip
  for the same case. This is the fix for the `is a directory` crash on committed
  `.claude/skills/quality` → … symlinks.

**R4 no silent empty evidence.** After packaging, if a _declared or inherited_
selector produced zero packaged files, the runner raises `source_unavailable`
(an existing failure category) for that area's work units instead of dispatching
judgment against an empty bundle. The root default (`.`) always contains at least
the QUALITY.md file, so the root never trips this; only a selector that matches
nothing does.

### Extension frontmatter (R6) — `internal/lint`

Today an unknown key becomes `RuleInvalidFrontmatter` at `SeverityError`
(`rules.go:61`), which sets `Valid = false` and blocks `model.Load`
(`lint.go:130`). Two structurally different cases are conflated:

1. a **known** model property placed on the wrong node (e.g. `ratingScale` on an
   area) — already split out as `RuleMisplacedRootKey` (`rules.go:57`) and a
   genuine structural error; **unchanged**; and
2. a key that is **no** model property anywhere — which `SPECIFICATION.md`
   §Extensions permits as an extension property.

For case 2, reclassify the finding to a warning-severity `unknown-key` rule that
does **not** set `Valid = false` and does **not** block `model.Load`. A
conforming extension document now lints valid; the unknown key still surfaces as
an advisory so a typo (`assessmnt:`) is not silently swallowed. This aligns lint
with the companion JSON schema, which is already open
(`specs/quality-schema-json.md`). The `--fix` path is unaffected — unknown keys
have no repair.

### Companion JSON schema (R7, R8) — `internal/schema`

- **R7 scalars.** `propertySchema` renders every `ScalarShape` as
  `{"type":"string"}` (`jsonschema.go:134`). Split on whether the scalar is a
  **name** or a **value**, which the Node model already encodes via `Pattern`:
  a scalar _with_ a `Pattern` is a name/ID and stays `string` + `pattern`; a
  scalar _without_ a `Pattern` (`assessment`, `criterion`, `ratings` values) is a
  content scalar and renders as any non-empty scalar —
  `{"anyOf":[{"type":"string","minLength":1},{"type":"number"},{"type":"boolean"}]}`.
  `assessment: 42` then passes both lint and the published schema.
- **R8 comment.** Drop `rating-level ordering` from the `jsonSchemaComment`
  enforcement list (`jsonschema.go:33`); no ordering check exists. Keep
  `uniqueness`, which `duplicate-level` does enforce. The regenerated
  `quality.schema.json` and its drift test update together.

## Spec response

- R1/R2 fall out of the shared `effectiveSource` resolver; R3/R5 are packaging
  changes; R4 is a post-packaging guard. All five stay inside `internal/runner`
  plus the extracted helper, so the runner's existing determinism/hashing
  invariants (`specs/evaluation/runner.md`) are preserved by construction.
- R6 is satisfied by severity, not by an allowlist, so _any_ extension property
  conforms — matching the spec's open rule rather than a curated set.
- R7/R8 keep the schema derived entirely from the Node definitions, so the
  generator still cannot encode a rule the linter does not.

## Alternatives

- **Default `""`→`.` inside `packageSource` only.** Simpler, but handles neither
  inheritance (R2) nor the status/runner divergence — it buries model semantics
  in the packager and leaves the two surfaces free to drift again. Rejected for
  the shared resolver.
- **R6 via an extension-namespace allowlist** (permit `x-*` / dotted keys, error
  on the rest). Rejected: `SPECIFICATION.md` permits _any_ additional property,
  not only namespaced ones; an allowlist re-closes the schema the spec leaves
  open. Warning severity is the faithful reading.
- **R6 by deleting unknown-key detection.** Rejected: loses the typo authoring
  aid that makes the warning worth keeping.
- **R7 leave scalars as `string`.** Rejected outright — it is the exact false
  negative (valid document, failing published schema) this case exists to remove.

## Trade-offs & risks

- **R1 + R5 are coupled.** Making a source-less root actually walk the file's
  directory (R1) re-introduces the whole-repo walk that hit the symlink; R5 must
  land in the same change so that walk is symlink-safe. They ship together.
- **R4 is intentionally louder.** Models that previously "passed" an area against
  empty evidence (e.g. a `source:` pointing at a not-yet-created path) will now
  fail with `source_unavailable`. That is the correction, but it is a behavior
  change for existing runs; the failure message must name the unresolved selector
  so the fix is obvious.
- **Shared resolver touches `status`.** `status` is already correct, so the risk
  is a refactor regression only; existing status tests pin the behavior.

## Resolved decisions

- **D1 — `**` via an in-tree matcher, no dependency.** The runner already walks
  recursively; matching file-relative paths against the pattern is a small
  addition. Adding doublestar would be the repo's first utility dependency
  against an otherwise minimal set, and supporting `*` without `**` is a
  confusing half-measure. ([Designing Go packages](../../../docs/guides/design-go-packages.md).)
- **D2 — Globs honor the skip list, with a literal-segment opt-in.** A recursive
  glob skips `skippedSourceDirs` by default so `**/*.md` stays sane; a skipped
  directory named as a literal prefix segment (`vendor/**`) is the author
  explicitly selecting it and is walked. Implemented by not calling `SkipDir`
  when the current directory matches a literal segment of the pattern.
- **D3 — Binary-only source is `source_unavailable`.** A selector that packages
  zero readable bytes is treated the same as one that matched nothing: no
  evidence is no evidence, and R4 exists to make that loud rather than silent.
- **D4 — The resolver lives in `internal/model`.** Source resolution is model
  semantics (`SPECIFICATION.md` §Source resolution); `model` is already imported
  by both `runner` and `status`. A new `internal/source` package would be a thin
  re-export that only adds an import edge. `EffectiveSource(spec, path)` returns
  `(selector, SourceState)`, serving `status` (which needs the state) and the
  runner (which needs the selector) from one implementation.
