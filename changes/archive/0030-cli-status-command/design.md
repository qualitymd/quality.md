---
type: Design Doc
title: CLI status snapshot command - design doc
description: How qualitymd status assembles a deterministic project-state snapshot from existing CLI-owned mechanics.
tags: [cli, command, wizard, design]
timestamp: 2026-06-19T00:00:00Z
---

# CLI status snapshot command - design doc

Design behind the
[CLI status snapshot command](../0030-cli-status-command.md) change and its
[functional spec](spec.md). The spec fixes _what_ the snapshot must report; this
doc covers the intended implementation approach.

## Context

Most of the data `qualitymd status` needs already exists behind deterministic
code paths:

- `internal/lint` can validate a model and emit stable finding data.
- `internal/model` can decode a valid model for tree traversal.
- `internal/evaluation` can resolve repository-local evaluation directories,
  load runs, count records, and compute reportability gaps.
- recommendation superseding is already represented in runtime record
  frontmatter.

The command should compose those mechanics rather than re-parse artifacts in the
CLI layer or in the `/quality` skill.

## Approach

Add an `internal/status` package that owns snapshot assembly and exposes one
entry point:

```go
func Snapshot(opts Options) (*SnapshotResult, error)
```

`internal/cli/status.go` should stay thin: parse `status [path]`, map errors to
the shared exit categories, render human output, and write the JSON result when
`--json` is passed.

### Snapshot flow

`Snapshot` should assemble the result in this order:

1. Resolve the selected model path, defaulting to `QUALITY.md`.
2. If the model path is absent, return a successful snapshot with
   `readiness: "missing-model"` and no model-shape or evaluation-history detail
   that depends on model bytes.
3. Run `lint.Check(path)` and copy the lint summary and findings into the
   snapshot.
4. If lint errors are present, return `readiness: "invalid-model"` without
   decoding model shape.
5. Decode the valid model through `lint.Load(path)` or the equivalent parsed
   document path, then traverse it once for counts and source coverage.
6. Resolve the repository root from the selected model path and locate the
   evaluation directory with the same `.quality/config.yaml` fallback used by
   evaluation commands.
7. List recognized run folders in deterministic order.
8. For each run, compare `model.md` bytes to the selected model bytes for
   staleness, then load the run and compute record counts, reportability gaps,
   and active recommendation counts.
9. Derive readiness and next actions from the assembled mechanical state.

The status package should keep the JSON structs explicit rather than reusing
`lint.Result` or `evaluation.Status` wholesale as top-level objects. Reusing
their nested public fields is fine; the status result is its own stable API.

### Model traversal

Traverse the valid model recursively with a small accumulator:

- count the root Model as the root Target;
- count each child entry under `targets` as a Target;
- count every Factor and nested sub-factor;
- count every Requirement under both Target and Factor scopes;
- carry the currently resolved Source down the Target tree; and
- emit one source-coverage row per Target in traversal order.

Traversal order should follow the deterministic ordering used elsewhere in the
tool: sort map keys lexicographically before walking children. This makes JSON
arrays and human summaries stable.

### Evaluation helpers

Keep run parsing and reportability rules in `internal/evaluation`. Add small
exported helpers only where they preserve that boundary, for example:

```go
func EvaluationDir(repoRoot, override string) (abs string, rel string, error)
func ListRunDirs(evalDirAbs string) ([]RunDir, error)
func (r *Run) ActiveRecommendationCount() int
```

`internal/status` can then compose these helpers without duplicating run-name
regular expressions, config parsing, or recommendation-superseding rules.

Malformed individual runs should produce a run summary with an inspection
problem. The command should continue with later runs so a single corrupt run does
not hide the rest of the history.

### Active recommendations

Active recommendation counting should use the same identity rules as
recommendation superseding elsewhere: a reference can name the file path, the
record id, or the conventional `recommendations/<id>.md` path. Process
recommendation records in deterministic file order and mark referenced earlier
records as superseded when the reference resolves.

Do not inspect `report.md`. If a future status field needs generated report
state, prefer `report.json` because it is a machine artifact; do not parse
human Markdown.

### Human output

Human output should be compact and route-oriented:

```text
Status
- QUALITY.md: present, valid
- Model: 3 targets, 5 factors, 12 requirements
- Evaluation history: 2 runs, 1 incomplete, 1 stale, 2 active recommendations
- Readiness: needs evaluation reconciliation

Next
- qualitymd evaluation show-status quality/evaluations/0002-subject-quality-eval
```

Detailed finding rows, full source coverage, and every run summary belong in
`--json` or in more specific commands.

### Tests

Add focused tests in `internal/status` for:

- missing model snapshots exit through success data;
- invalid model snapshots carry lint findings but no shape;
- valid model traversal counts root targets, child targets, factors,
  requirements, and source inheritance;
- absent evaluation directory reports zero runs;
- recognized runs are ordered by run number and name;
- stale detection compares run `model.md` bytes to the selected model bytes;
- malformed runs are reported without stopping later runs; and
- active recommendations ignore superseded records.

Add CLI tests in `internal/cli` for:

- `status --json` emits a snapshot on stdout;
- human `status` stays compact;
- `status -` exits as usage error; and
- snapshot states that report project problems still exit `0`.

## Alternatives

**Teach the wizard to keep probing with shell commands.** Rejected. That keeps
mechanical parsing in the skill and repeats the failure this change is meant to
fix.

**Make `status` shell out to existing `qualitymd` commands.** Rejected. It would
add process overhead, complicate error mapping, and make the command depend on
its own installed binary instead of shared internal packages.

**Embed full lint and show-status JSON unchanged.** Rejected as a top-level
shape. It would make the status schema an accidental nesting of other command
schemas. The status command should own its snapshot shape while reusing the same
public nested fields where that helps consumers.

**Read generated report Markdown for ratings and advice.** Rejected. Markdown is
the human artifact; status should read runtime records or machine artifacts and
never scrape report bodies.

## Trade-offs and risks

The main trade-off is a new package that composes several existing packages. That
is preferable to putting snapshot assembly into `internal/cli`, because the
snapshot needs unit tests independent of terminal rendering.

The biggest risk is schema creep: status could become a full evaluation report.
Keep the human output compact, keep JSON focused on routing signals, and point
to existing detailed commands for audit data.

Another risk is duplicating recommendation-active-state logic. The design avoids
that by placing active-count helpers in `internal/evaluation`, beside the run
loader and reportability rules.

## Open questions

- Whether `status --json` should include every run summary or cap historical
  runs in a later phase. Start with every recognized run because the JSON shape
  is agent-facing and deterministic.
- Whether status should surface latest root rating from `report.json`. The
  functional spec leaves this optional; defer until a consumer needs it.
