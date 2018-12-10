#!/usr/bin/env bash

deployName=$1
targetHost=$2
targetPath=$3
rpcPort=${4:-6979}
startHttp=${5:-false}
scp ./$deployName.linux.bin $targetHost:.
ssh -tt $targetHost "bash -s" << eeooff
mkdir -p $targetPath/
mv -f $deployName.linux.bin $targetPath/
cd $targetPath/
ps -ef|grep $deployName|grep -v grep|awk '{print \$2}'|xargs -r kill -9
nohup ./$deployName.linux.bin 2>&1 >> nohup.out &
exit
eeooff