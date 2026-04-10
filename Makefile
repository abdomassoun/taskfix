.PHONY: build test lint clean install

BINARY := taskfix
VERSION := 1.0.0
BUILD_FLAGS := -ldflags="-X github.com/taskfix/taskfix/cmd.Version=$(VERSION)"

## build: compile the binary to ./taskfix
build:
	go build $(BUILD_FLAGS) -o $(BINARY) .

## install: build and install to $GOPATH/bin
install:
	go install $(BUILD_FLAGS) .

## test: run all tests
test:
	go test ./... -v

## test-cover: run tests with coverage report
test-cover:
	go test ./... -coverprofile=coverage.out
	go tool cover -html=coverage.out -o coverage.html
	@echo "Coverage report: coverage.html"

## lint: run go vet
lint:
	go vet ./...

## tidy: tidy and verify dependencies
tidy:
	go mod tidy
	go mod verify

## clean: remove build artifacts
clean:
	rm -f $(BINARY) coverage.out coverage.html

## run: build and run with example input
run: build
	./$(BINARY) "user cant login when password wrong"

## deb: build binary and package taskfix-deb into a .deb file
deb: build
	@echo "Preparing deb package directory"
	mkdir -p taskfix-deb/usr/local/bin
	mkdir -p taskfix-deb/etc/taskfix
	mkdir -p taskfix-deb/etc/taskfix/config.d
	cp -f $(BINARY) taskfix-deb/usr/local/bin/$(BINARY)
	cp -f configs/*.json taskfix-deb/etc/taskfix/config.d/

	# Ensure config file has correct name (no .json extension)
	@if [ -f taskfix-deb/etc/taskfix/config.json ]; then \
		mv taskfix-deb/etc/taskfix/config.json taskfix-deb/etc/taskfix/config; \
	fi

	# Ensure control file version matches Makefile VERSION
	version=$$(echo $(VERSION) | sed 's/^v//') && \
		sed -i -e "s/^Version:.*$$/Version: $$version/" taskfix-deb/DEBIAN/control && \
		dpkg-deb --build taskfix-deb taskfix_$${version}_amd64.deb && \
		echo "Created taskfix_$${version}_amd64.deb"

## help: print this help
help:
	@grep -E '^##' Makefile | sed 's/## //'
