FROM golang:1.25

WORKDIR /app

COPY admin-service/go.mod admin-service/go.sum ./
COPY prismdbmodels /prismdbmodels
COPY prismlogger /prismlogger

RUN go mod download

COPY admin-service/ .

RUN CGO_ENABLED=0 GOOS=linux go build -o admin-service ./cmd/main.go

EXPOSE 8080

CMD ["./admin-service"]
