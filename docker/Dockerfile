# syntax=docker/dockerfile:1
FROM golang:1.16-alpine3.12 AS builder

WORKDIR /go/src/github.com/h2cone/dagagent
COPY go.mod .
COPY go.sum .
RUN go mod download -x

COPY . .

RUN go build -v

FROM alpine:3.12
RUN apk add curl

WORKDIR /root
COPY --from=builder /go/src/github.com/h2cone/dagagent/dagagent .
CMD ["./dagagent"]
