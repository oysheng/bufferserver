PACKAGES := $(shell go list ./... | grep -v '/vendor/')

.PHONY: all clean test

all:
	@mkdir -p target
	go build -o target/updater cmd/updater/main.go
	go build -o target/api cmd/api/main.go

test:
	@echo "Running go test"
	@go test $(PACKAGES)

intergation:
	@go run test/main.go config_travis.json

ci: test intergation

clean:
	rm -f target/updater
	rm -f target/api