package main

import (
	"fmt"
	"github.com/bingoohuang/go-utils"
	"github.com/dustin/go-humanize"
	"github.com/patrickmn/go-cache"
	"strings"
	"sync"
	"time"
)

var exLogCache *cache.Cache
var exLogCacheMux sync.Mutex

var erLogCache *cache.Cache
var erLogCacheMux sync.Mutex

func init() {
	exLogCache = cache.New(10*time.Minute, 30*time.Second)
	erLogCache = cache.New(10*time.Minute, 30*time.Second)
}

func blackcatAlertExLog(result *ExLogCommandResult) {
	blackcatAlertExLogMsg(result)
	blackcatAlertErLogMsg(result)
}

func blackcatAlertExLogMsg(result *ExLogCommandResult) {
	exLogCacheMux.Lock()
	defer exLogCacheMux.Unlock()

	interval := devopsConf.BlackcatThreshold.ExLogsCollapseInterval
	for _, log := range result.ExLogs {
		key := "ex" + NextID()
		hostname := result.Hostname
		log.MachineName = hostname
		WriteDbJson(exLogDb, key, log, 7*24*time.Hour)

		logger := log.Logger
		content := "Host: " + hostname + "\nTs: " + log.Normal + "\nLogger: " + logger +
			"\nLogTag: " + log.Normal + "\nFoundTs: " + result.Timestamp

		if len(log.Properties) > 0 {
			content += "\nProperties: " + go_utils.MapToString(log.Properties)
		}

		exNames := log.ExceptionNames
		cacheKey := hostname + "+" + logger + "+" + exNames
		_, found := exLogCache.Get(cacheKey)
		if found { continue }
		erLogCache.Set(cacheKey, "", time.Duration(interval)*time.Minute)
		content += "\n" + linkLogId(key) + "\nEx: " + exNames
		AddAlertMsg(log.MessageTargets, "发现异常啦~", content)
	}
}

func blackcatAlertErLogMsg(result *ExLogCommandResult) {
	erLogCacheMux.Lock()
	defer erLogCacheMux.Unlock()

	interval := devopsConf.BlackcatThreshold.ExLogsCollapseInterval
	if result.Error != "" {
		key := "er" + NextID()
		er := result.Error
		WriteDb(exLogDb, key, []byte(er), 7*24*time.Hour)

		_, found := erLogCache.Get(er)
		if found { return }
		erLogCache.Set(er, "", time.Duration(interval)*time.Minute)
		content := "\n" + linkLogId(key) + "\nEx: " + er
		AddAlertMsg(devopsConf.BlackcatThreshold.MessageTargets, "发现错误啦~", content)
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
			content = append(content, "负载告警。Load5 "+fmt.Sprintf("%.2f", result.Load5)+
				"高于"+fmt.Sprintf("%.2f", Load5Threshold))
		}

		if result.MemAvailable < threshold.MemAvailThresholdSize {
			content = append(content, "内存告警。 可用"+humanize.IBytes(result.MemAvailable)+
				"低于"+threshold.MemAvailThreshold)
		}
		memAvailRatio := 1 - result.MemUsedPercent/100
		if memAvailRatio < threshold.MemAvailRatioThreshold {
			content = append(content, "内存告警。比例"+fmt.Sprintf("%.2f", memAvailRatio)+
				"低于"+fmt.Sprintf("%.2f", threshold.MemAvailRatioThreshold))
		}

		for _, du := range result.DiskUsages {
			if du.Free < threshold.DiskAvailThresholdSize {
				content = append(content, "磁盘告警。"+du.Path+"可用"+humanize.IBytes(result.MemAvailable)+
					"低于"+threshold.MemAvailThreshold)
			}
			availRatio := 1 - du.UsedPercent/100
			if availRatio < threshold.DiskAvailRatioThreshold {
				content = append(content, "磁盘告警。"+du.Path+"可用比例"+fmt.Sprintf("%.2f", availRatio)+
					"低于"+fmt.Sprintf("%.2f", threshold.DiskAvailRatioThreshold))
			}
		}
	}

	AddAlertMsg(devopsConf.BlackcatThreshold.MessageTargets, "发来警报啦~", strings.Join(content, "\n"))
}

func linkLogId(key string) string {
	threshold := &devopsConf.BlackcatThreshold
	if threshold.ExLogViewUrlPrefix == "" {
		return `LogId: ` + key
	}

	return `<a href="` + threshold.ExLogViewUrlPrefix + `/exlog/` + key + `">LogId: ` + key + `</a>`
}

var exLogDb = OpenDb("./exlogdb")
