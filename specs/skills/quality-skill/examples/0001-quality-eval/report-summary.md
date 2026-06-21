# Quality Evaluation Summary

| Field | Value |
| --- | --- |
| Root area | Sparrow Payments API |
| Run | `0001-quality-eval` |
| Scope | Full evaluation |
| Rigor | not recorded |
| Evaluation verdict | Unacceptable |
| Full report | [report.md](report.md) |
| Machine report | [report.json](report.json) |

## Verdict

The root's own local rating binds. Ledger is rated target and Webhooks is rated minimum; if the local secret gap were cleared, Webhooks would become the next binding constraint.


| Area | Local rating | Aggregate rating | Rating basis |
| --- | --- | --- | --- |
| Sparrow Payments API | Unacceptable | Unacceptable | The root's own local rating binds. Ledger is rated target and Webhooks is rated minimum; if the local secret gap were cleared, Webhooks would become the next binding constraint. |
| Ledger | Target | Target | Ledger is a leaf area, so its aggregate equals its local rating. The assessed double-entry invariant reaches target while reconciliation remains not assessed. |
| Webhooks | Target | Minimum | The area's own work is rated target, but the Delivery child subtree is rated minimum and pulls the aggregate below the local rating. |
| Delivery | Minimum | Minimum | Delivery is a leaf area, so its aggregate equals its local rating. The bounded deduplication gap holds the target at minimum. |

## Selected Findings

1. **Critical**  
   A live payment-gateway API secret key is committed in plaintext; it matches the format of an active key and is not a placeholder.
   `internal/gateway/client.go:48`
   Assessment: `assessments/001-root-no-credentials-are-committed-to-the-repository.json`
2. **Medium**  
   A reconcile entrypoint exists, but no reconciliation job output, log, or report was available.
   `ledger/reconcile.go:31`
   Assessment: `assessments/005-ledger-reconciliation-runs-daily-and-flags-drift.json`
3. **Medium**  
   Suppression state is retained for only 24 hours, the same as the retry window, so a delivery re-enqueued late in that window can fire after its dedup record has expired.
   `webhooks/delivery/dedup.go:52`
   Assessment: `assessments/008-webhooks-delivery-a-redelivery-of-an-already-acknowledged-event-is-suppressed-for-that-endpoint.json`

## Recommended Actions

Primary next action: use `001-rotate-committed-gateway-key`.

| Recommendation ID | Priority | Recommendation | Done criterion |
| --- | --- | --- | --- |
| `001-rotate-committed-gateway-key` | 1 | [Rotate Committed Gateway Key](recommendations/001-rotate-committed-gateway-key.md) | The committed-credentials requirement reaches target; no live credential is present in the working tree and the previously exposed key has been revoked. |
| `002-produce-reconciliation-evidence` | 2 | [Produce Reconciliation Evidence](recommendations/002-produce-reconciliation-evidence.md) | The reconciliation requirement becomes assessable and reaches at least the acceptable floor. |
| `003-bound-webhook-dedup-window` | 3 | [Bound Webhook Dedup Window](recommendations/003-bound-webhook-dedup-window.md) | The webhook-delivery deduplication requirement reaches target for the full duration any retry could fire. |

## Scope & Limitations

Scope: **Full evaluation**

In scope: Sparrow Payments API; Ledger; Webhooks; Delivery

- Reconciliation runs daily and flags drift: Insufficient evidence to rate against the scale, so the requirement is recorded as not assessed rather than assigned a level
- Insufficient evidence to rate against the scale, so the requirement is recorded as not assessed rather than assigned a level
- The reconciliation requirement is not assessed, so the local rating rests on the assessed invariant evidence and is noted as incomplete rather than outstanding
- The assessed double-entry invariant reaches target while reconciliation remains not assessed
- The double-entry invariant holds in the recorded evidence; reconciliation is not assessed, so the factor rating is target but incomplete
