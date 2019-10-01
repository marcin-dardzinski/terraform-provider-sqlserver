PKG_NAME=terraform-provider-sqlserver

default: init

apply: init
	terraform apply

plan: init
	terraform plan

destroy: init
	terraform destroy

init: build
	terraform init

build:
	go build -o ${PKG_NAME}