SOURCES := $(shell find . -name '*.go')
BINARY := scanner-khulnasoft
IMAGE_TAG := dev
IMAGE := khulnasoft/harbor-scanner-khulnasoft:$(IMAGE_TAG)

.PHONY: build test test-integration test-component docker-build setup dev debug run

build: $(BINARY)

test: build
	go test -v -short -race -coverprofile=coverage.txt -covermode=atomic ./...

test-integration: build
	go test -count=1 -v -tags=integration ./test/integration/...

.PHONY: test-component
test-component: docker-build
	go test -count=1 -v -tags=component ./test/component/...

$(BINARY): $(SOURCES)
	GOOS=linux CGO_ENABLED=0 go build -o $(BINARY) cmd/scanner-khulnasoft/main.go

.PHONY: docker-build
docker-build: build
	docker build --no-cache -t $(IMAGE) .

lint:
	./bin/golangci-lint --build-tags component,integration run -v

setup:
	curl -sfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh| sh -s v1.21.0
	curl -sfL https://raw.githubusercontent.com/aquasecurity/trivy/main/contrib/install.sh | sh -s -- -b /usr/local/bin v0.48.3

submodule:
	git submodule update --init --recursive

dev:
	skaffold dev --tolerate-failures-until-deadline=true

debug:
	skaffold debug --tolerate-failures-until-deadline=true

run: export SCANNER_KHULNASOFT_CACHE_DIR = $(TMPDIR)harbor-scanner-khulnasoft/.cache/trivy
run: export SCANNER_KHULNASOFT_REPORTS_DIR=$(TMPDIR)harbor-scanner-khulnasoft/.cache/reports
run: export SCANNER_LOG_LEVEL=debug
run:
	@mkdir -p $(SCANNER_KHULNASOFT_CACHE_DIR) $(SCANNER_KHULNASOFT_REPORTS_DIR)
	@go run cmd/scanner-khulnasoft/main.go
