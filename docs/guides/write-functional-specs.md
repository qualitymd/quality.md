---
type: How-to Guide
title: Writing functional specs
description: How to write a functional spec for the qualitymd tooling.
tags: [specs, contributing]
timestamp: 2026-06-16T00:00:00Z
---

# Writing functional specs

A **functional spec** describes *what* a piece of `qualitymd` does and the
requirements it must meet — not how it's implemented. Our specs live in the
[`specs/`](../../specs/index.md) bundle, authored as [OKF](work-with-okf.md) concepts.

## Shape

Keep specs prose-first and skimmable. A typical spec has:

- **A title.**
- **A companion note** — what this spec governs, and a link to the source of
  truth it defers to (the format itself lives in
  [`SPECIFICATION.md`](../../SPECIFICATION.md)).
- **Scope** — what's covered now, and what's intentionally **deferred**, so an
  absence reads as deliberate rather than forgotten.
- **Requirements** — the normative content (see below).
- **Sub-specs** — split detail into child concepts when a spec grows; the parent
  carries the shared contract, children carry the specifics.

## Requirements

State requirements with RFC 2119 keywords — **MUST**, **SHOULD**, **MAY** — so
each is testable and its strength is unambiguous:

> A command **MUST** exit non-zero when it finds errors. Output **SHOULD** default
> to a human-readable form.

Declare the keyword convention once near the top:

> The key words "MUST", "SHOULD", "MAY" … are to be interpreted as described in
> IETF RFC 2119.

## Conventions

- **Specify behavior, not implementation.** Say *what* must hold; leave *how* to
  the code.
- **One source of truth.** Don't restate the format spec — link to it.
- **Don't invent requirements (YAGNI).** Specify only what's actually asked for
  or genuinely needed now. A spec is not the place to anticipate features,
  hedge against hypothetical needs, or add flags, formats, and edge cases
  "while we're here." Every MUST/SHOULD is a constraint someone has to implement
  and uphold — speculative ones cost more than they save. When a need is real
  but not yet, record it under **deferred** rather than specifying it.
- **Say it once.** This applies *within* a spec too: each requirement gets one
  home. An overview or principles list should *name* a property and link to the
  section that enforces it, not re-assert the requirement in full. A spec that
  introduces its properties up front and then repeats each as its own normative
  section is saying everything twice — collapse the overview to links.
- **Motivation in asides, not paragraphs.** Rationale earns its place (see
  [Shape](#shape)), but keep it to a clause or a `Note:` next to the requirement
  it justifies. When the *why* grows past a sentence or two, it's usually a sign
  the *what* has been buried.
- **Draft openly.** Mark early specs `Draft`; placeholders are fine — stub the
  sections to be filled and link back to the parent.
- **Keep OKF tidy.** New spec → add a `type` to its frontmatter and update the
  enclosing `index.md` and `log.md` (see the [OKF guide](work-with-okf.md)).

## Smells

If a spec feels bloated, look for these — each is over-specification, not
thoroughness:

- **The same requirement in two sections** — e.g. a "Design requirements"
  overview *and* an "Agent accessibility" section both asserting the JSON-output
  rule. Pick one home; link from the other.
- **A principles list that's really a second copy** of the normative body. If
  every bullet maps to a section that says the same thing, the list should be
  links.
- **Rationale that outweighs the rule** — a paragraph of justification wrapped
  around a one-line MUST.
- **Implementation detail leaking in** — exact env-var names, escape sequences,
  binary-naming arguments. State the behavior; let the code carry the mechanics.
