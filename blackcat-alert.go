package main

import (
	"errors"
	"fmt"
	"github.com/bingoohuang/go-utils"
	"github.com/dustin/go-humanize"
	"strings"
	"sync"
	"time"
)

func blackcatAlertExLog(result *ExLogCommandResult) {
	for _, log := range result.ExLogs {
		key := "ex" + NextID()
		log.MachineName = result.Hostname
		WriteDbJson(exLogDb, key, log, 7*24*time.Hour)

		content := "Host: " + result.Hostname + "\nTs: " + log.Normal + "\nLogger: " + log.Logger +
			"\nLogTag: " + log.Normal + "\nFoundTs: " + result.Timestamp

		if len(log.Properties) > 0 {
			content += "\nProperties: " + go_utils.MapToString(log.Properties)
		}

		content += "\n" + linkLogId(key) + "\nEx: " + log.ExceptionNames
		AddAlertMsg("发现异常啦~", content)
	}

	if result.Error != "" {
		key := "er" + NextID()
		WriteDb(exLogDb, key, []byte(result.Error), 7*24*time.Hour)
		content := "\n" + linkLogId(key) + "\nEx: " + result.Error
		AddAlertMsg("发现错误啦~", content)
	}
}

func blackcatAlertAgent(result *AgentCommandResult) {
	key := "ag" + NextID()
	WriteDbJson(exLogDb, key, result, 7*24*time.Hour)

	content := make([]string, 0)
	if result.Error != "" {
		content = append(content, "Error: "+result.Error)
	}

	content = append(content, "Host: "+result.Hostname)
	content = append(content, linkLogId(key))

	if result.MemTotal > 0 {
		threshold := &devopsConf.BlackcatThreshold
		Load5Threshold := threshold.Load5Threshold * float64(result.Cores)
		if result.Load5 > Load5Threshold {
			content = append(content, "负载告警。Load5 "+fmt.Sprintf("%.2f", result.Load5)+"高于"+fmt.Sprintf("%.2f", Load5Threshold))
		}

		if result.MemAvailable < threshold.MemAvailThresholdSize {
			content = append(content, "内存告警。 可用"+humanize.IBytes(result.MemAvailable)+"低于"+threshold.MemAvailThreshold)
		}
		memAvailRatio := 1 - result.MemUsedPercent/100
		if memAvailRatio < threshold.MemAvailRatioThreshold {
			content = append(content, "内存告警。比例"+fmt.Sprintf("%.2f", memAvailRatio)+"低于"+fmt.Sprintf("%.2f", threshold.MemAvailRatioThreshold))
		}

		for _, du := range result.DiskUsages {
			if du.Free < threshold.DiskAvailThresholdSize {
				content = append(content, "磁盘告警。"+du.Path+"可用"+humanize.IBytes(result.MemAvailable)+"低于"+threshold.MemAvailThreshold)
			}
			availRatio := 1 - du.UsedPercent/100
			if availRatio < threshold.DiskAvailRatioThreshold {
				content = append(content, "磁盘告警。"+du.Path+"可用比例"+fmt.Sprintf("%.2f", availRatio)+"低于"+fmt.Sprintf("%.2f", threshold.DiskAvailRatioThreshold))
			}
		}
	}

	AddAlertMsg("发来警报啦~", strings.Join(content, "\n"))
}

func linkLogId(key string) string {
	threshold := &devopsConf.BlackcatThreshold
	if threshold.ExLogViewUrlPrefix == "" {
		return `LogId: ` + key
	} else {
		return `<a href="` + threshold.ExLogViewUrlPrefix + `/exlog/` + key + `">LogId</a>: ` + key
	}
}

var exLogDb = OpenDb("./exlogdb")

type Msgs []Msg

func (msgs *Msgs) Clear() {
	*msgs = append([]Msg{})
}

func (msgs Msgs) firstHead() string {
	if len(msgs) > 0 {
		return msgs[0].Head
	}
	return ""
}

