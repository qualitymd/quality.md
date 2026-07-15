# Contributing

Thanks for your interest in improving QUALITY.md.

## Where process lives

- **Contributor setup, local tasks, pull request checks, and repo layout** live
  in this file.
- **Release operations** live in
  [Cut a release](docs/guides/cut-a-release.md).
- **Versioning and compatibility policy** lives in
  [Versioning](docs/reference/versioning.md).
- **Change Case workflow** lives in
  [Working with change cases](docs/guides/work-with-change-cases.md).
- **User install and bootstrap** lives in [Install QUALITY.md](install.md).

## Prerequisites

The project uses [mise](https://mise.jdx.dev) to pin Bun, Node, and Prettier and
define tasks. TypeScript, Effect, provider SDKs, the formatter, linter, and test
runner are pinned in `package.json` and `bun.lock`. With mise installed:

```sh
mise install        # install pinned tool versions
```

If you prefer not to use mise, install Bun 1.3.14, Node 22+, and Prettier 3.9.5
yourself. Run `bun install --frozen-lockfile` before the tasks below.

## Development tasks

```sh
mise run run -- init -   # run the current CLI scaffold command to stdout
mise run build           # build ./dist/qualitymd
mise run typecheck       # strict TypeScript check
mise run test            # Vitest suite
mise run lint            # oxlint with warnings denied
mise run fmt             # oxfmt and prettier writes
mise run check           # run the same gate as CI and git hooks
mise run schema          # regenerate quality.schema.json from TypeScript
mise run fmt-md-check    # prettier --check
mise run npm-pack-check  # verify npm package README packaging
mise run report-gallery  # regenerate checked-in example evaluation reports
mise run report-gallery-check  # verify generated report gallery is current
mise run cli-docs        # regenerate the generated mintlify/cli.mdx from the CLI
mise run sync-spec-docs  # regenerate the generated mintlify/specification.mdx
mise run docs-links-check  # check the Mintlify docs for broken links
mise run docs-deps       # install the pinned Mintlify CLI under mintlify/
mise run release-notes -- v0.3.0  # print curated release notes
mise run release-check -- v0.3.0  # run pre-tag release checks
mise run hooks           # install repo-managed git hooks
```

Please run `mise run fmt` before committing formatting-sensitive changes, then
`mise run check` before opening a pull request. For docs-only changes,
`mise run fmt-md-check` is the quick formatting gate.

### Test layout and targeted checks

Tests live in the package-local sibling `test/` tree and mirror `src/` by
boundary. Pure domain tests use ordinary Vitest. Effect workflows use
`@effect/vitest`; executable, real-filesystem, and cross-boundary contracts live
under `test/integration/`. See
[Write Effect TypeScript](docs/guides/effect-typescript-style.md).

Run one file while iterating:

```sh
mise run test -- test/domain/model/model.test.ts
mise run test -- test/services/source.test.ts
mise run test -- test/integration/cli.test.ts
```

Generated outputs are edited through their owning tasks: `schema`, `cli-docs`,
`sync-spec-docs`, and `report-gallery`. Their `*-check` tasks verify drift.
`mise run check` remains the one complete local/CI gate.

### Git hooks

Run this once after cloning if you want Git to run the repo's hooks:

```sh
mise run hooks
```

The repo-managed hooks run the same `mise run check` gate as CI. `pre-commit`
temporarily stashes unstaged files, regenerates and stages
`examples/report-gallery`, and checks the
staged snapshot. `pre-push`
temporarily stashes local changes and checks the committed snapshot.

### Testing the CLI from another directory

Some commands are sensitive to the current working directory. For example,
`qualitymd init` writes `QUALITY.md` to the process's current directory when no
path is supplied.

When testing that behavior against the latest local source, build a temporary
binary from the repo, then run that binary from the directory you want to test:

```sh
mise run build

mkdir -p /tmp/qualitymd-init-test
cd /tmp/qualitymd-init-test
<repo>/dist/qualitymd init
```

Running `mise run run -- init` from the repo root keeps the repo as its working
directory, so use the standalone executable when testing another directory.

## Project layout

```
src/main.ts            standalone executable entry point
src/cli                Effect CLI adapter and command tree
src/application        command use cases and evaluation runner
src/domain             pure model, evaluation, validation, and report logic
src/services           filesystem, evaluator, workspace, and output boundaries
src/adapters           provider SDK, process, and direct HTTP implementations
test/domain            pure deterministic unit tests
test/application       Effect workflow and selection tests
test/services          deterministic service/Layer tests
test/integration       filesystem, executable, architecture, and provider contracts
scripts/build-npm.mjs   assembles the npm distribution
scripts/extract-release-notes.mjs   extracts a tagged CHANGELOG.md section
scripts/check-release.mjs   runs strict pre-tag release checks
```

## Distribution and releases

Releases ship through three channels from a single git tag:

- **GitHub release archives** and a **Homebrew cask** — Bun compiles the
  standalone platform matrix through `scripts/build-release.ts`; the release
  workflow publishes checksum-verified archives and updates the tap cask.
- **npm / npx** — `scripts/build-npm.mjs` cross-compiles a native binary per
  platform into a `@qualitymd/cli-<os>-<arch>` package gated by npm `os`/`cpu`
  fields, with the `quality.md` launcher selecting the right one at runtime
  (the esbuild/Biome model — no postinstall download).
  `npm/quality.md/README.md` is generated from the repository `README.md`
  during package build/prepack so npm renders the same README as GitHub.
  Generated platform package READMEs explain that direct installs should use the
  `quality.md` package instead.

Versioning policy for the separately distributed CLI, `/quality` skill, and
`SPECIFICATION.md` lives in [Versioning](docs/reference/versioning.md). Release
preparation, tag publishing, verification, failure handling, release secrets, and
changelog guidance live in [Cut a release](docs/guides/cut-a-release.md).

### Local dry runs

```sh
mise run snapshot                 # build every standalone target, no publish
mise run npm-build                # assemble npm packages under npm/platforms, no publish
mise run release-notes -- v0.3.0  # preview GitHub Release notes from CHANGELOG.md
mise run release-check -- v0.3.0  # strict pre-tag release gate
```
