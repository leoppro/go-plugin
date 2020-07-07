### Makefile for ticdc
.PHONY: build test check clean fmt cdc kafka_consumer coverage \
	integration_test_build integration_test integration_test_mysql integration_test_kafka hashicorp_plugin

PROJECT=ticdc

FAIL_ON_STDOUT := awk '{ print  } END { if (NR > 0) { exit 1  }  }'

CURDIR := $(shell pwd)
path_to_add := $(addsuffix /bin,$(subst :,/bin:,$(GOPATH)))
export PATH := $(path_to_add):$(PATH)

TEST_DIR := /tmp/tidb_cdc_test
SHELL	 := /usr/bin/env bash

GO       := GO111MODULE=on go
GOBUILD  := CGO_ENABLED=1 $(GO) build $(BUILD_FLAG) -trimpath
ifeq ($(GOVERSION114), 1)
GOTEST   := CGO_ENABLED=1 $(GO) test -p 3 --race -gcflags=all=-d=checkptr=0
else
GOTEST   := CGO_ENABLED=1 $(GO) test -p 3 --race
endif

ARCH  := "`uname -s`"
LINUX := "Linux"
MAC   := "Darwin"
PACKAGE_LIST := go list ./...| grep -vE 'vendor|proto|ticdc\/tests'
PACKAGES  := $$($(PACKAGE_LIST))
PACKAGE_DIRECTORIES := $(PACKAGE_LIST) | sed 's|github.com/pingcap/$(PROJECT)/||'
FILES := $$(find . -name '*.go' -type f | grep -vE 'vendor')
CDC_PKG := github.com/pingcap/ticdc
FAILPOINT_DIR := $$(for p in $(PACKAGES); do echo $${p\#"github.com/pingcap/$(PROJECT)/"}|grep -v "github.com/pingcap/$(PROJECT)"; done)
FAILPOINT := bin/failpoint-ctl

FAILPOINT_ENABLE  := $$(echo $(FAILPOINT_DIR) | xargs $(FAILPOINT) enable >/dev/null)
FAILPOINT_DISABLE := $$(find $(FAILPOINT_DIR) | xargs $(FAILPOINT) disable >/dev/null)

LDFLAGS += -X "$(CDC_PKG)/pkg/util.BuildTS=$(shell date -u '+%Y-%m-%d %H:%M:%S')"
LDFLAGS += -X "$(CDC_PKG)/pkg/util.GitHash=$(shell git rev-parse HEAD)"
LDFLAGS += -X "$(CDC_PKG)/pkg/util.ReleaseVersion=$(shell git describe --tags --dirty="-dev")"
LDFLAGS += -X "$(CDC_PKG)/pkg/util.GitBranch=$(shell git rev-parse --abbrev-ref HEAD)"
LDFLAGS += -X "$(CDC_PKG)/pkg/util.GoVersion=$(shell go version)"

default: build buildsucc

buildsucc:
	@echo Build TiDB CDC successfully!

all: dev install

test: unit_test

build: std_plugin

run_std_plugin:
	$(GO) run -trimpath ./std-plugin/host/main.go ./std-plugin/host/plugin.so

std_plugin:
	$(GOBUILD) -buildmode=plugin -o ./std-plugin/host/plugin.so ./std-plugin/plugin/sink.go

run_hashicorp_plugin:
	$(GO) run -trimpath ./hashicorp_plugin/host/main.go ./hashicorp_plugin/host/plugin

hashicorp_plugin:
	$(GOBUILD) -o ./hashicorp_plugin/host/plugin ./hashicorp_plugin/plugin/sink.go

install:
	go install ./...

unit_test: check_failpoint_ctl
	mkdir -p "$(TEST_DIR)"
	$(FAILPOINT_ENABLE)
	@export log_level=error;\
	$(GOTEST) -cover -covermode=atomic -coverprofile="$(TEST_DIR)/cov.unit.out" $(PACKAGES) \
	|| { $(FAILPOINT_DISABLE); exit 1; }
	$(FAILPOINT_DISABLE)

fmt:
	@echo "gofmt (simplify)"
	@gofmt -s -l -w $(FILES) 2>&1 | $(FAIL_ON_STDOUT)

lint:tools/bin/revive
	@echo "linting"
	@tools/bin/revive -formatter friendly -config tools/check/revive.toml $(FILES)

vet:
	@echo "vet"
	$(GO) vet $(PACKAGES) 2>&1 | $(FAIL_ON_STDOUT)

tidy:
	@echo "go mod tidy"
	./tools/check/check-tidy.sh

check: check-copyright fmt lint check-static tidy

check-static: tools/bin/golangci-lint
	$(GO) mod vendor
	tools/bin/golangci-lint \
		run ./... # $$($(PACKAGE_DIRECTORIES))

clean:
	go clean -i ./...
	rm -rf *.out

tools/bin/revive: tools/check/go.mod
	cd tools/check; \
	$(GO) build -o ../bin/revive github.com/mgechev/revive

tools/bin/golangci-lint: tools/check/go.mod
	cd tools/check; \
	$(GO) build -o ../bin/golangci-lint github.com/golangci/golangci-lint/cmd/golangci-lint

