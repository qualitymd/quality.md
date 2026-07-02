package main

// The synthetic evaluation content: requirement assessments and ratings,
// factor and area analyses, and advice. Copy follows the evaluate workflow's
// finding-core field jobs; the synthetic-evidence caveat lives in the
// evaluation frame limits and the gallery README, not in the finding copy.

type criterionCase struct {
	Level string
	Text  string
}

type criteriaResultCase struct {
	Level     string
	Matched   bool
	Rationale string
}

type limitCase struct {
	ID          string
	Description string
	Impact      string
}

type findingCase struct {
	ID, Type, Severity, Confidence                string
	Statement, Condition                          string
	CriterionLevel, Criterion, CriterionRationale string
	BasisStatus, Basis                            string
	Effect, RatingEffect                          string
	EvidenceRef, Evidence                         string
	Driver                                        bool
}

type requirementCase struct {
	Area, Name, Title            string
	Factors                      []string
	AssessStatus, StatusReason   string
	Summary                      string
	Confidence, ConfidenceReason string
	RatingStatus, Rating         string
	RatingRationale              string
	AppliedCriteria              []criterionCase
	CriteriaResults              []criteriaResultCase
	EvidenceQuestion             string
	EvidenceRefs                 []string
	Limits                       []limitCase
	MissingEvidence              []limitCase
	Findings                     []findingCase
}

type factorCase struct {
	Area, Path, Title           string
	Children                    []string
	Status, StatusReason        string
	Rating, Confidence, Summary string
	Limits                      []limitCase
}

type areaCase struct {
	Name, Source                string
	Children                    []string
	Rating, Confidence, Summary string
	Limits                      []limitCase
	MissingEvidence             []limitCase
}

type recommendationCase struct {
	ID, Title, Description, Background string
	ExpectedValue, DoneCriterion       string
	Impact, Confidence                 string
	Traces                             []string
}

type rankedFindingCase struct {
	FindingID, Tier, Rationale string
}

var defaultMinimumCriterion = criterionCase{
	Level: "minimum",
	Text:  "Falls short of target with visible gaps or limited evidence, while remaining acceptable.",
}

