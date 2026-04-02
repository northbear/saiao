APP_NAME := saiao
CMD_PATH := ./cmd/saiao
BUILD_DIR := build
BINARY := $(BUILD_DIR)/$(APP_NAME)
VERSION := $(strip $(shell cat VERSION))
GIT_BRANCH := $(shell git branch --show-current 2>/dev/null)
VERSION_SUFFIX := $(if $(filter master,$(GIT_BRANCH)),,-dev)
EFFECTIVE_VERSION := $(VERSION)$(VERSION_SUFFIX)
COMMIT := $(shell git rev-parse --short HEAD 2>/dev/null || echo unknown)
LDFLAGS := -X 'saiao/internal/buildinfo.Version=$(EFFECTIVE_VERSION)' -X 'saiao/internal/buildinfo.Commit=$(COMMIT)'

DOCKER_IMAGE ?= saiao
DOCKER_TAG ?= $(EFFECTIVE_VERSION)
DOCKER_BUILD_CONTEXT ?= .

.PHONY: all test build docker-build clean

all: test build

test:
	go test ./...

build:
	mkdir -p $(BUILD_DIR)
	CGO_ENABLED=0 go build -ldflags "$(LDFLAGS)" -o $(BINARY) $(CMD_PATH)

docker-build:
	docker build --build-arg VERSION=$(VERSION) --build-arg GIT_BRANCH=$(GIT_BRANCH) --build-arg COMMIT=$(COMMIT) -t $(DOCKER_IMAGE):$(DOCKER_TAG) $(DOCKER_BUILD_CONTEXT)

clean:
	rm -rf $(BUILD_DIR)
