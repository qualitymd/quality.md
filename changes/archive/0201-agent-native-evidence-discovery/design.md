---
type: Design Doc
title: Agent-native evidence discovery — design
description: A simplified requirement work graph, isolated SDK inspection sessions, runner-sealed evidence manifests, and a clean removal of direct API evaluators.
tags: [evaluation, agents, cli, skill, evidence, security]
timestamp: 2026-07-14T00:00:00Z
---

# Agent-native evidence discovery — design

## Context

This design is implemented.

The current runner detects an area's source-selector kind, resolves or delegates
that selector, walks and sorts files, truncates each file and the total bundle,
hashes the package, and constructs one immutable area context. Requirement
sessions run in a neutral temporary directory with no tools and can judge only
that package. This gives the runner excellent control over bytes but poor control
over relevance: traversal order and area-level caps decide what the requirement
judge is allowed to know.

The current direct OpenAI and Anthropic adapters avoid an external agent runtime
but have no principled way to add discovery. Once tools are added, those adapters
would also need an agent loop, sandbox, authorization policy, context and
compaction behavior, tool-result normalization, cancellation, and provider
parity. That is a second coding-agent harness inside `qualitymd`.

The design makes the runner an evaluation harness without making it an agent
harness. Supported coding-agent SDKs own the iterative loop. The runner supplies
the task and authority boundary, then validates and seals the evidence used.

## Approach

### Revised execution boundary

The runtime flow becomes:

```text
QUALITY.md + requested scope
          |
          v
 deterministic plan and frames
          |
          v
 fresh requirement inspection session ----> read-only workspace
          |                                  optional sandboxed verification
          v
 judgment + proposed evidence manifest
          |
          v
 runner validation and evidence sealing
          |
          v
 atomic evaluation.json result
          |
          v
 tools-off synthesis, ranking, advice, reports
```

The evaluator is free to search broadly inside the authorized workspace for
supporting context, but the effective area `source` remains the judgment target.
For path and glob selectors, `evaluated` evidence must fall under the selected
material while `supporting` evidence may be elsewhere. For prose selectors, the
evaluator records the material it treated as the target and explains any limit;
all locators remain inside the authorized workspace in this change.

### Simplified work graph

`resolveSource` disappears. The graph retains deterministic framing and
downstream dependencies:

```text
frameEvaluation
  -> frameAreaEvaluation
    -> frameRequirementEvaluation
      -> assessRateRequirement
        -> analyzeFactor
          -> analyzeArea
            -> rankFindings
              -> recommend
                -> rankRecommendations
                  -> buildReports
```

Each `assessRateRequirement` unit opens one fresh session. Requirements remain
parallelizable up to the resolved concurrency cap. Factor and area dependencies
remain unchanged except that they wait directly on accepted requirement results
instead of an area source unit. Synthesis and advice requests receive persisted
payloads only and have no workspace tools.

Removing area-context construction also removes the misleading cache boundary.
The stable prefix of a requirement request can still contain task instructions,
schema, model vocabulary, area frame, and rating scale for provider caching, but
no reusable provider conversation or static evidence package is shared between
requirements.

### Inspection session abstraction

A new `InspectionSession` adapter boundary replaces `SourceBundle` and
`AreaContext`. Its input is a `RequirementInspectionRequest` containing:

- run, work-unit, area, factor, and requirement identity;
- the effective source selector and source form;
- the requirement frame, body guidance, and applied rating criteria;
- the structured output schema;
- an absolute authorized workspace root represented to the model as an explicit
  data boundary;
- the read, search, verification, network, write, time, turn, and output policy;
  and
- accepted upstream data needed for the requirement, but no preselected source
  content.

Its output is a `RequirementInspectionResponse` with a combined assessment and
rating payload, a proposed evidence manifest, evaluator identity and model,
usage when available, and non-sensitive session metadata.

The abstraction deliberately does not expose generic chat turns to application
code. Each provider adapter translates one bounded inspection request into its
SDK, streams cancellation and observable tool metadata, and returns one typed
response. Context management and compaction remain SDK behavior.

### SDK policy profiles

Both built-in agent adapters must prove the same policy capabilities, even when
their SDK options differ.

