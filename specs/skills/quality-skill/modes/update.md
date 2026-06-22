---
type: Functional Specification
title: /quality update
description: Behavioral component spec for orchestrating compatible /quality skill and qualitymd CLI updates.
tags: [skill, quality, mode, update]
timestamp: 2026-06-22T00:00:00Z
---

# /quality update

`update` is the `/quality` skill mode that diagnoses and orchestrates compatible
updates for the separately distributed `/quality` skill and `qualitymd` CLI. It
implements the shared contracts in the parent [/quality skill](../quality-skill.md)
spec and owns only the update-specific behavior below.

The runtime procedure lives at
[`skills/quality/modes/update.md`](../../../../skills/quality/modes/update.md).

This document uses BCP 14 keywords only for testable conformance requirements.
The key words "MUST", "MUST NOT", and "SHOULD" are to be interpreted as
described in [RFC 2119](../../../../docs/reference/rfc2119.md) and
[RFC 8174](../../../../docs/reference/rfc8174.md) when, and only when, they
appear in all capitals.

## Purpose and routing

`update` is selected when the user asks to update, upgrade, repair, or check
compatibility of the `/quality` skill, the `qualitymd` CLI, or their supported
version pair.

The mode's purpose is maintenance orchestration. It is not an Evaluation mode
and **MUST** stop before setup, evaluation, or improvement work.

## Mutation surface and artifacts

`update` may mutate installed tooling only after explicit confirmation. It
**MUST NOT** edit evaluated source, edit `QUALITY.md`, write evaluation records,
or write the quality log.

The skill **MUST NOT** edit installed skill files directly. Skill updates belong
to the Agent Skills installer or package manager. The skill **MUST NOT** replace
the `qualitymd` binary directly. CLI updates belong to `qualitymd update`, owner
package-manager commands, or documented manual guidance.

## Required flow

`update` **MUST** inspect the loaded skill metadata, inspect the visible
`qualitymd` CLI version, use `qualitymd update --check` when available, and
build an update plan before applying changes.

The plan **MUST** classify whether the `/quality` skill, the `qualitymd` CLI,
both, or neither need action. It **MUST** ask for explicit confirmation before
applying any update action.

After a CLI update, `update` **MUST** verify the visible `qualitymd` version
against the loaded skill's required CLI range. After a skill update, it **MUST**
tell the user that the active agent session may still be running previously
loaded skill instructions and may require restart, reload, or a new session.

## Stop conditions

`update` **MUST** stop before applying when it cannot identify the installed
skill/CLI pair, cannot determine a compatible remediation path, or lacks
explicit user confirmation.

If `qualitymd update --check` is unavailable because CLI support is too old,
`update` **MUST** report the install-aware manual or package-manager remediation
path rather than guessing or directly replacing binaries.

## Completion criteria

`update` is complete when it reports the inspected versions, planned or applied
actions, confirmation status, verification result, and any remaining restart,
reload, or manual remediation step.
