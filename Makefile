export GO111MODULE=on
export CGO_ENABLED=0

GO_CMD=go
BUILD_OPTS=-v -mod=vendor
PKG_NAME=github.com/nokrPOC
BUILD_PREFIX=hermes-
API_TARGET_NAME=worker

build-api:
	$(GO_CMD) build $(BUILD_OPTS) -o ./bin/$(BUILD_PREFIX)$(API_TARGET_NAME) $(PKG_NAME)/cmd/$(API_TARGET_NAME)

run-api: build-api
	cd bin && ./$(BUILD_PREFIX)$(API_TARGET_NAME) && cd ..

vendor: clean
	$(GO_CMD)  mod vendor

clean:
	rm -rf vendor