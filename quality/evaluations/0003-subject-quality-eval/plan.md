# Evaluation plan

## Effort

`standard` — every in-scope requirement assessed with targeted evidence; the
rating-binding findings re-checked before reporting.

## Requirement coverage (15 requirements across 3 targets)

### format-spec (`SPECIFICATION.md`) — 9 requirements

1. the format specification is complete
2. the format specification admits a single interpretation
3. the format specification separates rules from rationale
4. the format specification defines its terms before use
5. the format specification is internally consistent
6. each rule is observable or testable
7. the format's constructs are shown with valid and invalid examples
8. the format specifies its core and how it extends and evolves
9. the format specification is well-structured and readable

### readme (`README.md`, factor approachability) — 4 requirements

10. the README says what QUALITY.md is and who it's for
11. the README shows the format and its payoff by example
12. the README gets a newcomer to a first result quickly
13. the README reflects what the CLI and spec actually provide

### cli (`internal/cli`) — 2 requirements

14. the CLI follows its functional specifications
15. the CLI follows the project CLI design guide

## Analysis targets

`format-spec`, `readme`, `cli`, and the root model (empty `targetPath`).

## Rating-binding re-check

The following findings bind a headline rating and were independently re-verified
before writing report records:

- `cli` / determinism + consistency: `show-status --json` absolute path vs
  `create-run` relative path (re-run against the binary; checked against
  `create-run` and `cli.md` MUSTs).
- `cli` / consistency: empty-run gap kind `missing-analysis` vs spec-mandated
  `missing-root-analysis` (re-run `show-status`/`build-report` on an empty run;
  checked against `show-status` spec).
- `readme` / approachability: absence of any rendered `qualitymd` output sample
  in `README.md` (grep-confirmed) binding requirements 11 and 12.

## Out of scope / deferred areas

- Project dependencies not owned (Go toolchain, Cobra/Fang, release tooling).
- The `/quality` skill prompts and modes.
- Exhaustive README-schema-vs-SPECIFICATION field diff (spot-checked only).
