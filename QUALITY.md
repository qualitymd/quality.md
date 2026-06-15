---
factors:
  functionality:
    requirements:
      "the format specification is complete":
        target: "./SPECIFICATION.md"
        prompt: >
          Every part of the QUALITY.md format is specified: each frontmatter
          field and recommended body section has its shape, allowed values,
          whether it is required or optional, its cardinality, and any default.
          The spec states how a conforming reader treats malformed or omitted
          content, and addresses edge cases — empty file, empty values,
          duplicate keys, ordering, and case-sensitivity. An implementer could
          build a parser and an author could write a valid file from the spec
          alone, without having to invent undefined behavior.
      "the CLI specification is complete":
        target: "./specs/cli.md"
        prompt: >
          The CLI specification defines the qualitymd command surface well enough
          to implement and to test against: each command's purpose, arguments,
          flags, inputs, output, and exit codes, and how the deterministic CLI
          and the skill layer relate. A developer could implement the CLI, and
          write conformance tests for it, from the specification without
          inventing unspecified behavior.
      "the skill specification is complete":
        target: "./specs/skills.md"
        prompt: >
          The skill specification defines the qualitymd skills well enough to
          implement and to test against: each skill's purpose, the loop it runs
          over the CLI's model / evaluation / result resources, the CLI ↔ skill
          interface it relies on, and where judgment lives. A developer could
          implement the skills, and write conformance tests for them, from the
          specification without inventing unspecified behavior.
      "the CLI implementation conforms to its specification":
        target:
          - "./internal"
          - "./cmd"
        prompt: "./specs/cli.md"
  clarity:
    requirements:
      "the format specification admits a single interpretation":
        target: "./SPECIFICATION.md"
        prompt: >
          Each rule admits a single interpretation. Obligations are signaled
          with a consistent requirement vocabulary (must / should / may) and
          their strength is unambiguous, and no normative statement leans on a
          vague quantifier — "fast", "large", "reasonable" — without a stated
          bound.
      "the format specification separates rules from rationale":
        target: "./SPECIFICATION.md"
        prompt: >
          A reader can always tell whether a sentence states a binding rule or
          merely explains, motivates, or illustrates one, and a rule never
          appears only inside an example or an aside.
      "the format specification defines its terms before use":
        target: "./SPECIFICATION.md"
        prompt: >
          Every technical term the specification uses in a rule is defined
          before, or at, the point it is first used, so a reader never has to
          read ahead to learn what a rule means.
  consistency:
    requirements:
      "the format specification is internally consistent":
        target: "./SPECIFICATION.md"
        prompt: >
          No two statements contradict each other. One term denotes one concept
          throughout — the same idea is not renamed across sections, and one
          name is not reused for two ideas. Every example agrees with the rules
          it illustrates.
  verifiability:
    requirements:
      "each rule is observable or testable":
        target: "./SPECIFICATION.md"
        prompt: >
          Each rule maps to something a reader could observe or test about a
          file or an implementation, rather than a matter of taste, so that two
          independent readers would decide conformance the same way.
      "the format's constructs are shown with valid and invalid examples":
        target: "./SPECIFICATION.md"
        prompt: >
          The format's constructs are shown with worked examples that include
          both valid cases and invalid (counter-)examples, so the boundary of
          each rule is pinned down by example as well as by prose.
  extensibility:
    requirements:
      "the format specifies its core and how it extends and evolves":
        target: "./SPECIFICATION.md"
        prompt: >
          The spec names the minimal core every QUALITY.md must have and says
          whether and how authors may add their own factors, keys, or sections
          beyond it. It states how a conforming reader treats content it does
          not recognize (ignore vs. reject), and how the format is versioned or
          expected to evolve without silently breaking existing files.
  usability:
    requirements:
      "the format specification is well-structured and readable":
        target: "./SPECIFICATION.md"
        prompt: >
          Sections follow a logical order that introduces a concept before the
          detail that depends on it. Field tables and a clear schema make the
          structure scannable. Prose is plain and direct, and the document
          carries at least one minimal example and one fuller, realistic example
          a reader can copy and adapt.
  approachability:
    requirements:
      "the README says what QUALITY.md is and who it's for":
        target: "./README.md"
        prompt: >
          A first-time reader learns within the opening lines — before any
          install or usage detail — what a QUALITY.md file is (a plain-text
          quality model: YAML frontmatter plus a Markdown body), what problem
          the project solves, and who it is for (authors, coding agents, CI).
          The framing names a concrete subject and payoff rather than leaning
          on abstract adjectives.
      "the README shows the format and its payoff by example":
        target: "./README.md"
        prompt: >
          The README shows, not just describes, the format: it includes a
          realistic QUALITY.md excerpt — frontmatter and body together — and
          makes the payoff concrete by showing what running qualitymd against
          it produces, such as a sample of the grouped pass/fail report. A
          reader sees the artifact and its effect early, before the reference
          detail.
      "the README gets a newcomer to a first result quickly":
        target: "./README.md"
        prompt: >
          A newcomer can copy a short, self-contained sequence — install, then
          a single command — and see a real result, with the command shown
          together with representative output rather than the invocation alone.
          Where a command's exit code is meant to gate CI, the README says so.
      "the README reflects what the CLI and spec actually provide":
        target: "./README.md"
        prompt: >
          Every command, flag, and capability the README shows matches what the
          CLI provides today; anything specified but not yet built is marked as
          planned rather than presented as available, and any command the CLI
          still carries that the spec does not define — an off-spec placeholder —
          is shown as provisional rather than as a settled feature. The README
          states the project's maturity plainly and points to the authoritative
          spec(s) for full detail instead of restating rules that could drift
          from them.
