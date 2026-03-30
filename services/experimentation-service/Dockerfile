FROM golang:1.25

WORKDIR /app

COPY experimentation-service/go.mod experimentation-service/go.sum ./
COPY libs/prismdbmodels /libs/prismdbmodels
COPY libs/prismlogger /libs/prismlogger

RUN go mod download

COPY experimentation-service/ .

RUN CGO_ENABLED=0 GOOS=linux go build -o experimentation-service ./cmd/main.go

EXPOSE 8080

CMD ["./experimentation-service"]