// wx
func (msgs Msgs) wxContent() (ret string) {
	for _, m := range msgs {
		ret += m.wxContent()
	}
	return
}

// dingding
func (msgs Msgs) dingMarkdown() (ret string) {
	for _, m := range msgs {
		ret += m.dingMarkdown()
	}
	return
}

type MsgContext struct {
	lock *sync.RWMutex
	m    Msgs
}

func (m *MsgContext) PopOut() Msgs {
	m.lock.Lock()
	defer m.lock.Unlock()
	msgs := m.m
	m.m.Clear()
	return msgs
}

func (m *MsgContext) Add(msg Msg) {
	m.lock.Lock()
	defer m.lock.Unlock()
	m.m = append(m.m, msg)
}

func newMsgContext() *MsgContext {
	return &MsgContext{
		lock: new(sync.RWMutex),
		m:    make([]Msg, 0),
	}
}

type Msg struct {
	Head    string
	Content string
	Time    time.Time
}

// 每一行后加两个空格，用于Markdown换行
var msgTemplate = `
### 驻{hostname}黑猫{head}  
{content}  
at {time}
`

func (m Msg) dingMarkdown() (ret string) {
	ret = strings.Replace(msgTemplate, "{hostname}", hostname, -1)
	ret = strings.Replace(ret, "{head}", m.Head, -1)
	ret = strings.Replace(ret, "{content}", markdownNewlineSymbol(m.Content), -1)
	ret = strings.Replace(ret, "{time}", m.Time.Format("01月02日15:04:05"), -1)
	return
}

func (m Msg) wxContent() (ret string) {
	ret = strings.Replace(msgTemplate, "{hostname}", hostname, -1)
	ret = strings.Replace(ret, "{head}", m.Head, -1)
	ret = strings.Replace(ret, "{content}", m.Content, -1)
	ret = strings.Replace(ret, "{time}", m.Time.Format("01月02日15:04:05"), -1)
	return
}

func markdownNewlineSymbol(origin string) string {
	return strings.Replace(origin, "\n", "\n  ", -1)
}

var msgContext = newMsgContext()

func RunAlterMsgSender() {
	go func() {
		t := time.NewTicker(5 * time.Second)
		for {
			<-t.C
			msgs := msgContext.PopOut()
			if len(msgs) > 0 {
				sendAlterMsg(msgs)
			}
		}
	}()
}

func AddAlertMsg(head, content string) error {
	msgContext.Add(Msg{Head: head, Content: content, Time: time.Now()})
	return nil
}

func sendAlterMsg(msgs Msgs) error {
	wxErr := sendWxAlterMsg(msgs)
	dingErr := sendDingAlterMsg(msgs)

	if wxErr != nil || dingErr != nil {
		wxErrStr := ""
		if wxErr != nil {
			wxErrStr = wxErr.Error()
		}
		dingErrStr := ""
		if dingErr != nil {
			dingErrStr = dingErr.Error()
		}
		return errors.New(fmt.Sprintf("wxError: %s, dingErr: %s", wxErrStr, dingErrStr))
	}
	return nil
}

func sendWxAlterMsg(msgs Msgs) error {
	if qywxToken == "" {
		return nil
	}

	token := strings.Split(qywxToken, "/")
	_, err := go_utils.SendWxQyMsg(token[0], token[2], token[1], msgs.wxContent())
	return err
}

type j map[string]interface{}

func sendDingAlterMsg(msgs Msgs) error {
	if dingAccessToken == "" {
		return nil
	}
	msg := j{
		"msgtype": "markdown",
		"markdown": j{
			"title": msgs.firstHead(),
			"text":  msgs.dingMarkdown(),
		},
	}
	bytes, err := go_utils.HttpPost("https://oapi.dingtalk.com/robot/send?access_token="+dingAccessToken, msg)
	if err != nil {
		panic(err)
	}
	fmt.Println(string(bytes))
	return nil
}
