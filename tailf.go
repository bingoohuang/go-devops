package main

import (
	"bufio"
	"github.com/mitchellh/go-homedir"
	"os/exec"
	"time"
)

type Tailer interface {
	Line(line string)
	Loop()
	Error(err error)
}

func Tailf(logFile string, tailer Tailer, stop chan bool, exitFunc func()) {
	defer exitFunc()

	expanded, _ := homedir.Expand(logFile)
	cmd := exec.Command("bash", "-c", "tail -F "+expanded)
	stdout, err := cmd.StdoutPipe()
	if err != nil {
		tailer.Error(err)
		return
	}
	defer stdout.Close()

	if err := cmd.Start(); err != nil {
		tailer.Error(err)
		return
	}
	defer cmd.Process.Kill()
	defer cmd.Wait()

	reader := bufio.NewReader(stdout)
	timeout := make(chan bool, 1)
	go func() {
		for {
			line, err := reader.ReadString('\n')
			timeout <- false
			if err != nil {
				tailer.Error(err)
				stop <- true
				return
			}
			tailer.Line(line)
		}
	}()

	for {
		select {
		case <-stop:
			return
		case <-timeout:
		case <-time.After(1 * time.Second):
			tailer.Loop()
		}
	}

}
