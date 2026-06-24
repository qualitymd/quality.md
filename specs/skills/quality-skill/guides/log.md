# /quality Skill Guides Update Log

## 2026-06-24

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
