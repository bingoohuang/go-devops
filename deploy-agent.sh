#!/usr/bin/env bash

deployName=$1
targetHost=$2
targetPath=$3
rpcPort=${4:-6979}
startHttp=${5:-true}
scp ./$deployName.linux.bin.bz2 $targetHost:.
ssh -tt $targetHost "bash -s" << eeooff
mkdir -p $targetPath/
mv $deployName.linux.bin.bz2 $targetPath/
cd $targetPath/
ps -ef|grep $deployName|grep -v grep|awk '{print \$2}'|xargs -r kill -9
rm -f $deployName.linux.bin
bzip2 -d $deployName.linux.bin.bz2
nohup ./$deployName.linux.bin -rpcPort=$rpcPort -startHttp=$startHttp 2>&1 >> nohup.out &
exit
eeooff