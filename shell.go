package main

import (
	"bytes"
	"fmt"
	"io"
	"os/exec"
	"time"
)

type CommandsArg struct {
	Command string
	Timeout time.Duration
}

type CommandsResult struct {
	MachineName    string
	Error          string
	Stdout, Stderr string
	CostMillis     string
}

type ShellCommand int

func (t *ShellCommand) Execute(args *CommandsArg, result *CommandsResult) error {
	start := time.Now()

	stdout, stderr := ExecuteCommands(args.Command, args.Timeout)
	elapsed := time.Since(start)
	result.Stdout = stdout
	result.Stderr = stderr
	result.CostMillis = elapsed.String()
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
		return "", err.Error()
	}

	chStdout := goReadOut(stdout)
	chStderr := goReadOut(stderr)

	stdoutMsg, stderrMsg := waitCommandsOutput(chStdout, chStderr, cmd, timeout)

	fmt.Println("cmds:", cmds, "stdout:", stderrMsg, "stderr:", stderrMsg)

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

/*
https://superuser.com/questions/171858/how-do-i-interpret-the-results-of-the-ls-l-command

      +-permissions that apply to the owner
      |
      |     +-permissions that apply to all other users
      |     |
      |     |  +-number of hard links
      |     |  |
      |     |  |             +-size      +-last modification date and time
     _|_   _|_ |            _|__ ________|_______
    drwxr-xr-x 2 ataka root 4096 2008-11-04 16:58 ataka
        ___      _____ ____                       _____
         |         |    |                           |
         |         |    |                           +-name of file or directory
         |         |    |
         |         |    +-the group that the group permissions applies to
         |         |
         |         +-owner
         |
         +-permissions that apply to users who are members of the group
*/

func TestShell() {
	out, err := ExecuteCommands("ls\n"+"ps -ef|grep shell|grep -v grep\n"+"echo 'abc'", 3*time.Second)
	fmt.Print(out)
	fmt.Print(err)
}
