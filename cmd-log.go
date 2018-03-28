package main

import (
	"fmt"
	"github.com/dustin/go-humanize"
	"github.com/mitchellh/go-homedir"
	"github.com/valyala/fasttemplate"
	"net/rpc"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"sync"
	"time"
)

type LogFileArg struct {
	Options string
	LogPath string
	Ps      string
	Home    string
	Kill    string
	Start   string
	LogSeq  int
}

type LogFileInfoResult struct {
	MachineName  string
	Error        string
	LastModified string
	FileSize     string
	TailContent  string
	TailNextSeq  int
	CostTime     string
	ProcessInfo  string
}

type LogFileCommand int

func (t *LogFileCommand) TailFLog(args *LogFileArg, result *LogFileInfoResult) error {
	start := time.Now()
	tailContent, nextSeq := tail(args.LogPath, args.LogSeq)
	result.TailContent = string(tailContent)
	result.TailNextSeq = nextSeq

	result.CostTime = time.Since(start).String()
	return nil
}

func (t *LogFileCommand) RestartProcess(args *LogFileArg, result *LogFileInfoResult) error {
	start := time.Now()

	killTemplate := fasttemplate.New(args.Kill, "${", "}")
	killCommand := killTemplate.ExecuteString(map[string]interface{}{"ps": args.Ps})
	ExecuteCommands(killCommand, 500*time.Millisecond)

	argsHome, _ := homedir.Expand(args.Home)
	ExecuteCommands("cd "+argsHome+";"+args.Start, 500*time.Millisecond)

	err := ""
	result.ProcessInfo, err = ExecuteCommands(args.Ps, 500*time.Millisecond)
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
		stdout, stderr := ExecuteCommands("tail "+args.Options+" "+logPath, 500*time.Millisecond)
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
		ExecuteCommands("tail -100000 "+logPath+" > "+logPath+".tmp;"+
			"cat "+logPath+".tmp > "+logPath, 500*time.Millisecond)
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
		result.ProcessInfo, _ = ExecuteCommands(args.Ps, 500*time.Millisecond)
		humanizedPsOutput(result)
	}

	logPath, _ := homedir.Expand(args.LogPath)
	info, err := os.Stat(logPath)
	if err == nil {
		size := info.Size()
		if info.IsDir() {
			size, err = DirSize(args.LogPath)
			if err != nil {
				result.Error = err.Error()
			}
		}

		result.FileSize = humanize.IBytes(uint64(size))
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

func DirSize(path string) (int64, error) {
	var size int64
	err := filepath.Walk(path, func(_ string, info os.FileInfo, err error) error {
		if !info.IsDir() {
			size += info.Size()
		}
		return err
	})
	return size, err
}

func humanizedPsOutput(result *LogFileInfoResult) {
	fields := strings.Fields(result.ProcessInfo)
	if len(fields) < 6 {
		return
	}

	vszKib, _ := strconv.ParseUint(fields[4], 10, 64)
	vsz := humanize.IBytes(1024 * vszKib) // virtual memory usage of entire process (in KiB)
	vsz = strings.Replace(vsz, " ", "", 1)
	result.ProcessInfo = strings.Replace(result.ProcessInfo, fields[4], vsz, 1)

	rssKib, _ := strconv.ParseUint(fields[5], 10, 64)
	rss := humanize.IBytes(1024 * rssKib) // resident set size, the non-swapped physical memory that a task has used (in KiB)
	rss = strings.Replace(rss, " ", "", 1)
	result.ProcessInfo = strings.Replace(result.ProcessInfo, fields[5], rss, 1)
}

func CallLogFileCommand(wg *sync.WaitGroup, logMachineName string, log Log, resultChan chan *LogFileInfoResult,
	funcName string, processConfigRequired bool, options string, logSeq int) {
	if wg != nil {
		defer wg.Done()
	}

	found := fullFindLogMachineName(log, logMachineName)
	if !found {
		logMachineName, found = prefixFindLogMachineName(log, logMachineName)
	}

	if !found {
		fmt.Println(logMachineName, "is unknown")
		return
	}

	machineName, machineAddress, errorMsg := parseMachineNameAndAddress(logMachineName)
	reply := LogFileInfoResult{
		MachineName: machineName,
		Error:       errorMsg,
	}
	if errorMsg != "" {
		resultChan <- &reply
		return
	}

	process, ok := devopsConf.Processes[log.Process]
	if !ok {
		process = Process{Ps: log.Process}
	}

	if processConfigRequired && (process.Home == "" || process.Kill == "" || process.Start == "") {
		reply.Error = log.Path + " is not well configured"
		resultChan <- &reply
		return
	}

	c := make(chan LogFileInfoResult, 1)

	go func() {
		err := DialAndCall(machineAddress, func(client *rpc.Client) error {
			arg := &LogFileArg{
				LogPath: log.Path,
				Ps:      process.Ps,
				Home:    process.Home,
				Kill:    process.Kill,
				Start:   process.Start,
				Options: options,
				LogSeq:  logSeq,
			}
			return client.Call("LogFileCommand."+funcName, arg, &reply)
		})
		if err != nil {
			reply.Error = err.Error()
		}

		c <- reply
	}()

	select {
	case result := <-c:
		resultChan <- &result
	case <-time.After(1 * time.Second):
		reply.Error = "timeout"
		resultChan <- &reply
	}
}

func prefixFindLogMachineName(log Log, logMachineName string) (string, bool) {
	for _, configLogMachineName := range log.Machines {
		if strings.Index(configLogMachineName, logMachineName+":") == 0 {
			return configLogMachineName, true
		}
	}

	return "", false
}

func fullFindLogMachineName(log Log, logMachineName string) bool {
	for _, configLogMachineName := range log.Machines {
		if configLogMachineName == logMachineName {
			return true
		}
	}
	return false
}
