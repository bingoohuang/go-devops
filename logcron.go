package main

import (
	"fmt"
	"github.com/robfig/cron"
)

type LogRotate struct {
	Name       string
	Machines   []string
	Files      []string
	Crons      []string
	Type       string
	Parameters string
}

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

func dealLogCron(rotate LogRotate) {
	fmt.Println("run", rotate)

	for _, logMachineName := range rotate.Machines {
		machineName, nameAndAddress, err := parseMachineNameAndAddress(logMachineName)
		if err != "" {
			fmt.Println("unknown machine", err)
			continue
		}

		go RpcCall(machineName, nameAndAddress, &CronCommandArg{
			Files:      rotate.Files,
			Type:       rotate.Type,
			Parameters: rotate.Parameters,
		}, &CronCommandExecute{})
	}
}