---

# Quality model — QUALITY.md

## Overview

This model governs the quality of the **QUALITY.md project itself**: the
**format** — the plain-text quality-model format defined by
[`SPECIFICATION.md`](./SPECIFICATION.md) — the deterministic **`qualitymd` CLI**
that scaffolds and validates a model and records, rolls up, and reports
evaluation results, and the **skills** that carry the judgment the CLI
deliberately does not. The CLI and skills are specified under
[`specs/`](./specs/); the CLI is implemented under `internal/` and `cmd/`.

At this pre-1.0 stage the project's quality rests on its **design**, so the model
is scoped to the maturity of its specifications, the CLI's conformance to its
own, and the README that introduces the project. "Good" means an implementer
could build a parser and an evaluator, and an author could write a valid file,
from the specifications alone; the CLI behaves as its spec says; and a newcomer
can tell from the README what QUALITY.md is and reach a first result. The format
spec carries the most weight: it is the contract every reader, author,
implementation, and file depends on.

## Scope

This model covers the three deliverables the project owns: the **format**
([`SPECIFICATION.md`](./SPECIFICATION.md)), the **`qualitymd` CLI** (specified
under [`specs/`](./specs/), implemented under `internal/` and `cmd/`), and the
**skills** (specified under [`specs/`](./specs/)). Its weight is deliberately
uneven — the format spec is the contract every reader, author, implementation,
and file depends on, so it carries the most; the CLI and skill specs are held to
functional completeness alone at this stage (see **Known gaps**).

Out of scope by design — not deferred work that belongs here later:

- **Dependencies the project relies on but does not own** — the Go toolchain,
  Cobra/Fang, and release tooling. Their quality is their upstreams' concern, not
  this model's.

The CLI's own runtime product quality — reliability, performance, packaging, test
coverage — is *not* out of scope; it is in-scope work deferred while the
implementation is nascent, recorded under **Known gaps**.

## Needs

- **Format implementers** can build a parser and an evaluator that agree with
  every other conforming implementation, from `SPECIFICATION.md` alone.
- **Authors** can write a valid `QUALITY.md`, and know what each field and
  section means, without reverse-engineering an implementation.
- **Coding agents** can read the spec to author files and read a `QUALITY.md` to
  evaluate a subject, arriving at one interpretation of each rule.
- **CLI users** get a tool whose behavior matches what the specs and `--help`
  describe, with exit codes they can gate CI on.
- **Maintainers and contributors** can evolve the format, the CLI, and the skills
  without silently breaking existing files, because the core and its extension
  and versioning rules are written down.
- **Newcomers** can decide whether QUALITY.md is for them from the README alone —
  grasping what it is and who it's for, seeing the format work on a concrete
  example, reaching a first result of their own, and trusting that what they were
  shown is real.

## Risks

The format is the project's foundation, so its defects are the costliest:

- An **ambiguous or incomplete spec** is the worst outcome — implementations
  diverge, a file that passes one tool fails another, and the format's promise of
  a portable quality model breaks. This is what the model must guard against most.
- A **format with no specified path to grow** forces breaking changes as it
  evolves, stranding files authored against earlier drafts.
- An **implementation that drifts from its spec** — the CLI or the skills behaving
  differently from what their specs describe — gives users and agents behavior the
  documentation does not, eroding trust in the tool.

Prose and structure defects in the spec (a confusing order, a missing example)
are real but recoverable, and rank below these. So too a **README that fails to
land, or overstates what exists** — it turns newcomers away before the internals
are ever seen, and erodes trust when its claims outrun the tool, but the cost is
recoverable, not foundational.

## Factors

### Functionality

Does each deliverable do its job? Functionality is *functional completeness and
correctness* — whether a specification says enough, and accurately enough, to be
built and authored against without inventing undefined behavior, and whether an
implementation actually behaves as its specification says. It is the load-bearing
attribute: a deliverable that is clear, consistent, and verifiable but incomplete
or wrong still fails the people who depend on it. It carries the most weight here,
and in this model it is assessed separately for each deliverable the project owns.

### Clarity

Can each rule be read in only one way? A specification governs independent
implementations and agents that never confer, so a rule open to two readings
becomes two behaviors. Clarity is the attribute that closes that gap — every
obligation lands with one settled meaning and one settled force. It trades brevity
for precision: pinning terms and edge cases down makes a document longer, and
removes the guesswork *Functionality*'s completeness would otherwise leave to the
reader.

