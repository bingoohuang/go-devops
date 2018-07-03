package main

type ShellResultCommandArg struct {
	ShellKey string
}

type ShellResultCommandResult struct {
	MachineName string
	Error       string

	Ready      bool
	ShellKey   string
	Stdout     string
	Stderr     string
	CostMillis string
}

func (t *ShellResultCommandResult) GetMachineName() string {
	return t.MachineName
}

func (t *ShellResultCommandResult) SetMachineName(machineName string) {
	t.MachineName = machineName
}

func (t *ShellResultCommandResult) GetError() string {
	return t.Error
}

func (t *ShellResultCommandResult) SetError(err error) {
	if err != nil {
		t.Error += err.Error()
	}
}

type ShellResultCommandExecute struct {
}

func (t *ShellResultCommandExecute) CreateResult(err error) RpcResult {
	result := &ShellResultCommandResult{}
	result.SetError(err)
	return result
}

func (t *ShellResultCommandExecute) CommandName() string {
	return "ShellResultCommand"
}

type ShellResultCommand int

func (t *ShellResultCommand) Execute(a *ShellResultCommandArg, r *ShellResultCommandResult) error {
	r.MachineName = hostname
	r.ShellKey = a.ShellKey

	rs := &ResponseShell{}
	ok, err := ReadDbJson(exLogDb, a.ShellKey, rs)
	if err != nil {
		return err
	}

	if !ok {
		r.Ready = false
		return nil
	}

	r.Ready = true
	r.Stdout = rs.Stdout
	r.Stderr = rs.Stderr
	r.CostMillis = rs.End.Sub(rs.Start).String()

	return nil
}
