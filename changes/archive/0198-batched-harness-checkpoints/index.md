# Batched harness checkpoints

Children of the
[Batched harness checkpoints](../0198-batched-harness-checkpoints.md)
Change Case.

# Concepts

- [Functional spec](spec.md) - requirements for the rolling-window harness
  checkpoint request/result contract: keeping a bounded set of dependency-ready
  requests outstanding, per-member correlation and validation, partial-submission
  and resume behavior, and semantic parity with the sequential path.
- [Design doc](design.md) - how the runner delivers the window: pending state as
  a set, reuse of the concurrent scheduler's ready-frontier computation to top up
  and stream over discrete resume calls, streamed result intake,
  not-submitted-vs-failed accounting, and the skill loop.
