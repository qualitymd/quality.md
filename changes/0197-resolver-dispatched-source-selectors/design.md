---
type: Design Doc
title: Resolver-dispatched source selectors — design
description: How selector kinds are detected, prose resolution rides the harness checkpoint transport as first-class work units, and captured bundles gain provenance of record.
---

# Resolver-dispatched source selectors — design

## Context

Answers the [functional spec](spec.md) (R1–R6) for the
[Resolver-dispatched source selectors](../0197-resolver-dispatched-source-selectors.md)
case, under its settled questions: the format commits to non-filesystem
selectors (Q2), and kind is detected from the bare selector string (Q1). The
governing idea, from 0196's
[considerations sketch](../archive/0196-spec-faithful-model-reading/considerations.md):
resolution (gathering, per kind) is separated from judgment (rating against a
captured bundle), with the bounded, hashed evidence bundle as the contract
between them. The design's job is to add the prose kind and the harness
resolution path **around** the existing deterministic packaging without
touching it, reusing the checkpoint transport, unit state, and acceptance
machinery 0194 built rather than growing a second dispatch path.

## Approach

### Kind detection, pinned at run creation (R1)

Detection is runner-owned, in `internal/runner/source.go`, and reuses the
classification the packager already performs implicitly:

```
detectSourceKind(workspaceRoot, selector):
  if contains glob metacharacters (*, ?, [)  → glob      # existing rule
  if absolute, "..", or escaping             → path      # containment keeps it unresolvable
  if os.Stat(workspaceRoot/selector) exists  → path
  else                                       → prose
```

The escape check precedes the stat so `../shared` stays a filesystem selector
that fails loudly under 0196's workspace-containment rule — it never falls
back to prose. `model.EffectiveSource` is unchanged: which selector applies is
model semantics; what kind it is depends on filesystem state, which is runner
territory.

Detection runs once, at run creation, and the result is persisted in the run
artifact's new per-area `sources` record (below). Resume reads the pinned kind
instead of re-detecting, so a file appearing or vanishing mid-run cannot
silently re-dispatch a selector to a different resolver or change the graph
shape; switching kinds is only possible by starting a new run.

### Resolution as first-class work units (R2, R3)

For each in-scope area whose pinned kind is prose, `BuildGraph` emits one
`resolveSource:<areaRef>` unit — evaluator-backed, emitted alongside the
area's frame unit — and every source-consuming unit in the area (the
`assessRateRequirement` units) adds it as a dependency. Analysis and advice
units consume prior results, not source, and are unwired from it.

Making resolution a unit is the load-bearing choice: it inherits the entire
harness checkpoint machinery for free — `PendingEvaluatorCall` correlation,
deterministic request IDs, the input-hash guard, attempt counting, retry
classification, resume, and call logging all key off unit state
(`internal/runner/harness.go`), so no second pending-request mechanism is
needed. Harness-backed runs already resolve to `concurrency: 1`, so the
single-pending-checkpoint constraint holds unchanged.

`buildWorkRequest` gains the `resolveSource` kind: instructions to gather the
material the selector describes and return it as files; context carrying the
selector, its kind, and the area frame; an expected schema of
`{files: [{path, content}], minItems: 1}`; and an empty `Source` field — the
resolver is fed a description, never pre-gathered evidence. The instructions
direct the harness to return a classified `source_unavailable` failure when
the material the selector describes does not exist — including when the
selector reads like a filesystem path that names nothing (the Q1 typo case) —
rather than improvising evidence.

Acceptance validates the returned envelope through the shared paths, then
**captures**: non-empty unique file paths, the existing 64 KB/file and
512 KB/bundle caps with truncation marks, SHA-256 per file, and the same
bundle hash function as walked source. The captured bundle is persisted into
the run artifact before the unit completes, so every dependent judgment
request is built from persisted data. `areaSourceBundle` then dispatches on
the pinned kind: path/glob packages through the unchanged deterministic walk;
prose reads the captured bundle (its resolution unit is a completed
dependency, so absence is an internal error, never a silent empty bundle).

Because a judgment unit's `SourcePackageHash` is the captured bundle's hash,
the runner's guarantees hold by construction: the input-hash guard rejects a
result attached to different evidence, and re-judgment is bound to exactly the
captured material — reproducibility of record, independent of how gathering
went.

### Captured bundles and provenance live in `evaluation.json` (R5)

`evaluation.json` gains a per-area `sources` record, written at run creation
with `{selector, kind}` and completed as bundles materialize:

```json
"sources": {
  "area:api": {
    "selector": "open tickets in the support queue",
    "kind": "prose",
    "resolver": "harness",
    "harnessRuntime": "claude-code",
    "bundleHash": "…",
    "capturedAt": "2026-07-11T00:00:00Z",
    "files": [{"path": "…", "sha256": "…", "truncated": false, "content": "…"}]
  }
}
```

Deterministic areas record the same provenance shape (`resolver: "walk"`)
without `content` — their material is re-readable from the workspace, and
resume re-packages it as today. Harness-resolved areas keep `content`, because
the captured bundle **is** the evidence of record: resume must rebuild
dependent requests from it, not re-gather. This one record serves three
masters: kind pinning (R1), resume for prose areas (R2/R3), and audit
provenance (R5) — a reviewer reads the same shape for walked and
agent-gathered evidence and can tell which is which.

