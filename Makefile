APP_NAME := saiao
CMD_PATH := ./cmd/saiao
BUILD_DIR := build
BINARY := $(BUILD_DIR)/$(APP_NAME)

DOCKER_IMAGE ?= saiao
DOCKER_TAG ?= latest
DOCKER_BUILD_CONTEXT ?= .

.PHONY: all test build docker-build clean

all: test build

test:
	go test ./...

build:
	mkdir -p $(BUILD_DIR)
	CGO_ENABLED=0 go build -o $(BINARY) $(CMD_PATH)

docker-build:
	docker build -t $(DOCKER_IMAGE):$(DOCKER_TAG) $(DOCKER_BUILD_CONTEXT)

clean:
	rm -rf $(BUILD_DIR)
