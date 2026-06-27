---
type: Runtime Workflow
title: Update Workflow
description: Runtime workflow for updating the /quality skill and qualitymd CLI pair.
---

# Update Workflow

Use update to keep the installed `/quality` skill and `qualitymd` CLI
compatible. Update is maintenance orchestration: diagnose the pair, plan the
actions, ask before mutation, delegate mechanics, and verify afterward.

## Decision Tree

```text
Collect current state
- skill metadata present? capture version and requires-qualitymd-cli
- CLI present? capture version facts
- CLI update check available? capture recommended action
- skill installer can check/update? capture latest skill action

Build plan
- neither needs action? report compatible/current enough
- CLI only? plan CLI owner-channel action
- skill only? plan skill-installer action
- both? plan ordered skill/CLI actions and verification

Apply only after confirmation
- confirmed? delegate to owner commands
- not confirmed? stop after plan

Verify
- CLI visible and in area range?
- skill changed? tell user restart/reload/new session may be required
```

## Procedure

1. Emit the run frame before tool inspection:

   ```text
   **QUALITY.md · update**
   - **Model file:** n/a
   - **Scope:** installed `/quality` skill and `qualitymd` CLI compatibility
   - **Mutation:** installed tooling only, after confirmation
   - **Artifacts:** none in the current project by default
   - **Next gate:** version inspection, update plan, confirmation, verification
   ```

2. Read the loaded `SKILL.md` frontmatter. Record `metadata.version` and
   `metadata.requires-qualitymd-cli` when present. If either is missing, report
   that skill release metadata is unavailable and continue with the best
   available CLI facts.
3. Inspect the CLI:
   - Prefer `qualitymd version --json`.
   - Fall back to `qualitymd --version` when structured output is unavailable.
   - If the CLI is missing, classify the CLI action as install/repair.
4. Run `qualitymd update --check` when available. Use its recommended action
   and install-method facts for the CLI portion of the plan. If unavailable,
   fall back to documented install or package-manager guidance.
5. Check skill update support through the Agent Skills installer only when a
   supported check/update command is available in the current environment. Do
   not guess at private installer state. If unsupported, report that skill
   update automation is unavailable and show the manual reinstall command:

   ```sh
   npx skills add qualitymd/quality.md
   ```

6. Present a concise progress/status block after version inspection and before
   the mutation gate. Keep it factual and user-facing:

   ```text
   **Update status**

   **Skill:** <current version and required CLI range, or unavailable>
   **CLI:** <current version and in-range status, or missing>
   **Plan needed:** <none | skill | CLI | both | manual remediation>
   **Next:** review the update plan before anything changes.
   ```

7. Present a concise update plan before mutation. Make the primary call to
   action visually clear and use bold labels for the plan fields when the
   surface supports Markdown. Include:
   - current skill version and required CLI range when known;
   - current CLI version and whether it is in range;
   - whether the plan acts on the skill, the CLI, both, or neither;
   - exact owner command(s) that would be run, or manual commands when apply is
     unsupported;
   - restart/reload expectation if the skill changes.
8. Ask for explicit confirmation before applying any update action using a
   decision brief. Where the harness will itself prompt to authorize the owner
   command, render the confirmation through that native gate and keep the plan's
   teaching in the preceding message rather than stacking a second text `y`/`n`
   gate for the same mutation; never drop the confirmation, only the redundant
   gate. The text-fallback brief:

   ```text
   **Apply update plan?**
   <skill | CLI | both> — <what version moves where>

     [y] Run <owner command(s)> now  — recommended
     [n] Skip — leave versions as they are

   Reason: <why the update is needed, e.g. required CLI range>
   Done when: <how compatibility is re-checked afterward>
   ```

   The choices must name which owner command performs each mutation; the
   `Done when` line states how compatibility is checked afterward.
9. If confirmed, delegate only to owner commands:
   - CLI: `qualitymd update` when available and applicable, or the
     package-manager/install command recommended by `qualitymd update --check`.
   - Skill: the Agent Skills installer/package-manager command when available.
   - Otherwise stop with the manual command.
10. After a CLI action, re-run the CLI version check and verify the visible CLI
    satisfies the area `metadata.requires-qualitymd-cli` range when known.
11. After a skill action, tell the user the current agent session may still be
    using previously loaded skill instructions and may need restart, reload, or a
    new session before the updated skill is active.
12. Report the update closeout status-first: inspected versions, planned or
    applied actions, confirmation status, verification result, remaining restart
    or manual remediation step, and the recommended next action. Use this shape:

    ```text
    **Update complete**

    **Inspected:** <skill version/range and CLI version>
    **Applied:** <skill action | CLI action | none>
    **Verification:** <compatibility check result>
    **Restart/reload:** <needed | not needed | unknown>
    **Not changed:** no QUALITY.md edits, no evaluation records, no reports, no quality changelog
    **Next:** <recommended next action>
    ```

`update` does not create or edit `QUALITY.md`, create evaluation records, build
reports, rate the evaluated source, or apply quality recommendations.
