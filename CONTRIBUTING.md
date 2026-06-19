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

The project uses [mise](https://mise.jdx.dev) to pin tools (Go, dprint,
goreleaser) and define tasks. With mise installed:

```sh
mise install        # install pinned tool versions
```

If you prefer not to use mise, install Go 1.26+ and dprint yourself; the
`mise run …` commands below map to plain `go` and `dprint` invocations.

## Development tasks

```sh
mise run run -- init -   # run the current CLI scaffold command to stdout
mise run build           # build ./dist/qualitymd
mise run test            # go test ./...
mise run vet             # go vet ./...
mise run fmt             # gofmt -w . and dprint fmt
mise run fmt-md-check    # dprint check
mise run tidy            # go mod tidy
mise run hooks           # install repo-managed git hooks
```

Please run `mise run fmt`, `mise run test`, and `mise run vet` before opening a
pull request. For docs-only changes, `mise run fmt-md-check` is the quick
formatting gate.

### Git hooks

Run this once after cloning if you want Git to run the repo's pre-commit checks:

```sh
mise run hooks
```

The pre-commit hook checks staged Go files with `gofmt`, staged Markdown files
with `dprint`, verifies `go mod tidy` leaves `go.mod` and `go.sum` unchanged,
and runs `go test ./...`.

### Testing the CLI from another directory

Some commands are sensitive to the current working directory. For example,
`qualitymd init` writes `QUALITY.md` to the process's current directory when no
path is supplied.

When testing that behavior against the latest local source, build a temporary
binary from the repo, then run that binary from the directory you want to test:

```sh
go build -o /tmp/qualitymd-dev ./cmd/qualitymd

mkdir -p /tmp/qualitymd-init-test
cd /tmp/qualitymd-init-test
/tmp/qualitymd-dev init
```

Avoid `go run ./cmd/qualitymd init` from the repo root for this case: the command
would run with the repo as its working directory, so `init` would target the repo
root rather than your test directory.

## Project layout

```
cmd/qualitymd        entry point
internal/cli         Cobra commands, run through Charm Fang
internal/document    QUALITY.md frontmatter parsing, rendering, and file writes
internal/scaffold    embedded starter QUALITY.md used by qualitymd init
internal/model       typed QUALITY.md frontmatter model
scripts/build-npm.mjs   assembles the npm distribution
```

## Distribution and releases

Releases ship through three channels from a single git tag:

- **GitHub release archives** and a **Homebrew cask** — built by goreleaser
  (`.goreleaser.yaml`).
- **npm / npx** — `scripts/build-npm.mjs` cross-compiles a native binary per
  platform into a `@qualitymd/cli-<os>-<arch>` package gated by npm `os`/`cpu`
  fields, with the `quality.md` launcher selecting the right one at runtime
  (the esbuild/Biome model — no postinstall download).

Versioning policy for the separately distributed CLI, `/quality` skill, and
`SPECIFICATION.md` lives in [Versioning](docs/reference/versioning.md). Release
preparation, tag publishing, verification, failure handling, release secrets, and
changelog guidance live in [Cut a release](docs/guides/cut-a-release.md).

### Local dry runs

```sh
mise run snapshot                 # goreleaser build, no publish
mise run npm-build                # assemble npm packages under npm/platforms, no publish
```
