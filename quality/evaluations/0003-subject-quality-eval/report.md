# Evaluation Report

## Summary

- **Subject:** QUALITY.md
- **Altitude:** subject
- **Effort:** standard
- **Rating:** minimum
- **Rationale:** The model root declares no own requirements, so its aggregate reflects only its three children: format-spec aggregates to target, while readme and cli both aggregate to minimum. Two distinct binding constraints hold the whole model at the floor — the README tells rather than shows (no rendered qualitymd output, no quick first result), and the CLI has two confirmed MUST-level spec divergences in the automation surface (a machine-varying absolute path in the show-status payload and an empty-run gap kind that contradicts the documented missing-root-analysis contract). Both map directly to the model's stated Risks (an overstating/under-demonstrating README, and non-deterministic CLI behavior). Closing the README payoff/first-result gaps and the two CLI conformance gaps would lift the whole-model rating to target.

## Scope

standard` — every in-scope requirement assessed with targeted evidence; the rating-binding findings re-checked before reporting.

- **Narrowing:** whole recorded run
- **In scope:** QUALITY.md; cli; format-spec; readme
- **Out of scope:** Dependencies the project does not own (Go toolchain, Cobra/Fang, release tooling), per the model's Scope.; The `/quality` skill prompts and modes themselves (not modeled targets).; Exhaustive field-by-field diff of the README schema table against SPECIFICATION.md` (spot-checked only; recorded as a limitation).

## Top Risks and Limitations

- `assessments/001-format-spec-the-format-specification-is-complete.json` at `SPECIFICATION.md:112-130` [medium]: Gap: absent-frontmatter-entirely, multiple frontmatter blocks, and a valid-YAML-but-non-map top node are not explicitly handled.
- `assessments/001-format-spec-the-format-specification-is-complete.json` at `SPECIFICATION.md:137-152` [low]: Gap: scalar field types are abbreviated as <string>/<level-name> without stating allowed character sets, whitespace, or length bounds, leaving some edge cases undefined.
- `assessments/002-format-spec-the-format-specification-admits-a-single-interpretation.json` at `SPECIFICATION.md:345-347,384-385` [medium]: Gap: 'insufficient to rate' / 'too little has been assessed' qualifies the conformance-relevant not-assessed outcome (a MUST) with no stated bound, so different evaluators could draw the line differently.
- `assessments/002-format-spec-the-format-specification-admits-a-single-interpretation.json` at `SPECIFICATION.md:160-163` [low]: The two-level minimum is stated twice in adjacent sentences; consistent but a minor near-duplicate normative statement.
- `assessments/002-format-spec-the-format-specification-admits-a-single-interpretation.json` at `SPECIFICATION.md:387-388` [low]: 'most responsible' is a soft quantifier but sits under SHOULD (advisory), so it does not affect conformance and admits a reasonable single reading.
- `assessments/003-format-spec-the-format-specification-separates-rules-from-rationale.json` at `SPECIFICATION.md:252-256` [low]: Gap: a rationale sentence ('Referencing names that entity once instead of copying criteria that would drift') sits unmarked among normative statements, so it is normative by default though it reads as explanation.
- `assessments/004-format-spec-the-format-specification-defines-its-terms-before-use.json` at `SPECIFICATION.md:49-62` [low]: Gap: conformance-class terms (parser/linter/evaluator/report renderer) are used in the Conformance section before Terminology, though each is characterized operationally in place.
- `assessments/006-format-spec-each-rule-is-observable-or-testable.json` at `SPECIFICATION.md:160-161` [medium]: Gap: 'ordered from best to worst' is not machine-verifiable from the data alone; no comparable value encodes goodness, so it relies on author intent.
- Limitation: Two distinct binding constraints hold the whole model at the floor — the README tells rather than shows (no rendered qualitymd output, no quick first result), and the CLI has two confirmed MUST-level spec divergences in the automation surface (a machine-varying absolute path in the show-status payload and an empty-run gap kind that contradicts the documented missing-root-analysis contract)
- Limitation: Held below outstanding by edge gaps: undefined handling for absent/multiple frontmatter blocks and untyped scalar bounds
- Limitation: The show-don't-tell payoff is missing, holding the requirement at the floor
- Limitation: Most commands conform strongly (init, spec, lint, build-report gating and byte-stability), but two confirmed MUST-level spec divergences land squarely on the automation contract: a machine-varying absolute path in the show-status payload (determinism/consistency) and a gap kind that contradicts the documented missing-root-analysis contract (consistency)
- Limitation: The design-guide concerns (exit codes, conditional styling, non-interactivity, consistent grammar) are strongly met, but two confirmed MUST-level spec divergences in the automation surface — a machine-varying absolute path in the show-status payload and a gap kind contradicting the documented missing-root-analysis contract — bind the local rating to the floor
- Limitation: Accurate and clear on what/who, but the missing produced-output payoff and quick-first-result keep approachability at the floor
- Additional risks or limitations are available in `report.json`.

