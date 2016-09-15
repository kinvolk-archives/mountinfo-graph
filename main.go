package main

import (
	"github.com/kinvolk/mountinfo-graph/bindata"
	"github.com/kinvolk/mountinfo-graph/migraph"
	"html/template"
	"log"
	"net/http"
)

func generateFromTemplate(htmlTemplate string, w http.ResponseWriter, body string) {
	m := map[string]func() string{
		"index": bindata.Index,
		"show":  bindata.Show,
	}
	t, err := template.New(htmlTemplate).Parse(m[htmlTemplate]())
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	err = t.Execute(w, body)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func indexHandler(w http.ResponseWriter, r *http.Request) {
	generateFromTemplate("index", w, "Paste the contents of your mountinfo file below:")
}

func showHandler(w http.ResponseWriter, r *http.Request) {
	usrInput := r.FormValue("mountinfofile")
	j, err := migraph.GenerateJSON(usrInput)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	generateFromTemplate("show", w, string(j))
}

func main() {
	http.HandleFunc("/show", showHandler)
	http.HandleFunc("/", indexHandler)
	log.Fatal(http.ListenAndServe(":8000", nil))
}
