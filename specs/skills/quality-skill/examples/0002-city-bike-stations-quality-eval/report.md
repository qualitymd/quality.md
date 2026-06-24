# Evaluation Report

## Verdict

- **Root area:** City Bike Station Data
- **Evaluation level:** not recorded
- **Rigor:** standard
- **Evaluation verdict:** Unacceptable
- **Rationale:** The root's provenance gap binds the whole-model rating. Schema & structure is target, but Provenance & freshness is unacceptable for the same missing manifest metadata.

## Scope

Full evaluation of the illustrative city bike station data product.

- **Narrowing:** whole recorded run
- **In scope:** City Bike Station Data; Schema & structure; Provenance & freshness
- **Out of scope:** upstream live telemetry, trip records, pricing data, and the public website that consumes the station register

## Selected Findings and Limitations

- `assessments/002-root-each-snapshot-names-its-source-and-acquisition-time.json` at `metadata/manifest.yaml` [High]: The newest station snapshot lists the output file and row count but leaves the upstream source and acquisition timestamp blank.
- `assessments/001-root-active-stations-are-represented-with-plausible-locations.json` at `data/stations.csv` [Medium]: The station table includes 183 of 184 active stations from the operations roster; station ST-184 is active in the roster and absent from the published table.
- Limitation: The run reviewed a representative current release and two prior snapshots; it did not audit the operator's upstream live system.
- Limitation: Coordinate plausibility was judged against city-boundary and station-address evidence, not by field inspection.
- Limitation: Private maintainer knowledge was not counted as provenance evidence.

## Evidence Basis

- **source:** `data/stations.csv`
- **source:** `ops/active-stations-roster.csv`
- **review:** `manual roster-to-table comparison`
- **source:** `docs/station-addresses.md`
- **review:** `coordinate plausibility review`
- **source:** `metadata/manifest.yaml`
- **source:** `metadata/lineage.md`
- **review:** `manifest provenance re-check before report generation`
- **source:** `schema/stations.schema.yaml`
- **review:** `manual schema-to-table comparison`
- **review:** `identifier uniqueness profile`
- **review:** `prior-snapshot identifier comparison`
- **source:** `docs/releases.md`
- **review:** `snapshot timestamp and release-note review`

## Next Action

- [002-record-snapshot-provenance](recommendations/002-record-snapshot-provenance.md) - The snapshot-provenance requirement reaches target; every current snapshot names its upstream source, acquisition time, and transformation version in the manifest or lineage record.

## Area Breakdown

| Area | Path | Area Rating | Area + Sub-Areas Rating | Factors |
| --- | --- | --- | --- | --- |
| City Bike Station Data | `/` | Unacceptable | Unacceptable | Fitness for use: Minimum; Provenance: Unacceptable |
| Schema & structure | `schema-structure` | Target | Target | Structural validity: Target |
| Provenance & freshness | `provenance-freshness` | Unacceptable | Unacceptable | Provenance: Unacceptable; Freshness: Target |

## Area Details

### City Bike Station Data

- **Path:** /
- **Area rating:** Unacceptable
  - The root's provenance requirement is unacceptable because the newest snapshot cannot be traced to an upstream source or acquisition time. Fitness for use is minimum because one active station is missing from the published table.
- **+ Sub-Areas rating:** Unacceptable
  - The root's provenance gap binds the whole-model rating. Schema & structure is target, but Provenance & freshness is unacceptable for the same missing manifest metadata.
- **Factor Fitness for use:** Minimum
  - The register is usable for most planning and public-information decisions, but the missing active station keeps it below target.
- **Factor Provenance:** Unacceptable
  - The newest snapshot's source and acquisition time are absent, so a reader cannot audit the data without maintainer knowledge.
- **Analysis record:** `analysis/root.json`

### Schema & structure

- **Path:** schema-structure
- **Area rating:** Target
  - The required columns conform to the station schema, no duplicate identifiers are present, and reviewed identifier changes are documented.
- **+ Sub-Areas rating:** Target
  - Schema & structure is a leaf area, so its aggregate equals its local rating.
- **Factor Structural validity:** Target
  - Schema conformance and identifier stability both reach target over the reviewed release and prior snapshots.
- **Analysis record:** `analysis/schema-structure.json`

### Provenance & freshness

- **Path:** provenance-freshness
- **Area rating:** Unacceptable
  - Freshness reaches target, but the missing source and acquisition time for the newest snapshot holds the local area below the acceptable floor.
- **+ Sub-Areas rating:** Unacceptable
  - Provenance & freshness is a leaf area, so its aggregate equals its local rating; the provenance gap binds.
