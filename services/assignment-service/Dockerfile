FROM golang:1.25

WORKDIR /app

COPY assignment-service/go.mod assignment-service/go.sum ./
COPY libs/prismdbmodels /libs/prismdbmodels
COPY libs/prismlogger /libs/prismlogger

RUN go mod download

COPY assignment-service/ .

RUN CGO_ENABLED=0 GOOS=linux go build -o assignment-service ./cmd/main.go

EXPOSE 8080

CMD ["./assignment-service"]
