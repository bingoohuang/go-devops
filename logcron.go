package main

import (
	"github.com/robfig/cron"
	"log"
	"net/rpc"
)

var c *cron.Cron = nil

func loadCrons() {
	if c != nil {
		c.Stop()
	}
	c = cron.New()

	for key, logRotate := range devopsConf.Logrotates {
		addCron(key, logRotate)
	}
	c.Start()
}

func addCron(logRotateName string, logRotate LogRotate) {
	for _, cron := range logRotate.Crons {
		c.AddFunc(cron, func() {
			logRotate.Name = logRotateName
			dealLogCron(logRotate)
		})
	}
}

func dealLogCron(logRotate LogRotate) {
	log.Println("run", logRotate)

	for _, logMachineName := range logRotate.Machines {
		_, nameAndAddress, err := parseMachineNameAndAddress(logMachineName)
		if err != "" {
			log.Println("unknown machine", err)
			continue
		}

		go executeCron(nameAndAddress, logRotate)
	}
}
func executeCron(nameAndAddress string, rotate LogRotate) {
	var reply CronCommandResult
	err := DialAndCall(nameAndAddress, func(client *rpc.Client) error {
		return client.Call("CronCommand.ExecuteCron", &CronCommandArg{
			Files:      rotate.Files,
			Type:       rotate.Type,
			Parameters: rotate.Parameters,
		}, &reply)
	})

	if err != nil {
		log.Println("executeCron error", err.Error())
	}

	log.Println("reply", reply)
}
