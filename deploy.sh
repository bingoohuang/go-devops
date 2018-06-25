#!/usr/bin/env bash

# ./deploy.sh app@hd2.gw01 or ./deploy.sh app@hb2.gw01
targetHost=$1
deployName=go-devops
fast=$2

if [ "$fast" == "fast" ]; then
    echo "jump building in fast mode"
else
    echo "rebuilding"
    rm -fr $deployName.linux.bin $deployName.linux.bin.bz2
    ./gobin.sh
    env GOOS=linux GOARCH=amd64 go build -o $deployName.linux.bin
    upx $deployName.linux.bin
fi

rsync -avz --human-readable --progress -e "ssh -p 22" ./$deployName.linux.bin $targetHost:.
#scp ./$deployName.linux.bin $targetHost:./
scp ./deploy-agent.sh $targetHost:.

ssh -tt $targetHost "bash -s" << eeooff
mkdir -p ./app/$deployName/
cd ./app/$deployName/
mv -f ~/$deployName.linux.bin .
mv -f ~/deploy-agent.sh .
ps -ef|grep $deployName|grep -v grep|awk '{print \$2}'|xargs -r kill -9
chmod +x ./deploy-agent.sh
./deploy-agent.sh $deployName app01 app/$deployName
./deploy-agent.sh $deployName app02 app/$deployName
./deploy-agent.sh $deployName app04 app/$deployName
./deploy-agent.sh $deployName cp01 app/$deployName
./deploy-agent.sh $deployName cp02 app/$deployName
./deploy-agent.sh $deployName ino01 app/$deployName
./deploy-agent.sh $deployName ino02 app/$deployName
./deploy-agent.sh $deployName smc01 app/$deployName
./deploy-agent.sh $deployName bap01 app/$deployName
./deploy-agent.sh $deployName bam01 app/$deployName
./start-go-devops.sh
exit
eeooff
