package main

import (
	"fmt"
	"github.com/BurntSushi/toml"
)

type LogCronConf struct {
	LogCron       LogCron
	CopyTruncates map[string]CopyTruncate
	ClearLogs     map[string]ClearLog
	RotateLogs    map[string]RotateLog
}

type LogCron struct {
	At string
}

type CopyTruncate struct {
	File          string
	ThresholdSize string
}

type ClearLog struct {
	Dirs []string
}

type RotateLog struct {
	Dir      string
	Pattern  string
	DaysKeep int
}

var logCronConf LogCronConf

func loadLogCronConfig() {
	_, err := toml.DecodeFile("logcron.toml", &logCronConf)
	if err != nil {
		fmt.Println("DecodeFile error:", err)
		return
	}
}
