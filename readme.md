#   Ysql plugin for Hashicorp Vault's Dynamic Secrets: 

##  Steps to be followed to use the terminal

Admin's terminal to configure the database
```sh
#   Make sure that the go is added to the path
   export GOPATH=$HOME/go
   export PATH=$PATH:$GOROOT/bin:$GOPATH/bin

#   Clone and go to the database plugin directory
$   git clone https://github.com/yugabyte/hashicorp-vault-ysql-plugin

$   go build -o <build dir>/ysql-plugin  cmd/ysql-plugin/main.go

#   Add the VAULT_ADDR and VAULT_TOKEN
  export VAULT_ADDR="http://localhost:8200"
   export VAULT_TOKEN="root"

```

Run the vault server
```sh
#   Run the server 
$   vault server -dev -dev-root-token-id=root -dev-plugin-dir=<build dir> 

```

Register the plugin , config the database and create the role 
```sh
#   Register the plugin
$ export SHA256=$(sha256sum <build dir>//ysql-plugin  | cut -d' ' -f1)


$ vault secrets enable database

$ vault write sys/plugins/catalog/database/ysql-plugin \
    sha256=$SHA256 \
    command="ysql-plugin"

#   Add the database
$ vault write database/config/yugabytedb plugin_name=ysql-plugin  \
    host="127.0.0.1" \
    port=5433 \
    username="yugabyte" \
    password="yugabyte" \
    db="yugabyte" \
    allowed_roles="*"

#   Create the role
$ vault write database/roles/my-first-role \
    db_name=yugabytedb \
    creation_statements="CREATE ROLE \"{{username}}\" WITH PASSWORD '{{password}}' NOINHERIT LOGIN; \
       GRANT ALL ON DATABASE \"yugabyte\" TO \"{{username}}\";" \
    default_ttl="1h" \
    max_ttl="24h"

#   For managing the lease
$   vault lease lookup database/creds/my-first-role/MML1XWMjcJKXBlk47HHs6HrZ

$   vault lease renew  database/creds/my-first-role/MML1XWMjcJKXBlk47HHs6HrZ

$   vault lease revoke   database/creds/my-first-role/E8cCdoKTn9mvQjQAWd5aZohQ
```


-   Client/App code
Create the user 
```sh
$   vault read database/creds/my-first-role
```
docker  exec -it <docker id>  bash

##  Completion matrix
|API/TASK|Status|
|-|-|
| Initialize the plugin|✅|
| Create User |✅ |
| Delete User|✅|
| Update User|✅|
| Make File| |
| Create User -test| |
| Delete User -test| |
| Update User -test| |
| Blog| |
| Add the smart driver's feature|   |


##  Error with the revoke statement::
-   `failed to revoke lease: lease_id=database/creds/my-first-role/MML1XWMjcJKXBlk47HHs6HrZ error="failed to revoke entry: resp: (*logical.Response)(nil) err: unable to delete user: rpc error: code = Internal desc = unable to delete user: pq: role \"V_TOKEN_MY-FIRST-ROLE_HDZVDJXNAEYNDNWVW2IU_1649353280\" cannot be dropped because some objects depend on it"`



rm  /home/jayantanand/code/work/hashicorp/plugin_bin/ysql-plugin