# /quality Skill Guides Update Log

## 2026-06-26

- **Revision**: Updated the
  [recommendation follow-up guide contract](recommendation-follow-up-md.md) for
  [0128 - Agent-mediated skill alignment](../../../../changes/archive/0128-agent-mediated-skill-alignment.md).
  Recommendation follow-up now requires a first-output follow-up frame before
  recommendation inspection, outcome selection, local apply, issue creation, or
  quality-log writes.

- **Revision**: Updated the
  [recommendation follow-up guide contract](recommendation-follow-up-md.md) for
  [0123 - Render interactions through native affordances](../../../../changes/archive/0123-native-interaction-affordances.md).
  The apply-vs-hand-off outcome is now a single-select closed-choice intent
  rendered per the shared progressive-enhancement contract, with the numbered
  options as its text fallback.

- **Revision**: Updated the
  [Agent Harnessability guide contract](authoring/agent-harnessability.md) and
  [Top 10 checks contract](top-10-quality-md-checks-md.md) for early-alpha
  compatibility cleanup.
  The durable guide specs now require old `harnessability` factors to be reported
  as stale legacy naming, not accepted as current Agent Harnessability coverage.

- **Revision**: Updated the
  [recommendation follow-up contract](recommendation-follow-up-md.md) for
  [0121 - Scannable interaction hierarchy](../../../../changes/0121-scannable-interaction-hierarchy.md).
  Apply and issue-creation decisions must now lead with the question, render
  choices as a visually separated block with the alternative folded into the stop
  choice, and cap and demote supporting fields; result closeouts must lead with a
  primary outcome line and avoid stacking equally-weighted bold labels.

- **Revision**: Updated the durable authoring sub-guide specs for
  [0107 - Durable spec alignment](../../../../changes/archive/0107-durable-spec-alignment.md).
  The eight routed authoring guide contracts now declare their BCP 14 keyword
  convention and use explicit companion notes for the runtime guide each spec
  governs.

- **Revision**: Updated the
  [recommendation follow-up contract](recommendation-follow-up-md.md) for
  [0101 - Quality skill UX action clarity](../../../../changes/archive/0101-quality-skill-ux-action-clarity.md).
  Recommendation follow-up now requires numbered apply-vs-handoff outcome
  selection when the user has not already chosen, decision briefs before external
  issue creation, and result closeouts with `Next` plus boundary-sensitive
  `Not done`.

## 2026-06-25

- **Revision**: Updated durable [authoring](authoring.md), authoring
  [index](authoring/index.md), and
  [Requirements](authoring/requirements.md) guide specs for
  [0093 - Named Requirement identity](../../../../changes/archive/0093-requirement-identity.md).
  The durable guide contract now requires stable Requirement names,
  natural-language Requirement titles, and named Requirement authoring examples.

- **Revision**: Updated durable [Getting started](getting-started-md.md) and
  [Top 10 checks](top-10-quality-md-checks-md.md) guide specs for
  [0092 - Setup workflow scope trim](../../../../changes/archive/0092-setup-workflow-scope-trim.md).
  Getting-started now treats important model gaps as the first-run iteration
  starting point, and Top 10 now separates lifecycle state from model-usefulness
  findings without maturity or evaluation-readiness classifications.

## 2026-06-24

- **Revision**: Updated the
  [agent-harness Area guide contract](authoring/agent-harness.md),
  [Agent Harnessability guide contract](authoring/agent-harnessability.md),
  [model-structure guide contract](authoring/model-structure.md), and
  [Top 10 checks contract](top-10-quality-md-checks-md.md) for
  [0091 - Agent-harness holistic definition](../../../../changes/archive/0091-agent-harness-holistic-definition.md).
  The durable spec mirror now requires the holistic harness definition,
  checked-in governing-artifacts Area projection, mixed-artifact scoping rule,
  expanded requirement shapes across feedforward, feedback, and owned controls,
  and findings for instructions-only or unmodeled runtime-harness gaps.

- **Restructure**: Replaced the monolithic `authoring-md.md` guide contract with
  the [authoring guide family](authoring.md) and mirrored sub-specs under
  [authoring/](authoring/index.md) for
  [0090 - Skill-content OKF authoring split](../../../../changes/archive/0090-skill-content-okf-authoring-split.md).
  The durable spec tree now mirrors the runtime authoring sub-guide tree, with
  separate review surfaces for body authoring, model structure, factors,
  requirements, rating scale, Agent Harnessability, agent-harness Area modeling,
  and quality-log judgment.

