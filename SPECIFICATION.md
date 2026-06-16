# QUALITY.md Format

`QUALITY.md` is a plain-text quality model: YAML frontmatter containing the
structured model, followed by a Markdown body that explains the model's context
and rationale.

Throughout this document, **must** marks a hard rule for a conforming file or
reader, **should** marks a recommendation that may be departed from with good
reason, and **may** marks an author option. Examples illustrate rules without
adding new rules.

## Quality Model

The frontmatter is a single recursive **target** node. The file itself is the
apex target. Every child under `targets:` is another target node with the same
shape.

The model keeps three concepts separate:

- **Target** - a thing evaluated. A target node is bound to the material it is
  assessed from by `source`. Target names are open, user-chosen identifiers such
  as `source-code`, `payment-flows`, or `documentation`.
- **Factor** - a quality lens, such as `reliability` or `maintainability`.
  Declared on a target node, it is visible to that node and its descendants only.
- **Requirement** - an assessable expectation. Its `assessment` is performed
  against the target's `source`, producing a **finding**. The rating criteria are
  then applied to that finding to produce a **result**, whose recorded value is a
  **rating**.

A requirement may be placed directly under a target, where it is unlensed, or
under a factor, where it is assessed through that lens. A requirement may also
name secondary `factors` it supports so one result can contribute to several
factor views without repeating the requirement.

Example:

```yaml
---
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
  docs: ./docs
---
```

The `docs: ./docs` entry is shorthand for `docs: { source: ./docs }`.

### Schema

```yaml
ratings:                         # optional; inline scale sequence or path to a shared YAML scale
  - level: <level-name>           # levels listed best to worst; position is rank
    displayName: <string>         # optional human label
    criterion: <string>           # optional criterion for this rating level

source: <path | glob | URL | list> # optional; apex target source
requirements:
  <requirement-statement>:
    assessment: <text | path>     # required, single scalar assessment
    factors:                      # optional secondary factor names in scope
      - <factor-name>
    ratings:                      # optional per-requirement criteria by level name
      <level-name>: <criterion>
factors:
  <factor-name>:
    description: <string>         # recommended
    requirements:
      <requirement-statement>: <Requirement>
targets:
  <target-name>: <Target | source-shorthand>
```

The root mapping itself is the apex target. `source`, `requirements`, `factors`,
and `targets` are all optional on a target node; a node declares only
what it adds.

### Targets And Source

`targets` is a map of target name to target node. Position is lineage: a child
inherits all applicable declarations from its ancestors. Target names are open
vocabulary. A catalog may seed names or baseline
assessments, but a name with no catalog match is valid and simply starts with no
baseline content.

`source` identifies the material evaluated for that target. It may be a path, a
glob, a URL, or a list. Paths and globs resolve relative to the containing
`QUALITY.md` file. A list is applied in order; an entry beginning with `!`
excludes files matched earlier:

```yaml
source:
  - ./src/**
  - "!./src/generated/**"
```

When `source` is omitted, it defaults to the `QUALITY.md` file's directory and
all subdirectories, recursively. A grouping node may leave `source` implicit and
let children narrow it. A scalar child target is shorthand for a target node whose
only field is `source`.

### Factors

`factors` is a map of factor name to factor entry. A factor is a quality
attribute scoped to the declaring target's subtree. A factor declared at the apex
is project-wide; a factor declared on `targets.docs` applies only to that target
and its descendants.

Factor identity is local to its scope. The same factor name declared on two
unrelated targets denotes two distinct factors. Within a scope, a descendant may
refine an inherited factor by adding requirements under the same factor name, but
must not redefine it with a contradictory meaning.

A factor should have a `description` explaining the quality attribute itself:
what it means here, why it matters, and how it differs from sibling factors. The
description should not merely restate the requirements attached to it.

### Requirements

`requirements` is a map of requirement statement to requirement entry. The key is
the requirement's identity in reports and results. A requirement must
declare exactly one non-empty `assessment`. It produces exactly one result,
recorded against the target that declares it.

A direct requirement under a target is assessed against that target's `source`
without a primary factor. A requirement under a factor is assessed against the
same target source, but is also part of that factor's rollup.

`assessment` is a single scalar: either inline criteria text or a path to a
document containing those criteria. It is the instruction that produces a finding.
It is never a list of separate assessments. If one statement needs several
independent assessments, split it into several requirements.

The optional `factors` list names secondary factors the requirement supports.
Each name must resolve to a factor declared on the current target or one of its
ancestors. A name outside that visible scope is a structural diagnostic. Secondary
factors do not change the requirement's primary placement; they let one result
appear in additional factor views.

A requirement may optionally set its own rating criteria with a `ratings` map
keyed by the scale's level names. It changes only the criteria for this
requirement; it does not define levels, order, or display names.

### Containment

Containment is the only inheritance primitive. A target owns what it declares and
inherits applicable ancestor factors, requirements, and baseline content.

A requirement is **assessed once, at the target that declares it**, against that
target's `source` - one assessment, one finding, one result. Containment then
makes the requirement *govern* every descendant target: it joins their inherited
context and its result covers their subtree. Governing a subtree is not
re-assessment - an inherited requirement is never evaluated again against a
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

The optional top-level `ratings` value defines the scale shared by requirements.
It may be an inline sequence of levels or a path to a YAML file holding that
sequence. The sequence is ordered best to worst; position defines rank.

