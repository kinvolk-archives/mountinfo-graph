package main

import (
	"./mijson"
	"io"
	"net/http"
)

func loadPage(w http.ResponseWriter, r *http.Request) {
	j := mijson.GetJson()
	s := string(j)
	io.WriteString(w, s)
}

func main() {
	http.HandleFunc("/", loadPage)
	http.ListenAndServe(":8000", nil)
}
