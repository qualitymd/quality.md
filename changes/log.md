# Changes Update Log

## 2026-06-23

- **Creation + Implementation + Archival**: Created
  [0072 - Setup context checkpoint](archive/0072-setup-context-checkpoint.md)
  with its child folder
  ([index](archive/0072-setup-context-checkpoint/index.md),
  [functional spec](archive/0072-setup-context-checkpoint/spec.md),
  [design doc](archive/0072-setup-context-checkpoint/design.md)), advanced it
  through `Draft`, `Design`, `In-Progress`, and `In-Review` to `Done`, and
  archived it. Implemented skill-only with **no CLI/Go change**: setup now
  presents primary users/outcomes, maintainers/collaborators, other stakeholders,
  and missing or not-agent-accessible context as a compact correction-oriented
  checkpoint instead of four separate open-ended prompts, preserving provenance
  and recording omitted low-confidence gaps honestly. Synced the durable setup
  workflow spec, workflow/spec logs, runtime setup workflow, `CHANGELOG.md`, and
  archive index.

- **Creation + Implementation + Archival**: Created
  [0071 - Setup open-ended review gate](archive/0071-setup-open-ended-review-gate.md)
  with its child folder
  ([index](archive/0071-setup-open-ended-review-gate/index.md),
  [functional spec](archive/0071-setup-open-ended-review-gate/spec.md),
  [design doc](archive/0071-setup-open-ended-review-gate/design.md)), advanced
  it through `Draft`, `Design`, `In-Progress`, and `In-Review` to `Done`, and
  archived it. Implemented skill-only with **no CLI/Go change**: setup's final
  review recap now uses friendly, open-ended wording that preserves the
  `"looks good"` fast path while inviting priorities, worries, wording, edge
  cases, repo-invisible context, and other useful last-call input before
  authoring. Synced the durable setup workflow spec, workflow/spec logs, runtime
  setup workflow, `CHANGELOG.md`, and archive index.

- **Creation + Implementation + Archival**: Created
  [0070 - Setup missing-context provenance](archive/0070-setup-missing-context-provenance.md)
  with its child folder
  ([index](archive/0070-setup-missing-context-provenance/index.md),
  [functional spec](archive/0070-setup-missing-context-provenance/spec.md),
  [design doc](archive/0070-setup-missing-context-provenance/design.md)),
  advanced it through `Draft`, `Design`, `In-Progress`, and `In-Review` to
  `Done`, and archived it. Implemented skill-only with **no CLI/Go change**:
  setup missing-context discovery now treats material context as
  agent-accessible only when supported by repository/tool/source evidence or
  explicit setup-provided context, prohibits choices that assume low/no-evidence
  project facts are sufficiently understood, and preserves setup-provided
  provenance in authored `QUALITY.md` body context. Synced the durable setup
  workflow spec, workflow/spec logs, runtime setup workflow, `CHANGELOG.md`, and
  archive index.

- **Implementation + Archival**: Implemented and advanced
  [0069 - Setup review gate and discovery trim](archive/0069-setup-review-gate-and-pedagogy-trim.md)
  through `Design`, `In-Progress`, and `In-Review` to `Done`, moving it (parent
  and folder) into [`archive/`](archive/). Added the
  [design doc](archive/0069-setup-review-gate-and-pedagogy-trim/design.md), then
  implemented skill-only with **no CLI/Go change**: setup discovery now asks nine
  questions, removes modeling rigor and review posture as user-facing discovery
  questions, adds a Rating Scale confirmation question, trims per-question
  pedagogy to purpose/context only, and treats the final recap as a hard review
  gate before authoring. Synced the runtime setup workflow, durable setup spec,
  spec logs, and `CHANGELOG.md`. Removed it from the open-cases [index](index.md)
  and added it to the [archive index](archive/index.md).

- **Update**: Amended
  [0069 - Setup review gate and discovery trim](0069-setup-review-gate-and-pedagogy-trim.md)
  while still in `Draft` to add a rating-scale confirmation question. The new
  question teaches that Rating Levels are configurable model vocabulary, not
  baked into QUALITY.md, while recommending the standard
  `outstanding`/`target`/`minimum`/`unacceptable` scale and explaining the
  decision role of each level. Setup still must not ask the user to invent
  custom Rating Level names during discovery.

- **Implementation + Archival**: Implemented and advanced
  [0068 - Always-on setup feedback log](archive/0068-always-on-setup-feedback-log.md)
  through `In-Progress` and `In-Review` to `Done`, moving it (parent and folder)
  into [`archive/`](archive/). Implemented skill-only with **no CLI/Go change**:
  setup now creates the current run's feedback log during preflight after CLI
  support and the run frame, updates the current file for material
  workflow-experience events, and finalizes it at close with stable frontmatter
  metadata, lifecycle status, a timeline, and explicit no-notable-content notes.
  Synced the durable workflow feedback-log sub-spec, setup workflow spec, parent
  `/quality` skill spec, runtime skill files, CLI quick reference, spec logs, and
  `CHANGELOG.md`; no public durable docs changed. Removed it from the open-cases
  [index](index.md) and added it to the [archive index](archive/index.md).

- **In-Review**: Completed implementation of
  [0068 - Always-on setup feedback log](0068-always-on-setup-feedback-log.md).
  Implemented skill-only with **no CLI/Go change**: setup now creates the
  current run's feedback log during preflight after CLI support and the run
  frame, updates the current file for material workflow-experience events, and
  finalizes it at close. Synced the durable workflow feedback-log sub-spec,
  setup workflow spec, parent `/quality` skill spec, runtime skill files, CLI
  quick reference, spec logs, and `CHANGELOG.md`. No public durable docs changed.

- **In-Progress**: Advanced
  [0068 - Always-on setup feedback log](0068-always-on-setup-feedback-log.md)
  from `Design` to `In-Progress`; spec and design were settled and
  implementation began across durable specs, the runtime skill, and release
  notes.

- **Update**: Amended
  [0069 - Setup review gate and discovery trim](0069-setup-review-gate-and-pedagogy-trim.md)
  while still in `Draft` to include removing the modeling-rigor and
  review-posture discovery questions. Modeling rigor may remain an internal
  setup-brief inference, and review/loop expectations move to setup closeout
  next-step routing rather than discovery. No replacement question added.

- **Creation**: Added
  [0069 - Setup review gate and discovery trim](0069-setup-review-gate-and-pedagogy-trim.md)
  at `Draft` with its child folder
  ([index](0069-setup-review-gate-and-pedagogy-trim/index.md),
  [functional spec](0069-setup-review-gate-and-pedagogy-trim/spec.md)). The case
  makes `/quality setup` stop after discovery, present the final recap, and wait
  for an explicit user response before authoring `QUALITY.md`; structured
  question-tool completion does not satisfy that review gate. It also trims
  per-question teaching copy to purpose/context only, removing repeated
  "how to change it later" guidance while allowing one general living-document
  note. No CLI/Go change expected; skill + durable setup spec + changelog only.
  Design and implementation not started. Listed it under open cases in the
  bundle [index](index.md).

- **Advance**: Moved
  [0068 - Always-on setup feedback log](0068-always-on-setup-feedback-log.md) to
  `Design` and authored its
  [design doc](0068-always-on-setup-feedback-log/design.md). The design keeps the
  artifact skill-only and local, emits the run frame before the first feedback-log
  write, creates `.quality/logs/<started-at>-setup-feedback-log.md` immediately
  after preflight has CLI/model metadata, updates the current run's file in place
  for material workflow-experience events, and finalizes it at close. It records
  the current-run overwrite boundary, the timeline/body split, stop-handling edge
  cases, and the rejected alternatives of close-only creation, append-only events,
  a JSONL sidecar, and a CLI helper. Implementation not started.

- **Creation**: Added
  [0068 - Always-on setup feedback log](0068-always-on-setup-feedback-log.md) at
  `Draft` with its child folder
  ([index](0068-always-on-setup-feedback-log/index.md),
  [functional spec](0068-always-on-setup-feedback-log/spec.md)). The case changes
  `/quality setup` feedback logging from optional close-step authoring to an
  always-created run artifact under `.quality/logs/` that is created during
  preflight, updated as the workflow progresses, and finalized at close with
  stable frontmatter metadata and body sections. It remains skill-only, local,
  never transmitted, and bounded to workflow-experience feedback rather than
  `QUALITY.md` model rationale or evaluation records. Design and implementation
  not started. Listed it under open cases in the bundle [index](index.md).

