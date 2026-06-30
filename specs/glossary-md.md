---
type: Functional Specification
title: glossary.md
description: Artifact contract for the workspace-root glossary used by generated reports and readers.
tags: [glossary, reports]
timestamp: 2026-06-30T00:00:00Z
---

# glossary.md

`glossary.md` is the workspace-root human reference for shared QUALITY.md
terms, concepts, and fixed vocabularies used by generated Evaluation reports.

This document uses BCP 14 keywords only for testable conformance requirements.
The key words "MUST", "MUST NOT", and "SHOULD" are to be interpreted as
described in [RFC 2119](../docs/reference/rfc2119.md) and
[RFC 8174](../docs/reference/rfc8174.md) when, and only when, they appear in all
capitals.

## Background / motivation

Generated reports are easier to scan when they render semantic values directly
and do not repeat local legends in every artifact. A single glossary gives
readers one durable place to look up shared terms, fixed Evaluation vocabulary
values, and model-defined quality rating labels while keeping report tables
compact.

## Requirements

`glossary.md` **MUST** live at the QUALITY.md workspace root, beside the
selected `QUALITY.md`.

`glossary.md` **MUST** define shared QUALITY.md terms, concepts, and fixed
vocabularies in a flat alphabetical list of `##` entries.

Each fixed vocabulary entry **MUST** render a Markdown table with columns
`Label`, `Value`, and `Description`, in that order.

Fixed vocabulary table rows **MUST** preserve the source vocabulary order rather
than sorting values alphabetically.

The glossary introduction **MUST** state that fixed vocabulary `Value` cells are
canonical persisted values and `Label` cells are human display labels.

The glossary introduction **MUST** state that labels, markers, aliases, and case
variants are not accepted as structured data values.

`glossary.md` **MUST** include concept entries for Area, Factor, Finding,
Recommendation, and Requirement.

`glossary.md` **MUST** include a `Quality rating` entry rendered with the fixed
vocabulary table shape. Its rows **MUST** come from this repository's configured
`QUALITY.md` Rating Scale, and the entry **MUST** state that the labels and
values come from that Rating Scale.

`glossary.md` **MUST** include fixed Evaluation enum catalog entries for
Analysis status, Assessment status, Confidence, Data kind, Finding basis,
Finding coverage, Finding rank, Finding severity, Finding type, Rating result,
Rating status, Recommendation impact, Report kind, and Run gap kind.

Generated reports **MAY** link to `glossary.md`, but `glossary.md` **MUST NOT**
be treated as structured Evaluation source data.
