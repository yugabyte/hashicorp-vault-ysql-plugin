#!/bin/sh

#   Create a directory named as test_plugin
mkdir test_plugin
cd test_plugin

#   Clone the repositories required
##  Yugabytedb plugin for HashiCorp Vault
git clone  https://github.com/yugabyte/hashicorp-vault-ysql-plugin.git  && cd hashicorp-vault-ysql-plugin

#   Run the test
export LOGDATA=$(go test)
grep -q "FAIL" <<< "$LOGDATA";
if [ $? -eq 0 ] 
then 
echo -e "Test Failed \n $LOGDATA"
else 
echo -e "Test Passed \n $LOGDATA"
fi

#   Delete the test repository
cd ../.. && rm -rf test_plugin
