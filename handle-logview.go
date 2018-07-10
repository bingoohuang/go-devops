package main

import (
	"encoding/json"
	"fmt"
	"github.com/bingoohuang/go-utils"
	"github.com/dustin/go-humanize"
	"github.com/gorilla/mux"
	"html"
	"net/http"
	"strings"
)

func HandleExLog(w http.ResponseWriter, r *http.Request) {
	go_utils.HeadContentTypeHtml(w)

	vars := mux.Vars(r)
	exLogId := vars["exLogId"]
	log, err := ReadDb(exLogDb, exLogId)
	var index string

	if strings.HasPrefix(exLogId, "ex") {
		index = string(MustAsset("res/viewexlog.html"))
		index = exLogView(log, index, exLogId, err)
	} else if strings.HasPrefix(exLogId, "ag") {
		index = string(MustAsset("res/viewagent.html"))
		index = agentView(log, index, exLogId, err)
	} else if strings.HasPrefix(exLogId, "er") {
		index = string(MustAsset("res/viewerror.html"))
		index = agentError(log, index, exLogId, err)
	} else {
		index = string(MustAsset("res/viewerror.html"))
		index = strings.Replace(index, `<Error/>`, `LogId=`+exLogId+`'s format is unknown!`, -1)
	}
	w.Write([]byte(index))
}

func agentError(log []byte, index string, exLogId string, err error) string {
	if log != nil {
		index = strings.Replace(index, `<LogId/>`, exLogId, -1)
		return strings.Replace(index, `<Error/>`, string(log), -1)
	} else if err != nil {
		return strings.Replace(index, `<Error/>`, html.EscapeString(err.Error()), -1)
	} else {
		return strings.Replace(index, `<Error/>`, `LogId=`+exLogId+` Not Found!`, -1)
	}
}

func agentView(log []byte, index string, exLogId string, err error) string {
	exLog := &AgentCommandResult{}

	mergeScripts := go_utils.MergeJs(MustAsset, go_utils.FilterAssetNames(AssetNames(), ".js"))
	js := go_utils.MinifyJs(mergeScripts, devMode)
	index = strings.Replace(index, "${contextPath}", contextPath, -1)
	index = strings.Replace(index, "/*.SCRIPT*/", js, 1)

	if log != nil {
		json.Unmarshal(log, exLog)
		return buildAgentView(index, exLogId, exLog)
	} else if err != nil {
		return strings.Replace(index, `<Error/>`, html.EscapeString(err.Error()), -1)
	} else {
		return strings.Replace(index, `<Error/>`, `LogId=`+exLogId+` Not Found!`, -1)
	}
}

