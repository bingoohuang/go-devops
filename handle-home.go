package main

import (
	"github.com/bingoohuang/go-utils"
	"net/http"
	"strings"
)

func HandleHome(w http.ResponseWriter, r *http.Request) {
	go_utils.HeadContentTypeHtml(w)

	indexHtml := string(MustAsset("res/index.html"))

	html := go_utils.MinifyHtml(indexHtml, appConfig.DevMode)

	mergeCss := go_utils.MergeCss(MustAsset, FilterAssetNamesExcluded(AssetNames(), ".css",
		"res/codemirror.min", "res/jquery-confirm.min"))
	css := go_utils.MinifyCss(mergeCss, appConfig.DevMode)
	mergeScripts := go_utils.MergeJs(MustAsset, FilterAssetNamesExcluded(AssetNames(), ".js",
		"res/jquery.min", "res/codemirror.min", "res/toml.min", "res/jquery.contextMenu.min", "res/jquery-confirm.min"))
	js := go_utils.MinifyJs(mergeScripts, appConfig.DevMode)
	html = strings.Replace(html, "/*.CSS*/", css, 1)
	html = strings.Replace(html, "/*.SCRIPT*/", js, 1)
	html = strings.Replace(html, "${contextPath}", appConfig.ContextPath, -1)

	_, _ = w.Write([]byte(html))
}
