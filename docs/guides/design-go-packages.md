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

It is about _data and contract types_ — the structs that model a concept or
marshal to an agent-facing shape. Interfaces follow a different placement rule:
define them in the package that consumes them, not by the concept they abstract.
See [Go style](go-style.md) for that and the rest of the line-level conventions.

## Place by concept, not emitter

A type belongs to the concept it models, **not** to the command or package that
first emitted it. A type that several commands emit as part of an agent-facing
contract is its own concept; give it a home named for the contract concept
itself (e.g. a `receipt` package), not a generic bucket. Resist the urge to
park it in a catch-all `types`, `common`, or `util` package — those names tell
the reader nothing, attract unrelated dependencies, and reproduce the very
misplacement this guide warns against.

This repo faced exactly this choice. `lint` defined an `Action` type for its
findings, and when `init` needed a "next action" for its receipt, the tempting
shortcut was _"package `cli` already imports `lint`, so just reuse
`lint.Action`."_ An existing import justifies **convenience**, not
**ownership**. By that logic any type could live anywhere the import graph
already reaches — which is how a "next action" contract ends up named after, and
owned by, whichever feature happened to emit one first.

The reader would have paid for that on every read: `cli.InitReceipt` referencing
`lint.Action` reads as though `init` had something to do with linting. It does
not — the type is an _agent-receipt_ concept; `lint` was merely its first
emitter. So it got a neutral home instead: `Action` lives in `receipt`, and
`cli.InitReceipt` consumes `receipt.Action`.

## Two axes

The "wait until you have a third use" instinct (the rule of three) is sound — but
only for one of these axes. Keep them apart:

- **Abstraction extraction** — pulling shared _behavior_ into a common helper.
  Here the rule of three applies: don't build the shared abstraction until real
  duplication proves the shape. Defer.
- **Concept ownership** — deciding which package _owns_ a contract or
  protocol type. Here consumer count is irrelevant. A contract is its own
  concept with one consumer or three; counting emitters answers the wrong
  question. Place it correctly up front.

Conflating the two leads to "we only have two consumers, so leave the type
where it was born" — applying an abstraction-extraction rule to an
ownership decision.

## Contract vs. feature-internal

Ask: **would a second command emit this as part of its agent-facing output?** If
yes, it is a contract type and wants a neutral home. Signals:

- It marshals to a stable JSON shape agents depend on (a receipt element, an
  error object, a `schemaVersion`-bearing document).
- Its fields describe a _protocol_ (what the agent should do next, how an error
  is reported) rather than a _feature's internals_ (a lint rule, a model node).

Feature-internal types — ones only their own package emits and reasons about —
stay in the feature package.

## Cost of deferring

Where a Go struct lives is an internal detail. If it marshals to a stable JSON
shape, _moving the struct between packages does not change the wire contract_ —
agents never notice. So the mechanical cost of relocating a contract type later
is low.

That is exactly why you should place it correctly now rather than bank on the
cheap move: the refactor is safe but easy to defer forever, while the
readability cost (a type named for the wrong concept) is paid continuously until
someone does it. For a contract type, the right home is cheap to choose up front
and quietly expensive to leave wrong.
