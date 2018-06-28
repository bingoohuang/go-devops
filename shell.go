package main

import (
	"bytes"
	"context"
	"os/exec"
	"time"
)

func RunShellTimeout(cmds string, timeout time.Duration) (string, string) {
	return RunCommandTimeout(timeout, "bash", "-c", cmds)
}

// https://medium.com/@vCabbage/go-timeout-commands-with-os-exec-commandcontext-ba0c861ed738
func RunCommandTimeout(timeout time.Duration, name string, args ...string) (string, string) {
	// Create a new context and add a timeout to it
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel() // The cancel should be deferred so resources are cleaned up

	cmd := exec.CommandContext(ctx, name, args...)

	var out bytes.Buffer
	cmd.Stdout = &out
	var eout bytes.Buffer
	cmd.Stderr = &eout

	err := cmd.Run()
	if err != nil {
		return "", err.Error()
	}

	// We want to check the context error to see if the timeout was executed.
	// The error returned by cmd.Output() will be OS specific based on what
	// happens when a process is killed.
	if ctx.Err() == context.DeadlineExceeded {
		return out.String(), "timed out"
	}

	output := out.String()
	eoutput := eout.String()

	return output, eoutput
}
