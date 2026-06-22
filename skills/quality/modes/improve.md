# Improve Mode

Use improve to evaluate, recommend, apply an explicitly approved change, and
verify the result.

## Decision Tree

```text
Run evaluation first
- recommendations produced?
  - no? report no apply step
  - yes? ask for recommendation and option confirmation

User confirms specific option?
- no? stop after reporting evaluation and recommendations
- yes? apply only that option

After applying confirmed option
- create a new run and re-evaluate affected scope
```

## Procedure

1. Read and follow [`evaluate.md`](evaluate.md) for the evaluation and
   recommendation pass.
2. Before asking for confirmation, make the recommendation self-contained enough
   to apply and later re-evaluate: current state, proposed option, affected
   scope, risk, verification path, and re-evaluation criterion.
3. Before editing the root area or `QUALITY.md`, ask for explicit confirmation of
   the recommendation and option to apply using a decision brief:

   ```text
   Decision: apply <recommendation>?
   - Changes: <evaluated source | QUALITY.md | both>
   - Evidence/reason:
   - Recommended option:
   - Alternatives:
   1. Apply recommended option
   2. Defer and keep recommendation open
   3. Skip and record why
   - Done criterion / verification:
   ```

   If options differ in risk or coverage, say so explicitly. Do not treat an
   obvious recommendation as consent.
4. Apply only the confirmed option.
   - When the applied change is a **model change** to `QUALITY.md` (per the
     meaningful-change taxonomy in [`../guides/authoring.md`](../guides/authoring.md)),
     append one quality log entry under `quality/log/` cross-linking the
     evaluation run and recommendation it came from. This needs no confirmation
     beyond the user's existing confirmation of the change; its rationale is the
     rationale already in the decision brief. An evaluated-source fix that does
     not change the model gets no log entry. See the contract in
     [`../SKILL.md`](../SKILL.md).
5. Run a new evaluation in a new numbered folder and link it back to the prior
   run.
6. Check the done criterion against the new folder's rating.
7. Report the improvement delta:

   ```text
   Improvement result
   - Recommendation:
   - Applied option:
   - Changed artifacts:        (name the quality log entry when the model changed)
   - Before evidence:
   - After evidence:
   - Verification:
   - Rating movement:
   - Remaining gaps / limits:
   ```

   If the rating did not move, say why when knowable. If verification is
   incomplete, label the result as limited rather than fully confirmed.

If the user does not confirm a recommendation and option, stop after reporting
the evaluation and recommendations.
