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

func (t *CommandsResult) GetMachineName() string {
	return t.MachineName
}

func (t *CommandsResult) SetMachineName(machineName string) {
	t.MachineName = machineName
}

func (t *CommandsResult) GetError() string {
	return t.Error
}

func (t *CommandsResult) SetError(err error) {
	if err != nil {
		t.Error += err.Error()
	}
}

func (t *ShellCommandExecute) CreateResult(err error) RpcResult {
	result := &CommandsResult{}
	result.SetError(err)
	return result
}

func (t *ShellCommandExecute) CommandName() string {
	return "ShellCommand"
}

type ShellCommand int

type ShellCommandExecute struct {
}

func (t *ShellCommand) Execute(args *CommandsArg, result *CommandsResult) error {
	start := time.Now()

	stdout, stderr := ExecuteCommands(args.Command, args.Timeout)
	elapsed := time.Since(start)
	result.Stdout = stdout
	result.Stderr = stderr
	result.CostMillis = elapsed.String()
	return nil
}

func ExecCommands(cmds string) (string, string) {
	return ExecuteCommandsWithArgs(cmds, 500*time.Millisecond)
}

func ExecuteCommands(cmds string, timeout time.Duration) (string, string) {
	return ExecuteCommandsWithArgs(cmds, timeout)
}

func ExecuteCommandsWithArgs(cmds string, timeout time.Duration) (string, string) {
	return TimeoutExecuteCommands(timeout, "bash", "-c", cmds)
}

func TimeoutExecuteCommands(timeout time.Duration, name string, args ...string) (string, string) {
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

	chStdout := goReadOut(stdout)
	chStderr := goReadOut(stderr)

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
			cmd.Process.Kill()
			fmt.Println("Process Killed")
		}
	}

	cmd.Wait()
	fmt.Println("Process Waited")
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
