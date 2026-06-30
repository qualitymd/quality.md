---
type: Design Doc
title: lint command implementation — design doc
description: How qualitymd lint parses, traverses, and reports findings against a shared model.
tags: [cli, command, lint, design]
timestamp: 2026-06-17T00:00:00Z
---

# lint command implementation — design doc

Design behind the [Implement the lint command](../0003-implement-lint-command.md)
change and its [functional spec](spec.md). The durable
[`qualitymd lint` sub-spec](../../../specs/cli/lint.md) fixes the command-specific
behavior; this doc covers how the code should make it so while leaving room for
future deterministic model queries and writes.

## Context

`lint` needs more than the current `internal/spec.Load` path can provide. `Load`
decodes directly into typed structs, validates by returning the first error, and
then hands callers only a valid model. That is the right shape for consumers that
need a conforming model, but `lint` needs to:

- report multiple findings in one run;
- attach each finding to a stable `modelPath`, with source line/column when
  available;
- keep running non-dependent rules after one rule fails;
- report warnings as well as errors;
- apply fixable findings without corrupting authored content; and
- traverse targets, factors, and requirements in parent/child and
  ancestor/descendant relationships.

The same traversal needs will show up again in future query commands: listing
targets, finding a requirement's factors, walking ancestors or descendants, and
scoping work to self, children, descendants, or self-and-descendants. `lint`
should be the first consumer of that shared read model, not a one-off validator.

## Approach

### Split parsing from validation

Keep `internal/spec` as the package that owns the shared `QUALITY.md` document
model for this change. Do not introduce `internal/model` yet; the existing
package already owns the format vocabulary, and a package rename would add churn
before query/write commands exist. Split `internal/spec`'s responsibilities:

- `Parse(path)` reads the file, extracts a required frontmatter block, parses it
  as `yaml.Node`, and builds a document model with as much structure as is
  safely available.
- `Render(document)` renders the repaired frontmatter plus the original Markdown
  body bytes.
- `WriteAtomic(path, bytes, mode)` writes a complete replacement through a temp
  file in the target directory, preserves the original file mode where possible,
  and renames it over the target path.

`internal/spec` owns the document layer only — parsing, rendering, and atomic
writes — and holds **no rule logic and does not import `internal/lint`**.

`internal/lint` owns the lint rule catalog, reporting, and the "valid model"
convenience. It consumes the parsed `spec.Document` instead of reparsing YAML,
and owns rule-level `RepairOp`s:

- `lint.Load(path)` replaces the old `spec.Load`: it calls `spec.Parse`, runs the
  rule catalog, and returns an error when any error-severity finding exists. This
  is the convenience API for callers that need a valid model (today only the
  scaffold conformance test); they import `lint` rather than keep a second
  validator.

The dependency runs one way — `internal/lint` imports `internal/spec`, never the
reverse — so there is a single rule implementation and no `spec`↔`lint` import
cycle. Routing the valid-model check through the same rule catalog is what keeps
parser and rule behavior from drifting: `lint`, future query commands, and
callers that need a valid model all share one model and one set of rules.

### Preserve YAML identity beside typed values

The parsed document should retain both semantic values and source identity:

```go
type Document struct {
	Path        string
	Frontmatter *yaml.Node
	Body        []byte // original Markdown body, preserved verbatim
	Model       *Model
}

type NodeRef struct {
	Kind      NodeKind
	Name      string
	ModelPath []PathSegment
	Line      int
	Column    int
	Parent    *NodeRef
}
```

The exact Go types can vary, but every target, factor, requirement, rating level,
and relevant property needs a stable `modelPath` plus optional YAML source
position. `modelPath` is the machine contract from the lint spec; line and column
are advisory. Missing-key findings attach to the path where the key would appear.
The body is not a type detail to leave to the implementer: `Document` retains the
original Markdown body bytes so `Render` can reattach them byte-for-byte, as the
repair contract requires.

Build the model from `yaml.Node` rather than only from struct decoding. Struct
decoding is still useful once shapes are known, but lint rules need to distinguish
absent keys, empty values, wrong YAML shapes, unknown keys, duplicate level
names, and blocked child checks without losing source locations.

### Add traversal helpers now, selectors later

The first shared traversal API should be deliberately small:

- walk targets root-first;
- walk factors within each target, including sub-factors;
- walk requirements with their declaring target and primary factor, if any;
- expose parent, ancestors, children, and descendants for targets and factors;
- resolve a secondary factor name against the requirement's declaring target and
  ancestor targets; and
- provide deterministic child order.

That is enough for the current lint rules:

- `empty-target` needs target descendants;
- `empty-factor` needs factor descendants and tagged secondary requirements;
- `unknown-factor` needs ancestor-target factor lookup;
- location ordering needs a stable traversal order.

