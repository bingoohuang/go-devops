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

	mergeCss := go_utils.MergeCss(MustAsset, go_utils.FilterAssetNames(AssetNames(), ".css"))
	css := go_utils.MinifyCss(mergeCss, appConfig.DevMode)
	mergeScripts := go_utils.MergeJs(MustAsset, go_utils.FilterAssetNames(AssetNames(), ".js"))
	js := go_utils.MinifyJs(mergeScripts, appConfig.DevMode)
	html = strings.Replace(html, "/*.CSS*/", css, 1)
	html = strings.Replace(html, "/*.SCRIPT*/", js, 1)
	html = strings.Replace(html, "${contextPath}", appConfig.ContextPath, -1)

	_, _ = w.Write([]byte(html))
}
