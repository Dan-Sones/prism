FROM golang:1.25.9-alpine3.22 AS builder
WORKDIR /src

COPY services/clickhouse-writer/go.mod services/clickhouse-writer/go.sum ./services/clickhouse-writer/
COPY libs/prismdbmodels ./libs/prismdbmodels
COPY libs/prismlogger ./libs/prismlogger
COPY libs/prismmicrobatcher ./libs/prismmicrobatcher

WORKDIR /src/services/clickhouse-writer
RUN go mod download

COPY services/clickhouse-writer/ ./
RUN CGO_ENABLED=0 go build -o /clickhouse-writer ./cmd

FROM alpine:3.22
COPY --from=builder /clickhouse-writer /clickhouse-writer
ENTRYPOINT ["/clickhouse-writer"]