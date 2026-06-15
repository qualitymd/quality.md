---
ratings:
  - level: outstanding
    displayName: "Outstanding"
    promptCondition: "Exceeds the diagnostic requirement; satisfies it with margin to spare."
  - level: target
    displayName: "Target"
    promptCondition: "Satisfies the diagnostic requirement; no material gaps."
    bashCondition: "result.success"
  - level: minimum
    displayName: "Minimum"
    promptCondition: "Satisfies the core of the requirement but falls short of the goal; minor or scoped gaps remain at the acceptable floor."
  - level: unacceptable
    displayName: "Unacceptable"
    promptCondition: "Does not satisfy the diagnostic requirement; falls below the acceptable floor."
factors:
  functionality:
    factors:
      fitness for purpose:
        requirements:
          "model fulfills its declared functional purpose":
            prompt: >
              Taken as a whole, the model serves the purpose it declares. Evaluating
              it yields a meaningful verdict about whether the subject meets the
              quality expectations the model sets, supporting a real accept, reject,
              or improvement decision — not merely confirming that the file is
              well-formed.
          "requirements collectively realize the model's declared needs":
            prompt: >
              Every need and risk the model declares is realized by the
              requirements: each is guarded by at least one requirement whose
              failure would surface a violation of it, so no declared need is left
              unguarded. Conversely, no requirement stands without a need behind
              it. The set traces to the declared purpose in both directions — needs
              to requirements and back. (Whether the requirement set is internally
              sound — consistent, lean, free of placeholders — is judged under "the
              requirement set is well-formed as a whole", not here.)
          "the model yields correct, trustworthy verdicts":
            prompt: >
              For a subject that does or does not satisfy a requirement, the
              model's assessments and rating scale produce the verdict a
              knowledgeable reviewer would reach. The model does not systematically
              pass deficient subjects or fail adequate ones.
      model well-formedness:
        requirements:
          "each factor is individually well-formed":
            prompt: >
              Each factor the model declares exhibits these individual-factor
              characteristics:
              - a quality attribute: it names a dimension of the subject's quality —
                what good looks like along one axis — not a component, feature, or
                activity of the subject.
              - relevant: it is selected for the subject's declared needs and context,
                earning its place as a distinct dimension of quality — not a generic
                catalog entry with no bearing on this subject, and not an axis the
                subject's nature already settles, on which every plausible subject
                would land the same way.
              - distinct: it is scoped clearly enough that it is unambiguous which
                quality concern it covers and how it differs from its siblings.
              - grounded: it is drawn from a recognized, established quality vocabulary
                where one fits, tailored to the subject, rather than coined
                arbitrarily.
              - operationalized: it carries requirements — directly or through its
                sub-factors — that genuinely assess it; it is not left as a vague
                heading whose failure no requirement would surface.
              Rate on the scale: a single unmet characteristic caps the verdict at
              `minimum`; a factor that is not a quality attribute at all — a
              component, feature, or activity — is `unacceptable`. Localize the verdict:
              name each factor that falls short and the characteristic it misses,
              so the result points at a specific factor to fix rather than only
              that something is off.
          "the factor set is well-formed as a whole":
            prompt: >
              The full set of factors across the model exhibits these set-level
              characteristics:
              - complete: the factors together cover the quality concerns the
                subject's declared needs and risks imply; concerns intentionally left
                out are recorded explicitly — as out-of-scope in the Scope section, or
                as deferred in Known gaps — rather than silent omissions. A concern
                the model declares out of scope does not count against completeness.
              - non-overlapping: factors are mutually distinct and do not
                substantially overlap; the set partitions the subject's quality
                cleanly rather than double-counting one concern across several factors.
              - coherent: the set reflects a recognized, established quality model
                tailored to the subject rather than an arbitrary or generic list.
              - appropriately decomposed: nesting depth suits the subject —
                sub-factors break an attribute too broad to assess directly into
                assessable parts, rather than being added gratuitously.
              Rate on the scale: any unmet characteristic caps the verdict at
              `minimum`. Localize the verdict: name the specific gap, overlap, or
              incoherence in the set rather than only a bare rating.
          "each requirement is individually well-formed":
            prompt: >
              Each requirement the model declares exhibits these
              individual-requirement characteristics:
              - necessary: it defines an essential capability, characteristic,
                constraint, or quality factor whose removal would leave a deficiency
                no other requirement covers, and whose verdict is not a foregone
                conclusion — it discriminates sound subjects from deficient ones
                rather than restating an expectation the subject's context already
                guarantees or competent practice takes for granted. Obsolete, moot,
                or self-evident requirements are not present.
              - appropriate: its intent and level of detail suit the entity it
                applies to and avoid unnecessary constraints on the subject's
                architecture or design, with supporting detail (rationale,
                thresholds, method) in the assessment fields and Markdown body rather
                than baked into the statement.
              - unambiguous: it can be interpreted in only one way, simply and easy
                to understand.
              - complete: it describes its expected capability, characteristic,
                constraint, or quality factor well enough to be understood on its own.
              - concise: it says what it must in as few words as carry the meaning,
                spelling out neither what a competent reader already assumes nor
                common-sense qualifications that change no verdict. Brevity trims the
                self-evident, not the substantive — it serves clarity without costing
                completeness.
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
              Treat individual well-formedness as the single concern these
              characteristics jointly define — they are the conditions it must
              meet, not separate requirements. Rate on the scale: a single unmet
              characteristic caps the verdict at `minimum`; a requirement with no
              real assessment is `unacceptable`. Localize the verdict: name each
              requirement that falls short and the characteristic it misses, so a
              `minimum` says what to fix.
          "the requirement set is well-formed as a whole":
            prompt: >
              The full collection of requirements across the model's factors exhibits
              these set-level characteristics:
              - complete: the set stands alone — self-contained, free of unresolved
                TBD/TBS/TBR placeholders and of dangling references that would need
                further information to resolve; concerns intentionally left out are
                recorded explicitly — as out-of-scope in the Scope section, or as
                deferred in Known gaps — rather than silent omissions. A concern the
                model declares out of scope does not count against completeness.
                (Whether the requirements actually realize the model's declared
                needs is judged under Fitness for purpose, not here.)
              - consistent: requirements are unique and do not conflict with or
                overlap one another, and terminology, units, and measurement are used
                consistently throughout, with the same term meaning the same thing
                across the set.
              - lean: the set is no larger than the subject's needs and risks
                warrant; it does not multiply requirements past the point where each
                adds discriminating signal, so the verdict stays legible rather than
                diffused across redundant or self-evident checks.
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
              Rate on the scale: any unmet characteristic caps the verdict at
              `minimum`. Localize the verdict: name the specific conflict,
              redundancy, gap, or inconsistency rather than only a bare rating.
      documentation:
        requirements:
          "the Overview body section frames the subject":
            prompt: >
              Where present, the Overview section establishes the context the rest of
              the model depends on: what the subject is, who depends on it, and what
              "good" means for it. It frames the subject concretely enough that a
              reader can judge the model's prompt assessments consistently, rather
              than restating the factors or trailing off into generic description.
              The model's boundary — what it covers and deliberately excludes — is
              the Scope section's job, not the Overview's.
          "the Scope body section draws the model's boundary":
            prompt: >
              Where present, the Scope section states what the model covers and what
              it deliberately leaves out, including dependencies the subject relies on
              but does not own. Out-of-scope items are framed as exclusions by design
              — concerns outside the model's remit — not as deferred work that belongs
              here later; that distinction keeps Scope separate from Known gaps, and
              an item recorded as out of scope is not treated as a coverage gap. A
              subject whose boundary is unremarkable may omit the section; when
              present it draws the in/out line clearly rather than restating the
              Overview or duplicating Known gaps.
          "the Needs body section states what matters and to whom":
            prompt: >
              Where present, the Needs section states, in plain language, the outcomes
              stakeholders depend on the subject for — the progress each group is trying
              to make — and to whom each matters. Each need is expressed from the
              stakeholder's standpoint as a desired outcome, not as a feature, activity,
              want, or request, and not as a paraphrase of a requirement or a factor
              name; the test is consequence — for a genuine need you can name who
              suffers, and how, if it goes unmet or is met poorly. Together the needs
              give the requirements something concrete to realize. The section has
              considered the subject's full range
              of stakeholders rather than defaulting to end users alone — including,
              where they apply to this subject, the developers who build, review, and
              maintain it and the AI assistants or coding agents that build, operate,
              or consume it. These classes need not all appear, but a subject they
              plainly serve should not be silently overlooked.
          "the Risks body section states what failure costs and to whom":
            prompt: >
              Where the subject's needs carry material failure modes, the Risks
              section says what goes wrong, and for whom, if a need is not met, and
              conveys the relative severity that should shape the model's priorities —
              distinguishing the outcomes the model must guard against most from
              lesser, recoverable ones. A subject with no material risks may omit the
              section; when present it adds this consequence framing rather than
              restating the needs.
          "the Factors body section explains every declared factor":
            prompt: >
              The Factors section gives every factor and sub-factor declared in the
              frontmatter matching body prose that does its prescribed job: it
              characterizes the quality attribute itself — what the characteristic
              is, why it matters, and how it differs from its siblings — succinctly
              and at a level that would still hold if the factor's requirements or
              evaluated targets changed. It describes the attribute, not the
              requirements attached to it: it neither enumerates nor paraphrases
              them, and does not hard-code the particular targets they judge, so the
              same description could carry a different requirement set or apply to
              another subject unchanged. Explaining how an attribute is decomposed —
              its sub-factors, or why its requirements are grouped as they are — is
              part of this job; rehearsing what each individual requirement checks is
              not. The section mirrors the frontmatter's structure and reads as the
              rationale for why these are the right attributes — not a restatement of
              their names or their requirements.
          "the Known gaps body section records deferred concerns with reasons":
            prompt: >
              Where the model defers quality concerns that lie *inside* its scope —
              concerns it should answer but has not yet — the Known gaps section
              records each one with a brief reason, so the omission is explicit and
              intentional rather than a silent oversight. These are in-scope
              deferrals, distinct from the by-design exclusions the Scope section
              carries: a concern outside the model's remit belongs in Scope and is
              not a known gap. A model that addresses every in-scope concern may omit
              the section; when concerns are knowingly deferred, they appear here.
          "the Markdown body earns its length":
            prompt: >
              Across its sections, the Markdown body says what the model needs and
              no more. It supplies the subject-specific reasoning a reader could not
              supply themselves — what "good" means here, why these are the right
              factors and requirements — rather than general knowledge anyone
              familiar with this kind of subject already holds; it does not narrate
              the obvious or pad sections to look thorough. Each section stays short,
              preferring brevity to exhaustive documentation. Length is justified by
              content the reader lacks, not by a wish to appear complete.
      format conformance:
        requirements:
          "the model passes structural lint":
            bash: qualitymd lint
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
      format completeness:
        requirements:
          "the model includes everything the QUALITY.md spec prescribes":
            prompt: >
              The model includes everything the QUALITY.md spec prescribes, with
              nothing the format calls for left out. The required `factors`
              frontmatter is present and non-empty; the Markdown body carries the
              recommended spine — at minimum Overview, Needs, and Factors — and
              every factor and sub-factor declared in the frontmatter has a
              matching body subsection. The structured model and its
              human-readable rationale are both present, so the file delivers the
              whole of what the format intends rather than only its
              machine-readable half. Whether each present section is well-written
              is judged by the per-section body requirements, not here; this asks
              only that nothing prescribed is missing.
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
and the project code used as context for coverage-gap findings. Keeping these
criteria in a versioned, inspectable model — rather than an opaque prompt buried
in the evaluator — is deliberate: the bar a `QUALITY.md` is held to is itself
open to review and change.

The model carries two **factors**, the product-quality attributes of a quality
model treated as a working artifact: **Functionality** — does the model do its
job — and **Usability** — can the people and agents who must work with it
actually do so. A model can be functionally sound yet hard to use, so the two
are kept apart. Each requirement lands on the shared **outstanding / target /
minimum / unacceptable** scale, so a model with a minor, scoped gap is graded
`minimum` rather than forced onto a blunt pass-or-fail.

## Scope

The subject is the `QUALITY.md` artifact set out in the Overview. What this model
deliberately leaves **out** is the quality of the subject that `QUALITY.md`
*governs*: whether the underlying system is reliable, secure, or fast is the job
of evaluating that model against its subject, not of this one. The tools that run
an evaluation — the CLI and skills that execute and record assessments — are
likewise out of scope: this model relies on them but does not judge them. An item
placed out of scope here is excluded by design, not a coverage gap; in-scope
concerns deferred for now live under **Known gaps**.

## Needs

- A project can tell whether its `QUALITY.md` is a useful quality model rather
  than only a syntactically valid file.
- An agent can produce concrete defects and coverage gaps that improve the
  model before using it to evaluate the subject.

## Risks

This model governs whether a `QUALITY.md` can be trusted, so its own verdicts
carry cost:

- A **false acceptance** is the worst outcome — blessing a deficient quality model
  lets that model go on to govern its subject, so every evaluation run through
  it inherits the blind spot and the project gains false confidence. This is
  what the model must guard against most.
- A **false rejection** — flagging a sound model as deficient — is less
  damaging but erodes trust in the diagnostic and sends authors chasing
  non-defects.
- A **verdict that does not localize the defect** leaves a developer or agent
  unable to act, turning a failing result into a dead end rather than a next
  step.

## Factors

### Functionality

Does the model do its job? A quality model exists to govern its subject's
quality and produce a verdict that drives a decision. Functionality decomposes
into five sub-factors along two axes. The first asks whether the model serves
its **subject** — it is **fit for purpose**, its structured model is
**well-formed**, and its **documentation** does its job — so that a verdict
means something about the subject rather than only about the file. The second
asks whether the model is a faithful instance of the **format** it is written
in — **format conformance** and **format completeness** measured against the
QUALITY.md spec. The two axes mirror each other: each judges the model for
correctness and completeness, once against its subject and once against the
format.

#### Fitness for purpose

Does the model, taken as a whole, fulfill the purpose it declares? This is the
functional core, and it judges the *emergent whole* rather than the soundness of
the parts: a model can have well-formed factors, requirements, and prose and
still not hang together into something that drives a decision. The requirements
must collectively realize the needs the model sets out — every declared need
guarded by a requirement, and every requirement earning a need (functional
completeness) — and its assessments and rating scale must yield the verdict a
knowledgeable reviewer would reach (functional correctness), so that evaluating
the model drives a real accept, reject, or improve decision about the subject
rather than only confirming the file is well-formed. Whether each factor and
requirement is *individually* sound is judged under **Model well-formedness**,
not here.

#### Model well-formedness

Are the model's own factors and requirements soundly built? A model can serve
its purpose only if the dimensions along which it frames quality, and the
requirements beneath them, hold up on their own. Four requirements cover this —
two for the model's **factors** (the dimensions along which it frames quality)
and two for its **requirements**. Each pair splits the same way: one rolls up
the characteristics of an *individual* factor or requirement, the other the
characteristics of the *set as a whole* — the latter catching gaps, overlaps,
conflicts, and inconsistent terminology that no single-element check would
surface. The specific characteristics each one rolls up are spelled out in the
requirement prompts themselves. Because each requirement rolls many
characteristics into one verdict, each prompt grades on the shared scale and asks
the evaluator to localize the shortfall — naming the specific factor or
requirement and the characteristic it misses — so a `minimum` still says what to
fix rather than only that something is wrong. A structural breakdown — a factor
that is not a quality attribute, a requirement with no real assessment — is
`unacceptable`. (Whether a factor's *body prose* explains it is judged under
**Documentation**, not here.)

#### Documentation

Does the Markdown body do its prescribed job? The body is where quality is made
concrete: it frames what "good" means for the subject and supplies the grounding
a `prompt` assessment needs to be judged consistently, so a thin or generic body
undermines the model's functionality even when the frontmatter is sound. One
requirement per prescribed body section asks not whether the section is present
— **Format completeness** checks that — but whether, when present, it does the
job the spec assigns it; sections that apply only where the subject warrants
them — Scope, Risks, Known gaps — pass when appropriately absent. The **Scope**
and **Known gaps** requirements between them hold the in/out boundary: a concern
declared out of scope is judged as a by-design exclusion there, not a coverage
gap, which is what keeps the completeness checks above from flagging it. A final
requirement runs across the body
rather than singling out a section, asking that it earn its length — supplying
the subject-specific reasoning a reader lacks and no more, neither narrating the
obvious nor padding to appear thorough. What each section must do is defined in
its own requirement. These requirements own the quality of the body's prose,
which is why the well-formedness requirements above are scoped to the structured
model rather than its documentation.

#### Format conformance

Does the model use the QUALITY.md format correctly? The deterministic floor is
structural lint (`qualitymd lint`), captured as its own `bash` requirement: the
file must parse and satisfy the structural schema before any judgment is worth
making. Above that floor, conformance asks whether the model applies the format
as the spec defines it — single assessments, resolvable `target` and `prompt`
references, a well-ordered rating scale — and whether each construct is used for
its intended purpose rather than bent to another (method matched to the kind of
check, `target` scoping the right artifact, graded expectations carried by the
rating scale). A conforming model is a faithful application of the format, not
merely a file that lints clean.

#### Format completeness

Does the model include everything the spec prescribes? Format completeness asks
whether the required `factors` frontmatter and the recommended Markdown spine —
Overview, Needs, Factors, and matching prose for every declared factor and
sub-factor — are all present, so the file carries both the machine-readable
model and the human-readable rationale the format intends, not just one half.

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
can use mid-task rather than only a checklist scored at the end. Where
Functionality asks whether a requirement is *intrinsically* unambiguous and
verifiable, Agent usability asks whether *this agent, here* can act on it.

#### Developer usability

Can a human developer use the model? Developers author the model, review it,
maintain it as the subject evolves, and act on its verdicts. Developer usability
asks whether the model is learnable — its intent and rationale legible from the
Markdown body without reverse-engineering the frontmatter — whether its verdicts
are actionable, pointing clearly at what to change, and whether a maintainer can
extend or amend it through the format's extension points while keeping it
coherent. It is the human counterpart to Agent usability: the same artifact, read
and worked on by a person rather than an agent.

## Known gaps

- **The model's own maintainability is judged only through Developer usability,
  not as a standalone product-quality attribute.** A quality model is a living
  artifact, but at this stage its maintainability concerns — consistent
  conventions, ease of amendment — are covered well enough by *a developer can
  extend and maintain the model*; a separate factor would add structure without
  adding signal.
- **The structural-lint floor depends on an unbuilt command.** *The model passes
  structural lint* invokes `qualitymd lint`, which is specified but not yet
  implemented. Until it ships, that requirement cannot run and the structural
  floor rests on the format-conformance prompt and manual review.
