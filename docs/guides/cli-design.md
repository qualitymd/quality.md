---
type: How-to Guide
title: Designing CLI interfaces
description: How to design a qualitymd command's arguments, flags, output, and errors, with per-aspect guidelines adapted from clig.dev.
tags: [cli, design, contributing]
timestamp: 2026-06-18T00:00:00Z
---

# Designing CLI interfaces

Use this guide when you're designing a new `qualitymd` command — or reshaping an
existing one — to decide its arguments, flags, output, errors, and behavior. It
gives the principles to design by and a per-aspect checklist to work through. It
covers the _how_ of design; the binding _what_ lives in the
[CLI functional spec](../../specs/cli.md) and its command sub-specs. Where the two
meet, the spec wins — this guide is the reasoning that should make the spec feel
inevitable.

The guidance here is adapted from the [Command Line Interface Guidelines](https://clig.dev/)
(clig.dev), which we treat as the authoritative source for CLI design principles.
We adopt its principles wholesale and add one standing emphasis of our own:
`qualitymd` is driven by agents and automation as much as by people, so
**agent accessibility is a first-class design constraint**, not an afterthought.
Designing for a human at a terminal and designing for an agent in a pipe are the
same discipline done well.

## Philosophy

Nine principles, in priority order. When they conflict, the earlier one usually
wins.

1. **Human-first design.** Design for a person at a terminal first, even though
   programs and agents also consume the output. Machine-readability is layered on
   top of a clear human form, never the other way around.
2. **Simple parts that work together.** Prefer small, composable commands with
   clean interfaces over monoliths. Lean on the shared currency of the shell —
   pipes, exit codes, `stdout`/`stderr`, JSON — so a command drops into pipelines
   and scripts the author never anticipated.
3. **Consistency across the program.** Follow established conventions so the CLI
   is guessable. A flag means the same thing everywhere; output looks the same
   everywhere. Break a convention only deliberately, and say why where you do.
4. **Saying (just) enough.** Find the middle between cryptic silence and a wall of
   noise. Say what changed, what to do next, and nothing the user didn't ask for.
5. **Ease of discovery.** Make features findable from inside the tool: thorough
   `--help`, worked examples, and suggestions when the user is close but wrong.
   The CLI should teach itself.
6. **Conversation as the norm.** A session is a back-and-forth. The user tries
   something, the tool responds, they adjust. Good errors, confirmations, and
   next-action hints keep that loop tight.
7. **Robustness.** Handle unexpected input without falling over. Keep operations
   idempotent and recoverable where you can, and make the tool _feel_ solid —
   responsiveness and considered detail are part of robustness, not polish on top
   of it.
8. **Empathy.** Treat the CLI as something people _want_ to use. Thoughtful
   defaults, forgiving input, and clear language signal that you want the user to
   succeed.
9. **Chaos.** Conventions exist to be broken when following them would demonstrably
   hurt. Do it with intent and a stated reason — not by accident.

## The basics

A few things every command does, no exceptions:

- **Use a real argument parser.** Don't hand-roll flag parsing. Our stack is
  [Cobra](https://github.com/spf13/cobra) for command and flag structure,
  [Charm Fang](https://github.com/charmbracelet/fang) for the invocation harness,
  and [Lip Gloss](https://github.com/charmbracelet/lipgloss) for terminal output.
  A requirement that can only be met by fighting these libraries is a signal to
  reshape the requirement.
- **Exit zero on success, non-zero on failure.** And do it through stable,
  documented categories — see [Errors and exit codes](#errors-and-exit-codes).
- **Primary output to `stdout`, everything else to `stderr`.** The result is the
  payload; diagnostics, progress, and footers are not. This is what makes piping
  `stdout` safe.
- **Don't assume a TTY.** A command must run unattended — in CI, in a pipe, under
  an agent — without blocking on a prompt.

## Help

Help is the front door. Most users meet a command through its help before its
docs.

- Respond to both `-h` and `--help`.
- Default (no-args, where there's nothing to do) shows _concise_ help: a one-line
  description, one or two examples, the most common flags, and a pointer to
  `--help` for the rest.
- `--help` shows the full text: every flag, grouped and ordered by how often it's
  used — common things first, not alphabetically.
- **Lead with examples.** A couple of real invocations teach faster than a
  paragraph of prose. Show the common case first.
- **Document structured input, not just invocation.** When a command reads an
  author-supplied payload (JSON/YAML via `--file` or stdin), help must make the
  payload's shape discoverable — fields, types, required-ness, allowed values —
  and show at least one complete valid payload, not only how to call the command.
  A command whose only account of its input is "submit and read the rejection"
  fails the teach-itself test. See [Structured input](#structured-input).
- Format for the terminal: bold headings and structure via the styling layer, not
  raw escape codes baked into strings.
- Point to where help continues — the web docs, the repo — so a stuck user has a
  next step.

## Documentation

Help is for the terminal; docs are for depth.

- Provide web documentation that is searchable and linkable, and keep the
  terminal help in sync with it.
- Make docs reachable _from_ the tool — a link in help text, a `spec`-style
  command that emits the canonical artifact.
- Examples first here too. People copy, adapt, then read.
- **Make drift impossible, don't rely on diligence.** "Keep help in sync" is a
  wish, not a mechanism — stale examples are what happens when sync is a habit.
  Any example, payload, or schema shown to a user must be generated from, or
  golden-tested against, the implementation, so it cannot drift from the binary
  it documents.

## Output

Output is where human-first and machine-readable have to coexist without
fighting.

- **Human-readable by default.** The default form is for a person reading a
  terminal.
- **Machine-readable on request, via `--json`.** Offer structured output spelled
  `--json` — never `--format json` or a per-command variant — wherever it's
  meaningful, so a caller can reach for it without checking whether it exists. The
  carve-out is a command whose output already _is_ a verbatim artifact meant to be
  redirected; wrapping that in JSON adds nothing.
- **Do not add a second JSON mode for JSON artifacts.** A command whose stdout is
  already the JSON artifact, such as a schema, example payload, or stored JSON
  record, should not also offer a JSON result wrapper. It may recognize `--json`
  only to fail with a targeted usage error that says the command already emits
  JSON and should be rerun without `--json`.
- **No format auto-detection.** Passing `--json` is the only switch into JSON.
  Don't silently change the payload based on whether `stdout` is a pipe; changing
  _styling_ off a TTY is fine, changing _content_ is not.
- **Say what changed.** After a side effect, state what happened in plain terms.
  Under `--json`, a side-effecting command emits a result receipt describing what
  it did, so `--json` is meaningful even when the human output is just a
  confirmation.
- **Suggest the next step.** Where there's an obvious follow-up, offer it (see
  [Conversation and next actions](#conversation-and-next-actions)).
- **Color and glyphs, sparingly and conditionally.** Style is a terminal-only
  convenience over a canonical plain form. It may add color and a status glyph; it
  **must not** change the words, order, or facts. Apply it only when `stdout` is a
  terminal and `NO_COLOR` is unset — piped and agent-driven output is the plain
  form, byte for byte.
- **Page long artifacts, but never load-bearingly.** A long verbatim output may go
  through the user's pager when `stdout` is a terminal; when it isn't, or no pager
  exists, write directly instead.
- **`stderr` is not a logfile.** Don't dump levelled log lines (`INFO`, `DEBUG`)
  to `stderr` by default. Keep developer debugging behind an explicit flag.

## Structured input

Some commands take an author-supplied structured payload — JSON or YAML read from
`--file` or stdin — rather than a handful of flags. This is the mirror of
[Output](#output), and it deserves the same care: where output makes results both
human-readable and machine-readable, structured input must be discoverable,
validatable, and answerable in the caller's own terms. It is also where
agent-accessibility is won or lost, because an agent cannot infer a payload's
shape from context the way a person at a terminal sometimes can — it has to be
told.

- **The contract is discoverable from inside the tool.** Help — or a schema or
  dry-run affordance it points to — documents every field: name, type, whether
  it's required, and the allowed values for any enum. Discovery must not require
  reading source or guessing.
- **Show a complete, valid example.** At least one payload that would actually be
  accepted, with every required field present and real enum values — not a
  fragment. Callers copy, adapt, then read.
- **Offer a way to validate without committing.** A command that consumes a
  payload and writes something should let a caller check the payload first —
  `-n/--dry-run` or an equivalent — validating fully and reporting what would
  happen without persisting anything. This is "validate early, before
  irreversible work" made reachable, and it is what keeps schema discovery from
  meaning "write a junk record, then delete it."
- **Aggregate validation.** Report every problem with a payload in one pass, not
  just the first one encountered. Serial round-trips — fix one field, resubmit,
  discover the next — are the structured-input equivalent of a wall of noise:
  they make the caller do the tool's bookkeeping.
- **Answer in the input's own vocabulary.** A validation error names the field as
  the author wrote it — the JSON/YAML key — with the expected type and allowed
  values. Never echo an internal or language-level field name the caller never
  typed; making them re-map a struct name back to a key is the tool leaking its
  implementation.

## Arguments and flags

The distinction matters: **arguments** are positional (order carries meaning);
**flags** are named (`-x`, `--name`) and order-independent.

- **Prefer flags to positional arguments.** Flags are self-documenting and leave
  room to grow; a wall of positional args is hard to read and easy to get wrong.
  Multiple positional args are fine for the obvious case (a list of files);
  positional args with _different meanings_ are not.
- **Give every flag a full `--long` form;** add a one-letter short form only for
  the genuinely common ones.
- **Use the conventional names** when one exists, so muscle memory carries over:
  `-h/--help`, `--version`, `-o/--output`, `-q/--quiet`, `-v/--verbose`,
  `-f/--force`, `-n/--dry-run`, `--json`, `--no-color`.
- **Have sensible defaults.** The common case should need few or no flags. Don't
  make the user remember what they could have inferred.
- **Don't require interactive prompts.** Prompting is acceptable as a convenience
  when input is missing _and_ `stdin` is an interactive terminal, but there must
  always be a flag or argument that supplies the same input non-interactively.
- **Support `-` for stdin/stdout** where a command reads or writes a file, so it
  composes in pipelines.
- **Never take secrets from a flag.** Flags leak into shell history and process
  listings. Read secrets from a file or `stdin` instead.
- **Aim for order-independence** of flags and arguments where you reasonably can.

## Subcommands

Once a tool grows past a handful of actions, it becomes a `noun verb` (or
`verb noun`) tree. Consistency across that tree is the whole game.

- **Be consistent across subcommands:** the same flag name means the same thing,
  output is formatted the same way, errors read the same.
- **Pick one grammar and hold it.** `noun verb` (e.g. `evaluation report build`) is
  the more common shape; whichever you choose, don't mix.
- **No ambiguous or near-duplicate names,** and no catch-all subcommand that
  guesses an intent the user never stated.
- **Don't accept arbitrary abbreviations** of subcommand names — they quietly
  become an interface you can never change.

## Errors and exit codes

Errors are the most-read output a tool produces, because they're read when the
user is already stuck. Treat them as the conversation's most important turn.

- **Rewrite expected errors for humans.** Catch the common failures and say what
  went wrong, why, and what to try — not a raw stack trace.
- **Put the signal where the eye lands** — the important line at the end, not
  buried mid-scroll — and group similar errors instead of repeating them.
- **For author-supplied structured input, aggregate and speak the caller's
  vocabulary.** Validation errors report the whole set of problems at once and
  name each field as it appears in the payload, not as an internal symbol — see
  [Structured input](#structured-input).
- **For unexpected errors,** give enough to file a good bug report and make
  reporting easy.
- **Signal outcome through stable exit-code categories,** so a caller can branch
  without parsing text. `qualitymd` uses:

  | Code | Category               | Meaning                                                            |
  | ---- | ---------------------- | ------------------------------------------------------------------ |
  | `0`  | Success                | The command did its job.                                           |
  | `1`  | Ran but found problems | Completed normally, but the result is a reportable negative.       |
  | `2`  | Usage error            | Malformed invocation: unknown flag, bad argument, unknown command. |
  | `70` | Internal error         | Could not complete: I/O failure, unmet precondition, or a bug.     |

  This table is binding; it lives in the [CLI spec](../../specs/cli.md) and is
  reproduced here as part of the reasoning.

## Interactivity

- **Only prompt when `stdin` is an interactive TTY.** A pipe or redirect means no
  human is there to answer.
- **Honour a `--no-input` escape** that disables every interactive element, for
  scripts and agents.
- **Hide secret input** (disable echo when reading a password).
- **Keep Ctrl-C working.** The user must always be able to get out.

## Robustness

Robustness is both objective (it doesn't break) and subjective (it doesn't _feel_
fragile).

- **Validate input early,** before doing irreversible work.
- **Be responsive:** print _something_ within ~100ms, even if the real work takes
  longer. Show progress for anything slow, and an estimate where you can.
- **Set timeouts** on network operations with reasonable defaults.
- **Make operations recoverable and, where possible, idempotent** — safe to
  re-run after a transient failure.
- **Anticipate misuse:** being scripted, run twice at once, run in a half-broken
  environment. Degrade with a clear message rather than a crash.

## Determinism

Beyond clig.dev, because `qualitymd` is consumed by automation:

- **Same input and file state, same output.** No timestamps, ordering jitter, or
  sampling in the payload. Determinism is what lets an agent or CI step diff,
  cache, and assert against output without flakiness.
- This is what makes the human-readable form safe to also treat as data: nothing
  in it varies run to run except what the inputs varied.

## Conversation and next actions

A command may close with a short list of _next actions_ — the commands the caller
would most plausibly run next.

- **Concrete over vague:** a runnable command to copy
  (`qualitymd lint QUALITY.md`), not "consider validating your file."
- **Deterministic:** the same outcome yields the same suggestions, derived from
  the result rather than ranked or sampled.
- **Subordinate to the payload:** in human output they render as a footer on
  `stderr`, after the primary output, so piping `stdout` is never polluted. Under
  `--json` they appear in-band as a `nextActions` array, so agents receive them as
  data.
- **Useful on success and failure alike:** after a scaffold, point at linting;
  after a failed check, point at re-running it or opening the offending file.

## Configuration

Match the mechanism to how often, and at what scope, a setting varies:

- **Varies per invocation** → a flag (and sometimes an environment variable).
- **Stable per user or project** → environment variables or a config file.
- **Stable within a project, shared by everyone** → a version-controlled config
  file.
- **Apply precedence** flags → environment → project config → user config →
  system config, most specific winning.
- **Get consent before writing files you don't own.** Modifying a user's config or
  unrelated files is a cross-boundary action — make it explicit.

## Environment variables

- Respect the general-purpose ones the ecosystem already defines: `NO_COLOR`,
  `PAGER`, `EDITOR`, `TERM`, `HOME`, `HTTP_PROXY`, `DEBUG`.
- Name your own in `UPPER_SNAKE_CASE`, and don't commandeer standard POSIX names.
- **Never read secrets from the environment** — it leaks into logs and child
  processes too easily.

## Naming

- Pick a simple, memorable, lowercase name; dashes if needed, no surprises.
- Keep it short and comfortable to type — it's typed thousands of times.
- Avoid clashing with common existing commands.

## Distribution

- **Ship a single self-contained binary** where you can; it's the easiest thing to
  install, run, and reason about.
- Make uninstalling as clear as installing.

## Analytics

- **Don't phone home without explicit consent.** Prefer opt-in; if you ever
  collect anything, be transparent about what, why, and for how long. For a tool
  meant to run unattended in others' pipelines, the safe default is to collect
  nothing.

## Further reading

- [Command Line Interface Guidelines](https://clig.dev/) — the source these
  guidelines adapt.
- [The CLI functional spec](../../specs/cli.md) — the binding requirements for the
  `qualitymd` surface.
- [POSIX Utility Conventions](https://pubs.opengroup.org/onlinepubs/9699919799/basedefs/V1_chap12.html)
  and the [GNU Coding Standards](https://www.gnu.org/prep/standards/html_node/Command_002dLine-Interfaces.html).
