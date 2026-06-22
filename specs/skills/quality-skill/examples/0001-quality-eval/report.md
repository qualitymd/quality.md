# Evaluation Report

## Verdict

- **Root area:** Sparrow Payments API
- **Evaluation level:** not recorded
- **Rigor:** not recorded
- **Evaluation verdict:** Unacceptable
- **Rationale:** The root's own local rating binds. Ledger is rated target and Webhooks is rated minimum; if the local secret gap were cleared, Webhooks would become the next binding constraint.

## Scope

Evaluation scope reconstructed from the run's analysis and assessment result records.

- **Narrowing:** whole recorded run
- **In scope:** Sparrow Payments API; Ledger; Webhooks; Delivery
- **Out of scope:** not recorded
- **Metadata not recorded:** Rigor; Out-of-scope areas

## Selected Findings and Limitations

- `assessments/001-root-no-credentials-are-committed-to-the-repository.json` at `internal/gateway/client.go:48` [Critical]: A live payment-gateway API secret key is committed in plaintext; it matches the format of an active key and is not a placeholder.
- `assessments/005-ledger-reconciliation-runs-daily-and-flags-drift.json` at `ledger/reconcile.go:31` [Medium]: A reconcile entrypoint exists, but no reconciliation job output, log, or report was available.
- `assessments/008-webhooks-delivery-a-redelivery-of-an-already-acknowledged-event-is-suppressed-for-that-endpoint.json` at `webhooks/delivery/dedup.go:52` [Medium]: Suppression state is retained for only 24 hours, the same as the retry window, so a delivery re-enqueued late in that window can fire after its dedup record has expired.
- Limitation: Reconciliation runs daily and flags drift: Insufficient evidence to rate against the scale, so the requirement is recorded as not assessed rather than assigned a level
- Limitation: Insufficient evidence to rate against the scale, so the requirement is recorded as not assessed rather than assigned a level
- Limitation: The reconciliation requirement is not assessed, so the local rating rests on the assessed invariant evidence and is noted as incomplete rather than outstanding
- Limitation: The assessed double-entry invariant reaches target while reconciliation remains not assessed
- Limitation: The double-entry invariant holds in the recorded evidence; reconciliation is not assessed, so the factor rating is target but incomplete

## Evidence Basis

- **source:** `internal/gateway/client.go:48`
- **search:** `rg 'gateway secret key' internal/gateway/client.go`
- **search:** `rg 'gateway secret key' internal/gateway/client.go (rating-binding re-check)`
- **source:** `internal/gateway/client_test.go:12`
- **search:** `rg 'publishable key' internal/gateway/client_test.go`
- **source:** `cmd/api/routes.go`
- **search:** `rg 'RequireAuth|transfers|refunds|payouts' cmd/api/routes.go`
- **source:** `internal/gateway/client.go:51`
- **source:** `internal/transfer/handler.go:73`
- **search:** `rg 'Idempotency-Key' internal/transfer`
- **source:** `internal/transfer/handler_test.go:120`
- **test:** `go test ./internal/transfer -run TestTransferReplayReturnsOriginalResult`
- **source:** `ledger/posting.go:64`
- **search:** `rg 'sum|balanced|zero' ledger/posting.go`
- **source:** `ledger/posting_test.go:210`
- **test:** `go test ./ledger -run TestPostingPropertyBalancedTransfers`
- **source:** `ledger/reconcile.go:31`
- **search:** `rg 'reconcile' ledger`
- **search:** `rg 'reconciliation report|drift report|scheduled reconciliation' .`
- **source:** `webhooks/emit.go`
- **search:** `rg 'Sign|sign' webhooks/emit.go webhooks/sign.go`
- **source:** `webhooks/sign_test.go:48`
- **test:** `go test ./webhooks -run TestWebhookSignatureRejectsTamper`
- **source:** `webhooks/delivery/retry.go:37`
- **search:** `rg 'backoff|retry window|24h' webhooks/delivery/retry.go`
- **source:** `webhooks/delivery/retry_test.go:91`
- **test:** `go test ./webhooks/delivery -run TestRetryTransientFailure`
- **source:** `webhooks/delivery/dedup.go:52`
- **search:** `rg 'dedup' webhooks/delivery`
- **search:** `rg '24h|retention' webhooks/delivery`
- **source:** `docs/merchants/webhooks.md`
- **search:** `rg 'idempotent|duplicate' docs/merchants/webhooks.md`

