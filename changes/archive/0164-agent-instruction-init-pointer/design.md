---
type: Design Doc
title: Agent Instruction Init Pointer — design
description: Implementation approach for init-time agent instruction pointers.
tags: [cli, init, agents, setup]
timestamp: 2026-06-29T00:00:00Z
---

# Agent Instruction Init Pointer — design

Companion to the
[Agent Instruction Init Pointer functional spec](spec.md).

## Context

The change adds a small second write to the normal `qualitymd init` path: after
the model scaffold is written, the CLI makes `QUALITY.md` visible from local
agent instruction files. The stdout form remains a pure scaffold stream, and
`/quality setup` opts out because setup owns its own mutation boundary.

## Approach

Add a focused `internal/agentinstructions` package that owns instruction-file
target detection, pointer rendering, duplicate detection, and file writes. The
CLI package should consume it through a small options/result shape rather than
embedding Markdown update rules in `init.go`.

The updater runs only after the scaffold file write succeeds. It works from the
process current directory, not from parent repository discovery:

1. Build target candidates in this order: `AGENTS.md` with `Create: true`, then
   `CLAUDE.md` and `GEMINI.md` with `Create: false`.
2. Resolve existing paths through `filepath.EvalSymlinks` plus file identity from
   `os.Stat`; use absolute cleaned paths for missing targets. This de-dupes
   symlinked files such as this repo's `CLAUDE.md` and `GEMINI.md` links.
3. For each unique writable target, render the two-line block with a relative
   link from the target file's directory to the initialized model path.
4. Skip a file when it already contains the marker comment or an existing
   `QUALITY.md` quality-model pointer from this feature.
5. Append the block with a separating blank line, preserving existing content.

`internal/cli/init.go` gains `--no-agent-instructions` and a result field in the
JSON receipt. Human output reports instruction files only when the updater
created or changed a file. The stdout `-` branch returns before invoking the
updater, preserving current pipeline semantics.

The scaffold note belongs in the embedded `skeleton.md` and
`skeleton-minimal.md` bodies, not in CLI output, because users may inspect the
file later after the terminal output is gone. The note is intentionally brief and
points to `/quality setup` without teaching the whole setup workflow.

Setup runtime and durable specs switch their missing-file scaffold command to
`qualitymd init --no-agent-instructions [path]`, then retain the existing
read-before-authoring step. Existing-file handling becomes a maturity branch:
scaffold-only files can be replaced after the normal review gate, partially
authored files are updated conservatively, and mature files route toward review
or improve unless the user asked for setup-style reshaping.

## Spec response

The separate package keeps pointer detection and idempotence testable without
driving the Cobra command. Running it after a successful scaffold write prevents
agent instruction edits when `QUALITY.md` creation fails. Returning structured
result data lets human and JSON output share the same fact source.

Relative-link rendering satisfies custom path support while preserving the
visible `QUALITY.md` label. Marker-based duplicate detection covers normal
reruns; content-based detection for the pointer sentence catches same-feature
blocks even if spacing changes. Physical-file de-duping covers symlinked
instruction files.

## Alternatives

**Only print a next action.** Rejected because it preserves the current
remember-a-second-edit problem. The point of this change is to make the scaffold
agent-visible by default.

**Create every supported agent-specific file.** Rejected because it adds visible
project clutter. Creating `AGENTS.md` gives a neutral default; other agent files
are updated only when the project already owns them.

**Put the update logic in `internal/scaffold`.** Rejected because scaffold owns
the `QUALITY.md` artifact. Agent instruction files are a neighboring project
bootstrap concern with separate targets, idempotence, and reporting.

**Use begin/end managed markers.** Rejected because the block is deliberately
two lines and not managed after insertion. A single origin comment is enough.

## Trade-offs & risks

Creating `AGENTS.md` by default is a new side effect for `init`; the opt-out flag,
stdout exception, and terse reporting keep it explicit enough for automation.

The duplicate detector intentionally does not try to understand every possible
pre-existing human-written QUALITY.md mention. A hand-written sentence may remain
alongside the inserted pointer. That is preferable to accidentally suppressing
the origin-marked block because of an unrelated mention.

Relative links from agent instruction files outside the model directory can look
less tidy, but they keep the pointer correct for custom init paths.

## Open questions

None.
