---
type: Design Doc
title: agent accessibility — design doc
description: How categorized exit codes and the init --json receipt are threaded through the Cobra / Fang stack.
tags: [cli, specification, agent-accessibility, design]
timestamp: 2026-06-17T00:00:00Z
---

# agent accessibility — design doc

Design behind the [Specify and enforce agent accessibility](../0004-specify-agent-accessibility.md)
change and its [functional spec](spec.md). The spec fixes *what* the contract
requires; this doc covers *how* the code makes it so. The specification half of
the change (the durable [CLI spec](../../../specs/cli.md) edits) is prose, not
code, so this doc concerns the **conformance** half: categorized exit codes and
the `init --json` receipt.

## Context

Two baseline requirements need code that the current `internal/cli` does not
have:

- **Categorized exit codes.** [`cli.Execute`](../../../internal/cli/root.go)
  today does `os.Exit(1)` on any non-nil error from `fang.Execute`, so a `lint`
  that *found problems* is indistinguishable from a usage or internal failure.
- **`init --json`.** [`init`](../../../internal/cli/init.go) has no `--json`; the
  revised convention makes it owe a result receipt.

The enabling fact, from reading `fang@v1.0.0`: **`fang.Execute` returns the
command error and never calls `os.Exit` itself** — our wrapper does. Fang also
renders the error to stderr through a replaceable handler (`WithErrorHandler`,
delegating to the exported `DefaultErrorHandler`) and already prints plainly when
stderr is not a TTY. So the exit-code contract is entirely ours to impose at the
`Execute` boundary; we are not fighting the framework.

This shapes the whole design: returning the error rather than exiting *is* Fang's
intended model for letting the caller assign exit codes, and `WithErrorHandler`
is its sanctioned rendering hook. The conformance work should ride those two
extension points plus Cobra-native hooks (`FlagErrorFunc`, `Args`), not bolt a
parallel mechanism beside them. Fang exposes no exit-code option of its own
(confirmed against its `Option` set), so a thin boundary mapping is the idiomatic
path, not a workaround. The bar for every piece below is: does it use a hook the
framework already offers, or is it scaffolding we invented?

## Approach

### Exit codes: one typed error, mapped at the boundary

Define the categories and a coded error in `internal/cli`:

```go
const (
	ExitOK       = 0  // success
	ExitProblems = 1  // ran, produced a reportable negative (lint findings)
	ExitUsage    = 2  // malformed invocation
	ExitInternal = 70 // command could not complete the requested action
)

// CodedError carries an exit category out of a command. Silent suppresses
// Fang's error rendering for outcomes the command has already reported.
type CodedError struct {
	Code   int
	Silent bool
	Err    error
}

func (e *CodedError) Error() string { return e.Err.Error() }
func (e *CodedError) Unwrap() error { return e.Err }
```

`cli.Execute` maps the returned error to a code once, at the top:

```go
func Execute() {
	err := fang.Execute(ctx, newRootCmd(),
		fang.WithVersion(version), fang.WithCommit(commit),
		fang.WithErrorHandler(errorHandler),
	)
	os.Exit(codeFor(err))
}

func codeFor(err error) int {
	if err == nil {
		return ExitOK
	}
	var ce *CodedError
	if errors.As(err, &ce) {
		return ce.Code
	}
	if isUsageError(err) {
		return ExitUsage
	}
	return ExitInternal
}
```

Only commands that need a *non-default* category construct a `CodedError`.
Everything that returns a plain error falls through to `ExitInternal` — the right
default for an I/O failure or a bug, and what the existing I/O paths
(`scaffold.Create`, `lint.Check` on a missing file) already produce, so they need
no change.

### Tagging usage errors at their source, not by sniffing

Fang's `isUsageError` is an unexported string-prefix hack we cannot call and
should not duplicate as our primary path. Instead tag usage errors where Cobra
raises them, so they arrive already typed:

