---
type: How-to Guide
title: Modeling quality across domains
description: Why QUALITY.md's domain range is wide and how to keep examples agnostic — the society/sphere/quality-context model, the stress axes that make a domain a useful test, a role-tagged catalog of quality contexts, and a worked non-software example.
tags: [doctrine, domain-agnostic, examples, contributing]
timestamp: 2026-06-24T00:00:00Z
---

# Modeling quality across domains

## Why the range is wide

QUALITY.md is **quality-domain agnostic**: a Model can describe quality for
software, documentation, data sets, research, services, operations, processes, or
other evaluated entities. That breadth is not a neutral fact about a schema — it
follows from what the format is finally _for_.

The format exists to serve care: to help people tend and keep what is entrusted to
them in their work, and through that work, to serve the society their work holds up.
Care is wide, so the range is wide. It runs through every sphere of human life — the
learned professions, the trades, the arts, operations, teaching, research — and it
does not stop at paid work: the household is the first place most people tend
something in their charge, and the format means to serve that care no less than the
professional kind. [Why care, not conformance](#why-care-not-conformance) develops
the claim; the rest of the guide puts it to work.

Domain agnostic is not context neutral. QUALITY.md is agnostic about _what_ a Model
evaluates and opinionated about _how_ it is used — the agent- and skill-first
workflow described under [Agentic use context](#agentic-use-context).

## What this guide is for

This guide is for contributors writing example content, docs, specs, or skill
guidance. The agnosticism claim is easy to _assert_ and easy to leave _untested_ — a
claim exercised against one domain is no claim at all. The guide explains why a
domain makes a good stress test, gives a catalog of quality contexts to draw on,
fixes which to reach for, and demonstrates the claim with a worked non-software
example.

It is also the **authoritative home** for the repo-content rules that keep examples
domain agnostic and preserve the agentic use context — see
[Rules for domain-agnostic example content](#rules-for-domain-agnostic-example-content)
and [Agentic use context](#agentic-use-context). [`AGENTS.md`](../../AGENTS.md)
carries a short summary of each and links here; when the two diverge, this guide
governs.

## Why care, not conformance

To _work_ and to _keep_ is the double charge — "The Lord God took the man and put
him in the Garden of Eden to work it and take care of it" (Genesis 2:15). The
_keeping_ half is the quieter one, and it is the half a quality model serves. To
keep a garden is to tend it before it is overgrown — to meet what is small while it
is still small. The same counsel runs through the Tao Te Ching — "Magnify what is
small, increase what is meager"; "deal with great affairs while they are still small
matters"; "prepare for hardships while it is still easy to do so" — for "what has
not yet happened is easily prevented" (_Tao Te Ching_ 63–64, trans. David Bentley
Hart). A `QUALITY.md` is an instrument of that keeping: not a gate that catches
failure at the end, but a way of attending continuously, so concerns are seen and
met while they are still small.

And what is kept is tended by care, not conformance. Seeing why the occupations that
most formalize such care belong to one family explains where quality finally lives.

### One family of professions

Ministry, law, medicine, and engineering have long been treated together. Donald
Schön groups engineering with the "learned professions" of medicine and law under a
single epistemology of practice — knowing- and reflecting-in-action (_The Reflective
Practitioner_, 1983) — and the older professions tradition binds the same set by a
service whose "primary obligation is to the public good"
([J. R. Wilcox, _New Catholic Encyclopedia_](https://www.encyclopedia.com/religion/encyclopedias-almanacs-transcripts-and-maps/professional-ethics)).
Each tends a different entrusted good, but the shape repeats:

- **Ministry** tends the formation of persons — the "cure of souls" (_cura
  animarum_) that Gregory the Great's sixth-century _Pastoral Rule_ calls "the art
  of arts." Its excellence is discernment fitted to the person — "one and the same
  exhortation does not suit all" — never faithful repetition of a rite.
- **Medicine** tends the patient's good through a healing relationship; its
  excellence is the "right and good healing action for _this_ patient" (Pellegrino
  and Thomasma, _For the Patient's Good_, 1988) — the practical wisdom medicine calls
  _phronesis_, which guideline-concordant care can satisfy on paper and still miss.
- **Law** tends the rule of law and the administration of justice; its calling is to
  the public good, not to winning. At its best a contract sets clear, fair terms
  that order relationships toward justice. Law is not by nature adversarial:
  adversarial advocacy is one contingent means to that end, and other traditions reach it without (the civil-law _notaire_ advises
  all parties impartially). The lawyer is an officer of the court and fiduciary, not
  merely a hired advocate, and legal quality is fidelity to the law's purposes (Lon
  Fuller, _The Morality of Law_, 1964).
- **Engineering** tends a public it never meets; its codes require holding
  "paramount the safety, health, and welfare of the public" (NSPE Canon 1), and
  because engineering runs ahead of complete science, quality is judgment about
  failure modes — a design can meet every code and still fail the people who depend
  on it.

Read this way, an engineering `QUALITY.md` is not a compliance artifact within a
value-neutral science but the means by which one of these professions tends what is
in its care.

### Care under-determines conformance

The same lineage explains _why_ care, not conformance, is where quality lives. In
the practice tradition running from Heidegger's account of care (_Sorge_) through
Hubert Dreyfus's model of skill acquisition, expertise is _constituted_ by
involvement: the detached rule-follower is the novice, and mastery is the
situational coping that rules can only approximate. A `QUALITY.md` read as a
checklist captures only the beginner's stage of a craft — the requirements name the
rules, but the quality they point at is the involved care of someone who has made
the practice their own.

The religious traditions say the same in sharper terms. The same deed counts for
everything or nothing depending on the care behind it: eloquence, knowledge, even
giving away all one owns are "nothing" without love (1 Cor 13). The goal is "be
transformed, not conformed" (Rom 12). And the caution lands closest of all on a
written format like this one — "the letter kills, but the Spirit gives life" (2 Cor
3): a `QUALITY.md` held as literal requirements to satisfy becomes the very deadness
these texts condemn, while the same document held as living care orients
transformation. Quality here is known by _fruit_ over time, not surface conformance
at an instant. That the good is known by discernment — told from its lack — is not
confined to one lineage: "All who dwell under Heaven recognize beauty ... by
distinguishing it from deformity; all recognize the good by distinguishing it from
wickedness" (_Tao Te Ching_ 2).

Contemporary work sharpens the consequence for a format. B. Scot Rousse and David
Spivak describe care as _tending and attending_ — "the ongoing activity of
cultivation that brings things out at their best" — and stress that "you can't fully
specify in advance what caring for something will require"
([Notes on Care](https://withoutwhy.substack.com/p/notes-on-care-with-david-spivak-a)).
That is the deepest reason a `QUALITY.md` under-determines quality: not a gap to
close with more requirements, but the nature of care, which a written artifact can
point toward and never exhaust. Their caution also lands on this project's own use
context — they worry about systems that "process our values but do not care." A
Model handed to an agent or a team can be _processed_ without being _cared_ through;
the requirements read the same either way, and only the second is quality.

## Three axes that do not line up

It is tempting to collapse "the work," "the kind of quality in play," and "the
lenses a model uses" into one idea called a _domain_. Pulling them apart gives a
simple picture of where quality lives:

**society → spheres → quality contexts → a Model → factors**

- A **society** is held up by its **spheres** — family, the trades, commerce and
  industry, the academy, civic life, the arts, and the rest. A sphere is a domain of
  responsibility with an entrusted good of its own; it carries the _motivation and
  stakes_ — why the quality matters, and to whom.
- Within and across spheres, work happens in **quality contexts** — recurring
  _kinds_ of quality concern, or contexts of use: requirements quality, data
  quality, documentation, budgeting, caregiving, the quality of an analysis or a
  decision. This axis carries the _shape_ of the concern.
- A **Model** is a single `QUALITY.md` for a single entity, and it earns its
  **factors** — the specific quality lenses — from that entity's own risks and needs
  (see [the authoring guide](../../skills/quality/guides/authoring.md)).

The nesting orients, but it is not a tree, and the axes do not line up:

- A single quality context **recurs across spheres**. Budgeting is a quality context
  in the family (the household budget) and in commerce (corporate planning); data
  quality concerns the scientist, the clinician, and the logistics planner alike.
  This is why the [catalog below](#a-catalog-of-quality-contexts) is organized by
  context, not filed under spheres.
- Some familiar groupings are **cross-cutting, not spheres.** The _learned
  professions_ (medicine, law, ministry, engineering) thread through several spheres,
  bound by service to a public good ([above](#why-care-not-conformance)). _Knowledge
  work_ — cognitive, informational work — likewise cuts across spheres, and it is
  the attribute most tied to the [agentic use context](#agentic-use-context): an
  assistant is most often in the loop on knowledge-work contexts, whatever the
  sphere.
- **Factors are earned per Model**, never assigned by sphere or read off a context.

The lists in this guide are illustrative, non-exhaustive, and overlapping — not
taxonomies to fill in. That holds throughout; the
[rules below](#rules-for-domain-agnostic-example-content) make it normative for repo
content.

The quality-context axis is the one with the most external formalization, and that
is worth leaning on when _naming_ a context. Decades of measurement and evaluation
practice — the lineage QUALITY.md draws on (see [`SPECIFICATION.md`](../../SPECIFICATION.md),
Lineage) — grew up around specific contexts: product quality, data quality,
requirements quality, process quality, and service quality, among others. That each
earned its own body of practice is the best evidence that _quality context_ is a
real, recurring axis. Two cautions keep the borrowing honest. First, use that
lineage to bound the _context_, never to populate the _factors_: a standard's
characteristic list (reliability, security, maintainability, and the rest) is
factor-axis content and stays earned per entity, not adopted as defaults — importing
it wholesale is the exact anti-pattern the
[rules below](#rules-for-domain-agnostic-example-content) forbid. Second, those
models lean conformance-ward; QUALITY.md takes the boundaries and the vocabulary,
not the checklist stance.

## What makes a domain a useful test

Software is the project's primary illustrative domain — most contributors think in
it fluently — but it is _unrepresentative_ in ways that quietly bias the format's
abstractions. In a software Model the `source` is a readable code path, the
`assessment` can often lean on a runnable check, the audience is a roughly aligned
team, and "good" usually means correctness. A format tuned only against that domain
can develop hidden assumptions — that a source is always a file an agent can read,
that an assessment always has an executable oracle — without anyone noticing,
because the one domain in view never contradicts them.

A _secondary_ domain earns its place by contradicting one of those assumptions — by
pressing the abstractions along an axis where software does not push back. A few axes
matter most (illustrative and overlapping, not a closed taxonomy):

- **Source materiality.** Is the evaluated thing a readable artifact (code, prose),
  a data set, a _live process_ with no static artifact, or a personal system?
  Software sits at the readable-artifact end; a service has no file to read.
- **Assessment oracle.** Can a requirement be checked by something executable, only
  by expert _judgment_ against a rubric, or only against _values or tradition_? A
  second cut runs underneath: _internal_ quality (static properties of the artifact),
  _external_ quality (its behavior in operation), and _quality in use_ (whether
  someone accomplishes their goal with it in a real context) often need different
  oracles. Software's runnable check usually covers internal and external and stops
  short of quality in use — so a domain that lives in the in-use view stresses the
  oracle in a way software does not.
- **Constituency.** Does quality serve a single self, a roughly aligned team, or
  _multiple parties_ whose immediate interests differ but who answer to a shared
  good? A contract serves several parties under the law's public purpose; a personal
  system answers to one person.
- **Stakes.** Is the dominant stake correctness, fitness-for-purpose, persuasion, or
  _meaning_? Software stakes are usually correctness; a liturgy's stakes are meaning.

The skill already carries the machinery for _authoring_ in any domain once it is
named: the [authoring guide](../../skills/quality/guides/authoring.md) enumerates a
composite entity's constituents from stewardship concerns × audience/purpose, so the
same generator works in any domain — lean on it rather than inventing a parallel
taxonomy.

## A catalog of quality contexts

When an example needs a non-software illustration, draw from the catalog below. It is
organized by _quality context_, because a context recurs across spheres rather than
belonging to one — each entry notes the spheres it **recurs across** to show that.
Everything here is illustrative, non-exhaustive, and overlapping (as throughout); the
factor sketches are _not_ default factor sets — a real Model earns its factors (see
[the authoring guide](../../skills/quality/guides/authoring.md)). Two roles are
tagged: **cite-worthy** contexts are the ones to reach for in worked or substantial
examples; **range-finders** probe the edges of the format's reach.

### Cite-worthy — reach for these

These four keep a readable source but have no runnable oracle, so quality rests on
judgment. That is not a step down from software but a clarification of it: even a
software requirement's executable check is a _proxy_ for the judgment that finally
decides whether the thing is good. These remove the proxy and leave the judgment in
plain view. The [worked example below](#worked-example-a-documentation-set) develops
the first.

| Quality context                                 | Stresses                                | Recurs across                                         | Illustrative factor sketch                                      |
| ----------------------------------------------- | --------------------------------------- | ----------------------------------------------------- | --------------------------------------------------------------- |
| **Documentation / written corpus** _(flagship)_ | judgment oracle; prose source           | the professions, commerce, the academy, the household | accuracy, completeness, findability, clarity, currentness       |
| **Data set / data product**                     | source is data, not a code path         | science, commerce, civic and public life              | accuracy, completeness, provenance, timeliness, fitness for use |
| **Research / analytical report**                | credibility without an executable check | the academy, commerce, the professions                | rigor, reproducibility, sourcing, calibration, clarity          |
| **Service / operation**                         | no static artifact at all               | commerce, the professions, civic life                 | resolution quality, consistency, recoverability, responsiveness |

### Range-finders — probe the edges

These push constituency and stakes hardest. They are reasoning aids, not worked
models, and the project takes no position on any tradition, legal stance, or personal
arrangement:

- **Legal / contract.** Quality serves _multiple parties under a shared good_ — the
  law's public purpose — so a Model cannot reduce "good" to one party's preference;
  it names the good a fair ordering serves. Law is not by nature adversarial; legal
  quality is fidelity to the law's purposes, not winning. (Developed under
  [Why care, not conformance](#why-care-not-conformance).) Recurs
  across civic and commercial life.
- **Finance / budgeting.** A hard, objective floor (the books must balance; the plan
  must be solvent) sits _over_ value-laden priorities (what is worth spending on) —
  an executable-style constraint paired with judgment about fit-to-goals. Recurs
  across the household, commerce, and civic life.
- **Personal productivity.** The evaluated entity is a _system of one_ — a personal
  workflow — with no external audience and a legitimately self-defined bar. It tests
  whether the format works when the only constituent is the author. Recurs across any
  sphere, lived individually.
- **Devotional / religious practice.** Quality is defined by a tradition's account of
  the good and resists every reduction to checkable conformance; it is known by
  _fruit_ over time, not surface conformance at an instant. This is not a distant
  domain from engineering quality — it clarifies the center of all of them.
  (Developed under [Why care, not conformance](#why-care-not-conformance).)

### And across every sphere

The catalog is deliberately short; quality is tended far more widely than it samples.
A non-exhaustive sweep: in **commerce and industry**, the quality of sales research,
marketing content, a product spec, a support resolution, an analytics dashboard; in
**the academy**, a literature search, a study design, a reproducible analysis; in the
**trades and arts**, a joiner's fit and finish, an editor's ear; in **teaching and
care work**, a lesson, a clinical handoff, a care plan; and in the **household** —
the first place most people tend something in their charge — a budget, a meal plan, a
caregiving routine, home upkeep met while it is still small, the family calendar and
records, the formation of children. Care reaches from the most formalized profession
to the quietest thing kept at home.

## Worked example: a documentation set

Asserting agnosticism is cheap; the README and SPECIFICATION already show the model
_shape_ with a software service. What follows shows it holding in a different domain
— a complete documentation / knowledge-base Model. Two things make it a real stress
test rather than the software example relabeled: the `source` is a prose corpus, and
every `assessment` describes a human-judgment check, with no runnable oracle. It also
runs all three oracle views — internal (the page is well-formed), external (its
content is accurate to the shipping product), and quality in use (a reader of the
stated audience reaches the right page and follows it).

```markdown
---
title: Helios Platform Docs
ratingScale:
  - level: outstanding
    title: 🟢 Outstanding
    description: The documentation clearly exceeds the shared quality bar.
    criterion: "Consistently exceeds the requirement with clear margin."
  - level: target
    title: 🔵 Target
    description: The documentation meets the shared quality bar.
    criterion: "Meets the expected quality bar."
  - level: minimum
    title: 🟡 Minimum
    description: The documentation is usable, but has gaps worth improving.
    criterion: "Meets the lowest acceptable bar, with visible gaps."
  - level: unacceptable
    title: 🔴 Unacceptable
    description: The documentation is below the shared quality bar.
    criterion: "Falls below the minimum acceptable bar."
areas:
  reference:
    title: API Reference
    source: ./docs/reference
    factors:
      accuracy:
        title: Accuracy
        description: Documented behavior matches how the product actually works.
        requirements:
          documented-parameters-match-the-current-api:
            title: documented parameters match the current API
            assessment: >
              A reviewer checks each documented endpoint's parameters, types, and
              defaults against the current API and treats any divergence as a
              defect, paying closest attention to recently changed endpoints.
      findability:
        title: Findability
        description: A reader can reach the right page from a realistic starting point.
        requirements:
          common-tasks-are-reachable-in-a-few-steps:
            title: common tasks are reachable in a few steps
            assessment: >
              For a sample of common reader tasks, a reviewer starts from the docs
              entry point and confirms the relevant page is reachable through
              navigation or search within a few steps, without already knowing the
              page's title.
  guides:
    title: How-to Guides
    source: ./docs/guides
    factors:
      completeness:
        title: Completeness
        description: The guides cover the tasks readers actually need to accomplish.
        requirements:
          the-top-supported-tasks-each-have-a-guide:
            title: the top supported tasks each have a guide
            assessment: >
              A reviewer compares the guide set against the product's top supported
              tasks and confirms each has a guide, recording any task with no guide
              as a gap.
      clarity:
        title: Clarity
        description: A reader of the stated audience can follow a guide without help.
        requirements:
          steps-are-followable-by-the-stated-audience:
            title: steps are followable by the stated audience
            assessment: >
              A reviewer reads each guide as a member of its stated audience and
              confirms the steps can be followed in order without unstated
              prerequisites or undefined terms.
---

# Quality model: Helios Platform Docs

## Overview

This model describes the quality bar for the Helios platform documentation. Good
documentation here is accurate to the shipping product, lets a reader find the
right page quickly, covers the tasks people actually do, and can be followed by its
stated audience. None of these is settled by a runnable check; each is a judgment a
reviewer makes against the corpus.

## Scope

This model covers the reference and how-to documentation under `./docs`. It does
not cover marketing pages, the product UI's in-app copy, or the source code the
documentation describes.
```

Notice what changed and what did not. The `ratingScale`, the area/factor/requirement
structure, and the body sections are identical in shape — the format needed no
documentation-specific dialect. What changed is domain-carried: the factor family
(accuracy, findability, completeness, clarity) and the `assessment` strings, which
describe a reviewer's judgment rather than a test run. That is the demonstration: the
same Model carries a domain whose quality lives entirely in judgment.

## Rules for domain-agnostic example content

These rules are normative for repo content. [`AGENTS.md`](../../AGENTS.md) carries a
one-line summary and links here; this section governs.

Concrete quality model content in this repo is **illustrative** unless it defines
this project's own Model or states a normative format rule. That covers example
Areas, Factors, Requirements, Assessments, criteria, Rating Levels, Findings,
recommendations, and any quality-domain example. When in doubt, treat content as
illustrative and mark it so.

- **Pair the primary with a secondary.** For worked or substantial example content,
  anchor in software/product quality — the familiar domain contributors read
  fluently — and pair it with one **cite-worthy** secondary context (above), so the
  example _demonstrates invariance_: the same model shape, with only domain-carried
  content changing. One anchor plus one contrast is the sweet spot; avoid stacking
  multiple software cases, and keep the pairing balanced so software does not read as
  the privileged "real" example with the secondary as an afterthought. Brief inline
  illustrations need not carry two full domains — lead with the principle and cite a
  context or two non-exhaustively.
- **Mark illustrative status.** Make clear an example illustrates one domain and is
  not a default — unless it is a format rule or this project's own Model.
- **Frame lists as open.** When citing quality contexts or factor families, make
  clear the examples are brief, illustrative, overlapping, and not exhaustive.
- **Don't make software the default.** Prefer domain-neutral principles first, then
  examples from several domains when examples help. Use software quality as the
  default only for explicitly software topics; do not imply software product quality
  is the default use.
- **Keep the registers straight.** Domain-agnostic _model content_ is separate from
  the agent-first _use context_ (next section). AI-assistant and harness language is
  appropriate when describing how QUALITY.md is used or this project itself; it is
  not the default modeled domain.

## Agentic use context

These rules are normative for repo content. [`AGENTS.md`](../../AGENTS.md) carries a
one-line summary and links here; this section governs.

Domain agnostic does not mean context neutral. QUALITY.md is domain agnostic in
_what a quality model can describe_ — software, documents, data sets, research or
analytical reports, services, operations, processes, AI assistants, agent harnesses,
or other evaluated entities. This project is not context neutral in _how_ QUALITY.md
is used.

That distinction licenses explicit guidance for **use-context constituents**, not
for modeled domains. The agent harness and the QUALITY.md self-check recur from the
agentic use context, so project guidance may name their expected shapes, factor
families, and requirement patterns. A modeled domain — the thing a particular
`QUALITY.md` evaluates — still never receives a default factor roster; factors are
earned per Model from the entity's own risks and needs.

Even for a use-context constituent, keep its factors and requirements agnostic to
the **served domain**: the domain the project model is about, and that the
use-context constituent serves. For the agent harness, write requirements around
what harness governing artifacts do — orient an agent, point to how work is
verified, expose feedback, and bound permitted action — rather than around one
served domain's tools.

- Good: "Harness governing artifacts point an agent to how the project verifies
  work."
- Avoid: "Harness governing artifacts document the lint, type-check, and test
  commands."
  That assumes software; a documentation, data, service, research, or operational
  project verifies work differently.

The primary experience is agent- and skill-first. AI assistants and coding agents
read, author, evaluate, and improve `QUALITY.md`; the `/quality` skill carries
judgment; and the CLI provides deterministic support tooling. Preserve that agentic
use context in docs, specs, examples, and skill content.

- Do not remove references to AI assistants, coding agents, agent-accessible
  evidence, harnesses, skill workflows, or agent collaboration when they describe
  how QUALITY.md is used. Do not treat that operating context as the default modeled
  quality domain.
- Do not flag a phrase solely because the evaluated project is an AI assistant,
  coding agent, or harness. Those are valid project/use contexts. Flag only when the
  wording makes that domain sound inherent to QUALITY.md, normal for all QUALITY.md
  files, or the default model content.

Decision test — name which register a phrase is in before judging it:

- **Use context:** who uses QUALITY.md, through what workflow, with what tools.
  Agentic/AI language is appropriate and often preferred.
- **Model domain:** what a `QUALITY.md` evaluates. Keep this domain agnostic unless
  the example is explicitly scoped.
- **Project/use context:** the concrete project or harness being evaluated. AI
  assistant or harness language is acceptable when scoped to that project.
- **Project self-description:** this repository's own format, skill, CLI, docs, and
  examples. Agentic/AI language is appropriate when true for this project.

Examples:

- Good: "Use QUALITY.md with the `/quality` agent skill to align coding agents and
  teams."
- Good: "Record context that is agent-accessible."
- Good: "This example models an AI assistant harness."
- Good: "Evaluate and improve the quality of AI assistant projects and harnesses."
  Acceptable when describing that use case or this project's agentic tooling, not the
  universal scope of QUALITY.md.
- Good: "The CLI is support tooling for the agent-first workflow."
- Needs scoping: "QUALITY.md helps improve AI assistant quality." Better: "QUALITY.md
  can model AI assistant quality; it can also model other domains."
- Avoid: "A QUALITY.md normally evaluates a codebase or agent harness."
- Avoid: "QUALITY.md is for evaluating AI assistant projects and harnesses."
- Avoid: "Default factors include security, reliability, usability, and
  maintainability."
- Avoid: removing "agent-accessible" because it sounds AI-specific.

## Checklist

Before adding concrete quality-model example content, check:

- The example's quality domain — its sphere and context — is stated or implied.
- It is anchored in software/product and paired with one cite-worthy secondary
  context, balanced — unless it is a brief inline illustration.
- Illustrative status is clear, unless it is a format rule or this project's own
  Model.
- Brief lists are framed as non-exhaustive and possibly overlapping.
- Factors are earned for the example's entity, not imported from a standard's
  characteristic list.
- Software is the default only for explicitly software topics.
- The register is clear — model domain (agnostic) versus agentic use context
  (agent-first) — per the decision test above.
