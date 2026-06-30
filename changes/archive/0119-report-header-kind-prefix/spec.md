---
type: Functional Specification
title: Report header kind prefix and title-first layout
description: Prefix each generated Evaluation report's H1 with its kind, render the title before the navigation trails, drop the Path/Name identifier line, and lock the trail's root element to the model title.
status: Draft
tags: [evaluation, reports, markdown]
timestamp: 2026-06-26T00:00:00Z
---

# Report header kind prefix and title-first layout

This spec governs the **top of every generated Evaluation report** — the H1
title line, the kind it states, its position relative to the navigation trails,
and the identifier line beneath it. It is the delta for
[0119](../0119-report-header-kind-prefix.md). The durable contract it lands in is
[`specs/evaluation/reports/report-tree.md`](../../specs/evaluation/reports/report-tree.md)
(normative — its Navigation rules, per-report report lists, and Rendering Rules
bind here).

The key words **MUST**, **MUST NOT**, **SHOULD**, and **MAY** are to be
interpreted as described in BCP 14 when, and only when, they appear in all
capitals.

## Background / Motivation

Every report opens with its navigation trail and then a bare H1 — the subject's
display title with no statement of the report's _kind_. Because an Area, Factor,
and Requirement can share a title, a report opened cold or linked from outside
the run does not say plainly what it covers; the kind is only inferable from the
trail label and a secondary `Path:` / `Name:` line.

This change states the kind in the report's most prominent line — a kind prefix
on the H1 — and renders that title first, with the trails following it. That puts
the kind where the subject's identity already lives (the title), consistent with
the existing rule that the header _table_ must not repeat the kind as metadata
([0104](../../changes/archive/0104-evaluation-v2-report-header-navigation.md)).
With kind on the title and location on the trail, the `Path:` / `Name:`
identifier line is redundant — its remaining job, surfacing the canonical
structural ID, is already served by report filenames and `Data` links, which
derive from structural IDs (Report Paths), so the line is dropped.

The new title surfaces the root trail element more prominently, so the case also
pins that element to the Model `title` — already the renderer's behavior, now a
durable guard.

One consequence is accepted rather than fixed: on Area reports the kind prefix
and the trail label coincide (`# Area: X` above `Area: … / X`). Relabeling the
`Area:` trail is out of scope (see the
[parent scope](../0119-report-header-kind-prefix.md#scope)).

## Scope

Covered: the kind prefix on each report's H1, title-first ordering relative to
the trails, removal of the `Path:` / `Name:` line, and locking the trail's root
element to the Model `title`. Deferred / non-goals are recorded in the
[parent concept](../0119-report-header-kind-prefix.md#scope) — in short, no
relabeling of the `Area:` trail, no change to the no-title fallback, to trail
content/order/separators/targets, to summary-table columns or any other report
content, or to structured data, paths, filenames, and links.

## Requirements

### Kind-prefixed, title-first H1

Every generated report **MUST** render its H1 title line as the first content of
the report, before the `Area:` navigation trail and before any other context
line.

> > Rationale: a kind-prefixed title-first header states _what this is_ before
> > _where it sits_; leading with the trail made navigation chrome the opening
> > line. — 0119

Every generated report's H1 **MUST** be the subject's display title prefixed with
the report's kind label and a colon — `Area:` for the root and non-root Area
reports, `Factor:` for Factor reports, and `Requirement:` for Requirement
reports (for example, `# Requirement: Inputs are validated`).

> > Rationale: an Area, Factor, and Requirement can share a title; stating the
> > kind in the most prominent line lets a report opened cold or linked from
> > outside the run identify itself. The kind rides the title — where identity
> > already lives — not the state table. — 0104, 0119

### Trails follow the title

The `Area:` navigation trail, the Factor report `Factor:` trail, and the
Requirement report `Factors:` context line **MUST** render after the H1 title
line, preserving their existing relative order, content, separators, and link
targets.

> > Rationale: this change moves only the title's position relative to the trails;
> > the trail contract (root-through-current `Area:` links, the Factor trail, the
> > plural `Factors:` set) is unchanged. — 0119

### No separate identifier line

Reports **MUST NOT** render a separate `Path:` or `Name:` identifier line in the
header. The report's kind is carried by the H1 prefix and its location by the
navigation trail.

> > Rationale: with kind on the title and location on the trail, the line was
> > redundant. Its one unique payload — the canonical structural ID — is not lost:
> > report filenames and `Data` links already derive from structural IDs, so the
> > ID remains reachable without a duplicate header line. — 0119

### Trail root element

The `Area:` navigation trail's root element **MUST** render the Model `title`
when the model defines one.

> > Rationale: the title-first header surfaces the root element prominently (the
> > root report shows the model title in both its H1 and its trail). The renderer
> > already resolves this; pinning it to the durable contract guards against a
> > regression to a generic "Root" label. — 0119

## Durable spec changes

### To add

None

### To modify

- `specs/evaluation/reports/report-tree.md` — Navigation: replace the "Every
  report **MUST** start with an `Area:` navigation trail" rule with the
  title-first, kind-prefixed-H1 rule and the requirement that the trails follow
  the title; add the trail-root-element-renders-Model-`title` rule; add the
  prohibition on a `Path:` / `Name:` identifier line (per the kind-prefix,
  title-first, trails-follow, no-identifier-line, and trail-root requirements
  above). Per-report report lists: drop "and path" / "and name" from the Area,
  Factor, and Requirement entries and note the kind prefix (per the kind-prefix
  requirement above). Rendering Rules: reconcile the 0104 header-table rationale,
  whose cited `path/name line` is removed and whose "don't repeat the subject
  kind as metadata" now reads against a kind-prefixed title (per the kind-prefix
  and no-identifier-line requirements above).

### To rename

None

### To delete

None
