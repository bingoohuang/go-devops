package main

import (
	"bytes"
	"github.com/bingoohuang/go-utils"
	"net/http"
	"strings"
)

func HandleHome(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")

	indexHtml := string(MustAsset("res/index.html"))

	html := go_utils.MinifyHtml(indexHtml, devMode)

	css := go_utils.MinifyCss(mergeCss(), devMode)
	js := go_utils.MinifyJs(mergeScripts(), devMode)
	html = strings.Replace(html, "/*.CSS*/", css, 1)
	html = strings.Replace(html, "/*.SCRIPT*/", js, 1)
	html = strings.Replace(html, "${contextPath}", contextPath, -1)

	w.Write([]byte(html))
}

func mergeCss() string {
	return mergeStatic("index.css", "jquery.contextMenu.css", "codemirror-5.33.0.min.css")
}

func mergeScripts() string {
	return mergeStatic("jquery-3.2.1.min.js", "jquery.loading.js",
		"index.js", "jquery.contextMenu.js", "jquery.ui.position.js",
		"util.js", "machines.js", "contextMenu.js",
		"locateLog.js", "logSize.js", "tailFLog.js", "processInfo.js",
		"logs.js", "highlight.machine.js",
		"codemirror-5.33.0.min.js", "toml-5.33.0.min.js", "conf.js")
}

func mergeStatic(statics ...string) string {
	var scripts bytes.Buffer
	for _, static := range statics {
		scripts.Write(MustAsset("res/" + static))
		scripts.Write([]byte(";"))
	}

	return scripts.String()
}
