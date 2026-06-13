# QUALITY.md Format

QUALITY.md is a plain text representation of a quality model. It can be used to specify and evaluate the quality requirements for a software system or component.

A QUALITY.md file contains two parts: YAML frontmatter with the structured quality model and the markdown body.

## Quality Model

The quality model is embedded in the YAML front matter at the beginning of the file. The front matter block must begin with a line containing exactly --- and end with a line containing exactly ---. The YAML content between these delimiters is parsed according to the schema defined below.

The model keeps three concerns separate:

- **Requirement** — *what* must be true. A requirement is self-contained: a name, an optional target, and a single assessment. It holds nothing about scoring, so it can be lifted out of QUALITY.md and reused on its own — for example, referenced from an agent skill.
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
  <factor-name>:                  # a factor has requirements, sub-factors, or both
    requirements:
      <requirement-name>:
        target: <path | glob>     # optional; the artifact under evaluation
        # exactly one assessment — a single prompt OR a single bash command:
        prompt: <text | path>     # inferential (judged by a model/reviewer)
        bash: <command>           # computational (shell exit status)
    factors:                      # sub-factors, nested to any depth
      <factor-name>: <Factor>
```

The frontmatter is a single mapping. `factors` is required; `ratings` is an
optional sibling that customizes the rating scale.

**Factors.** `factors` is a map of factor name → factor. A factor is a named
quality attribute (for example `security` or `maintainability`). A factor may
list `requirements`, nest sub-`factors` to decompose into finer attributes, or
both — `requirements` assess the attribute directly while sub-`factors` break it
down further; a factor must carry at least one of the two. Sub-factors nest to
any depth, forming a tree, and requirements may hang off any factor in it, not
only the leaves. Names are the map keys and must be unique among siblings.

**Naming and describing factors.** The format does not constrain how you name or
describe a factor, but a few conventions keep a model legible (these are
recommendations, not rules):

- **A factor names a quality attribute, not a part of the system.** It captures
  one dimension of what *good* means — reliability, security, maintainability —
  not a component (*the API*, *the database*) or an activity (*testing*,
  *review*). A factor that names a thing or a task is miscast.
- **Draw factors from a recognizable vocabulary, and keep only the ones that
  matter here.** Familiar quality attributes — often the "-ilities": reliability,
  usability, maintainability, performance, security, portability — make a model's
  intent legible at a glance. Choose them because the subject's needs and risks
  call for them, not by adopting a standard list wholesale; record concerns you
  deliberately leave out under **Known gaps** rather than dropping them silently.
- **Pin the name down in prose.** The same attribute name means different things
  across systems, so a familiar label alone does not carry a definition. The
  factor's body section is where you say what it means for *this* subject, how you
  would know it is met, and what it trades off against its siblings.
- **Keep sibling factors distinct.** Factors under one parent should cover
  different concerns; substantial overlap is a sign that two factors are really
  one, or are cut along the wrong axis. Add sub-factors only to sharpen an
  attribute too broad to assess directly — not for tidiness.
- **Give every factor requirements.** A factor is operational only when
  requirements hang off it — directly or through its sub-factors — whose failure
  would reveal a real deficiency in that attribute. A factor with no requirements
  is a heading, not a part of the model.

**Requirements.** `requirements` is a map of requirement name → requirement,
where the name states the expectation (for example `"unit tests pass"`). A
requirement names an optional `target` and declares exactly one assessment. It
holds nothing about rating levels, which keeps it portable — a requirement, or
the file its assessment points to, can be referenced and evaluated on its own,
outside QUALITY.md.

A requirement carries a *single* assessment — one `prompt` or one `bash`
command, never several and never a list. A `prompt` is one body of criteria
(inline text or a single referenced document), not a collection of separate
prompts. When an expectation feels like it needs more than one prompt, that is
the signal to split it into separate requirements, each with its own single
assessment — keeping every requirement singular and independently evaluable.

The optional `target` is a path or glob pattern (relative to the QUALITY.md
file) identifying the file or directory the requirement is evaluated against. A
glob (for example `./src/**/*.ts`) selects a set of files. When `target` is
omitted, the requirement applies to the QUALITY.md file's directory.

The assessment is exactly one of:

- `prompt: <text | path>` — an *inferential* assessment, judged by a model (or
  a human reviewer). A single value: either inline text stating what the target
  must satisfy, or a path (relative to the QUALITY.md file) to a Markdown
  document holding the same — a checklist, a style guide, a specification to
  conform to. A path is the seam for reuse: the standard lives in its own file
  that both QUALITY.md and, say, an agent skill can load. A document referenced
  by `prompt` may be as long as it needs to be, but it is still one prompt for
  one requirement.
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

## Markdown Body

Below the frontmatter is a Markdown body that documents the model in prose. Where
the frontmatter is the machine-readable summary — the factors and the requirements
under them — the body is the reasoning that justifies it: what the system is, what
"good" means for it, and why these are the right requirements. A reader with only
the frontmatter knows *what* is checked; a reader with the body knows *why*.

The body also carries what the frontmatter cannot. Quality is not intrinsic: "fast
enough" or "reliable enough" mean nothing until you say for whom, doing what, and
under what assumptions. Capturing that context is the body's purpose — and it is
also the grounding a `prompt` assessment needs in order to be judged consistently.

The body is a flat sequence of named sections. The recommended sections are:

| Section | What it captures |
| --- | --- |
| **Overview** | What the system or component is, who depends on it, what "good" means here, and what the model covers — its target and boundary, including dependencies it relies on but does not own. |
| **Needs** | What matters, and to whom — the plain-language statements the requirements answer to. |
| **Risks** | What goes wrong, and for whom, if a need is not met. |
| **Factors** | One subsection per factor, mirroring the frontmatter: what each factor means for this system, how you would know it is met, and any trade-offs it carries against other factors. |
| **Known gaps** | Quality concerns known to matter but deliberately not addressed yet, each with a brief reason. |

**Overview**, **Needs**, and **Factors** are the recommended minimum — together
they make the file a quality *model* rather than a bare list of checks. **Risks**
and **Known gaps** are recommended where they apply. Every section is optional and
should stay short: the body is for shared understanding, not exhaustive
documentation. The format does not restrict the body to these sections — add your
own where a system needs them.

### Example

```markdown
# Quality model — Acme API

## Overview
The Acme API is the public HTTP interface our customers integrate against. It is
maintained by the platform team and depended on by every client app. "Good" here
means it behaves predictably under load and never silently corrupts data. This
model covers the API service and its data layer; the third-party auth provider and
the client SDKs are dependencies, not part of it.

## Needs
- Integrators can trust that a successful response means the data was saved.
- On-call engineers can find the cause of an incident from logs and metrics alone.

## Risks
A silent data-corruption bug is the worst outcome — it erodes customer trust and is
expensive to detect after the fact. A slow endpoint is a lesser problem with clear
workarounds.

## Factors
### Reliability
Customers build on our responses, so a confirmed write must be durable. You would
know it is reliable if a write is acknowledged only after it is committed, and
failures surface as errors rather than false successes.

### Security
The API handles customer data, so access must be authenticated and least-privilege.
When security and convenience conflict, security wins.

## Known gaps
- We do not yet test behavior under sustained peak load.
- Rate-limiting is enforced but not covered by an automated check.
```
