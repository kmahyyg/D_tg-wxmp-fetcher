package web

import (
	"context"
	"html/template"
	"net/http"
	"strings"

	"bitbucket.org/mutze5/wxfetcher/db"
)

var (
	redirectTemplate = template.Must(template.New("redirect").Parse(strings.NewReplacer("\t", "", "\n", "").Replace(`
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
			Redirecting...
			<script type="text/javascript">
				window.location.href="{{.Link}}";
			</script>
		</body>
	</html>`)))
)

func redirectHandler(w http.ResponseWriter, r *http.Request) {
	key := strings.Trim(r.URL.Path, "/")
	article, err := db.GetArticleMeta(context.Background(), key)
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		return
	}
	redirectTemplate.Execute(w, article)
}

// Serve creates a new web server
func Serve(listen string) {
	http.HandleFunc("/", redirectHandler)
	http.ListenAndServe(listen, nil)
}
