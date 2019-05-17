# CHANGELOG

## [0.0.19] - 2019-05-17
### Added
- Check Node approval before enrollment
- Store Node IP on enrollment

## [0.0.18] - 2019-05-17
### Added
- Version, source fields to App and VNF models

### Removed
- Image field from App and VNF models

## [0.0.17] - 2019-05-17
### Added
- Node App gRPC calls and tests
- Node App req/resp types

### Changed
- Refactored model/persistence interfaces

## [0.0.16] - 2019-05-16
### Added
- Token authentication login endpoint
- Require token authentication on all HTTP endpoints

### Fixed
- Update instances where HTTP response body was not being closed

## [0.0.15] - 2019-05-15
### Added
- Add TLS to gRPC Server
- Print self-signed CA at startup to output

## [0.0.14] - 2019-05-14
### Added
- Node App tests, application logic handler func framework

## [0.0.13] - 2019-05-14
### Added
- Node VNF traffic policy model + tests

## [0.0.12] - 2019-05-14
### Removed
- Unneeded JSON_EXTRACT usage in queries

## [0.0.11] - 2019-05-13
### Added
- StatsD service for stats telemetry support

## [0.0.10] - 2019-05-13
### Fix
- Run test-syslog without sudo

## [0.0.9] - 2019-05-13
### Added
- DNSService gRPC client

## [0.0.8] - 2019-05-09
### Added
- Syslog service for log telemetry support

## [0.0.7] - 2019-05-07
### Added
- gRPC listener and node auth gRPC endpoint

## [0.0.6] - 2019-05-02
### Added
- Self-sign CA on startup

## [0.0.5] - 2019-04-30
### Added
- Node <-> DNS config model

### Fixed
- JSON naming errors in traffic policy models

## [0.0.4] - 2019-04-30
### Added
- Dockerized MySQL DB, persistence service, HTTP server, API tests

## [0.0.3] - 2019-04-23
### Added
- Entity models and SQL schema

## [0.0.2] - 2019-04-18
### Added
- Service interfaces and application constants

## [0.0.1] - 2019-04-17
### Added
- `uuid` package to abstract implementation

## [0.0.0] - 2019-04-12
### Added
- README, CHANGELOG
