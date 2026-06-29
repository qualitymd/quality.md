# QUALITY.md domain terms

Use these terms distinctly:

| Term           | Meaning                                                                 | Prefer for                                                           |
| -------------- | ----------------------------------------------------------------------- | -------------------------------------------------------------------- |
| **Project**    | The quality subject and use context modeled by QUALITY.md.              | Value proposition, setup, evaluation intent, model meaning.          |
| **Workspace**  | The filesystem/tooling context resolved from one selected `QUALITY.md`. | `qualitymd status`, config, `.quality/`, evaluations, logs, updates. |
| **Repository** | A version-control and filesystem containment boundary.                  | Path safety, local discovery, Git-specific statements.               |

A project is not necessarily a repository. A repository may contain one
workspace, many workspaces, or only part of a project. A workspace root is where
tooling resolves paths and artifacts; it is not automatically the whole modeled
project.
