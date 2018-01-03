package main

import (
	"bufio"
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

	logQueue.cmd.Process.Kill()
}

func init() {
	// Create a cache with a default expiration time of 1 minutes, and which
	// purges expired items every 1 minutes
	tailCache = cache.New(1*time.Minute, 1*time.Minute)
	tailCache.OnEvicted(OnEvicted)
}

func tail(logFile string, seq int) ([]byte, int) {
	logQueue, found := tailCache.Get(logFile)
	newCreate := false
	if !found {
		logQueue, newCreate = createCache(logFile)
	}

	q := logQueue.(*CycleQueue)

	if newCreate {
		return nil, 0
	}

	return q.Get(seq)
}
func createCache(logFile string) (interface{}, bool) {
	cacheMux.Lock()
	defer cacheMux.Unlock()

	logQueue, found := tailCache.Get(logFile)
	if found {
		return logQueue, false
	}

	logQueue = NewQueue(100)
	tailCache.Set(logFile, logQueue, cache.DefaultExpiration)

	go startTail(logFile, logQueue.(*CycleQueue))

	return logQueue, true
}

func startTail(logFile string, logQueue *CycleQueue) {
	logQueue.cmd = exec.Command("tail", "-F", logFile)
	stdout, err := logQueue.cmd.StdoutPipe()
	if err != nil {
		log.Fatal(err)
	}

	reader := bufio.NewReader(stdout)
	if err := logQueue.cmd.Start(); err != nil {
		log.Fatal("Buffer Error:", err)
	}

	tmp := make([]byte, 10240)
	for {
		length, err := reader.Read(tmp)
		if err != nil && err != io.EOF {
			logQueue.Add(&Node{[]byte(err.Error())})
			break
		}

		if length > 0 {
			b := make([]byte, length)
			copy(b, tmp[0:length])
			logQueue.Add(&Node{b})
		}
	}

	stdout.Close()
}
