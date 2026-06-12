# QUALITY.md Format

QUALITY.md is a plain text representation of a quality model. It can be used to specify and evaluate the quality requirements for a software system or component.

A QUALITY.md file contains two parts: YAML frontmatter with the structured quality model and the markdown body.

## Quality Model
The quality model is embedded in the YAML front matter at the beginning of the file. The front matter block must begin with a line containing exactly --- and end with a line containing exactly ---. The YAML content between these delimiters is parsed according to the schema defined below.

Example:
```yaml

factors:
  security:
    requirements:
      "no secrets committed to the repository":
        rules: >
          No hard-coded credentials, API keys, private keys, or tokens
          appear in source, config, or fixtures. Secrets are loaded from
          environment variables or a secrets manager at runtime.
        rating:
          pass: "no literal secrets found; all secrets resolved at runtime"
          fail: "one or more credentials are committed to the repository"
      "meets application security standards":
        rules: "./standards/appsec-checklist.md"
        rating:
          pass: "every applicable item in the checklist is satisfied"
          fail: "any applicable checklist item is unmet"
  maintainability:
    factors:
      reusability:
        requirements:
          "shared domain types come from the common package":
            rules: >
              Domain models exchanged across modules are imported from
              @acme/common rather than redefined locally. Module-local
              duplicates of a shared type are not allowed.
            rating:
              pass: "all cross-module types are imported from @acme/common"
              fail: "a shared type is redefined instead of imported"
      testability:
        requirements:
          "unit tests pass":
            bash: "pnpm test:unit"
          "critical paths meet 80% line coverage":
            bash: "pnpm test:coverage --min 80"

```

### Schema

```yaml
ratings:                          # optional; defaults to pass / fail
  <level-name>:                   # listed best to worst
    displayName: <string>
    description: <string>         # optional
factors:
  <factor-name>:                  # a factor has requirements OR nested sub-factors
    requirements:
      <requirement-name>:
        rules: <text | path>      # one evaluator per requirement
        bash: <command>
        rating:                   # optional; keys are levels from `ratings`
          <level-name>: <criteria>
    factors:                      # sub-factors, nested to any depth
      <factor-name>: <Factor>
```

The frontmatter is a single mapping. `factors` is required; `ratings` is an
optional sibling that customizes the rating scale.

**Factors.** `factors` is a map of factor name → factor. A factor is a named
quality attribute (for example `security` or `maintainability`). A factor
either lists `requirements` directly, or nests sub-`factors` to decompose into
finer attributes — exactly one of the two, not both. Sub-factors nest to any
depth, forming a tree whose leaves carry requirements. Names are the map keys
and must be unique among siblings.

**Requirements.** `requirements` is a map of requirement name → requirement,
where the name states the expectation (for example `"unit tests pass"`). Each
requirement declares exactly one evaluator:

- `rules: <text | path>` — an *inferential* check, judged by a reviewer or
  model. Either inline text describing the conditions a target must satisfy,
  or a path (relative to the QUALITY.md file) to a Markdown guide.
- `bash: <command>` — a *computational* check. A shell command whose exit
  status is the verdict: zero passes, non-zero fails.

**Rating.** A requirement may also carry an optional `rating` whose keys are
levels from the rating scale, each mapped to the observable condition under
which the requirement earns that level. A rating documents the criteria behind
a `rules`-based check, so it normally pairs with `rules`; a `bash` requirement
derives its verdict from the exit status (the best level on a zero exit, the
worst otherwise) and needs none. Write each criterion so it follows from the
requirement's `rules` rather than restating the level name.

**Rating scale.** The optional top-level `ratings` map defines the scale
shared by every requirement. Each entry is a level name (the key used in a
requirement's `rating`) carrying a `displayName` and an optional
`description`. List levels best to worst; the order defines their ranking. When
`ratings` is omitted, the scale defaults to `pass` then `fail`. A custom scale
might read:

```yaml
ratings:
  A: { displayName: "Excellent" }
  B: { displayName: "Good" }
  C: { displayName: "Acceptable" }
  D: { displayName: "Poor" }
  E: { displayName: "Unacceptable", description: "Fails minimum bar" }
```
