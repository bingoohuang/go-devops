package main

import (
	"os/exec"
	"bufio"
	"io"
	"log"
	"sync"
)

var tailMap sync.Map

type TailInfo struct {
	logBlocks []string
}

func tail(logFile string, seq int) string {
	return ""
}

func startTail(logFile string, tailInfo *TailInfo) {
	cmd := exec.Command("tail", "-F", logFile)
	stdout, err := cmd.StdoutPipe()
	if err != nil {
		log.Fatal(err)
	}
	defer stdout.Close()

	reader := bufio.NewReader(stdout)
	if err := cmd.Start(); err != nil {
		log.Fatal("Buffer Error:", err)
	}

	dataChan := make(chan []byte)
	go readInput(reader, dataChan)
}

func readInput(reader *bufio.Reader, dataChan chan []byte) {
	tmp := make([]byte, 10240)
	for {
		length, err := reader.Read(tmp)
		if err != nil && err != io.EOF {
			log.Println("read error", err.Error())
			break
		}

		dataChan <- tmp[0:length]
	}
}
