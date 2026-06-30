# Evaluation Report

## Verdict

- **Subject:** QUALITY.md
- **Altitude:** subject
- **Rigor:** deep
- **Narrowing:** Full evaluation
- **Evaluation verdict:** 🔵 Target
- **Rationale:** The root is a grouping target. All three child deliverables now aggregate to target after the spec, README, CLI, and evaluation-history fixes, so the whole project meets the release quality bar.

## Scope

Full evaluation

- **Narrowing:** Full evaluation
- **In scope:** QUALITY.md; qualitymd CLI; Format specification; README
- **Out of scope:** Hosted CI and publication checks are deferred to the release runbook.; Exhaustive CLI matrix testing is deferred; focused behavior checks and package tests cover the release-relevant regressions found in this run.

## Selected findings and limitations

- `assessments/001-format-spec-the-format-specification-is-complete.json` at `SPECIFICATION.md:137` [low]: Scalar placeholders such as <string> and <level-name> intentionally do not define detailed character-set or length bounds, leaving edge cases to conforming tools.
- `assessments/002-format-spec-the-format-specification-admits-a-single-interpretation.json` at `SPECIFICATION.md:344` [low]: The not-assessed boundary still depends on evaluator judgment about evidence sufficiency, but the report distinction is explicit.
- `assessments/006-format-spec-each-rule-is-observable-or-testable.json` at `SPECIFICATION.md:356` [low]: Roll-up intentionally has no numeric aggregation formula, so exact rating inference remains evaluator judgment constrained by required distinctions.
- Limitation: The spec now includes valid minimal/suggested examples and explicit invalid counter-examples for missing required title, missing direct-requirement factors, and list-valued assessment

## Evidence basis

No command or source evidence basis was recorded in findings.

## Next action

No recommendation records exist for this run.

## Target summary

| Target | Local rating | Aggregate rating | Covered requirements | Note |
| --- | --- | --- | --- | --- |
| QUALITY.md | n/a (structural) | 🔵 Target | 0 | structural grouping target |
| qualitymd CLI | 🔵 Target | 🔵 Target | 2 |   |
| Format specification | 🔵 Target | 🔵 Target | 9 |   |
| README | 🔵 Target | 🔵 Target | 4 |   |

## Target details

### QUALITY.md

- **Path:** (root)
- **Local rating:** n/a (structural)
- **Aggregate rating:** 🔵 Target
  - The root is a grouping target. All three child deliverables now aggregate to target after the spec, README, CLI, and evaluation-history fixes, so the whole project meets the release quality bar.
- **Analysis record:** `analysis/root.json`

### qualitymd CLI

- **Path:** cli
- **Local rating:** 🔵 Target
  - Both CLI requirements meet target after fixing usage-error mapping and repo-relative evaluation status paths/gaps, with focused tests and behavior checks passing.
- **Aggregate rating:** 🔵 Target
  - Leaf target: aggregate equals local rating. The release-relevant CLI conformance gaps found in this run are fixed and verified.
- **Factor Usability:** 🔵 Target
  - Root welcome, help, errors, and next actions are concise and recoverable.
- **Factor Automation compatibility:** 🔵 Target
  - JSON surfaces are explicit, non-interactive input handling is deterministic, and stdout/stderr separation is preserved.
- **Factor Consistency:** 🔵 Target
  - Evaluation commands share the noun/verb grammar and status/write receipts now use consistent repo-relative paths.
- **Factor Determinism:** 🔵 Target
  - The fixed evaluation status output no longer includes machine-varying absolute paths, and report/status behavior is covered by focused tests.
- **Analysis record:** `analysis/cli.json`

### Format specification

- **Path:** format-spec
- **Local rating:** 🔵 Target
  - Nine requirements are covered. All meet target or better after the Appendix B and invalid-counter-example fixes; the remaining gaps are bounded edge precision and evaluator-judgment limits, not release blockers.
- **Aggregate rating:** 🔵 Target
  - Leaf target: aggregate equals local rating. The format spec satisfies every in-scope requirement at target or better.
- **Factor Clarity:** 🔵 Target
  - Single interpretation, rules/rationale separation, and terms-before-use all satisfy the model, with terms-before-use outstanding.
- **Factor Consistency:** 🔵 Target
  - The prior minimal-example contradiction is fixed; examples and normative rules now agree.
- **Factor Verifiability:** 🔵 Target
  - Structural rules are observable and the spec now includes invalid counter-examples; some semantic roll-up judgment remains intentionally non-formulaic.
- **Factor Extensibility:** 🔵 Target
  - The minimal core, extension boundaries, unrecognized-content preservation, and versioning reference provide a compatible evolution path.
- **Factor Usability:** 🔵 Target
  - The document is ordered concept-first, with scannable schemas and valid/invalid examples.
