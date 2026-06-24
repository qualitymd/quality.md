---
schemaVersion: 1
title: Record source and acquisition metadata for every snapshot
gap: The newest station snapshot omits upstream source and acquisition time.
evidenceLocators:
  - metadata/manifest.yaml
  - metadata/lineage.md
assessmentResultRecords:
  - assessments/002-root-each-snapshot-names-its-source-and-acquisition-time.json
  - assessments/005-provenance-freshness-each-snapshot-names-its-source-and-acquisition-time.json
remediationOptions:
  - Backfill source-system name, acquisition time, and transformation version for the newest snapshot, then make those fields required in the release checklist.
  - Add the metadata only to the lineage note.
  - Rely on maintainer memory for the missing source details.
recommendedOption: Backfill source-system name, acquisition time, and transformation version for the newest snapshot, then make those fields required in the release checklist.
doneCriterion: The snapshot-provenance requirement reaches target; every current snapshot names its upstream source, acquisition time, and transformation version in the manifest or lineage record.
---

# Record source and acquisition metadata for every snapshot

**Area / factor:** City Bike Station Data (root) -> Provenance; Provenance &
freshness -> Provenance
**In-scope requirement:** *Each snapshot names its source and acquisition time*
**Current rating:** Unacceptable - **binding constraint** on the whole-model
rating.

## Gap

The newest station snapshot lists its output file and row count, but it does not
name the upstream source or acquisition time. A reader cannot audit whether the
data is stale, trace an error to the source system, or reproduce the release
without asking the maintainer.

**Evidence**

- `metadata/manifest.yaml` - newest snapshot lacks source-system and acquisition
  timestamp fields.
- `metadata/lineage.md` - lineage note does not supply the missing source
  metadata.
- The manifest and lineage notes were re-checked before report generation
  because this gap binds the headline rating.

## Options

- **(a) Backfill the source-system name, acquisition time, and transformation
  version for the newest snapshot, then require those fields in the release
  checklist.**
- **(b) Add the metadata only to the lineage note.** Better than nothing, but it
  leaves the machine-readable manifest incomplete.
- **(c) Rely on maintainer memory.** This does not create auditable provenance.

## Recommended

**(a) Backfill the manifest and make the metadata required in the release
checklist.** The manifest is the durable handoff artifact; the release checklist
prevents the same gap from recurring.

## Done-criterion

The requirement *Each snapshot names its source and acquisition time* reaches
**Target**: every current snapshot names its upstream source, acquisition time,
and transformation version in the manifest or lineage record, and the next
release checklist requires those fields before publication.
