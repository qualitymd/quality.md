---
ratings:
  pass:
    displayName: "Pass"
    description: "The model fully satisfies the diagnostic requirement."
  fail:
    displayName: "Fail"
    description: "The model does not satisfy the diagnostic requirement."
factors:
  functionality:
    requirements:
      "model fulfills its declared functional purpose":
        prompt: >
          Taken as a whole, the model serves the purpose it declares. Evaluating
          it yields a meaningful verdict about whether the subject meets the
          quality expectations the model sets, supporting a real accept, reject,
          or improvement decision — not merely confirming that the file is
          well-formed.
      "each factor is individually well-formed":
        prompt: >
          Each factor the model declares exhibits these individual-factor
          characteristics:
          - a quality attribute: it names a dimension of the subject's quality —
            what good looks like along one axis — not a component, feature, or
            activity of the subject.
          - relevant: it is selected for the subject's declared needs and context,
            not a generic catalog entry with no bearing on this subject.
          - distinct: it is scoped clearly enough that it is unambiguous which
            quality concern it covers and how it differs from its siblings.
          - grounded: it is drawn from a recognized, established quality vocabulary
            where one fits, tailored to the subject, rather than coined
            arbitrarily.
          - operationalized: it carries requirements — directly or through its
            sub-factors — that genuinely assess it; it is not left as a vague
            heading whose failure no requirement would surface.
          - described: the Markdown body explains what the attribute means for
            this subject, how one would know it is met, and the trade-offs it
            carries against other factors.
      "the factor set is well-formed as a whole":
        prompt: >
          The full set of factors across the model exhibits these set-level
          characteristics:
          - complete: the factors together cover the quality concerns the
            subject's declared needs and risks imply; concerns intentionally left
            out are recorded as explicit known gaps rather than silent omissions.
          - non-overlapping: factors are mutually distinct and do not
            substantially overlap; the set partitions the subject's quality
            cleanly rather than double-counting one concern across several factors.
          - coherent: the set reflects a recognized, established quality model
            tailored to the subject rather than an arbitrary or generic list.
          - appropriately decomposed: nesting depth suits the subject —
            sub-factors break an attribute too broad to assess directly into
            assessable parts, rather than being added gratuitously.
      "requirements collectively realize the model's declared needs":
        prompt: >
          The requirements, taken together, address every need and risk the
          model declares as its purpose. No declared need is left without at
          least one requirement whose failure would surface a violation of it.
      "the model yields correct, trustworthy verdicts":
        prompt: >
          For a subject that does or does not satisfy a requirement, the
          model's assessments and rating scale produce the verdict a
          knowledgeable reviewer would reach. The model does not systematically
          pass deficient subjects or fail adequate ones.
      "each requirement is individually well-formed":
        prompt: >
          Each requirement the model declares exhibits these
          individual-requirement characteristics:
          - necessary: it defines an essential capability, characteristic,
            constraint, or quality factor whose removal would leave a deficiency
            no other requirement covers; obsolete or moot requirements are not
            present.
          - appropriate: its intent and level of detail suit the entity it
            applies to and avoid unnecessary constraints on the subject's
            architecture or design, with supporting detail (rationale,
            thresholds, method) in the assessment fields and Markdown body rather
            than baked into the statement.
          - unambiguous: it can be interpreted in only one way, simply and easy
            to understand.
          - complete: it describes its expected capability, characteristic,
            constraint, or quality factor well enough to be understood on its own.
          - singular: it states a single concern and is captured by a single
            assessment — one prompt or one bash command, never several and never
            a list. Multiple conditions under which that concern must hold are
            acceptable; bundling distinct concerns into one requirement is not.
          - feasible: it can be satisfied by the subject within its constraints
            (cost, schedule, technical) at acceptable risk.
          - verifiable: it is worded so its satisfaction can be proven and is
            paired with an assessment method that delivers that proof — bash when
            the verdict is deterministic and cheaply computable, prompt when it
            needs judgment; measurable, observable expectations are preferred.
          - correct: it accurately represents the stakeholder need or risk it was
            derived from.
          - conforming: it follows a consistent template and style for
            requirements in this model, where an applicable convention exists.
      "the requirement set is well-formed as a whole":
        prompt: >
          The full collection of requirements across the model's factors exhibits
          these set-level characteristics:
          - complete: the set stands alone, sufficiently describing the quality
            factors needed to meet the subject's needs without further
            information, with no unresolved TBD/TBS/TBR placeholders; concerns
            intentionally left out are recorded as explicit known gaps rather
            than silent omissions.
          - consistent: requirements are unique and do not conflict with or
            overlap one another, and terminology, units, and measurement are used
            consistently throughout, with the same term meaning the same thing
            across the set.
          - feasible: the complete set can be satisfied by the subject within its
            constraints (cost, schedule, technical) at acceptable risk.
          - comprehensible: the set makes clear what is expected of the subject
            and how it relates to the system it is part of, with enough rationale
            in the Markdown body for a human or agent to understand why these
            factors and requirements are the right ones.
          - able to be validated: it is practicable that satisfying the set would
            achieve the subject's needs within constraints — a subject that passes
            every requirement is genuinely good enough for the needs the model
            declares.
    factors:
      correctness:
        requirements:
          "the model passes structural lint":
            bash: qualitymd lint QUALITY.md
          "the model correctly applies the QUALITY.md format spec":
            prompt: >
              The model is a correct application of the QUALITY.md format as the
              spec defines it, not merely a file that parses. Factors are declared
              under `factors`, each carrying requirements and/or sub-factors;
              every requirement names a single assessment — one `prompt` or one
              `bash`, never several; `target` and `prompt` references point at
              artifacts that exist and that the requirement is genuinely meant to
              judge; and the `ratings` scale is well-shaped and ordered best to
              worst. Each format construct is used for its intended purpose: the
              assessment method suits the nature of the check (deterministic and
              cheaply computable → `bash`, judgment → `prompt`), `target` scopes
              the right artifact, and graded expectations live in the rating scale
              rather than being baked into requirement statements. Where the model
              departs from the spec it does so through the format's own extension
              points, not by misusing or contradicting a defined construct.
      completeness:
        requirements:
          "the model includes everything the QUALITY.md spec prescribes":
            prompt: >
              The model includes everything the QUALITY.md spec prescribes, with
              nothing the format calls for left out. The required `factors`
              frontmatter is present and non-empty; the Markdown body carries the
              recommended spine — at minimum Overview, Needs, and Factors — and
              every factor and sub-factor declared in the frontmatter has matching
              body prose explaining what it means for this subject. The structured
              model and its human-readable rationale are both present, so the file
              delivers the whole of what the format intends rather than only its
              machine-readable half.
  usability:
    factors:
      agent usability:
        requirements:
          "an agent can interpret each requirement and assessment unambiguously":
            prompt: >
              A coding agent reading the model arrives at a single reading of each
              requirement and its assessment. Statements and prompts are
              self-contained — they supply the scope, terms, and instruction an
              agent needs to render the verdict — so that independent agent
              evaluations converge rather than diverging on interpretation.
              Nothing essential to understanding what is being asked is left
              implicit or deferred to context the model does not provide.
          "an agent can execute every assessment the model declares":
            prompt: >
              Each assessment can be carried out by an agent as written. Every
              `bash` command is runnable in the subject's environment and its
              result maps cleanly to a verdict; every `prompt` gives the agent the
              instruction and scope it needs, and any `target` or referenced
              artifact it must read can be located and retrieved. No assessment
              requires a capability, input, or access the agent cannot obtain.
          "an agent can use the model as working context":
            prompt: >
              Beyond formal evaluation, an agent can read the model and build an
              accurate picture of the subject's quality expectations to guide work
              on the subject. The frontmatter is navigable and the body explains
              intent, so the model functions as actionable guidance an agent can
              apply while building or changing the subject, not only as a checklist
              run at the end.
      developer usability:
        requirements:
          "a developer can understand the model's intent":
            prompt: >
              A developer encountering the model can grasp what it expects and why
              from its Markdown body without reverse-engineering the frontmatter.
              The body conveys the purpose behind the chosen factors and
              requirements and what a passing or failing verdict means for the
              subject, so the model is learnable by a human who did not write it.
          "a developer can act on the model's verdicts":
            prompt: >
              When the model reports a result, a developer can tell what to do
              about it. A failing requirement points clearly enough at what is
              deficient, and where, that the developer knows what to change;
              verdicts are concrete and traceable to the factor and requirement
              that produced them, rather than opaque or generic.
          "a developer can extend and maintain the model":
            prompt: >
              A developer can amend the model — adding, refining, or removing
              factors and requirements — through the format's own extension points
              without restructuring what is already there. Naming, layout, and
              assessment conventions are consistent enough that a maintainer can
              follow the established pattern and keep the model coherent as it
              evolves.
