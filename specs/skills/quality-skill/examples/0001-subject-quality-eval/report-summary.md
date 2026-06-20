# Quality Evaluation Summary

| Field          | Value                       |
| -------------- | --------------------------- |
| Subject        | Sparrow Payments API        |
| Run            | `0001-subject-quality-eval` |
| Scope          | Full evaluation             |
| Rigor          | Standard                    |
| Overall rating | Unacceptable                |
| Full report    | [report.md](report.md)      |
| Machine report | [report.json](report.json)  |

## Summary

Held at unacceptable by a single binding constraint: a live payment-gateway
credential is committed to the repository. Rotating and removing it would lift
the root off the floor, but the Webhooks delivery deduplication gap would then
bind at minimum.

| Target               | Local rating | Overall rating | Driver                                                                                                                                                                                                                                                     |
| -------------------- | ------------ | -------------- | ---------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------- |
| Sparrow Payments API | Unacceptable | Unacceptable   | Held at unacceptable by a single binding constraint: a live payment-gateway credential is committed to the repository. Rotating and removing it would lift the root off the floor, but the Webhooks delivery deduplication gap would then bind at minimum. |
| Ledger               | Target       | Target         | The target requirements are satisfied by the available evidence.                                                                                                                                                                                           |
| Webhooks             | Target       | Minimum        | The target meets the floor but has a known bounded deduplication gap.                                                                                                                                                                                      |
| Delivery             | Minimum      | Minimum        | The target meets the floor but has a known bounded deduplication gap.                                                                                                                                                                                      |

## Top Issues

1. **critical**
   A live gateway secret key is committed in plaintext; value withheld.
   `internal/gateway/client.go:48`
   Assessment: `assessment-results/001-root-no-committed-credentials.json`
2. **missing-evidence**
   A reconcile entrypoint exists, but no scheduled run output, log, or report was
   available.
   `ledger/reconcile.go:31`
   Assessment: `assessment-results/005-ledger-reconciliation.json`
3. **medium**
   Deduplication state is retained for the same 24-hour window as retries, so
   late re-enqueue can duplicate delivery.
   `webhooks/delivery/dedup.go:52`
   Assessment: `assessment-results/008-delivery-redelivery-dedup.json`

## Recommendations

Primary next action: use `001-rotate-committed-gateway-key`.

| Recommendation ID                     | Priority | Recommendation                                                                            | Done criterion                                                                                                                                           |
| ------------------------------------- | -------: | ----------------------------------------------------------------------------------------- | -------------------------------------------------------------------------------------------------------------------------------------------------------- |
| `001-rotate-committed-gateway-key`    |        1 | [Rotate Committed Gateway Key](recommendations/001-rotate-committed-gateway-key.md)       | The committed-credentials requirement reaches target: no live credential is present in the working tree and the previously exposed key has been revoked. |
| `002-produce-reconciliation-evidence` |        2 | [Produce Reconciliation Evidence](recommendations/002-produce-reconciliation-evidence.md) | The reconciliation requirement becomes assessable and reaches at least the acceptable floor.                                                             |
| `003-bound-webhook-dedup-window`      |        3 | [Bound Webhook Dedup Window](recommendations/003-bound-webhook-dedup-window.md)           | The webhook-delivery deduplication requirement reaches target within the declared retry window.                                                          |

## Scope & Limitations

Scope: **Full evaluation**

In scope: Sparrow Payments API; Ledger; Webhooks; Delivery

- Standard rigor did not seek broader failure-injection and step-up-control
  evidence required for outstanding ratings.
- The committed-credential finding reflects the tracked working tree at the
  evaluated commit.
- Transfer idempotency and webhook deduplication each rest on one
  replay/redelivery test, not sustained concurrency or fault-injection soak.
