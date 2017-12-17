#!/usr/bin/env bash

deployName=$1
targetHost=$2
targetPath=$3
scp ./$deployName.linux.bin.bz2 $targetHost:.
ssh -tt $targetHost "bash -s" << eeooff
mkdir -p $targetPath/$deployName/
mv $deployName.linux.bin.bz2 $targetPath/$deployName/
cd $targetPath/$deployName/
ps -ef|grep $deployName|grep -v grep|awk '{print \$2}'|xargs -r kill -9
rm -f $deployName.linux.bin
bzip2 -d $deployName.linux.bin.bz2
nohup ./$deployName.linux.bin &
exit
eeooff