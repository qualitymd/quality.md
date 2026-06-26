# QUALITY.md Specification

**Specification version:** 0.4 (Draft)

This document specifies the QUALITY.md standard: a Markdown file with YAML
frontmatter that declares a quality model and a Markdown body that documents its
context. The specification is a reference for authors, parsers, linters,
evaluators, report renderers, and tools that need to exchange or interpret
QUALITY.md documents consistently.

This specification defines the document structure, model vocabulary,
frontmatter schema, evaluation semantics, and minimum report semantics of
QUALITY.md. Authoring advice, examples, and notes are informative unless
explicitly stated otherwise.

The specification version identifies the QUALITY.md document format and
evaluation semantics defined here. See
[Versioning](docs/reference/versioning.md#specification-version) for the
project's specification-version policy.

## Lineage

Note: This section is non-normative.

QUALITY.md is domain agnostic: a Model can describe quality for software,
documents, data sets, research or analytical reports, services, operations,
processes, or other evaluated entities. It is informed by established
software-quality, requirements,
measurement, testing, and evaluation traditions, including ISO/IEC and
ISO/IEC/IEEE standards, CISQ structural-quality work, and earlier software
quality models such as McCall and Dromey.

These are acknowledged as influences, not normative references or conformance
targets. QUALITY.md uses its own vocabulary and structure where that makes the
format more practical, readable, and accessible. It takes useful boundaries and
vocabulary from those traditions, not their characteristic lists as default
Factors.

## Conformance

Conforming uses and applications of QUALITY.md MUST fulfill all applicable
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

A QUALITY.md document conforms when it satisfies the document and frontmatter
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
for a root area and any child Areas.

**Entity**: A thing evaluated for quality.

**Model**: The root object in a QUALITY.md file. A Model is the root area plus
the model-wide Rating Scale.

**Area**: An entity or set of entities with quality requirements subject to
evaluation.

**Area name**: A single map key under `areas`, unique among sibling Areas in
that `areas` map.

**Area ID**: The ordered path of Area names from the root Area to an Area. The
root Area ID is the empty path.

**Area title**: The required human display label stored in an Area's `title`.

**Source**: A selector describing the material evaluated by an Area.

**Factor**: A quality characteristic or attribute through which an Area's
quality is described. A Factor groups connected Requirements and can be
decomposed into sub-factors.

**Factor name**: A single map key under `factors`, unique among sibling Factors
in that `factors` map.

**Factor ID**: The declaring Area ID plus the ordered path of Factor names from
that Area's `factors` map to the Factor.

**Requirement**: An assessable quality expectation. A Requirement has a stable
Requirement name, a title, an Assessment, zero or more explicit Factor
references, and optional per-level criterion overrides.

**Assessment**: The means for assessing an Area's Source against a Requirement,
stated inline or as a reference to an entity that defines those means. An
Assessment produces Findings.

**Finding**: A single observation produced by an Assessment. A Finding records
what was observed and is not itself rated.

**Rating Scale**: The ordered set of Rating Levels used by a Model.

**Rating Level**: A single level on a Rating Scale, with a stable meaning and a
default criterion for rating a Requirement's Findings.

**Rating Level ID**: The `level` value of a Rating Level, unique within the
Model's Rating Scale.

**Rating Result**: The outcome of rating a Requirement's Findings against the
Rating Scale: either one Rating Level or `not assessed`.

**Requirement name**: A Requirement's stable map key, unique within its
declaring Area.

**Model reference**: A text form used at human/tool boundaries to address an
Area, Factor, Requirement, or Rating Level.

**Display value**: A human-facing label for a known Model concept. Display
values are not model references unless a section explicitly says so.

**Evaluation Report**: The structured result of evaluating a Model, including
scope, Findings summaries, ratings, rationales, and advice.

## Names and Model References

Area names, Factor names, Requirement names, and Rating Level IDs MUST match:

```regex
^[A-Za-z0-9](?:[A-Za-z0-9_-]*[A-Za-z0-9])?$
```

The grammar excludes `/`, `:`, spaces, dots, and leading or trailing separators
so canonical model references are unambiguous and do not resemble filesystem
paths.

Qualified Area references use `area:<area-path>`. The root Area reference is
`area:root`; nested Area references join Area names with `/`, for example
`area:webhooks` and `area:webhooks/delivery`.

Qualified Factor references use
`factor:<declaring-area-path>::<factor-path>`. The root declaring Area is
written as `root`, for example `factor:root::security` and
`factor:root::security/secrets`. Nested declaring Areas and nested Factors use
`/` within each side of the `::` separator, for example
`factor:webhooks/delivery::reliability/retry-behavior`.

Qualified Requirement references use
`requirement:<declaring-area-path>::<requirement-name>`. The root declaring Area
is written as `root`, for example `requirement:root::release-notes-current` and
`requirement:webhooks/delivery::retry-window`.

Qualified Rating Level references use `rating:<rating-level-id>`, for example
`rating:target`.

Tools that render qualified model references MUST use the typed prefixes
`area:`, `factor:`, `requirement:`, and `rating:`. Tools that parse qualified
model references MUST reject references whose segments fail the strict name
grammar or whose referenced model element does not exist.

Tools MAY render or accept unqualified references that omit the type prefix only
where the surrounding context fixes the reference type. Unqualified Area
references render as `root` or `<area-path>`; unqualified Factor references
render as `<declaring-area-path>::<factor-path>`; unqualified Rating Level
references render as the Rating Level ID.

Tools MAY render display values in human-facing reports and UI. Area display
values render the root Area as `/` and nested Areas as `<area-path>`, such as
`webhooks/delivery`. `/` is a display value, not an Area reference; tools MUST
NOT parse `/` as a qualified or unqualified Area reference. Factor display
values render as `<factor-path>`; Rating Level display values render as the
Rating Level ID unless a title is resolved from the Model.

Tools MUST NOT render or accept unqualified references on mixed-reference
surfaces or anywhere the reference type must be recoverable from the value
alone. Tools MUST NOT persist unqualified references in evaluation routine data,
`EvaluationOutputResult`, or other durable machine-readable artifacts.
Mixed-reference surfaces MUST require qualified model references.

## Document Structure

A QUALITY.md document is a Markdown file containing:

1. A YAML frontmatter block containing the Model.
2. An optional Markdown body documenting the Model's context.

The document MUST begin with a valid YAML frontmatter block. The frontmatter
MUST contain a conforming Model.

The Markdown body can be empty. A conforming tool MUST preserve body content it
does not interpret. A conforming tool MUST NOT reject a document solely because
the body uses unrecognized headings or sections.

The location of a QUALITY.md document defines the default Source for the root
area: the directory containing the file and all descendants. A root area can
override that default by declaring `source`.

## Frontmatter Schema

Every property present in the frontmatter MUST use the YAML shape specified in
this section. Frontmatter that parses as YAML but does not conform to these
shapes is not a conforming QUALITY.md document.

Null or empty values do not satisfy required properties. A required property
with a null or empty value MUST be treated as absent.

> **Non-normative.** A companion JSON Schema for this frontmatter is available as
> `quality.schema.json` (emitted by `qualitymd schema`). It is structural-only —
> it describes the shapes in this section, not the semantic rules a conforming
> tool enforces — and is subordinate to this specification and to its durable
> spec (`specs/quality-schema-json.md`). Passing structural validation does not
> imply full conformance; the semantic layer stays with the tool.

#### Model

A Model is the root node of a QUALITY.md document. It has all Area
properties plus the model-wide `ratingScale`.

```yaml
title: <string>                 # Required
description: <string>           # Optional
ratingScale:                    # Required
  - level: <rating-level-id>    #   Required; unique within the scale
    title: <string>             #   Required
    description: <string>       #   Recommended
    criterion: <string>         #   Required
factors:                        # Optional*
  <factor-name>: <Factor>
requirements:                   # Optional*
  <requirement-name>: <Requirement>
areas:                          # Optional*
  <area-name>: <Area>
source: <string>                # Optional
```

An entry on either factors, requirements, or areas MUST be supplied.

`ratingScale` is unique to the Model. An Area MUST NOT declare `ratingScale`.

### Rating Scale

`ratingScale` MUST be a sequence of at least two Rating Levels ordered from best
to worst.

At least two rating levels MUST be supplied.

Each Rating Level MUST declare:

- `level`: a non-empty scalar Rating Level ID unique within the Rating Scale.
- `title`: a non-empty scalar human-readable label.
- `criterion`: a non-empty scalar default criterion for assigning that Rating
  Level to a Requirement's Findings.

Each Rating Level SHOULD declare:

- `description`: the stable meaning of the level across the Model.

A Rating Level's `description` and `criterion` have distinct semantics. The
`description` defines what the level means across the Model. The `criterion`
defines the default rule for deciding whether a Requirement's Findings are
rated at that level. Requirement-level `ratings` can override criteria, but
MUST NOT override a level's `description`, `title`, `level`, or ordering.

#### Area

An Area is the recursive node of the Model. Each entry under `areas` is an
Area.

```yaml
title: <string>                 # Required
description: <string>           # Optional
factors:                        # Optional*
  <factor-name>: <Factor>
requirements:                   # Optional*
  <requirement-name>: <Requirement>
areas:                          # Optional*
  <area-name>: <Area>
source: <string>                # Optional
```

An Area can declare no `factors` or `requirements` of its own when it is used
as a grouping node for child `areas`.

`title` is the Area's display name. The Area's map key remains its
Area name; its Area ID is the ordered path of Area names from the root Area to
that Area.

When present, `source` selects the entities evaluated by the Area. Relative
paths and globs resolve relative to the containing QUALITY.md file. When a
Area omits `source`, it inherits the Source of the nearest ancestor Area
that declares one; if no ancestor declares one, it inherits the document's
default Source.

Child Areas do not inherit parent Requirements. An ancestor Area's Source
can overlap with a descendant Area's Source; when that occurs, ancestor
Requirements still evaluate against the ancestor Source.

#### Factor

A Factor groups Requirements through a quality characteristic.

```yaml
title: <string>                 # Required
description: <string>           # Recommended
factors:                        # Optional
  <factor-name>: <Factor>
requirements:                   # Optional
  <requirement-name>: <Requirement>
```

`factors`, when present on a Factor, declares sub-factors. A sub-factor is a
Factor of the same shape, nested to any depth.

Factor identity is local to the Area on which the Factor is declared. Factors
with the same name on different Areas are distinct Factors.

`title` is the Factor's display name. The Factor's map key remains its stable
Factor name local to the Area where the Factor is declared. Its Factor ID is the
declaring Area ID plus the ordered path of Factor names from that Area's
`factors` map to the Factor.

#### Requirement

A Requirement is identified by its map key, the stable Requirement name. A
Requirement MUST declare a `title` and exactly one `assessment`.

```yaml
title: <string>                 # Required
description: <string>           # Optional
assessment: <string>            # Required
factors:                        # Optional; required for direct Area requirements
  - <factor-name>
ratings:                        # Optional
  <rating-level-id>: <criterion>
```

The Requirement name MUST match the strict name grammar and MUST be unique
within the declaring Area. For uniqueness, Requirements declared directly under
an Area and Requirements declared under that Area's Factors or sub-factors are
considered to belong to the declaring Area.

`title` is the Requirement's human-facing statement. The Requirement's map key
is its stable Requirement name.

`assessment` MUST be a single non-empty scalar. A missing, empty, null, or
list-valued `assessment` is invalid.

An `assessment` either states the means of assessing inline or references an
entity that defines them, such as a specification, guide, or checklist.
Referencing names that entity once instead of copying criteria that would drift
from their origin.

Note: This note is non-normative. A referenced entity may itself be an Area in
the Model. Referencing it by the same selector used as that Area's `source`
makes the dependency traceable from the Requirement to that Area without a
distinct link type.

Every Requirement MUST be connected to at least one Factor.

A Requirement declared under a Factor or sub-factor is connected by placement.
The containing Factor is its primary Factor. Such a Requirement can also declare
`factors`; those entries are secondary Factor references.

A Requirement declared directly under an Area is not connected by placement.
It MUST declare `factors` with at least one non-empty scalar entry.

Each explicit Factor reference MUST resolve to a Factor in scope. A Factor is in
scope when it is declared on the Area where the Requirement sits or on an
ancestor Area.

Missing `factors`, `factors: null`, `factors: []`, and sequences containing only
null or empty entries do not satisfy the Factor-reference requirement for a
direct Area Requirement.

`ratings`, when present, MUST be a map keyed by Rating Level IDs from the
Model's Rating Scale. Each value MUST be a non-empty scalar criterion. A
criterion override replaces only that level's criterion for that Requirement.

## Body Semantics

The Markdown body documents context for building, interpreting, using, and
evaluating the Model. The format does not require any body section names,
ordering, or content.

The body can document the root area, scope, stakeholder needs, risks, unknowns,
open questions, evidence context, or other important context for the Model.
Evaluators can use body content when judging model fit, importance, rationale,
and advice.

## Evaluation Semantics

Evaluation interprets a Model against selected Sources, produces Findings for
Requirements, rates Requirement Assessments, analyzes Factors and Areas, and
produces an Evaluation Report.

This specification defines the required observable semantics of Evaluation. An
implementation can use different internal algorithms when the resulting
interpretation is equivalent.

Evaluation proceeds through these semantic phases:

1. Define scope.
2. Frame Area and Requirement Evaluation.
3. Assess Requirements.
4. Rate Requirement Assessments.
5. Analyze Factors and Areas.
6. Report.

The current Evaluation workflow is specified in
[`specs/evaluation/`](specs/evaluation/index.md). It records structured
routine outputs under `data/` and renders reports deterministically from those
outputs.

### Define Scope

By default, an Evaluation's scope is the whole Model: every Area and every
Requirement within each Area.

An Evaluation can be narrowed by Area, by Factor, or by both. An Area filter
selects an Area and its subtree. A Factor filter selects Requirements connected
to the named Factor, including Requirements that reference it as a secondary
Factor.

For every Area in scope, the evaluator resolves the Area's Source according
to [Area](#area).

A scoped Evaluation is not a whole-model verdict. A conforming report renderer
MUST distinguish a scoped Evaluation from a whole-model Evaluation.

### Frame Area and Requirement Evaluation

For each in-scope Area, the evaluator frames the Area boundary before evaluating
local Requirements, local Factors, or child Areas.

For each local Requirement, the evaluator frames evidence targets, applied Rating
Level criteria, stop conditions, and known limits before inspecting assessment
evidence.

Applied Rating Level criteria MUST be adapted before evidence judgment. They
MUST NOT be adapted to observed evidence.

### Assess Requirements

Each in-scope Requirement is assessed once, against the Source of the Area on
which it is declared.

For each Requirement, the evaluator applies the Requirement's Assessment to the
Area Source, producing zero or more Findings, Unknowns, and Evaluation Limits.

Findings are evidence-backed assessment observations. Individual Findings are
not Rating Results.

### Rate Requirement Assessments

For each Requirement Assessment, the evaluator rates the assessment against the
Model's Rating Scale, using the Requirement's criterion overrides when present.

The evaluator produces one Requirement Rating Result when the assessment
evidence is sufficient to distinguish the applied criteria.

When there are no Findings, or when the available Findings are insufficient to
rate the Requirement against the Rating Scale, the Requirement's Rating Result
MUST be `not assessed`.

### Analyze Factors and Areas

Analysis infers ratings for Factors and Areas from Requirement Rating Results,
sub-factor analyses, child Area analyses, and the Model's body context where
relevant.

This specification does not define a numeric aggregation formula. A conforming
evaluator can use explicit weights, thresholds, or computed aggregation, but
the resulting ratings MUST preserve the relationships and distinctions defined
in this section.

For each Area in scope, processed bottom-up:

- Each Factor receives a Factor rating or `not assessed`.
- Each Area with its own Requirements receives a local rating or
  `not assessed`.
- Each Area receives an aggregate rating or `not assessed`.

A Factor rating characterizes the Factor considering, together:

- Rating Results of Requirements declared under that Factor or its sub-factors.
- Rating Results of Requirements that explicitly reference that Factor.
- Ratings of the Factor's sub-factors.

A local rating characterizes the Area considering all Requirements declared
on that Area, each counted once regardless of how many Factors it is connected
to.

A grouping Area with no Requirements of its own has no local rating.

An aggregate rating characterizes the Area considering its local rating, when
present, together with aggregate ratings of child Areas in scope. A leaf
Area's aggregate rating equals its local rating.

`not assessed` Requirement, Factor, local, or aggregate outcomes MUST remain
distinct from Rating Levels. When too little has been assessed to responsibly
infer a roll-up rating, the roll-up outcome MUST be `not assessed`.

Each roll-up SHOULD include a rationale naming the observations, constraints,
or gaps most responsible for the outcome.

### Advice

Advice identifies improvement information inferred from the Analysis. Advice does
not change any rating.

Advice and recommendations are optional extensions for Evaluation v0. When
produced, advice SHOULD identify:

- Key gaps: shortcomings most responsible for held-down ratings, including
  material `not assessed` areas.
- Options: available remediation options for key gaps.
- Recommendations: selected options and rationales.

## Report Semantics

An Evaluation Report is the structured result of an Evaluation.

A conforming Evaluation Report MUST include:

- The Evaluation scope.
- The rating or `not assessed` outcome for the in-scope root Area.
- A rationale for the in-scope root Area outcome.
- For each Area in scope, recursively:
  - each Requirement's Findings summary, Rating Result, and rationale;
  - each Factor's rating or `not assessed` outcome and rationale, including
    sub-factors;
  - the Area's local rating or `not assessed` outcome, when a local rating
    exists; and
  - the Area's aggregate rating or `not assessed` outcome.
- Advice, when produced.

`not assessed` outcomes MUST be shown wherever they occur and MUST be distinct
from Rating Levels.

When a Factor rating includes Requirements declared on descendant Areas
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
    description: "The acceptable floor: less than the target, but good enough to rely on."
    criterion: "Falls short of the target but remains acceptable."
  - level: unacceptable
    title: Unacceptable
    description: "Below the floor: not good enough to rely on."
    criterion: "Does not meet the requirement to an acceptable degree."
```

## Appendix B: Minimal Example

This appendix is non-normative.

This is an illustrative software product example, not a default domain or
factor set for QUALITY.md. The same Model shape applies across domains; only the
domain-carried Factors, Requirements, and Assessments change. For a worked
non-software example, see
[Modeling quality across domains](docs/guides/model-quality-across-domains.md#worked-example-a-documentation-set).

```markdown
---
title: Acme Checkout API
description: Public API for accepting and settling customer payments.
ratingScale:
  - level: target
    title: Target
    description: "Meets the agreed quality bar."
    criterion: "Satisfies the requirement."
  - level: unacceptable
    title: Unacceptable
    description: "Does not meet the agreed quality bar."
    criterion: "Does not satisfy the requirement."
factors:
  reliability:
    title: Reliability
    description: The API continues to accept and durably record orders.
    requirements:
      checkout-requests-durable:
        title: Checkout requests are durably recorded
        assessment: Review production write-path telemetry and recovery tests.
---

# Acme Checkout API

This model covers the checkout API and the payment write path it owns.
```

## Appendix C: Invalid Counter-Examples

This appendix is non-normative.

The following snippets illustrate invalid shapes. They are intentionally
incomplete and are not standalone QUALITY.md files.

### Missing rating-level title

```yaml
ratingScale:
  - level: target
    criterion: "Satisfies the requirement."
  - level: unacceptable
    title: Unacceptable
    criterion: "Does not satisfy the requirement."
```

Invalid because every Rating Level MUST declare `level`, `title`, and
`criterion`.

### Direct area requirement without factors

```yaml
requirements:
  checkout-requests-durable:
    title: Checkout requests are durably recorded
    assessment: Review production write-path telemetry and recovery tests.
```

Invalid because a Requirement declared directly under a Area MUST declare
`factors` with at least one non-empty scalar entry.

### List-valued assessment

```yaml
factors:
  reliability:
    title: Reliability
    requirements:
      checkout-requests-durable:
        title: Checkout requests are durably recorded
        assessment:
          - Review telemetry.
          - Review recovery tests.
```

Invalid because `assessment` MUST be a single non-empty scalar.
