---
type: How-to Guide
title: Designing Go packages
description: How to decide which package a type belongs in.
tags: [go, code, contributing]
timestamp: 2026-06-17T00:00:00Z
---

# Designing Go packages

The `qualitymd` CLI is Go under `cmd/qualitymd/` and `internal/`. This guide
covers one recurring decision: **which package a type belongs in.** It exists
because the easy answer — "wherever it's already imported" — is usually the
wrong one.

## Place a type by the concept it models

A type belongs to the concept it models, **not** to the command or package that
first emitted it. A type that several commands emit as part of an agent-facing
contract is its own concept; give it a home named for the contract concept
itself (e.g. a `receipt` package), not a generic bucket. Resist the urge to
park it in a catch-all `types`, `common`, or `util` package — those names tell
the reader nothing, attract unrelated dependencies, and reproduce the very
misplacement this guide warns against.

The tempting shortcut to resist: *"package `cli` already imports `lint`, so
`cli` can just reuse `lint.Action`."* An existing import justifies
**convenience**, not **ownership**. By that logic any type could live anywhere
the import graph already reaches — which is how a "next action" contract ends up
named after, and owned by, whichever feature happened to emit one first.

The reader pays for misplacement on every read: `cli.InitReceipt` referencing
`lint.Action` reads as though `init` has something to do with linting. It does
not. The type is an *agent-receipt* concept; `lint` was merely its first
emitter.

## Two axes, two different rules

The "wait until you have a third use" instinct (the rule of three) is sound — but
only for one of these axes. Keep them apart:

- **Abstraction extraction** — pulling shared *behavior* into a common helper.
  Here the rule of three applies: don't build the shared abstraction until real
  duplication proves the shape. Defer.
- **Concept ownership** — deciding which package *owns* a contract or
  protocol type. Here consumer count is irrelevant. A contract is its own
  concept with one consumer or three; counting emitters answers the wrong
  question. Place it correctly up front.

Conflating the two leads to "we only have two consumers, so leave the type
where it was born" — applying an abstraction-extraction rule to an
ownership decision.

## Is it a contract type or a feature-internal type?

Ask: **would a second command emit this as part of its agent-facing output?** If
yes, it is a contract type and wants a neutral home. Signals:

- It marshals to a stable JSON shape agents depend on (a receipt element, an
  error object, a `schemaVersion`-bearing document).
- Its fields describe a *protocol* (what the agent should do next, how an error
  is reported) rather than a *feature's internals* (a lint rule, a model node).

Feature-internal types — ones only their own package emits and reasons about —
stay in the feature package.

## Cost note: get contract types right up front

Where a Go struct lives is an internal detail. If it marshals to a stable JSON
shape, *moving the struct between packages does not change the wire contract* —
agents never notice. So the mechanical cost of relocating a contract type later
is low.

That is exactly why you should place it correctly now rather than bank on the
cheap move: the refactor is safe but easy to defer forever, while the
readability cost (a type named for the wrong concept) is paid continuously until
someone does it. For a contract type, the right home is cheap to choose up front
and quietly expensive to leave wrong.
