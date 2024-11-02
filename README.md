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

# Solution overview

## Running the services

The services are dockerized and can be run using docker-compose. You need to have `docker` and `docker-compose`
installed on your machine.

To run the services, clone the repository and run the following command in the root directory of the repository:

```bash
docker compose -f ./deployment/docker-compose.yaml up -d --build
```

The docker compose deployment will automatically start the asset and simulator services as well as their dependencies
(PostgreSQL, MongoDB, and RabbitMQ). Database migrations are run automatically on service startup.

## Testing the services

For easier local testing, we will use a Traefik as a reverse proxy to route the requests to the services based on the
domain name:

- `asset-service.localhost` for the asset service
- `simulator.localhost` for the simulator service

## Notes

What could be improved:

- Security: The services are not secured in any way. There is no authentication or authorization implemented.
- Pagination: The API endpoints do not support pagination. This could be a problem when the number of assets or
  measurements grows.

Compromises made:

- One of the compromises made in the implementation is omitting API model separation from domain models. Currently,
  domain model = API model. This is not a good practice in a real-world application, but it was done to save time.
- Guaranteeing proper unit conversions when generating and storing measurements. I've made the Power model with Unit and
  Measurement fields, but I didn't implement the conversion logic, as the requirements did not specify multiple units.

