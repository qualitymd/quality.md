# QUALITY.md Specification

`QUALITY.md` is a plain-text **quality model**: a file that declares the quality
requirements for a software system or component and records how it scores against
them. It pairs **YAML frontmatter** — the structured model — with a **Markdown
body** that explains the model's context and rationale. One file gives people and
agents a shared, persistent, machine-readable account of what *good* means for a
subject and how that judgment is reached, so the same standard can be read in
review, enforced in CI, and revisited as the subject evolves.

This specification defines the format: the shape of the frontmatter and what each
part means, the evaluation semantics a tool must honor to produce comparable
results, and the conventions of the Markdown body. Its purpose is interoperability —
so that independent files and the tools that read them agree on what a `QUALITY.md`
says and how it is judged.

The specification is organized around the format's concepts. After the
**Conformance** section — which fixes the reading conventions and gathers the
mechanically checkable rules — the **Frontmatter** sections introduce the model one
concept at a time, working from the whole file inward: the quality model, the
recursive target, the factor that lenses a target, the single assessable
requirement, and the scale a requirement's result is rated against. The **Markdown
Body** and **Extensibility and Versioning** sections close.

## Conformance

Throughout, **must** marks a hard rule, **should** a recommendation that may be
departed from with good reason, and **may** an option. This specification is
normative in its entirety: a passage that states the meaning of the format binds
even when it uses none of those keywords. Only three things carry no conformance
force, each marked as such — **examples**, which illustrate rules stated elsewhere
and add none; passages labeled *Non-normative*; and latitude granted with **may**.

Conformance has two subjects. A conforming **file** satisfies the file rules
(F-rules); a conforming **reader** — any validator, evaluator, or tool that consumes
a file — satisfies the reader rules (R-rules) and the evaluation semantics. The
F-rules and R-rules are the **mechanically checkable** floor: rules a tool can
decide by inspecting a file or its own behavior, without judgment. They are the
floor, not the ceiling — a tool may layer more on top (warnings, house style, a
stricter local policy, even errors) without ever contradicting them by treating a
conforming file as malformed or accepting what the F-rules forbid. The passages that
fix the *meaning* of the model and the *behavior* of an evaluator are equally
binding, but not reducible to a mechanical check.

Each F-rule is stated in context, in the concept section it governs, under a stable
label. This section indexes them and states the two reader obligations that belong
to no single concept.

| Rule   | What it requires                                                 | Stated under  |
| ------ | ---------------------------------------------------------------- | ------------- |
| **F1** | Frontmatter is a fenced YAML block whose content is the Model    | Quality Model |
| **F2** | An optional Markdown body may follow the closing fence           | Quality Model |
| **F3** | A Target is a mapping with recognized, optional keys             | Target        |
| **F4** | Target key types                                                 | Target        |
| **F5** | A requirement declares exactly one non-empty `assessment`        | Requirement   |
| **F6** | A requirement's secondary `factors` resolve in visible scope     | Requirement   |
| **F7** | A requirement's `ratings` keys name defined scale levels         | Requirement   |
| **F8** | The Model declares a `ratings` scale of ≥2 levels, best to worst | Rating Scale  |

**Reader obligations**

- **R1.** A conforming reader preserves and ignores unknown frontmatter keys and
  unknown body sections; it must not reject a file for containing them.
- **R2.** Malformed *recognized* content is an error, not an extension — a wrong key
  type, duplicate sibling keys, duplicate scale level names, a secondary factor
  outside visible scope, or an empty assessment. A reader may warn when an unknown
  key looks like a typo of a recognized one.

## Frontmatter

The file begins with a fenced YAML frontmatter block (`---` … `---`) whose content
is the structured quality model, and an optional Markdown body may follow the
closing fence. The sections below introduce the model from the whole file inward:
the **Quality Model** (the frontmatter as a whole), the **Target** (its recursive
backbone), the **Factor** (a lens on a target), the **Requirement** (the single
assessable unit), and the **Rating Scale** (what a requirement's result is rated
against).

### Quality Model

- **F1.** A conforming file begins with a fenced YAML frontmatter block
  (`---` … `---`) whose content is a single mapping: the **Model**. A file with no
  frontmatter or no closing fence is not a conforming `QUALITY.md`.
