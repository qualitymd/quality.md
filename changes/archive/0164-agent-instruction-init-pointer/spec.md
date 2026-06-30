---
type: Functional Specification
title: Agent Instruction Init Pointer — functional spec
description: Delta contract for adding concise agent instruction pointers during qualitymd init.
tags: [cli, init, agents, setup]
timestamp: 2026-06-29T00:00:00Z
---

# Agent Instruction Init Pointer — functional spec

Companion to the
[Agent Instruction Init Pointer](../0164-agent-instruction-init-pointer.md).
This spec states _what_ the case must do; the [design doc](design.md) covers
_how_.

The key words **MUST**, **MUST NOT**, **SHOULD**, and **MAY** are to be
interpreted as described in BCP 14 when, and only when, they appear in all
capitals.

## Background / Motivation

QUALITY.md is most useful when coding agents and AI assistants notice it during
ordinary project work. `qualitymd init` currently creates the model file but
leaves agent instruction files unaware of it, so the user must remember a second
manual edit before the scaffold participates in the agent-first experience. The
init-time pointer closes that discoverability gap with one concise line while
preserving `/quality setup` as the guided authoring workflow.

## Scope

Covered:

- default agent instruction pointer creation/update during file-writing
  `qualitymd init`;
- opt-out through `--no-agent-instructions`;
- idempotence across repeated runs and symlinked instruction files;
- JSON and human reporting of agent instruction effects;
- scaffold body guidance that routes standalone `init` users to `/quality setup`;
  and
- setup workflow guidance that uses the opt-out flag and handles existing model
  maturity.

Deferred:

- configurable agent instruction file names;
- detection beyond the working directory that receives the model file;
- removal, migration, or replacement of older pointer wording;
- any agent-specific setup beyond a single QUALITY.md pointer; and
- interactive prompting.

## Requirements

1. File-writing `qualitymd init` **MUST** add this exact Markdown block to
   eligible agent instruction files by default:

   ```text
   <!-- Added by qualitymd init. -->
   See [QUALITY.md](QUALITY.md) for this project's quality model.
   ```

   > Rationale: The pointer is intentionally small: it names origin and makes the
   > model discoverable without turning `init` into a workflow or duplicating
   > setup guidance.
   >
   > Durable spec: modify `specs/cli/init.md` to define the default pointer block
   > and its target files.

2. File-writing `qualitymd init` **MUST** ensure an `AGENTS.md` file in the
   current working directory contains the pointer block, creating `AGENTS.md`
   when it does not exist.

   > Durable spec: modify `specs/cli/init.md` to define the default `AGENTS.md`
   > target.

3. File-writing `qualitymd init` **MUST** update detected existing `CLAUDE.md`
   and `GEMINI.md` files in the current working directory with the same pointer
   block.

   > Rationale: Existing project instruction files for popular agents should
   > learn about `QUALITY.md`, but `init` should not create every agent-specific
   > file proactively.
   >
   > Durable spec: modify `specs/cli/init.md` to define detected-file targets.

4. `qualitymd init --no-agent-instructions` **MUST** write the selected scaffold
   without creating or updating agent instruction files.

   > Rationale: Setup and automation need the scaffold side effect without
   > widening their own mutation surface.
   >
   > Durable spec: modify `specs/cli/init.md` to add the flag.

5. `qualitymd init -` **MUST NOT** create or update agent instruction files.

   > Rationale: The stdout form is a pure artifact stream suitable for shell
   > redirection; side effects would make it unsafe in pipelines.
   >
   > Durable spec: modify `specs/cli/init.md` to keep stdout behavior pure.

6. Agent instruction pointer insertion **MUST** be idempotent: repeated `init`
   runs, including `--force`, must not add duplicate QUALITY.md pointers to the
   same instruction file.

   > Durable spec: modify `specs/cli/init.md` to define duplicate prevention.

7. When multiple targeted instruction paths resolve to the same physical file,
   `qualitymd init` **MUST** write that file at most once.

   > Rationale: This repository symlinks `CLAUDE.md` and `GEMINI.md` to
   > `AGENTS.md`; symlinked instruction files must not receive repeated blocks.
   >
   > Durable spec: modify `specs/cli/init.md` to define symlink/shared-target
   > handling.

8. When the initialized model path is not `QUALITY.md` beside the instruction
   file, `qualitymd init` **MUST** render the pointer link relative from the
   instruction file to the initialized model path while keeping the visible link
   label `QUALITY.md`.

   > Durable spec: modify `specs/cli/init.md` to define link rendering.

9. Human `qualitymd init` output **SHOULD** report which agent instruction files
   were created or updated, and **SHOULD** omit that line when no instruction file
   changed.

   > Durable spec: modify `specs/cli/init.md` to define human reporting.

10. `qualitymd init --json` **MUST** include agent instruction effects in its
    success receipt.

    > Durable spec: modify `specs/cli/init.md` to extend the JSON receipt.

11. Both guided and minimal `qualitymd init` scaffolds **MUST** include a brief
    body note telling users to invoke `/quality setup` when the file was created
    outside that workflow.

    > Rationale: Direct CLI initialization remains available, but the agent skill
    > is the intended path for turning a scaffold into a project-specific model.
    >
    > Durable spec: modify `specs/cli/init.md` to add scaffold content guidance.

12. The `/quality setup` workflow **MUST** scaffold missing model files with
    `qualitymd init --no-agent-instructions`.

    > Rationale: Setup already has a defined mutation boundary and should not
    > surprise users by editing agent instruction files while preparing to author
    > `QUALITY.md`.
    >
    > Durable spec: modify `specs/skills/quality-skill/workflows/setup.md`.

13. The `/quality setup` workflow **MUST** classify an existing `QUALITY.md` as
    scaffold-only, partially authored, or mature before planning edits, then
    preserve useful existing content unless the review gate names replacement.

    > Rationale: Setup is useful after a direct `init` run and after early manual
    > authoring; treating every existing file as the same update risks discarding
    > work or under-explaining why a rewrite is appropriate.
    >
    > Durable spec: modify `specs/skills/quality-skill/workflows/setup.md`.

## Durable spec changes

### To add

None

### To modify

- `specs/cli/init.md` - define the default agent instruction pointer, target
  detection, `--no-agent-instructions`, idempotence, symlink handling, relative
  links, reporting, JSON receipt, and scaffold setup note.
- `specs/skills/quality-skill/workflows/setup.md` - require setup to scaffold
  with `--no-agent-instructions` and classify existing model maturity before
  edits.
- `specs/skills/quality-skill/quality-skill.md` - update only if parent workflow
  summary needs to name the changed setup boundary.

### To rename

None

### To delete

None
