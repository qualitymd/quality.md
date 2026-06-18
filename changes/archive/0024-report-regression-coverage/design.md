---
type: Design Doc
title: Report regression coverage design
description: How report regression tests cover high-risk cases without committed fixtures.
tags: [evaluation, report, tests]
timestamp: 2026-06-18T00:00:00Z
---

# Report regression coverage design

## Context

The [Report regression coverage spec](spec.md) captures behaviors already
proven in the experiment program. The goal is to keep those behaviors from
regressing without adding repository smoke fixtures or large benchmark
snapshots.

## Approach

Use existing `internal/evaluation` temp-repo helpers to create short-lived run
folders during tests. Records are written through `AddRecord`, matching the CLI
record contract, and assertions inspect the generated `Report()` value plus
`report.md` / `report.json` when rendering behavior matters.

The tests should stay small and semantic:

- one adverse/safety report test for secret and prompt-injection style findings;
- one not-assessed/dotted-path test;
- reuse existing structural-root coverage rather than duplicating it.

## Alternatives

**Commit copied seeded reports as fixtures.** Rejected. It would add fixture
churn and overlap with external experiment artifacts.

**Use full subject repositories in tests.** Rejected. Unit tests should not
depend on external checkouts.

## Trade-offs and risks

Temp-run builders do not prove the full evaluator workflow, but they directly
protect deterministic report rendering, which is the behavior under test.
