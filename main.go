package main

import (
	"./mijson"
	"html/template"
	"net/http"
)

func indexHandler(w http.ResponseWriter, r *http.Request) {
	t, err := template.ParseFiles("index.html")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	body := "Paste the contents of your mountinfo file below:"
	err = t.Execute(w, body)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func showHandler(w http.ResponseWriter, r *http.Request) {
	usrInput := r.FormValue("mountinfofile")
	j, err := mijson.GetJson(usrInput)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	t, err := template.ParseFiles("show.html")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	err = t.Execute(w, string(j))
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func main() {
	http.HandleFunc("/show", showHandler)
	http.HandleFunc("/", indexHandler)
	http.ListenAndServe(":8000", nil)
	//TODO: handle errors
}
