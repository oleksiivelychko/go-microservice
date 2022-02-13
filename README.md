# go-microservice

### Completely ready-for-production API microservice is waiting to deploy on AWS cloud.

```
curl -v localhost:9090/products | jq
curl -v localhost:9090/products/ -X POST -d '{"name":"tea","description":"The cup of tea"}'
curl -v localhost:9090/products/3 -X PUT -d '{"name":"ice tea","description":"The coldest cup of tea","price":0.49}'
```
