# /quality skill

The companion evaluation skill for QUALITY.md: it carries the evaluative
judgment and drives the deterministic [CLI](../../cli.md) for every mechanical
step. This folder holds the skill's functional spec and the reference artifacts
that make it concrete. The installable skill artifact lives at
[`../../../skills/quality/SKILL.md`](../../../skills/quality/SKILL.md).

# Specs

- [/quality skill](quality-skill.md) - parent functional spec for shared
  contracts: operating model, invocation, CLI ownership, evaluation semantics,
  reporting, quality log, and cross-workflow invariants.
- [/quality evaluation workflow](evaluation.md) - component spec for evaluation
  conformance, workflow, grounding, QC, and rating judgment.
- [/quality reporting](reporting.md) - component spec for evaluation run
  artifacts, report outputs, records, recommendations, and correction behavior.
- [/quality quality log](quality-log.md) - component spec for dated
  `.quality/log/` model-change entries.
- [Workflow feedback log](workflow-feedback-log.md) - shared artifact contract
  for workflow feedback logs under `.quality/logs/`.
- [/quality recommendation follow-up](recommendation-follow-up.md) -
  post-evaluation follow-up workflow for applying or handing off evaluation
  recommendations.
- [workflows/](workflows/index.md) - behavioral component specs for the runtime
  workflows: setup, evaluate, and update.
- [guides/](guides/index.md) - specs for runtime guide resources bundled with
  the skill, including the routed authoring guide family, getting-started,
  recommendation follow-up, and checklist guides.
- [examples/](examples/index.md) - worked runtime reference artifacts produced by
  the skill.
- [Installable skill artifact](../../../skills/quality/SKILL.md) - the prompt
  distributed by Agent Skills tooling.

# Subfolders

- [examples/](examples/index.md) - worked reference artifacts produced by the
  skill (e.g. an example Evaluation Report), used to make the
  [Reporting](reporting.md#reporting) contract concrete.
- [guides/](guides/index.md) - contracts for runtime guide resources bundled
  with the skill.
- [workflows/](workflows/index.md) - behavioral component specs for the runtime
  workflows.
