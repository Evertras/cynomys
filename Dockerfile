# Builder image
FROM golang:1.19.3 AS builder

WORKDIR /app
COPY go.mod .
COPY go.sum .
RUN go mod download

COPY cmd/ cmd/
COPY pkg/ pkg/
RUN CGO_ENABLED=0 go build -o /cyn ./cmd/cyn/*.go

# We want to access some basic shell tools for debugging, but we want to be
# as tiny as possible...
FROM alpine:3.17.0
COPY --from=builder /cyn /usr/local/bin/cyn

ENTRYPOINT ["/usr/local/bin/cyn"]
