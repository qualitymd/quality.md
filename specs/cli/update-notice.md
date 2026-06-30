---
type: Functional Specification
title: qualitymd update notice
description: Cross-command ambient update notice contract for qualitymd commands.
tags: [cli, update, notice]
timestamp: 2026-06-22T00:00:00Z
---

# qualitymd update notice

This spec owns the ambient update notice that ordinary `qualitymd` commands may
surface from a local cache. The explicit update command and `--check` behavior
live in [`qualitymd update`](update.md).

This document uses BCP 14 keywords only for testable conformance requirements.
The key words "MUST", "MUST NOT", "SHOULD", and "MAY" are to be interpreted as
described in [RFC 2119](../../docs/reference/rfc2119.md) and
[RFC 8174](../../docs/reference/rfc8174.md) when, and only when, they appear in
all capitals.

## Ambient notice

Ordinary `qualitymd` commands **MAY** surface a one-line update-available notice
when a local cache indicates that a newer, ready release exists. The notice
**MUST** be written to stderr and **MUST NOT** appear in stdout or in any
machine-readable output. It **MUST** be suppressed when stderr is not an
interactive terminal, when CI is detected, when the documented opt-out is set,
and for development builds. It **MUST** name the current and latest versions and
the exact command to run, and **SHOULD** include the release-notes reference when
known.

The notice **MUST** be served from a local cache. Ordinary commands **MUST NOT**
block on a network fetch to produce it, and their exit code and primary output
**MUST NOT** be affected by the notice's presence, absence, or any failure to
produce it. The cache **MAY** be refreshed by a bounded, best-effort check no
more frequently than a documented interval; a failed refresh **MUST** be silent.
A refreshed value **MAY** surface on a later invocation rather than the one that
triggered the refresh.

`QUALITYMD_NO_UPDATE_CHECK=1` disables explicit update checks, the ambient cache
refresh, and the ambient notice.
