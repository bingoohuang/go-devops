package main

import (
	"log"
	"time"
)

var ImmediateShellChan = make(chan string)

func init() {
	go func() {
		for {
			autoShell := <-ImmediateShellChan
			log.Println("Run Auto Shell", autoShell)
			RunShellTimeout(autoShell, 30*time.Second)
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
				a.Stdout, a.Stderr = RunShellTimeout(a.Shell, a.Timeout)
				a.End = time.Now()
				WriteDbJson(exLogDb, a.ShellId, a, 24*time.Hour)
			}()
		}
	}()
}
