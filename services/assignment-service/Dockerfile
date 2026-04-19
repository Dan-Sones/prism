FROM golang:1.25.9-alpine3.22 AS builder
WORKDIR /src

COPY services/assignment-service/go.mod services/assignment-service/go.sum ./services/assignment-service/
COPY libs/prismdbmodels ./libs/prismdbmodels
COPY libs/prismlogger ./libs/prismlogger

WORKDIR /src/services/assignment-service
RUN go mod download

COPY services/assignment-service/ ./
RUN CGO_ENABLED=0 go build -o /assignment-service ./cmd

FROM alpine:3.22
COPY --from=builder /assignment-service /assignment-service
ENTRYPOINT ["/assignment-service"]