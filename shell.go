package main

import (
	"os/exec"
	"fmt"
	"io"
	"bytes"
	"log"
	"time"
)

func ExecuteCommands(cmds string, timeoutSeconds time.Duration) (string, string) {
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

	log.Println("pid:", cmd.Process.Pid)
	//
	chStdout := readOut(stdout)
	chStderr := readOut(stderr)

	quit := make(chan bool)
	time.AfterFunc(timeoutSeconds, func() { quit <- true })

	var bufferStdout bytes.Buffer
	var bufferStderr bytes.Buffer

LOOP:
	for {
		select {
		case s, ok := <-chStdout:
			if !ok {
				cmd.Process.Kill()
				break LOOP
			}
			bufferStdout.WriteString(s)
		case s, ok := <-chStderr:
			if !ok {
				break LOOP
			}
			bufferStderr.WriteString(s)
		case <-quit:
			cmd.Process.Kill()
		}
	}

	cmd.Wait()

	return bufferStdout.String(), bufferStderr.String()
}

func readOut(closer io.ReadCloser) chan string {
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

func main() {
	out, err := ExecuteCommands("ls", 3*time.Second)
	fmt.Print(out)
	fmt.Print(err)
}
