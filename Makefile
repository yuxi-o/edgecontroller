# Copyright 2019 Smart-Edge.com, Inc. All rights reserved.
#
# Licensed under the Apache License, Version 2.0 (the "License");
# you may not use this file except in compliance with the License.
# You may obtain a copy of the License at
#
#     http://www.apache.org/licenses/LICENSE-2.0
#
# Unless required by applicable law or agreed to in writing, software
# distributed under the License is distributed on an "AS IS" BASIS,
# WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
# See the License for the specific language governing permissions and
# limitations under the License.

export GO111MODULE = on

.PHONY: help clean build lint test \
	db-up db-reset db-down \
	statsd-up statsd-down syslog-up syslog-down \
	ui-up ui-down ui-test \
	test-unit test-syslog test-statsd test-api

help:
	@echo "Please use \`make <target>\` where <target> is one of"
	@echo "Building:"
	@echo "  clean            to clean up build artifacts and docker volumes"
	@echo "  build            to build the project to the ./dist/ folder"
	@echo ""
	@echo "Services:"
	@echo "  db-up            to start the MySQL database service"
	@echo "  db-reset         to start and reset the MySQL database service"
	@echo "  db-down          to stop the MySQL database service"
	@echo ""
	@echo "  syslog-up        to start the syslog service"
	@echo "  syslog-down      to stop the syslog service"
	@echo ""
	@echo "  statsd-up        to start the statsd service"
	@echo "  statsd-down      to stop the statsd service"
	@echo ""
	@echo "Testing:"
	@echo "  lint             to run linters and static analysis on the code"
	@echo "  test             to run all tests"
	@echo "  test-unit        to run unit tests"
	@echo "  test-api         to run api tests"
	@echo "  test-syslog      to run syslog tests"
	@echo "  test-statsd      to run statsd tests"
	@echo "  test             to run unit followed by api tests"
	@echo ""
	@echo "  ui-up            to start the production UI Container"
	@echo "  ui-down          to stop the production UI container"
	@echo "  ui-dev-up        to start local developer instance of the UI"
	@echo "  ui-test          run the UI project tests"

clean:
	rm -rf dist certificates statsd/stats.log syslog/logs

build:
	go build -o dist/cce ./cmd/cce
	go build -o dist/test-node ./test/node/grpc

lint:
	golangci-lint run

db-up:
	docker-compose up -d mysql
	until mysql -P 8083 --protocol tcp -uroot -pbeer -e '' 2>/dev/null; do \
		echo "Waiting for DB..."; \
		sleep 1; \
		done

db-reset: db-up
	mysql -P 8083 --protocol tcp -u root -pbeer < mysql/schema.sql

db-down:
	docker-compose stop mysql

statsd-up:
	docker-compose up -d statsd

statsd-down:
	docker-compose stop statsd

syslog-up:
	docker-compose up -d syslog

syslog-down:
	docker-compose stop syslog

ui-up:
	docker build -t cce-ui ./ui
	docker-compose up -d ui

ui-down:
	docker-compose stop ui

ui-dev-up:
	cd ui/ && yarn install && yarn start

ui-test:
	cd ui/ && yarn install && yarn build && yarn test

test-unit:
	ginkgo -v -r --randomizeAllSpecs --randomizeSuites \
		--skipPackage=vendor,statsd,syslog,cmd/cce

test-statsd: statsd-up
	ginkgo -v --randomizeAllSpecs --randomizeSuites statsd

test-syslog: syslog-up
	ginkgo -v --randomizeAllSpecs --randomizeSuites syslog

test-api: db-reset
	ginkgo -v --randomizeAllSpecs --randomizeSuites cmd/cce

test: test-unit test-statsd test-syslog test-api