func buildAgentView(index, exLogId string, exLog *AgentCommandResult) string {
	index = strings.Replace(index, `<LogId/>`, exLogId, -1)
	index = strings.Replace(index, `<Timestamp/>`, exLog.Timestamp, -1)
	index = strings.Replace(index, `<Hostname/>`, exLog.Hostname, -1)
	index = strings.Replace(index, `<Load1/>`, fmt.Sprintf("%.2f", exLog.Load1), -1)
	index = strings.Replace(index, `<Load5/>`, fmt.Sprintf("%.2f", exLog.Load5), -1)
	index = strings.Replace(index, `<Load15/>`, fmt.Sprintf("%.2f", exLog.Load15), -1)
	index = strings.Replace(index, `<MemTotal/>`, humanize.IBytes(exLog.MemTotal), -1)
	index = strings.Replace(index, `<MemAvailable/>`, humanize.IBytes(exLog.MemAvailable), -1)
	index = strings.Replace(index, `<MemUsed/>`, humanize.IBytes(exLog.MemUsed), -1)
	index = strings.Replace(index, `<MemUsedPercent/>`, fmt.Sprintf("%.2f", exLog.MemUsedPercent), -1)

	diskUsages := ""
	for _, du := range exLog.DiskUsages {
		if diskUsages == "" {
			diskUsages = `<table><thead><tr><td>Path</td><td>Fstype</td><td>Total</td>
				<td>Free</td><td>Used</td><td>UsedPercent</td></tr></thead><tbody>`
		}

		diskUsages += fmt.Sprintf("<tr><td>%s</td><td>%s</td><td>%s</td><td>%s</td><td>%s</td><td>%.2f</td></tr>",
			du.Path, du.Fstype, humanize.IBytes(du.Total), humanize.IBytes(du.Free), humanize.IBytes(du.Used), du.UsedPercent)
	}

	if diskUsages != "" {
		diskUsages += "</tbody></table>"
	}

	index = strings.Replace(index, `<DiskUsages/>`, diskUsages, -1)

	top := ""
	for _, t := range exLog.Top {
		if top == "" {
			top = `<table class="sortable"><thead><tr><td>User</td><td>Pid</td><td>Ppid</td><td>%Cpu</td><td>%Mem</td>
				<td>Vsz</td><td>Rss</td><td>Tty</td><td>Stat</td><td>Start</td><td>Time</td><td>Command</td></tr></thead><tbody>`
		}

		tr := `<tr><td>{User}</td><td>%{Pid}</td><td>{Ppid}</td><td>{Cpu}</td><td>{Mem}</td>
		<td sorttable_customkey="{Vsz}">{HumanizedVsz}</td><td sorttable_customkey="{Rss}">{HumanizedRss}</td>
		<td>{Tty}</td><td>{Stat}</td><td>{Start}</td><td>{Time}</td><td>{Command}</td></tr>`
		tr = strings.Replace(tr, "{User}", t.User, 1)
		tr = strings.Replace(tr, "{Pid}", t.Pid, 1)
		tr = strings.Replace(tr, "{Ppid}", t.Ppid, 1)
		tr = strings.Replace(tr, "{Cpu}", t.Cpu, 1)
		tr = strings.Replace(tr, "{Mem}", t.Mem, 1)
		tr = strings.Replace(tr, "{Vsz}", t.Vsz, 1)
		tr = strings.Replace(tr, "{HumanizedVsz}", HumanizedKib(t.Vsz), 1)
		tr = strings.Replace(tr, "{Rss}", t.Rss, 1)
		tr = strings.Replace(tr, "{HumanizedRss}", HumanizedKib(t.Rss), 1)
		tr = strings.Replace(tr, "{Tty}", t.Tty, 1)
		tr = strings.Replace(tr, "{Stat}", t.Stat, 1)
		tr = strings.Replace(tr, "{Start}", t.Start, 1)
		tr = strings.Replace(tr, "{Time}", t.Time, 1)
		tr = strings.Replace(tr, "{Command}", t.Command, 1)
		top += tr
	}

	if top != "" {
		top += "</tbody></table>"
	}

	return strings.Replace(index, `<Top/>`, top, -1)
}

func exLogView(log []byte, index string, exLogId string, err error) string {
	exLog := &ExLog{}
	if log != nil {
		json.Unmarshal(log, exLog)
		index = replaceIndex(index, exLogId, exLog)
		return strings.Replace(index, `<Error/>`, ``, -1)
	} else if err != nil {
		index = strings.Replace(index, `<Error/>`, html.EscapeString(err.Error()), -1)
		return replaceIndex(index, exLogId, exLog)
	} else {
		index = strings.Replace(index, `<Error/>`, `LogId=`+exLogId+` Not Found!`, -1)
		return replaceIndex(index, exLogId, exLog)
	}
}

func replaceIndex(index, exLogId string, log *ExLog) string {
	index = strings.Replace(index, `<LogId/>`, exLogId, -1)
	index = strings.Replace(index, `<Hostname/>`, log.MachineName, -1)
	index = strings.Replace(index, `<Logger/>`, log.Logger, -1)
	index = strings.Replace(index, `<Properties/>`, go_utils.MapToString(log.Properties), -1)
	index = strings.Replace(index, `<ExceptionNames/>`, html.EscapeString(log.ExceptionNames), -1)
	index = strings.Replace(index, `<Timestamp/>`, log.Normal, -1)
	index = strings.Replace(index, `<ContextLogs/>`, html.EscapeString(log.Context), -1)
	return index
}
