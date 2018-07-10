package main

import (
	"github.com/bingoohuang/go-utils"
	"github.com/docker/go-units"
	"github.com/mitchellh/go-homedir"
	"log"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"
)

type CronCommandArg struct {
	Files      []string
	Type       string
	Parameters string
}

type CronCommandResult struct {
	MachineName string
	Error       string
	CostMillis  string
}

type CronCommand int
type CronCommandExecute struct {
}

func (t *CronCommandResult) GetMachineName() string {
	return t.MachineName
}

func (t *CronCommandResult) SetMachineName(machineName string) {
	t.MachineName = machineName
}

func (t *CronCommandResult) GetError() string {
	return t.Error
}

func (t *CronCommandResult) SetError(err error) {
	if err != nil {
		t.Error += err.Error()
	}
}

func (t *CronCommandExecute) CreateResult(err error) RpcResult {
	result := &CronCommandResult{}
	result.SetError(err)
	return result
}

func (t *CronCommandExecute) CommandName() string {
	return "CronCommand"
}

func (t *CronCommand) Execute(arg *CronCommandArg, result *CronCommandResult) error {
	log.Println("ExecuteCron with arg ", *arg)

	start := time.Now()
	elapsed := time.Since(start)

	var executable CronExecutable
	if arg.Type == "CopyTruncate" {
		executable = &CopyTruncateCronExecutable{}
	} else if arg.Type == "Delete" {
		executable = &DeleteCronExecutable{}
	} else if arg.Type == "DeleteOlds" {
		executable = &DeleteOldsExecutable{}
	} else {
		result.Error = "unknown CronCommand Type " + arg.Type
		log.Println("unknown CronCommand Type", arg.Type)
		executable = nil
	}

	if executable != nil {
		executable.LoadParameters(arg.Parameters)
		executable.Execute(arg.Files)
	}

	result.CostMillis = elapsed.String()
	return nil
}

type CronExecutable interface {
	LoadParameters(parameters string)
	Execute(files []string)
}

type CopyTruncateCronExecutable struct {
	maxSize    int64
	maxSizeStr string
}

func (o *CopyTruncateCronExecutable) LoadParameters(parameters string) {
	m := go_utils.ParseMapString(parameters, ",", "=")

	humanMaxSize, ok := m["maxSize"]
	if !ok {
		humanMaxSize = "100M"
	}

	maxSize, err := units.FromHumanSize(humanMaxSize)
	if err != nil {
		log.Println("FromHumanSize " + humanMaxSize + "err" + err.Error())
		maxSize = 100 * 1024 * 1024
	}

	o.maxSize = maxSize
	o.maxSizeStr = strconv.FormatInt(maxSize, 10)
	log.Println(*o)
}

func (o *CopyTruncateCronExecutable) Execute(files []string) {
	for _, xfile := range files {
		file, _ := homedir.Expand(xfile)

		stat, err := os.Stat(file)
		if os.IsNotExist(err) {
			continue
		}

		if stat.IsDir() {
			filepath.Walk(file, func(path string, info os.FileInfo, err error) error {
				if !info.IsDir() && info.Size() > o.maxSize {
					o.tailMaxSize(path)
				}
				return nil
			})
		} else if stat.Size() > o.maxSize {
			o.tailMaxSize(file)
		}
	}
}

func (o *CopyTruncateCronExecutable) tailMaxSize(file string) {
	shell := `tail -c ` + o.maxSizeStr + ` ` + file + ` > ` + file + `.tmp; cat ` + file + `.tmp>` + file
	ImmediateShellChan <- shell
	log.Println("CopyTruncate ", file)
}

type DeleteCronExecutable struct {
}

func (o *DeleteCronExecutable) LoadParameters(parameters string) {

}

func (o *DeleteCronExecutable) Execute(files []string) {
	cmds := "rm -fr "
	for _, xfile := range files {
		file, _ := homedir.Expand(xfile)
		cmds += file + " "
	}

	go_utils.BashTimeout(cmds, 100*time.Millisecond)
	log.Println("delete files by ", cmds)
}

type DeleteOldsExecutable struct {
	days     int
	cutTime  time.Time
	patterns []string
}

func (o *DeleteOldsExecutable) LoadParameters(parameters string) {
	m := go_utils.ParseMapString(parameters, ",", "=")
	daysStr, ok := m["days"]
	if !ok {
		daysStr = "3"
	}

	o.days, _ = strconv.Atoi(daysStr)
	o.cutTime = time.Now().Truncate(24*time.Hour).AddDate(0, 0, -o.days)

	pattern, ok := m["pattern"]
	if !ok {
		log.Println("pattern required for DeleteOlds type")
		return
	}

	o.patterns = strings.Split(pattern, "|")

	log.Println("config:", *o)
}

func (o *DeleteOldsExecutable) Execute(files []string) {
	for _, xfile := range files {
		file, _ := homedir.Expand(xfile)
		stat, err := os.Stat(file)
		if os.IsNotExist(err) {
			continue
		}

		if stat.IsDir() {
			filepath.Walk(file, func(path string, info os.FileInfo, err error) error {
				if !info.IsDir() {
					o.deleteFile(path)
				}
				return nil
			})
		} else {
			o.deleteFile(file)
		}
	}
}

func (o *DeleteOldsExecutable) deleteFile(file string) {
	base := filepath.Base(file)

	time := o.fileTime(base)
	// fmt.Println("file", file, "'s time is", time)

	if !time.After(o.cutTime) {
		os.Remove(file)
		log.Println("removed file", file)
	}
}

func (o *DeleteOldsExecutable) fileTime(base string) time.Time {
	for _, pattern := range o.patterns {
		time, err := go_utils.ParseFmtDate(pattern, base)
		if err == nil {
			return time
		}
	}

	return time.Now()
}
