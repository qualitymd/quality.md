---
type: How-to Guide
title: Writing functional specs
description: How to write a functional spec for the qualitymd tooling.
tags: [specs, contributing]
timestamp: 2026-06-18T00:00:00Z
---

# Writing functional specs

A **functional spec** describes *what* a piece of `qualitymd` does and the
requirements it must meet — not how it's implemented. Our specs live in the
[`specs/`](../../specs/index.md) bundle, authored as [OKF](work-with-okf.md) concepts.

## Shape

Keep specs prose-first and skimmable, and order them for a reader, not a parser:
lead with enough context for the requirements to make sense, then the specifics.
A spec read top to bottom should land — a reader shouldn't have to reassemble it
out of order. The elements below are a palette, not a checklist (see
[Conventions](#conventions)); a typical spec draws on:

- **A title.**
- **A companion note** — what this spec governs, and a link to the source of
  truth it defers to (the format itself lives in
  [`SPECIFICATION.md`](../../SPECIFICATION.md)).
- **Background / Motivation** — a short prose section near the top stating the
  big-picture *why*: the problem or failure-mode the capability addresses, and
  any spec-scale lessons worth carrying forward. This is distinct from **Scope**
  (which says *what's* covered or deferred), from the companion note (which says
  what the spec governs), and from a **Scenario** (which names the self-contained
  case the spec solves, not the broader why). It exists so a spec's durable
  rationale lives with the spec, not only in the change that introduced it. Keep
  it to a paragraph or so; the fine-grained *why* goes on individual requirements
  (see below).
- **Scenario / use case** — when the spec exists to satisfy a self-contained use
  case, state it: the concrete thing a caller is trying to accomplish that this
  spec, on its own, makes possible. Draw the line against **Background /
  Motivation**: a scenario that's merely one step of a larger process or
  job-to-be-done *is* background — context for the big-picture *why*. But when the
  spec stands on its own as the answer to a particular case, state that case
  directly — a short walkthrough of who reaches for the capability and what they
  need to come away with — so the requirements that follow read as its solution.
  Keep it concrete and brief; many specs need only background.
- **Scope** — what's covered now, and what's left out *on purpose*, so an absence
  reads as deliberate rather than forgotten. Two kinds of absence belong here:
  **deferred** (real, but not yet — recorded so it isn't re-litigated) and
  **non-goals** (out of scope by design, e.g. *the CLI never calls a model*).
  Naming a non-goal kills the recurring "should it also…?" before it's asked.
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

### Per-requirement rationale

A requirement **may** carry a subordinate rationale annotation directly beneath
it — the fine-grained counterpart to **Background / Motivation**. Lead with the
testable sentence; put the *why* in a blockquote under it:

> A command **MUST** exit non-zero when it finds errors.
>
>> Rationale: CI gates on exit code alone; a zero exit on lint failure
>> green-lit broken files in practice. — 0012

The form: a blockquote led by `Rationale:` (the terser `Why:` is fine), one or
two sentences, optionally citing the originating change id (`— 0012`) for
provenance. The requirement stays the lead; the annotation is subordinate and
must never wrap around or bury it.

Annotate by this litmus: **when a future editor would otherwise repeat a mistake
or be misled.** A rule whose reason is obvious needs no note. Dead-end
alternatives and the full decision record stay in the (archived)
[design doc](write-design-docs.md); only the durable intent and the lesson get
promoted onto the requirement. Background carries the spec-scale *why*; an
annotation carries one requirement's — they must not restate each other, and
stale rationale is superseded, not left to accrete.

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
- **An unspecified case is a decision delegated.** A spec *decides* behavior; it
  doesn't just describe the happy path. Each error, conflict, empty input, or
  concurrent-use case left unstated isn't left flexible — it's handed silently to
  whoever writes the code (or the agent driving it), decided ad hoc and invisibly.
  For a deterministic tool that's the whole game: spec the cases where behavior
  could plausibly diverge so one input can't yield two defensible results. This is
  the complement to YAGNI, not its contradiction — decide the cases that *exist*
  (the run folder already collides; `--subject` can already point at a directory);
  don't invent ones that don't.
- **Sections are a palette, not a checklist.** The [Shape](#shape) list is what a
  spec *may* need, not boxes every spec must fill. Take only what a spec earns: a
  cross-cutting spec may use them all; a single-command spec is often just its
  companion note and requirements (see
  [`evaluation create-run`](../../specs/cli/evaluation-create-run.md)). Forcing
  every heading onto every spec is how specs bloat and authors stall. The OKF
  mechanics — frontmatter `type`, a tidy `index.md` and `log.md` — are the only
  non-negotiables.
- **Say it once.** This applies *within* a spec too: each requirement gets one
  home. An overview or principles list should *name* a property and link to the
  section that enforces it, not re-assert the requirement in full. A spec that
  introduces its properties up front and then repeats each as its own normative
  section is saying everything twice — collapse the overview to links.
- **Two whys, each in its place.** Rationale earns its place — durable specs
  should carry the reasons their requirements exist, or those reasons die in the
  archived change. Split it by grain. The big-picture *why* (the problem or
  failure-mode the capability addresses) goes in
  [Background / Motivation](#shape). A single requirement's *why* goes in a
  subordinate [annotation](#per-requirement-rationale) beneath it — a sentence or
  two, not a paragraph wrapped around the rule. The failure mode to avoid is
  rationale that *buries or outweighs* the requirement, not rationale itself:
  keep the requirement the lead, testable sentence and let the *why* sit under
  it.
- **Show, don't only tell.** A concrete example often pins a contract more
  precisely than prose and doubles as a test: a sample invocation, a fenced record
  or folder layout, a representative JSON receipt, an overview table. The sharpest
  specs here already do it — the run-folder layout in
  [Evaluation records](../../specs/evaluation-records.md), the `nextActions`
  receipt in the [CLI spec](../../specs/cli.md). Reach for an example where words
  get slippery; don't decorate.
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
- **Rationale that buries the rule** — justification wrapped around a one-line
  MUST until the requirement is no longer the lead, testable sentence. The cure
  is order and subordination, not deletion: lead with the rule, drop the *why*
  into a [`Rationale:` annotation](#per-requirement-rationale) (or
  [Background](#shape)) beneath it.
- **Rationale said twice** — a per-requirement annotation that restates what
  [Background / Motivation](#shape) already says, or vice versa. Background
  carries the spec-scale *why*; an annotation carries one requirement's. Say each
  once, in its own place, and supersede stale rationale rather than letting it
  accrete.
- **Implementation detail leaking in** — exact env-var names, escape sequences,
  binary-naming arguments. State the behavior; let the code carry the mechanics.
