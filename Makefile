BUILD_FLAGS := -ldflags "-s -w"
OUTPUT_PATH := ./bin/

.PHONY: default
default: build

.PHONY: build
build:
	go build -o $(OUTPUT_PATH) $(BUILD_FLAGS) ./cmd/...
	@echo "Build done"

.PHONY: clean
clean:
	rm -rf bin

.PHONY: test
test:
	go test
