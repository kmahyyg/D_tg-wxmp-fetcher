package web

import (
	"context"
	"html/template"
	"net/http"
	"strings"

	"bitbucket.org/mutongx/go-utils/log"
	"bitbucket.org/mutze5/wxfetcher/db"
)

var (
	previewTemplate = template.Must(template.New("redirect").Parse(strings.NewReplacer("\t", "", "\n", "").Replace(`
	<!DOCTYPE html>
	<html>
		<head>
			<meta content="text/html;charset=UTF-8" http-equiv="content-type">
			<meta content="{{.Title}}" property="og:title">
			<meta content="{{.Image}}" property="og:image">
			<meta content="{{.Brief}}" property="og:description">
			<meta content="{{.Author}}" property="og:site_name">
			<meta content="Powered by Telegram Bot @wxmpbot" property="tg:mutong">
			<title>{{.Title}}</title>
		</head>
		<body>
			<script type="text/javascript">
				window.location.href="{{.Link}}";
			</script>
		</body>
	</html>`)))
)

func redirectHandler(w http.ResponseWriter, r *http.Request) {
	var err error
	defer func() {
		if err != nil {
			log.Error("redirectHandler", "%v", err)
		}
	}()
	// Get article meta from URL Path
	key := strings.Trim(r.URL.Path, "/")
	meta, err := db.GetArticleMeta(context.Background(), key)
	if err != nil {
		http.NotFound(w, r)
		return
	}
	// http.Redirect(w, r, meta.Link, http.StatusMovedPermanently)
	previewTemplate.Execute(w, meta)
}

// Serve creates a new web server
func Serve(listen string) {
	http.HandleFunc("/", redirectHandler)
	http.ListenAndServe(listen, nil)
}
