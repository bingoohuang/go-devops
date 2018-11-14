package main

import (
	"bufio"
	"fmt"
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
	defer func() { _ = stdout.Close() }()

	if err := cmd.Start(); err != nil {
		tailer.Error(err)
		return
	}
	defer func() {
		_ = cmd.Process.Kill()
		_ = cmd.Wait()
	}()

	reader := bufio.NewReader(stdout)
	go func() {
		for {
			line, err := reader.ReadString('\n')
			if err != nil {
				fmt.Println(err)
				tailer.Error(err)
				stop <- true
				break
			}
			tailer.Line(line)
		}
	}()

	for {
		select {
		case <-stop:
			return
		case <-time.After(10 * time.Second):
			tailer.Loop()
		}
	}

}