var requirements = []requirementCase{
	// Root sub-factor requirements (agent harnessability).
	{
		Area: "", Name: "agents-reach-service-context-from-a-stable-entry-point",
		Title:        "a fresh agent reaches service context and deeper guidance from a stable entry point",
		Factors:      []string{"agent-harnessability/agent-accessibility"},
		AssessStatus: "assessed",
		Summary:      "The entry point routes a fresh agent to the contract, runbooks, and sensor catalog in one hop each.",
		Confidence:   "high", ConfidenceReason: "The entry point and its routed guidance were walked end to end from a clean session.",
		RatingStatus: "rated", Rating: "target",
		RatingRationale: "A fresh session reached the service purpose, contract, runbooks, and sensors without private context, meeting the target bar.",
		AppliedCriteria: []criterionCase{
			{Level: "target", Text: "A fresh agent reaches service purpose, contract, runbooks, and sensors from the entry point without private context."},
			defaultMinimumCriterion,
		},
		CriteriaResults: []criteriaResultCase{
			{Level: "target", Matched: true, Rationale: "Every routed destination resolved and matched its description."},
			{Level: "minimum", Matched: true, Rationale: "The entry point exists and routes."},
		},
		EvidenceQuestion: "Can a fresh agent session reach decision-relevant service context from the recorded entry point?",
		EvidenceRefs:     []string{"synthetic-source:agent-harness/entry-point"},
		Findings: []findingCase{{
			ID: "strength-001", Type: "strength", Confidence: "high", Driver: true,
			Statement:      "The agent entry point reaches the contract, runbooks, and sensor catalog in one hop each.",
			Condition:      "Each routed link resolved from a clean session, and the entry point stays under a page while linking deeper material.",
			CriterionLevel: "target", Criterion: "A fresh agent reaches service purpose, contract, runbooks, and sensors from the entry point without private context.",
			CriterionRationale: "Reachability from a cold start is the accessibility bar the model sets.",
			BasisStatus:        "verified", Basis: "The walk-through from a clean checkout reached every routed destination.",
			Effect: "Agents can orient without tribal knowledge, which supports the target accessibility rating.", RatingEffect: "supports target",
			EvidenceRef: "synthetic-source:agent-harness/entry-point",
			Evidence:    "The entry point names the service purpose and links the contract, runbooks, and sensor catalog; all links resolved.",
		}},
	},
	{
		Area: "", Name: "quality-loop-work-items-carry-done-criteria",
		Title:        "quality-loop work items carry scoped goals and done criteria",
		Factors:      []string{"agent-harnessability/task-specifiability"},
		AssessStatus: "assessed",
		Summary:      "Recent handoffs state goals, but most omit done criteria and the confirming sensor.",
		Confidence:   "high", ConfidenceReason: "The twelve most recent tracker handoffs were read directly.",
		RatingStatus: "rated", Rating: "minimum",
		RatingRationale: "Goals are stated, so work is startable, but missing done criteria leave completion a judgment call — short of target.",
		AppliedCriteria: []criterionCase{
			{Level: "target", Text: "Sampled handoffs state goals, non-goals, done criteria, and the sensor that confirms completion."},
			defaultMinimumCriterion,
		},
		CriteriaResults: []criteriaResultCase{
			{Level: "target", Matched: false, Rationale: "Nine of twelve sampled handoffs omit done criteria."},
			{Level: "minimum", Matched: true, Rationale: "Every sampled handoff states a scoped goal."},
		},
		EvidenceQuestion: "Do quality-loop handoffs give an agent explicit success and done criteria?",
		EvidenceRefs:     []string{"synthetic-source:tracker/quality-loop-handoffs"},
		Findings: []findingCase{{
			ID: "gap-001", Type: "gap", Severity: "medium", Confidence: "high", Driver: true,
			Statement:      "Most quality-loop handoffs omit done criteria and the confirming sensor.",
			Condition:      "Nine of the twelve most recent handoffs state a goal but no done criteria; none name the sensor that would confirm completion.",
			CriterionLevel: "target", Criterion: "Sampled handoffs state goals, non-goals, done criteria, and the sensor that confirms completion.",
			CriterionRationale: "Without done criteria an agent must guess when work is finished.",
			BasisStatus:        "verified", Basis: "The handoff sample was read in full; the omissions are countable.",
			Effect: "Agents declare completion on judgment rather than criteria, which constrains task specifiability to minimum.", RatingEffect: "constrains target",
			EvidenceRef: "synthetic-source:tracker/quality-loop-handoffs",
			Evidence:    "The sampled handoffs contain goal statements; done criteria appear in three, a confirming sensor in none.",
		}},
	},
	{
		Area: "", Name: "a-fresh-session-reaches-a-ready-to-work-environment",
		Title:        "a fresh agent session reaches a ready-to-work environment from recorded setup",
		Factors:      []string{"agent-harnessability/agent-operability"},
		AssessStatus: "assessed",
		Summary:      "Recorded setup builds the service and runs the sensors; one required credential is documented only in the private wiki.",
		Confidence:   "medium", ConfidenceReason: "The setup path was followed once from a clean checkout; the credential gap was confirmed against the recorded materials.",
		RatingStatus: "rated", Rating: "minimum",
		RatingRationale: "Setup works when the credential is already present, but a fresh session stalls on agent-inaccessible documentation, so the requirement sits at minimum.",
		AppliedCriteria: []criterionCase{
			{Level: "target", Text: "A fresh session reaches build, sensors, and required access from agent-accessible recorded materials alone."},
			defaultMinimumCriterion,
		},
		CriteriaResults: []criteriaResultCase{
			{Level: "target", Matched: false, Rationale: "The telemetry read credential is recorded only in the private wiki."},
			{Level: "minimum", Matched: true, Rationale: "Build and sensor steps are recorded and worked from a clean checkout."},
		},
		EvidenceQuestion: "Does the recorded setup take a fresh session to a ready-to-work environment?",
		EvidenceRefs:     []string{"synthetic-source:agent-harness/setup"},
		Findings: []findingCase{{
			ID: "gap-002", Type: "gap", Severity: "low", Confidence: "medium", Driver: true,
			Statement:      "The telemetry read credential needed for latency checks is documented only in the private wiki.",
			Condition:      "The recorded setup covers build and sensors, but the step granting telemetry read access points at a wiki page agents cannot reach.",
			CriterionLevel: "target", Criterion: "A fresh session reaches build, sensors, and required access from agent-accessible recorded materials alone.",
			BasisStatus: "verified", Basis: "The setup walk-through stalled at the credential step; every other step completed from recorded materials.",
			Effect: "Latency verification needs a human hand-off in an otherwise self-serve setup, constraining operability to minimum.", RatingEffect: "constrains target",
			EvidenceRef: "synthetic-source:agent-harness/setup",
			Evidence:    "The setup document's access section links the private wiki for the telemetry credential; no agent-accessible copy exists.",
		}},
	},
	{
		Area: "", Name: "handoffs-survive-session-loss",
		Title:        "in-flight work survives session loss through durable handoff records",
		Factors:      []string{"agent-harnessability/continuity"},
		AssessStatus: "assessed",
		Summary:      "Long-running work relies on chat scrollback; durable progress records exist for only one of four recent efforts.",
		Confidence:   "medium", ConfidenceReason: "Four recent multi-session efforts were traced; chat history could not be inspected directly, only its artifacts.",
		RatingStatus: "rated", Rating: "minimum",
		RatingRationale: "Work does resume today because the same people restart it, but the record trail would not survive a cold handoff — below target, above the floor.",
		AppliedCriteria: []criterionCase{
			{Level: "target", Text: "In-flight efforts keep durable records of decisions, remaining work, and verification status."},
			defaultMinimumCriterion,
		},
		CriteriaResults: []criteriaResultCase{
			{Level: "target", Matched: false, Rationale: "Three of four recent efforts left no durable progress record."},
			{Level: "minimum", Matched: true, Rationale: "Completed work landed with reviewable artifacts."},
		},
		EvidenceQuestion: "Would in-flight work survive losing the current session and its chat history?",
		EvidenceRefs:     []string{"synthetic-source:tracker/in-flight-work"},
		Findings: []findingCase{{
			ID: "risk-001", Type: "risk", Severity: "medium", Confidence: "medium", Driver: true,
			Statement:      "In-flight decisions live in chat scrollback and would not survive a cold handoff.",
			Condition:      "Three of four recent multi-session efforts kept decisions and remaining work only in conversation; one kept a durable progress note.",
			CriterionLevel: "target", Criterion: "In-flight efforts keep durable records of decisions, remaining work, and verification status.",
			BasisStatus: "plausible", Basis: "The missing records are confirmed; the impact projection assumes a handoff or compaction event occurs mid-effort.",
			Effect: "An interruption forces rediscovery of decisions already made, risking repeated or contradicted work.", RatingEffect: "constrains target",
			EvidenceRef: "synthetic-source:tracker/in-flight-work",
			Evidence:    "The traced efforts show tracker entries opened and closed with no interim state; the one durable note enabled a clean resume.",
		}},
	},
	{
		Area: "", Name: "sensors-return-pass-fail-with-remediation",
		Title:        "recorded sensors return objective pass/fail with remediation-bearing output",
		Factors:      []string{"agent-harnessability/self-verifiability"},
		AssessStatus: "assessed",
		Summary:      "The contract tests, invariant suite, and lint all run from the recorded commands and fail with the violated expectation named.",
		Confidence:   "high", ConfidenceReason: "Each sensor was run to both pass and forced-fail outcomes.",
		RatingStatus: "rated", Rating: "target",
		RatingRationale: "Every recorded sensor returned deterministic results whose failures name the violated expectation and point at the fix, meeting the target bar.",
		AppliedCriteria: []criterionCase{
			{Level: "target", Text: "Each recorded sensor returns deterministic pass/fail and failures name the violated expectation."},
			defaultMinimumCriterion,
		},
		CriteriaResults: []criteriaResultCase{
			{Level: "target", Matched: true, Rationale: "All three sensors met the bar under pass and forced-fail runs."},
			{Level: "minimum", Matched: true, Rationale: "Sensors exist and run from recorded commands."},
		},
		EvidenceQuestion: "Can an agent confirm its own work through the recorded sensors?",
		EvidenceRefs:     []string{"synthetic-source:agent-harness/sensor-catalog"},
		Findings: []findingCase{{
			ID: "strength-002", Type: "strength", Confidence: "high", Driver: true,
			Statement:      "All three recorded sensors are deterministic and their failures name the violated expectation.",
			Condition:      "The contract tests, ledger invariant suite, and lint each ran from the catalog command to a pass, and a seeded defect produced a failure naming the broken expectation and its location.",
			CriterionLevel: "target", Criterion: "Each recorded sensor returns deterministic pass/fail and failures name the violated expectation.",
			CriterionRationale: "Remediation-bearing failures are what let an agent fix its own work without a human interpreting the output.",
			BasisStatus:        "verified", Basis: "Pass and forced-fail runs were executed for each sensor.",
			Effect: "Agents can close their own verify loop, which supports the target self-verifiability rating.", RatingEffect: "supports target",
			EvidenceRef: "synthetic-source:agent-harness/sensor-catalog",
			Evidence:    "Sensor catalog commands ran as recorded; the seeded invariant defect failed with the invariant name, the offending entry pair, and the suite file to consult.",
		}},
	},
	{
		Area: "", Name: "standards-gate-nonconforming-changes",
		Title:        "core service standards are enforced by merge gates, not advisory prose",
		Factors:      []string{"agent-harnessability/enforcement-of-standards"},
		AssessStatus: "assessed",
		Summary:      "Lint gates merges; the contract tests and invariant suite run on merge but are advisory.",
		Confidence:   "high", ConfidenceReason: "The merge pipeline configuration was inspected directly.",
		RatingStatus: "rated", Rating: "minimum",
		RatingRationale: "The strongest sensors run on every merge but cannot block one, so standards hold by convention rather than gate — minimum, not target.",
		AppliedCriteria: []criterionCase{
			{Level: "target", Text: "Contract tests, the invariant suite, and lint block nonconforming merges or route them through reviewable exceptions."},
			defaultMinimumCriterion,
		},
		CriteriaResults: []criteriaResultCase{
			{Level: "target", Matched: false, Rationale: "Contract tests and the invariant suite are advisory at merge time."},
			{Level: "minimum", Matched: true, Rationale: "Lint blocks merges and the other sensors run visibly."},
		},
		EvidenceQuestion: "Do the recorded sensors actually gate nonconforming changes?",
		EvidenceRefs:     []string{"synthetic-source:ci/merge-pipeline"},
		Findings: []findingCase{{
			ID: "gap-003", Type: "gap", Severity: "medium", Confidence: "high", Driver: true,
			Statement:      "Contract tests and the invariant suite run at merge time but cannot block a merge.",
			Condition:      "The merge pipeline marks both sensor jobs as non-required; two merged changes in the sampled month had a failing contract-test run.",
			CriterionLevel: "target", Criterion: "Contract tests, the invariant suite, and lint block nonconforming merges or route them through reviewable exceptions.",
			CriterionRationale: "An advisory sensor protects nothing when the pressure is on; the gate is what makes the standard hold regardless of who is merging.",
			BasisStatus:        "verified", Basis: "Pipeline configuration and the two failing-but-merged runs were inspected.",
			Effect: "Contract conformance depends on reviewer diligence rather than the gate, constraining enforcement to minimum.", RatingEffect: "constrains target",
			EvidenceRef: "synthetic-source:ci/merge-pipeline",
			Evidence:    "The pipeline config lists lint as required and the contract-test and invariant jobs as informational; merge history shows two failing contract-test runs that merged.",
		}},
	},
	{
		Area: "", Name: "consequential-actions-require-approval",
		Title:        "money-moving and schema-changing actions require human approval",
		Factors:      []string{"agent-harnessability/containment-of-action"},
		AssessStatus: "assessed",
		Summary:      "Sandbox policy and deploy configuration put money movement and production schema changes behind named approval gates.",
		Confidence:   "high", ConfidenceReason: "Policy and deploy configuration were inspected; the approval gates are declarative and current.",
		RatingStatus: "rated", Rating: "target",
		RatingRationale: "An unattended run cannot move money, alter production schemas, or widen its own permissions; every such path ends at an approval gate.",
		AppliedCriteria: []criterionCase{
			{Level: "target", Text: "Unattended runs cannot move money, change production schemas, or escalate permissions without an approval gate."},
			defaultMinimumCriterion,
		},
		CriteriaResults: []criteriaResultCase{
			{Level: "target", Matched: true, Rationale: "All consequential paths terminate at declared approval gates."},
			{Level: "minimum", Matched: true, Rationale: "Sandbox limits exist and are enforced."},
		},
		EvidenceQuestion: "Can an unattended agent run take a consequential action without approval?",
		EvidenceRefs:     []string{"synthetic-source:agent-harness/sandbox-policy"},
		Findings: []findingCase{{
			ID: "strength-003", Type: "strength", Confidence: "high", Driver: true,
			Statement:      "Money movement and production schema changes are behind declarative approval gates.",
			Condition:      "The sandbox allowlist excludes payment execution and production migration commands, and the deploy pipeline requires a named approver for both.",
			CriterionLevel: "target", Criterion: "Unattended runs cannot move money, change production schemas, or escalate permissions without an approval gate.",
			BasisStatus: "verified", Basis: "The allowlist and pipeline approval rules were read directly and match the recorded policy.",
			Effect: "Agent autonomy on routine work does not expose the consequential actions, supporting the target containment rating.", RatingEffect: "supports target",
			EvidenceRef: "synthetic-source:agent-harness/sandbox-policy",
			Evidence:    "The policy denies payment-execution and migration commands to unattended sessions; the deploy config lists required approvers for both action classes.",
		}},
	},
	// Public API.
	{
		Area: "api", Name: "idempotent-mutations",
		Title:        "mutation endpoints are idempotent under retry",
		Factors:      []string{"correctness"},
		AssessStatus: "assessed",
		Summary:      "Duplicate-key replay is specified and tested; replay after an interrupted write is neither specified nor tested.",
		Confidence:   "medium", ConfidenceReason: "The contract-test suite and contract were compared directly; interrupted-write behavior is untested, so its risk is inferred.",
		RatingStatus: "rated", Rating: "minimum",
		RatingRationale: "Straightforward duplicate retries are safe and proven, but the interrupted-write replay path is unspecified, leaving the requirement short of target.",
		AppliedCriteria: []criterionCase{
			{Level: "target", Text: "Contract tests prove duplicate, replayed, and interrupted-write retries produce exactly-once ledger effects."},
			defaultMinimumCriterion,
		},
		CriteriaResults: []criteriaResultCase{
			{Level: "target", Matched: false, Rationale: "Interrupted-write replay has no contract clause and no test."},
			{Level: "minimum", Matched: true, Rationale: "Duplicate-key replay is specified and its tests pass."},
		},
		EvidenceQuestion: "Do retried mutations produce exactly-once ledger effects in every replay case the contract and tests cover?",
		EvidenceRefs:     []string{"synthetic-source:api/contract-tests", "synthetic-source:service-contract"},
		Findings: []findingCase{
			{
				ID: "gap-004", Type: "gap", Severity: "high", Confidence: "medium", Driver: true,
				Statement:      "Replay behavior after an interrupted write is unspecified and untested.",
				Condition:      "The contract defines duplicate-key replay, and the contract tests cover it; neither addresses a retry that arrives after a write failed mid-transaction.",
				CriterionLevel: "target", Criterion: "Contract tests prove duplicate, replayed, and interrupted-write retries produce exactly-once ledger effects.",
				CriterionRationale: "Interrupted writes are exactly the case integrators retry hardest against.",
				BasisStatus:        "verified", Basis: "The contract's retry section and the contract-test case list were both searched; the interrupted-write case appears in neither.",
				Effect: "Integrators cannot know whether an interrupted-write retry double-posts, which constrains correctness to minimum.", RatingEffect: "constrains target",
				EvidenceRef: "synthetic-source:service-contract",
				Evidence:    "The retry section defines duplicate-key semantics only; the contract-test suite has no interrupted-write replay case.",
			},
			{
				ID: "note-001", Type: "note", Confidence: "high",
				Statement:      "Replay traffic is a meaningful share of mutation volume.",
				Condition:      "The four-week telemetry window shows about 4% of mutation requests carrying a previously seen idempotency key.",
				CriterionLevel: "target", Criterion: "Contract tests prove duplicate, replayed, and interrupted-write retries produce exactly-once ledger effects.",
				CriterionRationale: "Volume context sizes how often the unspecified path is exercised.",
				BasisStatus:        "verified", Basis: "The telemetry rollup for the window was queried directly.",
				Effect: "The unspecified replay path is exercised daily, not theoretically.", RatingEffect: "informs severity",
				EvidenceRef: "synthetic-source:telemetry/mutation-replays",
				Evidence:    "The replay-rate panel reports 3.8-4.3% of mutation requests reusing an idempotency key across the window.",
			},
		},
	},
	{
		Area: "api", Name: "predictable-error-contracts",
		Title:        "error responses are predictable for callers",
		Factors:      []string{"operability"},
		AssessStatus: "assessed",
		Summary:      "Validation, authorization, and conflict responses all use the contract's error envelope, and handlers match it.",
		Confidence:   "high", ConfidenceReason: "The handler matrix was compared against the contract envelope for every endpoint in the index.",
		RatingStatus: "rated", Rating: "target",
		RatingRationale: "Every sampled failure mode returns the documented envelope with a stable code and retryable flag, meeting the target bar.",
		AppliedCriteria: []criterionCase{
			{Level: "target", Text: "Common failure modes return the contract's error envelope with stable codes callers can branch on."},
			defaultMinimumCriterion,
		},
		CriteriaResults: []criteriaResultCase{
			{Level: "target", Matched: true, Rationale: "Handlers and contract agree across the endpoint index."},
			{Level: "minimum", Matched: true, Rationale: "An error envelope exists and is documented."},
		},
		EvidenceQuestion: "Do error responses match the contract envelope across common failure modes?",
		EvidenceRefs:     []string{"synthetic-source:service-contract", "synthetic-source:api/handler-matrix"},
		Findings: []findingCase{{
			ID: "strength-004", Type: "strength", Confidence: "high", Driver: true,
			Statement:      "Error responses follow one documented envelope across all sampled failure modes.",
			Condition:      "Validation, authorization, and conflict cases for every endpoint in the contract index return the envelope's code, message, and retryable fields.",
			CriterionLevel: "target", Criterion: "Common failure modes return the contract's error envelope with stable codes callers can branch on.",
			BasisStatus: "verified", Basis: "The handler matrix and contract envelope were compared endpoint by endpoint.",
			Effect: "Callers can branch on stable error codes, supporting the target operability rating.", RatingEffect: "supports target",
			EvidenceRef: "synthetic-source:api/handler-matrix",
			Evidence:    "The matrix maps each failure mode to the envelope fields; no handler deviates from the documented codes.",
		}},
	},
	{
		Area: "api", Name: "p99-latency-within-budget",
		Title:        "p99 mutation latency stays within budget",
		Factors:      []string{"performance"},
		AssessStatus: "assessed",
		Summary:      "The four-week telemetry window puts mutation p99 at 262 ms, inside the 300 ms target band.",
		Confidence:   "high", ConfidenceReason: "The telemetry rollup covers the full window with no gaps; the band boundaries are unambiguous.",
		RatingStatus: "rated", Rating: "target",
		RatingRationale: "Measured p99 of 262 ms lands inside the 300 ms target band but above the 200 ms outstanding band.",
		AppliedCriteria: []criterionCase{
			{Level: "outstanding", Text: "p99 at or under 200 ms over the window."},
			{Level: "target", Text: "p99 at or under 300 ms over the window."},
			{Level: "minimum", Text: "p99 at or under 500 ms over the window."},
			{Level: "unacceptable", Text: "p99 above 500 ms over the window."},
		},
		CriteriaResults: []criteriaResultCase{
			{Level: "outstanding", Matched: false, Rationale: "262 ms is above the 200 ms outstanding band."},
			{Level: "target", Matched: true, Rationale: "262 ms is inside the 300 ms target band."},
			{Level: "minimum", Matched: true, Rationale: "262 ms is inside the 500 ms floor."},
		},
		EvidenceQuestion: "Where does mutation p99 land against the rating bands over the representative window?",
		EvidenceRefs:     []string{"synthetic-source:telemetry/latency-rollup"},
		Findings: []findingCase{{
			ID: "strength-005", Type: "strength", Confidence: "high", Driver: true,
			Statement:      "Mutation p99 measured 262 ms over the representative four-week window.",
			Condition:      "The rollup query across all mutation endpoints reports p99 between 248 ms and 276 ms week over week, aggregating to 262 ms.",
			CriterionLevel: "target", Criterion: "p99 at or under 300 ms over the window.",
			CriterionRationale: "The band was recalibrated to 300 ms after the caching work landed; the measurement tests the new floor.",
			BasisStatus:        "verified", Basis: "The rollup query was run against the recorded window with no missing weeks.",
			Effect: "Latency sits comfortably inside the recalibrated target band, supporting the target performance rating.", RatingEffect: "supports target",
			EvidenceRef: "synthetic-source:telemetry/latency-rollup",
			Evidence:    "Weekly p99 values of 248, 259, 276, and 264 ms across the window; no week exceeded the target band.",
		}},
	},
	// Service contract.
	{
		Area: "service-contract", Name: "contract-covers-mutation-semantics",
		Title:        "the contract defines retry, idempotency, and error semantics for every mutation endpoint",
		Factors:      []string{"completeness"},
		AssessStatus: "assessed",
		Summary:      "Twelve of fourteen mutation endpoints have complete semantics; two lack replay clauses.",
		Confidence:   "high", ConfidenceReason: "The endpoint index makes the population countable; each entry was checked.",
		RatingStatus: "rated", Rating: "minimum",
		RatingRationale: "The contract covers most of the population, but two endpoints with undefined replay semantics are exactly the absences the completeness bar exists to catch.",
		AppliedCriteria: []criterionCase{
			{Level: "target", Text: "Every endpoint in the contract's index defines retry, idempotency, and error semantics."},
			defaultMinimumCriterion,
		},
		CriteriaResults: []criteriaResultCase{
			{Level: "target", Matched: false, Rationale: "Two of fourteen indexed mutation endpoints lack replay semantics."},
			{Level: "minimum", Matched: true, Rationale: "Twelve endpoints are fully specified and the absences are countable."},
		},
		EvidenceQuestion: "Does every indexed mutation endpoint define retry, idempotency, and error semantics?",
		EvidenceRefs:     []string{"synthetic-source:service-contract"},
		Findings: []findingCase{{
			ID: "gap-005", Type: "gap", Severity: "medium", Confidence: "high", Driver: true,
			Statement:      "Two mutation endpoints have no replay semantics in the contract.",
			Condition:      "The reversal and adjustment endpoints appear in the contract's endpoint index with request and response shapes but no retry or idempotency clauses.",
			CriterionLevel: "target", Criterion: "Every endpoint in the contract's index defines retry, idempotency, and error semantics.",
			CriterionRationale: "The index is the population, so the totality claim is checkable and the absences are exact.",
			BasisStatus:        "verified", Basis: "All fourteen indexed endpoints were checked clause by clause.",
			Effect: "Integrators retrying reversals or adjustments are working from silence, which constrains completeness to minimum.", RatingEffect: "constrains target",
			EvidenceRef: "synthetic-source:service-contract",
			Evidence:    "The reversal and adjustment entries end at the response shape; the other twelve entries carry retry, idempotency, and error clauses.",
		}},
	},
	{
		Area: "service-contract", Name: "contract-matches-shipped-behavior",
		Title:        "contract semantics match shipped handler behavior",
		Factors:      []string{"consistency"},
		AssessStatus: "assessed",
		Summary:      "Where the contract speaks, handlers agree; one deprecated response field still ships that the contract no longer documents.",
		Confidence:   "high", ConfidenceReason: "The latest contract-test sensor run and handler matrix cover every specified clause.",
		RatingStatus: "rated", Rating: "target",
		RatingRationale: "Every specified clause matches shipped behavior in the latest sensor run; the deprecated-field remainder is noted but does not contradict a clause.",
		AppliedCriteria: []criterionCase{
			{Level: "target", Text: "Every specified contract clause matches shipped handler behavior in the latest contract-test run."},
			defaultMinimumCriterion,
		},
		CriteriaResults: []criteriaResultCase{
			{Level: "target", Matched: true, Rationale: "No specified clause diverges from shipped behavior."},
			{Level: "minimum", Matched: true, Rationale: "Divergences, where historical, were resolved in the contract's favor."},
		},
		EvidenceQuestion: "Does shipped handler behavior match every clause the contract specifies?",
		EvidenceRefs:     []string{"synthetic-source:service-contract", "synthetic-source:api/contract-tests"},
		Findings: []findingCase{
			{
				ID: "strength-006", Type: "strength", Confidence: "high", Driver: true,
				Statement:      "Shipped behavior matches every specified contract clause in the latest sensor run.",
				Condition:      "The contract-test suite exercises each specified clause against the running service and passed in full on the most recent run.",
				CriterionLevel: "target", Criterion: "Every specified contract clause matches shipped handler behavior in the latest contract-test run.",
				BasisStatus: "verified", Basis: "The sensor run report was read alongside the clause list; coverage is one test per clause.",
				Effect: "Judgments made against the contract hold for the shipped service, supporting the target consistency rating.", RatingEffect: "supports target",
				EvidenceRef: "synthetic-source:api/contract-tests",
				Evidence:    "The latest run reports all clause-mapped cases passing; the suite's clause coverage table shows no unmapped clause.",
			},
			{
				ID: "note-002", Type: "note", Confidence: "medium",
				Statement:      "The deprecated balance_after field still ships but is no longer documented.",
				Condition:      "Mutation responses include balance_after, removed from the contract in the v1.4 revision; two integrators are known to still read it.",
				CriterionLevel: "target", Criterion: "Every specified contract clause matches shipped handler behavior in the latest contract-test run.",
				CriterionRationale: "The field contradicts no clause, but undocumented shipped behavior is drift in the making.",
				BasisStatus:        "verified", Basis: "Response samples and the contract revision history were compared.",
				Effect: "Silent behavior outside the contract invites accidental breakage when the field is eventually removed.", RatingEffect: "informational",
				EvidenceRef: "synthetic-source:api/handler-matrix",
				Evidence:    "Response samples carry balance_after; the contract's changelog shows its documentation removed in v1.4 without a removal plan for the field itself.",
			},
		},
	},
	// Ledger persistence.
	{
		Area: "persistence", Name: "balance-invariants",
		Title:        "ledger mutations preserve balance invariants",
		Factors:      []string{"integrity"},
		AssessStatus: "assessed",
		Summary:      "The property-based suite passes across all mutation paths and the nightly reconciliation sensor shows zero unexplained drift over the window.",
		Confidence:   "high", ConfidenceReason: "Two independent sensors — the invariant suite and the reconciliation job — agree over the full window.",
		RatingStatus: "rated", Rating: "outstanding",
		RatingRationale: "Both sensors pass with margin: the property suite covers concurrent and failure paths beyond the required cases, and reconciliation drift is zero across the window, exceeding the target bar.",
		AppliedCriteria: []criterionCase{
			{Level: "outstanding", Text: "Invariant and reconciliation sensors pass over the window with coverage beyond the required mutation paths."},
			{Level: "target", Text: "Invariant tests pass across debit, credit, failed-write, and concurrent paths, and reconciliation shows no unexplained drift."},
			defaultMinimumCriterion,
			{Level: "unacceptable", Text: "Any reproducible invariant violation or unexplained reconciliation drift, however small."},
		},
		CriteriaResults: []criteriaResultCase{
			{Level: "outstanding", Matched: true, Rationale: "Coverage extends to interleaved concurrent reversals; drift is zero, not merely explained."},
			{Level: "target", Matched: true, Rationale: "All required paths pass and reconciliation is clean."},
			{Level: "minimum", Matched: true, Rationale: "Exceeded."},
		},
		EvidenceQuestion: "Do the invariant and reconciliation sensors show balance preservation across all mutation paths?",
		EvidenceRefs:     []string{"synthetic-source:persistence/invariant-suite", "synthetic-source:persistence/reconciliation"},
		Findings: []findingCase{{
			ID: "strength-007", Type: "strength", Confidence: "high", Driver: true,
			Statement:      "Balance invariants hold across every tested mutation path with zero reconciliation drift.",
			Condition:      "The property-based suite passes for debit, credit, failed-write, and interleaved concurrent reversal paths, and the nightly reconciliation report shows zero unexplained drift for the four-week window.",
			CriterionLevel: "outstanding", Criterion: "Invariant and reconciliation sensors pass over the window with coverage beyond the required mutation paths.",
			CriterionRationale: "This is the model's veto requirement; the sharpened unacceptable band is why two independent sensors are consulted.",
			BasisStatus:        "verified", Basis: "Suite results and the reconciliation drift report were both read for the window.",
			Effect: "The failure the model most exists to prevent shows no signal on either sensor, supporting an outstanding integrity rating.", RatingEffect: "supports outstanding",
			EvidenceRef: "synthetic-source:persistence/reconciliation",
			Evidence:    "The drift report reads zero unexplained cents across all accounts for the window; the suite's concurrent-reversal properties passed 10,000 generated cases.",
		}},
	},
	{
		Area: "persistence", Name: "migration-rollback",
		Title:        "migrations have rollback paths rehearsed against the current schema",
		Factors:      []string{"recoverability"},
		AssessStatus: "assessed",
		Summary:      "The runbook has rollback steps, but the last rehearsal predates the two most recent schema migrations.",
		Confidence:   "medium", ConfidenceReason: "The rehearsal record and migration history are unambiguous; whether the steps still work is untested inference.",
		RatingStatus: "rated", Rating: "minimum",
		RatingRationale: "Rollback guidance exists and once worked, but it has not been rehearsed against the current schema, so recoverability rests on hope at exactly the migrations most likely to need it.",
		AppliedCriteria: []criterionCase{
			{Level: "target", Text: "The most recent rollback rehearsal is newer than the latest schema change."},
			defaultMinimumCriterion,
		},
		CriteriaResults: []criteriaResultCase{
			{Level: "target", Matched: false, Rationale: "The last rehearsal predates migrations 041 and 042."},
			{Level: "minimum", Matched: true, Rationale: "Documented rollback steps exist and matched the schema when last rehearsed."},
		},
		EvidenceQuestion: "Is the rollback path rehearsed against the current schema?",
		EvidenceRefs:     []string{"synthetic-source:persistence/migration-runbook"},
		Findings: []findingCase{{
			ID: "risk-002", Type: "risk", Severity: "medium", Confidence: "medium", Driver: true,
			Statement:      "Rollback steps have not been rehearsed since the two most recent schema migrations landed.",
			Condition:      "The runbook's rehearsal record is dated before migrations 041 (ledger partitioning) and 042 (currency precision), both of which touch tables the rollback steps reorganize.",
			CriterionLevel: "target", Criterion: "The most recent rollback rehearsal is newer than the latest schema change.",
			CriterionRationale: "The criterion was tightened to rehearsal recency after two incidents where documented steps failed against drifted schemas.",
			BasisStatus:        "plausible", Basis: "The staleness is verified; that the steps would fail is projected from the incidents that motivated the tightened criterion.",
			Effect: "A failed release over the partitioned tables could not be confidently rolled back, constraining recoverability to minimum.", RatingEffect: "constrains target",
			EvidenceRef: "synthetic-source:persistence/migration-runbook",
			Evidence:    "The runbook rehearsal log's latest entry predates the migration history's entries for 041 and 042.",
		}},
	},
	// Operations.
	{
		Area: "operations", Name: "customer-impact-telemetry",
		Title:        "health signals explain customer impact",
		Factors:      []string{"observability"},
		AssessStatus: "assessed",
		Summary:      "Dashboards-as-code map service symptoms to failed customer actions, and the deployed dashboards match the definitions.",
		Confidence:   "medium", ConfidenceReason: "Definitions and deployed dashboards were compared; live incident use was not observed in the window.",
		RatingStatus: "rated", Rating: "target",
		RatingRationale: "Each major symptom class has a panel expressing customer impact, and the deployed dashboards match the committed definitions, meeting the target bar.",
		AppliedCriteria: []criterionCase{
			{Level: "target", Text: "Committed dashboard definitions map symptom classes to failed customer actions and match what is deployed."},
			defaultMinimumCriterion,
		},
		CriteriaResults: []criteriaResultCase{
			{Level: "target", Matched: true, Rationale: "Definitions cover the symptom classes and deployment drift is zero."},
			{Level: "minimum", Matched: true, Rationale: "Dashboards exist and are maintained as code."},
		},
		EvidenceQuestion: "Can an operator read customer impact from the health dashboards during an incident?",
		EvidenceRefs:     []string{"synthetic-source:operations/dashboards-as-code"},
		Findings: []findingCase{{
			ID: "strength-008", Type: "strength", Confidence: "medium", Driver: true,
			Statement:      "Health dashboards express failed customer actions, not just service symptoms.",
			Condition:      "The committed definitions include failed-mutation, retry-exhaustion, and queue-delay panels each denominated in affected customer actions, and the deployed dashboards match the definitions.",
			CriterionLevel: "target", Criterion: "Committed dashboard definitions map symptom classes to failed customer actions and match what is deployed.",
			BasisStatus: "verified", Basis: "The definitions were diffed against the deployed dashboards; the panel inventory was walked class by class.",
			Effect: "Operators can answer \"who is hurt\" directly from the dashboards, supporting the target observability rating.", RatingEffect: "supports target",
			EvidenceRef: "synthetic-source:operations/dashboards-as-code",
			Evidence:    "Panel definitions express failure counts as failed customer actions per minute; the deployment diff is empty.",
		}},
	},
	{
		Area: "operations", Name: "recovery-drill-ownership",
		Title:        "recovery drills have current owners and recent practice records",
		Factors:      []string{"recoverability"},
		AssessStatus: "not_assessed",
		StatusReason: "The recovery calendar and the incident playbook name different current owners, and no drill record exists after the ownership change, so current ownership could not be established from the available records.",
		Summary:      "Ownership records conflict; the requirement could not be assessed from the available material.",
		Confidence:   "none", ConfidenceReason: "The conflicting records cannot support a judgment either way.",
		RatingStatus:    "not_rated",
		RatingRationale: "Not assessed is recorded rather than a low rating: the evidence is missing, not failing.",
		AppliedCriteria: []criterionCase{
			{Level: "target", Text: "The calendar and playbook agree on the current drill owner and a drill record exists within the last quarter."},
			defaultMinimumCriterion,
		},
		EvidenceQuestion: "Who currently owns recovery drills, and when was the last practice?",
		EvidenceRefs:     []string{"synthetic-source:operations/recovery-calendar", "synthetic-source:operations/incident-playbook"},
		MissingEvidence: []limitCase{{
			ID:          "conflicting-ownership-records",
			Description: "The recovery calendar names the platform team while the incident playbook names a former individual owner; neither record postdates the reorg.",
			Impact:      "Current ownership and drill recency cannot be verified until the records are reconciled.",
		}},
		Findings: []findingCase{{
			ID: "note-003", Type: "note", Confidence: "low",
			Statement:      "Recovery drill ownership records contradict each other.",
			Condition:      "The calendar's owner field names the platform team; the playbook's escalation section names an engineer who left the rotation in the spring reorg; no drill record exists after the reorg.",
			CriterionLevel: "target", Criterion: "The calendar and playbook agree on the current drill owner and a drill record exists within the last quarter.",
			CriterionRationale: "The contradiction itself is the observation; it blocks assessment rather than failing it.",
			BasisStatus:        "not_assessed", Basis: "The records conflict, so neither can serve as the basis for a judgment.",
			Effect: "The requirement is recorded as not assessed; restoring assessability is the actionable next step.", RatingEffect: "blocks rating",
			EvidenceRef: "synthetic-source:operations/recovery-calendar",
			Evidence:    "The calendar and playbook disagree on the owner, and the drill log's latest entry predates the reorg.",
		}},
	},
	// Agent harness.
	{
		Area: "agent-harness", Name: "harness-orients-agents-and-routes-to-sensors",
		Title:        "the harness orients agents and routes them to runnable sensors",
		Factors:      []string{"completeness", "coherence", "currentness", "assessability"},
		AssessStatus: "assessed",
		Summary:      "The entry point, routed guidance, and sensor catalog cover setup through handoff and agree with the contract and runbooks; one sensor command name is stale.",
		Confidence:   "high", ConfidenceReason: "Every harness artifact was read and every catalog command executed.",
		RatingStatus: "rated", Rating: "target",
		RatingRationale: "The harness covers the work lifecycle, contradicts none of the material it routes to, and its sensors run — one stale command name is a blemish, not a bar failure.",
		AppliedCriteria: []criterionCase{
			{Level: "target", Text: "Harness artifacts cover setup, scoped work, verification, and handoff; agree with the contract and runbooks; and route to sensors that run as recorded."},
			defaultMinimumCriterion,
		},
		CriteriaResults: []criteriaResultCase{
			{Level: "target", Matched: true, Rationale: "Coverage, coherence, and runnable sensors all held; the stale name is noted."},
			{Level: "minimum", Matched: true, Rationale: "The harness exists and routes correctly."},
		},
		EvidenceQuestion: "Do the checked-in harness artifacts orient an agent and route to sensors that actually run?",
		EvidenceRefs:     []string{"synthetic-source:agent-harness"},
		Findings: []findingCase{
			{
				ID: "strength-009", Type: "strength", Confidence: "high", Driver: true,
				Statement:      "The harness covers the agent work lifecycle and its guidance agrees with the material it routes to.",
				Condition:      "The entry point, routed guides, and sensor catalog address setup, scoped work, verification, and handoff, and spot-checks against the contract and runbooks found no contradictions.",
				CriterionLevel: "target", Criterion: "Harness artifacts cover setup, scoped work, verification, and handoff; agree with the contract and runbooks; and route to sensors that run as recorded.",
				BasisStatus: "verified", Basis: "All harness artifacts were read; routed references were followed to their targets.",
				Effect: "Agents get complete, coherent steering, supporting the target rating across the harness factors.", RatingEffect: "supports target",
				EvidenceRef: "synthetic-source:agent-harness",
				Evidence:    "The lifecycle coverage table maps each work phase to a harness artifact; contradiction spot-checks across contract and runbook references came back clean.",
			},
			{
				ID: "note-004", Type: "note", Confidence: "high",
				Statement:      "The sensor catalog's invariant-suite command name is one rename behind.",
				Condition:      "The catalog records the pre-rename command; the current command name differs, and the old name fails with a not-found error rather than a redirect.",
				CriterionLevel: "target", Criterion: "Harness artifacts cover setup, scoped work, verification, and handoff; agree with the contract and runbooks; and route to sensors that run as recorded.",
				CriterionRationale: "A stale command turns a routed sensor into a dead end for a fresh session.",
				BasisStatus:        "verified", Basis: "Both the recorded and current command names were executed.",
				Effect: "A fresh agent following the catalog hits a not-found error on the invariant suite until it discovers the rename.", RatingEffect: "informational",
				EvidenceRef: "synthetic-source:agent-harness/sensor-catalog",
				Evidence:    "The catalog entry's command fails with not-found; the renamed command runs the suite successfully.",
			},
		},
	},
	// QUALITY.md self-check.
	{
		Area: "quality-md", Name: "the-model-follows-the-authoring-guide-family",
		Title:        "the quality model follows its authoring guide family",
		Factors:      []string{"context-grounding", "evaluability", "lifecycle-maintenance"},
		AssessStatus: "assessed",
		Summary:      "The model's structure, traceability, and changelog follow the guides; the body's unknowns and open questions have not kept up with the model's own changes.",
		Confidence:   "high", ConfidenceReason: "The model file and quality changelog are fully inspectable.",
		RatingStatus: "rated", Rating: "minimum",
		RatingRationale: "Structure and evaluability meet the guides, but the body's judgment context has drifted from the current model, and a stale body misleads the next evaluator — minimum until refreshed.",
		AppliedCriteria: []criterionCase{
			{Level: "target", Text: "Body context, factor traceability, requirement assessability, and changelog practice match the authoring guide family with current unknowns and open questions."},
			defaultMinimumCriterion,
		},
		CriteriaResults: []criteriaResultCase{
			{Level: "target", Matched: false, Rationale: "The body's unknowns and open questions predate the service-contract area."},
			{Level: "minimum", Matched: true, Rationale: "Structure, assessability, and changelog practice follow the guides."},
		},
		EvidenceQuestion: "Does the model artifact follow the authoring guide family it cites?",
		EvidenceRefs:     []string{"./QUALITY.md"},
		Findings: []findingCase{{
			ID: "gap-006", Type: "gap", Severity: "medium", Confidence: "high", Driver: true,
			Statement:      "The body's unknowns and open questions have not been revisited since the service-contract area was added.",
			Condition:      "The Risks section still carries the pre-contract open question about error-envelope deprecation as unresolved even though the contract area now owns it, and the Needs unknowns omit the new integrator-retry blind spot the idempotency work surfaced.",
			CriterionLevel: "target", Criterion: "Body context, factor traceability, requirement assessability, and changelog practice match the authoring guide family with current unknowns and open questions.",
			CriterionRationale: "Unknowns are judgment context; stale ones misdirect the next evaluation's attention.",
			BasisStatus:        "verified", Basis: "The body's section dates and the changelog's area-addition entry were compared.",
			Effect: "The next evaluator inherits judgment context that no longer matches the model, constraining context grounding to minimum.", RatingEffect: "constrains target",
			EvidenceRef: "./QUALITY.md",
			Evidence:    "The Risks open question predates the contract area's changelog entry; no body section's review line postdates it.",
		}},
	},
}

