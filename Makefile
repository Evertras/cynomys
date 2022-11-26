.PHONY: default
default:
	go test ./...

.PHONY: fmt
fmt:
	go fmt ./...
