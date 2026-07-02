---
date: 2026-05-06
kind: add
target: service-contract
---

Added a `service-contract` area so the contract is rated on its own axis
instead of only serving as the reference other areas are judged against.
Contract drift kept surfacing during API review: handlers and the contract
disagreed, and the model had no place to land that finding. An earlier
quality-loop cycle surfaced the pattern; its evaluation records predate the
runs retained in this gallery. The area carries `completeness` and
`consistency`, and `api` requirements now reference the contract by its source
selector so the dependency is traceable.