var factors = []factorCase{
	// Root umbrella and sub-factors.
	{
		Area: "", Path: "agent-harnessability", Title: "Agent Harnessability",
		Children: []string{
			"agent-harnessability/agent-accessibility",
			"agent-harnessability/task-specifiability",
			"agent-harnessability/agent-operability",
			"agent-harnessability/continuity",
			"agent-harnessability/self-verifiability",
			"agent-harnessability/enforcement-of-standards",
			"agent-harnessability/containment-of-action",
		},
		Rating: "minimum", Confidence: "medium",
		Summary: "Sensors, accessibility, and containment equip agents well, but handoff records, done criteria, and advisory merge gates hold the equipping capability at minimum.",
	},
	{Area: "", Path: "agent-harnessability/agent-accessibility", Title: "Agent Accessibility", Rating: "target", Confidence: "high",
		Summary: "A fresh agent reaches decision-relevant context from the stable entry point."},
	{Area: "", Path: "agent-harnessability/task-specifiability", Title: "Task Specifiability", Rating: "minimum", Confidence: "high",
		Summary: "Handoffs scope the work but mostly omit done criteria and the confirming sensor."},
	{Area: "", Path: "agent-harnessability/agent-operability", Title: "Agent Operability", Rating: "minimum", Confidence: "medium",
		Summary: "Recorded setup works except for one credential documented outside agent reach."},
	{Area: "", Path: "agent-harnessability/continuity", Title: "Continuity", Rating: "minimum", Confidence: "medium",
		Summary: "In-flight work mostly lives in chat scrollback rather than durable progress records."},
	{Area: "", Path: "agent-harnessability/self-verifiability", Title: "Self-Verifiability", Rating: "target", Confidence: "high",
		Summary: "The recorded sensors are deterministic and their failures carry remediation."},
	{Area: "", Path: "agent-harnessability/enforcement-of-standards", Title: "Enforcement of Standards", Rating: "minimum", Confidence: "high",
		Summary: "Lint gates merges, but the contract and invariant sensors remain advisory."},
	{Area: "", Path: "agent-harnessability/containment-of-action", Title: "Containment of Action", Rating: "target", Confidence: "high",
		Summary: "Money movement and schema changes sit behind declarative approval gates."},
	// Public API.
	{Area: "api", Path: "correctness", Title: "Correctness", Rating: "minimum", Confidence: "medium",
		Summary: "Duplicate replay is proven safe; the unspecified interrupted-write path holds correctness at minimum."},
	{Area: "api", Path: "operability", Title: "Operability", Rating: "target", Confidence: "high",
		Summary: "Error responses follow one documented envelope callers can branch on."},
	{Area: "api", Path: "performance", Title: "Performance", Rating: "target", Confidence: "high",
		Summary: "Mutation p99 of 262 ms sits inside the recalibrated 300 ms target band."},
	// Service contract.
	{Area: "service-contract", Path: "completeness", Title: "Completeness", Rating: "minimum", Confidence: "high",
		Summary: "Two of fourteen mutation endpoints lack replay semantics; the rest are fully specified."},
	{Area: "service-contract", Path: "consistency", Title: "Consistency", Rating: "target", Confidence: "high",
		Summary: "Shipped behavior matches every specified clause; one undocumented deprecated field ships on."},
	// Persistence.
	{Area: "persistence", Path: "integrity", Title: "Integrity", Rating: "outstanding", Confidence: "high",
		Summary: "Two independent sensors show balance preservation with margin: full-path property coverage and zero reconciliation drift."},
	{Area: "persistence", Path: "recoverability", Title: "Recoverability", Rating: "minimum", Confidence: "medium",
		Summary: "Rollback guidance exists but is unrehearsed against the two most recent schema migrations."},
	// Operations.
	{Area: "operations", Path: "observability", Title: "Observability", Rating: "target", Confidence: "medium",
		Summary: "Dashboards-as-code express customer impact and match what is deployed."},
	{
		Area: "operations", Path: "recoverability", Title: "Recoverability",
		Status:       "blocked",
		StatusReason: "The factor's only requirement is not assessed: recovery-ownership records conflict, so no judgment could be formed.",
		Confidence:   "none",
		Summary:      "Drill ownership records contradict each other; the factor awaits reconciled evidence.",
		Limits: []limitCase{{
			ID:          "conflicting-ownership-records",
			Description: "The recovery calendar and incident playbook name different current owners.",
			Impact:      "The factor cannot be analyzed until ownership records are reconciled.",
		}},
	},
	// Agent harness.
	{Area: "agent-harness", Path: "completeness", Title: "Completeness", Rating: "target", Confidence: "high",
		Summary: "Harness artifacts cover setup, scoped work, verification, and handoff."},
	{Area: "agent-harness", Path: "coherence", Title: "Coherence", Rating: "target", Confidence: "high",
		Summary: "Guidance agrees with the contract and runbooks it routes to."},
	{Area: "agent-harness", Path: "currentness", Title: "Currentness", Rating: "target", Confidence: "medium",
		Summary: "Guidance matches the current layout except one stale sensor command name."},
	{Area: "agent-harness", Path: "assessability", Title: "Assessability", Rating: "target", Confidence: "high",
		Summary: "Harness quality is checkable through inspectable artifacts and runnable sensors."},
	// QUALITY.md self-check.
	{Area: "quality-md", Path: "context-grounding", Title: "Context Grounding", Rating: "minimum", Confidence: "high",
		Summary: "The body's unknowns and open questions have drifted behind the model's own changes."},
	{Area: "quality-md", Path: "evaluability", Title: "Evaluability", Rating: "minimum", Confidence: "high",
		Summary: "Requirements are assessable from recorded evidence; the stale body context caps the shared requirement's rating."},
	{Area: "quality-md", Path: "lifecycle-maintenance", Title: "Lifecycle Maintenance", Rating: "minimum", Confidence: "high",
		Summary: "The changelog records why the model changed; the body has not kept pace with those changes."},
}

