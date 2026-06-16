# QUALITY.md Format

`QUALITY.md` is a plain-text quality model: YAML frontmatter containing the
structured model, followed by a Markdown body that explains the model's context
and rationale.

Throughout, **must** marks a hard rule, **should** a recommendation that may be
departed from with good reason, and **may** an option. This specification is
normative in its entirety: a passage that states the meaning of the format binds
even when it uses none of those keywords. Only three things carry no conformance
force, each marked as such — **examples**, which illustrate rules stated elsewhere
and add none; passages labeled *Non-normative*; and latitude granted with **may**.

Conformance has two subjects. A conforming **file** satisfies the file rules
(F-rules); a conforming **reader** — any validator, evaluator, or tool that consumes
a file — satisfies the reader rules (R-rules) and the evaluation semantics. The
**Conformance** section gathers the rules a tool can check mechanically, the floor
that lets independent validators agree; the sections after it fix the meaning of the
model and the behavior of an evaluator — equally binding, but not reducible to a
mechanical check.

## Conformance

These are the **mechanically checkable** rules — those a tool can decide by
inspecting a file or its own behavior, without judgment. They are the floor, not the
ceiling: a tool may layer more on top — warnings, house style, a stricter local
policy, even errors — and the sections below point to many such openings. What it
must not do is contradict them, treating a conforming file as malformed or accepting
what the F-rules forbid. The format prescribes the contract, not the tool; tools are
named here only by way of example.

**File — frontmatter and body**

- **F1.** A conforming file begins with a fenced YAML frontmatter block
  (`---` … `---`) whose content is a single mapping: the **Model**. A file with no
  frontmatter or no closing fence is not a conforming `QUALITY.md`.
- **F2.** An optional Markdown body may follow the closing fence.

**File — target**

- **F3.** A Target is a mapping. Its recognized keys are `source`, `requirements`,
  `factors`, and `targets`; the Model adds `ratings` on the root mapping and only
  there. All are optional and a target declares only what it adds, except that
  `ratings` is required on the Model (F8). Target and factor names are an open,
  case-sensitive vocabulary.
- **F4.** Key types: `source` is a string; `requirements` is a map of statement →
  Requirement; `factors` is a map of name → Factor; `targets` is a map of name →
  Target (recursively the same shape).

**File — requirement**

- **F5.** A requirement entry declares exactly one `assessment`: a non-empty
  scalar. A missing, empty, or list-valued `assessment` is invalid.
- **F6.** A requirement's optional secondary `factors` list names factors in
  scope; each name must resolve to a factor declared on the same target or one
  of its ancestors.
- **F7.** A requirement's optional `ratings` map is keyed by level names; every
  key must name a level defined by the Model's scale.

**File — rating scale**

- **F8.** The Model must declare `ratings`: a non-empty scale of at least two
  levels, listed best to worst, given inline. Each level has a unique `level` name
  and a `criterion` the evaluator judges against. A Model with no `ratings`, or a
  scale of fewer than two levels, is invalid. The format prescribes no default
  scale; an author chooses one to fit the subject.

**Reader obligations**

- **R1.** A conforming reader preserves and ignores unknown frontmatter keys and
  unknown body sections; it must not reject a file for containing them.
- **R2.** Malformed *recognized* content is an error, not an extension — a wrong key
  type, duplicate sibling keys, duplicate scale level names, a secondary factor
  outside visible scope, or an empty assessment. A reader may warn when an unknown
  key looks like a typo of a recognized one.

### Schema

The frontmatter is a **Model**: the file's root mapping, an apex **Target** plus the
one key that belongs to the file as a whole rather than to any single target —
`ratings`. The types below are the normative contract; each `<TypeName>` reference
resolves to the type of that name, and all keys are optional unless noted.

```yaml
# Model = Target + model-level keys. Appears once, as the root mapping.
ratings: <RatingScale>            # required (F8); the scale shared by all requirements
# ...plus every Target key below, applied to the apex target.
```

A **Target** is the recursive node type; `targets` nests Targets to any depth, and
none of them carry `ratings`:

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

A **Requirement** declares exactly one `assessment` (F5):

```yaml
# Requirement
assessment: <string>                # required; single non-empty scalar
factors:                            # optional secondary factor names in scope (F6)
  - <factor-name>
ratings:                            # optional per-requirement criteria (F7)
  <level-name>: <criterion>         #   keyed by a level name of the active scale
```

A **Factor** is a quality lens carrying its own requirements:

