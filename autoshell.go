package main

import (
	"github.com/bingoohuang/go-utils"
	"log"
	"time"
)

var ImmediateShellChan = make(chan string)

func init() {
	go func() {
		for {
			autoShell := <-ImmediateShellChan
			log.Println("Run Auto Shell", autoShell)
			go_utils.BashTimeout(autoShell, 30*time.Second)
		}
	}()
}

type ResponseShell struct {
	Shell   string
	Timeout time.Duration
	ShellId string
	Start   time.Time
	End     time.Time
	Stdout  string
	Stderr  string
}

var DelayShellChan = make(chan *ResponseShell)

func init() {
	go func() {
		for {
			a := <-DelayShellChan
			log.Println("Run Auto Response Shell", a.Shell)

			go func() {
				a.Start = time.Now()
				a.Stdout, a.Stderr = go_utils.BashTimeout(a.Shell, a.Timeout)
				a.End = time.Now()
				WriteDbJson(exLogDb, a.ShellId, a, 24*time.Hour)
			}()
		}
	}()
}
