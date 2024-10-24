# Asset measurements assignment

Your task is to implement two distinct services. The first service will manage asset configurations
and store their measurements, while the second service will simulate assets by generating real-time
measurement data for those assets.

## Requirements

- Use these technologies:

  - Golang,
  - Postgresql for relational data,
  - MongoDB for time series data, and
  - RabbitMQ for messaging.

- The finished solution should be runnable locally - Update this README file with instructions.

## Asset Service:

This service should hold information about assets and store their measurements.

- Asset entity looks like this:
  - id
  - name
  - description
  - type
  - enabled
- Create a rest api functions for getting and managing assets:
  - Create,
  - Update,
  - Delete,
  - Get by id, and
  - Get all (optionally filter by 'enabled' and 'type' fields, order by name).
- Consume asset measurements from RabbitMQ (from simulator service, read below) and save them to db
  (measurements should not be saved in db if asset is disabled)
- Create rest api functions for getting asset measurements:
  - Get latest measurements record of a specific asset.
  - Get measurements of a specific asset in the provided time interval (asc/desc).
  - Get average measurements of a specific asset in the provided time interval, grouped by minute,
    15min, hour (asc/desc).

## Simulator Service

Service for simulating power of assets.

- Asset simulation config should be stored in the database and should contain:
  - id
  - assetId (same as on asset service)
  - type (battery, solar, wind...)
  - measurement interval (seconds)
  - min, max power
  - max power step (if value is \<= 0, there is no limit on power step)
- Service should generate asset measurements (power \[W\], SOE \[%\]) according to the simulation
  config (above) and send them to asset service via RabbitMQ.
- Simulate SOE (state of energy) of assets according to generated power values.

## Optional Requirements:

- Create rest api for crud operations for asset simulation config.
- Track changes of asset simulation config, add api which would return state of config at certain
  timestamp.
- Tests.

