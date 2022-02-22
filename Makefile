check_install:
	which swagger || (brew tap go-swagger/go-swagger && brew install go-swagger)

swagger: check_install
	GO111MODULE=off swagger generate spec -o ./swagger.yaml --scan-models