- **Revision**: Updated the [authoring guide contract](authoring-md.md) and the
  [Top 10 checks contract](top-10-quality-md-checks-md.md) for
  [0089 - Agent-harness modeling guidance](../../../../changes/archive/0089-agent-harness-modeling-guidance.md).
  Agent Harnessability now has seven current sub-factors with `continuity` added
  for long-running work state preservation and resumption; self-verifiability now
  carries good-sensor and trace-evidence properties; the agent-harness area gains a
  domain-agnostic steering-materials factor and requirement template; and the Top
  10 checks now flag thinly factored or software-leaking harness areas while still
  recognizing legacy six-sub-factor harnessability as prior semantic coverage.

- **Revision**: Updated the [Top 10 checks contract](top-10-quality-md-checks-md.md)
  for [0088 - Domain-agnostic corpus alignment](../../../../changes/archive/0088-domain-agnostic-corpus-alignment.md).
  The area-and-factor-shape check now pairs its software constituent illustration
  with a data-product bracket — schema, provenance, freshness, and lineage metadata
  — so the readiness check does not imply software is the default modeled domain.

- **Revision**: Updated the [authoring guide contract](authoring-md.md) and the
  [Top 10 checks contract](top-10-quality-md-checks-md.md) for
  [0087 - Encode projection boundaries in the model](../../../../changes/0087-encode-projection-boundaries.md).
  The guide must now require the author to encode a concern's projection boundary in
  the emitted model when more than one projection is modeled — a YAML comment on each
  projection's node, plus a disambiguating `description` clause when both projections
  are rated nodes that surface in a report — with the Agent Harnessability factor vs.
  the agent-harness area as the canonical instance, and the Top 10 checks gain a
  matching missing-boundary-note finding.
- **Revision**: Corrected the [authoring guide contract](authoring-md.md) for
  [0086 - Umbrella factor roll-up framing](../../../../changes/0086-umbrella-factor-rollup-framing.md).
  The Agent Harnessability umbrella factor is no longer described as one that "does
  not roll up directly"; the contract now requires the guide to present it as
  carrying no requirements of its own and being rated by rolling up its
  sub-factors, matching the grouping-area roll-up semantics in `SPECIFICATION.md`.
- **Revision**: Updated the [authoring guide contract](authoring-md.md) and the
  [Top 10 checks contract](top-10-quality-md-checks-md.md) for
  [0085 - Agent Harnessability naming](../../../../changes/archive/0085-agent-harnessability-naming.md).
  The model-wide factor is now named Agent Harnessability with the recommended
  `agent-harnessability` key, an accountability-preserving definition that keeps
  human direction, review, and accountability explicit, and legacy handling that
  recognizes an existing `harnessability` factor with the expected six sub-factors
  as semantic coverage while recommending the new name during model-authoring work.
- **Revision**: Updated the
  [recommendation follow-up guide contract](recommendation-follow-up-md.md) for
  [0084 - Agent-mediated UX conformance](../../../../changes/archive/0084-agent-mediated-ux-conformance.md).
  Apply and issue-creation decisions now follow the shared agent-mediated UX
  contract: visually emphasized primary calls to action, scannable decision
  labels, explicit mutation boundaries, and status-first result reporting.
- **Revision**: Aligned the illustrative domain example lists in the bundled
  [authoring guide](../../../../skills/quality/guides/authoring.md) and
  [setup workflow](../../../../skills/quality/workflows/setup.md) with the canonical
  secondary-domain set for
  [0083 - Quality-domain agnosticism guide](../../../../changes/0083-quality-domain-agnosticism.md)
  — adding a research/analytical report to the named example domains. Illustrative
  example lists only; no guide-contract requirement changed.
- **Revision**: Updated the [authoring guide contract](authoring-md.md) for
  [0082 - Normalize QUALITY.md self-check roll-up](../../../../changes/archive/0082-normalize-quality-md-rollup.md).
  The QUALITY.md self-check must now be taught as an ordinary modeled Area using
  the `quality-md` key, `<Root Title> QUALITY.md` title shape, explicit
  path-based `source`, model-artifact Factors, and an authoring-guide Requirement.
  When in scope, it is assessed, analyzed, reported, and rolled up like any other
  Area; the guide must not keep it out of aggregate rating or on a separate
  evaluation axis solely because its source is `QUALITY.md`.
