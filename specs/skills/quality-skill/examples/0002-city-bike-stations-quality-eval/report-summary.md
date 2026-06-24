# Quality Evaluation Summary

| Field | Value |
| --- | --- |
| Root area | City Bike Station Data |
| Run | `0002-city-bike-stations-quality-eval` |
| Scope | Full evaluation |
| Rigor | standard |
| Evaluation verdict | Unacceptable |
| Full report | [report.md](report.md) |
| Machine report | [report.json](report.json) |

## Verdict

The root's provenance gap binds the whole-model rating. Schema & structure is
target, but Provenance & freshness is unacceptable for the same missing manifest
metadata.

## Area Breakdown

| Area | Path | Area Rating | Area + Sub-Areas Rating | Factors |
| --- | --- | --- | --- | --- |
| City Bike Station Data | `/` | Unacceptable | Unacceptable | Fitness for use: Minimum; Provenance: Unacceptable |
| Schema & structure | `schema-structure` | Target | Target | Structural validity: Target |
| Provenance & freshness | `provenance-freshness` | Unacceptable | Unacceptable | Provenance: Unacceptable; Freshness: Target |

## Selected Findings

1. **High**  
   The newest station snapshot lists the output file and row count but leaves the upstream source and acquisition timestamp blank.  
   `metadata/manifest.yaml`  
   Assessment: `assessments/002-root-each-snapshot-names-its-source-and-acquisition-time.json`
2. **Medium**  
   The station table includes 183 of 184 active stations from the operations roster; station ST-184 is active in the roster and absent from the published table.  
   `data/stations.csv`  
   Assessment: `assessments/001-root-active-stations-are-represented-with-plausible-locations.json`

## Recommended Actions

Primary next action: use `002-record-snapshot-provenance`.

| Recommendation ID | Priority | Recommendation | Done criterion |
| --- | --- | --- | --- |
| `002-record-snapshot-provenance` | 1 | [Record source and acquisition metadata for every snapshot](recommendations/002-record-snapshot-provenance.md) | The snapshot-provenance requirement reaches target; every current snapshot names its upstream source, acquisition time, and transformation version in the manifest or lineage record. |
| `001-backfill-missing-active-station` | 2 | [Backfill the missing active station](recommendations/001-backfill-missing-active-station.md) | The active-station coverage requirement reaches target; every active station in the operations roster is present in the station table with plausible coordinates. |

## Scope & Limitations

Scope: **Full evaluation**

In scope: City Bike Station Data; Schema & structure; Provenance & freshness

- The run reviewed a representative current release and two prior snapshots; it did not audit the operator's upstream live system.
- Coordinate plausibility was judged against city-boundary and station-address evidence, not by field inspection.
- Private maintainer knowledge was not counted as provenance evidence.
