SHELL := /bin/bash
DATE ?= $(shell date '+%Y-%m-%d %H:%M:%S')
BASE_PATH := $(shell dirname $(realpath $(lastword $(MAKEFILE_LIST))))
VERSION ?= $(shell git describe --tags --always --match=* 2> /dev/null || \
           			cat $(CURDIR)/.version 2> /dev/null || echo v0)
VERSION_HASH = $(shell git rev-parse HEAD)
BRANCH = $(shell git rev-parse --abbrev-ref HEAD)

go = GOGC=off go
MODULE = $(shell env GO111MODULE=on go list -m)

LDFLAGS += -X "$(MODULE)/version.Version=$(VERSION)" -X "$(MODULE)/version.CommitSHA=$(VERSION_HASH)"

.PHONY: build
build: ## Build
	$Q $(go) build -ldflags '$(LDFLAGS)' -o .