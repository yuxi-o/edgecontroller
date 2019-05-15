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

.PHONY: help clean build lint test

help:
	@echo "Please use \`make <target>\` where <target> is one of"
	@echo "  clean            to clean up build artifacts and docker"
	@echo "  build            to build the release docker image"
	@echo "  lint             to run linters and static analysis on the code"
	@echo "  db-up            to start the MySQL database using docker-compose"
	@echo "  db-reset         to start and reset the MySQL database using docker-compose"
	@echo "  db-down          to stop the MySQL database using docker-compose"
	@echo "  test-unit        to run unit tests"
	@echo "  test-api         to run api tests"
	@echo "  test             to run unit followed by api tests"

clean:
	rm -rf dist certificates statsd/stats.log syslog/logs

build:
	mkdir -p dist
	go build -o dist/cce ./cmd/cce
	go build -o dist/test-node ./test/node/grpc

lint:
	golangci-lint run

db-up:
	docker-compose up -d mysql
	until mysql -P 8083 --protocol tcp -uroot -pbeer -e '' 2>/dev/null; do echo "Waiting for DB..."; sleep 1; done

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
	cd ui/ && yarn install && yarn start

ui-test:
	cd ui/ && yarn install && yarn build && yarn test

test-unit:
	ginkgo -v -r --randomizeAllSpecs --randomizeSuites --skipPackage=vendor,statsd,syslog,cmd/cce

test-statsd: statsd-up
	ginkgo -v --randomizeAllSpecs --randomizeSuites statsd

test-syslog: syslog-up
	ginkgo -v --randomizeAllSpecs --randomizeSuites syslog

test-api: db-reset
	ginkgo -v --randomizeAllSpecs --randomizeSuites cmd/cce

test: test-unit test-statsd test-syslog test-api
