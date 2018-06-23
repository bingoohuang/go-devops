package main

import (
	"fmt"
	"reflect"
)

type ExLogCommandArg struct {
	LogFiles map[string]ExLogTailerConf
}

type ExLogCommandResult struct {
	ExLogs      []ExLog
	MachineName string
	Error       string
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

var exLogChanMap = make(map[string]*ExLogTailerRuntime)

func (t *ExLogCommand) Execute(a *ExLogCommandArg, r *ExLogCommandResult) error {
	fmt.Println("Execute ExLogCommand:", a)
	for k, v := range a.LogFiles {
		m, ok := exLogChanMap[k]
		if !ok {
			err := StartNewTailer(k, &v)
			fmt.Println("Start New Tailer")
			if err != nil {
				r.Error = err.Error()
				return err
			}
		} else {
			if !reflect.DeepEqual(m.Conf, &v) {
				fmt.Println("ReStart Tailer")
				m.Stop <- true
				StartNewTailer(k, &v)
			} else {
				fmt.Println("Resure old Tailer", m)
			}
		}
	}

	r.ExLogs = make([]ExLog, 0)
	for k, v := range exLogChanMap {
		for {
			select {
			case x, ok := <-v.ExLogChan:
				if ok {
					r.ExLogs = append(r.ExLogs, x)
				} else {
					delete(exLogChanMap, k)
					return nil
				}
			default:
				return nil
			}
		}
	}

	return nil
}

func StartNewTailer(k string, v *ExLogTailerConf) error {
	rt := ExLogTailerRuntime{
		Conf:      v,
		ExLogChan: make(chan ExLog, 10),
		Stop:      make(chan bool, 2),
	}
	exLogChanMap[k] = &rt
	tailer, err := NewExLogTailer(rt.ExLogChan, rt.Conf)
	if err != nil {
		return err
	}

	go Tailf(v.LogFileName, tailer, rt.Stop)
	return nil
}
