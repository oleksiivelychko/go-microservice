check-install:
	which swagger || (brew tap go-swagger/go-swagger && brew install go-swagger)

swagger:
	swagger generate spec -o ./swagger.yaml --scan-models

install-redoc:
	npm i -g redoc-cli

start: gen-client
	HOST=localhost PORT=9090 go run main.go

gen-client: swagger
	$(info swagger generate client --help)
	swagger generate client -A go-microservice -f swagger.yaml -t sdk