- **Analysis record:** `analysis/format-spec.json`

### README

- **Path:** readme
- **Local rating:** 🔵 Target
  - All four approachability requirements are covered and meet target after adding first-result output and report-output payoff examples.
- **Aggregate rating:** 🔵 Target
  - Leaf target: aggregate equals local rating. The front door now explains and demonstrates the tool adequately for release.
- **Factor Approachability:** 🔵 Target
  - The README says what QUALITY.md is, shows a realistic format example, shows first-result CLI output, and reflects the current command surface.
- **Analysis record:** `analysis/readme.json`

## Requirements

### the format specification is complete

- **State:** active
- **Target:** Format specification
- **Rating:** 🔵 Target
- **Assessment result record:** `assessments/001-format-spec-the-format-specification-is-complete.json`
- **Rationale:** The specification is complete enough for parser, linter, author, evaluator, and report-renderer implementation. Scalar bounds remain intentionally broad, so the result satisfies the requirement without reaching outstanding.

### the format specification admits a single interpretation

- **State:** active
- **Target:** Format specification
- **Rating:** 🔵 Target
- **Assessment result record:** `assessments/002-format-spec-the-format-specification-admits-a-single-interpretation.json`
- **Rationale:** Conformance classes, BCP 14 scope, normative vs non-normative distinctions, and schema wording give one settled reading for release-relevant rules.

### the format specification separates rules from rationale

- **State:** active
- **Target:** Format specification
- **Rating:** 🔵 Target
- **Assessment result record:** `assessments/003-format-spec-the-format-specification-separates-rules-from-rationale.json`
- **Rationale:** Binding rules appear in normative sections, while examples, notes, and appendices are marked non-normative.

### the format specification defines its terms before use

- **State:** active
- **Target:** Format specification
- **Rating:** 🟢 Outstanding
- **Assessment result record:** `assessments/004-format-spec-the-format-specification-defines-its-terms-before-use.json`
- **Rationale:** A dedicated terminology section defines the core vocabulary before the schema and evaluation rules use it.

### the format specification is internally consistent

- **State:** active
- **Target:** Format specification
- **Rating:** 🔵 Target
- **Assessment result record:** `assessments/005-format-spec-the-format-specification-is-internally-consistent.json`
- **Rationale:** The earlier Appendix B contradiction was fixed; current examples and rules use the same required title fields and terminology.

### each rule is observable or testable

- **State:** active
- **Target:** Format specification
- **Rating:** 🔵 Target
- **Assessment result record:** `assessments/006-format-spec-each-rule-is-observable-or-testable.json`
- **Rationale:** Structural conformance requirements map to observable document or record properties; judgment-bearing evaluation semantics are still explicitly semantic rather than mechanically formulaic.

### the format's constructs are shown with valid and invalid examples

- **State:** active
- **Target:** Format specification
- **Rating:** 🔵 Target
- **Assessment result record:** `assessments/007-format-spec-the-format-s-constructs-are-shown-with-valid-and-invalid-examples.json`
- **Rationale:** The spec now includes valid minimal/suggested examples and explicit invalid counter-examples for missing required title, missing direct-requirement factors, and list-valued assessment.

### the format specifies its core and how it extends and evolves

- **State:** active
- **Target:** Format specification
- **Rating:** 🔵 Target
- **Assessment result record:** `assessments/008-format-spec-the-format-specifies-its-core-and-how-it-extends-and-evolves.json`
- **Rationale:** The minimal core, extension boundaries, unrecognized-content preservation, and external versioning policy are defined sufficiently for compatible evolution.

### the format specification is well-structured and readable

- **State:** active
- **Target:** Format specification
- **Rating:** 🔵 Target
- **Assessment result record:** `assessments/009-format-spec-the-format-specification-is-well-structured-and-readable.json`
- **Rationale:** The document introduces concepts before dependent schema and semantics, uses scannable examples, and keeps prose direct.

### the README says what QUALITY.md is and who it's for

- **State:** active
- **Target:** README
- **Rating:** 🔵 Target
- **Assessment result record:** `assessments/010-readme-the-readme-says-what-quality-md-is-and-who-it-s-for.json`
- **Rationale:** The opening and Why section identify the format, companion skill/CLI, intended agent-project audience, and the quality-debt problem it addresses.

### the README shows the format and its payoff by example

- **State:** active
- **Target:** README
- **Rating:** 🔵 Target
- **Assessment result record:** `assessments/011-readme-the-readme-shows-the-format-and-its-payoff-by-example.json`
- **Rationale:** The README now pairs a realistic QUALITY.md excerpt with concrete CLI/report output before the reference schema.

### the README gets a newcomer to a first result quickly

