package web

import (
	"github.com/GoCollaborate/src/constants"
	"html/template"
	"net/http"
	"path/filepath"
	"regexp"
)

type DataModel struct {
	Title   string
	Header1 string
	Header2 string
	Header3 string
	Body    []string
}

var templates = template.Must(template.ParseFiles(makePath("index.html")))
var validPath = regexp.MustCompile("^/(index)*$")

func Index(w http.ResponseWriter, r *http.Request) {
	m := validPath.FindStringSubmatch(r.URL.Path)
	if m == nil {
		http.NotFound(w, r)
		return
	}
	renderTemplate(w, "index")
	return
}

func renderTemplate(w http.ResponseWriter, tmpl string) {
	err := templates.ExecuteTemplate(w, tmpl+".html", 0)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func makePath(path string) string {
	return filepath.Join(constants.LIB_UNIX_DIR+"web", "templates", path)
}