```yaml
# Factor
description: <string>               # recommended
requirements:                       # map of statement -> Requirement
  <requirement-statement>: <Requirement>
```

A **RatingScale** is an ordered sequence of levels, best to worst (F8):

```yaml
# RatingScale = list of RatingLevel
- level: <level-name>               # required; unique within the scale; position is rank
  title: <string>                   # optional human label
  criterion: <string>               # required; the criterion the evaluator judges against
```

## Model Semantics

This section fixes what the structured model *means*. It states no new F-rule, but
an evaluator that reads these elements differently does not interoperate.

The file itself is the apex target; every child under `targets:` is another Target of
the same shape. The model keeps three concepts separate:

- **Target** — a thing evaluated, bound to the material it is assessed from by
  `source`. Target names are open, user-chosen identifiers such as `source-code`,
  `payment-flows`, or `documentation`.
- **Factor** — a quality lens, such as `reliability` or `maintainability`,
  declared on a Target and visible to that target and its descendants only.
- **Requirement** — an assessable expectation. Its `assessment` is performed
  against the target's `source`, producing a **finding**. The rating criteria are
  then applied to that finding to produce a **result**, whose recorded value is a
  **rating**.

A requirement may sit directly under a target, where it is unlensed, or under a
factor, where it is assessed through that lens. It may also name secondary `factors`
it supports, so one result contributes to several factor views without repeating the
requirement.

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

A target should lead to at least one requirement — its own, one carried by a factor,
or one contributed by a descendant. A target whose subtree holds none assesses
nothing, though a pure grouping target stays meaningful as long as its descendants
carry requirements.

### Targets And Source

`targets` maps a target name to a Target. Position is lineage: a child inherits every
applicable declaration from its ancestors. Names are open and user-chosen, drawn
from no fixed set.

`source` identifies the material evaluated for that target — a single string,
conventionally a path, a glob, or a URL, though a reader may support whatever forms
suit it. Paths and globs resolve relative to the containing `QUALITY.md` file.

When `source` is omitted, it defaults to the file's own directory and all
subdirectories, recursively; a grouping target may leave it implicit and let
children narrow it.

### Factors

`factors` maps a factor name to a factor entry. A factor is a quality attribute
scoped to the declaring target's subtree: declared at the apex it is project-wide;
declared on `targets.docs` it applies only to that target and its descendants.

Factor identity is local to its scope — the same name on two unrelated targets
denotes two distinct factors. Within a scope a descendant may *refine* an inherited
factor by adding requirements under the same name; it should not *redefine* it with a
contradictory meaning. Because "contradictory" is a judgment, this is guidance a tool
may warn on rather than a deterministic rule.

A factor should carry at least one requirement — its own, one added by a refinement,
or one that names it as a secondary factor. A factor that nothing contributes to is a
lens over nothing.

A factor should have a `description` of the attribute itself: what it means here, why
it matters, and how it differs from its siblings. The description should not merely
restate the requirements attached to it.

### Requirements

`requirements` maps a requirement statement to a requirement entry. The key is the
requirement's identity in reports and results, and each requirement produces exactly
one result, recorded against the target that declares it.

A direct requirement is assessed against the target's `source` with no primary
factor. A requirement under a factor is assessed against the same source but also
joins that factor's rollup.

`assessment` is the instruction that produces a finding — inline criteria text or a
path to a document of criteria. It is never a list (F5); if one statement needs
several independent assessments, split it into several requirements.

The optional secondary `factors` list (F6) lets one result appear in additional
factor views without changing the requirement's primary placement. The optional
`ratings` map sets this requirement's own criteria, keyed by the scale's level names
(F7); it changes only the criteria, not the levels, their order, or their display
names.

### Containment And Evaluation

Containment describes how an evaluator treats the tree. These are **evaluator
obligations**, not file rules: they bind any reader that evaluates a model, even
though no file can be checked against them. An evaluator records and rolls up
results, while each individual assessment is performed by a judge — a person or a
model. Tools that disagree here produce results that cannot be compared.

Containment is the only inheritance primitive: a target owns what it declares and
inherits its ancestors' applicable factors and requirements. A requirement is
**assessed once, at the target that declares it**, against that target's `source` —
one assessment, one finding, one result. Containment then makes it *govern* every
descendant: the requirement joins their inherited context and its single result
covers their subtree, but it is never re-assessed against a descendant's narrower
source and never yields a second result.

