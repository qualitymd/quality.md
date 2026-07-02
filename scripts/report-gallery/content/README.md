# Software service report gallery

This gallery is a generated, illustrative QUALITY.md example for a fictional
software service named LedgerLite. It exists to make evaluation report design
and QUALITY.md authoring practice easy to inspect without cutting a release or
running a fresh /quality evaluate.

The [`QUALITY.md`](QUALITY.md) is written to demonstrate the /quality skill's
authoring guide family: body sections with unknowns, open questions, and review
provenance; a composite root with a model-wide agent harnessability factor; a
normative service-contract area; an implementation `codebase` area; a veto
requirement; a measured rating override; and assessments that demonstrate the
guide-and-sensor posture by reusing named computational sensors, reading sensor
results against guides or runbooks, and leaving inferential judgment visible
where no sensor exists yet.

Open the generated run report:

- [Evaluation report](.quality/evaluations/0001-full-eval/report.md)
- [Findings](.quality/evaluations/0001-full-eval/findings.md)
- [Recommendations](.quality/evaluations/0001-full-eval/recommendations.md)

The [quality changelog](.quality/changelog/) shows how the model matured
through earlier quality-loop cycles before this evaluation: areas added,
criteria recalibrated, and drift corrected. Earlier evaluation runs from those
cycles are not retained in the gallery; only the latest run is checked in.

The generated evaluation uses synthetic routine outputs. Findings, ratings,
roll-ups, recommendations, and synthetic-source:\* evidence references are
fictional and demonstrate report structure only; they are not an assessment of
a real system. The concrete source system and source-code tree are
intentionally omitted.

Regenerate this gallery from the repository root with
`mise run report-gallery`.

Do not edit files under .quality/ by hand; the generator owns the evaluation
run and the changelog entries.
