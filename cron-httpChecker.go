package main

import (
	"github.com/bingoohuang/go-utils"
	"github.com/metakeule/fmtdate"
	"github.com/robfig/cron"
	"time"
)

type BlackcatHttpChecker struct {
	Title          string
	Cron           string
	Url            string
	OK             string
	OkMsg          string
	FailMsg        string
	MessageTargets []string // 消息发送目标
}

func loadHttpCheckerCrons(blackcatCron *cron.Cron) {
	for _, httpChecker := range devopsConf.BlackcatHttpCheckers {
		checker := httpChecker
		_ = blackcatCron.AddFunc(httpChecker.Cron, func() {
			HttpCheck(checker)
		})
	}
}

func HttpCheck(checker BlackcatHttpChecker) {
	url := fmtdate.Format(checker.Url, time.Now())
	bytes, err := go_utils.HttpGet(url)
	if err != nil {
		AddAlertMsg(checker.MessageTargets, checker.Title, "有错误啦~\n"+err.Error())
		return
	}

	retMsg := string(bytes)
	if retMsg == checker.OK {
		AddAlertMsg(checker.MessageTargets, checker.Title, checker.OkMsg)
	} else {
		AddAlertMsg(checker.MessageTargets, checker.Title, checker.FailMsg+"\n"+retMsg)
	}
}
