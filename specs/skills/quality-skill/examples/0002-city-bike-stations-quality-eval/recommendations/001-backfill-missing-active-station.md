---
schemaVersion: 1
title: Backfill the missing active station
gap: One active station from the operations roster is absent from the published station table.
evidenceLocators:
  - data/stations.csv
  - ops/active-stations-roster.csv
assessmentResultRecords:
  - assessments/001-root-active-stations-are-represented-with-plausible-locations.json
remediationOptions:
  - Add ST-184 to the station table and rerun the roster-to-table review.
  - Mark ST-184 inactive in the operations roster if the roster is wrong.
  - Leave the gap documented in release notes only.
recommendedOption: Add ST-184 to the station table and rerun the roster-to-table review.
doneCriterion: The active-station coverage requirement reaches target; every active station in the operations roster is present in the station table with plausible coordinates.
---

# Backfill the missing active station

**Area / factor:** City Bike Station Data (root) -> Fitness for use
**In-scope requirement:** *Active stations are represented with plausible locations*
**Current rating:** Minimum.

## Gap

The station table omits one active station from the operations roster. The
dataset remains usable for most planning work, but a missing active station is a
public-information and operations risk because consumers can treat the station as
closed or nonexistent.

**Evidence**

- `data/stations.csv` - current published station table.
- `ops/active-stations-roster.csv` - station `ST-184` appears as active.
- Manual roster-to-table comparison found `183/184` active stations present.

## Options

- **(a) Add ST-184 to the station table and rerun the roster-to-table review.**
  This fixes the published data and verifies the gap is closed.
- **(b) Mark ST-184 inactive in the operations roster if the roster is wrong.**
  This is correct only if operations confirms the station is no longer active.
- **(c) Leave the gap documented in release notes only.** This preserves a known
  defect and keeps the requirement below target.

## Recommended

**(a) Add ST-184 to the station table and rerun the roster-to-table review.**
The defect is in the published table unless operations proves otherwise.

## Done-criterion

The requirement *Active stations are represented with plausible locations*
reaches **Target**: every active station in the operations roster is present in
the station table and coordinate review finds no unexplained outliers.
