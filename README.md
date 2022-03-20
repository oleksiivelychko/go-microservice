# go-microservice

### Completely ready-for-production API microservice is waiting to deploy on AWS cloud.

Test API using cURL:
```
curl -v localhost:9090/products | jq
curl -v localhost:9090/products -X POST -d '{"name":"tea","description":"The cup of tea","price":0.99,"SKU":"123-456-789"}'
curl -v localhost:9090/products/3 -X PUT -d '{"name":"ice tea","description":"The coldest cup of tea","price":0.49,"SKU":"123-456-000"}'
```
