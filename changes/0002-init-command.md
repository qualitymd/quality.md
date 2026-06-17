---
type: Change
title: Specify the init command
description: Settle what `qualitymd init` scaffolds, where it writes, and how it refuses to clobber existing files.
status: Design
tags: [cli, init]
timestamp: 2026-06-17T00:00:00Z
---

# Specify the init command

[`qualitymd init`](../specs/cli/init.md) is specified today only as a stub with a
"To be specified" list. This change settles that list: what the command
scaffolds, where it writes, and how it behaves when a file is already there. The
detail lives in its children:

- [Functional spec](0002-init-command/spec.md) — what `init` must do.
- [Design doc](0002-init-command/design.md) — how it's built, and why.

## Motivation

`init` is the first command an author runs, and the on-ramp to the whole format:
a good scaffold teaches the target → factor → requirement shape and the rating
scale by example, so the author edits rather than authors from scratch. Right now
nothing says what that scaffold contains or how the command behaves, so it can't
be built. Settling the *what* unblocks the design and implementation that follow.

## Scope

Covered: the scaffold's contents (a seeded rating scale, a commented
target → factor → requirement skeleton, and the recommended body sections as
headed stubs), the output target and stdout (`-`) piping, and refusing to
overwrite an existing file behind an explicit `--force`.

Deferred: cross-cutting CLI behavior (global flags, exit-code semantics, quiet
verbosity) belongs to the [CLI spec](../specs/cli.md), not here; interactive
prompting and selectable scaffold templates are left for later (see the spec's
[Deferred](0002-init-command/spec.md#deferred) list).

## Affected specs & docs

Updated before this change reaches **Done**:

- [ ] [`specs/cli/init.md`](../specs/cli/init.md) — replace the stub and its "To
      be specified" list with the durable `init` sub-spec, absorbing this
      change's functional spec.
- [ ] [`README.md`](../README.md) — drop the *(planned)* marker on `init` once
      the command ships.

## Status

`Design`. See the [status lifecycle](index.md#status-lifecycle). Implementation is
this change's **In-Progress** phase, taken up once the design is settled. When this
change reaches **Done**, move it and its [child folder](0002-init-command/) into
[`archive/`](archive/).
