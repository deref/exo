FROM golang:1.17 AS builder
LABEL org.opencontainers.image.source https://github.com/deref/exo

WORKDIR /go/src/app

COPY go.mod go.sum ./
RUN go mod download
RUN go mod graph | awk '{if ($1 !~ "@") print $2}' | xargs go get
COPY . .
RUN ["go", "build", "./cmd/mitm"]

FROM mitmproxy/mitmproxy:latest
COPY --from=builder /go/src/app/mitm .
CMD ["./mitm"]
