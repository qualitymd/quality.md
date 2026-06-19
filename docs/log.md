# Docs Update Log

## 2026-06-19

- **Creation**: Added [Cut a release](guides/cut-a-release.md), a `How-to Guide`
  for choosing a release version, checking the CLI/skill/specification versioned
  surfaces, preparing curated release notes, running dry runs, publishing a tag,
  verifying artifacts, handling failed releases, and tracking open process
  support items such as `CHANGELOG.md` adoption and release-note automation.

- **Revision**: Updated [Writing functional specs](guides/write-functional-specs.md)
  to use BCP 14 keywords sparingly: only where the keyword changes conformance
  meaning. Cleaned related workflow guides so they use ordinary prose for
  contribution-process advice rather than uppercase requirement words.

- **Creation**: Added [RFC 8174](reference/rfc8174.md) as a `Reference`
  concept beside [RFC 2119](reference/rfc2119.md), documenting the BCP 14
  uppercase-only clarification used by `SPECIFICATION.md`, and listed it in the
  [reference index](reference/index.md).

## 2026-06-18

- **Revision**: Made change cases account for durable-**spec** deltas in the
  functional spec, not the design doc.
  [Working with change cases](guides/work-with-change-cases.md) now requires a
  change-case `spec.md` to carry a `## Durable spec changes` section and splits
  the labor: the parent's **Affected specs & docs** is the index of every
  artifact touched, while the functional spec carries the substance of each
  durable-spec change (the `specs/` bundle and `SPECIFICATION.md`).
  [Writing functional specs](guides/write-functional-specs.md) defines that
  section's form — **To add** / **To modify** / **To delete**, each required to
  read a list or an explicit `None`, with entries pointing to the driving
  requirement rather than restating it.
  [Writing design docs](guides/write-design-docs.md) records that a design doc
  does not carry its own durable-spec/doc list; durable impact is the spec's job
  and the design covers only *how*.

- **Rename**: Renamed the change workflow guide to
  [Working with change cases](guides/work-with-change-cases.md), updated the
  `changes/` schema term to `Change Case`, and narrowed the `AGENTS.md` trigger:
  routine prompted edits do not require a Change Case.

- **Revision**: Updated [Use the /quality skill](guides/use-quality-skill.md)
  after removing the bundled model workflow. The CLI prerequisite check no
  longer includes `qualitymd models`, and the runnable examples now show
  subject evaluation rather than model-altitude evaluation.

- **Revision**: Expanded [Writing functional specs](guides/write-functional-specs.md)
  with lessons adapted from Joel Spolsky's *Painless Functional Specifications*
  (parts 2 and 4). Added a **Scenario / use case** shape entry — the self-contained
  case a spec solves, distinguished from the larger-process *why* that belongs in
  **Background / Motivation** — and broadened **Scope** to name **non-goals** (out
  of scope by design) alongside **deferred**. New conventions: *an unspecified case
  is a decision delegated* (spec the divergent cases or hand the decision silently
  to the implementer), *sections are a palette, not a checklist* (against
  template-bloat), and *show, don't only tell* (a concrete example pins a contract).
  Reframed the **Shape** intro to order a spec for a reader. Did not adopt Joel's
  humour or vivid personas — they fight the brevity ethos and the
  deterministic-contract purpose.
- **Revision**: Decoupled durable-spec/doc edits from the change lifecycle in
  [Working with change cases](guides/work-with-change-cases.md). Durable specs
  and docs **MAY** now be edited at any time, with or without a Change Case (the
  lifecycle gate governs **code**, not durable artifacts); a Change Case
  **MUST** list every durable spec or doc its work impacts in **Affected specs &
  docs** and **SHOULD** suggest any new durable specs worth creating, and
  **SHOULD** bring the listed docs into sync before **In-Review**. Softened the
  absorb-the-*why* step from a MUST gated on In-Review to a SHOULD encouraged
  whenever a Change Case updates a spec, and made reading
  [Writing functional specs](guides/write-functional-specs.md) a required
  precondition for creating or updating any change case spec. Mirrored the policy
  in the [changes index](../changes/index.md).
