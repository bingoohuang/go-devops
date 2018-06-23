package main

import (
	"fmt"
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
	Machines                []string
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
	threshold := devopsConf.BlackcatThreshold
	if threshold.DiskAvailThresholdSize == 0 {
		threshold.DiskAvailThresholdSize, _ = humanize.ParseBytes(threshold.DiskAvailThreshold)
	}
	if threshold.MemAvailThresholdSize == 0 {
		threshold.MemAvailThresholdSize, _ = humanize.ParseBytes(threshold.MemAvailThreshold)
	}

	if blackcatCron != nil {
		blackcatCron.Stop()
	}
	blackcatCron = cron.New()

	if len(threshold.Machines) > 0 {
		go cronAgent(threshold)
	}

	go cronExLog(threshold)

	blackcatCron.Start()
}

func cronExLog(threshold BlackcatThreshold) {
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

	fmt.Println("machineExLogConfs:", machineExLogConfs)

	blackcatCron.AddFunc(threshold.ExLogsCron, func() {
		for machineName, confs := range machineExLogConfs {
			logFiles := make(map[string]ExLogTailerConf)
			for _, conf := range confs {
				logFiles[conf.Logger] = conf
			}

			go RpcCallTimeout(machineName, "", "Execute",
				&ExLogCommandArg{LogFiles: logFiles},
				&ExLogCommandExecute{}, 3*time.Second, exLogChan)
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

func cronAgent(threshold BlackcatThreshold) {
	resultChan := make(chan RpcResult)
	fmt.Println("threshold:", threshold)
	for _, machineName := range threshold.Machines {
		processes := make(map[string][]string)

		for processName, processConfig := range devopsConf.BlackcatProcesses {
			if go_utils.IndexOf(machineName, processConfig.Machines) >= 0 {
				processes[processName] = processConfig.Keywords
			}
		}

		fmt.Println("machineName:", machineName, "processes:", processes)
		blackcatCron.AddFunc(threshold.ThresholdCron, func() {
			go RpcCallTimeout(machineName, "", "Execute",
				&AgentCommandArg{Processes: processes},
				&AgentCommandExeucte{}, 3*time.Second, resultChan)
		})
	}

	for x := range resultChan {
		agentResult := x.(*AgentCommandResult)
		if agentResult.Error != "" || beyondThreadhold(agentResult, threshold) {
			blackcatAlertAgent(agentResult)
		}
	}
}

func beyondThreadhold(result *AgentCommandResult, threshold BlackcatThreshold) bool {
	return result.Load5 >= threshold.Load5Threshold ||
		result.MemAvailable <= threshold.MemAvailThresholdSize ||
		1-result.MemUsedPercent <= threshold.MemAvailRatioThreshold
}