- **State:** active
- **Target:** README
- **Rating:** 🔵 Target
- **Assessment result record:** `assessments/012-readme-the-readme-gets-a-newcomer-to-a-first-result-quickly.json`
- **Rationale:** The README now gives a short install/use path and representative first-result output for init and lint.

### the README reflects what the CLI and spec actually provide

- **State:** active
- **Target:** README
- **Rating:** 🔵 Target
- **Assessment result record:** `assessments/013-readme-the-readme-reflects-what-the-cli-and-spec-actually-provide.json`
- **Rationale:** The README command list and format summary align with the current v0.4 command surface and specification.

### the CLI follows its functional specifications

- **State:** active
- **Target:** qualitymd CLI
- **Rating:** 🔵 Target
- **Assessment result record:** `assessments/014-cli-the-cli-follows-its-functional-specifications.json`
- **Rationale:** Focused checks and tests show the current CLI satisfies the release-relevant functional contracts, including the issues fixed during this evaluation.

### the CLI follows the project CLI design guide

- **State:** active
- **Target:** qualitymd CLI
- **Rating:** 🔵 Target
- **Assessment result record:** `assessments/015-cli-the-cli-follows-the-project-cli-design-guide.json`
- **Rationale:** The command surface is examples-led, non-interactive, stdout/stderr-safe, deterministic under JSON/plain output, and uses consistent noun/verb evaluation commands.

## Findings

