# 003 — Spec describes invalidity but never shows invalid examples

**Target:** format-spec → verifiability
**Related requirement:** *the format's constructs are shown with valid and
invalid examples*
**Severity:** medium.

## Gap

`SPECIFICATION.md` illustrates constructs richly on the valid side (schema
blocks, a ratings-override example, a body example, the Appendix A report) and
describes invalidity in prose, but shows no worked invalid counter-examples a
reader can compare against. The requirement asks for both valid cases and
invalid counter-examples.

## Evidence

- `SPECIFICATION.md:48-63,191-201,225-257,323-419` — valid examples (good).
- `SPECIFICATION.md:38-42,40,185` — invalidity described in prose only
  (non-conforming shapes; missing/empty/list-valued assessment; null treated as
  absent).

## Options

1. Add short invalid counter-examples beside the existing valid blocks — e.g. a
   frontmatter with a list-valued `assessment`, a model with only one rating
   level, a target whose subtree has no requirements — each with a one-line "why
   this is invalid."
2. Add a single consolidated "Invalid examples" subsection collecting the common
   conformance failures.

## Recommended

**Option 1** — place each invalid counter-example next to the construct it
violates, so the contrast is local and teachable.

## Done criterion

`the format's constructs are shown with valid and invalid examples` reaches at
least **target**: the principal constructs each show at least one invalid
counter-example with a stated reason. Re-evaluate in a new numbered run.
