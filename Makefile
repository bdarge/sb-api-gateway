THIS_FILE := $(lastword $(MAKEFILE_LIST))

proto:
	 mkdir -p out && protoc ./pb/*.proto  --go_out=:. --go-grpc_out=:. \
 	--go_opt=Mpb/model.proto=github.com/bdarge/api-gateway/out/model \
	--go_opt=Mpb/transaction.proto=github.com/bdarge/api-gateway/out/transaction \
	--go_opt=Mpb/customer.proto=github.com/bdarge/api-gateway/out/customer \
	--go_opt=Mpb/auth.proto=github.com/bdarge/api-gateway/out/auth \
	--go_opt=Mpb/profile.proto=github.com/bdarge/api-gateway/out/profile \
	--go_opt=module=github.com/bdarge/api-gateway \
	--go-grpc_opt=Mpb/model.proto=github.com/bdarge/api-gateway/out/model \
	--go-grpc_opt=Mpb/transaction.proto=github.com/bdarge/api-gateway/out/transaction \
	--go-grpc_opt=Mpb/customer.proto=github.com/bdarge/api-gateway/out/customer \
	--go-grpc_opt=Mpb/auth.proto=github.com/bdarge/api-gateway/out/auth \
	--go-grpc_opt=Mpb/profile.proto=github.com/bdarge/api-gateway/out/profile \
	--go-grpc_opt=module=github.com/bdarge/api-gateway

doc: ## create api doc
	cd cmd; swag init --parseDependency cd -

server:
	go run cmd/main.go

build:
	@$(MAKE) -f $(THIS_FILE) proto; docker build -f Dockerfile.grpc -t api-gateway --target dev . --no-cache