- **Revision**: Taught three how-to guides to keep durable rationale in the
  spec, for change
  [0025 — Durable spec rationale](../changes/archive/0025-durable-spec-rationale.md).
  [Writing functional specs](guides/write-functional-specs.md) gains a
  **Background / Motivation** shape entry and a per-requirement `Rationale:`
  annotation convention (with form, the annotate-when litmus, and a say-it-once
  rule); its "Motivation in asides" convention is rewritten as a two-whys split,
  and its rationale smells now target rationale that *buries* the rule and
  rationale said twice. [Working with change cases](guides/work-with-change-cases.md) now
  requires absorbing a landing Change Case's enduring *why* — its motivation and the
  design doc's rationale — into the durable spec's Background and annotations,
  tied to the **Before setting In-Review** gate.
  [Writing design docs](guides/write-design-docs.md) records that the design
  doc's rationale is promoted into the spec on landing while the doc remains the
  archived record of alternatives and trade-offs.

- **Creation**: Added [Designing CLI interfaces](guides/cli-design.md), a
  `How-to Guide` for designing a `qualitymd` command's arguments, flags, output,
  and errors. It carries the design principles and a per-aspect checklist
  (philosophy, basics, help, output, arguments and flags, subcommands, errors and
  exit codes, robustness, determinism, next actions, configuration, environment,
  naming, distribution, analytics) adapted from clig.dev and cross-linked to the
  [CLI spec](../specs/cli.md). Listed it in the [guides index](guides/index.md)
  and the `AGENTS.md` guides table.
- **Revision**: Updated [Use the /quality skill](guides/use-quality-skill.md)
  so the CLI prerequisite check includes the evaluation commands the skill now
  uses for run creation, record writes, status checks, and report rendering.

## 2026-06-17

- **Creation**: Added [Use the /quality skill](guides/use-quality-skill.md), a
  how-to guide for installing the skill with `npx skills add
  qualitymd/quality.md`, verifying the `qualitymd` CLI prerequisite, running
  setup/wizard/evaluation modes, and configuring `.quality/config.yaml`
  `evaluationDir`.

- **Rename**: Renamed the changes-workflow guide from *Proposing a change*
  (`propose-a-change.md`) to [Working with change cases](guides/work-with-change-cases.md)
  (`work-with-change-cases.md`), because the old title named only the first step while
  the guide covers the whole lifecycle (create → spec → design → implement →
  done → archive). Updated the title, heading, and intro; the
  [guides index](guides/index.md); and the `AGENTS.md` guides table. Earlier log
  entries keep their original wording; their links now resolve to the new path.
- **Revision**: Hardened the lifecycle gate in
  [Proposing a change](guides/work-with-change-cases.md): added a phase-authorization
  rule stating that each phase modifies only what it authorizes — **Draft** the
  `spec.md`, **Design** the `design.md` (each MAY add supporting files within the
  change folder), **In-Progress** the implementation — and that implementation
  does not begin until the change is In-Progress. Framed it as a whitelist (what
  does this phase permit?), not a blacklist scoped to the **Affected specs &
  docs** list, with the change folder `changes/NNNN-<slug>/` as the hard boundary
  before In-Progress: nothing outside it is touched until then.
- **Creation**: Added the [Proposing a change](guides/work-with-change-cases.md)
  `How-to Guide` covering the `changes/` workflow — numbering, the
  Change Case plus spec/design shape, the status lifecycle, and archiving.
- **Revision**: Extended [Proposing a change](guides/work-with-change-cases.md) to
  account for the durable specs and docs a change touches — a new
  "Account for the specs and docs it touches" section, the **Affected specs &
  docs** step in the parent concept, and a Done gate requiring the listed
  enduring specs/docs be brought into sync before archiving.

## 2026-06-16

- **Conversion**: Restructured `docs/` into a single OKF knowledge bundle
  organized by the four [Diátaxis](https://diataxis.fr/) modes. Added the bundle
  [index](index.md) (`okf_version: "0.1"`), [`schema.md`](schema.md) registering
  the mode types (`Tutorial`, `How-to Guide`, `Reference`, `Explanation`), and
  listing-only indexes for [`guides/`](guides/index.md) and
  [`reference/`](reference/index.md).
- **Move**: Relocated the editing guides into [`guides/`](guides/) as
  `How-to Guide` concepts — [Working with OKF](guides/work-with-okf.md),
  [Writing functional specs](guides/write-functional-specs.md), and
  [Writing design docs](guides/write-design-docs.md) — adding frontmatter and
  fixing cross-links.
- **Conversion**: Turned `reference/rfc2119.txt` into the
  [RFC 2119](reference/rfc2119.md) `Reference` concept with OKF frontmatter, the
  RFC text preserved in a fenced block, and a citation to the canonical RFC
  Editor source.
