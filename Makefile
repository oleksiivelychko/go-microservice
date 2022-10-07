generate-swagger:
	rm -rf sdk/client && rm -rf sdk/models
	swagger generate spec -o ./sdk/swagger.yaml --scan-models

generate-client: generate-swagger
	$(info 'swagger generate client --help')
	swagger generate client -A go-microservice -f ./sdk/swagger.yaml -t ./sdk

install-aws:
	sudo -S rm /usr/local/bin/aws
	curl "https://awscli.amazonaws.com/AWSCLIV2.pkg" -o "AWSCLIV2.pkg"
	sudo -S installer -pkg ./AWSCLIV2.pkg -target /
	rm AWSCLIV2.pkg
	which aws
	aws --version

install-protobuf:
	brew install protobuf
	go install google.golang.org/protobuf/cmd/protoc-gen-go@latest

install-swagger:
	which swagger || (brew tap go-swagger/go-swagger && brew install go-swagger)

install-redoc:
	npm i -g redoc-cli

.PHONY: protos
protos:
	protoc --go_out=protos --go-grpc_out=protos protos/*.proto

start: generate-client
	HOST=localhost PORT=9090 go run main.go
