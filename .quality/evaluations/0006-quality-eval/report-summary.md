# Quality evaluation summary

| Field | Value |
| --- | --- |
| Root area | QUALITY.md Project |
| Run | `0006-quality-eval` |
| Scope | Standard. |
| Rigor | Standard |
| Evaluation verdict | 🟡 Minimum |
| Full report | [report.md](report.md) |
| Machine report | [report.json](report.json) |

## Verdict

The root local rating is target, but aggregate quality is held to minimum by minimum-rated child areas for docs, distribution, and evaluation history.

## Area breakdown

| Area | Path | Area Rating | Area + Sub-Areas Rating | Factors |
| --- | --- | --- | --- | --- |
| QUALITY.md Project | `/` | 🔵 Target | 🟡 Minimum | Agent Harnessability: 🔵 Target; Agent Harnessability / Agent Accessibility: 🔵 Target; Agent Harnessability / Task Specifiability: 🔵 Target; Agent Harnessability / Agent Operability: 🔵 Target; Agent Harnessability / Continuity: 🔵 Target; Agent Harnessability / Self-Verifiability: 🔵 Target; Agent Harnessability / Enforcement of Standards: 🔵 Target; Agent Harnessability / Containment of Action: 🔵 Target |
| Agent Harness | `agent-harness` | 🔵 Target | 🔵 Target | Completeness: 🔵 Target; Coherence: 🔵 Target; Currentness: 🔵 Target; Assessability: 🔵 Target |
| qualitymd CLI | `cli` | 🔵 Target | 🔵 Target | Specification Conformance: 🔵 Target; Automation Compatibility: 🔵 Target; Usability: 🔵 Target; Maintainability: 🔵 Target |
| Installation and Distribution | `distribution` | 🟡 Minimum | 🟡 Minimum | Installability: 🟡 Minimum; Release Readiness: 🔵 Target |
| Documentation and Examples | `docs` | 🟡 Minimum | 🟡 Minimum | Approachability: 🔵 Target; Domain Range: 🔵 Target; Currentness: 🟡 Minimum |
| Evaluation History | `evaluation-history` | 🟡 Minimum | 🟡 Minimum | Reportability: 🟡 Minimum; Traceability: 🔵 Target |
| Format Specification | `format-spec` | 🔵 Target | 🔵 Target | Clarity: 🔵 Target; Consistency: 🔵 Target; Completeness: 🔵 Target; Verifiability: 🔵 Target; Domain Agnosticism: 🔵 Target |
| QUALITY.md Project QUALITY.md | `quality-md` | 🔵 Target | 🔵 Target | Context Grounding: 🔵 Target; Model Structure: 🔵 Target; Evaluability: 🔵 Target; Lifecycle Maintenance: 🔵 Target |
| /quality Skill | `quality-skill` | 🔵 Target | 🔵 Target | Judgment Grounding: 🔵 Target; Mutation Safety: 🔵 Target; Agent-Mediated UX: 🔵 Target; Workflow Completeness: 🔵 Target |
| Tooling Specs Bundle | `specs-bundle` | 🔵 Target | 🔵 Target | Traceability: 🔵 Target; Consistency: 🔵 Target |

## Selected findings

1. **Low**  
   AGENTS.md names README.md and CONTRIBUTING.md as required context and links task-specific guides.
   `AGENTS.md:5`
   Assessment: `assessments/001-root-agents-can-reach-the-minimum-project-context-and-deeper-routed-guidance.json`
2. **Low**  
   Routine changes are scoped directly, while guide routing and change-case rules define when durable design history is needed.
   `AGENTS.md:33`
   Assessment: `assessments/002-root-project-guidance-makes-task-boundaries-and-done-criteria-explicit.json`
3. **Low**  
   The check task aggregates format, tidy, vet, lint, test, Markdown, and npm package checks.
   `mise.toml:38`
   Assessment: `assessments/003-root-agents-can-discover-the-tools-and-commands-needed-to-work-the-project.json`
4. **Low**  
   The skill requires evaluate to write numbered records through the CLI and feedback logs under .quality/logs/.
   `skills/quality/SKILL.md:47`
   Assessment: `assessments/004-root-quality-workflows-preserve-useful-state-in-durable-local-artifacts.json`
5. **Low**  
   The full local check passed across gofmt, tidy, dprint, vet, tests, npm package verification, and golangci-lint.
   `command:mise run check`
   Assessment: `assessments/005-root-verification-signals-are-discoverable-and-remediation-bearing.json`

## Recommended actions

Primary next action: use `001-refresh-pinned-installer-examples`.

| Recommendation ID | Priority | Recommendation | Done criterion |
| --- | --- | --- | --- |
| `001-refresh-pinned-installer-examples` | 1 | [Refresh Pinned Installer Examples](recommendations/001-refresh-pinned-installer-examples.md) | The install guide no longer contains stale v0.5.1 pins, pinned examples cannot age silently, and the docs and distribution assessments reach target on re-evaluation. |
| `002-reconcile-stale-evaluation-history` | 2 | [Reconcile Stale Evaluation History](recommendations/002-reconcile-stale-evaluation-history.md) | qualitymd status QUALITY.md --json no longer reports a stale non-reportable prior run, or the repository has an explicit documented archive state that keeps stale runs out of current-readiness gaps. |

## Scope and limitations

Scope: **Standard.**

In scope: QUALITY.md Project; Agent Harness; qualitymd CLI; Installation and Distribution; Documentation and Examples; Evaluation History; Format Specification; QUALITY.md Project QUALITY.md; /quality Skill; Tooling Specs Bundle

Limitations: none recorded.
