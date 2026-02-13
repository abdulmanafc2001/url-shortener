# variables
IMAGE_NAME = cumulus-edge

# Tagging: git commit by default, or user provided
GIT_COMMIT := $(shell git rev-parse --short HEAD)
IMAGE_TAG ?= $(GIT_COMMIT)

# Binaries
BINARY_NAME = bin/manager
GOOS = linux
GOARCH = amd64

# Show build info
info:
	@echo "Git Commit: $(GIT_COMMIT)"
	@echo "Image Tag:  $(IMAGE_TAG)"
	@echo "Images:"
	@echo "  - $(IMAGE_NAME):$(IMAGE_TAG)"

build:
	CGO_ENABLED=0 GOOS=$(GOOS) GOARCH=$(GOARCH) go build -o $(BINARY_NAME) cmd/main.go