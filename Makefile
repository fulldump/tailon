
VERSION = $(shell git describe --tags --always)
FLAGS = -ldflags "\
  -X main.VERSION=$(VERSION) \
"

test:
	go test -cover ./...

run:
	go run $(FLAGS) ./cmd/tailon/...

build:
	go build $(FLAGS) -o bin/ ./cmd/tailon/...

.PHONY: release
release: clean
	CGO_ENABLED=0 GOOS=linux   GOARCH=arm64 go build $(FLAGS) -o bin/tailon.linux.arm64 ./cmd/...
	CGO_ENABLED=0 GOOS=linux   GOARCH=amd64 go build $(FLAGS) -o bin/tailon.linux.amd64 ./cmd/...
	CGO_ENABLED=0 GOOS=windows GOARCH=arm64 go build $(FLAGS) -o bin/tailon.win.arm64.exe ./cmd/...
	CGO_ENABLED=0 GOOS=windows GOARCH=amd64 go build $(FLAGS) -o bin/tailon.win.amd64.exe ./cmd/...
	CGO_ENABLED=0 GOOS=darwin  GOARCH=arm64 go build $(FLAGS) -o bin/tailon.mac.arm64 ./cmd/...
	CGO_ENABLED=0 GOOS=darwin  GOARCH=amd64 go build $(FLAGS) -o bin/tailon.mac.amd64 ./cmd/...
	md5sum bin/* > bin/checksum

.PHONY: clean
clean:
	rm -f bin/*

.PHONY: deps
deps:
	go mod tidy -v;
	go mod download;
	go mod vendor;

.PHONY: doc
doc:
	go clean -testcache
	API_EXAMPLES_PATH="../doc/examples" go test ./api/...

.PHONY: version
version:
	@echo $(VERSION)