- **Factor Provenance:** Unacceptable
  - Snapshot source and acquisition metadata are missing for the newest published release.
- **Factor Freshness:** Target
  - The newest snapshot is within the documented weekly cadence and no unannounced outage was found.
- **Analysis record:** `analysis/provenance-freshness.json`

## Requirements

### Active stations are represented with plausible locations

- **State:** active
- **Area:** City Bike Station Data
- **Rating:** Minimum
- **Assessment result record:** `assessments/001-root-active-stations-are-represented-with-plausible-locations.json`
- **Rationale:** The current register covers most active stations and all reviewed coordinates are in plausible city bounds, but one active station in the operations roster is absent from the published table.

### Each snapshot names its source and acquisition time

- **State:** active
- **Area:** City Bike Station Data
- **Rating:** Unacceptable
- **Assessment result record:** `assessments/002-root-each-snapshot-names-its-source-and-acquisition-time.json`
- **Rationale:** The current manifest names the generated file but omits upstream source and acquisition time for the newest snapshot.

### Required columns conform to the station schema

- **State:** active
- **Area:** Schema & structure
- **Rating:** Target
- **Assessment result record:** `assessments/003-schema-structure-required-columns-conform-to-the-station-schema.json`
- **Rationale:** All required columns are present, required fields are populated in the reviewed sample, and profiled values match the documented types and ranges.

### Station identifiers are unique and stable across snapshots

- **State:** active
- **Area:** Schema & structure
- **Rating:** Target
- **Assessment result record:** `assessments/004-schema-structure-station-identifiers-are-unique-and-stable-across-snapshots.json`
- **Rationale:** No duplicate station identifiers are present in the current table, and identifier changes across the prior two snapshots are tied to documented station retirements.

### Each snapshot names its source and acquisition time

- **State:** active
- **Area:** Provenance & freshness
- **Rating:** Unacceptable
- **Assessment result record:** `assessments/005-provenance-freshness-each-snapshot-names-its-source-and-acquisition-time.json`
- **Rationale:** The manifest does not name the upstream source or acquisition time for the newest snapshot.

### Published data is refreshed within the documented weekly cadence

- **State:** active
- **Area:** Provenance & freshness
- **Rating:** Target
- **Assessment result record:** `assessments/006-provenance-freshness-published-data-is-refreshed-within-the-documented-weekly-cadence.json`
- **Rationale:** The newest published snapshot is six days old against the documented weekly cadence, and no unannounced outage is present in the release notes.

## Findings

- `assessments/001-root-active-stations-are-represented-with-plausible-locations.json` at `data/stations.csv`: The station table includes 183 of 184 active stations from the operations roster; station ST-184 is active in the roster and absent from the published table.
- `assessments/001-root-active-stations-are-represented-with-plausible-locations.json` at `data/stations.csv`: Reviewed station coordinates fall inside the city boundary and match their documented station addresses within the accepted review tolerance.
- `assessments/002-root-each-snapshot-names-its-source-and-acquisition-time.json` at `metadata/manifest.yaml`: The newest station snapshot lists the output file and row count but leaves the upstream source and acquisition timestamp blank.
- `assessments/003-schema-structure-required-columns-conform-to-the-station-schema.json` at `schema/stations.schema.yaml`: The schema declares station_id, name, status, latitude, longitude, and capacity as required columns, and the station table includes each column.
- `assessments/004-schema-structure-station-identifiers-are-unique-and-stable-across-snapshots.json` at `data/stations.csv`: The current table has no duplicate station_id values.
- `assessments/004-schema-structure-station-identifiers-are-unique-and-stable-across-snapshots.json` at `metadata/lineage.md`: Identifier changes found across the prior two snapshots correspond to documented station retirements and one station relocation.
- `assessments/005-provenance-freshness-each-snapshot-names-its-source-and-acquisition-time.json` at `metadata/manifest.yaml`: The newest station snapshot has no upstream source name or acquisition timestamp in the manifest.
- `assessments/006-provenance-freshness-published-data-is-refreshed-within-the-documented-weekly-cadence.json` at `metadata/manifest.yaml`: The newest station snapshot is six days old, within the documented weekly refresh cadence.

## Advice

- [002-record-snapshot-provenance](recommendations/002-record-snapshot-provenance.md) [active] - The snapshot-provenance requirement reaches target; every current snapshot names its upstream source, acquisition time, and transformation version in the manifest or lineage record.
- [001-backfill-missing-active-station](recommendations/001-backfill-missing-active-station.md) [active] - The active-station coverage requirement reaches target; every active station in the operations roster is present in the station table with plausible coordinates.
