generate-swagger:
	rm -rf sdk/client && rm -rf sdk/models
	swagger generate spec -o ./sdk/swagger.yaml --scan-models

generate-client: generate-swagger
	$(info 'swagger generate client --help')
	swagger generate client -A go-microservice -f ./sdk/swagger.yaml -t ./sdk

install-aws-cli:
	sudo -S rm /usr/local/bin/aws
	curl "https://awscli.amazonaws.com/AWSCLIV2.pkg" -o "AWSCLIV2.pkg"
	sudo -S installer -pkg AWSCLIV2.pkg -target /
	rm AWSCLIV2.pkg
	which aws
	aws --version

install-swagger:
	which swagger || (brew tap go-swagger/go-swagger && brew install go-swagger)

install-redoc:
	npm i -g redoc-cli

start: generate-client
	HOST=localhost PORT=9090 GRPC_PORT=9091 go run main.go

create-docker-network:
	docker network create -d bridge gonet

create-docker-volume:
	docker volume create mysql-data

run-mysql-server:
	docker run --name mysql-server \
		--network gonet \
		-v mysql-data:/var/lib/mysql \
		-p 3306:3306 \
		-e MYSQL_ROOT_PASSWORD=secret \
		mysql:8.0.31

stop-mysql:
	docker stop mysql-server
	docker rm mysql-server

run-mysql-client:
	docker run -it --network gonet mysql:8.0.31 mysql -hmysql-server -uroot -p
