package main

import (
	"github.com/bingoohuang/go-utils"
	"regexp"
	"strings"
	"time"
)

type MessageTargetConf struct {
	MessagingType string // like dingtalk-robot/qywx-app
	Properties    string
}

type Messaging interface {
	sendMessage(title, message string)
}

type DingtalkRobotMessaging struct {
	accessToken string
}

func NewDingtalkRobotMessaging(parameters string) *DingtalkRobotMessaging {
	return &DingtalkRobotMessaging{
		accessToken: parameters,
	}
}

var hrefReg = regexp.MustCompile(`<a.+href="(.*?)">(.*?)</a>`)

// 将文本中的<a href="www.example.com">Text</a>替换文markdown文法, [Text](www.example.com)
func fixMarkdown(origin string) string {
	markdown := hrefReg.ReplaceAllString(origin, `[$2]($1)`)

	// 将文本中的换行符转换为markdown换行符，即将"\n"替换为"\n\n"
	return strings.Replace(markdown, "\n", "\n\n", -1)
}

func (t *DingtalkRobotMessaging) sendMessage(title, message string) {
	msg := map[string]interface{}{
		"msgtype": "markdown",
		"markdown": map[string]interface{}{
			"title": title,
			"text":  fixMarkdown(message),
		},
	}

	_, _ = go_utils.HttpPost("https://oapi.dingtalk.com/robot/send?access_token="+t.accessToken, msg)
}

type QywxAppMessaging struct {
	CorpId     string
	CorpSecret string
	AgentId    string
}

func NewQywxAppMessaging(parameters string) *QywxAppMessaging {
	token := strings.Split(parameters, "/")

	return &QywxAppMessaging{CorpId: token[0], CorpSecret: token[2], AgentId: token[1]}
}

func (t *QywxAppMessaging) sendMessage(title, message string) {
	_, _ = go_utils.SendWxQyMsg(t.CorpId, t.CorpSecret, t.AgentId, message)
}

type MessageItem struct {
	Head    string
	Content string
	Time    time.Time
}

type MessagingQueue struct {
	messaging Messaging
	itemsChan chan MessageItem
}

// 每一行后加两个空格，用于Markdown换行
var msgTemplate = `
### 驻{hostname}黑猫{head}  
{content}  
at {time}
`

func (m MessageItem) makeMessageContent() string {
	ret := strings.Replace(msgTemplate, "{hostname}", hostname, -1)
	ret = strings.Replace(ret, "{head}", m.Head, -1)
	ret = strings.Replace(ret, "{content}", m.Content, -1)
	ret = strings.Replace(ret, "{time}", m.Time.Format("01月02日15:04:05"), -1)
	return ret
}

var messageTargetMap = make(map[string]*MessagingQueue)

var quiteAllMessageTarget = make(chan bool)
var quiteAllMessageTargetNum int

func InitMessageTargets() {
	if quiteAllMessageTargetNum > 0 {
		for i := 0; i < quiteAllMessageTargetNum; i++ {
			quiteAllMessageTarget <- true
		}
	}

	quiteAllMessageTargetNum = len(devopsConf.MessageTargets)
	for k, v := range devopsConf.MessageTargets {
		CreateMessageTarget(k, v)
	}
}

func CreateMessageTarget(name string, messageTargetConf MessageTargetConf) {
	var queue *MessagingQueue
	switch messageTargetConf.MessagingType {
	case "dingtalk-robot":
		queue = &MessagingQueue{
			messaging: NewDingtalkRobotMessaging(messageTargetConf.Properties),
		}
	case "qywx-app":
		queue = &MessagingQueue{
			messaging: NewQywxAppMessaging(messageTargetConf.Properties),
		}
	}

	if queue == nil {
		return
	}

	queue.itemsChan = make(chan MessageItem, 100)
	messageTargetMap[name] = queue

	go func() {
		t := time.NewTicker(10 * time.Second)
		messageTitle := ""
		messageContent := ""
		messageCount := 0
		for {
			select {
			case <-t.C:
				if messageCount > 0 {
					queue.messaging.sendMessage(messageTitle, messageContent)
					messageTitle = ""
					messageContent = ""
					messageCount = 0
				}
			case msg := <-queue.itemsChan:
				if messageCount == 0 {
					messageTitle = msg.Head

				}
				messageCount++
				messageContent += msg.makeMessageContent()
			case <-quiteAllMessageTarget:
				return
			}
		}
	}()
}

func AddAlertMsg(messageTargets []string, head, content string) {
	for _, messageTarget := range messageTargets {
		queue, ok := messageTargetMap[messageTarget]
		if ok {
			queue.itemsChan <- MessageItem{Head: head, Content: content, Time: time.Now()}
		}
	}
}
