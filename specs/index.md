---
okf_version: "0.1"
---

# Specifications

Working specifications for the `qualitymd` tooling, as an OKF knowledge bundle.
The QUALITY.md format itself is specified in the repository-root
[`SPECIFICATION.md`](../SPECIFICATION.md), the source of truth for the file
format and evaluation semantics.

This bundle holds the _cumulative_ source of truth for the tool's current
behavior. The [`changes/`](../changes/index.md) bundle proposes _deltas_ to it:
each change lists the specs here it creates or updates, and brings them into sync
before it lands.

# Specs

- [glossary.md](glossary-md.md) - artifact contract for the workspace-root glossary used by generated reports and readers.
- [qualitymd CLI](cli.md) - high-level requirements for the deterministic command-line surface.
- [Evaluation](evaluation/evaluation.md) - replacement evaluation
  workflow, protocol, structured data, orchestration, and reports.
- [quality.schema.json](quality-schema-json.md) - companion structural JSON Schema
  for QUALITY.md frontmatter, derived from the linter's schema.

# Bundle

- [schema.md](schema.md) - the concept types used in this bundle.

# Subfolders

- [cli/](cli/) - per-command sub-specs.
- [evaluation/](evaluation/) - replacement evaluation workflow, routine,
  record, orchestration, and report contracts.
- [skills/](skills/) - companion skills that carry judgment around QUALITY.md (`/quality`).
