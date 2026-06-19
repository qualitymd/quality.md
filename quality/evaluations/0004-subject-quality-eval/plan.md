# Evaluation plan

## Effort

`standard` — every in-scope Requirement assessed against targeted evidence, with
the rating-binding findings re-checked before reporting.

## Coverage

All 15 in-scope Requirements across the three evaluable Targets:

### format-spec (9)

- target-level: "the format specification is complete"
- clarity: single interpretation; separates rules from rationale; defines terms
  before use
- consistency: internally consistent
- verifiability: each rule observable/testable; constructs shown with valid and
  invalid examples
- extensibility: specifies core and how it extends/evolves
- usability: well-structured and readable

### readme (4, all under Approachability)

- says what QUALITY.md is and who it's for
- shows the format and its payoff by example
- gets a newcomer to a first result quickly
- reflects what the CLI and spec actually provide

### cli (2 target-level)

- follows its functional specifications (specs/cli.md, specs/cli/)
- follows the project CLI design guide (docs/guides/cli-design.md)

## Analysis targets

Four analysis records: `format-spec`, `readme`, `cli`, and the structural root
(aggregate only).

## Limitations

- Standard effort: targeted evidence, not a full adversarial line-by-line read of
  every Go file. CLI evidence relied on a read-only subagent sweep plus direct
  orchestrator re-checks of the rating-binding findings.
- Roll-ups are judgment, not a numeric formula, per SPECIFICATION.md.
