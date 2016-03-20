package slackboard

import (
	"fmt"
	"html/template"
	"net/http"
	"os"
)

type Topic struct {
	Tag   SectionTag
	Count uint64
}

func SetupUI() {
	root := ConfSlackboard.UI.Root
	index := fmt.Sprintf("%s/index.html", root)
	_, err := os.Stat(index)
	if err != nil {
		LogError.Warn(fmt.Sprintf("%s is not found", index))
		return
	}

	IndexTemplate = template.Must(template.ParseFiles(index))

	http.HandleFunc("/ui", UIHandler)
	cssDir := fmt.Sprintf("%s/css", root)
	jsDir := fmt.Sprintf("%s/js", root)
	http.Handle("/css/", http.StripPrefix("/css/", http.FileServer(http.Dir(cssDir))))
	http.Handle("/js/", http.StripPrefix("/js/", http.FileServer(http.Dir(jsDir))))
}

func UIHandler(w http.ResponseWriter, r *http.Request) {
	LogAcceptedRequest(r, "")
	w.Header().Set("Content-Type", "text/html")
	w.Header().Set("Server", fmt.Sprintf("slackboard/%s", Version))
	IndexTemplate.Execute(w, Topics)
}
