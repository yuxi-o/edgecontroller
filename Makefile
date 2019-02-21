export GO111MODULE = on

.PHONY: help clean build lint test

help:
	@echo "Please use \`make <target>\` where <target> is one of"
	@echo "  clean          to clean up build artifacts and docker"
	@echo "  build          to build the release docker image"
	@echo "  lint           to run linters and static analysis on the code"
	@echo "  test           to run unit tests"

clean:
	rm -rf dist

build:
	mkdir -p dist
	go build -o dist/helloapp ./cmd/helloapp

lint:
	gometalinter --config=.gometalinter-config.json ./...

test:
	ginkgo -v -r --randomizeAllSpecs --randomizeSuites --failOnPending --skipPackage=vendor
