SOURCES := $(shell find . -name '*.go')
BINARY := scanner-adapter
IMAGE_TAG := dev
IMAGE := khulnasoft/harbor-scanner-khulnasoft:$(IMAGE_TAG)

build: $(BINARY)

test: build
	GO111MODULE=on go test -v -short -race -coverprofile=coverage.txt -covermode=atomic ./...

test-integration: build
	GO111MODULE=on go test -count=1 -v -tags=integration ./test/integration/...

$(BINARY): $(SOURCES)
	GOOS=linux GO111MODULE=on CGO_ENABLED=0 go build -o $(BINARY) cmd/scanner-adapter/main.go

docker-build: build
	docker build --no-cache -t $(IMAGE) .
