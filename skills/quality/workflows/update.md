---
type: Runtime Workflow
title: Update Workflow
description: Runtime workflow for updating the /quality skill and qualitymd CLI pair.
---

# Update Mode

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

1. Read the loaded `SKILL.md` frontmatter. Record `metadata.version` and
   `metadata.requires-qualitymd-cli` when present. If either is missing, report
   that skill release metadata is unavailable and continue with the best
   available CLI facts.
2. Inspect the CLI:
   - Prefer `qualitymd version --json`.
   - Fall back to `qualitymd --version` when structured output is unavailable.
   - If the CLI is missing, classify the CLI action as install/repair.
3. Run `qualitymd update --check` when available. Use its recommended action
   and install-method facts for the CLI portion of the plan. If unavailable,
   fall back to documented install or package-manager guidance.
4. Check skill update support through the Agent Skills installer only when a
   supported check/update command is available in the current environment. Do
   not guess at private installer state. If unsupported, report that skill
   update automation is unavailable and show the manual reinstall command:

   ```sh
   npx skills add qualitymd/quality.md
   ```

5. Present a concise update plan before mutation. Make the primary call to
   action visually clear and use bold labels for the plan fields when the
   surface supports Markdown. Include:
   - current skill version and required CLI range when known;
   - current CLI version and whether it is in range;
   - whether the plan acts on the skill, the CLI, both, or neither;
   - exact owner command(s) that would be run, or manual commands when apply is
     unsupported;
   - restart/reload expectation if the skill changes.
6. Ask for explicit confirmation before applying any update action using a
   decision brief:

   ```text
   **Apply update plan?**

   **Changes:** <skill | CLI | both>
   **Evidence/reason:**
   **Recommended option:**
   **Alternatives:**
   **Done criterion / verification:**
   ```

   The brief must name which owner command performs each mutation and how
   compatibility will be checked afterward.
7. If confirmed, delegate only to owner commands:
   - CLI: `qualitymd update` when available and applicable, or the
     package-manager/install command recommended by `qualitymd update --check`.
   - Skill: the Agent Skills installer/package-manager command when available.
   - Otherwise stop with the manual command.
8. After a CLI action, re-run the CLI version check and verify the visible CLI
   satisfies the area `metadata.requires-qualitymd-cli` range when known.
9. After a skill action, tell the user the current agent session may still be
   using previously loaded skill instructions and may need restart, reload, or a
   new session before the updated skill is active.
10. Report the update closeout status-first: inspected versions, planned or
    applied actions, confirmation status, verification result, remaining restart
    or manual remediation step, and the recommended next action.

`update` does not create or edit `QUALITY.md`, create evaluation records, build
reports, rate the evaluated source, or apply quality recommendations.
