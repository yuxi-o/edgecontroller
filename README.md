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

`./dist/test-node -port 8082`


	// application orchestration mode
	flag.StringVar(&orchMode, "orchestration-mode", "native", "Orchestration mode. options [native, kubernetes] ")

	// k8s
	flag.StringVar(&k8sClientCAPath, "k8s-client-ca-path", "", "Kubernetes root certificate path")
	flag.StringVar(&k8sClientCertPath, "k8s-client-cert-path", "", "Kubernetes client certificate path")
	flag.StringVar(&k8sClientKeyPath, "k8s-client-key-path", "", "Kubernetes client private key path")
	flag.StringVar(&k8sHost, "k8s-master-host", "", "Kubernetes master host")
	flag.StringVar(&k8sAPIPath, "k8s-api-path", "", "Kubernetes api path")
	flag.StringVar(&k8sUser, "k8s-master-user", "", "Kubernetes default user")

### Running the Controller in Docker Native Mode

```
./dist/cce -orchestration-mode "native" -dsn "root:<db_pass>@tcp(:8083)/controller_ce" -adminPass <admin_pass> -httpPort 8080 -grpcPort 8081`
```

Replace the following **required** variables before executing the run command:

- <db_pass>: The password for database user `root`.
- <admin_pass>: The password for the API user `admin`.

### Running the Controller in Kubernetes Mode

```
./dist/cce -dsn "root:<db_pass>@tcp(:8083)/controller_ce" -adminPass <admin_pass> -httpPort 8080 -grpcPort 8081 -orchestration-mode "kubernetes" -k8s-client-ca-path <k8s_client_ca_path> -k8s-client-cert-path <k8s_client_cert_path> -k8s-client-key-path <k8s_client_key_path> -k8s-master-host <k8s_master_host> -k8s-api-path <k8s_api_path> -k8s-master-user <k8s_master_user>`
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

##### Prerequisites

- Node & NPM installed (v10.15.3, or V10 LTS)
  - recommended to use NVM https://github.com/nvm-sh/nvm to manage your Node versions
- Yarn installed globally `npm install -g yarn`
- install dependencies via `yarn install` within the project

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
