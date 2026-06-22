# /quality Skill Guides Update Log

## 2026-06-22

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
