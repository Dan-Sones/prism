.PHONY: fmt docker-build

fmt:
	go fmt ./...

docker-build:
	cd .. && docker build -f admin-service/Dockerfile -t admin-service .

