#!/usr/bin/env bash

# ./deploy.sh app@hd2.gw01 or ./deploy.sh app@hb2.gw01
targetHost=$1
deployName=go-devops

rm -fr $deployName.linux.bin $deployName.linux.bin.bz2
./gobin.sh
env GOOS=linux GOARCH=amd64 go build -o $deployName.linux.bin
bzip2 $deployName.linux.bin
rsync -avz --human-readable --progress -e "ssh -p 22" ./$deployName.linux.bin.bz2 $targetHost:.
#scp ./$deployName.linux.bin.bz2 $targetHost:./
scp ./deploy-agent.sh $targetHost:.

ssh -tt $targetHost "bash -s" << eeooff
mkdir -p ./app/$deployName/
cd ./app/$deployName/
mv ~/$deployName.linux.bin.bz2 .
mv ~/deploy-agent.sh .
ps -ef|grep $deployName|grep -v grep|awk '{print \$2}'|xargs -r kill -9
chmod +x ./deploy-agent.sh
./deploy-agent.sh $deployName app01 app/$deployName
./deploy-agent.sh $deployName app02 app/$deployName
./deploy-agent.sh $deployName cp01 app/$deployName
./deploy-agent.sh $deployName cp02 app/$deployName
./deploy-agent.sh $deployName ino01 app/$deployName
./deploy-agent.sh $deployName ino02 app/$deployName
./deploy-agent.sh $deployName smc01 app/$deployName
./deploy-agent.sh $deployName bap01 app/$deployName
./deploy-agent.sh $deployName bam01 app/$deployName
rm -f $deployName.linux.bin
bzip2 -d $deployName.linux.bin.bz2
nohup ./$deployName.linux.bin -devMode -authBasic 2>&1 >> nohup.out &
exit
eeooff