## Evidence Basis

- **command:** `grep -nEi 'report.md|report.json|rating result' README.md`
- **command:** `qualitymd init && qualitymd lint QUALITY.md (in temp dir)`
- **command:** `qualitymd --help`
- **command:** `qualitymd evaluation --help`
- **command:** `qualitymd evaluation build-report --help`
- **command:** `qualitymd evaluation create-run --json (relative) vs show-status --json (absolute)`
- **command:** `qualitymd evaluation build-report <empty-run>  # 'Run is not reportable: missing-analysis'`
- **command:** `NO_COLOR=1 qualitymd lint QUALITY.md | cat -v ; qualitymd lint QUALITY.md 2>/dev/null`

## Next Action

- [001-show-what-running-qualitymd-produces-not-just-the-model](recommendations/001-show-what-running-qualitymd-produces-not-just-the-model.md) - The requirement 'the README shows the format and its payoff by example' reaches at least target: the README displays representative qualitymd output produced from the shown example. Re-evaluate in a new numbered run.

## Target Summary

| Target      | Local            | Aggregate | Covered Requirements | Note                       |
| ----------- | ---------------- | --------- | -------------------- | -------------------------- |
| QUALITY.md  | n/a (structural) | minimum   | 0                    | structural grouping target |
| cli         | minimum          | minimum   | 2                    |                            |
| format-spec | target           | target    | 9                    |                            |
| readme      | minimum          | minimum   | 4                    |                            |

## Target Details

### QUALITY.md

- **Path:** (root)
- **Local rating:** n/a (structural)
  - Structural grouping target; local rating does not apply.
- **Aggregate rating:** minimum
  - The model root declares no own requirements, so its aggregate reflects only its three children: format-spec aggregates to target, while readme and cli both aggregate to minimum. Two distinct binding constraints hold the whole model at the floor — the README tells rather than shows (no rendered qualitymd output, no quick first result), and the CLI has two confirmed MUST-level spec divergences in the automation surface (a machine-varying absolute path in the show-status payload and an empty-run gap kind that contradicts the documented missing-root-analysis contract). Both map directly to the model's stated Risks (an overstating/under-demonstrating README, and non-deterministic CLI behavior). Closing the README payoff/first-result gaps and the two CLI conformance gaps would lift the whole-model rating to target.
- **Analysis record:** `analysis/quality-md.json`

### cli

- **Path:** cli
- **Local rating:** minimum
  - Two requirements: 'follows its functional specifications' at minimum and 'follows the project CLI design guide' at target. The design-guide concerns (exit codes, conditional styling, non-interactivity, consistent grammar) are strongly met, but two confirmed MUST-level spec divergences in the automation surface — a machine-varying absolute path in the show-status payload and a gap kind contradicting the documented missing-root-analysis contract — bind the local rating to the floor.
- **Aggregate rating:** minimum
  - Leaf target: aggregate equals local rating. Strong design-guide adherence and automation ergonomics are pulled to the floor by two confirmed spec-conformance divergences in determinism and consistency.
- **Factor consistency:** minimum
  - Command grammar, flags, and help are uniform (a strength), but payload paths are inconsistent across sibling commands (show-status absolute vs create-run relative) and the empty-run gap kind diverges from the documented contract.
- **Factor determinism:** minimum
  - build-report is byte-stable and ANSI-free (strong), but show-status --json emits a machine-varying absolute path, a determinism MUST violation in a machine-consumed payload field.
- **Factor automation-compatibility:** target
  - Strict JSON decoding, --json everywhere, non-interactive stdin handling, stdout/stderr separation, and exit-code categories all hold; held from outstanding by the duplicate human ERROR block on init --json refusals.
- **Factor usability:** target
  - Help, next-action footers, and error categories are clear and recoverable; a non-concise bare-args root help is the only low-severity deviation.
- **Analysis record:** `analysis/cli.json`

### format-spec

- **Path:** format-spec
- **Local rating:** target
  - Nine own requirements considered together: seven at target (completeness; all of clarity; consistency; the core verifiability rule; usability) and two at minimum (no shown invalid counter-examples; forward-evolution/versioning under-specified). The two gaps are real but contained documentation shortfalls, not failures of the spec's core contract, so the whole-set verdict holds at target.
- **Aggregate rating:** target
  - Leaf target: aggregate equals local rating. The spec satisfies its completeness, clarity, consistency, core-verifiability, and usability bars; extensibility is the weakest factor at minimum, alongside the unmet invalid-examples requirement.
- **Factor clarity:** target
  - All three clarity requirements meet target: BCP 14 keywords scoped to all-caps with an explicit normative split, rule/rationale separation declared and tagged, and a front-loaded glossary plus point-of-use definitions. Held from outstanding by a soft 'insufficient to rate' qualifier and one unmarked rationale aside.
