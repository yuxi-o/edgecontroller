```text
SPDX-License-Identifier: Apache-2.0
Copyright © 2019 Intel Corporation and Smart-Edge.com, Inc.
```

# CHANGELOG

## [0.0.59] - 2019-07-03
### Changed
- Updated version of `github.com/smartedgemec/log` used.

## [0.0.58] - 2019-07-01
### Changed
- Updated version of `github.com/smartedgemec/log` used.

### Fixed
- Building the docker image no longer clones `github.com/smartedgemec/log` since `go.mod` is used.

## [0.0.57] - 2019-06-29
### Changed
- Updated `ansible/README.md` with more information and define host requirements.
- Updated `ansible/README.md` to explain how to install tag release of controller.
- Ansible provides configuration with `.env` file.
- Ansible brings up the controller using `make`.
- Ansible installs to `/opt/controller` by default.
- Ansible sets a random password for admin and MySQL. Credentials are stored
    under `credentials/` from where you run `ansible-playbook`.

### Fixed
- Ansible installs missing dependencies such as `git`, `make`, and `epel-release`.
- `make` no longer displays comments

## [0.0.56] - 2019-06-28
### Changed
- Backend API times out HTTP requests after 120 seconds
- Backend API rejects requests payloads larger than 64 KB
- Backend API times out non-singular DB queries after 10 seconds

### Fixed
- Auth (/auth) endpoint was not setting the content type
- DELETE /nodes/{node_id}/apps/{app_id} endpoint did not safely return in an error case
- README now makes the build requirements clearer when editing the `.env` file

## [0.0.55] - 2019-06-27
### Fixed
- Persistence on disk now has more organized directory structure
- Controller backend was not persisting CA key+cert and the log/stats files
- Some `make` commands were behaving non-deterministically

## [0.0.54] - 2019-06-27
### Changed
- UI API request timeouts set to 120seconds

## [0.0.53] - 2019-06-27
### Added
- Limits for application cores, memory, and ports
- (INTC-762) Add unit size to App memory

## [0.0.52] - 2019-06-26
### Added
- Added ansible playbook for controller deployment

## [0.0.51] - 2019-06-26
### Added
- Configuration options to support private Git repositories if there are private dependencies

### Changed
- Backend build inside of container now supports GitHub tokens instead of SSH keys

## [0.0.50] - 2019-06-26
### Added
- Comprehensive README additions with complementary commands from the Makefile

## [0.0.49] - 2019-06-25
### Fixed
- Expose Syslog and StatsD service ports in `docker-compose.yml`
- Persist telemetry logs on Docker host
- UI AddNode & AddApp hangs on 400 response from API
- UI not loading page on very first load

### Removed
- Hard-coded MySQL passwords
- Hard-coded admin user password

### Changed
- Moved CCE UI and CUPS UI environment configuration to root project `.env` file
- CCE run flags are configurable via the root project `.env` file

## [0.0.48] - 2019-06-25
### Fixed
- PATCH /nodes/{node_id}/dns endpoint now returns a 501 when forwarders are provided, as they are currently unsupported

## [0.0.47] - 2019-06-25
### Fixed
- UI: Form schema in TrafficPolicy had incorrect MacAddress type
- UI: Form schema on GTPFilter.IMSIs now of type String

## [0.0.46] - 2019-06-24
### Fixed
- UI: Traffic Policies on Interfaces & Applications did not show state

## [0.0.45] - 2019-06-23
### Fixed
- The following endpoints are now Swagger schema compliant:
    - GET, POST /nodes
    - GET, PATCH, DELETE /nodes/{node_id}
    - GET, PATCH, DELETE /nodes/{node_id}/interfaces/{interface_id}/policy
    - GET, PATCH, DELETE /nodes/{node_id}/dns
    - GET, POST /nodes/{node_id}/apps
    - GET, PATCH, DELETE /nodes/{node_id}/apps/{app_id}
- User interface is compliant to the Swagger schema

## [0.0.44] - 2019-06-19
### Added
- Controller API now builds and runs in a container

## [0.0.43] - 2019-06-17
### Fixed
- The following endpoints are now Swagger schema compliant:
    - GET, PATCH, DELETE /nodes/{node_id}/apps/{app_id}/policy

## [0.0.42] - 2019-06-17
### Fixed
- The following endpoints are now Swagger schema compliant:
    - GET, PATCH /nodes/{node_id}/interfaces
    - GET /nodes/{node_id}/interfaces/{interface_id}

