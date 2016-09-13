package main

import (
	"./mijson"
	"fmt"
	"html/template"
	"net/http"
)

func indexHandler(w http.ResponseWriter, r *http.Request) {
	t, _ := template.ParseFiles("index.html")
	//TODO: handle errors
	body := "Hello world"
	t.Execute(w, body)
	//TODO: handle errors
}

func showHandler(w http.ResponseWriter, r *http.Request) {
	usrInput := r.FormValue("mountinfofile")
	fmt.Println(usrInput)
	j := mijson.GetJson(usrInput)
	t, _ := template.ParseFiles("show.html")
	//TODO: handle errors
	t.Execute(w, string(j))
	//TODO: handle errors
}

func main() {
	http.HandleFunc("/show", showHandler)
	http.HandleFunc("/", indexHandler)
	http.ListenAndServe(":8000", nil)
}
