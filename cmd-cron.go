package main

import (
	"github.com/docker/go-units"
	"github.com/metakeule/fmtdate"
	"log"
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
	MachineName    string
	Error          string
	Stdout, Stderr string
	CostMillis     string
}

type CronCommand int

func (t *CronCommand) ExecuteCron(arg *CronCommandArg, result *CronCommandResult) error {
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
	m := parseCommaSeparatedKeyVales(parameters)

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
}

func (o *CopyTruncateCronExecutable) Execute(files []string) {
	for _, file := range files {
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
	ExecCommands("tail -c " + o.maxSizeStr + " > " + file + ".tmp; cat " + file + ".tmp > " + file)
	log.Println("CopyTruncate ", file)
}

type DeleteCronExecutable struct {
}

func (o *DeleteCronExecutable) LoadParameters(parameters string) {

}

func (o *DeleteCronExecutable) Execute(files []string) {
	cmds := "rm -fr "
	for _, file := range files {
		cmds += file + " "
	}

	ExecCommands(cmds)
	log.Println("delete files by ", cmds)
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
		log.Println("pattern required for DeleteOlds type")
		return
	}

	o.pattern = pattern
}

func (o *DeleteOldsExecutable) Execute(files []string) {
	for _, file := range files {
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
	time, err := fmtdate.Parse(o.pattern, filepath.Base(file))
	if err == nil && time.Before(o.cutTime) {
		os.Remove(file)
		log.Println("removed file", file)
	}
}
