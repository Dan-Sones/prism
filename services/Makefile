.PHONY: docker-build docker-build-admin docker-build-assignment docker-build-events

docker-build: docker-build-admin docker-build-assignment docker-build-events

docker-build-admin:
	docker build -f admin-service/Dockerfile -t admin-service .

docker-build-assignment:
	docker build -f assignment-service/Dockerfile -t assignment-service .

docker-build-events:
	cd events-service && ./mvnw package -DskipTests && docker build -t events-service .
