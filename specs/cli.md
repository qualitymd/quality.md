# CLI specification

Functional specification for the `qualitymd` CLI surface.

> **Status:** high-level functional spec. Exact JSON field names and payload
> details are illustrative; the authoritative payload contract will live in a
> forthcoming evaluation-lifecycle spec.

The CLI is deterministic. It never calls a model. Skills carry judgment and
orchestration (see [`skills.md`](./skills.md)).

Detail docs:

- [`cli-init.md`](./cli-init.md) - scaffold a starter target-tree `QUALITY.md`.
- [`skills.md`](./skills.md) - skill-layer orchestration.

Forthcoming (not yet written):

- A deterministic structural validator (`lint`) spec.
- An `evaluation` and `result` lifecycle spec: report rollup and CLI-to-skill
  payloads.
- A federation spec: multiple models as one grafted target tree.

## Split

- **CLI:** parse and inspect the model, resolve the target tree and `source`
  manifests, persist evaluations, record verdicts, roll up targets and factors,
  render reports.
- **Skills:** choose work order, perform each requirement's `assessment`, produce
  a finding, choose the rating level whose criterion the finding satisfies, and
  write evidence through `result set`.

The format's chain is:

```text
assessment -> finding -> rating criteria -> result (rating)
```

## Command Surface

### Top-Level

| Command                 | Purpose                                                                            | Output         |
| ----------------------- | ---------------------------------------------------------------------------------- | -------------- |
| `qualitymd init`        | Scaffold a starter `QUALITY.md` target tree.                                       | `./QUALITY.md` |
| `qualitymd lint [file]` | Validate structure: Targets, scoped factors, requirements, references, body shape. | JSON findings  |

### `qualitymd model`

| Command                                         | Purpose                                                                                                                                               | Output |
| ----------------------------------------------- | ----------------------------------------------------------------------------------------------------------------------------------------------------- | ------ |
| `model show [--requirement <address>] [--json]` | Parsed model: recursive target tree, requirements by target-tree address, resolved `source` manifests, loaded `assessment` text, active rating scale. | JSON   |

### `qualitymd evaluation` (`eval`)

| Command                                                              | Purpose                                                                                                          | Transition |
| -------------------------------------------------------------------- | ---------------------------------------------------------------------------------------------------------------- | ---------- |
| `evaluation create [--model <path>] [--target <path>] [--from <id>]` | Create or re-enter the living run for one model and CLI run target; enumerate selected target-tree requirements. | open       |
| `evaluation list`                                                    | List active and archived runs.                                                                                   | -          |
| `evaluation show [<id>]`                                             | Run manifest: status, result set, rollup, verdict.                                                               | -          |
| `evaluation report [<id>] [--fail-on <level>]`                       | Render report and fail when the rollup trips the gate.                                                           | -          |
| `evaluation archive [<id>] --as <name>`                              | Snapshot a run.                                                                                                  | archived   |
| `evaluation delete <id>`                                             | Remove a run.                                                                                                    | abandoned  |

### `qualitymd result`

| Command                                                | Purpose                                                               | Transition |
| ------------------------------------------------------ | --------------------------------------------------------------------- | ---------- |
| `result list [--status pending,stale,...] [--json]`    | Query results by state. There is no `next` cursor; skills order work. | -          |
| `result show <address> [--json]`                       | Resolved target-tree requirement payload for a skill.                 | -          |
| `result set <address> --rating <level> --evidence ...` | Record the skill's finding and rating.                                | recorded   |
| `result skip <address> --reason ...`                   | Mark deliberately not assessed.                                       | skipped    |
| `result reset <address>`                               | Return to pending.                                                    | pending    |

## Two Meanings Of Target

The CLI has a run-selection flag named `--target <path>`: it selects the subject
path for an evaluation run. The schema has `targets:`: recursive **Targets**
inside a `QUALITY.md`. The CLI run target selects which subject tree is being
evaluated; schema Targets describe how the model decomposes that subject.
Docs use **CLI run target** for the flag and **Target** for the schema type.

## Structural Tier vs. Management Tier

- **`lint`** asks whether the file is structurally valid: Target shape,
  scoped factor identity, requirement `assessment`s, `source` and referenced
  assessment paths, rating scales, and body headings.
- **`model` / `evaluation` / `result`** manage the recorded state of evaluation.
  They do not judge requirements; skills do.

The intended order is:

```sh
qualitymd init
qualitymd lint
qualitymd evaluation create
qualitymd result list --status pending,stale --json
qualitymd result show <address> --json
qualitymd result set <address> --rating <level> --evidence ...
qualitymd evaluation report --fail-on unacceptable
```

There is no separate settled `check` verb. The current binary's `check` command is
a placeholder predating this surface.

## Evaluation Lifecycle

The run model will be specified in a forthcoming evaluation-lifecycle spec. In
brief:

- one living run per `(model, CLI run target)`, re-entered in place;
- stored under `.quality/evaluations/<slug>/`;
- always mutable; git history is the audit layer;
- result states are `pending`, `recorded`, `skipped`, and `stale`;
- reports roll requirements through factors and Targets to an overall
  verdict.

## Conventions

- Resource commands are singular: `evaluation create`, `result show`.
- `-f, --file <path>` selects one `QUALITY.md`; default is `./QUALITY.md`.
- With no explicit file and multiple discovered models, federation rules apply.
- Paths in model files resolve relative to the containing `QUALITY.md`.
- `--json` emits schema-versioned machine-readable output. `lint` and
  `model show` default to JSON.
- Non-interactive mode is used when stdin/stdout is not a TTY, `--json` is
  passed, or `--non-interactive` is passed.

## Exit Codes

| Code | Meaning                                                                      |
| ---- | ---------------------------------------------------------------------------- |
| `0`  | Command ran and any gate passed.                                             |
| `1`  | Gate failure: lint errors or `evaluation report --fail-on` tripped.          |
| `2`  | Tool failure: bad flags, unreadable files, parse failure, or internal error. |

Recording a low rating with `result set` exits `0`; the gate is
`evaluation report`.

## `.quality/`

```text
.quality/
  config.yaml
  evaluations/
    <slug>/
    archive/
      <name>/
```

Everything under `.quality/` is committed unless a later policy says otherwise.
In a federation, runs live beside each model and shared configuration lives at the
scan root. The built-in quality meta-model ships as a normal recursive
Target-tree `QUALITY.md` file and is used by `improve-quality-md`.

## Open Questions

- Final field names for CLI-to-skill payloads and staleness serialization.
- Whether the CLI emits PR comments/check annotations directly.
- Whether to add a `spec` verb that prints the format spec for agents.
- Whether the bundled meta-model can later be extended or replaced by a project.
