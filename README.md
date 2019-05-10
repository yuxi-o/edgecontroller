# Controller CE

This is the project for the Controller Community Edition.

## Overview

This project uses a MySQL database running inside a Docker container. Please
make sure you have Docker installed locally. To start/reset the database, run:

`make db-reset`

To run unit and integration tests:

`make test`

To run only unit tests:

`make test-unit`

To run only integration tests (this will call `make db-reset`):

`make test-api`

## Building and Running

To build the binaries:

`make build`

To run the test node which simulates an appliance listening as a gRPC server:

`./dist/test-node -port 8082`

To run the Controller CE:

`./dist/cce -dsn "root:<db_pass>@tcp(:8083)/controller_ce" -httpPort 8080 -grpcPort 8081`

## Project Layout

See [this link](https://medium.com/@benbjohnson/standard-package-layout-7cdbc8391fc1)
for guidelines.

## Architecture Notes

- Domain models (types that implement the Entity interface) decouple application
  logic from PB types