```yaml
ratings:
  - { level: A, displayName: Excellent, criterion: "Fully satisfies the assessment; no material gaps." }
  - { level: B, displayName: Good }
  - { level: C, displayName: Acceptable, criterion: "Satisfies the core assessment with minor gaps." }
  - { level: D, displayName: Poor }
  - { level: E, displayName: Unacceptable, criterion: "Does not satisfy the assessment." }
```

The evaluator applies criteria top-down and records the best level whose
criterion the finding satisfies. A level without a criterion is an intermediate
band: it may be assigned when the finding clearly does better than the level below
but does not meet the stated criterion above.

When `ratings` is omitted, the default scale is:

```yaml
ratings:
  - { level: outstanding,  displayName: Outstanding,  criterion: "Exceeds the requirement; satisfies it with margin to spare." }
  - { level: target,       displayName: Target,       criterion: "Satisfies the requirement." }
  - { level: minimum,      displayName: Minimum,      criterion: "Falls short of the goal but stays at the acceptable floor." }
  - { level: unacceptable, displayName: Unacceptable, criterion: "Falls below the acceptable floor." }
```

The default vocabulary and best-to-worst framing are adapted from the Agile
Landing Zone pattern: **outstanding** exceeds the goal, **target** meets it,
**minimum** is the acceptable floor, and **unacceptable** is below that floor.

### Custom Rating Criteria

Set `ratings` on a requirement when the default criteria cannot express the
gradient that matters. The custom criteria should name ordered, mutually distinct
levels and the evaluator assigns the best level met.

A measured-bound example:

```yaml
requirements:
  "evaluation completes fast enough to sit in a pre-commit hook":
    assessment: >
      Wall-clock time for `qualitymd lint` on a typical model, measured cold on CI
      hardware, p95 over 20 runs.
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

```yaml
requirements:
  "input is validated":
    assessment:
      - "Query parameters are validated." # invalid: assessment must be one scalar
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
          - clarity # invalid: sibling target's factor is out of scope
```

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
            # invalid: redefines inherited factor meaning instead of refining it
```

```yaml
ratings:
  - { level: target, criterion: "Meets the requirement." }
requirements:
  "tests cover critical paths":
    assessment: Critical paths are covered.
    ratings:
      gold: "Exceptional coverage." # invalid: `gold` is not a scale level
```

An empty `targets: {}` map is valid but meaningless and a tool may warn. An empty
target node is structurally valid as a grouping node, but a useful model should
eventually declare `source`, `requirements`, `factors`, or child `targets`.

## Markdown Body

The Markdown body documents why the structured model is the right one. It gives
the context an evaluator needs to interpret assessments consistently.

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

Example:

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

The format does not restrict the body to these sections. Unknown sections are
allowed and preserved by tools.

## Federation

A repository may contain more than one `QUALITY.md`. Federation grafts models
into one target tree using the same containment rule as in-file `targets:`.

The rules are:

1. **Open target vocabulary.** Target names are user-driven. A catalog may seed
   names or baseline assessments; it is not a closed enum.
2. **Factor identity is scoped.** A factor is defined where it is declared and is
   visible only to that target and descendants. Descendants may refine inherited
   factors by adding requirements, not redefine them.
3. **Containment inheritance.** A target inherits ancestor factors and
   requirements. Requirements declared on a node apply there and flow down.
   Inheritance is additive: a descendant adds factors and requirements but does
   not remove inherited ones.
4. **Baseline is the rolling root ancestor.** Shipped baseline assessments are
   the outermost target tree. Improved baseline assessments reach everyone; they
   are always evaluated and visible rather than version-pinned away.
5. **Nest vs. federate.** Nest sub-targets when parts share the node's factors.
   Federate into a separate model when a part warrants its own ownership or
   factors. Federation grafts that model as a target subtree.

Operational discovery, evaluation runs, and reports are specified in
[`specs/cli-federation.md`](specs/cli-federation.md).

## Conformance, Extensibility, And Versioning

### Minimal Core

The minimal structural core is a fenced YAML frontmatter block whose content is a
target node mapping. Because all target fields are optional, `---\n{}\n---` is
structurally valid but not useful; tools should warn when a model declares no
requirements, factors, or child targets. The Markdown body and custom `ratings`
scale are optional.

### Unrecognized Content

The format grows through use. A conforming reader ignores and preserves unknown
frontmatter keys and unknown body sections rather than rejecting them. A tool may
warn when an unknown key looks like a typo of a known key. Malformed known content
is an error, not an extension.

### Edge Cases

| Case | Treatment |
| --- | --- |
| **No frontmatter, or no closing `---`** | Invalid: there is no model to read. |
| **Empty `targets`** | Valid but usually a warning; it declares no child targets. |
| **Empty target node** | Valid as a grouping placeholder; a tool may warn if it contributes nothing. |
| **Factor refinement that changes description incompatibly** | Invalid redefinition when it contradicts the inherited meaning; valid refinement when it adds requirements or compatible detail. |
| **Secondary `factors` entry outside visible scope** | Invalid: the factor name must resolve on the current target or an ancestor. |
| **Empty assessment value** | Invalid: the assessment is the criteria and must be non-empty. |
| **Duplicate sibling names** | Invalid: YAML duplicate keys are rejected, and duplicate `ratings[].level` names are rejected structurally. |
| **Ordering** | Significant within `ratings`; otherwise evaluation does not depend on authoring order, though tools may preserve it for display. |
| **Case** | Schema keys and author-defined names are case-sensitive. |

### Versioning

The format evolves additively. New optional keys and sections may be introduced
without invalidating existing conforming files. A change that would invalidate an
existing conforming file is breaking and reserved for a new major format version.
