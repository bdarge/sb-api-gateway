proto:
	protoc pkg/**/pb/*.proto --go_out=:. --go-grpc_out=:.

doc: ## create api doc
	cd cmd; swag init --parseDependency cd -

server:
	go run cmd/main.go