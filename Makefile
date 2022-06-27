generate-client: generate-swagger
	$(info swagger generate client --help)
	swagger generate client -A go-microservice -f swagger.yaml -t sdk

generate-swagger:
	swagger generate spec -o ./swagger.yaml --scan-models

install-aws:
	sudo -S rm /usr/local/bin/aws
	curl "https://awscli.amazonaws.com/AWSCLIV2.pkg" -o "AWSCLIV2.pkg"
	sudo -S installer -pkg ./AWSCLIV2.pkg -target /
	rm AWSCLIV2.pkg
	which aws
	aws --version

install-swagger:
	which swagger || (brew tap go-swagger/go-swagger && brew install go-swagger)

install-redoc:
	npm i -g redoc-cli

start: generate-client
	HOST=localhost PORT=9090 go run main.go