- **Flag-parse errors** — set a `FlagErrorFunc` on the root (inherited by
  subcommands) that wraps the error as `&CodedError{Code: ExitUsage}`. This
  catches unknown/ambiguous flags and missing flag arguments.
- **Argument errors** — wrap the `Args` validator so a positional-arg violation
  is tagged too: `Args: usage(cobra.MaximumNArgs(1))`, where `usage` maps a
  non-nil result to `ExitUsage`.

`codeFor`'s `isUsageError` fallback is then a *thin defensive net* we own — a
small prefix check over the residual Cobra-generated strings (notably
`"unknown command"`) for cases without a typed hook — not the main mechanism. If
that net misses, the outcome degrades to `ExitInternal`, never to a wrong
*success*. This keeps the typed path authoritative and the fragile string match
secondary.

### Silent rendering for already-reported outcomes

`lint` prints its findings (human or `--json`) to stdout and must still exit
non-zero. Returning an error re-triggers Fang's handler, which would print a
redundant styled "Error: …" to stderr — and pollute the stderr channel an agent
reads. So the found-problems error is **silent**:

```go
// lint.go, after rendering result:
if !result.Valid {
	return &CodedError{Code: ExitProblems, Silent: true, Err: result.Err()}
}
```

The custom handler honors that flag and otherwise defers to Fang's default, so
usage and internal errors still render normally:

```go
func errorHandler(w io.Writer, styles fang.Styles, err error) {
	var ce *CodedError
	if errors.As(err, &ce) && ce.Silent {
		return
	}
	fang.DefaultErrorHandler(w, styles, err)
}
```

`internal/lint` stays exit-code-agnostic: it reports `Valid`/findings, and the
CLI layer (which already holds the `Result`) decides the category. The dependency
direction — `cli` → `lint` — is unchanged; no new package is needed.

### `init --json`: a result receipt on stdout

Add `--json` to `init`. Without it, behavior is unchanged (scaffold to file,
`Created …` confirmation to stderr). With it, `init` emits a receipt to stdout
and writes nothing to stderr on success:

```go
type InitReceipt struct {
	SchemaVersion int              `json:"schemaVersion"`
	Path          string           `json:"path"`
	Created       bool             `json:"created"`
	NextActions   []receipt.Action `json:"nextActions"`
}
```

The receipt describes the *action*, not the scaffold contents — that is what
makes `--json` meaningful for a command whose human output is a confirmation. The
`nextActions` (lint the new file) appear in-band, satisfying the convention that
agents receive next actions as data under `--json`.

**Move the next-action element to a neutral home; do not redefine it and do not
leave it in `lint`.** The element today lives as
[`lint.Action`](../../../internal/lint/result.go) (`{id, label, command}`) only
because `lint` was the first command to emit one. But a "next action" is an
*agent-receipt* concept, not a lint concept: `init` emits the same `{id, label,
command}` shape for a non-lint purpose. Reusing `lint.Action` from `internal/cli`
would read as though `init` has something to do with linting, and the only
argument for it — `internal/cli` already imports `internal/lint`, so reuse adds
no edge — justifies *convenience*, not *ownership*. So this change extracts the
type to a neutral `internal/receipt` package that both `lint` and `cli` import;
`Result.NextActions` and `InitReceipt.NextActions` both carry `[]receipt.Action`.

This is a **concept-ownership** decision, not an abstraction-extraction one, so
the rule-of-three "wait for a third consumer" instinct does not apply — see
[Designing Go packages](../../../docs/guides/design-go-packages.md), which this
change motivated. The contract type is its own concept with one consumer or two;
counting emitters answers the wrong question. The move is also cheap and safe:
the JSON wire shape is unchanged, so no agent notices the Go type's new home.
(`internal/diagnostics` is a docs-only meta-model folder, not a Go package, so it
is not a candidate home — hence a fresh `internal/receipt`.)

