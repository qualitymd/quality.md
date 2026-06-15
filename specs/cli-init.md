# CLI: `init`

> Detail doc for the scaffolding command. See [`cli.md`](./cli.md) for the full
> command surface and shared conventions, [`cli-lint.md`](./cli-lint.md) for the
> structural tier, and [`cli-evaluate.md`](./cli-evaluate.md) for the deep
> semantic tier.

```bash
qualitymd init                       # write a starter ./QUALITY.md
```

## Purpose

`init` answers the bootstrapping question every other command assumes away:
**where does the first `QUALITY.md` come from?** `lint`, `evaluate`, and
`evaluate-model` all take an existing model as input; `init` produces one.

`init` is **deterministic and offline**: it writes a minimal, schema-commented
starter `QUALITY.md` — a few example factors and a Markdown body skeleton, on
the default Outstanding/Target/Minimum/Unacceptable scale — that the author then
fills in. It has no opinion
about your codebase and does no agentic work; it gives you a well-formed file to
edit and a clear next step.

Drafting a model that is actually *tailored* to a codebase is not `init`'s job —
that falls out of the normal authoring loop below. Running
`evaluate-model` against the scaffold surfaces the requirements and
factors a real codebase warrants, as structured `proposedChange` suggestions
(see [`cli-evaluate.md`](./cli-evaluate.md#one-engine-two-targets) and the
agent-friendly next-action pattern in
[`cli.md`](./cli.md#agent-friendly-ci-patterns)). So the cheap deterministic
scaffold and the expensive agentic tailoring stay cleanly separated across two
commands, rather than hiding an agent inside `init`.

## The authoring loop

`init` is step zero of the intended workflow; the other commands refine what it
produces, each suggesting the next via structured next-action output (see
[`cli.md`](./cli.md#agent-friendly-ci-patterns)):

```text
qualitymd init                    # scaffold a starter model                  (deterministic)
qualitymd lint                    # is the file well-formed?                  (structural)
qualitymd evaluate-model          # which requirements does the subject want? (semantic)
qualitymd evaluate                # does the subject satisfy them?            (semantic)
```

The scaffold is explicitly a *starting point*, not a finished model. `init` says
so in the file it writes (a header comment) and in the next-action it emits:
"edit the placeholders, then run `qualitymd lint`."

## What it writes

A single `QUALITY.md` with commented frontmatter and a short body skeleton:

```markdown
---
# QUALITY.md — quality model for <project>. See SPECIFICATION.md.
# Each factor has `requirements`, nested `factors`, or both.
# Each requirement declares a single `prompt` assessment.
# Ratings default to outstanding / target / minimum / unacceptable; add a `ratings:` block to customize.
factors:
  functionality:
    requirements:
      "<what must be true>":
        target: "./src"           # optional; path or glob under evaluation
        prompt: "<inline criteria, or ./path/to/standard.md>"
  testability:
    requirements:
      "<what good testing looks like here>":
        prompt: "<inline criteria, or ./path/to/standard.md>"
---

# Quality model — <project>

## Overview
<!-- What this system or component is, who depends on it, and what "good" means
     here. -->

## Scope
<!-- What the model covers and what it deliberately leaves out (out of scope),
     including dependencies it relies on but does not own. -->

## Needs
<!-- What matters, and to whom — the plain-language statements the requirements
     answer to. -->

## Risks
<!-- What goes wrong, and for whom, if a need is not met. -->

## Factors
<!-- One subsection per factor, mirroring the frontmatter. -->
### Functionality
<!-- What this factor means here, how you would know it is met, any trade-offs. -->
### Testability

## Known gaps
<!-- Quality concerns known to matter but deliberately not addressed yet, each
     with a brief reason. -->
```

The skeleton seeds the body with the spec's recommended sections (Overview,
Scope, Needs, Risks, Factors, Known gaps), leaning on prose to match the project's
"prose over tokens" bet: the comments prompt the author to *write the
reasoning*, not just fill in token values. The factor subsections mirror the
example factors in the frontmatter.

### Optional config (`--config`)

Also scaffold the project quality home (see
[`cli.md`](./cli.md#configuration--quality)):

```text
./.quality/
  config.yaml        # commented defaults: rigor, output dir, ignore globs
```

Off by default — `init` writes only `QUALITY.md` unless `--config` is passed.

## Interactive by default

`init` is the one human-facing, human-authored entry point in the CLI, so it
**prompts by default** when run in a terminal — matching the convention of
scaffolding commands like `npm init` and `gh repo create`. The prompts ask for a
few fields (project name, an initial factor or two) and weave the
answers into the scaffold in place of the placeholders; everything not asked about
(the body skeleton, the remaining placeholders) is still seeded for the author to
edit.

Interactivity follows the CLI's shared rule (see
[`cli.md`](./cli.md#shared-conventions)): `init` runs the **non-interactive** path
— no prompts, the fixed placeholder scaffold, never blocking — whenever stdin/stdout
is **not a TTY** (piped, redirected, or in CI), or `--non-interactive` or `--json`
is passed. So `qualitymd init` on a terminal is interactive, while an agent, a
pipe, or a CI step gets the deterministic file with no prompt and no hang. Prompts
are written to stderr, so `--json` result output on stdout is never contaminated.

Both paths produce the same well-formed file on the default
Outstanding/Target/Minimum/Unacceptable scale and emit the same next-action ("edit the placeholders, then run `qualitymd lint`").

## Flags, exit codes

`init` always writes `./QUALITY.md` in the current directory — there is no path
argument and no stdout mode. This keeps the bootstrapping step a single,
predictable action.

`init` writes one fixed starter scaffold (see [What it writes](#what-it-writes)):
the example factors above, relying on the default
Outstanding/Target/Minimum/Unacceptable scale. There is
no template or rating-scale choice — the author edits the file afterward.

Flags (shared flags are in [`cli.md`](./cli.md#shared-conventions)):

- `--config` — also scaffold `./.quality/config.yaml`.

Exit codes follow the shared three-code convention (see
[`cli.md`](./cli.md#machine-readable-result-contract)). `init` has no gate, so it
never emits the gate-verdict code `1`:

- **`0`** on a successful write.
- **`2`** — **tool failure:** `init` could not complete the write. Either
  `./QUALITY.md` already exists (`init` never clobbers an existing model) or the
  write itself failed (permissions, unwritable directory).
