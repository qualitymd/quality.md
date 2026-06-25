---
title: City Bike Station Data
ratingScale:
  - level: outstanding
    title: Outstanding
    description: "The stretch band - reached only with significant extra effort."
    criterion: "Exceeds the requirement; satisfies it with margin to spare."
  - level: target
    title: Target
    description: "The level to aim for - achievable at reasonable cost and effort."
    criterion: "Satisfies the requirement."
  - level: minimum
    title: Minimum
    description: "The acceptable floor - less than you'd aim for, but consciously agreed as good enough to rely on."
    criterion: "Falls short of the target but remains acceptable."
  - level: unacceptable
    title: Unacceptable
    description: "Below the floor - not good enough to rely on."
    criterion: "Does not meet the requirement to an acceptable degree."
factors:
  fitness-for-use:
    title: Fitness for use
    description: >
      Fitness for use is the degree to which planners and public-information
      maintainers can rely on the station data to answer ordinary service
      questions: where stations are, whether they are active, and how current
      each snapshot is. It is the "can this data support the intended decisions"
      lens, distinct from Provenance's "can a reader trace where the data came
      from and when".
    requirements:
      active-stations-are-represented-with-plausible-locations:
        title: Active stations are represented with plausible locations
        assessment: >
          Review the station register against the published operations roster
          and profile each active station's latitude/longitude against city
          boundaries and known station addresses. The requirement is met when
          active stations are present and coordinate outliers are explained or
          corrected; the reviewer records missing stations and implausible
          locations as defects.
  provenance:
    title: Provenance
    description: >
      Provenance is the degree to which a maintainer can tell which upstream
      source, acquisition time, and transformation produced a data snapshot; it
      matters because decisions based on stale or unattributed data cannot be
      audited or corrected. It is the "where did this data come from" lens,
      distinct from Fitness for use's "does it support the decision".
    requirements:
      each-snapshot-names-its-source-and-acquisition-time:
        title: Each snapshot names its source and acquisition time
        assessment: >
          Inspect the dataset manifest and sample snapshots for source-system
          name, acquisition timestamp, and transformation version. The
          requirement is met when every reviewed snapshot can be traced to its
          source and acquisition time without relying on tribal knowledge.
areas:
  schema-structure:
    title: Schema & structure
    description: >
      The tabular station register and its schema definition: station
      identifiers, required fields, field types, and documented value ranges.
    source: ./data/stations.csv and ./schema/stations.schema.yaml
    factors:
      structural-validity:
        title: Structural validity
        description: >
          Structural validity is the degree to which the data's columns and
          values match the published schema closely enough for downstream users
          to interpret every row consistently.
        requirements:
          required-columns-conform-to-the-station-schema:
            title: Required columns conform to the station schema
            assessment: >
              Compare the station table with the published schema and profile
              the required fields for missing values, out-of-range values, and
              type mismatches. The requirement is met when every required column
              is present and deviations are either absent or explained by a
              documented exception.
          station-identifiers-are-unique-and-stable-across-snapshots:
            title: Station identifiers are unique and stable across snapshots
            assessment: >
              Compare the current station identifiers with the previous two
              snapshots and the change log. The requirement is met when no
              duplicate identifiers exist and any identifier change is tied to a
              documented station merge, split, or retirement.
  provenance-freshness:
    title: Provenance & freshness
    description: >
      The source manifest, acquisition metadata, and update cadence that tell a
      reader where the station data came from and whether it is current enough
      for planning and public-information use.
    source: ./metadata/manifest.yaml and ./metadata/lineage.md
    factors:
      provenance:
        title: Provenance
        description: >
          Provenance here is the traceability of the dataset's upstream source,
          acquisition time, and transformation notes for each snapshot.
        requirements:
          each-snapshot-names-its-source-and-acquisition-time-2:
            title: Each snapshot names its source and acquisition time
            assessment: >
              Inspect the manifest and lineage note for every snapshot in the
              current release. The requirement is met when a reader can identify
              the upstream source, acquisition time, and transformation version
              for each snapshot without asking the data maintainer.
      freshness:
        title: Freshness
        description: >
          Freshness is the degree to which the published snapshot cadence and
          visible timestamps make the dataset current enough for station
          planning and public-information updates.
        requirements:
          published-data-is-refreshed-within-the-documented-weekly-cadence:
            title: Published data is refreshed within the documented weekly cadence
            assessment: >
              Compare the newest snapshot timestamp with the documented weekly
              cadence and review release notes for any announced outage. The
              requirement is met when the latest snapshot is no more than seven
              days old or the delay is explicitly announced with a recovery date.
---

# City Bike Station Data - Quality model

## Overview

This illustrative model describes a city bike-share station data product: a
station register, its schema, and the manifest that records where each snapshot
came from. Good means the dataset supports ordinary planning and
public-information decisions without forcing a reviewer to guess whether a row,
coordinate, timestamp, or source is trustworthy.

## Scope

This model covers the checked-in station table (`./data/stations.csv`), schema
definition (`./schema/stations.schema.yaml`), and metadata/lineage files under
`./metadata`. It does not cover the bike-share operator's live telemetry system,
trip records, pricing data, or the public website that consumes the station
register.

## Needs

- Planners can see all active stations with plausible coordinates.
- Data consumers can parse the table consistently from the published schema.
- Maintainers can trace each snapshot to an upstream source and acquisition time.
- Public-information updates are based on data that is current enough for
  weekly station changes.

## Risks

A missing active station or implausible coordinate can send maintenance crews or
riders to the wrong place. A schema drift can break downstream imports without a
clear owner. Missing provenance makes stale or erroneous data hard to correct
because the maintainer cannot tell which source or transformation produced it.