| Concern       | Codex adapter                                                                             | Claude adapter                                                                                                  |
| ------------- | ----------------------------------------------------------------------------------------- | --------------------------------------------------------------------------------------------------------------- |
| Neutral start | Temporary working directory; workspace supplied as an additional read-only directory      | Temporary working directory; workspace supplied as an additional directory; filesystem setting sources disabled |
| Instructions  | Runner-owned system/task context; repository instructions are not discovered as authority | Runner-owned system/task context; `settingSources: []`; project skills, hooks, and `CLAUDE.md` loading disabled |
| Workspace     | SDK `read-only` sandbox, explicit workspace root                                          | Restricted read/search tools and SDK sandbox with workspace writes denied                                       |
| Network       | Disabled                                                                                  | Disabled                                                                                                        |
| Approval      | `never`                                                                                   | Non-interactive deny outside the declared tool policy                                                           |
| Verification  | Unavailable until the SDK exposes equivalent policy and sealable event metadata           | Unavailable until the SDK exposes equivalent policy and sealable event metadata                                 |
| Output        | JSON Schema constrained                                                                   | JSON Schema constrained                                                                                         |
| Cancellation  | Abort the active thread/session and await cleanup                                         | Abort the active query/session and await cleanup                                                                |

The adapter readiness probe becomes capability-based. A runtime is usable only
when it can provide structured output, fresh context, a neutral instruction
boundary, read-only workspace access, disabled network, non-interactive policy,
cancellation, and evidence-bearing responses. Both adapters currently report
verification as unavailable: a requirement that needs it becomes limited,
partially assessed, or not rated rather than receiving shell access.

The implementation acceptance work includes adversarial workspace fixtures with
`AGENTS.md`, `CLAUDE.md`, local settings, hooks, skills, symlinks, and prompt-like
file content. The fixtures must prove that these files can be inspected as data
but do not change tool authority or evaluator instructions. If the pinned SDK
cannot establish that boundary, its adapter is incompatible until the SDK is
upgraded or the adapter can supply an equivalent isolated configuration.

### Evidence proposal and sealing

Evidence selection is agentic; evidence acceptance is deterministic. A
requirement response proposes records with stable response-local IDs. The
runner normalizes and seals them before it accepts the paired payloads.

The schema-version-8 artifact replaces top-level `sources` with top-level
`evidence`, keyed by requirement work-unit ID. A representative record is:

```json
{
  "evidence": {
    "assessRateRequirement:requirement:registry-auth": {
      "requirementId": "requirement:registry-auth",
      "source": { "selector": "apis/registry", "kind": "path" },
      "observations": [
        {
          "id": "ev-001",
          "kind": "file",
          "role": "evaluated",
          "path": "apis/registry/auth.ts",
          "locator": { "startLine": 18, "endLine": 61 },
          "sha256": "...",
          "capturedAt": "2026-07-14T00:00:00Z"
        },
        {
          "id": "ev-002",
          "kind": "file",
          "role": "supporting",
          "path": "docs/design/authentication.md",
          "locator": { "heading": "Registry authentication" },
          "sha256": "...",
          "capturedAt": "2026-07-14T00:00:00Z"
        }
      ],
      "limits": [],
      "manifestHash": "..."
    }
  }
}
```

The exact locator union supports file paths with a line range or heading
anchor. The artifact stores locators and digests, not full file bodies. The
assessment's evidence entries use a manifest selector such as
`evidence[ev-001]` as `sourceRef`, alongside the evaluator's concise evidence
statement.

For file observations, the runner:

1. resolves the proposed path against the authorized root and rejects escapes,
   symlink escapes, absolute paths, URLs, and unreadable files;
2. validates that concrete `evaluated` evidence belongs to the source selector;
3. resolves and validates the proposed line or heading locator;
4. reads the selected material after the session, computes the canonical digest,
   and discards any evaluator-supplied digest that does not match;
5. confirms every finding `sourceRef` names an observation in this manifest; and
6. hashes the canonical manifest and persists it in the same serialized store
   write as the assessment and rating.

Command observations are not accepted in this release because neither adapter
provides a proven mediated event path with the required isolation and
provenance. The request says verification is unavailable; evaluators record the
resulting limit. A future command-observation union requires a new durable
contract and runner-observed metadata rather than model-authored provenance.

