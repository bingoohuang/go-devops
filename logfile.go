package main

import (
	"github.com/dustin/go-humanize"
	"github.com/mitchellh/go-homedir"
	"github.com/valyala/fasttemplate"
	"net/rpc"
	"os"
	"time"
)

type LogFileArg struct {
	Options string
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

	ExecuteCommands(killCommand, 500*time.Millisecond, true)

	argsHome, _ := homedir.Expand(args.Home)
	ExecuteCommands("cd "+argsHome+";"+args.Start, 1000*time.Millisecond, false)
	//randomShellName := RandStringBytesMaskImpr(16) + ".sh"
	//ExecuteCommands("cd "+args.Home+"\n"+
	//	"echo \""+args.Start+"\">"+randomShellName+"\n"+
	//	"chmod +x "+randomShellName+"\n"+
	//	"./"+randomShellName+"\n"+
	//	"rm "+randomShellName, 500*time.Millisecond)

	err := ""
	result.ProcessInfo, err = ExecuteCommands(args.Ps, 500*time.Millisecond, true)
	if err != "" {
		result.Error = err
	}

	result.CostTime = time.Since(start).String()
	return nil
}

func (t *LogFileCommand) TailLogFile(args *LogFileArg, result *LogFileInfoResult) error {
	start := time.Now()

	logPath, _ := homedir.Expand(args.LogPath)
	_, err := os.Stat(logPath)
	if err == nil {
		stdout, stderr := ExecuteCommands("tail "+args.Options+" "+logPath, 500*time.Millisecond, true)
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

	logPath, _ := homedir.Expand(args.LogPath)
	_, err := os.Stat(logPath)
	if err == nil {
		ExecuteCommands("> "+logPath, 500*time.Millisecond, true)
		info, _ := os.Stat(logPath)

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
		result.ProcessInfo, _ = ExecuteCommands(args.Ps, 500*time.Millisecond, true)
	}

	logPath, _ := homedir.Expand(args.LogPath)
	info, err := os.Stat(logPath)
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

func TimeoutCallLogFileCommand(machineName string, log Log, resultChan chan LogFileInfoResult,
	funcName string, processConfigRequired bool, options string) {
	c := make(chan LogFileInfoResult, 1)

	reply := LogFileInfoResult{
		MachineName: machineName,
	}

	machine, ok := devopsConf.Machines[machineName]
	if !ok {
		reply.Error = machineName + " is unknown"
		resultChan <- reply
		return
	}

	process, ok := devopsConf.Processes[log.Process]
	if !ok {
		process = Process{Ps: log.Process}
	}

	if processConfigRequired && (process.Home == "" || process.Kill == "" || process.Start == "") {
		reply.Error = log.Path + " is not well configured"
		resultChan <- reply
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
					Options: options,
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
