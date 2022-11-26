# Build for local
bin/cyn: ./cmd/cyn/*.go
	go build -o bin/cyn ./cmd/cyn/*.go

.PHONY: docker
docker:
	docker build -t evertras/cyn:latest .

.PHONY: test
test:
	go test ./...

.PHONY: fmt
fmt:
	go fmt ./...
