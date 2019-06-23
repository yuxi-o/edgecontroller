# Controller CE

This is the project for the Controller Community Edition.

## Overview

To start the whole controller stack with default settings:

`make all-up`

This project uses a MySQL database running inside a Docker container. Please
make sure you have Docker installed locally. To start/reset the database, run:

`make db-reset`

To run unit and integration tests:

`make test`

To run only unit tests:

`make test-unit`

To run only integration tests (this will call `make db-reset`):

`make test-api`

## Building

To build the binaries:

`make build`

## Running

The Controller runs in two orchestration modes with the default being Docker
Native, where the appliance directly uses the Docker daemon to control
container instances of applications. The other is Kubernetes mode where an
external Kubernetes master is used to orchestrate the applications and is
controlled by its API.

### Simulating the Appliance

To run the test node which simulates an appliance listening as a gRPC server:

```
./dist/test-node -ela-port 42101 -eva-port 42102
```

### Running the Controller in Docker Native Mode

Configure the needed environment variables

- MYSQL_ROOT_PASSWORD: The password for database user `root`.
- CCE_ADMIN_PASSWORD: The password for the controller UI user `admin`.

These can be exported to the system evironment variables, or defined in the `./.env` file.

To start only the controller backend API:

```
make cce-up
```

### Running the Controller in Kubernetes Mode

```
./dist/cce -dsn "root:<db_pass>@tcp(:8083)/controller_ce" -adminPass <admin_pass> -httpPort 8080 -grpcPort 8081 -elaPort 42101 -evaPort 42102 -orchestration-mode "kubernetes" -k8s-client-ca-path <k8s_client_ca_path> -k8s-client-cert-path <k8s_client_cert_path> -k8s-client-key-path <k8s_client_key_path> -k8s-master-host <k8s_master_host> -k8s-api-path <k8s_api_path> -k8s-master-user <k8s_master_user>`
```

The flags are:

- <db_pass>: The password for database user `root` (**required**).
- <admin_pass>: The password for the API user `admin` (**required**).
- <k8s_client_ca_path>: Kubernetes root certificate path (**required**).
- <k8s_client_cert_path>: Kubernetes client certificate path (**required**).
- <k8s_client_key_path>: Kubernetes client private key path (**required**).
- <k8s_master_host>: Kubernetes master host (**required**).
- <k8s_api_path>: Kubernetes api path (default: `/`)
- <k8s_master_user>: Kubernetes default user (default: none)

### Running the CCE-UI & CUPS-UI
- You must have the corresponding .env.production environment files configured per the Readme's in the UI projects
- Use `Make` in order to spin up the production level docker-containers via docker-compose
- eg `make cups-ui-up && make ui-up`

## Project Layout

See [this link](https://medium.com/@benbjohnson/standard-package-layout-7cdbc8391fc1)
for guidelines.

## Architecture Notes

- Domain models (types that implement the Persistable interface) decouple application logic from PB types

### Controller CE UI

#### Development Environment

For instructions on how to setup your development environment, see `ui/controller/README.md`.

#### Production Deployment

##### Prerequisites

- For production deployments, create an `.env.production` file in the
  `ui/controller` project and define the required variables. See the
  `ui/controller/README.md` for details.

##### Running

```
make ui-up
```

##### Stopping

```
make ui-down
```

### CUPS UI

#### Development Environment

For instructions on how to setup your development environment, see `ui/cups/README.md`.

#### Production Deployment

##### Prerequisites
- For production deployments, create an `.env.production` file in the `ui/cups`
  project and define the required variables. See the `ui/cups/README.md` for
  details.

##### Running

```
make cups-ui-up
```

##### Stopping

```
make cups-ui-down
```

### Log Service

There is a syslog-ng service that receives log events via the Syslog protocol
over UDP on port 514 and stores the logs in `syslog/log/messages-kv.log` on
the Controller host. The syslog-ng
[Docker](https://github.com/balabit/syslog-ng-docker) container is started along
with the other services in the docker-compose.yml.

### Statistics Service

There is a StatsD service that receives statistics via the StatsD protocol over
UDP on port 8125 and stores the statistics in `statsd/stats.log` on the
Controller host. The gostatsd [Docker](https://github.com/atlassian/gostatsd)
container is started along with the other services in the docker-compose.yml.

## Test Coverage

### Kubernetes App Deployment

Kubernetes app deployment testing is done in 2 different layers.

- Testing the wrapper to the Kubernetes API with a minikube backend
- Testing the Controller REST API with a Controller in Kubernetes mode and a
  minikube backend

To run k8s tests locally, you need to install `minikube` and `kubectl` and add
them to PATH and start minikube with `minikube start`
