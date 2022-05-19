setup: # Setup dependencies
	go mod tidy
	go generate ./...
.PHONY: setup

format: # Run gofmt
	go fmt ./...
.PHONY: format

lint: # Run linter
	golangci-lint run ./...
.PHONY: lint

test: # Test uses race and coverage
	go clean -testcache && go test -race -coverprofile=coverage.out -covermode=atomic
.PHONY: test

test-v: # Test with -v
	go clean -testcache && go test -race -v $$(go list ./... | $(excluded)) -coverprofile=coverage.out -covermode=atomic
.PHONY: test-v

cover: test # Run all the tests and opens the coverage report
	go tool cover -html=coverage.out
.PHONY: cover

bench: # Runs benchmarks
	go test -benchmem -bench .
.PHONY: bench

mock: # Make mocks keeping directory tree
	rm -rf mocks \
	&& mockery --name=Encoder --recursive --exported=true --output=./mocks \
	&& mockery --name=RedisStore --recursive --exported=true --output=./mocks
.PHONY: mock

doc: # Run go doc
	godoc -http localhost:8080
.PHONY: doc

all: # Make format, lint and test
	$(MAKE) format
	$(MAKE) lint
	$(MAKE) test
.PHONY: all

todo: # Show to-do items per file
	$(Q) grep \
		--exclude=Makefile.util \
		--exclude-dir=vendor \
		--exclude-dir=.idea \
		--text \
		--color \
		-nRo \
		-E '\S*[^\.]TODO.*' \
		.
.PHONY: todo

help: # Display this help
	$(Q) awk 'BEGIN {FS = ":.*#"; printf "Usage: make \033[36m<target>\033[0m\n"} /^[a-zA-Z_-]+:.*?#/ { printf "  \033[36m%-15s\033[0m %s\n", $$1, $$2 } /^##@/ { printf "\n\033[1m%s\033[0m\n", substr($$0, 5) } ' $(MAKEFILE_LIST)
.PHONY: help
