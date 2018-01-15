package main

import (
	"fmt"
	"github.com/BurntSushi/toml"
	"github.com/jasonlvhit/gocron"
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

func startLogCron() {
	// refer https://github.com/jasonlvhit/gocron
	gocron.Clear()
	if logCronConf.LogCron.At != "" {
		gocron.Every(1).Day().At(logCronConf.LogCron.At).Do(dealLogCron)
	}
}

func dealLogCron() {

}
