PKG_NAME=terraform-provider-sqlserver
VERSION=1.0.0

default: init

apply: init
	terraform apply

plan: init
	terraform plan

destroy: init
	terraform destroy

init: build
	rm -f .terraform.lock.hcl
	terraform init

build:
	go build -o "terraform.d/plugins/local/local/sqlserver/${VERSION}/linux_amd64/${PKG_NAME}"