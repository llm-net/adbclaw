BINARY_NAME := adbclaw
BUILD_DIR := ../bin
MODULE := github.com/llm-net/adbclaw
VERSION ?= $(shell git describe --tags --always --dirty 2>/dev/null || echo "dev")
GO := go
LDFLAGS := -ldflags "-X $(MODULE)/cmd.Version=$(VERSION)"

DEX_SRC := ../helper/ADBClawMonitor.java
DEX_DST := pkg/monitor/classes.dex
AUDIO_DEX_SRC := ../helper/ADBClawAudio.java
AUDIO_DEX_DST := pkg/audio/classes.dex
JAVAC := javac
ANDROID_JAR ?= $(ANDROID_HOME)/platforms/android-33/android.jar
D8 := d8

.PHONY: build test lint clean tidy dex audio-dex

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

audio-dex:
	@echo "Building audio capture DEX..."
	@mkdir -p /tmp/adbclaw-audio-dex
	$(JAVAC) --release 11 -classpath $(ANDROID_JAR) -d /tmp/adbclaw-audio-dex $(AUDIO_DEX_SRC)
	$(D8) --output /tmp/adbclaw-audio-dex /tmp/adbclaw-audio-dex/ADBClawAudio.class
	cp /tmp/adbclaw-audio-dex/classes.dex $(AUDIO_DEX_DST)
	@rm -rf /tmp/adbclaw-audio-dex
	@echo "Audio DEX built: $(AUDIO_DEX_DST)"
