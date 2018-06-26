package main

import (
	"bytes"
	"fmt"
	"io"
	"os/exec"
	"time"
)

func RunShell(cmds string) (string, string) {
	return RunCommandTimeout(500*time.Millisecond, "bash", "-c", cmds)
}

func RunShellTimeout(cmds string, timeout time.Duration) (string, string) {
	return RunCommandTimeout(timeout, "bash", "-c", cmds)
}

func RunCommandTimeout(timeout time.Duration, name string, args ...string) (string, string) {
	start := time.Now()
	cmd := exec.Command(name, args...)

	stderr, err := cmd.StderrPipe()
	if err != nil {
		fmt.Println("err:", err.Error())
		return "", err.Error()
	}

	stdout, err := cmd.StdoutPipe()
	if err != nil {
		fmt.Println("err:", err.Error())
		return "", err.Error()
	}

	if err := cmd.Start(); err != nil {
		fmt.Println("err:", err.Error())
		return "", err.Error()
	}

	chStdout := CreateReaderChan(stdout)
	chStderr := CreateReaderChan(stderr)

	stdoutMsg, stderrMsg := waitCommandsOutput(chStdout, chStderr, cmd, timeout)

	elapsed := time.Since(start)
	fmt.Println(hostname, time.Now(), "cost:", elapsed, "name:", name,
		"args:", args, "stdout:", stderrMsg, "stderr:", stderrMsg)

	return stdoutMsg, stderrMsg
}

func waitCommandsOutput(chStdout, chStderr <-chan string, cmd *exec.Cmd, timeout time.Duration) (string, string) {
	quit := make(chan bool)
	time.AfterFunc(timeout, func() { quit <- true })

	var bufStdout bytes.Buffer
	var bufStderr bytes.Buffer
LOOP:
	for {
		select {
		case s, ok := <-chStdout:
			if !ok {
				break LOOP
			}
			bufStdout.WriteString(s)
		case s, ok := <-chStderr:
			if !ok {
				break LOOP
			}
			bufStderr.WriteString(s)
		case <-quit:
			fmt.Println("Process Killed")
			break LOOP
		}
	}

	cmd.Process.Kill()
	cmd.Wait()
	fmt.Println("Process Waited")
	return bufStdout.String(), bufStderr.String()
}

func CreateReaderChan(closer io.Reader) <-chan string {
	ch := make(chan string)
	go func() {
		buf := make([]byte, 1024)
		for {
			n, err := closer.Read(buf)
			if n != 0 {
				ch <- string(buf[:n])
			}
			if err != nil {
				break
			}
		}
		close(ch)
	}()

	return ch
}
