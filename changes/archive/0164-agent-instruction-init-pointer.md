---
type: Change Case
title: Agent Instruction Init Pointer
description: Make qualitymd init surface the new QUALITY.md to agent instruction files while preserving setup's narrow mutation boundary.
status: Done
tags: [cli, init, agents, setup]
timestamp: 2026-06-29T00:00:00Z
---

# Agent Instruction Init Pointer

A **Change Case** capturing the *why* and *status*; the detail lives in its
children:

- [Functional spec](0164-agent-instruction-init-pointer/spec.md) - what the case must do.
- [Design doc](0164-agent-instruction-init-pointer/design.md) - how it is built,
  and why.

## Motivation

`qualitymd init` can create a valid `QUALITY.md`, but the agent-facing harness
may still miss it. Since QUALITY.md is normally consumed through agents and
skills, a newly scaffolded model should be discoverable from the project's agent
instruction files without asking the user to remember a second manual edit.

The pointer must stay small and unsurprising. `init` is still a deterministic
scaffold command, not a setup workflow; the injected text should only make
`QUALITY.md` visible to agents, while `/quality setup` remains the guided path
for authoring a useful first model.

## Scope

Covered:

- Add a default `qualitymd init` side effect that creates or updates concise
  agent instruction pointers.
- Add `--no-agent-instructions` to disable that side effect.
- Keep stdout scaffold mode pure.
- Keep injection idempotent across repeated runs and symlinked instruction files.
- Update the starter scaffold to tell users to invoke `/quality setup` when init
  was run outside setup.
- Update `/quality setup` guidance so setup scaffolds through
  `qualitymd init --no-agent-instructions` and handles existing model maturity
  explicitly.

Deferred:

- Configurable agent instruction filename lists.
- Removing or migrating previously injected pointer blocks.
- Interactive prompts for agent instruction updates.
- Repository-wide parent directory discovery for agent instruction files.

## Affected artifacts

Derived by sweeping for `qualitymd init`, `--minimal`, setup scaffolding,
`AGENTS.md`, and agent instruction files.

**Code**

- [x] `internal/agentinstructions/` - add the instruction-file pointer updater.
- [x] `internal/cli/init.go` - add `--no-agent-instructions`, call the updater,
      report updates, and extend the JSON receipt.
- [x] `internal/cli/init_test.go` - cover default injection, opt-out, stdout,
      symlink/idempotence, force, custom paths, and JSON reporting.
- [x] `internal/scaffold/skeleton.md` and `internal/scaffold/skeleton-minimal.md`
      - add the `/quality setup` handoff note.
- [x] `internal/scaffold/scaffold_test.go` - assert the handoff note remains in
      scaffolds.

**Durable specs**

- [x] `specs/cli/init.md` - update the cumulative `init` contract.
- [x] `specs/skills/quality-skill/workflows/setup.md` - mirror setup's init
      opt-out and existing-file maturity handling.
- [x] `specs/skills/quality-skill/quality-skill.md` - no edit needed; the parent
      setup summary remains accurate.

**Bundled skill / runtime guidance**

- [x] `skills/quality/workflows/setup.md` - update runtime setup instructions.

**Durable docs / distribution**

- [x] `README.md` and `install.md` - mention the behavior only where the init
      path is documented.
- [x] `CHANGELOG.md` - note the user-visible CLI and setup behavior.
- [x] `changes/index.md`, `changes/archive/index.md`, and `changes/log.md` -
      Change Case lifecycle.

## Status

`Done`. Implemented and archived. `qualitymd init` now writes a concise,
idempotent pointer to local agent instruction files by default, supports
`--no-agent-instructions`, keeps stdout scaffold mode pure, reports pointer
effects in human and JSON output, points direct init users to `/quality setup`,
and setup uses the opt-out while handling existing model maturity explicitly.
