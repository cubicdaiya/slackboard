package slackboard

import (
	"bytes"
	"fmt"
	"html/template"
	"io"
	"mime"
	"net/http"
	"path"
)

type Topic struct {
	Tag   SectionTag
	Count uint64
}

func SetupUI() {
	index := "ui/index.html"
	bs, err := Asset(index)
	if err != nil {
		LogError.Warn(fmt.Sprintf("%s is not found", index))
		return
	}

	IndexTemplate = template.New("index")
	IndexTemplate, err = IndexTemplate.Parse(string(bs))
	if err != nil {
		LogError.Warn(fmt.Sprintf("template: %s could not be parsed", index))
		return
	}

	http.HandleFunc("/ui/", UIHandler)
}

func UIHandler(w http.ResponseWriter, r *http.Request) {
	LogAcceptedRequest(r, "")
	w.Header().Set("Server", fmt.Sprintf("slackboard/%s", Version))

	p := r.URL.Path
	switch p {
	case "/ui/":
		fallthrough
	case "/ui/index.html":
		w.Header().Set("Content-Type", "text/html")
		IndexTemplate.Execute(w, Topics)
	default:
		bs, err := Asset(p[1:])
		if err != nil {
			http.Error(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)
			return
		}
		ctype := mime.TypeByExtension(path.Ext(p))
		if ctype != "" {
			w.Header().Set("Content-Type", ctype)
		} else {
			w.Header().Set("Content-Type", "text/plain")
		}
		io.Copy(w, bytes.NewBuffer(bs))
	}

}
