FROM golang:1.25.9-alpine3.22 AS builder
WORKDIR /src

COPY services/experimentation-service/go.mod services/experimentation-service/go.sum ./services/experimentation-service/
COPY libs/prismdbmodels ./libs/prismdbmodels
COPY libs/prismlogger ./libs/prismlogger
COPY libs/prismhash ./libs/prismhash

WORKDIR /src/services/experimentation-service
RUN go mod download

COPY services/experimentation-service/ ./
RUN CGO_ENABLED=0 go build -o /experimentation-service ./cmd

FROM alpine:3.22
COPY --from=builder /experimentation-service /experimentation-service
ENTRYPOINT ["/experimentation-service"]