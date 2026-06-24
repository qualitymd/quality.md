---
type: Change Case
title: View command
description: Add a qualitymd view presentation command — a layered human/agent rendering surface over the QUALITY.md workspace, starting with a pretty document render and a model outline.
status: Draft
tags: [cli, view, presentation, skill, setup]
timestamp: 2026-06-24T00:00:00Z
---

# View command

A **Change Case** to add `qualitymd view`, a read-only presentation surface that
renders a `QUALITY.md` workspace for humans and agents. `view` is structured along
two orthogonal axes: a **lens** (what you are looking at — the document, the model
outline, and later ratings/coverage/trends/recommendations) and a **surface** (how
it is rendered — terminal text, JSON, a diagram source, and later static HTML or an
interactive server). Every renderer consumes the same deterministic workspace data,
so richer surfaces can be added later without re-architecting.

This case lands the first slice: the default document render and the `outline`
lens, on the `text`, `json`, and `mermaid` surfaces. The remaining lenses and
surfaces are designed here but deferred.

Detail lives in:

- [Functional spec](0078-view-command/spec.md) - what this case must build.
- Design doc - the holistic two-axis design, the view/evaluation boundary, the
  workspace-graph data contract, and the interactive future (authored in the
  Design phase).

## Motivation

`/quality setup` closes by reporting status and next steps, but the user never
sees the model their discovery answers produced — it is authored body-first and
frontmatter-second and then summarized only in prose. A rendered structure view at
setup close turns abstract discovery into a concrete "here is what we built," which
fits setup's teaching-first philosophy and doubles as a soft correction gate.

More broadly, a `QUALITY.md` workspace is multi-dimensional (areas × factors ×
ratings × time) and the only ways to see it today are reading raw YAML/Markdown or
scraping `status`. `status` is deliberately a compact routing snapshot, not a
presentation surface. A dedicated `view` command gives the project one home for
rendering the model and (later) its evaluation state — serving the format's
**alignment** purpose for humans and stakeholders outside the agent loop, and
giving the agent a deterministic structure/diagram surface to embed in reports.

## Scope

Covered:

- A new read-only `qualitymd view [path]` command.
- Default lens (no subcommand): pretty-render the model's Markdown body with the
  shared glamour renderer (the one `qualitymd spec` uses), plus a readable model
  header — never a raw YAML dump.
- `view outline` lens: the model hierarchy (root area → areas → factors) with
  requirement counts and a rating-scale legend, at a counts-plus-top-levels depth.
- Surfaces via `--format`: `text` (default), `json`, `mermaid`; `--json` as an
  alias for `--format json`.
- Reuse of the deterministic model-shape computation `status` already performs.
- Wire `/quality setup` closeout to render the outline via `view`.

Deferred / non-goals:

- Additional lenses: `health` (ratings overlay), `coverage` (source mapping),
  `trends` (ratings over time), `recommendations`.
- Additional surfaces: `dot`, static `html`, and an interactive `--serve` server.
- A standalone durable spec for the workspace-graph data contract (suggested
  below; the contract lives in the `view` spec for now).
- No model-quality judgment, rating recomputation, or report-body scraping — that
  stays with `evaluation`.
- No change to the QUALITY.md format or schema.

## Affected artifacts

Derived by sweeping the repo for the command surface and the setup-closeout flow;
empty kinds are deliberate.

### Code

- [ ] `internal/cli/` - add the `view` command (new file) and register it in
      `root.go`; render the default lens with the shared glamour/style path used by
      `spec.go`.
- [ ] `internal/status/` (or shared model-shape helper) - reuse the existing
      deterministic Area/Factor/Requirement shape computation for `outline`.
- [ ] `internal/cli/` tests - golden/behavior tests for `view`, `view outline`,
      and the `text`/`json`/`mermaid` surfaces.

### Durable specs

- [ ] `specs/cli/view.md` - **new** durable command spec (To add).
- [ ] `specs/cli/index.md` - register the `view` spec.
- [ ] `specs/cli.md` - list `view` in the command set if it enumerates commands.
- [ ] `specs/cli/status.md` - optional cross-reference clarifying the
      status (routing) vs. view (presentation) boundary.
- [ ] `specs/skills/quality-skill/workflows/setup.md` - specify that setup
      closeout renders the model outline via `view`.

Suggested new durable spec (not a precondition for landing): a workspace-graph
data-contract spec once a second non-text surface consumes the JSON.

### Format spec

- [ ] `SPECIFICATION.md` - no change; presentation is not a format concern.

### Durable docs (README and bundled skill)

- [ ] `skills/quality/workflows/setup.md` - render the outline at setup close.
- [ ] `README.md` - mention `qualitymd view` in the CLI surface if the command
      list warrants it.

### Release

- [ ] `CHANGELOG.md` - add an Unreleased CLI note for `view`.

## Children

- [Functional spec](0078-view-command/spec.md) - required `view` command behavior
  for this case's slice.

## Status

`Draft`. Functional spec authored; design doc and implementation not started.
Code is gated to In-Progress.
