package main

import (
	"fmt"
	"github.com/docker/go-units"
	"github.com/mitchellh/go-homedir"
	"os"
	"path/filepath"
	"strconv"
	"time"
)

type CronCommandArg struct {
	Files      []string
	Type       string
	Parameters string
}

type CronCommandResult struct {
	Error      string
	CostMillis string
}

type CronCommand struct {
}

func (t *CronCommand) ExecuteCron(arg *CronCommandArg, result *CronCommandResult) error {
	fmt.Println("ExecuteCron with arg ", *arg)

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
		fmt.Println("unknown CronCommand Type", arg.Type)
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
	m := parseCommaSeparatedKeyVales(parameters)

	humanMaxSize, ok := m["maxSize"]
	if !ok {
		humanMaxSize = "100M"
	}

	maxSize, err := units.FromHumanSize(humanMaxSize)
	if err != nil {
		fmt.Println("FromHumanSize " + humanMaxSize + "err" + err.Error())
		maxSize = 100 * 1024 * 1024
	}

	o.maxSize = maxSize
	o.maxSizeStr = strconv.FormatInt(maxSize, 10)
	fmt.Println(*o)
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
	ExecuteCommands("tail -c "+o.maxSizeStr+" "+file+" > "+file+".tmp; cat "+file+".tmp > "+file, 5*time.Minute)
	fmt.Println("CopyTruncate ", file)
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

	ExecCommands(cmds)
	fmt.Println("delete files by ", cmds)
}

type DeleteOldsExecutable struct {
	days    int
	cutTime time.Time
	pattern string
}

func (o *DeleteOldsExecutable) LoadParameters(parameters string) {
	m := parseCommaSeparatedKeyVales(parameters)
	daysStr, ok := m["days"]
	if !ok {
		daysStr = "3"
	}

	o.days, _ = strconv.Atoi(daysStr)
	o.cutTime = time.Now().Truncate(24*time.Hour).AddDate(0, 0, -o.days)

	pattern, ok := m["pattern"]
	if !ok {
		fmt.Println("pattern required for DeleteOlds type")
		return
	}

	o.pattern = pattern

	fmt.Println("config:", *o)
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
	time, err := ParseFmtDate(o.pattern, base)
	if err != nil {
		//fmt.Println("parse error ", o.pattern, "for base", base, err.Error())
		return
	}

	//fmt.Println("file", file, "'s time is", time)

	if !time.After(o.cutTime) {
		os.Remove(file)
		fmt.Println("removed file", file)
	}
}
