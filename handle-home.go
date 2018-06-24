package main

import (
	"github.com/bingoohuang/go-utils"
	"net/http"
	"strings"
)

func HandleHome(w http.ResponseWriter, r *http.Request) {
	go_utils.HeadContentTypeHtml(w)

	indexHtml := string(MustAsset("res/index.html"))

	html := go_utils.MinifyHtml(indexHtml, devMode)

	mergeCss := go_utils.MergeCss(MustAsset, go_utils.FilterAssetNames(AssetNames(), ".css"))
	css := go_utils.MinifyCss(mergeCss, devMode)
	mergeScripts := go_utils.MergeJs(MustAsset, go_utils.FilterAssetNames(AssetNames(), ".js"))
	js := go_utils.MinifyJs(mergeScripts, devMode)
	html = strings.Replace(html, "/*.CSS*/", css, 1)
	html = strings.Replace(html, "/*.SCRIPT*/", js, 1)
	html = strings.Replace(html, "${contextPath}", contextPath, -1)

	w.Write([]byte(html))
}
