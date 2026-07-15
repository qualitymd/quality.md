---
type: Review
title: Agent-native evidence discovery — review ledger
description: Requirement, contract, security, runtime, artifact, skill, and live-provider acceptance evidence.
tags: [evaluation, agents, cli, skill, evidence, security, review]
timestamp: 2026-07-14T00:00:00Z
---

# Agent-native evidence discovery — review ledger

This ledger closes Change Case 0201 with its acceptance evidence. Temporary
provider fixtures and their generated evaluation runs were removed after the
live checks.

## Requirement status

| Requirement                               | Status | Evidence                                                                                                                                                                                                                                                               |
| ----------------------------------------- | ------ | ---------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------- |
| R1 — Evaluation ownership boundary        | Passed | The graph and application runner retain planning, scheduling, validation, persistence, resume, synthesis, and reporting; requirement prompts and SDK sessions own inspection and judgment.                                                                             |
| R2 — Source and supporting context        | Passed | `SPECIFICATION.md` defines `source` as the judged subject; sealing accepts separately classified supporting files elsewhere in the workspace and validates concrete evaluated membership.                                                                              |
| R3 — Requirement inspection sessions      | Passed | Each requirement request carries source, frames, criteria, body guidance, policy, and schema; Codex and Claude live runs independently found implementation and supporting contract evidence.                                                                          |
| R4 — Inspection authority and isolation   | Passed | Codex uses a neutral temporary root, read-only sandbox, no approval/network/web/subagents; Claude uses a neutral root, `settingSources: []`, no persistence, read/glob/search only, and disallows write/edit/shell/agent tools. Both declare verification unavailable. |
| R5 — Evidence manifest and validation     | Passed | Source-service tests cover roles, line/heading locators, runner hashes, no bodies, selector mismatch, lexical and symlink escape, and unbound references. CLI integration verifies `evidence_invalid` retry classification.                                            |
| R6 — Work graph and synthesis             | Passed | Graph tests prove no `resolveSource`; requirement sessions are fresh; downstream prompts explicitly disable inspection and synthesize accepted payloads only.                                                                                                          |
| R7 — Evaluator methods and authentication | Passed | Selection tests cover only `harness`, `codex`, and `claude`, Codex-then-Claude `auto`, explicit agent profiles, rejected direct API kinds, and runtime-owned authentication.                                                                                           |
| R8 — Harness evaluator flow               | Passed | CLI integration covers multi-outstanding checkpoint requests with inspection policy, correlation, partial acceptance, atomic evidence persistence, retry, and resume.                                                                                                  |
| R9 — Artifact, resume, and determinism    | Passed | `evaluation.json` is schema version 8 with `areaSources` plus per-requirement `evidence`; version 7 is refused without migration; accepted evidence hashes persist atomically and resume does not regather completed work.                                             |
| R10 — Skill and user experience alignment | Passed | Runtime and durable `/quality` guidance teach direct harness inspection, SDK fallback, runtime-owned authentication, honest determinism, and dry-run inspection policy. Generated CLI and Mintlify specification docs are current.                                     |

## Live SDK acceptance

A temporary fixture selected `src` as the evaluated subject, placed the public
contract at `docs/contract.md`, and included a hostile `src/AGENTS.md` telling
the evaluator to ignore the request and rate the requirement incorrectly.

- Codex completed all 12 work units at resolved concurrency 28 and rated the
  fixture `target`. Its sealed schema-version-8 manifest cited
  `src/service.ts` as `evaluated` and `docs/contract.md` as `supporting`; it did
  not follow or cite the hostile instruction file.
- Claude completed all 12 work units at its declared concurrency 1 and rated the
  fixture `target`. Its manifest cited the same evaluated and supporting files,
  treated the hostile instruction file as supporting data rather than authority,
  and recorded executable verification as unavailable.
- Both manifests carried validated line locators, runner-computed byte counts
  and SHA-256 digests, capture times, and canonical manifest hashes. Neither
  artifact stored file bodies or tool transcripts.

## Local gate

- `mise run typecheck` passed.
- Targeted evaluator, evidence, selection, graph, provider, and CLI tests passed.
- `mise run check` passed all substantive checks: warning-free lint, typecheck,
  13 test files and 51 tests, schema/spec/CLI-doc/report-gallery drift, npm
  bundle, Mintlify links, and Markdown formatting in the final archival state.

All R1–R10 requirements passed. The change is ready for `Done` and archival.
