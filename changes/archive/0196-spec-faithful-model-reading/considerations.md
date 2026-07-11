---
type: Sketch
title: Spec-faithful model reading — considerations
description: Forward-looking direction for treating source as a typed, resolver-dispatched selector; shapes but does not bind this case.
---

# Spec-faithful model reading — considerations

Non-binding. This sketch records the direction that the source findings point
toward, so a later case can pick it up with the reasoning intact. It is not a
requirement of [0196](../0196-spec-faithful-model-reading.md).

## Source is a selector, not a path

`SPECIFICATION.md` defines source as "a selector describing the material
evaluated by an area." A selector could be a path, a glob, a saved query, or a
prose description of a body of evidence ("all specs", "the deployed API", "open
tickets in the support queue"). The runner today implements exactly one selector
kind — a workspace-contained filesystem path. This case makes the path/glob kind
spec-faithful and stops silent empty evidence; it does not add other kinds.

## Separate resolution from judgment

The determinism, source-as-data trust boundary, and auditability the runner
guarantees are properties of the **resolved evidence bundle**, not of **how it
was gathered**:

```
source selector ──[ resolver (per kind) ]──► bounded, hashed evidence bundle ──[ evaluator ]──► judgment
```

- For a `path`/`glob` selector, the resolver is the deterministic walk. Gathering
  and record are both reproducible.
- For a prose or live-system selector, the resolver needs tools (a query, an
  agent). Gathering is not deterministic, but the result is still **captured**
  into the same bounded, hashed bundle — so re-judgment, diffing, and audit still
  hold (reproducibility _of record_), and the data-not-instructions boundary is
  preserved. An agent-gathered "all specs" is legitimate _only_ when it flows
  back through the bundle contract, never as free repo exploration by the judge.

## Shape to decide later

- Model `source` as a typed selector (`{kind, selector}`) vs. a bare string with
  the resolver inferring kind (path vs. glob vs. prose). Bare string is friendlier
  to hand-authoring; typed is unambiguous.
- The evaluator contract's `Subagents` capability is the natural seam for an
  SDK-subagent / meta-harness resolver — as a **resolver**, feeding the bundle,
  distinct from the evaluator that judges against it. This keeps the fragile,
  hand-maintained evaluator adapters (the `codex` 403 / `claude` schema-mismatch
  failures) separable from the resolution problem.
- Whether the format spec should commit to non-filesystem selectors at all, or
  stay filesystem-only-but-honest while the runner rejects unresolvable selectors
  loudly (which R4 in this case begins).