- `assessments/001-format-spec-the-format-specification-is-complete.json` at `SPECIFICATION.md:105`: The document structure, frontmatter schema, body semantics, evaluation semantics, report semantics, and extensions are all defined in one specification.
- `assessments/001-format-spec-the-format-specification-is-complete.json` at `SPECIFICATION.md:123`: Model, Rating Scale, Target, Factor, and Requirement shapes define requiredness, cardinality, inheritance/default source behavior, and malformed-value handling.
- `assessments/001-format-spec-the-format-specification-is-complete.json` at `SPECIFICATION.md:137`: Scalar placeholders such as <string> and <level-name> intentionally do not define detailed character-set or length bounds, leaving edge cases to conforming tools.
- `assessments/002-format-spec-the-format-specification-admits-a-single-interpretation.json` at `SPECIFICATION.md:21`: BCP 14 keywords are scoped to all-capital uses, reducing ambiguity around obligation strength.
- `assessments/002-format-spec-the-format-specification-admits-a-single-interpretation.json` at `SPECIFICATION.md:29`: Examples and notes are explicitly non-normative and do not add conformance requirements.
- `assessments/002-format-spec-the-format-specification-admits-a-single-interpretation.json` at `SPECIFICATION.md:344`: The not-assessed boundary still depends on evaluator judgment about evidence sufficiency, but the report distinction is explicit.
- `assessments/003-format-spec-the-format-specification-separates-rules-from-rationale.json` at `SPECIFICATION.md:27`: The specification states that all content is normative except passages explicitly marked otherwise.
- `assessments/003-format-spec-the-format-specification-separates-rules-from-rationale.json` at `SPECIFICATION.md:445`: Appendix A, Appendix B, and Appendix C are each explicitly non-normative before their examples.
- `assessments/004-format-spec-the-format-specification-defines-its-terms-before-use.json` at `SPECIFICATION.md:56`: Quality Model, Entity, Model, Target, Source, Factor, Requirement, Assessment, Finding, Rating Scale, Rating Level, Rating Result, and Evaluation Report are defined before schema rules.
- `assessments/004-format-spec-the-format-specification-defines-its-terms-before-use.json` at `SPECIFICATION.md:123`: Later schema and evaluation sections reuse those defined terms consistently.
- `assessments/005-format-spec-the-format-specification-is-internally-consistent.json` at `SPECIFICATION.md:158`: Rating Level rules require level, title, and criterion; the suggested scale and minimal example now include those fields.
- `assessments/005-format-spec-the-format-specification-is-internally-consistent.json` at `SPECIFICATION.md:485`: Appendix B now includes title fields for the target and unacceptable rating levels and for the reliability factor.
- `assessments/006-format-spec-each-rule-is-observable-or-testable.json` at `SPECIFICATION.md:112`: Frontmatter validity, required properties, list cardinality, factor references, and assessment scalar shape are all observable from the document.
- `assessments/006-format-spec-each-rule-is-observable-or-testable.json` at `SPECIFICATION.md:356`: Roll-up intentionally has no numeric aggregation formula, so exact rating inference remains evaluator judgment constrained by required distinctions.
- `assessments/007-format-spec-the-format-s-constructs-are-shown-with-valid-and-invalid-examples.json` at `SPECIFICATION.md:445`: Appendix A and Appendix B provide copyable valid examples for the rating scale and a minimal file.
- `assessments/007-format-spec-the-format-s-constructs-are-shown-with-valid-and-invalid-examples.json` at `SPECIFICATION.md:509`: Appendix C provides worked invalid counter-examples for a missing rating-level title, direct target requirement without factors, and list-valued assessment.
- `assessments/008-format-spec-the-format-specifies-its-core-and-how-it-extends-and-evolves.json` at `SPECIFICATION.md:132`: The root model requires ratingScale plus at least one of factors, requirements, or targets.
- `assessments/008-format-spec-the-format-specifies-its-core-and-how-it-extends-and-evolves.json` at `SPECIFICATION.md:431`: Extensions cannot change conforming-property meaning, and tools should preserve unrecognized extension fields and body content.
- `assessments/008-format-spec-the-format-specifies-its-core-and-how-it-extends-and-evolves.json` at `SPECIFICATION.md:16`: The specification identifies its version and points to the project versioning policy for forward-versioning behavior.
- `assessments/009-format-spec-the-format-specification-is-well-structured-and-readable.json` at `SPECIFICATION.md:56`: Terminology precedes document structure, schema, and evaluation semantics.
- `assessments/009-format-spec-the-format-specification-is-well-structured-and-readable.json` at `SPECIFICATION.md:137`: Annotated YAML blocks make required fields and nesting shapes quick to scan.
- `assessments/009-format-spec-the-format-specification-is-well-structured-and-readable.json` at `SPECIFICATION.md:445`: Valid and invalid appendices give copy-and-adapt examples after the normative sections.
- `assessments/010-readme-the-readme-says-what-quality-md-is-and-who-it-s-for.json` at `README.md:3`: The first paragraph says QUALITY.md is an agent-friendly file format with companion skill and CLI for coding-agent and AI assistant projects/harnesses.
- `assessments/010-readme-the-readme-says-what-quality-md-is-and-who-it-s-for.json` at `README.md:7`: The Why section frames the problem as technical, cognitive, and intent debt becoming explicit and checkable.
- `assessments/011-readme-the-readme-shows-the-format-and-its-payoff-by-example.json` at `README.md:66`: A realistic Support Inbox QUALITY.md excerpt shows ratingScale, target source, factors, requirements, and body context.
- `assessments/011-readme-the-readme-shows-the-format-and-its-payoff-by-example.json` at `README.md:129`: The README shows the evaluation report build output and a report-summary excerpt, making the CLI payoff visible rather than only described.
- `assessments/012-readme-the-readme-gets-a-newcomer-to-a-first-result-quickly.json` at `README.md:25`: Install remains a two-step skill and CLI install sequence.
- `assessments/012-readme-the-readme-gets-a-newcomer-to-a-first-result-quickly.json` at `README.md:37`: The Usage section gives `qualitymd init` then `qualitymd lint`, followed by expected output showing `Created QUALITY.md` and `QUALITY.md is valid.`
- `assessments/013-readme-the-readme-reflects-what-the-cli-and-spec-actually-provide.json` at `README.md:196`: The CLI section lists current commands: init, lint, spec, status, version, upgrade, and the evaluation run-record surface.
- `assessments/013-readme-the-readme-reflects-what-the-cli-and-spec-actually-provide.json` at `README.md:216`: The typical local loop uses existing commands: spec, init, lint, and status --json.
- `assessments/013-readme-the-readme-reflects-what-the-cli-and-spec-actually-provide.json` at `go run ./cmd/qualitymd --help`: The current command tree exposes the commands named in README.md.
- `assessments/014-cli-the-cli-follows-its-functional-specifications.json` at `go run ./cmd/qualitymd evaluation status quality/evaluations/0005-subject-quality-eval --json`: evaluation status now emits a repository-relative `path`, repository-relative next action, and `missing-root-analysis` for an empty run.
- `assessments/014-cli-the-cli-follows-its-functional-specifications.json` at `go run ./cmd/qualitymd lint -`: lint '-' now maps to the usage category; go run reports `exit status 2`, matching the CLI exit-code contract.
- `assessments/014-cli-the-cli-follows-its-functional-specifications.json` at `go test ./internal/cli ./internal/evaluation ./internal/status`: Focused regression tests pass for CLI behavior, evaluation record/report behavior, and status snapshots.
- `assessments/015-cli-the-cli-follows-the-project-cli-design-guide.json` at `internal/cli/style.go:68`: The no-argument root command renders concise examples-first welcome output rather than full help.
- `assessments/015-cli-the-cli-follows-the-project-cli-design-guide.json` at `internal/cli/evaluation.go:16`: The evaluation surface follows a consistent noun/verb tree: create, list, status, assessment add/list, analysis set/list, recommendation add/list, report build/gate.
- `assessments/015-cli-the-cli-follows-the-project-cli-design-guide.json` at `internal/cli/style.go:44`: Human styling is conditional on terminal output and NO_COLOR, keeping redirected output plain.

## Advice

No recommendation records exist for this run.
