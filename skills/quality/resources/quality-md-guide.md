<!-- AGENTS: Only minor grammatical, spelling, and formatting changes may be made. Do not alter prose. -->

# QUALITY.md How-to Guide

Working with QUALITY.md files.

## The QUALITY.md file

The entire `QUALITY.md` file represents a single apex quality target (defined below).

## Rating Scale

## Targets

Targets are the entities subject to evaluation by quality requirements associated
with different quality factors for that target, identified by their source.

The target `source` identifies the entity/entities that will be subject to evaluation. In the typical case where `QUALITY.md` does not define a `source` property, it is assumed to be the project or workspace context defined by the contents of the same directory that the `QUALITY.md` file is present in and all subdirectories (`**/*`).

Targets may have child targets to facilitate segregating requirements or factors
that are distinctly associated with the child target.

### Selecting targets

- [ ] Prefer primary artifacts over secondary/supporting artifacts/entities.

### Defining target sources

### Specifying target requirements

- [ ] Keep target

## Requirements

### Related Factors

Requirements may reference multiple other factors within the scope of their
target. This is helpful when a single assessment spans multiple factors or has
secondary factors in addition to the factor it is listed under.

### Assessments

Requirement assessments may be articulated inline or reference a source of
assessment (such as docs, specs, checklists, etc.).

- [ ] Name each referenced assessment once, referencing related factors.
- [ ] Don't extract, summarize, or duplicate assessment content from specs, docs
      or other possible assessments.

## Best Practices

## When to update QUALITY.md

- [ ] Discoveries are made that inform the context or content of the quality evaluation