**`schemaVersion` is a per-contract namespace, not a shared counter.** `init`'s
receipt and `lint`'s result each version their *own* JSON contract, so they keep
independent `schemaVersion` values (`lint` today holds a private
`const schemaVersion = 1`); `init` gets its own. What is shared is the
*convention* — the field name and that it is a per-contract integer — not the
number. The two must not be wired to a single constant.

**One source for the next action, two renderings.** `init` today hardcodes the
human line `"Next: qualitymd lint <path>"`. Under `--json` the same next action
becomes a `receipt.Action`. Build the `[]receipt.Action` once and render the human
confirmation line from it, so the prose and the structured receipt cannot drift —
the same discipline `lint` already uses, rendering its human output from the
`Result`.

Two interactions to pin down, both required by the spec:

- **`-` plus `--json` is a usage error.** `-` selects raw-scaffold passthrough to
  stdout; `--json` selects receipt output to stdout. Both claim stdout with
  incompatible payloads, so the combination is rejected with `ExitUsage` rather
  than given a contrived precedence. Deterministic and easy to explain.
- **Overwrite refusal under `--json`.** When the target exists and `--force` is
  absent, `init` refuses. This is not a malformed invocation and not a bug; it is
  the command *declining to complete the requested action*. It maps to
  `ExitInternal` under the broadened reading below, and under `--json` it emits a
  distinct **error object** to stderr — `{schemaVersion, path, reason}`, not a
  success receipt and not an `{ok:false}`-wrapped success shape. Success goes bare
  to stdout; the error is told apart by exit code and channel, not by a body flag
  (see the receipt-shape resolution under Open questions).

### One refinement back to the spec

Implementing the above surfaces that the *internal error* category is better read
as **"the command could not complete the requested action"** (I/O failure,
unmet precondition such as overwrite refusal, or a bug) than the spec's narrower
"I/O failure, a bug." The functional spec and the durable CLI-spec section should
adopt the broader phrasing so guarded refusals have a home. This is the one place
the design feeds a wording change back into the spec before it lands.

## Alternatives

- **Sniff `os.Args` / error strings as the primary usage detector.** Rejected as
  the main path: it duplicates Fang's acknowledged hack and drifts with Cobra's
  message wording. Typed tagging at the `FlagErrorFunc` / `Args` source is stable;
  the string check survives only as a labelled fallback.
- **Let `lint` return a normal error and accept the extra stderr line.** Rejected.
  A redundant styled error on stderr pollutes the channel agents read and
  duplicates information already on stdout, especially under `--json`. The silent
  flag costs one branch in the handler.
- **A dedicated `ExitProblems`-style sentinel package imported by `internal/lint`.**
  Rejected as unnecessary. The CLI command already has the `Result`; pushing
  exit-code semantics down into `internal/lint` would invert the dependency for
  no gain.
- **Map overwrite refusal to a new "refusal" code or to `ExitProblems`.**
  Rejected. A fifth category overcomplicates the contract for one case, and
  `ExitProblems` is reserved for a *reported negative result* (findings), not an
  action the command declined. Broadening *internal error* to "could not complete
  the action" covers it within the four categories.
- **Give `-`+`--json` a precedence instead of erroring.** Rejected. Any precedence
  silently discards one requested output; a usage error is the honest, agent-
  legible response.
- **Pick `3` for internal error instead of `70`.** Plausible, but `70`
  (`EX_SOFTWARE` from `sysexits.h`) conventionally means "internal software
  error," reinforcing that an internal-error exit signals a tool failure rather
  than an expected negative. Kept `0`/`1`/`2` for the common outcomes (success,
  findings, usage) where their conventional meanings already align with grep- and
  getopt-style tools.

## Trade-offs & risks

- **The usage-error fallback is still string-based.** Mitigated by making typed
  tagging primary and the string net a thin, owned fallback that fails safe to
  `ExitInternal`. Tests should pin each category, including an unknown-command
  invocation, so drift in Cobra's wording is caught.
