.PHONY: fmt docker-build

fmt:
	go fmt ./...

docker-build:
	cd .. && docker build -f experimentation-service/Dockerfile -t experimentation-service .

