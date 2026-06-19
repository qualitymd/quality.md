# QUALITY.md Specification

**Specification version:** 0.1 (Draft)

This document specifies the `QUALITY.md` standard: a Markdown file with YAML
frontmatter that declares a quality model and a Markdown body that documents its
context. The specification is a reference for authors, parsers, linters,
evaluators, report renderers, and tools that need to exchange or interpret
`QUALITY.md` documents consistently.

This specification defines the document structure, model vocabulary,
frontmatter schema, evaluation semantics, and minimum report semantics of
`QUALITY.md`. Authoring advice, examples, and notes are informative unless
explicitly stated otherwise.

The specification version identifies the `QUALITY.md` document format and
evaluation semantics defined here. See
[Versioning](docs/reference/versioning.md#specification-version) for the
project's specification-version policy.

## Conformance

Conforming uses and applications of `QUALITY.md` MUST fulfill all applicable
normative requirements in this specification.

The key words "MUST", "MUST NOT", "REQUIRED", "SHALL", "SHALL NOT", "SHOULD",
"SHOULD NOT", "RECOMMENDED", "NOT RECOMMENDED", "MAY", and "OPTIONAL" in this
document are to be interpreted as described in BCP 14,
[RFC 2119](docs/reference/rfc2119.md), and
[RFC 8174](docs/reference/rfc8174.md) when, and only when, they appear in all
capitals, as shown here.

All content in this specification is normative except sections or passages
explicitly marked as non-normative, informative, examples, or notes.

Examples and notes are non-normative. They illustrate intended interpretation,
motivation, or common edge cases; they do not add conformance requirements.

A conforming application can provide additional functionality, fields, output
formats, filters, aggregation methods, or authoring aids, but it MUST NOT do so
where explicitly disallowed or where doing so would make the application
non-conforming.

### Conformance Classes

A `QUALITY.md` document conforms when it satisfies the document and frontmatter
requirements in this specification.

A parser conforms when it accepts conforming documents, rejects frontmatter that
is not valid YAML, and exposes the parsed frontmatter and body without changing
their meaning.

A linter conforms when it reports violations of the frontmatter schema and
structural semantic requirements defined in this specification. A linter can
report additional diagnostics.

An evaluator conforms when it interprets a conforming model according to the
evaluation semantics defined in this specification.

A report renderer conforms when it preserves the required report information
and distinctions defined in [Report Semantics](#report-semantics), regardless of
the concrete rendering format.

## Terminology

**Quality Model**: A structured, declarative description of what quality means
for a subject.

**Entity**: A thing evaluated for quality.

**Model**: The root object in a `QUALITY.md` file. A Model is the apex Target
plus the model-wide Rating Scale.

**Target**: An entity or set of entities with quality requirements subject to
evaluation.

**Source**: A selector describing the entities evaluated by a Target.

**Factor**: A quality characteristic or attribute through which a Target's
quality is described. A Factor groups connected Requirements and can be
decomposed into sub-factors.

**Requirement**: An assessable quality expectation. A Requirement has a
statement, an Assessment, zero or more explicit Factor references, and optional
per-level criterion overrides.

**Assessment**: The means for assessing a Target's Source against a Requirement,
stated inline or as a reference to an entity that defines those means. An
Assessment produces Findings.

**Finding**: A single observation produced by an Assessment. A Finding records
what was observed and is not itself rated.

**Rating Scale**: The ordered set of Rating Levels used by a Model.

**Rating Level**: A single level on a Rating Scale, with a stable meaning and a
default criterion for rating a Requirement's Findings.

**Rating Result**: The outcome of rating a Requirement's Findings against the
Rating Scale: either one Rating Level or `not assessed`.

**Evaluation Report**: The structured result of evaluating a Model, including
scope, Findings summaries, ratings, rationales, and advice.

## Document Structure

A `QUALITY.md` document is a Markdown file containing:

1. A YAML frontmatter block containing the Model.
2. An optional Markdown body documenting the Model's context.

The document MUST begin with a valid YAML frontmatter block. The frontmatter
MUST contain a conforming Model.

The Markdown body can be empty. A conforming tool MUST preserve body content it
does not interpret. A conforming tool MUST NOT reject a document solely because
the body uses unrecognized headings or sections.

The location of a `QUALITY.md` document defines the default Source for the root
Model: the directory containing the file and all descendants. A root Model can
override that default by declaring `source`.

## Frontmatter Schema

Every property present in the frontmatter MUST use the YAML shape specified in
this section. Frontmatter that parses as YAML but does not conform to these
shapes is not a conforming `QUALITY.md` document.

Null or empty values do not satisfy required properties. A required property
with a null or empty value MUST be treated as absent.

#### Model

A Model is the root node of a `QUALITY.md` document. It has all Target
properties plus the model-wide `ratingScale`.

```yaml
title: <string>                 # Recommended
description: <string>           # Optional
ratingScale:                    # Required
  - level: <level-name>         #   Required; unique within the scale
    title: <string>             #   Optional
    description: <string>       #   Recommended
    criterion: <string>         #   Required
factors:                        # Optional*
  <factor-name>: <Factor>
requirements:                   # Optional*
  <requirement-statement>: <Requirement>
targets:                        # Optional*
  <target-name>: <Target>
source: <string>                # Optional
```

An entry on either factors, requirements, or targets MUST be supplied.

`ratingScale` is unique to the Model. A Target MUST NOT declare `ratingScale`.

### Rating Scale

`ratingScale` MUST be a sequence of at least two Rating Levels ordered from best
to worst.

At least two rating levels MUST be supplied.

Each Rating Level MUST declare:

- `level`: a non-empty scalar name unique within the Rating Scale.
- `criterion`: a non-empty scalar default criterion for assigning that Rating
  Level to a Requirement's Findings.

Each Rating Level SHOULD declare:

- `description`: the stable meaning of the level across the Model.

Each Rating Level can declare:

- `title`: a human-readable label.

A Rating Level's `description` and `criterion` have distinct semantics. The
`description` defines what the level means across the Model. The `criterion`
defines the default rule for deciding whether a Requirement's Findings are
rated at that level. Requirement-level `ratings` can override criteria, but
MUST NOT override a level's `description`, `title`, `level`, or ordering.

#### Target

A Target is the recursive node of the Model. Each entry under `targets` is a
Target.

```yaml
title: <string>                 # Recommended
description: <string>           # Optional
factors:                        # Optional*
  <factor-name>: <Factor>
requirements:                   # Optional*
  <requirement-statement>: <Requirement>
targets:                        # Optional*
  <target-name>: <Target>
source: <string>                # Optional
```

A Target can declare no `factors` or `requirements` of its own when it is used
as a grouping node for child `targets`.

When present, `title` is the Target's display name. The Target's map key remains
its identifier.

When present, `source` selects the entities evaluated by the Target. Relative
paths and globs resolve relative to the containing `QUALITY.md` file. When a
Target omits `source`, it inherits the Source of the nearest ancestor Target
that declares one; if no ancestor declares one, it inherits the document's
default Source.

Child Targets do not inherit parent Requirements. An ancestor Target's Source
can overlap with a descendant Target's Source; when that occurs, ancestor
Requirements still evaluate against the ancestor Source.

#### Factor

A Factor groups Requirements through a quality characteristic.

```yaml
description: <string>           # Recommended
factors:                        # Optional
  <factor-name>: <Factor>
requirements:                   # Optional
  <requirement-statement>: <Requirement>
```

`factors`, when present on a Factor, declares sub-factors. A sub-factor is a
Factor of the same shape, nested to any depth.

Factor identity is local to the Target on which the Factor is declared. Factors
with the same name on different Targets are distinct Factors.

#### Requirement

A Requirement is identified by its map key, the Requirement statement. A
Requirement MUST declare exactly one `assessment`.

```yaml
assessment: <string>            # Required
factors:                        # Optional; required for direct Target requirements
  - <factor-name>
ratings:                        # Optional
  <level-name>: <criterion>
```

`assessment` MUST be a single non-empty scalar. A missing, empty, null, or
list-valued `assessment` is invalid.

An `assessment` either states the means of assessing inline or references an
entity that defines them, such as a specification, guide, or checklist.
Referencing names that entity once instead of copying criteria that would drift
from their origin.

Note: This note is non-normative. A referenced entity may itself be a Target in
the Model. Referencing it by the same selector used as that Target's `source`
makes the dependency traceable from the Requirement to that Target without a
distinct link type.

Every Requirement MUST be connected to at least one Factor.

A Requirement declared under a Factor or sub-factor is connected by placement.
The containing Factor is its primary Factor. Such a Requirement can also declare
`factors`; those entries are secondary Factor references.

A Requirement declared directly under a Target is not connected by placement.
It MUST declare `factors` with at least one non-empty scalar entry.

Each explicit Factor reference MUST resolve to a Factor in scope. A Factor is in
scope when it is declared on the Target where the Requirement sits or on an
ancestor Target.

Missing `factors`, `factors: null`, `factors: []`, and sequences containing only
null or empty entries do not satisfy the Factor-reference requirement for a
direct Target Requirement.

`ratings`, when present, MUST be a map keyed by Rating Level names from the
Model's Rating Scale. Each value MUST be a non-empty scalar criterion. A
criterion override replaces only that level's criterion for that Requirement.

## Body Semantics

The Markdown body documents context for interpreting the Model. The format does
not require any body section names, ordering, or content.

The body can document the subject, scope, stakeholder needs, risks, known gaps,
or other rationale for the Model. Evaluators can use body content when judging
importance, rationale, and advice.

Known gaps documented in the body are author-declared model context. They are
distinct from a `not assessed` Rating Result, which is produced during an
Evaluation when evidence is absent or insufficient.

## Evaluation Semantics

Evaluation interprets a Model against selected Sources, produces Findings for
Requirements, rates those Findings, rolls ratings up through Factors and
Targets, and produces an Evaluation Report.

This specification defines the required observable semantics of Evaluation. An
implementation can use different internal algorithms when the resulting
interpretation is equivalent.

Evaluation proceeds through these semantic phases:

1. Define scope.
2. Assess and rate Requirements.
3. Analyze roll-ups.
4. Advise.
5. Report.

### Define Scope

By default, an Evaluation's scope is the whole Model: every Target and every
Requirement within each Target.

An Evaluation can be narrowed by Target, by Factor, or by both. A Target filter
selects a Target and its subtree. A Factor filter selects Requirements connected
to the named Factor, including Requirements that reference it as a secondary
Factor.

For every Target in scope, the evaluator resolves the Target's Source according
to [Target](#target).

A scoped Evaluation is not a whole-model verdict. A conforming report renderer
MUST distinguish a scoped Evaluation from a whole-model Evaluation.

### Assess and Rate Requirements

Each in-scope Requirement is assessed once, against the Source of the Target on
which it is declared.

For each Requirement:

1. The evaluator applies the Requirement's Assessment to the Target Source,
   producing zero or more Findings.
2. The evaluator rates the Findings together against the Model's Rating Scale,
   using the Requirement's criterion overrides when present.
3. The evaluator produces one Rating Result for the Requirement.

Findings are evidence. Individual Findings are not Rating Results.

When there are no Findings, or when the available Findings are insufficient to
rate the Requirement against the Rating Scale, the Requirement's Rating Result
MUST be `not assessed`.

### Analyze Roll-Ups

Roll-up infers ratings for Factors and Targets from Requirement Rating Results,
sub-factor ratings, child Target ratings, and the Model's body context where
relevant.

This specification does not define a numeric aggregation formula. A conforming
evaluator can use explicit weights, thresholds, or computed aggregation, but
the resulting ratings MUST preserve the relationships and distinctions defined
in this section.

For each Target in scope, processed child Targets before their ancestors:

- Each Factor receives a Factor rating or `not assessed`.
- Each Target with its own Requirements receives a local rating or
  `not assessed`.
- Each Target receives an aggregate rating or `not assessed`.

A Factor rating characterizes the Factor considering, together:

- Rating Results of Requirements declared under that Factor or its sub-factors.
- Rating Results of Requirements that explicitly reference that Factor.
- Ratings of the Factor's sub-factors.

A local rating characterizes the Target considering all Requirements declared
on that Target, each counted once regardless of how many Factors it is connected
to.

A grouping Target with no Requirements of its own has no local rating.

An aggregate rating characterizes the Target considering its local rating, when
present, together with aggregate ratings of child Targets in scope. A leaf
Target's aggregate rating equals its local rating.

`not assessed` Requirement, Factor, local, or aggregate outcomes MUST remain
distinct from Rating Levels. When too little has been assessed to responsibly
infer a roll-up rating, the roll-up outcome MUST be `not assessed`.

Each roll-up SHOULD include a rationale naming the observations, constraints,
or gaps most responsible for the outcome.

### Advise

Advice identifies improvement information inferred from the Analysis. Advice
does not change any rating.

An Evaluation Report SHOULD identify:

- Key gaps: shortcomings most responsible for held-down ratings, including
  material `not assessed` areas.
- Options: available remediation options for key gaps.
- Recommendations: selected options and rationales.

## Report Semantics

An Evaluation Report is the structured result of an Evaluation.

A conforming Evaluation Report MUST include:

- The Evaluation scope.
- The rating or `not assessed` outcome for the in-scope root Target.
- A rationale for the in-scope root Target outcome.
- For each Target in scope, recursively:
  - each Requirement's Findings summary, Rating Result, and rationale;
  - each Factor's rating or `not assessed` outcome and rationale, including
    sub-factors;
  - the Target's local rating or `not assessed` outcome, when a local rating
    exists; and
  - the Target's aggregate rating or `not assessed` outcome.
- Advice, when produced.

`not assessed` outcomes MUST be shown wherever they occur and MUST be distinct
from Rating Levels.

When a Factor rating includes Requirements declared on descendant Targets
because those Requirements reference the Factor, the report SHOULD make those
contributing Requirements identifiable.

A report renderer can render the Evaluation Report as Markdown, JSON, terminal
output, a gate result, or another format, provided the required information and
distinctions are preserved.

## Extensions

Applications can define additional frontmatter properties, report fields,
filters, output formats, or evaluation mechanics unless this specification
disallows them.

Extensions MUST NOT change the meaning of conforming properties defined in this
specification. Extensions should use names unlikely to conflict with future
versions of this specification.

A conforming tool MUST preserve body content it does not understand. A tool that
rewrites frontmatter should preserve unrecognized extension fields unless doing
so would produce invalid YAML or a non-conforming document.

## Appendix A: Suggested Rating Scale

This appendix is non-normative.

When a graded scale fits and an author has no stronger domain-specific scale, a
tool can seed the following four-level scale:

```yaml
ratingScale:
  - level: outstanding
    title: Outstanding
    description: "The stretch band: reached only with significant extra effort."
    criterion: "Exceeds the requirement; satisfies it with margin to spare."
  - level: target
    title: Target
    description: "The level to aim for: achievable at reasonable cost and effort."
    criterion: "Satisfies the requirement."
  - level: minimum
    title: Minimum
    description: "The acceptable floor: less than the target, but good enough to ship."
    criterion: "Falls short of the target but remains acceptable."
  - level: unacceptable
    title: Unacceptable
    description: "Below the floor: not good enough to ship."
    criterion: "Does not meet the requirement to an acceptable degree."
```

## Appendix B: Minimal Example

This appendix is non-normative.

```markdown
---
title: Acme Checkout API
description: Public API for accepting and settling customer payments.
ratingScale:
  - level: target
    description: "Meets the agreed quality bar."
    criterion: "Satisfies the requirement."
  - level: unacceptable
    description: "Does not meet the agreed quality bar."
    criterion: "Does not satisfy the requirement."
factors:
  reliability:
    description: The API continues to accept and durably record orders.
    requirements:
      "Checkout requests are durably recorded":
        assessment: Review production write-path telemetry and recovery tests.
---

# Acme Checkout API

This model covers the checkout API and the payment write path it owns.
```
