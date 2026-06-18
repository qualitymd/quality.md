---
type: Design Doc
title: Evaluation status and report build — design doc
description: How qualitymd inspects a run's renderability and derives report.md and report.json deterministically from its records, with a CI gate.
tags: [evaluation, cli, design, skill]
timestamp: 2026-06-17T00:00:00Z
---

# Evaluation status and report build — design doc

Design behind the
[Evaluation status and report build](../0015-evaluation-report-build.md) change
and its [functional spec](spec.md). The spec fixes *what*
`qualitymd evaluation show-status` and `qualitymd evaluation build-report` must
do; this doc covers *how* the code makes it so, and why this way. The on-disk
record contract these commands read is the
[record format spec](../0012-evaluation-record-format/spec.md); the records they
consume are written by [0014](../0014-evaluation-record-write.md). The
cross-cutting CLI contract — exit categories, `--json`, stdout/stderr split,
terminal rendering — lives in the [CLI spec](../../../specs/cli.md) and is already
realized in `internal/cli`.

## Context

These are the **read side** of the evaluation run. By the time they run, the run
folder, `assessments/*.json`, `analysis/*.json`, and `recommendations/*.md`
already exist (0013, 0014). Two distinct jobs:

- `show-status` is a **non-destructive probe**: read the records, report counts
  and whether the run is renderable, write nothing, and exit `0` even when the
  run is not yet reportable.
- `build-report` is a **deterministic renderer**: read the same records, derive
  `report.md` and `report.json`, write both, and (with `--fail-at-or-below`) map
  the rendered root rating to an exit code for CI.

The spec's hard line is *render, never judge*: every rating, rationale, and *not
assessed* outcome already lives in the records; these commands transcribe them
and **MUST NOT** infer or recompute. That makes both commands pure functions of
the record set, which is what lets idempotency and a stable CI gate hold.

No evaluation code exists yet (no `internal/evaluation/`, no `evaluation`
command — confirmed against the tree). This change is the first to read a run, so
it introduces the run-read layer that 0013/0014 will share. It must reuse the
existing CLI mechanics rather than reinvent them:

- exit categories via `codedError` / `codeFor` in `internal/cli/root.go`
  (`ExitOK 0`, `ExitProblems 1`, `ExitUsage 2`, `ExitInternal 70`), with
  `usageError(...)` and `silentProblems(...)` constructors;
- terminal-vs-verbatim Markdown via `writeMarkdown` / `renderMarkdown` /
  `colorEnabled` in `internal/cli/spec.go` and `style.go` (Glamour on a TTY,
  raw bytes otherwise);
- atomic writes via `document.WriteAtomic`;
- the `--json` + `schemaVersion` + `receipt.Action` result-struct pattern from
  `internal/lint` and `internal/cli/init.go`.

## Approach

### A run-read package below the CLI

Add `internal/evaluation` to own everything below `internal/cli`: loading a run
folder, validating completeness, and rendering the two reports. The CLI layer
(`internal/cli/evaluation.go`, registered from `root.go` like the other
commands) stays thin — argument parsing, `--json`, `--fail-at-or-below`, output
routing, and exit status — exactly as `lint.go` is thin over `internal/lint`.
This mirrors the lint split (rules/model below, CLI on top) and gives 0013/0014's
writers a single run model to reuse rather than each command re-deriving the
folder layout.

The package centers on a loaded run:

```go
// internal/evaluation
type Run struct {
    Path            string
    Assessments     []Assessment     // assessments/*.json, filename order
    Analyses        map[string]Analysis // analysis/<target>.json, keyed by target slug
    Recommendations []Recommendation // recommendations/NNN-*.md, filename order
    Scale           model.RatingScale // resolved from the evaluated model
}

func Load(path string) (*Run, error) // reads + JSON/frontmatter-decodes every record
```

`Load` is the shared front door: it reads the directory, decodes each record
against the 0012 schema (`schemaVersion`, required fields), and returns a typed
`Run` or a `ParseError`-style failure naming the offending file. Decode failures
and unreadable/missing targets are I/O/shape errors, surfaced to the CLI as
`ExitInternal` (70) — they are not "found problems," they are "cannot read the
run." Record structs are exported Go structs whose JSON tags double as the
read contract, the same convention `internal/lint` uses for its result.

### show-status: renderability without writing

`show-status` calls `Load`, then a single completeness check
(`func (r *Run) Renderable() ([]Gap, error)` — see below) and reports. It never
writes. Because the spec wants it to exit `0` even when the run is *not*
reportable — missing records are the *payload*, not a failure — the command
distinguishes two outcomes:

- **Inspectable but incomplete** → exit `0`, with `reportable: false` and the
  gap list. This is the normal "not done yet" answer a skill polls for.
- **Not inspectable** (folder absent, unreadable record, malformed JSON) →
  `ExitInternal` (70). `Load` returning an error drives this; a successful
  `Load` always yields exit `0`.

Output structs follow the lint/init pattern:

```go
type Status struct {
    SchemaVersion int             `json:"schemaVersion"`
    Path          string          `json:"path"`
    Reportable    bool            `json:"reportable"`
    Counts        Counts          `json:"counts"` // assessments/analyses/recommendations
    Gaps          []Gap           `json:"gaps"`   // missing/inconsistent record refs
    NextActions   []receipt.Action `json:"nextActions"`
}
```

Human output goes through the brand styles/glyphs in `style.go`; `--json`
marshals `Status` to stdout. `nextActions` points the caller at the fix (run
`build-report` when reportable, or `add-record` for the named gap).

### Renderability is one shared completeness check

The same predicate decides `show-status.reportable` and whether `build-report`
may proceed — defining it once is what keeps the probe honest (a green
`show-status` must mean `build-report` will succeed). `Renderable` walks the
record graph for the closure conditions the spec enumerates:

- the `analysis/` roll-up is present and complete — every target referenced by
  another record has an `analysis/<target>.json`;
- exactly one analysis record represents the in-scope root target, identified by
  an empty `targetPath`;
- every `assessmentRecords` reference in an analysis record resolves to a present
  `assessments/*.json`;
- every `recommendations/*.md` referenced by an assessment is present and has
  parseable runtime frontmatter.

Each unmet condition becomes a `Gap{kind, ref, detail}`. A *not assessed*
outcome recorded in a record is **not** a gap — it is valid content and renders
as such; the check only flags *absent* or *dangling* records, never a recorded
judgment. `show-status` reports the gaps and exits `0`; `build-report` treats a
non-empty gap list as fatal: `ExitInternal` (70), the offending record named on
stderr, and **no partial report written**.

### build-report: derive both files from one in-memory report

`build-report` builds a single intermediate value and renders both files from it,
so the human and machine reports cannot drift:

```
Load(run) → Renderable (fail 70 if gaps) → Report (in-memory) → {report.md, report.json}
```

