# CLI: the evaluation lifecycle (`evaluation` / `result`)

This is the detail doc for the deterministic `evaluation` and `result`
resources: data model, on-disk layout, command behavior, report rollup, and the
authoritative CLI-to-skill payloads. The CLI never calls a model; skills judge
requirements and write verdicts through these resources.

## Resources

| Resource | What it is |
| --- | --- |
| `evaluation` (`eval`) | A living run for one `(model, run target)` pair. The run target is the CLI's `--target <path>` selection, distinct from schema `targets:` nodes. |
| `result` | One requirement result inside that run, addressed by a target-tree address. |

## Boundary

The CLI performs mechanics: parse the model, resolve the target tree, expand
`source` globs, persist runs, record ratings, detect staleness, and roll up
reports. Judgment belongs to skills: they read `result show`, perform the
requirement's `assessment`, produce a finding, choose the rating level whose
criterion the finding meets, and call `result set`.

## Data Model

A run carries:

- identity: model path, CLI run target path, and deterministic slug;
- selected requirement set, each with a target-tree address and result state;
- derived status and rollup.

A result carries:

- **address** - target-tree address of the requirement:
  `targetPath`, optional `factorPath`, and `requirement`;
- state: `pending`, `recorded`, `skipped`, or `stale`;
- rating level and structured evidence once recorded;
- staleness provenance for the resolved definition and source contents.

The result is the diffable artifact a reviewer reads.

## States

Run status is derived:

| Status | Meaning |
| --- | --- |
| `open` | At least one result is `pending` or `stale`. |
| `complete` | Every result is `recorded` or `skipped`. |
| `archived` | Frozen copy created by `evaluation archive --as <name>`. |

Result states:

| State | Meaning |
| --- | --- |
| `pending` | Enumerated but not assessed. |
| `recorded` | Carries rating and evidence. |
| `skipped` | Deliberately not assessed, with reason. |
| `stale` | Was recorded, but hash inputs changed. |

On re-run, only stale results return to pending. A low rating is still a recorded
result; gating happens in `evaluation report`.

## On-Disk Layout

```text
.quality/
  config.yaml
  evaluations/
    <slug>/
      evaluation.json
      results/
        <result-id>.json
      report.md
      .run/
        meta.json
    archive/
      <name>/
```

`evaluation.json`, `results/*.json`, and `report.md` are deterministic,
reviewable artifacts. Volatile metadata such as timestamps and durations lives in
`.run/`.

## Commands

### `evaluation create [--model <path>] [--target <path>] [--from <id>]`

Create or re-enter the living run. It parses the model, resolves the recursive
target tree, applies the CLI run target selection, and enumerates every selected
direct and lensed requirement as a result. On an existing run it re-hashes
recorded results and marks changed ones `stale`.

### `evaluation list`

List living and archived runs.

### `evaluation show [<id>] [--json]`

Show run identity, status, selected result set, and rollup.

### `evaluation report [<id>] [--fail-on <level>] [--json]`

Render `report.md` and the JSON rollup. `--fail-on` exits non-zero when the
overall rating lands at or below the named level. The default is `unacceptable`.

### `evaluation archive [<id>] --as <name>`

Snapshot a run to `.quality/evaluations/archive/<name>/`.

### `evaluation delete <id>`

Discard a living or archived run.

### `result list [--status pending,stale,...] [--json]`

List results by state. There is no `next` cursor; skills choose work order.

### `result show <address> [--json]`

Emit the fully resolved requirement payload a skill consumes.

### `result set <address> --rating <level> --evidence ...`

Record the skill's verdict. The rating must be a level in the resolved scale.

### `result skip <address> --reason ...`

Mark a result skipped with a reason.

### `result reset <address>`

Return a result to pending.

## Addressing

A requirement address is a stable target-tree locator:

```json
{
  "targetPath": ["api", "handlers"],
  "factorPath": ["security"],
  "requirement": "malformed input is rejected"
}
```

`targetPath` is empty for the apex target. `factorPath` is empty for direct
unlensed requirements. A string form may use dotted segments with quoted map keys,
for example:

```text
targets.api.targets.handlers.factors.security.requirements."malformed input is rejected"
```

This resolves the former open question about federated addressing: a federated
model is grafted into the same target tree, so its model path plus target-tree
address uniquely identifies a result.

## Interface Payloads