### Consistency

Does the document agree with itself? Consistency is the attribute of internal
agreement — across a whole specification one concept keeps one name, no two
statements pull in opposite directions, and illustrations track the rules they
illustrate. It is distinct from *Clarity*: a rule can be locally unambiguous yet
contradict a far-off section or quietly reuse a term defined elsewhere. It is the
attribute most easily lost as a draft grows, because each addition is one more
thing the rest of the document must still square with.

### Verifiability

Can conformance be decided rather than argued? Verifiability is the attribute
that makes a rule checkable — its satisfaction turns on something a reader can
observe or test, so two reviewers reach the same verdict instead of trading
opinions. It is the precondition for tooling: a conformance checker can only
enforce rules whose satisfaction is decidable in the first place, which is why an
unverifiable rule is one no tool and no two readers can hold a subject to. It
splits decidability from demonstration: a rule's satisfaction must turn on an
observable test, and the boundary where conformance flips must be pinned by
worked example — valid cases *and* invalid counter-examples. Those
counter-examples fix a rule by instance rather than by prose, a verifiability
concern distinct from the copy-and-adapt examples *Usability* asks for.

### Extensibility

Can the format grow without breaking what already exists? Extensibility is the
attribute of room to evolve — a stable, minimal core that every file shares, with
defined room around it for authors to add what their own subjects need and for the
format to version forward without stranding files written against earlier drafts.
It is a deliberate bet on growing through users rather than spec churn, and it
trades the simplicity of a closed format for the cost of having to define how the
open edges behave.

### Usability

Can a reader find and follow what they need? Where *Clarity* asks whether a rule
has one meaning, Usability asks whether a reader can locate that rule and learn
the document as a whole — navigability and readability rather than per-rule
precision. The two are independent: a document can be locally exact yet so poorly
ordered or unscannable that a newcomer cannot assemble the picture. A complete,
precise specification no one can find their way through still fails its readers.
It owns the document's navigability and the onboarding examples a reader learns
from — the copy-and-adapt specimens, distinct from the valid/invalid
counter-examples *Verifiability* uses to pin rule boundaries.

### Approachability

Does the front door bring a newcomer in? A project's internals can be impeccable
and still go unread if its entry point — the first and often only thing a newcomer
sees — does not land. Approachability is the attribute of that entry point: how
readily a first-time reader can grasp what the thing is and who it's for, see it
work on a concrete example, reach a first result of their own, and trust that what
they were shown is real. Where *Usability* asks whether a reader can navigate the
document they're already in, Approachability asks whether the entry point earns
them the next five minutes at all. It is the narrowest in scope — the project's
front door, not its internals — but that door has to do several jobs at once
(orient the reader, show the format working, get them to a first result, and stay
honest about what exists), so it carries more requirements than its scope alone
would suggest. It is the attribute a visitor meets first.

## Known gaps

- **The CLI's runtime quality is not modeled** — reliability, performance,
  packaging, and test coverage are deferred while the implementation is
  nascent. The one command built today, `check`, is an off-spec placeholder that
  predates the current spec and implements none of its surface (`init`, `lint`,
  and the `model` / `evaluation` / `result` resources); the spec has since
  settled `check` away entirely, gating on `lint` and `evaluation report
  --fail-on` instead. So *the CLI implementation conforms to its specification*
  is expected to rate Unacceptable until that surface lands — a not-yet-built
  gap, not spec drift the model failed to catch. These come into scope as the CLI
  matures.
- **The skills are specified but not yet built, so only their spec is modeled.**
  *The skill specification is complete* holds `specs/skills.md` to functional
  completeness, but no requirement yet judges a skill *implementation* against it
  — there is none to judge. A skills-conformance requirement joins the model once
  the skills exist, mirroring the CLI's.
- **Quality attributes beyond completeness are scoped to the format spec.**
  Clarity, consistency, verifiability, and usability are assessed only against
  `SPECIFICATION.md`; the CLI and skill specifications (`specs/cli.md`,
  `specs/skills.md`) are held to functional completeness alone. This is
  deliberate — the format spec is the contract every reader, author,
  implementation, and file depends on, so it carries the most weight — not an
  oversight; the other specs come under those attributes as they earn more of
  that weight.
- **No structural self-lint requirement yet.** Once `qualitymd lint` ships, add a
  requirement under *Functionality* whose `prompt` asks that `qualitymd lint`
  report no errors, so this file is held to the structural floor it describes —
  mirroring the built-in meta-model.
- **Implementation conformance is judged against the umbrella `specs/cli.md`**,
  which references the per-document specs (`cli-init.md`, `cli-lint.md`,
  `cli-evaluate.md`, `cli-federation.md`, and `skills.md`)
  rather than restating them; deep per-command and per-skill conformance is not
  yet a separate requirement.
