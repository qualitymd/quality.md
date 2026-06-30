---
type: How-to Guide
title: Go style
description: Judgment-based Go conventions that the linters do not enforce.
tags: [go, code, contributing]
timestamp: 2026-06-20T00:00:00Z
---

# Go style

The `qualitymd` CLI is Go under `cmd/qualitymd/` and `internal/`. This guide
covers the **judgment calls a reviewer makes** — the conventions that shape
readable Go but that no tool can decide for you.

It deliberately does **not** repeat anything the automated gate already
enforces. `mise run check` (run in CI and the git hooks) covers formatting and
imports (`gofmt`, `goimports`), suspicious constructs (`go vet`), unchecked
errors (`errcheck`), dead code and unused parameters (`unused`, `unparam`),
ineffective assignments and needless conversions (`ineffassign`, `unconvert`),
correctness, simplification, and some style conventions (`staticcheck`),
misspellings (`misspell`), and size/complexity limits (`cyclop`, `gocognit`,
`nestif`, `funlen`). If a rule below would duplicate one of those, it has been
left out on purpose — let the tool be the source of truth.

`staticcheck` here is the golangci-lint v2 linter, which folds in the former
`stylecheck` (ST) checks — so a few _style_ conventions are enforced too, not
just correctness: notably `self`/`this` receiver names (ST1006) and
error-string punctuation (ST1005). Don't restate those below.

For _which package a type belongs in_, see
[Designing Go packages](design-go-packages.md). This guide is about everything
else.

## Naming

- **Don't stutter.** The package qualifies the name, so the name should not
  repeat it: `evaluation.Load`, not `evaluation.LoadEvaluation`;
  `lint.Result`, not `lint.LintResult`.
- **No `Get` prefix on accessors.** A getter is named for the thing it returns:
  `cfg.Timeout()`, not `cfg.GetTimeout()`. Reserve verb prefixes for methods
  that _do_ something.
- **Package names are short, lowercase, singular, and meaningful.** Name the
  package for the concept it provides; never a catch-all (`util`, `common`,
  `helpers`, `misc`). [Designing Go packages](design-go-packages.md) covers why
  a catch-all bucket is the wrong home for a type.
- **Receiver names are short and consistent.** One or two letters, the same
  letter for every method on a type (`e *Evaluation`). The `self`/`this` ban is
  enforced by the gate (see above), so it lives there, not here.
- **Error values and types follow the convention.** Sentinel values are
  `ErrNotFound`; custom error types are `ParseError`.

## Errors

- **Add context as you return.** Wrap with `fmt.Errorf("loading %s: %w", path,
err)`. Skip noise like `"failed to"` / `"error while"` — the fact that it is
  an error is already clear, and the chain reads as a path: `loading config:
parsing yaml: unexpected EOF`.
- **`%w` only when the caller should inspect the cause.** Use `%w` to preserve
  a cause that callers match with `errors.Is` / `errors.As`. Use `%v` to flatten
  an error at a trust boundary where the underlying type is an implementation
  detail you don't want to leak.
- **Handle an error once.** Either log it or return it — not both. Logging and
  returning produces the same failure reported at every layer. The code that
  decides what to do owns the logging.
- **Compare with `errors.Is` / `errors.As`, never string matching.** Define a
  sentinel or a typed error so callers have something stable to check.
- **Libraries return errors; they don't `panic` or exit.** Reserve `panic` for
  genuine programmer error (unreachable invariants). Only `main` and top-level
  command wiring call `os.Exit` / `log.Fatal`.

## Interfaces and API shape

- **Accept interfaces, return concrete types.** A function takes the narrow
  interface it needs and hands back the concrete value it built, so callers keep
  the full type.
- **Define interfaces where they're consumed,** not next to the implementation.
  The consumer knows the minimal method set it actually requires; keep it small.
- **Assert compliance at compile time** when a type must satisfy an interface:
  `var _ Loader = (*FileLoader)(nil)`. It turns an accidental break into a
  build error.
- **Prefer composition over embedding in exported structs.** Embedding promotes
  the inner type's methods into your public surface, leaking detail you may not
  want to commit to.
- **Many parameters → a config struct or functional options.** Once a
  constructor grows past a few arguments, an options struct (named fields,
  obvious defaults, no positional swaps) or variadic `Option` functions beats a
  long positional list.
- **Context is the first parameter and is never stored in a struct.** Thread
  `ctx context.Context` through call chains. In Cobra commands use
  `cmd.Context()`, not `context.Background()`.

## Concurrency

- **Every goroutine has a defined end.** Before you spawn one, know how it stops
  and how the caller waits for it (a `WaitGroup`, a done channel, context
  cancellation). No fire-and-forget.
- **Don't start goroutines in `init()`.** Background work belongs behind an
  explicit constructor that returns something closeable, so lifetime is the
  caller's to manage.

## Values and state

- **Treat a `nil` slice or map as empty.** Return `nil` rather than `[]T{}` and
  test emptiness with `len(s) == 0`. Don't distinguish nil from empty in an API.
- **Copy slices and maps at boundaries** when you store a caller-supplied one
  (or return an internal one) and the two sides must not alias. Sharing the
  backing array invites spooky action at a distance.
- **Use the comma-ok form for type assertions:** `v, ok := x.(T)`. The
  single-return form panics on a miss.
- **No mutable package-level globals.** Inject dependencies through
  constructors. Keep `init()` deterministic — no I/O, no environment reads.

## Documentation comments

- **Document every exported identifier, starting with its name:** `// Load reads
…`, `// Evaluation represents …`. This is the form `go doc` expects.
- **Document the contract, not the obvious.** Call out what a reader can't infer
  from the signature: nil-on-error behavior, which sentinel errors come back,
  whether a type is safe for concurrent use, and any cleanup the caller owes
  (`Close` / `Stop`).

## Tests

- **Table-driven, with named fields.** Use a slice of structs with named fields
  (not positional literals) so a reader can tell which value is which, and run
  cases as subtests with `t.Run`.
- **Mark helpers with `t.Helper()`** so failures report the caller's line, not
  the helper's.
- **Never call `t.Fatal` from a goroutine.** It must run on the test's own
  goroutine; signal the failure back and fail there.

## Sources

These conventions are drawn from the Go community's standard references, narrowed
to what the linters don't already enforce:

- [Effective Go](https://go.dev/doc/effective_go)
- [Go Code Review Comments](https://go.dev/wiki/CodeReviewComments)
- [Google Go Style Guide](https://google.github.io/styleguide/go/)
- [Uber Go Style Guide](https://github.com/uber-go/guide/blob/master/style.md)
