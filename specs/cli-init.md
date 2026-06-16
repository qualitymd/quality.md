# CLI: `init`

`init` scaffolds the first `QUALITY.md`.

```bash
qualitymd init
```

## Purpose

`init` is deterministic and offline. It writes a starter recursive Target-tree
model plus a Markdown body skeleton. It does not inspect the codebase and does
not call a model; tailoring the model belongs to the authoring skills and normal
edit/lint loop.

## Authoring Loop

```text
qualitymd init
qualitymd lint
evaluate-quality / improve-quality-md through skills
```

The scaffold is a starting point, not a finished quality model.

## What It Writes

```markdown
---
# QUALITY.md - quality model for <project>. See SPECIFICATION.md.
# The frontmatter root is the Model (the apex Target).
# `source` defaults to this file's directory recursively when omitted.
# Child `targets:` narrow or decompose the subject; position is lineage.
# Requirements use one `assessment`; ratings default to Outstanding/Target/Minimum/Unacceptable.
targets:
  source-code:
    source: ./src
    requirements:
      "<what this target must accomplish>":
        assessment: "<inline criteria, or ./path/to/standard.md>"
    factors:
      maintainability:
        description: "<what maintainability means for this target>"
        requirements:
          "<what maintainable code looks like here>":
            assessment: "<inline criteria, or ./path/to/standard.md>"
  docs: ./docs
---

# Quality model - <project>

## Overview
<!-- What this system or component is, who depends on it, and what "good" means here. -->

## Scope
<!-- What this model covers and deliberately leaves out. -->

## Needs
<!-- Stakeholder outcomes the requirements answer to. -->

## Risks
<!-- What goes wrong, and for whom, if a need is not met. -->

## Targets and factors
<!-- Mirror the target tree and explain the factors declared on each target. -->

### source-code
<!-- Why this target exists and why Maintainability is the right lens here. -->

### docs
<!-- What quality concerns belong to documentation, if any. -->

## Known gaps
<!-- In-scope concerns intentionally deferred, each with a reason. -->
```

The scaffold is valid under the revised schema and is intended to pass `lint`
after placeholders are filled with non-empty criteria.

## Optional Config

`--config` also creates:

```text
./.quality/
  config.yaml
```

Configuration is off by default.

## Interactive Behavior

In a TTY, `init` may prompt for project name and initial target/factor labels.
With `--json`, `--non-interactive`, or non-TTY stdin/stdout, it writes the fixed
placeholder scaffold without prompting. Prompts go to stderr.

## Flags And Exit Codes

`init` writes `./QUALITY.md` and never clobbers an existing file.

- `--config` - also scaffold `.quality/config.yaml`.

Exit codes:

- `0` - successful write.
- `2` - tool failure, including an existing `QUALITY.md` or write error.
