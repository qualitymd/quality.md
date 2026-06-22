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
- **Sub-specs** — split detail into child concepts when a spec grows or when a
  durable behavior is independently reviewable. The parent carries the shared
  contract; children carry the specifics. Before reshaping a large spec,
  inventory its major headings and classify each one by durable contract:
  shared invariant, behavioral component, artifact contract, example, or
  deferred work. A heading classified as a behavioral component or artifact
  contract needs either a child spec or a short written reason it stays in the
  parent.

## Requirements

State requirements as clear, testable obligations. Use BCP 14 keywords
(`MUST`, `SHOULD`, `MAY`) only when the keyword changes conformance meaning:
required behavior, prohibited behavior, a default with a real escape hatch, or
optional behavior that affects compatibility. Normal prose can use lowercase
"must", "should", "may", and "can" with their ordinary English meanings.

> A command **MUST** exit non-zero when it finds errors. Output **SHOULD** default
> to a human-readable form.

When a spec uses BCP 14 keywords, declare the convention once near the top:

> The key words "MUST", "SHOULD", and "MAY" are to be interpreted as described
> in BCP 14 when, and only when, they appear in all capitals.

Do not add BCP 14 keywords just to make a sentence feel spec-like. A requirement
can be normative without a keyword when the obligation is still clear and
testable.

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

## Durable spec changes

When a functional spec is a change case's `spec.md` (not a durable
[`specs/`](../../specs/index.md) concept), it also accounts for the durable
**specs** the change rewrites — the bridge from this delta to the cumulative
source of truth. Unlike the [Shape](#shape) sections, which are a palette, this
section is **required** for a change-case spec: a silent omission is how a
contract change lands undocumented, and the *what* of a spec change belongs with
the spec, not buried in the design doc. (See
[Working with change cases](work-with-change-cases.md#account-for-the-artifacts-it-touches)
for how it divides labor with the parent's **Affected artifacts** index.)

Give it four subsections, in order — **To add** (new durable specs), **To
modify**, **To rename**, **To delete** — and make each one present and explicit:
a list, or the single word `None`. A subsection is never left blank or dropped; a
written `None` is a deliberate "nothing of this kind," not an oversight. **To
rename** carries any durable spec whose path changes (`old → new`): a rename is a
delete-plus-add, so without its own home it hides — the new path slips into
**To add** and the old path goes unaccounted in **To delete**. Naming the rename
once keeps both ends honest.

Each entry names the durable spec, says what changes, and links to the
requirement above that drives it. It **MUST NOT** restate the normative text —
that lives once, in [Requirements](#requirements); the entry only maps a
requirement to the durable contract it lands in.

```markdown
## Durable spec changes

### To add

None

### To modify

- `specs/cli/lint.md` — add the `uncharacterized-requirement` error row (per the
  lint requirement above).
- `SPECIFICATION.md` — require factor characterization; reframe "secondary
  factor" (per the characterization requirements above).

### To rename

- `specs/cli/check.md` → `specs/cli/lint.md` — track the command rename so the
  old path's removal and the new path's arrival are both accounted (per the
  rename requirement above).

### To delete

None
```

Scope is durable **specs** only — the [`specs/`](../../specs/index.md) bundle and
the format spec [`SPECIFICATION.md`](../../SPECIFICATION.md). The other artifact
kinds — durable *docs* (README, guides, scaffold, the bundled skill) and *code* —
are tracked in the change case's parent **Affected artifacts** index, not here.

## Conventions

- **Specify behavior, not implementation.** Say *what* must hold; leave *how* to
  the code.
- **One source of truth.** Don't restate the format spec — link to it.
- **Inventory before splitting.** Before reshaping a large spec, list its major
  headings and classify each one as a shared invariant, behavioral component,
  artifact contract, example, or deferred work. A behavioral component or
  artifact contract should become a child spec unless it is too small or too
  entangled to review independently; when it stays in the parent, say why. This
  keeps the split from stopping at the first obvious boundary while another
  full workflow remains buried in the parent.
- **Split by durable contract, not by file layout.** A child spec earns its
  place when the contract is independently understandable and reviewable: a
  command, workflow, artifact, evaluator phase, lifecycle, or orchestration path
  with its own purpose, inputs, outputs, mutation boundary, safety rules,
  compatibility surface, or done criteria. Do not split only because the
  implementation has a separate file; implementation files can stay governed by
  a parent spec when they are mostly procedure under shared rules.
- **Keep shared invariants in the parent.** Parent specs hold vocabulary,
  global safety rules, argument models, common output posture, shared artifact
  relationships, and cross-component invariants. Child specs hold the behavior a
  reader would naturally review in isolation: routing, state transitions,
  component-specific stop conditions, mutation rules, required artifacts, and
  verification.
- **Use fictional examples to test the boundary.** In a fictional "billing
  assistant" spec, the parent might own account vocabulary, authentication
  rules, and common audit-log invariants. Separate child specs could own the
  "invoice reconciliation" workflow, the "refund approval" lifecycle, and the
  `refund-receipt.json` artifact. The point is not the folder shape; it is that
  each child carries a contract a reviewer can understand and verify on its own.
- **Name 1:1 artifact specs after the artifact.** When a durable spec is scoped
  to one concrete generated file or artifact, preserve the artifact filename in
  the spec filename by replacing `.` with `-` and then using the normal `.md`
  concept extension. Examples: the `report-summary.md` artifact is specified by
  `report-summary-md.md`; `report.json` by `report-json.md`; `report.md` by
  `report-md.md`. This keeps the artifact identity visible without creating
  filenames with multiple operational meanings.
- **Name behavioral component specs after the capability.** When a durable spec
  governs behavior rather than one concrete artifact, name it for the command,
  workflow, phase, lifecycle, or component it specifies. For example, a
  fictional spec for refund approval behavior should use a capability name such
  as `refund-approval.md`; reserve artifact-normalized names such as
  `refund-approval-md.md` for a contract over a literal generated or bundled
  `refund-approval.md` file.
- **Don't invent requirements (YAGNI).** Specify only what's actually asked for
  or genuinely needed now. A spec is not the place to anticipate features,
  hedge against hypothetical needs, or add flags, formats, and edge cases
  "while we're here." Every requirement is a constraint someone has to implement
  and uphold; speculative ones cost more than they save. When a need is real
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
  [`evaluation create`](../../specs/cli/evaluation-create.md)). Forcing
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
  requirement until it is no longer the lead, testable sentence. The cure
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
