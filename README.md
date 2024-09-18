#   YSQL plugin for Hashicorp Vault 
##  About YugabyteDB:
YugabyteDB is a high-performance, cloud-native distributed SQL database that aims to support all PostgreSQL features. It is best to fit for cloud-native OLTP (i.e. real-time, business-critical) applications that need absolute data correctness and require at least one of the following: scalability, high tolerance to failures, or globally-distributed deployments.

### What makes YugabyteDB unique?
YugabyteDB is a transactional database that brings together 4 must-have needs of cloud native apps, namely SQL as a flexible query language, low-latency performance, continuous availability and globally-distributed scalability. Other databases do not serve all 4 of these needs simultaneously.

Monolithic SQL databases offer SQL and low-latency reads but neither have ability to tolerate failures nor can scale writes across multiple nodes, zones, regions and clouds.
Distributed NoSQL databases offer read performance, high availability and write scalability but give up on SQL features such as relational data modeling and ACID transactions.

Read more about YugabyteDB in our [Docs](https://docs.yugabyte.com/preview/faq/general/).

##  About HashiCorp Vault:
HashiCorp Vault is designed to help organizations manage access to secrets and transmit them safely within an organization. 
Secrets are defined as any form of sensitive credentials that need to be tightly controlled and monitored and can be used to unlock sensitive information. 
Secrets could be in the form of passwords, API keys, SSH keys, RSA tokens, or OTP.

### Dynamic Secrets:
A dynamic secret is generated on demand and is unique to a client, instead of a static secret, which is defined ahead of time and shared. 
Vault associates each dynamic secret with a lease and automatically destroys the credentials when the lease expires.
In this example, a client is requesting a database credential. Vault connects to the database with a private, root level credential and creates a new username and password. This new set of credentials are provided back to the client with a lease of 7 days. A week later, Vault will connect to the database with its privileged credentials and delete the newly created username.

Using Dynamic Secrets means we don’t have to be concerned about them having the shared PEM when a developer or operator leaves the organization. It also gives us a better break glass procedure should these credentials leak, as the credentials are localized to an individual resource reducing the attack vector, and the credentials are also issued with a time to live, meaning that Vault will automatically revoke them after a predetermined duration. In addition to this, by leveraging Vault Auth and Dynamic Secrets, you also gain full access logs directly tying a SSH session to an individual user.

![ alt text for screen readers source: HashiCorp](https://www.datocms-assets.com/2885/1519774324-dynamic-secret-img-001.jpeg?fit=max&q=80&w=2500)

###  Ysql-plugin for Hashicorp Vault:
-   ysql-plugin provides APIs for using the HashiCorp Vault's Dynamic Secrets for the yugabyteDB.
-   The APIs that can be used are as follows:  
    -   Add yugabyteDB to the manage secrets i.e. enabling `write database` for yugabyteDB(ysql) while using vault.
    -   To create new users i.e. enabling `write` roles and `read` roles commands for yugabyteDB(ysql) while using vault.
    -   Mangae lease related to the yugabyteDB(ysql) i.e. enabling `lease lookup` , `lease renew` and `lease revoke` for yugabyteDB (ysql) while using vault.
-   Why seperate plugin for yugabyteDB(ysql):
    -   YugabyteDB Go driver can be used for connecting with the database.
  This will allow us to use the added [smart features](https://docs.yugabyte.com/preview/reference/drivers/ysql-client-drivers/#yugabytedb-pgx-smart-driver), providing a high tolerance towards failures.
        

###  Before using the vault follow the below steps:
-   Make sure that the go is added to the path
```sh
export GOPATH=$HOME/go
export PATH=$PATH:$GOROOT/bin:$GOPATH/bin
```
-   Get the plugin binary:

    -   One can clone and build:
      
        -   Clone and go to the database plugin directory
          
        ```sh 
        git clone https://github.com/yugabyte/hashicorp-vault-ysql-plugin && cd hashicorp-vault-ysql-plugin  
        ```
        -   Build the plugin
          
        ```sh
        go build -o <build dir>/ysql-plugin cmd/ysql-plugin/main.go
        ```

    -    Alternatively, download the binary directly from GitHub: 
 
         -    Pre-built binary can be found at the [releases page](https://github.com/yugabyte/hashicorp-vault-ysql-plugin/releases). Download, unzip the file and place the binary `ysql-plugin` in \<build dir\>.

-   For production mode register the plugin:
```sh
export SHA256=$(sha256sum <build dir>/ysql-plugin  | cut -d' ' -f1)

vault write sys/plugins/catalog/database/ysql-plugin \
    sha256=$SHA256 \
    command="ysql-plugin"
```
-   For using the vault in the development mode add the default Vault address and Vault token
```sh
#   Add the VAULT_ADDR and VAULT_TOKEN
export VAULT_ADDR="http://localhost:8200"
export VAULT_TOKEN="root"
```

###  Using Vault

-   Running the server in the development mode
    -   For running the vault server in development mode `dev` flag is used.
    -   The `dev-root-token` informs the vault to use the default vault token of `root` to login.
        In case of production mode this token is required to be set.   
        Token policies are discussed [here](https://www.vaultproject.io/docs/commands/login).
    -   While running in the development mode vault will automatically register the plugin if 
        the directory of the binary of the plugin is provided as an input with the dev-plugin-dir flag as shown below.
```sh
vault server -dev -dev-root-token-id=root -dev-plugin-dir=<build dir> 
```

-   Enable the database's secrets:
```sh
vault secrets enable database
```
-   Add the database

One can enter the credentials:
```sh
vault write database/config/yugabytedb plugin_name=ysql-plugin  \
host="127.0.0.1" \
port=5433 \
username="yugabyte" \
password="yugabyte" \
db="yugabyte" \
load_balance=true \
yb_servers_refresh_interval=0 \
allowed_roles="*"
``` 
or use connection string:	
```sh
vault write database/config/yugabytedb \
plugin_name=ysql-plugin \
connection_url="postgres://{{username}}:{{password}}@localhost:5433/yugabyte?sslmode=disable&load_balance=true&yb_servers_refresh_interval=0" \
allowed_roles="*" \
username="yugabyte" \
password="yugabyte"
```

-   Write the role 
```sh
vault write database/roles/my-first-role \
db_name=yugabytedb \
creation_statements="CREATE ROLE \"{{username}}\" WITH PASSWORD '{{password}}' VALID UNTIL '{{expiration}}' NOINHERIT LOGIN; \
    GRANT ALL ON DATABASE \"yugabyte\" TO \"{{username}}\";" \
default_ttl="1h" \
max_ttl="24h"
```
-   Create the user 
```sh
vault read database/creds/my-first-role
```

-   Lookup the details about the lease
```sh 
vault lease lookup  <leaseid>
```
-   Renew the lease
```sh
vault lease renew   <leaseid>
```    
-   Revoke the lease
```sh
vault lease revoke  <leaseid>
```

## Configure SSL/TLS
To allow YSQL Hashicorp Vault plugin to communicate securely over SSL with YugabyteDB database, you need the root certificate (`ca.crt`) of the YugabyteDB cluster. To generate these certificates and install them while launching the cluster, follow the instructions in [Create server certificates](https://docs.yugabyte.com/preview/secure/tls-encryption/server-certificates).
Because a YugabyteDB Aeon cluster is always configured with SSL/TLS, you don't have to generate any certificate but only set the client-side SSL configuration. To fetch your root certificate, refer to [CA certificate](https://docs.yugabyte.com/preview/yugabyte-cloud/cloud-secure-clusters/cloud-authentication/#download-your-cluster-certificate).

To start a secure local YugabyteDB cluster using `yugabyted`, refer to [Create a local multi-node cluster](https://docs.yugabyte.com/preview/reference/configuration/yugabyted/#create-a-local-multi-node-cluster).

For a YugabyteDB Aeon cluster, or a local YugabyteDB cluster with SSL/TLS enabled, set the SSL-related connection parameters along with other connection information while adding the database by either of the following ways:

- Provide the connection information in DSN format:

    ```sh
    vault write database/config/yugabytedb plugin_name=ysql-plugin  \
    host="127.0.0.1" \
    port=5433 \
    username="yugabyte" \
    password="yugabyte" \
    db="yugabyte" \
    load_balance=true \
    yb_servers_refresh_interval=0 \
    sslmode="verify-full" \
    sslrootcert="path/to/.crt-file" \
    allowed_roles="*"
    ```

- Provide the connection information as a connection string:

    ```sh
    vault write database/config/yugabytedb \
    plugin_name=ysql-plugin \
    connection_url="postgres://{{username}}:{{password}}@localhost:5433/yugabyte?sslmode=verify-full&load_balance=true&yb_servers_refresh_interval=0&sslrootcert=path/to/.crt-file" \
    allowed_roles="*" \
    username="yugabyte" \
    password="yugabyte"
    ```

### SSL modes

The following table summarizes the SSL modes:

| SSL Mode | Client Driver Behavior | YugabyteDB Support |
| :------- | :--------------------- | ------------------ |
| disable  | SSL disabled | Supported
| allow    | SSL enabled only if server requires SSL connection | Supported
| prefer (default) | SSL enabled only if server requires SSL connection | Supported
| require | SSL enabled for data encryption and Server identity is not verified | Supported
| verify-ca | SSL enabled for data encryption and Server CA is verified | Supported
| verify-full | SSL enabled for data encryption. Both CA and hostname of the certificate are verified | Supported

YugabyteDB Aeon requires SSL/TLS, and connections using SSL mode `disable` will fail.

## Apart from Dynamic roles ysql-plugin also supports [Static roles](https://developer.hashicorp.com/vault/tutorials/db-credentials/database-creds-rotation), [Root credential rotation](https://developer.hashicorp.com/vault/tutorials/db-credentials/database-root-rotation) and [Username customization](https://developer.hashicorp.com/vault/tutorials/secrets-management/username-templating).

## Known issues

When executing vault operations, the internal query may fail with the following error:

```output
ERROR: The catalog snapshot used for this transaction has been invalidated: expected: 2, got: 1: MISMATCHED_SCHEMA (SQLSTATE 40001)
```

A DML query in YSQL may touch multiple servers, and each server has a Catalog Version which is used to track schema changes. When a DDL statement runs in the middle of the DML query, the Catalog Version is changed and the query has a mismatch, causing it to fail.

For such cases, the database aborts the query and returns a 40001 error code. Operations failing with this code can be safely retried.

For more information, refer to [How to troubleshoot Schema or Catalog version mismatch database errors](https://support.yugabyte.com/hc/en-us/articles/4406287763597-How-to-troubleshoot-Schema-or-Catalog-version-mismatch-database-errors).

##  For testing:
go test can be used for testing the ysql-plugin
Use:: `go test github.com/yugabyte/hashicorp-vault-ysql-plugin`
For individual cases
-   For Initialize:
    `go test -run ^TestYsql_Initialize$ github.com/yugabyte/hashicorp-vault-ysql-plugin`
-   For Create User:
    `go test -run ^TestYsql_NewUser$ github.com/yugabyte/hashicorp-vault-ysql-plugin`
-   For Update User Password:
    `go test -run ^TestUpdateUser_Password$ github.com/yugabyte/hashicorp-vault-ysql-plugin`
-   For Update User Expiration:
    `go test -run ^TestUpdateUser_Expiration$ github.com/yugabyte/hashicorp-vault-ysql-plugin`
-   For Delete User:
    `go test -run ^TestDeleteUser$ github.com/yugabyte/hashicorp-vault-ysql-plugin`

## How to use the Makefile:
-   Set the BUILD_DIR in the Makefile
-   For building the plugin, registering it and running it in development mode use `make`.
-   For enabling the plugin and creating a basic role named 'my-first-role' use `make enable`.
-   To read user  use `vault read database/creds/my-first-role`.
-   Use `make clean` to remove the build and `make test` to test the plugin. 
