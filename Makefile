SHELL := /bin/bash
include .env
export


ifneq ("$(wildcard .env)","")
else
    $(shell cp ./doc/example-dot-env .env)
endif

.PHONY: analyze init pre-commit-init pre-commit-run help choco-pack choco-push

# Default target executed when no arguments are given to make.
all: help

######################
# Go
######################
lint:
	golangci-lint run

test:
	go test -v ./...

run:
	make build
	go run main.go get chatbots --name hr

build:
	go get -v -t ./...
	go build -o smarter main.go

analyze:
	cloc . --exclude-ext=svg,json,zip --fullpath --not-match-d=smarter/smarter/static/assets/ --vcs=git

choco-pack:
	choco pack

choco-push: choco-pack
	choco push smarter.0.2.0.nupkg --source https://push.chocolatey.org/

# initialize local development environment.
# takes around 5 minutes to complete
init:
	@echo "Initializing local development environment..." & \
	GO111MODULE=off go install golang.org/x/tools/cmd/goimports & \
	GO111MODULE=on go install github.com/golangci/golangci-lint/cmd/golangci-lint@v1.57.2 & \
	make pre-commit-init	# install and configure pre-commit & \
	npm install


python-init:
	python3 -m venv venv
	source venv/bin/activate

pre-commit-init:
	pre-commit install
	pre-commit autoupdate

pre-commit-run:
	pre-commit run --all-files

release:
	git commit -m "fix: force a new release" --allow-empty && git push

######################
# HELP
######################

help:
	@echo '===================================================================='
	@echo 'analyze                - Analyze the project with cloc'
	@echo 'init                   - Initialize local and Docker environments'
	@echo 'lint            		  - run all code linters'
	@echo 'test                   - run all tests'
	@echo 'run                    - run the main.go file'
	@echo 'build                  - build the smarter cli'
	@echo 'choco-pack             - package the smarter cli for chocolatey'
	@echo 'choco-push             - push the smarter cli package to chocolatey'
	@echo 'pre-commit-init        - install and configure pre-commit'
	@echo 'pre-commit-run         - runs all pre-commit hooks on all files'
	@echo 'release				  - force a new release'
	@echo '===================================================================='
