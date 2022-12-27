FROM golang:1.18.3-buster as builder

WORKDIR /go/src/app

RUN go install github.com/swaggo/swag/cmd/swag@latest

COPY go.mod go.sum ./

RUN go mod download

COPY . .

RUN go build -v -o app ./main.go

FROM debian:stable-slim

RUN apt-get update && \
    apt-get install -y ca-certificates apache2-utils && \
    rm -rf /var/cache/apt/*

WORKDIR /go/src/app

COPY --from=builder /go/src/app/ /go/src/app/

CMD ["./app"]
