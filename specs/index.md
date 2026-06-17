---
okf_version: "0.1"
---

# Specifications

Working specifications for the `qualitymd` tooling, as an OKF knowledge bundle.
The `QUALITY.md` format itself is specified in the repository-root
[`SPECIFICATION.md`](../SPECIFICATION.md), the source of truth for the file
format and evaluation semantics.

This bundle holds the *cumulative* source of truth for the tool's current
behavior. The [`changes/`](../changes/index.md) bundle proposes *deltas* to it:
each change lists the specs here it creates or updates, and brings them into sync
before it lands.

# Specs

- [qualitymd CLI](cli.md) - high-level requirements for the deterministic command-line surface.

# Bundle

- [schema.md](schema.md) - the concept types used in this bundle.

# Subfolders

- [cli/](cli/) - per-command sub-specs (`init`, `lint`, `spec`).
