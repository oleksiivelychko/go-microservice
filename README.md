# go-microservice

### Demo API microservice covered by OpenAPI/Swagger-generated documentation. Communicates with gRPC service as well.

📌 CRUD actions:
```
curl -v localhost:9090/products | jq
curl -v "localhost:9090/products?currency=USD" | jq
curl -v "localhost:9090/products/1?currency=USD" | jq
curl -v localhost:9090/products -X POST -d '{"name":"tea","price":0.99,"SKU":"123-456-789"}' | jq
curl -v localhost:9090/products/3 -X PUT -d '{"name":"ice tea","price":0.49,"SKU":"123-456-000"}' | jq
curl -v localhost:9090/products/3 -X DELETE
```
📌 Upload file as binary data:
```
curl --request POST --data-binary "@public/files/sample.png" localhost:9090/files/1/sample.png
```
📌 Upload compressed file:
```
curl -v localhost:9090/files/1/sample.png --compressed --output public/files/1/sample-gzip.png
```
📌 Send multipart/form-data:
```
curl -v localhost:9090/products-form -X POST -F 'id=1' -F 'name=ice tea' -F 'price=0.99' -F 'SKU=123-456-789' -F 'image=@public/files/sample.png'
```
📌 OpenAPI/Swagger-generated API Documentation based on Swagger UI is available by [localhost:9090/swagger](http://localhost:9090/swagger)
![OpenAPI/Swagger-generated API Documentation based on Swagger UI](public/swagger_ui.png)

📌 OpenAPI/Swagger-generated API Documentation based on Redoc is available by [localhost:9090/redoc](http://localhost:9090/redoc)
![OpenAPI/Swagger-generated API Documentation based on Redoc UI](public/redoc_ui.png)

⚠️ [gRPC server](https://github.com/oleksiivelychko/go-grpc-service) must be running before.

⚠️ Install **swagger** locally before generate:
```
git clone https://github.com/go-swagger/go-swagger && cd go-swagger
git checkout v0.30.4
go install -ldflags "-X github.com/go-swagger/go-swagger/cmd/swagger/commands.Version=$(git describe --tags) -X github.com/go-swagger/go-swagger/cmd/swagger/commands.Commit=$(git rev-parse HEAD)" ./cmd/swagger
```

🎥 Thanks [Nic Jackson](https://www.youtube.com/c/NicJackson) for sharing his knowledge.