- **Exit-code numbers are now a public contract.** Once agents and CI branch on
  `1` vs `2` vs `70`, changing them is breaking. They belong in the durable CLI
  spec (this change puts them there) and in tests, not just in code.
- **`init --json` adds a stable schema to maintain.** Small here (`path`,
  `created`, `nextActions`), but it is a standing obligation as future
  side-effecting commands adopt the receipt shape; keep the receipt fields and the
  `schemaVersion` convention consistent with `lint`'s result.
- **Silent errors can hide real failures if misapplied.** The `Silent` flag must
  be set *only* where the command has already reported the outcome on stdout
  (today: `lint` found-problems). A test should assert that usage and internal
  errors still render to stderr.

## Open questions

- **Receipt envelope shape** *(settled).* No shared `{ok, …}` envelope. Success
  goes **bare** to stdout (`{schemaVersion, path, created, nextActions}`, matching
  `lint`'s bare result); failure goes as a **distinct error object** to stderr
  (`{schemaVersion, path, reason}`); and the **exit code** is the success/failure
  discriminator. An agent never branches on a body flag.

  This follows the dominant CLI convention rather than the API one. `gh`, `kubectl
  -o json`, and `aws` all emit the bare resource on success with no `{ok}` wrapper,
  and `gh` ships a *categorized* exit-code contract (`0` ok, `1` error, `2`
  cancelled, `4` auth-required) — the same shape as this change's exit categories,
  confirming the exit code is the idiomatic discriminator. AWS CLI's opt-in
  structured error output is the precedent for the stderr error object: even the
  CLI that went furthest toward machine-readable errors kept success and error as
  *separate shapes told apart by channel + exit code*, not merged into one `{ok}`
  envelope. The `{ok}` / `result`-xor-`error` envelope is an HTTP/RPC convention
  (Slack Web API, JSON-RPC, GraphQL) forced by the absence of an out-of-band status
  channel; a CLI has exit code and stderr, so the envelope would only restate the
  exit code in the body — a second source of the same truth that can drift from it.
  For `lint` specifically a boolean `ok` is worse than redundant: the
  ran-with-findings (`ExitProblems`) vs could-not-complete (`ExitInternal`)
  distinction this change introduces is exactly what a single bool collapses.

  What *is* shared across commands is the **convention set**, not an envelope: the
  `schemaVersion` field (per-contract integer), the `nextActions: []receipt.Action`
  shape, and a consistent error-object shape when an error body is emitted. A
  unified `status`/envelope is reconsidered only if a future command's outcome
  genuinely cannot be expressed by exit code + bare body.
- **Where the exit-code constants live** *(settled).* In `internal/cli`, beside
  the mapping. `internal/lint` stays exit-code-agnostic; no shared `internal/exit`
  package until a second producer of coded errors exists outside `internal/cli`.
- **Where `Action` lives** *(settled).* Extract it to a neutral
  `internal/receipt` package that both `lint` and `cli` import; do not redefine
  it and do not leave it owned by `lint`. A "next action" is an agent-receipt
  concept, not a lint concept — `lint` was just its first emitter. This is a
  concept-ownership decision, so the rule-of-three "wait for a third consumer"
  guard does not apply; the move happens now (see the `init --json` section and
  [Designing Go packages](../../../docs/guides/design-go-packages.md)). The wire
  shape is unchanged, so the relocation breaks no agent.
- **Unknown-subcommand detection** *(settled).* Keep the thin, owned prefix check
  so `qualitymd bogus` maps to `ExitUsage` — option (a). Flag and argument usage
  errors are tagged through the idiomatic `FlagErrorFunc` / wrapped-`Args` hooks;
  the unknown-subcommand case has no Cobra tagging hook, so a small string check
  we own (not Fang's unexported one) covers it, failing safe to `ExitInternal` if
  Cobra's wording drifts. A test pins the unknown-command exit code so that drift
  is caught.
