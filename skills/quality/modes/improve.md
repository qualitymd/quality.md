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
2. Before editing the subject or `QUALITY.md`, ask for explicit confirmation of
   the recommendation and option to apply.
3. Apply only the confirmed option.
4. Run a new evaluation in a new numbered folder and link it back to the prior
   run.
5. Check the done criterion against the new folder's rating.

If the user does not confirm a recommendation and option, stop after reporting
the evaluation and recommendations.
