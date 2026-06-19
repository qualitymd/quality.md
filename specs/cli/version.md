---
type: Functional Specification
title: qualitymd version
description: Show structured CLI version metadata without network access.
tags: [cli, command, versioning]
timestamp: 2026-06-19T00:00:00Z
---

# qualitymd version

`qualitymd version` reports detailed local version metadata. It complements the
invocation-wide `qualitymd --version` shortcut, which remains available through
the CLI harness.

This document uses BCP 14 keywords only for testable conformance requirements.
The key words "MUST", "MUST NOT", "SHOULD", and "MAY" are to be interpreted as
described in [RFC 2119](../../docs/reference/rfc2119.md) and
[RFC 8174](../../docs/reference/rfc8174.md) when, and only when, they appear in
all capitals.

`qualitymd version` **MUST NOT** contact the network.

`qualitymd version --json` **MUST** emit a stable JSON object with:

- `schemaVersion`;
- `version`;
- `commit`, when known;
- `developmentBuild`; and
- `specificationVersion`, the version line from the bundled `SPECIFICATION.md`.

Human output **SHOULD** include the same facts in concise text.

`qualitymd --version` **MUST** remain available as the terse human version
shortcut.
