# /quality Skill Workflows

Behavioral component specs for the `/quality` runtime workflows. The parent
[/quality skill](../quality-skill.md) spec owns shared contracts; these specs own
workflow-specific routing, mutation surfaces, required artifacts, stop
conditions, and completion criteria. Each workflow is dispatched as a mode.

# Workflows

- [setup](setup.md) - bootstrap or populate a `QUALITY.md` model.
- [evaluate](evaluate.md) - create evaluation records and reports.
- [update](update.md) - orchestrate compatible `/quality` skill and CLI updates.

# Bundle

- [log.md](log.md) - changes to this workflow-spec folder.
