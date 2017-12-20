package main

import (
	"bytes"
	"fmt"
	"io"
	"log"
	"os/exec"
	"time"
)

type CommandsArg struct {
	Command string
	Timeout time.Duration
}

type CommandsResult struct {
	Error          string
	Stdout, Stderr string
	CostMillis     int64
}

type ShellCommand int

func (t *ShellCommand) Execute(args *CommandsArg, result *CommandsResult) error {
	start := time.Now()

	stdout, stderr := ExecuteCommands(args.Command, args.Timeout)
	elapsed := time.Since(start)
	result.Stdout = stdout
	result.Stderr = stderr
	result.CostMillis = elapsed.Nanoseconds() / 1e6
	return nil
}

func ExecuteCommands(cmds string, timeout time.Duration) (string, string) {
	cmd := exec.Command("bash", "-c", cmds)

	stderr, err := cmd.StderrPipe()
	if err != nil {
		return "", err.Error()
	}

	stdout, err := cmd.StdoutPipe()
	if err != nil {
		return "", err.Error()
	}

	if err := cmd.Start(); err != nil {
		log.Fatal(err)
	}

	chStdout := goReadOut(stdout)
	chStderr := goReadOut(stderr)

	return waitCommandsOutput(chStdout, chStderr, cmd, timeout)
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
			cmd.Process.Kill()
		}
	}
	cmd.Wait()
	return bufStdout.String(), bufStderr.String()
}

func goReadOut(closer io.ReadCloser) <-chan string {
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

func TestShell() {
	out, err := ExecuteCommands("ls\n"+"ps -ef|grep shell|grep -v grep\n"+"echo 'abc'", 3*time.Second)
	fmt.Print(out)
	fmt.Print(err)
}
