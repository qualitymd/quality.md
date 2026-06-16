# Contributing

Thanks for your interest in improving quality.md.

## Prerequisites

The project uses [mise](https://mise.jdx.dev) to pin tools (Go, goreleaser) and
define tasks. With mise installed:

```sh
mise install        # install pinned tool versions
```

If you prefer not to use mise, install Go 1.26+ yourself; the `mise run …`
commands below map to plain `go` invocations.

## Development tasks

```sh
mise run run -- check   # run the placeholder CLI (go run ./cmd/qualitymd)
mise run build          # build ./dist/qualitymd
mise run test           # go test ./...
mise run vet            # go vet ./...
mise run fmt            # gofmt -w .
mise run tidy           # go mod tidy
```

Please run `mise run fmt` and `mise run vet` before opening a pull request.

## Project layout

```
cmd/qualitymd        entry point
internal/cli         Cobra commands, run through Charm Fang
internal/spec        current QUALITY.md frontmatter model loader
internal/eval        placeholder assessment traversal
internal/report      Lip Gloss terminal output
scripts/build-npm.mjs   assembles the npm distribution
```

## Distribution channels

Releases ship through three channels from a single git tag:

- **GitHub release archives** and a **Homebrew cask** — built by goreleaser
  (`.goreleaser.yaml`).
- **npm / npx** — `scripts/build-npm.mjs` cross-compiles a native binary per
  platform into a `@qualitymd/cli-<os>-<arch>` package gated by npm `os`/`cpu`
  fields, with the `quality.md` launcher selecting the right one at runtime
  (the esbuild/Biome model — no postinstall download).

### Local dry runs

```sh
mise run snapshot                 # goreleaser build, no publish
mise run npm-build                # assemble npm packages under npm/platforms, no publish
```

### Cutting a release

Releases are automated by `.github/workflows/release.yml` on any `v*` tag:

```sh
git tag v1.2.3
git push origin v1.2.3
```

The workflow runs goreleaser and publishes the npm packages. It requires two
repository secrets:

- `NPM_TOKEN` — npm automation token with publish access to the `quality.md`
  package and the `@qualitymd` scope.
- `HOMEBREW_TAP_GITHUB_TOKEN` — token with write access to
  `qualitymd/homebrew-tap`.

> The Homebrew cask strips the macOS quarantine attribute because the binaries
> are currently unsigned. Remove that step once the binaries are signed and
> notarized.