- **Factor consistency:** target
  - No contradictions found; one term per concept and worked examples track the rules they illustrate.
- **Factor verifiability:** target
  - The binding question 'can conformance be decided' is answered yes: requirement-shape rules are statically checkable. Held from outstanding by judgment-bound rules (scale ordering, rating sufficiency, roll-up) and by the separate minimum-rated requirement that invalid counter-examples are not shown.
- **Factor extensibility:** minimum
  - Minimal core and structural extension are well specified, but 'how it evolves' is pointer-only (versioning policy delegated to an external doc) and unrecognized frontmatter keys lack an explicit reader rule. The format spec's weakest factor.
- **Factor usability:** target
  - Concept-before-detail ordering, scannable schema blocks, plain prose, and realistic appendices; minor structural nits (annotation-vs-table, heading hierarchy) keep it from outstanding.
- **Analysis record:** `analysis/format-spec.json`

### readme

- **Path:** readme
- **Local rating:** minimum
  - Four approachability requirements: two at target (the README says what QUALITY.md is and who it's for; it now accurately reflects what the CLI and spec provide, resolving the prior overstatement) and two at minimum (it does not show what running qualitymd produces, and it offers no quick install-then-one-command first result with representative output). The binding constraints are the two show-don't-tell gaps, holding the target at the floor.
- **Aggregate rating:** minimum
  - Leaf target: aggregate equals local rating. Accurate and clear on what/who, but the missing produced-output payoff and quick-first-result keep approachability at the floor.
- **Factor approachability:** minimum
  - The front door clearly states what QUALITY.md is and who it's for and is now factually accurate, but it tells rather than shows: there is no rendered qualitymd output and no copyable first-result quickstart, so a newcomer cannot see the payoff or reach a real result quickly. Held at minimum by those two gaps.
- **Analysis record:** `analysis/readme.json`

## Requirements

### the format specification is complete

- **State:** active
- **Target:** format-spec
- **Rating:** target
- **Assessment record:** `assessments/001-format-spec-the-format-specification-is-complete.json`
- **Rationale:** The frontmatter and body are specified completely enough to parse and author from, including the principal malformed/omitted handling, satisfying the requirement. Held below outstanding by edge gaps: undefined handling for absent/multiple frontmatter blocks and untyped scalar bounds.

### the format specification admits a single interpretation

- **State:** active
- **Target:** format-spec
- **Rating:** target
- **Assessment record:** `assessments/002-format-spec-the-format-specification-admits-a-single-interpretation.json`
- **Rationale:** Rules generally admit one reading with disciplined BCP 14 usage; the one conformance-relevant soft term ('insufficient to rate') is inherently an evaluator judgment. Satisfies the requirement; held from outstanding by that imprecision.

### the format specification separates rules from rationale

- **State:** active
- **Target:** format-spec
- **Rating:** target
- **Assessment record:** `assessments/003-format-spec-the-format-specification-separates-rules-from-rationale.json`
- **Rationale:** A default classification rule plus consistent explicit non-normative tagging let a reader tell rule from rationale, and rules are not hidden in examples. Held from outstanding by one unmarked rationale aside.

### the format specification defines its terms before use

- **State:** active
- **Target:** format-spec
- **Rating:** target
- **Assessment record:** `assessments/004-format-spec-the-format-specification-defines-its-terms-before-use.json`
- **Rationale:** Technical terms are defined before or at first use via a glossary plus point-of-use definitions, satisfying the requirement; a minor ordering wrinkle for the conformance-class terms keeps it from outstanding.

### the format specification is internally consistent

- **State:** active
- **Target:** format-spec
- **Rating:** target
- **Assessment record:** `assessments/005-format-spec-the-format-specification-is-internally-consistent.json`
- **Rationale:** No two statements contradict; one term denotes one concept throughout and worked examples track the rules they illustrate across a large document. Satisfies the requirement at target.

### each rule is observable or testable

- **State:** active
- **Target:** format-spec
- **Rating:** target
- **Assessment record:** `assessments/006-format-spec-each-rule-is-observable-or-testable.json`
- **Rationale:** The bulk of conformance rules map to something a reader can observe or linter-decide, so independent readers decide the same way. Held from outstanding by several rules (ordering, sufficiency, roll-up) that are judgment-bound rather than reproducibly testable.

### the format's constructs are shown with valid and invalid examples

- **State:** active
- **Target:** format-spec
- **Rating:** minimum
- **Assessment record:** `assessments/007-format-spec-the-format-s-constructs-are-shown-with-valid-and-invalid-examples.json`
- **Rationale:** Realistic valid examples hold the floor, but the requirement explicitly wants both valid and invalid examples: invalid counter-examples are not shown as examples and construct coverage is sparse, so the requirement is not satisfied. Rated minimum.

### the format specifies its core and how it extends and evolves

- **State:** active
- **Target:** format-spec
- **Rating:** minimum
- **Assessment record:** `assessments/008-format-spec-the-format-specifies-its-core-and-how-it-extends-and-evolves.json`
- **Rationale:** Core and structural extension are now well specified, but 'how it evolves' remains pointer-only and unrecognized-frontmatter handling has a gap, so the conjunction core+extends+evolves is not fully satisfied. This is the format spec's weakest requirement; rated minimum.

### the format specification is well-structured and readable

- **State:** active
- **Target:** format-spec
- **Rating:** target
- **Assessment record:** `assessments/009-format-spec-the-format-specification-is-well-structured-and-readable.json`
- **Rationale:** Dependency-respecting structure, scannable schema blocks, plain prose, and realistic examples satisfy the requirement; a few structural nits (annotation-vs-table, heading hierarchy) keep it from outstanding.

### the README says what QUALITY.md is and who it's for

- **State:** active
- **Target:** readme
- **Rating:** target
- **Assessment record:** `assessments/010-readme-the-readme-says-what-quality-md-is-and-who-it-s-for.json`
- **Rationale:** A first-time reader learns within the opening lines what QUALITY.md is and who it is for, satisfying the requirement; held from outstanding by the audience/example mismatch and an implicit problem statement.

### the README shows the format and its payoff by example

- **State:** active
- **Target:** readme
- **Rating:** minimum
- **Assessment record:** `assessments/011-readme-the-readme-shows-the-format-and-its-payoff-by-example.json`
- **Rationale:** The format-excerpt half is well met, but the requirement explicitly asks for what running qualitymd produces, and no produced output is shown anywhere. The show-don't-tell payoff is missing, holding the requirement at the floor. Rated minimum.

### the README gets a newcomer to a first result quickly

- **State:** active
- **Target:** readme
- **Rating:** minimum
- **Assessment record:** `assessments/012-readme-the-readme-gets-a-newcomer-to-a-first-result-quickly.json`
- **Rationale:** Install commands hold the floor, but there is no short install-then-one-command sequence with representative output that lets a newcomer see a real first result. The requirement's core is unmet despite the underlying CLI flow working well. Rated minimum.

### the README reflects what the CLI and spec actually provide

- **State:** active
- **Target:** readme
- **Rating:** target
- **Assessment record:** `assessments/013-readme-the-readme-reflects-what-the-cli-and-spec-actually-provide.json`
- **Rationale:** Every shown command, flag, and capability matches what the CLI provides today; planned work is marked planned and placeholders provisional, resolving the prior overstatement. Held from outstanding only by the improve-mode omission and an unexhaustive schema cross-check. Rated target.

### the CLI follows its functional specifications

- **State:** active
- **Target:** cli
- **Rating:** minimum
- **Assessment record:** `assessments/014-cli-the-cli-follows-its-functional-specifications.json`
- **Rationale:** Most commands conform strongly (init, spec, lint, build-report gating and byte-stability), but two confirmed MUST-level spec divergences land squarely on the automation contract: a machine-varying absolute path in the show-status payload (determinism/consistency) and a gap kind that contradicts the documented missing-root-analysis contract (consistency). These hold the requirement at the floor rather than target. Rated minimum.

### the CLI follows the project CLI design guide

- **State:** active
- **Target:** cli
- **Rating:** target
- **Assessment record:** `assessments/015-cli-the-cli-follows-the-project-cli-design-guide.json`
- **Rationale:** The CLI strongly follows the design guide on the load-bearing concerns: exit-code categories, conditional styling with no auto-detection, non-interactivity, and a consistent command grammar. Held from outstanding by low-severity deviations (inconsistent '-' support with a miscategorized exit code, and a non-concise root help). Rated target.

## Findings

- `assessments/001-format-spec-the-format-specification-is-complete.json` at `SPECIFICATION.md:137-152`: Model/Target/Factor/Requirement are each shown as annotated YAML skeletons stating requiredness (Required/Recommended/Optional), uniqueness, and conditional cardinality; an implementer could build a parser and an author write a valid file from these.
- `assessments/001-format-spec-the-format-specification-is-complete.json` at `SPECIFICATION.md:112-130`: Principal malformed/omitted handling is stated: a document MUST begin with valid frontmatter containing a conforming Model; null/empty required values MUST be treated as absent; frontmatter that parses as YAML but does not conform is non-conforming.
- `assessments/001-format-spec-the-format-specification-is-complete.json` at `SPECIFICATION.md:283-290`: The Markdown body is fully open (no required section names/order/content); recommended sections are listed as non-normative, so an author can write a valid file with an empty body.
- `assessments/001-format-spec-the-format-specification-is-complete.json` at `SPECIFICATION.md:112-130`: Gap: absent-frontmatter-entirely, multiple frontmatter blocks, and a valid-YAML-but-non-map top node are not explicitly handled.
- `assessments/001-format-spec-the-format-specification-is-complete.json` at `SPECIFICATION.md:137-152`: Gap: scalar field types are abbreviated as <string>/<level-name> without stating allowed character sets, whitespace, or length bounds, leaving some edge cases undefined.
- `assessments/002-format-spec-the-format-specification-admits-a-single-interpretation.json` at `SPECIFICATION.md:26-34`: BCP 14 keywords are formally invoked and scoped to all-caps usage, with an explicit normative/non-normative split, anchoring obligation strength to a recognized convention.
- `assessments/002-format-spec-the-format-specification-admits-a-single-interpretation.json` at `SPECIFICATION.md:345-347,384-385`: Gap: 'insufficient to rate' / 'too little has been assessed' qualifies the conformance-relevant not-assessed outcome (a MUST) with no stated bound, so different evaluators could draw the line differently.
- `assessments/002-format-spec-the-format-specification-admits-a-single-interpretation.json` at `SPECIFICATION.md:160-163`: The two-level minimum is stated twice in adjacent sentences; consistent but a minor near-duplicate normative statement.
- `assessments/002-format-spec-the-format-specification-admits-a-single-interpretation.json` at `SPECIFICATION.md:387-388`: 'most responsible' is a soft quantifier but sits under SHOULD (advisory), so it does not affect conformance and admits a reasonable single reading.
- `assessments/003-format-spec-the-format-specification-separates-rules-from-rationale.json` at `SPECIFICATION.md:33-37`: A global rule/rationale separation is declared: all content normative except passages explicitly marked non-normative; examples and notes are non-normative.
- `assessments/003-format-spec-the-format-specification-separates-rules-from-rationale.json` at `SPECIFICATION.md:257-260,447,474`: Notes and both appendices are explicitly tagged non-normative, keeping rationale visibly separated from rules.
- `assessments/003-format-spec-the-format-specification-separates-rules-from-rationale.json` at `SPECIFICATION.md:252-256`: Gap: a rationale sentence ('Referencing names that entity once instead of copying criteria that would drift') sits unmarked among normative statements, so it is normative by default though it reads as explanation.
- `assessments/003-format-spec-the-format-specification-separates-rules-from-rationale.json` at `SPECIFICATION.md:476-498`: The only complete worked file example lives in non-normative Appendix B, so no normative rule depends solely on an example.
- `assessments/004-format-spec-the-format-specification-defines-its-terms-before-use.json` at `SPECIFICATION.md:64-103`: A dedicated Terminology section front-loads the core terms (Model, Target, Source, Factor, Requirement, Assessment, Finding, Rating Scale/Level/Result) before the normative sections that use them.
- `assessments/004-format-spec-the-format-specification-defines-its-terms-before-use.json` at `SPECIFICATION.md:367-381`: Factor rating, local rating, and aggregate rating are each defined inline immediately before/at their first normative use.
- `assessments/004-format-spec-the-format-specification-defines-its-terms-before-use.json` at `SPECIFICATION.md:99-100`: 'not assessed' is defined at first use as a Rating Result outcome and used consistently thereafter.
- `assessments/004-format-spec-the-format-specification-defines-its-terms-before-use.json` at `SPECIFICATION.md:49-62`: Gap: conformance-class terms (parser/linter/evaluator/report renderer) are used in the Conformance section before Terminology, though each is characterized operationally in place.
- `assessments/005-format-spec-the-format-specification-is-internally-consistent.json` at `SPECIFICATION.md:156-157,183,190-200`: ratingScale ownership is consistent across statements: a Target MUST NOT declare ratingScale and the Target schema omits it.
- `assessments/005-format-spec-the-format-specification-is-internally-consistent.json` at `SPECIFICATION.md:262-269,275-277`: Factor-connection rules agree across their positive and negative statements (every requirement connected to a factor; direct-target requirements MUST declare factors).
- `assessments/005-format-spec-the-format-specification-is-internally-consistent.json` at `SPECIFICATION.md:377-381`: Leaf and grouping roll-up definitions compose without contradiction (leaf aggregate equals local; grouping aggregate derives from children).
- `assessments/005-format-spec-the-format-specification-is-internally-consistent.json` at `SPECIFICATION.md:476-498`: The Appendix B example agrees with the rules it illustrates: two rating levels (>=2 minimum), required level/criterion present, optional title omitted, requirement placement-connected to its factor.
- `assessments/006-format-spec-each-rule-is-observable-or-testable.json` at `SPECIFICATION.md:239-281`: Requirement-shape rules are statically checkable (exactly one assessment as a non-empty scalar; ratings keyed by scale levels; factor references resolve to an in-scope factor).
- `assessments/006-format-spec-each-rule-is-observable-or-testable.json` at `SPECIFICATION.md:160-161`: Gap: 'ordered from best to worst' is not machine-verifiable from the data alone; no comparable value encodes goodness, so it relies on author intent.
- `assessments/006-format-spec-each-rule-is-observable-or-testable.json` at `SPECIFICATION.md:345-347`: Gap: 'insufficient to rate' is a judgment, not a determinate test, so the not-assessed boundary case is not reproducibly decidable across readers.
- `assessments/006-format-spec-each-rule-is-observable-or-testable.json` at `SPECIFICATION.md:355-358`: Gap (intentional): roll-up defines no numeric aggregation formula, only 'preserve relationships and distinctions', so exact roll-up outcomes are not testable for correctness.
- `assessments/007-format-spec-the-format-s-constructs-are-shown-with-valid-and-invalid-examples.json` at `SPECIFICATION.md:445-498`: Valid worked examples exist: Appendix A a sample rating scale and Appendix B a complete realistic file (frontmatter plus body).
- `assessments/007-format-spec-the-format-s-constructs-are-shown-with-valid-and-invalid-examples.json` at `SPECIFICATION.md:249-250,275-277`: Gap: invalid cases are enumerated only in prose (missing/empty/null/list-valued assessment; factors null/[]) and never shown as worked invalid YAML snippets.
- `assessments/007-format-spec-the-format-s-constructs-are-shown-with-valid-and-invalid-examples.json` at `SPECIFICATION.md:476-498`: Gap: many constructs (direct-target requirement requiring explicit factors, sub-factors, multi-target trees, source/glob selectors, ratings criterion overrides) have no worked example, valid or invalid.
- `assessments/008-format-spec-the-format-specifies-its-core-and-how-it-extends-and-evolves.json` at `SPECIFICATION.md:132-154`: The minimal core is specified (ratingScale plus at least one of factors/requirements/targets), though it is assembled across several sections rather than stated as one consolidated minimal-file rule.
- `assessments/008-format-spec-the-format-specifies-its-core-and-how-it-extends-and-evolves.json` at `SPECIFICATION.md:431-443`: An Extensions section covers added frontmatter properties, collision-resistant naming, no change to defined-property meaning, and preservation of unrecognized body content.
- `assessments/008-format-spec-the-format-specifies-its-core-and-how-it-extends-and-evolves.json` at `SPECIFICATION.md:113-117,125-127,437`: Gap: an unrecognized frontmatter key that is neither defined nor a sanctioned extension has no explicit reader rule and leans toward making the file non-conforming.
- `assessments/008-format-spec-the-format-specifies-its-core-and-how-it-extends-and-evolves.json` at `SPECIFICATION.md:3,16-19`: Gap: the forward-evolution/versioning policy is delegated to an external docs/reference/versioning.md; the spec alone does not state how the format versions forward (compatibility rules, file-version field).
- `assessments/009-format-spec-the-format-specification-is-well-structured-and-readable.json` at `SPECIFICATION.md:105-296`: Document ordering is logical and concept-before-detail: Conformance, Terminology, Document Structure, Frontmatter Schema, Body Semantics, Evaluation Semantics, Report Semantics, Extensions, Appendices.
- `assessments/009-format-spec-the-format-specification-is-well-structured-and-readable.json` at `SPECIFICATION.md:445-498`: Two appendices provide minimal and realistic examples, clearly marked non-normative; prose throughout is plain and declarative.
- `assessments/009-format-spec-the-format-specification-is-well-structured-and-readable.json` at `SPECIFICATION.md:137-152`: Gap: requiredness is conveyed via YAML code-comment annotations rather than true field tables with type/required/cardinality/default columns, so allowed-values and defaults are not uniformly tabulated.
- `assessments/009-format-spec-the-format-specification-is-well-structured-and-readable.json` at `SPECIFICATION.md:132,158,185-236`: Gap: heading levels are inconsistent (Model/Target/Factor/Requirement at H4 directly under an H2, skipping an intervening level).
- `assessments/010-readme-the-readme-says-what-quality-md-is-and-who-it-s-for.json` at `README.md:3-4`: The opening sentence states what QUALITY.md is (an agent-friendly file format plus companion skill and CLI) and who it is for (people building coding-agent/AI-assistant projects/harnesses) in the first lines.
- `assessments/010-readme-the-readme-says-what-quality-md-is-and-who-it-s-for.json` at `README.md:35-39`: The 'The Format' section gives a clear jargon-light definition (a structured quality model plus plain-language rationale), reinforcing the what.
- `assessments/010-readme-the-readme-says-what-quality-md-is-and-who-it-s-for.json` at `README.md:4,49`: Gap: the stated audience (coding-agent/AI harnesses) mismatches the worked 'Support Inbox' customer-support example, which can confuse a first-time reader about who the tool is for.
- `assessments/010-readme-the-readme-says-what-quality-md-is-and-who-it-s-for.json` at `README.md:3-4`: Gap: the problem solved is only implicit and somewhat circular ('continuously improving quality'); a reader does not learn what pain it removes.
- `assessments/011-readme-the-readme-shows-the-format-and-its-payoff-by-example.json` at `README.md:47-99`: A realistic, complete QUALITY.md excerpt (frontmatter ratingScale/targets/factors/requirements plus a Markdown body) is shown before the reference schema, matching 'before reference detail'.
- `assessments/011-readme-the-readme-shows-the-format-and-its-payoff-by-example.json` at `README.md:101-102`: Gap: the payoff is described only in prose; the README never shows what running qualitymd against the example produces (no lint output, init receipt, report.md/report.json excerpt, or rating result).
- `assessments/011-readme-the-readme-shows-the-format-and-its-payoff-by-example.json` at `README.md:174-176`: A conceptual assessment -> findings -> rating result diagram conveys the idea but is not a real tool-output example.
- `assessments/011-readme-the-readme-shows-the-format-and-its-payoff-by-example.json` at `grep -nEi 'report\.md|report\.json|\$ qualitymd|Findings:' README.md`: Re-checked: no rendered qualitymd output or report excerpt appears anywhere in README.md (only the conceptual diagram and a build-report bullet).
- `assessments/012-readme-the-readme-gets-a-newcomer-to-a-first-result-quickly.json` at `README.md:6-18`: The Install section gives two short copyable commands (skill via npx, CLI via npm).
- `assessments/012-readme-the-readme-gets-a-newcomer-to-a-first-result-quickly.json` at `README.md:20-31`: Gap: the Usage section routes the newcomer to the /quality skill, not a copyable install-then-one-command CLI sequence that shows a real result.
- `assessments/012-readme-the-readme-gets-a-newcomer-to-a-first-result-quickly.json` at `qualitymd init; qualitymd lint`: Re-checked: the natural CLI first-result path works end-to-end (init creates QUALITY.md and points to lint; lint prints 'QUALITY.md is valid.'), but the README never shows it as a quickstart nor shows representative output. The gap is omission, not inaccuracy.
- `assessments/012-readme-the-readme-gets-a-newcomer-to-a-first-result-quickly.json` at `README.md:194,204-205`: CI exit-code behavior is referenced and accurate (lint non-zero on errors; build-report --fail-at-or-below), but no representative output sample is shown for the newcomer.
- `assessments/013-readme-the-readme-reflects-what-the-cli-and-spec-actually-provide.json` at `README.md:180-188`: The CLI status note is accurate and marked provisional ('an early work in progress'; ships init, lint, spec, and the evaluation surface) and the 'deterministic, never calls a model' framing matches the implemented surface; verified against --help.
- `assessments/013-readme-the-readme-reflects-what-the-cli-and-spec-actually-provide.json` at `README.md:192-205`: Every command and flag in the deterministic-surface list (init, lint, spec, evaluation create-run/add-record/set-planned-coverage/show-status/build-report, --fail-at-or-below) was verified to exist with matching behavior.
- `assessments/013-readme-the-readme-reflects-what-the-cli-and-spec-actually-provide.json` at `README.md:24-31`: The /quality skill usage examples are backed by real skill modes; minor: the improve mode exists in the skill but is not surfaced in the README usage list.
- `assessments/013-readme-the-readme-reflects-what-the-cli-and-spec-actually-provide.json` at `README.md:127-151`: The condensed Model Schema aligns with qualitymd spec; not exhaustively field-diffed against SPECIFICATION.md (limitation).
- `assessments/014-cli-the-cli-follows-its-functional-specifications.json` at `qualitymd evaluation show-status <run> --json`: Re-checked: show-status --json emits an ABSOLUTE filesystem path in the 'path' field while create-run --json emits a repository-relative path. This is a machine-varying payload value, violating specs/cli/evaluation-create-run.md:33 (path MUST be repository-relative) and specs/cli.md:99-100 (deterministic output with no machine-varying values). add-record and set-planned-coverage receipts share the absolute-path behavior.
- `assessments/014-cli-the-cli-follows-its-functional-specifications.json` at `qualitymd evaluation show-status <empty-run> --json`: Re-checked: an empty run reports gap kind 'missing-analysis', but specs/cli/evaluation-show-status.md:29-32 mandates a 'missing-root-analysis' gap when no analysis record has an empty targetPath. A caller branching on the documented gap kind misses the empty-run case.
- `assessments/014-cli-the-cli-follows-its-functional-specifications.json` at `qualitymd evaluation build-report <run> --fail-at-or-below`: The gate matches specs/cli/evaluation-build-report.md:96-99 exactly: gating at the root rating exits 1 (equal-or-worse), gating at a better level exits 0, an off-scale level exits 2; report.md/report.json are still written before gating.
- `assessments/014-cli-the-cli-follows-its-functional-specifications.json` at `qualitymd evaluation build-report <run>`: build-report is deterministic and idempotent: report.md is byte-identical on rebuild and contains zero ANSI/terminal escape sequences; the human 'Wrote ...' line is stderr-only while stdout stays empty without --json.
- `assessments/014-cli-the-cli-follows-its-functional-specifications.json` at `internal/cli/init.go:44-65; internal/cli/spec.go:15-27`: init and spec conform strongly to their specs: init output-target/-/--force/overwrite-protection rules hold and the scaffold lints clean; spec emits the bundled specification verbatim to stdout and refuses --json/extra args with usage exit 2.
- `assessments/014-cli-the-cli-follows-its-functional-specifications.json` at `qualitymd init --json (existing file)`: Gap: on an overwrite refusal init --json emits the spec'd JSON error object on stderr but ALSO a duplicate styled ERROR block, so a JSON consumer reading stderr sees two representations of the same failure.
- `assessments/015-cli-the-cli-follows-the-project-cli-design-guide.json` at `internal/cli/style.go:33-42`: Color/glyphs are applied only when the writer is a real TTY and NO_COLOR is unset; piped/redirected output is the plain canonical form byte-for-byte, with no format auto-detection (docs/guides/cli-design.md:118-131).
- `assessments/015-cli-the-cli-follows-the-project-cli-design-guide.json` at `qualitymd frobnicate; qualitymd lint --bogus; qualitymd lint <missing>`: Exit-code categories match the binding table (docs/guides/cli-design.md:191-196): success 0, lint findings 1, unknown command/flag and usage validation 2, I/O and unmet preconditions 70.
- `assessments/015-cli-the-cli-follows-the-project-cli-design-guide.json` at `internal/cli/evaluation.go:16-27,71-97`: Consistent 'noun verb' grammar (evaluation create-run / add-record <kind> <run> / show-status <run>), discoverable record subcommands, uniform --file/--json flags, both -h and --help, and no catch-all subcommand (docs/guides/cli-design.md:164-176).
- `assessments/015-cli-the-cli-follows-the-project-cli-design-guide.json` at `internal/cli/evaluation.go:198-210`: Record-writing commands never block on a prompt: they accept --file <path>, --file -, or piped stdin and fail with a usage error when input is absent on a TTY; strict JSON decoding rejects CLI-owned fields (docs/guides/cli-design.md:154-160,202-208).
- `assessments/015-cli-the-cli-follows-the-project-cli-design-guide.json` at `qualitymd lint - ; qualitymd init -`: Gap: init supports '-' for stdin/stdout but lint does not (docs/guides/cli-design.md:157-158), an inconsistency; and the lint '-' refusal maps to exit 70 rather than the usage category 2.
- `assessments/015-cli-the-cli-follows-the-project-cli-design-guide.json` at `qualitymd (bare, no args)`: Gap: a bare invocation prints the full Fang help rather than the guide's concise, example-led default for the nothing-to-do case (docs/guides/cli-design.md:85-86).

## Advice

- [001-show-what-running-qualitymd-produces-not-just-the-model](recommendations/001-show-what-running-qualitymd-produces-not-just-the-model.md) [active] - The requirement 'the README shows the format and its payoff by example' reaches at least target: the README displays representative qualitymd output produced from the shown example. Re-evaluate in a new numbered run.
- [002-add-a-copyable-cli-quickstart-with-a-representative-first-result](recommendations/002-add-a-copyable-cli-quickstart-with-a-representative-first-result.md) [active] - The requirement 'the README gets a newcomer to a first result quickly' reaches at least target: a copyable install-then-one-command sequence with representative output appears early in the README. Re-evaluate in a new numbered run.
- [003-make-evaluation-command-payload-paths-repository-relative](recommendations/003-make-evaluation-command-payload-paths-repository-relative.md) [active] - show-status --json (and the other evaluation receipts) emit a repository-relative 'path' identical in form to create-run, with no machine-varying value in any --json payload. Re-evaluate in a new numbered run.
- [004-reconcile-the-empty-run-gap-kind-with-the-missing-root-analysis-contract](recommendations/004-reconcile-the-empty-run-gap-kind-with-the-missing-root-analysis-contract.md) [active] - An empty run's show-status gap kind matches the spec exactly (either the CLI emits 'missing-root-analysis' or the spec documents 'missing-analysis' for the zero case). Re-evaluate in a new numbered run.
- [005-show-invalid-counter-examples-and-cover-more-constructs-in-the-spec](recommendations/005-show-invalid-counter-examples-and-cover-more-constructs-in-the-spec.md) [active] - The requirement 'the format's constructs are shown with valid and invalid examples' reaches at least target: worked invalid counter-examples accompany the principal constructs. Re-evaluate in a new numbered run.
- [006-state-the-format-s-forward-evolution-and-unrecognized-key-handling-in-the-spec](recommendations/006-state-the-format-s-forward-evolution-and-unrecognized-key-handling-in-the-spec.md) [active] - The requirement 'the format specifies its core and how it extends and evolves' reaches at least target: the spec states how the format versions forward and how readers treat unrecognized frontmatter content. Re-evaluate in a new numbered run.