var areas = []areaCase{
	{
		Name: "api", Source: "synthetic-source:api",
		Rating: "minimum", Confidence: "medium",
		Summary: "Error contracts and latency meet target, but unspecified interrupted-write replay holds the API at minimum.",
	},
	{
		Name: "service-contract", Source: "synthetic-source:service-contract",
		Rating: "minimum", Confidence: "high",
		Summary: "The contract matches shipped behavior where it speaks; two mutation endpoints without replay semantics cap completeness.",
	},
	{
		Name: "persistence", Source: "synthetic-source:persistence",
		Rating: "minimum", Confidence: "medium",
		Summary: "Balance integrity is outstanding on two independent sensors; unrehearsed rollback against the current schema caps the area at minimum.",
	},
	{
		Name: "operations", Source: "synthetic-source:operations",
		Rating: "target", Confidence: "low",
		Summary: "Customer-impact telemetry meets target; drill ownership could not be assessed, which limits confidence in the area's recoverability.",
		Limits: []limitCase{{
			ID:          "drill-ownership-unassessed",
			Description: "Recovery drill ownership is not assessed because its records conflict.",
			Impact:      "The area rating reflects observability evidence only; recoverability is missing, not weak.",
		}},
		MissingEvidence: []limitCase{{
			ID:          "reconciled-ownership-records",
			Description: "A reconciled owner record and a post-reorg drill entry are needed before recoverability can be judged.",
			Impact:      "Until then the area's recoverability contributes no rating signal.",
		}},
	},
	{
		Name: "agent-harness", Source: "synthetic-source:agent-harness",
		Rating: "target", Confidence: "high",
		Summary: "The harness orients agents and routes to runnable sensors, with one stale command name noted for repair.",
	},
	{
		Name: "quality-md", Source: "./QUALITY.md",
		Rating: "minimum", Confidence: "high",
		Summary: "The model's structure and changelog practice follow the guides; its body judgment context has drifted behind the model.",
	},
}

