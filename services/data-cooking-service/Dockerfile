FROM golang:1.25.9-alpine3.22 AS builder
WORKDIR /src

COPY services/data-cooking-service/go.mod services/data-cooking-service/go.sum ./services/data-cooking-service/
COPY libs/prismdbmodels ./libs/prismdbmodels
COPY libs/prismlogger ./libs/prismlogger
COPY libs/prismmicrobatcher ./libs/prismmicrobatcher
COPY libs/prismhash ./libs/prismhash

WORKDIR /src/services/data-cooking-service
RUN go mod download

COPY services/data-cooking-service/ ./
RUN CGO_ENABLED=0 go build -o /data-cooking-service ./cmd

FROM alpine:3.22
COPY --from=builder /data-cooking-service /data-cooking-service
ENTRYPOINT ["/clickhouse-writer"]