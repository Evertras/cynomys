.PHONY: build
build: bin/cyn

.PHONY: build-all
build-all: bin/cyn bin/cyn-linux bin/cyn-mac bin/cyn-windows

# Build for local
bin/cyn: ./cmd/cyn/*.go ./pkg/listener/*.go ./pkg/sender/*.go
	go build -o bin/cyn ./cmd/cyn/*.go

# Build for other OSes
bin/cyn-linux: ./cmd/cyn/*.go ./pkg/listener/*.go ./pkg/sender/*.go
	CGO_ENABLED=0 GOOS=linux go build -o bin/cyn-linux ./cmd/cyn/*.go
bin/cyn-mac: ./cmd/cyn/*.go ./pkg/listener/*.go ./pkg/sender/*.go
	CGO_ENABLED=0 GOOS=darwin go build -o bin/cyn-mac ./cmd/cyn/*.go
bin/cyn-windows: ./cmd/cyn/*.go ./pkg/listener/*.go ./pkg/sender/*.go
	CGO_ENABLED=0 GOOS=windows go build -o bin/cyn-windows ./cmd/cyn/*.go

.PHONY: docker
docker:
	docker build -t evertras/cyn:latest .

.PHONY: test
test:
	go test ./pkg/...

.PHONY: bdd
bdd: bin/cyn
	go test -race -v ./tests/...

.PHONY: fmt
fmt:
	go fmt ./...
