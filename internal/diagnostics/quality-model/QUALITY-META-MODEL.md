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
      "the model's factors decompose the subject's quality appropriately":
        prompt: >
          The model organizes quality into factors (and subfactors) that suit
          the subject. Each factor names a distinct quality attribute that is
          relevant to the subject's declared needs and context, scoped clearly
          enough that it is unambiguous which quality concern it covers and how
          it differs from its siblings. The set as a whole reflects a recognized
          quality model (e.g., ISO/IEC 25010) tailored to the subject rather
          than an arbitrary or generic list, and the factors do not substantially
          overlap. Each factor is operationalized by requirements that genuinely
          assess it — no factor is left as a vague heading without requirements
          whose failure would surface a deficiency in it.
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
  requirement quality:
    requirements:
      "each requirement is necessary":
        prompt: >
          Every requirement defines an essential capability, characteristic,
          constraint, or quality factor: removing it would leave a deficiency no
          other requirement covers. Obsolete requirements, or ones made moot by
          the passage of time, are not present.
      "each requirement is appropriate":
        prompt: >
          Each requirement's intent and level of detail suit the entity it
          applies to, avoiding unnecessary constraints on the subject's
          architecture or design. Supporting detail (rationale, thresholds,
          method) lives in the assessment fields and Markdown body, not baked
          into the requirement statement.
      "each requirement is unambiguous":
        prompt: >
          Each requirement is stated so it can be interpreted in only one way,
          simply and easy to understand.
      "each requirement is complete":
        prompt: >
          Each requirement describes its expected capability, characteristic,
          constraint, or quality factor well enough to be understood on its own,
          without needing other information to interpret it.
      "each requirement is singular":
        prompt: >
          Each requirement states a single capability, characteristic,
          constraint, or quality factor. Multiple conditions under which that
          single concern must hold are acceptable; bundling several distinct
          concerns into one requirement is not.
      "each requirement is feasible":
        prompt: >
          Each requirement can be satisfied by the subject within its
          constraints (cost, schedule, technical) at acceptable risk.
      "each requirement is verifiable":
        prompt: >
          Each requirement is worded so its satisfaction can be proven for the
          subject, and is paired with an assessment method that delivers that
          proof — bash when the verdict is deterministic and cheaply computable,
          prompt when it needs judgment. Measurable, observable expectations are
          preferred.
      "each requirement is correct":
        prompt: >
          Each requirement is an accurate representation of the stakeholder need
          or risk it was derived from.
      "each requirement is conforming":
        prompt: >
          Each requirement follows a consistent template and style for
          requirements in this model, where an applicable convention exists.
  requirement set quality:
    requirements:
      "the requirement set is complete":
        prompt: >
          The set of requirements stands alone, sufficiently describing the
          quality factors needed to meet the subject's needs without further
          information, and contains no unresolved TBD/TBS/TBR placeholders.
          Quality concerns intentionally left out are recorded as explicit known
          gaps rather than silent omissions.
      "the requirement set is consistent":
        prompt: >
          Requirements in the set are unique and do not conflict with or overlap
          one another; terminology, units, and measurement are used consistently
          throughout, with the same term meaning the same thing across the set.
      "the requirement set is feasible":
        prompt: >
          The complete set of requirements can be satisfied by the subject
          within its constraints (cost, schedule, technical) at acceptable risk.
      "the requirement set is comprehensible":
        prompt: >
          The set makes clear what is expected of the subject and how it relates
          to the system the subject is part of. The Markdown body carries enough
          rationale for a human or agent to understand why these factors and
          requirements are the right ones.
      "the requirement set is able to be validated":
        prompt: >
          It is practicable that satisfying the requirement set would achieve the
          subject's needs within constraints — i.e., a subject that passes every
          requirement is genuinely good enough for the needs the model declares.
---

# Quality meta model

## Overview

This is the CLI's built-in diagnostic model for evaluating a project's
`QUALITY.md` file. The file is parsed with the same schema as a normal quality
model, but it is not a separate public root-file convention. Its subject is a
quality model: the project's `QUALITY.md`, its referenced prompts or standards,
and the project code used as context for coverage-gap findings.

The model is read at two levels. Its **factors** are product-quality attributes
of the quality model treated as a working artifact — in the ISO/IEC 25010 style
we use for any subject — beginning with **Functionality**: does the model
actually fulfill its purpose. Within that, the well-formedness of the model's
requirements is judged against the requirement-quality characteristics of
**ISO/IEC/IEEE 29148** (§5.2.5–5.2.6): the *individual*-requirement
characteristics become per-requirement checks (the **Requirement quality**
factor), and the *set* characteristics become checks on the requirement set as a
whole (the **Requirement set quality** factor). The 29148 characteristics are
expressed as requirements, not as factors of their own.

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
recognized quality model (ISO/IEC 25010 style) and tailored to the subject's
needs rather than chosen ad hoc, each operationalized by the requirements
beneath it; its requirements collectively realize the needs it declares
(functional completeness); and its assessments and rating scale yield verdicts a
knowledgeable reviewer would trust (functional correctness), so a result means
something about the subject rather than only about the file.

The factor check draws on the SQuaRE guidance for establishing quality
characteristics — they are *selected* for relevance to the subject's needs and
context (ISO/IEC 25010 §4.2; ISO/IEC 25030 §7.3.2, §7.4.2; ISO/IEC 25040
§5.3.3.2), grounded in a recognized model used as building blocks rather than an
arbitrary list (ISO/IEC 25030 Annex B), and decomposed so each attribute is
actually assessable rather than a vague heading (the "definitional gap" of
Dromey and Quamoco).

### Requirement quality

Are the model's *individual* requirements well-formed? This factor applies the
ISO/IEC/IEEE 29148 §5.2.5 characteristics of an individual requirement —
necessary, appropriate, unambiguous, complete, singular, feasible, verifiable,
correct, conforming — to each requirement the model declares. These are atomic
properties of a single requirement, which is why they are requirements of this
factor rather than factors in their own right.

### Requirement set quality

Is the model's requirement set sound *as a whole*? This factor applies the
ISO/IEC/IEEE 29148 §5.2.6 characteristics of a set of requirements — complete,
consistent, feasible, comprehensible, able to be validated — to the full
collection of requirements across the model's factors, catching defects (gaps,
conflicts, inconsistent terminology, an unvalidatable whole) that no
single-requirement check would surface.