Malformed or unverifiable provenance fails result acceptance as
`evidence_invalid`. Missing or inconclusive evidence that the evaluator reports
honestly remains a blocked, partially assessed, or not-rated domain outcome.
The old `missing_api_key`, `selector_unsupported`, and
`insufficient_evidence` infrastructure failures disappear; authentication
readiness remains `evaluator_unauthenticated`.

### Harness checkpoints

Harness transport keeps the existing multi-outstanding checkpoint and partial
submission model. It emits `RequirementInspectionRequest` envelopes directly;
there is no gather checkpoint. The `/quality` skill uses its own authorized
workspace tools to perform the same discovery and returns the same combined
payload plus proposed manifest.

The request hash covers all immutable request fields and policy but cannot cover
evidence that has not yet been discovered. On submission, evidence sealing
produces the response and manifest hashes stored with the accepted unit. The
first result still binds the run to one harness identity, and a resumed caller
must match it. Rejected or missing submissions remain retryable without
discarding other accepted results.

### Persistence and resume

`evaluation.json` advances directly from schema version 7 to 8. The new envelope
contains `manifest`, `state`, `evidence`, `results`, and `outputs`. There is no
dual `sources`/`evidence` shape and no reader for in-progress version-7 runs.
The CLI reports that an incompatible run must be restarted.

An accepted requirement unit is durable only after its normalized payloads and
sealed evidence manifest land in one atomic write. Resume treats both as an
inseparable immutable input to synthesis. A failed or cancelled requirement
without an accepted result opens a new session on retry. Pending harness calls
are rebuilt from the model snapshot and request fields; they never contain raw
workspace content.

The evidence record provides reproducibility of the claim and observation, not
a promise that a new agent run will choose the same context or rating. Capture
times and digests make workspace drift visible. A completed result remains the
record of what was accepted at that time; the runner does not silently replace
its evidence during resume.

### Evaluator and configuration cutover

The public runnable set becomes:

- `harness` — explicit outer-agent checkpoint transport;
- `codex` — Codex SDK plus a ready authenticated Codex runtime;
- `claude` — Claude Agent SDK plus a ready authenticated Claude runtime; and
- `auto` — deterministic readiness selection, Codex then Claude.

Configured profiles may set `kind: codex | claude`, `command`, and `model`, plus
future SDK-runtime options only when the contract needs them. The CLI removes
`openai`, `anthropic`, `shell`, and `manual` from built-in and reserved names;
removes `apiKeyEnv` and `baseUrl`; deletes direct HTTP evaluator adapters; and
stops searching configured API profiles during `auto` selection.

Authentication is observed only as runtime readiness. The inherited environment
may let an agent runtime use its documented login, subscription, or API-key
mechanism, but `qualitymd` neither names nor validates a provider key. Environment
sanitization ensures verification commands do not inherit those credentials.

### Dry-run, progress, and logs

Dry-run reports the selected evaluator/profile, readiness evidence, required and
optional inspection capabilities, concurrency, requirement count, and each
area's effective source selector. It no longer reports source kind dispatch,
resolver choice, file counts, byte caps, or bundle hashes.

Progress treats workspace inspection as part of each requirement call. Logs may
record work-unit identity, evaluator/model, duration, usage, tool names, command
metadata, evidence counts and roles, manifest hash, retry, cancellation, and
failure category. They do not record prompts, file bodies, raw tool output,
hidden reasoning, full transcripts, or credentials.

### Implementation and acceptance sequence

The implementation proceeds as one clean contract cutover:

1. update the cumulative format, evaluation, CLI, and skill specs to the new
   ownership and artifact contract;
2. replace source-package domain types and graph units with inspection requests,
   responses, and evidence manifests;
3. implement evidence sealing and schema-version-8 persistence before accepting
   agent-produced evidence;
4. configure and prove the neutral Codex and Claude SDK session policies,
   explicitly reporting verification unavailable;
5. update harness checkpoints and the runtime skill;
6. remove direct API adapters, profile fields, evaluator names, failures, and
   tests rather than retaining compatibility branches;
7. update reports, generated schemas/examples/docs, CLI output, and release
   notes; and
8. run the full gate plus authenticated end-to-end evaluations against a fixture
   whose requirements need different evidence from one shared area.

The acceptance fixture deliberately orders irrelevant files before relevant
ones, includes more material than the former bundle cap, places supporting
context outside the area's concrete source, and contains hostile repository
instructions. Codex, Claude, and harness runs must each discover
requirement-specific evidence, keep supporting context classified, reject
unsafe access, persist valid manifests, record verification limits when needed,
complete tools-off synthesis, and resume without regathering accepted work.

