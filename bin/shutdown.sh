#!/bin/bash

cd `dirname $0`
target_dir=`pwd`

pid=`ps ax | grep -i opengateway.opengateway | grep ${target_dir}  | grep -v grep | awk '{print $1}'`
if [ -z "$pid" ] ; then
        echo "No open-gateway-service running."
        exit -1;
fi

echo "The open-gateway-service (${pid}) is running..."

kill -9 ${pid}

echo "Send shutdown request to open-gateway-service(${pid}) OK"
