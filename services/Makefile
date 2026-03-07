.PHONY: docker-build docker-build-experimentation docker-build-assignment docker-build-events

docker-build: docker-build-experimentation docker-build-assignment docker-build-events

docker-build-experimentation:
	docker build -f experimentation-service/Dockerfile -t experimentation-service .

docker-build-assignment:
	docker build -f assignment-service/Dockerfile -t assignment-service .

docker-build-events:
	cd events-service && ./mvnw package -DskipTests && docker build -t events-service .
