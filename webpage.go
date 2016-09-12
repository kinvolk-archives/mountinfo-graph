package main

import (
	"./mijson"
	"html/template"
	"net/http"
)

func indexHandler(w http.ResponseWriter, r *http.Request) {
	j := mijson.GetJson()
	t, _ := template.ParseFiles("index.html")
	t.Execute(w, j)
}

func main() {
	http.HandleFunc("/", indexHandler)
	http.ListenAndServe(":8000", nil)
}
