# Specs Update Log

## 2026-06-16

- **Convention**: Added a bundle-root [`schema.md`](schema.md) (`type: Schema`)
  registering the bundle's concept types (`Schema`, `Functional Specification`)
  in frontmatter, and listed it from the [index](index.md). Retyped the CLI
  spec and command sub-specs from `Specification` / `Command Specification` to a
  single `Functional Specification` type.

- **Initialization**: Created the `specs/` OKF bundle with a root
  [CLI spec](cli.md) capturing the high-level CLI requirements (design
  properties, global flags, output formats, exit codes, agent accessibility).
- **Creation**: Added placeholder command sub-specs for
  [`init`](cli/init.md), [`lint`](cli/lint.md), and [`spec`](cli/spec.md), plus
  the [`cli/` index](cli/index.md).
- **Revision**: Reduced the [CLI spec](cli.md) to a stub. No requirements had
  actually been requested, so the speculative design properties, global flags,
  output formats, exit codes, and agent-accessibility requirements were stripped
  out, leaving scope, a command list, and a "to be specified" outline — matching
  the placeholder command sub-specs.
