

.PHONY: migration-init
migration-init:
	goose -dir "./migrations" create "" sql

.PHONY: migration-up
migration-up:
	goose -dir "./migrations" postgres "user=user password=password dbname=postgres host=localhost port=5432" up
.PHONY: migration-down
migration-down:
	goose -dir "./migrations" postgres "user=postgres password=postgres dbname=postgres host=localhost port=5432" down

.PHONY: proto
proto:
	protoc --go_out=internal/app/service --go-grpc_out=internal/app/service --proto_path=proto proto/service.proto
