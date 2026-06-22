# Evaluation plan

Whole-model subject evaluation of the QUALITY.md project.

1. Resolve config (none present → default `quality/evaluations/`) and scope
   (whole model, standard effort).
2. Lint `QUALITY.md` (passed) and ground rules with `qualitymd spec`.
3. Assess 13 requirements across two targets:
   - `format-spec` (9): completeness + clarity (3), consistency (1),
     verifiability (2), extensibility (1), usability (1).
   - `readme` (4): approachability.
4. Write one assessment record per requirement under `assessments/`.
5. Roll up per target under `analysis/`, then the root.
6. Produce `report.md` + `report.json` and per-gap `recommendations/`.