- **Implementation + Archival**: Implemented and advanced
  [0067 - Setup discovery pedagogy](archive/0067-setup-discovery-pedagogy.md)
  through `In-Progress` and `In-Review` to `Done`, moving it (parent and folder)
  into [`archive/`](archive/). Implemented skill-only with **no code change**:
  added authored per-question background and how-to-change-later copy inline in
  the runtime [setup workflow](../skills/quality/workflows/setup.md) with a
  teaching-over-round-trips framing, made setup ask every one of the ten questions
  every run (removed the accept-all-and-skip escape; kept a per-question fast
  confirm and show-all-at-once), relabeled confidence from
  `strongly inferred`/`weakly inferred`/`assumed` to `Low`/`Med`/`High` (evidence
  note retained, no-evidence → `Low`), and added a final review recap before
  authoring. Synced the durable
  [setup spec](../specs/skills/quality-skill/workflows/setup.md) with promoted
  rationale annotations (the parent skill spec's generic "confidence-labeled
  defaults" phrasing reviewed and left unchanged), `CHANGELOG.md`, and the
  spec log. Reconciled the **Affected artifacts** list: no `qualitymd` CLI/Go
  change and the `status` lifecycle `readiness` field is unchanged; the listed
  public `use-quality-skill.md` guide was removed independently by concurrent docs
  cleanup, so 0067 lands with no durable-docs delta. Skill `metadata.version`
  bump deferred to the release cut. The In-Review gate was
  collapsed at the user's explicit direction. Removed it from the open-cases
  [index](index.md) and added it to the [archive index](archive/index.md).

- **Implementation + Archival**: Implemented and advanced
  [0066 - Setup feedback log](archive/0066-setup-feedback-log.md) to `Done`,
  moving it (parent and folder) into [`archive/`](archive/). Implemented
  skill-only with **no code change**: added the new durable
  [workflow feedback log](../specs/skills/quality-skill/workflows/setup/feedback-log.md)
  sub-spec under a new `workflows/setup/` folder, made
  [`setup`](../specs/skills/quality-skill/workflows/setup.md) its parent with an
  amended (narrowly widened) mutation surface and close-step authoring, and
  recorded the shared artifact plus its redaction/no-transmission boundary in the
  parent [`/quality` skill](../specs/skills/quality-skill/quality-skill.md) spec.
  Synced the runtime skill
  ([`SKILL.md`](../skills/quality/SKILL.md),
  [`workflows/setup.md`](../skills/quality/workflows/setup.md),
  [`cli-quick-reference.md`](../skills/quality/resources/cli-quick-reference.md)),
  the [use-the-skill guide](../docs/guides/use-quality-skill.md), `CHANGELOG.md`,
  and the spec/doc logs. Skill `metadata.version` bump deferred to the release
  cut. Removed it from the open-cases [index](index.md) and added it to the
  [archive index](archive/index.md).

- **Advance**: Moved
  [0067 - Setup discovery pedagogy](0067-setup-discovery-pedagogy.md) to `Design`
  and authored its [design doc](0067-setup-discovery-pedagogy/design.md): inline
  per-question teaching copy composed with the 0065 presentation tiers, the
  accept-all-and-skip escape removed in favor of a per-question fast confirm,
  the `assumed`→`Low (no signal)` mapping that lets `Low`/`Med`/`High` plus the
  evidence note carry the old vocabulary's meaning, and the final recap inserted
  between discovery and authoring. No code involved.

- **Creation**: Added
  [0067 - Setup discovery pedagogy](0067-setup-discovery-pedagogy.md) at `Draft`
  with its child folder ([index](0067-setup-discovery-pedagogy/index.md),
  [functional spec](0067-setup-discovery-pedagogy/spec.md)). The case
  repositions `/quality setup` discovery as teaching-first: authored
  per-question background and how-to-change-later copy inline in the skill, asks
  every discovery question even at high confidence (removing/revising the
  accept-all-defaults-and-skip escape), relabels the confidence vocabulary to
  `Low`/`Med`/`High` (retaining the evidence note), and adds a final review
  recap of the full question/answer set with a last-chance comment before
  writing `QUALITY.md`. No CLI/Go change; skill + specs + docs only. Revisits
  the confidence vocabulary and accept-all escape that archived
  [0065](archive/0065-setup-discovery-and-close-refinements.md) deliberately
  left alone, and is independent of [0066](0066-setup-feedback-log.md). Spec
  authored; design doc deferred until the spec settles. Listed it under open
  cases in the bundle [index](index.md).

- **Creation**: Added
  [0066 - Setup feedback log](0066-setup-feedback-log.md) at `Design` with its
  child folder ([index](0066-setup-feedback-log/index.md),
  [functional spec](0066-setup-feedback-log/spec.md),
  [design doc](0066-setup-feedback-log/design.md)). The case adds a
  hand-authored, **skill-only** workflow feedback log under `.quality/logs/`
  (`<timestamp>-setup-feedback-log.md`) that records setup-experience friction,
  UX/AX, and efficiency signals for improving the skill — distinct from
  evaluation's per-run `debug-log.md` and from the quality log under
  `.quality/log/`. No CLI/Go change: the directory is created on demand (as
  evaluation already does), the log is recorded locally and never transmitted, so
  no opt-in is needed; secrets and raw prompt-injection text are never written
  and sensitive project context is sanitized. Spec and design authored; no code
  involved. Listed it under open cases in the bundle [index](index.md).

- **Archival**: Advanced
  [0065 - Setup discovery and close refinements](archive/0065-setup-discovery-and-close-refinements.md)
  to `Done` and moved it (parent and folder) into [`archive/`](archive/).
  Implementation landed: made setup discovery agent-agnostic (present all ten
  questions, iterate one at a time without a structured question affordance, page
  through one when available, escapes on request), added the
  read-the-`qualitymd init`-scaffold-before-authoring step, disentangled model
  maturity (`starter`/`immature`/`evaluation-ready`) from the CLI's lifecycle
  `readiness` in the setup close and the Top 10 guide, and renamed the skill
  `modes/` folder to `workflows/` across the runtime skill and `specs/` mirror
  with all live path references updated and append-only logs left frozen. Updated
  durable specs, public docs, and `CHANGELOG.md`; verified `mise run check`. The
  In-Review review gate was collapsed at the user's explicit direction. Updated
  the bundle [index](index.md) and [archive index](archive/index.md).

- **In-Progress**: Advanced
  [0065 - Setup discovery and close refinements](archive/0065-setup-discovery-and-close-refinements.md)
  from `Design` to `In-Progress`; spec and design were settled and implementation
  began across the runtime skill, durable specs, docs, and the folder rename.

- **Archival**: Advanced
  [0062 - Remove wizard mode](archive/0062-remove-wizard-mode.md),
  [0063 - Contextual setup flow](archive/0063-contextual-setup-flow.md), and
  [0064 - Structured setup workflow](archive/0064-structured-setup-workflow.md)
  from `In-Review` to `Done` and moved each parent concept and child folder into
  [`archive/`](archive/index.md). Added their archive [index](archive/index.md)
  entries and removed the open-cases entries from the bundle
  [index](index.md). Repointed the live [0065](0065-setup-discovery-and-close-refinements.md)
  "Relationship to 0064" link into `archive/` and updated its now-stale "0064 is
  In-Review" note. Append-only `log.md` references under
  [`specs/`](../specs/log.md), `specs/skills/quality-skill/`, and
  [`docs/`](../docs/log.md) stay frozen at their original paths as historical
  record.

- **Creation**: Added change
  [0065 - Setup discovery and close refinements](0065-setup-discovery-and-close-refinements.md)
  (`status: Design`) with its
  [functional spec](0065-setup-discovery-and-close-refinements/spec.md),
  [design doc](0065-setup-discovery-and-close-refinements/design.md), and
  [index](0065-setup-discovery-and-close-refinements/index.md). The case captures
  four frictions from a first field run of `/quality setup`: make discovery
  agent-agnostic and present all ten questions (iterating one at a time when no
  structured question affordance exists), read the `qualitymd init` scaffold
  before authoring it, disentangle the skill's model-maturity judgment from the
  CLI's lifecycle `readiness`, and take up the `modes/` → `workflows/` folder
  rename 0064 deferred. Records the affected runtime skill, durable specs, docs,
  and packaging; notes that append-only `log.md` files keep historical `modes/`
  references frozen. Spec and design are settled; no code, runtime, or durable
  spec edits made yet. Updated the open-cases entry in the bundle
  [index](index.md).

- **In-Review**: Completed implementation of
  [0064 - Structured setup workflow](0064-structured-setup-workflow.md) and
  advanced it from `In-Progress` to `In-Review`. Rewrote runtime setup guidance
  as an operator workflow, added a setup brief and ten concrete discovery
  questions, aligned durable setup and parent skill specs, updated
  getting-started guidance, public README copies, specs logs, and changelog, and
  preserved setup's `QUALITY.md`-only mutation boundary. Verified
  `mise run check`.

- **In-Progress**: Advanced
  [0064 - Structured setup workflow](0064-structured-setup-workflow.md) from
  `Design` to `In-Progress`. The functional spec and design doc are settled;
  implementation begins across runtime setup guidance, durable skill specs,
  public docs, logs, and changelog.

- **Design**: Advanced
  [0064 - Structured setup workflow](0064-structured-setup-workflow.md) from
  `Draft` to `Design` and added its
  [design doc](0064-structured-setup-workflow/design.md). The design keeps
  existing dispatch paths stable, rewrites runtime setup guidance as an operator
  playbook, adds a setup brief template, defaults to one compact confirmation
  prompt with all ten discovery questions, and preserves the `QUALITY.md`-only
  mutation boundary. Updated the open-cases entry in the bundle
  [index](index.md). Code not started.

- **Draft**: Created
  [0064 - Structured setup workflow](0064-structured-setup-workflow.md)
  (`Draft`) with its
  [functional spec](0064-structured-setup-workflow/spec.md) and
  [child index](0064-structured-setup-workflow/index.md). The case turns
  `/quality setup` guidance into an explicit setup workflow with a setup brief,
  ten concrete discovery questions, confidence-labeled defaults, prompt framing,
  workflow terminology, and the existing `QUALITY.md`-only mutation boundary.
  Added the case to the open-cases list in the bundle [index](index.md). Design
  and code not started.

- **In-Review**: Completed implementation of
  [0063 - Contextual setup flow](0063-contextual-setup-flow.md) and advanced it
  from `In-Progress` to `In-Review`. Updated runtime setup guidance, durable
  setup and quality-log contracts, getting-started and Top 10 checklist
  guidance, public docs, changelog, and OKF logs so setup analyzes context,
  asks confidence-labeled discovery questions, writes only `QUALITY.md`,
  validates/readiness-checks the model, and offers next-step choices without
  running evaluation, writing the quality log, creating issues, or configuring
  integrations. Verified `mise run check`.

- **In-Progress**: Advanced
  [0063 - Contextual setup flow](0063-contextual-setup-flow.md) from `Design` to
  `In-Progress`. The functional spec and design doc are settled;
  implementation begins across runtime setup guidance, durable skill specs,
  quality-log contracts, public docs, logs, and changelog.

- **Design**: Advanced
  [0063 - Contextual setup flow](0063-contextual-setup-flow.md) from `Draft` to
  `Design` and added its
  [design doc](0063-contextual-setup-flow/design.md). The design keeps setup
  skill-driven, uses a bounded context-analysis fact sheet, asks compact
  discovery questions with confidence-labeled defaults, writes only
  `QUALITY.md`, validates with lint plus Top 10 readiness inspection, and offers
  next-step choices without running evaluation or configuring integrations.
  Updated the open-cases entry in the bundle [index](index.md). Code not
  started.

- **Draft**: Created
  [0063 - Contextual setup flow](0063-contextual-setup-flow.md) (`Draft`) with
  its [functional spec](0063-contextual-setup-flow/spec.md) and
  [child index](0063-contextual-setup-flow/index.md). The case reworks
  `/quality setup` into a context-informed discovery flow that writes only
  `QUALITY.md`, validates/readiness-checks the model, and offers next-step
  choices without running evaluation, writing the quality log, creating issues,
  or configuring recurring review automation. Added the case to the open-cases
  list in the bundle [index](index.md). Design and code not started.

- **In-Review**: Completed implementation of
  [0062 - Remove wizard mode](0062-remove-wizard-mode.md) and advanced it from
  `In-Progress` to `In-Review`. Removed runtime and durable wizard mode files,
  folded bare/ambiguous `/quality` handling into read-only orientation, removed
  wizard from public docs and setup handoffs, updated quality-log/checklist
  wording, reconciled indexes/logs/changelog, and verified `mise run check`.

- **In-Progress**: Advanced
  [0062 - Remove wizard mode](0062-remove-wizard-mode.md) from `Design` to
  `In-Progress`. The functional spec and design doc are settled; implementation
  begins across runtime `/quality` guidance, durable skill specs, public docs,
  indexes, logs, and changelog.

- **Design**: Advanced
  [0062 - Remove wizard mode](0062-remove-wizard-mode.md) from `Draft` to
  `Design` and added its
  [design doc](0062-remove-wizard-mode/design.md). The design treats this as a
  surface reduction rather than a rename, absorbs safe read-only orientation
  into the parent skill routing contract, deletes public wizard mode files, and
  keeps `/quality status` and `/quality next` out of the public contract. Updated
  the open-cases entry in the bundle [index](index.md). Code not started.

- **Draft**: Created
  [0062 - Remove wizard mode](0062-remove-wizard-mode.md) (`Draft`) with its
  [functional spec](0062-remove-wizard-mode/spec.md) and
  [child index](0062-remove-wizard-mode/index.md). The case removes `wizard`
  from the `/quality` public contract without promoting `status` or `next` as
  replacement modes, while preserving read-only orientation for bare or
  ambiguous requests. Added the case to the open-cases list in the bundle
  [index](index.md). Design and code not started.

- **Done**: Landed and archived
  [0061 - Natural scope labels](archive/0061-natural-scope-labels.md) —
  advanced it through implementation to `Done`, moved the parent concept and
  child folder into [`archive/`](archive/index.md), added it to the archive
  [index](archive/index.md), and removed it from the open-cases [index](index.md).
  The implementation updates README examples, runtime `/quality` scope
  resolution guidance, durable `/quality` skill specs, specs logs, and the
  changelog so natural Area and Factor labels are the primary documented scoped
  input while qualified references remain exact-addressing syntax.

- **Design**: Advanced
  [0061 - Natural scope labels](archive/0061-natural-scope-labels.md) from
  `Draft` to `Design` and added its
  [design doc](archive/0061-natural-scope-labels/design.md). The design keeps
  label resolution in the `/quality` skill, treats natural labels as human-edge
  input only, preserves qualified references and stable artifact identifiers,
  and records exact matching plus targeted ambiguity prompts as the
  implementation shape. Updated the open-cases entry in the bundle
  [index](index.md). Code not started.

- **Draft**: Created
  [0061 - Natural scope labels](archive/0061-natural-scope-labels.md) (`Draft`)
  with its [functional spec](archive/0061-natural-scope-labels/spec.md) and
  [child index](archive/0061-natural-scope-labels/index.md). The case makes
  natural Area and Factor labels the primary documented scoped-evaluation input
  for `/quality evaluate`, keeps qualified model references as exact-addressing
  syntax, and preserves stable IDs in durable evaluation artifacts. Added the
  case to the open-cases list in the bundle [index](index.md). Design and code
  not started.

## 2026-06-22

- **Done**: Landed and archived
  [0059 - Unqualified model references](archive/0059-unqualified-model-references.md)
  and [0060 - Friendly path display](archive/0060-friendly-path-display.md) —
  advanced both to `Done`, moved their parent concepts and child folders into
  [`archive/`](archive/index.md), added them to the archive [index](archive/index.md),
  and removed them from the open-cases [index](index.md). Both cases are part of
  the v0.9.0 release state.

- **In-Review**: Completed implementation of
  [0060 - Friendly path display](archive/0060-friendly-path-display.md) and advanced it
  from `In-Progress` to `In-Review`. Separated Area/Factor/Rating display
  helpers from reference helpers; rendered `/` for root Area paths in human
  Markdown report path fields; preserved `area:root`, `root`, `root::factor`,
  and structured `report.json` identifiers; aligned durable specs, runtime
  `/quality` guidance, generated examples, logs, and changelog. Verified
  `go test ./internal/evaluation` and `mise run check`.

- **In-Progress**: Advanced
  [0060 - Friendly path display](archive/0060-friendly-path-display.md) from `Design`
  to `In-Progress`. The functional spec and design doc are settled;
  implementation begins across display/reference helper separation, report path
  rendering, durable specs, runtime skill guidance, generated examples, and
  changelog.

- **Design**: Advanced
  [0060 - Friendly path display](archive/0060-friendly-path-display.md) from `Draft` to
  `Design` and added its
  [design doc](archive/0060-friendly-path-display/design.md). The design separates
  display helpers from qualified and unqualified reference helpers, keeps `/`
  out of reference parsing, and updates report rendering to use display values
  in human path fields. Updated the open-cases entry in the bundle
  [index](index.md). Code not started.

- **Draft**: Created
  [0060 - Friendly path display](archive/0060-friendly-path-display.md) (`Draft`) with
  its [functional spec](archive/0060-friendly-path-display/spec.md) and
  [child index](archive/0060-friendly-path-display/index.md). The case separates
  display values from qualified and unqualified model references, keeps
  reference grammar stable, and proposes rendering the root Area path as `/` in
  human report display contexts. Added the case to the open-cases list in the
  bundle [index](index.md). Design and code not started.

- **In-Review**: Completed implementation of
  [0059 - Unqualified model references](archive/0059-unqualified-model-references.md)
  and advanced it from `In-Progress` to `In-Review`. Added unqualified Area,
  Factor, and Rating reference helpers; type-specific unqualified parsers while
  preserving strict qualified parsing; unqualified Area Breakdown `Path`
  rendering in `report-summary.md` and `report.md`; generated example updates;
  durable spec alignment; runtime `/quality` guidance; and changelog. Verified
  `go test ./...` and `mise run check`.

- **In-Progress**: Advanced
  [0059 - Unqualified model references](archive/0059-unqualified-model-references.md)
  from `Design` to `In-Progress`. The functional spec and design doc are
  settled; implementation begins across unqualified reference helpers,
  type-specific parsing, shared Area Breakdown rendering, durable specs, runtime
  skill guidance, generated examples, and changelog.

- **Done**: Landed and archived
  [0058 - Model reference identifiers](archive/0058-model-reference-identifiers.md)
  — advanced it to `Done` and moved the parent concept and its
  [folder](archive/0058-model-reference-identifiers/index.md) into
  [`archive/`](archive/index.md). The case defined strict Area name, Factor
  name, and Rating Level ID grammar; canonical typed model references;
  edge-only shorthand boundaries; lint diagnostics; JSON Schema patterns;
  revised Area Breakdown columns; durable specs; runtime `/quality` guidance;
  scaffold updates; docs; and changelog. Updated the archive [index](archive/index.md)
  and removed the open-cases entry from the bundle [index](index.md).

- **Refinement**: Updated
  [0059 - Unqualified model references](archive/0059-unqualified-model-references.md)
  to explicitly include the durable `/quality` reporting spec in the affected
  artifacts, so agent-facing report guidance is updated alongside
  `report-summary.md`, `report.md`, `report.json`, generated examples, and the
  shared report renderer.

- **Design**: Created
  [0059 - Unqualified model references](archive/0059-unqualified-model-references.md)
  (`Design`) with its
  [functional spec](archive/0059-unqualified-model-references/spec.md),
  [design doc](archive/0059-unqualified-model-references/design.md), and
  [child index](archive/0059-unqualified-model-references/index.md). The case defines
  unqualified references as a bounded fixed-type form for Areas, Factors, and
  Rating Levels; preserves qualified references for mixed-reference and
  machine-readable surfaces; and plans named helper functions plus Area
  Breakdown rendering updates. Added the case to the open-cases list in the
  bundle [index](index.md). Code not started.

- **In-Review**: Completed implementation of
  [0058 - Model reference identifiers](archive/0058-model-reference-identifiers.md) and
  advanced it from `In-Progress` to `In-Review`. Added strict Area name, Factor
  name, and Rating Level ID validation; generated JSON Schema patterns where
  JSON Schema can express the strict grammar; canonical model-reference
  render/parse helpers; revised Area Breakdown columns; updated generated
  example reports, scaffold placeholders, durable specs, runtime and durable
  `/quality` guidance, authoring guidance, README, and changelog. Verified
  `mise run check` and `go run ./cmd/qualitymd lint --json QUALITY.md`.

- **In-Progress**: Advanced
  [0058 - Model reference identifiers](archive/0058-model-reference-identifiers.md)
  from `Design` to `In-Progress`. The functional spec and design doc are
  settled; implementation begins across strict model-name lint rules, structural
  schema support, canonical model-reference helpers, Area Breakdown rendering,
  durable specs, runtime skill files, docs, and changelog.

- **Design**: Advanced
  [0058 - Model reference identifiers](archive/0058-model-reference-identifiers.md)
  from `Draft` to `Design` and added its
  [design doc](archive/0058-model-reference-identifiers/design.md). The design keeps
  `areaPath` and `factorPath` arrays as durable machine data, adds canonical
  typed reference helpers at human/tool boundaries, enforces strict local names
  through named lint rules with JSON Schema pattern support where structural
  support is possible,
  and updates the shared Area Breakdown renderer to separate Area titles from
  stable Area references. Updated the open-cases entry in the bundle
  [index](index.md). Code not started.

- **Done**: Landed and archived
  [0057 - Quality data directory](archive/0057-quality-data-directory.md) —
  advanced it to `Done` and moved the parent concept and its
  [folder](archive/0057-quality-data-directory/index.md) into
  [`archive/`](archive/index.md). The case defined shared QUALITY.md workspace
  resolution, moved evaluation and quality-log defaults under `.quality/`,
  added the root `config` tooling pointer with lint validation, updated durable
  specs, runtime skill guidance, docs, and release notes, and moved existing
  project quality data into `.quality/`.

- **In-Progress**: Advanced
  [0057 - Quality data directory](archive/0057-quality-data-directory.md) from `Design`
  to `In-Progress`. The functional spec and design doc are settled;
  implementation begins across shared workspace resolution, evaluation/status
  path defaults, lint handling for the root `config` tooling key, durable specs,
  runtime skill guidance, docs, and changelog.

- **Refinement**: Updated
  [0058 - Model reference identifiers](archive/0058-model-reference-identifiers.md) to
  explicitly list [`specs/cli/lint.md`](../specs/cli/lint.md) alongside
  [`specs/cli/lint-rules.md`](../specs/cli/lint-rules.md), so the lint command
  contract and lint rule catalog both account for strict Area name, Factor name,
  and Rating Level ID validation with named rule IDs.

- **Refinement**: Updated
  [0058 - Model reference identifiers](archive/0058-model-reference-identifiers.md) to
  use `rating:` as the canonical model-reference prefix for Rating Level IDs,
  while keeping the formal identifier term "Rating Level ID". This keeps the
  CLI/user-facing reference vocabulary aligned with Area and Factor references
  without renaming the underlying `ratingScale[].level` field.

- **Draft**: Created
  [0058 - Model reference identifiers](archive/0058-model-reference-identifiers.md)
  (`Draft`) with its
  [functional spec](archive/0058-model-reference-identifiers/spec.md) and
  [child index](archive/0058-model-reference-identifiers/index.md). The case defines
  strict Area names, Factor names, and Rating Level IDs; formal Area, Factor,
  and Rating Level IDs; canonical typed model references; edge-only shorthand;
  and clearer report summary Area Breakdown columns that separate Area title,
  stable Area reference, Area-only rating, aggregate rating, and compact Factor
  ratings. Added the case to the open-cases list in the bundle
  [index](index.md).

- **Design**: Advanced
  [0057 - Quality data directory](archive/0057-quality-data-directory.md) from `Draft`
  to `Design` and added its
  [design doc](archive/0057-quality-data-directory/design.md). The design introduces
  `internal/workspace` as the shared resolver for selected model path,
  repository root, config file, quality data directory, evaluation directory,
  and quality log directory; keeps `config` out of the normative Model and JSON
  Schema; and makes lint's unknown-key handling internally rule-configurable
  while defaulting to strict errors except for root `config`. Updated the
  open-cases entry in the bundle [index](index.md). Code not started.

- **Draft**: Created
  [0057 - Quality data directory](archive/0057-quality-data-directory.md) (`Draft`) with
  its [functional spec](archive/0057-quality-data-directory/spec.md) and
  [child index](archive/0057-quality-data-directory/index.md). The case defines the
  QUALITY.md workspace as the resolved operating context for one model file,
  uses `.quality/` as the quality data directory, moves default evaluation runs
  and the quality log under that directory, adds a root `config` pointer for
  tooling config resolution, and keeps lint strict through internally
  configurable unknown-key rule options. Added the case to the open-cases list
  in the bundle [index](index.md).

- **Done**: Landed and archived
  [0056 - Prospective evaluation plan artifacts](archive/0056-prospective-evaluation-plan-artifacts.md) —
  advanced it to `Done` and moved the parent concept and its
  [folder](archive/0056-prospective-evaluation-plan-artifacts/index.md) into
  [`archive/`](archive/index.md). The case made `design.md` and the initial
  `plan.md` prospective `/quality evaluate` artifacts authored before assessment
  begins, with later scope, coverage, rigor, or evidence-strategy changes
  recorded as plan amendments. Updated the archive [index](archive/index.md) and
  removed the entry from the open-cases list in the bundle [index](index.md).

- **In-Review**: Completed implementation of
  [0056 - Prospective evaluation plan artifacts](archive/0056-prospective-evaluation-plan-artifacts.md)
  and advanced it from `In-Progress` to `In-Review`. Added the
  [`design.md`](../specs/evaluation-records/design-md.md) artifact spec,
  clarified [`plan.md`](../specs/evaluation-records/plan-md.md) as prospective
  execution planning with explicit amendments, aligned durable `/quality`
  evaluation workflow, evaluate mode, and reporting specs, updated the runtime
  evaluate prompt and quick reference, and added the unreleased changelog entry.
  Verified `mise run fmt-md-check`, `git diff --check`, and
  `mise run npm-pack-check`.

- **In-Progress**: Advanced
  [0056 - Prospective evaluation plan artifacts](archive/0056-prospective-evaluation-plan-artifacts.md)
  from `Design` to `In-Progress`. The functional spec and design doc are
  settled; implementation begins across the durable evaluation-record specs,
  durable `/quality` evaluation workflow specs, and runtime `/quality evaluate`
  guidance.

- **Done**: Landed and archived
  [0055 - Self-describing evaluation record input](archive/0055-evaluation-input-ergonomics.md) —
  advanced it to `Done` and moved the parent concept and its
  [folder](archive/0055-evaluation-input-ergonomics/index.md) into
  [`archive/`](archive/index.md). The case made evaluation record payloads
  discoverable and validatable through payload-documenting help, no-persist
  dry-runs, aggregated key-named validation for record payloads and `plan.md`
  coverage, synced runtime and durable `/quality` skill surfaces, and added a
  published-skill relative-link package guard. Updated the archive [index](archive/index.md)
  and removed the entry from the open-cases list in the bundle [index](index.md).

- **Design**: Advanced
  [0056 - Prospective evaluation plan artifacts](archive/0056-prospective-evaluation-plan-artifacts.md)
  from `Draft` to `Design` and added its
  [design doc](archive/0056-prospective-evaluation-plan-artifacts/design.md). The design
  lands the change as a contract and prompt repair rather than a CLI behavior
  change: add a planning checkpoint after run creation and before assessment,
  split `design.md`, `plan.md`, `debug-log.md`, and formal records by job, add a
  small `design.md` artifact spec, and keep later scope or coverage changes as
  explicit `plan.md` amendments. Updated the open-cases entry in the bundle
  [index](index.md). Code not started.

- **Draft**: Created
  [0056 - Prospective evaluation plan artifacts](archive/0056-prospective-evaluation-plan-artifacts.md)
  (`Draft`) with its
  [functional spec](archive/0056-prospective-evaluation-plan-artifacts/spec.md)
  and [child index](archive/0056-prospective-evaluation-plan-artifacts/index.md). The case
  tightens `/quality evaluate` so `design.md` and the initial `plan.md` are
  authored immediately after run creation and before assessment begins, separates
  intended evidence planning from actual findings and rating rationale, and
  requires later scope or coverage changes to be explicit amendments. Added the
  case to the open-cases list in the bundle [index](index.md).

- **In-Review**: Completed implementation of
  [0055 - Self-describing evaluation record input](archive/0055-evaluation-input-ergonomics.md)
  and advanced it from `In-Progress` to `In-Review`. Added payload-documenting
  help and canonical examples for evaluation write commands, `-n/--dry-run`
  receipts that do not persist records, aggregated JSON-key validation for record
  payloads and `plan.md` coverage, seeded planned-coverage shape guidance, updated
  the runtime and durable `/quality` skill surfaces, and added a published-skill
  relative-link package guard. Verified with `mise run check`.

- **In-Progress**: Advanced
  [0055 - Self-describing evaluation record input](archive/0055-evaluation-input-ergonomics.md)
  from `Design` to `In-Progress`. The functional spec and design doc are settled;
  implementation begins across the evaluation write commands, validation path,
  status coverage checks, runtime skill surfaces, package guard, and durable
  specs/docs listed in the case.

- **Design**: Advanced
  [0055 - Self-describing evaluation record input](archive/0055-evaluation-input-ergonomics.md)
  from `Draft` to `Design` and added its
  [design doc](archive/0055-evaluation-input-ergonomics/design.md). The design splits each
  evaluation write into decode → validate → plan → commit (the seam that lets
  `-n/--dry-run` validate without persisting and report intended paths), replaces
  first-error validation with a key-named accumulator that also folds the decoder's
  unknown-field and type-mismatch errors into JSON-key vocabulary, and drift-proofs
  the surfaced payloads with one golden-tested canonical example per kind embedded
  in both help and the skill quick reference. Rejected a standalone `schema`
  command, a full example generator, and a descriptor-table validator rewrite.
  Updated the open-cases entry in the bundle [index](index.md). Code not started.

- **Draft**: Created
  [0055 - Self-describing evaluation record input](archive/0055-evaluation-input-ergonomics.md)
  (`Draft`) with its [functional spec](archive/0055-evaluation-input-ergonomics/spec.md).
  Motivated by field run `0001-quality-eval`, the case makes the `qualitymd
  evaluation` record-writing surface self-describing — payload-documenting help, a
  no-persist `-n/--dry-run`, aggregated key-named validation — and repairs the
  `/quality` skill surfaces that drifted from the binary (the unshipped
  source-of-truth citation, the stale quick-reference payloads) plus a published-
  bundle link guard. The motivating CLI design-guide additions (a new **Structured
  input** section and Help/Documentation/Errors edits) landed alongside as durable
  docs running ahead of code. Added the case to the open-cases list in the bundle
  [index](index.md).

- **Done**: Landed and archived
  [0054 - Remove improve mode](archive/0054-remove-improve-mode.md) — advanced
  it to `Done` and moved the parent concept and its
  [folder](archive/0054-remove-improve-mode/index.md) into
  [`archive/`](archive/index.md). The case removed `/quality improve` as a public
  mode, added recommendation follow-up with apply-now and issue-tracker handoff
  outcomes, updated runtime skill guidance and durable skill specs, and removed
  the improve mode files. Updated the archive [index](archive/index.md) and
  emptied the open-cases list in the bundle [index](index.md).

- **In-Review**: Completed implementation of
  [0054 - Remove improve mode](archive/0054-remove-improve-mode.md) and advanced
  it from `In-Progress` to `In-Review`. Removed runtime and durable `improve`
  mode files, added recommendation follow-up runtime and durable guidance,
  updated wizard/evaluate/update routing and quality-log ownership, updated user
  docs and examples, and verified `mise run fmt-md-check`, `git diff --check`,
  and targeted stale-reference searches.

- **In-Progress**: Advanced
  [0054 - Remove improve mode](archive/0054-remove-improve-mode.md) from `Design` to
  `In-Progress`. No design doc is required: the implementation is a mechanical
  skill/spec/doc surface change that removes the public `improve` mode and keeps
  recommendation follow-up as a non-mode workflow.

- **Design**: Advanced
  [0054 - Remove improve mode](archive/0054-remove-improve-mode.md) from `Draft` to
  `Design`. The functional spec is settled enough to work through design for
  removing the public `/quality improve` mode while preserving recommendation
  follow-up and issue-tracker handoff.

- **Creation**: Added
  [0054 - Remove improve mode](archive/0054-remove-improve-mode.md) (`status:
  Draft`) with its
  [functional spec](archive/0054-remove-improve-mode/spec.md) and
  [child index](archive/0054-remove-improve-mode/index.md). The case simplifies
  the `/quality` skill surface by removing the separate improve mode while
  preserving recommendation follow-up with apply-now and issue-tracker handoff
  outcomes. Updated the bundle [index](index.md).

- **Done**: Landed and archived
  [0053 - Align remaining durable specs](archive/0053-align-remaining-durable-specs.md) —
  advanced it to `Done` and moved the parent concept and its
  [folder](archive/0053-align-remaining-durable-specs/index.md) into
  [`archive/`](archive/index.md). The case split remaining large durable specs
  for evaluation records, lint, and ambient update notices into parent and
  component/artifact contracts. Updated the archive [index](archive/index.md) and
  emptied the open-cases list in the bundle [index](index.md).

- **In-Review**: Completed implementation of
  [0053 - Align remaining durable specs](archive/0053-align-remaining-durable-specs.md)
  and advanced it from `In-Progress` to `In-Review`. Split evaluation-records
  runtime contracts into child specs for the run folder, records, artifacts, and
  report outputs; split the lint command from lint rules and output schema; and
  split ambient update-notice behavior from the explicit update command. Updated
  affected links, indexes, and spec logs. `mise run fmt-md-check` and
  `git diff --check` pass.

- **Creation**: Added
  [0053 - Align remaining durable specs](archive/0053-align-remaining-durable-specs.md)
  (`status: In-Progress`) with its
  [functional spec](archive/0053-align-remaining-durable-specs/spec.md) and
  [child index](archive/0053-align-remaining-durable-specs/index.md). The case applies
  the revised durable-spec granularity guidance to evaluation records, lint, and
  ambient update notice behavior while keeping `SPECIFICATION.md` out of scope as
  the single primary format deliverable. Updated the bundle [index](index.md).

- **Done**: Landed and archived
  [0052 - Durable spec alignment](archive/0052-durable-spec-alignment.md) —
  advanced it to `Done` and moved the parent concept and its
  [folder](archive/0052-durable-spec-alignment/index.md) into [`archive/`](archive/index.md).
  The case aligned durable specs with artifact-spec versus behavioral-component
  guidance, added `/quality` child specs for modes, evaluation workflow,
  reporting, and quality log, narrowed the parent skill spec to shared contracts
  and links, and strengthened the general spec-splitting guidance with a heading
  inventory and fictional examples. Updated the archive [index](archive/index.md)
  and emptied the open-cases list in the bundle [index](index.md).

- **Review correction**: Reopened
  [0052 - Durable spec alignment](archive/0052-durable-spec-alignment.md) from
  `In-Review` to `In-Progress` after review found the parent `/quality` skill
  spec still retained large independently reviewable contracts. Extended the
  functional spec and affected artifacts to split the evaluation workflow,
  reporting contract, and quality log into child component specs before archive.

- **In-Review**: Completed implementation of
  [0052 - Durable spec alignment](archive/0052-durable-spec-alignment.md) and advanced it
  from `In-Progress` to `In-Review`. Added behavioral component specs under
  [`specs/skills/quality-skill/modes/`](../specs/skills/quality-skill/modes/index.md)
  for setup, wizard, evaluate, improve, and update; narrowed the parent
  [`/quality` skill spec](../specs/skills/quality-skill/quality-skill.md) to
  shared contracts plus mode summaries; updated the skill-spec
  [index](../specs/skills/quality-skill/index.md), mode [index](../specs/skills/quality-skill/modes/index.md),
  mode [log](../specs/skills/quality-skill/modes/log.md), and
  [`specs/log.md`](../specs/log.md). Reconciled the affected-artifacts list:
  no code, format spec, runtime skill files, install/scaffold files, or generated
  artifact formats changed. `mise run fmt-md-check` passes.

- **In-Progress**: Advanced
  [0052 - Durable spec alignment](archive/0052-durable-spec-alignment.md) from `Draft` to
  `In-Progress`. Its functional spec is settled and no design doc is required:
  the implementation is a mechanical durable-spec restructuring that adds
  behavioral component specs for the `/quality` modes and narrows the parent
  skill spec to shared contracts plus mode links. Updated the bundle
  [index](index.md).

- **Creation**: Added
  [0052 - Durable spec alignment](archive/0052-durable-spec-alignment.md) (`status: Draft`)
  with its [functional spec](archive/0052-durable-spec-alignment/spec.md) and
  [child index](archive/0052-durable-spec-alignment/index.md). The case aligns durable
  specs with the updated artifact-spec versus behavioral-component guidance,
  starting with child specs for the `/quality` modes while keeping 1:1 artifact
  specs named after their artifacts. Updated the bundle [index](index.md).

- **Done**: Landed and archived
  [0051 - Setup quality-md Area](archive/0051-setup-quality-md-area.md) —
  advanced from `In-Review` to `Done` and moved the parent concept and its
  [folder](archive/0051-setup-quality-md-area/index.md) into [`archive/`](archive/index.md).
  The case added the setup-authored `quality-md` Area pattern, kept
  `qualitymd init` and the CLI scaffold generic, strengthened the authoring guide
  and guide spec around quality-attribute Factor names and one referenced
  assessment across multiple Factors, and synced setup mode plus the durable
  skill spec with the concrete Area shape. Updated the archive [index](archive/index.md)
  and emptied the open-cases list in the bundle [index](index.md).

- **In-Review**: Completed implementation of
  [0051 - Setup quality-md Area](archive/0051-setup-quality-md-area.md) and advanced it
  from `In-Progress` to `In-Review`. Synced the durable skill spec and runtime
  setup mode on the concrete `quality-md` Area shape (`quality-md` key,
  `<Root Title> QUALITY.md` title, Area `description`, path-based `source`, YAML
  comments, and one guide-backed Requirement across Factors). Synced the
  authoring-guide durable spec and runtime guide on quality-attribute Factor names
  and single referenced assessments across multiple Factors. Verified the
  affected-artifacts list: no Go code, CLI scaffold, format spec, or durable docs
  were changed. `mise run fmt-md-check` passes.

- **In-Progress**: Advanced
  [0051 - Setup quality-md Area](archive/0051-setup-quality-md-area.md) from `Design` to
  `In-Progress`. Implementation is limited to the durable skill specs and bundled
  skill prompt/guide files: the CLI scaffold and Go code remain out of scope.

- **Design**: Advanced
  [0051 - Setup quality-md Area](archive/0051-setup-quality-md-area.md) from `Draft` to
  `Design` and added its [design doc](archive/0051-setup-quality-md-area/design.md).
  The design keeps `qualitymd init` generic, puts the `quality-md` Area in skill
  setup's guided population phase, uses normal path-based `source`, adds concise
  YAML comments to distinguish `source` from `assessment`, and records why one
  authoring-guide Requirement can feed multiple Factor roll-ups. Updated the
  [child index](archive/0051-setup-quality-md-area/index.md).

- **Creation**: Added
  [0051 - Setup quality-md Area](archive/0051-setup-quality-md-area.md) (`status: Draft`)
  with its [functional spec](archive/0051-setup-quality-md-area/spec.md) and
  [child index](archive/0051-setup-quality-md-area/index.md). The case proposes a
  setup-authored `quality-md` Area that evaluates the active `QUALITY.md` artifact
  itself against the active authoring guide, keeps `qualitymd init` generic, and
  strengthens the authoring guide around quality-attribute Factor names plus one
  referenced assessment connected to multiple Factors. Functional spec lists
  `specs/skills/quality-skill/quality-skill.md` and
  `specs/skills/quality-skill/guides/authoring-md.md` under **To modify**.
  Updated the bundle [index](index.md).

- **Done**: Landed and archived
  [0050 - Quality log](archive/0050-quality-log.md) — advanced from `In-Progress`
  through `In-Review` to `Done` and moved the parent concept and its
  [folder](archive/0050-quality-log/index.md) into [`archive/`](archive/index.md).
  The case added the convention-first quality log: dated `quality/log/` entries the
  `/quality` skill writes (`setup` seeds an inaugural entry, `improve` appends one
  per confirmed model change), with the format contract in `skills/quality/SKILL.md`,
  the meaningful-change taxonomy in `skills/quality/guides/authoring.md`, the
  inaugural-seed step in `modes/setup.md`, the model-change entry in `modes/improve.md`,
  and the read-only model-history/reconciliation surface in `modes/wizard.md`. Synced
  the durable `/quality` skill spec with a new `## Quality log` section and a
  deferred-CLI bullet, logged it in `specs/log.md`, and added a quality-log mention
  to `docs/guides/use-quality-skill.md`. No Go code: the `qualitymd log` CLI
  command is explicitly deferred. Updated the archive [index](archive/index.md) and
  emptied the open-cases list in the bundle [index](index.md).

- **In-Review**: Completed implementation of
  [0050 - Quality log](archive/0050-quality-log.md) and advanced it from
  `In-Progress` to `In-Review`. Reconciled the Affected artifacts list with reality
  — the only doc beyond the listed durable spec and bundled skill files was
  `docs/guides/use-quality-skill.md`, which already enumerated skill outputs and
  would have gone stale.

- **In-Progress**: Advanced [0050 - Quality log](archive/0050-quality-log.md) from
  `Draft` to `In-Progress`. The functional spec is settled and needs no design
  doc, so implementation of the convention-first quality log begins: the durable
  quality-skill spec subsection plus the bundled skill edits (`SKILL.md`,
  `guides/authoring.md`, `modes/setup.md`, `modes/improve.md`, `modes/wizard.md`).

- **Done**: Landed and archived
  [0049 - Companion JSON Schema](archive/0049-companion-json-schema.md) — advanced
  from `In-Review` to `Done` and moved the parent concept and its
  [folder](archive/0049-companion-json-schema/index.md) into
  [`archive/`](archive/index.md). The case published a structural, non-normative
  JSON Schema for QUALITY.md frontmatter (`quality.schema.json`), generated from
  `internal/schema` by `GenerateJSON()` and guarded against drift by a
  consistency test, embedded via a new root `schema.go`, and emitted by the new
  `qualitymd schema` command; it added the durable `specs/quality-schema-json.md`
  and `specs/cli/schema.md` specs. Updated the archive [index](archive/index.md)
  and removed the entry from the bundle [index](index.md).

- **In-Review**: Completed implementation of
  [0049 - Companion JSON Schema](0049-companion-json-schema.md) and advanced it
  from `In-Progress` to `In-Review`. Landed `internal/schema/jsonschema.go`
  (`GenerateJSON()`), the `internal/schema/gen` `go:generate` entrypoint, the
  generated repo-root `quality.schema.json`, the root `schema.go` embed +
  `Schema()`, the `qualitymd schema` command (verbatim plain output;
  chroma-highlighted + paged on a TTY), command registration, the chroma direct
  dependency, a root no-drift consistency test, and generator unit tests. Synced
  the durable specs/docs: new `specs/quality-schema-json.md` and
  `specs/cli/schema.md`; registered in `specs/index.md`, `specs/cli/index.md`,
  and `specs/log.md`; added the non-normative note in `SPECIFICATION.md`, the
  deferral clarification in `specs/cli/spec.md`, and the `README.md`
  quick-reference row. Build, tests, `go vet`, and `gofmt` clean; `go generate`
  idempotent; redirect round-trip byte-identical. Corrected the functional spec
  and the new artifact-spec to scope the "at least one of" `anyOf` rule to the
  Model only (Area emptiness is the warning-level, semantic `empty-area` check).
  Not committed or archived. Updated the bundle [index](index.md).

- **In-Progress**: Advanced
  [0049 - Companion JSON Schema](0049-companion-json-schema.md) from `Design` to
  `In-Progress` and began implementation. Resolved the one pending external
  input — the schema `$id` domain is `getquality.md`
  (`https://getquality.md/quality.schema.json`), matching the live docs site the
  CLI links to. Work spans `internal/schema` (JSON Schema generation), the
  generated repo-root `quality.schema.json` embedded via a new root `schema.go`,
  the `qualitymd schema` command, and the durable specs/docs being synced in
  parallel. Updated the bundle [index](index.md).

- **Done**: Landed and archived
  [0048 - Area factor report breakdown](archive/0048-area-factor-report-breakdown.md)
  — advanced from `In-Review` to `Done` and moved the parent concept and its
  [folder](archive/0048-area-factor-report-breakdown/index.md) into
  [`archive/`](archive/index.md). The case exposed a compact Area-by-Factor
  breakdown from a first-class report model across `report-summary.md`,
  `report.md`, and `report.json`, renamed the Area rating fields, and landed the
  durable `specs/reports/` artifact specs. Updated the archive
  [index](archive/index.md) and removed the entry from the bundle [index](index.md).

- **Creation**: Added [0050 - Quality log](0050-quality-log.md) (`status: Draft`)
  with its [functional spec](0050-quality-log/spec.md) and
  [child index](0050-quality-log/index.md). The case proposes a curated,
  evidence-linked **quality log** — dated entries under `quality/log/` recording
  meaningful changes to a QUALITY.md model, written by the `/quality` skill
  (`setup` seeds, `improve` appends, `wizard` reconciles drift). Convention-first:
  no `qualitymd log` CLI command or standalone artifact-spec yet. Functional spec
  lists `specs/skills/quality-skill/quality-skill.md` under **To modify**. Updated
  the bundle [index](index.md).

- **Creation**: Added
  [0049 - Companion JSON Schema](0049-companion-json-schema.md) (`status: Draft`)
  with its [functional spec](0049-companion-json-schema/spec.md) and
  [child index](0049-companion-json-schema/index.md). The case proposes a
  structural JSON Schema for QUALITY.md frontmatter — derived from
  [`internal/schema`](../internal/schema/schema.go) so it cannot drift,
  non-normative and subordinate to [`SPECIFICATION.md`](../SPECIFICATION.md) —
  plus a `qualitymd schema` verbatim-artifact command that emits it. Functional
  spec lists new durable specs `specs/quality-schema-json.md` and
  `specs/cli/schema.md`. Updated the bundle [index](index.md).

- **Design**: Advanced
  [0049 - Companion JSON Schema](0049-companion-json-schema.md) from `Draft` to
  `Design` and added its [design doc](0049-companion-json-schema/design.md).
  Decided terminal JSON highlighting uses
  [`chroma`](https://github.com/alecthomas/chroma) directly (promoted from an
  indirect dep, byte-safe on the redirect path), rejecting a glamour code-fence
  (reflows content) and a hand-rolled lipgloss tokenizer (reinvents a lexer).
  Decided the artifact lives at the repo root (`quality.schema.json`, a sibling
  of [`SPECIFICATION.md`](../SPECIFICATION.md)), embedded via a new root
  `schema.go` mirroring [`specification.go`](../specification.go) — over
  co-locating under `internal/schema/` or a dedicated `schema/` dir. Updated the
  parent artifacts (added `go.mod` chroma promotion) and the bundle
  [index](index.md).

- **Design**: Closed the remaining
  [0049 - Companion JSON Schema](0049-companion-json-schema.md) open questions.
  Generation is a `go:generate` tool writing the committed root file (the embed
  *is* the golden, guarded by a consistency test re-running an exported
  `GenerateJSON()`) over runtime generation — keeping schema changes visible as a
  reviewable diff. The schema declares JSON Schema draft 2020-12 and an
  unversioned `$id` of `https://quality.md/quality.schema.json` (identity, not
  hosting; GitHub raw-root URL as fallback if `quality.md` is not the canonical
  domain). No design questions remain; the case is ready for **In-Progress**.

- **Implementation**: Completed
  [0048 - Area factor report breakdown](0048-area-factor-report-breakdown.md)
  code in `internal/evaluation/` and advanced it to `In-Review`. Renamed the
  report-model Area rating fields to
  `areaRatingState` / `areaRatingResult` / `areaWithDescendantsRatingResult` on
  `areaSummary` and `areas`, dropped the `structural` bool and the
  structural-grouping note, added `factorRatingResults` to the compact
  `areaSummary` layer, and rendered a shared `## Area Breakdown` table (absolute
  Area display paths, path-aware Factor labels, `(area group)`, `not assessed`,
  and empty-Factor states) in `report-summary.md` and `report.md`. Strengthened
  analysis-write validation to reject duplicate and vocabulary-unresolvable
  Factor paths, added regression tests, and regenerated the three
  `0001-quality-eval` golden report fixtures. `go test ./...` green.

- **Status**: Advanced
  [0048 - Area factor report breakdown](0048-area-factor-report-breakdown.md)
  from `Design` to `In-Progress` to begin code implementation in
  `internal/evaluation/`. Durable specs and guide-spec renames already landed.

- **Design refinement**: Sharpened
  [0048 - Area factor report breakdown](0048-area-factor-report-breakdown.md)
  for long-term structure: adopted renaming the opaque Area rating fields to
  `areaRatingResult` / `areaWithDescendantsRatingResult` / `areaRatingState` and
  collapsing the redundant `structural` bool and derived note into the typed
  state (spec and design); asserted that `report.json` element arrays are the
  canonical identifiers while display paths are derived (with a separator-escaping
  non-goal); moved the guide-spec renames into a new `To rename` durable-spec
  subsection; and stated the parent status as two clocks (code vs. durable specs).

- **Durable specs**: Applied the artifact-spec filename convention to existing
  `/quality` runtime guide contracts:
  [`authoring.md`](../specs/skills/quality-skill/guides/authoring-md.md),
  [`getting-started.md`](../specs/skills/quality-skill/guides/getting-started-md.md),
  and
  [`top-10-quality-md-checks.md`](../specs/skills/quality-skill/guides/top-10-quality-md-checks-md.md).
  Updated spec links while leaving the runtime guide artifacts themselves
  unchanged.

- **Durable specs**: Added 1:1 report artifact specs for
  [0048 - Area factor report breakdown](0048-area-factor-report-breakdown.md):
  [`report-summary.md`](../specs/reports/report-summary-md.md),
  [`report.md`](../specs/reports/report-md.md), and
  [`report.json`](../specs/reports/report-json.md). Updated the specs index,
  specs log, evaluation-report command contract, and shared evaluation-records
  report-output contract to point at the new artifact specs.

- **Design refinement**: Updated
  [0048 - Area factor report breakdown](0048-area-factor-report-breakdown.md)
  to use `(area group)` as the human Markdown label for Areas with child Areas
  but no direct requirements, while preserving the typed structural/grouping
  state in machine-readable report data.

- **Design refinement**: Updated
  [0048 - Area factor report breakdown](0048-area-factor-report-breakdown.md)
  to use absolute Area display paths in the example and requirements: the root
  renders as `/ (<root title>)`, descendants render with a leading `/`, and the
  breakdown table's first column is labeled `Path`.

- **Design refinement**: Clarified the rating vocabulary for
  [0048 - Area factor report breakdown](0048-area-factor-report-breakdown.md).
  The functional spec and design now distinguish Area-only ratings from
  Area-with-descendants ratings and recommend concise Markdown labels of `Area`
  and `+ Sub-Areas` when both ratings appear in the breakdown.

- **Design refinement**: Tightened
  [0048 - Area factor report breakdown](0048-area-factor-report-breakdown.md)
  while keeping it in `Design`. The functional spec now explicitly requires
  path-aware Area and Factor labels and keeps detailed rationales in the full
  report; the design narrows cleanup to stale summary-basis helpers and records
  the `areaSummary` naming trade-off.

- **Design**: Advanced
  [0048 - Area factor report breakdown](0048-area-factor-report-breakdown.md)
  from `Draft` to `Design` and added its
  [design doc](0048-area-factor-report-breakdown/design.md). The design keeps
  `areaSummary` as the canonical compact report layer, adds Factor ratings to
  that shape, reuses the same breakdown rendering in `report-summary.md` and
  `report.md`, and strengthens new analysis writes with duplicate and
  model-aware Factor path validation. Updated the [child index](0048-area-factor-report-breakdown/index.md)
  and bundle [index](index.md).

- **Creation**: Opened
  [0048 - Area factor report breakdown](0048-area-factor-report-breakdown.md)
  (`status: Draft`) with its
  [functional spec](0048-area-factor-report-breakdown/spec.md) and
  [child index](0048-area-factor-report-breakdown/index.md). The case strengthens
  generated evaluation reports so `report-summary.md`, `report.md`, and
  `report.json` expose an at-a-glance Area-by-Factor breakdown from the
  assembled report model, with nested Area and Factor paths, structural and
  not-assessed states, and tests for stable machine identifiers. Updated the
  bundle [index](index.md).

## 2026-06-21

- **Done**: Archived
  [0047 - Area terminology changeover](archive/0047-area-terminology.md) after
  implementation and verification. Moved the parent concept and child folder
  into [`archive/`](archive/), set status to `Done`, added it to the
  [archive index](archive/index.md), and removed it from the open
  [changes index](index.md).

- **In-Review**: Completed implementation for
  [0047 - Area terminology changeover](archive/0047-area-terminology.md). The live
  schema, typed model, lint/status surfaces, evaluation records, reports, CLI
  create flag and run naming, durable specs, `/quality` skill guidance,
  scaffold, dogfood model, README/npm README, changelog, and maintained Sparrow
  example bundle now use Area/`areas:`/`areaPath` terminology while preserving
  the default `target` / `Target` rating level. Verified `go test ./...`.

- **In-Progress**: Advanced
  [0047 - Area terminology changeover](archive/0047-area-terminology.md) from `Design`
  to `In-Progress` to implement the full no-compatibility Target to Area
  changeover across the schema, evaluation records, reports, CLI, `/quality`
  skill, scaffold, dogfood model, maintained examples, and docs. Updated the
  bundle [index](index.md).

- **Design refinement**: Updated
  [0047 - Area terminology changeover](archive/0047-area-terminology.md) to keep
  `source` as the Area selector property. The
  [functional spec](archive/0047-area-terminology/spec.md) now explicitly rejects
  renaming `source` and asks prose to distinguish `source` from source code; the
  [design doc](archive/0047-area-terminology/design.md) records the rejected alternatives.

- **Design**: Advanced
  [0047 - Area terminology changeover](archive/0047-area-terminology.md) from `Draft` to
  `Design` and added its
  [design doc](archive/0047-area-terminology/design.md). The design uses a big-bang
  schema/type/record rename from Target to Area, replaces user-facing Subject
  labels with root area or model-file wording, renames evaluation-create
  `--subject` to `--model`, drops the subject altitude from new run folders, and
  guards record decoding so legacy `targetPath` records cannot be mistaken for
  root-area records. Updated the [child index](archive/0047-area-terminology/index.md)
  and bundle [index](index.md).

- **Creation**: Opened
  [0047 - Area terminology changeover](archive/0047-area-terminology.md)
  (`status: Draft`) with its
  [functional spec](archive/0047-area-terminology/spec.md) and
  [child index](archive/0047-area-terminology/index.md). The case replaces the formal
  Target model-node vocabulary with Area, introduces root area as the formal root
  descriptor, rejects legacy `targets:` / `targetPath` compatibility, and scopes
  the change across schema, records, reports, CLI, skill, scaffold, examples,
  and docs. Updated the bundle [index](index.md).

- **Done**: Archived
  [0046 - Evaluation debug log](archive/0046-evaluation-debug-log.md) after
  implementation and verification. Moved the parent concept and child folder
  into [`archive/`](archive/), set status to `Done`, added it to the
  [archive index](archive/index.md), and removed it from the open
  [changes index](index.md).

- **In-Review**: Completed implementation and durable artifact synchronization
  for [0046 - Evaluation debug log](archive/0046-evaluation-debug-log.md). New
  evaluation runs seed a process-only `debug-log.md`; the record specs, CLI
  create contract, `/quality` skill guidance, reference fixture, docs,
  changelog, and skill compatibility metadata now preserve the boundary between
  evaluation-process events and formal subject-quality evidence. Verified
  `go test ./...` and `mise run check`. Updated the bundle [index](index.md).

- **In-Progress**: Advanced
  [0046 - Evaluation debug log](archive/0046-evaluation-debug-log.md) from `Design` to
  `In-Progress`. Implementation will seed `debug-log.md` in evaluation runs,
  update the runtime and CLI contracts, align the `/quality` skill guidance, and
  refresh tests, examples, and release notes. Updated the bundle
  [index](index.md).

- **Design**: Advanced
  [0046 - Evaluation debug log](archive/0046-evaluation-debug-log.md) from `Draft` to
  `Design` and added its
  [design doc](archive/0046-evaluation-debug-log/design.md). The design seeds a small
  run-root `debug-log.md` through `qualitymd evaluation create`, keeps report
  assembly independent of debug prose, and puts the process-only boundary in the
  `/quality` skill guidance. Updated the bundle [index](index.md).

- **Creation**: Opened
  [0046 - Evaluation debug log](archive/0046-evaluation-debug-log.md) (`status: Draft`)
  with its [functional spec](archive/0046-evaluation-debug-log/spec.md) and
  [child index](archive/0046-evaluation-debug-log/index.md). The case adds a
  process-only `debug-log.md` artifact to evaluation runs while keeping
  assessments, analysis, recommendations, and reports authoritative for
  subject-quality judgment. Updated the bundle [index](index.md).

- **Done**: Archived completed change cases
  [0042 - Typed report model](archive/0042-typed-report-model.md),
  [0043 - Evaluation history compatibility](archive/0043-evaluation-history-compatibility.md),
  [0044 - Section unknowns and open questions](archive/0044-section-unknowns-open-questions.md),
  and [0045 - Evaluable body context](archive/0045-evaluable-body-context.md)
  after review. Moved each parent concept and child folder into
  [`archive/`](archive/), set their statuses to `Done`, added them to the
  [archive index](archive/index.md), and removed them from the open
  [changes index](index.md).

- **In-Review**: Completed implementation and durable artifact synchronization
  for [0045 - Evaluable body context](archive/0045-evaluable-body-context.md). The
  authoring guide, format spec, guide contracts, getting-started/top-10/setup
  guidance, README summary, and scaffold now treat the Markdown body as
  evaluable, agent-accessible judgment context; the scaffold test asserts the
  new marker. Reviewed the dogfood `QUALITY.md` and active eval model for
  concrete access-gap fallout. Verified `go test ./...` and `mise run check`.
  Updated the [child index](archive/0045-evaluable-body-context/index.md) and bundle
  [index](index.md).

- **In-Progress**: Advanced
  [0045 - Evaluable body context](archive/0045-evaluable-body-context.md) from `Draft`
  through `Design` (no design doc needed) to `In-Progress`. Implementation will
  update the authoring guide, its durable spec contract, body-context checks,
  and any scaffold or setup guidance needed to treat the Markdown body as
  evaluable, agent-accessible judgment context. Updated the [child index](archive/0045-evaluable-body-context/index.md)
  and bundle [index](index.md).

- **Creation**: Opened
  [0045 - Evaluable body context](archive/0045-evaluable-body-context.md)
  (`status: Draft`) with its
  [functional spec](archive/0045-evaluable-body-context/spec.md) and
  [child index](archive/0045-evaluable-body-context/index.md). The case clarifies that
  the Markdown body is evaluable judgment context for building, justifying, and
  evaluating model quality; that body sections should be concise,
  self-explanatory, and progressively disclosed; and that material support that
  is not agent-accessible is a first-class limitation captured in the relevant
  section's unknowns or open questions. Updated the bundle [index](index.md).

## 2026-06-20

- **In-Review**: Completed
  [0044 - Section unknowns and open questions](archive/0044-section-unknowns-open-questions.md).
  Retired the standalone Known gaps body section in favor of per-section unknowns,
  open questions, and a human/agent review state line across the format spec,
  authoring guide, `init` scaffold (and its test), skill setup/getting-started/
  top-10 checks, the durable specs, the example fixtures, and the dogfood
  `QUALITY.md` and active eval model. Verified `go test ./...` and `mise run check`.
  Updated the [child index](archive/0044-section-unknowns-open-questions/index.md) and
  bundle [index](index.md).

- **In-Progress**: Created and advanced
  [0044 - Section unknowns and open questions](archive/0044-section-unknowns-open-questions.md)
  from `Draft` through `Design` (no design doc needed) to `In-Progress`. The case
  replaces the standalone Known gaps body section with a common per-section shape,
  per-section unknowns and open questions, and a human/agent review state line,
  propagating across the format spec, authoring guide, scaffold, skill checks, and
  dogfood instances. Added the [functional spec](archive/0044-section-unknowns-open-questions/spec.md),
  [child index](archive/0044-section-unknowns-open-questions/index.md), and bundle
  [index](index.md) entry.

- **In-Review**: Completed implementation and durable artifact synchronization
  for [0043 - Evaluation history compatibility](archive/0043-evaluation-history-compatibility.md).
  Evaluation history inspection now surfaces malformed, unsupported, and
  incomplete historical records as typed non-reportable gaps; list/status/latest
  workflows remain usable; report build/gate refuse incompatible selected runs
  with status-oriented diagnostics; and the `/quality` skill guidance treats
  incompatible records as history status rather than subject quality evidence.
  Verified `go test ./...` and `mise run check`. Updated the [child index](archive/0043-evaluation-history-compatibility/index.md)
  and bundle [index](index.md).

- **In-Progress**: Advanced
  [0043 - Evaluation history compatibility](archive/0043-evaluation-history-compatibility.md)
  from `Design` to `In-Progress` to implement tolerant evaluation-history
  inspection, compatibility gaps, and graceful report/list/status behavior.
  Updated the [child index](archive/0043-evaluation-history-compatibility/index.md) and
  bundle [index](index.md).

- **Design**: Advanced
  [0043 - Evaluation history compatibility](archive/0043-evaluation-history-compatibility.md)
  from `Draft` to `Design` and added its
  [design doc](archive/0043-evaluation-history-compatibility/design.md). The design uses
  a tolerant run-inspection layer for status/list/history commands, records
  incompatible files as reportability gaps, keeps record writers strict, and
  gates report build/gate through compatibility status before trusted report
  assembly. Updated the [child index](archive/0043-evaluation-history-compatibility/index.md)
  and bundle [index](index.md).

- **Creation**: Opened
  [0043 - Evaluation history compatibility](archive/0043-evaluation-history-compatibility.md)
  (`status: Draft`) with its
  [functional spec](archive/0043-evaluation-history-compatibility/spec.md) and
  [child index](archive/0043-evaluation-history-compatibility/index.md). The case
  captures the strict-writer / tolerant-reader posture for evaluation history:
  historical or hand-edited runs can become non-reportable compatibility gaps
  without breaking ordinary status, list, latest-run, or fresh-evaluation
  workflows. Updated the bundle [index](index.md).

- **In-Review**: Completed implementation and durable artifact synchronization
  for [0042 - Typed report model](archive/0042-typed-report-model.md). Evaluation reports
  now use typed rating-result, local-rating, next-step, lifecycle,
  missing-metadata, rigor, evaluation-level, path, and gap concepts; report JSON
  exposes explicit state objects; existing invalid rating/severity records become
  non-reportable gaps; and the Sparrow fixture reports were regenerated.
  Verified `mise run check`. Updated the bundle [index](index.md).

- **In-Progress**: Opened
  [0042 - Typed report model](archive/0042-typed-report-model.md) to replace stringly
  typed and implicit evaluation-report states with explicit typed concepts for
  rating results, local target ratings, next steps, lifecycle state, run gaps,
  rigor, evaluation level, missing metadata, and path identities. Added the
  parent case, functional spec, design doc, and child index; updated the bundle
  [index](index.md).

- **Done**: Set status `Done` and archived
  [0041 - Update command and improvements](archive/0041-update-command.md) after
  publishing and verifying `v0.5.0`. Moved the parent concept and child folder
  into [`archive/`](archive/), added the entry to the [archive index](archive/index.md),
  and removed it from the open [changes index](index.md).

- **In-Review**: Completed implementation and durable artifact synchronization
  for [0041 - Update command and improvements](archive/0041-update-command.md). The CLI
  now exposes apply-by-default `qualitymd update` with `--check`, readiness and
  release-notes fields, managed standalone apply, post-apply version
  verification, update-check opt-out, and a cached ambient notice. Renamed the
  `/quality` maintenance mode and durable skill/spec/docs references to
  `update`. Verified `mise run check`, a Windows compile-only check for
  `internal/cli`, and CLI smoke checks for `update` and removed `upgrade`.
  Updated the bundle [index](index.md).

- **In-Progress**: Advanced
  [0041 - Update command and improvements](archive/0041-update-command.md) from `Design`
  to `In-Progress` to begin implementation of the apply-by-default
  `qualitymd update` command, ambient cached update notice, and paired
  `/quality update` skill-mode rename. Updated the bundle [index](index.md).

- **Re-characterization**: Re-characterized
  [0041 - Update command and improvements](archive/0041-update-command.md) as the
  upgrade→update rename plus its improvements, dropping the earlier framing and
  renaming the case from slug `0041-codex-aligned-update` to `0041-update-command`
  (parent, child folder, and the same-day entries below repointed to the new
  path). Expanded scope to rename the paired `/quality upgrade` skill mode to
  `/quality update`: the [functional spec](archive/0041-update-command/spec.md) gains a
  paired skill-mode-rename requirement and a durable-spec change for
  `specs/skills/quality-skill/quality-skill.md`, and the parent's
  **Affected artifacts** now lists the skill spec, the runtime
  `skills/quality/modes/upgrade.md` → `update.md` rename, `SKILL.md` routing,
  `wizard.md`, the CLI quick reference, and the top-10-checks route token. Updated
  the bundle [index](index.md).

- **Redesign**: Reshaped 0041 to an apply-by-default `update` command and renamed
  the case to [0041 - Update command and improvements](archive/0041-update-command.md)
  (slug `0041-upgrade-apply-and-readiness` → `0041-update-command`; earlier
  entries below repointed to the new path). Per the chosen direction, the
  [functional spec](archive/0041-update-command/spec.md) and
  [design doc](archive/0041-update-command/design.md) now rename `upgrade`→`update`
  with apply-by-default and a `--check` advisory (deprecated `upgrade` alias for
  one cycle), and add an ambient cached update notice on ordinary commands. The
  notice deliberately reverses 0032's "ordinary commands MUST NOT check the
  network" rule; it is fenced by strict rails — stderr only, never in
  stdout/`--json`, suppressed off a terminal, in CI, behind
  `QUALITYMD_NO_UPDATE_CHECK`, and for dev builds — served from a cache under
  `$QUALITYMD_HOME` refreshed by a detached, non-blocking subprocess. The managed
  standalone self-apply, readiness gating, and release-notes reference carry
  forward onto the new command shape. Expands the affected-artifact footprint
  (the durable `specs/cli/upgrade.md` is renamed to `specs/cli/update.md`;
  `specs/cli.md`, versioning docs, and the `/quality` skill files all change).
  Updated the bundle [index](index.md).

- **Design**: Advanced
  [0041 - Upgrade self-apply, readiness, and release notes](archive/0041-update-command.md)
  from `Draft` to `Design` and added its
  [design doc](archive/0041-update-command/design.md). The design lands all
  three deltas inside the existing `internal/cli/upgrade.go` seams: widen
  `latestVersionProvider` to return a `{version, ready, releaseNotesURL}` struct
  (so readiness and notes ride the single injectable, offline-testable network
  call); resolve `Ready` from the `assets[]`/`html_url` already in the GitHub
  `releases/latest` response (npm's registry latest is ready by definition); gate
  reported availability and `--apply` on readiness; and add managed standalone to
  the `applySupported`/`upgradeCommand` tables, invoking the existing idempotent
  installer non-interactively via `QUALITYMD_NO_INPUT=1`. Records the Homebrew
  latest-provider quirk and a possible `releaseReady` JSON field as open
  questions. Updated the [child index](archive/0041-update-command/index.md)
  and the bundle [index](index.md).

- **Creation**: Opened
  [0041 - Upgrade self-apply, readiness, and release notes](archive/0041-update-command.md)
  (`status: Draft`) with its [functional spec](archive/0041-update-command/spec.md)
  and [child index](archive/0041-update-command/index.md). The case captures
  three improvements drawn from comparing `qualitymd upgrade` with a conventional
  CLI update flow: extend `--apply` to self-update managed standalone installs (the
  channel the project owns yet 0032 left unable to apply), gate "update
  available" and `--apply` on the target release actually being downloadable, and
  surface a release-notes reference in advisory and `--json` output. Records
  [`internal/cli/upgrade.go`](../internal/cli/upgrade.go),
  [`specs/cli/upgrade.md`](../specs/cli/upgrade.md), the `/quality` upgrade-mode
  skill files, and versioning docs as affected. Added the open-case entry to the
  bundle [index](index.md).

## 2026-06-19

- **Done**: Set status `Done` and archived
  [0040 - Readable report summary](archive/0040-readable-report-summary.md).
  Moved the parent concept and child folder into [`archive/`](archive/), added
  the entry to the [archive index](archive/index.md), and removed it from the
  open [changes index](index.md). Verified `go test ./...`, targeted
  `dprint check`, and `git diff --check`.

- **In-Review**: Completed implementation and durable artifact synchronization
  for [0040 - Readable report summary](archive/0040-readable-report-summary.md). The
  summary renderer now emits the key-details, Summary, Top Issues,
  Recommendations, and Scope & Limitations outline; uses "Full evaluation" and
  "Overall rating" in human-facing summary output; surfaces copyable
  Recommendation IDs; and keeps `report.json` unchanged. Updated durable report
  specs, the `/quality` skill contract/runtime wording, tests, and the worked
  summary example. Verified `go test ./...` and targeted `dprint check`.

- **Implementation**: Advanced
  [0040 - Readable report summary](archive/0040-readable-report-summary.md)
  from `Draft` to `In-Progress` and added its
  [design doc](archive/0040-readable-report-summary/design.md). The design keeps the
  existing `EvaluationReportDocument` and JSON schema, reshaping only the concise
  Markdown renderer into a decision-brief outline with display-time wording for
  "Full evaluation" and "Overall rating".

- **Creation**: Opened
  [0040 - Readable report summary](archive/0040-readable-report-summary.md)
  (`status: Draft`) with its child
  [index](archive/0040-readable-report-summary/index.md) and
  [functional spec](archive/0040-readable-report-summary/spec.md). The spec proposes the
  revised `report-summary.md` outline: key details, Summary, Top Issues,
  Recommendations, and Scope & Limitations; updates human-facing labels to
  "Full evaluation" and "Overall rating"; and makes active Recommendation IDs
  prominent for follow-up prompts.

- **Done**: Set status `Done` and archived
  [0039 — Evaluation command surface redesign](archive/0039-evaluation-command-surface.md).
  Moved the parent concept and child folder into [`archive/`](archive/), fixed
  archive-relative links, added the entry to the [archive index](archive/index.md),
  and removed it from the open [changes index](index.md). Verified `go test
  ./...`, `go vet ./...`, targeted `dprint check`, and CLI smoke checks for the
  new and removed evaluation command surfaces.

- **Design**: Reconciled the
  [0039 — Evaluation command surface redesign](archive/0039-evaluation-command-surface.md)
  impact list and renamed its section from "Affected specs & docs" to **Affected
  artifacts**. Added an **Affected code** subsection (the `internal/cli` command
  tree, the `internal/evaluation/*` backends incl. `planned_coverage.go`, and
  `internal/status`) and the previously-missing artifacts: the skill spec
  `specs/skills/quality-skill/quality-skill.md`, plus `skills/quality/SKILL.md`,
  `install.md`, `docs/guides/use-quality-skill.md`, `docs/guides/cli-design.md`,
  `docs/guides/write-functional-specs.md`, and `CHANGELOG.md`.

- **Design**: Advanced
  [0039 — Evaluation command surface redesign](archive/0039-evaluation-command-surface.md)
  from `Draft` to `Design` and added its
  [design doc](archive/0039-evaluation-command-surface/design.md): the cobra command
  tree, shared `resolveRun`/payload-batching/output-stream helpers, the
  `plan.md`-folded coverage with read-time validation, the report `build`/`gate`
  split, and the altitude removal — with rejected alternatives and three open
  design calls (flat vs subfolder specs, `list --state` scope, malformed-coverage
  gap name). Updated the [changes index](index.md) and child
  [index](archive/0039-evaluation-command-surface/index.md).

- **Creation**: Opened
  [0039 — Evaluation command surface redesign](archive/0039-evaluation-command-surface.md)
  (`status: Draft`) with its child [index](archive/0039-evaluation-command-surface/index.md)
  and [functional spec](archive/0039-evaluation-command-surface/spec.md). The spec sets a
  single noun/verb rule for the `qualitymd evaluation` surface, renames the
  run-lifecycle verbs, promotes the record kinds and the report to nouns with
  honest verbs, adds run/record `list`, folds planned coverage into `plan.md`
  frontmatter (deleting `set-planned-coverage` and `planned-coverage.json`),
  separates `report gate` from `report build`, and removes the altitude residue.
  Added the case to the open [changes index](index.md).

- **Done**: Set status `Done` and archived
  [0038 — /quality skill interaction UX](archive/0038-quality-skill-interaction-ux.md).
  Moved the parent concept and child folder into [`archive/`](archive/), fixed
  repo-relative links for the deeper path, added the entry to the
  [archive index](archive/index.md), and removed it from the open
  [changes index](index.md).

- **In-Review**: Completed implementation and durable artifact synchronization
  for [0038 — /quality skill interaction UX](archive/0038-quality-skill-interaction-ux.md).
  Added the durable `User interaction contract` section to the `/quality` skill
  spec, added compact shared interaction rules to root `SKILL.md`, and updated
  wizard/evaluate/improve/setup/upgrade mode prompts for run frames, decision
  briefs, stop/reroute behavior, history context, improvement delta reporting,
  and status-first output. Verified targeted Markdown formatting with
  `dprint check`. The case is ready to archive per the requested goal.

- **Implementation**: Advanced change
  [0038 — /quality skill interaction UX](archive/0038-quality-skill-interaction-ux.md)
  from `Design` to `In-Progress` so the settled interaction contract can be
  implemented in the durable `/quality` skill spec and runtime skill files.

- **Design**: Advanced change
  [0038 — /quality skill interaction UX](archive/0038-quality-skill-interaction-ux.md)
  from `Draft` to `Design` and added its
  [design doc](archive/0038-quality-skill-interaction-ux/design.md). The design adds a
  durable `User interaction contract` section to the `/quality` skill spec,
  keeps shared run-frame and decision-brief shapes compact in root `SKILL.md`,
  applies the behavior at mode boundaries, uses existing status/evaluation
  history surfaces rather than new storage, and keeps improvement delta reports
  as human output rather than a new evaluation artifact.

- **Creation**: Added change
  [0038 — /quality skill interaction UX](archive/0038-quality-skill-interaction-ux.md)
  in `Draft` with its
  [functional spec](archive/0038-quality-skill-interaction-ux/spec.md). The change
  proposes a durable interaction contract for the `/quality` skill covering run
  frames, decision briefs, stop/reroute behavior, history-aware operation,
  improvement delta reports, and status-first output, while keeping the existing
  skill/CLI boundary and evaluation artifact format intact. Updated the bundle
  [index](index.md).

- **Done**: Set status `Done` and archived
  [0036 — Harden install scripts and upgrade idiomatics](archive/0036-harden-install-scripts.md)
  and
  [0037 — Render rating-level titles in evaluation reports](archive/0037-report-rating-titles.md).
  Moved both parent concepts and child folders into [`archive/`](archive/),
  fixed repo-relative links for the deeper path, added the entries to the
  [archive index](archive/index.md), and removed them from the open
  [changes index](index.md).

- **In-Review**: Completed implementation and durable artifact synchronization
  for [0037 — Render rating-level titles in evaluation reports](archive/0037-report-rating-titles.md).
  Human Markdown reports now resolve rating labels through the run's rating-scale
  titles with a level-id fallback, while `report.json`, `BuildResult`, and
  `--fail-at-or-below` continue using level ids. Added emoji rating titles to
  `QUALITY.md`, clarified the durable build-report spec, and updated evaluation
  tests for title rendering, JSON id preservation, fallback behavior, and
  non-rating states. Verified `go test ./...`, `go vet ./...`, and
  `dprint check`. The case remains open in [`changes/`](index.md) for review; it
  is not archived until it lands.

- **Implementation**: Advanced change
  [0037 — Render rating-level titles in evaluation reports](archive/0037-report-rating-titles.md)
  from `Draft` to `In-Progress`. No design doc is needed for this localized
  renderer change, so implementation can begin from the settled
  [functional spec](archive/0037-report-rating-titles/spec.md).

- **Creation**: Added change
  [0037 — Render rating-level titles in evaluation reports](archive/0037-report-rating-titles.md)
  in `Draft` with its
  [functional spec](archive/0037-report-rating-titles/spec.md). The change makes the human
  reports (`report.md`, `report-summary.md`) display each rating level's `title`
  instead of its `level` id — bringing the renderer into conformance with the
  existing build-report SHOULD so emoji-bearing titles read in reports — while
  keeping `level` ids in `report.json`, `BuildResult`, and the
  `--fail-at-or-below` gate, and dogfoods emoji titles in `QUALITY.md`. Omits a
  design doc. Updated the bundle [index](index.md).

- **In-Review**: Completed implementation and durable artifact synchronization
  for [0036 — Harden install scripts and upgrade idiomatics](archive/0036-harden-install-scripts.md).
  Added the three-tool checksum fallback with a non-silent skip and the
  print-the-export-line PATH guidance to `install/install.sh`; added the TLS 1.2
  shim, per-user PATH mutation, and `-NonInteractive` gating to
  `install/install.ps1`; made `updateAvailable` SemVer-correct via
  `golang.org/x/mod/semver` in `internal/cli/upgrade.go` with regression
  coverage; commented the intentional Homebrew cask in `.goreleaser.yaml`; and
  synced the durable upgrade spec, install docs, contributor guide, and release
  guide. Verified `go test ./...`, `go vet`, `golangci-lint`, `gofmt`,
  `shellcheck install/install.sh`, `dprint check`, and `goreleaser check`. The
  case remains open in [`changes/`](index.md) for review; it is not archived
  until it lands.

- **Implementation**: Advanced change
  [0036 — Harden install scripts and upgrade idiomatics](archive/0036-harden-install-scripts.md)
  from `Design` to `In-Progress` so the settled installer and upgrade-check fixes
  can be implemented and synced into the durable upgrade spec and install docs.

- **Design**: Advanced change
  [0036 — Harden install scripts and upgrade idiomatics](archive/0036-harden-install-scripts.md)
  from `Draft` to `Design` and added its
  [design doc](archive/0036-harden-install-scripts/design.md). The design settles a
  three-tool checksum fallback with a non-silent skip, a `-bor` TLS 1.2 shim for
  PowerShell 5.1, a deliberately asymmetric PATH model (per-user PATH mutation on
  Windows, print-the-export-line on Unix), `--non-interactive` as a verbosity
  gate rather than a phantom prompt, SemVer-correct update detection via
  `golang.org/x/mod/semver`, and keeping the Homebrew cask with documented
  rationale (rejecting the deprecated formula path).

- **Creation**: Added change
  [0036 — Harden install scripts and upgrade idiomatics](archive/0036-harden-install-scripts.md)
  in `Draft` with its
  [functional spec](archive/0036-harden-install-scripts/spec.md). The change fixes five
  portability/convention gaps surfaced by an install-surface review — dead
  checksum verification on stock Linux, a missing TLS 1.2 pin on Windows
  PowerShell, absent/asymmetric PATH integration, a string-compare standing in
  for a SemVer update check, and a no-op `--non-interactive` flag — and records
  that the Homebrew **cask** (not a formula) is the idiomatic distribution path
  after the "convert to formula" review item was investigated and reversed.
  Updated the bundle [index](index.md).

- **Done**: Set status `Done` and archived
  [0035 — /quality upgrade mode](archive/0035-quality-skill-upgrade-mode.md).
  Moved the parent concept and child folder into [`archive/`](archive/), fixed
  repo-relative links for the deeper path, added the entry to the
  [archive index](archive/index.md), and removed it from the open
  [changes index](index.md).

- **In-Review**: Completed implementation and durable artifact synchronization
  for [0035 — /quality upgrade mode](archive/0035-quality-skill-upgrade-mode.md). Added
  the runtime `skills/quality/modes/upgrade.md` procedure, routed `upgrade` from
  `SKILL.md`, taught wizard to recommend it for stale/incompatible skill/CLI
  state, updated the durable `/quality` skill spec, documented the existing
  install maintenance flow, and verified targeted Markdown formatting. The case
  remains open in [`changes/`](index.md) for review; it is not archived until it
  lands.

- **Implementation**: Advanced change
  [0035 — /quality upgrade mode](archive/0035-quality-skill-upgrade-mode.md) from
  `Design` to `In-Progress` so the settled upgrade-mode spec and design can be
  implemented and synced into the durable skill spec, runtime skill files, and
  install/versioning docs.

- **Design**: Advanced change
  [0035 — /quality upgrade mode](archive/0035-quality-skill-upgrade-mode.md) from
  `Draft` to `Design` and added its
  [design doc](archive/0035-quality-skill-upgrade-mode/design.md). The design adds a
  mode-specific upgrade procedure that snapshots skill and CLI versions, builds a
  plan before mutation, delegates CLI changes to `qualitymd upgrade`, delegates
  skill changes to the Agent Skills installer when available, verifies the
  visible CLI afterward, and warns that skill upgrades may require a restarted
  agent session.

- **Creation**: Added change
  [0035 — /quality upgrade mode](archive/0035-quality-skill-upgrade-mode.md) in `Draft`
  with its [functional spec](archive/0035-quality-skill-upgrade-mode/spec.md). The
  change proposes a skill mode that checks the installed `/quality` skill and
  `qualitymd` CLI pair, diagnoses compatibility and available updates, plans
  skill and CLI upgrade actions, asks before mutation, delegates mechanics to
  the Agent Skills installer and `qualitymd upgrade`, and reports any required
  agent restart or reload. Updated the bundle [index](index.md).

- **Done**: Set status `Done` and archived
  [0031 — Evaluation report summary artifact](archive/0031-report-summary-artifact.md),
  [0032 — CLI managed upgrades](archive/0032-cli-managed-upgrades.md),
  [0033 — Required display titles](archive/0033-required-display-titles.md), and
  [0034 — Skill release metadata](archive/0034-skill-release-metadata.md).
  Moved each parent concept and child folder into [`archive/`](archive/), updated
  the root and archive indexes, and left no open change cases.

- **In-Review**: Completed implementation and durable artifact synchronization
  for [0034 — Skill release metadata](archive/0034-skill-release-metadata.md). The case
  remains open in [`changes/`](index.md) for review; it is not archived until it
  lands.

- **Implementation**: Advanced
  [0034 — Skill release metadata](archive/0034-skill-release-metadata.md) from `Design`
  to `In-Progress` so its settled metadata and release-check design can be
  implemented.

- **In-Review**: Completed implementation and durable artifact synchronization
  for [0031 — Evaluation report summary artifact](archive/0031-report-summary-artifact.md),
  [0032 — CLI managed upgrades](archive/0032-cli-managed-upgrades.md), and
  [0033 — Required display titles](archive/0033-required-display-titles.md). The cases
  remain open in [`changes/`](index.md) for review; they are not archived until
  they land.

- **Design**: Advanced change
  [0034 — Skill release metadata](archive/0034-skill-release-metadata.md) from `Draft`
  to `Design` and added its
  [design doc](archive/0034-skill-release-metadata/design.md). The design uses
  Agent Skills `metadata.version` and `metadata.requires-qualitymd-cli`, mirrors
  the range in `compatibility`, adds release-check validation against the tag and
  changelog, updates runtime/docs wording, and leaves installer enforcement for a
  future package contract.

- **Creation**: Added change
  [0034 — Skill release metadata](archive/0034-skill-release-metadata.md) in `Draft`
  with its [functional spec](archive/0034-skill-release-metadata/spec.md). The change
  proposes project-owned Agent Skills metadata in `skills/quality/SKILL.md` for
  the `/quality` skill SemVer and required `qualitymd` CLI range, mirrored by
  `compatibility` prose and curated release notes, with release-check validation
  and installer enforcement explicitly deferred. Updated the bundle
  [index](index.md).

- **Implementation**: Advanced changes
  [0031 — Evaluation report summary artifact](archive/0031-report-summary-artifact.md),
  [0032 — CLI managed upgrades](archive/0032-cli-managed-upgrades.md), and
  [0033 — Required display titles](archive/0033-required-display-titles.md) to
  `In-Progress` so their settled specs/designs can be implemented and synced
  into durable specs, docs, tests, and examples. Updated the bundle
  [index](index.md).

- **Creation**: Added change
  [0033 — Required display titles](archive/0033-required-display-titles.md) in `Draft`
  with its [functional spec](archive/0033-required-display-titles/spec.md). The change
  proposes required `title` properties on the Model, every Target, every Factor,
  and every Rating Level; adds `Factor.title`; keeps Requirements title-free;
  makes `missing-title` an error across those nodes; and records the affected
  format, lint, init, report, status, skill, README, guide, scaffold, and example
  updates. Updated the bundle [index](index.md).

- **Done**: Set status `Done` and archived
  [0030 — CLI status snapshot command](archive/0030-cli-status-command.md).
  Moved the parent concept and child folder into [`archive/`](archive/), fixed
  repo-relative links for the deeper path, added the entry to the
  [archive index](archive/index.md), and removed it from the open
  [changes index](index.md).

- **Status**: Advanced change
  [0030 — CLI status snapshot command](archive/0030-cli-status-command.md) to
  `In-Review` after adding `qualitymd status [path] [--json]`, the
  `internal/status` snapshot assembler, evaluation helpers for run listing and
  active recommendation counts, CLI tests and status-package tests, durable CLI
  specs, README command docs, and `/quality` skill updates. Verified targeted
  Markdown formatting, `go test ./...`, `mise run vet`, and a smoke run of
  `go run ./cmd/qualitymd status --json`.

- **Implementation**: Advanced change
  [0030 — CLI status snapshot command](archive/0030-cli-status-command.md) from `Design`
  to `In-Progress` so the settled status-command spec and design can be
  implemented and synced into the durable CLI docs and `/quality` skill
  consumers.

- **Design**: Advanced change
  [0030 — CLI status snapshot command](archive/0030-cli-status-command.md) from `Draft`
  to `Design` and added its
  [design doc](archive/0030-cli-status-command/design.md). The design introduces an
  `internal/status` snapshot assembler, keeps CLI rendering thin, reuses lint and
  evaluation mechanics, compares run `model.md` snapshots for staleness, counts
  active recommendations through evaluation-owned helpers, and keeps report-body
  scraping out of the command.

- **Draft**: Replaced the placeholder for
  [0030 — CLI status snapshot command](archive/0030-cli-status-command.md) with a full
  [functional spec](archive/0030-cli-status-command/spec.md). The spec defines the
  read-only `qualitymd status [path] [--json]` invocation, lint validity and
  model-shape snapshot, source coverage, evaluation history and staleness
  signals, active recommendation counts, readiness states, deterministic
  next-action data, and exit behavior. Updated the case and bundle listings.

- **Design**: Advanced change
  [0031 — Evaluation report summary artifact](archive/0031-report-summary-artifact.md)
  from `Draft` to `Design` and added its
  [design doc](archive/0031-report-summary-artifact/design.md). The design reuses the
  existing `ReportJSON` assembly as the single report model, adds a
  `renderReportSummaryMarkdown` projection, extends `BuildResult` with the
  summary path, and keeps `report-summary.md` generated from the same recorded
  run data as `report.md` and `report.json`. Updated the bundle [index](index.md).

- **Design**: Advanced change
  [0032 — CLI managed upgrades](archive/0032-cli-managed-upgrades.md) from `Draft` to
  `Design` and added its
  [design doc](archive/0032-cli-managed-upgrades/design.md). The design stages the work
  through structured version metadata, install-context detection, explicit
  upgrade checks, guarded apply behavior, and GitHub-hosted managed installer
  entrypoints under top-level `install/`.

- **Creation**: Added change
  [0032 — CLI managed upgrades](archive/0032-cli-managed-upgrades.md) in `Draft` with
  its [functional spec](archive/0032-cli-managed-upgrades/spec.md). The change proposes
  structured version metadata, explicit upgrade checks, safe install-method
  detection, advisory output by default, guarded `--apply` behavior, npm launcher
  marking, and a long-term managed standalone installer path. Records affected
  CLI specs, install/versioning docs, release guidance, npm launcher, and
  `/quality` skill consumers. Updated the bundle [index](index.md).

- **Creation**: Added change
  [0031 — Evaluation report summary artifact](archive/0031-report-summary-artifact.md)
  in `Draft` with its
  [functional spec](archive/0031-report-summary-artifact/spec.md). The change proposes
  generating `report-summary.md` beside `report.md` and `report.json` during
  `qualitymd evaluation build-report`, derived from the same recorded run data
  and summary layer, for PR/CI/stakeholder triage without replacing the full
  Evaluation Report. Records affected CLI specs, evaluation record contract,
  skill reporting spec, README, and example bundles. Updated the bundle
  [index](index.md).

- **Done**: Set status `Done` and archived
  [0028 — Require factor references](archive/0028-require-characterized-requirements.md)
  and
  [0029 — Sharpen assessment references and traceability](archive/0029-sharpen-assessment-references.md).
  Moved each parent concept and child folder into [`archive/`](archive/),
  fixed repo-relative links for the deeper path, added both entries to the
  [archive index](archive/index.md), and left
  [0030 — CLI status snapshot command](archive/0030-cli-status-command.md) as the only
  open change case.

- **Creation**: Queued change
  [0030 — CLI status snapshot command](archive/0030-cli-status-command.md) in `Draft`
  with a placeholder [functional spec](archive/0030-cli-status-command/spec.md). The case
  proposes a read-only `qualitymd status [--json]` command that emits a
  deterministic project-state snapshot (model validity and shape, evaluation run
  history, open recommendation counts) so the `/quality` wizard routes from
  structured data instead of hand-parsing `QUALITY.md` and reading report bodies —
  restoring the CLI-owns-mechanical-work split. Records the affected CLI specs,
  README, and `/quality` skill consumers. Spec is a placeholder until the case is
  picked up. Updated the bundle [index](index.md).

- **Status**: Advanced change
  [0029 — Sharpen assessment references and traceability](archive/0029-sharpen-assessment-references.md)
  to `In-Review` after extending `SPECIFICATION.md`'s **Assessment** terminology
  and Requirement section with the inline-or-reference framing and a non-normative
  traceability note, applying the six authoring-guide edits (reserve "source" for
  `Target.source`, traceability-graph job, entity gloss, target/assessment-
  reference duality, split-by-claim job, and the renamed "Reference an external
  assessment" job), nudging the scaffold to "reference" wording, and verifying the
  Go test suite, vet, and Markdown formatting.

- **Implementation**: Advanced change
  [0029 — Sharpen assessment references and traceability](archive/0029-sharpen-assessment-references.md)
  from `Draft` to `In-Progress` (no design doc) so the durable `SPECIFICATION.md`,
  authoring guide, and scaffold edits can be made from the settled spec.

- **Creation**: Added change
  [0029 — Sharpen assessment references and traceability](archive/0029-sharpen-assessment-references.md)
  in `Draft` with its
  [functional spec](archive/0029-sharpen-assessment-references/spec.md). The change frames
  a requirement's `assessment` as either stated inline or a reference to another
  entity, reserves "source" for `Target.source`, extends the "reference"
  terminology 0028 set for factors to the requirement→entity edge, and makes the
  model's traceability graph an authoring concern — across `SPECIFICATION.md`, the
  authoring guide, and the scaffold, with no schema or lint change. Omits a design
  doc. Updated the bundle [index](index.md).

## 2026-06-18

- **Status**: Advanced change
  [0028 — Require factor references](archive/0028-require-characterized-requirements.md)
  to `In-Review` after adding the `missing-factor-reference` lint error,
  updating factor-reference terminology, syncing durable specs/docs/scaffold, and
  verifying the Go test suite.

- **Implementation**: Advanced change
  [0028 — Require factor references](archive/0028-require-characterized-requirements.md)
  from `Design` to `In-Progress` so the settled lint rule, terminology updates,
  durable specs, README, and scaffold guidance can be implemented.

- **Alignment**: Brought change
  [0028 — Require factor references](archive/0028-require-characterized-requirements.md)
  into alignment with the current change-case guides by adding the required
  [Durable spec changes](archive/0028-require-characterized-requirements/spec.md#durable-spec-changes)
  section to its functional spec and moving durable-doc accounting out of its
  design doc.

- **Done**: Set status `Done` and archived
  [0026 — Authoring guide replaces meta-model workflow](archive/0026-authoring-guide-remove-meta-model.md)
  and [0027 — Modularize quality skill modes](archive/0027-modularize-quality-skill.md).
  Moved each parent concept and child folder into [`archive/`](archive/),
  updated repo-relative links for the deeper path, added both entries to the
  [archive index](archive/index.md), and left
  [0028 — Require factor references](archive/0028-require-characterized-requirements.md)
  as the only open change case.

- **Template**: Added a required `## Durable spec changes` section (**To add** /
  **To modify** / **To delete**, each a list or `None`) to the example template
  [spec](archive/0001-example-change/spec.md), so copies account for the durable
  specs a change rewrites. See
  [Writing functional specs](../docs/guides/write-functional-specs.md#durable-spec-changes).

- **Design**: Advanced change
  [0028 — Require factor references](archive/0028-require-characterized-requirements.md)
  from `Draft` to `Design` and added its
  [design doc](archive/0028-require-characterized-requirements/design.md). The design
  keeps `factors` structurally optional, adds a context-sensitive
  `missing-factor-reference` lint error for direct target-level requirements
  without factor references, renames secondary-factor internals to neutral
  factor-reference wording, and records why `missing-factor` and a schema-level
  required `factors` property were rejected.

- **Creation**: Added change
  [0028 — Require factor references](archive/0028-require-characterized-requirements.md)
  in `Draft` with its
  [functional spec](archive/0028-require-characterized-requirements/spec.md). The change
  makes requirements invalid unless they reference at least one factor, keeps
  "lens" available as shorthand, and distinguishes direct target-level
  `factors` references from secondary factors on requirements already nested
  under a factor. Updated the bundle [index](index.md).

- **Schema migration**: Renamed the `changes/` parent concept type from
  `Change` to `Change Case`, updated existing parent concepts and the
  [changes schema](schema.md), renamed the contributor guide to
  [Working with change cases](../docs/guides/work-with-change-cases.md), and
  narrowed `AGENTS.md` so routine prompted edits do not require a Change Case.

- **Status**: Advanced change
  [0027 — Modularize quality skill modes](archive/0027-modularize-quality-skill.md) to
  `In-Review` after keeping `SKILL.md` as the root router/global contract,
  adding setup, wizard, evaluate, and improve mode files under
  `skills/quality/modes/`, renaming supporting skill docs to `resources/`,
  syncing the durable skill spec, and verifying the test suite.

- **Implementation**: Added change
  [0027 — Modularize quality skill modes](archive/0027-modularize-quality-skill.md) in
  `In-Progress` with its
  [functional spec](archive/0027-modularize-quality-skill/spec.md).
  The change keeps `SKILL.md` as the `/quality` router and moves setup, wizard,
  evaluate, and improve procedures into separate files under
  `skills/quality/modes/`, with supporting docs under `skills/quality/resources/`.

- **Status**: Advanced change
  [0026 — Authoring guide replaces meta-model workflow](archive/0026-authoring-guide-remove-meta-model.md)
  to `In-Review` after replacing the skill-facing meta-model reference with
  [authoring.md](../skills/quality/guides/authoring.md),
  removing the bundled `models` CLI/package, making evaluation run creation
  subject-only, syncing durable specs and docs, and verifying the Go test suite.

- **Implementation**: Added change
  [0026 — Authoring guide replaces meta-model workflow](archive/0026-authoring-guide-remove-meta-model.md)
  in `In-Progress` with its
  [functional spec](archive/0026-authoring-guide-remove-meta-model/spec.md) and
  [design doc](archive/0026-authoring-guide-remove-meta-model/design.md). The change
  replaces the skill-facing meta-model reference with an authoring guide
  authoring guide, removes the public `qualitymd models` / model-altitude
  workflow, and syncs durable specs and docs around subject-only evaluation.

- **Done**: Set status `Done` and archived the full in-review set —
  [0012 — Evaluation record format](archive/0012-evaluation-record-format.md),
  [0013 — Evaluation run scaffold](archive/0013-evaluation-run-scaffold.md),
  [0014 — Evaluation record write](archive/0014-evaluation-record-write.md),
  [0015 — Evaluation status and report build](archive/0015-evaluation-report-build.md),
  [0016 — Skill consumes evaluation CLI](archive/0016-skill-consume-eval-cli.md),
  [0017 — Skill rigor and efficiency](archive/0017-skill-rigor-efficiency.md),
  [0018 — Evaluation report UX](archive/0018-evaluation-report-ux.md),
  [0019 — Duplicate assessment status](archive/0019-duplicate-assessment-status.md),
  [0020 — Planned coverage status](archive/0020-planned-coverage-status.md),
  [0021 — Recommendation superseding](archive/0021-recommendation-superseding.md),
  [0022 — Create-run subject validation](archive/0022-create-run-subject-validation.md),
  [0023 — Assessment superseding](archive/0023-assessment-superseding.md),
  [0024 — Report regression coverage](archive/0024-report-regression-coverage.md),
  and [0025 — Durable spec rationale](archive/0025-durable-spec-rationale.md).
  Moved each parent concept and its child folder into
  [`archive/`](archive/), fixed their repo-relative links for the deeper path,
  added them to the [archive index](archive/index.md), and emptied the
  open-changes [index](index.md).

- **Status**: Advanced change
  [0025 — Durable spec rationale](archive/0025-durable-spec-rationale.md) to `In-Review`
  after teaching the three durable contributor guides to keep rationale in the
  spec: a **Background / Motivation** shape entry and per-requirement `Rationale:`
  annotation convention (with litmus and say-it-once rule) in
  [write-functional-specs.md](../docs/guides/write-functional-specs.md), the
  rewritten two-whys convention and refined rationale smells there, the
  absorb-the-why step gated on **Before setting In-Review** in
  [work-with-change-cases.md](../docs/guides/work-with-change-cases.md), and the
  rationale-is-promoted note in
  [write-design-docs.md](../docs/guides/write-design-docs.md). Recorded the guide
  edits in the [docs log](../docs/log.md).

- **Implementation**: Advanced change
  [0025 — Durable spec rationale](archive/0025-durable-spec-rationale.md) from `Design`
  to `In-Progress` so the three durable contributor guides
  ([write-functional-specs.md](../docs/guides/write-functional-specs.md),
  [work-with-change-cases.md](../docs/guides/work-with-change-cases.md),
  [write-design-docs.md](../docs/guides/write-design-docs.md)) can be edited from
  the settled spec and design.

- **Design**: Advanced change
  [0025 — Durable spec rationale](archive/0025-durable-spec-rationale.md) from `Draft`
  to `Design` and added its
  [design doc](archive/0025-durable-spec-rationale/design.md). The design settles a
  two-layer, co-located in-spec rationale convention — a Background/Motivation
  section plus subordinate per-requirement `Rationale:` annotations governed by
  a litmus and a say-it-once rule — over the rejected alternatives (a separate
  Diátaxis explanation doc, design-intent-only depth, and a full ADR embedded in
  the spec), with spec bloat the headline risk mitigated by keeping the
  requirement the lead sentence.

- **Creation**: Added change
  [0025 — Durable spec rationale](archive/0025-durable-spec-rationale.md) in `Draft`
  with its [functional spec](archive/0025-durable-spec-rationale/spec.md). The change
  targets the contributor guides: durable specs inherit a requirement when a
  case archives but lose the case's motivation and the design doc's
  rationale, so editors re-litigate settled lessons and "simplify" rules back
  into the bugs they fixed. The spec states a two-layer in-spec rationale
  convention — a spec-level Background/Motivation section and per-requirement
  `Rationale:` annotations — plus the litmus for when to annotate and an
  absorb-the-why step on landing, and dogfoods the convention itself. Updated
  the bundle [index](index.md).

- **Refinement**: Folded the E49 TypeScript SDK recommendation-quality finding
  into [0017 — Skill rigor and efficiency](archive/0017-skill-rigor-efficiency.md):
  recommendations should name inferable route hints such as affected package,
  path, workflow, maintainer surface, or verification route in existing text
  fields rather than adding a schema field.

- **Status**: Advanced change
  [0024 — Report regression coverage](archive/0024-report-regression-coverage.md) to
  `In-Review` after adding focused temp-run tests for secret-style,
  prompt-injection-style, not-assessed, dotted-path, structural-root, and
  empty-recommendation report behavior.

- **Implementation**: Advanced change
  [0024 — Report regression coverage](archive/0024-report-regression-coverage.md) from
  `Design` to `In-Progress` to turn repeated report-rendering experiment findings
  into focused automated tests without committing benchmark fixture snapshots.

- **Design**: Advanced change
  [0024 — Report regression coverage](archive/0024-report-regression-coverage.md) from
  `Draft` to `Design` and added its
  [design doc](archive/0024-report-regression-coverage/design.md). The design builds
  temporary evaluation runs in tests and asserts high-risk rendered `report.md`
  and `report.json` properties without committing benchmark fixture snapshots.

- **Creation**: Added change
  [0024 — Report regression coverage](archive/0024-report-regression-coverage.md) in
  `Draft` after the experiment program repeatedly found report-rendering
  regressions around seeded safety cases, prompt-injection handling,
  not-assessed propagation, dotted-path limitation extraction, structural roots,
  and empty recommendation arrays.

- **Status**: Advanced change
  [0023 — Assessment superseding](archive/0023-assessment-superseding.md) to
  `In-Review` after implementing assessment `supersedes` metadata,
  superseding status gaps, active/superseded report rendering, durable specs,
  and skill guidance.

- **Implementation**: Advanced change
  [0023 — Assessment superseding](archive/0023-assessment-superseding.md) from `Design`
  to `In-Progress` to close the remaining correction-workflow gap after
  recommendation superseding. The change requires analyses to reference active
  corrected assessments rather than superseded records.

- **Design**: Advanced change
  [0023 — Assessment superseding](archive/0023-assessment-superseding.md) from `Draft`
  to `Design` and added its
  [design doc](archive/0023-assessment-superseding/design.md). The design adds an
  optional `supersedes` list to assessment records, validates superseding and
  stale-analysis references in status, and renders active versus superseded
  assessments while requiring analyses to reference active records.

- **Creation**: Added change
  [0023 — Assessment superseding](archive/0023-assessment-superseding.md) in `Draft`
  after recommendation superseding (E28) left no ergonomic way to correct an
  assessment inside a run while keeping analysis roll-ups bound to active
  judgment.

- **Status**: Advanced change
  [0022 — Create-run subject validation](archive/0022-create-run-subject-validation.md)
  to `In-Review` after validating subject paths before run-folder creation,
  syncing the durable create-run spec, and verifying the failed `--subject .`
  scenario leaves no partial evaluation artifacts.

- **Implementation**: Advanced change
  [0022 — Create-run subject validation](archive/0022-create-run-subject-validation.md)
  from `Design` to `In-Progress` after the E14/E29 CLI UX finding showed that
  invalid `create-run --subject` values can fail after creating an empty run
  skeleton.

- **Design**: Advanced change
  [0022 — Create-run subject validation](archive/0022-create-run-subject-validation.md)
  from `Draft` to `Design` and added its
  [design doc](archive/0022-create-run-subject-validation/design.md). The design
  validates the subject path before creating the evaluation directory or run
  folder so invalid subjects consume no run number and leave no partial
  artifacts.

- **Creation**: Added change
  [0022 — Create-run subject validation](archive/0022-create-run-subject-validation.md)
  in `Draft` after the E14 improve/re-evaluate experiment found that `qualitymd
  evaluation create-run --subject .` failed after creating an empty run
  skeleton.

- **Status**: Advanced change
  [0021 — Recommendation superseding](archive/0021-recommendation-superseding.md) to
  `In-Review` after implementing recommendation `supersedes` metadata,
  dangling-reference status gaps, active/superseded report rendering, durable
  specs, and skill guidance.

- **Implementation**: Advanced change
  [0021 — Recommendation superseding](archive/0021-recommendation-superseding.md) from
  `Design` to `In-Progress` so the CLI record schema, status checks, report
  rendering, durable specs, and skill guidance can be updated from the settled
  spec and design.

- **Design**: Advanced change
  [0021 — Recommendation superseding](archive/0021-recommendation-superseding.md) from
  `Draft` to `Design` and added its
  [design doc](archive/0021-recommendation-superseding/design.md). The design uses an
  optional `supersedes` list on recommendation records, validates dangling
  references in status, and keeps superseded recommendations visible while
  choosing Next Action from active recommendations.

- **Creation**: Added change
  [0021 — Recommendation superseding](archive/0021-recommendation-superseding.md) in
  `Draft` after the E15 recommendation-correction trial showed that append-only
  correction records can leave the report's primary Next Action pointing at a
  stale recommendation.

- **Status**: Advanced change
  [0020 — Planned coverage status](archive/0020-planned-coverage-status.md) to
  `In-Review` after implementing `qualitymd evaluation set-planned-coverage`,
  planned-coverage status gaps, durable specs/docs, and skill prompt guidance.

- **Implementation**: Advanced change
  [0020 — Planned coverage status](archive/0020-planned-coverage-status.md) from
  `Design` to `In-Progress` so the CLI writer, status checks, durable specs,
  and skill prompt can be implemented from the settled spec and design.

- **Design**: Advanced change
  [0020 — Planned coverage status](archive/0020-planned-coverage-status.md) from
  `Draft` to `Design` and added its
  [design doc](archive/0020-planned-coverage-status/design.md). The design uses an
  optional CLI-owned `planned-coverage.json` artifact plus status set
  comparisons so interrupted runs can identify missing planned assessments,
  missing planned analyses, and unexpected records without changing current
  behavior for runs that omit the artifact.

- **Creation**: Added change
  [0020 — Planned coverage status](archive/0020-planned-coverage-status.md) in `Draft`
  after the E11 interruption/resume experiment and planned-coverage prototype
  showed that `show-status` can report missing analysis but cannot name missing
  planned assessments without structured planned coverage metadata.

- **Status**: Advanced change
  [0019 — Duplicate assessment status](archive/0019-duplicate-assessment-status.md) to
  `In-Review` after implementing duplicate-assessment renderability checks,
  syncing durable specs, and verifying the command-boundary duplicate trial.

- **Implementation**: Advanced change
  [0019 — Duplicate assessment status](archive/0019-duplicate-assessment-status.md) from
  `Design` to `In-Progress` to implement the `duplicate-assessment` renderability
  gap, sync durable specs, and update skill guidance from the settled design.

- **Design**: Advanced change
  [0019 — Duplicate assessment status](archive/0019-duplicate-assessment-status.md) from
  `Draft` to `Design` and added its
  [design doc](archive/0019-duplicate-assessment-status/design.md). The design detects
  assessment records that share a target path and requirement, reports them as
  `duplicate-assessment` gaps, and fails report rendering through the existing
  renderability gate.

- **Creation**: Added change
  [0019 — Duplicate assessment status](archive/0019-duplicate-assessment-status.md) in
  `Draft` after the experiment program found that re-adding a corrected
  assessment appends a conflicting duplicate record while status still reports
  the run as reportable.

- **Refinement**: Tightened reportability for change
  [0015 — Evaluation status and report build](archive/0015-evaluation-report-build.md).
  `show-status` now reports `missing-root-analysis` unless exactly one analysis
  record has an empty `targetPath`, and `build-report` refuses to render instead
  of silently using a child target as the headline.

- **Refinement**: Extended change
  [0018 — Evaluation report UX](archive/0018-evaluation-report-ux.md) to read report
  context from bounded `design.md` and `plan.md` sections before falling back to
  folder names or rationale text. This makes scope descriptions and planned
  limitations populate `report.md` and `report.json` directly when the skill
  records them.

- **Refinement**: Updated change
  [0018 — Evaluation report UX](archive/0018-evaluation-report-ux.md) after the ESLint
  standard-run experiment to deduplicate equivalent summary limitations and
  accept the skill's `Scope description` / `Narrowing` resolved-parameter
  labels.

- **Refinement**: Updated change
  [0018 — Evaluation report UX](archive/0018-evaluation-report-ux.md) after seeded
  stability repeats to preserve dotted file paths when deriving limitation
  summaries from recorded rationales.

- **Status**: Advanced change
  [0018 — Evaluation report UX](archive/0018-evaluation-report-ux.md) from
  `In-Progress` to `In-Review` after implementing summary-first report
  rendering, syncing durable specs/docs and the skill prompt, and verifying the
  renderer on copied ESLint and DataLoader runs.

- **Implementation**: Advanced change
  [0018 — Evaluation report UX](archive/0018-evaluation-report-ux.md) from `Design` to
  `In-Progress` so the report renderer, durable specs/docs, and skill prompt can
  be updated from the settled functional spec and design.

- **Design**: Created change
  [0018 — Evaluation report UX](archive/0018-evaluation-report-ux.md) in `Design`.
  The change turns the experiment-backed V1 report-shape recommendation into a
  functional spec and design doc for summary-first `report.md`, clearer
  `report.json`, explicit scope/limitations/evidence basis, grouping-target
  rendering, and stable empty recommendation arrays.

- **Status**: Advanced changes
  [0012](archive/0012-evaluation-record-format.md),
  [0013](archive/0013-evaluation-run-scaffold.md),
  [0014](archive/0014-evaluation-record-write.md),
  [0015](archive/0015-evaluation-report-build.md),
  [0016](archive/0016-skill-consume-eval-cli.md), and
  [0017](archive/0017-skill-rigor-efficiency.md) from `In-Progress` to `In-Review`
  after implementing the evaluation CLI surface and syncing the listed durable
  specs/docs.

- **Status**: Advanced changes
  [0012](archive/0012-evaluation-record-format.md),
  [0013](archive/0013-evaluation-run-scaffold.md),
  [0014](archive/0014-evaluation-record-write.md),
  [0015](archive/0015-evaluation-report-build.md),
  [0016](archive/0016-skill-consume-eval-cli.md), and
  [0017](archive/0017-skill-rigor-efficiency.md) from `Design` to `In-Progress` so code
  and durable specs/docs are now phase-authorized.

- **Refinement**: Adopted implementation-readiness review fixes across the open
  changes. Updated lifecycle wording to gate durable spec/doc sync before
  `In-Review`, aligned affected-doc sections with the current `Design` phase,
  renamed the planned durable `show-status` spec path to
  `specs/cli/evaluation-show-status.md`, clarified `show-status` gap semantics,
  kept recommendation Markdown bodies CLI-rendered from structured payload fields,
  expanded durable-doc coverage for the `/quality` skill spec and reference
  examples, and fixed the design-guide link below.

- **Design**: Advanced the evaluation-workflow set and the skill rigor pass from
  `Draft` to `Design`, adding a
  [design doc](../docs/guides/write-design-docs.md) to each (drafted in parallel,
  one per change). The designs settle the *how* against the settled specs:
  - [0012 — Evaluation record format](archive/0012-evaluation-record-format.md) - the
    contract lives as one enduring bundle-root concept `specs/evaluation-records.md`
    (not under `cli/`), registered the normal OKF way; the skill's switch from
    inlined prose to a reference is deferred to In-Progress.
  - [0013 — Evaluation run scaffold](archive/0013-evaluation-run-scaffold.md) - a new
    `evaluation` Cobra group with `create-run`, a thin CLI over a new
    `internal/evaluation` package, collision-fixing numbering by a single
    directory scan (max+1, create-or-fail), and three-level
    `--evaluation-dir` → `.quality/config.yaml` → `quality/evaluations/`
    resolution.
  - [0014 — Evaluation record write](archive/0014-evaluation-record-write.md) - the
    `internal/evaluation` package owns the record layer; three subcommands share
    one decode→validate→place→write pipeline with strict `DisallowUnknownFields`
    rejection of CLI-owned fields, stateless scan-based numbering with one retry,
    and deterministic recommendation Markdown.
  - [0015 — Evaluation status and report build](archive/0015-evaluation-report-build.md) -
    a typed run-read layer with a shared `Renderable` predicate so green
    `show-status` guarantees `build-report`; one in-memory report renders both
    byte-deterministic files (Glamour kept out of the write path), and
    `--fail-at-or-below` reuses the existing coded-exit mechanism.
  - [0016 — Skill consumes evaluation CLI](archive/0016-skill-consume-eval-cli.md) - the
    skill's evaluation flow maps onto the four commands, judgment JSON piped over
    stdin (CLI stamps/numbers), the inlined Artifact Contract replaced by a
    reference, and the prerequisite probe extended to the evaluation commands.
  - [0017 — Skill rigor and efficiency](archive/0017-skill-rigor-efficiency.md) - rigor
    rules enforced as durable artifact constraints (applied breadth recorded in
    `plan.md`, evidence/locator rigor on existing fields with no schema bump,
    a rating-binding re-check that re-runs the verifying command, compute-then-batch
    writes, and confined `deep` fan-out).
    Updated the bundle [index](index.md). Durable `specs/`/skill/README sync still
    happens per change at In-Progress.

- **Refinement**: Adopted the full explicit verb-object evaluation command
  surface and the CLI-guideline follow-ups. Renamed run creation to
  `qualitymd evaluation create-run`, report rendering to
  `qualitymd evaluation build-report`, added
  `qualitymd evaluation show-status`, renamed the gate flag to
  `--fail-at-or-below`, added `--evaluation-dir` to run creation, added `--file`
  input support to `add-record`, specified atomic numbered writes with one
  recompute retry, and fixed deterministic recommendation rendering order.

- **Refinement**: Renamed the record writer command to
  `qualitymd evaluation add-record assessment|analysis|recommendation`, keeping
  record writes in the evaluation namespace while making the object explicit in
  the command name. Updated dependent wording in
  [0013](archive/0013-evaluation-run-scaffold.md),
  [0014](archive/0014-evaluation-record-write.md),
  [0015](archive/0015-evaluation-report-build.md), and
  [0016](archive/0016-skill-consume-eval-cli.md), plus the bundle index.

- **Refinement**: Kept recommendation artifacts as human-readable Markdown while
  making them first-class CLI-written runtime records. The evaluation-record
  contract now says `recommendations/*.md` carries runtime YAML frontmatter with
  `schemaVersion: 1` and machine-readable metadata, without making the run folder
  an OKF bundle. Change
  [0014 — Evaluation record write](archive/0014-evaluation-record-write.md) now includes
  `qualitymd evaluation add-record recommendation <run>`, and dependent report/skill
  wording reads recommendation records mechanically instead of hand-authoring
  them.

- **Refinement**: Revised change
  [0014 — Evaluation record write](archive/0014-evaluation-record-write.md) and its
  dependents to place record writes under the evaluation namespace:
  `qualitymd evaluation add-record assessment|analysis|recommendation` instead of a separate
  `qualitymd result add` top-level command. Updated dependent wording in
  [0013](archive/0013-evaluation-run-scaffold.md),
  [0015](archive/0015-evaluation-report-build.md), and
  [0016](archive/0016-skill-consume-eval-cli.md), plus the bundle index.

- **Creation**: Added a coordinated set of six changes to sharpen the evaluation
  workflow, drafted in parallel and consolidated here. The seam: the
  deterministic `qualitymd` CLI **writes** every evaluation record (numbering,
  schema stamping, report rendering) while the skill supplies **judgment**.
  - [0012 — Evaluation record format](archive/0012-evaluation-record-format.md)
    (`Draft`) - the keystone: move the artifact contract from the skill prompt
    into an enduring `specs/` spec both the CLI and skill consume.
  - [0013 — Evaluation run scaffold](archive/0013-evaluation-run-scaffold.md) (`Draft`) -
    `qualitymd evaluation create-run`, deterministic shared run numbering (fixes
    a real collision), seeds `model.md`/`design.md`/`plan.md`.
  - [0014 — Evaluation record write](archive/0014-evaluation-record-write.md) (`Draft`) -
    `qualitymd evaluation add-record assessment|analysis|recommendation`,
    JSON judgment payload from `--file` or stdin, inherent validation, atomic
    rejection.
  - [0015 — Evaluation status and report build](archive/0015-evaluation-report-build.md) (`Draft`) -
    `qualitymd evaluation show-status` and `qualitymd evaluation build-report`,
    idempotent rendering of `report.md`/`report.json` from records,
    `--fail-at-or-below` CI gate; renders ratings, never infers them.
  - [0016 — Skill consumes evaluation CLI](archive/0016-skill-consume-eval-cli.md)
    (`Draft`) - the skill drives the CLI instead of hand-authoring run folders,
    records, or reports; replaces its inlined Artifact Contract with a reference.
  - [0017 — Skill rigor and efficiency](archive/0017-skill-rigor-efficiency.md)
    (`Draft`) - operationalized effort levels, verified evidence and pinned
    locators, rating-binding re-check, batched writes, optional deep fan-out;
    independent of the CLI work.
    Updated the bundle [index](index.md). Each spec is the change's *delta*;
    durable `specs/`/skill/README sync happens per change at In-Progress.

- **Done**: Archived change
  [0011 — CLI human output polish](archive/0011-cli-human-output-polish.md)
  after styling `models list`, adding lint next actions to JSON and human
  output, making dev version output include a short VCS revision when available,
  syncing the durable CLI specs, and adding focused output-gate tests.

- **Implementation**: Created and advanced change
  [0011 — CLI human output polish](archive/0011-cli-human-output-polish.md) to
  `In-Progress`. The change covers terminal styling for `models list`, lint
  next actions in JSON and human output, dev version output that includes a VCS
  revision, and broader output-gate tests, with durable CLI specs to be synced
  during implementation.

## 2026-06-17

- **Done**: Archived change
  [0010 — Implement the /quality skill](archive/0010-implement-quality-skill.md)
  after landing the skill artifact, `qualitymd models` CLI surface, durable
  specs/docs, raw runtime example bundle, and verification. Removed it from the
  open-changes index and added it to the archive index.

- **Implementation**: Implemented change
  [0010 — Implement the /quality skill](archive/0010-implement-quality-skill.md):
  added
  `skills/quality/SKILL.md`, implemented `qualitymd models list/view` with Markdown
  and JSON output plus `--source`, moved the bundled quality meta-model under
  `internal/models`, added skill-first install/docs, synced the durable CLI and
  skill specs, re-captured the example as raw runtime artifacts with JSON
  assessment/analysis/report files, ignored default dogfood runs under
  `quality/evaluations/`, and verified the skill/CLI surfaces locally.

- **Implementation**: Advanced change
  [0010 — Implement the /quality skill](archive/0010-implement-quality-skill.md) from
  `Design` to `In-Progress`. The functional spec and design doc are settled, so
  implementation files and durable specs/docs can now be updated: the
  `skills/quality/SKILL.md` artifact, `qualitymd models` CLI surface,
  `.quality/config.yaml` behavior, raw JSON example artifacts, and related durable
  documentation.

- **Refinement**: Added a comprehensive acceptance checklist to change
  [0010 — Implement the /quality skill](archive/0010-implement-quality-skill.md), covering
  skill packaging/install, CLI prerequisite handling including dev builds,
  `qualitymd models` Markdown/JSON behavior, `.quality/config.yaml` validation,
  default dogfood-output ignoring, quick model-altitude dogfooding, JSON artifact
  parsing and shape, minimal `report.json` finding summaries, example re-capture,
  and durable spec/doc sync before **Done**. Optional installer/UI metadata such as
  `agents/openai.yaml` is explicitly deferred until the installer or target agent
  docs require it.

- **Refinement**: Settled the final `SKILL.md` description text for change
  [0010 — Implement the /quality skill](archive/0010-implement-quality-skill.md): "Use when
  a user wants setup, wizard guidance, evaluation, or improvement for quality
  management of a project/entity or one of its components/areas. Trigger for
  requests about quality factors, characteristics, attributes, criteria, areas,
  factors, requirements, improving a quality factor such as
  security/reliability/usability, evaluating a subject against quality criteria, or
  evaluating/improving the QUALITY.md model itself."

- **Refinement**: Added evaluation-directory configuration to change
  [0010 — Implement the /quality skill](archive/0010-implement-quality-skill.md). The skill
  now reads repository-local `.quality/config.yaml` with `evaluationDir` to choose
  the parent directory for numbered evaluation runs, defaulting to
  `quality/evaluations/` when absent; the config is framed as shared qualitymd
  system config that future CLI evaluation commands should also honor. The path must
  be repository-relative, normalized, and unable to escape the repository; unknown
  keys are warned and ignored. Broader configuration ideas (default target file,
  effort, output formats, thresholds, retention, install commands) are deferred
  until they have a concrete need.

- **Refinement**: Added trigger-description requirements for change
  [0010 — Implement the /quality skill](archive/0010-implement-quality-skill.md). The skill
  description must now cover quality management/evaluation/improvement prompts even
  when the user does not mention `QUALITY.md` (for example, improving security
  quality), include mode trigger terms (`setup`, `wizard`, `evaluate`, `improve`)
  and generic quality vocabulary (factors, characteristics, attributes, criteria),
  and stay bounded away from generic copyediting or one-off "higher quality"
  requests. The design records initial `SKILL.md` description text that names
  Targets, Factors, Requirements, subject evaluation, model evaluation/improvement,
  and broad project/entity plus component/target quality, while leaving CLI
  implementation details to the skill body. The design now records the criteria
  used to derive that description and the durable spec sync explicitly carries those
  criteria into `quality-skill.md`'s Frontmatter and metadata section.

- **Refinement**: Added dogfooding guidance to change
  [0010 — Implement the /quality skill](archive/0010-implement-quality-skill.md). The design
  now requires an In-Progress verification pass that installs the skill from the
  working tree, accepts a local development `qualitymd` binary when it exposes the
  required commands, runs a quick model-altitude evaluation against this repo's
  `QUALITY.md`, checks the generated artifact shape, and avoids committing ad hoc
  `quality/evaluations/` output unless deliberately re-captured as a durable
  example.

- **Refinement**: Resolved the remaining open questions for change
  [0010 — Implement the /quality skill](archive/0010-implement-quality-skill.md). Root
  `install.md` now uses a verification-first install flow (`qualitymd --version`,
  documented CLI install/upgrade, verify again, `npx skills add qualitymd/quality.md`,
  `npx skills list`) with the exact package-manager command filled in when the first
  public CLI release channel exists. `report.json` now inlines only minimal generic
  finding summaries by reference for single-file gate/dashboard consumers, while
  full finding detail remains in `assessments/*.json`. Future
  `qualitymd outline QUALITY.md --json`, CLI evaluation record/gate commands,
  and deep-run subagent fanout are explicitly deferred rather than open design
  questions.

- **Refinement**: Generalized the structured finding shape for change
  [0010 — Implement the /quality skill](archive/0010-implement-quality-skill.md). Replaced
  the sample's bespoke top-level `credentialType` with a generic finding object:
  `locator`, `observation`, open `category`, optional `severity`, supporting
  `evidence`, and optional `attributes` for domain-specific metadata. Added the
  broader JSON-shape rule that public top-level fields stay tied to the evaluation
  workflow, while factor- or requirement-specific details live under
  `attributes`.

- **Refinement**: Updated change
  [0010 — Implement the /quality skill](archive/0010-implement-quality-skill.md) so
  `qualitymd models view <name>` supports `--json`. The default output remains
  Markdown with the same terminal-rendered vs plain/verbatim split as
  `qualitymd spec`, preserving byte-for-byte `model.md` snapshots while giving
  humans a readable TTY view. `--json` emits the same source-rewritten bundled model
  as structured data (`model` plus `bodyMarkdown`) for agents and gates that should
  not have to reparse Markdown/YAML. Updated the functional spec, design doc, and
  design log wording.

- **Refinement**: Corrected the onboarding model for change
  [0010 — Implement the /quality skill](archive/0010-implement-quality-skill.md). The skill
  is now the primary entry point, installed from this repo with
  `npx skills add qualitymd/quality.md`; the `qualitymd` CLI is a prerequisite that
  `setup` and `wizard` detect, version-check, and help install or upgrade before
  running CLI-dependent work. Updated the functional spec, design doc, parent
  change, and indexes to remove the plugin-first assumption, added root
  `install.md` to affected docs, and kept Claude plugin/marketplace packaging as a
  possible secondary channel rather than this change's primary distribution path.

- **Design**: Advanced change
  [0010 — Implement the /quality skill](archive/0010-implement-quality-skill.md) from
  `Draft` to `Design` and added its
  [design doc](archive/0010-implement-quality-skill/design.md). Confirmed the three
  **blocking** open items at their recommended resolutions and worked out the *how*:
  the skill ships from `skills/quality/SKILL.md` as an Agent Skills artifact
  installable with `npx skills add qualitymd/quality.md`, while `setup`/`wizard`
  verify the separate `qualitymd` CLI prerequisite before doing CLI-dependent work.
  Specified the `qualitymd models` surface (`list` with `--json`; `view <name>
  [--source]` as verbatim Markdown by default and structured JSON with `--json`,
  reusing the `lint --fix` node-rewrite to re-point the meta-model's apex
  `source`), homed in a new `internal/models` package for the bundled-model
  catalog. Settled the raw, non-OKF evaluation
  artifacts — JSON `assessments/`/`analysis/` source-of-record records and a
  `report.json` rendered over them beside the human `report.md`, altitude-first
  folder naming, and `improve`'s new-folder re-evaluation. Recorded the alternatives
  (plugin-marketplace-first/CLI-installer distribution, inline vs referencing
  `report.json`, meta-model embed home) and planned the In-Progress durable sync.
  Updated the change [index](archive/0010-implement-quality-skill/index.md), bundle
  [index](index.md), and the parent concept's status.

- **Refinement**: Tightened change
  [0010 — Implement the /quality skill](archive/0010-implement-quality-skill.md) after
  review. Reconciled the open items' conflicting lifecycle timing — they are now
  **surfaced in Draft**, with the **blocking** ones resolved before **Design** and
  the rest during **In-Progress**, all before **Done** (replacing contradictory
  "settle before Design" / "while Draft" / "before Done" wording across the spec and
  parent). Settled that **evaluation artifacts are raw runtime outputs, not OKF
  concepts**: JSON assessment/analysis records plus a `report.json`, with the worked
  example re-captured raw and its now-unused concept types retired from
  [`specs/schema.md`](../specs/schema.md). Brought the deterministic
  `qualitymd models` command into the change's **Covered** scope, since the model
  altitude depends on it and the skill drives the CLI rather than reimplementing it.
  Renamed the open-items anchor and aligned the parent's `Q1` references to the
  spec's item numbering. Updated the [log](log.md).

- **Creation**: Added change
  [0010 — Implement the /quality skill](archive/0010-implement-quality-skill.md)
  (`status: Draft`) with its
  [functional spec](archive/0010-implement-quality-skill/spec.md) to build the
  specified-but-unimplemented
  [`/quality` skill](../specs/skills/quality-skill/quality-skill.md). The spec
  **defers the behavioral contract** to the durable skill spec and states only the
  delta — package an invocable skill that conforms to it and drives the `qualitymd`
  CLI for every mechanical step — plus the open items and gaps a review of the
  skill spec surfaced for this change to settle: where the skill is packaged; where the
  **model** altitude draws its criteria (the built-in
  [meta-model](../internal/models/quality-meta-model.md) is
  neither referenced nor CLI-exposed); what `setup` does beyond `init`; how the
  default target file resolves (a CLI convention still "to be specified"); and
  whether `improve`'s post-apply re-evaluation writes a new evaluation folder.
  Records [`quality-skill.md`](../specs/skills/quality-skill/quality-skill.md) and
  user docs as affected durable artifacts, to sync once the questions resolve.
  Updated the bundle [index](index.md).

- **Completion**: Implemented and archived
  [0009 — Diagnose rating-scale soundness in the meta-model](archive/0009-rating-scale-diagnostic.md),
  adding the *rating scale and any overrides are well-formed and meaningful*
  requirement to the [meta-model](../internal/models/quality-meta-model.md)'s
  Functionality factor. The meta-model previously assessed the rating scale only
  structurally (lint's "well-shaped" check) and as one clause in a conformance
  list, despite the scale being what turns assessments into verdicts and despite
  per-requirement `ratings` overrides giving authors room to miscalibrate a
  threshold or quietly redefine a level. Synced the Functionality summary and the
  diagnostic coverage checklist, and confirmed `qualitymd lint` still reports the
  model valid. No durable specs/docs were affected — the requirement traces to the
  rating-scale semantics already in [`SPECIFICATION.md`](../SPECIFICATION.md),
  which is unchanged.

- **Creation**: Added change
  [0009 — Diagnose rating-scale soundness in the meta-model](archive/0009-rating-scale-diagnostic.md)
  (`status: Draft`) with its
  [functional spec](archive/0009-rating-scale-diagnostic/spec.md). Prompted by the
  meta-model's coverage asymmetry — six requirements for the Markdown body
  sections, but the rating scale assessed only for structural shape. The change
  adds a single Functionality requirement covering level meaning, band
  separability, floor placement against the model's needs and risks, and sound
  per-requirement criterion overrides, written to pass trivially for a model that
  inherits the suggested scale unchanged. Omits a design doc as a one-requirement
  content change; records no affected durable specs/docs.

- **Completion**: Implemented and archived
  [0008 — Describe targets with title and description](archive/0008-target-display-fields.md),
  adding `title` and `description` to Target, `description` to Model, and the
  matching durable [`SPECIFICATION.md`](../SPECIFICATION.md) prose and
  [`lint`](../specs/cli/lint.md) rule row. The structural schema now accepts
  target display fields, `misplaced-root-key` flags only nested `ratingScale`,
  and focused tests cover accepted nested target `title`/`description` plus the
  still-rejected nested `ratingScale`.

- **Implementation**: Advanced change
  [0008 — Describe targets with title and description](archive/0008-target-display-fields.md)
  from `Design` to `In-Progress` so the schema, linter, and durable
  specification updates can land.

- **Design**: Refined change
  [0008 — Describe targets with title and description](archive/0008-target-display-fields.md):
  made `Model.description` **Optional** (was `Recommended`), matching
  `Target.description`, so `description` reads uniformly across the tree. Updated
  the [functional spec](archive/0008-target-display-fields/spec.md) (Model schema now shows
  `description` as Optional) and the
  [design doc](archive/0008-target-display-fields/design.md): the `OptionalPresence`
  addition, the `# Optional` Model snippet, the composition alternative (now closer
  since `title`/`description` presence matches, still rejected on the `model-content`
  `RequiredAny` group and the mid-list `ratingScale` splice), and the trade-off note.

- **Design**: Advanced change
  [0008 — Describe targets with title and description](archive/0008-target-display-fields.md)
  from `Draft` to `Design` and added its
  [design doc](archive/0008-target-display-fields/design.md). Reading the code showed the
  change is almost entirely schema + prose: three property additions in
  [`internal/schema`](../internal/schema/schema.go) (`Target` gains `title`/
  `description`, `Model` gains `description`) drive everything, because
  `misplaced-root-key` already fires on "a key valid on `Model` but not on this
  target" and so **self-narrows** to `ratingScale` with no rule-logic change, and
  the data-driven spec↔schema consistency test passes once the declarations and
  `SPECIFICATION.md` snippets move in lockstep. The only test edit flips
  `rules_test.go`'s "nested target title" case from a finding to a valid model.
  Recorded the decisions: keep `Model`/`Target` as explicit property lists rather
  than composing `Model` from `Target` (their `title` presence diverges —
  `Recommended` on the root, `Optional` on a target — and the consistency test
  already guards drift); and **not** add a `missing-target-description` warning
  (`RecommendedPresence` is documentary, not auto-enforced), leaving it as a noted
  follow-up. Updated the change [index](archive/0008-target-display-fields/index.md) and
  bundle [index](index.md).

- **Creation**: Added change
  [0008 — Describe targets with title and description](archive/0008-target-display-fields.md)
  (`status: Draft`) with its
  [functional spec](archive/0008-target-display-fields/spec.md). A target's only
  human-facing label today is its map key; the change lets every target carry an
  optional `title` (display name) and a recommended `description` (what the target
  *is*), and adds `description` to the model root. It also reframes the root as a
  **Model** — the model-wide `ratingScale` plus all **Target** properties — so
  the difference from a nested target is a type distinction (`ratingScale` is the
  one Model-only property) rather than the awkward "a non-root target MUST NOT
  declare …" prohibition, matching how `internal/schema` already models the two
  as distinct nodes. Records the
  [`SPECIFICATION.md`](../SPECIFICATION.md) and
  [`lint` sub-spec](../specs/cli/lint.md) deltas plus the `internal/schema` and
  `internal/lint` conformance (the `misplaced-root-key` rule and its
  "nested target title" test case) as affected. Updated the bundle
  [index](index.md). A design doc follows at **Design** for how
  `misplaced-root-key` narrows to `ratingScale` alone.

- **Completion**: Implemented and archived
  [0007 — Delightful human CLI output](archive/0007-delightful-cli-output.md),
  giving the human surface a single brand palette shared with the Fang harness, a
  styled [`lint`](../specs/cli/lint.md) finding list (severity glyphs, color,
  clickable `file:line`, colored summary) and [`init`](../specs/cli/init.md)
  confirmation, runnable `--help` examples on all three commands,
  [`spec`](../specs/cli/spec.md) paging through `$PAGER`/`less`, and an
  informative `--version` recovered from the Go toolchain's embedded build info.
  All of it sits behind the TTY/`NO_COLOR` gate, so the agent-facing plain and
  `--json` paths are byte-for-byte unchanged. Added the **Human output styling**
  and **Binary version** conventions to the [CLI spec](../specs/cli.md), the
  paging clause to the [`spec` sub-spec](../specs/cli/spec.md), the shared
  `internal/cli/style.go`, and focused tests; the styling consolidates onto one
  `colorEnabled` gate, retiring `spec`'s bespoke `shouldRenderSpec`.

- **Completion**: Implemented and archived
  [0006 — Specify and implement the spec command](archive/0006-spec-command.md),
  replacing the placeholder [`spec` sub-spec](../specs/cli/spec.md), adding the
  [design doc](archive/0006-spec-command/design.md), registering
  `qualitymd spec`, embedding [`SPECIFICATION.md`](../SPECIFICATION.md) in the
  binary, rendering Markdown for TTY output with Glamour while preserving
  byte-for-byte Markdown for redirected/plain output, and updating the
  [`README.md`](../README.md) command status.

- **Completion**: Implemented and archived
  [0005 — Single source of truth for the structural schema](archive/0005-schema-source-of-truth.md),
  adding `internal/schema` as the typed structural schema declaration, deriving
  lint's unknown-key, shape, required-property, model-content, and rating-scale
  minimum checks from it, and adding tests that compare the declaration against
  [`SPECIFICATION.md`](../SPECIFICATION.md). Reconciled the public format
  snippet's `title` presence to `Recommended`, updated
  [`lint`](../specs/cli/lint.md) to record schema-derived structural validation,
  and added the [design doc](archive/0005-schema-source-of-truth/design.md).

- **Refinement**: Added a human-readable rendering to change
  [0006 — Specify and implement the spec command](0006-spec-command.md). `spec`
  now **SHOULD** render the specification formatted (via the stack's terminal
  renderer) when stdout is a terminal, while still writing **verbatim Markdown**
  when output must be plain (non-terminal or `NO_COLOR`) — so a redirect still
  reproduces the artifact byte-for-byte and the `--json` carve-out is unaffected.
  The rendered/verbatim split needs no flag: it rides the
  [baseline](../specs/cli.md#baseline)'s existing terminal-detection rule, exactly
  as color does. Updated the [functional spec](0006-spec-command/spec.md) and the
  change's scope.

- **Creation**: Added change
  [0006 — Specify and implement the spec command](0006-spec-command.md)
  (`status: Draft`) with its
  [functional spec](0006-spec-command/spec.md). Picks up the `spec` command that
  [0004](archive/0004-specify-agent-accessibility.md) deferred as "a separate
  change that inherits this baseline." The change settles the still-stub
  [`specs/cli/spec.md`](../specs/cli/spec.md) — whose open questions predate and
  now conflict with the agent-accessibility work (it floats a `--format json`
  form the settled [`--json` convention](../specs/cli.md#conventions) forbids) —
  and lands the command: emit the bundled format specification verbatim as
  Markdown to stdout, no arguments, no `--json` (the verbatim-artifact carve-out),
  full baseline conformance. Records [`specs/cli/spec.md`](../specs/cli/spec.md)
  and [`README.md`](../README.md) as affected; structured forms, sub-views, and
  `spec`-specific flags are deferred. A design doc follows at **Design** for how
  the root-level specification is embedded. Updated the bundle [index](index.md).

- **Refinement**: Recorded the schema-source direction for change
  [0005 — Single source of truth for the structural schema](archive/0005-schema-source-of-truth.md):
  a **typed Go declaration** the linter derives from directly (over an embedded
  data file or a `specs/` concept), with spec/linter consistency enforced by a
  test checking [`SPECIFICATION.md`](../SPECIFICATION.md) against it rather than by
  generating docs — lowest-machinery path that meets the spec's requirements, with
  data-file/generation revisited only if a second consumer appears. Left
  unknown-key typo suggestions as a deferred follow-up, untouched. Detail lands in
  the design doc at **Design**.

- **Creation**: Added change
  [0005 — Single source of truth for the structural schema](archive/0005-schema-source-of-truth.md)
  (`status: Draft`) with its
  [functional spec](archive/0005-schema-source-of-truth/spec.md). Prompted by reviewing
  design.md's linter, which derives its structural rules from one schema artifact:
  our structural schema is encoded twice — implicitly in `internal/lint/rules.go`
  and again in prose in [`SPECIFICATION.md`](../SPECIFICATION.md) and
  [`specs/`](../specs/index.md) — so the two can drift. The change requires a
  single authoritative definition of valid keys, required properties, and the
  rating-scale shape that the linter derives from, as a behavior-preserving
  refactor; records [`SPECIFICATION.md`](../SPECIFICATION.md) and
  [`specs/cli/lint.md`](../specs/cli/lint.md) as affected durable docs; and defers
  doc generation, runtime configuration, and unknown-key typo suggestions. Updated
  the bundle [index](index.md).

- **Completion**: Implemented and archived
  [0004 — Specify and enforce agent accessibility](archive/0004-specify-agent-accessibility.md),
  adding the durable [CLI spec](../specs/cli.md) agent-accessibility contract,
  categorized exit codes (`0`/`1`/`2`/`70`), the broadened `--json` convention,
  and the [`init --json`](../specs/cli/init.md) receipt contract. The
  implementation maps exit categories through Fang, suppresses duplicate stderr
  for already-reported lint findings, adds the neutral `internal/receipt` action
  type, and tests the exit categories plus receipt and overwrite-refusal shapes.
- **Design**: Advanced change
  [0004 — Specify and enforce agent accessibility](archive/0004-specify-agent-accessibility.md)
  from `Draft` to `Design` and added its
  [design doc](archive/0004-specify-agent-accessibility/design.md). Reading
  `fang@v1.0.0` confirmed `fang.Execute` returns the command error and never
  exits, so categorized exit codes ride Fang's intended model: a thin boundary
  mapping in `cli.Execute` (`errors.As` → category, default `ExitInternal`), a
  `CodedError` carrying the category plus a `Silent` flag, a `WithErrorHandler`
  that suppresses the already-reported `lint` found-problems error, and
  Cobra-native `FlagErrorFunc`/`Args` hooks to tag usage errors at their source —
  with only an unknown-subcommand string fallback left as an explicit open
  decision. Picked `0`/`1`/`2`/`70` for success / found-problems / usage /
  internal, broadened *internal error* to "could not complete the requested
  action" so guarded refusals (e.g. `init` overwrite without `--force`) have a
  home, and ruled `-` plus `--json` a usage error. Settled the one open
  design decision — keep a thin, owned prefix check so an unknown subcommand maps
  to `ExitUsage` (option a), failing safe to `ExitInternal` and pinned by a test.
  Reconciled the
  [functional spec](archive/0004-specify-agent-accessibility/spec.md)'s internal-error
  definition with that broadening, and updated the change
  [index](archive/0004-specify-agent-accessibility/index.md) and bundle
  [index](index.md).

- **Scope**: Broadened the `--json` convention within change
  [0004 — Specify and enforce agent accessibility](archive/0004-specify-agent-accessibility.md)
  after comparing against agent-first CLIs (e.g. Basecamp, where nearly every
  command accepts `--json`). The change now revises the
  [CLI spec](../specs/cli.md)'s `--json` convention from a narrow should-offer
  gate to a SHOULD-by-default: commands SHOULD offer `--json`, human rendering
  stays the default (no auto-JSON), a command MAY omit it only when its output is
  a verbatim artifact that is the payload (e.g. `spec`), and under `--json` a
  side-effecting command emits a **result receipt** rather than human prose. The
  conformance work gains `init --json` (a receipt of the written path / created
  flag / `nextActions`), and [`specs/cli/init.md`](../specs/cli/init.md) joins
  the affected durable docs — its "offers no `--json`" statement is replaced by
  the receipt contract. `spec` stays the deliberate carve-out. Updated the
  functional [spec](archive/0004-specify-agent-accessibility/spec.md), the change
  [index](archive/0004-specify-agent-accessibility/index.md), and the bundle
  [index](index.md).

- **Scope**: Expanded change
  [0004 — Specify and enforce agent accessibility](archive/0004-specify-agent-accessibility.md)
  from a spec-only change to spec **plus** conformance after auditing the shipped
  commands. `internal/cli` exits `1` on every error path, so a `lint` that *found
  problems* is indistinguishable from a usage or internal error — a baseline
  violation. The change now settles the **exit-code categories** (success,
  ran-but-found-problems, usage error, internal error) concretely, removes both
  *Agent-accessibility and CI requirements* and *Exit-code semantics* from the
  [CLI spec](../specs/cli.md)'s "To be specified" list, and brings `init` and
  `lint` into compliance with tests. Threading distinct codes through Fang is a
  real design question, so the change now carries a forthcoming
  [design doc](archive/0004-specify-agent-accessibility/design.md) (added at **Design**);
  the unimplemented [`spec`](../specs/cli/spec.md) command is scoped out as a
  separate change that inherits the baseline. Retitled accordingly and updated
  the change [index](archive/0004-specify-agent-accessibility/index.md) and bundle
  [index](index.md).

- **Creation**: Added change
  [0004 — Specify agent accessibility](archive/0004-specify-agent-accessibility.md)
  (`status: Draft`) with its
  [functional spec](archive/0004-specify-agent-accessibility/spec.md). The change settles
  the *Agent-accessibility and CI requirements* item on the
  [CLI spec](../specs/cli.md)'s "To be specified" list by adding an **Agent
  accessibility** section framed as two tiers: a baseline every in-scope command
  owes (non-interactivity, stdout-is-payload/stderr-is-everything-else,
  determinism, categorized exit codes, plain non-TTY output) and the opt-in
  capabilities (`--json`, `nextActions`, quiet/verbosity) gated by criteria and
  cross-referenced to the existing conventions. Records
  [`specs/cli.md`](../specs/cli.md) as the only affected durable doc; no command
  behavior changes. Omits a design doc as spec-only work. Updated the bundle
  [index](index.md).

- **Completion**: Implemented and archived
  [0003 — Implement the lint command](archive/0003-implement-lint-command.md),
  adding `qualitymd lint`, the shared lint rule catalog, JSON and human output,
  deterministic finding ordering, in-place `--fix` repair for fixable findings,
  parser/render/write support in `internal/spec`, and focused tests for the rule
  set, output shape, exit behavior, and repair behavior. Updated the README
  status and moved the change into [`archive/`](archive/).

- **Revision**: Hardened [0003's design doc](archive/0003-implement-lint-command/design.md)
  after review. Gave `internal/spec` a one-way dependency — it owns the document
  layer and no longer imports `internal/lint`, which now owns the rule catalog and
  the valid-model convenience (`lint.Load` replacing `spec.Load`) — removing a
  `spec`↔`lint` import cycle and a duplicate validator. Routed misplaced root-only
  keys (`title`/`ratingScale` on a non-root target) to `misplaced-root-key`
  instead of `invalid-frontmatter`; added the original Markdown body to the
  `Document` model so `Render` preserves it byte-for-byte; noted that `Load`'s
  acceptance tightens under the required-frontmatter parser; had `lint` reject a
  bare `-` this phase; clarified that post-repair `summary` counts (including
  `fixable`) reflect the re-lint; and reframed Resolved Questions as Open
  questions with the parent-CLI invocation as the one genuinely-open item.
  Recorded the provisional `lint [path]` shape as deliberately not durably specced
  in the [change](archive/0003-implement-lint-command.md).
- **Revision**: Worked down the open questions and risks in
  [0003's design doc](archive/0003-implement-lint-command/design.md): kept the shared
  document/model code in `internal/spec`, assigned rule-level repair operations
  to `internal/lint` and rendering/atomic replacement to `internal/spec`,
  resolved unknown frontmatter keys as `invalid-frontmatter` in this phase,
  confirmed `lint [path]` as the local invocation shape, and added mitigations
  for YAML round-tripping, deterministic ordering, atomic replacement, and
  symlink paths.
- **Revision**: Scoped `--fix` into change
  [0003 — Implement the lint command](archive/0003-implement-lint-command.md) after
  reviewing fixable-rule behavior. The durable lint spec, implementation spec,
  and design now require deterministic in-place repair of fixable findings,
  transactional per-file writes, post-repair linting, and JSON repair reporting,
  while keeping suppression, rule selection, and patch/full-file repair output
  modes deferred.
- **Design**: Advanced change
  [0003 — Implement the lint command](archive/0003-implement-lint-command.md) from
  `Draft` to `Design` and added its
  [design doc](archive/0003-implement-lint-command/design.md): `lint` parses once into a
  shared document/model graph with stable `modelPath` locations and optional
  source positions, runs narrow rule visitors from `internal/lint`, exposes the
  traversal primitives needed by current rules and future query commands, and
  adds a narrow repair writer for `lint --fix`. The design uses `lint [path]`,
  defaulting to `QUALITY.md`, as the minimum invocation shape while the parent
  CLI spec continues to own the broader file/stdin convention. Updated the
  change [index](archive/0003-implement-lint-command/index.md).

- **Creation**: Added change
  [0003 — Implement the lint command](archive/0003-implement-lint-command.md)
  (`status: Draft`) with a child
  [functional spec](archive/0003-implement-lint-command/spec.md). The change defers
  command-specific behavior to the completed durable
  [`lint` sub-spec](../specs/cli/lint.md), records README status updates as the
  durable docs work before Done, and calls out the remaining cross-cutting CLI
  invocation/file-argument convention as a dependency to settle before Design.
  Updated the bundle [index](index.md).

- **Archival**: Retired the placeholder [0001 — Example change](archive/0001-example-change.md)
  into [`archive/`](archive/) now that the bundle has real changes to follow,
  keeping it as the reference template the
  [propose-a-change guide](../docs/guides/work-with-change-cases.md) points to. Set its
  status to `Done`, fixed the relative links for the deeper path, and updated the
  bundle [index](index.md) and the [archive index](archive/index.md).

- **Completion**: Implemented and archived
  [0002 — Specify the init command](archive/0002-init-command.md), adding
  `qualitymd init`, replacing the durable [`init` sub-spec](../specs/cli/init.md),
  and updating the README status.

- **Refinement**: Tightened change [0002 — Specify the init command](archive/0002-init-command.md)
  after review: framed implementation as the change's own **In-Progress** phase
  rather than deferred work, specified that a successful `init` writes its
  confirmation to standard error (keeping stdout clean for `-` piping), recorded
  the non-atomic `--force` overwrite as a [design](archive/0002-init-command/design.md)
  risk, and trimmed the `--json` note in the
  [functional spec](archive/0002-init-command/spec.md) to a pointer to the
  [CLI spec](../specs/cli.md) convention.

- **Design**: Advanced change [0002 — Specify the init command](archive/0002-init-command.md)
  from `Draft` to `Design` and added its [design doc](archive/0002-init-command/design.md):
  the scaffold ships as a static `//go:embed` asset (comments and body prose can't
  round-trip through YAML struct marshalling), overwrite protection rides on an
  atomic `O_CREATE|O_EXCL` open, and a conformance test runs the embedded skeleton
  through `spec.Load`. Updated the change [index](archive/0002-init-command/index.md).

- **Creation**: Added change [0002 — Specify the init command](archive/0002-init-command.md)
  (`status: Draft`) with its [functional spec](archive/0002-init-command/spec.md), settling
  the "To be specified" list on the [`init` sub-spec](../specs/cli/init.md): the
  scaffold contents (seeded rating scale, a commented target → factor → requirement
  skeleton, recommended body sections as headed stubs), the output target and
  stdout (`-`) piping, and `--force` overwrite protection. Records
  [`specs/cli/init.md`](../specs/cli/init.md) and [`README.md`](../README.md) as
  affected. Updated the bundle [index](index.md).

- **Process**: Defined the relationship between `changes/` and the enduring
  [`specs/`](../specs/index.md) bundle (replacing the "independent for now"
  note) — a Change Case states a *delta* and is archived, while `specs/` and
  [`SPECIFICATION.md`](../SPECIFICATION.md) hold the *cumulative* source of
  truth. Added an **Affected specs & docs** section to the
  [Change Case concept](archive/0001-example-change.md) so each change records the durable
  specs and docs it creates or updates, brought into sync before `Done`.

## 2026-06-16

- **Initialization**: Created the `changes/` OKF bundle — a home for incremental
  work, independent of [`specs/`](../specs/index.md) for now. Added the bundle
  [index](index.md), [`schema.md`](schema.md) (`type: Schema`) registering the
  `Change Case`, `Functional Specification`, and `Design Doc` types, and an
  [`archive/`](archive/) folder for completed changes.
- **Creation**: Added a placeholder [Example change](archive/0001-example-change.md)
  (`status: Draft`) with child [spec](archive/0001-example-change/spec.md) and
  [design](archive/0001-example-change/design.md) concepts showing the intended shape.
