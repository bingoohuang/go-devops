package main

import (
	"log"
	"time"
)

var AutoShellChan = make(chan string)

func init() {
	go func() {
		for {
			autoShell := <-AutoShellChan
			log.Println("Run Auto Shell", autoShell)
			RunShellTimeout(autoShell, 30*time.Second)
		}
	}()
}
