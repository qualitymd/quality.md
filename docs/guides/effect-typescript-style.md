---
type: How-to Guide
title: Write Effect TypeScript
description: Keep qualitymd TypeScript pure at the center, explicit at runtime boundaries, deterministic under tests, and easy for agents to validate.
tags: [typescript, effect, testing, agents]
timestamp: 2026-07-14T00:00:00Z
---

# Write Effect TypeScript

Use this guide when adding or reshaping CLI runtime code, tests, generators, or
release TypeScript.

## Put code in the owning layer

- `src/domain/` owns deterministic values, parsing, validation, planning, and
  rendering. Pass time and generated IDs in as values. Do not read the host.
- `src/application/` composes use cases from domain functions and Effect
  services. It owns workflow order, not Bun, terminal, network, or provider
  details.
- `src/services/` defines repository-owned capabilities and live/test Layers.
- `src/adapters/` translates provider SDKs and direct HTTP/process behavior into
  service contracts.
- `src/cli/` owns command grammar and the final mapping to stdout, stderr, and
  exit codes.

The architecture contract test rejects direct host APIs in domain/application
code and inward-pointing imports from domain or services. Extend that focused
test when an observed boundary violation can be detected reliably.

## Keep pure work pure

Use ordinary functions for deterministic transformations. Return an Effect when
the operation needs a service, a typed failure channel, scope, interruption,
concurrency, retry, or time. Do not wrap a pure calculation in Effect merely for
uniformity.

Use namespace imports for Effect modules, for example `import * as Effect from
"effect/Effect"`. Keep exported types explicit. Preserve exact dependency pins
while Effect v4 and its CLI/testing packages are pre-stable.

## Model runtime capabilities as services

Define the smallest repository-owned capability that makes a workflow
deterministic. Provide live Bun/SDK behavior at `src/main.ts` or through an
adapter; provide a small test Layer in service tests. Prefer Effect's standard
filesystem, path, clock, random, and scope facilities before inventing another
abstraction.

Expected operational failures belong in the typed error channel. A `throw`
inside `Effect.gen` is a defect, not an expected failure; use `Effect.fail`,
`Effect.try`, or `Effect.tryPromise`. Translate foreign SDK/process errors once
at the adapter boundary. Keep CLI exit-code decisions out of lower layers.

Acquire child processes, provider sessions, streams, and temporary resources in
scope. Register cleanup at acquisition time. Test interruption where losing
accepted work or leaving a child behind would be costly.

## Place and write tests

Keep the package-local sibling tree:

```text
src/domain/...       test/domain/...
src/application/...  test/application/...
src/services/...     test/services/...
src/adapters/...     test/adapters/...
                     test/integration/...
```

Use ordinary Vitest for pure functions. Use `@effect/vitest` and `it.effect` for
tests that return an Effect; return the Effect instead of calling
`Effect.runPromise`. Supply deterministic Layers for host behavior. Put real
filesystem, subprocess, compiled-executable, provider-runtime, and
cross-boundary checks under `test/integration/` or the owning hosted workflow.

Choose the lowest useful test boundary. Protect public command grammar,
stdout/stderr, exit categories, JSON, and generated files with executable
contract tests. Protect retry, cancellation, accepted-result durability, and
resume at the workflow boundary. A migration ledger may classify retired tests,
but active contracts need active tests.

## Regenerate from the owner

Do not hand-edit a generated artifact. Use its owner, then run its drift check.

| Generated artifact                | Owner                     | Drift check                     |
| --------------------------------- | ------------------------- | ------------------------------- |
| `quality.schema.json`             | `mise run schema`         | `mise run schema-check`         |
| `mintlify/cli.mdx`                | `mise run cli-docs`       | `mise run cli-docs-check`       |
| `mintlify/specification.mdx`      | `mise run sync-spec-docs` | `mise run sync-spec-docs-check` |
| `examples/report-gallery/`        | `mise run report-gallery` | `mise run report-gallery-check` |
| npm package README/platform trees | `mise run npm-build`      | `mise run npm-pack-check`       |

## Validate at the right phase

During an edit, run the narrowest relevant command:

```sh
mise run test -- test/domain/model/model.test.ts
mise run typecheck
mise run schema-check
mise run cli-docs-check
```

Before handoff or commit, run `mise run fmt`, then the canonical full gate:

```sh
mise run check
```

For distribution or release changes, also run `mise run snapshot`, `mise run
npm-build -- <version>`, and the release guide's exact pre-tag checks. Native
execution and published-channel verification belong to their hosted workflows;
do not infer them from a source-mode test.
