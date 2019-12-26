# SPDX-License-Identifier: Apache-2.0
# Copyright (c) 2019 Intel Corporation

# Source user configured environment file.
include .env

export GO111MODULE = on
export GOPROXY = "https://goproxy.cn"
export MINIKUBE_WANTUPDATENOTIFICATION=false
export MINIKUBE_WANTREPORTERRORPROMPT=false
export MINIKUBE_HOME=$(HOME)
export CHANGE_MINIKUBE_NONE_USER=true
export KUBECONFIG=$(HOME)/.kube/config

# Build CCE run flags to be passed to docker-compose. This must be declared
# before any command declarations since docker-compose depends on this variable
# to always be set.
define CCE_FLAGS_BASE
	-adminPass $(CCE_ADMIN_PASSWORD) \
	-dsn root:$(MYSQL_ROOT_PASSWORD)@tcp(mysql:3306)/controller_ce \
	-log-level $(CCE_LOG_LEVEL)
endef

# Pass kubernetes related flags if and only if the user specified kubernetes
# as the orchestration mode. Otherwise, assume native orchestration mode and
# pass in base flags.
ifeq ($(CCE_ORCHESTRATION_MODE),$(filter $(CCE_ORCHESTRATION_MODE),kubernetes kubernetes-ovn))
define CCE_FLAGS
	$(CCE_FLAGS_BASE) \
	-orchestration-mode $(CCE_ORCHESTRATION_MODE) \
	-k8s-client-ca-path /artifacts/k8s/ca.pem \
	-k8s-client-cert-path /artifacts/k8s/cert.pem \
	-k8s-client-key-path /artifacts/k8s/key.pem \
	-k8s-master-host $(CCE_K8S_MASTER_HOST) \
	-k8s-api-path $(CCE_K8S_API_PATH) \
	-k8s-master-user $(CCE_K8S_MASTER_USER)
endef
	export CCE_FLAGS
else
	export CCE_FLAGS=$(CCE_FLAGS_BASE)
endif

define bring_ui_up 
	docker-compose up -d ui
	docker-compose up -d cups-ui
	docker-compose up -d cnca-ui
endef

.PHONY: help all-up all-down clean build build-dnscli lint test \
	db-up db-reset db-down \
	minikube-install kubectl-install minikube-wait \
	ui-up ui-down ui-test \
	test-k8s test-api-k8s \
	test-unit test-api test-dnscli

help:
	@echo "Please use \`make <target>\` where <target> is one of"
	@echo "Building:"
	@echo "  clean            to clean up build artifacts and docker volumes"
	@echo "  build            to build the project to the ./dist/ folder"
	@echo "  build-ifsvccli   to build interfaceservice CLI to the ./dist/ folder"
	@echo "  build-dnscli     to build edgednscli to the ./dist/ folder"
	@echo ""
	@echo "Services:"
	@echo "  all-up           to start the full controller stack"
	@echo "  all-down         to stop the full controller stack"
	@echo ""
	@echo "  db-up            to start the MySQL database service"
	@echo "  db-reset         to start and reset the MySQL database service"
	@echo "  db-down          to stop the MySQL database service"
	@echo ""
	@echo "  ui-up            to start the production UI Container"
	@echo "  ui-down          to stop the production UI container"
	@echo "  ui-dev-up        to start local developer instance of the UI"
	@echo "  ui-test          run the UI project tests"
	@echo ""
	@echo "  cups-ui-up       to start the production UI Container"
	@echo "  cups-ui-down     to stop the production UI container"
	@echo "  cups-ui-dev-up   to start local developer instance of the UI"
	@echo "  cups-ui-test     run the UI project tests"
	@echo ""
	@echo "  cnca-ui-up       to start the production UI Container"
	@echo "  cnca-ui-down     to stop the production UI container"
	@echo "  cnca-ui-dev-up   to start local developer instance of the UI"
	@echo "  cnca-ui-test     run the UI project tests"
	@echo ""
	@echo "  kubectl-install  to install kubectl"
	@echo "  minikube-install to install minikube"
	@echo "  minikube-start   to start minikube"
	@echo "  minikube-wait    to wait for minikube to be ready"
	@echo "  minikube-stop    to stop minikube"
	@echo ""
	@echo "Testing:"
	@echo "  lint             to run linters and static analysis on the code"
	@echo "  test             to run all tests"
	@echo "  test-unit        to run unit tests"
	@echo "  test-api         to run api tests"
	@echo "  test-api-k8s     to run k8s app deployment api tests"
	@echo "  test-k8s         to run kubernetes orchestration tests"
	@echo "  test-dnscli      to run edgednscli tests"
	@echo "  test             to run unit followed by api tests"

clean:
	@docker-compose stop
	@docker-compose rm
	rm -rf dist certificates artifacts

all-up: db-up cce-up ui-up

all-down: db-down cce-down ui-down

build:
	docker-compose build

	@# TODO: Remove the following when the test node is built as a Docker image and add it to the docker-compose.yml
	@# and add details to the README about running a test node.
	@###########################
	@# go build -o dist/test-node ./test/node/grpc
	@###########################

build-ifsvccli:
	go build -o dist/interfaceservicecli ./cmd/interfaceservicecli

lint:
	golangci-lint run

