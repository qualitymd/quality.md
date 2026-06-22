---
type: Functional Specification
title: qualitymd lint output
description: Finding, JSON, repair-result, and human-output contract for qualitymd lint.
tags: [cli, command, lint, output]
timestamp: 2026-06-22T00:00:00Z
---

# qualitymd lint output

This spec owns the finding and output contract for [`qualitymd lint`](lint.md):
the JSON result shape, finding and repair objects, locations, and human output.
The rule system and rule catalog live in [`qualitymd lint rules`](lint-rules.md).
The command spec owns invocation, flags, repair execution, exit status, and
ordering behavior.

This document uses BCP 14 keywords only for testable conformance requirements.
The key words "MUST", "MUST NOT", "SHOULD", and "MAY" are to be interpreted as
described in [RFC 2119](../../docs/reference/rfc2119.md) and
[RFC 8174](../../docs/reference/rfc8174.md) when, and only when, they appear in
all capitals.

## Findings and output

`lint` emits zero or more **findings**. Human-readable output and JSON output
MUST report the same findings for the same input; `--json` changes only the
format.

### Finding schema

Under `--json`, `lint` MUST emit one JSON document on stdout with this shape:

```json
{
  "schemaVersion": 1,
  "path": "QUALITY.md",
  "valid": false,
  "summary": {
    "errors": 1,
    "warnings": 0,
    "info": 0,
    "fixable": 0,
    "fixed": 0
  },
  "findings": [
    {
      "ruleId": "missing-rating-scale",
      "severity": "error",
      "message": "The model root declares no `ratingScale`; a QUALITY.md model requires one rating scale.",
      "location": {
        "path": "QUALITY.md",
        "modelPath": ["ratingScale"],
        "label": "ratingScale"
      },
      "fixable": false
    }
  ],
  "repairs": [],
  "nextActions": []
}
```

Fields are stable public API:

- `schemaVersion` **MUST** be `1` until the JSON shape changes incompatibly.
- `path` **MUST** identify the linted file using the caller-facing path.
- `valid` **MUST** be `true` when there are no `error` findings and `false`
  otherwise.
- `summary.errors`, `summary.warnings`, and `summary.info` **MUST** equal the
  number of findings at each severity.
- `summary.fixable` **MUST** equal the number of findings whose `fixable` value
  is `true`.
- `summary.fixed` **MUST** equal the number of repairs applied during this run;
  it is `0` when `--fix` is not passed.
- `findings` **MUST** contain every emitted finding in deterministic order.
- `repairs` **MUST** contain every applied repair in deterministic order. It is
  an empty array when `--fix` is not passed or no repairs were applied.
- `nextActions` follows the [CLI spec's convention](../cli.md#conventions).

Each finding object **MUST** contain:

- `ruleId` — the stable rule id from the [Rules](lint-rules.md#rules) table.
- `severity` — one of `error`, `warning`, or `info`.
- `message` — the per-instance message described in [Messages](lint-rules.md#messages).
- `location` — the location object described in [Locations](#locations).
- `fixable` — whether this finding has a deterministic repair under
  [Fixability](lint-rules.md#fixability).

Each repair object **MUST** contain:

- `ruleId` — the rule id whose finding was repaired.
- `message` — a deterministic description of the edit that was applied.
- `location` — the location object for the repaired finding.

The JSON object **MUST NOT** include human styling, terminal control sequences,
or implementation-only fields. Additional documented fields can be added later
when they do not change the meaning of the required fields.

### Locations

A finding location names the smallest stable place the finding can be attached
to. It is not required to be a byte-perfect source span.

Under `--json`, `location` **MUST** contain:

- `path` — the linted file path, matching the top-level `path`.
- `modelPath` — an array path into the frontmatter model, using strings for map
  keys and numbers for list indexes. For example, the first rating level's
  missing `criterion` is `["ratingScale", 0, "criterion"]`; a requirement under
  a factor is `["factors", "<factor-name>", "requirements",
  "<requirement-statement>", "assessment"]`.
- `label` — a concise human-readable rendering of the same location, suitable
  for human output.

When a finding attaches to the frontmatter block as a whole,
`modelPath` **MUST** be an empty array and `label` **MUST** be `frontmatter`.
When a finding attaches to a missing key, `modelPath` **MUST** include the
missing key at the place it would appear.

If source position is available, `location` **SHOULD** also include 1-based
`line` and `column` fields for the start of the relevant YAML node. Source
positions are advisory: callers must treat `modelPath` as the stable machine
location.

### Human output

Human-readable output **MUST** include, for each finding, the severity, rule id,
message, and location label. It should summarize the total errors and
warnings. When `--fix` applies repairs, human-readable output **MUST** also
report how many repairs were applied. Styling and exact layout are governed by
the CLI's output conventions and are not part of this sub-spec.

When the result is invalid, human-readable output **MUST** render a
deterministic next-action footer on stderr using the same action data emitted in
the JSON result. The footer should prefer `qualitymd lint --fix <path>` when
at least one remaining finding is fixable, and otherwise should point to
rerunning `qualitymd lint <path>`.

When there are no findings, human-readable output should report that the file
is valid. Under `--json`, a valid file is represented by `"valid": true`, zero
counts in `summary.errors`, `summary.warnings`, and `summary.info`, and an empty
`findings` array.
