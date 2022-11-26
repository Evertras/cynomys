# Build for local
bin/cyn: ./cmd/cyn/*.go ./pkg/listener/*.go
	go build -o bin/cyn ./cmd/cyn/*.go

.PHONY: docker
docker:
	docker build -t evertras/cyn:latest .

.PHONY: test
test:
	go test ./pkg/...

.PHONY: bdd
bdd: bin/cyn
	go test ./tests/...

.PHONY: fmt
fmt:
	go fmt ./...
