PKG_NAME=terraform-provider-sqlserver

default: build

install: 

build:
	go build -o ${PKG_NAME}