- **F2.** An optional Markdown body may follow the closing fence.

The **Model** is the frontmatter's root mapping, and it is two things at once: the
**apex Target** of the model — so every Target key (next section) applies to it —
and the carrier of the one key that belongs to the file as a whole rather than to
any single target, `ratings` (the [Rating Scale](#rating-scale), F8). `ratings` is
the only key that distinguishes the root from any other target.

```yaml
# Model = Target + model-level keys. Appears once, as the root mapping.
ratings: <RatingScale>            # required (F8); the scale shared by all requirements
# ...plus every Target key, applied to the apex target.
```

The whole frontmatter, then, is a tree of Targets sharing a single Rating Scale at
its root. The named types are **Target**, **Factor**, **Requirement**, and
**RatingScale**, each defined in its own section below. The Target keys are optional,
so a Model carrying only `ratings` is structurally valid but assesses nothing; the
apex target should lead to at least one requirement somewhere in its subtree (see
[Target](#target)).

### Target

A **Target** is a thing evaluated, bound to the material it is assessed from by
`source`. It is the recursive node of the model: the apex Target is the Model
itself, and every child under `targets:` is another Target of the same shape, nested
to any depth. Target names are open, user-chosen identifiers such as `source-code`,
`payment-flows`, or `documentation`.

- **F3.** A Target is a mapping. Its recognized keys are `source`, `requirements`,
  `factors`, and `targets`; the Model adds `ratings` on the root mapping and only
  there. All are optional and a target declares only what it adds, except that
  `ratings` is required on the Model (F8). Target and factor names are an open,
  case-sensitive vocabulary.
- **F4.** Key types: `source` is a string; `requirements` is a map of statement →
  Requirement; `factors` is a map of name → Factor; `targets` is a map of name →
  Target (recursively the same shape).

```yaml
# Target — the recursive node type.
source: <string>                    # material this target is assessed from
requirements:                       # map of statement -> Requirement
  <requirement-statement>: <Requirement>
factors:                            # map of name -> Factor
  <factor-name>: <Factor>
targets:                            # map of name -> Target (recursively)
  <target-name>: <Target>
```

**Source.** `source` identifies the material evaluated for a target — a single
string, conventionally a path, a glob, or a URL, though a reader may support
whatever forms suit it. Paths and globs resolve relative to the containing
`QUALITY.md` file. When `source` is omitted, it defaults to the file's own directory
and all subdirectories, recursively; a grouping target may leave it implicit and let
children narrow it.

**Lineage and inheritance.** Position in `targets` is lineage: a child inherits
every applicable declaration from its ancestors. Containment is the only inheritance
primitive — a target owns what it declares and inherits its ancestors' applicable
factors and requirements. Inheritance is purely additive: a descendant may add
factors and requirements but cannot remove an inherited one. To assess a concern at
a finer grain, declare a requirement at that lower target — declaration altitude is
assessment altitude. (How a single requirement is assessed once yet governs every
descendant is fixed under [Requirement](#requirement).)

A target should lead to at least one requirement — its own, one carried by a factor,
or one contributed by a descendant. A target whose subtree holds none assesses
nothing, though a pure grouping target stays meaningful as long as its descendants
carry requirements.

### Factor

A **Factor** is a quality lens — such as `reliability` or `maintainability` —
declared on a Target and scoped to that target's subtree: declared at the apex it is
project-wide; declared on `targets.docs` it applies only to that target and its
descendants. A factor carries its own requirements, which are assessed through that
lens.

```yaml
# Factor
description: <string>               # recommended
requirements:                       # map of statement -> Requirement
  <requirement-statement>: <Requirement>
```

Factor identity is local to its scope — the same name on two unrelated targets
denotes two distinct factors. Within a scope a descendant may *refine* an inherited
factor by adding requirements under the same name; it should not *redefine* it with a
contradictory meaning. Because "contradictory" is a judgment, this is guidance a tool
may warn on rather than a deterministic rule.

A factor should carry at least one requirement — its own, one added by a refinement,
or one that names it as a secondary factor (F6). A factor that nothing contributes to
is a lens over nothing.

A factor should have a `description` of the attribute itself: what it means here, why
it matters, and how it differs from its siblings. The description should not merely
restate the requirements attached to it.

### Requirement

A **Requirement** is an assessable expectation — the single unit the model is built
to judge. Its `assessment` is performed against the target's `source`, producing a
**finding**; the active rating criteria are then applied to that finding to produce a
**result**, whose recorded value is a **rating**. Each requirement produces exactly
one result, recorded against the target that declares it, and the requirement
statement (the map key) is its identity in reports and results.

A requirement may sit directly under a target, where it is **unlensed**, or under a
factor, where it is assessed through that lens and joins that factor's rollup. It is
assessed against the target's `source` either way.

- **F5.** A requirement entry declares exactly one `assessment`: a non-empty scalar.
  A missing, empty, or list-valued `assessment` is invalid.
- **F6.** A requirement's optional secondary `factors` list names factors in scope;
  each name must resolve to a factor declared on the same target or one of its
  ancestors.
- **F7.** A requirement's optional `ratings` map is keyed by level names; every key
  must name a level defined by the Model's scale.

```yaml
# Requirement
assessment: <string>                # required; single non-empty scalar
factors:                            # optional secondary factor names in scope (F6)
  - <factor-name>
ratings:                            # optional per-requirement criteria (F7)
  <level-name>: <criterion>         #   keyed by a level name of the active scale
```

`assessment` is the instruction that produces a finding — inline criteria text or a
path to a document of criteria. It is never a list (F5); if one statement needs
several independent assessments, split it into several requirements.

The optional secondary `factors` list (F6) lets one result appear in additional
factor views without changing the requirement's primary placement, so one result can
contribute to several factor rollups without repeating the requirement. The optional
`ratings` map sets this requirement's own criteria (F7); it changes only the
criteria, not the levels, their order, or their display names. See
[Custom rating criteria](#custom-rating-criteria) for when and how to use it.

**Assessed once.** A requirement is **assessed once, at the target that declares
it**, against that target's `source` — one assessment, one finding, one result.
Containment then makes it *govern* every descendant: the requirement joins their
inherited context and its single result covers their subtree, but it is never
re-assessed against a descendant's narrower source and never yields a second result.
Because a target's `source` ordinarily spans its descendants', that one assessment
routinely inspects artifacts that also belong to sub-targets; that is expected — a
finding may cite material anywhere in the declaring target's source, including files
a descendant also selects, without splitting into multiple results.

These are **evaluator obligations**: they bind any reader that evaluates a model,
even though no file can be checked against them. An evaluator records and rolls up
results, while each individual assessment is performed by a judge — a person or a
model. Tools that disagree here produce results that cannot be compared.

**Invalid examples.** Each illustrates the rule it violates and adds none of its own.

```yaml
requirements:
  "input is validated":
    assessment:
      - "Query parameters are validated." # invalid (F5): assessment must be one scalar
      - "Request bodies are validated."
```

```yaml
targets:
  docs:
    factors:
      clarity:
        description: Documentation is understandable.
  api:
    requirements:
      "errors are readable":
        assessment: Errors explain what happened.
        factors:
          - clarity # invalid (F6): sibling target's factor is out of scope
```

```yaml
ratings:
  - { level: target, criterion: "Meets the requirement." }
  - { level: unacceptable, criterion: "Does not meet the requirement." }
requirements:
  "tests cover critical paths":
    assessment: Critical paths are covered.
    ratings:
      gold: "Exceptional coverage." # invalid (F7): `gold` is not a scale level
```

**A worked example.** The pieces together — a target with a direct requirement, two
scoped factors each carrying a requirement, and one requirement that names a
secondary factor:

```yaml
---
ratings:
  - { level: outstanding,  criterion: "Exceeds the requirement with margin to spare." }
  - { level: target,       criterion: "Satisfies the requirement." }
  - { level: minimum,      criterion: "Falls short of the goal but holds the floor." }
  - { level: unacceptable, criterion: "Falls below the acceptable floor." }
targets:
  api:
    source: ./internal/api
    requirements:
      "the API preserves accepted orders":
        assessment: >
          Accepted orders are durably stored before a success is returned; failures
          surface as errors rather than false success responses.
    factors:
      reliability:
        description: >
          The API behaves predictably under ordinary and failure conditions.
        requirements:
          "the write path is covered end-to-end":
            assessment: >
              Automated tests exercise a successful write, a storage failure, and a
              retry path, and would fail if an acknowledged order could be lost.
      security:
        description: >
          Customer data and privileged operations are protected from unauthorized use.
        requirements:
          "no secrets are committed":
            assessment: >
              No credentials, API keys, private keys, or tokens appear in source,
              config, or fixtures; secrets are loaded from the runtime environment or
              a secrets manager.
            factors:
              - reliability
---
```

### Rating Scale

A **RatingScale** is the shared yardstick every requirement's result is rated
against. The Model declares it once, at the root, in `ratings`.

- **F8.** The Model must declare `ratings`: a non-empty scale of at least two
  levels, listed best to worst, given inline. Each level has a unique `level` name
  and a `criterion` the evaluator judges against. A Model with no `ratings`, or a
  scale of fewer than two levels, is invalid. The format prescribes no default
  scale; an author chooses one to fit the subject.

```yaml
# RatingScale = list of RatingLevel, best to worst (F8)
- level: <level-name>               # required; unique within the scale; position is rank
  title: <string>                   # optional human label
  criterion: <string>               # required; the criterion the evaluator judges against
```

`ratings` is an inline sequence of levels ordered best to worst, where position
defines rank. The evaluator applies the criteria top-down and records the best level
whose criterion the finding satisfies.

A scale need not be elaborate. The minimum is two levels — a plain pass/fail gate:

```yaml
ratings:
  - { level: pass, criterion: "Satisfies the assessment." }
  - { level: fail, criterion: "Does not satisfy the assessment." }
```

#### A suggested scale: the landing zone

*Non-normative.* The format prescribes no default (F8), but when a graded scale fits
and you have no strong preference, this four-level scale is a reasonable starting
point, and a scaffolding tool may seed it. Its vocabulary and best-to-worst framing
are adapted from the Agile Landing Zone pattern — **outstanding** exceeds the goal,
**target** meets it, **minimum** holds the acceptable floor, and **unacceptable**
falls below it:

```yaml
ratings:
  - { level: outstanding,  title: Outstanding,  criterion: "Exceeds the requirement; satisfies it with margin to spare." }
  - { level: target,       title: Target,       criterion: "Satisfies the requirement." }
  - { level: minimum,      title: Minimum,      criterion: "Falls short of the goal but stays at the acceptable floor." }
  - { level: unacceptable, title: Unacceptable, criterion: "Falls below the acceptable floor." }
```

#### Custom rating criteria

Set `ratings` on a requirement (F7) when the scale's shared criteria cannot express
the gradient that matters. The custom criteria should name ordered, mutually distinct
levels of the active scale, and the evaluator assigns the best level met.

A measured-bound example:

```yaml
requirements:
  "evaluation completes fast enough to sit in a pre-commit hook":
    assessment: >
      Wall-clock time to evaluate a typical model, measured cold on CI hardware,
      p95 over 20 runs.
    ratings:
      outstanding: "p95 under 100 ms."
      target: "p95 under 250 ms."
      minimum: "p95 under 500 ms."
      unacceptable: "p95 at or above 500 ms, or it cannot run in the hook."
```

A judged-spectrum example:

```yaml
requirements:
  "the suite covers what matters and is pyramid-balanced":
    assessment: >
      Taken as a whole, the suite exercises the behavior that matters, with most
      checks at the unit level and fewer, broader tests above.
    ratings:
      outstanding: "Every critical path is covered; the suite is fast and unit-dominant."
      target: "Important behavior is covered and the shape is roughly pyramidal."
      minimum: "Important paths are partly covered, but notable gaps or slowness remain."
      unacceptable: "A critical behavior is untested, or the suite cannot be trusted."
```

Do not customize `ratings` merely to restate "met" and "not met"; spend that text on
the `assessment` instead.

## Markdown Body

The Markdown body documents *why* this is the right model — the context an evaluator
needs to interpret assessments consistently. It is optional (F2), and the format does
not restrict it to any fixed set of sections.

*Non-normative.* The sections below are a recommended starting point, not a required
structure; teams and tools may use, rename, or replace them, and a reader preserves
sections it does not recognize (R1).

| Section                 | What it captures                                                                                                                  |
| ----------------------- | --------------------------------------------------------------------------------------------------------------------------------- |
| **Overview**            | What the subject is, who depends on it, and what "good" means here.                                                               |
| **Scope**               | The model boundary: what it covers and deliberately leaves out. Out-of-scope concerns are exclusions by design, not deficiencies. |
| **Needs**               | Stakeholder outcomes the requirements answer to.                                                                                  |
| **Risks**               | What goes wrong, and for whom, if a need is not met.                                                                              |
| **Targets and factors** | A prose mirror of the target tree: each target's role, the scoped factors declared there, and why those lenses belong there.      |
| **Known gaps**          | In-scope quality concerns deliberately deferred, each with a brief reason.                                                        |

Applicability is structural: a factor declared on `targets.docs` applies there and
below, not to unrelated targets. The body should explain that structure rather than
argue exceptions in prose. **Scope** is for concerns outside the model's remit;
**Known gaps** is for in-scope concerns deliberately deferred.

```markdown
# Quality model - Orders platform

## Overview

The Orders platform receives, stores, and exposes customer orders. Good means an
accepted order is durable, observable, and protected from unauthorized access.

## Scope

This model covers the API service, worker, and storage layer. The external payment
provider is out of scope because this project does not own it.

## Needs

- Customers can trust that an accepted order will not disappear.
- On-call engineers can diagnose an incident from logs and metrics.

## Risks

Silent order loss is the worst outcome because customers cannot repair it
themselves and support cannot reconstruct intent reliably.

## Targets and factors

### api

The API target covers the request boundary and declares Reliability and Security.
Reliability matters here because success responses create customer commitments;
Security matters because this boundary receives customer data.

### worker

The worker target carries Operability because delayed or stuck asynchronous work is
visible first through operational signals.

## Known gaps

- Sustained peak-load behavior is in scope but not modeled yet.
```

## Extensibility and Versioning

The minimal structural core is a frontmatter Model with its required `ratings` scale
(F1, F8). Because the Target fields are optional, a Model carrying only `ratings` is
structurally valid but not useful; tools should warn when one declares no
requirements, factors, or child targets. The body is optional. The format grows
through use: a reader preserves unknown frontmatter keys and unknown body sections
(R1) rather than rejecting them, while malformed recognized content is an error, not
an extension (R2).

### Extending the format

*Non-normative.* Target and factor names are an open vocabulary, so producers may
extend freely: add frontmatter keys and body sections for your own tooling, kept
additive so other readers can ignore them (R1). A tool may enforce a stricter local
policy that rejects what the format merely discourages, provided it never treats a
conforming file as malformed (Conformance). A scaffolder may seed opinionated
defaults — the landing-zone scale, the recommended body sections — as a starting
point an author is free to replace.

### Edge cases

How the rules above resolve specific cases; each row applies a rule already stated:

| Case                                                        | Treatment                                                                                                                             |
| ----------------------------------------------------------- | ------------------------------------------------------------------------------------------------------------------------------------- |
| **No frontmatter, or no closing `---`**                     | Invalid (F1): there is no model to read.                                                                                              |
| **Empty `targets`**                                         | Valid but usually a warning; it declares no child targets.                                                                            |
| **Empty Target**                                            | Valid as a grouping placeholder; a tool may warn if it contributes nothing.                                                           |
| **Factor refinement that changes description incompatibly** | Discouraged, not invalid; a tool may warn. Adding requirements or compatible detail is fine.                                          |
| **Secondary `factors` entry outside visible scope**         | Invalid (F6): the factor name must resolve on the current target or an ancestor.                                                      |
| **Empty assessment value**                                  | Invalid (F5): the assessment is the criteria and must be non-empty.                                                                   |
| **Duplicate sibling names**                                 | Invalid (R2): YAML duplicate keys and duplicate `ratings[].level` names are rejected structurally.                                    |
| **Ordering**                                                | Significant within `ratings` (F8); otherwise evaluation does not depend on authoring order, though tools may preserve it for display. |
| **Case**                                                    | Schema keys and author-defined names are case-sensitive.                                                                              |

### Versioning

The format evolves additively. New optional keys and sections may be introduced
without invalidating existing conforming files. A change that would invalidate an
existing conforming file is breaking and reserved for a new major format version.
