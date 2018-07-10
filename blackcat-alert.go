package main

import (
	"fmt"
	"github.com/bingoohuang/go-utils"
	"github.com/dustin/go-humanize"
	"strings"
	"time"
)

func blackcatAlertExLog(result *ExLogCommandResult) {
	for _, log := range result.ExLogs {
		key := "ex" + NextID()
		log.MachineName = result.Hostname
		WriteDbJson(exLogDb, key, log, 7*24*time.Hour)

		content := "Host: " + result.Hostname + "\nTs: " + log.Normal + "\nLogger: " + log.Logger +
			"\nLogTag: " + log.Normal + "\nFoundTs: " + result.Timestamp

		if len(log.Properties) > 0 {
			content += "\nProperties: " + go_utils.MapToString(log.Properties)
		}

		content += "\n" + linkLogId(key) + "\nEx: " + log.ExceptionNames
		SendAlertMsg("发现异常啦~", content)
	}

	if result.Error != "" {
		key := "er" + NextID()
		WriteDb(exLogDb, key, []byte(result.Error), 7*24*time.Hour)
		content := "\n" + linkLogId(key) + "\nEx: " + result.Error
		SendAlertMsg("发现错误啦~", content)
	}
}

func blackcatAlertAgent(result *AgentCommandResult) {
	key := "ag" + NextID()
	WriteDbJson(exLogDb, key, result, 7*24*time.Hour)

	content := make([]string, 0)
	if result.Error != "" {
		content = append(content, "Error: "+result.Error)
	}

	content = append(content, "Host: "+result.Hostname)
	content = append(content, linkLogId(key))

	if result.MemTotal > 0 {
		threshold := &devopsConf.BlackcatThreshold
		Load5Threshold := threshold.Load5Threshold * float64(result.Cores)
		if result.Load5 > Load5Threshold {
			content = append(content, "负载告警。Load5 "+fmt.Sprintf("%.2f", result.Load5)+"高于"+fmt.Sprintf("%.2f", Load5Threshold))
		}

		if result.MemAvailable < threshold.MemAvailThresholdSize {
			content = append(content, "内存告警。 可用"+humanize.IBytes(result.MemAvailable)+"低于"+threshold.MemAvailThreshold)
		}
		memAvailRatio := 1 - result.MemUsedPercent/100
		if memAvailRatio < threshold.MemAvailRatioThreshold {
			content = append(content, "内存告警。比例"+fmt.Sprintf("%.2f", memAvailRatio)+"低于"+fmt.Sprintf("%.2f", threshold.MemAvailRatioThreshold))
		}

		for _, du := range result.DiskUsages {
			if du.Free < threshold.DiskAvailThresholdSize {
				content = append(content, "磁盘告警。"+du.Path+"可用"+humanize.IBytes(result.MemAvailable)+"低于"+threshold.MemAvailThreshold)
			}
			availRatio := 1 - du.UsedPercent/100
			if availRatio < threshold.DiskAvailRatioThreshold {
				content = append(content, "磁盘告警。"+du.Path+"可用比例"+fmt.Sprintf("%.2f", availRatio)+"低于"+fmt.Sprintf("%.2f", threshold.DiskAvailRatioThreshold))
			}
		}
	}

	SendAlertMsg("发来警报啦~", strings.Join(content, "\n"))
}

func linkLogId(key string) string {
	threshold := &devopsConf.BlackcatThreshold
	if threshold.ExLogViewUrlPrefix == "" {
		return `LogId: ` + key
	} else {
		return `<a href="` + threshold.ExLogViewUrlPrefix + `/exlog/` + key + `">LogId</a>: ` + key
	}
}

var exLogDb = OpenDb("./exlogdb")

func SendAlertMsg(head, content string) {
	if qywxToken == "" {
		return
	}

	token := strings.Split(qywxToken, "/")
	msg := "驻" + hostname + "黑猫" + head + "\n" + content + "\nat " + time.Now().Format("01月02日15:04:05")
	go_utils.SendWxQyMsg(token[0], token[2], token[1], msg)
}
