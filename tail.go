package main

import (
	"bufio"
	"bytes"
	"github.com/bingoohuang/go-utils"
	"github.com/mitchellh/go-homedir"
	"github.com/patrickmn/go-cache"
	"log"
	"os/exec"
	"sync"
	"time"
)

var tailCache *cache.Cache
var cacheMux sync.Mutex

func OnEvicted(key string, val interface{}) {
	logQueue := val.(*go_utils.CycleQueue)

	attach := logQueue.Attach.(*CycleQueueAttach)
	go func() {
		attach.Stop <- true
	}()
	_ = attach.Exec.Process.Kill()
	_, _ = attach.Exec.Process.Wait()
}

func init() {
	// Create a cache with a default expiration time of 1 minutes, and which
	// purges expired items every 1 minutes
	tailCache = cache.New(1*time.Minute, 30*time.Second)
	tailCache.OnEvicted(OnEvicted)
}

type CycleQueueAttach struct {
	Stop chan bool
	Exec *exec.Cmd
}

func tail(logFile string, seq int) ([]byte, int) {
	cacheMux.Lock()
	defer cacheMux.Unlock()

	logQueue, found := tailCache.Get(logFile)
	if !found {
		cycleQueue := go_utils.NewCycleQueue(100)
		cycleQueue.Attach = &CycleQueueAttach{
			Stop: make(chan bool),
		}
		logQueue = cycleQueue
		go startTail(logFile, cycleQueue, func() { tailCache.Delete(logFile) })
	}

	// reset expiration
	tailCache.Set(logFile, logQueue, 1*time.Minute)
	if !found {
		return nil, 0
	}

	q := logQueue.(*go_utils.CycleQueue)

	nodes, index := q.FetchAll(seq)

	var tailBytes bytes.Buffer
	for _, node := range nodes {
		tailBytes.Write(node.([]byte))
	}

	return tailBytes.Bytes(), index
}

func startTail(logFile string, logQueue *go_utils.CycleQueue, exitFunc func()) {
	defer exitFunc()
	expanded, _ := homedir.Expand(logFile)

	cmd := exec.Command("bash", "-c", "tail -F "+expanded)
	logQueue.Attach.(*CycleQueueAttach).Exec = cmd
	stdout, err := cmd.StdoutPipe()
	if err != nil {
		log.Fatal(err)
	}

	reader := bufio.NewReader(stdout)
	if err := cmd.Start(); err != nil {
		log.Fatal("Buffer Error:", err)
	}

	log.Println("start to tail -F ", expanded)

	defer func() {
		_ = cmd.Process.Kill()
		_ = cmd.Wait()
	}()

	tmp := make([]byte, 10240)
Loop:
	for {
		select {
		case <-logQueue.Attach.(*CycleQueueAttach).Stop:
			break Loop
		default:
		}

		length, err := reader.Read(tmp)
		if err != nil {
			logQueue.Add([]byte(err.Error()))
			break
		}

		if length == 0 {
			time.Sleep(300 * time.Millisecond)
			continue
		}

		b := make([]byte, length)
		copy(b, tmp[0:length])
		logQueue.Add(b)
	}
	log.Println("Exit tail -F " + logFile)

	_ = stdout.Close()
}
