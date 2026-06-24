# Docs Update Log

## 2026-06-24

- **Revision**: Updated
  [Modeling quality across domains](guides/model-quality-across-domains.md) for
  [0089 - Agent-harness modeling guidance](../changes/archive/0089-agent-harness-modeling-guidance.md).
  The guide now distinguishes explicit modeling guidance for recurring
  use-context constituents (the agent harness and QUALITY.md self-check) from
  modeled-domain defaults, defines the served-domain guardrail, and gives a
  good/avoid pair for harness requirements that stay domain-neutral.

- **Revision**: Updated README/front-door domain positioning for
  [0088 - Domain-agnostic corpus alignment](../changes/archive/0088-domain-agnostic-corpus-alignment.md).
  The README opening now presents QUALITY.md's modeled domain as broad (software,
  documentation, data, services, operations, and other tended entities) while keeping
  the `/quality` agent skill as the primary experience. The Agent Harnessability
  section now reads as an illustrative, earned factor family for
  agent-collaborated entities, not a built-in default for every QUALITY.md, and the
  Example QUALITY.md section links to the domain-agnostic guide's worked
  non-software example. Also corrected this log's 0083 description so it no longer
  names an obsolete Anthropic sample.

- **Creation**: Added
  [Modeling quality across domains](guides/model-quality-across-domains.md), a
  contributor-doctrine guide for keeping QUALITY.md examples quality-domain
  agnostic, for
  [0083 - Quality-domain agnosticism guide](../changes/0083-quality-domain-agnosticism.md).
  The guide defines the stress axes (source materiality, assessment oracle —
  including the ISO/IEC 25010 internal/external/quality-in-use split — constituency,
  stakes), grounds the quality-context axis in ISO lineage
  (product/data/requirements/process/service quality) while keeping ISO out of the
  earned factor axis, names a canonical set of secondary knowledge-work domains
  with illustrative factor families, adds range-finder illustrations (legal,
  budgeting, personal productivity, devotional/religious), carries a small
  illustrative catalog of quality contexts spanning professional knowledge work
  and the household (budget, meal planning, caregiving, upkeep, records, formation),
  cross-references the authoring constituent generator, and carries one full worked
  non-software (documentation) example. Linked it from the
  [guides index](guides/index.md).
- **Creation**: Added
  [Designing agent-mediated UX](guides/agent-mediated-ux.md), a contributor
  guide for workflows users experience through AI assistants or coding agents.
  The guide defines agent-mediated UX, gives formatting and emphasis rules,
  requires the primary question or call to action to stand out, and covers
  progress, discovery questions, checkpoints, decision gates, closeout, emoji,
  and tone. Linked it from the [guides index](guides/index.md).

## 2026-06-23

- **Removal**: Deleted the `Use the /quality skill` how-to guide
  (`guides/use-quality-skill.md`) and its index entry. The guide was the lone
  end-user doc in a contributor-facing folder and largely duplicated the
  README's install and usage sections; the skill itself and the README now
  carry that operational guidance.

- **Revision**: Documented the workflow feedback log in
  [Use the /quality skill](guides/use-quality-skill.md) for
  [0066 - Setup feedback log](../changes/archive/0066-setup-feedback-log.md):
  setup may record a hand-authored `.quality/logs/<timestamp>-setup-feedback-log.md`
  about the run experience, distinct from the quality log and evaluation
  `debug-log.md`, recorded locally and never transmitted, with secrets and
  prompt-injection text excluded and sensitive context sanitized — shareable by
  an explicit user action.

- **Revision**: Clarified agentic use context in [Use the /quality skill](guides/use-quality-skill.md)
  and [Install QUALITY.md](../install.md). The docs now preserve the
  agent-first experience while keeping modeled quality domains broad.

- **Revision**: Updated [Use the /quality skill](guides/use-quality-skill.md)
  for [0063 - Contextual setup flow](../changes/0063-contextual-setup-flow.md).
  Public guidance now describes setup as a context-informed `QUALITY.md` authoring
  flow that writes only the model, validates readiness, and leaves evaluation,
  quality-log writes, recommendation handoff, and recurring-review automation to
  follow-on workflows.

- **Revision**: Updated [Use the /quality skill](guides/use-quality-skill.md)
  and the [guides index](guides/index.md) for
  [0062 - Remove wizard mode](../changes/0062-remove-wizard-mode.md). Public
  guidance now lists `setup`, bare read-only `/quality`, `evaluate`, and
  `update` without advertising wizard.

## 2026-06-22

- **Revision**: Updated [Use the /quality skill](guides/use-quality-skill.md)
  and the top-level [install guide](../install.md) for
  [0057 - Quality data directory](../changes/0057-quality-data-directory.md):
  evaluation runs now default under `.quality/evaluations/`, the quality log
  defaults under `.quality/log/`, and root `config` frontmatter can point to a
  non-default workspace config file.

- **Revision**: Clarified functional-spec granularity in
  [Writing functional specs](guides/write-functional-specs.md): child specs can
  be justified by durable behavioral components, not only 1:1 artifact
  contracts, while parent specs keep shared invariants. Removed artifact-spec
  naming guidance from the purpose-agnostic [Working with OKF](guides/work-with-okf.md)
  guide and updated [Working with change cases](guides/work-with-change-cases.md)
  to account for both artifact specs and behavioral component specs.
