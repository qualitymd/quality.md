# QUALITY.md Format

`QUALITY.md` is a plain-text quality model: YAML frontmatter containing the
structured model, followed by a Markdown body that explains the model's context
and rationale.

Throughout this document, **must** marks a hard rule for a conforming file or
reader, **should** marks a recommendation that may be departed from with good
reason, and **may** marks an option.

This specification is **normative in its entirety**. Normative strength is carried
by those keywords — not by a section's placement — and a passage that states the
meaning of the format is binding even when it uses none of them. Three things carry
no conformance force, and each is marked as such: **examples**, which illustrate
rules stated elsewhere and add none; passages and sections labeled *Non-normative*,
which a conforming tool may wholly ignore; and latitude granted explicitly with
**may**.

Conformance has two subjects. A conforming **file** satisfies the file rules
(F-rules); a conforming **reader** — a validator, evaluator, or any tool that
consumes a file — satisfies the reader rules (R-rules) and the evaluation
semantics. The F- and R-rules in the **Conformance** section below are the subset a
tool can check mechanically, gathered so independent validators agree on one floor.
The sections after it are equally normative: they fix the meaning of the model and
the behavior of an evaluator — the shared semantics that make findings and ratings
portable between tools — but cannot be reduced to a mechanical file check.

## Conformance

This section gathers the **mechanically checkable** rules — those a tool can decide
by inspecting a file or its own behavior, without judgment. A conforming **file**
satisfies the file rules (F-rules); a conforming **reader** — a validator,
evaluator, or any tool that consumes the file — satisfies the reader rules
(R-rules). They are the floor that lets independent validators agree. The semantic
and evaluation rules in the sections that follow are equally normative but cannot be
reduced to a mechanical check; together they form the full contract.

Conformance is the floor, not the ceiling. A tool may layer additional checks
on top of these rules — warnings, suggestions, stylistic opinions, even errors
under a stricter local policy — and the sections below point to many such
opportunities. What a tool must not do is contradict them: it must not treat a
conforming file as malformed, nor accept as valid what the F-rules forbid. The
format prescribes the contract, not the tool; this document names tools only by
way of example.

**File — frontmatter and body**

- **F1.** A conforming file begins with a fenced YAML frontmatter block
  (`---` … `---`) whose content is a single mapping: the **Model**. The Model is
  the apex **Target** extended with model-level keys. A file with no frontmatter
  or no closing fence is not a conforming `QUALITY.md`.
- **F2.** An optional Markdown body may follow the closing fence. `ratings` is a
  Model key, not a Target key: it appears only on the root mapping, never on a
  nested target.

**File — target**

- **F3.** A Target is a mapping. Its recognized keys are `source`,
  `requirements`, `factors`, and `targets`; the Model adds `ratings` to that set
  on the root mapping. The Target keys are optional and a target declares only
  what it adds; `ratings` is required on the Model (F8). Target and factor names
  are an open, case-sensitive vocabulary.
- **F4.** Key types: `source` is a string; `requirements`
  is a map of statement → Requirement; `factors` is a map of name → Factor;
  `targets` is a map of name → Target (recursively the same shape).

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
  levels listed best to worst, given inline. Each level has a unique `level` name
  and a `criterion` the evaluator judges against. A Model with no `ratings`, or a
  scale of fewer than two levels, is invalid. The format prescribes no default
  scale; an author chooses one to fit the subject.

**Reader obligations**

- **R1.** A conforming reader preserves and ignores unknown frontmatter keys and
  unknown body sections; it must not reject a file for containing them.
- **R2.** Malformed *recognized* content is an error, not an extension —
  including a wrong key type, duplicate sibling keys, duplicate scale level
  names, a secondary factor outside visible scope, or an empty assessment. A
  reader may warn when an unknown key looks like a typo of a recognized one.

### Schema

The frontmatter is a **Model**: an apex **Target** extended with model-level keys.
The types below are the normative contract; each `<TypeName>` reference resolves to
the type of that name. All keys are optional unless noted.

