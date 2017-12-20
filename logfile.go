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

	info, _ := os.Stat(args.LogPath)

	elapsed := time.Since(start)
	result.Error = ""
	result.FileSize = humanize.IBytes(uint64(info.Size()))
	result.LastModified = humanize.Time(info.ModTime())
	result.CostTime = elapsed.String()
	return nil
}

func TimeoutCallLogFileInfo(machineName, logPath string, resultChan chan LogFileInfoResult) {
	c := make(chan LogFileInfoResult, 1)
	machine := config.Machines[machineName]
	go func() {
		c <- DialAndCall(machine, CallLogFileInfo, &LogFileInfoArg{LogPath: logPath}).(LogFileInfoResult)
	}()
	select {
	case result := <-c:
		result.MachineName = machineName
		resultChan <- result
	case <-time.After(1 * time.Second):
		resultChan <- LogFileInfoResult{
			MachineName: machineName,
			Error:       "timeout",
		}
	}
}

func CallLogFileInfo(client *rpc.Client, args interface{}) interface{} {
	var reply LogFileInfoResult
	err := client.Call("LogFileInfoCommand.LogFileInfo", args, &reply)
	if err != nil {
		return MachineCommandResult{
			Error: err.Error(),
		}
	}

	return reply
}
