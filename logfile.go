package main

import (
	"github.com/dustin/go-humanize"
	"net/rpc"
	"os"
	"time"
)

type LogFileInfoArg struct {
	LogPath string
}

type LogFileInfoResult struct {
	MachineName  string
	Error        string
	LastModified string
	FileSize     string
	CostTime     string
}

type LogFileInfoCommand int

func (t *LogFileInfoCommand) LogFileInfo(args *LogFileInfoArg, result *LogFileInfoResult) error {
	start := time.Now()

	info, err := os.Stat(args.LogPath)
	if err != nil {
		return err
	}

	result.FileSize = humanize.IBytes(uint64(info.Size()))
	result.LastModified = humanize.Time(info.ModTime())
	result.CostTime = time.Since(start).String()
	return nil
}

func TimeoutCallLogFileInfo(machineName, logPath string, resultChan chan LogFileInfoResult) {
	c := make(chan LogFileInfoResult, 1)
	machine := config.Machines[machineName]
	reply := LogFileInfoResult{
		MachineName: machineName,
	}

	go func() {
		err := DialAndCall(machine, func(client *rpc.Client) error {
			return client.Call("LogFileInfoCommand.LogFileInfo", &LogFileInfoArg{LogPath: logPath}, &reply)
		})
		if err != nil {
			reply.Error = err.Error()
		}

		c <- reply
	}()
	select {
	case result := <-c:
		resultChan <- result
	case <-time.After(1 * time.Second):
		reply.Error = "timeout"
		resultChan <- reply
	}
}