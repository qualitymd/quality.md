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

The same requirements bar applies in two contexts:

- **Durable specs** in `specs/` and `SPECIFICATION.md` are the cumulative source
  of truth for current behavior. They carry enduring rationale and stable
  contracts a future maintainer can evaluate without reading archived changes.
- **Change-case specs** in `changes/NNNN-*/spec.md` are delta contracts. They say
  what this change must alter, prove the requirements are good enough to design
  against, and map any durable spec edits through
  [Durable spec changes](#durable-spec-changes).

The principles below apply to both. Only the job-specific sections call out
extra obligations for one context.

## Shape

Keep specs prose-first and skimmable, and order them for a reader, not a parser:
lead with enough context for the requirements to make sense, then the specifics.
A spec read top to bottom should land — a reader shouldn't have to reassemble it
out of order. The elements below are a palette, not a checklist (see
[Conventions](#conventions)); a typical spec draws on:

- **A title.**
- **A companion note** — what this spec governs, and a link to the source of
  truth it defers to (the format itself lives in
  [`SPECIFICATION.md`](../../SPECIFICATION.md)). Distinguish **normative**
  references — the binding sources of truth this spec defers to, whose rules it
  inherits — from **informational** ones that supply only context; a reader
  should be able to tell which links *bind*.
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
- **Assumptions & dependencies** — external facts the requirements rest on that
  this spec does not control: another component's behavior, a file or folder
  layout, an environment guarantee. List them when a *change* to one would
  invalidate a requirement here. This is distinct from **Scope** (what's covered)
  and **Background** (the *why*): an assumption is a load-bearing fact that, if it
  shifts, should flag the requirements that depend on it rather than let them fail
  silently.
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

Phrase each requirement in the active voice with an explicit subject — the actor
or surface under obligation — followed by the keyword and the result:
*"`qualitymd lint` MUST exit non-zero…"*, not *"a non-zero exit is required"*.
Keep the BCP 14 keyword as the obligation verb; do not rewrite `MUST` as `shall`.

### Requirement quality bar

Every requirement should survive this pass before it drives design or
implementation:

- **One obligation.** A reader can point to the single thing the requirement
  requires or forbids.
- **Named surface.** The actor, command, workflow, artifact, file, or model
  surface under obligation is explicit.
- **Unambiguous.** The requirement admits one interpretation. Vague predicates
  ("better", "robust", "user-friendly", "as needed") are the usual tell — replace
  them with the observable result.
- **Bounded condition.** The trigger, input state, lifecycle phase, or scope is
  named when behavior depends on it.
- **Observable result.** The required outcome is visible in output, persisted
  data, source state, documented behavior, or a user interaction.
- **Divergence handled.** Existing error, conflict, empty-input, stale-state, or
  unsafe cases that could yield plausible different behavior are decided or
  deliberately deferred.
- **Verification path.** The author can name how to check it: a test, fixture,
  command run, rendered artifact, source inspection, or documented review.
- **Earned constraint.** The requirement follows from the change's motivation,
  an existing contract, a real user or agent need, a safety boundary, or an
  observed failure mode — not from speculative completeness. A future maintainer
  should be able to tell which of those reasons earns the imperative.
- **Public contract, not private method.** Implementation details appear only
  when they are themselves the compatibility surface.

If a requirement cannot pass the bar, either sharpen it, move the uncertainty to
**Open questions** or **Deferred**, or delete it. A weak requirement often hides
as a vague verb:

```markdown
Weak: Status output **MUST** be better for incomplete runs.

Strong: `qualitymd evaluation status <run>` **MUST** report a run as incomplete
when any planned in-scope Requirement lacks a Requirement Rating Result, and the
JSON output **MUST** include the missing Requirement IDs.
```

When the reason would not be obvious from the nearby Background / Motivation,
scenario, or surrounding requirements, add a short
[per-requirement rationale](#per-requirement-rationale). This is especially
important for prohibitions, surprising defaults, consent gates, compatibility
boundaries, naming choices, agent workflow constraints, and rules that choose
one plausible behavior over another.

### The requirement set

The bar above tests one requirement at a time; a set of individually sound
requirements can still fail *as a set*. Before a spec drives design, check the
requirements together for:

- **Consistent.** No requirement conflicts with or duplicates another, and one
  term means one thing throughout — the same surface, state, or artifact is named
  the same way every time.
- **Complete.** The set covers the need with no unresolved "to be decided"
  placeholders left inline; a genuinely open point goes under **Open questions**
  or **Deferred**, not buried in a requirement.
- **Able to be validated.** Satisfying the whole set would actually achieve the
  motivation — the requirements add up to the outcome, with nothing essential
  missing and nothing pulling against it.

This complements the per-requirement bar; it is the requirements-engineering
distinction between well-formed *requirements* and a well-formed *set*. For a
change-case spec, the set-level pass pairs
with the [validation check](work-with-change-cases.md#write-the-spec-then-the-design)
at the Draft→Design boundary.

### A statement template (optional)

When a requirement is awkward to phrase — especially one that fires on a trigger
or handles an edge case — a small set of statement patterns gives a reliable
shape. They are a **recommended aid, not a required form**: reach for one when it
sharpens a requirement, and leave prose alone when prose is already clear. Each
pattern puts a BCP 14 keyword in the response slot.

| Pattern            | Shape                                                          | Use for                                         |
| ------------------ | -------------------------------------------------------------- | ----------------------------------------------- |
| Ubiquitous         | `<surface>` **MUST** `<result>`                                | an always-on obligation                         |
| State-driven       | **While** `<state>`, `<surface>` **MUST** `<result>`           | behavior that holds during a condition          |
| Event-driven       | **When** `<trigger>`, `<surface>` **MUST** `<result>`          | a response to an event or input                 |
| Optional-feature   | **Where** `<feature present>`, `<surface>` **MUST** `<result>` | behavior tied to an optional flag or capability |
| Unwanted-behaviour | **If** `<condition>`, **then** `<surface>` **MUST** `<result>` | errors, conflicts, empty input, unsafe state    |
| Complex            | combine the above                                              | a trigger qualified by a state                  |

Two of these map straight onto the [quality bar](#requirement-quality-bar): the
**When** form is its **Bounded condition** made explicit, and the **If/then**
form is **Divergence handled** — the syntactic home for the cases
["an unspecified case is a decision delegated"](#conventions) tells you to decide.

> When `qualitymd lint` finds errors, it **MUST** exit non-zero.
>
> If the run folder already exists, then `qualitymd evaluation create` **MUST**
> refuse and report the collision.

### Per-requirement rationale

A requirement **may** carry a subordinate rationale annotation directly beneath
it — the fine-grained counterpart to **Background / Motivation**. Use one when
the requirement's reason is load-bearing but not obvious from the nearby prose.
Lead with the testable sentence; put the *why* in a blockquote under it:

> A command **MUST** exit non-zero when it finds errors.
>
>> Rationale: CI gates on exit code alone; a zero exit on lint failure
>> green-lit broken files in practice. — 0012

The form: a blockquote led by `Rationale:` (the terser `Why:` is fine), one or
two sentences, optionally citing the originating change id (`— 0012`) for
provenance. The requirement stays the lead; the annotation is subordinate and
must never wrap around or bury it. A change-case spec requirement carries a second
subordinate annotation alongside this one — the `Durable spec:` line that records
its spec impact (see [Durable spec changes](#durable-spec-changes)); the two are
distinct (one carries *why*, the other maps to the durable contract) and each sits
on its own blockquote.

Annotate by this litmus: **when a future editor would otherwise repeat a mistake,
miss the user or agent need, break a safety or consent boundary, weaken
compatibility, or be misled.** A rule whose reason is obvious needs no note, but
an imperative with no visible reason should be treated as suspect: add rationale,
weaken it, or remove it. Dead-end alternatives and the full decision record stay
in the (archived) [design doc](write-design-docs.md); only the durable intent and
the lesson get promoted onto the requirement. Background carries the spec-scale
*why*; an annotation carries one requirement's — they must not restate each other,
and stale rationale is superseded, not left to accrete.

## Durable spec changes

When a functional spec is a change case's `spec.md` (not a durable
[`specs/`](../../specs/index.md) concept), it has one extra job: account for the
durable **specs** the change rewrites. This section is the bridge from the
change-case delta to the cumulative source of truth.

Unlike the [Shape](#shape) sections, which are a palette, this section is
**required** for a change-case spec. A silent omission is how a contract change
lands undocumented, and the *what* of a spec change belongs with the spec, not
buried in the design doc. Durable specs do not carry this section; they absorb
the resulting contract and rationale directly. (See
[Working with change cases](work-with-change-cases.md#account-for-the-artifacts-it-touches)
for how it divides labor with the parent's **Affected artifacts** index.)

Attribute the substance **per requirement first**. Each requirement carries a
subordinate `Durable spec:` annotation — parallel to [`Rationale:`](#per-requirement-rationale),
a blockquote directly beneath the testable sentence — naming the durable-spec
change it drives: the verb (**add** / **modify** / **rename** / **delete**), the
spec path, and what changes. A requirement that drives no durable-spec change
says `Durable spec: none`, so it is explicit which requirements result in a spec
change and which deliberately do not — never a silent omission. The annotation
**MUST NOT** restate the normative text; it maps the requirement to the durable
contract it lands in.

When a requirement creates a new public, durable, named artifact, first apply
the artifact-contract test: will readers, tools, generated outputs, or workflows
rely on this file by name, location, shape, lifecycle, or link target? If yes,
the requirement's annotation should usually be `Durable spec: add <spec path>`
for a 1:1 artifact spec, or it must name the existing durable spec that owns the
complete artifact contract and why no new spec is needed. A new artifact listed
only as a durable doc, with no durable-spec owner or justification, is
under-accounted.

> The single-kind form **MUST** emit a self-contained schema.
>
>> Durable spec: modify `specs/cli/evaluation-data.md` — replace "rooted at that
>> kind" with the self-contained-legibility requirement.

Then **roll those annotations up** into a `## Durable spec changes` section with
four subsections, in order — **To add** (new durable specs), **To modify**, **To
rename**, **To delete** — each present and explicit: a list, or the single word
`None`. The rollup is the at-a-glance footprint and the completeness check; it
also houses any change no single requirement owns (a rename or delete that serves
the requirement set as a whole). A subsection is never left blank or dropped; a
written `None` is a deliberate "nothing of this kind," not an oversight. **To
rename** carries any durable spec whose path changes (`old → new`): a rename is a
delete-plus-add, so without its own home it hides — the new path slips into
**To add** and the old path goes unaccounted in **To delete**. Naming the rename
once keeps both ends honest. The per-requirement annotations are authoritative;
each rollup entry back-links the requirement(s) that drive it and **MUST NOT**
restate the normative text — that lives once, in [Requirements](#requirements).

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
the format spec [`SPECIFICATION.md`](../../SPECIFICATION.md), including the skill's
functional spec under [`specs/skills/`](../../specs/index.md) when a change alters
skill *behavior*. The other artifact kinds — durable *docs* (README, guides,
scaffold, the bundled skill's runtime content under `skills/`) and *code* — are
tracked in the change case's parent **Affected artifacts** index, not here.

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
  concept extension. Example: an `evaluation-output-result.json` artifact spec
  would use `evaluation-output-result-json.md`. This keeps the artifact identity
  visible without creating filenames with multiple operational meanings.
- **Do not bury new artifact contracts in adjacent specs.** If a change adds a
  named artifact such as `glossary.md`, prefer a 1:1 artifact spec such as
  `glossary-md.md` unless an existing durable spec truly owns the whole artifact:
  purpose, location, structure, maintenance source, link behavior, and reader
  contract. If you choose the existing-spec route, state that ownership
  explicitly in the change case.
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
  specs here already do it — the Evaluation routine-output examples in
  [Evaluation](../../specs/evaluation/evaluation.md), the `nextActions`
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