db-up:
	docker-compose up -d mysql
	@until mysql -P 8083 --protocol tcp -uroot -p$(MYSQL_ROOT_PASSWORD) -e '' 2>/dev/null; do \
		echo "Bringing up DB (this may take a moment...)"; \
		sleep 1; \
		done

	@# Either the DB already exists or it should run the schema.sql to create the DB
	@mysql -P 8083 --protocol tcp -u root -p$(MYSQL_ROOT_PASSWORD) -e '' controller_ce >/dev/null 2>&1 || \
	mysql -P 8083 --protocol tcp -u root -p$(MYSQL_ROOT_PASSWORD) < mysql/schema.sql >/dev/null 2>&1

db-reset:
	@# Checks if accessing the MySQL engine exits 0 (success); if so, try to drop the database
ifeq ($(shell mysql -P 8083 --protocol tcp -u root -p$(MYSQL_ROOT_PASSWORD) -e '' >/dev/null 2>&1; echo $$?),0)
	@mysql -P 8083 --protocol tcp -u root -p$(MYSQL_ROOT_PASSWORD) -e "DROP DATABASE IF EXISTS controller_ce;" >/dev/null 2>&1
endif

db-down:
	docker-compose stop mysql

kubectl-install:
ifeq ($(shell uname -s),Darwin)
	brew install kubernetes-cli
else
	curl -Lo kubectl \
		https://storage.googleapis.com/kubernetes-release/release/v1.14.2/bin/linux/amd64/kubectl \
		&& sudo install kubectl /usr/local/bin/
endif

minikube-install:
ifeq ($(shell uname -s),Darwin)
	brew cask install minikube
else
	curl -Lo minikube \
		https://storage.googleapis.com/minikube/releases/latest/minikube-linux-amd64 \
		&& sudo install minikube /usr/local/bin/
	mkdir -p $(HOME)/.kube $(HOME)/.minikube
	touch $(KUBECONFIG)
endif

minikube-start:
ifeq ($(shell uname -s),Darwin)
	minikube start
else
	sudo -E minikube start --vm-driver=none
	sudo chown -R travis: $(HOME)/.minikube/
endif

minikube-wait:
	kubectl cluster-info
	@# kube-addon-manager is responsible for managing other kubernetes components, such as kube-dns, dashboard, storage-provisioner.
	until kubectl -n kube-system get pods -lcomponent=kube-addon-manager -o jsonpath="{range .items[*]}{@.metadata.name}:{range @.status.conditions[*]}{@.type}={@.status};{end}{end}" 2>&1 | grep -q "Ready=True"; do \
		sleep 1; \
		echo "waiting for kube-addon-manager to be available"; \
		kubectl get pods --all-namespaces; \
	done
	@# Wait for kube-dns to be ready.
	until kubectl -n kube-system get pods -lk8s-app=kube-dns -o jsonpath="{range .items[*]}{@.metadata.name}:{range @.status.conditions[*]}{@.type}={@.status};{end}{end}" 2>&1 | grep -q "Ready=True"; do \
		sleep 1; \
		echo "waiting for kube-dns to be available"; \
		kubectl get pods --all-namespaces; \
	done

minikube-stop:
ifeq ($(shell uname -s),Darwin)
	minikube delete
else
	sudo -E minikube delete
endif

cce-up:
ifeq ($(CCE_ORCHESTRATION_MODE),$(filter $(CCE_ORCHESTRATION_MODE),kubernetes kubernetes-ovn))
	mkdir -p ./artifacts/controller/k8s
	cp ${CCE_K8S_CLIENT_CA_PATH} ./artifacts/controller/k8s/ca.pem
	cp ${CCE_K8S_CLIENT_CERT_PATH} ./artifacts/controller/k8s/cert.pem
	cp ${CCE_K8S_CLIENT_KEY_PATH} ./artifacts/controller/k8s/key.pem
endif
	docker-compose up -d cce

cce-down:
	docker-compose stop cce

ui-up:
	$(call bring_ui_up)

ui-down:
	docker-compose stop ui

ui-dev-up:
	cd ui/controller && yarn install && yarn start

ui-test:
	cd ui/controller && yarn install && yarn build && yarn test

cups-ui-up:
	docker-compose up -d cups-ui

cups-ui-down:
	docker-compose stop cups-ui

cups-ui-dev-up:
	cd ui/cups && yarn install && yarn start

cups-ui-test:
	cd ui/cups && yarn install && yarn build && yarn test

cnca-ui-up:
	docker-compose up -d cnca-ui

cnca-ui-down:
	docker-compose stop cnca-ui

cnca-ui-dev-up:
	cd ui/cnca && yarn install && yarn start

cnca-ui-test:
	cd ui/cnca && yarn install && yarn build && yarn test

build-dnscli:
	go build -o dist/edgednscli ./cmd/edgednscli

test-unit:
	ginkgo -v -r --randomizeAllSpecs --randomizeSuites \
		--skipPackage=vendor,k8s,cmd/cce,cmd/cce/k8s

test-api:
	$(MAKE) db-reset
	$(MAKE) db-up
	ginkgo -v --randomizeAllSpecs --randomizeSuites cmd/cce

test-api-k8s:
	$(MAKE) db-reset
	$(MAKE) db-up
	docker pull nginx:1.12
	ginkgo -v --randomizeAllSpecs --randomizeSuites cmd/cce/k8s

test-api-kubeovn:
	$(MAKE) db-reset
	$(MAKE) db-up
	docker pull nginx:1.12
	ginkgo -v --randomizeAllSpecs --randomizeSuites cmd/cce/kubeovn

test-k8s:
	ginkgo -v -r --randomizeAllSpecs --randomizeSuites k8s

test-dnscli:
	ginkgo -v -r --randomizeAllSpecs --randomizeSuites edgednscli

test: test-unit test-api test-k8s test-api-k8s test-api-kubeovn test-dnscli