`Report` is the assembled Evaluation Report: the in-scope root aggregate rating +
rationale, the scope, per-target results (root first, recursive: each
requirement's findings summary / rating / rationale; each factor + every
sub-factor at every depth; the target's local and aggregate ratings), and the
advice from the recommendation records. It carries *only recorded values* —
ratings and rationales are copied verbatim from `analysis/*.json` and
`assessments/*.json`. *Not assessed* is a first-class state on every rating-
bearing node (mirroring the `notAssessed` flag in the records), rendered
distinctly at every level, never collapsed to a missing or default rating.

Two renderers consume `Report`:

- **report.md** — the human Evaluation Report per
  [`SPECIFICATION.md` → Report](../../../SPECIFICATION.md#report), with
  [Appendix A](../../../SPECIFICATION.md#appendix-a-sample-evaluation-report) as the
  reference rendering. The renderer produces **plain Markdown bytes** and
  `document.WriteAtomic`s them to `report.md`. Terminal styling is *not* applied
  to the file — the file is the deterministic artifact and must be byte-stable.
  (When a future command *displays* a report on a TTY it can route through the
  existing `writeMarkdown`/Glamour path, exactly as `spec` and `models view` do;
  build-report itself only writes the file.)
- **report.json** — the machine rendering conforming to
  [report.json](../0012-evaluation-record-format/spec.md#reportjson):
  `schemaVersion`, root rating + rationale, scope, per-target results. Findings
  are **referenced by record** (`assessments/*.json` path) with only a minimal
  inline summary; full finding detail stays in the assessment records.
  Marshaled with `json.MarshalIndent(..., "", "  ")` (the lint convention) for a
  stable, diff-friendly serialization, then `WriteAtomic`d.

Both writes replace any existing report files; `WriteAtomic` makes each
replacement crash-safe (temp file in the run dir, rename into place).

### Determinism and idempotency

The reports are pure functions of the decoded records, so byte-stability reduces
to pinning every ordering and excluding every non-record input:

- **Target order** is the structural order of the run's target tree (root first,
  recursive), derived from `targetPath` arrays in the records — not map iteration
  order. Within a target, requirements and factors follow that same structural
  order.
- **Recommendation/advice order** is deterministic filename order by
  `recommendations/NNN-<slug>.md`, as the spec mandates (until a contract adds an
  explicit ordering field).
- **Assessment order** within a requirement follows `assessments/NNN-*`
  filename order.
- **No volatile content**: no timestamps, no host paths beyond the run-relative
  references, no rendering jitter. Glamour is *never* in the file-write path, so
  renderer/version drift cannot change the artifact.

Idempotency then falls out: re-running over unchanged records reproduces both
files byte-for-byte, and `report.md`/`report.json` themselves are not inputs to
`Load`, so re-rendering over prior output is a no-op on content. Tests pin this
by rendering twice and diffing the bytes.

### The `--fail-at-or-below` gate maps to existing exit categories

The gate changes only the exit code; the files are always written first. After a
successful render, the command compares the **in-scope root aggregate rating**
(the same rating `report.md` shows) against `<level>`:

1. Resolve `<level>` against the evaluated model's rating scale. The scale is a
   slice ordered best→worst (`model.RatingScale`, normative per
   `SPECIFICATION.md`); comparison is by index, larger index = worse. An unknown
   `<level>` → `usageError(...)` → `ExitUsage` (2), before or after the write is
   immaterial since it is a caller mistake, but resolved up front so a bad gate
   value fails fast.
2. Root rating **at or below** `<level>` (index ≥, i.e. equal or worse) →
   `silentProblems(...)` → `ExitProblems` (1). `silent` is right here: the human
   report already went to its files and the gate line to stderr, so Fang should
   not re-render an error blob — this matches how `lint` returns
   `silentProblems` after printing findings.
3. Root rating strictly better than `<level>` → `ExitOK` (0).
4. Root rating *not assessed* → `ExitProblems` (1): an unrated root cannot clear
   a bar. This reuses the same `silentProblems` path; *not assessed* is treated
   as failing only for the gate, while still rendering as *not assessed* in both
   files.

Without `--fail-at-or-below`, a successful render is always `ExitOK` (0)
regardless of rating. The gate `SHOULD` print, to stderr, which rating it
compared against which `<level>` and the verdict — using the brand
success/error styles — so a CI log explains itself. Reusing `codedError` means
the gate needs no new exit machinery: it just returns the right constructor.

### Skill consumption

`skills/quality/SKILL.md` (updated in the In-Progress phase, not now) stops
hand-authoring `report.md`/`report.json` and instead calls `build-report` after
its `add-record` writes. The skill keeps owning *judgment* (it wrote the
ratings into the records); the CLI now owns *rendering*. `--json` from both
commands gives the skill a structured status/receipt to act on without parsing
human text.

## Alternatives

- **Render report.md and report.json independently, each straight from the
  records.** Rejected: two code paths over the same records is exactly the drift
  the change exists to remove (the skill's hand-authoring had the same flaw).
  Building one in-memory `Report` and rendering both from it makes divergence
  structurally impossible.
- **Compute or roll up ratings in build-report.** Rejected outright by the spec:
  ratings are skill judgment recorded by 0014. Recomputing would let the report
  disagree with the records, and the records are authoritative. build-report
  transcribes.
- **Let build-report tolerate missing records and emit a partial report.**
  Rejected: a partial report is a misleading artifact, and a CI gate over it
  would be meaningless. Incompleteness is surfaced by `show-status` (exit 0,
  `reportable: false`) and is fatal for `build-report` (exit 70, nothing
  written). This keeps the probe/derive split clean.
- **Embed full finding detail inline in report.json.** Rejected: the 0012
  contract says reference findings by record with minimal inline summaries; the
  full detail already lives in `assessments/*.json`. Inlining would duplicate it
  and bloat the machine artifact.
- **Apply Glamour/terminal styling when writing report.md.** Rejected: the file
  must be byte-deterministic and idempotent, and Glamour output varies with
  terminal width, color scheme, and renderer version. The TTY-rendering path
  (`writeMarkdown`) belongs to commands that *display* a report, not to the one
  that *writes* it.
- **Fold `show-status` into `build-report --dry-run`.** Rejected: probing must
  guarantee zero writes and exit `0` on an incomplete-but-readable run, whereas
  build-report's contract is to write and to gate. Two commands keep those two
  contracts unambiguous, and a polling skill wants the read-only one by name.
- **A new exit code for "gate failed."** Rejected: a failed gate *is* the CLI's
  "ran but found problems" category (`ExitProblems` 1), the same code `lint`
  uses. Reusing it keeps CI semantics consistent across commands.

## Trade-offs & risks

- **New `internal/evaluation` package surface.** It is the first run-read code
  and 0013/0014 will lean on its `Run`/`Load`. Keep `Load` and the record
  structs the shared contract; resist letting CLI-only concerns leak below the
  package boundary (the lint design's main lesson).
- **Renderability check duplicates record-relationship knowledge the writers
  also hold.** Worth it: the check is the single definition of "done," shared by
  the probe and the renderer. If 0013/0014 grow their own notion of a complete
  run, consolidate on this predicate rather than forking it.
- **Byte-stability is a standing test obligation.** Map iteration, JSON key
  order, and any incidental ordering can silently break idempotency. Pin order
  explicitly (target tree, filename sequences), marshal with stable indentation,
  and keep a render-twice-and-diff test as the guard.
- **Rating-scale resolution depends on the evaluated model.** The gate needs the
  scale that produced the run's ratings. Resolve it from the run's recorded model
  rather than from the working tree's `QUALITY.md`, so a gate run in CI compares
  against the same scale the records were rated under.
- **report.md must track the SPECIFICATION Report contract.** Appendix A is the
  reference rendering; if the Report phase changes, the renderer follows. Keeping
  report.md a thin projection of the in-memory `Report` localizes that change.

## Open questions

- **Scale source when records and working tree differ.** The design resolves the
  rating scale from the run's recorded model. If a run ever omits an embedded
  scale, the fallback (working-tree `QUALITY.md`, or fail) needs settling against
  the 0013 run-scaffold contract; deferred to whichever change pins the run's
  `model.md` scale capture.
- **Inline finding-summary shape in report.json.** The spec requires "minimal"
  summaries referencing records; the exact fields (locator + category only, vs.
  a one-line observation) can be pinned during implementation against Appendix A,
  and recorded in the durable build-report sub-spec.