## Spec response

- **Ownership and source semantics (R1–R2):** the simplified flow makes context
  discovery part of requirement judgment while keeping `source` as the target
  and preserving runner authority over artifacts.
- **Sessions and safety (R3–R4):** fresh SDK sessions, neutral temporary roots,
  explicit read-only workspace access, disabled settings/network/writes, and a
  mediated verification path constrain agentic inspection.
- **Evidence and orchestration (R5–R6):** proposed manifests plus deterministic
  sealing replace static bundles; downstream work remains tools-off and consumes
  only accepted results.
- **Methods and harness (R7–R8):** Codex and Claude SDKs supply agent loops, the
  invoking skill supplies the explicit harness loop, and no raw API evaluator
  remains.
- **Artifacts and UX (R9–R10):** schema version 8, immutable accepted manifests,
  honest determinism language, and synchronized skill/CLI guidance create one
  comprehensible clean-break contract.

## Alternatives

### Keep deterministic source bundles but increase their caps

Rejected because a larger area-wide package still chooses evidence before the
requirement is interpreted, wastes context on irrelevant material, and fails on
cross-references. It makes the symptom less frequent without correcting
ownership.

### Let the runner gather a requirement-specific bundle

Rejected because the runner would still need semantic search and iterative
inspection. Encoding that behavior in deterministic collection rules either
recreates an agent or becomes a growing set of repository-specific heuristics.

### Add tools to direct model API evaluators

Rejected because tool schemas are the small part of the work. The CLI would own
the tool loop, authorization, context lifecycle, compaction, retries, provider
differences, and sandbox—precisely the agent-runtime mechanics the SDKs provide.

### Give agent sessions the workspace as their normal project root

Rejected because provider runtimes may auto-load repository instructions,
settings, skills, hooks, or memory. A neutral root with the repository mounted
as explicitly authorized data makes the trust boundary testable.

### Trust evaluator-cited paths without runner validation

Rejected because an agent could cite an escaped path, stale content, an
unobserved command, or a fabricated locator. Deterministic sealing is the point
where semantic selection becomes durable evidence.

### Persist complete inspected content and transcripts

Rejected because it expands artifact size and privacy exposure and may capture
secrets or irrelevant source. Validated locators, hashes, concise statements,
and command digests provide audit provenance without turning the run artifact
into a repository snapshot or reasoning log.

## Trade-offs and risks

- A fresh agent session per requirement costs more time and tokens than one
  shared area package. Bounded parallelism, cacheable stable prompt prefixes,
  and avoiding irrelevant bundle tokens offset part of that cost.
- Evidence selection and ratings are intentionally less repeatable at the byte
  level. The runner can reproduce orchestration and preserve what was accepted,
  but it cannot promise identical judgment from an agentic process.
- Locator-and-hash provenance is compact but does not preserve the underlying
  content after the workspace changes. Historical review may need the matching
  repository revision; reports still retain the evaluator's evidence statement.
- Neutral instruction isolation depends on provider SDK behavior. Capability
  tests and fail-closed readiness avoid silently running with a weaker boundary,
  but an SDK upgrade may be required before one adapter can ship.
- Sandboxed verification is less portable than static reading. Unsupported or
  unsafe commands reduce assessment completeness rather than triggering a
  permissive fallback.
- A user can mutate the workspace between an SDK read and runner sealing. The
  post-session digest records the accepted state and rejects a claimed digest
  mismatch; it cannot lock arbitrary external edits.
- Removing raw API methods narrows headless deployment choices. Users retain the
  supported SDK runtimes, whose own authentication may still use API keys, and
  the explicit harness transport.

## Decisions closed during implementation

The pinned SDKs can expose the workspace from a neutral temporary root with
repository settings disabled or non-authoritative, read/search-only tools,
disabled network, and no approval escalation. Neither provides a command path
with sufficient cross-provider isolation and runner-observed provenance, so
verification is explicitly unavailable rather than weakening the boundary.

File locators support whole-file references, 1-based line ranges, and exact
Markdown heading anchors. The runner computes whole-file SHA-256 and byte count
after the session and hashes the canonical manifest.