---

# Quality meta model

## Overview

This is the CLI's built-in diagnostic model for evaluating a project's
`QUALITY.md` file. The file is parsed with the same schema as a normal quality
model, but it is not a separate public root-file convention. Its subject is a
quality model: the project's `QUALITY.md`, its referenced prompts or standards,
and the project code used as context for coverage-gap findings.

The model's **factors** are product-quality attributes of the quality model
treated as a working artifact, beginning with **Functionality**: does the model
actually fulfill its purpose. The well-formedness of the model's own **factors**
and **requirements** is expressed as two requirement pairs *within* Functionality
rather than as factors of their own: for each, one requirement rolls up the
characteristics of an *individual* factor or requirement and the other the
characteristics of the *set* as a whole. Well-formed factors and requirements are
part of what makes the model functional, so they belong under Functionality
rather than beside it.

Alongside Functionality the model carries a **Usability** factor that judges
whether the model can actually be used, decomposed into **Agent usability** and
**Developer usability** for its two classes of user: the coding agents that run
it and read it as context, and the developers who author, maintain, and act on
it. A model can be functionally sound yet hard to use; this factor keeps that
concern explicit.

## Needs

- A project can tell whether its `QUALITY.md` is a useful quality model rather
  than only a syntactically valid file.
- An agent can produce concrete defects and coverage gaps that improve the
  model before using it to evaluate the subject.