A **Model** is the file's root mapping — an apex `Target` plus the keys that belong
to the file as a whole rather than to any one target. `ratings` is the only such
key, and it appears only here (F2, F8):

```yaml
# Model = Target + model-level keys. Appears once, as the root mapping.
ratings: <RatingScale>            # required (F8); the scale shared by all requirements
# ...plus every Target key below, applied to the apex target.
```

A **Target** is the recursive node type. `targets` nests Targets to any depth;
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

The apex `Target` carries the file's requirements, factors, and source directly, so
a Model both *is* the top target and *contains* the target tree. `ratings` is the
sole key that distinguishes a Model from any other `Target`.

## Model Semantics

This section fixes what the structured model *means*. It states no new F-rule, but
the meaning it gives to targets, factors, and requirements is normative: an
evaluator that reads these elements differently does not interoperate.

The frontmatter is a **Model**: the apex **Target** plus the file-level `ratings`
scale. The file itself is the apex target, and every child under `targets:` is
another Target with the same shape. The Model keeps three concepts separate:

- **Target** — a thing evaluated, bound to the material it is assessed from by
  `source`. Target names are open, user-chosen identifiers such as `source-code`,
  `payment-flows`, or `documentation`.
- **Factor** — a quality lens, such as `reliability` or `maintainability`,
  declared on a Target and visible to that target and its descendants only.
- **Requirement** — an assessable expectation. Its `assessment` is performed
  against the target's `source`, producing a **finding**. The rating criteria are
  then applied to that finding to produce a **result**, whose recorded value is a
  **rating**.

A requirement may be placed directly under a target, where it is unlensed, or
under a factor, where it is assessed through that lens. A requirement may also
name secondary `factors` it supports so one result can contribute to several
factor views without repeating the requirement.

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

A target should lead to at least one requirement — declared on it, carried by one
of its factors, or contributed by a descendant target. A target whose subtree
holds no requirement assesses nothing. A pure grouping target stays meaningful as
long as its descendants carry requirements.

### Targets And Source

`targets` is a map of target name to Target. Position is
lineage: a child inherits all applicable declarations from its ancestors. A
catalog may seed names or baseline assessments, but a name with no catalog match
is valid and simply starts with no baseline content.

`source` identifies the material evaluated for that target. It is a single
string, conventionally interpreted as a path, a glob, or a URL; a reader is free
to support whatever forms suit it. Paths and globs resolve relative to the
containing `QUALITY.md` file.

When `source` is omitted, it defaults to the `QUALITY.md` file's directory and
all subdirectories, recursively. A grouping target may leave `source` implicit and
let children narrow it.

### Factors

`factors` is a map of factor name to factor entry. A factor is a
quality attribute scoped to the declaring target's subtree. A factor declared at
the apex is project-wide; a factor declared on `targets.docs` applies only to
that target and its descendants.

Factor identity is local to its scope. The same factor name declared on two
unrelated targets denotes two distinct factors. Within a scope, a descendant may
refine an inherited factor by adding requirements under the same factor name. It
should not redefine an inherited factor with a contradictory meaning; because
"contradictory" is a judgment, this is guidance a tool may warn on rather than a
deterministic rule.

A factor should carry at least one requirement — declared under it, added by a
descendant refinement, or named as a secondary factor by a requirement in its
scope. A factor that no requirement contributes to is a lens over nothing.

A factor should have a `description` explaining the quality attribute itself:
what it means here, why it matters, and how it differs from sibling factors. The
description should not merely restate the requirements attached to it.

### Requirements

`requirements` is a map of requirement statement to requirement
entry. The key is the requirement's identity in reports and results. Each
requirement produces exactly one result, recorded against the target that
declares it.

A direct requirement under a target is assessed against that target's `source`
without a primary factor. A requirement under a factor is assessed against the
same target source, but is also part of that factor's rollup.

`assessment` is the instruction that produces a finding — inline criteria text or
a path to a document of criteria. It is never a list of separate assessments
(F5). If one statement needs several independent assessments, split it into
several requirements.

