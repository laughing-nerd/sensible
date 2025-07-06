BINARY_NAME := sensible
BUILD_DIR := ./build
PLATFORMS := linux windows darwin
ARCHS := amd64 arm64
VERSION ?= dev

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
			GOOS=$$os GOARCH=$$arch go build -trimpath -ldflags="-s -w -X 'main.version=$(VERSION)'" -o $$outdir/$(BINARY_NAME)$$( [ $$os = windows ] && echo .exe ); \
		done \
	done
	@echo "Builds completed. Check the $(BUILD_DIR) directory."
