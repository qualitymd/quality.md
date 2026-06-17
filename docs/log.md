# Docs Update Log

## 2026-06-17

- **Creation**: Added the [Proposing a change](guides/propose-a-change.md)
  `How-to Guide` covering the `changes/` workflow — numbering, the
  Change-plus-spec-plus-design shape, the status lifecycle, and archiving.

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
