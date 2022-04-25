# go-microservice

### Completely ready-for-production API microservice is waiting to deploy on AWS cloud.

Test API using cURL:
```
curl -v localhost:9090/products | jq
curl -v localhost:9090/products -X POST -d '{"name":"tea","description":"The cup of tea","price":0.99,"SKU":"123-456-789"}'
curl -v localhost:9090/products/3 -X PUT -d '{"name":"ice tea","description":"The coldest cup of tea","price":0.49,"SKU":"123-456-000"}'
curl -v localhost:9090/files/1/unsplash.jpeg --output -
curl --request POST --data-binary "@public/files/unsplash.jpeg" localhost:9090/files/1/unsplash.jpeg
```

P.S. Photo by <a href="https://unsplash.com/@nate_dumlao?utm_source=unsplash&utm_medium=referral&utm_content=creditCopyText">Nathan Dumlao</a> on <a href="https://unsplash.com/s/photos/coffee?utm_source=unsplash&utm_medium=referral&utm_content=creditCopyText">Unsplash</a>
