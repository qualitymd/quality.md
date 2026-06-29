---
type: Functional Specification
title: Evaluation Enum Catalogs - functional spec
description: Requirements for centralizing fixed Evaluation enum values and display metadata.
tags: [evaluation, schema, reports, enums]
timestamp: 2026-06-29T00:00:00Z
---

# Evaluation Enum Catalogs - functional spec

Companion to
[Evaluation Enum Catalogs](../0173-evaluation-enum-catalogs.md). This spec
states the delta for fixed Evaluation vocabularies in structured data and
generated reports. The durable source of truth is absorbed into the Evaluation
record, report, CLI data, and `/quality` skill specs.

The key words **MUST**, **MUST NOT**, and **SHOULD** are to be interpreted as
described in BCP 14 when, and only when, they appear in all capitals.

## Background / Motivation

Evaluation payloads are strict JSON contracts, while generated Markdown reports
are human review surfaces. Both need the same fixed vocabularies, but the
current implementation keeps enum literals, labels, markers, and report sort
orders in separate tables. That makes drift easy: a value can be accepted by the
CLI but render with a fallback label, or a report can carry a stale display
entry for a value that is no longer valid.

The change should preserve the stable persisted values and make the code treat
display metadata as part of the same contract as validation. Rating Level IDs
remain model-defined values resolved from the run's model snapshot, not fixed
Evaluation enums.

## Scope

Covered:

- fixed Evaluation data kind, report kind, status, confidence, finding type,
  finding severity, finding basis status, recommendation impact, finding
  ranking tier, finding coverage disposition, run gap kind, and rating-result
  kind vocabularies;
- CLI data validation and generated JSON Schema enum lists;
- generated Markdown report labels, markers, and enum ordering;
- data examples that use fixed enum values;
- checked-in report-gallery Markdown affected by generated report display
  changes;
- durable specs and skill guidance that name allowed fixed Evaluation values.

Deferred:

- persisted value renames;
- backwards-compatible aliases or normalization;
- model-defined Rating Level IDs and Rating Level display titles;
- user-configurable labels or markers;
- historical archives and released changelog entries.

## Requirements

1. Fixed Evaluation vocabularies **MUST** have one typed source for allowed
   persisted values.

   > Rationale: Validation, schema generation, examples, and rendering should
   > not copy string lists independently. Durable spec: modify
   > `specs/evaluation/records/payload-kinds.md`,
   > `specs/evaluation/records/json-conventions.md`, and
   > `specs/cli/evaluation-data.md`.

2. `qualitymd evaluation data set` and `qualitymd evaluation data verify`
   **MUST** reject out-of-vocabulary fixed enum values before a payload is
   written or reported valid.

   > Durable spec: modify `specs/cli/evaluation-data.md`.

3. Generated Evaluation JSON Schemas **MUST** derive fixed enum lists from the
   same fixed vocabulary source used by CLI validation.

   > Rationale: Agents often inspect schema before writing payloads; schema and
   > write validation must not disagree. Durable spec: modify
   > `specs/cli/evaluation-data.md`.

4. Every known fixed Evaluation enum value that can appear in generated
   Markdown reports **MUST** have an explicit human label and marker, except
   report frontmatter `type` strings, which remain plain artifact taxonomy.

   > Rationale: Reports are the human review surface, and every known value
   > should scan consistently instead of falling through to accidental
   > title-casing. Durable spec: modify
   > `specs/evaluation/reports/report-tree.md`.

5. Generated Markdown reports **MUST** render known fixed enum values with the
   shared marker-plus-label display and **MUST** preserve raw enum values in
   routine JSON, `EvaluationOutputResult`, schemas, and receipts.

   > Durable spec: modify `specs/evaluation/reports/report-tree.md`.

6. Generated Markdown reports **MUST** use shared rank metadata for ordered
   enum-like values where report ordering depends on value severity or impact.

   > Rationale: Display labels and sort behavior are two facets of the same
   > vocabulary; separating them caused previous severity and impact drift.
   > Durable spec: modify `specs/evaluation/reports/report-tree.md`.

7. The implementation **MUST NOT** accept aliases, case-insensitive values,
   display labels, markers, or legacy values as substitutes for canonical fixed
   enum values.

   > Rationale: Data set is a strict agent write boundary; friendly display
   > values are not parseable data. Durable spec: modify
   > `specs/cli/evaluation-data.md` and
   > `specs/evaluation/records/payload-kinds.md`.

8. Finding type display metadata **MUST** cover only active Finding Core `type`
   values unless a value is reintroduced into the payload contract.

   > Rationale: Stale display-only values make reports imply a broader contract
   > than the CLI validates. Durable spec: modify
   > `specs/evaluation/reports/report-tree.md`.

9. Report-reference kinds **MUST** include explicit labels and markers for
   `run`, `area`, `factor`, `requirement`, `findings`, `recommendations`, and
   `recommendation`.

   > Durable spec: modify `specs/evaluation/records/json-conventions.md` and
   > `specs/evaluation/reports/report-tree.md`.

10. `/quality` evaluation guidance **MUST** continue to author canonical enum
    values, not display labels or markers, when writing Evaluation data.

    > Durable spec: modify `specs/skills/quality-skill/evaluation.md`.

## Acceptance criteria

- Fixed enum values are declared as typed vocabularies rather than repeated raw
  string lists in data contracts and report sort/display helpers.
- `qualitymd evaluation data set --dry-run` rejects invalid values for status,
  confidence, finding type, severity, basis status, recommendation impact,
  ranking tier, coverage disposition, data kind, and report kind fields.
- `qualitymd evaluation data schema` emits enum lists that match the typed
  vocabularies.
- Report rendering uses explicit labels and markers for all known report-rendered
  fixed enum values, including basis status, ranking tier, coverage disposition
  when rendered, and recommendation/report kinds.
- Stale display-only finding type values are removed or made valid by the data
  contract.
- Focused Go tests and repository checks pass.
- `mise run report-gallery-check` passes with checked-in generated reports.

## Durable spec changes

### To add

None

### To modify

- `specs/evaluation/records/payload-kinds.md` - clarify fixed vocabularies,
  canonical values, and strict rejection of aliases/display values
  (requirements 1, 7).
- `specs/evaluation/records/json-conventions.md` - align report-reference kind
  and confidence vocabulary notes with fixed enum catalogs (requirements 1, 9).
- `specs/evaluation/reports/report-tree.md` - define shared marker-plus-label
  rendering and ordering for known fixed report-rendered enum values
  (requirements 4-6, 8-9).
- `specs/cli/evaluation-data.md` - require schema and data validation to derive
  enum sets from the same fixed vocabularies and reject aliases/display values
  (requirements 1-3, 7).
- `specs/skills/quality-skill/evaluation.md` - require skill-authored
  Evaluation data to use canonical fixed enum values (requirement 10).

### To rename

None

### To delete

None
