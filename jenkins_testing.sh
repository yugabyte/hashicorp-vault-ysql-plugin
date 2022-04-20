#!/bin/sh

#   Set the log directory
if [ "$LOGFILE" == "" ] 
then
export LOGFILE=$(pwd)/hashicorp_plugin.log
fi

#   Clone if not their
git clone  https://github.com/yugabyte/hashicorp-vault-ysql-plugin.git 
if [ $? -eq 0 ]; then
    export DeleteDir=true
else
    echo Unable to Clone the Dir
fi

cd hashicorp-vault-ysql-plugin 

#   Run the test
go test github.com/yugabyte/hashicorp-vault-ysql-plugin  >> $LOGFILE

grep -q "FAIL" "$LOGFILE"; 
if [ $? -eq 0 ] 
then 
echo "Test Failed" 
echo "For more detail see the log files $LOGFILE" 
else 
echo "Test Passed"
rm $LOGFILE
fi

#   Delete Dir
if [ "$DeleteDir" == "true"  ]
then 
cd ..  && rm -rf hashicorp-vault-ysql-plugin
else 
cd ..
fi

export DeleteDir
export LOGFILE 