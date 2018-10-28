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

	bizCron = cron.New()
	cron := "0 30 20 1/1 * ?" // "@every 10s"
	bizCron.AddFunc(cron, func() {
		msg := createRemindMsg()
		AddAlertMsg("统计上课提醒啦~", msg)
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
