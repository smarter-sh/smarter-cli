SHELL := /bin/bash
include .env
export


ifneq ("$(wildcard .env)","")
else
    $(shell cp ./doc/example-dot-env .env)
endif

.PHONY: analyze init pre-commit-init pre-commit-run help

# Default target executed when no arguments are given to make.
all: help

analyze:
	cloc . --exclude-ext=svg,json,zip --fullpath --not-match-d=smarter/smarter/static/assets/ --vcs=git

# initialize local development environment.
# takes around 5 minutes to complete
init:
	@echo "Initializing local development environment..." & \
	GO111MODULE=off go install golang.org/x/tools/cmd/goimports & \
	GO111MODULE=on go install github.com/golangci/golangci-lint/cmd/golangci-lint@v1.57.2 & \
	make pre-commit-init	# install and configure pre-commit


pre-commit-init:
	pre-commit install
	pre-commit autoupdate
	pre-commit run --all-files

pre-commit-run:
	pre-commit run --all-files

lint:
	golangci-lint run

######################
# Go lang
######################
run:
	go run main.go

######################
# HELP
######################

help:
	@echo '===================================================================='
	@echo 'init                   - Initialize local and Docker environments'
	@echo 'pre-commit-init        - install and configure pre-commit'
	@echo 'pre-commit-run         - runs all pre-commit hooks on all files'
	@echo '===================================================================='
