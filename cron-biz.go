package main

import (
	"github.com/bingoohuang/go-utils"
	"github.com/robfig/cron"
	"time"
)

var bizCron *cron.Cron = nil

func startBizCron() {
	if redisServer == nil {
		return
	}

	if bizCron != nil {
		bizCron.Stop()
	}

	bizCron = cron.New()
	c := "0 30 20 1/1 * ?"
	//c := "@every 10s"
	_ = bizCron.AddFunc(c, func() {
		msg := createRemindMsg()
		AddAlertMsg(devopsConf.BlackcatThreshold.MessageTargets, "统计上课提醒啦~", msg)
	})
	bizCron.Start()
}

func createRemindMsg() string {
	key := "YogaClassNotificationCount:" + time.Now().Format("20060102")
	coach, _ := RedisGet(key + "coach")
	member, _ := RedisGet(key + "member")
	return "今日发送上课提醒: 教练" + go_utils.EmptyThen(coach, "0") + "条;会员" +
		go_utils.EmptyThen(member, "0") + "条!"
}
