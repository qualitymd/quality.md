# /quality Skill Modes Update Log

## 2026-06-22

- **Mode removal**: Removed the `improve` mode spec from this folder. Applying or
  handing off evaluation recommendations is now governed by the non-mode
  [`recommendation follow-up`](../recommendation-follow-up.md) spec.

- **Creation**: Originally added behavioral component specs for the `/quality` runtime modes:
  [`setup`](setup.md), [`wizard`](wizard.md), [`evaluate`](evaluate.md),
  `improve`, and [`update`](update.md). The parent
  [`/quality` skill](../quality-skill.md) spec keeps shared contracts and links
  to these mode-specific contracts.
