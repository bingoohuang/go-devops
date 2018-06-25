package main

import (
	"github.com/bingoohuang/go-utils"
	"github.com/dustin/go-humanize"
	"github.com/robfig/cron"
	"strings"
	"time"
)

type BlackcatThreshold struct {
	Load5Threshold          float64
	DiskAvailThreshold      string
	DiskAvailThresholdSize  uint64
	DiskAvailRatioThreshold float64
	MemAvailThreshold       string
	MemAvailThresholdSize   uint64
	MemAvailRatioThreshold  float64
	ThresholdCron           string
	ExLogsCron              string
	PatrolCron              string
	Machines                []string
	Topn                    int
}

type BlackcatProcessConf struct {
	Keywords []string
	Machines []string
}

type BlackcatExLogConf struct {
	DirectRegex    bool
	NormalRegex    string
	ExceptionRegex string
	Ignores        []string
	LogFileName    string
	Properties     []string
	Machines       []string
}

var blackcatCron *cron.Cron = nil

func loadBlackcatCrons() {
	ExLogClearAll()

	threshold := &devopsConf.BlackcatThreshold
	fixBlackcatConfig(threshold)

	if blackcatCron != nil {
		blackcatCron.Stop()
	}
	blackcatCron = cron.New()

	if len(threshold.Machines) > 0 {
		go cronAgent(threshold)
	}

	go cronExLog(threshold)

	hourlyTip()

	blackcatCron.Start()
}

func ExLogClearAll() {
	for k := range devopsConf.Machines {
		go RpcCall(k, "ClearAll", &ExLogCommandArg{}, &ExLogCommandExecute{})
	}
}

func fixBlackcatConfig(t *BlackcatThreshold) {
	if t.DiskAvailThresholdSize == 0 {
		t.DiskAvailThresholdSize, _ = humanize.ParseBytes(t.DiskAvailThreshold)
	}
	if t.DiskAvailThreshold == "" {
		t.DiskAvailThreshold = humanize.IBytes(t.DiskAvailThresholdSize)
	}
	if t.MemAvailThresholdSize == 0 {
		t.MemAvailThresholdSize, _ = humanize.ParseBytes(t.MemAvailThreshold)
	}
	if t.MemAvailThreshold == "" {
		t.MemAvailThreshold = humanize.IBytes(t.MemAvailThresholdSize)
	}
	if t.PatrolCron == "" {
		t.PatrolCron = "@hourly"
	}
	if t.Topn == 0 {
		t.Topn = 30
	}
}

func hourlyTip() {
	blackcatCron.AddFunc("@hourly", func() {
		SendAlertMsg("黑猫正在巡逻中~", "敬请及时关注信息~")
	})
}

func cronExLog(threshold *BlackcatThreshold) {
	exLogChan := make(chan RpcResult)
	machineExLogConfs := make(map[string][]ExLogTailerConf)
	for logger, exLogConf := range devopsConf.BlackcatExLogs {
		for _, machineName := range exLogConf.Machines {
			confs, ok := machineExLogConfs[machineName]
			if !ok {
				confs = make([]ExLogTailerConf, 0)
				machineExLogConfs[machineName] = confs
			}
			machineExLogConfs[machineName] = append(confs, createExLogTailerConf(logger, exLogConf))
		}
	}

	blackcatCron.AddFunc(threshold.ExLogsCron, func() {
		for machineName, confs := range machineExLogConfs {
			logFiles := make(map[string]ExLogTailerConf)
			for _, conf := range confs {
				logFiles[conf.Logger] = conf
			}

			go RpcExecuteTimeout(machineName, &ExLogCommandArg{LogFiles: logFiles}, &ExLogCommandExecute{}, 3*time.Second, exLogChan)
		}
	})

	for x := range exLogChan {
		exLogResult := x.(*ExLogCommandResult)
		if exLogResult.Error != "" || len(exLogResult.ExLogs) != 0 {
			blackcatAlertExLog(exLogResult)
		}
	}
}

func createExLogTailerConf(logger string, conf BlackcatExLogConf) ExLogTailerConf {
	return ExLogTailerConf{
		DirectRegex:    conf.DirectRegex,
		NormalRegex:    conf.NormalRegex,
		ExceptionRegex: conf.ExceptionRegex,
		Ignores:        strings.Join(conf.Ignores, ","),
		Logger:         logger,
		LogFileName:    conf.LogFileName,
		Properties:     go_utils.MapOf(conf.Properties),
	}
}

func cronAgent(t *BlackcatThreshold) {
	resultChan := make(chan RpcResult)
	for _, machineName := range t.Machines {
		processes := make(map[string][]string)

		for processName, processConfig := range devopsConf.BlackcatProcesses {
			if go_utils.IndexOf(machineName, processConfig.Machines) >= 0 {
				processes[processName] = processConfig.Keywords
			}
		}

		local := machineName // 本行是为了在每一次循环内新建变量，以方便下面的闭包引用
		blackcatCron.AddFunc(t.ThresholdCron, func() {
			go RpcExecuteTimeout(local, &AgentCommandArg{Processes: processes, Topn: t.Topn}, &AgentCommandExeucte{}, 3*time.Second, resultChan)
		})
	}

	for x := range resultChan {
		r := x.(*AgentCommandResult)

		if r.Error != "" || beyondThreshold(r, t) {
			blackcatAlertAgent(r)
		}
	}
}

func beyondThreshold(r *AgentCommandResult, t *BlackcatThreshold) bool {
	return r.Load5 > t.Load5Threshold*float64(r.Cores) ||
		r.MemAvailable < t.MemAvailThresholdSize || 1-r.MemUsedPercent/100 < t.MemAvailRatioThreshold ||
		diskBeyondThreshold(r, t)
}

func diskBeyondThreshold(r *AgentCommandResult, t *BlackcatThreshold) bool {
	for _, du := range r.DiskUsages {
		if du.Free < t.DiskAvailThresholdSize || (1-du.UsedPercent/100) < t.DiskAvailRatioThreshold {
			return true
		}
	}

	return false
}
