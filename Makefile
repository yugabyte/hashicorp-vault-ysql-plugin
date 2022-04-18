GOARCH = amd64

BUILD_DIR = /home/jayantanand/code/work/hashicorp/plugin_bin

UNAME = $(shell uname -s)

ifndef OS
	ifeq ($(UNAME), Linux)
		OS = linux
	else ifeq ($(UNAME), Darwin)
		OS = darwin
	endif
endif




.DEFAULT_GOAL := all


all: fmt build start-dev

build:
	GOOS=$(OS) GOARCH="$(GOARCH)" go build -o $(BUILD_DIR)/ysql-plugin cmd/ysql-plugin/main.go

start-dev:
	vault server -dev -dev-root-token-id=root -dev-plugin-dir=$(BUILD_DIR)

enable:
	vault secrets enable database

	vault write database/config/yugabytedb \
    plugin_name=ysql-plugin  host="127.0.0.1" \
    port=5433 \
    username="yugabyte" \
    password="yugabyte" \
    db="yugabyte" \
    allowed_roles="*"

	vault write database/roles/my-first-role \
    db_name=yugabytedb \
    creation_statements="CREATE ROLE \"{{username}}\" WITH PASSWORD '{{password}}' NOINHERIT LOGIN; \
       GRANT ALL ON DATABASE \"yugabyte\" TO \"{{username}}\";" \
    default_ttl="1h" \
    max_ttl="24h"


clean:
	rm -f $(BUILD_DIR)/ysql-plugin

fmt:
	go fmt $$(go list ./...)

.PHONY: build clean fmt start-dev enable write