Do not add a full query language or selector grammar in this change. Future query
commands can layer selection modes such as `self`, `children`, `descendants`,
`ancestors`, and `self-and-descendants` on top of these primitives once their
user-facing contract exists.

### Add a narrow repair writer

`--fix` should be implemented as a small repair layer over the parsed YAML
document, not as a general model editing API. A repair is an operation associated
with a specific finding location:

```go
type RepairOp struct {
	RuleID   string
	Message  string
	Location lint.Location
	Apply    func(*spec.Document) error
}
```

Rules that are marked fixable produce repairs alongside their findings. The
initial writer only needs one repair shape: remove an empty optional property
without touching unrelated YAML nodes or the Markdown body.

The repair flow is:

1. Parse and lint the original file.
2. Refuse `--fix` when the target path is a symbolic link.
3. Collect all repairs for findings whose `fixable` value is `true`.
4. Apply repairs to the in-memory YAML document, detecting conflicts before any
   disk write.
5. Write a complete replacement to a temporary file in the same directory, then
   replace the target path.
6. Parse and lint the repaired file again.
7. Render the post-repair result. The `findings` array and every `summary` count
   except `fixed` — including `fixable` — reflect the post-repair lint from step
   6; only `summary.fixed` and `repairs` come from the repairs applied in step 3.
   A clean `--fix` therefore reports `fixable: 0` and `fixed: N`.

This gives `lint --fix` transactional behavior per file: either every compatible
repair is written, or the original file is left alone. It also avoids inventing
patch output, full-file stdout output, or arbitrary model mutation commands in
this change.

If conflict detection or the atomic replacement fails, return a command error
instead of rendering a normal lint result; the original file remains the source
of truth, and there is no post-repair result to report.

Body preservation is strict: split the original file into frontmatter and body,
rewrite only the frontmatter, and append the original body bytes unchanged. YAML
preservation is best-effort around the edited nodes: keep map order, comments,
and scalar style that `yaml.Node` can preserve, and avoid reordering unrelated
keys.

Only `internal/spec` writes bytes to disk. `internal/lint` decides _what_ repair
operations exist; `internal/spec` owns rendering and atomic replacement. That
keeps lint rules from growing filesystem behavior and gives future write
commands a single low-level write path to reuse or replace.

### Implement lint as rule visitors

Create an `internal/lint` package with a small rule catalog:

```go
type Rule struct {
	ID          string
	Severity    Severity
	Description string
	Fixable     bool
	Run         func(*spec.Document, *Reporter)
}
```

Rules emit findings through a `Reporter`, which centralizes message construction,
location attachment, summary counts, deterministic sorting, and blocked-rule
behavior. Rules should be narrow and named one-to-one with the durable lint
sub-spec. Fixable rules also register a repair with the reporter for each
finding whose fix can be applied deterministically.

Parsing and gross shape failures emit `invalid-frontmatter` and mark dependent
model rules blocked. Shape failures inside a parent node block only rules that
need that malformed shape; they do not suppress unrelated checks elsewhere in
the document.

Genuinely unrecognized keys and wrong YAML shapes are treated as model-shape
failures under `invalid-frontmatter` for this change; do not add a narrower
unknown-key rule until the durable lint sub-spec calls for one. Known root-only
keys are the exception: `title` or `ratingScale` on a non-root target is not an
unknown key but a named defect, and the model must route it to
`misplaced-root-key` — located at the offending key, non-blocking — rather than
fold it into `invalid-frontmatter`. Build shape detection from the `yaml.Node`
model, not from struct decoding with known-field enforcement, so that a condition
with its own rule stays a locatable finding instead of a blocking decode error.

### Wire the CLI thinly

Add `newLintCmd()` under `internal/cli` and register it from `root.go`.

For this change, use a minimal invocation shape: `qualitymd lint [path]`,
defaulting to `QUALITY.md`. Stdin is deferred: `init` already uses `-` as a
stdout sentinel, but the shared file/stdin convention is parent-CLI work, so
`lint` does not add `-` handling yet. To avoid a confusing failure — a literal
attempt to open a file named `-` — `lint` rejects a bare `-` with a clear error
this phase, reserving the token for the parent [CLI spec](../../../specs/cli.md).

The command owns only argument parsing, `--json`, `--fix`, output routing, and
exit status. Parsing, linting, repair, sorting, and JSON document construction
live below `internal/cli`.

### Output structs are the JSON contract

Represent the lint result with exported structs in `internal/lint`. The repair
records in this result are data, separate from the executable `RepairOp` values
used internally by the writer:

```go
type Result struct {
	SchemaVersion int       `json:"schemaVersion"`
	Path          string    `json:"path"`
	Valid         bool      `json:"valid"`
	Summary       Summary   `json:"summary"`
	Findings      []Finding `json:"findings"`
	Repairs       []RepairRecord `json:"repairs"`
	NextActions   []Action  `json:"nextActions"`
}
```

