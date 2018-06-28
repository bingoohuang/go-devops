package main

import (
	"github.com/dustin/go-humanize"
	"github.com/mitchellh/go-homedir"
	"github.com/valyala/fasttemplate"
	"log"
	"os"
	"path/filepath"
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

func (t *LogFileInfoResult) GetMachineName() string {
	return t.MachineName
}

func (t *LogFileInfoResult) SetMachineName(machineName string) {
	t.MachineName = machineName
}

func (t *LogFileInfoResult) GetError() string {
	return t.Error
}

func (t *LogFileInfoResult) SetError(err error) {
	if err != nil {
		t.Error += err.Error()
	}
}

func (t *LogFileCommandExecute) CreateResult(err error) RpcResult {
	result := &LogFileInfoResult{}
	result.SetError(err)
	return result
}

func (t *LogFileCommandExecute) CommandName() string {
	return "LogFileCommand"
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

type LogFileCommandExecute struct {
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
	RunShellTimeout(killCommand, 500*time.Millisecond)

	argsHome, _ := homedir.Expand(args.Home)
	RunShellTimeout("cd "+argsHome+";"+args.Start, 500*time.Millisecond)

	err := ""
	result.ProcessInfo, err = RunShellTimeout(args.Ps, 500*time.Millisecond)
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
		stdout, stderr := RunShellTimeout("tail "+args.Options+" "+logPath, 500*time.Millisecond)
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

func (t *LogFileCommand) TruncateLogFile(a *LogFileArg, r *LogFileInfoResult) error {
	start := time.Now()

	logPath, _ := homedir.Expand(a.LogPath)
	_, err := os.Stat(logPath)
	if err == nil {
		shell := `tail -100000 ` + logPath + ` > ` + logPath + `.tmp;` + `cat ` + logPath + `.tmp > ` + logPath
		AutoShellChan <- shell
		info, _ := os.Stat(logPath)

		r.FileSize = humanize.IBytes(uint64(info.Size()))
		r.LastModified = humanize.Time(info.ModTime())
	} else {
		if os.IsNotExist(err) {
			r.Error = "Log file does not exist"
		} else {
			r.Error = err.Error()
		}
	}

	r.CostTime = time.Since(start).String()
	return nil
}

func (t *LogFileCommand) LogFileInfo(args *LogFileArg, result *LogFileInfoResult) error {
	start := time.Now()

	if args.Ps != "" {
		result.ProcessInfo, _ = RunShellTimeout(args.Ps, 500*time.Millisecond)
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

	fields[4] = HumanizedKib(fields[4]) // virtual memory usage of entire process (in KiB)
	fields[5] = HumanizedKib(fields[5]) // resident set size, the non-swapped physical memory that a task has used (in KiB)

	result.ProcessInfo = strings.Join(fields, " ")
}

func GoCallLogFileCommand(wg *sync.WaitGroup, logMachineName string, log Log, resultChan chan RpcResult,
	funcName string, processConfigRequired bool, options string, logSeq int) {
	go CallLogFileCommandWait(wg, logMachineName, log, resultChan, funcName, processConfigRequired, options, logSeq)
}

func CallLogFileCommandWait(wg *sync.WaitGroup, logMachineName string, log Log, resultChan chan RpcResult,
	funcName string, processConfigRequired bool, options string, logSeq int) {
	defer wg.Done()

	CallLogFileCommand(logMachineName, log, resultChan, funcName, processConfigRequired, options, logSeq)
}

func CallLogFileCommand(logMachineName string, logf Log, resultChan chan RpcResult,
	funcName string, processConfigRequired bool, options string, logSeq int) {

	found := fullFindLogMachineName(logf, logMachineName)
	if !found {
		logMachineName, found = prefixFindLogMachineName(logf, logMachineName)
	}
	if !found {
		log.Println(logMachineName, "is unknown")
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

	process, ok := devopsConf.Processes[logf.Process]
	if !ok {
		process = Process{Ps: logf.Process}
	}

	if processConfigRequired && (process.Home == "" || process.Kill == "" || process.Start == "") {
		reply.Error = logf.Path + " is not well configured"
		resultChan <- &reply
		return
	}

	arg := &LogFileArg{
		LogPath: logf.Path,
		Ps:      process.Ps,
		Home:    process.Home,
		Kill:    process.Kill,
		Start:   process.Start,
		Options: options,
		LogSeq:  logSeq,
	}
	RpcAddrCallTimeout(machineName, machineAddress, funcName, arg, &LogFileCommandExecute{}, 1*time.Second, resultChan)
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