- **Revision**: Updated the [authoring guide contract](authoring-md.md) and the
  [Top 10 checks contract](top-10-quality-md-checks-md.md) for
  [0081 - Harnessability factor](../../../../changes/archive/0081-harnessability-factor.md).
  The guide now must teach harnessability as the model-wide factor projection of
  the agent-collaboration concern for agent-collaborated composite roots, distinct
  from the agent harness constituent and agent audience. It is a non-directly-rated
  umbrella decomposed into six sub-factors (agent-accessibility,
  task-specifiability, agent-operability, self-verifiability,
  enforcement-of-standards, containment-of-action), proposed by default and never
  omitted for harness thinness. The Top 10 area-and-factor-shape check now flags
  missing harnessability coverage unless the factor is not germane.
- **Revision**: Updated the [authoring guide contract](authoring-md.md) and the
  [Top 10 checks contract](top-10-quality-md-checks-md.md) for
  [0080 - Model constituents by default](../../../../changes/0080-model-constituents-by-default.md).
  The constituent-coverage contract is now **model-by-default**: enumerate the
  implied kinds and model each as its own area unless one of two disqualifiers
  holds (no distinct concerns → fold; not germane / outside the boundary → out of
  scope). A germane concern is never omitted in prose — an absent or thin artifact
  is surfaced as a ratable gap (a minimal area with a missing-anchor finding, or a
  requirement on an existing area), with a routing criterion between the two.
  Deferral is demoted to a blocker-only exception, and the Top 10 check flags a
  non-disqualified constituent left unmodeled or merely deferred. Replaces the
  earn-it inclusion test; the 0076/0077 generators and 0079 vocabulary discipline
  are unchanged.
- **Revision**: Updated the [authoring guide contract](authoring-md.md) for
  [0079 - Stewardship vocabulary discipline](../../../../changes/0079-stewardship-vocabulary-discipline.md).
  The three-projections requirement now also requires the guide to keep the
  motivation-layer stewardship/care vocabulary from modifying or replacing a
  taxonomy noun — a concern is the source a factor projects from, not a kind of
  factor — and to name the root's recurring factors as model-wide (cross-cutting)
  factors rather than "stewardship factors" or "stewardship lenses." The singular
  "a factor is a quality lens" gloss is preserved; the 0076/0077 grounding is
  unchanged.
- **Revision**: Updated the [authoring guide contract](authoring-md.md) for
  [0077 - Care-grounded stewardship concerns](../../../../changes/0077-stewardship-care-grounding.md).
  The guide's stewardship-concern generator must now read as care — an activity of
  tending whose artifact is its trace — so the claim that earns a constituent comes
  from a Need or Risk rather than the generator list, a present artifact is evidence
  (an area) distinct from whether the tending is done well (a factor), and the
  protective pair is cross-cutting stewardship under vulnerability. Framing only;
  the nine concerns, the two axes, and the earn-it test are unchanged.
- **Revision**: Updated the [authoring guide contract](authoring-md.md) and the
  [top-10-quality-md-checks](top-10-quality-md-checks-md.md) guide contract for
  [0076 - Domain constituent kinds and stewardship concerns](../../../../changes/0076-domain-constituent-kinds.md).
  The authoring guide must now teach enumerating a composite root's domain
  constituents by **constituent kind** using two generators — a
  stewardship-concern axis (a lifecycle band plus the protective pair secure and
  safeguard) and an audience×purpose axis (Diátaxis on the *enable* concern) —
  the three-projections rule (a concern projects as factor, constituent, and
  audience without double-counting), and the rule that germane high-leverage kinds
  are carried as areas even when thin or missing, earned not as a roster. The Top
  10 area-and-factor-shape check must flag a domain whose implied constituent
  kinds are neither modeled nor accounted for, earned-not-roster.

- **Revision**: Updated the [authoring guide contract](authoring-md.md) for
  [0075 - Rating title emoji defaults](../../../../changes/0075-rating-title-emoji-defaults.md).
  The guide must now present `🟢 Outstanding`, `🔵 Target`, `🟡 Minimum`, and
  `🔴 Unacceptable` as the recommended display titles for the standard four-level
  Rating Scale while keeping stable Rating Level IDs plain and treating emoji as
  a human scanning aid rather than semantics.

