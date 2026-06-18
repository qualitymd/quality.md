# Docs Update Log

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
