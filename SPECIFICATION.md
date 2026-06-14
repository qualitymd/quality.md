# QUALITY.md Format

QUALITY.md is a plain text representation of a quality model. It can be used to specify and evaluate the quality requirements for a software system or component.

A QUALITY.md file contains two parts: YAML frontmatter with the structured quality model and the Markdown body.

Throughout this document, **must** marks a hard rule that a conforming file or reader has to obey, **should** a recommendation that may be departed from with good reason, and **may** an option left to the author. Prose that explains or motivates a rule is rationale, not an additional rule, and examples illustrate rules without extending them.

## Quality Model

The quality model is embedded in the YAML frontmatter at the beginning of the file. The frontmatter block must begin with a line containing exactly --- and end with a line containing exactly ---. The YAML content between these delimiters is parsed according to the schema defined below.

The model keeps three concerns separate:

- **Requirement** — *what* must be true. A requirement is self-contained: a name, an optional target, and a single assessment. It names no rating level, so it can be lifted out of QUALITY.md and reused on its own — for example, referenced from an agent skill. It may optionally override the scale's *conditions* for its own result, an escape hatch for when one shared condition cannot fit every requirement.
- **Assessment** — *how* the requirement is checked. Either inferential (`prompt`) or computational (`bash`); the key you choose is the declaration of method.
- **Rating** — *how the result is classified*. A single, requirement-agnostic scale (`ratings`) names the levels a result can land on.

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
ratings:                          # optional; an inline scale (below) OR a path to a shared scale file; defaults to pass / fail
  <level-name>:                   # listed best to worst
    displayName: <string>         # optional; human label for the level
    promptCondition: <string>     # optional; criterion a `prompt` is judged against
    bashCondition: <CEL boolean>  # optional; predicate a `bash` result is classified against
