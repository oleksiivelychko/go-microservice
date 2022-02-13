# go-microservice

### Completely ready-for-production API microservice is waiting to deploy on AWS cloud.

```
curl -v localhost:9090 | jq
curl -v localhost:9090 -d '{"name":"tea","description":"The cup of tea"}'
curl -v localhost:9090/3 -XPUT -d '{"name":"ice tea","description":"The coldest cup of tea","price":0.49}'
```
