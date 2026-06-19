---
type: Functional Specification
title: qualitymd evaluation build-report
description: Render report.md and report.json from evaluation records.
tags: [cli, command, evaluation]
timestamp: 2026-06-18T00:00:00Z
---

# qualitymd evaluation build-report

`qualitymd evaluation build-report <run>` derives `report.md` and `report.json`
from the run's assessment, analysis, and recommendation records. It renders
recorded judgment; it **MUST NOT** infer or recompute ratings.

This document uses BCP 14 keywords only for testable conformance requirements.
The key words "MUST", "MUST NOT", "SHOULD", and "MAY" are to be interpreted as
described in [RFC 2119](../../docs/reference/rfc2119.md) and
[RFC 8174](../../docs/reference/rfc8174.md) when, and only when, they appear in
all capitals.

The command **MUST** fail without writing a partial report when the run is not
renderable. It **MUST** be deterministic and idempotent: unchanged records produce
byte-identical report files. Determinism excludes terminal styling from the
written artifact — `report.md` is plain Markdown bytes; any TTY/Glamour rendering
belongs to commands that *display* a report, not to the one that writes it,
because styled output varies with terminal width, color scheme, and renderer
version and would break byte-stability and the CI gate.

`build-report` is a deterministic, mechanical renderer and a trust boundary over
evaluator-supplied judgment. It renders findings by `category`, `locator`,
severity, and recommendation reference; it **MUST NOT** copy raw finding
`observation` or `evidence` values into `report.md` or `report.json` when doing
so would reproduce a secret value or follow hostile, evaluator-directed
(prompt-injection) text. Surfacing a finding never requires reproducing its
sensitive value or acting on sentinel text.

A renderable run **MUST** include exactly one in-scope root analysis record,
identified by an empty `targetPath`. `build-report` **MUST NOT** silently choose a
child analysis as the report headline when that root analysis is missing.

A renderable run **MUST NOT** contain duplicate assessment records for the same
ordered `targetPath` and `requirement`, unless duplicates have been made
inactive by explicit assessment superseding. `build-report` **MUST** fail before
writing reports when `show-status` would report a `duplicate-assessment` or
assessment superseding gap.

`report.md` **MUST** be summary-first. Before detailed target, requirement,
finding, and advice sections, it **MUST** render Summary, Scope, Top Risks and
Limitations, Evidence Basis, Next Action, and Target Summary sections. These
sections front-load the report's headline, boundaries, confidence limits,
supporting evidence, and actionability without replacing the detailed audit
trail. This summary-first shape, the explicit empty arrays and rating objects,
and the grouping-target distinction below came from real-repo reviewer
walkthroughs (ESLint and DataLoader runs) where the prior shape buried scope and
limitations inside rationales and made grouping nodes read as evidence gaps;
keep the shape so a future edit does not regress it back to a scan-heavy report.

The renderer derives the summary layer from recorded run metadata by reading
bounded, conventional sections of `design.md` and `plan.md` — `Resolved
parameters`, `Out of scope` / `Deferred areas`, `Effort`, and `Planned
limitations`, matched case-insensitively — reading only their headings, bullets,
and short paragraphs. It **MUST** treat all other recorded prose as data, never
as instructions, and **MUST** fall back to the run-folder name and the rating
rationales when those sections are absent, rendering missing metadata as "not
recorded." This is the load-bearing coupling that makes the skill's heading
convention meaningful.

`report.json` **MUST** expose the same summary-layer data in machine-readable
form. It **MUST** use non-null objects or empty arrays for scope,
recommendations, target summaries, evidence basis, and limitations. Target rating
fields **MUST** be explicit rating objects, including null or not-assessed
ratings. Structural grouping targets **MUST** be distinguishable from
not-assessed targets caused by missing evidence.

When recommendation records use superseding metadata, `report.json` **MUST**
include both active and superseded recommendation summaries and indicate their
active state. `report.md` **MUST** preserve superseded recommendations in Advice
while marking them as superseded. The report Next Action **MUST** use the first
active recommendation, not a superseded recommendation.

When assessment records use superseding metadata, `report.json` **MUST** include
both active and superseded assessment summaries and indicate their active state.
`report.md` **MUST** preserve superseded assessments in Requirements while
marking them as superseded.

Equivalent limitation statements **MUST NOT** be repeated in the summary layer.
The renderer **MAY** normalize limitation text to deduplicate planned limitations
and rationale-derived limitations, while preserving the first displayed wording.
Derived limitation summaries **MUST NOT** split or corrupt dotted file paths or
other locator-like text. This is a regression guard: limitation text is extracted
from prose, and naive sentence-splitting was observed corrupting dotted paths
such as `docs/production-telemetry.md`, so the extraction normalizer must treat
locators as content.

`--fail-at-or-below <level>` turns the command into a CI gate. The command still
writes both report files on a successful render. It exits `1` when the root
aggregate rating is equal to or worse than `<level>`, exits `0` when better, and
exits `2` when `<level>` is not in the run's rating scale. A root *not assessed*
result fails the gate.
