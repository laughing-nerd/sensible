# Project name (binary output)
BINARY_NAME := sensible

# Output directory for builds
BUILD_DIR := ./build

# Platforms and architectures to build for
PLATFORMS := linux windows darwin
ARCHS := amd64 arm64

.PHONY: all clean

all: clean build

clean:
	@rm -rf $(BUILD_DIR)

build:
	@mkdir -p $(BUILD_DIR)
	@echo "Building for platforms: $(PLATFORMS), architectures: $(ARCHS)"
	@for os in $(PLATFORMS); do \
		for arch in $(ARCHS); do \
			outdir=$(BUILD_DIR)/$$os-$$arch; \
			mkdir -p $$outdir; \
			echo "Building for $$os/$$arch..."; \
			GOOS=$$os GOARCH=$$arch go build -trimpath -ldflags="-s -w" -o $$outdir/$(BINARY_NAME)$$( [ $$os = windows ] && echo .exe ); \
		done \
	done
	@echo "Builds completed. Check the $(BUILD_DIR) directory."

