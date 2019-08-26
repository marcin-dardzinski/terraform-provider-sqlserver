PKG_NAME=terraform-plugin-sqlserver

default: build

install: 

build:
	go build -o ${PKG_NAME}