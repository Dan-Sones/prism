.PHONY: docker-build docker-build-experimentation docker-build-assignment docker-build-events

docker-build: docker-build-experimentation docker-build-assignment docker-build-events

docker-build-experimentation:
	docker build -f services/experimentation-service/Dockerfile -t experimentation-service .

docker-build-assignment:
	docker build -f services/assignment-service/Dockerfile -t assignment-service .

docker-build-events:
	cd services/events-service && ./mvnw package -DskipTests && docker build -t events-service .
