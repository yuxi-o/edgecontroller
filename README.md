# Controller CE

This is the project for the Controller Community Edition.

## Overview

This project uses a MySQL database. To create the schema in a local database,
run this command:

`mysql -u root -p < schema.sql`

To run the test node:

`go run github.com/smartedgemec/controller-ce/test/node/grpc -port 8081`

To run unit tests using the test node:

`make lint test`

## Project Layout

See [this link](https://medium.com/@benbjohnson/standard-package-layout-7cdbc8391fc1)
for guidelines.

## Architecture Notes

- Domain models (types that implement the Entity interface) decouple application
  logic from PB types