The optional `factors` list names secondary factors the requirement supports;
each name must be in visible scope (F6). Secondary factors do not change the
requirement's primary placement; they let one result appear in additional factor
views.

A requirement may optionally set its own rating criteria with a `ratings` map
keyed by the scale's level names (F7). It changes only the criteria for this
requirement; it does not define levels, order, or display names.

### Containment And Evaluation

Containment describes how an evaluator treats the tree. These are **evaluator
obligations**, not file rules: they bind any conforming reader that evaluates a
model, even though no file can be checked against them. An evaluator records and
rolls up results, while each individual assessment is performed by a judge — a
person or a model. Tools that disagree here produce results that cannot be compared.

Containment is the only inheritance primitive. A target owns what it declares and
inherits applicable ancestor factors, requirements, and baseline content.

A requirement is **assessed once, at the target that declares it**, against that
target's `source` — one assessment, one finding, one result. Containment then
makes the requirement *govern* every descendant target: it joins their inherited
context and its result covers their subtree. Governing a subtree is not
re-assessment — an inherited requirement is never evaluated again against a
descendant's narrower source and never produces a second result.

Because a target's `source` ordinarily spans its descendants', that single
assessment routinely inspects artifacts that also belong to sub-targets. That is
expected: a finding may cite material anywhere in the declaring target's source,
including files a descendant target also selects, without splitting into multiple
results.

Inheritance is purely additive: a descendant may add factors and requirements,
but it cannot remove an inherited requirement and does not re-assess it. To assess
a concern at a finer grain, declare a requirement at that lower target;
declaration altitude is assessment altitude.

## Rating Scale

The structural rule for scales is F8; this section fixes how a scale is written and
applied.

The Model's required `ratings` value defines the scale shared by requirements. It
is an inline sequence of levels, ordered best to worst; position defines rank.
The format prescribes no default — the author picks a scale that fits the
subject, whether a binary gate, a graded rubric, or a maturity ladder.

```yaml
ratings:
  - { level: A, title: Excellent, criterion: "Fully satisfies the assessment; no material gaps." }
  - { level: B, title: Good, criterion: "Satisfies the assessment with only trivial gaps." }
  - { level: C, title: Acceptable, criterion: "Satisfies the core assessment with minor gaps." }
  - { level: D, title: Poor, criterion: "Partly satisfies the assessment; significant gaps remain." }
  - { level: E, title: Unacceptable, criterion: "Does not satisfy the assessment." }
```

The evaluator applies criteria top-down and records the best level whose
criterion the finding satisfies.

A scale need not be elaborate. The minimum is two levels — a plain pass/fail
gate:

```yaml
ratings:
  - { level: pass, criterion: "Satisfies the assessment." }
  - { level: fail, criterion: "Does not satisfy the assessment." }
```

### A suggested scale: the landing zone

*Non-normative.* The format prescribes no default scale (F8). When a graded scale
fits but you have no strong preference, the following four-level scale is a
reasonable starting point, and a scaffolding tool may seed it. Its vocabulary and
best-to-worst framing are adapted from the Agile Landing Zone pattern —
**outstanding** exceeds the goal, **target** meets it, **minimum** is the acceptable
floor, and **unacceptable** is below that floor:

```yaml
ratings:
  - { level: outstanding,  title: Outstanding,  criterion: "Exceeds the requirement; satisfies it with margin to spare." }
  - { level: target,       title: Target,       criterion: "Satisfies the requirement." }
  - { level: minimum,      title: Minimum,      criterion: "Falls short of the goal but stays at the acceptable floor." }
  - { level: unacceptable, title: Unacceptable, criterion: "Falls below the acceptable floor." }
```

### Custom Rating Criteria

Set `ratings` on a requirement when the default criteria cannot
express the gradient that matters. The custom criteria should name ordered,
mutually distinct levels of the active scale (F7), and the evaluator assigns the
best level met.

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

Do not customize `ratings` merely to restate "met" and "not met"; spend that text
on the `assessment` instead.

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

A descendant that re-declares an inherited factor with a contradictory meaning is
discouraged rather than invalid; a tool may warn:

