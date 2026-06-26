# /quality Runtime Skill Update Log

## 2026-06-26

- **Revision**: Retitled the run frame and unified workflow vocabulary for 0110 -
  Run frame title and workflow vocabulary.
  The run-frame template header is now `**Quality · <workflow>**` instead of
  `**/quality run**`, the `Mode:` field is dropped (the workflow name moves into
  the header), and Arguments/Workflow Dispatch wording now says "workflow" rather
  than "mode".

- **Revision**: Updated the root runtime interaction contract for 0106 - Binary
  confirmation UX.
  Runtime guidance now distinguishes non-binary closed choices, which keep `1` as
  the shortest accept path, from true binary mutation confirmations, which show
  `y`/`n`.

- **Revision**: Updated the root runtime interaction contract for 0101 - Quality
  skill UX action clarity.
  Runtime guidance now requires explicit shortest-answer paths for user
  interactions, code spans for concrete operational examples, and numbered
  ambiguity choices for scoped evaluation prompts.

## 2026-06-24

- **Restructure**: Started the runtime skill content as an OKF-shaped bundle with
  root `index.md`, `schema.md`, and `log.md`; added guide indexes/logs; and split
  authoring guidance into routed sub-guides.
