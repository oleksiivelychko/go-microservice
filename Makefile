generate-swagger:
	rm -rf sdk/client && rm -rf sdk/models
	swagger generate spec -o ./sdk/swagger.yaml --scan-models

generate-client: generate-swagger
	$(info 'swagger generate client --help')
	swagger generate client -A go-microservice -f ./sdk/swagger.yaml -t ./sdk

install-swagger:
	which swagger || (brew tap go-swagger/go-swagger && brew install go-swagger)

install-redoc:
	npm i -g redoc-cli

start: generate-client
	HOST=localhost PORT=9090 GRPC_PORT=9091 go run main.go

docker-network:
	docker network inspect go-network >/dev/null 2>&1 || docker network create --driver bridge go-network

docker-volume:
	docker volume inspect mysql-data >/dev/null 2>&1 || docker volume create mysql-data

mysql-client-run: docker-network
	docker run -it --network go-network mysql:8.0.31 mysql -hmysql-server -uroot -p

mysql-server-run: docker-network docker-volume
	docker run --name mysql-server \
		--network gonet \
		-v mysql-data:/var/lib/mysql \
		-p 3306:3306 \
		-e MYSQL_ROOT_PASSWORD=secret \
		mysql:8.0.31

mysql-hard-stop:
	docker stop mysql-server
	docker rm mysql-server

