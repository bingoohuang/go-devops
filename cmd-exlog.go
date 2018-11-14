package main

import (
	"reflect"
	"sync"
	"time"
)

type ExLogCommandArg struct {
	LogFiles map[string]ExLogTailerConf
}

type ExLogCommandResult struct {
	ExLogs         []ExLog
	MachineName    string
	Error          string
	Hostname       string
	Timestamp      string
	MessageTargets []string // 消息发送目标
}

type ExLogCommand int

type ExLogCommandExecute struct {
}

func (t *ExLogCommandResult) GetMachineName() string {
	return t.MachineName
}

func (t *ExLogCommandResult) SetMachineName(machineName string) {
	t.MachineName = machineName
}

func (t *ExLogCommandResult) GetError() string {
	return t.Error
}

func (t *ExLogCommandResult) SetError(err error) {
	if err != nil {
		t.Error += err.Error()
	}
}

func (t *ExLogCommandExecute) CreateResult(err error) RpcResult {
	result := &ExLogCommandResult{}
	result.SetError(err)
	return result
}

func (t *ExLogCommandExecute) CommandName() string {
	return "ExLogCommand"
}

type ExLogTailerRuntime struct {
	Conf      *ExLogTailerConf
	ExLogChan chan ExLog
	Stop      chan bool
}

var exLogChanMap sync.Map

func (t *ExLogCommand) ClearAll(a *ExLogCommandArg, r *ExLogCommandResult) error {
	exLogChanMap.Range(func(k, v interface{}) bool {
		v.(*ExLogTailerRuntime).Stop <- true
		exLogChanMap.Delete(k)
		return true
	})

	return nil
}

func (t *ExLogCommand) Execute(a *ExLogCommandArg, r *ExLogCommandResult) error {
	r.Hostname = hostname
	r.Timestamp = time.Now().Format("2006-01-02 15:04:05")

	for k, v := range a.LogFiles {
		m, ok := exLogChanMap.Load(k)
		if ok {
			rt := m.(*ExLogTailerRuntime)
			if !reflect.DeepEqual(rt.Conf, &v) {
				rt.Stop <- true

				exLogChanMap.Delete(k)
				_ = StartNewTailer(k, v)
			}
		} else {
			err := StartNewTailer(k, v)
			if err != nil {
				r.Error = err.Error()
				return err
			}
		}
	}

	r.ExLogs = make([]ExLog, 0)
	exLogChanMap.Range(func(k, v interface{}) bool {
		exLogChan := v.(*ExLogTailerRuntime).ExLogChan

		for {
			select {
			case x, ok := <-exLogChan:
				if ok {
					r.ExLogs = append(r.ExLogs, x)
				} else {
					exLogChanMap.Delete(k)
					return true
				}
			default:
				return true
			}
		}
	})

	return nil
}

func StartNewTailer(k string, v ExLogTailerConf) error {
	rt := ExLogTailerRuntime{
		Conf:      &v,
		ExLogChan: make(chan ExLog, 10),
		Stop:      make(chan bool, 2),
	}

	_, existed := exLogChanMap.LoadOrStore(k, &rt)
	if existed {
		return nil
	}

	tailer, err := NewExLogTailer(rt.ExLogChan, rt.Conf)
	if err != nil {
		return err
	}

	go Tailf(v.LogFileName, tailer, rt.Stop, func() { exLogChanMap.Delete(k) })
	return nil
}
