package main

import (
	"github.com/bingoohuang/go-utils"
	"regexp"
	"strings"
	"unicode"
)

type ExLog struct {
	Properties     map[string]string
	ExceptionNames string
	Logger         string
	Normal         string
	Context        string
	Err            string
}

type ExLogTailer struct {
	// config
	ExLogChan  chan<- ExLog
	Normal     *regexp.Regexp
	Exception  *regexp.Regexp
	Properties map[string]*regexp.Regexp
	Ignores    []string

	Logger string

	// temp data
	Previous  *go_utils.FifoQueue
	Following []string
}

type ExLogTailerConf struct {
	DirectRegex    bool
	NormalRegex    string
	ExceptionRegex string
	Ignores        string
	Logger         string
	LogFileName    string
	Properties     map[string]string
}

func NewExLogTailer(exLogChan chan<- ExLog, conf *ExLogTailerConf) (*ExLogTailer, error) {
	var normal *regexp.Regexp
	var err error

	if conf.DirectRegex {
		normal, err = regexp.Compile(conf.NormalRegex)
	} else {
		expr := ""
		for _, char := range conf.NormalRegex {
			if unicode.IsDigit(char) {
				expr += `\d`
			} else {
				expr += `\` + string(char)
			}
		}
		normal, err = regexp.Compile(expr)
	}
	if err != nil {
		return nil, err
	}

	var exception *regexp.Regexp
	exception, err = regexp.Compile(conf.ExceptionRegex)
	if err != nil {
		return nil, err
	}

	prop := make(map[string]*regexp.Regexp)
	if conf.Properties != nil {
		var propRegex *regexp.Regexp
		for k, v := range conf.Properties {
			propRegex, err = regexp.Compile(v)
			if err != nil {
				return nil, err
			}

			prop[k] = propRegex
		}
	}

	ignoreArr := make([]string, 0)
	if conf.Ignores != "" {
		ignoreItems := strings.Split(conf.Ignores, ",")
		for _, ignoreItem := range ignoreItems {
			ignoreItem = strings.TrimSpace(ignoreItem)
			if ignoreItem != "" {
				ignoreArr = append(ignoreArr, ignoreItem)
			}
		}
	}

	return &ExLogTailer{
		ExLogChan:  exLogChan,
		Normal:     normal,
		Exception:  exception,
		Properties: prop,
		Ignores:    ignoreArr,
		Logger:     conf.Logger,

		Previous:  go_utils.NewFifoQueue(80),
		Following: make([]string, 0),
	}, nil
}

func (t *ExLogTailer) Loop() {
	if !t.Previous.Empty() && len(t.Following) > 0 {
		t.evictEx()
		t.resetTailer()
	}
}
func (t *ExLogTailer) Line(line string) {
	if t.Normal.MatchString(line) {
		t.Loop()
		t.Previous.Append(line)
	} else if !t.Previous.Empty() {
		t.Following = append(t.Following, line)
	}
}

func (t *ExLogTailer) Error(err error) {
	t.ExLogChan <- ExLog{
		Err: err.Error(),
	}
}

func (t *ExLogTailer) resetTailer() {
	t.Previous = go_utils.NewFifoQueue(t.Previous.Capacity())
	t.Following = make([]string, 0)
}

func (t *ExLogTailer) evictEx() {
	exceptionNames := t.createExceptionNames()
	if t.isIgnored(exceptionNames) {
		return
	}

	pop := t.Previous.Pop().(string)
	normal := t.Normal.FindString(pop)
	properties := t.createProperties(pop)
	context := t.createContext(pop)

	t.ExLogChan <- ExLog{
		Properties:     properties,
		ExceptionNames: exceptionNames,
		Context:        context,
		Normal:         normal,
		Logger:         t.Logger,
	}
}
func (t *ExLogTailer) isIgnored(exceptionNames string) bool {
	for _, ignore := range t.Ignores {
		if strings.Index(exceptionNames, ignore) >= 0 {
			return true
		}
	}

	return false
}

func (t *ExLogTailer) createExceptionNames() string {
	exceptionNames := ""
	for _, l := range t.Following {
		if t.Exception.MatchString(l) {
			exceptionNames += l
		}
	}
	return exceptionNames
}

func (t *ExLogTailer) createProperties(pop string) map[string]string {
	properties := make(map[string]string)
	for k, v := range t.Properties {
		sub := v.FindStringSubmatch(pop)
		if sub != nil && len(sub) >= 2 {
			properties[k] = sub[1]
		}
	}
	return properties
}

func (t *ExLogTailer) createContext(pop string) string {
	context := ""
	for {
		l := t.Previous.Shift()
		if l == nil {
			break
		}

		context += l.(string)
	}
	context += pop

	for _, l := range t.Following {
		context += l
	}

	return context
}
