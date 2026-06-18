# 002 — README tells the payoff but never shows it

**Target:** readme → approachability
**Related requirements:** *the README shows the format and its payoff by
example*; *the README gets a newcomer to a first result quickly*
**Severity:** medium — two approachability requirements held at minimum.

## Gap

The README shows a realistic QUALITY.md excerpt but never shows what running
`qualitymd` against it produces: no report excerpt, no CLI output, no rendered
result, and no command-plus-exit-code demonstration. The payoff and the "first
result" are described in prose only, so a newcomer cannot see the tool work.

## Evidence

- `README.md:21-103` — realistic input excerpt shown (good).
- `README.md:105-109,140-159` — payoff described in prose, no output shown.
- `README.md:111-132` — install-then-command sequence with no representative
  output.
- `README.md:140-141` — CI exit-code behavior mentioned but not demonstrated.

## Options

1. Add a short "what you get" block showing real output from a built command a
   newcomer can run today — e.g. `qualitymd lint` output and its exit code, and
   a trimmed `qualitymd spec`/`models list` sample.
2. Add an illustrative evaluation-report excerpt (drawn from Appendix A of the
   spec), clearly labelled as produced by the `/quality` skill, to show the
   headline payoff even before the evaluation CLI surface lands.
3. Both: a runnable CLI result now, plus a labelled report excerpt for the
   eventual evaluation payoff.

## Recommended

**Option 3** — show a real, runnable result from a built command (honest today)
*and* a labelled report excerpt for the headline payoff. This satisfies both
"show the payoff by example" and "first result quickly" without overstating what
is built.

## Done criterion

Both requirements reach at least **target**: the README shows produced output
for at least one runnable command and shows what an evaluation produces, with
planned vs. built clearly distinguished. Re-evaluate in a new numbered run.
