# Writing functional specs

A **functional spec** describes *what* a piece of `qualitymd` does and the
requirements it must meet — not how it's implemented. Our specs live in the
[`specs/`](../../specs/index.md) bundle, authored as [OKF](okf.md) concepts.

## Shape

Keep specs prose-first and skimmable. A typical spec has:

- **A title and a one-line draft version** — e.g. `**Version 0.1 — Draft**`.
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
- **Draft openly.** Mark early specs `Draft`; placeholders are fine — stub the
  sections to be filled and link back to the parent.
- **Keep OKF tidy.** New spec → add a `type` to its frontmatter and update the
  enclosing `index.md` and `log.md` (see the [OKF guide](okf.md)).
