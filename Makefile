#!/usr/bin/make -f

PACKAGES=$(shell go list ./...)
VERSION := $(shell echo $(shell git describe --tags) | sed 's/^v//')
COMMIT := $(shell git log -1 --format='%H')
GOBIN ?= $(GOPATH)/bin

export GO111MODULE = on

# process linker flags

ldflags = -X github.com/QOSGroup/qos/version.Version=$(VERSION) \
  -X github.com/QOSGroup/qos/version.Commit=$(COMMIT)

ldflags += $(LDFLAGS)
ldflags := $(strip $(ldflags))

BUILD_FLAGS := -ldflags '$(ldflags)'

all: install test

########################################
### Build/Install

build:
ifeq ($(OS),Windows_NT)
	go build -mod=readonly $(BUILD_FLAGS) -o build/qosd.exe ./cmd/qosd
	go build -mod=readonly $(BUILD_FLAGS) -o build/qoscli.exe ./cmd/qoscli
else
	go build -mod=readonly $(BUILD_FLAGS) -o build/qosd ./cmd/qosd
	go build -mod=readonly $(BUILD_FLAGS) -o build/qoscli ./cmd/qoscli
endif

build-linux:
	LEDGER_ENABLED=false GOOS=linux GOARCH=amd64 $(MAKE) build

install:
	go install -mod=readonly $(BUILD_FLAGS) ./cmd/qosd
	go install -mod=readonly $(BUILD_FLAGS) ./cmd/qoscli

########################################
### Testing

test: test_unit

test_unit:
	@VERSION=$(VERSION) go test -mod=readonly $(PACKAGES) -tags='ledger test_ledger_mock'

test_race:
	@VERSION=$(VERSION) go test -mod=readonly -race $(PACKAGES)

format:
	find . -name '*.go' -type f -not -path "*.git*" | xargs gofmt -w -s