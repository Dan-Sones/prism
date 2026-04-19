.PHONY: docker-build docker-build-experimentation docker-build-assignment docker-build-events docker-build-clickhouse-writer docker-push docker-push-experimentation docker-push-assignment docker-push-events docker-push-clickhouse-writer

docker-build: docker-build-experimentation docker-build-assignment docker-build-events docker-build-clickhouse-writer

docker-build-experimentation:
	docker build -f services/experimentation-service/Dockerfile -t ghcr.io/dan-sones/experimentation-service .

docker-build-assignment:
	docker build -f services/assignment-service/Dockerfile -t ghcr.io/dan-sones/assignment-service .

docker-build-events:
	cd services/events-service && ./mvnw package -DskipTests && docker build -t ghcr.io/dan-sones/events-service .

docker-build-clickhouse-writer:
	docker build -f services/clickhouse-writer/Dockerfile -t ghcr.io/dan-sones/clickhouse-writer .

docker-push: docker-push-experimentation docker-push-assignment docker-push-events docker-push-clickhouse-writer

docker-push-experimentation:
	docker push ghcr.io/dan-sones/experimentation-service

docker-push-assignment:
	docker push ghcr.io/dan-sones/assignment-service	

docker-push-events:
	docker push ghcr.io/dan-sones/events-service

docker-push-clickhouse-writer:
	docker push ghcr.io/dan-sones/clickhouse-writer