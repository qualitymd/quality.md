---
date: 2026-06-17
kind: drift-correction
target: persistence/migration-rollback
---

Tightened `migration-rollback` from "rollback steps are documented" to
"rollback paths are rehearsed against the current schema". Two release
incidents showed documented steps failing against migrations added after the
last rehearsal, so the old criterion had silently fallen out of step with what
recoverability requires. This is a drift correction: the model now demands
rehearsal evidence newer than the latest schema change, which the runbook
alone never guaranteed.
