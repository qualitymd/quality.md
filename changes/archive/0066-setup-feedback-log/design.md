---
type: Design Doc
title: Setup feedback log - design doc
description: Design for a hand-authored, skill-only workflow feedback log under .quality/logs/, with setup as the first adopter.
tags: [skill, setup, logging, feedback]
timestamp: 2026-06-23T00:00:00Z
---

# Setup feedback log - design doc

Design behind the [Setup feedback log](../0066-setup-feedback-log.md) and its
[functional spec](spec.md).

## Context

The artifact is improvement feedback, not an audit trail. A maintainer reads it
to find out where a `/quality` workflow was slow, confusing, or wrong, so they
can improve the skill, CLI, and prompts. It is recorded locally and never
transmitted; the user may share it deliberately.

Three facts from the codebase shape the design:

- **Directory creation is already on demand.** `internal/evaluation/create.go`
  `MkdirAll`s `.quality/evaluations/...` at the moment it writes a run, and
  `qualitymd init` does not touch `.quality/` at all — it only writes
  `QUALITY.md`. The `internal/workspace` package *computes* paths; it never
  creates them. So nothing needs to "ensure" `.quality/logs/`; whatever writes a
  feedback file creates the directory the same way evaluation does.
- **The CLI carries no consent or location need here.** Location is already
  configurable (the root-model frontmatter `config` pointer and the
  `evaluationDir` config key), and the log is never transmitted, so there is no
  consent decision to gate. That removes the only two jobs a CLI flag could have
  done.
- **`.quality/log/` already exists** as the quality-log directory
  (`DefaultQualityLogDir`). The feedback directory is the plural `.quality/logs/`,
  chosen deliberately despite the one-character proximity.

## Approach

Skill-only. When `setup` reaches its close step, it hand-authors a feedback log
and writes it to `.quality/logs/<timestamp>-setup-feedback-log.md`, creating the
directory on demand. No CLI flag, no new command, no config field, no Go change.

- **Artifact shape** (generic, so evaluate/update can adopt it later):
  - Path: `.quality/logs/<timestamp>-<workflow>-feedback-log.md`; setup writes
    `<timestamp>-setup-feedback-log.md`. "log" is kept in the name so humans can
    refer to it plainly.
  - Environment header (frontmatter): workflow, UTC timestamp, agent/model, skill
    version, CLI version, platform, model-file-pre-existed, outcome (setup
    maturity), rough effort signal. The header is what makes an out-of-context
    note actionable; without "which agent/model, which version," a "this was
    slow" line is noise.
  - Body sections: friction/errors, UX/AX observations, efficiency/speed, what
    worked, suggested improvements, redaction note.

- **Timestamp.** The skill formats a sortable, UTC, filesystem-safe timestamp
  (for example `2026-06-23T154233Z`) so logs sort and read consistently — a SHOULD
  for consistency, not a rigid format. The spec defines no collision rule: at
  second granularity in a single interactive workflow a clash is vanishingly
  rare, so the agent appends a short disambiguator ad hoc only if one ever occurs.
  No CLI seeder is needed.

- **Redaction.** Secrets and raw prompt-injection text are never written
  (absolute). Sanitizing sensitive project context — proprietary source,
  identifying data, sensitive names/paths/domain specifics — is a SHOULD because
  nothing transmits the log automatically; it matters when the user shares it.

- **Sharing.** Neither skill nor CLI transmits the file. Sharing is an explicit
  user action (paste or attach). No phone-home.

## Alternatives

- **CLI flag on `init` (consent + dir scaffold).** Rejected. The codebase
  creates `.quality/` directories on demand, so a flag to create the dir solves a
  non-problem, and with no transmission there is no consent to record. `init`
  also runs only once on the create-new path, so it could not seed per-run files
  anyway.
- **General-purpose `--data-dir` flag.** Rejected for this case. It conflates
  *location* with the feedback feature; location is already configurable via the
  frontmatter `config` pointer and `evaluationDir`, and the directory is
  auto-created. A general data-dir override could be a separate ergonomic change,
  but it should not stand in for or ride on this work.
- **Small per-run CLI helper.** Rejected. The feedback log is hand-authored prose
  with no machine contract, no numbered sequence, and no snapshot to coordinate,
  so the determinism a CLI seeder buys is not worth a new command. If feedback
  logs later need a stable schema or a second adopter standardizes them, a
  `qualitymd feedback` helper can be added against observed data — the same staged
  path the `debug-log.md` design left open.
- **Per-run folder like evaluation.** Rejected. Setup is short, single-pass, and
  not interruption/resume-shaped, so it does not need a numbered run container. A
  flat `.quality/logs/` of timestamped files fits the longitudinal,
  improvement-feedback purpose.
- **Reuse evaluation's `debug-log.md`.** Rejected. That artifact is a per-run
  local audit inside an evaluation run folder with a process-vs-judgment boundary.
  Overloading it would blur a per-run audit with a central, cross-run feedback
  artifact.
- **`.quality/feedback/` directory.** Considered to avoid the `.quality/log/`
  proximity, but `.quality/logs/` was chosen so the directory reads as a general
  home for workflow logs. The collision risk is accepted and noted.
- **Top-level `specs/workflow-feedback-log.md`.** Considered, but the contract is
  placed as a **sub-spec of the setup workflow spec** (under
  `specs/skills/quality-skill/workflows/setup/`) instead, scoping it to its first
  and only adopter. It can be lifted to a shared spec when a second workflow
  adopts the artifact, rather than starting general before there is a second user.
- **Fold notes into the close report.** Rejected as the primary home: the
  "Setup complete / Not done / gaps" block is a terse user-facing status summary,
  not an environment-stamped, redacted improvement report. The feedback log is the
  persisted, expanded sibling of that idea, kept separate so the summary stays
  terse.

## Trade-offs & risks

- **Mutation-surface creep.** This is the first widening of setup's mutation
  boundary. Mitigation: narrow — only `.quality/logs/` — and every other setup
  prohibition stays in force.
- **Redaction.** A recorded artifact the user might share should not leak.
  Mitigation: absolute MUST-NOTs for secrets and injection text, plus a SHOULD to
  sanitize sensitive project context and note it.
- **Naming proximity.** `.quality/logs/` vs `.quality/log/`. Mitigation: the
  `*-feedback-log.md` filename, the explicit "distinct from the quality log" rule, and
  the "not an OKF `log.md`" rule.
- **Timestamp determinism.** Mitigation: sortable-and-unique is sufficient, with a
  disambiguator on collision; revisit a CLI seeder only if real logs collide.
- **Empty value.** If runs rarely produce notable feedback, the artifact adds
  little. Mitigation: setup MAY omit a log when nothing is notable.

## Open questions

None. The spec home (setup sub-spec), filename (`*-setup-feedback-log.md`), and
timestamp posture (sortable UTC SHOULD, no normative collision rule) are settled
above. The only implementation detail to confirm is the exact sub-spec filename
under `setup/`, which follows the artifact-spec naming convention.
