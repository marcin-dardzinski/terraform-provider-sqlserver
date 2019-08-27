PKG_NAME=terraform-provider-sqlserver

default: build

init: build
	terraform init

build:
	go build -o ${PKG_NAME}