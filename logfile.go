package main

import (
	"github.com/dustin/go-humanize"
	"github.com/valyala/fasttemplate"
	"net/rpc"
	"os"
	"time"
)

type LogFileArg struct {
	LogPath string
	Ps      string
	Home    string
	Kill    string
	Start   string
}

type LogFileInfoResult struct {
	MachineName  string
	Error        string
	LastModified string
	FileSize     string
	TailContent  string
	CostTime     string
	ProcessInfo  string
}

type LogFileCommand int

func (t *LogFileCommand) RestartProcess(args *LogFileArg, result *LogFileInfoResult) error {
	start := time.Now()

	killTemplate := fasttemplate.New(args.Kill, "${", "}")
	killCommand := killTemplate.ExecuteString(map[string]interface{}{"ps": args.Ps})

	ExecuteCommands(killCommand, 100*time.Millisecond)
	ExecuteCommands("cd "+args.Home+"\n"+args.Start, 100*time.Millisecond)

	err := ""
	result.ProcessInfo, err = ExecuteCommands(args.Ps, 100*time.Millisecond)
	if err != "" {
		result.Error = err
	}

	result.CostTime = time.Since(start).String()
	return nil
}

func (t *LogFileCommand) TailLogFile(args *LogFileArg, result *LogFileInfoResult) error {
	start := time.Now()

	_, err := os.Stat(args.LogPath)
	if err == nil {
		stdout, stderr := ExecuteCommands("tail "+args.LogPath, 500*time.Millisecond)
		result.TailContent = stdout
		if stderr != "" {
			result.Error = stderr
		}
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

	if args.Ps != "" {
		result.ProcessInfo, _ = ExecuteCommands(args.Ps, 100*time.Millisecond)
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

func TimeoutCallLogFileCommand(machineName string, log Log, resultChan chan LogFileInfoResult, funcName string, processConfigRequired bool) {
	c := make(chan LogFileInfoResult, 1)
	machine := devopsConf.Machines[machineName]
	reply := LogFileInfoResult{
		MachineName: machineName,
	}

	process, ok := devopsConf.Processes[log.Process]
	if !ok {
		process = Process{Ps: log.Process}

		if processConfigRequired {
			reply.Error = log.Process + " is not configured"
			return
		}
	}

	if process.Home == "" || process.Kill == "" || process.Start == "" {
		reply.Error = log.Process + " is not well configured"
		return
	}

	go func() {
		err := DialAndCall(machine, func(client *rpc.Client) error {
			return client.Call("LogFileCommand."+funcName,
				&LogFileArg{
					LogPath: log.Path,
					Ps:      process.Ps,
					Home:    process.Home,
					Kill:    process.Kill,
					Start:   process.Start,
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