Because a target's `source` ordinarily spans its descendants', that one assessment
routinely inspects artifacts that also belong to sub-targets. That is expected: a
finding may cite material anywhere in the declaring target's source — including files
a descendant also selects — without splitting into multiple results.

Inheritance is purely additive. A descendant may add factors and requirements but
cannot remove an inherited one. To assess a concern at a finer grain, declare a
requirement at that lower target: declaration altitude is assessment altitude.

## Rating Scale

The structural rule for scales is F8; this section fixes how one is written and
applied.

The Model's `ratings` is an inline sequence of levels ordered best to worst, where
position defines rank. Each level pairs a unique `level` identifier with the
`criterion` the evaluator judges against, plus an optional `title` for display. The
evaluator applies the criteria top-down and records the best level whose criterion
the finding satisfies.

A scale need not be elaborate. The minimum is two levels — a plain pass/fail gate:

```yaml
ratings:
  - { level: pass, criterion: "Satisfies the assessment." }
  - { level: fail, criterion: "Does not satisfy the assessment." }
```

### A suggested scale: the landing zone

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

### Custom Rating Criteria

Set `ratings` on a requirement when the scale's shared criteria cannot express the
gradient that matters. The custom criteria should name ordered, mutually distinct
levels of the active scale (F7), and the evaluator assigns the best level met.

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

## Invalid Examples

Each example illustrates a rule it violates and adds none of its own.

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

## Markdown Body

The Markdown body documents *why* this is the right model — the context an evaluator
needs to interpret assessments consistently. It is optional, and the format does not
restrict it to any fixed set of sections.

*Non-normative.* The sections below are a recommended starting point, not a required
structure; teams and tools may use, rename, or replace them, and a reader preserves
sections it does not recognize (R1).

| Section | What it captures |
| --- | --- |
| **Overview** | What the subject is, who depends on it, and what "good" means here. |
| **Scope** | The model boundary: what it covers and deliberately leaves out. Out-of-scope concerns are exclusions by design, not deficiencies. |
| **Needs** | Stakeholder outcomes the requirements answer to. |
| **Risks** | What goes wrong, and for whom, if a need is not met. |
| **Targets and factors** | A prose mirror of the target tree: each target's role, the scoped factors declared there, and why those lenses belong there. |
| **Known gaps** | In-scope quality concerns deliberately deferred, each with a brief reason. |

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

## Extensibility And Versioning

The minimal structural core is a frontmatter Model with its required `ratings` scale
(F1, F8). Because the Target fields are optional, a Model carrying only `ratings` is
structurally valid but not useful; tools should warn when one declares no
requirements, factors, or child targets. The body is optional. The format grows
through use: a reader preserves unknown frontmatter keys and unknown body sections
(R1) rather than rejecting them, while malformed recognized content is an error, not
an extension (R2).

### Extending The Format

*Non-normative.* Target and factor names are an open vocabulary, so producers may
extend freely: add frontmatter keys and body sections for your own tooling, kept
additive so other readers can ignore them (R1). A tool may enforce a stricter local
policy that rejects what the format merely discourages, provided it never treats a
conforming file as malformed (Conformance). A scaffolder may seed opinionated
defaults — the landing-zone scale, the recommended body sections — as a starting
point an author is free to replace.

### Edge Cases

How the rules above resolve specific cases; each row applies a rule already stated:

| Case | Treatment |
| --- | --- |
| **No frontmatter, or no closing `---`** | Invalid (F1): there is no model to read. |
| **Empty `targets`** | Valid but usually a warning; it declares no child targets. |
| **Empty Target** | Valid as a grouping placeholder; a tool may warn if it contributes nothing. |
| **Factor refinement that changes description incompatibly** | Discouraged, not invalid; a tool may warn. Adding requirements or compatible detail is fine. |
| **Secondary `factors` entry outside visible scope** | Invalid (F6): the factor name must resolve on the current target or an ancestor. |
| **Empty assessment value** | Invalid (F5): the assessment is the criteria and must be non-empty. |
| **Duplicate sibling names** | Invalid (R2): YAML duplicate keys and duplicate `ratings[].level` names are rejected structurally. |
| **Ordering** | Significant within `ratings` (F8); otherwise evaluation does not depend on authoring order, though tools may preserve it for display. |
| **Case** | Schema keys and author-defined names are case-sensitive. |

### Versioning

The format evolves additively. New optional keys and sections may be introduced
without invalidating existing conforming files. A change that would invalidate an
existing conforming file is breaking and reserved for a new major format version.
