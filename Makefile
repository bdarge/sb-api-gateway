proto:
	 mkdir -p out && protoc ./pb/*.proto  --go_out=:. --go-grpc_out=:. \
	--go_opt=Mpb/disposition.proto=github.com/bdarge/api-gateway/out/disposition \
	--go_opt=Mpb/customer.proto=github.com/bdarge/api-gateway/out/customer \
	--go_opt=Mpb/auth.proto=github.com/bdarge/api-gateway/out/auth \
	--go_opt=module=github.com/bdarge/api-gateway \
	--go-grpc_opt=Mpb/disposition.proto=github.com/bdarge/api-gateway/out/disposition \
	--go-grpc_opt=Mpb/customer.proto=github.com/bdarge/api-gateway/out/customer \
	--go-grpc_opt=Mpb/auth.proto=github.com/bdarge/api-gateway/out/auth \
	--go-grpc_opt=module=github.com/bdarge/api-gateway

doc: ## create api doc
	cd cmd; swag init --parseDependency cd -

server:
	go run cmd/main.go