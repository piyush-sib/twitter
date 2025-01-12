.DEFAULT_GOAL=build

GO_BUILD_DIR=build
.PHONY: build
build:
	mkdir -p $(GO_BUILD_DIR)
	go build -v -ldflags="-n -X main.version=$(VERSION)" -o $(GO_BUILD_DIR) ./cmd/...



.PHONY: lint-rules
lint-rules::
	# Mandatory files.
	[ -e .dockerignore ]
	[ -e Dockerfile ]