- The diagnostic criteria stay inspectable and versioned instead of being an
  opaque prompt hidden inside the evaluator.

## Factors

### Functionality

Does the model do its job? A quality model exists to govern its subject's
quality and produce a verdict that drives a decision. Functionality asks whether
the model, run as a whole, fulfills that purpose: it frames quality with factors
that suit the subject — distinct, clearly scoped attributes drawn from a
recognized, established quality model and tailored to the subject's needs rather
than chosen ad hoc, each operationalized by the requirements beneath it; its
requirements collectively realize the needs it declares (functional
completeness); and its assessments and rating scale yield verdicts a
knowledgeable reviewer would trust (functional correctness), so a result means
something about the subject rather than only about the file.

Two Functionality requirements cover the well-formedness of the model's factors —
the dimensions along which it frames quality:

- *each factor is individually well-formed* rolls up the characteristics of an
  individual factor — a quality attribute (not a component or activity), relevant,
  distinct, grounded, operationalized, described — applied to each factor the
  model declares.
- *the factor set is well-formed as a whole* rolls up the set-level
  characteristics — complete coverage, non-overlapping, coherent, appropriately
  decomposed — applied to the full set of factors, catching gaps and overlaps that
  no single-factor check would surface.

Two further Functionality requirements cover the well-formedness of the model's
own requirements:

- *each requirement is individually well-formed* rolls up the characteristics of
  an individual requirement — necessary, appropriate, unambiguous, complete,
  singular, feasible, verifiable, correct, conforming — applied to each
  requirement the model declares.
- *the requirement set is well-formed as a whole* rolls up the characteristics of
  a set of requirements — complete, consistent, feasible, comprehensible, able to
  be validated — applied to the full collection of requirements across the
  model's factors, catching defects (gaps, conflicts, inconsistent terminology,
  an unvalidatable whole) that no single-requirement check would surface.

Both pairs live under Functionality rather than as standalone factors because
well-formed factors and requirements are a precondition for the model
functioning, not a separate product-quality attribute of it.

Functionality also decomposes into two sub-factors that judge the model's
conformance to the QUALITY.md spec itself — its *format*, not its subject. Where
the Functionality requirements above ask whether the model serves its subject and
needs, these ask whether the model is a correct and complete instance of the
QUALITY.md format. They are the meta-model's expression of functional correctness
and functional completeness measured against the spec rather than the subject.

#### Correctness

Does the model use the QUALITY.md format correctly? The deterministic floor is
structural lint (`qualitymd lint`), captured as its own `bash` requirement: the
file must parse and satisfy the structural schema before any judgment is worth
making. Above that floor, correctness asks whether the model applies the format
as the spec defines it — single assessments, resolvable `target` and `prompt`
references, a well-ordered rating scale — and whether each construct is used for
its intended purpose rather than bent to another (method matched to the kind of
check, `target` scoping the right artifact, graded expectations carried by the
rating scale). A correct model is a faithful application of the format, not
merely a file that lints clean.

#### Completeness

Does the model include everything the spec prescribes? Completeness asks whether
the required `factors` frontmatter and the recommended Markdown spine — Overview,
Needs, Factors, and matching prose for every declared factor and sub-factor — are
all present, so the file carries both the machine-readable model and the
human-readable rationale the format intends, not just one half.

### Usability

Can the model be used? Functionality asks whether the model frames and judges
quality correctly; Usability asks whether the people and agents who must work
with it can actually do so. Its two classes of user have different needs, so it
decomposes into two sub-factors — **Agent usability** and **Developer
usability** — each judging the same artifact from the standpoint of one user. A
model can be functionally sound yet leave a user unable to act on it without
guesswork; Usability keeps that concern explicit.

#### Agent usability

Can a coding agent use the model? The model is consumed by agents two ways: as
the artifact an evaluator agent runs to score the subject, and as context an
agent reads while building or changing the subject. Agent usability asks whether
an agent can interpret, execute, and apply it on its own — that each requirement
has a single clear reading, that every assessment is runnable with the inputs and
references it needs, and that the model doubles as actionable guidance an agent
can use mid-task rather than only a checklist scored at the end.

#### Developer usability

Can a human developer use the model? Developers author the model, review it,
maintain it as the subject evolves, and act on its verdicts. Developer usability
asks whether the model is learnable — its intent and rationale legible from the
Markdown body without reverse-engineering the frontmatter — whether its verdicts
are actionable, pointing clearly at what to change, and whether a maintainer can
extend or amend it through the format's extension points while keeping it
coherent. It is the human counterpart to Agent usability: the same artifact, read
and worked on by a person rather than an agent.
