---
type: Design Doc
title: Evaluation record write — design doc
description: How the evaluation add-record subcommands map skill judgment to validated, atomically-numbered evaluation records.
tags: [evaluation, cli, skill, design]
timestamp: 2026-06-17T00:00:00Z
---

# Evaluation record write — design doc

Design behind the [Evaluation record write](../0014-evaluation-record-write.md)
change and its [functional spec](spec.md). The spec fixes _what_
`qualitymd evaluation add-record assessment|analysis|recommendation <run>` must
do — the input channel, numbering, `schemaVersion` stamping, validation, atomic
write, and outputs. This doc covers _how_ the code makes it so, on the
established Go + Cobra + Fang stack, and why a new package owns the record layer
rather than reusing `internal/schema` or `internal/document`.

## Context

`add-record` is the **only** writer of evaluation records, so validation is
inherent: a payload that would produce a non-conformant record is rejected at
write time (see the [spec](spec.md#validation-and-rejection)). The record
contract — fields, `schemaVersion`, run-folder layout, the CLI-writes /
skill-judges split — is fixed by the
[record format spec](../0012-evaluation-record-format/spec.md); the run folder
this writes into is scaffolded by
[change 0013](../0013-evaluation-run-scaffold/spec.md). This change persists one
record per invocation into that already-existing run.

Three forces shape the design:

1. **Records are not QUALITY.md documents.** The existing `internal/schema`,
   `internal/model`, and `internal/document` packages model the _quality model_
   frontmatter — targets, factors, requirements, rating scale. A record is a
   different artifact: a flat JSON judgment object (assessment, analysis) or a
   Markdown-with-frontmatter advice file (recommendation). Their field sets do
   not overlap, so neither `internal/schema`'s structural node schema nor
   `internal/model`'s typed `Spec` validates a record.

2. **The atomic-write helper does not fit.** `document.WriteAtomic` is built for
   _replacing an existing_ QUALITY.md: it `Lstat`s the target (erroring if it is
   absent), preserves the prior mode, and refuses symlinks. Records are new
   files in a fresh run; assessment and recommendation writes must additionally
   _fail_ if the computed `NNN` path already exists (collision detection), which
   a replace-in-place helper cannot express.

3. **The CLI owns numbering and stamping; the skill owns judgment.** The payload
   carries only judgment content and must be rejected if it smuggles a CLI-owned
   field (`schemaVersion`, any `NNN`). That boundary has to be enforced
   structurally, not by trusting input.

## Approach

### A new `internal/evaluation` package owns the record layer

Add an `internal/evaluation` package that owns the record contract end to end:
the typed record shapes, payload decoding, validation, `schemaVersion` stamping,
numbering, filename derivation, deterministic serialization, and the atomic
numbered write. `internal/cli` stays thin wiring — argument parsing, input-source
selection, `--json` routing, and exit status — mirroring how `lint`/`init`
delegate to `internal/lint`/`internal/scaffold`.

This package is deliberately separate from `internal/schema`/`internal/model`:
those define the QUALITY.md _model_; `internal/evaluation` defines the _run
outputs_. It depends on `internal/model` only to read the run's rating scale (see
[rating validation](#rating-validation-reads-the-runs-modelmd)), and on
`internal/receipt` for the shared `Action` type. The dependency runs one way —
`evaluation` may import `model`/`document`/`receipt`; none import `evaluation`.

```text
internal/evaluation/
  record.go      # typed payloads + records (Assessment, Analysis, Recommendation)
  payload.go     # strict JSON decode of the judgment payload
  validate.go    # inherent validation against the record contract
  number.go      # next-NNN scan, slug derivation, filename construction
  write.go       # atomic create-or-fail / create-or-replace, collision retry
  recommendation.go # deterministic Markdown+frontmatter rendering
  result.go      # the add-record receipt struct (the --json contract)
```

### Three subcommands over one shared write pipeline

`evaluation add-record` is a Cobra parent with three child subcommands —
`assessment`, `analysis`, `recommendation` — each taking the run folder as its
sole positional argument. The kind is structural (a subcommand, not a flag)
because the three payloads have different required fields and land in different
subdirectories; see the [spec rationale](spec.md#command-surface). Each
subcommand differs only in which payload type it decodes and which placement
strategy it uses; they share one pipeline:

1. **Resolve and verify the run.** The positional run folder must exist and look
   like a run (its `assessments/`, `analysis/`, `recommendations/`
   subdirectories present, per the
   [run-folder contract](../0012-evaluation-record-format/spec.md#run-folder)). A
   missing or non-run target is an **internal error** (exit 70), per the
   [spec](spec.md#outputs-and-exit-codes) — the writer never scaffolds.
2. **Read the payload** (see [input channel](#input-source-file-or-stdin)).
3. **Decode strictly** into the kind's typed payload, rejecting unknown and
   CLI-owned fields.
4. **Validate** against the record contract (see [validation](#validation-is-inherent)).
   On any failure, return a **usage error** (exit 2) naming the offending field;
   nothing is written.
5. **Place**: derive the subdirectory-local filename — assigning `NNN` for
   assessment/recommendation, the target slug for analysis — and stamp
   `schemaVersion`.
6. **Serialize** deterministically and **write atomically** (see
   [atomic write](#atomic-write-and-numbering-collisions)).
7. **Report** the written path on stderr, or a result receipt on stdout under
   `--json`.

Steps 4–7 are the only kind-specific seams; a small per-kind handler interface
keeps the shared driver in one place:

```go
type recordKind interface {
    Decode(raw []byte) (record, error)   // strict decode + payload-level reject
    Validate(rec record, run *Run) error // contract validation
    Place(rec record, run *Run) (placement, error) // dir-local name + NNN/slug
    Marshal(rec record) ([]byte, error)  // deterministic bytes, schemaVersion stamped
}
```

### Input source: `--file` or stdin

The input rules are entirely the parent-CLI input convention applied here:
`--file <path>` reads that file; `--file -` and an absent flag with a
non-terminal stdin read standard input; an absent flag with a terminal stdin is
a usage error telling the caller to pass `--file` or pipe JSON
([spec](spec.md#input-channel)). The CLI layer resolves the source to a single
`io.Reader` and hands the bytes to the package; terminal detection uses the same
`term.IsTerminal` check the rest of the surface uses on the command's `InOrStdin`.

Decoding is strict and single-document. Use `json.Decoder` with
`DisallowUnknownFields`, decode exactly one value, then assert the stream is at
`io.EOF` — empty input, trailing garbage, and multiple documents are all usage
errors ([spec](spec.md#input-channel)). `DisallowUnknownFields` is what makes a
payload carrying a CLI-owned field (`schemaVersion`, an `NNN`) a structural
rejection: those fields are simply absent from the payload structs, so any
attempt to supply them fails decode rather than being silently honored or
dropped. This is cheaper and harder to bypass than a hand-maintained denylist.

### Validation is inherent

`Validate` enforces the record contract for the record's kind. There is no
separate `validate` command and no path to disk that skips it — validation runs
in the pipeline before any byte is written. It covers at least the spec's
[minimum rejections](spec.md#validation-and-rejection): required fields present;
an assessment's `rating` is `null` whenever `notAssessed` is true; each finding
carries `locator`, `observation`, and `category`; a recommendation carries its
required recommendation fields needed to render frontmatter and body content; and
the `rating` is a level the run defines (next section).

Validation is expressed as plain Go checks over the decoded typed payload,
collecting a clear per-field diagnostic, rather than a generic JSON-Schema engine.
The record field sets are small, fixed by the contract, and kind-specific, so
typed checks are simpler to read, give precise messages, and avoid a schema-doc
artifact that would itself need stamping and drift control. (Reusing
`internal/schema` was considered and rejected — see [Alternatives](#alternatives).)

### Rating validation reads the run's `model.md`

An assessment's `rating` must be a level "defined by the run's rating scale." The
authoritative scale is the model the run was created against, snapshotted by
change 0013 into the run's `model.md`. So rating validation parses
`<run>/model.md` with `document.Parse` + `model.Decode` and checks the rating
against `Spec.RatingScale` level names (plus the contract's allowance for `null`
under `notAssessed`). Reading the per-run snapshot, not the working-directory
`QUALITY.md`, keeps a run internally consistent even if the live model changes
mid-run. A `model.md` that is missing or unparseable is an internal error (the
run is malformed), distinct from a usage error in the payload.

### Numbering, slugs, and placement

Numbering is per-subdirectory and recomputed at write time by scanning the
destination directory — there is no counter file to keep in sync:

- **Assessment / recommendation** get an independent zero-padded `NNN`, one past
  the highest `NNN` already present in `assessments/` / `recommendations/`. The
  scan parses the leading `NNN-` of each matching entry, ignores non-matching
  entries, and takes max+1 (empty directory → `001`). Padding width follows the
  contract's zero-padding; the implementation pads to a fixed width and the scan
  parses leading digits irrespective of width so it stays robust.
- **Analysis** is keyed by target, not numbered: the filename is `<target>.json`.

The slug rule is a single deterministic function shared by every segment that
needs one (`<target>` and `<requirement>` in assessment names, `<target>` in
analysis names, `<slug>` from `title` in recommendation names): lowercase,
collapse each run of non-`[a-z0-9]` to a single hyphen, trim leading/trailing
hyphens ([spec](spec.md#numbering-naming-and-placement)). One function, one
behavior, filesystem-safe everywhere.

### Atomic write and numbering collisions

Writes are atomic from the caller's perspective: write a complete temp file in
the **destination directory** (so the rename is same-filesystem) and rename it
into place. This is the temp-then-rename shape `document.WriteAtomic` already
uses, but with create-semantics that helper cannot express, so
`internal/evaluation` carries its own small writer rather than bending the
document one:

- **Numbered records (assessment, recommendation)** must not clobber. The final
  step uses an atomic exclusive rename target: create the destination via
  `O_CREATE|O_EXCL` (or rename onto a path checked to not exist and treat an
  existing target as a collision). On a detected `NNN` collision — a concurrent
  writer took the number between scan and rename — the command **recomputes the
  next `NNN` once and retries** the rename to the new path. A second collision is
  an **internal error** naming the contended directory
  ([spec](spec.md#validation-and-rejection)). One retry is enough for the
  expected single-concurrent-writer case while still bounding the work and
  surfacing genuinely pathological contention.
- **Analysis records** are keyed by target and intentionally overwrite: the write
  is create-or-replace via temp-then-rename, and the command reports whether it
  **created** or **replaced** the file (the receipt's `created` flag), determined
  by stat-before-rename.

Rejection is atomic by construction: validation and decode happen entirely before
step 5, so a rejected payload consumes no `NNN`, writes no temp file, and leaves
no placeholder. The numbering sequence advances only when a record actually lands.

### Deterministic recommendation rendering

A recommendation record is Markdown with YAML frontmatter, both rendered by the
CLI from the judgment payload. Determinism matters because the file is both
human-reviewed and CLI-read, and identical input must produce identical bytes:

- **Frontmatter** is emitted in a fixed key order — `schemaVersion`, `title`,
  `gap`, `evidenceLocators`, `assessmentRecords`, `remediationOptions`,
  `recommendedOption`, `doneCriterion` — matching the
  [contract's field list](../0012-evaluation-record-format/spec.md#recommendation-record).
  A fixed order (not Go map iteration) is what makes the output stable. Render via
  a `yaml.Node` / ordered struct so key order is pinned.
- **Body** is rendered in a fixed section order — gap, evidence locators,
  remediation options, recommended option, done criterion — under stable headings
  ([contract](../0012-evaluation-record-format/spec.md#recommendation-record)),
  from the same payload fields. The body is generated, not passed through, so the
  skill supplies structured judgment and the CLI owns the human layout.

Assessment and analysis JSON are serialized with a stable, indented encoder and a
fixed top-level field order via struct field order, with `schemaVersion` stamped
as a top-level field.

### CLI wiring and the receipt

`internal/cli/evaluation.go` adds `newEvaluationCmd()` (the `evaluation` parent),
registered from `root.go`, with `add-record` and its three kind subcommands
beneath it. The same parent will later host `create-run` (0013) and
`build-report` (0015), keeping the run lifecycle in one namespace. Each subcommand
owns only flag parsing (`--file`, `--json`), input-source resolution, and mapping
the package's typed errors to the right exit category via the existing
`usageError`/`codedError` machinery (decode/validation → `ExitUsage`; missing or
malformed run, I/O, double collision → `ExitInternal`).

The `--json` receipt is an exported struct in `internal/evaluation` (so the human
and JSON paths render the same value), carrying at least `schemaVersion`, `path`,
`kind`, and — for analysis — `created`, plus optional `nextActions`
([spec](spec.md#outputs-and-exit-codes)). On success without `--json`, the human
confirmation (the written path) goes to **stderr**, leaving stdout clean,
consistent with `init`/`lint` and the
[CLI baseline](../../../specs/cli.md). No success receipt is emitted when the command
refuses to write.

## Alternatives

- **Validate records with `internal/schema`.** Rejected. That package models the
  QUALITY.md frontmatter node schema (targets/factors/requirements/rating
  levels); record field sets are disjoint from it, so it would have to grow a
  parallel, unrelated schema vocabulary. Typed per-kind validation in
  `internal/evaluation` is smaller, gives precise field-level messages, and keeps
  the two artifact families' contracts from entangling.
- **Reuse `document.WriteAtomic` for record writes.** Rejected. It is built to
  _replace an existing_ file (it `Lstat`s the target and errors when absent, and
  refuses symlinks). Records are new files, and numbered records must _fail_ on an
  existing target to detect collisions — the opposite of replace-in-place. A small
  create-or-fail / create-or-replace writer in `internal/evaluation` expresses
  both placements; the temp-then-rename mechanics are shared by convention, not by
  forcing one helper to do both jobs.
- **A persistent counter file per subdirectory for `NNN`.** Rejected. Scanning the
  directory for the highest `NNN` is stateless, self-healing, and matches how 0013
  computes the run number; a counter file is a second source of truth that can
  drift from the actual files and needs its own atomic update.
- **Trust `schemaVersion`/`NNN` from the payload (or strip them).** Rejected. The
  CLI-writes / skill-judges split must be enforced, not honored or silently
  cleaned. `DisallowUnknownFields` over payload structs that simply omit those
  fields turns smuggling into a structural decode rejection — no denylist to
  maintain.
- **A single `add-record --kind` flag instead of three subcommands.** Rejected,
  per the [spec rationale](spec.md#command-surface): the payloads have different
  required fields and destinations, so a structural subcommand keeps each
  contract distinct and the help discoverable.
- **Pass the recommendation Markdown body through from the payload.** Rejected.
  Rendering the body from structured fields in a fixed section order is what makes
  output deterministic and keeps the human layout a CLI concern; a pass-through
  body would let the skill own serialization, against the contract.
- **Read the rating scale from the working-directory `QUALITY.md`.** Rejected. The
  run's `model.md` snapshot is the authority for _that run_; reading the live file
  would let a mid-run model edit silently change what ratings are valid.

## Trade-offs & risks

- **Two atomic-write code paths.** `internal/document` and `internal/evaluation`
  each carry temp-then-rename logic. The duplication is small and the semantics
  genuinely differ (replace-existing vs. create-or-fail / create-or-replace);
  collapsing them into one over-parameterized helper would be harder to read than
  two focused ones. If a third writer appears, factor a shared low-level
  temp-then-rename primitive then.
- **One-retry collision window.** Recompute-once-and-retry handles the expected
  single concurrent writer; a burst of concurrent writers to the same directory
  can still fail the second attempt. That is the spec's chosen behavior (fail
  loudly, naming the contended directory) rather than an unbounded retry loop that
  could mask corruption.
- **Numbering is O(dir entries) per write.** Each numbered write scans its
  destination directory. Runs are small (one assessment per requirement), so this
  is negligible; it is not built for tens of thousands of records.
- **`model.md` parse coupling.** Rating validation depends on the run's snapshot
  being a parseable QUALITY.md. If 0013 ever seeds a non-QUALITY model snapshot,
  rating validation needs revisiting. Treating an unparseable `model.md` as an
  internal error (malformed run) keeps the failure legible.
- **Deterministic-bytes promise needs tests.** Stable frontmatter key order, body
  section order, and JSON field order are easy to regress through a struct or
  encoder change. Pin them with golden-file tests over fixed payloads.

## Open questions

- **Shared next-actions machinery.** Like `init`/`lint`, the receipt's
  `nextActions` and the stderr footer are ad hoc until the
  [CLI spec](../../../specs/cli.md#conventions) settles the shared next-actions
  convention; `add-record` should migrate to it when it lands rather than
  entrench a bespoke footer.
- **Where the create-or-fail / create-or-replace primitive ultimately lives.**
  Started in `internal/evaluation`; if `build-report` (0015) needs the same
  create-new semantics, promote it to a shared low-level writer at that point
  rather than now.
