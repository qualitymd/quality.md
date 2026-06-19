# Quality Evaluation Summary

**Run:** `0001-subject-quality-eval`
**Subject:** `Sparrow Payments API`
**Scope:** whole model
**Effort:** standard
**Root rating:** Unacceptable
**Full report:** [report.md](report.md)
**Machine report:** [report.json](report.json)

## Headline

Held at Unacceptable by a committed live payment-gateway credential. Removing
and rotating the credential lifts the root only to Minimum until the webhook
delivery deduplication gap is also closed.

## Top Risks

1. **critical** - A live gateway credential is committed in plaintext; value
   withheld and referenced by credential type and locator only.
2. **high** - Webhook delivery redelivery behavior is not bounded by durable
   deduplication evidence.
3. **missing-evidence** - Ledger reconciliation cannot be fully assessed because
   scheduled reconciliation evidence is absent.

## Rating Summary

| Target               | Aggregate rating | Reason                                                                             |
| -------------------- | ---------------- | ---------------------------------------------------------------------------------- |
| Sparrow Payments API | Unacceptable     | Root security finding binds the whole-model rating.                                |
| Ledger               | Target           | Ledger requirements are at target except the not-assessed reconciliation evidence. |
| Webhooks             | Minimum          | Delivery child target holds the Webhooks aggregate below local Target.             |
| Delivery             | Minimum          | Redelivery deduplication needs a bounded durable guarantee.                        |

## Limitations

- One Ledger requirement is not assessed because scheduled reconciliation
  evidence was not available.
- This is a standard-effort evaluation over representative evidence, not a full
  source audit.

## Next Action

Rotate the exposed gateway credential and remove it from version control.

See active recommendations:

- [001-rotate-committed-gateway-key](recommendations/001-rotate-committed-gateway-key.md) - credential rotation is complete and a clean secret scan passes.
- [002-produce-reconciliation-evidence](recommendations/002-produce-reconciliation-evidence.md) - reconciliation evidence exists and the requirement can be rated.
- [003-bound-webhook-dedup-window](recommendations/003-bound-webhook-dedup-window.md) - webhook redelivery deduplication is bounded and verified.
