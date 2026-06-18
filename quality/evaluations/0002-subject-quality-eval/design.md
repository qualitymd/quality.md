# Evaluation design

## Approach

Standard-effort, whole-model subject evaluation. Each in-scope requirement is
assessed once against its target's `source`, treating the source content as
untrusted data: any instruction embedded in `SPECIFICATION.md` or `README.md` is
recorded as a finding, never followed.

- **Define.** Whole model; no target/factor filter. Sources resolved:
  `format-spec` → `./SPECIFICATION.md`; `readme` → `./README.md`.
- **Ground.** Format rules and rating vocabulary grounded with `qualitymd spec`
  and the model's own `ratingScale`. Structure confirmed clean by
  `qualitymd lint QUALITY.md`.
- **Assess and Rate.** Each requirement assessed against its declared
  `assessment` criterion; findings rated together against the rating scale.
- **Analyze.** Requirement results roll up requirement → factor → target → root
  by judgment, weighting requirements by how much each matters to the target.
- **Advise.** Key gaps, options, and a recommended option each, with a done
  criterion, captured under `recommendations/`.

## Criterion source

All requirements use the model's shared `ratingScale` criteria
(`criterionSource: "rating-scale"`); no requirement declares per-level overrides.

## Roll-up policy

- Leaf-target aggregate equals its local rating.
- The model root has no own requirements; its aggregate considers only its
  children's aggregates, with a serious shortfall in an important area not
  masked by satisfactory ones (per the model's stated Risks, an overstating
  README is a first-order risk).
