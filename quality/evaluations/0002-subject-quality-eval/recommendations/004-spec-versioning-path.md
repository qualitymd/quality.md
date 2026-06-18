# 004 — Spec does not specify how the format versions forward

**Target:** format-spec → extensibility (also material to completeness)
**Related requirements:** *the format specifies its core and how it extends and
evolves*; *the format specification is complete*
**Severity:** medium — extensibility is the format spec's weakest factor.

## Gap

The spec names a minimal core and supports structural extension (nesting, body
sections, preserved unrecognized body content), but says nothing about how the
*format itself* evolves: there is no file schema-version field, no
forward/backward compatibility policy, and no statement of how an older reader
treats a newer file or how unrecognized frontmatter keys are handled (today they
lean toward making a file non-conforming). The document is versioned "0.1 —
Draft" but the file format is not.

## Evidence

- `SPECIFICATION.md:51,65,71` — minimal core specified (good).
- `SPECIFICATION.md:107-118,143-167,207` — structural extension specified
  (good).
- `SPECIFICATION.md:13,38,42` — unrecognized frontmatter keys lean toward
  rejection; no extension path for new keys.
- `SPECIFICATION.md:3` — document version only; no file-format versioning.

## Options

1. Add a short "Versioning and evolution" section: an optional file
   schema/format-version field, a forward-compatibility rule for unrecognized
   frontmatter keys (ignore-and-preserve vs. reject), and a deprecation/version
   bump policy.
2. Minimal: state a compatibility policy for unrecognized frontmatter keys and a
   version-bump convention, without yet adding a file-version field.
3. Defer explicitly: record this as a Known gap in the spec while pre-1.0, with
   the intended direction noted.

## Recommended

**Option 1** — a small dedicated section closes the "how it evolves" half of the
requirement and removes the ambiguity around unrecognized keys, which also lifts
the completeness requirement. If the team prefers to stay minimal pre-1.0,
Option 3 is an honest interim that should still name the intended versioning
direction.

## Done criterion

`the format specifies its core and how it extends and evolves` reaches at least
**target**: the spec states a forward-evolution/versioning path and how readers
treat unrecognized frontmatter content. Re-evaluate in a new numbered run.
