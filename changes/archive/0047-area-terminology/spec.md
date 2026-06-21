---
type: Functional Specification
title: Area terminology changeover - functional spec
description: Replace the formal Target model-node vocabulary with Area everywhere the live schema, records, reports, CLI, skill, scaffold, examples, and docs expose it.
tags: [terminology, schema, evaluation, report, cli, skill]
timestamp: 2026-06-21T00:00:00Z
---

# Area terminology changeover - functional spec

Companion to [Area terminology changeover](../0047-area-terminology.md). This
spec states the no-backward-compatibility delta for replacing the formal
QUALITY.md model-node term **Target** with **Area**.

This document uses BCP 14 keywords only for testable conformance requirements.
The key words "MUST", "MUST NOT", "SHOULD", and "MAY" are to be interpreted as
described in [RFC 2119](../../../docs/reference/rfc2119.md) and
[RFC 8174](../../../docs/reference/rfc8174.md) when, and only when, they appear in
all capitals.

## Background / Motivation

The current recursive model-node term, Target, collides with the default
`Target` rating level and produces awkward report and documentation language.
Area gives authors and readers a friendlier noun for "the part of the model being
evaluated" while still supporting code, docs, systems, services, test suites,
and other evaluatable slices. The vocabulary should stay small: Area is the
canonical concept, and the top-level Area is the root area.

This change intentionally does not preserve the old schema names. The project is
still in draft format territory, and carrying both vocabularies would make the
spec, skill, records, and examples harder to learn.

## Scope

Covered: model schema terminology, machine field names, record names, report
model fields, report Markdown headings, CLI help/diagnostics, evaluation run
naming, lint/status semantics, `/quality` skill guidance, scaffold comments,
examples, fixtures, dogfood model, generated package README, and public docs.

Non-goals: keeping `targets:` or `targetPath` accepted as legacy inputs,
auto-migrating old evaluation runs, adding aliases, or renaming the default
rating level `target` / `Target`.

## Requirements

The formal recursive model-node concept **MUST** be named **Area**.

The document's top-level Area **MUST** be called the **root area** in formal spec
and report prose.

User-facing prose **MUST NOT** introduce a second canonical noun for Area.
Contextual phrases such as "the whole repository" or "everything covered by this
QUALITY.md" may describe a concrete root area, but the named concept remains
Area.

Human-facing report and CLI output **MUST NOT** use Subject as a canonical label
for the evaluated thing. Labels that currently say Subject for the evaluated
model node **MUST** become root area or a direct file/model label, depending on
whether they name the Area or the `QUALITY.md` file being read.

The frontmatter schema **MUST** use `areas:` for child Area maps. Each Area entry
MUST keep the current recursive node semantics: title, description, factors,
requirements, child areas, and source.

The frontmatter schema **MUST NOT** accept `targets:` as a conforming property.

> Rationale: This is a full draft-format changeover. Accepting both names would
> teach two shapes and would keep old terminology alive in examples, lint
> behavior, and agent guidance.

The format specification **MUST** distinguish Area from Factor: an Area is what
is evaluated; a Factor is the quality lens used to organize Requirements for an
Area.

The `Source` definition **MUST** describe the material assessed for an Area.

The Area selector property **MUST** remain `source`. This change **MUST NOT**
rename `source` to material, content, evidence, scope, location, input,
resource, corpus, basis, or selector.

> Rationale: `source` is the shortest field name that still reads naturally for
> paths, globs, URLs, and other selectors. Evidence is a report/assessment
> concept, scope is already an evaluation concept, and material/content are less
> precise as schema keys.

User-facing prose **SHOULD** distinguish `source` from source code when ambiguity
is likely. Use `Source` or `` `source` `` for the Area property, and "source
code" for code artifacts.

Evaluation scope and narrowing **MUST** use Area terminology. A narrowed
evaluation can select an Area and its subtree, a Factor, or an Area/Factor pair.

Evaluation record JSON **MUST** rename `targetPath` to `areaPath` everywhere the
ordered model-node path appears, including assessment result records, analysis
records, planned coverage frontmatter, report JSON, and any related status gaps.

Evaluation analysis records **MUST** be described as one JSON file per Area, and
analysis filenames **SHOULD** use Area-derived slugs.

Report JSON **MUST** rename Target-derived collections and fields to
Area-derived names, including target summaries/details and any TargetPath typed
state.

Human Markdown reports **MUST** use Area labels and headings, including `Area
Ratings` or equivalent, `Area Details`, Area rows, and Area-specific requirement
labels.

`report-summary.md` **MUST** show the root area row when a run contains a root
area analysis, even when the model has no child areas.

> Rationale: A fresh model can be only a root area with root-level Factors.
> Reporting "No area ratings were recorded" hides the most important rating from
> the summary.

CLI help, diagnostics, receipts, and examples **MUST** use Area terminology for
model-node selection, lint findings, status output, evaluation reports, and
record examples.

`qualitymd evaluation create` **MUST** use `--model <path>` as the flag for the
repository-relative `QUALITY.md` file to snapshot. It **MUST NOT** accept
`--subject`.

New evaluation run folder names **MUST NOT** include `subject` as an altitude
segment. An unnarrowed run should use `NNNN-quality-eval`; a narrowed run should
use `NNNN-<narrowing>-quality-eval`.

The `/quality` skill **MUST** parse and speak Area terminology for scope
resolution, run frames, wizard options, evaluation instructions, report
interpretation, and improvement recommendations.

The bundled scaffold **MUST** introduce optional child `areas:` and root-area
language; it **MUST NOT** include commented `targets:` guidance.

Maintained examples and fixtures **MUST** use `areas:` and `areaPath` in live
records and generated reports. Historical archived Change Cases and append-only
logs may retain old terminology when they describe past behavior.

The default rating level id/title `target` / `Target` **MUST NOT** be renamed by
this change.

> Rationale: The pain is the model-node term colliding with the rating label.
> Renaming the model node to Area removes that collision without changing the
> rating scale contract.

Old evaluation runs or records that still use `targetPath` **MUST** be treated as
historical incompatible records by current readers. Implementations **MUST NOT**
silently translate old record fields into the new report model.

## Durable spec changes

### To add

None

### To modify

- `SPECIFICATION.md` - replace the formal Target concept with Area, introduce
  root area, replace `targets:` with `areas:`, and update evaluation semantics
  for Area scope, Source resolution, roll-up, and report obligations (per the
  schema and evaluation requirements above).
- `specs/evaluation-records.md` - replace `targetPath`, target analysis records,
  target summaries/details, report sections, and planned coverage keys with
  Area equivalents, with no legacy-field acceptance (per the record/report
  requirements above).
- `specs/cli.md` and relevant `specs/cli/` children - update CLI contracts that
  expose model-node names in help, diagnostics, receipts, status, lint,
  evaluation create/list/status/report, run naming, the model-file selection
  flag, and examples (per the CLI requirement above).
- `specs/skills/quality-skill/quality-skill.md` and related
  `specs/skills/quality-skill/` guide/example contracts - align the skill's
  scope-resolution, reporting, and example-evaluation contracts with Area
  terminology (per the skill and fixture requirements above).

### To delete

None
