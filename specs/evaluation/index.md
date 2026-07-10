# Evaluation

Durable specifications for the evaluation workflow, protocol, deterministic
runner, structured data, orchestration model, and generated report tree.

# Concepts

- [Evaluation](evaluation.md) - shared invariants and replacement scope.
- [Protocol](protocol.md) - evaluation phases, traversal, and routine ordering.
- [Orchestration](orchestration.md) - runner-owned work graph, scheduling,
  persistence, resume, retry, and cancellation rules.
- [Runner](runner.md) - CLI-owned deterministic evaluation engine, execution
  strategy, run-local logging, and failure taxonomy.
- [Evaluator contract](evaluator-contract.md) - capability, envelope, and
  configuration contract for pluggable evaluators.
- [evaluation.json](evaluation-json.md) - artifact contract for the
  authoritative runner run artifact.

# Subfolders

- [routines/](routines/) - prompt-style routine contracts.
- [records/](records/) - structured JSON payload and `data/` layout contracts.
- [reports/](reports/) - deterministic Markdown report tree contracts.