Prose file paths are labels for gathered material (a ticket ID, a URL, a
repo-relative path), not workspace paths; they are recorded and hashed
verbatim.

### Unsupported selectors fail loudly, early, and distinctly (R4)

`evaluator.Capabilities` gains `SourceResolution bool`; only the harness
evaluator declares it. A new failure category `selector_unsupported` is added
to the taxonomy, tripped at run creation: if any in-scope area's pinned kind
has no resolver under the selected evaluator — prose without
`SourceResolution` — the run fails before any judgment is dispatched, naming
the selector, its detected kind, and the remedy (evaluate through harness
dispatch, or change the selector). Dry-run surfaces each area's detected kind
and its resolution units, so the dispatch plan is visible before anything
runs. A per-unit guard backstops the plan-time check.

The classification boundary: `selector_unsupported` means _this run cannot
resolve this kind of selector_; `source_unavailable` keeps meaning _the
material is missing_ — whether a path matched nothing or a harness resolver
reported the described material does not exist.

### Source-as-data across resolvers (R6)

Captured files flow into `WorkRequest.Source` exactly like walked files, so
the prompt rendering and standing safety instructions of the source-as-data
invariant apply without a special case. Nothing to build; the design keeps it
true by routing everything through one bundle type.

## Spec response

- **R1** — detection is a thin classifier in front of the untouched 0196
  packaging path; path/glob behavior is preserved by not moving a line of it.
  Pinning at creation satisfies the resume-stability clause.
- **R2** — judgment units read only `areaSourceBundle`, which serves either a
  deterministic package or a captured bundle; there is no third path by which
  evidence reaches a judge.
- **R3** — resolution rides the existing checkpoint transport as a unit kind,
  and capture-before-dependents is enforced by graph dependency, not
  convention.
- **R4** — the capability plus plan-time check makes the failure precede any
  judgment work, with the remedy in the message.
- **R5** — the `sources` record is the single provenance surface, uniform
  across resolvers.
- **R6** — by construction, one bundle type through one rendering path.

## Alternatives

- **Typed or hybrid `{kind, selector}` frontmatter (Q1).** Unambiguous — a
  typo'd path could never become a prose selector — but it changes the
  frontmatter shape and `quality.schema.json`, taxes hand-authoring of the
  overwhelmingly common path/glob case, and (in the hybrid form) creates two
  spellings for the same filesystem selector. Rejected with Q1: detection
  order makes the filesystem interpretation always win, and the typo hazard is
  mitigated rather than structural (see trade-offs).
- **Lazy resolution inside `areaSourceBundle`, no graph unit.** Resolving on
  first use looks smaller, but it would need its own pending-request state,
  retry accounting, and resume correlation outside unit state — a second,
  weaker copy of the harness machinery — and resolution work would be
  invisible to dry-run, status, and the work-unit ledger. Rejected.
- **Captured bundles as run-local sidecar files** (`<run>/source/…`), with
  `evaluation.json` holding only hashes. Keeps the artifact lean, but splits
  authoritative resume state across files, forfeits the store's atomic
  temp-file-plus-rename write, and adds a data-layout contract for what the
  512 KB cap already bounds. Rejected: `evaluation.json` remains the one
  authoritative artifact.
- **Reusing the `Subagents` capability as the resolver seam** (the 0196
  sketch's lean). `Subagents` describes an evaluator's internal parallelism
  for judgment work; serving a distinct request kind is a different promise,
  and conflating them would make every subagent-capable evaluator implicitly
  claim resolution. A dedicated `SourceResolution` capability keeps the
  declaration honest.
- **Re-detecting kind on resume instead of pinning.** Fresher, but a file
  created mid-run (often by the evaluating agent itself) would flip a prose
  area to path, change the graph shape under a pending checkpoint, and strand
  the captured bundle. Pinning keeps resume coherent; a new run picks up the
  new state.

## Trade-offs & risks

- **A mistyped path dispatches an agent** (accepted with Q1). `docs/guids`
  stats to nothing and detects as prose. Mitigations: dry-run and the run
  receipt name every prose detection before judgment; the resolution
  instructions direct the harness to fail path-like selectors that name
  nothing as `source_unavailable` rather than gathering by guesswork; and the
  provenance record makes what was gathered, and by whom, auditable after the
  fact.
- **Gathering is not reproducible.** Two runs over the same model may capture
  different evidence for the same prose selector. This is the case's accepted
  premise — reproducibility of record, not of gathering — but it means run
  comparisons for prose areas diff captured bundles, not selectors.
- **Captured bundles are frozen.** Workspace writes do not invalidate a
  prose area's evidence (no re-stat), which is exactly the churn relief the
  0196 re-validation asked for — but it is also a staleness risk: re-gathering
  requires a new run, and readers of the report must trust `capturedAt`-era
  material.
- **`evaluation.json` grows raw gathered content** — up to 512 KB per prose
  area. The artifact is workspace-local and may be committed; gathered
  material (tickets, live-system output) lands in it verbatim. The cap
  bounds the size; the sensitivity is inherent to capturing evidence of
  record and is the same class of exposure as quoted evidence in reports.
- **Kind detection depends on filesystem state at creation.** An author who
  writes a prose selector and _later_ adds a matching path gets prose behavior
  until a new run. The `sources` record and dry-run make the pinned kind
  visible, which is the honest remedy.

## Open questions

None. The gating format questions were settled before this design; resolution
instruction wording and the exact `sources` field names are implementation
detail within the shapes above.
