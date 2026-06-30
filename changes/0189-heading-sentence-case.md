---
type: Change Case
title: Heading sentence case
description: Align active Markdown and MDX headings with the repository sentence-case convention while preserving proper names, model titles, and historical records.
status: Design
tags: [docs, specs, reports, skill, mintlify]
timestamp: 2026-06-30T00:00:00Z
---

# Heading sentence case

The repository already says headings should use sentence case, but active docs,
specs, runtime skill files, generated report artifacts, and generated Mintlify
spec docs still contain many title-cased headings. This case turns that
convention into a scoped cleanup across active surfaces.

- [Functional spec](0189-heading-sentence-case/spec.md) - what the case must do.
- [Design doc](0189-heading-sentence-case/design.md) - how the alignment will be
  applied without introducing automation.

## Motivation

Heading case is a small presentation rule, but it is visible across every
reader-facing surface. Inconsistent title case makes the documentation set feel
less deliberate, obscures which names are true proper nouns or formal concepts,
and lets generated report artifacts drift from the same style expected of
authored docs.

The desired rule is not "lowercase everything." Proper names, acronyms, formal
QUALITY.md terms used as type names, user/model-provided titles, and command
names keep their intended casing. The cleanup should make ordinary prose in
headings sentence-cased while preserving those names.

## Scope

**Covered:** active Markdown and MDX headings in repository docs, guides,
durable specs, the formal `SPECIFICATION.md`, generated Mintlify specification
output, runtime `/quality` skill content, generated report artifacts, checked-in
report-gallery output, top-level docs, and active `.quality` evaluation
artifacts when they are maintained as current project records.

**Deferred:** archived Change Cases under `changes/archive/`, archived
evaluation records under `.quality/evaluations/archive/`, broad table-header
normalization, frontmatter `type` normalization, and any automated heading-case
lint or checker.

## Affected artifacts

**Repository guidance**

- `AGENTS.md` - broadens the heading-capitalization rule to every active
  Markdown/MDX surface and records the historical/archive boundary.
- `docs/guides/work-with-okf.md` - updates OKF heading examples and concept
  authoring guidance.
- `docs/guides/write-functional-specs.md` - updates spec section examples and
  durable-spec-change examples.
- `docs/guides/write-design-docs.md` - updates design-doc section examples.
- `docs/guides/work-with-change-cases.md` - updates Change Case examples and
  guidance.
- `docs/guides/reporting-design.md` - updates generated report heading
  guidance and examples.
- `docs/guides/cut-a-release.md` - updates release-guide headings and
  release-note heading examples.

**Durable specs**

- `SPECIFICATION.md` - sentence-cases active formal spec headings while
  preserving formal concept names and appendix labels.
- `specs/` - sentence-cases active durable spec headings, including report-tree,
  evaluation, CLI, and skill specs.
- `specs/log.md` and nested active `log.md` files - may receive a current log
  entry; historical dated entries are not rewritten unless they are current
  headings maintained as live navigation.

**Generated docs**

- `mintlify/specification.mdx` - regenerated from `SPECIFICATION.md`.
- Other `mintlify/*.mdx` files - reviewed and updated only where active headings
  are out of alignment.

**Bundled skill**

- `skills/quality/` - sentence-cases runtime skill, workflow, guide, resource,
  and active log/index headings where those files are maintained as current
  operational guidance.
- `specs/skills/quality-skill/` - sentence-cases the matching durable skill
  specs.

**Generated report artifacts**

- `internal/evaluation/report_tree.go` and affected tests - update generated
  report section labels and Contents labels to sentence case.
- `docs/guides/reporting-design.md`, `specs/evaluation/reports/report-tree.md`,
  and `specs/cli/evaluation-report.md` - update generated report contracts and
  examples.
- `examples/report-gallery/software-service/.quality/evaluations/0001-full-eval/`
  - regenerated checked-in gallery output.

**Top-level and install docs**

- `install.md`, `QUALITY.md`, `CHANGELOG.md`, and any other active root docs
  with out-of-alignment headings. Historical release sections in `CHANGELOG.md`
  are treated as append-only unless a current unreleased heading is being edited.

**Not affected**

- No new script, lint rule, CI check, pre-commit hook, or generated allowlist is
  added by this case.

## Status

`Design`. The functional requirements and implementation approach are captured;
code and document edits have not started. See the [status lifecycle](index.md#status-lifecycle).
