package main

import (
	"github.com/dustin/go-humanize"
	"net/rpc"
	"os"
	"time"
)

type LogFileArg struct {
	LogPath string
	Process string
}

type LogFileInfoResult struct {
	MachineName  string
	Error        string
	LastModified string
	FileSize     string
	CostTime     string
	ProcessInfo  string
}

type LogFileCommand int

func (t *LogFileCommand) TruncateLogFile(args *LogFileArg, result *LogFileInfoResult) error {
	start := time.Now()

	_, err := os.Stat(args.LogPath)
	if err == nil {
		ExecuteCommands("> "+args.LogPath, 100*time.Millisecond)
		info, _ := os.Stat(args.LogPath)

		result.FileSize = humanize.IBytes(uint64(info.Size()))
		result.LastModified = humanize.Time(info.ModTime())
	} else {
		if os.IsNotExist(err) {
			result.Error = "Log file does not exist"
		} else {
			result.Error = err.Error()
		}
	}

	result.CostTime = time.Since(start).String()
	return nil
}

func (t *LogFileCommand) LogFileInfo(args *LogFileArg, result *LogFileInfoResult) error {
	start := time.Now()

	if args.Process != "" {
		result.ProcessInfo, _ = ExecuteCommands(args.Process, 100*time.Millisecond)
	}

	info, err := os.Stat(args.LogPath)
	if err == nil {
		result.FileSize = humanize.IBytes(uint64(info.Size()))
		result.LastModified = humanize.Time(info.ModTime())
	} else {
		if os.IsNotExist(err) {
			result.Error = "Log file does not exist"
		} else {
			result.Error = err.Error()
		}
	}

	result.CostTime = time.Since(start).String()
	return nil
}

func TimeoutCallLogFileCommand(machineName string, log Log, resultChan chan LogFileInfoResult, funcName string) {
	c := make(chan LogFileInfoResult, 1)
	machine := config.Machines[machineName]
	reply := LogFileInfoResult{
		MachineName: machineName,
	}

	go func() {
		err := DialAndCall(machine, func(client *rpc.Client) error {
			return client.Call("LogFileCommand."+funcName,
				&LogFileArg{
					LogPath: log.Path,
					Process: log.Process,
				}, &reply)
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
