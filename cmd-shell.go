package main

import (
	"time"
)

type CommandsArg struct {
	Command string
	Timeout time.Duration
}

type CommandsResult struct {
	MachineName string
	Error       string
	ShellKey    string
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
	result.ShellKey = "sh" + NextID()

	DelayShellChan <- &ResponseShell{
		Shell:   args.Command,
		Timeout: args.Timeout,
		ShellId: result.ShellKey,
	}

	return nil
}
