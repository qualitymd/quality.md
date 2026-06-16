# QUALITY.md Specification

This is the specification for the `QUALITY.md` standard: a file format and set of conventions to help humans and agents model, evaluate, and improve quality. This specification uses terminology from the software development context, but the `QUALITY.md` standard can work in any operational context.

### Conformance

Conforming tools and agents must fulfill all normative requirements. Conformance requirements are described in this document via both descriptive assertions and key words with clearly defined meanings.

The key words “MUST”, “MUST NOT”, “REQUIRED”, “SHALL”, “SHALL NOT”, “SHOULD”, “SHOULD NOT”, “RECOMMENDED”, “MAY”, and “OPTIONAL” in the normative portions of this document are to be interpreted as described in IETF RFC 2119. These key words may appear in lowercase and still retain their meaning unless explicitly declared as non-normative.

A conforming use or application of QUALITY.md may provide additional functionality, but must not where explicitly disallowed or would otherwise result in non-conformance.

### Key Terms

**Quality Model**: a structured, declarative description of what
quality means for a given entity.
**Entity** a thing that is evaluated for quality.
**Factor** a quality (sub)characteristic or attribute for describing the quality of an entity.
**Requirement**: a quality requirement for assesing and rating the quality of an entity.
**Finding**: A single observation produced by assessing the source entities against a requirement — a unit of evidence such as a measured value, an inspection note, or a diagnostic result. A finding records *what was observed* and is not itself rated; the **findings** of a requirement are rated together.
**Assessment**: The means for assessing an entity - measurement, specifications, inspection, checklists, diagnostics, etc.
**Rating Scale**: A defined set of rating levels for a quality model
**Rating Level**: A single level on a rating scale, providing the default criterion for rating a requirement's findings
**Rating Result**: The outcome of rating a requirement's findings against the rating scale — a single rating level (or a *not assessed* outcome) assigned to the requirement, considering all of its findings together.
**Target**: An entity or set of entities with quality requirements subject to evaluatioe
**Source**: The scope of entities defined by a target

## QUALITY.md File

A `QUALITY.md` file is a markdown file with YAML frontmatter with a structured quality model and a markdown body.

The presence of a `QUALITY.md` file in a directory MUST imply that the directory and all its sub-directories and their contents are the implied source of any quality evaluation for the contained quality model (unless the optional root `source` defines otherwise.)

TODO (move these to evaluation)
The presence of a `QUALITY.md` within sub-directory
parent quality mds

### YAML Frontmatter

`QUALITY.md` files MUST begin with a valid YAML frontmatter block containing the required **Model** properties specified below.

When authoring `QUALITY.md` front-matter, null or empty optional properties SHOULD be omitted.

#### Model

The model represents a quality model: what things (**Targets**) are evaluated for quality, their important quality characteristics (**Factors**), measurable quality **Requirements**, assessment criteria, and **Rating** criteria for determining the level of quality.

```yaml
title: <string>                   # Required, the title of the entity whose quality is being modeled
ratings:                    # Required, the rating scale 
- level: <level-name>               # required; 
  title: <string>                   # optional human label
  criterion: <string>               # required; used to rate requirement assessment findings
factors:                          # Optional*  
  <factor-name>: <Factor>
requirements:                     # Optional*
  <requirement-statement>: <Requirement>
targets: 
  <target-name>: <Target>         # Optional*
source: <string>                # Optional
```

*An entry on either factors, requirements, or targets MUST be supplied.

**Title**: This is the title of the entity whose quality is modeled. For software projects, this is typically the name of the product, system, or library.

**Rating Scale**: This is the rating scale that provides the default criterion for how requirement assessments should be judged to arrive at a rating level result. Each **Rating Level** contains a level name which MUST be unique within the scale and an optional **title** for improved readability. Rating levels MUST be ordered from best (first) to worst (last). At least two rating levels MUST be supplied.

**Factors**: quality characteristics or attributes that matter most for evaluating the overall quality of the entity.

**Requirements**: quality requirements that will be used to assess the quality of the entity. These are typically nested under a factor or target, but may be defined at the model root level when it is simpler to define a single requirement at the root and cross reference to multiple quality attributes.

**Targets**: more focused quality modeling for possible target entities. Not required but useful when a distinct set of factors or requirements would be more cohesively defined around a narrower target of evaluation than scope implied by the source of the entire quality model.

**Source**:

- Authors SHOULD specify factors for critical quality characteristics required to satisfy the needs of stakeholders factoring in their needs (e.g. Security, Safety)
- Authors SHOULD NOT specify factors that would be irrelevant for the purpose of the project (e.g. Maintainability for a spike/POC)
- Authors SHOULD NOT specify factors

#### Target

#### Factor

#### Requirement

#### Rating Scale

### Markdown Body

## Discovery

## Evaluation

Evaluation assesses a target's source entities against each requirement and rates the resulting evidence against the rating scale.

For every target:

1. Resolve the **source** entities to be evaluated from the target's `source`.
2. For each requirement:
   1. **Assess** the source entities using the requirement's **assessment**, producing the **findings** — the evidence for this requirement. Each finding records an observation (e.g., a measured value, inspection note, or diagnostic result) and is not itself rated.
   2. **Rate** the findings *together* against the rating scale criterion (or the requirement's criterion overrides), producing the requirement's **Rating Result**: a single rating level. When there are no findings or the evidence is insufficient to rate against the scale, the requirement MUST be recorded as **not assessed** rather than assigned a rating level.

> The steps above produce one Rating Result per requirement. Aggregating those results up the model tree (requirement → factor → target → root), including weighting, pass thresholds, and how *not assessed* propagates, is specified separately in **Roll-up** (TODO).
