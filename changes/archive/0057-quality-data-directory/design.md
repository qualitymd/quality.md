---
type: Design Doc
title: Quality data directory - design doc
description: How workspace resolution, .quality/ defaults, root config, and strict lint rule options should be implemented.
tags: [workspace, config, evaluation, lint, skill]
timestamp: 2026-06-22T00:00:00Z
---

# Quality data directory - design doc

## Context

This design answers the
[Quality data directory functional spec](spec.md). The change has two related
jobs: give path-owning commands one shared workspace resolver, and let lint
accept the root `config` tooling key while keeping unknown-key diagnostics strict
by default.

The key vocabulary distinction is deliberate: a QUALITY.md workspace is the
resolved operating context for one selected `QUALITY.md`; `.quality/` is only
that workspace's default quality data directory.

## Approach

Add `internal/workspace` as the single path-resolution boundary for qualitymd
tooling paths. The package should not know evaluation-record semantics or lint
rule semantics; it should only resolve safe paths and parse the small config
surface needed by the CLI and skill.

The central type should carry both absolute and repository-relative path forms:

```go
type PathRef struct {
    Abs string
    Rel string
}

type Workspace struct {
    RepoRoot PathRef
    Model    PathRef

    Config        PathRef
    ConfigPresent bool

    DataDir     PathRef
    Evaluations PathRef
    Log         PathRef
}
```

Resolution should flow in one direction:

1. Start from the selected model path, defaulting to `QUALITY.md`.
2. Find the repository root from that model path.
3. Parse the model frontmatter only far enough to read root `config`.
4. Resolve `config`, or default it to `.quality/config.yaml`.
5. Read the config file when present.
6. Resolve `.quality/` as the default quality data directory.
7. Resolve `evaluationDir`, honoring command override, then config, then
   `.quality/evaluations`.
8. Resolve the log directory as `.quality/log`.

All path resolution should reuse one helper that rejects absolute paths and
repository escapes, then returns slash-normalized relative paths. Existing
evaluation path helpers can move into `internal/workspace` or delegate to it.

`internal/evaluation` should stop reading `.quality/config.yaml` directly.
`CreateRun`, `ResolveDir`, run listing, and command receipts should consume the
workspace's `Evaluations` path. `--evaluation-dir` remains an explicit override
passed into workspace resolution.

`internal/status` should resolve the same workspace for the selected model path
and use `Workspace.Evaluations` for history. This keeps status and evaluation
creation from drifting.

The runtime `/quality` skill does not need a separate implementation package,
but its instructions should mirror the same resolution order and terminology:
workspace, quality data directory, config file, evaluation directory, quality
log directory.

### Lint

Keep `config` out of the normative model structs and the companion JSON Schema.
It is a qualitymd tooling key, not a Model property.

Add a small internal rule-configuration shape to lint, scoped to the current
need:

```go
type Options struct {
    Rules RuleOptions
}

type RuleOptions struct {
    UnknownKey UnknownKeyOptions
}

type UnknownKeyOptions struct {
    AllowedRootKeys map[string]bool
}
```

The default lint options should allow root `config` and no other non-schema root
keys. Unknown nested keys continue to use the existing structural schema checks.

Add a dedicated `invalid-config` rule for bad root `config` values. Let the
unknown-key rule decide whether the key is allowed; let `invalid-config`
validate that an allowed `config` value is a non-empty scalar safe relative
path. This avoids treating "allowed key" as "arbitrary YAML is fine."

The public `qualitymd lint` command should keep using default options. No CLI
flag or config-file surface is added in this change.

## Alternatives

**Keep path resolution in `internal/evaluation`.** Rejected because status,
evaluation creation, and the skill already need the same answer. Leaving config
resolution in evaluation keeps the current drift point.

**Add `config` to the model schema.** Rejected because it would blur the
normative QUALITY.md Model with qualitymd-specific tooling behavior and would
also pressure `qualitymd schema` to describe a non-normative key.

**Create a general extension registry.** Rejected as too broad. The immediate
need is rule-configurable allowed keys, not a semantic extension framework.

**Expose lint rule configuration now.** Rejected because no user-facing syntax
has been designed. The implementation should have the seam, but the public
surface should wait for a concrete use case.

**Name `.quality/` the workspace.** Rejected because workspace usefully names
the whole resolved operating context for one model file. The directory inside
that context is narrower: the quality data directory.

## Trade-offs & Risks

Parsing root `config` before full lint means workspace resolution needs a small,
tolerant frontmatter read path. That path should only inspect one scalar key and
should not attempt to decode or validate the full Model.

The selected model path and repository root relationship should stay explicit.
Commands that currently assume the process working directory is the repo root
may need tests around nested invocation and `--model`.

There is a risk of two path helpers drifting if `internal/evaluation` keeps
some current helpers. Prefer moving shared helpers into `internal/workspace` and
leaving thin compatibility wrappers only where call sites would otherwise churn.

Keeping unknown keys as errors is intentionally stricter than the broad
extension allowance in the formal spec. Durable lint docs should describe this
as the default qualitymd lint profile, not as a universal statement that such a
document cannot be a valid extension-bearing QUALITY.md file.

## Open Questions

- Should `Workspace.DataDir` ever be configurable directly, or should only child
  paths such as `evaluationDir` be configurable?
- Should status JSON expose the resolved config path and quality data directory,
  or is the evaluation history path enough for this phase?
- Should the quality log eventually share the same resolver through a `logDir`
  config key, or remain fixed under `.quality/log/` until the skill has a
  queryable log command?
