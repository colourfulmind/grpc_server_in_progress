.PHONY: build_server
build_server:
	go build -o server cmd/server/main.go
	./main --config=./config/config.yaml

.PHONY: build_client
build_client:
	go build -o client cmd/client/main.go
	./client --config=./config/config.yaml

.PHONY: clean
clean:
	rm -rf server client amazing_logo.png

generate_files:
	protoc -I protos/proto protos/proto/blog/blog.proto --go_out=./protos/gen/go --go_opt=paths=source_relative --go-grpc_out=./protos/gen/go --go-grpc_opt=paths=source_relative

#migrate:
#	go run ./cmd/migrator --storage-path=./storage/sso.db --migrations-path=./migrations