var rootLocalAnalysis = areaCase{
	Name: "", Rating: "minimum", Confidence: "medium",
	Summary: "Agent equipping is strong on sensors, accessibility, and containment, but handoff records, done criteria, and advisory merge gates sit below target.",
}

var rootAggregateAnalysis = areaCase{
	Name: "", Rating: "minimum", Confidence: "medium",
	Summary: "LedgerLite is money-safe today — balance integrity is outstanding — but unspecified replay semantics, unrehearsed rollback, and advisory merge gates hold the money-touching areas below the target margin the model's body requires.",
}

var recommendations = []recommendationCase{
	{
		ID:            "qrec_replaycontract",
		Title:         "Specify and test replay semantics for interrupted mutations",
		Description:   "Extend the service contract's retry section to define replay outcomes after interrupted writes for all fourteen mutation endpoints, including the reversal and adjustment endpoints that currently lack replay clauses, and add contract-test cases for duplicate, replayed, and interrupted-write retries.",
		Background:    "Interrupted-write replay is unspecified and untested while replay traffic runs near 4% of mutation volume, and two endpoints lack replay semantics entirely. Both findings trace to the same contract section, so one change closes both.",
		ExpectedValue: "Integrators and agents can verify retry behavior from the contract and its sensor instead of inferring undocumented recovery semantics; the API correctness and contract completeness ratings can reach target.",
		DoneCriterion: "Every endpoint in the contract's index defines retry, idempotency, and error semantics including interrupted-write replay, and the contract-test suite covers duplicate, replayed, and interrupted-write cases for each mutation endpoint.",
		Impact:        "high", Confidence: "high",
		Traces: []string{"gap-004", "gap-005"},
	},
	{
		ID:            "qrec_rollbackrehearsal",
		Title:         "Rehearse migration rollback against the current schema",
		Description:   "Run the runbook's rollback path against a copy of the current schema — including migrations 041 and 042 — and record the rehearsal with date, schema version, and any step corrections.",
		Background:    "The rollback steps have not been rehearsed since the partitioning and currency-precision migrations landed, and the criterion was tightened to rehearsal recency precisely because documented-but-unrehearsed steps failed in two past incidents.",
		ExpectedValue: "Release risk drops because rollback instructions are proven against the schema they would actually run on, and the recoverability rating can return to target.",
		DoneCriterion: "The runbook's rehearsal log contains an entry newer than migration 042 with the rehearsed steps matching the current runbook text.",
		Impact:        "high", Confidence: "medium",
		Traces: []string{"risk-002"},
	},
	{
		ID:            "qrec_gatesensors",
		Title:         "Make the contract-test and invariant sensors required at merge",
		Description:   "Change the merge pipeline to mark the contract-test and invariant-suite jobs as required, with a reviewable exception path for genuinely unrelated failures.",
		Background:    "Both sensors already run on every merge and their failures carry remediation, but they cannot block a merge; two failing contract-test runs merged in the sampled month. The sensors are trustworthy enough to gate on.",
		ExpectedValue: "Contract conformance and ledger invariants hold regardless of reviewer attention, converting the two strongest sensors from advisory signals into enforced standards.",
		DoneCriterion: "The merge pipeline lists both jobs as required, and a nonconforming test change is demonstrably blocked or routed through the recorded exception path.",
		Impact:        "high", Confidence: "medium",
		Traces: []string{"gap-003"},
	},
	{
		ID:            "qrec_drillownership",
		Title:         "Reconcile recovery drill ownership and restore assessability",
		Description:   "Agree the current drill owner, update the recovery calendar and incident playbook to name the same owner, and schedule the next drill so a post-reorg record exists.",
		Background:    "The ownership requirement could not be assessed because the calendar and playbook contradict each other and no drill record postdates the reorg. The blocker is record reconciliation, not drill quality.",
		ExpectedValue: "The recoverability factor becomes assessable again, and the next evaluation can rate drill practice on evidence instead of recording missing evidence.",
		DoneCriterion: "The calendar and playbook name the same current owner, and the drill log contains a post-reorg entry or a scheduled date.",
		Impact:        "medium", Confidence: "high",
		Traces: []string{"note-003"},
	},
	{
		ID:            "qrec_durablehandoffs",
		Title:         "Record done criteria and progress in durable handoff notes",
		Description:   "Add done criteria and the confirming sensor to the quality-loop handoff template, and keep in-flight decisions and remaining work in a durable progress note rather than chat scrollback.",
		Background:    "Most handoffs omit done criteria, and in-flight decisions live in conversation history that a session loss would erase. One template change and one recording habit address both weaknesses.",
		ExpectedValue: "Agents can declare completion against criteria and resume interrupted work from records, lifting task specifiability and continuity toward target.",
		DoneCriterion: "The handoff template carries done-criteria and sensor fields, and the next three multi-session efforts each leave a durable progress note that names decisions, remaining work, and verification status.",
		Impact:        "medium", Confidence: "medium",
		Traces: []string{"gap-001", "risk-001"},
	},
	{
		ID:            "qrec_refreshmodelbody",
		Title:         "Refresh the model body's unknowns and open questions",
		Description:   "Revisit each body section's unknowns and open questions against the current model — retiring the resolved error-envelope question now owned by the contract area, adding the integrator-retry blind spot — and advance the review provenance lines.",
		Background:    "The body's judgment context predates the service-contract area, so the next evaluator inherits attention pointers that no longer match the model.",
		ExpectedValue: "Evaluations start from current judgment context, and the QUALITY.md self-check can return to target.",
		DoneCriterion: "Every body section's unknowns and open questions reflect the current model, and each section's review line postdates the contract area's changelog entry.",
		Impact:        "low", Confidence: "high",
		Traces: []string{"gap-006"},
	},
}

