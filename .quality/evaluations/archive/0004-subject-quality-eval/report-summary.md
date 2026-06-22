# Quality Evaluation Summary

**Run:** `0004-subject-quality-eval`
**Subject:** `QUALITY.md`
**Scope:** standard` — every in-scope Requirement assessed against targeted evidence, with the rating-binding findings re-checked before reporting.
**Effort:** standard
**Root rating:** 🟡 Minimum
**Full report:** [report.md](report.md)
**Machine report:** [report.json](report.json)

## Headline

Structural root with no local requirements (source missing by design). Aggregate over three child targets: format-spec at target (strongest), readme at minimum, cli at minimum. Two of three deliverables sit at the acceptable floor, so the whole-project aggregate is minimum: the format specification is solid, while the README does not demonstrate the tool working and the CLI has a confirmed determinism MUST-violation in its automation surface.

## Top Risks

1. **low** - Held below outstanding: scalar field types are abbreviated as <string>/<level-name> without character-set, whitespace, or length bounds, leaving some edge cases undefined. (`assessments/001-format-spec-the-format-specification-is-complete.json` at `SPECIFICATION.md:137-152`)
2. **medium** - Appendix B's minimal example omits the required 'title' on both rating levels (target/unacceptable carry only level/description/criterion) and on the 'reliability' factor, contradicting the rule the example illustrates. (`assessments/005-format-spec-the-format-specification-is-internally-consistent.json` at `SPECIFICATION.md:481-493`)
3. **medium** - Invalidity is described only in prose ('A missing, empty, null, or list-valued assessment is invalid'; 'Missing factors, factors: null, factors: [] ... do not satisfy'); there are no worked invalid counter-example blocks showing what a malformed construct looks like. (`assessments/007-format-spec-the-format-s-constructs-are-shown-with-valid-and-invalid-examples.json` at `SPECIFICATION.md:250-278`)
4. **low** - Forward versioning is addressed via the specification-version field referencing the versioning policy; held from outstanding because the versioning mechanics live in an external doc rather than the spec body. (`assessments/008-format-spec-the-format-specifies-its-core-and-how-it-extends-and-evolves.json` at `SPECIFICATION.md:16-19`)
5. **medium** - The payoff is described only in prose ('An agent that reads this file can evaluate support work ... produce findings, and rate the results'); the README never shows what running qualitymd against the example produces (no lint output, init receipt, or report excerpt). (`assessments/011-readme-the-readme-shows-the-format-and-its-payoff-by-example.json` at `README.md:120-121`)

## Rating Summary

| Target               | Aggregate rating | Reason                                                                                                                                                                                                                                                                                                                                                                                                                                                         |
| -------------------- | ---------------- | -------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------- |
| QUALITY.md           | 🟡 Minimum       | Structural root with no local requirements (source missing by design). Aggregate over three child targets: format-spec at target (strongest), readme at minimum, cli at minimum. Two of three deliverables sit at the acceptable floor, so the whole-project aggregate is minimum: the format specification is solid, while the README does not demonstrate the tool working and the CLI has a confirmed determinism MUST-violation in its automation surface. |
| qualitymd CLI        | 🟡 Minimum       | Leaf target: aggregate equals local rating. Strong design-guide adherence and automation ergonomics are pulled to the floor by one confirmed determinism MUST-violation and an exit-code consistency gap.                                                                                                                                                                                                                                                      |
| Format specification | 🔵 Target        | Leaf target: aggregate equals local rating. The strongest target in the model; complete, clear, extensible, and usable, with two bounded documentation gaps in consistency and verifiability.                                                                                                                                                                                                                                                                  |
| README               | 🟡 Minimum       | Leaf target: aggregate equals local rating. The model's weakest target; the front door explains itself well but never demonstrates the tool working.                                                                                                                                                                                                                                                                                                           |

## Limitations

- Structural root with no local requirements (source missing by design)
- Normative statements and terminology are mutually consistent, but the requirement explicitly tests that every example agrees with the rule it illustrates, and the Appendix B minimal example violates required-field rules (missing title on rating levels and factor)
- A newcomer gets a copyable install-then-command sequence, but the requirement also asks for representative output and CI exit-code behavior, which are absent
- Satisfies the requirement; held below outstanding by a missing --no-color flag and the determinism deviation carried from the functional-spec finding

## Next Action

show-status --json emits a repo-relative path for both 'path' and nextActions[].command given a relative run argument, and the requirement 'the CLI follows its functional specifications' reaches at least target. Re-evaluate in a new numbered run.

See active recommendations:

- [001-make-show-status-json-emit-repo-relative-paths](recommendations/001-make-show-status-json-emit-repo-relative-paths.md) - show-status --json emits a repo-relative path for both 'path' and nextActions[].command given a relative run argument, and the requirement 'the CLI follows its functional specifications' reaches at least target. Re-evaluate in a new numbered run.
- [002-map-lint-to-a-usage-error-exit-2](recommendations/002-map-lint-to-a-usage-error-exit-2.md) - qualitymd lint - exits 2 with a usage diagnostic, matching status; a test asserts the exit code. Re-evaluate in a new numbered run.
- [003-show-what-running-qualitymd-produces-not-just-the-model](recommendations/003-show-what-running-qualitymd-produces-not-just-the-model.md) - The requirement 'the README shows the format and its payoff by example' reaches at least target: the README displays representative qualitymd output produced from the shown example. Re-evaluate in a new numbered run.
- [004-get-a-newcomer-to-a-visible-first-result-with-output-and-ci-exit-codes](recommendations/004-get-a-newcomer-to-a-visible-first-result-with-output-and-ci-exit-codes.md) - The requirement 'the README gets a newcomer to a first result quickly' reaches at least target: representative output and CI exit-code behavior are shown. Re-evaluate in a new numbered run.
- [005-fix-the-appendix-b-minimal-example-to-satisfy-required-fields](recommendations/005-fix-the-appendix-b-minimal-example-to-satisfy-required-fields.md) - The Appendix B example includes all required title fields and would pass lint; the requirement 'the format specification is internally consistent' reaches at least target. Re-evaluate in a new numbered run.
- [006-add-worked-invalid-counter-examples-to-the-specification](recommendations/006-add-worked-invalid-counter-examples-to-the-specification.md) - The spec shows worked invalid counter-examples for its key constructs; the requirement 'the format's constructs are shown with valid and invalid examples' reaches at least target. Re-evaluate in a new numbered run.
