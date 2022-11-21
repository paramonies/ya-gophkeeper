.PHONY: env_up
env_up:
	docker-compose up -d
	docker-compose ps
	./build/wait.sh
	sql-migrate up -env=local
	sql-migrate status -env=local

.PHONY: env_down
env_down:
	docker-compose down -v --rmi local --remove-orphans

.PHONY: proto
proto:
	buf generate --path api/gophkeeper/v1


.PHONY: bootstrap-deps
bootstrap-deps:
	cd tools && \
	go mod tidy && \
	go get -v github.com/rubenv/sql-migrate/... && \
	go install \
        github.com/grpc-ecosystem/grpc-gateway/v2/protoc-gen-grpc-gateway \
        github.com/grpc-ecosystem/grpc-gateway/v2/protoc-gen-openapiv2 \
        google.golang.org/protobuf/cmd/protoc-gen-go \
        google.golang.org/grpc/cmd/protoc-gen-go-grpc \
        github.com/envoyproxy/protoc-gen-validate

.PHONY: proto-test
proto-test:
	protoc --go_out=. \
	--go_opt=paths=source_relative \
	--go-grpc_out=. \
	--go-grpc_opt=paths=source_relative \
	./api/gophkeeper/v1/gophkeeper_service.proto

.PHONY: fmt
fmt:
	goimports -local "github.com/paramonies/ya-gophkeeper" -w cmd internal pkg/logger

.PHONY: server-run
server-run:
	go run cmd/server/main.go

.PHONY: client-register
client-register:
	go run cmd/client/main.go registerUser --login=test4 --password=123456