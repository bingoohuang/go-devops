package main

import (
	"github.com/jasonlvhit/gocron"
)

func startLogCron() {
	// refer https://github.com/jasonlvhit/gocron
	gocron.Clear()
	//gocron.Every(1).Day().At(logCronConf.LogCron.At).Do(dealLogCron)
}

func dealLogCron() {

}
