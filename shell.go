package main

import (
	"os/exec"
	"io/ioutil"
	"os"
	"fmt"
	"io"
	"bytes"
	"time"
	"log"
)

func ExecuteCommands(cmds string) (string, string) {
	pwd, err := os.Getwd()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	file, err := ioutil.TempFile(pwd, "TempShellFile")
	if err != nil {
		return "", err.Error()
	}

	fileName := file.Name()
	fmt.Println("TempShellFile:", fileName)
	//defer os.Remove(fileName)

	file.WriteString("#!/usr/bin/env bash\n\n")
	file.WriteString(cmds)
	file.Close()

	cmd := exec.Command("bash", "-c", fileName)
	//output, err := cmd.Output()

	//return string(output), ""

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
	//
	chStdout := readOut(stdout)
	chStderr := readOut(stderr)

	quit := make(chan bool)
	time.AfterFunc(3*time.Second, func() { quit <- true })

	var bufferStdout bytes.Buffer
	var bufferStderr bytes.Buffer

LOOP:
	for {
		select {
		case s, ok := <-chStdout:
			if !ok {
				break LOOP
			}
			fmt.Print(s)
			bufferStdout.WriteString(s)
		case s, ok := <-chStderr:
			if !ok {
				break LOOP
			}
			fmt.Print(s)
			bufferStderr.WriteString(s)
		case <-quit:
			cmd.Process.Kill()
		}
	}
	fmt.Println("all done!")

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
		fmt.Println("Goroutine finished")
		close(ch)
	}()

	return ch
}

func main() {
	out, err := ExecuteCommands("ls\n")
	fmt.Print(out)
	fmt.Print(err)
}
