# go-microservice

### Completely ready-for-production API microservice waits to deploy on AWS cloud.

Test API using cURL:
```
# CRUD operations.
curl -v localhost:9090/products | jq
curl -v localhost:9090/products/1 | jq
curl -v localhost:9090/products -X POST -d '{"name":"tea","description":"The cup of tea","price":0.99,"SKU":"123-456-789"}'
curl -v localhost:9090/products/3 -X PUT -d '{"name":"ice tea","description":"The coldest cup of tea","price":0.49,"SKU":"123-456-000"}'
curl -v localhost:9090/products/3 -X DELETE

# Upload file as binary data.
curl --request POST --data-binary "@public/files/sample.png" localhost:9090/files/1/sample.png

# Upload compressed file.
curl -v localhost:9090/files/1/sample.png --compressed --output public/files/1/sample-gzip.png

# Post form data.
curl -v localhost:9090/products-form -X POST -F 'id=1' -F 'name=ice tea' -F 'description=The cup of tea' -F 'price=0.99' -F 'SKU=123-456-789' -F 'image=@public/files/sample.png'
```

🎥 Thanks <a href="https://www.youtube.com/c/NicJackson">Nic Jackson</a> for sharing his knowledge.
