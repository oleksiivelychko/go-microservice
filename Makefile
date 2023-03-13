generate-spec:
	$(info 'swagger generate spec --help')
	swagger generate spec --scan-models --output=./sdk/swagger.yaml

generate-model:
	$(info 'swagger generate model --help')
	swagger generate model --spec=./sdk/swagger.yaml --target=./sdk

generate-client:
	$(info 'swagger generate client --help')
	swagger generate client -A go-microservice --spec=./sdk/swagger.yaml --target=./sdk

start:
	HOST=localhost PORT=9090 GRPC_PORT=9091 go run main.go
