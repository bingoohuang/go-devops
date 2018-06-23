package main

type AgentCommandArg struct {
	Processes map[string][]string
}

type AgentCommandResult struct {
	Load1  float64
	Load5  float64
	Load15 float64

	MemTotal       uint64
	MemAvailable   uint64
	MemUsed        uint64
	MemUsedPercent float64

	Processes map[string]PsAuxItem
	Top       []PsAuxItem

	MachineName string
	Error       string
}

type AgentCommand int

type AgentCommandExeucte struct {
}

func (t *AgentCommandResult) GetMachineName() string {
	return t.MachineName
}

func (t *AgentCommandResult) SetMachineName(machineName string) {
	t.MachineName = machineName
}

func (t *AgentCommandResult) GetError() string {
	return t.Error
}

func (t *AgentCommandResult) SetError(err error) {
	if err != nil {
		t.Error += err.Error()
	}
}

func (t *AgentCommandExeucte) CreateResult(err error) RpcResult {
	result := &AgentCommandResult{}
	result.SetError(err)
	return result
}

func (t *AgentCommandExeucte) CommandName() string {
	return "AgentCommand"
}

func (t *AgentCommand) Execute(a *AgentCommandArg, r *AgentCommandResult) error {
	load := Load()
	r.Load1 = load.Load1
	r.Load5 = load.Load5
	r.Load15 = load.Load15

	memory := Memory()
	r.MemTotal = memory.Total
	r.MemAvailable = memory.Available
	r.MemUsed = memory.Used
	r.MemUsedPercent = memory.UsedPercent

	processes := make(map[string]PsAuxItem)
	for k, v := range a.Processes {
		grep := PsAuxGrep(v...)
		if len(grep) > 0 {
			processes[k] = *grep[0]
		}
	}

	r.Processes = processes
	top := PsAuxTop(10)
	r.Top = make([]PsAuxItem, 0)
	for _, i := range top {
		r.Top = append(r.Top, *i)
	}

	return nil
}
