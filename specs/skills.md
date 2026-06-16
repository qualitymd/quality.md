# Skills specification

The skill layer supplies judgment and orchestration on top of the deterministic
`qualitymd` CLI.

## Purpose

- **CLI:** parse target-tree models, resolve `source` manifests, persist
  evaluations, record findings and ratings, roll up targets/factors, render
  reports.
- **Skills:** decide what to inspect, perform each requirement's `assessment`,
  produce a finding, choose the best rating level whose criterion is met, and
  write evidence back through `result set`.

Earlier drafts put the agentic engine inside CLI commands. After the inversion,
agentic work lives in skills; the CLI exposes `init`, `lint`, `model`,
`evaluation`, and `result`.

## Skill Set

| Skill                | What it does                                                                                                                                            | Model used                      |
| -------------------- | ------------------------------------------------------------------------------------------------------------------------------------------------------- | ------------------------------- |
| `setup-quality-md`   | Onboards a project: runs/scaffolds `qualitymd init`, drafts a first target-tree `QUALITY.md` grounded in real needs and risks, and sets up `.quality/`. | none; authors the model         |
| `evaluate-quality`   | Evaluates the subject against its `QUALITY.md`.                                                                                                         | the project's `QUALITY.md`      |
| `improve-quality-md` | Diagnoses then improves the `QUALITY.md` artifact itself.                                                                                               | the built-in quality meta-model |

`improve-quality-md` diagnose mode is the same evaluation loop with the bundled
meta-model as model and the user's `QUALITY.md` as the CLI run target, followed
by an edit phase.

## Orchestration Contract

The CLI never emits an agent prompt and never chooses the next result. A skill
runs:

1. `qualitymd evaluation create`
2. `qualitymd result list --status pending,stale --json`
3. For each result:
   - `qualitymd result show <address> --json`
   - perform the `assessment` against the resolved `source`;
   - record the finding and rating with `qualitymd result set <address> ...`.
4. `qualitymd evaluation report --fail-on unacceptable`

The result address is the target-tree address, to be specified in a forthcoming
evaluation-lifecycle spec.

## CLI To Skill Interface

The payloads will be specified authoritatively in a forthcoming
evaluation-lifecycle spec:

- `result show` gives the target-tree address, resolved `assessment`, target
  `source` manifest, primary and secondary factor context, active rating
  criteria, result state, and sufficiency guidance.
- `result set` records the chosen rating and structured finding evidence.
- A staleness hash defines when a recorded result reopens because the model,
  assessment, criteria, source selection, or source contents changed.

## Per-Skill Detail

### `setup-quality-md`

Owns authoring judgment for the first model: what Targets the project
needs, which factors belong on which target subtrees, which direct and lensed
requirements answer the project's needs, and what body context explains it.
It uses `init` and `lint` but does not evaluate the subject.

### `evaluate-quality`

Owns subject-evaluation judgment: read the resolved target source, perform each
assessment, decide when evidence is sufficient, rate against criteria, and record
findings. It does not author or repair the model.

### `improve-quality-md`

Diagnose phase:

- uses the built-in meta-model at
  `internal/diagnostics/quality-model/QUALITY-META-MODEL.md`;
- evaluates the user's `QUALITY.md` as the subject;
- inspects whether Targets, factors, requirements, rating criteria, and body
  sections form a useful model.

Improve phase:

- proposes and applies fixes to the target tree and Markdown body;
- re-runs `lint`;
- re-runs diagnosis until the recorded issues are resolved or deliberately
  accepted.

## Relationship To Other Specs

- [`cli.md`](./cli.md) - resource command surface.
- [`cli-init.md`](./cli-init.md) - operates on the recursive Target type.
- Forthcoming: an evaluation-lifecycle spec (lifecycle and payload contract) and
  a structural-validator (`lint`) spec.
- [`../SPECIFICATION.md`](../SPECIFICATION.md) - the format: targets, scoped
  factors, requirements, assessment/finding/result, and rating criteria.

## Open Questions

- Shared judging methodology: evidence standards, saturation, and rigor levels.
- Exact meta-model wiring and how users replace or extend it, if ever.
- Improve-phase edit protocol.
- How much setup infers from the codebase versus asking the user.
- Skill packaging and repository home.
