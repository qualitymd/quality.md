---
type: Design Doc
title: Evaluation run scaffold — design doc
description: How qualitymd evaluation create-run resolves the evaluation directory, numbers a run, and seeds the run folder.
tags: [evaluation, cli, command, design]
timestamp: 2026-06-17T00:00:00Z
---

# Evaluation run scaffold — design doc

Design behind the [Evaluation run scaffold](../0013-evaluation-run-scaffold.md)
change and its [functional spec](spec.md). The spec fixes *what*
`qualitymd evaluation create-run` does; the
[evaluation-record contract](../0012-evaluation-record-format/spec.md) fixes the
run-folder layout, naming, and shared numbering. This doc covers *how* the code
makes it so, on the established Go + Cobra + Fang stack.

## Context

`create-run` is the first command of the deterministic evaluation surface. It
must:

- resolve the evaluation directory by a three-level precedence (`--evaluation-dir`,
  then `.quality/config.yaml`'s `evaluationDir`, then `quality/evaluations/`),
  normalizing the result to a repository-relative path and rejecting absolute or
  escaping paths;
- compute `NNNN` deterministically as one past the highest existing run, across
  **both** altitudes (a single shared sequence), fixing the real collision the
  hand-rolled skill numbering produced;
- create the run folder `NNNN-<altitude>[-<narrowing>]-quality-eval` with its
  `assessments/`, `analysis/`, and `recommendations/` subdirectories, refusing to
  overwrite an existing folder of the same name;
- seed `model.md` (snapshotted by the CLI), `design.md`, and `plan.md`; and
- emit a human confirmation on stderr, or a JSON receipt under `--json`, with
  next actions pointing at change 0014.

The current code has none of this scaffolding. There is no evaluation package,
no `.quality/config.yaml` reader, and no repository-root resolution anywhere
under `internal/`. `document.WriteAtomic` exists but `Lstat`s the target first,
so it replaces existing files rather than creating new ones — it is not the seed
writer this command needs. This change therefore introduces the evaluation
domain rather than extending an existing one.

## Approach

### A new `evaluation` command group, wired thinly

Add `internal/cli/evaluation.go` with `newEvaluationCmd()` — a parent command
carrying the `create-run` subcommand — and register it in `root.go` alongside
`newInitCmd`, `newLintCmd`, `newModelsCmd`, and `newSpecCmd`. The shape mirrors
`newModelsCmd()`, which is already a parent (`models`) with `list`/`view`
children:

```go
func newEvaluationCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "evaluation",
		Short: "Work with QUALITY.md evaluation runs",
	}
	cmd.AddCommand(newEvaluationCreateRunCmd())
	return cmd
}
```

The subcommand owns only flag parsing, `--json` routing, output routing, and
exit status, exactly as `init` and `models view` do. Flags map directly to the
spec: `--altitude` (required), `--narrowing`, `--subject`, `--evaluation-dir`,
plus `--json`. `--altitude` validates against `{subject, model}` and
`--narrowing` against a path-safe-slug check; both surface a usage error
(`usageError`, exit 2) on a bad value, consistent with how `models view` routes
bad input through `usageError`.

### Domain logic lives in a new `internal/evaluation` package

Per [Designing Go packages](../../../docs/guides/design-go-packages.md), the run
folder, its numbering, and its layout are a concept in their own right — they are
named for the evaluation run, not for the command that first emits one. So the
resolution, numbering, and folder-creation logic go in a new
`internal/evaluation` package, keeping `internal/cli` thin. This is the same
split as `init` → `internal/scaffold` and `models` → `internal/models`.

The package's surface is small and deterministic:

```go
// CreateRun resolves the directory, numbers the run, creates the folder and
// subdirectories, seeds the three files, and returns what was created.
func CreateRun(opts Options) (*Run, error)

type Options struct {
	RepoRoot      string // resolved repository root
	EvaluationDir string // --evaluation-dir override, empty if unset
	Altitude      string // "subject" | "model"
	Narrowing     string // optional slug
	Subject       string // --subject path, defaults to QUALITY.md
}

type Run struct {
	Path     string // repository-relative run-folder path
	Number   int
	Altitude string
}
```

The CLI builds `Options`, calls `CreateRun`, and renders `Run` either as a human
line or as a receipt. `Run` is feature-internal; the wire contract is the
receipt struct in `internal/cli` (below).

### Resolving the evaluation directory

Resolution composes three steps the package owns:

1. **Repository root.** Walk up from the working directory to the nearest
   ancestor containing `.git` (a directory or file). This is new shared behavior;
   keep it as a small `evaluation`-package helper for now rather than a premature
   `internal/repo` package — `create-run` is its only consumer. If a second
   command needs it, promote it then.
2. **`evaluationDir` precedence.** Take `--evaluation-dir` when set; else read
   `.quality/config.yaml` at the repo root and use its `evaluationDir` when
   present; else default to `quality/evaluations/`. The config read is a narrow
   typed decode of a single field — `struct { EvaluationDir string`yaml:"evaluationDir"`}` —
   not a general config loader, since that is the only field this change needs.
3. **Normalization and escape check.** Reject an absolute path or one that
   `filepath.Clean`/`filepath.Rel` shows escapes the repo root, returning an
   internal-error-category error (exit 70) with a diagnostic, per the spec. The
   accepted value is repo-relative and cleaned.

The resolved directory is created if absent (`os.MkdirAll`); a missing evaluation
directory is normal for the first run, not an error.

### Deterministic numbering by directory scan

The next number is computed by listing the evaluation directory once, matching
each entry against the run-folder name pattern
`^(\d{4})-(subject|model)(-[a-z0-9-]+)?-quality-eval$`, taking the max of the
captured `NNNN` across **all** altitudes, and adding one. An empty or freshly
created directory yields `0001`. Entries that do not match the pattern are
ignored. Because the sequence is shared, a `subject` run can never reuse a number
a `model` run holds, and vice versa.

This scan-then-create is the deterministic mechanism that fixes the collision:
the number is derived from on-disk state at run time, not tracked separately, so
there is no counter to drift. The folder is then created with
`os.Mkdir` using `O_EXCL` semantics — `os.Mkdir` fails with `EEXIST` if the
computed-name folder already exists. The command does **not** retry with the next
number on collision; per the spec, an existing folder at the computed name
signals concurrent or corrupt state, so it fails with the internal-error category
and leaves the folder untouched. (There is an inherent non-atomic gap between
scan and create under true concurrency; see [Trade-offs](#trade-offs--risks).)

### Seeding the run folder

After the run folder and its three subdirectories are created, seed the files
with plain `os.WriteFile` into the new (own) directory — `document.WriteAtomic`
is for *replacing* existing files and `Lstat`s the target, so it does not apply
to creating files in a folder this command just made.

- `model.md` is snapshotted by the CLI. For `--altitude model`, reuse
  `internal/models` to obtain the bundled meta-model bytes
  (`models.Markdown("quality-meta-model", "")`) — the same content the skill gets
  via `qualitymd models view quality-meta-model`. For `--altitude subject`, read
  the resolved `--subject` file (default `QUALITY.md`); a missing subject file is
  an internal error (exit 70) per the spec.
- `design.md` and `plan.md` are stubs — a minimal heading only. The command
  must not invent their judgment content; the skill authors them.

Per the contract, the run folder is **not** an OKF bundle: the command seeds no
`index.md`/`log.md`/`schema.md`.

### Receipt is the agent contract; reuse `internal/receipt`

Follow the `InitReceipt` pattern. Define the receipt struct in `internal/cli`
(it is the command's wire output) and reuse the shared `receipt.Action` for next
actions:

```go
type CreateRunReceipt struct {
	SchemaVersion int              `json:"schemaVersion"`
	Path          string           `json:"path"`
	Altitude      string           `json:"altitude"`
	NextActions   []receipt.Action `json:"nextActions"`
}
```

Under `--json`, marshal this to stdout (mirroring `writeInitReceipt`); otherwise
print the created run-folder path to stderr so stdout stays clean, per the CLI
baseline. The next action points at recording results — change 0014's
`qualitymd evaluation add-record`. On failure under `--json`, emit a small error
object to stderr the way `init` does with `writeInitError`.

Exit codes map straight onto the existing constants in `root.go`: `ExitOK` (0),
`ExitUsage` (2) via `usageError`, and `ExitInternal` (70) for resolution
failures, collisions, missing subject, and I/O — the default for any non-coded
error returned from `RunE`.

## Alternatives

- **Put the logic in `internal/cli` directly.** Rejected. The numbering,
  resolution, and layout rules are domain logic with their own tests and a clear
  concept boundary; the established pattern (`scaffold`, `models`, `lint`) keeps
  `cli` as wiring. Inlining would make `cli` the de facto owner of the run-folder
  concept.
- **Name the package `internal/run` or fold it into a future `internal/eval`
  report package.** Deferred. `evaluation` names the domain the whole 0012–0015
  surface shares; later commands (`add-record`, `build-report`) can live in the
  same package or sibling packages under it. Choosing `evaluation` now avoids a
  rename when those land.
- **Introduce a general `internal/config` loader for `.quality/config.yaml`.**
  Rejected for this change. Only one field (`evaluationDir`) is read; a narrow
  typed decode is honest about scope. Generalize when a second setting exists,
  per the rule of three for shared behavior.
- **Introduce a shared `internal/repo` root-resolution package.** Deferred for
  the same reason — one consumer today. Keep the `.git`-walk helper local to
  `evaluation` and promote it on the second consumer.
- **Track the next number in a counter file or config.** Rejected. A separate
  counter can drift from the actual folders on disk and reintroduces exactly the
  collision class this change fixes. Deriving `NNNN` from a directory scan keeps
  the on-disk state the single source of truth.
- **Retry with the next free number on a folder-name collision.** Rejected per
  the spec: with a deterministic shared sequence, a collision means concurrent or
  corrupt state, which should surface as an error rather than be silently papered
  over by picking another number.
- **Use `document.WriteAtomic` for the seed files.** Rejected: it is built to
  replace an existing file (it `Lstat`s first and refuses symlinks) and offers no
  benefit when writing fresh files into a folder this command just created.

## Trade-offs & risks

- **Scan-then-create is not atomic under concurrency.** Two `create-run`
  invocations racing on the same evaluation directory could both compute the same
  `NNNN`; the loser's `os.Mkdir` fails `EEXIST` and the command errors out (exit
  70) rather than corrupting state. That is the spec's intended behavior —
  surface the race, don't hide it — but it means the command is safe, not
  serializable. A lock file is deferred until a real concurrent-use need appears.
- **Repository-root detection depends on `.git`.** Resolving the root by walking
  to `.git` is conventional but assumes a Git checkout. A non-Git working tree has
  no root to make paths relative to; the command should fail with a legible
  diagnostic rather than guess. This is acceptable for the tool's expected
  context and can be revisited if non-Git use is specified.
- **Two snapshot sources for `model.md`.** The `model` altitude snapshots the
  bundled meta-model and the `subject` altitude copies the user's file. These are
  genuinely different inputs, so the branch is inherent, not incidental; the risk
  is only that the `model`-altitude snapshot stays in step with what
  `models view` emits — reusing `internal/models` for both keeps them identical.
- **Slug validation is the command's gate.** `--narrowing` flows straight into a
  directory name, so its path-safe-slug check is a correctness boundary, not a
  nicety; an under-strict pattern could let a slug create or escape directories.
  Keep the accepted character class narrow (lowercase, digits, hyphen) and reject
  anything else as a usage error.

## Open questions

- **Shared next-actions and config infrastructure.** This change hand-rolls a
  single-field `.quality/config.yaml` decode and a local `.git`-walk. Both are
  candidates for shared packages (`internal/config`, `internal/repo`) once a
  second consumer exists; recorded here so the promotion is deliberate rather
  than forgotten.
- **Where sibling evaluation commands live.** Changes 0014 (`add-record`) and
  0015 (`build-report`) will share the run-folder layout and `schemaVersion`
  conventions. Whether they extend `internal/evaluation` or sit in sibling
  packages under it is left to those changes' designs; this one only needs to not
  foreclose either.

Settled during design:

- **Package placement.** Domain logic in a new `internal/evaluation` package;
  `internal/cli/evaluation.go` is thin wiring. The receipt struct is the wire
  contract and lives in `internal/cli`, reusing `internal/receipt.Action`.
- **Numbering mechanism.** Directory scan for the max matching `NNNN` across both
  altitudes, plus one; create with `os.Mkdir` and fail on `EEXIST`. No counter
  file, no collision retry.
- **Config and root resolution.** Narrow single-field config decode and a local
  `.git`-walk root resolver, both kept inside `internal/evaluation` until a
  second consumer justifies promotion.
