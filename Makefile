APP := coffee-machine
CMD := ./cmd/coffee-machine
BUILD_DIR := build

.PHONY: fmt test vet build run tidy clean

fmt:
	gofmt -w ./cmd ./internal

test:
	go test ./...

vet:
	go vet ./...

build:
	mkdir -p $(BUILD_DIR)
	go build -o $(BUILD_DIR)/$(APP) $(CMD)

run:
	go run $(CMD)

tidy:
	go mod tidy

clean:
	rm -rf $(BUILD_DIR)