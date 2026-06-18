---
type: Functional Specification
title: Duplicate assessment status
description: Detect duplicate assessment records before evaluation reports are rendered.
tags: [evaluation, status, cli]
timestamp: 2026-06-18T00:00:00Z
---

# Duplicate assessment status

The key words **MUST**, **MUST NOT**, **SHOULD**, and **MAY** are to be
interpreted as described in IETF RFC 2119.

## Scope

This change covers duplicate assessment detection during run status and report
renderability checks.

It does not add record replacement, deletion, superseding metadata, or a new
`add-record` mode.

## Requirements

An evaluation run **MUST NOT** be reportable when two or more assessment records
cover the same assessed requirement for the same target path.

The duplicate identity is the ordered `targetPath` plus `requirement`. The
display `target` title is not part of the identity because target titles can be
renamed while target paths remain the stable model location.

`qualitymd evaluation show-status` **MUST** include a `duplicate-assessment` gap
for each duplicate assessment after the first loaded record for that identity.
The gap **MUST** reference the later duplicate record and **SHOULD** identify the
first conflicting assessment in its detail.

`qualitymd evaluation build-report` **MUST** fail through the existing
renderability gate when duplicate assessments are present, and **MUST NOT**
write a partial report for that run.

Analysis records **MAY** continue to replace by target slug. This change does
not alter analysis replacement behavior.