- **Revision**: Added a **To rename** subsection to the
  [`## Durable spec changes`](guides/write-functional-specs.md#durable-spec-changes)
  convention in [Writing functional specs](guides/write-functional-specs.md), so a
  durable spec rename (`old → new`) is accounted in one place instead of splitting
  across To add and To delete. Mirrored the reference in
  [Working with change cases](guides/work-with-change-cases.md).
- **Revision**: Documented the two-clock status convention in
  [Working with change cases](guides/work-with-change-cases.md): `status` names the
  code clock, so when durable specs and docs advance ahead of code the status block
  should name both clocks.

## 2026-06-21

- **Creation**: Added [Go style](guides/go-style.md), a contributor guide for
  judgment-based Go conventions that the deterministic check gate does not
  enforce. Linked it from the [guide index](guides/index.md) and
  [`AGENTS.md`](../AGENTS.md), clarified the Go source scope in
  [Designing Go packages](guides/design-go-packages.md), and updated the dogfood
  [`QUALITY.md`](../QUALITY.md) model to evaluate CLI maintainability against the
  new guide.
- **Revision**: Updated [Versioning](reference/versioning.md) to mirror the
  `v0.7.1` `/quality` skill metadata while keeping the compatible
  `qualitymd >=0.7.0 <0.8.0` range.
- **Revision**: Updated [Use the /quality skill](guides/use-quality-skill.md)
  to mention the process-only `debug-log.md` included in evaluation runs.
- **Revision**: Updated [Versioning](reference/versioning.md) to mirror the
  `v0.6.0` `/quality` skill metadata and `qualitymd >=0.6.0 <0.7.0`
  compatibility range.

## 2026-06-19

- **Revision**: Broadened the change-case "Affected specs & docs" framing to
  **Affected artifacts** across [Working with change cases](guides/work-with-change-cases.md),
  [Writing functional specs](guides/write-functional-specs.md), and
  [Writing design docs](guides/write-design-docs.md) (and the
  [changes index](../changes/index.md)). A Change Case now accounts for every
  artifact kind it touches — code, durable specs, durable docs, the bundled
  skill, scaffold — so "specs and docs" no longer pre-narrows the list and leaves
  code or skill files unaccounted; code stays gated to `In-Progress` while
  durable specs/docs remain editable in any phase. Also added a "find them by
  analysis, not recall" step: derive the list from a repo-wide sweep over changed
  names/symbols/paths (and the names renamed away from) across all bundles, follow
  inbound cross-references, and triage live artifacts from frozen historical ones.

- **Revision**: Updated [Versioning](reference/versioning.md),
  [Use the /quality skill](guides/use-quality-skill.md), and the top-level
  [install guide](../install.md) for `/quality upgrade` as the paired
  skill/CLI maintenance flow, including manual fallbacks and restart/reload
  guidance after skill upgrades.

- **Revision**: Updated [Versioning](reference/versioning.md),
  [Use the /quality skill](guides/use-quality-skill.md), and
  [Cut a release](guides/cut-a-release.md) to make
  `skills/quality/SKILL.md` metadata the source of truth for the `/quality`
  skill version and required `qualitymd` CLI range, with release notes as the
  mirror and installer enforcement deferred.

- **Revision**: Updated [Versioning](reference/versioning.md),
  [Use the /quality skill](guides/use-quality-skill.md), and
  [Cut a release](guides/cut-a-release.md) for structured CLI version metadata,
  explicit upgrade checks, and managed installer release verification.

- **Revision**: Worked through the remaining release-guide process items. The
  [Cut a release](guides/cut-a-release.md) guide now records the current support
  boundary: keep release preparation manual unless repeated mistakes justify a
  focused helper, and keep `/quality` skill compatibility in release notes and
  docs until the skill installer defines package metadata.

- **Revision**: Clarified [Cut a release](guides/cut-a-release.md) after the
  `v0.2.2` release. The guide now makes hosted `main` CI an explicit pre-tag
  gate, notes that pre-tag Goreleaser snapshots can mention the previous tag,
  treats trailing-newline-only release-note diffs as equivalent, and keeps
  release-prep manual unless repeated mistakes justify more mechanics.

- **Revision**: Adopted root [`CHANGELOG.md`](../CHANGELOG.md) as the canonical
  curated release-note source. Updated
  [Cut a release](guides/cut-a-release.md) to remove the temporary changelog
  adoption language, require release prep to update the changelog, and keep only
  the remaining automation and skill-metadata support items open.

- **Automation**: Added release helpers for curated notes and pre-tag checks.
  `mise run release-notes -- vX.Y.Z` extracts the matching
  [`CHANGELOG.md`](../CHANGELOG.md) section, `mise run release-check -- vX.Y.Z`
  validates a prepared release commit and runs the release dry-run gates, and the
  release workflow now replaces Goreleaser's generated GitHub Release body with
  the curated changelog section.

- **Revision**: Made [Cut a release](guides/cut-a-release.md) the authoritative
  release runbook and narrowed [Contributing](../CONTRIBUTING.md) to contributor
  setup, local tasks, repo layout, and pointers to the release and versioning
  docs.

- **Creation**: Added [Cut a release](guides/cut-a-release.md), a `How-to Guide`
  for choosing a release version, checking the CLI/skill/specification versioned
  surfaces, preparing curated release notes, running dry runs, publishing a tag,
  verifying artifacts, handling failed releases, and tracking open process
  support items such as `CHANGELOG.md` adoption and release-note automation.

- **Creation**: Added [Versioning](reference/versioning.md), a `Reference`
  concept defining the separately versioned `qualitymd` CLI, `/quality` skill,
  and QUALITY.md specification surfaces, with CLI SemVer as the current
  skill-compatibility boundary.

- **Revision**: Updated [Use the /quality skill](guides/use-quality-skill.md)
  and the top-level [install guide](../install.md) to point CLI prerequisite
  checks at the versioning policy and the current evaluation command surface.

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
