---
type: Design Doc
title: Model constituents by default — design
description: How the /quality skill's constituent-coverage guidance is reframed to model-by-default — the disqualifier set, the no-silent-omission rule, the (a)-vs-(b) routing, and the enforcement bar.
tags: [skill, authoring, setup, constituents, coverage]
timestamp: 2026-06-24T00:00:00Z
---

# Model constituents by default — design

## Context

Answers the [functional spec](spec.md) for the
[0080 change case](../0080-model-constituents-by-default.md). The work is a
guidance reframe across the bundled skill and its spec mirror — no code. The
design question is not *what* to build but *how to word the inversion* so it
reverses the bias without losing the legitimate guard it replaces.

## The problem with the old shape

The old guidance is structurally lopsided. "Cover the domain's constituent kinds"
lists four accounting options as peers — *model it / defer in Scope / out of Scope
/ unknown* — then spends its emphasis ("a prompt, not a quota," "earn each area,"
"not a roster," "a throwaway script earns almost none") guarding the *upper* bound
on areas. There is no matching guard on the *lower* bound. The inclusion test is
phrased as a gate the constituent must pass (owned artifact AND divergent factors
AND traced to Need/Risk), so the natural reading under uncertainty is "don't add
the area." And because deferral is a peer option, writing "deferred — next
iteration" *feels* like discharging the obligation while actually dropping
coverage.

The result observed in the field: a composite monorepo modeled as a flat root,
every real constituent deferred, and a clean pass through the maturity bar.

## Approach

Three coordinated moves, applied identically in the runtime guide and the spec
mirror so they cannot drift.

### 1. Invert the default and collapse the gate into two skip conditions

Replace "earn each area" with "model each constituent as its own area by default."
The old three-part gate becomes a short **disqualifier list** — the same tests,
read as reasons *not* to spin out an area:

- old "divergent factor family" → **No distinct concerns** (fold).
- old "in this entity's boundary" → **Not germane / outside the boundary** (out of
  scope).
- old "traced to a Need or Risk" → **dropped as a gate.** It governs which
  *factors* a modeled area carries, not whether the area exists; keeping it as an
  inclusion gate is what let "I can't immediately name a risk" suppress a real
  constituent. Two disqualifiers is the "more simple" target.

The "prompt, not a roster" caution is *kept* but **re-aimed**: it no longer
opposes modeling, it scopes the enumeration to *this* entity (don't import generic
best-practice areas a throwaway entity would never carry). Model-by-default and
not-a-roster stop being opposing pressures and become one rule: enumerate what the
entity's domain implies, model each unless a disqualifier fires.

### 2. Make omission impossible for a germane concern

The sharp correction. "No artifact yet" is removed from the skip conditions
entirely, because for a *germane* concern the absent artifact is the highest-value
early finding — exactly what the model should surface, not hide. The rule: a
germane concern is **always** surfaced as a *ratable* element (something that
produces a rating), never as prose. A Scope/"deferred" note is explicitly *not*
surfacing.

Two surfacing forms, routed by the *same* distinct-concerns test that decides
area-vs-fold for present constituents:

- **(a) minimal area + missing-anchor finding** — when the kind would carry its
  own factor family once it exists (high-leverage/first-class, or a
  whole-constituent gap). The empty area rates poorly on its own axis and reserves
  the shape.
- **(b) requirement on an existing area** — when the concern folds into an existing
  area's factors, the gap is partial/degree-based, or a standalone area would be a
  single-finding stub. The gap rates under the parent.

This reuses the existing "carry a germane, high-leverage kind as an area even when
thin or missing" rule (already in the guide) and generalizes it: that rule was the
(a) branch; (b) is its complement for concerns that don't warrant their own axis.

### 3. Make the default enforceable, not advisory

Inverting the prose is not enough — the maturity check still passed deferrals. So:

- **Setup Verify-and-Close / completion criteria:** a non-disqualified constituent
  left unmodeled or merely deferred is an important gap that bars
  `evaluation-ready`.
- **Top 10 check:** flip from "treat constituent presence as earned; do not flag a
  missing one" to "flag a non-disqualified constituent that is unmodeled or only
  deferred; a bare deferral note does not satisfy the check." The existing
  not-germane carve-out (don't flag what a throwaway entity wouldn't carry) is
  preserved — it is now just the not-germane disqualifier.

Deferral survives only as a narrow, blocker-recorded exception, so the report
categories that mention "deferred areas" (`evaluate.md`, `reporting.md`) stay
valid without change — they now describe the rare blocked case, not a routine one.

## Alternatives

- **Keep deferral as a peer but add a counterweight sentence.** Rejected: the
  field failure came precisely from deferral reading as cost-free. A sentence
  against a structural peer-option loses; the option itself has to be demoted.
- **Three disqualifiers (keep "no distinct artifact yet").** Rejected on the
  user's point: for a germane concern that is the worst case, not a skip reason.
  Folding it into the no-silent-omission rule is both more correct and simpler.
- **Drop "no distinct concerns" too (model even overlapping constituents).**
  Rejected: it reintroduces double-counting — two areas rating the same factors —
  which the three-projections rule already warns against. The distinct-concerns
  test does real work (it also routes (a) vs (b)), so it stays.
- **Push the principle into `SPECIFICATION.md`.** Rejected: constituent coverage is
  authoring *judgment*, not format *semantics*. `SPECIFICATION.md` carries none of
  the earn-it/defer language today; keeping it out preserves the layering (format =
  what a model is; skill guide = how to author one well).
- **Remove in-scope deferral entirely.** Rejected: a genuinely blocked constituent
  (undecided boundary) needs an honest disposition that isn't a forced wrong area.
  Narrowing it to blocker-recorded keeps the escape hatch without the loophole.

## Trade-offs & risks

- **More areas, including empty ones.** First models get larger and some areas
  carry only a missing-anchor finding. That is the intended cost — an empty area
  with a true finding is a real signal, and the not-germane disqualifier plus the
  "prompt, not a roster" scoping keep it from ballooning into generic
  best-practice areas.
- **Routing judgment ((a) vs (b)) is a soft call.** Mitigated by tying it to the
  distinct-concerns test the author already applies elsewhere, rather than a new
  axis of judgment.
- **"Blocked" could become the new "next iteration."** Mitigated by requiring the
  specific blocker to be recorded and by the maturity bar treating an unrecorded
  deferral as a gap.

## Open questions

- Whether the dogfooded `QUALITY.md` in this repo should gain or split areas under
  the new rule is left to the follow-up pass named in the change case, so this case
  stays a pure guidance change.
