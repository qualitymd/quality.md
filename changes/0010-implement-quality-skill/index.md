# Implement the /quality skill

Children of the
[Implement the /quality skill](../0010-implement-quality-skill.md) change.

# Concepts

- [Functional spec](spec.md) - what the change must do; defers the behavioral
  contract to the
  [`/quality` skill spec](../../specs/skills/quality-skill/quality-skill.md) and
  records the open items and gaps to settle.
- [Design doc](design.md) - how the change is built: the skill packaged for Agent
  Skills installation, the CLI prerequisite check in `setup`/`wizard`, the
  `qualitymd models` CLI surface, the raw JSON evaluation artifacts, and how the
  open items resolve into the durable spec.
