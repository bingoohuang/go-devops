package main

import (
	"bufio"
	"fmt"
	"github.com/mitchellh/go-homedir"
	"github.com/patrickmn/go-cache"
	"io"
	"log"
	"os/exec"
	"sync"
	"time"
)

var tailMap sync.Map

var tailCache *cache.Cache
var cacheMux sync.Mutex

func OnEvicted(key string, val interface{}) {
	logQueue := val.(*CycleQueue)

	logQueue.stop <- true
	logQueue.cmd.Process.Kill()
	logQueue.cmd.Process.Wait()
}

func init() {
	// Create a cache with a default expiration time of 1 minutes, and which
	// purges expired items every 1 minutes
	tailCache = cache.New(1*time.Minute, 30*time.Second)
	tailCache.OnEvicted(OnEvicted)
}

func tail(logFile string, seq int) ([]byte, int) {
	cacheMux.Lock()
	defer cacheMux.Unlock()

	logQueue, found := tailCache.Get(logFile)
	if !found {
		logQueue = NewQueue(100)
		go startTail(logFile, logQueue.(*CycleQueue))
	}

	// reset expiration
	tailCache.Set(logFile, logQueue, 10*time.Minute)
	if !found {
		return nil, 0
	}

	q := logQueue.(*CycleQueue)
	return q.Get(seq)
}

func startTail(logFile string, logQueue *CycleQueue) {
	fullPathLogFile, _ := homedir.Expand(logFile)
	logQueue.cmd = exec.Command("tail", "-F", fullPathLogFile)
	stdout, err := logQueue.cmd.StdoutPipe()
	if err != nil {
		log.Fatal(err)
	}

	reader := bufio.NewReader(stdout)
	if err := logQueue.cmd.Start(); err != nil {
		log.Fatal("Buffer Error:", err)
	}

	fmt.Println("start to tail -F", fullPathLogFile)

	tmp := make([]byte, 10240)
Loop:
	for {
		select {
		case <-logQueue.stop:
			break Loop
		default:
		}

		length, err := reader.Read(tmp)
		if err != nil && err != io.EOF {
			logQueue.Add(&Node{[]byte(err.Error())})
			break
		}

		if length == 0 {
			time.Sleep(300 * time.Millisecond)
			continue
		}

		b := make([]byte, length)
		copy(b, tmp[0:length])
		logQueue.Add(&Node{b})
	}
	log.Println("Exit tail -F " + logFile)

	stdout.Close()
}