## Next Action

- [001-rotate-committed-gateway-key](recommendations/001-rotate-committed-gateway-key.md) - The committed-credentials requirement reaches target; no live credential is present in the working tree and the previously exposed key has been revoked.

## Area Breakdown

| Area | Path | Area Rating | Area + Sub-Areas Rating | Factors |
| --- | --- | --- | --- | --- |
| Sparrow Payments API | `/` | Unacceptable | Unacceptable | Security: Unacceptable; Reliability: Target |
| Ledger | `ledger` | Target | Target | Correctness: Target |
| Webhooks | `webhooks` | Target | Minimum | Security: Target |
| Delivery | `webhooks/delivery` | Minimum | Minimum | Reliability: Minimum |

## Area Details

### Sparrow Payments API

- **Path:** /
- **Area rating:** Unacceptable
  - Over the root's three own requirements, the committed-credential shortfall is security-critical and holds the local rating at the floor regardless of the two requirements rated target.
- **+ Sub-Areas rating:** Unacceptable
  - The root's own local rating binds. Ledger is rated target and Webhooks is rated minimum; if the local secret gap were cleared, Webhooks would become the next binding constraint.
- **Factor Security:** Unacceptable
  - Bound by the Secrets handling sub-factor because a live gateway credential is committed in plaintext; Access control is rated target.
- **Factor Reliability:** Target
  - Transfer idempotency holds, and the Ledger double-entry invariant tags Reliability as a secondary factor.
- **Analysis record:** `analysis/root.json`

### Ledger

- **Path:** ledger
- **Area rating:** Target
  - The double-entry invariant is rated target. The reconciliation requirement is not assessed, so the local rating rests on the assessed invariant evidence and is noted as incomplete rather than outstanding.
- **+ Sub-Areas rating:** Target
  - Ledger is a leaf area, so its aggregate equals its local rating. The assessed double-entry invariant reaches target while reconciliation remains not assessed.
- **Factor Correctness:** Target
  - The double-entry invariant holds in the recorded evidence; reconciliation is not assessed, so the factor rating is target but incomplete.
- **Analysis record:** `analysis/ledger.json`
- **Not assessed:** Reconciliation runs daily and flags drift

### Webhooks

- **Path:** webhooks
- **Area rating:** Target
  - The Webhooks target's own signing requirement is rated target.
- **+ Sub-Areas rating:** Minimum
  - The area's own work is rated target, but the Delivery child subtree is rated minimum and pulls the aggregate below the local rating.
- **Factor Security:** Target
  - Every emit path signs webhook payloads and no unsigned path was found.
- **Analysis record:** `analysis/webhooks.json`

### Delivery

- **Path:** webhooks/delivery
- **Area rating:** Minimum
  - The retry requirement is rated target, but the deduplication requirement is rated minimum and sets the local rating at the floor.
- **+ Sub-Areas rating:** Minimum
  - Delivery is a leaf area, so its aggregate equals its local rating. The bounded deduplication gap holds the target at minimum.
- **Factor Reliability:** Minimum
  - Retry-with-backoff reaches target, but redelivery deduplication is bounded and best-effort, so the reliability factor is minimum.
- **Analysis record:** `analysis/webhooks-delivery.json`

## Requirements

### No credentials are committed to the repository

