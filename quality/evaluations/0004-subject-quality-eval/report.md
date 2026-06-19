# Evaluation Report

## Summary

- **Subject:** QUALITY.md
- **Altitude:** subject
- **Effort:** standard
- **Rating:** 🟡 Minimum
- **Rationale:** Structural root with no local requirements (source missing by design). Aggregate over three child targets: format-spec at target (strongest), readme at minimum, cli at minimum. Two of three deliverables sit at the acceptable floor, so the whole-project aggregate is minimum: the format specification is solid, while the README does not demonstrate the tool working and the CLI has a confirmed determinism MUST-violation in its automation surface.

## Scope

standard` — every in-scope Requirement assessed against targeted evidence, with the rating-binding findings re-checked before reporting.

- **Narrowing:** whole recorded run
- **In scope:** QUALITY.md; qualitymd CLI; Format specification; README
- **Out of scope:** Dependencies the project does not own (Go toolchain, Cobra/Fang/Lip Gloss, release tooling) — excluded by the model's Scope section.; The `/quality` skill prompt and CLI internals beyond `internal/cli` behavior.; Deep adversarial fan-out (this is a `standard`-effort run, not `deep`).

## Top Risks and Limitations

- `assessments/001-format-spec-the-format-specification-is-complete.json` at `SPECIFICATION.md:137-152` [low]: Held below outstanding: scalar field types are abbreviated as <string>/<level-name> without character-set, whitespace, or length bounds, leaving some edge cases undefined.
- `assessments/005-format-spec-the-format-specification-is-internally-consistent.json` at `SPECIFICATION.md:481-493` [medium]: Appendix B's minimal example omits the required 'title' on both rating levels (target/unacceptable carry only level/description/criterion) and on the 'reliability' factor, contradicting the rule the example illustrates.
- `assessments/007-format-spec-the-format-s-constructs-are-shown-with-valid-and-invalid-examples.json` at `SPECIFICATION.md:250-278` [medium]: Invalidity is described only in prose ('A missing, empty, null, or list-valued assessment is invalid'; 'Missing factors, factors: null, factors: [] ... do not satisfy'); there are no worked invalid counter-example blocks showing what a malformed construct looks like.
- `assessments/008-format-spec-the-format-specifies-its-core-and-how-it-extends-and-evolves.json` at `SPECIFICATION.md:16-19` [low]: Forward versioning is addressed via the specification-version field referencing the versioning policy; held from outstanding because the versioning mechanics live in an external doc rather than the spec body.
- `assessments/011-readme-the-readme-shows-the-format-and-its-payoff-by-example.json` at `README.md:120-121` [medium]: The payoff is described only in prose ('An agent that reads this file can evaluate support work ... produce findings, and rate the results'); the README never shows what running qualitymd against the example produces (no lint output, init receipt, or report excerpt).
- `assessments/012-readme-the-readme-gets-a-newcomer-to-a-first-result-quickly.json` at `README.md:207-230` [medium]: No representative output is shown for any command, and CI exit-code behavior is mentioned only in passing ('can fail CI when ratings fall below your chosen bar') without an example, so a newcomer runs commands blind to what success looks like.
- `assessments/013-readme-the-readme-reflects-what-the-cli-and-spec-actually-provide.json` at `README.md:200-202` [low]: The CLI is clearly marked 'an early work in progress', but its 'Today the binary ships' sentence lists only init/lint/spec/status and the evaluation surface, omitting version and upgrade, which do ship and appear in the table below it (internal understatement, not an overstatement).
- `assessments/014-cli-the-cli-follows-its-functional-specifications.json` at `internal/cli/evaluation.go:129-149` [high]: 'qualitymd evaluation show-status <run> --json' emits a machine-varying absolute path in the 'path' field and inside nextActions[].command even when given a repo-relative run path (verified: a relative input was echoed as /Users/craig/Code/qualitymd/quality.md/quality/evaluations/...). This violates payload determinism for a machine-consumed field.
- Limitation: Structural root with no local requirements (source missing by design)
- Limitation: Normative statements and terminology are mutually consistent, but the requirement explicitly tests that every example agrees with the rule it illustrates, and the Appendix B minimal example violates required-field rules (missing title on rating levels and factor)
- Limitation: A newcomer gets a copyable install-then-command sequence, but the requirement also asks for representative output and CI exit-code behavior, which are absent
- Limitation: Satisfies the requirement; held below outstanding by a missing --no-color flag and the determinism deviation carried from the functional-spec finding
- Additional risks or limitations are available in `report.json`.

## Evidence Basis

No command or source evidence basis was recorded in findings.

## Next Action

- [001-make-show-status-json-emit-repo-relative-paths](recommendations/001-make-show-status-json-emit-repo-relative-paths.md) - show-status --json emits a repo-relative path for both 'path' and nextActions[].command given a relative run argument, and the requirement 'the CLI follows its functional specifications' reaches at least target. Re-evaluate in a new numbered run.

## Target Summary

| Target               | Local            | Aggregate  | Covered Requirements | Note                       |
| -------------------- | ---------------- | ---------- | -------------------- | -------------------------- |
| QUALITY.md           | n/a (structural) | 🟡 Minimum | 0                    | structural grouping target |
| qualitymd CLI        | 🟡 Minimum       | 🟡 Minimum | 2                    |                            |
| Format specification | 🔵 Target        | 🔵 Target  | 9                    |                            |
| README               | 🟡 Minimum       | 🟡 Minimum | 4                    |                            |

## Target Details

### QUALITY.md

- **Path:** (root)
- **Local rating:** n/a (structural)
  - Structural grouping target; local rating does not apply.
- **Aggregate rating:** 🟡 Minimum
  - Structural root with no local requirements (source missing by design). Aggregate over three child targets: format-spec at target (strongest), readme at minimum, cli at minimum. Two of three deliverables sit at the acceptable floor, so the whole-project aggregate is minimum: the format specification is solid, while the README does not demonstrate the tool working and the CLI has a confirmed determinism MUST-violation in its automation surface.
- **Analysis record:** `analysis/quality-md.json`

### qualitymd CLI

- **Path:** cli
- **Local rating:** 🟡 Minimum
  - Two requirements: 'follows its functional specifications' at minimum and 'follows the project CLI design guide' at target. Design-guide adherence (exit codes, stream separation, conditional styling, non-interactivity, consistent grammar, examples-first help) is strong, but a confirmed determinism MUST-violation (show-status --json absolute path) plus a 'lint -' exit-code inconsistency in the automation surface bind the local rating to the floor.
- **Aggregate rating:** 🟡 Minimum
  - Leaf target: aggregate equals local rating. Strong design-guide adherence and automation ergonomics are pulled to the floor by one confirmed determinism MUST-violation and an exit-code consistency gap.
- **Factor Consistency:** 🟡 Minimum
  - Command grammar, flags, and help are uniform, but sibling commands disagree on path form (show-status absolute vs create-run relative) and 'lint -' exits 70 where 'status -' exits 2.
- **Factor Determinism:** 🟡 Minimum
  - build-report is byte-stable and ANSI-free, but show-status --json emits a machine-varying absolute path, a determinism MUST violation in a machine-consumed payload field (verified).
- **Factor Automation compatibility:** 🔵 Target
  - Strict JSON decoding, --json everywhere, non-interactive stdin handling, stdout/stderr separation, and exit-code categories all hold.
- **Factor Usability:** 🔵 Target
  - Help leads with examples, next-action footers and error categories are clear and recoverable.
- **Analysis record:** `analysis/cli.json`

### Format specification

- **Path:** format-spec
- **Local rating:** 🔵 Target
  - Nine requirements: the target-level completeness requirement and Clarity/Extensibility/Usability all at target, define-terms-before-use at outstanding, with Consistency and Verifiability at minimum. The spec satisfies its core job of letting an implementer build a parser and an author write a valid file; the two minimum factors are bounded documentation gaps (one flawed non-normative example; invalid cases described in prose rather than shown).
- **Aggregate rating:** 🔵 Target
  - Leaf target: aggregate equals local rating. The strongest target in the model; complete, clear, extensible, and usable, with two bounded documentation gaps in consistency and verifiability.
- **Factor Clarity:** 🔵 Target
  - Single interpretation and rules-vs-rationale at target; define-terms-before-use at outstanding via a dedicated glossary preceding all normative use.
- **Factor Consistency:** 🟡 Minimum
  - Normative statements and terminology are mutually consistent, but the Appendix B minimal example omits required title fields, contradicting the rule it illustrates.
- **Factor Verifiability:** 🟡 Minimum
  - Rules are observable/testable (target), but constructs are not shown with worked invalid counter-examples; invalidity is only described in prose.
- **Factor Extensibility:** 🔵 Target
  - Core, extension rules, unrecognized-content handling, and forward versioning are all specified; versioning detail is delegated to an external policy doc.
- **Factor Usability:** 🔵 Target
  - Logical ordering, scannable annotated schema blocks, plain prose, and copy-ready examples make the spec navigable.
- **Analysis record:** `analysis/format-spec.json`

### README

- **Path:** readme
- **Local rating:** 🟡 Minimum
  - Four Approachability requirements: what-it-is/who-for and reflects-CLI/spec at target, but shows-payoff-by-example and gets-to-first-result both at minimum. The README is a clear, accurate front door, but does not show what running qualitymd produces or representative output, so a newcomer cannot see the show-don't-tell payoff or what success looks like.
- **Aggregate rating:** 🟡 Minimum
  - Leaf target: aggregate equals local rating. The model's weakest target; the front door explains itself well but never demonstrates the tool working.
- **Factor Approachability:** 🟡 Minimum
  - Strong on what/why and accuracy, but two of four requirements fall to the floor because the README describes the payoff and first result in prose without showing any rendered command output.
- **Analysis record:** `analysis/readme.json`

## Requirements

### the format specification is complete

- **State:** active
- **Target:** Format specification
- **Rating:** 🔵 Target
- **Assessment record:** `assessments/001-format-spec-the-format-specification-is-complete.json`
- **Rationale:** Frontmatter and body are specified completely enough to parse and author from, including principal malformed/omitted handling, satisfying the requirement. Held below outstanding only by untyped scalar bounds.

### the format specification admits a single interpretation

- **State:** active
- **Target:** Format specification
- **Rating:** 🔵 Target
- **Assessment record:** `assessments/002-format-spec-the-format-specification-admits-a-single-interpretation.json`
- **Rationale:** Each rule admits a single reading: obligation strength is explicit, BCP 14 keywords are used consistently, and quantifiers are bounded. Satisfies the requirement.

### the format specification separates rules from rationale

- **State:** active
- **Target:** Format specification
- **Rating:** 🔵 Target
- **Assessment record:** `assessments/003-format-spec-the-format-specification-separates-rules-from-rationale.json`
- **Rationale:** A reader can always tell a binding rule from explanation via the explicit normativity convention and marked notes/examples; no rule appears only in an example. Satisfies the requirement.

### the format specification defines its terms before use

- **State:** active
- **Target:** Format specification
- **Rating:** 🟢 Outstanding
- **Assessment record:** `assessments/004-format-spec-the-format-specification-defines-its-terms-before-use.json`
- **Rationale:** A complete glossary precedes all normative use, exceeding a baseline of define-before-use by centralizing definitions ahead of every dependent section. Exceeds the requirement with margin.

### the format specification is internally consistent

- **State:** active
- **Target:** Format specification
- **Rating:** 🟡 Minimum
- **Assessment record:** `assessments/005-format-spec-the-format-specification-is-internally-consistent.json`
- **Rationale:** Normative statements and terminology are mutually consistent, but the requirement explicitly tests that every example agrees with the rule it illustrates, and the Appendix B minimal example violates required-field rules (missing title on rating levels and factor). Falls short of full satisfaction; holds the acceptable floor because the defect is confined to one non-normative example.

### each rule is observable or testable

- **State:** active
- **Target:** Format specification
- **Rating:** 🔵 Target
- **Assessment record:** `assessments/006-format-spec-each-rule-is-observable-or-testable.json`
- **Rationale:** Each normative rule maps to something a reader or tool can observe or test about a file or implementation, so independent readers decide conformance the same way. Satisfies the requirement.

### the format's constructs are shown with valid and invalid examples

- **State:** active
- **Target:** Format specification
- **Rating:** 🟡 Minimum
- **Assessment record:** `assessments/007-format-spec-the-format-s-constructs-are-shown-with-valid-and-invalid-examples.json`
- **Rationale:** Constructs are shown with valid worked examples and invalidity is well-described in prose, but the requirement asks for worked examples that include invalid counter-examples, which the spec does not provide. Falls short of target; holds the floor because invalidity is at least specified textually.

### the format specifies its core and how it extends and evolves

- **State:** active
- **Target:** Format specification
- **Rating:** 🔵 Target
- **Assessment record:** `assessments/008-format-spec-the-format-specifies-its-core-and-how-it-extends-and-evolves.json`
- **Rationale:** The spec names the minimal core, how authors add factors/keys/sections, how readers treat unrecognized content, and how the format versions forward. Satisfies the requirement; held below outstanding because versioning detail is delegated to an external policy doc.

### the format specification is well-structured and readable

- **State:** active
- **Target:** Format specification
- **Rating:** 🔵 Target
- **Assessment record:** `assessments/009-format-spec-the-format-specification-is-well-structured-and-readable.json`
- **Rationale:** Logical ordering, scannable annotated schema blocks, plain prose, and copy-ready examples make the spec navigable and readable. Satisfies the requirement.

### the README says what QUALITY.md is and who it's for

- **State:** active
- **Target:** README
- **Rating:** 🔵 Target
- **Assessment record:** `assessments/010-readme-the-readme-says-what-quality-md-is-and-who-it-s-for.json`
- **Rationale:** A first-time reader learns within the opening lines what a QUALITY.md file is, the problem it solves, and who it is for. Satisfies the requirement.

### the README shows the format and its payoff by example

- **State:** active
- **Target:** README
- **Rating:** 🟡 Minimum
- **Assessment record:** `assessments/011-readme-the-readme-shows-the-format-and-its-payoff-by-example.json`
- **Rationale:** The format is shown by a realistic example, but the requirement also asks for what running qualitymd produces, which appears only as prose, not rendered output. Falls short of target; holds the floor because the example itself is strong.

### the README gets a newcomer to a first result quickly

- **State:** active
- **Target:** README
- **Rating:** 🟡 Minimum
- **Assessment record:** `assessments/012-readme-the-readme-gets-a-newcomer-to-a-first-result-quickly.json`
- **Rationale:** A newcomer gets a copyable install-then-command sequence, but the requirement also asks for representative output and CI exit-code behavior, which are absent. Falls short of target; holds the floor because the command path itself is short and correct.

### the README reflects what the CLI and spec actually provide

- **State:** active
- **Target:** README
- **Rating:** 🔵 Target
- **Assessment record:** `assessments/013-readme-the-readme-reflects-what-the-cli-and-spec-actually-provide.json`
- **Rationale:** Every command and capability shown matches what the CLI provides today, planned/WIP state is marked, and the project's key risk (overstating what exists) is avoided. Satisfies the requirement; held below outstanding by the minor prose/table mismatch on version/upgrade.

### the CLI follows its functional specifications

- **State:** active
- **Target:** qualitymd CLI
- **Rating:** 🟡 Minimum
- **Assessment record:** `assessments/014-cli-the-cli-follows-its-functional-specifications.json`
- **Rationale:** The implementation follows the functional specs broadly and is well-tested, but a confirmed determinism MUST-violation in the show-status --json payload (machine-varying absolute path) plus an exit-code inconsistency for 'lint -' sit in the automation surface the specs govern. The verified determinism defect binds the rating to the floor.

### the CLI follows the project CLI design guide

- **State:** active
- **Target:** qualitymd CLI
- **Rating:** 🔵 Target
- **Assessment record:** `assessments/015-cli-the-cli-follows-the-project-cli-design-guide.json`
- **Rationale:** Stream separation, exit-code categories, conditional styling, non-interactivity, examples-first help, and consistent grammar follow the design guide closely. Satisfies the requirement; held below outstanding by a missing --no-color flag and the determinism deviation carried from the functional-spec finding.

## Findings

- `assessments/001-format-spec-the-format-specification-is-complete.json` at `SPECIFICATION.md:123-282`: The Frontmatter Schema section gives annotated YAML skeletons for Model, Rating Scale, Target, Factor, and Requirement, each stating requiredness, uniqueness, conditional cardinality ('at least two', 'exactly one'), and defaults (source inheritance). An implementer could build a parser and an author write a valid file from these.
- `assessments/001-format-spec-the-format-specification-is-complete.json` at `SPECIFICATION.md:112-130`: Malformed/omitted handling is stated: a document MUST begin with valid frontmatter containing a conforming Model; null/empty required values MUST be treated as absent; YAML that parses but does not conform is non-conforming.
- `assessments/001-format-spec-the-format-specification-is-complete.json` at `SPECIFICATION.md:284-291`: The Markdown body is fully open (no required section names/order); recommended sections are non-normative, so a valid file may have an empty body.
- `assessments/001-format-spec-the-format-specification-is-complete.json` at `SPECIFICATION.md:137-152`: Held below outstanding: scalar field types are abbreviated as <string>/<level-name> without character-set, whitespace, or length bounds, leaving some edge cases undefined.
- `assessments/002-format-spec-the-format-specification-admits-a-single-interpretation.json` at `SPECIFICATION.md:26-31`: BCP 14 keywords are defined once and used in all-caps only with normative force, per the RFC 2119/8174 convention; obligation strength is explicit (MUST/MUST NOT/SHOULD/MAY) throughout.
- `assessments/002-format-spec-the-format-specification-admits-a-single-interpretation.json` at `SPECIFICATION.md:160-170`: Quantifiers are bounded ('at least two Rating Levels', 'exactly one assessment'), so rules do not lean on vague quantifiers without a stated bound.
- `assessments/002-format-spec-the-format-specification-admits-a-single-interpretation.json` at `SPECIFICATION.md:199-213`: Permissive lowercase 'can'/'may' phrasing is used for non-normative latitude and is distinguishable from normative all-caps keywords; each rule reads one way.
- `assessments/003-format-spec-the-format-specification-separates-rules-from-rationale.json` at `SPECIFICATION.md:33-37`: An explicit normativity rule states all content is normative except passages marked non-normative; examples and notes are marked non-normative.
- `assessments/003-format-spec-the-format-specification-separates-rules-from-rationale.json` at `SPECIFICATION.md:258-261`: Inline notes are explicitly flagged ('Note: This note is non-normative'), and appendices are marked 'This appendix is non-normative' (lines 448, 475).
- `assessments/003-format-spec-the-format-specification-separates-rules-from-rationale.json` at `SPECIFICATION.md:176-181`: Rationale ('A Rating Level's description and criterion have distinct semantics') is given alongside the binding rule without a rule appearing only inside an example.
- `assessments/004-format-spec-the-format-specification-defines-its-terms-before-use.json` at `SPECIFICATION.md:64-103`: A dedicated Terminology section defines every model term (Quality Model, Entity, Model, Target, Source, Factor, Requirement, Assessment, Finding, Rating Scale, Rating Level, Rating Result, Evaluation Report) before any normative section uses them.
- `assessments/004-format-spec-the-format-specification-defines-its-terms-before-use.json` at `SPECIFICATION.md:123-282`: The Frontmatter Schema, Evaluation Semantics, and Report Semantics sections that follow rely only on terms already defined in Terminology; no rule introduces an undefined technical term.
- `assessments/005-format-spec-the-format-specification-is-internally-consistent.json` at `SPECIFICATION.md:165-170`: The schema marks Rating Level 'title' as MUST (a non-empty scalar human-readable label) and Factor 'title' as Required (line 220).
- `assessments/005-format-spec-the-format-specification-is-internally-consistent.json` at `SPECIFICATION.md:481-493`: Appendix B's minimal example omits the required 'title' on both rating levels (target/unacceptable carry only level/description/criterion) and on the 'reliability' factor, contradicting the rule the example illustrates.
- `assessments/005-format-spec-the-format-specification-is-internally-consistent.json` at `SPECIFICATION.md:64-103`: Outside that example, terminology is used consistently and no two normative statements contradict each other; the Appendix A suggested scale agrees with the schema.
- `assessments/006-format-spec-each-rule-is-observable-or-testable.json` at `SPECIFICATION.md:44-62`: Conformance Classes define observable behavior for documents, parsers, linters, evaluators, and report renderers, each phrased as a checkable acceptance/rejection criterion.
- `assessments/006-format-spec-each-rule-is-observable-or-testable.json` at `SPECIFICATION.md:250-282`: Structural requirements (assessment must be a single non-empty scalar; direct-target requirements must declare resolvable factors) map to observable file properties a linter can decide.
- `assessments/006-format-spec-each-rule-is-observable-or-testable.json` at `SPECIFICATION.md:384-389`: Some roll-up requirements are softer ('SHOULD include a rationale'), but they are recommendations, not normative obligations, so they do not undermine decidable conformance.
- `assessments/007-format-spec-the-format-s-constructs-are-shown-with-valid-and-invalid-examples.json` at `SPECIFICATION.md:446-498`: Valid worked examples are provided: Appendix A suggested rating scale and Appendix B a minimal complete document.
- `assessments/007-format-spec-the-format-s-constructs-are-shown-with-valid-and-invalid-examples.json` at `SPECIFICATION.md:250-278`: Invalidity is described only in prose ('A missing, empty, null, or list-valued assessment is invalid'; 'Missing factors, factors: null, factors: [] ... do not satisfy'); there are no worked invalid counter-example blocks showing what a malformed construct looks like.
- `assessments/008-format-spec-the-format-specifies-its-core-and-how-it-extends-and-evolves.json` at `SPECIFICATION.md:154`: The minimal core is named: an entry on factors, requirements, or targets MUST be supplied, with required fields marked throughout the schema.
- `assessments/008-format-spec-the-format-specifies-its-core-and-how-it-extends-and-evolves.json` at `SPECIFICATION.md:432-444`: An Extensions section says applications may add frontmatter properties, report fields, filters, and formats, constrained so they MUST NOT change conforming meaning and SHOULD use non-conflicting names.
- `assessments/008-format-spec-the-format-specifies-its-core-and-how-it-extends-and-evolves.json` at `SPECIFICATION.md:442-444`: Treatment of unrecognized content is defined: tools MUST preserve body content they do not understand and should preserve unrecognized extension fields.
- `assessments/008-format-spec-the-format-specifies-its-core-and-how-it-extends-and-evolves.json` at `SPECIFICATION.md:16-19`: Forward versioning is addressed via the specification-version field referencing the versioning policy; held from outstanding because the versioning mechanics live in an external doc rather than the spec body.
- `assessments/009-format-spec-the-format-specification-is-well-structured-and-readable.json` at `SPECIFICATION.md:21-103`: Sections introduce concepts before dependent detail: Conformance and Terminology precede Frontmatter Schema, Evaluation Semantics, and Report Semantics.
- `assessments/009-format-spec-the-format-specification-is-well-structured-and-readable.json` at `SPECIFICATION.md:137-248`: Annotated YAML skeletons with inline Required/Optional/Recommended comments make field structure scannable in place of dense prose.
- `assessments/009-format-spec-the-format-specification-is-well-structured-and-readable.json` at `SPECIFICATION.md:446-498`: The document carries minimal, realistic examples (suggested scale, minimal document) that a reader can copy and adapt; prose is plain.
- `assessments/010-readme-the-readme-says-what-quality-md-is-and-who-it-s-for.json` at `README.md:3-5`: The opening sentence states what it is: 'QUALITY.md is an agent-friendly file format and companion agent skill and CLI for continuously improving the quality of coding agent and AI assistant projects/harnesses.'
- `assessments/010-readme-the-readme-says-what-quality-md-is-and-who-it-s-for.json` at `README.md:7-17`: A 'Why QUALITY.md' section names the problem solved (technical, cognitive, and intent debt) in the opening lines.
- `assessments/010-readme-the-readme-says-what-quality-md-is-and-who-it-s-for.json` at `README.md:3-17`: The audience is identified (coding-agent/AI-assistant projects and software development teams).
- `assessments/011-readme-the-readme-shows-the-format-and-its-payoff-by-example.json` at `README.md:61-118`: A realistic 'Example QUALITY.md' (Support Inbox) excerpt is shown before reference detail.
- `assessments/011-readme-the-readme-shows-the-format-and-its-payoff-by-example.json` at `README.md:120-121`: The payoff is described only in prose ('An agent that reads this file can evaluate support work ... produce findings, and rate the results'); the README never shows what running qualitymd against the example produces (no lint output, init receipt, or report excerpt).
- `assessments/012-readme-the-readme-gets-a-newcomer-to-a-first-result-quickly.json` at `README.md:22-35`: A short install sequence is given (npx skills add; npm install -g).
- `assessments/012-readme-the-readme-gets-a-newcomer-to-a-first-result-quickly.json` at `README.md:223-230`: A 'Typical local loop' shows a copyable command sequence (spec, init, lint, status --json).
- `assessments/012-readme-the-readme-gets-a-newcomer-to-a-first-result-quickly.json` at `README.md:207-230`: No representative output is shown for any command, and CI exit-code behavior is mentioned only in passing ('can fail CI when ratings fall below your chosen bar') without an example, so a newcomer runs commands blind to what success looks like.
- `assessments/013-readme-the-readme-reflects-what-the-cli-and-spec-actually-provide.json` at `README.md:212-221`: The common-commands table matches shipped commands verified via 'qualitymd --help': spec, init, lint, lint --fix, status --json, version --json, upgrade --check, <command> --help all exist.
- `assessments/013-readme-the-readme-reflects-what-the-cli-and-spec-actually-provide.json` at `README.md:200-202`: The CLI is clearly marked 'an early work in progress', but its 'Today the binary ships' sentence lists only init/lint/spec/status and the evaluation surface, omitting version and upgrade, which do ship and appear in the table below it (internal understatement, not an overstatement).
- `assessments/013-readme-the-readme-reflects-what-the-cli-and-spec-actually-provide.json` at `README.md:207-208`: The claim that the CLI 'can fail CI when ratings fall below your chosen bar' matches the verified build-report --fail-at-or-below gate.
- `assessments/014-cli-the-cli-follows-its-functional-specifications.json` at `internal/cli/evaluation.go:129-149`: 'qualitymd evaluation show-status <run> --json' emits a machine-varying absolute path in the 'path' field and inside nextActions[].command even when given a repo-relative run path (verified: a relative input was echoed as /Users/craig/Code/qualitymd/quality.md/quality/evaluations/...). This violates payload determinism for a machine-consumed field.
- `assessments/014-cli-the-cli-follows-its-functional-specifications.json` at `internal/cli/lint.go:31-33`: 'qualitymd lint -' returns a plain error and exits 70 (internal error), while the sibling 'qualitymd status -' uses usageError and exits 2 (verified by running both). A bad stdin argument should be a usage error per the exit-code table.
- `assessments/014-cli-the-cli-follows-its-functional-specifications.json` at `internal/cli/evaluation.go:60-127`: Command behavior, flags, records, and exit codes otherwise match specs/cli.md and the specs/cli/ sub-specs across init, lint, spec, status, create-run, add-record, set-planned-coverage, show-status, build-report, version, and upgrade, with strong *_test.go coverage of stream separation, exit categories, and JSON shapes.
- `assessments/015-cli-the-cli-follows-the-project-cli-design-guide.json` at `internal/cli/root.go:19-24`: Stable exit-code categories 0/1/2/70 match the design guide's binding table; payload goes to stdout and confirmations/footers/JSON errors to stderr across init, lint, status, and evaluation commands.
- `assessments/015-cli-the-cli-follows-the-project-cli-design-guide.json` at `internal/cli/style.go:33-42`: Styling is terminal-only and additive over a canonical plain form and honors NO_COLOR; --json is spelled identically everywhere with the documented 'spec' carve-out; commands hold a consistent noun-verb grammar (evaluation create-run, etc.).
- `assessments/015-cli-the-cli-follows-the-project-cli-design-guide.json` at `internal/cli/evaluation.go:198-210`: Commands are non-interactive: readPayload fails with a usage diagnostic when stdin is a TTY and no --file is given, and '-' stdin is supported where applicable; help leads with examples on every leaf command.
- `assessments/015-cli-the-cli-follows-the-project-cli-design-guide.json` at `docs/guides/cli-design.md:151`: Deviations: no --no-color flag (only the NO_COLOR env var is honored) though the guide lists --no-color among conventional flags; the determinism principle (lines 224-232) is undercut by the show-status absolute-path payload noted in the functional-spec assessment.

## Advice

- [001-make-show-status-json-emit-repo-relative-paths](recommendations/001-make-show-status-json-emit-repo-relative-paths.md) [active] - show-status --json emits a repo-relative path for both 'path' and nextActions[].command given a relative run argument, and the requirement 'the CLI follows its functional specifications' reaches at least target. Re-evaluate in a new numbered run.
- [002-map-lint-to-a-usage-error-exit-2](recommendations/002-map-lint-to-a-usage-error-exit-2.md) [active] - qualitymd lint - exits 2 with a usage diagnostic, matching status; a test asserts the exit code. Re-evaluate in a new numbered run.
- [003-show-what-running-qualitymd-produces-not-just-the-model](recommendations/003-show-what-running-qualitymd-produces-not-just-the-model.md) [active] - The requirement 'the README shows the format and its payoff by example' reaches at least target: the README displays representative qualitymd output produced from the shown example. Re-evaluate in a new numbered run.
- [004-get-a-newcomer-to-a-visible-first-result-with-output-and-ci-exit-codes](recommendations/004-get-a-newcomer-to-a-visible-first-result-with-output-and-ci-exit-codes.md) [active] - The requirement 'the README gets a newcomer to a first result quickly' reaches at least target: representative output and CI exit-code behavior are shown. Re-evaluate in a new numbered run.
- [005-fix-the-appendix-b-minimal-example-to-satisfy-required-fields](recommendations/005-fix-the-appendix-b-minimal-example-to-satisfy-required-fields.md) [active] - The Appendix B example includes all required title fields and would pass lint; the requirement 'the format specification is internally consistent' reaches at least target. Re-evaluate in a new numbered run.
- [006-add-worked-invalid-counter-examples-to-the-specification](recommendations/006-add-worked-invalid-counter-examples-to-the-specification.md) [active] - The spec shows worked invalid counter-examples for its key constructs; the requirement 'the format's constructs are shown with valid and invalid examples' reaches at least target. Re-evaluate in a new numbered run.
