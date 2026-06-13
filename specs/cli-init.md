# CLI: `init`

> Detail doc for the scaffolding command. See [`cli.md`](./cli.md) for the full
> command surface and shared conventions, [`cli-lint.md`](./cli-lint.md) for the
> structural tier, and [`cli-evaluate.md`](./cli-evaluate.md) for the deep
> semantic tier.

```bash
qualitymd init                       # write a starter ./QUALITY.md
qualitymd init docs/QUALITY.md       # write to a chosen path
qualitymd init -                     # print to stdout, write nothing
```

## Purpose

`init` answers the bootstrapping question every other command assumes away:
**where does the first `QUALITY.md` come from?** `lint`, `evaluate`, and
`evaluate-model` all take an existing model as input; `init` produces one.

`init` is **deterministic and offline**: it writes a minimal, schema-commented
starter `QUALITY.md` — a few example factors, the default `pass`/`fail` scale,
and a Markdown body skeleton — that the author then fills in. It has no opinion
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
# QUALITY.md — quality model for <project>. See docs/spec.md.
# Each factor has `requirements`, nested `factors`, or both.
# Each requirement declares exactly one assessment: `prompt` or `bash`.
ratings:                  # optional; omit to default to pass / fail
  pass: { displayName: "Pass" }
  fail: { displayName: "Fail" }
factors:
  functionality:
    requirements:
      "<what must be true>":
        target: "./src"           # optional; path or glob under evaluation
        prompt: "<inline criteria, or ./path/to/standard.md>"
  testability:
    requirements:
      "unit tests pass":
        bash: "<your test command>"
---

# Quality model

Document the rationale here — why these factors, what "good" means for this
system, what is deliberately out of scope. Prose is the heart of the model;
the frontmatter is the machine-readable summary of it.
```

The template leans on prose, matching the project's "prose over tokens" bet:
placeholders prompt the author to *write the reasoning*, not just fill in token
values.

### Optional config (`--config`)

Also scaffold the project quality home (see
[`cli.md`](./cli.md#configuration--quality)):

```text
./.quality/
  config.yaml        # commented defaults: rigor, output dir, ignore globs
```

Off by default — `init` writes only `QUALITY.md` unless `--config` is passed.

## Flags, exit codes

Flags (shared flags are in [`cli.md`](./cli.md#shared-conventions)):

- `path` — positional output path, or `-` for stdout. Defaults to `./QUALITY.md`
  / `-f`.
- `--template minimal|example` — starter content: `minimal` (a near-empty
  scaffold) or `example` (a few illustrative factors, the default).
- `--ratings pass-fail|letter` — which rating scale to seed (`pass-fail` default;
  `letter` writes the A–E scale from `docs/spec.md`).
- `--config` — also scaffold `./.quality/config.yaml`.
- `--force` — overwrite an existing target file (see below).

Exit codes:

- **`0`** on a successful write.
- **`1`** if the target file already exists and `--force` was not given — `init`
  never clobbers a model by default. (Writing to `-`/stdout never touches disk
  and so never trips this.)
- **`1`** on a write error.

## Open questions

- **Interactivity.** Pure flags-only (current assumption, agent-friendly and
  scriptable) vs. an interactive prompt that asks for project name, factors, and
  test command. Interactive is friendlier for humans but cuts against the
  agent-consumable, non-interactive grain of the rest of the CLI.
- **Body skeleton vs. frontmatter only.** Whether template mode emits the
  Markdown body skeleton (current) or just the frontmatter, leaving prose
  entirely to the author. The "prose over tokens" philosophy argues for seeding
  the body.
- **Template registry.** `--template` is a fixed enum in v1. Whether to support
  user/shared templates (e.g. a `web-service` or `library` starter) is deferred.
