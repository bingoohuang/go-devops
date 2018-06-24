package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/bingoohuang/go-utils"
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

	index := string(MustAsset("res/exlogview.html"))

	exLog := &ExLog{}
	if log != nil {
		json.Unmarshal(log, exLog)
		fmt.Println("exlog:%v", exLog)
		index = replaceIndex(index, exLogId, exLog)
		index = strings.Replace(index, `<Error/>`, ``, -1)
	} else if err != nil {
		index = strings.Replace(index, `<Error/>`, html.EscapeString(err.Error()), -1)
		index = replaceIndex(index, exLogId, exLog)
	} else {
		index = strings.Replace(index, `<Error/>`, `LogId=`+exLogId+` Not Found!`, -1)
		index = replaceIndex(index, exLogId, exLog)
	}

	w.Write([]byte(index))
}

func replaceIndex(index, exLogId string, log *ExLog) string {
	index = strings.Replace(index, `<LogId/>`, exLogId, -1)
	index = strings.Replace(index, `<Hostname/>`, log.MachineName, -1)
	index = strings.Replace(index, `<Logger/>`, log.Logger, -1)
	index = strings.Replace(index, `<Properties/>`, CreateKeyValuePairs(log.Properties), -1)
	index = strings.Replace(index, `<ExceptionNames/>`, html.EscapeString(log.ExceptionNames), -1)
	index = strings.Replace(index, `<Timestamp/>`, log.Normal, -1)
	index = strings.Replace(index, `<ContextLogs/>`, html.EscapeString(log.Context), -1)
	return index
}

func CreateKeyValuePairs(m map[string]string) string {
	b := new(bytes.Buffer)
	fmt.Fprintf(b, "%v", m)
	return b.String()
}
