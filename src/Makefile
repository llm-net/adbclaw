BINARY_NAME := adbclaw
BUILD_DIR := ../bin
MODULE := github.com/llm-net/adbclaw
VERSION ?= $(shell git describe --tags --always --dirty 2>/dev/null || echo "dev")
GO := go
LDFLAGS := -ldflags "-X $(MODULE)/cmd.Version=$(VERSION)"

DEX_SRC := ../helper/ADBClawMonitor.java
DEX_DST := pkg/monitor/classes.dex
JAVAC := javac
ANDROID_JAR ?= $(ANDROID_HOME)/platforms/android-33/android.jar
D8 := d8

.PHONY: build test lint clean tidy dex

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

dex:
	@echo "Building monitor DEX..."
	@mkdir -p /tmp/adbclaw-dex
	$(JAVAC) --release 11 -classpath $(ANDROID_JAR) -d /tmp/adbclaw-dex $(DEX_SRC)
	$(D8) --output /tmp/adbclaw-dex /tmp/adbclaw-dex/ADBClawMonitor.class
	cp /tmp/adbclaw-dex/classes.dex $(DEX_DST)
	@rm -rf /tmp/adbclaw-dex
	@echo "DEX built: $(DEX_DST)"
