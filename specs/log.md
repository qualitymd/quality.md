# Specs Update Log

## 2026-06-17

- **Revision**: Added a "Technical requirements" section to the
  [CLI spec](cli.md) requiring that every functional requirement be satisfiable
  through the idiomatic capabilities of the chosen stack (Go + Cobra + Charm
  Fang + Lip Gloss), rather than working against the grain of those libraries.
- **Revision**: Added a "Conventions" section to the [CLI spec](cli.md)
  establishing `--json` as the spelling for machine-readable output wherever a
  command offers one (not a requirement that every command do so), with criteria
  for when a command should offer `--json` and worked examples across the current
  commands. Updated [`lint`](cli/lint.md) to reference `--json` in place of its
  earlier `--format json`.
- **Convention**: Added a "Suggested next actions" convention to the
  [CLI spec](cli.md): commands may close with advisory, deterministic next-action
  suggestions — a stderr footer in human output, an in-band `nextActions` array
  under `--json` — that never affect behavior or the exit code.

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
