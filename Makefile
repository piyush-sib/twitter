.DEFAULT_GOAL=build

GO_BUILD_DIR=build
.PHONY: build
build:
	mkdir -p $(GO_BUILD_DIR)
	go build -v -ldflags="-n -X main.version=$(VERSION)" -o $(GO_BUILD_DIR) ./cmd/...

.PHONY: test
test:
	go test -v -cover -coverprofile=coverage.out ./...
	go tool cover -func=coverage.out -o=coverage.txt
	cat coverage.txt
	go tool cover -html=coverage.out -o=coverage.html

.PHONY: start
start:  build
	./build/twitter-backend -structured-log-file=./log -environment=testing
