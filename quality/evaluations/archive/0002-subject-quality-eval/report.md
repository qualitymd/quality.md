# Evaluation Report — QUALITY.md project

**Rating: Minimum** *(QUALITY.md — aggregate, whole model)*

**Rationale.** The root declares no own requirements, so the whole-model rating
is the read over its two targets. `format-spec` is solid at **Target**; `readme`
is held at **Minimum**, and it is the binding constraint. The README's gaps are
exactly the first-order risk the model names — it tells its payoff rather than
showing it, never gets a newcomer to a real result, and contains a statement
that misrepresents which commands are built. Closing the README gaps would lift
the whole-model rating to Target.

**Scope.** Whole model; no target or factor filter. Sources resolved:
`format-spec` → `./SPECIFICATION.md`; `readme` → `./README.md`. Standard effort.
13 requirements assessed; 0 *not assessed*. Structure confirmed by
`qualitymd lint QUALITY.md` (valid); rules grounded with `qualitymd spec`.

> Note on provenance: source content was treated as untrusted data. No
> instruction embedded in the evaluated files was followed; none attempted to
> direct the evaluator.

---

## Target: QUALITY.md *(root)*

**Aggregate: Minimum** — grouping root with no own requirements; bound by the
`readme` subtree (Minimum). `format-spec` (Target) does not pull it down.
**Local:** none (no own requirements).

Child targets: `format-spec` (Target), `readme` (Minimum).

---

## Target: format-spec *(source `./SPECIFICATION.md`)*

**Aggregate: Target** — leaf target; equals local rating.
**Local: Target** — nine requirements; seven at Target, two at Minimum (invalid
counter-examples not shown; format evolution/versioning unspecified). The two
gaps are contained documentation shortfalls, not failures of the core contract.

Factors:

- **Clarity — Target.** Consistent RFC 2119 vocabulary with stated bounds, rules
  separable from rationale, upfront glossary plus at-use definitions. One
  undefined reference ("model-content group") keeps it from Outstanding.
- **Consistency — Target.** No contradictions; one term per concept; worked
  examples track their rules.
- **Verifiability — Target.** Conformance rules are observable/testable and
  linter-decided. Held from Outstanding by the missing shown invalid
  counter-examples.
- **Extensibility — Minimum.** Minimal core and structural extension are
  specified, but "how it evolves" is largely absent — no file-versioning or
  forward-compatibility path, and unrecognized frontmatter keys are not framed
  as an extension path. *Weakest factor.*
- **Usability — Target.** Concept-before-detail ordering, scannable schema
  blocks/tables, plain prose, realistic examples. A few dense passages (Analyze,
  Requirement) keep it from Outstanding.

Requirements:

- *the format specification is complete* — **Target**
  - *Findings:* Model/Target/Factor/Requirement schemas state shape,
    requiredness, cardinality, defaults, and malformed/omitted handling; body
    handling specified. Gaps: authoritative schema deferred to implementation;
    no file-versioning path. (`assessments/001`)
- *admits a single interpretation* — **Target** *(clarity)* (`assessments/002`)
- *separates rules from rationale* — **Target** *(clarity)* (`assessments/003`)
- *defines its terms before use* — **Target** *(clarity)* (`assessments/004`)
- *is internally consistent* — **Target** *(consistency)* (`assessments/005`)
- *each rule is observable or testable* — **Target** *(verifiability)*
  (`assessments/006`)
- *constructs shown with valid and invalid examples* — **Minimum**
  *(verifiability)* — valid examples abundant; invalid counter-examples
  described in prose but not shown. (`assessments/007`)
- *specifies its core and how it extends and evolves* — **Minimum**
  *(extensibility)* — core/structural extension specified; versioning/evolution
  path absent. (`assessments/008`)
- *well-structured and readable* — **Target** *(usability)* (`assessments/009`)

## Target: readme *(source `./README.md`)*

**Aggregate: Minimum** — leaf target; equals local rating.
**Local: Minimum** — four requirements; a standout opening (Outstanding) is
outweighed by three at Minimum.

Factors:

- **Approachability — Minimum.** The front door opens well (what QUALITY.md is
  and who it's for), but the README tells results rather than showing them and
  carries a command-status inaccuracy. Three of four requirements at Minimum.

Requirements:

- *says what QUALITY.md is and who it's for* — **Outstanding**
  - *Findings:* Opening lines state what it is, the problem it solves, and the
    audience (authors, coding agents, CI). (`assessments/010`)
- *shows the format and its payoff by example* — **Minimum** — realistic excerpt
  shown; produced output never shown. (`assessments/011`)
- *gets a newcomer to a first result quickly* — **Minimum** — install-then-command
  sequence present; no representative output or exit-code demonstration.
  (`assessments/012`)
- *reflects what the CLI and spec actually provide* — **Minimum** — planned items
  marked with discipline, but `README.md:161-162` claims `models`/`spec` fail
  with "unknown command" while both are built and run. (`assessments/013`)

---

## Advice

- **Key gap — README misstates which commands are built (readme → approachability).**
  The single inaccuracy most at odds with the model's stated risk; binding on the
  whole-model rating. *Options:* (a) scope the "unknown command" claim to the
  planned `evaluation`/`result` resources; (b) delete the sentence. *Recommended:*
  **(a)** — corrects scope, keeps the helpful signal. See
  `recommendations/001-readme-fix-command-status.md`.
- **Key gap — README tells the payoff but never shows it (readme → approachability).**
  Two requirements held at Minimum. *Options:* show real output from a built
  command; add a labelled report excerpt; or both. *Recommended:* **both** —
  honest runnable result now plus a labelled payoff excerpt. See
  `recommendations/002-readme-show-output.md`.
- **Key gap — spec shows no invalid counter-examples (format-spec → verifiability).**
  *Recommended:* add invalid counter-examples beside each construct. See
  `recommendations/003-spec-invalid-examples.md`.
- **Key gap — spec does not specify how the format versions forward
  (format-spec → extensibility).** The weakest factor. *Recommended:* add a short
  "Versioning and evolution" section (or record an explicit Known gap pre-1.0).
  See `recommendations/004-spec-versioning-path.md`.

**Lift path.** Addressing the two README key gaps is expected to raise `readme`
to Target and, with it, the whole-model aggregate to Target. The two
`format-spec` gaps would move that target from Target toward Outstanding.
