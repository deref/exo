FROM golang:1.17

WORKDIR /go/src/app

COPY go.mod .
RUN go mod download
COPY . .
RUN ["go", "build", "-o=main", "."]
CMD ["./main"]
