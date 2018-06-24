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

	if err := cmd.Start(); err != nil {
		tailer.Error(err)
		return
	}
	defer cmd.Wait()

	reader := bufio.NewReader(stdout)
	timeoutCh := make(chan bool)
	go func() {
		for {
			line, err := reader.ReadString('\n')
			timeoutCh <- true
			if err != nil {
				tailer.Error(err)
				stop <- true
				return
			}
			tailer.Line(line)
		}

		close(timeoutCh)
	}()

	go func() {
		for {
			select {
			case _, ok := <-timeoutCh:
				if !ok {
					return
				}
			case <-time.After(1 * time.Second):
				tailer.Loop()
			}
		}
	}()

	<-stop
	stdout.Close()
}
