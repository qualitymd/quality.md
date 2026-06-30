---
type: Functional Specification
title: Heading sentence case - functional spec
description: Requirements for aligning active Markdown and MDX headings with sentence case.
tags: [docs, specs, reports, skill, mintlify]
timestamp: 2026-06-30T00:00:00Z
---

# Heading sentence case - functional spec

Companion to the
[Heading sentence case](../0189-heading-sentence-case.md) change case. This
spec states *what* the heading-case alignment must do.

The key words **MUST**, **SHOULD**, and **MAY** are to be interpreted as
described in BCP 14 when, and only when, they appear in all capitals.

## Background / motivation

The repository's agent instructions already require sentence case for headings,
but the rule is both too narrowly scoped and inconsistently followed. The
current mismatch appears in authored docs, formal specs, runtime skill files,
generated report artifacts, and generated Mintlify spec output.

The cleanup needs to align visible headings without weakening meaningful casing.
Proper names, acronyms, command names, formal QUALITY.md terms used as type
names, appendix labels, and user/model-provided titles still carry their own
casing. The goal is consistent editorial style, not destructive recasing.

## Scope

Covered active surfaces:

- top-level docs such as `AGENTS.md`, `install.md`, `QUALITY.md`, and the
  current editable parts of `CHANGELOG.md`;
- `docs/` guides and reference docs;
- the formal `SPECIFICATION.md`;
- active durable specs under `specs/`;
- active runtime `/quality` skill content under `skills/quality/`;
- active Mintlify pages under `mintlify/`, with `mintlify/specification.mdx`
  regenerated from `SPECIFICATION.md`;
- generated Evaluation report headings and Contents labels emitted by the CLI;
- checked-in generated report-gallery output; and
- active `.quality` evaluation artifacts when maintained as current project
  records.

Out of scope:

- archived Change Cases under `changes/archive/`;
- archived evaluation records under `.quality/evaluations/archive/`;
- blanket normalization of table headers, frontmatter `type` values, enum
  display labels, YAML titles, or structured metadata;
- user/model-provided display titles embedded inside headings; and
- any automated heading-case lint, checker, CI gate, pre-commit hook, or
  allowlist.

## Requirements

### Convention scope

- Repository guidance **MUST** define the sentence-case heading convention for
  all active Markdown and MDX surfaces, not only README, docs, guides, and specs.
  > Rationale: The visible drift is cross-surface; keeping the rule scoped to a
  > subset leaves generated reports, runtime skill docs, and Mintlify output
  > outside the stated convention.
  >
  >> Durable spec: none. This is repository editorial guidance, not a format or
  >> tooling contract.

- Repository guidance **MUST** state that ordinary prose in headings uses
  sentence case: capitalize the first word, the first word after a colon, and
  proper names or proper nouns.
  > Rationale: The existing rule is correct but needs to be carried into the
  > broader scope.
  >
  >> Durable spec: none.

### Preserved casing

- Heading edits **MUST** preserve proper names, acronyms, command names, file
  names, literal values, and formal QUALITY.md concepts when used as type names
  or terms of art.
  > Rationale: Sentence case should distinguish real names from generic prose,
  > not erase meaningful casing.
  >
  >> Durable spec: modify `SPECIFICATION.md` and `specs/` headings only where
  >> current formal concept names are embedded in headings.

- Heading edits **MUST** preserve user-provided or model-provided display titles
  embedded in generated report headings, such as the `<Area title>` portion of
  `Area: <Area title>`.
  > Rationale: The model owns display titles; the report renderer should not
  > silently rewrite a user's Model labels.
  >
  >> Durable spec: modify `specs/evaluation/reports/report-tree.md` and
  >> `specs/cli/evaluation-report.md` to state the preserved-title boundary.

### Active authored content

- Active authored Markdown and MDX headings in top-level docs, `docs/`,
  `specs/`, `skills/quality/`, `mintlify/`, and active `changes/` files **MUST**
  be brought into sentence-case alignment.
  > Rationale: The repository presents these files as current sources of truth or
  > current operational guidance.
  >
  >> Durable spec: modify active `specs/` headings where out of alignment.

- Current generated `mintlify/specification.mdx` output **MUST** align with the
  `SPECIFICATION.md` source and **MUST NOT** be hand-edited as the source of
  truth.
  > Rationale: The generated page should stay a projection of the formal spec,
  > not become an independent docs surface.
  >
  >> Durable spec: none.

### Generated reports

- Generated Evaluation report section headings and generated Contents labels
  **MUST** use sentence case for generic section names.
  > Rationale: Reports are active reader-facing Markdown artifacts and should
  > follow the same heading convention as authored docs.
  >
  >> Durable spec: modify `specs/evaluation/reports/report-tree.md` and
  >> `specs/cli/evaluation-report.md` to require sentence-cased generated report
  >> section headings and Contents labels.

- Generated report metadata fields, report `type` values, enum display labels,
  table headers, and structured JSON/YAML field values **MUST NOT** be recased
  solely to satisfy the heading convention.
  > Rationale: Those surfaces have separate compatibility and readability jobs;
  > this case is scoped to headings and generated Contents labels.
  >
  >> Durable spec: modify `specs/evaluation/reports/report-tree.md` to keep the
  >> boundary explicit.

- Checked-in report-gallery output **MUST** be regenerated after renderer
  changes so examples match the current generated report contract.
  > Rationale: The gallery is checked-in example output and should not preserve
  > stale generated headings after the renderer changes.
  >
  >> Durable spec: none.

### Historical boundary

- Archived Change Cases under `changes/archive/` and archived evaluation records
  under `.quality/evaluations/archive/` **MUST NOT** be rewritten solely for
  heading sentence case.
  > Rationale: Those files are historical records; changing their headings would
  > create churn without improving current guidance or generated output.
  >
  >> Durable spec: none.

- Append-only logs and historical changelog sections **SHOULD** keep historical
  entries unchanged unless an entry is currently being edited for another
  reason; new or current log headings **MUST** follow the updated convention.
  > Rationale: Logs preserve past wording, while new navigation should stop
  > adding title-case drift.
  >
  >> Durable spec: none.

### No automation

- This change **MUST NOT** add a heading-case lint rule, checker script, CI gate,
  pre-commit hook, generated allowlist, or other automated enforcement.
  > Rationale: The user explicitly accepted the cleanup but rejected automation;
  > false positives around proper nouns and model titles are also likely.
  >
  >> Durable spec: none.

## Durable spec changes

### To add

None.

### To modify

- `specs/evaluation/reports/report-tree.md` - require sentence-cased generated
  report section headings and Contents labels; preserve model display titles and
  keep metadata/table-header surfaces out of scope.
- `specs/cli/evaluation-report.md` - align the CLI report-build contract with
  the generated report heading convention.
- Active `specs/` files with out-of-alignment headings - sentence-case headings
  without changing their substantive requirements.

### To rename

None.

### To delete

None.

## Validation check

The change is complete when:

- the affected active Markdown/MDX surfaces have been reviewed and aligned;
- `mintlify/specification.mdx` is regenerated from `SPECIFICATION.md`;
- checked-in report-gallery output is regenerated after report renderer changes;
- existing formatting, tests, report-gallery, generated-doc, and docs-link
  checks pass as applicable; and
- no heading-case automation has been added.
