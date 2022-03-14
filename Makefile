check-install:
	which swagger || (brew tap go-swagger/go-swagger && brew install go-swagger)

swagger:
	swagger generate spec -o ./swagger.yaml --scan-models

install-redoc:
	npm i -g redoc-cli

run: swagger
	PORT=9090 go run main.go