```yaml
targets:
  source-code:
    factors:
      maintainability:
        description: How readily code can be changed.
    targets:
      generated:
        factors:
          maintainability:
            description: How quickly generated code can be regenerated.
            # discouraged: redefines an inherited factor instead of refining it
```

An empty `targets: {}` map is valid but meaningless and a tool may warn. An empty
Target is structurally valid as a grouping target, but a useful model should
eventually declare `source`, `requirements`, `factors`, or child `targets`.

## Markdown Body

The Markdown body documents why the structured model is the right
one. It gives the context an evaluator needs to interpret assessments
consistently. The body is optional, and the format does not restrict it to any
fixed set of sections.

*Non-normative.* The sections below are a recommended starting point, not a
required structure; teams and tools may use, rename, or replace them. A reader
preserves sections it does not recognize (R1).

Recommended sections:

| Section | What it captures |
| --- | --- |
| **Overview** | What the subject is, who depends on it, and what "good" means here. |
| **Scope** | The model boundary: what it covers and deliberately leaves out. Out-of-scope concerns are exclusions by design, not deficiencies. |
| **Needs** | Stakeholder outcomes the requirements answer to. |
| **Risks** | What goes wrong, and for whom, if a need is not met. |
| **Targets and factors** | A prose mirror of the target tree: each target's role, the scoped factors declared there, and why those lenses belong there. |
| **Known gaps** | In-scope quality concerns deliberately deferred, each with a brief reason. |

Applicability is structural: if a factor is declared on `targets.docs`, it applies
there and below, not to unrelated targets. The body should explain that structure
rather than argue exceptions in prose. **Scope** is for concerns outside the
model's remit; **Known gaps** is for concerns inside the model that are deferred.

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

Unknown sections are allowed and preserved by tools (R1).

## Federation

A repository may contain more than one `QUALITY.md`. Federation
grafts models into one target tree using the same containment rule as in-file
`targets:`. The composition rule below is normative; it is a composition
convention, not a new file rule.

1. **Open target vocabulary.** Target names are user-driven. A catalog may seed
   names or baseline assessments; it is not a closed enum.
2. **Factor identity is scoped.** A factor is defined where it is declared and is
   visible only to that target and descendants. Descendants may refine inherited
   factors by adding requirements, not redefine them.
3. **Containment inheritance.** A target inherits ancestor factors and
   requirements. Requirements declared on a target apply there and flow down.
   Inheritance is additive: a descendant adds factors and requirements but does
   not remove inherited ones.
4. **Baseline is the rolling root ancestor.** Shipped baseline assessments are
   the outermost target tree. Improved baseline assessments reach everyone; they
   are always evaluated and visible rather than version-pinned away.
5. **Nest vs. federate.** Nest sub-targets when parts share the target's factors.
   Federate into a separate model when a part warrants its own ownership or
   factors. Federation grafts that model as a target subtree.

How a tool discovers, evaluates, and reports across federated models is left to
the tool; the format fixes only the composition rule above.

## Extensibility And Versioning

The minimal structural core is a fenced YAML frontmatter block whose content is a
Model mapping — an apex Target plus its required `ratings` scale (F1, F8). Because
the Target fields are optional, a Model carrying only `ratings` is structurally
valid but not useful; tools should warn when a Model declares no requirements,
factors, or child targets. The Markdown body is optional.

The format grows through use. A conforming reader ignores and preserves unknown
frontmatter keys and unknown body sections (R1) rather than rejecting them.
Malformed recognized content is an error, not an extension (R2).

### Extending The Format

*Non-normative.* The rules above are what make extension safe; this is guidance on
using them. Target and factor names are an open vocabulary — a catalog may seed
names without closing the set. A tool may layer checks on top of conformance:
warnings, house style, or a stricter local policy that rejects what the format
merely discourages, provided it never treats a conforming file as malformed
(Conformance). Producers may add frontmatter keys and body sections for their own
tooling; keep them additive so other readers can ignore them (R1). A scaffolder may
seed opinionated defaults — the landing-zone scale, the recommended body sections —
as a starting point an author is free to replace.

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