- **Revision**: Updated the [authoring guide contract](authoring-md.md) and the
  [top-10-quality-md-checks](top-10-quality-md-checks-md.md) guide contract for
  [0074 - Composite root areas and use-context constituents](../../../../changes/0074-composite-root-areas.md).
  The authoring guide must now teach three recursive, composable decomposition
  shapes (primary-subject, collection, composite), the two recurring use-context
  constituents (agent harness and QUALITY.md self-check, the latter kept out of
  the entity roll-up), and the factor-coverage aim scoped per primary-subject
  node. The Top 10 area-and-factor-shape check must flag a composite entity
  flattened into one root and a missing expected use-context constituent, treated
  as earned defaults rather than a required roster.

## 2026-06-23

- **Revision**: Updated the [authoring guide contract](authoring-md.md) so the
  runtime guide must include a worked Markdown body section example and cite the
  agent surface plus model in each `agent-reviewed` state line.

- **Revision**: Updated the
  [top-10-quality-md-checks](top-10-quality-md-checks-md.md) guide contract for
  [0065 - Setup discovery and close refinements](../../../../changes/0065-setup-discovery-and-close-refinements.md).
  The checklist now keeps two axes distinct — lifecycle state (owned by
  `qualitymd status`) and model maturity (`starter`, `immature`,
  `evaluation-ready`) — instead of one blended classification, and must include a
  condensed close checklist the setup workflow uses to classify maturity without
  reading every check.

- **Revision**: Updated the getting-started guide contract for
  [0064 - Structured setup workflow](../../../../changes/0064-structured-setup-workflow.md).
  Getting-started now treats setup assumptions from the structured workflow —
  including root Area, domain, review posture, and handoff posture — as the
  starting point for first-run iteration.

- **Revision**: Updated the getting-started guide contract for
  [0063 - Contextual setup flow](../../../../changes/0063-contextual-setup-flow.md).
  Getting-started is now post-setup iteration guidance for starter or immature
  models, including setup assumptions, stakeholder needs, missing context,
  recurring review or handoff posture, and next-step choices.

- **Revision**: Expanded the top-10-checks guide contract for
  [0063 - Contextual setup flow](../../../../changes/0063-contextual-setup-flow.md).
  The checklist now treats durable setup assumptions such as project posture,
  stakeholder needs, agent/collaboration fit, missing context, and quality-loop
  expectations as ongoing model-readiness concerns.

- **Revision**: Updated the getting-started and top-10-checks guide contracts for
  [0062 - Remove wizard mode](../../../../changes/0062-remove-wizard-mode.md).
  Getting-started now recommends public follow-on workflows directly, and the
  top-10 checklist supports read-only orientation and model review rather than
  wizard-specific routing.

## 2026-06-22

- **Revision**: Updated the [authoring guide contract](authoring-md.md) for
  [0058 - Model reference identifiers](../../../../changes/archive/0058-model-reference-identifiers.md)
  so the guide must teach strict Area names, Factor names, Rating Level IDs,
  structured IDs, and canonical model references while keeping Requirement
  statements natural-language keys.

- **Creation**: Added the
  [`recommendation-follow-up.md`](recommendation-follow-up-md.md) guide contract
  for the runtime recommendation follow-up guide.

- **Revision**: Renamed the guide contract specs to follow the 1:1
  artifact-spec filename convention: `authoring-md.md`,
  `getting-started-md.md`, and `top-10-quality-md-checks-md.md`. Runtime guide
  artifact filenames remain `authoring.md`, `getting-started.md`, and
  `top-10-quality-md-checks.md`.

## 2026-06-21

- **Revision**: Clarified that the authoring, getting-started, and top-10-checks
  guide contracts treat the Markdown body as evaluable judgment context. Body
  sections should be concise, self-explanatory, and grounded in
  agent-accessible support, with material inaccessible support captured in the
  relevant section's unknowns or open questions.

## 2026-06-19

- **Revision**: Clarified guide boundaries: authoring is the best-practices
  prerequisite and getting-started is the first-run process/outcomes guide.

- **Creation**: Added the Top 10 QUALITY.md checks guide contract for quick
  read-only model/lifecycle inspection findings used by wizard and related
  modes.

- **Revision**: Clarified that getting-started Known gaps includes known
  unknowns: missing context, unresolved questions, and evidence gaps.

- **Revision**: Added desired outcomes for each getting-started Markdown body
  section so the body can better support initial model authoring.

- **Revision**: Updated the getting-started guide contract so the rating scale
  follows the Markdown body before the rest of the model tree is expanded.

- **Revision**: Tightened the getting-started guide contract so first-run
  authoring fills the Markdown body before expanding the quality model tree.

- **Creation**: Added the guides subfolder, moved the authoring guide contract
  into it, and added the getting-started guide contract for first-run model
  population after `qualitymd init`.
