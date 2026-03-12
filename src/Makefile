BINARY_NAME := adbclaw
BUILD_DIR := ../bin
MODULE := github.com/llm-net/adbclaw
VERSION ?= $(shell git describe --tags --always --dirty 2>/dev/null || echo "dev")
GO := go
LDFLAGS := -ldflags "-X $(MODULE)/cmd.Version=$(VERSION)"

.PHONY: build test lint clean tidy

build: tidy
	@mkdir -p $(BUILD_DIR)
	$(GO) build $(LDFLAGS) -o $(BUILD_DIR)/$(BINARY_NAME) .

tidy:
	$(GO) mod tidy

test:
	$(GO) test ./...

lint:
	$(GO) vet ./...

clean:
	rm -rf $(BUILD_DIR)