These payloads are the authoritative CLI-to-skill contract. Field names are still
provisional, but the information is required.

### `result show` Output

```json
{
  "schemaVersion": 1,
  "model": "QUALITY.md",
  "address": {
    "targetPath": ["api"],
    "factorPath": ["security"],
    "requirement": "malformed input is rejected"
  },
  "target": {
    "name": "api",
    "source": {
      "patterns": ["./internal/api/**/*.go"],
      "files": [
        { "path": "internal/api/orders.go" },
        { "path": "internal/api/validation.go" }
      ]
    }
  },
  "factor": {
    "primary": "security",
    "secondary": ["reliability"],
    "visibleFactors": ["security", "reliability"]
  },
  "assessment": "Malformed requests are rejected at the boundary with a clear error and no partial side effects.",
  "ratings": {
    "order": "bestToWorst",
    "levels": [
      { "level": "outstanding", "criterion": "Exceeds the requirement; satisfies it with margin to spare." },
      { "level": "target", "criterion": "Satisfies the requirement." },
      { "level": "minimum", "criterion": "Falls short of the goal but stays at the acceptable floor." },
      { "level": "unacceptable", "criterion": "Falls below the acceptable floor." }
    ]
  },
  "ratingOverrides": null,
  "state": "pending",
  "sufficiency": "Rate after inspecting every file in the resolved source manifest that can handle input."
}
```

The skill receives resolved assessment text, source patterns and file manifest,
the in-scope factor context, scale criteria, current state, and sufficiency
guidance. File contents may be inlined by option later, but paths are enough for
the skill to load evidence itself.

### `result set` Input

```json
{
  "schemaVersion": 1,
  "address": {
    "targetPath": ["api"],
    "factorPath": ["security"],
    "requirement": "malformed input is rejected"
  },
  "rating": "minimum",
  "finding": {
    "summary": "JSON bodies are validated, but query parameters reach handlers unchecked.",
    "items": [
      { "location": "internal/api/orders.go:42", "note": "validates JSON body before use" },
      { "location": "internal/api/validation.go:17", "note": "query parameters are parsed without validation" }
    ]
  },
  "rationale": "The boundary has partial validation, so the result falls below target but above unacceptable."
}
```

On disk, the CLI serializes this shape deterministically and keeps volatile
metadata outside the result file.

### Staleness Hash

A recorded result becomes stale when any of these inputs change:

1. **Resolved requirement definition** - model path, target-tree address,
   requirement statement, factor placement, secondary factors, resolved
   `assessment` text, active rating scale criteria, and per-requirement criteria.
2. **Resolved source selection** - `source` patterns and expanded manifest for
   the target where the requirement is assessed.
3. **Resolved source contents** - bytes or configured content hash for files in
   the manifest.

Canonical serialization of the hash input is an implementation detail, but it
must be deterministic. Whether the content contribution is raw bytes or a coarser
revision signal is still open.

## Carry-Forward

`evaluation create --from <id>` copies still-valid recorded results from another
run when the target-tree address exists in both runs and the staleness hash still
matches. Otherwise the destination result starts pending.

## Report Rollup

`evaluation report` rolls recorded ratings through the target tree:

1. Requirements roll up to their primary factor when lensed.
2. Requirements with secondary `factors` also surface under those factor views
   without duplicating the recorded result.
3. Direct requirements roll up to their target.
4. Factors roll up to their declaring target and descendant targets that inherit
   them.
5. Targets roll up through parent targets to the apex overall verdict.

The default rollup is worst-wins. `skipped` results are reported as coverage gaps
and excluded from the rating. `pending` or `stale` results make the rollup
incomplete; the report states incompleteness rather than scoring around it.

`report.md` includes definition (model, CLI run target, selected target nodes),
overall verdict, per-target rollup, per-factor rollup, per-requirement findings,
secondary-factor appearances, and coverage gaps.

## Exit Codes

| Code | Meaning |
| --- | --- |
| `0` | Command ran and any gate passed. `result set` exits `0` even for a low rating. |
| `1` | `evaluation report --fail-on` gate tripped. |
| `2` | Tool failure: bad flags, unreadable model, parse failure, or internal error. |

## Open Questions

- Exact staleness serialization and whether source contents hash raw bytes or a
  coarser content signal.
- Whether `.run/` is committed or gitignored.
- Whether rollup remains fixed worst-wins or gains configured strategies.
- Final field names for the interface payloads.