## [0.0.41] - 2019-06-14
### Fixed
- The following endpoints are now Swagger schema compliant:
    - GET, POST /apps
    - GET, PATCH, DELETE /apps/{app_id}
    - GET, PATCH /nodes/{node_id}/interfaces
    - GET /nodes/{node_id}/interfaces/{interface_id}

## [0.0.40] - 2019-06-14
### Fixed
- Missing .env.development file in `ui/controller` which caused running the CCE
  UI in development mode to fail.
- Apps view improperly handling fatal errors when an empty response is
  received from API.

#### Changed
- Updated nodes, apps, policies, node apps, and node interfaces view to expect
  an object wrapped array from CCE API calls to resource collections for schema
  compliance.

#### Added
- Traffic policies view fetches traffic policies from API.

## [0.0.39] - 2019-06-14
### Changed
- Move protobuf app message types to EVA namespace

## [0.0.38] - 2019-06-14
### Fixed
- Edge Node gRPC credentials must have client key usage

## [0.0.37] - 2019-06-13
### Fixed
- Content Security Policy (CSP) did not whitelist the Controller and CUPS API as a `connect-src`.
- Controller API properly handles Cross-Origin Resource Sharing (CORS).
- CUPS UI used development CUPS API in production environments.

### Removed
- CUPS UI mock client and server due to unimplemented status.

### Changed
- Moved CUPS UI project from `cups-ui` to `ui/cups`.
- Moved Controller UI project from `ui` to `ui/controller`.
- Replaced UI production NGINX servers with a Node.js server with nonce injection for CSP compliance.
- Base Docker image for CCE UI production bundle builder from `node:latest` to `node:lts-alpine`.
- Base Docker image for CUPS UI production bundle builder from `node:latest` to `node:lts-alpine`.
- Base Docker image for CCE UI production server from `nginx:alpine` to `node:lts-alpine`.
- Base Docker image for CUPS UI production server from `nginx:alpine` to `node:lts-alpine`.

## [0.0.36] - 2019-06-12
### Added
- Dial edge node on separate ELA and EVA ports
- Dial edge node with TLS
- elaPort and evaPort configuration flags

### Changed
- Update protobuf for latest schemas

## [0.0.35] - 2019-06-11
### Fixed
- Increased node_grpc_targets.grpc_target max length to 47 to handle IPv6 addresses

## [0.0.34] - 2019-06-11
### Fixed
- Passwords in an auth request payload were not scrubbed from logs

## [0.0.33] - 2019-06-10
### Fixed
- UI did not notify user of successful or failed operations

## [0.0.32] - 2019-06-07
### Added
- Inital Controller UI

## [0.0.31] - 2019-06-06
### Added
- Explicitly whitelist allowed filter queries in GET /{resource} endpoints
- Improve security of default NGINX config
- Increase minimum TLS version to TLS1.2
- Limit TLS to single secure cipher suite

## [0.0.30] - 2019-06-04
### Added
- Node Interfaces and Interface Traffic Policies

## [0.0.29] - 2019-06-03
### Added
- GetContainerByIP gRPC endpoint

## [0.0.28] - 2019-05-31
### Added
- Syslog ingress server
- StatsD ingress server

## [0.0.27] - 2019-05-29
### Changed
- Move GetStatus invocation from Application DeploymentServiceClient to ApplicationLifecycleServiceClient

## [0.0.26] - 2019-05-23
### Added
- Support for apps to specify ports and protocols that need to be exposed

## [0.0.25] - 2019-05-24
### Added
- Kubernetes App deployment

## [0.0.24] - 2019-05-22
### Changed
- Change EC key generation from P256 to P384

## [0.0.23] - 2019-05-22
### Added
- Add -log-level flag

### Changed
- Use github.com/smartedgemec/log

## [0.0.22] - 2019-05-22
### Removed
- VNF models/endpoints (deferred from June release)

## [0.0.21] - 2019-05-22
### Changed
- Auth token signer uses a unique key to sign tokens

## [0.0.20] - 2019-05-18
### Added
- Node VNF gRPC calls and API tests
- Node VNF req/resp types
- Node DNS config gRPC calls
- Node App Traffic policies gRPC calls and API tests
- Node VNF Traffic policies API tests

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
- Node App gRPC calls and API tests
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
