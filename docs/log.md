# Docs Update Log

## 2026-06-18

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
  [Working with changes](guides/work-with-changes.md). Durable specs and docs
  **MAY** now be edited at any time, with or without a change (the lifecycle gate
  governs **code**, not durable artifacts); a change **MUST** list every durable
  spec or doc its work impacts in **Affected specs & docs** and **SHOULD** suggest
  any new durable specs worth creating, and **SHOULD** bring the listed docs into
  sync before **In-Review**. Softened the absorb-the-*why* step from a MUST gated
  on In-Review to a SHOULD encouraged whenever a change updates a spec, and made
  reading [Writing functional specs](guides/write-functional-specs.md) a required
  precondition for creating or updating any change spec. Mirrored the policy in
  the [changes index](../changes/index.md).
- **Revision**: Taught three how-to guides to keep durable rationale in the
  spec, for change
  [0025 — Durable spec rationale](../changes/archive/0025-durable-spec-rationale.md).
  [Writing functional specs](guides/write-functional-specs.md) gains a
  **Background / Motivation** shape entry and a per-requirement `Rationale:`
  annotation convention (with form, the annotate-when litmus, and a say-it-once
  rule); its "Motivation in asides" convention is rewritten as a two-whys split,
  and its rationale smells now target rationale that *buries* the rule and
  rationale said twice. [Working with changes](guides/work-with-changes.md) now
  requires absorbing a landing change's enduring *why* — its motivation and the
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
  (`propose-a-change.md`) to [Working with changes](guides/work-with-changes.md)
  (`work-with-changes.md`), because the old title named only the first step while
  the guide covers the whole lifecycle (create → spec → design → implement →
  done → archive). Updated the title, heading, and intro; the
  [guides index](guides/index.md); and the `AGENTS.md` guides table. Earlier log
  entries keep their original wording; their links now resolve to the new path.
- **Revision**: Hardened the lifecycle gate in
  [Proposing a change](guides/work-with-changes.md): added a phase-authorization
  rule stating that each phase modifies only what it authorizes — **Draft** the
  `spec.md`, **Design** the `design.md` (each MAY add supporting files within the
  change folder), **In-Progress** the implementation — and that implementation
  does not begin until the change is In-Progress. Framed it as a whitelist (what
  does this phase permit?), not a blacklist scoped to the **Affected specs &
  docs** list, with the change folder `changes/NNNN-<slug>/` as the hard boundary
  before In-Progress: nothing outside it is touched until then.
- **Creation**: Added the [Proposing a change](guides/work-with-changes.md)
  `How-to Guide` covering the `changes/` workflow — numbering, the
  Change-plus-spec-plus-design shape, the status lifecycle, and archiving.
- **Revision**: Extended [Proposing a change](guides/work-with-changes.md) to
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
