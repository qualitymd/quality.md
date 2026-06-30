# /quality skill workflows

Behavioral component specs for the `/quality` runtime workflows. The parent
[/quality skill](../quality-skill.md) spec owns shared contracts; these specs own
workflow-specific routing, mutation surfaces, required artifacts, stop
conditions, and completion criteria. Each workflow is dispatched from the root
prompt.

# Workflows

- [setup](setup.md) - bootstrap or populate a `QUALITY.md` model.
- [evaluate](evaluate.md) - create evaluation records and reports.
- [review](review.md) - review an Evaluation result, the model, or a quality concern.
- [improve](improve.md) - improve the work or model after confirming focus and mutation surface.
- [update](update.md) - orchestrate compatible `/quality` skill and CLI updates.

# Subfolders

- [evaluate/](evaluate/index.md) - sub-specs owned by the evaluate workflow,
  including the [evaluate feedback log](evaluate/feedback-log.md).
- [setup/](setup/index.md) - sub-specs owned by the setup workflow, including the
  [setup feedback log](setup/feedback-log.md).

# Bundle

- [log.md](log.md) - changes to this workflow-spec folder.
