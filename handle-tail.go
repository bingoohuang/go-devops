package main

import (
	"encoding/json"
	"github.com/bingoohuang/go-utils"
	"github.com/gorilla/mux"
	"net/http"
	"strconv"
	"strings"
	"sync"
)

func HandleTailLogFile(w http.ResponseWriter, r *http.Request) {
	go_utils.HeadContentTypeJson(w)
	vars := mux.Vars(r)
	loggerName := vars["loggerName"]
	lines := vars["lines"]
	log := devopsConf.Logs[loggerName]

	resultChan := make(chan *LogFileInfoResult, len(log.Machines))
	var wg sync.WaitGroup
	for _, logMachineName := range log.Machines {
		wg.Add(1)
		go CallLogFileCommand(&wg, logMachineName, log, resultChan,
			"TailLogFile", false, "-"+lines, 0)
	}

	wg.Wait()
	close(resultChan)

	resultsMap := make(map[string]*LogFileInfoResult)
	for commandResult := range resultChan {
		resultsMap[commandResult.MachineName] = commandResult
	}

	logs := createLogsResult(log, resultsMap)

	json.NewEncoder(w).Encode(logs)
}

func HandleTailFLog(w http.ResponseWriter, r *http.Request) {
	go_utils.HeadContentTypeJson(w)
	vars := mux.Vars(r)
	loggerName := vars["loggerName"]
	traceMobile := vars["traceMobile"]
	logSeq, _ := vars["logSeq"]

	machineLogSeqMap := parseMachineSeqs(logSeq)
	log := devopsConf.Logs[loggerName]
	if traceMobile != "0" {
		lastSlash := strings.LastIndex(log.Path, "/")
		if lastSlash >= 0 {
			log.Path = log.Path[:lastSlash] + "/" + traceMobile + ".log"
		}
	}
	machinesNum := len(log.Machines)

	newSeqMap := make(map[string]int)
	resultChan := make(chan *LogFileInfoResult, machinesNum)
	var wg sync.WaitGroup
	for _, logMachineName := range log.Machines {
		wg.Add(1)
		machineName, seq := findSeq(machineLogSeqMap, logMachineName)
		newSeqMap[machineName] = seq
		go CallLogFileCommand(&wg, logMachineName, log, resultChan,
			"TailFLog", false, "", seq)
	}
	wg.Wait()
	close(resultChan)

	resultsMap := make(map[string]*LogFileInfoResult)

	for commandResult := range resultChan {
		resultsMap[commandResult.MachineName] = commandResult
		if commandResult.Error == "" {
			newSeqMap[commandResult.MachineName] = commandResult.TailNextSeq
		}
	}

	logs := createLogsResult(log, resultsMap)
	json.NewEncoder(w).Encode(
		struct {
			Results   []*LogFileInfoResult
			NewLogSeq string
		}{
			Results:   logs,
			NewLogSeq: createMachineSeqs(newSeqMap),
		})
}

func findSeq(machineLogSeqMap map[string]int, logMachineName string) (string, int) {
	machineName := findMachineName(logMachineName)
	machineLogSeq, ok := machineLogSeqMap[machineName]
	if ok {
		return machineName, machineLogSeq
	}
	return machineName, -1
}

func createMachineSeqs(newSeqMap map[string]int) string {
	newLogSeq := ""
	for key, value := range newSeqMap {
		if newLogSeq != "" {
			newLogSeq += ","
		}
		newLogSeq += key + "|" + strconv.Itoa(value)
	}

	return newLogSeq
}

func parseMachineSeqs(logSeq string) map[string]int {
	machineLogSeqMap := make(map[string]int)
	if logSeq != "init" {
		ss := strings.Split(logSeq, ",")
		for _, pair := range ss {
			z := strings.Split(pair, "|")
			machineLogSeqMap[z[0]], _ = strconv.Atoi(z[1])
		}
	}

	return machineLogSeqMap
}
