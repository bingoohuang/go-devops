package main

import (
	"encoding/json"
	"fmt"
	"github.com/bingoohuang/go-utils"
	"github.com/dustin/go-humanize"
	"strings"
	"time"
)

func blackcatAlertExLog(exLogResult *ExLogCommandResult) {
	for _, log := range exLogResult.ExLogs {
		key := "ex" + NextID()
		fmt.Println("===exLog Id:", key)
		log.MachineName = exLogResult.GetMachineName()

		json, _ := json.Marshal(log)
		WriteDb(exLogDb, key, json, 7*24*time.Hour)

		content := "Host: " + log.MachineName + "\nTs: " + log.Normal + "\nLogger: " + log.Logger +
			"\nProperties: " + CreateKeyValuePairs(log.Properties) + "\nLogId: " + key +
			"\nEx: " + log.ExceptionNames

		SendAlertMsg("黑猫发现异常啦~", content)
	}

	if exLogResult.Error != "" {
		key := "er" + NextID()
		fmt.Println("===exLog Id:", key)

		WriteDb(exLogDb, key, []byte(exLogResult.Error), 7*24*time.Hour)
		content := "\nLogId: " + key + "\nEx: " + exLogResult.Error
		SendAlertMsg("黑猫发现错误啦~", content)
	}
}

func blackcatAlertAgent(result *AgentCommandResult) {
	key := "ag" + NextID()
	fmt.Println("===agent Id:", key)

	json, _ := json.Marshal(result)
	WriteDb(exLogDb, key, json, 7*24*time.Hour)

	content := make([]string, 0)
	if result.Error != "" {
		content = append(content, "Error: "+result.Error)
	}

	content = append(content, "Host: "+result.Hostname)
	content = append(content, "LogId: "+key)

	threshold := &devopsConf.BlackcatThreshold
	Load5Threshold := threshold.Load5Threshold * float64(result.Cores)
	if result.Load5 >= Load5Threshold {
		content = append(content, "Load5 "+fmt.Sprintf("%.2f", result.Load5)+" >= "+fmt.Sprintf("%.2f", Load5Threshold))
	}
	if result.MemAvailable <= threshold.MemAvailThresholdSize {
		content = append(content, "MemAvail "+humanize.IBytes(result.MemAvailable)+" <= "+threshold.MemAvailThreshold)
	}

	memAvailRatio := 1 - result.MemUsedPercent/100.
	if memAvailRatio <= threshold.MemAvailRatioThreshold {
		content = append(content, "MemAvailPercent "+fmt.Sprintf("%.2f", memAvailRatio)+" <= "+fmt.Sprintf("%.2f", threshold.MemAvailRatioThreshold))
	}

	SendAlertMsg("黑猫发来警报啦~", strings.Join(content, "\n"))
}

var exLogDb = OpenDb("./exlogdb")

func SendAlertMsg(head, content string) {
	if qywxToken == "" {
		return
	}

	token := strings.Split(qywxToken, "/")
	msg := head + "\n" + content + "\nat " + time.Now().Format("01月02日15:04:05")
	go_utils.SendWxQyMsg(token[0], token[2], token[1], msg)
}