var rankedFindings = []rankedFindingCase{
	{FindingID: "gap-004", Tier: "P1", Rationale: "Unspecified interrupted-write replay on a money-moving API is the highest-exposure gap, exercised daily by real replay traffic."},
	{FindingID: "risk-002", Tier: "P1", Rationale: "An unrehearsed rollback over freshly partitioned ledger tables is the failure mode two past incidents already demonstrated."},
	{FindingID: "gap-003", Tier: "P2", Rationale: "Advisory merge gates let the two strongest sensors be ignored under pressure."},
	{FindingID: "gap-005", Tier: "P2", Rationale: "Two endpoints with silent replay semantics share a root cause and a fix with the P1 replay gap."},
	{FindingID: "note-003", Tier: "P2", Rationale: "Missing evidence on drill ownership blocks the recoverability judgment entirely; restoring assessability is cheap."},
	{FindingID: "gap-001", Tier: "P3", Rationale: "Missing done criteria degrade every agent handoff, but each instance is recoverable."},
	{FindingID: "risk-001", Tier: "P3", Rationale: "Chat-bound progress records only hurt when an interruption lands mid-effort."},
	{FindingID: "gap-006", Tier: "P3", Rationale: "Stale body context misdirects the next evaluation but does not affect the service itself."},
	{FindingID: "gap-002", Tier: "P3", Rationale: "One wiki-bound credential is a small, well-understood setup snag."},
	{FindingID: "note-004", Tier: "P4", Rationale: "The stale sensor command name is a one-line catalog fix."},
	{FindingID: "note-001", Tier: "P4", Rationale: "Replay-volume context informs the P1 gap's severity; no separate action."},
	{FindingID: "note-002", Tier: "P4", Rationale: "The deprecated field is drift-in-waiting worth watching, not acting on yet."},
	{FindingID: "strength-007", Tier: "P4", Rationale: "Outstanding balance integrity is the evaluation's anchor strength."},
	{FindingID: "strength-002", Tier: "P4", Rationale: "Remediation-bearing sensors underpin agent self-verification."},
	{FindingID: "strength-003", Tier: "P4", Rationale: "Containment gates keep agent autonomy safe on consequential actions."},
	{FindingID: "strength-001", Tier: "P4", Rationale: "The entry point makes the rest of the harness reachable."},
	{FindingID: "strength-004", Tier: "P4", Rationale: "Stable error contracts support integrator trust."},
	{FindingID: "strength-005", Tier: "P4", Rationale: "Latency sits inside the recalibrated band."},
	{FindingID: "strength-006", Tier: "P4", Rationale: "Contract-to-behavior conformance holds where specified."},
	{FindingID: "strength-008", Tier: "P4", Rationale: "Dashboards answer customer impact directly."},
	{FindingID: "strength-009", Tier: "P4", Rationale: "Harness coverage and coherence hold up under inspection."},
}
