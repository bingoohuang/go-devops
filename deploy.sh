#!/usr/bin/env bash

# ./deploy.sh app@hd2.gw01 or ./deploy.sh app@hb2.gw01
targetHost=$1
deployName=go-devops

rm -fr $deployName.linux.bin $deployName.linux.bin.bz2
env GOOS=linux GOARCH=amd64 go build -o $deployName.linux.bin
bzip2 $deployName.linux.bin
rsync -avz --human-readable --progress -e "ssh -p 22" ./$deployName.linux.bin.bz2 $targetHost:./
#scp ./$deployName.linux.bin.bz2 $targetHost:./
scp ./deploy-agent.sh $targetHost:./
ssh -tt $targetHost "bash -s" << eeooff
chmod +x ./deploy-agent.sh
./deploy-agent.sh $deployName app01 app
rm -f ./deploy-agent.sh
exit
eeooff