- **State:** active
- **Area:** Sparrow Payments API
- **Rating:** Unacceptable
- **Assessment result record:** `assessments/001-root-no-credentials-are-committed-to-the-repository.json`
- **Rationale:** Rated against this requirement's ratings overrides. Unacceptable applies because a live credential is present in the working tree; a single live secret satisfies that criterion regardless of clean findings elsewhere. The secret value is withheld and referenced only by locator and credential type.

### Every money-moving endpoint enforces authentication

- **State:** active
- **Area:** Sparrow Payments API
- **Rating:** Target
- **Assessment result record:** `assessments/002-root-every-money-moving-endpoint-enforces-authentication.json`
- **Rationale:** Full route coverage with no unauthenticated money-moving exceptions satisfies the requirement. No evidence of step-up or re-authentication controls was recorded, so the result does not reach outstanding.

### Transfers are idempotent on retry

- **State:** active
- **Area:** Sparrow Payments API
- **Rating:** Target
- **Assessment result record:** `assessments/003-root-transfers-are-idempotent-on-retry.json`
- **Rationale:** The transfer entrypoint requires an idempotency key and the replay test shows no duplicate debit on retry. Broader failure-injection evidence was not recorded, so the result remains at target.

### Every transfer debits and credits to a net zero (double-entry invariant)

- **State:** active
- **Area:** Ledger
- **Rating:** Target
- **Assessment result record:** `assessments/004-ledger-every-transfer-debits-and-credits-to-a-net-zero-double-entry-invariant.json`
- **Rationale:** The invariant is enforced before commit and covered by a property test over observed transfer paths. Multi-currency rounding paths were not exercised, so the result remains below outstanding.

### Reconciliation runs daily and flags drift

- **State:** active
- **Area:** Ledger
- **Rating:** not assessed
- **Assessment result record:** `assessments/005-ledger-reconciliation-runs-daily-and-flags-drift.json`
- **Rationale:** Insufficient evidence to rate against the scale, so the requirement is recorded as not assessed rather than assigned a level. The Ledger local rating rests on its assessed requirement and is noted as incomplete.

### Every outbound webhook is signed so merchants can verify its origin

- **State:** active
- **Area:** Webhooks
- **Rating:** Target
- **Assessment result record:** `assessments/006-webhooks-every-outbound-webhook-is-signed-so-merchants-can-verify-its-origin.json`
- **Rationale:** Every recorded emit path signs the payload with origin and replay protection, and tamper rejection is tested. Per-merchant signing-secret rotation was not recorded, so the result remains below outstanding.

### Failed deliveries retry with exponential backoff until acknowledged or the retry window expires

- **State:** active
- **Area:** Delivery
- **Rating:** Target
- **Assessment result record:** `assessments/007-webhooks-delivery-failed-deliveries-retry-with-exponential-backoff-until-acknowledged-or-the-retry-window-expires.json`
- **Rationale:** Failures retry with exponential backoff until acknowledgment or the bounded retry window expires, and exhausted deliveries are recorded rather than silently dropped. Broader failure-injection evidence was not recorded, so the result remains below outstanding.

### A redelivery of an already-acknowledged event is suppressed for that endpoint

- **State:** active
- **Area:** Delivery
- **Rating:** Minimum
- **Assessment result record:** `assessments/008-webhooks-delivery-a-redelivery-of-an-already-acknowledged-event-is-suppressed-for-that-endpoint.json`
- **Rationale:** Deduplication exists and covers the common case but is bounded and best-effort. A duplicate is reachable on a known path, so the requirement falls short of the area's intent while remaining at the acceptable floor.

## Findings

