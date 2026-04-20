.PHONY: docker-build docker-build-experimentation docker-build-assignment docker-build-events docker-build-clickhouse-writer docker-push docker-push-experimentation docker-push-assignment docker-push-events docker-push-clickhouse-writer

docker-build: docker-build-experimentation docker-build-assignment docker-build-events docker-build-clickhouse-writer

docker-build-experimentation:
	docker buildx build --platform linux/amd64,linux/arm64 -f services/experimentation-service/Dockerfile -t ghcr.io/dan-sones/prism-experimentation-service:latest .

docker-build-assignment:
	docker buildx build --platform linux/amd64,linux/arm64 -f services/assignment-service/Dockerfile -t ghcr.io/dan-sones/prism-assignment-service:latest .

docker-build-events:
	cd services/events-service && ./mvnw package -DskipTests
	docker buildx build --platform linux/amd64,linux/arm64 -f services/events-service/Dockerfile -t ghcr.io/dan-sones/prism-events-service:latest services/events-service

docker-build-clickhouse-writer:
	docker buildx build --platform linux/amd64,linux/arm64 -f services/clickhouse-writer/Dockerfile -t ghcr.io/dan-sones/prism-clickhouse-writer:latest .

docker-push: docker-push-experimentation docker-push-assignment docker-push-events docker-push-clickhouse-writer

docker-push-experimentation:
	docker push ghcr.io/dan-sones/prism-experimentation-service:latest

docker-push-assignment:
	docker push ghcr.io/dan-sones/prism-assignment-service:latest

docker-push-events:
	docker push ghcr.io/dan-sones/prism-events-service:latest

docker-push-clickhouse-writer:
	docker push ghcr.io/dan-sones/prism-clickhouse-writer:latest
