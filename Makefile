HOST_NAME := localhost
HOST_PORT := 9090
HOST_ADDR := $(HOST_NAME):$(HOST_PORT)

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
	HOST=$(HOST_NAME) PORT=$(HOST_PORT) PORT_GRPC=9091 go run main.go

curl-list:
	curl -v $(HOST_ADDR)/products | jq

curl-list-query:
	curl -v "$(HOST_ADDR)/products?currency=USD" | jq

curl-get:
	curl -v "$(HOST_ADDR)/products/1?currency=USD" | jq

curl-post:
	curl -v $(HOST_ADDR)/products -X POST -d '{"name":"tea","price":0.99,"SKU":"123-456-789"}' | jq

curl-put:
	curl -v $(HOST_ADDR)/products/3 -X PUT -d '{"name":"ice tea","price":0.49,"SKU":"123-456-000"}' | jq

curl-delete:
	curl -v $(HOST_ADDR)/products/3 -X DELETE

curl-upload-file-binary:
	curl --request POST --data-binary "@public/files/sample.png" $(HOST_ADDR)/files/1/sample.png

curl-upload-file-compressed:
	curl -v $(HOST_ADDR)/files/1/sample.png --compressed --output public/files/1/sample-gzip.png

curl-multipart-form-data:
	curl -v $(HOST_ADDR)/products-form -X POST -F 'id=1' -F 'name=ice tea' -F 'price=0.99' -F 'SKU=123-456-789' -F 'image=@public/files/sample.png'
