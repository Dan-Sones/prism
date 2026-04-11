EXPERIMENTATION_SERVICE_DIR := ../experimentation-service
ASSIGNMENT_SERVICE_DIR := ../assignment-service
EVENTS_SERVICE_DIR := ../events-service
DATA_COOKING_SERVICE_DIR := ../data-cooking-service
BUCKET_FINDER_SERVICE_DIR := ../../tools/bucket-finder
STATS_ENGINE_DIR := ../stats-engine

EXPERIMENTATION_GEN_DIR := $(EXPERIMENTATION_SERVICE_DIR)/internal/grpc/generated
ASSIGNMENT_GEN_DIR := $(ASSIGNMENT_SERVICE_DIR)/internal/grpc/generated
EVENTS_GEN_DIR := $(EVENTS_SERVICE_DIR)/src/main/java
DATA_COOKING_GEN_DIR := $(DATA_COOKING_SERVICE_DIR)/internal/grpc/generated
BUCKET_FINDER_GEN_DIR := $(BUCKET_FINDER_SERVICE_DIR)/internal/grpc/generated
STATS_ENGINE_GEN_DIR := $(STATS_ENGINE_DIR)

.PHONY: generate clean

generate:
	@echo "Generating protos for experimentation-service..."
	@mkdir -p $(EXPERIMENTATION_GEN_DIR)
	protoc \
		--proto_path=. \
		--go_out=$(EXPERIMENTATION_GEN_DIR) \
		--go_opt=paths=source_relative \
		--go-grpc_out=$(EXPERIMENTATION_GEN_DIR) \
		--go-grpc_opt=paths=source_relative \
		experimentation_service_assignment/v1/*.proto \
		events_catalog/v1/*.proto \
		stats_engine/v1/*.proto

	@echo "Files Generated in $(EXPERIMENTATION_GEN_DIR)"
	@echo "Generating protos for assignment-service..."
	@mkdir -p $(ASSIGNMENT_GEN_DIR)
	protoc \
		--proto_path=. \
		--go_out=$(ASSIGNMENT_GEN_DIR) \
		--go_opt=paths=source_relative \
		--go-grpc_out=$(ASSIGNMENT_GEN_DIR) \
		--go-grpc_opt=paths=source_relative \
		experimentation_service_assignment/v1/*.proto \
		assignment_service/v1/*.proto
	@echo "Files Generated in $(ASSIGNMENT_GEN_DIR)"
	@echo "Generating protos for events-service..."
	protoc \
		--proto_path=. \
		--java_out=$(EVENTS_GEN_DIR) \
		--grpc-java_out=$(EVENTS_GEN_DIR) \
		events_catalog/v1/*.proto
	@echo "Files Generated in $(EVENTS_GEN_DIR)"
	@echo "Generating protos for data-cooking-service..."
	@mkdir -p $(DATA_COOKING_GEN_DIR)
	protoc \
		--proto_path=. \
		--go_out=$(DATA_COOKING_GEN_DIR) \
		--go_opt=paths=source_relative \
		--go-grpc_out=$(DATA_COOKING_GEN_DIR) \
		--go-grpc_opt=paths=source_relative \
		assignment_service/v1/*.proto
	@echo "Files Generated in $(DATA_COOKING_GEN_DIR)"
	@echo "Generating protos for bucket-finder..."
	@mkdir -p $(BUCKET_FINDER_GEN_DIR)
	protoc \
		--proto_path=. \
		--go_out=$(BUCKET_FINDER_GEN_DIR) \
		--go_opt=paths=source_relative \
		--go-grpc_out=$(BUCKET_FINDER_GEN_DIR) \
		--go-grpc_opt=paths=source_relative \
		assignment_service/v1/*.proto
	@echo "Files Generated in $(BUCKET_FINDER_GEN_DIR)"
	@echo "Generating protos for stats-engine..."
	@mkdir -p $(STATS_ENGINE_GEN_DIR)
	python-grpc-tools-protoc \
		--proto_path=. \
		--python_out=$(STATS_ENGINE_GEN_DIR) \
		--grpc_python_out=$(STATS_ENGINE_GEN_DIR) \
		stats_engine/v1/*.proto
	@echo "Files Generated in $(STATS_ENGINE_GEN_DIR)"


clean:
	rm -rf $(EXPERIMENTATION_GEN_DIR)/*
	rm -rf $(ASSIGNMENT_GEN_DIR)/*
	rm -rf $(EVENTS_GEN_DIR)/org/prism/eventsservice/grpc
	@echo "Cleaned generated files"