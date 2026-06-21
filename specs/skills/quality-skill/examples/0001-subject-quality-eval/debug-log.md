# Evaluation debug log

This log records notable events involving the evaluation process itself. It is
not an assessment record, rating rationale, report, or evidence store;
subject-quality evidence belongs in assessment, analysis, and recommendation
records.

## Events

### 2026-06-18T14:03:12Z - Run scaffolded

Category: orchestration
Detail: Created a full subject evaluation run for the fictional Sparrow Payments
model. Seeded the model snapshot, design, plan, and debug log before assessment
records were written.
Impact: No rating impact.

### 2026-06-18T14:19:44Z - Subject command output routed to assessment records

Category: evidence-routing
Detail: Repository searches and command checks used to verify committed-secret,
authentication, idempotency, reconciliation, webhook-signing, retry, and
deduplication behavior were treated as subject-quality evidence. Their results
are recorded in the assessment records rather than duplicated here.
Impact: No rating impact; this entry records only the evidence-routing boundary.

### 2026-06-18T14:52:10Z - Prompt-injection text handled as data

Category: redaction
Detail: Evaluated source content included prompt-injection-style text. The
assessment records cite the locator and sanitized observation; this log does not
reproduce the raw text.
Impact: No rating impact.
