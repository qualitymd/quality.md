# QUALITY.md specification

**Specification version:** 0.11 (Draft)

This document specifies the QUALITY.md standard: a Markdown file with YAML
frontmatter that declares a quality model and a Markdown body that documents its
context. The specification is a reference for authors and for tools that read,
write, or exchange QUALITY.md documents and need to interpret them
consistently.

This specification defines the document structure, model vocabulary,
frontmatter schema, and model semantics of QUALITY.md: what a conforming
document is, and what a conforming document means. It does not define how a
model is to be evaluated. Evaluation methods, aggregation approaches, report
formats, and tool behavior are outside its scope;
[Appendix D](#appendix-d-an-illustrative-evaluation-method) illustrates one
evaluation method without making it normative. Authoring advice, examples, and
notes are informative unless explicitly stated otherwise.

The specification version identifies the QUALITY.md document format and
semantics defined here. See
[Versioning](docs/reference/versioning.md#specification-version) for the
project's specification-version policy.

## Lineage

Note: This section is non-normative.

QUALITY.md is domain agnostic: a model can describe quality for software,
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
factors.

## Conformance

Conformance is a property of QUALITY.md documents. A QUALITY.md document
conforms when it satisfies all applicable normative requirements in this
specification.

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

This specification defines no conformance classes for tools. Parsers, linters,
evaluators, and report renderers are expected to interpret conforming documents
according to the vocabulary and semantics defined here, but how a model is
evaluated, aggregated, or reported is a method choice this specification does
not constrain. An interpretation that contradicts
[Model semantics](#model-semantics) is a misreading of the document; an unsound
evaluation method applied to a correctly read document is outside this
specification's scope.

A document can carry additional frontmatter properties as described in
[Extensions](#extensions).

## Terminology

In this specification, bare "requirement" always means the model object defined
below; conformance obligations of this specification are called "normative
requirements" or "conformance requirements".

**Quality model**: A structured, declarative description of what quality means
for a root area and any child areas.

**Entity**: A thing evaluated for quality.

**Model**: The root object in a QUALITY.md file. A model is the root area plus
the model-wide rating scale.

**Area**: An entity or set of entities with quality requirements subject to
evaluation.

**Area name**: A single map key under `areas`, unique among sibling areas in
that `areas` map.

**Area ID**: The ordered path of area names from the root area to an area. The
root area ID is the empty path.

**Area title**: The required human display label stored in an area's `title`.

**Source**: A selector describing the material evaluated by an area.

**Factor**: A quality characteristic or attribute through which an area's
quality is described. A factor groups connected requirements and can be
decomposed into sub-factors.

**Factor name**: A single map key under `factors`, unique among sibling factors
in that `factors` map.

**Factor ID**: The declaring area ID plus the ordered path of factor names from
that area's `factors` map to the factor.

**Requirement**: An assessable quality expectation. A requirement has a stable
requirement name, a title, an assessment, zero or more explicit factor
references, and optional per-level criterion overrides.

**Assessment**: The means for assessing an area's source against a requirement,
stated inline or as a reference to an entity that defines those means.

**Rating scale**: The ordered set of rating levels used by a model.

**Rating level**: A single level on a rating scale, with a stable meaning and a
default criterion for judging a requirement at that level.

**Rating level ID**: The `level` value of a rating level, unique within the
model's rating scale.

**Requirement name**: A requirement's stable map key, unique within its
declaring area.

**Model reference**: A text form used at human/tool boundaries to address an
area, factor, requirement, or rating level.

**Display value**: A human-facing label for a known model concept. Display
values are not model references unless a section explicitly says so.

## Names and model references

Area names, factor names, requirement names, and rating level IDs MUST match:

```regex
^[A-Za-z0-9](?:[A-Za-z0-9_-]*[A-Za-z0-9])?$
```

The grammar excludes `/`, `:`, spaces, dots, and leading or trailing separators
so canonical model references are unambiguous and do not resemble filesystem
paths.

Area names MUST NOT be `root`. The token `root` is reserved for the root area in
area, factor, and requirement references.

Qualified area references use `area:<area-path>`. The root area reference is
`area:root`; nested area references join area names with `/`, for example
`area:webhooks` and `area:webhooks/delivery`.

Qualified factor references use
`factor:<declaring-area-path>::<factor-path>`. The root declaring area is
written as `root`, for example `factor:root::security` and
`factor:root::security/secrets`. Nested declaring areas and nested factors use
`/` within each side of the `::` separator, for example
`factor:webhooks/delivery::reliability/retry-behavior`.

Qualified requirement references use
`requirement:<declaring-area-path>::<requirement-name>`. The root declaring area
is written as `root`, for example `requirement:root::release-notes-current` and
`requirement:webhooks/delivery::retry-window`.

Qualified rating level references use `rating:<rating-level-id>`, for example
`rating:target`.

A qualified model reference is valid only when it uses one of the typed
prefixes `area:`, `factor:`, `requirement:`, or `rating:`, every segment
matches the strict name grammar, and the referenced model element is declared
in the model. A value that fails any of these conditions is not a qualified
model reference.

An unqualified reference omits the type prefix: unqualified area references
render as `root` or `<area-path>`, unqualified factor references render as
`<declaring-area-path>::<factor-path>`, and unqualified rating level references
render as the rating level ID. An unqualified reference is meaningful only
where the surrounding context fixes the reference type. On mixed-reference
surfaces, and anywhere the reference type must be recoverable from the value
alone — including durable machine-readable artifacts — only qualified
references are unambiguous.

Display values are human-facing labels, not references. The root area's display
value is `/`; nested areas display as `<area-path>`, such as
`webhooks/delivery`. `/` is never a model reference, qualified or unqualified.
Factor display values render as `<factor-path>`; rating level display values
render as the rating level ID unless a title is resolved from the model.

## Document structure

A QUALITY.md document is a Markdown file containing:

1. A YAML frontmatter block containing the model.
2. An optional Markdown body documenting the model's context.

The document MUST begin with a valid YAML frontmatter block. The frontmatter
MUST contain a conforming model.

The Markdown body can be empty. The body is free-form Markdown: no headings,
sections, or ordering are required, and none are reserved. Unrecognized
headings or sections do not affect a document's conformance.

The location of a QUALITY.md document defines the default source for the root
area: the directory containing the file and all descendants. A root area can
override that default by declaring `source`.

## Frontmatter schema

Every property defined in this section MUST, when present, use the YAML shape
specified here. Frontmatter whose defined properties do not conform to these
shapes is not a conforming QUALITY.md document. Properties this section does
not define are extension properties; see [Extensions](#extensions).

Null or empty values do not satisfy required properties. A required property
with a null or empty value MUST be treated as absent.

> **Non-normative.** A companion JSON Schema for this frontmatter is available as
> `quality.schema.json` (emitted by `qualitymd schema`). It is structural-only —
> it describes the shapes in this section, not the semantic rules that follow
> from this specification — and is subordinate to this specification and to its
> durable spec (`specs/quality-schema-json.md`). Passing structural validation
> does not imply document conformance.

### Model

A model is the root node of a QUALITY.md document. It has all area
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

An entry on either `factors`, `requirements`, or `areas` MUST be supplied.

`ratingScale` is unique to the model. An area MUST NOT declare `ratingScale`.

### Rating scale

`ratingScale` MUST be a sequence of at least two rating levels ordered from best
to worst.

Each rating level MUST declare:

- `level`: a non-empty scalar rating level ID unique within the rating scale.
- `title`: a non-empty scalar human-readable label.
- `criterion`: a non-empty scalar default criterion for judging whether a
  requirement sits at that rating level.

Each rating level SHOULD declare:

- `description`: the stable meaning of the level across the model.

The meaning of `description`, `criterion`, scale order, and criterion overrides
is defined in [Rating scale semantics](#rating-scale-semantics).

### Area

An area is the recursive node of the model. Each entry under `areas` is an
area.

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

An area can declare no `factors` or `requirements` of its own when it is used
as a grouping node for child `areas`.

`title` is the area's display name. The area's map key remains its
area name; its area ID is the ordered path of area names from the root area to
that area.

When present, `source` is a single string selecting the material evaluated by
the area. A selector is not limited to filesystem material: relative paths
and globs resolve relative to the containing QUALITY.md file, and any other
selector describes the evaluated material in prose — a saved query, a live
system, a body of records ("open tickets in the support queue") — resolved by
the evaluating tool's own means. The selector's kind is detected from the
string, never declared; see
[Source resolution](#source-resolution). Source inheritance and the scope of
an area's requirements are defined in [Model semantics](#model-semantics).

### Factor

A factor groups requirements through a quality characteristic.

```yaml
title: <string>                 # Required
description: <string>           # Recommended
factors:                        # Optional
  <factor-name>: <Factor>
requirements:                   # Optional
  <requirement-name>: <Requirement>
```

`factors`, when present on a factor, declares sub-factors. A sub-factor is a
factor of the same shape, nested to any depth.

Factor identity is local to the area on which the factor is declared. Factors
with the same name on different areas are distinct factors.

Within one area, authors SHOULD avoid reusing the same factor name anywhere in
that area's recursive factor tree. Canonical factor IDs remain path-based, but
explicit requirement factor references use scalar names, so repeated names in
one area can make authored references ambiguous to readers and tools.

`title` is the factor's display name. The factor's map key remains its stable
factor name local to the area where the factor is declared. Its factor ID is the
declaring area ID plus the ordered path of factor names from that area's
`factors` map to the factor.

### Requirement

A requirement is identified by its map key, the stable requirement name. A
requirement MUST declare a `title` and exactly one `assessment`.

```yaml
title: <string>                 # Required
description: <string>           # Optional
assessment: <string>            # Required
factors:                        # Optional; required for direct area requirements
  - <factor-name>
ratings:                        # Optional
  <rating-level-id>: <criterion>
```

The requirement name MUST match the strict name grammar and MUST be unique
within the declaring area. For uniqueness, requirements declared directly under
an area and requirements declared under that area's factors or sub-factors are
considered to belong to the declaring area.

`title` is the requirement's human-facing statement. The requirement's map key
is its stable requirement name.

`assessment` MUST be a single non-empty scalar. A missing, empty, null, or
list-valued `assessment` is invalid.

An `assessment` either states the means of assessing inline or references an
entity that defines them, such as a specification, guide, or checklist.
Referencing names that entity once instead of copying criteria that would drift
from their origin.

Note: This note is non-normative. A referenced entity may itself be an area in
the model. Referencing it by the same selector used as that area's `source`
makes the dependency traceable from the requirement to that area without a
distinct link type.

Every requirement MUST be connected to at least one factor.

A requirement declared under a factor or sub-factor is connected by placement.
Such a requirement can also declare `factors`; those entries are same-area
secondary factor references. A requirement declared directly under an area is
not connected by placement. It MUST declare `factors` with at least one
non-empty scalar entry.

Each explicit factor reference MUST resolve to a factor declared on the area
where the requirement sits. A factor declared on an ancestor area, sibling area,
descendant area, or unrelated area is not in scope for that requirement's
`factors` entries.

Missing `factors`, `factors: null`, `factors: []`, and sequences containing only
null or empty entries do not connect a direct area requirement to any factor,
so a requirement declared with only such values does not conform.

`ratings`, when present, MUST be a map keyed by rating level IDs from the
model's rating scale. Each value MUST be a non-empty scalar criterion. The
meaning of a criterion override is defined in
[Rating scale semantics](#rating-scale-semantics).

## Model semantics

The rules in this section define what a conforming model means. They bind every
interpretation of the document, by humans or tools, but they do not prescribe
any evaluation method, aggregation approach, or report format.

### Source resolution

An area that declares `source` selects its own evaluated material. An area that
omits `source` inherits the source of the nearest ancestor area that declares
one; if no ancestor declares one, it inherits the document's default source
defined in [Document structure](#document-structure).

An effective source selector is interpreted by kind, detected from the bare
selector string in this order: a selector containing supported glob
metacharacters (`*`, `?`, `[`) is a glob only when it parses under the supported
glob grammar; otherwise a selector that names an existing filesystem entry
relative to the containing QUALITY.md file is a path; otherwise the selector is
prose — a description of the evaluated material that filesystem resolution
cannot serve. Malformed glob syntax does not make a glob. An absolute selector
or one that lexically escapes the workspace remains filesystem intent and fails
containment rather than falling back to prose. An empty glob result, unreadable
existing path, or deterministic collection failure remains a filesystem-source
failure. Tools resolve path and glob selectors deterministically and MAY resolve
prose selectors through their own gathering means; a tool that cannot resolve a
selector's kind MUST report that inability distinctly, rather than treating the
named material as missing or evaluating the area against no material.

### Requirement scope

A requirement is scoped to the area on which it is declared: it states an
expectation about that area's source, and only that source.

Child areas do not inherit ancestor requirements. Declaring a child area
narrows attention to that area's source; it does not re-apply, copy, or extend
any ancestor expectation to the child.

An ancestor area's source can overlap a descendant area's source. Overlap does
not reassign, merge, or deduplicate anything: the ancestor's requirements still
read against the ancestor's source, and the descendant's requirements read
against the descendant's source. The same material can therefore be subject to
requirements declared on more than one area, each read at its own declaring
area's scope.

### Factor connection

Factor connection relates requirements to the factors that describe them, on
one area. A requirement declared under a factor or sub-factor is connected to
that factor by placement; the containing factor is its primary factor. Entries
in a requirement's `factors` list connect it to further factors declared on the
same area as same-area secondary references.

A factor's meaning spans, together:

- requirements declared under that factor or its sub-factors, and
- requirements declared on the same area that reference it explicitly.

Factor connection is a many-to-many relation over a single set of requirements.
A requirement connected to several factors remains one requirement with one
assessment; connection describes it from several quality perspectives, it does
not multiply it.

### Rating scale semantics

The rating scale's sequence order is part of its meaning: levels are declared
from best to worst.

A rating level's `description` and `criterion` have distinct semantics. The
`description` defines what the level means across the model. The `criterion`
defines the default rule for deciding whether a requirement sits at that
level.

A requirement's `ratings` entry overrides only that level's criterion for that
requirement. An override MUST NOT be read as altering a level's `description`,
`title`, `level`, or position in the scale order.

Rating levels have no meaning beyond what the model declares. Level names such
as `target`, `pass`, or `fail` carry no built-in semantics; an interpretation
MUST NOT assume fixed rating level meanings the declared scale does not state.

## Body semantics

The Markdown body documents context for building, interpreting, using, and
evaluating the model. The format does not require any body section names,
ordering, or content.

The body can document the root area, scope, stakeholder needs, risks, unknowns,
open questions, evidence context, or other important context for anyone
interpreting or applying the model.

## Extensions

A document MAY include frontmatter properties beyond those defined in
[Frontmatter schema](#frontmatter-schema). Extension properties MUST NOT change
the meaning of properties defined in this specification and should use names
unlikely to conflict with future versions of this specification.

Note: This note is non-normative. Tools that rewrite QUALITY.md documents are
expected to preserve body content and extension properties they do not
interpret; dropping or reinterpreting them changes a document the tool did not
author. Tools can freely build additional capability on top of the format —
filters, output formats, aggregation methods, authoring aids — as long as they
read the document according to [Model semantics](#model-semantics).

## Appendix A: Suggested rating scale

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

## Appendix B: Minimal example

This appendix is non-normative.

This is an illustrative software product example, not a default domain or
factor set for QUALITY.md. The same model shape applies across domains; only the
domain-carried factors, requirements, and assessments change. For a worked
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

## Appendix C: Invalid counter-examples

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

Invalid because every rating level MUST declare `level`, `title`, and
`criterion`.

### Direct area requirement without factors

```yaml
requirements:
  checkout-requests-durable:
    title: Checkout requests are durably recorded
    assessment: Review production write-path telemetry and recovery tests.
```

Invalid because a requirement declared directly under an area MUST declare
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

## Appendix D: An illustrative evaluation method

This appendix is non-normative.

QUALITY.md does not prescribe an evaluation method. Any process that reads a
model according to [Model semantics](#model-semantics) can evaluate it: a team
working through requirements in a review meeting, a scheduled audit, a CI
check, or an AI agent. This appendix sketches one concrete method — the one
used by the `/quality` agent skill — to illustrate how the declared semantics
play out in practice. Nothing in this outline adds conformance requirements.

The skill evaluates a model in phases:

1. **Scope.** The evaluation covers the whole model by default, or is narrowed
   by area or factor. A narrowed evaluation is recorded and reported as scoped,
   never presented as a whole-model verdict. Each in-scope area's source is
   resolved per [Source resolution](#source-resolution).
2. **Assess.** Each in-scope requirement is assessed once, against the source
   of its declaring area, by applying its `assessment`. This produces
   evidence-backed findings — gaps, risks, strengths, and notes — along with
   any unknowns and limits the evaluator hit.
3. **Rate.** Each requirement's findings are judged against each rating level's
   criterion, with the requirement's `ratings` overrides applied. The outcome
   is the one level whose criterion the findings satisfy, or `not assessed`
   when the findings cannot responsibly distinguish the criteria — never a
   guessed level.
4. **Analyze.** Ratings roll up bottom-up: each factor is characterized from
   its connected requirements and sub-factors, and each area from its own
   requirements and its child areas. Every roll-up carries a rationale naming
   what drove it, and a roll-up over too little assessed material stays
   `not assessed`.
5. **Advise.** Findings and ratings are distilled into recommendations for
   improvement or for review of the quality bar itself. Advice never changes a
   rating.
6. **Report.** The report presents the scope, each requirement's findings and
   rating with rationale, the factor and area roll-ups, and the advice, with
   `not assessed` outcomes shown distinctly from rating levels throughout.

`not assessed` is this method's own outcome convention, not part of the
document format: the skill reserves it so that unjudged material stays visibly
distinct from rated material instead of being guessed into a rating level.

Other methods are equally legitimate. A method might rate only requirements and
skip roll-ups, apply numeric scoring, or re-evaluate a single area on a
schedule. What makes any of them an evaluation of the same model is that it
reads each requirement against its declaring area's source, honors the declared
criteria and scale order, and asserts no meanings the model does not declare.
