---
date: 2026-06-24
kind: add
target: codebase
---

Added the `codebase` area and expanded the security and auditability coverage
across `api`, `persistence`, and `operations`. The previous model judged the
contract, data, operations, harness, and model artifact, but it had no area for
the implementation that maintainers and coding agents change weekly. That left
maintainability, architecture-boundary consistency, implementation security,
tenant access, break-glass review, and audit-event explainability either
implicit or folded into neighboring areas.

The new shape gives the implementation its own maintainability family
(`analyzability`, `modifiability`, `testability`), local `consistency`, and
local `security`, while the money-touching areas now carry explicit security
and auditability requirements where those concerns are judged directly.
