# QUALITY.md Format

QUALITY.md is a plain text representation of a quality model. It can be used to specify and evaluate the quality requirements for a software system or component.

A QUALITY.md file contains two parts: YAML frontmatter with the structured quality model and the markdown body.

## Quality Model

The quality model is embedded in the YAML front matter at the beginning of the file. The front matter block must begin with a line containing exactly --- and end with a line containing exactly ---. The YAML content between these delimiters is parsed according to the schema defined below.

The model keeps three concerns separate:

- **Requirement** — *what* must be true. A requirement is self-contained: a name, an optional target, and one assessment. It holds nothing about scoring, so it can be lifted out of QUALITY.md and reused on its own — for example, referenced from an agent skill.
- **Evaluation** — *how* a requirement is assessed. Either inferential (`prompt`) or computational (`bash`); the key you choose is the declaration of method.
- **Rating** — *how the result is classified*. A single, requirement-agnostic scale (`ratings`) names the levels an evaluation can land on.

Example:
```yaml
factors:
  functionality:
    requirements:
      "CLI matches its specification":
        target: "./internal/cli"
        prompt: "./specs/cli.md"
  security:
    requirements:
      "no secrets committed to the repository":
        prompt: >
          No hard-coded credentials, API keys, private keys, or tokens
          appear in source, config, or fixtures. Secrets are loaded from
          environment variables or a secrets manager at runtime.
      "meets application security standards":
        prompt: "./standards/appsec-checklist.md"
  maintainability:
    factors:
      reusability:
        requirements:
          "shared domain types come from the common package":
            target: "./src/**/*.ts"
            prompt: >
              Domain models exchanged across modules are imported from
              @acme/common rather than redefined locally. Module-local
              duplicates of a shared type are not allowed.
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
        target: <path | glob>     # optional; the artifact under evaluation
        # exactly one assessment:
        prompt: <text | path>     # inferential (judged by a model/reviewer)
        bash: <command>           # computational (shell exit status)
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
where the name states the expectation (for example `"unit tests pass"`). A
requirement names an optional `target` and declares exactly one assessment. It
holds nothing about rating levels, which keeps it portable — a requirement, or
the file its assessment points to, can be referenced and evaluated on its own,
outside QUALITY.md.

The optional `target` is a path or glob pattern (relative to the QUALITY.md
file) identifying the file or directory the requirement is evaluated against. A
glob (for example `./src/**/*.ts`) selects a set of files. When `target` is
omitted, the requirement applies to the QUALITY.md file's directory.

The assessment is exactly one of:

- `prompt: <text | path>` — an *inferential* assessment, judged by a model (or
  a human reviewer). Either inline text stating what the target must satisfy,
  or a path (relative to the QUALITY.md file) to a Markdown document holding the
  same — a checklist, a style guide, a specification to conform to. A path is
  the seam for reuse: the standard lives in its own file that both QUALITY.md
  and, say, an agent skill can load.
- `bash: <command>` — a *computational* assessment. A shell command whose exit
  status is the verdict.

The key is the declaration of method: `prompt` says "assess this by judgment,"
`bash` says "assess this by running a command." There is no separate field
naming the criteria — the criteria *are* the prompt or the command.

**Rating.** Evaluating a requirement produces a rating: the level its result
lands on, drawn from the `ratings` scale below. A `prompt` is judged against the
scale — the target earns the best level when it fully satisfies the prompt, and
lower levels as it falls short. A `bash` assessment earns the best level on a
zero exit and the worst otherwise. The rating is produced by evaluation; it is
never declared on the requirement. When a requirement needs graded expectations
(for example coverage bands), state the gradations in the `prompt` prose itself
— never as a map keyed on level names — so the requirement stays self-contained
and scale-independent.

**Rating scale.** The optional top-level `ratings` map defines the single scale
shared by every requirement. Each entry is a level name carrying a `displayName`
and an optional `description`. Write the descriptions generically — in terms of
how fully an evaluation meets its assessment — so the one scale applies to every
requirement. List levels best to worst; the order defines their ranking. When
`ratings` is omitted, the scale defaults to `pass` then `fail`. A custom scale
might read:

```yaml
ratings:
  A: { displayName: "Excellent",    description: "Fully satisfies the assessment; no gaps" }
  B: { displayName: "Good" }
  C: { displayName: "Acceptable",   description: "Satisfies the core of the assessment; minor gaps" }
  D: { displayName: "Poor" }
  E: { displayName: "Unacceptable", description: "Does not satisfy the assessment" }
```
