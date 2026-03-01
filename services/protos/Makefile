ADMIN_SERVICE_DIR := ../admin-service
ASSIGNMENT_SERVICE_DIR := ../assignment-service
EVENTS_SERVICE_DIR := ../events-service

ADMIN_GEN_DIR := $(ADMIN_SERVICE_DIR)/internal/grpc/generated
ASSIGNMENT_GEN_DIR := $(ASSIGNMENT_SERVICE_DIR)/internal/grpc/generated
EVENTS_GEN_DIR := $(EVENTS_SERVICE_DIR)/src/main/java

.PHONY: generate clean

generate:
	@echo "Generating protos for admin-service..."
	@mkdir -p $(ADMIN_GEN_DIR)
	protoc \
		--proto_path=. \
		--go_out=$(ADMIN_GEN_DIR) \
		--go_opt=paths=source_relative \
		--go-grpc_out=$(ADMIN_GEN_DIR) \
		--go-grpc_opt=paths=source_relative \
		assignment/v1/*.proto \
		events_catalog/v1/*.proto
	@echo "Files Generated in $(ADMIN_GEN_DIR)"
	@echo "Generating protos for assignment-service..."
	@mkdir -p $(ASSIGNMENT_GEN_DIR)
	protoc \
		--proto_path=. \
		--go_out=$(ASSIGNMENT_GEN_DIR) \
		--go_opt=paths=source_relative \
		--go-grpc_out=$(ASSIGNMENT_GEN_DIR) \
		--go-grpc_opt=paths=source_relative \
		assignment/v1/*.proto
	@echo "Files Generated in $(ASSIGNMENT_GEN_DIR)"


	@echo "Generating protos for events-service..."
	protoc \
		--proto_path=. \
		--java_out=$(EVENTS_GEN_DIR) \
		--grpc-java_out=$(EVENTS_GEN_DIR) \
		events_catalog/v1/*.proto
	@echo "Files Generated in $(EVENTS_GEN_DIR)"

clean:
	rm -rf $(ADMIN_GEN_DIR)/*
	rm -rf $(ASSIGNMENT_GEN_DIR)/*
	rm -rf $(EVENTS_GEN_DIR)/org/prism/eventsservice/grpc
	@echo "Cleaned generated files"