.PHONY: proto-gen proto-clean fmt

proto-gen:
	protoc --go_out=. --go_opt=paths=source_relative \
	       --go-grpc_out=. --go-grpc_opt=paths=source_relative \
	       server/api/proto/assignment/v1/assignment.proto

proto-clean:
	find server/api/proto -name "*.pb.go" -type f -delete

fmt:
	go fmt ./...