The CLI JSON path should marshal this value directly. Human output should render
the same `Result`, not rerun lint or use a separate finding shape.

## Alternatives

- **Extend `spec.Load` validation in place.** Rejected. A fail-fast validator
  cannot produce the multi-finding, warning-aware, location-rich report required
  by the lint spec. It also would not produce the traversal primitives future
  query commands need.
- **Build `lint` directly over raw YAML with no shared model.** Rejected. That
  would ship faster for this one command, but the next query or evaluation
  feature would need to rebuild target/factor/requirement traversal from
  scratch, creating drift in factor resolution and target tree semantics.
- **Create a broad query engine now.** Rejected. The need for traversal is real,
  but the user-facing query surface is not specified yet. Build the graph and
  walkers now; design selectors and output formats when a query command is
  proposed.
- **Create a general model writer now.** Rejected. `lint --fix` needs a narrow
  repair writer, not a broad mutation API for arbitrary future commands. The
  writer should apply rule-owned repairs against stable model paths; future write
  commands can generalize from that once their contracts exist.
- **Emit patches instead of writing in place.** Rejected for this phase. Patch
  output is useful for review workflows, but it adds another output contract and
  does not replace the common "clean this file" use case. The durable lint spec
  scopes in-place writes and defers patch/full-file repair modes.

## Trade-offs & risks

- **More upfront structure than a simple validator.** Building a document model
  with node references costs more than direct struct validation, but it pays for
  itself immediately in locations, multiple findings, traversal-dependent rules,
  and repair targeting. Keep the first pass focused on the nodes the current
  rule set needs; do not model evaluation results or future report data.
- **`yaml.Node` can leak implementation concerns.** Keep YAML-specific details at
  the parser boundary. Rule code should mostly work with model nodes,
  `modelPath`, and source positions, not raw YAML traversal. The exception is
  repair code, which may touch YAML nodes through a narrow `RepairOp` interface.
- **Ordering has to be pinned.** YAML map order is input order, not sorted order.
  The model should preserve source order for source-position ordering and use
  explicit tie-breakers so output remains deterministic. Add tests where input
  map order and rule id tie-breakers would otherwise produce unstable output.
- **YAML round-tripping can still cause churn.** `yaml.Node` preserves more than
  struct marshalling, but not every byte of original frontmatter formatting.
  Tests should pin the important promise: body bytes are unchanged, unrelated
  keys stay in order, and repairs do not rewrite more of the frontmatter than
  needed.
- **Atomic replacement has platform edges.** Writing a temp file and renaming it
  is the right default, but permissions, symlinks, and cross-device paths need
  careful handling. Keep temp files in the target directory, preserve the
  original file mode where possible, and refuse `--fix` when the target path is a
  symlink until symlink semantics are specified.
- **`Load`'s acceptance tightens.** Today's loader treats a fence-less file as
  YAML in its entirety; `Parse` requires a frontmatter block, so fence-less input
  now reports `invalid-frontmatter`. This is the intended direction — the lint
  spec requires a leading frontmatter block — but it changes the existing load
  path. The scaffold conformance test that runs the embedded skeleton through the
  loader must stay clean under both the stricter parser and the full error-rule
  set (it must not, for example, trip `empty-model`).
- **The input convention is still not a durable CLI-wide contract.** This design
  uses `lint [path]` with a default of `QUALITY.md` so implementation can move.
  The parent CLI spec still needs to settle the shared file/stdin convention, and
  `lint` may need a small follow-up if that contract chooses different stdin
  semantics.

## Open questions

- **Parent CLI invocation contract** _(open)._ This change implements
  `lint [path]`, defaulting to `QUALITY.md`, as a provisional shape and leaves
  stdin and the shared file-argument convention to the parent
  [CLI spec](../../../specs/cli.md), where they remain "To be specified." `lint` may
  need a small follow-up if that contract settles different stdin semantics. The
  shape is deliberately not recorded in a durable spec by this change (see the
  [change's affected-specs note](../0003-implement-lint-command.md)).

Settled during design:

- **Package name.** Keep the shared document/model code in `internal/spec` for
  this change. Revisit `internal/model` only when a second command needs the same
  APIs and the package boundary becomes misleading.
- **Package boundary.** `internal/spec` is the document layer and does not import
  `internal/lint`. `internal/lint` imports `spec` and owns the rule catalog plus
  the valid-model convenience (`lint.Load`). The dependency is one-way, so there
  is a single rule implementation and no import cycle.
- **Repair ownership.** `internal/lint` owns rule-level repair operations;
  `internal/spec` owns rendering and atomic file replacement.
- **Unknown keys.** Report genuinely unrecognized keys and wrong YAML shapes as
  `invalid-frontmatter` in this phase. Known root-only keys (`title`,
  `ratingScale`) on a non-root target are not unknown keys: they are reported as
  `misplaced-root-key`, located at the offending key, not folded into
  `invalid-frontmatter`.