factors:
  <factor-name>:                  # a factor has requirements, sub-factors, or both
    requirements:
      <requirement-name>:
        target: <path | glob | list>  # optional; artifact(s) under evaluation; a list narrows, a "!"-prefixed entry excludes
        # exactly one assessment — a single prompt OR a single bash command:
        prompt: <text | path>     # inferential (judged by a model/reviewer)
        bash: <command>           # computational (shell exit status)
        ratings:                  # optional; override the scale's conditions for this requirement
          <level-name>: <condition>   # CEL boolean (bash) or judging criterion (prompt); level from the scale
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
- **Define the attribute, not its requirements.** The same attribute name means
  different things across systems, so a familiar label alone does not carry a
  definition. The factor's body section says what the characteristic is, why it
  matters, and what distinguishes it from its siblings — describing the attribute
  itself, not restating the requirements you hang off it. A good description still
  reads true if those requirements change or the factor is pointed at a different
  target; the subject-specific grounding belongs in Overview, Needs, and Risks and
  in the requirements.
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
names no rating level, which keeps it portable — a requirement, or the file its
assessment points to, can be referenced and evaluated on its own, outside
QUALITY.md. A requirement may optionally override the scale's conditions for
itself (see [Per-requirement rating overrides](#per-requirement-rating-overrides));
doing so trades some of that portability for a per-requirement signal.

A requirement carries a *single* assessment — one `prompt` or one `bash`
command, never several and never a list. A `prompt` is one body of criteria
(inline text or a single referenced document), not a collection of separate
prompts. When an expectation feels like it needs more than one prompt, that is
the signal to split it into separate requirements, each with its own single
assessment — keeping every requirement singular and independently evaluable.

A requirement **should** state its concern concisely, and **should not** restate
an expectation that is self-evident or already guaranteed by the subject's
context — a requirement whose verdict is a foregone conclusion discriminates
nothing and only dilutes the set.

The optional `target` is a path or glob pattern (relative to the QUALITY.md
file) identifying the file or directory the requirement is evaluated against. A
glob (for example `./src/**/*.ts`) selects a set of files. When `target` is
omitted, the requirement applies to the QUALITY.md file's directory (and, when
models compose, its subtree — see [Federation](#federation)).

`target` may also be a **list** of patterns, applied in order, where an entry
prefixed with `!` **excludes** previously matched files (the gitignore
convention). This is how a requirement narrows itself — by selection (a positive
glob) or by exclusion:

```yaml
target:
  - "./**"                      # everything under this model…
  - "!./packages/sandbox/**"    # …except the sandbox package
```

An explicit `target` only ever *narrows* a requirement's reach, and stays within
the model's own directory subtree.

The assessment is exactly one of:

- `prompt: <text | path>` — an *inferential* assessment, judged by a model (or
  a human reviewer). A single value: either inline text stating what the target
  must satisfy, or a path (relative to the QUALITY.md file) to a Markdown
  document holding the same — a checklist, a style guide, a specification to
  conform to. A path is the seam for reuse: the standard lives in its own file
  that both QUALITY.md and, say, an agent skill can load. A document referenced
  by `prompt` may be as long as it needs to be, but it is still one prompt for
  one requirement.
- `bash: <command>` — a *computational* assessment. A shell command whose result
  is classified against the rating scale; by default a zero exit earns the best
  level (see [Computational rating](#computational-rating)).

The key is the declaration of method: `prompt` says "assess this by judgment,"
`bash` says "assess this by running a command." There is no separate field
naming the criteria — the criteria *are* the prompt or the command.

**Rating.** Evaluating a requirement produces a rating: the level its result
lands on, drawn from the `ratings` scale below. The two assessment methods are
classified differently, and a level carries one register for each:

- A `prompt` is **judged** against the scale's `promptCondition`s — the target
  earns the best level when it fully satisfies the requirement's prompt, and lower
  levels as it falls short.
- A `bash` result is **classified** against the scale's `bashCondition`s — the
  levels are tested best to worst and the result takes the first level whose
  `bashCondition` is true (see [Computational rating](#computational-rating)).

The rating is produced by evaluation; the *outcome* is never declared on the
requirement — you never pin a result to a level. Per-level criteria live on the
scale by default — a level's `promptCondition` and `bashCondition`. A requirement
may, however, **override** those criteria with its own `ratings` map keyed on the
scale's level names (see [Per-requirement rating overrides](#per-requirement-rating-overrides)).
The override changes only the conditions by which *this* requirement lands on a
level; the levels themselves — their names, order, and `displayName` — still come
from the shared scale. Because the override names the scale's levels, it couples
that requirement to the scale, so reach for it only when a shared condition
genuinely cannot classify the requirement's result.

**Rating scale.** The optional top-level `ratings` map defines the single scale
shared by every requirement. Each entry is a level name carrying an optional
`displayName`, an optional `promptCondition` (the criterion a `prompt` is judged
against), and an optional `bashCondition` (a CEL boolean a `bash` result is
classified against). Write the conditions generically — in terms of how fully an
evaluation meets its assessment — so the one scale applies to every requirement.
List levels best to worst; the order defines their ranking. A level's conditions
are optional, and the two registers are reached differently: for a `prompt`, the
evaluator reads the levels as an ordered ladder and assigns the highest one the
target satisfies, so a level with no `promptCondition` is an intermediate band —
assigned when a result clearly does better than the level below it but does not
meet the labeled level above; for a `bash` result, a level with no `bashCondition`
is never matched directly and is reachable only as the worst-level fallback (see
[Computational rating](#computational-rating)). When `ratings` is
omitted, the scale defaults to `pass` then `fail`, with `pass` defined as
`bashCondition: "result.success"` (a zero exit). A custom scale might read:

```yaml
ratings:
  A: { displayName: "Excellent",    promptCondition: "Fully satisfies the assessment; no gaps" }
  B: { displayName: "Good" }
  C: { displayName: "Acceptable",   promptCondition: "Satisfies the core of the assessment; minor gaps" }
  D: { displayName: "Poor" }
  E: { displayName: "Unacceptable", promptCondition: "Does not satisfy the assessment" }
```

The `ratings` value may be given **inline** (the map above) or as a **path** to a
YAML file holding that same map. The path resolves relative to the QUALITY.md
file, like a `prompt` or `target` path, and is the seam for sharing one scale
across files so that composing models rate on a common scale (see
[Federation](#federation)).

### Per-requirement rating overrides

The shared scale is written generically so one set of conditions applies to every
requirement. That works when requirements meet their assessments in comparable
ways, but some do not — most often a `bash` requirement whose command emits a
signal the scale's condition cannot interpret. For these, a requirement may carry
an optional `ratings` map that overrides the scale's conditions for itself alone.

The map is keyed on the **scale's** level names; each value is a single condition,
in the register implied by the requirement's assessment:

- under a `bash` requirement, a value is a [CEL](https://cel.dev) boolean over
  `result`, classified exactly as a `bashCondition` (see below);
- under a `prompt` requirement, a value is a judging criterion, applied exactly as
  a `promptCondition`.

Only the conditions are overridden. The level set, its order, and each level's
`displayName` still come from the shared scale, and classification proceeds as it
otherwise would — best to worst, first match wins, worst level as the fallback. A
level the override omits keeps the scale's condition for that level, if any. A
level name not present in the scale is a configuration error.

For example, a scale bands a coverage percentage one command prints, while a
second `bash` requirement reports a different signal and overrides the bands for
itself:

```yaml
ratings:
  A: { bashCondition: "double(result.stdout.trim()) >= 90" }
  B: { bashCondition: "double(result.stdout.trim()) >= 80" }
  C: { bashCondition: "double(result.stdout.trim()) >= 70" }
  fail: {}
factors:
  maintainability:
    requirements:
      "line coverage":
        bash: "pnpm coverage:lines --print"   # uses the scale's bands as-is
      "mutation score":
        bash: "pnpm mutation --print"          # a different signal; rebands for itself
        ratings:
          A: "double(result.stdout.trim()) >= 75"
          B: "double(result.stdout.trim()) >= 60"
          C: "double(result.stdout.trim()) >= 50"
```

This keeps a requirement portable by default — it names no level until it opts
into an override — and confines the coupling to the requirements that need it.

### Computational rating

A `bash` requirement is classified by running its command and evaluating the
scale's `bashCondition` expressions against the result. It is deterministic — no
model judgment is involved.

Each `bashCondition` is a [CEL](https://cel.dev) (Common Expression Language)
boolean, evaluated against a single `result` describing the command run:

| Field | Type | Meaning |
| --- | --- | --- |
| `result.success` | bool | the command exited zero |
| `result.exit` | int | exit status |
| `result.stdout` | string | captured standard output |
| `result.stderr` | string | captured standard error |

Because CEL does not coerce between types, a command's raw text output is bridged
to the value a condition tests through CEL's standard library, written
receiver-style:

- **String operations** come from CEL's strings extension as member calls —
  `result.stdout.trim()` (command output carries trailing newlines),
  `.lowerAscii()`, `.contains(s)`, `.startsWith(s)`, `.endsWith(s)`,
  `.matches(re)` (regex), `.size()`.
- **Numeric parsing** uses CEL's standard conversions, which parse from a string —
  `double(result.stdout.trim())`, `int(result.stdout.trim())`.
- **JSON** has no standard parse in CEL, so the evaluator provides one convenience,
  receiver-style: `result.stdout.json()` parses the output into a value (map, list,
  number, …) for further indexing.

**Classification.** The levels are tested in order, best to worst; the result
takes the **first** level whose `bashCondition` is true. If no level matches, the
result takes the **worst** level — the scale denies by default. A level with no
`bashCondition` is never selected by computation; it is reachable only as that
worst-level fallback (which is why the default `fail` needs none). A
`bashCondition` that fails to evaluate — `.json()` on output that is not JSON, or
`double()` on output that is not a number, say — is a configuration error in the
model, surfaced as such rather than silently scored.

**The default scale, made explicit.** Omitting `ratings` is equivalent to:

```yaml
ratings:
  pass: { bashCondition: "result.success" }   # zero exit
  fail: {}                                    # the fallback
```

So the default `bash` behavior — best level on a zero exit, worst otherwise — is
just this scale. An author can refine the pass condition on the scale without
touching any requirement:

```yaml
ratings:
  pass: { bashCondition: "result.success && result.stderr == ''" }
  fail: {}
```

or band a numeric signal a command prints on stdout:

```yaml
ratings:
  A: { bashCondition: "double(result.stdout.trim()) >= 90" }
  B: { bashCondition: "double(result.stdout.trim()) >= 80" }
  C: { bashCondition: "double(result.stdout.trim()) >= 70" }
  fail: {}
```

A scale that bands on `result.stdout` like this assumes every `bash` requirement
under it emits a comparable value. When one does not, that requirement can supply
its own bands with a [per-requirement rating override](#per-requirement-rating-overrides)
rather than forcing every command onto one signal.

Classification yields the rating; whether a given rating should *gate* — fail a
build, block a change — is the evaluating tool's concern, not the format's.

### Invalid examples

The constructs above are easier to pin down against cases that fall just outside
them. Each requirement below is malformed for the reason its comment gives.

```yaml
factors:
  security:
    requirements:
      "no secrets committed":
        prompt: "No credentials appear in source."
        bash: "git secrets --scan"   # ✗ two assessments — declare exactly one of prompt/bash
      "dependencies are current": {}  # ✗ no assessment — every requirement needs exactly one
      "input is validated":
        prompt:                       # ✗ a list — an assessment is a single value, never a list
          - "Query parameters are validated."
          - "Request bodies are validated."
  performance: {}                     # ✗ a factor with neither requirements nor sub-factors
```

A per-requirement `ratings` override may only re-state levels the scale already
declares; naming one it does not is a configuration error, not a new level:

```yaml
ratings:
  pass: { bashCondition: "result.success" }
  fail: {}
factors:
  maintainability:
    requirements:
      "line coverage":
        bash: "pnpm coverage:lines --print"
        ratings:
          gold: "double(result.stdout.trim()) >= 90"   # ✗ "gold" is not a level in the scale
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
| **Factors** | One subsection per factor, mirroring the frontmatter: what each factor's quality attribute is and why it matters, and how it differs from sibling factors — describing the attribute itself, not restating the requirements attached to it. |
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

## Federation

A repository may contain more than one `QUALITY.md`. Rather than one file growing
to cover an entire system, each component can carry its own model, and the files
**compose**. This is the same nesting convention as `AGENTS.md`: place a model
where it applies, and the nearest one carries the local detail. It suits
monorepos, and lets a large model be decomposed into per-component models that are
each legible on their own.

Composition follows four rules:

- **Ownership stays at the defining file.** Each `QUALITY.md` is a complete,
  standalone model; no file's content is merged into another. Composition is a
  reading of the set, not a rewrite of any file.
- **Scope spans downward.** A model's requirements apply to its whole directory
  subtree, *including the subtrees of nested models*. A nested model **adds**
  requirements for its part of the tree; it does not remove what an enclosing
  model already requires there.
- **`target` narrows.** A requirement with no `target` spans the model's whole
  subtree; an explicit `target` narrows it to the files it matches — by selection
  or by exclusion (above). A model exempts part of its subtree from a requirement
  by excluding those paths from that requirement's `target`, so the exemption is
  visible on the requirement that owns it.
- **Scales are shared by reference, not inherited.** A nested model that omits
  `ratings` uses the default `pass`/`fail` scale, not the enclosing model's; to
  rate on a common scale, models reference one shared scale file (above).

So two levers compose a system's quality model, each with one job: **placement
adds** (a nested model contributes requirements for its subtree) and **`target`
narrows** (a requirement limits or exempts its reach). A conforming tool
discovers the set, evaluates each model over its scope, and reports the result as
a tree; that operational detail is specified with the CLI in
[`specs/cli-federation.md`](specs/cli-federation.md).

## Conformance, extensibility, and versioning

### Minimal core

Every QUALITY.md needs only one thing: frontmatter containing a `factors` map
with at least one factor, each factor carrying at least one `requirements` entry
or sub-`factor`. The `ratings` scale and the whole Markdown body are optional.
That minimal core is the stable contract; everything else is additive.

### Unrecognized content

The format is meant to grow through use, so a conforming reader **ignores content
it does not recognize** rather than rejecting it: unknown frontmatter keys and
unknown body sections are preserved and skipped, never errors. Recognizing a
genuinely custom extension is not optional — ignore, don't reject — though a tool
*may* warn about a key that looks like a typo of a known one (`factor:` for
`factors:`). Malformed *known* content is the opposite case: it is an error, not
an extension (see below).

### Edge cases

A conforming reader resolves these the same way:

| Case | Treatment |
| --- | --- |
| **No frontmatter, or no closing `---`** | Invalid — there is no model to read. A QUALITY.md must carry a fenced frontmatter block. |
| **Empty top-level `factors`** | Invalid — a model with no factors declares nothing. |
| **Empty `requirements` / `factors` under a factor** | A degenerate branch with nothing to assess: permitted but meaningless, and a tool may warn. A factor whose collections are all empty fails the "at least one" rule. |
| **Empty assessment value** (`prompt:` or `bash:` with no value) | Invalid — the assessment *is* its own criteria, so it must be non-empty. |
| **Duplicate keys among siblings** (factors, requirements, or levels sharing a name) | Invalid — names must be unique among siblings; a reader rejects duplicates rather than silently choosing one. |
| **Ordering** | Significant only within `ratings`, where levels are listed best to worst. The order of factors, sub-factors, and requirements does not affect evaluation; a tool may preserve authoring order for display. |
| **Case** | Schema keys are lower-case and case-sensitive (`factors`, not `Factors`). Author-defined names — factors, requirements, and rating levels — are case-sensitive and matched exactly, so a per-requirement override must name a level using the same case as the scale. |

### Versioning

The format evolves **additively**. The minimal core above is fixed, and new
capability arrives as new *optional* keys and sections. Because a reader ignores
content it does not recognize, a file written against a newer revision still reads
under an older tool — the unknown parts are skipped — and a file written against
an older revision stays valid under a newer one. A change that would invalidate an
existing conforming file is a breaking change, reserved for a new major version of
the format: the aim is to grow through use, not to churn the files already written.
