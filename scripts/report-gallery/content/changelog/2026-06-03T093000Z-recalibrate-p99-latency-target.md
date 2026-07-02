---
date: 2026-06-03
kind: recalibrate
target: api/p99-latency-within-budget
---

Raised the p99 mutation-latency `target` band from 450 ms to 300 ms (and
`outstanding` from 300 ms to 200 ms) after the read-path caching work landed
and four consecutive weeks measured p99 under 280 ms. This is deliberate
recalibration, not drift correction: the service proved the tighter band is
achievable, and the model resets the criterion so the new floor sticks rather
than letting the old band flatter future results.
