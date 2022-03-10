check-install:
	which swagger || (brew tap go-swagger/go-swagger && brew install go-swagger)

swagger: check-install
	GO111MODULE=off swagger generate spec -o ./swagger.yaml --scan-models

install-redoc:
	npm i -g redoc-cli
