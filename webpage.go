package main

import (
	"io"
	"net/http"
)

func loadPage(w http.ResponseWriter, r *http.Request) {
	io.WriteString(w, "Hello, world")
}

func main() {
	http.HandleFunc("/", loadPage)
	http.ListenAndServe(":8000", nil)
}