- `assessments/001-root-no-credentials-are-committed-to-the-repository.json` at `internal/gateway/client.go:48`: A live payment-gateway API secret key is committed in plaintext; it matches the format of an active key and is not a placeholder.
- `assessments/001-root-no-credentials-are-committed-to-the-repository.json` at `internal/gateway/client_test.go:12`: A test publishable key is present and is non-secret by design.
- `assessments/002-root-every-money-moving-endpoint-enforces-authentication.json` at `cmd/api/routes.go`: All 18 transfer, refund, and payout routes resolve through the RequireAuth middleware before reaching their handlers; no unauthenticated money-moving route was found.
- `assessments/002-root-every-money-moving-endpoint-enforces-authentication.json` at `internal/gateway/client.go:51`: Evaluator-directed source text resembling prompt injection was present and treated as data, not followed.
- `assessments/003-root-transfers-are-idempotent-on-retry.json` at `internal/transfer/handler.go:73`: POST /transfers requires an Idempotency-Key before processing a transfer.
- `assessments/003-root-transfers-are-idempotent-on-retry.json` at `internal/transfer/handler_test.go:120`: A replay test confirms a retried key returns the original result without a second debit.
- `assessments/004-ledger-every-transfer-debits-and-credits-to-a-net-zero-double-entry-invariant.json` at `ledger/posting.go:64`: The posting routine rejects entries whose sum is not zero before commit.
- `assessments/004-ledger-every-transfer-debits-and-credits-to-a-net-zero-double-entry-invariant.json` at `ledger/posting_test.go:210`: A property test over 10k generated transfers found zero imbalanced postings.
- `assessments/005-ledger-reconciliation-runs-daily-and-flags-drift.json` at `ledger/reconcile.go:31`: A reconcile entrypoint exists, but no reconciliation job output, log, or report was available.
- `assessments/006-webhooks-every-outbound-webhook-is-signed-so-merchants-can-verify-its-origin.json` at `webhooks/emit.go`: All six webhook emit paths sign through the shared signing path before delivery.
- `assessments/006-webhooks-every-outbound-webhook-is-signed-so-merchants-can-verify-its-origin.json` at `webhooks/sign_test.go:48`: A signing test confirms tampered or unsigned payloads are rejected.
- `assessments/007-webhooks-delivery-failed-deliveries-retry-with-exponential-backoff-until-acknowledged-or-the-retry-window-expires.json` at `webhooks/delivery/retry.go:37`: Failed deliveries re-enqueue on an exponential backoff schedule within a bounded 24-hour retry window.
- `assessments/007-webhooks-delivery-failed-deliveries-retry-with-exponential-backoff-until-acknowledged-or-the-retry-window-expires.json` at `webhooks/delivery/retry_test.go:91`: A retry test confirms a transiently failing endpoint is retried and no event is silently lost.
- `assessments/008-webhooks-delivery-a-redelivery-of-an-already-acknowledged-event-is-suppressed-for-that-endpoint.json` at `webhooks/delivery/dedup.go:52`: Each event carries a stable delivery id, and the engine suppresses a redelivery to an endpoint that already acknowledged it.
- `assessments/008-webhooks-delivery-a-redelivery-of-an-already-acknowledged-event-is-suppressed-for-that-endpoint.json` at `webhooks/delivery/dedup.go:52`: Suppression state is retained for only 24 hours, the same as the retry window, so a delivery re-enqueued late in that window can fire after its dedup record has expired.
- `assessments/008-webhooks-delivery-a-redelivery-of-an-already-acknowledged-event-is-suppressed-for-that-endpoint.json` at `docs/merchants/webhooks.md`: Merchant docs say to consume events idempotently, confirming duplicates are possible by design.

## Advice

- [001-rotate-committed-gateway-key](recommendations/001-rotate-committed-gateway-key.md) [active] - The committed-credentials requirement reaches target; no live credential is present in the working tree and the previously exposed key has been revoked.
- [002-produce-reconciliation-evidence](recommendations/002-produce-reconciliation-evidence.md) [active] - The reconciliation requirement becomes assessable and reaches at least the acceptable floor.
- [003-bound-webhook-dedup-window](recommendations/003-bound-webhook-dedup-window.md) [active] - The webhook-delivery deduplication requirement reaches target for the full duration any retry could fire.
