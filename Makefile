.PHONY: docker-build docker-build-experimentation docker-build-assignment docker-build-events docker-build-clickhouse-writer docker-build-push docker-build-push-experimentation docker-build-push-assignment docker-build-push-events docker-build-push-clickhouse-writer

docker-build: docker-build-experimentation docker-build-assignment docker-build-events docker-build-clickhouse-writer

docker-build-experimentation:
	docker build -f services/experimentation-service/Dockerfile -t ghcr.io/dan-sones/prism-experimentation-service:latest .

docker-build-assignment:
	docker build -f services/assignment-service/Dockerfile -t ghcr.io/dan-sones/prism-assignment-service:latest .

docker-build-events:
	cd services/events-service && ./mvnw package -DskipTests
	docker build -f services/events-service/Dockerfile -t ghcr.io/dan-sones/prism-events-service:latest services/events-service

docker-build-clickhouse-writer:
	docker build -f services/clickhouse-writer/Dockerfile -t ghcr.io/dan-sones/prism-clickhouse-writer:latest .

docker-build-push: docker-build-push-experimentation docker-build-push-assignment docker-build-push-events docker-build-push-clickhouse-writer

docker-build-push-experimentation:
	docker buildx build --platform linux/amd64,linux/arm64 --push -f services/experimentation-service/Dockerfile -t ghcr.io/dan-sones/prism-experimentation-service:latest .

docker-build-push-assignment:
	docker buildx build --platform linux/amd64,linux/arm64 --push -f services/assignment-service/Dockerfile -t ghcr.io/dan-sones/prism-assignment-service:latest .

docker-build-push-events:
	cd services/events-service && ./mvnw package -DskipTests
	docker buildx build --platform linux/amd64,linux/arm64 --push -f services/events-service/Dockerfile -t ghcr.io/dan-sones/prism-events-service:latest services/events-service

docker-build-push-clickhouse-writer:
	docker buildx build --platform linux/amd64,linux/arm64 --push -f services/clickhouse-writer/Dockerfile -t ghcr.io/dan-sones/prism-clickhouse-writer:latest .
