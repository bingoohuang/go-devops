package main

import (
	"github.com/robfig/cron"
	"log"
	"regexp"
)

var cronRegexp = regexp.MustCompile(`(?i)Every\s+(\d+)\s+(Second|Minute|Hour|Day|Monday|Tuesday|Wednesday|Thursday|Friday|Saturday|Sunday|Weeks)s?(\s+at\s+(\d\d:\d\d))?`)

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


}
