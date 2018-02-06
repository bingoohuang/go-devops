package main

import (
	"time"
	"strings"
	"github.com/docker/go-units"
	"log"
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
}

func (o *CopyTruncateCronExecutable) LoadParameters(parameters string) {

}

func (o *CopyTruncateCronExecutable) Execute(files []string) {

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
}

type DeleteOldsExecutable struct {
	maxSize int64
}

func (o *DeleteOldsExecutable) LoadParameters(parameters string) {
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
}

func parseCommaSeparatedKeyVales(str string) map[string]string {
	parts := strings.Split(str, ",")

	m := make(map[string]string)
	for _, part := range parts {
		p := strings.TrimSpace(part)
		if p == "" {
			continue
		}

		index := strings.Index(p, "=")
		if index > 0 {
			key := p[0:index]
			val := p[index+1:]
			k := strings.TrimSpace(key)
			v := strings.TrimSpace(val)

			if k != "" {
				m[k] = v
			}
		} else if index < 0 {
			m[p] = ""
		}
	}

	return m
}

func (o *DeleteOldsExecutable) Execute(files []string) {
	
}
